package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Session struct {
	Mode      string        `json:"mode"`
	Activity  string        `json:"activity"`
	Started   time.Time     `json:"started"`
	Duration  time.Duration `json:"duration"`
	Completed bool          `json:"completed"`
}

type Data struct {
	Sessions []Session `json:"sessions"`
}

type Storage struct {
	path string
	data Data
}

func New() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dirPath := filepath.Join(homeDir, ".workpulse")
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, err
	}

	filePath := filepath.Join(dirPath, "activities.json")

	s := &Storage{
		path: filePath,
		data: Data{Sessions: []Session{}},
	}

	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return s, nil
}

func (s *Storage) load() error {
	file, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, &s.data)
}

func (s *Storage) save() error {
	file, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, file, 0644)
}

func (s *Storage) SaveSession(mode, activity string, duration time.Duration, started time.Time) error {
	session := Session{
		Mode:      mode,
		Activity:  activity,
		Started:   started,
		Duration:  duration,
		Completed: true,
	}

	s.data.Sessions = append(s.data.Sessions, session)
	return s.save()
}

func (s *Storage) GetTodaySessions() []Session {
	today := time.Now().Truncate(24 * time.Hour)
	var todaySessions []Session

	for _, session := range s.data.Sessions {
		if session.Started.After(today) || session.Started.Equal(today) {
			todaySessions = append(todaySessions, session)
		}
	}

	return todaySessions
}

func (s *Storage) GetAllSessions() []Session {
	result := make([]Session, len(s.data.Sessions))
	copy(result, s.data.Sessions)
	return result
}

func (s *Storage) GetSessionsByDate(date time.Time) []Session {
	startOfDay := date.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var sessions []Session
	for _, session := range s.data.Sessions {
		if session.Started.After(startOfDay) && session.Started.Before(endOfDay) {
			sessions = append(sessions, session)
		}
	}

	return sessions
}

func (s *Storage) GetWeekSessions() []Session {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	startOfWeek := now.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)

	var sessions []Session
	for _, session := range s.data.Sessions {
		if session.Started.After(startOfWeek) || session.Started.Equal(startOfWeek) {
			sessions = append(sessions, session)
		}
	}

	return sessions
}

func (s *Storage) GetMonthSessions() []Session {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var sessions []Session
	for _, session := range s.data.Sessions {
		if session.Started.After(startOfMonth) || session.Started.Equal(startOfMonth) {
			sessions = append(sessions, session)
		}
	}

	return sessions
}

func (s *Storage) ExportCSV(sessions []Session) string {
	var result string
	result = "mode,activity,started,duration_minutes,completed\n"

	for _, session := range sessions {
		activity := session.Activity
		if activity == "" {
			activity = "-"
		}
		durationMins := session.Duration.Minutes()
		completed := "false"
		if session.Completed {
			completed = "true"
		}
		result += session.Mode + "," +
			escapeCSV(activity) + "," +
			session.Started.Format(time.RFC3339) + "," +
			formatFloat(durationMins) + "," +
			completed + "\n"
	}

	return result
}

func escapeCSV(s string) string {
	needsQuotes := false
	for _, c := range s {
		if c == ',' || c == '"' || c == '\n' {
			needsQuotes = true
			break
		}
	}
	if !needsQuotes {
		return s
	}
	escaped := ""
	for _, c := range s {
		if c == '"' {
			escaped += "\"\""
		} else {
			escaped += string(c)
		}
	}
	return "\"" + escaped + "\""
}

func formatFloat(f float64) string {
	if f == float64(int(f)) {
		return fmt.Sprintf("%d", int(f))
	}
	return fmt.Sprintf("%.2f", f)
}

func (s *Storage) GetTodayWorkDuration() time.Duration {
	var total time.Duration
	for _, session := range s.GetTodaySessions() {
		if session.Mode == "work" && session.Completed {
			total += session.Duration
		}
	}
	return total
}

func (s *Storage) GetStreak() int {
	if len(s.data.Sessions) == 0 {
		return 0
	}

	daysWithWork := make(map[string]bool)
	for _, session := range s.data.Sessions {
		if session.Mode == "work" && session.Completed {
			dayKey := session.Started.Format("2006-01-02")
			daysWithWork[dayKey] = true
		}
	}

	if len(daysWithWork) == 0 {
		return 0
	}

	streak := 0
	today := time.Now()

	for i := 0; i < 365; i++ {
		checkDate := today.AddDate(0, 0, -i)
		dayKey := checkDate.Format("2006-01-02")

		if daysWithWork[dayKey] {
			streak++
		} else if i > 0 {
			break
		}
	}

	return streak
}

func (s *Storage) GetLongestStreak() int {
	if len(s.data.Sessions) == 0 {
		return 0
	}

	daysWithWork := make(map[string]bool)
	var minDate, maxDate time.Time

	for _, session := range s.data.Sessions {
		if session.Mode == "work" && session.Completed {
			dayKey := session.Started.Format("2006-01-02")
			daysWithWork[dayKey] = true

			if minDate.IsZero() || session.Started.Before(minDate) {
				minDate = session.Started
			}
			if maxDate.IsZero() || session.Started.After(maxDate) {
				maxDate = session.Started
			}
		}
	}

	if len(daysWithWork) == 0 {
		return 0
	}

	longest := 0
	current := 0

	for d := minDate; !d.After(maxDate.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
		dayKey := d.Format("2006-01-02")
		if daysWithWork[dayKey] {
			current++
			if current > longest {
				longest = current
			}
		} else {
			current = 0
		}
	}

	return longest
}
