# Adlaire CI — Setup 詳細仕様

`setup` owner component の詳細本文責務は、[`docs/details/setup.md`](setup.md) を正本とする。

詳細本文境界管理条件は [`docs/DETAIL_INDEX.md`](../DETAIL_INDEX.md) 詳細仕様入口責務 §0b.1 に従う。`setup` owner component の主本文であり、collaborator component の仕様は配置対象、状態初期化、admin 配布、service health、fixture、検証観点として参照する。

[`docs/details/setup.md`](setup.md) 詳細本文責務は、バイナリ配布、配置、systemd、セットアップ、アップデート、リリース成果物検証を定義する。runner / api / sdk / ui / admin の個別機能本文は各 owner component 別の [`docs/details/*.md`](../details/) 詳細本文責務を参照する。状態判断は [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務を参照し、fixture、fake、expected / effects、実装検証証跡は [`docs/details/fixture.md`](fixture.md) fixture 証跡責務を参照する。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `setup` |
| collaborator component | `runner`、`api`、`statefile`、`admin` |
| 持つ内容 | `setup` owner が主本文として定義するバイナリ配布、配置、systemd、セットアップ、アップデート、リリース成果物検証。 |
| 持たない内容 | runner / api / sdk / ui / admin の個別機能本文、状態 schema、API endpoint、SDK method、UI DOM、fixture 証跡責務、外部依存追加。 |

---

## 26. セットアップ・アップデート手順

該当節は、Go 版 Adlaire CI のセットアップ手順を定義する。

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
| `adlaire-ci-build-$OS_ARCH` | 初回セットアップ、アップデート | `components/builder.go` から生成した Markdown → 静的 Web サイトビルドバイナリ。 |
| `adlaire-ci-runner-$OS_ARCH` | 初回セットアップ、アップデート | `components/runner.go` から生成した CI ランナーバイナリ。 |
| `adlaire-ci-api-$OS_ARCH` | 管理 API 導入手順、管理 API 導入後のアップデート | `components/api.go` から生成した管理 API サーバーバイナリ。 |
| `admin-ui.tar.gz` | 管理 API 導入手順、管理 API 導入後のアップデート | [`docs/details/admin.md`](admin.md) A1 の管理 UI 配布物。 |
| `SHA256SUMS` | Release 添付ファイル取得時 | Release 添付ファイルの SHA-256 checksum 一覧。 |

Release asset 名は上表の文字列と完全一致させる。`$OS_ARCH` は `linux-amd64` だけを初期標準とし、未知 OS/arch を指定した場合は取得前に `unsupported OS_ARCH: {OS_ARCH}` を stderr へ出力して終了コード `2` とする。`SHA256SUMS` は `"{sha256}  {filename}"` 形式の LF 区切り text とし、対象 filename が 1 回だけ出現することを必須とする。対象行が 0 件または 2 件以上の場合は checksum 検証失敗とする。

### §26.2b セットアップ・アップデート機能単位

セットアップ・アップデート実装は、以下の機能単位に分離する。各機能は前段の出力だけを入力として受け取り、失敗時は後続機能を実行しない。

| 機能 | 入力 | 出力 | 失敗条件 | 失敗時の終了状態 |
|------|------|------|----------|------------------|
| Release asset resolver | `VERSION`、`OS_ARCH`、取得対象成果物名、GitHub Release URL | `DOWNLOAD_DIR` 内の取得済みファイル | `VERSION` / `OS_ARCH` 空、HTTP status 非 2xx、取得ファイル 0 byte | 取得済みファイルを配置せず終了 |
| checksum verifier | `SHA256SUMS`、取得済み成果物 | 検証済み成果物一覧 | `SHA256SUMS` 不在、対象行不在、SHA-256 不一致 | バイナリ配置を実行せず終了 |
| binary installer | 検証済みバイナリ、`BIN_DIR` | `adlaire-ci-build`、`adlaire-ci-runner`、API 導入対象の実装では `adlaire-ci-api` | 入力バイナリ不在、実行権限付与失敗、`install` 失敗 | systemd 変更を実行せず終了 |
| secret initializer | PAT 入力、`INSTALL_DIR` | `.github_token` mode `0600` | PAT 空、書き込み失敗、mode 補正失敗 | systemd 変更を実行せず終了 |
| state initializer | `INSTALL_DIR` | `.last_sha`、build log 保存対象の実装では `.build_logs/`、snapshot 保存対象の実装では `.snapshots/` | 書き込み失敗、mode 補正失敗 | systemd 変更を実行せず終了 |
| systemd unit writer | unit 内容、`SERVICE_USER`、`INSTALL_DIR`、`BIN_DIR` | `/etc/systemd/system/adlaire-ci.service`、`adlaire-ci.timer`、API 導入対象の実装では `adlaire-ci-api.service` | unit 書き込み失敗、`systemctl daemon-reload` 失敗 | enable/start を実行せず終了 |
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

`admin-ui.tar.gz` は checksum 検証後に一時展開ディレクトリへ展開する。配布物の中身、必須 file、拒否する archive entry は [`docs/details/admin.md`](admin.md) A1〜A2 を参照する。検証に失敗した場合は、既存 `$INSTALL_DIR/admin` を変更しない。

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

| 失敗箇所 | 保持するもの | 変更対象 | 禁止事項 |
|----------|--------------|------------------|----------|
| download / checksum | 既存 binary、既存 systemd、既存 state、既存 admin UI | `DOWNLOAD_DIR` 内の取得済みファイル | 未検証 asset の配置、service restart。 |
| binary 配置前 | 既存 binary、既存 service 稼働状態 | `DOWNLOAD_DIR` | systemd unit 書換、state 書換。 |
| binary 配置後 / restart 前 | 配置済み新 binary または rollback 対象旧 binary | rollback executor が対象 binary だけを復元する。 | state、history、secret、admin UI の巻き戻し。 |
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
| 最終確認 | なし | [`docs/details/setup.md`](setup.md) §26.3、[`docs/details/setup.md`](setup.md) §26.3b、[`docs/details/setup.md`](setup.md) §26.5 の確認項目がすべて成功。 | 成功報告しない。確認失敗箇所と journal 確認対象を出力する。 |

setup / update 実装は、各段階の開始と成功を stderr または stdout に固定文言で 1 行ずつ出力する。PAT、password、session token、API token、Webhook secret、SMTP password、Release URL の credential 部分は出力してはならない。secret file が既に存在する場合は、個別手順で上書きを明記している場合を除き、既存値を保持する。特に `.github_token`、`.admin_credentials`、`.webhook_secret`、`.smtp_secret` は、アップデートで自動上書きしない。

`systemctl daemon-reload` 成功だけではセットアップ成功と扱わない。`enable --now`、`restart`、`is-active`、API 導入時の `/api/health` 確認まで完了して初めて成功とする。確認コマンドが利用環境に存在しない場合は、Go `net/http` client または systemd D-Bus / `systemctl show` で同じ確認項目を検証し、実装検証証跡に代替コマンド、期待値、実測値を記録する。未確認のまま成功扱いにしない。

### §26.3 Go 版初回セットアップ手順

対象は Go 版の `components/builder.go` と `components/runner.go` から生成した `adlaire-ci-build`、`adlaire-ci-runner`、`adlaire-ci.service`、`adlaire-ci.timer` とする。

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
| build binary | `$BIN_DIR/adlaire-ci-build --version` | exit `0`、stdout が `adlaire-ci-build v3` を含む。 |
| runner binary | `$BIN_DIR/adlaire-ci-runner --version` | exit `0`、stdout が `adlaire-ci-runner v3` を含む。 |
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
# runner systemd unit の内容を /etc/systemd/system/ に配置した上で:
systemctl daemon-reload

# ── 7. タイマー有効化・起動 ───────────────────────────
systemctl enable --now adlaire-ci.timer

# ── 8. 起動確認 ───────────────────────────────────────
systemctl status adlaire-ci.timer
```

Go 版初回セットアップでは以下を実行しない。

| 対象 | 理由 |
|------|------|
| `/usr/local/bin/adlaire-ci-api --init-credentials --state-dir "$INSTALL_DIR"` | 初回セットアップ対象は runner と build バイナリに限定し、API 認証情報生成は [`docs/details/setup.md`](setup.md) §26.3b で実行する。 |
| `systemctl enable --now adlaire-ci-api` | API service は [`docs/details/setup.md`](setup.md) §26.3b の API バイナリ配置、認証情報生成、unit 配置がすべて成功した後にのみ起動する。 |
| `.build_logs/` 作成 | runner 初期導入ではビルド実行時に必要な状態だけを初期化し、api が参照する履歴ディレクトリは [`docs/details/setup.md`](setup.md) §26.3b で作成する。 |
| `.snapshots/` 作成 | snapshot 参照・rollback API と組み合わせて使うため、[`docs/details/setup.md`](setup.md) §26.3b の管理 API 導入時に作成する。 |

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
| API binary | `$BIN_DIR/adlaire-ci-api --version` | exit `0`、stdout が `adlaire-ci-api v3` を含む。 |
| credentials | `stat -c '%a' "$INSTALL_DIR/.admin_credentials"` | `600`。 |
| admin UI | `test -f "$INSTALL_DIR/admin/index.html"` / `test -f "$INSTALL_DIR/admin/adlaire-ci-sdk.js"` | 両方成功。 |
| API service | `systemctl is-active adlaire-ci-api` | `active`。 |
| local health | `curl -fsS http://127.0.0.1:8765/api/health` | local health endpoint の response 契約は [`docs/details/api.md`](api.md) §22.0e を参照し、setup 側は local API 到達と応答取得を確認する。 |

`curl` が利用できない環境では、Go 実装変更の検証で `net/http` client または同等のローカル HTTP 確認を行う。HTTP status、response body、必須 key の具体契約は [`docs/details/api.md`](api.md) §22.0e を参照し、未確認のまま API 導入確認を満たした扱いにしてはならない。

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
# 管理 API systemd unit の内容を /etc/systemd/system/adlaire-ci-api.service に配置した上で:
systemctl daemon-reload

# ── 7. サービス有効化・起動 ───────────────────────────
systemctl enable --now adlaire-ci-api

# ── 8. 起動確認 ───────────────────────────────────────
systemctl status adlaire-ci-api
```

### §26.4 systemd サービスファイル

#### §26.4.1 Go 版 runner の systemd ファイル

**`/etc/systemd/system/adlaire-ci.service`**（`components/runner.go`）：

```ini
[Unit]
Description=Adlaire CI Runner

[Service]
Type=oneshot
User=root
WorkingDirectory=/opt/adlaire-builder
ExecStart=/usr/local/bin/adlaire-ci-runner --state-dir /opt/adlaire-builder
```

**`/etc/systemd/system/adlaire-ci.timer`**（`components/runner.go` 定期起動タイマー）：

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

`User` / `WorkingDirectory` / `ExecStart` のパスは [`docs/details/setup.md`](setup.md) §26.2 の設定変数に合わせて変更する。

systemd unit は上記キー以外を初期標準で追加しない。`Environment=`、`EnvironmentFile=`、`ExecStartPre=`、`ExecStartPost=` を追加する場合は、先に該当節へ対象変数、secret 扱い、失敗時挙動を定義する。API service は `127.0.0.1:8765` bind を標準とし、外部公開 bind は [`docs/details/setup.md`](setup.md) 詳細本文責務で未定義のため設定しない。

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
| local API | 確認しない。 | `GET /api/health` の endpoint 契約は [`docs/details/api.md`](api.md) §22.0e を参照し、API service が local health check に応答する。 |
| state preservation | `.github_token`、`.last_sha`、`.build_state`、`.build_history` の mtime と内容が更新対象操作と無関係に変わっていない。 | 左記に加え `.admin_credentials` が存在する場合は mode `600` と内容が保持される。 |

確認失敗時はアップデート失敗として扱う。binary 配置や restart が成功していても、確認失敗を成功報告してはならない。local API 確認で `curl` がない場合は Go 実装変更の検証で `net/http` client による `GET /api/health` 相当の確認を行う。HTTP status、JSON object、`status` key の具体契約は [`docs/details/api.md`](api.md) §22.0e を参照し、実装検証証跡形式は [`docs/details/fixture.md`](fixture.md) §0g.8-F を参照する。未確認のまま合格扱いにしない。

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

setup / release / update の詳細実装確認では、下表の受け入れ条件をすべて満たす。実装対象外のコンポーネントは「未実装」として明記し、確認済み扱いにしない。

| 対象 | 必須コマンド / 確認 | 合格条件 |
|------|---------------------|----------|
| Go 共通 | `gofmt -l <実装対象Goファイル>` | 実装対象 Go ファイルが存在する場合、出力が空。未作成ファイルはコマンド対象に含めない。 |
| Go test | `go test ./...` | Go module が存在する場合に成功する。Go module が存在しない場合は、その理由を実装確認結果に明記する。 |
| build script | `adlaire-ci-build --src <sample.md> --out <tmp.html>` | exit code `0`、HTML 出力あり、`[REPORT]` の `status` が `success`。 |
| runner | `adlaire-ci-runner --state-dir <tmp-state>` | 必須 secret 未設定時の exit code / ERROR log が [`docs/details/runner.md`](runner.md) §12 と一致し、`.build_lock` が残らない。 |
| API | API service の起動、local health check、admin UI から到達可能な endpoint 境界を確認する。 | API endpoint、request / response、HTTP status、body は [`docs/details/api.md`](api.md) §22.0 / [`docs/details/api.md`](api.md) §22.0e を参照する。 |
| SDK | admin UI 配布物に SDK 静的ファイルが含まれ、browser runtime から読み込めることを確認する。 | SDK method、transport、error、stream、timeout は [`docs/details/sdk.md`](sdk.md) §23 を参照する。 |
| UI | admin UI 配布物が静的配信され、ログイン画面と主要 panel へ到達できることを確認する。 | UI DOM、disabled、成功表示、失敗表示、再取得、秘密情報消去は [`docs/details/ui.md`](ui.md) §24 を参照する。 |
| setup | [`docs/details/setup.md`](setup.md) §26.3 または [`docs/details/setup.md`](setup.md) §26.3b の手順を fresh 環境で実行する。 | unit 配置、権限、`systemctl is-active`、secret mode が仕様どおり。 |
| update | [`docs/details/setup.md`](setup.md) §26.5 の手順を前版バイナリから新 tag のリリースバイナリへ実行する。 | 旧バイナリ退避、新バイナリ配置、restart、失敗時 rollback 条件が仕様どおり。 |
| security | secret 値を含む入力後、stdout、stderr、journal、API response、UI 表示の漏えい有無を確認する。 | 漏えい禁止対象と security 処理本文は [`docs/details/security.md`](security.md) §25、[`docs/details/security.md`](security.md) §27.42〜§27.47 を参照する。setup 側は配置・保持・権限・log 出力を確認する。 |

**認証 / セットアップ fixture 固定：**

| fixture | 入力 | 合格条件 |
|---------|------|----------|
| credentials init success | 空 state dir で `adlaire-ci-api --init-credentials --state-dir <abs>` | `.admin_credentials` mode `600`、schema 全 key、`must_change=true`、stdout 固定文言。 |
| credentials init existing | `.admin_credentials` 既存 | exit `2`、stderr `credentials already exist`、既存ファイル差分なし。 |
| login success | 初期 password `admin` | `must_change:"prompt"`、TOTP 無効時 token 発行、hash/salt 非表示、`.access_log` 成功行。 |
| login failure lock | password 連続 10 回失敗 | lockout 期間、HTTP status、error body は [`docs/details/security.md`](security.md) §25 と [`docs/details/api.md`](api.md) §22.0e を参照する。setup 側は password 詳細非表示と secret 非保存を確認する。 |
| change password | current 正、新 password 有効 | `.admin_credentials` の salt/hash 更新、`must_change=false`、現 session 以外破棄。 |
| session restart | token 発行後に API process restart | restart 後の認証失敗表現は [`docs/details/security.md`](security.md) §25 と [`docs/details/api.md`](api.md) §22.0e を参照する。setup 側は session file が存在しないことを確認する。 |
| setup checksum mismatch | Release asset と `SHA256SUMS` 不一致 | バイナリ配置なし、systemd 変更なし、終了コード `1`。 |
| update restart failure | 新バイナリ配置後に service restart 失敗 | 旧バイナリ復元を 1 回だけ行い、state/history/secret は巻き戻さない。 |

**セットアップ / アップデート fixture 固定：**

| fixture | 入力 | 合格条件 |
|---------|------|----------|
| setup variable invalid | `INSTALL_DIR=/`、`BIN_DIR=/`、または `OS_ARCH=darwin-arm64` | 終了コード `2`、download / directory 作成 / 配置なし。 |
| setup download failure | build binary 取得が失敗 | 終了コード `1`、runner binary 取得済みでも配置なし、systemd 変更なし。HTTP status の具体値は fake / fixture 側で固定する。 |
| setup checksum duplicate | `SHA256SUMS` に同一 asset 行が 2 件 | 終了コード `1`、checksum failure、配置なし。 |
| setup symlink target | `$BIN_DIR/adlaire-ci-build` が symlink | 終了コード `1`、symlink 参照先を上書きしない。 |
| setup pat empty | PAT 入力が空 | 終了コード `2`、`.github_token` 作成なし、systemd 変更なし。 |
| setup success | 正常 asset、正常 PAT、fresh 環境 | binary mode `755`、`.github_token`/`.last_sha` mode `600`、timer active、固定確認すべて成功。 |
| setup admin release layout | `admin-ui.tar.gz` と API binary を含む正常 Release | `admin-ui.tar.gz` は [`docs/details/admin.md`](admin.md) A1 の root layout と一致し、API binary と admin UI の version が同一 `VERSION`。 |
| api setup credentials existing | `.admin_credentials` 既存で API 導入 | `--init-credentials` を再実行せず既存 credentials を保持し、API service 起動確認まで進む。 |
| api setup admin archive unsafe | `admin-ui.tar.gz` に `../x` または symlink entry | 終了コード `1`、既存 admin UI 維持、API service start なし。 |
| api setup admin archive extra file | `admin-ui.tar.gz` に A1 未定義 file、重複必須 file、`admin/` wrapper directory | 終了コード `1`、archive 展開なし、既存 admin UI 維持、API service start なし。 |
| api setup health failure | API service active だが `/api/health` が非 200 | セットアップ失敗扱い、runner timer は停止しない、journal 確認対象を出力。 |
| setup secret preservation | `.github_token`、`.admin_credentials`、`.webhook_secret`、`.smtp_secret` が既存の状態で update | 明示対象外の secret は content / mode / mtime を保持し、stdout/stderr/journal に secret 原文を出さない。 |
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
| snapshot rollback running | `.build_state.running=true` で `POST /api/history/{id}/rollback` | HTTP status / error body は [`docs/details/api.md`](api.md) §22.0e を参照する。queue 追加なし、history 追記なし。 |
| maintenance enable no-op | 同一 reason で enable を 2 回実行 | API response は [`docs/details/api.md`](api.md) §22.0e を参照する。`.config_log` 追記なし。 |
| access-control deny | allow に接続元以外を設定して API 呼び出し | HTTP status / error body は [`docs/details/api.md`](api.md) §22.0e と [`docs/details/security.md`](security.md) §27.42〜§27.47 を参照する。password/token 検証なし。 |
| hook pre abort | `pre` hook が exit `1`、`abort_on_failure=true` | pipeline 未実行、build status `hook_error`、hook log 保存。 |
| alert duplicate | 同一 alert rule を 2 回作成 | HTTP status / error body は [`docs/details/api.md`](api.md) §22.0e を参照する。`.alert_rules` 差分なし。 |
| tag rule invalid | 破損 condition を含む `.tag_rules` で build | build failure、ERROR log `TAG_RULE_INVALID`、SHA cache 更新なし。 |
| verify-output no history | 成功履歴なしで `POST /api/verify-output` | HTTP status / error body は [`docs/details/api.md`](api.md) §22.0e を参照する。状態ファイル変更なし。 |
| pipeline config reserved arg | `extra_args:["--src","x"]` | HTTP status は [`docs/details/api.md`](api.md) §22.0e を参照する。`.pipeline_config` 差分なし。 |
| notes no-op | 同一 content を 2 回保存 | API response は [`docs/details/api.md`](api.md) §22.0e を参照する。`.config_log` 追記なし。 |
| smtp secret mask | password 付き `POST /api/smtp-config` 後に GET / backup / log 確認 | password 本体は返らず、mask / secret 表示契約は [`docs/details/api.md`](api.md) §22.0e と [`docs/details/security.md`](security.md) §25 を参照する。 |
| queue disabled | `queue_max_size=0`、build running 中に `POST /api/build` | HTTP status / error body は [`docs/details/api.md`](api.md) §22.0e を参照する。`.build_state.queued` は空。 |
| dashboard duplicate widget | widgets に重複 id を指定 | HTTP status は [`docs/details/api.md`](api.md) §22.0e を参照する。`.dashboard_layout` 差分なし。 |

**関連責務参照：**

[`docs/details/api.md`](api.md) §22、[`docs/details/statefile.md`](statefile.md) §22.0a / [`docs/details/statefile.md`](statefile.md) §22.0c、[`docs/details/sdk.md`](sdk.md) §23、[`docs/details/ui.md`](ui.md) §24、[`docs/details/security.md`](security.md) §25 / [`docs/details/security.md`](security.md) §27.42〜§27.47、[`docs/details/setup.md`](setup.md) §26 にまたがる API、状態ファイル、SDK、UI、認証、setup / update の整合は、各 owner component 別の詳細本文責務と fixture 証跡責務を同時に参照する。該当節は setup / release / update の実行条件だけを扱い、API endpoint、SDK method、UI 操作、認証方式、状態 schema、fixture 名、実装検証証跡項目を重複定義しない。

| 対象 | 主本文 | setup 側の確認範囲 |
|------|--------|--------------------|
| API endpoint / response / read-write 境界 | [`docs/details/api.md`](api.md) §22 | API service の配置、起動、health check、systemd 連携だけを確認する。 |
| 状態ファイル schema / atomic write / 破損時処理 | [`docs/details/statefile.md`](statefile.md) §22.0a、[`docs/details/statefile.md`](statefile.md) §22.0c | setup が初期作成または保持する path、権限、既存 state / secret 保護だけを確認する。 |
| SDK method / transport / error | [`docs/details/sdk.md`](sdk.md) §23 | admin UI 配布時に SDK 静的ファイルを配置することだけを確認する。 |
| UI DOM / 操作 / 表示状態 | [`docs/details/ui.md`](ui.md) §24 | admin UI 配布物と静的配信境界だけを確認する。 |
| 認証 / session / token / TOTP / audit | [`docs/details/security.md`](security.md) §25、[`docs/details/security.md`](security.md) §27.42〜§27.47 | secret / credential file の配置、保持、権限、漏えい防止だけを確認する。 |
| fixture / fake / expected / 実装検証証跡 | [`docs/details/fixture.md`](fixture.md) §0g.8-F、[`docs/details/fixture.md`](fixture.md) §22-F、[`docs/details/fixture.md`](fixture.md) §27-F | setup / release / update に関わる証跡の記録先だけを確認する。 |

**setup / update 失敗時副作用固定契約：**

| ケース | 固定結果 |
|--------|----------|
| setup partial failure | 既存 binary、state、secret、admin UI、systemd を、[`docs/details/setup.md`](setup.md) §26 で許可した対象以外は変更しない。 |
| update rollback failure | 追加推測復旧を行わず、失敗箇所、退避先、現在配置済みファイル、journal 確認対象を報告する。 |
| checksum / download / unsafe archive failure | binary、admin UI、systemd、state、secret を変更せず、失敗箇所と再実行条件を記録する。 |

**setup / update 実装前・実装後確認固定契約：**

| 段階 | 確認 | 合格条件 |
|------|------|----------|
| 実装前 | setup / update 対象 | 配置対象 binary、admin UI asset、systemd unit、state / secret 保持対象、rollback 対象が [`docs/details/setup.md`](setup.md) §26 に定義済み。 |
| 実装前 | secret handling | setup / update が触る secret file の保存先、権限、保持条件、log 禁止が定義済み。 |
| 実装後 | setup/update | checksum、unsafe archive、restart failure、rollback failure、health failure が [`docs/details/setup.md`](setup.md) §26 の fixture と一致する。 |

Phase 順序、実装変更単位、実装着手条件、判定責務は [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務 §4.1 と [`docs/SPEC.md`](../SPEC.md) ポリシー責務 §0f を参照する。

Phase 別の fixture、fake、expected / effects、実装検証証跡、不足時の扱いは [`docs/details/fixture.md`](fixture.md) §0g.8-F を参照する。`setup` 詳細では Phase 別の受け入れ条件、検証記録、fixture 証跡項目を重複定義しない。

**setup / release 未実行検証の代替条件：**

| 未実行対象 | 代替として認める確認 | 確認済み扱い |
|------------|----------------------|----------|
| systemd 実起動 | service file の内容確認、起動 command / user / working directory / environment の静的確認、API handler の Go test。 | 実 target host での起動確認が仕様化された Phase では `未実行`。仕様化されていない場合は静的確認で可。 |
| `go test` または `gofmt -l` が実行不能 | 実行不能理由、未実行 command、再実行条件の記録。 | 認めない。setup / update に必要な検証が未実行の場合は確認済み扱い不可。 |
| release asset checksum 検証不能 | 対象 asset、取得元、期待 checksum、検証不能理由、再実行条件を記録する。 | 認めない。checksum 検証確認まで setup / update 確認済み扱い不可。 |

**fixture 期待値更新固定：**

| 更新対象 | 更新条件 | 必須確認 |
|----------|----------|----------|
| `expected/` 内の生成物 | setup / release / update の配置、起動、保持、rollback、終了コード、標準出力、標準エラー、systemd、ファイル権限に関わる期待値が変更された場合のみ更新する。API response、UI DOM、SDK return、状態 schema の期待値更新条件は、それぞれの owner component 別詳細本文責務と [`docs/details/fixture.md`](fixture.md) §0g.8-F を参照する。 | setup / release / update の期待値と該当節の setup 責務が一致し、他 owner component の expected は該当する詳細本文責務と一致すること。 |
| fake transcript | setup / release / update が直接実行する Release asset 取得、checksum 検証、systemd、権限確認、file 配置、rollback の呼び出し仕様が変更された場合のみ更新する。API / SDK / UI / security の fake 更新条件は該当する owner component 別詳細本文責務と [`docs/details/fixture.md`](fixture.md) §0g.8-F を参照する。 | secret / token / password 原文が transcript に存在しないこと。 |
| DOM assertion | setup / release / update が admin UI 静的ファイルの配置と到達確認を変更した場合のみ、setup 側の検証入口として更新する。UI DOM、panel、表示文言、disabled / loading / success / error 条件の具体契約は [`docs/details/ui.md`](ui.md) §24 と [`docs/details/fixture.md`](fixture.md) §0g.8-F を参照する。 | setup 側は admin UI 配布物の存在と到達確認だけを固定し、DOM assertion の具体値は UI 詳細本文責務と一致すること。 |
| error expected | setup / release / update が直接返す exit code、stderr prefix、rollback 結果、配置失敗結果が変更された場合のみ更新する。HTTP status、API error body、`AdlaireCIError.code` の具体契約は [`docs/details/api.md`](api.md) §22.0e、[`docs/details/sdk.md`](sdk.md) §23、[`docs/details/fixture.md`](fixture.md) §0g.8-F を参照する。 | setup / release / update の正常系 fixture と異常系 fixture の両方で setup 責務の期待値が固定され、API / SDK / UI の期待値は各 owner component の詳細本文責務と一致すること。 |

Phase 判定の実装検証証跡テンプレート、必須記載項目、不足時の扱いは [`docs/details/fixture.md`](fixture.md) §0g.8-F を参照する。該当節は setup / release / update の実行条件、setup / release 未実行検証の代替条件、fixture 期待値更新条件だけを定義し、実装検証証跡項目を重複定義しない。

setup / release / update に関わる受け入れ結果は、[`docs/details/fixture.md`](fixture.md) §0g.8-F の形式で実装検証証跡に記録する。失敗、未実行、環境都合で省略した項目がある場合、setup / release / update を確認済み扱いにしてはならない。

---

### §26.8 Setup / Admin 配布実装確認ゲート

セットアップ、アップデート、管理 API 導入、admin UI 配布の詳細実装確認では、[`docs/details/setup.md`](setup.md) §26.1〜§26.7 の本文に加えて下表を満たす。該当節は実装時の確認粒度を固定するための詳細であり、未定義の成果物、未定義の service、未定義の rollback 対象を追加する根拠にしてはならない。

| 段階 | 必須入力 | 成功確定条件 | 失敗時固定結果 | fixture 必須 |
|------|----------|--------------|----------------|--------------|
| 変数検証 | `VERSION`、`OS_ARCH`、`INSTALL_DIR`、`BIN_DIR`、`DOWNLOAD_DIR` | 空値なし、危険 path なし、`OS_ARCH=linux-amd64`。 | 終了コード `2`。directory、download、配置なし。 | `setup variable invalid` |
| Release 取得 | asset URL、`SHA256SUMS` | 対象 asset と `SHA256SUMS` が HTTP 2xx、size > 0。 | 終了コード `1`。未検証 asset を配置しない。 | `setup download failure` |
| checksum | asset、`SHA256SUMS` | 対象 filename が 1 行だけ存在し、SHA-256 が一致する。 | 終了コード `1`。binary、admin、systemd、state 差分なし。 | `setup checksum duplicate`、`update checksum before change` |
| binary 配置 | 検証済み binary | symlink でない通常 file へ `0755` で配置し、`--version` が期待値を返す。 | systemd を変更しない。restart 前失敗なら旧 binary を保持する。 | `setup symlink target`、`setup success` |
| secret / state 初期化 | PAT、初期 state | secret `0600`、`.last_sha` `0600`、LF 付き JSON、fsync 完了。 | systemd を変更しない。secret 値を出力しない。 | `setup pat empty`、`setup success` |
| admin archive 展開 | `admin-ui.tar.gz` | [`docs/details/admin.md`](admin.md) A1〜A2 を満たし、一時 directory 検証後に差し替える。 | 既存 `$INSTALL_DIR/admin` を変更しない。API service を起動 / restart しない。 | `api setup admin archive unsafe` |
| systemd 配置 | unit file 内容 | unit 書込、mode、`daemon-reload`、enable/start/restart、`is-active` が成功する。 | enable/start/restart を成功扱いしない。journal 確認対象を出力する。 | `api setup health failure`、`update runner restart failure` |
| rollback | 旧 binary / 旧 admin backup | 定義済み対象だけ 1 回復元し、対象 service を 1 回 restart する。 | 追加推測復旧を行わず、現在配置済み path と journal 確認対象を出力する。 | `update rollback failure` |
| 最終確認 | 配置済み binary、state、service、admin UI | [`docs/details/setup.md`](setup.md) §26.3 / [`docs/details/setup.md`](setup.md) §26.3b / [`docs/details/setup.md`](setup.md) §26.5 の固定確認がすべて成功する。 | 成功報告しない。未確認項目を `未実行` として記録する。 | `setup success`、`update api absent` |

**setup / admin / release 連動 fixture 固定：**

| fixture 群 | 対象 component | 必須 input | 必須 expected | 合格条件 |
|------------|----------------|------------|---------------|----------|
| `setup-admin-release-layout` | `setup`、`admin` | Release asset 一式、`SHA256SUMS`、`admin-ui.tar.gz`、fake download response。 | `expected/effects.json`、admin archive file list、`expected/security.json`。 | 対象 asset 名、checksum、admin archive root layout、必須 file、任意 file、file mode、directory mode が [`docs/details/setup.md`](setup.md) §26.2a と [`docs/details/admin.md`](admin.md) A1〜A2 に一致する。 |
| `setup-admin-archive-boundary` | `setup`、`admin` | unsafe archive、既存 `$INSTALL_DIR/admin`、既存 API binary、API service fake。 | `expected/effects.json.unchanged_paths`、`forbidden_writes`、`forbidden_calls`、`expected/stderr.txt`。 | archive 検証失敗時に既存 admin UI、API binary、credentials、runner state を変更せず、API service start / restart を呼ばない。 |
| `setup-systemd-rollback-boundary` | `setup`、`runner`、`api` | systemd fake、旧 binary backup、旧 admin backup、restart failure。 | `expected/effects.json.write_order`、`updated_paths`、`unchanged_paths`、`forbidden_writes`、`commands`。 | rollback 対象は失敗段階で定義済みの binary / admin UI だけで、state、history、secret、runner timer を未定義に巻き戻さない。 |
| `admin-static-serving-security` | `admin`、`api` | static request、secret/state/log/snapshot path、method variation。 | `expected/response.json`、`expected/security.json`、`expected/effects.json`。 | A3 の status、header、body 有無に一致し、secret / state / log / snapshot / directory listing を返さず、request body を読まない。 |
| `setup-secret-preservation` | `setup`、`security` | 既存 secret files、update input、failure fake。 | `expected/security.json`、`expected/effects.json.unchanged_paths`、`forbidden_writes`。 | 明示対象外の secret content / mode / mtime を保持し、stdout、stderr、journal、fixture expected に secret 原文を残さない。 |

**setup / update 実装者向け出力固定：**

| 出力先 | 必須内容 | 禁止内容 |
|--------|----------|----------|
| stdout | 段階開始、段階成功、最終成功、配置 binary version。 | PAT、password、token、Webhook secret、SMTP password、Release URL credential。 |
| stderr | 固定 error prefix、失敗段階、終了コード、確認すべき journal / path。 | secret 原文、checksum 対象 file の内容、環境変数全量 dump。 |
| 実装検証証跡 | 実行 command、exit code、重要 stdout/stderr、差分あり / なし、未実行理由。 | secret 原文、credential 付き URL、ローカル固有 token。 |

**setup / update 差分確認固定：**

| ケース | 必須差分確認 | 合格条件 |
|--------|--------------|----------|
| fresh setup success | `$BIN_DIR`、`$INSTALL_DIR`、systemd unit | 仕様で許可された binary、secret、state、unit だけが作成される。 |
| fresh setup failure | `$BIN_DIR`、`$INSTALL_DIR`、systemd unit | 失敗段階より後の対象に差分がない。 |
| API 導入 success | API binary、admin directory、API unit、credentials | runner timer と runner state は不要に変更されない。 |
| API 導入 failure | admin backup、API binary、API unit | admin 展開失敗では既存 admin directory に差分がない。 |
| update success | 対象 binary、admin UI、systemd restart 記録 | 既存 state、history、secret は保持される。 |
| update failure | rollback 対象、journal 確認対象 | rollback 表で許可した対象以外に差分がない。 |

`setup` 実装変更は、上表の fixture、差分確認、secret 非表示確認、終了コード確認を記録する。いずれかが未実行の場合、対象段階を確認済み扱いにせず、未実行理由と再実行条件を記録する。
