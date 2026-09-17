# Adlaire CI — Admin 詳細仕様

本ファイルは `docs/DETAIL_INDEX.md` から分割した `admin` owner component の詳細仕様である。

本ファイルに、方針、ポリシー、正本関係、実装状態、ロードマップ状態、実装可否の上位判断を記載してはならない。方針、ポリシー、正本関係は `docs/SPEC.md`、実装状態、ロードマップ状態、実装可否は `docs/ROADMAP.md` を正とする。

本ファイルを読む前に、`docs/SPEC.md` で方針とポリシーを確認し、`docs/ROADMAP.md` で実装状態と実装可否を確認し、`docs/DETAIL_INDEX.md` §0〜§0j で共通固定値、責務 component、詳細節対応表、リポジトリ内ソース配置を確認する。本ファイルは `admin` owner component の主本文であり、collaborator component の仕様は呼び出し境界、配布境界、検証観点として参照する。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `admin` |
| collaborator component | `api`、`ui`、`sdk`、`setup` |
| 持つ内容 | `admin` owner が主本文として定義する管理 UI 静的ファイルの配布物構成、配置、検証、HTTP 静的配信境界。 |
| 持たない内容 | UI DOM 詳細、SDK method 実装、API endpoint 実装、状態 schema、systemd 導入手順、release asset 取得手順、fixture / PR 証跡正本。 |

`admin` は、管理 UI 静的ファイルの中身を生成・変更してはならない。`ui` の仕様は `docs/details/ui.md` を正とし、`sdk` の仕様は `docs/details/sdk.md` を正とする。

---

## A1. Admin Static File Boundary

標準管理 UI の配布物は、次の path に固定する。

| 配布元 path | 配置先 path | 必須 | 内容の正本 |
|-------------|-------------|------|------------|
| `admin/index.html` | `$INSTALL_DIR/admin/index.html` | 必須 | `docs/details/ui.md` |
| `admin/adlaire-ci-sdk.js` | `$INSTALL_DIR/admin/adlaire-ci-sdk.js` | 必須 | `docs/details/sdk.md` |
| `admin/style.css` | `$INSTALL_DIR/admin/style.css` | 任意 | `docs/details/ui.md` |
| `admin/app.js` | `$INSTALL_DIR/admin/app.js` | 任意 | `docs/details/ui.md` |

配布物に上表以外のファイルを含める場合は、先に本表へ path、必須区分、内容の正本を追加する。未記載ファイルを暗黙に配布してはならない。

`admin-ui.tar.gz` の archive root は `admin/` directory を含めず、展開直後の root 直下に `index.html` と `adlaire-ci-sdk.js` が存在する形式とする。setup は検証済み archive を `$INSTALL_DIR/admin/` へ配置する。

---

## A2. Admin Archive Validation

admin archive の検証は以下の順序に固定する。

1. archive entry を全件列挙する。
2. entry path が相対 path であることを確認する。
3. entry path に空文字、`.`、`..`、絶対 path、backslash、NUL byte を含まないことを確認する。
4. symlink、hardlink、device file、FIFO、socket を拒否する。
5. `index.html` と `adlaire-ci-sdk.js` が root 直下に 1 件ずつ存在することを確認する。
6. 任意 file が存在する場合は、A1 の表に定義された path だけであることを確認する。
7. 通常 file の展開後 mode を `0644`、directory mode を `0755` に固定する。

いずれかの検証に失敗した場合は、既存 `$INSTALL_DIR/admin/` を変更しない。

---

## A3. Static Serving Contract

`api` が管理 UI を配信する場合、`admin` は静的 file 解決と response header 決定だけを担当する。

| request path | file | Content-Type | Cache-Control |
|--------------|------|--------------|---------------|
| `/` | `$INSTALL_DIR/admin/index.html` | `text/html; charset=utf-8` | `no-store` |
| `/admin/` | `$INSTALL_DIR/admin/index.html` | `text/html; charset=utf-8` | `no-store` |
| `/admin/index.html` | `$INSTALL_DIR/admin/index.html` | `text/html; charset=utf-8` | `no-store` |
| `/admin/adlaire-ci-sdk.js` | `$INSTALL_DIR/admin/adlaire-ci-sdk.js` | `text/javascript; charset=utf-8` | `no-cache` |
| `/admin/style.css` | `$INSTALL_DIR/admin/style.css` | `text/css; charset=utf-8` | `no-cache` |
| `/admin/app.js` | `$INSTALL_DIR/admin/app.js` | `text/javascript; charset=utf-8` | `no-cache` |

未定義 path、directory listing、path traversal、hidden file、状態ファイル、secret file へのアクセスは `404` とする。認証前に配信する file は上表の静的 file だけとし、API response、状態ファイル、credential、build log、snapshot を静的配信してはならない。

静的配信処理は request body を読まない。`GET` と `HEAD` 以外の method は `405` を返す。

---

## A4. Setup Boundary

`docs/details/setup.md` は、release asset 取得、checksum 検証、systemd、配置順序、rollback を扱う。

本ファイルは、setup が扱う `admin-ui.tar.gz` の中身、展開後の必須 file、静的配信 path、拒否すべき archive entry を定義する。

setup が admin UI を配置する場合は、以下を満たす。

| 確認 | 合格条件 |
|------|----------|
| archive 内容 | A1 の配布物だけを含む。 |
| 必須 file | `index.html` と `adlaire-ci-sdk.js` が root 直下に存在する。 |
| path 安全性 | A2 の拒否条件に該当しない。 |
| 配置結果 | `$INSTALL_DIR/admin/index.html` と `$INSTALL_DIR/admin/adlaire-ci-sdk.js` が通常 file として存在する。 |
| 失敗時 | 既存 `$INSTALL_DIR/admin/` を変更しない。 |

---

## A5. Acceptance Criteria

`admin` の実装完了には、以下をすべて満たす。

| 観点 | 合格条件 |
|------|----------|
| file boundary | A1 の必須 file を検証し、未定義 file を暗黙配布しない。 |
| archive safety | A2 の危険 entry をすべて拒否する。 |
| serving | A3 の path、Content-Type、Cache-Control、method、404 / 405 が一致する。 |
| no mutation | UI / SDK file 内容、状態ファイル、credential、build log、snapshot を変更しない。 |
| no secret exposure | `.admin_credentials`、`.github_token`、`.server_config`、`.build_logs`、`.snapshots` を静的配信しない。 |
| setup integration | `docs/details/setup.md` §26 の配置・rollback 条件と矛盾しない。 |
