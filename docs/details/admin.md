# Adlaire CI — Admin 詳細仕様

---

<a id="0-責務境界"></a>

**0. 責務境界：**

| 項目 | 内容 |
|------|------|
| owner component | `admin` |
| 実装主体 | 単独の Go artifact は持たない。配布対象は [`admin/index.html`](../../admin/index.html) と [`admin/adlaire-ci-sdk.js`](../../admin/adlaire-ci-sdk.js)、archive生成は`components/release.go`、配置は`components/setup.go`、HTTP静的配信は[`components/api.go`](../../components/api.go)とする。archive内容は本詳細本文、生成手順は[`docs/details/release.md`](release.md)、配置挙動は[`docs/details/setup.md` 詳細本文責務 §26.2a](setup.md#sec-26-2a)と[§26.8](setup.md#sec-26-8)を正本とする。 |
| 持つ内容 | `admin` owner が主本文として定義する管理 UI 静的ファイルの配布物構成、配置、検証、HTTP 静的配信境界。 |

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
## A7. CLI 管理クライアント

CLI 管理クライアントの owner は `admin` とする。

配布 binary 名は `adlaire-ci-admin` とする。

実装主体は `components/admin.go` とする。

起動入口は `main.go` の basename dispatch とする。

CLI 管理クライアントは `api` owner の endpoint を呼び出す client であり、server side の状態、認証、認可、endpoint response を再定義しない。

**CLI 形式：**

```text
adlaire-ci-admin --api-url <url> --token <token> <command> [--json]
adlaire-ci-admin --help
adlaire-ci-admin --version
```

| command | 呼び出す API | stdout |
|---------|--------------|--------|
| `status` | `GET /api/status` | status summary。`--json` 指定時は API response JSON。 |
| `queue` | `GET /api/queue` | queue summary。`--json` 指定時は API response JSON。 |
| `history` | `GET /api/history` | history summary。`--json` 指定時は API response JSON。 |
| `trigger-build` | `POST /api/builds` | created build id。`--json` 指定時は API response JSON。 |
| `cancel-queue` | `DELETE /api/queue/{queue_id}` | fixed message。`--json` 指定時は API response JSON。 |
| `config-snapshot` | `POST /api/config-snapshots` | snapshot id。`--json` 指定時は API response JSON。 |
| `events` | `GET /api/events` | event summary。`--json` 指定時は API response JSON。 |

`--token` の値を stdout、stderr、server log、fixture expected に出力してはならない。

未知 command は stdout 空、stderr `unknown command: <command>` + LF、終了 code `2` とする。

API が `4xx` または `5xx` を返した場合、CLI は stdout 空、stderr `api error: <status>` + LF、終了 code `1` とする。

`--json` 指定時でも error response body を stderr に出力してはならない。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| help | 終了 code `0`、状態変更なし。 |
| version | 終了 code `0`、状態変更なし。 |
| token redaction | token が stdout / stderr / log / fixture expected に出現しない。 |
| unknown command | 終了 code `2`。 |
| api 401 | 終了 code `1`、token を出さない。 |
