# Passkey Sandbox

WebAuthn/Passkeyの動作を確認できるローカル開発環境です。

## 機能

- Passkeyによるユーザー登録
- Passkeyによるログイン認証
- リアルタイムログ出力
- シンプルなWebインターフェース
- 詳細なシーケンス図ドキュメント

## 必要な環境

- Docker
- Docker Compose

## 使用方法

### 1. Docker Composeでサーバーを起動

```bash
docker-compose up --build
```

### 2. ブラウザでアクセス

http://localhost:8080 にアクセスしてください。

### 3. Passkeyでユーザー登録

1. ユーザー名を入力
2. 「登録 (Register)」ボタンをクリック
3. ブラウザのPasskey登録プロンプトに従って操作
4. 登録成功メッセージを確認

### 4. Passkeyでログイン

1. 登録済みのユーザー名を入力
2. 「ログイン (Login)」ボタンをクリック
3. ブラウザのPasskey認証プロンプトに従って操作
4. ログイン成功メッセージを確認

## 注意事項

- このサンプルはローカル開発・検証用です
- HTTPSでない環境のため、実際のブラウザでは一部制限があります
- データはメモリ上にのみ保存され、サーバー再起動時に消去されます
- ブラウザのコンソールでも詳細なログを確認できます

## 実装の詳細

WebAuthnの実装についての詳細なシーケンス図とAPIフローについては、[PASSKEY_SEQUENCE.md](./PASSKEY_SEQUENCE.md)を参照してください。

## ディレクトリ構造

```
passkey_sandbox/
├── docker-compose.yml
├── README.md
├── PASSKEY_SEQUENCE.md          # WebAuthnシーケンス図
└── server/
    ├── Dockerfile
    ├── go.mod
    ├── go.sum
    ├── main.go                  # サーバーエントリーポイント
    ├── internal/
    │   ├── handlers/
    │   │   └── webauthn.go      # WebAuthnハンドラー
    │   ├── models/
    │   │   └── user.go          # ユーザーモデル
    │   └── storage/
    │       └── memory.go        # インメモリストレージ
    └── templates/
        ├── index.html           # メインページ
        └── app.js               # クライアントサイドJS
```

## API エンドポイント

- `GET /` - メインページ（ユーザー登録・ログイン機能）
- `GET /app.js` - クライアントサイドJavaScript
- `POST /register/begin` - Passkey登録開始（チャレンジ生成）
- `POST /register/finish` - Passkey登録完了（クレデンシャル検証・保存）
- `POST /login/begin` - Passkeyログイン開始（チャレンジ生成）
- `POST /login/finish` - Passkeyログイン完了（署名検証）

## 開発について

### ローカル開発（Docker不使用）

```bash
cd server
go mod tidy
go run main.go
```

### 依存関係

- Go 1.23+
- github.com/go-webauthn/webauthn v0.11.2

### アーキテクチャ

- **Clean Architecture**: 責務を分離したパッケージ構造
- **Handler Layer**: HTTP リクエスト処理
- **Storage Layer**: データ永続化（インメモリ実装）
- **Model Layer**: ドメインモデル定義
- **Client-side**: バニラJavaScript（フレームワーク不使用）
