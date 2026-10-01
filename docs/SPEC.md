# Adlaire CI — 仕様ドキュメント

**標準実装 artifact：** [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3) のディレクトリ構成 tree を参照
**実装 artifact / 機能現在状態：** [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を参照
**Builder 標準出力形式：** 静的 Web サイト（HTML / CSS / JavaScript / search index）
**仕様世代：** Go 初期仕様
**仕様バージョン：** V.N（正式リリース前の暫定表記）/ **リリースバージョン：** V.X.N（正式リリース前の暫定表記）。[`docs/SPEC.md` ポリシー責務 §1](SPEC.md#policy-versioning) 参照。
**更新履歴：** 日付本文を正本化しない。仕様変更の時系列は Git 履歴と Pull Request を正とする。

---

> **Adlaire CI** とは、最初から Go を前提として仕様策定するビルド・CI・管理システムの総称である。

<a id="spec-document-responsibility"></a>

## 文書責務

[`AGENTS.md`](../AGENTS.md) と [`docs/SPEC.md`](SPEC.md) は、本リポジトリにおける最上位文書である。

[`AGENTS.md`](../AGENTS.md) は作業ルールの最上位正本、[`docs/SPEC.md`](SPEC.md) は仕様、方針、ポリシー、禁止事項、完了扱い、実装着手可否の最上位正本である。各文書の責務範囲と矛盾時の正本は、本文書の「責務文書構成」だけで定義する。

[`docs/SPEC.md`](SPEC.md) は「なぜ必要か」と「何を禁止するか」を定義し、[`AGENTS.md`](../AGENTS.md) は承認、実行コマンド、Git 操作、Pull Request 操作、検証手順を定義する。仕様判断と作業手順を同一の正本で定義してはならない。

方針・ポリシーの配置、生成静的 Web サイトと標準管理 UI のデザイン関係の例外、詳細本文との分離、`技術方針` 表と `ディレクトリ構成` tree の削除禁止は、[`docs/SPEC.md` 方針責務 §4.2a](SPEC.md#sec-4-2a) の責務ベース明示的原則を唯一の本文とする。[`docs/SPEC.md`](SPEC.md) 方針責務の入口では同じ禁止事項を再定義しない。

[`docs/SPEC.md`](SPEC.md) 内の番号付き節を参照する場合は、表示文言に `方針責務` または `ポリシー責務` と節番号を必ず併記する。両責務は独立した節番号体系を持つため、責務名を伴わない節番号、または `docs/SPEC.md §番号` だけの参照を仕様判断に使用してはならない。番号を持たない前置き節のうち、[文書責務](#spec-document-responsibility)、[状態参照方針](#spec-state-reference-policy)、[責務文書構成](#document-responsibility-map) は、表示文言に対象責務名を併記し、ここで指定した固定 anchor へリンクする場合に限り参照できる。この例外を番号付き節、その他の見出し、または自動生成 anchor だけの参照へ拡張してはならない。

<a id="spec-state-reference-policy"></a>

## 状態参照方針

状態語彙、状態定義、実装可否、状態遷移条件は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity)、実装 artifact と各機能の現在状態、Phase、依存順序、将来計画は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を唯一の正本とする。

---

## 責務文書構成

<a id="document-responsibility-map"></a>
| 責務 | 正本 | 責務の問い | 記載する内容 |
|------|------|-----------|------------|
| 方針責務 | [`docs/SPEC.md`](SPEC.md) | **なぜ・何を** | 目的、設計思想、方向性の原則。変更頻度が低く、判断の拠り所となる指針。 |
| ポリシー責務 | [`docs/SPEC.md`](SPEC.md) | **しなければならない／してはならない** | 遵守義務のある規則、制約、禁止事項、セキュリティ要件、運用ルール、バージョン管理規則。 |
| 状態・計画責務 | [`docs/ROADMAP.md`](ROADMAP.md) | **現在どの状態か・いつ・どれを** | 実装 artifact と各機能の現在状態、Phase、機能インベントリ、将来計画。状態語彙と遷移条件は再定義しない。 |
| 詳細仕様入口責務 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) | **どこから読むか** | 詳細仕様参照入口、共通固定値、owner 対応表、collaborator 境界参照入口。 |
| owner component 別詳細本文責務 | [`docs/details/*.md`](details/)。ただし [`docs/details/fixture.md`](details/fixture.md) を除く。 | **どのように実装するか** | owner component 別の入出力、状態、処理順序、異常系、検証条件。 |
| fixture 証跡責務 | [`docs/details/fixture.md`](details/fixture.md) | **何で検証するか** | fixture、expected、fake、実装検証証跡、acceptance checklist、差し戻し条件。 |
| デザイン責務 | [`docs/DESIGN.md`](DESIGN.md) | **どう見せるか** | 生成静的 Web サイトと標準管理 UI のデザイン関係。 |
| 文書・実装ファイル所在の索引責務 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) | **どこにあるか** | 文書、実装ファイル、testdata、未作成 path の所在。 |
| 利用入口責務 | [`README.md`](../README.md) | **どう始めるか** | 利用者向け入口、概要、参照先。 |

新しい記載内容は「この内容はどの責務の問いに答えるか」を基準に、責務を持つ正本を決定する。

新しい判断対象の正本は、[責務文書構成表](#document-responsibility-map) の責務の問いと記載内容だけで決定する。表の分担に一致しない他文書を、例外の正本としてはならない。

実装者は、実装前に以下の参照順序方針に従う。

1. [`docs/SPEC.md`](SPEC.md) の「文書責務」と「状態参照方針」で、正本範囲を確認する。
2. [`docs/ROADMAP.md`](ROADMAP.md) で対象の現在状態、Phase、将来計画該当有無を確認し、実装可否は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity) で判定する。
3. [`docs/SPEC.md` 方針責務 §4.1](SPEC.md#sec-4-1)〜[`docs/SPEC.md` 方針責務 §4.10](SPEC.md#sec-4-10) で、ゼロ依存、責務ベース明示的原則、ディレクトリ構成、完全仕様詳細化、成熟度、着手ゲート、完了判定、Go 正本方針を確認する。
4. [`docs/SPEC.md`](SPEC.md) のポリシー責務で、対象領域の禁止事項、セキュリティ、バージョン、外部依存を確認する。
5. 実装または検証を扱う場合は、[`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test) の意味のあるテストポリシーを確認する。
6. 生成静的 Web サイトまたは標準管理 UI のデザイン関係を扱う場合は、[`docs/DESIGN.md`](DESIGN.md) デザイン責務で視覚仕様を確認する。
7. [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務で詳細仕様参照入口、共通固定値、owner 対応表、collaborator 境界参照入口を確認し、該当する owner component 別の [`docs/details/*.md`](details/) 詳細本文責務で実装に必要な入出力、状態、異常系、検証条件を確認する。Phase 11 を扱う場合も専用詳細ファイルを作らず、状態は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry)、fixture は [`docs/details/fixture.md` fixture 証跡責務 Phase 11 fixture harness 参照](details/fixture.md#phase-11-fixture-harness-reference)、実装契約は該当 owner 詳細本文責務を確認する。文書と実装ファイルの実在所在は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務を確認する。

---

# 方針責務
> 目的・設計思想・方向性の原則を定める。「なぜこう作るか」に答える。

## 1. 目的

Adlaire CI は、Markdown からの静的 Web サイト生成、source 変更検出と build / deploy、管理 interface、setup / update、バイナリ release を、Go を中心とする自己管理型システムとして一貫して提供する。

個別機能の実装契約は owner component 別の [`docs/details/*.md`](details/) 詳細本文責務、実装 artifact と機能の現在状態は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、利用者向け release asset の形式は [`docs/SPEC.md` ポリシー責務 §1](SPEC.md#policy-versioning) を正本とする。

## 2. 開発方針

本プロジェクトは **仕様駆動開発（Spec-Driven Development）** を採用する。

- **責務正本群が判断基準**：実装の追加・変更は、変更対象の責務を持つ正本への反映を先行させる
- **仕様から実装への順序**：実装は、[`docs/SPEC.md`](SPEC.md)、[`docs/ROADMAP.md`](ROADMAP.md)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務の内容に基づいて修正する
- **責務正本との乖離は不整合**：乖離が生じた場合も、責務を持つ正本を先に改訂し、その正本に基づいて実装を修正する

<a id="direction-technical"></a>

## 4. 技術方針

以下の技術方針表は Adlaire CI の技術選定を定義する。詳細仕様はこの表に従い、具体的な入出力、処理順序、例外条件、検証条件を owner component 別の [`docs/details/*.md`](details/) 詳細本文責務へ記載する。

| 領域 | 方針 |
|---|---|
| ランタイム | 正式公開済みの Go stable toolchain。最低バージョンは [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](DETAIL_INDEX.md#0d-共通固定値) の共通固定値とし、beta、release candidate、開発 snapshot は使用しない。 |
| 言語 | Go（ビルド、ランナー、管理 API、セットアップ、リリース）/ JavaScript（SDK）/ HTML・CSS・Vanilla JavaScript（UI） |
| HTTP | Go 標準ライブラリ `net/http` |
| データベース | なし（ファイルベース） |
| Git・GitHub 参照 | `runner` の source 取得は GitHub REST API の Blobs / Trees API を Go 標準ライブラリ `net/http` 経由で行う。`release` の local commit / tag / source 参照は read-only Git command、remote repository / branch / tag / Release 参照は GitHub REST API を使用する。実行時の Git 状態変更、ref 作成・更新・削除、fetch、checkout、merge、commit、push は行わない。 |
| フロントエンド | HTML / CSS / Vanilla JavaScript |
| JavaScript 検証ランタイム | JavaScript 系実装 artifact の構文・型検証は Deno stable runtime を標準とする。Node.js、npm、bundler、transpiler を標準検証ランタイムとして使用しない。Deno は検証ランタイムであり、Adlaire CI の production dependency、配布物同梱 dependency、外部ライブラリ許可条件として扱ってはならない。 |
| 標準運用 | systemd を使用する自己管理 Linux CI サーバーでビルドし、別ホストの静的コンテンツ配信サーバーへ SSH で転送する 2 サーバー構成 |
| データ交換形式 | 構造化 API request / response、Adlaire CI が読み書きまたは解釈する状態ファイル、設定ファイル、build 拡張設定は JSON に統一し、YAML・CSV・XML 等の代替構造化形式を使用しない。この禁止は Adlaire CI の runtime / API / state / build 定義に適用し、GitHub Actions / GitHub workflow の `.github/workflows/*.yml` / `.github/workflows/*.yaml` は対象外とする。例外は、保存済み build log の有限 SSE response、保存済み snapshot の binary response、管理ツールの静的 asset response だけとする。例外となる endpoint、media type、body、失敗時 JSON response は [`docs/details/api.md`](details/api.md) 詳細本文責務、SDK の受信・変換契約は [`docs/details/sdk.md`](details/sdk.md) 詳細本文責務を正本とする。詳細本文責務に明示されていない非 JSON API 交換、状態設定、build 定義を実装してはならない |

<a id="41-ゼロ依存フルインハウス原則"></a>

<a id="sec-4-1"></a>
**4.1 ゼロ依存・フルインハウス原則：**

Adlaire CI は、ゼロ依存・フルインハウスを技術哲学の中核とする。

[`docs/SPEC.md` 方針責務 §4.1](SPEC.md#sec-4-1) のゼロ依存とは、各コンポーネントが外部ライブラリ、外部フレームワーク、外部ビルドツール、外部ホスティング実行基盤に機能成立を依存しないことを意味する。[`docs/SPEC.md` 方針責務 §4.1](SPEC.md#sec-4-1) のフルインハウスとは、Markdown 変換、CI 実行、管理 API、SDK、標準管理ツール、状態管理、認証、ログ、通知、セットアップ、リリースの主要機能を本リポジトリ内で仕様化し、内製コードとして理解、検証、保守できる状態を意味する。

実装者は、便利さ、実装速度、一般的な慣習を理由に外部ライブラリで未定義機能を補完してはならない。外部依存がなければ成立しない設計は設計不備として扱い、先に仕様を見直す。

各コンポーネントの自律性は以下を満たす。

| owner component | 実装内包先 / 配布対象 | 自律性の条件 |
|-----------------|---------------------|--------------|
| `builder` | `components/builder.go` | Go 標準ライブラリだけで Markdown 解析、HTML/CSS/JS/search index 生成、検証レポート出力を行う。外部 Markdown parser、template engine、syntax highlight library、search library に依存しない。 |
| `runner` | `components/runner.go` | Go 標準ライブラリと OS 標準コマンドだけで GitHub API polling、SHA 比較、ビルド起動、ログ、通知、SSH 転送、snapshot、lock、retry を処理する。外部 CI サービス、job queue、scheduler library に依存しない。 |
| `api` | `components/api.go` | Go 標準ライブラリ `net/http` を基本に、認証、session、状態ファイル CRUD、入力検証、API response を内製実装する。外部 web framework、router、ORM、database driver に依存しない。 |
| `setup` | `components/setup.go` | Go 標準ライブラリと明示した OS 標準 command だけで Release asset 取得、checksum 検証、配置、systemd 操作、更新、rollback を実装する。外部 installer framework、package manager、shell script を実装主体にしない。 |
| `release` | `components/release.go` | Go 標準ライブラリ、Go toolchain、Git command、GitHub REST API だけで再現可能な成果物生成、checksum、draft upload、再取得検証、正式公開を実装する。外部 release framework、archive tool、checksum tool に依存しない。 |
| `sdk` | `admin/adlaire-ci-sdk.js` | 単一 ES Module とし、browser 標準 API のみで API client、error handling、streaming、timeout を実装する。npm package、bundler、polyfill、framework に依存しない。 |
| `ui` | `admin/index.html` | HTML / CSS / Vanilla JavaScript だけで標準管理ツールを構成し、SDK 経由で通信する。React、Vue、Svelte、CSS framework、icon package、chart library に依存しない。 |
| `admin` | `admin/index.html`、`admin/adlaire-ci-sdk.js`、`components/admin.go` | 管理 UI 静的配布物の構成、検証、配信境界、および CLI 管理クライアントを本リポジトリ内で完結させる。外部 asset pipeline、CDN、archive framework、CLI framework に依存しない。 |
| `statefile` | `components/runner.go`、`components/api.go` | Go 標準ライブラリだけで状態 schema、lock、atomic write、JSON Lines、破損検出を処理する。外部 database、storage engine、serialization library に依存しない。 |
| `archive` | `components/runner.go`、`components/api.go` | Go 標準ライブラリだけで archive、snapshot、展開、検証、世代管理を処理する。外部 archive tool、snapshot service、object storage SDK に依存しない。 |
| `commitstatus` | `components/runner.go` | Go 標準ライブラリ `net/http` だけで GitHub Commit Status payload、送信、応答検証を処理する。外部 GitHub client library、CI status service に依存しない。 |
| `security` | `components/api.go` | Go 標準ライブラリだけで認証、認可、token、session、TOTP、rate limit、audit、secret 処理を実装する。外部 authentication framework、secret management SDK に依存しない。 |
| `mcp` | `components/mcp.go` | Go 標準ライブラリを前提とし、MCP 通信、JSON-RPC 処理、API bridge、監査ログを内製する。外部 MCP framework に依存する前提で仕様化しない。現在状態は [`docs/ROADMAP.md`](ROADMAP.md)、詳細仕様は [`docs/details/mcp.md`](details/mcp.md)、実在所在は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) を参照する。 |

「実装内包先 / 配布対象」列は、owner component の責務を実装または配布する artifact の所在を示す。同列を、owner component と artifact の同一視、単独専用 artifact の存在保証、現在状態、実装着手許可の根拠として使用してはならない。実装 artifact の実在所在は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務、実装 artifact と機能の現在状態は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を正本とする。

[`docs/SPEC.md` 方針責務 §4.1](SPEC.md#sec-4-1) のコンポーネント自律性表はコンポーネント自律性の方針であり、関数単位の処理、入出力、状態、異常系、検証条件を定義するものではない。各コンポーネントの具体的な実装契約は、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を正本とする。

外部依存の禁止条件、例外採用条件、許可外部ライブラリ一覧、登録要件は [`docs/SPEC.md` ポリシー責務 §4](SPEC.md#policy-dependencies) だけを正本とする。[`docs/SPEC.md` 方針責務 §4.1](SPEC.md#sec-4-1) はコンポーネント自律性だけを定義し、外部依存の採否手順または許可条件を再定義しない。

<a id="42-core-非採用共通責務コンポーネント方針"></a>

<a id="sec-4-2"></a>
**4.2 Core 非採用・共通責務コンポーネント方針：**

Adlaire CI は、`core`、`adlaire-ci-core`、`internal/core`、`common`、`base`、`foundation`、`utils` のような中心化・汎用置き場化する概念を採用しない。

横断的に利用される処理は、単一責務、入出力、状態、異常系、検証条件、依存方向が仕様化されている場合に限り、共通責務コンポーネントとして定義できる。ただし、共通責務コンポーネントは他コンポーネントより上位ではなく、同格の独立コンポーネントとして扱う。共通責務コンポーネントを、正本、中心、基盤、親、上位レイヤーとして扱ってはならない。

共通責務コンポーネントを追加する場合は、以下をすべて満たす。

- 単一責務である。
- [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務に参照入口があり、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務に入力、出力、状態、異常系、検証条件が仕様化されている。
- 利用元コンポーネントとの依存方向が明確である。
- アプリケーション固有の判断、業務判断、UI 判断、API endpoint 判断、ビルド対象判断を含まない。
- 何でも置き場として使用しない。
- `core`、`common`、`base`、`foundation`、`utils` など中心性または汎用置き場を示す名称を使わない。
- 追加時に [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、[`docs/SPEC.md` ポリシー責務 §4](SPEC.md#policy-dependencies)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を整合する。

共通責務コンポーネントの命名候補を仕様化する場合は、以下の形式で責務を明示する。

| 命名候補 | 責務 |
|----|------|
| `state_store` | 状態ファイルの atomic write、lock、JSON 読み書き。 |
| `secret_masker` | PAT、Webhook Secret、SMTP password、session token、API token のマスク処理。 |
| `event_log` | 構造化イベントログの形式、出力、読み取り。 |
| `github_client` | GitHub REST API 呼び出し、rate limit、retry 境界。 |

[`docs/SPEC.md` 方針責務 §4.2](SPEC.md#sec-4-2) の共通責務コンポーネント命名候補表は、共通責務コンポーネントを追加する場合の命名・責務明示ポリシーである。[`docs/SPEC.md` 方針責務 §4.2](SPEC.md#sec-4-2) の共通責務コンポーネント命名候補表に含まれる名称だけを理由に、実装ファイル、package、directory、状態項目、placeholder を作成してはならない。追加可否と状態語彙は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity)、追加後の実装 artifact と機能の現在状態は [`docs/ROADMAP.md`](ROADMAP.md)、所在は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) を正本とする。

禁止される設計は以下とする。

| 禁止対象 | 理由 |
|----------|------|
| `core` / `adlaire-ci-core` / `internal/core` | 中心コンポーネント化し、正本や上位概念と誤読されるため。 |
| `common` / `base` / `foundation` / `utils` | 責務が曖昧になり、何でも置き場化しやすいため。 |
| 全共通処理の一括集約 package | 単一責務を失い、コンポーネント境界を壊すため。 |
| 業務判断を含む共通処理 | 利用元コンポーネントの責務を横断基盤へ漏らすため。 |

共通責務コンポーネントは、コード共有のためだけに追加してはならない。重複削減より、責務境界、仕様の明確さ、依存方向の追跡可能性を優先する。

<a id="42a-責務ベース明示的原則"></a>

<a id="sec-4-2a"></a>
**4.2a 責務ベース明示的原則：**

Adlaire CI の仕様体系は、責務ベース明示的原則を仕様全般の最上位方針として採用する。

[`docs/SPEC.md` 方針責務 §4.2a](SPEC.md#sec-4-2a) の責務ベース明示的原則は [`docs/SPEC.md`](SPEC.md) 全体に適用する。[`docs/SPEC.md`](SPEC.md) 内の各記載は、方針責務、ポリシー責務、状態・計画責務の参照、詳細仕様入口責務の参照、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務の参照、[`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務の参照、文書・実装ファイル所在の索引責務の参照、利用入口責務の参照のいずれかとして読めなければならない。

責務ベース明示的原則とは、方針・ポリシー・状態語彙・状態定義・状態遷移条件、実装 artifact と機能の現在状態・Phase・将来計画、詳細仕様本文、fixture・expected・fake・検証証跡、文書・実装所在を、それぞれ異なる責務として明示的に分離する原則である。

<a id="spec-global-no-duplicate-principle"></a>
**仕様全般重複記載禁止原則：**

仕様上の一つの判断対象は、一つの責務正本だけが本文を持たなければならない。この原則は、[`docs/SPEC.md`](SPEC.md)、[`docs/DESIGN.md`](DESIGN.md)、[`docs/ROADMAP.md`](ROADMAP.md)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、owner component 別の [`docs/details/*.md`](details/)、[`docs/details/fixture.md`](details/fixture.md)、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、[`README.md`](../README.md) の仕様全般に適用する。

同一判断対象について、完全一致する記載、言い換え、要約、部分転載、表または箇条書きへの再掲、具体値、条件、禁止事項、schema、処理順序、異常系、検証条件の再定義を重複記載として禁止する。複数箇所の内容が一致する場合は重複違反、内容が一致しない場合は重複違反かつ仕様矛盾として扱う。

正本以外の文書は、正本の内容を説明、補足、要約、緩和、例外化、再解釈してはならない。責務外の判断対象が必要な場合は、責務名付き Markdown link で正本の該当 anchor を参照し、自分の責務に固有の内容だけを記載する。参照先に必要な仕様が存在しない場合は、参照元へ仮置きまたは補足せず、責務を持つ正本を先に改訂する。責務を持つ正本で確定できない内容は未確定として扱い、仕様化済みとして扱ってはならない。

分かりやすさ、参照性、一覧性、実装者向け補足、または文書単体で理解できることを、重複記載の理由として認めない。これらは責務名、責務正本、正本範囲、参照 anchor を明示した Markdown link によって担保する。

重複記載を発見した場合は、責務ベース明示的原則によって唯一の正本を確定し、正本だけに本文を残す。正本以外の重複本文は削除し、必要な参照だけを責務名付き Markdown link へ置き換える。正本を確定できない状態、重複記載が残る状態、または参照へ置き換えた結果として実装判断に必要な情報が失われる状態を、仕様整合完了として扱ってはならない。

次の記載だけは、同一判断対象の正本本文を持たず、正本の意味を追加または変更しない場合に限り、重複記載として扱わない。

- Markdown link の参照対象を識別するために必要な機能名、責務名、owner component 名、ファイル名、節名、anchor 名。
- 索引責務または詳細仕様入口責務が管理する path、存在区分、owner 対応、参照先。
- collaborator component が自分の責務として持つ接続、入力受け渡し、出力受け渡し、変換、失敗伝播の固有契約。
- fixture 証跡責務が正本へのリンクとともに記録する入力、操作、期待結果、assertion、実装検証証跡。これらは仕様本文の正本ではなく、正本との一致を検証する証跡としてだけ扱う。

正本参照先は、必ず責務名とファイル名で示す。文書を章構成、便宜分類、または他文書の従属章として扱ってはならない。[`docs/SPEC.md`](SPEC.md) 内部の見出しも責務名で示し、`Part` 名称で扱ってはならない。[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、owner component 別の [`docs/details/*.md`](details/)、[`docs/details/fixture.md`](details/fixture.md)、[`docs/ROADMAP.md`](ROADMAP.md)、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、[`README.md`](../README.md) を [`docs/SPEC.md`](SPEC.md) の章として扱ってはならない。

参照はリンク化を必須とする。文書間参照、節参照、表参照、責務正本参照、実装ファイル参照、fixture 参照を説明文として書く場合は、Markdown link を用いて参照先へ移動できる形にする。単なるファイル名、裸の節番号、裸の見出し名、または `参照` という文字だけで参照先を示した扱いにしてはならない。

テスト方針、テスト完了可否、fixture 証跡 schema、assertion、test gap inventory、test improvement batch closure、mutation、mutation selection、race trigger、drift、skip、fixture root closure、PR 証跡への参照は、owner 詳細本文、作業ルール、Pull Request 本文、状態・計画責務、文書・実装ファイル所在の索引責務、利用入口責務、または他の責務文書から行う場合、必ず [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 横断テスト証跡共通入口](DETAIL_INDEX.md#cross-test-evidence-route) を経由する。ただし、[`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test) がテスト方針と完了可否の本文を持つ場合、[`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務が fixture 証跡本文を持つ場合、または [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務が所在索引として固定 anchor を列挙する場合は、それぞれの責務範囲内で直接リンクしてよい。owner 詳細本文、作業ルール、Pull Request 本文、状態・計画責務が同じテスト証跡の schema、必須 key、記録単位、例外条件、完了可否を個別に再掲または再定義することを禁止する。

契約値そのものを記録するコードブロック、ディレクトリ tree、JSON schema、CLI 例、HTTP path、設定値、生成物名、状態ファイル名、および所在索引表の path セルは、Markdown link 化によって契約文字列が変わるためリンク化対象外とする。この例外は説明文中の参照には適用しない。説明文から実ファイルまたは文書へ移動させる目的がある場合は、同じ段落または表の参照列に Markdown link を併記する。

リンク化する参照は、初出または単独参照の表示文言に責務名とファイル名を必ず含め、節または本文ラベルを参照する場合は固定 anchor へリンクする。同じ段落、同じ表セル、または同じ箇条書き内で、直前のリンクと同一ファイル・同一責務を指す連続参照だけは、表示文言を節番号、本文ラベル、または対象名へ短縮できる。短縮参照であっても自動生成 anchor へリンクしてはならない。例として、状態・計画責務を参照する場合は [`docs/ROADMAP.md`](ROADMAP.md)、詳細仕様入口責務を参照する場合は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、owner component 別詳細本文責務を参照する場合は [`docs/details/*.md`](details/) 詳細本文責務、fixture 証跡責務を参照する場合は [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務、文書・実装ファイル所在の索引責務を参照する場合は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、利用入口責務を参照する場合は [`README.md`](../README.md) のように記載する。リンク化できない生成物、PR 本文、外部ツール出力で参照を記録する場合でも、参照先ファイル名、責務名、節番号または固定 anchor 名を省略してはならない。

細目を見出し化してはならない。詳細仕様、ロードマップ、索引、fixture 証跡では、実装項目、正常系、異常系、入力、出力、処理順序、検証条件、固定契約、確認ゲート、表名、schema 名、fixture 名、機能別補足を Markdown 見出しとして作成してはならない。これらは表の行、箇条書き、または本文内ラベルとして記載する。

見出しは文書構造の最低限の区切りに限定し、責務、正本、禁止事項、判断原則、固定契約を見出しで定義してはならない。見出しをルール本文、責務本文、正本本文、参照本文の代替として扱ってはならない。

見出し内に Markdown link を置いてはならない。参照が必要な場合は、見出しではなく直後の本文、表、箇条書きに責務名付き Markdown link を置く。見出し内リンクにより生成される anchor を正本参照先として扱ってはならない。

本文内ラベルをリンク参照先にする必要がある場合は、Markdown 見出しを作成せず、当該ラベル直前に `<a id="..."></a>` 形式の本文アンカーを置く。本文アンカーの id は、同一ファイル内で重複してはならない。本文アンカーは参照先を固定するためだけに使用し、責務、正本、禁止事項、判断原則、固定契約の本文を置き換えてはならない。

方針とポリシーは [`docs/SPEC.md`](SPEC.md) だけに記載する。ただし、生成静的 Web サイトと標準管理 UI のデザイン関係は [`docs/DESIGN.md`](DESIGN.md) だけに記載する。詳細仕様、ロードマップ、索引、README、実装ファイル、fixture、PR 本文は、方針またはポリシーを本文として定義、補足、緩和、例外化、再解釈してはならない。

方針またはポリシーに該当する記載を [`docs/SPEC.md`](SPEC.md) から削除し、他文書への参照だけに置き換えてはならない。方針またはポリシーを整理する場合は、[`docs/SPEC.md`](SPEC.md) 内で責務名、適用範囲、禁止事項、参照先を明示して整える。

[`docs/SPEC.md`](SPEC.md) の `技術方針` 表と `ディレクトリ構成` tree は、方針責務本文である。これらを実装詳細、状態一覧、索引本文とみなして削除、他文書へ移動、または参照だけへ置き換えてはならない。

owner component は対象機能の詳細本文を持つ。collaborator component は、境界、接続、入力受け渡し、出力受け渡し、検証観点として参照される。collaborator component は、owner component の本文を置き換えたり、同じ判断対象を別正本として再定義したりしてはならない。

owner component と実装 artifact は別の判断対象とする。owner component は機能契約の責務境界であり、単独の実装ファイルが存在することを意味しない。実装 artifact は [`main.go`](../main.go)、`components/*.go`、`admin/*`、実行バイナリ、配布物、または運用実行物のように、リポジトリまたはリリースで実在と実行可否を確認できる対象をいう。owner component が持つ機能の現在状態を、同名の実装 artifact が存在するという理由だけで判定してはならない。

[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務は、実装 artifact の現在状態と owner component が持つ機能の現在状態を区別して記録する。単独の専用 Go 実装 artifact を持たない `admin`、`statefile`、`archive`、`commitstatus`、`security` の owner component は、機能インベントリの現在状態で管理する。実装 artifact の path、起動名、入力 interface のいずれかが必要であるにもかかわらず未定義の機能は、実装可能な仕様として扱ってはならない。

[責務文書構成表](#document-responsibility-map) の owner component 別詳細本文責務で [`docs/details/*.md`](details/) を総称として参照する場合でも、[`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務を含めて読んではならない。

仕様全般を整理する場合は、[`docs/SPEC.md` 方針責務 §4.2a 仕様全般重複記載禁止原則](SPEC.md#spec-global-no-duplicate-principle) を、各文書の個別整理規則より上位の方針として適用する。

<a id="43-ディレクトリ構成"></a>

<a id="sec-4-3"></a>
**4.3 ディレクトリ構成：**

Adlaire CI のディレクトリ構成は、責務ベースで整理する。

ディレクトリ構成 tree は以下とする。[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) は所在索引として [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3) に従う。

```text
.
├── main.go
├── main_test.go
│
├── components/
│   ├── builder.go
│   ├── builder_test.go
│   ├── runner.go
│   ├── runner_test.go
│   ├── api.go
│   ├── api_test.go
│   ├── admin.go
│   ├── admin_test.go
│   ├── setup.go
│   ├── setup_test.go
│   ├── release.go
│   ├── release_test.go
│   ├── mcp.go
│   └── mcp_test.go
│
├── admin/
│   ├── index.html
│   └── adlaire-ci-sdk.js
│
├── testdata/
│   ├── builder/
│   │   ├── single/
│   │   │   └── expected/
│   │   ├── site/
│   │   │   └── expected/
│   │   ├── empty-dir/
│   │   │   └── expected/
│   │   ├── strict/
│   │   │   └── expected/
│   │   └── safe/
│   │       └── expected/
│   ├── runner/
│   ├── api/
│   ├── admin/
│   ├── sdk/
│   ├── ui/
│   ├── statefile/
│   ├── archive/
│   ├── commitstatus/
│   ├── security/
│   ├── setup/
│   ├── release/
│   └── mcp/
│
├── docs/
│   ├── SPEC.md
│   ├── ROADMAP.md
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
│   │   ├── release.md
│   │   ├── statefile.md
│   │   ├── security.md
│   │   ├── archive.md
│   │   ├── commitstatus.md
│   │   ├── mcp.md
│   │   └── fixture.md
│   └── examples/
│
├── sdk_contract_test.go
├── ui_contract_test.go
├── README.md
├── AGENTS.md
└── go.mod
```

[`main.go`](../main.go) は 1 ファイルとし、起動入口、実行ファイル名の exact 判定、引数受け取り、binary version 注入値の受け渡し、対象 owner component 呼び出しだけを担当する。[`main.go`](../main.go) に Markdown 変換、CI 実行、HTTP handler、状態ファイル操作、archive 処理、GitHub Commit Status 送信、setup、release、MCP 処理の実装詳細を書いてはならない。未知の実行ファイル名を既定 owner component へ fallback してはならない。

`components/` は、1 標準 Go 実装対象 = 1 Go ファイルとする。この 1 ファイル原則は実装 artifact の配置規則であり、owner component と実装 artifact を同一概念にする規則ではない。Go 実装ファイルの所在は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務、現在状態は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を参照する。`admin` は `admin/` 配下の静的配布物と、CLI 管理クライアント用の `components/admin.go` を所有する責務境界として扱う。`setup` と `release` の標準 Go 実装 artifact はそれぞれ `components/setup.go` と `components/release.go` とする。`statefile`、`archive`、`commitstatus`、`security` は詳細仕様上の責務境界であり、単独 Go ファイルを作成する場合は該当 Phase または追加実装 PR で仕様状態と索引を更新してから追加する。`mcp.go` の現在状態は [`docs/ROADMAP.md`](ROADMAP.md)、実装可否と追加条件は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity) と該当 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を参照する。

`admin/` は標準管理 UI の静的ファイルを配置する。`testdata/` は責務別 fixture を配置する。`docs/examples/` は利用例、設定例、サンプル構成を配置する。

ディレクトリ構成は方針上の到達形を示す。実ファイルの有無と未作成 path は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、実装 artifact と各機能の現在状態、Phase、将来計画への割当は [`docs/ROADMAP.md`](ROADMAP.md) を正本とする。ディレクトリ構成に含まれることだけを理由に、未実装ファイル、将来追加予定 path、空ディレクトリ、placeholder を作成してはならない。

<a id="44-詳細仕様方針"></a>

<a id="sec-4-4"></a>
**4.4 詳細仕様方針：**

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) は、詳細仕様入口責務として、実装者が追加判断なしに該当する詳細本文へ到達できる粒度で記載する。

owner component 別の [`docs/details/*.md`](details/) は、詳細本文責務として、抽象的な方針や目的の再掲ではなく、実装時に必要な具体値、処理順序、入出力、状態、失敗時の扱いを定義する。ただし [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務は owner component 別詳細本文責務に含めない。Phase 専用詳細仕様ファイルを作成して、owner component 別詳細本文責務を置き換えてはならない。

仕様化済み・未実装の項目であっても、実装予定として扱う場合は実装者が迷わない粒度まで詳細化する。実装時期、設計判断、具体値が未確定の内容は、実装可能な仕様として扱わず、未仕様化または将来計画として明示する。

詳細仕様入口責務と詳細本文責務の組み合わせは、[`docs/SPEC.md` ポリシー責務 §0 詳細仕様必須項目](SPEC.md#detail-contract-required-fields) の全項目を追加判断なしで特定できる状態を満たす。[`docs/SPEC.md` 方針責務 §4.4](SPEC.md#sec-4-4) に必須項目を再掲して別の判定表としてはならない。

<a id="complete-detail-specification-policy"></a>
完全仕様詳細化は、全 Phase、全 owner component、全実装対象機能に適用する。完全仕様詳細化とは、実装者が設計判断、仕様補完、例外判断、検証条件の推測、既存実装への後追い合わせを行わずに、責務正本だけから実装、検証、完了判定を進められる状態をいう。

完全仕様詳細化が完了していない機能は、一定の仕様が存在していても実装対象として扱わない。実装しながら仕様を決めること、テスト結果で仕様を後から確定すること、既存実装の挙動を理由に詳細本文の不足を補うことを、実装着手方針として認めない。

<a id="45-仕様成熟度方針"></a>

<a id="sec-4-5"></a>
**4.5 仕様成熟度方針：**

仕様項目は、実装可否を判断できるように成熟度を明確に区分する。

仕様成熟度方針は、成熟度を必要とする理由と判断方針を定める。状態名、状態定義、実装可否、昇格条件は、[`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity) の仕様成熟度ポリシーを正本とする。

成熟度は、構想の有無ではなく、実装者が実装に着手できるだけの情報が揃っているかで判定する。

将来計画や未確定アイデアは、実装対象として扱わない。実装対象にする場合は、先に詳細仕様を整え、仕様化済み・未実装へ昇格させる。

実装中に仕様不足、実装者判断に依存する分岐、未定義の入出力、未定義の状態ファイル、未定義の異常系を発見した場合は、実装判断で補完せず、仕様改訂へ戻す。

<a id="46-phase-実装単位方針"></a>

<a id="sec-4-6"></a>
**4.6 Phase 実装単位方針：**

Phase は、対象 owner component、実装範囲、依存条件、完了条件、検証条件を一体として管理する実装単位である。Phase 単位の必須操作、active Phase、優先度ラベルの使用禁止、途中追加禁止は [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit) を唯一の正本とする。Phase の一覧、現在状態、順序、依存関係は [`docs/ROADMAP.md` 状態・計画責務 §4](ROADMAP.md#roadmap-phase-plan)、完了判定方針は [`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8)、個別の実装契約は owner component 別の [`docs/details/*.md`](details/) 詳細本文責務、fixture の一般形式と実装検証証跡は [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務を参照する。Phase 専用の詳細仕様ファイルを作成してはならない。Phase は実装順序と完了境界を持つだけであり、詳細本文の責務正本を所有しない。

実装作業の Pull Request は Phase を境界としなければならない。Phase を Pull Request の境界にする理由は、仕様根拠、対象 owner、実装範囲、検証条件、完了判定を一つの責務単位に固定し、複数 Phase の混在、同一 Phase の並行分割、後続 Phase の先取り、仕様根拠のない実装補完を禁止するためである。

Phase を Pull Request の境界にすることは、Phase の一部分だけを完了扱いにすることを意味しない。実装作業の Pull Request は、対象 Phase 全体の実装、検証、証跡、状態整合が完了した単位で扱う。Phase 内に未実装、未検証、仕様不整合、証跡不足、状態更新不足が残る状態を、Pull Request 作成可能、review ready、merge 可能、または完了済みとして扱ってはならない。

<a id="47-実装着手ゲート方針"></a>

<a id="sec-4-7"></a>
**4.7 実装着手ゲート方針：**

実装着手は、[`docs/SPEC.md` ポリシー責務 §0d](SPEC.md#policy-spec-freeze) の仕様凍結条件と [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit) の active Phase 条件をともに満たす対象に限る。条件不足を実装判断で補完してはならない。

<a id="48-完了判定方針"></a>

<a id="sec-4-8"></a>
**4.8 完了判定方針：**

`実装済み` への遷移は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity)、仕様変更の完了は [`docs/SPEC.md` ポリシー責務 §0b](SPEC.md#policy-spec-pr-completion)、Phase の完了単位は [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit) の条件で判定する。Phase 11 の横断 acceptance gate、差し戻し条件、未残条件は [`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8) を正本とし、fixture の一般形式と実装検証証跡は [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務を正本とする。

バグ修正ゼロ化とは、実装済み機能、実装中・検証未完了機能、または Phase 11 対象として [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務に割り当てた検証基盤について、既知の仕様不整合、未検証分岐、未固定の副作用、secret 漏えい可能性、状態 schema 揺れ、fixture 不足、環境依存の合格条件、実装後の追加修正前提を残さない状態をいう。バグ修正ゼロ化は品質目標であり、仕様外の新機能追加、状態語彙の緩和、検証省略、または fixture 期待値の弱体化を許可する理由にしてはならない。

Phase 11 は、バグ修正ゼロ化そのものを目的とする Phase である。Phase 11 は新機能追加 Phase、将来計画実装 Phase、仕様外補完 Phase、検証省略 Phase、または品質目標の一般論を記載する Phase ではない。Phase 11 対象は、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務で対象 owner、対象機能、依存 Phase、現在状態を明示し、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) から対象 owner 詳細本文と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務へ到達できなければならない。Phase 11 の source-code audit inventory、横断 owner 割当、acceptance gate、差し戻し条件、完了時の未残条件は、専用詳細ファイルではなく、[`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8)、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry)、対象 owner 詳細本文、[`docs/details/fixture.md` fixture 証跡責務 Phase 11 fixture harness 参照](details/fixture.md#phase-11-fixture-harness-reference) に分担して記載する。個別の fixture 名、expected file、fake、実装検証証跡の記録形式は [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務、statefile の lock、atomic write、strict schema、JSON Lines、read-only no mutation は [`docs/details/statefile.md`](details/statefile.md) 詳細本文責務を正本とする。

Phase 11 では、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務で Phase 11 の対象として明示されていない新規機能を実装してはならない。Phase 11 の実装対象は、[`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan) の Phase 11 行、および [`docs/ROADMAP.md` 状態・計画責務 §5](ROADMAP.md#522-統合ロードマップ表) で現在状態が `仕様化済み・未実装` かつ `詳細入口 / 次の扱い` が [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) または [`docs/details/fixture.md` fixture 証跡責務 Phase 11 fixture harness 参照](details/fixture.md#phase-11-fixture-harness-reference) へ到達する検証基盤項目だけとする。現在状態が `改訂予定`、`将来計画`、`未仕様化`、または Phase 11 参照入口へ到達しない項目を、Phase 11 で実装、endpoint 化、SDK method 化、UI 操作化、状態 schema 化、配布物化、または外部連携化してはならない。

Phase 11 のバグ修正ゼロ化は、文書上の未完了一覧だけでなく、[`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3) の標準実装 artifact と [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務に実在する [`main.go`](../main.go)、[`components/*.go`](../components/)、[`admin/adlaire-ci-sdk.js`](../admin/adlaire-ci-sdk.js)、[`admin/index.html`](../admin/index.html) の関数、状態 I/O、JSON / JSON Lines 処理、builder output publish / restore、generated site validation、setup install / rollback、release reproducibility / GitHub boundary、admin CLI / SDK / UI client boundary、API / MCP listener lifecycle、graceful shutdown、HTTP header / CORS / cookie boundary、clock / timer / entropy / request ID / time ID、goroutine / channel / worker ordering / cancel / timeout、MCP state bridge / read-only mutation、出力成果物検査、外部 I/O、認証・認可、secret 処理、queue / finalizer、fixture harness 接続を棚卸し対象に含める。棚卸しで検出した差分は、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) で対象 owner と fixture へ割り当て、該当する owner 詳細本文へ実装契約を置き、[`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務へ検証証跡条件を置くか、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務で Phase 11 対象外として到達可能にするまで、完了扱いにしてはならない。既存 test の成功、既存実装済み状態、または別変更で解消するという説明だけで、source audit 由来の artifact 未分類、未固定分岐、直接状態書込、重複 algorithm、best-effort 読込、破損黙殺、secret 応答順序、実 OS listener / signal / sleep / random 依存、未接続 fixture root を残してはならない。

Phase 11 の仕様全般完了は、`source-code audit residual zero`、`artifact coverage zero gap`、`fixture root identity zero duplicate`、`test gap inventory zero open item`、`test / contract drift zero`、`仕様全般完了` がすべて pass し、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務、[`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務のいずれにも Phase 11 対象の未割当、未接続、未検証、重複 fixture、孤立 test、孤立 assertion、対象外理由未到達が残らない場合だけ認める。Phase 11 の完了証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 11 fixture harness 参照](details/fixture.md#phase-11-fixture-harness-reference) へ記録する。仕様全般完了は、実装完了、検証完了、または Phase 11 の `実装済み` 遷移を意味しない。現在状態の変更は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務と [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity) に従う。

<a id="49-仕様策定単位方針"></a>

<a id="sec-4-9"></a>
**4.9 仕様策定単位方針：**

仕様策定は、実装者が 1 つの責務単位として読める範囲でまとめる。

同一の横断契約変更に関係する API、SDK、UI、状態ファイル、認証、セットアップ、検証条件は、一つの仕様策定単位として扱う。ただし、各本文はそれぞれの owner component が持つ責務正本へ記載し、単一文書へ集約してはならない。

仕様変更の統合単位、分割可否、同一 Pull Request 条件、既存 Pull Request への統合、競合解消条件は [`docs/SPEC.md` ポリシー責務 §0c](SPEC.md#policy-spec-change-unit) だけを正本とする。仕様成熟度と未確定事項の扱いは [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity) だけを正本とする。[`docs/SPEC.md` 方針責務 §4.9](SPEC.md#sec-4-9) は仕様策定単位の目的と責務境界だけを定義し、これらの強制条件を再定義しない。

<a id="410-go-正本策定方針"></a>

<a id="sec-4-10"></a>
**4.10 Go 正本策定方針：**

[`docs/SPEC.md`](SPEC.md) 方針責務は、Adlaire CI を最初から Go 言語で設計・実装する前提で策定する。

実装者は、[`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務の Go 前提と禁止事項に従い、配置は [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3)、現在状態は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、実在所在は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務、実装契約は owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を確認する。過去の実装、試作、他言語スクリプト、既存ファイル名、既存 CLI、既存ログ、既存生成物、既存状態ファイルを仕様上の根拠にしてはならない。

Go 実装の判断基準は以下とする。

- 標準 Go 実装 artifact は [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3) のディレクトリ構成 tree に従う。
- `adlaire-ci-build`、`adlaire-ci-runner`、`adlaire-ci-api`、`adlaire-ci-setup`、`adlaire-ci-release`、`adlaire-ci-admin`、`adlaire-ci-mcp` を標準実行バイナリ名とする。
- 仕様未記載の自動変換処理、暗黙の読み替え処理、仕様外分岐を実装判断で追加してはならない。
- 該当する責務正本に記載されていない挙動は、仕様対象外として扱う。

## 5. CI ランナー方針

<a id="51-ci-ランナーの目的"></a>

<a id="sec-5-1"></a>
**5.1 CI ランナーの目的：**

GitHub API を定期的にポーリングし、対象変更を検出してビルドパイプラインを自動実行する自己ホスト型 CI ランナー（[`components/runner.go`](../components/runner.go)）。標準の変更検出・実行経路は GitHub Actions、Webhook、外部 CI サービスのいずれにも依存しない。管理 API の任意の GitHub Webhook 受信機能は補助 trigger 経路であり、無効または未実装でも runner の polling 経路と定期実行は単独で成立しなければならない。

- 変更検出の具体的な API、比較値、保存先は [`docs/details/runner.md`](details/runner.md) 詳細本文責務を正本とする
- Go 標準ライブラリを基本とし、外部依存を追加する場合は [`docs/SPEC.md` ポリシー責務 §4](SPEC.md#policy-dependencies) の例外承認を必須とする

<a id="52-ci-ランナーの開発方針"></a>

<a id="sec-5-2"></a>
**5.2 CI ランナーの開発方針：**

- **単一責務実装**：[`components/runner.go`](../components/runner.go) は CI ランナー責務に限定し、Markdown 変換と管理 API を内包しない
- **実行境界**：oneshot、差分検出、多重実行防止は [`docs/SPEC.md` ポリシー責務 §6](SPEC.md#policy-runner-execution) に従う

<a id="53-github-actions-非依存方針"></a>

<a id="sec-5-3"></a>
**5.3 GitHub Actions 非依存方針：**

Adlaire CI は GitHub Actions の workflow、hosted runner、Marketplace action、push webhook を前提にしない。内製ランナーの起動方式、pipeline 実行方式、secret 参照、変更検出、通知、転送は [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/setup.md`](details/setup.md) 詳細本文責務を正本とする。

<a id="direction-management-tools"></a>

## 6. 管理ツール方針

<a id="61-管理ツールの目的"></a>

<a id="sec-6-1"></a>
**6.1 管理ツールの目的：**

Adlaire CI の状態確認・操作を行う管理インターフェース。ヘッドレスアーキテクチャにより、フロントエンドとバックエンドを明確に分離する。

[`docs/SPEC.md` 方針責務 §6](SPEC.md#direction-management-tools) 以降の管理ツール・管理 API・SDK に関する記載は、現在状態を固定しない機能方針である。対象ファイルの存在は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、現在状態は [`docs/ROADMAP.md`](ROADMAP.md)、実装可否は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity)、検証結果は [`docs/details/fixture.md`](details/fixture.md) を参照する。

<a id="62-ヘッドレスアーキテクチャ方針"></a>

<a id="sec-6-2"></a>
**6.2 ヘッドレスアーキテクチャ方針：**

Adlaire CI と管理ツールは API を介して通信する。フロントエンドとバックエンドを完全に分離し、管理ツールの実装・置き換えをバックエンドから独立させる。

- バックエンド（Adlaire CI）は API を公開する
- フロントエンド（管理ツール）は API のみを通じてバックエンドと通信する
- 直接のファイル操作・プロセス呼び出しは管理ツールから行わない

<a id="63-sdk-方針"></a>

<a id="sec-6-3"></a>
**6.3 SDK 方針：**

API は SDK として提供し、管理ツール実装者が直接 HTTP 通信を記述しなくてよい抽象化レイヤーを提供する。

SDK の実装言語と技術は [`docs/SPEC.md` 方針責務 §4 技術方針表](SPEC.md#direction-technical)、外部依存の禁止条件は [`docs/SPEC.md` ポリシー責務 §4](SPEC.md#policy-dependencies)、API 変更との同期と標準管理ツールからの通信境界は [`docs/SPEC.md` ポリシー責務 §8](SPEC.md#policy-sdk-transport) および [`docs/SPEC.md` ポリシー責務 §9](SPEC.md#policy-9-admin-ui-sdk-boundary)、公開 method と transport の実装契約は [`docs/details/sdk.md`](details/sdk.md) 詳細本文責務だけを正本とする。[`docs/SPEC.md` 方針責務 §6.3](SPEC.md#sec-6-3) は SDK を設ける目的だけを定義し、これらの強制条件を再定義しない。

<a id="64-標準管理ツール方針"></a>

<a id="sec-6-4"></a>
**6.4 標準管理ツール方針：**

Adlaire CI はすぐに使える標準管理ツールを同梱する。

標準管理ツールは、利用者が管理機能を直ちに操作でき、同じ SDK 境界を利用する別クライアントへ差し替えられる標準クライアントとして位置付ける。実装技術は [`docs/SPEC.md` 方針責務 §4 技術方針表](SPEC.md#direction-technical)、SDK 通信とカスタマイズ境界の強制条件は [`docs/SPEC.md` ポリシー責務 §9](SPEC.md#policy-9-admin-ui-sdk-boundary)、視覚仕様は [`docs/DESIGN.md`](DESIGN.md) デザイン責務、DOM・操作・状態の実装契約は [`docs/details/ui.md`](details/ui.md) 詳細本文責務だけを正本とする。[`docs/SPEC.md` 方針責務 §6.4](SPEC.md#sec-6-4) は標準管理ツールを同梱する目的だけを定義し、これらの強制条件を再定義しない。

## 7. 機能インベントリ・ロードマップ参照

機能インベントリ、各機能の現在状態、Phase 一覧、将来計画一覧の本文は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を参照する。状態語彙、実装可否、昇格条件は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity) を参照する。

[`docs/SPEC.md`](SPEC.md) では、機能をどの方針とポリシーで扱うか、Phase 単位でのみ実装する理由、昇格時に満たすべき方針とポリシーだけを定義する。

---

# ポリシー責務
> 遵守義務のある規則と制約を定める。「何をしなければならないか／してはならないか」に答える。

<a id="policy-detail-contract"></a>

## 0. 詳細仕様記載ポリシー

owner component 別の [`docs/details/*.md`](details/) の各仕様項目は、実装者が実装レベルで迷わない記載にしなければならない。[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) は、対象機能から該当する詳細本文と fixture 証跡へ迷わず到達できる参照入口にしなければならない。

詳細仕様、状態管理、fixture、索引、入口文書を記載または整理する場合は、[`docs/SPEC.md` 方針責務 §4.2a 仕様全般重複記載禁止原則](SPEC.md#spec-global-no-duplicate-principle) に従わなければならない。[`docs/SPEC.md` ポリシー責務 §0](SPEC.md#policy-detail-contract) は、責務別の配置、重複記載の定義、参照方法、未確定時の扱いを再定義せず、詳細仕様に必要な実装情報だけを追加定義する。

機能、API、設定、状態ファイル、UI、SDK method、運用手順を仕様化する場合は、該当する必須項目を owner component の詳細本文で特定できる状態にしなければならない。

<a id="detail-contract-required-fields"></a>
**詳細仕様必須項目：**

| 項目 | 必ず特定する内容 |
|------|------------------------|
| 目的 | 解決対象、利用者、成功時に得られる結果。 |
| owner / collaborator | owner component 一件、必要な collaborator component、入出力の受渡し境界。 |
| 実装主体 | 実装 artifact の path または配布物名、起動名、呼出し interface。実行可能な機能でこれらが未定義の場合は仕様化済みとしない。 |
| 入力 | 型、必須性、既定値、許容値、入力元、HTTP method、request 形式。 |
| 出力 | 型、形式、保存先、公開先、戻り値、HTTP status、response 形式。 |
| 設定 | 設定 key、既定値、型、許容範囲、保存先、更新タイミング。 |
| 状態 | ファイル path、形式、schema、初期値、read / write 責務、更新タイミング、破損時処理。 |
| 正常系 | 処理順序、分岐条件、アルゴリズム、成功条件、副作用の確定順序。 |
| 異常系 | エラー条件、応答、終了コード、固定メッセージ、ログ level、通知条件、継続可否。 |
| 制御 | 再試行、lock、排他制御、冪等性、timeout、partial failure、再実行時の扱い。 |
| security | 認証、認可、token、秘密情報、権限、公開境界、出力禁止情報。 |
| 検証 | [`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test) を満たす意味のあるテスト、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 横断テスト証跡共通入口](DETAIL_INDEX.md#cross-test-evidence-route) から到達できる test gap inventory、test improvement batch closure、test evidence package、test verification closure record set、必須 fixture、assertion、requirement coverage、oracle、failure diagnostics、boundary / failure matrix、isolation、determinism、race trigger、concurrency / race、mutation selection、mutation test、harness self-verification、contract drift、skip / 未実行証跡、構文確認、実行確認、生成物確認、必須証跡。 |

「適切に処理する」「必要に応じて対応する」「安全に扱う」のように実装判断を実装者へ委ねる表現を単独で完了仕様として扱ってはならない。使用する場合は、具体的な条件、処理、値、禁止事項、確認方法を併記する。

許可値、許可 path、許可副作用、許可依存、owner 候補、fixture component、schema key を列挙する契約では、`など`、`等`、`任意の同等物`、例示だけの列挙を使用して集合を開いてはならない。許可集合を固定列挙するか、追加値を許可する判定条件、登録先の正本、未登録値の拒否結果を同じ契約で明示する。禁止対象の理解を助ける例示は使用できるが、例示に含まれない対象を許可する意味に読めないことを明記する。

<a id="complete-detail-specification-gate"></a>
**完全仕様詳細化義務：**

すべての実装対象機能は、実装着手前に完全仕様詳細化を完了しなければならない。完全仕様詳細化は [`docs/SPEC.md` 方針責務 §4.4](SPEC.md#complete-detail-specification-policy) の定義、[`docs/SPEC.md` ポリシー責務 §0 詳細仕様必須項目](SPEC.md#detail-contract-required-fields)、[`docs/SPEC.md` ポリシー責務 §0d](SPEC.md#policy-spec-freeze) の仕様凍結条件をすべて満たすことをいう。

以下のいずれかが責務正本で固定されていない場合、その機能は完全仕様詳細化未完了として扱い、実装着手、実装 PR の review ready 報告、merge 可能報告、`実装済み` 判定を禁止する。

| 判定対象 | 固定が必要な内容 |
|----------|------------------|
| 入出力 | 入力元、型、必須性、既定値、許容値、HTTP method、request / response 形式、戻り値、終了コード、HTTP status。 |
| 状態 | 状態ファイル path、schema、初期値、read / write 責務、更新順序、atomicity、lock、破損時処理、read-only no mutation。 |
| 処理 | 正常系順序、分岐条件、境界値、アルゴリズム、成功条件、副作用の確定順序、再実行時の扱い。 |
| 異常系 | 失敗条件、固定 error、固定 message、状態不変条件、partial failure、cleanup failure、rollback 要否、継続可否。 |
| security | 認証、認可、scope、secret 保存禁止、secret 出力禁止、token / password / TOTP / session の露出境界。 |
| 外部境界 | GitHub API、SSH、SMTP、systemd、webhook、browser fetch / SSE、MCP client、command 実行、filesystem の成功 / 失敗 / timeout / malformed 条件。 |
| 制御 | retry、timeout、cancel、clock、timer、entropy、ID generation、goroutine、channel、worker ordering、並行更新、衝突時処理。 |
| 検証証跡 | fixture、fake、expected、effects、security expected、意味のあるテスト条件、test gap inventory、test improvement batch closure、test evidence package、test verification closure record set、requirement coverage ledger、oracle evidence、failure diagnostics evidence、boundary / failure matrix evidence、isolation evidence、determinism evidence、race trigger matrix、concurrency / race evidence、mutation selection ledger、mutation test 条件、harness self-verification evidence、contract drift report、skip / 未実行時の扱い、pass / fail 条件、完了証跡。 |

実装中に完全仕様詳細化未完了の事項を発見した場合、実装者はコード判断で補完してはならない。該当箇所の実装を停止し、責務正本を先に改訂して完全仕様詳細化を完了させてから実装を再開しなければならない。

<a id="policy-spec-maturity"></a>

## 0a. 仕様成熟度ポリシー

仕様項目の状態は、以下の定義に従って管理する。

| 状態 | 定義 | 実装可否 |
|------|------|----------|
| 未仕様化 | 要求、目的、責務、入出力、処理、状態、検証条件のいずれかが実装判断に必要な粒度で定義されていない状態。 | 実装不可 |
| 将来計画 | 将来的な方向性または候補として記録した状態。実装時期、整理順序、詳細仕様は未確定でもよい。整理順序は実装単位、PR 単位、完了判定単位ではない。 | 実装不可 |
| 改訂予定 | 将来計画または未仕様化の項目を仕様化対象へ昇格した状態。詳細仕様の作成・改訂作業中であり、実装条件はまだ満たしていない。 | 実装不可 |
| 仕様化済み・未実装 | [`docs/SPEC.md` ポリシー責務 §0 完全仕様詳細化義務](SPEC.md#complete-detail-specification-gate) を満たし、実装 artifact または実装コードが未作成・未反映の状態。 | 実装可 |
| 実装中・検証未完了 | 仕様に基づくコード変更へ着手済みだが、必須検証、証跡、または関連文書の整合確認が未完了の状態。 | 検証待ち |
| 実装済み | [`docs/SPEC.md` ポリシー責務 §0a 実装完了条件](SPEC.md#implementation-completion-transition) をすべて満たした状態。 | 完了済み |

`仕様化済み・未実装` へ昇格するには、[`docs/SPEC.md` ポリシー責務 §0 詳細仕様必須項目](SPEC.md#detail-contract-required-fields) のうち該当機能に適用する全項目が owner component 詳細本文で特定され、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務から該当本文と fixture 証跡へ到達でき、[`docs/SPEC.md` ポリシー責務 §0 完全仕様詳細化義務](SPEC.md#complete-detail-specification-gate) を満たし、実装判断に必要な未確定事項が残っていないことを必須とする。

`仕様化済み・未実装` は詳細仕様の成熟度を示す状態であり、単独では実装着手許可、Pull Request 作成許可、Phase 完了可能性を意味しない。実装着手には、対象機能が [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務で active Phase に割り当てられ、[`docs/SPEC.md` ポリシー責務 §0d](SPEC.md#policy-spec-freeze) の仕様凍結条件と [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit) の active Phase 条件を同時に満たす必要がある。

詳細仕様入口、owner 詳細本文、または fixture 証跡の参照先が存在するだけでは、`仕様化済み・未実装` へ昇格してはならない。実装 artifact、入力、出力、状態、異常系、検証証跡、fixture root、Phase 割当、実装 PR 境界、完了条件のいずれかが未固定である場合は、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務で `改訂予定`、`将来計画`、または `未仕様化` として扱う。`詳細入口 / 次の扱い` の列にある link は、状態語彙または実装可否を上書きしない。

実装着手は、対象項目が `仕様化済み・未実装` の状態に到達し、かつ [`docs/SPEC.md` ポリシー責務 §0d](SPEC.md#policy-spec-freeze) の仕様凍結条件を満たしている場合に限る。完全仕様詳細化が未完了の項目は、`改訂予定` または `未仕様化` として扱い、実装不可とする。

<a id="implementation-completion-transition"></a>
実装完了は、コード変更だけでは成立しない。仕様との差分確認、構文確認、実行確認または生成物確認、[`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test) を満たす意味のあるテスト、mutation test、必須 fixture、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 横断テスト証跡共通入口](DETAIL_INDEX.md#cross-test-evidence-route) から到達できる test evidence package、test verification closure record set、`final_open_item_count=0` の実装検証証跡、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務の現在状態と [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務の更新要否確認を完了した場合にのみ `実装済み` と扱う。

API、SDK、標準管理ツールのいずれかを変更する場合は、API 仕様、SDK メソッド、UI 操作、詳細仕様の整合を同時に確認する。いずれか一方だけを変更して完了扱いにしてはならない。

<a id="policy-spec-pr-completion"></a>

## 0b. 仕様 PR 完了ポリシー

仕様策定または仕様改訂の Pull Request は、以下を満たすまで完了扱いにしてはならない。

| 対象 | 完了条件 |
|------|----------|
| 方針責務・ポリシー責務 | 方針、ポリシー、実装着手ゲート、完了判定が [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務に明記されている。 |
| 状態・計画責務 | 実装 artifact と各機能の現在状態、Phase、将来計画が [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務に明記され、状態語彙、実装可否、昇格条件が [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity) と矛盾していない。 |
| 詳細仕様入口責務 | 対象機能の owner 対応表、owner 詳細本文から collaborator 境界へ到達する参照、詳細本文参照先、fixture 証跡参照先が [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務に明記されている。 |
| owner component 別詳細本文責務 | 実装に必要な具体値、入出力、状態、処理順序、異常系、検証条件が該当する [`docs/details/*.md`](details/) 詳細本文責務に明記されている。検証条件は [`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test) を満たす。Phase 11 横断対象の場合も、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) から該当 owner 詳細本文と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務へ到達できる。 |
| 横断契約 | API、SDK、UI、状態ファイル、認証、セットアップの対応関係が該当する詳細本文責務で同期している。 |
| 重複記載 | [`docs/SPEC.md` 方針責務 §4.2a 仕様全般重複記載禁止原則](SPEC.md#spec-global-no-duplicate-principle) への適合確認が完了し、未解消違反が 0 件である。 |
| 仕様全般不備一括棚卸し | 仕様全般の再整備、重複箇所、問題点、改善点、または残存しない全件洗い出しを目的にする仕様 PR では、検出した各不備を `重複記載`、`責務外本文`、`参照切れ`、`責務正本未確定`、`状態不整合`、`索引不整合`、`詳細仕様不足`、`fixture 証跡不足`、`デザイン責務混入`、`作業ルール混入`、`実装 artifact 所在不整合` のいずれかに分類し、各件について唯一の責務正本、処置先、処置内容、参照化または削除または正本本文化の結果、解消状態を同一 PR 内で閉じる。完了証跡は [`docs/SPEC.md` ポリシー責務 仕様全般不備 inventory record 固定契約](SPEC.md#spec-deficiency-inventory-record-contract) と [`docs/SPEC.md` ポリシー責務 仕様全般不備 batch closure 固定契約](SPEC.md#spec-deficiency-batch-closure-contract) の両方を満たす。テスト関連の不備は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 横断テスト証跡共通入口](DETAIL_INDEX.md#cross-test-evidence-route) を経由して fixture 証跡責務へ接続する。未分類、責務正本未確定、未解消、対象外理由 anchor 不足、`record_count` 不一致、`final_unresolved_count=0` 未達、category routing 違反、または `source=spec-gap` と `test_gap_connection` の未接続が 1 件でも残る場合は完了扱いにしてはならない。 |
| 索引責務 | ファイル名、正本参照先、実装対象の変更がある場合、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務の更新要否を確認している。 |
| owner 網羅 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) に列挙された owner component のうち、[`docs/ROADMAP.md`](ROADMAP.md) で `未仕様化` または `将来計画` 以外の機能を持つ owner は、少なくとも 1 件の機能から owner 詳細本文と fixture 証跡へ到達できる。実装 artifact の行で owner 機能の現在状態を代用していない。 |
| 列挙閉包 | 許可値、許可 path、許可副作用、許可依存、fixture component、schema key の集合が固定列挙または明示的な登録条件で閉じており、未登録値の扱いが確定している。 |
| デザイン責務 | 生成静的 Web サイトと標準管理 UI の視覚値は [`docs/DESIGN.md`](DESIGN.md) にあり、owner 詳細本文は DOM、selector、状態、操作境界だけを持つ。 |
| fixture 配置 | [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3) の `testdata/` tree、[`docs/details/fixture.md`](details/fixture.md) の配置契約、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) の実在所在が一致している。 |

<a id="spec-deficiency-inventory-record-contract"></a>
**仕様全般不備 inventory record 固定契約：**

仕様全般不備 inventory record は、仕様全般の再整備、重複箇所、問題点、改善点、または残存しない全件洗い出しを目的にする仕様 PR で検出した不備 1 件につき 1 record 作成する。record の正本は [`docs/SPEC.md` ポリシー責務 §0b](SPEC.md#policy-spec-pr-completion) とし、owner 詳細本文、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、[`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務、または Pull Request 本文で同じ schema を再定義してはならない。

| field | 固定値 / 形式 | 未完了条件 |
|-------|---------------|------------|
| `defect_id` | `specdef-<responsibility>-<scope>-<number>` の lowercase kebab-case。 | 空、重複、または責務と scope を識別できない。 |
| `category` | `重複記載`、`責務外本文`、`参照切れ`、`責務正本未確定`、`状態不整合`、`索引不整合`、`詳細仕様不足`、`fixture 証跡不足`、`デザイン責務混入`、`作業ルール混入`、`実装 artifact 所在不整合` のいずれか。 | 分類なし、未登録分類、または複数分類を 1 record に混在している。 |
| `canonical_responsibility` | 不備を最終判断する唯一の責務正本への責務名付き Markdown link。 | 正本なし、裸のファイル名、または複数正本を並列にしている。 |
| `detected_location` | 不備を検出した file path と固定 anchor、または Pull Request evidence label。 | 検出位置が不明、または行番号だけで責務 anchor がない。 |
| `disposition` | `delete`、`reference`、`move-to-canonical`、`define-in-canonical`、`mark-not-applicable` のいずれか。 | 処置が自由記述だけ、または残件を後続 PR 前提にしている。 |
| `target_location` | 処置先の file path と固定 anchor、または削除対象の責務 anchor。 | 処置先が不明、または実在しない path / anchor を指す。 |
| `test_gap_link` | テスト固有の仕様不足を含む場合は [`docs/details/fixture.md` fixture 証跡責務 test gap inventory record 固定契約](details/fixture.md#test-gap-inventory-record-contract) の `gap_id`。含まない場合は `not_applicable` と対象外理由。 | テスト固有の不備なのに `gap_id` がない、または非テスト不備を test gap record だけで閉じている。 |
| `closure_evidence` | PR 本文、差分、または責務正本 anchor への link。 | 処置結果へ到達できない、または説明文だけで完了扱いにしている。 |
| `status` | `open`、`closed`、`not_applicable` のいずれか。 | 完了時に `open` が残る、または `not_applicable` に責務正本 anchor がない。 |

`canonical_responsibility` は、次の category routing に従って 1 件だけ選定する。複数の文書に症状が見える場合でも、record の `canonical_responsibility` は最終的な仕様判断または処置結果を所有する責務正本 1 件に固定する。補助的に参照する文書は `closure_evidence` に置き、`canonical_responsibility` へ複数正本を並べてはならない。

| `category` | `canonical_responsibility` 選定規則 |
|------------|--------------------------------------|
| `重複記載` | 重複している判断対象の正本を [`docs/SPEC.md` 方針責務 §4.2a 仕様全般重複記載禁止原則](SPEC.md#spec-global-no-duplicate-principle) と [`docs/SPEC.md` 責務文書構成表](SPEC.md#document-responsibility-map) で 1 件に確定する。 |
| `責務外本文` | 本文を所有すべき正本を [`docs/SPEC.md` 責務文書構成表](SPEC.md#document-responsibility-map) で 1 件に確定する。 |
| `参照切れ` | 参照対象が文書・実装 artifact の所在である場合は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務、参照対象が仕様本文の anchor である場合はその anchor を所有する責務正本を 1 件に確定する。 |
| `責務正本未確定` | 正本分担の判断は [`docs/SPEC.md` 責務文書構成表](SPEC.md#document-responsibility-map) を `canonical_responsibility` とする。 |
| `状態不整合` | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を `canonical_responsibility` とする。 |
| `索引不整合` | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務を `canonical_responsibility` とする。 |
| `詳細仕様不足` | 不足している実装契約を所有する owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を `canonical_responsibility` とする。 |
| `fixture 証跡不足` | [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務を `canonical_responsibility` とする。 |
| `デザイン責務混入` | [`docs/DESIGN.md`](DESIGN.md) デザイン責務を `canonical_responsibility` とする。 |
| `作業ルール混入` | [`AGENTS.md`](../AGENTS.md) 最上位ルールブックを `canonical_responsibility` とする。 |
| `実装 artifact 所在不整合` | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務を `canonical_responsibility` とする。 |

仕様全般不備 inventory record は、全 record の `status` が `closed` または `not_applicable`、`open` 件数が `0`、テスト固有の仕様不足がある場合は `test_gap_link` から [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 横断テスト証跡共通入口](DETAIL_INDEX.md#cross-test-evidence-route) へ到達し、かつ [`docs/SPEC.md` ポリシー責務 仕様全般不備 batch closure 固定契約](SPEC.md#spec-deficiency-batch-closure-contract) の `test_gap_connection` で当該 `defect_id` と `gap_id` の 1 対 1 対応へ到達できる場合だけ閉じる。仕様全般不備を fixture 証跡責務の `source=spec-gap` だけで閉じること、またはテスト固有の仕様不足を本 record だけで閉じることを禁止する。

<a id="spec-deficiency-batch-closure-contract"></a>
**仕様全般不備 batch closure 固定契約：**

仕様全般不備 batch closure は、仕様全般の再整備、重複箇所、問題点、改善点、または残存しない全件洗い出しを目的にする仕様 PR 1 本につき 1 組だけ作成する。batch closure の正本は [`docs/SPEC.md` ポリシー責務 §0b](SPEC.md#policy-spec-pr-completion) とし、Pull Request 本文、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、[`docs/ROADMAP.md`](ROADMAP.md)、owner 詳細本文、または fixture 証跡責務で同じ schema を再定義してはならない。

| field | 固定値 / 形式 | 未完了条件 |
|-------|---------------|------------|
| `batch_id` | `specbatch-<scope>-<number>` の lowercase kebab-case。 | 空、重複、または対象 scope を識別できない。 |
| `inventory_location` | 仕様全般不備 inventory record 群の所在。Pull Request 本文に置く場合は `pr-verification`、文書内に置く場合は責務名付き Markdown link。 | 所在がない、口頭説明だけ、または実在しない path / anchor を指す。 |
| `record_count` | inventory record の総数を整数で記録する。検出 0 件の場合も `0` を明示する。 | 件数未記録、または inventory record 数と一致しない。 |
| `category_summary` | 使用した category ごとの件数と、未使用 category の `0` を記録する。 | category ごとの件数が不明、または未登録 category を含む。 |
| `canonical_responsibility_summary` | 全 record が本節の category routing に従い、唯一の `canonical_responsibility` を持つことを記録する。 | 複数正本、正本未確定、裸のファイル名、または routing 不一致が残る。 |
| `test_gap_connection` | テスト固有の仕様不足を含む場合は、対応する `defect_id` と [`docs/details/fixture.md` fixture 証跡責務 test gap inventory record 固定契約](details/fixture.md#test-gap-inventory-record-contract) の `gap_id` を 1 対 1 で列挙する。含まない場合は `not_applicable` と対象外理由。 | `source=spec-gap` の test gap と `defect_id` の対応がない、または非テスト不備を test gap だけで閉じている。 |
| `disposition_summary` | `delete`、`reference`、`move-to-canonical`、`define-in-canonical`、`mark-not-applicable` ごとの件数と処置先への link。 | 処置件数が不明、処置先へ到達できない、または後続 PR 前提の処置を含む。 |
| `final_unresolved_count` | 未解消 record 件数を整数で記録し、完了扱いでは `0` に固定する。 | 件数未記録、`0` 以外、または残件を別変更で閉じる説明がある。 |
| `status` | `closed` または `not_applicable` のいずれか。 | `open`、未登録値、または `not_applicable` に対象外範囲と責務正本 anchor がない。 |

仕様全般不備 batch closure は、全 inventory record の `status` が `closed` または `not_applicable`、`record_count` と実 record 数が一致、`final_unresolved_count=0`、category routing 違反 `0`、テスト固有 `source=spec-gap` の `test_gap_link` 未接続 `0`、`test_gap_connection` 未接続 `0`、処置先未到達 `0` の場合だけ閉じる。最終 record より後に残件、暫定対応、後続 PR 前提、または未確認事項を追記した batch closure は完了証跡として扱わない。

仕様 PR は、未確定事項を「推奨」「検討」「適切に」等の表現だけで残してはならない。未確定事項を残す場合は、実装不可の `未仕様化` または `将来計画` として明示する。

<a id="policy-spec-change-unit"></a>

## 0c. 仕様責務分割ポリシー

仕様変更は、責務ベース明示的原則に基づき、変更責務と統合順序が明確な単位で作成する。同一目的、同一仕様領域、同一ファイル群の仕様変更は必ず 1 本の変更単位に集約し、競合状態、未確認状態、重複状態を完了扱いにしてはならない。

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

仕様 PR 作成前の確認コマンド、Git 操作順序、PR 操作手順は [`AGENTS.md` Git 運用ルール](../AGENTS.md#agents-git-operations) を正本とする。

競合解消後は、競合マーカーが残っていないこと、`git diff --check` が成功すること、open PR が同一仕様領域で重複していないこと、統合先 PR の merge 状態が `CLEAN` であることを確認する。`UNKNOWN` は merge 可能として扱わない。

<a id="policy-spec-freeze"></a>

## 0d. 仕様凍結ポリシー

実装着手可能な仕様として扱うには、対象項目を一時的に凍結する。

仕様凍結は、以下をすべて満たす状態をいう。

| 判定項目 | 条件 |
|----------|------|
| 成熟度 | 対象項目が `仕様化済み・未実装` である。 |
| 詳細仕様 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務から対象機能の owner component 詳細本文と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務へ到達でき、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務に入力、出力、状態、処理順序、異常系、検証条件が明記されている。 |
| 必須項目 | [`docs/SPEC.md` ポリシー責務 §0 詳細仕様必須項目](SPEC.md#detail-contract-required-fields) を満たしている。 |
| 完全仕様詳細化 | [`docs/SPEC.md` ポリシー責務 §0 完全仕様詳細化義務](SPEC.md#complete-detail-specification-gate) を満たしている。 |
| 対応表 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.1](DETAIL_INDEX.md#0i1-builder--静的-web-サイト出力)〜[§0i.5](DETAIL_INDEX.md#0i5-setup--release) の詳細節対応表に対象機能が記載され、参照先の owner component 詳細本文と fixture 証跡の対象が一致している。 |
| 横断整合 | API、SDK、UI、状態ファイル、セットアップ、受け入れ条件が矛盾していない。 |
| 未確定事項 | 実装者が設計判断、仕様補完、例外判断、検証条件の推測、既存実装への後追い合わせを行う余地が残っていない。 |
| 変更境界 | 実装 PR で変更してよい範囲と変更してはならない範囲が明確である。 |

仕様凍結後、実装中に仕様不足を発見した場合は、実装 PR 内で独自判断による補完を行わず、仕様改訂 PR または同一 PR 内の仕様改訂コミットで凍結状態を更新する。

仕様凍結は永久固定ではない。変更する場合は、凍結解除理由、変更対象、影響範囲、再検証条件を PR 本文に明記する。

## 0e. 初期実装スコープ確定ポリシー

Go 版初期実装で新規実装へ着手できる対象は、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務で `仕様化済み・未実装` と割り当てられ、かつ [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務の対応表と該当 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務に具体的な実装詳細が存在する範囲に限定する。`実装中・検証未完了` は既着手範囲の継続、修正、検証対象として扱う。`実装済み` は初期実装スコープの現状確認対象には含められるが、新規実装の着手許可を意味しない。

初期実装スコープの実装可否と判定方針は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity)、実装 artifact と個別機能の現在状態、対象 Phase、将来計画該当有無は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、詳細本文への対応は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務、実装契約は owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を正本とする。

初期実装中の対象追加は、[`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit) の途中追加禁止と [`docs/SPEC.md` ポリシー責務 §0d](SPEC.md#policy-spec-freeze) の仕様凍結条件に従う。

初期実装スコープの実装順序は [`docs/ROADMAP.md` 状態・計画責務 §4](ROADMAP.md#roadmap-phase-plan) の Phase 割り当て、PR 分割は [`docs/SPEC.md` ポリシー責務 §0c](SPEC.md#policy-spec-change-unit) と [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit) に従う。

初期実装スコープの完了判定は [`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8) と [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit) に従い、対象 Phase は [`docs/ROADMAP.md` 状態・計画責務 §4](ROADMAP.md#roadmap-phase-plan)、対象機能から詳細本文と fixture 証跡への参照入口は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務を正本とする。

<a id="policy-phase-unit"></a>

## 0f. Phase 実装単位ポリシー

実装順序、実装計画、実装 PR、完了判定は Phase 単位で行わなければならない。

実装を主目的とする Pull Request は、1 本につき 1 Phase だけを対象にしなければならない。複数 Phase の実装変更を 1 本の Pull Request に混在させてはならない。同一 Phase の実装変更を複数の並行 Pull Request へ分割してはならない。active Phase に対応する open Pull Request が既に存在する場合、同じ Phase の追加実装、修正、検証、仕様根拠の補強は新規 Pull Request を作成せず、既存の該当 Phase Pull Request へ統合しなければならない。

対象 Phase 全体が完了するまで、実装を主目的とする Pull Request を作成してはならない。既に active Phase に対応する open Pull Request が存在する場合、その Pull Request は Phase 全体完了まで work in progress として扱い、review ready、merge 可能、完了済みとして報告してはならない。

Phase 全体完了とは、[`docs/ROADMAP.md` 状態・計画責務 §4](ROADMAP.md#roadmap-phase-plan) で対象 Phase に割り当てられた全 owner、全機能、全依存条件について、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務の実装契約、[`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務の必須 assertion と実装検証証跡、[`docs/SPEC.md` ポリシー責務 §0a 実装完了条件](SPEC.md#implementation-completion-transition)、[`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test) の意味のあるテストポリシー、[`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8) の完了判定を満たすことをいう。Phase 11 横断対象は、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) から該当 owner 詳細本文と fixture 証跡へ到達できなければならない。

Phase 内に `仕様化済み・未実装`、`実装中・検証未完了`、未実行の必須検証、未記録の実装検証証跡、未解消の仕様不整合、未反映の状態・索引更新が残る場合、実装者は同一作業ブランチで実装、検証、不整合修正、再検証を繰り返さなければならない。この反復を省略して Pull Request 作成、完了報告、または merge 可能報告を行ってはならない。

仕様全般に基づく実装とは、[`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務、[`docs/ROADMAP.md` 状態・計画責務 §4](ROADMAP.md#roadmap-phase-plan)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務、[`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に到達し、その全てと矛盾しない実装だけを行うことをいう。これらのいずれかで対象 Phase、対象機能、入力、出力、状態、副作用、異常系、検証条件、完了条件が未定義または矛盾している場合、実装者はコード判断で補完してはならない。先に該当する責務正本を改訂し、仕様根拠を確定してから実装しなければならない。

対象 Phase に含まれる全機能は、実装着手前に [`docs/SPEC.md` ポリシー責務 §0 完全仕様詳細化義務](SPEC.md#complete-detail-specification-gate) と [`docs/SPEC.md` ポリシー責務 §0d](SPEC.md#policy-spec-freeze) を満たさなければならない。Phase の一部だけが完全仕様詳細化済みである状態、または一定の仕様だけを固定した状態で、その Phase の実装作業を開始してはならない。

`P0`、`P1`、`P2〜P5` などの優先度ラベル、抽象段階、API 内部分類、fixture 分類を、実装単位、PR 単位、完了判定単位として使ってはならない。

Phase は [`docs/ROADMAP.md` 状態・計画責務 §4](ROADMAP.md#roadmap-phase-plan) に割り当てられた対象 owner、順序、依存関係に従わなければならない。実装契約は owner component 詳細本文、完了判定は [`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8) を正本とする。

新規実装作業の active Phase は同時に一つだけとする。active Phase は、[`docs/ROADMAP.md` 状態・計画責務 §4](ROADMAP.md#roadmap-phase-plan) の順序で最初に `実装済み` でない Phase とする。依存する Phase が `実装済み` でない後続 Phase で、新規機能実装、`仕様化済み・未実装` 機能の実装着手、または Phase 完了判定を行ってはならない。

後続 Phase に既存コードが存在する場合、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務はその事実に基づく現在状態を記録する。ただし、現在状態が `実装中・検証未完了` であることは active Phase であることを意味しない。後続 Phase の既存コードは、active Phase の固定契約を保つために必要な不整合修正、回帰修正、または検証だけを許可し、後続 Phase の機能拡張は許可しない。

Phase の途中で未仕様化、将来計画、改訂予定の機能を追加してはならない。追加する場合は、先に現在状態の割当、詳細仕様、検証条件、受け入れ条件を更新し、仕様凍結を再実施しなければならない。

API の内部説明や fixture 名に既存の段階名が残る場合でも、それらは検証分類としてのみ扱い、実装順序、実装 PR、完了判定の正本にしてはならない。

<a id="policy-meaningful-test"></a>

## 0g. 意味のあるテストポリシー

[`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test) は、実装変更、検証変更、fixture 変更、実装完了判定、Phase 完了判定に適用する。テスト方針と完了可否の本文は本節を正本とし、具体的な fixture、expected、fake、assertion、実行証跡は [`docs/details/fixture.md` fixture 証跡責務](details/fixture.md)、owner 固有の入出力と異常系は owner component 別の [`docs/details/*.md` 詳細本文責務](details/) を正本とする。横断的なテスト証跡の参照入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 横断テスト証跡共通入口](DETAIL_INDEX.md#cross-test-evidence-route) とし、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務はテスト方針、完了可否、fixture 証跡 schema、必須 key、例外条件を再定義してはならない。

意味のあるテストとは、仕様違反、実装の分岐誤り、境界値誤り、副作用漏れ、認証・認可 bypass、secret 漏えい、状態破損、並行処理順序誤り、外部境界の失敗処理漏れ、または回帰を検出できるテストをいう。実装を一度実行するだけの smoke test、assertion が存在しないテスト、戻り値や副作用を検証しないテスト、正常系だけのテスト、実装詳細の存在だけを確認するテスト、coverage percentage だけを満たすテストは、意味のあるテストとして扱ってはならない。

意味のあるテストの完了証跡は、対象 owner、対象仕様 anchor、assertion、実行証跡、未実行理由、対象外理由、mutation 判定、drift 判定、closure 状態へ到達できる記録でなければならない。実行した command 名、pass 件数、coverage percentage、fixture directory の存在、または Pull Request の説明文だけを完了証跡として扱ってはならない。完了証跡の記録単位と記録先は [`docs/details/fixture.md` fixture 証跡責務 test evidence package 記録先固定契約](details/fixture.md#test-evidence-package-record-location-contract)、assertion id は [`docs/details/fixture.md` fixture 証跡責務 assertion id 固定契約](details/fixture.md#test-assertion-id-contract)、skip / 未実行は [`docs/details/fixture.md` fixture 証跡責務 skip / 未実行証跡固定契約](details/fixture.md#test-skip-evidence-contract)、drift report は [`docs/details/fixture.md` fixture 証跡責務 test / contract drift report schema 固定契約](details/fixture.md#test-contract-drift-report-schema-contract) を正本とする。

テスト関連改善作業とは、test、fixture、expected、fake、harness、checker、assertion、mutation、race / concurrency、contract drift、実装検証証跡、または Phase 完了判定のいずれかを変更または評価する作業をいう。テスト関連改善作業は、発見した問題点を個別に散発修正してはならない。最初に全件棚卸しを作成し、各問題点を owner、artifact、仕様 anchor、fixture root、証跡種別、必要対応、closure 記録へ接続し、open item を `0` にする一括 closure まで完了させなければならない。対象変更単位の test evidence package 記録先、問題点の棚卸し schema、batch closure の記録条件、mutation 選定台帳、race trigger 判定表は [`docs/details/fixture.md` fixture 証跡責務 test evidence package 記録先固定契約](details/fixture.md#test-evidence-package-record-location-contract)、[`test gap inventory record 固定契約`](details/fixture.md#test-gap-inventory-record-contract)、[`test improvement batch closure 固定契約`](details/fixture.md#test-improvement-batch-closure-contract)、[`mutation selection ledger 固定契約`](details/fixture.md#mutation-selection-ledger-contract)、[`race trigger matrix 固定契約`](details/fixture.md#race-trigger-matrix-contract) を正本とする。

テスト関連改善作業と仕様全般不備の洗い出しを同一変更で扱う場合、テスト方針と完了可否は [`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test)、仕様全般不備の分類と完了可否は [`docs/SPEC.md` ポリシー責務 仕様全般不備 inventory record 固定契約](SPEC.md#spec-deficiency-inventory-record-contract) と [`docs/SPEC.md` ポリシー責務 仕様全般不備 batch closure 固定契約](SPEC.md#spec-deficiency-batch-closure-contract) を正本とする。テスト固有ではない仕様不備を、test gap inventory の owner component record へ無理に混在させてはならない。テスト固有の仕様不足は `source=spec-gap` として fixture 証跡責務へ接続し、仕様全般不備 batch closure の `test_gap_connection` で `defect_id` と `gap_id` の 1 対 1 対応を閉じる。文書責務の重複、参照切れ、状態不整合、索引不整合、作業ルール混入のような仕様全般不備は [`docs/SPEC.md` ポリシー責務 仕様全般不備 inventory record 固定契約](SPEC.md#spec-deficiency-inventory-record-contract) と [`docs/SPEC.md` ポリシー責務 仕様全般不備 batch closure 固定契約](SPEC.md#spec-deficiency-batch-closure-contract) の完了条件で閉じる。

テスト関連改善作業で open item、未分類 artifact、未接続仕様 anchor、未割当 fixture root、未判定 mutation class、未判定 race trigger、または証跡不足が 1 件でも残る場合は、実装完了、Phase 全体完了、review ready、merge 可能として扱ってはならない。将来計画、対象外、または not applicable とする場合でも、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務または該当 owner 詳細本文の責務正本 anchor へ到達できる記録を必須とし、説明文だけで残件を閉じてはならない。

skip、`t.Skip`、環境機能不足、Docker / Deno / Go toolchain / race detector の未使用、または local 環境都合による未実行は、成功として扱ってはならない。必須検証を実行しない場合は、対象 owner、対象 requirement、未実行理由、代替証跡、再実行条件、完了可否への影響、責務正本 anchor を記録しなければならない。代替証跡または仕様上の対象外理由へ到達できない skip は open item として扱い、`実装済み`、Phase 全体完了、review ready、merge 可能として扱ってはならない。

実装対象機能、実装変更、検証基盤変更は、以下をすべて満たさなければ完了扱いにしてはならない。

| 判定対象 | 必須条件 |
|----------|----------|
| 仕様追跡性 | 各 test、fixture、expected、assertion は、対象 owner、対象仕様、対象 anchor、検出したい仕様違反を特定できる。仕様に存在しない期待値をテストだけへ埋め込まない。 |
| 検証要求網羅 | owner 詳細本文、fixture 証跡、状態・計画責務で実装対象に要求された検証条件は、test artifact、fixture root、expected、assertion、mutation class、closure item のいずれかへ接続し、未検証契約を 0 件にする。 |
| assertion 強度 | 成功条件だけでなく、失敗時 response、stderr、終了コード、状態差分、副作用有無、禁止出力、禁止外部通信を検証する。 |
| failure diagnostics | 各 assertion は、assertion id、対象仕様 anchor、期待値、実値、差分、failure reason、再現条件を特定できる。fixture 名、test 名、pass / fail 件数だけを診断根拠にしない。 |
| 検証基盤自己検証 | test harness、checker、assertion、expected 比較、security assertion、state diff assertion は、検出すべき不正を fail として検出できることを negative control と positive control の両方で証明する。 |
| 正常系 / 異常系 | 正常系、入力不正、欠損、型不一致、境界値、権限不足、状態破損、外部失敗、timeout、partial failure、rollback、cleanup failure を対象機能に応じて固定する。 |
| 境界値 | 空、最小、最大、上限超過、重複、順序差、path、文字コード、時刻、ID、JSON key、HTTP header、CLI option、state schema、archive entry の境界を検証する。 |
| 決定性 | 時刻、乱数、file order、network、GitHub API、systemd、process、並行実行、timer、sleep に依存して結果が揺れない。必要な場合は fake clock、fake entropy、fake filesystem、fake HTTP、固定 fixture を使用する。 |
| 並行処理 / race | goroutine、channel、worker、lock、listener、timer、file lock、queue、shutdown、共有状態更新、並行 request、同時刻 event を扱う変更は、data race、deadlock、goroutine leak、lost update、二重 commit、順序依存、cleanup 漏れを検出できる。 |
| 契約横断 | API、SDK、UI、CLI、statefile、archive、release asset、MCP、setup の境界をまたぐ機能は、呼び出し元と呼び出し先の契約を同じ変更で検証する。 |
| 副作用 | 期待する副作用だけでなく、禁止された状態 write、file 作成、外部通信、通知、log、secret 出力、cache 書込、audit 欠落が発生しないことを検証する。 |
| 回帰 | 修正した不具合、Phase 11 のバグ修正ゼロ化対象、過去に検出した仕様不整合は、同種の再発で失敗する regression test を持つ。 |
| 完了証跡 | 実行結果、未実行項目、対象外理由、実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務](details/fixture.md) へ到達できる形で記録する。 |

owner 詳細本文の検証条件、fixture 証跡条件、状態・計画責務の Phase 対象、既存 test / fixture / expected / assertion のいずれかが test requirement coverage ledger へ到達できない場合は、意味のあるテストの完了証跡として扱ってはならない。test requirement coverage ledger の固定契約、記録項目、対象外理由、open item の扱いは [`docs/details/fixture.md` fixture 証跡責務 test requirement coverage ledger 固定契約](details/fixture.md#test-requirement-coverage-ledger-contract) を正本とする。

弱い oracle、実装結果の丸写し、snapshot の無条件受け入れ、fixture 存在だけの確認、status code だけの確認、stdout / stderr の空確認だけ、状態差分または副作用を確認しない expected、禁止出力または禁止外部通信を確認しない test は、意味のあるテストとして扱ってはならない。test oracle は、正本 anchor、入力、期待 response、期待 error、終了 code、状態差分、effects、security expected、禁止副作用、失敗時 no mutation を対象機能に応じて固定しなければならない。test oracle の証跡 schema、記録項目、expected / actual 比較、禁止副作用、fixture assertion との対応は [`docs/details/fixture.md` fixture 証跡責務 test oracle evidence set 固定契約](details/fixture.md#test-oracle-evidence-set-contract) を正本とする。

assertion id がない test、対象仕様 anchor へ到達できない test、失敗時に期待値、実値、差分、failure reason、再現条件を特定できない test、generic な failure message だけを返す test、pass / fail 件数だけを記録する test、secret を診断出力へ露出する test は、意味のあるテストとして扱ってはならない。failure diagnostics は、assertion id、対象 owner、対象仕様 anchor、expected / actual / diff、failure reason、reproduction command、secret-safe diagnostics を対象機能に応じて固定しなければならない。failure diagnostics の証跡 schema、記録項目、assertion identity、failure reason、secret-safe diagnostics、closure 接続は [`docs/details/fixture.md` fixture 証跡責務 test assertion identity / failure diagnostics evidence set 固定契約](details/fixture.md#test-assertion-failure-diagnostics-evidence-set-contract) を正本とする。

assertion id は、test 名、fixture 名、line number、subtest の表示名、実行順、または自動採番だけで代用してはならない。assertion id は対象 owner、feature、case、assertion の意味を持つ安定識別子とし、名前変更、行番号変更、test 順序変更、fixture directory の移動だけで意味が変わってはならない。assertion id の形式、禁止形式、fixture manifest assertion との接続は [`docs/details/fixture.md` fixture 証跡責務 assertion id 固定契約](details/fixture.md#test-assertion-id-contract) を正本とする。

正常系だけの test、代表的な異常系だけの test、境界値を 1 点だけ確認する test、HTTP status / exit code / error class だけを確認する test、partial failure、rollback、cleanup failure、retry 上限、timeout 境界、size / count / path / ID / schema / header / option の上限下限を固定しない test は、意味のあるテストとして扱ってはならない。boundary / failure matrix は、入力 class、limit、error taxonomy、partial failure、rollback、cleanup failure、retry / recovery、失敗時 no mutation を対象機能に応じて固定しなければならない。boundary / failure matrix の証跡 schema、記録項目、error taxonomy、partial failure、rollback、cleanup failure、closure 接続は [`docs/details/fixture.md` fixture 証跡責務 test boundary / failure matrix evidence set 固定契約](details/fixture.md#test-boundary-failure-matrix-evidence-set-contract) を正本とする。

test 間で共有状態を汚染する test、実行順序に依存する test、fixture / expected を実行中に直接変更する test、環境変数、working directory、temp root、state dir、listener、goroutine、process、timer、file lock、network fake を残留させる test は、意味のあるテストとして扱ってはならない。test isolation は、test 単位の隔離境界、順序入替結果、共有状態の初期化、cleanup failure、parallel 実行可否、残留 resource 検出を対象機能に応じて固定しなければならない。test isolation の証跡 schema、記録項目、順序入替、共有状態、残留 resource、cleanup failure との対応は [`docs/details/fixture.md` fixture 証跡責務 test isolation evidence set 固定契約](details/fixture.md#test-isolation-evidence-set-contract) を正本とする。

非決定的な test、flaky test、retry で偶然成功した test、実時間、乱数、file order、map order、network、GitHub API、systemd、process scheduling、goroutine scheduling、timer、sleep、host 固有 path、OS 差分、外部サービス応答に依存して結果が変わる test は、意味のあるテストとして扱ってはならない。失敗後の再実行で成功した結果、一定回数中の成功率、手元環境での成功、または CI 上の偶発的成功を完了根拠にしてはならない。決定性の証跡 schema、記録項目、fake adapter との対応、再実行一致条件は [`docs/details/fixture.md` fixture 証跡責務 test determinism evidence set 固定契約](details/fixture.md#test-determinism-evidence-set-contract) を正本とする。

data race を検出できない test、goroutine leak を残す test、lock / channel / worker の終了条件を確認しない test、並行 request や同時刻 event の競合結果を固定しない test、共有状態の lost update、二重 commit、二重 cleanup、重複通知、重複 audit、file lock 競合、shutdown 中 request の確定結果を検証しない test は、意味のあるテストとして扱ってはならない。並行処理 / race の証跡 schema、記録項目、race detector、schedule / interleaving、lock / channel / goroutine lifecycle、conflict outcome、atomicity、cleanup との対応は [`docs/details/fixture.md` fixture 証跡責務 test concurrency / race evidence set 固定契約](details/fixture.md#test-concurrency-race-evidence-set-contract) を正本とする。

Go 実装で goroutine、channel、lock、listener、timer、file lock、queue、shutdown、共有状態、並行 request、同時刻 event を扱う変更は、race detector の実行証跡または race detector を対象外にできる責務正本 anchor 付き理由を持たなければならない。race detector を実行しないまま concurrency / race 対象変更を完了扱いにしてはならない。race detector が実行できない環境では、未実行理由と代替 interleaving 証跡を記録しても、対象外根拠がない限り open item として扱う。

test harness、checker、contract drift checker、fixture assertion、expected 比較、security assertion、state diff assertion を変更する場合、またはこれらを実装完了・Phase 完了の根拠として使用する場合は、検出すべき不正 fixture、欠損 expected、禁止副作用、secret leak、fake transcript 不一致、cleanup failure、assertion 無効化を fail として検出できなければならない。常に fail する検証基盤、常に pass する検証基盤、negative control だけの検証、positive control だけの検証、または failure reason を特定できない検証基盤は、意味のあるテストの根拠として扱ってはならない。検証基盤自己検証の証跡 schema、negative control、positive control、fake transcript 検証、closure 接続は [`docs/details/fixture.md` fixture 証跡責務 test harness self-verification evidence set 固定契約](details/fixture.md#test-harness-self-verification-evidence-set-contract) を正本とする。

mutation test（ミューテーションテスト）は必須とする。実装コード、test harness、fixture assertion、expected 比較、security assertion、state diff assertion を変更する場合、対象変更が検出すべき代表的な mutation を定義し、適用可能な mutation を kill しなければならない。mutation test を実施できない実装変更、または適用可能な mutation が生存する実装変更は、`実装済み`、Phase 全体完了、review ready、merge 可能として扱ってはならない。mutation test の証跡 schema、記録項目、fixture manifest との対応は [`docs/details/fixture.md` fixture 証跡責務 mutation test 証跡固定契約](details/fixture.md#mutation-test-evidence-contract) を正本とする。

mutation test は、ゼロ依存・フルインハウス原則に従い、本リポジトリで所有する Go 標準ライブラリ実装または既存の標準検証ランタイムだけで再現できなければならない。外部 mutation testing service、外部 hosted runner、許可外部ライブラリ、npm package、外部 framework、手作業の目視確認、coverage percentage だけを mutation test の完了根拠として使用してはならない。

mutation test は、対象変更単位ごとに mutation operation、変異前、変異後、適用方法、実行 command、期待 failure、実際の failure、判定、集計へ到達できる再現可能な証跡を持たなければならない。mutation の目視確認、説明文だけの mutation、手元で一時的に試しただけの mutation、正本 artifact を直接変更したまま残す mutation、または再実行できない mutation は完了証跡として扱ってはならない。

mutation test は、少なくとも次の mutation class を対象機能に応じて検出対象へ含める。

| mutation class | 検出すべき誤り |
|----------------|----------------|
| 条件反転 | `==` / `!=`、`<` / `<=`、`>` / `>=`、nil 判定、空判定、feature flag、read-only 判定の反転。 |
| 境界値変更 | 上限、下限、timeout、retry 回数、size、件数、permission scope、path depth、token 長、queue 順序の off-by-one。 |
| エラー無視 | parse error、I/O error、HTTP status error、JSON decode error、state lock error、archive validation error、cleanup failure の握りつぶし。 |
| 認証・認可 bypass | session、API token、scope、role、confirmation、read-only、loopback、CORS、CSRF 相当境界の bypass。 |
| 状態更新漏れ | state write skip、atomic rename skip、lock skip、audit skip、metrics skip、history append skip、rollback skip。 |
| 副作用過剰 | 禁止 file write、禁止 external call、禁止 notification、禁止 log、secret 平文出力、read-only 実行時 mutation。 |
| 順序変更 | validation 優先順位、state write 順序、response 構築順、worker order、finalizer、cleanup、再取得順の変更。 |
| response / schema 変更 | HTTP status、JSON key、型、必須 key、error code、CLI stdout / stderr、exit code、SDK error class、MCP JSON-RPC error の変更。 |
| assertion 無効化 | expected file 比較、state diff、effects、security expected、manifest assertion、fixture presence check の削除または常時成功化。 |

mutation の扱いは以下に固定する。

| 判定 | 扱い |
|------|------|
| killed | 意味のあるテストが mutation を検出し、期待どおり失敗した状態。完了条件に使用できる。 |
| survived | mutation が検出されず test が成功した状態。テスト不足または仕様不足として扱い、完了不可とする。 |
| invalid | mutation が構文上成立しない、または対象仕様の観測可能挙動を作れない状態。理由、対象 file、対象 mutation class を証跡へ記録した場合だけ完了判定から除外できる。 |
| equivalent | 責務正本に照らして観測可能挙動が完全に同一であることを、対象 anchor と理由で証明できる状態。証明できない場合は survived とする。 |

coverage は参考指標に限る。line coverage、branch coverage、function coverage、statement coverage のいずれも、意味のあるテスト、mutation test、fixture 証跡、異常系、境界値、契約横断、副作用検証の代替にしてはならない。coverage が高い場合でも、mutation が生存する、assertion が弱い、失敗系がない、仕様追跡性がない、fixture 証跡がない場合は完了不可とする。

テスト未整備の状態で実装を完了扱いにしてはならない。対象機能に対して意味のあるテストまたは mutation test を定義できない場合は、実装判断で対象外にせず、仕様不足として扱い、owner component 別の [`docs/details/*.md` 詳細本文責務](details/) または [`docs/details/fixture.md` fixture 証跡責務](details/fixture.md) を先に改訂しなければならない。仕様上明示された対象外だけは、対象外理由と正本 anchor を証跡へ記録した場合に限り、未実施テストとして扱わない。

<a id="policy-versioning"></a>

## 1. バージョン管理

正式リリース前は、ヘッダーの `仕様バージョン: V.N` と `リリースバージョン: V.X.N` を暫定表記として扱う。暫定表記は、バージョン体系そのものを示す placeholder であり、実在する release tag、GitHub Release、実装済みバージョンを意味しない。暫定表記の期間は、仕様変更、実装変更、ビルド、将来計画追記のいずれによっても `V.N` または `V.X.N` を別の placeholder や推測した数値へ変更してはならない。[`docs/SPEC.md` ポリシー責務 §1](SPEC.md#policy-versioning) のインクリメント規則は、以下の実数運用開始手順を完了した後だけ適用する。

実数運用を開始する場合は、同一 PR で以下をすべて実施する。

1. ヘッダーの `V.N` と `V.X.N` を具体値へ置換する。
2. 具体値に対応する tag / GitHub Release の作成条件を満たしているか確認する。
3. [`README.md`](../README.md)、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) にバージョン表記がある場合は整合させる。
4. 実数運用開始後は placeholder 表記へ戻さない。

<a id="仕様バージョン-vn"></a>

**仕様バージョン V.N：**

[`docs/SPEC.md`](SPEC.md) の仕様バージョンを管理する。

| 項目 | 内容 |
|------|------|
| 形式 | `V.N`（`V` は固定、`N` は正の整数） |
| 更新方針 | 実数運用開始後、[`docs/SPEC.md`](SPEC.md) の変更・追記のたびに `N` を 1 以上インクリメントする。暫定表記中は適用しない。 |
| リセット禁止 | `N` はリセット禁止。`V.1` に戻してはならない |
| 例 | `V.205` → `V.206` → `V.207` |

<a id="リリースバージョン-vxn"></a>

**リリースバージョン V.X.N：**

開発・ビルド・安定版リリースのバージョンを管理する。`N` は開発・ビルドのたびに、`X` は安定版リリースのたびにインクリメントする。

| 項目 | 内容 |
|------|------|
| 形式 | `V.X.N`（`V` は固定、`X` は安定版リリースの正の整数、`N` は開発・ビルド番号） |
| `N` 更新方針 | 実数運用開始後、開発・ビルドのたびに 1 以上インクリメントする。暫定表記中は適用しない。 |
| `X` 更新方針 | 実数運用開始後、安定版リリースのたびに 1 以上インクリメントする。暫定表記中は適用しない。 |
| リセット禁止 | `X` と `N` はリセット禁止。`V.1` に戻してはならない |
| 例 | `V.1.100` → `V.1.101` → `V.2.102`（X は安定版リリースのたびに、N はビルドのたびにインクリメント） |

バイナリが公開するバージョンはリリースバージョンと同一の識別子とする。正式リリース用バイナリは対応する tag と完全一致する具体的な `V.X.N` を返さなければならない。バージョン未注入の local / 検証用開発ビルドは `V.0.0-dev` を返し、この値を release tag、GitHub Release、Release asset、setup / update の導入可能バージョンとして扱ってはならない。実数運用開始後のバージョン付き開発ビルドは、当該ビルドに割り当てた具体的な `V.X.N` を返す。`V.0.0-dev` はヘッダーの暫定表記 `V.X.N` を置換する値ではない。バイナリ出力形式と検証方法は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](DETAIL_INDEX.md#0d-共通固定値) を正本とする。

<a id="github-リリースポリシー"></a>

**GitHub リリースポリシー：**

| 項目 | ルール |
|------|--------|
| リリース作成条件 | 安定版リリース（`X` インクリメント時）のみ GitHub Release を作成する。開発・ビルド（`N` インクリメントのみ）では作成しない |
| タグ形式 | `V.X.N`（リリースバージョンと一致させる）例：`V.2.102` |
| リリースタイトル | タグ名と同一にする |
| リリース形式 | バイナリ配布を標準とする。利用者は GitHub Release から OS/arch 別の実行バイナリを取得し、ソースからのビルドを標準導入手順に含めない |
| バイナリバージョン | Release 用実行バイナリ asset の `--version` が、バイナリ名に続く第 2 token として tag と同一のリリースバージョンを返す。`V.0.0-dev`、tag と不一致の値、未注入値を含む実行バイナリ asset は公開禁止とする |
| 標準 OS/arch | 初期標準は Linux x86_64（`linux-amd64`）とする。追加 OS/arch は将来のリリース対象として個別に仕様化する |
| 添付ファイル | builder、runner、管理API、setup、管理 CLI、MCP の各`linux-amd64`実行バイナリ、管理UI archive、checksum manifestの8件だけを添付する。exact asset名と生成条件は[`docs/details/release.md` 詳細本文責務 §R3](details/release.md#release-asset-contract)を正本とする。`adlaire-ci-release`、生成静的Webサイト、debug binaryをRelease assetとして添付しない |
| 管理 UI 配布物 | 管理UIは実行バイナリと分離した単独のarchive assetとして全安定版Releaseへ添付する。exact asset名と生成条件は[`docs/details/release.md` 詳細本文責務 §R3](details/release.md#release-asset-contract)、archive内容は[`docs/details/admin.md` 詳細本文責務 §A1](details/admin.md#a1-管理-ui-静的ファイル境界)、setup側の安全検証は[`docs/details/admin.md` 詳細本文責務 §A2](details/admin.md#a2-管理-ui-archive-検証)を正本とする |
| checksum | checksum manifest 自身を除く全 Release asset の SHA-256 checksum を checksum manifest で提供する。checksum manifest 自身を再帰的な checksum 対象にしてはならない。exact asset 名、生成形式、公開前検証は [`docs/details/release.md` 詳細本文責務 §R3](details/release.md#release-asset-contract) を正本とし、受け入れ・配置前検証は [`docs/details/setup.md` 詳細本文責務 §26.2a](details/setup.md#sec-26-2a) を正本とする |
| プレリリースフラグ | 安定版リリースでは `Pre-release` にチェックを入れない |
| ドラフト公開禁止 | Draft Release のまま公開しない |

<a id="仕様変更とバージョン更新条件"></a>

**仕様変更とバージョン更新条件：**

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
| 将来拡張 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 | 未実装の設定拡張は現在状態を先に確定する。 |

<a id="policy-dependencies"></a>

## 4. 外部ライブラリ・フレームワーク方針

<a id="基本原則"></a>

**基本原則：**

ゼロ依存・フルインハウスの定義、各 component の自律性、仕様不足時の扱いは [`docs/SPEC.md` 方針責務 §4.1](SPEC.md#sec-4-1) を正本とする。この節は例外採用の禁止条件と許可一覧だけを追加定義する。

<a id="外部フレームワーク"></a>

**外部フレームワーク：**

本リポジトリが所有する実装 artifact、標準 SDK、標準管理 UI、build/runtime 依存、配布物へ外部フレームワークを導入することは、**いかなる条件でも禁止する。** 例外なし。

外部 consumer application が [`admin/adlaire-ci-sdk.js`](../admin/adlaire-ci-sdk.js) を ES Module として読み込む環境は、この禁止対象ではない。React、Vue、Svelte を含む consumer 側 framework から SDK を利用できるが、その framework、adapter、型定義 package、bundler、polyfill を本リポジトリの依存、配布物、検証前提へ追加してはならない。

<a id="外部ライブラリ"></a>

**外部ライブラリ：**

外部ライブラリは禁止する。例外採用は、以下の条件をすべて満たし、かつ [`docs/SPEC.md` ポリシー責務 §4](SPEC.md#policy-dependencies) の許可外部ライブラリ一覧に登録した場合に限る。

- 内製化が技術的に困難であり、標準ライブラリだけでは安全性または正確性を担保できない
- 採用範囲が単一責務に限定され、コンポーネント全体の自律性を壊さない
- 採用理由、代替困難性、責務範囲、削除方針、検証条件を [`docs/SPEC.md`](SPEC.md) ポリシー責務に明記している

開発コスト短縮、実装の容易さ、流行、一般的なベストプラクティスだけを理由にした採用は認めない。許可リスト外のライブラリ使用は認めない。

<a id="内製共通処理"></a>

**内製共通処理：**

内製共通処理の採用可否、禁止名称、依存方向、同格性は [`docs/SPEC.md` 方針責務 §4.2](SPEC.md#sec-4-2) を正本とする。責務、入力、出力、状態、異常系、検証条件が未定義の共通処理は採用できない。

<a id="内製実装管理ポリシー"></a>

**内製実装管理ポリシー：**

内製実装の配置は [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3)、実在所在は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、現在状態と Phase 割当ては [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を正本とする。実装着手可否は、[`docs/SPEC.md` 方針責務 §4.7](SPEC.md#sec-4-7)の着手ゲート、[`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity)の実装可否、[`docs/SPEC.md` ポリシー責務 §0d](SPEC.md#policy-spec-freeze)の凍結条件、[`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit)の active Phase 条件、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務の現在状態と Phase 割当てをすべて使用して判定する。この節では再定義しない。

<a id="許可外部ライブラリ一覧"></a>

**許可外部ライブラリ一覧：**

内製スクリプト・ライブラリは許可外部ライブラリ一覧に記載しない。内製の管理は内製実装管理ポリシーで行う。

| ライブラリ | 用途 | 許可理由 |
|-----------|------|---------|
| （なし） | — | — |

> 許可外部ライブラリは存在しない。

---

<a id="policy-runner-secrets"></a>

## 5. CI ランナー秘密情報・公開境界ポリシー

- GitHub PAT（Personal Access Token）はスクリプト内にハードコードしてはならない
- PAT は最小権限とし、標準構成では対象リポジトリの `Contents: Read` だけを許可する。任意機能の Commit Status を有効にした場合だけ `Commit statuses: Write` を追加でき、他の GitHub 書き込み API へ流用してはならない
- PAT の GitHub repository permission と API 呼び出し条件は [`docs/details/runner.md` 詳細本文責務 §17](details/runner.md#17-github-連携前提) と [`docs/details/commitstatus.md` 詳細本文責務 §27.1](details/commitstatus.md#sec-27-1)、保存 path と file 境界は [`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、読み取り、mask、漏えい禁止は [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/security.md`](details/security.md) 詳細本文責務を正本とする
- ランナーは外部公開エンドポイントを持たない。サーバーから GitHub API への送信のみで動作する

<a id="policy-runner-execution"></a>

## 6. CI ランナー実行境界ポリシー

- 変更がない場合はビルドをスキップしなければならない
- `adlaire-ci-runner` は oneshot 実行とし、多重実行を防止しなければならない
- 成功、失敗、再試行、状態更新、ログ記録の具体条件は [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/statefile.md`](details/statefile.md) 詳細本文責務を正本とする

## 7. CI ランナー branch target ポリシー

- ランナーは 1 件以上のブランチターゲットを扱える方針とする
- ブランチターゲットの fields、既定値、処理順序、並列可否、検出間隔は [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/statefile.md`](details/statefile.md) 詳細本文責務を正本とする
- ブランチターゲットを変更する場合は、[`docs/details/runner.md`](details/runner.md) 詳細本文責務、[`docs/details/statefile.md`](details/statefile.md) 詳細本文責務、[`docs/details/setup.md`](details/setup.md) 詳細本文責務を同時に整合しなければならない

<a id="policy-sdk-transport"></a>

## 8. SDK 通信契約ポリシー

[`docs/SPEC.md` ポリシー責務 §8](SPEC.md#policy-sdk-transport) は、[`admin/adlaire-ci-sdk.js`](../admin/adlaire-ci-sdk.js) に適用する。

SDK の実装言語と技術は [`docs/SPEC.md` 方針責務 §4 技術方針表](SPEC.md#direction-technical)、外部依存の禁止条件は [`docs/SPEC.md` ポリシー責務 §4](SPEC.md#policy-dependencies) だけを正本とする。この節は SDK の通信契約だけを定義し、実装技術または外部依存条件を再定義しない。

- SDK の対応言語を変更する場合は、[`docs/SPEC.md` 方針責務 §4 技術方針表](SPEC.md#direction-technical) の改訂を実装より先に完了する
- バックエンド API の変更は SDK の更新を伴う
- SDK の module 形式、公開 API、error class、timeout、streaming 契約は [`docs/details/sdk.md`](details/sdk.md) 詳細本文責務を正本とする
- SDK は API response と利用者操作の境界を変更してはならない。禁止する補完、再試行、永続化、global 汚染の具体条件は [`docs/details/sdk.md`](details/sdk.md) 詳細本文責務を正本とする

<a id="policy-9-admin-ui-sdk-boundary"></a>

## 9. 標準管理ツール UI / SDK 境界ポリシー

[`docs/SPEC.md` ポリシー責務 §9](SPEC.md#policy-9-admin-ui-sdk-boundary) は、[`admin/index.html`](../admin/index.html) に適用する。

標準管理ツールの実装技術は [`docs/SPEC.md` 方針責務 §4 技術方針表](SPEC.md#direction-technical)、外部依存の禁止条件は [`docs/SPEC.md` ポリシー責務 §4](SPEC.md#policy-dependencies) だけを正本とする。この節は UI と SDK の境界だけを定義し、実装技術または外部依存条件を再定義しない。

- バックエンドとの通信はすべて SDK 経由とする
- カスタマイズを妨げる密結合な実装を禁止する。SDK 境界、DOM 境界、状態管理境界を満たさない UI 実装は完了扱いにしてはならない
- DOM、form、初期ロード順、イベント処理順、成功/失敗表示、秘密情報消去条件は [`docs/details/ui.md`](details/ui.md) 詳細本文責務を正本とする
- token の保持・消去と API transport 禁止対象の具体条件は [`docs/details/ui.md`](details/ui.md) 詳細本文責務を正本とする

## 10. 状態ファイル永続化ポリシー

データベースは使用せず、状態はファイルベースで永続化する。構造化状態の標準形式は JSON object、JSON array、または JSON Lines とする。例外は、[`docs/details/statefile.md`](details/statefile.md) 詳細本文責務の状態ファイル固定表に path と形式を明示した UTF-8 text、gzip 圧縮 JSON、directory、および [`docs/details/archive.md`](details/archive.md) 詳細本文責務に形式を明示した tar.gz artifact だけとする。固定表または archive 詳細本文責務に明示されていない非 JSON 形式を追加してはならない。具体的な状態ファイル一覧、schema、権限、更新順序、破損時処理は [`docs/details/statefile.md`](details/statefile.md) 詳細本文責務を正本とする。技術選択を変更する場合は、[`docs/SPEC.md` 方針責務 §4 技術方針表](SPEC.md#direction-technical) の改訂と外部依存の許可判断を実装より先に完了しなければならない。

<a id="policy-single-user-auth"></a>

## 11. シングルユーザー認証ポリシー

[`docs/SPEC.md` ポリシー責務 §11](SPEC.md#policy-single-user-auth) は、管理 API サーバーおよび標準管理ツールに適用する。

- 初期構成の既定はシングルユーザーとする。マルチユーザー対応の現在状態は [`docs/ROADMAP.md`](ROADMAP.md)、詳細仕様は [`docs/details/security.md`](details/security.md) を正本とする
- 初期認証、初回変更、強制変更、session、保存形式の具体条件は [`docs/details/api.md`](details/api.md) 詳細本文責務と [`docs/details/security.md`](details/security.md) 詳細本文責務を正本とする
- パスワードは平文保存禁止とし、保存が必要な認証情報はハッシュ化または secret として扱う

<a id="policy-api-exposure"></a>

## 12. 管理 API 公開境界・session ポリシー

[`docs/SPEC.md` ポリシー責務 §12](SPEC.md#policy-api-exposure) は、[`components/api.go`](../components/api.go) に適用する。

- 管理 API サーバーを外部へ直接公開してはならない
- TLS 終端、listen host、session token、認証除外 endpoint の具体条件は [`docs/details/api.md`](details/api.md) 詳細本文責務と [`docs/details/security.md`](details/security.md) 詳細本文責務を正本とする
- セッショントークンは永続化してはならない
- API エンドポイントは、詳細本文責務で明示的に認証除外されたものを除き、認証必須とする
