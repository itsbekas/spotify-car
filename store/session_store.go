package store

import (
	"sync"
	"time"
)

type SessionData struct {
	SessionHash  string
	AccessToken  string
	RefreshToken string
	Expiration   time.Time // The token's expiry time
	CreatedAt    time.Time // The token's created, used to clean old tokens
}

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]SessionData
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]SessionData),
	}
}

func (s *SessionStore) SetInitialData(picoID, sessionHash string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[picoID] = SessionData{SessionHash: sessionHash, CreatedAt: time.Now()}
}

func (s *SessionStore) SetTokens(picoID, accessToken, refreshToken string, expiresIn int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data := s.sessions[picoID]
	data.AccessToken = accessToken
	data.RefreshToken = refreshToken
	data.Expiration = time.Now().Add(time.Second * time.Duration(expiresIn))

	s.sessions[picoID] = data
}

func (s *SessionStore) GetTokens(picoID, sessionHash string) (string, string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, ok := s.sessions[picoID]

	if !ok || data.SessionHash != sessionHash || time.Now().After(data.Expiration) {
		return "", "", false
	}

	return data.AccessToken, data.RefreshToken, true
}
