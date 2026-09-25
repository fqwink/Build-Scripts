# Adlaire CI — Runner 詳細仕様

[`docs/details/runner.md`](runner.md) は `runner` owner component の詳細本文責務として、`runner` が主本文として持つ実装契約だけを扱う。

owner / collaborator 境界管理は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0b.1](../DETAIL_INDEX.md#0b1-owner-component-別-owner-collaborator-境界管理) に従う。`runner` owner component の主本文であり、collaborator component の仕様は呼び出し境界、schema、setup、security、検証観点として参照する。fixture、expected、fake、実装検証証跡は [`docs/details/fixture.md`](fixture.md) fixture 証跡責務を参照する。

runner 拡張機能の owner / collaborator は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.2](../DETAIL_INDEX.md#0i2-runner--ci-実行) を入口とする。runner の処理順と副作用は [`docs/details/runner.md`](runner.md)、状態 schema は [`docs/details/statefile.md`](statefile.md)、archive 実体は [`docs/details/archive.md`](archive.md)、Commit Status は [`docs/details/commitstatus.md`](commitstatus.md)、API / SDK / UI 接続は各 owner 詳細本文、検証証跡は [`docs/details/fixture.md`](fixture.md) を正本とする。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `runner` |
| 実装主体 | [`components/runner.go`](../../components/runner.go)。起動入口は [`main.go`](../../main.go)、実行バイナリ名は `adlaire-ci-runner` とする。 |
| 持つ内容 | `runner` owner が主本文として定義する GitHub 監視、設定読取、状態ファイル更新呼び出し、pipeline、deploy、snapshot 作成トリガー、通知、runner 検証条件、runner owner 追加機能。 |
| 持たない内容 | API endpoint の認証・応答本文、SDK method 実装、UI DOM 詳細、builder の変換処理、admin 静的配信、security 主本文、状態 schema、setup / update 手順、release 生成・公開手順、fixture 証跡責務。 |

---

## 10. CI ランナー 要件

| 項目 | 内容 |
|------|------|
| Go ランタイム | 最小バージョンは [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の共通固定値を参照する。 |
| 対象 OS | Linux（systemd 対応環境） |
| ネットワーク | GitHub API は `api.github.com` への HTTPS、通知は検証済みの設定先への HTTP(S) または SMTP、deploy と remote build は設定済み host への SSH だけを許可する。接続先、port、timeout、認証は各機能契約に定義された値だけを許可し、未定義の送信先へ接続しない。 |

---

## 10a. CI ランナー 実装対象

[`docs/details/runner.md` 詳細本文責務 §10a](runner.md#10a-ci-ランナー-実装対象) は、`runner` として実装する CI ランナー機能を定義する。

`runner` は `adlaire-ci-runner` バイナリとして実行する。起動形式は systemd timer から呼び出される oneshot 実行とし、1 回の起動で対象ブランチ設定を読み込み、変更検出、ビルド起動、ログ保存、通知、転送、後処理を完了して終了する。

実装時は、対象項目ごとに [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0c](../DETAIL_INDEX.md#0c-実装前確認項目) の実装前確認項目を満たしていることを確認する。未充足の項目が 1 つでもある場合の実装着手可否は [`docs/SPEC.md` ポリシー責務 §0a](../SPEC.md#0a-仕様成熟度ポリシー)〜[§0f](../SPEC.md#0f-phase-実装単位ポリシー) を参照し、詳細本文の不足は先に [`docs/details/runner.md` 詳細本文責務 §10a](runner.md#10a-ci-ランナー-実装対象) を改訂する。

| 項目 | 関連節 | 実装内容 |
|------|--------|------------|
| JSON 形式の SHA キャッシュ | [`docs/details/runner.md` 詳細本文責務 §11](runner.md#11-ci-ランナー-ファイル構成)〜[`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) | 各ターゲットの `sha_file` を `{"sha": "..."}` JSON 形式で読み書きする。 |
| `BRANCH_TARGETS` | [`docs/details/runner.md` 詳細本文責務 §12](runner.md#12-設定値runner)〜[`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) | 複数ブランチ、複数 target file、複数出力先を 1 つの設定リストとして処理する。 |
| GitHub API リトライ | [`docs/details/runner.md` 詳細本文責務 §12](runner.md#12-設定値runner)〜[`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) | `API_RETRY_MAX`、`API_RETRY_BASE_SECONDS`、指数バックオフ、レート制限待機を実装する。 |
| ビルドロック | [`docs/details/runner.md` 詳細本文責務 §11](runner.md#11-ci-ランナー-ファイル構成)〜[`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) | `.build_lock` による多重起動防止を実装する。 |
| ビルドクールダウン | [`docs/details/runner.md` 詳細本文責務 §12](runner.md#12-設定値runner)〜[`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) | `.server_config.build_cooldown_seconds` による起動抑制を実装する。 |
| 強制再ビルド間隔 | [`docs/details/runner.md` 詳細本文責務 §12](runner.md#12-設定値runner)〜[`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) | `.server_config.force_build_interval_hours` による変更なし時の定期強制ビルドを実装する。 |
| コミット情報記録 | [`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) | ビルドトリガー commit の SHA、message、author、date をビルドログへ記録する。 |
| 事前チェック | [`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) | ディスク空き容量、`adlaire-ci-build` 実行可否、Go 版ビルドバイナリ配置を確認する。 |
| Webhook 通知 | [`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) | `.notify_config` 読み込み、成功/失敗/転送失敗/週次サマリー通知、`.notify_pending` 再送を実装する。 |
| ビルドログ保存 | [`docs/details/runner.md` 詳細本文責務 §11](runner.md#11-ci-ランナー-ファイル構成)〜[`docs/details/runner.md` 詳細本文責務 §15](runner.md#15-ログ) | `.build_logs/{id}.json` へ stdout/stderr、変換レポート、所要時間を保存する。 |
| SSH 転送 | [`docs/details/runner.md` 詳細本文責務 §14a](runner.md#14a-ssh-サイト転送) | SHA256 差分検出、stdin パイプ転送、転送後整合性検証、ペンディングキューを実装する。 |
| スナップショット | [`docs/details/runner.md` 詳細本文責務 §14b](runner.md#14b-スナップショット管理) | 作成条件を判定し、`archive` owner の snapshot save を呼び出し、結果を build log / history へ反映する。保存形式、検証、世代削除、download、delete、rollback 実体は [`docs/details/archive.md` 詳細本文責務 §27.15](archive.md#sec-27-15) を参照する。 |
| サーキットブレーカー | [`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) | `.build_circuit_state` による連続失敗停止を実装する。 |
| ログ世代管理 | [`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) | `LOG_KEEP_N` による `.build_logs/` 削除を実装する。 |
| 出力サイズ警告 | [`docs/details/runner.md` 詳細本文責務 §12](runner.md#12-設定値runner)〜[`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) | `OUTPUT_SIZE_WARN_MB` による WARN ログと `size_warn` 記録を実装する。 |

**初期実装対象外の連携範囲：**

[`docs/details/runner.md` 詳細本文責務 §10](runner.md#10-ci-ランナー-要件)〜[`docs/details/runner.md` 詳細本文責務 §20](runner.md#20-ci-ランナー-既知の制限) には、`runner` 単体の責務ではなく管理 API、標準管理ツール、追加の運用機能と結合して成立する項目が含まれる。これらは、API・SDK・UI の対象節に、呼び出し元、呼び出し先、状態ファイル、失敗時応答、検証条件が定義されるまで `runner` 単体で実装しない。

| 項目 | 理由 |
|------|------|
| API 経由の動的ブランチ設定 | `runner` 単体では設定 API を持たないため、`api` 実装と合わせて扱う。 |
| API 経由のロールバック | `POST /api/history/{id}/rollback` は `api` のエンドポイント実装が前提となる。 |
| 管理画面からのスケジュール操作 | systemd timer の変更 API と標準管理ツール UI が前提となる。 |

---

## 11. CI ランナー ファイル構成

[`docs/details/runner.md` 詳細本文責務 §11](runner.md#11-ci-ランナー-ファイル構成) は、Go 版 CI ランナー固有の実行物、入力 checkout、出力先だけを示す。runtime 状態ファイルの path、形式、初期値、更新責務、破損時処理は [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a)、リポジトリ内ファイルの実在有無と役割は [`docs/DOCUMENT_INDEX.md`](../DOCUMENT_INDEX.md)、インストール先と unit 配置は [`docs/details/setup.md` 詳細本文責務 §26](setup.md#26-セットアップアップデート手順) を正本とする。

**Go 版 CI ランナーで使用するファイル：**

| パス | 用途 |
|------|------|
| `/usr/local/bin/adlaire-ci-runner` | `components/runner.go` から生成する CI ランナーバイナリ。 |
| `/usr/local/bin/adlaire-ci-build` | `components/builder.go` から生成する Markdown → 静的 Web サイトビルドバイナリ。 |
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

**runtime 状態参照：**

runner は [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) に登録され、各 runner 処理節の読取対象または更新対象に明示された状態だけを読み書きする。runner 詳細本文責務で状態ファイル一覧、API / SDK / UI 配置、schema を再定義してはならない。

**出力先・配信先：**

```
/opt/adlaire-builder/dist/
└── site/                # 静的 Web サイト出力先

/var/www/html/           # SSH 転送先
```

**systemd 配置参照：**

以下は runner が起動される運用上の配置である。unit 本文、配置、enable、restart、更新、rollback は [`docs/details/setup.md` 詳細本文責務 §26](setup.md#26-セットアップアップデート手順) を参照する。

```
/etc/systemd/system/
├── adlaire-ci.service      # systemd ユニット（oneshot）
├── adlaire-ci.timer        # systemd タイマー（定期実行）
└── adlaire-ci-api.service  # 管理 API サーバー（常駐）
```

**リポジトリ側：**

```
<repo>/
├── docs   # ソース Markdown（GitHub 上のマスター）
└── .ci/
    └── pipeline.sh      # ビルド手順定義（実行権限付き）
```

---

## 12. 設定値（`runner`）

`runner` は [`docs/details/runner.md` 詳細本文責務 §12](runner.md#12-設定値runner) の設定値を基準とする。設定値は Go 構造体の既定値、設定ファイル、または CLI 引数で与える。どの入力経路を採用する場合でも、内部表現は [`docs/details/runner.md` 詳細本文責務 §12](runner.md#12-設定値runner) のキー名・型・既定値に従う。

**関連型：**

```go
type RunnerConfig struct {
    StateDir                    string
    RepositoryOwner            string
    RepositoryName             string
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
    WatchMode                   string
    TagFilter                   TagFilterConfig
    BuildCacheEnabled           bool
    DeployParallelism           int
    RemoteBuild                 RemoteBuildConfig
    ApprovalTimeoutSeconds      int
    BranchTargets               []BranchTarget
}

type BranchTarget struct {
    Branch           string
    TargetFile       string
    TargetFiles      []string
    SHAFile          string
    Src              string
    Out              string
    ApprovalRequired bool
    Env              map[string]string
    DeployTargets    []DeployTarget
}

type DeployTarget struct {
    ID      string
    Host    string
    User    string
    DestDir string
}

type TagFilterConfig struct {
    Enabled  bool
    Patterns []string
}

type RemoteBuildConfig struct {
    Enabled      bool
    Host         *string
    User         *string
    WorkDir      *string
    CommandArgs  []string
    ArtifactPath *string
}
```

**設定値検証：**

| 項目 | 条件 | 不正時 |
|------|------|--------|
| `StateDir` | 絶対パス。空文字不可。 | 終了コード `2` |
| `RepositoryOwner` | [`docs/details/statefile.md` 詳細本文責務 `.repo_config` schema](statefile.md#repo-config-schema) の `owner` 許容値に完全一致する。 | 終了コード `2` |
| `RepositoryName` | [`docs/details/statefile.md` 詳細本文責務 `.repo_config` schema](statefile.md#repo-config-schema) の `repo` 許容値に完全一致する。 | 終了コード `2` |
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
| `WatchMode` | `"github"` または `"local"`。 | 終了コード `2` |
| `TagFilter` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の TagFilter object。`WatchMode="local"` と `Enabled=true` の併用禁止。 | 終了コード `2` |
| `DeployParallelism` | 1〜16。 | 終了コード `2` |
| `RemoteBuild` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の RemoteBuildConfig object。 | 終了コード `2` |
| `ApprovalTimeoutSeconds` | 60〜2592000。 | 終了コード `2` |
| `BranchTargets` | 1 件以上。 | 終了コード `2` |
| `BranchTarget.Branch` | 空文字不可。`..`、`~`、制御文字禁止。 | 終了コード `2` |
| `BranchTarget.TargetFile` | 相対パス。絶対パス、`..`、先頭 `/` 禁止。 | 終了コード `2` |
| `BranchTarget.TargetFiles` | 1〜100 件。各 path は `TargetFile` と同じ制約を満たし、`TargetFile` を含む。 | 終了コード `2` |
| `BranchTarget.SHAFile` / `Src` / `Out` | 絶対パス。空文字不可。 | 終了コード `2` |
| `BranchTarget.Env` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.branch_config.env` schema。 | 終了コード `2` |
| `DeployTarget.ID` / `Host` / `User` / `DestDir` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.branch_config.deploy_targets[]` schema に完全一致する。 | 当該 deploy target を無効として ERROR ログ。全 target 無効なら終了コード `2` |

`.repo_config`、`.server_config`、`.branch_config`、CLI 引数から読み込んだ値は `RunnerConfig` に正規化してから使用する。正規化後の `RunnerConfig` にないキーを処理フローで直接参照してはならない。

**設定入力の優先順位：**

1. CLI 引数
2. `.repo_config` / `.server_config` / `.branch_config` など状態ファイルの保存値
3. [`docs/details/runner.md` 詳細本文責務 §12](runner.md#12-設定値runner) の既定値

同一キーが複数の入力経路に存在する場合は、上位の値だけを採用する。採用しなかった値を混合してはならない。未知 CLI 引数は終了コード `2` とし、状態ファイルの未知 key は [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の schema 厳格化契約と [`docs/details/runner.md` 詳細本文責務 §27.10](runner.md#sec-27-10) の起動時整合性チェックに従って破損として扱う。未知値を WARN だけで無視して処理を続行してはならない。

**runner CLI 引数仕様：**

| 引数 | 必須 | 既定値 | 説明 |
|------|------|--------|------|
| `--state-dir <path>` | 任意 | `/opt/adlaire-builder` | 状態ファイル、repo、dist、admin の基準ディレクトリ。相対パスは禁止。 |
| `--once` | 任意 | `true` | 1 回だけ実行して終了する。Go 版 runner は oneshot 固定のため、指定してもしなくても同じ挙動とする。 |
| `--dry-run` | 任意 | `false` | 設定、状態、GitHub target、SHA 差分、起動可否だけを検証し、ビルド、deploy、通知、状態ファイル更新を行わず終了する。 |
| `--version` | 任意 | なし | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の CLI 共通固定契約に従う。 |
| `--help` | 任意 | なし | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の CLI 共通固定契約に従う。 |

未知引数、値欠落、相対 `--state-dir` は終了コード `2` とし、ビルド処理を開始しない。

**runner CLI パース固定仕様：**

- runner CLI の parse 形式と短縮 option 禁止は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の CLI 共通固定契約に従う。
- 同一引数が複数回指定された場合は最後の値を採用する。ただし `--once` は指定有無にかかわらず `true` として扱う。`--dry-run` は 1 回以上指定されれば `true` とする。
- `--help` と `--version` は他の引数より優先し、`.github_token` 読み込み、lock 作成、状態ファイル読み込み、GitHub API 呼び出しを行わない。
- stdout / stderr 行末と単一エラー出力は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の CLI 共通固定契約に従う。
- runner 固有の重複 option、`--once`、`--dry-run` の扱いは [`docs/details/runner.md` 詳細本文責務 §12](runner.md#12-設定値runner) を正本とする。

**runner 設定正規化契約：**

runner は CLI、`.repo_config`、`.server_config`、`.branch_config`、既定値を読み込んだ後、処理開始前に 1 回だけ `RunnerConfig` へ正規化する。正規化前の map、JSON raw message、環境変数、CLI flag 値を、GitHub API、pipeline、deploy、snapshot、通知処理から直接参照してはならない。

| 入力 | 正規化先 | 正規化ルール |
|------|----------|--------------|
| `--state-dir` | `RunnerConfig.StateDir` | `filepath.Clean` 後も絶対パスであることを確認する。末尾 `/` の有無で別パス扱いしない。 |
| `readRepoConfig().owner` | `RunnerConfig.RepositoryOwner` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) で schema 検証済みの値をそのまま採用する。 |
| `readRepoConfig().repo` | `RunnerConfig.RepositoryName` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) で schema 検証済みの値をそのまま採用する。 |
| `PENDING_FILE` | `RunnerConfig.PendingFile` | CLI で `StateDir` が変更された場合、明示設定がない限り `{StateDir}/.pending_transfers` に再解決する。 |
| `.branch_config.branch_targets[]` | `RunnerConfig.BranchTargets` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) に従い、`target_files`、`approval_required`、`env`、deploy target id を含む全定義 key を採用する。同節が定める branch 一意性と正規化済みの branch 名昇順、deploy target id 昇順を変更せず処理順に使用する。 |
| `BRANCH_TARGETS` 既定値 | `RunnerConfig.BranchTargets` | `.branch_config` が存在しない場合だけ使用する。`.branch_config` が破損復旧で不在化された場合も同じ既定値へ fallback する。 |
| `.server_config` の数値 | `RunnerConfig` の数値 field | JSON number が整数でない場合は設定不正とする。文字列数値の暗黙変換は禁止する。 |
| `.server_config` の boolean | `RunnerConfig` の boolean field | JSON boolean のみ許可する。`"true"`、`1`、`"yes"` は不正値とする。 |
| `.server_config.watch_mode` | `RunnerConfig.WatchMode` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の既定値と列挙値を適用する。 |
| `.server_config.tag_filter` | `RunnerConfig.TagFilter` | pattern の入力順を保持する。 |
| `.server_config.build_cache_enabled` | `RunnerConfig.BuildCacheEnabled` | boolean をそのまま採用する。 |
| `.server_config.deploy_parallelism` | `RunnerConfig.DeployParallelism` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の schema 検証済み値を採用する。 |
| `.server_config.remote_build` | `RunnerConfig.RemoteBuild` | object 全体を 1 単位として正規化し、field ごとに別 input source を混合しない。 |
| `.server_config.approval_timeout_seconds` | `RunnerConfig.ApprovalTimeoutSeconds` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の schema 検証済み値を採用する。 |

正規化後は、全 path を絶対パス文字列として保持する。`BranchTarget.TargetFile` だけは GitHub repository 内の相対パスとして保持し、[`docs/details/statefile.md` 詳細本文責務 `.branch_config` schema](statefile.md#branch-config-schema) の path 検証に一致しない場合は終了コード `2` とする。`BranchTarget.Src` と `BranchTarget.Out` が同一または親子関係になる場合は終了コード `2` とし、ERROR ログ `CONFIG_PATH_CONFLICT: src={src} out={out}` を出す。`.repo_config` が破損または読込不能の場合は終了コード `2`、ERROR `REPO_CONFIG_INVALID` とし、GitHub API、queue 取得、build id 採番、pipeline、deploy、snapshot を開始しない。

**状態ディレクトリ構造検証契約：**

`--state-dir` 検証後、runner は以下を処理開始前に確認する。

| 対象 | 条件 | 不正時 |
|------|------|--------|
| `StateDir` | directory、owner が実行ユーザーまたは書き込み可能、mode に owner write がある。 | 終了コード `2`、ERROR `STATE_DIR_INVALID: path={path}`。 |
| `{StateDir}/repo` | 不在なら作成する。file の場合は不正。 | 作成失敗または file の場合、終了コード `2`。 |
| `{StateDir}/dist` | 不在なら作成する。file の場合は不正。 | 作成失敗または file の場合、終了コード `2`。 |
| `{StateDir}/.build_logs` | 不在なら `0700` で作成する。 | 作成失敗時は終了コード `2`。 |

対象 directory 作成は dry-run では実行しない。dry-run では作成予定を [`docs/details/runner.md` 詳細本文責務 §27.2](runner.md#sec-27-2) の stdout JSON `warnings[]` に `code="DRY_RUN_WOULD_CREATE_DIR"` として出し、`would_write` に `"state_dir"` を追加してはならない。ただし既存 path が file の場合は dry-run でも終了コード `2` とし、stdout JSON `errors[]` に原因を出す。

**固定出力：**

| 条件 | stdout |
|------|--------|
| `--help` | `Usage: adlaire-ci-runner [--state-dir path] [--once] [--dry-run] [--version] [--help]` |
| `--version` | `adlaire-ci-runner v3 go=<runtime.Version()>` |

**CLI 異常系：**

| 条件 | 終了コード | stderr |
|------|------------|--------|
| 未知の引数 | `2` | `unknown option: <name>` |
| `--state-dir` 値欠落 | `2` | `missing value: --state-dir` |
| `--state-dir` が空文字 | `2` | `state directory must not be empty` |
| `--state-dir` が相対パス | `2` | `state directory must be absolute: <path>` |
| `--state-dir` が存在しない | `2` | `state directory not found: <path>` |
| `--state-dir` がディレクトリではない | `2` | `state path is not directory: <path>` |

以下の `BRANCH_TARGETS`、`PENDING_FILE`、`API_RETRY_MAX`、`.server_config.build_cooldown_seconds`、`.server_config.snapshots_keep`、`.server_config.force_build_interval_hours`、`LOG_KEEP_N`、`API_CIRCUIT_BREAKER_THRESHOLD`、`OUTPUT_SIZE_WARN_MB`、`WEEKLY_SUMMARY_*` を Go 版 runner の標準設定とする。`.server_config` の正確な key、型、既定値、許容範囲、0 の意味は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) を正本とする。

```text
PENDING_FILE           = "/opt/adlaire-builder/.pending_transfers"   # SSH 転送ペンディングキュー（JSON）
API_RETRY_MAX          = 5    # GitHub API 失敗時の最大再試行回数（指数バックオフ）
API_RETRY_BASE_SECONDS = 1    # バックオフ基底秒数（1→2→4→8→16 秒。0 = リトライ無効）
server_config.build_cooldown_seconds = 60   # 前回ビルド完了から次ビルドまでの最小間隔（秒。0 = 無効）
server_config.force_build_interval_hours = 0 # 強制再ビルド間隔（時間。0 = 無効）
LOG_KEEP_N             = 50   # ビルドログ保持件数（0 = 無制限）
API_CIRCUIT_BREAKER_THRESHOLD = 3    # 全ブランチ連続失敗の許容周回数（0 = 無効）
OUTPUT_SIZE_WARN_MB           = 5    # 出力サイト合計サイズ警告閾値（MB。0 = 無効）
WEEKLY_SUMMARY_ENABLED        = true # 週次サマリー Webhook の有効/無効
WEEKLY_SUMMARY_DAY            = 0    # 送信曜日（0=月曜〜6=日曜）
WEEKLY_SUMMARY_HOUR           = 9    # 送信時刻（0〜23、ローカル時刻）

# ブランチターゲット設定
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

<a id="runner-終了コード"></a>
**runner 終了コード：**

| 終了コード | 条件 |
|------------|------|
| `0` | 起動、ロック確認、対象処理が完了した。変更なし、クールダウン、既存ロックによるスキップも正常終了に含める。 |
| `1` | 1 件以上のビルドまたは転送が失敗したが、runner 自体は最後まで処理できた。 |
| `2` | 設定不正、必須ファイル不足、CLI 引数不正。 |
| `3` | GitHub API など外部サービスへの全再試行が失敗し、全ターゲットが処理不能。 |
| `4` | `.build_lock` 作成に失敗し、既存 PID の実行中確認もできない。 |

systemd timer からの再実行を妨げないため、終了コード `1` と `3` でも、[処理フロー](#13-処理フロー) の failure 別書込順と [ビルド通知連携](#sec-27-32) に従い、該当する log / history / status / state / notification pending の書込と所有確認付き lock 解放を各契約の最大回数で実行してから終了する。本文だけを根拠に追加 retry を行ってはならない。

<a id="atomic-write-共通契約"></a>
**atomic write 共通契約：**

runner が JSON object、JSON array、SHA cache、`.build_state`、`.build_circuit_state`、`.notify_pending`、`.pending_transfers`、`.server_config`、`.branch_config` を更新する場合、path、schema、lock、tmp、mode、rename、file sync、parent directory sync、失敗時の target 維持は [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の状態ファイル更新手順に従う。runner は独自の atomic write 手順、tmp 名、sync 成否判定、lock 手順を定義または実装してはならない。

runner は業務処理上の複数ファイル書込順と、statefile write が返した失敗後に実行する処理または停止する処理だけを各機能節で定義する。statefile write が失敗を返した場合、runner は対象ファイルを更新済みとして扱わず、成功扱いへ読み替えない。

<a id="lock-ファイル契約"></a>
**lock ファイル契約：**

`.build_lock` は UTF-8 text で、内容は `pid={pid}\nstarted_at={UTC_ISO8601}\n` とする。

| 状態 | 処理 |
|------|------|
| lock なし | `os.OpenFile` を `O_CREATE`、`O_EXCL`、`O_WRONLY` の各 flag と mode `0600` で呼び出して作成する。 |
| lock あり・PID 実行中 | INFO ログ `BUILD_SKIP: already running (PID {pid})` を出し終了コード `0`。 |
| lock あり・PID 不在 | WARN ログ `STALE_LOCK: pid={pid}` を出し、lock を削除してから再作成する。 |
| lock あり・PID 解析不能 | ERROR ログ `LOCK_CORRUPT: path={path}` を出し終了コード `4`。 |
| lock 作成失敗 | ERROR ログ `LOCK_CREATE_FAILED: {reason}` を出し終了コード `4`。 |

PID 実行中判定は Linux の `/proc/{pid}` 存在確認で行う。`/proc` を読めない場合は PID 実行中確認不能として終了コード `4` とする。

`.build_lock` の作成、stale 再確認、所有確認付き解放は runner owner だけが実装する。API owner が rollback または snapshot delete を実行する場合は runner owner の coordinator を呼び出し、API owner が lock file を直接書き換えない。

| 呼出用途 | lock 取得後の固定処理 | conflict 時 |
|----------|--------------------------|-------------|
| 通常 build | startup integrity、queue / polling、build lifecycle へ進む。 | [runner 終了コード](#runner-終了コード) に従う。 |
| rollback | [`docs/details/runner.md` 詳細本文責務 §14b](runner.md#14b-スナップショット管理) の rollback accepted 確定順へ進む。 | runner owner が conflict を返し、API owner が `409` へ変換する。 |
| snapshot delete 排他 | lock 取得後に `.build_state` を再読込し、`running=false`、`current_build_id=null` を確認した場合だけ一回性 guard を API owner へ返す。`active_queue_entry` と waiting queue は判定対象外とし、変更しない。 | lock 存在、lock 形式不正、PID 判定不能、state running は guard を返さず conflict とする。state 破損または読取失敗は state error とする。新規 lock を取得済みの場合はいずれも所有確認後に解放する。 |

lock 解放時は、開始時に自身が書いた `pid` と `started_at` の完全一致を再読取で確認した場合だけ削除する。不一致、読取不能、削除失敗は他操作の lock とみなして上書きせず、caller へ lock release failure を返す。snapshot delete guard は archive delete、`.config_log`、`.audit_log` の成否にかかわらず最終処理で解放し、lock release failure で先行する snapshot 削除または log 追記を巻き戻さない。

<a id="github-token-読み込み契約"></a>
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

<a id="runner-状態ファイル参照契約"></a>
**runner 状態ファイル参照契約：**

状態ファイルの path、形式、初期値、schema、破損時の扱い、atomic write、adapter は [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) を参照する。`runner` 詳細では、runner がどの処理段階で状態を読むか、いつ更新するか、失敗時に後続処理を止めるかだけを定義する。API endpoint ごとの状態読取順と response 算出は [`docs/details/api.md` 詳細本文責務 §22.0c.1](api.md#sec-22-0c-1) を参照する。

runner 実装は [`docs/details/statefile.md`](statefile.md) 詳細本文責務に未定義の状態ファイル、永続 key、queue entry key、notification entry key、pending transfer entry key を追加してはならない。

`.branch_config` が存在しない場合は [`docs/details/runner.md` 詳細本文責務 §12](runner.md#12-設定値runner) の `BRANCH_TARGETS` 既定値を使用する。存在する場合の schema、空配列の扱い、永続 key、API 表示名との境界は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.branch_config` schema を参照する。

runner は起動時の設定正規化で `.branch_config` を 1 回だけ読み、正規化後の `RunnerConfig.BranchTargets` を当該起動中の唯一の branch target 情報として使用する。同一 runner 起動中に `.branch_config` を再読込してはならない。API による `.branch_config` 更新、削除、default 復帰は、既に実行中の runner には反映せず、次回 runner 起動から反映する。

`.pending_transfers` の entry schema、投入順、転送先識別子、新しい build による置換、成功時削除、成果物 checksum 不一致時の扱いは [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.pending_transfers` schema を参照する。runner はその schema を独自に拡張、別名化、index 化してはならない。

<a id="sha-cache-読み書き契約"></a>
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
| 未知 key あり | schema 破損として `failure_decode` で当該 target を失敗扱いし、sha_file を更新しない。 |

SHA cache の更新は、pipeline 成功後、deploy 前に行う。複数 target のうち一部 target が成功した場合は、成功 target の `sha_file` だけを更新する。失敗 target、skip target、branch target 設定不正 target の `sha_file` を更新してはならない。

<a id="状態ファイル権限契約"></a>
**状態ファイル権限契約：**

状態ファイルと lock file、状態ディレクトリの mode は [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の固定値を使用する。runner は起動時に mode を確認し、`.github_token`、`.notify_config`、`.notify_pending`、`.pending_transfers` が固定値より広い場合は自動補正せず終了コード `2` で停止する。それ以外の runner 状態ファイルが固定値より広い場合は WARN `STATE_FILE_INSECURE_MODE: path={path} mode={mode}` を出し、statefile の固定値へ chmod する。chmod 失敗時は終了コード `2` とする。

<a id="設定ファイル起動時整合性チェック"></a>
**設定ファイル起動時整合性チェック：**

対象ファイル、実行順序、判定分類、復旧、停止条件、通知、状態更新禁止条件、検証条件は [`docs/details/runner.md` 詳細本文責務 §27.10](runner.md#sec-27-10) にのみ定義する。[`docs/details/runner.md` 詳細本文責務 §12](runner.md#12-設定値runner) では再定義しない。

<a id="build-id-契約"></a>
**build id 契約：**

runner が生成する build id の base は UTC 時刻ベースの `b{YYYYMMDDHHmmss}` とし、衝突 suffix は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の時刻ベース ID 契約に従う。build id は `.build_logs/{id}.json`、`.build_history`、`.snapshots/{id}/` で同一値を使用し、suffix 上限到達時は既存 ID を上書きせず終了コード `1` とする。

---

## 13. 処理フロー

[`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) の処理フローは、`runner` の標準フローである。

**状態更新順序の規範：**

1. `.build_lock` を作成する。
2. 設定ファイル起動時整合性チェックを実行する。
3. `.build_status.json` に `status="running"`、`trigger`、`started_at`、`current_build_id` を atomic write で保存する。
4. `.build_state.running=true`、`current_build_id`、`last_started_at` を atomic write で保存する。
5. build id 確定後、GitHub API、pipeline、deploy などの外部副作用を開始する前に、`.build_logs/{id}.json` を `status="running"` の途中保存形で atomic write する。
6. GitHub API、Blob 書き出し、事前チェック、`pipeline.sh` 実行を行う。
7. ビルド成功時のみ `sha_file` を新 SHA に更新し、deploy を実行する。
8. deploy target がないか、全 deploy target の転送と検証が成功した場合だけ、[`docs/details/runner.md` 詳細本文責務 §14b](runner.md#14b-スナップショット管理) に従って `archive` owner の snapshot save を呼び出す。
9. stdout/stderr、`[REPORT]`、警告、転送結果、snapshot 結果、`trigger` を含む最終形で、同じ build id の running log を atomic replace する。
10. `.build_history` に同じ `id` の最終要約行を JSON Lines で 1 回だけ追記する。
11. `.build_status.json` に最終状態を atomic write で保存する。
12. `.build_state.running=false`、`last_finished_at` を保存する。
13. `.build_lock` を削除する。

途中失敗時は、失敗が発生した段階以降の成功前提更新を行わない。例えば `pipeline.sh` 失敗時は `sha_file`、snapshot、転送成功履歴を更新しない。ただし `.build_logs/{id}.json`、`.build_history`、`.build_state.running=false`、通知 pending は失敗記録として保存する。

<a id="runner-target-status-selection"></a>
**ターゲット結果分類：**

runner は `BRANCH_TARGETS` の各 entry について、[`docs/details/statefile.md` 詳細本文責務 §22.0c runner 結果値 schema](statefile.md#runner-result-schema) で build log に許可された値から、最終的に 1 つの `target_status` を確定する。基本 build lifecycle の選択条件は [`docs/details/runner.md` 詳細本文責務 ターゲット結果分類表](runner.md#runner-target-status-selection) を正本とする。追加機能は runner 結果値 schema の許可値と、対象 owner component の機能契約に定めた選択条件の両方に従う。

| `target_status` | 条件 | runner 終了コードへの影響 |
|-----------------|------|---------------------------|
| `skipped_no_change` | SHA 一致かつ強制ビルド条件未達。 | 失敗扱いしない。 |
| `skipped_cooldown` | クールダウンで全体処理を開始しない。 | `0`。 |
| `success` | pipeline 成功、SHA 更新成功、deploy target がないか全 deploy target が成功。 | `0` 候補。 |
| `success_deploy_pending` | pipeline 成功、SHA 更新成功、1 件以上の deploy が pending。 | `1`。 |
| `failure_api` | GitHub API が全再試行失敗。 | 全 target がこれなら `3`、一部なら `1`。 |
| `failure_decode` | Blob Base64 decode または Markdown 書き出し失敗。 | `1`。 |
| `failure_precheck` | 事前チェック失敗。 | `1`。 |
| `failure_build` | `pipeline.sh` または必須 pipeline step が非 0 終了。 | `1`。 |
| `failure_timeout` | builder、pipeline step、remote build が timeout。 | `1`。 |
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
| lock 実行中 PID | `lock_skipped` | lock が有効な実行中 PID を指す場合に限り、`.build_status.json` に `status="skipped"`、`last_target_status="lock_skipped"` を保存する。 | `.build_state`、`.build_logs/`、`.build_history`、`sha_file`、deploy、snapshot。 | lock を削除しない。 | `0` |
| startup config permission error | `config_error` | `.build_status.json` に `status="failure"`、`last_target_status="config_error"`、ERROR log。 | `.build_state.running=true`、`.build_logs/`、`.build_history`、`sha_file`、deploy、snapshot。 | lock 作成済みなら削除する。 | `2` |
| maintenance state unavailable | `config_error` | [復旧 record 固定契約](#sec-27-10) の status failure 条件が成立する場合だけ、`.build_status.json` に `status="failure"`、`last_target_status="config_error"`、`last_trigger="startup_config_integrity"`、ERROR log を保存する。 | `.maintenance` 自動修復 / backup / 上書き、pending retry、`.build_state.running=true`、`.build_logs/`、`.build_history`、`.last_sha`、`.sha_cache/`、deploy、snapshot。 | lock 作成済みなら削除する。 | `2` |
| status start write failure | `failure_state_write` | ERROR log。 | `.build_state.running=true`、`.build_logs/`、`.build_history`、`sha_file`、deploy、snapshot。 | lock 作成済みなら削除する。 | `1` |
| state start write failure | `failure_state_write` | `.build_status.json` に `failure`、ERROR log。 | `.build_logs/`、`.build_history`、`sha_file`、deploy、snapshot。 | lock を削除する。 | `1` |
| GitHub tree / blob API failure | `failure_api` | `.build_logs/{id}.json`、`.build_history`、`.build_status.json`、`.build_state.running=false`。 | `sha_file`、deploy、snapshot、materialized src の確定置換。 | 実行する。 | 全 target 失敗なら `3`、一部なら `1` |
| blob decode / materialize failure | `failure_decode` | `.build_logs/{id}.json`、`.build_history`、`.build_status.json`、`.build_state.running=false`。 | `sha_file`、pipeline、deploy、snapshot。 | 実行する。 | `1` |
| precheck failure | `failure_precheck` | `.build_logs/{id}.json`、`.build_history`、`.build_status.json`、`.build_state.running=false`。 | `sha_file`、pipeline、deploy、snapshot。 | 実行する。 | `1` |
| pipeline non-zero | `failure_build` | `.build_logs/{id}.json`、`.build_history`、`.build_status.json`、`.build_state.running=false`。 | `sha_file`、deploy、snapshot。 | 実行する。 | `1` |
| builder / pipeline / remote timeout | `failure_timeout` | `.build_logs/{id}.json`、`.build_history`、`.build_status.json`、`.build_state.running=false`。 | `sha_file`、deploy、snapshot。 | 実行する。 | `1` |
| build log start write failure | `failure_state_write` | `.build_status.json`、`.build_state.running=false`、ERROR log。 | `.build_history`、`sha_file`、deploy、snapshot。 | 実行する。 | `1` |
| build log final write failure | `failure_state_write` | 既存 running log、`.build_status.json`、`.build_state.running=false`、ERROR log。 | `.build_history` 追記、完了値の推測補完。 | 既に完了した `sha_file`、deploy、snapshot は巻き戻さず、status finalizer、state finalizer、所有確認付き lock 解放をそれぞれ 1 回だけ実行する。 | `1` |
| build history append failure | `failure_state_write` | 最終形の `.build_logs/{id}.json`、`.build_status.json`、`.build_state.running=false`、ERROR log。 | 同じ id の history 再追記、完了値の推測補完。 | 既に完了した `sha_file`、deploy、snapshot は巻き戻さず、status finalizer、state finalizer、所有確認付き lock 解放をそれぞれ 1 回だけ実行する。 | `1` |
| SHA write failure | `failure_state_write` | `.build_logs/{id}.json`、`.build_history`、`.build_status.json`、`.build_state.running=false`。 | deploy、snapshot。 | 実行する。 | `1` |
| deploy failure | `success_deploy_pending` | `sha_file`、`.pending_transfers`、`.build_logs/{id}.json`、`.build_history`、`.build_status.json`、`.build_state.running=false`。 | snapshot。 | 実行する。 | `1` |
| snapshot failure | `success` | `sha_file`、deploy 成功、`.build_logs/{id}.json` に WARN、`.build_history`、`.build_status.json`、`.build_state.running=false`。 | snapshot 成功扱い、`.pending_transfers` 追加。 | 実行する。 | `0` |
| status finalizer write failure | 実行結果に従う | `.build_logs/{id}.json`、`.build_history`、`.build_state.running=false`、ERROR log。 | 部分 `.build_status.json`。 | 継続する。 | 最低 `1` |
| build_state finalizer write failure | 実行結果に従う | `.build_logs/{id}.json`、`.build_history`、`.build_status.json`、ERROR log。 | 正常終了扱い。 | 所有確認付き lock 解放を 1 回だけ実行する。 | `1` |

[`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) の固定表の保存必須に含まれる状態ファイルは、保存失敗時に `failure_state_write` へ分類する。ただし status finalizer と build_state finalizer は、既に確定した target の log / history を取り消さない。保存禁止に含まれる処理を実行した場合は仕様違反とし、実装変更の fixture で失敗として扱う。

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

`runner` は、[`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) の固定表の機能単位で状態を更新する。各機能単位は、Write 列にない状態ファイルを更新してはならない。

| 機能単位 | Read | Write | 成功条件 | 失敗時更新 |
|----------|------|-------|----------|------------|
| lock acquisition | `.build_lock` | `.build_lock` | lock file を `O_CREATE` と `O_EXCL` の両 flag で作成し、PID と started_at を書き込む。 | 実行中 PID がある場合は終了コード `0`。形式不正 lock は終了コード `4`、上書き禁止。 |
| startup config integrity | `.branch_config`, `.notify_config`, `.build_state`, `.build_circuit_state`, `.pending_transfers`, `.notify_pending`, `.maintenance` | 復旧対象 6 file と corrupt backup。`.maintenance` は書込禁止。 | 対象 JSON の parse、top-level type、必須 key、型、値範囲、未知 key を検証する。`.maintenance` 不在は disabled として継続する。 | 復旧対象 6 file の回復可能な破損は backup と初期化後に継続。`.maintenance` 破損、unknown key、読取不能は自動修復せず終了コード `2`。その他の権限エラーは `2`、IO エラーは `1`。 |
| pending transfer retry | `.pending_transfers` | `.pending_transfers`, `.build_logs/{id}.json` | 成功 entry を削除し、失敗 entry は `retry_count` を +1 して保持する。 | queue ファイル破損時は [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) に従い退避して `[]` で再生成する。 |
| notify pending retry | `.notify_pending`, `.notify_config` | `.notify_pending`, `.notify_log` | channel 別の成功条件を満たした entry を削除し、失敗 entry は `attempts` を +1 して保持する。 | `.notify_config` 不在時は送信せず entry を保持し、ERROR ログを記録する。 |
| circuit gate | `.build_circuit_state`, `.build_state` | `.build_status.json`、必要時だけ `.build_state` | closed なら queue 取得または polling 判定へ進む。open なら pending retry 以外の build 処理を開始しない。 | open 時は active / waiting を保持し、runtime flag が残る場合だけ `running=false`、`current_build_id=null` へ atomic write する。保存失敗は終了コード `1`。 |
| status start | `.build_state`, `.build_circuit_state`, `.pending_transfers`, `.notify_pending` | `.build_status.json` | circuit closed を確認し、queue entry または polling target を確定した後、`status="running"` と当該起動の `trigger` を保存する。 | 書き込み失敗は終了コード `1`。`.build_state.running=true` へ進まない。 |
| state start | `.build_state` | `.build_state` | queue entry は waiting から active への移動と `running=true`、`current_build_id`、`last_started_at` を同じ atomic write で保存する。polling は active / waiting を変更せず runtime field だけを保存する。 | `.build_lock` を削除し、終了コード `1`。 |
| build log start | build id, `.build_state`, `.build_status.json` | `.build_logs/{id}.json` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の途中保存 schema を atomic write する。 | GitHub API、pipeline、deploy、snapshot を開始せず、`failure_state_write` とする。 |
| cooldown gate | `.build_state`, `.server_config` | none | cooldown 範囲外なら target 処理へ進む。 | cooldown 中は `target_status=skipped_cooldown`、終了コード `0`。 |
| GitHub tree resolve | `.github_token`, branch target | `.build_logs/{id}.json` | target file の blob SHA を取得する。 | 全 retry 失敗は `failure_api` を記録し、SHA を更新しない。 |
| SHA decision | `sha_file`, `.server_config` | none | 変更あり、force 条件成立、または webhook/force queue payload により build 対象を確定する。 | 変更なしは `skipped_no_change`。状態ファイルを更新しない。 |
| blob materializer | GitHub blob API, branch target | `src` | Base64 decode 後、対象 Markdown を atomic write する。 | decode/write 失敗は `failure_decode`。SHA を更新しない。 |
| precheck | branch target, output path, build binary | `.build_logs/{id}.json` | disk、binary、version がすべて合格する。 | `failure_precheck` を記録し、pipeline を起動しない。 |
| pipeline executor | `src`, `.pipeline_config` | none | timeout 前に exit code `0` と stdout / stderr を最終結果 writer へ返す。 | 非 0 は `failure_build`、timeout は `failure_timeout` を返す。SHA、deploy、snapshot を更新しない。 |
| report importer | pipeline stdout | none | `[REPORT]` を 1 行だけ検出し schema へ変換し、最終結果 writer へ返す。 | `[REPORT]` なしは `report:null` と `REPORT_MISSING` warning を返す。build 成否は pipeline exit code に従う。 |
| SHA updater | `sha_file`, new SHA | `sha_file` | `{"sha":"<new_sha>"}` を atomic write する。 | `failure_state_write`。deploy、snapshot を実行しない。 |
| deploy executor | `deploy_targets`, output site | `.pending_transfers`, `.build_logs/{id}.json` | 全 deploy target の転送と検証が成功する。 | 失敗 target を `.pending_transfers` に保存し、`target_status=success_deploy_pending`。 |
| snapshot coordinator | output site, `.server_config`, build id, output manifest SHA-256 | none | `snapshots_keep > 0` で deploy phase 成功時だけ `archive` owner を呼び出し、`saved` または `exists_valid` を `snapshot_id=build_id` として最終結果 writer へ返す。 | `failed` は build 成功を取り消さず `SNAPSHOT_SAVE_FAILED`、`snapshot_id=null` を返す。 |
| build result writer | running build log, target result, report, deploy result, snapshot result | `.build_logs/{id}.json`, `.build_history` | running log を最終 schema で atomic replace した後、同じ id の history を 1 行追記する。 | final log 失敗では history を追記しない。history 失敗では final log を保持し、自動再追記しない。 |
| status finalizer | target results, `.build_state`, `.build_circuit_state`, `.pending_transfers`, `.notify_pending` | `.build_status.json` | runner 起動 1 回の最終要約を保存する。 | 書き込み失敗は ERROR ログを出し、runner 終了コードを最低 `1` にする。 |
| finalizer | target results, `.build_state` | `.build_state`, `.build_lock` | `running=false`、`last_finished_at` を保存し、lock を削除する。 | lock 削除失敗は WARN。`.build_state.running=false` 保存失敗は終了コード `1`。 |

`.build_logs/{id}.json.status` は正規化状態、`.build_logs/{id}.json.target_status` は詳細結果として同時に保存する。`failure_api`、`failure_decode`、`failure_precheck`、`failure_build`、`failure_state_write` および追加機能の failure 値は `status="failure"` へ写像し、詳細値を `target_status` から失わせてはならない。`error` は [`docs/details/runner.md` 詳細本文責務 §15](runner.md#15-ログ) の固定文言を保存する。API の正規化 `status` と詳細結果の response field は [`docs/details/api.md` 詳細本文責務 §22.0c.1](api.md#sec-22-0c-1) を参照する。

**ビルドトリガー種別契約：**

`trigger` は runner が処理を開始した原因を表す固定文字列であり、runner は `.build_logs/{id}.json`、`.build_history`、`.build_status.json` に同じ値を保存する。実装者判断で `auto`、`force`、`scheduled` など別名を追加してはならない。

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

`.build_status.json` の schema、許容値、初期値は [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) を参照する。API endpoint の読取順と response 算出は [`docs/details/api.md` 詳細本文責務 §22.0c.1](api.md#sec-22-0c-1) を参照する。`runner` 詳細では runner が `.build_status.json` を更新するタイミングと、更新失敗時の runner 挙動だけを定義する。

`.build_status.json` の更新タイミングは次に固定する。

| タイミング | `status` | `last_target_status` | `trigger` | 備考 |
|------------|----------|----------------------|-----------|------|
| lock 取得後、target 処理前 | `running` | `null` | 判定済み trigger | `current_build_id` を設定する。 |
| 実行中 lock 検出 | `skipped` | `lock_skipped` | `polling` | `.build_state` と build log は変更しない。 |
| `.maintenance` 破損 / 読取不能 | `failure` | `config_error` | `startup_config_integrity` | [`docs/details/runner.md` 詳細本文責務 §27.10](runner.md#sec-27-10) の status failure 固定契約を 1 回だけ適用し、build log / history は作成しない。 |
| cooldown skip | `skipped` | `skipped_cooldown` | `polling` | `last_build_id` は既存値を保持する。 |
| SHA 変更なし | `skipped` | `skipped_no_change` | `polling` | build log / history は作成しない。 |
| tag filter 不一致 | `skipped` | `skipped_tag_filter` | 実行 trigger | build log / history は作成しない。 |
| maintenance mode | `skipped` | `skipped_maintenance` | `polling` | build log / history は作成しない。 |
| force interval build 成功 | `success` | `success` | `force_interval` | build log / history と同じ build id を保存する。 |
| build 成功・deploy 成功 | `success` | `success` | 実行 trigger | `last_deploy_at` に deploy result 確定時刻、`last_deploy_status="success"` を保存する。 |
| build 成功・deploy pending | `success` | `success_deploy_pending` | 実行 trigger | `last_deploy_at` に pending 確定時刻、`last_deploy_status="pending"` と pending 件数を更新する。 |
| rollback prepared | `running` | `null` | `rollback` | rollback coordinator が status / state / running log をすべて保存した後のみ保持する。 |
| rollback 成功 | `success` | `success` | `rollback` | rollback build log / history と同じ build id、`rollback_from`、`snapshot_id` に対応させる。 |
| rollback deploy pending | `success` | `success_deploy_pending` | `rollback` | pending 保存成功後だけ pending 件数と `last_deploy_status="pending"` を更新する。 |
| rollback 失敗 | `failure` | `failure_build` または `failure_state_write` | `rollback` | 元 snapshot、元 log、過去 history、`.last_sha` を変更しない。 |
| build 失敗 | `failure` | 確定した failure 値 | 実行 trigger | `last_error` に固定エラー文言を保存する。 |
| build cancel | `cancelled` | `cancelled` | 実行 trigger | build log / history と同じ build id を保存する。 |
| circuit open skip | `skipped` | `circuit_open` | `polling` | GitHub API 呼び出し前に保存する。 |
| startup config 復旧あり | `warning` | `config_recovered` | `startup_config_integrity` | build log / history は作成しない。復旧後に通常 build へ進む場合、次の `running` 更新で上書きする。 |
| startup config 停止 | `failure` | `config_error` | `startup_config_integrity` | `.build_state.running=true` へ進まない。 |
| pending transfer retry のみ | `success` または `failure` | 既存値を保持 | `retry_pending_transfer` | build が発生しない場合でも pending 件数を更新する。 |

`.build_status.json` は atomic write 対象であり、書き込み失敗時に部分ファイルを残してはならない。未知 key を追加してはならない。API は GET request で `.build_status.json` を自動修復してはならない。

**queue entry 実行契約：**

runner は process lock、起動時整合性チェック、pending transfer retry、notify pending retry、circuit closed 判定を完了した後にだけ queue entry を取得する。`.build_state.active_queue_entry` が非 `null` ならその entry を再実行対象とする。active が `null` で `.build_state.queued` に entry がある場合は、通常ポーリング対象の前に [`docs/details/runner.md` 詳細本文責務 §27.35](runner.md#sec-27-35) の priority、created_seq、id 順で 1 件だけ選択し、同じ `.build_state` atomic write で waiting queue から削除して `active_queue_entry` へ移し、`running=true`、`current_build_id`、`last_started_at` を確定する。queue entry 処理が成功または失敗として build log と `.build_history` に確定した後、finalizer は active id が実行対象 id と一致することを再確認し、同じ atomic write で `active_queue_entry=null`、`running=false`、`current_build_id=null`、`last_finished_at` を保存する。runner 起動 1 回で複数 queue entry を連続処理してはならない。queue entry の `trigger` が `"manual"` かつ `payload.force=true` の場合は SHA 比較を行わず build を実行する。`trigger` が `"webhook"` の場合は payload の `branch` と `sha` を使用し、branch target に一致しない entry は `failure_api` として記録した後に finalizer で active を消去する。

process crash、強制終了、start state 以後の必須保存失敗により build log / history が最終結果まで確定しなかった場合、`active_queue_entry` を保持する。finalizer を実行できる場合は `running=false` と `current_build_id=null` だけを同じ write で保存し、active を消去しない。次回 runner は stale process lock を処理した後、waiting queue より先に active entry を再実行する。これは at-least-once 実行契約であり、API または runner は未確定 active entry を成功扱いで削除してはならない。再実行では新しい build id を採番し、前回の未確定 `current_build_id` を再利用しない。

API の runner 起動要求は process 起動の契機にすぎない。runner は起動元を trigger として保存せず、queue entry の `trigger` だけを使用する。API が systemd 起動要求を送信した時点では process lock、build id、`running=true` を確定したものとして扱わない。

queue entry の全必須 key、`queued_at`、`requested_by`、priority、created_seq、trigger 別 payload は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の Queue entry schema を正本とする。不正 entry が active または処理順の先頭にある場合、runner はその entry を `failure_decode` として build log / history に記録し、最終結果の保存成功後だけ active を消去して、次回起動まで次 waiting entry を処理しない。queue 全体が JSON として破損している場合は [`docs/details/runner.md` 詳細本文責務 §12](runner.md#12-設定値runner) の `.build_state` 破損処理に従う。

**cooldown / force build 判定契約：**

runner は起動ごとに `.server_config.build_cooldown_seconds` と `.server_config.force_build_interval_hours` を 1 回読み、当該起動中の cooldown / force build 判定に使用する。同一 runner 起動中に `.server_config` を再読込して判定値を変更してはならない。

判定順は pending transfer retry、notify pending retry、circuit breaker、active / waiting queue 取得、cooldown、SHA decision の順とする。circuit open 中は queue entry を waiting から active へ移さず、既存 active も実行しない。manual queue entry の `payload.force=true` は cooldown を無視する。webhook queue entry は cooldown を適用する。`force_build_interval_hours` は SHA 一致時だけ評価し、SHA 不一致時は常に通常 build とする。

| 条件 | 結果 |
|------|------|
| `.build_state.last_finished_at=null` | cooldown は適用しない。 |
| `.server_config.build_cooldown_seconds=0` | cooldown は無効。 |
| `now - last_finished_at < build_cooldown_seconds` | `skipped_cooldown`。GitHub API、pipeline、deploy、snapshot は実行しない。 |
| SHA 一致かつ `.server_config.force_build_interval_hours=0` | `skipped_no_change`。 |
| SHA 一致かつ `force_build_interval_hours>0` かつ直近成功 build から指定時間未満 | `skipped_no_change`。 |
| SHA 一致かつ `force_build_interval_hours>0` かつ直近成功 build から指定時間以上 | `force_interval` として build を実行する。 |

force interval の直近成功 build は `.build_history` のうち同じ `branch` と `target_file` で status が `success` または `success_deploy_pending` の行を `finished_at` 降順、同時刻は `id` 降順に並べた先頭とする。`.build_history` が存在しない、または該当行がない場合は force interval 条件成立として build する。

**runner finalizer 固定契約：**

runner は lock 取得後、正常終了、失敗終了、panic 相当の recover、context timeout のいずれでも finalizer を実行する。finalizer は次の順序に固定する。

1. 未保存の `.build_logs/{id}.json` がある場合は、取得済みで schema を満たす値だけを使って最終形を保存する。未取得値を推測して補完してはならない。
2. `.build_history` へ追記対象の build がある場合は 1 行だけ追記する。同じ `id` が既に存在する場合は追記せず、ERROR ログ `BUILD_HISTORY_DUPLICATE: id={id}` を出す。
3. `.build_status.json` を最終状態へ更新する。
4. queue entry 実行で build log と history の最終結果が確定した場合は、実行対象 id と `active_queue_entry.id` の一致を確認し、`active_queue_entry=null`、`running=false`、`current_build_id=null`、`last_finished_at={now}` を同じ atomic write で保存する。
5. queue entry 実行で手順 1〜3 の必須保存が確定しなかった場合は、`active_queue_entry` を保持する。失敗対象が `.build_state` の読取、path 検証、または atomic write ではなく、実行対象 id と `active_queue_entry.id` が一致する場合に限り、`active_queue_entry` と `queued` を保持したまま `running=false`、`current_build_id=null`、`last_finished_at={now}` を保存する atomic write を 1 回だけ試行する。この試行が失敗した場合は再試行せず `BUILD_STATE_FINALIZE_FAILED` を ERROR で記録し、終了コードを最低 `1` とする。失敗対象が `.build_state` の読取、path 検証、または atomic write である場合、および active id 不一致時は追加 write を行わず、状態を変更せず終了コードを最低 `1` とする。
6. `.build_lock` を削除する。
7. 通知対象 event がある場合は通知または `.notify_pending` 追記を行う。

finalizer 中に複数失敗が発生した場合、終了コードは最も重い値を採用する。`.build_state.running=false` の保存失敗は終了コード `1` 固定とし、lock 削除だけ成功しても正常終了扱いにしない。結果未確定時の `active_queue_entry` 消去、結果確定時の active id 不一致、waiting queue の変更は禁止する。lock 削除失敗は WARN とし、他失敗がなければ終了コードを変更しない。

**build log / history 書き込み契約：**

`.build_logs/{id}.json` は 1 build id につき 1 file とする。build 開始時に [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の途中保存 schema で atomic write し、同じ runner 実行が保持する build id と一致する場合だけ、最終 schema で atomic replace する。build id 採番時点で同名 file が先に存在する場合は上書きせず、次の suffix 付き build id を採番し直す。他実行が作成した file、既に最終状態の file、build id 不一致 file を置換してはならない。`.build_history` は JSON Lines とし、最終 build log の保存成功後に 1 行だけ追記する。追記前に既存 file の末尾が LF で終わることを確認し、LF がない場合は 1 個だけ LF を追加してから新規行を追記する。

`.build_history` の保存 key、nullable 条件、詳細結果値は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.build_history` JSON Lines schema を正とする。build log を伴う行の `status` は同じ id の `target_status` と一致させる。approval event と chain dependency skip は同じ正本 schema で許可された history-only 結果を保存し、対応する build log を作成しない。JSON Lines の壊れた既存行は読み取り対象から除外するが、追記時に既存 file 全体を書き換えてはならない。

**サーキットブレーカー更新契約：**

`.build_circuit_state.open=true` の場合、runner は pending transfer retry と notify pending retry を完了した後、GitHub API、queue 取得、build id 採番、pipeline、deploy、snapshot を実行せず、`.build_status.json` に `status="skipped"`、`last_target_status="circuit_open"`、`running=false`、`current_build_id=null` を保存して終了する。active / waiting queue は変更しない。`.build_state.running=true` または `current_build_id != null` が残っている場合だけ、process lock 保持中に `running=false`、`current_build_id=null` を同じ atomic write で保存し、`active_queue_entry`、`queued`、`last_started_at`、`last_finished_at` を保持する。状態保存に成功した場合の終了コードは `0`、失敗は `1` とする。circuit open を queue entry の成功、失敗、完了として log / history へ記録してはならず、active entry を消去してはならない。

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

runner が読み込む JSON object / JSON array の状態ファイルが破損している場合は、[`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の破損時の扱いに従う。JSON Lines は壊れた行だけを無視し、ファイル全体を破棄してはならない。破損退避ファイル名は `{original}.corrupt.{YYYYMMDDHHMMSS}.bak` とする。

```
runner 起動（systemd timer、または API の `systemctl start --no-block adlaire-ci.service` から呼び出し）
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
    │   ├─ §27.10 の対象 7 ファイルを固定順序で検証
    │   ├─ 復旧可能な破損 → corrupt backup、初期化または default fallback、[`docs/details/runner.md` 詳細本文責務 §27.10](runner.md#sec-27-10) の通知条件を満たす場合は config_corrupt 通知
    │   ├─ unknown key → corrupt backup、ファイル別の初期化または default fallback
    │   └─ permission error / IO error → .build_state.running=true にせず終了
    │
    ├─ [ペンディングキュー再試行] PENDING_FILE が存在する場合
    │   └─ ペンディングエントリごとに SSH 転送を再試行
    │       ├─ 成功 → エントリを PENDING_FILE から削除
    │       └─ 失敗 → ERROR ログ、エントリを保持（次回起動時に再試行）
    │
    ├─ [Webhook 通知ペンディング再試行] .notify_pending が存在する場合
    │   └─ ペンディングエントリごとに Webhook 送信を再試行
    │       ├─ 成功（HTTP 2xx）→ INFO ログ、エントリを .notify_pending から削除
    │       └─ 失敗（HTTP エラー・接続失敗）→ ERROR ログ、エントリを保持（次回起動時に再試行）
    │
    ├─ [サーキットブレーカー開始判定] .build_circuit_state.open = true の場合
    │   ├─ .build_status.json に skipped / circuit_open / running=false を保存
    │   ├─ runtime flag が残る場合だけ .build_state.running=false / current_build_id=null を保存
    │   ├─ active_queue_entry と waiting queue は変更しない
    │   └─ GitHub API、queue 取得、build id 採番、pipeline、deploy、snapshot を実行せず終了
    │
    ├─ [queue entry 取得] circuit closed の場合
    │   ├─ active_queue_entry あり → 同じ entry を再実行対象にする
    │   ├─ active なし / waiting あり → priority / created_seq / id 順の 1 件を active へ atomic move
    │   └─ active / waiting なし → 通常 polling 判定へ進む
    │
    ├─ [クールダウンチェック] build_cooldown_seconds > 0 の場合
    │   └─ .build_state.last_finished_at から build_cooldown_seconds 秒未満
    │       → INFO ログ（`COOLDOWN: skip, last_build N秒前`）、正常終了
    │
    ├─ BRANCH_TARGETS の各エントリを順次処理：
    │   │
    │   ├─ Step 1: Git Trees API
    │   │   GET /repos/{RunnerConfig.RepositoryOwner}/{RunnerConfig.RepositoryName}/git/trees/{branch}?recursive=1
    │   │   → target_file の blob SHA を取得
    │   │   └─ API 失敗時：API_RETRY_MAX 回まで指数バックオフ（API_RETRY_BASE_SECONDS × 2^n 秒）で再試行
    │   │        ├─ Trees API 全試行失敗時：ERROR ログ、このエントリを failure(api_error) として記録し、SHA を更新せず次エントリへ進む
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
    │   │   GET /repos/{RunnerConfig.RepositoryOwner}/{RunnerConfig.RepositoryName}/git/blobs/{sha}
    │   │   → Base64 デコード → src パスへ書き出し
    │   │   └─ API 失敗時：API_RETRY_MAX 回まで指数バックオフで再試行
    │   │        ├─ Blobs API 全試行失敗時：ERROR ログ、このエントリを failure(api_error) として記録し、SHA を更新せず次エントリへ進む
    │   │        └─ レスポンスヘッダー GitHub-Authentication-Token-Expiration が存在する場合は Step 1 と同様に PAT 有効期限チェックを行う
    │   │
    │   ├─ [コミット情報取得] ビルドトリガーとなったコミット情報を取得し、ビルドログへ記録する
    │   │   GET /repos/{RunnerConfig.RepositoryOwner}/{RunnerConfig.RepositoryName}/commits?path={BRANCH_TARGET.src}&sha={branch}&per_page=1
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
    │   │   └─ `/usr/local/bin/adlaire-ci-build --version` が終了コード 0 で、stdout に `adlaire-ci-build` と `v3` を含むこと
    │   │
    │   ├─ pipeline.sh 実行（bash {src の親ディレクトリ}/.ci/pipeline.sh）
    │   │   ├─ 成功（exit 0）：INFO ログ
    │   │   │   └─ [Webhook 通知送信] on: ["success"] 設定時
    │   │   │       → .notify_config の Webhook 宛先へ POST（ペイロード: event="success", branch, build_id 等）
    │   │   │       → 送信失敗（HTTP エラー・接続失敗・タイムアウト）の場合：
    │   │   │           ERROR ログ（`NOTIFY_FAIL: url={url} status={code}`）
    │   │   │           .notify_pending へ statefile 正本の NotificationPending object を追記
    │   │   └─ 失敗（exit ≠ 0）：ERROR ログ、sha_file 更新せず、このエントリを failure(build_failed) として記録し、次エントリへ進む
    │   │       └─ [Webhook 通知送信] on: ["failure"] 設定時
    │   │           → .notify_config の Webhook 宛先へ POST（ペイロード: event="failure", branch, build_id 等）
    │   │           → 送信失敗の場合：ERROR ログ、.notify_pending へ statefile 正本の NotificationPending object を追記
    │   │
    │   └─ sha_file を新 SHA で更新（`{"sha": "<new_sha>"}` を JSON 書き込み）
    │        └─ SSH サイト転送（deploy_targets リストの各エントリへ転送）
    │             ├─ [転送後整合性検証] 引用済み remote path を `sha256sum --zero --` で検証
    │             │   ├─ 全ファイルのローカル sha256 と一致 → 転送成功
    │             │   └─ 不一致またはコマンド失敗 → ERROR ログ、ペンディングキューへ再投入（[`docs/details/runner.md` 詳細本文責務 §14a](runner.md#14a-ssh-サイト転送)）
    │             └─ 整合性検証成功後 → §14b の条件判定後に archive owner の snapshot save を呼び出す
    │
    │        [出力サイズチェック] OUTPUT_SIZE_WARN_MB > 0 の場合
    │        出力サイト配下の通常ファイル合計サイズを取得し、閾値と比較：
    │            size_mb = total_site_bytes / (1024 * 1024)
    │            size_mb > OUTPUT_SIZE_WARN_MB の場合：
    │            → WARN ログ（`OUTPUT_SIZE_WARN: size={size_mb:.1f}MB threshold={OUTPUT_SIZE_WARN_MB}MB`）
    │            → ビルドログの size_warn フィールドを true に設定（[`docs/details/builder.md` 詳細本文責務 §8](builder.md#8-実行方法)）
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
        └─ 送信失敗（HTTP エラー・接続失敗）→ ERROR ログ、retry 対象の場合だけ statefile 正本の NotificationPending object を追記
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
- `pipeline.sh` は `ADLAIRE_CI_OUT` を出力サイトディレクトリとし、`index.html`、`assets/style.css`、`assets/app.js`、`assets/search-index.json` を必ず生成しなければならない
- Markdown ディレクトリ入力時の `pages/*.html` を含む出力 file set は、[`docs/details/builder.md` 詳細本文責務 §2a](builder.md#2a-入力収集出力パス決定) と [§5](builder.md#5-静的-web-サイト出力構造) の出力契約と一致しなければならない

**最小固定内容：**
```bash
#!/usr/bin/env bash
set -euo pipefail
/usr/local/bin/adlaire-ci-build --src "$ADLAIRE_CI_SRC" --out "$ADLAIRE_CI_OUT"
```

ビルド実行コマンドは `pipeline.sh` 内に直接記述する（`runner` は参照しない）。`adlaire-ci-build` は [`components/builder.go`](../../components/builder.go) から生成した Go 版バイナリである。

**runner からの実行契約：**

| 項目 | 仕様 |
|------|------|
| command | `bash {srcParent}/.ci/pipeline.sh` |
| working directory | `{srcParent}` |
| timeout | `.server_config.build_timeout_seconds` が存在すればその値、なければ `300` 秒。 |
| stdout/stderr | それぞれ最大 1 MiB までメモリに保持し、超過分は末尾 1 MiB を保存する。超過時は `stdout_truncated` / `stderr_truncated` を `true` にする。 |
| 環境変数 | 親 process の環境を引き継ぎ、[`docs/details/runner.md` 詳細本文責務 §14](runner.md#14-pipelinesh) の固定表の値で上書きする。 |

| 環境変数 | 値 |
|----------|----|
| `ADLAIRE_CI_SRC` | `BranchTarget.Src` |
| `ADLAIRE_CI_OUT` | 出力サイトディレクトリである `BranchTarget.Out` |
| `ADLAIRE_CI_BRANCH` | `BranchTarget.Branch` |
| `ADLAIRE_CI_BUILD_ID` | build id |
| `ADLAIRE_CI_TARGET_FILE` | `BranchTarget.TargetFile` |
| `ADLAIRE_CI_STATE_DIR` | `RunnerConfig.StateDir` |

`pipeline.sh` が timeout した場合、runner は process group を終了し、`target_status="failure_timeout"`、`error="pipeline timeout"` として記録する。timeout 時も stdout/stderr の取得済み内容は `.build_logs/{id}.json` に保存する。

**GitHub API 固定契約：**

| 項目 | 仕様 |
|------|------|
| repository identity | runner 起動時に `readRepoConfig()` から取得し、`RunnerConfig.RepositoryOwner` / `RepositoryName` に正規化した 1 組を当該 process の全 GitHub API request で使用する。実行中に `.repo_config` を再読込みしない。 |
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

`POST /api/repo-config` の成功後に起動する次回 runner process から新しい repository identity を使用する。すでに起動済みの runner、実行中 build、active queue entry の repository identity を API 更新によって途中変更しない。`.repo_config` の `branch` / `target_file` を参照する互換処理は実装せず、branch と監視対象は `.branch_config` だけから取得する。

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
| exit `1`〜`125` | 実際の終了コード | `failure_build` | `pipeline failed` | [`docs/details/runner.md` 詳細本文責務 §27.3](runner.md#sec-27-3) の retry 対象。 |
| exit `126` | `126` | `failure_precheck` | `precheck failed` | retry しない。 |
| exit `127` | `127` | `failure_precheck` | `precheck failed` | retry しない。 |
| signal 終了 | `128 + signal` | `failure_build` | `pipeline failed` | retry 対象。 |
| timeout | `null` | `failure_timeout` | `pipeline timeout` | retry 対象。 |
| stdout 上限超過 | 実際の終了コード | exit code に従う | exit code に従う | exit code に従う。 |
| stderr 上限超過 | 実際の終了コード | exit code に従う | exit code に従う | exit code に従う。 |

runner は stdout / stderr の CRLF を LF に正規化して保存する。NUL byte は `\u0000` 文字列へ置換する。保存する stdout / stderr は UTF-8 不正 byte を `�` に置換する。secret mask は保存前に適用し、PAT、Webhook Secret、SMTP password、API token、session token、TOTP secret に一致する値を `***` に置換する。

---

## 14a. SSH サイト転送

[`docs/details/runner.md` 詳細本文責務 §14a](runner.md#14a-ssh-サイト転送) は、runner owner の SSH 転送詳細本文責務である。

`runner` は `pipeline.sh` 成功後に、出力サイトディレクトリ配下の全ファイルを SSH 経由で静的コンテンツ配信サーバーへ転送する。scp・rsync は使用しない。local process は `ssh` バイナリを `exec.CommandContext` で直接起動し、local の `/bin/sh -c` を使わない。OpenSSH の remote command は remote login shell に 1 文字列として渡るため、動的値は [remote 引数引用固定契約](#runner-remote-argument-quoting-contract) で必ず引用する。

**設定値：**

`BRANCH_TARGETS` 各エントリの `deploy_targets` リスト内で管理する。

| フィールド | 説明 | 例 |
|-----------|------|-----|
| `host` | 配信サーバーのホスト名 / IP | `"192.0.2.1"` |
| `user` | SSH 接続ユーザー | `"deploy"` |
| `dest_dir` | 配信サーバー上の転送先ディレクトリ | `"/var/www/html"` |

転送対象は `BRANCH_TARGETS[n]["out"]` のディレクトリ配下にある通常ファイルすべてとする。`deploy_targets` に複数エントリを定義した場合は全ての転送先へ順次転送する。

**差分検出：**

転送前にリモートサーバーで対象ファイルごとの SHA256 ハッシュを取得し、ローカルファイルのハッシュと比較する。

remote command は `sha256sum --zero -- {quoted_remote_path}` とする。`{quoted_remote_path}` は [remote 引数引用固定契約](#runner-remote-argument-quoting-contract) の `quoteRemoteArg` を適用した 1 引数である。

- ハッシュが一致 → 当該ファイルをスキップ（`SKIP` ログを記録）
- ハッシュが不一致、またはリモートにファイルが存在しない → 当該ファイルを転送する

relative path は `out` からの相対 path とし、`filepath.Rel` 後に `/` 区切りへ変換して保存する。空文字、`.`、`..` を含む path、先頭 `/`、NUL byte、制御文字を含む path は転送対象から除外し、`failure_precheck` とする。symbolic link、directory、device file、FIFO は転送しない。symbolic link を検出した場合は ERROR `DEPLOY_UNSUPPORTED_FILE: path={path}` を出し、当該 target を `failure_precheck` とする。

**転送：**

stdin パイプ経由で SSH 転送する。

runner は local file を開き、SSH process の stdin へ `io.Copy` で送る。リモート側 stdout は使用せず破棄し、stderr は失敗理由として `.build_logs/{id}.json.error` と ERROR ログへ記録する。

<a id="runner-remote-argument-quoting-contract"></a>
remote 引数引用固定契約は次のとおりとする。`quoteRemoteArg(value)` は `"'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"` と完全一致する文字列を返す。remote command に含めるすべての動的値にこの変換を個別適用し、未引用値、二重引用、backslash escape、文字列置換後の再解釈を禁止する。schema 検証を通過していない host、user、path、argument で remote command を生成してはならない。

remote path は Go `path.Join(dest_dir, relativePath)` で生成する。生成後に、`dest_dir="/"` の場合は remote path が `/` で開始すること、それ以外は remote path が `dest_dir + "/"` で開始することを再確認し、不一致は `failure_precheck` とする。directory 作成は `mkdir -p -- {quoted_remote_dir}`、file 転送は `tee -- {quoted_remote_path}` の 2 回の remote command に分け、local 起動はどちらも `exec.CommandContext(ctx, "ssh", "--", user+"@"+host, remoteCommand)` とする。`--` は local `ssh` option の終端として必須とし、接続先を option として解釈させない。`remoteCommand` の固定 command 名、固定 option、区切り空白以外を動的値から生成してはならない。remote host は POSIX-compatible login shell、`mkdir`、`tee`、GNU `sha256sum` の `--zero` を提供しなければならず、command 不在または非対応は転送失敗とする。

転送は file 単位で行い、1 file の転送 timeout は 60 秒とする。timeout 時は SSH process group を終了し、当該 deploy target を pending とする。1 deploy target 内で 1 file でも転送または検証に失敗した場合、その deploy target 全体を pending とし、snapshot は作成しない。

**ペンディングキュー：**

転送失敗時（接続エラー、認証失敗、timeout、remote checksum 不一致）は `PENDING_FILE` へ PendingTransfer object を保存する。通常 build 由来は `trigger="deploy"`、`source_kind="output"`、`out=<当該 branch target の絶対 path>`、`rollback_from=null`、`snapshot_id=null` とする。entry の key、型、初期 `retry_count`、転送先重複判定、新しい build による置換は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.pending_transfers` schema を参照する。

- `runner` 起動時の `BRANCH_TARGETS` 処理前に `PENDING_FILE` を読み込み、array 順に再試行する。
- `source_kind="output"` は `out` の現在の成果物 manifest SHA-256 を [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) で再計算し、`output_sha256` と一致する場合だけ転送する。不一致は SSH を実行せず、`retry_count` を 1 増やし、`failed_at` を更新し、`last_error="pending_output_changed"` で entry を保持する。
- `source_kind="snapshot"` は runner owner が entry を archive owner へ渡し、archive owner が保存済み snapshot の検証、`meta.json.output_sha256` と entry の `output_sha256` 一致確認、新規 temporary directory への再展開、当該 1 target への転送、temporary cleanup を [`docs/details/archive.md` 詳細本文責務 §27.15](archive.md#sec-27-15) に従って実行する。runner owner は `out` を推測または復元しない。
- snapshot 不在は `last_error="pending_snapshot_missing"`、検証失敗は `last_error="pending_snapshot_corrupt"`、checksum 不一致は `last_error="pending_snapshot_changed"` とし、SSH を実行せず `retry_count` と `failed_at` を更新して entry を保持する。
- 再試行成功時は entry を削除する。再試行失敗時は `retry_count` を 1 増やし、`failed_at` と `last_error` を更新する。
- SSH 転送失敗は `.notify_config` で `deploy_failure` が有効な channel へ通知する。通知失敗は PendingTransfer object の保存成否を変更しない。
- `branch_idx`、`deploy_idx`、`attempt` を読み替えたり保存したりしない。設定配列の index で pending target を復元しない。

**転送後整合性検証：**

SSH 転送完了後に、リモートファイルの SHA-256 チェックサムをローカルのものと照合する。

**検証コマンド：**

local 起動は `exec.CommandContext(ctx, "ssh", "--", user+"@"+host, "sha256sum --zero -- "+quoteRemoteArg(remotePath))` とする。remote stdout は `64 文字 lowercase hex + "  " + remotePath の UTF-8 byte + NUL` と完全一致する 1 record だけを成功とし、余分な byte、大文字 hex、path 不一致、NUL 不足は検証失敗とする。hash を local `crypto/sha256` の lowercase hex digest と比較する。

| 項目 | 仕様 |
|---|---|
| タイムアウト | 30 秒（SSH 転送タイムアウトとは独立） |
| 検証失敗時 | ERROR ログ＋ペンディングキューへ再投入。スナップショット保存はしない |
| ログフィールド | `transfer_verified: false`（`.build_logs/{id}.json` に記録） |
| 正常時 | `transfer_verified: true`（`.build_logs/{id}.json` に記録） |

remote `sha256sum` の終了コード非 0、空出力、stderr のみの出力、または `64 文字 lowercase hex + "  " + remotePath の UTF-8 byte + NUL` の 1 record 完全一致に失敗した出力は検証失敗とする。local checksum は転送直前に読んだ file 内容ではなく、転送後に local file を再読込して計算する。

**デプロイログ：**

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

[`docs/details/runner.md` 詳細本文責務 §14b](runner.md#14b-スナップショット管理) は、`runner` owner の snapshot 作成トリガー判定、`archive` owner への入力引渡し、build 結果への反映だけを定義する。snapshot の保存形式、archive 作成、検証、世代削除、一覧、download、delete、展開、rollback 転送実体は [`docs/details/archive.md` 詳細本文責務 §27.15](archive.md#sec-27-15) を正本とし、runner owner 側で重複定義しない。

**snapshot 呼出条件固定契約：**

| 条件 | runner の処理 |
|------|-----------------|
| `.server_config.snapshots_keep=0` | `archive` owner を呼び出さず、`snapshot_id=null` とする。`.snapshots/` の存在有無は判定に使わない。 |
| pipeline 失敗、SHA 書込失敗、deploy pending / failure | `archive` owner を呼び出さず、`snapshot_id=null` とする。snapshot 保存失敗としては記録しない。 |
| pipeline 成功、SHA 書込成功、deploy target なし | deploy phase の成功 no-op とし、`archive` owner を 1 回呼び出す。 |
| pipeline 成功、SHA 書込成功、全 deploy target の転送と検証に成功 | `archive` owner を 1 回呼び出す。 |

**snapshot save 呼出契約：**

| 入力 | 値 |
|------|----|
| `build_id` | 当該 build log / history と同じ build id。 |
| `site_dir` | 当該 branch target の確定済み output site 絶対 path。 |
| `saved_at` | snapshot save 呼出直前の fake clock 対応 UTC ISO 8601 秒精度。 |
| `output_sha256` | 当該 output site に対する build history と同じ manifest SHA-256。算出不能な契約済みケースだけ `null`。 |
| `snapshots_keep` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) で検証済みの `.server_config.snapshots_keep`。 |

`runner` は archive 内部の file 作成、tar/gzip 処理、tmp rename、prune を直接実行しない。呼出結果は次の固定値とし、最終 build log と同じ id の history に同じ `snapshot_id` を保存する。

| archive 結果 | `snapshot_id` | warning | build / deploy への影響 |
|----------------|---------------|---------|-------------------------|
| `saved` | `build_id` | archive owner が返した `SNAPSHOT_PRUNE_FAILED`、`SNAPSHOT_DIRECTORY_SYNC_FAILED`、`SNAPSHOT_TMP_CLEANUP_FAILED` を返却順のまま追加する。 | 成功を維持する。 |
| `exists_valid` | `build_id` | `SNAPSHOT_EXISTS` | 冪等成功とし、成功を維持する。 |
| `failed` | `null` | `SNAPSHOT_SAVE_FAILED` | build 成功と deploy 成功を取り消さず、runner 終了コードを変更しない。 |

`failed` は `.pending_transfers` 追加理由にしない。publish 後の directory sync、tmp cleanup、または prune だけが失敗した場合は archive 結果を `saved` のままとし、archive owner が返した固定 warning を build log へ保存する。snapshot root の作成、archive 作成、事前検証、prune は [`docs/details/archive.md` 詳細本文責務 §27.15](archive.md#sec-27-15) のみが定義し、runner は `.snapshots/` を初期化または直接変更しない。

**rollback coordinator 固定契約：**

API 経由の rollback で runner owner は、rollback の排他、build id 採番、開始状態、非同期実行、最終 log / history / pending / status / state、lock 解放を担当する。archive owner は [`docs/details/archive.md` 詳細本文責務 §27.15](archive.md#sec-27-15) に従い、保存済み archive の事前検証、temporary directory への展開、deploy、cleanup と結果返却だけを担当する。

rollback coordinator は snapshot id と trigger actor を受け取り、元 history の `branch` と `target_file` に完全一致する現在の正規化済み branch target を 1 件選択する。元 history、snapshot、一致 target、または deploy target が存在しない場合は rollback を開始しない。現在の deploy target を使い、元 history に保存された host、user、dest_dir を復元しない。

**rollback accepted 確定順：**

1. `.build_lock` を [`docs/details/runner.md` 詳細本文責務 lock ファイル契約](#lock-ファイル契約) で取得する。有効な running lock、形式不正、PID 判定不能は conflict とし、状態を変更しない。
2. 新規 build id と `started_at` を確定する。
3. archive owner に `snapshot_id`、`new_build_id`、現在の deploy target 一覧を渡し、archive 検証と temporary directory への展開を完了させる。不在、破損、展開失敗では persistent state を書き込まず、archive owner の temporary cleanup と runner owner の所有確認付き lock 解放をそれぞれ 1 回だけ実行する。
4. `.build_status.json` に `status="running"`、`running=true`、`current_build_id=<new_build_id>`、`last_trigger="rollback"`、`last_started_at` を保存する。
5. `.build_state` に `running=true`、`current_build_id=<new_build_id>`、`last_started_at` を保存し、`active_queue_entry` と `queued` を変更しない。
6. `.build_logs/{new_build_id}.json` に `status="running"`、`target_status=null`、`trigger="rollback"`、`rollback_from=<snapshot_id>`、`snapshot_id=<snapshot_id>`、元 history の `branch` / `target_file`、`output_sha256`、`trigger_actor` を保存する。
7. worker をまだ開始せず、api owner へプロセス内だけで有効な一回性 `prepared` handle と `new_build_id` を返す。手順 6 まで完了していない場合は `prepared` を返さない。
8. api owner が同じ request id、trigger actor、`target_type="build"`、`target_id=<new_build_id>` の `build_trigger` audit を追記する。
9. audit 成功後にだけ、api owner が `StartPreparedRollback(prepared)` を 1 回呼び出す。runner owner は handle を即時無効化し、request context から切り離した rollback worker を開始し、`accepted`、`new_build_id` を返す。

手順 4〜6 の書込失敗では `prepared` を返さず、audit、worker、SSH、deploy を実行しない。`build_trigger` audit 追記失敗時、api owner は worker を開始せず `AbortPreparedRollback(prepared, "rollback audit failed")` を 1 回呼び出す。runner owner は呼出開始時に handle を無効化し、以降の start、abort、再利用を拒否する。API は audit 失敗を `500` で返し、補償成否で元の HTTP 結果を変更しない。

**rollback worker 開始前補償固定契約：**

補償 reason は、手順 4〜6 の失敗では `rollback start state write failed`、`AbortPreparedRollback` では引数の `rollback audit failed` とする。runner owner は次の順に進め、各手順を最大 1 回だけ実行する。一つの手順が失敗しても後続手順を省略せず、失敗した手順を再試行しない。

1. `.build_logs/{new_build_id}.json` を read adapter で 1 回読む。`id=new_build_id`、`status="running"`、`trigger="rollback"` に完全一致する schema-valid log が存在する場合だけ、他 field を維持したまま `status="failure"`、`target_status="failure_state_write"`、`finished_at=<fake clock 対応の補償開始時刻>`、`duration_seconds=max(0, finished_at-started_at)`、`error=<reason>` で atomic replace する。`failure_category` と `failure_evidence` は元の失敗事実を [失敗原因の自動分類](#sec-27-36) に適用した値とし、証跡がない場合は `failure_category="unknown"`、`failure_evidence=[]` とする。log 不在、破損、id / status / trigger 不一致では log を作成、修復、上書きしない。置換失敗は `ROLLBACK_PREWORKER_LOG_FINALIZE_FAILED` を ERROR で記録する。
2. 手順 1 の final log 置換が成功した場合だけ、同じ id の history を 1 行追記する。`status="failure_state_write"`、`trigger="rollback"`、`error=<reason>` とし、その他は final log の要約値を使う。追記失敗は `ROLLBACK_PREWORKER_HISTORY_WRITE_FAILED` を ERROR で記録し、同じ id を再追記しない。
3. `.build_status.json` は read adapter で得た schema-valid object の `running=true`、`current_build_id=new_build_id`、`last_trigger="rollback"` がすべて一致する場合だけ、`status="failure"`、`running=false`、`current_build_id=null`、`last_target_status="failure_state_write"`、`last_finished_at=<補償開始時刻>`、`last_duration_seconds=max(0, last_finished_at-last_started_at)`、`last_error=<reason>` へ atomic write する。手順 1 の final log 置換成功時だけ `last_build_id=new_build_id` とし、final log が確定していない場合は既存 `last_build_id` を維持する。不在、破損、不一致では作成または上書きしない。write 失敗は `ROLLBACK_PREWORKER_STATUS_FINALIZE_FAILED` を ERROR で記録する。
4. `.build_state` は read adapter で得た schema-valid object の `running=true`、`current_build_id=new_build_id` が一致する場合だけ、`running=false`、`current_build_id=null`、`last_finished_at=<補償開始時刻>` へ atomic write する。`active_queue_entry`、`queued`、その他の field は維持する。不在、破損、不一致では作成または上書きしない。write 失敗は `ROLLBACK_PREWORKER_STATE_FINALIZE_FAILED` を ERROR で記録する。
5. archive owner が [`docs/details/archive.md` 詳細本文責務 §27.15](archive.md#sec-27-15) に従い rollback temporary directory を 1 回 cleanup する。失敗は `SNAPSHOT_ROLLBACK_TMP_CLEANUP_FAILED` warning とし、補償済み state を巻き戻さない。
6. runner owner が [lock ファイル契約](#lock-ファイル契約) に従い、所有確認付き lock 解放を 1 回実行する。

この補償は [`docs/details/statefile.md` 詳細本文責務 状態ファイル更新手順](statefile.md#statefile-update-procedure) の明示的 lifecycle 補償であり、失敗した開始 payload の retry ではない。手順 4〜6 の失敗は元の coordinator error、audit 失敗は元の audit error を primary error とし、補償失敗で置き換えない。

**rollback worker 書込順：**

1. archive owner から `success`、`failure_build`、または `success_deploy_pending`、deploy target 別結果、pending 候補、warning を受け取る。
2. `success_deploy_pending` の場合だけ、pending 候補を `trigger="rollback"`、`source_kind="snapshot"`、`out=null`、`rollback_from=snapshot_id`、`output_sha256=<meta.json.output_sha256>` とする [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の PendingTransfer object として `.pending_transfers` へ保存する。snapshot の `output_sha256` が `null` の場合は再送可能な pending を作成せず `failure_build` とする。保存失敗時は rollback 結果を `failure_state_write` へ変更し、pending 成功を記録しない。
3. 同じ build id の running log を最終 schema で atomic replace する。`target_status` は確定済み結果、`rollback_from` と `snapshot_id` はどちらも元 snapshot id、`finished_at`、`duration_seconds`、deploy 結果、warning を含める。
4. `.build_history` に同じ id、trigger、target status、rollback_from、snapshot_id、output_sha256 を 1 行追記する。
5. `.build_status.json` を最終状態へ更新する。
6. `.build_state` の `running=false`、`current_build_id=null`、`last_finished_at` を保存し、`active_queue_entry` と `queued` は変更しない。
7. `.build_lock` を解放する。

final log 書込失敗時は history を追記せず、running log を維持したまま status finalizer、state finalizer、所有確認付き lock 解放をそれぞれ 1 回だけ実行する。history 追記失敗時は final log を維持し、同じ id を自動再追記せず、status finalizer、state finalizer、所有確認付き lock 解放をそれぞれ 1 回だけ実行する。status または state finalizer が失敗しても後続手順を実行し、失敗した write を再試行しない。accepted 後の panic、context timeout、API server の graceful shutdown による cancel は `failure_state_write` と `ROLLBACK_INTERRUPTED` を最終 log / history に記録し、request context の cancel だけで worker を中断しない。

process crash 後、次の runner 起動が stale lock、`.build_state.running=true`、`current_build_id`、および同じ id の `status="running"` かつ `trigger="rollback"` の log をすべて確認できた場合に限り、新しい build を開始する前に当該 log を `failure_state_write`、`error="rollback interrupted"` で最終化し、history、status、state を rollback worker 書込順で各 1 回だけ書き込み、archive owner が `.snapshots/.rollback.{current_build_id}.tmp` を 1 回 cleanup し、runner owner が所有確認付き lock 解放を 1 回実行する。各失敗で後続手順を省略せず、失敗した write は再試行しない。これは rollback coordinator が作成した同一 build id の途中状態だけを失敗として閉じる固定例外であり、他の stale 状態を修復する一般規則としてはならない。必要な値の不足または破損がある場合は推測修復せず、通常 build を開始しない。

---

## 15. ログ

[`docs/details/runner.md` 詳細本文責務 §15](runner.md#15-ログ) は、`runner` の stdout ログと構造化ビルドログを定義する。

**stdout ログ：**

| レベル | 出力条件 |
|--------|---------|
| `INFO` | 起動、変更なしスキップ、ビルド開始・完了、SHA 更新 |
| `WARNING` | — |
| `ERROR` | トークン読み込み失敗、API 失敗、ビルド失敗 |
| `DEBUG` | API レスポンス詳細等（`LOG_LEVEL = "DEBUG"` 時のみ） |

stdout は Go 標準ライブラリ `log/slog` で出力し、systemd が journald に転送する。独自 logger 実装を使用してはならない。

**構造化ビルドログ：**

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

これらのログ項目を実装対象に含める時点で、[`docs/details/runner.md` 詳細本文責務 §10a](runner.md#10a-ci-ランナー-実装対象) の実装対象、[`docs/details/runner.md` 詳細本文責務 §12](runner.md#12-設定値runner) の設定値、[`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) の処理フロー、[`docs/details/api.md` 詳細本文責務 §22](api.md#22-バックエンド-api-仕様) の API レスポンス仕様と整合させる。

**`.build_logs/{id}.json` schema 参照：**

`.build_logs/{id}.json` の保存 key、型、必須条件、Report object、Attempt object、CommitStatus object、BuildMeta object は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.build_logs/{id}.json` schema を参照する。

runner owner component は、build log の生成タイミング、stdout / stderr 取り込み、`[REPORT]` 変換、WARN 取り込み、secret mask、最終状態保存、書き込み失敗時の後続停止だけを担当する。schema key の追加、削除、型変更、未知 key 保存は [`docs/details/runner.md`](runner.md) 詳細本文責務で行ってはならない。

**ビルドログ最終形契約：**

| 状況 | `finished_at` | `duration_seconds` | `pipeline` | `report` | `deploy` | `snapshot_id` |
|------|---------------|--------------------|------------|----------|----------|---------------|
| GitHub API 失敗 | 保存時刻 | 0 以上 | `exit_code:null`, stdout/stderr 空 | `null` | `[]` | `null` |
| blob decode 失敗 | 保存時刻 | 0 以上 | `exit_code:null`, stdout/stderr 空 | `null` | `[]` | `null` |
| precheck 失敗 | 保存時刻 | 0 以上 | `exit_code:null`, stdout/stderr 空 | `null` | `[]` | `null` |
| pipeline timeout | timeout 検出時刻 | 0 以上 | 取得済み stdout/stderr、`exit_code:null` | parse できた場合のみ object | `[]` | `null` |
| pipeline 非 0 | process 終了時刻 | 0 以上 | 実 exit code と取得済み stdout/stderr | parse できた場合のみ object | `[]` | `null` |
| pipeline 成功 / deploy なし / snapshot 無効 | snapshot 非呼出確定時刻 | 0 以上 | `exit_code:0` | object または `null` | `[]` | `null` |
| pipeline 成功 / deploy なし / snapshot 成功 | snapshot 保存時刻 | 0 以上 | `exit_code:0` | object または `null` | `[]` | build id |
| pipeline 成功 / deploy pending | deploy 判定時刻 | 0 以上 | `exit_code:0` | object または `null` | pending entry | `null` |
| pipeline 成功 / deploy 成功 / snapshot 成功 | snapshot 保存時刻 | 0 以上 | `exit_code:0` | object または `null` | success entry | build id |
| snapshot 失敗 | snapshot 失敗時刻 | 0 以上 | `exit_code:0` | object または `null` | `[]` または success entry | `null` |
| rollback 成功 | rollback deploy 完了時刻 | 0 以上 | `exit_code:null`、stdout/stderr 空 | `null` | success entry | 元 snapshot id |
| rollback deploy pending | pending 保存完了時刻 | 0 以上 | `exit_code:null`、stdout/stderr 空 | `null` | pending entry | 元 snapshot id |
| rollback 失敗 | 失敗確定時刻 | 0 以上 | `exit_code:null`、stdout/stderr 空 | `null` | 取得済み deploy result | 元 snapshot id |

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

`.build_history` の保存 key、型、必須条件、許容値は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.build_history` JSON Lines schema を参照する。

runner は build 結果確定後、`.build_history` へ 1 build につき 1 行だけ追記する。`status` は runner の最終結果、`trigger` は [`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) の有効値、`output_sha256` は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の出力成果物 manifest SHA-256 とする。manifest 生成に失敗した場合のみ `output_sha256:null` を許可する。JSON Lines 追記は `O_APPEND|O_CREATE|O_WRONLY` で行い、1 行全体を書き込んでから file sync する。

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
| rollback audit 失敗 | `rollback audit failed` |
| rollback worker 中断 | `rollback interrupted` |

---

<a id="15a-runner-受け入れ検証条件"></a>
**15a. `runner` 受け入れ検証条件：**

`runner` の初期実装は、[`docs/details/runner.md` 詳細本文責務 §15a](runner.md#15a-runner-受け入れ検証条件) の検証条件と [`docs/details/fixture.md`](fixture.md) fixture 証跡責務の fixture をすべて満たすまで完了として扱わない。`testdata/runner/` は runner fixture の配置予定 path であり、現時点で未作成の場合は現行実体として扱わない。fixture ファイルは [`docs/details/fixture.md`](fixture.md) fixture 証跡責務に従う実装変更で `testdata/runner/` 配下へ追加する。外部 GitHub API と SSH サーバーへ実接続するテストは初期 fixture に含めず、HTTP test server と fake `ssh` executable で再現する。

<a id="sec-15a-0"></a>
**15a.0 runner fixture 共通検証観点：**

[`docs/details/runner.md` 詳細本文責務 §15a.0](runner.md#sec-15a-0) は runner owner の共通検証観点だけを示す。fixture 名、入力状態、expected、fake GitHub / fake ssh / fake notifier / fake filesystem、実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) を正本とする。

| 共通検証観点 | 確認内容 | fixture 正本 |
|--------------|----------|--------------|
| no external execution | 対象 fixture が禁止する GitHub API、pipeline、deploy、snapshot を実行しない。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) |
| no log history creation | 対象 fixture が禁止する `.build_logs/` と `.build_history` を作成しない。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) |
| clean final state | `.build_state.running=false`、`current_build_id=null`、`.build_lock` 不在で終了する。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) |
| finalizer failure | finalizer 保存失敗を正常扱いせず、所有確認付き lock 解放を 1 回だけ実行し、ERROR 証跡を残す。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) |

<a id="sec-15a-1"></a>
**15a.1 runner 受け入れ fixture catalog 参照：**

`runner` 初期実装の受け入れ fixture catalog は [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) を正本とする。`runner` 詳細本文では、各 fixture の前提状態、fake response、expected file、実行 command、状態差分を再定義しない。

| fixture | runner owner 検証観点 | fixture 正本 |
|---------|----------------------|--------------|
| R1 | CLI 異常系、help、未知 option、state-dir validation。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R1` |
| R2 | 変更なし skip、SHA cache 維持、log/history 非作成。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R2` |
| R3 | 変更あり build 成功、deploy なし、log/history/state finalizer。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R3` |
| R4 | pipeline 失敗、SHA 非更新、deploy/snapshot 非実行。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R4` |
| R5 | deploy pending、pending transfer 保存、snapshot 非作成。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R5` |
| R6 | lock 競合、既存状態非変更。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R6` |
| R7 | 状態破損退避、再生成、継続処理。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R7` |
| R8 | pipeline timeout、旧 SHA 維持、clean final state。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R8` |
| R9 | GitHub API 全再試行失敗、pipeline/deploy/snapshot 非実行。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R9` |
| R10 | 通知失敗を build 成功へ反転しない。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R10` |
| R11 | `[REPORT]` 重複時の採用行と警告。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R11` |
| R12 | finalizer state write failure。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R12` |
| R13 | GitHub token mode 不正、secret 非出力。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R13` |
| R14 | dry-run directory 非作成、外部副作用なし。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R14` |
| R15 | SHA cache 破損、pipeline/deploy/snapshot 非実行。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R15` |
| R16 | GitHub rate limit reset 不正、長時間待機禁止。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R16` |
| R17 | cooldown skip と manual force queue。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R17` |
| R18 | SSH checksum mismatch pending。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R18` |
| R19 | pending transfer 重複統合。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R19` |
| R20 | snapshot atomic save and prune。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R20` |
| R21 | status start write failure、外部副作用なし。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R21` |
| R22 | running build log 開始時保存失敗、外部副作用なし、history 非追記、SHA 非更新。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R22` |
| R23 | history append failure、最終 log と完了済み SHA / deploy / snapshot 維持、history 自動再追記禁止。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R23` |
| R24 | multi target partial failure continues。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R24` |
| R25 | all targets GitHub API failure。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R25` |
| R26 | snapshot failure remains success。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R26` |
| R27 | build_state finalizer failure keeps failure。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R27` |
| R28 | process crash 後の active entry 再実行。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R28` |
| R29 | queue 結果保存失敗時の active entry 保持。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R29` |
| R30 | circuit open 中の pending retry、queue 保持、runtime flag 終了状態。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R30` |
| R31 | final build log atomic replace 失敗、running log 維持、history 非追記、完了済み副作用非巻き戻し。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R31` |
| R32 | repository identity の既定値、API 更新後の次回起動適用、破損時の GitHub API 呼出し禁止。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) `Fixture R32` |

---

<a id="16-systemd-タイマー参照"></a>
**16. systemd タイマー参照：**

systemd unit 本文、配置先、起動手順、更新手順は setup owner component の責務とし、[`docs/details/setup.md` 詳細本文責務 §26.4.1](setup.md#sec-26-4-1)、[`docs/details/setup.md` 詳細本文責務 §26.5](setup.md#sec-26-5) を参照する。

runner owner component は、`adlaire-ci-runner --state-dir /opt/adlaire-builder` として oneshot 実行された場合の処理、終了コード、状態ファイル更新、ログ出力だけを定義する。

runner 実装は systemd unit file を生成、配置、更新、enable、restart してはならない。systemd 操作が必要な機能は [`docs/details/setup.md`](setup.md) 詳細本文責務または [`docs/details/api.md`](api.md) 詳細本文責務の owner component 別詳細本文責務で定義する。

runner が journal へ出力する内容は [`docs/details/runner.md` 詳細本文責務 §15](runner.md#15-ログ) のログ仕様を参照する。`systemctl`、`journalctl` の操作手順は `runner` 詳細では定義しない。

---

## 17. GitHub 連携前提

| 項目 | 内容 |
|------|------|
| PAT 基本権限 | Fine-grained PAT の対象リポジトリに `Contents: Read` を付与する。Commit Status が無効な標準構成では、GitHub repository permission をこれより広げてはならない。 |
| Commit Status 有効時の追加権限 | `.server_config.commit_status_enabled=true` の場合だけ、同じ対象リポジトリに `Commit statuses: Write` を追加する。`false` の場合は付与しない。これ以外の GitHub 書き込み権限を追加してはならない。 |
| PAT の種類 | Fine-grained PAT（特定リポジトリのみ許可）を使用する。 |
| Webhook 設定（ポーリング方式） | **不要**（デフォルト。`BRANCH_TARGETS` によるポーリングのみ使用する場合） |
| Webhook 受信方式の前提 | `POST /api/webhook` endpoint と Webhook Secret が必要。endpoint、署名検証、request / response は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) と [`docs/details/api.md` 詳細本文責務 Webhook 受信境界](api.md#webhook-receive-overview) を参照する。 |
| 外部公開境界 | GitHub から webhook を受信する場合の外部公開、TLS 終端、リバースプロキシ構成は runner 詳細本文責務では定義しない。 |

---

<a id="18-初回セットアップ手順参照"></a>
**18. 初回セットアップ手順参照：**

初回セットアップ、Release asset 取得、checksum 検証、バイナリ配置、secret 初期化、状態ファイル初期化、systemd unit 書き込み、service 起動、管理 API 導入、管理 UI 配置は setup owner component の責務とし、[`docs/details/setup.md` 詳細本文責務 §26.1](setup.md#sec-26-1)〜[§26.4](setup.md#sec-26-4) を参照する。

runner owner component は、セットアップ済み環境で `/usr/local/bin/adlaire-ci-runner` が起動された後の処理だけを定義する。

runner 実装は以下を行ってはならない。

| 禁止条件 | 理由 |
|----------|------|
| OS user 作成、directory 作成、chown / chmod の初期設定 | setup owner component の責務。 |
| Release asset 取得、checksum 検証、バイナリ配置 | setup owner component の責務。 |
| `.github_token` の新規生成または対話入力 | setup owner component の secret initializer の責務。 |
| `.admin_credentials` 初期化、API service 配置、管理 UI 配置 | api / admin / setup owner component の責務。 |
| systemd unit file の配置、enable、restart | setup owner component の責務。ただし API endpoint が systemd timer を変更する機能は [`docs/details/api.md` 詳細本文責務 §27.11](api.md#sec-27-11) を参照する。 |

runner が起動時に必要ファイル不足または権限不備を検出した場合は、[`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15a](runner.md#15a-runner-受け入れ検証条件)、[`docs/details/runner.md` 詳細本文責務 §20](runner.md#20-ci-ランナー-既知の制限) の異常系に従い、セットアップ手順を自動実行せずに失敗として記録する。

---

<a id="19-管理-api-サーバー制限参照"></a>
**19. 管理 API サーバー制限参照：**

管理 API サーバーの HTTP listener、認証、session、rate limit、TLS 非対応、外部認証非対応、worker pool 非採用の制限は api owner component の責務とし、[`docs/details/api.md` 詳細本文責務 §21a](api.md#21a-管理-api-サーバー制限) および [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) を参照する。

runner owner component は、管理 API サーバーの起動、listener、session、認証、HTTP response、rate limit を実装してはならない。

runner と api が同じ状態ファイルを参照する場合でも、runner は API session、API token、TOTP、rate limit、HTTP access log を読み書きしない。runner が読み書きする状態ファイルは [`docs/details/runner.md` 詳細本文責務 §11](runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](runner.md#15-ログ)、[`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a)、[`docs/details/api.md` 詳細本文責務 §22.0d](api.md#sec-22-0d)、および対象の [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) 機能契約に明記されたものだけとする。

---

## 20. CI ランナー 既知の制限

<a id="sec-20-1"></a>
**20.1 Go 版 runner の制限：**

| 制限 | 詳細 |
|------|------|
| ポーリング遅延 | 変更検出は systemd timer の実行間隔に依存する。即時反応が必要な場合は `POST /api/webhook` を併用する。 |
| `pipeline.sh` 起動 | ビルド起動は `src` と同じディレクトリ配下の `.ci/pipeline.sh` を標準とする。YAML 形式のパイプライン定義は `runner` 詳細では定義しない。 |
| `BRANCH_TARGETS` 直列処理 | 複数エントリはリスト順に順次処理する。並列処理は行わない。1 件の処理が失敗しても、失敗をログと `.build_logs/{id}.json` に記録した上で次エントリへ進む。 |
| GitHub API リトライ | GitHub API 失敗時は `API_RETRY_MAX` 回まで指数バックオフで再試行する。全試行失敗時は ERROR ログを記録し、当該ターゲットのビルドをスキップする。SHA は更新しない。 |
| ビルド失敗時の扱い | `pipeline.sh` が非 0 で終了した場合は ERROR ログを出し、SHA を更新しない。次回実行では同じ blob SHA を再検出して再度ビルド対象になる。 |

<a id="sec-20-2"></a>
**20.2 標準機能の制限：**

以下は管理 API または追加拡張と連携する場合の制限である。

| 制限 | 詳細 |
|------|------|
| Webhook 受信の外部公開 | `POST /api/webhook` は `api`（`127.0.0.1` バインド）で受信するため、GitHub から直接受信する構成ではリバースプロキシと TLS 終端が必要。 |
| ペンディングキュー | ペンディング再試行が失敗した場合、`retry_count` を 1 増やしてエントリを保持する。runner による自動放棄は行わない。削除は転送成功時、または管理 API / 手動運用で明示的に削除する場合に限定する。 |
| ペンディングキュー肥大化 | `queue_max_size` を超えた新規投入は ERROR ログを記録し、新規エントリを追加しない。既存エントリは削除しない。 |
| サーキットブレーカー | 連続失敗回数が `API_CIRCUIT_BREAKER_THRESHOLD` 以上になった場合はポーリングを停止し、`POST /api/circuit-breaker/reset` でのみ復帰する。 |

---

## 27. Runner owner 追加仕様化機能 詳細仕様

<a id="runner-27-api--sdk--ui-共通参照先"></a>
**[`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) API / SDK / UI 共通参照先：**

[`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) の各機能で API endpoint、HTTP status、request / response、warning、SDK method、UI 表示、filter、error body を述べる場合、API 契約は [`docs/details/api.md` 詳細本文責務 §22.0c.1](api.md#sec-22-0c-1) / [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e)、SDK 契約は [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様)、UI 契約は [`docs/details/ui.md` 詳細本文責務 §24](ui.md#24-標準管理ツール-仕様) を共通参照先とする。runner 詳細本文では、runner が保存する値、処理順、状態差分、失敗時副作用だけを定義する。

<a id="sec-27"></a>
**[`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) owner / collaborator 境界参照先：**

[`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) の各機能契約で owner / collaborator を宣言する場合、その宣言は対象機能の実装境界確認であり、境界管理は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0b.1](../DETAIL_INDEX.md#0b1-owner-component-別-owner-collaborator-境界管理) を正本とする。同一 collaborator 組み合わせの機能契約でも、機能単位の境界確認として維持する。

<a id="sec-27-1"></a>
**27.1 GitHub Commit Status API：**

[`docs/details/runner.md` 詳細本文責務 §27.1](runner.md#sec-27-1) の runner 側境界本文は [`docs/details/commitstatus.md` 詳細本文責務 §27.1](commitstatus.md#sec-27-1) を参照する。owner component は `commitstatus`、collaborator component は `runner`、`statefile` とする。

runner は、commit SHA 確定、build id 採番、build 開始前の pending 送信呼び出し、pipeline / deploy / snapshot / history の最終結果確定後の final 送信呼び出しだけを担当する。GitHub Commit Status API payload、送信順、送信失敗時の非反転、保存値、secret mask、検証条件は [`docs/details/commitstatus.md` 詳細本文責務 §27.1](commitstatus.md#sec-27-1) を基準とし、[`docs/details/runner.md`](runner.md) 詳細本文責務へ重複定義してはならない。

<a id="sec-27-2"></a>
**27.2 ドライラン実行モード：**

[`docs/details/runner.md` 詳細本文責務 §27.2](runner.md#sec-27-2) の境界は owner component `runner`、collaborator component `statefile` とする。

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
      "current_commit_sha": "89abcdef0123456789abcdef0123456789abcdef",
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
| `targets[].previous_sha` | `string` または `null` | dry-run 開始時点の保存済み SHA。存在しない場合は `null`。 |
| `targets[].current_blob_sha` | `string` または `null` | fake GitHub read で取得した blob SHA。GitHub read 失敗または precheck failure では `null`。 |
| `targets[].current_commit_sha` | `string` または `null` | fake GitHub read で取得した commit SHA。GitHub read 失敗または precheck failure では `null`。 |
| `targets[].would_build` | boolean | 非 dry-run なら build を開始する場合だけ `true`。`github_error`、`config_error`、`precheck_error` では `false`。 |
| `targets[].trigger` | string | 判定 trigger。既定は `"polling"`。 |
| `targets[].reason` | string | `"sha_changed"`、`"no_change"`、`"cooldown"`、`"circuit_open"`、`"config_error"`、`"github_error"`、`"precheck_error"` のいずれか。 |
| `targets[].warnings` | array[object] | target 固有 warning。形式は `warnings[]` と同じ。 |
| `would_call` | array[string] | 実行予定の外部 API / command 種別だけを固定文字列で返す。secret、URL query 全体、token は含めない。 |
| `would_write` | array[string] | 非 dry-run なら発生する論理書込予定を固定文字列で返す。許可値は `"lock"`、`"build_log"`、`"history"`、`"status"`、`"sha_cache"`、`"notification"`、`"deploy"`、`"snapshot"`、`"commit_status"`。実書込が発生したことを意味しない。 |
| `secrets_masked` | boolean | stdout、stderr、expected、effects に secret 平文が残らない検証を通した場合だけ `true`。 |
| `errors` | array[object] | `code` と `message` は string、`target` は string または `null` の object。 |
| `warnings` | array[object] | `code` と `message` は string、`target` は string または `null` の object。 |

`would_call` の許可値は `"github_tree"`、`"github_blob"`、`"github_commit"`、`"github_rate_limit"` に限定する。dry-run は fake GitHub read だけを実行対象にし、実 GitHub write API、GitHub Commit Status、SSH、pipeline、deploy、snapshot、notification、systemd、hook、remote build を `would_call` に含めてはならない。

dry-run は、破損 state の backup、初期値作成、lock 作成、通知、GitHub Commit Status、deploy、archive、cleanup、quarantine、状態正規化を実行しない。stdout 以外の状態差分が発生した場合は実装不合格とする。

終了コードは [`docs/details/runner.md` 詳細本文責務 §27.2](runner.md#sec-27-2) の固定表に固定する。`--help` / `--version` と同時指定された場合は `--help` / `--version` を優先し、dry-run JSON を出力しない。

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

<a id="sec-27-3"></a>
**27.3 ビルド失敗時の自動リトライ：**

[`docs/details/runner.md` 詳細本文責務 §27.3](runner.md#sec-27-3) の境界は owner component `runner`、collaborator component `statefile` とする。

runner は `.server_config.build_retry_max > 0` の場合、retry 対象失敗だけを同一 build id 内で最大 `build_retry_max` 回追加試行する。総試行回数は `1 + build_retry_max` とする。

retry 対象は以下に限定する。

| 対象 | 条件 |
|------|------|
| GitHub API | network error、HTTP 429、HTTP 500〜599。 |
| pipeline | timeout。exit code 非 0 は retry しない。 |
| deploy | SSH 接続失敗、検証用 checksum 取得失敗、network timeout。checksum mismatch は retry しない。 |

retry 待機秒数は `build_retry_base_seconds * attempt` とする。初回失敗後の retry 1 回目は `base * 1`、retry 2 回目は `base * 2`。待機中に process が終了した場合、未完了 retry を再開してはならない。

attempt ごとの結果は `.build_logs/{id}.json.attempts[]` に必ず保存する。最終 attempt が成功した場合、`.build_history.status` は `"success"` または deploy 結果に応じた `"success_deploy_pending"` とし、`retry_count` に追加 retry 回数を保存する。全 attempt 失敗時は最終失敗段階に対応する詳細結果を保存し、総称 `"failure"` を history へ保存してはならない。SHA 更新、snapshot、deploy 成功記録は最終成功時だけ行う。

**attempt 更新固定契約：**

Attempt object の key、型、列挙値は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) を正とする。runner は `attempt` を 1 から連番化し、stage、開始・終了時刻、成功 / 失敗、retry 可否、固定 error を実行順で保存する。secret、token、command 全文を `error` へ含めてはならない。

retry 待機中に SIGTERM、context timeout、lock 喪失を検出した場合は待機を中断し、未実行 attempt を作成せず、現在までの attempts だけを保存する。最終成功時だけ `.last_sha`、snapshot、deploy success、history success を確定する。途中失敗 attempt で SHA cache を更新してはならない。

検証条件:

| ケース | 期待結果 |
|--------|----------|
| API 429 後成功 | attempts 2 件、history success、retry_count 1。 |
| pipeline timeout 後成功 | attempts 2 件、SHA は最終成功後のみ更新。 |
| pipeline exit 1 | retry なし、history `failure_build`、retry_count 0。 |
| deploy checksum mismatch | retry なし、pending transfer 記録。 |
| retry 上限到達 | history は最終失敗段階の詳細結果、attempts は `1 + build_retry_max` 件。 |
| retry 中断 | 未実行 attempt を作らず終了する。 |
| secret error | attempts[].error に secret 平文が出ない。 |

<a id="sec-27-8"></a>
**27.8 ビルドステータスファイル出力：**

本機能の目的は、runner の現在状態と直近結果を `.build_status.json` に集約し、api、sdk、ui が同じ read-only 情報を参照できるようにすることである。

owner component は `runner` とする。collaborator component は `api`、`statefile` とする。`runner` は `.build_status.json` の唯一の通常更新責務を持つ。`api` は `GET /api/status`、`GET /api/dashboard`、`GET /api/health` で read-only 参照する。api は `.build_status.json` を自動修復してはならない。

**入力：**

| 入力 | 説明 |
|------|------|
| runner 起動状態 | lock 取得、起動時整合性チェック結果、target 処理結果、pending transfer / notify 件数、circuit 状態。 |
| `.build_state` | `running`、`current_build_id`、`active_queue_entry`、`last_started_at`、`last_finished_at`、`queued`。 |
| `.build_history` | 直近 build id、status、trigger、duration。 |
| `.pending_transfers` | pending transfer 件数。 |
| `.notify_pending` | pending notify 件数。 |
| `.build_circuit_state` | circuit open 状態、連続失敗回数。 |

**出力：**

`.build_status.json` は [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) / [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の schema、文字コード、改行、file mode、atomic write 契約に従う。runner は状態値と更新タイミングだけを決定し、保存形式を再定義しない。

**更新タイミング：**

| タイミング | 必須値 |
|------------|--------|
| 起動時整合性チェックで復旧または停止が発生した直後 | `status="warning"` または `"failure"`、`last_trigger="startup_config_integrity"`、`running=false`。 |
| queue build 開始前 | pending retry と circuit closed 判定後、waiting entry を `active_queue_entry` へ移し、`status="running"`、`running=true`、`current_build_id`、`last_trigger`、`last_started_at` を保存する。自動 polling build は active を変更しない。 |
| 変更なし skip | `status="skipped"`、`last_target_status="skipped_no_change"`、`running=false`。build id は更新しない。 |
| cooldown skip | `status="skipped"`、`last_target_status="skipped_cooldown"`、`running=false`。 |
| build 成功 | `status="success"`、`last_build_id`、`last_finished_at`、`last_duration_seconds`、`output_sha256` を保存する。 |
| deploy result 確定 | deploy attempt を実行した場合だけ `last_deploy_at` に確定時刻、`last_deploy_status` に `success` / `failure` / `pending` を保存する。deploy 未実行 build では既存の 2 field を保持する。 |
| deploy pending | `status="success"`、`last_deploy_at`、`last_deploy_status="pending"`、`pending_transfers_count` を保存する。 |
| build 失敗 | `status="failure"`、`last_error`、`last_target_status`、`last_finished_at` を保存する。 |
| runner finalizer | `running=false`、`current_build_id=null` を必ず保存する。queue 結果確定時だけ active を `null`、未確定時は active を保持する。 |

`status`、`last_target_status`、`last_deploy_status` の列挙値は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.build_status.json` schema を正とする。runner は正規化状態を `status`、詳細結果を `last_target_status` へ保存し、同一 field へ混在させてはならない。

**異常系：**

| 条件 | 処理 |
|------|------|
| `.build_status.json` 書き込み失敗 | ERROR ログ `BUILD_STATUS_WRITE_FAILED: path={path} error={reason}` を出し、runner 終了コードを最低 `1` にする。build 成功後に発生した場合も終了コードは `1` とする。 |
| `.build_status.json` 破損を API が検出 | runner は破損 status file を自動修復しない。 |
| `.build_status.json` 不在 | runner は API 読取のために status file を作成しない。 |
| pending 件数読取失敗 | 件数を `null` にせず `0` として返してはならない。status 書き込み時はエラー扱いにし、`last_error` に固定文言を保存する。 |

**セキュリティ：**

`.build_status.json` に GitHub PAT、Webhook URL secret、SMTP password、API token、session token、request body を保存してはならない。`last_error` は最大 500 文字に切り詰め、改行は `\n` 文字列へ escape する。

**status 算出・不一致固定契約：**

| 項目 | 仕様 |
|------|------|
| `running=true` 不一致 | `.build_lock` が存在しないのに `running=true` の場合、API は自動修復しない。runner は process lock 取得後、circuit open gate または当該起動の state start / finalizer が定める runtime field だけを更新し、active / waiting を推測変更しない。 |
| active retry | `active_queue_entry != null` かつ valid running lock がない場合、runner は active を削除または waiting へ戻さず、process lock と circuit closed 判定後に active を最優先で再実行する。circuit open なら保持する。 |
| pending 件数 | `.pending_transfers` と `.notify_pending` は読めた場合だけ件数算出対象にする。読取失敗時も runner は件数を補完保存しない。 |
| history 不一致 | `.build_status.json.last_build_id` と [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の共通順序で選んだ最新の status summary 対象 history id が異なる場合、runner は status file を書き換えない。history-only 行は比較対象にしない。 |
| finalizer | finalizer は既存 `last_build_id`、`last_finished_at` を消さず、`running=false` と `current_build_id=null` だけを最低更新する。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| build 開始 | `status="running"`、`running=true`、`current_build_id` が保存される。 |
| queue build 開始 | waiting から選択した entry が同じ atomic write で active へ移り、同じ id は waiting から消える。 |
| queue result 確定 | log / history 確定後に active が消え、running / current build が clear される。 |
| queue result 未確定 | active が保持され、次回 runner が waiting より先に新しい build id で再実行する。 |
| build 成功 | `status="success"`、`running=false`、`last_build_id` と `last_duration_seconds` が保存される。 |
| 変更なし | build log / history を作らず、`.build_status.json` は `status="skipped"`、`last_target_status="skipped_no_change"` を保持する。 |
| 起動時整合性復旧 | `last_trigger="startup_config_integrity"`、復旧内容が `last_error` または warning として確認できる。 |
| 書込失敗 | runner 終了コードが最低 `1`、ERROR ログが出る。 |
| API read | runner は `.build_status.json` を保存入力として提供する。 |
| running stale | [`docs/details/runner.md` 詳細本文責務 §14b](runner.md#14b-スナップショット管理) の全条件に一致する中断 rollback だけを固定失敗として閉じる。それ以外の stale 状態は自動修復しない。 |
| finalizer | `last_build_id` を消さない。 |

<a id="sec-27-9"></a>
**27.9 ビルドトリガー種別の記録：**

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
| `startup_config_integrity` | 起動時整合性チェックで復旧または停止が発生する。 | build log / history は作成しない。 |
| `rollback` | `POST /api/history/{id}/rollback` により snapshot を再転送する。 | 新しい build id を作成する。 |
| `local_watch` | `watch_mode="local"` の local SHA 差分により build する。 | GitHub API を呼ばない。 |
| `approval` | `POST /api/approvals/{id}/approve` 由来の queue entry を処理する。 | 承認済み entry のみ。 |

runner は [`docs/details/runner.md` 詳細本文責務 §27.9](runner.md#sec-27-9) の固定表以外の値を保存してはならない。特に `"auto"`、`"force"`、`"scheduled"`、`"timer"` は runner 保存値として使用禁止とする。

**保存先：**

| 保存先 | 必須条件 |
|--------|----------|
| `.build_logs/{id}.json.trigger` | build log を作成する全処理で必須。 |
| `.build_history.trigger` | build history を追記する全処理で必須。 |
| `.build_status.json.last_trigger` | build、skip、復旧、rollback の最終 trigger を保存する。 |
| Queue entry `trigger` | `"manual"`、`"webhook"`、`"approval"` のみ許可する。 |
| API response / SDK / UI 連携 | runner は API / SDK / UI に提供する同じ trigger 値を保存する。 |

**判定順序：**

1. 起動引数が `--help` または `--version` の場合、trigger を確定しない。
2. 起動時整合性チェックで復旧または停止が発生した場合、`startup_config_integrity` を `.build_status.json` に記録する。
3. queue entry が存在する場合、entry の `trigger` を採用する。
4. pending transfer の再送だけで終了する起動は `retry_pending_transfer` とする。
5. `watch_mode="local"` で local SHA 差分がある場合は `local_watch` とする。
6. SHA 差分がある通常起動は `polling` とする。
7. SHA 差分がなく force interval 条件を満たす場合は `force_interval` とする。
8. rollback API が作成する処理は `rollback` とする。

複数条件が同時に成立した場合は、定義済み順序で最初に該当した trigger を採用する。1 回の runner 起動で複数 branch target を処理する場合、target ごとに同じ trigger を保存する。ただし queue entry が target を指定する場合は、対象 target のみにその trigger を適用する。

**api / sdk / ui 参照：**

`GET /api/history` の `trigger` query、HTTP status、response warning、SDK `getHistory({trigger})`、UI filter 表示は [`docs/details/runner.md` 詳細本文責務 §27 API / SDK / UI 共通参照先](#runner-27-api--sdk--ui-共通参照先) を参照する。`runner` 詳細では runner が保存する trigger 値、保存先、判定順序だけを定義する。

**異常系：**

| 条件 | 処理 |
|------|------|
| queue entry の trigger が不正 | queue entry を処理せず ERROR ログ `INVALID_TRIGGER: id={id} trigger={value}` を出す。 |
| 既存 history に未知 trigger がある | runner は新規保存で未知値を禁止する。 |
| build log と history の trigger 不一致 | runner は不一致を状態不整合として扱う。 |

**trigger 保存・queue 固定契約：**

| 項目 | 仕様 |
|------|------|
| queue 取り出し | queue entry の trigger は取り出し時に validation し、不正なら entry を削除せず処理を中断する。 |
| rollback | rollback API は queue を経由しない場合でも build log / history に `rollback` を保存する。 |
| startup | `startup_config_integrity` は build id を採番しない。 |
| unknown 既存値 | runner は新規保存で未知値を作らない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| polling build | log、history、status が `polling` で一致する。 |
| manual force | queue は `manual`、payload は `force=true`、保存 trigger は `manual`。 |
| webhook | 署名検証成功時だけ `webhook` が保存される。 |
| force interval | SHA 差分なしで `force_interval` が保存される。 |
| filter | runner は filter 対象となる `trigger` を保存する。 |
| 不正 trigger | runner は不正 queue entry を処理しない。 |
| startup trigger | build log / history を作らない。 |
| mismatch | runner は不一致を状態不整合として扱う。 |

<a id="sec-27-10"></a>
**27.10 設定ファイル起動時整合性チェック：**

本機能の目的は、runner が build 処理に入る前に、runner が読む状態ファイルの破損、型不一致、必須 key 不足、権限不備を検出し、規定どおり復旧または停止することである。

owner component は `runner` とする。collaborator component は `statefile` とする。管理 API、sdk、ui は本機能の実行責務を持たない。api が同じ状態ファイルを読む場合も、起動時整合性チェックを代行してはならない。

**対象ファイル：**

| 順序 | ファイル | 不在時 | 破損時 | unknown key |
|------|----------|--------|--------|-------------|
| 1 | `.branch_config` | 作成せず default 採用 | backup 後、不在扱い | 破損として backup 後、不在扱い |
| 2 | `.notify_config` | 初期値作成 | backup 後、初期値作成 | 破損として backup 後、初期値作成 |
| 3 | `.build_state` | 初期値作成 | backup 後、初期値作成 | 破損として backup 後、初期値作成 |
| 4 | `.build_circuit_state` | 初期値作成 | backup 後、初期値作成 | 破損として backup 後、初期値作成 |
| 5 | `.pending_transfers` | `[]` 作成 | backup 後、`[]` 作成 | 破損として backup 後、`[]` 作成 |
| 6 | `.notify_pending` | `[]` 作成 | backup 後、`[]` 作成 | 破損として backup 後、`[]` 作成 |
| 7 | `.maintenance` | file を作成せず disabled 扱い | backup、修復、上書きを行わず停止 | 破損扱いで差分なし停止 |

**実行順序：**

1. CLI 引数を検証する。
2. `--help` または `--version` の場合は設定ファイル起動時整合性チェックを実行しない。
3. `--dry-run` の場合は検証だけ行い、backup、初期化、復旧、通知、状態更新を行わない。
4. `StateDir` が絶対パスかつ既存ディレクトリであることを確認する。
5. `.build_lock` を取得する。
6. [`docs/details/runner.md` 詳細本文責務 §27.10](runner.md#sec-27-10) の固定表の順序で対象ファイルを検証する。
7. 復旧可能な問題は backup 後、ファイル別の初期状態へ復旧する。
8. 復旧結果に応じて `.build_status.json.last_trigger="startup_config_integrity"` を保存する。
9. 復旧通知条件を満たす場合は `config_corrupt` 通知を 1 回だけ送信する。
10. 停止条件がなければ pending retry、cooldown、target 処理へ進む。

**判定分類と停止条件：**

| 分類 | 処理 | runner 終了コード |
|------|------|------------------|
| `missing_optional` | `.branch_config` を作らず default 採用。 | 継続、最終結果に従う |
| `missing_required` | 初期値を atomic write。 | 継続、最終結果に従う |
| `parse_error` | `.maintenance` 以外は corrupt backup 後、ファイル別復旧。`.maintenance` は差分なしで停止。 | 通常は継続。`.maintenance` は `2` |
| `top_level_type_mismatch` | `.maintenance` 以外は corrupt backup 後、ファイル別復旧。`.maintenance` は差分なしで停止。 | 通常は継続。`.maintenance` は `2` |
| `required_key_missing` | `.maintenance` 以外は corrupt backup 後、ファイル別復旧。`.maintenance` は差分なしで停止。 | 通常は継続。`.maintenance` は `2` |
| `required_key_type_mismatch` | `.maintenance` 以外は corrupt backup 後、ファイル別復旧。`.maintenance` は差分なしで停止。 | 通常は継続。`.maintenance` は `2` |
| `invalid_value` | `.maintenance` 以外は corrupt backup 後、ファイル別復旧。`.maintenance` は差分なしで停止。 | 通常は継続。`.maintenance` は `2` |
| `unknown_key` | `.maintenance` 以外は破損元を backup 後、ファイル別の破損時処理を行う。未知 key だけを除去した保存は禁止する。`.maintenance` は差分なしで停止。 | 通常は継続。`.maintenance` は `2` |
| `permission_error` | 自動復旧しない。`.build_state.running` を変更しない。 | `2` |
| `io_error` | 自動復旧しない。`.build_state.running` を変更しない。 | `1` |

backup 名は `{original}.corrupt.{YYYYMMDDHHMMSS}.bak` とする。UTC 秒単位で衝突する場合は `{original}.corrupt.{YYYYMMDDHHMMSS}.{n}.bak` とし、`n` は `2` から始める。

**状態更新禁止条件：**

設定ファイル起動時整合性チェックだけで `.build_logs/{id}.json` と `.build_history` を作成してはならない。`startup_config_integrity` は `.build_status.json` の `last_trigger` にだけ記録する。ただし、設定ファイル起動時整合性チェック後に通常 build が発生する場合、通常 build の log / history は実際の build trigger を保存する。

`--dry-run` では、破損検出結果を dry-run JSON の `errors[]` または `warnings[]` に出力するだけとし、backup、初期化、復旧、通知、`.build_status.json` 更新を行わない。

**復旧通知：**

復旧通知は `.notify_config` の検証完了後、復旧対象に `.notify_config` と `.notify_pending` 以外のファイルが 1 件以上ある場合だけ送信する。送信イベントは `config_corrupt` とする。`.notify_pending` が破損復旧された場合、通知失敗時の pending 追記は行わない。

**復旧 record 固定契約：**

| 項目 | 仕様 |
|------|------|
| status warning | 復旧して継続する場合、`.build_status.json.status="warning"`、`last_error` に固定文言を保存する。 |
| status failure | permission / io error で停止する場合、失敗対象が `.build_status.json` ではなく state directory の path 検証、open、write が成功可能な状態なら、`.build_status.json.status="failure"` の atomic write を 1 回だけ試行する。失敗対象が `.build_status.json`、state directory 自体の permission / io error、またはこの 1 回の atomic write が失敗した場合は追加試行しない。atomic write 失敗時は `BUILD_STATUS_WRITE_FAILED` を ERROR で記録し、最初に確定した終了コードを変更しない。 |
| backup 内容 | 破損元 file を byte 単位でコピーする。整形、マスク、改行変換を行わない。 |
| unknown key | schema 未定義 key は破損として扱い、破損元を byte 単位で backup してからファイル別の破損時処理を行う。未知 key だけを削除した保存は禁止する。 |
| 通知 payload | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の NotificationPayload object とし、`event="config_corrupt"`、event 別 key `files`、`recovered` を設定する。secret 値と破損内容は含めない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 必須ファイル不在 | 初期値作成、終了コード `0`、追加 backup なし。 |
| `.branch_config` 破損 | backup 後に `.branch_config` 不在、default 採用。 |
| `.build_state` 破損 | backup、初期値作成、`startup_config_integrity` 記録。 |
| `.maintenance` 不在 | file を作成せず disabled として継続。 |
| `.maintenance` 破損 / 読取不能 | 自動修復、backup、pending retry、build log、history、SHA cache、deploy、snapshot 差分なし。status failure 固定契約を 1 回だけ適用し、終了コード `2`。 |
| unknown key | 破損元を backup 後、ファイル別の破損時処理を行い、未知 key だけを削除した保存を行わない。 |
| permission error | 自動復旧なし、終了コード `2`、running 未変更。 |
| dry-run | 差分なし、backup なし、dry-run JSON に検出結果。 |
| 通知失敗 | `.notify_pending` が正常な場合だけ pending 追記。 |
| backup byte | 破損元と backup が byte 単位一致。 |
| unknown only | 他の破損と同様に backup 後、ファイル別の破損時処理。 |

<a id="sec-27-14"></a>
**27.14 ビルド所要時間の記録と統計入力：**

本機能の目的は、build ごとの開始・終了・所要時間を構造化ログへ保存し、統計 API が参照する duration 入力を確定することである。

owner component は `runner` とする。collaborator component は `api`、`statefile`、`archive` とする。

**記録仕様：**

`runner` は `.build_logs/{id}.json` に `started_at`、`finished_at`、`duration_seconds` を必ず保存する。`started_at` は build id 採番直後、`finished_at` は最終 target status 確定直後とする。`duration_seconds` は `finished_at - started_at` を秒単位で切り上げず整数化し、1 秒未満は `0` とする。

`.build_history.duration_seconds` は `.build_logs/{id}.json.duration_seconds` と同じ値にする。失敗、deploy pending、rollback でも記録する。変更なし skip で build log を作らない場合は記録しない。

**統計読取連携：**

runner owner は、API が読む `.build_logs/` と `.build_logs/archive/` の `duration_seconds`、`finished_at`、`status` を保存する責務を持つ。

**duration 統計入力固定契約：**

| 項目 | 仕様 |
|------|------|
| 選択入力 | `finished_at` 降順、同時刻は build id 降順で最新 N 件を選べるよう、`id`、`finished_at`、`duration_seconds`、`status`、`target_status` を build log に保存する。 |
| archive 重複 | 同じ id の通常 log と archive は同一 build を表し、API が通常 log を優先できる。 |
| 公開契約 | `GET /api/stats/build-duration` の query、集計、丸め、response key、不在・破損時処理は [`docs/details/api.md`](api.md) 詳細本文責務の `BuildDurationStats` 契約を正本とする。 |

**duration 記録・統計実装確認固定契約：**

| 項目 | 仕様 |
|------|------|
| started_at | build id 採番直後、`.build_status.json.status="running"` 保存前の UTC 秒精度時刻。 |
| finished_at | finalizer が最終 status、deploy pending、rollback pending、failure を確定した時刻。 |
| duration_seconds | `finished_at - started_at` を秒単位で切り捨てる。負値になる fake clock / clock drift は `0` に丸め、WARN `DURATION_CLOCK_DRIFT` を出す。 |
| 保存順 | build log 最終更新 → `.build_history` 追記 → `.build_status.json` finalizer の順で同じ duration を保存する。 |
| deploy pending | deploy pending でも build 処理自体の finished_at を保存し、pending retry の所要時間を合算しない。 |
| rollback | rollback build log も duration を保存する。元 snapshot の duration は変更しない。 |
| skip | `skipped_no_change`、`skipped_cooldown`、`skipped_tag_filter`、`skipped_maintenance`、`circuit_open`、`lock_skipped` は build log / history / duration を作らない。chain の `skipped_dependency_failed` は history だけを作成し `duration_seconds=null` とする。 |
| read-only stats | stats 系 API が読む history、logs、archive、trend、status は runner 側で read-only 入力として提供される。runner は統計 API 呼び出しで状態を書き換えない。 |

**統計入力分類・丸め固定契約：**

| 参照元 | 対象 | 計算 |
|-----|------|------|
| `.build_history` | `finished_at` が API 指定範囲内の行。 | `success` / `success_deploy_pending` を成功、`failure_` prefix / `cancelled` / `hook_error` を失敗、`approval_` prefix / `skipped_dependency_failed` を build 集計から除外する。 |
| `.build_history` UTC 日付 bucket | API 指定範囲内の UTC 日付 bucket。 | 日付降順。build 0 件の日は返却対象に含めない。 |
| `.build_logs/` / `.build_logs/archive/` | duration あり完了 log。 | latest 判定後 API 指定件数、avg は小数第 2 位、min / max は整数。 |

破損 build log、破損 history 行、gzip 展開失敗は統計対象から除外し、固定 WARN code だけを出す。破損内容、secret 風値、stdout/stderr 本文は WARN に含めない。

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

<a id="sec-27-19"></a>
**27.19 週次ビルドサマリー Webhook：**

本機能の目的は、過去 7 日間の build 結果を指定曜日・時刻に集計し、Webhook へ定期通知することである。

owner component は `runner` とする。collaborator component は `api`、`statefile` とする。runner は自動送信、api は設定表示・手動送信を担当する。

**設定：**

`.notify_config.summary.enabled=true` の場合だけ有効とする。`interval` は `"weekly"`、`hour` は 0〜23、`day_of_week` は 0〜6 とする。タイムゾーンは UTC 固定。

**自動送信条件：**

runner 起動時に、現在 UTC の曜日と時が設定値に一致し、`.build_state.weekly_summary_sent_date` が当日でない場合に送信する。送信成功時だけ `weekly_summary_last_sent_at` と `weekly_summary_sent_date` を更新する。

**集計対象：**

`.build_history` のうち、現在時刻から過去 7 日以内の行を対象とする。`success` / `success_deploy_pending` を成功、`failure_` prefix / `cancelled` / `hook_error` を失敗として数える。`approval_` prefix / `skipped_dependency_failed` は build 件数から除外する。所要時間は `duration_seconds != null` の行だけ平均対象にする。

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
| 手動 weekly summary 連携 | API から受け取った対象時刻で、対象 channel 抽出 → 集計 → 通知送信 → `.notify_log` 追記を行い、構造化した結果または失敗を API caller へ返す。 | sent date と `.build_state` は更新しない。HTTP status、response、認証は [`docs/details/api.md`](api.md) 詳細本文責務の手動 weekly summary 契約に従う。 |
| 同日二重自動 | `.build_state.weekly_summary_sent_date` が現在 UTC 日付と一致する場合は送信しない | `.notify_log`、`.notify_pending`、`.build_state` を変更しない。 |

weekly summary payload は secret、repository token、SMTP password、Webhook secret、API token、session token を含めてはならない。`success_rate` は小数第 2 位まで `math.Round(x*100)/100` 相当で丸める。集計対象 0 件の場合は `success_count=0`、`failure_count=0`、`success_rate=0`、`avg_duration_seconds=null` とする。

**weekly 集計固定契約：**

| 項目 | 仕様 |
|------|------|
| 期間 | `now - 7*24h <= finished_at <= now`。timezone は UTC。 |
| 成功 | `status="success"`、`"success_deploy_pending"`。 |
| 失敗 | `status` が `failure_` prefix、`"cancelled"`、`"hook_error"`。 |
| 除外 | `approval_` prefix、`"skipped_dependency_failed"`、duration 欠落の平均対象。 |
| 最大 duration | payload に `max_duration_seconds`、`max_duration_build_id` を含める。対象なしは `null`。 |
| 手動 caller 返却値 | runner 側の共通処理は `period_from`、`period_to`、`success_count`、`failure_count`、`success_rate`、`avg_duration_seconds`、`max_duration_seconds`、`max_duration_build_id`、channel 別送信結果を API caller へ返す。HTTP response に公開する key は [`docs/details/api.md`](api.md) 詳細本文責務の固定契約に従う。 |

**weekly summary 実装確認固定契約：**

| 項目 | 仕様 |
|------|------|
| 実行位置 | runner 起動時の startup integrity、pending transfer retry、maintenance 判定後、通常 GitHub polling / local watch 差分検出前に判定する。 |
| channel 抽出 | `.notify_config.channels[]` の `enabled=true` かつ `on` に `weekly_summary` を含む channel を配列順に使う。legacy `webhooks` / root `on` だけがある場合は互換 channel として配列順に正規化する。 |
| 自動宛先なし | WARN `WEEKLY_SUMMARY_NO_CHANNEL` を出し、`.notify_pending` と sent date を変更しない。build status は変更しない。 |
| 手動宛先なし | runner は `.notify_log`、`.notify_pending`、sent date を変更しない。 |
| payload | `event`、`period_days`、`period_from`、`period_to`、`success_count`、`failure_count`、`success_rate`、`avg_duration_seconds`、`max_duration_seconds`、`max_duration_build_id` を含む。 |
| success_rate | 分母が 0 の場合 `0`。それ以外は `success_count / (success_count + failure_count) * 100` を小数第 2 位で丸める。 |
| avg_duration_seconds | 対象 duration がない場合 `null`。ある場合は秒の平均を小数第 2 位で丸める。 |
| notify log | channel ごとに [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の NotificationLog record を 1 件追記する。payload 本文全体は保存しない。 |
| pending | retry 対象 channel の送信失敗だけ `.notify_pending` に追加する。retry 非対象 channel は `.notify_log` だけに失敗を残す。 |
| secret mask | webhook secret、SMTP password、API token、session token、repository token、Authorization header は payload、notify log、pending、server log に保存しない。 |

手動 weekly summary は build lock を取得しない。自動 weekly summary も通常 build の `.build_lock` を取得しない。ただし `.notify_log`、`.notify_pending`、`.build_state` の書き込みでは [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の atomic write / lock 契約に従う。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 条件一致 | Webhook 送信、`.notify_log` 追記、sent date 更新。 |
| 同日二重起動 | 2 回目は送信しない。 |
| 宛先なし | 自動送信は WARN。 |
| 手動送信 | 手動送信時も sent date は変更しない。 |
| 送信失敗 | sent date を更新せず、retry 対象なら `.notify_pending` に追加。 |
| deploy pending | 成功として数える。 |
| skipped | 集計から除外する。 |

<a id="sec-27-21"></a>
**27.21 複数ファイル監視：**

本機能の目的は、単一 `target_file` 前提を拡張し、複数 Markdown ファイルまたは Markdown ディレクトリを 1 回の runner 起動で監視、差分判定、ビルド対象決定できるようにすることである。

owner component は `runner` とする。collaborator component は `builder`、`api`、`statefile` とする。runner は差分検出と build target 決定、builder は複数入力の静的サイト生成、api は設定表示・更新を担当する。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 schema | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.branch_config.branch_targets[].target_files`。型、既定値、上限、path 検証、保存順は同 schema を正とする。 |
| SHA cache | `.sha_cache/{branch_safe}/{target_hash}.sha` に target file 単位で保存する。file 内容の形式と不正値の扱いは [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.sha_cache/{branch_safe}/{target_hash}.sha` schema を参照する。 |
| build log | `.build_logs/{id}.json.changed_targets[]` に [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の ChangedTarget object を保存する。 |

**正常系：**

1. runner 起動時に `.branch_config` を読み、`target_file` と `target_files` を正規化する。
2. `target_files` に [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の正規化を適用し、正規化後の保存順で処理する。
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

**multi-file 実装確認固定契約：**

| 項目 | 仕様 |
|------|------|
| 差分検出単位 | target path ごとに before SHA、after SHA、source `github` / `local`、result `changed` / `unchanged` / `missing` / `error` を memory 上で確定してから build 可否を決める。 |
| build id | 1 runner 起動で複数 target が変更されても build id は 1 件だけ採番する。target ごとに build id を分けない。 |
| builder 入力 | builder へ渡す `--src` は branch target の `src` 1 件だけとし、target_files を複数 `--src` に展開しない。変更 target list は `ADLAIRE_CHANGED_TARGETS` だけで渡す。 |
| SHA 更新順 | build success finalizer 後に、changed target と force build 対象 target の SHA cache を target path 辞書順で更新する。 |
| 部分更新禁止 | build failure、deploy failure before success、pipeline failure、hook pre abort、SHA 部分取得失敗では target SHA cache を 1 件も更新しない。 |
| log | `.build_logs/{id}.json.changed_targets[]` は target path 辞書順で保存し、`before_sha` が不明な場合は `null`、`after_sha` が missing の場合は `null` とする。 |
| status | target が missing の場合は build を開始せず `failure_target_missing` とし、missing target を build log に保存する。 |
| dry-run | `--dry-run` では SHA cache、build log、history、status、snapshot、deploy、notification を変更せず、stdout JSON に target 判定だけを出す。 |

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
| 不正 path | runner は終了コード `2`。 |
| force build | 変更なしでも build 実行、成功時に全 SHA cache 更新。 |
| SHA 部分失敗 | build なし、SHA cache 差分なし。 |

<a id="sec-27-22"></a>
**27.22 ビルドパイプライン YAML 定義：**

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

`.pipeline_config` schema は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) を参照する。API による保存、request / response、HTTP status は [`docs/details/api.md` 詳細本文責務 ビルドパイプライン設定](api.md#pipeline-config-api) を参照する。

runner は build 開始後、builder command または pipeline step command を組み立てる直前に `.pipeline_config` を 1 回だけ読む。同一 build 中に `.pipeline_config` を再読込してはならない。

| 項目 | 仕様 |
|------|------|
| 読込タイミング | build id 採番、target 確定、`running=true` 保存後、builder / step command 組み立て直前。 |
| `extra_args` | 標準 builder command を使う場合にだけ、固定引数の後ろへ配列順で追加する。YAML step command には追加しない。 |
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
| block scalar <code>&#124;</code> / `>` | parse error。 |
| inline comment | quoted string 外の `#` は、行頭 comment 以外 parse error。 |

**pipeline 実装確認固定契約：**

| 項目 | 仕様 |
|------|------|
| source 優先順位 | `.pipeline.yml` が存在する場合は常に file を優先し、`.pipeline_config.inline_yaml` は読まない。file 不在時だけ inline YAML を読む。 |
| no pipeline | file と inline YAML がない場合は標準 builder command を使う。標準 builder command でも `.pipeline_config.extra_args` と `.pipeline_config.env` は適用する。 |
| pre-build failure | YAML parse、schema validation、禁止引数、command 解決失敗は build 本体、deploy、snapshot、SHA cache 更新を開始せず、build log に `failure_pipeline_config` を残す。 |
| step log | `pipeline_steps[]` は定義順で保存し、未実行 step は `status:"not_run"`、`exit_code:null`、`stdout:""`、`stderr:""` として保存する。 |
| stdout/stderr | 各 step の stdout/stderr は最大 64 KiB まで保存し、超過時は末尾を切り詰めて `truncated=true` を保存する。 |
| required failure | required step が failure / timeout の場合、以降の step は実行しない。ただし未実行 step は `not_run` として保存する。 |
| optional failure | optional failure は build status を反転しないが、`warnings[]` と `[WARN]` に固定 code を残す。 |
| secret | env 値、command args 内 secret 風値、stdout/stderr 内 secret 値は保存前に mask する。mask 不能なら build を失敗させ、平文を保存しない。 |

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

<a id="sec-27-23"></a>
**27.23 ローカルファイル監視モード：**

[`docs/details/runner.md` 詳細本文責務 §27.23](runner.md#sec-27-23) の境界は owner component `runner`、collaborator component `statefile` とする。

本機能の目的は、GitHub API を使わない環境で、ローカル Markdown 入力の変更を SHA-256 snapshot により検出することである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 schema | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.server_config.watch_mode`。型、許容値、既定値は同 schema を正とする。 |
| local root | branch target の `src`。絶対 path 必須。 |
| 状態 | `.local_watch_state.json` |
| trigger | local 差分起動時は `"local_watch"`。 |

`.local_watch_state.json` の object schema、path 制約、field 型、mode は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) を参照する。runner owner は走査、比較、更新タイミングだけを定義する。

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

**local watch 実装確認固定契約：**

| 項目 | 仕様 |
|------|------|
| GitHub 呼び出し禁止 | `watch_mode="local"` では GitHub Trees / Blobs / Commits / Tags / Rate Limit API、Commit Status API、PAT verify を 1 件も呼ばない。 |
| scan 順 | directory walk 結果は保存前に相対 path 辞書順へ sort する。filesystem の列挙順に依存しない。 |
| symlink | symlink file / directory は走査対象外とし、WARN `LOCAL_WATCH_SYMLINK_SKIPPED` を出す。 |
| state 破損 | `.local_watch_state.json` 破損は full build 扱いにするが、build 成功まで既存破損 file を上書きしない。 |
| out 除外 | `out` が `src` 配下の場合も、`out` 配下は必ず除外する。`out` が `src` 外の場合は除外 path として追加しない。 |
| 削除検知 | 削除だけの差分でも build を実行し、成功時に削除済み path を state から取り除く。 |
| trigger | local 差分 build の build log / history / status の `trigger` は必ず `local_watch` とする。 |

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

<a id="sec-27-24"></a>
**27.24 タグ付きコミットのみビルド：**

[`docs/details/runner.md` 詳細本文責務 §27.24](runner.md#sec-27-24) の境界は owner component `runner`、collaborator component `api`、`statefile` とする。

本機能の目的は、release tag が付いた commit だけを build 対象にする filter を提供することである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 schema | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.server_config.tag_filter` と TagFilter object。型、許容 pattern、既定値は同 schema を正とする。 |
| log | `.build_logs/{id}.json.matched_tags[]` |

**正常系：**

1. SHA 差分を検出する。
2. `tag_filter.enabled=true` の場合、対象 commit に紐付く tags を GitHub refs API から取得する。
3. `patterns` が空の場合は「任意の tag が 1 件以上」を条件とする。
4. tag が条件に一致した場合だけ build を実行する。
5. 不一致の場合は `.build_status.json.status="skipped"`、`last_target_status="skipped_tag_filter"` とし、SHA cache は更新しない。

**tag 判定固定契約：**

| 項目 | 仕様 |
|------|------|
| tag 正規化 | `refs/tags/` prefix を除いた tag 名で pattern 判定する。 |
| 取得順 | GitHub refs API response の順序を維持し、`matched_tags[]` には一致した tag を最大 100 件まで保存する。 |
| pattern `*` | suffix `*` は prefix match。`*` 単体は任意 tag に一致する。中間 `*`、正規表現、glob は禁止。 |
| skip 副作用 | tag 不一致 skip では `.build_status.json` だけ更新し、`.build_logs/{id}.json`、`.build_history`、SHA cache、snapshot、deploy、notify は更新しない。 |
| local mode | `watch_mode="local"` かつ `tag_filter.enabled=true` は設定不整合として終了コード `2`。 |

**tag filter 実装確認固定契約：**

| 項目 | 仕様 |
|------|------|
| 判定位置 | SHA 差分検出後、build id 採番前、pipeline / hook / builder 起動前に判定する。 |
| tag API | 対象 commit SHA に到達する tag だけを評価する。branch の最新 tag や repository 全 tag を無条件一致として扱わない。 |
| patterns 空 | `enabled=true` かつ `patterns=[]` の場合、tag が 1 件以上あれば一致とする。 |
| matched_tags | build 実行時だけ build log に保存する。skip 時は `.build_status.json.status="skipped"`、`last_target_status="skipped_tag_filter"` を保存する。未定義の `skip_reason` key は追加しない。 |
| cache | tag 不一致、tag API failure、pattern 不正、local mode conflict では `.last_sha`、`.sha_cache/*`、`.build_history` を更新しない。 |
| retry | tags API の `429` / timeout / 5xx は retry 対象、`404` / validation failure は nonretryable とする。 |
| secret | tag 名は secret として扱わない。ただし API error body、Authorization header、repository token は log に保存しない。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| GitHub tags API 失敗 | retry 対象。最終失敗時は build なし、終了コード `3`。 |
| pattern 不正 | runner は終了コード `2`。 |
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

<a id="sec-27-26"></a>
**27.26 並列マルチターゲットビルド：**

[`docs/details/runner.md` 詳細本文責務 §27.26](runner.md#sec-27-26) の境界は owner component `runner`、collaborator component `statefile` とする。

本機能の目的は、複数 deploy target への転送を bounded parallelism で処理し、遅い target が全体を不必要に止めないようにすることである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 schema | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.server_config.deploy_parallelism`。型、許容範囲、既定値は同 schema を正とする。 |
| 対象 | `branch_targets[].deploy_targets[]` |
| build log | `target_results[]` に [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の TargetResult object を保存する。 |

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
| pending 重複 | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の転送先識別子が一致する pending は重複追加せず、新しい build の値で置換する。 |
| 成功 target | 失敗 target があっても成功 target は pending に入れない。 |
| status | 1 件以上 pending があれば `success_deploy_pending`、全件成功なら `success`。 |
| 保存順 | `.build_logs/{id}.json.target_results` → `.pending_transfers` → `.build_history` → `.build_status.json`。 |

**parallel deploy 実装確認固定契約：**

| 項目 | 仕様 |
|------|------|
| worker 入力 | deploy target queue は設定順で作成し、各 target に固定 `target_id` を割り当てる。target id 未定義時は `branch_index-target_index` 形式を使う。 |
| timeout | target ごとの timeout は既存 deploy timeout を使う。timeout target は `status:"failure"`、`error_code:"deploy_timeout"` とする。 |
| result 保存 | `started_at` / `finished_at` は target 単位で UTC 秒精度。未開始 target は作らない。 |
| pending entry | 失敗 target は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の PendingTransfer object として保存し、初期 `retry_count=0` とする。 |
| 全失敗 | build 本体が成功している限り、全 deploy target 失敗でも status は `success_deploy_pending` とし、SHA cache は [`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) の build success 契約に従って更新できる。 |
| notify | deploy failure 通知は target_results と pending 保存後に送信する。通知失敗は build status を反転しない。 |
| panic 相当 | worker 内部 error は該当 target failure として扱い、他 worker を cancel しない。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| parallelism 不正 | runner は終了コード `2`。 |
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

<a id="sec-27-27"></a>
**27.27 ビルド前後フック：**

本機能の目的は、build 前後に登録済み command を安全に実行し、外部 shell 文字列に依存しない拡張点を提供することである。

[`docs/details/runner.md` 詳細本文責務 §27.27](runner.md#sec-27-27) の境界は owner component `runner`、collaborator component `api`、`statefile` とする。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 状態 | `.hooks` JSON object。mode は [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の状態ファイル固定値。 |
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
| hook log | 1 実行 1 JSON object とし、[`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の HookLog object schema に一致する値だけを保存する。runner はここで実行結果と保存タイミングを決定し、別 schema を定義しない。 |
| log 保存失敗 | pre hook の log 保存失敗は build を開始せず failure。post hook の log 保存失敗は build status を維持し runner 終了コードを最低 `1`。 |
| shell 禁止 | `command_args` を `exec.Command` 相当で実行し、shell 展開、変数展開、glob 展開を行わない。 |

**hook 実装確認固定契約：**

| 項目 | 仕様 |
|------|------|
| pre 実行位置 | build id 採番、status running 保存後、builder / pipeline / remote build 起動前に実行する。 |
| post 実行位置 | build / deploy / snapshot / history の最終 status 確定後、notification 送信前に実行する。 |
| abort | pre hook abort では builder、pipeline、remote build、deploy、snapshot、SHA cache 更新を行わない。history には `hook_error` を追記する。 |
| post failure | post hook failure は build status を反転しないが、hook log、build log warning、runner 終了コード最低 `1` を固定する。 |
| output limit | stdout/stderr は各 64 KiB まで保存し、超過時 `truncated=true`。secret mask は切り詰め前に適用する。 |
| process kill | timeout 時は process group 全体を kill し、kill 失敗は ERROR とする。shell は使わない。 |
| disabled hook | `enabled=false` の hook は実行せず、hook log も作らない。build log の `skipped_hook_ids` に id だけを昇順で保存する。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| pre hook 失敗かつ `abort_on_failure=true` | build 本体を実行せず build log `status="failure"`、`target_status="hook_error"`、history `status="hook_error"`。 |
| pre hook 失敗かつ `abort_on_failure=false` | WARN、build 継続。 |
| command 不正 | runner は該当 hook failure。 |
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

<a id="sec-27-29"></a>
**27.29 リモートビルド対応：**

owner component は `runner` とする。collaborator component は `api`、`archive`、`statefile` とする。

本機能の目的は、runner が SSH 先で build を実行し、成果物を archive と manifest で回収できるようにすることである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 schema | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.server_config.remote_build` と RemoteBuildConfig object。型、必須条件、許容値、既定値は同 schema を正とする。 |
| command_args | SSH 先で実行する論理 argv。各要素は [§14a の remote 引数引用固定契約](#14a-ssh-サイト転送) で個別引用し、値の shell 展開を禁止する。 |
| artifact | tar.gz。必須ファイル `site/`、`manifest.json`。 |
| log | `.build_logs/{id}.json.remote_build` |

**正常系：**

1. remote build enabled の場合、local builder を起動しない。
2. 各 remote build attempt の開始時に `build_timeout_seconds` の共通 deadline を作成する。
3. `work_dir` へ移動する固定 prefix と引用済み `command_args` で remote command を実行する。
4. remote command が終了コード `0` の場合だけ、同じ deadline 内で `artifact_path` を一時 file へ取得する。
5. `manifest.json` の SHA-256 と展開 file を検証する。
6. 検証成功後、既存 deploy 処理へ渡す。

**remote artifact 固定契約：**

| 項目 | 仕様 |
|------|------|
| 一時取得先 | state dir 配下 `.remote_artifacts/{build_id}/artifact.tar.gz.tmp`。directory は mode `0700`、file は exclusive create、mode `0600` とする。同じ build id の path が既に存在する場合は上書きせず失敗とする。 |
| artifact 取得 | `cat -- {quoted_artifact_path}` を remote command とし、SSH stdout の byte を変換せず `.tmp` へ stream する。remote host は `cat` の `--` を提供しなければならない。SSH stderr は artifact に混ぜず、失敗時の runtime ERROR log だけに使用する。 |
| 取得確定 | SSH 終了コード `0` の後、`.tmp` を `fsync`、close、同一 directory 内の `artifact.tar.gz` へ rename、directory `fsync` の順で確定する。途中失敗時は展開しない。 |
| 一時展開先 | `.remote_artifacts/{build_id}/site/` と `.remote_artifacts/{build_id}/manifest.json` のみ。既存 output directory へ直接展開しない。 |
| tar.gz entry | `site/` と `manifest.json` だけを root 直下必須とする。entry path の `..`、絶対 path、NUL、symlink、device は拒否する。 |
| manifest schema | `{ "files": [{"path": string, "sha256": string, "size": integer}] }`。 |
| 検証順 | tar.gz 展開前 entry 検査 → 一時展開 → manifest parse → file 存在 / size / sha256 検証 → deploy へ渡す。 |
| process output | build command の stdout/stderr は deadlock を防ぐため並行して drain する。stdout は破棄する。stderr と artifact fetch stderr は secret mask、CRLF → LF、不正 UTF-8 置換の後に末尾 65536 bytes までを失敗時の runtime ERROR log へ出す。`RemoteBuild` object、history、status、pending に stdout/stderr を保存しない。 |
| secret | artifact fetch の stdout は artifact byte stream としてのみ使用し、log に保存しない。SSH 秘密鍵 path、token 値、remote credential、command 引数を log に保存しない。`RemoteBuild.error` は入力値や remote output を連結しない固定文字列とする。 |
| cleanup | `.remote_artifacts/{build_id}` 作成後は、成功・失敗・timeout に関係なく、処理終了前に当該 directory だけを `os.RemoveAll` で 1 回削除する。削除失敗は `REMOTE_ARTIFACT_CLEANUP_FAILED` を WARN で記録し、削除を再試行せず、先に確定した remote build / deploy 結果と終了コードを変更しない。 |

remote build command は `cd {quoted_work_dir} && exec {quoted_arg_0} {quoted_arg_1} ...` の形だけを許可する。`cd`、`&&`、`exec`、区切り空白は固定文字列とし、`work_dir` と `command_args` の各要素は [§14a の `quoteRemoteArg`](#14a-ssh-サイト転送) を個別適用する。local 起動は `exec.CommandContext(ctx, "ssh", "--", user+"@"+host, remoteCommand)` とし、local shell は使用しない。引数内の空白、single quote、`$`、backtick、semicolon、glob 文字はすべて引数の一部として伝達し、shell 構文として展開してはならない。

**remote build 実装確認固定契約：**

| 項目 | 仕様 |
|------|------|
| local build 禁止 | `remote_build.enabled=true` の build では local builder command または local pipeline step を実行しない。pre/post hook は通常契約どおり実行する。 |
| remote command | 論理 argv の各要素を `quoteRemoteArg` で個別引用した固定形 remote command だけを実行する。未引用値、任意 shell fragment、glob、環境変数展開を禁止する。 |
| artifact fetch | remote command 成功後だけ、引用済み `artifact_path` に対する固定 `cat --` command で artifact を取得する。remote command 失敗時は artifact が存在しても取得しない。 |
| deploy 境界 | artifact 検証成功後だけ既存 deploy / snapshot / history 処理へ渡す。検証前に公開 output や snapshot を置換しない。 |
| log | `.build_logs/{id}.json.remote_build` は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の RemoteBuild object とし、credential、secret、SSH key path、command 引数を保存しない。 |
| cleanup failure | `.remote_artifacts/{build_id}` 削除失敗は WARN とし、build 成否を反転しない。cleanup 対象 path は state dir 配下だけに限定する。 |
| retry | SSH 接続、remote timeout、artifact fetch timeout は retry 対象。manifest mismatch、unsafe archive、schema error は nonretryable とする。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| SSH 接続失敗 | retry 対象。最終失敗で `target_status="failure_remote_build"`。 |
| remote command exit 非 0 | retry せず `target_status="failure_remote_build"`。artifact 取得と deploy を行わない。 |
| artifact 不在 | failure、deploy しない。 |
| manifest 不一致 | failure、deploy しない。 |
| remote command timeout | process kill、failure_timeout。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| remote 成功 | artifact 検証後 deploy。 |
| manifest 不一致 | deploy なし、failure。 |
| SSH 一時失敗 | retry 後成功なら success。 |
| shell metachar argument | remote 引数内の文字として完全一致で渡り、展開、command 置換、追加 command 実行が 0 件。 |
| unsafe tar entry | 展開中止、deploy なし。 |
| cleanup 失敗 | build 結果維持、WARN。 |

<a id="sec-27-30"></a>
**27.30 ビルド承認フロー：**

owner component は `runner` とする。collaborator component は `api`、`security`、`statefile`、`sdk`、`ui` とする。

本機能の目的は、`approval_required` な target を通常 build として即時実行せず、人間承認後の queue entry だけを build / deploy 実行対象にすることである。

API endpoint、approve / reject の request / response、sdk / ui 操作境界は [`docs/details/api.md` 詳細本文責務 §27.30](api.md#sec-27-30) を参照する。`.approval_queue` record schema は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) を参照する。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 key | `branch_targets[].approval_required` |
| timeout schema | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.server_config.approval_timeout_seconds`。型、許容範囲、既定値は同 schema を正とする。 |
| 承認状態 | `.approval_queue` JSON Lines。 |
| 実行 queue | `.build_state.queued[]` の `trigger:"approval"` entry。 |
| 通知設定 | `.notify_config` の `approval_required` event。 |

**正常系：**

1. runner は target 判定時、`approval_required=true` かつ通常 build 条件成立の場合、pipeline / deploy / snapshot を開始しない。
2. runner は `.approval_queue` lock を取得し、branch / sha / target / requested_trigger / requested_force / delivery_id がすべて一致する最新 `pending` record を確認する。
3. 重複 pending がない場合、runner は `pending` record を `.approval_queue` へ追記する。
4. runner は対象 approval id の `approval_pending` / `result:"success"` audit がない場合だけ `.audit_log` へ追記する。既存 audit は重複追記しない。
5. audit 成功後、対象 approval id と channel id の同一 payload に対する `.notify_log` または `.notify_pending` がない channel だけ approval request 通知を送信する。
6. 通知成功時は `.notify_log` に結果を記録し、通知失敗時は `.notify_pending` に retry entry を追記する。
7. runner は起動時に期限超過 pending を検出し、`expired` record、`.build_history.status="approval_expired"`、`.audit_log.action="approval_expired"` をこの順で追記する。
8. runner は approve API 由来の queue entry を取り出す前に、対応 approval id の最新 status が `approved` で、record の `queue_id` が entry id と一致することを確認し、一致時だけ `trigger="approval"` で build / deploy を実行する。

**pending 作成固定契約：**

| 項目 | 仕様 |
|------|------|
| pending id | base は `appr{YYYYMMDDHHmmss}` とし、衝突処理は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の時刻ベース ID 契約に従う。既存 id は再利用しない。 |
| requested force | `force_interval` または `POST /api/build/force` 由来は `requested_force=true`。通常 polling、通常 manual、webhook、local watch は `false`。 |
| 重複 pending | branch / sha / target / requested_trigger / requested_force / delivery_id がすべて一致し、最新 status が `pending` の record。 |
| 重複時 | 新規 record を作成せず、既存 pending id を使用する。成功済み audit と channel ごとの通知証跡を確認し、不足副作用だけを定義順で再開する。 |
| 通知 payload | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の NotificationPayload object とし、`event="approval_required"` の event 別 key を設定する。secret、token、path secret は含めない。 |
| build 抑止 | pending 作成成功または重複 pending 検出時、当該 target の build は開始しない。 |

**approval runner 状態遷移：**

| 現在 status | runner 操作 | 次 status | 副作用 |
|-------------|-------------|-----------|--------|
| なし | approval_required 検出 | `pending` | `.approval_queue` に pending record、`.audit_log` に `approval_pending` を追記し、両方の保存成功後にだけ [ビルド通知連携](#sec-27-32) の `approval_required` event を 1 件生成する。送信回数、retry、pending 保存は同節の channel 契約に従う。 |
| `pending` | 重複検出 | `pending` | 新規 record / build は発生させない。不足した `approval_pending` audit または channel 単位の通知処理だけを再開する。 |
| `pending` | timeout | `expired` | `.build_history` に `status:"approval_expired"`、`.audit_log` に `approval_expired` を追記する。queue は追加しない。 |
| `approved` | queue 取り出し | 変更なし | `trigger:"approval"` として build / deploy を実行する。 |
| `rejected` | runner 起動 | 変更なし | build しない。 |
| `expired` | runner 起動 | 変更なし | build しない。 |

**approval runner 更新順：**

| 操作 | 更新順 | 失敗時 |
|------|--------|--------|
| pending 作成 / 再開 | `.approval_queue` lock → 重複確認 →不在時だけ pending record append →成功済み `approval_pending` audit 有無確認 →不在時だけ audit append → channel ごとの通知証跡確認 →不足 channel だけ送信 → `.notify_log` / `.notify_pending` 更新 | pending append または audit 失敗時は build と後続通知を開始せず runner failure。audit 失敗時も pending record は戻さない。通知失敗でも pending と audit は残す。 |
| timeout | runner 起動時に `.approval_queue` lock → expires_at 超過 pending を created_at 昇順で抽出 → expired record append → `.build_history` append → `approval_expired` audit | history または audit 失敗時も保存済み expired record / history は戻さず、runner は ERROR を出して継続する。 |
| approved queue 実行 | `.build_state` lock → `trigger:"approval"` entry 候補確認 → `.approval_queue` 最新 record 読込 → `status="approved"` かつ `queue_id=entry.id` 一致確認 → entry 取り出し → requested_force 適用 → build / deploy 実行 → history / status finalizer | approval が pending、不在、queue id 不一致の場合は entry と queue 全体を変更せず `APPROVAL_QUEUE_NOT_FINALIZED` を記録し、build を開始せず終了コード `1`。queue entry schema 自体が不正な場合は `failure_decode`。 |

runner は `approval_required=true` の target に対して、approval queue 以外の経路で build を開始してはならない。manual force、webhook、force interval、local watch のいずれであっても、対象 target が approval_required の場合は pending 作成を優先し、承認済み queue entry になるまで pipeline を起動しない。

**approval 実装確認固定契約：**

| 項目 | 仕様 |
|------|------|
| pending record | key、型、status 別 nullable 条件は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.approval_queue` schema に従う。 |
| approved queue | approve API が作る queue entry には pending record の `branch`、`target`、`sha`、`requested_trigger`、`requested_force`、`delivery_id` と `approval_id` を payload へコピーする。runner は approval latest record と queue id の一致確認後、queue entry の値だけを使って build する。 |
| force build | force build でも `approval_required=true` なら pending 作成だけを行い、pipeline 起動は approval 後まで行わない。SHA cache は approval 前にクリアしない。approval 後は `requested_force=true` により SHA 一致 skip だけを bypass し、build 成功後に新 SHA を保存する。 |
| webhook | webhook 由来でも `approval_required=true` なら webhook queue を直接 build せず、approval pending へ変換する。delivery id は `.approval_queue.delivery_id` と approved queue payload に残す。 |
| timeout | timeout 判定は UTC fake clock で行い、`expires_at <= now` を expired とする。expired 後に approve された場合は API 側で `409`。 |
| audit | pending 作成は `approval_pending`、approve は `approval_approved`、reject は `approval_rejected`、expire は `approval_expired`。runner event は `actor_type:"system"`、`actor_id:"system"`、`request_id:null`、API event は request actor / request id を使用する。target は `target_type:"approval"`、`target_id:approval id`。失敗時は [`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) を参照する。 |
| secret | approval payload、notify payload、history、audit には token、Authorization header、repository secret を保存しない。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| `.approval_queue` lock 取得失敗 | build を開始せず終了コード `1`。 |
| pending append 失敗 | build を開始せず終了コード `1`。 |
| pending audit 失敗 | pending は維持し、通知と build を開始せず終了コード `1`。 |
| audit 失敗後の同一 pending 再検出 | 新規 pending を作成せず、不足 audit 追記から再開する。audit 成功後に証跡がない channel だけ通知する。 |
| approval queue の latest status 不一致 | queue entry を取り出さず、build を開始せず終了コード `1`。 |
| 通知失敗 | pending は維持し、`.notify_pending` に追記する。 |
| `.notify_pending` 追記失敗 | pending は維持し、ERROR ログを出して runner 終了コードを最低 `1` にする。 |
| timeout history append 失敗 | expired record は維持し、ERROR ログを出して継続する。 |
| timeout audit 失敗 | expired record と追記済み history は維持し、ERROR ログを出して継続する。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| approval required | build せず pending 作成、`approval_pending` audit 後に通知。 |
| duplicate pending 完了済み | 新規 record / audit / 通知なし、build なし。 |
| duplicate pending 部分失敗 | 新規 record なし、不足 audit / channel 通知だけ実行、build なし。 |
| notify failure | pending 維持、`.notify_pending` 追記。 |
| timeout | expired、history `approval_expired`、audit `approval_expired`、build なし。 |
| approved queue | `trigger="approval"` で build / deploy 実行。 |
| unfinalized approval queue | queue 差分なし、build なし、終了コード `1`。 |

<a id="sec-27-31"></a>
**27.31 ブランチ別環境変数：**

[`docs/details/runner.md` 詳細本文責務 §27.31](runner.md#sec-27-31) の境界は owner component `runner`、collaborator component `api`、`statefile` とする。

本機能の目的は、branch target ごとに build process へ注入する環境変数を定義し、branch や deploy 先ごとの差分を、保存前検証、注入対象固定、secret mask、log 保存禁止値によって扱うことである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 schema | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.branch_config.branch_targets[].env`。型、key、value、件数、保存順は同 schema を正とする。 |
| secret key | key に `TOKEN`、`SECRET`、`PASSWORD`、`PAT` を含むもの。 |

**正常系：**

1. runner は `.branch_config` へ保存された env key/value を読み取り、保存値以外を補完しない。
2. runner は build process の environment に branch env を追加する。
3. 同名 key が system env に存在する場合、branch env を優先する。
4. build log には env key 一覧だけを保存し、value は保存しない。
5. secret key は stdout/stderr の mask 対象に追加する。

**env 正規化・mask 固定契約：**

| 項目 | 仕様 |
|------|------|
| schema 適用 | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の schema に一致する保存値だけを使用する。schema 不一致時は build process を開始せず、状態ファイルを推測補正しない。 |
| secret 判定 | key に `TOKEN`、`SECRET`、`PASSWORD`、`PAT` を含む場合は大文字小文字を区別せず secret。 |
| log 保存 | `.build_logs/{id}.json.environment.env_keys` に key 名だけを保存する。value、value length、hash は保存しない。 |
| process env | branch env は builder、pipeline step、hook、command notification に渡す。通知 payload には値を含めない。 |
| mask failure | mask 対象値を保存前に置換できない場合、build を失敗扱いにし、平文を保存しない。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| key 不正 | runner は終了コード `2`。 |
| value 上限超過 | `422`。 |
| secret mask 漏れ | 実装不合格。該当 build は成功扱いにしない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| env 注入 | build command が指定値を参照できる。 |
| system env と同名 | branch env が優先される。 |
| secret stdout 出力 | log では `"***"` に置換。 |
| 不正 key | 保存不可、状態差分なし。 |
| lower secret key | `my_token` も secret 扱い。 |
| mask failure | build 成功扱いにしない。 |

**branch env 実装確認固定契約：**

| 項目 | 合格条件 |
|------|----------|
| 設定保存 | `.branch_config.branch_targets[].env` は key 検証、value 検証、ASCII 昇順保存、上限 100 key をすべて満たす。 |
| 注入範囲 | builder、pipeline step、hook、command notification の process env に同一値を渡し、system env 同名 key より branch env を優先する。 |
| secret mask | secret key の value は response、stdout、stderr、build log、history、notify log、pending、UI 表示、fixture effects に平文で残さない。 |
| ログ保存 | `.build_logs/{id}.json.environment.env_keys` には key 名だけを保存し、value、length、hash、部分文字列を保存しない。 |
| 失敗境界 | key 不正、value 不正、mask failure では build または保存を成功扱いにせず、失敗地点以降の状態更新を行わない。 |
| 確認条件 | fixture は env 注入、system env 上書き、secret stdout mask、不正 key、lower secret key、mask failure をすべて固定する。 |

<a id="sec-27-32"></a>
**27.32 ビルド通知連携：**

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
| notify id | `.notify_log` と `.notify_pending` の ID prefix、形式、衝突処理は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の各 schema に従う。 |
| payload | 共通 key、event 別追加 key、型、未知 key 禁止、canonical JSON は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の NotificationPayload object に従う。 |
| channel 順 | channel id 昇順。送信失敗しても次 channel を継続する。 |
| timeout | webhook と command は 30 秒。SMTP は 60 秒。timeout は retry 対象か [通知送信・pending 固定契約](#runner-notification-result-map) に従う。 |
| pending payload | mask 後 payload だけを保存し、secret は retry 送信直前に状態ファイルから再読込する。 |
| pending 保存順 | `.notify_log` 追記成功後に `.notify_pending` を保存する。log 失敗時は pending を追加しない。 |

**通知 channel 正規化契約：**

| 項目 | 仕様 |
|------|------|
| channel id | 未指定時は `n` + 6 桁連番。重複 id は `422`。 |
| event 判定 | channel `on` が空の場合は top-level `on` を使用する。channel `on` に `"*"` があれば全 event 対象。 |
| webhook config | `config.url` 必須。`http` / `https` のみ許可。userinfo、fragment、空 host は `422`。 |
| email config | `.smtp_config` / `.smtp_secret` を基準とする。channel `config.to` がある場合は `.smtp_config.to` より優先して送信先に使う。 |
| command config | `config.command_args` 必須。shell 経由は禁止。stdout/stderr は `.notify_log` に secret mask 後で保存する。 |
| secret mask | `secret`、`password`、`token`、`smtp_password`、Webhook secret、SMTP password は GET、backup、log、pending、UI 表示で `"***"`。 |

<a id="runner-notification-result-map"></a>
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

pending entry の key、型、保存順、payload schema は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.notify_pending` schema に従う。`event`、`channel_id`、secret mask 後 `payload` の canonical JSON が同一の pending entry が既にある場合は重複追加しない。Webhook secret や SMTP password は pending payload から復元せず、送信直前に該当状態ファイルから再読込する。

runner 起動時の pending retry は `next_attempt_at <= now` の entry を `created_at` 昇順、同時刻は `id` 昇順で処理する。成功した entry は削除する。失敗した entry は `attempts += 1`、`next_attempt_at = now + retry_interval_seconds` として保存する。初回送信失敗時に channel の `retry_count=0` なら pending を作成しない。retry 失敗後の `attempts > retry_count` になった entry は `.notify_log` に `result:"dropped"` を追記して pending から削除する。

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

**notification 実装確認固定契約：**

| 項目 | 合格条件 |
|------|----------|
| channel 選択 | event、enabled、channel `on`、top-level `on`、`"*"` の優先順位が固定契約どおり一致する。 |
| 送信順 | channel id 昇順で送信し、1 channel 失敗後も後続 channel を継続する。 |
| retry 境界 | webhook 5xx / timeout だけを pending 化し、4xx、SMTP 未設定、command 非 0、command timeout は pending 化しない。 |
| secret mask | GET response、backup、notify log、pending、command stdout/stderr、UI 表示に secret 平文を残さない。 |
| 状態保存順 | `.notify_log` 追記後に `.notify_pending` を保存し、log 失敗時は pending 追加判定を行わない。 |
| build 独立性 | 通知失敗、pending 保存失敗、retry exhausted は build status を反転しない。 |
| 確認条件 | fixture は複数 channel 成功、webhook pending、disabled/no target no-op、secret mask、retry exhausted、pending duplicate をすべて固定する。 |

<a id="sec-27-33"></a>
**27.33 ビルド時間トレンド記録：**

[`docs/details/runner.md` 詳細本文責務 §27.33](runner.md#sec-27-33) の境界は owner component `runner`、collaborator component `api`、`statefile` とする。

本機能の目的は、build 所要時間の統計を蓄積し、性能傾向と回帰検知の基準を提供することである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 状態 | `.build_trends.json` |
| API 参照 | runner は `GET /api/stats/build-trends?n=N` が読む `.build_trends.json` の入力値を提供する。 |
| sample | `{build_id, finished_at, branch, trigger, duration_seconds, status, target_status, anomaly}`。field schema は [`docs/details/statefile.md` 詳細本文責務](statefile.md) の `.build_trends.json` schema を正とする。 |
| 保持件数 | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.server_config.build_trend_keep_count`。型、既定値、許容範囲は同 schema を正とする。 |

**正常系：**

1. build 完了時、`duration_seconds != null` の場合だけ sample を追加する。
2. `finished_at` 昇順で保存し、保持件数超過分は古い順に削除する。
3. summary に `count`、`avg_seconds`、`median_seconds`、`p95_seconds`、`anomaly_count` を保存する。
4. runner は `build_trend_keep_count` 件を上限として sample と全 sample の summary を保存する。API query `n` は保存件数を変更しない。

**統計算出固定契約：**

| 項目 | 仕様 |
|------|------|
| 対象 sample | `duration_seconds` が 0 以上の数値で、`status` が `"success"`、`"failure"`、`"cancelled"` のいずれかであり、`target_status` との写像が [`docs/details/statefile.md` 詳細本文責務](statefile.md) の `.build_trends.json` schema に一致する sample。 |
| `avg_seconds` | 対象 duration の合計を対象件数で割り、小数第 3 位を四捨五入して小数第 2 位まで保存する。対象 0 件の場合は `null`。 |
| `median_seconds` | duration 昇順の中央値。偶数件の場合は中央 2 件の平均を使い、小数第 3 位を四捨五入して小数第 2 位まで保存する。 |
| `p95_seconds` | duration 昇順配列の `ceil(count * 0.95) - 1` 番目を採用する。`count=0` の場合は `null`。 |
| `anomaly_count` | `sample.anomaly == true` の件数。 |
| 同一 build id | 既存 sample と同じ `build_id` を追加する場合は append せず既存 sample を置換し、`finished_at` 昇順に再整列する。 |
| API response 参照 | runner は `.build_trends.json` の sample と summary を保存する。 |

`.build_trends.json` は sample 置換または追加後に summary を再計算して atomic write する。summary だけの部分更新は禁止する。

**異常系：**

| 条件 | 処理 |
|------|------|
| `.build_trends.json` 破損 | `.build_trends.json.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避後、`.build_history` の有効行から再集計する。 |
| 再集計不能 | 初期値で作成し、WARN を出す。 |
| API `n` 不正 | runner は API query を処理しない。API は状態更新なしで `422` を返す。 |
| `.build_history` に duration 欠落 | 当該行は再集計対象外とし、復旧操作 1 回につき server log へ WARN code `TREND_SAMPLE_SKIPPED` を 1 回だけ記録する。state と API response に warning key を追加しない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| build success | sample 追加、summary 更新。 |
| build failure | duration があれば sample 追加。 |
| 保持件数超過 | 古い sample だけ削除。 |
| API n=10 | 最新 10 件を返す。 |
| 同一 build id 再記録 | sample は 1 件のまま値が置換される。 |
| p95 算出 | `ceil(count * 0.95) - 1` の値と一致する。 |

**trend 実装確認固定契約：**

| 項目 | 合格条件 |
|------|----------|
| sample 追加 | build 完了時に duration が有効な場合だけ sample を追加し、`finished_at` 昇順で保存する。 |
| summary 再計算 | sample 追加、置換、prune、復旧後は summary を部分更新せず全再計算する。 |
| 数値丸め | avg、median は小数第 3 位を四捨五入して小数第 2 位、p95 は `ceil(count * 0.95) - 1` を使う。 |
| 同一 build id | 同じ `build_id` は append せず置換し、重複 sample を残さない。 |
| 破損復旧 | `.build_trends.json` 破損時は backup 作成後、`.build_history` 有効行だけから再集計し、skip warning を固定する。 |
| API 入力と response | `n`、選択順、response summary、HTTP error は [`docs/details/api.md`](api.md) 詳細本文責務の `BuildTrendStats` 契約に従う。runner は GET により状態を変更しない。 |
| 確認条件 | fixture は summary 更新、同一 build id 置換、保持件数 prune、破損復旧、API `n` 不正をすべて固定する。 |

<a id="sec-27-34"></a>
**27.34 ビルド依存チェーン：**

[`docs/details/runner.md` 詳細本文責務 §27.34](runner.md#sec-27-34) の境界は owner component `runner`、collaborator component `api`、`statefile` とする。

本機能の目的は、複数 build job の依存関係を DAG として定義し、依存 job 成功後だけ後続 job を実行することである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 状態 | `.build_chain_config` |
| API | `GET /api/build-chain-config` / `POST /api/build-chain-config` |
| 設定 schema | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.build_chain_config`。構造、job key、件数、参照整合、循環禁止は同 schema を正とする。 |

**正常系：**

1. API 保存時に schema、参照整合、循環を検証する。
2. runner は chain enabled job を topological order で実行する。
3. 依存 job が失敗した場合、後続 required job は `skipped_dependency_failed` として history に記録する。
4. `required=false` の依存失敗は WARN とし、後続 job を継続できる。
5. `.build_logs/{id}.json.chain` に job id、depends_on、chain_index を保存する。

**chain 実行固定契約：**

| 項目 | 仕様 |
|------|------|
| `chain_run_id` | chain 起動ごとの base は `chain{YYYYMMDDHHmmss}` とし、衝突処理は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の時刻ベース ID 契約に従う。chain 内の全 job log / history で同じ値を使う。 |
| topological 同順位 | 依存数が同じで同時に実行可能な job は、`.build_chain_config.chains[]` の出現順で実行する。 |
| disabled job | `enabled=false` の job は DAG から除外する。enabled job が disabled job に依存する場合、API 保存時に `422`。 |
| job build id | 各 job は独立した build id を持つ。build id は通常採番規則に従い、chain job であることを id 文字列へ埋め込まない。 |
| skip log | dependency failure により skip した job も `.build_history` に 1 行追記し、`.build_logs/{id}.json` は作成しない。 |
| chain summary | 最終 job 処理後、`.build_logs/{last_id}.json.chain_summary` に [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の ChainSummary object を保存する。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| 循環依存 | runner は chain 無効化して通常 build。 |
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

**chain 実装確認固定契約：**

| 項目 | 合格条件 |
|------|----------|
| 保存時検証 | job id、参照整合、disabled 依存、循環、件数上限を保存前に検証し、不合格時は状態差分を残さない。 |
| 実行順 | topological order と同順位の config 出現順が一致し、全 job log / history に同一 `chain_run_id` を保存する。 |
| skip 境界 | required dependency failure は後続 job を `skipped_dependency_failed` として history へ記録し、build log は作成しない。 |
| optional 境界 | `required=false` の依存失敗は WARN とし、後続 job の実行可否を固定契約どおり判定する。 |
| summary 保存 | 最終 job 後に `chain_summary` を保存し、保存失敗時は job 結果を維持したまま runner 終了コードを最低 `1` にする。 |
| 再開境界 | runner 停止後は完了済み job だけを既存履歴として扱い、未実行 job は次回再判定する。 |
| 確認条件 | fixture は DAG 実行、required skip、optional 継続、cycle 失敗、disabled 依存、summary 保存失敗をすべて固定する。 |

<a id="sec-27-35"></a>
**27.35 ビルド優先度キュー：**

[`docs/details/runner.md` 詳細本文責務 §27.35](runner.md#sec-27-35) の境界は owner component `runner`、collaborator component `api`、`statefile` とする。

本機能の目的は、manual、webhook、approval などの queue entry を優先度順に処理し、緊急 build を先に実行できるようにすることである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 状態 | `.build_state.queued[]` |
| queue schema | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の Queue entry schema。priority の許容値と既定値、created_seq の型は同 schema を正とする。 |
| 数値順 | urgent=0、high=10、normal=20、low=30 |
| created_seq | queue 追加時に単調増加する整数。 |

**正常系：**

1. queue 追加時に priority と created_seq を保存する。
2. circuit closed かつ active entry がある場合は priority に関係なく active を再実行し、active がない場合だけ priority 数値昇順、同一 priority では created_seq 昇順で waiting から 1 件を active へ移す。circuit open では active / waiting を変更しない。
3. runner が読む waiting queue は priority / created_seq / id 順で固定する。
4. runner は waiting queue を clear しない。API による clear 完了後は保存状態から消失した waiting entry を実行対象に戻さず、保持された active entry だけを引き続き実行対象にできる。clear の HTTP 契約は [`docs/details/api.md`](api.md) 詳細本文責務の `DELETE /api/queue` を正本とする。

**priority / created_seq 固定契約：**

| 項目 | 仕様 |
|------|------|
| 採番元 | `.build_state.queued[].created_seq` の最大値。存在しない場合は `0`。次 entry は最大 + 1。 |
| 旧 entry | `created_seq` 欠落 entry は GET 表示時だけ末尾扱いとし、runner 取り出し前に `.build_state` lock 内で `created_seq` を正規化保存する。正規化値は採番元の規則だけで決定し、entry 内容、priority、id、時刻から推測しない。 |
| priority 省略 | API / webhook / approval の queue 追加時は Queue entry schema の既定値を保存する。 |
| 表示順 | priority 数値昇順、created_seq 昇順、同値なら id 昇順。 |
| active 優先 | circuit closed かつ `active_queue_entry != null` の場合は waiting の priority に関係なく active を先に再実行する。circuit open では保持する。 |
| 取り出し順 | circuit closed かつ active がない場合に限り waiting の表示順と同一とし、選択 entry を waiting から active へ atomic move する。 |
| clear 後入力 | API が保存した clear 後の `.build_state` をそのまま読み、消失した waiting entry を再生成しない。active が保持されていれば active 優先契約を適用する。 |
| queue full | priority が高くても既存 entry を押し出さない。 |

runner が旧 entry の `created_seq` 正規化保存に失敗した場合、build を開始せず終了コード `1` とする。正規化前の推測順で build を開始してはならない。

**異常系：**

| 条件 | 処理 |
|------|------|
| priority 不正 | runner は対象 queue entry を処理しない。 |
| created_seq 欠落の旧 entry | GET 表示時は末尾扱いにする。runner 取り出し前または queue 更新時に `.build_state` lock 内で正規化保存する。正規化保存失敗時は build を開始しない。 |
| queue full | `429`。priority による上書き削除はしない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| urgent 後投入 | normal より先に処理。 |
| 同一 priority | FIFO。 |
| active あり | circuit closed では新規 urgent より active を先に再実行する。circuit open では active と urgent waiting の両方を保持する。 |
| 不正 priority | 状態差分なしで `422`。 |
| created_seq 欠落 | 正規化保存後に順序判定し、正規化保存失敗なら build なし。 |
| urgent queue full | `429`、既存 low entry も削除しない。 |

**priority queue 実装確認固定契約：**

| 項目 | 合格条件 |
|------|----------|
| 採番 | queue 追加時に `.build_state.queued[].created_seq` 最大値 + 1 を保存し、priority 省略時は `"normal"` を保存する。 |
| 表示順 | runner が読む queue 順は priority 数値昇順、created_seq 昇順、id 昇順とする。 |
| 取り出し順 | circuit closed かつ active がなければ runner は表示順と同じ順序で 1 件だけ waiting から active へ移す。active があれば waiting を移動しない。circuit open ではどちらも移動しない。 |
| 旧 entry | `created_seq` 欠落 entry は runner 取り出し前に lock 内で正規化保存し、保存失敗時は build を開始しない。 |
| clear 連携 | runner は API clear 後の waiting entry を再生成せず、保持済み active entry を waiting へ戻さない。 |
| queue full | priority が高くても既存 entry を削除せず `429` とする。 |
| 確認条件 | fixture は urgent 優先、同一 priority FIFO、created_seq 正規化、不正 priority、queue full、clear をすべて固定する。 |

<a id="sec-27-36"></a>
**27.36 失敗原因の自動分類：**

[`docs/details/runner.md` 詳細本文責務 §27.36](runner.md#sec-27-36) の境界は owner component `runner`、collaborator component `api`、`statefile` とする。

本機能の目的は、build failure と deploy pending の原因を固定カテゴリへ分類し、調査開始点を build log、history、UI に残すことである。

<a id="runner-failure-category-map"></a>
**分類値：**

| category | 判定条件 |
|----------|----------|
| `github_api` | GitHub API 最終失敗、rate limit 復帰不可、認証失敗。 |
| `pipeline_timeout` | build step timeout。 |
| `pipeline_exit` | build command exit code 非 0。 |
| `deploy_failure` | SSH 転送、remote checksum、pending transfer 発生。`target_status="success_deploy_pending"` でもこの値を保存する。 |
| `hook_error` | pre hook abort または hook timeout。 |
| `config_error` | 設定 parse、validation、必須値不足。 |
| `resource_error` | disk 不足、binary 不在、権限エラー。 |
| `unknown` | [分類値表](#runner-failure-category-map) に列挙した他の category に該当しない failure。 |

**正常系：**

1. 正規化状態 `failure` または `target_status="success_deploy_pending"` の確定時に、分類優先順位で category を決定する。
2. `.build_logs/{id}.json.failure_category` と `failure_evidence[]` を保存する。
3. `.build_history.failure_category` に同じ値を保存する。
4. runner は保存済み category を API が read-only で参照できる形で保存する。

**`failure_evidence[]` 保存境界：**

FailureCategory 値と FailureEvidence object の key、型、列挙値は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) を正とする。runner は分類条件と evidence の選択だけを担当し、schema を追加または別名化してはならない。

分類時は `failure_evidence[]` を最大 10 件まで保存する。10 件を超える場合は分類に使った evidence を先頭に残し、残りは発生順で 9 件まで保存する。

**API filter 参照：**

`GET /api/history` の `failure_category` query、`status="success"` 履歴の `null`、`status="success_deploy_pending"` 履歴の `deploy_failure`、未知 query の `422`、既存未知値と `warnings` の扱いは [`docs/details/api.md`](api.md) 詳細本文責務の `HistoryPageObject` 契約を正本とする。runner は HTTP query、status、response を定義または生成しない。

**異常系：**

| 条件 | 処理 |
|------|------|
| 複数カテゴリ該当 | [分類値表](#runner-failure-category-map) の上から最初に該当した category を採用する。 |
| evidence 抽出不能 | category は保存し、evidence は空配列。 |
| 既存 history に未知 category | runner は新規保存で未知 category を作らない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| timeout | `pipeline_timeout`。 |
| SSH 失敗 | `deploy_failure`。 |
| 設定不正 | `config_error`。 |
| filter | category 指定で該当履歴だけ返る。 |
| evidence 上限超過 | 最大 10 件で保存され、secret 平文を含まない。 |
| history `status="success"` | `failure_category:null`。 |
| history `status="success_deploy_pending"` | API の正規化状態は `success`、`failure_category:"deploy_failure"`。 |

**failure classification 実装確認固定契約：**

| 項目 | 合格条件 |
|------|----------|
| 分類優先順位 | 複数カテゴリに該当する場合、[分類値表](#runner-failure-category-map) の上から最初に該当した category を保存する。 |
| 保存先一致 | `.build_logs/{id}.json.failure_category`、`failure_evidence[]`、`.build_history.failure_category` が同じ分類結果を参照する。 |
| evidence 安全性 | `failure_evidence[].message` は最大 300 文字で、secret、token、path 全体、入力値連結を含めない。 |
| 上限 | evidence は最大 10 件とし、分類に使った evidence を先頭に残す。 |
| API 連携 | API が schema-valid 保存値をそのまま参照でき、runner は API filter のために category を書き換えない。 |
| 成功系履歴 | history `status="success"` は `failure_category:null`、history `status="success_deploy_pending"` は `failure_category:"deploy_failure"` とし、deploy 障害を filter 対象から除外しない。 |
| 確認条件 | fixture は timeout、deploy failure、config error、filter、unknown warning、evidence mask、上限超過をすべて固定する。 |

<a id="sec-27-37"></a>
**27.37 ビルド実行環境の記録：**

[`docs/details/runner.md` 詳細本文責務 §27.37](runner.md#sec-27-37) の境界は owner component `runner`、collaborator component `statefile` とする。

本機能の目的は、build 時点の実行環境を記録し、後から再現性と障害原因を確認できるようにすることである。

**記録先：**

`.build_logs/{id}.json.environment` に [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の Environment object を保存する。runner は取得時点、取得手段、継続可否、secret 非保存だけを担当し、Environment object の key、型、nullable 条件を再定義してはならない。

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
| 保存失敗 | environment 保存失敗は build を開始せず、`.build_status.json` に `status="failure"`、`last_target_status="failure_state_write"` を保存し、終了コード `1`。 |

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

**environment record 実装確認固定契約：**

| 項目 | 合格条件 |
|------|----------|
| 取得時点 | build id 採番直後、builder 起動前に environment snapshot を取得して保存する。 |
| secret 非保存 | 環境変数 value、token、secret、PATH 全体、VCS revision を保存しない。 |
| version 取得 | builder version は 2 秒 timeout、stdout 先頭行の最初の token だけを保存し、失敗時は `"unknown"`。 |
| path 境界 | home 配下 state dir は basename だけ、それ以外は絶対 path を保存する。 |
| 保存失敗 | environment 保存失敗時は build 本体を起動せず、`.build_status.json` に `status="failure"`、`last_target_status="failure_state_write"` を保存する。 |
| 継続可能失敗 | hostname、disk stat、builder version 取得不能は WARN または unknown/null とし、build を継続する。 |
| 確認条件 | fixture は通常保存、builder version timeout、home path 短縮、secret 非保存、environment write failure をすべて固定する。 |

<a id="sec-27-38"></a>
**27.38 ビルド所要時間の異常検知：**

[`docs/details/runner.md` 詳細本文責務 §27.38](runner.md#sec-27-38) の境界は owner component `runner`、collaborator component `api`、`statefile` とする。

本機能の目的は、過去 trend と比較して異常に遅い build を検出し、性能劣化を WARN、history flag、通知で可視化することである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 schema | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.server_config.duration_anomaly` と DurationAnomaly object。型、許容範囲、既定値は同 schema を正とする。 |
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
| notify payload | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の NotificationPayload object とし、`event="duration_anomaly"` の event 別 key を設定する。 |
| 設定 validation | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の DurationAnomaly object schema に一致する場合だけ判定を有効にする。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| trend 破損 | [`docs/details/runner.md` 詳細本文責務 §27.33](runner.md#sec-27-33) の復旧後、復旧できた場合だけ判定する。 |
| avg / p95 が null | 判定しない。 |
| 通知失敗 | build status は変更せず [`docs/details/runner.md` 詳細本文責務 §27.32](runner.md#sec-27-32) の notify retry 契約に従う。 |
| 設定値不正 | runner は既定値ではなく機能無効として扱う。 |
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

**duration anomaly 実装確認固定契約：**

| 項目 | 合格条件 |
|------|----------|
| 判定順 | 今回 sample を `.build_trends.json` に追加する前の summary で判定し、判定後に trend を更新する。 |
| 対象 status | `"success"` と `"success_deploy_pending"` だけを anomaly 判定対象とし、failure は trend sample 追加だけを行う。 |
| 閾値 | avg 超過、p95 超過、両方超過を区別し、`threshold_source` を固定値で保存・通知する。 |
| history 更新 | anomaly 時は WARN、`flagged=true`、`tags += ["duration_anomaly"]` を保存し、重複 tag を作らない。 |
| 通知独立性 | duration_anomaly 通知失敗は build status を変更せず、[`docs/details/runner.md` 詳細本文責務 §27.32](runner.md#sec-27-32) の notify retry 契約に従う。 |
| 設定不正 | runner は既定値へ補正せず機能無効として扱う。 |
| 確認条件 | fixture は sample 不足 no-op、avg 超過、p95 境界、通知失敗、failure build、tag 重複、設定不正をすべて固定する。 |

<a id="sec-27-21-2"></a>
**[`docs/details/runner.md` 詳細本文責務 §27.21〜`docs/details/runner.md` 詳細本文責務 §27.38 runner / statefile 連動実装確認ゲート](runner.md#sec-27-21-2)：**

[`docs/details/runner.md` 詳細本文責務 §27.21](runner.md#sec-27-21)〜[`docs/details/runner.md` 詳細本文責務 §27.38](runner.md#sec-27-38) の runner owner 機能は、対象機能契約の確認条件に加えて [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) runner owner 機能横断ゲート固定表を満たす。runner owner 機能横断ゲートは runner が状態更新を呼び出す順序と失敗時境界を固定するものであり、状態ファイル schema 本文は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c)、fixture 本文は [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) を参照する。

| ゲート | 合格条件 | 禁止条件 |
|--------|----------|----------|
| state target | 対象の [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) 機能契約に列挙された状態ファイル、[`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) / [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の schema、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](../DETAIL_INDEX.md#0i-詳細節対応表) の collaborator だけを使用する。 | 未定義状態ファイル、未知 key、空 placeholder file、component 固有でない汎用 state file を追加する。 |
| write order | 通常 build は [`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) の状態更新順序、rollback は [`docs/details/runner.md` 詳細本文責務 §14b](runner.md#14b-スナップショット管理) の accepted 確定順と worker 書込順、追加機能は該当 owner 節の固定順を fixture でそれぞれ検証する。 | 異なる lifecycle の書込を 1 つの抽象順にまとめる、並列処理の完了順を永続保存順にする、fixture に `write_order` を持たない複数書込を確認済み扱いにする。 |
| partial failure | 失敗地点より後の write / external call / notification は実行しない。失敗地点より前に成功済みの状態は、対象の [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) 機能契約が rollback を明記しない限り戻さない。 | 保存済み build log / history / status を実装者判断で削除、巻き戻し、再分類する。 |
| no-op / skip | 変更なし、cooldown、tag 不一致、disabled、sample 不足、queue duplicate 等は、対象の [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) 機能契約が定める runner 結果、保存値、終了コード、副作用ゼロまたは指定最小副作用で固定する。HTTP response は api owner だけが定義する。 | no-op で config log、audit、notify、build log、history、SHA cache を暗黙更新する。 |
| dry-run | dry-run は設定検証、対象判定、差分判定、実行可否の出力だけを行い、状態ファイル、lock、log、history、cache、notification、external write を作成・更新しない。 | dry-run 結果を後続実行用 cache として保存する。 |
| secret mask | PAT、Webhook secret、SMTP password、branch env secret、hook output secret、remote credential は読込直後に mask 登録し、stdout / stderr / build log / history / status / pending / fixture expected に平文を残さない。 | mask 登録前にログ保存する、secret 長・hash・prefix・suffix を保存する。 |
| schema strictness | runner が保存する object は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の key だけを持ち、nullable、UTC 時刻、列挙値、配列順を満たす。 | SDK / UI 用の表示名、計算済み label、未定義 fallback key を state に保存する。 |
| fixture evidence | 対象 [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) 機能の fixture は success、failure、partial、no-op、idempotency、secret mask、corrupt state のうち該当するケースを持つ。 | 正常系だけの fixture で runner / statefile 連動を確認済み扱いにする。 |
