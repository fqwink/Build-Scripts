# Build-Scripts Document Index

このファイルは、Build-Scripts リポジトリ内の文書・実装ファイルの参照先と役割を整理する索引である。

## Documents

| ファイル | 役割 |
|---------|------|
| `ADLAIRE_CI_SPEC.md` | Adlaire CI の方針、ポリシー、実装状態、正本関係を定めるマスター仕様書正本。 |
| `ADLAIRE_CI_DETAIL_SPEC.md` | `ADLAIRE_CI_SPEC.md` の Part 3 詳細仕様の入口。索引、共通固定値、実装前確認項目、検証マトリクス、Phase、詳細節対応表、リポジトリ内ソース配置、責務 component 別詳細仕様ファイル管理仕様、§27.38a の横断補足契約を持つ。状態ファイル詳細は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md`、fixture 共通契約は `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` を正とする。方針、ポリシー、実装状態、ロードマップ状態、実装可否は記載しない。 |
| `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` | `builder` owner component の詳細仕様。Markdown 変換、静的 Web サイト出力、HTML / CSS / JavaScript、theme component、builder fixture、builder owner 追加機能を扱う。 |
| `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` | `runner` owner component の詳細仕様。GitHub 監視、状態ファイル更新、pipeline、deploy、snapshot、通知、runner fixture、runner owner 追加機能を扱う。 |
| `ADLAIRE_CI_DETAIL_API_SPEC.md` | `api` owner component の詳細仕様。HTTP 共通契約、endpoint、状態ファイル read/write の呼び出し境界、認証連携、API owner 追加機能を扱う。状態ファイル schema と更新手順は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md`、API fixture は `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §22-F を正とする。 |
| `ADLAIRE_CI_DETAIL_ADMIN_SPEC.md` | `admin` owner component の詳細仕様。管理 UI 静的ファイルの配布物構成、配置、検証、HTTP 静的配信境界を扱う。 |
| `ADLAIRE_CI_DETAIL_SDK_SPEC.md` | `sdk` owner component の詳細仕様。SDK class、method、HTTP 対応、error、stream、token 破棄を扱う。 |
| `ADLAIRE_CI_DETAIL_UI_SPEC.md` | `ui` owner component の詳細仕様。DOM id、panel、操作、表示状態、SDK 呼び出し、秘密情報消去を扱う。 |
| `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` | `setup` owner component の詳細仕様。バイナリ配布、配置、systemd、セットアップ、アップデート、リリース成果物検証を扱う。 |
| `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` | `statefile` owner component の詳細仕様。状態ファイル共通仕様、lock、atomic write、JSON Lines、破損時処理、状態読取 adapter、主要 schema を扱う。 |
| `ADLAIRE_CI_DETAIL_ARCHIVE_SPEC.md` | `archive` owner component の詳細仕様。build log archive、snapshot、download、delete、rollback、cleanup を扱う。 |
| `ADLAIRE_CI_DETAIL_COMMITSTATUS_SPEC.md` | `commitstatus` owner component の詳細仕様。GitHub Commit Status API payload、送信順、失敗時非反転、保存値、secret mask を扱う。 |
| `ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` | `security` owner component の詳細仕様。API token scope、API key、audit、session timeout、TOTP、rate limit、漏えい禁止、security 横断順序を扱う。 |
| `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` | fixture / fake / testdata / assertion / PR 証跡の詳細仕様。Phase 別 fixture 配置、API P0〜P5 fixture、api / sdk / ui / statefile cross fixture、§27 fixture カタログ、manifest、expected/effects、受け入れゲート、実装 PR 完了証跡を扱う。 |
| `DESIGN.md` | 生成静的 Web サイトのデザイン仕様。レイアウト、色、タイポグラフィ、TOC、コードブロック等の視覚仕様を整理する。 |
| `AGENTS.md` | エージェント作業ルールブック。承認、仕様書管理、実装管理、Git 運用、文書整合の最上位ルール。 |
| `DOCUMENT_INDEX.md` | 本索引。リポジトリ内の文書・実装ファイルの役割と正本関係を示す。仕様正本ではない。 |

## Detail Spec Management

詳細仕様は、責務 component 別の分割済み詳細仕様ファイルとして管理する。

`ADLAIRE_CI_DETAIL_SPEC.md` は、詳細仕様の入口、索引、共通固定値、実装前確認項目、検証マトリクス、Phase、詳細節対応表、リポジトリ内ソース配置、責務 component 別詳細仕様ファイル管理仕様、横断補足契約を持つ。各 component の入出力、処理順序、状態、異常系、検証条件の本文は下表の owner component 別詳細仕様ファイルを正とする。

fixture、fake、testdata、expected / effects、PR 証跡、acceptance checklist、差し戻し条件、実装 PR 完了証跡は `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` を正とする。`ADLAIRE_CI_DETAIL_SPEC.md` §0e、§0g、§0i は完了判定の入口であり、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 は setup / release / Phase 判定の実行条件である。両ファイルは fixture 名、fake 動作、PR 証跡項目、差し戻し条件を重複定義しない。

| ファイル | 状態 | 役割 |
|----------|------|------|
| `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` | 分割済み | `builder` owner の Markdown 変換、静的 Web サイト出力、HTML / CSS / JavaScript、theme component、builder fixture。 |
| `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` | 分割済み | `runner` owner の GitHub 監視、状態ファイル更新、pipeline、deploy、snapshot、通知、runner fixture。 |
| `ADLAIRE_CI_DETAIL_API_SPEC.md` | 分割済み | `api` owner の HTTP 共通契約、endpoint、状態ファイル read/write、認証連携。API fixture は `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §22-F。 |
| `ADLAIRE_CI_DETAIL_ADMIN_SPEC.md` | 分割済み | `admin` owner の管理 UI 静的ファイル配布物構成、配置、検証、HTTP 静的配信境界。 |
| `ADLAIRE_CI_DETAIL_SDK_SPEC.md` | 分割済み | `sdk` owner の SDK class、method、HTTP 対応、error、stream、token 破棄。 |
| `ADLAIRE_CI_DETAIL_UI_SPEC.md` | 分割済み | `ui` owner の DOM id、panel、操作、表示状態、SDK 呼び出し、秘密情報消去。 |
| `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` | 分割済み | `setup` owner のバイナリ配布、配置、systemd、セットアップ、アップデート、リリース成果物検証。 |
| `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` | 分割済み | `statefile` owner の状態ファイル共通仕様、lock、atomic write、JSON Lines、破損時処理、状態読取 adapter、主要 schema。 |
| `ADLAIRE_CI_DETAIL_ARCHIVE_SPEC.md` | 分割済み | `archive` owner の build log archive、snapshot、download、delete、rollback、cleanup。 |
| `ADLAIRE_CI_DETAIL_COMMITSTATUS_SPEC.md` | 分割済み | `commitstatus` owner の GitHub Commit Status API payload、送信順、失敗時非反転、保存値、secret mask。 |
| `ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` | 分割済み | `security` owner の API token scope、API key、audit、session timeout、TOTP、rate limit、漏えい禁止、security 横断順序。 |
| `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` | 分割済み | fixture manifest、assertion、fake、testdata、API P0〜P5 fixture、api / sdk / ui / statefile cross fixture、受け入れ fixture 共通契約、PR 証跡テンプレート、実装 PR 完了証跡。 |

責務 component 別詳細仕様ファイルは、各ファイルの `## 0. 責務境界` を実装前に確認する。`owner component` は主本文を持つ component、`collaborator component` は呼び出し境界、schema、fixture、security、setup、表示、受け入れ条件を参照する component として扱う。collaborator 側の参照は、owner component の主本文を上書きしない。

各詳細仕様ファイルの `持つ内容` と `持たない内容` が本文、`ADLAIRE_CI_DETAIL_SPEC.md` §0b、または本索引と矛盾する場合、その項目は実装判断に使わず、先に文書整合を行う。

`ADLAIRE_CI_DETAIL_SPEC.md` §0e、§0g、§0i は、同じ owner component / collaborator component 境界で読む。Phase の主対象は owner component とする。statefile、security、archive、commitstatus、admin、fixture、setup は、owner component の該当機能が状態ファイル、認証・監査、snapshot / log archive、commit status、admin 配布物、fixture、setup / update 手順を参照または変更する場合に collaborator component として検証対象、schema 参照、setup 参照、security 参照、fixture 参照、配布境界確認を提供する。collaborator component は owner component の入出力、状態、endpoint、SDK method、UI DOM、fixture を追加定義しない。

実装者が詳細仕様を読む順序は、`ADLAIRE_CI_SPEC.md` で実装状態と実装可否を確認し、`ADLAIRE_CI_DETAIL_SPEC.md` §0〜§0j で共通固定値、責務 component、詳細節対応表、リポジトリ内ソース配置を確認し、owner component の分割先詳細仕様ファイルを主本文として読む順に固定する。collaborator component の分割先詳細仕様ファイルは、呼び出し境界、schema、表示、security、setup、fixture、検証観点として参照し、owner component の主本文を上書きしない。

`COMMON`、`CORE`、`BASE`、`SHARED`、`FOUNDATION`、その他の横断共通基盤ファイルは作成しない。横断する固定値、読み順、対応表、横断補足契約は `ADLAIRE_CI_DETAIL_SPEC.md` の入口・索引・共通固定値・管理仕様として扱い、component として扱わない。

`ADLAIRE_CI_DETAIL_SPEC.md` §27.38a は、runner / builder / api / sdk / ui / statefile / archive にまたがる横断補足契約であり、責務 component 別詳細仕様ファイルへ分割しない。§27.21〜§27.38 または api / sdk / ui / statefile の横断連動を実装する場合は、owner component の分割先詳細仕様ファイルと `ADLAIRE_CI_DETAIL_SPEC.md` §27.38a を同時に確認する。§27.38a は個別機能の入出力、状態 schema、endpoint、SDK method、UI DOM、fixture schema を定義しない。これらは owner component 別詳細仕様ファイルを正とし、§27.38a は横断処理順、同期禁止、成功後再取得、失敗時固定、横断受け入れ観点だけを補足する。

## Specified Components

`ADLAIRE_CI_SPEC.md` と `ADLAIRE_CI_DETAIL_SPEC.md` では、以下のコンポーネントも仕様化されている。

リポジトリ内ソース配置の標準構成は、`ADLAIRE_CI_SPEC.md` Part 1 §4.3 と `ADLAIRE_CI_DETAIL_SPEC.md` §0j を参照する。

下表は、現行リポジトリに存在する実装ファイルと、標準配置で仕様化済みの未実装ファイルを示す。`components/` 標準配置への移行前は、現行ファイルを実装実体として扱う。`components/builder.go` と `components/runner.go` は標準配置名であり、現行実装実体はそれぞれ `build_spec.go` と `runner.go` である。

標準配置図に含まれる未実装 path は、該当 owner component が実装対象になった PR で追加する。標準配置図に含まれていることだけを理由に、未実装ファイル、将来計画ファイル、空ディレクトリ、placeholder を作成しない。

| パス | component | 状態 | 役割 |
|------|-----------|------|------|
| `build_spec.go` | `builder` | 実装済み | Go 版静的 Web サイトビルドスクリプト。`adlaire-ci-build` バイナリとして実行する。標準配置名 `components/builder.go` の現行実装実体。 |
| `build_spec_test.go` | `builder` | 実装済み | `build_spec.go` の Phase 1 fixture テスト。標準移行後のテスト配置は実装 PR で決定する。 |
| `go.mod` | `-` | 実装済み | Go module 定義。外部 module は追加しない。 |
| `testdata/build_spec/` | `builder` | 実装済み | Phase 1 の受け入れ fixture 入力。標準移行後の配置は `testdata/builder/`。 |
| `runner.go` | `runner` | 実装済み | Go 版 CI ランナー。`adlaire-ci-runner` バイナリとして実行する。Phase 2 完了判定パスを対象とする。標準配置名 `components/runner.go` の現行実装実体。 |
| `runner_test.go` | `runner` | 実装済み | `runner.go` の Phase 2 fixture、hardening、完了判定パステスト。標準移行後のテスト配置は実装 PR で決定する。 |
| `components/api.go` | `api` | 仕様化済み・未実装 | 管理 API サーバー。常駐 HTTP サーバーとして Adlaire CI の状態確認・操作 API を提供する。 |
| `components/admin.go` | `admin` | 仕様化済み・未実装 | 管理 UI 静的ファイルの配布物構成、配置、検証、HTTP 静的配信境界を提供する。 |
| `admin/adlaire-ci-sdk.js` | `sdk` | 仕様化済み・未実装 | 管理ツール用 JavaScript SDK。管理 API 通信を抽象化する。 |
| `admin/index.html` | `ui` | 仕様化済み・未実装 | 標準管理ツール UI。SDK 経由で API と通信する。 |
| `components/mcp.go` | `mcp` | 将来計画 | MCP サーバー。現時点では実装可能な詳細仕様を持たず、MCP 専用詳細仕様が新設されるまで実装対象ではない。 |

## Source Of Truth

方針、ポリシー、実装状態、正本関係の判断では `ADLAIRE_CI_SPEC.md` を正とする。

実装の具体的詳細の判断では、入口、索引、共通固定値、実装前確認項目、検証マトリクス、Phase、詳細節対応表、リポジトリ内ソース配置、責務 component 別詳細仕様ファイル管理仕様、横断補足契約は `ADLAIRE_CI_DETAIL_SPEC.md`、各 component の詳細本文は owner component 別の `ADLAIRE_CI_DETAIL_*_SPEC.md` を正とする。

デザイン判断では、`ADLAIRE_CI_SPEC.md` と矛盾しない範囲で `DESIGN.md` を参照する。

Go 版実装ファイルの挙動が `ADLAIRE_CI_SPEC.md`、`ADLAIRE_CI_DETAIL_SPEC.md`、または該当する owner component 別の `ADLAIRE_CI_DETAIL_*_SPEC.md` と矛盾する場合は、仕様と実装の不整合として扱う。

仕様を変更する場合は、先に該当する仕様書を更新し、その内容に基づいて実装ファイルを更新する。

`DESIGN.md` はデザイン仕様の補助文書であり、Adlaire CI 全体の機能仕様・運用仕様の正本ではない。

`DOCUMENT_INDEX.md` は索引であり、仕様・デザイン・実装判断の正本ではない。

`AGENTS.md` と他ファイルが作業ルール上矛盾する場合は、`AGENTS.md` を正とする。

## Consistency Notes

現時点では、`ADLAIRE_CI_SPEC.md`、`ADLAIRE_CI_DETAIL_SPEC.md`、owner component 別の `ADLAIRE_CI_DETAIL_*_SPEC.md` に記載された一部コンポーネントや機能は仕様化済みだが、リポジトリ内に実装ファイルが存在しない。

仕様化済みだが未実装の内容は、実装済み機能として扱わない。

現行配置では、`build_spec.go` は標準配置名 `components/builder.go` の現行実装実体として Phase 1 実装済みであり、`runner.go` は標準配置名 `components/runner.go` の現行実装実体として Phase 2 完了判定パス実装済みである。Go toolchain による `gofmt` と `go test` の検証を完了している。

Go 版コンポーネントの実装状態は、`ADLAIRE_CI_SPEC.md` を正とする。`ADLAIRE_CI_DETAIL_SPEC.md` は入口、索引、共通固定値、実装前確認項目、検証マトリクス、Phase、詳細節対応表、リポジトリ内ソース配置、責務 component 別詳細仕様ファイル管理仕様、横断補足契約として扱い、実装詳細本文と検証条件は owner component 別の `ADLAIRE_CI_DETAIL_*_SPEC.md` を参照する。

標準配置への移行完了後は、`build_spec.go`、`runner.go`、`build_spec_test.go`、`runner_test.go`、`testdata/build_spec/` を現行実装実体として扱わない。標準配置への移行完了条件は `ADLAIRE_CI_DETAIL_SPEC.md` §0j を正とする。
