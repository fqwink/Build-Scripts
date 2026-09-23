# Adlaire CI — UI 詳細仕様

`ui` owner component の詳細本文責務は、[`docs/details/ui.md`](ui.md) を正本とする。

owner / collaborator 境界管理は [`docs/DETAIL_INDEX.md`](../DETAIL_INDEX.md) 詳細仕様入口責務 §0b.1 に従う。`ui` owner component の主本文であり、collaborator component の仕様は SDK method、API response、security、admin 配布、fixture、検証観点として参照する。

UI が呼び出す SDK method、戻り値、error、stream、token 破棄は [`docs/details/sdk.md`](sdk.md) §23 を参照する。[`docs/details/ui.md`](ui.md) 詳細本文責務は UI 側の DOM id、panel、操作、表示状態、SDK 呼び出し、秘密情報消去を定義する。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `ui` |
| collaborator component | `sdk`、`api`、`security` |
| 持つ内容 | `ui` owner が主本文として定義する DOM id、panel、操作、表示状態、SDK 呼び出し、秘密情報消去。 |
| 持たない内容 | SDK method 実装、API endpoint 実装、状態 schema、状態ファイル直接操作、admin 静的配信、setup / release 手順、fixture 証跡責務。 |

---

## 24. 標準管理ツール 仕様

[`docs/details/ui.md`](ui.md) §24 は、ui owner の `admin/index.html` 詳細本文責務である。

**ファイル構成：**
```
/opt/adlaire-builder/admin/
├── index.html          # 管理画面（単一ファイル完結）
└── adlaire-ci-sdk.js   # SDK（標準管理ツールに同梱）
```

**DOM / section / form field 命名契約：**

標準管理ツールは、[`docs/details/ui.md`](ui.md) §24 の DOM 固定表の DOM id、`data-panel`、form field name を使用する。表にない主要パネル id、主要 form name、主要 button id を追加してはならない。表示・非表示は `hidden` 属性で制御し、DOM 要素の生成順は [`docs/details/ui.md`](ui.md) §24 の DOM 固定表の順序とする。

| パネル | section id | data-panel | 主フォーム id | 主要 field name | 主要 button id |
|--------|------------|------------|---------------|-----------------|----------------|
| ログイン | `panel-login` | `login` | `form-login` | `password`, `totp_code` | `btn-login`, `btn-login-totp` |
| パスワード変更 | `panel-password` | `password` | `form-password` | `current_password`, `new_password` | `btn-change-password` |
| ステータス | `panel-status` | `status` | なし | なし | `btn-refresh-status`, `btn-save-dashboard-layout` |
| 手動実行 | `panel-build` | `build` | なし | なし | `btn-build`, `btn-build-force`, `btn-build-cancel`, `btn-stream-close`, `btn-queue-clear` |
| 承認待ち | `panel-approvals` | `approvals` | なし | なし | `btn-refresh-approvals` |
| ログビューア | `panel-logs` | `logs` | `form-log-search` | `n`, `q`, `from`, `to`, `level` | `btn-load-logs`, `btn-search-logs`, `btn-export-logs`, `btn-cleanup-logs`, `btn-archive-logs` |
| ビルド履歴 | `panel-history` | `history` | `form-history-filter` | `page`, `per_page`, `trigger`, `tag`, `flagged` | `btn-export-history` |
| システム情報 | `panel-system` | `system` | `form-pat` | `token`, `pat_expires_at` | `btn-pat-verify`, `btn-pat-update` |
| 通知設定 | `panel-notify` | `notify` | `form-notify` | `webhooks`, `channels`, `on`, `summary`, `email`, `secret`, `smtp_password` | `btn-save-notify`, `btn-notify-test`, `btn-weekly-summary`, `btn-save-webhook-secret`, `btn-save-smtp`, `btn-smtp-test` |
| 設定 | `panel-config` | `config` | `form-config` | `log_max_lines`, `history_max_count`, `build_timeout_seconds`, `log_retention_days`, `log_archive_after_days`, `log_level`, `queue_max_size`, `snapshots_keep`, `build_retry_max`, `build_retry_base_seconds`, `commit_status_enabled`, `commit_status_context`, `commit_status_target_url`, `build_trend_keep_count`, `duration_anomaly_enabled`, `duration_anomaly_min_samples`, `duration_anomaly_avg_multiplier`, `duration_anomaly_p95_multiplier` | `btn-save-config`, `btn-validate-config`, `btn-set-log-level` |
| セキュリティ | `panel-security` | `security` | `form-security` | `totp_code`, `api_rate_enabled`, `api_rate_group`, `api_rate_window_seconds`, `api_rate_max_requests`, `session_timeout_seconds` | `btn-load-security`, `btn-totp-setup`, `btn-totp-confirm`, `btn-totp-disable`, `btn-save-api-rate-limit`, `btn-save-session-timeout` |
| アクセスログ | `panel-access-log` | `access-log` | `form-api-access-log-filter` | `limit`, `offset`, `method`, `path`, `status` | `btn-load-access-log`, `btn-load-api-access-log` |
| 監査ログ | `panel-audit-log` | `audit-log` | `form-audit-log-filter` | `limit`, `offset`, `actor`, `action`, `result` | `btn-load-audit-log` |
| 統計 | `panel-stats` | `stats` | なし | なし | `btn-load-stats` |
| リポジトリ情報 | `panel-repo` | `repo` | `form-repo` | `owner`, `repo`, `branch`, `target_file`, `interval_seconds`, `allowed_from`, `allowed_to`, `maintenance_reason` | `btn-save-repo`, `btn-save-branch-config`, `btn-pause-schedule`, `btn-resume-schedule`, `btn-enable-maintenance`, `btn-disable-maintenance` |
| セッション管理 | `panel-sessions` | `sessions` | なし | なし | `btn-load-sessions`, `btn-revoke-sessions` |
| システム診断 | `panel-diagnostics` | `diagnostics` | なし | なし | `btn-run-diagnostics`, `btn-verify-output` |
| ビルド比較 | `panel-compare` | `compare` | `form-compare` | `left_build_id`, `right_build_id` | `btn-compare-builds` |
| API トークン管理 | `panel-tokens` | `tokens` | `form-token` | `label`, `scopes`, `expires_at` | `btn-create-token` |
| 運用ノート | `panel-notes` | `notes` | `form-notes` | `content` | `btn-save-notes` |
| スナップショット | `panel-snapshots` | `snapshots` | なし | なし | `btn-load-snapshots` |
| メンテナンス | `panel-maintenance` | `maintenance` | `form-maintenance` | `reason` | `btn-maintenance-enable`, `btn-maintenance-disable` |
| アクセス制御 | `panel-access-control` | `access-control` | `form-access-control` | `allow` | `btn-save-access-control` |
| フック | `panel-hooks` | `hooks` | `form-hook` | `phase`, `command_args`, `abort_on_failure` | `btn-add-hook` |

共通領域の DOM id は、`app-root`、`nav-panels`、`global-banner`、`global-error`、`global-success`、`maintenance-banner`、`build-log-stream`、`build-queue-summary`、`issued-token-once` とする。エラー表示要素は各 panel 内に `id="{section-id}-error"`、成功表示要素は `id="{section-id}-success"` を置く。`label[for]` と input `id` は `field-{field_name}` 形式で一致させる。複数行・配列入力は `textarea` または table row で表現し、保存直前に SDK 引数の型へ変換する。

**画面構成：**

| パネル | 表示内容 | 表示条件 |
|-------|---------|---------|
| ログイン | パスワード入力フォーム。TOTP が必要な場合は同じ panel 内で TOTP code 入力へ切り替える。 | 未ログイン時のみ |
| パスワード変更 | 現在・新パスワード入力フォーム | `must_change: "prompt"` または `"forced"` 時（`"forced"` 時は他パネル非表示） |
| ステータス | 最終ビルド時刻・SHA・成否・出力サイトリンク・ダッシュボードウィジェット編集モード（表示するウィジェットをチェックボックスで選択・並び替え・保存） | ログイン済み |
| 手動実行 | ビルドトリガーボタン・SHA リセットを含む強制ビルドボタン・キャンセルボタン（実行中のみ有効）・実行結果表示・リアルタイムログ表示エリア（SSE ストリーミング）・キュー状態表示（待機中件数・クリアボタン） | ログイン済み |
| 承認待ち | 承認待ち build の一覧（id・branch・sha・target・status・created_at・expires_at）・承認ボタン・却下ボタン・再取得ボタン。pending 以外の entry は操作ボタンを disabled にする。 | ログイン済み |
| ログビューア | 最新ビルドログ（n 行・キーワードフィルター・ログレベルフィルターボタン（INFO / WARNING / ERROR / DEBUG）・JSON エクスポートボタン・横断検索フォーム（期間指定）・検索結果一覧） | ログイン済み |
| ビルド履歴 | 過去ビルド一覧（日時・SHA・成否・トリガー種別・所要時間・重要フラグ列・タグ列・ログ表示リンク・コメント入力欄）・フラグ付きのみ表示フィルター・タグフィルター・ページネーション UI・JSON エクスポートボタン | ログイン済み |
| システム情報 | 出力サイトサイズ・更新日時・稼働時間・ディスク使用量（ログ合計・出力サイト）・PAT 即時検証ボタン・PAT 更新フォーム・PAT 有効期限表示（設定フォーム・期限切れ間近で警告表示）・GitHub API レート制限表示 | ログイン済み |
| 通知設定     | Webhook 一覧（追加/削除/ラベル/有効無効切り替え/リトライ回数・間隔設定/シークレット入力欄）・通知条件設定（ビルド開始時・成功時・失敗時）・各 Webhook ペイロードテンプレート編集フォーム（変数一覧表示）・テスト送信ボタン・定期サマリー設定（間隔・時刻・曜日・即時送信ボタン）・送信履歴（試行回数・エラー内容列含む）・メール通知セクション（SMTP 設定フォーム・宛先リスト・通知条件・テスト送信ボタン） | ログイン済み |
| 設定         | ログ保持行数・履歴保持件数の設定変更・ビルドタイムアウト設定・ログレベル変更（INFO / DEBUG）・ログ保持期間（日数、0 = 無制限）・スナップショット保持世代数設定・ビルドキュー最大長設定・手動クリーンアップボタン・設定変更履歴（変更日時・項目・変更前後の値）・IP アクセス制限セクション（許可 IP / CIDR 一覧・追加フォーム・削除ボタン）・フック設定セクション（pre / post フック一覧・command_args 入力フォーム・実行ログリンク・有効無効切り替え）・アラートルール設定セクション（メトリクス・演算子・しきい値・レベル・メッセージの入力フォーム・ルール一覧・削除ボタン）・自動タグ付けルールセクション（条件式・タグ入力フォーム・ルール一覧・削除ボタン）・パイプライン設定セクション（追加引数入力欄・環境変数テーブル） | ログイン済み |
| セキュリティ | 単一 admin の TOTP 状態、TOTP setup secret / otpauth URI の一回表示、TOTP 確認/無効化、API rate limit policy、session timeout 設定。QR code 生成は初期実装対象外。 | ログイン済み |
| アクセスログ | ログイン履歴（日時・成否）・API アクセスログ（method、path、status、duration、actor、remote_addr）・API アクセスログフィルター | ログイン済み |
| 監査ログ | 設定変更、認証、token、build trigger、権限拒否の監査イベント一覧と actor / action / result フィルター。secret、password、token 本体は表示しない。 | ログイン済み |
| 統計         | ビルド回数・成功率・平均間隔・平均・最大ビルド時間・中央値・p95・異常件数・日別時系列データ（グラフ表示対応） | ログイン済み |
| リポジトリ情報 | 監視対象リポジトリ・ブランチ・ファイルの確認・設定変更フォーム（OWNER / REPO / BRANCH / TARGET_FILE）・ポーリング間隔変更フォーム・ポーリング一時停止／再開ボタン・許可時間帯設定（from〜to、解除ボタン）・メンテナンスモード有効化フォーム（理由入力）・解除ボタン・現在の状態表示 | ログイン済み |
| セッション管理 | 有効セッション一覧・全セッション強制無効化ボタン | ログイン済み |
| システム診断   | PAT・GitHub API・出力サイト・systemd・Webhook の診断項目一覧（ok / warn / error）・診断実行ボタン・アラートバッジ（`GET /api/dashboard` の `alerts` に基づき warn / error を表示）・出力整合性チェック項目（`POST /api/verify-output` 結果表示）・メンテナンスモード中はバナーを全パネル上部に表示 | ログイン済み |
| ビルド比較     | ビルド履歴から 2 件を選択・ログ並列表示・差分ハイライト（クライアントサイド処理） | ログイン済み |
| API トークン管理 | 発行済みトークン一覧（ラベル・スコープ・作成日時・最終使用日時・有効期限・失効日時）・新規発行フォーム（ラベル入力・scope 複数選択・有効期限）・発行時のみトークン文字列を表示・失効ボタン | ログイン済み |
| 運用ノート     | Markdown レンダリング表示・編集モード切替・保存ボタン・最終更新日時表示 | ログイン済み |
| スナップショット | ビルド成果物の世代一覧（最大`snapshots_keep`件）・個別ダウンロード・削除 | ログイン済み |
| メンテナンス   | メンテナンスモードの有効化（理由テキスト付き）・無効化・状態・開始時刻表示 | ログイン済み |
| アクセス制御   | 許可 IP / CIDR 一覧・CIDR 追加フォーム・削除ボタン（ブロック時は 403 を返す） | ログイン済み |
| フック         | Pre/Post ビルドフック一覧・command_args 追加・削除・実行ログ（直近 N 件）確認 | ログイン済み |

**UI パネル初期取得契約：**

各パネルを表示する時は、[`docs/details/ui.md`](ui.md) §24 の固定表の SDK method を上から順に呼び出す。表示済み panel へ再遷移した場合も、ユーザー操作で表示した時点で同じ順序で再取得する。空配列はエラーではなく空状態として表示する。

| パネル | 初期取得 SDK method | 空状態表示 |
|--------|---------------------|------------|
| ステータス | `getDashboard()`, `getDashboardLayout()`, `getQueue()` | alerts なし、queue なしを 1 行で表示する。 |
| 手動実行 | `getStatus()`, `getQueue()`, `getMaintenance()` | queue なしを 1 行で表示する。 |
| ログビューア | `getLogs(100,"")` | ログなしを 1 行で表示する。 |
| ビルド履歴 | `getHistory({page:1,perPage:20})` | 履歴なしを 1 行で表示する。 |
| システム情報 | `getSysinfo()`, `getPatStatus()`, `getRateLimit()`, `getDiskUsage()` | 取得不能項目は該当ブロックにエラー表示し、他ブロックは表示する。 |
| 通知設定 | `getNotifyConfig()`, `getWebhookConfig()`, `getSmtpConfig()`, `getNotifyLog()` | webhook / email / log なしをそれぞれ 1 行で表示する。 |
| 設定 | `getConfig()`, `getConfigLog()`, `getAccessControl()`, `getHooks()`, `getAlertRules()`, `getTagRules()`, `getPipelineConfig()` | 各一覧なしを 1 行で表示する。 |
| セキュリティ | `getTotpStatus()`, `getApiRateLimit()`, `getConfig()` | TOTP disabled、rate limit disabled は通常状態として表示する。admin 権限なしの `403` は rate limit と session timeout 編集欄だけ非表示にする。 |
| アクセスログ | `getAccessLog()`, `getApiAccessLog({limit:100,offset:0})` | ログなしを 1 行で表示する。 |
| 監査ログ | `getAuditLog({limit:100,offset:0})` | log なしを 1 行で表示する。admin 権限なしの `403` は panel を非表示にする。 |
| 統計 | `getStats(7)`, `getStatsTimeline(30)`, `getStatsBuildDuration(20)` | 統計対象なしを 1 行で表示する。 |
| リポジトリ情報 | `getRepoInfo()`, `getBranchConfig()`, `getSchedule()`, `getMaintenance()` | branch target なしを 1 行で表示する。 |
| セッション管理 | `getSessions()` | 現 session だけの場合も通常一覧として表示する。 |
| システム診断 | `getDiagnostics()` | 診断 item なしは `No diagnostics` と表示する。 |
| ビルド比較 | `getHistory({page:1,perPage:100})` | 比較対象 2 件未満は compare button を disabled にする。 |
| API トークン管理 | `getTokens()` | token なしを 1 行で表示する。 |
| 運用ノート | `getNotes()` | content 空文字は空 editor として表示する。 |
| スナップショット | `getSnapshots()` | snapshot なしを 1 行で表示する。 |
| メンテナンス | `getMaintenance()` | disabled 状態を通常表示する。 |
| アクセス制御 | `getAccessControl()` | `allow:[]` は制限なしとして表示する。 |
| フック | `getHooks()` | hook なしを 1 行で表示する。 |

UI は、初期取得で一部 API が失敗した場合、ログイン状態を維持し、該当 panel の error 領域に失敗を表示する。ただし `401` は全 panel 表示を中止してログイン画面へ戻す。

**パスワード変更フロー：**
- `must_change: "prompt"`: パスワード変更パネルを表示。他パネルも操作可能
- `must_change: "forced"`: パスワード変更パネルのみ表示。変更完了後に通常画面へ遷移

**UI 操作契約表：**

標準管理ツールは、[`docs/details/ui.md`](ui.md) §24 の固定表の SDK method 以外を直接呼び出してはならない。ファイル操作、`fetch()` の直接呼び出し、`systemctl` 実行、`runner` 直接起動は禁止する。成功時表示は対象パネル内に 1 行で表示し、失敗時表示は `AdlaireCIError.message` と `details` を同じパネル内に表示する。

| パネル | 操作 | SDK method | 成功時表示 | 成功後再取得 | disabled 条件 |
|--------|------|------------|------------|--------------|---------------|
| ログイン | ログイン | `login(password)` | 通常画面へ遷移、または TOTP 入力へ切替 | `getDashboard()`, `getStatus()`, `getQueue()` | password 空欄、送信中 |
| ログイン | TOTP 確認 | `loginTotp(ticket,code)` | 通常画面へ遷移 | `getDashboard()`, `getStatus()`, `getQueue()` | ticket なし、code 空欄、送信中 |
| パスワード変更 | 変更保存 | `changePassword(currentPassword,newPassword)` | `Password changed` | なし | 入力不足、送信中 |
| 全パネル共通 | ログアウト | `logout()` | ログイン画面へ戻る | なし | 送信中 |
| ステータス | 再読み込み | `getDashboard()` | 最終更新時刻を表示 | `getDashboard()` | 読み込み中 |
| ステータス | ウィジェット保存 | `setDashboardLayout(widgets)` | `Dashboard layout updated` | `getDashboardLayout()`, `getDashboard()` | widgets 空配列、送信中 |
| 手動実行 | ビルド開始 | `triggerBuild()` | `Build queued` または `Build started` | `getStatus()`, `getQueue()` | running true、maintenance enabled、送信中 |
| 手動実行 | 強制ビルド | `buildForce()` | `Force build queued` または `Force build started` | `getStatus()`, `getQueue()` | maintenance enabled、送信中 |
| 手動実行 | キャンセル | `cancelBuild()` | `Build cancelled` | `getStatus()`, `getQueue()` | running false、送信中 |
| 手動実行 | リアルタイムログ開始 | `streamBuild(onLine,onEnd)` | ログ行を追記 | `getStatus()`, `getQueue()` on end | token なし、接続中 |
| 手動実行 | キュー削除 | `clearQueue()` | `Queue cleared` | `getQueue()`, `getStatus()` | queue 空、送信中 |
| ログビューア | 最新ログ取得 | `getLogs(n,q)` | ログ行を表示 | なし | 読み込み中 |
| ログビューア | 横断検索 | `searchLogs(q,from,to)` | 検索結果件数を表示 | なし | 日付不正、検索中 |
| ログビューア | JSON export | `exportLogs()` | export 完了表示 | なし | 処理中 |
| ログビューア | cleanup | `cleanupLogs()` | 削除件数を表示 | `getLogs()`, `getDiskUsage()` | 処理中 |
| ログビューア | archive | `archiveLogs()` | 圧縮件数を表示 | `getLogs()`, `getDiskUsage()` | 処理中 |
| ビルド履歴 | 一覧取得 | `getHistory({page,perPage,trigger,tag,flagged})` | 件数を表示 | なし | 読み込み中 |
| ビルド履歴 | コメント保存 | `setHistoryComment(id,comment)` | `Comment saved` | `getHistory()`, `getHistoryComment(id)` | id 未選択、送信中 |
| ビルド履歴 | 重要フラグ切替 | `setHistoryFlag(id,flagged)` | `Flag updated` | `getHistory()` | id 未選択、送信中 |
| ビルド履歴 | タグ保存 | `setHistoryTags(id,tags)` | `Tags updated` | `getHistory()` | id 未選択、送信中 |
| ビルド履歴 | rollback | `rollbackHistory(id)` | `Rollback started` | `getHistory()`, `getStatus()` | running true、snapshot なし、送信中 |
| システム情報 | PAT 検証 | `patVerify()` | 検証結果を表示 | `getPatStatus()`, `getDiagnostics()` | 送信中 |
| システム情報 | PAT 更新 | `updatePat(token)` | `PAT updated` | `getPatStatus()`, `getDiagnostics()` | token 空欄、送信中 |
| 通知設定 | 通知設定保存 | `setNotifyConfig(config)` | `Notify config updated` | `getNotifyConfig()` | 入力不正、送信中 |
| 通知設定 | Webhook test | `notifyTest()` | 送信結果を表示 | `getNotifyLog()` | webhook 未設定、送信中 |
| 通知設定 | 週次 summary 送信 | `notifyWeeklySummary()` | success/failure 件数を表示 | `getNotifyLog()` | 宛先なし、送信中 |
| 通知設定 | Webhook secret 保存 | `setWebhookConfig(secret)` | `Webhook secret updated` | `getWebhookConfig()` | secret 空欄、送信中 |
| 通知設定 | SMTP 保存 | `setSmtpConfig(config)` | `SMTP config updated` | `getSmtpConfig()`, `getNotifyConfig()` | 入力不正、送信中 |
| 通知設定 | SMTP test | `smtpTest()` | 送信結果を表示 | `getNotifyLog()` | SMTP disabled、送信中 |
| 設定 | サーバー設定保存 | `setConfig(config)` | `Config updated` | `getConfig()`, `getConfigLog()` | 入力不正、送信中 |
| 設定 | サーバー設定検証 | `validateConfig(config)` | 検証結果を表示 | なし | 入力不正、送信中 |
| 設定 | log level 変更 | `setLogLevel(level)` | `Log level changed` | `getConfig()` | level 未選択、送信中 |
| 設定 | access control 保存 | `setAccessControl(allowList)` | `Access control updated` | `getAccessControl()`, `getConfigLog()` | CIDR 不正、送信中 |
| 設定 | hook 追加/削除 | `addHook()`, `deleteHook(id)` | hook 件数を表示 | `getHooks()`, `getConfigLog()` | 入力不正、送信中 |
| 設定 | alert rule 追加/削除 | `addAlertRule()`, `deleteAlertRule(id)` | rule 件数を表示 | `getAlertRules()`, `getDashboard()` | 入力不正、送信中 |
| 設定 | tag rule 追加/削除 | `addTagRule()`, `deleteTagRule(id)` | rule 件数を表示 | `getTagRules()` | 入力不正、送信中 |
| 設定 | pipeline config 保存 | `setPipelineConfig(config)` | `Pipeline config updated` | `getPipelineConfig()`, `getConfigLog()` | 入力不正、送信中 |
| セキュリティ | TOTP setup | `setupTotp()` | secret と otpauth URI を一回表示 | なし | TOTP enabled、送信中 |
| セキュリティ | TOTP confirm | `confirmTotp(code)` | `TOTP enabled` | `getTotpStatus()`, `getAuditLog()` | code 空欄、送信中 |
| セキュリティ | TOTP disable | `disableTotp(code)` | `TOTP disabled` | `getTotpStatus()`, `getAuditLog()` | code 空欄、送信中 |
| セキュリティ | API rate limit 保存 | `setApiRateLimit(policy)` | `API rate limit updated` | `getApiRateLimit()`, `getAuditLog()` | 入力不正、送信中 |
| セキュリティ | session timeout 保存 | `setConfig({session_timeout_seconds})` | `Session timeout updated` | `getConfig()`, `getConfigLog()` | 範囲外、送信中 |
| アクセスログ | 取得 | `getAccessLog()` | 件数を表示 | なし | 読み込み中 |
| アクセスログ | API アクセスログ取得 | `getApiAccessLog({limit,offset,method,path,status})` | 件数を表示 | なし | 読み込み中 |
| 監査ログ | 取得 | `getAuditLog({limit,offset,actor,action,result})` | 件数を表示 | なし | 読み込み中 |
| 統計 | 取得 | `getStats()`, `getStatsTimeline()`, `getStatsBuildDuration()` | グラフと数値を更新 | なし | 読み込み中 |
| リポジトリ情報 | repo 保存 | `setRepoConfig(config)` | `Repo config updated` | `getRepoInfo()`, `getConfigLog()` | 入力不正、送信中 |
| リポジトリ情報 | branch config 保存 | `setBranchConfig(branches)` | `Branch config updated` | `getBranchConfig()`, `getConfigLog()` | 入力不正、送信中 |
| リポジトリ情報 | schedule 変更 | `setScheduleInterval()`, `pauseSchedule()`, `resumeSchedule()`, `setAllowedHours()`, `clearAllowedHours()`, `setForceInterval()`, `setBuildCooldown()` | schedule 変更結果を表示 | `getSchedule()`, `getConfigLog()` | 入力不正、送信中 |
| セッション管理 | 全 revoke | `revokeAllSessions()` | revoke 件数を表示 | `getSessions()` | 送信中 |
| システム診断 | 診断実行 | `getDiagnostics()` | 診断結果を表示 | なし | 読み込み中 |
| システム診断 | 出力 checksum 検証 | `verifyOutput()` | match 結果を表示 | なし | 送信中 |
| API トークン管理 | token 発行 | `createToken(label,scopes,expiresAt)` | token 本体を 1 回だけ表示 | `getTokens()`, `getAuditLog()` | label 空欄、scope 未選択、送信中 |
| API トークン管理 | token 失効 | `revokeToken(id)` | `Token revoked` | `getTokens()` | id 未選択、送信中 |
| 運用ノート | 保存 | `setNotes(content)` | `Notes updated` | `getNotes()` | 送信中 |
| スナップショット | 一覧取得 | `getSnapshots()` | 件数を表示 | なし | 読み込み中 |
| スナップショット | download | `downloadSnapshot(id)` | browser download を開始 | なし | id 未選択、送信中 |
| スナップショット | 削除 | `deleteSnapshot(id)` | `Snapshot deleted` | `getSnapshots()` | id 未選択、送信中 |
| メンテナンス | 有効化 | `enableMaintenance(reason)` | `Maintenance mode enabled` | `getMaintenance()`, `getDashboard()` | reason 空欄、送信中 |
| メンテナンス | 無効化 | `disableMaintenance()` | `Maintenance mode disabled` | `getMaintenance()`, `getDashboard()` | 送信中 |

**UI 共通動作契約：**

| 項目 | 仕様 |
|------|------|
| 初期表示 | `localStorage` から token を復元しない。画面読み込み時は未ログイン状態から開始する。 |
| API 経路 | UI は必ず `AdlaireCI` SDK method を呼び出す。`fetch()`、`XMLHttpRequest`、`EventSource`、`ReadableStream` reader の直接生成は禁止する。ただし SDK の `streamBuild()` が返した `StreamHandle.close()` を呼ぶ操作は許可する。 |
| API 呼び出し中 | 対象ボタンを disabled にし、同一操作の二重送信を防ぐ。完了または失敗後に元へ戻す。 |
| 成功表示 | 変更系操作は成功時にパネル内へ 1 行の成功メッセージを表示し、関連 GET API を再取得する。 |
| 失敗表示 | SDK が投げた `AdlaireCIError.message` をパネル内エラー領域に表示する。`details` が配列の場合は各 `field` のフォーム項目に `message` を紐付け、該当項目が存在しない場合はパネル内エラー領域へ箇条書きで表示する。`responseBody`、token、secret、PAT は表示しない。 |
| `401` | token を破棄し、ログインパネルへ戻す。直前の入力値のうち秘密情報は消去する。 |
| `403` | 操作権限なしとしてエラー表示し、ログアウトはしない。 |
| `409` | 状態競合としてエラー表示し、ステータス・キュー・スケジュールを再取得する。 |
| `422` | 入力エラーとして該当フォーム項目へエラーを紐付ける。 |
| `429` | rate limit または login lock としてエラー表示し、同一操作ボタンを 10 秒間 disabled にする。ログインロックの場合は password field を空にする。 |
| `503` | メンテナンスバナーを表示し、ビルド操作ボタンを disabled にする。 |
| 秘密情報入力 | password、TOTP code、TOTP secret、PAT、Webhook Secret、SMTP password、発行直後 token は画面遷移、成功表示、再取得後にフォーム値から消去する。 |
| 自動更新 | ステータス、キュー、SSE 以外のパネルは自動ポーリングしない。ユーザー操作または画面表示時に取得する。 |
| SSE 切断 | `streamBuild()` が error になった場合はリアルタイム表示を停止し、`getStatus()` と `getQueue()` を再取得する。ユーザーが停止した場合はエラー表示しない。 |
| フォーム保存 | 保存 API が成功するまで UI 上の表示値を確定表示にしない。失敗時は入力値を保持する。 |
| 入力検証 | UI は送信前に必須入力、数値範囲、配列空、URL、CIDR、日付形式を検証する。UI 検証に通っても API 側検証は省略しない。 |
| 破壊的操作 | snapshot 削除、token 失効、queue clear、session revoke all、rollback はクリック後に確認ダイアログを 1 回表示する。確認文には対象 ID または件数を含める。 |
| 表示時刻 | API から受け取った UTC ISO 8601 はブラウザのローカル時刻で表示する。ただし data 属性または title 属性に元の ISO 8601 文字列を保持する。 |
| 一覧の空状態 | 配列が空の場合は、空表ではなくパネル内に 1 行の空状態メッセージを表示する。空状態はエラーとして扱わない。 |
| focus / aria | `422` は最初の invalid field へ focus する。`401` はログイン password field へ focus する。SSE ログ領域は `aria-live="polite"` とし、エラー領域は `role="alert"` とする。 |

**UI 初期ロード / イベント処理順序：**

標準管理ツールは、`DOMContentLoaded` 後に以下の順で初期化する。順序を入れ替えてはならない。

1. `app-root`、`nav-panels`、`global-error`、`global-success`、各 `panel-*` の存在を検査する。欠落時は `global-error` に `UI initialization failed` を表示し、以降の API 呼び出しを行わない。
2. `window.AdlaireCI` 等の global 参照を使わず、`./adlaire-ci-sdk.js` から `AdlaireCI` と `AdlaireCIError` を ES Module import する。
3. `AdlaireCI` を `new AdlaireCI({baseUrl})` で 1 回だけ生成する。`baseUrl` は同一 origin の `/api` を既定値とし、外部 origin は ui 詳細本文責務では許可しない。
4. すべての panel を `hidden=true` にし、`panel-login` だけを表示する。
5. form submit と button click の event listener を登録する。登録対象は [`docs/details/ui.md`](ui.md) §24 の DOM / section / form field 命名契約表の id に限定する。
6. `localStorage`、`sessionStorage`、Cookie から token を読み込まない。
7. `global-error`、`global-success`、各 panel error/success を空にする。
8. login password field へ focus する。

ログイン成功後の初期取得順は、`getDashboard()` → `getStatus()` → `getQueue()` → `getMaintenance()` → `getDashboardLayout()` とする。途中で `401` を受信した場合は残りの取得を中止してログイン画面へ戻す。`getMaintenance()` が `enabled=true` を返した場合は `maintenance-banner` を表示し、ビルド開始、強制ビルド、rollback、hook 追加、設定変更系ボタンを disabled にする。

イベント処理は、各操作につき以下の順で行う。

1. 対象 panel の error/success を空にする。
2. UI 側入力検証を行う。失敗時は SDK method を呼ばない。
3. 対象 button と同一操作グループを disabled にする。
4. SDK method を呼ぶ。
5. 成功時は成功メッセージを表示し、[`docs/details/ui.md`](ui.md) §24 UI 操作契約表の成功後再取得を左から順に実行する。
6. 失敗時は `AdlaireCIError` として表示する。`TypeError` は UI 実装エラーとして `global-error` に `Client error` を表示する。
7. 秘密情報 field を消去する。
8. disabled を解除する。ただし `401`、`503`、SSE 接続中、メンテナンス中、または仕様上 disabled 条件が継続する場合は解除しない。

秘密情報 field は、`password`、`current_password`、`new_password`、`totp_code`、`token`、`secret`、`smtp_password`、`issued-token-once`、`totp-secret-once` とする。これらは成功、失敗、画面遷移、`401`、`logout()`、`revokeAllSessions()` のいずれの場合も DOM 値を空にする。発行直後 token は `issued-token-once`、TOTP setup secret は `totp-secret-once` に 1 回だけ表示し、次の任意の user action で消去する。

**UI 操作完全性検証契約：**

標準管理ツールの詳細実装確認では、[`docs/details/ui.md`](ui.md) §24 の DOM / section / form field 命名契約表と UI 操作契約表を照合し、[`docs/details/ui.md`](ui.md) §24 の固定表を満たす。

| 検証項目 | 合格条件 |
|----------|----------|
| panel coverage | DOM / section / form field 命名契約表の `section id` がすべて `index.html` に存在する。 |
| button coverage | UI 操作契約表の各操作に対応する button または form submit が存在し、event listener が 1 つだけ登録される。 |
| SDK only | UI 操作契約表の SDK method 以外を UI から呼び出していない。直接 `fetch()`、`XMLHttpRequest`、`EventSource` を使用していない。 |
| success refresh | 成功後再取得列に複数 method がある場合、左から順に await し、途中失敗時は残りを中止して error 表示する。 |
| disabled restore | 操作失敗時も、継続条件がない限り disabled を解除する。`401`、`503`、SSE 接続中、メンテナンス中は解除しない。 |
| secret clearing | [`docs/details/ui.md`](ui.md) §24 の秘密情報 field が、成功、失敗、画面遷移、`401`、logout、revoke all の全経路で空になる。 |
| empty state | UI パネル初期取得契約の空状態表示が、各 panel 内に 1 行で表示される。 |
| global error | 初期化失敗、SDK constructor 失敗、想定外 `TypeError` は `global-error` に固定文言 `Client error` または `UI initialization failed` を表示する。 |

**UI 状態遷移固定契約：**

| 状態 | 表示 / 処理 |
|------|-------------|
| 未ログイン | `panel-login` だけを表示し、nav item は disabled。API token、session token、TOTP ticket を永続化しない。 |
| password 成功 / TOTP 必須 | login password field を消去し、同じ login panel 内で TOTP field を表示する。`ticket` はメモリだけに保持する。 |
| login 完了 | `panel-status` を表示し、ログイン成功後の初期取得順を実行する。 |
| must_change prompt | パスワード変更 panel を表示するが、他 panel 操作を許可する。 |
| must_change forced | パスワード変更 panel 以外を hidden または disabled にする。 |
| maintenance enabled | `maintenance-banner` を表示し、build、force build、rollback、hook 追加、設定変更系操作を disabled にする。 |
| session expired | token と ticket を破棄し、秘密情報 field を消去し、`panel-login` に戻す。 |
| fatal UI init error | API 呼び出しを行わず、`global-error` に `UI initialization failed` を表示する。 |

**UI 成功後再取得失敗契約：**

成功後再取得列に複数 SDK method がある操作では、成功 message を先に表示せず、再取得がすべて成功した後に成功 message を表示する。途中の再取得が失敗した場合は、対象操作自体は成功済みとして扱い、`global-success` に操作成功の固定文言、該当 panel error に再取得失敗を表示する。再取得失敗を理由に同じ変更 API を自動再実行してはならない。

**UI 操作状態固定契約：**

標準管理ツールは、同一操作の多重実行、API 成功前の確定表示、秘密情報の残存を防ぐため、各操作を [`docs/details/ui.md`](ui.md) §24 の固定表の状態で管理する。

| 状態 | 開始条件 | UI 表示 | 許可される遷移 |
|------|----------|---------|----------------|
| `idle` | 初期状態、前操作完了後 | 操作可能。継続 disabled 条件がある場合は対象だけ disabled。 | `validating` |
| `validating` | ユーザー操作直後 | panel error / success を空にする。field error を初期化する。 | `blocked` / `sending` |
| `blocked` | UI 入力検証失敗、確認 dialog cancel | SDK method を呼ばない。該当 field error または表示変更なし。 | `idle` |
| `sending` | SDK method 呼び出し開始 | 同一操作 button / submit だけ disabled。spinner または loading 表示は対象 panel 内だけ。 | `refreshing` / `failed` |
| `refreshing` | 変更 API 成功後、成功後再取得が必要 | success を未表示のまま、再取得を左から順に実行する。 | `succeeded` / `refresh_failed` |
| `succeeded` | 変更 API と必要な再取得が成功 | panel success を 1 行表示し、secret field を消去する。 | `idle` |
| `refresh_failed` | 変更 API は成功したが再取得に失敗 | `global-success` に操作成功、panel error に再取得失敗を表示する。secret field を消去する。 | `idle` |
| `failed` | SDK method が `AdlaireCIError` を投げる | panel error と field error を表示する。secret field を消去する。 | `idle` / `unauthorized` |
| `unauthorized` | `401` | token / ticket / secret field を消去し、全 panel を hidden、`panel-login` だけ表示する。 | `idle` |

同一 panel に複数操作がある場合でも、`sending` による disabled は同一操作グループに限定する。ただし [`docs/details/ui.md`](ui.md) §24 の disabled 優先順位で maintenance、forced password change、SSE 接続中、`401` が上位条件として残る場合は、その上位条件に従う。UI は `sending` または `refreshing` の間、同じ SDK method を再実行してはならない。

秘密情報 field は、`blocked` のうち確認 dialog cancel を除き、`succeeded`、`refresh_failed`、`failed`、`unauthorized`、panel 遷移、logout、revoke all のいずれでも空にする。API 成功前に入力欄以外の確定表示、一覧更新、badge 更新、設定値反映を行ってはならない。入力中の form 値は、secret を除き、`failed` と `refresh_failed` では保持する。

**UI DOM 更新契約：**

| 項目 | 仕様 |
|------|------|
| 一覧描画 | API response 配列の順序を保持する。UI 独自 sort は行わない。 |
| 件数表示 | `total` がある API は `total` を表示する。配列長を total の代替にしない。 |
| 日時表示 | 表示テキストはローカル時刻でよいが、`datetime` または `title` に元の UTC ISO 8601 を保持する。 |
| disabled | 送信中 disabled は操作単位で行う。同一 panel の無関係 button は disabled にしない。ただし maintenance、forced password change、SSE 接続中は仕様上の対象をまとめて disabled にする。 |
| error 領域 | panel ごとに 1 つの error 領域を使う。field error は該当 field に紐付け、panel error にも summary を 1 行表示する。 |
| success 領域 | 変更系操作成功時だけ更新する。GET 再読み込みだけでは success を表示しない。 |
| secret one-time 表示 | issued token と TOTP secret は専用領域に 1 回だけ表示し、任意の次 user action、panel 遷移、logout、`401` で消去する。 |

**破壊的操作確認文言：**

| 操作 | 確認文 |
|------|--------|
| snapshot 削除 | `Delete snapshot {id}?` |
| token 失効 | `Revoke token {id}?` |
| queue clear | `Clear {count} queued builds?` |
| session revoke all | `Revoke all other sessions?` |
| rollback | `Rollback from build {id}?` |

確認 dialog で cancel した場合は SDK method を呼ばず、success / error 表示を変更しない。

**UI fixture 固定：**

| fixture | 入力 | 合格条件 |
|---------|------|----------|
| login totp | `login()` が `totp_required:true` を返す | password 消去、TOTP field 表示、ticket は DOM に表示しない。 |
| forced password | `must_change:"forced"` | password panel 以外が操作不可。変更成功後に通常初期取得を行う。 |
| refresh failure | 変更 API 成功後の再取得 2 件目が失敗 | 変更成功は維持し、再取得失敗だけ panel error に表示する。 |
| destructive cancel | 確認 dialog cancel | SDK method 呼び出し 0 回、表示差分なし。 |
| secret clearing | token 発行、TOTP setup、PAT 更新、Webhook secret 保存 | 次 user action または遷移で秘密情報 field と一回表示が消える。 |

**Phase 3 UI 操作固定契約：**

Phase 3 UI は、ビルド状態確認、手動ビルド、強制ビルド、キャンセル、SSE ログ表示、履歴、ログ、キュー、circuit breaker reset だけを最小運用操作として固定する。UI は SDK response に存在しない状態を推測せず、API / SDK の error status と message に基づいて表示を分岐する。

| 操作 | 使用 SDK method | 成功時表示 | 成功後再取得 | 失敗時表示 / disabled |
|------|-----------------|------------|--------------|------------------------|
| 初期状態取得 | `getStatus()`, `getQueue()` | status badge、last build、queue 件数を表示する。 | なし | `500` は status panel error。`401` は login panel へ戻す。 |
| 手動 build | `triggerBuild()` | `Build queued` または `Build started` を表示する。 | `getStatus()`, `getQueue()` | `409` は競合表示後に status/queue 再取得。`429` は build button を 10 秒 disabled。`503` は maintenance/circuit 表示。 |
| 強制 build | `buildForce()` | `Force build queued` または `Force build started` を表示する。 | `getStatus()`, `getQueue()` | `409`、`429`、`503` は手動 build と同じ扱い。 |
| cancel | `cancelBuild()` | `Build cancel requested` を表示する。 | `getStatus()`, `getQueue()` | `409` は「実行中 build なし」として表示し、status/queue を再取得する。 |
| SSE 表示開始 | `streamBuild(onLine,onEnd)` | 接続中は stream indicator を表示し、受信行を append する。 | `onEnd` 後に `getStatus()`, `getQueue()`, `getLogs()` | invalid frame / network error は stream error 表示後に status/queue 再取得。ユーザー停止は error 表示なし。 |
| ログ取得 | `getLogs(n,q)` | 行順を保持して log panel に表示する。空配列は空状態表示。 | なし | `500` は log panel error。検索条件は保持する。 |
| 履歴取得 | `getHistory(options)` | `total`、`pages`、history 行を API 順序で表示する。 | なし | `422` は該当 filter field、`500` は history panel error。 |
| 履歴 log 表示 | `getHistoryLog(id)` | 選択 build の log detail を表示する。 | なし | `404` は選択解除して not found 表示。`500` は detail error。 |
| キュー取得 | `getQueue()` | queue 件数、max size、各 entry を API 順序で表示する。 | なし | `500` は queue panel error。 |
| circuit reset | `resetCircuitBreaker()` | `Circuit breaker reset` を表示する。 | `getStatus()`, `getQueue()` | `500` は circuit error。成功後に build を自動開始しない。 |

Phase 3 UI の disabled 条件は以下に固定する。

| 条件 | disabled 対象 | 解除条件 |
|------|---------------|----------|
| SDK method 実行中 | 同一操作 button と同一 form submit | 成功または失敗後。ただし次の継続条件がある場合は解除しない。 |
| `streamBuild()` 接続中 | stream 開始 button、manual build button、force build button | `end` frame、stream error、またはユーザー停止。 |
| `getStatus().running=true` | manual build button。force build は仕様上許可される場合のみ有効。rollback は disabled。 | 次回 `getStatus().running=false`。 |
| `429` | 同一操作 button | 10 秒経過後に status/queue 再取得し、継続条件がなければ解除する。 |
| `503` maintenance / circuit | build、force build、cancel 以外の状態変更操作。circuit reset は有効。 | maintenance disabled または circuit reset 成功後の再取得。 |
| `401` | 全 authenticated 操作 | login 成功後。 |

**Phase 3 UI fixture 固定：**

| fixture | fake SDK 入力 | 合格条件 |
|---------|---------------|----------|
| ui phase3 initial status error | `getStatus()` が `AdlaireCIError(status=500,message="State file is corrupted")` | status panel error に固定 message を表示し、build button を成功扱いにしない。 |
| ui phase3 manual build conflict | `triggerBuild()` が `409 Conflict` | error 表示、`getStatus()` と `getQueue()` をこの順で再取得、同じ build request を再送しない。 |
| ui phase3 queue full | `triggerBuild()` が `429 queue_full` | build button を 10 秒 disabled、password や secret field は変更しない。 |
| ui phase3 stream success | `streamBuild()` が log 2 件と end 1 件を返す | log 行 2 件を append、end 後に status、queue、logs を順に再取得、stream indicator を消す。 |
| ui phase3 stream user close | ユーザーが `StreamHandle.close()` を押す | error 表示なし、closed 表示、status/queue 再取得あり。 |
| ui phase3 history validation | `getHistory()` が `422 details` を返す | 該当 filter field に message を紐付け、history rows を前回表示のまま維持する。 |
| ui phase3 log not found | `getHistoryLog(id)` が `404 Not found` | detail panel に not found を表示し、履歴一覧は再取得しない。 |
| ui phase3 circuit reset | `resetCircuitBreaker()` 成功 | circuit 表示を閉じ、status/queue を再取得し、build を自動開始しない。 |
| ui phase3 unauthorized | 任意操作が `401` | token/ticket/secret field を消去し、`panel-login` だけ表示する。 |

**Phase 4 UI 操作固定契約：**

Phase 4 UI は、[`docs/details/ui.md`](ui.md) §24 UI 操作契約表の SDK method だけを呼び出す。UI は API / SDK response の補完、状態ファイル直接操作、未定義 endpoint 呼び出し、保存成功前の確定表示を行ってはならない。

| 機能群 | 主操作 | 成功時表示 | 成功後再取得 | 失敗時表示 / disabled |
|--------|--------|------------|--------------|------------------------|
| config / repo / branch | config 保存、repo 保存、branch config 保存、config validate | API message を表示する。validate は `valid` と errors / warnings を表示する。 | 保存系は `getConfig()` または対象 GET と `getConfigLog()`。validate は再取得なし。 | `422` は field error。`500` は panel error。入力値は保持する。 |
| schedule | interval、pause、resume、allowed hours、force interval、cooldown | 変更後値を表示する。 | `getSchedule()`, `getConfigLog()` | systemd 失敗 `500` は schedule panel error とし、再取得で保存済み値を表示する。 |
| diagnostics / dashboard | diagnostics、dashboard、rate limit、disk、output meta 取得 | item ごとの ok/warn/error を表示する。 | なし | HTTP error は panel error。item warn/error を HTTP error として扱わない。 |
| notify / SMTP / webhook | notify config 保存、webhook secret 保存、SMTP 保存、test、weekly summary、webhook events 表示 | API message、送信結果、件数を表示する。 | 保存系は対象 GET と `getConfigLog()`。test / summary は `getNotifyLog()`。 | secret 入力は成功・失敗の両方で消去する。未設定 `422` / `501` は panel error。 |
| snapshots / rollback | snapshot list、download、delete、rollback | list 件数、download 開始、delete 完了、rollback 開始を表示する。 | delete は `getSnapshots()`。rollback は `getHistory()`, `getStatus()`。 | delete / rollback は確認 dialog 必須。running `409` は status 再取得。 |
| maintenance / access / hooks | maintenance enable/disable、access 保存、hook 追加/削除 | 固定成功文言と件数または状態を表示する。 | `getMaintenance()` / `getAccessControl()` / `getHooks()` と `getConfigLog()`。 | maintenance enabled 中は build / rollback / 設定変更系を disabled。hook 追加失敗時は command 入力を保持する。 |
| alert / tag / pipeline / notes / layout | rule 追加/削除、pipeline 保存、notes 保存、dashboard layout 保存 | 固定成功文言を表示する。 | 対象 GET。dashboard layout 保存後は `getDashboard()`、設定ログ対象操作後は `getConfigLog()`。 | duplicate `409` は競合表示。validation `422` は field error。no-op は成功表示のみ。 |
| tokens / sessions / audit | token 発行/失効、session revoke、audit/API access log 表示 | token 発行時は token 本体を一回表示する。失効/revoke は固定成功文言。 | token 操作は `getTokens()`, `getAuditLog()`。session revoke は `getSessions()`。 | token 本体は次 user action、panel 遷移、logout、`401` で消去する。`403` は logout しない。 |

Phase 4 UI の秘密情報消去条件は以下に固定する。

| 対象 field / 表示 | 消去タイミング |
|-------------------|----------------|
| PAT、Webhook Secret、SMTP password | 保存成功、保存失敗、panel 遷移、logout、`401`。 |
| 発行直後 API token | 次 user action、copy button 押下後、panel 遷移、logout、`401`。 |
| TOTP secret / ticket / code | confirm 成功、confirm 失敗、panel 遷移、logout、`401`。 |
| password / current_password / new_password | login / change 成功、login / change 失敗、logout、`401`。 |

**Phase 4 UI fixture 固定：**

| fixture | fake SDK 入力 | 合格条件 |
|---------|---------------|----------|
| ui phase4 config validation | `setConfig()` が `422 details` | 該当 field に error、panel error summary 1 行、入力値保持、`getConfig()` を呼ばない。 |
| ui phase4 schedule save failure | `setScheduleInterval()` が `500` | panel error 表示後に `getSchedule()` を 1 回呼び、保存済み値を表示する。 |
| ui phase4 secret save failure | `setWebhookConfig()` または `setSmtpConfig()` が `500` | secret field を消去し、secret 平文を error 表示しない。 |
| ui phase4 notify test | `notifyTest()` 成功 | 結果表示後に `getNotifyLog()` を呼び、通知設定を自動保存しない。 |
| ui p4 snapshot delete cancel | delete 確認 dialog cancel | SDK method 呼び出し 0 回、success / error 表示差分なし。 |
| ui p4 rollback conflict | `rollbackHistory()` が `409 Build is running` | error 表示、`getStatus()` を呼ぶ、rollback request を再送しない。 |
| ui p4 maintenance enabled | `getMaintenance()` が enabled | maintenance banner 表示、build / rollback / 設定変更系 disabled、disable maintenance は enabled。 |
| ui p5 token issue once | `createToken()` 成功 | token 本体を一回表示し、`getTokens()` 後の一覧には token 本体を表示しない。 |
| ui p5 duplicate rule | `addAlertRule()` が `409 Conflict` | 競合表示、rule list は前回表示を保持し、自動 retry しない。 |
| ui p5 layout invalid | `setDashboardLayout()` が `422 details` | 該当 widget field error、dashboard 表示順を変更しない。 |

**UI 表示データ固定契約：**

| 表示対象 | 表示元 | 表示順 | 補完禁止 |
|----------|--------|--------|----------|
| dashboard widget | `getDashboardLayout().widgets`、`getDashboard()` | layout の `widgets` 順。未知 widget は表示せず、panel error に `Unknown dashboard widget` を 1 行表示する。 | UI が widget を自動追加、並び替え、既定復元してはならない。 |
| build history | `getHistory()` | API response の `history` 配列順。 | UI 独自 sort、欠落 duration の算出、status 名の言い換えは禁止。 |
| build compare | `getHistory({page:1,perPage:100})`、選択後 `getHistoryLog(left)`, `getHistoryLog(right)` | 左選択、右選択の順。ログ行は各 response の行順。 | API にない diff 結果を保存しない。比較結果は DOM 上の一時表示だけとする。 |
| approvals | `getApprovals()` | API response の `approvals` 配列順。`pending` 以外は操作 button disabled。 | 期限切れ判定を UI 時刻だけで確定しない。API status を基準とする。 |
| tokens | `getTokens()`、`createToken()` | 一覧は API 配列順。発行直後 token は `issued-token-once` だけへ表示する。 | `GET /api/tokens` の record に token 本体を合成しない。 |
| notes | `getNotes()` | `content` を editor へそのまま入れる。表示 preview は HTML escape 後の簡易 Markdown 表示に限定する。 | UI が保存前に trim、整形、Markdown 拡張を行わない。 |
| hooks / rules | `getHooks()`、`getAlertRules()`、`getTagRules()` | API 配列順。 | UI 側で重複排除、無効化推測、command 文字列結合を行わない。 |
| pipeline config | `getPipelineConfig()` | `extra_args` と `env` を response 順で表示する。 | reserved arg の削除、env key の補完、inline YAML の再整形を行わない。 |

**UI 入力正規化固定契約：**

| 入力 | UI 正規化 | SDK 送信値 | 禁止条件 |
|------|-----------|------------|----------|
| 数値 field | ASCII 数字だけを整数化する。空文字は未指定として扱う。 | number または未指定。 | `Number("") == 0` 扱い、範囲外丸め、自動既定値保存は禁止。 |
| checkbox 群 | checked の DOM 順で配列化する。 | string array または boolean。 | 未選択時に UI が勝手に全選択へ戻さない。 |
| comma separated tag / scope | `,` で分割し、各要素の前後 ASCII whitespace だけ除去し、空要素を捨てる。 | string array。 | 重複削除、大小文字変換、未知値削除は API に委ねる。 |
| URL / webhook / SMTP host | 前後 whitespace を除去する。 | string。 | scheme 補完、host 置換、password 埋め込みは禁止。 |
| CIDR / IP | 前後 whitespace を除去し、空行を捨てる。 | string array。 | UI 側で CIDR 正規化や範囲展開を行わない。 |
| command_args | 1 行 1 引数として配列化する。空行は捨てる。 | string array。 | shell 文字列結合、quote 展開、環境変数展開は禁止。 |
| notes content | 入力値をそのまま送る。 | `{content}`。 | trim、改行正規化、Markdown 整形は禁止。 |
| date / expires_at | 空文字は `null`、入力ありは browser が返す ISO 互換文字列を送る。 | string/null。 | UI が現在時刻を補完しない。 |

**UI error / disabled 優先順位固定：**

| 優先 | 条件 | UI 処理 |
|------|------|---------|
| 1 | fatal UI init error | 全 API 呼び出しを停止し、`global-error` に `UI initialization failed`。 |
| 2 | `401` | token / ticket / secret field を消去し、全 panel を hidden、`panel-login` だけ表示。 |
| 3 | must_change forced | password panel 以外を hidden または disabled。 |
| 4 | maintenance enabled | build、force build、rollback、hook 追加、設定変更系を disabled。maintenance disable は有効。 |
| 5 | SSE 接続中 | stream 開始、manual build、force build を disabled。stream close は有効。 |
| 6 | 対象 SDK method 実行中 | 同一操作 button / submit だけ disabled。 |
| 7 | `429` | 同一操作 button を 10 秒 disabled。10 秒後に必要な再取得を行い、上位条件がなければ解除。 |
| 8 | `422 details` | 最初の invalid field へ focus し、field error と panel summary を表示。 |
| 9 | `403` / `409` / `500` | panel error を表示し、仕様上の再取得だけ実行。logout や自動 retry は行わない。 |

上位条件が残っている場合、下位条件の解除処理で button を有効化してはならない。複数 error が同時に発生した場合は、最上位の条件だけを global 表示し、下位の詳細は対象 panel error に残す。

**UI 詳細 fixture 固定：**

| fixture | fake SDK 入力 | 合格条件 |
|---------|---------------|----------|
| ui dashboard unknown widget | `getDashboardLayout()` が `["status","unknown","stats"]` を返す。 | `status`、`stats` だけ表示し、順序保持。`unknown` は panel error 1 行。layout 保存を自動実行しない。 |
| ui compare two builds | history 2 件選択後、左右の `getHistoryLog()` が異なる stdout を返す。 | 左右ログを別 column で API 行順表示し、差分 class は DOM 一時表示だけ。状態保存 API を呼ばない。 |
| ui compare missing build | 右側 `getHistoryLog()` が `404`。 | compare panel error、左側表示は維持、history 再取得なし、選択値は保持。 |
| ui approvals expired | `getApprovals()` が `status:"expired"` を含む。 | approve / reject button disabled、期限切れ表示、UI が pending へ戻さない。 |
| ui approval approve conflict | `approveBuild(id)` が `409`。 | error 表示後に `getApprovals()` を 1 回呼び、同じ approve を再送しない。 |
| ui notes preserve content | notes に前後空白と連続改行を含めて保存。 | `setNotes(content)` へ入力値そのまま送信し、trim しない。 |
| ui hook command args | 3 行の command args を入力し、中央行が空。 | 空行を除いた配列を `addHook()` に渡し、shell 文字列を作らない。 |
| ui pipeline reserved arg | `setPipelineConfig()` が `422 details`。 | field error を表示し、入力値を保持し、`getPipelineConfig()` を呼ばない。 |
| ui token issued clear | `createToken()` が token 本体を返す。 | `issued-token-once` に 1 回表示し、次 user action で消去。`getTokens()` の一覧に token 本体を表示しない。 |
| ui disabled priority | maintenance enabled 中に `429` が発生し 10 秒経過。 | maintenance が継続する限り build / rollback / 設定変更系は disabled のまま。 |

**§27.21〜§27.47 UI 連動実装確認固定契約：**

[`docs/details/ui.md`](ui.md) §27.21〜§27.47 の追加仕様化機能で UI の詳細実装確認を満たすには、[`docs/details/ui.md`](ui.md) §24 の DOM / section / form field 命名契約、UI 操作契約表、UI 共通動作契約、UI 操作完全性検証契約、UI error / disabled 優先順位固定、[`docs/details/sdk.md`](sdk.md) §23 の SDK 連動実装確認固定契約、[`docs/details/fixture.md`](fixture.md) §27-F を同時に満たす。UI は SDK response に存在しない key を補完せず、状態ファイルを直接読まず、API endpoint を直接呼ばず、成功前に確定表示を行わない。

| 対象 | UI 表示 / 操作 | 使用 SDK method | 成功後再取得 | 固定する確認条件 |
|------|----------------|-----------------|--------------|------------------|
| [`docs/details/runner.md`](runner.md) §27.21 / [`docs/details/runner.md`](runner.md) §27.31 branch target / env | リポジトリ情報 panel に target files と branch env を表示 / 保存する。secret env value は入力欄以外へ表示しない。 | `getBranchConfig()`, `setBranchConfig(branches)`, `getConfig()` | `getBranchConfig()`, `getConfigLog()` | API 配列順を保持し、env key / target path を UI が正規化しない。保存失敗時は secret を消去し、その他入力値を保持する。 |
| [`docs/details/runner.md`](runner.md) §27.22 pipeline | 設定 panel の pipeline config を表示 / 保存する。reserved arg や inline YAML を UI が削除・整形しない。 | `getPipelineConfig()`, `setPipelineConfig(config)` | `getPipelineConfig()`, `getConfigLog()` | `422 details` は該当 field error、成功前に画面上の確定 config を更新しない。 |
| [`docs/details/runner.md`](runner.md) §27.23〜§27.26 local watch / tag / cache / parallel | 設定、履歴、status、build result 表示に API response の watch / tag / cache / target result を表示する。 | `getConfig()`, `setConfig(config)`, `getStatus()`, `getHistory()`, `getHistoryLog(id)` | 操作ごとの表に従う。 | UI は変更検出、tag match、cache hit、parallel result を再計算しない。API response の順序と status を基準とする。 |
| [`docs/details/runner.md`](runner.md) §27.27 hook | フック panel で `command_args` を 1 行 1 引数として表示 / 保存する。 | `getHooks()`, `addHook()`, `deleteHook(id)`, `getHookLog(id)` | `getHooks()`, `getConfigLog()` | 空行だけ除外し、shell 文字列化、quote 展開、環境変数展開を行わない。失敗時は command 入力を保持する。 |
| [`docs/details/api.md`](api.md) §27.30 / [`docs/details/runner.md`](runner.md) §27.30 approval | 承認待ち panel に approval record を API 順で表示し、pending だけ approve / reject を有効にする。 | `getApprovals()`, `approveBuild(id)`, `rejectBuild(id)`, `getQueue()` | `getApprovals()`, `getQueue()` | UI 時刻だけで expired 判定を確定しない。`409` 後は一覧再取得だけ行い、同じ approve / reject を再送しない。 |
| [`docs/details/runner.md`](runner.md) §27.32 notification | 通知設定 panel に channel / notify log / SMTP / webhook を表示 / 保存する。secret は入力欄と mask 表示だけに限定する。 | notify / SMTP / webhook methods | 対象 GET と `getNotifyLog()` / `getConfigLog()` | secret 保存成功・失敗の両方で secret field を消去し、error に secret 平文を表示しない。test は設定を自動保存しない。 |
| [`docs/details/runner.md`](runner.md) §27.33 / [`docs/details/runner.md`](runner.md) §27.38 trend / anomaly | 統計 panel と dashboard alert に trend summary、sample、anomaly を表示する。 | `getBuildTrends()`, `getStatsBuildDuration()`, `getDashboard()`, `getConfig()`, `setConfig(config)` | config 保存時は `getConfig()`, `getConfigLog()`。表示取得は再取得なし。 | avg / median / p95 / anomaly tag を UI が再計算しない。API warnings は panel 内 warning として表示する。 |
| [`docs/details/runner.md`](runner.md) §27.34 / [`docs/details/runner.md`](runner.md) §27.35 chain / queue | chain 設定、queue 表示、queue clear、manual build priority 表示を扱う。 | `getBuildChainConfig()`, `setBuildChainConfig(chains)`, `getQueue()`, `clearQueue()`, `triggerBuild()`, `buildForce()` | chain 保存は `getBuildChainConfig()`, `getConfigLog()`。queue clear は `getQueue()`, `getStatus()`。 | priority / created_seq の並びを UI が変更しない。queue full `429` は同一操作だけ 10 秒 disabled。 |
| [`docs/details/runner.md`](runner.md) §27.36 / [`docs/details/runner.md`](runner.md) §27.37 failure category / environment | 履歴 panel、履歴 detail、システム情報に failure category、evidence、environment を表示する。 | `getHistory()`, `getHistoryLog(id)`, `getOutputMeta()`, `getStatus()` | なし | category ラベル変換、environment fallback 補完、evidence secret 表示を行わない。unknown category warning は warning として表示する。 |
| [`docs/details/security.md`](security.md) §27.42 / [`docs/details/security.md`](security.md) §27.43 token scope / token | API token 管理 panel で scope 複数選択、発行 token 一回表示、失効を扱う。 | `getTokens()`, `createToken()`, `revokeToken(id)`, `getAuditLog()` | token 操作は `getTokens()`, `getAuditLog()` | token 本体は `issued-token-once` に 1 回だけ表示し、一覧へ合成しない。`403` は権限不足表示で logout しない。 |
| [`docs/details/security.md`](security.md) §27.44 audit | 監査ログ panel に actor / action / result filter と結果を表示する。 | `getAuditLog({limit,offset,actor,action,result})` | なし | secret、request body、Authorization header、token hash を表示しない。壊れた行の内容を UI に表示しない。 |
| [`docs/details/security.md`](security.md) §27.45 session timeout / sessions | セキュリティ panel とセッション管理 panel で session timeout、session list、revoke all を扱う。 | `getConfig()`, `setConfig({session_timeout_seconds})`, `getSessions()`, `revokeAllSessions()` | timeout 保存は `getConfig()`, `getConfigLog()`。revoke all は `getSessions()`。 | timeout 更新後も UI が既存 session の期限を再計算しない。revoke all 後は secret field を消去する。 |
| [`docs/details/security.md`](security.md) §27.46 TOTP | セキュリティ / login panel で setup、confirm、disable、login TOTP を扱う。 | `getTotpStatus()`, `setupTotp()`, `confirmTotp(code)`, `disableTotp(code)`, `loginTotp(ticket,code)` | confirm / disable は `getTotpStatus()`, `getAuditLog()` | secret と otpauth URI は一回表示だけ。ticket は DOM に表示しない。code 成功・失敗・panel 遷移・`401` で消去する。 |
| [`docs/details/security.md`](security.md) §27.47 rate limit | セキュリティ panel に policy と state summary を表示 / 保存する。 | `getApiRateLimit()`, `setApiRateLimit(policy)` | `getApiRateLimit()`, `getAuditLog()` | UI は reset_at、count、group を API 値で表示し、window / count を再計算しない。`429` は自動 retry しない。 |

**§27.21〜§27.47 UI 合格ゲート：**

| ゲート | 合格条件 |
|--------|----------|
| SDK only | [`docs/details/ui.md`](ui.md) §24 の固定表の全操作が `AdlaireCI` public method だけを呼び、直接 `fetch()` / `XMLHttpRequest` / `EventSource` / 状態ファイル操作を行わない。 |
| refresh order | 成功後再取得は表の左から順に await し、途中失敗時は変更成功を維持したまま再取得失敗だけを panel error に表示する。 |
| no speculative state | UI が status、queue、approval、token、rate limit、trend、failure category、environment、TOTP 状態を API response なしに確定しない。 |
| secret clearing | password、PAT、Webhook secret、SMTP password、発行 token、TOTP secret、ticket、TOTP code は成功、失敗、panel 遷移、logout、`401`、revoke all で消去される。 |
| error discipline | `401` は login へ戻す。`403` は logout しない。`409` は仕様上の再取得だけ行う。`422` は field error。`429` は同一操作だけ 10 秒 disabled。 |
| one-time display | 発行 token、TOTP secret、otpauth URI は専用領域に 1 回だけ表示し、次 user action、copy、panel 遷移、logout、`401` で消去する。 |
| fixture evidence | [`docs/details/fixture.md`](fixture.md) §27-F の UI 関連 fixture で、SDK only、refresh order、disabled priority、secret clearing、one-time display、no speculative state が確認される。 |

**§27.21〜§27.47 UI 連動 fixture 必須証跡：**

UI 実装変更は、対象 [`docs/details/ui.md`](ui.md) §27 機能ごとに [`docs/details/ui.md`](ui.md) §24 の固定表の証跡を fixture で固定する。UI は API / SDK の返却値を表示する補助層であり、状態確定、補完、保存、再試行を独自判断で行わない。

| UI 証跡 | 固定する内容 | 合格条件 | 禁止条件 |
|---------|--------------|----------|----------|
| SDK only call trace | user action ごとの SDK method 名、引数、呼び出し順。 | [`docs/details/ui.md`](ui.md) §24 の固定表の使用 SDK method だけを呼ぶ。直接 `fetch()`、`XMLHttpRequest`、`EventSource`、状態ファイル操作が 0 件。 | API endpoint を UI から直接呼ぶ、SDK にない method を仮実装する。 |
| refresh order | 成功後再取得、`409` / `429` / `500` 後の再取得、再取得失敗時表示。 | 表の左から順に await し、途中失敗時は変更成功を維持して panel error に固定文言を表示する。 | 再取得失敗を理由に同じ変更 API を再送する。 |
| disabled priority | maintenance、SSE 接続中、送信中、`429`、validation error の優先順位。 | [`docs/details/ui.md`](ui.md) §24 の UI error / disabled 優先順位固定に従い、上位条件が残る限り下位解除で有効化しない。 | `429` timer 終了で maintenance disabled を無視して button を有効化する。 |
| one-time / secret clearing | password、PAT、Webhook secret、SMTP password、発行 token、TOTP secret、otpauth URI、ticket、TOTP code。 | 成功、失敗、panel 遷移、logout、`401`、revoke all、次 user action で対象値が DOM から消える。 | token / secret を一覧、hidden field、data attribute、error message、clipboard 履歴表示へ残す。 |
| no speculative display | status、queue、approval、token、rate limit、trend、failure category、environment、TOTP 状態。 | API / SDK response に存在する値だけを表示し、未知値は panel error または空状態で表現する。 | UI 時刻だけで expired を確定、avg / p95 / anomaly / rate limit count を再計算する。 |
| field error mapping | `422 details` の `field` と panel error summary。 | 該当 field が存在する場合は field error と panel summary、存在しない場合は panel error へ表示する。入力値は保持し secret だけ消去する。 | `422` 後に対象 GET を呼んで入力値を上書きする。 |

**UI 設定値契約：**

| 設定値 | 取得元 | 既定値 | 仕様 |
|--------|--------|--------|------|
| SDK `baseUrl` | `index.html` 内の `data-api-base-url` 属性 | `/api` | 空文字の場合は `/api` を使用する。外部 origin の URL は ui 詳細本文責務では使用しない。 |
| 初期表示 panel | 固定値 | `panel-login` | token 永続化を行わないため、画面読み込み直後は常にログイン panel を表示する。 |
| theme token | `:root` CSS custom property | [`docs/DESIGN.md`](../DESIGN.md) デザイン責務 §2 の値 | JavaScript は theme token を変更しない。UI 操作で theme 切替を実装しない。 |
| panel 表示制御 | `hidden` 属性 | 全 panel hidden、`panel-login` のみ表示 | DOM 削除ではなく `hidden` で切り替える。 |
| API 呼び出し経路 | `AdlaireCI` instance | 1 instance | panel ごとに SDK instance を作らず、画面全体で 1 つの `AdlaireCI` instance を共有する。 |

---
