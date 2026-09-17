# Adlaire CI — Runner 詳細仕様

本ファイルは `ADLAIRE_CI_DETAIL_SPEC.md` から分割した `runner` owner component の詳細仕様である。

本ファイルに、方針、ポリシー、実装状態、正本関係、ロードマップ状態、実装可否の上位判断を記載してはならない。これらは `ADLAIRE_CI_SPEC.md` を正とする。

本ファイルを読む前に、`ADLAIRE_CI_SPEC.md` で実装状態と実装可否を確認し、`ADLAIRE_CI_DETAIL_SPEC.md` §0〜§0j で共通固定値、責務 component、詳細節対応表、リポジトリ内ソース配置を確認する。本ファイルは `runner` owner component の主本文であり、collaborator component の仕様は呼び出し境界、schema、setup、security、fixture、検証観点として参照する。

`ADLAIRE_CI_DETAIL_SPEC.md` §27.38a は runner / builder / api / sdk / ui にまたがる横断補足契約であり、本ファイルへ移動しない。runner 拡張機能を実装する場合は、本ファイルの個別節を正本とし、横断する処理順、状態ファイル保存責務、api / sdk / ui 連動条件、受け入れ fixture の同期確認として `ADLAIRE_CI_DETAIL_SPEC.md` §27.38a を同時に確認する。§27.38a は本ファイルの個別節を上書きせず、§27.38a の内容を本ファイルへ重複定義してはならない。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `runner` |
| collaborator component | `builder`、`statefile`、`commitstatus`、`api`、`archive` |
| 持つ内容 | `runner` owner が主本文として定義する GitHub 監視、設定読取、状態ファイル更新呼び出し、pipeline、deploy、snapshot 作成トリガー、通知、runner fixture、runner owner 追加機能。 |
| 持たない内容 | API endpoint の認証・応答本文、SDK method 実装、UI DOM 詳細、builder の変換処理、admin 静的配信、security 主本文、状態 schema、setup / release 手順、fixture / PR 証跡正本。 |

---

### — CI ランナー —

## 10. CI ランナー 要件

| 項目 | 内容 |
|------|------|
| Go バージョン | Go `1.22` 以上。 |
| 外部依存 | **なし** — Go 標準ライブラリ（`net/http`、`encoding/json`、`os`、`os/exec`、`log/slog` 等）を使用する |
| 対象 OS | Linux（systemd 対応環境） |
| ネットワーク | サーバーから `api.github.com` への HTTPS 送信のみ |

---

## 10a. CI ランナー 実装対象

本節は、`runner` として実装する CI ランナー機能を定義する。

`runner` は `adlaire-ci-runner` バイナリとして実行する。起動形式は systemd timer から呼び出される oneshot 実行とし、1 回の起動で対象ブランチ設定を読み込み、変更検出、ビルド起動、ログ保存、通知、転送、後処理を完了して終了する。

実装時は、対象項目ごとに §0c の実装前確認項目を満たしていることを確認する。未充足の項目が 1 つでもある場合は、実装を開始せず、先に本ファイルの該当節を改訂する。

| 項目 | 関連節 | 実装内容 |
|------|--------|------------|
| JSON 形式の SHA キャッシュ | §11〜§13 | 各ターゲットの `sha_file` を `{"sha": "..."}` JSON 形式で読み書きする。 |
| `BRANCH_TARGETS` | §12〜§13 | 複数ブランチ、複数 target file、複数出力先を 1 つの設定リストとして処理する。 |
| GitHub API リトライ | §12〜§13 | `API_RETRY_MAX`、`API_RETRY_BASE_SECONDS`、指数バックオフ、レート制限待機を実装する。 |
| ビルドロック | §11〜§13 | `.build_lock` による多重起動防止を実装する。 |
| ビルドクールダウン | §12〜§13 | `.server_config.build_cooldown_seconds` による起動抑制を実装する。 |
| 強制再ビルド間隔 | §12〜§13 | `.server_config.force_build_interval_hours` による変更なし時の定期強制ビルドを実装する。 |
| コミット情報記録 | §13 | ビルドトリガー commit の SHA、message、author、date をビルドログへ記録する。 |
| 事前チェック | §13 | ディスク空き容量、`adlaire-ci-build` 実行可否、Go 版ビルドバイナリ配置を確認する。 |
| Webhook 通知 | §13 | `.notify_config` 読み込み、成功/失敗/転送失敗/週次サマリー通知、`.notify_pending` 再送を実装する。 |
| ビルドログ保存 | §11〜§15 | `.build_logs/{id}.json` へ stdout/stderr、変換レポート、所要時間を保存する。 |
| SSH 転送 | §14a | SHA256 差分検出、stdin パイプ転送、転送後整合性検証、ペンディングキューを実装する。 |
| スナップショット | §14b | `.snapshots/` への成果物保存、世代管理、ロールバック連携を実装する。 |
| サーキットブレーカー | §13 | `.build_circuit_state` による連続失敗停止を実装する。 |
| ログ世代管理 | §13 | `LOG_KEEP_N` による `.build_logs/` 削除を実装する。 |
| 出力サイズ警告 | §12〜§13 | `OUTPUT_SIZE_WARN_MB` による WARN ログと `size_warn` 記録を実装する。 |

### 初期実装対象外の連携範囲

§10〜§20 には、`runner` 単体の責務ではなく管理 API、標準管理ツール、追加の運用機能と結合して成立する項目が含まれる。これらは、API・SDK・UI の対象節に、呼び出し元、呼び出し先、状態ファイル、失敗時応答、検証条件が定義されるまで `runner` 単体で実装しない。

| 項目 | 理由 |
|------|------|
| API 経由の動的ブランチ設定 | `runner` 単体では設定 API を持たないため、`api` 実装と合わせて扱う。 |
| API 経由のロールバック | `POST /api/history/{id}/rollback` は `api` のエンドポイント実装が前提となる。 |
| 管理画面からのスケジュール操作 | systemd timer の変更 API と標準管理ツール UI が前提となる。 |

---

## 11. CI ランナー ファイル構成

本節のファイル構成は、Go 版 CI ランナーで使用するファイル、管理 API / sdk / ui 側のファイル、出力先、systemd ファイルを分離して示す。

### Go 版 CI ランナーで使用するファイル

| パス | 用途 |
|------|------|
| `/usr/local/bin/adlaire-ci-runner` | `components/runner.go` から生成する CI ランナーバイナリ。標準配置への移行完了前は、現行実装実体 `runner.go` から生成する同等バイナリとして扱う。 |
| `/usr/local/bin/adlaire-ci-build` | `components/builder.go` から生成する Markdown → 静的 Web サイトビルドバイナリ。標準配置への移行完了前は、現行実装実体 `build_spec.go` から生成する同等バイナリとして扱う。 |
| `/opt/adlaire-builder/.github_token` | GitHub PAT。`runner` が読み込む。 |
| `/opt/adlaire-builder/.last_sha` | 前回取得した blob SHA。JSON 形式で保存する。 |
| `/opt/adlaire-builder/repo/docs/` | GitHub Blobs API から取得した Markdown の書き出し先。単一 Markdown の場合も本ディレクトリ内へ保存する。 |
| `/opt/adlaire-builder/repo/.ci/pipeline.sh` | `runner` が `bash` で起動するビルド手順。 |

```
/opt/adlaire-builder/
├── .github_token
├── .last_sha
└── repo/
    ├── docs/
    └── .ci/
        └── pipeline.sh
```

### CI ランナー状態ファイル

以下は §10a の実装対象に対応するファイルである。`runner` は本仕様に従って作成・読み書きする。

| パス | 用途 |
|------|------|
| `/opt/adlaire-builder/.build_history` | ビルド履歴。 |
| `/opt/adlaire-builder/.notify_config` | Webhook 通知設定。runner は起動時整合性チェックと送信読み込みを行い、管理 API は設定更新を行う。 |
| `/opt/adlaire-builder/.notify_log` | Webhook 送信履歴。 |
| `/opt/adlaire-builder/.notify_pending` | Webhook 通知失敗時の再送キュー。 |
| `/opt/adlaire-builder/.pending_transfers` | SSH 転送失敗時の再送キュー。 |
| `/opt/adlaire-builder/.build_lock` | 実行中ビルドの PID ロック。 |
| `/opt/adlaire-builder/.branch_config` | ブランチターゲット設定。永続 JSON key は `branch_targets` とする。 |
| `/opt/adlaire-builder/.build_state` | ビルド実行状態、週次サマリー送信日等。 |
| `/opt/adlaire-builder/.build_status.json` | runner 現在状態と直近結果の要約。api / ui / mcp の read-only 参照元。 |
| `/opt/adlaire-builder/.build_circuit_state` | サーキットブレーカー状態。 |
| `/opt/adlaire-builder/.local_watch_state.json` | ローカルファイル監視モードの SHA-256 snapshot。 |
| `/opt/adlaire-builder/.build_cache.json` | ビルドキャッシュの manifest。 |
| `/opt/adlaire-builder/.build_cache/` | ビルドキャッシュの page fragment 保存先。 |
| `/opt/adlaire-builder/.dependency_manifest.json` | Markdown 入力と依存ファイルの対応 manifest。 |
| `/opt/adlaire-builder/.approval_queue` | ビルド承認フローの保留 queue。 |
| `/opt/adlaire-builder/.build_trends.json` | build 所要時間 trend と異常検知基準。 |
| `/opt/adlaire-builder/.build_chain_config` | build job 依存チェーン設定。 |
| `/opt/adlaire-builder/.build_logs/` | ビルドごとの個別ログ。 |
| `/opt/adlaire-builder/.snapshots/` | ビルド成果物スナップショット。 |

### 管理 API / sdk / ui 側ファイル

以下は `api`、`sdk`、`ui` の仕様に属する。CI ランナー拡張と連携するものを含むが、`runner` 単体の実装対象範囲には含めない。

```
/opt/adlaire-builder/
├── .admin_credentials   # 認証情報ファイル（JSON、パーミッション 600）
├── .server_config       # サーバー設定（JSON）
├── .access_log          # ログイン履歴（JSON）
├── .repo_config         # リポジトリ監視設定（JSON）
├── .config_log          # 設定変更履歴（JSON）
├── .access_control      # IP アクセス制限設定（JSON）
├── .hooks               # ビルド前後フック設定（JSON）
├── .maintenance         # メンテナンスモード状態（JSON）
├── .api_tokens          # API トークン管理（JSON、トークン本体はハッシュのみ保存）
├── .alert_rules         # カスタムアラートルール（JSON）
├── .tag_rules           # 自動タグ付けルール（JSON）
├── .pipeline_config     # パイプライン設定（JSON）
├── .notes               # 運用ノート（テキスト）
├── .smtp_config         # SMTP 設定（JSON、パスワード除く）
├── .smtp_secret         # SMTP パスワード（プレーンテキスト、パーミッション 600）
├── .dashboard_layout    # ダッシュボードウィジェットレイアウト（JSON）
├── .webhook_events.json # Webhook 受信イベントログ（JSON Lines 形式、1行1イベント）
├── .approval_queue      # ビルド承認フロー保留 queue（JSON Lines）
├── .build_trends.json   # ビルド時間 trend（JSON）
├── .build_chain_config  # ビルド依存チェーン設定（JSON）
└── admin/
    ├── index.html           # 管理画面（単一ファイル完結）
    └── adlaire-ci-sdk.js    # JavaScript SDK（管理画面に同梱）
```

管理 API サーバーの実行ファイルは `/usr/local/bin/adlaire-ci-api` とし、`components/api.go` から生成する。

### 出力先・配信先

```
/opt/adlaire-builder/dist/
└── site/                # 静的 Web サイト出力先

/var/www/html/           # SSH 転送先
```

### systemd

以下は runner が起動される運用上の配置である。unit 本文、配置、enable、restart、更新、rollback は `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 を正とする。

```
/etc/systemd/system/
├── adlaire-ci.service      # systemd ユニット（oneshot）
├── adlaire-ci.timer        # systemd タイマー（定期実行）
└── adlaire-ci-api.service  # 管理 API サーバー（常駐）
```

### リポジトリ側

```
<repo>/
├── docs   # ソース Markdown（GitHub 上のマスター）
└── .ci/
    └── pipeline.sh      # ビルド手順定義（実行権限付き）
```

---

## 12. 設定値（`runner`）

`runner` は本節の設定値を正とする。設定値は Go 構造体の既定値、設定ファイル、または CLI 引数で与える。どの入力経路を採用する場合でも、内部表現は本節のキー名・型・既定値に従う。

**関連型：**

```go
type RunnerConfig struct {
    StateDir                    string
    PendingFile                 string
    APIRetryMax                 int
    APIRetryBaseSeconds         int
    BuildCooldownSeconds        int
    HistoryKeepN                int
    ForceBuildIntervalHours     int
    LogKeepN                    int
    APICircuitBreakerThreshold  int
    OutputSizeWarnMB            int
    WeeklySummaryEnabled        bool
    WeeklySummaryDay            int
    WeeklySummaryHour           int
    BranchTargets               []BranchTarget
}

type BranchTarget struct {
    Branch        string
    TargetFile    string
    SHAFile       string
    Src           string
    Out           string
    DeployTargets []DeployTarget
}

type DeployTarget struct {
    Host    string
    User    string
    DestDir string
}
```

**設定値検証：**

| 項目 | 条件 | 不正時 |
|------|------|--------|
| `StateDir` | 絶対パス。空文字不可。 | 終了コード `2` |
| `PendingFile` | 絶対パス。空文字不可。 | 終了コード `2` |
| `APIRetryMax` | 0 以上。 | 終了コード `2` |
| `APIRetryBaseSeconds` | 0 以上。 | 終了コード `2` |
| `BuildCooldownSeconds` | 0 以上。 | 終了コード `2` |
| `HistoryKeepN` | 0 以上。 | 終了コード `2` |
| `ForceBuildIntervalHours` | 0 以上。 | 終了コード `2` |
| `LogKeepN` | 0 以上。 | 終了コード `2` |
| `APICircuitBreakerThreshold` | 0 以上。 | 終了コード `2` |
| `OutputSizeWarnMB` | 0 以上。 | 終了コード `2` |
| `WeeklySummaryDay` | 0〜6。 | 終了コード `2` |
| `WeeklySummaryHour` | 0〜23。 | 終了コード `2` |
| `BranchTargets` | 1 件以上。 | 終了コード `2` |
| `BranchTarget.Branch` | 空文字不可。`..`、`~`、制御文字禁止。 | 終了コード `2` |
| `BranchTarget.TargetFile` | 相対パス。絶対パス、`..`、先頭 `/` 禁止。 | 終了コード `2` |
| `BranchTarget.SHAFile` / `Src` / `Out` | 絶対パス。空文字不可。 | 終了コード `2` |
| `DeployTarget.Host` / `User` / `DestDir` | 空文字不可。`DestDir` は絶対パス。 | 当該 deploy target を無効として ERROR ログ。全 target 無効なら終了コード `2` |

`.server_config`、`.branch_config`、CLI 引数から読み込んだ値は `RunnerConfig` に正規化してから使用する。正規化後の `RunnerConfig` にないキーを処理フローで直接参照してはならない。

**設定入力の優先順位：**

1. CLI 引数
2. `.server_config` / `.branch_config` など状態ファイルの保存値
3. 本節の既定値

同一キーが複数の入力経路に存在する場合は、上位の値だけを採用する。採用しなかった値を混合してはならない。未知キーは WARN ログ `CONFIG_UNKNOWN_KEY: key={key}` を出して無視する。

**runner CLI 引数仕様：**

| 引数 | 必須 | 既定値 | 説明 |
|------|------|--------|------|
| `--state-dir <path>` | 任意 | `/opt/adlaire-builder` | 状態ファイル、repo、dist、admin の基準ディレクトリ。相対パスは禁止。 |
| `--once` | 任意 | `true` | 1 回だけ実行して終了する。Go 版 runner は oneshot 固定のため、指定してもしなくても同じ挙動とする。 |
| `--dry-run` | 任意 | `false` | 設定、状態、GitHub target、SHA 差分、起動可否だけを検証し、ビルド、deploy、通知、状態ファイル更新を行わず終了する。 |
| `--version` | 任意 | なし | バイナリ名、仕様名、Go build 情報を 1 行で標準出力へ表示して終了する。 |
| `--help` | 任意 | なし | 引数一覧を標準出力へ表示して終了する。 |

未知引数、値欠落、相対 `--state-dir` は終了コード `2` とし、ビルド処理を開始しない。

**runner CLI パース固定仕様：**

- 引数は `flag` package 互換の `--name value` と `--name=value` の両方を許可する。
- 短縮オプション（例：`-s`、`-o`）は禁止する。指定された場合は未知の引数として扱う。
- 同一引数が複数回指定された場合は最後の値を採用する。ただし `--once` は指定有無にかかわらず `true` として扱う。`--dry-run` は 1 回以上指定されれば `true` とする。
- `--help` と `--version` は他の引数より優先し、`.github_token` 読み込み、lock 作成、状態ファイル読み込み、GitHub API 呼び出しを行わない。
- stderr のエラー行は末尾に改行 1 つを付ける。複数エラーをまとめて出力せず、最初に検出したエラー 1 件で終了する。

**runner 設定正規化契約：**

runner は CLI、`.server_config`、`.branch_config`、既定値を読み込んだ後、処理開始前に 1 回だけ `RunnerConfig` へ正規化する。正規化前の map、JSON raw message、環境変数、CLI flag 値を、GitHub API、pipeline、deploy、snapshot、通知処理から直接参照してはならない。

| 入力 | 正規化先 | 正規化ルール |
|------|----------|--------------|
| `--state-dir` | `RunnerConfig.StateDir` | `filepath.Clean` 後も絶対パスであることを確認する。末尾 `/` の有無で別パス扱いしない。 |
| `PENDING_FILE` | `RunnerConfig.PendingFile` | CLI で `StateDir` が変更された場合、明示設定がない限り `{StateDir}/.pending_transfers` に再解決する。 |
| `.branch_config.branch_targets[]` | `RunnerConfig.BranchTargets` | 配列順を保持する。各 entry は `branch`、`target_file`、`sha_file`、`src`、`out`、`deploy_targets` だけを採用する。 |
| `BRANCH_TARGETS` 既定値 | `RunnerConfig.BranchTargets` | `.branch_config` が存在しない場合だけ使用する。`.branch_config` が破損復旧で不在化された場合も同じ既定値へ fallback する。 |
| `.server_config` の数値 | `RunnerConfig` の数値 field | JSON number が整数でない場合は設定不正とする。文字列数値の暗黙変換は禁止する。 |
| `.server_config` の boolean | `RunnerConfig` の boolean field | JSON boolean のみ許可する。`"true"`、`1`、`"yes"` は不正値とする。 |

正規化後は、全 path を絶対パス文字列として保持する。`BranchTarget.TargetFile` だけは GitHub repository 内の相対パスとして保持し、`filepath.Clean` 後に `.`、空文字、`..` を含む path、先頭 `/`、NUL byte、制御文字を禁止する。`BranchTarget.Src` と `BranchTarget.Out` が同一または親子関係になる場合は終了コード `2` とし、ERROR ログ `CONFIG_PATH_CONFLICT: src={src} out={out}` を出す。

**状態ディレクトリ構造検証契約：**

`--state-dir` 検証後、runner は以下を処理開始前に確認する。

| 対象 | 条件 | 不正時 |
|------|------|--------|
| `StateDir` | directory、owner が実行ユーザーまたは書き込み可能、mode に owner write がある。 | 終了コード `2`、ERROR `STATE_DIR_INVALID: path={path}`。 |
| `{StateDir}/repo` | 不在なら作成する。file の場合は不正。 | 作成失敗または file の場合、終了コード `2`。 |
| `{StateDir}/dist` | 不在なら作成する。file の場合は不正。 | 作成失敗または file の場合、終了コード `2`。 |
| `{StateDir}/.build_logs` | 不在なら `0700` で作成する。 | 作成失敗時は終了コード `2`。 |
| `{StateDir}/.snapshots` | 不在なら `0700` で作成する。ただし snapshot 無効時も directory 作成は許可する。 | 作成失敗時は終了コード `2`。 |

上記 directory 作成は dry-run では実行しない。dry-run では作成予定を `§27.2` の stdout JSON `warnings[]` に `code="DRY_RUN_WOULD_CREATE_DIR"` として出し、`would_write` に `"state_dir"` を追加してはならない。ただし既存 path が file の場合は dry-run でも終了コード `2` とし、stdout JSON `errors[]` に原因を出す。

**固定出力：**

| 条件 | stdout |
|------|--------|
| `--help` | `Usage: adlaire-ci-runner [--state-dir path] [--once] [--dry-run] [--version] [--help]` |
| `--version` | `adlaire-ci-runner ADLAIRE_CI_SPEC go=<runtime.Version()>` |

**CLI 異常系：**

| 条件 | 終了コード | stderr |
|------|------------|--------|
| 未知の引数 | `2` | `unknown option: <name>` |
| `--state-dir` 値欠落 | `2` | `missing value: --state-dir` |
| `--state-dir` が空文字 | `2` | `state directory must not be empty` |
| `--state-dir` が相対パス | `2` | `state directory must be absolute: <path>` |
| `--state-dir` が存在しない | `2` | `state directory not found: <path>` |
| `--state-dir` がディレクトリではない | `2` | `state path is not directory: <path>` |

以下の `BRANCH_TARGETS`、`PENDING_FILE`、`API_RETRY_MAX`、`.server_config.build_cooldown_seconds`、`HISTORY_KEEP_N`、`.server_config.force_build_interval_hours`、`LOG_KEEP_N`、`API_CIRCUIT_BREAKER_THRESHOLD`、`OUTPUT_SIZE_WARN_MB`、`WEEKLY_SUMMARY_*` を Go 版 runner の標準設定とする。

```text
PENDING_FILE           = "/opt/adlaire-builder/.pending_transfers"   # SSH 転送ペンディングキュー（JSON）
API_RETRY_MAX          = 5    # GitHub API 失敗時の最大再試行回数（指数バックオフ）
API_RETRY_BASE_SECONDS = 1    # バックオフ基底秒数（1→2→4→8→16 秒。0 = リトライ無効）
server_config.build_cooldown_seconds = 60   # 前回ビルド完了から次ビルドまでの最小間隔（秒。0 = 無効）→ §13
HISTORY_KEEP_N         = 10   # スナップショット保持世代数（0 = 無制限）→ §14b
server_config.force_build_interval_hours = 0 # 強制再ビルド間隔（時間。0 = 無効）→ §13
LOG_KEEP_N             = 50   # ビルドログ保持件数（0 = 無制限）→ §13
API_CIRCUIT_BREAKER_THRESHOLD = 3    # 全ブランチ連続失敗の許容周回数（0 = 無効）→ §13
OUTPUT_SIZE_WARN_MB           = 5    # 出力サイト合計サイズ警告閾値（MB。0 = 無効）→ §13・§8
WEEKLY_SUMMARY_ENABLED        = true # 週次サマリー Webhook の有効/無効 → §13
WEEKLY_SUMMARY_DAY            = 0    # 送信曜日（0=月曜〜6=日曜） → §13
WEEKLY_SUMMARY_HOUR           = 9    # 送信時刻（0〜23、ローカル時刻） → §13

# ブランチターゲット設定（→ §14a）
# 複数エントリを定義した場合はリスト順に順次処理する（並列処理は対象外）
BRANCH_TARGETS = [
    {
        "branch":      "main",                                        # 監視対象ブランチ
        "target_file": "docs",                                        # 監視対象 Markdown ファイルまたはディレクトリ
        "sha_file":    "/opt/adlaire-builder/.last_sha",             # blob SHA キャッシュ（JSON 形式: {"sha": "..."}）
        "src":         "/opt/adlaire-builder/repo/docs",              # Blobs API 書き出し先
        "out":         "/opt/adlaire-builder/dist/site",              # ビルド成果物ディレクトリ
        "deploy_targets": [
            {
                "host":     "<配信サーバーIP>",                       # SSH 転送先ホスト
                "user":     "deploy",                                 # SSH 接続ユーザー
                "dest_dir": "/var/www/html",                          # 転送先ディレクトリ
            }
        ],
    }
]
```

`BRANCH_TARGETS` が空の場合、`runner` は ERROR ログを出力し、ビルドを実行せず終了コード `2` で終了する。

**runner 終了コード：**

| 終了コード | 条件 |
|------------|------|
| `0` | 起動、ロック確認、対象処理が完了した。変更なし、クールダウン、既存ロックによるスキップも正常終了に含める。 |
| `1` | 1 件以上のビルドまたは転送が失敗したが、runner 自体は最後まで処理できた。 |
| `2` | 設定不正、必須ファイル不足、CLI 引数不正。 |
| `3` | GitHub API など外部サービスへの全再試行が失敗し、全ターゲットが処理不能。 |
| `4` | `.build_lock` 作成に失敗し、既存 PID の実行中確認もできない。 |

systemd timer からの再実行を妨げないため、終了コード `1` と `3` でもロック削除、ログ保存、通知キュー保存を試行してから終了する。

**atomic write 共通契約：**

runner が JSON object、JSON array、SHA cache、`.build_state`、`.build_circuit_state`、`.notify_pending`、`.pending_transfers`、`.server_config`、`.branch_config` を更新する場合は、以下の順序で atomic write を行う。

1. 対象ファイルと同じディレクトリに `.{basename}.tmp.{pid}` を作成する。
2. JSON は末尾改行付き UTF-8 として書き込む。
3. `file.Sync()` を実行してから close する。
4. `os.Rename(tmp, target)` で置換する。
5. 親ディレクトリを open できる場合は directory sync を実行する。directory sync が `EINVAL` 等で未対応の場合は WARN ログ `DIR_SYNC_UNSUPPORTED: path={dir}` を出し、処理は成功扱いとする。

atomic write 失敗時は対象ファイルを更新済みとして扱わない。tmp ファイルが残った場合は削除を試行し、削除失敗時は WARN ログ `TMP_CLEANUP_FAILED: path={tmp}` を出す。

**lock ファイル契約：**

`.build_lock` は UTF-8 text で、内容は `pid={pid}\nstarted_at={UTC_ISO8601}\n` とする。

| 状態 | 処理 |
|------|------|
| lock なし | `os.OpenFile(path, O_CREATE|O_EXCL|O_WRONLY, 0600)` で作成する。 |
| lock あり・PID 実行中 | INFO ログ `BUILD_SKIP: already running (PID {pid})` を出し終了コード `0`。 |
| lock あり・PID 不在 | WARN ログ `STALE_LOCK: pid={pid}` を出し、lock を削除してから再作成する。 |
| lock あり・PID 解析不能 | ERROR ログ `LOCK_CORRUPT: path={path}` を出し終了コード `4`。 |
| lock 作成失敗 | ERROR ログ `LOCK_CREATE_FAILED: {reason}` を出し終了コード `4`。 |

PID 実行中判定は Linux の `/proc/{pid}` 存在確認で行う。`/proc` を読めない場合は PID 実行中確認不能として終了コード `4` とする。

**GitHub token 読み込み契約：**

`.github_token` は `StateDir` 直下の通常ファイルだけを認める。symbolic link、directory、device file、FIFO は禁止する。

| 条件 | 処理 |
|------|------|
| ファイル不在 | ERROR `GITHUB_TOKEN_MISSING: path={path}`、終了コード `2`。 |
| 読み込み権限なし | ERROR `GITHUB_TOKEN_PERMISSION: path={path}`、終了コード `2`。 |
| mode が `0600` より広い | ERROR `GITHUB_TOKEN_INSECURE_MODE: path={path} mode={mode}`、終了コード `2`。 |
| UTF-8 不正 | ERROR `GITHUB_TOKEN_INVALID: reason=utf8`、終了コード `2`。 |
| trim 後空文字 | ERROR `GITHUB_TOKEN_INVALID: reason=empty`、終了コード `2`。 |
| trim 後に改行、空白、NUL、制御文字を含む | ERROR `GITHUB_TOKEN_INVALID: reason=character`、終了コード `2`。 |

token は `strings.TrimSpace` 後の値だけを HTTP Authorization header に使用する。token の値、先頭文字、末尾文字、長さ、hash は stdout、stderr、`.build_logs/{id}.json`、`.build_history`、`.notify_pending`、`.notify_log`、snapshot、fixture expected output に保存してはならない。secret mask は token 読み込み成功直後に登録し、以降の全ログ保存処理より前に適用する。

**runner 状態ファイル参照契約：**

状態ファイルの path、形式、初期値、schema、破損時の扱い、atomic write、adapter、読取 priority は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a、§22.0c を正とする。本ファイルでは、runner がどの処理段階で状態を読むか、いつ更新するか、失敗時に後続処理を止めるかだけを定義する。

runner 実装は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` に未定義の状態ファイル、永続 key、queue entry key、notification entry key、pending transfer entry key を追加してはならない。

`.branch_config` が存在しない場合は §12 の `BRANCH_TARGETS` 既定値を使用する。存在する場合の schema、空配列の扱い、永続 key、API 表示名との境界は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0c の `.branch_config` schema を正とする。

runner は起動時の設定正規化で `.branch_config` を 1 回だけ読み、正規化後の `RunnerConfig.BranchTargets` を当該起動中の唯一の branch target 情報として使用する。同一 runner 起動中に `.branch_config` を再読込してはならない。API による `.branch_config` 更新、削除、default 復帰は、既に実行中の runner には反映せず、次回 runner 起動から反映する。

`.pending_transfers` entry は §14a の形式と `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0c の状態 schema を同時に満たす。JSON array 内の entry は投入順を保持し、再試行も投入順で処理する。重複統合は `out`、`host`、`user`、`dest_dir` の 4 項目完全一致で判定する。

**SHA cache 読み書き契約：**

`sha_file` は target ごとの処理済み Git blob SHA を保存する JSON file である。runner は legacy text 形式を自動変換してはならない。

| 状態 | 処理 |
|------|------|
| 不在 | 初回実行として `previous_blob_sha=""` を扱う。build 成功時に `{"sha":"<new_sha>"}` を atomic write する。 |
| `{"sha":""}` | 初回実行として扱う。 |
| `{"sha":"<value>"}` | `<value>` を前回 SHA として比較する。 |
| JSON 破損 | `failure_state_write` ではなく `failure_decode` として当該 target を失敗扱いし、sha_file を更新しない。 |
| object 以外 | JSON 破損と同じ扱い。 |
| `sha` key 不在または string 以外 | JSON 破損と同じ扱い。 |
| 未知 key あり | 未知 key を除去し、build 成功時に `sha` だけの object で上書きする。 |

SHA cache の更新は、pipeline 成功後、deploy 前に行う。複数 target のうち一部 target が成功した場合は、成功 target の `sha_file` だけを更新する。失敗 target、skip target、branch target 設定不正 target の `sha_file` を更新してはならない。

**状態ファイル権限契約：**

runner が新規作成する状態ファイルは、JSON object、JSON array、SHA cache、lock file、pending transfer file、build log、build history の分類に関係なく `0600` とする。directory は `0700` とする。既存ファイルの mode が広い場合、secret を含む `.github_token`、`.notify_config`、`.notify_pending`、`.pending_transfers` は停止条件とし、それ以外の runner 状態ファイルは WARN `STATE_FILE_INSECURE_MODE: path={path} mode={mode}` を出して `0600` へ chmod する。chmod 失敗時は終了コード `2` とする。

**設定ファイル起動時整合性チェック：**

本機能の目的は、runner 起動時に状態ファイルの破損、型不一致、必須 key 不足、権限不備を検出し、ビルド処理開始前に復旧または停止することである。owner component は `runner` とし、collaborator component は `statefile` とする。管理 API の HTTP endpoint、sdk、ui は本機能の実行責務を持たない。

対象ファイルは次の 6 件に固定する。実装者判断で対象ファイルを追加または除外してはならない。

| 順序 | ファイル | 必須性 | 正常時の扱い | 不在時の扱い |
|------|----------|--------|--------------|--------------|
| 1 | `.branch_config` | 任意 | JSON object として検証し、`branch_targets` を `RunnerConfig.BranchTargets` へ正規化する。 | ファイルを作成せず、§12 の `BRANCH_TARGETS` 既定値を使用する。 |
| 2 | `.notify_config` | 任意 | JSON object として検証し、通知送信時の設定として使用する。 | 初期値を atomic write で作成する。 |
| 3 | `.build_state` | 必須 | JSON object として検証し、`queued`、`last_finished_at`、週次サマリー状態を保持する。 | 初期値を atomic write で作成する。 |
| 4 | `.build_circuit_state` | 必須 | JSON object として検証し、circuit breaker 状態を保持する。 | 初期値を atomic write で作成する。 |
| 5 | `.pending_transfers` | 必須 | JSON array として検証し、entry の投入順を保持する。 | `[]` を atomic write で作成する。 |
| 6 | `.notify_pending` | 必須 | JSON array として検証し、entry の投入順を保持する。 | `[]` を atomic write で作成する。 |

実行順序は固定とする。

1. CLI 引数を検証し、`--help` または `--version` の場合は本チェックを実行しない。
2. `StateDir` が絶対パスかつ既存ディレクトリであることを確認する。
3. `.build_lock` を取得する。実行中 PID がある場合は本チェックを実行せず終了コード `0` で終了する。
4. 本チェックを上表の順序で実行する。
5. 本チェックが復旧可能な問題のみで完了した場合は、`.github_token` 読み込み、pending retry、cooldown、target 処理へ進む。
6. 本チェックが停止条件に該当した場合は、`.build_state.running` を `true` にせず、`.build_logs/{id}.json` と `.build_history` を作成せず、`.build_lock` を削除して終了する。

判定分類は次のとおり固定する。

| 分類 | 条件 | 処理 |
|------|------|------|
| `missing_optional` | `.branch_config` が存在しない。 | ファイルを作成せず、`BRANCH_TARGETS` 既定値を採用する。ログは出さない。 |
| `missing_required` | `.notify_config`、`.build_state`、`.build_circuit_state`、`.pending_transfers`、`.notify_pending` が存在しない。 | 初期値を atomic write で作成し、WARN ログ `CONFIG_INIT: path={path}` を出す。 |
| `empty_file` | ファイルサイズ 0 byte。 | `parse_error` と同じ扱い。 |
| `parse_error` | UTF-8 として読めない、または JSON parse に失敗する。 | corrupt backup へ退避し、ファイル別初期化を行う。 |
| `top_level_type_mismatch` | object 必須のファイルが array/string/null、array 必須のファイルが object/string/null。 | corrupt backup へ退避し、ファイル別初期化を行う。 |
| `required_key_missing` | 必須 key が存在しない。 | corrupt backup へ退避し、ファイル別初期化を行う。 |
| `required_key_type_mismatch` | 必須 key の型が schema と異なる。 | corrupt backup へ退避し、ファイル別初期化を行う。 |
| `invalid_value` | 値が許容範囲外。例: 負の `retry_count`、不正な ISO 8601、相対 `sha_file`。 | corrupt backup へ退避し、ファイル別初期化を行う。 |
| `unknown_key` | schema にない key が存在する。 | backup せず、未知 key を除去した正規化 JSON を atomic write し、WARN ログ `CONFIG_UNKNOWN_KEY: path={path} key={key}` を出す。 |
| `permission_error` | 読み込み、rename、chmod、親ディレクトリ sync、lock 作成のいずれかが権限エラー。 | 自動退避せず、ERROR ログ `CONFIG_PERMISSION_ERROR: path={path} op={op}` を出し、終了コード `2`。 |
| `io_error` | 権限以外の read/write/rename/sync 失敗。 | 自動退避せず、ERROR ログ `CONFIG_IO_ERROR: path={path} op={op}` を出し、終了コード `1`。 |

corrupt backup のファイル名は `{original}.corrupt.{YYYYMMDDHHMMSS}.bak` とする。timestamp は UTC、秒単位、ゼロ埋め固定とする。同一秒内に同じファイルの backup 名が衝突した場合は、2 件目以降を `{original}.corrupt.{YYYYMMDDHHMMSS}.{n}.bak` とし、`n` は `2` から始める。

ファイル別の復旧結果は次のとおり固定する。

| ファイル | 復旧時の書き込み内容 | 復旧後の処理継続 | 備考 |
|----------|----------------------|------------------|------|
| `.branch_config` | 書き込まない。破損元を backup した後、`.branch_config` は不在状態にする。 | 継続する。 | §12 の `BRANCH_TARGETS` 既定値へ fallback する。 |
| `.notify_config` | §22.0a の初期値。 | 継続する。 | このファイル自体が破損していた場合、`config_corrupt` 通知は送信しない。 |
| `.build_state` | §22.0a の初期値。 | 継続する。 | `queued` は失われるため、ERROR ログと通知対象に含める。 |
| `.build_circuit_state` | §22.0a の初期値。 | 継続する。 | circuit open 状態は解除されるため、ERROR ログと通知対象に含める。 |
| `.pending_transfers` | `[]`。 | 継続する。 | 未再送転送は失われるため、ERROR ログと通知対象に含める。 |
| `.notify_pending` | `[]`。 | 継続する。 | 未送信通知は失われるため、ERROR ログを出す。通知 queue 自体が失われるため追加 queue は行わない。 |

復旧時ログは次の 3 種に固定する。

| 条件 | ログ level | メッセージ |
|------|------------|------------|
| 初期作成 | WARN | `CONFIG_INIT: path={path}` |
| corrupt backup 成功 | ERROR | `CONFIG_CORRUPT_BACKUP: path={path} backup={backup_path} reason={reason}` |
| 正規化書き戻し | WARN | `CONFIG_NORMALIZED: path={path}` |

復旧通知は `.notify_config` の復旧と検証が完了した後に 1 回だけ送信する。送信条件は、復旧対象に `.notify_config` と `.notify_pending` 以外のファイルが 1 件以上含まれ、かつ `.notify_config.channels[]` または互換 `.notify_config.webhooks[]` のうち `enabled=true` で `on` に `"config_corrupt"` または `"*"` を含む宛先が存在する場合とする。payload は以下の JSON object に固定する。

```json
{
  "event": "config_corrupt",
  "reason": "config_corrupt",
  "recovered_at": "2026-09-16T00:00:00Z",
  "files": [
    {
      "path": "/opt/adlaire-builder/.build_state",
      "action": "backup_and_init",
      "backup": "/opt/adlaire-builder/.build_state.corrupt.20260916000000.bak",
      "reason": "parse_error"
    }
  ]
}
```

通知送信に失敗した場合は、`.notify_pending` が本チェックで破損していない場合に限り、`event="config_corrupt"` の entry を `.notify_pending` へ追記する。`.notify_pending` が本チェックで初期化された場合は、二重消失を避けるため追記せず、ERROR ログ `NOTIFY_FAILED: event=config_corrupt url={url}` のみ出す。

schema 検証では次を必須とする。

- `.branch_config` は `branch_targets` のみを永続 key とし、`branches` だけを持つファイルは `required_key_missing` として扱う。API request / response の alias は永続ファイルへ保存する前に `branch_targets` へ変換する。
- `.notify_config.channels[].on` と `.notify_config.webhooks[].on` は `start`、`success`、`failure`、`deploy_failure`、`weekly_summary`、`approval_required`、`duration_anomaly`、`config_corrupt`、`*` のみ許可する。
- `.build_state` は `weekly_summary_sent_date` を必須 key とする。値は `null` または `YYYY-MM-DD` とする。
- `.pending_transfers[]` は §14a の pending entry schema と一致すること。1 件でも不正 entry がある場合はファイル全体を corrupt として扱う。
- `.notify_pending[]` は `event`、`url`、`payload`、`queued_at`、`retry_count`、`last_error` を必須 key とする。`payload` は JSON object、`retry_count` は 0 以上の integer とする。

本チェックは冪等でなければならない。初期化または正規化済みの状態で runner を再起動した場合、追加 backup、追加通知、追加 WARN/ERROR は発生しない。同一破損ファイルが復旧失敗後に残っている場合だけ、次回起動時に再度同じ判定を行う。

検証 fixture は以下を必須とする。

| fixture | 入力状態 | 期待結果 |
|---------|----------|----------|
| `config-startup/missing-required` | 対象必須ファイルが存在しない。 | 初期値が作成され、終了コード `0`、ビルド処理へ進む。 |
| `config-startup/corrupt-build-state` | `.build_state` が `{bad json`。 | `.build_state.corrupt.{timestamp}.bak` へ退避、初期値作成、`config_corrupt` 通知対象、終了コード `0`。 |
| `config-startup/corrupt-branch-config` | `.branch_config` が JSON array。 | backup 後 `.branch_config` は不在、既定 `BRANCH_TARGETS` 採用、終了コード `0`。 |
| `config-startup/unknown-key` | `.build_circuit_state` に未知 key がある。 | backup なしで未知 key を除去、`CONFIG_NORMALIZED`、終了コード `0`。 |
| `config-startup/invalid-pending-entry` | `.pending_transfers` に必須 key 不足 entry がある。 | backup 後 `[]` 作成、終了コード `0`。 |
| `config-startup/permission-error` | 対象ファイルが読み込み不可。 | 自動退避なし、終了コード `2`、`.build_state.running` 未変更。 |
| `config-startup/help-version-skip` | `--help` または `--version`。 | 対象ファイルを読まず、変更しない。 |

完了条件は、上記 fixture を Go test で検証し、`ADLAIRE_CI_SPEC.md` の状態表、`DOCUMENT_INDEX.md` の実装状態、PR 本文の検証結果が一致していることとする。

**build id 契約：**

runner が生成する build id は UTC 時刻ベースの `b{YYYYMMDDHHmmss}` とする。同一秒内に複数 target のビルドログが必要な場合は、2 件目以降を `b{YYYYMMDDHHmmss}-2`、`-3` とする。build id は `.build_logs/{id}.json`、`.build_history`、`.snapshots/{id}/` で同一値を使用する。

---

## 13. 処理フロー

本節の処理フローは、`runner` の標準フローである。

**状態更新順序の規範：**

1. `.build_lock` を作成する。
2. 設定ファイル起動時整合性チェックを実行する。
3. `.build_status.json` に `status="running"`、`trigger`、`started_at`、`current_build_id` を atomic write で保存する。
4. `.build_state.running=true`、`current_build_id`、`last_started_at` を atomic write で保存する。
5. GitHub API、Blob 書き出し、事前チェック、`pipeline.sh` 実行を行う。
6. ビルド成功時のみ `sha_file` を新 SHA に更新する。
7. `.build_logs/{id}.json` を作成し、stdout/stderr、`[REPORT]`、警告、転送結果、`trigger` を保存する。
8. `.build_history` に同じ `id` の要約行を JSON Lines で追記する。
9. 転送成功後に `.snapshots/` を更新する。
10. `.build_status.json` に最終状態を atomic write で保存する。
11. `.build_state.running=false`、`last_finished_at` を保存する。
12. `.build_lock` を削除する。

途中失敗時は、失敗が発生した段階以降の成功前提更新を行わない。例えば `pipeline.sh` 失敗時は `sha_file`、snapshot、転送成功履歴を更新しない。ただし `.build_logs/{id}.json`、`.build_history`、`.build_state.running=false`、通知 pending は失敗記録として保存する。

**ターゲット結果分類：**

runner は `BRANCH_TARGETS` の各 entry について、最終的に次のいずれか 1 つの `target_status` を確定する。

| `target_status` | 条件 | runner 終了コードへの影響 |
|-----------------|------|---------------------------|
| `skipped_no_change` | SHA 一致かつ強制ビルド条件未達。 | 失敗扱いしない。 |
| `skipped_cooldown` | クールダウンで全体処理を開始しない。 | `0`。 |
| `success` | pipeline 成功、SHA 更新成功、deploy target がないか全 deploy target が成功。 | `0` 候補。 |
| `success_deploy_pending` | pipeline 成功、SHA 更新成功、1 件以上の deploy が pending。 | `1`。 |
| `failure_api` | GitHub API が全再試行失敗。 | 全 target がこれなら `3`、一部なら `1`。 |
| `failure_decode` | Blob Base64 decode または Markdown 書き出し失敗。 | `1`。 |
| `failure_precheck` | 事前チェック失敗。 | `1`。 |
| `failure_build` | `pipeline.sh` が非 0 終了または timeout。 | `1`。 |
| `failure_state_write` | SHA、ログ、履歴、状態ファイルの必須更新に失敗。 | `1`。 |

`BRANCH_TARGETS` が複数ある場合、runner は設定不正を除き、1 target の失敗で全体処理を中断しない。全 target 処理後、最も重い終了コードを採用する。重さは `4 > 3 > 2 > 1 > 0` とする。

**更新可否マトリクス：**

| 状況 | `sha_file` | `.build_logs/{id}.json` | `.build_history` | deploy | snapshot |
|------|------------|-------------------------|------------------|--------|----------|
| 変更なし | 更新しない | 作成しない | 追記しない | 実行しない | 作成しない |
| 強制ビルド成功 | 新 SHA または既存 SHA を保存 | 作成する | 追記する | 実行する | deploy 成功時のみ作成 |
| API 失敗 | 更新しない | 作成する | 追記する | 実行しない | 作成しない |
| Blob decode 失敗 | 更新しない | 作成する | 追記する | 実行しない | 作成しない |
| 事前チェック失敗 | 更新しない | 作成する | 追記する | 実行しない | 作成しない |
| pipeline 失敗 | 更新しない | 作成する | 追記する | 実行しない | 作成しない |
| pipeline 成功・deploy 成功 | 更新する | 作成する | 追記する | 実行する | 作成する |
| pipeline 成功・deploy pending | 更新する | 作成する | 追記する | pending へ投入 | 作成しない |

`sha_file` は pipeline 成功後、deploy 実行前に更新する。理由は、ビルド成果物生成が成功した時点で入力 SHA の処理は完了しており、deploy 失敗は `.pending_transfers` の責務で再試行するためである。

**runner 失敗段階別副作用固定契約：**

| 失敗段階 | `target_status` | 保存必須 | 保存禁止 | finalizer | 終了コード |
|----------|-----------------|----------|----------|-----------|------------|
| lock 実行中 PID | `lock_skipped` | lock が有効な実行中 PID を指す場合に限り、`.build_status.json` に `lock_skipped` を保存する。 | `.build_state`、`.build_logs/`、`.build_history`、`sha_file`、deploy、snapshot。 | lock を削除しない。 | `0` |
| startup config permission error | `config_error` | `.build_status.json` に `config_error`、ERROR log。 | `.build_state.running=true`、`.build_logs/`、`.build_history`、`sha_file`、deploy、snapshot。 | lock 作成済みなら削除する。 | `2` |
| status start write failure | `failure_state_write` | ERROR log。 | `.build_state.running=true`、`.build_logs/`、`.build_history`、`sha_file`、deploy、snapshot。 | lock 作成済みなら削除する。 | `1` |
| state start write failure | `failure_state_write` | `.build_status.json` に `failure`、ERROR log。 | `.build_logs/`、`.build_history`、`sha_file`、deploy、snapshot。 | lock を削除する。 | `1` |
| GitHub tree / blob API failure | `failure_api` | `.build_logs/{id}.json`、`.build_history`、`.build_status.json`、`.build_state.running=false`。 | `sha_file`、deploy、snapshot、materialized src の確定置換。 | 実行する。 | 全 target 失敗なら `3`、一部なら `1` |
| blob decode / materialize failure | `failure_decode` | `.build_logs/{id}.json`、`.build_history`、`.build_status.json`、`.build_state.running=false`。 | `sha_file`、pipeline、deploy、snapshot。 | 実行する。 | `1` |
| precheck failure | `failure_precheck` | `.build_logs/{id}.json`、`.build_history`、`.build_status.json`、`.build_state.running=false`。 | `sha_file`、pipeline、deploy、snapshot。 | 実行する。 | `1` |
| pipeline non-zero / timeout | `failure_build` | `.build_logs/{id}.json`、`.build_history`、`.build_status.json`、`.build_state.running=false`。 | `sha_file`、deploy、snapshot。 | 実行する。 | `1` |
| build log write failure | `failure_state_write` | `.build_status.json`、`.build_state.running=false`、ERROR log。 | `.build_history`、`sha_file`、deploy、snapshot。 | 実行する。 | `1` |
| build history append failure | `failure_state_write` | `.build_logs/{id}.json`、`.build_status.json`、`.build_state.running=false`、ERROR log。 | `sha_file`、deploy、snapshot。 | 実行する。 | `1` |
| SHA write failure | `failure_state_write` | `.build_logs/{id}.json`、`.build_history`、`.build_status.json`、`.build_state.running=false`。 | deploy、snapshot。 | 実行する。 | `1` |
| deploy failure | `success_deploy_pending` | `sha_file`、`.pending_transfers`、`.build_logs/{id}.json`、`.build_history`、`.build_status.json`、`.build_state.running=false`。 | snapshot。 | 実行する。 | `1` |
| snapshot failure | `success` | `sha_file`、deploy 成功、`.build_logs/{id}.json` に WARN、`.build_history`、`.build_status.json`、`.build_state.running=false`。 | snapshot 成功扱い、`.pending_transfers` 追加。 | 実行する。 | `0` |
| status finalizer write failure | 実行結果に従う | `.build_logs/{id}.json`、`.build_history`、`.build_state.running=false`、ERROR log。 | 部分 `.build_status.json`。 | 継続する。 | 最低 `1` |
| build_state finalizer write failure | 実行結果に従う | `.build_logs/{id}.json`、`.build_history`、`.build_status.json`、ERROR log。 | 正常終了扱い。 | lock 削除を試みる。 | `1` |

上表の保存必須に含まれる状態ファイルは、保存失敗時に `failure_state_write` へ分類する。ただし status finalizer と build_state finalizer は、既に確定した target の log / history を取り消さない。保存禁止に含まれる処理を実行した場合は仕様違反とし、実装 PR の fixture で失敗として扱う。

**複数 target 継続 / 中断固定契約：**

| 条件 | 継続可否 | 次 target への影響 |
|------|----------|--------------------|
| 1 target の `failure_api`、`failure_decode`、`failure_precheck`、`failure_build` | 継続する。 | 当該 target の `sha_file` は更新せず、次 target は通常判定する。 |
| 1 target の `success_deploy_pending` | 継続する。 | `.pending_transfers` に当該 target を保存し、次 target は通常判定する。 |
| 1 target の snapshot failure | 継続する。 | 当該 target は `success` とし、次 target は通常判定する。 |
| `.build_logs/{id}.json` 保存失敗 | 継続しない。 | 同一 runner 起動内の後続 target を開始しない。 |
| `.build_history` 追記失敗 | 継続しない。 | 同一 runner 起動内の後続 target を開始しない。 |
| `.build_status.json` start 保存失敗 | 継続しない。 | `.build_state.running=true` へ進まない。 |
| `.build_state.running=true` 保存失敗 | 継続しない。 | target 処理へ進まない。 |
| `.github_token` 不備 / mode 不正 | 継続しない。 | GitHub API、pipeline、deploy を一切実行しない。 |
| `.branch_config` 全体破損かつ復旧不能 | 継続しない。 | target を推測して実行しない。 |

複数 target の終了コードは、処理済み target の最大重大度で決める。重大度は `4`（lock 形式不正など実行継続不能） > `3`（全 target GitHub API 失敗） > `2`（設定・secret・権限不正） > `1`（target 失敗または deploy pending） > `0`（成功または通常 skip）とする。`failure_api` が一部 target だけの場合は `1`、全処理 target が `failure_api` の場合だけ `3` とする。

**runner 機能単位契約：**

`runner` は、下表の機能単位で状態を更新する。各機能単位は、Write 列にない状態ファイルを更新してはならない。

| 機能単位 | Read | Write | 成功条件 | 失敗時更新 |
|----------|------|-------|----------|------------|
| lock acquisition | `.build_lock` | `.build_lock` | lock file を `O_CREATE|O_EXCL` で作成し、PID と started_at を書き込む。 | 実行中 PID がある場合は終了コード `0`。形式不正 lock は終了コード `4`、上書き禁止。 |
| startup config integrity | `.branch_config`, `.notify_config`, `.build_state`, `.build_circuit_state`, `.pending_transfers`, `.notify_pending` | 同左、corrupt backup | 対象 JSON の parse、top-level type、必須 key、型、値範囲、未知 key 正規化が完了する。 | 復旧可能な破損は backup と初期化後に継続。権限エラーは終了コード `2`、IO エラーは終了コード `1`。 |
| status start | `.build_state`, `.build_circuit_state`, `.pending_transfers`, `.notify_pending` | `.build_status.json` | `status="running"` と当該起動の `trigger` を保存する。 | 書き込み失敗は終了コード `1`。`.build_state.running=true` へ進まない。 |
| state start | `.build_state` | `.build_state` | `running=true`、`current_build_id`、`last_started_at` を保存する。 | `.build_lock` を削除し、終了コード `1`。 |
| pending transfer retry | `.pending_transfers` | `.pending_transfers`, `.build_logs/{id}.json` | 成功 entry を削除し、失敗 entry は `retry_count` を +1 して保持する。 | queue ファイル破損時は §22.0a に従い退避して `[]` で再生成する。 |
| notify pending retry | `.notify_pending`, `.notify_config` | `.notify_pending`, `.notify_log` | HTTP 2xx の entry を削除し、失敗 entry は `retry_count` を +1 して保持する。 | `.notify_config` 不在時は送信せず entry を保持し、ERROR ログを記録する。 |
| cooldown gate | `.build_state`, `.server_config` | none | cooldown 範囲外なら target 処理へ進む。 | cooldown 中は `target_status=skipped_cooldown`、終了コード `0`。 |
| GitHub tree resolve | `.github_token`, branch target | `.build_logs/{id}.json` | target file の blob SHA を取得する。 | 全 retry 失敗は `failure_api` を記録し、SHA を更新しない。 |
| SHA decision | `sha_file`, `.server_config` | none | 変更あり、force 条件成立、または webhook/force queue payload により build 対象を確定する。 | 変更なしは `skipped_no_change`。状態ファイルを更新しない。 |
| blob materializer | GitHub blob API, branch target | `src` | Base64 decode 後、対象 Markdown を atomic write する。 | decode/write 失敗は `failure_decode`。SHA を更新しない。 |
| precheck | branch target, output path, build binary | `.build_logs/{id}.json` | disk、binary、version がすべて合格する。 | `failure_precheck` を記録し、pipeline を起動しない。 |
| pipeline executor | `src`, `.pipeline_config` | `.build_logs/{id}.json`, `.build_history` | timeout 前に exit code `0`。 | 非 0 / timeout は `failure_build`。SHA、deploy、snapshot を更新しない。 |
| report importer | pipeline stdout | `.build_logs/{id}.json` | `[REPORT]` を 1 行だけ検出し schema へ変換する。 | `[REPORT]` なしは warning として `report:null` を保存する。build 成否は pipeline exit code に従う。 |
| SHA updater | `sha_file`, new SHA | `sha_file` | `{"sha":"<new_sha>"}` を atomic write する。 | `failure_state_write`。deploy、snapshot を実行しない。 |
| deploy executor | `deploy_targets`, output site | `.pending_transfers`, `.build_logs/{id}.json` | 全 deploy target の転送と検証が成功する。 | 失敗 target を `.pending_transfers` に保存し、`target_status=success_deploy_pending`。 |
| snapshot writer | output site, `.server_config` | `.snapshots/`, `.build_logs/{id}.json` | `snapshots_keep > 0` の場合に新 snapshot を保存し、超過世代を削除する。 | snapshot 失敗は build 成功を取り消さず WARN として記録する。 |
| status finalizer | target results, `.build_state`, `.build_circuit_state`, `.pending_transfers`, `.notify_pending` | `.build_status.json` | runner 起動 1 回の最終要約を保存する。 | 書き込み失敗は ERROR ログを出し、runner 終了コードを最低 `1` にする。 |
| finalizer | target results, `.build_state` | `.build_state`, `.build_lock` | `running=false`、`last_finished_at` を保存し、lock を削除する。 | lock 削除失敗は WARN。`.build_state.running=false` 保存失敗は終了コード `1`。 |

`failure_api`、`failure_decode`、`failure_precheck`、`failure_build`、`failure_state_write` は、同じ build id の `.build_logs/{id}.json.status` では `"failure"` として保存し、詳細理由は `error` に固定文言で保存する。`target_status` は runner 内部分類および `.build_logs/{id}.json.error` の詳細判定に使用し、API response の `status` 値としては返さない。

**ビルドトリガー種別契約：**

`trigger` は runner が処理を開始した原因を表す固定文字列であり、`.build_logs/{id}.json`、`.build_history`、`.build_status.json`、API response、SDK 型、UI 表示で同じ値を使用する。実装者判断で `auto`、`force`、`scheduled` など別名を追加してはならない。

| `trigger` | 発生条件 | build log | history | status |
|-----------|----------|-----------|---------|--------|
| `polling` | systemd timer 等の通常起動で SHA 差分により build する。 | 記録する | 記録する | 記録する |
| `force_interval` | SHA 差分なしだが `.server_config.force_build_interval_hours` 条件成立により build する。 | 記録する | 記録する | 記録する |
| `manual` | `POST /api/build` により `.build_state.queued` へ投入された entry を処理する。 | 記録する | 記録する | 記録する |
| `webhook` | GitHub Webhook 受信により `.build_state.queued` へ投入された entry を処理する。 | 記録する | 記録する | 記録する |
| `retry_pending_transfer` | `.pending_transfers` の再送のみを実行する。 | 再送専用 log を作成する場合に記録する | 再送履歴を追記する場合に記録する | 記録する |
| `startup_config_integrity` | 設定ファイル起動時整合性チェックで破損復旧または停止が発生する。 | 作成しない | 追記しない | 記録する |
| `rollback` | `POST /api/history/{id}/rollback` により snapshot を再転送する。 | 記録する | 記録する | 記録する |
| `local_watch` | `watch_mode="local"` の local SHA 差分により build する。 | 記録する | 記録する | 記録する |
| `approval` | `POST /api/approvals/{id}/approve` により承認済み queue entry を処理する。 | 記録する | 記録する | 記録する |

queue entry の `trigger` は `"manual"`、`"webhook"`、`"approval"` のみ許可する。`"force"` は使用せず、強制実行 API は queue 保存時に `"manual"` と `payload.force=true` を保存する。`force_interval` と `local_watch` は runner が設定値と差分検出結果から内部判定する場合のみ使用する。

**`.build_status.json` 更新契約：**

`.build_status.json` の schema、許容値、初期値、api / ui / mcp の読取 priority は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a、§22.0c を正とする。本ファイルでは runner が `.build_status.json` を更新するタイミングと、更新失敗時の runner 挙動だけを定義する。

`.build_status.json` の更新タイミングは次に固定する。

| タイミング | `status` | `trigger` | 備考 |
|------------|----------|-----------|------|
| lock 取得後、target 処理前 | `running` | 判定済み trigger | `current_build_id` を設定する。 |
| 実行中 lock 検出 | `lock_skipped` | `polling` | `.build_state` と build log は変更しない。 |
| cooldown skip | `skipped_cooldown` | `polling` | `last_build_id` は既存値を保持する。 |
| SHA 変更なし | `skipped_no_change` | `polling` | build log / history は作成しない。 |
| force interval build 成功 | `success` | `force_interval` | build log / history と同じ build id を保存する。 |
| build 成功・deploy 成功 | `success` | 実行 trigger | `last_deploy_status="success"`。 |
| build 成功・deploy pending | `success_deploy_pending` | 実行 trigger | `pending_transfers_count` を更新する。 |
| build 失敗 | `failure` | 実行 trigger | `last_error` に固定エラー文言を保存する。 |
| circuit open skip | `circuit_open` | `polling` | GitHub API 呼び出し前に保存する。 |
| startup config 復旧あり | `config_recovered` | `startup_config_integrity` | build log / history は作成しない。復旧後に通常 build へ進む場合、次の `running` 更新で上書きする。 |
| startup config 停止 | `config_error` | `startup_config_integrity` | `.build_state.running=true` へ進まない。 |
| pending transfer retry のみ | `success` または `failure` | `retry_pending_transfer` | build が発生しない場合でも pending 件数を更新する。 |

`.build_status.json` は atomic write 対象であり、書き込み失敗時に部分ファイルを残してはならない。未知 key を追加してはならない。API は GET request で `.build_status.json` を自動修復してはならない。

**queue entry 実行契約：**

`.build_state.queued` に entry がある場合、runner は通常ポーリング対象の前に queue を FIFO で 1 件だけ取り出して処理する。queue entry 処理が成功または失敗として `.build_history` に記録された場合、その entry を queue から削除する。runner 起動 1 回で複数 queue entry を連続処理してはならない。queue entry の `trigger` が `"manual"` かつ `payload.force=true` の場合は SHA 比較を行わず build を実行する。`trigger` が `"webhook"` の場合は payload の `ref` と `sha` を優先し、branch target に一致しない entry は `failure_api` として記録した後に queue から削除する。

queue entry は JSON object とし、最低限 `id`、`trigger`、`created_at`、`payload` を持つ。`id` は queue 内で一意、`created_at` は UTC ISO 8601、`payload` は JSON object とする。不正 entry が先頭にある場合、runner はその entry を `failure_decode` として build log / history に記録して queue から削除し、次回起動まで次 entry は処理しない。queue 全体が JSON として破損している場合は §12 の `.build_state` 破損処理に従う。

**cooldown / force build 判定契約：**

runner は起動ごとに `.server_config.build_cooldown_seconds` と `.server_config.force_build_interval_hours` を 1 回読み、当該起動中の cooldown / force build 判定に使用する。同一 runner 起動中に `.server_config` を再読込して判定値を変更してはならない。

判定順は queue、pending retry、circuit breaker、cooldown、SHA decision の順とする。manual queue entry の `payload.force=true` は cooldown を無視する。webhook queue entry は cooldown を適用する。`force_build_interval_hours` は SHA 一致時だけ評価し、SHA 不一致時は常に通常 build とする。

| 条件 | 結果 |
|------|------|
| `.build_state.last_finished_at=null` | cooldown は適用しない。 |
| `.server_config.build_cooldown_seconds=0` | cooldown は無効。 |
| `now - last_finished_at < build_cooldown_seconds` | `skipped_cooldown`。GitHub API、pipeline、deploy、snapshot は実行しない。 |
| SHA 一致かつ `.server_config.force_build_interval_hours=0` | `skipped_no_change`。 |
| SHA 一致かつ `force_build_interval_hours>0` かつ直近成功 build から指定時間未満 | `skipped_no_change`。 |
| SHA 一致かつ `force_build_interval_hours>0` かつ直近成功 build から指定時間以上 | `force_interval` として build を実行する。 |

force interval の直近成功 build は `.build_history` のうち同じ `branch` と `target_file` で status が `success` または `success_deploy_pending` の最新行とする。`.build_history` が存在しない、または該当行がない場合は force interval 条件成立として build する。

**runner finalizer 固定契約：**

runner は lock 取得後、正常終了、失敗終了、panic 相当の recover、context timeout のいずれでも finalizer を実行する。finalizer は次の順序に固定する。

1. 未保存の `.build_logs/{id}.json` がある場合は、取得済みで schema を満たす値だけを使って最終形を保存する。未取得値を推測して補完してはならない。
2. `.build_history` へ追記対象の build がある場合は 1 行だけ追記する。同じ `id` が既に存在する場合は追記せず、ERROR ログ `BUILD_HISTORY_DUPLICATE: id={id}` を出す。
3. `.build_status.json` を最終状態へ更新する。
4. `.build_state.running=false`、`current_build_id=null`、`last_finished_at={now}` を保存する。
5. `.build_lock` を削除する。
6. 通知対象 event がある場合は通知または `.notify_pending` 追記を行う。

finalizer 中に複数失敗が発生した場合、終了コードは最も重い値を採用する。`.build_state.running=false` の保存失敗は終了コード `1` 固定とし、lock 削除だけ成功しても正常終了扱いにしない。lock 削除失敗は WARN とし、他失敗がなければ終了コードを変更しない。

**build log / history 書き込み契約：**

`.build_logs/{id}.json` は atomic write で 1 build id につき 1 file だけ作成する。既に同名 file が存在する場合は上書きせず、次の suffix 付き build id を採番し直す。`.build_history` は JSON Lines とし、追記前に既存 file の末尾が LF で終わることを確認する。LF がない場合は 1 個だけ LF を追加してから新規行を追記する。

`.build_history` の 1 行は `.build_logs/{id}.json` の要約であり、少なくとも `id`、`status`、`trigger`、`branch`、`target_file`、`started_at`、`finished_at`、`duration_seconds`、`commit_sha`、`blob_sha`、`warnings`、`error` を含む。history へ保存する `status` は `target_status` と同じ値を使用する。JSON Lines の壊れた既存行は読み取り対象から除外するが、追記時に既存 file 全体を書き換えてはならない。

**サーキットブレーカー更新契約：**

`.build_circuit_state.open=true` の場合、runner は GitHub API、pipeline、deploy、snapshot を実行せず、`.build_status.json` に `status="circuit_open"` を保存して終了コード `0` で終了する。pending transfer retry と notify pending retry は circuit open 判定前の再試行処理として実行する。

連続失敗数を増やす対象は `failure_build`、`failure_precheck`、`failure_decode`、`failure_state_write`、`success_deploy_pending` とする。`failure_api`、`skipped_no_change`、`skipped_cooldown`、`lock_skipped`、通知失敗だけの成功 build は連続失敗数を増やさない。いずれかの target が `success` になった場合だけ、`consecutive_failures` は 0 に戻す。

**SHA 更新禁止条件：**

runner は以下のいずれかに該当する場合、`sha_file` を更新してはならない。

| 条件 | 理由 |
|------|------|
| GitHub Trees API 失敗 | 対象 SHA が確定していない。 |
| GitHub Blob API 失敗 | 入力 Markdown が取得できていない。 |
| blob decode / src 書込失敗 | builder へ渡す入力が確定していない。 |
| precheck 失敗 | pipeline を実行していない。 |
| pipeline timeout / 非 0 | 出力成果物が成功状態ではない。 |
| `[REPORT]` 不在かつ pipeline 非 0 | 成功確認できない。 |
| status / log / history の必須保存失敗 | 実行結果を追跡できない。 |

pipeline が終了コード `0` で `[REPORT]` が不在の場合、SHA を更新する。ただし `.build_logs/{id}.json.report=null`、`warnings` に `REPORT_MISSING` を追加し、`.build_history.warnings` に 1 を加算する。

**状態ファイル破損時の処理：**

runner が読み込む JSON object / JSON array の状態ファイルが破損している場合は、§22.0a の破損時の扱いに従う。JSON Lines は壊れた行だけを無視し、ファイル全体を破棄してはならない。破損退避ファイル名は `{original}.corrupt.{YYYYMMDDHHMMSS}.bak` とする。

```
runner 起動（systemd タイマーから呼び出し）
    │
    ├─ .github_token 読み込み（不在、空、改行除去後 1 文字未満の場合は ERROR ログ、終了コード 2）
    │
    ├─ [重複チェック] .build_lock が存在する場合
    │   ├─ ファイル内 PID が実行中 → INFO ログ（`BUILD_SKIP: already running (PID N)`）、正常終了
    │   └─ PID が存在しない（前回の異常終了） → .build_lock を削除して続行
    │
    ├─ .build_lock に自プロセスの PID と started_at を書き込み
    │   （以降、正常終了・例外終了いずれの場合も defer で .build_lock を削除）
    │
    ├─ [設定ファイル起動時整合性チェック]
    │   ├─ .branch_config / .notify_config / .build_state / .build_circuit_state / .pending_transfers / .notify_pending を固定順序で検証
    │   ├─ 復旧可能な破損 → corrupt backup、初期化または default fallback、§27.10 の通知条件を満たす場合は config_corrupt 通知
    │   ├─ unknown key のみ → backup せず正規化書き戻し
    │   └─ permission error / IO error → .build_state.running=true にせず終了
    │
    ├─ [ペンディングキュー再試行] PENDING_FILE が存在する場合（→ §14a）
    │   └─ ペンディングエントリごとに SSH 転送を再試行
    │       ├─ 成功 → エントリを PENDING_FILE から削除
    │       └─ 失敗 → ERROR ログ、エントリを保持（次回起動時に再試行）
    │
    ├─ [Webhook 通知ペンディング再試行] .notify_pending が存在する場合
    │   └─ ペンディングエントリごとに Webhook 送信を再試行
    │       ├─ 成功（HTTP 2xx）→ INFO ログ、エントリを .notify_pending から削除
    │       └─ 失敗（HTTP エラー・接続失敗）→ ERROR ログ、エントリを保持（次回起動時に再試行）
    │
    ├─ [クールダウンチェック] build_cooldown_seconds > 0 の場合
    │   └─ .build_state.last_finished_at から build_cooldown_seconds 秒未満
    │       → INFO ログ（`COOLDOWN: skip, last_build N秒前`）、正常終了
    │
    ├─ BRANCH_TARGETS の各エントリを順次処理：
    │   │
    │   ├─ Step 1: Git Trees API
    │   │   GET /repos/{OWNER}/{REPO}/git/trees/{branch}?recursive=1
    │   │   → target_file の blob SHA を取得
    │   │   └─ API 失敗時：API_RETRY_MAX 回まで指数バックオフ（API_RETRY_BASE_SECONDS × 2^n 秒）で再試行
    │   │        ├─ 全試行失敗時：ERROR ログ、このエントリを failure(api_error) として記録し、SHA を更新せず次エントリへ進む
    │   │        └─ レスポンスヘッダー X-RateLimit-Remaining = 0 の場合
    │   │             → X-RateLimit-Reset（Unix 時刻）まで待機してから再試行
    │   │               INFO ログ（`RATE_LIMIT: waiting until {reset_time}`）
    │   │   [PAT 有効期限チェック] レスポンスヘッダー GitHub-Authentication-Token-Expiration が存在する場合
    │   │        ├─ 残日数 ≤ 7 日 → WARN ログ（`PAT_EXPIRY_WARN: expires_at={date} remaining_days={N}`）
    │   │        └─ 残日数 > 7 日 → 処理継続（チェックのみ）
    │   │
    │   ├─ SHA 比較（sha_file の前回 SHA と比較）
    │   │   ├─ 一致（変更なし）かつ force_build_interval_hours = 0 → INFO ログ、このエントリをスキップ
    │   │   ├─ 一致（変更なし）かつ force_build_interval_hours > 0 → 前回ビルドから指定時間以上経過していれば強制ビルド続行
    │   │   └─ 不一致（変更あり）→ 続行
    │   │
    │   ├─ Step 2: Git Blobs API
    │   │   GET /repos/{OWNER}/{REPO}/git/blobs/{sha}
    │   │   → Base64 デコード → src パスへ書き出し
    │   │   └─ API 失敗時：API_RETRY_MAX 回まで指数バックオフで再試行
    │   │        ├─ 全試行失敗時：ERROR ログ、このエントリを failure(api_error) として記録し、SHA を更新せず次エントリへ進む
    │   │        └─ レスポンスヘッダー GitHub-Authentication-Token-Expiration が存在する場合は Step 1 と同様に PAT 有効期限チェックを行う
    │   │
    │   ├─ [コミット情報取得] ビルドトリガーとなったコミット情報を取得し、ビルドログへ記録する
    │   │   GET /repos/{OWNER}/{REPO}/commits?path={BRANCH_TARGET.src}&sha={branch}&per_page=1
    │   │   → 先頭エントリから取得：
    │   │       commit_sha    = commit.sha
    │   │       commit_message = commit.commit.message（先頭1行のみ）
    │   │       commit_author  = commit.commit.author.name
    │   │       commit_at      = commit.commit.author.date
    │   │   └─ API 失敗時：各フィールドを null として記録し、処理続行（ビルドは妨げない）
    │   │
    │   ├─ [事前チェック] pipeline.sh 実行前に以下を確認し、不足時は ERROR ログ＋deploy_failure Webhook 通知、このエントリをスキップ
    │   │   ├─ ディスク空き容量 ≥ max(出力サイト推定サイズ × 3, 64MiB)。取得は `syscall.Statfs(outDir)` を使用する
    │   │   ├─ `adlaire-ci-build` が存在し実行可能であること（`os.Stat` と mode bit）
    │   │   └─ `/usr/local/bin/adlaire-ci-build --version` が終了コード 0 で、stdout に `adlaire-ci-build` と `ADLAIRE_CI_SPEC` を含むこと
    │   │
    │   ├─ pipeline.sh 実行（bash {src の親ディレクトリ}/.ci/pipeline.sh）
    │   │   ├─ 成功（exit 0）：INFO ログ
    │   │   │   └─ [Webhook 通知送信] on: ["success"] 設定時
    │   │   │       → .notify_config の Webhook 宛先へ POST（ペイロード: event="success", branch, build_id 等）
    │   │   │       → 送信失敗（HTTP エラー・接続失敗・タイムアウト）の場合：
    │   │   │           ERROR ログ（`NOTIFY_FAIL: url={url} status={code}`）
    │   │   │           .notify_pending へ `{"event":"success","url":"...","payload":{...},"queued_at":"<ISO8601>"}` を追記
    │   │   └─ 失敗（exit ≠ 0）：ERROR ログ、sha_file 更新せず、このエントリを failure(build_failed) として記録し、次エントリへ進む
    │   │       └─ [Webhook 通知送信] on: ["failure"] 設定時
    │   │           → .notify_config の Webhook 宛先へ POST（ペイロード: event="failure", branch, build_id 等）
    │   │           → 送信失敗の場合：ERROR ログ、.notify_pending へキューイング（success と同一形式）
    │   │
    │   └─ sha_file を新 SHA で更新（`{"sha": "<new_sha>"}` を JSON 書き込み）
    │        └─ SSH サイト転送（deploy_targets リストの各エントリへ転送 → §14a）
    │             ├─ [転送後整合性検証] ssh user@host "sha256sum /dest/<relative-path>" でリモート SHA を取得
    │             │   ├─ 全ファイルのローカル sha256 と一致 → 転送成功
    │             │   └─ 不一致またはコマンド失敗 → ERROR ログ、ペンディングキューへ再投入（§14a）
    │             └─ 整合性検証成功後 → スナップショット保存（→ §14b）
    │
    │        [出力サイズチェック] OUTPUT_SIZE_WARN_MB > 0 の場合
    │        出力サイト配下の通常ファイル合計サイズを取得し、閾値と比較：
    │            size_mb = total_site_bytes / (1024 * 1024)
    │            size_mb > OUTPUT_SIZE_WARN_MB の場合：
    │            → WARN ログ（`OUTPUT_SIZE_WARN: size={size_mb:.1f}MB threshold={OUTPUT_SIZE_WARN_MB}MB`）
    │            → ビルドログの size_warn フィールドを true に設定（§8）
    │
    └─ [サーキットブレーカー判定] API_CIRCUIT_BREAKER_THRESHOLD > 0 の場合
        全ブランチの今周回結果を集計し、全ブランチが失敗（API エラー・スキップを除くビルド失敗）の場合：
        → .build_circuit_state（JSON）の consecutive_failures を +1
        いずれか成功した場合：
        → consecutive_failures を 0 にリセット
        consecutive_failures ≥ API_CIRCUIT_BREAKER_THRESHOLD の場合：
        → ERROR ログ（`CIRCUIT_OPEN: consecutive_failures={N} threshold={API_CIRCUIT_BREAKER_THRESHOLD}`）
        → deploy_failure Webhook 通知（reason="circuit_open"）
        → .build_circuit_state に `{"open": true, "consecutive_failures": N, "opened_at": "<ISO8601>"}` を書き込み
        → 以降のポーリング周回はビルドをスキップ（SHA チェックも行わない）
        → POST /api/circuit-breaker/reset でリセット可能（open: false、consecutive_failures: 0 に戻す）

    ├─ [ログクリーンアップ] LOG_KEEP_N > 0 の場合
    │   .build_logs/ 内の {id}.json を mtime 昇順でソートし、
    │   件数が LOG_KEEP_N を超えた分を古いものから削除
    │
    └─ [週次サマリー判定] WEEKLY_SUMMARY_ENABLED = true かつ on: ["weekly_summary"] 設定の Webhook 宛先が存在する場合
        現在の曜日が WEEKLY_SUMMARY_DAY かつ現在時刻が WEEKLY_SUMMARY_HOUR:00±30分以内の場合：
        ├─ 二重送信防止チェック：.build_state の weekly_summary_sent_date と当日の日付（YYYY-MM-DD）を比較
        │   → 一致（本日送信済み）→ スキップ（INFO ログ: `WEEKLY_SUMMARY_SKIP: already sent today`）
        │   → 不一致（未送信）→ 以下を実行
        ├─ 過去 7 日間の .build_logs/{id}.json を集計：
        │       成功件数・失敗件数・成功率（%）・平均ビルド時間（秒）・最長ビルド（秒・ID）
        ├─ on: ["weekly_summary"] 設定の Webhook 宛先へ POST
        │       ペイロード: event="weekly_summary", period="7d", success_count, failure_count, success_rate, avg_duration_seconds, max_duration_seconds, max_duration_build_id
        ├─ 送信成功 → .build_state に `weekly_summary_sent_date: "YYYY-MM-DD"` を記録（INFO ログ）
        └─ 送信失敗（HTTP エラー・接続失敗）→ ERROR ログ、.notify_pending へキューイング（他の Webhook 失敗と同一形式）
```

---

## 14. `pipeline.sh`

リポジトリの `.ci/pipeline.sh` にビルド手順を記述する。

- 実行権限（`chmod +x`）が必要
- 先頭の shebang は `#!/usr/bin/env bash` または `#!/bin/bash` とする
- `set -euo pipefail` を shebang 直後に記述する
- 終了コード `0` で成功、`0` 以外で失敗とみなす
- runner は `bash {src の親ディレクトリ}/.ci/pipeline.sh` として実行し、作業ディレクトリは `src` の親ディレクトリに設定する
- 環境変数として `ADLAIRE_CI_SRC`、`ADLAIRE_CI_OUT`、`ADLAIRE_CI_BRANCH`、`ADLAIRE_CI_BUILD_ID` を渡す
- `pipeline.sh` は最終的に `ADLAIRE_CI_OUT` のパスへ HTML を生成しなければならない

**例：**
```bash
#!/usr/bin/env bash
set -euo pipefail
/usr/local/bin/adlaire-ci-build --src "$ADLAIRE_CI_SRC" --out "$ADLAIRE_CI_OUT"
```

ビルド実行コマンドは `pipeline.sh` 内に直接記述する（`runner` は参照しない）。`adlaire-ci-build` は `components/builder.go` から生成した Go 版バイナリである。標準配置への移行完了前は、現行実装実体 `build_spec.go` から生成した同等バイナリとして扱う。

**runner からの実行契約：**

| 項目 | 仕様 |
|------|------|
| command | `bash {srcParent}/.ci/pipeline.sh` |
| working directory | `{srcParent}` |
| timeout | `.server_config.build_timeout_seconds` が存在すればその値、なければ `300` 秒。 |
| stdout/stderr | それぞれ最大 1 MiB までメモリに保持し、超過分は末尾 1 MiB を保存する。超過時は `stdout_truncated` / `stderr_truncated` を `true` にする。 |
| 環境変数 | 親 process の環境を引き継ぎ、下表の値で上書きする。 |

| 環境変数 | 値 |
|----------|----|
| `ADLAIRE_CI_SRC` | `BranchTarget.Src` |
| `ADLAIRE_CI_OUT` | `BranchTarget.Out` |
| `ADLAIRE_CI_BRANCH` | `BranchTarget.Branch` |
| `ADLAIRE_CI_BUILD_ID` | build id |
| `ADLAIRE_CI_TARGET_FILE` | `BranchTarget.TargetFile` |
| `ADLAIRE_CI_STATE_DIR` | `RunnerConfig.StateDir` |

`pipeline.sh` が timeout した場合、runner は process group を終了し、`target_status="failure_build"`、`error="pipeline timeout"` として記録する。timeout 時も stdout/stderr の取得済み内容は `.build_logs/{id}.json` に保存する。

**GitHub API 固定契約：**

| 項目 | 仕様 |
|------|------|
| Base URL | `https://api.github.com` |
| User-Agent | `adlaire-ci-runner` |
| Authorization | `Bearer {token}` |
| Accept | `application/vnd.github+json` |
| API version header | `X-GitHub-Api-Version: 2022-11-28` |
| timeout | 30 秒 |
| retry 対象 | network error、timeout、HTTP `429`、`500`、`502`、`503`、`504` |
| retry 対象外 | HTTP `400`、`401`、`403`（rate limit を除く）、`404`、`422` |
| backoff | `API_RETRY_BASE_SECONDS * 2^attempt` 秒。attempt は 0 始まり。 |

HTTP `401` は `failure_api` とし、ERROR ログ `GITHUB_AUTH_FAILED` を出す。HTTP `404` は `target_file` または branch 設定不正として `failure_api` とし、ERROR ログ `GITHUB_NOT_FOUND: branch={branch} target={target_file}` を出す。HTTP `403` で `X-RateLimit-Remaining: 0` の場合のみ rate limit として reset まで待機する。

**GitHub API response 処理契約：**

| 対象 | 条件 | 処理 |
|------|------|------|
| Trees API | `truncated=true` | `failure_api`。ERROR `GITHUB_TREE_TRUNCATED: branch={branch}` を出し、Blob API へ進まない。 |
| Trees API | `target_file` が file として一致 | 対象 blob SHA を使用する。 |
| Trees API | `target_file` が directory として一致 | 配下 `.md` file の blob SHA を path 昇順に連結し、SHA-256 hex を target digest とする。Blob API は各 `.md` file に対して実行する。 |
| Trees API | `target_file` が見つからない | `failure_api`。ERROR `GITHUB_TARGET_MISSING: branch={branch} target={target_file}`。 |
| Blob API | `encoding!="base64"` | `failure_decode`。 |
| Blob API | Base64 decode 失敗 | `failure_decode`。 |
| Blob API | decode 後 UTF-8 不正 | `failure_decode`。 |
| Commits API | 取得失敗 | build は継続し、commit fields を `null` にする。 |

directory target の materialize では、GitHub path から `target_file` prefix を取り除いた相対 path を `BranchTarget.Src` 配下に再現する。対象外 file、hidden directory、`.git`、`.ci` は書き出さない。書き出し前に `BranchTarget.Src` 配下の前回 materialized Markdown を一時 directory へ置換し、途中失敗時は既存 `src` を保持する。

rate limit 待機は `X-RateLimit-Reset` が現在時刻より未来かつ 3600 秒以内の場合だけ実行する。3600 秒を超える場合、または header が不正な場合は待機せず `failure_api` とする。待機中に context timeout または SIGTERM を受けた場合は `failure_api` として finalizer へ進む。

**pipeline 実行結果分類：**

| 条件 | `pipeline.exit_code` | `target_status` | `error` | retry |
|------|----------------------|-----------------|---------|-------|
| exit `0` | `0` | `success` 候補 | `null` | なし |
| exit `1`〜`125` | 実際の終了コード | `failure_build` | `pipeline failed` | §27.3 の retry 対象。 |
| exit `126` | `126` | `failure_precheck` | `precheck failed` | retry しない。 |
| exit `127` | `127` | `failure_precheck` | `precheck failed` | retry しない。 |
| signal 終了 | `128 + signal` | `failure_build` | `pipeline failed` | retry 対象。 |
| timeout | `null` | `failure_build` | `pipeline timeout` | retry 対象。 |
| stdout 上限超過 | 実際の終了コード | exit code に従う | exit code に従う | exit code に従う。 |
| stderr 上限超過 | 実際の終了コード | exit code に従う | exit code に従う | exit code に従う。 |

runner は stdout / stderr の CRLF を LF に正規化して保存する。NUL byte は `\u0000` 文字列へ置換する。保存する stdout / stderr は UTF-8 不正 byte を `�` に置換する。secret mask は保存前に適用し、PAT、Webhook Secret、SMTP password、API token、session token、TOTP secret に一致する値を `***` に置換する。

---

## 14a. SSH サイト転送

本節は、Go 版 CI ランナーの SSH 転送標準仕様である。

`runner` は、`pipeline.sh` 成功後に、出力サイトディレクトリを SSH 経由で静的コンテンツ配信サーバーへ転送する。本節を SSH 転送の正本仕様とする。

`runner` は `pipeline.sh` 成功後に、出力サイトディレクトリ配下の全ファイルを SSH 経由で静的コンテンツ配信サーバーへ転送する。scp・rsync は使用しない。SSH コマンドは `ssh` バイナリを `exec.CommandContext` で直接起動し、`/bin/sh -c` を使わない。

### 設定値

`BRANCH_TARGETS` 各エントリの `deploy_targets` リスト内で管理する（→ §12）。

| フィールド | 説明 | 例 |
|-----------|------|-----|
| `host` | 配信サーバーのホスト名 / IP | `"192.0.2.1"` |
| `user` | SSH 接続ユーザー | `"deploy"` |
| `dest_dir` | 配信サーバー上の転送先ディレクトリ | `"/var/www/html"` |

転送対象は `BRANCH_TARGETS[n]["out"]` のディレクトリ配下にある通常ファイルすべてとする。`deploy_targets` に複数エントリを定義した場合は全ての転送先へ順次転送する。

### 差分検出

転送前にリモートサーバーで対象ファイルごとの SHA256 ハッシュを取得し、ローカルファイルのハッシュと比較する。

```bash
# runner が os/exec 経由で実行
ssh <user>@<host> sha256sum <dest_dir>/<relative-path>
```

- ハッシュが一致 → 当該ファイルをスキップ（`SKIP` ログを記録）
- ハッシュが不一致、またはリモートにファイルが存在しない → 当該ファイルを転送する

relative path は `out` からの相対 path とし、`filepath.Rel` 後に `/` 区切りへ変換して保存する。空文字、`.`、`..` を含む path、先頭 `/`、NUL byte、制御文字を含む path は転送対象から除外し、`failure_precheck` とする。symbolic link、directory、device file、FIFO は転送しない。symbolic link を検出した場合は ERROR `DEPLOY_UNSUPPORTED_FILE: path={path}` を出し、当該 target を `failure_precheck` とする。

### 転送

stdin パイプ経由で SSH 転送する。

```bash
# runner が os/exec（StdinPipe）経由で実行
ssh <user>@<host> 'mkdir -p <dest_dir>/<relative-dir> && tee <dest_dir>/<relative-path>'
```

runner は local file を開き、SSH process の stdin へ `io.Copy` で送る。リモート側 stdout は使用せず破棄し、stderr は失敗理由として `.build_logs/{id}.json.error` と ERROR ログへ記録する。

SSH command は local shell 文字列を組み立てず、`exec.CommandContext` の argv として分離して起動する。directory 作成は `exec.CommandContext(ctx, "ssh", user+"@"+host, "mkdir", "-p", "--", remoteDir)`、file 転送は `exec.CommandContext(ctx, "ssh", user+"@"+host, "tee", "--", remotePath)` の 2 段階に分ける。`dest_dir`、relative path、host、user を `/bin/sh -c` 用の 1 文字列へ連結して渡してはならない。host と user は `^[A-Za-z0-9._-]+$` に一致する値だけ許可する。

転送は file 単位で行い、1 file の転送 timeout は 60 秒とする。timeout 時は SSH process group を終了し、当該 deploy target を pending とする。1 deploy target 内で 1 file でも転送または検証に失敗した場合、その deploy target 全体を pending とし、snapshot は作成しない。

### ペンディングキュー

転送失敗時（接続エラー・認証失敗等）は `PENDING_FILE`（JSON）へエントリを追記する。

```json
[
  {
    "branch_idx": 0,
    "deploy_idx": 0,
    "out": "/opt/adlaire-builder/dist/site",
    "host": "192.0.2.1",
    "user": "deploy",
    "dest_dir": "/var/www/html",
    "failed_at": "2026-09-15T10:00:00Z",
    "retry_count": 1
  }
]
```

- `runner` 起動時（`BRANCH_TARGETS` 処理前）に `PENDING_FILE` を読み込み、エントリごとに再試行する（→ §13 処理フロー）
- 再試行成功時にエントリを削除する。失敗時は `retry_count` をインクリメントして保持する
- SSH 転送失敗 Webhook 通知（`deploy_failure` イベント）を送信する（on: `["deploy_failure"]` 設定時）
- 同一 `out`、`host`、`user`、`dest_dir` の pending エントリが既に存在する場合は新規追記せず、既存エントリの `retry_count` を +1 し、`failed_at` を最新時刻へ更新する

pending entry は次の key だけを保存する。未知 key を保存してはならない。

| key | 型 | 内容 |
|-----|----|------|
| `branch_idx` | integer | `RunnerConfig.BranchTargets` の 0 始まり index。 |
| `deploy_idx` | integer | `DeployTargets` の 0 始まり index。 |
| `out` | string | local output directory の絶対パス。 |
| `host` | string | deploy target host。 |
| `user` | string | deploy target user。 |
| `dest_dir` | string | remote destination directory。 |
| `failed_at` | string | UTC ISO 8601。 |
| `retry_count` | integer | 1 以上。 |
| `last_error` | string | secret mask 済みの短い失敗理由。最大 500 文字。 |

pending retry 時に元の `branch_idx` または `deploy_idx` が現在設定範囲外の場合は、その entry を削除せず `retry_count` を +1 し、ERROR `PENDING_TARGET_MISSING: branch_idx={branch_idx} deploy_idx={deploy_idx}` を出す。現在設定の `out`、`host`、`user`、`dest_dir` が pending entry と異なる場合は、entry の保存値を優先して再送する。

### 転送後整合性検証

SSH 転送完了後に、リモートファイルの SHA-256 チェックサムをローカルのものと照合する。

**検証コマンド：**
```
ssh {user}@{host} sha256sum {dest_dir}/{filename}
```
出力形式 `{hash}  {filename}` の最初のフィールドをローカル `crypto/sha256` の hex digest と比較する。

| 項目 | 仕様 |
|---|---|
| タイムアウト | 30 秒（SSH 転送タイムアウトとは独立） |
| 検証失敗時 | ERROR ログ＋ペンディングキューへ再投入。スナップショット保存はしない |
| ログフィールド | `transfer_verified: false`（`.build_logs/{id}.json` に記録） |
| 正常時 | `transfer_verified: true`（`.build_logs/{id}.json` に記録） |

remote `sha256sum` 出力は 1 行目の先頭 field だけを採用し、hex 64 文字以外は検証失敗とする。複数行出力、空出力、stderr 出力のみ、終了コード非 0 は検証失敗とする。local checksum は転送直前に読んだ file 内容ではなく、転送後に local file を再読込して計算する。

### ログ

| 状態 | ログレベル | メッセージ例 |
|------|-----------|------------|
| スキップ（差分なし） | `INFO` | `SKIP site: no change` |
| 転送成功 | `INFO` | `DEPLOY site → 192.0.2.1` |
| 転送失敗→キューイング | `ERROR` | `DEPLOY FAILED site: <reason> (queued)` |
| ペンディング再試行成功 | `INFO` | `PENDING RETRY OK site → 192.0.2.1` |
| ペンディング再試行失敗 | `ERROR` | `PENDING RETRY FAILED site: <reason>` |
| 整合性検証失敗→再投入 | `ERROR` | `VERIFY FAILED site → 192.0.2.1: checksum mismatch (queued)` |

---

## 14b. スナップショット管理

本節は、Go 版 CI ランナーのスナップショット標準仕様である。

`runner` は、SSH 転送成功後に `.snapshots/` ディレクトリへ成果物を保存する。本節をスナップショット保存、世代管理、ロールバック連携の正本仕様とする。

`runner` は SSH 転送成功後に、ビルド成果物を `.snapshots/` ディレクトリへアーカイブする。`HISTORY_KEEP_N = 0` の場合はスナップショット世代削除を行わず、無制限保持とする。

### ディレクトリ構造

```
/opt/adlaire-builder/
└── .snapshots/
    ├── b20260915100000/           # build_id = b{YYYYMMDDHHmmss}
    │   └── site  # ビルド成果物のコピー
    ├── b20260914180000/
    │   └── site
    └── ...
```

- ディレクトリ名は `b{YYYYMMDDHHmmss}` 形式のビルド ID（`.build_history` の `id` と一致する）
- `BRANCH_TARGETS` に複数エントリがある場合は、同一ビルド ID ディレクトリ内に各エントリの成果物をまとめて保存する

### 世代管理

- スナップショット保存後、`.snapshots/` 内のディレクトリ数が `HISTORY_KEEP_N` を超えた場合、最古のディレクトリから順に削除する
- 削除対象ディレクトリの特定は作成日時降順ソートで行う（ディレクトリ名の辞書順 = 時系列順）

snapshot 保存は `{StateDir}/.snapshots/{build_id}.tmp.{pid}` へ copy した後、`{StateDir}/.snapshots/{build_id}` へ rename する。同じ snapshot id が既に存在する場合は上書きせず、WARN `SNAPSHOT_EXISTS: id={id}` を出して snapshot 保存を skip する。snapshot 内には output site 配下の通常ファイルだけを含め、`.github_token`、runner 状態ファイル、`.git`、lock、pending queue を含めてはならない。

**snapshot 保存対象固定契約：**

| 対象 | 扱い |
|------|------|
| 通常ファイル | output site 配下の相対 path を保持して copy する。file mode は実行 bit だけを保持し、setuid / setgid bit は落とす。 |
| directory | 必要な directory だけ作成し、mode は最大 `0755` とする。 |
| symlink | file / directory を問わず保存しない。WARN `SNAPSHOT_SKIP_SYMLINK: path={path}` を出す。 |
| hidden file | output site 配下の通常ファイルで、かつ `.git` directory 配下ではない場合だけ保存対象に含める。 |
| path traversal | snapshot 内相対 path に `..`、絶対 path、空 segment、NUL を含む場合は snapshot 保存失敗とする。 |
| runner 状態ファイル | `.github_token`、`.build_lock`、`.pending_transfers`、`.notify_pending`、`.build_state`、`.build_status.json`、`.admin_credentials`、`.api_tokens` は保存禁止。 |
| tmp directory | `{build_id}.tmp.{pid}` は成功時に残してはならない。失敗時も削除を 1 回試行し、失敗時は WARN `SNAPSHOT_TMP_CLEANUP_FAILED`。 |
| prune | 新 snapshot rename 成功後に実行する。prune 失敗は WARN とし、snapshot 成功を取り消さない。 |

snapshot copy 中に読み取り失敗、書き込み失敗、path 検証失敗、tmp rename 失敗が発生した場合、当該 snapshot は作成失敗とし、build 成功を取り消さない。`.build_logs/{id}.json.warnings` に `SNAPSHOT_SAVE_FAILED` を追加し、`snapshot_id=null` とする。snapshot 失敗を理由に `.pending_transfers` を作成してはならない。

### ロールバック

`POST /api/history/{id}/rollback`（→ §22）で指定ビルド ID のスナップショットから SSH 転送を再実行する。

- ロールバック API は `api` の実装を前提とする。`api` が実装されるまでは、API 経由のロールバックはとして扱う
- `.snapshots/{id}/` が存在しない場合は `404` を返す
- 転送成功時は `.build_history` に rollback エントリを追記する

### ログ

| 状態 | ログレベル | メッセージ例 |
|------|-----------|------------|
| スナップショット保存 | `INFO` | `SNAPSHOT b20260915100000 saved` |
| 古世代削除 | `INFO` | `SNAPSHOT b20260910000000 pruned (keep_n=10)` |
| ロールバック成功 | `INFO` | `ROLLBACK b20260914180000 → 192.0.2.1 OK` |
| ロールバック失敗 | `ERROR` | `ROLLBACK b20260914180000: <reason>` |

---

## 15. ログ

本節は、`runner` の stdout ログと構造化ビルドログを定義する。

### stdout ログ

| レベル | 出力条件 |
|--------|---------|
| `INFO` | 起動、変更なしスキップ、ビルド開始・完了、SHA 更新 |
| `WARNING` | — |
| `ERROR` | トークン読み込み失敗、API 失敗、ビルド失敗 |
| `DEBUG` | API レスポンス詳細等（`LOG_LEVEL = "DEBUG"` 時のみ） |

stdout は Go 標準ライブラリ `log/slog` で出力し、systemd が journald に転送する。独自 logger 実装を使用してはならない。

### 構造化ビルドログ

`runner` は、ビルドごとに `.build_logs/{id}.json` を作成する。

| 項目 | 内容 |
|------|------|
| ビルドログファイル | ビルドごとに `.build_logs/{id}.json` を作成する。 |
| stdout / stderr 保存 | `pipeline.sh` の標準出力・標準エラーをビルドログへ保存する。 |
| 変換レポート取り込み | `builder` が出力する `[REPORT]` 行をパースし、`tables_count`、`code_blocks_count` 等へ変換して保存する。 |
| 警告取り込み | `[WARN]` 行を配列として保存し、`warnings` 件数と整合させる。 |
| ビルド所要時間 | `started_at`、`finished_at`、`duration_seconds` を保存する。 |
| コミット情報 | ビルド対象 commit の SHA、message、author、date を保存する。 |
| 転送検証結果 | SSH 転送後整合性検証の結果として `transfer_verified` を保存する。 |
| 出力サイズ警告 | `OUTPUT_SIZE_WARN_MB` 超過時に `size_warn: true` を保存する。 |
| ログ世代管理 | `LOG_KEEP_N` を超過した `.build_logs/{id}.json` を古いものから削除する。 |

これらのログ項目を実装対象に含める時点で、§10a の実装対象、§12 の設定値、§13 の処理フロー、§22 の API レスポンス仕様と整合させる。

**`.build_logs/{id}.json` schema 参照：**

`.build_logs/{id}.json` の保存 key、型、必須条件、Report object、Attempt object、CommitStatus object、BuildMeta object は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0c の `.build_logs/{id}.json` schema を正とする。

runner owner component は、build log の生成タイミング、stdout / stderr 取り込み、`[REPORT]` 変換、WARN 取り込み、secret mask、最終状態保存、書き込み失敗時の後続停止だけを担当する。schema key の追加、削除、型変更、未知 key 保存は本ファイルで行ってはならない。

**ビルドログ最終形契約：**

| 状況 | `finished_at` | `duration_seconds` | `pipeline` | `report` | `deploy` | `snapshot_id` |
|------|---------------|--------------------|------------|----------|----------|---------------|
| GitHub API 失敗 | 保存時刻 | 0 以上 | `exit_code:null`, stdout/stderr 空 | `null` | `[]` | `null` |
| blob decode 失敗 | 保存時刻 | 0 以上 | `exit_code:null`, stdout/stderr 空 | `null` | `[]` | `null` |
| precheck 失敗 | 保存時刻 | 0 以上 | `exit_code:null`, stdout/stderr 空 | `null` | `[]` | `null` |
| pipeline timeout | timeout 検出時刻 | 0 以上 | 取得済み stdout/stderr、`exit_code:null` | parse できた場合のみ object | `[]` | `null` |
| pipeline 非 0 | process 終了時刻 | 0 以上 | 実 exit code と取得済み stdout/stderr | parse できた場合のみ object | `[]` | `null` |
| pipeline 成功 / deploy なし | process 終了時刻 | 0 以上 | `exit_code:0` | object または `null` | `[]` | `null` |
| pipeline 成功 / deploy pending | deploy 判定時刻 | 0 以上 | `exit_code:0` | object または `null` | pending entry | `null` |
| pipeline 成功 / deploy 成功 / snapshot 成功 | snapshot 保存時刻 | 0 以上 | `exit_code:0` | object または `null` | success entry | build id |
| snapshot 失敗 | snapshot 失敗時刻 | 0 以上 | `exit_code:0` | object または `null` | success entry | `null` |

`finished_at` は `started_at` より前にしてはならない。同一 build id のログを複数回保存する場合は、最後の保存が完全 schema を満たすように全 key を含める。途中保存で欠けた key がある状態を最終状態として残してはならない。

**ログ行・secret mask 契約：**

| 項目 | 仕様 |
|------|------|
| stdout/stderr 行長 | 1 行最大 4000 文字。超過分は末尾を切り捨て、`...[truncated]` を付ける。 |
| 配列化 | API 表示用の `stdout` / `stderr` 配列は LF で分割し、空末尾行は除外する。 |
| raw 保存 | `.build_logs/{id}.json.pipeline.stdout` / `stderr` は LF 正規化後の文字列として保存する。 |
| secret mask | 保存前に既知 secret 値を長い順で置換する。空文字 secret は mask 対象にしない。 |
| mask 対象 | `.github_token`、`.webhook_secret`、`.smtp_secret`、`.api_tokens` の有効 token hash 元値は保存しないため対象外、実行時に保持する平文 token、session token、TOTP secret。 |
| WARN 取り込み | stdout 行頭が `[WARN] ` または `[WARNING] ` の行だけを `warnings` に取り込む。stderr は WARN 取り込み対象外。 |
| REPORT 重複 | `[REPORT]` が複数ある場合は最初の 1 行を採用し、`warnings` に `REPORT_DUPLICATE` を追加する。 |
| REPORT parse 失敗 | `report:null` とし、`warnings` に `REPORT_PARSE_FAILED` を追加する。pipeline exit code は変更しない。 |

**`.build_history` JSON Lines 追記契約：**

`.build_history` の保存 key、型、必須条件、許容値は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0c の `.build_history` JSON Lines schema を正とする。

runner は build 結果確定後、`.build_history` へ 1 build につき 1 行だけ追記する。`status` は runner の最終結果、`trigger` は §13 の有効値、`output_sha256` は出力サイト全体 manifest の SHA-256 hex とする。manifest 生成に失敗した場合のみ `output_sha256:null` を許可する。JSON Lines 追記は `O_APPEND|O_CREATE|O_WRONLY` で行い、1 行全体を書き込んでから file sync する。

**固定エラー文言：**

| 条件 | `error` |
|------|---------|
| GitHub 認証失敗 | `github authentication failed` |
| GitHub target 不在 | `github target not found` |
| GitHub API 全再試行失敗 | `github api failed` |
| Blob decode 失敗 | `blob decode failed` |
| Markdown 書き出し失敗 | `source write failed` |
| 事前チェック失敗 | `precheck failed` |
| pipeline timeout | `pipeline timeout` |
| pipeline 非 0 | `pipeline failed` |
| deploy pending | `deploy pending` |
| 状態ファイル書き込み失敗 | `state write failed` |

---

## 15a. `runner` 受け入れ fixture

`runner` の初期実装は、本節の fixture をすべて満たすまで完了として扱わない。fixture ファイルは実装 PR で `testdata/runner/` 配下へ追加する。外部 GitHub API と SSH サーバーへ実接続するテストは初期 fixture に含めず、HTTP test server と fake `ssh` executable で再現する。

### Fixture R1: CLI 異常系

| 実行 | 終了コード | stdout | stderr |
|------|------------|--------|--------|
| `adlaire-ci-runner --help` | `0` | `Usage: adlaire-ci-runner [--state-dir path] [--once] [--dry-run] [--version] [--help]` | 空 |
| `adlaire-ci-runner --state-dir relative` | `2` | 空 | `state directory must be absolute: relative` |
| `adlaire-ci-runner --unknown` | `2` | 空 | `unknown option: --unknown` |

### Fixture R2: 変更なし skip

**前提状態：**

```text
testdata/runner/r2/state/.github_token
testdata/runner/r2/state/.last_sha
testdata/runner/r2/state/.branch_config
```

`.last_sha`:

```json
{"sha":"blob-1"}
```

GitHub Trees API fake response は `target_file=docs` の SHA として `blob-1` を返す。

**期待結果：**

- 終了コード `0`。
- `.last_sha` は変更しない。
- `.build_logs/` に新規 log を作成しない。
- `.build_history` に追記しない。
- stdout slog に INFO ログ `NO_CHANGE: branch=main target=docs sha=blob-1` を 1 件出力する。

### Fixture R3: 変更あり build 成功 deploy なし

**前提状態：**

- `.last_sha` は `{"sha":"old-blob"}`。
- GitHub Trees API fake response は `new-blob` を返す。
- GitHub Blobs API fake response は UTF-8 Markdown を Base64 で返す。
- `deploy_targets` は空配列。
- fake `.ci/pipeline.sh` は終了コード `0` で、stdout に `adlaire-ci-build` の `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default` を出力する。

**期待結果：**

- 終了コード `0`。
- `.last_sha` は `{"sha":"new-blob"}` に atomic write される。
- `.build_logs/{id}.json` が作成され、`target_status="success"`、`pipeline.exit_code=0`、`report.pages=1`、`report.tables_count=0`、`deploy=[]`、`error=null` を含む。
- `.build_history` に同じ `id` の JSON Lines が 1 行追記される。
- `.build_state.running` は終了時 `false`、`current_build_id` は `null`。
- `.build_lock` は終了時に存在しない。

### Fixture R4: pipeline 失敗

**前提状態：**

- GitHub fake response は変更ありを返す。
- fake `.ci/pipeline.sh` は終了コード `7`、stdout `before fail`、stderr `failed` を出力する。

**期待結果：**

- 終了コード `1`。
- `.last_sha` は旧 SHA のまま。
- `.build_logs/{id}.json` は `target_status="failure_build"`、`pipeline.exit_code=7`、`pipeline.stdout="before fail"`、`pipeline.stderr="failed"`、`error="pipeline failed"` を含む。
- `.build_history` に `status="failure_build"` の行を追記する。
- deploy、snapshot は実行しない。

### Fixture R5: deploy pending

**前提状態：**

- pipeline は成功する。
- `deploy_targets` は 1 件。
- fake `ssh` は転送時に終了コード `255` を返す。

**期待結果：**

- 終了コード `1`。
- `.last_sha` は新 SHA に更新する。
- `.pending_transfers` に `out`、`host`、`user`、`dest_dir`、`failed_at`、`retry_count=1` を持つ entry を 1 件追加する。
- `.build_logs/{id}.json` は `target_status="success_deploy_pending"`、`deploy[0].status="pending"`、`deploy[0].transfer_verified=false`、`error="deploy pending"` を含む。
- snapshot は作成しない。

### Fixture R6: lock 競合

`.build_lock` が存在し、`pid` が `/proc/{pid}` に存在する実行中 PID を指す場合、runner は終了コード `0` で終了し、`.build_state`、`.build_logs/`、`.build_history` を変更しない。stdout slog に `BUILD_SKIP: already running (PID {pid})` を出力する。

### Fixture R7: 状態破損

`.notify_pending` が JSON として壊れている場合、runner は `.notify_pending.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、`.notify_pending` を `[]` で再生成する。その後、通常処理を継続する。退避ファイル名の timestamp は UTC とし、秒単位で固定する。

### Fixture R8: pipeline timeout

**前提状態：**

- GitHub fake response は変更ありを返す。
- `.server_config.build_timeout_seconds` は `1`。
- fake `.ci/pipeline.sh` は stdout に `started` を出力後、timeout まで終了しない。

**期待結果：**

- 終了コード `1`。
- `.last_sha` は旧 SHA のまま。
- `.build_logs/{id}.json` は `target_status="failure_build"`、`pipeline.exit_code=null`、`pipeline.stdout="started\n"`、`error="pipeline timeout"` を含む。
- `.build_history` に `status="failure_build"` を 1 行だけ追記する。
- `.build_state.running=false`、`current_build_id=null`、`.build_lock` 不在で終了する。

### Fixture R9: GitHub API 全再試行失敗

**前提状態：**

- GitHub Trees API fake server は retry 対象の HTTP `503` を返し続ける。
- `.last_sha` は `{"sha":"old-blob"}`。

**期待結果：**

- 終了コード `3`。
- `.last_sha` は旧 SHA のまま。
- `.build_logs/{id}.json` は `target_status="failure_api"`、`blob_sha=null`、`pipeline.exit_code=null`、`error="github api failed"` を含む。
- `.build_history` に `status="failure_api"` を 1 行だけ追記する。
- pipeline、deploy、snapshot は実行しない。

### Fixture R10: 通知失敗は build 成功を反転しない

**前提状態：**

- pipeline は成功する。
- `.notify_config` は `on:["success"]` の webhook channel を 1 件持つ。
- fake webhook endpoint は HTTP `500` を返す。

**期待結果：**

- 終了コード `0`。
- `.last_sha` は新 SHA に更新する。
- `.build_logs/{id}.json.target_status` は `"success"`。
- `.notify_pending` に `event="success"`、`retry_count=1` の entry を保存する。
- `.notify_log` に失敗記録を残す。通知失敗を理由に `.build_history.status` を failure にしない。

### Fixture R11: REPORT 重複

**前提状態：**

- pipeline は終了コード `0`。
- stdout に `[REPORT]` 行が 2 行ある。

**期待結果：**

- 終了コード `0`。
- 1 行目の `[REPORT]` だけを `report` に保存する。
- `.build_logs/{id}.json.warnings` に `REPORT_DUPLICATE` を含める。
- `.build_history.warnings` は `[WARN]` 行数に `REPORT_DUPLICATE` 分を加えた値にする。

### Fixture R12: finalizer state write failure

**前提状態：**

- pipeline は成功する。
- `.build_state` の atomic write が finalizer 時だけ失敗する fake filesystem を使用する。

**期待結果：**

- 終了コード `1`。
- `.build_logs/{id}.json` と `.build_history` は成功結果を保存済み。
- `.build_status.json` は最終状態を保存済み。
- ERROR ログ `BUILD_STATE_FINALIZE_FAILED` を出す。
- `.build_lock` は削除を試みる。削除成功/失敗に関わらず、`.build_state.running=false` 保存失敗を正常扱いにしない。

### Fixture R13: token 権限不正

**前提状態：**

- `.github_token` が存在し、mode が `0644`。
- `.build_lock` は存在しない。

**期待結果：**

- 終了コード `2`。
- ERROR ログ `GITHUB_TOKEN_INSECURE_MODE` を出す。
- `.build_state.running` を `true` にしない。
- `.build_logs/` と `.build_history` を作成しない。
- token 値、token 長、token hash を stdout、stderr、状態ファイルへ出力しない。

### Fixture R14: dry-run directory 作成なし

**前提状態：**

- `--state-dir` は既存 directory。
- `{StateDir}/repo`、`{StateDir}/dist`、`{StateDir}/.build_logs`、`{StateDir}/.snapshots` は存在しない。

**実行：**

```bash
adlaire-ci-runner --state-dir <state> --dry-run
```

**期待結果：**

- 終了コード `0`。
- 上記 directory を作成しない。
- stdout は `§27.2` の dry-run JSON 1 件だけを出し、`warnings[]` に `code="DRY_RUN_WOULD_CREATE_DIR"` を対象 directory ごとに出す。
- `.github_token`、`.build_lock`、GitHub API、pipeline、deploy、通知を実行しない。

### Fixture R15: SHA cache 破損

**前提状態：**

- `.last_sha` が `{bad json`。
- GitHub Trees API fake response は `new-blob` を返す。

**期待結果：**

- 終了コード `1`。
- `.last_sha` は変更しない。
- `.build_logs/{id}.json` は `target_status="failure_decode"`、`previous_blob_sha=""`、`error="sha cache invalid"` を含む。
- pipeline、deploy、snapshot は実行しない。

### Fixture R16: GitHub rate limit reset 不正

**前提状態：**

- GitHub Trees API fake server は HTTP `403`、`X-RateLimit-Remaining: 0`、不正な `X-RateLimit-Reset` を返す。
- `.last_sha` は `{"sha":"old-blob"}`。

**期待結果：**

- 終了コード `3`。
- `.last_sha` は旧 SHA のまま。
- `.build_logs/{id}.json` は `target_status="failure_api"`、`error="github api failed"` を含む。
- runner は reset header を無視して長時間待機しない。

### Fixture R17: cooldown skip と manual force

**前提状態：**

- `.build_state.last_finished_at` が現在時刻から `build_cooldown_seconds` 未満。
- `.build_state.queued` は空。

**期待結果 A: polling**

- 終了コード `0`。
- `.build_status.json.status` は `skipped_cooldown`。
- GitHub API、pipeline、deploy、snapshot を実行しない。

**期待結果 B: manual force queue**

- `.build_state.queued[0].trigger="manual"`、`payload.force=true` の場合、cooldown を無視して build を実行する。
- build log / history の `trigger` は `manual`。
- 処理済み queue entry は `.build_state.queued` から削除する。

### Fixture R18: SSH checksum mismatch pending

**前提状態：**

- pipeline は成功する。
- fake `ssh` は転送コマンドを成功させる。
- fake `ssh sha256sum` は local checksum と異なる hash を返す。

**期待結果：**

- 終了コード `1`。
- `.last_sha` は新 SHA に更新する。
- `.pending_transfers` に `last_error` を含む entry を 1 件保存する。
- `.build_logs/{id}.json.deploy[0].transfer_verified=false`、`target_status="success_deploy_pending"`、`error="deploy pending"`。
- snapshot は作成しない。

### Fixture R19: pending 重複統合

**前提状態：**

- `.pending_transfers` に `out`、`host`、`user`、`dest_dir` が同一の entry が 1 件存在する。
- 新規 deploy 失敗も同じ `out`、`host`、`user`、`dest_dir`。

**期待結果：**

- `.pending_transfers` の件数は増えない。
- 既存 entry の `retry_count` が +1 され、`failed_at` と `last_error` が最新値に更新される。
- entry の投入順は保持する。

### Fixture R20: snapshot atomic save and prune

**前提状態：**

- pipeline と deploy は成功する。
- `HISTORY_KEEP_N=2`。
- `.snapshots/` に古い snapshot directory が 2 件存在する。

**期待結果：**

- `{StateDir}/.snapshots/{build_id}` が作成される。
- 一時 directory `{build_id}.tmp.{pid}` は残らない。
- snapshot 内に通常ファイルだけが保存され、`.github_token`、`.build_lock`、`.pending_transfers` を含まない。
- snapshot は 2 件だけ残り、最古 snapshot が削除される。

### Fixture R21: status start write failure

**前提状態：**

- `.build_lock` は存在しない。
- `.build_status.json` の atomic write だけが失敗する fake filesystem を使用する。

**期待結果：**

- 終了コード `1`。
- `.build_state.running` を `true` にしない。
- GitHub API、pipeline、deploy、snapshot を実行しない。
- `.build_logs/` と `.build_history` を作成しない。
- `.build_lock` は削除される。

### Fixture R22: build log write failure

**前提状態：**

- GitHub fake response は変更ありを返す。
- pipeline は成功する。
- `.build_logs/{id}.json` の atomic write だけが失敗する fake filesystem を使用する。

**期待結果：**

- 終了コード `1`。
- `.build_history` は追記しない。
- `.last_sha` は旧 SHA のまま。
- deploy、snapshot は実行しない。
- `.build_status.json` は `status="failure"`、`last_target_status="failure_state_write"`、`last_error="state write failed"` を含む。
- `.build_state.running=false`、`current_build_id=null`、`.build_lock` 不在で終了する。

### Fixture R23: history append failure

**前提状態：**

- GitHub fake response は変更ありを返す。
- pipeline は成功する。
- `.build_logs/{id}.json` は保存成功する。
- `.build_history` の追記だけが失敗する fake filesystem を使用する。

**期待結果：**

- 終了コード `1`。
- `.build_logs/{id}.json` は成功結果を保持する。
- `.last_sha` は旧 SHA のまま。
- deploy、snapshot は実行しない。
- `.build_status.json.last_target_status` は `failure_state_write`。
- `.build_state.running=false`、`.build_lock` 不在で終了する。

### Fixture R24: multi target partial failure continues

**前提状態：**

- `BRANCH_TARGETS` が 2 件。
- 1 件目の GitHub Trees API は HTTP `503` を返し続ける。
- 2 件目は変更あり、pipeline 成功、deploy なし。

**期待結果：**

- 終了コード `1`。
- 1 件目は `.build_logs/{id1}.json.target_status="failure_api"`、`.build_history.status="failure_api"`。
- 1 件目の `sha_file` は更新しない。
- 2 件目は `.build_logs/{id2}.json.target_status="success"`、`.build_history.status="success"`、`sha_file` を更新する。
- runner は 1 件目の失敗で中断しない。

### Fixture R25: all targets GitHub API failure

**前提状態：**

- `BRANCH_TARGETS` が 2 件。
- 両方の GitHub Trees API が retry 対象 HTTP `503` を返し続ける。

**期待結果：**

- 終了コード `3`。
- 各 target の `.build_logs/{id}.json.target_status` は `failure_api`。
- 各 target の `.build_history.status` は `failure_api`。
- すべての `sha_file` は旧値のまま。
- pipeline、deploy、snapshot は実行しない。

### Fixture R26: snapshot failure remains success

**前提状態：**

- GitHub fake response は変更ありを返す。
- pipeline と deploy は成功する。
- snapshot writer だけが失敗する fake filesystem を使用する。

**期待結果：**

- 終了コード `0`。
- `.last_sha` は新 SHA に更新する。
- `.build_logs/{id}.json.target_status="success"`。
- `.build_logs/{id}.json.warnings` に `SNAPSHOT_SAVE_FAILED` を含める。
- `.build_history.status="success"`。
- `.pending_transfers` は追加しない。
- `.build_status.json.status="success"`。

### Fixture R27: build_state finalizer failure keeps failure

**前提状態：**

- pipeline は成功する。
- `.build_logs/{id}.json`、`.build_history`、`.build_status.json` は保存成功する。
- finalizer の `.build_state` atomic write だけが失敗する fake filesystem を使用する。

**期待結果：**

- 終了コード `1`。
- `.build_logs/{id}.json.target_status="success"` と `.build_history.status="success"` は保持する。
- `.build_status.json.status="success"` は保持する。
- ERROR ログ `BUILD_STATE_FINALIZE_FAILED` を出す。
- `.build_lock` は削除を試みる。
- `.build_state.running=false` 保存失敗を正常扱いにしない。

---

## 16. systemd タイマー参照

systemd unit 本文、配置先、起動手順、更新手順は setup owner component の責務とし、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.4.1、§26.5 を正とする。

runner owner component は、`adlaire-ci-runner --state-dir /opt/adlaire-builder` として oneshot 実行された場合の処理、終了コード、状態ファイル更新、ログ出力だけを定義する。

runner 実装は systemd unit file を生成、配置、更新、enable、restart してはならない。systemd 操作が必要な機能は `setup` または `api` owner component の詳細仕様で定義する。

runner が journal へ出力する内容は §15 のログ仕様を正とする。`systemctl`、`journalctl` の操作手順は本ファイルでは定義しない。

---

## 17. GitHub 側設定

| 項目 | 内容 |
|------|------|
| PAT スコープ | `contents: read`（読み取り専用）のみ |
| PAT の種類 | Fine-grained PAT（特定リポジトリのみ許可）を使用する。 |
| Webhook 設定（ポーリング方式） | **不要**（デフォルト。`BRANCH_TARGETS` によるポーリングのみ使用する場合） |
| Webhook 設定（受信方式） | GitHub リポジトリ設定 → Webhooks → Add webhook で `POST /api/webhook` の URL・Secret を設定する（→ §22）。イベントは `push` のみ選択する。**外部公開エンドポイントが必要**（リバースプロキシ経由） |

---

## 18. 初回セットアップ手順参照

初回セットアップ、Release asset 取得、checksum 検証、バイナリ配置、secret 初期化、状態ファイル初期化、systemd unit 書き込み、service 起動、管理 API 導入、管理 UI 配置は setup owner component の責務とし、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.1〜§26.4 を正とする。

runner owner component は、セットアップ済み環境で `/usr/local/bin/adlaire-ci-runner` が起動された後の処理だけを定義する。

runner 実装は以下を行ってはならない。

| 禁止事項 | 理由 |
|----------|------|
| OS user 作成、directory 作成、chown / chmod の初期設定 | setup owner component の責務。 |
| Release asset 取得、checksum 検証、バイナリ配置 | setup owner component の責務。 |
| `.github_token` の新規生成または対話入力 | setup owner component の secret initializer の責務。 |
| `.admin_credentials` 初期化、API service 配置、管理 UI 配置 | api / admin / setup owner component の責務。 |
| systemd unit file の配置、enable、restart | setup owner component の責務。ただし API endpoint が systemd timer を変更する機能は `ADLAIRE_CI_DETAIL_API_SPEC.md` の該当節を正とする。 |

runner が起動時に必要ファイル不足または権限不備を検出した場合は、§13、§15a、§20 の異常系に従い、セットアップ手順を自動実行せずに失敗として記録する。

---

## 19. 管理 API サーバー制限参照

管理 API サーバーの HTTP listener、認証、session、rate limit、TLS 非対応、外部認証非対応、worker pool 非採用の制限は api owner component の責務とし、`ADLAIRE_CI_DETAIL_API_SPEC.md` §21a および `ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` §27.42〜§27.47 を正とする。

runner owner component は、管理 API サーバーの起動、listener、session、認証、HTTP response、rate limit を実装してはならない。

runner と api が同じ状態ファイルを参照する場合でも、runner は API session、API token、TOTP、rate limit、HTTP access log を読み書きしない。runner が読み書きする状態ファイルは §11、§13、§15、§22.0a、§22.0d、および runner owner の個別 §27.x に明記されたものだけとする。

---

## 20. CI ランナー 既知の制限

### 20.1 Go 版 runner の制限

| 制限 | 詳細 |
|------|------|
| ポーリング遅延 | 変更検出は systemd timer の実行間隔に依存する。即時反応が必要な場合は `POST /api/webhook` を併用する。 |
| `pipeline.sh` 起動 | ビルド起動は `src` と同じディレクトリ配下の `.ci/pipeline.sh` を標準とする。YAML 形式のパイプライン定義は本ファイルでは定義しない。 |
| `BRANCH_TARGETS` 直列処理 | 複数エントリはリスト順に順次処理する。並列処理は行わない。1 件の処理が失敗しても、失敗をログと `.build_logs/{id}.json` に記録した上で次エントリへ進む。 |
| GitHub API リトライ | GitHub API 失敗時は `API_RETRY_MAX` 回まで指数バックオフで再試行する。全試行失敗時は ERROR ログを記録し、当該ターゲットのビルドをスキップする。SHA は更新しない。 |
| ビルド失敗時の扱い | `pipeline.sh` が非 0 で終了した場合は ERROR ログを出し、SHA を更新しない。次回実行では同じ blob SHA を再検出して再度ビルド対象になる。 |

### 20.2 標準機能の制限

以下は管理 API または追加拡張と連携する場合の制限である。

| 制限 | 詳細 |
|------|------|
| Webhook 受信の外部公開 | `POST /api/webhook` は `api`（`127.0.0.1` バインド）で受信するため、GitHub から直接受信する構成ではリバースプロキシと TLS 終端が必要。 |
| ペンディングキュー | ペンディング再試行が失敗した場合、`retry_count` を 1 増やしてエントリを保持する。runner による自動放棄は行わない。削除は転送成功時、または管理 API / 手動運用で明示的に削除する場合に限定する。 |
| ペンディングキュー肥大化 | `queue_max_size` を超えた新規投入は ERROR ログを記録し、新規エントリを追加しない。既存エントリは削除しない。 |
| サーキットブレーカー | 連続失敗回数が `API_CIRCUIT_BREAKER_THRESHOLD` 以上になった場合はポーリングを停止し、`POST /api/circuit-breaker/reset` でのみ復帰する。 |

---

---

## 27. Runner owner 追加仕様化機能 詳細仕様

### 27.1 GitHub Commit Status API

本節の主本文は `ADLAIRE_CI_DETAIL_COMMITSTATUS_SPEC.md` §27.1 を正とする。owner component は `commitstatus`、collaborator component は `runner`、`statefile` とする。

runner は、commit SHA 確定、build id 採番、build 開始前の pending 送信呼び出し、pipeline / deploy / snapshot / history の最終結果確定後の final 送信呼び出しだけを担当する。GitHub Commit Status API payload、送信順、送信失敗時の非反転、保存値、secret mask、検証条件は `ADLAIRE_CI_DETAIL_COMMITSTATUS_SPEC.md` §27.1 を正とし、本ファイルへ重複定義してはならない。

### 27.2 ドライラン実行モード

owner component は `runner` とする。collaborator component は `statefile` とする。

`adlaire-ci-runner --dry-run` は、実行計画を検証する読み取り専用モードである。dry-run は `.build_lock`、`.build_state`、`.build_status.json`、`.build_history`、`.build_logs/`、`.pending_transfers`、`.notify_*`、`.snapshots/`、GitHub Commit Status、deploy 先を変更してはならない。

dry-run の stdout は JSON object 1 件と末尾改行に固定する。

```json
{
  "mode": "dry-run",
  "dry_run": true,
  "state_dir": "/opt/adlaire-builder",
  "targets": [
    {
      "branch": "main",
      "target_file": "docs",
      "previous_sha": "old",
      "current_blob_sha": "new",
      "current_commit_sha": "abcdef",
      "would_build": true,
      "trigger": "polling",
      "reason": "sha_changed",
      "warnings": []
    }
  ],
  "would_write": [
    "build_log",
    "history",
    "status",
    "sha_cache"
  ],
  "would_call": ["github_tree", "github_blob"],
  "secrets_masked": true,
  "errors": [],
  "warnings": []
}
```

`targets[].reason` は `"sha_changed"`、`"no_change"`、`"cooldown"`、`"circuit_open"`、`"config_error"`、`"github_error"`、`"precheck_error"` のいずれかとする。`would_write` は実際に書き込んだ path ではなく、非 dry-run 実行で書込対象になる論理種別の予告である。dry-run は `would_write` に値がある場合でも実ファイルを作成、更新、削除してはならない。

**dry-run 出力 schema 固定契約：**

| key | 型 | 仕様 |
|-----|----|------|
| `mode` | string | 常に `"dry-run"`。 |
| `dry_run` | boolean | 常に `true`。 |
| `state_dir` | string | 解決後の state directory を返す。secret、token、URL credential を含めてはならない。 |
| `targets` | array | branch target 正規化後の処理順で返す。複数 target は branch 名昇順、同一 branch は target path 昇順。 |
| `targets[].branch` | string | 対象 branch 名。 |
| `targets[].target_file` | string | 対象 Markdown file または directory。 |
| `targets[].previous_sha` | string|null | dry-run 開始時点の保存済み SHA。存在しない場合は `null`。 |
| `targets[].current_blob_sha` | string|null | fake GitHub read で取得した blob SHA。GitHub read 失敗または precheck failure では `null`。 |
| `targets[].current_commit_sha` | string|null | fake GitHub read で取得した commit SHA。GitHub read 失敗または precheck failure では `null`。 |
| `targets[].would_build` | boolean | 非 dry-run なら build を開始する場合だけ `true`。`github_error`、`config_error`、`precheck_error` では `false`。 |
| `targets[].trigger` | string | 判定 trigger。既定は `"polling"`。 |
| `targets[].reason` | string | `"sha_changed"`、`"no_change"`、`"cooldown"`、`"circuit_open"`、`"config_error"`、`"github_error"`、`"precheck_error"` のいずれか。 |
| `targets[].warnings` | array[object] | target 固有 warning。形式は `warnings[]` と同じ。 |
| `would_call` | array[string] | 実行予定の外部 API / command 種別だけを固定文字列で返す。secret、URL query 全体、token は含めない。 |
| `would_write` | array[string] | 非 dry-run なら発生する論理書込予定を固定文字列で返す。許可値は `"lock"`、`"build_log"`、`"history"`、`"status"`、`"sha_cache"`、`"notification"`、`"deploy"`、`"snapshot"`、`"commit_status"`。実書込が発生したことを意味しない。 |
| `secrets_masked` | boolean | stdout、stderr、expected、effects に secret 平文が残らない検証を通した場合だけ `true`。 |
| `errors` | array[object] | `{ "code": string, "message": string, "target": string|null }`。 |
| `warnings` | array[object] | `{ "code": string, "message": string, "target": string|null }`。 |

`would_call` の許可値は `"github_tree"`、`"github_blob"`、`"github_commit"`、`"github_rate_limit"` に限定する。dry-run は fake GitHub read だけを実行対象にし、実 GitHub write API、GitHub Commit Status、SSH、pipeline、deploy、snapshot、notification、systemd、hook、remote build を `would_call` に含めてはならない。

dry-run は、破損 state の backup、初期値作成、lock 作成、通知、GitHub Commit Status、deploy、archive、cleanup、quarantine、状態正規化を実行しない。stdout 以外の状態差分が発生した場合は実装不合格とする。

終了コードは下表に固定する。`--help` / `--version` と同時指定された場合は `--help` / `--version` を優先し、dry-run JSON を出力しない。

| 終了コード | 条件 | stdout / stderr | 副作用 |
|------------|------|-----------------|--------|
| `0` | dry-run 検証が完了した。`would_build=true`、`would_build=false`、warning あり、cooldown、circuit open を含む。 | stdout に dry-run JSON 1 件。stderr は空、または環境依存でなく secret mask 済み warning だけ。 | 状態 / log / cache / lock / 通知 / deploy / status 差分なし。 |
| `2` | CLI 引数不正、設定不正、設定破損、path 不正、schema 不正、precheck failure。GitHub read 前に確定する failure。 | stdout に dry-run JSON 1 件、`errors[]` に原因、該当 target の `reason="config_error"` または `"precheck_error"`。stderr は secret mask 済み。 | backup、初期化、正規化、quarantine、rewrite を含めて差分なし。 |
| `3` | fake GitHub read の最終失敗、rate limit、network failure、GitHub response schema 不正。 | stdout に dry-run JSON 1 件、`errors[]` に原因、該当 target の `reason="github_error"`。stderr は secret mask 済み。 | pipeline、deploy、notification、commit status、状態更新を含めて差分なし。 |

検証条件:

| ケース | 期待結果 |
|--------|----------|
| SHA 差分あり | `would_build=true`、`reason="sha_changed"`、`would_write` に非 dry-run 時の論理書込予定、状態ファイル差分なし。 |
| SHA 差分なし | `would_build=false`、`reason="no_change"`。 |
| cooldown | `would_build=false`、`reason="cooldown"`。 |
| fake GitHub read 失敗 | 終了コード `3`、`reason="github_error"`、`errors[]` に理由、状態ファイル差分なし、pipeline / deploy / notification / commit status 呼び出しなし。 |
| 設定破損 | 終了コード `2`、`reason="config_error"`、破損ファイルの退避、再生成、正規化、quarantine、rewrite を行わない。 |
| 複数 target | branch / target path の固定順で返る。 |
| secret 設定済み | stdout JSON、stderr、expected、effects に secret 平文が出ず、`secrets_masked=true`。 |

### 27.3 ビルド失敗時の自動リトライ

owner component は `runner` とする。collaborator component は `statefile` とする。

runner は `.server_config.build_retry_max > 0` の場合、retry 対象失敗だけを同一 build id 内で最大 `build_retry_max` 回追加試行する。総試行回数は `1 + build_retry_max` とする。

retry 対象は以下に限定する。

| 対象 | 条件 |
|------|------|
| GitHub API | network error、HTTP 429、HTTP 500〜599。 |
| pipeline | timeout。exit code 非 0 は retry しない。 |
| deploy | SSH 接続失敗、検証用 checksum 取得失敗、network timeout。checksum mismatch は retry しない。 |

retry 待機秒数は `build_retry_base_seconds * attempt` とする。初回失敗後の retry 1 回目は `base * 1`、retry 2 回目は `base * 2`。待機中に process が終了した場合、未完了 retry を再開してはならない。

attempt ごとの結果は `.build_logs/{id}.json.attempts[]` に必ず保存する。最終 attempt が成功した場合、`.build_history.status` は `"success"` とし、`retry_count` に追加 retry 回数を保存する。全 attempt 失敗時は `"failure"` とする。SHA 更新、snapshot、deploy 成功記録は最終成功時だけ行う。

**attempt schema / 更新固定契約：**

| key | 型 | 仕様 |
|-----|----|------|
| `attempt` | integer | 初回を `1` とする。 |
| `started_at` / `finished_at` | string | UTC ISO 8601 秒精度。 |
| `target_status` | string | §13 の固定値。 |
| `retryable` | boolean | 次 attempt の対象なら `true`。 |
| `error` | string/null | 固定文言。secret、token、command 全文は含めない。 |

retry 待機中に SIGTERM、context timeout、lock 喪失を検出した場合は待機を中断し、未実行 attempt を作成せず、現在までの attempts だけを保存する。最終成功時だけ `.last_sha`、snapshot、deploy success、history success を確定する。途中失敗 attempt で SHA cache を更新してはならない。

検証条件:

| ケース | 期待結果 |
|--------|----------|
| API 429 後成功 | attempts 2 件、history success、retry_count 1。 |
| pipeline timeout 後成功 | attempts 2 件、SHA は最終成功後のみ更新。 |
| pipeline exit 1 | retry なし、history failure、retry_count 0。 |
| deploy checksum mismatch | retry なし、pending transfer 記録。 |
| retry 上限到達 | history failure、attempts は `1 + build_retry_max` 件。 |
| retry 中断 | 未実行 attempt を作らず終了する。 |
| secret error | attempts[].error に secret 平文が出ない。 |

### 27.8 ビルドステータスファイル出力

本機能の目的は、runner の現在状態と直近結果を `.build_status.json` に集約し、api、sdk、ui が同じ read-only 情報を参照できるようにすることである。

owner component は `runner` とする。collaborator component は `api`、`statefile` とする。`runner` は `.build_status.json` の唯一の通常更新責務を持つ。`api` は `GET /api/status`、`GET /api/dashboard`、`GET /api/health` で read-only 参照する。api は `.build_status.json` を自動修復してはならない。

**入力：**

| 入力 | 説明 |
|------|------|
| runner 起動状態 | lock 取得、起動時整合性チェック結果、target 処理結果、pending transfer / notify 件数、circuit 状態。 |
| `.build_state` | `running`、`current_build_id`、`last_started_at`、`last_finished_at`、`queued`。 |
| `.build_history` | 直近 build id、status、trigger、duration。 |
| `.pending_transfers` | pending transfer 件数。 |
| `.notify_pending` | pending notify 件数。 |
| `.build_circuit_state` | circuit open 状態、連続失敗回数。 |

**出力：**

`.build_status.json` は §22.0a / §22.0c の schema に従う JSON object とする。文字コードは UTF-8、改行は末尾 1 つ、ファイル mode は `600` とする。更新は同一ディレクトリ一時ファイルへの書き込み、`fsync`、`os.Rename`、親ディレクトリ `fsync` の順で atomic write する。

**更新タイミング：**

| タイミング | 必須値 |
|------------|--------|
| 起動時整合性チェックで復旧または停止が発生した直後 | `status="warning"` または `"failure"`、`last_trigger="startup_config_integrity"`、`running=false`。 |
| build 開始前 | `status="running"`、`running=true`、`current_build_id`、`last_trigger`、`last_started_at` を保存する。 |
| 変更なし skip | `status="skipped"`、`last_target_status="skipped_no_change"`、`running=false`。build id は更新しない。 |
| cooldown skip | `status="skipped"`、`last_target_status="skipped_cooldown"`、`running=false`。 |
| build 成功 | `status="success"`、`last_build_id`、`last_finished_at`、`last_duration_seconds`、`output_sha256` を保存する。 |
| deploy pending | `status="success"`、`last_deploy_status="pending"`、`pending_transfers_count` を保存する。 |
| build 失敗 | `status="failure"`、`last_error`、`last_target_status`、`last_finished_at` を保存する。 |
| runner finalizer | `running=false`、`current_build_id=null` を必ず保存する。 |

`status` の許容値は `"none"`、`"running"`、`"success"`、`"failure"`、`"skipped"`、`"warning"` に固定する。`last_target_status` は §13 の `target_status` 値、または `null` とする。`last_deploy_status` は `"success"`、`"failure"`、`"pending"`、`"skipped"`、`"none"`、`null` のいずれかとする。

**異常系：**

| 条件 | 処理 |
|------|------|
| `.build_status.json` 書き込み失敗 | ERROR ログ `BUILD_STATUS_WRITE_FAILED: path={path} error={reason}` を出し、runner 終了コードを最低 `1` にする。build 成功後に発生した場合も終了コードは `1` とする。 |
| `.build_status.json` 破損を API が検出 | `GET /api/status` と `GET /api/dashboard` は `500` を返す。`GET /api/health` は `status="degraded"` を返し、破損を `checks[]` に含める。 |
| `.build_status.json` 不在 | API は `.build_state`、`.build_history`、`.build_lock` から後方互換値を算出して返す。ただしファイル作成はしない。 |
| pending 件数読取失敗 | 件数を `null` にせず `0` として返してはならない。status 書き込み時はエラー扱いにし、`last_error` に固定文言を保存する。 |

**セキュリティ：**

`.build_status.json` に GitHub PAT、Webhook URL secret、SMTP password、API token、session token、request body を保存してはならない。`last_error` は最大 500 文字に切り詰め、改行は `\n` 文字列へ escape する。

**status 算出・不一致固定契約：**

| 項目 | 仕様 |
|------|------|
| `running=true` 不一致 | `.build_lock` が存在しないのに `running=true` の場合、API は `degraded` として返し、自動修復しない。 |
| pending 件数 | `.pending_transfers` と `.notify_pending` は読めた場合だけ件数を返す。読取失敗は `checks[]` に記録する。 |
| history 不一致 | `.build_status.json.last_build_id` と最新 history id が異なる場合、API は status file を優先し、warnings に `history_status_mismatch` を含める。 |
| finalizer | finalizer は既存 `last_build_id`、`last_finished_at` を消さず、`running=false` と `current_build_id=null` だけを最低更新する。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| build 開始 | `status="running"`、`running=true`、`current_build_id` が保存される。 |
| build 成功 | `status="success"`、`running=false`、`last_build_id` と `last_duration_seconds` が保存される。 |
| 変更なし | build log / history を作らず、`.build_status.json` は `skipped_no_change` を保持する。 |
| 起動時整合性復旧 | `last_trigger="startup_config_integrity"`、復旧内容が `last_error` または warning として確認できる。 |
| 書込失敗 | runner 終了コードが最低 `1`、ERROR ログが出る。 |
| API read | `GET /api/status`、`GET /api/dashboard` が `.build_status.json` を第一参照元にする。 |
| running stale | API が自動修復せず degraded を返す。 |
| finalizer | `last_build_id` を消さない。 |

### 27.9 ビルドトリガー種別の記録

本機能の目的は、runner がなぜ build または関連処理を開始したかを、履歴、ログ、状態、api、sdk、ui で同一の固定値として扱うことである。

owner component は `runner` とする。collaborator component は `api`、`sdk`、`ui`、`statefile` とする。`runner` は trigger の確定と永続化を担当し、api / sdk / ui は既存値の表示と filter のみを担当する。

**trigger 許容値：**

| 値 | 発生条件 | 補足 |
|----|----------|------|
| `polling` | systemd timer 等の通常起動で SHA 差分がある。 | 既定の自動ビルド。 |
| `force_interval` | SHA 差分なし、かつ `force_build_interval_hours` 条件を満たす。 | 手動 force には使わない。 |
| `manual` | `POST /api/build` または `POST /api/build/force` 由来の queue entry を処理する。 | force は `payload.force=true` で表す。 |
| `webhook` | `POST /api/webhook` 由来の queue entry を処理する。 | 署名検証成功済み event のみ。 |
| `retry_pending_transfer` | `.pending_transfers` の再送のみを実行する。 | 通常 build とは別 trigger。 |
| `startup_config_integrity` | 起動時整合性チェックで復旧、正規化、または停止が発生する。 | build log / history は作成しない。 |
| `rollback` | `POST /api/history/{id}/rollback` により snapshot を再転送する。 | 新しい build id を作成する。 |
| `local_watch` | `watch_mode="local"` の local SHA 差分により build する。 | GitHub API を呼ばない。 |
| `approval` | `POST /api/approvals/{id}/approve` 由来の queue entry を処理する。 | 承認済み entry のみ。 |

上表以外の値を保存、返却、表示してはならない。特に `"auto"`、`"force"`、`"scheduled"`、`"timer"` は使用禁止とする。

**保存先：**

| 保存先 | 必須条件 |
|--------|----------|
| `.build_logs/{id}.json.trigger` | build log を作成する全処理で必須。 |
| `.build_history.trigger` | build history を追記する全処理で必須。 |
| `.build_status.json.last_trigger` | build、skip、復旧、rollback の最終 trigger を保存する。 |
| Queue entry `trigger` | `"manual"`、`"webhook"`、`"approval"` のみ許可する。 |
| API response | `StatusObject.last_trigger`、`HistoryRecord.trigger`、`DashboardObject.status.last_trigger` で同じ値を返す。 |

**判定順序：**

1. 起動引数が `--help` または `--version` の場合、trigger を確定しない。
2. 起動時整合性チェックで復旧、正規化、停止が発生した場合、`startup_config_integrity` を `.build_status.json` に記録する。
3. queue entry が存在する場合、entry の `trigger` を採用する。
4. pending transfer の再送だけで終了する起動は `retry_pending_transfer` とする。
5. `watch_mode="local"` で local SHA 差分がある場合は `local_watch` とする。
6. SHA 差分がある通常起動は `polling` とする。
7. SHA 差分がなく force interval 条件を満たす場合は `force_interval` とする。
8. rollback API が作成する処理は `rollback` とする。

複数条件が同時に成立した場合は、上記順序で最初に該当した trigger を採用する。1 回の runner 起動で複数 branch target を処理する場合、target ごとに同じ trigger を保存する。ただし queue entry が target を指定する場合は、対象 target のみにその trigger を適用する。

**api / sdk / ui：**

`GET /api/history` の `trigger` query は上表の値だけを受け付ける。不正値は `422` を返す。SDK `getHistory({trigger})` は値を変換せず送信する。ui は filter の選択肢を上表の 9 件に固定し、未知 trigger を受け取った場合は `Unknown` へ丸めず、該当行に `invalid trigger` エラーを表示する。

**異常系：**

| 条件 | 処理 |
|------|------|
| queue entry の trigger が不正 | queue entry を処理せず ERROR ログ `INVALID_TRIGGER: id={id} trigger={value}`、HTTP API 由来なら queue 作成時に `422`。 |
| 既存 history に未知 trigger がある | API はその行を返すが、`warnings[]` に `unknown_trigger` を含める。新規保存では未知値を禁止する。 |
| build log と history の trigger 不一致 | API は `500` を返し、server log に `TRIGGER_MISMATCH: id={id}` を出す。 |

**trigger 保存・queue 固定契約：**

| 項目 | 仕様 |
|------|------|
| queue 取り出し | queue entry の trigger は取り出し時に validation し、不正なら entry を削除せず処理を中断する。 |
| rollback | rollback API は queue を経由しない場合でも build log / history に `rollback` を保存する。 |
| startup | `startup_config_integrity` は build id を採番しない。 |
| unknown 既存値 | API response に既存値をそのまま返し、UI が検出できるよう warnings を付ける。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| polling build | log、history、status が `polling` で一致する。 |
| manual force | queue は `manual`、payload は `force=true`、保存 trigger は `manual`。 |
| webhook | 署名検証成功時だけ `webhook` が保存される。 |
| force interval | SHA 差分なしで `force_interval` が保存される。 |
| filter | `GET /api/history?trigger=manual` が manual のみ返す。 |
| 不正 trigger | queue 作成または history filter が `422`。 |
| startup trigger | build log / history を作らない。 |
| mismatch | API は `500`。 |

### 27.10 設定ファイル起動時整合性チェック

本機能の目的は、runner が build 処理に入る前に、runner が読む状態ファイルの破損、型不一致、必須 key 不足、権限不備を検出し、規定どおり復旧または停止することである。

owner component は `runner` とする。collaborator component は `statefile` とする。管理 API、sdk、ui は本機能の実行責務を持たない。api が同じ状態ファイルを読む場合も、起動時整合性チェックを代行してはならない。

**対象ファイル：**

| 順序 | ファイル | 不在時 | 破損時 | unknown key |
|------|----------|--------|--------|-------------|
| 1 | `.branch_config` | 作成せず default 採用 | backup 後、不在扱い | 除去して正規化 |
| 2 | `.notify_config` | 初期値作成 | backup 後、初期値作成 | 除去して正規化 |
| 3 | `.build_state` | 初期値作成 | backup 後、初期値作成 | 除去して正規化 |
| 4 | `.build_circuit_state` | 初期値作成 | backup 後、初期値作成 | 除去して正規化 |
| 5 | `.pending_transfers` | `[]` 作成 | backup 後、`[]` 作成 | 除去して正規化 |
| 6 | `.notify_pending` | `[]` 作成 | backup 後、`[]` 作成 | 除去して正規化 |

**実行順序：**

1. CLI 引数を検証する。
2. `--help` または `--version` の場合は本チェックを実行しない。
3. `--dry-run` の場合は検証だけ行い、backup、初期化、正規化、通知、状態更新を行わない。
4. `StateDir` が絶対パスかつ既存ディレクトリであることを確認する。
5. `.build_lock` を取得する。
6. 上表の順序で対象ファイルを検証する。
7. 復旧可能な問題は backup、初期化、正規化を行う。
8. 復旧結果に応じて `.build_status.json.last_trigger="startup_config_integrity"` を保存する。
9. 復旧通知条件を満たす場合は `config_corrupt` 通知を 1 回だけ送信する。
10. 停止条件がなければ pending retry、cooldown、target 処理へ進む。

**判定分類と停止条件：**

| 分類 | 処理 | runner 終了コード |
|------|------|------------------|
| `missing_optional` | `.branch_config` を作らず default 採用。 | 継続、最終結果に従う |
| `missing_required` | 初期値を atomic write。 | 継続、最終結果に従う |
| `parse_error` | corrupt backup 後、ファイル別復旧。 | 継続、最終結果に従う |
| `top_level_type_mismatch` | corrupt backup 後、ファイル別復旧。 | 継続、最終結果に従う |
| `required_key_missing` | corrupt backup 後、ファイル別復旧。 | 継続、最終結果に従う |
| `required_key_type_mismatch` | corrupt backup 後、ファイル別復旧。 | 継続、最終結果に従う |
| `invalid_value` | corrupt backup 後、ファイル別復旧。 | 継続、最終結果に従う |
| `unknown_key` | backup せず未知 key を除去して atomic write。 | 継続、最終結果に従う |
| `permission_error` | 自動復旧しない。`.build_state.running` を変更しない。 | `2` |
| `io_error` | 自動復旧しない。`.build_state.running` を変更しない。 | `1` |

backup 名は `{original}.corrupt.{YYYYMMDDHHMMSS}.bak` とする。UTC 秒単位で衝突する場合は `{original}.corrupt.{YYYYMMDDHHMMSS}.{n}.bak` とし、`n` は `2` から始める。

**状態更新禁止事項：**

本チェックだけで `.build_logs/{id}.json` と `.build_history` を作成してはならない。`startup_config_integrity` は `.build_status.json` の `last_trigger` にだけ記録する。ただし、本チェック後に通常 build が発生する場合、通常 build の log / history は実際の build trigger を保存する。

`--dry-run` では、破損検出結果を dry-run JSON の `errors[]` または `warnings[]` に出力するだけとし、backup、初期化、正規化、通知、`.build_status.json` 更新を行わない。

**復旧通知：**

復旧通知は `.notify_config` の検証完了後、復旧対象に `.notify_config` と `.notify_pending` 以外のファイルが 1 件以上ある場合だけ送信する。送信イベントは `config_corrupt` とする。`.notify_pending` が破損復旧された場合、通知失敗時の pending 追記は行わない。

**復旧 record 固定契約：**

| 項目 | 仕様 |
|------|------|
| status warning | 復旧して継続する場合、`.build_status.json.status="warning"`、`last_error` に固定文言を保存する。 |
| status failure | permission / io error で停止する場合、可能なら `.build_status.json.status="failure"` を保存する。保存不能でも終了コードを優先する。 |
| backup 内容 | 破損元 file を byte 単位でコピーする。整形、マスク、改行変換を行わない。 |
| unknown key 正規化 | unknown key のみの場合は corrupt backup を作らない。 |
| 通知 payload | `{event:"config_corrupt", files:[...], recovered:boolean}`。secret 値と破損内容は含めない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 必須ファイル不在 | 初期値作成、終了コード `0`、追加 backup なし。 |
| `.branch_config` 破損 | backup 後に `.branch_config` 不在、default 採用。 |
| `.build_state` 破損 | backup、初期値作成、`startup_config_integrity` 記録。 |
| unknown key | backup なしで正規化、2 回目起動では追加 WARN なし。 |
| permission error | 自動復旧なし、終了コード `2`、running 未変更。 |
| dry-run | 差分なし、backup なし、dry-run JSON に検出結果。 |
| 通知失敗 | `.notify_pending` が正常な場合だけ pending 追記。 |
| backup byte | 破損元と backup が byte 単位一致。 |
| unknown only | backup なしで正規化。 |

### 27.14 ビルド所要時間の記録と統計 API

本機能の目的は、build ごとの開始・終了・所要時間を構造化ログへ保存し、統計 API で直近 N 件の平均、最小、最大を返すことである。

owner component は `runner` とする。collaborator component は `api`、`statefile`、`archive` とする。

**記録仕様：**

`runner` は `.build_logs/{id}.json` に `started_at`、`finished_at`、`duration_seconds` を必ず保存する。`started_at` は build id 採番直後、`finished_at` は最終 target status 確定直後とする。`duration_seconds` は `finished_at - started_at` を秒単位で切り上げず整数化し、1 秒未満は `0` とする。

`.build_history.duration_seconds` は `.build_logs/{id}.json.duration_seconds` と同じ値にする。失敗、deploy pending、rollback でも記録する。変更なし skip で build log を作らない場合は記録しない。

**統計 API：**

`GET /api/stats/build-duration?n=N` は `.build_logs/` と `.build_logs/archive/` を読み、完了済み build log の `duration_seconds` が `null` でない最新 N 件を集計する。`N` は 1〜1000、既定値 20 とする。

Response は `BuildDurationStats` とし、`count=0` の場合は `avg_seconds`、`min_seconds`、`max_seconds` を `null`、`recent` を `[]` とする。

**duration stats 固定契約：**

| 項目 | 仕様 |
|------|------|
| latest 判定 | `finished_at` 降順、同時刻は build id 昇順で最新 N 件を選ぶ。 |
| avg | 合計 / count を小数第 3 位で四捨五入し小数第 2 位まで返す。 |
| recent | `{build_id, finished_at, duration_seconds, status}` を返す。 |
| archive 優先 | 同じ id が通常 log と archive にある場合、通常 log を採用する。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| `n` 不正 | `422`。 |
| build log 破損 | 対象 log を除外し、server log に WARN。 |
| archive gzip 展開失敗 | 対象 log を除外し、server log に WARN。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 成功 build | log/history に同一 duration。 |
| 失敗 build | duration を記録する。 |
| 統計対象なし | `count=0`、平均/最小/最大 `null`。 |
| archive 含む | 通常 log と archive log を横断集計する。 |
| duplicate id | 通常 log を優先し二重集計しない。 |

### 27.19 週次ビルドサマリー Webhook

本機能の目的は、過去 7 日間の build 結果を指定曜日・時刻に集計し、Webhook へ定期通知することである。

owner component は `runner` とする。collaborator component は `api`、`statefile` とする。runner は自動送信、api は設定表示・手動送信を担当する。

**設定：**

`.notify_config.summary.enabled=true` の場合だけ有効とする。`interval` は `"weekly"`、`hour` は 0〜23、`day_of_week` は 0〜6 とする。タイムゾーンは UTC 固定。

**自動送信条件：**

runner 起動時に、現在 UTC の曜日と時が設定値に一致し、`.build_state.weekly_summary_sent_date` が当日でない場合に送信する。送信成功時だけ `weekly_summary_last_sent_at` と `weekly_summary_sent_date` を更新する。

**集計対象：**

`.build_history` のうち、現在時刻から過去 7 日以内の行を対象とする。`status="success"` を成功、`"failure"`、`"cancelled"`、`"hook_error"` を失敗として数える。所要時間は `duration_seconds != null` の行だけ平均対象にする。

**通知 payload：**

```json
{
  "event": "weekly_summary",
  "period_days": 7,
  "success_count": 10,
  "failure_count": 2,
  "success_rate": 83.33,
  "avg_duration_seconds": 42
}
```

**手動送信 API：**

`POST /api/notify/weekly-summary` は同じ集計を即時送信する。手動送信は `weekly_summary_sent_date` を更新しない。

**送信・状態更新固定契約：**

| 操作 | 更新順 | 失敗時 |
|------|--------|--------|
| 自動 weekly summary | 集計 → 対象 channel 抽出 → 通知送信 → `.notify_log` 追記 → `.build_state.weekly_summary_last_sent_at` / `weekly_summary_sent_date` 保存 | 送信または `.notify_log` 追記失敗時は sent date を更新しない。build status は変更しない。 |
| 手動 weekly summary | 認証 → 対象 channel 抽出 → 集計 → 通知送信 → `.notify_log` 追記 → response | 宛先なしは `422`。送信失敗は `500`、sent date は更新しない。 |
| 同日二重自動 | `.build_state.weekly_summary_sent_date` が現在 UTC 日付と一致する場合は送信しない | `.notify_log`、`.notify_pending`、`.build_state` を変更しない。 |

weekly summary payload は secret、repository token、SMTP password、Webhook secret、API token、session token を含めてはならない。`success_rate` は小数第 2 位まで `math.Round(x*100)/100` 相当で丸める。集計対象 0 件の場合は `success_count=0`、`failure_count=0`、`success_rate=0`、`avg_duration_seconds=null` とする。

**weekly 集計固定契約：**

| 項目 | 仕様 |
|------|------|
| 期間 | `now - 7*24h <= finished_at <= now`。timezone は UTC。 |
| 成功 | `status="success"`、`"success_deploy_pending"`。 |
| 失敗 | `status="failure"`、`"cancelled"`、`"hook_error"`、`"failure_remote_build"`。 |
| 除外 | `skipped_*`、`approval_*`、duration 欠落の平均対象。 |
| 最大 duration | payload に `max_duration_seconds`、`max_duration_build_id` を含める。対象なしは `null`。 |
| 手動 response | 送信 payload と送信結果 `{sent:boolean, channel_results:[]}` を返す。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 条件一致 | Webhook 送信、`.notify_log` 追記、sent date 更新。 |
| 同日二重起動 | 2 回目は送信しない。 |
| 宛先なし | 自動送信は WARN、手動 API は `422`。 |
| 手動送信 | payload を返し、sent date は変更しない。 |
| 送信失敗 | sent date を更新せず、retry 対象なら `.notify_pending` に追加。 |
| deploy pending | 成功として数える。 |
| skipped | 集計から除外する。 |

### 27.21 複数ファイル監視

本機能の目的は、単一 `target_file` 前提を拡張し、複数 Markdown ファイルまたは Markdown ディレクトリを 1 回の runner 起動で監視、差分判定、ビルド対象決定できるようにすることである。

owner component は `runner` とする。collaborator component は `builder`、`api`、`statefile` とする。runner は差分検出と build target 決定、builder は複数入力の静的サイト生成、api は設定表示・更新を担当する。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 key | `.branch_config.branch_targets[].target_files` |
| 型 | string 配列。要素は UTF-8 path。 |
| 既定値 | 既存 `target_file` を 1 要素配列へ正規化する。 |
| 上限 | branch target ごとに 100 件。 |
| 許容 path | 相対 path のみ。空文字、絶対 path、`..`、NUL、改行は禁止。 |
| SHA cache | `.sha_cache/{branch}/{target_hash}.sha` に target file 単位で保存する。 |
| build log | `.build_logs/{id}.json.changed_targets[]` に `{target_file,before_sha,after_sha}` を保存する。 |

**正常系：**

1. runner 起動時に `.branch_config` を読み、`target_file` と `target_files` を正規化する。
2. `target_files` を辞書順に重複排除する。
3. 各 target の GitHub content SHA または local SHA-256 を取得する。
4. SHA cache と比較し、変更 target だけ `changed_targets` に追加する。
5. `changed_targets` が空で force 条件もない場合は `skipped_no_change` とする。
6. 1 件以上変更がある場合は builder に `--src` として branch target の `src` を渡し、対象一覧を `ADLAIRE_CHANGED_TARGETS` 環境変数の JSON array で渡す。
7. build 成功時だけ対象 target の SHA cache を更新する。

**target_files 正規化・SHA cache 固定契約：**

| 項目 | 仕様 |
|------|------|
| 正規化順 | `target_file` を 1 要素配列化 → `target_files` と統合 → `/` 区切りへ変換 → `.` segment 除去 → 重複除去 → 辞書順 sort。 |
| 重複判定 | 大文字小文字を区別する。`docs/a.md` と `Docs/a.md` は別 target として扱う。 |
| `target_hash` | 正規化済み target path の SHA-256 hex 先頭 32 文字。 |
| SHA cache path | branch 名と `target_hash` を URL encode せず、branch 名は `/` を `_` に置換して `.sha_cache/{branch_safe}/{target_hash}.sha` に保存する。 |
| force build | force 条件では `changed_targets` が空でも build を実行し、`changed_targets=[]` を build log に保存する。SHA cache は build 成功時に全 target 分を更新する。 |
| 部分失敗 | 1 target でも SHA 取得に最終失敗した場合、build は開始せず、成功取得済み target の SHA cache も更新しない。 |
| 環境変数 | `ADLAIRE_CHANGED_TARGETS` は JSON array string。要素順は正規化済み target path の辞書順。secret は含めない。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| `target_files` が空 | `422`。API 保存不可。runner は設定エラーで終了コード `2`。 |
| target が GitHub 上に存在しない | 該当 target を `missing` として build log に記録し、全体 status は `failure`。SHA cache は更新しない。 |
| 一部 target の SHA 取得失敗 | retry 対象。最終失敗時は build 実行しない。 |
| SHA cache 破損 | 該当 target は変更ありとして扱い、成功時に上書きする。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 1 件変更 | `changed_targets` が 1 件、該当 SHA cache だけ更新。 |
| 複数変更 | 辞書順で記録、build は 1 回だけ実行。 |
| 変更なし | build なし、status `skipped_no_change`。 |
| 不正 path | API は `422`、runner は終了コード `2`。 |
| force build | 変更なしでも build 実行、成功時に全 SHA cache 更新。 |
| SHA 部分失敗 | build なし、SHA cache 差分なし。 |

### 27.22 ビルドパイプライン YAML 定義

本機能の目的は、固定 `pipeline.sh` 依存をなくし、内製 YAML subset で build step を明示定義できるようにすることである。

owner component は `runner` とする。collaborator component は `api`、`statefile` とする。外部 YAML ライブラリは使用しない。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定ファイル | `.pipeline.yml` または `.pipeline_config.inline_yaml` |
| API | `GET /api/pipeline-config` / `POST /api/pipeline-config` |
| YAML root | `version: 1`、`steps:` のみ許可。 |
| step key | `name`、`phase`、`command`、`args`、`env`、`timeout_seconds`、`required`。 |
| phase | `"precheck"`、`"build"`、`"test"`、`"deploy"`、`"post"`。 |
| command | 絶対 path または PATH 解決可能なコマンド名。shell 文字列は禁止。 |
| args | string 配列。空文字、NUL、改行は禁止。 |
| env | string:string object。key は `^[A-Z_][A-Z0-9_]{0,63}$`。 |
| timeout_seconds | 1〜86400。省略時は `build_timeout_seconds`。 |

**YAML subset：**

対応する構文は、2 space indent、string scalar、integer scalar、boolean scalar、string array、object array のみとする。anchor、alias、複数 document、flow style、tag、複数行 string、コメント行以外の inline comment は禁止する。禁止構文を検出した場合は parse error とする。

**`.pipeline_config` 適用固定契約：**

`.pipeline_config` schema は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0c を正とする。API による保存、request / response、HTTP status は `ADLAIRE_CI_DETAIL_API_SPEC.md` §15D を正とする。

runner は build 開始後、builder command または pipeline step command を組み立てる直前に `.pipeline_config` を 1 回だけ読む。同一 build 中に `.pipeline_config` を再読込してはならない。

| 項目 | 仕様 |
|------|------|
| 読込タイミング | build id 採番、target 確定、`running=true` 保存後、builder / step command 組み立て直前。 |
| `extra_args` | legacy builder command を使う場合にだけ、固定引数の後ろへ配列順で追加する。YAML step command には追加しない。 |
| `env` | runner 基本 env → branch env → `.pipeline_config.env` → YAML step env の順で上書きする。 |
| 読込不能 | build 本体を開始せず `failure_pipeline_config`、終了コード `2`。SHA cache、deploy、snapshot は更新しない。 |
| schema 不正 | build 本体を開始せず `failure_pipeline_config`、終了コード `2`。SHA cache、deploy、snapshot は更新しない。 |
| secret mask | `.pipeline_config.env` の secret key 値は stdout/stderr、hook log、notify payload、pipeline step log へ保存前に mask する。 |

`.pipeline_config.extra_args` は、builder の固定引数である `--src`、`--out`、`--build-id`、`--commit-sha`、`--build-at`、`--version`、`--help` を上書きまたは追加してはならない。禁止引数を検出した場合は build 本体を開始せず `failure_pipeline_config` とする。

**正常系：**

1. `.pipeline.yml` があれば優先し、なければ `.pipeline_config.inline_yaml` を使用する。
2. YAML subset parser で `PipelineConfig` に変換する。
3. step を定義順に実行する。
4. `required=false` の step 失敗は WARN として継続し、build log に `optional_failed` を記録する。
5. `required=true` または省略 step の失敗は build を中断する。
6. 各 step の stdout/stderr、exit_code、duration_seconds を `.build_logs/{id}.json.pipeline_steps[]` に保存する。

**step 実行・保存固定契約：**

| 項目 | 仕様 |
|------|------|
| step id | 保存時は 0 始まりの `index` と `name` を保存する。`name` が空の場合は API 保存時 `422`。 |
| env merge | runner 基本 env → branch env → `.pipeline_config.env` → pipeline step env の順で上書きする。 |
| secret mask | branch env と step env の secret key 値を stdout/stderr、hook log、notify payload、pipeline step log へ保存前に mask する。 |
| optional failure | `required=false` の step が失敗した場合、`status="optional_failed"` として保存し、後続 step を継続する。全体 status は後続 required step の結果で決める。 |
| timeout | step timeout 時は process group を kill し、`exit_code:null`、`status:"timeout"`、`error:"step timeout"` を保存する。 |
| 保存順 | step 完了ごとにメモリへ結果を追加し、build 終了時に `.build_logs/{id}.json.pipeline_steps[]` へ定義順で保存する。完了順で並べ替えない。 |

**YAML parser 禁止構文固定：**

| 構文 | 処理 |
|------|------|
| tab indent | parse error。 |
| anchor / alias | parse error。 |
| `---` / `...` | parse error。 |
| flow style `{}` / `[]` | parse error。 |
| block scalar `|` / `>` | parse error。 |
| inline comment | quoted string 外の `#` は、行頭 comment 以外 parse error。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| YAML parse error | build 実行前に `failure_pipeline_config`、終了コード `2`。 |
| command 不正 | `failure_pipeline_config`。 |
| timeout | step を kill し、`failure_timeout`。 |
| `.pipeline.yml` 読み取り権限エラー | 終了コード `1`、状態更新なし。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| build step 成功 | step log と build success が記録される。 |
| required step 失敗 | 後続 step を実行せず failure。 |
| optional step 失敗 | WARN、後続 step 継続。 |
| 禁止 YAML 構文 | parse error、build なし。 |
| secret env stdout | pipeline step log では `"***"`。 |
| step timeout | process kill、後続 required step なし。 |

### 27.23 ローカルファイル監視モード

owner component は `runner` とする。collaborator component は `statefile` とする。

本機能の目的は、GitHub API を使わない環境で、ローカル Markdown 入力の変更を SHA-256 snapshot により検出することである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 key | `.server_config.watch_mode` |
| 値 | `"github"` または `"local"`。既定値 `"github"`。 |
| local root | branch target の `src`。絶対 path 必須。 |
| 状態 | `.local_watch_state.json` |
| trigger | local 差分起動時は `"local_watch"`。 |

`.local_watch_state.json` は `{ "files": { "<relative_path>": { "sha256": "...", "mtime_unix": 0, "size": 0 } } }` とする。mode は `600`。

**正常系：**

1. `watch_mode="local"` の場合、GitHub API、PAT、rate limit 処理を呼ばない。
2. `src` 配下の `.md` と `.markdown` を辞書順に列挙する。
3. hidden directory、`.git`、出力先 `out`、`.snapshots` は走査対象外とする。
4. SHA-256 manifest を作成し、前回 `.local_watch_state.json` と比較する。
5. 差分があれば build を実行し、成功時だけ state を更新する。

**local scan 固定契約：**

| 項目 | 仕様 |
|------|------|
| 相対 path | `src` からの相対 path を `/` 区切りで保存する。先頭 `/`、`..`、空 segment は保存しない。 |
| 除外 directory | `.git`、`.snapshots`、`.build_cache`、`.sha_cache`、出力先 `out` 配下、名前が `.` で始まる directory。 |
| 対象拡張子 | `.md`、`.markdown`。大文字拡張子は対象外。 |
| 削除検知 | 前回 state に存在し今回 scan に存在しない path は差分として扱う。 |
| state 更新 | build 成功時に今回 scan 結果へ置換する。skip、failure、dry-run では更新しない。 |
| dry-run | `.local_watch_state.json` を作成・更新せず、差分結果だけ stdout JSON に含める。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| `src` 不在 | `failure_precheck`、終了コード `2`。 |
| state 破損 | full build 扱い、成功時に state 再作成。 |
| ファイル読み取り失敗 | build なし、終了コード `1`。 |
| `watch_mode` 不正 | 起動時設定エラー、終了コード `2`。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 初回 | full build、state 作成。 |
| 変更なし | `skipped_no_change`。 |
| 1 ファイル変更 | build 実行、該当 SHA 更新。 |
| GitHub token 不在 | local mode では失敗しない。 |
| ファイル削除 | build 実行、成功時に state から削除。 |
| out 配下変更 | 差分対象外。 |

### 27.24 タグ付きコミットのみビルド

owner component は `runner` とする。collaborator component は `api`、`statefile` とする。

本機能の目的は、release tag が付いた commit だけを build 対象にする filter を提供することである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 key | `.server_config.tag_filter` |
| schema | `{ "enabled": boolean, "patterns": string[] }` |
| pattern | `*` suffix の prefix match または完全一致のみ。正規表現は禁止。 |
| 既定値 | `{ "enabled": false, "patterns": [] }` |
| log | `.build_logs/{id}.json.matched_tags[]` |

**正常系：**

1. SHA 差分を検出する。
2. `tag_filter.enabled=true` の場合、対象 commit に紐付く tags を GitHub refs API から取得する。
3. `patterns` が空の場合は「任意の tag が 1 件以上」を条件とする。
4. tag が条件に一致した場合だけ build を実行する。
5. 不一致の場合は status `skipped_tag_filter` とし、SHA cache は更新しない。

**tag 判定固定契約：**

| 項目 | 仕様 |
|------|------|
| tag 正規化 | `refs/tags/` prefix を除いた tag 名で pattern 判定する。 |
| 取得順 | GitHub refs API response の順序を維持し、`matched_tags[]` には一致した tag を最大 100 件まで保存する。 |
| pattern `*` | suffix `*` は prefix match。`*` 単体は任意 tag に一致する。中間 `*`、正規表現、glob は禁止。 |
| skip 副作用 | tag 不一致 skip では `.build_status.json` だけ更新し、`.build_logs/{id}.json`、`.build_history`、SHA cache、snapshot、deploy、notify は更新しない。 |
| local mode | `watch_mode="local"` かつ `tag_filter.enabled=true` は設定不整合として終了コード `2`。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| GitHub tags API 失敗 | retry 対象。最終失敗時は build なし、終了コード `3`。 |
| pattern 不正 | API は `422`、runner は終了コード `2`。 |
| tag 数が 1000 超 | 先頭 1000 件だけ評価し、WARN を記録する。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| tag 一致 | build 実行、matched_tags 記録。 |
| tag 不一致 | build skip、SHA cache 未更新。 |
| patterns 空で tag あり | build 実行。 |
| API 失敗 | retry 後 failure、build なし。 |
| `v*` pattern | `v1.0.0` は一致、`release/v1` は不一致。 |
| local mode 併用 | build なし、終了コード `2`。 |

### 27.26 並列マルチターゲットビルド

owner component は `runner` とする。collaborator component は `statefile` とする。

本機能の目的は、複数 deploy target への転送を bounded parallelism で処理し、遅い target が全体を不必要に止めないようにすることである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 key | `.server_config.deploy_parallelism` |
| 許容値 | 1〜16。既定値 1。 |
| 対象 | `branch_targets[].deploy_targets[]` |
| build log | `target_results[]` に target id、status、started_at、finished_at、error を保存。 |

**正常系：**

1. build 成功後、deploy target を設定順に queue へ入れる。
2. worker 数は `min(deploy_parallelism, len(targets))` とする。
3. target ごとに SSH 転送と remote checksum 検証を行う。
4. target 成功/失敗を個別に記録する。
5. 1 target 以上失敗した場合、全体 status は `success_deploy_pending` とし、失敗 target だけ `.pending_transfers` に追加する。

**parallel deploy 保存固定契約：**

| 項目 | 仕様 |
|------|------|
| result 順序 | `target_results[]` は設定順で保存する。完了順では保存しない。 |
| worker 上限 | `deploy_parallelism` が target 数を超える場合も worker 数は target 数まで。 |
| pending 重複 | 同一 build id、target id、dest path の pending が既にある場合は重複追加しない。 |
| 成功 target | 失敗 target があっても成功 target は pending に入れない。 |
| status | 1 件以上 pending があれば `success_deploy_pending`、全件成功なら `success`。 |
| 保存順 | `.build_logs/{id}.json.target_results` → `.pending_transfers` → `.build_history` → `.build_status.json`。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| parallelism 不正 | API は `422`、runner は終了コード `2`。 |
| worker panic 相当の内部エラー | 該当 target failure、他 target は継続。 |
| 全 target 失敗 | status `success_deploy_pending`、pending に全件追加。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 3 target / parallelism 2 | 同時実行最大 2、全 target result 記録。 |
| 1 target 失敗 | pending は 1 件、成功 target は再投入しない。 |
| parallelism 1 | 既存順次処理と同じ結果。 |
| result order | 完了順に関係なく設定順で保存。 |
| pending duplicate | 同一 pending は 1 件。 |

### 27.27 ビルド前後フック

本機能の目的は、build 前後に登録済み command を安全に実行し、外部 shell 文字列に依存しない拡張点を提供することである。

owner component は `runner` とする。collaborator component は `api`、`statefile` とする。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 状態 | `.hooks` JSON object。mode `600`。 |
| API | `GET /api/hooks`、`POST /api/hooks`、`DELETE /api/hooks/{id}`、`GET /api/hooks/{id}/log` |
| phase | `"pre"` または `"post"`。 |
| command_args | string 配列。shell 経由禁止。 |
| timeout_seconds | 1〜3600。省略時 300。 |
| log | `.build_logs/{build_id}_hook_{hook_id}.json` |

**正常系：**

1. pre hook を id 昇順に実行する。
2. pre hook が成功した場合だけ build 本体へ進む。
3. build 終了後、post hook を id 昇順に実行する。
4. hook ごとに stdout/stderr、exit_code、duration_seconds を保存する。
5. post hook 失敗は build status を変更しない。

**hook 保存・mask 固定契約：**

| 項目 | 仕様 |
|------|------|
| hook id | `^[A-Za-z0-9_-]{1,64}$`。重複 id は API 保存時 `422`。 |
| 実行順 | phase ごとに id 昇順。pre 全件後に build、build 後に post。 |
| env | branch env と hook 固有 env を渡す。secret key の値は hook log 保存前に mask する。 |
| hook log | 1 実行 1 JSON object とし、`hook_id`、`build_id`、`phase`、`status`、`started_at`、`finished_at`、`duration_seconds`、`stdout`、`stderr`、`exit_code`、`timed_out`、`truncated` を保存する。 |
| log 保存失敗 | pre hook の log 保存失敗は build を開始せず failure。post hook の log 保存失敗は build status を維持し runner 終了コードを最低 `1`。 |
| shell 禁止 | `command_args` を `exec.Command` 相当で実行し、shell 展開、変数展開、glob 展開を行わない。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| pre hook 失敗かつ `abort_on_failure=true` | build 本体を実行せず status `hook_error`。 |
| pre hook 失敗かつ `abort_on_failure=false` | WARN、build 継続。 |
| command 不正 | API は `422`、runner は該当 hook failure。 |
| timeout | process kill、exit_code `null`、status `timeout`。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| pre success | build 本体が実行される。 |
| pre abort | build なし、history `hook_error`。 |
| post failure | build 結果維持、hook log 記録。 |
| shell metachar | command_args として渡され、shell 展開されない。 |
| hook log write failure | pre は build なし、post は build 結果維持。 |
| secret stdout | hook log では `"***"`。 |

### 27.29 リモートビルド対応

owner component は `runner` とする。collaborator component は `api`、`archive`、`statefile` とする。

本機能の目的は、runner が SSH 先で build を実行し、成果物を archive と manifest で回収できるようにすることである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 key | `.server_config.remote_build` |
| schema | `{ "enabled": boolean, "host": string, "user": string, "work_dir": string, "command_args": string[], "artifact_path": string }` |
| command_args | shell 経由禁止。SSH 先で実行する argv。 |
| artifact | tar.gz。必須ファイル `site/`、`manifest.json`。 |
| log | `.build_logs/{id}.json.remote_build` |

**正常系：**

1. remote build enabled の場合、local builder を起動しない。
2. SSH で remote work dir を確認する。
3. `command_args` を remote で実行する。
4. `artifact_path` を取得し、一時 directory へ展開する。
5. `manifest.json` の SHA-256 と展開 file を検証する。
6. 検証成功後、既存 deploy 処理へ渡す。

**remote artifact 固定契約：**

| 項目 | 仕様 |
|------|------|
| 一時展開先 | state dir 配下 `.remote_artifacts/{build_id}/`。既存 output directory へ直接展開しない。 |
| tar.gz entry | `site/` と `manifest.json` だけを root 直下必須とする。entry path の `..`、絶対 path、NUL、symlink、device は拒否する。 |
| manifest schema | `{ "files": [{"path": string, "sha256": string, "size": integer}] }`。 |
| 検証順 | tar.gz 展開前 entry 検査 → 一時展開 → manifest parse → file 存在 / size / sha256 検証 → deploy へ渡す。 |
| secret | remote command stdout/stderr は保存前に mask する。SSH 秘密鍵 path や token 値は log に保存しない。 |
| cleanup | 成功・失敗に関係なく、検証後に一時展開先の削除を試みる。削除失敗は WARN。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| SSH 接続失敗 | retry 対象。最終失敗で status `failure_remote_build`。 |
| artifact 不在 | failure、deploy しない。 |
| manifest 不一致 | failure、deploy しない。 |
| remote command timeout | process kill、failure_timeout。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| remote 成功 | artifact 検証後 deploy。 |
| manifest 不一致 | deploy なし、failure。 |
| SSH 一時失敗 | retry 後成功なら success。 |
| unsafe tar entry | 展開中止、deploy なし。 |
| cleanup 失敗 | build 結果維持、WARN。 |

### 27.30 ビルド承認フロー

owner component は `runner` とする。collaborator component は `api`、`statefile`、`sdk`、`ui` とする。

本機能の目的は、`approval_required` な target を通常 build として即時実行せず、人間承認後の queue entry だけを build / deploy 実行対象にすることである。

API endpoint、approve / reject の request / response、sdk / ui 操作境界は `ADLAIRE_CI_DETAIL_API_SPEC.md` §27.30 を正とする。`.approval_queue` record schema は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0c を正とする。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 key | `branch_targets[].approval_required` |
| timeout | `.server_config.approval_timeout_seconds`。既定値 86400。 |
| 承認状態 | `.approval_queue` JSON Lines。 |
| 実行 queue | `.build_state.queued[]` の `trigger:"approval"` entry。 |
| 通知設定 | `.notify_config` の `approval_required` event。 |

**正常系：**

1. runner は target 判定時、`approval_required=true` かつ通常 build 条件成立の場合、pipeline / deploy / snapshot を開始しない。
2. runner は `.approval_queue` lock を取得し、同一 branch / sha / target の最新 `pending` record を確認する。
3. 重複 pending がない場合、runner は `pending` record を `.approval_queue` へ追記する。
4. runner は `.notify_config` に従い approval request 通知を送信する。
5. 通知成功時は `.notify_log` に結果を記録し、通知失敗時は `.notify_pending` に retry entry を追記する。
6. runner は起動時に期限超過 pending を検出し、`expired` record と `.build_history.status="approval_expired"` を追記する。
7. runner は approve API 由来の queue entry を通常 queue 処理として取り出し、`trigger="approval"` で build / deploy を実行する。

**pending 作成固定契約：**

| 項目 | 仕様 |
|------|------|
| pending id | `appr{YYYYMMDDHHmmss}`、同秒衝突時は `-001` から連番。既存 id は再利用しない。 |
| 重複 pending | 同一 branch / sha / target の最新 status が `pending` の record。 |
| 重複時 | 新規 record を作成せず、既存 pending id を使用する。新規通知は送らない。 |
| 通知 payload | `{event:"approval_required", approval_id, branch, sha, target, expires_at}`。secret、token、path secret は含めない。 |
| build 抑止 | pending 作成成功または重複 pending 検出時、当該 target の build は開始しない。 |

**approval runner 状態遷移：**

| 現在 status | runner 操作 | 次 status | 副作用 |
|-------------|-------------|-----------|--------|
| なし | approval_required 検出 | `pending` | `.approval_queue` に pending record を追記し、通知を試行する。 |
| `pending` | 重複検出 | `pending` | 新規 record / 通知 / build を発生させない。 |
| `pending` | timeout | `expired` | `.build_history` に `status:"approval_expired"` を追記する。queue は追加しない。 |
| `approved` | queue 取り出し | 変更なし | `trigger:"approval"` として build / deploy を実行する。 |
| `rejected` | runner 起動 | 変更なし | build しない。 |
| `expired` | runner 起動 | 変更なし | build しない。 |

**approval runner 更新順：**

| 操作 | 更新順 | 失敗時 |
|------|--------|--------|
| pending 作成 | `.approval_queue` lock → 重複確認 → pending record append → 通知送信 → `.notify_log` / `.notify_pending` 更新 | 通知失敗でも pending は残す。pending append 失敗時は build を開始せず runner failure。 |
| timeout | runner 起動時に `.approval_queue` lock → expires_at 超過 pending を created_at 昇順で抽出 → expired record append → `.build_history` append | history append 失敗時も expired record は残し、runner は ERROR を出して継続する。 |
| approved queue 実行 | `.build_state` lock → `trigger:"approval"` entry を 1 件取り出し → build / deploy 実行 → history / status finalizer | queue entry が不正な場合は `failure_decode` として記録し、次 entry は次回起動まで処理しない。 |

runner は `approval_required=true` の target に対して、approval queue 以外の経路で build を開始してはならない。manual force、webhook、force interval、local watch のいずれであっても、対象 target が approval_required の場合は pending 作成を優先し、承認済み queue entry になるまで pipeline を起動しない。

**異常系：**

| 条件 | 処理 |
|------|------|
| `.approval_queue` lock 取得失敗 | build を開始せず終了コード `1`。 |
| pending append 失敗 | build を開始せず終了コード `1`。 |
| 通知失敗 | pending は維持し、`.notify_pending` に追記する。 |
| `.notify_pending` 追記失敗 | pending は維持し、ERROR ログを出して runner 終了コードを最低 `1` にする。 |
| timeout history append 失敗 | expired record は維持し、ERROR ログを出して継続する。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| approval required | build せず pending 作成。 |
| duplicate pending | 新規 record / 通知なし、build なし。 |
| notify failure | pending 維持、`.notify_pending` 追記。 |
| timeout | expired、history `approval_expired`、build なし。 |
| approved queue | `trigger="approval"` で build / deploy 実行。 |

### 27.31 ブランチ別環境変数

owner component は `runner` とする。collaborator component は `api`、`statefile` とする。

本機能の目的は、branch target ごとに build process へ注入する環境変数を定義し、branch や deploy 先ごとの差分を、保存前検証、注入対象固定、secret mask、log 保存禁止値によって扱うことである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 key | `branch_targets[].env` |
| 型 | string:string object |
| key | `^[A-Z_][A-Z0-9_]{0,63}$` |
| value | UTF-8 string。最大 4096 bytes。NUL 禁止。 |
| secret key | key に `TOKEN`、`SECRET`、`PASSWORD`、`PAT` を含むもの。 |
| 上限 | branch target ごとに 100 key。 |

**正常系：**

1. API は env key/value を検証して `.branch_config` に保存する。
2. runner は build process の environment に branch env を追加する。
3. 同名 key が system env に存在する場合、branch env を優先する。
4. build log には env key 一覧だけを保存し、value は保存しない。
5. secret key は stdout/stderr の mask 対象に追加する。

**env 正規化・mask 固定契約：**

| 項目 | 仕様 |
|------|------|
| 保存順 | env key は ASCII 昇順で保存する。 |
| value 正規化 | UTF-8 不正、NUL、改行を含む値は禁止。前後空白は保持する。 |
| secret 判定 | key に `TOKEN`、`SECRET`、`PASSWORD`、`PAT` を含む場合は大文字小文字を区別せず secret。 |
| log 保存 | `.build_logs/{id}.json.environment.env_keys` に key 名だけを保存する。value、value length、hash は保存しない。 |
| process env | branch env は builder、pipeline step、hook、command notification に渡す。通知 payload には値を含めない。 |
| mask failure | mask 対象値を保存前に置換できない場合、build を失敗扱いにし、平文を保存しない。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| key 不正 | API は `422`、runner は終了コード `2`。 |
| value 上限超過 | `422`。 |
| secret mask 漏れ | 実装不合格。該当 build は完了扱いにしない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| env 注入 | build command が指定値を参照できる。 |
| system env と同名 | branch env が優先される。 |
| secret stdout 出力 | log では `"***"` に置換。 |
| 不正 key | 保存不可、状態差分なし。 |
| lower secret key | `my_token` も secret 扱い。 |
| mask failure | build 完了扱いにしない。 |

### 27.32 ビルド通知連携

本機能の目的は、build lifecycle event を複数通知 channel へ同一契約で送信し、通知の成功、失敗、再試行、監査を固定仕様で扱えるようにすることである。

owner component は `runner` とする。collaborator component は `api`、`sdk`、`ui`、`statefile` とする。runner は送信、api は設定・履歴表示、sdk / ui は設定操作と履歴表示を担当する。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 | `.notify_config.channels[]` |
| channel type | `"webhook"`、`"email"`、`"command"` |
| event | `"start"`、`"success"`、`"failure"`、`"deploy_failure"`、`"weekly_summary"`、`"approval_required"`、`"duration_anomaly"`、`"config_corrupt"` |
| log | `.notify_log` JSON Lines |
| pending | `.notify_pending` JSON array |
| secret | webhook secret、SMTP password、command env secret は GET response と log で必ず mask。 |

`command` channel は `command_args` 配列だけを許可し、shell 文字列は禁止する。`command_args[0]` は絶対 path または PATH 解決可能なコマンド名とする。

**正常系：**

1. build event 発生時に `.notify_config` を読み込む。
2. `enabled=true` かつ event が `on[]` に含まれる channel を抽出する。
3. channel id 昇順で送信する。
4. 各送信結果を `.notify_log` へ JSON Lines で追記する。
5. retry 対象失敗は `.notify_pending` に追加する。
6. build 自体の status は通知失敗で変更しない。

**通知 payload / id 固定契約：**

| 項目 | 仕様 |
|------|------|
| notify id | `.notify_log` は `ntfy{YYYYMMDDHHmmss}`、pending は `np{YYYYMMDDHHmmss}`、衝突時 `-001`。 |
| payload 共通 key | `event`、`build_id`、`status`、`branch`、`trigger`、`created_at` を含める。取得できない値は `null`。 |
| channel 順 | channel id 昇順。送信失敗しても次 channel を継続する。 |
| timeout | webhook と command は 30 秒。SMTP は 60 秒。timeout は retry 対象か個別表に従う。 |
| pending payload | mask 後 payload だけを保存し、secret は retry 送信直前に状態ファイルから再読込する。 |
| pending 保存順 | `.notify_log` 追記成功後に `.notify_pending` を保存する。log 失敗時は pending を追加しない。 |

**通知 channel 正規化契約：**

| 項目 | 仕様 |
|------|------|
| channel id | 未指定時は `n` + 6 桁連番。重複 id は `422`。 |
| event 判定 | channel `on` が空の場合は top-level `on` を使用する。channel `on` に `"*"` があれば全 event 対象。 |
| webhook config | `config.url` 必須。`http` / `https` のみ許可。userinfo、fragment、空 host は `422`。 |
| email config | `.smtp_config` / `.smtp_secret` を正とする。channel `config.to` がある場合は `.smtp_config.to` より優先して送信先に使う。 |
| command config | `config.command_args` 必須。shell 経由は禁止。stdout/stderr は `.notify_log` に secret mask 後で保存する。 |
| secret mask | `secret`、`password`、`token`、`smtp_password`、Webhook secret、SMTP password は GET、backup、log、pending、UI 表示で `"***"`。 |

**通知送信・pending 固定契約：**

| ケース | `.notify_log` | `.notify_pending` | build status |
|--------|---------------|-------------------|--------------|
| 送信成功 | `result:"success"` を追記 | 変更なし | 変更しない |
| webhook 5xx / timeout | `result:"failure"` を追記 | retry entry 追加 | 変更しない |
| webhook 4xx | `result:"failure"` を追記 | 追加しない | 変更しない |
| email SMTP 未設定 | `result:"not_configured"` を追記 | 追加しない | 変更しない |
| command exit 非 0 | `result:"failure"` を追記 | 追加しない | 変更しない |
| `.notify_log` 追記失敗 | server log に固定コード | pending 追加判定は実行しない | 変更しない |
| `.notify_pending` 保存失敗 | `result:"failure"` を追記済み | 追加なし | 変更しない |

pending entry は `{ "id", "event", "channel_id", "channel_type", "payload", "attempts", "next_attempt_at", "last_error", "created_at" }` を必須 key とする。`payload` は secret mask 済み JSON object とし、送信時に secret を復元しない。Webhook secret や SMTP password は送信直前に状態ファイルから再読込する。`event`、`channel_id`、mask 後 `payload` が同一の pending entry が既にある場合は重複追加しない。

runner 起動時の pending retry は `next_attempt_at <= now` の entry を `created_at` 昇順で処理する。成功した entry は削除する。失敗した entry は `attempts += 1`、`next_attempt_at = now + retry_interval_seconds` として保存する。`attempts > retry_count` になった entry は `.notify_log` に `result:"dropped"` を追記して pending から削除する。

**異常系：**

| 条件 | 処理 |
|------|------|
| `.notify_config` 破損 | 起動時整合性チェックに従い復旧し、通知送信は skip。 |
| webhook HTTP 5xx / timeout | retry 対象として `.notify_pending` へ追加。 |
| webhook HTTP 4xx | retry しない。`.notify_log` に failure。 |
| SMTP 未設定 | email channel は `not_configured` として log、retry しない。 |
| command timeout | process kill、retry 対象外、failure log。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| success event | 対象 channel へ送信、notify log 追記。 |
| webhook 5xx | pending 追加。 |
| secret 設定済み | GET / log / UI で値が `"***"`。 |
| command channel | shell 展開されず argv 実行。 |
| pending duplicate | 同一 event/channel/payload の pending は 1 件だけ。 |
| retry exhausted | dropped log 追記後 pending から削除。 |

### 27.33 ビルド時間トレンド記録

owner component は `runner` とする。collaborator component は `api`、`statefile` とする。

本機能の目的は、build 所要時間の統計を蓄積し、性能傾向と回帰検知の基準を提供することである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 状態 | `.build_trends.json` |
| API | `GET /api/stats/build-trends?n=N` |
| sample | `{build_id, finished_at, branch, trigger, duration_seconds, status, anomaly}` |
| 保持件数 | `.server_config.build_trend_keep_count`。既定値 1000、許容値 10〜10000。 |

**正常系：**

1. build 完了時、`duration_seconds != null` の場合だけ sample を追加する。
2. `finished_at` 昇順で保存し、保持件数超過分は古い順に削除する。
3. summary に `count`、`avg_seconds`、`median_seconds`、`p95_seconds`、`anomaly_count` を保存する。
4. API は `n` の最新 sample と summary を返す。`n` は 1〜1000、既定値 100。

**統計算出固定契約：**

| 項目 | 仕様 |
|------|------|
| 対象 sample | `duration_seconds` が 0 以上の数値で、`status` が `"success"`、`"failure"`、`"success_deploy_pending"` のいずれかである sample。 |
| `avg_seconds` | 対象 duration の合計を対象件数で割り、小数第 3 位を四捨五入して小数第 2 位まで保存する。対象 0 件の場合は `null`。 |
| `median_seconds` | duration 昇順の中央値。偶数件の場合は中央 2 件の平均を使い、小数第 3 位を四捨五入して小数第 2 位まで保存する。 |
| `p95_seconds` | duration 昇順配列の `ceil(count * 0.95) - 1` 番目を採用する。`count=0` の場合は `null`。 |
| `anomaly_count` | `sample.anomaly == true` の件数。 |
| 同一 build id | 既存 sample と同じ `build_id` を追加する場合は append せず既存 sample を置換し、`finished_at` 昇順に再整列する。 |
| API response | `{ "samples": TrendSample[], "summary": TrendSummary, "warnings": string[] }` を返す。warnings がない場合は空配列。 |

`.build_trends.json` は sample 置換または追加後に summary を再計算して atomic write する。summary だけの部分更新は禁止する。

**異常系：**

| 条件 | 処理 |
|------|------|
| `.build_trends.json` 破損 | `.build_trends.json.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避後、`.build_history` の有効行から再集計する。 |
| 再集計不能 | 初期値で作成し、WARN を出す。 |
| `n` 不正 | API は `422`。 |
| `.build_history` に duration 欠落 | 当該行は再集計対象外とし、warnings に `trend_sample_skipped` を 1 回だけ含める。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| build success | sample 追加、summary 更新。 |
| build failure | duration があれば sample 追加。 |
| 保持件数超過 | 古い sample だけ削除。 |
| API n=10 | 最新 10 件を返す。 |
| 同一 build id 再記録 | sample は 1 件のまま値が置換される。 |
| p95 算出 | `ceil(count * 0.95) - 1` の値と一致する。 |

### 27.34 ビルド依存チェーン

owner component は `runner` とする。collaborator component は `api`、`statefile` とする。

本機能の目的は、複数 build job の依存関係を DAG として定義し、依存 job 成功後だけ後続 job を実行することである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 状態 | `.build_chain_config` |
| API | `GET /api/build-chain-config` / `POST /api/build-chain-config` |
| schema | `{ "chains": ChainJob[] }` |
| job key | `id`, `branch`, `target_file`, `depends_on`, `required`, `enabled` |
| 上限 | job 100 件、依存 20 件 / job。 |

`id` は `^[a-zA-Z0-9_-]{1,64}$` とする。`depends_on` は同一 config 内の job id のみ許可する。循環依存は禁止する。

**正常系：**

1. API 保存時に schema、参照整合、循環を検証する。
2. runner は chain enabled job を topological order で実行する。
3. 依存 job が失敗した場合、後続 required job は `skipped_dependency_failed` として history に記録する。
4. `required=false` の依存失敗は WARN とし、後続 job を継続できる。
5. `.build_logs/{id}.json.chain` に job id、depends_on、chain_index を保存する。

**chain 実行固定契約：**

| 項目 | 仕様 |
|------|------|
| `chain_run_id` | chain 起動ごとに `chain{YYYYMMDDHHmmss}`、衝突時 `-001` を付与する。chain 内の全 job log / history で同じ値を使う。 |
| topological 同順位 | 依存数が同じで同時に実行可能な job は、`.build_chain_config.chains[]` の出現順で実行する。 |
| disabled job | `enabled=false` の job は DAG から除外する。enabled job が disabled job に依存する場合、API 保存時に `422`。 |
| job build id | 各 job は独立した build id を持つ。build id は通常採番規則に従い、chain job であることを id 文字列へ埋め込まない。 |
| skip log | dependency failure により skip した job も `.build_history` に 1 行追記し、`.build_logs/{id}.json` は作成しない。 |
| chain summary | 最終 job 処理後、`.build_logs/{last_id}.json.chain_summary` に `chain_run_id`、`total_jobs`、`success_count`、`failure_count`、`skipped_count` を保存する。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| 循環依存 | API は `422`、runner は chain 無効化して通常 build。 |
| 依存 job 不在 | `422`。 |
| job 実行中に runner 停止 | 完了済み job だけ history に残し、未実行 job は次回再判定。 |
| chain summary 保存失敗 | job 結果は維持し、ERROR ログを出して runner 終了コードを最低 `1` にする。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| A -> B | A 成功後 B 実行。 |
| A 失敗 / B required | B は `skipped_dependency_failed`。 |
| 循環 | 保存不可。 |
| optional 依存失敗 | 後続 job 継続。 |
| 同順位 job | config 出現順で実行される。 |
| disabled 依存 | 保存時 `422`、状態差分なし。 |

### 27.35 ビルド優先度キュー

owner component は `runner` とする。collaborator component は `api`、`statefile` とする。

本機能の目的は、manual、webhook、approval などの queue entry を優先度順に処理し、緊急 build を先に実行できるようにすることである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 状態 | `.build_state.queued[]` |
| priority | `"low"`、`"normal"`、`"high"`、`"urgent"` |
| 数値順 | urgent=0、high=10、normal=20、low=30 |
| created_seq | queue 追加時に単調増加する整数。 |
| 既定値 | priority `"normal"` |

**正常系：**

1. queue 追加時に priority と created_seq を保存する。
2. runner は priority 数値昇順、同一 priority では created_seq 昇順で 1 件だけ処理する。
3. `GET /api/queue` は並び替え後の queue を返す。
4. `DELETE /api/queue` は waiting entry 全件を削除し、実行中 build は停止しない。

**priority / created_seq 固定契約：**

| 項目 | 仕様 |
|------|------|
| 採番元 | `.build_state.queued[].created_seq` の最大値。存在しない場合は `0`。次 entry は最大 + 1。 |
| 旧 entry | `created_seq` 欠落 entry は GET 表示時だけ末尾扱いとし、runner 取り出し前に `.build_state` lock 内で `created_seq` を正規化保存する。正規化値は採番元の規則だけで決定し、entry 内容、priority、id、時刻から推測しない。 |
| priority 省略 | API / webhook / approval の queue 追加時は `"normal"` を保存する。 |
| 表示順 | priority 数値昇順、created_seq 昇順、同値なら id 昇順。 |
| 取り出し順 | 表示順と同一。 |
| clear | priority に関係なく waiting entry 全件を削除する。 |
| queue full | priority が高くても既存 entry を押し出さない。 |

runner が旧 entry の `created_seq` 正規化保存に失敗した場合、build を開始せず終了コード `1` とする。正規化前の推測順で build を開始してはならない。

**異常系：**

| 条件 | 処理 |
|------|------|
| priority 不正 | API は `422`。 |
| created_seq 欠落の旧 entry | GET 表示時は末尾扱いにする。runner 取り出し前または queue 更新時に `.build_state` lock 内で正規化保存する。正規化保存失敗時は build を開始しない。 |
| queue full | `429`。priority による上書き削除はしない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| urgent 後投入 | normal より先に処理。 |
| 同一 priority | FIFO。 |
| 不正 priority | 状態差分なしで `422`。 |
| created_seq 欠落 | 正規化保存後に順序判定し、正規化保存失敗なら build なし。 |
| urgent queue full | `429`、既存 low entry も削除しない。 |

### 27.36 失敗原因の自動分類

owner component は `runner` とする。collaborator component は `api`、`statefile` とする。

本機能の目的は、build failure を固定カテゴリへ分類し、調査開始点を build log、history、UI に残すことである。

**分類値：**

| category | 判定条件 |
|----------|----------|
| `github_api` | GitHub API 最終失敗、rate limit 復帰不可、認証失敗。 |
| `pipeline_timeout` | build step timeout。 |
| `pipeline_exit` | build command exit code 非 0。 |
| `deploy_failure` | SSH 転送、remote checksum、pending transfer 発生。 |
| `hook_error` | pre hook abort または hook timeout。 |
| `config_error` | 設定 parse、validation、必須値不足。 |
| `resource_error` | disk 不足、binary 不在、権限エラー。 |
| `unknown` | 上記に該当しない failure。 |

**正常系：**

1. failure 確定時に分類優先順位で category を決定する。
2. `.build_logs/{id}.json.failure_category` と `failure_evidence[]` を保存する。
3. `.build_history.failure_category` に同じ値を保存する。
4. api / ui は category で filter できる。未知 query は `422`。

**`failure_evidence[]` schema：**

| key | 型 | 必須 | 仕様 |
|-----|----|------|------|
| `source` | string | 必須 | `"github_api"`、`"pipeline"`、`"deploy"`、`"hook"`、`"config"`、`"resource"`、`"runner"`。 |
| `code` | string | 必須 | 固定コード。例: `http_500`、`timeout`、`exit_nonzero`、`checksum_mismatch`。 |
| `message` | string | 必須 | 固定文言。入力値、secret、token、path 全体を連結しない。最大 300 文字。 |
| `at` | string | 必須 | UTC ISO 8601。 |

分類時は `failure_evidence[]` を最大 10 件まで保存する。10 件を超える場合は分類に使った evidence を先頭に残し、残りは発生順で 9 件まで保存する。

**API filter 固定契約：**

| 項目 | 仕様 |
|------|------|
| query | `failure_category` を受け付ける。空文字は未指定扱い。 |
| 未知 category | `422 {"error":"Invalid failure category"}`。 |
| 成功 history | `failure_category:null` として返す。 |
| 既存未知値 | response に含め、top-level `warnings:["unknown_failure_category"]` を 1 回だけ返す。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| 複数カテゴリ該当 | 表の上から最初に該当した category を採用する。 |
| evidence 抽出不能 | category は保存し、evidence は空配列。 |
| 既存 history に未知 category | API は返すが warnings に `unknown_failure_category`。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| timeout | `pipeline_timeout`。 |
| SSH 失敗 | `deploy_failure`。 |
| 設定不正 | `config_error`。 |
| filter | category 指定で該当履歴だけ返る。 |
| evidence 上限超過 | 最大 10 件で保存され、secret 平文を含まない。 |
| 成功履歴 | `failure_category:null`。 |

### 27.37 ビルド実行環境の記録

owner component は `runner` とする。collaborator component は `statefile` とする。

本機能の目的は、build 時点の実行環境を記録し、後から再現性と障害原因を確認できるようにすることである。

**記録先：**

`.build_logs/{id}.json.environment` に以下を保存する。

| key | 型 | 説明 |
|-----|----|------|
| `os` | string | `runtime.GOOS`。 |
| `arch` | string | `runtime.GOARCH`。 |
| `go_version` | string | `runtime.Version()`。 |
| `runner_version` | string | build info または `"unknown"`。 |
| `builder_version` | string | `adlaire-ci-build --version` の先頭 token。 |
| `hostname` | string | OS hostname。 |
| `state_dir` | string | `--state-dir`。 |
| `disk_free_bytes` | integer/null | state dir filesystem の空き容量。 |
| `captured_at` | string | UTC ISO 8601。 |

**正常系：**

1. build id 採番直後に environment snapshot を取得する。
2. builder 起動前に `.build_logs/{id}.json.environment` へ保存する。
3. 取得不能項目は `null` または `"unknown"` とし、build は継続する。
4. 環境変数の値、token、secret、PATH 全体は保存しない。

**取得・保存固定契約：**

| 項目 | 仕様 |
|------|------|
| `builder_version` | `adlaire-ci-build --version` を最大 2 秒で実行し、stdout 先頭行の最初の空白区切り token を保存する。stderr は保存しない。 |
| `runner_version` | Go build info の main version が空の場合は `"unknown"`。VCS revision は保存しない。 |
| `hostname` | 255 文字を超える場合は 255 文字で切り詰める。取得失敗時は `"unknown"`。 |
| `state_dir` | `--state-dir` が home directory 配下の場合は basename だけ保存する。それ以外は絶対 path を保存する。 |
| 保存失敗 | environment 保存失敗は build を開始せず、`.build_status.json` に `failure_state_write` を保存し、終了コード `1`。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| hostname 取得失敗 | `"unknown"`。 |
| disk stat 失敗 | `disk_free_bytes=null`、WARN。 |
| builder version 取得 timeout | `"unknown"`、build 継続。 |
| environment 保存失敗 | build 本体を実行せず failure。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 通常 build | environment object が保存される。 |
| builder version 失敗 | build 継続、unknown。 |
| secret env 存在 | log に値が出ない。 |
| home 配下 state dir | basename だけ保存される。 |
| environment write failure | pipeline を起動しない。 |

### 27.38 ビルド所要時間の異常検知

owner component は `runner` とする。collaborator component は `api`、`statefile` とする。

本機能の目的は、過去 trend と比較して異常に遅い build を検出し、性能劣化を WARN、history flag、通知で可視化することである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 key | `.server_config.duration_anomaly` |
| schema | `{ "enabled": boolean, "min_samples": integer, "avg_multiplier": number, "p95_multiplier": number }` |
| 既定値 | `{ "enabled": false, "min_samples": 20, "avg_multiplier": 2.0, "p95_multiplier": 1.5 }` |
| 参照状態 | `.build_trends.json` |
| 通知 event | `duration_anomaly` |

**正常系：**

1. build 完了後、今回 duration を trend 更新前の summary と比較する。
2. sample 数が `min_samples` 未満の場合は判定しない。
3. `duration > avg_seconds * avg_multiplier` または `duration > p95_seconds * p95_multiplier` なら anomaly とする。
4. anomaly の場合、WARN log、`.build_history.flagged=true`、`.build_history.tags += ["duration_anomaly"]` を保存する。
5. `.notify_config` に `duration_anomaly` event 対象 channel がある場合は通知する。
6. 最後に `.build_trends.json` へ今回 sample を追加する。

**判定・通知固定契約：**

| 項目 | 仕様 |
|------|------|
| 判定対象 | build status が `"success"` または `"success_deploy_pending"` で、`duration_seconds` が 0 以上の場合だけ判定する。failure build は trend sample には含めるが anomaly 判定しない。 |
| 比較基準 | 今回 sample を追加する前の `.build_trends.json.summary` を使う。復旧再集計が発生した場合も復旧後、追加前の summary を使う。 |
| tag 更新 | 既存 `tags` に `"duration_anomaly"` がある場合は追加しない。 |
| notify payload | `{event:"duration_anomaly", build_id, branch, duration_seconds, avg_seconds, p95_seconds, threshold_source}`。`threshold_source` は `"avg"`、`"p95"`、`"avg_and_p95"`。 |
| 設定 validation | `min_samples` は 1〜10000、`avg_multiplier` と `p95_multiplier` は 1.0〜100.0。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| trend 破損 | §27.33 の復旧後、復旧できた場合だけ判定する。 |
| avg / p95 が null | 判定しない。 |
| 通知失敗 | build status は変更せず notify retry 契約に従う。 |
| 設定値不正 | API は `422`、runner は既定値ではなく機能無効として扱う。 |
| trend 保存失敗 | anomaly 判定結果と history は維持し、runner 終了コードを最低 `1` にする。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| sample 不足 | anomaly 判定なし。 |
| 平均 2 倍超 | WARN、flag、tag、通知 event。 |
| p95 以内 | anomaly なし。 |
| 通知失敗 | build success 維持、pending 追加。 |
| failure build | trend sample 追加、anomaly 判定なし。 |
| tag 重複 | `duration_anomaly` が 1 件だけ。 |
