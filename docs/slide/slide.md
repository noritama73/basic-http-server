---
title: 実装しながら学ぶ！HTTPサーバーの基本 / 技育CAMPアカデミア
author: Takuma Kobayashi @ Finatext
header: 実装しながら学ぶ！HTTPサーバーの基本 / 技育CAMPアカデミア
slide: true
marp: true
theme: meta
---

<script type="module">
  import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid@latest/dist/mermaid.esm.min.mjs';
  mermaid.initialize({ startOnLoad: true });
  window.addEventListener('vscode.markdown.updateContent', function() { mermaid.init() });
</script>

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
- 2000 年生まれの 23 卒
- Software Engineer at [Finatext](https://finatext.com/)
  - １年目からシステムの詳細設計・実装・運用を経験
  - 主な技術スタック
    - Go, AWS, Terraform
- [2024 Japan AWS Jr.Champions](https://aws.amazon.com/jp/blogs/psa/2024-japan-aws-jr-champions-report/)
- 趣味は料理・マラソン

## ![bg right:40%](./TAKUMA.jpeg)

---

## アイスブレイク

皆さんがプログラミングするときに意識していることは何ですか？

---

## 今日の目的

- HTTP サーバーの実装を通して、実践的な開発手法について学ぶ
  - "ただ動くだけのシステム"を脱却する
- 単なる技術解説ではなく、**なぜ**その技術が必要なのかを理解する
- 実装を通して、実際の開発現場でも役立つ知識とスキルを身につける

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

## HTTP サーバーの処理を分割してみよう

---

Hello, World! を返す HTTP サーバーの実装例

```go
func main() {
  mux := http.NewServeMux()

  mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, World!"))
  })

  log.Fatal(http.ListenAndServe(":8080", mux))
}
```

```bash
$ curl http://localhost:8080/
Hello, World!%
```

---

クエリパラメータで名前を受け取り、その名前に対応するメールアドレスを返す HTTP サーバーの実装例

```go
mux.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
  name := r.URL.Query().Get("name")

  email := getEmailByNameFromDB(name)

  w.Write([]byte(email))
})
```

```bash
$ curl "http://localhost:8080/profile?name=hoge"
hoge@example.com%
```

---

HTTP サーバーは何をするシステム？🤔

---

### 抽象化された Web サーバーの処理

<pre class="mermaid">
sequenceDiagram
    participant Client as クライアント
    participant WebServer as サーバー
    participant Database as データベース

    Client->>+WebServer: リクエスト
    WebServer->>WebServer: リクエストからモデルを生成
    WebServer->>+Database: クエリ発行
    Database-->>-WebServer: データ返却
    WebServer->>WebServer: データからモデルを生成
    WebServer-->>-Client: レスポンス
</pre>

---

### Web サーバーの動作はいくつかの処理に分けられる

1. リクエストを受け取る
2. サーバーのデータモデルに読み替える
3. 出来たモデルに対してクエリを組み立て、発行する
4. クエリの結果を受け取り、モデルに変換する
5. モデルからレスポンスを生成する

---

### 分けられた処理をグルーピングしてみる

1. リクエストを受け取る <span style="color: #ff0000">①</span>
2. サーバーのデータモデルに読み替える <span style="color:rgb(255, 140, 60)">②</span>
3. 出来たモデルに対してクエリを組み立て、発行する <span style="color:rgb(50, 0, 255)">③</span>
4. クエリの結果を受け取り、モデルに変換する <span style="color: rgb(255, 140, 60)">②</span>
5. モデルからレスポンスを生成する <span style="color: #ff0000">①</span>

<span style="color: #ff0000">クライアントとの IF 処理</span>
<span style="color:rgb(255, 140, 60)">クライアントとの IF とアプリケーション内データモデルの変換</span>
<span style="color:rgb(50, 0, 255)">アプリケーション内データモデルと DB とのやりとり</span>

---

ということで、それぞれの責務を考えてみましょう

---

## 関数の責務を捉える

---

責務...

アーキテクチャの話...？🤔

---

![bg 70% blur:3px opacity:.4](./CleanArchitecture.jpg)

**クリーンアーキテクチャとは誰々が提唱した概念で〜**
**...という話は今日はしません！**

> > > 画像： <https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html>

---

## なぜ関数の責務に着目するのか？

- その関数の関心ごとを明確にする
  - 設計を変える場合、依存する部分が最小限になる
  - 同じユースケースにおいて再利用できる
  - テストも簡単

---

### 関数の関心ごとを明確にする

```go
// DBに対してクエリを発行し、プロフィールを取得する
func (r *userRepository) GetProfile(name string) (*Profile, error) {
  return r.db.Query("SELECT * FROM profiles WHERE name = ?", name)
}
```

```go
// リクエストからユーザ名を取得し、対応するプロフィールを返す
func (h *handler)GetProfileHandler(w http.ResponseWriter, r *http.Request) {
  name := r.URL.Query().Get("name")

  profile, err := h.UserRepository.GetProfile(name)
  if err != nil {
    http.Error(w, "Profile not found", http.StatusNotFound)
    return
  }

  w.Write([]byte(profile))
}
```

---

## 設計を変える場合、依存する部分が最小限になる

### 処理が関数に分かれていない場合

```go
mux.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
  name := r.URL.Query().Get("name")

  profile := db.Query("SELECT * FROM profiles WHERE name = ?", name)

  w.Write([]byte(profile))
})
```

```go
mux.HandleFunc("/article", func(w http.ResponseWriter, r *http.Request) {
  name := r.URL.Query().Get("name")

  profile := db.Query("SELECT * FROM profiles WHERE name = ?", name) // 実装が重複している
  article := db.Query("SELECT * FROM articles WHERE author = ?", name)

  w.Write([]byte(profile + article))
})
```

---

ユーザー名からプロフィールを取得する処理が重複している
ではプロフィールの項目が増えた場合はどうなるでしょう？
テーブルの構造が変わった場合は？
etc.

---

[関数の関心ごとを明確にする](#関数の関心ごとを明確にする)で紹介したサンプルコードにおいては、

- API のリクエストパラメータはレスポンス形式を変更したい場合は`GetProfileHandler` を
- DB のクエリを変更したい場合は `GetProfile` を

変更すれば良い！

---

### レイヤー分けによる責務の分離

このプロジェクトでは、コードを以下のレイヤーに分けています：

1. **ドメイン層**（`internal/domain`）

   - データモデルの構造や、その振る舞いを定義
   - 例：`User`構造体

2. **ユースケース層**（`internal/usecase`）

   - データモデルを操作したり、外部とのやり取りを行うロジックを定義
   - 例：`Register`、`Login`関数

3. **インターフェース層**（`internal/interface`）

   - 外部(クライアント、DB)とのやり取りを担当
   - 例：`UserHandler`（HTTP）、`InMemoryUserRepository`（データ保存）

---

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

---

### 同じユースケースにおいて再利用できる

- アプリケーションの中で、特定の処理を代表させる
  - その関数に必要なルールを集約できる
  - 例：
    - ユーザー登録の処理を`Register`関数に集約
    - ユーザー名は 3 文字以上、メールアドレスは正しい形式であることを確認
    - ユーザー登録に必ずこの関数を使うようにすれば、上記のチェックを漏らすことがなくなる

---

### テストが簡単

例えば先述のコード

```go
mux.HandleFunc("/article", func(w http.ResponseWriter, r *http.Request) {
  name := r.URL.Query().Get("name")

  profile := db.Query("SELECT * FROM profiles WHERE name = ?", name) // 実装が重複している
  article := db.Query("SELECT * FROM articles WHERE author = ?", name)

  w.Write([]byte(profile + article))
})
```

ここで、

- article テーブルの author カラムをを削除し、代わりに profile テーブルの profile_id カラムを参照するように変更したい場合
- profile 取得のクエリが他の API でも使われている場合

を考えます

---

- article テーブルの author カラムをを削除し、代わりに profile テーブルの profile_id カラムを参照するように変更したい場合
- profile 取得のクエリが他の API でも使われている場合

テーブル構造を変更した影響が、上記クエリを利用しているすべての API のテストコードに波及してしまう
→ 直すのが大変 🫠

---

## 適切なリクエスト処理とバリデーション

---

### ちょっと脱線

プログラムの特性を一言で表すと何だと思いますか？

---

それは「**書いた通りにしか動かない**」です。

正常に処理できないリクエストを受け取った際に、想定していない動作をしたり、その結果データが壊れたり、不正なリクエストを外部に送ったりしてしまう可能背があります。

---

## ハンドリングが必要なデータの例

- 携帯電話番号
  - 日本国内に限って言えば、
    - 10 桁の数字
    - 最初の 3 桁は、060, 070, 080, 090 のいずれか
- マイナンバー
  - 12 桁の数字
  - チェックディジットの計算が必要

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

### バリデーションの階層化が重要な理由

- **防衛的プログラミング**：複数の層でチェックすることで安全性が向上
- **責務の明確化**：各層は自分の責任範囲内でのみ検証を行う
- **エラーの適切な処理**：発生場所に応じた適切なエラーハンドリングが可能
- **セキュリティの多層化**：一つの層のバグや見落としがあっても他の層でカバー

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

### セキュリティ対策が不十分だと何が起きるか？

- **データ漏洩**：顧客情報や機密情報の流出
- **サービス停止**：DoS 攻撃によるシステムダウン
- **データ改ざん**：不正なデータ操作
- **権限昇格**：一般ユーザーが管理者権限を取得
- **風評被害**：セキュリティインシデントによる信頼喪失

---

## 認証

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

### 今日学んだこと

1. **関数の責務を捉える**

   - 単一責任の原則に基づいたコード設計
   - レイヤー分けによる関心事の分離

2. **適切なリクエスト処理とバリデーション**

   - 多層的なバリデーションの重要性
   - セキュリティを考慮したエラーハンドリング

3. **認証の基本と実装方法**
   - JWT を使ったステートレス認証
   - セキュアな認証フローの実装

### 次のステップ

- **学んだ概念を自分のプロジェクトに適用してみる**
- **他の認証方式（OAuth、多要素認証など）について学ぶ**
- **より複雑なユースケースでの実装を試みる**
- **テストの書き方を学び、堅牢なコードを目指す**

---

## 宣伝

---

![bg 80%](./summer_internship_summer_internship_engineer.png)

---

**サマーインターンやります！**

**7/6 締め切りなので急いで！**

![bg 80% opacity:.5](./summer_internship_summer_internship_engineer.png)

---

## ご清聴ありがとうございました！

質問やフィードバックがあればお気軽にどうぞ！
例）

- 「自分なりにコード書いてみたのでレビューしてください！」
- 「この部分の実装が難しいのでアドバイスください！」

X の DM でも受け付けています！→
![bg right:35% width:450px](./x.png)

教材のリポジトリ：<https://github.com/noritama73/basic-http-server>
