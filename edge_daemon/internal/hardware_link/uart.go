package hardware_link

import (
	"math"
	"sync"
	"time"

	"github.com/mtepenner/dtn-rover-edge/edge_daemon/internal/dtn"
)

type UART struct {
	mu        sync.Mutex
	step      int
	position  [2]float64
	target    [2]float64
	battery   float64
	lastCmdID string
}

func NewUART() *UART {
	return &UART{target: [2]float64{12, 3}, battery: 100}
}

func (uart *UART) NextTelemetry(mode string) dtn.TelemetryPayload {
	uart.mu.Lock()
	defer uart.mu.Unlock()

	t := float64(uart.step) * 0.35
	uart.step++
	dx := uart.target[0] - uart.position[0]
	dy := uart.target[1] - uart.position[1]
	uart.position[0] += dx * 0.045
	uart.position[1] += dy * 0.045
	clearance := 0.58 + 0.18*math.Sin(t) + 0.06*math.Cos(t*0.4)
	tilt := 4.0*math.Sin(t*0.6) + 1.2*math.Cos(t*0.3)
	uart.battery = math.Max(68.0, uart.battery-0.08)

	hazard := clearance < 0.31 || math.Abs(tilt) > 10.5
	return dtn.TelemetryPayload{
		Timestamp:  time.Now().UTC(),
		PositionM:  uart.position,
		HeadingDeg: math.Mod(180+28*math.Sin(t*0.3), 360),
		ClearanceM: clearance,
		TiltDeg:    tilt,
		BatteryPct: uart.battery,
		HazardStop: hazard,
		Mode:       mode,
	}
}

func (uart *UART) ApplyCommand(command dtn.CommandPayload) {
	uart.mu.Lock()
	defer uart.mu.Unlock()
	uart.target = command.WaypointM
	uart.lastCmdID = command.ID
}

func (uart *UART) LastCommandID() string {
	uart.mu.Lock()
	defer uart.mu.Unlock()
	return uart.lastCmdID
}
