# Build-Scripts - 最上位ルールブック

## 0. 絶対原則

`AGENTS.md` は、本リポジトリにおける最上位ルールブックである。

このリポジトリで作業するすべてのエージェントは、調査、設計、仕様改訂、実装、検証、Git 操作、Pull Request 作成、レビュー対応を含む全作業において、最初に本ファイルを必ず読む。

本ファイルの読了を完了するまで、調査、設計、仕様改訂、実装、検証、ファイル作成、編集、移動、削除、リネーム、整形、生成物更新、Git 操作、Pull Request 作成、レビュー対応を含む、リポジトリに関するすべての作業を開始してはならない。

`AGENTS.md` を確認しただけで、作業判断に必要な確認を完了したと扱ってはならない。

本リポジトリの仕様判断は、`build_spec_v3_spec.md` を正本として行う。

`DESIGN.md` は、`Adlaire-db-spec.html` のデザイン仕様を整理する補助文書である。機能仕様、運用仕様、API 仕様、CI 仕様の正本ではない。

`DOCUMENT_INDEX.md` は、文書・実装ファイルの役割を整理する索引である。仕様正本ではない。

`AGENTS.md` と他ファイルが作業ルール上矛盾する場合は、`AGENTS.md` を正とする。

`build_spec_v3_spec.md` と実装ファイルが仕様上矛盾する場合は、仕様と実装の不整合として扱う。仕様を変更する場合は、先に `build_spec_v3_spec.md` を改訂し、その内容に基づいて実装を更新する。

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

提示済み範囲を超える変更が必要になった場合は、追加の変更内容を提示し、別途 `承認` を得る。

---

## 2. 仕様書管理ルール

`build_spec_v3_spec.md` は、Adlaire CI のマスター仕様書正本である。

`build_spec_v3.py`、`runner.py`、将来コンポーネントである `api_server.py`、`adlaire-ci-sdk.js`、`admin/index.html`、`mcp_server.py` は、`build_spec_v3_spec.md` に基づいて更新する。

`DESIGN.md` は、出力 HTML のデザイン仕様を整理する補助文書である。`build_spec_v3_spec.md` と矛盾する場合は、`build_spec_v3_spec.md` を優先する。

仕様改訂では、既存仕様、`DOCUMENT_INDEX.md`、実装ファイルとの整合性を確認する。

`build_spec_v3_spec.md` に記載された一部コンポーネントや機能は、仕様化済みであっても未実装の場合がある。リポジトリ内に実装ファイルまたは実装コードが存在しない内容を、実装済み機能として扱ってはならない。

一時ファイル、退避ファイル、比較用ファイルは、整合性確認が完了するまで削除しない。

不要ファイル削除は、削除対象、削除理由、影響範囲を提示し、別途 `承認` を得てから行う。

---

## 3. 実装管理ルール

現行実装ファイルは以下とする。

| ファイル | 役割 |
|---------|------|
| `build_spec_v3.py` | Adlaire DB 仕様書 Markdown を単一 HTML へ変換するビルドスクリプト。 |
| `runner.py` | GitHub API で対象 Markdown の変更を検出し、ビルドパイプラインを実行する CI ランナー。 |

仕様化済みだが未実装の主なコンポーネントは以下とする。

| ファイル | 状態 |
|---------|------|
| `api_server.py` | 未実装 |
| `adlaire-ci-sdk.js` | 未実装 |
| `admin/index.html` | 未実装 |
| `mcp_server.py` | 将来計画 |

未実装コンポーネントを追加する場合は、`build_spec_v3_spec.md` の該当仕様、`DOCUMENT_INDEX.md`、本ファイルを必要に応じて整合させる。

実装変更後は、変更範囲に応じて構文確認、実行確認、生成物確認を行う。

Python 実装の構文確認では、環境に応じて `PYTHONPYCACHEPREFIX=/tmp/codex-pycache python3 -m py_compile ...` を使用してよい。

---

## 4. Git 運用ルール

`main` は保護対象ブランチとする。

`main` への直接 push を禁止する。

変更作業は、必ず作業ブランチで行う。

作業ブランチは一本化し、ドキュメント変更作業と実装変更作業の両方で同じ作業ブランチを使用する。

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
- main branch protection: 設定対象（下記の初期標準を適用）

GitHub 設定の初期適用方針は以下とする。

- `delete_branch_on_merge=true` は即時設定対象とする。
- `main` branch protection は、初期標準として Pull Request 必須、force push 禁止、branch deletion 禁止を設定する。
- `main` branch protection の required approvals は初期値 `0` とする。
- 運用が安定した後、必要に応じて required approvals を `1` へ引き上げる。

標準 GitHub リポジトリ設定のうち、自動化に関わる設定が未確認の場合は、現在の設定状態を確認する。

自動化に関わる設定が未設定または標準値と異なる場合は、変更対象、変更内容、影響範囲を提示し、ユーザーから `承認` を得たうえで標準値へ設定する。

自動化に関わる設定には、少なくとも `delete_branch_on_merge=true` を含める。その他の自動化設定が `build_spec_v3_spec.md` または本ルールブックで標準化された場合も同様に扱う。

`delete_branch_on_merge=true` は、remote branch 自動削除の必須設定とする。

エージェントは、ユーザー承認なしに GitHub リポジトリ設定を変更、無効化、初期化してはならない。

remote branch は、GitHub リポジトリ設定 `delete_branch_on_merge=true` により、Pull Request merge 後に GitHub 側で自動削除する。

remote branch 自動削除の対象は、merge 済み Pull Request の head branch に限定する。

local branch は、GitHub 側の自動削除では削除されない。

エージェントは、ユーザーによる Pull Request の merge 完了を確認できた場合、追加承認なしで対応する local branch の削除を自動実行する。

local branch 削除の対象は、merge 済み Pull Request の head branch と同名の local branch に限定する。

local branch 削除では、`main` へ移動した後に対象 local branch を削除する。

`main`、merge 未完了の作業ブランチ、merge 状態を確認できないブランチ、Pull Request と対応しないブランチは削除してはならない。

`.gitignore` は作成・使用しない。

`.gitignore` が必要になる生成物、一時ファイル、実行時データ、ビルド成果物が発生した場合は、除外設定で隠蔽せず、生成先、運用、または実装を見直す。

承認済み変更作業が完了した場合、エージェントはユーザーからの追加指示および追加承認なしで、作業ブランチでのcommit、remoteへのpush、Pull Requestの作成または既存Pull Requestの更新まで自動実行する。

Pull Request作成自動化では、`main`への直接pushを行ってはならない。

Pull Request作成自動化では、Pull Requestのmergeを行ってはならない。mergeはユーザーが行う。

Pull Request作成自動化は、承認済み変更作業の範囲内で行うGit操作に限る。未承認のファイル作成、編集、移動、削除、リネーム、整形、生成物更新を含めてはならない。

---

## 5. 外部依存ルール

外部フレームワークおよび外部ライブラリは、原則として採用しない。

機能実現は、Python 標準ライブラリ、Vanilla JavaScript、内製実装、例外承認済み外部ライブラリの順で検討する。

外部依存は最小限に抑え、可能な範囲で内製化を重視する。

ただし、開発コスト、実装難易度、安全性、保守性、暗号・認証等の専門性を考慮し、外部ライブラリの採用を例外として許可する場合がある。

例外として外部ライブラリを採用する場合は、採用理由、対象範囲、影響範囲、代替困難性、保守方針を明示する。

例外採用は、`build_spec_v3_spec.md` または承認済み変更範囲に明記された場合のみ有効とする。

外部依存を追加、削除、更新、置換する作業は変更作業として扱い、事前承認を必須とする。

---

## 6. 文書整合ルール

`DOCUMENT_INDEX.md` は、リポジトリ内の文書・実装ファイルの役割を示す索引として維持する。

ファイル名、正本関係、実装コンポーネントの追加・削除・リネームが発生した場合は、`DOCUMENT_INDEX.md` の更新要否を確認する。

`build_spec_v3_spec.md` を改訂した場合は、`DESIGN.md`、`DOCUMENT_INDEX.md`、実装ファイルへの影響を確認する。

`DESIGN.md` を改訂した場合は、`build_spec_v3.py` 内の HTML / CSS / JavaScript テンプレートとの整合性を確認する。

仕様化済み項目を実装した場合は、`build_spec_v3_spec.md` 内の状態表現、`DOCUMENT_INDEX.md` の Planned Components、実装ファイルの存在を整合させる。
