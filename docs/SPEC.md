# Adlaire CI — 仕様ドキュメント

**仕様対象コンポーネント：** `components/builder.go` / `components/runner.go` / `components/api.go` / `admin/adlaire-ci-sdk.js` / `admin/index.html`
**将来計画コンポーネント：** `components/mcp.go`
**出力形式：** 静的 Web サイト（HTML / CSS / JavaScript / search index）
**スクリプトバージョン：** v3
**仕様バージョン：** V.N（正式リリース前の暫定表記）/ **リリースバージョン：** V.X.N（正式リリース前の暫定表記）。[`docs/SPEC.md`](SPEC.md) ポリシー責務 §1 参照。
**最終更新：** 2026-09-15

---

> **Adlaire CI** とは、最初から Go を前提として仕様策定するビルド・CI・管理システムの総称である。仕様対象コンポーネントは `components/builder.go`、`components/runner.go`、`components/api.go`、`admin/index.html`、`admin/adlaire-ci-sdk.js` とする。`components/mcp.go` は将来計画コンポーネントであり、実装対象コンポーネントとして扱わない。各 component の実装状態、実装可否、将来計画状態は [`docs/ROADMAP.md`](ROADMAP.md) の状態・計画責務を正本とする。

## 文書責務

[`AGENTS.md`](../AGENTS.md) と [`docs/SPEC.md`](SPEC.md) は、本リポジトリにおける最上位文書である。

[`AGENTS.md`](../AGENTS.md) は、作業ルール、承認、Git 操作、Pull Request 作成、レビュー対応、検証手順、エージェント実行手順の最上位ルールブックである。

[`docs/SPEC.md`](SPEC.md) は、仕様、方針、ポリシー、正本参照先、禁止事項、リリース判断、実装着手可否の最上位仕様書である。

作業ルール上の矛盾は [`AGENTS.md`](../AGENTS.md) を正とする。

仕様、方針、ポリシー、正本参照先、禁止事項、リリース判断、実装着手可否の矛盾は [`docs/SPEC.md`](SPEC.md) を正とする。

[`docs/SPEC.md`](SPEC.md) は、なぜその運用が必要か、仕様として何を禁止するか、どの状態を完了扱いにしてはならないかを定義する。[`AGENTS.md`](../AGENTS.md) は、承認、実行コマンド、Git 操作、PR 作成、検証手順をどう実行するかを定義する。

完了扱いの可否、実装着手可否、仕様変更単位、禁止事項の判断で迷う場合は [`docs/SPEC.md`](SPEC.md) を正とする。実行順序、コマンド、確認手順、PR 操作手順の判断で迷う場合は [`AGENTS.md`](../AGENTS.md) を正とする。

[`docs/SPEC.md`](SPEC.md) は、Adlaire CI の方針責務・ポリシー責務の正本である。

[`docs/SPEC.md`](SPEC.md) は、目的、設計方針、禁止事項、成熟度、実装着手可否、仕様判断の原則を扱う。ただし、生成 HTML のデザイン関係は [`docs/DESIGN.md`](DESIGN.md) を正本とする。

方針またはポリシーに該当する内容は、必ず [`docs/SPEC.md`](SPEC.md) に記載する。ただし、生成 HTML のデザイン関係は例外として [`docs/DESIGN.md`](DESIGN.md) に記載する。これら以外の方針またはポリシーに該当する内容を、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、owner component 別の [`docs/details/*.md`](details/)、[`docs/ROADMAP.md`](ROADMAP.md)、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、[`README.md`](../README.md)、実装ファイル、fixture、PR 本文へ移してはならない。

[`docs/SPEC.md`](SPEC.md) の `技術方針` 表と `ディレクトリ構成` tree は、Adlaire CI の方針責務本文である。削除、他文書への移動、参照だけへの置き換えを禁止する。これらは実装詳細、状態一覧、索引本文として扱わない。

[`docs/SPEC.md`](SPEC.md) は、関数単位の処理、HTTP response schema、状態ファイル schema、SDK method、UI DOM、fixture assertion、具体的な実行手順を定義しない。これらの本文は、責務ベース明示的原則、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務に従って、責務を持つ正本へ分離する。

文書を整理する場合も、[`docs/SPEC.md`](SPEC.md) の正本範囲を越えてはならない。[`docs/SPEC.md`](SPEC.md) は、詳細仕様本文、実装状態、ロードマップ状態、文書所在を重複定義しない。

## 状態参照方針

状態責務は [`docs/ROADMAP.md`](ROADMAP.md) を正本とする。

[`docs/SPEC.md`](SPEC.md) では、状態分類、Phase 実装単位、実装着手可否の方針とポリシーだけを定義する。個別 component の現在状態、Phase 一覧、将来計画一覧、昇格手順は重複定義しない。

---

## 責務文書構成

| 責務 | 正本 | 責務の問い | 記載する内容 |
|------|------|-----------|------------|
| 方針責務 | [`docs/SPEC.md`](SPEC.md) | **なぜ・何を** | 目的、設計思想、方向性の原則。変更頻度が低く、判断の拠り所となる指針。 |
| ポリシー責務 | [`docs/SPEC.md`](SPEC.md) | **しなければならない／してはならない** | 遵守義務のある規則、制約、禁止事項、セキュリティ要件、運用ルール、バージョン管理規則。 |
| 状態・計画責務 | [`docs/ROADMAP.md`](ROADMAP.md) | **いつ・どれを** | 実装状態、実装可否、Phase、機能インベントリ、将来計画、昇格手順。 |
| 詳細仕様入口責務 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) | **どこから読むか** | 詳細仕様の入口、読み順、共通固定値、詳細節対応表、リポジトリ内ソース配置。 |
| owner component 別詳細本文責務 | [`docs/details/*.md`](details/) | **どのように実装するか** | owner component 別の入出力、状態、処理順序、異常系、検証条件。 |
| デザイン責務 | [`docs/DESIGN.md`](DESIGN.md) | **どう見せるか** | 生成 HTML のデザイン関係。 |
| 文書・実装ファイル所在の索引責務 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) | **どこにあるか** | 文書、実装ファイル、生成物、将来追加予定 path の所在。 |
| 利用入口責務 | [`README.md`](../README.md) | **どう始めるか** | 利用者向け入口、概要、参照先。 |

新しい記載内容は「この内容はどの責務の問いに答えるか」を基準に、責務を持つ正本を決定する。

判断対象が方針またはポリシーである場合、責務を持つ正本は必ず [`docs/SPEC.md`](SPEC.md) とする。ただし、判断対象が生成 HTML のデザイン関係である場合、責務を持つ正本は [`docs/DESIGN.md`](DESIGN.md) とする。判断対象が実装詳細、実装状態、文書所在、利用入口、検証証跡のいずれかである場合だけ、該当責務の正本を参照する。

実装者は、実装前に以下の参照順序方針に従う。

1. [`docs/SPEC.md`](SPEC.md) の「文書責務」と「状態参照方針」で、正本範囲を確認する。
2. [`docs/ROADMAP.md`](ROADMAP.md) で、対象の状態、実装可否、Phase、将来計画該当有無を確認する。
3. [`docs/SPEC.md`](SPEC.md) 方針責務 §4.1〜§4.10 で、ゼロ依存、責務ベース明示的原則、ディレクトリ構成、詳細仕様粒度、成熟度、着手ゲート、完了判定、Go 正本方針を確認する。
4. [`docs/SPEC.md`](SPEC.md) のポリシー責務で、対象領域の禁止事項、セキュリティ、バージョン、外部依存を確認する。
5. 生成 HTML のデザイン関係を扱う場合は、[`docs/DESIGN.md`](DESIGN.md) デザイン責務で視覚仕様を確認する。
6. [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務で読み順、共通固定値、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1〜§0i.4 の対応表、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0j のリポジトリ内ソース配置を確認し、該当する owner component 別の [`docs/details/*.md`](details/) 詳細本文責務で実装に必要な入出力、状態、異常系、検証条件、配置を確認する。

---

# 方針責務
> 目的・設計思想・方向性の原則を定める。「なぜこう作るか」に答える。

方針責務の節を参照する場合は、必ず [`docs/SPEC.md`](SPEC.md) 方針責務 §番号 の形式で記載する。番号だけ、または `docs/SPEC.md §番号` だけで参照してはならない。

## 1. 目的

`components/builder.go` は、任意の UTF-8 Markdown ファイルまたは Markdown ディレクトリを、静的配信可能な Web サイトへ変換する Go プログラムである。[`docs/SPEC.md`](SPEC.md) 方針責務では、Go 実装を最初からの正本として定義する。

- 大規模 Markdown 仕様書、複数 Markdown ドキュメント、運用メモを、静的 Web サイトへ変換する
- `index.html`、ページ HTML、共通 CSS、共通 JavaScript、検索 index を出力ディレクトリへ生成する

## 2. 開発方針

本プロジェクトは **仕様駆動開発（Spec-Driven Development）** を採用する。

- **責務正本群が判断基準**：実装の追加・変更は、変更対象の責務を持つ正本への反映を先行させる
- **仕様から実装への順序**：実装は、[`docs/SPEC.md`](SPEC.md)、[`docs/ROADMAP.md`](ROADMAP.md)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務の内容に基づいて修正する
- **責務正本との乖離は不整合**：乖離が生じた場合も、責務を持つ正本を先に改訂し、その正本に基づいて実装を修正する

## 4. 技術方針

本表は Adlaire CI の技術選定に関する方針責務本文である。削除、他文書への移動、参照だけへの置き換えを禁止する。詳細仕様は本表の方針に従い、具体的な入出力、処理順序、例外条件、検証条件を owner component 別の [`docs/details/*.md`](details/) 詳細本文責務へ記載する。

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

### 4.1 ゼロ依存・フルインハウス原則

Adlaire CI は、ゼロ依存・フルインハウスを技術哲学の中核とする。

本原則におけるゼロ依存とは、各コンポーネントが外部ライブラリ、外部フレームワーク、外部ビルドツール、外部ホスティング実行基盤に機能成立を依存しないことを意味する。本原則におけるフルインハウスとは、Markdown 変換、CI 実行、管理 API、SDK、標準管理ツール、状態管理、認証、ログ、通知、セットアップの主要機能を本リポジトリ内で仕様化し、内製コードとして理解、検証、保守できる状態を意味する。

実装者は、便利さ、実装速度、一般的な慣習を理由に外部ライブラリで未定義機能を補完してはならない。外部依存がなければ成立しない設計は設計不備として扱い、先に仕様を見直す。

各コンポーネントの自律性は以下を満たす。

| コンポーネント | 自律性の条件 |
|----------------|--------------|
| `components/builder.go` | Go 標準ライブラリだけで Markdown 解析、HTML/CSS/JS/search index 生成、検証レポート出力を行う。外部 Markdown parser、template engine、syntax highlight library、search library に依存しない。 |
| `components/runner.go` | Go 標準ライブラリと OS 標準コマンドだけで GitHub API polling、SHA 比較、ビルド起動、ログ、通知、SSH 転送、snapshot、lock、retry を処理する。外部 CI サービス、job queue、scheduler library に依存しない。 |
| `components/api.go` | Go 標準ライブラリ `net/http` を基本に、認証、session、状態ファイル CRUD、入力検証、API response を内製実装する。外部 web framework、router、ORM、database driver に依存しない。 |
| `admin/adlaire-ci-sdk.js` | 単一 ES Module とし、browser 標準 API のみで API client、error handling、streaming、timeout を実装する。npm package、bundler、polyfill、framework に依存しない。 |
| `admin/index.html` | HTML / CSS / Vanilla JavaScript だけで標準管理ツールを構成し、SDK 経由で通信する。React、Vue、Svelte、CSS framework、icon package、chart library に依存しない。 |
| `components/mcp.go` | 将来計画コンポーネントとしての自律性方針を示す。将来計画の段階でも Go 標準ライブラリを前提とし、MCP 通信、JSON-RPC 処理、API bridge、監査ログを内製する。外部 MCP framework に依存する前提で仕様化しない。本行は実装着手許可、詳細仕様成立、実ファイル作成許可を意味しない。 |

本表はコンポーネント自律性の方針であり、関数単位の処理、入出力、状態、異常系、検証条件を定義するものではない。各コンポーネントの具体的な実装契約は、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を正本とする。

外部依存を例外採用する場合は、[`docs/SPEC.md`](SPEC.md) ポリシー責務 §4 の許可外部ライブラリ一覧へ登録し、採用理由、代替困難性、責務範囲、削除方針、検証条件を同一 PR で明記する。許可リストにない外部依存は、実装済みとして受け入れない。

### 4.2 Core 非採用・共通責務コンポーネント方針

Adlaire CI は、`core`、`adlaire-ci-core`、`internal/core`、`common`、`base`、`foundation`、`utils` のような中心化・汎用置き場化する概念を採用しない。

横断的に利用される処理は、単一責務、入出力、状態、異常系、検証条件、依存方向が仕様化されている場合に限り、共通責務コンポーネントとして定義できる。ただし、共通責務コンポーネントは他コンポーネントより上位ではなく、同格の独立コンポーネントとして扱う。共通責務コンポーネントを、正本、中心、基盤、親、上位レイヤーとして扱ってはならない。

共通責務コンポーネントを追加する場合は、以下をすべて満たす。

- 単一責務である。
- [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務に参照入口があり、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務に入力、出力、状態、異常系、検証条件が仕様化されている。
- 利用元コンポーネントとの依存方向が明確である。
- アプリケーション固有の判断、業務判断、UI 判断、API endpoint 判断、ビルド対象判断を含まない。
- 何でも置き場として使用しない。
- `core`、`common`、`base`、`foundation`、`utils` など中心性または汎用置き場を示す名称を使わない。
- 追加時に [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、[`docs/SPEC.md`](SPEC.md) ポリシー責務 §4、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を整合する。

共通責務コンポーネントの命名候補を仕様化する場合は、以下の形式で責務を明示する。

| 命名候補 | 責務 |
|----|------|
| `state_store` | 状態ファイルの atomic write、lock、JSON 読み書き。 |
| `secret_masker` | PAT、Webhook Secret、SMTP password、session token、API token のマスク処理。 |
| `event_log` | 構造化イベントログの形式、出力、読み取り。 |
| `github_client` | GitHub REST API 呼び出し、rate limit、retry 境界。 |

本表は、共通責務コンポーネントを追加する場合の命名・責務明示ポリシーである。本表に含まれる名称だけを理由に、実装ファイル、package、directory、状態項目、placeholder を作成してはならない。追加可否、実装状態、所在は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務と [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務を正本とする。

禁止される設計は以下とする。

| 禁止対象 | 理由 |
|----------|------|
| `core` / `adlaire-ci-core` / `internal/core` | 中心コンポーネント化し、正本や上位概念と誤読されるため。 |
| `common` / `base` / `foundation` / `utils` | 責務が曖昧になり、何でも置き場化しやすいため。 |
| 全共通処理の一括集約 package | 単一責務を失い、コンポーネント境界を壊すため。 |
| 業務判断を含む共通処理 | 利用元コンポーネントの責務を横断基盤へ漏らすため。 |

共通責務コンポーネントは、コード共有のためだけに追加してはならない。重複削減より、責務境界、仕様の明確さ、依存方向の追跡可能性を優先する。

### 4.2a 責務ベース明示的原則

Adlaire CI の仕様体系は、責務ベース明示的原則を仕様全般の最上位方針として採用する。

本原則は [`docs/SPEC.md`](SPEC.md) 全体に適用する。[`docs/SPEC.md`](SPEC.md) 内の各記載は、方針責務、ポリシー責務、状態・計画責務の参照、詳細仕様入口責務の参照、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務の参照、文書・実装ファイル所在の索引責務の参照、利用入口責務の参照のいずれかとして読めなければならない。

責務ベース明示的原則とは、方針、ポリシー、実装状態、実装可否、Phase、将来計画、詳細仕様本文、fixture、expected、fake、検証証跡、文書索引を、それぞれ異なる責務として明示的に分離し、同一判断対象を複数文書で正本化しない原則である。

正本参照先は、必ず責務名とファイル名で示す。文書を章構成、便宜分類、または他文書の従属章として扱ってはならない。[`docs/SPEC.md`](SPEC.md) 内部の見出しも責務名で示し、`Part` 名称で扱ってはならない。[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、[`docs/details/*.md`](details/)、[`docs/ROADMAP.md`](ROADMAP.md)、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、[`README.md`](../README.md) を [`docs/SPEC.md`](SPEC.md) の章として扱ってはならない。

参照はリンク化を必須とする。文書間参照、節参照、表参照、責務正本参照、実装ファイル参照、fixture 参照を書く場合は、Markdown link を用いて参照先へ移動できる形にする。単なるファイル名、裸の節番号、裸の見出し名、または `参照` という文字だけで参照先を示した扱いにしてはならない。

リンク化する参照は、表示文言に責務名、ファイル名、節番号、または対象名を含める。例として、状態・計画責務を参照する場合は [`docs/ROADMAP.md`](ROADMAP.md)、詳細仕様入口責務を参照する場合は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、owner component 別詳細本文責務を参照する場合は [`docs/details/*.md`](details/)、文書・実装ファイル所在の索引責務を参照する場合は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、利用入口責務を参照する場合は [`README.md`](../README.md) のように記載する。リンク化できない生成物、PR 本文、外部ツール出力で参照を記録する場合でも、参照先ファイル名と節番号を省略してはならない。

方針とポリシーは [`docs/SPEC.md`](SPEC.md) だけに記載する。ただし、生成 HTML のデザイン関係は [`docs/DESIGN.md`](DESIGN.md) だけに記載する。詳細仕様、ロードマップ、索引、README、実装ファイル、fixture、PR 本文は、方針またはポリシーを本文として定義、補足、緩和、例外化、再解釈してはならない。

方針またはポリシーに該当する記載を [`docs/SPEC.md`](SPEC.md) から削除し、他文書への参照だけに置き換えてはならない。ただし、生成 HTML のデザイン関係は [`docs/SPEC.md`](SPEC.md) に本文を置かず、[`docs/DESIGN.md`](DESIGN.md) に記載する。方針またはポリシーを整理する場合は、[`docs/SPEC.md`](SPEC.md) 内で責務名、適用範囲、禁止事項、参照先を明示して整える。

[`docs/SPEC.md`](SPEC.md) の `技術方針` 表と `ディレクトリ構成` tree は、方針責務本文である。これらを実装詳細、状態一覧、索引本文とみなして削除、他文書へ移動、または参照だけへ置き換えてはならない。

詳細仕様は実装詳細だけを記載する。ロードマップは状態、実装可否、Phase、将来計画だけを記載する。索引は文書と実装ファイルの所在だけを記載する。README は利用入口だけを記載する。fixture、expected、fake、検証証跡、完了判定は検証責務だけを記載する。

各文書は、自分の責務外の内容を本文として記載してはならない。責務外の内容が必要な場合は、責務を持つ正本と節番号を参照する。参照先に必要な内容が存在しない場合は、参照元へ補足を書かず、責務を持つ正本側を改訂する。責務を持つ正本側で確定できない内容は未確定として扱い、仕様化済みとして扱ってはならない。

分かりやすさは、同じ説明を複数文書へ重複記載することで担保してはならない。分かりやすさは、責務名、責務を持つ正本、正本範囲、参照節を明示することで担保する。責務名を使わない正本説明、判断に迷う記載、二重に読める記載、例外に見える記載、参照先を持たない責務外説明は禁止する。

owner component は対象機能の詳細本文を持つ。collaborator component は、境界、接続、入力受け渡し、出力受け渡し、検証観点として参照される。collaborator component は、owner component の本文を置き換えたり、同じ判断対象を別正本として再定義したりしてはならない。

fixture、expected、fake、検証証跡、完了判定は検証責務として扱い、[`docs/details/fixture.md`](details/fixture.md) を正本とする。実装状態、実装可否、Phase、将来計画は状態責務として扱い、[`docs/ROADMAP.md`](ROADMAP.md) を正本とする。詳細仕様の入口、対応表、共通固定値、リポジトリ内ソース配置は詳細仕様入口責務として扱い、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) を正本とする。生成 HTML のデザイン関係はデザイン責務として扱い、[`docs/DESIGN.md`](DESIGN.md) を正本とする。文書配置と実装ファイル所在は、文書・実装ファイル所在の索引責務として扱い、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) を正本とする。

本原則は、[`docs/DESIGN.md`](DESIGN.md)、[`docs/ROADMAP.md`](ROADMAP.md)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、owner component 別の [`docs/details/*.md`](details/)、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、[`README.md`](../README.md) の記載整理より上位の方針である。これらの文書を整理する場合は、本原則に従い、重複本文を増やさず、責務と参照先を明示する。

本原則への違反が残る状態では、仕様整合完了、仕様 PR 完了、実装着手、実装済み判定を行ってはならない。

### 4.3 ディレクトリ構成

Adlaire CI のディレクトリ構成は、責務ベースで整理する。

本節のディレクトリ構成 tree は、Adlaire CI のディレクトリ構成を明示する方針責務本文である。削除、他文書への移動、参照だけへの置き換えを禁止する。[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) は所在索引として本節に従う。

ディレクトリ構成は以下とする。

```text
.
├── main.go
│
├── components/
│   ├── builder.go
│   ├── runner.go
│   └── api.go
│
├── admin/
│   ├── index.html
│   └── adlaire-ci-sdk.js
│
├── testdata/
│   ├── builder/
│   └── runner/
│
├── docs/
│   ├── SPEC.md
│   ├── DETAIL_INDEX.md
│   ├── DOCUMENT_INDEX.md
│   ├── DESIGN.md
│   ├── details/
│   │   ├── builder.md
│   │   ├── runner.md
│   │   ├── api.md
│   │   ├── sdk.md
│   │   ├── ui.md
│   │   ├── admin.md
│   │   ├── setup.md
│   │   ├── statefile.md
│   │   ├── security.md
│   │   ├── archive.md
│   │   ├── commitstatus.md
│   │   └── fixture.md
│   └── examples/
│
├── README.md
├── AGENTS.md
└── go.mod
```

`main.go` は 1 ファイルとし、起動入口、実行ファイル名判定、引数受け取り、対象 owner component 呼び出しだけを担当する。`main.go` に Markdown 変換、CI 実行、HTTP handler、状態ファイル操作、archive 処理、GitHub Commit Status 送信、MCP 処理の実装詳細を書いてはならない。

`components/` は、1 実装対象コンポーネント = 1 Go ファイルとする。Go 実装ファイルの所在は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務、実装状態は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を正本とする。`admin` は Go コンポーネントではなく `admin/` 配下の静的配布物として扱う。`statefile`、`archive`、`commitstatus` は詳細仕様上の責務境界であり、単独 Go ファイルを作成する場合は該当 Phase または追加実装 PR で仕様状態と索引を更新してから追加する。`mcp.go` の状態、実装可否、追加条件は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務と該当 owner component の [`docs/details/*.md`](details/) 詳細本文責務を正本とする。

`admin/` は標準管理 UI の静的ファイルを配置する。`testdata/` は責務別 fixture を配置する。`docs/examples/` は利用例、設定例、サンプル構成を配置する。

ディレクトリ構成は方針上の到達形を示す。実ファイルの有無、追加予定 path、追加時期、状態分類は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務と [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務を正本とする。ディレクトリ構成に含まれることだけを理由に、未実装ファイル、将来追加予定 path、空ディレクトリ、placeholder を作成してはならない。

### 4.4 詳細仕様方針

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) は、詳細仕様入口責務として、実装者が追加判断なしに該当する詳細本文へ到達できる粒度で記載する。

owner component 別の [`docs/details/*.md`](details/) は、詳細本文責務として、抽象的な方針や目的の再掲ではなく、実装時に必要な具体値、処理順序、入出力、状態、失敗時の扱いを定義する。

仕様化済み・未実装の項目であっても、実装予定として扱う場合は実装者が迷わない粒度まで詳細化する。実装時期、設計判断、具体値が未確定の内容は、実装可能な仕様として扱わず、未仕様化または将来計画として明示する。

詳細仕様入口責務と詳細本文責務の組み合わせは、少なくとも以下の問いに答えられる状態を満たす。

- どのコンポーネントが責務を持つか
- どのファイル、API、関数、設定値、状態ファイルを使用するか
- 入力、出力、データ構造、既定値、許容値は何か
- 正常系の処理順序は何か
- 異常系、再試行、ロック、冪等性、タイムアウトをどう扱うか
- セキュリティ上の禁止事項、保存してよい情報、保存してはならない情報は何か
- 実装完了をどの確認条件で判定するか

### 4.5 仕様成熟度方針

仕様項目は、実装可否を判断できるように成熟度を明確に区分する。

本節は成熟度を必要とする理由と判断方針を定める。状態名、状態定義、実装可否、昇格条件は、[`docs/SPEC.md`](SPEC.md) ポリシー責務 §0a の仕様成熟度ポリシーを正本とする。

成熟度は、構想の有無ではなく、実装者が実装に着手できるだけの情報が揃っているかで判定する。

将来計画や未確定アイデアは、実装対象として扱わない。実装対象にする場合は、先に詳細仕様を整え、仕様化済み・未実装へ昇格させる。

実装中に仕様不足、実装者判断に依存する分岐、未定義の入出力、未定義の状態ファイル、未定義の異常系を発見した場合は、実装判断で補完せず、仕様改訂へ戻す。

### 4.6 Phase 実装単位方針

Adlaire CI の実装順序、実装計画、実装 PR、完了判定は Phase 単位でのみ管理する。

本節は Phase を実装単位とする方針を定める。Phase 単位で行わなければならない事項、優先度ラベルを実装単位として使ってはならない事項、Phase 途中追加の禁止事項は、[`docs/SPEC.md`](SPEC.md) ポリシー責務 §0f の Phase 実装単位ポリシーを正本とする。

Phase は、対象 owner component、実装範囲、依存条件、完了条件、検証条件が明確な実装単位である。Phase の一覧、順序、依存条件、完了条件は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §4 を正本とする。

`P0`、`P1`、`P2〜P5` などの優先度ラベル、抽象段階、API 内部分類、fixture 分類を、実装単位、PR 単位、完了判定単位として扱ってはならない。

個別 Phase の対象範囲、順序、依存条件、完了条件は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を正本とする。[`docs/SPEC.md`](SPEC.md) は、Phase 単位でのみ実装する方針と、優先度ラベルを実装単位として扱わないポリシーだけを定義する。

将来計画、改訂予定、未仕様化の機能は Phase に含めない。対象機能を Phase に含める場合は、先に詳細仕様、検証条件、受け入れ条件を整え、`仕様化済み・未実装` へ昇格させる。

### 4.7 実装着手ゲート方針

実装者は、対象コンポーネントについて以下をすべて満たすまで実装に着手してはならない。

実装順序、実装計画、実装 PR、完了判定は、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.6 の Phase 実装単位方針と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0f の Phase 実装単位ポリシーに従う。

| 判定項目 | 着手条件 |
|----------|----------|
| 正本確認 | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務の該当節を確認済みである。 |
| 状態分類 | 対象が `仕様化済み・未実装`、`実装中・検証未完了`、または `実装済み` に分類され、`未仕様化`、`将来計画`、`改訂予定` ではない。 |
| 詳細節対応 | 対象機能が [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1〜§0i.4 の詳細節対応表に含まれ、該当する詳細仕様節と受け入れ条件を確認済みである。 |
| テンプレート充足 | 対象機能の詳細仕様が [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h の機能仕様テンプレートに必要な項目を満たしている。 |
| 責務境界 | 対象コンポーネント、呼び出し元、呼び出し先、状態ファイル、外部接続先が明確である。 |
| 契約同期 | API、SDK、UI、状態ファイル、セットアップ、検証条件のうち関係する仕様が同時に整合している。 |
| 禁止事項 | 外部依存、秘密情報、直接 API 呼び出し、互換処理、暗黙フォールバックなどの禁止事項が明確である。 |
| 完了条件 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §4、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0e、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0f、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1〜§0i.4、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0j および [`docs/details/setup.md`](details/setup.md) 詳細本文責務 §26.7 の受け入れ条件で、実装完了、実装順序、配置を判定できる。 |

着手条件を満たさない場合、実装者はコードで補完せず、先に仕様改訂を行う。実装 PR では、着手前に参照した詳細仕様節を PR 本文へ明記する。

### 4.8 完了判定方針

実装完了は、実装ファイルの作成やテスト成功だけでは成立しない。以下をすべて満たした場合にのみ `実装済み` と扱う。

本節は実装済み判定の方針を定める。仕様成熟度、仕様 PR 完了、仕様凍結、初期実装スコープ、Phase 実装単位の具体的な必須条件は、[`docs/SPEC.md`](SPEC.md) ポリシー責務 §0a〜§0f を正本とする。

- 実装が owner component 別の [`docs/details/*.md`](details/) 詳細本文責務の入力、出力、状態、処理順序、異常系、セキュリティ制約と一致している。
- [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0e の完全実装検証マトリクスと [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0f の仕様策定完了チェックを満たしている。
- [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §4 の Phase 実装計画に従い、対象 Phase の依存条件、完了条件、PR 分割条件を満たしている。
- 対象機能が [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h の機能仕様テンプレートを満たし、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1〜§0i.4 の詳細節対応表に記載された受け入れ条件を満たしている。
- セットアップまたは運用手順に影響する場合、[`docs/details/setup.md`](details/setup.md) 詳細本文責務 §26.7 の実装受け入れ条件を満たしている。
- API、SDK、UI のいずれかを変更した場合、[`docs/details/api.md`](details/api.md) 詳細本文責務 §22、[`docs/details/sdk.md`](details/sdk.md) 詳細本文責務 §23、[`docs/details/ui.md`](details/ui.md) 詳細本文責務 §24 の対応関係が崩れていない。
- [`docs/ROADMAP.md`](ROADMAP.md) の実装状態、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、[`docs/SPEC.md`](SPEC.md)、詳細仕様の更新要否を確認済みである。
- 実装 PR 本文に、対象、実行コマンド、期待結果、実結果、判定を記録している。

検証不能な項目、未実行の項目、環境都合で省略した項目が残る場合、そのコンポーネントを `実装済み` として扱ってはならない。

### 4.9 仕様策定単位方針

仕様策定は、実装者が 1 つの責務単位として読める範囲でまとめる。

同一責務に属する API、SDK、UI、状態ファイル、セットアップ、検証条件は、同じ仕様策定単位として扱う。これらを別々の仕様 PR に分割して、片方だけが先に merge される状態を作ってはならない。

仕様策定単位は以下のいずれかに分類する。

| 単位 | 含める内容 | 分割可否 |
|------|------------|----------|
| コンポーネント単位 | `components/builder.go`、`components/runner.go`、`components/api.go`、`admin/adlaire-ci-sdk.js`、`admin/index.html` の単一責務変更 | 同一コンポーネント内で完結する場合のみ単独 PR 可 |
| 横断契約単位 | API / SDK / UI / 状態ファイル / 認証 / セットアップの対応関係 | 分割不可。同一 PR で同期する |
| 運用単位 | systemd、セットアップ、アップデート、rollback、受け入れ条件 | 分割不可。手順と検証を同一 PR に含める |
| 将来計画単位 | 実装時期未定の方向性、候補、未確定案 | 実装可能仕様と混在禁止。将来計画として明示する |

仕様策定時に詳細が不足する項目は、推測で仕様化済みへ昇格してはならない。具体値、処理順序、異常系、検証条件を確定できない場合は、`将来計画` または `未仕様化` として残す。

### 4.10 Go 正本策定方針

[`docs/SPEC.md`](SPEC.md) 方針責務は、Adlaire CI を最初から Go 言語で設計・実装する前提で策定する。

実装者は、方針、ポリシー、Go 前提、禁止事項を [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務、実装状態と標準配置の状態判断を [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、詳細仕様の入口とソース配置を [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務、各 component の入出力、状態、処理順序、異常系、検証条件を owner component 別の [`docs/details/*.md`](details/) 詳細本文責務に従って判断する。過去の実装、試作、他言語スクリプト、既存ファイル名、既存 CLI、既存ログ、既存生成物、既存状態ファイルを前提にしてはならない。

Go 実装の判断基準は以下とする。

- `components/builder.go`、`components/runner.go`、`components/api.go` を Go 仕様対象コンポーネントとして扱う。
- `adlaire-ci-build`、`adlaire-ci-runner`、`adlaire-ci-api` を標準実行バイナリ名とする。
- 仕様未記載の自動変換処理、暗黙の読み替え処理、仕様外分岐を実装判断で追加してはならない。
- 該当する責務正本に記載されていない挙動は、仕様対象外として扱う。

## 5. CI ランナー方針

### 5.1 CI ランナーの目的

GitHub API を定期的にポーリングし、対象変更を検出してビルドパイプラインを自動実行する自己ホスト型 CI ランナー（`components/runner.go`）。GitHub Actions・Webhook・外部 CI サービスへの依存をゼロにする。

- 変更検出の具体的な API、比較値、保存先は [`docs/details/runner.md`](details/runner.md) 詳細本文責務を正本とする
- 対象変更が存在する場合だけビルドを実行する
- 外部公開エンドポイント・リバースプロキシ不要
- Go 標準ライブラリを基本とし、外部依存を追加する場合は [`docs/SPEC.md`](SPEC.md) ポリシー責務 §4 の例外承認を必須とする

### 5.2 CI ランナーの開発方針

- **単一責務実装**：`components/runner.go` は CI ランナー責務に限定し、Markdown 変換と管理 API を内包しない
- **シンプル性優先**：HTTP サーバー不要。1 回実行して終了する oneshot 設計
- **差分検出**：変更がない場合はビルドをスキップ

### 5.3 GitHub Actions 非依存方針

Adlaire CI は GitHub Actions の workflow、hosted runner、Marketplace action、push webhook を前提にしない。内製ランナーの起動方式、pipeline 実行方式、secret 参照、変更検出、通知、転送は [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/setup.md`](details/setup.md) 詳細本文責務を正本とする。

## 6. 管理ツール方針

### 6.1 管理ツールの目的

Adlaire CI の状態確認・操作を行う管理インターフェース。ヘッドレスアーキテクチャにより、フロントエンドとバックエンドを明確に分離する。

この節以降の管理ツール・管理 API・SDK に関する記載は、実装状態を固定しない機能方針である。対象ファイルの存在、実装状態、実装可否、分類、検証結果は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を正本とする。

### 6.2 ヘッドレスアーキテクチャ方針

Adlaire CI と管理ツールは API を介して通信する。フロントエンドとバックエンドを完全に分離し、管理ツールの実装・置き換えをバックエンドから独立させる。

- バックエンド（Adlaire CI）は API を公開する
- フロントエンド（管理ツール）は API のみを通じてバックエンドと通信する
- 直接のファイル操作・プロセス呼び出しは管理ツールから行わない

### 6.3 SDK 方針

API は SDK として提供し、管理ツール実装者が直接 HTTP 通信を記述しなくてよい抽象化レイヤーを提供する。

- **初期対応言語**：JavaScript のみ
- **フレームワーク非依存**：バニラ JS・React・Vue・Svelte 等、いずれの環境でも利用可能
- **内製 SDK**：外部ライブラリへの依存はゼロ（[`docs/SPEC.md`](SPEC.md) ポリシー責務 §4 参照）
- 標準管理ツールも本 SDK を経由して通信する

### 6.4 標準管理ツール方針

Adlaire CI はすぐに使える標準管理ツールを同梱する。

- **実装技術**：HTML / CSS / JavaScript（バニラ）。外部フレームワーク不使用
- **SDK 経由**：バックエンドとの通信はすべて SDK を介する
- **カスタマイズ基盤**：標準管理ツールをベースとしたカスタマイズを前提とした設計とする。上書き・差し替えが容易な構造を仕様条件として固定する

## 7. 機能インベントリ参照

機能インベントリ責務は [`docs/ROADMAP.md`](ROADMAP.md) を正本とする。

[`docs/SPEC.md`](SPEC.md) では、機能をどの方針とポリシーで扱うかだけを定義し、機能一覧本文を重複定義しない。

## 8. ロードマップ参照

ロードマップ責務は [`docs/ROADMAP.md`](ROADMAP.md) を正本とする。

[`docs/SPEC.md`](SPEC.md) では、状態分類、Phase 実装単位、昇格時の方針とポリシーだけを定義し、ロードマップ本文を重複定義しない。

---

# ポリシー責務
> 遵守義務のある規則と制約を定める。「何をしなければならないか／してはならないか」に答える。

ポリシー責務の節を参照する場合は、必ず [`docs/SPEC.md`](SPEC.md) ポリシー責務 §番号 の形式で記載する。番号だけ、または `docs/SPEC.md §番号` だけで参照してはならない。

方針責務とポリシー責務は、それぞれ独立した節番号体系を持つ。`§5`、`§6`、`§7` などの番号が重複して見える場合でも、責務名を省略して参照してはならない。責務名を伴わない節番号参照は、仕様判断に使用してはならない。

## 0. 詳細仕様記載ポリシー

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) の各仕様項目は、実装者が実装レベルで迷わない記載にしなければならない。

詳細仕様、状態管理、fixture、索引、入口文書を記載または整理する場合は、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.2a の責務ベース明示的原則に従わなければならない。

方針、ポリシー、実装状態、Phase、将来計画、詳細仕様本文、fixture、expected、fake、検証証跡、文書索引を同じ本文内で重複正本化してはならない。

方針とポリシーは [`docs/SPEC.md`](SPEC.md) だけに記載しなければならない。ただし、生成 HTML のデザイン関係は [`docs/DESIGN.md`](DESIGN.md) だけに記載しなければならない。[`docs/ROADMAP.md`](ROADMAP.md)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、owner component 別の [`docs/details/*.md`](details/)、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、[`README.md`](../README.md)、実装ファイル、fixture、PR 本文へ方針またはポリシーを本文として記載してはならない。

詳細仕様は実装詳細だけを記載しなければならない。ロードマップは状態、実装可否、Phase、将来計画だけを記載しなければならない。索引は文書と実装ファイルの所在だけを記載しなければならない。README は利用入口だけを記載しなければならない。

他責務の内容を説明する必要がある場合は、本文を再定義せず、責務を持つ正本と節番号を参照する。参照先に必要な内容が存在しない場合は、参照元へ補足を書かず、責務を持つ正本側を先に改訂しなければならない。責務を持つ正本側で確定できない内容は未確定として扱い、仕様化済みとして扱ってはならない。

判断に迷う記載、二重に読める記載、例外に見える記載、参照先を持たない責務外説明を記載してはならない。これらの違反が残る場合、仕様整合完了、仕様 PR 完了、実装着手、実装済み判定を行ってはならない。

機能、API、設定、状態ファイル、UI、SDK メソッド、運用手順を仕様化する場合は、該当する以下の項目を必ず明記する。

- 対象コンポーネントと責務境界
- 入力値、出力値、戻り値、HTTP メソッド、ステータスコード
- 設定キー、既定値、型、許容範囲、保存先
- ファイルパス、ファイル形式、JSON スキーマ、更新タイミング
- 処理フロー、分岐条件、アルゴリズム、疑似コード
- エラー種別、エラーメッセージ、ログレベル、通知条件
- 再試行、ロック、排他制御、冪等性、タイムアウト
- 認証、認可、トークン、秘密情報、アクセス制御
- 構文確認、実行確認、API 確認、生成物確認などの検証条件

「適切に処理する」「必要に応じて対応する」「安全に扱う」など、実装判断を実装者へ委ねる表現を単独で完了仕様として扱ってはならない。使用する場合は、具体的な条件、処理、値、禁止事項、確認方法を併記する。

## 0a. 仕様成熟度ポリシー

仕様項目の状態は、以下の定義に従って管理する。

| 状態 | 定義 | 実装可否 |
|------|------|----------|
| 未仕様化 | 要求、目的、責務、入出力、処理、状態、検証条件のいずれかが実装判断に必要な粒度で定義されていない状態。 | 実装不可 |
| 将来計画 | 将来的な方向性または候補として記録した状態。実装時期、整理順序、詳細仕様は未確定でもよい。整理順序は実装単位、PR 単位、完了判定単位ではない。 | 実装不可 |
| 改訂予定 | 将来計画または未仕様化の項目を仕様化対象へ昇格した状態。詳細仕様の作成・改訂作業中であり、実装条件はまだ満たしていない。 | 実装不可 |
| 仕様化済み・未実装 | 対象コンポーネント、入出力、設定値、状態、処理順序、異常系、セキュリティ制約、検証条件が実装可能な粒度で定義済みだが、実装ファイルまたは実装コードが未作成・未反映の状態。 | 実装可 |
| 実装中・検証未完了 | 仕様に基づくコード変更へ着手済みだが、必須検証、証跡、または関連文書の整合確認が未完了の状態。 | 検証待ち |
| 実装済み | 仕様に基づくコード変更が完了し、仕様との差分確認、必要な構文確認、実行確認または生成物確認、関連文書の整合確認が完了した状態。 | 完了済み |

`仕様化済み・未実装` へ昇格するには、少なくとも以下を満たさなければならない。

- 対象コンポーネントと責務境界が明確である
- 入力、出力、設定値、状態ファイル、データ構造が明確である
- 正常系の処理順序と分岐条件が明確である
- 異常系、ログ、通知、再試行、ロック、タイムアウトの扱いが必要範囲で明確である
- 認証、認可、秘密情報、外部公開可否などのセキュリティ制約が明確である
- 実装完了を判定する検証条件が明確である

実装着手は、対象項目が `仕様化済み・未実装` の状態に到達している場合に限る。

実装完了は、コード変更だけでは成立しない。仕様との差分確認、構文確認、実行確認または生成物確認、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)・[`docs/ROADMAP.md`](ROADMAP.md) の実装状態・関連仕様の更新要否確認を完了した場合にのみ `実装済み` と扱う。

API、SDK、標準管理ツールのいずれかを変更する場合は、API 仕様、SDK メソッド、UI 操作、詳細仕様の整合を同時に確認する。いずれか一方だけを変更して完了扱いにしてはならない。

## 0b. 仕様 PR 完了ポリシー

仕様策定または仕様改訂の Pull Request は、以下を満たすまで完了扱いにしてはならない。

| 対象 | 完了条件 |
|------|----------|
| 方針責務・ポリシー責務 | 方針、ポリシー、実装着手ゲート、完了判定が [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務に明記されている。 |
| 状態・計画責務 | 実装状態、実装可否、Phase、将来計画、昇格条件が [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務に明記されている。 |
| 詳細仕様入口責務 | 対象機能の読み順、対応表、参照先、受け入れ条件入口が [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務に明記されている。 |
| owner component 別詳細本文責務 | 実装に必要な具体値、入出力、状態、処理順序、異常系、検証条件が該当する [`docs/details/*.md`](details/) 詳細本文責務に明記されている。 |
| 横断契約 | API、SDK、UI、状態ファイル、認証、セットアップの対応関係が該当する詳細本文責務で同期している。 |
| 索引責務 | ファイル名、正本参照先、実装対象の変更がある場合、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務の更新要否を確認している。 |

仕様 PR は、未確定事項を「推奨」「検討」「適切に」等の表現だけで残してはならない。未確定事項を残す場合は、実装不可の `未仕様化` または `将来計画` として明示する。

## 0c. 仕様責務分割ポリシー

仕様変更は、変更責務と統合順序が明確な単位で作成する。競合防止を優先し、同一目的、同一仕様領域、同一ファイル群の仕様変更は必ず 1 本の変更単位にまとめる。

本節は責務ベース明示的原則に基づく必須ポリシーである。仕様変更は同一責務単位で集約し、競合状態、未確認状態、重複状態を完了扱いにしてはならない。

| 分類 | ルール |
|------|--------|
| 同一ファイル変更 | 同じ仕様領域で同一ファイルを編集する変更は、必ず 1 本の PR にまとめる。 |
| API / SDK / UI | endpoint、SDK method、UI 操作のいずれかを変更する場合、対応する仕様を同一 PR に含める。 |
| 状態ファイル | 状態ファイル schema、read/write 対応、破損時処理、権限、atomic write を同一 PR に含める。 |
| セットアップ | systemd、配置パス、権限、初期化、アップデート、rollback、受け入れ条件を同一 PR に含める。 |
| 将来計画 | 実装不可の将来計画は、実装可能仕様と同一表で混在させる場合でも状態を明示する。 |
| 別変更許可 | 変更対象ファイル、責務、統合順序が完全に独立し、片方だけ反映されても仕様矛盾、参照切れ、状態不一致、未定義の依存関係が起きない場合のみ別変更を許可する。 |
| 既存変更優先 | 既存 open PR と同じ仕様領域または同じファイルを変更する場合、新規 PR を作成せず既存 PR へ統合する。 |
| 競合変更 | `DIRTY`、同一ファイル編集、同一仕様領域、統合順序依存の PR が存在する場合、先に 1 本の PR へ統合し、重複 PR を close する。 |

同一目的の仕様変更を複数 PR に分けてはならない。複数 PR 間で同一ファイルまたは同一仕様領域を編集している場合は、最新の作業ブランチへ統合し、1 本の PR にまとめる。

仕様 PR 作成前には、open PR 一覧、変更ファイル一覧、merge 状態を確認する。[`docs/SPEC.md`](SPEC.md) は、競合状態、未確認状態、重複状態を完了扱いにしてはならないことを定義する。具体的な確認コマンド、Git 操作順序、PR 操作手順は [`AGENTS.md`](../AGENTS.md) の Git 運用ルールを正本とする。

競合解消後は、競合マーカーが残っていないこと、`git diff --check` が成功すること、open PR が同一仕様領域で重複していないこと、統合先 PR の merge 状態が `CLEAN` であることを確認する。`UNKNOWN` は merge 可能として扱わない。

## 0d. 仕様凍結ポリシー

実装着手可能な仕様として扱うには、対象仕様を一時的に凍結する。

仕様凍結は、以下をすべて満たす状態をいう。

| 判定項目 | 条件 |
|----------|------|
| 成熟度 | 対象項目が `仕様化済み・未実装` である。 |
| 詳細仕様 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務に対象機能の詳細本文参照と受け入れ条件が明記され、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務に入力、出力、状態、処理順序、異常系、検証条件が明記されている。 |
| テンプレート | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h の機能仕様テンプレートの必須項目を満たしている。 |
| 対応表 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1〜§0i.4 の詳細節対応表に対象機能が記載され、参照先の詳細仕様節と受け入れ条件が一致している。 |
| 横断整合 | API、SDK、UI、状態ファイル、セットアップ、受け入れ条件が矛盾していない。 |
| 未確定事項 | 実装判断に必要な未確定事項が残っていない。 |
| 変更境界 | 実装 PR で変更してよい範囲と変更してはならない範囲が明確である。 |

仕様凍結後、実装中に仕様不足を発見した場合は、実装 PR 内で独自判断による補完を行わず、仕様改訂 PR または同一 PR 内の仕様改訂コミットで凍結状態を更新する。

仕様凍結は永久固定ではない。変更する場合は、凍結解除理由、変更対象、影響範囲、再検証条件を PR 本文に明記する。

## 0e. 初期実装スコープ確定ポリシー

Go 版初期実装では、実装対象を [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務で `仕様化済み・未実装` または `実装済み` と判定し、かつ [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務の対応表と該当 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務に具体的な実装詳細が存在する範囲に限定する。

初期実装の対象範囲は、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務の Phase 実装計画と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務の対応表を正本とする。

[`docs/SPEC.md`](SPEC.md) は、初期実装スコープの判定方針だけを定義する。個別 component の対象可否、対象 Phase、詳細節番号、将来計画該当有無は重複定義しない。

初期実装 PR では、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務で対象外、将来計画、未仕様化、改訂予定に分類された項目を実装してはならない。

初期実装中に対象範囲へ追加したい機能を発見した場合は、先に [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務、該当 owner component の [`docs/details/*.md`](details/) 詳細本文責務、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務を更新し、仕様凍結を再実施する。

初期実装スコープの実装順序と PR 分割は、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §4 の Phase 実装計画に従う。

初期実装スコープの完了判定は、対象コンポーネントごとに [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §4 の対象 Phase 完了条件、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0f の仕様策定完了チェック、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h の機能仕様テンプレート、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1〜§0i.4 の詳細節対応表の受け入れ条件を満たしていることを条件とする。チェック、Phase 完了条件、テンプレート、対応表のいずれかを満たさない対象は、実装済みとして扱ってはならない。

## 0f. Phase 実装単位ポリシー

実装順序、実装計画、実装 PR、完了判定は Phase 単位で行わなければならない。

`P0`、`P1`、`P2〜P5` などの優先度ラベル、抽象段階、API 内部分類、fixture 分類を、実装単位、PR 単位、完了判定単位として使ってはならない。

Phase は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §4 に定義された対象、実装範囲、依存条件、完了条件、検証条件に従わなければならない。

Phase の途中で未仕様化、将来計画、改訂予定の機能を追加してはならない。追加する場合は、先に状態分類、詳細仕様、検証条件、受け入れ条件を更新し、仕様凍結を再実施しなければならない。

API の内部説明や fixture 名に既存の段階名が残る場合でも、それらは検証分類としてのみ扱い、実装順序、実装 PR、完了判定の正本にしてはならない。

## 1. バージョン管理

正式リリース前は、ヘッダーの `仕様バージョン: V.N` と `リリースバージョン: V.X.N` を暫定表記として扱う。暫定表記は、バージョン体系そのものを示す placeholder であり、実在する release tag、GitHub Release、実装済みバージョンを意味しない。

実数運用を開始する場合は、同一 PR で以下をすべて実施する。

1. ヘッダーの `V.N` と `V.X.N` を具体値へ置換する。
2. 具体値に対応する tag / GitHub Release の作成条件を満たしているか確認する。
3. [`README.md`](../README.md)、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) にバージョン表記がある場合は整合させる。
4. 実数運用開始後は placeholder 表記へ戻さない。

### 仕様バージョン V.N

[`docs/SPEC.md`](SPEC.md) の仕様バージョンを管理する。

| 項目 | 内容 |
|------|------|
| 形式 | `V.N`（`V` は固定、`N` は正の整数） |
| 更新方針 | [`docs/SPEC.md`](SPEC.md) の変更・追記のたびに `N` を 1 以上インクリメントする |
| リセット禁止 | `N` はリセット禁止。`V.1` に戻してはならない |
| 例 | `V.205` → `V.206` → `V.207` |

### リリースバージョン V.X.N

開発・ビルド・安定版リリースのバージョンを管理する。`N` は開発・ビルドのたびに、`X` は安定版リリースのたびにインクリメントする。

| 項目 | 内容 |
|------|------|
| 形式 | `V.X.N`（`V` は固定、`X` は安定版リリースの正の整数、`N` は開発・ビルド番号） |
| `N` 更新方針 | 開発・ビルドのたびに 1 以上インクリメントする |
| `X` 更新方針 | 安定版リリースのたびに 1 以上インクリメントする |
| リセット禁止 | `X` と `N` はリセット禁止。`V.1` に戻してはならない |
| 例 | `V.1.100` → `V.1.101` → `V.2.102`（X は安定版リリースのたびに、N はビルドのたびにインクリメント） |

### GitHub リリースポリシー

| 項目 | ルール |
|------|--------|
| リリース作成条件 | 安定版リリース（`X` インクリメント時）のみ GitHub Release を作成する。開発・ビルド（`N` インクリメントのみ）では作成しない |
| タグ形式 | `V.X.N`（リリースバージョンと一致させる）例：`V.2.102` |
| リリースタイトル | タグ名と同一にする |
| リリース形式 | バイナリ配布を標準とする。利用者は GitHub Release から OS/arch 別の実行バイナリを取得し、ソースからのビルドを標準導入手順に含めない |
| 標準 OS/arch | 初期標準は Linux x86_64（`linux-amd64`）とする。追加 OS/arch は将来のリリース対象として個別に仕様化する |
| 添付ファイル | `adlaire-ci-build`、`adlaire-ci-runner`、管理 API 導入後は `adlaire-ci-api` の実行バイナリを添付する。将来 `adlaire-ci-cli` を実装した場合のみ CLI バイナリを追加する。静的 Web サイト生成物を release archive として添付しない |
| 管理 UI 配布物 | `admin/adlaire-ci-sdk.js` と `admin/index.html` はバイナリではなく管理 UI 配布物として扱い、管理 API 導入後の release archive に含める |
| checksum | Release 添付ファイルごとに SHA-256 checksum を提供する。セットアップ手順では配置前に checksum を検証する |
| プレリリースフラグ | 安定版リリースでは `Pre-release` にチェックを入れない |
| ドラフト公開禁止 | Draft Release のまま公開しない |

### 仕様変更とバージョン更新条件

| 変更種別 | 仕様バージョン | リリースバージョン | 備考 |
|----------|----------------|--------------------|------|
| 方針責務・ポリシー責務の変更 | `V.N` を 1 以上更新 | 変更しない | 実装物が変わらない仕様策定のみの変更。 |
| 詳細仕様の実装契約変更 | `V.N` を 1 以上更新 | 変更しない | 実装着手条件や受け入れ条件が変わる。 |
| 仕様影響ありの実装コード変更 | `V.N` を 1 以上更新 | `N` を 1 以上更新 | 実装に合わせて仕様を変更する場合。仕様 PR と実装 PR の両方の完了条件を満たす。 |
| 仕様準拠のみの実装コード変更 | 変更しない | `N` を 1 以上更新 | 既存仕様に完全準拠する実装のみ。仕様文書の変更は不要だが、検証結果は実装 PR に記録する。 |
| 実装内部のみの修正 | 変更しない | `N` を 1 以上更新 | 仕様上の入出力、状態、API、UI、運用手順に影響しない内部修正。 |
| 安定版リリース | 未反映の仕様変更がある場合のみ `V.N` を更新 | `X` と `N` を 1 以上更新 | GitHub Release と tag を作成する。 |
| 将来計画の追記 | `V.N` を 1 以上更新 | 変更しない | 実装可能仕様として扱わない。 |

仕様策定 PR で実装コードを変更しない場合、リリースバージョンを更新してはならない。実装 PR と仕様 PR を同一 PR にまとめる場合は、仕様変更と実装変更の両方の完了条件を満たす。

## 3. カスタマイズ可能範囲

以下は、設定変更時に改訂すべき責務正本を示す参照表である。具体値、既定値、処理本文は owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を正本とする。

| 設定面 | 責務参照 | 変更時の扱い |
|---------|------|---------|
| 入出力 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務 | CLI、設定値、入出力 path、優先順位を同時に整合する。 |
| 将来拡張 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 | 未実装の設定拡張は状態分類を先に確定する。 |

## 4. 外部ライブラリ・フレームワーク方針

### 基本原則

[`docs/SPEC.md`](SPEC.md) 方針責務 §4.1 のゼロ依存・フルインハウス原則を正本とする。開発言語の**標準ライブラリのみ**を採用し、主要機能は本リポジトリ内の仕様と内製実装で完結させる。

外部依存を追加しなければ実装できない機能は、仕様不足または設計不備として扱う。実装者は外部依存の追加で不足仕様を補完してはならない。

### 外部フレームワーク

**いかなる条件でも禁止する。** 例外なし。

### 外部ライブラリ

原則禁止とする。例外採用は、以下の条件をすべて満たす場合に限る。

- 内製化が技術的に困難であり、標準ライブラリだけでは安全性または正確性を担保できない
- 採用範囲が単一責務に限定され、コンポーネント全体の自律性を壊さない
- 採用理由、代替困難性、責務範囲、削除方針、検証条件を [`docs/SPEC.md`](SPEC.md) ポリシー責務に明記している
- [`docs/SPEC.md`](SPEC.md) ポリシー責務 §4 の許可外部ライブラリ一覧へ登録している

開発コスト短縮、実装の容易さ、流行、一般的なベストプラクティスだけを理由にした採用は認めない。許可リスト外のライブラリ使用は認めない。

### 内製ライブラリ・フレームワーク

内製化したライブラリ・フレームワークは、責務、入力、出力、検証条件が [`docs/SPEC.md`](SPEC.md) ポリシー責務に記載されている場合に限り採用できる。内製であっても、仕様未記載の共通基盤を暗黙に追加してはならない。

横断的な内製共通処理を追加する場合は、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.2 に従い、Core ではなく同格の共通責務コンポーネントとして仕様化する。`core`、`common`、`base`、`foundation`、`utils` の名称を使う package、directory、binary、component を作成してはならない。

### 内製実装管理ポリシー

内製スクリプト、標準実装ファイル、実装状態、将来追加予定 path の一覧は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務と [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務を正本とする。

本節では、内製実装を採用するポリシーだけを定義し、個別ファイルの状態表を重複定義しない。新規スクリプトを追加する場合は、先に [`docs/ROADMAP.md`](ROADMAP.md) の状態分類、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) の実装ファイル索引、該当する owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を整合させる。
> 内製 Go コンポーネントは Go 標準ライブラリを基本とする。外部依存は許可リスト登録を必須とする。
> 仕様化済み・未実装、実装中・検証未完了、または将来計画のスクリプトは、実装ファイル、詳細仕様、検証結果、関連文書が実装済みとして整合するまで実装済みとして扱わない。

### 許可外部ライブラリ一覧

内製スクリプト・ライブラリは許可外部ライブラリ一覧に記載しない。内製の管理は内製実装管理ポリシーで行う。

| ライブラリ | 用途 | 許可理由 |
|-----------|------|---------|
| （なし） | — | — |

> 許可外部ライブラリは存在しない。

---

## 5. CI ランナー秘密情報・公開境界ポリシー

- GitHub PAT（Personal Access Token）はスクリプト内にハードコードしてはならない
- PAT は最小権限とし、書き込み権限を必要とする設計を標準としてはならない
- PAT の保存先、permission、読み取り方法、失敗時の扱いは [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/security.md`](details/security.md) 詳細本文責務を正本とする
- ランナーは外部公開エンドポイントを持たない。サーバーから GitHub API への送信のみで動作する

## 6. CI ランナー実行境界ポリシー

- 変更がない場合はビルドをスキップしなければならない
- `adlaire-ci-runner` は oneshot 実行とし、多重実行を防止しなければならない
- 成功、失敗、再試行、状態更新、ログ記録の具体条件は [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/statefile.md`](details/statefile.md) 詳細本文責務を正本とする

## 7. CI ランナー branch target ポリシー

- ランナーは 1 件以上のブランチターゲットを扱える方針とする
- ブランチターゲットの fields、既定値、処理順序、並列可否、検出間隔は [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/statefile.md`](details/statefile.md) 詳細本文責務を正本とする
- ブランチターゲットを変更する場合は、[`docs/details/runner.md`](details/runner.md) 詳細本文責務、[`docs/details/statefile.md`](details/statefile.md) 詳細本文責務、[`docs/details/setup.md`](details/setup.md) 詳細本文責務を同時に整合しなければならない

## 8. SDK 通信契約ポリシー

この節は、`admin/adlaire-ci-sdk.js` に適用する。

- SDK は内製とし、外部ライブラリに依存しない（[`docs/SPEC.md`](SPEC.md) ポリシー責務 §4 参照）
- SDK の対応言語追加は [`docs/SPEC.md`](SPEC.md) ポリシー責務への記載を先行させる
- バックエンド API の変更は SDK の更新を伴う
- SDK の module 形式、公開 API、error class、timeout、streaming 契約は [`docs/details/sdk.md`](details/sdk.md) 詳細本文責務を正本とする
- SDK は自動 retry、戻り値補完、token 永続化、global 代入を行ってはならない

## 9. 標準管理ツール UI / SDK 境界ポリシー

この節は、`admin/index.html` に適用する。

- バニラ HTML / CSS / JavaScript のみで実装する。外部フレームワーク・外部ライブラリは使用しない（[`docs/SPEC.md`](SPEC.md) ポリシー責務 §4 参照）
- バックエンドとの通信はすべて SDK 経由とする。SDK を迂回した直接 API 呼び出しは行わない
- カスタマイズを妨げる密結合な実装を禁止する。SDK 境界、DOM 境界、状態管理境界を満たさない UI 実装は完了扱いにしてはならない
- DOM、form、初期ロード順、イベント処理順、成功/失敗表示、秘密情報消去条件は [`docs/details/ui.md`](details/ui.md) 詳細本文責務を正本とする
- UI は `localStorage`、`sessionStorage`、Cookie から token を復元してはならない
- UI は API transport を直接生成してはならない。禁止対象の具体 API、例外条件、検証条件は [`docs/details/ui.md`](details/ui.md) 詳細本文責務を正本とする

## 10. 状態ファイル永続化ポリシー

- **初期方針：データベース不使用。** 状態はファイルで管理する。具体的な状態ファイル一覧、schema、権限、更新順序、破損時処理は [`docs/details/statefile.md`](details/statefile.md) 詳細本文責務を正本とする
- RDBMS・NoSQL・組み込み DB（SQLite 等）を問わず、初期仕様ではいかなるデータベースも採用しない
- 将来的にデータベースを採用する場合は、[`docs/SPEC.md`](SPEC.md) ポリシー責務への仕様追記と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §4 の許可外部ライブラリ一覧の更新を先行させる
- **データ形式：JSON 形式**を標準とする
- JSON schema、ファイル分割、ネスト制約、encoding、atomic write は [`docs/details/statefile.md`](details/statefile.md) 詳細本文責務を正本とする

## 11. シングルユーザー認証ポリシー

この節は、管理 API サーバーおよび標準管理ツールに適用する。

- 初期構成はシングルユーザーとし、マルチユーザーは将来計画として扱う
- 初期認証、初回変更、強制変更、session、保存形式の具体条件は [`docs/details/api.md`](details/api.md) 詳細本文責務と [`docs/details/security.md`](details/security.md) 詳細本文責務を正本とする
- パスワードは平文保存禁止とし、保存が必要な認証情報はハッシュ化または secret として扱う

## 12. 管理 API 公開境界・session ポリシー

この節は、`components/api.go` に適用する。

- 管理 API サーバーを外部へ直接公開してはならない
- TLS 終端、listen host、session token、認証除外 endpoint の具体条件は [`docs/details/api.md`](details/api.md) 詳細本文責務と [`docs/details/security.md`](details/security.md) 詳細本文責務を正本とする
- セッショントークンは永続化してはならない
- API エンドポイントは、詳細本文責務で明示的に認証除外されたものを除き、認証必須とする
