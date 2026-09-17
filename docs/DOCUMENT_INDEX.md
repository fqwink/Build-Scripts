# Build-Scripts Document Index

このファイルは、Build-Scripts リポジトリ内の文書・実装ファイルの参照先と役割を整理する索引である。

仕様・詳細仕様・補助文書・本索引は `docs/` 配下に集約する。ルールブック `AGENTS.md` と入口文書 `README.md` はリポジトリ root に置く。

## Specification Structure

仕様文書の構造は、作業ルール、入口、索引、正本、詳細入口、責務 component 別本文、補助文書に分ける。

| 階層 | ファイル | 位置付け | 主な用途 |
|------|----------|----------|----------|
| 作業ルール | `AGENTS.md` | 最上位ルールブック | 承認、Git 運用、仕様書管理、文書整合ルールを判断する。 |
| 入口 | `README.md` | 初見向け入口 | リポジトリ概要、読む順番、主要ファイルを把握する。 |
| 索引 | `docs/DOCUMENT_INDEX.md` | 文書・実装ファイル索引 | 文書配置、正本関係、実装ファイル所在を確認する。 |
| マスター仕様 | `docs/SPEC.md` | 方針・ポリシー正本 | 実装状態、実装可否、ロードマップ、禁止事項を判断する。 |
| 詳細仕様入口 | `docs/DETAIL_INDEX.md` | Part 3 詳細仕様セット入口 | 詳細仕様の読み方、共通固定値、対応表、横断補足契約を確認する。 |
| 詳細仕様本文 | `docs/details/*.md` | owner component 別本文 | 対象 component の入出力、状態、処理順序、異常系、検証条件を実装単位で確認する。 |
| 補助文書 | `docs/DESIGN.md` | 生成静的 Web サイトのデザイン補助 | レイアウト、色、タイポグラフィ、TOC、コードブロック等の視覚仕様を確認する。 |

仕様判断では、上位階層が下位階層を置き換えるのではなく、各階層の正本範囲だけを参照する。`docs/DOCUMENT_INDEX.md` は所在と役割の索引であり、仕様本文を定義しない。

## Reading Order

| 順序 | ファイル | 目的 |
|------|----------|------|
| 1 | `AGENTS.md` | 作業ルール、承認、Git 運用、文書整合ルールを確認する。 |
| 2 | `README.md` | 初見向けの概要、読む順番、主要ファイルを確認する。 |
| 3 | `docs/DOCUMENT_INDEX.md` | 文書と実装ファイルの役割、正本関係、配置を確認する。 |
| 4 | `docs/SPEC.md` | 方針、ポリシー、実装状態、実装可否、ロードマップを確認する。 |
| 5 | `docs/DETAIL_INDEX.md` | 詳細仕様の入口、読み順、共通固定値、対応表、ソース配置を確認する。 |
| 6 | `docs/details/*.md` | 対象 owner component の入出力、状態、処理順序、異常系、検証条件を確認する。 |

上記の順序は、文書整理、仕様改訂、実装、検証、PR 作成のすべてで共通とする。`docs/DOCUMENT_INDEX.md` は索引であり、仕様判断の正本ではない。

## Documents

| ファイル | 役割 |
|---------|------|
| `README.md` | 初見向け入口。概要、読む順番、主要ファイルを示す。 |
| `docs/SPEC.md` | 方針、ポリシー、実装状態、実装可否、ロードマップを示す。 |
| `docs/DETAIL_INDEX.md` | 詳細仕様の入口、読み順、共通固定値、対応表、リポジトリ内ソース配置、横断補足契約を示す。 |
| `docs/details/builder.md` | `builder` owner component の詳細仕様。 |
| `docs/details/runner.md` | `runner` owner component の詳細仕様。 |
| `docs/details/api.md` | `api` owner component の詳細仕様。 |
| `docs/details/admin.md` | `admin` owner component の詳細仕様。 |
| `docs/details/sdk.md` | `sdk` owner component の詳細仕様。 |
| `docs/details/ui.md` | `ui` owner component の詳細仕様。 |
| `docs/details/setup.md` | `setup` owner component の詳細仕様。 |
| `docs/details/statefile.md` | `statefile` owner component の詳細仕様。 |
| `docs/details/archive.md` | `archive` owner component の詳細仕様。 |
| `docs/details/commitstatus.md` | `commitstatus` owner component の詳細仕様。 |
| `docs/details/security.md` | `security` owner component の詳細仕様。 |
| `docs/details/fixture.md` | `fixture` owner component の詳細仕様。 |
| `docs/DESIGN.md` | 生成静的 Web サイトのデザイン仕様。レイアウト、色、タイポグラフィ、TOC、コードブロック等の視覚仕様を整理する。 |
| `AGENTS.md` | エージェント作業ルールブック。承認、仕様書管理、実装管理、Git 運用、文書整合の最上位ルール。 |
| `docs/DOCUMENT_INDEX.md` | 本索引。文書・実装ファイルの役割と所在を示す。仕様本文を定義しない。 |

## Detail Spec Management

詳細仕様は、責務 component 別の `docs/details/*.md` を本文として管理する。`docs/DETAIL_INDEX.md` は、詳細仕様本文を集約せず、入口、読み順、共通固定値、対応表、リポジトリ内ソース配置、横断補足契約だけを持つ。

下表は、詳細仕様本文の配置先を示す索引である。各ファイルの責務境界、持つ内容、持たない内容、owner / collaborator の扱いは `docs/DETAIL_INDEX.md` §0b および §0b.1 を正とする。

| ファイル | 状態 | 役割 |
|----------|------|------|
| `docs/details/builder.md` | 分割済み | `builder` owner の Markdown 変換、静的 Web サイト出力、HTML / CSS / JavaScript、theme component、builder fixture、builder owner 追加機能。 |
| `docs/details/runner.md` | 分割済み | `runner` owner の GitHub 監視、設定読取、状態ファイル更新呼び出し、pipeline、deploy、snapshot 作成トリガー、通知、runner fixture、runner owner 追加機能。 |
| `docs/details/api.md` | 分割済み | `api` owner の HTTP 共通契約、endpoint、request / response、状態ファイル read/write 呼び出し境界、認証連携、api owner 追加機能。API fixture は `docs/details/fixture.md` §22-F。 |
| `docs/details/admin.md` | 分割済み | `admin` owner の管理 UI 静的ファイル配布物構成、配置、検証、HTTP 静的配信境界。 |
| `docs/details/sdk.md` | 分割済み | `sdk` owner の SDK class、method、HTTP 対応、query / body 生成、error、stream、token 破棄。 |
| `docs/details/ui.md` | 分割済み | `ui` owner の DOM id、panel、操作、表示状態、SDK 呼び出し、秘密情報消去。 |
| `docs/details/setup.md` | 分割済み | `setup` owner のバイナリ配布、配置、systemd、セットアップ、アップデート、リリース成果物検証。 |
| `docs/details/statefile.md` | 分割済み | `statefile` owner の状態ファイル共通仕様、lock、atomic write、JSON Lines、破損時処理、状態読取 adapter、主要 schema。 |
| `docs/details/archive.md` | 分割済み | `archive` owner の build log archive、snapshot、download、delete、rollback、cleanup。 |
| `docs/details/commitstatus.md` | 分割済み | `commitstatus` owner の GitHub Commit Status API payload、送信順、失敗時非反転、保存値、secret mask、検証条件。 |
| `docs/details/security.md` | 分割済み | `security` owner の API token scope、API key、audit、session timeout、TOTP、rate limit、漏えい禁止、security 横断順序。 |
| `docs/details/fixture.md` | 分割済み | fixture manifest、assertion、fake、testdata、expected / effects、Phase 3 / Phase 4 API fixture、api / sdk / ui / statefile cross fixture、受け入れ fixture 共通契約、PR 証跡テンプレート、acceptance checklist、差し戻し条件、実装 PR 完了証跡。 |

詳細仕様を改訂する場合は、`docs/SPEC.md` で実装状態と実装可否を確認し、`docs/DETAIL_INDEX.md` の対応表から owner component を特定し、該当する `docs/details/*.md` を本文として更新する。`docs/DOCUMENT_INDEX.md` は配置と役割の索引に限定し、仕様本文、詳細仕様本文、実装状態の最終判断を定義しない。

## Specified Components

`docs/SPEC.md` と `docs/DETAIL_INDEX.md` では、以下のコンポーネントも仕様化されている。

リポジトリ内ソース配置の標準構成は、`docs/SPEC.md` Part 1 §4.3 と `docs/DETAIL_INDEX.md` §0j を参照する。

下表は、現行リポジトリに存在する実装ファイルと、実装不可の将来計画ファイルを示す。標準ソース配置への移行は完了済みであり、`main.go`、現存する `components/*.go`、`admin/` 配下の静的 UI ファイル、`testdata/<component>/` を現行配置として扱う。

標準配置図に含まれる未実装 path は、該当 owner component が実装対象になった PR で追加する。標準配置図に含まれていることだけを理由に、未実装ファイル、将来計画ファイル、空ディレクトリ、placeholder を作成しない。

| パス | component | 状態 | 役割 |
|------|-----------|------|------|
| `main.go` | `-` | 実装済み | 起動入口。実行ファイル名に応じて `builder` または `runner` component を呼び出す。 |
| `components/builder.go` | `builder` | 実装済み | Go 版静的 Web サイトビルドスクリプト。`adlaire-ci-build` バイナリとして実行する。 |
| `components/builder_test.go` | `builder` | 実装済み | `components/builder.go` の Phase 1 fixture テスト。 |
| `go.mod` | `-` | 実装済み | Go module 定義。外部 module は追加しない。 |
| `testdata/builder/` | `builder` | 実装済み | Phase 1 の受け入れ fixture 入力。 |
| `components/runner.go` | `runner` | 実装済み | Go 版 CI ランナー。`adlaire-ci-runner` バイナリとして実行する。Phase 2 完了判定パスを対象とする。 |
| `components/runner_test.go` | `runner` | 実装済み | `components/runner.go` の Phase 2 fixture、hardening、完了判定パステスト。 |
| `components/api.go` | `api` | 実装済み | 管理 API サーバー。常駐 HTTP サーバーとして Adlaire CI の状態確認・操作 API を提供する。 |
| `admin/adlaire-ci-sdk.js` | `sdk` | 実装済み | 管理ツール用 JavaScript SDK。管理 API 通信を抽象化する。 |
| `admin/index.html` | `ui` | 実装済み | 標準管理ツール UI。SDK 経由で API と通信する。 |
| `components/mcp.go` | `mcp` | 将来計画 | MCP サーバー。現時点では実装可能な詳細仕様を持たず、MCP 専用詳細仕様が新設されるまで実装対象ではない。 |

## Source Of Truth

| 判断対象 | 正本 |
|----------|------|
| 作業ルール、承認、Git 運用、文書整合 | `AGENTS.md` |
| 方針、ポリシー、実装状態、実装可否、ロードマップ、正本関係 | `docs/SPEC.md` |
| 詳細仕様の入口、索引、共通固定値、実装前確認項目、検証マトリクス、Phase、詳細節対応表、リポジトリ内ソース配置、横断補足契約 | `docs/DETAIL_INDEX.md` |
| 各 component の入出力、状態、処理順序、異常系、検証条件の本文 | owner component 別の `docs/details/*.md` |
| 生成静的 Web サイトのデザイン補助 | `docs/DESIGN.md` |
| 文書・実装ファイルの参照先と役割 | `docs/DOCUMENT_INDEX.md` |
| 実装ファイル、テスト、fixture の所在 | `docs/DOCUMENT_INDEX.md` の Specified Components |

`docs/DOCUMENT_INDEX.md` は索引であり、仕様・デザイン・実装判断の正本ではない。仕様を変更する場合は、先に該当する正本仕様書を更新し、その内容に基づいて実装ファイルを更新する。

## Consistency Notes

現時点では、`docs/SPEC.md`、`docs/DETAIL_INDEX.md`、owner component 別の `docs/details/*.md` に記載された一部コンポーネントや機能は仕様化済みだが、リポジトリ内に実装ファイルが存在しない。

仕様化済みだが未実装の内容は、実装済み機能として扱わない。

標準ソース配置への移行は完了済みである。`components/builder.go` は Phase 1 実装済みであり、`components/runner.go` は Phase 2 完了判定パス実装済みである。Go toolchain による `gofmt` と `go test` の検証対象である。

`build_spec.go`、`runner.go`、`build_spec_test.go`、`runner_test.go`、`testdata/build_spec/` は現行実装実体として扱わない。標準配置への移行完了条件は `docs/DETAIL_INDEX.md` §0j を正とする。
