# Passkey Sandbox

WebAuthn/Passkeyの動作を確認できるローカル開発環境です。

## 機能

- Passkeyによるユーザー登録
- Passkeyによる認証
- リアルタイムログ出力とスクロール機能
- ログクリア機能
- ホットリロード対応（Air使用）
- シンプルなWebインターフェース
- 詳細なシーケンス図ドキュメント

## 必要な環境

- Go 1.23+
- Docker & Docker Compose （オプション）
- Make （推奨）

## クイックスタート

### 推奨: ローカル開発モード（ホットリロード）

```bash
# Airをインストール（初回のみ）
make install-air

# 開発サーバーを起動
make up
```

### Docker使用

```bash
# Docker Composeで起動
make docker-up

# または直接
docker-compose up --build
```

## 利用可能なコマンド

### 開発関連（推奨）
- `make up` - ローカル開発サーバーを起動（Air使用、ホットリロード）
- `make down` - ローカル開発サーバーを停止
- `make install-air` - Airをインストール

### Docker関連
- `make docker-up` - Docker Composeでアプリケーションを起動
- `make docker-down` - Docker Composeでアプリケーションを停止
- `make build` - Docker Composeでイメージをビルドして起動
- `make logs` - Docker Composeのログを表示
- `make clean` - Dockerイメージとボリュームを削除

## 使用方法

### 1. アプリケーションを起動

```bash
# 推奨: ホットリロード開発モード
make up

# または Docker使用
make docker-up
```

### 2. ブラウザでアクセス

http://localhost:8080 にアクセスしてください。

### 3. Passkeyでユーザー登録

1. ユーザー名を入力
2. 「登録」ボタンをクリック
3. ブラウザのPasskey登録プロンプトに従って操作
4. ログエリアで登録成功メッセージを確認

### 4. Passkeyで認証

1. 登録済みのユーザー名を入力
2. 「認証」ボタンをクリック
3. ブラウザのPasskey認証プロンプトに従って操作
4. ログエリアで認証成功メッセージを確認

### 5. ログの管理

- ログエリアはスクロール可能で、新しいメッセージが自動的に下部に表示されます
- 「ログクリア」ボタンでログをクリアできます
- ブラウザのコンソールでも詳細なログを確認できます

## 注意事項

- このサンプルはローカル開発・検証用です
- HTTPSでない環境のため、実際のブラウザでは一部制限があります
- データはメモリ上にのみ保存され、サーバー再起動時に消去されます
- ブラウザのコンソールでも詳細なログを確認できます

## 開発について

### ホットリロード開発

```bash
# 開発サーバーを起動（推奨）
make up

# ファイル変更を検知して自動再起動
# - Go ファイル (.go)
# - HTML ファイル (.html)  
# - CSS ファイル (.css)
# - JavaScript ファイル (.js)
```

### 従来のローカル開発（Docker不使用）

```bash
cd server
go mod tidy
go run main.go
```

### 依存関係

- Go 1.23+
- github.com/go-webauthn/webauthn v0.11.2
- github.com/air-verse/air@latest（開発時）

### アーキテクチャ

- **Clean Architecture**: 責務を分離したパッケージ構造
- **Handler Layer**: HTTP リクエスト処理
- **Storage Layer**: データ永続化（インメモリ実装）
- **Model Layer**: ドメインモデル定義
- **Config Layer**: WebAuthn設定管理
- **Client-side**: バニラJavaScript（フレームワーク不使用）

## ディレクトリ構造

```
passkey_sandbox/
├── docker-compose.yml              # Docker Compose設定
├── Makefile                        # 開発用コマンド定義
├── README.md                       # このファイル
├── PASSKEY_SEQUENCE.md             # WebAuthnシーケンス図
└── server/
    ├── .air.toml                   # Air設定ファイル
    ├── Dockerfile                  # Dockerビルド設定
    ├── go.mod                      # Go モジュール定義
    ├── go.sum                      # Go 依存関係チェックサム
    ├── main.go                     # サーバーエントリーポイント
    ├── internal/
    │   ├── config/
    │   │   └── webauthn.go         # WebAuthn設定管理
    │   ├── handlers/
    │   │   └── webauthn.go         # WebAuthnハンドラー
    │   ├── models/
    │   │   └── user.go             # ユーザーモデル
    │   └── storage/
    │       └── memory.go           # インメモリストレージ
    └── templates/
        ├── index.html              # メインページ
        ├── app.js                  # クライアントサイドJS
        └── style.css               # スタイルシート
```

## API エンドポイント

- `GET /` - メインページ（ユーザー登録・認証機能）
- `GET /static/app.js` - クライアントサイドJavaScript
- `GET /static/style.css` - スタイルシート
- `POST /register/begin` - Passkey登録開始（チャレンジ生成）
- `POST /register/finish` - Passkey登録完了（クレデンシャル検証・保存）
- `POST /login/begin` - Passkey認証開始（チャレンジ生成）
- `POST /login/finish` - Passkey認証完了（署名検証）

## 実装の詳細

WebAuthnの実装についての詳細なシーケンス図とAPIフローについては、[PASSKEY_SEQUENCE.md](./PASSKEY_SEQUENCE.md)を参照してください。

## トラブルシューティング

### ホットリロードが動作しない場合

```bash
# Airを再インストール
make install-air

# 開発サーバーを再起動
make down
make up
```

### Dockerビルドエラーの場合

```bash
# クリーンビルド
make clean
make build
```

### ポート8080が使用中の場合

```bash
# 使用中のプロセスを確認
lsof -i :8080

# 既存のサーバーを停止
make down
make docker-down
```
