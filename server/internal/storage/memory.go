package storage

import (
	"fmt"
	"sync"

	"passkey_sandbox/internal/models"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

// InMemoryStore インメモリデータベース
type InMemoryStore struct {
	users    map[string]*models.User
	sessions map[string]*models.CustomSessionData
	mutex    *sync.Mutex
}

// NewInMemoryStore 新しいインメモリストアを作成
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		users:    make(map[string]*models.User),
		sessions: make(map[string]*models.CustomSessionData),
		mutex:    &sync.Mutex{},
	}
}

// ユーザー関連のメソッド

// GetUser ユーザーを取得
func (s *InMemoryStore) GetUser(username string) (*models.User, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	user, exists := s.users[username]
	return user, exists
}

// CreateUser 新しいユーザーを作成
func (s *InMemoryStore) CreateUser(username string) *models.User {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user := &models.User{
		ID:          []byte(uuid.NewString()),
		Name:        username,
		DisplayName: username,
		Credentials: []webauthn.Credential{},
	}
	s.users[username] = user
	return user
}

// UpdateUserCredentials ユーザーのクレデンシャルを更新
func (s *InMemoryStore) UpdateUserCredentials(username string, credential webauthn.Credential) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, exists := s.users[username]
	if !exists {
		return fmt.Errorf("user not found: %s", username)
	}

	user.Credentials = append(user.Credentials, credential)
	return nil
}

// セッション関連のメソッド

// CreateSession 新しいセッションを作成
func (s *InMemoryStore) CreateSession(sessionData *models.CustomSessionData) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	sessionID := uuid.NewString()
	s.sessions[sessionID] = sessionData
	return sessionID
}

// GetSession セッションを取得
func (s *InMemoryStore) GetSession(sessionID string) (*models.CustomSessionData, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	session, exists := s.sessions[sessionID]
	return session, exists
}

// DeleteSession セッションを削除
func (s *InMemoryStore) DeleteSession(sessionID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.sessions, sessionID)
}
