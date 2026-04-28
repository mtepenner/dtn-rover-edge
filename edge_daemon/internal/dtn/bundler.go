package dtn

import (
	"encoding/json"
	"fmt"
	"time"
)

type Bundler struct {
	nodeID   string
	sequence int
}

func NewBundler(nodeID string) *Bundler {
	return &Bundler{nodeID: nodeID}
}

func (bundler *Bundler) TelemetryBundle(payload TelemetryPayload) (Bundle, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return Bundle{}, err
	}
	bundler.sequence++
	return Bundle{
		ID:          fmt.Sprintf("%s-tm-%04d", bundler.nodeID, bundler.sequence),
		Sequence:    bundler.sequence,
		Source:      bundler.nodeID,
		Destination: "earth-mission-control",
		Type:        "telemetry",
		Priority:    "bulk",
		CreatedAt:   time.Now().UTC(),
		TTLMinutes:  720,
		Payload:     body,
	}, nil
}
