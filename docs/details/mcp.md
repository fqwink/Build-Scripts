# MCP 詳細仕様

本書は `mcp` owner component の詳細仕様正本である。

本書は [docs/SPEC.md](../SPEC.md)、[docs/ROADMAP.md](../ROADMAP.md)、[docs/DETAIL_INDEX.md](../DETAIL_INDEX.md)、[docs/details/fixture.md](fixture.md)、[docs/DOCUMENT_INDEX.md](../DOCUMENT_INDEX.md) を参照する。

本書は方針、ポリシー、状態語彙、現在状態、Phase、将来計画を再定義しない。

<a id="sec-29-0"></a>
## 29.0 実装境界

| 項目 | 仕様 |
|------|------|
| owner component | `mcp` |
| 実装対象ファイル | `components/mcp.go` |
| 起動入口 | `main.go` の command dispatch |
| 配布 binary | `adlaire-ci-mcp` |
| 通信方式 | loopback HTTP JSON-RPC 2.0 |
| 既定 listen address | `127.0.0.1:8766` |
| 外部公開 | 既定禁止 |
| 永続化 | `statefile` owner の固定 state path を使用 |
| 監査ログ | `components/mcp.go` が `statefile` owner の固定 state path へ追記 |
| API 連携 | `api` owner の endpoint 契約を再利用 |
| runner 連携 | `runner` owner の queue / build / history 契約を再利用 |

`components/mcp.go` は MCP 専用の protocol adapter として動作する。

`components/mcp.go` は build、deploy、archive、commit status、API、SDK、UI、statefile、security の詳細仕様を再定義してはならない。

`components/mcp.go` が他 owner の機能を操作する場合は、該当 owner の公開済み関数または state 契約だけを使用する。

<a id="sec-29-1"></a>
## 29.1 起動 CLI

`adlaire-ci-mcp` は以下の CLI を持つ。

```text
adlaire-ci-mcp --state-dir <path> [--addr <host:port>] [--read-only] [--client-token <token>]
adlaire-ci-mcp --help
adlaire-ci-mcp --version
```

| option | 必須 | 既定値 | 仕様 |
|--------|------|--------|------|
| `--state-dir` | yes | none | statefile root。存在しない場合は起動失敗 |
| `--addr` | no | `127.0.0.1:8766` | bind address。loopback 以外は `--allow-non-loopback` がない限り起動失敗 |
| `--read-only` | no | `false` | 副作用 tool を登録しない |
| `--client-token` | no | none | 指定時は `Authorization: Bearer <token>` を必須にする |
| `--allow-non-loopback` | no | `false` | 明示指定時のみ loopback 以外 bind を許可 |

`--help` は終了 code `0` で usage を stdout へ出力する。

`--version` は終了 code `0` で binary name、version、commit、build time を stdout へ出力する。

起動失敗時は stderr に `error: <reason>` を出力し、終了 code `2` とする。

<a id="sec-29-2"></a>
## 29.2 HTTP endpoint

| method | path | 用途 | 成功 status | 失敗 status |
|--------|------|------|-------------|-------------|
| `POST` | `/mcp` | JSON-RPC 2.0 request | `200` | `400`, `401`, `403`, `422`, `500` |
| `GET` | `/mcp/events` | MCP notification SSE | `200` | `401`, `403`, `500` |
| `GET` | `/health` | health check | `200` | `500` |

`/mcp` request は `Content-Type: application/json` を必須とする。

`/mcp/events` response は `Content-Type: text/event-stream` を返す。

`/health` response body は以下とする。

```json
{
  "status": "ok",
  "component": "mcp"
}
```

<a id="sec-29-3"></a>
## 29.3 JSON-RPC 共通契約

request object は以下とする。

```json
{
  "jsonrpc": "2.0",
  "id": "string-or-number",
  "method": "string",
  "params": {}
}
```

success response object は以下とする。

```json
{
  "jsonrpc": "2.0",
  "id": "string-or-number",
  "result": {}
}
```

error response object は以下とする。

```json
{
  "jsonrpc": "2.0",
  "id": "string-or-number-or-null",
  "error": {
    "code": -32600,
    "message": "Invalid Request",
    "data": {
      "reason": "string"
    }
  }
}
```

| code | message | 条件 |
|------|---------|------|
| `-32700` | `Parse error` | JSON parse 失敗 |
| `-32600` | `Invalid Request` | JSON-RPC 形式不正 |
| `-32601` | `Method not found` | 未登録 method |
| `-32602` | `Invalid params` | params schema 不一致 |
| `-32603` | `Internal error` | 想定外エラー |
| `-32001` | `Unauthorized` | token 不一致 |
| `-32002` | `Forbidden` | scope 不足または read-only |
| `-32003` | `Timeout` | tool timeout 超過 |
| `-32004` | `Confirmation required` | elicitation confirmation 未完了 |

**JSON-RPC method 固定表：**

| method | params | result | 副作用 |
|--------|--------|--------|--------|
| `initialize` | [`docs/details/mcp.md` 詳細本文責務 §29.4](mcp.md#sec-29-4) | serverInfo / capabilities | `.mcp_client_log` 追記 |
| `notifications/initialized` | body なし | response なし | connection を initialized 済みにする |
| `tools/list` | `{}` | `{ "tools": ToolDescriptor[] }` | なし |
| `tools/call` | `{ "name": string, "arguments": object }` | [`docs/details/mcp.md` 詳細本文責務 §29.6](mcp.md#sec-29-6) の wrapper | tool ごとの契約に従う |
| `resources/list` | `{}` | `{ "resources": ResourceDescriptor[] }` | なし |
| `resources/read` | `{ "uri": string }` | `{ "contents": ResourceContent[] }` | なし |
| `resources/subscribe` | `{ "uri": string }` | `{}` | connection 内 subscription 追加 |
| `resources/unsubscribe` | `{ "uri": string }` | `{}` | connection 内 subscription 削除 |
| `prompts/list` | `{}` | `{ "prompts": PromptDescriptor[] }` | なし |
| `prompts/get` | `{ "name": string, "arguments": object }` | `{ "messages": PromptMessage[] }` | なし |

`initialize` 成功前に `initialize` 以外の method を受けた場合は `-32600 Invalid Request` とする。`notifications/initialized` は JSON-RPC notification とし、`id` を持つ request として送られた場合は `-32600 Invalid Request` とする。batch request は受け付けず `-32600 Invalid Request` とする。

<a id="sec-29-4"></a>
## 29.4 MCP initialize

`initialize` params は以下とする。

```json
{
  "protocolVersion": "2025-06-18",
  "clientInfo": {
    "name": "string",
    "version": "string"
  },
  "capabilities": {}
}
```

`initialize` result は以下とする。

```json
{
  "protocolVersion": "2025-06-18",
  "serverInfo": {
    "name": "adlaire-ci-mcp",
    "version": "string"
  },
  "capabilities": {
    "tools": {},
    "resources": {
      "subscribe": true
    },
    "prompts": {},
    "logging": {},
    "sampling": {},
    "elicitation": {}
  }
}
```

`initialize` 成功時は client name、client version、remote address、capabilities、connected_at を `.mcp_client_log` へ追記する。

<a id="sec-29-5"></a>
## 29.5 Tool 一覧

| tool name | scope | 副作用 | owner 参照 |
|-----------|-------|--------|------------|
| `adlaire.getStatus` | `read:status` | no | `api`, `runner` |
| `adlaire.getQueue` | `read:queue` | no | `runner` |
| `adlaire.triggerBuild` | `write:build` | yes | `runner` |
| `adlaire.cancelQueueEntry` | `write:queue` | yes | `runner` |
| `adlaire.getHistory` | `read:history` | no | `runner`, `archive` |
| `adlaire.getBuildLog` | `read:history` | no | `runner` |
| `adlaire.analyzeBuildError` | `read:history` | no | `mcp` |
| `adlaire.getConfig` | `read:config` | no | `api`, `statefile` |
| `adlaire.setConfig` | `write:config` | yes | `api`, `statefile` |
| `adlaire.getMcpConfig` | `read:mcp` | no | `mcp`, `statefile` |
| `adlaire.setMcpConfig` | `write:mcp` | yes | `mcp`, `statefile` |
| `adlaire.createConfigSnapshot` | `write:config` | yes | `statefile` |
| `adlaire.diffConfigSnapshots` | `read:config` | no | `statefile` |
| `adlaire.restoreConfigSnapshot` | `write:config` | yes | `statefile` |
| `adlaire.getMetrics` | `read:metrics` | no | `api` |
| `adlaire.getAuditLog` | `read:audit` | no | `statefile`, `security` |
| `adlaire.resendWebhook` | `write:notification` | yes | `runner` |

`--read-only` 指定時は副作用 `yes` の tool を `tools/list` に含めない。

副作用 `yes` の tool は [`docs/details/mcp.md` 詳細本文責務 §29.15](mcp.md#sec-29-15) の elicitation を必須とする。

<a id="sec-29-6"></a>
<a id="296-tool-schema"></a>
## 29.6 Tool schema

`tools/list` の `ToolDescriptor` は `name`、`description`、`inputSchema`、`annotations` を持つ。`annotations.readOnlyHint` は副作用 `no` の tool だけ `true`、副作用 `yes` の tool は `false` とする。`inputSchema` は JSON Schema draft 非依存の object とし、`type`、`required`、`properties`、`additionalProperties:false` だけを使用する。

| tool name | params | result `data` | 固定条件 |
|-----------|--------|---------------|----------|
| `adlaire.getStatus` | `{}` | `GET /api/status` 相当の object | 状態を変更しない。 |
| `adlaire.getQueue` | `{}` | `GET /api/queue` 相当の object | 状態を変更しない。 |
| `adlaire.triggerBuild` | `target`, `source`, `options` | queue id と dispatch | elicitation 必須。`POST /api/build` 境界を使用する。 |
| `adlaire.cancelQueueEntry` | `queue_id` | cancelled queue id | elicitation 必須。running entry は error。 |
| `adlaire.getHistory` | `limit`, `offset` | history list | `limit` 1〜100、`offset` 0 以上。 |
| `adlaire.getBuildLog` | `build_id` | log object | 存在しない build は `-32602`。 |
| `adlaire.analyzeBuildError` | `build_id` | sampling summary | 外部 AI API を server から直接呼ばない。 |
| `adlaire.getConfig` | `{}` | effective Adlaire config | `.mcp_config` を返さない。 |
| `adlaire.setConfig` | `path`, `value` | updated path | elicitation 必須。`.server_config` の許可 key path だけ。 |
| `adlaire.getMcpConfig` | `{}` | `.mcp_config` | token 本体を返さない。 |
| `adlaire.setMcpConfig` | `tool_timeout_ms`, `sampling_timeout_ms`, `scopes` | updated MCP config | elicitation 必須。token 本体は禁止。 |
| `adlaire.createConfigSnapshot` | `label` | snapshot id | elicitation 必須。 |
| `adlaire.diffConfigSnapshots` | `left_id`, `right_id` | diff object | secret は `"***"`。 |
| `adlaire.restoreConfigSnapshot` | `snapshot_id` | restored paths | elicitation 必須。 |
| `adlaire.getMetrics` | `{}` | metrics object | `.mcp_metrics` と API metrics を読む。 |
| `adlaire.getAuditLog` | `limit`, `offset` | audit list | MCP audit tail と API audit を混在させない。 |
| `adlaire.resendWebhook` | `delivery_id` | resend accepted | elicitation 必須。 |

tool 別 params schema は以下に固定する。表にない key は `-32602 Invalid params` とし、tool 実行、audit、metrics 更新を行わない。

| tool name | required | validation |
|-----------|----------|------------|
| `adlaire.getStatus` | なし | params は `{}`。 |
| `adlaire.getQueue` | なし | params は `{}`。 |
| `adlaire.triggerBuild` | `target`, `source` | `target` と `source` は 1〜128 byte UTF-8。`options` は object、省略時 `{}`。 |
| `adlaire.cancelQueueEntry` | `queue_id` | `queue_id` は 1〜128 byte UTF-8、slash、NUL、CR、LF 禁止。 |
| `adlaire.getHistory` | なし | `limit` は 1〜100、省略時 50。`offset` は 0 以上、省略時 0。 |
| `adlaire.getBuildLog` | `build_id` | `build_id` は 1〜128 byte UTF-8、slash、NUL、CR、LF 禁止。 |
| `adlaire.analyzeBuildError` | `build_id` | `build_id` は `adlaire.getBuildLog` と同じ。 |
| `adlaire.getConfig` | なし | params は `{}`。 |
| `adlaire.setConfig` | `path`, `value` | `path` は `.server_config` の許可 key path。`value` は JSON value。 |
| `adlaire.getMcpConfig` | なし | params は `{}`。 |
| `adlaire.setMcpConfig` | なし | `tool_timeout_ms`、`sampling_timeout_ms`、`scopes` の 1 key 以上。各 key は [`docs/details/statefile.md` 詳細本文責務 §22.0d](statefile.md#sec-22-0d) の `McpConfig` に従う。 |
| `adlaire.createConfigSnapshot` | なし | `label` は string または `null`、省略時 `null`。 |
| `adlaire.diffConfigSnapshots` | `left_id`, `right_id` | id は `cfgsnap_` prefix 必須。同一 id は `-32602`。 |
| `adlaire.restoreConfigSnapshot` | `snapshot_id` | `snapshot_id` は `cfgsnap_` prefix 必須。 |
| `adlaire.getMetrics` | なし | params は `{}`。 |
| `adlaire.getAuditLog` | なし | `limit` は 1〜100、省略時 50。`offset` は 0 以上、省略時 0。 |
| `adlaire.resendWebhook` | `delivery_id` | `delivery_id` は 1〜128 byte UTF-8、slash、NUL、CR、LF 禁止。 |

`adlaire.triggerBuild` params は以下とする。

```json
{
  "target": "string",
  "source": "string",
  "options": {}
}
```

`target` は `builder` 詳細仕様に定義された build target と一致する。

`source` は local path または repository ref を表す文字列とする。

`adlaire.cancelQueueEntry` params は以下とする。

```json
{
  "queue_id": "string"
}
```

`adlaire.setConfig` params は以下とする。

```json
{
  "path": "string",
  "value": {}
}
```

`path` は `.server_config` 内の許可 key path のみ指定できる。

`adlaire.restoreConfigSnapshot` params は以下とする。

```json
{
  "snapshot_id": "string"
}
```

`adlaire.resendWebhook` params は以下とする。

```json
{
  "delivery_id": "string"
}
```

tool result は必ず以下の wrapper を返す。

```json
{
  "ok": true,
  "data": {},
  "warnings": []
}
```

失敗は JSON-RPC error とし、`ok:false` result は返さない。

tool params に未知 key がある場合、必須 key 不足、型不一致、範囲外、secret 値を許可しない field への secret 形状入力は `-32602 Invalid params` とし、tool 実行、audit、metrics 更新を行わない。tool 実行開始後の失敗は `.mcp_audit_log` と `.mcp_metrics` に失敗結果を記録する。ただし audit 追記不能時は副作用 tool を失敗扱いにし、対象 owner の状態変更を開始しない。

<a id="sec-29-7"></a>
## 29.7 Resources

| resource URI | 内容 | 更新通知 |
|--------------|------|----------|
| `adlaire://status` | 現在 status | yes |
| `adlaire://queue` | build queue | yes |
| `adlaire://history` | build history summary | yes |
| `adlaire://history/{build_id}` | build detail | no |
| `adlaire://logs/{build_id}` | build log | no |
| `adlaire://config` | effective config | yes |
| `adlaire://metrics` | tool and build metrics | yes |
| `adlaire://audit` | MCP audit log tail | yes |

`resources/list` は上表の URI と name、description、mimeType を返す。

`resources/read` は存在しない URI に `-32602 Invalid params` を返す。

`resources/read` の `ResourceContent` は `uri`、`mimeType`、`text` を持つ。`adlaire://logs/{build_id}` だけは `mimeType:"text/plain"`、その他は `mimeType:"application/json"` とする。JSON resource の `text` は UTF-8 JSON object 文字列とし、secret、token、Authorization header を含めない。

<a id="sec-29-8"></a>
## 29.8 Resource subscription

`resources/subscribe` params は以下とする。

```json
{
  "uri": "adlaire://status"
}
```

subscription は client connection 単位で保持する。

subscription は `.mcp_subscriptions` へ永続化しない。

`resources/unsubscribe` は未購読 URI に対しても `{}` を返す。存在しない URI は `-32602 Invalid params` とする。

resource 更新時は `/mcp/events` へ以下を送信する。

```text
event: resource-updated
data: {"uri":"adlaire://status","updated_at":"2026-09-28T00:00:00Z"}
```

<a id="sec-29-9"></a>
## 29.9 Prompts

| prompt name | params | 用途 |
|-------------|--------|------|
| `adlaire.buildFailureTriage` | `{ "build_id": "string" }` | build 失敗原因の確認手順生成 |
| `adlaire.releaseReadiness` | `{ "target": "string" }` | release 前確認項目生成 |
| `adlaire.configReview` | `{ "snapshot_id": "string" }` | 設定 snapshot の差分確認 |

`prompts/get` result は message list を返す。

`prompts/list` は `name`、`description`、`arguments` を返す。`prompts/get` は `arguments` を prompt ごとの params と照合し、必須 key 不足、未知 key、型不一致を `-32602 Invalid params` とする。`PromptMessage` は `role:"user"` と `content.type:"text"` だけを使用する。prompt text は実行を促す手順文だけを生成し、server 側で build、config、queue を変更しない。

prompt は実行を伴わない。

<a id="sec-29-10"></a>
## 29.10 Sampling

`adlaire.analyzeBuildError` は sampling 対応 tool とする。

server は外部 AI API を直接呼び出してはならない。

server は MCP client へ `sampling/createMessage` request を送信する。

sampling request timeout は `.mcp_config.sampling_timeout_ms` を使用する。

timeout 時は `-32003 Timeout` を返す。

sampling result は `.mcp_audit_log` へ prompt hash、build_id、client name、duration_ms、status だけを記録する。

sampling result の本文は statefile へ保存しない。

<a id="sec-29-11"></a>
## 29.11 Notifications

MCP server は以下の notification を送信する。

| notification | 条件 |
|--------------|------|
| `notifications/tools/list_changed` | tool scope または read-only 状態が変わった |
| `notifications/resources/list_changed` | resource 一覧が変わった |
| `notifications/prompts/list_changed` | prompt 一覧が変わった |
| `notifications/message` | build、queue、config snapshot、audit のイベントが発生した |

notification は `/mcp/events` の SSE で送信する。

SSE client が未接続の場合、notification は破棄する。

notification の破棄はエラーとして扱わない。

<a id="sec-29-12"></a>
## 29.12 HTTP SSE transport

`/mcp/events` は keepalive として 30 秒ごとに comment frame を送信する。

```text
: keepalive

```

client disconnect は正常終了として扱う。

server shutdown 時は以下を送信して connection を閉じる。

```text
event: shutdown
data: {"reason":"server_shutdown"}
```

<a id="sec-29-13"></a>
## 29.13 Tool scope

scope は `.mcp_config.scopes` に定義する。

token 未指定起動時は local trusted mode とし、loopback 接続に限り全 scope を許可する。

token 指定起動時は token record に紐づく scope だけを許可する。

scope 不足時は `-32002 Forbidden` を返す。

`write:*` scope は `read:*` scope を暗黙に含めない。

<a id="sec-29-14"></a>
## 29.14 Audit / client / metrics

副作用 tool 実行時は `.mcp_audit_log` へ以下を追記する。

```json
{
  "id": "string",
  "timestamp": "2026-09-28T00:00:00Z",
  "client_name": "string",
  "tool": "adlaire.triggerBuild",
  "params_hash": "sha256",
  "scope": "write:build",
  "confirmation_id": "string",
  "status": "success",
  "duration_ms": 123
}
```

全 tool 実行は `.mcp_metrics` へ aggregate する。

metrics key は `tool_name`、`status`、`count`、`last_duration_ms`、`last_at` とする。

client 接続は `.mcp_client_log` へ append-only で記録する。

<a id="sec-29-15"></a>
## 29.15 Timeout / config CRUD / elicitation

`.mcp_config` schema は [`docs/details/statefile.md` 詳細本文責務 §22.0d](statefile.md#sec-22-0d) を正本とする。

tool timeout は `.mcp_config.tool_timeout_ms` を使用する。

`adlaire.getMcpConfig` は `.mcp_config` 全体を返す。

`adlaire.setMcpConfig` は `.mcp_config` の許可 key のみ更新する。

副作用 tool は以下の elicitation flow を必須とする。

1. 初回 request で `confirmation_id` がない場合、`-32004 Confirmation required` を返し、`data.confirmation_id` と `data.summary` を含める。
2. client は同じ params に `confirmation_id` を付けて再要求する。
3. server は params hash、confirmation_id、有効期限、tool name が一致する場合だけ実行する。
4. confirmation 有効期限は 5 分とする。

confirmation_id は memory only とし、statefile へ保存しない。

<a id="sec-29-16"></a>
## 29.16 仕様化済み対象

本書は以下を仕様化済み・未実装の詳細仕様として定義する。

| 機能 | 詳細仕様 |
|------|----------|
| MCP サーバー実装 | [`docs/details/mcp.md` 詳細本文責務 §29.0](mcp.md#sec-29-0) から [§29.4](mcp.md#sec-29-4) |
| MCP ツール・リソース公開 | [`docs/details/mcp.md` 詳細本文責務 §29.5](mcp.md#sec-29-5) から [§29.8](mcp.md#sec-29-8) |
| AI 支援ビルドエラー分析 | [`docs/details/mcp.md` 詳細本文責務 §29.10](mcp.md#sec-29-10) |
| MCP Prompts 定義 | [`docs/details/mcp.md` 詳細本文責務 §29.9](mcp.md#sec-29-9) |
| MCP Sampling によるビルドログ自動分析 | [`docs/details/mcp.md` 詳細本文責務 §29.10](mcp.md#sec-29-10) |
| MCP Notifications | [`docs/details/mcp.md` 詳細本文責務 §29.11](mcp.md#sec-29-11) |
| MCP HTTP SSE transport 対応 | [`docs/details/mcp.md` 詳細本文責務 §29.12](mcp.md#sec-29-12) |
| MCP ツールスコープ細分化 | [`docs/details/mcp.md` 詳細本文責務 §29.13](mcp.md#sec-29-13) |
| MCP ツール呼び出し監査ログ | [`docs/details/mcp.md` 詳細本文責務 §29.14](mcp.md#sec-29-14) |
| MCP リソース購読 | [`docs/details/mcp.md` 詳細本文責務 §29.8](mcp.md#sec-29-8) |
| MCP クライアント情報ログ | [`docs/details/mcp.md` 詳細本文責務 §29.14](mcp.md#sec-29-14) |
| MCP ツール実行統計 | [`docs/details/mcp.md` 詳細本文責務 §29.14](mcp.md#sec-29-14) |
| MCP ツール実行タイムアウト設定 | [`docs/details/mcp.md` 詳細本文責務 §29.15](mcp.md#sec-29-15) |
| MCP 設定 CRUD ツール | [`docs/details/mcp.md` 詳細本文責務 §29.15](mcp.md#sec-29-15) |
| MCP Elicitation による副作用操作の確認 | [`docs/details/mcp.md` 詳細本文責務 §29.15](mcp.md#sec-29-15) |

<a id="sec-29-17"></a>
## 29.17 検証条件

実装完了判定には以下を必須とする。

- `adlaire-ci-mcp --help` が終了 code `0` で usage を出力する。
- `adlaire-ci-mcp --version` が終了 code `0` で version を出力する。
- `adlaire-ci-mcp --state-dir <fixture>` が loopback で起動する。
- `/health` が `{"status":"ok","component":"mcp"}` を返す。
- `initialize` が [`docs/details/mcp.md` 詳細本文責務 §29.4](mcp.md#sec-29-4) の capabilities を返す。
- `tools/list` が read-only 指定の有無に応じて副作用 tool の有無を切り替える。
- `resources/list` が [`docs/details/mcp.md` 詳細本文責務 §29.7](mcp.md#sec-29-7) の URI を返す。
- `prompts/list` が [`docs/details/mcp.md` 詳細本文責務 §29.9](mcp.md#sec-29-9) の prompt を返す。
- 副作用 tool が confirmation なしで実行されない。
- confirmation 付き副作用 tool が `.mcp_audit_log` を追記する。
- tool timeout 超過が `-32003 Timeout` を返す。
- SSE が keepalive、resource update、shutdown を送信する。
- sampling tool が server から外部 AI API を直接呼び出さない。
- `.mcp_client_log`、`.mcp_metrics`、`.mcp_audit_log` の fixture assertion が存在する。
