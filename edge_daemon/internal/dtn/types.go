package dtn

import (
	"encoding/json"
	"time"
)

type TelemetryPayload struct {
	Timestamp  time.Time  `json:"timestamp"`
	PositionM  [2]float64 `json:"position_m"`
	HeadingDeg float64    `json:"heading_deg"`
	ClearanceM float64    `json:"clearance_m"`
	TiltDeg    float64    `json:"tilt_deg"`
	BatteryPct float64    `json:"battery_pct"`
	HazardStop bool       `json:"hazard_stop"`
	Mode       string     `json:"mode"`
}

type CommandPayload struct {
	ID         string             `json:"id"`
	Action     string             `json:"action"`
	WaypointM  [2]float64         `json:"waypoint_m"`
	Parameters map[string]float64 `json:"parameters"`
	CreatedAt  time.Time          `json:"created_at"`
}

type Bundle struct {
	ID          string          `json:"id"`
	Sequence    int             `json:"sequence"`
	Source      string          `json:"source"`
	Destination string          `json:"destination"`
	Type        string          `json:"type"`
	Priority    string          `json:"priority"`
	CreatedAt   time.Time       `json:"created_at"`
	TTLMinutes  int             `json:"ttl_minutes"`
	Payload     json.RawMessage `json:"payload"`
}

type WindowStatus struct {
	Active             bool      `json:"active"`
	CurrentTime        time.Time `json:"current_time"`
	NextWindowStart    time.Time `json:"next_window_start"`
	CurrentWindowEnd   time.Time `json:"current_window_end"`
	OneWayDelaySeconds int       `json:"one_way_delay_seconds"`
	PacketLossRate     float64   `json:"packet_loss_rate"`
}
