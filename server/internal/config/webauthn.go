package config

import (
	"github.com/go-webauthn/webauthn/metadata/providers/cached"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// WebAuthnConfig WebAuthn設定を生成
func NewWebAuthnConfig() *webauthn.Config {
	return &webauthn.Config{
		RPDisplayName: "Go WebAuthn Example",             // Relying Party Name
		RPID:          "localhost",                       // Relying Party ID (ドメイン)
		RPOrigins:     []string{"http://localhost:8080"}, // 許可するオリジン
		// AuthenticatorSelectionを明示的に設定
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			AuthenticatorAttachment: protocol.AuthenticatorAttachment("platform"), // "platform" or "cross-platform"
			ResidentKey:             protocol.ResidentKeyRequirementPreferred,     // "discouraged", "preferred", "required"
			RequireResidentKey:      protocol.ResidentKeyNotRequired(),            // Boolean pointer
			UserVerification:        protocol.VerificationPreferred,               // "discouraged", "preferred", "required"
		},
	}
}

// WebAuthnConfigOptions 設定オプションの構造体
type WebAuthnConfigOptions struct {
	RPDisplayName           string
	RPID                    string
	RPOrigins               []string
	AuthenticatorAttachment string
	ResidentKey             protocol.ResidentKeyRequirement
	RequireResidentKey      *bool
	UserVerification        protocol.UserVerificationRequirement
	MDSProviderOptions      []cached.Option
}

// NewWebAuthnConfigWithOptions カスタムオプションでWebAuthn設定を生成
func NewWebAuthnConfigWithOptions(opts WebAuthnConfigOptions) (*webauthn.Config, error) {
	config := &webauthn.Config{
		RPDisplayName: opts.RPDisplayName,
		RPID:          opts.RPID,
		RPOrigins:     opts.RPOrigins,
	}

	// AuthenticatorSelectionの設定
	authenticatorSelection := protocol.AuthenticatorSelection{
		ResidentKey:      opts.ResidentKey,
		UserVerification: opts.UserVerification,
	}

	// AuthenticatorAttachmentが指定されている場合のみ設定
	if opts.AuthenticatorAttachment != "" {
		authenticatorSelection.AuthenticatorAttachment = protocol.AuthenticatorAttachment(opts.AuthenticatorAttachment)
	}

	// RequireResidentKeyが指定されている場合のみ設定
	if opts.RequireResidentKey != nil {
		if *opts.RequireResidentKey {
			authenticatorSelection.RequireResidentKey = protocol.ResidentKeyRequired()
		} else {
			authenticatorSelection.RequireResidentKey = protocol.ResidentKeyNotRequired()
		}
	}

	var err error
	if config.MDS, err = cached.New(opts.MDSProviderOptions...); err != nil {
		return nil, err
	}

	config.AuthenticatorSelection = authenticatorSelection
	return config, nil
}

// GetDefaultWebAuthnOptions デフォルトの設定オプションを取得
func GetDefaultWebAuthnOptions() WebAuthnConfigOptions {
	return WebAuthnConfigOptions{
		RPDisplayName:           "Go WebAuthn Example",
		RPID:                    "localhost",
		RPOrigins:               []string{"http://localhost:8080"},
		AuthenticatorAttachment: "platform",
		ResidentKey:             protocol.ResidentKeyRequirementPreferred,
		RequireResidentKey:      protocol.ResidentKeyNotRequired(),
		UserVerification:        protocol.VerificationPreferred,
		MDSProviderOptions: []cached.Option{
			// MDSキャッシュファイルのパスを指定
			cached.WithPath(".cache/fidoalliance_mds.jwt"),
		},
	}
}

// GetPasskeyOptimizedOptions Passkeyに最適化された設定オプションを取得
func GetPasskeyOptimizedOptions() WebAuthnConfigOptions {
	return WebAuthnConfigOptions{
		RPDisplayName:           "Go WebAuthn Example",
		RPID:                    "localhost",
		RPOrigins:               []string{"http://localhost:8080"},
		AuthenticatorAttachment: "platform",
		ResidentKey:             protocol.ResidentKeyRequirementRequired,
		RequireResidentKey:      protocol.ResidentKeyRequired(),
		UserVerification:        protocol.VerificationRequired,
		MDSProviderOptions: []cached.Option{
			cached.WithPath(".cache/fidoalliance_mds.jwt"),
		},
	}
}

// GetCompatibilityOptions 互換性重視の設定オプションを取得
func GetCompatibilityOptions() WebAuthnConfigOptions {
	return WebAuthnConfigOptions{
		RPDisplayName:           "Go WebAuthn Example",
		RPID:                    "localhost",
		RPOrigins:               []string{"http://localhost:8080"},
		AuthenticatorAttachment: "", // 指定なし（どちらでも可）
		ResidentKey:             protocol.ResidentKeyRequirementPreferred,
		RequireResidentKey:      protocol.ResidentKeyNotRequired(),
		UserVerification:        protocol.VerificationPreferred,
		MDSProviderOptions: []cached.Option{
			cached.WithPath(".cache/fidoalliance_mds.jwt"),
		},
	}
}
