# Adlaire CI — 仕様ドキュメント

**対象コンポーネント：** `components/builder.go`（ビルドスクリプト、実装済み）/ `components/runner.go`（CI ランナー、実装済み）/ `components/api.go`（管理 API サーバー、実装中・検証未完了）/ `admin/adlaire-ci-sdk.js`（JavaScript SDK、仕様化済み・未実装）/ `admin/index.html`（標準管理ツール、仕様化済み・未実装）/ `components/mcp.go`（MCP サーバー、将来計画）
**出力形式：** 静的 Web サイト（HTML / CSS / JavaScript / search index）
**スクリプトバージョン：** v3（Adlaire Design System ブルートークン正式採用）
**仕様バージョン：** V.N（正式リリース前の暫定表記）/ **リリースバージョン：** V.X.N（正式リリース前の暫定表記） → Part 2 §2 参照
**最終更新：** 2026-09-15

---

> **Adlaire CI** とは、最初から Go を前提として仕様策定するビルド・CI・管理システムの総称である。正本コンポーネントは `components/builder.go`、`components/runner.go`、`components/api.go`、`admin/index.html`、`admin/adlaire-ci-sdk.js` とする。将来的には `components/mcp.go`（MCP サーバー）を加えた構成へ拡張予定（→ §13 将来計画 MCP サーバー実装）。

## 文書責務

Adlaire CI の仕様判断では、次の責務分担を固定する。

| 文書 | 正本範囲 | 記載する内容 | 記載しない内容 |
|------|----------|--------------|----------------|
| `ADLAIRE_CI_SPEC.md` | 方針、ポリシー、実装状態、ロードマップ、リリース判断 | 目的、設計方針、禁止事項、成熟度、実装可否、仕様昇格手順 | 関数単位の処理、HTTP response schema、状態ファイル schema、具体的な実行手順 |
| `ADLAIRE_CI_DETAIL_SPEC.md` | 実装詳細の入口 | 詳細仕様の読み方、共通固定値、実装前確認項目、検証マトリクス、Phase、詳細節対応表、リポジトリ内ソース配置、責務 component 別詳細仕様ファイル管理仕様、横断補足契約 | 各 component の詳細な処理本文、方針、ポリシー、実装状態、実装可否、ロードマップ状態、PR 分割判断 |
| `ADLAIRE_CI_DETAIL_*_SPEC.md` | owner component 別の実装詳細 | CLI、API、SDK、UI、状態ファイル、処理順序、異常系、セットアップ、受け入れ条件の本文 | 方針、ポリシー、実装状態、実装可否、ロードマップ状態、PR 分割判断 |
| `DOCUMENT_INDEX.md` | 文書・実装ファイル索引 | ファイルの役割、正本関係、実装ファイルの所在 | 仕様本文、詳細仕様、実装状態の最終判断 |
| `DESIGN.md` | 生成静的 Web サイトのデザイン補助 | 見た目、レイアウト、デザイントークン参照 | CI、API、運用、実装可否の判断 |

同じ内容を複数文書に重複定義してはならない。方針や実装可否は本ファイルを正とし、入口、読み順、共通固定値、対応表、横断補足契約は `ADLAIRE_CI_DETAIL_SPEC.md`、入出力・状態・処理・検証の本文は owner component 別の `ADLAIRE_CI_DETAIL_*_SPEC.md` を正とする。

文書を整理する場合も、上表の正本範囲を越えてはならない。読み順、索引、参照先の整理は許可するが、詳細仕様本文、方針、実装状態、ロードマップ状態を別文書へ重複定義しない。

## 実装状態

本ドキュメントでは、仕様化済みの内容と実装済みの内容を区別して扱う。

下表のコンポーネント名は、Part 1 §4.3 の標準ソース配置に基づく。標準ソース配置への実装移行は完了済みであり、`main.go`、`components/*.go`、`testdata/<component>/` を現行配置として扱う。

| コンポーネント | 状態 | 備考 |
|---------------|------|------|
| `components/builder.go` | 実装済み | Go 版 Markdown → 静的 Web サイトビルドスクリプトとして Phase 1 の `gofmt` と `go test` 検証済み。 |
| `components/runner.go` | 実装済み | Go 版 CI ランナーとして Phase 2 完了判定パスの `gofmt` と `go test` 検証済み。 |
| `components/api.go` | 実装中・検証未完了 | Go 版管理 API サーバー。実装ファイルとテストは存在するが、実装状態表、詳細仕様、SDK/UI 連携、検証結果がすべて実装済みとして整合するまでは実装済みへ昇格しない。 |
| `admin/adlaire-ci-sdk.js` | 仕様化済み・未実装 | 管理ツール用 JavaScript SDK。仕様は本ドキュメントに定義するが、リポジトリには実装ファイルが存在しない。 |
| `admin/index.html` | 仕様化済み・未実装 | 標準管理ツール UI。仕様は本ドキュメントに定義するが、リポジトリには実装ファイルが存在しない。 |
| `components/mcp.go` | 将来計画 | Go 版 MCP サーバー。将来計画として管理し、実装済みとは扱わない。 |

仕様化済み・未実装、または将来計画の項目を、実装済み機能として扱ってはならない。

---

## 本ドキュメントの構成

| Part | 名称 | 責務の問い | 記載する内容 |
|------|------|-----------|------------|
| Part 1 | 方針 | **なぜ・何を** | 目的、設計思想、方向性の原則。変更頻度が低く、判断の拠り所となる指針 |
| Part 2 | ポリシー | **しなければならない／してはならない** | 遵守義務のある規則・制約・禁止事項。セキュリティ要件・運用ルール・バージョン管理規則 |
| Part 3 | 仕様 | **どのように** | 実装の具体的詳細。入口、対応表、共通固定値、横断補足契約は `ADLAIRE_CI_DETAIL_SPEC.md`、要件・構成・API・アルゴリズム・設定値・手順の本文は owner component 別の `ADLAIRE_CI_DETAIL_*_SPEC.md` を正とする |

新しい記載内容は「この内容はどの責務の問いに答えるか」を基準に Part を決定する。

実装者は、実装前に以下の順で読む。

1. 本ファイルの「文書責務」と「実装状態」で、対象がどの文書・状態に属するかを確認する。
2. Part 1 §4a〜§4f で、詳細仕様粒度、成熟度、着手ゲート、完了判定を確認する。
3. Part 1 §12 と §13 で、対象機能が実装可か将来計画かを確認する。
4. Part 2 で、対象領域の禁止事項、セキュリティ、バージョン、外部依存を確認する。
5. `ADLAIRE_CI_DETAIL_SPEC.md` で読み順、共通固定値、§0i.1〜§0i.4 の対応表、§0j のリポジトリ内ソース配置を確認し、該当する owner component 別の `ADLAIRE_CI_DETAIL_*_SPEC.md` で実装に必要な入出力、状態、異常系、検証条件、配置を確認する。

---

# Part 1 — 方針
> 目的・設計思想・方向性の原則を定める。「なぜこう作るか」に答える。

## 1. 目的

`components/builder.go` は、任意の UTF-8 Markdown ファイルまたは Markdown ディレクトリを、静的配信可能な Web サイトへ変換する Go プログラムである。本仕様では、Go 実装を最初からの正本として定義する。

- 大規模 Markdown 仕様書、複数 Markdown ドキュメント、運用メモを、快適に閲覧できる静的 Web サイトへ変換する
- `index.html`、ページ HTML、共通 CSS、共通 JavaScript、検索 index を出力ディレクトリへ生成する
- 初期テーマ `adlaire-default` と固定テーマコンポーネントにより、一貫したデザイン言語を維持する
- Adlaire Design System（ADS）のトークンを採用し、テーマの見た目は ADS 準拠の範囲内で管理する

## 2. 開発方針

本プロジェクトは **仕様駆動開発（Spec-Driven Development）** を採用する。

- **仕様書が唯一の真実（Single Source of Truth）**：実装の追加・変更はすべて本ドキュメントへの反映を先行させる
- **仕様→実装の順序**：仕様書に基づいてスクリプトを修正する
- **仕様との乖離は不整合**：乖離が生じた場合も仕様書を先に改訂し、仕様書に基づいて実装を修正する

## 3. デザイン方針

docs.rs / MDN に倣った技術ドキュメントレイアウト。14,000 行超の仕様書を快適に閲覧するため、**構造の明快さ**と**情報密度への耐性**を最優先とする。

- ヘッダーのみアクセントカラーを使う。コンテンツ・サイドバーは中性色ベース
- CSS カスタムプロパティは [Adlaire Design System](https://github.com/fqwink/Adlaire-Design-System)（`Tokens/`）定義の `--adlaire-*` トークンのみ使用
- **ライトモード固定**（`prefers-color-scheme` 非対応、ダークモードなし）
- 外部フォント不使用。システムフォントスタックで日本語環境の可読性を確保

## 4. 技術方針

| 領域 | 方針 |
|---|---|
| ランタイム | Go stable release |
| 言語 | Go（ビルド、ランナー、管理 API）/ JavaScript（SDK）/ HTML・CSS・Vanilla JavaScript（UI） |
| HTTP | Go 標準ライブラリ `net/http` |
| データベース | なし（ファイルベース） |
| Git 操作 | GitHub REST API（Blobs API）を Go 標準ライブラリ `net/http` 経由で呼び出す |
| フロントエンド | HTML / CSS / Vanilla JavaScript |
| 推奨運用 | CI サーバー（VPS 等）でビルドし、静的コンテンツ配信サーバーへ SSH で転送する 2 サーバー構成 |
| データ交換形式 | JSON に統一する。エクスポート・インポートを含む全 API データ交換に CSV・XML 等の非 JSON 形式を使用しない |

### 4.1 ゼロ依存・フルインハウス原則

Adlaire CI は、ゼロ依存・フルインハウスを技術哲学の中核とする。

本原則におけるゼロ依存とは、各コンポーネントが外部ライブラリ、外部フレームワーク、外部ビルドツール、外部ホスティング実行基盤に機能成立を依存しないことを意味する。本原則におけるフルインハウスとは、Markdown 変換、CI 実行、管理 API、SDK、標準管理ツール、状態管理、認証、ログ、通知、セットアップの主要機能を本リポジトリ内で仕様化し、内製コードとして理解、検証、保守できる状態を意味する。

実装者は、便利さ、実装速度、一般的な慣習を理由に外部ライブラリで未定義機能を補完してはならない。外部依存がなければ成立しない設計は、原則として設計不備として扱い、先に仕様を見直す。

各コンポーネントの自律性は以下を満たす。

| コンポーネント | 自律性の条件 |
|----------------|--------------|
| `components/builder.go` | Go 標準ライブラリだけで Markdown 解析、HTML/CSS/JS/search index 生成、テーマコンポーネント出力、検証レポート出力を行う。外部 Markdown parser、template engine、syntax highlight library、search library に依存しない。 |
| `components/runner.go` | Go 標準ライブラリと OS 標準コマンドだけで GitHub API polling、SHA 比較、ビルド起動、ログ、通知、SSH 転送、snapshot、lock、retry を処理する。外部 CI サービス、job queue、scheduler library に依存しない。 |
| `components/api.go` | Go 標準ライブラリ `net/http` を基本に、認証、session、状態ファイル CRUD、入力検証、API response を内製実装する。外部 web framework、router、ORM、database driver に依存しない。 |
| `admin/adlaire-ci-sdk.js` | 単一 ES Module とし、browser 標準 API のみで API client、error handling、streaming、timeout を実装する。npm package、bundler、polyfill、framework に依存しない。 |
| `admin/index.html` | HTML / CSS / Vanilla JavaScript だけで標準管理ツールを構成し、SDK 経由で通信する。React、Vue、Svelte、CSS framework、icon package、chart library に依存しない。 |
| `components/mcp.go` | 将来計画の段階でも Go 標準ライブラリを前提とし、MCP 通信、JSON-RPC 処理、API bridge、監査ログを内製する。外部 MCP framework に依存する前提で仕様化しない。 |

外部依存を例外採用する場合は、Part 2 §4 の許可外部ライブラリ一覧へ登録し、採用理由、代替困難性、責務範囲、削除方針、検証条件を同一 PR で明記する。許可リストにない外部依存は、実装済みとして受け入れない。

### 4.2 Core 非採用・共通責務コンポーネント方針

Adlaire CI は、`core`、`adlaire-ci-core`、`internal/core`、`common`、`base`、`foundation`、`utils` のような中心化・汎用置き場化する概念を採用しない。

横断的に利用される処理は、必要に応じて共通責務コンポーネントとして定義できる。ただし、共通責務コンポーネントは他コンポーネントより上位ではなく、同格の独立コンポーネントとして扱う。共通責務コンポーネントを、正本、中心、基盤、親、上位レイヤーとして扱ってはならない。

共通責務コンポーネントを追加する場合は、以下をすべて満たす。

- 単一責務である。
- 入力、出力、状態、異常系、検証条件が `ADLAIRE_CI_DETAIL_SPEC.md` に仕様化されている。
- 利用元コンポーネントとの依存方向が明確である。
- アプリケーション固有の判断、業務判断、UI 判断、API endpoint 判断、ビルド対象判断を含まない。
- 何でも置き場として使用しない。
- `core`、`common`、`base`、`foundation`、`utils` など中心性または汎用置き場を示す名称を使わない。
- 追加時に `DOCUMENT_INDEX.md`、本ファイルの実装状態表、Part 2 §4 の内製スクリプト一覧、`ADLAIRE_CI_DETAIL_SPEC.md` の詳細仕様と受け入れ条件を整合する。

許可される共通責務コンポーネントの例は以下とする。

| 例 | 責務 |
|----|------|
| `state_store` | 状態ファイルの atomic write、lock、JSON 読み書き。 |
| `secret_masker` | PAT、Webhook Secret、SMTP password、session token、API token のマスク処理。 |
| `event_log` | 構造化イベントログの形式、出力、読み取り。 |
| `github_client` | GitHub REST API 呼び出し、rate limit、retry 境界。 |

禁止される設計は以下とする。

| 禁止対象 | 理由 |
|----------|------|
| `core` / `adlaire-ci-core` / `internal/core` | 中心コンポーネント化し、正本や上位概念と誤読されるため。 |
| `common` / `base` / `foundation` / `utils` | 責務が曖昧になり、何でも置き場化しやすいため。 |
| 全共通処理の一括集約 package | 単一責務を失い、コンポーネント境界を壊すため。 |
| 業務判断を含む共通処理 | 利用元コンポーネントの責務を横断基盤へ漏らすため。 |

共通責務コンポーネントは、コード共有のためだけに追加してはならない。重複削減より、責務境界、仕様の明確さ、依存方向の追跡可能性を優先する。

### 4.3 リポジトリ構成方針

Adlaire CI のリポジトリ内ソース構成は、責務ベースで整理する。

標準構成は以下とする。

```text
.
├── main.go
│
├── components/
│   ├── builder.go
│   ├── runner.go
│   ├── api.go
│   ├── admin.go
│   ├── statefile.go
│   ├── archive.go
│   ├── commitstatus.go
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
│   ├── admin/
│   ├── statefile/
│   ├── archive/
│   ├── commitstatus/
│   └── mcp/
│
├── docs/
│   └── examples/
│
├── ADLAIRE_CI_SPEC.md
├── ADLAIRE_CI_DETAIL_SPEC.md
├── DOCUMENT_INDEX.md
├── DESIGN.md
├── README.md
├── AGENTS.md
└── go.mod
```

`main.go` は 1 ファイルとし、起動入口、サブコマンド判定、引数受け取り、対象コンポーネント呼び出しだけを担当する。`main.go` に Markdown 変換、CI 実行、HTTP handler、状態ファイル操作、archive 処理、GitHub Commit Status 送信、MCP 処理の実装詳細を書いてはならない。

`components/` は、1 コンポーネント = 1 Go ファイルとする。ファイル名は責務名を表し、`builder.go`、`runner.go`、`api.go`、`admin.go`、`statefile.go`、`archive.go`、`commitstatus.go`、`mcp.go` を標準コンポーネントとする。

`admin/` は標準管理 UI の静的ファイルを配置する。`testdata/` は責務別 fixture を配置する。`docs/examples/` は利用例、設定例、サンプル構成を配置する。

標準構成は移行後の最終形を示す。仕様化済み・未実装または将来計画の path は、該当 owner component が実装対象になった PR で追加する。標準構成に含まれることだけを理由に、未実装ファイル、将来計画ファイル、空ディレクトリ、placeholder を作成してはならない。

## 4c. Go 正本策定方針

本仕様は、Adlaire CI を最初から Go 言語で設計・実装する前提で策定する。

実装者は、`ADLAIRE_CI_SPEC.md` と `ADLAIRE_CI_DETAIL_SPEC.md` に記載された Go 仕様を正とする。過去の実装、試作、他言語スクリプト、既存ファイル名、既存 CLI、既存ログ、既存生成物、既存状態ファイルを前提にしてはならない。

Go 実装の判断基準は以下とする。

- `components/builder.go`、`components/runner.go`、`components/api.go` を Go 正本コンポーネントとして扱う。
- `adlaire-ci-build`、`adlaire-ci-runner`、`adlaire-ci-api` を標準実行バイナリ名とする。
- 仕様未記載の自動変換処理、暗黙の読み替え処理、仕様外分岐を実装判断で追加してはならない。
- 本仕様に記載されていない挙動は、仕様対象外として扱う。

## 4a. 詳細仕様方針

`ADLAIRE_CI_DETAIL_SPEC.md` は、実装者が追加判断なしに実装へ着手できる粒度で記載する。

詳細仕様は、抽象的な方針や目的の再掲ではなく、実装時に必要な具体値、処理順序、入出力、状態、失敗時の扱いを定義する。

仕様化済み・未実装の項目であっても、実装予定として扱う場合は実装者が迷わない粒度まで詳細化する。実装時期、設計判断、具体値が未確定の内容は、実装可能な仕様として扱わず、未仕様化または将来計画として明示する。

詳細仕様は、少なくとも以下の問いに答えられる状態を維持する。

- どのコンポーネントが責務を持つか
- どのファイル、API、関数、設定値、状態ファイルを使用するか
- 入力、出力、データ構造、既定値、許容値は何か
- 正常系の処理順序は何か
- 異常系、再試行、ロック、冪等性、タイムアウトをどう扱うか
- セキュリティ上の禁止事項、保存してよい情報、保存してはならない情報は何か
- 実装完了をどの確認条件で判定するか

## 4b. 仕様成熟度方針

仕様項目は、実装可否を判断できるように成熟度を明確に区分する。

成熟度は、構想の有無ではなく、実装者が実装に着手できるだけの情報が揃っているかで判定する。

将来計画や未確定アイデアは、実装対象として扱わない。実装対象にする場合は、先に詳細仕様を整え、仕様化済み・未実装へ昇格させる。

実装中に仕様不足、実装者判断に依存する分岐、未定義の入出力、未定義の状態ファイル、未定義の異常系を発見した場合は、実装判断で補完せず、仕様改訂へ戻す。

## 4c. Phase 実装単位方針

Adlaire CI の実装順序、実装計画、実装 PR、完了判定は Phase 単位でのみ管理する。

Phase は、対象 owner component、実装範囲、依存条件、完了条件、検証条件が明確な実装単位である。Phase の一覧、順序、依存条件、完了条件は `ADLAIRE_CI_DETAIL_SPEC.md` §0g を正とする。

`P0`、`P1`、`P2〜P5` などの優先度ラベル、抽象段階、API 内部分類、fixture 分類を、実装単位、PR 単位、完了判定単位として扱ってはならない。

API の実装範囲は Phase 3 と Phase 4 に分けて扱う。Phase 3 は API 基盤、認証、状態 read/write、運用基本操作を対象とする。Phase 4 は API 拡張運用操作を対象とする。

将来計画、改訂予定、未仕様化の機能は Phase に含めない。対象機能を Phase に含める場合は、先に詳細仕様、検証条件、受け入れ条件を整え、`仕様化済み・未実装` へ昇格させる。

## 4d. 実装着手ゲート方針

実装者は、対象コンポーネントについて以下をすべて満たすまで実装に着手してはならない。

実装順序、実装計画、実装 PR、完了判定は、本ファイル §4c の Phase 実装単位方針と Part 2 §0f の Phase 実装単位ポリシーに従う。

| 判定項目 | 着手条件 |
|----------|----------|
| 正本確認 | `ADLAIRE_CI_SPEC.md` と `ADLAIRE_CI_DETAIL_SPEC.md` の該当節を確認済みである。 |
| 状態分類 | 対象が `仕様化済み・未実装` または `実装済み` に分類され、`未仕様化`、`将来計画`、`改訂予定` ではない。 |
| 詳細節対応 | 対象機能が `ADLAIRE_CI_DETAIL_SPEC.md` §0i.1〜§0i.4 の詳細節対応表に含まれ、該当する詳細仕様節と受け入れ条件を確認済みである。 |
| テンプレート充足 | 対象機能の詳細仕様が `ADLAIRE_CI_DETAIL_SPEC.md` §0h の機能仕様テンプレートに必要な項目を満たしている。 |
| 責務境界 | 対象コンポーネント、呼び出し元、呼び出し先、状態ファイル、外部接続先が明確である。 |
| 契約同期 | API、SDK、UI、状態ファイル、セットアップ、検証条件のうち関係する仕様が同時に整合している。 |
| 禁止事項 | 外部依存、秘密情報、直接 API 呼び出し、互換処理、暗黙フォールバックなどの禁止事項が明確である。 |
| 完了条件 | `ADLAIRE_CI_DETAIL_SPEC.md` §0e、§0f、§0g、§0h、§0i.1〜§0i.4、§0j および §26.7 の受け入れ条件で、実装完了、実装順序、配置を判定できる。 |

着手条件を満たさない場合、実装者はコードで補完せず、先に仕様改訂を行う。実装 PR では、着手前に参照した詳細仕様節を PR 本文へ明記する。

## 4e. 完了判定方針

実装完了は、実装ファイルの作成やテスト成功だけでは成立しない。以下をすべて満たした場合にのみ `実装済み` と扱う。

- 実装が `ADLAIRE_CI_DETAIL_SPEC.md` の入力、出力、状態、処理順序、異常系、セキュリティ制約と一致している。
- `ADLAIRE_CI_DETAIL_SPEC.md` §0e の完全実装検証マトリクスと §0f の仕様策定完了チェックを満たしている。
- `ADLAIRE_CI_DETAIL_SPEC.md` §0g の初期実装 Phase 分割に従い、対象 Phase の依存条件、完了条件、PR 分割条件を満たしている。
- 対象機能が `ADLAIRE_CI_DETAIL_SPEC.md` §0h の機能仕様テンプレートを満たし、§0i.1〜§0i.4 の詳細節対応表に記載された受け入れ条件を満たしている。
- セットアップまたは運用手順に影響する場合、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.7 の実装受け入れ条件を満たしている。
- API、SDK、UI のいずれかを変更した場合、§22、§23、§24 の対応関係が崩れていない。
- 実装状態表、`DOCUMENT_INDEX.md`、本ドキュメント、詳細仕様の更新要否を確認済みである。
- 実装 PR 本文に、対象、実行コマンド、期待結果、実結果、判定を記録している。

検証不能な項目、未実行の項目、環境都合で省略した項目が残る場合、そのコンポーネントを `実装済み` として扱ってはならない。

## 4f. 仕様策定単位方針

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

### — CI ランナー —

## 5. CI ランナーの目的

GitHub API（Git Blobs API）を定期的にポーリングし、対象ファイルの変更を検出してビルドパイプラインを自動実行する自己ホスト型 CI ランナー（`components/runner.go`）。GitHub Actions・Webhook・外部 CI サービスへの依存をゼロにする。

- GitHub Git Trees API で対象ファイルの blob SHA を取得し、前回 SHA と比較して変更を検出する
- 変更検出時のみ Git Blobs API でファイル本文を取得し、ビルドを実行する
- 外部公開エンドポイント・リバースプロキシ不要
- Go 標準ライブラリを基本とし、外部依存を追加する場合は Part 2 §4 の例外承認を必須とする

## 6. CI ランナーの開発方針

- **単一責務実装**：`components/runner.go` は CI ランナー責務に限定し、Markdown 変換と管理 API を内包しない
- **シンプル性優先**：HTTP サーバー不要。1 回実行して終了する oneshot 設計
- **差分検出**：SHA キャッシュにより変更がない場合はビルドをスキップ

## 7. GitHub Actions との対応関係

| GitHub Actions | 内製ランナー |
|---|---|
| `.github/workflows/*.yml` | `.ci/pipeline.sh` |
| `runs-on: ubuntu-latest` | 自前サーバー（固定） |
| `steps:` の各ステップ | shell の各コマンド |
| `secrets.*` | PAT ファイル（`.github_token`）・サーバー上の環境変数 |
| Actions Marketplace | なし（自前実装のみ） |
| push トリガー | systemd タイマーによる定期ポーリング + SHA 差分検出 |

### — 管理ツール —

## 8. 管理ツールの目的

Adlaire CI の状態確認・操作を行う管理インターフェース。ヘッドレスアーキテクチャにより、フロントエンドとバックエンドを明確に分離する。

本節以降の管理ツール・管理 API・SDK に関する記載は、仕様化済み・未実装の内容である。対象ファイルの存在、実装状態表、詳細仕様、検証結果が実装済みとして整合するまでは、実装済み機能として扱わない。

## 9. ヘッドレスアーキテクチャ方針

Adlaire CI と管理ツールは API を介して通信する。フロントエンドとバックエンドを完全に分離し、管理ツールの実装・置き換えをバックエンドから独立させる。

- バックエンド（Adlaire CI）は API を公開する
- フロントエンド（管理ツール）は API のみを通じてバックエンドと通信する
- 直接のファイル操作・プロセス呼び出しは管理ツールから行わない

## 10. SDK 方針

API は SDK として提供し、管理ツール実装者が直接 HTTP 通信を記述しなくてよい抽象化レイヤーを提供する。

- **初期対応言語**：JavaScript のみ
- **フレームワーク非依存**：バニラ JS・React・Vue・Svelte 等、いずれの環境でも利用可能
- **内製 SDK**：外部ライブラリへの依存はゼロ（→ Part 2 §4）
- 標準管理ツールも本 SDK を経由して通信する

## 11. 標準管理ツール方針

Adlaire CI はすぐに使える標準管理ツールを同梱する。

- **実装技術**：HTML / CSS / JavaScript（バニラ）。外部フレームワーク不使用
- **SDK 経由**：バックエンドとの通信はすべて SDK を介する
- **カスタマイズ基盤**：標準管理ツールをベースとしたカスタマイズを前提とした設計とする。上書き・差し替えが容易な構造を維持する

---

## 12. 機能一覧

本仕様が定義する全機能の一覧。各機能の仕様詳細は `ADLAIRE_CI_DETAIL_SPEC.md` を参照。実装状態は、本ドキュメント冒頭の「実装状態」と本ファイル Part 2 §0a の仕様成熟度ポリシーに従って判定する。

### ビルド・CI ランナー（components/runner.go）

**Go 版で実装済みの初期範囲：**

- `--state-dir`、`--once`、`--version`、`--help` の CLI 契約
- `.branch_config` と `.last_sha` による branch target / SHA cache 読み込み
- GitHub Trees API / Blobs API による対象 Markdown 取得
- SHA 一致時の変更なし skip
- `pipeline.sh` 起動、stdout/stderr 収集、`[REPORT]` / `[WARN]` 取り込み
- `.build_logs/{id}.json`、`.build_history`、`.build_status.json`、`.build_state`、`.build_lock` の作成・更新
- deploy 失敗時の `.pending_transfers` 追加
- `.notify_pending` 破損時の退避と `[]` 再生成
- GitHub API retry / rate limit 待機
- `adlaire-ci-build` 実行可否と disk 空き容量の precheck
- `.snapshots/{build_id}/site` 保存と世代 pruning
- `.notify_pending` の HTTP 再送と成功時削除
- `.build_circuit_state` による circuit open skip と失敗回数記録
- `.pending_transfers` の起動時再試行と成功時削除
- SSH 転送時の checksum 比較、未変更ファイル skip、転送後 checksum 検証
- 複数 branch target の順次処理
- `BUILD_COOLDOWN_SECONDS` による cooldown skip
- `FORCE_BUILD_INTERVAL` による変更なし時の定期強制ビルド
- GitHub commits API によるトリガー commit 情報の build log 記録
- GitHub PAT 有効期限ヘッダーの 7 日以内 WARN ログ
- `.notify_config` に基づく成功・失敗・転送失敗 Webhook 通知送信
- `OUTPUT_SIZE_WARN_MB` による出力サイトサイズ警告
- Phase 2 fixture R1〜R7 と完了判定パステスト

**Go 版で仕様化済みの全体範囲：**

- GitHub リポジトリの対象ファイルを定期ポーリング（systemd timer）
- blob SHA による差分検出（変更なし時はビルドをスキップ）
- Markdown → 静的 Web サイト変換（`adlaire-ci-build` を呼び出し）
- ビルド成功後に SHA キャッシュを更新する
- ビルド失敗時は SHA キャッシュを更新せず、次回起動時に再試行可能な状態を維持する
- systemd oneshot ユニットとして動作（`adlaire-ci.service`）

- ビルド結果を `.build_history` に記録（ID 形式：`b{YYYYMMDDHHmmss}`）
- ビルドごとのログを `.build_logs/{id}.json` に保存
- ビルド成功・失敗時に Webhook 通知を送信（`.notify_config` を読み込み送信。送信責務は `components/runner.go`。通知 API は設定の読み書きのみ）
- ビルド成功後、出力サイトディレクトリを SSH 経由で静的コンテンツ配信サーバーへ転送する（差分転送・`DEPLOY_TARGETS` 複数先対応 → §14a）
- SSH 転送失敗時は `.pending_transfers` へキューイングし、次回起動時に自動再試行する（→ §14a ペンディングキュー）
- 転送成功後、出力サイトディレクトリを `.snapshots/` へアーカイブし `HISTORY_KEEP_N` 世代を超過分から自動削除する（→ §14b）
- SSH 転送失敗時に Webhook 通知を送信する（`on: ["deploy_failure"]` 設定時）
- `BRANCH_TARGETS` リストで複数ブランチを順次ポーリング・ビルドする（→ §12 設定値）
- `FORCE_BUILD_INTERVAL` 設定時、前回ビルドから指定時間経過で変更なしでも強制ビルドする

### 管理 API エンドポイント（components/api.go）

| カテゴリ | エンドポイント |
|---------|--------------|
| 認証 | `POST /api/login` / `POST /api/logout` / `POST /api/change-password` |
| 死活監視 | `GET /api/health` |
| ビルド操作 | `POST /api/build` / `POST /api/build/force` / `POST /api/build/cancel` / `GET /api/build/stream` / `POST /api/circuit-breaker/reset` |
| ステータス | `GET /api/status` / `GET /api/dashboard` |
| ログ | `GET /api/logs` / `GET /api/logs/export` / `GET /api/logs/search` / `POST /api/logs/cleanup` |
| ビルド履歴 | `GET /api/history` / `GET /api/history/{id}/log` / `GET /api/history/{id}/comment` / `POST /api/history/{id}/comment` / `GET /api/history/export` / `POST /api/history/{id}/flag` / `POST /api/history/{id}/tags` / `POST /api/history/{id}/rollback` |
| スケジュール | `GET /api/schedule` / `POST /api/schedule/interval` / `POST /api/schedule/pause` / `POST /api/schedule/resume` / `POST /api/schedule/allowed-hours` / `POST /api/schedule/force-interval` / `POST /api/schedule/cooldown` |
| 通知 | `GET /api/notify-config` / `POST /api/notify-config` / `POST /api/notify-test` / `GET /api/notify-log` / `POST /api/notify/weekly-summary` |
| システム情報 | `GET /api/sysinfo` / `GET /api/output-meta` / `GET /api/diagnostics` / `GET /api/rate-limit` / `GET /api/disk-usage` |
| 統計 | `GET /api/stats` / `GET /api/stats/timeline` / `GET /api/stats/build-duration` |
| リポジトリ | `GET /api/repo-info` / `POST /api/repo-config` / `GET /api/branch-config` / `POST /api/branch-config` |
| PAT 管理 | `GET /api/pat-status` / `POST /api/pat-verify` / `POST /api/pat-update` |
| 設定 | `GET /api/config` / `POST /api/config` / `POST /api/log-level` / `GET /api/config-log` |
| アクセスログ | `GET /api/access-log` |
| バックアップ | `GET /api/backup` / `POST /api/restore` |
| セッション管理 | `GET /api/sessions` / `POST /api/sessions/revoke-all` |
| API トークン | `GET /api/tokens` / `POST /api/tokens` / `DELETE /api/tokens/{id}` |
| スナップショット | `GET /api/snapshots` / `GET /api/snapshots/{id}/download` / `DELETE /api/snapshots/{id}` |
| メンテナンス | `GET /api/maintenance` / `POST /api/maintenance/enable` / `POST /api/maintenance/disable` |
| アクセス制御 | `GET /api/access-control` / `POST /api/access-control` |
| フック | `GET /api/hooks` / `POST /api/hooks` / `DELETE /api/hooks/{id}` / `GET /api/hooks/{id}/log` |
| アラートルール | `GET /api/alert-rules` / `POST /api/alert-rules` / `DELETE /api/alert-rules/{id}` |
| Webhook 受信 | `POST /api/webhook` / `GET /api/webhook-events` / `GET /api/webhook-config` / `POST /api/webhook-config` |
| 自動タグ付け | `GET /api/tag-rules` / `POST /api/tag-rules` / `DELETE /api/tag-rules/{id}` |
| パイプライン | `GET /api/pipeline-config` / `POST /api/pipeline-config` / `POST /api/verify-output` |
| 運用ノート | `GET /api/notes` / `POST /api/notes` |
| メール通知 | `GET /api/smtp-config` / `POST /api/smtp-config` / `POST /api/smtp-test` |
| ビルドキュー | `GET /api/queue` / `DELETE /api/queue` |
| ダッシュボードレイアウト | `GET /api/dashboard-layout` / `POST /api/dashboard-layout` |

### SDK メソッド（adlaire-ci-sdk.js）

| カテゴリ | メソッド |
|---------|--------|
| 認証 | `login()` / `logout()` / `changePassword()` |
| ビルド操作 | `triggerBuild()` / `buildForce()` / `cancelBuild()` / `streamBuild()` / `resetCircuitBreaker()` |
| ステータス | `getStatus()` / `getDashboard()` |
| ログ | `getLogs()` / `exportLogs()` / `searchLogs()` / `cleanupLogs()` |
| ビルド履歴 | `getHistory()` / `getHistoryLog()` / `getHistoryComment()` / `setHistoryComment()` / `exportHistory()` / `setHistoryFlag()` / `setHistoryTags()` / `rollbackHistory()` |
| スケジュール | `getSchedule()` / `setScheduleInterval()` / `pauseSchedule()` / `resumeSchedule()` / `setAllowedHours()` / `clearAllowedHours()` / `setForceInterval()` / `setBuildCooldown()` |
| Webhook 受信 | `getWebhookEvents()` / `getWebhookConfig()` / `setWebhookConfig()` |
| 通知 | `getNotifyConfig()` / `setNotifyConfig()` / `notifyTest()` / `getNotifyLog()` / `notifyWeeklySummary()` |
| システム情報 | `getSysinfo()` / `getOutputMeta()` / `getDiagnostics()` / `getRateLimit()` / `getDiskUsage()` |
| 統計 | `getStats()` / `getStatsTimeline()` / `getStatsBuildDuration()` |
| リポジトリ | `getRepoInfo()` / `setRepoConfig()` / `getBranchConfig()` / `setBranchConfig()` |
| PAT 管理 | `getPatStatus()` / `patVerify()` / `updatePat()` |
| 設定 | `getConfig()` / `setConfig()` / `setLogLevel()` / `getConfigLog()` |
| アクセスログ | `getAccessLog()` |
| バックアップ | `backup()` / `restore()` |
| セッション管理 | `getSessions()` / `revokeAllSessions()` |
| API トークン | `getTokens()` / `createToken()` / `revokeToken()` |
| スナップショット | `getSnapshots()` / `downloadSnapshot()` / `deleteSnapshot()` |
| メンテナンス | `getMaintenance()` / `enableMaintenance()` / `disableMaintenance()` |
| アクセス制御 | `getAccessControl()` / `setAccessControl()` |
| フック | `getHooks()` / `addHook()` / `deleteHook()` / `getHookLog()` |
| アラートルール | `getAlertRules()` / `addAlertRule()` / `deleteAlertRule()` |
| 自動タグ付け | `getTagRules()` / `addTagRule()` / `deleteTagRule()` |
| パイプライン | `getPipelineConfig()` / `setPipelineConfig()` / `verifyOutput()` |
| 運用ノート | `getNotes()` / `setNotes()` |
| メール通知 | `getSmtpConfig()` / `setSmtpConfig()` / `smtpTest()` |
| ビルドキュー | `getQueue()` / `clearQueue()` |
| ダッシュボードレイアウト | `getDashboardLayout()` / `setDashboardLayout()` |
| 死活監視 | `health()` |

ES Module・外部依存なし。全メソッドは `Promise` を返す。`streamBuild` は SSE 接続確立後に `Promise<StreamHandle>` として resolve し、`StreamHandle` は `{ close(): void, closed: boolean }` を持つ。`constructor` を除く合計は 96 メソッド。

### 標準管理ツール パネル（admin/index.html）

| パネル | 主な機能 |
|-------|---------|
| ログイン | パスワード認証 |
| パスワード変更 | 強制変更フロー対応（5 回目以降は他パネルを非表示） |
| ステータス | 最終ビルド情報・出力サイトリンク・メンテナンスバナー表示（モード中） |
| 手動実行 | ビルド起動・強制ビルド・キャンセル・キュー状態表示・キューのクリア |
| ログビューア | ログ閲覧・キーワードフィルター・ログレベルフィルター・JSON エクスポート・横断検索（期間指定） |
| ビルド履歴 | 過去ビルド一覧・タグ列・フラグ列・ページネーション・タグ/フラグフィルター・JSON エクスポート・個別ログ参照 |
| システム情報 | ファイルサイズ・稼働時間・ディスク使用量・PAT 検証・PAT 更新フォーム・PAT 有効期限表示・GitHub API レート制限表示 |
| 通知設定 | 複数 Webhook 設定・ペイロードテンプレート編集・テスト送信・定期サマリー設定・送信履歴・メール通知設定（SMTP連携）・Webhook 署名シークレット設定 |
| 設定 | 保持行数・件数・タイムアウト・ログレベル変更・ログ保持期間（日数）・手動クリーンアップ・設定変更履歴・キュー最大サイズ設定・スナップショット保持世代数設定 |
| アクセスログ | ログイン履歴（日時・成否） |
| 統計 | 成功率・平均/最大ビルド時間・時系列グラフ |
| リポジトリ情報 | 監視設定確認・ポーリング間隔変更フォーム・ポーリング一時停止/再開・許可時間帯設定・強制再ビルド間隔設定・Webhook 受信 Secret 設定・ブランチターゲット設定（`.branch_config` 編集） |
| セッション管理 | セッション一覧・全セッション強制無効化 |
| システム診断 | PAT・GitHub API・ファイル・systemd・Webhook・出力整合性 一括診断・アラートバッジ表示・メンテナンスバナー表示 |
| ビルド比較 | ビルド履歴から 2 件を選択してログを並列差分表示 |
| API トークン管理 | 読み取り専用トークンの発行・一覧・失効 |
| スナップショット | ビルド成果物の世代一覧・ダウンロード・削除・ロールバック |
| メンテナンス | メンテナンスモードの有効化・解除・状態表示 |
| アクセス制御 | 許可 IP / CIDR 一覧・追加・削除 |
| フック | Pre/Post ビルドフック設定・実行ログ確認 |
| 運用ノート | Markdown 記述の運用メモ閲覧・編集 |

---

## 13. 拡張ポイント・将来計画

本仕様に対する実装済み項目、実装中・検証未完了項目、仕様化済み・未実装項目、改訂予定項目、将来計画項目を統合管理するロードマップ。
仕様化する際は `ADLAIRE_CI_SPEC.md`、`ADLAIRE_CI_DETAIL_SPEC.md` の対応表、該当 owner component の詳細仕様本文への追記を先行させる。

将来計画は、実装対象ではない。将来計画内の「検討」「予定」「候補」「推奨」は、実装可能な仕様を意味しない。将来計画を実装対象にする場合は、先に対象項目を `改訂予定` へ昇格し、`ADLAIRE_CI_DETAIL_SPEC.md` の対応表と該当 owner component の詳細仕様本文に実装可能な詳細仕様を追加したうえで `仕様化済み` とする。

**担当領域：** `CI ランナー` / `管理ツール・API` / `MCP サーバー` / `ビルドスクリプト`

### 13.1 統合ロードマップの読み方

13章の項目は、単一の統合ロードマップ表で管理する。実装可否は `状態` 列で判断し、担当領域や機能名だけで実装対象と判断してはならない。

| 状態 | 実装可否 | 意味 | 次アクション |
|------|----------|------|--------------|
| 実装済み | 完了済み | ソースコード実装と検証が完了した項目。 | 実装ファイル、検証結果、`DOCUMENT_INDEX.md` を維持する。 |
| 実装中・検証未完了 | 検証待ち | ソースコード実装に着手済みだが、必須検証が未完了の項目。 | 必須検証を実行し、不足が残る場合は `実装済み` へ移動しない。 |
| 改訂予定 | 実装不可 | 将来計画から格上げ済みだが、詳細仕様作成中の項目。 | `ADLAIRE_CI_DETAIL_SPEC.md` の対応表と該当 owner component の詳細仕様本文を作成し、仕様化条件を満たす。 |
| 仕様化済み・未実装 | 実装可 | 正本仕様と詳細仕様があり、実装対象として扱える項目。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1〜§0i.4 と該当 owner component の詳細仕様本文を確認して実装する。 |
| 将来計画 | 実装不可 | `components/mcp.go` など、将来構想として管理する項目。 | 本節 13.3 の手順で `改訂予定` へ昇格する。 |

将来計画、改訂予定の項目は、実装着手可能な仕様ではない。実装対象にする場合は、先に 13.3 の手順で `仕様化済み・未実装` へ昇格させる。

---

### 13.2 統合ロードマップ表

本表は、§13 の全項目を状態別に統合した唯一の一覧である。項目を追加、削除、昇格、実装完了する場合は、本表の `状態`、`実装可否`、`次アクション` を同時に更新する。

MCP サーバー領域の行は、現時点ではすべて将来構想例であり、実装契約、API 契約、状態ファイル契約、起動手順、検証条件を定義しない。`components/mcp.go`、MCP tools、MCP resources、MCP prompts、HTTP SSE transport、MCP audit / stats / config CRUD は、MCP 専用詳細仕様を新設し、`改訂予定` を経て `仕様化済み・未実装` へ昇格するまで実装してはならない。

| 状態 | 実装可否 | 担当領域 | 機能 | 概要 | 次アクション |
|------|----------|----------|------|------|--------------|
| 実装済み | 完了済み | CI ランナー | Phase 2 完了判定パス | GitHub API polling、SHA 差分検出、ビルド起動、ログ、履歴、snapshot、lock、precheck、retry、rate limit、circuit breaker、通知、転送、cooldown、force interval、commit info、PAT 期限警告、出力サイズ警告を `components/runner.go` で実装済み。 | Go test で Phase 2 fixture、hardening、完了判定パスを検証済み。 |
| 改訂予定 | 実装不可 | 全領域 | （なし） | 現時点で、将来計画から格上げ済みの仕様作成中項目はない。 | 格上げ時に元状態、格上げ日、整理順序（実装単位ではない）、詳細仕様作成先を概要へ記録する。 |
| 実装済み | 完了済み | CI ランナー | ビルドタイムアウト | Go 標準ライブラリ `context.WithTimeout` と `os/exec` で長時間ビルドを強制終了する。API 経由の `build_timeout_seconds` 動的変更は管理 API 実装対象として残す（→ §22 `GET /api/config`）。 | Go test と runner 回帰検証で pipeline 起動経路を検証済み。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ポーリング間隔の動的変更 | systemd タイマーの `OnUnitActiveSec` を変更して間隔を調整（→ §22 `POST /api/schedule/interval`） | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.3、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26、`ADLAIRE_CI_DETAIL_API_SPEC.md` §27.11 に従って実装する。 |
| 実装済み | 完了済み | CI ランナー | ビルドログのファイル保存 | `os/exec` で起動したビルドプロセスの stdout/stderr を `.build_logs/{id}.json` に記録（→ §11 ファイル構成）。 | Go test で build log 生成と pipeline stdout/stderr 記録を検証済み。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | GitHub Webhook 受信 | 定期ポーリングと併用可能な即時検出方式。`POST /api/webhook` で GitHub push イベントを受信し即時ビルドをトリガーする（HMAC-SHA256 署名検証付き → §22） | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.3、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22-W、`ADLAIRE_CI_DETAIL_API_SPEC.md` §27.12 に従って実装する。 |
| 実装済み | 完了済み | CI ランナー | ネットワーク断時の再試行 | GitHub API 失敗時に指数バックオフ（`API_RETRY_BASE_SECONDS × 2^n`、最大 `API_RETRY_MAX` 回）で再試行する（→ §12・§13）。 | Go test で fake GitHub 一時失敗からの retry 成功を検証済み。 |
| 実装済み | 完了済み | CI ランナー | GitHub API レート制限自動待機 | `X-RateLimit-Remaining: 0` 検出時に `X-RateLimit-Reset` まで待機してから再試行する（→ §13）。 | Go test 対象の retry 経路と同じ GitHub API retry 実装で検証済み。 |
| 実装済み | 完了済み | CI ランナー | 転送後リモート整合性検証 | SSH 転送後に `sha256sum` でリモートファイルを検証し、不一致時はペンディングキューへ再投入する（→ §13・§14a）。 | Go test で転送失敗時 pending 化と pending 再試行成功を検証済み。 |
| 実装済み | 完了済み | CI ランナー | マルチブランチビルド | `BRANCH_TARGETS` リストで複数ブランチを順次ポーリング・ビルド・転送する（→ §12 設定値・§13 処理フロー）。 | Go test で branch target 設定に基づく処理経路を検証済み。 |
| 実装済み | 完了済み | CI ランナー | ビルドログ世代管理 | `LOG_KEEP_N` 件を超えた `.build_logs/{id}.json` を古いものから自動削除する（→ §12・§13）。 | `go test ./...` で runner 回帰検証済み。 |
| 実装済み | 完了済み | CI ランナー | ビルド出力の外部転送 | ビルド成功時に生成静的 Web サイトを SSH 経由（差分転送・複数ファイル対応）で静的コンテンツ配信サーバーへ自動転送する（→ §14a）。 | Go test で SSH command 差し替えによる pending / retry 経路を検証済み。 |
| 実装済み | 完了済み | CI ランナー | ビルドクールダウン | 前回ビルド完了から `BUILD_COOLDOWN_SECONDS` 秒以内の起動はビルドをスキップする（Webhook 二重トリガー防止 → §12・§13）。 | Go test で cooldown 中の build skip を検証済み。 |
| 実装済み | 完了済み | CI ランナー | ビルド前の事前チェック | `pipeline.sh` 実行前にディスク空き容量・`adlaire-ci-build` 実行可否・`components/builder.go` 由来のビルドバイナリ配置を確認し、不足時は `failure_precheck` として記録する（→ §13）。 | Go test で fake binary 不在時の `failure_precheck` を検証済み。 |
| 実装済み | 完了済み | CI ランナー | 定期強制ビルド | `FORCE_BUILD_INTERVAL`（時間単位）設定時、変更なしでも前回ビルドから経過時間超過で強制ビルドする。API 経由の動的変更は管理 API 実装対象として残す（→ §12・§13）。 | Go test で同一 SHA かつ force interval 超過時の build 実行を検証済み。 |
| 実装済み | 完了済み | CI ランナー | ビルド中重複スキップ | `.build_lock` に PID を記録し、起動時に実行中ビルドを検出したらスキップする（→ §11・§13）。 | Go test で lock conflict 時の build skip を検証済み。 |
| 実装済み | 完了済み | CI ランナー | GitHub PAT 有効期限の事前警告 | GitHub API レスポンスの `GitHub-Authentication-Token-Expiration` ヘッダーを解析し、7 日以内の期限切れを WARN ログで通知する（→ §13）。 | Go test で GitHub API response header 処理を検証済み。 |
| 実装済み | 完了済み | CI ランナー | コミット情報のビルドログ記録 | ビルドトリガーとなったコミットの SHA・メッセージ・作者名・コミット日時を `.build_logs/{id}.json` に記録する（→ §13）。 | Go test で fake commits API からの commit info 記録を検証済み。 |
| 実装済み | 完了済み | CI ランナー | GitHub API 連続失敗によるサーキットブレーカー | 連続失敗が `API_CIRCUIT_BREAKER_THRESHOLD` 周回以上になった場合に `.build_circuit_state.open=true` とし、open 中はポーリングをスキップする（→ §11・§12・§13）。 | Go test で circuit open 時の polling skip を検証済み。 |
| 実装済み | 完了済み | CI ランナー | 出力サイトサイズ警告閾値 | ビルド後の出力サイト合計サイズが `OUTPUT_SIZE_WARN_MB` を超えた場合に WARN ログを出力する。§8 変換レポートに `size_warn` フラグを追加（→ §8・§12・§13・§22）。 | Go test で閾値超過時の `size_warn=true` と WARN 記録を検証済み。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | Webhook イベントログ | 受信した Webhook push イベントを `.webhook_events.json` に JSON Lines 形式で追記記録する。`delivery_id`・`event`・`ref`・`sha`・`build_triggered` を保存（→ §11・§22） | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.3、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §11、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22-W、`ADLAIRE_CI_DETAIL_API_SPEC.md` §27.13 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド所要時間の記録と統計 API | `.build_logs/{id}.json` に `started_at`・`finished_at`・`duration_seconds` を記録し、`GET /api/stats/build-duration` で過去 N 件の平均・最小・最大を提供する（→ §22） | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.4、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.14 に従って実装する。 |
| 実装済み | 完了済み | CI ランナー | ビルドアーティファクト世代管理 | `HISTORY_KEEP_N` 世代分を `.snapshots/` に自動保持し超過分を削除する。`POST /api/history/{id}/rollback` による再転送は API 実装対象として残す（→ §14b）。 | Go test で build 成功時の `.snapshots/{id}/site` 作成を検証済み。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | ビルドアーティファクト管理 | スナップショット一覧・ダウンロード・削除・ロールバック（→ §14b・§22 `POST /api/history/{id}/rollback`） | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.4、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §14b、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23、`ADLAIRE_CI_DETAIL_UI_SPEC.md` §24、`ADLAIRE_CI_DETAIL_ARCHIVE_SPEC.md` §27.15 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | ヘルスチェックエンドポイント | `GET /api/health` を拡充。最終ビルド時刻・最終ビルド結果・最終転送結果・稼働秒数を返す（→ §22） | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.3、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_API_SPEC.md` §27.16 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | Webhook イベント一覧取得 API | `.webhook_events.json` をページネーション付きで返す `GET /api/webhook-events` を追加する（→ §22） | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.3、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23、`ADLAIRE_CI_DETAIL_UI_SPEC.md` §24、`ADLAIRE_CI_DETAIL_API_SPEC.md` §27.13 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | ビルドログ重大度フィルター | 既存の `GET /api/logs/search` に `level=warn\|error` パラメータを追加し、重大度別に絞り込む（→ §22） | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.3、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23、`ADLAIRE_CI_DETAIL_UI_SPEC.md` §24、`ADLAIRE_CI_DETAIL_API_SPEC.md` §27.17 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 変換レポート出力 | ビルド完了後に変換統計（見出し数・テーブル数・コードブロック数・警告）を stdout 出力する。`components/runner.go` が取り込み `GET /api/output-meta` で参照可（→ §8・§22） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | シンタックスハイライト | コードブロックに言語別色分けを `assets/app.js` で適用する。対応言語：`python`・`bash`・`json`・`sql`・`ini`・`diff`（→ §7.8） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 本文内全文検索 | ビルド時に `assets/search-index.json` を生成し、`assets/app.js` の検索 UI と統合して本文ヒット箇所へジャンプ（→ §7.9） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | アンカーリンク自動検証 | 生成 HTML 内の `#anchor` リンクが実際の見出しスラグと一致するか検証し、不整合を `[WARN] BROKEN_LINK` として警告出力する（→ §4.3・§8） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | コードブロックの折りたたみ | 30 行超のコードブロックを初期折りたたみ。「全 N 行を表示」リンクで展開（→ §7.10） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 印刷スタイル（`@media print`） | サイドバー・ヘッダー・ボタン類を非表示、コードブロック展開、リンク URL 末尾表示（→ §6） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 静的 Web サイト出力 | Markdown ファイルまたは Markdown ディレクトリから `index.html`、ページ HTML、`assets/style.css`、`assets/app.js`、`assets/search-index.json` を生成する（→ §5） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | テーマコンポーネント | 初期テーマ `adlaire-default` の header / sidebar / breadcrumb / toc / search / footer / codeblock / table / pagination を内製テンプレートとして提供する（→ §5） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 外部リンクの自動処理 | 外部リンク（`http://`・`https://`）に `target="_blank" rel="noopener noreferrer"` を付与し、内部リンクと区別する（→ §4.3） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 読み取り進捗バー | スクロール位置に応じた 3px プログレスバーをページ上端に固定表示する（→ §7.13） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | コードブロックのコピーボタン | コードブロック右上にワンクリックコピーボタンを配置する（→ §7.6） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 見出しアンカーリンクコピー | ホバーで表示される `.hn-link` ボタンクリックでアンカー URL をクリップボードにコピー（→ §3・§7.11） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | TOC 開閉状態の永続化 | TOC グループの展開／折りたたみ状態を `localStorage` に保存し、リロード後も復元する（→ §7.3） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 見出しスラグ重複解決 | 同一テキストの見出しが複数存在する場合に 2 番目以降のスラグへ `-2`・`-3` を付与して一意にする。TOC・アンカーコピー・全文検索と整合させる（→ §4.5） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 前後章ナビゲーションボタン | h2 見出し単位で「← 前の章」「次の章 →」ボタンを各章末尾に静的生成する（→ §4.5・§5・§7.15） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 内部リンク整合性チェック | `[label](#anchor)` 形式の内部リンクが実際のスラグと一致するか変換時に検証し、不一致を `[WARN]` で報告。§8 変換レポートの `broken_links` フィールドに件数を記録する（→ §4.3・§8） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 見出し階層スキップ警告 | h1→h3 のような見出しレベルの 2 段以上のスキップを `[WARN]` で報告。§8 変換レポートの `heading_skips` フィールドに件数を記録する（→ §4.5・§8） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 読了時間推計と表示 | 本文文字数（コードブロック・タグ除く）から読了時間（分、200文字/分・切り上げ）を算出し、固定ヘッダーに静的埋め込みする。§8 変換レポートの `reading_time` フィールドに記録する（→ §4.5・§5・§6・§8） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | CI ランナー | Webhook 通知失敗リトライキュー | `.notify_pending`（JSON）を起動時に再送し、HTTP 2xx 成功時に削除、失敗時に `retry_count` と `last_error` を更新して保持する（→ §11・§13）。 | Go test で fake HTTP endpoint への再送成功と queue 空化を検証済み。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ブランチ設定の動的変更 API | `BRANCH_TARGETS` を外部 JSON（`.branch_config`）で管理し `GET /api/branch-config` / `POST /api/branch-config` で API 経由変更可能にする。`components/runner.go` 再起動不要（→ §11・§12・§22） | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.3、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §11、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §12、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_API_SPEC.md` §27.18 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | 週次ビルドサマリー Webhook | 指定曜日・時刻に過去 7 日間の成功率・平均ビルド時間・エラー件数をまとめた定期通知を送信する（→ §12・§13・§22） | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.3、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §12、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §16、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.19 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | 設定変更の詳細 diff 記録 | `.config_log` の各エントリに変更前後の値の diff 文字列を付加し `GET /api/config-log` レスポンスに含める（→ §22） | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.3、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_API_SPEC.md` §27.20 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | テーブルのソート機能 | 列ヘッダークリックで昇順/降順ソートができるインタラクティブテーブル。`aria-sort` 属性と CSS `::after` でインジケーター表示（→ §7.14） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | キーボードショートカット | `/` で検索フォーカス・`Escape` で検索クリア・`t` でページ先頭へスクロール（→ §7.12） | `ADLAIRE_CI_DETAIL_SPEC.md` §0h・§0i.1 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | 複数ファイル監視 | `target_files` で複数 Markdown ファイルまたは Markdown ディレクトリを監視し、変更対象ごとに build target を決定する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §11、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §12、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.21 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | GitHub Commit Status API | ビルド開始時・成功時・失敗時・pending 時に GitHub Commit Status API へ `pending` / `success` / `failure` を送信し、対象 commit に Adlaire CI の結果を紐付ける。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §12、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0c、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_COMMITSTATUS_SPEC.md` §27.1 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルドパイプライン YAML 定義 | 内製 YAML subset parser で `.pipeline.yml` を読み込み、build / test / deploy step を固定順で実行する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §12、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.22 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ドライラン実行モード | `adlaire-ci-runner --dry-run` で設定・状態・GitHub target・SHA 差分・起動可否を検証し、ビルド、deploy、通知、状態更新を行わず結果を出力する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §11、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §12、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.2 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルドログのアーカイブ圧縮 | 保持期間を超えた `.build_logs/{id}.json` を `.build_logs/archive/{id}.json.gz` に gzip 圧縮し、通常ログ API は圧縮済みログも透過的に参照する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.4、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §12、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0c、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_ARCHIVE_SPEC.md` §27.7 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ローカルファイル監視モード | GitHub API を使わず、ローカル状態 snapshot の SHA-256 差分で対象 Markdown の変更を検出する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §11、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §12、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.23 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | タグ付きコミットのみビルド | 監視対象 commit に許可 tag pattern が付いている場合だけ build を実行する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §12、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.24 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルドキャッシュ | 入力ファイル単位の SHA-256 manifest を `.build_cache.json` に保存し、未変更ページの変換結果を再利用する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.1、`ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §5、`ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §8、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §11、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §27.25 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド通知連携 | `components/runner.go` が `.notify_config.channels` に基づき Webhook / email / command 通知を同一通知イベント契約で送信する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §16、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.32 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド時間トレンド記録 | build ごとの所要時間を `.build_trends.json` に集計保存し、移動平均、中央値、p95、直近件数を API / UI から参照できるようにする。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.33 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド失敗時の自動リトライ | 一時的な GitHub API / network / pipeline timeout / deploy 検証失敗を対象に、設定回数まで同一 build id 内で再試行し、attempt ごとの結果を build log と history に記録する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §12、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0c、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.3 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | 並列マルチターゲットビルド | branch target 内の複数 deploy target を bounded worker で並列処理し、target ごとの結果を build log に記録する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §12、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §14a、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.26 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド前後フック | `.hooks` の pre / post hook を shell 経由なしで実行し、abort 条件と hook log を固定仕様どおり扱う。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.27 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | 依存ファイルトラッキング | Markdown 内の相対リンク・画像・include 対象を依存 manifest として記録し、関連 target だけを再ビルドする。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.1、`ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §4.3、`ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §5、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §11、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §27.28 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | リモートビルド対応 | SSH 経由で remote build command を実行し、成果物 archive と manifest を取得して既存 deploy / log 契約へ接続する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §12、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §14a、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.29 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルドステータスファイル出力 | runner の現在状態、最終ビルド、最終 deploy、pending 件数、circuit 状態、最終 trigger を `.build_status.json` に JSON object として出力し、API / UI / MCP の read-only 参照元にする。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §11、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0c、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.8 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド承認フロー | approval_required な target を `.approval_queue` に保留し、API 承認後だけ build / deploy を継続する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.3、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §11、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §16、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.30 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ブランチ別環境変数 | branch target ごとに許可済み環境変数を build process へ注入し、secret を log に出さない。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §12、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.31 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド依存チェーン | `.build_chain_config` で job 間依存を定義し、依存成功後だけ後続 job を実行する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §11、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.34 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド優先度キュー | `.build_state.queued` に priority / created_seq を保存し、優先度順かつ同一優先度 FIFO で build queue を処理する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §11、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.35 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | 失敗原因の自動分類 | stdout / stderr / exit code / runner error を固定分類ルールで解析し、failure_category と evidence を build log / history に記録する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.36 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド実行環境の記録 | build 開始時の OS、arch、Go version、binary version、disk usage、hostname を `.build_logs/{id}.json.environment` に記録する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.37 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルドトリガー種別の記録 | `polling`、`force_interval`、`manual`、`webhook`、`retry_pending_transfer`、`startup_config_integrity`、`rollback`、`local_watch`、`approval` を `.build_logs/{id}.json`、`.build_history`、`.build_status.json` に記録し、履歴 filter と状態表示で同じ値を使う。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0c、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23、`ADLAIRE_CI_DETAIL_UI_SPEC.md` §24、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.9、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.23、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.30 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド所要時間の異常検知 | `.build_trends.json` の移動平均と p95 を基準に異常に遅い build を検出し、WARN、history flag、通知へ反映する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §16、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.38 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | 設定ファイル起動時整合性チェック | runner 起動時に `.branch_config`、`.notify_config`、`.build_state`、`.pending_transfers`、`.notify_pending`、`.build_circuit_state` の JSON 整合性を検証し、破損・型不一致・必須 key 不足を規定どおり退避、初期化、通知、または停止する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.2、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §11、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §12、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0c、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.10 に従って実装する。 |
| 将来計画 | 実装不可 | 管理ツール・API | マルチユーザー対応 | 初期仕様の単一 admin 認証を、複数ユーザー・ユーザー別セッション・ユーザー別監査へ拡張する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 通知先の拡張 | 管理画面・API からメール・Slack・Discord 等の通知チャンネルを設定・追加できるようにする。`components/runner.go` 側のフック実装は → ビルド通知連携 | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | データストア切り替え | 大量ビルド履歴・ログ運用に備え、フラットファイルから SQLite 等への切り替えを検討する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 外部認証連携 | SSO・OAuth 等の外部認証基盤との連携を検討する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | TOTP 二要素認証 | RFC 6238 TOTP を Go 標準ライブラリだけで検証し、login を password verified / totp required / session issued の二段階に分離する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.4、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23、`ADLAIRE_CI_DETAIL_UI_SPEC.md` §24、`ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` §27.46 に従って実装する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 統計データの CSV エクスポート | `GET /api/stats/timeline` の日別データを CSV 形式でダウンロードできるエンドポイントを追加する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | キュー内個別エントリのキャンセル | `DELETE /api/queue/{id}` で特定エントリのみキャンセルする（初期仕様では全クリアのみ） | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | 設定バリデーション API | `POST /api/config/validate` で `.server_config` 互換の設定差分を保存前に検証し、正規化後設定、警告、エラー位置を返す。状態ファイルは変更しない。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.3、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0c、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23、`ADLAIRE_CI_DETAIL_UI_SPEC.md` §24、`ADLAIRE_CI_DETAIL_API_SPEC.md` §27.5 に従って実装する。 |
| 将来計画 | 実装不可 | 管理ツール・API | Prometheus メトリクスエンドポイント | `GET /api/metrics` で Prometheus 形式のメトリクス（ビルド数・成功率・ディスク使用量等）を返す | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | CLI 管理クライアント | Go 標準ライブラリのみで実装した `adlaire-ci-cli` で API を CUI 操作できるツール | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 設定の自動スナップショット | `POST /api/config` 変更時に自動で設定バックアップを世代保存する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ステータスバッジ生成 | `GET /api/badge` で最終ビルド結果を SVG バッジとして返す（README 埋め込み用） | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ビルド履歴の自動削除設定 | `history_retention_days` 設定でビルド履歴エントリを自動削除する（ログの `log_retention_days` に対応する履歴版） | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | セッションタイムアウト変更設定 | `.server_config.session_timeout_seconds` で session 有効期限を 5 分〜30 日の範囲で変更し、既存 session の扱いを固定する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.4、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0c、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23、`ADLAIRE_CI_DETAIL_UI_SPEC.md` §24、`ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` §27.45 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | ビルドトリガー専用 API スコープ | API token scope に `trigger` を追加し、build 起動系だけを許可する最小権限 token を発行できるようにする。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.4、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23、`ADLAIRE_CI_DETAIL_UI_SPEC.md` §24、`ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` §27.42 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | 監査ログ | `.audit_log` に設定変更、認証、token、build trigger、承認、権限拒否を actor 付き JSON Lines で記録する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.4、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23、`ADLAIRE_CI_DETAIL_UI_SPEC.md` §24、`ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` §27.44 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | API レート制限 | `.api_rate_state` で actor / IP / endpoint group ごとの固定窓 rate limit を管理し、超過時 `429` を返す。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.4、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23、`ADLAIRE_CI_DETAIL_UI_SPEC.md` §24、`ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` §27.47 に従って実装する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ロールベースアクセス制御 | 複数ユーザー対応後に、管理者・オペレーター・閲覧者等の役割ごとに API 権限を分ける | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 成果物ダウンロード API | 生成静的 Web サイトを archive として API エンドポイント経由で直接ダウンロードできるようにする | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 設定スナップショット差分表示 | 保存済みスナップショット間の設定変更点を diff 形式で確認できる API | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 複数プロジェクト管理 | 単一インスタンスで複数リポジトリ／プロジェクトを切り替え管理する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | API キー管理 | `.api_tokens` で API key の hash、scopes、expires_at、revoked_at を管理し、発行時だけ token 本体を返す。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.4、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23、`ADLAIRE_CI_DETAIL_UI_SPEC.md` §24、`ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` §27.43 に従って実装する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ビルドログのリアルタイム配信 | 実行中ビルドのログを SSE / WebSocket でストリーミング配信するエンドポイント | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 設定のインポート／エクスポート | 設定全体を JSON でエクスポートし、別環境へそのままインポートする | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ビルド統計ダッシュボード | 成功率・平均ビルド時間・エラー分布等を可視化する管理画面を生成する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ユーザー管理 API | 複数ユーザー対応後に、管理者アカウントの追加・削除・パスワード変更を API で操作する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | IP アドレス制限 | 管理 API へのアクセスを許可 IP レンジに限定するフィルタリング | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | API バージョニング | `/api/v1/` 等のバージョンプレフィックスで API 世代を明示する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | Webhook 署名検証 | 受信 Webhook の HMAC 署名を検証し、なりすましリクエストを拒否する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | API ドキュメント自動生成 | OpenAPI / Swagger 仕様を自動生成し、インタラクティブなドキュメントとして提供する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 通知チャンネル管理 | Slack / Discord / メール等の通知先を管理画面から追加・削除・テスト送信する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ビルドキューの手動並び替え | 管理画面からキュー内ジョブの実行順序を変更する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 設定テンプレート | よく使う設定パターンをテンプレートとして保存・再利用する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | API アクセスログ | 認証後 API と Webhook の呼び出しを `.api_access_log` に JSON Lines で記録し、`GET /api/api-access-log` でページング参照する。認証ログ `.access_log` とは分離する。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.3、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0c、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23、`ADLAIRE_CI_DETAIL_UI_SPEC.md` §24、`ADLAIRE_CI_DETAIL_API_SPEC.md` §27.6 に従って実装する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 管理者向けイベントフィード | ビルド完了・エラー・設定変更等のシステムイベントをリアルタイムで流す管理ページ | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ビルドキュー可視化 | キューに積まれたビルドの状態一覧を静的 HTML ステータスページとして出力する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | メンテナンスモード | 管理 API から即時にメンテナンスモードへ切り替え、ビルドキューを一時停止する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 読み取り専用共有リンク | ビルドステータス・統計を外部に公開する期限付き読み取り専用リンクを発行する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | アラート閾値設定 | ビルド失敗率・所要時間等が設定閾値を超えた際に自動アラートを発火する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | バックアップ／リストア | 設定・ビルド履歴・ログ等の全データをアーカイブ化してリストアできる機能 | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | API レスポンスキャッシュ制御 | 頻繁に参照される統計・ログ API のキャッシュ TTL を設定から変更する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | スナップショット間サイト差分 API | 2 つのスナップショット ID を指定し、出力サイトの追加/削除行数・変更率を返す `GET /api/snapshots/{id1}/diff/{id2}` を追加する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | Webhook 送信履歴の手動再送 API | `GET /api/notify-log` の各エントリに対して `POST /api/notify-log/{id}/retry` で同一ペイロードを即時再送できる手動リトライ API。`.notify_pending` 自動再試行とは別に特定通知だけ個別再送できる運用機能 | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 差分ビルド | 変更箇所のみ処理し、大規模 MD の変換を高速化する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 複数出力形式 | HTML に加えて PDF・ePub 等の出力形式をサポートする | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | Markdown 拡張記法サポート | アドモニション（`> [!NOTE]`）・カラーバッジ等の独自拡張記法に対応する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | コードブロック行番号表示 | コードブロック左端に行番号を表示するオプションを追加する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 見出しの自動採番 | h2 以下の見出しに `1.1`・`1.2` 等の番号を自動付与するオプション | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | セクション折りたたみ | 見出しクリックでコンテンツを折りたたむ機能（デフォルト展開） | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | TOC 深さ制御 | TOC に含める見出しレベルを設定で指定する（例：h2–h3 のみ） | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 最終更新日の自動埋め込み | ソースの git コミットタイムスタンプをフッターに自動出力する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | diff ハイライト | `+`/`-` で始まる行を git diff スタイルで緑/赤に色分けするコードブロックオプション | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 画像の遅延読み込み | `<img>` に `loading="lazy"` を付与し、初期表示を高速化する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | カスタムメタタグ注入 | OGP / Twitter Card 等のメタタグをビルド設定から生成・埋め込む | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | ダークモード対応 | `prefers-color-scheme` に応じたライト／ダーク切り替えを実装する（現在はライト固定） | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | コードブロックのファイル名表示 | ` ```go:filename.go ` や ` ```言語名:filename.ext ` 記法でコードブロック上部にファイル名ラベルを表示する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | テンプレート変数展開 | ビルド設定に定義した変数を `{{ VERSION }}` 等の記法で Markdown 本文中に展開する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | HTML ミニファイ | 生成 HTML のホワイトスペース・コメントを除去してファイルサイズを削減する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | TOC ハイライト追従 | スクロール位置に応じてサイドバー TOC の現在セクションを自動ハイライトする | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | Mermaid ダイアグラム描画 | ` ```mermaid ` コードブロックをフローチャート・シーケンス図として SVG 描画する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 脚注サポート | `[^1]` 記法の脚注をページ末尾に自動レンダリングする | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | インライン数式レンダリング | `$...$` / `$$...$$` 記法の数式を KaTeX 等で描画する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | ページ内ナビゲーション履歴 | ブラウザの戻る/進むに対応したハッシュベースの履歴管理を実装する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 読み上げ対応（アクセシビリティ） | `aria-label`・`role` 属性の付与対象、値、検証方法を詳細仕様で定義した上でスクリーンリーダー閲覧に対応する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 画像ライトボックス | 画像クリックでモーダル拡大表示する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | 出力サイトへのビルドメタ埋め込み | `adlaire-ci-build` が build id、commit SHA、build at を受け取り、生成 HTML の `<head>` に固定 meta として埋め込む。`GET /api/output-meta` は同値を返す。 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.1、`ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §2、`ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §5、`ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §8、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §13、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、`ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §27.4 に従って実装する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 印刷時 QR コード挿入 | `@media print` で元ページの URL を QR コードとしてフッターに埋め込む | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | 定義リストサポート | `term\n: definition` 記法を `<dl>/<dt>/<dd>` タグにレンダリングする | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | ビルドスクリプト | タスクリストサポート | `- [ ]` / `- [x]` 記法をチェックボックス付きリストとして描画する | 13.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP サーバー実装 | `components/mcp.go` を将来追加コンポーネントとして追加。MCP プロトコル（JSON-RPC over stdio）で Claude Desktop 等の AI クライアントから直接接続可能にする。内部では `components/api.go` REST API に Go 標準ライブラリ `net/http` でローカル接続するラッパー設計（`encoding/json` + `os.Stdin` / `os.Stdout` + `net/http`、ゼロ外部依存）。認証は `.mcp_token` に専用 API トークンを保存し、スコープ（`read` のみ / `trigger` 許可）をトークン単位で選択可能。Claude Desktop の `mcpServers` 設定に `/usr/local/bin/adlaire-ci-mcp` を指定して接続する | 13.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP ツール・リソース公開 | MCP サーバーが公開するツール：`get_status`（ビルド状態・CB 状態・PAT 残日数）/ `get_history(n)`（直近 N 件）/ `search_logs(query, level?, from?, to?)`（ログ全文検索）/ `get_build_log(id)`（個別ビルドログ）/ `trigger_build(force?)`（ビルドトリガー、`trigger` スコープ必須）/ `reset_circuit_breaker`（CB リセット、`trigger` スコープ必須）。リソース：`adlaire://status` / `adlaire://history` / `adlaire://logs/{id}` / `adlaire://config`（→ MCP サーバー実装） | 13.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | AI 支援ビルドエラー分析 | ビルド失敗時、AI クライアント（Claude Desktop 等）が `get_build_log` / `search_logs` ツールを自律的に呼び出してエラーログを取得し、原因推定と修正提案を生成できる設計。AI 側が pull するため `components/runner.go` のゼロ依存を完全維持。将来的には Webhook 通知をトリガーに AI が自動分析を開始する構成も検討可（→ MCP サーバー実装） | 13.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP Prompts 定義 | よく使う分析シナリオを MCP Prompts として定義し、Claude Desktop 等のプロンプトメニューから即時呼び出し可能にする。例：「先週のビルド失敗率をまとめて」「最後のエラーの原因を分析して」「PAT 有効期限が近いか確認して」。`get_history` / `search_logs` ツールと組み合わせて定型分析を 1 クリックで実行できる（→ MCP サーバー実装） | 13.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP Sampling によるビルドログ自動分析 | MCP Sampling 機能を使い、ビルド失敗時に `components/mcp.go` が AI クライアントへ sampling リクエストを送って原因推定テキストを生成し `.build_logs/{id}.json` の `ai_analysis` フィールドへ自動記録する。`components/runner.go` は MCP サーバーへソケット通知を送るだけで Anthropic API キーは `components/mcp.go` も保持しない（→ MCP サーバー実装） | 13.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP Notifications（イベントプッシュ） | ビルド完了・失敗・CB 開放等のシステムイベントを MCP Notifications として接続中の AI クライアントへリアルタイムプッシュする。`components/runner.go` が `components/api.go` 経由でイベントをキューに積み `components/mcp.go` が接続クライアントへ転送する設計（→ MCP サーバー実装） | 13.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP HTTP SSE transport 対応 | 初期仕様の stdio transport に加えて HTTP + SSE transport をサポートし、リモートマシンや複数クライアントからの同時接続を可能にする。Go 標準ライブラリ `net/http` で実装しゼロ依存を維持。`components/api.go` と同一プロセス統合か独立ポート起動かを設定で選択可能（→ MCP サーバー実装） | 13.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP ツールスコープ細分化 | 初期仕様の `read` / `trigger` 2 スコープを `read`（参照のみ）/ `trigger`（ビルド起動）/ `admin`（設定変更・CB リセット・ブランチ設定変更）の 3 スコープに細分化する。`.mcp_token` の各トークンにスコープを紐付け、ツール呼び出し時にスコープ検証を行う（→ MCP サーバー実装） | 13.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP ツール呼び出し監査ログ | MCP 経由で呼び出されたツールの履歴（呼び出し日時・ツール名・引数サマリー・成否）を `.mcp_access_log` に記録する。`GET /api/mcp-access-log` で参照可能にし、AI クライアントがどの操作をいつ実行したかを追跡できる（→ MCP サーバー実装） | 13.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP リソース購読（Resource Subscriptions） | クライアントが `adlaire://status` 等のリソースを subscribe し、状態変化時に `notifications/resources/updated` を自動受信できる MCP 標準機能。MCP Notifications（イベント起点プッシュ）とは異なりリソース変更起点のプッシュで、クライアントがポーリングなしに最新状態を保持できる（→ MCP サーバー実装） | 13.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP クライアント情報ログ | initialize リクエストの `clientInfo`（クライアント名・バージョン）を `.mcp_access_log` の接続エントリとして記録する。Claude Desktop / Cursor / 自作クライアント等どのツールから接続されたかを追跡し、監査と動作確認に利用する（→ MCP サーバー実装） | 13.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP ツール実行統計 | ツールごとの呼び出し回数・平均応答時間（ms）・エラー率を `.mcp_stats` に蓄積する。`GET /api/mcp-stats` で参照可能にし、どのツールが頻用されているか・ボトルネックがどこかを可視化する（→ MCP サーバー実装） | 13.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP ツール実行タイムアウト設定 | `.mcp_config` にツールごとのタイムアウト秒数を設定可能にする（例：`search_logs: 10`、`trigger_build: 30`）。タイムアウト超過時は JSON-RPC エラーを返し、MCP サーバーが無応答になる事態を防ぐ（→ MCP サーバー実装） | 13.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP 設定 CRUD ツール | `get_config(section?)` / `set_config(key, value)` ツールを `admin` スコープで公開する。AI クライアントから直接 `.server_config` / `.notify_config` 等を読み書きでき、チャット上で「ポーリング間隔を 30 秒に変更して」と指示するだけで設定変更が完結する（`components/runner.go` 再起動不要）。`admin` スコープの定義は → MCP ツールスコープ細分化（→ MCP サーバー実装） | 13.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP Elicitation による副作用操作の確認 | MCP Elicitation に対応し、`trigger_build(force=true)` / `reset_circuit_breaker` 等の副作用操作の実行前に AI クライアントへ確認プロンプトを送信して明示的な承認を得てから実行する。JSON-RPC メッセージの追加のみで実装しゼロ依存を維持。意図しない操作の誤実行を防ぐ安全機構（→ MCP サーバー実装） | 13.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |

### 13.3 状態変更・昇格手順

統合ロードマップ表の `将来計画` の項目を実装対象にする場合は、以下の順で進める。

1. 対象項目の `状態` を `改訂予定` に変更し、元状態、格上げ日、担当領域、整理順序（実装単位ではない）、ステータスを同じ行の概要または次アクションへ記録する。
2. `ADLAIRE_CI_SPEC.md` の該当する方針、ポリシー、実装状態、機能一覧を改訂する。
3. `ADLAIRE_CI_DETAIL_SPEC.md` §0i.1〜§0i.4 の詳細節対応表に、対象機能、対象コンポーネント、詳細仕様節、受け入れ条件を追加する。
4. 該当 owner component の詳細仕様本文に、§0h の機能仕様テンプレートを満たす目的、責務、入出力、状態、処理順序、異常系、セキュリティ、検証条件を追加する。MCP サーバー領域を昇格する場合は、MCP 専用詳細仕様を新設し、`ADLAIRE_CI_DETAIL_SPEC.md` の責務 component 対応表へ追加する。
5. `DOCUMENT_INDEX.md` と `AGENTS.md` の更新要否を確認する。
6. 仕様凍結条件を満たした後、対象項目の `状態` を `仕様化済み・未実装`、`実装可否` を `実装可` に変更する。
7. 実装と検証が完了した後、対象項目の `状態` を `実装済み`、`実装可否` を `完了済み` に変更する。

上記手順を完了していない項目は、実装者判断で実装してはならない。

---

# Part 2 — ポリシー
> 遵守義務のある規則と制約を定める。「何をしなければならないか／してはならないか」に答える。

## 0. 詳細仕様記載ポリシー

`ADLAIRE_CI_DETAIL_SPEC.md` の各仕様項目は、実装者が実装レベルで迷わない記載にしなければならない。

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
| 実装済み | 仕様に基づくコード変更が完了し、仕様との差分確認、必要な構文確認、実行確認または生成物確認、関連文書の整合確認が完了した状態。 | 完了済み |

`仕様化済み・未実装` へ昇格するには、少なくとも以下を満たさなければならない。

- 対象コンポーネントと責務境界が明確である
- 入力、出力、設定値、状態ファイル、データ構造が明確である
- 正常系の処理順序と分岐条件が明確である
- 異常系、ログ、通知、再試行、ロック、タイムアウトの扱いが必要範囲で明確である
- 認証、認可、秘密情報、外部公開可否などのセキュリティ制約が明確である
- 実装完了を判定する検証条件が明確である

実装着手は、対象項目が `仕様化済み・未実装` の状態に到達している場合に限る。

実装完了は、コード変更だけでは成立しない。仕様との差分確認、構文確認、実行確認または生成物確認、`DOCUMENT_INDEX.md`・実装状態表・関連仕様の更新要否確認を完了した場合にのみ `実装済み` と扱う。

API、SDK、標準管理ツールのいずれかを変更する場合は、API 仕様、SDK メソッド、UI 操作、詳細仕様の整合を同時に確認する。いずれか一方だけを変更して完了扱いにしてはならない。

## 0b. 仕様 PR 完了ポリシー

仕様策定または仕様改訂の Pull Request は、以下を満たすまで完了扱いにしてはならない。

| 対象 | 完了条件 |
|------|----------|
| マスター仕様 | 方針、ポリシー、実装状態、実装着手ゲート、完了判定が `ADLAIRE_CI_SPEC.md` に明記されている。 |
| 詳細仕様 | 実装に必要な具体値、入出力、状態、処理順序、異常系、検証条件が `ADLAIRE_CI_DETAIL_SPEC.md` に明記されている。 |
| API / SDK / UI | いずれかを変更した場合、endpoint、SDK method、UI 操作、エラー表示、再取得、秘密情報消去が同期している。 |
| セットアップ | バイナリ名、配置パス、systemd unit、権限、初期化順、失敗時停止条件、アップデート rollback 方針が同期している。 |
| 実装状態 | 仕様化済み・未実装、実装済み、将来計画の区分が矛盾していない。 |
| 索引 | ファイル名、正本関係、実装対象の変更がある場合、`DOCUMENT_INDEX.md` の更新要否を確認している。 |

仕様 PR は、未確定事項を「推奨」「検討」「適切に」等の表現だけで残してはならない。未確定事項を残す場合は、実装不可の `未仕様化` または `将来計画` として明示する。

## 0c. 仕様 PR 分割ポリシー

仕様 PR は、変更責務と merge 順序が明確な単位で作成する。競合防止を優先し、同一目的、同一仕様領域、同一ファイル群の仕様変更は必ず 1 本の PR にまとめる。

| 分類 | ルール |
|------|--------|
| 同一ファイル変更 | 同じ仕様領域で同一ファイルを編集する変更は、必ず 1 本の PR にまとめる。 |
| API / SDK / UI | endpoint、SDK method、UI 操作のいずれかを変更する場合、対応する仕様を同一 PR に含める。 |
| 状態ファイル | 状態ファイル schema、read/write 対応、破損時処理、権限、atomic write を同一 PR に含める。 |
| セットアップ | systemd、配置パス、権限、初期化、アップデート、rollback、受け入れ条件を同一 PR に含める。 |
| 将来計画 | 実装不可の将来計画は、実装可能仕様と同一表で混在させる場合でも状態を明示する。 |
| 別 PR 許可 | 変更対象ファイル、責務、merge 順序が完全に独立し、片方だけ merge されても仕様矛盾、参照切れ、状態不一致、未定義の依存関係が起きない場合のみ別 PR を許可する。 |
| 既存 PR 優先 | 既存 open PR と同じ仕様領域または同じファイルを変更する場合、新規 PR を作成せず既存 PR へ統合する。 |
| 競合 PR | `DIRTY`、同一ファイル編集、同一仕様領域、merge 順序依存の PR が存在する場合、先に 1 本の PR へ統合し、重複 PR を close する。 |

同一目的の仕様 PR を複数に分割してはならない。分割済みの PR 間で同一ファイルまたは同一仕様領域を編集している場合は、最新の作業ブランチへ統合し、1 本の PR にまとめる。

仕様 PR 作成前には、open PR 一覧、変更ファイル一覧、merge 状態を確認する。確認対象は、少なくとも `gh pr list --state open`、`git fetch origin`、`git diff --name-status origin/main...HEAD` とする。

競合解消後は、競合マーカーが残っていないこと、`git diff --check` が成功すること、open PR が同一仕様領域で重複していないこと、統合先 PR の merge 状態が `CLEAN` であることを確認する。`UNKNOWN` は merge 可能として扱わない。

## 0d. 仕様凍結ポリシー

実装着手可能な仕様として扱うには、対象仕様を一時的に凍結する。

仕様凍結は、以下をすべて満たす状態をいう。

| 判定項目 | 条件 |
|----------|------|
| 成熟度 | 対象項目が `仕様化済み・未実装` である。 |
| 詳細仕様 | `ADLAIRE_CI_DETAIL_SPEC.md` に入力、出力、状態、処理順序、異常系、検証条件が明記されている。 |
| テンプレート | `ADLAIRE_CI_DETAIL_SPEC.md` §0h の機能仕様テンプレートの必須項目を満たしている。 |
| 対応表 | `ADLAIRE_CI_DETAIL_SPEC.md` §0i.1〜§0i.4 の詳細節対応表に対象機能が記載され、参照先の詳細仕様節と受け入れ条件が一致している。 |
| 横断整合 | API、SDK、UI、状態ファイル、セットアップ、受け入れ条件が矛盾していない。 |
| 未確定事項 | 実装判断に必要な未確定事項が残っていない。 |
| 変更境界 | 実装 PR で変更してよい範囲と変更してはならない範囲が明確である。 |

仕様凍結後、実装中に仕様不足を発見した場合は、実装 PR 内で独自判断による補完を行わず、仕様改訂 PR または同一 PR 内の仕様改訂コミットで凍結状態を更新する。

仕様凍結は永久固定ではない。変更する場合は、凍結解除理由、変更対象、影響範囲、再検証条件を PR 本文に明記する。

## 0e. 初期実装スコープ確定ポリシー

Go 版初期実装では、実装対象を本ファイルで `仕様化済み・未実装` と判定し、かつ `ADLAIRE_CI_DETAIL_SPEC.md` の対応表と該当 owner component 別詳細仕様に具体的な実装詳細が存在する範囲に限定する。

初期実装の対象範囲は以下とする。

| 対象 | 実装対象 | 境界 |
|------|----------|------|
| `components/builder.go` | 対象 | `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §1〜§9 に記載された CLI、Markdown 変換、静的 Web サイト出力、テーマコンポーネント、検証条件。 |
| `components/runner.go` | 対象 | `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §10〜§20 に記載された CI ランナー、状態ファイル、ビルド起動、転送、通知、ログ保存。 |
| `components/api.go` | 対象 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §21〜§22、§25 および `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 に記載された管理 API、認証、状態ファイル、セットアップ。 |
| `admin/adlaire-ci-sdk.js` | 対象 | `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 に記載された API 呼び出し契約、戻り値、エラー処理。 |
| `admin/index.html` | 対象 | `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 に記載された標準管理ツール UI、操作、表示、秘密情報消去。 |
| `components/mcp.go` | 対象外 | 将来計画。詳細仕様、起動手順、認証、ツール定義を別途仕様化するまで実装不可。 |

初期実装 PR では、上表の対象外項目、将来計画、未仕様化項目、改訂予定項目を実装してはならない。

初期実装中に対象範囲へ追加したい機能を発見した場合は、先に本節の表、`ADLAIRE_CI_DETAIL_SPEC.md` の該当詳細仕様、`DOCUMENT_INDEX.md` の状態表現を更新し、仕様凍結を再実施する。

初期実装スコープの実装順序と PR 分割は、`ADLAIRE_CI_DETAIL_SPEC.md` §0g の初期実装 Phase 分割に従う。

初期実装スコープの完了判定は、対象コンポーネントごとに `ADLAIRE_CI_DETAIL_SPEC.md` §0f の仕様策定完了チェック、§0g の対象 Phase 完了条件、§0h の機能仕様テンプレート、§0i.1〜§0i.4 の詳細節対応表の受け入れ条件を満たしていることを条件とする。チェック、Phase 完了条件、テンプレート、対応表のいずれかを満たさない対象は、実装済みとして扱ってはならない。

## 0f. Phase 実装単位ポリシー

実装順序、実装計画、実装 PR、完了判定は Phase 単位で行わなければならない。

`P0`、`P1`、`P2〜P5` などの優先度ラベル、抽象段階、API 内部分類、fixture 分類を、実装単位、PR 単位、完了判定単位として使ってはならない。

Phase は `ADLAIRE_CI_DETAIL_SPEC.md` §0g に定義された対象、実装範囲、依存条件、完了条件、検証条件に従わなければならない。

Phase の途中で未仕様化、将来計画、改訂予定の機能を追加してはならない。追加する場合は、先に状態分類、詳細仕様、検証条件、受け入れ条件を更新し、仕様凍結を再実施しなければならない。

API の内部説明や fixture 名に既存の段階名が残る場合でも、それらは検証分類としてのみ扱い、実装順序、実装 PR、完了判定の正本にしてはならない。

## 1. デザイントークン準拠

> ⚠️ **準拠義務：** `--adlaire-*` トークンの値は [Adlaire Design System](https://github.com/fqwink/Adlaire-Design-System)（`Tokens/` ディレクトリ）で定義された値に準拠すること。

- スクリプト側での独自トークンの追加・変更は行わない
- ライトモード固定のため、ダークモード用トークンブロックは不要
- トークン値の変更は ADS 側のアップデートに追従する形でのみ実施する

## 2. バージョン管理

正式リリース前は、ヘッダーの `仕様バージョン: V.N` と `リリースバージョン: V.X.N` を暫定表記として扱う。暫定表記は、バージョン体系そのものを示す placeholder であり、実在する release tag、GitHub Release、実装済みバージョンを意味しない。

実数運用を開始する場合は、同一 PR で以下をすべて実施する。

1. ヘッダーの `V.N` と `V.X.N` を具体値へ置換する。
2. 具体値に対応する tag / GitHub Release の作成条件を満たしているか確認する。
3. `README.md`、`DOCUMENT_INDEX.md`、`ADLAIRE_CI_DETAIL_SPEC.md` にバージョン表記がある場合は整合させる。
4. 実数運用開始後は placeholder 表記へ戻さない。

### 仕様バージョン V.N

本仕様書自体のバージョンを管理する。

| 項目 | 内容 |
|------|------|
| 形式 | `V.N`（`V` は固定、`N` は正の整数） |
| 更新方針 | 仕様書の変更・追記のたびに `N` を 1 以上インクリメントする |
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
| マスター仕様の方針・ポリシー変更 | `V.N` を 1 以上更新 | 変更しない | 実装物が変わらない仕様策定のみの変更。 |
| 詳細仕様の実装契約変更 | `V.N` を 1 以上更新 | 変更しない | 実装着手条件や受け入れ条件が変わる。 |
| 仕様影響ありの実装コード変更 | `V.N` を 1 以上更新 | `N` を 1 以上更新 | 実装に合わせて仕様を変更する場合。仕様 PR と実装 PR の両方の完了条件を満たす。 |
| 仕様準拠のみの実装コード変更 | 変更しない | `N` を 1 以上更新 | 既存仕様に完全準拠する実装のみ。仕様文書の変更は不要だが、検証結果は実装 PR に記録する。 |
| 実装内部のみの修正 | 変更しない | `N` を 1 以上更新 | 仕様上の入出力、状態、API、UI、運用手順に影響しない内部修正。 |
| 安定版リリース | 未反映の仕様変更がある場合のみ `V.N` を更新 | `X` と `N` を 1 以上更新 | GitHub Release と tag を作成する。 |
| 将来計画の追記 | `V.N` を 1 以上更新 | 変更しない | 実装可能仕様として扱わない。 |

仕様策定 PR で実装コードを変更しない場合、リリースバージョンを更新してはならない。実装 PR と仕様 PR を同一 PR にまとめる場合は、仕様変更と実装変更の両方の完了条件を満たす。

## 3. カスタマイズ可能範囲

以下の項目はスクリプト内で変更可能な設定ポイントである。

| 設定項目 | 場所 | 変更方法 |
|---------|------|---------|
| 入出力パス | `adlaire-ci-build --src` / `--out`、または `DefaultBuildConfig` | `--src` は Markdown ファイルまたは Markdown ディレクトリ、`--out` は出力サイトディレクトリ。CLI 引数を優先し、既定値変更時は `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §2 と整合させる |
| テーマ | `adlaire-ci-build --theme`、または `DefaultBuildConfig.Theme` | 初期仕様では `adlaire-default` のみ許可。カスタムテーマ、外部テンプレート、テーマパッケージは将来計画とする |
| デザイントークン値 | `adlaire-default` の `style.css` が定義する `:root { }` ブロック | ADS 準拠の範囲内で変更し、`DESIGN.md` と整合させる |
| ドキュメントタイトル | `PageData.Title` / `SiteData.Title` | `PageData.Title` が空の場合は `SiteData.Title` を使用する。`SiteData.Title` は空文字禁止。変更時は `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §5 `PageData` / `SiteData` 契約に従う |
| ヘッダー表示名 | `<span id="doc-title">` に出力する表示タイトル | `PageData.Title` が空の場合は `SiteData.Title` を表示し、別名を持たせない |
| バージョンバッジ | 安定版リリース情報を表示する場合の `PageData` 拡張 | `V.X.N` 形式。追加する場合は先に `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §5 の `PageData` にフィールドを追加する |
| TOC 対象見出しレベル | `buildTOC(headings []Heading)` | 初期仕様では h1〜h3 固定。変更する場合は §4.4、§6、§7.3〜§7.5 を同時に改訂する |

## 4. 外部ライブラリ・フレームワーク方針

### 基本原則

Part 1 §4.1 のゼロ依存・フルインハウス原則を正とする。開発言語の**標準ライブラリのみ**を採用し、主要機能は本リポジトリ内の仕様と内製実装で完結させる。

外部依存を追加しなければ実装できない機能は、原則として仕様不足または設計不備として扱う。実装者は外部依存の追加で不足仕様を補完してはならない。

### 外部フレームワーク

**いかなる条件でも禁止する。** 例外なし。

### 外部ライブラリ

原則禁止とする。例外採用は、以下の条件をすべて満たす場合に限る。

- 内製化が技術的に困難であり、標準ライブラリだけでは安全性または正確性を担保できない
- 採用範囲が単一責務に限定され、コンポーネント全体の自律性を壊さない
- 採用理由、代替困難性、責務範囲、削除方針、検証条件を本仕様に明記している
- §4 許可外部ライブラリ一覧へ登録している

開発コスト短縮、実装の容易さ、流行、一般的なベストプラクティスだけを理由にした採用は認めない。許可リスト外のライブラリ使用は認めない。

### 内製ライブラリ・フレームワーク

内製化したライブラリ・フレームワークは、責務、入力、出力、検証条件が本仕様に記載されている場合に限り採用できる。内製であっても、仕様未記載の共通基盤を暗黙に追加してはならない。

横断的な内製共通処理を追加する場合は、Part 1 §4.2 に従い、Core ではなく同格の共通責務コンポーネントとして仕様化する。`core`、`common`、`base`、`foundation`、`utils` の名称を使う package、directory、binary、component を作成してはならない。

### 内製スクリプト一覧

| スクリプト | 状態 | 役割 |
|-----------|------|------|
| `components/builder.go` | 実装済み | Go 版ビルドスクリプト（Markdown → 静的 Web サイト変換）として Phase 1 の `gofmt` と `go test` 検証済み。 |
| `components/runner.go` | 実装済み | Go 版 CI ランナー（変更検出・ビルド起動・通知・転送）として Phase 2 完了判定パスの `gofmt` と `go test` 検証済み。 |
| `components/api.go` | 実装中・検証未完了 | Go 版管理 API サーバー（常駐 HTTP サーバー） |
| `admin/adlaire-ci-sdk.js` | 仕様化済み・未実装 | JavaScript SDK（管理ツール用 API クライアント） |
| `admin/index.html` | 仕様化済み・未実装 | 標準管理ツール UI |
| `components/mcp.go` | 将来計画 | Go 版 MCP サーバー（将来追加予定 → §13 将来計画 MCP サーバー実装） |

> 内製スクリプト・ライブラリは §4 方針に基づき積極的に採用する。新規スクリプトを追加する場合は本一覧へ登録する。
> 内製 Go コンポーネントは Go 標準ライブラリを基本とする。外部依存は許可リスト登録を必須とする。
> 仕様化済み・未実装、実装中・検証未完了、または将来計画のスクリプトは、実装ファイル、詳細仕様、検証結果、関連文書が実装済みとして整合するまで実装済みとして扱わない。

### 許可外部ライブラリ一覧

内製スクリプト・ライブラリは許可外部ライブラリ一覧に記載しない。内製の管理は内製スクリプト一覧で行う。

| ライブラリ | 用途 | 許可理由 |
|-----------|------|---------|
| （なし） | — | — |

> 許可外部ライブラリは存在しない。

---

### — CI ランナー —

## 5. CI ランナー セキュリティポリシー

- GitHub PAT（Personal Access Token）はファイル（`.github_token`）に保存し、パーミッションを `600` に設定する。スクリプト内にハードコードしない
- PAT のスコープは `contents: read`（読み取り専用）のみ付与する。書き込みスコープは不要
- ランナーは外部公開エンドポイントを持たない。サーバーから GitHub API への送信のみで動作する

## 6. CI ランナー 実行ポリシー

- blob SHA が前回実行時と同一の場合はビルドをスキップする（差分なしと判断）
- `components/runner.go` から生成する `adlaire-ci-runner` は 1 回実行して終了する oneshot 設計とし、多重起動は systemd タイマーの設定（`OnUnitActiveSec`）で防ぐ
- `pipeline.sh` の終了コードが `0` 以外の場合はビルド失敗としてログに記録する
- SHA ファイルはビルド成功後にのみ更新する。ビルド失敗時は前回 SHA を保持し、次回起動時に再試行する

## 7. CI ランナー ブランチポリシー

- `BRANCH_TARGETS` リストで 1 件以上のブランチターゲットを定義する。デフォルトは `main` ブランチの 1 エントリ構成
- 各ブランチへのマージ後、次回ポーリングサイクル（最大 5 分以内）で変更を検出しビルドを実行する
- 複数エントリを定義した場合はリスト順に順次処理する（並列処理は対象外）
- `BRANCH_TARGETS` の各エントリは `branch`・`target_file`・`sha_file`・`src`・`out`・`deploy_targets` を持つ（→ §12 設定値）

### — 管理ツール —

## 8. SDK ポリシー

本節は、仕様化済み・未実装の `admin/adlaire-ci-sdk.js` に適用する。

- SDK は内製とし、外部ライブラリに依存しない（→ Part 2 §4）
- SDK の対応言語追加は本ドキュメントへの記載を先行させる
- バックエンド API の変更は SDK の更新を伴う
- SDK は ES Module とし、`AdlaireCI` と `AdlaireCIError` を明示 export する
- SDK の実行環境、timeout、error class、`streamBuild()` の `StreamHandle` 契約は `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 を正とする
- SDK は自動 retry、戻り値補完、token 永続化、global 代入を行ってはならない

## 9. 標準管理ツール ポリシー

本節は、仕様化済み・未実装の `admin/index.html` に適用する。

- バニラ HTML / CSS / JavaScript のみで実装する。外部フレームワーク・外部ライブラリは使用しない（→ Part 2 §4）
- バックエンドとの通信はすべて SDK 経由とする。SDK を迂回した直接 API 呼び出しは行わない
- カスタマイズを妨げる密結合な実装は避ける
- DOM id、`data-panel`、form field name、初期ロード順、イベント処理順、成功/失敗表示、秘密情報消去条件は `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 を正とする
- UI は `localStorage`、`sessionStorage`、Cookie から token を復元してはならない
- `fetch()`、`XMLHttpRequest`、`EventSource`、`ReadableStream` reader を UI から直接生成してはならない

## 10. データ永続化ポリシー

- **初期方針：データベース不使用。** 状態はファイルで管理する（`.last_sha`・`.admin_credentials`・`.totp_secret`・`.build_history`・`.build_status.json`・`.notify_config`・`.notify_log`・`.notify_pending`・`.server_config`・`.access_log`・`.audit_log`・`.api_rate_state`・`.repo_config`・`.config_log`・`.access_control`・`.hooks`・`.maintenance`・`.alert_rules`・`.tag_rules`・`.pipeline_config`・`.notes`・`.smtp_config`・`.smtp_secret`・`.dashboard_layout`・`.pending_transfers`・`.build_lock`・`.build_state`・`.build_circuit_state`・`.branch_config`・`.webhook_events.json`・`.local_watch_state.json`・`.build_cache.json`・`.build_cache/`・`.dependency_manifest.json`・`.approval_queue`・`.build_trends.json`・`.build_chain_config`・`.api_tokens`・`.build_logs/`・`.snapshots/`・ビルドログ等）
- RDBMS・NoSQL・組み込み DB（SQLite 等）を問わず、初期仕様ではいかなるデータベースも採用しない
- 将来的にデータベースを採用する場合は、本ドキュメントへの仕様追記と §4 許可外部ライブラリ一覧の更新を先行させる
- **データ形式：フラットファイル JSON 形式**を標準とする
- ネストは最小限に抑え、1 ファイル 1 用途とする
- ファイルエンコーディングは UTF-8 とする

## 11. 管理ツール 認証ポリシー

本節は、仕様化済み・未実装の管理 API サーバーおよび標準管理ツールに適用する。

- **初期構成：シングルユーザー（`admin`）**
- 初期パスワードは `admin` とする
- 初回ログイン時はパスワード変更を促す通知を表示する
- **5 回目のログイン時はパスワード変更を強制する**（変更完了まで管理画面の操作を制限する）
- パスワードは平文保存禁止。ハッシュ化して保存する（§10 方針に基づきフラットファイル JSON 形式でファイル管理 → `ADLAIRE_CI_DETAIL_API_SPEC.md` §25）

## 12. 管理 API サーバー セキュリティポリシー

本節は、実装中・検証未完了の `components/api.go` に適用する。

- `HOST` は `127.0.0.1` に固定し、外部へ直接公開しない
- HTTPS は nginx 等のリバースプロキシでターミネートする。`components/api.go` 自体に TLS を実装しない
- セッショントークンはインメモリで管理し、ファイルに書き出さない
- API エンドポイントはすべて認証必須とする（`GET /api/health` を除く）

---

## 詳細仕様

Part 3 の実装詳細は、入口、索引、共通固定値、実装前確認項目、検証マトリクス、Phase、詳細節対応表、リポジトリ内ソース配置、責務 component 別詳細仕様ファイル管理仕様、横断補足契約については `ADLAIRE_CI_DETAIL_SPEC.md`、各 owner component の詳細本文については `ADLAIRE_CI_DETAIL_*_SPEC.md` を正とする。

`ADLAIRE_CI_DETAIL_SPEC.md` には、方針、ポリシー、実装状態、実装可否、ロードマップ状態、PR 分割判断を記載しない。

`ADLAIRE_CI_DETAIL_SPEC.md` と owner component 別の `ADLAIRE_CI_DETAIL_*_SPEC.md` は、Part 1 §4a と Part 2 §0 に定める範囲に従い、実装者が迷わない入出力、状態、処理順序、異常系、検証条件だけを維持する。
