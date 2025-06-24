---
title: 実装しながら学ぶ！HTTPサーバーの基本 / 技育CAMPアカデミア
author: Takuma Kobayashi @ Finatext
header: 実装しながら学ぶ！HTTPサーバーの基本 / 技育CAMPアカデミア
slide: true
marp: true
theme: meta
---

# 実装しながら学ぶ！HTTP サーバー実装の基本

Takuma Kobayashi ([@takuma5884rbb](https://x.com/takuma5884rbb))
株式会社 Finatext

---

## 目次

<!-- paginate: true -->

- 自己紹介
- 今日学んで欲しいこと
- リクエストを受け取る
- 異常系を考える
- 認証
- まとめ

---

## 自己紹介

- 小林拓磨
  - X: [@takuma5884rbb](https://x.com/takuma5884rbb)
- 2000 年生まれ
- Software Engineer at [Finatext](https://finatext.com/)
  - １年目からシステムの詳細設計・実装・運用を経験
  - 主な技術スタック
    - Go, AWS, Terraform
- [2024 Japan AWS Jr.Champions](https://aws.amazon.com/jp/blogs/psa/2024-japan-aws-jr-champions-report/)
- 趣味は料理・マラソン

![bg right:45%](./TAKUMA.jpeg)

---

## 今日の目的

- HTTP サーバーの実装を通して、実践的な開発手法について学ぶ
  - "ただ動くだけのシステム"を脱却する
- 使用する技術スタック
  - Go 言語（シンプルで読みやすい構文が特徴）
  - 標準ライブラリのみを使用（フレームワークに依存しない基本を学ぶ）

---

## 今日学んで欲しいこと

- 関数の責務を捉える重要性
  - 一つの関数は一つの責任を持つべき
  - コードの可読性と保守性の向上
- リクエスト処理とバリデーション
  - 不正なデータからシステムを守る方法
  - ユーザー入力を常に疑う姿勢
- 認証の基本と実装方法
  - セキュリティの重要性
  - JWT を使った認証の仕組み

---

## 目次：リクエストを受け取る

- HTTP リクエストの基本構造
  - メソッド（GET, POST, PUT, DELETE）
  - パス（URL）
  - ヘッダー（メタデータ）
  - ボディ（データ）
- サーバーの初期化と設定
  - ポートの設定
  - ルーティング（どのパスでどの処理を行うか）
- ミドルウェアの活用
  - リクエストの前処理・後処理
  - ログ記録、認証チェックなど

```go
// cmd/main.go から抜粋
mux := http.NewServeMux()  // ルーティングを管理するマルチプレクサを作成
userHandler.RegisterRoutes(mux)  // ルートを登録
handler := logMiddleware(mux)  // ログ記録用のミドルウェアを適用
```

---

## 目次：異常系を考える

- エラーハンドリングの重要性
  - ユーザーに適切なエラーメッセージを返す
  - システム内部のエラーは隠す
- 想定されるエラーパターン
  - 不正な入力データ
  - リソースが見つからない
  - 認証・認可エラー
  - サーバー内部エラー
- エラーの種類に応じた HTTP ステータスコード
  - 400: Bad Request（クライアントのリクエストに問題）
  - 401: Unauthorized（認証が必要）
  - 404: Not Found（リソースが存在しない）
  - 500: Internal Server Error（サーバー内部のエラー）

```go
// internal/interface/handler/user_handler.go から抜粋
if err != nil {
    if err == usecase.ErrInvalidCredentials {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
    } else {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    return
}
```

---

## 目次：認証

- 認証とは何か
  - ユーザーが本人であることを確認するプロセス
  - システムのセキュリティの基盤
- JWT（JSON Web Token）の仕組み
  - ヘッダー、ペイロード、署名の 3 つの部分からなる
  - 改ざん検知が可能
  - ステートレス（サーバーにセッション情報を保存しない）
- 認証フロー
  - ユーザー登録
  - ログイン（認証情報の検証とトークン発行）
  - 保護されたリソースへのアクセス（トークン検証）

```go
// internal/util/jwt.go から抜粋
// Generate creates a new JWT token for a user
func (m *JWTManager) Generate(userID, username string) (string, error) {
    now := time.Now()
    claims := JWTClaims{
        UserID:    userID,
        Username:  username,
        ExpiresAt: now.Add(m.expiry).Unix(),
        IssuedAt:  now.Unix(),
    }
    // ... トークン生成処理 ...
}
```

---

## 今日学んで欲しいこと

---

## 今日学んで欲しいこと

1. **関数の責務を捉える**

   - 関心事の分離による保守性の向上
   - 責任範囲の明確化

2. **適切なリクエスト処理とバリデーション**

   - 不正なリクエストからサーバーを守る
   - エラーハンドリングの重要性

3. **認証の基本と実装方法**
   - セキュアな認証の実現
   - トークンベース認証のメリット・デメリット

---

## 関数の責務を捉える

---

責務...

アーキテクチャの話...？🤔

---

**クリーンアーキテクチャとは誰々が提唱した概念で〜**

![bg 70% blur:3px opacity:.4](./CleanArchitecture.jpg)

**...という話は今日はしません！**

> > > 画像： https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

---

## 関数の責務を捉える（1/5）

### 関数の責務とは？

- **責務**：その関数が果たすべき役割や責任
- **単一責任の原則**：一つの関数は一つのことだけを行うべき

### なぜ重要か？

- コードの**可読性**が向上する
- **テスト**がしやすくなる
- **再利用**しやすくなる
- **バグ**が少なくなる

```go
// 良い例：関数名が責務を明確に表している
func (r *InMemoryUserRepository) FindByUsername(username string) (*domain.User, error) {
    // ユーザー名からユーザーを検索する処理のみを行う
}
```

---

## 関数の責務を捉える（2/5）

### レイヤー分けによる責務の分離

このプロジェクトでは、コードを以下のレイヤーに分けています：

1. **ドメイン層**（`internal/domain`）

   - ビジネスエンティティとルールを定義
   - 例：`User`構造体

2. **ユースケース層**（`internal/usecase`）

   - ビジネスロジックを実装
   - 例：`Register`、`Login`関数

3. **インターフェース層**（`internal/interface`）

   - 外部とのやり取りを担当
   - 例：`UserHandler`（HTTP）、`InMemoryUserRepository`（データ保存）

4. **インフラ層**（`cmd`）
   - アプリケーションの起動と設定
   - 例：`main`関数

---

## 関数の責務を捉える（3/5）

### 実例：ユーザー登録の処理

1. **ハンドラー**（HTTP 処理）

```go
// internal/interface/handler/user_handler.go
func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
    // HTTPリクエストの解析とレスポンス返却のみを担当
    body, err := io.ReadAll(r.Body)
    // ...
    var req usecase.RegisterRequest
    if err := json.Unmarshal(body, &req); err != nil {
        http.Error(w, "Invalid request format", http.StatusBadRequest)
        return
    }
    // ビジネスロジックはユースケースに委譲
    resp, err := h.userUseCase.Register(req)
    // ...
}
```

---

## 関数の責務を捉える（4/5）

### 実例：ユーザー登録の処理（続き）

2. **ユースケース**（ビジネスロジック）

```go
// internal/usecase/user_usecase.go
func (uc *UserUseCase) Register(req RegisterRequest) (*RegisterResponse, error) {
    // 入力検証
    if req.Username == "" || req.Password == "" || req.Email == "" {
        return nil, ErrInvalidInput
    }

    // ビジネスルール適用（ユーザー名の重複チェックなど）
    _, err := uc.userRepo.FindByUsername(req.Username)
    if err == nil {
        return nil, errors.New("username already taken")
    }

    // パスワードハッシュ化などのビジネスロジック
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    // ...

    // データ保存はリポジトリに委譲
    if err := uc.userRepo.Store(user); err != nil {
        return nil, err
    }
    // ...
}
```

---

## 関数の責務を捉える（5/5）

### 実例：ユーザー登録の処理（続き）

3. **リポジトリ**（データアクセス）

```go
// internal/interface/repository/user_repository.go
func (r *InMemoryUserRepository) Store(user *domain.User) error {
    // データ保存のみを担当
    r.mutex.Lock()
    defer r.mutex.Unlock()

    // 重複チェック
    if _, exists := r.users[user.ID]; exists {
        return ErrUserExists
    }

    // データ保存
    r.users[user.ID] = user
    return nil
}
```

### ポイント

- 各関数は明確な責務を持ち、他の層の詳細を知らない
- これにより、コードの変更が他の部分に影響しにくくなる

---

## 適切なリクエスト処理とバリデーション

---

### ちょっと脱線

プログラムの特性を一言で表すと何だと思いますか？

---

それは「**書いた通りにしか動かない**」です。

正常に処理できないリクエストを受け取った際に、想定していない動作をしたり、その結果データが壊れたり、不正なリクエストを外部に送ったりしてしまう可能背があります。

---

## 適切なリクエスト処理とバリデーション（1/5）

### リクエスト処理の基本

- HTTP リクエストの構成要素
  - **メソッド**：GET, POST, PUT, DELETE など
  - **パス**：リソースの場所（例：`/users/123`）
  - **ヘッダー**：メタデータ（例：`Content-Type`, `Authorization`）
  - **ボディ**：送信データ（JSON, フォームデータなど）

```go
// internal/interface/handler/user_handler.go
func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
    // メソッドチェック
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // ボディ読み取り
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    // ...
}
```

---

## 適切なリクエスト処理とバリデーション（2/5）

### バリデーションの重要性

- **バリデーション**：入力データが期待通りかを確認すること
- なぜ重要か？
  - **セキュリティ**：悪意ある入力からシステムを守る
  - **データ整合性**：不正なデータがシステムに入らないようにする
  - **ユーザー体験**：早期にエラーを検出し、適切なフィードバックを提供

### バリデーションのレベル

1. **構文的バリデーション**：データ形式が正しいか

   ```go
   if err := json.Unmarshal(body, &req); err != nil {
       http.Error(w, "Invalid request format", http.StatusBadRequest)
       return
   }
   ```

2. **意味的バリデーション**：ビジネスルールに沿っているか
   ```go
   if req.Username == "" || req.Password == "" || req.Email == "" {
       return nil, ErrInvalidInput
   }
   ```

---

## 適切なリクエスト処理とバリデーション（3/5）

### バリデーションの場所

- **ハンドラー層**：HTTP リクエストの形式チェック

  ```go
  // internal/interface/handler/user_handler.go
  if r.Method != http.MethodPost {
      http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
      return
  }
  ```

- **ユースケース層**：ビジネスルールに基づくチェック

  ```go
  // internal/usecase/user_usecase.go
  if req.Username == "" || req.Password == "" || req.Email == "" {
      return nil, ErrInvalidInput
  }
  ```

- **リポジトリ層**：データ整合性のチェック
  ```go
  // internal/interface/repository/user_repository.go
  if _, exists := r.users[user.ID]; exists {
      return ErrUserExists
  }
  ```

---

## 適切なリクエスト処理とバリデーション（4/5）

### エラーハンドリング

- エラーの種類に応じた適切な HTTP ステータスコードを返す
  - **400 Bad Request**：クライアントのリクエストに問題がある
  - **401 Unauthorized**：認証が必要
  - **403 Forbidden**：権限がない
  - **404 Not Found**：リソースが存在しない
  - **500 Internal Server Error**：サーバー内部のエラー

```go
// internal/interface/handler/user_handler.go
func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
    // ...
    resp, err := h.userUseCase.Login(req)
    if err != nil {
        if err == usecase.ErrInvalidCredentials {
            // 認証エラーは401
            http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        } else {
            // その他のエラーは500
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }
    // ...
}
```

---

## 適切なリクエスト処理とバリデーション（5/5）

### セキュリティに関する注意点

- **入力は常に疑う**：すべてのユーザー入力は潜在的に危険
- **最小権限の原則**：必要最小限の権限だけを与える
- **内部エラーの詳細は隠す**：攻撃者に情報を与えない

  ```go
  // 悪い例
  if err != nil {
      http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
  }

  // 良い例
  if err != nil {
      log.Printf("Database error: %v", err) // 内部でログ記録
      http.Error(w, "Internal server error", http.StatusInternalServerError) // ユーザーには最小限の情報
  }
  ```

- **機密データの保護**：パスワードなどの機密情報は適切に保護
  ```go
  // パスワードのハッシュ化
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
  ```

---

## 認証の基本と実装方法（1/5）

### 認証とは

- **認証（Authentication）**：ユーザーが本人であることを確認するプロセス

  - 「あなたは誰ですか？」という質問に答える

- **認可（Authorization）**：ユーザーが特定のリソースにアクセスする権限があるかを確認するプロセス
  - 「あなたは何ができますか？」という質問に答える

---

## 認証の基本と実装方法（2/5）

### JWT（JSON Web Token）とは

- **JWT**：JSON 形式のデータを安全に転送するための標準規格
- **構成**：
  1. **ヘッダー**：トークンのタイプと使用しているアルゴリズム
  2. **ペイロード**：クレーム（ユーザー ID などの情報）
  3. **署名**：トークンが改ざんされていないことを確認するための署名

```go
// internal/util/jwt.go
type JWTClaims struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    // 標準クレーム
    ExpiresAt int64 `json:"exp"`
    IssuedAt  int64 `json:"iat"`
}
```

---

### JWT のメリット

- **ステートレス**：サーバーにセッション情報を保存する必要がない
- **スケーラビリティ**：複数のサーバー間で認証情報を共有しやすい
- **クロスドメイン**：異なるドメイン間でも使用可能

---

## 認証の基本と実装方法（3/5）

### JWT の生成

```go
func (m *JWTManager) Generate(userID, username string) (string, error) {
    now := time.Now()
    claims := JWTClaims{
        UserID:    userID,
        Username:  username,
        ExpiresAt: now.Add(m.expiry).Unix(),  // トークンの有効期限
        IssuedAt:  now.Unix(),                // トークンの発行時刻
    }

    header := map[string]string{
        "alg": "HS256",  // 署名アルゴリズム
        "typ": "JWT",    // トークンタイプ
    }

    // ヘッダーとペイロードをBase64エンコード
    headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)
    payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)

    // 署名の作成
    signatureInput := headerEncoded + "." + payloadEncoded
    h := hmac.New(sha256.New, m.secretKey)
    h.Write([]byte(signatureInput))
    signature := h.Sum(nil)
    signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)

    // JWTの組み立て
    token := headerEncoded + "." + payloadEncoded + "." + signatureEncoded
    return token, nil
}
```

---

## 認証の基本と実装方法（4/5）

### JWT の検証

```go
func (m *JWTManager) Verify(token string) (*JWTClaims, error) {
    // トークンを3つの部分に分割
    parts := strings.Split(token, ".")
    if len(parts) != 3 {
        return nil, ErrInvalidToken
    }

    // 署名の検証
    signatureInput := parts[0] + "." + parts[1]
    h := hmac.New(sha256.New, m.secretKey)
    h.Write([]byte(signatureInput))
    signature := h.Sum(nil)
    expectedSignature := base64.RawURLEncoding.EncodeToString(signature)

    if parts[2] != expectedSignature {
        return nil, ErrInvalidToken  // 署名が一致しない場合はエラー
    }

    // ペイロードのデコードとパース
    payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
    // ...

    // 有効期限のチェック
    if time.Now().Unix() > claims.ExpiresAt {
        return nil, ErrExpiredToken  // 有効期限切れの場合はエラー
    }

    return &claims, nil  // 検証成功
}
```

---

## 認証の基本と実装方法（5/5）

### 認証フロー

1. **ユーザー登録**

   ```go
   // internal/usecase/user_usecase.go
   func (uc *UserUseCase) Register(req RegisterRequest) (*RegisterResponse, error) {
       // パスワードのハッシュ化
       hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
       // ユーザー情報の保存
       // ...
   }
   ```

---

2. **ログイン（認証）**

   ```go
   // internal/usecase/user_usecase.go
   func (uc *UserUseCase) Login(req LoginRequest) (*LoginResponse, error) {
       // ユーザー検索
       user, err := uc.userRepo.FindByUsername(req.Username)
       // パスワード検証
       if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
           return nil, ErrInvalidCredentials
       }
       // JWTトークン生成
       token, err := uc.jwtManager.Generate(user.ID, user.Username)
       // ...
   }
   ```

---

3. **保護されたリソースへのアクセス**

   ```go
   // internal/interface/handler/user_handler.go
   func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
       // トークンの抽出と検証
       claims, err := h.jwtManager.Verify(token)
       // アクセス制御
       if claims.UserID != userID {
           http.Error(w, "Unauthorized access", http.StatusForbidden)
           return
       }
       // ...
   }
   ```

---

## まとめ

---

## 宣伝

---

**サマーインターンやります！**

**7/6 締め切りなので急いで！**

![bg 80% opacity:.5](./summer_internship_summer_internship_engineer.png)

---

## ご清聴ありがとうございました！

質問やフィードバックがあればお気軽にどうぞ！
X の DM でも受け付けています！→
![bg right:35% width:450px](./x.png)

教材のリポジトリ：<https://github.com/noritama73/basic-http-server>
