# go-layered-architecture

- An example of layered architecture in Go and notes
- Go 言語におけるレイヤードアーキテクチャの一例とメモ

```
go-layered-architecture/
├── cmd
│   └── main.go         // アプリケーションのエントリーポイント。各レイヤーを組み合わせてサーバーを起動する
├── internal
│   ├── domain          // ドメイン層：ビジネスモデル・エンティティを定義
│   │   └── user.go     // ユーザーのエンティティを定義
│   ├── repository      // リポジトリ層：データアクセスを担当
│   │   └── user_repository.go  // ユーザーのデータ保存／取得の実装（ここでは InMemory の例）
│   ├── service         // サービス層：ビジネスロジックを担当
│   │   └── user_service.go     // ユーザーに関するビジネスロジックの実装
│   └── handler         // ハンドラー層：HTTP リクエストの受付とレスポンス処理を担当
│       └── user_handler.go     // ユーザーに関する API エンドポイントの実装
└── go.mod              // Go モジュール定義ファイル
```

## 解説

1. ドメイン層 (internal/domain)
   - ユーザーというエンティティ（ビジネスモデル）を定義
   - ここでは User 構造体を用意
2. リポジトリ層 (internal/repository)
   - ユーザーの保存や取得など、データの永続化に関する操作を抽象化するための
     インターフェースと、その実装（ここではメモリ上にデータを保持する
     InMemoryUserRepository）を実装
3. サービス層 (internal/service)
   - リポジトリ層を利用して、ビジネスロジックを実装
   - たとえば、ユーザーの取得や作成の処理を行う
4. ハンドラー層 (internal/handler)
   - HTTP リクエストを受け付け、サービス層の処理を呼び出してレスポンスを返
   - ここでは GET リクエストでユーザーの取得、POST リクエストでユーザーの作成を行う
5. エントリーポイント (cmd/main.go)
   - 各層のコンポーネントを初期化し、HTTP サーバーのルーティングを設定してサーバーを起動

## cURL

`-i` や `-v` でより詳しく

### ユーザー作成 (POST)

```
curl -X POST -H "Content-Type: application/json" \
-d '{"name": "John Doe", "email": "john@example.com"}' \
http://localhost:8080/users
```

### ユーザー取得 (GET)

```
curl -X GET 'http://localhost:8080/users?id=1'
```
