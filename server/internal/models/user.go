package models

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

// User 構造体 (インメモリDB用)
type User struct {
	ID          []byte
	Name        string
	DisplayName string
	Credentials []webauthn.Credential
}

// webauthn.User インターフェースを満たすためのメソッド
func (u User) WebAuthnID() []byte                         { return u.ID }
func (u User) WebAuthnName() string                       { return u.Name }
func (u User) WebAuthnDisplayName() string                { return u.DisplayName }
func (u User) WebAuthnIcon() string                       { return "" }
func (u User) WebAuthnCredentials() []webauthn.Credential { return u.Credentials }

// CustomSessionData セッション管理用
type CustomSessionData struct {
	Username        string
	WebAuthnSession webauthn.SessionData
}
