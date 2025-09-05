package handlers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"

	"passkey_sandbox/internal/models"
	"passkey_sandbox/internal/storage"
)

// WebAuthnHandlers WebAuthn関連のハンドラー
type WebAuthnHandlers struct {
	webauthn *webauthn.WebAuthn
	store    *storage.InMemoryStore
}

// NewWebAuthnHandlers 新しいWebAuthnハンドラーを作成
func NewWebAuthnHandlers(w *webauthn.WebAuthn, store *storage.InMemoryStore) *WebAuthnHandlers {
	return &WebAuthnHandlers{
		webauthn: w,
		store:    store,
	}
}

// ServeIndex インデックスページを提供
func (h *WebAuthnHandlers) ServeIndex(rw http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(rw, nil)
}

// BeginRegistration 登録開始 (チャレンジ生成)
func (h *WebAuthnHandlers) BeginRegistration(rw http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		h.jsonResponse(rw, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, exists := h.store.GetUser(reqBody.Username)
	if !exists {
		user = h.store.CreateUser(reqBody.Username)
	}

	options, sessionData, err := h.webauthn.BeginRegistration(user)
	if err != nil {
		h.jsonResponse(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID := h.store.CreateSession(&models.CustomSessionData{
		Username:        user.Name,
		WebAuthnSession: *sessionData,
	})

	rw.Header().Set("X-Session-ID", sessionID)
	h.jsonResponse(rw, options, http.StatusOK)
}

// FinishRegistration 登録完了 (クレデンシャル検証・保存)
func (h *WebAuthnHandlers) FinishRegistration(rw http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("X-Session-ID")
	if sessionID == "" {
		h.jsonResponse(rw, "Session ID required", http.StatusBadRequest)
		return
	}

	customSession, exists := h.store.GetSession(sessionID)
	if !exists {
		h.jsonResponse(rw, "Session not found", http.StatusBadRequest)
		return
	}
	h.store.DeleteSession(sessionID)

	parsedResponse, err := protocol.ParseCredentialCreationResponse(r)
	if err != nil {
		h.jsonResponse(rw, err.Error(), http.StatusBadRequest)
		return
	}

	user, exists := h.store.GetUser(customSession.Username)
	if !exists {
		h.jsonResponse(rw, "User not found", http.StatusBadRequest)
		return
	}

	credential, err := h.webauthn.CreateCredential(user, customSession.WebAuthnSession, parsedResponse)
	if err != nil {
		h.jsonResponse(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.store.UpdateUserCredentials(customSession.Username, *credential)
	if err != nil {
		h.jsonResponse(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	h.jsonResponse(rw, "Registration Success", http.StatusOK)
}

// BeginLogin ログイン開始 (チャレンジ生成)
func (h *WebAuthnHandlers) BeginLogin(rw http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		h.jsonResponse(rw, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, exists := h.store.GetUser(reqBody.Username)
	if !exists {
		h.jsonResponse(rw, "User not found", http.StatusBadRequest)
		return
	}

	options, sessionData, err := h.webauthn.BeginLogin(user)
	if err != nil {
		h.jsonResponse(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID := h.store.CreateSession(&models.CustomSessionData{
		Username:        user.Name,
		WebAuthnSession: *sessionData,
	})

	rw.Header().Set("X-Session-ID", sessionID)
	h.jsonResponse(rw, options, http.StatusOK)
}

// FinishLogin ログイン完了 (署名検証)
func (h *WebAuthnHandlers) FinishLogin(rw http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("X-Session-ID")
	if sessionID == "" {
		h.jsonResponse(rw, "Session ID required", http.StatusBadRequest)
		return
	}

	customSession, exists := h.store.GetSession(sessionID)
	if !exists {
		h.jsonResponse(rw, "Session not found", http.StatusBadRequest)
		return
	}
	h.store.DeleteSession(sessionID)

	parsedResponse, err := protocol.ParseCredentialRequestResponse(r)
	if err != nil {
		h.jsonResponse(rw, err.Error(), http.StatusBadRequest)
		return
	}

	user, exists := h.store.GetUser(customSession.Username)
	if !exists {
		h.jsonResponse(rw, "User not found", http.StatusBadRequest)
		return
	}

	_, err = h.webauthn.ValidateLogin(user, customSession.WebAuthnSession, parsedResponse)
	if err != nil {
		h.jsonResponse(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	h.jsonResponse(rw, "Login Success", http.StatusOK)
}

// jsonResponse JSON レスポンスを送信するヘルパー関数
func (h *WebAuthnHandlers) jsonResponse(rw http.ResponseWriter, data interface{}, status int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		http.Error(rw, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
		return
	}
	rw.Write(buf.Bytes())
}
