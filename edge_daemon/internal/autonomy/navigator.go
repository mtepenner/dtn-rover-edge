package autonomy

import (
	"math"

	"github.com/mtepenner/dtn-rover-edge/edge_daemon/internal/dtn"
)

type State struct {
	Mode              string     `json:"mode"`
	TargetWaypointM   [2]float64 `json:"target_waypoint_m"`
	DistanceRemaining float64    `json:"distance_remaining_m"`
	Recommendation    string     `json:"recommendation"`
}

type Navigator struct {
	target [2]float64
}

func NewNavigator() *Navigator {
	return &Navigator{target: [2]float64{12, 3}}
}

func (navigator *Navigator) ApplyCommand(command dtn.CommandPayload) {
	navigator.target = command.WaypointM
}

func (navigator *Navigator) Evaluate(sample dtn.TelemetryPayload) State {
	dx := navigator.target[0] - sample.PositionM[0]
	dy := navigator.target[1] - sample.PositionM[1]
	distance := math.Hypot(dx, dy)
	state := State{
		Mode:              sample.Mode,
		TargetWaypointM:   navigator.target,
		DistanceRemaining: distance,
		Recommendation:    "continue-traverse",
	}
	if sample.HazardStop {
		state.Recommendation = "pause-for-hazard"
	}
	if distance < 1.2 {
		state.Recommendation = "prepare-sample-collection"
	}
	return state
}
