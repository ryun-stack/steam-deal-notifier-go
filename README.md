## Getting Started

本プロジェクトをローカル環境で実行・デバッグするための手順です。

### 前提条件
* **Go** (1.26以上)
* **Azure Functions Core Tools** (v4.x)
* **Azurite** (ローカルストレージエミュレータ。VS Code拡張機能またはDocker)
* **Discordサーバ**
* **Microsoft AI Founrary**

---

### 1. リポジトリのクローンと依存関係の解決
リポジトリをクローンし、依存パッケージをダウンロードします。

```bash
git clone https://github.com/ryun-stack/steam-deal-discord-notifier.git
cd steam-deal-discord-notifier
go mod download
```

### 2. 環境変数の設定 (local.settings.json の作成)
ルートディレクトリに `local.settings.json` を作成し、各種APIキーと接続情報を設定します。

```json
{
  "IsEncrypted": false,
  "Values": {
    "AzureWebJobsStorage": "UseDevelopmentStorage=true",
    "FUNCTIONS_WORKER_RUNTIME": "native",
    "ITAD_API_KEY": "",
    "DISCORD_WEBHOOK_URL": "",
    "AOAI_CHAT_COMPLETIONS_MODEL": "",
    "AZURE_OPENAI_ENDPOINT": "",
    "AZURE_OPENAI_API_KEY": "",
    "TWITCH_CLIENT_ID": "",
    "TWITCH_CLIENT_SECRET": "",
    "TWITCH_OAUTH_ENDPOINT": "",
    "ITAD_FILTER_TAGS": ""
  }
}

```

#### 各種キーの取得方法について
* **DISCORD_WEBHOOK_URL**: 通知先のDiscordサーバーのチャンネル設定 ➔ 連携サービス ➔ ウェブフック から作成・取得してください。
* **OPENAI_API_KEY**: Microsoft AI Foundry等でデプロイしたモデルのAPIキーを設定してください。
* **ITAD_API_KEY**: IsThereAnyDealの[一般アプリ登録画面](https://isthereanydeal.com/apps/)からユーザー登録し、APIキーを取得してください。
* **TWITCH_CLIENT_ID / SECRET**: IGDBの利用にはTwitch Developerアカウントが必要です。[公式開発者ドキュメント](https://api-docs.igdb.com/#getting-started)の手順に従ってアプリケーションを登録し、クライアントIDとシークレットを取得してください。


### 3. ストレージエミュレータ (Azurite) の起動
Timer Triggerの実行ロックを管理するため、バックグラウンドでAzuriteを起動します。

コマンドパレットから `Azurite: Start` を実行してください。

### 4. アプリケーションのビルドと起動
Goをビルドし、Functionsを起動します。

#### Windows の場合
```bash
func start
```