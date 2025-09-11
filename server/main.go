package main

import (
	"log"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"

	"passkey_sandbox/internal/config"
	"passkey_sandbox/internal/handlers"
	"passkey_sandbox/internal/storage"
)

func main() {
	// WebAuthn設定（設定ファイルから取得）
	// 用途に応じて以下から選択:
	// config.NewWebAuthnConfig()                           // デフォルト設定
	// config.NewWebAuthnConfigWithOptions(config.GetCompatibilityOptions())    // 互換性重視
	webAuthnConfig, err := config.NewWebAuthnConfigWithOptions(config.GetPasskeyOptimizedOptions()) // Passkey最適化
	if err != nil {
		log.Fatalf("Failed to create WebAuthn config: %v", err)
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
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/register/begin", webauthnHandlers.BeginRegistration)
	http.HandleFunc("/register/finish", webauthnHandlers.FinishRegistration)
	http.HandleFunc("/login/begin", webauthnHandlers.BeginLogin)
	http.HandleFunc("/login/finish", webauthnHandlers.FinishLogin)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
