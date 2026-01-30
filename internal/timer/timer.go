package timer

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TickMsg time.Time

func Tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type Timer struct {
	Mode      Mode
	Elapsed   time.Duration
	Running   bool
	Completed bool
}

func New(mode Mode) Timer {
	return Timer{
		Mode:      mode,
		Elapsed:   0,
		Running:   false,
		Completed: false,
	}
}

func (t Timer) Start() Timer {
	return Timer{
		Mode:      t.Mode,
		Elapsed:   t.Elapsed,
		Running:   true,
		Completed: false,
	}
}

func (t Timer) Pause() Timer {
	return Timer{
		Mode:      t.Mode,
		Elapsed:   t.Elapsed,
		Running:   false,
		Completed: t.Completed,
	}
}

func (t Timer) Reset() Timer {
	return Timer{
		Mode:      t.Mode,
		Elapsed:   0,
		Running:   false,
		Completed: false,
	}
}

func (t Timer) Tick() Timer {
	newElapsed := t.Elapsed + time.Second
	completed := newElapsed >= t.Mode.Duration

	return Timer{
		Mode:      t.Mode,
		Elapsed:   newElapsed,
		Running:   !completed && t.Running,
		Completed: completed,
	}
}

func (t Timer) SetMode(mode Mode) Timer {
	return Timer{
		Mode:      mode,
		Elapsed:   0,
		Running:   false,
		Completed: false,
	}
}

func (t Timer) Remaining() time.Duration {
	remaining := t.Mode.Duration - t.Elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (t Timer) Progress() float64 {
	if t.Mode.Duration == 0 {
		return 0
	}
	progress := float64(t.Elapsed) / float64(t.Mode.Duration)
	if progress > 1 {
		return 1
	}
	return progress
}

func FormatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return formatTime(minutes, seconds)
}

func formatTime(minutes, seconds int) string {
	return padZero(minutes) + ":" + padZero(seconds)
}

func padZero(n int) string {
	if n < 10 {
		return "0" + itoa(n)
	}
	return itoa(n)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
