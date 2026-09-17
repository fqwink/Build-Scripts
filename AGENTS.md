# Build-Scripts - 最上位ルールブック

## 0. 絶対原則

`AGENTS.md` と `docs/SPEC.md` は、本リポジトリにおける最上位文書である。

`AGENTS.md` は、作業ルール、承認、Git 操作、Pull Request 作成、レビュー対応、検証手順、エージェント実行手順の最上位ルールブックである。

`docs/SPEC.md` は、仕様、方針、ポリシー、正本関係、禁止事項、リリース判断、実装着手可否の最上位仕様書である。

このリポジトリで作業するすべてのエージェントは、調査、設計、仕様改訂、実装、検証、Git 操作、Pull Request 作成、レビュー対応を含む全作業において、最初に `AGENTS.md` と `docs/SPEC.md` の両方を必ず読む。

`AGENTS.md` と `docs/SPEC.md` の読了を完了するまで、調査、設計、仕様改訂、実装、検証、ファイル作成、編集、移動、削除、リネーム、整形、生成物更新、Git 操作、Pull Request 作成、レビュー対応を含む、リポジトリに関するすべての作業を開始してはならない。

`AGENTS.md` または `docs/SPEC.md` の片方だけを確認した状態で、作業判断に必要な確認を完了したと扱ってはならない。

本リポジトリの仕様判断は、方針、ポリシー、正本関係、禁止事項、リリース判断は `docs/SPEC.md`、実装状態、実装可否、Phase、機能インベントリ、将来計画、昇格手順、§27 追加仕様化機能参照、横断補足契約は `docs/ROADMAP.md`、詳細仕様の入口、索引、共通固定値、実装前確認項目、検証マトリクス、詳細節対応表、リポジトリ内ソース配置、責務 component 別詳細仕様ファイル管理仕様は `docs/DETAIL_INDEX.md`、各 owner component の詳細本文は `docs/details/*.md` を正本として行う。

`docs/DESIGN.md` は、生成静的 Web サイトのデザイン仕様を整理する補助文書である。機能仕様、運用仕様、API 仕様、CI 仕様の正本ではない。

`docs/DOCUMENT_INDEX.md` は、文書・実装ファイルの役割を整理する索引である。仕様正本ではない。

`AGENTS.md` と他ファイルが作業ルール上矛盾する場合は、`AGENTS.md` を正とする。

`docs/SPEC.md` と他ファイルが仕様、方針、ポリシー、正本関係、禁止事項、リリース判断、実装着手可否で矛盾する場合は、`docs/SPEC.md` を正とする。

`docs/SPEC.md`、`docs/ROADMAP.md`、`docs/DETAIL_INDEX.md`、または該当する owner component 別の `docs/details/*.md` と実装ファイルが仕様上矛盾する場合は、仕様と実装の不整合として扱う。仕様を変更する場合は、先に該当する仕様書を改訂し、その内容に基づいて実装を更新する。

---

## 0a. SPEC 方針・ポリシー統合ルール

本ルールブックは、作業者が `docs/SPEC.md` の方針・ポリシー正本を読み落とさないよう、`docs/SPEC.md` の方針責務・ポリシー責務の要点を本節へ統合して扱う。

`docs/SPEC.md` は、Adlaire CI の方針責務・ポリシー責務の正本である。

`docs/SPEC.md` は、目的、設計方針、禁止事項、成熟度、実装着手可否、仕様判断の原則を扱う。

方針またはポリシーに該当する内容は、必ず `docs/SPEC.md` に記載する。方針またはポリシーに該当する内容を、`docs/DETAIL_INDEX.md`、owner component 別の `docs/details/*.md`、`docs/ROADMAP.md`、`docs/DOCUMENT_INDEX.md`、`README.md`、実装ファイル、fixture、PR 本文へ移してはならない。

`docs/SPEC.md` の `技術方針` 表と `リポジトリ構成方針` のディレクトリ構成 tree は、Adlaire CI の明示的方針であり、削除してはならない。これらは実装詳細、状態一覧、索引本文として扱わない。

`docs/SPEC.md` は、関数単位の処理、HTTP response schema、状態ファイル schema、SDK method、UI DOM、fixture assertion、具体的な実行手順を定義しない。これらの本文は、責務ベース明示的原則、詳細仕様入口責務、owner component 詳細本文責務に従って、責務を持つ正本へ分離する。

状態責務は `docs/ROADMAP.md` を正本とする。

`docs/SPEC.md` では、状態分類、Phase 実装単位、実装着手可否の方針とポリシーだけを定義する。個別 component の現在状態、Phase 一覧、将来計画一覧、昇格手順は重複定義しない。

責務文書構成は以下とする。

| 責務 | 正本 | 責務の問い | 記載する内容 |
|------|------|-----------|------------|
| 方針責務 | `docs/SPEC.md` | **なぜ・何を** | 目的、設計思想、方向性の原則。変更頻度が低く、判断の拠り所となる指針。 |
| ポリシー責務 | `docs/SPEC.md` | **しなければならない／してはならない** | 遵守義務のある規則、制約、禁止事項、セキュリティ要件、運用ルール、バージョン管理規則。 |
| 状態・計画責務 | `docs/ROADMAP.md` | **いつ・どれを** | 実装状態、実装可否、Phase、機能インベントリ、将来計画、昇格手順。 |
| 詳細仕様入口責務 | `docs/DETAIL_INDEX.md` | **どこから読むか** | 詳細仕様の入口、読み順、共通固定値、詳細節対応表、リポジトリ内ソース配置。 |
| owner component 詳細本文責務 | `docs/details/*.md` | **どのように実装するか** | owner component 別の入出力、状態、処理順序、異常系、検証条件。 |
| 文書・実装ファイル所在の索引責務 | `docs/DOCUMENT_INDEX.md` | **どこにあるか** | 文書、実装ファイル、生成物、将来追加予定 path の所在。 |
| 利用入口責務 | `README.md` | **どう始めるか** | 利用者向け入口、概要、参照先。 |

新しい記載内容は「この内容はどの責務の問いに答えるか」を基準に、責務を持つ正本を決定する。

判断対象が方針またはポリシーである場合、責務を持つ正本は必ず `docs/SPEC.md` とする。判断対象が実装詳細、実装状態、文書所在、利用入口、検証証跡のいずれかである場合だけ、該当責務の正本を参照する。

作業者は、実装前に以下の参照順序方針に従う。

1. `docs/SPEC.md` の文書責務と状態参照方針で、正本範囲を確認する。
2. `docs/ROADMAP.md` で、対象の状態、実装可否、Phase、将来計画該当有無を確認する。
3. `docs/SPEC.md` 方針責務で、ゼロ依存、責務ベース明示的原則、リポジトリ構成、詳細仕様粒度、成熟度、着手ゲート、完了判定、Go 正本方針を確認する。
4. `docs/SPEC.md` ポリシー責務で、対象領域の禁止事項、セキュリティ、バージョン、外部依存を確認する。
5. `docs/DETAIL_INDEX.md` 詳細仕様入口責務で読み順、共通固定値、詳細節対応表、リポジトリ内ソース配置を確認し、該当する owner component 別の `docs/details/*.md` 詳細本文責務で実装に必要な入出力、状態、異常系、検証条件、配置を確認する。

## 0b. 技術方針・構成方針

Adlaire CI の技術選定は以下とする。本表は明示的方針であり、削除または他文書への移動を禁止する。

| 領域 | 方針 |
|---|---|
| ランタイム | Go stable release |
| 言語 | Go（ビルド、ランナー、管理 API）/ JavaScript（SDK）/ HTML・CSS・Vanilla JavaScript（UI） |
| HTTP | Go 標準ライブラリ `net/http` |
| データベース | なし（ファイルベース） |
| Git 操作 | GitHub REST API（Blobs API）を Go 標準ライブラリ `net/http` 経由で呼び出す |
| フロントエンド | HTML / CSS / Vanilla JavaScript |
| 標準運用 | CI サーバー（VPS 等）でビルドし、静的コンテンツ配信サーバーへ SSH で転送する 2 サーバー構成 |
| データ交換形式 | JSON に統一する。エクスポート・インポートを含む全 API データ交換に CSV・XML 等の非 JSON 形式を使用しない |

Adlaire CI は、ゼロ依存・フルインハウスを技術哲学の中核とする。

本原則におけるゼロ依存とは、各コンポーネントが外部ライブラリ、外部フレームワーク、外部ビルドツール、外部ホスティング実行基盤に機能成立を依存しないことを意味する。本原則におけるフルインハウスとは、Markdown 変換、CI 実行、管理 API、SDK、標準管理ツール、状態管理、認証、ログ、通知、セットアップの主要機能を本リポジトリ内で仕様化し、内製コードとして理解、検証、保守できる状態を意味する。

実装者は、便利さ、実装速度、一般的な慣習を理由に外部ライブラリで未定義機能を補完してはならない。外部依存がなければ成立しない設計は設計不備として扱い、先に仕様を見直す。

Adlaire CI は、`core`、`adlaire-ci-core`、`internal/core`、`common`、`base`、`foundation`、`utils` のような中心化・汎用置き場化する概念を採用しない。

横断的に利用される処理は、単一責務、入出力、状態、異常系、検証条件、依存方向が仕様化されている場合に限り、共通責務コンポーネントとして定義できる。ただし、共通責務コンポーネントは他コンポーネントより上位ではなく、同格の独立コンポーネントとして扱う。共通責務コンポーネントを、正本、中心、基盤、親、上位レイヤーとして扱ってはならない。

共通責務コンポーネントは、コード共有のためだけに追加してはならない。重複削減より、責務境界、仕様の明確さ、依存方向の追跡可能性を優先する。

Adlaire CI のリポジトリ内ソース構成は、責務ベースで整理する。本節のディレクトリ構成 tree は、Adlaire CI の構成方針を明示する正本であり、削除または他文書への移動を禁止する。

標準構成は以下とする。

```text
.
├── main.go
│
├── components/
│   ├── builder.go
│   ├── runner.go
│   ├── api.go
│   └── mcp.go
│
├── admin/
│   ├── index.html
│   ├── adlaire-ci-sdk.js
│   ├── style.css
│   └── app.js
│
├── testdata/
│   ├── builder/
│   ├── runner/
│   ├── api/
│   ├── statefile/
│   ├── archive/
│   ├── commitstatus/
│   └── fixtures/
│
├── docs/
│   ├── SPEC.md
│   ├── ROADMAP.md
│   ├── DETAIL_INDEX.md
│   ├── DOCUMENT_INDEX.md
│   ├── DESIGN.md
│   ├── examples/
│   └── details/
│       ├── builder.md
│       ├── runner.md
│       ├── api.md
│       ├── sdk.md
│       ├── ui.md
│       ├── setup.md
│       ├── statefile.md
│       ├── archive.md
│       ├── fixture.md
│       └── commitstatus.md
│
├── README.md
├── AGENTS.md
└── go.mod
```

`main.go` は 1 ファイルとし、起動入口、実行ファイル名判定、引数受け取り、対象 owner component 呼び出しだけを担当する。`main.go` に Markdown 変換、CI 実行、HTTP handler、状態ファイル操作、archive 処理、GitHub Commit Status 送信、MCP 処理の実装詳細を書いてはならない。

`components/` は、1 実装対象コンポーネント = 1 Go ファイルとする。`admin` は Go コンポーネントではなく `admin/` 配下の静的配布物として扱う。標準構成に含まれることだけを理由に、未実装ファイル、将来追加予定 path、空ディレクトリ、placeholder を作成してはならない。

## 0c. 責務ベース明示的原則

Adlaire CI の仕様体系は、責務ベース明示的原則を仕様全般の最上位方針として採用する。

責務ベース明示的原則とは、方針、ポリシー、実装状態、実装可否、Phase、将来計画、詳細仕様本文、fixture、expected、fake、検証証跡、文書索引を、それぞれ異なる責務として明示的に分離し、同一判断対象を複数文書で正本化しない原則である。

正本関係は、必ず責務名とファイル名で示す。文書を章構成、便宜分類、または他文書の従属章として扱ってはならない。`docs/SPEC.md` 内部の見出しも責務名で示し、`Part` 名称で扱ってはならない。`docs/DETAIL_INDEX.md`、`docs/details/*.md`、`docs/ROADMAP.md`、`docs/DOCUMENT_INDEX.md`、`README.md` を `docs/SPEC.md` の章として扱ってはならない。

方針とポリシーは `docs/SPEC.md` だけに記載する。詳細仕様、ロードマップ、索引、README、実装ファイル、fixture、PR 本文は、方針またはポリシーを本文として定義、補足、緩和、例外化、再解釈してはならない。

方針またはポリシーに該当する記載を `docs/SPEC.md` から削除し、他文書への参照だけに置き換えてはならない。方針またはポリシーを整理する場合は、`docs/SPEC.md` 内で責務名、適用範囲、禁止事項、参照先を明示して整える。

詳細仕様は実装詳細だけを記載する。ロードマップは状態、実装可否、Phase、将来計画だけを記載する。索引は文書と実装ファイルの所在だけを記載する。README は利用入口だけを記載する。fixture、expected、fake、検証証跡、完了判定は検証責務だけを記載する。

各文書は、自分の責務外の内容を本文として記載してはならない。責務外の内容が必要な場合は、責務を持つ正本と節番号を参照する。参照先に必要な内容が存在しない場合は、参照元へ補足を書かず、責務を持つ正本側を改訂する。責務を持つ正本側で確定できない内容は未確定として扱い、仕様化済みとして扱ってはならない。

分かりやすさは、同じ説明を複数文書へ重複記載することで担保してはならない。分かりやすさは、責務名、責務を持つ正本、正本範囲、参照節を明示することで担保する。責務名を使わない正本説明、判断に迷う記載、二重に読める記載、例外に見える記載、参照先を持たない責務外説明は禁止する。

本原則への違反が残る状態では、仕様整合完了、仕様 PR 完了、実装着手、実装済み判定を行ってはならない。

## 0d. Phase・着手・完了判定方針

仕様項目は、実装可否を判断できるように成熟度を明確に区分する。

成熟度は、構想の有無ではなく、実装者が実装に着手できるだけの情報が揃っているかで判定する。

将来計画や未確定アイデアは、実装対象として扱わない。実装対象にする場合は、先に詳細仕様を整え、仕様化済み・未実装へ昇格させる。

Adlaire CI の実装順序、実装計画、実装 PR、完了判定は Phase 単位でのみ管理する。

Phase は、対象 owner component、実装範囲、依存条件、完了条件、検証条件が明確な実装単位である。Phase の一覧、順序、依存条件、完了条件は `docs/ROADMAP.md` 状態・計画責務を正本とする。

`P0`、`P1`、`P2〜P5` などの優先度ラベル、抽象段階、API 内部分類、fixture 分類を、実装単位、PR 単位、完了判定単位として扱ってはならない。

実装者は、対象コンポーネントについて以下をすべて満たすまで実装に着手してはならない。

- `docs/SPEC.md` 方針責務・ポリシー責務、`docs/ROADMAP.md` 状態・計画責務、`docs/DETAIL_INDEX.md` 詳細仕様入口責務、owner component 詳細本文責務の該当節を確認済みである。
- 対象が `仕様化済み・未実装`、`実装中・検証未完了`、または `実装済み` に分類され、`未仕様化`、`将来計画`、`改訂予定` ではない。
- 対象機能が `docs/DETAIL_INDEX.md` 詳細仕様入口責務の詳細節対応表に含まれ、該当する詳細仕様節と受け入れ条件を確認済みである。
- 対象機能の詳細仕様が、機能仕様テンプレートに必要な項目を満たしている。
- 対象コンポーネント、呼び出し元、呼び出し先、状態ファイル、外部接続先が明確である。
- API、SDK、UI、状態ファイル、セットアップ、検証条件のうち関係する仕様が同時に整合している。
- 外部依存、秘密情報、直接 API 呼び出し、互換処理、暗黙フォールバックなどの禁止事項が明確である。
- `docs/ROADMAP.md`、`docs/DETAIL_INDEX.md`、owner component 詳細本文責務の受け入れ条件で、実装完了、実装順序、配置を判定できる。

着手条件を満たさない場合、実装者はコードで補完せず、先に仕様改訂を行う。

実装完了は、実装ファイルの作成やテスト成功だけでは成立しない。以下をすべて満たした場合にのみ `実装済み` と扱う。

- 実装が owner component 詳細本文責務の入力、出力、状態、処理順序、異常系、セキュリティ制約と一致している。
- `docs/DETAIL_INDEX.md` 詳細仕様入口責務の完全実装検証マトリクスと仕様策定完了チェックを満たしている。
- `docs/ROADMAP.md` 状態・計画責務の Phase 実装計画に従い、対象 Phase の依存条件、完了条件、PR 分割条件を満たしている。
- 対象機能が `docs/DETAIL_INDEX.md` 詳細仕様入口責務の機能仕様テンプレートと詳細節対応表の受け入れ条件を満たしている。
- API、SDK、UI のいずれかを変更した場合、`docs/details/api.md`、`docs/details/sdk.md`、`docs/details/ui.md` の対応関係が崩れていない。
- `docs/ROADMAP.md` の実装状態、`docs/DOCUMENT_INDEX.md`、`docs/SPEC.md`、詳細仕様の更新要否を確認済みである。
- 実装 PR 本文に、対象、実行コマンド、期待結果、実結果、判定を記録している。

検証不能な項目、未実行の項目、環境都合で省略した項目が残る場合、そのコンポーネントを `実装済み` として扱ってはならない。

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

承認済み変更の整合に直接必要な `docs/DOCUMENT_INDEX.md` 等の更新は、最初に提示した影響範囲に含まれている場合に限り、追加承認なしで行ってよい。

---

## 2. 仕様書管理ルール

`docs/SPEC.md` は、Adlaire CI の方針、ポリシー、正本関係、禁止事項、リリース判断を定めるマスター仕様書正本である。

`docs/ROADMAP.md` は、Adlaire CI の実装状態、実装可否、Phase、機能インベントリ、将来計画、昇格手順、§27 追加仕様化機能参照、横断補足契約を定めるロードマップ正本である。

`docs/DETAIL_INDEX.md` は、詳細仕様の入口、索引、共通固定値、実装前確認項目、検証マトリクス、詳細節対応表、リポジトリ内ソース配置、責務 component 別詳細仕様ファイル管理仕様を持つ詳細仕様入口正本である。

owner component 別の `docs/details/*.md` は、各 component の詳細仕様本文に関する正本である。

標準ソース配置は `main.go` と `components/*.go`、および `admin/` 配下の管理 UI ファイルとする。

標準ソース配置への実装移行は完了済みである。現行実装実体は `main.go`、`components/*.go`、`admin/` 配下の管理 UI 静的ファイル、`testdata/<component>/` とする。

`build_spec.go`、`runner.go`、`build_spec_test.go`、`runner_test.go`、`testdata/build_spec/` を現行実装実体として扱ってはならない。

実装済みコンポーネントである `components/api.go`、`admin/adlaire-ci-sdk.js`、`admin/index.html` は、`docs/SPEC.md`、`docs/ROADMAP.md`、`docs/DETAIL_INDEX.md`、該当する owner component 別の `docs/details/*.md` に基づいて更新する。

`components/mcp.go` は将来計画コンポーネントであり、実装状態、実装可否、ロードマップ状態は `docs/ROADMAP.md` を正とする。`docs/DETAIL_INDEX.md` の責務 component 対応表と該当する owner component 別詳細仕様に入出力、状態、起動手順、検証条件が定義されるまでは実装対象として扱わない。

`COMMON`、`CORE`、`BASE`、`SHARED`、`FOUNDATION`、その他の横断共通基盤ファイルは、詳細仕様ファイルとして作成してはならない。横断する固定値、読み順、対応表は `docs/DETAIL_INDEX.md` の入口・索引・共通固定値・管理仕様として扱う。§27 追加仕様化機能の横断補足契約は `docs/ROADMAP.md` §6 を正とする。いずれも component として扱わない。

`docs/DESIGN.md` は、出力 HTML のデザイン仕様を整理する補助文書である。機能仕様、実装状態、実装可否、API 仕様、状態 schema、builder 処理本文は定義しない。方針は `docs/SPEC.md`、状態判断は `docs/ROADMAP.md`、builder の視覚実装詳細は `docs/details/builder.md` を優先する。

仕様改訂では、既存仕様、`docs/DOCUMENT_INDEX.md`、実装ファイルとの整合性を確認する。

文書整理では、`README.md`、`docs/DOCUMENT_INDEX.md`、`docs/SPEC.md`、`docs/ROADMAP.md`、`docs/DETAIL_INDEX.md`、owner component 別の `docs/details/*.md` の読み順と正本範囲を維持する。

文書整理だけを目的とする作業では、機能仕様、実装状態、実装可否、ロードマップ状態を変更してはならない。状態変更が必要な場合は、変更対象、変更理由、影響範囲を別途提示し、承認を得る。

`README.md` は初見向けの入口、`docs/DOCUMENT_INDEX.md` は索引、`docs/SPEC.md` は方針・ポリシーの正本、`docs/ROADMAP.md` は状態・Phase・将来計画・§27 追加仕様化機能参照・横断補足契約の正本、`docs/DETAIL_INDEX.md` は詳細仕様入口、owner component 別の `docs/details/*.md` は詳細仕様本文として扱う。

`docs/ROADMAP.md`、`docs/DETAIL_INDEX.md`、owner component 別の `docs/details/*.md` に記載された一部コンポーネントや機能は、仕様化済みであっても未実装の場合がある。リポジトリ内に実装ファイルまたは実装コードが存在しない内容を、実装済み機能として扱ってはならない。

`docs/DETAIL_INDEX.md` または owner component 別の `docs/details/*.md` を改訂する場合は、方針、ポリシー、実装状態、実装可否、ロードマップ状態、PR 分割判断を記載してはならない。これらは `docs/SPEC.md` または `docs/ROADMAP.md` の該当正本範囲を正とする。

`docs/DETAIL_INDEX.md` と owner component 別の `docs/details/*.md` には、実装者が実装時に必要とする対象コンポーネント、入出力、設定値、データ構造、処理順序、異常系、状態管理、セキュリティ制約、検証条件だけを記載する。

未確定の内容を実装可能な詳細仕様として扱ってはならない。実装判断に必要な具体値、条件、処理が未確定の場合は、`docs/ROADMAP.md` の状態分類を確認し、`docs/DETAIL_INDEX.md` または owner component 別の `docs/details/*.md` へ推測で具体値を記載してはならない。

仕様項目の成熟度と実装可否は、`docs/SPEC.md` の仕様成熟度方針・仕様成熟度ポリシーに従って判定する。

作業開始時には、対象機能・対象コンポーネントについて `docs/SPEC.md`、`docs/ROADMAP.md`、`docs/DETAIL_INDEX.md`、該当する owner component 別の `docs/details/*.md` の該当節と実ファイルの存在を確認する。

実ファイルの存在確認には `rg --files` を使用する。

仕様上の状態は、以下のいずれかに分類して扱う。

- 未仕様化
- 将来計画
- 改訂予定
- 仕様化済み・未実装
- 実装中・検証未完了
- 実装済み

仕様化済み・未実装、実装中・検証未完了、改訂予定、将来計画、未仕様化の項目を、実装済みとして報告してはならない。

未仕様化、将来計画、改訂予定の項目は、実装着手可能な状態として扱ってはならない。

実装着手は、対象項目が仕様化済み・未実装の状態に到達している場合に限る。

実装中に仕様不足、未定義の入出力、未定義の状態ファイル、未定義の異常系、未定義の検証条件を発見した場合は、実装判断で補完せず、先に仕様を改訂する。

実装済みとして扱うには、コード変更、仕様との差分確認、必要な構文確認、実行確認または生成物確認、`docs/DOCUMENT_INDEX.md`・`docs/ROADMAP.md` の実装状態・関連仕様の更新要否確認を完了していなければならない。

API、SDK、標準管理ツールのいずれかを変更する場合は、API 仕様、SDK メソッド、UI 操作、詳細仕様の整合を同時に確認する。

未実装項目を実装する場合は、`docs/DOCUMENT_INDEX.md` の Specified Components と本ファイルの実装管理ルールの更新要否を確認する。

仕様の記載と実ファイルの存在が矛盾する場合は、先に仕様・索引・ルールブックの整合を取る。

一時ファイル、退避ファイル、比較用ファイルは、整合性確認が完了するまで削除しない。

不要ファイル削除は、削除対象、削除理由、影響範囲を提示し、別途 `承認` を得てから行う。

---

## 3. 実装管理ルール

Go 実装対象ファイルは以下とする。

| ファイル | 役割 |
|---------|------|
| `main.go` | 起動入口。現時点では実行ファイル名に応じて `builder` または `runner` component を呼び出す。 |
| `components/builder.go` | Markdown ファイルまたは Markdown ディレクトリを静的 Web サイトへ変換する Go 版ビルドスクリプト。 |
| `components/runner.go` | GitHub API で対象 Markdown の変更を検出し、ビルドパイプラインを実行する Go 版 CI ランナー。 |

実装済みコンポーネントおよび将来追加予定 path は以下とする。

| ファイル | 状態 |
|---------|------|
| `components/api.go` | 実装済み |
| `admin/adlaire-ci-sdk.js` | 実装済み |
| `admin/index.html` | 実装済み |
| `components/mcp.go` | 将来追加予定 path（未作成） |

未実装コンポーネントを追加する場合は、`docs/SPEC.md`、`docs/ROADMAP.md`、`docs/DETAIL_INDEX.md`、該当する owner component 別の `docs/details/*.md`、`docs/DOCUMENT_INDEX.md`、本ファイルを必要に応じて整合させる。

実装変更後は、変更範囲に応じて構文確認、実行確認、生成物確認を行う。

Go 実装の構文確認では、対象ファイルに対して `gofmt -l ...` を実行し、Go module が存在する場合は `go test ./...` を実行する。

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

上記確認で、同一ファイル、同一仕様領域、同一責務、または merge 順序依存の open Pull Request が見つかった場合は、新規 Pull Request を作成してはならない。既存 Pull Request への統合、または一本化 Pull Request への集約を先に完了する。

Pull Request の `mergeStateStatus` が `DIRTY`、`UNKNOWN`、または確認不能の場合は、merge 可能と報告してはならない。`UNKNOWN` の場合は GitHub の再計算後に再確認し、最終的に `CLEAN` を確認する。

競合解消時は、競合マーカーの除去だけで完了としてはならない。`git diff --check`、競合マーカー検索、変更対象文書の正本関係確認、open Pull Request 一覧確認を完了条件とする。

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

GitHub 設定を確認する場合は、少なくとも以下を確認する。

- `gh api repos/fqwink/Build-Scripts` で、`delete_branch_on_merge`、`allow_auto_merge`、`allow_update_branch`、`allow_merge_commit`、`allow_squash_merge`、`allow_rebase_merge`、`has_issues`、`has_projects`、`has_wiki`、`has_discussions`、`security_and_analysis` を確認する。
- `gh api repos/fqwink/Build-Scripts/branches/main/protection` で、`main` branch protection を確認する。
- `main` branch protection の確認で `Branch not protected` が返る場合は、未設定として扱う。

GitHub 設定を変更する前には、以下を必ず提示する。

- 変更対象
- 現在値
- 推奨値
- 影響範囲

GitHub 設定を変更した後は、GitHub API で再取得し、`AGENTS.md` の標準設定との差分がないかを確認して報告する。

標準 GitHub リポジトリ設定のうち、自動化に関わる設定が未確認の場合は、現在の設定状態を確認する。

自動化に関わる設定が未設定または標準値と異なる場合は、変更対象、変更内容、影響範囲を提示し、ユーザーから `承認` を得たうえで標準値へ設定する。

自動化に関わる設定には、少なくとも `delete_branch_on_merge=true` を含める。その他の自動化設定が `docs/SPEC.md` または本ルールブックで標準化された場合も同様に扱う。

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

上記手順で fast-forward できない場合、merge 対象やローカル変更の状態を確認し、勝手に履歴を書き換えてはならない。

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
- open Pull Request の `mergeStateStatus` が `DIRTY` または `UNKNOWN` の場合は、競合状態または未確認状態として扱い、`CLEAN` を確認するまで完了報告しない。
- 競合解消または PR 一本化を行った場合は、重複 PR が open のまま残っていないことを `gh pr list --state open` で確認する。
- 文書変更では、`rg` で不要になった名称、矛盾参照、不要になったファイル名が残っていないか確認する。
- 文書変更では、`git diff --stat` で変更範囲を確認する。
- ファイル追加、削除、リネームを含む場合は、`git diff --cached --summary` で Git 上の扱いを確認する。
- 実装変更では、対象言語に応じた構文確認を行う。Go 実装では `gofmt -l ...` を標準の整形確認とし、Go module が存在する場合は `go test ./...` を標準の確認とする。
- 実装変更では、必要に応じて対象スクリプトの実行確認または生成物確認を行う。
- 仕様変更では、`docs/SPEC.md`、`docs/ROADMAP.md`、`docs/DETAIL_INDEX.md`、該当する owner component 別の `docs/details/*.md`、`docs/DOCUMENT_INDEX.md`、`docs/DESIGN.md`、実装ファイルの整合を確認する。

Pull Request 本文には、少なくとも以下を記載する。

- `Summary`
- `Verification`
- 競合防止確認
- 未実施の確認がある場合は、その理由

---

## 5. 外部依存ルール

外部フレームワークおよび外部ライブラリは、原則として採用しない。

機能実現は、Go 標準ライブラリ、Vanilla JavaScript、内製実装、例外承認済み外部ライブラリの順で検討する。

外部依存は最小限に抑え、可能な範囲で内製化を重視する。

ただし、開発コスト、実装難易度、安全性、保守性、暗号・認証等の専門性を考慮し、外部ライブラリの採用を例外として許可する場合がある。

例外として外部ライブラリを採用する場合は、採用理由、対象範囲、影響範囲、代替困難性、保守方針を明示する。

例外採用は、`docs/SPEC.md` または承認済み変更範囲に明記された場合のみ有効とする。

外部依存を追加、削除、更新、置換する作業は変更作業として扱い、事前承認を必須とする。

---

## 6. 文書整合ルール

`docs/DOCUMENT_INDEX.md` は、リポジトリ内の文書・実装ファイルの役割を示す索引として維持する。

ファイル名、正本関係、実装コンポーネントの追加・削除・リネームが発生した場合は、`docs/DOCUMENT_INDEX.md` の更新要否を確認する。

`docs/SPEC.md`、`docs/ROADMAP.md`、`docs/DETAIL_INDEX.md`、または owner component 別の `docs/details/*.md` を改訂した場合は、`docs/DESIGN.md`、`docs/DOCUMENT_INDEX.md`、実装ファイルへの影響を確認する。

`docs/DESIGN.md` を改訂した場合は、`components/builder.go` 内の HTML / CSS / JavaScript / theme component テンプレートとの整合性を確認する。

仕様化済み項目を実装した場合は、`docs/ROADMAP.md` の状態分類と実装状態、`docs/DOCUMENT_INDEX.md` の Specified Components、`docs/DETAIL_INDEX.md` の対応表、owner component 別の `docs/details/*.md` の検証条件、実装ファイルの存在を整合させる。
