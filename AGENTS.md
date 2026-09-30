# Build-Scripts - 最上位ルールブック

## 0. 絶対原則

[AGENTS.md](AGENTS.md) と [docs/SPEC.md](docs/SPEC.md) は、本リポジトリにおける最上位文書である。

[AGENTS.md](AGENTS.md) は、作業ルール、承認、Git 操作、Pull Request 作成、レビュー対応、検証手順、エージェント実行手順の最上位ルールブックである。

[docs/SPEC.md](docs/SPEC.md) は、仕様、方針、ポリシー、正本参照先、禁止事項、リリース判断、実装着手可否の最上位仕様書である。

このリポジトリで作業するすべてのエージェントは、最初に [AGENTS.md](AGENTS.md) と [docs/SPEC.md](docs/SPEC.md) の両方を必ず読む。両方の読了が完了するまで、調査、設計、仕様改訂、実装、検証、ファイル操作、生成物更新、Git 操作、Pull Request 作成、レビュー対応を含むリポジトリに関するすべての作業を開始してはならない。片方だけの読了を、必要な確認の完了と扱ってはならない。

本リポジトリの仕様判断は、[`docs/SPEC.md` 責務文書構成表](docs/SPEC.md#document-responsibility-map) で対象判断を所有する責務正本を確定して行う。[AGENTS.md](AGENTS.md) で仕様責務の分担を再定義してはならない。

[AGENTS.md](AGENTS.md) と他ファイルが作業ルール上矛盾する場合は、[AGENTS.md](AGENTS.md) を正とする。

仕様文書間で矛盾する場合は、[`docs/SPEC.md` 責務文書構成表](docs/SPEC.md#document-responsibility-map) で該当判断を所有する責務正本を正とする。

[docs/SPEC.md](docs/SPEC.md)、[docs/ROADMAP.md](docs/ROADMAP.md)、[docs/DETAIL_INDEX.md](docs/DETAIL_INDEX.md)、または該当する owner component 別の [docs/details/*.md](docs/details/) と実装ファイルが仕様上矛盾する場合は、仕様と実装の不整合として扱う。仕様を変更する場合は、先に該当する仕様書を改訂し、その内容に基づいて実装を更新する。

---

<a id="agents-approval-rules"></a>

## 1. 承認ルール

変更作業では、承認工程を省略してはならない。

ファイル作成、編集、移動、削除、リネーム、整形、生成物更新、Git 操作など、リポジトリ内の状態を変更する作業は、いかなる場合もユーザーから事前承認を得るまで実行しない。

変更作業前には、いかなる場合も以下を提示する。

- 変更対象
- 変更内容
- 影響範囲

ユーザー承認は、ユーザーの返信に `承認` という単語が明示された場合のみ有効とする。

`OK`、`はい`、`お願いします`、`進めて`、その他の類似表現は、変更作業の承認として扱わない。

承認後は、提示済みの変更作業内容の範囲内でのみ作業する。

承認後の作業中も、現在の作業が承認済み範囲内かを継続して確認する。

作業中に新たな不整合、改善候補、設定差分を発見した場合でも、承認済み範囲外であれば編集、移動、削除、リネーム、生成物更新、設定変更を行ってはならない。

提示済み範囲を超える変更が必要になった場合は、追加の変更内容を提示し、別途 `承認` を得る。

追加承認を求める場合は、以下を提示する。

- 追加変更対象
- 追加変更内容
- 影響範囲
- 今その変更を行う必要性

検証、読取、検索、差分確認など、リポジトリ状態を変更しない調査は、承認済み作業の判断材料として実行してよい。

承認済み変更の整合に直接必要な [docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md) 等の更新は、最初に提示した影響範囲に含まれている場合に限り、追加承認なしで行ってよい。

---

## 2. 仕様書管理ルール

仕様変更では、最初に [`docs/SPEC.md` 責務文書構成表](docs/SPEC.md#document-responsibility-map) で対象判断の責務正本を確定し、[`docs/SPEC.md` 方針責務 §4.2a](docs/SPEC.md#sec-4-2a) の記載範囲と禁止事項を適用する。[AGENTS.md](AGENTS.md) で同じ判断対象の正本分担または禁止事項を再定義してはならない。

仕様変更の影響確認は、[文書整合ルール](#agents-document-consistency-rules) に従う。

仕様変更の編集前と編集後に、[docs/SPEC.md 方針責務 §4.2a 仕様全般重複記載禁止原則](docs/SPEC.md#spec-global-no-duplicate-principle) への適合を確認する。完全一致する本文だけでなく、言い換え、要約、部分転載、表と本文の再掲、owner と collaborator 間の意味上の重複を確認する。確認では、判断対象、唯一の責務正本、重複候補の所在、削除または参照化の処置を特定し、未解消件数が 0 になるまで仕様変更を完了扱いにしてはならない。

文書整理だけを目的とする変更では、機能契約、現在状態、実装可否、Phase、将来計画を変更してはならない。ただし、実装と必須証跡を確認した結果、既存の状態記載が事実と矛盾すると判明した場合は、承認済み範囲内で [docs/ROADMAP.md](docs/ROADMAP.md) の現在状態を事実へ一致させる。

標準ディレクトリ構成は [docs/SPEC.md 方針責務 §4.3](docs/SPEC.md#sec-4-3)、実在所在は [docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md) を確認する。未作成 path を実在ファイルとして扱ってはならない。

---

## 3. 実装管理ルール

実装作業前には、対象機能について [docs/SPEC.md](docs/SPEC.md)、[docs/ROADMAP.md](docs/ROADMAP.md)、[docs/DETAIL_INDEX.md](docs/DETAIL_INDEX.md)、対象 owner / collaborator 詳細本文、[docs/details/fixture.md](docs/details/fixture.md)、実在ファイルを確認する。

実在ファイルの確認には hidden fixture を含めて列挙できる `rg --files --hidden -g '!.git/**'` を使用する。

新規実装の着手可否は、[`docs/SPEC.md` ポリシー責務 §0a](docs/SPEC.md#policy-spec-maturity) の実装可否、[`docs/SPEC.md` ポリシー責務 §0 完全仕様詳細化義務](docs/SPEC.md#complete-detail-specification-gate)、[`docs/SPEC.md` ポリシー責務 §0d](docs/SPEC.md#policy-spec-freeze) の凍結条件、[`docs/SPEC.md` ポリシー責務 §0f](docs/SPEC.md#policy-phase-unit) の active Phase 条件、[`docs/ROADMAP.md`](docs/ROADMAP.md) 状態・計画責務の現在状態によって判定する。

実装中に未定義の入力、出力、状態、異常系、セキュリティ条件、検証条件、fixture、fake、expected、完了条件を発見した場合は、実装判断で補完せず、先に該当する責務正本を改訂する。

完全仕様詳細化が未完了の機能、または一定の仕様だけを固定した機能は、実装着手不可として扱う。実装しながら仕様を決めること、既存実装やテスト結果に合わせて仕様を後追い確定すること、実装者判断で未定義事項を補うことを行ってはならない。

Go 実装の標準配置は [docs/SPEC.md 方針責務 §4.3](docs/SPEC.md#sec-4-3) を参照する。[`main.go`](main.go) は起動入口 artifact、[`components/builder.go`](components/builder.go)、[`components/runner.go`](components/runner.go)、[`components/api.go`](components/api.go) はそれぞれ `builder`、`runner`、`api` owner component の標準 Go 実装 artifact、`components/admin.go` は `admin` owner component の CLI 管理クライアント用 Go 実装 artifact、`components/mcp.go` は `mcp` owner component の Go 実装 artifact、[`admin/adlaire-ci-sdk.js`](admin/adlaire-ci-sdk.js) と [`admin/index.html`](admin/index.html) はそれぞれ `sdk`、`ui` owner component の標準管理クライアント実装 artifact として扱う。owner component と実装 artifact を同一概念として扱ってはならない。各実装 artifact の現在状態は [docs/ROADMAP.md](docs/ROADMAP.md)、実在所在は [docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md) を参照する。

実装変更後は、変更範囲に応じて構文確認、単体確認、実行確認、生成物確認、異常系確認、必須 fixture 確認を行う。

実装作業では、[`docs/SPEC.md` ポリシー責務 §0f](docs/SPEC.md#policy-phase-unit) に従い、active Phase 全体が完了するまで同一作業ブランチで実装、検証、不整合修正、再検証を繰り返す。

active Phase 内に未実装、未検証、仕様不整合、証跡不足、状態更新不足が残る場合は、実装作業を完了扱いにしてはならない。既存 Pull Request がある場合も Phase 全体完了まで work in progress として扱い、review ready、merge 可能、または完了済みと報告してはならない。

Go 実装では、対象ファイルに `gofmt -l ...` を実行し、Go module が存在する場合は `go test ./...` を実行する。JavaScript 系実装では、[`docs/SPEC.md` 方針責務 §4 技術方針表](docs/SPEC.md#direction-technical) に従い、Deno stable runtime で対象 JavaScript file に `deno check ...` を実行する。Node.js、npm、bundler、transpiler を JavaScript 系実装の標準検証コマンドとして代替使用してはならない。実行できない確認は、未実施理由を Pull Request 本文へ記録する。

API、SDK、UI のいずれかを変更する場合は、対応する endpoint、SDK method、UI 操作、状態副作用、認証・認可、成功後再取得、失敗時固定、fixture 証跡を同じ変更で確認する。

`実装済み` への状態変更は、[`docs/SPEC.md` ポリシー責務 §0a 実装完了条件](docs/SPEC.md#implementation-completion-transition) と [`docs/SPEC.md` 方針責務 §4.8](docs/SPEC.md#sec-4-8) に従って判定する。

---
<a id="agents-git-operations"></a>

## 4. Git 運用ルール

`main` は保護対象ブランチとする。

`main` への直接 push を禁止する。

変更作業は、必ず作業ブランチで行う。

作業ブランチは一本化し、ドキュメント変更作業と実装変更作業の両方で同じ作業ブランチを使用する。

同一目的、同一仕様領域、同一ファイル群に対する変更は、必ず 1 本の作業ブランチと 1 本の Pull Request にまとめる。

同一目的の変更を複数の積み上げ Pull Request に分割してはならない。

複数の Pull Request に分ける場合は、変更対象ファイル、責務、merge 順序が明確に分離でき、相互に同一ファイルを編集せず、片方だけが merge されても仕様矛盾、参照切れ、状態不一致、未定義の依存関係が発生しない場合に限る。

既存の open Pull Request と同じファイルまたは同じ仕様領域を変更する必要がある場合は、別 Pull Request を作成せず、既存 Pull Request へ変更を統合する。

積み上げ Pull Request、同一ファイル編集の並行 Pull Request、merge 順序依存の Pull Request、または GitHub 上で `DIRTY` / conflict 状態の Pull Request が発生した場合は、競合解消作業として扱う。競合解消作業では、最新 `origin/main` から一本化ブランチを作成するか、最も包括的な既存 Pull Request の branch を統合先とし、必要な変更を 1 本の Pull Request に統合する。

一本化後、重複する既存 Pull Request は、統合先 Pull Request を明記したコメントを残して close する。

<a id="git-conflict-prevention-check"></a>
競合防止確認として、作業開始前と Pull Request 作成前に以下を必ず実行する。

1. `git fetch origin`
2. `gh pr list --state open --json number,title,headRefName,baseRefName,mergeStateStatus,url`
3. `git diff --name-status origin/main...HEAD`

この競合防止確認で、同一ファイル、同一仕様領域、同一責務、または merge 順序依存の open Pull Request が見つかった場合は、新規 Pull Request を作成してはならない。既存 Pull Request への統合、または一本化 Pull Request への集約を先に完了する。

Pull Request の `mergeStateStatus` が `DIRTY`、`UNKNOWN`、または確認不能の場合は、merge 可能と報告してはならない。`UNKNOWN` の場合は GitHub の再計算後に再確認し、最終的に `CLEAN` を確認する。

競合解消時は、競合マーカーの除去だけで完了としてはならない。`git diff --check`、競合マーカー検索、変更対象文書の正本参照先確認、open Pull Request 一覧確認を確認条件とする。

`main` への反映は、Pull Request 経由で行う。

Pull Request の merge はユーザーが行う。

エージェントは Pull Request の merge を行ってはならない。

GitHub リポジトリ設定の変更は変更作業として扱い、事前に変更対象、変更内容、影響範囲を提示し、ユーザーから `承認` を得るまで実行してはならない。

Build-Scripts の標準 GitHub リポジトリ設定は以下とする。

- visibility: `public`
- default branch: `main`
- delete branch on merge: `true`
- allow merge commit: `true`
- allow squash merge: `true`
- allow rebase merge: `true`
- allow auto merge: `false`
- allow update branch: `false`
- issues: `true`
- projects: `true`
- wiki: `true`
- discussions: `false`
- secret scanning: `enabled`
- secret scanning push protection: `enabled`
- Dependabot security updates: `disabled`
- main branch protection: 設定対象（GitHub 設定の初期適用方針を適用）

GitHub 設定の初期適用方針は以下とする。

- `main` branch protection は、初期標準として Pull Request 必須、force push 禁止、branch deletion 禁止を設定する。
- `main` branch protection の required approvals は初期値 `0` とする。
- required approvals を `1` へ引き上げる場合は、運用安定後の別変更として、変更対象、現在値、標準値、影響範囲を提示して承認を得る。

GitHub 設定を確認する場合は、少なくとも以下を確認する。

- `gh api repos/fqwink/Build-Scripts` で、`delete_branch_on_merge`、`allow_auto_merge`、`allow_update_branch`、`allow_merge_commit`、`allow_squash_merge`、`allow_rebase_merge`、`has_issues`、`has_projects`、`has_wiki`、`has_discussions`、`security_and_analysis` を確認する。
- `gh api repos/fqwink/Build-Scripts/branches/main/protection` で、`main` branch protection を確認する。
- `main` branch protection の確認で `Branch not protected` が返る場合は、未設定として扱う。

GitHub 設定を変更する前には、以下を必ず提示する。

- 変更対象
- 現在値
- 標準値
- 影響範囲

GitHub 設定を変更した後は、GitHub API で再取得し、[AGENTS.md](AGENTS.md) の標準設定との差分がないかを確認して報告する。

標準 GitHub リポジトリ設定のうち、自動化に関わる設定が未確認の場合は、現在の設定状態を確認する。

自動化に関わる設定が未設定または標準値と異なる場合は、変更対象、変更内容、影響範囲を提示し、ユーザーから `承認` を得たうえで標準値へ設定する。

`delete_branch_on_merge=true` は remote branch 自動削除の必須自動化設定であり、GitHub 設定の初期適用で即時設定する。その他の自動化設定が [docs/SPEC.md](docs/SPEC.md) または本ルールブックで標準化された場合も同様に扱う。

エージェントは、ユーザー承認なしに GitHub リポジトリ設定を変更、無効化、初期化してはならない。

remote branch は、GitHub リポジトリ設定 `delete_branch_on_merge=true` により、Pull Request merge 後に GitHub 側で自動削除する。

remote branch 自動削除の対象は merge 済み Pull Request の head branch に限定する。local branch は GitHub 側の自動削除では削除されないため、エージェントはユーザーによる merge 完了を確認できた場合に限り、追加承認なしで同名の local branch を削除する。

local branch 削除では、`main` へ移動した後に対象 local branch を削除する。

Pull Request merge 後のローカル同期は、以下の手順を標準とする。

1. `git fetch --prune`
2. `git switch main`
3. `git merge --ff-only origin/main`
4. merge 済み Pull Request の head branch と同名の local branch を `git branch -d <branch>` で削除する
5. `git status --short --branch` で `main` と `origin/main` が一致し、作業ツリーが clean であることを確認する

このローカル同期手順で fast-forward できない場合、merge 対象やローカル変更の状態を確認し、勝手に履歴を書き換えてはならない。

`main`、merge 未完了の作業ブランチ、merge 状態を確認できないブランチ、Pull Request と対応しないブランチは削除してはならない。

`.gitignore` は作成・使用しない。

`.gitignore` が必要になる生成物、一時ファイル、実行時データ、ビルド成果物が発生した場合は、除外設定で隠蔽せず、生成先、運用、または実装を見直す。

承認済み変更作業が完了した場合、エージェントはユーザーからの追加指示および追加承認なしで、作業ブランチでのcommit、remoteへのpush、Pull Requestの作成または既存Pull Requestの更新まで自動実行する。

Pull Request 作成自動化は、[承認ルール](#agents-approval-rules) の承認済み範囲と、[本節](#agents-git-operations) の `main` 直接 push 禁止およびエージェントによる merge 禁止を例外なく適用する。

Pull Request 作成前には、変更内容に応じて以下を確認する。

- [競合防止確認](#git-conflict-prevention-check) を再実行し、変更対象が承認済み範囲内であり、統合先と重複 PR の状態が確定し、Pull Request の `mergeStateStatus` が `CLEAN` であることを確認する。
- 文書変更では、`rg` で不要になった名称、矛盾参照、不要になったファイル名が残っていないか確認する。
- 文書変更では、`git diff --stat` で変更範囲を確認する。
- ファイル追加、削除、リネームを含む場合は、`git diff --cached --summary` で Git 上の扱いを確認する。
- 実装変更では、[`docs/SPEC.md` ポリシー責務 §0f](docs/SPEC.md#policy-phase-unit) に従い、対象 Phase 全体が完了し、未実装、未検証、仕様不整合、証跡不足、状態更新不足が残っていないことを確認する。
- 実装変更では、対象言語に応じた構文確認を行う。Go 実装では `gofmt -l ...` を標準の整形確認とし、Go module が存在する場合は `go test ./...` を標準の確認とする。JavaScript 系実装では Deno stable runtime の `deno check ...` を標準の確認とする。
- 実装変更では、変更した実装が実行可能な場合は対象スクリプトの実行確認または生成物確認を行う。実行不能な場合は理由を Pull Request 本文に記録する。
- 仕様変更では、[文書整合ルール](#agents-document-consistency-rules) に従って責務正本、索引、デザイン、実装への影響を確認する。

Pull Request 本文には、少なくとも以下を記載する。

- `Summary`
- `Verification`
- 競合防止確認
- 未実施の確認がある場合は、その理由

---

## 5. 外部依存変更ルール

外部依存の採否、禁止条件、例外条件、許可範囲、許可外部ライブラリ一覧は、[docs/SPEC.md 方針責務 §4.1](docs/SPEC.md#sec-4-1) と [docs/SPEC.md ポリシー責務 §4](docs/SPEC.md#policy-dependencies) だけを正本とする。[AGENTS.md](AGENTS.md) で同じ方針または許可条件を再定義してはならない。

外部依存を追加、削除、更新、置換する前に、対象実装、[`go.mod`](go.mod)、配布物、セットアップ、検証手段への影響と、[docs/SPEC.md ポリシー責務 §4](docs/SPEC.md#policy-dependencies) の許可外部ライブラリ一覧を確認する。

許可外部ライブラリ一覧にない依存を実装へ追加してはならない。追加が必要な場合は、依存名、採用理由、代替困難性、責務範囲、影響範囲、保守・削除方針、検証条件、[docs/SPEC.md](docs/SPEC.md) の変更内容を提示し、実装変更と仕様変更の両方について事前承認を得る。

外部依存の追加、削除、更新、置換は変更作業として扱う。変更後は、[docs/SPEC.md ポリシー責務 §4](docs/SPEC.md#policy-dependencies) の許可一覧、[`go.mod`](go.mod)、実装 import、配布・セットアップ手順、検証結果が一致していることを確認する。

---

<a id="agents-document-consistency-rules"></a>

## 6. 文書整合ルール

[docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md) は、リポジトリ内の文書・実装ファイルの役割を示す索引として崩してはならない。

ファイル名、正本参照先、実装 artifact の追加・削除・リネームが発生した場合は、[docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md) の更新要否を確認する。

[docs/SPEC.md](docs/SPEC.md)、[docs/ROADMAP.md](docs/ROADMAP.md)、[docs/DETAIL_INDEX.md](docs/DETAIL_INDEX.md)、または owner component 別の [docs/details/*.md](docs/details/) を改訂した場合は、[docs/DESIGN.md](docs/DESIGN.md)、[docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md)、実装ファイルへの影響を確認する。

[docs/DESIGN.md](docs/DESIGN.md) を改訂した場合は、生成静的 Web サイトの変更では [`components/builder.go`](components/builder.go) 内の HTML / CSS / JavaScript / theme component テンプレート、標準管理 UI の変更では [`admin/index.html`](admin/index.html) の HTML / CSS との整合性を確認する。

仕様化済み項目を実装した場合は、[docs/ROADMAP.md](docs/ROADMAP.md) の現在状態、[docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md) の「実装ファイル一覧」、[docs/DETAIL_INDEX.md](docs/DETAIL_INDEX.md) の対応表、owner component 別の [docs/details/*.md](docs/details/) の検証条件、実装ファイルの存在を整合させる。
