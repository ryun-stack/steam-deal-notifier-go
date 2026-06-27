# Credits & Third-Party Notices

本プロジェクト（`steam-deal-discord-notifier`）の開発および運行にあたり、以下のサードパーティAPI、クラウドサービス、およびオープンソースライブラリを利用しています。素晴らしいサービスとコミュニティの成果物に感謝いたします。

---

## 1. 3rd-Party APIs & Services

本システムは、以下の外部サービスとAPI連携を行っています。各サービスの利用規約およびガイドラインを遵守して運用されています。

* **IsThereAnyDeal API**
  * ゲームのセール情報および価格比較データの取得に利用。
  * Data provided by [IsThereAnyDeal](https://isthereanydeal.com/) (Subject to [ITAD Terms of Service](https://docs.isthereanydeal.com/#section/Terms-of-Service))
* **IGDB API**
  * ゲームの評価、レビュー数、サムネイルなどの詳細メタデータの取得に利用。
  * Game data provided by [IGDB](https://www.igdb.com/) (Subject to [Twitch Developer Agreement](https://legal.twitch.com/legal/developer-agreement/))
* **OpenAI API**
  * 取得したゲーム概要を自然な日本語に要約・動記生成するために利用。
  * Text generation powered by [OpenAI API](https://openai.com/)
* **Discord API**
  * Discord Embed（カード形式）を用いたリッチなセール情報の定期通知に利用。
  * Bot integration provided by [Discord](https://discord.com/)

---

## 2. Open Source Software (Go Modules)

本アプリケーションの構築には、以下のオープンソースソフトウェア（OSS）が使用されています。

### ── Direct Dependencies (直接依存)

| パッケージ名 | ライセンス | 用途 |
| :--- | :--- | :--- |
| `github.com/azure/azure-functions-golang-worker` | MIT License | Azure Functions 上で Go をカスタムハンドラーとして動作させるためのワーカーランタイム |
| `github.com/openai/openai-go/v2` | Apache-2.0 | OpenAI API を Go からセキュアに呼び出すための公式SDK |

### ── Indirect Dependencies (間接依存・ユーティリティ)

| パッケージ名 | ライセンス | 用途 |
| :--- | :--- | :--- |
| `github.com/tidwall/gjson` / `sjson` / `match` / `pretty` | MIT License | APIから取得したレスポンス（JSON）を、型を過度に共有せず軽量かつ柔軟にパース・操作するためのライブラリ群 |
| `github.com/spf13/pflag` | BSD-3-Clause | コマンドライン引数のパース（内部利用） |
| `golang.org/x/net` / `sys` / `text` | BSD-3-Clause | Go公式によるネットワーク・システム・テキスト処理用の拡張ライブラリ |
| `google.golang.org/grpc` / `protobuf` / `genproto` | Apache-2.0 | インフラ層やワーカー内の通信・シリアライズ用コンポーネント |