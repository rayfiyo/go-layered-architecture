# go-layered-architecture

- An example of layered architecture + DDD in Go and notes
- Go 言語におけるレイヤードアーキテクチャ + DDD の一例とメモ
  - ユーザーの作成と取得の API を実装した
  - DB 周りは sqlx を用いて実装した

# 関連記事

- https://qiita.com/os1ma/items/7a229585ebdd8b7d86c2
  - アーキテクチャの基礎知識
  - 勘違いしがちな周辺知識の関係について記載がある
    > 実は、MVC、MVP、MVVM といったものは、
    > 全てプレゼンテーション層のアーキテクチャなのです。
    >
    > なので、アプリケーションの構成を検討するときに、「MVC にするぜ」というだけでなく、
    > 「全体としては 3 層で、プレゼンテーション層は MVC にする」という話になるわけです。
    >
    > 多くの例では、ビジネスロジック層を「アプリケーション層」と「ドメイン層」の
    > 2 つに分離しています。
    >
    > レイヤー構成 ... 3 層、ヘキサゴナル、クリーンなどから選ぶ
    > プレゼンテーション層 ... MVC、MVP、MVVM などから選ぶ
    > ビジネスロジック層 ... トランザクションスクリプト、ドメインモデルから選ぶ
- https://qiita.com/tono-maron/items/345c433b86f74d314c8d
  - 非常に詳しく レイヤードアーキテクチャ + DDD の実装が書かれている
  - 別の実装方法: https://engineerblog.mynavi.jp/technology/goddd-layered-architecture/

# 構成

```
go-layered-architecture/
├── cmd/
│   └── main.go // アプリケーションのエントリーポイント
│               // 各レイヤーを組み合わせてサーバーを起動する
├── internal/
│   ├── application/    // アプリケーション層：ビジネスロジックを担当
│   │   └── service/
│   │       └── user_service.go
│   │
│   ├── domain/         // ドメイン層: ビジネスモデル・エンティティの定義を担当
│   │   ├── user.go
│   │   └── user_repository.go
│   │
│   ├── infrastructure/ // インフラ層：外部技術の実装を担当
│   │   └── repository/
│   │       └── user_repository_sqlx.go
│   │
│   └── presentation    // プレゼンテーション層：HTTP リクエスト受付とレスポンス処理を担当
│       └── handler
│           ├── router.go
│           └── user_handler.go
│
└── go.mod // Go モジュール定義ファイル
```

以下、実装を行う順に説明する。

## 1. ドメイン層 (`internal/domain/`)

- ビジネスモデル・エンティティを定義する
  - 今回は、ユーザーのエンティティを定義した
- 愚直なレイヤードアーキテクチャであればドメイン層はインフラ層に依存するが、
  今回はそれを防ぐ考えの DDD レイヤードアーキテクチャなので、
  `user_repo.go` では `interface` （抽象型） を使って実装した
  - DDD ではドメイン層はどこの層にも依存せず単体で完結する
- 回は小規模なので `domain/` のみだが、`domain/repository/user_repo.go` などにしても良い

## 2. インフラ層（`internal/infrastructure/`）

- 外部技術の実装を行う
  - 例えば、ユーザーの保存や取得など、
    データの永続化に関する操作を抽象化するための抽象型と、
    その実装を行う（レポジトリ層という）
  - 例えば、外部 API を呼び出す抽象型と、その実装を行う
- DDD でないレイヤードアーキテクチャであれば、ドメイン層よりも先に実装する
  - これは、ドメイン層がインフラ層に依存するためである
  - DDD はこれを良しとしない設計思想なので、
    抽象型を使った抽象化で依存性逆転の原則 (DIP) を行う
    - これによって、インフラ層がドメイン層に依存するようにできる
- 今回は小規模なので、データの保存取得の `repository/` （レポジトリ層）のみになった
  - このような場合、インフラ層ではなく単にレポジトリ層と言ったり、
    `internal/repository/` と実装することもある
  - 逆に、大規模になると `persistence/` `cache/` `external/` `queue/` `api/` などを生やす

## 3. アプリケーション層（`internal/application/`）

- リポジトリ層を利用してビジネスロジックの実装を行う
  - アプリケーション処理のコードを整理するための層
  - 多くはユーザー操作単位の処理（ユースケース）を実装する
  - 例えば、ユーザーの取得や作成、商品購入などの処理を行う
- アプリケーション層の多くはユースケース (`usecase`) なので、ユースケース層ともいう
  - また、ユースケース層は概念名なので、実装名のサービス層とも言う
    - つまり、ユースケース層を `usecase/` ではなく `service/` と実装することもある
    - 特に、OOP の文脈では `usecase/` を用いて実装するが、
      OOP 思想の薄い Go では `service/` を用いることが多い
    - また、これは Application Service の service である
- 今回は小規模なので、ユーザー登録の `service/` （ユースケース層）のみになった
  - このような場合、アプリケーション層ではなく単にユースケース層と言ったり、
    `internal/service/` と実装することもある
  - 逆に、大規模になると `dot/` などを生やす
    - トランザクション管理や認可をこの層に置くこともある

## 4. プレゼンテーション層（`internal/presentation/`）

- ユーザー向けの入力と出力を実装する
  - 例えば、HTTP リクエストを受け付けて、アプリケーション層の処理を呼び出し、
    レスポンスを返すなどの処理を行う
    - 出力（レスポンス）も同様
  - その他にも、HTML、JSON、UI、API などの request/response 実装を行うこともある
  - 今回は HTTP サーバーの実装を行った
- 特にクリーンアーキテクチャでは、インターフェースアダプタ層とも呼ばれる
  - これは、アプリケーション外からの入力を
    アプリケーション層の処理に適した形式に変換する役割を持つためである
  - そのため、`interfaces/` や `adapter/` などと実装することもある
  - また、レイヤードアーキテクチャと比べると、より抽象的なニュアンスになる
    - これは、ユーザーに見えない外部の入出力、
      例えば gRPC、WebSocket、Message queue、Batch job なども扱うため、
      より抽象的な概念にしようという意図がある
  - なお interfaces と presentation と handler などはあまり区別されないことも多い
    - 特に Go では、HTTP サーバーの実装を handler と呼ぶことが多い
      - 逆に、MVC/Java/Rails などの文脈では、
        HTTP サーバーの実装を controller と呼ぶことが多い
- 今回は小規模なので、HTTP サーバーの `handler/` のみになった
  - このような場合、プレゼンテーション層ではなく単にハンドラー層と言ったり、
    `internal/handler/` と実装することもある
  - 逆に、大規模になると `middleware/` などを生やす

## その他

- エントリーポイント (cmd/main.go)
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
