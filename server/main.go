package main

import (
	"log"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"

	"passkey_sandbox/internal/handlers"
	"passkey_sandbox/internal/storage"
)

func main() {
	// WebAuthn設定
	webAuthnConfig := &webauthn.Config{
		RPDisplayName: "Go WebAuthn Example",             // Relying Party Name
		RPID:          "localhost",                       // Relying Party ID (ドメイン)
		RPOrigins:     []string{"http://localhost:8080"}, // 許可するオリジン
	}

	w, err := webauthn.New(webAuthnConfig)
	if err != nil {
		log.Fatalf("Failed to create WebAuthn from config: %v", err)
	}

	// ストレージとハンドラーを初期化
	store := storage.NewInMemoryStore()
	webauthnHandlers := handlers.NewWebAuthnHandlers(w, store)

	// ルーティング設定
	http.HandleFunc("/", webauthnHandlers.ServeIndex)

	// 静的ファイル（JS、CSS等）を提供
	fs := http.FileServer(http.Dir("templates/"))
	http.Handle("/app.js", http.StripPrefix("/", fs))

	http.HandleFunc("/register/begin", webauthnHandlers.BeginRegistration)
	http.HandleFunc("/register/finish", webauthnHandlers.FinishRegistration)
	http.HandleFunc("/login/begin", webauthnHandlers.BeginLogin)
	http.HandleFunc("/login/finish", webauthnHandlers.FinishLogin)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
