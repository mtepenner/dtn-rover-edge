package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/mtepenner/dtn-rover-edge/edge_daemon/internal/autonomy"
	"github.com/mtepenner/dtn-rover-edge/edge_daemon/internal/dtn"
	"github.com/mtepenner/dtn-rover-edge/edge_daemon/internal/hardware_link"
)

type agentState struct {
	mu          sync.RWMutex
	telemetry   dtn.TelemetryPayload
	navigator   autonomy.State
	window      dtn.WindowStatus
	lastSync    time.Time
	lastCommand string
	lastResult  string
}

func main() {
	storage, err := dtn.NewStorage("data/pending_bundles.json")
	if err != nil {
		log.Fatal(err)
	}

	nodeID := envOrDefault("ROVER_NODE_ID", "titan-rover-7")
	uart := hardware_link.NewUART()
	navigator := autonomy.NewNavigator()
	bundler := dtn.NewBundler(nodeID)
	transceiver := dtn.NewTransceiver(envOrDefault("DEEP_SPACE_LINK_URL", "http://127.0.0.1:8082"))
	state := &agentState{}

	go telemetryLoop(uart, navigator, bundler, storage, state)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(writer http.ResponseWriter, _ *http.Request) {
		respondJSON(writer, http.StatusOK, map[string]any{
			"service":         "edge-daemon",
			"status":          "ok",
			"pending_bundles": len(storage.Pending()),
		})
	})
	mux.HandleFunc("/api/state", func(writer http.ResponseWriter, _ *http.Request) {
		state.mu.RLock()
		payload := map[string]any{
			"telemetry":       state.telemetry,
			"navigator":       state.navigator,
			"window":          state.window,
			"last_sync":       state.lastSync,
			"last_command_id": state.lastCommand,
			"last_result":     state.lastResult,
		}
		state.mu.RUnlock()
		respondJSON(writer, http.StatusOK, payload)
	})
	mux.HandleFunc("/api/bundles", func(writer http.ResponseWriter, _ *http.Request) {
		respondJSON(writer, http.StatusOK, storage.Pending())
	})
	mux.HandleFunc("/api/sync", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		result := syncOnce(nodeID, transceiver, storage, uart, navigator, state)
		respondJSON(writer, http.StatusOK, result)
	})

	log.Println("edge daemon listening on http://127.0.0.1:8081")
	if err := http.ListenAndServe(":8081", withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func telemetryLoop(uart *hardware_link.UART, navigator *autonomy.Navigator, bundler *dtn.Bundler, storage *dtn.Storage, state *agentState) {
	ticker := time.NewTicker(900 * time.Millisecond)
	defer ticker.Stop()

	mode := "autonomy-cruise"
	publishTelemetry(uart, navigator, bundler, storage, state, mode)
	for range ticker.C {
		publishTelemetry(uart, navigator, bundler, storage, state, mode)
	}
}

func publishTelemetry(uart *hardware_link.UART, navigator *autonomy.Navigator, bundler *dtn.Bundler, storage *dtn.Storage, state *agentState, mode string) {
	sample := uart.NextTelemetry(mode)
	navState := navigator.Evaluate(sample)
	bundle, err := bundler.TelemetryBundle(sample)
	if err == nil {
		_ = storage.Add(bundle)
	}
	state.mu.Lock()
	state.telemetry = sample
	state.navigator = navState
	state.lastCommand = uart.LastCommandID()
	state.mu.Unlock()
}

func syncOnce(nodeID string, transceiver *dtn.Transceiver, storage *dtn.Storage, uart *hardware_link.UART, navigator *autonomy.Navigator, state *agentState) map[string]any {
	window, windowErr := transceiver.Window()
	acceptedIDs := make([]string, 0)
	commands := make([]dtn.CommandPayload, 0)

	if windowErr == nil {
		if window.Active {
			acceptedIDs, _ = transceiver.SendDownlink(nodeID, storage.Pending())
			_ = storage.Remove(acceptedIDs)
			commands, _ = transceiver.PollRover(nodeID)
			for _, command := range commands {
				uart.ApplyCommand(command)
				navigator.ApplyCommand(command)
			}
		}
		state.mu.Lock()
		state.window = window
		state.lastSync = time.Now().UTC()
		state.lastResult = "sync-complete"
		if len(commands) > 0 {
			state.lastCommand = commands[len(commands)-1].ID
		}
		state.mu.Unlock()
	}

	return map[string]any{
		"window":               window,
		"delivered_bundle_ids": acceptedIDs,
		"received_commands":    commands,
		"pending_after_sync":   len(storage.Pending()),
		"window_error":         errorString(windowErr),
	}
}

func respondJSON(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(payload)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func errorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
