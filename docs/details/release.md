# Adlaire CI — Release 詳細仕様

owner / collaborator 境界管理は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0b.1](../DETAIL_INDEX.md#0b1-owner-component-別-owner-collaborator-境界管理) に従う。Release の作成条件、バージョン体系、標準 OS/arch は [`docs/SPEC.md` ポリシー責務 §1 GitHub リリースポリシー](../SPEC.md#github-リリースポリシー)、asset の受け入れと配置は [`docs/details/setup.md`](setup.md)、admin archive の内容は [`docs/details/admin.md`](admin.md)、fixture と実装検証証跡は [`docs/details/fixture.md`](fixture.md) fixture 証跡責務を参照する。

---

<a id="release-responsibility-boundary"></a>
**責務境界：**

| 項目 | 内容 |
|------|------|
| owner component | `release` |
| 実装主体 | `components/release.go`。起動入口は `main.go`、実行ファイル名は `adlaire-ci-release`。 |
| 持つ内容 | Release 用バイナリ生成、admin archive 生成、checksum manifest 生成、再現性確認、GitHub draft Release 作成、asset upload・再取得検証、正式公開、失敗時 draft 削除。 |

`adlaire-ci-release` は保守者がリポジトリ checkout で実行する Release 作成用バイナリであり、利用者向け Release asset に含めない。利用者向け setup 実行バイナリは `adlaire-ci-setup-linux-amd64` として公開する。

---

<a id="release-cli-contract"></a>
**R1. CLI 固定契約：**

実行形式は次の exact 形式とする。

```text
adlaire-ci-release --repository owner/repository --tag V.X.N --commit 40-hex-sha --notes-file path --token-file path --out path
```

位置引数、短縮 option、`--name=value`、未知 option、同一 option の重複を禁止する。全値 option は `--name value` の 2 token 形式だけを許可する。`--help` と `--version` は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 CLI 共通固定契約](../DETAIL_INDEX.md#common-cli-contract) に従う。

| option | 必須 | 固定条件 |
|--------|------|----------|
| `--repository` | 必須 | `^[A-Za-z0-9_.-]{1,100}/[A-Za-z0-9_.-]{1,100}$`。`.`、`..` を owner または repository として許可しない。 |
| `--tag` | 必須 | `^V\.[1-9][0-9]*\.[0-9]+$`。前後空白、改行を禁止する。 |
| `--commit` | 必須 | lowercase hexadecimal 40 文字。短縮 SHA、大文字、branch 名を禁止する。 |
| `--notes-file` | 必須 | checkout root 内の絶対 path。path全segmentと対象fileがsymlinkでない通常file。UTF-8、1〜262144 bytes、NUL / CRなし、改行LF、末尾LF exact 1個。 |
| `--token-file` | 必須 | checkout root の内外を問わない絶対 path。全 path segment と対象 file が symlink でない通常 file、mode exact `0600`、実行 user 所有。 |
| `--out` | 必須 | checkout root 外の絶対 path。対象pathと`<out>.tmp`は存在してはならない。親 directory は既存の symlink でない書込可能 directory とする。`/`、checkout root、checkout root の親、symlinkを含む親segmentを禁止する。 |

固定 help は `Usage: adlaire-ci-release --repository owner/repository --tag V.X.N --commit 40-hex-sha --notes-file path --token-file path --out path [--version] [--help]` + LF とする。version 出力の binary name は `adlaire-ci-release` とする。

入力検証順序は、argv parse、必須 option、repository、tag、commit、checkout、notes file、token file、out path、local Git 状態、tag / commit 対応、GitHub 事前状態の順とする。最初の不合格だけを stderr に出し、成果物生成、GitHub write、directory 作成を開始しない。

| 終了コード | 条件 |
|------------|------|
| `0` | 成果物生成、再現性確認、GitHub asset 再取得検証、正式公開がすべて成功。 |
| `1` | file I/O、Go test/build、archive、checksum、再現性確認、draft cleanup の失敗。 |
| `2` | CLI、path、Git、tag、commit、既存 Release、token file の事前検証不合格。 |
| `3` | GitHub API、asset upload/download、外部 HTTPS 通信の失敗。 |

CLI parse error は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 CLI 共通固定契約](../DETAIL_INDEX.md#common-cli-contract) の固定文言を使用する。parse 完了後の失敗は、次の表で error code と終了コードを一意に決定する。複数条件が同時に成立する場合は [R1 入力検証順序](#release-cli-contract) または [R5 公開順序](#release-publish-contract) で最初に検出した条件だけを返す。

| error code | 条件 | 終了コード |
|------------|------|------------|
| `INVALID_INPUT` | repository、tag、commit、path、file type、mode、owner、UTF-8、size の不合格。 | `2` |
| `DIRTY_WORKTREE` | tracked または untracked の checkout 差分が 1 件以上。 | `2` |
| `TAG_MISMATCH` | HEAD、local tag、remote tag、`--commit`、remote default branch ancestor の不一致または対象 ref 不在。 | `2` |
| `RELEASE_EXISTS` | 同一 tag の draft または published Release が存在する。 | `2` |
| `FORMAT_FAILED` | `gofmt -d` が差分を返す、timeout、起動失敗。 | `1` |
| `TEST_FAILED` | `go test ./...` が非 0、timeout、起動失敗。 | `1` |
| `BUILD_FAILED` | `go build`、成果物 type / mode / size、binary 起動が不合格。 | `1` |
| `VERSION_MISMATCH` | 生成 binary の `--version` が [R4 ビルド・再現性固定契約](#release-build-contract) の固定値と不一致。 | `1` |
| `NON_REPRODUCIBLE` | 2 回生成した同名 asset の byte が不一致。 | `1` |
| `ARCHIVE_FAILED` | admin archive の生成、metadata、entry、再読込検証が不合格。 | `1` |
| `CHECKSUM_FAILED` | checksum manifest の生成または再検証が不合格。 | `1` |
| `OUTPUT_FAILED` | sibling staging作成、copy、mode、file sync、directory sync、atomic rename、parent syncの失敗。 | `1` |
| `GITHUB_READ_FAILED` | GitHub read-only API が HTTP、protocol、timeout、size 上限で失敗。 | `3` |
| `DRAFT_CREATE_FAILED` | draft Release create が失敗、または response が固定契約に不合格。 | `3` |
| `ASSET_UPLOAD_FAILED` | asset upload が失敗、または upload response が固定契約に不合格。 | `3` |
| `ASSET_VERIFY_FAILED` | asset list、再取得 size、digest のいずれかが不一致。 | `3` |
| `PUBLISH_FAILED` | metadata 再確認、publish、公開後再確認のいずれかが失敗。 | `3` |
| `DRAFT_CLEANUP_FAILED` | 正式公開前失敗後の draft DELETE が失敗。 | `1` |

---

<a id="release-precondition-contract"></a>
**R2. 実行前固定契約：**

release 実装は checkout root を `git rev-parse --show-toplevel` で確定し、以下を上から順に検証する。

1. current working directory が checkout root 内にある。
2. `git status --porcelain=v1 --untracked-files=all` が空である。
3. `git rev-parse HEAD` が `--commit` と exact 一致する。
4. `git rev-parse refs/tags/<tag>^{commit}` が `--commit` と exact 一致する。
5. GitHub repository の `full_name` と `default_branch`、remote default branch head ref を取得し、GitHub compare API で `--commit` がその head commit の ancestor または同一 commit であることを確認する。local `origin/main` のみで remote default branch 包含を判定しない。
6. GitHub の `refs/tags/<tag>` が存在し、peeled commit が `--commit` と一致する。
7. 同一 tag の GitHub Release が存在しない。draft を含む既存 Release がある場合は上書き、再利用、削除せず終了コード `2` とする。

release 実装は tag を作成、移動、削除、push してはならない。`git fetch`、branch 切替、merge、commit、push も行わない。

token file は root directory から各 path segment を no-follow の directory file descriptor として順に開き、確認済み parent descriptor に対して対象 file を no-follow で 1 回だけ開く。path 文字列への `Lstat` 後に通常の `Open` を行う実装は禁止する。open 後の file descriptor から type、device、inode、mode、owner、size、mtime を取得する。4097 bytes 目まで読み、UTF-8、末尾 LF は 0 または 1、trim 後 1〜4096 bytesとし、trim 後の token byte に空白・NUL・CR・LF・制御文字がないことと、読取前後の identity / mode / owner / size / mtime 不変を確認する。合格した token byte は memory だけに保持し、path を再 open しない。token を argv、environment、stdout、stderr、生成物、Release notes、fixture expected に保存しない。

notes file は checkout root descriptor から各相対 path segment を no-follow で順に開き、確認済み parent descriptor に対して対象 file を no-follow で 1 回だけ開く。open 後の file descriptor から type、device、inode、size、mtime を取得する。262145 bytes 目まで読み、[R1 CLI 固定契約](#release-cli-contract) の encoding、size、改行条件と読取前後の identity / size / mtime 不変を確認する。合格した byte を memory に 1 回だけ保持し、draft body、公開前照合、公開後照合のすべてで同じ byte を使用する。path の再 open と再読込みを禁止する。

`--out` は root directory から parent の各 path segment を no-follow の directory file descriptor として順に開き、最終 parent descriptor の device / inode を記録して処理終了まで保持する。staging directory 作成直前とatomic rename直前に、root directoryから同じno-follow手順でparent pathを再解決し、保持済みdescriptorとdevice / inodeが一致することを確認する。不一致は`OUTPUT_FAILED`とし、staging作成前はwrite 0件、作成後は保持済みdescriptor相対でstagingだけをcleanupする。`--out` と `<out>.tmp` の存在確認、staging directory 作成、rename、parent sync、cleanup はすべて保持済みparent descriptor に対する basename 相対操作で行う。path 文字列への検証後に通常のpath操作へ戻る実装、symlink を追跡する存在確認、異なる parent descriptor 間の rename を禁止する。

release 実装は永続設定ファイルと業務状態ファイルを持たず、`--commit` から生成した immutable source snapshot、保持済み notes byte、検証済み token、remote repository / default branch / tag / Release、`--out` を実行入力として扱う。live checkout file を build、test、archive 入力に使用しない。lock file は作成しない。同一 checkout または同一 tag に対する並行実行を成功扱いせず、GitHub write 直前の同一 tag Release 再確認と draft create の競合応答により片方だけを継続する。既存 Release、既存 draft、既存 asset を再利用、更新、削除して再実行してはならない。

外部 process と HTTP の timeout は次に固定する。timeout は process または request context を cancel し、終了を待ってから後続処理へ進む。shell を介さない。`git archive` 以外の child process の stdout / stderr は各 1 MiB を上限として読み、超過を当該 command の失敗とする。`git archive` の stdout は 128 MiB + 1 byte まで stream し、128 MiB 超過を `BUILD_FAILED`、stderr は 1 MiB 超過を `BUILD_FAILED` とする。

| 対象 | 1 回の timeout | timeout 時 error |
|------|----------------|------------------|
| `git rev-parse --show-toplevel` | `30s` | `INVALID_INPUT`。 |
| `git status --porcelain=v1 --untracked-files=all` | `30s` | `DIRTY_WORKTREE`。 |
| `git rev-parse HEAD` / local tag | `30s` | `TAG_MISMATCH`。 |
| `git show` commit timestamp | `30s` | `BUILD_FAILED`。 |
| `git archive` | `2m` | `BUILD_FAILED`。 |
| `gofmt -d` | `2m` | `FORMAT_FAILED`。 |
| `go test ./...` | `15m` | `TEST_FAILED`。 |
| 1 回の `go build` | `10m` | `BUILD_FAILED`。 |
| 生成 binary の `--version` | `10s` | `BUILD_FAILED`。 |
| GitHub metadata GET / POST / PATCH / DELETE | `30s` | 操作に対応する GitHub error code。 |
| asset 1 件の upload / download | `5m` | upload は `ASSET_UPLOAD_FAILED`、download は `ASSET_VERIFY_FAILED`。 |

---

<a id="release-asset-contract"></a>
**R3. Release asset 固定契約：**

安定版 Release は次の 6 asset だけを持つ。追加 asset、source archive の独自 upload、debug binary、`latest` alias を禁止する。GitHub が自動提供する source code archive は本表の Release asset に含めない。

| asset 名 | 内容 | mode |
|----------|------|------|
| `adlaire-ci-build-linux-amd64` | root `main` package を build し、basename dispatch で `builder` owner を起動する実行バイナリ。 | `0755` |
| `adlaire-ci-runner-linux-amd64` | root `main` package を build し、basename dispatch で `runner` owner を起動する実行バイナリ。 | `0755` |
| `adlaire-ci-api-linux-amd64` | root `main` package を build し、basename dispatch で `api` owner を起動する実行バイナリ。 | `0755` |
| `adlaire-ci-setup-linux-amd64` | root `main` package を build し、basename dispatch で `setup` owner を起動する実行バイナリ。 | `0755` |
| `admin-ui.tar.gz` | [`docs/details/admin.md` 詳細本文責務 §A1](admin.md#a1-管理-ui-静的ファイル境界) の 2 file を archive root 直下に持つ gzip 圧縮 tar。 | archive `0644` |
| `SHA256SUMS` | 前 5 asset の SHA-256 一覧。 | `0644` |

実行バイナリは component の単一 source file を直接 build して生成してはならない。4 binary はすべて checkout root の `main` package を入力とし、実行時の basename exact 一致によって owner component を選択する。部分一致、default fallback、未知 basename の builder 扱いを禁止する。

`SHA256SUMS` は `SHA256SUMS` 自身を除く 5 asset を filename の ASCII 昇順で並べ、各行を lowercase SHA-256 64 文字、ASCII space 2 文字、filename、LF の順とする。空行、comment、絶対 path、directory、重複 filename、未掲載 asset を禁止する。

---

<a id="release-build-contract"></a>
**R4. ビルド・再現性固定契約：**

release 実装は、checkout を変更せず、`--out` と OS temporary directory だけへ書き込む。Go source の整形確認、test、build は live checkout ではなく `--commit` の immutable source snapshot を使用し、次の順序に固定する。

1. mode `0700` の独立した source temporary directory A と B を作成する。
2. A / B それぞれに対し、`git archive --format=tar <commit>` を shell を介さず独立に 1 回ずつ実行する。Go 標準ライブラリ `archive/tar` で stdout を stream 展開し、entry 名は UTF-8 の checkout root 相対 path だけ、type は directory または regular file だけを許可する。絶対 path、空 segment、`.`、`..`、backslash、NUL、CR、LF、symlink、hardlink、device、FIFO、socket、重複 entry、A / B root 外への脱出、既存 path への上書きを `BUILD_FAILED` とする。regular file は Git archive mode の executable bit に応じて `0644` または `0755`、directory は `0755` とする。
3. A / B の展開 file set、mode、size、SHA-256 を relative path の ASCII 昇順で比較し、完全一致を確認する。`go.mod`、`main.go`、Release asset に必要な実装 artifact、admin 配布物は A / B 内の通常 file として再確認する。
4. A の snapshot 内にある `.go` regular file を relative path の ASCII 昇順で列挙し、`gofmt -d <path...>` を shell を介さず 1 回実行する。stdout が空、stderr が空、終了コード `0` の場合だけ合格とする。対象 0 件は `FORMAT_FAILED` とする。
5. `go test ./...` を A で実行し、終了コード `0` を確認する。
6. A / B それぞれで [R3 Release asset 固定契約](#release-asset-contract) の 4 binary を同一環境・同一引数で生成し、admin archive を生成する。
7. A と B の同名 asset を byte 単位で比較する。
8. R2 で保持した `--out` parent descriptor に対し、basename `<out>.tmp`をmode`0700`のsibling staging directoryとして排他的に作成し、一致したAの5成果物を固定modeでcopyして各fileをsync / closeする。
9. staging directory内に`SHA256SUMS`を生成してsync / closeし、全6fileの名前、type、mode、size、digestを再検証する。
10. `git status --porcelain=v1 --untracked-files=all`、`git rev-parse HEAD`、`git rev-parse refs/tags/<tag>^{commit}` を再実行し、実行開始時と同じ clean / commit / tag 条件を確認する。dirty は `DIRTY_WORKTREE`、HEAD / tag 不一致は `TAG_MISMATCH` とし、`<out>.tmp` を削除して `--out` と GitHub write を作成しない。
11. staging directoryをsync / closeしてから、保持済みparent descriptorに対するbasename相対操作でstaging directoryを`--out`へatomic renameし、同じparent descriptorをsyncする。

commit timestampは`git show -s --format=%ct <commit>`をshellを介さず1回実行し、stdoutがpositive decimal Unix seconds + LF、stderr空、exit`0`の場合だけ使用する。A / Bのbinary出力basenameは [R3 Release asset 固定契約](#release-asset-contract) のasset名とexact一致させる。atomic rename前の失敗は`<out>.tmp`だけを削除し、`--out`を作成しない。rename成功後の親directory sync失敗は、完全検証済みの`--out`を削除せず`OUTPUT_FAILED`とし、GitHub writeを開始しない。全caseでA / Bのtemporary directoryを削除する。cleanup失敗でも元error codeを維持し、pathを出力しない。残存`<out>.tmp`または`--out`は次回 [R1 CLI 固定契約](#release-cli-contract) の検証で`INVALID_INPUT`として拒否し、自動再利用または上書きしない。

Go format / test / build 環境は `CGO_ENABLED=0`、`GOOS=linux`、`GOARCH=amd64`、`GOFLAGS=`、`GOENV=off`、`GOWORK=off`、`GOTOOLCHAIN=local`、`GOPROXY=off`、`GOSUMDB=off` に固定する。build argv は `go`、`build`、`-trimpath`、`-buildvcs=false`、`-ldflags`、`-s -w -X main.binaryVersion=<tag>`、`-o`、`<output>`、`.` の 9 token とする。`<tag>` を含む ldflags 全体と出力 path はそれぞれ 1 argv 要素として渡し、shell 文字列へ連結しない。既存 environment の同名 key によって固定値を上書きしてはならない。`GOPATH`、`GOMODCACHE`、`PATH`、`HOME`、`TMPDIR` は process 開始時の値を A / B で同一に保持し、片方だけを変更してはならない。その他の `GO*` 環境変数は child へ継承しない。

各 binary は実行して `--version` を確認する。stdout は exact `<binary-name> <tag> go=<non-empty>` + LF、stderr は空、終了コードは `0` とする。binary name、tag、Go version の不一致、`V.0.0-dev` は公開前失敗とする。

`admin-ui.tar.gz` は Go 標準ライブラリ `archive/tar` と `compress/gzip` で生成する。entry は `adlaire-ci-sdk.js`、`index.html` の ASCII 昇順、typeflag `tar.TypeReg`、format `tar.FormatUSTAR`、mode `0644`、uid/gid `0`、uname/gname 空とし、PAX record、access time、change time、extended headerを持たせない。tar header のmodification timeとgzip headerのmodification timeは`--commit`のUTC commit timestamp秒精度、gzip compression levelは`9`、name/comment/extraは空、OS byteは`255`とする。archive entry の内容は対応する A / B source snapshot 内の対象 file と byte 単位で一致させる。

再現性比較が不一致の場合は、[R6 出力・副作用固定契約](#release-output-contract) に従って stderr に `release: NON_REPRODUCIBLE` + LF だけを出力する。不一致 asset 名、asset byte、temporary path、host user pathは出力せず終了コード`1`とし、GitHub write は開始しない。

---

<a id="release-publish-contract"></a>
**R5. GitHub Release 公開固定契約：**

GitHub API は `https://api.github.com` と、create response が返す scheme `https`、host exact `uploads.github.com` の upload URL だけを使用する。userinfo、fragment、未知 query key を禁止する。metadataとwrite操作のredirectは追従しない。asset再取得だけは後述する1回のdownload redirectを許可する。TLS verification を無効化してはならない。metadataとGitHub write requestは`Authorization: Bearer <token>`、`Accept: application/vnd.github+json`、`X-GitHub-Api-Version: 2022-11-28`、credential を含まない `User-Agent: adlaire-ci-release` を持つ。asset再取得requestは`Accept: application/octet-stream`へ置き換える。metadata response body は 4 MiB、error response body は 1 MiB を上限とし、超過は当該 GitHub 操作の失敗とする。asset download は response `Content-Length` が local size と一致する場合だけ local size + 1 byte まで読み、超過、短縮、未知長を `ASSET_VERIFY_FAILED` とする。

`--token-file` の credential は対象 repository への Release 作成・更新・削除と asset upload を許可されている必要がある。fine-grained personal access token または GitHub App token の最小 repository permission は `Contents: write` とする。対象 commit が repository の default branch に対して `.github/workflows/` 配下を追加または変更する場合は `Workflows: write` も必須とする。classic personal access token を使用する場合も同等の repository write 権限と、該当時は `workflow` scope を必須とする。権限不足を別 token への fallback、無認証 request、既存 Release 再利用で回避しない。read-only API の `401` / `403` は `GITHUB_READ_FAILED`、draft create の `401` / `403` / 権限不足による `404` は `DRAFT_CREATE_FAILED`、upload / PATCH / DELETE の権限不足は各操作の固定 error code とする。

`owner`、`repository`、default branch、`tag`、asset name、object SHA、release id、asset idは各path segmentまたはquery valueとして1回だけpercent encodeする。JSON requestは`Content-Type: application/json`、asset uploadは`Content-Type: application/octet-stream`とlocal sizeと同じ`Content-Length`を必須とし、chunked uploadを禁止する。次のmethod、endpoint、成功status、必須response以外のGitHub API requestを行ってはならない。

| 操作 | method / endpoint | 成功status / 必須response |
|------|-------------------|---------------------------|
| repository確認 | `GET /repos/{owner}/{repository}` | `200`、`full_name`が`owner/repository`とexact一致。`default_branch`は1〜255 bytesのUTF-8、前後空白なし、NUL / CR / LFなし、`refs/`非prefix、`.` / `..` 以外の非空値。 |
| remote default branch head確認 | `GET /repos/{owner}/{repository}/git/ref/heads/{default-branch}` | `200`、`ref`が`refs/heads/<default-branch>`とexact一致、`object.type="commit"`、`object.sha`が40文字lowercase hex。`404`は`TAG_MISMATCH`。 |
| remote default branch包含確認 | `GET /repos/{owner}/{repository}/compare/{commit}...{default-head-sha}` | `200`、`base_commit.sha=<commit>`、`merge_base_commit.sha=<commit>`、`status`が`ahead`または`identical`。`404` / `409` または値不一致は`TAG_MISMATCH`。 |
| remote tag確認 | `GET /repos/{owner}/{repository}/git/ref/tags/{tag}` | `200`、`object.type`が`commit`または`tag`、`object.sha`が40文字lowercase hex。`404`は`TAG_MISMATCH`。 |
| annotated tag解決 | `GET /repos/{owner}/{repository}/git/tags/{sha}` | `200`、`object.type`が`commit`または`tag`、最大5回でcommitへ到達し、最終SHAが`--commit`と一致。cycle、6段目、未知typeは`TAG_MISMATCH`。 |
| 既存Release確認 | `GET /repos/{owner}/{repository}/releases/tags/{tag}` | `404`だけを不在とする。`200`は`RELEASE_EXISTS`、その他は`GITHUB_READ_FAILED`。 |
| draft作成 | `POST /repos/{owner}/{repository}/releases` | `201`、positive integer `id`、`draft=true`、`prerelease=false`、tag / target / name / body一致、安全な`upload_url`と`html_url`。 |
| asset upload | create responseの`upload_url`へquery `name=<asset-name>`だけを追加した`POST` | `201`、positive integer `id`、name、size一致、state=`uploaded`。 |
| asset一覧 | `GET /repos/{owner}/{repository}/releases/{release-id}/assets?per_page=100&page=1` | `200`、6件、id / name / size一意、次pageを示す`Link`なし。 |
| asset再取得 | `GET /repos/{owner}/{repository}/releases/assets/{asset-id}`、`Accept: application/octet-stream` | `200`、または許可download redirect 1回後の`200`。最終responseのContent-Length / byte数 / SHA-256一致。 |
| Release再取得 | `GET /repos/{owner}/{repository}/releases/{release-id}` | `200`、id、tag、target、name、body、draft、prerelease、asset件数一致。 |
| 正式公開 | `PATCH /repos/{owner}/{repository}/releases/{release-id}` body `{"draft":false,"prerelease":false}` | `200`、`draft=false`、`prerelease=false`、その他の固定値一致。 |
| draft cleanup | `DELETE /repos/{owner}/{repository}/releases/{release-id}` | `204`、body 0 byte。 |

`upload_url`のtemplate suffix `{?name,label}`はexact一致で1回だけ除去し、その他のtemplate、既存query、label queryを拒否する。`html_url`はscheme`https`、host exact`github.com`、userinfo / query / fragmentなし、path exact`/{owner}/{repository}/releases/tag/{tag}`とし、成功JSONの`release_url`に使用する。GitHub JSONの未知keyは無視できるが、上表の必須key欠落、型不一致、値不一致は操作別errorとする。

remote default branch head SHA、remote tag peeled commit、既存 Release 不在は、実行前の初回確認値を cache したまま公開に使用しない。GitHub write 直前に repository、default branch head、compare、tag、既存 Release 不在をすべて再取得し、初回確認後の default branch head 進行は compare 再合格時だけ許可する。tag peeled commit の変更、`--commit` の default branch 非包含化、repository / default branch の不一致、既存 Release 出現は GitHub write を開始せず失敗する。

許可download redirectはstatus`302`、Locationのscheme`https`、host exact`release-assets.githubusercontent.com`、userinfo / fragmentなしの場合だけ1回追従する。redirect先requestへ`Authorization`、`Cookie`、`Referer`、GitHub API version headerを転送せず、`User-Agent`だけを送る。2回目のredirect、相対Location、未知host、Location欠落は`ASSET_VERIFY_FAILED`とする。

公開順序は次に固定する。

1. local checkout の clean、HEAD、local tag と、remote repository、default branch head、`--commit` 包含、tag peeled commit、既存 Release 不在を再確認する。この確認は [R2 実行前固定契約](#release-precondition-contract) の事前確認後に状態が変化した競合を検出する 2 回目の確認であり、省略しない。
2. `tag_name=<tag>`、`target_commitish=<commit>`、`name=<tag>`、`body=<notes-file byte>`、`draft=true`、`prerelease=false`、`generate_release_notes=false` で draft Release を1件作成する。
3. draft 作成直後に remote default branch head、`--commit` 包含、tag peeled commit、Release metadata を再取得する。不一致は `PUBLISH_FAILED` とし、asset upload を開始せず draft cleanup を実行する。
4. [R3 Release asset 固定契約](#release-asset-contract) の順で6 assetを1件ずつuploadする。
5. draft Releaseのasset一覧を再取得し、名前、件数、sizeをlocal成果物と一致させる。
6. 各uploaded assetをGitHubから一時directoryへ再取得し、local SHA-256と一致させる。
7. repository、remote default branch head、`--commit` 包含、tag peeled commit、Releaseのtag、target commit、title、notes、draft、prerelease、asset件数を再取得して一致させる。不一致は `PUBLISH_FAILED` とし、publish せず draft cleanup を実行する。
8. `draft=false`、`prerelease=false`へ更新して正式公開する。
9. 公開済みReleaseを再取得し、`draft=false`、`prerelease=false`、asset 6件を確認する。

upload は `Content-Type: application/octet-stream` とし、asset name は URL query へ 1 回だけ percent encode する。既存 asset の削除・置換、同名 upload の retry、公開後の本文変更を行わない。通信失敗時の自動 retry は、response body を送信していない read-only GET に限り最大2回、1秒、2秒で許可する。create、upload、PATCH、DELETE を自動再送してはならない。

draft 作成後、正式公開前に失敗した場合は、作成した draft の id が確定している場合だけ DELETE を1回実行する。DELETE成功後は終了コードを元の失敗種別のままとする。DELETE失敗時は終了コード`1`、error code`DRAFT_CLEANUP_FAILED`とし、URL、draft id、元のresponse body、元のerror文字列をstdout / stderrへ出さない。tagとlocal成果物は削除しない。正式公開後は自動削除・rollbackを行わない。

---

<a id="release-output-contract"></a>
**R6. 出力・副作用固定契約：**

成功時のstdoutはUTF-8 JSON object 1行とLFだけとし、key順を次に固定する。

```json
{"tag":"V.X.N","commit":"40-hex-sha","release_url":"https://github.com/owner/repository/releases/tag/V.X.N","assets":["adlaire-ci-build-linux-amd64","adlaire-ci-runner-linux-amd64","adlaire-ci-api-linux-amd64","adlaire-ci-setup-linux-amd64","admin-ui.tar.gz","SHA256SUMS"],"published":true}
```

実値へ置換した後もkey順とasset順を維持する。成功時stderrは空とする。失敗時stdoutは空、stderrは`release: <error-code>` + LFの1行とする。error codeは [R1 CLI 固定契約](#release-cli-contract) の固定表に列挙した値だけを許可する。追加説明、path、URL、draft id、HTTP body、Go error、command outputをstdout / stderrへ出してはならない。

許可する永続副作用は、`--out`への検証済み成果物作成とGitHub draft / asset /正式公開だけとする。checkout、Git index、tag、branch、repository設定、issue、PR、workflow、secret、environment、systemdを変更してはならない。

---

<a id="release-acceptance-contract"></a>
**R7. 実装受け入れ条件：**

| 観点 | 合格条件 |
|------|----------|
| CLI | 全 option、検証順、help、version、終了コード、固定errorがfixtureと一致する。 |
| source | clean checkout、HEAD、local tag、remote tag、remote default branch ancestor、immutable source snapshot、GitHub write 前の local / remote 再検証が固定契約と一致する。 |
| build | Go test、format、固定build環境、version注入、4 binaryのbasename dispatchが一致する。 |
| reproducibility | 2回生成した5 assetがbyte単位で一致し、不一致時にGitHub writeが0件である。 |
| archive | file set、順序、mode、owner、timestamp、gzip headerが固定値に一致する。 |
| checksum | 5 asset、ASCII順、lowercase SHA-256、space 2文字、LF、自己行なしが一致する。 |
| GitHub | 既存Releaseを上書きせず、draftで全assetを再取得検証してから公開する。 |
| failure | 公開前失敗ではdraftを1回だけ削除し、tag、checkout、既存Releaseを変更しない。 |
| control | timeout、response size、再試行、並行実行、既存Release、partial failureが [R2 実行前固定契約](#release-precondition-contract) と [R5 GitHub Release 公開固定契約](#release-publish-contract) の固定境界に一致する。 |
| secret | token平文がstdout、stderr、notes、asset、URL、fixture expectedへ出ない。 |
| evidence | [`docs/details/fixture.md` fixture 証跡責務 Release fixture 固定契約](fixture.md#release-fixture-contract) の全fixtureが合格する。 |

実装完了と状態遷移の判定は [`docs/SPEC.md` ポリシー責務 §0a](../SPEC.md#0a-仕様成熟度ポリシー)、現在状態は [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務を正本とし、[`docs/details/release.md`](release.md) 詳細本文責務で再定義しない。
