package session

import "time"

// Ratio represents a work/break duration pair
type Ratio struct {
	Work  time.Duration
	Break time.Duration
}

// Preset ratios
var (
	RatioStandard   = Ratio{Work: 25 * time.Minute, Break: 5 * time.Minute}
	RatioShortBurst = Ratio{Work: 15 * time.Minute, Break: 10 * time.Minute}
	RatioDeepWork   = Ratio{Work: 45 * time.Minute, Break: 15 * time.Minute}
)

// PresetName returns the name of a preset ratio
func PresetName(r Ratio) string {
	switch {
	case r == RatioStandard:
		return "standard"
	case r == RatioShortBurst:
		return "short-burst"
	case r == RatioDeepWork:
		return "deep-work"
	default:
		return "custom"
	}
}

// GetPreset returns a ratio by preset name
func GetPreset(name string) (Ratio, bool) {
	switch name {
	case "standard":
		return RatioStandard, true
	case "short-burst":
		return RatioShortBurst, true
	case "deep-work":
		return RatioDeepWork, true
	default:
		return Ratio{}, false
	}
}

// AvailablePresets returns all available preset names
func AvailablePresets() []string {
	return []string{"standard", "short-burst", "deep-work"}
}

// CycleDuration returns the total duration of one work+break cycle
func (r Ratio) CycleDuration() time.Duration {
	return r.Work + r.Break
}

// Plan represents a planned session with total goal and progress
type Plan struct {
	TotalGoal    time.Duration
	Ratio        Ratio
	TotalCycles  int
	CurrentCycle int
	ElapsedTotal time.Duration
	Active       bool
}

// NewPlan creates a new session plan
func NewPlan(totalGoal time.Duration, ratio Ratio) *Plan {
	cycleDuration := ratio.CycleDuration()
	totalCycles := int(totalGoal / cycleDuration)
	if totalCycles < 1 {
		totalCycles = 1
	}

	return &Plan{
		TotalGoal:    totalGoal,
		Ratio:        ratio,
		TotalCycles:  totalCycles,
		CurrentCycle: 1,
		ElapsedTotal: 0,
		Active:       true,
	}
}

// Start activates the session plan
func (p Plan) Start() Plan {
	return Plan{
		TotalGoal:    p.TotalGoal,
		Ratio:        p.Ratio,
		TotalCycles:  p.TotalCycles,
		CurrentCycle: p.CurrentCycle,
		ElapsedTotal: p.ElapsedTotal,
		Active:       true,
	}
}

// AddElapsed adds elapsed time to the session
func (p Plan) AddElapsed(d time.Duration) Plan {
	return Plan{
		TotalGoal:    p.TotalGoal,
		Ratio:        p.Ratio,
		TotalCycles:  p.TotalCycles,
		CurrentCycle: p.CurrentCycle,
		ElapsedTotal: p.ElapsedTotal + d,
		Active:       p.Active,
	}
}

// NextCycle advances to the next cycle
func (p Plan) NextCycle() Plan {
	return Plan{
		TotalGoal:    p.TotalGoal,
		Ratio:        p.Ratio,
		TotalCycles:  p.TotalCycles,
		CurrentCycle: p.CurrentCycle + 1,
		ElapsedTotal: p.ElapsedTotal,
		Active:       p.Active,
	}
}

// Cancel deactivates the session plan
func (p Plan) Cancel() Plan {
	return Plan{
		TotalGoal:    p.TotalGoal,
		Ratio:        p.Ratio,
		TotalCycles:  p.TotalCycles,
		CurrentCycle: p.CurrentCycle,
		ElapsedTotal: p.ElapsedTotal,
		Active:       false,
	}
}

// IsComplete returns true if all cycles are done
func (p Plan) IsComplete() bool {
	return p.CurrentCycle > p.TotalCycles
}

// Progress returns the overall session progress (0.0 to 1.0)
func (p Plan) Progress() float64 {
	if p.TotalGoal == 0 {
		return 0
	}
	progress := float64(p.ElapsedTotal) / float64(p.TotalGoal)
	if progress > 1 {
		return 1
	}
	return progress
}

// RemainingTime returns the estimated remaining time
func (p Plan) RemainingTime() time.Duration {
	remaining := p.TotalGoal - p.ElapsedTotal
	if remaining < 0 {
		return 0
	}
	return remaining
}
