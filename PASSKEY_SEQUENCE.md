# パスキー認証シークエンス図

このドキュメントでは、WebAuthn（パスキー）による認証フローの詳細なシークエンス図と実装関数を示します。

## API エンドポイント

### 登録関連
- **`POST /register/begin`** → `BeginRegistration()`
- **`POST /register/finish`** → `FinishRegistration()`

### 認証関連  
- **`POST /login/begin`** → `BeginLogin()`
- **`POST /login/finish`** → `FinishLogin()`

### その他
- **`GET /`** → `ServeIndex()`

## 登録フロー (Registration Flow)

```
┌──────────┐       ┌─────────────┐       ┌──────────────┐       ┌────────────┐
│  Client  │       │   Server    │       │   Browser    │       │Authenticator│
└─────┬────┘       └──────┬──────┘       └──────┬───────┘       └─────┬──────┘
      │                   │                     │                     │
      │ 1. POST /register/begin                 │                     │
      ├──────────────────→│                     │                     │
      │   {username}      │                     │                     │
      │                   │ 2. Generate Challenge                     │
      │                   ├─────────────────────┤                     │
      │ 3. Challenge & Options                  │                     │
      │←──────────────────┤                     │                     │
      │                   │                     │                     │
      │ 4. navigator.credentials.create()       │                     │
      ├─────────────────────────────────────────→│                     │
      │                   │                     │ 5. Create Credential│
      │                   │                     ├────────────────────→│
      │                   │                     │ 6. Signed Response  │
      │                   │                     │←────────────────────┤
      │ 7. Credential Response                  │                     │
      │←─────────────────────────────────────────┤                     │
      │                   │                     │                     │
      │ 8. POST /register/finish                │                     │
      ├──────────────────→│                     │                     │
      │   Credential      │ 9. Verify & Store   │                     │
      │                   ├─────────────────────┤                     │
      │ 10. Success       │                     │                     │
      │←──────────────────┤                     │                     │
```

### 登録フローの詳細

1. **登録開始リクエスト**: クライアントがユーザー名を含む登録開始リクエストを送信
   - **実装関数**: `BeginRegistration()`
   - リクエストボディの解析とバリデーション

2. **チャレンジ生成**: サーバーがランダムなチャレンジと登録オプションを生成
   - **実装関数**: `webauthn.BeginRegistration()`（WebAuthnライブラリ）
   - ユーザー情報の取得または新規作成: `GetUser()`, `CreateUser()`

3. **チャレンジ送信**: サーバーがクライアントにチャレンジとオプションを返却
   - **実装関数**: `CreateSession()`, `jsonResponse()`
   - セッションIDをヘッダーに設定

4. **クレデンシャル作成要求**: ブラウザが WebAuthn API を使用してクレデンシャル作成を要求
   - **WebAuthn API**: `navigator.credentials.create()`

5. **クレデンシャル作成**: 認証デバイス（生体認証、セキュリティキーなど）がクレデンシャルを作成
   - デバイス側の処理（実装関数なし）

6. **署名済みレスポンス**: 認証デバイスが秘密鍵で署名したレスポンスをブラウザに返却
   - デバイス側の処理（実装関数なし）

7. **クレデンシャルレスポンス**: ブラウザがクライアントにクレデンシャルレスポンスを返却
   - ブラウザ側の処理（実装関数なし）

8. **登録完了リクエスト**: クライアントが署名済みクレデンシャルをサーバーに送信
   - **実装関数**: `FinishRegistration()`
   - セッション検証: `GetSession()`, `DeleteSession()`

9. **検証と保存**: サーバーがクレデンシャルを検証し、ユーザーデータベースに保存
   - **実装関数**: `protocol.ParseCredentialCreationResponse()`, `webauthn.CreateCredential()`
   - クレデンシャル保存: `UpdateUserCredentials()`

10. **成功レスポンス**: サーバーが登録成功をクライアントに通知
    - **実装関数**: `jsonResponse()`

## 認証フロー (Authentication Flow)

```
┌──────────┐       ┌─────────────┐       ┌──────────────┐       ┌────────────┐
│  Client  │       │   Server    │       │   Browser    │       │Authenticator│
└─────┬────┘       └──────┬──────┘       └──────┬───────┘       └─────┬──────┘
      │                   │                     │                     │
      │ 1. POST /login/begin                    │                     │
      ├──────────────────→│                     │                     │
      │   {username}      │                     │                     │
      │                   │ 2. Generate Challenge                     │
      │                   ├─────────────────────┤                     │
      │ 3. Challenge & Options                  │                     │
      │←──────────────────┤                     │                     │
      │                   │                     │                     │
      │ 4. navigator.credentials.get()          │                     │
      ├─────────────────────────────────────────→│                     │
      │                   │                     │ 5. Sign Challenge   │
      │                   │                     ├────────────────────→│
      │                   │                     │ 6. Signed Response  │
      │                   │                     │←────────────────────┤
      │ 7. Authentication Response              │                     │
      │←─────────────────────────────────────────┤                     │
      │                   │                     │                     │
      │ 8. POST /login/finish                   │                     │
      ├──────────────────→│                     │                     │
      │   Assertion       │ 9. Verify Signature │                     │
      │                   ├─────────────────────┤                     │
      │ 10. Success       │                     │                     │
      │←──────────────────┤                     │                     │
```

### 認証フローの詳細

1. **ログイン開始リクエスト**: クライアントがユーザー名を含むログイン開始リクエストを送信
   - **実装関数**: `BeginLogin()`
   - リクエストボディの解析とバリデーション

2. **チャレンジ生成**: サーバーがランダムなチャレンジと認証オプションを生成
   - **実装関数**: `webauthn.BeginLogin()`（WebAuthnライブラリ）
   - ユーザー情報の取得: `GetUser()`

3. **チャレンジ送信**: サーバーがクライアントにチャレンジとオプションを返却
   - **実装関数**: `CreateSession()`, `jsonResponse()`
   - セッションIDをヘッダーに設定

4. **認証要求**: ブラウザが WebAuthn API を使用して認証を要求
   - **WebAuthn API**: `navigator.credentials.get()`

5. **チャレンジ署名**: 認証デバイスが保存された秘密鍵でチャレンジに署名
   - デバイス側の処理（実装関数なし）

6. **署名済みレスポンス**: 認証デバイスが署名済みアサーションをブラウザに返却
   - デバイス側の処理（実装関数なし）

7. **認証レスポンス**: ブラウザがクライアントに認証レスポンスを返却
   - ブラウザ側の処理（実装関数なし）

8. **認証完了リクエスト**: クライアントが署名済みアサーションをサーバーに送信
   - **実装関数**: `FinishLogin()`
   - セッション検証: `GetSession()`, `DeleteSession()`

9. **署名検証**: サーバーが保存された公開鍵を使用して署名を検証
   - **実装関数**: `protocol.ParseCredentialRequestResponse()`, `webauthn.ValidateLogin()`
   - ユーザー情報の取得: `GetUser()`

10. **成功レスポンス**: サーバーが認証成功をクライアントに通知
    - **実装関数**: `jsonResponse()`

## 実装関数の詳細

### ハンドラー関数（`internal/handlers/webauthn.go`）

- **`BeginRegistration()`**: 登録フローの開始処理
  - リクエストの解析、ユーザー作成/取得、WebAuthnチャレンジ生成
- **`FinishRegistration()`**: 登録フローの完了処理  
  - セッション検証、クレデンシャル解析・検証、データベース保存
- **`BeginLogin()`**: ログインフローの開始処理
  - リクエストの解析、ユーザー検索、認証チャレンジ生成
- **`FinishLogin()`**: ログインフローの完了処理
  - セッション検証、アサーション解析・検証、署名検証
- **`ServeIndex()`**: HTMLページの提供
- **`jsonResponse()`**: JSON レスポンスの送信（共通ヘルパー）

### ストレージ関数（`internal/storage/memory.go`）

- **`NewInMemoryStore()`**: インメモリストアの初期化
- **`GetUser(username)`**: ユーザー情報の取得
- **`CreateUser(username)`**: 新規ユーザーの作成
- **`UpdateUserCredentials(username, credential)`**: ユーザーのクレデンシャル更新
- **`CreateSession(sessionData)`**: セッションの作成とID生成
- **`GetSession(sessionID)`**: セッション情報の取得
- **`DeleteSession(sessionID)`**: セッションの削除

### モデル関数（`internal/models/user.go`）

- **`WebAuthnID()`**: WebAuthnライブラリ用のユーザーID取得
- **`WebAuthnName()`**: WebAuthnライブラリ用のユーザー名取得
- **`WebAuthnDisplayName()`**: WebAuthnライブラリ用の表示名取得
- **`WebAuthnIcon()`**: WebAuthnライブラリ用のアイコン取得
- **`WebAuthnCredentials()`**: WebAuthnライブラリ用のクレデンシャル一覧取得

### 外部ライブラリ関数

- **`webauthn.BeginRegistration(user)`**: WebAuthn登録チャレンジ生成
- **`webauthn.CreateCredential(user, session, response)`**: クレデンシャル作成・検証
- **`webauthn.BeginLogin(user)`**: WebAuthn認証チャレンジ生成
- **`webauthn.ValidateLogin(user, session, response)`**: ログイン署名検証
- **`protocol.ParseCredentialCreationResponse(r)`**: 登録レスポンス解析
- **`protocol.ParseCredentialRequestResponse(r)`**: 認証レスポンス解析

## セキュリティ特徴

### パスキーの利点

- **フィッシング耐性**: オリジンに紐付けられているため、偽サイトでは動作しない
- **パスワード不要**: 生体認証やPINで安全に認証可能
- **リプレイ攻撃耐性**: チャレンジベースの認証により過去の認証情報は再利用不可
- **クライアントサイド認証**: 秘密鍵はデバイスから外部に送信されない

### 実装における重要なポイント

- **チャレンジの一意性**: 各認証で異なるランダムなチャレンジを生成
- **オリジン検証**: レスポンスが正しいオリジンから来ていることを確認
- **セッション管理**: チャレンジとセッションの適切な管理と検証
- **タイムアウト**: セッションに適切な有効期限を設定

## 技術仕様

- **プロトコル**: WebAuthn (Web Authentication API)
- **暗号方式**: 楕円曲線暗号 (通常 ES256)
- **認証デバイス**: 生体認証、セキュリティキー、TPM など
- **ブラウザサポート**: Chrome, Firefox, Safari, Edge など主要ブラウザ

## 参考リンク

- [WebAuthn仕様](https://www.w3.org/TR/webauthn-2/)
- [FIDO Alliance](https://fidoalliance.org/)
- [MDN WebAuthn API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Authentication_API)
