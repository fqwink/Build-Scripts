# Adlaire CI — Setup 詳細仕様

本ファイルは `ADLAIRE_CI_DETAIL_SPEC.md` から分割した `setup` owner component の詳細仕様である。

本ファイルに、方針、ポリシー、実装状態、正本関係、ロードマップ状態、実装可否の上位判断を記載してはならない。これらは `ADLAIRE_CI_SPEC.md` を正とする。

本ファイルを読む前に、`ADLAIRE_CI_SPEC.md` で実装状態と実装可否を確認し、`ADLAIRE_CI_DETAIL_SPEC.md` §0〜§0j で共通固定値、責務 component、詳細節対応表、リポジトリ内ソース配置を確認する。本ファイルは `setup` owner component の主本文であり、collaborator component の仕様は配置対象、状態初期化、admin 配布、service health、fixture、検証観点として参照する。

本ファイルは、バイナリ配布、配置、systemd、セットアップ、アップデート、リリース成果物検証、Phase 完了判定 fixture 記録を定義する。runner / api / sdk / ui / admin の個別機能本文は各 owner component の詳細仕様ファイルを正とする。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `setup` |
| collaborator component | `runner`、`api`、`statefile`、`admin` |
| 持つ内容 | バイナリ配布、配置、systemd、セットアップ、アップデート、リリース成果物検証。 |
| 持たない内容 | runner / api / sdk / ui / admin の個別機能本文、状態 schema の暗黙変更、外部依存追加。 |

---

## 26. セットアップ・アップデート手順

本節は、Go 版 Adlaire CI のセットアップ手順を定義する。

### §26.1 要件

| 項目 | 要件 |
|------|------|
| Go 版バイナリ | `adlaire-ci-build`、`adlaire-ci-runner`。管理 API 導入時は `adlaire-ci-api` も配置する。 |
| 配布形式 | GitHub Release に添付された OS/arch 別の実行バイナリを標準とする。初期標準は Linux x86_64（`linux-amd64`）。 |
| Go toolchain | 利用環境には不要。リリースバイナリをそのまま配置し、利用環境で `go build` しない。 |
| checksum | Release 添付ファイルごとの SHA-256 checksum を取得し、配置前に必ず検証する。 |
| init システム | systemd（Linux） |
| バージョン管理 | GitHub Releases のタグ付き安定版を使用する。利用環境でリポジトリ checkout を更新経路にしない。 |
| ネットワーク | GitHub API への HTTPS 送信。SSH 転送機能を実装した場合のみデプロイ先への SSH 接続。 |

### §26.2 設定変数

スクリプト内で以下の変数をカスタマイズする。

| 変数 | デフォルト値 | 説明 |
|------|------------|------|
| `INSTALL_DIR` | `/opt/adlaire-builder` | インストール先ディレクトリ |
| `BIN_DIR` | `/usr/local/bin` | Go 版バイナリ配置先 |
| `SERVICE_USER` | `root` | systemd サービスの実行ユーザー |
| `VERSION` | —（必須） | セットアップ・アップデート対象の安定版タグ（例：`V.1.100`） |
| `OS_ARCH` | `linux-amd64` | 取得するリリースバイナリの OS/arch。初期標準は `linux-amd64` のみ |
| `DOWNLOAD_DIR` | `/tmp/adlaire-ci-release-$VERSION` | Release 添付ファイルの一時取得先 |

### §26.2a リリース成果物

セットアップ手順は、以下の Release 添付ファイルを取得対象とする。

| 成果物 | 取得タイミング | 説明 |
|--------|----------------|------|
| `adlaire-ci-build-$OS_ARCH` | 初回セットアップ、アップデート | `components/builder.go` から生成した Markdown → 静的 Web サイトビルドバイナリ。標準配置への移行完了前は、現行実装実体 `build_spec.go` から生成する同等バイナリ。 |
| `adlaire-ci-runner-$OS_ARCH` | 初回セットアップ、アップデート | `components/runner.go` から生成した CI ランナーバイナリ。標準配置への移行完了前は、現行実装実体 `runner.go` から生成する同等バイナリ。 |
| `adlaire-ci-api-$OS_ARCH` | 管理 API 導入手順、管理 API 導入後のアップデート | `components/api.go` から生成した管理 API サーバーバイナリ。 |
| `admin-ui.tar.gz` | 管理 API 導入手順、管理 API 導入後のアップデート | `ADLAIRE_CI_DETAIL_ADMIN_SPEC.md` A1 の管理 UI 配布物。 |
| `SHA256SUMS` | Release 添付ファイル取得時 | Release 添付ファイルの SHA-256 checksum 一覧。 |

Release asset 名は上表の文字列と完全一致させる。`$OS_ARCH` は `linux-amd64` だけを初期標準とし、未知 OS/arch を指定した場合は取得前に `unsupported OS_ARCH: {OS_ARCH}` を stderr へ出力して終了コード `2` とする。`SHA256SUMS` は `"{sha256}  {filename}"` 形式の LF 区切り text とし、対象 filename が 1 回だけ出現することを必須とする。対象行が 0 件または 2 件以上の場合は checksum 検証失敗とする。

### §26.2b セットアップ・アップデート機能単位

セットアップ・アップデート実装は、以下の機能単位に分割する。各機能は前段の出力だけを入力として受け取り、失敗時は後続機能を実行しない。

| 機能 | 入力 | 出力 | 失敗条件 | 失敗時の終了状態 |
|------|------|------|----------|------------------|
| Release asset resolver | `VERSION`、`OS_ARCH`、取得対象成果物名、GitHub Release URL | `DOWNLOAD_DIR` 内の取得済みファイル | `VERSION` / `OS_ARCH` 空、HTTP status 非 2xx、取得ファイル 0 byte | 取得済みファイルを配置せず終了 |
| checksum verifier | `SHA256SUMS`、取得済み成果物 | 検証済み成果物一覧 | `SHA256SUMS` 不在、対象行不在、SHA-256 不一致 | バイナリ配置を実行せず終了 |
| binary installer | 検証済みバイナリ、`BIN_DIR` | `adlaire-ci-build`、`adlaire-ci-runner`、必要時 `adlaire-ci-api` | 入力バイナリ不在、実行権限付与失敗、`install` 失敗 | systemd 変更を実行せず終了 |
| secret initializer | PAT 入力、`INSTALL_DIR` | `.github_token` mode `0600` | PAT 空、書き込み失敗、mode 補正失敗 | systemd 変更を実行せず終了 |
| state initializer | `INSTALL_DIR` | `.last_sha`、必要時 `.build_logs/`、`.snapshots/` | 書き込み失敗、mode 補正失敗 | systemd 変更を実行せず終了 |
| systemd unit writer | unit 内容、`SERVICE_USER`、`INSTALL_DIR`、`BIN_DIR` | `/etc/systemd/system/adlaire-ci.service`、`adlaire-ci.timer`、必要時 `adlaire-ci-api.service` | unit 書き込み失敗、`systemctl daemon-reload` 失敗 | enable/start を実行せず終了 |
| service activator | systemd unit 名 | active な timer / service | `enable --now` 失敗、`is-active` 非 `active` | 直前の journal 確認コマンドを出力して終了 |
| admin UI installer | `admin-ui.tar.gz`、`INSTALL_DIR` | `$INSTALL_DIR/admin/index.html`、`$INSTALL_DIR/admin/adlaire-ci-sdk.js` | archive 不在、checksum 不一致、展開後必須ファイル不在 | API service 起動を実行せず終了 |
| rollback executor | `BACKUP_DIR`、`BIN_DIR`、再起動対象 unit | 旧バイナリ復元済み状態 | 旧バイナリ不在、復元失敗、復元後 restart 失敗 | 自動復旧を継続せず journal 確認対象を出力 |

**セットアップ / アップデート共通終了コード：**

| 終了コード | 条件 |
|------------|------|
| `0` | 全手順成功。 |
| `1` | 取得失敗、checksum 不一致、配置失敗、systemd 操作失敗、権限補正失敗、rollback 失敗。 |
| `2` | 変数不正、unsupported OS/arch、必須入力空、既存 credentials あり、実行前検証不合格。 |

各手順は失敗時に固定文言を stderr へ 1 行以上出力する。secret 値、PAT、token、password、Release URL に埋め込まれた認証情報を stderr/stdout に出してはならない。

**Release asset 取得・検証固定契約：**

セットアップ、管理 API 導入、アップデートはいずれも下表の順序で Release asset を扱う。順序を入れ替えてはならない。取得済みファイルは checksum 検証が成功するまで配置対象として扱わない。

| 手順 | 入力 | 成功条件 | 失敗時 |
|------|------|----------|--------|
| 1. 変数検証 | `VERSION`, `OS_ARCH`, `DOWNLOAD_DIR`, 取得対象 asset 名 | 空値なし、`OS_ARCH=linux-amd64`、`DOWNLOAD_DIR` が `/` でない。 | 終了コード `2`。directory 作成、download、配置を行わない。 |
| 2. download dir 作成 | `DOWNLOAD_DIR` | directory が存在し mode `0755` 以上で書込可能。 | 終了コード `1`。配置、systemd 操作を行わない。 |
| 3. asset 取得 | Release URL、asset 名 | HTTP 2xx、取得ファイル size > 0。 | 終了コード `1`。取得済み未検証ファイルを配置しない。 |
| 4. SHA256SUMS 取得 | Release URL | HTTP 2xx、size > 0、LF text。 | 終了コード `1`。asset を配置しない。 |
| 5. 対象行確認 | `SHA256SUMS`、asset 名 | 対象 filename が 1 回だけ出現する。 | 終了コード `1`。0 件、2 件以上はいずれも checksum failure。 |
| 6. checksum 検証 | 対象 asset、対象 SHA-256 | 実ファイル digest が一致する。 | 終了コード `1`。asset を配置しない。 |
| 7. 実行権限付与前確認 | 検証済み binary asset | 通常ファイルであり、directory / symlink ではない。 | 終了コード `1`。配置しない。 |

`admin-ui.tar.gz` は checksum 検証後に一時展開ディレクトリへ展開する。配布物の中身、必須 file、拒否する archive entry は `ADLAIRE_CI_DETAIL_ADMIN_SPEC.md` A1〜A2 を正とする。検証に失敗した場合は、既存 `$INSTALL_DIR/admin` を変更しない。

**配置・権限固定契約：**

| 対象 | 配置方法 | mode | 失敗時 |
|------|----------|------|--------|
| `adlaire-ci-build` | 検証済み asset を `install -m 0755` で `$BIN_DIR/adlaire-ci-build` へ配置する。 | `0755` | systemd restart を行わない。旧 binary がある場合は保持する。 |
| `adlaire-ci-runner` | 検証済み asset を `install -m 0755` で `$BIN_DIR/adlaire-ci-runner` へ配置する。 | `0755` | systemd restart を行わない。旧 binary がある場合は保持する。 |
| `adlaire-ci-api` | 検証済み asset を `install -m 0755` で `$BIN_DIR/adlaire-ci-api` へ配置する。 | `0755` | API service を restart / start しない。 |
| `.github_token` | 一時ファイルへ書込後、mode `0600`、rename、fsync。 | `0600` | systemd unit を変更しない。secret 平文を stderr/stdout に出さない。 |
| `.last_sha` | `{"sha":""}` + LF を一時ファイルへ書込後、mode `0600`、rename、fsync。 | `0600` | systemd unit を変更しない。 |
| `.admin_credentials` | `adlaire-ci-api --init-credentials` の固定手順で生成する。 | `0600` | API service を enable/start しない。 |
| `admin/` | checksum 検証済み archive を一時 directory へ展開し、必須ファイル確認後に差し替える。 | directory `0755`、file `0644` | 既存 `admin/` を変更しない。 |

通常ファイル配置では symlink を最終配置先として許可しない。既存配置先が symlink の場合は `1` で停止し、symlink の参照先を上書きしてはならない。`BIN_DIR`、`INSTALL_DIR`、`DOWNLOAD_DIR` が同一 path、親子関係で危険な組み合わせ、またはいずれかが `/` の場合は `2` で停止する。

**setup / update 失敗時の状態保持契約：**

| 失敗箇所 | 保持するもの | 変更してよいもの | 禁止事項 |
|----------|--------------|------------------|----------|
| download / checksum | 既存 binary、既存 systemd、既存 state、既存 admin UI | `DOWNLOAD_DIR` 内の取得済みファイル | 未検証 asset の配置、service restart。 |
| binary 配置前 | 既存 binary、既存 service 稼働状態 | `DOWNLOAD_DIR` | systemd unit 書換、state 書換。 |
| binary 配置後 / restart 前 | 配置済み新 binary または rollback 対象旧 binary | rollback executor が対象 binary だけ復元してよい。 | state、history、secret、admin UI の巻き戻し。 |
| runner restart 失敗 | `.github_token`、`.last_sha`、history、snapshot、admin UI | build / runner binary の旧版復元、runner restart 1 回 | API credentials や admin UI の変更。 |
| API setup 失敗 | runner binary、runner timer、runner state | API binary、admin 一時展開 directory | runner timer 停止、`.github_token` 変更。 |
| admin UI 差し替え失敗 | 旧 admin UI、API binary、runner state | admin 一時 directory / backup directory | API restart、credentials 変更。 |
| rollback 失敗 | 現在配置済み binary、state、secret | journal 確認対象の報告 | 追加 rollback 推測、state/history/secret 巻き戻し。 |

**セットアップ / アップデート副作用固定契約：**

セットアップ、管理 API 導入、アップデートは、下表の副作用境界を超えてはならない。実装者判断で部分成功を成功報告したり、secret、state、systemd、admin UI をまとめて巻き戻したりしてはならない。

| 段階 | 変更可能対象 | 成功確定条件 | 失敗時固定動作 |
|------|--------------|--------------|----------------|
| 変数検証 | なし | すべての必須変数が空でなく、危険 path でない。 | 終了コード `2`。directory 作成、download、systemd 操作を行わない。 |
| download | `DOWNLOAD_DIR` 配下だけ | 対象 asset と `SHA256SUMS` を取得し、size > 0。 | 既存 binary、state、secret、systemd、admin UI を変更しない。 |
| checksum | `DOWNLOAD_DIR` 配下だけ | 対象 filename が `SHA256SUMS` に 1 回だけ存在し、SHA-256 が一致する。 | 未検証 asset を配置しない。 |
| binary 配置 | `$BIN_DIR` の対象 binary だけ | 通常ファイルへ `0755` で配置し、`--version` が対象 version を返す。 | systemd restart を行わない。配置済み新 binary は rollback 表に従う。 |
| secret / state 初期化 | 対象 secret / state file だけ | 一時ファイル、mode、rename、fsync、親 directory sync が成功する。 | systemd unit を変更しない。secret 平文を出力しない。 |
| admin UI 展開 | admin 一時 directory、成功時のみ `$INSTALL_DIR/admin` | archive 安全検査、必須ファイル確認、差し替えがすべて成功する。 | 既存 admin UI を維持する。API restart を行わない。 |
| systemd unit 配置 | 対象 unit file だけ | unit 書込、mode、`systemctl daemon-reload` が成功する。 | enable / restart / start を行わない。 |
| service 起動 / 再起動 | 対象 unit だけ | `systemctl is-active` が `active`。api は health check も成功。 | rollback 表に従い、追加推測復旧を行わない。 |
| 最終確認 | なし | §26.3、§26.3b、§26.5 の確認項目がすべて成功。 | 成功報告しない。確認失敗箇所と journal 確認対象を出力する。 |

setup / update 実装は、各段階の開始と成功を stderr または stdout に固定文言で 1 行ずつ出してよいが、PAT、password、session token、API token、Webhook secret、SMTP password、Release URL の credential 部分は出力してはならない。secret file が既に存在する場合は、個別手順で上書きを明記している場合を除き、既存値を保持する。特に `.github_token`、`.admin_credentials`、`.webhook_secret`、`.smtp_secret` は、アップデートで自動上書きしない。

`systemctl daemon-reload` 成功だけではセットアップ成功と扱わない。`enable --now`、`restart`、`is-active`、API 導入時の `/api/health` 確認まで完了して初めて成功とする。確認コマンドが利用環境に存在しない場合は、同等確認を実装 PR の検証で実施し、未確認のまま成功扱いにしない。

### §26.3 Go 版初回セットアップ手順

対象は Go 版の `components/builder.go` と `components/runner.go` から生成した `adlaire-ci-build`、`adlaire-ci-runner`、`adlaire-ci.service`、`adlaire-ci.timer` とする。標準配置への移行完了前は、`components/builder.go` を現行実装実体 `build_spec.go`、`components/runner.go` を現行実装実体 `runner.go` と読み替える。

初回セットアップは以下の停止条件に従う。各手順は直前の手順が成功した場合のみ実行する。失敗時に後続手順を継続してはならない。

| 手順 | 停止条件 | 失敗時の扱い |
|------|----------|--------------|
| 変数検証 | `INSTALL_DIR`、`BIN_DIR`、`VERSION`、`OS_ARCH`、`DOWNLOAD_DIR` が空、`INSTALL_DIR` が `/`、`BIN_DIR` が `/`、`DOWNLOAD_DIR` が `/` | 何も変更せず終了する。 |
| バイナリ取得 | Release バイナリまたは `SHA256SUMS` の取得に失敗、checksum 検証に失敗 | バイナリを配置せず、systemd 設定を変更せず終了する。 |
| バイナリ配置 | checksum 検証済みバイナリが存在しない、または `install` が失敗 | systemd 設定を変更せず終了する。 |
| secret 保存 | PAT が空 | `.github_token` を作成せず終了する。 |
| systemd 配置 | unit ファイル生成または `systemctl daemon-reload` が失敗 | timer を enable せず終了する。 |
| 起動確認 | `systemctl is-active adlaire-ci.timer` が `active` でない | 失敗として扱い、直前のログ確認コマンドを表示する。 |

初回セットアップが中断した場合、作成済みの通常ディレクトリと展開済みリリース資産は自動削除しない。秘密情報ファイルを作成した後に失敗した場合は、`.github_token` の mode が `0600` であることを確認し、mode 補正に失敗した場合はその場で停止する。

**初回セットアップ後の固定確認：**

| 確認 | コマンド | 合格条件 |
|------|----------|----------|
| build binary | `$BIN_DIR/adlaire-ci-build --version` | exit `0`、stdout が `adlaire-ci-build ADLAIRE_CI_SPEC` を含む。 |
| runner binary | `$BIN_DIR/adlaire-ci-runner --version` | exit `0`、stdout が `adlaire-ci-runner ADLAIRE_CI_SPEC` を含む。 |
| PAT file | `stat -c '%a' "$INSTALL_DIR/.github_token"` | `600`。 |
| SHA cache | `cat "$INSTALL_DIR/.last_sha"` | `{"sha":""}` + LF。 |
| timer | `systemctl is-active adlaire-ci.timer` | `active`。 |

確認のいずれかが失敗した場合、セットアップは失敗扱いとする。ただし自動削除や状態ファイル巻き戻しは行わない。

```bash
# ── 変数設定 ──────────────────────────────────────────
INSTALL_DIR="/opt/adlaire-builder"
BIN_DIR="/usr/local/bin"
VERSION="V.1.100"
OS_ARCH="linux-amd64"
DOWNLOAD_DIR="/tmp/adlaire-ci-release-$VERSION"
SERVICE_USER="root"

# ── 1. インストール先作成 ─────────────────────────────
mkdir -p "$INSTALL_DIR"
chmod 0755 "$INSTALL_DIR"

# ── 2. Release バイナリ取得・checksum 検証 ────────────
mkdir -p "$DOWNLOAD_DIR"
cd "$DOWNLOAD_DIR"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$VERSION/adlaire-ci-build-$OS_ARCH"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$VERSION/adlaire-ci-runner-$OS_ARCH"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$VERSION/SHA256SUMS"
grep "  adlaire-ci-build-$OS_ARCH$" SHA256SUMS | sha256sum -c -
grep "  adlaire-ci-runner-$OS_ARCH$" SHA256SUMS | sha256sum -c -

# ── 3. Go 版バイナリ配置 ──────────────────────────────
install -m 0755 "adlaire-ci-build-$OS_ARCH"  "$BIN_DIR/adlaire-ci-build"
install -m 0755 "adlaire-ci-runner-$OS_ARCH" "$BIN_DIR/adlaire-ci-runner"

# ── 4. GitHub PAT 保存 ────────────────────────────────
printf '%s\n' "<PAT>" > "$INSTALL_DIR/.github_token"
chmod 600 "$INSTALL_DIR/.github_token"

# ── 5. SHA キャッシュ初期化 ───────────────────────────
printf '%s\n' '{"sha":""}' > "$INSTALL_DIR/.last_sha"
chmod 600 "$INSTALL_DIR/.last_sha"

# ── 6. systemd サービスファイル配置 ───────────────────
# §26.4.1 のファイル内容を /etc/systemd/system/ に配置した上で:
systemctl daemon-reload

# ── 7. タイマー有効化・起動 ───────────────────────────
systemctl enable --now adlaire-ci.timer

# ── 8. 起動確認 ───────────────────────────────────────
systemctl status adlaire-ci.timer
```

Go 版初回セットアップでは以下を実行しない。

| 対象 | 理由 |
|------|------|
| `/usr/local/bin/adlaire-ci-api --init-credentials --state-dir "$INSTALL_DIR"` | 初回セットアップ対象は runner と build バイナリに限定し、API 認証情報生成は §26.3b で実行する。 |
| `systemctl enable --now adlaire-ci-api` | API service は §26.3b の API バイナリ配置、認証情報生成、unit 配置がすべて成功した後にのみ起動する。 |
| `.build_logs/` 作成 | runner 初期導入ではビルド実行時に必要な状態だけを初期化し、api が参照する履歴ディレクトリは §26.3b で作成する。 |
| `.snapshots/` 作成 | snapshot 参照・rollback API と組み合わせて使うため、§26.3b の管理 API 導入時に作成する。 |

### §26.3b 管理 API 導入後の追加セットアップ手順

`api`、`ui`、`sdk` を実装した後にのみ本手順を実行する。

管理 API 導入手順は、runner の既存稼働状態を壊してはならない。`adlaire-ci-api` の配置、認証情報生成、systemd enable のいずれかが失敗した場合でも、`adlaire-ci.timer` は停止しない。`.admin_credentials` が既に存在する場合は `--init-credentials` を再実行せず、既存 credentials を維持する。

管理 API 導入手順は以下の停止条件に従う。

| 手順 | 停止条件 | 失敗時の扱い |
|------|----------|--------------|
| ディレクトリ作成 | `$INSTALL_DIR/.build_logs`、`$INSTALL_DIR/.snapshots`、`$INSTALL_DIR/admin` の作成に失敗 | runner timer を変更せず終了する。 |
| API バイナリ取得 | `adlaire-ci-api-$OS_ARCH` または `SHA256SUMS` の取得、checksum 検証に失敗 | API バイナリを配置せず終了する。 |
| 管理 UI 取得 | `admin-ui.tar.gz` の取得、checksum 検証、展開に失敗 | API service を起動せず終了する。 |
| 管理 UI 必須ファイル確認 | `$INSTALL_DIR/admin/index.html` または `$INSTALL_DIR/admin/adlaire-ci-sdk.js` が存在しない | API service を起動せず終了する。 |
| API バイナリ配置 | checksum 検証済み API バイナリ不在、または `install` 失敗 | API service を起動せず終了する。 |
| 認証情報生成 | `.admin_credentials` 新規生成に失敗。ただし既存ファイルがある場合は成功扱い | API service を起動せず終了する。 |
| systemd 配置 | unit 書き込みまたは `systemctl daemon-reload` 失敗 | API service を enable/start せず終了する。 |
| 起動確認 | `systemctl is-active adlaire-ci-api` が `active` でない | runner timer を停止せず、API の journal 確認コマンドを出力して終了する。 |

**管理 API 導入後の固定確認：**

| 確認 | コマンド | 合格条件 |
|------|----------|----------|
| API binary | `$BIN_DIR/adlaire-ci-api --version` | exit `0`、stdout が `adlaire-ci-api ADLAIRE_CI_SPEC` を含む。 |
| credentials | `stat -c '%a' "$INSTALL_DIR/.admin_credentials"` | `600`。 |
| admin UI | `test -f "$INSTALL_DIR/admin/index.html"` / `test -f "$INSTALL_DIR/admin/adlaire-ci-sdk.js"` | 両方成功。 |
| API service | `systemctl is-active adlaire-ci-api` | `active`。 |
| local health | `curl -fsS http://127.0.0.1:8765/api/health` | HTTP `200`、JSON object。 |

`curl` が利用できない環境では、Go 実装 PR の検証で `net/http` client または同等のローカル HTTP 確認を行う。未確認のまま API 導入完了扱いにしてはならない。

```bash
# ── 1. 拡張用ディレクトリ作成 ─────────────────────────
mkdir -p "$INSTALL_DIR/.build_logs"
mkdir -p "$INSTALL_DIR/.snapshots"
mkdir -p "$INSTALL_DIR/admin"

# ── 2. Release バイナリ取得・checksum 検証 ────────────
mkdir -p "$DOWNLOAD_DIR"
cd "$DOWNLOAD_DIR"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$VERSION/adlaire-ci-api-$OS_ARCH"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$VERSION/admin-ui.tar.gz"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$VERSION/SHA256SUMS"
grep "  adlaire-ci-api-$OS_ARCH$" SHA256SUMS | sha256sum -c -
grep "  admin-ui.tar.gz$" SHA256SUMS | sha256sum -c -

# ── 3. Go 版 API バイナリ配置 ─────────────────────────
install -m 0755 "adlaire-ci-api-$OS_ARCH" "$BIN_DIR/adlaire-ci-api"

# ── 4. 管理 UI 配布物展開 ────────────────────────────
tar -xzf admin-ui.tar.gz -C "$INSTALL_DIR/admin"
test -f "$INSTALL_DIR/admin/index.html"
test -f "$INSTALL_DIR/admin/adlaire-ci-sdk.js"

# ── 5. 初期認証情報生成（初期パスワード: admin）────────
/usr/local/bin/adlaire-ci-api --init-credentials --state-dir "$INSTALL_DIR"
chmod 600 "$INSTALL_DIR/.admin_credentials"

# ── 6. 管理 API systemd サービス配置 ─────────────────
# §26.4.2 のファイル内容を /etc/systemd/system/adlaire-ci-api.service に配置した上で:
systemctl daemon-reload

# ── 7. サービス有効化・起動 ───────────────────────────
systemctl enable --now adlaire-ci-api

# ── 8. 起動確認 ───────────────────────────────────────
systemctl status adlaire-ci-api
```

### §26.4 systemd サービスファイル

#### §26.4.1 Go 版 runner の systemd ファイル

**`/etc/systemd/system/adlaire-ci.service`**（`components/runner.go`。標準配置への移行完了前の現行実装実体は `runner.go`）：

```ini
[Unit]
Description=Adlaire CI Runner

[Service]
Type=oneshot
User=root
WorkingDirectory=/opt/adlaire-builder
ExecStart=/usr/local/bin/adlaire-ci-runner --state-dir /opt/adlaire-builder
```

**`/etc/systemd/system/adlaire-ci.timer`**（`components/runner.go` 定期起動タイマー。標準配置への移行完了前の現行実装実体は `runner.go`）：

```ini
[Unit]
Description=Adlaire CI Runner Timer

[Timer]
OnBootSec=1min
OnUnitActiveSec=5min
Unit=adlaire-ci.service

[Install]
WantedBy=timers.target
```

#### §26.4.2 管理 API 導入後の systemd ファイル

**`/etc/systemd/system/adlaire-ci-api.service`**（`components/api.go`）：

```ini
[Unit]
Description=Adlaire CI API Server
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/adlaire-builder
ExecStart=/usr/local/bin/adlaire-ci-api --addr 127.0.0.1:8765 --state-dir /opt/adlaire-builder
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

`User` / `WorkingDirectory` / `ExecStart` のパスは §26.2 の設定変数に合わせて変更する。

systemd unit は上記キー以外を初期標準で追加しない。`Environment=`、`EnvironmentFile=`、`ExecStartPre=`、`ExecStartPost=` を追加する場合は、先に本節へ対象変数、secret 扱い、失敗時挙動を定義する。API service は `127.0.0.1:8765` bind を標準とし、外部公開 bind は本ファイルで未定義のため設定しない。

### §26.5 アップデート手順

`git pull`、利用環境での `go build`、開発ブランチ checkout は使用しない。タグ付き安定版のリリースバイナリを配置し、サービスを再起動する。管理 API を導入していない構成では、管理 API サービスは再起動対象に含めない。

アップデートは以下の順序で実行し、途中失敗時は表の rollback 条件に従う。

| 手順 | 成功条件 | 失敗時 rollback / 停止条件 |
|------|----------|-----------------------------|
| 現在版記録 | 既存バイナリを退避し、退避先パスを保持する。 | 更新を開始しない。 |
| バイナリ取得 | 新 Release バイナリと `SHA256SUMS` の取得、checksum 検証が成功する。 | 退避済み旧バイナリを維持して終了する。 |
| バイナリ更新 | checksum 検証済みの新バイナリを `install -m 0755` で配置できる。 | 退避済み旧バイナリを元へ戻し、サービスを再起動しない。 |
| runner 再起動 | `systemctl restart adlaire-ci.timer` と `systemctl is-active adlaire-ci.timer` が成功する。 | 旧バイナリを戻し、再度 `systemctl restart adlaire-ci.timer` を 1 回だけ実行する。 |
| API 再起動 | API 導入済みの場合のみ `systemctl restart adlaire-ci-api` と `systemctl is-active adlaire-ci-api` が成功する。 | 旧バイナリを戻し、runner と API の再起動を 1 回だけ実行する。 |
| 管理 UI 更新 | API 導入済みの場合のみ `admin-ui.tar.gz` の取得、checksum 検証、一時ディレクトリへの展開、必須ファイル確認、旧 `admin/` との差し替えが成功する。 | 旧 `admin/` を維持または退避先から復元し、API 再起動を実行しない。 |

rollback 後も service が active にならない場合は、自動復旧を継続せず、`journalctl -u adlaire-ci.service -n 100`、API 導入済みなら `journalctl -u adlaire-ci-api -n 100` を確認対象として報告する。rollback はバイナリ差し戻しと service restart のみを行い、状態ファイル、履歴、ログ、secret を巻き戻してはならない。

**アップデート rollback 固定契約：**

| 失敗箇所 | rollback 対象 | rollback 後に実行する確認 | 禁止事項 |
|----------|---------------|----------------------------|----------|
| checksum 検証前 | なし | 旧 service active 確認のみ | 取得済み未検証ファイルを配置しない。 |
| build / runner 配置失敗 | 配置に成功した新バイナリだけ旧版へ戻す。 | `adlaire-ci-build --version`、`adlaire-ci-runner --version` | systemd restart しない。 |
| runner restart 失敗 | build / runner 旧版復元 | `systemctl is-active adlaire-ci.timer` | state、history、secret を戻さない。 |
| API binary 配置失敗 | API 旧版復元。runner は戻さない。 | `adlaire-ci-api --version` | runner service を restart しない。 |
| API restart 失敗 | API 旧版復元 | `systemctl is-active adlaire-ci-api` | `.admin_credentials`、admin UI を戻さない。ただし UI 更新前に失敗した場合。 |
| admin UI 展開失敗 | 旧 `admin/` 維持 | 必須ファイル確認 | API restart しない。 |
| admin UI 差し替え後 API restart 失敗 | 旧 `admin/` 復元、API 旧版復元 | API service active 確認 | runner state を戻さない。 |

rollback は 1 回だけ実行する。rollback 自体が失敗した場合は、追加の推測復旧を行わず、失敗箇所、退避先、現在配置済みファイル、journal 確認コマンドを報告対象として固定する。

**アップデート実行判定固定契約：**

アップデート実装は、以下の判定順で対象を決定する。判定結果を実装者判断で短絡、統合、または省略してはならない。

1. 既存 `$BIN_DIR/adlaire-ci-build` と `$BIN_DIR/adlaire-ci-runner` の存在を確認する。どちらかが不在の場合は終了コード `2` とし、更新を開始しない。
2. API 導入済み判定は `$BIN_DIR/adlaire-ci-api` が通常ファイルとして存在し、`systemctl is-enabled adlaire-ci-api` が `enabled` または `static` を返す場合だけ `true` とする。
3. API 導入済みでない場合、`adlaire-ci-api-$OS_ARCH` と `admin-ui.tar.gz` は取得しない。
4. API 導入済みの場合、build / runner / api binary と admin UI を同じ `NEW_VERSION` の asset から取得する。version 混在は禁止する。
5. すべての対象 asset の checksum 検証が成功するまで、既存 binary、既存 admin UI、systemd unit を変更しない。
6. binary 配置後の version 確認に失敗した場合は、その binary を配置失敗として rollback 対象に含める。
7. runner restart が失敗した場合、API restart と admin UI 更新へ進まない。
8. API restart が失敗した場合、admin UI 更新へ進まない。
9. admin UI 差し替え後に API restart が失敗した場合、admin UI と API binary だけ rollback 対象とし、build / runner binary と runner timer は戻さない。

**アップデート後確認固定契約：**

| 確認 | API 未導入 | API 導入済み |
|------|------------|--------------|
| binary version | `adlaire-ci-build --version`、`adlaire-ci-runner --version` が `NEW_VERSION` を含む。 | 左記に加え `adlaire-ci-api --version` が `NEW_VERSION` を含む。 |
| service | `systemctl is-active adlaire-ci.timer` が `active`。 | 左記に加え `systemctl is-active adlaire-ci-api` が `active`。 |
| admin UI | 確認しない。 | `$INSTALL_DIR/admin/index.html` と `$INSTALL_DIR/admin/adlaire-ci-sdk.js` が存在する。 |
| local API | 確認しない。 | `GET /api/health` が HTTP `200` JSON object を返す。 |
| state preservation | `.github_token`、`.last_sha`、`.build_state`、`.build_history` の mtime と内容が更新対象操作と無関係に変わっていない。 | 左記に加え `.admin_credentials` が存在する場合は mode `600` と内容が保持される。 |

確認失敗時はアップデート失敗として扱う。binary 配置や restart が成功していても、確認失敗を成功報告してはならない。local API 確認で `curl` がない場合は Go 実装 PR の検証で `net/http` client による同等確認を行い、未確認のまま合格扱いにしない。

```bash
# ── 変数設定 ──────────────────────────────────────────
BIN_DIR="/usr/local/bin"
NEW_VERSION="V.2.102"
OS_ARCH="linux-amd64"
DOWNLOAD_DIR="/tmp/adlaire-ci-release-$NEW_VERSION"
BACKUP_DIR="/tmp/adlaire-ci-bin-backup-${NEW_VERSION}"
ADMIN_BACKUP_DIR="/tmp/adlaire-ci-admin-backup-${NEW_VERSION}"
ADMIN_TMP_DIR="/tmp/adlaire-ci-admin-new-${NEW_VERSION}"

# ── 1. 既存バイナリ退避 ──────────────────────────────
mkdir -p "$BACKUP_DIR"
cp "$BIN_DIR/adlaire-ci-build"  "$BACKUP_DIR/adlaire-ci-build"
cp "$BIN_DIR/adlaire-ci-runner" "$BACKUP_DIR/adlaire-ci-runner"
if [ -d "/opt/adlaire-builder/admin" ]; then
  cp -a "/opt/adlaire-builder/admin" "$ADMIN_BACKUP_DIR"
fi

# ── 2. Release バイナリ取得・checksum 検証 ────────────
mkdir -p "$DOWNLOAD_DIR"
cd "$DOWNLOAD_DIR"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$NEW_VERSION/adlaire-ci-build-$OS_ARCH"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$NEW_VERSION/adlaire-ci-runner-$OS_ARCH"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$NEW_VERSION/SHA256SUMS"
grep "  adlaire-ci-build-$OS_ARCH$" SHA256SUMS | sha256sum -c -
grep "  adlaire-ci-runner-$OS_ARCH$" SHA256SUMS | sha256sum -c -

# ── 3. Go 版バイナリ更新 ──────────────────────────────
install -m 0755 "adlaire-ci-build-$OS_ARCH"  "$BIN_DIR/adlaire-ci-build"
install -m 0755 "adlaire-ci-runner-$OS_ARCH" "$BIN_DIR/adlaire-ci-runner"

# ── 4. runner サービス再起動 ─────────────────────────
systemctl restart adlaire-ci.timer

# ── 5. 起動確認 ───────────────────────────────────────
systemctl status adlaire-ci.timer
```

管理 API 導入後は、追加で `adlaire-ci-api` を再起動する。

```bash
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$NEW_VERSION/adlaire-ci-api-$OS_ARCH"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$NEW_VERSION/admin-ui.tar.gz"
grep "  adlaire-ci-api-$OS_ARCH$" SHA256SUMS | sha256sum -c -
grep "  admin-ui.tar.gz$" SHA256SUMS | sha256sum -c -
install -m 0755 "adlaire-ci-api-$OS_ARCH" "$BIN_DIR/adlaire-ci-api"
mkdir -p "$ADMIN_TMP_DIR"
tar -xzf admin-ui.tar.gz -C "$ADMIN_TMP_DIR"
test -f "$ADMIN_TMP_DIR/index.html"
test -f "$ADMIN_TMP_DIR/adlaire-ci-sdk.js"
if [ -d "/opt/adlaire-builder/admin" ]; then
  mv "/opt/adlaire-builder/admin" "$ADMIN_BACKUP_DIR"
fi
mv "$ADMIN_TMP_DIR" "/opt/adlaire-builder/admin"
systemctl restart adlaire-ci-api
systemctl status adlaire-ci-api
```

### §26.6 サービス操作リファレンス

#### Go 版 runner

| 操作 | コマンド |
|------|---------|
| 状態確認 | `systemctl status adlaire-ci.timer` |
| 起動 | `systemctl start adlaire-ci.timer` |
| 停止 | `systemctl stop adlaire-ci.timer` |
| 再起動 | `systemctl restart adlaire-ci.timer` |
| ログ確認（runner） | `journalctl -u adlaire-ci.service -f` |

#### 管理 API 導入後

| 操作 | コマンド |
|------|---------|
| 状態確認 | `systemctl status adlaire-ci-api` |
| 起動 | `systemctl start adlaire-ci-api` |
| 停止 | `systemctl stop adlaire-ci-api` |
| 再起動 | `systemctl restart adlaire-ci-api` |
| ログ確認（API） | `journalctl -u adlaire-ci-api -f` |

### §26.7 実装受け入れ条件

実装 PR は、下表の受け入れ条件をすべて満たすまで完了扱いにしてはならない。実装対象外のコンポーネントは「未実装」として明記し、合格扱いにしない。

| 対象 | 必須コマンド / 確認 | 合格条件 |
|------|---------------------|----------|
| Go 共通 | `gofmt -l <実装対象Goファイル>` | 実装対象 Go ファイルが存在する場合、出力が空。未作成ファイルはコマンド対象に含めない。 |
| Go test | `go test ./...` | Go module が存在する場合に成功する。Go module が存在しない場合は、その理由を実装完了報告に明記する。 |
| build script | `adlaire-ci-build --src <sample.md> --out <tmp.html>` | exit code `0`、HTML 出力あり、`[REPORT]` の `status` が `success`。 |
| runner | `adlaire-ci-runner --state-dir <tmp-state>` | 必須 secret 未設定時の exit code / ERROR log が §12 と一致し、`.build_lock` が残らない。 |
| API | `POST /api/login`、`GET /api/status`、未知 path、body 禁止 endpoint、JSON 不正、認証なし | §22.0 / §22.0e の status code と body に一致する。 |
| SDK | browser runtime で `login()`、`getStatus()`、`streamBuild()`、HTTP error、timeout を確認する。 | `AdlaireCIError`、`StreamHandle`、token 破棄、timeout が §23 と一致する。 |
| UI | login、manual build、SSE 表示、config 保存、token 発行、logout を確認する。 | §24 の DOM id、disabled、成功表示、失敗表示、再取得、秘密情報消去に一致する。 |
| setup | §26.3 または §26.3b の手順を fresh 環境で実行する。 | unit 配置、権限、`systemctl is-active`、secret mode が仕様どおり。 |
| update | §26.5 の手順を前版バイナリから新 tag のリリースバイナリへ実行する。 | 旧バイナリ退避、新バイナリ配置、restart、失敗時 rollback 条件が仕様どおり。 |
| security | secret 値を含む入力後、stdout、stderr、journal、API response、UI 表示を確認する。 | PAT、Webhook Secret、SMTP password、session token、API token 本体が平文で出ない。 |

**認証 / セットアップ fixture 固定：**

| fixture | 入力 | 合格条件 |
|---------|------|----------|
| credentials init success | 空 state dir で `adlaire-ci-api --init-credentials --state-dir <abs>` | `.admin_credentials` mode `600`、schema 全 key、`must_change=true`、stdout 固定文言。 |
| credentials init existing | `.admin_credentials` 既存 | exit `2`、stderr `credentials already exist`、既存ファイル差分なし。 |
| login success | 初期 password `admin` | `must_change:"prompt"`、TOTP 無効時 token 発行、hash/salt 非表示、`.access_log` 成功行。 |
| login failure lock | password 連続 10 回失敗 | 10 回目後、10 分間 `429 {"error":"Too many attempts"}`、password 詳細非表示。 |
| change password | current 正、新 password 有効 | `.admin_credentials` の salt/hash 更新、`must_change=false`、現 session 以外破棄。 |
| session restart | token 発行後に API process restart | 旧 token は `401`、session file は存在しない。 |
| setup checksum mismatch | Release asset と `SHA256SUMS` 不一致 | バイナリ配置なし、systemd 変更なし、終了コード `1`。 |
| update restart failure | 新バイナリ配置後に service restart 失敗 | 旧バイナリ復元を 1 回だけ行い、state/history/secret は巻き戻さない。 |

**セットアップ / アップデート fixture 固定：**

| fixture | 入力 | 合格条件 |
|---------|------|----------|
| setup variable invalid | `INSTALL_DIR=/`、`BIN_DIR=/`、または `OS_ARCH=darwin-arm64` | 終了コード `2`、download / directory 作成 / 配置なし。 |
| setup download failure | build binary 取得が HTTP `404` | 終了コード `1`、runner binary 取得済みでも配置なし、systemd 変更なし。 |
| setup checksum duplicate | `SHA256SUMS` に同一 asset 行が 2 件 | 終了コード `1`、checksum failure、配置なし。 |
| setup symlink target | `$BIN_DIR/adlaire-ci-build` が symlink | 終了コード `1`、symlink 参照先を上書きしない。 |
| setup pat empty | PAT 入力が空 | 終了コード `2`、`.github_token` 作成なし、systemd 変更なし。 |
| setup success | 正常 asset、正常 PAT、fresh 環境 | binary mode `755`、`.github_token`/`.last_sha` mode `600`、timer active、固定確認すべて成功。 |
| api setup credentials existing | `.admin_credentials` 既存で API 導入 | `--init-credentials` を再実行せず既存 credentials を保持し、API service 起動確認まで進む。 |
| api setup admin archive unsafe | `admin-ui.tar.gz` に `../x` または symlink entry | 終了コード `1`、既存 admin UI 維持、API service start なし。 |
| api setup health failure | API service active だが `/api/health` が非 200 | セットアップ失敗扱い、runner timer は停止しない、journal 確認対象を出力。 |
| update api absent | `adlaire-ci-api` 未導入環境で update | build / runner asset だけ取得・検証・配置し、API asset と admin UI を取得しない。 |
| update checksum before change | update 対象 asset の checksum 不一致 | 既存 binary、admin UI、systemd、state に差分なし。 |
| update runner restart failure | runner restart fake failure | build / runner 旧版復元を 1 回だけ実行し、API restart と admin UI 更新へ進まない。 |
| update api restart failure | API restart fake failure | API binary 旧版復元、runner は戻さない、admin UI 更新へ進まない。 |
| update admin restart failure | admin UI 差し替え後の API restart fake failure | 旧 admin UI と API 旧版を復元し、runner state / history / secret は巻き戻さない。 |
| update rollback failure | 旧 binary 不在で rollback executor 失敗 | 追加復旧せず、失敗箇所、退避先、現在配置済みファイル、journal 確認コマンドを報告対象にする。 |

**運用 API fixture 固定：**

| fixture | 入力 | 合格条件 |
|---------|------|----------|
| snapshot save and prune | `snapshots_keep=2` で build success を 3 回実行 | 最新 2 世代だけ残り、各 snapshot に `site.tar.gz` と `meta.json` が存在する。 |
| snapshot rollback running | `.build_state.running=true` で `POST /api/history/{id}/rollback` | `409 {"error":"Build is running"}`、queue 追加なし、history 追記なし。 |
| maintenance enable no-op | 同一 reason で enable を 2 回実行 | 2 回目は `No changes`、`.config_log` 追記なし。 |
| access-control deny | allow に接続元以外を設定して API 呼び出し | 認証前に `403 {"error":"Forbidden"}`、password/token 検証なし。 |
| hook pre abort | `pre` hook が exit `1`、`abort_on_failure=true` | pipeline 未実行、build status `hook_error`、hook log 保存。 |
| alert duplicate | 同一 alert rule を 2 回作成 | 2 回目は `409 {"error":"Conflict"}`、`.alert_rules` 差分なし。 |
| tag rule invalid | 破損 condition を含む `.tag_rules` で build | build failure、ERROR log `TAG_RULE_INVALID`、SHA cache 更新なし。 |
| verify-output no history | 成功履歴なしで `POST /api/verify-output` | `404 {"error":"Not Found"}`、状態ファイル変更なし。 |
| pipeline config reserved arg | `extra_args:["--src","x"]` | `422`、`.pipeline_config` 差分なし。 |
| notes no-op | 同一 content を 2 回保存 | 2 回目は `No changes`、`.config_log` 追記なし。 |
| smtp secret mask | password 付き `POST /api/smtp-config` 後に GET / backup / log 確認 | password 本体は返らず、mask または `password_set:true` だけ表示。 |
| queue disabled | `queue_max_size=0`、build running 中に `POST /api/build` | `429 {"error":"queue_full"}`、`.build_state.queued` は空。 |
| dashboard duplicate widget | widgets に重複 id を指定 | `422`、`.dashboard_layout` 差分なし。 |

**§22〜§26 api / sdk / ui / 認証 / セットアップ 実装完全性固定契約：**

§22〜§26 のコンポーネントは、owner component の各節本文、endpoint 表、SDK method 表、UI 操作契約、fixture を主本文とし、下表を横断受け入れ確認として満たした場合だけ実装完了とする。下表は既存機能の詳細実装を確認する表であり、未定義 endpoint、未定義 UI、未定義認証方式、将来計画機能を追加する根拠にしてはならない。

| 節 | 機能 | 入力 | 出力 | 状態ファイル / 外部副作用 | 失敗時副作用 | 必須 fixture |
|----|------|------|------|---------------------------|--------------|--------------|
| §22.0 | API 共通 | HTTP method/path/header/body、remote addr。 | 固定 status、固定 error body、security header。 | `.api_access_log` 以外は endpoint 契約に従う。 | path/method/body/JSON/auth/scope/validation 失敗時は endpoint 固有処理を開始しない。 | unknown path、method mismatch、body 禁止、JSON 不正、401、403、422、500 mask。 |
| §22.0a〜§22.0c | 状態ファイル / schema | state dir、JSON / JSON Lines / text state。 | typed adapter result、固定初期値、破損時 error。 | atomic write、lock、chmod、fsync、corrupt backup。 | read-only API は状態を修復しない。write 失敗は target を部分更新しない。 | corrupt JSON、unknown key、nullable 違反、lock timeout、chmod failure、GET no write。 |
| §22.0d〜§22.0e | endpoint 契約 | endpoint ごとの request/query/body/path。 | endpoint ごとの response、sdk / ui 対応。 | Read / Write 列に明記された状態だけ扱う。 | 個別 status と共通 error 優先順位に従う。未定義 endpoint を追加しない。 | P0/P1、P2〜P5、pagination、SSE、binary、no-op、partial failure。 |
| §23 | SDK | public method 引数、token、fake fetch response。 | Promise return、`AdlaireCIError`、`StreamHandle`、Blob。 | token は memory のみ。DOM / state file / storage を変更しない。 | `401` だけ token 破棄。`403`、`429`、`500`、network、timeout は token 維持。 | request shape、error shape、timeout、invalid JSON、invalid SSE、body 禁止、401 purge。 |
| §24 | 標準管理 UI | DOM event、form value、SDK return/error。 | DOM 表示、disabled/loading、success/error、secret 消去。 | SDK method だけを呼ぶ。直接 API、状態ファイル、systemd を触らない。 | API 成功前に確定表示しない。失敗時は secret を消し、非 secret 入力は保持する。 | login/TOTP、manual build、stream、refresh failure、secret clearing、disabled priority、direct fetch absence。 |
| §25 | 認証 | password、TOTP code、session token、API token。 | session token、ticket、auth error、access/audit log。 | `.admin_credentials`、`.totp_secret`、memory session/ticket、logs。 | token/ticket は必要ログ成功まで返さない。hash/salt/secret/token 本体を保存しない。 | init、login success/failure、lock、change password、session restart、TOTP replay、audit failure。 |
| §26 | setup/update | release asset、checksum、INSTALL_DIR、BIN_DIR、systemd。 | binary 配置、admin UI、unit、service active、health。 | 検証済み asset だけ配置。secret/state/systemd は段階順に変更する。 | checksum/download/unsafe archive/restart 失敗で後続段階に進まない。rollback は定義範囲だけ 1 回。 | setup invalid、download failure、checksum duplicate、symlink target、API setup、update rollback、health failure。 |

**§22〜§26 横断失敗時副作用固定契約：**

| ケース | 固定結果 |
|--------|----------|
| API validation failure | 状態ファイル、外部 API、systemd、runner、hook、通知を変更しない。`.api_access_log` だけ通常記録対象とする。 |
| API write success / log failure | 個別節が巻き戻しを明記していない限り、保存済み状態は巻き戻さず `500` を返す。 |
| SDK network / timeout | `AdlaireCIError(status=0)` とし、自動 retry、自動 refresh、token 破棄を行わない。 |
| UI refresh failure after success | 操作成功は保持し、再取得失敗だけ panel error に表示する。同じ変更 API を自動再実行しない。 |
| auth log failure before token response | session token、login ticket、API token 本体を response しない。 |
| setup partial failure | 既存 binary、state、secret、admin UI、systemd を、表で許可した対象以外は変更しない。 |
| update rollback failure | 追加推測復旧を行わず、失敗箇所、退避先、現在配置済みファイル、journal 確認対象を報告する。 |

**§22〜§26 実装前・実装後確認固定契約：**

| 段階 | 確認 | 合格条件 |
|------|------|----------|
| 実装前 | endpoint / SDK / UI 対応 | §22.0e の API、§23 の SDK method、§24 の UI 操作が同一機能でそろっている。欠落時は先に仕様改訂する。 |
| 実装前 | state schema | 使用する状態ファイルが §22.0a / §22.0c にあり、型、初期値、破損時処理、更新責務が定義済み。 |
| 実装前 | secret handling | secret 値の保存先、mask、response 禁止、log 禁止、UI 消去条件が定義済み。 |
| 実装後 | common error | unknown path、method mismatch、body 禁止、JSON 不正、401、403、422、500 が固定 body と一致する。 |
| 実装後 | state side effect | 成功、validation failure、conflict、write failure、log failure の状態差分が fixture expected と一致する。 |
| 実装後 | client behavior | SDK error、UI disabled、success/error、refresh、secret 消去、direct fetch 不在が固定契約どおり。 |
| 実装後 | setup/update | checksum、unsafe archive、restart failure、rollback failure、health failure が §26 fixture と一致する。 |

Phase 別の実装受け入れ条件は以下とする。

| Phase | 必須確認 | 合格条件 |
|-------|----------|----------|
| Phase 1 | §0g.1、§8a、§26.7 build script | `adlaire-ci-build` の CLI、出力構造、report、fixture、冪等性が固定され、Phase 2 が追加判断なしに呼び出せる。 |
| Phase 2 | §0g.2、§15a、§26.7 runner | runner 状態ファイル、履歴、ログ、pending queue、snapshot、通知、lock、終了コードが固定され、Phase 3 が状態ファイルを追加判断なしに読める。 |
| Phase 3 | §0g.3、§22.0f P0/P1、§25、§26.7 API | API 共通契約、認証、session、P0/P1 endpoint、error body、状態ファイル read/write が固定される。 |
| Phase 4 | §0g.4、§22.0f P2〜P5、§22.0e | 全 endpoint 契約が SDK 実装に渡せる粒度で固定され、P0/P1 の互換を壊していない。 |
| Phase 5 | §0g.5、§23、§26.7 SDK | 全 SDK method が固定済み endpoint 契約に一致し、Phase 6 が UI 実装に使える error / stream / return 契約を持つ。 |
| Phase 6 | §0g.6、§24、§26.7 UI/security | 全 UI 操作が SDK 経由で成立し、成功/失敗/disabled/loading/secret 消去が仕様どおりである。 |

Phase 別の実装検証記録は以下の単位で行う。

| Phase | 必須記録 | 失敗時の扱い |
|-------|----------|--------------|
| Phase 1 | CLI 引数、fixture A〜D、生成物一覧、`[REPORT]`、冪等性、strict / non-strict の結果。 | `builder` を完了扱いにせず、Phase 2 着手禁止。 |
| Phase 2 | secret 不足、lock、GitHub fake、pipeline fake、deploy fake、snapshot、notify、状態ファイル schema の結果。 | `runner` を完了扱いにせず、Phase 3 着手禁止。 |
| Phase 3 | P0 / P1 endpoint、認証、session、error body、SSE、状態 read/write、systemd service の結果。 | API P0 / P1 を完了扱いにせず、Phase 4 着手禁止。 |
| Phase 4 | P2〜P5 endpoint、secret mask、rollback、maintenance、hook、token、P0 / P1 互換確認の結果。 | endpoint 契約を固定扱いにせず、Phase 5 着手禁止。 |
| Phase 5 | method 対応表、fake fetch、HTTP error、timeout、stream、`401` token 破棄、body 禁止の結果。 | SDK 契約を固定扱いにせず、Phase 6 着手禁止。 |
| Phase 6 | DOM id、panel、SDK 呼び出し、success/error/loading/disabled、stream、secret 消去、直接 API 呼び出し不存在の結果。 | 初期実装完了扱いにせず、未充足 UI 仕様を解消する。 |

Phase 別受け入れ条件のいずれかが未実行、失敗、または環境都合で省略された場合、その Phase を完了扱いにしてはならない。後続 Phase の実装 PR を開始する前に、先行 Phase の未充足条件を仕様または実装で解消する。

**Phase 完了判定 fixture 記録固定：**

| Phase | fixture 記録 | fake 記録 | 期待値記録 |
|-------|--------------|-----------|------------|
| Phase 1 | 使用した `testdata/builder/...` path、入力 Markdown、テーマ設定、asset 入力。 | fake filesystem を使った場合のみ、注入した失敗と対象 path。 | 生成物一覧、期待 HTML / CSS / JS / search index、stdout / stderr、終了コード、冪等性結果。 |
| Phase 2 | 使用した `testdata/runner/r*` path、初期 `state/`、GitHub 応答、pipeline / deploy / notify fixture。 | fake GitHub server、fake ssh executable、fake notifier、fake filesystem の呼び出し記録。 | `.last_sha`、queue、snapshot、通知結果、retry / circuit breaker、状態ファイル更新有無。 |
| Phase 3 | 使用した `testdata/api/p0-p1/...` path、request、初期 `state/`。 | API handler への request 記録、SSE frame 記録、状態 read/write 記録。 | HTTP status、response body、error body、SSE event、状態ファイル更新後 expected。 |
| Phase 4 | 使用した `testdata/api/p2-p5/...` path、request、secret mask fixture。 | hook / notify / token / maintenance の fake 呼び出し記録。 | HTTP status、response body、mask 後 payload、token 再取得不可結果、P0 / P1 回帰結果。 |
| Phase 5 | 使用した `testdata/sdk/...` path、fake fetch transcript、stream chunk。 | fake fetch の method、URL、headers、body 有無、abort / timeout 記録。 | SDK return、`AdlaireCIError`、stream event、token 破棄、body 禁止 endpoint の request 不成立。 |
| Phase 6 | 使用した `testdata/ui/...` path、初期 DOM、操作手順。 | fake SDK の method、引数、呼び出し順、stream unsubscribe 記録。 | DOM assertion、success / error / loading / disabled、secret 消去、直接 `fetch()` 不存在。 |

**未実行検証の代替条件：**

| 未実行対象 | 代替として認める確認 | 完了扱い |
|------------|----------------------|----------|
| 実 GitHub API、実 SSH、実通知先 | §0g.8 の fake GitHub server、fake ssh executable、fake notifier の fixture 検証。 | 認める。実外部接続は Phase 完了条件に含めない。 |
| systemd 実起動 | service file の内容確認、起動 command / user / working directory / environment の静的確認、API handler の Go test。 | 実 target host での起動確認が仕様化された Phase では `未実行`。仕様化されていない場合は静的確認で可。 |
| browser 実機操作 | Vanilla JS が動作する browser runtime または同等 DOM 環境での fake SDK / DOM assertion。 | 認める。ただし直接 `fetch()` 不存在、secret 消去、loading / disabled は必須。 |
| `go test` または `gofmt -l` が実行不能 | 実行不能理由、未実行 command、再実行条件の記録。 | 認めない。該当 Phase は `未実行` を含むため完了扱い不可。 |
| fixture expected の欠落 | 欠落 path、期待値を定義できない理由の記録。 | 認めない。fixture expected を追加するまで完了扱い不可。 |

**fixture 期待値更新固定：**

| 更新対象 | 更新条件 | 必須確認 |
|----------|----------|----------|
| `expected/` 内の生成物 | 仕様本文の出力契約、schema、error body、DOM id、SDK return のいずれかが変更された場合のみ更新する。 | 仕様本文の該当節と fixture expected が同じ値を示すこと。 |
| fake transcript | 外部依存の呼び出し method、path、payload、retry、timeout の仕様が変更された場合のみ更新する。 | secret / token / password 原文が transcript に存在しないこと。 |
| DOM assertion | UI の DOM id、panel、表示文言、disabled / loading / success / error 条件が変更された場合のみ更新する。 | SDK method 対応表と DOM assertion が一致すること。 |
| error expected | HTTP status、exit code、`AdlaireCIError.code`、stderr prefix が変更された場合のみ更新する。 | 正常系 fixture と異常系 fixture の両方で期待値が固定されていること。 |

Phase 完了判定テンプレートは以下とする。実装 PR 本文では、対象 Phase ごとに本テンプレートの項目を埋める。

```text
## Phase 完了判定

- 対象 Phase:
- owner component:
- collaborator component:
- 実装対象ファイル:
- 追加 fixture / testdata:
- 固定契約:
- 実装対象外:
- 後続 Phase への影響:

| 対象 | コマンド / 確認 | 期待結果 | 実結果 | 判定 |
|------|------------------|----------|--------|------|
| <対象> | <実行内容> | <仕様上の期待結果> | <実際の結果> | PASS / FAIL / 未実行 |
```

`判定` が `FAIL` または `未実行` の行を含む場合、その Phase は完了扱いにしてはならない。環境制約により確認できない項目がある場合も `未実行` とし、完了扱いにするには代替検証を仕様化してから再実行する。

受け入れ結果は、実装 PR 本文に `対象 / コマンド / 期待結果 / 実結果 / 判定` の形式で記録する。失敗、未実行、環境都合で省略した項目がある場合、そのコンポーネントを完了扱いにしてはならない。

---
