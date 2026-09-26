# Adlaire CI — SDK 詳細仕様

[`docs/details/sdk.md`](sdk.md) は `sdk` owner component の詳細本文責務として、`sdk` が主本文として持つ実装契約だけを扱う。

owner / collaborator 境界管理は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0b.1](../DETAIL_INDEX.md#0b1-owner-component-別-owner-collaborator-境界管理) に従う。`sdk` owner component の主本文であり、collaborator component の仕様は endpoint、response、error、security、UI 呼び出し境界、検証観点として参照する。fixture、expected、fake、実装検証証跡は [`docs/details/fixture.md`](fixture.md) fixture 証跡責務を参照する。

SDK が呼び出す API endpoint の method、path、request、response、error、認証要否は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) を参照する。[`docs/details/sdk.md`](sdk.md) 詳細本文責務は SDK 側の class、method、引数変換、transport、error、stream、token 破棄を定義する。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `sdk` |
| 実装主体 | [`admin/adlaire-ci-sdk.js`](../../admin/adlaire-ci-sdk.js) の単一 ES Module。 |
| 持つ内容 | `sdk` owner が主本文として定義する SDK class、method、HTTP 対応、query / body 生成、error、stream、token 破棄。 |
| 持たない内容 | API endpoint 実装、API endpoint の状態ファイル更新責務、UI DOM 詳細、状態 schema、状態ファイル直接操作、admin 静的配信、setup / update 手順、release 生成・公開手順、fixture 証跡責務。 |

---

## 23. JavaScript SDK 仕様

[`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) は、sdk owner の JavaScript SDK 詳細本文責務である。実装 artifact は [責務境界](#0-責務境界)、module 形式は以下の SDK 実行環境契約を参照する。

**SDK 実行環境契約：**

| 項目 | 仕様 |
|------|------|
| JavaScript | ECMAScript 2022 以上を前提とする。transpile、bundle、polyfill は SDK 詳細本文責務に含めない。 |
| module | `admin/adlaire-ci-sdk.js` は ES Module とし、`export { AdlaireCI, AdlaireCIError }` を必須 export とする。default export は定義しない。 |
| browser API | `fetch`、`AbortController`、`ReadableStream.getReader()`、`TextDecoder` が存在する browser を必須環境とする。いずれかが存在しない場合、`AdlaireCI` constructor は `TypeError("Unsupported browser runtime")` を投げる。 |
| 非 browser runtime | browser API 行の必須 API が存在しない実行環境では、runtime 名を判定分岐せず、`AdlaireCI` constructor が `TypeError("Unsupported browser runtime")` を投げる。Node.js 専用 API、npm package、bundler、polyfill による補完は行わない。 |
| global 汚染 | `window.AdlaireCI` 等の global 代入を行わない。標準管理ツールは ES Module import で SDK を読み込む。 |
| 外部 consumer | 必須 browser API を提供する外部 consumer application は、利用側の framework または bundler から本 ES Module を import してよい。consumer 側の framework adapter、package manifest、bundler 設定、polyfill、framework runtime を本リポジトリ、SDK 配布物、標準管理 UI 配布物へ追加してはならない。 |
| stream 前提 | `streamBuild()` は native `EventSource` を使用しない。Authorization header を付与できる `fetch` streaming を必須実装とする。 |
| API 対応範囲 | [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の全 endpoint のうち、GitHub が直接送信する `POST /api/webhook` は SDK method 0 件、`POST /api/schedule/allowed-hours` は request body ありの `setAllowedHours()` と `{from:null,to:null}` を送る `clearAllowedHours()` の 2 件、それ以外は endpoint 表の SDK 列に記載された 1 method と対応させる。表外 method、対応 0 件、明示例外以外の複数 method を禁止する。 |

```js
class AdlaireCI {
  constructor({ baseUrl })
  // this._token でセッショントークンを管理。login() 後の全リクエストに自動付与

  login(password)                               // POST /api/login → {token?, must_change, totp_required?, ticket?}; token がある場合は this._token にセット
  loginTotp(ticket, code)                       // POST /api/login/totp → {token, must_change}; this._token にセット
  logout()                                      // POST /api/logout → Promise<{message: string}>; this._token をクリア
  changePassword(currentPassword, newPassword)  // POST /api/change-password → Promise<{message: string}>
  getAuditLog({ limit = 100, offset = 0, actor = null, action = null, result = null } = {}) // GET /api/audit-log → Promise<{log: AuditLogRecord[], total: number}>

  getStatus()               // GET /api/status               → Promise<StatusObject>
  triggerBuild()            // POST /api/build               → Promise<{message: string, queue_id: string, queued: true, dispatch: "requested"|"timer_fallback"}>
  getLogs(n = 100, q = '')  // GET /api/logs?n={n}&q={q}     → Promise<{lines: string[]}>
  getHistory({ page = 1, perPage = 20, trigger = null, tag = null, flagged = null, failureCategory = null } = {}) // GET /api/history?page={page}&per_page={perPage}&trigger={trigger}&tag={tag}&flagged={flagged}&failure_category={failureCategory} → Promise<HistoryPageObject>
  getSysinfo()              // GET /api/sysinfo              → Promise<SysinfoObject>
  getSchedule()             // GET /api/schedule             → Promise<ScheduleObject>
  getNotifyConfig()         // GET /api/notify-config        → Promise<NotifyConfig>
  setNotifyConfig(config)   // POST /api/notify-config       → Promise<{message: string}>
  getConfig()               // GET /api/config               → Promise<ConfigObject>
  validateConfig(config)    // POST /api/config/validate     → Promise<ConfigValidationObject>
  setConfig(config)         // POST /api/config              → Promise<{message: string, config: ConfigObject}>
  health()                  // GET /api/health               → Promise<HealthObject>
  getPatStatus()            // GET /api/pat-status           → Promise<PatStatusObject>
  getAccessLog({ limit = 100, offset = 0 } = {}) // GET /api/access-log?limit={limit}&offset={offset} → Promise<{log: AccessRecord[]}>
  getApiAccessLog({ limit = 100, offset = 0, method = null, path = null, status = null } = {}) // GET /api/api-access-log → Promise<{log: ApiAccessRecord[], total: number}>
  getStats(days = 7)        // GET /api/stats?days={days}    → Promise<StatsObject>
  exportLogs()              // GET /api/logs/export          → Promise<{exported_at: string, lines: string[]}>
  cleanupLogs()             // POST /api/logs/cleanup        → Promise<{message: string, deleted_count: number, failed_count: number}>
  archiveLogs()             // POST /api/logs/archive        → Promise<{message: string, archived_count: number}>
  getRepoInfo()             // GET /api/repo-info            → Promise<RepoInfoObject>
  backup()                  // GET /api/backup               → Promise<BackupObject>
  restore(config)           // POST /api/restore RestoreObject → Promise<{message: string}>
  notifyTest()              // POST /api/notify-test         → Promise<{message: string, channel_id: string}>
  buildForce()              // POST /api/build/force         → Promise<{message: string, queue_id: string, queued: true, dispatch: "requested"|"timer_fallback"}>
  patVerify()               // POST /api/pat-verify          → Promise<PatVerifyObject>
  getHistoryLog(id)         // GET /api/history/{id}/log     → Promise<HistoryLogObject>
  cancelBuild()             // POST /api/build/cancel        → Promise<{message: string}>
  resetCircuitBreaker()     // POST /api/circuit-breaker/reset → Promise<{message: string, open: boolean, consecutive_failures: number}>
  streamBuild(onLine, onEnd) // GET /api/build/stream (保存済み log の有限 SSE) → Promise<StreamHandle>（onLine(line), onEnd({status, duration_seconds:number|null}) コールバック、done で終端を監視）
  setLogLevel(level)        // POST /api/log-level           → Promise<{message: string, level: string}>
  updatePat(token)          // POST /api/pat-update          → Promise<{message: string}>
  getDashboard()            // GET /api/dashboard            → Promise<DashboardObject>
  getNotifyLog({ limit = 100, offset = 0 } = {}) // GET /api/notify-log?limit={limit}&offset={offset} → Promise<{log: NotifyRecord[]}>
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
  searchLogs(q = '', from = '', to = '', level = undefined) // GET /api/logs/search?q={q}&from={from}&to={to}&level={level} → Promise<SearchResult>
  getOutputMeta()              // GET /api/output-meta       → Promise<OutputMetaObject>
  getStatsTimeline(days = 30)  // GET /api/stats/timeline?days={days} → Promise<TimelineObject>
  getStatsBuildDuration(n = 10) // GET /api/stats/build-duration?n={n} → Promise<BuildDurationStats>
  getBuildTrends(n = 100)       // GET /api/stats/build-trends?n={n} → Promise<BuildTrendStats>
  getDiagnostics()             // GET /api/diagnostics       → Promise<DiagnosticsObject>
  getRateLimit()               // GET /api/rate-limit        → Promise<RateLimitObject>
  getDiskUsage()               // GET /api/disk-usage        → Promise<DiskUsageObject>
  getConfigLog({ limit = 100, offset = 0 } = {}) // GET /api/config-log?limit={limit}&offset={offset} → Promise<{log: ConfigLogRecord[]}>
  getApiRateLimit()            // GET /api/api-rate-limit    → Promise<ApiRateLimitPolicyResponse>
  setApiRateLimit(policy)      // POST /api/api-rate-limit ApiRateLimitPolicyInput → Promise<ApiRateLimitPolicyResponse>
  getBranchConfig()            // GET /api/branch-config     → Promise<{source: string, branches: BranchTargetRecord[]}>
  setBranchConfig(branches)    // POST /api/branch-config    → Promise<{message: string, branches_count: number}>
  notifyWeeklySummary()        // POST /api/notify/weekly-summary → Promise<{message: string, period: string, success_count: number, failure_count: number, success_rate: number}>
  getWebhookEvents(limit = 50, offset = 0) // GET /api/webhook-events?limit={limit}&offset={offset} → Promise<{events: WebhookEventRecord[], total: number}>
  getWebhookConfig()          // GET /api/webhook-config    → Promise<{configured: boolean}>
  setWebhookConfig(secret)    // POST /api/webhook-config   → Promise<{message: string}>
  getBuildChainConfig()       // GET /api/build-chain-config → Promise<BuildChainConfig>
  setBuildChainConfig(chains) // POST /api/build-chain-config → Promise<{message: string, chains_count: number}>
  getApprovals()              // GET /api/approvals         → Promise<{approvals: ApprovalRecord[]}>
  approveBuild(id)            // POST /api/approvals/{id}/approve → Promise<{message: string, queued: true, queue_id: string, dispatch: "requested"|"timer_fallback"}>
  rejectBuild(id)             // POST /api/approvals/{id}/reject  → Promise<{message: string}>
  getHistoryComment(id)        // GET /api/history/{id}/comment  → Promise<CommentObject>
  setHistoryComment(id, comment) // POST /api/history/{id}/comment → Promise<{message: string}>
  setRepoConfig({ owner = undefined, repo = undefined } = {}) // POST /api/repo-config RepoConfigPatch → Promise<{message: string}>
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
  getQueue()                   // GET /api/queue             → Promise<{active: QueueEntry|null, queued: QueueEntry[], max_size: number}>
  clearQueue()                 // DELETE /api/queue          → Promise<{message: string, cleared_count: number}>
  // 16D ダッシュボードレイアウト
  getDashboardLayout()         // GET /api/dashboard-layout  → Promise<{widgets: string[]}>
  setDashboardLayout(widgets)  // POST /api/dashboard-layout → Promise<{message: string}>
}

export { AdlaireCI, AdlaireCIError };
```

全メソッドは `Promise` を返す。`streamBuild` は SSE 接続確立後に `StreamHandle` で resolve し、接続前エラーは `AdlaireCIError` で reject する。HTTP エラー（4xx / 5xx）は `AdlaireCIError` としてスローする。`401` 受信時はセッション期限切れとして `this._token` をクリアする。constructor、private method、helper 関数を除く public method は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の SDK 列と完全一致させる。

**SDK 共通実装契約：**

| 項目 | 仕様 |
|------|------|
| `baseUrl` | API origin だけを表す absolute `http` / `https` URL とする。path は `/`、userinfo、query、fragment はなしとし、保持時は末尾 `/` を除去する。空文字、relative URL、`null`、`undefined`、その他の形式違反は HTTP 送信前に `TypeError("Invalid argument: baseUrl")` とする。`baseUrl` に `/api` を含めない。 |
| URL 組み立て | パスは `/api/...` をそのまま連結する。query key は各 method の固定名をそのまま使用し、query value は `encodeURIComponent(String(value))` の結果を使用する。空白は `%20` とし、`+` への変換を禁止する。 |
| 認証ヘッダー | `this._token` が存在する場合のみ `Authorization: Bearer ${token}` を付与する。 |
| JSON 送信 | `POST` / `DELETE` で body を送る場合は `Content-Type: application/json` を付与し、`JSON.stringify` した body を送信する。 |
| JSON 受信 | `Content-Type` が JSON の場合のみ `response.json()` を呼ぶ。[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) で JSON response を定義した endpoint の成功時に空 body を受信した場合は protocol error として `AdlaireCIError(status=0, message="Empty JSON response")` を投げる。 |
| `AdlaireCIError` | `name="AdlaireCIError"`、`status`、`message`、`details`、`responseBody` を持つ `Error` 派生クラスとする。constructor は `new AdlaireCIError({status, message, details = null, responseBody = null})` とし、network error、timeout、protocol error は `status=0` とする。`message` は API error response の `error`、network error は `"Network error"`、timeout は `"Request timeout"`、protocol error は固定文言を使用する。 |
| `logout()` | API 呼び出しが失敗しても `finally` で `this._token` をクリアする。 |
| request timeout | 通常 API の 30 秒は `fetch()` 開始から JSON text または Blob の body 読取完了までとする。response header 受信時に timeout を解除しない。SDK 自身の timeout で header 待機または body 読取が abort した場合は `AdlaireCIError(status=0, message="Request timeout")` を投げる。`streamBuild()` は接続確立まで 30 秒、接続確立後は timeout なしとし、利用者が `StreamHandle.close()` で停止する。 |
| `streamBuild()` | `onLine` と `onEnd` は function 必須とし、不正時は HTTP 送信前に `TypeError` とする。token がない場合は接続前に `AdlaireCIError(status=401, message="Unauthorized")` を投げる。native `EventSource` は Authorization header を付与できないため使用禁止とする。SDK は `fetch()`、`AbortController`、`ReadableStream` reader を使用し、`Accept: text/event-stream` と `Authorization` header を付与する。`2xx`、media type `text/event-stream`、readable body の確認後だけ `StreamHandle` で resolve する。`2xx` で media type または readable body が不正な場合は接続前に `AdlaireCIError(status=0,message="Invalid SSE response")` で reject する。 |
| SSE parser | `TextDecoder("utf-8",{fatal:true})` の streaming decode を使用し、chunk 境界をまたぐ byte、Unicode 文字、frame を buffer する。frame separator は LF 2 個の `\n\n` だけとし、1 frame は `data: ` で始まる 1 行と compact JSON 1 object だけを許可する。CRLF、追加行、空 frame、無効 UTF-8、JSON parse 不能、未知 key、未知 `type`、必須 key 不足、型不一致は `Invalid SSE frame` とする。 |
| SSE frame 適用 | `type="log"` は key が `type,line,at` だけ、`line` が string、`at` が UTC ISO 8601 秒精度の場合だけ `onLine(line)` を呼ぶ。`type="end"` は key が `type,status,duration_seconds` だけで、status と duration の組合せが API 契約に一致する場合だけ一時保持する。`end` は最終 frame 1 件だけを許可し、その後に EOF と空 buffer を確認してから `onEnd(summary)` を 1 回呼ぶ。実行中 snapshot の `duration_seconds:null` を保持し、完了と推測しない。EOF 前の `end` 不在、`end` 後の byte / frame、EOF 時の未完 frame は `Invalid SSE frame` とする。`onLine` または `onEnd` が例外を投げた場合は reader を停止し、`AdlaireCIError(status=0,message="Stream callback failed")` へ置き換える。 |
| `StreamHandle` | `StreamEnd.status` は `"running"`、`"success"`、`"failure"`、`"cancelled"` のいずれか、`StreamEnd.duration_seconds` は number または `null` とし、exact object は `{status,duration_seconds}` とする。`streamBuild()` の戻り値は `close()`、`closed`、`error`、`done` だけを持ち、`error` は `AdlaireCIError` または `null`、`done` は `StreamEnd` または `null` で resolve する Promise とする。正常終了時は `onEnd(summary)` 後に `done` を同じ summary で resolve し、`closed=true`、`error=null` とする。接続後 error は `closed=true`、`error` に同じ `AdlaireCIError` を保持し、`done` をその error で reject する。`close()` は reader と AbortController を停止し、複数回呼んでも例外を投げず、`onEnd` を呼ばず、`closed=true`、`error=null`、`done` は `null` で resolve する。normal end、error、user close のうち最初の terminal transition だけが状態と `done` を確定し、後続 callback と再 settle を禁止する。 |
| Blob レスポンス | `downloadSnapshot(id)` の `2xx` のみ `response.blob()` を使用する。成功 response の `Content-Type` を media type と parameter に分解し、media type が case-insensitive で `application/octet-stream` と完全一致する場合だけ Blob を返す。不一致は `AdlaireCIError(status=0,message="Invalid binary response")` とする。その他の成功 response は JSON とする。 |
| メソッド引数検証 | SDK 側でも必須引数の空値、配列型、数値範囲を検証し、HTTP 送信前に `TypeError` を投げる。 |
| endpoint 対応 | SDK method は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の SDK 列に存在する endpoint だけを呼び出す。[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) にない endpoint を SDK 独自判断で追加してはならない。 |
| body なし endpoint | [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の `Request` が `none` の場合、SDK は `fetch` に `body` を設定しない。`{}` も送信しない。 |
| token 保存 | セッショントークンはメモリ上の `this._token` のみに保持する。`localStorage`、`sessionStorage`、Cookie へ保存しない。 |
| 秘密情報引数 | `updatePat(token)`、`setWebhookConfig(secret)`、SMTP password、`createToken()` の返却 token は console 出力しない。 |
| query 生成 | `undefined`、`null`、空文字の任意 query は送信しない。ただし `q`、`from`、`to` は endpoint 仕様で空文字を有効値として定義している場合だけ、空文字を query value として送信する。 |
| 戻り値補完禁止 | 成功時は API response にない key を追加しない。失敗時は HTTP status、API error、details、responseBody 以外を推測しない。fallback 値は API response に含まれる値だけを返し、表示用加工は UI 側で行う。 |
| retry | SDK は自動 retry を行わない。ユーザー操作による再実行、または UI の明示的な再取得のみを許可する。 |

**SDK transport / error 固定契約：**

SDK の内部 request helper は、すべての public method で [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) の固定表の処理順に固定する。public method ごとに個別 fetch 処理を複製してはならない。

| 順序 | 処理 | 固定仕様 |
|------|------|----------|
| 1 | 引数検証 | 必須引数、型、空配列、数値範囲を検証する。失敗時は `TypeError` を投げ、`fetch()` を呼ばない。 |
| 2 | URL 生成 | `baseUrl + path + query` を生成する。query key は [SDK 引数変換契約](#sdk-argument-contract) の対象 method 行に記載された順序で追加する。 |
| 3 | body 生成 | `Request=none` では `body` と `Content-Type` を設定しない。JSON body ありの場合だけ `JSON.stringify()` する。 |
| 4 | header 生成 | `Accept` を常に付与する。JSON body を送信する場合だけ `Content-Type` を付与する。token が空でない場合だけ `Authorization` を付与し、token が空の場合は `Authorization` header を付けない。 |
| 5 | timeout 設定 | 通常 request は `AbortController` で 30 秒 timeout を開始し、手順 7 の body 読取またはそれ以前の失敗が確定するまで維持する。`streamBuild()` は接続確立まで 30 秒。 |
| 6 | `fetch()` 実行 | 手順 5 で SDK 自身が設定した timeout により `AbortController` が abort した reject だけを `AdlaireCIError(status=0,message="Request timeout")` とする。timeout 以外の理由で `fetch()` が reject した場合は、例外名や browser 文言を公開せず `AdlaireCIError(status=0,message="Network error")` とする。`streamBuild().close()` による user close は reject として扱わず、StreamHandle 固定契約の terminal transition を適用する。 |
| 7 | response body 読取 / parse | 成功 / 失敗に関わらず、JSON は `response.text()`、binary は `response.blob()` の完了まで timeout 対象とする。body 読取が SDK timeout で reject した場合は `Request timeout`、その他の理由で reject した場合は `Network error` とする。JSON error body がある場合は parse し、parse 不能 error body は `responseBody` に text を保持し、`message` は HTTP status 固定文言とする。 |
| 8 | timeout 解除 | body 読取と parse、または先行失敗が確定した terminal path の `finally` で 1 回だけ解除する。response header 受信直後の解除を禁止する。 |
| 9 | token 変化 | `401` の場合だけ `this._token = null`。`403`、`429`、`500`、network error、timeout では token を破棄しない。 |
| 10 | return / throw | `2xx` は endpoint の型で返す。`4xx` / `5xx` は `AdlaireCIError` を投げる。 |

HTTP status と SDK error の対応は [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) の固定表に固定する。

| 条件 | `AdlaireCIError.status` | `message` | `details` | token |
|------|-------------------------|-----------|-----------|-------|
| API JSON error | HTTP status | response の `error` | response の `details` が配列なら保持 | `401` のみ破棄 |
| API error body なし | HTTP status | `HTTP {status}` | `null` | `401` のみ破棄 |
| API error JSON parse 不能 | HTTP status | `HTTP {status}` | `null` | `401` のみ破棄 |
| network error | `0` | `Network error` | `null` | 維持 |
| timeout | `0` | `Request timeout` | `null` | 維持 |
| empty success JSON | `0` | `Empty JSON response` | `null` | 維持 |
| invalid success JSON | `0` | `Invalid JSON response` | `null` | 維持 |
| invalid binary response | `0` | `Invalid binary response` | `null` | 維持 |
| invalid SSE response | `0` | `Invalid SSE response` | `null` | 維持 |
| invalid SSE frame | `0` | `Invalid SSE frame` | `null` | 維持 |
| stream callback failure | `0` | `Stream callback failed` | `null` | 維持 |

`429` は SDK で自動待機、自動再送、自動 refresh を行わない。binary response は `downloadSnapshot(id)` の `2xx` のみ `Blob` とする。`4xx` / `5xx` では `Content-Type` を media type と parameter に分解し、type と subtype を case-insensitive で判定する。type が `application` で subtype が `json`、または subtype が 1 文字以上の prefix と `+json` から成る場合だけ JSON error として parse する。parse 成功時は response の `error` / `details` を保持した `AdlaireCIError` を投げる。JSON parse 不能、JSON 以外の error body、空 body の場合は body text を `responseBody` に保持し、`message` は `HTTP {status}` とする。`streamBuild()` は接続後の `close()` をユーザー停止として扱い、`AdlaireCIError` を投げない。接続後の network / frame / callback error は `StreamHandle.error` と `StreamHandle.done` の reject で同じ `AdlaireCIError` を通知する。callback が投げた元の error 本文は `responseBody`、console、UI へ転写しない。

<a id="sdk-method-implementation-contract"></a>
**SDK メソッド実装固定契約：**

| 項目 | 仕様 |
|------|------|
| public method 定義順 | class 内の public method は [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) の一覧順に定義する。追加 public method を末尾に置くことは禁止し、先に [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) と本一覧を更新する。 |
| private helper | private helper は `_request`, `_json`, `_query`, `_requireToken`, `_validateId`, `_clearTokenOn401` だけを定義する。helper を export しない。 |
| TypeError 文言 | SDK 側引数検証の `TypeError.message` は `"Invalid argument: <name>"` に固定する。複数不正がある場合は最初に検出した引数だけを返す。 |
| path parameter | `id` を path に入れる method は、`id` が string かつ [`docs/details/api.md` 詳細本文責務 §22.0b](api.md#sec-22-0b) の `^[A-Za-z0-9_-]{1,64}$` に完全一致することを SDK 側で検証する。不一致は HTTP 送信前に `TypeError("Invalid argument: id")` とする。検証成功後の値に `encodeURIComponent(id)` を 1 回だけ適用し、1 segment として連結する。 |
| query parameter | query key は [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) SDK 引数変換契約の表記順で生成する。value は各々 `encodeURIComponent(String(value))` を 1 回だけ適用し、空白を `%20` とする。`URLSearchParams` 等による `+` 変換、2 重 encode、並べ替えを禁止する。任意 query が未指定の場合、`?` 自体を付けない。 |
| body parameter | body object の key 順は [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) SDK 引数変換契約の送信値順とする。未知 key を SDK が追加しない。 |
| token mutation | `login()` と `loginTotp()` は response に `token` が存在する場合だけ `this._token` を更新する。`totp_required:true` かつ token なしの場合は既存 token を保持せず `null` にする。 |
| logout failure | `logout()` は network error、`401`、`500` のいずれでも `finally` で `this._token=null` にする。 |
| response passthrough | 成功 response は clone、整形、既定値 merge を行わず、そのまま返す。Blob と StreamHandle は例外とする。 |
| error details | API error response の `details` が配列なら `AdlaireCIError.details` に同じ配列を保持する。配列でなければ `null`。 |
| responseBody | JSON parse できた error body は object のまま、parse 不能 error body は先頭 4000 文字の string として `responseBody` に保持する。 |

**SDK 検証 fixture 参照：**

SDK の fixture 名、fake fetch 入力、expected、error shape、stream frame、binary response、実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約)、[`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約)、[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) を正本とする。[`docs/details/sdk.md`](sdk.md) 詳細本文責務では、SDK method、引数変換、token mutation、response passthrough、error 変換、stream handle の実装契約だけを扱う。

<a id="sec-3"></a>
**基本運用 SDK 操作固定契約：**

基本運用 SDK は、以下の SDK method の request、response、error 変換を固定する。SDK は成功時 response を endpoint schema の範囲でそのまま返し、失敗時は HTTP status、API error、details、responseBody を保持した `AdlaireCIError` へ変換する。UI が必要とする表示用既定値、並べ替え、ラベル変換は SDK で行わない。

| SDK method | HTTP | 成功時 | 失敗時 | 追加禁止条件 |
|------------|------|--------|--------|--------------|
| `login(password)` | `POST /api/login` | `token` がある場合だけ `this._token` へ保存する。`totp_required:true` の場合は token を保存せず response を返す。 | `401`、`429`、`500` は `AdlaireCIError`。`401` で既存 token を破棄する。 | password を console、error、responseBody 加工結果へ出さない。 |
| `logout()` | `POST /api/logout` | response に関わらず `finally` で token を破棄する。 | network error、`401`、`500` でも token 破棄後に error を投げる。 | logout 失敗を理由に token を保持しない。 |
| `getStatus()` | `GET /api/status` | `StatusObject` をそのまま返す。 | `500 State file is corrupted` / `State file read failed` を message として保持する。 | `.build_status.json` 不在時の fallback 値を SDK が推測しない。 |
| `triggerBuild()` | `POST /api/build` | `{message,queue_id,queued:true,dispatch}` をそのまま返す。 | `409`、`429`、`503` は `AdlaireCIError.status` に HTTP status を保持し、`message` は API response の `error` を保持する。 | dispatch から running / build id を推測せず、runner 起動を SDK から再実行しない。 |
| `buildForce()` | `POST /api/build/force` | `{message,queue_id,queued:true,dispatch}` をそのまま返す。 | `409`、`429`、`503` を `AdlaireCIError`。 | force 可否、SHA 更新、build 開始を SDK 側で状態推測しない。 |
| `cancelBuild()` | `POST /api/build/cancel` | `{message}` を返す。 | running なしの `409` を `AdlaireCIError(status=409)`。 | cancel 後に SDK が自動 `getStatus()` を呼ばない。 |
| `getLogs(n,q)` | `GET /api/logs` | `{lines}` を返す。`lines` は API 順序を保持する。 | `500` は固定 error message を保持する。 | line を結合、trim、level 分類しない。 |
| `getHistory(options)` | `GET /api/history` | `{total,page,per_page,pages,warnings,history}` を返す。 | query 不正 `422` は `details` を保持する。 | `pages`、`total`、`warnings` を SDK 側で再計算または補完しない。 |
| `getHistoryLog(id)` | `GET /api/history/{id}/log` | log object を返す。 | `404`、`500` を `AdlaireCIError`。 | archive fallback を SDK 側で再試行しない。 |
| `getQueue()` | `GET /api/queue` | `{active,queued,max_size}` を返し、active と waiting の区分を保持する。 | `.build_state` 破損の `500` を固定 message で保持する。 | active を queued に混ぜる、queue 並び替え、重複排除、上限補正をしない。 |
| `resetCircuitBreaker()` | `POST /api/circuit-breaker/reset` | `{message,open,consecutive_failures}` を返す。 | 破損状態 `500` を `AdlaireCIError`。 | reset 成功後に SDK が自動 build を開始しない。 |
| `streamBuild(onLine,onEnd)` | `GET /api/build/stream` | `StreamHandle` を返し、`log` frame を `onLine`、検証済み最終 `end` を `onEnd` と `done` へ渡す。 | 接続前 `401`、`404`、`500`、timeout、invalid response は `streamBuild()` を reject。接続後 error は `done` を reject。 | `EventSource`、自動 reconnect、log 永続化を行わない。 |

<a id="sec-3-2"></a>
**基本運用 SDK fixture 参照：**

基本運用 SDK の fake fetch 入力、expected request、expected return、`AdlaireCIError`、stream frame、unauthorized token clear は [`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) を正本とする。

<a id="sec-4"></a>
**拡張運用 SDK 操作固定契約：**

拡張運用 SDK は、[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の endpoint 契約と [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) SDK 引数変換契約だけに従う。SDK は保存前検証の一部を `TypeError` で行う場合でも、検証対象は必須引数、型、範囲、path parameter 形式に限定する。API response の補完、no-op 判定、secret mask 変換、状態ファイル由来値の再計算を行ってはならない。

| 機能群 | SDK method | 成功時 | 失敗時 | 追加禁止条件 |
|--------|------------|--------|--------|--------------|
| config / repo / branch | `getConfig()`, `setConfig(config)`, `validateConfig(config)`, `getRepoInfo()`, `setRepoConfig({owner,repo})`, `getBranchConfig()`, `setBranchConfig(branches)` | API response をそのまま返す。`set*` は success message と count / config を保持する。 | `422 details` は `AdlaireCIError.details` に保持する。`500` は固定 message。 | `getRepoInfo()` は `RepoInfoObject` を補完せず返す。`setRepoConfig` に branch / target field を合成せず、`setBranchConfig` に owner / repo を合成しない。SDK 側で未知 key を削除したり既定値 merge したりしない。 |
| schedule | `getSchedule()`, `setScheduleInterval()`, `pauseSchedule()`, `resumeSchedule()`, `setAllowedHours()`, `clearAllowedHours()`, `setForceInterval()`, `setBuildCooldown()` | response の interval / hours / seconds / allowed_hours をそのまま返す。 | systemd 更新失敗の `500` を `AdlaireCIError` にする。 | timer 状態を SDK が推測しない。 |
| diagnostics / dashboard | `getDiagnostics()`, `getDashboard()`, `getRateLimit()`, `getDiskUsage()`, `getOutputMeta()` | read-only response をそのまま返す。 | item 単位 warn/error は成功 response として返し、HTTP error だけ例外にする。 | read-only response を SDK が保存・集計しない。 |
| notify / SMTP / webhook | `getNotifyConfig()`, `setNotifyConfig()`, `notifyTest()`, `notifyWeeklySummary()`, `getWebhookEvents()`, `getWebhookConfig()`, `setWebhookConfig()`, `getSmtpConfig()`, `setSmtpConfig()`, `smtpTest()` | secret は API が mask した値だけ返す。 | 未設定 `422` / `501`、送信失敗 `500` を保持する。 | secret 平文を console、throw message、responseBody 加工結果へ出さない。 |
| snapshots / rollback | `getSnapshots()`, `downloadSnapshot(id)`, `deleteSnapshot(id)`, `rollbackHistory(id)` | list は API 順序、download は Blob、rollback は `{message,build_id}`。 | `404`、running `409` を `AdlaireCIError`。 | SDK が snapshot 存在確認や rollback 可否を事前推測しない。 |
| maintenance / access / hooks | `getMaintenance()`, `enableMaintenance()`, `disableMaintenance()`, `getAccessControl()`, `setAccessControl()`, `getHooks()`, `addHook()`, `deleteHook()`, `getHookLog()` | API response をそのまま返す。 | access deny `403`、validation `422`、hook timeout `500` を保持する。 | command args を文字列結合しない。CIDR を SDK 独自正規化しない。 |
| alert / tag / pipeline / notes / dashboard layout | `getAlertRules()`, `addAlertRule()`, `deleteAlertRule()`, `getTagRules()`, `addTagRule()`, `deleteTagRule()`, `getPipelineConfig()`, `setPipelineConfig()`, `getNotes()`, `setNotes()`, `getDashboardLayout()`, `setDashboardLayout()` | API response の rules / config / notes / widgets をそのまま返す。 | duplicate `409`、validation `422`、read failure `500` を保持する。 | rule 重複排除、widget 補完、notes trim を行わない。 |
| tokens / sessions / audit | `getTokens()`, `createToken()`, `revokeToken()`, `getSessions()`, `revokeAllSessions()`, `getAuditLog()`, `getApiAccessLog()` | `createToken()` の token 本体は response として 1 回だけ返す。 | `403`、`404`、`422`、`429` を status 付きで保持する。 | token 本体を保存しない。token list に作成時 token を合成しない。 |

<a id="sec-4-2"></a>
**拡張運用 SDK fixture 参照：**

拡張運用 SDK の fake fetch 入力、expected request、expected return、secret mask、webhook paging、snapshot binary response は [`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) および [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) を正本とする。

<a id="sdk-argument-contract"></a>
**SDK 引数変換契約：**

SDK method は、[`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) の固定表の通りに引数を path、query、body へ変換する。[`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) の固定表にない引数、既定値、body key を追加してはならない。

| SDK method | 引数 | 変換先 | 送信値 |
|------------|------|--------|--------|
| `login(password)` | `password` | body | `{password}` |
| `loginTotp(ticket,code)` | `ticket`, `code` | body | `{ticket,code}` |
| `changePassword(currentPassword,newPassword)` | `currentPassword`, `newPassword` | body | `{current_password: currentPassword, new_password: newPassword}` |
| `getAuditLog({limit,offset,actor,action,result})` | `limit=100`, `offset=0`, `actor=null`, `action=null`, `result=null` | query | `limit`、`offset` は常に送信する。任意値は `null` の場合送信しない。 |
| `getAccessLog({limit,offset})` | `limit=100`, `offset=0` | query | `limit`、`offset` を常に送信する。 |
| `getLogs(n,q)` | `n=100`, `q=""` | query | `n`、`q`。`q` は空文字でも送信する。 |
| `getHistory({page,perPage,trigger,tag,flagged,failureCategory})` | `page=1`, `perPage=20`, `trigger=null`, `tag=null`, `flagged=null`, `failureCategory=null` | query | `page`、`per_page: perPage` は常に送信する。`trigger`、`tag`、`flagged`、`failure_category: failureCategory` は `null` の場合は送信しない。 |
| `getHistoryLog(id)` | `id` | path | path `{id}`。body と query は送信しない。 |
| `streamBuild(onLine,onEnd)` | `onLine`, `onEnd` | local callback | callback は request に含めず、GET body と query は送信しない。 |
| `getStats(days)` | `days=7` | query | `days`。 |
| `getStatsTimeline(days)` | `days=30` | query | `days`。 |
| `getStatsBuildDuration(n)` | `n=10` | query | `n`。 |
| `getBuildTrends(n)` | `n=100` | query | `n`。 |
| `getApiAccessLog({limit,offset,method,path,status})` | `limit=100`, `offset=0`, `method=null`, `path=null`, `status=null` | query | `limit`、`offset` は常に送信する。`method`、`path`、`status` は `null` の場合は送信しない。 |
| `getNotifyLog({limit,offset})` | `limit=100`, `offset=0` | query | `limit`、`offset` を常に送信する。 |
| `getConfigLog({limit,offset})` | `limit=100`, `offset=0` | query | `limit`、`offset` を常に送信する。 |
| `validateConfig(config)` | `config` | body | `config` をそのまま送信する。保存は API が行わない。 |
| `setNotifyConfig(config)` | `config` | body | `config` をそのまま送信する。 |
| `setConfig(config)` | `config` | body | `config` をそのまま送信する。未知 key は送信前に削除せず、API の `422` に委ねる。 |
| `setLogLevel(level)` | `level` | body | `{level}` |
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
| `restore(config)` | `config` | body | `RestoreObject` として `config` をそのまま送信する。 |
| `setBranchConfig(branches)` | `branches` | body | `{branches}` |
| `setBuildChainConfig(chains)` | `chains` | body | `{chains}` |
| `approveBuild(id)` | `id` | path | path `{id}`。body は送信しない。 |
| `rejectBuild(id)` | `id` | path | path `{id}`。body は送信しない。 |
| `getHistoryComment(id)` | `id` | path | path `{id}`。body と query は送信しない。 |
| `setRepoConfig({owner,repo})` | `owner=undefined`, `repo=undefined` | body | `undefined` でない key だけで `RepoConfigPatch` を作る。両方 `undefined` は HTTP 送信前に `TypeError`。branch / target key は追加しない。 |
| `rollbackHistory(id)` | `id` | path | path `{id}`。body は送信しない。 |
| `revokeToken(id)` | `id` | path | path `{id}`。body は送信しない。 |
| `downloadSnapshot(id)` | `id` | path | path `{id}`。body と query は送信しない。 |
| `deleteSnapshot(id)` | `id` | path | path `{id}`。body は送信しない。 |
| `enableMaintenance(reason)` | `reason` | body | `{reason}` |
| `setAccessControl(allowList)` | `allowList` | body | `{allow: allowList}` |
| `deleteHook(id)` | `id` | path | path `{id}`。body は送信しない。 |
| `getHookLog(id)` | `id` | path | path `{id}`。body と query は送信しない。 |
| `deleteAlertRule(id)` | `id` | path | path `{id}`。body は送信しない。 |
| `deleteTagRule(id)` | `id` | path | path `{id}`。body は送信しない。 |

path に入る `id` の型、形式、失敗、encode 回数は [SDK メソッド実装固定契約](#sdk-method-implementation-contract) の `path parameter` 行に従う。本引数変換一覧は、各 method の `id` が path 引数であることだけを定義し、検証契約を再定義しない。

**SDK method 完全性検証契約：**

SDK 詳細実装確認では、[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の SDK 列に記載された method 名と `AdlaireCI.prototype` の public method 名が一致しなければならない。constructor、private method、helper 関数、`AdlaireCIError` は比較対象外とする。

| 検証項目 | 合格条件 |
|----------|----------|
| endpoint coverage | [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の SDK 列で `none` 以外の method がすべて `AdlaireCI.prototype` に存在する。 |
| extra method | [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の SDK 列に存在しない public method がない。 |
| request shape | 各 method が [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の Request と [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) SDK 引数変換契約どおりの path / query / body を生成する。 |
| response handling | JSON endpoint は JSON object を返し、binary endpoint は `Blob`、stream endpoint は `StreamHandle` を返す。 |
| error handling | 4xx / 5xx、network error、timeout、empty JSON、invalid SSE frame が `AdlaireCIError` になる。 |
| token handling | `login()` 成功で token を保持し、`logout()` と `401` で token を破棄する。 |

<a id="sec-27-21"></a>
<a id="sec-27-47"></a>
**[§27.21〜§27.47 SDK 連動実装確認固定契約](sdk.md#sec-27-21)：**

[`docs/details/sdk.md` 詳細本文責務 §27.21](sdk.md#sec-27-21)〜[§27.47](sdk.md#sec-27-47) の追加仕様化機能で SDK の詳細実装確認を満たすには、対象 owner component 別の [`docs/details/*.md`](../details/) 詳細本文責務、[`docs/details/api.md` 詳細本文責務 §27](api.md#27-api-owner-追加仕様化機能-詳細仕様) の連動参照表、[`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) SDK 引数変換契約、SDK method 完全性検証契約、[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) を同時に満たす。SDK は API の補助層であり、API response の補完、状態推測、保存済み値の再計算、UI 表示用変換、自動 retry、自動 refresh、状態ファイル直接操作を行ってはならない。

| 対象 | SDK method | request 固定 | success 固定 | error 固定 | 禁止条件 |
|------|------------|--------------|---------------|------------|----------|
| [`docs/details/runner.md` 詳細本文責務 §27.21](runner.md#sec-27-21) / [`docs/details/runner.md` 詳細本文責務 §27.31](runner.md#sec-27-31) branch target / env | `getBranchConfig()`, `setBranchConfig(branches)`, `getConfig()`, `setConfig(config)` | `branches` または `config` を指定 key のまま送信する。`target_files`、`env`、secret 風 key を削除しない。 | API response の branch 配列、`source`、件数をそのまま返す。 | validation `422` details を保持する。 | target path 正規化、env key 並び替え、secret mask 判定を SDK が行わない。 |
| [`docs/details/runner.md` 詳細本文責務 §27.22](runner.md#sec-27-22) / [`docs/details/runner.md` 詳細本文責務 §27.27](runner.md#sec-27-27) pipeline / hook | `getPipelineConfig()`, `setPipelineConfig(config)`, `getHooks()`, `addHook()`, `deleteHook(id)`, `getHookLog(id)` | pipeline config はそのまま送る。hook `commandArgs` は配列のまま送る。 | config / hook / log response をそのまま返す。 | `409`、`422`、`500` を `AdlaireCIError` として保持する。 | shell 文字列結合、quote 展開、reserved arg 除去、env 補完を行わない。 |
| [`docs/details/runner.md` 詳細本文責務 §27.23](runner.md#sec-27-23)〜[§27.26](runner.md#sec-27-26) local watch / tag / cache / parallel | `getConfig()`, `setConfig(config)`, `getStatus()`, `getHistory()`, `getHistoryLog(id)` | config と query を表どおり送る。 | status / history / log に含まれる追加 key を削除しない。 | `422` details、`500` message を保持する。 | local watch 可否、tag match、cache hit、parallel result を SDK が再判定しない。 |
| [`docs/details/api.md` 詳細本文責務 §27.30](api.md#sec-27-30) / [`docs/details/runner.md` 詳細本文責務 §27.30](runner.md#sec-27-30) approval | `getApprovals()`, `approveBuild(id)`, `rejectBuild(id)`, `getQueue()` | `id` は path parameter 契約で encode し、approve / reject に body を送らない。 | approve は `{message,queued,queue_id,dispatch}`、reject は `{message}`、list は `requested_trigger` / `requested_force` を含む API 配列順で返す。 | `409` / `429` / `500` を保持し、自動再送しない。 | expired / rejected / approved / requested_force / dispatch を SDK が書き換えない。approve 後に SDK が自動 queue 取得しない。 |
| [`docs/details/runner.md` 詳細本文責務 §27.32](runner.md#sec-27-32) notification | `getNotifyConfig()`, `setNotifyConfig(config)`, `getNotifyLog()`, `notifyTest()`, `notifyWeeklySummary()`, `getSmtpConfig()`, `setSmtpConfig(config)`, `smtpTest()` | config は unknown key を削除せず API へ送る。password 未指定時だけ `setSmtpConfig()` は password key を送らない。 | API が返した mask 値、log、pending 状態をそのまま返す。 | 未設定 `422` / `501`、送信失敗 `500` を保持する。 | secret 平文を console、error message、responseBody 加工結果、SDK field に保存しない。 |
| [`docs/details/runner.md` 詳細本文責務 §27.33](runner.md#sec-27-33) / [`docs/details/runner.md` 詳細本文責務 §27.38](runner.md#sec-27-38) trend / anomaly | `getBuildTrends(n)`, `getStatsBuildDuration(n)`, `getDashboard()`, `getConfig()`, `setConfig(config)` | `n` は number として query へ送る。config はそのまま送る。 | summary、samples、anomaly flag を API 値のまま返す。`BuildTrendStats` に `warnings` を補完しない。 | invalid `n` / config `422` を保持する。 | avg / median / p95 / anomaly を SDK が再計算しない。 |
| [`docs/details/runner.md` 詳細本文責務 §27.34](runner.md#sec-27-34) / [`docs/details/runner.md` 詳細本文責務 §27.35](runner.md#sec-27-35) chain / queue | `getBuildChainConfig()`, `setBuildChainConfig(chains)`, `getQueue()`, `clearQueue()`, `triggerBuild()`, `buildForce()` | chains は配列のまま送る。queue clear は body を送らない。 | active / waiting 区分、queue priority / created_seq、dispatch、chain config を API 順序のまま返す。 | queue full `429`、conflict `409`、validation `422` を保持する。 | active を waiting へ混在、priority 並び替え、queue 重複排除、chain DAG 検証を SDK が行わない。 |
| [`docs/details/runner.md` 詳細本文責務 §27.36](runner.md#sec-27-36) / [`docs/details/runner.md` 詳細本文責務 §27.37](runner.md#sec-27-37) failure category / environment | `getHistory()`, `getHistoryLog(id)` | `getHistory()` は `failureCategory != null` の場合だけ `failure_category` query を送る。 | `getHistory()` は history の `failure_category` と `HistoryPageObject.warnings`、`getHistoryLog(id)` は build log の `failure_category`、`failure_evidence`、`environment` を API response のまま返す。 | 未知 filter `422` と log 不在 `404` を保持する。 | category 分類、warnings 生成、environment fallback、secret mask 判定を SDK が行わない。 |
| [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) security | `getTokens()`, `createToken()`, `revokeToken()`, `getAuditLog()`, `getSessions()`, `revokeAllSessions()`, `getTotpStatus()`, `setupTotp()`, `confirmTotp(code)`, `disableTotp(code)`, `getApiRateLimit()`, `setApiRateLimit(policy)` | token / TOTP / rate limit request は [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) 引数変換契約どおり送る。 | token 本体、TOTP secret、otpauth URI は response として返すだけで SDK 内部に保存しない。 | `401` は token 破棄、`403` は token 維持、`429` は自動待機なし、`500` は message 保持。 | scope 推測、rate limit 待機、token list への token 合成、TOTP code 再送、audit 補完を行わない。 |

<a id="sec-27-21-2"></a>
**[§27.21〜§27.47 SDK 合格ゲート](sdk.md#sec-27-21-2)：**

| ゲート | 合格条件 |
|--------|----------|
| method coverage | [`docs/details/sdk.md` 詳細本文責務 §27.21](sdk.md#sec-27-21)〜[§27.47](sdk.md#sec-27-47) の api / sdk / ui 連動参照表で SDK method が記載された項目について、該当 public method が存在する。 |
| request exactness | fake fetch fixture で method、path、query key 順、body key 順、body なし endpoint が [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) と一致する。 |
| response passthrough | 成功 response に存在する追加 key を削除せず、存在しない key を追加しない。 |
| security handling | token、password、PAT、Webhook secret、SMTP password、TOTP secret、ticket、Authorization header を SDK property、console、error message に保存しない。 |
| error stability | `401` / `403` / `409` / `422` / `429` / `500` / network / timeout が固定 `AdlaireCIError` になり、自動 retry、自動 refresh、自動 logout は仕様に記載された場合だけ行う。 |
| binary / stream | snapshot download は `Blob`、SSE は `StreamHandle` とし、JSON response と混同しない。 |
| fixture evidence | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の SDK / UI 関連 fixture で、request shape、error shape、secret leak、token mutation、no retry、no response補完が確認される。 |

<a id="sec-27-21-3"></a>
**[§27.21〜§27.47 SDK 連動 fixture 証跡参照](sdk.md#sec-27-21-3)：**

SDK 連動 fixture 名、入力、expected、合格条件、禁止条件、実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) API / SDK / UI 連動 fixture 固定契約を正本とする。[`docs/details/sdk.md`](sdk.md) 詳細本文責務では、SDK method、request 生成、response passthrough、error object、token mutation、secret leak、side-effect 境界の実装契約だけを扱う。

**SDK response schema 参照契約：**

[`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) の public method 一覧に記載する `StatusObject`、`HistoryPageObject`、その他の型名は、API response を識別するための名称であり、response schema の正本ではない。HTTP response のキー、型、必須性、nullable、配列、許容値は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の Response 列と対象 endpoint 契約を正本とし、状態由来 record は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の schema を参照する。

SDK は成功 response を補完、削除、rename、既定値 merge、再集計せず、そのまま返す。`AdlaireCIError` と `StreamHandle` だけは API response ではなく SDK が生成する型であるため、[`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) の SDK 共通実装契約を正本とする。API response の required / nullable / array key を SDK 文書へ別表として再掲してはならない。

---
