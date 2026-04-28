package dtn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Transceiver struct {
	baseURL string
	client  *http.Client
}

func NewTransceiver(baseURL string) *Transceiver {
	return &Transceiver{baseURL: baseURL, client: &http.Client{Timeout: 4 * time.Second}}
}

func (transceiver *Transceiver) Window() (WindowStatus, error) {
	response, err := transceiver.client.Get(transceiver.baseURL + "/window")
	if err != nil {
		return WindowStatus{}, err
	}
	defer response.Body.Close()
	var status WindowStatus
	if err := json.NewDecoder(response.Body).Decode(&status); err != nil {
		return WindowStatus{}, err
	}
	return status, nil
}

func (transceiver *Transceiver) SendDownlink(nodeID string, bundles []Bundle) ([]string, error) {
	if len(bundles) == 0 {
		return nil, nil
	}
	body, err := json.Marshal(map[string]any{"node_id": nodeID, "bundles": bundles})
	if err != nil {
		return nil, err
	}
	response, err := transceiver.client.Post(transceiver.baseURL+"/relay/downlink", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var payload struct {
		Accepted []string `json:"accepted_bundle_ids"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return payload.Accepted, nil
}

func (transceiver *Transceiver) PollRover(nodeID string) ([]CommandPayload, error) {
	endpoint := fmt.Sprintf("%s/poll/rover?node_id=%s", transceiver.baseURL, url.QueryEscape(nodeID))
	response, err := transceiver.client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var payload struct {
		Bundles []Bundle `json:"bundles"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, err
	}
	commands := make([]CommandPayload, 0, len(payload.Bundles))
	for _, bundle := range payload.Bundles {
		var command CommandPayload
		if err := json.Unmarshal(bundle.Payload, &command); err != nil {
			continue
		}
		commands = append(commands, command)
	}
	return commands, nil
}
