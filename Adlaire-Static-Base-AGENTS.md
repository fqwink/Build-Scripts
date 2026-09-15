# Adlaire-Static-Base - 最上位ルールブック

## 0. 絶対原則

`AGENTS.md` は、本リポジトリにおける最上位ルールブックである。

このリポジトリで作業するすべてのエージェントは、調査、設計、仕様改訂、実装、検証、Git 操作、Pull Request 作成、レビュー対応を含む全作業において、最初に本ファイルを必ず読む。

本ファイルの読了を完了するまで、調査、設計、仕様改訂、実装、検証、ファイル作成、編集、移動、削除、リネーム、整形、生成物更新、Git 操作、Pull Request 作成、レビュー対応を含む、リポジトリに関するすべての作業を開始してはならない。

`AGENTS.md` を確認しただけで、作業判断に必要な確認を完了したと扱ってはならない。

本リポジトリの仕様判断は、`ASB-spec.md` を正本として行う。

`ASB-spec.html` は、`ASB-spec.md` に基づいて更新・再生成するHTML版仕様書である。

`AGENTS.md` と他ファイルが作業ルール上矛盾する場合は、`AGENTS.md` を正とする。

`ASB-spec.md` と `ASB-spec.html` が仕様上矛盾する場合は、`ASB-spec.md` を正とする。

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

`ASB-spec.md` は、本リポジトリのマスター仕様書正本である。

`ASB-spec.html` は、`ASB-spec.md` に基づいて更新・再生成するHTML版仕様書である。

`ASB-spec.md` を変更した場合は、`ASB-spec.html` の更新・再生成要否を確認する。

`ASB-spec.md` と `ASB-spec.html` が仕様上矛盾する場合は、`ASB-spec.md` を正とする。

仕様改訂では、既存仕様との整合性を確認する。

一時ファイル、退避ファイル、比較用ファイルは、整合性確認が完了するまで削除しない。

不要ファイル削除は、削除対象、削除理由、影響範囲を提示し、別途 `承認` を得てから行う。

---

## 3. Git 運用ルール

`main` は保護対象ブランチとする。

`main` への直接 push を禁止する。

変更作業は、必ず作業ブランチで行う。

作業ブランチは一本化し、ドキュメント変更作業と実装変更作業の両方で同じ作業ブランチを使用する。

`main` への反映は、Pull Request 経由で行う。

Pull Request の merge はユーザーが行う。

エージェントは Pull Request の merge を行ってはならない。

GitHub リポジトリ設定は、ASB の Git 運用前提として管理する。

GitHub リポジトリ設定の変更は変更作業として扱い、事前に変更対象、変更内容、影響範囲を提示し、ユーザーから `承認` を得るまで実行してはならない。

ASB の標準 GitHub リポジトリ設定は以下とする。

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
- main branch protection: 設定対象

`delete_branch_on_merge=true` は、remote branch 自動削除の必須設定とする。

エージェントは、ユーザー承認なしに GitHub リポジトリ設定を変更、無効化、初期化してはならない。

remote branch は、GitHub リポジトリ設定 `delete_branch_on_merge=true` により、Pull Request merge 後に GitHub 側で自動削除する。

remote branch 自動削除の対象は、merge 済み Pull Request の head branch に限定する。

local branch は、GitHub 側の自動削除では削除されない。

エージェントは、ユーザーによる Pull Request の merge 完了を確認できた場合、追加承認なしで対応する local branch の削除を自動実行する。

local branch 削除の対象は、merge 済み Pull Request の head branch と同名の local branch に限定する。

local branch 削除では、`main` へ移動した後に対象 local branch を削除する。

`main`、merge 未完了の作業ブランチ、merge 状態を確認できないブランチ、Pull Request と対応しないブランチは削除してはならない。

SSH URL は `origin` に設定し、HTTPS URL はバックアップ remote として保持する。

`.gitignore` は作成・使用しない。

`.gitignore` が必要になる生成物、一時ファイル、実行時データ、ビルド成果物が発生した場合は、除外設定で隠蔽せず、生成先、運用、または実装を見直す。

承認済み変更作業が完了した場合、エージェントはユーザーからの追加指示および追加承認なしで、作業ブランチでのcommit、remoteへのpush、Pull Requestの作成または既存Pull Requestの更新まで自動実行する。

Pull Request作成自動化では、`main`への直接pushを行ってはならない。

Pull Request作成自動化では、Pull Requestのmergeを行ってはならない。mergeはユーザーが行う。

Pull Request作成自動化は、承認済み変更作業の範囲内で行うGit操作に限る。未承認のファイル作成、編集、移動、削除、リネーム、整形、生成物更新を含めてはならない。

---

## 4. 外部依存ルール

外部フレームワークおよび外部ライブラリは、原則として採用しない。

機能実現は、Go標準ライブラリ、内製実装、例外承認済み外部ライブラリの順で検討する。

外部依存は最小限に抑え、可能な範囲で内製化を重視する。

ただし、開発コスト、実装難易度、安全性、保守性、暗号・認証等の専門性を考慮し、外部ライブラリの採用を例外として許可する場合がある。

例外として外部ライブラリを採用する場合は、採用理由、対象範囲、影響範囲、代替困難性、保守方針を明示する。

例外採用は、`ASB-spec.md` または承認済み変更範囲に明記された場合のみ有効とする。

外部依存を追加、削除、更新、置換する作業は変更作業として扱い、事前承認を必須とする。

---

## 5. テンプレート運用ルール

`templates/` 配下のファイルは、新規文書、仕様節、実装フェーズ、変更履歴、ルールブックを作成または改訂する際の雛形である。

テンプレートは、現行リポジトリの仕様正本、作業ルール、実装タスク、Git運用ルールを直接変更するものではない。

`templates/rulebook/AGENTS.md` は、新規リポジトリまたは派生リポジトリ向けの最上位ルールブック雛形であり、本リポジトリの作業判断には使用しない。

本リポジトリの作業判断では、常に現行の `AGENTS.md` を正とする。

`templates/rulebook/` 配下のテンプレートと現行 `AGENTS.md` が矛盾する場合は、現行 `AGENTS.md` を正とする。

テンプレートを変更しただけでは、本リポジトリの作業ルールは変更されない。

本リポジトリの作業ルールを変更する場合は、必ず現行 `AGENTS.md` 本体を変更する。

テンプレートを理由に、仕様、API、設定項目、保存JSON、ディレクトリ、外部依存、実行時データ、Git運用、承認ルールを追加または変更してはならない。

テンプレートから新規文書を作成する場合も、ファイル作成、編集、移動、削除、リネーム、整形、生成物更新として扱い、事前承認を必須とする。

テンプレートから作成した一時ファイル、比較用ファイル、生成途中ファイルを、承認なしに開発リポジトリ内へ残してはならない。

テンプレート更新時は、`ASB-spec.md`、`DOCUMENT_INDEX.md`、`README.md`、`IMPLEMENTATION_TASKS.md` との整合を確認する。

---

## 6. 仕様改訂後の整合自動化ルール

`ASB-spec.md` は、本リポジトリの仕様正本である。

仕様の改訂は、必ず `ASB-spec.md` のみで行う。

`ASB-spec.html` および `IMPLEMENTATION_TASKS.md` を、仕様改訂の入力元として扱ってはならない。

`ASB-spec.md` を改訂した場合は、改訂後の `ASB-spec.md` に記載された本書バージョンを確認する。

`ASB-spec.md` の本書バージョンを確認した後、`ASB-spec.html` を `ASB-spec.md` に基づいて更新・再生成する。

`ASB-spec.md` の本書バージョンを確認した後、`IMPLEMENTATION_TASKS.md` を `ASB-spec.md` に基づいて更新する。

`ASB-spec.html` および `IMPLEMENTATION_TASKS.md` の更新は、`ASB-spec.md` 改訂後の整合作業として扱う。

この整合作業は、`ASB-spec.md` 改訂後に省略してはならない。

ただし、`ASB-spec.md`、`ASB-spec.html`、`IMPLEMENTATION_TASKS.md` の作成、編集、再生成、更新、Git 操作は、すべて本ルールブックの承認ルールに従う。
