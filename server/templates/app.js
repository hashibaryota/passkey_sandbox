// -----------------------------------------------------------------------------
// Passkey WebAuthn Client-side Implementation
// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
// 1. DOM要素の取得とイベントリスナーの設定
// -----------------------------------------------------------------------------
const usernameInput = document.getElementById('username');
const registerBtn = document.getElementById('registerBtn');
const loginBtn = document.getElementById('loginBtn');
const logOutput = document.getElementById('log');

registerBtn.addEventListener('click', handleRegister);
loginBtn.addEventListener('click', handleLogin);

// ログ出力用ヘルパー
function logMessage(message) {
    console.log(message);
    const now = new Date().toLocaleTimeString();
    logOutput.textContent += `[${now}] ${message}\n\n`;
}

// -----------------------------------------------------------------------------
// 2. Base64URL <-> ArrayBuffer 変換ヘルパー
// これらの変換はWebAuthn APIとサーバー間でデータをやり取りするために必須
// -----------------------------------------------------------------------------

function base64urlToBuffer(base64url) {
    const base64 = base64url.replace(/-/g, '+').replace(/_/g, '/');
    const raw = window.atob(base64);
    const buffer = new Uint8Array(raw.length);
    for (let i = 0; i < raw.length; i++) {
        buffer[i] = raw.charCodeAt(i);
    }
    return buffer;
}

function bufferToBase64url(buffer) {
    const str = String.fromCharCode.apply(null, new Uint8Array(buffer));
    return window.btoa(str).replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '');
}

// -----------------------------------------------------------------------------
// 3. 登録処理 (Registration)
// -----------------------------------------------------------------------------
async function handleRegister() {
    const username = usernameInput.value;
    if (!username) {
        logMessage("Error: ユーザー名を入力してください。");
        return;
    }

    try {
        // Step 1: サーバーに登録開始をリクエストし、チャレンジを取得
        logMessage(`[CLIENT-REGISTER-1] サーバーに登録開始をリクエスト (User: ${username})`);
        const beginRes = await fetch('/register/begin', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ username })
        });
        
        if (!beginRes.ok) throw new Error(`サーバーエラー: ${await beginRes.text()}`);

        const options = await beginRes.json();
        const sessionID = beginRes.headers.get('X-Session-ID');
        logMessage(`[CLIENT-REGISTER-2] サーバーからチャレンジを含むオプション受信. SessionID: ${sessionID}`);
        logMessage(`Options: ${JSON.stringify(options, null, 2)}`);

        // Step 2: サーバーから受け取ったデータをWebAuthn APIが扱える形式に変換
        options.challenge = base64urlToBuffer(options.challenge);
        options.user.id = base64urlToBuffer(options.user.id);
        if (options.excludeCredentials) {
            options.excludeCredentials.forEach(c => {
                c.id = base64urlToBuffer(c.id);
            });
        }

        // Step 3: WebAuthn APIを呼び出し、認証器にクレデンシャル作成を依頼
        logMessage("[CLIENT-REGISTER-3] navigator.credentials.create() を呼び出し...");
        const credential = await navigator.credentials.create({ publicKey: options });
        logMessage("[CLIENT-REGISTER-4] クレデンシャル作成成功。");

        // Step 4: 作成されたクレデンシャルをサーバーが検証できる形式に変換
        const credentialForServer = {
            id: credential.id,
            rawId: bufferToBase64url(credential.rawId),
            type: credential.type,
            response: {
                attestationObject: bufferToBase64url(credential.response.attestationObject),
                clientDataJSON: bufferToBase64url(credential.response.clientDataJSON),
            },
        };
        logMessage(`[CLIENT-REGISTER-5] サーバーに送信するクレデンシャル: ${JSON.stringify(credentialForServer, null, 2)}`);

        // Step 5: 作成されたクレデンシャルをサーバーに送信して検証・保存を依頼
        logMessage("[CLIENT-REGISTER-6] サーバーに登録完了をリクエスト。");
        const finishRes = await fetch('/register/finish', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Session-ID': sessionID,
            },
            body: JSON.stringify(credentialForServer)
        });

        if (!finishRes.ok) throw new Error(`サーバーエラー: ${await finishRes.text()}`);
        
        const result = await finishRes.json();
        logMessage(`[CLIENT-REGISTER-7] 登録成功！ サーバーからの応答: ${result}`);

    } catch (err) {
        logMessage(`Error: 登録に失敗しました - ${err.message}`);
    }
}

// -----------------------------------------------------------------------------
// 4. ログイン処理 (Authentication)
// -----------------------------------------------------------------------------
async function handleLogin() {
    const username = usernameInput.value;
    if (!username) {
        logMessage("Error: ユーザー名を入力してください。");
        return;
    }

    try {
        // Step 1: サーバーにログイン開始をリクエストし、チャレンジを取得
        logMessage(`[CLIENT-LOGIN-1] サーバーにログイン開始をリクエスト (User: ${username})`);
        const beginRes = await fetch('/login/begin', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ username })
        });

        if (!beginRes.ok) throw new Error(`サーバーエラー: ${await beginRes.text()}`);

        const options = await beginRes.json();
        const sessionID = beginRes.headers.get('X-Session-ID');
        logMessage(`[CLIENT-LOGIN-2] サーバーからチャレンジを含むオプション受信. SessionID: ${sessionID}`);
        logMessage(`Options: ${JSON.stringify(options, null, 2)}`);
        
        // Step 2: データをWebAuthn APIが扱える形式に変換
        options.challenge = base64urlToBuffer(options.challenge);
        if (options.allowCredentials) {
            options.allowCredentials.forEach(c => {
                c.id = base64urlToBuffer(c.id);
            });
        }

        // Step 3: WebAuthn APIを呼び出し、認証器に署名作成を依頼
        logMessage("[CLIENT-LOGIN-3] navigator.credentials.get() を呼び出し...");
        const assertion = await navigator.credentials.get({ publicKey: options });
        logMessage("[CLIENT-LOGIN-4] 署名アサーション取得成功。");

        // Step 4: 作成された署名をサーバーが検証できる形式に変換
        const assertionForServer = {
            id: assertion.id,
            rawId: bufferToBase64url(assertion.rawId),
            type: assertion.type,
            response: {
                authenticatorData: bufferToBase64url(assertion.response.authenticatorData),
                clientDataJSON: bufferToBase64url(assertion.response.clientDataJSON),
                signature: bufferToBase64url(assertion.response.signature),
                userHandle: assertion.response.userHandle ? bufferToBase64url(assertion.response.userHandle) : null,
            },
        };
        logMessage(`[CLIENT-LOGIN-5] サーバーに送信するアサーション: ${JSON.stringify(assertionForServer, null, 2)}`);

        // Step 5: 作成された署名をサーバーに送信して検証を依頼
        logMessage("[CLIENT-LOGIN-6] サーバーにログイン完了をリクエスト。");
        const finishRes = await fetch('/login/finish', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Session-ID': sessionID,
            },
            body: JSON.stringify(assertionForServer)
        });

        if (!finishRes.ok) throw new Error(`サーバーエラー: ${await finishRes.text()}`);

        const result = await finishRes.json();
        logMessage(`[CLIENT-LOGIN-7] ログイン成功！ サーバーからの応答: ${result}`);

    } catch (err) {
        logMessage(`Error: ログインに失敗しました - ${err.message}`);
    }
}
