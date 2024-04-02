package session

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

type Session struct {
	ID     string
	Values map[string]interface{}
}

type StoreSession struct {
	Sessions map[string]*Session
	mu       sync.Mutex
}

type SessionNotFound struct{}

func (e *SessionNotFound) Error() string {
	return "session not found"
}

func (s *StoreSession) Get(id string) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if session, ok := s.Sessions[id]; ok {
		return session, nil
	}
	return nil, &SessionNotFound{}
}

func (s *StoreSession) Create() *Session {
	id := GenerateRandomString(128)
	session := &Session{
		ID:     id,
		Values: make(map[string]interface{}),
	}
	s.Sessions[id] = session
	return session
}

func (s *StoreSession) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Sessions, id)
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	var result strings.Builder
	for i := 0; i < length; i++ {
		randomIndex := seededRand.Intn(len(charset))
		result.WriteString(string(charset[randomIndex]))
	}
	return result.String()
}
