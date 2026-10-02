# Adlaire CI — Admin 詳細仕様

---

<a id="0-責務境界"></a>

**0. 責務境界：**

| 項目 | 内容 |
|------|------|
| owner component | `admin` |
| 実装主体 | 配布対象は [`admin/index.html`](../../admin/index.html) と [`admin/adlaire-ci-sdk.js`](../../admin/adlaire-ci-sdk.js)、CLI 管理クライアントは [`components/admin/`](../../components/admin/)、archive 生成は [`components/release/`](../../components/release/)、配置は [`components/setup/`](../../components/setup/)、HTTP 静的配信は [`components/api/`](../../components/api/) とする。archive 内容は本詳細本文、生成手順は [`docs/details/release.md`](release.md)、配置挙動は [`docs/details/setup.md` 詳細本文責務 §26.2a](setup.md#sec-26-2a) と [§26.8](setup.md#sec-26-8) を正本とする。 |
| 持つ内容 | `admin` owner が主本文として定義する管理 UI 静的ファイルの配布物構成、配置、検証、HTTP 静的配信境界、および CLI 管理クライアント。 |
| 検証接続 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 owner 詳細本文 検証接続共通入口](../DETAIL_INDEX.md#owner-detail-verification-route) を参照する。 |

`admin` は、管理 UI 静的ファイルの中身を生成・変更してはならない。`ui` の DOM と動作は [`docs/details/ui.md`](ui.md) 詳細本文責務、標準管理 UI の視覚値は [`docs/DESIGN.md` デザイン責務 標準管理 UI 視覚契約](../DESIGN.md#admin-ui-visual-contract)、`sdk` の仕様は [`docs/details/sdk.md`](sdk.md) 詳細本文責務を参照する。

---

<a id="a1-管理-ui-静的ファイル境界"></a>
**A1. 管理 UI 静的ファイル境界：**

標準管理 UI の配布物は、次の path に固定する。

| 配布元 path | 配置先 path | 必須 | 内容確認先 |
|-------------|-------------|------|------------|
| `admin/index.html` | `$INSTALL_DIR/admin/index.html` | 必須 | DOM と動作は [`docs/details/ui.md`](ui.md) 詳細本文責務、視覚値は [`docs/DESIGN.md` デザイン責務 標準管理 UI 視覚契約](../DESIGN.md#admin-ui-visual-contract) |
| `admin/adlaire-ci-sdk.js` | `$INSTALL_DIR/admin/adlaire-ci-sdk.js` | 必須 | [`docs/details/sdk.md`](sdk.md) 詳細本文責務 |

配布物に [`docs/details/admin.md` 詳細本文責務 §A1](admin.md#a1-管理-ui-静的ファイル境界) の配布物固定表以外のファイルを含める場合は、先に [`docs/details/admin.md` 詳細本文責務 §A1](admin.md#a1-管理-ui-静的ファイル境界) の配布物固定表へ path、必須区分、内容確認先を追加する。未記載ファイルを暗黙に配布してはならない。

`admin-ui.tar.gz` の archive root は `admin/` directory を含めず、展開直後の root 直下に `index.html` と `adlaire-ci-sdk.js` が存在する形式とする。setup は検証済み archive を `$INSTALL_DIR/admin/` へ配置する。

---

<a id="a2-管理-ui-archive-検証"></a>
**A2. 管理 UI Archive 検証：**

admin archive の検証は以下の順序に固定する。

1. archive entry を全件列挙する。
2. entry path が相対 path であることを確認する。
3. entry path に空文字、`.`、`..`、絶対 path、backslash、NUL byte を含まないことを確認する。
4. symlink、hardlink、device file、FIFO、socket を拒否する。
5. `index.html` と `adlaire-ci-sdk.js` が root 直下に 1 件ずつ存在することを確認する。
6. `index.html` と `adlaire-ci-sdk.js` 以外の entry が存在しないことを確認する。
7. 通常 file の展開後 mode を `0644`、directory mode を `0755` に固定する。

いずれかの検証に失敗した場合は、既存 `$INSTALL_DIR/admin/` を変更しない。

---

<a id="a3-静的配信契約"></a>
**A3. 静的配信契約：**

`api` が管理 UI を配信する場合、`admin` は静的 file 解決と response header 決定だけを担当する。

| request path | file | Content-Type | Cache-Control |
|--------------|------|--------------|---------------|
| `/` | `$INSTALL_DIR/admin/index.html` | `text/html; charset=utf-8` | `no-store` |
| `/admin/` | `$INSTALL_DIR/admin/index.html` | `text/html; charset=utf-8` | `no-store` |
| `/admin/index.html` | `$INSTALL_DIR/admin/index.html` | `text/html; charset=utf-8` | `no-store` |
| `/admin/adlaire-ci-sdk.js` | `$INSTALL_DIR/admin/adlaire-ci-sdk.js` | `text/javascript; charset=utf-8` | `no-cache` |

未定義 path、directory listing、path traversal、hidden file、状態ファイル、secret file へのアクセスは `404` とする。認証前に配信する file は [`docs/details/admin.md` 詳細本文責務 §A3](admin.md#a3-静的配信契約) の固定表の静的 file だけとし、API response、状態ファイル、credential、build log、snapshot を静的配信してはならない。

静的配信処理は request body を読まない。`GET` と `HEAD` 以外の method は `405` を返す。

---

<a id="a4-setup-連携境界"></a>
**A4. Setup 連携境界：**

setup 手順本文、Release asset 取得手順、rollback 手順は [`docs/details/setup.md`](setup.md) 詳細本文責務を参照する。

[`docs/details/admin.md`](admin.md) 詳細本文責務は、`admin` owner の配布境界として、setup が扱う `admin-ui.tar.gz` の中身、展開後の必須 file、静的配信 path、拒否すべき archive entry を定義する。

setup が admin UI を配置する場合、admin owner の正本本文は [`docs/details/admin.md` 詳細本文責務 §A1](admin.md#a1-管理-ui-静的ファイル境界) の配布物境界と [`docs/details/admin.md` 詳細本文責務 §A2](admin.md#a2-管理-ui-archive-検証) の archive 検証だけとする。配置順、rollback、既存 file 保護、終了コードは [`docs/details/setup.md`](setup.md) 詳細本文責務を参照する。

---

<a id="a5-受け入れ条件"></a>
**A5. 受け入れ条件：**

`admin` の詳細実装確認では、以下をすべて満たす。

| 観点 | 合格条件 |
|------|----------|
| file boundary | [`docs/details/admin.md` 詳細本文責務 §A1](admin.md#a1-管理-ui-静的ファイル境界) の必須 file を検証し、未定義 file を暗黙配布しない。 |
| archive safety | [`docs/details/admin.md` 詳細本文責務 §A2](admin.md#a2-管理-ui-archive-検証) の危険 entry をすべて拒否する。 |
| serving | [`docs/details/admin.md` 詳細本文責務 §A3](admin.md#a3-静的配信契約) の path、Content-Type、Cache-Control、method、404 / 405 が一致する。 |
| no mutation | UI / SDK file 内容、状態ファイル、credential、build log、snapshot を変更しない。 |
| no secret exposure | `.admin_credentials`、`.github_token`、`.server_config`、`.build_logs`、`.snapshots` を静的配信しない。 |
| setup integration | [`docs/details/setup.md` 詳細本文責務 §26](setup.md#26-セットアップアップデート手順) の配置・rollback 条件と矛盾しない。 |

---

<a id="a6-admin-fixture-参照契約"></a>
**A6. Admin fixture 参照契約：**

`admin` owner component は、配布物検証、archive 安全性、静的配信、no mutation を fixture で確認できる状態にする。`admin` 詳細では UI DOM、SDK method、API endpoint、fixture 入力、expected、fake、実装検証証跡を再定義せず、admin 配布境界だけを確認する。

Admin fixture の fixture 名、入力、操作、expected file、禁止副作用は [`docs/details/fixture.md` fixture 証跡責務 §27-F setup / admin / Release asset 連動 fixture 固定契約](fixture.md#sec-27-f-19) を正本とする。

Admin 実装確認は [`docs/details/admin.md` 詳細本文責務 §A5](admin.md#a5-受け入れ条件) の全条件と、[`docs/details/fixture.md` fixture 証跡責務 §27-F setup / admin / Release asset 連動 fixture 固定契約](fixture.md#sec-27-f-19) の該当条件を同時に満たした場合だけ合格とする。fixture 名、input、expected、fake、禁止副作用を本節で再掲してはならない。

<a id="sec-a7"></a>
**A7. CLI 管理クライアント：**

CLI 管理クライアントの owner は `admin` とする。

配布 binary 名は `adlaire-ci-admin` とする。

実装主体は `components/admin/` とする。

起動入口は `main.go` の basename dispatch とする。

CLI 管理クライアントは `api` owner の endpoint を呼び出す client であり、server side の状態、認証、認可、endpoint response を再定義しない。

**CLI 形式：**

```text
adlaire-ci-admin --api-url <url> --token <token> [--json] <command> [command-args]
adlaire-ci-admin --help
adlaire-ci-admin --version
```

CLI parse の共通優先順位、argv token safety、`--help`、`--version` は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 CLI 共通固定契約](../DETAIL_INDEX.md#common-cli-contract) を使用する。`--help` の stdout は `Usage: adlaire-ci-admin --api-url url --token token [--json] command [command-args]` + LF とする。`--version` の stdout は `adlaire-ci-admin <binary-version> go=<runtime.Version()>` + LF とする。`--help` と `--version` は API URL 検証、token 検証、network、状態 file read/write、乱数取得を行わない。

option は `--name value` の 2 token 形式だけを許可する。`--name=value`、短縮 option、位置引数による option 値、同一 option の重複、未知 option は parse error とする。`--json` は値を取らない boolean option とし、複数回指定は parse error とする。`--api-url` は `http://` または `https://` の absolute URL とし、host 必須、userinfo、query、fragment を禁止する。path は空、または `/` から始まる clean path だけを許可し、`..` segment、重複 slash、backslash、NUL byte を禁止する。末尾 `/` は 1 個だけ除去し、`/api` を暗黙追加しない。`--token` は 1〜4096 byte の UTF-8 text とし、空文字、NUL、CR、LF を禁止する。

CLI 管理クライアントの parse / validation は次の順序に固定する。各順序で不合格を検出した場合は最初の 1 件だけを返し、後続順序を評価してはならない。network、state read/write、file read/write、乱数取得、HTTP request 作成、request body 構築は、次表の全順序が合格した後だけ開始する。

| 順序 | 判定 | 不合格時 |
|------|------|----------|
| 1 | exact `--help` / exact `--version` を [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 CLI 共通固定契約](../DETAIL_INDEX.md#common-cli-contract) の優先順位で判定する。 | `--help` または `--version` の成功終了とし、他 token を検証しない。 |
| 2 | argv token safety を [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 CLI 共通固定契約](../DETAIL_INDEX.md#common-cli-contract) で判定する。 | stdout 空、stderr `invalid command line token` + LF、終了 code `2`。 |
| 3 | command より前の option token を判定する。許可 option は `--api-url`、`--token`、`--json` だけとする。`--api-url` と `--token` は `--name value` だけを許可し、`--json` は値を取らない。 | 共通契約が固定する unknown option / missing value は共通契約の stderr、同一 option 重複と `--json` 複数指定は `usage error`。 |
| 4 | `--api-url`、`--token`、command の存在を判定する。 | stdout 空、stderr `usage error` + LF、終了 code `2`。 |
| 5 | `--api-url` の URL scheme、host、userinfo、query、fragment、path、末尾 `/` 正規化を判定する。 | stdout 空、stderr `usage error` + LF、終了 code `2`。 |
| 6 | `--token` の byte 長、UTF-8、禁止文字を判定する。 | stdout 空、stderr `usage error` + LF、終了 code `2`。 |
| 7 | command 固定表の 7 command と一致するかを判定する。 | stdout 空、stderr `unknown command: <command>` + LF、終了 code `2`。 |
| 8 | command-args の個数、`cancel-queue` の `<queue_id>`、`config-snapshot` の `label` を判定する。 | stdout 空、stderr `usage error` + LF、終了 code `2`。 |
| 9 | command、正規化済み `--api-url`、token、command-args から HTTP request method、path、header、body を決定する。 | この順序では parse error を発生させない。 |

parse / validation 失敗時は stdout 空、stderr 固定 1 行、終了 code `2` とし、API URL、token、command-args、HTTP request body、Go error、OS error、絶対 path、network response を stdout、stderr、log、fixture expected に出力してはならない。`--api-url` の正規化は順序 5 だけで行い、順序 9 で path clean、percent decode、query 付与、`/api` 重複除去を行ってはならない。

API request path は、正規化後の `--api-url` path prefix と command 固定 path を byte 連結して作る。例として `--api-url http://127.0.0.1:8765` の `status` は `http://127.0.0.1:8765/api/status`、`--api-url http://127.0.0.1:8765/adlaire` の `status` は `http://127.0.0.1:8765/adlaire/api/status` とする。path 連結時に percent decode、path clean、`/api` の重複除去、query 付与を行ってはならない。

| command | command-args | 呼び出す API | request body | stdout |
|---------|--------------|--------------|--------------|--------|
| `status` | なし | `GET /api/status` | なし | `status=<last_build_status> running=<running>` + LF。`--json` 指定時は API response JSON + LF。 |
| `queue` | なし | `GET /api/queue` | なし | `active=<id-or-none> queued=<count>` + LF。`--json` 指定時は API response JSON + LF。 |
| `history` | なし | `GET /api/history` | なし | `total=<total> latest=<id-or-none>` + LF。`--json` 指定時は API response JSON + LF。 |
| `trigger-build` | なし | `POST /api/build` | `{}` | `queued=<queue_id>` + LF。`--json` 指定時は API response JSON + LF。 |
| `cancel-queue` | `<queue_id>` | `DELETE /api/queue/{queue_id}` | なし | `queue cancelled` + LF。`--json` 指定時は API response JSON + LF。 |
| `config-snapshot` | `[label]` | `POST /api/config-snapshots` | 未指定時 `{"label":null}` / 指定時 `{"label":"<label>"}` | `snapshot=<id>` + LF。`--json` 指定時は API response JSON + LF。 |
| `events` | なし | `GET /api/events` | なし | `events=<total>` + LF。`--json` 指定時は API response JSON + LF。 |

`cancel-queue` の `<queue_id>` は path parameter として 1 回だけ percent encode する。空文字、`/`、`..`、NUL byte を含む値は API 呼び出し前に parse error とする。`config-snapshot` の `label` は未指定なら `null`、指定時は 1〜128 Unicode scalar values とし、改行、NUL byte、BOM を禁止する。

`trigger-build` の request body は byte 列 `{}` とする。`config-snapshot` の request body は未指定時 `{"label":null}`、指定時 `{"label":"<label>"}` とし、JSON string escape は RFC 8259 に従う。request body は UTF-8、末尾 LF なし、余分な空白なしとする。

CLI 管理クライアントは command 固定表の 7 command だけを実装する。API 側に存在する他 endpoint を CLI command として追加する場合は、先に [`docs/details/admin.md` 詳細本文責務 §A7](admin.md#sec-a7) の command 固定表、stdout 固定表、検証条件、[`docs/details/fixture.md` fixture 証跡責務 Admin CLI fixture 固定契約](fixture.md#admin-cli-fixture-contract) を同じ変更で改訂する。未記載 endpoint を汎用 passthrough、任意 path、任意 method、任意 JSON body として呼び出してはならない。

`--token` の値を stdout、stderr、server log、fixture expected に出力してはならない。

全 API 呼び出しは `Authorization: Bearer <token>`、`Accept: application/json`、`User-Agent: adlaire-ci-admin/<binary-version>` を送信する。body を持つ command は `Content-Type: application/json` を送信する。Cookie、Referer、X-Forwarded-*、環境変数由来 proxy、redirect 追従、retry、connection reuse 前提の状態保持を使用してはならない。timeout は dial、TLS handshake、request body write、response header read、response body read の全体で 30 秒固定とする。

HTTP response は body 全体を上限 1 MiB まで読む。1 MiB を超える場合は stdout 空、stderr `api error: invalid response` + LF、終了 code `1` とし、body の残り、token、URL query、header 値を出力しない。success は HTTP status `200`〜`299` だけとする。`300`〜`599` は stdout 空、stderr `api error: <status>` + LF、終了 code `1` とする。`1xx`、status なし、redirect response、HTTP protocol error は network error と同じ扱いにする。

success response は `Content-Type` が `application/json` または `application/json; charset=utf-8` であることを必須とする。不一致、複数 `Content-Type`、body 空、JSON として単一値でない body は stdout 空、stderr `api error: invalid response` + LF、終了 code `1` とする。

`--json` 指定時は、検証済み API response body から前後の ASCII whitespace だけを除去し、残った byte 列をそのまま stdout へ出し、最後に LF 1 個を付ける。object key order、number 表現、string escape は API wire response を保持し、CLI 側で再 encode しない。`--json` 指定時でも command ごとの成功 response shape が [`docs/details/api.md`](api.md) 詳細本文責務の対象 response 契約を満たさない場合は invalid response とする。`--json` 未指定時は CLI が response JSON を parse し、上表の stdout へ写像する。必要 key 欠落、型不一致、`null` 不許可 field、余分な stderr 出力が必要になる body は invalid response とする。response JSON parse 失敗は stdout 空、stderr `api error: invalid response` + LF、終了 code `1` とする。

human stdout の field は固定順とする。`status` は `last_build_status` を string としてそのまま出し、`running` は JSON boolean を `true` / `false` の lowercase で出す。`queue` は `active` が `null` の場合だけ `active=none` とし、`active` が object の場合は `active.id` の string を出す。`queue` の `queued` は API response の `queued` array length だけを 10 進数で出し、数値 `queued`、`active` object の `id` 欠落、`active.id` の string 以外、`queued` の array 以外は invalid response とする。`history` は `total` を 10 進数で出し、`latest` は `history` array の先頭要素 `id` string、空配列の場合だけ `none` とする。`total` の integer 以外、`history` の array 以外、非空 `history[0].id` の string 以外は invalid response とする。`trigger-build` は `queue_id` の string を必須とし、`queued` は boolean `true`、`dispatch` は `requested` または `timer_fallback` だけを許可する。`queue_id` 欠落、空文字、`queued:false`、未知 `dispatch` は invalid response とし、`none` へ置換しない。`cancel-queue` は response body の `message` に依存せず固定 `queue cancelled` とする。`config-snapshot` は `id` の string を必須とし、欠落、空文字、string 以外は invalid response とする。`events` は `total` integer を 10 進数で出し、`events` array は paging 後の要素列として存在を検証するが stdout の件数には使用しない。`events` array length と `total` が異なることは、[`docs/details/api.md` 詳細本文責務 §27.66](api.md#sec-27-66) の paging 契約上 valid とする。`total` の integer 以外、`events` の array 以外は invalid response とする。

未知 command は stdout 空、stderr `unknown command: <command>` + LF、終了 code `2` とする。未知 command の `<command>` は argv token safety 合格後、かつ `--token` の option 値ではない command token だけを使用する。

[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 CLI 共通固定契約](../DETAIL_INDEX.md#common-cli-contract) が固定する argv token safety、unknown option、missing value、禁止位置引数の stderr は共通契約に従う。共通契約で stderr が固定されない admin 固有の同一 option 重複、`--json` 複数指定、引数不足、引数過多、`--api-url` 不正、`--token` 不正、command 固有引数不正は stdout 空、stderr `usage error` + LF、終了 code `2` とする。

API が `2xx` 以外を返した場合、CLI は stdout 空、stderr `api error: <status>` + LF、終了 code `1` とする。

`--json` 指定時でも error response body を stderr に出力してはならない。

network error、TLS error、timeout、connection close before response は stdout 空、stderr `api error: connection failed` + LF、終了 code `1` とする。timeout は 30 秒固定とし、retry しない。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| help | 終了 code `0`、状態変更なし。 |
| version | 終了 code `0`、状態変更なし。 |
| token redaction | token が stdout / stderr / log / fixture expected に出現しない。 |
| unknown command | 終了 code `2`。 |
| api 401 | 終了 code `1`、token を出さない。 |
| api redirect | 終了 code `1`、redirect 追従なし、Location を出さない。 |
| trigger-build | `POST /api/build` だけを呼び、`POST /api/builds` を呼ばない。 |
| cancel-queue path encode | queue id を 1 回だけ percent encode する。 |
| json mode | success 時だけ API response JSON を stdout へ出し、error body は出さない。 |
| human output mapping | command ごとの固定 field、固定順、LF 1 個、stderr 0 byte を照合する。 |
| response body limit | 1 MiB 超過 response を invalid response とし、body 断片を出さない。 |
| fixture linkage | [`docs/details/fixture.md` fixture 証跡責務 Admin CLI fixture 固定契約](fixture.md#admin-cli-fixture-contract) の fixture を満たす。 |

---

<a id="phase-13-admin-alignment-contract"></a>
**Phase 13 admin 実装整合契約：**

[`docs/details/admin.md`](admin.md) 詳細本文責務は、Phase 13 で CLI 管理クライアントの command、API URL 連結、token 取得、stdout / stderr、API route 照合、管理 UI 配布境界を所有する。Phase 13 の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](../DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry)、完了証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 13 実装整合・品質改善証跡](fixture.md#phase-13-implementation-alignment-quality-evidence) を参照する。

| 対象 | 固定契約 |
|------|----------|
| token 取得 | Phase 13 実装では token を command line 引数の平文値として受け取る経路を廃止し、token file または stdin だけから取得する。token file の安全読込、secret mask、保存禁止は [`docs/details/security.md` 詳細本文責務 Phase 13 security 実装整合契約](security.md#phase-13-security-alignment-contract) を参照する。 |
| URL / route 照合 | CLI command 固定表の全 method、path、query、body、response mapping を [`docs/details/api.md` 詳細本文責務 Phase 13 API 実装整合契約](api.md#phase-13-api-alignment-contract) の API route 正本と自動照合する。表にない URL へ送信する汎用 passthrough を作らない。 |
| stdout / stderr | setup/API credential 初期化と連動する admin 出力は [`docs/details/setup.md` 詳細本文責務 Phase 13 setup 実装整合契約](setup.md#phase-13-setup-alignment-contract) と一致させる。成功 stdout、失敗 stderr、終了 code、secret 非表示を fixture で照合する。 |
| HTTP safety | redirect 追従、環境変数 proxy、Cookie、Referer、token query 化、retry、1 MiB 超過 body 出力を禁止する。timeout、body 上限、invalid response は固定 error へ写像する。 |
| 配布境界 | 管理 UI 静的 file の存在、archive safety、no mutation、secret exposure 禁止は [`docs/details/admin.md` 詳細本文責務 §A1](admin.md#a1-管理-ui-静的ファイル境界) から [`docs/details/admin.md` 詳細本文責務 §A6](admin.md#a6-admin-fixture-参照契約) を維持し、Phase 13 では API route parity と token 取得経路だけを追加確認する。 |

Phase 13 の Admin CLI command は `input/contract_inventory.json` の `client_bindings.admin_command` によって API route contract と 1 対 1 に紐づく。Admin CLI にだけ存在する command、API route に紐づかない URL、または複数 API route へ曖昧に展開される command は `phase13_contract_mismatch_count` に計上する。

Phase 13 の `admin` 実装は、[`docs/details/fixture.md` fixture 証跡責務 Phase 13 実装整合・品質改善証跡](fixture.md#phase-13-implementation-alignment-quality-evidence) の `phase13_contract_mismatch_count=0`、`phase13_token_arg_open_count=0`、`phase13_external_boundary_open_count=0` を満たすまで完了扱いにしてはならない。
