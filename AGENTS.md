# Build-Scripts - 作業ルールブック

## 0. 最上位文書

`AGENTS.md` と `docs/SPEC.md` は、本リポジトリにおける最上位文書である。

`AGENTS.md` は、作業ルール、承認、Git 操作、Pull Request 作成、レビュー対応、検証手順、エージェント実行手順の最上位ルールブックである。

`docs/SPEC.md` は、仕様、方針、ポリシー、正本関係、禁止事項、リリース判断、実装着手可否の最上位仕様書である。

このリポジトリで作業するすべてのエージェントは、調査、設計、仕様改訂、実装、検証、Git 操作、Pull Request 作成、レビュー対応を含む全作業において、最初に `AGENTS.md` と `docs/SPEC.md` の両方を必ず読む。

作業ルール上の矛盾は `AGENTS.md` を正とする。仕様、方針、ポリシー、正本関係、禁止事項、リリース判断、実装着手可否の矛盾は `docs/SPEC.md` を正とする。

## 1. 承認ルール

変更作業では、承認工程を省略してはならない。

ファイル作成、編集、移動、削除、リネーム、整形、生成物更新、Git 操作、GitHub 設定変更など、リポジトリまたはリモート状態を変更する作業は、ユーザーから事前承認を得るまで実行してはならない。

変更作業前には、必ず変更対象、変更内容、影響範囲を提示する。

ユーザー承認は、ユーザーの返信に `承認` という単語が明示された場合のみ有効とする。

`OK`、`はい`、`お願いします`、`進めて`、その他の類似表現は、変更作業の承認として扱わない。

承認後は、提示済みの変更作業内容の範囲内でのみ作業する。

作業中に新たな不整合、改善候補、設定差分を発見した場合でも、承認済み範囲外であれば編集、移動、削除、リネーム、生成物更新、設定変更を行ってはならない。

提示済み範囲を超える変更が必要になった場合は、追加変更対象、追加変更内容、影響範囲、今その変更を行う必要性を提示し、別途 `承認` を得る。

読取、検索、差分確認、検証など、リポジトリ状態を変更しない調査は、承認済み作業の判断材料として実行してよい。

## 2. 作業開始ルール

作業開始時には、対象機能・対象コンポーネントについて以下を確認する。

- `docs/SPEC.md` の方針、ポリシー、禁止事項
- `docs/ROADMAP.md` の状態、実装可否、Phase、将来計画
- `docs/DETAIL_INDEX.md` の詳細仕様入口、読み順、対応表
- 該当する owner component 別の `docs/details/*.md`
- `docs/DOCUMENT_INDEX.md` の文書・実装ファイル所在
- 実ファイルの存在

実ファイルの存在確認には `rg --files` を使用する。

仕様に関する判断は、`docs/SPEC.md` の責務分離と正本関係に従う。

`docs/SPEC.md`、`docs/ROADMAP.md`、`docs/DETAIL_INDEX.md`、owner component 別の `docs/details/*.md` と実装ファイルが矛盾する場合は、仕様と実装の不整合として扱う。仕様を変更する場合は、先に該当する仕様書を改訂し、その内容に基づいて実装を更新する。

## 3. 実装管理ルール

標準ソース配置は `main.go`、`components/*.go`、`admin/` 配下の管理 UI ファイル、`testdata/<component>/` とする。

Go 実装対象ファイルは以下とする。

| ファイル | 役割 |
|---------|------|
| `main.go` | 起動入口。実行ファイル名に応じて対象 component を呼び出す。 |
| `components/builder.go` | Markdown ファイルまたは Markdown ディレクトリを静的 Web サイトへ変換する Go 版ビルドスクリプト。 |
| `components/runner.go` | GitHub API で対象 Markdown の変更を検出し、ビルドパイプラインを実行する Go 版 CI ランナー。 |
| `components/api.go` | 管理 API サーバー。 |
| `components/mcp.go` | 将来計画コンポーネント。 |

管理 UI は `admin/index.html`、`admin/adlaire-ci-sdk.js`、`admin/style.css`、`admin/app.js` を標準配置とする。

`build_spec.go`、`runner.go`、`build_spec_test.go`、`runner_test.go`、`testdata/build_spec/` を現行実装実体として扱ってはならない。

実装変更後は、変更範囲に応じて構文確認、実行確認、生成物確認を行う。

Go 実装の構文確認では、対象ファイルに対して `gofmt -l ...` を実行し、Go module が存在する場合は `go test ./...` を実行する。

## 4. Git 運用ルール

`main` は保護対象ブランチとする。

`main` への直接 push を禁止する。

変更作業は、必ず作業ブランチで行う。

同一目的、同一仕様領域、同一ファイル群に対する変更は、必ず 1 本の作業ブランチと 1 本の Pull Request にまとめる。

同一目的の変更を複数の積み上げ Pull Request に分割してはならない。

既存の open Pull Request と同じファイルまたは同じ仕様領域を変更する必要がある場合は、新規 Pull Request を作成せず、既存 Pull Request へ変更を統合する。

競合防止のため、作業開始前と Pull Request 作成前に以下を必ず実行する。

1. `git fetch origin`
2. `gh pr list --state open --json number,title,headRefName,baseRefName,mergeStateStatus,url`
3. `git diff --name-status origin/main...HEAD`

Pull Request の `mergeStateStatus` が `DIRTY`、`UNKNOWN`、または確認不能の場合は、merge 可能と報告してはならない。`UNKNOWN` の場合は GitHub の再計算後に再確認し、最終的に `CLEAN` を確認する。

競合解消時は、競合マーカーの除去だけで完了としてはならない。`git diff --check`、競合マーカー検索、変更対象文書の正本関係確認、open Pull Request 一覧確認を完了条件とする。

`main` への反映は Pull Request 経由で行う。

Pull Request の merge はユーザーが行う。エージェントは Pull Request の merge を行ってはならない。

承認済み変更作業が完了した場合、エージェントはユーザーからの追加指示および追加承認なしで、作業ブランチでの commit、remote への push、Pull Request の作成または既存 Pull Request の更新まで自動実行する。

Pull Request 本文には、少なくとも `Summary`、`Verification`、競合防止確認、未実施の確認がある場合はその理由を記載する。

## 5. GitHub 設定ルール

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
- main branch protection: Pull Request 必須、force push 禁止、branch deletion 禁止

GitHub 設定を変更した後は、GitHub API で再取得し、標準設定との差分がないかを確認して報告する。

remote branch は、GitHub リポジトリ設定 `delete_branch_on_merge=true` により、Pull Request merge 後に GitHub 側で自動削除する。

local branch は、GitHub 側の自動削除では削除されない。ユーザーによる Pull Request の merge 完了を確認できた場合、追加承認なしで対応する local branch を削除してよい。

## 6. 文書整合ルール

`docs/DOCUMENT_INDEX.md` は、リポジトリ内の文書・実装ファイルの役割を示す索引として維持する。

ファイル名、正本関係、実装コンポーネントの追加・削除・リネームが発生した場合は、`docs/DOCUMENT_INDEX.md` の更新要否を確認する。

`docs/SPEC.md`、`docs/ROADMAP.md`、`docs/DETAIL_INDEX.md`、または owner component 別の `docs/details/*.md` を改訂した場合は、`docs/DESIGN.md`、`docs/DOCUMENT_INDEX.md`、実装ファイルへの影響を確認する。

`.gitignore` は作成・使用しない。

不要ファイル削除は、削除対象、削除理由、影響範囲を提示し、別途 `承認` を得てから行う。
