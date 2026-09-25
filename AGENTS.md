# Build-Scripts - 最上位ルールブック

## 0. 絶対原則

[AGENTS.md](AGENTS.md) と [docs/SPEC.md](docs/SPEC.md) は、本リポジトリにおける最上位文書である。

[AGENTS.md](AGENTS.md) は、作業ルール、承認、Git 操作、Pull Request 作成、レビュー対応、検証手順、エージェント実行手順の最上位ルールブックである。

[docs/SPEC.md](docs/SPEC.md) は、仕様、方針、ポリシー、正本参照先、禁止事項、リリース判断、実装着手可否の最上位仕様書である。

このリポジトリで作業するすべてのエージェントは、調査、設計、仕様改訂、実装、検証、Git 操作、Pull Request 作成、レビュー対応を含む全作業において、最初に [AGENTS.md](AGENTS.md) と [docs/SPEC.md](docs/SPEC.md) の両方を必ず読む。

[AGENTS.md](AGENTS.md) と [docs/SPEC.md](docs/SPEC.md) の読了を完了するまで、調査、設計、仕様改訂、実装、検証、ファイル作成、編集、移動、削除、リネーム、整形、生成物更新、Git 操作、Pull Request 作成、レビュー対応を含む、リポジトリに関するすべての作業を開始してはならない。

[AGENTS.md](AGENTS.md) または [docs/SPEC.md](docs/SPEC.md) の片方だけを確認した状態で、作業判断に必要な確認を完了したと扱ってはならない。

本リポジトリの仕様判断は、方針、ポリシー、状態語彙、状態定義、状態遷移条件、禁止事項、リリース判断、実装着手可否は [docs/SPEC.md](docs/SPEC.md) 方針責務・ポリシー責務、生成 HTML のデザイン関係は [docs/DESIGN.md](docs/DESIGN.md) デザイン責務、実装 artifact と各機能の現在状態、Phase、機能インベントリ、将来計画は [docs/ROADMAP.md](docs/ROADMAP.md) 状態・計画責務、詳細仕様参照入口、共通固定値、owner / collaborator 対応表は [docs/DETAIL_INDEX.md](docs/DETAIL_INDEX.md) 詳細仕様入口責務、owner component 別の詳細本文は [docs/details/*.md](docs/details/) 詳細本文責務、fixture、expected、fake、実装検証証跡は [docs/details/fixture.md](docs/details/fixture.md) fixture 証跡責務、文書と実装ファイルの実在所在は [docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務を正本として行う。

[docs/DESIGN.md](docs/DESIGN.md) は、生成静的 Web サイトのデザイン関係の正本である。機能仕様、運用仕様、API 仕様、CI 仕様、状態語彙、状態定義、実装 artifact と機能の現在状態の正本ではない。

[docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md) は、文書・実装ファイルの役割を整理する索引である。仕様本文の正本ではない。

[AGENTS.md](AGENTS.md) と他ファイルが作業ルール上矛盾する場合は、[AGENTS.md](AGENTS.md) を正とする。

[docs/SPEC.md](docs/SPEC.md) 方針責務・ポリシー責務と他ファイルが仕様、方針、ポリシー、状態語彙、状態定義、状態遷移条件、正本参照先、禁止事項、リリース判断、実装着手可否で矛盾する場合は、[docs/SPEC.md](docs/SPEC.md) 方針責務・ポリシー責務を正とする。実装 artifact と各機能へ割り当てた現在状態は [docs/ROADMAP.md](docs/ROADMAP.md) 状態・計画責務を正とする。ただし、生成 HTML のデザイン関係は [docs/DESIGN.md](docs/DESIGN.md) デザイン責務を正とする。

[docs/SPEC.md](docs/SPEC.md)、[docs/ROADMAP.md](docs/ROADMAP.md)、[docs/DETAIL_INDEX.md](docs/DETAIL_INDEX.md)、または該当する owner component 別の [docs/details/*.md](docs/details/) と実装ファイルが仕様上矛盾する場合は、仕様と実装の不整合として扱う。仕様を変更する場合は、先に該当する仕様書を改訂し、その内容に基づいて実装を更新する。

---

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

仕様変更では、最初に対象判断の正本を確定する。

| 判断対象 | 正本 |
|----------|------|
| 方針、ポリシー、状態語彙、状態定義、状態遷移条件、禁止事項、実装着手可否 | [docs/SPEC.md](docs/SPEC.md) |
| 実装 artifact と機能の現在状態、Phase、機能インベントリ、将来計画 | [docs/ROADMAP.md](docs/ROADMAP.md) |
| 詳細仕様入口、共通固定値、owner / collaborator 対応 | [docs/DETAIL_INDEX.md](docs/DETAIL_INDEX.md) |
| owner component 別の入出力、状態、処理順序、異常系、検証条件 | [docs/details/*.md](docs/details/) |
| fixture、expected、fake、assertion、実装検証証跡 | [docs/details/fixture.md](docs/details/fixture.md) |
| 生成 HTML のデザイン | [docs/DESIGN.md](docs/DESIGN.md) |
| 文書、実装、testdata、未作成 path の実在所在 | [docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md) |

[docs/DETAIL_INDEX.md](docs/DETAIL_INDEX.md) と owner component 別の [docs/details/*.md](docs/details/) に、方針、ポリシー、状態語彙、現在状態、Phase、将来計画、Git 運用、PR 分割判断を記載してはならない。

[docs/ROADMAP.md](docs/ROADMAP.md) に、状態語彙、状態定義、状態遷移条件、endpoint、schema、SDK method、UI DOM、処理順序、fixture assertion を再定義してはならない。

[docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md) に、方針、現在状態、実装詳細、検証条件を記載してはならない。

実装ファイルの存在だけで `実装済み` と判定してはならない。判定には [docs/SPEC.md](docs/SPEC.md) の状態遷移条件、[docs/ROADMAP.md](docs/ROADMAP.md) の現在状態、owner component 詳細本文、[docs/details/fixture.md](docs/details/fixture.md) の必須証跡を使用する。

仕様変更では、変更した責務正本から参照される [docs/ROADMAP.md](docs/ROADMAP.md)、[docs/DETAIL_INDEX.md](docs/DETAIL_INDEX.md)、[docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md)、[docs/DESIGN.md](docs/DESIGN.md)、owner / collaborator 詳細本文、fixture 証跡、実装ファイルへの影響を確認する。

文書整理だけを目的とする変更では、機能契約、現在状態、実装可否、Phase、将来計画を変更してはならない。ただし、実装と必須証跡を確認した結果、既存の状態記載が事実と矛盾すると判明した場合は、承認済み範囲内で [docs/ROADMAP.md](docs/ROADMAP.md) の現在状態を事実へ一致させる。

標準ディレクトリ構成は [docs/SPEC.md 方針責務 §4.3](docs/SPEC.md#sec-4-3)、実在所在は [docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md) を確認する。未作成 path を実在ファイルとして扱ってはならない。

---

## 3. 実装管理ルール

実装作業前には、対象機能について [docs/SPEC.md](docs/SPEC.md)、[docs/ROADMAP.md](docs/ROADMAP.md)、[docs/DETAIL_INDEX.md](docs/DETAIL_INDEX.md)、対象 owner / collaborator 詳細本文、[docs/details/fixture.md](docs/details/fixture.md)、実在ファイルを確認する。

実在ファイルの確認には hidden fixture を含めて列挙できる `rg --files --hidden -g '!.git/**'` を使用する。

新規実装は、[docs/ROADMAP.md](docs/ROADMAP.md) の現在状態が `仕様化済み・未実装` であり、[docs/SPEC.md](docs/SPEC.md) の着手条件を満たす対象だけに行う。`実装中・検証未完了` は既着手範囲の継続、修正、検証だけを許可する。`未仕様化`、`将来計画`、`改訂予定` へ実装着手してはならない。

実装中に未定義の入力、出力、状態、異常系、セキュリティ条件、検証条件を発見した場合は、実装判断で補完せず、先に該当する責務正本を改訂する。

Go 実装の標準配置は [docs/SPEC.md 方針責務 §4.3](docs/SPEC.md#sec-4-3) を参照する。[`main.go`](main.go) は起動入口、[`components/builder.go`](components/builder.go)、[`components/runner.go`](components/runner.go)、[`components/api.go`](components/api.go) は各 owner component、[`admin/adlaire-ci-sdk.js`](admin/adlaire-ci-sdk.js) と [`admin/index.html`](admin/index.html) は管理クライアント実装として扱う。[`components/mcp.go` 将来追加予定 path](docs/ROADMAP.md) は [docs/ROADMAP.md](docs/ROADMAP.md) が将来計画の間は作成してはならない。

実装変更後は、変更範囲に応じて構文確認、単体確認、実行確認、生成物確認、異常系確認、必須 fixture 確認を行う。

Go 実装では、対象ファイルに `gofmt -l ...` を実行し、Go module が存在する場合は `go test ./...` を実行する。実行できない確認は、未実施理由を Pull Request 本文へ記録する。

API、SDK、UI のいずれかを変更する場合は、対応する endpoint、SDK method、UI 操作、状態副作用、認証・認可、成功後再取得、失敗時固定、fixture 証跡を同じ変更で確認する。

実装済みへの状態変更は、コード、仕様差分、構文確認、実行または生成物確認、必須 fixture、実装検証証跡、[docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md) の所在更新要否をすべて確認した後に限る。

---
## 4. Git 運用ルール

`main` は保護対象ブランチとする。

`main` への直接 push を禁止する。

変更作業は、必ず作業ブランチで行う。

作業ブランチは一本化し、ドキュメント変更作業と実装変更作業の両方で同じ作業ブランチを使用する。

同一目的、同一仕様領域、同一ファイル群に対する変更は、必ず 1 本の作業ブランチと 1 本の Pull Request にまとめる。

同一目的の変更を複数の積み上げ Pull Request に分割してはならない。

複数の Pull Request に分ける場合は、変更対象ファイル、責務、merge 順序が明確に分離でき、相互に同一ファイルを編集せず、片方だけが merge されても仕様矛盾、参照切れ、状態不一致、未定義の依存関係が発生しない場合に限る。

既存の open Pull Request と同じファイルまたは同じ仕様領域を変更する必要がある場合は、新規 Pull Request を作成せず、既存 Pull Request へ変更を統合する。

既存の open Pull Request と同じファイルまたは同じ仕様領域を変更する必要があるにもかかわらず、別 Pull Request を作成することを禁止する。

積み上げ Pull Request、同一ファイル編集の並行 Pull Request、merge 順序依存の Pull Request、または GitHub 上で `DIRTY` / conflict 状態の Pull Request が発生した場合は、競合解消作業として扱う。競合解消作業では、最新 `origin/main` から一本化ブランチを作成するか、最も包括的な既存 Pull Request の branch を統合先とし、必要な変更を 1 本の Pull Request に統合する。

一本化後、重複する既存 Pull Request は、統合先 Pull Request を明記したコメントを残して close する。

競合防止のため、作業開始前と Pull Request 作成前に以下を必ず実行する。

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

- `delete_branch_on_merge=true` は即時設定対象とする。
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

自動化に関わる設定には、少なくとも `delete_branch_on_merge=true` を含める。その他の自動化設定が [docs/SPEC.md](docs/SPEC.md) または本ルールブックで標準化された場合も同様に扱う。

`delete_branch_on_merge=true` は、remote branch 自動削除の必須設定とする。

エージェントは、ユーザー承認なしに GitHub リポジトリ設定を変更、無効化、初期化してはならない。

remote branch は、GitHub リポジトリ設定 `delete_branch_on_merge=true` により、Pull Request merge 後に GitHub 側で自動削除する。

remote branch 自動削除の対象は、merge 済み Pull Request の head branch に限定する。

local branch は、GitHub 側の自動削除では削除されない。

エージェントは、ユーザーによる Pull Request の merge 完了を確認できた場合、追加承認なしで対応する local branch の削除を自動実行する。

local branch 削除の対象は、merge 済み Pull Request の head branch と同名の local branch に限定する。

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

Pull Request作成自動化では、`main`への直接pushを行ってはならない。

Pull Request作成自動化では、Pull Requestのmergeを行ってはならない。mergeはユーザーが行う。

Pull Request作成自動化は、承認済み変更作業の範囲内で行うGit操作に限る。未承認のファイル作成、編集、移動、削除、リネーム、整形、生成物更新を含めてはならない。

Pull Request 作成前には、変更内容に応じて以下を確認する。

- `git fetch origin` を実行し、最新 `origin/main` を取得する。
- `git diff --name-status origin/main...HEAD` で、変更対象が承認済み範囲内であることを確認する。
- open Pull Request を確認し、同一ファイルまたは同一仕様領域を変更する Pull Request が存在しないことを確認する。
- 同一ファイルまたは同一仕様領域の open Pull Request が存在する場合は、新規 Pull Request ではなく既存 Pull Request への統合、または最新 `origin/main` 起点の一本化 Pull Request を作成する。
- open Pull Request の `mergeStateStatus` が `DIRTY` または `UNKNOWN` の場合は、競合状態または未確認状態として扱い、`CLEAN` を確認するまで報告しない。
- 競合解消または PR 一本化を行った場合は、重複 PR が open のまま残っていないことを `gh pr list --state open` で確認する。
- 文書変更では、`rg` で不要になった名称、矛盾参照、不要になったファイル名が残っていないか確認する。
- 文書変更では、`git diff --stat` で変更範囲を確認する。
- ファイル追加、削除、リネームを含む場合は、`git diff --cached --summary` で Git 上の扱いを確認する。
- 実装変更では、対象言語に応じた構文確認を行う。Go 実装では `gofmt -l ...` を標準の整形確認とし、Go module が存在する場合は `go test ./...` を標準の確認とする。
- 実装変更では、変更した実装が実行可能な場合は対象スクリプトの実行確認または生成物確認を行う。実行不能な場合は理由を Pull Request 本文に記録する。
- 仕様変更では、[docs/SPEC.md](docs/SPEC.md)、[docs/ROADMAP.md](docs/ROADMAP.md)、[docs/DETAIL_INDEX.md](docs/DETAIL_INDEX.md)、該当する owner component 別の [docs/details/*.md](docs/details/)、[docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md)、[docs/DESIGN.md](docs/DESIGN.md)、実装ファイルの整合を確認する。

Pull Request 本文には、少なくとも以下を記載する。

- `Summary`
- `Verification`
- 競合防止確認
- 未実施の確認がある場合は、その理由

---

## 5. 外部依存変更ルール

外部依存の採否、禁止条件、例外条件、許可範囲、許可外部ライブラリ一覧は、[docs/SPEC.md 方針責務 §4.1](docs/SPEC.md#sec-4-1) と [docs/SPEC.md ポリシー責務 §4](docs/SPEC.md#4-外部ライブラリフレームワーク方針) だけを正本とする。[AGENTS.md](AGENTS.md) で同じ方針または許可条件を再定義してはならない。

外部依存を追加、削除、更新、置換する前に、対象実装、[`go.mod`](go.mod)、配布物、セットアップ、検証手段への影響と、[docs/SPEC.md ポリシー責務 §4](docs/SPEC.md#4-外部ライブラリフレームワーク方針) の許可外部ライブラリ一覧を確認する。

許可外部ライブラリ一覧にない依存を実装へ追加してはならない。追加が必要な場合は、依存名、採用理由、代替困難性、責務範囲、影響範囲、保守・削除方針、検証条件、[docs/SPEC.md](docs/SPEC.md) の変更内容を提示し、実装変更と仕様変更の両方について事前承認を得る。

外部依存の追加、削除、更新、置換は変更作業として扱う。変更後は、[docs/SPEC.md ポリシー責務 §4](docs/SPEC.md#4-外部ライブラリフレームワーク方針) の許可一覧、[`go.mod`](go.mod)、実装 import、配布・セットアップ手順、検証結果が一致していることを確認する。

---

## 6. 文書整合ルール

[docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md) は、リポジトリ内の文書・実装ファイルの役割を示す索引として崩してはならない。

ファイル名、正本参照先、実装 artifact の追加・削除・リネームが発生した場合は、[docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md) の更新要否を確認する。

[docs/SPEC.md](docs/SPEC.md)、[docs/ROADMAP.md](docs/ROADMAP.md)、[docs/DETAIL_INDEX.md](docs/DETAIL_INDEX.md)、または owner component 別の [docs/details/*.md](docs/details/) を改訂した場合は、[docs/DESIGN.md](docs/DESIGN.md)、[docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md)、実装ファイルへの影響を確認する。

[docs/DESIGN.md](docs/DESIGN.md) を改訂した場合は、[`components/builder.go`](components/builder.go) 内の HTML / CSS / JavaScript / theme component テンプレートとの整合性を確認する。

仕様化済み項目を実装した場合は、[docs/ROADMAP.md](docs/ROADMAP.md) の現在状態、[docs/DOCUMENT_INDEX.md](docs/DOCUMENT_INDEX.md) の「実装ファイル一覧」、[docs/DETAIL_INDEX.md](docs/DETAIL_INDEX.md) の対応表、owner component 別の [docs/details/*.md](docs/details/) の検証条件、実装ファイルの存在を整合させる。
