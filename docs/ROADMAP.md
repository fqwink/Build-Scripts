# Adlaire CI — Roadmap

[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務は、実装 artifact と各機能の現在状態・未完了理由、Phase の順序・対象 owner・依存関係、機能インベントリ、将来計画を管理する唯一の正本である。状態語彙、状態定義、実装可否、状態遷移条件、`実装済み` への遷移条件は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity)、実装着手条件は [`docs/SPEC.md` 方針責務 §4.7](SPEC.md#sec-4-7)、[`docs/SPEC.md` ポリシー責務 §0d](SPEC.md#policy-spec-freeze)、[`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit)、仕様策定・改訂の完了条件は [`docs/SPEC.md` ポリシー責務 §0b](SPEC.md#policy-spec-pr-completion)、Phase 完了単位は [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit) を参照し、本文書では再定義しない。

<a id="状態計画責務-共通参照先"></a>
**状態・計画責務 共通参照先：**

本文書の状態セルは現在の割り当てだけを示す。endpoint、schema、SDK method、DOM、処理順序、fixture 本文は記載しない。

## 3. 実装 artifact 現在状態

この表は [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3) のディレクトリ構成で定義する実装 artifact だけを対象とする。owner component が持つ個別機能の現在状態は [統合機能インベントリ](#522-統合ロードマップ表) を参照する。実装不一致の技術的な確認内容は [`docs/details/fixture.md` fixture 証跡責務 現行実装整合証跡](details/fixture.md#current-implementation-alignment-evidence) を参照する。

| 実装 artifact | 現在状態 | 未完了理由 / 証跡 |
|-----------------|----------|-------------------|
| [`main.go`](../main.go) | 実装済み | 標準実行バイナリ 8 件の dispatch と version 受け渡しは実装・検証済みである。最終 Phase の実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 15 Obsidian local vault 同期証跡](details/fixture.md#phase-15-obsidian-local-sync-evidence) を参照する。 |
| [`components/builder/`](../components/builder/) | 実装済み | Phase 1 builder 実装は [PR #74](https://github.com/fqwink/Build-Scripts/pull/74) で merge 済みである。Phase 12 で owner package 5 ファイル固定へ移行済みである。実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 1 builder 実装検証証跡](details/fixture.md#phase-1-builder-implementation-evidence) と [`docs/details/fixture.md` fixture 証跡責務 Phase 12 実装品質ゲート再構築証跡](details/fixture.md#phase-12-quality-gate-evidence) を参照する。 |
| [`components/runner/`](../components/runner/) | 実装済み | Phase 2 runner 実装は [PR #75](https://github.com/fqwink/Build-Scripts/pull/75) で実装・検証済みである。Phase 12 で owner package 5 ファイル固定へ移行済みである。実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 2 runner 実装検証証跡](details/fixture.md#phase-2-runner-implementation-evidence) と [`docs/details/fixture.md` fixture 証跡責務 Phase 12 実装品質ゲート再構築証跡](details/fixture.md#phase-12-quality-gate-evidence) を参照する。 |
| [`components/api/`](../components/api/) | 実装済み | Phase 4 `api` operations は実装・検証済みである。Phase 12 で owner package 5 ファイル固定へ移行済みである。実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 4 api operations 実装検証証跡](details/fixture.md#phase-4-api-operations-implementation-evidence) と [`docs/details/fixture.md` fixture 証跡責務 Phase 12 実装品質ゲート再構築証跡](details/fixture.md#phase-12-quality-gate-evidence) を参照する。追加管理 API、SDK、UI、admin、mcp、正式 fixture directory harness の現在状態は各 owner の行と [統合機能インベントリ](#522-統合ロードマップ表) を参照する。 |
| [`components/admin/`](../components/admin/) | 実装済み | Phase 7 `admin` CLI 管理クライアントは実装・検証済みである。Phase 12 で owner package 5 ファイル固定へ移行済みである。実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 7 admin CLI 実装検証証跡](details/fixture.md#phase-7-admin-cli-implementation-evidence) と [`docs/details/fixture.md` fixture 証跡責務 Phase 12 実装品質ゲート再構築証跡](details/fixture.md#phase-12-quality-gate-evidence) を参照する。 |
| [`components/setup/`](../components/setup/) | 実装済み | Phase 8 `setup` は実装・検証済みである。Phase 12 で owner package 5 ファイル固定へ移行済みである。実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 8 setup 実装検証証跡](details/fixture.md#phase-8-setup-implementation-evidence) と [`docs/details/fixture.md` fixture 証跡責務 Phase 12 実装品質ゲート再構築証跡](details/fixture.md#phase-12-quality-gate-evidence) を参照する。 |
| [`components/release/`](../components/release/) | 実装済み | Phase 9 `release` は実装・検証済みである。Phase 12 で owner package 5 ファイル固定へ移行済みである。実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 9 release 実装検証証跡](details/fixture.md#phase-9-release-implementation-evidence) と [`docs/details/fixture.md` fixture 証跡責務 Phase 12 実装品質ゲート再構築証跡](details/fixture.md#phase-12-quality-gate-evidence) を参照する。 |
| [`admin/adlaire-ci-sdk.js`](../admin/adlaire-ci-sdk.js) | 実装済み | Phase 5 SDK 実装は実装・検証済みである。実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 5 sdk 実装検証証跡](details/fixture.md#phase-5-sdk-implementation-evidence) を参照する。 |
| [`admin/index.html`](../admin/index.html) | 実装済み | Phase 6 UI 実装は実装・検証済みである。実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 6 ui 実装検証証跡](details/fixture.md#phase-6-ui-implementation-evidence) を参照する。正式 fixture directory harness の現在状態は [`ALIGN-07`](details/fixture.md#align-07)、[`ALIGN-24`](details/fixture.md#align-24) を参照する。 |
| [`components/mcp/`](../components/mcp/) | 実装済み | Phase 10 `mcp` は実装・検証済みである。Phase 12 で owner package 5 ファイル固定へ移行済みである。実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 10 mcp 実装検証証跡](details/fixture.md#phase-10-mcp-implementation-evidence) と [`docs/details/fixture.md` fixture 証跡責務 Phase 12 実装品質ゲート再構築証跡](details/fixture.md#phase-12-quality-gate-evidence) を参照する。 |

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
| Phase 11 | バグ修正ゼロ化。source-code audit、横断 regression、正式 fixture harness、意味のあるテスト、test gap inventory / batch closure、race trigger、mutation selection / mutation zero survivor、contract drift の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) を参照する。 | 実装済み | Phase 10 |
| Phase 12 | 実装品質ゲート再構築。契約不整合、状態安全性、実行型 fixture harness、mutation / race、owner package 5 ファイル固定、release 再現性の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 12 実装品質ゲート再構築参照](DETAIL_INDEX.md#phase-12-quality-gate-entry) を参照する。 | 実装済み | Phase 11 |
| Phase 13 | 実装整合・品質改善。owner package 5 ファイル責務純度、状態安全性、security / archive / commitstatus 責務集約、queue / finalizer / recovery、MCP 実動作化、外部境界 hardening、fixture / mutation / fault / E2E / CI / release governance の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry) を参照する。 | 実装済み | Phase 12 |
| Phase 14 | Obsidian Vault 連携。local vault 読取、wikilink / embed / tag / asset 正規化、YAML frontmatter 拒否、builder handoff の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 14 Obsidian Vault 連携参照](DETAIL_INDEX.md#phase-14-obsidian-vault-integration-entry) を参照する。 | 実装済み | Phase 13 |
| Phase 15 | Obsidian local vault 同期。plan / apply / rollback、conflict、tombstone、state schema、atomic write、`adlaire-ci-obsidian` 配布連携の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 15 Obsidian local vault 同期参照](DETAIL_INDEX.md#phase-15-obsidian-local-sync-entry) を参照する。 | 実装済み | Phase 14 |
| Phase 16 | 実装済み品質証跡実体化・追加検証候補 closure。宣言型証跡の実行型証跡化、`ALIGN-*` 追加検証候補の分類・完了、ignored error 分類、determinism、filesystem durability、skip、workflow hardening、巨大 owner リスク管理の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 参照](DETAIL_INDEX.md#phase-16-quality-evidence-closure-entry) を参照する。 | 仕様化済み・未実装 | Phase 15 |

現在の active Phase は Phase 16 とする。初期実装 Phase 1 から Phase 15 まではすべて `実装済み` である。Phase 11 は [`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8) のバグ修正ゼロ化そのものを実装・検証する Phase として、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) を入口に扱う。source-code audit により Phase 11 対象へ割り当てた各項目は、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) から対象 owner 詳細本文、[`docs/details/fixture.md` fixture 証跡責務 Phase 11 fixture harness 参照](details/fixture.md#phase-11-fixture-harness-reference)、または [`docs/ROADMAP.md` 状態・計画責務 §5](ROADMAP.md#522-統合ロードマップ表) の対象外理由へ到達する。Phase 11 は [`docs/details/fixture.md` fixture 証跡責務 Phase 11 完了 closure](details/fixture.md#phase-11-quality-gate-closure) に記録された closure record set の `final_open_item_count=0`、mutation `survived=0`、race trigger open `0`、fixture identity duplicate `0` を満たす。Phase 12 は [`docs/details/fixture.md` fixture 証跡責務 Phase 12 実装品質ゲート再構築証跡](details/fixture.md#phase-12-quality-gate-evidence) に記録された closure counter の `final_open_item_count=0`、`phase12_mutation_survived_count=0`、`phase12_race_trigger_open_count=0`、`phase12_owner_package_violation_count=0`、`phase12_ci_required_check_open_count=0` を満たす。Phase 13 は [`docs/details/fixture.md` fixture 証跡責務 Phase 13 実装整合・品質改善証跡](details/fixture.md#phase-13-implementation-alignment-quality-evidence) に記録された全 closure counter の完了値を満たす。Phase 14 は [`docs/details/fixture.md` fixture 証跡責務 Phase 14 Obsidian Vault 連携証跡](details/fixture.md#phase-14-obsidian-vault-integration-evidence) に記録された全 closure counter の完了値を満たす。Phase 15 は [`docs/details/fixture.md` fixture 証跡責務 Phase 15 Obsidian local vault 同期証跡](details/fixture.md#phase-15-obsidian-local-sync-evidence) に記録された全 closure counter の完了値を満たす。Phase 16 は [`docs/details/fixture.md` fixture 証跡責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 証跡](details/fixture.md#phase-16-quality-evidence-closure-evidence) に記録する全 closure counter の完了値を満たすまで `実装済み` に遷移してはならない。active Phase の決定条件と後続 Phase の禁止事項は [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit)、Phase 専用詳細仕様ファイルの扱いは [`docs/SPEC.md` 方針責務 §4.4](SPEC.md#sec-4-4) を参照する。

各行の owner 詳細本文は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0b](DETAIL_INDEX.md#0b-詳細仕様参照表)、機能別の詳細節は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表)、fixture 証跡は [`docs/details/fixture.md` fixture 証跡責務](details/fixture.md) を参照する。

将来計画または改訂予定の現在状態が `仕様化済み・未実装` または `実装中・検証未完了` であっても、active Phase でない対象の新規実装着手可否は [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit) を参照する。

---

## 5. 統合機能インベントリ

状態の意味と遷移条件は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity)、実装詳細の入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) を参照する。以下の表は実装可能な機能、実装中の機能、未実装機能、将来計画を同じ現在状態語彙で管理する。

`改訂予定` の行に記載する詳細入口は、仕様策定先または再評価先を示す。`改訂予定` の詳細入口は、実装着手根拠、Phase 11 実装対象、または `仕様化済み・未実装` の代替根拠として扱わない。`仕様化済み・未実装` の行であっても、実装着手可否は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity)、[`docs/SPEC.md` ポリシー責務 §0d](SPEC.md#policy-spec-freeze)、[`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit) を同時に満たす場合だけ成立する。

<a id="522-統合ロードマップ表"></a>
**統合機能インベントリ：** 以下を全機能の現在状態に関する唯一の一覧とする。

| 現在状態 | 担当領域 | 機能 | 詳細入口 / 次の扱い |
|----------|----------|------|----------------------|
| 実装済み | 管理ツール・SDK | JavaScript SDK 公開契約 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・UI | 標準管理ツール UI 契約 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 状態管理 | 状態ファイル共通永続化契約 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 検証基盤 | バグ修正ゼロ化 / 仕様全般完了 gate / cross-owner regression gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 実装済み | 検証基盤 | 全標準実装 artifact source-code audit / artifact coverage zero gap | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 実装済み | 検証基盤 | 意味のあるテスト / test gap inventory / batch closure / race trigger / mutation selection / mutation zero survivor gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 実装済み | 検証基盤 | 仕様全般不備一括棚卸し / 重複ゼロ / inspection pass / spec-gap closure gate | [`docs/SPEC.md` ポリシー責務 仕様全般不備 inventory record 固定契約](SPEC.md#spec-deficiency-inventory-record-contract)、[`docs/SPEC.md` ポリシー責務 仕様全般不備 batch closure 固定契約](SPEC.md#spec-deficiency-batch-closure-contract)、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 横断テスト証跡共通入口](DETAIL_INDEX.md#cross-test-evidence-route) |
| 実装済み | 検証基盤 | main dispatch / binary version / output manifest / file tree consistency gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 実装済み | 検証基盤 | source-code audit inventory / direct I/O elimination / statefile common persistence gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 実装済み | 検証基盤 | builder parser / config / output / generated site regression gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 実装済み | 検証基盤 | runner source / pipeline / deploy / notification / corrupt state regression gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 実装済み | 検証基盤 | API auth / config / backup / webhook / read model / filesystem regression gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 実装済み | 検証基盤 | SDK / UI / admin CLI / MCP bridge / route parity regression gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 実装済み | 検証基盤 | setup / release distribution / filesystem / external boundary regression gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| 実装済み | 検証基盤 | 正式 fixture directory harness / fixture identity / fixture manifest / contract drift zero gate | [`docs/details/fixture.md` fixture 証跡責務 Phase 11 fixture harness 参照](details/fixture.md#phase-11-fixture-harness-reference) |
| 実装済み | 検証基盤 | Phase 12 契約不整合・状態安全性 gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 12 実装品質ゲート再構築参照](DETAIL_INDEX.md#phase-12-quality-gate-entry) |
| 実装済み | 検証基盤 | Phase 12 実行型 fixture harness / dependency injection / production entrypoint / expected 比較 gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 12 実装品質ゲート再構築参照](DETAIL_INDEX.md#phase-12-quality-gate-entry) |
| 実装済み | 検証基盤 | Phase 12 mutation / race / queue state machine / CI required checks gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 12 実装品質ゲート再構築参照](DETAIL_INDEX.md#phase-12-quality-gate-entry) |
| 実装済み | 実装構造 | Phase 12 owner package 5 ファイル固定 / 巨大コンポーネント分割 gate | [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3)、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 12 実装品質ゲート再構築参照](DETAIL_INDEX.md#phase-12-quality-gate-entry) |
| 実装済み | 配布・リリース | Phase 12 version tag / release notes / release reproducibility gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 12 実装品質ゲート再構築参照](DETAIL_INDEX.md#phase-12-quality-gate-entry) |
| 実装済み | 実装構造 | Phase 13 owner package 5 ファイル責務純度 / 6 ファイル目・サブディレクトリ・空ファイル・ダミー実装ゼロ gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry) |
| 実装済み | 状態管理 | Phase 13 statefile 直接更新排除 / process lock / atomic write / JSON Lines 破損検出・隔離・復旧 gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry) |
| 実装済み | セキュリティ | Phase 13 credential / token / Go 標準ライブラリ内製 KDF / secret mask / file safety / required log write gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry) |
| 実装済み | 管理ツール・API | Phase 13 API route / Admin / SDK / UI client binding / MCP bridge / setup stdout / credential 初期化契約統一 gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry) |
| 実装済み | CI ランナー | Phase 13 queue state machine / finalizer / active recovery / at-least-once / backup-restore transaction gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry) |
| 実装済み | MCP サーバー | Phase 13 MCP resendWebhook / subscribe / unsubscribe / sampling 実動作化と未実装 error gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry) |
| 実装済み | 外部境界 | Phase 13 HTTP timeout / body 上限 / graceful shutdown / Webhook SSRF 防止 / SSH strict / systemd 最小権限 gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry) |
| 実装済み | 検証基盤 | Phase 13 fixture execution / mutation / race / fault injection / integration / E2E / CI required check gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry) |
| 実装済み | 配布・リリース | Phase 13 Git tag / GitHub Release / SHA256SUMS / signature / SBOM / reproducible build evidence / recovery procedure gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry) |
| 実装済み | 文書整合 | Phase 13 旧 path / 重複仕様 / 実装済み表記 drift ゼロ gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry) |
| 実装済み | Obsidian 連携 | Phase 14 Obsidian Vault 連携 / local vault 読取 / wikilink / embed / tag / asset 正規化 / builder handoff gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 14 Obsidian Vault 連携参照](DETAIL_INDEX.md#phase-14-obsidian-vault-integration-entry) |
| 実装済み | Obsidian 同期 | Phase 15 Obsidian local vault 同期 / plan / apply / rollback / conflict / tombstone / atomic write / `adlaire-ci-obsidian` 配布連携 gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 15 Obsidian local vault 同期参照](DETAIL_INDEX.md#phase-15-obsidian-local-sync-entry) |
| 仕様化済み・未実装 | 検証基盤 | Phase 16 source coverage set / evidence package manifest / inventory / expected / record schema / checker 再導出 gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 参照](DETAIL_INDEX.md#phase-16-quality-evidence-closure-entry) |
| 仕様化済み・未実装 | 検証基盤 | Phase 16 宣言型証跡の実行型証跡化 / production entrypoint 実行 / actual expected 比較 gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 参照](DETAIL_INDEX.md#phase-16-quality-evidence-closure-entry) |
| 仕様化済み・未実装 | 検証基盤 | Phase 16 `ALIGN-*` 追加検証候補分類・完了 / 未分類 0 gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 参照](DETAIL_INDEX.md#phase-16-quality-evidence-closure-entry) |
| 仕様化済み・未実装 | 実装品質 | Phase 16 ignored error 分類 / required write failure / panic / skip closure gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 参照](DETAIL_INDEX.md#phase-16-quality-evidence-closure-entry) |
| 仕様化済み・未実装 | 実装品質 | Phase 16 clock / sleep / timeout / entropy determinism gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 参照](DETAIL_INDEX.md#phase-16-quality-evidence-closure-entry) |
| 仕様化済み・未実装 | 状態管理 | Phase 16 filesystem durability parity / atomic write / fsync / parent directory fsync / symlink 非追従 gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 参照](DETAIL_INDEX.md#phase-16-quality-evidence-closure-entry) |
| 仕様化済み・未実装 | CI / 配布 | Phase 16 Phase 12 workflow hardening / Docker 検証手順固定 / Deno stable JavaScript 検証 / release rehearsal gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 参照](DETAIL_INDEX.md#phase-16-quality-evidence-closure-entry) |
| 仕様化済み・未実装 | 保守性 | Phase 16 巨大 owner file risk ledger / 5 ファイル原則維持 / 内部責務区画検査 gate | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 参照](DETAIL_INDEX.md#phase-16-quality-evidence-closure-entry) |
| 改訂予定 | 管理ツール・配布 | 管理 UI 静的配布物構成・Archive 検証 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 改訂予定 | 管理ツール・配布 | 管理 UI 静的 HTTP 配信 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
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
| 改訂予定 | ビルドスクリプト | ビルドキャッシュ | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
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
| 改訂予定 | 管理ツール・API | マルチユーザー対応 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 改訂予定 | 管理ツール・API | データストア切り替え | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 改訂予定 | 管理ツール・API | 外部認証連携 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | TOTP 二要素認証 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 改訂予定 | 管理ツール・API | 統計データの JSON エクスポート | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 改訂予定 | 管理ツール・API | キュー内個別エントリのキャンセル | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | 設定バリデーション API | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 改訂予定 | 管理ツール・API | Prometheus メトリクスエンドポイント | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | CLI 管理クライアント | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 改訂予定 | 管理ツール・API | 設定の自動スナップショット | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 改訂予定 | 管理ツール・API | ステータスバッジ生成 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 改訂予定 | 管理ツール・API | ビルド履歴の自動削除設定 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | セッションタイムアウト変更設定 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | ビルドトリガー専用 API スコープ | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | 監査ログ | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | API レート制限 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 改訂予定 | 管理ツール・API | ロールベースアクセス制御 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 改訂予定 | 管理ツール・API | 設定スナップショット差分表示 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 改訂予定 | 管理ツール・API | 複数プロジェクト管理 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | API キー管理 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | ビルドログの保存済み有限 SSE 配信 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | ビルド統計ダッシュボード | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 改訂予定 | 管理ツール・API | ユーザー管理 API | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | IP アドレス制限 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 改訂予定 | 管理ツール・API | API バージョニング | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 改訂予定 | 管理ツール・API | API ドキュメント自動生成 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | 通知チャンネル管理 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 改訂予定 | 管理ツール・API | ビルドキューの手動並び替え | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 改訂予定 | 管理ツール・API | 設定テンプレート | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | API アクセスログ | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 改訂予定 | 管理ツール・API | 管理者向けイベントフィード | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | ビルドキュー可視化 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | メンテナンスモード | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 改訂予定 | 管理ツール・API | 読み取り専用共有リンク | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 実装済み | 管理ツール・API | アラート閾値設定 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 実装済み | 管理ツール・API | バックアップ／リストア | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](DETAIL_INDEX.md#0i-詳細節対応表) |
| 改訂予定 | 管理ツール・API | API レスポンスキャッシュ制御 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 改訂予定 | 管理ツール・API | スナップショット間サイト差分 API | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
| 改訂予定 | 管理ツール・API | Webhook 送信履歴の手動再送 API | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.6](DETAIL_INDEX.md#0i6-追加管理api機能) |
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
