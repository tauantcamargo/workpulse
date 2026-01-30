package storage

import (
	"encoding/json"
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
