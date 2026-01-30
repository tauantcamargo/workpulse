package timer

import "time"

type ModeType string

type DurationConfig struct {
	Work       time.Duration
	ShortBreak time.Duration
	LongBreak  time.Duration
	Walk       time.Duration
	Water      time.Duration
	Video      time.Duration
}

func DefaultDurations() DurationConfig {
	return DurationConfig{
		Work:       25 * time.Minute,
		ShortBreak: 5 * time.Minute,
		LongBreak:  15 * time.Minute,
		Walk:       10 * time.Minute,
		Water:      2 * time.Minute,
		Video:      20 * time.Minute,
	}
}

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

type Modes struct {
	modes map[ModeType]Mode
}

func NewModes(cfg DurationConfig) *Modes {
	return &Modes{
		modes: map[ModeType]Mode{
			ModeWork: {
				Type:     ModeWork,
				Name:     "WORK",
				Duration: cfg.Work,
				Emoji:    "🎯",
				Message:  "Stay focused, you got this!",
			},
			ModeShortBreak: {
				Type:     ModeShortBreak,
				Name:     "SHORT BREAK",
				Duration: cfg.ShortBreak,
				Emoji:    "☕",
				Message:  "Take a breather, you've earned it!",
			},
			ModeLongBreak: {
				Type:     ModeLongBreak,
				Name:     "LONG BREAK",
				Duration: cfg.LongBreak,
				Emoji:    "🌴",
				Message:  "Time to relax and recharge!",
			},
			ModeWalk: {
				Type:     ModeWalk,
				Name:     "WALK",
				Duration: cfg.Walk,
				Emoji:    "🚶",
				Message:  "Get moving, stretch those legs!",
			},
			ModeWater: {
				Type:     ModeWater,
				Name:     "WATER",
				Duration: cfg.Water,
				Emoji:    "💧",
				Message:  "Stay hydrated!",
			},
			ModeVideo: {
				Type:     ModeVideo,
				Name:     "VIDEO",
				Duration: cfg.Video,
				Emoji:    "🎬",
				Message:  "Enjoy your video break!",
			},
		},
	}
}

func (m *Modes) Get(t ModeType) Mode {
	if mode, ok := m.modes[t]; ok {
		return mode
	}
	return m.modes[ModeWork]
}

func (m *Modes) Next(current ModeType) Mode {
	switch current {
	case ModeWork:
		return m.modes[ModeShortBreak]
	case ModeShortBreak:
		return m.modes[ModeWork]
	case ModeLongBreak:
		return m.modes[ModeWork]
	default:
		return m.modes[ModeWork]
	}
}
