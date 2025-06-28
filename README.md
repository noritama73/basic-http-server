# Academia Basic HTTP Server

このプロジェクトは、エンジニアを目指す学生向けの基本的な HTTP サーバー実装のサンプルアプリケーションです。シンプルなレイヤードアーキテクチャを採用し、認証には JWT を利用しています。

## アーキテクチャ

このアプリケーションは以下のレイヤードアーキテクチャを採用しています：

1. **ドメイン層** (`internal/domain/`): エンティティとリポジトリインターフェースを定義
2. **ユースケース層** (`internal/usecase/`): ビジネスロジックを実装
3. **インターフェース層** (`internal/interface/`):
   - `handler/`: HTTP リクエストの処理
   - `repository/`: データアクセスの実装
4. **ユーティリティ** (`internal/util/`): JWT 認証などの共通機能

## 機能

- ユーザー登録
- ユーザーログイン（JWT 認証）
- ユーザー情報の取得（認証済みユーザーのみ）

## ディレクトリ構成

```txt
academia-basic-http-server/
├── cmd/                # エントリーポイント
│   └── main.go         # アプリケーションの起動コード
├── internal/           # アプリケーションの内部ロジック
│   ├── domain/         # ドメイン層（エンティティやビジネスロジック）
│   │   └── user.go     # ユーザーエンティティ
│   ├── usecase/        # ユースケース層（アプリケーションのユースケース）
│   │   └── user_usecase.go # ユーザー関連のユースケース
│   ├── interface/      # インターフェース層（ハンドラーやリポジトリ）
│   │   ├── handler/    # HTTPハンドラー
│   │   │   └── user_handler.go
│   │   └── repository/ # データアクセス
│   │       └── user_repository.go
│   └── util/           # 再利用可能なユーティリティ
│       └── jwt.go      # JWT関連のユーティリティ
└── go.mod              # Goモジュールファイル
```

## 実行方法

```bash
# サーバーの起動
go run cmd/main.go

# 環境変数でポートとJWTシークレットを指定することも可能
PORT=3000 JWT_SECRET=your-secret-key go run cmd/main.go
```

デフォルトでは、サーバーは 8080 ポートで起動します。

## API エンドポイント

### ユーザー登録

```bash
POST /register
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123",
  "email": "test@example.com"
}
```

### ログイン

```bash
POST /login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

レスポンスで JWT トークンが返されます。

### ユーザー情報の取得

```bash
GET /users/{user_id}
Authorization: Bearer {jwt_token}
```

## 学習ポイント

このサンプルアプリケーションでは以下の点を学ぶことができます：

1. **レイヤードアーキテクチャ**: 関心事の分離と責任の明確化
2. **インターフェース**: 依存性の逆転と疎結合な設計
3. **認証**: JWT を使った認証の基本的な実装
4. **HTTP サーバー**: 標準ライブラリのみを使った HTTP サーバーの実装
5. **エラーハンドリング**: 適切な HTTP ステータスコードの返却
6. **ミドルウェア**: ロギングなどの横断的関心事の実装

## 注意点

このサンプルアプリケーションは学習目的で作成されており、以下の点に注意してください：

- インメモリリポジトリを使用しているため、サーバー再起動でデータは失われます
- 本番環境では適切なデータベースの使用を検討してください
- JWT シークレットは環境変数から取得するようにし、ソースコードにハードコードしないでください

## 作動例

```bash
[takuma.kobayashi@fn]: ~/noritama73/academia-basic-http-server
$ curl -X POST -H "Content-Type: application/json" -d '{"username":"testuser","password":"password123","email":"test@example.com"}' http://localhost:8080/register
{"id":"d4fd81f023f26b7acef2adc602ad0e77","username":"testuser","email":"test@example.com","created_at":"2025-06-09T21:03:30.870816+09:00"}

[takuma.kobayashi@fn]: ~/noritama73/academia-basic-http-server
$ curl -X POST -H "Content-Type: application/json" -d '{"username":"testuser","password":"password123"}' http://localhost:8080/login
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZDRmZDgxZjAyM2YyNmI3YWNlZjJhZGM2MDJhZDBlNzciLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzQ5NTU3MDI3LCJpYXQiOjE3NDk0NzA2Mjd9.3i3UwL1-dsQ7luMwS6vOYaVMiU_60TMbg6II7Sb0umA"}

[takuma.kobayashi@fn]: ~/noritama73/academia-basic-http-server
$ curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZDRmZDgxZjAyM2YyNmI3YWNlZjJhZGM2MDJhZDBlNzciLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhw
IjoxNzQ5NTU3MDI3LCJpYXQiOjE3NDk0NzA2Mjd9.3i3UwL1-dsQ7luMwS6vOYaVMiU_60TMbg6II7Sb0umA" http://localhost:8080/users/d4fd81f023f26b7acef2adc602ad0e77
{"id":"d4fd81f023f26b7acef2adc602ad0e77","username":"testuser","email":"test@example.com","created_at":"2025-06-09T21:03:30.870816+09:00","updated_at":"2025-06-09T21:03:30.870816+09:00"}
```
