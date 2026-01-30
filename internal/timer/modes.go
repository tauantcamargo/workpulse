package timer

import "time"

type ModeType string

const (
	ModeWork       ModeType = "work"
	ModeShortBreak ModeType = "short_break"
	ModeLongBreak  ModeType = "long_break"
	ModeWalk       ModeType = "walk"
	ModeWater      ModeType = "water"
	ModeVideo      ModeType = "video"
)

type Mode struct {
	Type     ModeType
	Name     string
	Duration time.Duration
	Emoji    string
	Message  string
}

var DefaultModes = map[ModeType]Mode{
	ModeWork: {
		Type:     ModeWork,
		Name:     "WORK",
		Duration: 25 * time.Minute,
		Emoji:    "🎯",
		Message:  "Stay focused, you got this!",
	},
	ModeShortBreak: {
		Type:     ModeShortBreak,
		Name:     "SHORT BREAK",
		Duration: 5 * time.Minute,
		Emoji:    "☕",
		Message:  "Take a breather, you've earned it!",
	},
	ModeLongBreak: {
		Type:     ModeLongBreak,
		Name:     "LONG BREAK",
		Duration: 15 * time.Minute,
		Emoji:    "🌴",
		Message:  "Time to relax and recharge!",
	},
	ModeWalk: {
		Type:     ModeWalk,
		Name:     "WALK",
		Duration: 10 * time.Minute,
		Emoji:    "🚶",
		Message:  "Get moving, stretch those legs!",
	},
	ModeWater: {
		Type:     ModeWater,
		Name:     "WATER",
		Duration: 2 * time.Minute,
		Emoji:    "💧",
		Message:  "Stay hydrated!",
	},
	ModeVideo: {
		Type:     ModeVideo,
		Name:     "VIDEO",
		Duration: 20 * time.Minute,
		Emoji:    "🎬",
		Message:  "Enjoy your video break!",
	},
}

var ModeOrder = []ModeType{
	ModeWork,
	ModeShortBreak,
	ModeLongBreak,
	ModeWalk,
	ModeWater,
	ModeVideo,
}

func GetMode(t ModeType) Mode {
	if mode, ok := DefaultModes[t]; ok {
		return mode
	}
	return DefaultModes[ModeWork]
}

func GetModeByIndex(index int) Mode {
	if index < 0 || index >= len(ModeOrder) {
		return DefaultModes[ModeWork]
	}
	return DefaultModes[ModeOrder[index]]
}

func NextMode(current ModeType) Mode {
	switch current {
	case ModeWork:
		return DefaultModes[ModeShortBreak]
	case ModeShortBreak:
		return DefaultModes[ModeWork]
	case ModeLongBreak:
		return DefaultModes[ModeWork]
	default:
		return DefaultModes[ModeWork]
	}
}
