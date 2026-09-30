# Adlaire CI — Roadmap

[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務は、実装 artifact と各機能の現在状態・未完了理由、Phase の順序・対象 owner・依存関係、機能インベントリ、将来計画を管理する唯一の正本である。状態語彙、状態定義、実装可否、状態遷移条件、`実装済み` への遷移条件は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity)、実装着手条件は [`docs/SPEC.md` 方針責務 §4.7](SPEC.md#sec-4-7)、[`docs/SPEC.md` ポリシー責務 §0d](SPEC.md#policy-spec-freeze)、[`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit)、仕様策定・改訂の完了条件は [`docs/SPEC.md` ポリシー責務 §0b](SPEC.md#policy-spec-pr-completion)、Phase 完了単位は [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit) を参照し、本文書では再定義しない。

<a id="状態計画責務-共通参照先"></a>
**状態・計画責務 共通参照先：**

本文書の状態セルは現在の割り当てだけを示す。endpoint、schema、SDK method、DOM、処理順序、fixture 本文は記載しない。

## 3. 実装 artifact 現在状態

この表は [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3) のディレクトリ構成で定義する実装 artifact だけを対象とする。owner component が持つ個別機能の現在状態は [統合機能インベントリ](#522-統合ロードマップ表) を参照する。実装不一致の技術的な確認内容は [`docs/details/fixture.md` fixture 証跡責務 現行実装整合証跡](details/fixture.md#current-implementation-alignment-evidence) を参照する。

| 実装 artifact | 現在状態 | 未完了理由 / 証跡 |
|-----------------|----------|-------------------|
| [`main.go`](../main.go) | 実装済み | 標準実行バイナリ 7 件の dispatch と version 受け渡しは実装・検証済みである。最終 Phase の実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 10 mcp 実装検証証跡](details/fixture.md#phase-10-mcp-implementation-evidence) を参照する。 |
| [`components/builder.go`](../components/builder.go) | 実装済み | Phase 1 builder 実装は [PR #74](https://github.com/fqwink/Build-Scripts/pull/74) で merge 済みである。実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 1 builder 実装検証証跡](details/fixture.md#phase-1-builder-implementation-evidence) を参照する。 |
| [`components/runner.go`](../components/runner.go) | 実装済み | Phase 2 runner 実装は [PR #75](https://github.com/fqwink/Build-Scripts/pull/75) で実装・検証済みである。実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 2 runner 実装検証証跡](details/fixture.md#phase-2-runner-implementation-evidence) を参照する。 |
| [`components/api.go`](../components/api.go) | 実装済み | Phase 4 `api` operations は実装・検証済みである。実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 4 api operations 実装検証証跡](details/fixture.md#phase-4-api-operations-implementation-evidence) を参照する。追加管理 API、SDK、UI、admin、mcp、正式 fixture directory harness の現在状態は各 owner の行と [統合機能インベントリ](#522-統合ロードマップ表) を参照する。 |
| [`components/admin.go`](../components/admin.go) | 実装済み | Phase 7 `admin` CLI 管理クライアントは実装・検証済みである。実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 7 admin CLI 実装検証証跡](details/fixture.md#phase-7-admin-cli-implementation-evidence) を参照する。 |
| [`components/setup.go`](../components/setup.go) | 実装済み | Phase 8 `setup` は実装・検証済みである。実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 8 setup 実装検証証跡](details/fixture.md#phase-8-setup-implementation-evidence) を参照する。 |
| [`components/release.go`](../components/release.go) | 実装済み | Phase 9 `release` は実装・検証済みである。実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 9 release 実装検証証跡](details/fixture.md#phase-9-release-implementation-evidence) を参照する。 |
| [`admin/adlaire-ci-sdk.js`](../admin/adlaire-ci-sdk.js) | 実装済み | Phase 5 SDK 実装は実装・検証済みである。実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 5 sdk 実装検証証跡](details/fixture.md#phase-5-sdk-implementation-evidence) を参照する。 |
| [`admin/index.html`](../admin/index.html) | 実装済み | Phase 6 UI 実装は実装・検証済みである。実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 6 ui 実装検証証跡](details/fixture.md#phase-6-ui-implementation-evidence) を参照する。正式 fixture directory harness の現在状態は [`ALIGN-07`](details/fixture.md#align-07)、[`ALIGN-24`](details/fixture.md#align-24) を参照する。 |
| [`components/mcp.go`](../components/mcp.go) | 実装済み | Phase 10 `mcp` は実装・検証済みである。実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 10 mcp 実装検証証跡](details/fixture.md#phase-10-mcp-implementation-evidence) を参照する。 |

<a id="roadmap-phase-plan"></a>

## 4. Phase 実装計画

Phase の実装単位、禁止事項、着手条件、完了判定方針は [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit) を参照する。[`docs/ROADMAP.md` 状態・計画責務 §4](ROADMAP.md#roadmap-phase-plan) は Phase の割り当て、現在状態、順序、依存関係だけを管理する。

<a id="roadmap-initial-phase-plan"></a>

## 4.1 初期実装 Phase 単位

| Phase | owner / scope | 現在状態 | 依存する Phase |
|-------|---------------|----------|----------------|
| Phase 1 | `builder` | 実装済み | なし |
| Phase 2 | `runner` | 実装済み | Phase 1 |
| Phase 3 | `api` request lifecycle | 実装済み | Phase 2 |
| Phase 4 | `api` operations | 実装済み | Phase 3 |
| Phase 5 | `sdk` | 実装済み | Phase 4 |
| Phase 6 | `ui` | 実装済み | Phase 5 |
| Phase 7 | `admin` CLI 管理クライアント | 実装済み | Phase 6 |
| Phase 8 | `setup` | 実装済み | Phase 7 |
| Phase 9 | `release` | 実装済み | Phase 8 |
| Phase 10 | `mcp` | 実装済み | Phase 9 |
| Phase 11 | バグ修正ゼロ化 / 仕様全般完了 gate / 全標準実装 artifact source-code audit / main dispatch / builder parser・config・output / runner source・pipeline・deploy・notify・corrupt recovery / API auth・config・backup・webhook・read model・filesystem・listener lifecycle・HTTP boundary / runtime clock・timer・entropy・parallel worker 決定性 / SDK・UI runtime / setup・release distribution・filesystem / admin・MCP bridge・protocol・lifecycle / admin CLI・API・SDK・UI route parity / output manifest / secret・queue・cross-owner regression / 正式 fixture directory harness / fixture manifest・catalog closure / test・contract drift zero | 仕様化済み・未実装 | Phase 10 |

現在の active Phase は Phase 11 である。初期実装 Phase 1 から Phase 10 まではすべて `実装済み` である。Phase 11 は [`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8) のバグ修正ゼロ化そのものを実装・検証する Phase として、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) を入口に扱う。Phase 11 に専用詳細仕様ファイルは置かない。source-code audit により Phase 11 対象へ割り当てた各項目は、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) から対象 owner 詳細本文、[`docs/details/fixture.md` fixture 証跡責務 Phase 11 fixture harness 参照](details/fixture.md#phase-11-fixture-harness-reference)、または [`docs/ROADMAP.md` 状態・計画責務 §5](ROADMAP.md#522-統合ロードマップ表) の対象外理由へ到達させる。仕様全般完了 gate は Phase 11 の仕様策定上の閉じ条件であり、`実装済み` への状態変更ではない。active Phase の決定条件と後続 Phase の禁止事項は [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit) を参照する。

各行の owner 詳細本文は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0b](DETAIL_INDEX.md#0b-詳細仕様参照表)、機能別の詳細節は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表)、fixture 証跡は [`docs/details/fixture.md`](details/fixture.md) を参照する。

後続 Phase の現在状態が `仕様化済み・未実装` または `実装中・検証未完了` であっても、active Phase でない Phase の新規実装着手可否は [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit) を参照する。

---

## 5. 統合機能インベントリ

状態の意味と遷移条件は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity)、実装詳細の入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) を参照する。以下の表は実装可能な機能、実装中の機能、未実装機能、将来計画を同じ現在状態語彙で管理する。

<a id="522-統合ロードマップ表"></a>
**統合機能インベントリ：** 以下を全機能の現在状態に関する唯一の一覧とする。

| 現在状態 | 担当領域 | 機能 | 詳細入口 / 次の扱い |
|----------|----------|------|----------------------|
| 実装済み | 管理ツール・SDK | JavaScript SDK 公開契約 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・UI | 標準管理ツール UI 契約 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | 状態管理 | 状態ファイル共通永続化契約 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | 検証基盤 | バグ修正ゼロ化 / 仕様全般完了 gate / cross-owner regression gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 仕様化済み・未実装 | 検証基盤 | 全標準実装 artifact source-code audit / artifact coverage zero gap | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 仕様化済み・未実装 | 検証基盤 | main dispatch / binary version / output manifest / file tree consistency gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 仕様化済み・未実装 | 検証基盤 | source-code audit inventory / direct I/O elimination / statefile common persistence gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 仕様化済み・未実装 | 検証基盤 | builder parser / config / output / generated site regression gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 仕様化済み・未実装 | 検証基盤 | runner source / pipeline / deploy / notification / corrupt state regression gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 仕様化済み・未実装 | 検証基盤 | API auth / config / backup / webhook / read model / filesystem regression gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 仕様化済み・未実装 | 検証基盤 | SDK / UI / admin CLI / MCP bridge / route parity regression gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 仕様化済み・未実装 | 検証基盤 | setup / release distribution / filesystem / external boundary regression gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 仕様化済み・未実装 | 検証基盤 | 正式 fixture directory harness / fixture identity / fixture manifest / contract drift zero gate | [`docs/details/fixture.md` fixture 証跡責務 Phase 11 fixture harness 参照](details/fixture.md#phase-11-fixture-harness-reference) |
| 仕様化済み・未実装 | 管理ツール・配布 | 管理 UI 静的配布物構成・Archive 検証 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | 管理ツール・配布 | 管理 UI 静的 HTTP 配信 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルドタイムアウト | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | ポーリング間隔の動的変更 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルドログのファイル保存 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | GitHub Webhook 受信 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ネットワーク断時の再試行 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | GitHub API レート制限自動待機 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | 転送後リモート整合性検証 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | マルチブランチビルド | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルドログ世代管理 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルド出力の外部転送 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルドクールダウン | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルド前の事前チェック | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | 定期強制ビルド | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルド中重複スキップ | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | GitHub PAT 有効期限の事前警告 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | コミット情報のビルドログ記録 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | GitHub API 連続失敗によるサーキットブレーカー | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | 出力サイトサイズ警告閾値 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | Webhook イベントログ | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルド所要時間の記録と統計入力 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルドアーティファクト世代管理 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | ビルドアーティファクト管理 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | ヘルスチェックエンドポイント | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | Webhook イベント一覧取得 API | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | ビルドログ重大度フィルター | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 変換レポート出力 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | builder 起動入口 / version | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | builder CLI parse / validation order | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | builder 入力 path / symlink / base-dir | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | サイドバー開閉 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | TOC 検索フィルター | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | トップへ戻るボタン | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | シンタックスハイライト | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 本文内全文検索 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | アンカーリンク自動検証 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | コードブロックの折りたたみ | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 印刷スタイル（`@media print`） | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 静的 Web サイト出力 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | テーマコンポーネント | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 外部リンクの自動処理 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 読み取り進捗バー | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | コードブロックのコピーボタン | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 見出しアンカーリンクコピー | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | TOC 開閉状態の永続化 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 見出しスラグ重複解決 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 前後章ナビゲーションボタン | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 内部リンク整合性チェック | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 見出し階層スキップ警告 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 読了時間推計と表示 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | Webhook 通知失敗リトライキュー | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | ブランチ設定の動的変更 API | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | 週次ビルドサマリー Webhook | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | 設定変更の詳細 diff 記録 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | テーブルのソート機能 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | キーボードショートカット | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | 複数ファイル監視 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | GitHub Commit Status API | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | 標準 builder command 拡張設定 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ドライラン実行モード | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルドログのアーカイブ圧縮 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ローカルファイル監視モード | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | タグ付きコミットのみビルド | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | ビルドキャッシュ | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルド通知連携 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルド時間トレンド記録 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルド失敗時の自動リトライ | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | 並列マルチターゲットビルド | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルド前後フック | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | 依存ファイルトラッキング | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | リモートビルド対応 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルドステータスファイル出力 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルド承認フロー | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ブランチ別環境変数 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルド依存チェーン | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルド優先度キュー | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | 失敗原因の自動分類 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルド実行環境の記録 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルドトリガー種別の記録 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | ビルド所要時間の異常検知 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | CI ランナー | 設定ファイル起動時整合性チェック | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | 管理ツール・API | マルチユーザー対応 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 仕様化済み・未実装 | 管理ツール・API | データストア切り替え | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 仕様化済み・未実装 | 管理ツール・API | 外部認証連携 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | TOTP 二要素認証 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | 管理ツール・API | 統計データの JSON エクスポート | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 仕様化済み・未実装 | 管理ツール・API | キュー内個別エントリのキャンセル | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | 設定バリデーション API | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | 管理ツール・API | Prometheus メトリクスエンドポイント | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | CLI 管理クライアント | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 仕様化済み・未実装 | 管理ツール・API | 設定の自動スナップショット | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 仕様化済み・未実装 | 管理ツール・API | ステータスバッジ生成 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 仕様化済み・未実装 | 管理ツール・API | ビルド履歴の自動削除設定 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | セッションタイムアウト変更設定 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | ビルドトリガー専用 API スコープ | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | 監査ログ | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | API レート制限 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | 管理ツール・API | ロールベースアクセス制御 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 仕様化済み・未実装 | 管理ツール・API | 設定スナップショット差分表示 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 仕様化済み・未実装 | 管理ツール・API | 複数プロジェクト管理 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | API キー管理 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | ビルドログの保存済み有限 SSE 配信 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | ビルド統計ダッシュボード | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | 管理ツール・API | ユーザー管理 API | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | IP アドレス制限 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | 管理ツール・API | API バージョニング | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 仕様化済み・未実装 | 管理ツール・API | API ドキュメント自動生成 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | 通知チャンネル管理 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | 管理ツール・API | ビルドキューの手動並び替え | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 仕様化済み・未実装 | 管理ツール・API | 設定テンプレート | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | API アクセスログ | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | 管理ツール・API | 管理者向けイベントフィード | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | ビルドキュー可視化 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | メンテナンスモード | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | 管理ツール・API | 読み取り専用共有リンク | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | アラート閾値設定 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | バックアップ／リストア | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | 管理ツール・API | API レスポンスキャッシュ制御 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 仕様化済み・未実装 | 管理ツール・API | スナップショット間サイト差分 API | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 仕様化済み・未実装 | 管理ツール・API | Webhook 送信履歴の手動再送 API | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | ビルドスクリプト | 差分ビルド | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 複数出力形式 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | Markdown 拡張記法サポート | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | コードブロック行番号表示 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 見出しの自動採番 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | セクション折りたたみ | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | TOC 深さ制御 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 最終更新日の自動埋め込み | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | diff ハイライト | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 画像の遅延読み込み | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | カスタムメタタグ注入 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | ライトモード固定 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | コードブロックのファイル名表示 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | テンプレート変数展開 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | HTML ミニファイ | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | TOC ハイライト追従（アクティブ見出し追跡） | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | Mermaid ダイアグラム描画 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 脚注サポート | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | インライン数式レンダリング | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | ページ内ナビゲーション履歴 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 読み上げ対応（アクセシビリティ） | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 画像ライトボックス | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 出力サイトへのビルドメタ埋め込み | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 印刷時 QR コード挿入 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | 定義リストサポート | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | ビルドスクリプト | タスクリストサポート | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 配布・セットアップ | 初回セットアップ | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.5](DETAIL_INDEX.md#0i5-setup--release) |
| 実装済み | 配布・セットアップ | 管理 API 導入 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.5](DETAIL_INDEX.md#0i5-setup--release) |
| 実装済み | 配布・セットアップ | バイナリアップデート | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.5](DETAIL_INDEX.md#0i5-setup--release) |
| 実装済み | 配布・セットアップ | アップデート rollback | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.5](DETAIL_INDEX.md#0i5-setup--release) |
| 実装済み | リリース | GitHub Release 成果物生成・公開前検証・公開 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.5](DETAIL_INDEX.md#0i5-setup--release) |
| 実装済み | MCP サーバー | MCP サーバー実装 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.7](DETAIL_INDEX.md#0i7-mcp) |
| 実装済み | MCP サーバー | MCP ツール・リソース公開 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.7](DETAIL_INDEX.md#0i7-mcp) |
| 実装済み | MCP サーバー | AI 支援ビルドエラー分析 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.7](DETAIL_INDEX.md#0i7-mcp) |
| 実装済み | MCP サーバー | MCP Prompts 定義 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.7](DETAIL_INDEX.md#0i7-mcp) |
| 実装済み | MCP サーバー | MCP Sampling によるビルドログ自動分析 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.7](DETAIL_INDEX.md#0i7-mcp) |
| 実装済み | MCP サーバー | MCP Notifications（イベントプッシュ） | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.7](DETAIL_INDEX.md#0i7-mcp) |
| 実装済み | MCP サーバー | MCP HTTP SSE transport 対応 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.7](DETAIL_INDEX.md#0i7-mcp) |
| 実装済み | MCP サーバー | MCP ツールスコープ細分化 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.7](DETAIL_INDEX.md#0i7-mcp) |
| 実装済み | MCP サーバー | MCP ツール呼び出し監査ログ | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.7](DETAIL_INDEX.md#0i7-mcp) |
| 実装済み | MCP サーバー | MCP リソース購読（Resource Subscriptions） | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.7](DETAIL_INDEX.md#0i7-mcp) |
| 実装済み | MCP サーバー | MCP クライアント情報ログ | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.7](DETAIL_INDEX.md#0i7-mcp) |
| 実装済み | MCP サーバー | MCP ツール実行統計 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.7](DETAIL_INDEX.md#0i7-mcp) |
| 実装済み | MCP サーバー | MCP ツール実行タイムアウト設定 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.7](DETAIL_INDEX.md#0i7-mcp) |
| 実装済み | MCP サーバー | MCP 設定 CRUD ツール | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.7](DETAIL_INDEX.md#0i7-mcp) |
| 実装済み | MCP サーバー | MCP Elicitation による副作用操作の確認 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.7](DETAIL_INDEX.md#0i7-mcp) |
