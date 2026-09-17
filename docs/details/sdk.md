# Adlaire CI — SDK 詳細仕様

本ファイルは `docs/DETAIL_INDEX.md` から分割した `sdk` owner component の詳細仕様である。

本ファイルに、方針、ポリシー、正本関係、実装状態、ロードマップ状態、実装可否の上位判断を記載してはならない。方針、ポリシー、正本関係は `docs/SPEC.md`、実装状態、ロードマップ状態、実装可否は `docs/ROADMAP.md` を正とする。

本ファイルを読む前に、`docs/SPEC.md` で方針とポリシーを確認し、`docs/ROADMAP.md` で実装状態と実装可否を確認し、`docs/DETAIL_INDEX.md` §0〜§0j で共通固定値、責務 component、詳細節対応表、リポジトリ内ソース配置を確認する。本ファイルは `sdk` owner component の主本文であり、collaborator component の仕様は endpoint、response、error、security、UI 呼び出し境界、fixture、検証観点として参照する。

SDK が呼び出す API endpoint の method、path、request、response、error、認証要否は `docs/details/api.md` §22.0e を正とする。本ファイルは SDK 側の class、method、引数変換、transport、error、stream、token 破棄を定義する。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `sdk` |
| collaborator component | `api`、`ui`、`security` |
| 持つ内容 | `sdk` owner が主本文として定義する SDK class、method、HTTP 対応、query / body 生成、error、stream、token 破棄。 |
| 持たない内容 | API endpoint 実装、API endpoint の状態ファイル更新責務、UI DOM 詳細、状態 schema、状態ファイル直接操作、admin 静的配信、setup / release 手順、fixture / PR 証跡正本。 |

---

## 23. JavaScript SDK 仕様

本節は、`sdk` に関する仕様である。

**ファイル：** `admin/adlaire-ci-sdk.js`（単一ファイル、外部依存なし）
**モジュール形式：** ES Module（`import` / `export`）

**SDK 実行環境契約：**

| 項目 | 仕様 |
|------|------|
| JavaScript | ECMAScript 2022 以上を前提とする。transpile、bundle、polyfill は標準仕様に含めない。 |
| module | `admin/adlaire-ci-sdk.js` は ES Module とし、`export { AdlaireCI, AdlaireCIError }` を必須 export とする。default export は定義しない。 |
| browser API | `fetch`、`AbortController`、`ReadableStream.getReader()`、`TextDecoder`、`URLSearchParams` が存在する browser を必須環境とする。いずれかが存在しない場合、`AdlaireCI` constructor は `TypeError("Unsupported browser runtime")` を投げる。 |
| 非 browser runtime | browser API 行の必須 API が存在しない実行環境では、runtime 名を判定分岐せず、`AdlaireCI` constructor が `TypeError("Unsupported browser runtime")` を投げる。Node.js 専用 API、npm package、bundler、polyfill による補完は行わない。 |
| 外部依存 | npm package、CDN script、framework、build tool を使用してはならない。 |
| global 汚染 | `window.AdlaireCI` 等の global 代入を行わない。標準管理ツールは ES Module import で SDK を読み込む。 |
| stream 前提 | `streamBuild()` は native `EventSource` を使用しない。Authorization header を付与できる `fetch` streaming を必須実装とする。 |

```js
class AdlaireCI {
  constructor({ baseUrl })
  // this._token でセッショントークンを管理。login() 後の全リクエストに自動付与

  login(password)                               // POST /api/login → {token?, must_change, totp_required?, ticket?}; token がある場合は this._token にセット
  loginTotp(ticket, code)                       // POST /api/login/totp → {token, must_change}; this._token にセット
  logout()                                      // POST /api/logout; this._token をクリア
  changePassword(currentPassword, newPassword)  // POST /api/change-password
  getAuditLog({ limit = 100, offset = 0, actor = null, action = null, result = null } = {}) // GET /api/audit-log → Promise<{log: AuditLogRecord[], total: number}>

  getStatus()               // GET /api/status               → Promise<StatusObject>
  triggerBuild()            // POST /api/build               → Promise<{message: string, build_id?: string, queued?: boolean}>
  getLogs(n = 100, q = '')  // GET /api/logs?n={n}&q={q}     → Promise<{lines: string[]}>
  getHistory({ page = 1, perPage = 20, trigger = null, tag = null, flagged = null } = {}) // GET /api/history?page={page}&per_page={perPage}&trigger={trigger}&tag={tag}&flagged={flagged} → Promise<HistoryPageObject>
  getSysinfo()              // GET /api/sysinfo              → Promise<SysinfoObject>
  getSchedule()             // GET /api/schedule             → Promise<ScheduleObject>
  getNotifyConfig()         // GET /api/notify-config        → Promise<NotifyConfig>
  setNotifyConfig(config)   // POST /api/notify-config       → Promise<{message: string}>
  getConfig()               // GET /api/config               → Promise<ConfigObject>
  validateConfig(config)    // POST /api/config/validate     → Promise<ConfigValidationObject>
  setConfig(config)         // POST /api/config              → Promise<{message: string, config: ConfigObject}>
  health()                  // GET /api/health               → Promise<{status: string}>
  getPatStatus()            // GET /api/pat-status           → Promise<PatStatusObject>
  getAccessLog()            // GET /api/access-log           → Promise<{log: AccessRecord[]}>
  getApiAccessLog({ limit = 100, offset = 0, method = null, path = null, status = null } = {}) // GET /api/api-access-log → Promise<{log: ApiAccessRecord[], total: number}>
  getStats(days = 7)        // GET /api/stats?days={days}    → Promise<StatsObject>
  exportLogs()              // GET /api/logs/export          → Promise<{exported_at: string, lines: string[]}>
  cleanupLogs()             // POST /api/logs/cleanup        → Promise<{message: string, deleted_count: number}>
  archiveLogs()             // POST /api/logs/archive        → Promise<{message: string, archived_count: number}>
  getRepoInfo()             // GET /api/repo-info            → Promise<RepoInfoObject>
  backup()                  // GET /api/backup               → Promise<BackupObject>
  restore(config)           // POST /api/restore             → Promise<{message: string}>
  notifyTest()              // POST /api/notify-test         → Promise<{message: string, webhook_url: string}>
  buildForce()              // POST /api/build/force         → Promise<{message: string, build_id?: string, queued?: boolean}>
  patVerify()               // POST /api/pat-verify          → Promise<PatVerifyObject>
  getHistoryLog(id)         // GET /api/history/{id}/log     → Promise<HistoryLogObject>
  cancelBuild()             // POST /api/build/cancel        → Promise<{message: string}>
  resetCircuitBreaker()     // POST /api/circuit-breaker/reset → Promise<{message: string, open: boolean, consecutive_failures: number}>
  streamBuild(onLine, onEnd) // GET /api/build/stream (SSE)  → Promise<StreamHandle>（onLine(line), onEnd({status, duration_seconds}) コールバック）
  setLogLevel(level)        // POST /api/log-level           → Promise<{message: string, level: string}>
  updatePat(token)          // POST /api/pat-update          → Promise<{message: string}>
  getDashboard()            // GET /api/dashboard            → Promise<DashboardObject>
  getNotifyLog()            // GET /api/notify-log           → Promise<{log: NotifyRecord[]}>
  getSessions()             // GET /api/sessions             → Promise<{sessions: SessionRecord[]}>
  revokeAllSessions()       // POST /api/sessions/revoke-all → Promise<{message: string, revoked_count: number}>
  getTotpStatus()           // GET /api/auth/totp-status     → Promise<TotpStatus>
  setupTotp()               // POST /api/auth/totp-setup     → Promise<{secret: string, otpauth_uri: string}>
  confirmTotp(code)         // POST /api/auth/totp-confirm   → Promise<TotpStatus>
  disableTotp(code)         // DELETE /api/auth/totp         → Promise<TotpStatus>
  setScheduleInterval(seconds) // POST /api/schedule/interval → Promise<{message: string, interval_seconds: number}>
  pauseSchedule()              // POST /api/schedule/pause   → Promise<{message: string}>
  resumeSchedule()             // POST /api/schedule/resume  → Promise<{message: string}>
  setAllowedHours(from, to)    // POST /api/schedule/allowed-hours {from, to} → Promise<{message: string, allowed_hours: {from: number, to: number}}>
  clearAllowedHours()          // POST /api/schedule/allowed-hours {from:null, to:null} → Promise<{message: string, allowed_hours: null}>
  setForceInterval(hours)      // POST /api/schedule/force-interval → Promise<{message: string, hours: number}>
  setBuildCooldown(seconds)    // POST /api/schedule/cooldown → Promise<{message: string, seconds: number}>
  searchLogs(q = '', from = '', to = '') // GET /api/logs/search?q={q}&from={from}&to={to} → Promise<SearchResult>
  getOutputMeta()              // GET /api/output-meta       → Promise<OutputMetaObject>
  getStatsTimeline(days = 30)  // GET /api/stats/timeline?days={days} → Promise<TimelineObject>
  getStatsBuildDuration(n = 10) // GET /api/stats/build-duration?n={n} → Promise<BuildDurationStats>
  getBuildTrends(n = 100)       // GET /api/stats/build-trends?n={n} → Promise<BuildTrendStats>
  getDiagnostics()             // GET /api/diagnostics       → Promise<DiagnosticsObject>
  getRateLimit()               // GET /api/rate-limit        → Promise<RateLimitObject>
  getDiskUsage()               // GET /api/disk-usage        → Promise<DiskUsageObject>
  getConfigLog()               // GET /api/config-log        → Promise<{log: ConfigLogRecord[]}>
  getApiRateLimit()            // GET /api/api-rate-limit    → Promise<ApiRateLimitPolicy>
  setApiRateLimit(policy)      // POST /api/api-rate-limit   → Promise<ApiRateLimitPolicy>
  getBranchConfig()            // GET /api/branch-config     → Promise<{source: string, branches: BranchTargetRecord[]}>
  setBranchConfig(branches)    // POST /api/branch-config    → Promise<{message: string, branches_count: number}>
  notifyWeeklySummary()        // POST /api/notify/weekly-summary → Promise<{message: string, period: string, success_count: number, failure_count: number, success_rate: number}>
  getWebhookEvents(limit = 50, offset = 0) // GET /api/webhook-events?limit={limit}&offset={offset} → Promise<{events: WebhookEventRecord[], total: number}>
  getWebhookConfig()          // GET /api/webhook-config    → Promise<{configured: boolean}>
  setWebhookConfig(secret)    // POST /api/webhook-config   → Promise<{message: string}>
  getBuildChainConfig()       // GET /api/build-chain-config → Promise<BuildChainConfig>
  setBuildChainConfig(chains) // POST /api/build-chain-config → Promise<{message: string, chains_count: number}>
  getApprovals()              // GET /api/approvals         → Promise<{approvals: ApprovalRecord[]}>
  approveBuild(id)            // POST /api/approvals/{id}/approve → Promise<{message: string, queued: boolean}>
  rejectBuild(id)             // POST /api/approvals/{id}/reject  → Promise<{message: string}>
  getHistoryComment(id)        // GET /api/history/{id}/comment  → Promise<CommentObject>
  setHistoryComment(id, comment) // POST /api/history/{id}/comment → Promise<{message: string}>
  setRepoConfig(config)        // POST /api/repo-config      → Promise<{message: string}>
  exportHistory()              // GET /api/history/export    → Promise<ExportObject>（JSON）
  setHistoryFlag(id, flagged)  // POST /api/history/{id}/flag → Promise<{message: string}>
  setHistoryTags(id, tags)     // POST /api/history/{id}/tags → Promise<{message: string}>
  rollbackHistory(id)          // POST /api/history/{id}/rollback → Promise<{message: string, build_id: string}>
  getTokens()                  // GET /api/tokens            → Promise<{tokens: TokenRecord[]}>
  createToken(label, scopes = ['read'], expiresAt = null) // POST /api/tokens → Promise<TokenCreateResult>
  revokeToken(id)              // DELETE /api/tokens/{id}    → Promise<{message: string}>
  // 14A スナップショット
  getSnapshots()               // GET /api/snapshots         → Promise<{snapshots: SnapshotRecord[]}>
  downloadSnapshot(id)         // GET /api/snapshots/{id}/download → Promise<Blob>
  deleteSnapshot(id)           // DELETE /api/snapshots/{id} → Promise<{message: string}>
  // 14B メンテナンスモード
  getMaintenance()             // GET /api/maintenance       → Promise<MaintenanceObject>
  enableMaintenance(reason)    // POST /api/maintenance/enable → Promise<{message: string, since: string}>
  disableMaintenance()         // POST /api/maintenance/disable → Promise<{message: string}>
  // 14C IP アクセス制限
  getAccessControl()           // GET /api/access-control    → Promise<{allow: string[]}>
  setAccessControl(allowList)  // POST /api/access-control   → Promise<{message: string, allow: string[]}>
  // 14E フック
  getHooks()                   // GET /api/hooks             → Promise<{hooks: HookRecord[]}>
  addHook(phase, commandArgs, abortOnFailure = true, timeoutSeconds = 300) // POST /api/hooks → Promise<HookRecord>
  deleteHook(id)               // DELETE /api/hooks/{id}     → Promise<{message: string}>
  getHookLog(id)               // GET /api/hooks/{id}/log    → Promise<{id: string, runs: HookRunRecord[]}>
  // 15A アラートルール
  getAlertRules()              // GET /api/alert-rules       → Promise<{rules: AlertRule[]}>
  addAlertRule(metric, operator, threshold, level, message) // POST /api/alert-rules → Promise<AlertRule>
  deleteAlertRule(id)          // DELETE /api/alert-rules/{id} → Promise<{message: string}>
  // 15B 自動タグ付けルール
  getTagRules()                // GET /api/tag-rules         → Promise<{rules: TagRule[]}>
  addTagRule(condition, tags)  // POST /api/tag-rules        → Promise<TagRule>
  deleteTagRule(id)            // DELETE /api/tag-rules/{id} → Promise<{message: string}>
  // 15C チェックサム
  verifyOutput()               // POST /api/verify-output    → Promise<{match: boolean, expected: string, actual: string}>
  // 15D パイプライン設定
  getPipelineConfig()          // GET /api/pipeline-config   → Promise<PipelineConfig>
  setPipelineConfig(config)    // POST /api/pipeline-config  → Promise<{message: string}>
  // 15E 運用ノート
  getNotes()                   // GET /api/notes             → Promise<{content: string, updated_at: string|null}>
  setNotes(content)            // POST /api/notes            → Promise<{message: string, updated_at: string|null}>
  // 16B メール通知
  getSmtpConfig()              // GET /api/smtp-config       → Promise<SmtpConfig>
  setSmtpConfig(config)        // POST /api/smtp-config      → Promise<{message: string}>
  smtpTest()                   // POST /api/smtp-test        → Promise<{result: string, message: string}>
  // 16C ビルドキュー
  getQueue()                   // GET /api/queue             → Promise<{queued: QueueEntry[], max_size: number}>
  clearQueue()                 // DELETE /api/queue          → Promise<{message: string, cleared_count: number}>
  // 16D ダッシュボードレイアウト
  getDashboardLayout()         // GET /api/dashboard-layout  → Promise<{widgets: string[]}>
  setDashboardLayout(widgets)  // POST /api/dashboard-layout → Promise<{message: string}>
}

export { AdlaireCI, AdlaireCIError };
```

全メソッドは `Promise` を返す。`streamBuild` は SSE 接続確立後に `StreamHandle` で resolve し、接続前エラーは `AdlaireCIError` で reject する。HTTP エラー（4xx / 5xx）は `AdlaireCIError` としてスローする。`401` 受信時はセッション期限切れとして `this._token` をクリアする。constructor、private method、helper 関数を除く public method は §22.0e の SDK 列と完全一致させる。

**SDK 共通実装契約：**

| 項目 | 仕様 |
|------|------|
| `baseUrl` | 末尾 `/` を除去して保持する。空文字、`null`、`undefined` は `TypeError`。 |
| URL 組み立て | パスは `/api/...` をそのまま連結し、クエリ値は `encodeURIComponent` でエンコードする。 |
| 認証ヘッダー | `this._token` が存在する場合のみ `Authorization: Bearer ${token}` を付与する。 |
| JSON 送信 | `POST` / `DELETE` で body を送る場合は `Content-Type: application/json` を付与し、`JSON.stringify` した body を送信する。 |
| JSON 受信 | `Content-Type` が JSON の場合のみ `response.json()` を呼ぶ。§22.0e で JSON response を定義した endpoint の成功時に空 body を受信した場合は protocol error として `AdlaireCIError(status=0, message="Empty JSON response")` を投げる。 |
| `AdlaireCIError` | `name="AdlaireCIError"`、`status`、`message`、`details`、`responseBody` を持つ `Error` 派生クラスとする。constructor は `new AdlaireCIError({status, message, details = null, responseBody = null})` とし、network error、timeout、protocol error は `status=0` とする。`message` は API error response の `error`、network error は `"Network error"`、timeout は `"Request timeout"`、protocol error は固定文言を使用する。 |
| `logout()` | API 呼び出しが失敗しても `finally` で `this._token` をクリアする。 |
| request timeout | 通常 API は 30 秒で abort し、`AdlaireCIError(status=0, message="Request timeout")` を投げる。`streamBuild()` は接続確立まで 30 秒、接続確立後は timeout なしとし、利用者が `StreamHandle.close()` で停止する。 |
| `streamBuild()` | token がない場合は接続前に `AdlaireCIError(status=401, message="Unauthorized")` を投げる。native `EventSource` は Authorization header を付与できないため使用禁止とする。SDK は `fetch()`、`AbortController`、`ReadableStream` reader を使用し、`Accept: text/event-stream` と `Authorization` header を付与して SSE frame を解析する。`data:` 行の JSON を parse し、`type="log"` は `onLine(line)`、`type="end"` は `onEnd({status, duration_seconds})` を呼んで reader を close する。parse 不能 frame は `AdlaireCIError(status=0, message="Invalid SSE frame")` として stream error にする。 |
| `StreamHandle` | `streamBuild()` の戻り値は `{ close(): void, closed: boolean }` とする。`close()` は AbortController を abort し、複数回呼んでも例外を投げない。`closed` は `end` 受信、error、または `close()` 後に `true` になる。 |
| Blob レスポンス | `downloadSnapshot(id)` のみ `response.blob()` を使用する。その他は JSON とする。 |
| メソッド引数検証 | SDK 側でも必須引数の空値、配列型、数値範囲を検証し、HTTP 送信前に `TypeError` を投げる。 |
| endpoint 対応 | SDK method は §22.0e の SDK 列に存在する endpoint だけを呼び出す。§22.0e にない endpoint を SDK 独自判断で追加してはならない。 |
| body なし endpoint | §22.0e の `Request` が `none` の場合、SDK は `fetch` に `body` を設定しない。`{}` も送信しない。 |
| token 保存 | セッショントークンはメモリ上の `this._token` のみに保持する。`localStorage`、`sessionStorage`、Cookie へ保存しない。 |
| 秘密情報引数 | `updatePat(token)`、`setWebhookConfig(secret)`、SMTP password、`createToken()` の返却 token は console 出力しない。 |
| query 生成 | `undefined`、`null`、空文字の任意 query は送信しない。ただし `q`、`from`、`to` は endpoint 仕様で空文字を有効値として定義している場合だけ、空文字を query value として送信する。 |
| 戻り値補完禁止 | 成功時は API response にない key を追加しない。失敗時は HTTP status、API error、details、responseBody 以外を推測しない。fallback 値は API response に含まれる値だけを返し、表示用加工は UI 側で行う。 |
| retry | SDK は自動 retry を行わない。ユーザー操作による再実行、または UI の明示的な再取得のみを許可する。 |

**SDK transport / error 固定契約：**

SDK の内部 request helper は、すべての public method で下表の処理順に固定する。public method ごとに個別 fetch 処理を複製してはならない。

| 順序 | 処理 | 固定仕様 |
|------|------|----------|
| 1 | 引数検証 | 必須引数、型、空配列、数値範囲を検証する。失敗時は `TypeError` を投げ、`fetch()` を呼ばない。 |
| 2 | URL 生成 | `baseUrl + path + query` を生成する。query key は method 契約表の順序で追加する。 |
| 3 | body 生成 | `Request=none` では `body` と `Content-Type` を設定しない。JSON body ありの場合だけ `JSON.stringify()` する。 |
| 4 | header 生成 | `Accept` を常に付与する。JSON body を送信する場合だけ `Content-Type` を付与する。token が空でない場合だけ `Authorization` を付与し、token が空の場合は `Authorization` header を付けない。 |
| 5 | timeout 設定 | 通常 request は `AbortController` で 30 秒 timeout。`streamBuild()` は接続確立まで 30 秒。 |
| 6 | `fetch()` 実行 | network error、abort、CORS 等の失敗はすべて `AdlaireCIError(status=0,message="Network error")`、timeout だけ `"Request timeout"` とする。 |
| 7 | response parse | 成功 / 失敗に関わらず JSON error body がある場合は parse する。parse 不能 error body は `responseBody` に text を保持し、`message` は HTTP status 固定文言とする。 |
| 8 | token 変化 | `401` の場合だけ `this._token = null`。`403`、`429`、`500`、network error、timeout では token を破棄しない。 |
| 9 | return / throw | `2xx` は endpoint の型で返す。`4xx` / `5xx` は `AdlaireCIError` を投げる。 |

HTTP status と SDK error の対応は下表に固定する。

| 条件 | `AdlaireCIError.status` | `message` | `details` | token |
|------|-------------------------|-----------|-----------|-------|
| API JSON error | HTTP status | response の `error` | response の `details` が配列なら保持 | `401` のみ破棄 |
| API error body なし | HTTP status | `HTTP {status}` | `null` | `401` のみ破棄 |
| API error JSON parse 不能 | HTTP status | `HTTP {status}` | `null` | `401` のみ破棄 |
| network error | `0` | `Network error` | `null` | 維持 |
| timeout | `0` | `Request timeout` | `null` | 維持 |
| empty success JSON | `0` | `Empty JSON response` | `null` | 維持 |
| invalid success JSON | `0` | `Invalid JSON response` | `null` | 維持 |
| invalid SSE frame | `0` | `Invalid SSE frame` | `null` | 維持 |

`429` は SDK で自動待機、自動再送、自動 refresh を行わない。binary response は `downloadSnapshot(id)` の `2xx` のみ `Blob` とする。`4xx` / `5xx` では `Content-Type` が `application/json` または `+json` で終わる場合だけ JSON error として parse し、parse 成功時は response の `error` / `details` を保持した `AdlaireCIError` を投げる。JSON parse 不能、JSON 以外の error body、空 body の場合は body text を `responseBody` に保持し、`message` は `HTTP {status}` とする。`streamBuild()` は接続後の `close()` をユーザー停止として扱い、`AdlaireCIError` を投げない。接続後に network error または invalid frame が発生した場合は `StreamHandle.closed=true`、`StreamHandle.error` に `AdlaireCIError(status=0,message="Network error")` または `AdlaireCIError(status=0,message="Invalid SSE frame")` を保存する。

**SDK メソッド実装固定契約：**

| 項目 | 仕様 |
|------|------|
| public method 定義順 | class 内の public method は §23 の一覧順に定義する。追加 public method を末尾に置くことは禁止し、先に §22.0e と本一覧を更新する。 |
| private helper | private helper は `_request`, `_json`, `_query`, `_requireToken`, `_validateId`, `_clearTokenOn401` だけを定義する。helper を export しない。 |
| TypeError 文言 | SDK 側引数検証の `TypeError.message` は `"Invalid argument: <name>"` に固定する。複数不正がある場合は最初に検出した引数だけを返す。 |
| path parameter | `id` を path に入れる method は、SDK 側で `encodeURIComponent(id)` を必ず行う。`/`、`.`、`..`、空文字は送信前に `TypeError`。 |
| query parameter | query key は §23 SDK 引数変換契約の表記順で生成する。任意 query が未指定の場合、`?` 自体を付けない。 |
| body parameter | body object の key 順は §23 SDK 引数変換契約の送信値順とする。未知 key を SDK が追加しない。 |
| token mutation | `login()` と `loginTotp()` は response に `token` が存在する場合だけ `this._token` を更新する。`totp_required:true` かつ token なしの場合は既存 token を保持せず `null` にする。 |
| logout failure | `logout()` は network error、`401`、`500` のいずれでも `finally` で `this._token=null` にする。 |
| response passthrough | 成功 response は clone、整形、既定値 merge を行わず、そのまま返す。Blob と StreamHandle は例外とする。 |
| error details | API error response の `details` が配列なら `AdlaireCIError.details` に同じ配列を保持する。配列でなければ `null`。 |
| responseBody | JSON parse できた error body は object のまま、parse 不能 error body は先頭 4000 文字の string として `responseBody` に保持する。 |

**SDK 検証 fixture：**

| fixture | 入力 | 合格条件 |
|---------|------|----------|
| auth token flow | `login()`、`loginTotp()`、`logout()`、`401` response | token set / clear が仕様どおり。localStorage、sessionStorage、Cookie を使わない。 |
| request shape | 全 public method を fake fetch で呼ぶ | method、path、query、body、headers が §22.0e と §23 引数変換契約に一致する。 |
| error shape | `400`、`401`、`403`、`422 details`、`500`、network error、timeout | `AdlaireCIError` の `status`、`message`、`details`、`responseBody` が固定値になる。 |
| stream | log frame、end frame、invalid frame、client close | callback、closed、error が仕様どおり。EventSource を使用しない。 |
| binary | `downloadSnapshot(id)` | `Blob` を返し、JSON parse を試みない。 |

**Phase 3 SDK 操作固定契約：**

Phase 3 実装では、下表の SDK method を最小運用範囲として固定する。SDK は成功時 response を endpoint schema の範囲でそのまま返し、失敗時は HTTP status、API error、details、responseBody を保持した `AdlaireCIError` へ変換する。UI が必要とする表示用既定値、並べ替え、ラベル変換は SDK で行わない。

| SDK method | HTTP | 成功時 | 失敗時 | 追加禁止事項 |
|------------|------|--------|--------|--------------|
| `login(password)` | `POST /api/login` | `token` がある場合だけ `this._token` へ保存する。`totp_required:true` の場合は token を保存せず response を返す。 | `401`、`429`、`500` は `AdlaireCIError`。`401` で既存 token を破棄する。 | password を console、error、responseBody 加工結果へ出さない。 |
| `logout()` | `POST /api/logout` | response に関わらず `finally` で token を破棄する。 | network error、`401`、`500` でも token 破棄後に error を投げる。 | logout 失敗を理由に token を保持しない。 |
| `getStatus()` | `GET /api/status` | `StatusObject` をそのまま返す。 | `500 State file is corrupted` / `State file read failed` を message として保持する。 | `.build_status.json` 不在時の fallback 値を SDK が推測しない。 |
| `triggerBuild()` | `POST /api/build` | `{message, build_id?, queued?}` を返す。 | `409`、`429`、`503` は `AdlaireCIError.status` に HTTP status を保持し、`message` は API response の `error` を保持する。 | running / queue / circuit を SDK 側で事前判定しない。 |
| `buildForce()` | `POST /api/build/force` | `{message, build_id?, queued?}` を返す。 | `409`、`429`、`503` を `AdlaireCIError`。 | force 可否を SDK 側で状態推測しない。 |
| `cancelBuild()` | `POST /api/build/cancel` | `{message}` を返す。 | running なしの `409` を `AdlaireCIError(status=409)`。 | cancel 後に SDK が自動 `getStatus()` を呼ばない。 |
| `getLogs(n,q)` | `GET /api/logs` | `{lines}` を返す。`lines` は API 順序を保持する。 | `500` は固定 error message を保持する。 | line を結合、trim、level 分類しない。 |
| `getHistory(options)` | `GET /api/history` | `{total,page,per_page,pages,history}` を返す。 | query 不正 `422` は `details` を保持する。 | `pages`、`total` を SDK 側で再計算しない。 |
| `getHistoryLog(id)` | `GET /api/history/{id}/log` | log object を返す。 | `404`、`500` を `AdlaireCIError`。 | archive fallback を SDK 側で再試行しない。 |
| `getQueue()` | `GET /api/queue` | `{queued,max_size}` を返す。`running` が response に含まれる場合も削除しない。 | `.build_state` 破損の `500` を固定 message で保持する。 | queue 並び替え、重複排除、上限補正をしない。 |
| `resetCircuitBreaker()` | `POST /api/circuit-breaker/reset` | `{message,open,consecutive_failures}` を返す。 | 破損状態 `500` を `AdlaireCIError`。 | reset 成功後に SDK が自動 build を開始しない。 |
| `streamBuild(onLine,onEnd)` | `GET /api/build/stream` | `StreamHandle` を返し、`log` frame を `onLine`、`end` frame を `onEnd` へ渡す。 | 接続前 `401`、`404`、timeout、invalid frame を `AdlaireCIError`。 | `EventSource`、自動 reconnect、log 永続化を行わない。 |

**Phase 3 SDK fixture 固定：**

| fixture | fake fetch 入力 | 合格条件 |
|---------|-----------------|----------|
| sdk phase3 status corrupted | `GET /api/status` が `500 {"error":"State file is corrupted"}` | `AdlaireCIError.status=500`、`message="State file is corrupted"`、token 維持。 |
| sdk phase3 build conflict | `POST /api/build` が `409 {"error":"Conflict"}` | `AdlaireCIError.status=409`、自動 retry なし、自動 `getStatus()` 呼び出しなし。 |
| sdk phase3 queue full | `POST /api/build` が `429 {"error":"queue_full"}` | `AdlaireCIError.status=429`、`message="queue_full"`、body 再送なし。 |
| sdk phase3 history paging | `getHistory({page:2,perPage:20,trigger:"manual"})` | query は `page=2&per_page=20&trigger=manual`、`total` と `pages` は API 値をそのまま返す。 |
| sdk phase3 log not found | `GET /api/history/{id}/log` が `404 {"error":"Not found"}` | `AdlaireCIError.status=404`、`id` は `encodeURIComponent` 済み。 |
| sdk phase3 stream end | `log` frame 2 件、`end` frame 1 件 | `onLine` 2 回、`onEnd` 1 回、`StreamHandle.closed=true`。 |
| sdk phase3 stream invalid | `data:` 行が JSON parse 不能 | `AdlaireCIError(status=0,message="Invalid SSE frame")`、`closed=true`。 |
| sdk phase3 unauthorized | 任意 Phase 3 endpoint が `401` | `this._token=null`、次 request に Authorization header を付けない。 |

**Phase 4 SDK 操作固定契約：**

Phase 4 SDK は、§22.0e の endpoint 契約と §23 SDK 引数変換契約だけに従う。SDK は保存前検証の一部を `TypeError` で行う場合でも、検証対象は必須引数、型、範囲、path parameter 形式に限定する。API response の補完、no-op 判定、secret mask 変換、状態ファイル由来値の再計算を行ってはならない。

| 機能群 | SDK method | 成功時 | 失敗時 | 追加禁止事項 |
|--------|------------|--------|--------|--------------|
| config / repo / branch | `getConfig()`, `setConfig(config)`, `validateConfig(config)`, `setRepoConfig(config)`, `getBranchConfig()`, `setBranchConfig(branches)` | API response をそのまま返す。`set*` は success message と count / config を保持する。 | `422 details` は `AdlaireCIError.details` に保持する。`500` は固定 message。 | SDK 側で未知 key を削除しない。既定値 merge しない。 |
| schedule | `setScheduleInterval()`, `pauseSchedule()`, `resumeSchedule()`, `setAllowedHours()`, `clearAllowedHours()`, `setForceInterval()`, `setBuildCooldown()` | response の interval / hours / seconds / allowed_hours をそのまま返す。 | systemd 更新失敗の `500` を `AdlaireCIError` にする。 | timer 状態を SDK が推測しない。 |
| diagnostics / dashboard | `getDiagnostics()`, `getDashboard()`, `getRateLimit()`, `getDiskUsage()`, `getOutputMeta()` | read-only response をそのまま返す。 | item 単位 warn/error は成功 response として返し、HTTP error だけ例外にする。 | read-only response を SDK が保存・集計しない。 |
| notify / SMTP / webhook | `getNotifyConfig()`, `setNotifyConfig()`, `notifyTest()`, `notifyWeeklySummary()`, `getWebhookEvents()`, `getWebhookConfig()`, `setWebhookConfig()`, `getSmtpConfig()`, `setSmtpConfig()`, `smtpTest()` | secret は API が mask した値だけ返す。 | 未設定 `422` / `501`、送信失敗 `500` を保持する。 | secret 平文を console、throw message、responseBody 加工結果へ出さない。 |
| snapshots / rollback | `getSnapshots()`, `downloadSnapshot(id)`, `deleteSnapshot(id)`, `rollbackHistory(id)` | list は API 順序、download は Blob、rollback は `{message,build_id}`。 | `404`、running `409` を `AdlaireCIError`。 | SDK が snapshot 存在確認や rollback 可否を事前推測しない。 |
| maintenance / access / hooks | `getMaintenance()`, `enableMaintenance()`, `disableMaintenance()`, `getAccessControl()`, `setAccessControl()`, `getHooks()`, `addHook()`, `deleteHook()`, `getHookLog()` | API response をそのまま返す。 | access deny `403`、validation `422`、hook timeout `500` を保持する。 | command args を文字列結合しない。CIDR を SDK 独自正規化しない。 |
| alert / tag / pipeline / notes / dashboard layout | `getAlertRules()`, `addAlertRule()`, `deleteAlertRule()`, `getTagRules()`, `addTagRule()`, `deleteTagRule()`, `getPipelineConfig()`, `setPipelineConfig()`, `getNotes()`, `setNotes()`, `getDashboardLayout()`, `setDashboardLayout()` | API response の rules / config / notes / widgets をそのまま返す。 | duplicate `409`、validation `422`、read failure `500` を保持する。 | rule 重複排除、widget 補完、notes trim を行わない。 |
| tokens / sessions / audit | `getTokens()`, `createToken()`, `revokeToken()`, `getSessions()`, `revokeAllSessions()`, `getAuditLog()`, `getApiAccessLog()` | `createToken()` の token 本体は response として 1 回だけ返す。 | `403`、`404`、`422`、`429` を status 付きで保持する。 | token 本体を保存しない。token list に作成時 token を合成しない。 |

**Phase 4 SDK fixture 固定：**

| fixture | fake fetch 入力 | 合格条件 |
|---------|-----------------|----------|
| sdk phase4 config validation | `POST /api/config` が `422 details` | `AdlaireCIError.status=422`、`details` 配列保持、送信 body の未知 key は削除されていない。 |
| sdk phase4 schedule failure | `POST /api/schedule/interval` が `500 {"error":"Internal server error"}` | error を投げ、SDK が timer 再試行や rollback request を行わない。 |
| sdk phase4 secret mask | `GET /api/notify-config` と `GET /api/smtp-config` が mask 値を返す | mask 値をそのまま返し、secret 平文を生成しない。 |
| sdk phase4 webhook events paging | `getWebhookEvents(20,40)` | query は `limit=20&offset=40`、`total` は API 値をそのまま返す。 |
| sdk phase4 snapshot binary | `downloadSnapshot(id)` が binary response | `Blob` を返し、JSON parse を試みない。 |
| sdk phase4 rollback conflict | `rollbackHistory(id)` が `409 {"error":"Build is running"}` | `AdlaireCIError.status=409`、自動 `getStatus()` 呼び出しなし。 |
| sdk phase4 token issue | `createToken()` が `{token:"..."}` を返す | token を response として返すだけで、SDK 内部保存、console 出力、token list 合成をしない。 |
| sdk phase4 duplicate rule | `addAlertRule()` または `addTagRule()` が `409 Conflict` | `AdlaireCIError.status=409`、自動 retry なし。 |

**SDK 引数変換契約：**

SDK method は、下表の通りに引数を path、query、body へ変換する。下表にない引数、既定値、body key を追加してはならない。

| SDK method | 引数 | 変換先 | 送信値 |
|------------|------|--------|--------|
| `login(password)` | `password` | body | `{password}` |
| `loginTotp(ticket,code)` | `ticket`, `code` | body | `{ticket,code}` |
| `changePassword(currentPassword,newPassword)` | `currentPassword`, `newPassword` | body | `{current_password: currentPassword, new_password: newPassword}` |
| `getAuditLog({limit,offset,actor,action,result})` | `limit=100`, `offset=0`, `actor=null`, `action=null`, `result=null` | query | `limit`、`offset` は常に送信する。任意値は `null` の場合送信しない。 |
| `getLogs(n,q)` | `n=100`, `q=""` | query | `n`、`q`。`q` は空文字でも送信する。 |
| `getHistory({page,perPage,trigger,tag,flagged})` | `page=1`, `perPage=20`, `trigger=null`, `tag=null`, `flagged=null` | query | `page`、`per_page: perPage` は常に送信する。`trigger`、`tag`、`flagged` は `null` の場合は送信しない。 |
| `getApiAccessLog({limit,offset,method,path,status})` | `limit=100`, `offset=0`, `method=null`, `path=null`, `status=null` | query | `limit`、`offset` は常に送信する。`method`、`path`、`status` は `null` の場合は送信しない。 |
| `validateConfig(config)` | `config` | body | `config` をそのまま送信する。保存は API が行わない。 |
| `setNotifyConfig(config)` | `config` | body | `config` をそのまま送信する。 |
| `setConfig(config)` | `config` | body | `config` をそのまま送信する。未知 key は送信前に削除せず、API の `422` に委ねる。 |
| `updatePat(token)` | `token` | body | `{token}` |
| `setScheduleInterval(seconds)` | `seconds` | body | `{interval_seconds: seconds}` |
| `setAllowedHours(from,to)` | `from`, `to` | body | `{from,to}` |
| `clearAllowedHours()` | なし | body | `{from:null,to:null}` |
| `setForceInterval(hours)` | `hours` | body | `{hours}` |
| `setBuildCooldown(seconds)` | `seconds` | body | `{seconds}` |
| `searchLogs(q,from,to,level)` | `q=""`, `from=""`, `to=""`, `level=undefined` | query | `q`、`from`、`to` は空文字でも送信する。`level` は指定時のみ送信する。 |
| `getWebhookEvents(limit,offset)` | `limit=50`, `offset=0` | query | `limit`、`offset` |
| `setWebhookConfig(secret)` | `secret` | body | `{secret}` |
| `confirmTotp(code)` | `code` | body | `{code}` |
| `disableTotp(code)` | `code` | body | `{code}` |
| `setApiRateLimit(policy)` | `policy` | body | `policy` をそのまま送信する。 |
| `setHistoryComment(id,comment)` | `id`, `comment` | path/body | path `{id}`、body `{comment}` |
| `setHistoryFlag(id,flagged)` | `id`, `flagged` | path/body | path `{id}`、body `{flagged}` |
| `setHistoryTags(id,tags)` | `id`, `tags` | path/body | path `{id}`、body `{tags}` |
| `createToken(label,scopes,expiresAt)` | `label`, `scopes=["read"]`, `expiresAt=null` | body | `{label,scopes,expires_at: expiresAt}` |
| `addHook(phase,commandArgs,abortOnFailure,timeoutSeconds)` | `phase`, `commandArgs`, `abortOnFailure=true`, `timeoutSeconds=300` | body | `{phase,command_args: commandArgs, abort_on_failure: abortOnFailure, timeout_seconds: timeoutSeconds}` |
| `addAlertRule(metric,operator,threshold,level,message)` | 各引数 | body | `{metric,operator,threshold,level,message}` |
| `addTagRule(condition,tags)` | `condition`, `tags` | body | `{condition,tags}` |
| `setPipelineConfig(config)` | `config` | body | `config` をそのまま送信する。 |
| `setNotes(content)` | `content` | body | `{content}` |
| `setSmtpConfig(config)` | `config` | body | `config` をそのまま送信する。`password` が未指定の場合は送信しない。 |
| `setDashboardLayout(widgets)` | `widgets` | body | `{widgets}` |

path に入る `id` は `encodeURIComponent` したうえで 1 segment として連結する。`/`、`.`、空文字を含む `id` は HTTP 送信前に `TypeError` とする。

**SDK method 完全性検証契約：**

SDK 実装完了時は、§22.0e の SDK 列に記載された method 名と `AdlaireCI.prototype` の public method 名が一致しなければならない。constructor、private method、helper 関数、`AdlaireCIError` は比較対象外とする。

| 検証項目 | 合格条件 |
|----------|----------|
| endpoint coverage | §22.0e の SDK 列で `none` 以外の method がすべて `AdlaireCI.prototype` に存在する。 |
| extra method | §22.0e の SDK 列に存在しない public method がない。 |
| request shape | 各 method が §22.0e の Request と §23 SDK 引数変換契約どおりの path / query / body を生成する。 |
| response handling | JSON endpoint は JSON object を返し、binary endpoint は `Blob`、stream endpoint は `StreamHandle` を返す。 |
| error handling | 4xx / 5xx、network error、timeout、empty JSON、invalid SSE frame が `AdlaireCIError` になる。 |
| token handling | `login()` 成功で token を保持し、`logout()` と `401` で token を破棄する。 |

**§27.21〜§27.47 SDK 連動実装完了固定契約：**

§27.21〜§27.47 の追加仕様化機能を SDK で実装完了と扱うには、対象 owner component の詳細仕様、`docs/details/api.md` §27 の連動参照表、§23 SDK 引数変換契約、SDK method 完全性検証契約、`docs/details/fixture.md` §27-F を同時に満たす。SDK は API の補助層であり、API response の補完、状態推測、保存済み値の再計算、UI 表示用変換、自動 retry、自動 refresh、状態ファイル直接操作を行ってはならない。

| 対象 | SDK method | request 固定 | success 固定 | error 固定 | 禁止事項 |
|------|------------|--------------|---------------|------------|----------|
| §27.21 branch target / env | `getBranchConfig()`, `setBranchConfig(branches)`, `getConfig()`, `setConfig(config)` | `branches` または `config` を指定 key のまま送信する。`target_files`、`env`、secret 風 key を削除しない。 | API response の branch 配列、`source`、件数をそのまま返す。 | validation `422` details を保持する。 | target path 正規化、env key 並び替え、secret mask 判定を SDK が行わない。 |
| §27.22 / §27.27 pipeline / hook | `getPipelineConfig()`, `setPipelineConfig(config)`, `getHooks()`, `addHook()`, `deleteHook(id)`, `getHookLog(id)` | pipeline config はそのまま送る。hook `commandArgs` は配列のまま送る。 | config / hook / log response をそのまま返す。 | `409`、`422`、`500` を `AdlaireCIError` として保持する。 | shell 文字列結合、quote 展開、reserved arg 除去、env 補完を行わない。 |
| §27.23〜§27.26 local watch / tag / cache / parallel | `getConfig()`, `setConfig(config)`, `getStatus()`, `getHistory()`, `getHistoryLog(id)` | config と query を表どおり送る。 | status / history / log に含まれる追加 key を削除しない。 | `422` details、`500` message を保持する。 | local watch 可否、tag match、cache hit、parallel result を SDK が再判定しない。 |
| §27.30 approval | `getApprovals()`, `approveBuild(id)`, `rejectBuild(id)`, `getQueue()` | `id` は path parameter 契約で encode し、approve / reject に body を送らない。 | approve は `{message,queued}`、reject は `{message}`、list は API 配列順で返す。 | `409` / `429` / `500` を保持し、自動再送しない。 | expired / rejected / approved を SDK が書き換えない。approve 後に SDK が自動 queue 取得しない。 |
| §27.32 notification | `getNotifyConfig()`, `setNotifyConfig(config)`, `getNotifyLog()`, `notifyTest()`, `notifyWeeklySummary()`, `getSmtpConfig()`, `setSmtpConfig(config)`, `smtpTest()` | config は unknown key を削除せず API へ送る。password 未指定時だけ `setSmtpConfig()` は password key を送らない。 | API が返した mask 値、log、pending 状態をそのまま返す。 | 未設定 `422` / `501`、送信失敗 `500` を保持する。 | secret 平文を console、error message、responseBody 加工結果、SDK field に保存しない。 |
| §27.33 / §27.38 trend / anomaly | `getBuildTrends(n)`, `getStatsBuildDuration(n)`, `getDashboard()`, `getConfig()`, `setConfig(config)` | `n` は number として query へ送る。config はそのまま送る。 | summary、samples、warnings、anomaly flag を API 値のまま返す。 | invalid `n` / config `422` を保持する。 | avg / median / p95 / anomaly を SDK が再計算しない。 |
| §27.34 / §27.35 chain / queue | `getBuildChainConfig()`, `setBuildChainConfig(chains)`, `getQueue()`, `clearQueue()`, `triggerBuild()`, `buildForce()` | chains は配列のまま送る。queue clear は body を送らない。 | queue priority / created_seq / chain config を API 順序のまま返す。 | queue full `429`、conflict `409`、validation `422` を保持する。 | priority 並び替え、queue 重複排除、chain DAG 検証を SDK が行わない。 |
| §27.36 / §27.37 failure category / environment | `getHistory()`, `getHistoryLog(id)`, `getOutputMeta()`, `getStatus()` | filter query は指定された値だけ送る。 | `failure_category`、`failure_evidence`、environment object を削除しない。 | 未知 filter `422` を保持する。 | category 分類、environment fallback、secret mask 判定を SDK が行わない。 |
| §27.42〜§27.47 security | `getTokens()`, `createToken()`, `revokeToken()`, `getAuditLog()`, `getSessions()`, `revokeAllSessions()`, `getTotpStatus()`, `setupTotp()`, `confirmTotp(code)`, `disableTotp(code)`, `getApiRateLimit()`, `setApiRateLimit(policy)` | token / TOTP / rate limit request は §23 引数変換契約どおり送る。 | token 本体、TOTP secret、otpauth URI は response として返すだけで SDK 内部に保存しない。 | `401` は token 破棄、`403` は token 維持、`429` は自動待機なし、`500` は message 保持。 | scope 推測、rate limit 待機、token list への token 合成、TOTP code 再送、audit 補完を行わない。 |

**§27.21〜§27.47 SDK 合格ゲート：**

| ゲート | 合格条件 |
|--------|----------|
| method coverage | §27.21〜§27.47 の api / sdk / ui 連動参照表で SDK method が記載された項目について、該当 public method が存在する。 |
| request exactness | fake fetch fixture で method、path、query key 順、body key 順、body なし endpoint が §23 と一致する。 |
| response passthrough | 成功 response に存在する追加 key を削除せず、存在しない key を追加しない。 |
| security handling | token、password、PAT、Webhook secret、SMTP password、TOTP secret、ticket、Authorization header を SDK property、console、error message に保存しない。 |
| error stability | `401` / `403` / `409` / `422` / `429` / `500` / network / timeout が固定 `AdlaireCIError` になり、自動 retry、自動 refresh、自動 logout は仕様に記載された場合だけ行う。 |
| binary / stream | snapshot download は `Blob`、SSE は `StreamHandle` とし、JSON response と混同しない。 |
| fixture evidence | `docs/details/fixture.md` §27-F の SDK / UI 関連 fixture で、request shape、error shape、secret leak、token mutation、no retry、no response補完が確認される。 |

**§27.21〜§27.47 SDK 連動 fixture 必須証跡：**

SDK 実装 PR は、対象 §27 機能ごとに下表の証跡を fixture で固定する。下表の証跡がない場合、SDK method が存在していても実装完了としない。

| 証跡 | 固定する内容 | 合格条件 | 禁止事項 |
|------|--------------|----------|----------|
| request trace | `method`、`path`、query key 順、body key、body なし endpoint。 | §22.0e と §23 SDK 引数変換契約に完全一致する。 | body なし endpoint へ `{}` を送る、query 未指定時に `?` を付ける。 |
| response passthrough | API success body、binary body、SSE frame。 | SDK は存在 key を削除せず、存在しない key を追加しない。binary は `Blob`、SSE は `StreamHandle`。 | UI 用 label、集計値、既定値、token list、rate limit reset を SDK が合成する。 |
| error object | `401`、`403`、`409`、`422 details`、`429`、`500`、network、timeout、protocol error。 | すべて `AdlaireCIError` になり、`status`、`message`、`details`、`responseBody` が固定される。 | `403` で token を破棄する、`409` / `429` を自動 retry する。 |
| token mutation | login / logout / `401` / `403` / token create。 | login 成功だけ `_token` を設定し、logout と `401` だけ破棄する。createToken の token 本体は保存しない。 | `localStorage`、`sessionStorage`、Cookie、console、token list への保存。 |
| secret leak | PAT、Webhook secret、SMTP password、TOTP secret、ticket、Authorization header。 | SDK property、throw message、console、加工済み response に平文が残らない。 | `responseBody` を UI 用に文字列加工して secret を露出する。 |
| no side-effect helper | 自動 refresh、自動 retry、自動 logout、自動 queue fetch。 | 仕様で明記された `logout()` finally と `401` token 破棄以外の副作用を行わない。 | approve 後に `getQueue()` を SDK が自動実行するなど UI 責務を代行する。 |

**SDK 型定義表：**

本表は SDK が返す object 型の正本である。`nullable` は `null` を許可することを示す。配列は API response に `[]` として存在する場合だけ `[]` を返し、SDK が未取得配列を生成してはならない。API response に存在しないキーを SDK が補完してはならない。ただし `GET /api/config`、`GET /api/notify-config`、`GET /api/dashboard-layout` の既定値 merge は API 側の責務とする。

| 型名 | 必須キー | nullable キー | 配列キー | 対応 API |
|------|----------|---------------|----------|----------|
| `StatusObject` | `last_sha`, `last_build_at`, `last_build_status`, `last_trigger`, `last_deploy_status`, `pending_transfers_count`, `notify_pending_count`, `circuit_open`, `output_url`, `running` | `last_sha`, `last_build_at`, `last_trigger`, `last_deploy_status`, `output_url` | なし | `GET /api/status` |
| `HistoryPageObject` | `total`, `page`, `per_page`, `pages`, `history` | なし | `history: HistoryRecord[]` | `GET /api/history` |
| `HistoryRecord` | `id`, `build_at`, `sha`, `status`, `trigger`, `output_size_bytes`, `duration_seconds`, `flagged`, `tags` | `sha`, `output_size_bytes`, `duration_seconds`, `comment`, `output_sha256`, `rollback_from` | `tags: string[]` | `.build_history` |
| `HistoryLogObject` | `.build_logs/{id}.json` の全必須キー、`lines` | `.build_logs/{id}.json` の nullable キー | `stdout`, `stderr`, `warnings`, `tags`, `lines` | `GET /api/history/{id}/log` |
| `CommentObject` | `id`, `comment`, `updated_at` | `comment`, `updated_at` | なし | `GET /api/history/{id}/comment` |
| `ConfigObject` | `.server_config` schema の全キー | `pat_expires_at`, `allowed_hours` | なし | `GET /api/config` |
| `ConfigValidationObject` | `valid`, `config`, `errors`, `warnings` | なし | `errors`, `warnings` | `POST /api/config/validate` |
| `AuditLogRecord` | `.audit_log` schema の全必須キー | `actor_id`, `target_id`, `remote_addr`, `message` | なし | `GET /api/audit-log` |
| `TotpStatus` | `enabled`, `confirmed_at` | `confirmed_at` | なし | `GET /api/auth/totp-status`, `POST /api/auth/totp-confirm`, `DELETE /api/auth/totp` |
| `ApiRateLimitPolicy` | `enabled`, `groups`, `state_summary` | なし | `groups`, `state_summary` | `GET/POST /api/api-rate-limit` |
| `NotifyConfig` | `webhooks`, `channels`, `on`, `summary`, `email` | `webhooks[].payload_template`, `webhooks[].secret` | `webhooks`, `channels`, `on`, `email.to`, `email.on` | `GET /api/notify-config` |
| `RepoInfoObject` | `owner`, `repo`, `branch`, `target_file` | なし | なし | `GET /api/repo-info` |
| `BranchTargetRecord` | `branch`, `target_file`, `sha_file`, `src`, `out`, `deploy_targets` | なし | `deploy_targets` | `GET /api/branch-config` |
| `SysinfoObject` | `output_size_bytes`, `output_mtime`, `uptime_seconds` | `output_mtime` | なし | `GET /api/sysinfo` |
| `HealthObject` | `status`, `last_build_at`, `last_build_status`, `last_deploy_at`, `last_deploy_status`, `pending_transfers`, `uptime_seconds` | `last_build_at`, `last_deploy_at` | なし | `GET /api/health` |
| `ScheduleObject` | `next_run_at`, `interval`, `paused`, `allowed_hours` | `next_run_at`, `allowed_hours` | なし | `GET /api/schedule` |
| `StatsObject` | `days`, `total_builds`, `success_count`, `failure_count`, `success_rate`, `avg_interval_minutes`, `avg_duration_seconds`, `max_duration_seconds` | `avg_interval_minutes`, `avg_duration_seconds`, `max_duration_seconds` | なし | `GET /api/stats` |
| `TimelineObject` | `days`, `timeline` | なし | `timeline` | `GET /api/stats/timeline` |
| `BuildDurationStats` | `n`, `count`, `avg_seconds`, `min_seconds`, `max_seconds`, `recent` | `avg_seconds`, `min_seconds`, `max_seconds` | `recent` | `GET /api/stats/build-duration` |
| `OutputMetaObject` | `size_bytes`, `mtime`, `heading_count`, `size_diff_bytes`, `tables_count`, `code_blocks_count`, `build_warnings`, `size_warn`, `sha256`, `build_id`, `commit_sha`, `build_at` | `mtime`, `size_diff_bytes`, `tables_count`, `code_blocks_count`, `sha256` | `build_warnings` | `GET /api/output-meta` |
| `DashboardObject` | `status`, `sysinfo`, `stats`, `schedule`, `alerts` | なし | `alerts` | `GET /api/dashboard` |
| `DiagnosticsObject` | `checked_at`, `items` | なし | `items` | `GET /api/diagnostics` |
| `RateLimitObject` | `limit`, `remaining`, `reset_at`, `used` | なし | なし | `GET /api/rate-limit` |
| `DiskUsageObject` | `build_logs_bytes`, `build_logs_count`, `build_logs_archive_bytes`, `build_logs_archive_count`, `output_file_bytes`, `total_bytes` | なし | なし | `GET /api/disk-usage` |
| `AccessRecord` | `.access_log` schema の全必須キー | `session_id`, `token_id`, `remote_addr`, `reason` | なし | `GET /api/access-log` |
| `ApiAccessRecord` | `.api_access_log` schema の全必須キー | `actor`, `remote_addr`, `user_agent`, `error` | なし | `GET /api/api-access-log` |
| `ConfigLogRecord` | `.config_log` schema の全必須キー | なし | なし | `GET /api/config-log` |
| `NotifyRecord` | `at`, `event`, `http_status`, `result`, `attempt`, `error` | `http_status`, `error` | なし | `GET /api/notify-log` |
| `SessionRecord` | `created_at`, `expires_at`, `last_used_at`, `current` | `last_used_at` | なし | `GET /api/sessions` |
| `WebhookEventRecord` | `timestamp`, `delivery_id`, `event`, `ref`, `sha`, `build_triggered` | `delivery_id`, `sha` | なし | `GET /api/webhook-events` |
| `SnapshotRecord` | `id`, `build_id`, `saved_at`, `size_bytes` | なし | なし | `GET /api/snapshots` |
| `MaintenanceObject` | `enabled`, `reason`, `since` | `reason`, `since` | なし | `GET /api/maintenance` |
| `QueueEntry` | `id`, `trigger`, `queued_at`, `requested_by`, `payload` | なし | なし | `GET /api/queue` |
| `TokenRecord` | `id`, `label`, `scopes`, `created_at`, `last_used_at`, `expires_at`, `revoked_at` | `last_used_at`, `expires_at`, `revoked_at` | `scopes` | `GET /api/tokens` |
| `TokenCreateResult` | `id`, `token`, `label`, `scopes`, `created_at`, `expires_at` | `expires_at` | `scopes` | `POST /api/tokens` |
| `HookRecord` | `id`, `phase`, `command_args`, `enabled`, `abort_on_failure`, `timeout_seconds` | なし | `command_args` | `GET/POST /api/hooks` |
| `HookRunRecord` | `build_id`, `ran_at`, `exit_code`, `output` | なし | なし | `GET /api/hooks/{id}/log` |
| `AlertRule` | `id`, `metric`, `operator`, `threshold`, `level`, `message` | なし | なし | `GET/POST /api/alert-rules` |
| `TagRule` | `id`, `condition`, `tags` | なし | `tags` | `GET/POST /api/tag-rules` |
| `PipelineConfig` | `extra_args`, `env` | なし | `extra_args` | `GET /api/pipeline-config` |
| `BuildChainConfig` | `chains` | なし | `chains` | `GET /api/build-chain-config` |
| `BuildTrendStats` | `count`, `avg_seconds`, `median_seconds`, `p95_seconds`, `anomaly_count`, `samples` | `avg_seconds`, `median_seconds`, `p95_seconds` | `samples` | `GET /api/stats/build-trends` |
| `SmtpConfig` | `host`, `port`, `user`, `tls`, `from`, `to`, `on`, `enabled`, `password_set` | `host`, `user`, `from` | `to`, `on` | `GET /api/smtp-config` |
| `ApprovalRecord` | `id`, `status`, `branch`, `sha`, `target`, `created_at`, `expires_at` | なし | なし | `GET /api/approvals` |
| `BackupObject` | `exported_at`, `server_config`, `notify_config` | なし | 設定内容に従う | `GET /api/backup` |
| `ExportObject` | `exported_at`, `history` | なし | `history` | `GET /api/history/export` |

---
