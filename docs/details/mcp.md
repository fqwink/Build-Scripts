# Adlaire CI — MCP 詳細仕様

[`docs/details/mcp.md`](mcp.md) 詳細本文責務は、`mcp` owner component の詳細仕様正本である。

本書は [`docs/SPEC.md` 方針責務・ポリシー責務](../SPEC.md)、[`docs/ROADMAP.md` 状態・計画責務](../ROADMAP.md)、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務](../DETAIL_INDEX.md)、[`docs/details/fixture.md` fixture 証跡責務](fixture.md)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務](../DOCUMENT_INDEX.md) を参照する。

本書は方針、ポリシー、状態語彙、現在状態、Phase、将来計画を再定義しない。

<a id="sec-29-0"></a>
**29.0 実装境界：**

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
| 検証接続 | 検証方針と完了可否は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、fixture / expected / fake / mutation test 証跡 / test・contract drift 証跡の記録先は [`docs/details/fixture.md` fixture 証跡責務](fixture.md)、[mutation test 証跡固定契約](fixture.md#mutation-test-evidence-contract)、[test / contract drift 証跡固定契約](fixture.md#test-contract-drift-evidence-contract) を参照する。 |

`components/mcp.go` は MCP 専用の protocol adapter として動作する。

`components/mcp.go` は build、deploy、archive、commit status、API、SDK、UI、statefile、security の詳細仕様を再定義してはならない。

`components/mcp.go` が他 owner の機能を操作する場合は、該当 owner の公開済み関数または state 契約だけを使用する。

<a id="sec-29-1"></a>
**29.1 起動 CLI：**

`adlaire-ci-mcp` は以下の CLI を持つ。

```text
adlaire-ci-mcp --state-dir <path> [--addr <host:port>] [--read-only] [--client-token <token>] [--allow-non-loopback]
adlaire-ci-mcp --help
adlaire-ci-mcp --version
```

| option | 必須 | 既定値 | 仕様 |
|--------|------|--------|------|
| `--state-dir <path>` | yes | none | statefile root。空でない絶対 path、既存 directory、symlink でないことを必須とする。 |
| `--addr <host:port>` | no | `127.0.0.1:8766` | host は IPv4 literal または exact `localhost` だけを許可する。port は `1`〜`65535` の 10 進数。 |
| `--read-only` | no | `false` | 副作用 tool を `tools/list` から除外し、既存 connection の副作用 tool call を `-32002 Forbidden` にする。 |
| `--client-token <token>` | no | none | 指定時は `/mcp` と `/mcp/events` に `Authorization: Bearer <token>` を必須にする。`/health` では要求しない。 |
| `--allow-non-loopback` | no | `false` | 明示指定時のみ `--addr` の non-loopback IPv4 unicast host を許可する。指定がない場合、`127.0.0.0/8` と exact `localhost` 以外は拒否する。 |
| `--help` | no | none | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 CLI 共通固定契約](../DETAIL_INDEX.md#common-cli-contract) に従う。 |
| `--version` | no | none | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 CLI 共通固定契約](../DETAIL_INDEX.md#common-cli-contract) に従う。 |

CLI parse は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 CLI 共通固定契約](../DETAIL_INDEX.md#common-cli-contract) を使用する。短縮 option、未定義 option、未許可の `--name=value`、未定義位置引数を禁止する。`--state-dir`、`--addr`、`--client-token`、`--read-only`、`--allow-non-loopback` は同一 option の重複を禁止する。値 option と boolean option のどちらも、2 回目以降の出現を parse error とし、最後の値採用、重複 boolean の黙認、重複値の merge を行ってはならない。

`--addr` の host 判定は文字列 parse 後、名前解決を行わずに実施する。exact `localhost` は loopback として扱い、その他の host 名、IPv6 literal、空 host、wildcard、zone identifier を禁止する。IPv4 literal は 4 octet dotted decimal だけを許可し、octet の空文字、符号、16 進表記、8 進表記、余分な空白、先頭 `+`、NUL、CR、LF を禁止する。`--allow-non-loopback` の有無にかかわらず、`0.0.0.0`、`255.255.255.255`、`224.0.0.0/4`、`240.0.0.0/4` は bind 対象として禁止する。`--allow-non-loopback` なしの場合は `127.0.0.0/8` と exact `localhost` だけを許可する。`--allow-non-loopback` ありの場合は、上記の禁止 host を除く IPv4 unicast literal と exact `localhost` を許可する。

| 条件 | stdout | stderr | 終了 code | 副作用 |
|------|--------|--------|-----------|--------|
| `--help` | `Usage: adlaire-ci-mcp --state-dir path [--addr host:port] [--read-only] [--client-token token] [--allow-non-loopback] [--version] [--help]` + LF | 空 | `0` | 状態、listener、client log、metrics、audit に触れない。 |
| `--version` | `adlaire-ci-mcp <binary-version> go=<runtime.Version()>` + LF | 空 | `0` | 同上。 |
| option 重複 | 空 | `duplicate option: <option>` + LF | `2` | 状態 read/write、listener、client log、metrics、audit を開始しない。 |
| `--state-dir` 未指定 | 空 | `state directory is required` + LF | `2` | listener を起動しない。 |
| `--state-dir` が symlink | 空 | `state directory must not be symlink: <path>` + LF | `2` | 同上。 |
| `--addr` 形式不正 | 空 | `invalid listen address: <address>` + LF | `2` | 状態 read/write を開始しない。 |
| non-loopback かつ `--allow-non-loopback` なし | 空 | `non-loopback address is not allowed: <address>` + LF | `2` | 同上。 |
| `--client-token` が空文字 | 空 | `client token must not be empty` + LF | `2` | 同上。 |
| listener 起動失敗 | 空 | `listen failed` + LF | `1` | statefile を変更しない。 |

CLI 検証順は、共通 option mode 確定、argv token safety、option parse、option 重複、`--state-dir` 未指定、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 CLI state directory 共通固定契約](../DETAIL_INDEX.md#common-state-dir-contract) の空文字・相対 path・不在・directory 判定、symlink 判定、`--addr` 形式、禁止 host 判定、loopback 判定、`--client-token` 検証、listener 起動の順に固定する。

<a id="sec-29-2"></a>
**29.2 HTTP endpoint：**

| method | path | 用途 | 成功 status | 失敗 status |
|--------|------|------|-------------|-------------|
| `POST` | `/mcp` | JSON-RPC 2.0 request | `200` | `400`, `401`, `403`, `413`, `422`, `500` |
| `GET` | `/mcp/events` | MCP notification SSE | `200` | `401`, `403`, `500` |
| `GET` | `/health` | health check | `200` | `500` |

`/mcp` request は `Content-Type: application/json` または `application/json; charset=utf-8` を必須とする。未指定、不一致、複数値は `400` とし、JSON-RPC error object ではなく `{"error":"Invalid content type"}` を返す。

`/mcp` request body の上限は 1 MiB とする。超過時は `413` と `{"error":"Payload too large"}` を返し、JSON parse、tool 実行、audit、metrics 更新を行わない。

`/mcp/events` response は `Content-Type: text/event-stream` を返す。

`/mcp/events` は `Cache-Control: no-store` と `X-Accel-Buffering: no` を返す。SSE 接続確立前の認証失敗は JSON `{"error":"Unauthorized"}` または `{"error":"Forbidden"}` を返し、SSE frame を送信しない。

`/mcp` と `/mcp/events` の token 認証は `Authorization` header だけを使用する。header 値は exact `Bearer ` prefix と token 文字列 1 個で構成し、前後空白、複数 header、空 token、token 後続空白、query parameter、Cookie、Basic 認証をすべて拒否する。`--client-token` 未指定時でも non-loopback 接続は起動時に拒否済みでなければならず、request handler で外部公開を補完許可してはならない。

`POST /mcp` 以外の method で `/mcp` を呼んだ場合、request body を読まず `405` と `{"error":"Method not allowed"}` を返す。`GET /mcp/events` 以外の method で `/mcp/events` を呼んだ場合、SSE 接続を確立せず `405` と `{"error":"Method not allowed"}` を返す。`/mcp/events` は request body を読まない。

HTTP response の JSON body は UTF-8、LF なしの compact JSON とする。成功時、JSON-RPC error 時、HTTP error 時のいずれも token、Authorization header、raw params の secret 値、state-dir absolute path、Go error 文字列を response body へ含めてはならない。

`/health` response body は以下とする。

```json
{
  "status": "ok",
  "component": "mcp"
}
```

`/health` は状態 file を読まず、`.mcp_client_log`、`.mcp_audit_log`、`.mcp_metrics` を更新しない。method 不一致は `405` と `{"error":"Method not allowed"}`、未知 path は `404` と `{"error":"Not found"}` を返す。

<a id="sec-29-3"></a>
**29.3 JSON-RPC 共通契約：**

request object は以下とする。

```json
{
  "jsonrpc": "2.0",
  "id": "string-or-number-or-null",
  "method": "string",
  "params": {}
}
```

`id` は string、number、または `null` だけを許可する。number は JSON number とし、NaN、Infinity は JSON として受け付けない。object、array、boolean の `id` は `-32600 Invalid Request` とする。

`method` は 1〜128 byte の UTF-8 string とし、NUL、CR、LF を禁止する。`params` は object、または省略時 `{}` と同じ扱いに固定する。`params:null`、array、string、number、boolean は `-32602 Invalid params` とする。

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

JSON parse 失敗時の HTTP status は `400`、JSON-RPC error response の `id` は `null` とする。JSON-RPC 形式不正、method 不在、params 不正、初期化順序違反は HTTP status `200` で JSON-RPC error object を返す。認証失敗だけは JSON-RPC error object を返さず、HTTP status `401` または `403` と `{"error":"Unauthorized"}` / `{"error":"Forbidden"}` を返す。

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

`initialize` 成功前に `initialize` 以外の method を受けた場合は `-32600 Invalid Request` とする。`notifications/initialized` は JSON-RPC notification とし、`id` を持つ request として送られた場合は `-32600 Invalid Request` とする。notification 成功時は HTTP status `202`、body 空、状態副作用は connection 初期化状態の memory 更新だけとする。batch request は受け付けず `-32600 Invalid Request` とする。

JSON-RPC request object の top-level key は `jsonrpc`、`id`、`method`、`params` だけを許可する。未知 key、重複 key、`jsonrpc` が exact `"2.0"` でない値、`method` 欠落、`method` 空文字は `-32600 Invalid Request` とする。JSON object 内の重複 key は最初の key を採用せず、request 全体を不正として扱う。

JSON-RPC success response と error response の key order は `jsonrpc`、`id`、`result` または `error` の順に固定する。`error` object の key order は `code`、`message`、`data` の順とし、`data.reason` は secret、path、Go error を含まない固定分類文字列だけにする。notification に対する成功 response body は送信しない。

<a id="sec-29-4"></a>
**29.4 MCP initialize：**

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

`clientInfo.name` と `clientInfo.version` は 1〜128 byte UTF-8、NUL、CR、LF 禁止とする。不正時は `-32602 Invalid params` とし、`.mcp_client_log` を更新しない。

`initialize` 成功時は [`docs/details/statefile.md` 詳細本文責務 §22.0d](statefile.md#sec-22-0d) の `McpClientRecord` として client name、client version、protocol version、remote address、capabilities、connected_at を `.mcp_client_log` へ追記する。`.mcp_client_log` 追記失敗時は `-32603 Internal error` を返し、connection を initialized 済みにしない。

<a id="sec-29-5"></a>
**29.5 Tool 一覧：**

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

`tools/list` の tool order は [`docs/details/mcp.md` 詳細本文責務 §29.5](mcp.md#sec-29-5) の固定表の上から下の順とする。`--read-only` 指定時は副作用 `yes` の行だけを除外し、残る tool の相対順序を変更しない。scope 不足の tool は `tools/list` から除外せず、`tools/call` 時に `-32002 Forbidden` を返す。

<a id="sec-29-6"></a>
**29.6 Tool schema：**

`tools/list` の `ToolDescriptor` は `name`、`description`、`inputSchema`、`annotations` を持つ。これ以外の key を返してはならない。`name` は [`docs/details/mcp.md` 詳細本文責務 §29.5](mcp.md#sec-29-5) の tool name と完全一致させる。`description` は下記の `ToolDescriptor.description` 固定表の文字列だけを返す。`annotations` は `readOnlyHint` だけを持ち、副作用 `no` の tool は `true`、副作用 `yes` の tool は `false` とする。`inputSchema` は JSON Schema draft 非依存の object とし、root は `type:"object"`、`required`、`properties`、`additionalProperties:false` の 4 key だけを持つ。

`ToolDescriptor.inputSchema.required` は tool 別 params schema 表の `required` 列を、表記順の JSON array で返す。必須 key がない場合は空 array を返す。`properties` は当該 tool で許可される params key だけを持ち、未知 key を schema に含めてはならない。string 値は `{"type":"string"}`、integer 値は `{"type":"integer"}`、boolean 値は `{"type":"boolean"}`、object 値は `{"type":"object","required":[],"properties":{},"additionalProperties":false}`、全 JSON 型を受ける `value` だけは `{"type":["object","array","string","number","integer","boolean","null"]}` とする。`options` は予約 field とし、初期仕様では空 object だけを許可する。`inputSchema` は型、必須 key、未知 key 禁止だけを表し、byte 長、prefix、path、scope、confirmation は下表と本節本文の validation で判定する。

`ToolDescriptor.description` は以下に固定する。

| tool name | description |
|-----------|-------------|
| `adlaire.getStatus` | `Read current Adlaire CI status.` |
| `adlaire.getQueue` | `Read current build queue.` |
| `adlaire.triggerBuild` | `Request a new build through the durable queue.` |
| `adlaire.cancelQueueEntry` | `Cancel a waiting build queue entry.` |
| `adlaire.getHistory` | `Read build history summary.` |
| `adlaire.getBuildLog` | `Read one saved build log.` |
| `adlaire.analyzeBuildError` | `Request client-side sampling for a saved build error.` |
| `adlaire.getConfig` | `Read effective Adlaire CI configuration.` |
| `adlaire.setConfig` | `Update one allowed configuration path.` |
| `adlaire.getMcpConfig` | `Read MCP configuration without token secrets.` |
| `adlaire.setMcpConfig` | `Update MCP timeout or scope configuration.` |
| `adlaire.createConfigSnapshot` | `Create a configuration snapshot.` |
| `adlaire.diffConfigSnapshots` | `Diff two configuration snapshots with secrets masked.` |
| `adlaire.restoreConfigSnapshot` | `Restore one configuration snapshot.` |
| `adlaire.getMetrics` | `Read MCP and API metrics snapshot.` |
| `adlaire.getAuditLog` | `Read MCP audit log entries.` |
| `adlaire.resendWebhook` | `Request webhook delivery resend.` |

| tool name | params | result `data` | 固定条件 |
|-----------|--------|---------------|----------|
| `adlaire.getStatus` | `{}` | [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の `GET /api/status` response body object | 状態を変更しない。 |
| `adlaire.getQueue` | `{}` | [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の `GET /api/queue` response body object | 状態を変更しない。 |
| `adlaire.triggerBuild` | `target`, `source`, `options` | queue id と dispatch | elicitation 必須。`POST /api/build` 境界を使用する。`options` は初期仕様では空 object だけを許可する。 |
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
| `adlaire.getMetrics` | `{}` | metrics object | `.mcp_metrics` と API metrics を読む。返却 snapshot は現在の `adlaire.getMetrics` 呼び出し分を含めず、response data 確定後に現在呼び出しの metrics を更新する。 |
| `adlaire.getAuditLog` | `limit`, `offset` | audit list | MCP audit tail と API audit を混在させない。 |
| `adlaire.resendWebhook` | `delivery_id` | resend accepted | elicitation 必須。 |

tool 別 params schema は以下に固定する。表にない key は `-32602 Invalid params` とし、tool 実行、audit、metrics 更新を行わない。

| tool name | required | validation |
|-----------|----------|------------|
| `adlaire.getStatus` | なし | params は `{}`。 |
| `adlaire.getQueue` | なし | params は `{}`。 |
| `adlaire.triggerBuild` | `target`, `source` | `target` と `source` は 1〜128 byte UTF-8。`options` は object、省略時 `{}`。`options` に key がある場合は `-32602 Invalid params`。 |
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

`target` は `builder` 詳細仕様に定義された build target と一致する。空文字、`/` だけ、NUL、CR、LF、`..` segment を禁止する。初期仕様では `target` を queue entry の `target_id` として保存し、MCP 側で別名展開、default target 補完、複数 target 展開を行ってはならない。

`source` は `local:<absolute-path>` または `git:<ref>` のいずれかとする。`local:` の path は絶対 path、NUL / CR / LF 禁止、`..` segment 禁止とする。`git:` の ref は 1〜128 byte、空白、NUL、CR、LF、`..`、`@{` を禁止する。MCP 側は `source` を Git checkout、fetch、local file copy の実行根拠にしてはならず、API / runner 境界へ渡す queue payload の値としてだけ扱う。

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

`path` は `.server_config` 内の許可 key path のみ指定できる。許可 key path は `polling_interval_minutes`、`github_pat_expires_at`、`cooldown_seconds`、`maintenance_mode`、`notification`、`dashboard_layout`、`api_rate_limit` に固定する。未知 path、空 path、`..`、slash 始まりは `-32602 Invalid params` とする。

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

tool params に未知 key がある場合、必須 key 不足、型不一致、範囲外、secret 値を許可しない field への secret 形状入力は `-32602 Invalid params` とし、tool 実行、audit、metrics 更新を行わない。

`tools/call` params は exact `{ "name": string, "arguments": object }` とする。`arguments` 省略、`arguments:null`、`arguments` が object 以外、`name` が空文字、`name` が 128 byte 超過、params の未知 key は `-32602 Invalid params` とする。tool 別 params schema の「params は `{}`」は `arguments` が空 object であることを意味し、`arguments` 省略を許可しない。

`tools/call` は以下の順序で処理する。

1. tool name を [`docs/details/mcp.md` 詳細本文責務 §29.5](mcp.md#sec-29-5) の固定表に照合する。
2. params を本節の tool 別 params schema に照合する。
3. `--read-only`、scope、Authorization を確認する。
4. 副作用 tool は [`docs/details/mcp.md` 詳細本文責務 §29.15](mcp.md#sec-29-15) の confirmation を確認する。
5. 副作用 tool は対象 owner component を呼び出す前に [`docs/details/statefile.md` 詳細本文責務 §22.0d](statefile.md#sec-22-0d) の `McpAuditRecord` を `status:"accepted"`、`duration_ms:0` で `.mcp_audit_log` へ追記する。
6. 対象 owner component または sampling request を実行する。
7. `adlaire.analyzeBuildError` は sampling result を [`docs/details/statefile.md` 詳細本文責務 §22.0d](statefile.md#sec-22-0d) の `McpAuditRecord` として `.mcp_audit_log` へ追記する。
8. 実行結果を [`docs/details/statefile.md` 詳細本文責務 §22.0d](statefile.md#sec-22-0d) の `McpMetrics` へ反映する。
9. JSON-RPC success response または error response を返す。

上記 1〜4 の失敗では、tool 実行、audit、metrics 更新、対象 owner state 更新を行わない。副作用 tool の audit 追記に失敗した場合は `-32603 Internal error` を返し、対象 owner component を呼び出さない。sampling result の audit 追記に失敗した場合は `-32603 Internal error` を返し、metrics を更新しない。対象 owner component または sampling request の開始後に失敗した場合は、結果 status を `.mcp_metrics` に反映する。metrics 更新失敗は `-32603 Internal error` とし、完了済みの対象 owner state を rollback しない。

<a id="sec-29-7"></a>
**29.7 Resources：**

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

`resources/list` の resource order は [`docs/details/mcp.md` 詳細本文責務 §29.7](mcp.md#sec-29-7) の固定表の上から下の順とする。`resources/read` は statefile を補完、修復、初期化、削除、書き戻ししてはならない。破損 state を読んだ場合は `-32603 Internal error` とし、破損内容、absolute path、Go error を response、SSE、audit、metrics へ出力しない。

<a id="sec-29-8"></a>
**29.8 Resource subscription：**

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

resource update frame は `event: resource-updated\n`、`data: <compact-json>\n`、空行 `\n` の 3 要素をこの順で送信する。`data:` は 1 行だけとし、複数 `data:`、CRLF、`retry:`、event id、comment 混在を使用しない。frame 全体の write が完了した場合だけ flush し、partial write、client disconnect、flush error では以後の frame を送信しない。client disconnect は状態変更、audit、metrics 更新の失敗として扱わない。

<a id="sec-29-9"></a>
**29.9 Prompts：**

| prompt name | params | 用途 |
|-------------|--------|------|
| `adlaire.buildFailureTriage` | `{ "build_id": "string" }` | build 失敗原因の確認手順生成 |
| `adlaire.releaseReadiness` | `{ "target": "string" }` | release 前確認項目生成 |
| `adlaire.configReview` | `{ "snapshot_id": "string" }` | 設定 snapshot の差分確認 |

`prompts/get` result は message list を返す。

`prompts/list` は `name`、`description`、`arguments` を返す。`prompts/get` は `arguments` を prompt ごとの params と照合し、必須 key 不足、未知 key、型不一致を `-32602 Invalid params` とする。`PromptMessage` は `role:"user"` と `content.type:"text"` だけを使用する。prompt text は実行を促す手順文だけを生成し、server 側で build、config、queue を変更しない。

prompt は実行を伴わない。

<a id="sec-29-10"></a>
**29.10 Sampling：**

`adlaire.analyzeBuildError` は sampling 対応 tool とする。

server は外部 AI API を直接呼び出してはならない。

server は MCP client へ `sampling/createMessage` request を送信する。

sampling request timeout は `.mcp_config.sampling_timeout_ms` を使用する。

timeout 時は `-32003 Timeout` を返す。

sampling result は [`docs/details/statefile.md` 詳細本文責務 §22.0d](statefile.md#sec-22-0d) の `McpAuditRecord` として `.mcp_audit_log` へ追記する。`tool` は `adlaire.analyzeBuildError`、`confirmation_id` は `null`、`build_id` と `prompt_hash` は必須、`status` は `success`、`error`、`timeout` のいずれかとする。

sampling result の本文は statefile へ保存しない。

<a id="sec-29-11"></a>
**29.11 Notifications：**

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
**29.12 HTTP SSE transport：**

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

keepalive、resource update、notification、shutdown の各 SSE frame は 1 frame ごとに write と flush を完了させる。frame 生成時刻、送信順、disconnect reason は fixture の `expected/events.json` と `expected/effects.json` で検証できるように固定する。SSE 接続確立後に認証状態を変更しない。server shutdown frame の送信に失敗した場合でも、追加 JSON response、audit、metrics、statefile write を行わない。

<a id="sec-29-13"></a>
**29.13 Tool scope：**

scope は [`docs/details/statefile.md` 詳細本文責務 §22.0d](statefile.md#sec-22-0d) の `.mcp_config.scopes` に定義する。

token 未指定起動時は local trusted mode とし、loopback 接続に限り全 scope を許可する。この場合 `.mcp_config.scopes` は参照しない。

token 指定起動時は `--client-token` 値の SHA-256 lowercase hex を `.mcp_config.scopes[].token_hash` と照合し、一致した record の `scopes` だけを許可する。該当 record が存在しない場合は `-32002 Forbidden` を返す。

scope 不足時は `-32002 Forbidden` を返す。

`write:*` scope は `read:*` scope を暗黙に含めない。

<a id="sec-29-14"></a>
**29.14 Audit / client / metrics：**

副作用 tool 実行時は対象 owner component を呼び出す前に `.mcp_audit_log` へ以下を追記する。

```json
{
  "id": "string",
  "timestamp": "2026-09-28T00:00:00Z",
  "client_name": "string",
  "tool": "adlaire.triggerBuild",
  "params_hash": "sha256",
  "scope": "write:build",
  "confirmation_id": "string",
  "status": "accepted",
  "duration_ms": 0,
  "jsonrpc_error_code": null,
  "build_id": null,
  "prompt_hash": null
}
```

全 tool 実行は [`docs/details/statefile.md` 詳細本文責務 §22.0d](statefile.md#sec-22-0d) の `McpMetrics` へ aggregate する。

metrics は tool name を key とし、value は `success_count`、`error_count`、`timeout_count`、`last_status`、`last_duration_ms`、`last_at` を持つ。

client 接続は [`docs/details/statefile.md` 詳細本文責務 §22.0d](statefile.md#sec-22-0d) の `McpClientRecord` として `.mcp_client_log` へ append-only で記録する。

`.mcp_audit_log`、`.mcp_client_log`、`.mcp_metrics` の read / write / append は statefile owner の lock、atomic write、append 契約に従う。audit または client log の追記失敗時に対象 owner の state を先に変更してはならない。metrics 更新失敗時は、対象 owner の実行が完了済みの場合だけ rollback せず、JSON-RPC `-32603 Internal error` を返す。metrics 更新失敗を success response に隠してはならない。

<a id="sec-29-15"></a>
**29.15 Timeout / config CRUD / elicitation：**

`.mcp_config` schema は [`docs/details/statefile.md` 詳細本文責務 §22.0d](statefile.md#sec-22-0d) を正本とする。

tool timeout は `.mcp_config.tool_timeout_ms` を使用する。

`adlaire.getMcpConfig` は `.mcp_config` 全体を返す。

`adlaire.setMcpConfig` は `.mcp_config` の許可 key のみ更新する。`scopes` 更新時は `token_hash` と `scopes` を含む `McpScopeRecord[]` 全体を置換し、token 本体または Authorization header 値を受け取った場合は `-32602 Invalid params` とする。

副作用 tool は以下の elicitation flow を必須とする。

1. 初回 request で `confirmation_id` がない場合、`-32004 Confirmation required` を返し、`data.confirmation_id` と `data.summary` を含める。
2. client は同じ params に `confirmation_id` を付けて再要求する。
3. server は params hash、confirmation_id、有効期限、tool name が一致する場合だけ実行する。
4. confirmation 有効期限は 5 分とする。

confirmation_id は memory only とし、statefile へ保存しない。

confirmation_id は `mcpconf_` + 128 bit 以上の乱数を Crockford Base32 26 文字で表現する。params hash は confirmation_id を除いた canonical JSON の SHA-256 lowercase hex とする。期限切れ、tool 名不一致、params hash 不一致、read-only mode、scope 不足は tool を実行せず、audit と metrics を更新しない。

<a id="sec-29-16"></a>
**29.16 詳細仕様対象：**

本書は以下の `mcp` owner component 詳細仕様を定義する。現在状態は [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務を参照する。

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
**29.17 検証条件：**

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
- confirmation 付き副作用 tool が対象 owner component 呼び出し前に `status:"accepted"` の `.mcp_audit_log` を追記し、実行結果を `.mcp_metrics` の `success_count`、`error_count`、`timeout_count` のいずれかへ反映する。
- tool timeout 超過が `-32003 Timeout` を返す。
- SSE が keepalive、resource update、shutdown を送信する。
- sampling tool が server から外部 AI API を直接呼び出さない。
- `.mcp_client_log`、`.mcp_metrics`、`.mcp_audit_log` の fixture assertion が存在する。
- invalid CLI option、invalid addr、non-loopback 拒否、missing state-dir が固定 stderr と終了 code `2` を返す。
- `/mcp` の invalid content type、body 上限超過、JSON parse 失敗、batch request、未初期化 method、params 不正が固定 error を返し、tool 副作用を開始しない。
- `--client-token` 指定時、`/mcp` と `/mcp/events` は token 不一致を `401` で拒否し、`/health` は token 不要で応答する。
- read-only mode では副作用 tool が `tools/list` に出ず、直接 `tools/call` されても `-32002 Forbidden` で状態を変更しない。
- confirmation_id は memory only、5 分で期限切れ、params hash 不一致時に tool を実行しない。
- [`docs/details/fixture.md` fixture 証跡責務 §30-F MCP fixture 固定契約](fixture.md#mcp-fixture-contract) の expected、category、assertions、no-write、secret-mask、order を満たす。
