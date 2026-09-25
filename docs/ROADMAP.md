# Adlaire CI — Roadmap

[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務は、Adlaire CI の実装 artifact と各機能へ割り当てる現在状態、Phase、機能インベントリ、将来計画を管理する正本である。owner component と実装 artifact の区別、状態語彙、状態定義、実装可否、着手条件、状態遷移条件、完了条件は [`docs/SPEC.md` 方針責務 §4.2a](SPEC.md#sec-4-2a) と [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) を正本とし、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務では再定義しない。

方針・ポリシーは [`docs/SPEC.md`](SPEC.md)、生成静的 Web サイトと標準管理 UI のデザインは [`docs/DESIGN.md`](DESIGN.md)、詳細仕様入口・owner 対応表・collaborator 境界参照入口は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、owner component 別詳細本文は [`docs/details/`](details/)、fixture と実装検証証跡は [`docs/details/fixture.md`](details/fixture.md)、実在所在は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) を参照する。

<a id="状態計画責務-共通参照先"></a>
**状態・計画責務 共通参照先：**

現在状態は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) の判定条件を適用して [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務へ割り当てる。[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務は割当結果と未完了理由だけを記録し、endpoint、schema、SDK method、DOM、処理順序、fixture 本文を再掲しない。

## 1. ロードマップ責務

| 管理対象 | 本文で管理する内容 |
|----------|------------------|
| 現在状態 | 実装 artifact と機能へ割り当てた現在の状態、および未完了理由。 |
| Phase | Phase 順序、対象 owner component、依存関係、現在状態。 |
| 機能インベントリ | 全機能の担当領域と現在状態。 |
| 将来計画 | 将来計画へ割り当てた機能名と担当領域。 |

## 2. 状態語彙参照

状態語彙と判定条件は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) を参照する。[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務の状態セルは現在の割り当てだけを示す。

## 3. 実装 artifact 現在状態

この表は [`docs/SPEC.md` 方針責務 §4.2a](SPEC.md#sec-4-2a) の実装 artifact だけを対象とする。owner component が持つ個別機能の現在状態は [統合機能インベントリ](#522-統合ロードマップ表) を参照する。実装不一致の技術的な確認内容は [`docs/details/fixture.md` fixture 証跡責務 現行実装整合証跡](details/fixture.md#current-implementation-alignment-evidence) を参照する。

| 実装 artifact | 現在状態 | 未完了理由 / 証跡 |
|-----------------|----------|-------------------|
| [`main.go`](../main.go) | 実装中・検証未完了 | 起動入口の必須実装・証跡が未完了である。[`ALIGN-01`](details/fixture.md#align-01) |
| [`components/builder.go`](../components/builder.go) | 実装中・検証未完了 | デザイン整合、Markdown block 契約、report 契約、必須 fixture が未完了である。[`ALIGN-06`](details/fixture.md#align-06)、[`ALIGN-07`](details/fixture.md#align-07)、[`ALIGN-08`](details/fixture.md#align-08)、[`ALIGN-13`](details/fixture.md#align-13) |
| [`components/runner.go`](../components/runner.go) | 実装中・検証未完了 | owner 契約の必須実装・証跡が未完了である。[`ALIGN-07`](details/fixture.md#align-07)、[`ALIGN-09`](details/fixture.md#align-09)、[`ALIGN-12`](details/fixture.md#align-12) |
| [`components/api.go`](../components/api.go) | 実装中・検証未完了 | owner 契約の必須実装・証跡が未完了である。[`ALIGN-01`](details/fixture.md#align-01)〜[`ALIGN-04`](details/fixture.md#align-04)、[`ALIGN-07`](details/fixture.md#align-07)、[`ALIGN-09`](details/fixture.md#align-09)〜[`ALIGN-12`](details/fixture.md#align-12) |
| [`admin/adlaire-ci-sdk.js`](../admin/adlaire-ci-sdk.js) | 実装中・検証未完了 | owner 契約の必須実装・証跡が未完了である。[`ALIGN-03`](details/fixture.md#align-03)、[`ALIGN-07`](details/fixture.md#align-07)、[`ALIGN-10`](details/fixture.md#align-10) |
| [`admin/index.html`](../admin/index.html) | 実装中・検証未完了 | owner 契約の必須実装・証跡が未完了である。[`ALIGN-05`](details/fixture.md#align-05)、[`ALIGN-07`](details/fixture.md#align-07)、[`ALIGN-10`](details/fixture.md#align-10) |
| `components/mcp.go` | 将来計画 | 実装ファイルと MCP 専用詳細仕様が存在しない。 |

## 4. Phase 実装計画

Phase の実装単位、禁止事項、着手条件、完了判定方針は [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#0f-phase-実装単位ポリシー) を参照する。[`docs/ROADMAP.md` 状態・計画責務 §4](ROADMAP.md#4-phase-実装計画) は Phase の割り当て、現在状態、順序、依存関係だけを管理する。

## 4.1 初期実装 Phase 単位

| Phase | owner / scope | 現在状態 | 依存する Phase |
|-------|---------------|----------|----------------|
| Phase 1 | `builder` | 実装中・検証未完了 | なし |
| Phase 2 | `runner` | 実装中・検証未完了 | Phase 1 |
| Phase 3 | `api` request lifecycle | 実装中・検証未完了 | Phase 2 |
| Phase 4 | `api` operations | 実装中・検証未完了 | Phase 3 |
| Phase 5 | `sdk` | 実装中・検証未完了 | Phase 4 |
| Phase 6 | `ui` | 実装中・検証未完了 | Phase 5 |

現在の active Phase は `Phase 1` である。active Phase の決定条件と後続 Phase の禁止事項は [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#0f-phase-実装単位ポリシー) を参照する。

各行の owner 詳細本文は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0b](DETAIL_INDEX.md#0b-詳細仕様参照表)、機能別の詳細節は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表)、fixture 証跡は [`docs/details/fixture.md`](details/fixture.md) を参照する。

`setup` と `release` の機能は現在状態が `改訂予定` であり、現在は Phase 未割当である。状態遷移条件は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー)、Phase 追加条件は [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#0f-phase-実装単位ポリシー) を参照する。

---

## 5. 統合機能インベントリ

状態の意味と遷移条件は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー)、実装詳細の入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) を参照する。以下の表は実装可能な機能、実装中の機能、未実装機能、将来計画を同じ現在状態語彙で管理する。

<a id="522-統合ロードマップ表"></a>
**統合機能インベントリ：** 以下を全機能の現在状態に関する唯一の一覧とする。

| 現在状態 | 担当領域 | 機能 | 詳細入口 / 次の扱い |
|----------|----------|------|----------------------|
| 実装中・検証未完了 | 管理ツール・SDK | JavaScript SDK 公開契約 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | 管理ツール・UI | 標準管理ツール UI 契約 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | 状態管理 | 状態ファイル共通永続化契約 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 改訂予定 | 管理ツール・配布 | 管理 UI 静的配布物構成・Archive 検証 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表)。setup 実行 artifact の path、起動名、入力 interface を確定するまで実装不可。 |
| 仕様化済み・未実装 | 管理ツール・配布 | 管理 UI 静的 HTTP 配信 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | ビルドタイムアウト | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | ポーリング間隔の動的変更 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | ビルドログのファイル保存 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | GitHub Webhook 受信 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | ネットワーク断時の再試行 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | GitHub API レート制限自動待機 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | 転送後リモート整合性検証 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | マルチブランチビルド | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | ビルドログ世代管理 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | ビルド出力の外部転送 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | ビルドクールダウン | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | ビルド前の事前チェック | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | 定期強制ビルド | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | ビルド中重複スキップ | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | GitHub PAT 有効期限の事前警告 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | コミット情報のビルドログ記録 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | GitHub API 連続失敗によるサーキットブレーカー | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | 出力サイトサイズ警告閾値 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | Webhook イベントログ | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | ビルド所要時間の記録と統計 API | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | ビルドアーティファクト世代管理 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | 管理ツール・API | ビルドアーティファクト管理 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | 管理ツール・API | ヘルスチェックエンドポイント | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | 管理ツール・API | Webhook イベント一覧取得 API | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | 管理ツール・API | ビルドログ重大度フィルター | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | 変換レポート出力 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | シンタックスハイライト | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | 本文内全文検索 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | アンカーリンク自動検証 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | コードブロックの折りたたみ | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | 印刷スタイル（`@media print`） | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | 静的 Web サイト出力 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | テーマコンポーネント | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | 外部リンクの自動処理 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | 読み取り進捗バー | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | コードブロックのコピーボタン | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | 見出しアンカーリンクコピー | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | TOC 開閉状態の永続化 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | 見出しスラグ重複解決 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | 前後章ナビゲーションボタン | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | 内部リンク整合性チェック | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | 見出し階層スキップ警告 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | 読了時間推計と表示 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | Webhook 通知失敗リトライキュー | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | ブランチ設定の動的変更 API | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | 週次ビルドサマリー Webhook | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | 管理ツール・API | 設定変更の詳細 diff 記録 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | テーブルのソート機能 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | キーボードショートカット | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | 複数ファイル監視 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | GitHub Commit Status API | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | ビルドパイプライン YAML 定義 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | ドライラン実行モード | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | ビルドログのアーカイブ圧縮 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | ローカルファイル監視モード | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | タグ付きコミットのみビルド | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | ビルドキャッシュ | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | ビルド通知連携 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | ビルド時間トレンド記録 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | ビルド失敗時の自動リトライ | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | 並列マルチターゲットビルド | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | ビルド前後フック | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | 依存ファイルトラッキング | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | リモートビルド対応 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | ビルドステータスファイル出力 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | ビルド承認フロー | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | ブランチ別環境変数 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | ビルド依存チェーン | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | CI ランナー | ビルド優先度キュー | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | 失敗原因の自動分類 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | ビルド実行環境の記録 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | ビルドトリガー種別の記録 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | ビルド所要時間の異常検知 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | CI ランナー | 設定ファイル起動時整合性チェック | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 将来計画 | 管理ツール・API | マルチユーザー対応 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | 管理ツール・API | データストア切り替え | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | 管理ツール・API | 外部認証連携 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 実装中・検証未完了 | 管理ツール・API | TOTP 二要素認証 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 将来計画 | 管理ツール・API | 統計データの JSON エクスポート | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | 管理ツール・API | キュー内個別エントリのキャンセル | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 実装中・検証未完了 | 管理ツール・API | 設定バリデーション API | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 将来計画 | 管理ツール・API | Prometheus メトリクスエンドポイント | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | 管理ツール・API | CLI 管理クライアント | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | 管理ツール・API | 設定の自動スナップショット | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | 管理ツール・API | ステータスバッジ生成 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | 管理ツール・API | ビルド履歴の自動削除設定 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 実装中・検証未完了 | 管理ツール・API | セッションタイムアウト変更設定 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | 管理ツール・API | ビルドトリガー専用 API スコープ | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | 管理ツール・API | 監査ログ | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | 管理ツール・API | API レート制限 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 将来計画 | 管理ツール・API | ロールベースアクセス制御 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | 管理ツール・API | 設定スナップショット差分表示 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | 管理ツール・API | 複数プロジェクト管理 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 実装中・検証未完了 | 管理ツール・API | API キー管理 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | 管理ツール・API | ビルドログの保存済み有限 SSE 配信 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | 管理ツール・API | ビルド統計ダッシュボード | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 将来計画 | 管理ツール・API | ユーザー管理 API | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 実装中・検証未完了 | 管理ツール・API | IP アドレス制限 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 将来計画 | 管理ツール・API | API バージョニング | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | 管理ツール・API | API ドキュメント自動生成 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 実装中・検証未完了 | 管理ツール・API | 通知チャンネル管理 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 将来計画 | 管理ツール・API | ビルドキューの手動並び替え | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | 管理ツール・API | 設定テンプレート | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 実装中・検証未完了 | 管理ツール・API | API アクセスログ | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 将来計画 | 管理ツール・API | 管理者向けイベントフィード | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 実装中・検証未完了 | 管理ツール・API | ビルドキュー可視化 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | 管理ツール・API | メンテナンスモード | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 将来計画 | 管理ツール・API | 読み取り専用共有リンク | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 実装中・検証未完了 | 管理ツール・API | アラート閾値設定 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | 管理ツール・API | バックアップ／リストア | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 将来計画 | 管理ツール・API | API レスポンスキャッシュ制御 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | 管理ツール・API | スナップショット間サイト差分 API | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | 管理ツール・API | Webhook 送信履歴の手動再送 API | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 仕様化済み・未実装 | ビルドスクリプト | 差分ビルド | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | 複数出力形式 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | Markdown 拡張記法サポート | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | コードブロック行番号表示 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | 見出しの自動採番 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | セクション折りたたみ | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | TOC 深さ制御 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | 最終更新日の自動埋め込み | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | diff ハイライト | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | 画像の遅延読み込み | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | カスタムメタタグ注入 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | ライトモード固定 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | コードブロックのファイル名表示 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | テンプレート変数展開 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | HTML ミニファイ | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | TOC ハイライト追従 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | Mermaid ダイアグラム描画 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | 脚注サポート | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | インライン数式レンダリング | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | ページ内ナビゲーション履歴 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | 読み上げ対応（アクセシビリティ） | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | 画像ライトボックス | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | 出力サイトへのビルドメタ埋め込み | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | 印刷時 QR コード挿入 | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 仕様化済み・未実装 | ビルドスクリプト | 定義リストサポート | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装中・検証未完了 | ビルドスクリプト | タスクリストサポート | [`docs/DETAIL_INDEX.md` §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 改訂予定 | 配布・セットアップ | 初回セットアップ | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.5](DETAIL_INDEX.md#0i5-setup--release) |
| 改訂予定 | 配布・セットアップ | 管理 API 導入 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.5](DETAIL_INDEX.md#0i5-setup--release) |
| 改訂予定 | 配布・セットアップ | バイナリアップデート | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.5](DETAIL_INDEX.md#0i5-setup--release) |
| 改訂予定 | 配布・セットアップ | アップデート rollback | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.5](DETAIL_INDEX.md#0i5-setup--release) |
| 改訂予定 | リリース | GitHub Release 成果物生成・公開前検証・公開 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.5](DETAIL_INDEX.md#0i5-setup--release) |
| 将来計画 | MCP サーバー | MCP サーバー実装 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | MCP サーバー | MCP ツール・リソース公開 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | MCP サーバー | AI 支援ビルドエラー分析 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | MCP サーバー | MCP Prompts 定義 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | MCP サーバー | MCP Sampling によるビルドログ自動分析 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | MCP サーバー | MCP Notifications（イベントプッシュ） | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | MCP サーバー | MCP HTTP SSE transport 対応 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | MCP サーバー | MCP ツールスコープ細分化 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | MCP サーバー | MCP ツール呼び出し監査ログ | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | MCP サーバー | MCP リソース購読（Resource Subscriptions） | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | MCP サーバー | MCP クライアント情報ログ | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | MCP サーバー | MCP ツール実行統計 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | MCP サーバー | MCP ツール実行タイムアウト設定 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | MCP サーバー | MCP 設定 CRUD ツール | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
| 将来計画 | MCP サーバー | MCP Elicitation による副作用操作の確認 | [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) |
