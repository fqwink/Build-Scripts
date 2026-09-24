# Adlaire CI — Admin 詳細仕様

本ファイルは `admin` owner component の詳細本文責務として、`admin` が主本文として持つ実装契約だけを扱う。

owner / collaborator 境界管理は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0b.1](../DETAIL_INDEX.md#0b1-owner-component-別-owner-collaborator-境界管理) に従う。`admin` owner component の主本文であり、collaborator component の仕様は呼び出し境界、配布境界、検証観点として参照する。fixture、expected、fake、実装検証証跡は [`docs/details/fixture.md`](fixture.md) fixture 証跡責務を参照する。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `admin` |
| collaborator component | `api`、`ui`、`sdk`、`setup` |
| 持つ内容 | `admin` owner が主本文として定義する管理 UI 静的ファイルの配布物構成、配置、検証、HTTP 静的配信境界。 |
| 持たない内容 | UI DOM 詳細、SDK method 実装、API endpoint 実装、状態 schema、systemd 導入手順、release asset 取得手順、fixture 証跡責務。 |

`admin` は、管理 UI 静的ファイルの中身を生成・変更してはならない。`ui` の仕様は [`docs/details/ui.md`](ui.md) 詳細本文責務を基準とし、`sdk` の仕様は [`docs/details/sdk.md`](sdk.md) 詳細本文責務を参照する。

---

## A1. 管理 UI 静的ファイル境界

標準管理 UI の配布物は、次の path に固定する。

| 配布元 path | 配置先 path | 必須 | 内容確認先 |
|-------------|-------------|------|------------|
| `admin/index.html` | `$INSTALL_DIR/admin/index.html` | 必須 | [`docs/details/ui.md`](ui.md) 詳細本文責務 |
| `admin/adlaire-ci-sdk.js` | `$INSTALL_DIR/admin/adlaire-ci-sdk.js` | 必須 | [`docs/details/sdk.md`](sdk.md) 詳細本文責務 |
| `admin/style.css` | `$INSTALL_DIR/admin/style.css` | 任意 | [`docs/details/ui.md`](ui.md) 詳細本文責務 |
| `admin/app.js` | `$INSTALL_DIR/admin/app.js` | 任意 | [`docs/details/ui.md`](ui.md) 詳細本文責務 |

配布物に [`docs/details/admin.md`](admin.md) 詳細本文責務 §A1 の配布物固定表以外のファイルを含める場合は、先に [`docs/details/admin.md`](admin.md) 詳細本文責務 §A1 の配布物固定表へ path、必須区分、内容確認先を追加する。未記載ファイルを暗黙に配布してはならない。

`admin-ui.tar.gz` の archive root は `admin/` directory を含めず、展開直後の root 直下に `index.html` と `adlaire-ci-sdk.js` が存在する形式とする。setup は検証済み archive を `$INSTALL_DIR/admin/` へ配置する。

---

## A2. 管理 UI Archive 検証

admin archive の検証は以下の順序に固定する。

1. archive entry を全件列挙する。
2. entry path が相対 path であることを確認する。
3. entry path に空文字、`.`、`..`、絶対 path、backslash、NUL byte を含まないことを確認する。
4. symlink、hardlink、device file、FIFO、socket を拒否する。
5. `index.html` と `adlaire-ci-sdk.js` が root 直下に 1 件ずつ存在することを確認する。
6. 任意 file が存在する場合は、[`docs/details/admin.md`](admin.md) 詳細本文責務 A1 の表に定義された path だけであることを確認する。
7. 通常 file の展開後 mode を `0644`、directory mode を `0755` に固定する。

いずれかの検証に失敗した場合は、既存 `$INSTALL_DIR/admin/` を変更しない。

---

## A3. 静的配信契約

`api` が管理 UI を配信する場合、`admin` は静的 file 解決と response header 決定だけを担当する。

| request path | file | Content-Type | Cache-Control |
|--------------|------|--------------|---------------|
| `/` | `$INSTALL_DIR/admin/index.html` | `text/html; charset=utf-8` | `no-store` |
| `/admin/` | `$INSTALL_DIR/admin/index.html` | `text/html; charset=utf-8` | `no-store` |
| `/admin/index.html` | `$INSTALL_DIR/admin/index.html` | `text/html; charset=utf-8` | `no-store` |
| `/admin/adlaire-ci-sdk.js` | `$INSTALL_DIR/admin/adlaire-ci-sdk.js` | `text/javascript; charset=utf-8` | `no-cache` |
| `/admin/style.css` | `$INSTALL_DIR/admin/style.css` | `text/css; charset=utf-8` | `no-cache` |
| `/admin/app.js` | `$INSTALL_DIR/admin/app.js` | `text/javascript; charset=utf-8` | `no-cache` |

未定義 path、directory listing、path traversal、hidden file、状態ファイル、secret file へのアクセスは `404` とする。認証前に配信する file は [`docs/details/admin.md`](admin.md) 詳細本文責務 §A3 の固定表の静的 file だけとし、API response、状態ファイル、credential、build log、snapshot を静的配信してはならない。

静的配信処理は request body を読まない。`GET` と `HEAD` 以外の method は `405` を返す。

---

## A4. Setup 連携境界

setup 手順本文、release asset 取得手順、rollback 手順は [`docs/details/setup.md`](setup.md) 詳細本文責務を参照する。

[`docs/details/admin.md`](admin.md) 詳細本文責務は、`admin` owner の配布境界として、setup が扱う `admin-ui.tar.gz` の中身、展開後の必須 file、静的配信 path、拒否すべき archive entry を定義する。

setup が admin UI を配置する場合は、以下を満たす。

| 確認 | 合格条件 |
|------|----------|
| archive 内容 | [`docs/details/admin.md`](admin.md) 詳細本文責務 A1 の配布物だけを含む。 |
| 必須 file | `index.html` と `adlaire-ci-sdk.js` が root 直下に存在する。 |
| path 安全性 | [`docs/details/admin.md`](admin.md) 詳細本文責務 A2 の拒否条件に該当しない。 |
| 配置結果 | `$INSTALL_DIR/admin/index.html` と `$INSTALL_DIR/admin/adlaire-ci-sdk.js` が通常 file として存在する。 |
| 失敗時 | 既存 `$INSTALL_DIR/admin/` を変更しない。 |

---

## A5. 受け入れ条件

`admin` の詳細実装確認では、以下をすべて満たす。

| 観点 | 合格条件 |
|------|----------|
| file boundary | [`docs/details/admin.md`](admin.md) 詳細本文責務 A1 の必須 file を検証し、未定義 file を暗黙配布しない。 |
| archive safety | [`docs/details/admin.md`](admin.md) 詳細本文責務 A2 の危険 entry をすべて拒否する。 |
| serving | [`docs/details/admin.md`](admin.md) 詳細本文責務 A3 の path、Content-Type、Cache-Control、method、404 / 405 が一致する。 |
| no mutation | UI / SDK file 内容、状態ファイル、credential、build log、snapshot を変更しない。 |
| no secret exposure | `.admin_credentials`、`.github_token`、`.server_config`、`.build_logs`、`.snapshots` を静的配信しない。 |
| setup integration | [`docs/details/setup.md` 詳細本文責務 §26](setup.md#26-セットアップアップデート手順) の配置・rollback 条件と矛盾しない。 |

---

## A6. Admin fixture 参照契約

`admin` owner component は、配布物検証、archive 安全性、静的配信、no mutation を fixture で確認できる状態にする。`admin` 詳細では UI DOM、SDK method、API endpoint、fixture 入力、expected、fake、実装検証証跡を再定義せず、admin 配布境界だけを確認する。

Admin fixture の fixture 名、入力、操作、expected file、禁止副作用は [`docs/details/fixture.md`](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) fixture 証跡責務 §27-F setup / admin / release 連動 fixture 固定契約を正本とする。

**Admin 実装確認ゲート：**

| 観点 | 合格条件 |
|------|----------|
| archive validation | [`docs/details/admin.md`](admin.md) 詳細本文責務 A2 と [`docs/details/fixture.md`](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) fixture 証跡責務 §27-F の admin archive fixture が成功し、失敗時に既存 admin directory 差分がない。 |
| static serving | [`docs/details/admin.md`](admin.md) 詳細本文責務 A3 と [`docs/details/fixture.md`](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) fixture 証跡責務 §27-F の admin request fixture が status、header、body 有無、method 制限に一致する。 |
| secret isolation | secret、state、log、snapshot path への direct request がすべて `404` で、response body に secret 原文を含まない。 |
| no generation | admin は UI / SDK file 内容を生成・整形・書換しない。配布と配信だけを行う。 |
| setup integration | [`docs/details/setup.md` 詳細本文責務 §26.8](setup.md#sec-26-8) の admin archive 展開、差分確認、rollback 条件と同じ expected を参照する。 |
| fixture integration | [`docs/details/fixture.md`](fixture.md) fixture 証跡責務の `setup-admin-release-layout`、`setup-admin-archive-boundary`、`setup-systemd-rollback-boundary`、`admin-static-serving-security`、`setup-secret-preservation` と fixture 名、expected file、禁止副作用が一致する。 |
