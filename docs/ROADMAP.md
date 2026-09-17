# Adlaire CI — Roadmap

本ファイルは、Adlaire CI の実装状態、実装可否、Phase、機能インベントリ、将来計画、昇格手順を管理するロードマップ正本である。

方針、ポリシー、正本関係、禁止事項、リリース判断は `docs/SPEC.md` を正とする。詳細仕様の入口、owner component 対応表、共通固定値、横断補足契約は `docs/DETAIL_INDEX.md` を正とする。owner component 別の入出力、状態、処理順序、異常系、検証条件は `docs/details/*.md` を正とする。

本ファイルは状態と計画を扱う。HTTP response schema、状態ファイル schema、SDK method の実装詳細、UI DOM、fixture assertion、具体的な処理順序は定義しない。

## 1. ロードマップ責務

| 管理対象 | 本ファイルで扱う内容 | 正本外の内容 |
|----------|----------------------|--------------|
| 実装状態 | component ごとの `実装済み`、`仕様化済み・未実装`、`改訂予定`、`将来計画` の分類。 | 実装コード本文、詳細仕様本文。 |
| Phase | Phase 順序、対象 owner component、依存条件、完了条件、引き継ぎ契約。 | 各 component の関数、endpoint、schema、DOM、fixture 本文。 |
| 機能インベントリ | 機能の分類、実装可否、詳細仕様参照先。 | API request / response の完全表、SDK method の完全実装条件。 |
| 将来計画 | 実装不可の構想、昇格条件、実装禁止条件。 | 実装可能な詳細仕様としての具体値。 |

## 2. 状態分類

| 状態 | 実装可否 | 意味 | 実装者の扱い |
|------|----------|------|--------------|
| 実装済み | 完了済み | ソースコード実装と必須検証が完了した項目。 | 仕様、実装、検証、索引の整合を維持する。 |
| 実装中・検証未完了 | 検証待ち | 実装に着手済みだが、必須検証または証跡が未完了の項目。 | 未完了項目を明記し、実装済みとして扱わない。 |
| 仕様化済み・未実装 | 実装可 | `docs/DETAIL_INDEX.md` と owner component 詳細仕様に実装可能な詳細が揃っている項目。 | Phase と詳細仕様に従って実装できる。 |
| 改訂予定 | 実装不可 | 将来計画から格上げ済みだが、詳細仕様作成中の項目。 | 詳細仕様が完了するまで実装しない。 |
| 将来計画 | 実装不可 | 構想として保持するが、実装契約がない項目。 | 実装対象にしない。 |
| 未仕様化 | 実装不可 | 本ロードマップと詳細仕様に存在しない項目。 | 推測で実装しない。 |

## 3. 実装状態

本ファイルでは、仕様化済みの内容と実装済みの内容を区別して扱う。

下表のコンポーネント名は、`docs/SPEC.md` Part 1 §4.3 の標準ソース配置に基づく。標準ソース配置への実装移行は完了済みであり、`main.go`、現存する `components/*.go`、`admin/` 配下の管理 UI 静的ファイル、`testdata/<component>/` を現行配置として扱う。

| コンポーネント | 状態 | 備考 |
|---------------|------|------|
| `components/builder.go` | 実装済み | Go 版 Markdown → 静的 Web サイトビルドスクリプトとして Phase 1 の `gofmt` と `go test` 検証済み。 |
| `components/runner.go` | 実装済み | Go 版 CI ランナーとして Phase 2 完了判定パスの `gofmt` と `go test` 検証済み。 |
| `components/api.go` | 実装済み | Go 版管理 API サーバー。API 完全契約表の endpoint、認証、状態ファイル、管理操作、検証テストを実装済み。 |
| `admin/adlaire-ci-sdk.js` | 実装済み | 管理ツール用 JavaScript SDK。Phase 5 の SDK class、method、HTTP error、token 破棄、query / body 生成を実装済み。 |
| `admin/index.html` | 実装済み | 標準管理ツール UI。Phase 6 の DOM id、panel、SDK 呼び出し、表示状態、秘密情報消去を実装済み。 |
| `components/mcp.go` | 将来計画 | Go 版 MCP サーバー。将来計画として管理し、実装済みとは扱わない。 |

仕様化済み・未実装、または将来計画の項目を、実装済み機能として扱ってはならない。

## 4. Phase 実装計画

## 4.1 初期実装 Phase 分割

Go 版初期実装は、`docs/SPEC.md` §0e の対象範囲を一括実装せず、下表の Phase 順に進める。上位 Phase の完了判定を満たす前に、下位 Phase の実装 PR を開始してはならない。

実装単位、実装 PR 単位、完了判定単位は Phase のみとする。`P0`、`P1`、`P2〜P5` などの優先度ラベル、抽象段階、API 内部段階名を実装単位として使ってはならない。API の実装範囲は、Phase 3 を「API 基盤・認証・状態 read/write・運用基本操作」、Phase 4 を「API 拡張運用操作」として扱う。

各 Phase の `対象` は、その Phase の owner component を示す。状態ファイル、security、archive、commitstatus、admin、fixture、setup が関わる場合も、それらは collaborator component として該当 Phase の完了条件に含める。collaborator component の詳細仕様に未充足がある場合は、owner component の実装で補完せず、先に該当する責務 component 別詳細仕様ファイルを改訂する。

| Phase | 対象 | 実装範囲 | 依存条件 | 完了条件 |
|-------|------|----------|----------|----------|
| Phase 1 | `builder` | §2〜§9 の CLI、Markdown 変換、静的 Web サイト出力、テーマコンポーネント、生成物確認。 | なし。 | §0e の `builder` 必須検証と §0f の `builder` 完了判定を満たす。 |
| Phase 2 | `runner` | §10〜§20 の CI ランナー、GitHub API 連携、SHA キャッシュ、pipeline 起動、SSH 転送、snapshot、通知、ログ、systemd / setup 参照境界。systemd unit 本文と配置手順は `docs/details/setup.md` §26 を正とする。 | Phase 1 が完了し、`adlaire-ci-build` の CLI 契約が固定されている。 | §0e の `runner` 必須検証と §0f の `runner` 完了判定を満たす。 |
| Phase 3 | `api` | `docs/details/api.md` §21〜§22、§21a、§25 と `docs/details/setup.md` §26 のうち、認証、セッション、共通エラー、状態ファイル読み書き、ビルド操作、status、logs、history、queue、circuit breaker。 | Phase 2 が完了し、runner が書き込む状態ファイル schema が固定されている。 | API 基盤・認証・状態 read/write・運用基本操作の必須検証、§0e の `api` API 共通・状態ファイル検証、§0f の `api` 完了判定の該当範囲を満たす。 |
| Phase 4 | `api` | `docs/details/api.md` §22.0e、§22.0f のうち、config、repo、branch、schedule、notify、snapshot、rollback、maintenance、access control、hooks、tokens 等の拡張運用操作。 | Phase 3 が完了し、API 共通処理と認証が固定されている。 | API 拡張運用操作の必須検証と §0e の `api` endpoint 契約を満たす。 |
| Phase 5 | `sdk` | `docs/details/sdk.md` §23 の SDK class、method、戻り値、HTTP error、token 破棄、query / body 生成。 | Phase 3 と Phase 4 が完了し、`docs/details/api.md` §22.0e の endpoint 契約が固定されている。 | §0e の `sdk` 契約と §0f の `sdk` 完了判定を満たす。 |
| Phase 6 | `ui` | `docs/details/ui.md` §24 の標準管理ツール UI、DOM id、panel、操作、SDK 呼び出し、成功表示、失敗表示、disabled、再取得、秘密情報消去。 | Phase 5 が完了し、SDK method 契約が固定されている。 | §0e の `ui` 契約と §0f の `ui` 完了判定を満たす。 |

### 4.1.1 Phase 1 完全仕様ゲート（`builder`）

Phase 1 の実装詳細本文は `docs/details/builder.md` を正とする。本ファイルでは、Phase 1 の対象、依存条件、完了条件、後続 Phase への引き継ぎ確認だけを扱う。

| 確認 | 参照先 |
|------|--------|
| CLI、Markdown 変換、静的 Web サイト出力、theme component、生成物確認 | `docs/details/builder.md` §1〜§9、§8a |
| Phase 1 fixture、testdata、PR 証跡 | `docs/details/fixture.md` §0g.8-F |
| release / setup 受け入れ条件 | `docs/details/setup.md` §26.7 |

### 4.1.2 Phase 2 完全仕様ゲート（`runner`）

Phase 2 の実装詳細本文は `docs/details/runner.md` を正とする。本ファイルでは、Phase 2 が Phase 1 の `adlaire-ci-build` 契約に依存し、後続 API が読む runner 状態契約を固定することだけを扱う。

| 確認 | 参照先 |
|------|--------|
| CI runner、GitHub API 連携、SHA cache、pipeline、deploy、snapshot、通知、systemd | `docs/details/runner.md` §10〜§20、§15a |
| runner fixture、fake GitHub、fake ssh / notifier、PR 証跡 | `docs/details/fixture.md` §0g.8-F |
| release / setup 受け入れ条件 | `docs/details/setup.md` §26.7 |

### 4.1.3 Phase 3 完全仕様ゲート（`api`）

Phase 3 の実装詳細本文は `docs/details/api.md` と `docs/details/setup.md` を正とする。本ファイルでは、API 共通契約、認証、状態 read/write、運用基本 endpoint が Phase 4〜6 の前提になることだけを扱う。

| 確認 | 参照先 |
|------|--------|
| API 共通処理、API server 制限、認証、運用基本 endpoint、状態 read/write | `docs/details/api.md` §21〜§22、§21a、§25、§22.0f |
| 状態ファイル schema、lock、atomic write | `docs/details/statefile.md` §22.0a、§22.0c |
| 管理 API 導入、systemd、release 受け入れ条件 | `docs/details/setup.md` §26 |

### 4.1.4 Phase 4 完全仕様ゲート（`api`）

Phase 4 の実装詳細本文は `docs/details/api.md` を正とする。本ファイルでは、拡張運用 endpoint が SDK / UI の最終入力契約になることだけを扱う。

| 確認 | 参照先 |
|------|--------|
| 拡張運用 endpoint、request / response、error、auth、secret mask | `docs/details/api.md` §22.0e、§22.0f |
| security 連携 | `docs/details/security.md` §27.42〜§27.47 |
| API fixture、endpoint 証跡 | `docs/details/fixture.md` §0g.8-F、§22-F、§27-F |

### 4.1.5 Phase 5 完全仕様ゲート（`sdk`）

Phase 5 の実装詳細本文は `docs/details/sdk.md` を正とする。本ファイルでは、SDK が固定済み API endpoint だけを呼び、UI 表示判断を持たないことだけを扱う。

| 確認 | 参照先 |
|------|--------|
| SDK class、method、HTTP 対応、query / body 生成、error、stream、token 破棄 | `docs/details/sdk.md` §23 |
| API endpoint 対応 | `docs/details/api.md` §22.0e |
| SDK fixture、fake fetch / stream | `docs/details/fixture.md` §0g.8-F |

### 4.1.6 Phase 6 完全仕様ゲート（`ui`）

Phase 6 の実装詳細本文は `docs/details/ui.md` を正とする。本ファイルでは、UI が SDK 経由だけで API と通信し、秘密情報を DOM に残さないことだけを扱う。

| 確認 | 参照先 |
|------|--------|
| DOM id、panel、操作、表示状態、SDK 呼び出し、秘密情報消去 | `docs/details/ui.md` §24 |
| SDK method 契約 | `docs/details/sdk.md` §23 |
| UI fixture、fake SDK、直接 API 呼び出し禁止確認 | `docs/details/fixture.md` §0g.8-F |

### 4.1.7 Phase 間引き継ぎ契約

各 Phase の完了時は、次 Phase が依存する契約を変更不可として扱う。後続 Phase で変更が必要になった場合は、後続 Phase の実装で吸収せず、契約を定義した owner 詳細仕様へ戻す。

| 引き継ぎ元 | 引き継ぎ先 | 固定する契約 | 正本 |
|------------|------------|--------------|------|
| Phase 1 | Phase 2 | `adlaire-ci-build` CLI、終了コード、stdout / stderr、`[REPORT]`、出力サイト構造。 | `docs/details/builder.md` |
| Phase 2 | Phase 3 | runner 状態ファイル schema、lock、history/log、pending queue、circuit breaker、snapshot、通知ログ。 | `docs/details/runner.md` / `docs/details/statefile.md` |
| Phase 3 | Phase 4 | API 共通契約、認証、session、error body、validation、lock error、SSE 基本形式。 | `docs/details/api.md` |
| Phase 4 | Phase 5 | 全 endpoint の method、path、query、request、response、error、認証要否。 | `docs/details/api.md` |
| Phase 5 | Phase 6 | SDK method 名、引数、戻り値、error object、stream handle、token 破棄条件。 | `docs/details/sdk.md` |
| Phase 6 | 初期実装完了 | UI 操作、表示状態、secret 消去、SDK 経由通信、実装完了検証結果。 | `docs/details/ui.md` |

### 4.1.8 Phase 別 実装 PR 成果物チェックリスト

Phase fixture / testdata 配置、fake 実装、実装 PR 証跡の詳細は `docs/details/fixture.md` §0g.8-F を正とする。本ファイルでは、Phase ごとの成果物参照先だけを保持し、fixture 名、expected / effects、fake 動作、PR 証跡項目を重複定義しない。

| Phase | 実装対象 | 成果物・fixture 正本 | 受け入れ条件 |
|-------|----------|----------------------|--------------|
| Phase 1 | `builder` | `docs/details/fixture.md` §0g.8-F | `docs/details/builder.md` と `docs/details/setup.md` §26.7 を満たす。 |
| Phase 2 | `runner` | `docs/details/fixture.md` §0g.8-F | `docs/details/runner.md` と `docs/details/setup.md` §26.7 を満たす。 |
| Phase 3 | `api` | `docs/details/fixture.md` §0g.8-F、§27-F | `docs/details/api.md` の API 基盤・認証・状態 read/write・運用基本操作と setup API 導入条件を満たす。 |
| Phase 4 | `api` | `docs/details/fixture.md` §0g.8-F、§27-F | `docs/details/api.md` の API 拡張運用操作を満たす。 |
| Phase 5 | `sdk` | `docs/details/fixture.md` §0g.8-F | `docs/details/sdk.md` §23 を満たす。 |
| Phase 6 | `ui` | `docs/details/fixture.md` §0g.8-F | `docs/details/ui.md` §24 を満たす。 |

---


## 5. 機能インベントリと統合ロードマップ

## 5.1 機能一覧

本ロードマップが管理する機能インベントリである。各機能の仕様詳細は `docs/DETAIL_INDEX.md` と owner component 別の `docs/details/*.md` を参照する。実装状態は本ファイル §2 と §3、および `docs/SPEC.md` Part 2 §0a の仕様成熟度ポリシーに従って判定する。

### ビルド・CI ランナー（components/runner.go）

**Go 版で実装済みの初期範囲：**

- `--state-dir`、`--once`、`--version`、`--help` の CLI 契約
- `.branch_config` と `.last_sha` による branch target / SHA cache 読み込み
- GitHub Trees API / Blobs API による対象 Markdown 取得
- SHA 一致時の変更なし skip
- `pipeline.sh` 起動、stdout/stderr 収集、`[REPORT]` / `[WARN]` 取り込み
- `.build_logs/{id}.json`、`.build_history`、`.build_status.json`、`.build_state`、`.build_lock` の作成・更新
- deploy 失敗時の `.pending_transfers` 追加
- `.notify_pending` 破損時の退避と `[]` 再生成
- GitHub API retry / rate limit 待機
- `adlaire-ci-build` 実行可否と disk 空き容量の precheck
- `.snapshots/{build_id}/site` 保存と世代 pruning
- `.notify_pending` の HTTP 再送と成功時削除
- `.build_circuit_state` による circuit open skip と失敗回数記録
- `.pending_transfers` の起動時再試行と成功時削除
- SSH 転送時の checksum 比較、未変更ファイル skip、転送後 checksum 検証
- 複数 branch target の順次処理
- `BUILD_COOLDOWN_SECONDS` による cooldown skip
- `FORCE_BUILD_INTERVAL` による変更なし時の定期強制ビルド
- GitHub commits API によるトリガー commit 情報の build log 記録
- GitHub PAT 有効期限ヘッダーの 7 日以内 WARN ログ
- `.notify_config` に基づく成功・失敗・転送失敗 Webhook 通知送信
- `OUTPUT_SIZE_WARN_MB` による出力サイトサイズ警告
- Phase 2 fixture R1〜R7 と完了判定パステスト

**Go 版で仕様化済みの全体範囲：**

- GitHub リポジトリの対象ファイルを定期ポーリング（systemd timer）
- blob SHA による差分検出（変更なし時はビルドをスキップ）
- Markdown → 静的 Web サイト変換（`adlaire-ci-build` を呼び出し）
- ビルド成功後に SHA キャッシュを更新する
- ビルド失敗時は SHA キャッシュを更新せず、次回起動時に再試行可能な状態を維持する
- systemd oneshot ユニットとして動作（`adlaire-ci.service`）

- ビルド結果を `.build_history` に記録（ID 形式：`b{YYYYMMDDHHmmss}`）
- ビルドごとのログを `.build_logs/{id}.json` に保存
- ビルド成功・失敗時に Webhook 通知を送信（`.notify_config` を読み込み送信。送信責務は `components/runner.go`。通知 API は設定の読み書きのみ）
- ビルド成功後、出力サイトディレクトリを SSH 経由で静的コンテンツ配信サーバーへ転送する（差分転送・`DEPLOY_TARGETS` 複数先対応 → §14a）
- SSH 転送失敗時は `.pending_transfers` へキューイングし、次回起動時に自動再試行する（→ §14a ペンディングキュー）
- 転送成功後、出力サイトディレクトリを `.snapshots/` へアーカイブし `HISTORY_KEEP_N` 世代を超過分から自動削除する（→ §14b）
- SSH 転送失敗時に Webhook 通知を送信する（`on: ["deploy_failure"]` 設定時）
- `BRANCH_TARGETS` リストで複数ブランチを順次ポーリング・ビルドする（→ §12 設定値）
- `FORCE_BUILD_INTERVAL` 設定時、前回ビルドから指定時間経過で変更なしでも強制ビルドする

### 管理 API エンドポイント（components/api.go）

| カテゴリ | エンドポイント |
|---------|--------------|
| 認証 | `POST /api/login` / `POST /api/logout` / `POST /api/change-password` |
| 死活監視 | `GET /api/health` |
| ビルド操作 | `POST /api/build` / `POST /api/build/force` / `POST /api/build/cancel` / `GET /api/build/stream` / `POST /api/circuit-breaker/reset` |
| ステータス | `GET /api/status` / `GET /api/dashboard` |
| ログ | `GET /api/logs` / `GET /api/logs/export` / `GET /api/logs/search` / `POST /api/logs/cleanup` |
| ビルド履歴 | `GET /api/history` / `GET /api/history/{id}/log` / `GET /api/history/{id}/comment` / `POST /api/history/{id}/comment` / `GET /api/history/export` / `POST /api/history/{id}/flag` / `POST /api/history/{id}/tags` / `POST /api/history/{id}/rollback` |
| スケジュール | `GET /api/schedule` / `POST /api/schedule/interval` / `POST /api/schedule/pause` / `POST /api/schedule/resume` / `POST /api/schedule/allowed-hours` / `POST /api/schedule/force-interval` / `POST /api/schedule/cooldown` |
| 通知 | `GET /api/notify-config` / `POST /api/notify-config` / `POST /api/notify-test` / `GET /api/notify-log` / `POST /api/notify/weekly-summary` |
| システム情報 | `GET /api/sysinfo` / `GET /api/output-meta` / `GET /api/diagnostics` / `GET /api/rate-limit` / `GET /api/disk-usage` |
| 統計 | `GET /api/stats` / `GET /api/stats/timeline` / `GET /api/stats/build-duration` |
| リポジトリ | `GET /api/repo-info` / `POST /api/repo-config` / `GET /api/branch-config` / `POST /api/branch-config` |
| PAT 管理 | `GET /api/pat-status` / `POST /api/pat-verify` / `POST /api/pat-update` |
| 設定 | `GET /api/config` / `POST /api/config` / `POST /api/log-level` / `GET /api/config-log` |
| アクセスログ | `GET /api/access-log` |
| バックアップ | `GET /api/backup` / `POST /api/restore` |
| セッション管理 | `GET /api/sessions` / `POST /api/sessions/revoke-all` |
| API トークン | `GET /api/tokens` / `POST /api/tokens` / `DELETE /api/tokens/{id}` |
| スナップショット | `GET /api/snapshots` / `GET /api/snapshots/{id}/download` / `DELETE /api/snapshots/{id}` |
| メンテナンス | `GET /api/maintenance` / `POST /api/maintenance/enable` / `POST /api/maintenance/disable` |
| アクセス制御 | `GET /api/access-control` / `POST /api/access-control` |
| フック | `GET /api/hooks` / `POST /api/hooks` / `DELETE /api/hooks/{id}` / `GET /api/hooks/{id}/log` |
| アラートルール | `GET /api/alert-rules` / `POST /api/alert-rules` / `DELETE /api/alert-rules/{id}` |
| Webhook 受信 | `POST /api/webhook` / `GET /api/webhook-events` / `GET /api/webhook-config` / `POST /api/webhook-config` |
| 自動タグ付け | `GET /api/tag-rules` / `POST /api/tag-rules` / `DELETE /api/tag-rules/{id}` |
| パイプライン | `GET /api/pipeline-config` / `POST /api/pipeline-config` / `POST /api/verify-output` |
| 運用ノート | `GET /api/notes` / `POST /api/notes` |
| メール通知 | `GET /api/smtp-config` / `POST /api/smtp-config` / `POST /api/smtp-test` |
| ビルドキュー | `GET /api/queue` / `DELETE /api/queue` |
| ダッシュボードレイアウト | `GET /api/dashboard-layout` / `POST /api/dashboard-layout` |

### SDK メソッド（adlaire-ci-sdk.js）

| カテゴリ | メソッド |
|---------|--------|
| 認証 | `login()` / `logout()` / `changePassword()` |
| ビルド操作 | `triggerBuild()` / `buildForce()` / `cancelBuild()` / `streamBuild()` / `resetCircuitBreaker()` |
| ステータス | `getStatus()` / `getDashboard()` |
| ログ | `getLogs()` / `exportLogs()` / `searchLogs()` / `cleanupLogs()` |
| ビルド履歴 | `getHistory()` / `getHistoryLog()` / `getHistoryComment()` / `setHistoryComment()` / `exportHistory()` / `setHistoryFlag()` / `setHistoryTags()` / `rollbackHistory()` |
| スケジュール | `getSchedule()` / `setScheduleInterval()` / `pauseSchedule()` / `resumeSchedule()` / `setAllowedHours()` / `clearAllowedHours()` / `setForceInterval()` / `setBuildCooldown()` |
| Webhook 受信 | `getWebhookEvents()` / `getWebhookConfig()` / `setWebhookConfig()` |
| 通知 | `getNotifyConfig()` / `setNotifyConfig()` / `notifyTest()` / `getNotifyLog()` / `notifyWeeklySummary()` |
| システム情報 | `getSysinfo()` / `getOutputMeta()` / `getDiagnostics()` / `getRateLimit()` / `getDiskUsage()` |
| 統計 | `getStats()` / `getStatsTimeline()` / `getStatsBuildDuration()` |
| リポジトリ | `getRepoInfo()` / `setRepoConfig()` / `getBranchConfig()` / `setBranchConfig()` |
| PAT 管理 | `getPatStatus()` / `patVerify()` / `updatePat()` |
| 設定 | `getConfig()` / `setConfig()` / `setLogLevel()` / `getConfigLog()` |
| アクセスログ | `getAccessLog()` |
| バックアップ | `backup()` / `restore()` |
| セッション管理 | `getSessions()` / `revokeAllSessions()` |
| API トークン | `getTokens()` / `createToken()` / `revokeToken()` |
| スナップショット | `getSnapshots()` / `downloadSnapshot()` / `deleteSnapshot()` |
| メンテナンス | `getMaintenance()` / `enableMaintenance()` / `disableMaintenance()` |
| アクセス制御 | `getAccessControl()` / `setAccessControl()` |
| フック | `getHooks()` / `addHook()` / `deleteHook()` / `getHookLog()` |
| アラートルール | `getAlertRules()` / `addAlertRule()` / `deleteAlertRule()` |
| 自動タグ付け | `getTagRules()` / `addTagRule()` / `deleteTagRule()` |
| パイプライン | `getPipelineConfig()` / `setPipelineConfig()` / `verifyOutput()` |
| 運用ノート | `getNotes()` / `setNotes()` |
| メール通知 | `getSmtpConfig()` / `setSmtpConfig()` / `smtpTest()` |
| ビルドキュー | `getQueue()` / `clearQueue()` |
| ダッシュボードレイアウト | `getDashboardLayout()` / `setDashboardLayout()` |
| 死活監視 | `health()` |

ES Module・外部依存なし。全メソッドは `Promise` を返す。`streamBuild` は SSE 接続確立後に `Promise<StreamHandle>` として resolve し、`StreamHandle` は `{ close(): void, closed: boolean }` を持つ。`constructor` を除く合計は 96 メソッド。

### 標準管理ツール パネル（admin/index.html）

| パネル | 主な機能 |
|-------|---------|
| ログイン | パスワード認証 |
| パスワード変更 | 強制変更フロー対応（5 回目以降は他パネルを非表示） |
| ステータス | 最終ビルド情報・出力サイトリンク・メンテナンスバナー表示（モード中） |
| 手動実行 | ビルド起動・強制ビルド・キャンセル・キュー状態表示・キューのクリア |
| ログビューア | ログ閲覧・キーワードフィルター・ログレベルフィルター・JSON エクスポート・横断検索（期間指定） |
| ビルド履歴 | 過去ビルド一覧・タグ列・フラグ列・ページネーション・タグ/フラグフィルター・JSON エクスポート・個別ログ参照 |
| システム情報 | ファイルサイズ・稼働時間・ディスク使用量・PAT 検証・PAT 更新フォーム・PAT 有効期限表示・GitHub API レート制限表示 |
| 通知設定 | 複数 Webhook 設定・ペイロードテンプレート編集・テスト送信・定期サマリー設定・送信履歴・メール通知設定（SMTP連携）・Webhook 署名シークレット設定 |
| 設定 | 保持行数・件数・タイムアウト・ログレベル変更・ログ保持期間（日数）・手動クリーンアップ・設定変更履歴・キュー最大サイズ設定・スナップショット保持世代数設定 |
| アクセスログ | ログイン履歴（日時・成否） |
| 統計 | 成功率・平均/最大ビルド時間・時系列グラフ |
| リポジトリ情報 | 監視設定確認・ポーリング間隔変更フォーム・ポーリング一時停止/再開・許可時間帯設定・強制再ビルド間隔設定・Webhook 受信 Secret 設定・ブランチターゲット設定（`.branch_config` 編集） |
| セッション管理 | セッション一覧・全セッション強制無効化 |
| システム診断 | PAT・GitHub API・ファイル・systemd・Webhook・出力整合性 一括診断・アラートバッジ表示・メンテナンスバナー表示 |
| ビルド比較 | ビルド履歴から 2 件を選択してログを並列差分表示 |
| API トークン管理 | 読み取り専用トークンの発行・一覧・失効 |
| スナップショット | ビルド成果物の世代一覧・ダウンロード・削除・ロールバック |
| メンテナンス | メンテナンスモードの有効化・解除・状態表示 |
| アクセス制御 | 許可 IP / CIDR 一覧・追加・削除 |
| フック | Pre/Post ビルドフック設定・実行ログ確認 |
| 運用ノート | Markdown 記述の運用メモ閲覧・編集 |

---

## 5.2 拡張ポイント・将来計画

Adlaire CI の実装済み項目、実装中・検証未完了項目、仕様化済み・未実装項目、改訂予定項目、将来計画項目を統合管理する。
仕様化する際は `docs/SPEC.md` の方針・ポリシー、`docs/DETAIL_INDEX.md` の対応表、該当 owner component の詳細仕様本文への追記を先行させる。

将来計画は、実装対象ではない。将来計画内の「検討」「予定」「候補」「推奨」は、実装可能な仕様を意味しない。将来計画を実装対象にする場合は、先に対象項目を `改訂予定` へ昇格し、`docs/DETAIL_INDEX.md` の対応表と該当 owner component の詳細仕様本文に実装可能な詳細仕様を追加したうえで `仕様化済み` とする。

**担当領域：** `CI ランナー` / `管理ツール・API` / `MCP サーバー` / `ビルドスクリプト`

### 5.2.1 統合ロードマップの読み方

本章の項目は、単一の統合ロードマップ表で管理する。実装可否は `状態` 列で判断し、担当領域や機能名だけで実装対象と判断してはならない。

| 状態 | 実装可否 | 意味 | 次アクション |
|------|----------|------|--------------|
| 実装済み | 完了済み | ソースコード実装と検証が完了した項目。 | 実装ファイル、検証結果、`docs/DOCUMENT_INDEX.md` を維持する。 |
| 実装中・検証未完了 | 検証待ち | ソースコード実装に着手済みだが、必須検証が未完了の項目。 | 必須検証を実行し、不足が残る場合は `実装済み` へ移動しない。 |
| 改訂予定 | 実装不可 | 将来計画から格上げ済みだが、詳細仕様作成中の項目。 | `docs/DETAIL_INDEX.md` の対応表と該当 owner component の詳細仕様本文を作成し、仕様化条件を満たす。 |
| 仕様化済み・未実装 | 実装可 | 正本仕様と詳細仕様があり、実装対象として扱える項目。 | `docs/DETAIL_INDEX.md` §0h・§0i.1〜§0i.4 と該当 owner component の詳細仕様本文を確認して実装する。 |
| 将来計画 | 実装不可 | `components/mcp.go` など、将来構想として管理する項目。 | 本節 5.2.3 の手順で `改訂予定` へ昇格する。 |

将来計画、改訂予定の項目は、実装着手可能な仕様ではない。実装対象にする場合は、先に 5.2.3 の手順で `仕様化済み・未実装` へ昇格させる。

---

### 5.2.2 統合ロードマップ表

本表は、`docs/ROADMAP.md` §5.2 の全項目を状態別に統合した唯一の一覧である。項目を追加、削除、昇格、実装完了する場合は、本表の `状態`、`実装可否`、`次アクション` を同時に更新する。

MCP サーバー領域の行は、現時点ではすべて将来構想例であり、実装契約、API 契約、状態ファイル契約、起動手順、検証条件を定義しない。`components/mcp.go`、MCP tools、MCP resources、MCP prompts、HTTP SSE transport、MCP audit / stats / config CRUD は、MCP 専用詳細仕様を新設し、`改訂予定` を経て `仕様化済み・未実装` へ昇格するまで実装してはならない。

| 状態 | 実装可否 | 担当領域 | 機能 | 概要 | 次アクション |
|------|----------|----------|------|------|--------------|
| 実装済み | 完了済み | CI ランナー | Phase 2 完了判定パス | GitHub API polling、SHA 差分検出、ビルド起動、ログ、履歴、snapshot、lock、precheck、retry、rate limit、circuit breaker、通知、転送、cooldown、force interval、commit info、PAT 期限警告、出力サイズ警告を `components/runner.go` で実装済み。 | Go test で Phase 2 fixture、hardening、完了判定パスを検証済み。 |
| 改訂予定 | 実装不可 | 全領域 | （なし） | 現時点で、将来計画から格上げ済みの仕様作成中項目はない。 | 格上げ時に元状態、格上げ日、整理順序（実装単位ではない）、詳細仕様作成先を概要へ記録する。 |
| 実装済み | 完了済み | CI ランナー | ビルドタイムアウト | Go 標準ライブラリ `context.WithTimeout` と `os/exec` で長時間ビルドを強制終了する。API 経由の `build_timeout_seconds` 動的変更は管理 API 実装対象として残す（→ §22 `GET /api/config`）。 | Go test と runner 回帰検証で pipeline 起動経路を検証済み。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ポーリング間隔の動的変更 | systemd タイマーの `OnUnitActiveSec` を変更して間隔を調整（→ §22 `POST /api/schedule/interval`） | `docs/DETAIL_INDEX.md` §0i.3、`docs/details/api.md` §22.0e、`docs/details/setup.md` §26、`docs/details/api.md` §27.11 に従って実装する。 |
| 実装済み | 完了済み | CI ランナー | ビルドログのファイル保存 | `os/exec` で起動したビルドプロセスの stdout/stderr を `.build_logs/{id}.json` に記録（→ §11 ファイル構成）。 | Go test で build log 生成と pipeline stdout/stderr 記録を検証済み。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | GitHub Webhook 受信 | 定期ポーリングと併用可能な即時検出方式。`POST /api/webhook` で GitHub push イベントを受信し即時ビルドをトリガーする（HMAC-SHA256 署名検証付き → §22） | `docs/DETAIL_INDEX.md` §0i.3、`docs/details/runner.md` §13、`docs/details/api.md` §22.0e、`docs/details/api.md` §22-W、`docs/details/api.md` §27.12 に従って実装する。 |
| 実装済み | 完了済み | CI ランナー | ネットワーク断時の再試行 | GitHub API 失敗時に指数バックオフ（`API_RETRY_BASE_SECONDS × 2^n`、最大 `API_RETRY_MAX` 回）で再試行する（→ §12・`docs/ROADMAP.md` §5.2）。 | Go test で fake GitHub 一時失敗からの retry 成功を検証済み。 |
| 実装済み | 完了済み | CI ランナー | GitHub API レート制限自動待機 | `X-RateLimit-Remaining: 0` 検出時に `X-RateLimit-Reset` まで待機してから再試行する（→ `docs/ROADMAP.md` §5.2）。 | Go test 対象の retry 経路と同じ GitHub API retry 実装で検証済み。 |
| 実装済み | 完了済み | CI ランナー | 転送後リモート整合性検証 | SSH 転送後に `sha256sum` でリモートファイルを検証し、不一致時はペンディングキューへ再投入する（→ `docs/ROADMAP.md` §5.2・§14a）。 | Go test で転送失敗時 pending 化と pending 再試行成功を検証済み。 |
| 実装済み | 完了済み | CI ランナー | マルチブランチビルド | `BRANCH_TARGETS` リストで複数ブランチを順次ポーリング・ビルド・転送する（→ §12 設定値・`docs/ROADMAP.md` §5.2 処理フロー）。 | Go test で branch target 設定に基づく処理経路を検証済み。 |
| 実装済み | 完了済み | CI ランナー | ビルドログ世代管理 | `LOG_KEEP_N` 件を超えた `.build_logs/{id}.json` を古いものから自動削除する（→ §12・`docs/ROADMAP.md` §5.2）。 | `go test ./...` で runner 回帰検証済み。 |
| 実装済み | 完了済み | CI ランナー | ビルド出力の外部転送 | ビルド成功時に生成静的 Web サイトを SSH 経由（差分転送・複数ファイル対応）で静的コンテンツ配信サーバーへ自動転送する（→ §14a）。 | Go test で SSH command 差し替えによる pending / retry 経路を検証済み。 |
| 実装済み | 完了済み | CI ランナー | ビルドクールダウン | 前回ビルド完了から `BUILD_COOLDOWN_SECONDS` 秒以内の起動はビルドをスキップする（Webhook 二重トリガー防止 → §12・`docs/ROADMAP.md` §5.2）。 | Go test で cooldown 中の build skip を検証済み。 |
| 実装済み | 完了済み | CI ランナー | ビルド前の事前チェック | `pipeline.sh` 実行前にディスク空き容量・`adlaire-ci-build` 実行可否・`components/builder.go` 由来のビルドバイナリ配置を確認し、不足時は `failure_precheck` として記録する（→ `docs/ROADMAP.md` §5.2）。 | Go test で fake binary 不在時の `failure_precheck` を検証済み。 |
| 実装済み | 完了済み | CI ランナー | 定期強制ビルド | `FORCE_BUILD_INTERVAL`（時間単位）設定時、変更なしでも前回ビルドから経過時間超過で強制ビルドする。API 経由の動的変更は管理 API 実装対象として残す（→ §12・`docs/ROADMAP.md` §5.2）。 | Go test で同一 SHA かつ force interval 超過時の build 実行を検証済み。 |
| 実装済み | 完了済み | CI ランナー | ビルド中重複スキップ | `.build_lock` に PID を記録し、起動時に実行中ビルドを検出したらスキップする（→ §11・`docs/ROADMAP.md` §5.2）。 | Go test で lock conflict 時の build skip を検証済み。 |
| 実装済み | 完了済み | CI ランナー | GitHub PAT 有効期限の事前警告 | GitHub API レスポンスの `GitHub-Authentication-Token-Expiration` ヘッダーを解析し、7 日以内の期限切れを WARN ログで通知する（→ `docs/ROADMAP.md` §5.2）。 | Go test で GitHub API response header 処理を検証済み。 |
| 実装済み | 完了済み | CI ランナー | コミット情報のビルドログ記録 | ビルドトリガーとなったコミットの SHA・メッセージ・作者名・コミット日時を `.build_logs/{id}.json` に記録する（→ `docs/ROADMAP.md` §5.2）。 | Go test で fake commits API からの commit info 記録を検証済み。 |
| 実装済み | 完了済み | CI ランナー | GitHub API 連続失敗によるサーキットブレーカー | 連続失敗が `API_CIRCUIT_BREAKER_THRESHOLD` 周回以上になった場合に `.build_circuit_state.open=true` とし、open 中はポーリングをスキップする（→ §11・§12・`docs/ROADMAP.md` §5.2）。 | Go test で circuit open 時の polling skip を検証済み。 |
| 実装済み | 完了済み | CI ランナー | 出力サイトサイズ警告閾値 | ビルド後の出力サイト合計サイズが `OUTPUT_SIZE_WARN_MB` を超えた場合に WARN ログを出力する。§8 変換レポートに `size_warn` フラグを追加（→ §8・§12・`docs/ROADMAP.md` §5.2・§22）。 | Go test で閾値超過時の `size_warn=true` と WARN 記録を検証済み。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | Webhook イベントログ | 受信した Webhook push イベントを `.webhook_events.json` に JSON Lines 形式で追記記録する。`delivery_id`・`event`・`ref`・`sha`・`build_triggered` を保存（→ §11・§22） | `docs/DETAIL_INDEX.md` §0i.3、`docs/details/runner.md` §11、`docs/details/api.md` §22.0e、`docs/details/api.md` §22-W、`docs/details/api.md` §27.13 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド所要時間の記録と統計 API | `.build_logs/{id}.json` に `started_at`・`finished_at`・`duration_seconds` を記録し、`GET /api/stats/build-duration` で過去 N 件の平均・最小・最大を提供する（→ §22） | `docs/DETAIL_INDEX.md` §0i.4、`docs/details/runner.md` §15、`docs/details/api.md` §22.0e、`docs/details/runner.md` §27.14 に従って実装する。 |
| 実装済み | 完了済み | CI ランナー | ビルドアーティファクト世代管理 | `HISTORY_KEEP_N` 世代分を `.snapshots/` に自動保持し超過分を削除する。`POST /api/history/{id}/rollback` による再転送は API 実装対象として残す（→ §14b）。 | Go test で build 成功時の `.snapshots/{id}/site` 作成を検証済み。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | ビルドアーティファクト管理 | スナップショット一覧・ダウンロード・削除・ロールバック（→ §14b・§22 `POST /api/history/{id}/rollback`） | `docs/DETAIL_INDEX.md` §0i.4、`docs/details/runner.md` §14b、`docs/details/api.md` §22.0e、`docs/details/sdk.md` §23、`docs/details/ui.md` §24、`docs/details/archive.md` §27.15 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | ヘルスチェックエンドポイント | `GET /api/health` を拡充。最終ビルド時刻・最終ビルド結果・最終転送結果・稼働秒数を返す（→ §22） | `docs/DETAIL_INDEX.md` §0i.3、`docs/details/api.md` §22.0e、`docs/details/api.md` §27.16 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | Webhook イベント一覧取得 API | `.webhook_events.json` をページネーション付きで返す `GET /api/webhook-events` を追加する（→ §22） | `docs/DETAIL_INDEX.md` §0i.3、`docs/details/api.md` §22.0e、`docs/details/sdk.md` §23、`docs/details/ui.md` §24、`docs/details/api.md` §27.13 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | ビルドログ重大度フィルター | 既存の `GET /api/logs/search` に `level=warn\|error` パラメータを追加し、重大度別に絞り込む（→ §22） | `docs/DETAIL_INDEX.md` §0i.3、`docs/details/api.md` §22.0e、`docs/details/sdk.md` §23、`docs/details/ui.md` §24、`docs/details/api.md` §27.17 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 変換レポート出力 | ビルド完了後に変換統計（見出し数・テーブル数・コードブロック数・警告）を stdout 出力する。`components/runner.go` が取り込み `GET /api/output-meta` で参照可（→ §8・§22） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | シンタックスハイライト | コードブロックに言語別色分けを `assets/app.js` で適用する。対応言語：`python`・`bash`・`json`・`sql`・`ini`・`diff`（→ §7.8） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 本文内全文検索 | ビルド時に `assets/search-index.json` を生成し、`assets/app.js` の検索 UI と統合して本文ヒット箇所へジャンプ（→ §7.9） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | アンカーリンク自動検証 | 生成 HTML 内の `#anchor` リンクが実際の見出しスラグと一致するか検証し、不整合を `[WARN] BROKEN_LINK` として警告出力する（→ §4.3・§8） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | コードブロックの折りたたみ | 30 行超のコードブロックを初期折りたたみ。「全 N 行を表示」リンクで展開（→ §7.10） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 印刷スタイル（`@media print`） | サイドバー・ヘッダー・ボタン類を非表示、コードブロック展開、リンク URL 末尾表示（→ §6） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 静的 Web サイト出力 | Markdown ファイルまたは Markdown ディレクトリから `index.html`、ページ HTML、`assets/style.css`、`assets/app.js`、`assets/search-index.json` を生成する（→ §5） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | テーマコンポーネント | 初期テーマ `adlaire-default` の header / sidebar / breadcrumb / toc / search / footer / codeblock / table / pagination を内製テンプレートとして提供する（→ §5） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 外部リンクの自動処理 | 外部リンク（`http://`・`https://`）に `target="_blank" rel="noopener noreferrer"` を付与し、内部リンクと区別する（→ §4.3） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 読み取り進捗バー | スクロール位置に応じた 3px プログレスバーをページ上端に固定表示する（→ §7.13） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | コードブロックのコピーボタン | コードブロック右上にワンクリックコピーボタンを配置する（→ §7.6） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 見出しアンカーリンクコピー | ホバーで表示される `.hn-link` ボタンクリックでアンカー URL をクリップボードにコピー（→ §3・§7.11） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | TOC 開閉状態の永続化 | TOC グループの展開／折りたたみ状態を `localStorage` に保存し、リロード後も復元する（→ §7.3） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 見出しスラグ重複解決 | 同一テキストの見出しが複数存在する場合に 2 番目以降のスラグへ `-2`・`-3` を付与して一意にする。TOC・アンカーコピー・全文検索と整合させる（→ §4.5） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 前後章ナビゲーションボタン | h2 見出し単位で「← 前の章」「次の章 →」ボタンを各章末尾に静的生成する（→ §4.5・§5・§7.15） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 内部リンク整合性チェック | `[label](#anchor)` 形式の内部リンクが実際のスラグと一致するか変換時に検証し、不一致を `[WARN]` で報告。§8 変換レポートの `broken_links` フィールドに件数を記録する（→ §4.3・§8） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 見出し階層スキップ警告 | h1→h3 のような見出しレベルの 2 段以上のスキップを `[WARN]` で報告。§8 変換レポートの `heading_skips` フィールドに件数を記録する（→ §4.5・§8） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 読了時間推計と表示 | 本文文字数（コードブロック・タグ除く）から読了時間（分、200文字/分・切り上げ）を算出し、固定ヘッダーに静的埋め込みする。§8 変換レポートの `reading_time` フィールドに記録する（→ §4.5・§5・§6・§8） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | CI ランナー | Webhook 通知失敗リトライキュー | `.notify_pending`（JSON）を起動時に再送し、HTTP 2xx 成功時に削除、失敗時に `retry_count` と `last_error` を更新して保持する（→ §11・`docs/ROADMAP.md` §5.2）。 | Go test で fake HTTP endpoint への再送成功と queue 空化を検証済み。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ブランチ設定の動的変更 API | `BRANCH_TARGETS` を外部 JSON（`.branch_config`）で管理し `GET /api/branch-config` / `POST /api/branch-config` で API 経由変更可能にする。`components/runner.go` 再起動不要（→ §11・§12・§22） | `docs/DETAIL_INDEX.md` §0i.3、`docs/details/runner.md` §11、`docs/details/runner.md` §12、`docs/details/api.md` §22.0e、`docs/details/api.md` §27.18 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | 週次ビルドサマリー Webhook | 指定曜日・時刻に過去 7 日間の成功率・平均ビルド時間・エラー件数をまとめた定期通知を送信する（→ §12・`docs/ROADMAP.md` §5.2・§22） | `docs/DETAIL_INDEX.md` §0i.3、`docs/details/runner.md` §12、`docs/details/runner.md` §13、`docs/details/runner.md` §16、`docs/details/api.md` §22.0e、`docs/details/runner.md` §27.19 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | 設定変更の詳細 diff 記録 | `.config_log` の各エントリに変更前後の値の diff 文字列を付加し `GET /api/config-log` レスポンスに含める（→ §22） | `docs/DETAIL_INDEX.md` §0i.3、`docs/details/statefile.md` §22.0a、`docs/details/api.md` §22.0e、`docs/details/api.md` §27.20 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | テーブルのソート機能 | 列ヘッダークリックで昇順/降順ソートができるインタラクティブテーブル。`aria-sort` 属性と CSS `::after` でインジケーター表示（→ §7.14） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | キーボードショートカット | `/` で検索フォーカス・`Escape` で検索クリア・`t` でページ先頭へスクロール（→ §7.12） | `docs/DETAIL_INDEX.md` §0h・§0i.1 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | 複数ファイル監視 | `target_files` で複数 Markdown ファイルまたは Markdown ディレクトリを監視し、変更対象ごとに build target を決定する。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §11、`docs/details/runner.md` §12、`docs/details/runner.md` §13、`docs/details/runner.md` §15、`docs/details/runner.md` §27.21 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | GitHub Commit Status API | ビルド開始時・成功時・失敗時・pending 時に GitHub Commit Status API へ `pending` / `success` / `failure` を送信し、対象 commit に Adlaire CI の結果を紐付ける。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §12、`docs/details/runner.md` §13、`docs/details/runner.md` §15、`docs/details/statefile.md` §22.0c、`docs/details/api.md` §22.0e、`docs/details/commitstatus.md` §27.1 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルドパイプライン YAML 定義 | 内製 YAML subset parser で `.pipeline.yml` を読み込み、build / test / deploy step を固定順で実行する。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §12、`docs/details/runner.md` §13、`docs/details/runner.md` §15、`docs/details/api.md` §22.0e、`docs/details/runner.md` §27.22 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ドライラン実行モード | `adlaire-ci-runner --dry-run` で設定・状態・GitHub target・SHA 差分・起動可否を検証し、ビルド、deploy、通知、状態更新を行わず結果を出力する。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §11、`docs/details/runner.md` §12、`docs/details/runner.md` §13、`docs/details/runner.md` §15、`docs/details/runner.md` §27.2 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルドログのアーカイブ圧縮 | 保持期間を超えた `.build_logs/{id}.json` を `.build_logs/archive/{id}.json.gz` に gzip 圧縮し、通常ログ API は圧縮済みログも透過的に参照する。 | `docs/DETAIL_INDEX.md` §0i.4、`docs/details/runner.md` §12、`docs/details/runner.md` §13、`docs/details/runner.md` §15、`docs/details/statefile.md` §22.0c、`docs/details/api.md` §22.0e、`docs/details/archive.md` §27.7 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ローカルファイル監視モード | GitHub API を使わず、ローカル状態 snapshot の SHA-256 差分で対象 Markdown の変更を検出する。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §11、`docs/details/runner.md` §12、`docs/details/runner.md` §13、`docs/details/runner.md` §27.23 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | タグ付きコミットのみビルド | 監視対象 commit に許可 tag pattern が付いている場合だけ build を実行する。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §12、`docs/details/runner.md` §13、`docs/details/runner.md` §15、`docs/details/api.md` §22.0e、`docs/details/runner.md` §27.24 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルドキャッシュ | 入力ファイル単位の SHA-256 manifest を `.build_cache.json` に保存し、未変更ページの変換結果を再利用する。 | `docs/DETAIL_INDEX.md` §0i.1、`docs/details/builder.md` §5、`docs/details/builder.md` §8、`docs/details/runner.md` §11、`docs/details/runner.md` §13、`docs/details/builder.md` §27.25 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド通知連携 | `components/runner.go` が `.notify_config.channels` に基づき Webhook / email / command 通知を同一通知イベント契約で送信する。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §13、`docs/details/runner.md` §15、`docs/details/runner.md` §16、`docs/details/api.md` §22.0e、`docs/details/runner.md` §27.32 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド時間トレンド記録 | build ごとの所要時間を `.build_trends.json` に集計保存し、移動平均、中央値、p95、直近件数を API / UI から参照できるようにする。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §13、`docs/details/runner.md` §15、`docs/details/api.md` §22.0e、`docs/details/runner.md` §27.33 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド失敗時の自動リトライ | 一時的な GitHub API / network / pipeline timeout / deploy 検証失敗を対象に、設定回数まで同一 build id 内で再試行し、attempt ごとの結果を build log と history に記録する。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §12、`docs/details/runner.md` §13、`docs/details/runner.md` §15、`docs/details/statefile.md` §22.0c、`docs/details/runner.md` §27.3 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | 並列マルチターゲットビルド | branch target 内の複数 deploy target を bounded worker で並列処理し、target ごとの結果を build log に記録する。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §12、`docs/details/runner.md` §13、`docs/details/runner.md` §14a、`docs/details/runner.md` §15、`docs/details/runner.md` §27.26 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド前後フック | `.hooks` の pre / post hook を shell 経由なしで実行し、abort 条件と hook log を固定仕様どおり扱う。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §13、`docs/details/runner.md` §15、`docs/details/api.md` §22.0e、`docs/details/runner.md` §27.27 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | 依存ファイルトラッキング | Markdown 内の相対リンク・画像・include 対象を依存 manifest として記録し、関連 target だけを再ビルドする。 | `docs/DETAIL_INDEX.md` §0i.1、`docs/details/builder.md` §4.3、`docs/details/builder.md` §5、`docs/details/runner.md` §11、`docs/details/runner.md` §13、`docs/details/builder.md` §27.28 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | リモートビルド対応 | SSH 経由で remote build command を実行し、成果物 archive と manifest を取得して既存 deploy / log 契約へ接続する。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §12、`docs/details/runner.md` §13、`docs/details/runner.md` §14a、`docs/details/runner.md` §15、`docs/details/runner.md` §27.29 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルドステータスファイル出力 | runner の現在状態、最終ビルド、最終 deploy、pending 件数、circuit 状態、最終 trigger を `.build_status.json` に JSON object として出力し、API / UI / MCP の read-only 参照元にする。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §11、`docs/details/runner.md` §13、`docs/details/runner.md` §15、`docs/details/statefile.md` §22.0a、`docs/details/statefile.md` §22.0c、`docs/details/api.md` §22.0e、`docs/details/runner.md` §27.8 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド承認フロー | approval_required な target を `.approval_queue` に保留し、API 承認後だけ build / deploy を継続する。 | `docs/DETAIL_INDEX.md` §0i.3、`docs/details/runner.md` §11、`docs/details/runner.md` §13、`docs/details/runner.md` §15、`docs/details/runner.md` §16、`docs/details/api.md` §22.0e、`docs/details/runner.md` §27.30 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ブランチ別環境変数 | branch target ごとに許可済み環境変数を build process へ注入し、secret を log に出さない。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §12、`docs/details/runner.md` §13、`docs/details/runner.md` §15、`docs/details/api.md` §22.0e、`docs/details/runner.md` §27.31 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド依存チェーン | `.build_chain_config` で job 間依存を定義し、依存成功後だけ後続 job を実行する。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §11、`docs/details/runner.md` §13、`docs/details/runner.md` §15、`docs/details/api.md` §22.0e、`docs/details/runner.md` §27.34 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド優先度キュー | `.build_state.queued` に priority / created_seq を保存し、優先度順かつ同一優先度 FIFO で build queue を処理する。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §11、`docs/details/runner.md` §13、`docs/details/api.md` §22.0e、`docs/details/runner.md` §27.35 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | 失敗原因の自動分類 | stdout / stderr / exit code / runner error を固定分類ルールで解析し、failure_category と evidence を build log / history に記録する。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §13、`docs/details/runner.md` §15、`docs/details/api.md` §22.0e、`docs/details/runner.md` §27.36 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド実行環境の記録 | build 開始時の OS、arch、Go version、binary version、disk usage、hostname を `.build_logs/{id}.json.environment` に記録する。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §13、`docs/details/runner.md` §15、`docs/details/runner.md` §27.37 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルドトリガー種別の記録 | `polling`、`force_interval`、`manual`、`webhook`、`retry_pending_transfer`、`startup_config_integrity`、`rollback`、`local_watch`、`approval` を `.build_logs/{id}.json`、`.build_history`、`.build_status.json` に記録し、履歴 filter と状態表示で同じ値を使う。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §13、`docs/details/runner.md` §15、`docs/details/statefile.md` §22.0c、`docs/details/api.md` §22.0e、`docs/details/sdk.md` §23、`docs/details/ui.md` §24、`docs/details/runner.md` §27.9、`docs/details/runner.md` §27.23、`docs/details/runner.md` §27.30 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド所要時間の異常検知 | `.build_trends.json` の移動平均と p95 を基準に異常に遅い build を検出し、WARN、history flag、通知へ反映する。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §13、`docs/details/runner.md` §15、`docs/details/runner.md` §16、`docs/details/api.md` §22.0e、`docs/details/runner.md` §27.38 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | 設定ファイル起動時整合性チェック | runner 起動時に `.branch_config`、`.notify_config`、`.build_state`、`.pending_transfers`、`.notify_pending`、`.build_circuit_state` の JSON 整合性を検証し、破損・型不一致・必須 key 不足を規定どおり退避、初期化、通知、または停止する。 | `docs/DETAIL_INDEX.md` §0i.2、`docs/details/runner.md` §11、`docs/details/runner.md` §12、`docs/details/runner.md` §13、`docs/details/statefile.md` §22.0a、`docs/details/statefile.md` §22.0c、`docs/details/runner.md` §27.10 に従って実装する。 |
| 将来計画 | 実装不可 | 管理ツール・API | マルチユーザー対応 | 初期仕様の単一 admin 認証を、複数ユーザー・ユーザー別セッション・ユーザー別監査へ拡張する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 通知先の拡張 | 管理画面・API からメール・Slack・Discord 等の通知チャンネルを設定・追加できるようにする。`components/runner.go` 側のフック実装は → ビルド通知連携 | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | データストア切り替え | 大量ビルド履歴・ログ運用に備え、フラットファイルから SQLite 等への切り替えを検討する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 外部認証連携 | SSO・OAuth 等の外部認証基盤との連携を検討する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | TOTP 二要素認証 | RFC 6238 TOTP を Go 標準ライブラリだけで検証し、login を password verified / totp required / session issued の二段階に分離する。 | `docs/DETAIL_INDEX.md` §0i.4、`docs/details/statefile.md` §22.0a、`docs/details/api.md` §22.0e、`docs/details/sdk.md` §23、`docs/details/ui.md` §24、`docs/details/security.md` §27.46 に従って実装する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 統計データの CSV エクスポート | `GET /api/stats/timeline` の日別データを CSV 形式でダウンロードできるエンドポイントを追加する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | キュー内個別エントリのキャンセル | `DELETE /api/queue/{id}` で特定エントリのみキャンセルする（初期仕様では全クリアのみ） | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | 設定バリデーション API | `POST /api/config/validate` で `.server_config` 互換の設定差分を保存前に検証し、正規化後設定、警告、エラー位置を返す。状態ファイルは変更しない。 | `docs/DETAIL_INDEX.md` §0i.3、`docs/details/statefile.md` §22.0c、`docs/details/api.md` §22.0e、`docs/details/sdk.md` §23、`docs/details/ui.md` §24、`docs/details/api.md` §27.5 に従って実装する。 |
| 将来計画 | 実装不可 | 管理ツール・API | Prometheus メトリクスエンドポイント | `GET /api/metrics` で Prometheus 形式のメトリクス（ビルド数・成功率・ディスク使用量等）を返す | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | CLI 管理クライアント | Go 標準ライブラリのみで実装した `adlaire-ci-cli` で API を CUI 操作できるツール | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 設定の自動スナップショット | `POST /api/config` 変更時に自動で設定バックアップを世代保存する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ステータスバッジ生成 | `GET /api/badge` で最終ビルド結果を SVG バッジとして返す（README 埋め込み用） | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ビルド履歴の自動削除設定 | `history_retention_days` 設定でビルド履歴エントリを自動削除する（ログの `log_retention_days` に対応する履歴版） | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | セッションタイムアウト変更設定 | `.server_config.session_timeout_seconds` で session 有効期限を 5 分〜30 日の範囲で変更し、既存 session の扱いを固定する。 | `docs/DETAIL_INDEX.md` §0i.4、`docs/details/statefile.md` §22.0c、`docs/details/api.md` §22.0e、`docs/details/sdk.md` §23、`docs/details/ui.md` §24、`docs/details/security.md` §27.45 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | ビルドトリガー専用 API スコープ | API token scope に `trigger` を追加し、build 起動系だけを許可する最小権限 token を発行できるようにする。 | `docs/DETAIL_INDEX.md` §0i.4、`docs/details/statefile.md` §22.0a、`docs/details/api.md` §22.0e、`docs/details/sdk.md` §23、`docs/details/ui.md` §24、`docs/details/security.md` §27.42 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | 監査ログ | `.audit_log` に設定変更、認証、token、build trigger、承認、権限拒否を actor 付き JSON Lines で記録する。 | `docs/DETAIL_INDEX.md` §0i.4、`docs/details/statefile.md` §22.0a、`docs/details/api.md` §22.0e、`docs/details/sdk.md` §23、`docs/details/ui.md` §24、`docs/details/security.md` §27.44 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | API レート制限 | `.api_rate_state` で actor / IP / endpoint group ごとの固定窓 rate limit を管理し、超過時 `429` を返す。 | `docs/DETAIL_INDEX.md` §0i.4、`docs/details/statefile.md` §22.0a、`docs/details/api.md` §22.0e、`docs/details/sdk.md` §23、`docs/details/ui.md` §24、`docs/details/security.md` §27.47 に従って実装する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ロールベースアクセス制御 | 複数ユーザー対応後に、管理者・オペレーター・閲覧者等の役割ごとに API 権限を分ける | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 成果物ダウンロード API | 生成静的 Web サイトを archive として API エンドポイント経由で直接ダウンロードできるようにする | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 設定スナップショット差分表示 | 保存済みスナップショット間の設定変更点を diff 形式で確認できる API | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 複数プロジェクト管理 | 単一インスタンスで複数リポジトリ／プロジェクトを切り替え管理する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | API キー管理 | `.api_tokens` で API key の hash、scopes、expires_at、revoked_at を管理し、発行時だけ token 本体を返す。 | `docs/DETAIL_INDEX.md` §0i.4、`docs/details/statefile.md` §22.0a、`docs/details/api.md` §22.0e、`docs/details/sdk.md` §23、`docs/details/ui.md` §24、`docs/details/security.md` §27.43 に従って実装する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ビルドログのリアルタイム配信 | 実行中ビルドのログを SSE / WebSocket でストリーミング配信するエンドポイント | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 設定のインポート／エクスポート | 設定全体を JSON でエクスポートし、別環境へそのままインポートする | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ビルド統計ダッシュボード | 成功率・平均ビルド時間・エラー分布等を可視化する管理画面を生成する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ユーザー管理 API | 複数ユーザー対応後に、管理者アカウントの追加・削除・パスワード変更を API で操作する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | IP アドレス制限 | 管理 API へのアクセスを許可 IP レンジに限定するフィルタリング | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | API バージョニング | `/api/v1/` 等のバージョンプレフィックスで API 世代を明示する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | Webhook 署名検証 | 受信 Webhook の HMAC 署名を検証し、なりすましリクエストを拒否する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | API ドキュメント自動生成 | OpenAPI / Swagger 仕様を自動生成し、インタラクティブなドキュメントとして提供する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 通知チャンネル管理 | Slack / Discord / メール等の通知先を管理画面から追加・削除・テスト送信する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ビルドキューの手動並び替え | 管理画面からキュー内ジョブの実行順序を変更する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 設定テンプレート | よく使う設定パターンをテンプレートとして保存・再利用する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | API アクセスログ | 認証後 API と Webhook の呼び出しを `.api_access_log` に JSON Lines で記録し、`GET /api/api-access-log` でページング参照する。認証ログ `.access_log` とは分離する。 | `docs/DETAIL_INDEX.md` §0i.3、`docs/details/statefile.md` §22.0a、`docs/details/statefile.md` §22.0c、`docs/details/api.md` §22.0e、`docs/details/sdk.md` §23、`docs/details/ui.md` §24、`docs/details/api.md` §27.6 に従って実装する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 管理者向けイベントフィード | ビルド完了・エラー・設定変更等のシステムイベントをリアルタイムで流す管理ページ | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ビルドキュー可視化 | キューに積まれたビルドの状態一覧を静的 HTML ステータスページとして出力する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | メンテナンスモード | 管理 API から即時にメンテナンスモードへ切り替え、ビルドキューを一時停止する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 読み取り専用共有リンク | ビルドステータス・統計を外部に公開する期限付き読み取り専用リンクを発行する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | アラート閾値設定 | ビルド失敗率・所要時間等が設定閾値を超えた際に自動アラートを発火する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | バックアップ／リストア | 設定・ビルド履歴・ログ等の全データをアーカイブ化してリストアできる機能 | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | API レスポンスキャッシュ制御 | 頻繁に参照される統計・ログ API のキャッシュ TTL を設定から変更する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | スナップショット間サイト差分 API | 2 つのスナップショット ID を指定し、出力サイトの追加/削除行数・変更率を返す `GET /api/snapshots/{id1}/diff/{id2}` を追加する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | Webhook 送信履歴の手動再送 API | `GET /api/notify-log` の各エントリに対して `POST /api/notify-log/{id}/retry` で同一ペイロードを即時再送できる手動リトライ API。`.notify_pending` 自動再試行とは別に特定通知だけ個別再送できる運用機能 | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 差分ビルド | 変更箇所のみ処理し、大規模 MD の変換を高速化する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 複数出力形式 | HTML に加えて PDF・ePub 等の出力形式をサポートする | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | Markdown 拡張記法サポート | アドモニション（`> [!NOTE]`）・カラーバッジ等の独自拡張記法に対応する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | コードブロック行番号表示 | コードブロック左端に行番号を表示するオプションを追加する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 見出しの自動採番 | h2 以下の見出しに `1.1`・`1.2` 等の番号を自動付与するオプション | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | セクション折りたたみ | 見出しクリックでコンテンツを折りたたむ機能（デフォルト展開） | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | TOC 深さ制御 | TOC に含める見出しレベルを設定で指定する（例：h2–h3 のみ） | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 最終更新日の自動埋め込み | ソースの git コミットタイムスタンプをフッターに自動出力する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | diff ハイライト | `+`/`-` で始まる行を git diff スタイルで緑/赤に色分けするコードブロックオプション | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 画像の遅延読み込み | `<img>` に `loading="lazy"` を付与し、初期表示を高速化する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | カスタムメタタグ注入 | OGP / Twitter Card 等のメタタグをビルド設定から生成・埋め込む | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | ダークモード対応 | `prefers-color-scheme` に応じたライト／ダーク切り替えを実装する（現在はライト固定） | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | コードブロックのファイル名表示 | ` ```go:filename.go ` や ` ```言語名:filename.ext ` 記法でコードブロック上部にファイル名ラベルを表示する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | テンプレート変数展開 | ビルド設定に定義した変数を `{{ VERSION }}` 等の記法で Markdown 本文中に展開する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | HTML ミニファイ | 生成 HTML のホワイトスペース・コメントを除去してファイルサイズを削減する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | TOC ハイライト追従 | スクロール位置に応じてサイドバー TOC の現在セクションを自動ハイライトする | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | Mermaid ダイアグラム描画 | ` ```mermaid ` コードブロックをフローチャート・シーケンス図として SVG 描画する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 脚注サポート | `[^1]` 記法の脚注をページ末尾に自動レンダリングする | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | インライン数式レンダリング | `$...$` / `$$...$$` 記法の数式を KaTeX 等で描画する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | ページ内ナビゲーション履歴 | ブラウザの戻る/進むに対応したハッシュベースの履歴管理を実装する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 読み上げ対応（アクセシビリティ） | `aria-label`・`role` 属性の付与対象、値、検証方法を詳細仕様で定義した上でスクリーンリーダー閲覧に対応する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 画像ライトボックス | 画像クリックでモーダル拡大表示する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | 出力サイトへのビルドメタ埋め込み | `adlaire-ci-build` が build id、commit SHA、build at を受け取り、生成 HTML の `<head>` に固定 meta として埋め込む。`GET /api/output-meta` は同値を返す。 | `docs/DETAIL_INDEX.md` §0i.1、`docs/details/builder.md` §2、`docs/details/builder.md` §5、`docs/details/builder.md` §8、`docs/details/runner.md` §13、`docs/details/api.md` §22.0e、`docs/details/builder.md` §27.4 に従って実装する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 印刷時 QR コード挿入 | `@media print` で元ページの URL を QR コードとしてフッターに埋め込む | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 定義リストサポート | `term\n: definition` 記法を `<dl>/<dt>/<dd>` タグにレンダリングする | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | タスクリストサポート | `- [ ]` / `- [x]` 記法をチェックボックス付きリストとして描画する | 5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP サーバー実装 | `components/mcp.go` を将来追加コンポーネントとして追加。MCP プロトコル（JSON-RPC over stdio）で Claude Desktop 等の AI クライアントから直接接続可能にする。内部では `components/api.go` REST API に Go 標準ライブラリ `net/http` でローカル接続するラッパー設計（`encoding/json` + `os.Stdin` / `os.Stdout` + `net/http`、ゼロ外部依存）。認証は `.mcp_token` に専用 API トークンを保存し、スコープ（`read` のみ / `trigger` 許可）をトークン単位で選択可能。Claude Desktop の `mcpServers` 設定に `/usr/local/bin/adlaire-ci-mcp` を指定して接続する | 5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP ツール・リソース公開 | MCP サーバーが公開するツール：`get_status`（ビルド状態・CB 状態・PAT 残日数）/ `get_history(n)`（直近 N 件）/ `search_logs(query, level?, from?, to?)`（ログ全文検索）/ `get_build_log(id)`（個別ビルドログ）/ `trigger_build(force?)`（ビルドトリガー、`trigger` スコープ必須）/ `reset_circuit_breaker`（CB リセット、`trigger` スコープ必須）。リソース：`adlaire://status` / `adlaire://history` / `adlaire://logs/{id}` / `adlaire://config`（→ MCP サーバー実装） | 5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | AI 支援ビルドエラー分析 | ビルド失敗時、AI クライアント（Claude Desktop 等）が `get_build_log` / `search_logs` ツールを自律的に呼び出してエラーログを取得し、原因推定と修正提案を生成できる設計。AI 側が pull するため `components/runner.go` のゼロ依存を完全維持。将来的には Webhook 通知をトリガーに AI が自動分析を開始する構成も検討可（→ MCP サーバー実装） | 5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP Prompts 定義 | よく使う分析シナリオを MCP Prompts として定義し、Claude Desktop 等のプロンプトメニューから即時呼び出し可能にする。例：「先週のビルド失敗率をまとめて」「最後のエラーの原因を分析して」「PAT 有効期限が近いか確認して」。`get_history` / `search_logs` ツールと組み合わせて定型分析を 1 クリックで実行できる（→ MCP サーバー実装） | 5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP Sampling によるビルドログ自動分析 | MCP Sampling 機能を使い、ビルド失敗時に `components/mcp.go` が AI クライアントへ sampling リクエストを送って原因推定テキストを生成し `.build_logs/{id}.json` の `ai_analysis` フィールドへ自動記録する。`components/runner.go` は MCP サーバーへソケット通知を送るだけで Anthropic API キーは `components/mcp.go` も保持しない（→ MCP サーバー実装） | 5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP Notifications（イベントプッシュ） | ビルド完了・失敗・CB 開放等のシステムイベントを MCP Notifications として接続中の AI クライアントへリアルタイムプッシュする。`components/runner.go` が `components/api.go` 経由でイベントをキューに積み `components/mcp.go` が接続クライアントへ転送する設計（→ MCP サーバー実装） | 5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP HTTP SSE transport 対応 | 初期仕様の stdio transport に加えて HTTP + SSE transport をサポートし、リモートマシンや複数クライアントからの同時接続を可能にする。Go 標準ライブラリ `net/http` で実装しゼロ依存を維持。`components/api.go` と同一プロセス統合か独立ポート起動かを設定で選択可能（→ MCP サーバー実装） | 5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP ツールスコープ細分化 | 初期仕様の `read` / `trigger` 2 スコープを `read`（参照のみ）/ `trigger`（ビルド起動）/ `admin`（設定変更・CB リセット・ブランチ設定変更）の 3 スコープに細分化する。`.mcp_token` の各トークンにスコープを紐付け、ツール呼び出し時にスコープ検証を行う（→ MCP サーバー実装） | 5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP ツール呼び出し監査ログ | MCP 経由で呼び出されたツールの履歴（呼び出し日時・ツール名・引数サマリー・成否）を `.mcp_access_log` に記録する。`GET /api/mcp-access-log` で参照可能にし、AI クライアントがどの操作をいつ実行したかを追跡できる（→ MCP サーバー実装） | 5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP リソース購読（Resource Subscriptions） | クライアントが `adlaire://status` 等のリソースを subscribe し、状態変化時に `notifications/resources/updated` を自動受信できる MCP 標準機能。MCP Notifications（イベント起点プッシュ）とは異なりリソース変更起点のプッシュで、クライアントがポーリングなしに最新状態を保持できる（→ MCP サーバー実装） | 5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP クライアント情報ログ | initialize リクエストの `clientInfo`（クライアント名・バージョン）を `.mcp_access_log` の接続エントリとして記録する。Claude Desktop / Cursor / 自作クライアント等どのツールから接続されたかを追跡し、監査と動作確認に利用する（→ MCP サーバー実装） | 5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP ツール実行統計 | ツールごとの呼び出し回数・平均応答時間（ms）・エラー率を `.mcp_stats` に蓄積する。`GET /api/mcp-stats` で参照可能にし、どのツールが頻用されているか・ボトルネックがどこかを可視化する（→ MCP サーバー実装） | 5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP ツール実行タイムアウト設定 | `.mcp_config` にツールごとのタイムアウト秒数を設定可能にする（例：`search_logs: 10`、`trigger_build: 30`）。タイムアウト超過時は JSON-RPC エラーを返し、MCP サーバーが無応答になる事態を防ぐ（→ MCP サーバー実装） | 5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP 設定 CRUD ツール | `get_config(section?)` / `set_config(key, value)` ツールを `admin` スコープで公開する。AI クライアントから直接 `.server_config` / `.notify_config` 等を読み書きでき、チャット上で「ポーリング間隔を 30 秒に変更して」と指示するだけで設定変更が完結する（`components/runner.go` 再起動不要）。`admin` スコープの定義は → MCP ツールスコープ細分化（→ MCP サーバー実装） | 5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP Elicitation による副作用操作の確認 | MCP Elicitation に対応し、`trigger_build(force=true)` / `reset_circuit_breaker` 等の副作用操作の実行前に AI クライアントへ確認プロンプトを送信して明示的な承認を得てから実行する。JSON-RPC メッセージの追加のみで実装しゼロ依存を維持。意図しない操作の誤実行を防ぐ安全機構（→ MCP サーバー実装） | 5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |

### 5.2.3 状態変更・昇格手順

統合ロードマップ表の `将来計画` の項目を実装対象にする場合は、以下の順で進める。

1. 対象項目の `状態` を `改訂予定` に変更し、元状態、格上げ日、担当領域、整理順序（実装単位ではない）、ステータスを同じ行の概要または次アクションへ記録する。
2. `docs/SPEC.md` の該当する方針、ポリシーと、本ファイルの実装状態、機能一覧を改訂する。
3. `docs/DETAIL_INDEX.md` §0i.1〜§0i.4 の詳細節対応表に、対象機能、対象コンポーネント、詳細仕様節、受け入れ条件を追加する。
4. 該当 owner component の詳細仕様本文に、§0h の機能仕様テンプレートを満たす目的、責務、入出力、状態、処理順序、異常系、セキュリティ、検証条件を追加する。MCP サーバー領域を昇格する場合は、MCP 専用詳細仕様を新設し、`docs/DETAIL_INDEX.md` の責務 component 対応表へ追加する。
5. `docs/DOCUMENT_INDEX.md` と `AGENTS.md` の更新要否を確認する。
6. 仕様凍結条件を満たした後、対象項目の `状態` を `仕様化済み・未実装`、`実装可否` を `実装可` に変更する。
7. 実装と検証が完了した後、対象項目の `状態` を `実装済み`、`実装可否` を `完了済み` に変更する。

上記手順を完了していない項目は、実装者判断で実装してはならない。

## 6. 追加仕様化機能実装参照

本節は、追加仕様化機能の実装状態、実装可否、owner、主本文、collaborator を確認するための参照である。各機能の主本文は、owner component の分割先詳細仕様ファイルを正とする。本節に定義された機能は、`docs/DETAIL_INDEX.md` §0i と owner 詳細仕様を入口として、owner 詳細仕様の主本文と必要な collaborator 詳細仕様の確認項目を組み合わせて実装可否を判定する。

### 6.1 追加仕様化機能 共通実装契約

§27 の主本文は、owner component の分割先詳細仕様ファイルを正とする。本節では、§27 の実装時に共通して確認する参照順、越境禁止、PR 証跡の入口だけを示す。

| 確認 | 固定内容 |
|------|----------|
| 実装対象判定 | `docs/ROADMAP.md` で実装状態と実装可否を確認し、将来計画、実装不可、未仕様化、MCP 専用機能を実装対象にしない。 |
| owner 確定 | §0b と §0i で owner component を 1 件に確定し、主本文は owner の分割先詳細仕様ファイルで確認する。 |
| collaborator 確認 | collaborator がある場合は、§27.1〜§27.47 の参照インデックスに列挙された component の分割先ファイルを schema、呼び出し境界、表示、security、setup、fixture、検証観点として読む。 |
| 補完禁止 | 個別節または分割先詳細仕様に存在しない endpoint、状態ファイル、設定 key、UI 操作、SDK method、外部依存を実装判断で追加しない。追加が必要な場合は owner component の詳細仕様、関連 collaborator 詳細仕様、fixture catalog、必要な対応表を先に更新する。 |
| 状態更新 | 状態ファイル更新は `docs/details/statefile.md` §22.0a、§22.0c を正とし、lock、atomic write、JSON Lines、破損時処理を独自定義しない。 |
| security | secret mask、token、session、scope、audit、rate limit は `docs/details/security.md` を正とし、平文保存・平文表示を行わない。 |
| fixture / PR 証跡 | fixture manifest、expected/effects、assertion、PR 証跡、受け入れゲートは `docs/details/fixture.md` §27-F を正とする。 |
| api / sdk / ui 同期 | API endpoint、SDK method、UI 操作が同一機能に関わる場合は、endpoint は `docs/details/api.md`、SDK method は `docs/details/sdk.md`、UI 操作は `docs/details/ui.md` をそれぞれ正本とし、名称、引数、response、error、表示、成功後再取得、失敗時固定が食い違わないことを確認する。 |

§27 の機能を実装した PR は、対象節、owner 詳細仕様、collaborator 詳細仕様、fixture、secret mask、失敗時副作用、実装対象外を PR 本文に記録する。記録が不足する場合は、実装完了として扱わない。

§27 の PR 分割、dry-run 固定契約、fixture 完了条件の詳細は、owner 詳細仕様と `docs/details/fixture.md` §27-F を正とする。`docs/DETAIL_INDEX.md` に同じ fixture schema、expected/effects、個別機能本文を重複定義しない。

### 6.2 追加仕様化機能 参照索引

本節は、§27 機能の参照先を一覧化するインデックスである。個別機能の入力、出力、状態、処理順序、異常系、endpoint、SDK method、UI DOM、fixture は下表の「主本文」に記載された owner component 詳細仕様を正とする。`docs/DETAIL_INDEX.md` は、下表に記載された主本文、owner component、collaborator component を置き換えない。

| 節 | 機能 | owner | 主本文 | collaborator | `docs/DETAIL_INDEX.md` 側の扱い |
|----|------|-------|--------|--------------|--------------------|
| §27.1 | GitHub Commit Status API | `commitstatus` | `docs/details/commitstatus.md` §27.1 | `runner`、`statefile` | Commit Status payload と送信順の参照先だけを示す。runner の build 実行、commit SHA 確定、build id 採番、pipeline / deploy / snapshot / history の最終結果確定は `docs/details/runner.md` を正とする。 |
| §27.2 | ドライラン実行モード | `runner` | `docs/details/runner.md` §27.2 | `statefile` | 状態ファイル非更新、ログ、history、deploy、通知の扱いを`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.3 | ビルド失敗時の自動リトライ | `runner` | `docs/details/runner.md` §27.3 | `statefile` | retry 対象、回数、backoff、SHA 更新禁止条件を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.4 | 出力サイトへのビルドメタ埋め込み | `builder` | `docs/details/builder.md` §27.4 | `runner`、`api`、`statefile` | HTML meta、REPORT、API 表示、状態反映の境界だけを確認する。 |
| §27.5 | 設定バリデーション API | `api` | `docs/details/api.md` §27.5 | `sdk`、`ui`、`statefile` | validate の保存禁止、response、SDK/UI 対応を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.6 | API アクセスログ | `api` | `docs/details/api.md` §27.6 | `sdk`、`ui`、`statefile` | `.api_access_log` schema と一覧 API の主本文を`docs/DETAIL_INDEX.md` へ複製しない。 |
| §27.7 | ビルドログのアーカイブ圧縮 | `archive` | `docs/details/archive.md` §27.7 | `runner`、`api`、`statefile` | gzip archive、cleanup、参照順の実体処理を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.8 | ビルドステータスファイル出力 | `runner` | `docs/details/runner.md` §27.8 | `api`、`statefile` | `.build_status.json` schema と更新タイミングを`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.9 | ビルドトリガー種別の記録 | `runner` | `docs/details/runner.md` §27.9 | `api`、`sdk`、`ui`、`statefile` | trigger 有効値、判定条件、UI 表示の同期確認だけを扱う。 |
| §27.10 | 設定ファイル起動時整合性チェック | `runner` | `docs/details/runner.md` §27.10 | `statefile` | JSON 破損、退避、初期化、終了コードを`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.11 | ポーリング間隔の動的変更 | `api` | `docs/details/api.md` §27.11 | `runner`、`statefile` | systemd timer 反映の導入・検証手順は `docs/details/setup.md` §26 を確認する。 |
| §27.12 | GitHub Webhook 受信 | `api` | `docs/details/api.md` §27.12 | `runner`、`statefile` | HMAC、event 記録、queue 投入、エラー応答を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.13 | Webhook イベントログ / 一覧取得 API | `api` | `docs/details/api.md` §27.13 | `sdk`、`ui`、`statefile` | webhook event schema、一覧 API、SDK/UI 対応を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.14 | ビルド所要時間の記録と統計 API | `runner` | `docs/details/runner.md` §27.14 | `api`、`statefile`、`archive` | duration 計測、統計値、archive 連携の境界だけを確認する。 |
| §27.15 | ビルドアーティファクト管理 | `archive` | `docs/details/archive.md` §27.15 | `api`、`sdk`、`ui`、`runner`、`statefile` | snapshot 作成トリガーは `docs/details/runner.md` §14b を正とする。 |
| §27.16 | ヘルスチェックエンドポイント | `api` | `docs/details/api.md` §27.16 | `statefile` | health response とエラー応答を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.17 | ビルドログ重大度フィルター | `api` | `docs/details/api.md` §27.17 | `sdk`、`ui`、`archive`、`statefile` | log query、検索結果、UI filter の同期確認だけを扱う。 |
| §27.18 | ブランチ設定の動的変更 API | `api` | `docs/details/api.md` §27.18 | `runner`、`statefile` | `.branch_config`、GET/POST API、runner 反映条件を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.19 | 週次ビルドサマリー Webhook | `runner` | `docs/details/runner.md` §27.19 | `api`、`statefile` | 週次集計、通知 payload、手動送信 API の境界だけを確認する。 |
| §27.20 | 設定変更の詳細 diff 記録 | `api` | `docs/details/api.md` §27.20 | `statefile` | `.config_log` diff 形式と mask 条件を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.21 | 複数ファイル監視 | `runner` | `docs/details/runner.md` §27.21 | `builder`、`api`、`statefile` | target_files、SHA 差分、build target、API 表示の同期確認だけを扱う。 |
| §27.22 | ビルドパイプライン YAML 定義 | `runner` | `docs/details/runner.md` §27.22 | `api`、`statefile` | pipeline subset、step 実行順、timeout、env を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.23 | ローカルファイル監視モード | `runner` | `docs/details/runner.md` §27.23 | `statefile` | GitHub API 非使用条件と local snapshot 差分検出を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.24 | タグ付きコミットのみビルド | `runner` | `docs/details/runner.md` §27.24 | `api`、`statefile` | tag pattern、skip 条件、history/log 反映を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.25 | ビルドキャッシュ | `builder` | `docs/details/builder.md` §27.25 | `runner`、`statefile` | cache manifest、再利用条件、無効化条件を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.26 | 並列マルチターゲットビルド | `runner` | `docs/details/runner.md` §27.26 | `statefile` | worker 上限、target 別 status、ログ順序を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.27 | ビルド前後フック | `runner` | `docs/details/runner.md` §27.27 | `api`、`statefile` | hook schema、pre/post 実行、abort 条件を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.28 | 依存ファイルトラッキング | `builder` | `docs/details/builder.md` §27.28 | `runner`、`statefile` | dependency manifest、関連 target 判定、full build 条件を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.29 | リモートビルド対応 | `runner` | `docs/details/runner.md` §27.29 | `api`、`archive`、`statefile` | remote command、archive 取得、manifest 検証を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.30 | ビルド承認フロー | `api` | `docs/details/api.md` §27.30 | `runner`、`sdk`、`ui`、`statefile` | approval queue、承認/却下 API、通知、UI 操作の同期確認だけを扱う。 |
| §27.31 | ブランチ別環境変数 | `runner` | `docs/details/runner.md` §27.31 | `api`、`statefile` | branch env schema、許可 key、secret mask、process env 注入を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.32 | ビルド通知連携 | `runner` | `docs/details/runner.md` §27.32 | `api`、`sdk`、`ui`、`statefile` | 通知 event、channel schema、retry、notify log の境界だけを確認する。 |
| §27.33 | ビルド時間トレンド記録 | `runner` | `docs/details/runner.md` §27.33 | `api`、`statefile` | `.build_trends.json`、移動平均、中央値、p95 を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.34 | ビルド依存チェーン | `runner` | `docs/details/runner.md` §27.34 | `api`、`statefile` | DAG 検証、実行順、skip / failure status を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.35 | ビルド優先度キュー | `runner` | `docs/details/runner.md` §27.35 | `api`、`statefile` | queue priority、created_seq、FIFO を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.36 | 失敗原因の自動分類 | `runner` | `docs/details/runner.md` §27.36 | `api`、`statefile` | failure_category、evidence、分類優先順位を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.37 | ビルド実行環境の記録 | `runner` | `docs/details/runner.md` §27.37 | `statefile` | environment snapshot と secret 非含有条件を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.38 | ビルド所要時間の異常検知 | `runner` | `docs/details/runner.md` §27.38 | `api`、`statefile` | trend 基準、異常判定、通知 payload、設定値を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.42 | ビルドトリガー専用 API スコープ | `security` | `docs/details/security.md` §27.42 | `api`、`sdk`、`ui`、`statefile` | scope 判定、拒否条件、API/SDK/UI 表示を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.43 | API キー管理 | `security` | `docs/details/security.md` §27.43 | `api`、`sdk`、`ui`、`statefile` | API key の一回表示、hash 保存、scope、期限、失効、監査を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.44 | 監査ログ | `security` | `docs/details/security.md` §27.44 | `api`、`statefile` | `.audit_log` schema、対象操作、mask、検索 API、UI 表示を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.45 | セッションタイムアウト変更設定 | `security` | `docs/details/security.md` §27.45 | `api`、`sdk`、`ui`、`statefile` | timeout 範囲、保存、既存 session、新規 session 期限を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.46 | TOTP 二要素認証 | `security` | `docs/details/security.md` §27.46 | `api`、`sdk`、`ui`、`statefile` | TOTP、二段階 login、secret 保存、確認、無効化、UI 操作を`docs/DETAIL_INDEX.md` で再定義しない。 |
| §27.47 | API レート制限 | `security` | `docs/details/security.md` §27.47 | `api`、`sdk`、`ui`、`statefile` | 固定窓制限、`429`、状態保存、設定 API、UI 表示を`docs/DETAIL_INDEX.md` で再定義しない。 |

### 6.3 / §27.38a 横断連動・Runner 拡張機能 実装補足契約

本節は、責務 component 別詳細仕様へ分割せず、追加仕様化機能の横断補足契約として扱う。§27.21〜§27.38 および api / sdk / ui / statefile の横断連動は runner、builder、api、sdk、ui、statefile、archive にまたがるため、実装者は owner 詳細仕様を主本文とし、本節を横断確認として同時に確認する。

本節の正本範囲は、横断確認、同期禁止、横断処理順、成功後再取得、失敗時固定、実装完了時の横断受け入れ観点に限定する。本節は、個別機能の処理本文、入力、出力、状態 schema、fixture schema、endpoint 詳細、SDK method 詳細、UI DOM 詳細を持たない。これらは各 owner / collaborator の分割先詳細仕様ファイルを正とする。

本節と owner component 別詳細仕様ファイルの内容が矛盾する場合は、個別機能の入出力、状態、処理、異常系、endpoint、SDK、UI、fixture は owner component 別詳細仕様ファイルを正とし、横断処理順、API / SDK / UI / statefile 同期、成功後再取得、失敗時固定だけを本節で確認する。§27.38a を理由に、owner 詳細仕様に存在しない endpoint、SDK method、UI 操作、状態ファイル、設定 key、fixture を追加してはならない。

| 確認 | 固定内容 |
|------|----------|
| runner 起点 | §27.21〜§27.38 の多くは runner の build 実行、queue、history、log、notification に影響するため、`docs/details/runner.md` の該当 §27 節を先に確認する。 |
| builder 連携 | cache、dependency、output meta、生成物に関わる場合は `docs/details/builder.md` の該当 §27 節を同時に確認する。 |
| API 連携 | 設定保存、queue、approval、history、stats、snapshot、rollback、hook、notify、search に関わる場合は `docs/details/api.md` の endpoint / state read-write 契約を同時に確認する。 |
| SDK / UI 連携 | API を管理画面から操作する機能は、SDK method は `docs/details/sdk.md` §23、DOM / 表示条件は `docs/details/ui.md` §24 を正本として確認する。 |
| statefile | 状態 schema、lock、atomic write、JSON Lines、破損時処理、保存順は `docs/details/statefile.md` を正とする。 |
| archive | snapshot、artifact、download、delete、rollback、log archive は `docs/details/archive.md` を正とする。 |
| fixture | §27.21〜§27.38 の受け入れ fixture、secret mask、effects、PR 証跡は `docs/details/fixture.md` §27-F を正とする。 |

**api / sdk / ui / statefile 横断連動契約：**

下表は横断確認表であり、新しい API endpoint、SDK method、UI 操作、状態ファイル副作用を定義する表ではない。下表の機能群は、各 owner component 別詳細仕様ファイルに定義済みの API endpoint、SDK method、UI 操作、状態ファイル副作用を同じ実装単位でそろえる。API だけ、SDK だけ、UI だけを先行して仕様外の仮実装にしてはならない。UI が未実装の Phase では、UI 列は fixture の期待操作として固定し、実装完了扱いには含めない。

| 機能群 | API | SDK | UI 操作 | 状態ファイル副作用 | 成功後再取得 | 失敗時固定 |
|--------|-----|-----|---------|--------------------|--------------|------------|
| login / session | `POST /api/login`, `POST /api/login/totp`, `POST /api/logout`, `GET /api/sessions`, `POST /api/sessions/revoke-all` | `login()`, `loginTotp()`, `logout()`, `getSessions()`, `revokeAllSessions()` | login、TOTP 確認、logout、session 一括失効。 | `.admin_credentials`、`.access_log`、`.audit_log`、memory session。token 本体は永続化しない。 | login 成功後は `getDashboard()` → `getStatus()` → `getQueue()`。session 失効後は `getSessions()`。 | `401` は SDK token / ticket / secret field を破棄し、UI は `panel-login` のみ表示する。 |
| build control | `POST /api/build`, `POST /api/build/force`, `POST /api/build/cancel`, `GET /api/build/stream`, `GET /api/queue`, `DELETE /api/queue` | `triggerBuild()`, `buildForce()`, `cancelBuild()`, `streamBuild()`, `getQueue()`, `clearQueue()` | manual build、force build、cancel、stream 開始/停止、queue clear。 | `.build_state`、`.build_lock`、`.build_logs/{id}.json`、queue entry。force build 時だけ SHA cache 更新。 | `getStatus()` → `getQueue()`、stream end 後は `getLogs()` も実行。 | running は `409`、queue full は `429`、maintenance / circuit は `503` または仕様済み `409`。UI は同じ build request を自動再送しない。 |
| logs / history | `GET /api/logs`, `GET /api/logs/search`, `GET /api/history`, `GET /api/history/{id}/log`, `POST /api/history/{id}/comment`, `POST /api/history/{id}/flag`, `POST /api/history/{id}/tags` | `getLogs()`, `searchLogs()`, `getHistory()`, `getHistoryLog()`, `setHistoryComment()`, `setHistoryFlag()`, `setHistoryTags()` | log 表示、検索、history 表示、comment / flag / tag 保存。 | read-only GET は副作用なし。comment / flag / tag は `.build_logs/{id}.json` と、owner 詳細仕様で履歴・設定ログ更新が定義された場合に限り `.build_history`、`.config_log`。 | 変更系は `getHistory()`、comment は `getHistoryComment(id)` も実行。 | `404` は選択解除または not found 表示。`422 details` は field error。壊れた JSON Lines は response に含めない。 |
| config / repo / branch / schedule | `GET/POST /api/config`, `POST /api/config/validate`, `GET/POST /api/repo-config`, `GET/POST /api/branch-config`, `GET /api/schedule`, `POST /api/schedule/*` | `getConfig()`, `setConfig()`, `validateConfig()`, `setRepoConfig()`, `getBranchConfig()`, `setBranchConfig()`, `getSchedule()`, schedule 系 method | config 保存、validate、repo 保存、branch 保存、schedule 変更。 | `.server_config`、`.repo_config`、`.branch_config`、`.config_log`。validate は保存なし。systemd 反映失敗時も保存済み値は戻さない。 | 保存系は対象 GET → `getConfigLog()`。validate は再取得なし。 | `422` は書込前停止。systemd 失敗 `500` は保存済み値を再取得して表示する。no-op は状態ファイルと log を変更しない。 |
| notify / SMTP / webhook | `GET/POST /api/notify-config`, `POST /api/notify-test`, `GET /api/notify-log`, `GET/POST /api/smtp-config`, `POST /api/smtp-test`, `GET/POST /api/webhook-config`, `GET /api/webhook-events`, `POST /api/notify/weekly-summary` | `getNotifyConfig()`, `setNotifyConfig()`, `notifyTest()`, `getNotifyLog()`, `getSmtpConfig()`, `setSmtpConfig()`, `smtpTest()`, `getWebhookConfig()`, `setWebhookConfig()`, `getWebhookEvents()`, `notifyWeeklySummary()` | notify 保存、test、SMTP 保存/test、webhook secret 保存、event 表示、weekly summary。 | `.notify_config`、`.smtp_config`、`.smtp_secret`、`.webhook_secret`、`.notify_log`、`.webhook_events.json`、`.config_log`。 | 保存系は対象 GET → `getConfigLog()`。test / summary は `getNotifyLog()`。 | secret は response / log / fixture へ平文出力しない。保存成功・失敗とも UI secret field を消去する。 |
| snapshots / rollback / maintenance | `GET /api/snapshots`, `GET /api/snapshots/{id}/download`, `DELETE /api/snapshots/{id}`, `POST /api/history/{id}/rollback`, `GET /api/maintenance`, `POST /api/maintenance/enable`, `POST /api/maintenance/disable` | `getSnapshots()`, `downloadSnapshot()`, `deleteSnapshot()`, `rollbackHistory()`, `getMaintenance()`, `enableMaintenance()`, `disableMaintenance()` | snapshot list/download/delete、rollback、maintenance enable/disable。 | `.snapshots/`、`.build_history`、`.build_logs/{new_id}.json`、`.maintenance`、`.config_log`。download は副作用なし。 | delete は `getSnapshots()`。rollback は `getHistory()` → `getStatus()`。maintenance は `getMaintenance()`。 | delete / rollback は確認必須。running rollback は `409`。maintenance enabled 中は build / rollback / 設定変更系を disabled。 |
| access / hooks / rules / pipeline / notes / layout | access、hooks、alert rules、tag rules、pipeline config、notes、dashboard layout の GET/POST/DELETE endpoint | 対応する §23 SDK method | 保存、追加、削除、notes 保存、dashboard layout 保存。 | `.access_control`、`.hooks`、`.alert_rules`、`.tag_rules`、`.pipeline_config`、`.notes`、`.dashboard_layout`、`.config_log`。 | 対象 GET → 変更系で config log 対象の場合は `getConfigLog()`。layout は `getDashboardLayout()` → `getDashboard()`。 | duplicate `409` は競合表示。validation `422` は field error。削除対象不在は `404`。 |
| tokens / audit / access logs / rate limit | `GET /api/tokens`, `POST /api/tokens`, `DELETE /api/tokens/{id}`, `GET /api/audit-log`, `GET /api/access-log`, `GET /api/api-access-log`, `GET/POST /api/api-rate-limit` | `getTokens()`, `createToken()`, `revokeToken()`, `getAuditLog()`, `getAccessLog()`, `getApiAccessLog()`, `getApiRateLimit()`, `setApiRateLimit()` | token 発行/失効、audit / access log 表示、rate limit 保存。 | `.api_tokens`、`.audit_log`、`.access_log`、`.api_access_log`、`.api_rate_state`、`.server_config`、`.config_log`。token 本体は作成時 response のみ。 | token 操作は `getTokens()` → `getAuditLog()`。rate limit は `getApiRateLimit()`。 | token 本体は再取得不可。`403` は logout しない。`429` は rate limit 表示し、同一操作を自動 retry しない。 |

**横断処理順契約：**

| 処理種別 | 固定順序 |
|----------|----------|
| 認証必須 JSON API | method / path 判定 → body 禁止判定 → JSON parse → 認証 / scope → rate limit → endpoint 固有 validation → read → write 計画 → atomic write → JSON Lines 追記 → response。 |
| read-only API | method / path 判定 → body 禁止判定 → 認証 / scope → query validation → read → 壊れた任意行除外 → response。read-only API は状態ファイルを書き換えない。 |
| UI 変更操作 | panel error / success 消去 → UI 入力検証 → 対象操作 disabled → SDK 呼び出し → 成功後再取得 → success 表示 → secret 消去 → disabled 再評価。 |
| UI 取得操作 | panel error 消去 → 対象操作 disabled → SDK 呼び出し → DOM 更新 → empty state 判定 → disabled 再評価。success 表示は行わない。 |
| SDK request | 引数検証 → path / query / body 生成 → Authorization 付与 → timeout 設定 → fetch → status 判定 → response parse → token 変化適用 → return / throw。 |
| multi-file write | 全入力検証 → 全対象 read → 全 write payload 生成 → `docs/details/api.md` §22.0d の Write 順に atomic write → JSON Lines 追記 → response。途中失敗時は未処理ファイルを書かない。 |

§27.21〜§27.38 の実装では、owner 詳細仕様にない状態ファイル、endpoint、SDK method、UI 操作、外部公開構成を追加してはならない。追加が必要な場合は、`docs/DETAIL_INDEX.md` ではなく、該当 owner / collaborator の分割先詳細仕様を先に改訂する。
