# Adlaire CI — 詳細仕様入口

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務は、実装者が owner component 別詳細本文責務へ到達するための入口、共通固定値、実装前確認項目、検証マトリクス、詳細節対応表、リポジトリ内ソース配置だけを扱う。責務分離と重複禁止の方針は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.2a を正本とする。

方針、ポリシー、正本参照先は [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務、状態分類、実装可否、Phase、将来計画は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、各 owner component の入出力、状態、処理順序、異常系、検証条件は owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を参照する。

---

## 詳細仕様管理

詳細仕様入口責務は、管理対象と責務を持つ正本への確認入口を示す。

| 管理対象 | 参照先 | 詳細仕様入口責務での扱い |
|----------|------|--------------------|
| 詳細仕様参照入口 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 実装者が詳細仕様本文へ到達するための入口を示す。 |
| 共通固定値 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | component 間で共有する固定値だけを示す。 |
| 詳細節対応表 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 対象機能、詳細仕様節、受け入れ条件の入口を示す。 |
| 状態分類、実装可否、Phase、将来計画 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 | 確認入口。 |
| [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6 追加仕様化機能参照 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6 | 確認入口。 |
| owner component 本文 | owner component 別の [`docs/details/*.md`](details/) 詳細本文責務 | 参照入口。 |
| 方針、ポリシー、正本参照先 | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | 確認入口。 |

## 詳細仕様入口責務
> 実装の具体的詳細へ到達するための入口を定める。「どこから詳細仕様を読むか」に答える。

---

## 詳細仕様参照入口

詳細仕様を読む前に、[`AGENTS.md`](../AGENTS.md) と [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務を必ず確認する。文書所在は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務を確認する。

対象機能ごとの参照入口は以下とする。

1. [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務で対象の状態分類、実装可否、Phase、将来計画該当有無を確認する。
2. 生成 HTML のデザイン関係を扱う場合は、[`docs/DESIGN.md`](DESIGN.md) デザイン責務で視覚仕様を確認する。
3. [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b.0 と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b で owner component を確定する。
4. [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i で、対象機能に対応する詳細仕様節と受け入れ条件を特定する。
5. [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0a〜§0h で、記載基準、共通固定値、実装前確認項目、検証入口、Phase 詳細仕様参照を確認する。
6. owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を主本文として読み、入力、出力、状態、正常系、異常系、セキュリティ、検証条件を確認する。
7. collaborator component がある場合は、該当する [`docs/details/*.md`](details/) 詳細本文責務を呼び出し境界、schema、表示、security、setup、fixture、検証観点として確認する。
8. [`docs/details/setup.md`](details/setup.md) §26 のセットアップ・アップデート手順と [`docs/details/setup.md`](details/setup.md) §26.7 の受け入れ条件に影響がある場合は、実装変更の検証対象に含める。

詳細仕様節に [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h の項目不足がある場合は、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.7 と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0d を参照する。[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務では、不足を検出した参照入口と対応表の更新要否だけを扱う。

| 範囲 | 役割 |
|------|------|
| [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0〜§0j | 詳細仕様の記載基準、実装前確認項目、共通固定値、検証、詳細節対応表、リポジトリ内ソース配置 |
| [`docs/details/builder.md`](details/builder.md) §1〜§9 | [`docs/details/builder.md`](details/builder.md)。`builder` / `adlaire-ci-build` の詳細仕様 |
| [`docs/details/runner.md`](details/runner.md) §10〜§20 | [`docs/details/runner.md`](details/runner.md)。`runner` / `adlaire-ci-runner` の詳細仕様 |
| [`docs/details/api.md`](details/api.md) §21〜§22 | [`docs/details/api.md`](details/api.md)。`api` / `adlaire-ci-api` の詳細仕様。[`docs/details/api.md`](details/api.md) §21a は管理 API サーバー制限を定義する。 |
| [`docs/details/sdk.md`](details/sdk.md) §23 | [`docs/details/sdk.md`](details/sdk.md)。`sdk` の詳細仕様 |
| [`docs/details/ui.md`](details/ui.md) §24 | [`docs/details/ui.md`](details/ui.md)。`ui` の詳細仕様 |
| [`docs/details/api.md`](details/api.md) §25 | [`docs/details/api.md`](details/api.md)。認証の実装仕様 |
| [`docs/details/setup.md`](details/setup.md) §26 | [`docs/details/setup.md`](details/setup.md)。バイナリ配布前提のセットアップ、アップデート、受け入れ条件 |
| `admin` | [`docs/details/admin.md`](details/admin.md)。管理 UI 静的ファイルの配布物構成、配置、HTTP 静的配信境界 |

owner component 別の [`docs/details/*.md`](details/) 詳細本文の owner / collaborator 境界は、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b.1 に従う。

## 0a. 詳細仕様の記載基準

詳細仕様項目の粒度は、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.4 と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0 を参照する。詳細仕様入口責務は、owner component 別詳細本文責務で確認する項目の入口表である。

仕様項目を追加または改訂する場合は、対象範囲に応じて以下の項目を owner component 別の [`docs/details/*.md`](details/) 詳細本文責務で確認する。

| 項目 | 記載する内容 |
|------|-------------|
| 対象 | owner component、collaborator component、対象ファイル、対象機能、責務境界 |
| 入力 | 引数、HTTP リクエスト、設定値、読み込みファイル、環境前提 |
| 出力 | 戻り値、HTTP レスポンス、生成ファイル、ログ、通知、標準出力 |
| 状態 | 状態ファイル、メモリ上の状態、ロック、キャッシュ、更新タイミング |
| データ構造 | JSON キー、型、必須/任意、既定値、許容値、例 |
| 処理順序 | 正常系フロー、分岐条件、ループ条件、疑似コード |
| 異常系 | エラー条件、例外処理、ステータスコード、ログレベル、通知条件 |
| 運用条件 | タイムアウト、再試行、冪等性、排他制御、世代管理、削除条件 |
| セキュリティ | 認証、認可、秘密情報の保存禁止、権限、外部公開可否 |
| 検証 | 構文確認、実行確認、API 確認、生成物確認、整合性確認 |

未確定内容、実装可能性、推測補完の扱いは [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0 と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0a を参照する。[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務では、未確定項目の検出先と参照先を示す。

対象範囲の具体化先は、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務とする。

---

## 0b.0 詳細仕様選択フロー

対象機能の詳細仕様を読む場合は、以下の順で参照先を確定する。

1. [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b.0.1 のカテゴリ表で、対象機能が属する大まかな領域を特定する。
2. [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b の詳細仕様参照表で、owner component を 1 件に確定する。
3. owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を主本文として読む。
4. [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b の参照表または owner component 本文に collaborator component が示されている場合だけ、該当する [`docs/details/*.md`](details/) 詳細本文責務を読む。
5. 状態ファイル、認証、fixture、setup、archive、Commit Status など横断参照が必要な場合は、owner component 本文から明示された参照先に限定して読む。

カテゴリは入口であり、owner component ではない。owner component の確定は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b の詳細仕様参照表で行う。

## 0b.0.1 詳細仕様カテゴリ

カテゴリ索引は、実装者が対象機能から読むべき owner component 別の詳細本文責務を特定するための検索入口である。カテゴリは owner component、状態分類、実装可否、ロードマップ状態を変更しない。

| カテゴリ | 対象 component | 主な確認対象 | 読む詳細仕様 |
|----------|----------------|--------------|--------------|
| Build 系 | `builder` | Markdown 入力、静的 Web サイト出力、theme component、検索 index、変換 report。 | [`docs/details/builder.md`](details/builder.md) |
| CI / 運用系 | `runner`、`commitstatus`、`archive` | GitHub 監視、pipeline、deploy、snapshot、通知、Commit Status、archive / rollback / cleanup。 | [`docs/details/runner.md`](details/runner.md)、[`docs/details/commitstatus.md`](details/commitstatus.md)、[`docs/details/archive.md`](details/archive.md) |
| 管理系 | `api`、`sdk`、`ui`、`admin` | 管理 API、JavaScript SDK、標準管理 UI、管理 UI 静的ファイル配布と配信境界。 | [`docs/details/api.md`](details/api.md)、[`docs/details/sdk.md`](details/sdk.md)、[`docs/details/ui.md`](details/ui.md)、[`docs/details/admin.md`](details/admin.md) |
| 状態 / 安全系 | `statefile`、`security` | 状態ファイル、lock、atomic write、schema、token、scope、audit、session、TOTP、rate limit。 | [`docs/details/statefile.md`](details/statefile.md)、[`docs/details/security.md`](details/security.md) |
| 配布 / 検証系 | `setup`、`fixture` | バイナリ配布、systemd、セットアップ、アップデート、fixture、fake、実装検証証跡、acceptance checklist。 | [`docs/details/setup.md`](details/setup.md)、[`docs/details/fixture.md`](details/fixture.md) |
| 将来計画 | `mcp` | MCP サーバー。現時点では実装可能な詳細仕様を持たない。 | 未定義。将来計画状態は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.2 を確認する。 |

カテゴリをまたぐ機能では、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b の owner component 参照表を優先して owner component を 1 件に確定する。

---

## 0b. 詳細仕様参照表

詳細仕様参照入口は、owner component ごとに参照する詳細仕様節を示す。状態分類、実装可否、ロードマップ状態は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を確認する。

| owner component | 詳細仕様節 | 主な確認対象 |
|--------------------|------------|--------------|
| `builder` | [`docs/details/builder.md`](details/builder.md) §1〜§9、[`docs/details/builder.md`](details/builder.md) §8a、[`docs/details/builder.md`](details/builder.md) §27.4、[`docs/details/builder.md`](details/builder.md) §27.25、[`docs/details/builder.md`](details/builder.md) §27.28 | CLI、入力 Markdown、静的 Web サイト出力、HTML / CSS / JavaScript、theme component、変換 report、builder fixture、builder owner 追加機能。 |
| `runner` | [`docs/details/runner.md`](details/runner.md) §10〜§20、[`docs/details/runner.md`](details/runner.md) §15a、[`docs/details/runner.md`](details/runner.md) §27.2〜§27.3、[`docs/details/runner.md`](details/runner.md) §27.8〜§27.10、[`docs/details/runner.md`](details/runner.md) §27.14、[`docs/details/runner.md`](details/runner.md) §27.19、[`docs/details/runner.md`](details/runner.md) §27.21〜§27.24、[`docs/details/runner.md`](details/runner.md) §27.26〜§27.27、[`docs/details/runner.md`](details/runner.md) §27.29、[`docs/details/runner.md`](details/runner.md) §27.31〜§27.38、[`docs/details/commitstatus.md`](details/commitstatus.md) §27.1、[`docs/details/setup.md`](details/setup.md) §26 | GitHub 監視、設定読取、状態ファイル更新呼び出し、pipeline、deploy、snapshot 作成トリガー、通知、runner fixture、runner owner 追加機能、Commit Status 呼び出し境界、systemd / setup 参照境界。 |
| `api` | [`docs/details/api.md`](details/api.md) §21〜§22、[`docs/details/api.md`](details/api.md) §21a、[`docs/details/api.md`](details/api.md) §25、[`docs/details/api.md`](details/api.md) §27.5〜§27.6、[`docs/details/api.md`](details/api.md) §27.11〜§27.13、[`docs/details/api.md`](details/api.md) §27.16〜§27.18、[`docs/details/api.md`](details/api.md) §27.20、[`docs/details/api.md`](details/api.md) §27.30、[`docs/details/api.md`](details/api.md) §27.42〜§27.47、[`docs/details/security.md`](details/security.md) §27.42〜§27.47、[`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/statefile.md`](details/statefile.md) §22.0c | HTTP 共通契約、API server 制限、endpoint、request / response、状態ファイル read/write 呼び出し境界、認証連携、security 呼び出し境界、api owner 追加機能。 |
| `admin` | [`docs/details/admin.md`](details/admin.md) §0、A1〜A6、[`docs/details/setup.md`](details/setup.md) §26.8、[`docs/details/fixture.md`](details/fixture.md) §27-F setup / admin / release 連動 fixture 固定契約 | 管理 UI 静的ファイルの配布物構成、配置、検証、HTTP 静的配信境界、setup / admin / release 連動 fixture。 |
| `sdk` | [`docs/details/sdk.md`](details/sdk.md) §23 | SDK class、method、HTTP 対応、query / body 生成、error、stream、token 破棄。 |
| `ui` | [`docs/details/ui.md`](details/ui.md) §24 | 画面構成、DOM id、panel、SDK 呼び出し、表示状態、秘密情報消去。 |
| `setup` | [`docs/details/setup.md`](details/setup.md) §26 | バイナリ配布、配置、systemd、セットアップ、アップデート、リリース成果物検証。 |
| `statefile` | [`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/statefile.md`](details/statefile.md) §22.0c | 状態ファイル共通仕様、lock、atomic write、JSON Lines、破損時処理、状態読取 adapter、主要 schema。 |
| `archive` | [`docs/details/archive.md`](details/archive.md) §27.7、[`docs/details/archive.md`](details/archive.md) §27.15 | build log archive、snapshot、download、delete、rollback、cleanup。 |
| `commitstatus` | [`docs/details/commitstatus.md`](details/commitstatus.md) §27.1 | GitHub Commit Status API payload、送信順、失敗時非反転、保存値、secret mask、検証条件。 |
| `security` | [`docs/details/security.md`](details/security.md) §27.42〜§27.47 | API token scope、API key、audit、session timeout、TOTP、rate limit、漏えい禁止、security 横断順序。 |
| `mcp` | 未定義。将来計画状態は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.2 を確認する。 | 将来計画。現時点では実装可能な入出力、状態、起動手順、ツール定義、検証条件の本文を持たない。 |

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b の詳細仕様参照表の `詳細仕様節` は参照入口であり、主本文の owner component を変更しない。複数ファイルを参照する行では、対象機能の owner component のファイルを主本文とし、他ファイルは collaborator の境界、schema、fixture、security、setup、受け入れ条件を確認するために読む。同じ HTTP body、状態 schema、DOM id、SDK method、fixture assertion の重複定義禁止は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.2a を参照する。

---

## 0b.1 owner component 別 owner / collaborator 境界管理

owner / collaborator 境界管理は、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務の境界を維持する。責務境界の変更、詳細本文配置変更、参照先更新では、方針・ポリシーは [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務、状態分類・実装可否・ロードマップ状態は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を確認する。

owner component 別の [`docs/details/*.md`](details/) 詳細本文責務は、owner component の主本文を持つ。方針、ポリシー、正本参照先は [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務、状態分類、ロードマップ状態、実装可否は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を確認する。各 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務の冒頭では、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b.1 への参照と自ファイルの owner / collaborator 境界だけを示す。

詳細仕様本文の配置単位は owner component を第一基準とする。複数 component が関わる機能は、owner component のファイルに主本文を置き、collaborator component のファイルには呼び出し境界、schema、表示、security、setup、fixture、検証観点だけを置く。

owner component 別の [`docs/details/*.md`](details/) 詳細本文責務は以下に固定する。`COMMON`、`CORE`、`BASE`、`SHARED`、`FOUNDATION`、その他の横断共通基盤ファイルの作成禁止は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.2 と [`docs/SPEC.md`](SPEC.md) 方針責務 §4.2a を参照する。

| ファイル | 持つ内容 | 持たない内容 |
|----------|----------|--------------|
| [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 詳細仕様参照入口、共通固定値、実装前確認項目、検証マトリクス、詳細節対応表、リポジトリ内ソース配置、owner component 別詳細本文責務の owner / collaborator 境界管理。 | 各 component の詳細な処理本文、fixture 証跡責務、状態 schema、API endpoint 詳細、SDK method 実装、UI DOM 詳細、setup / release 手順。 |
| [`docs/details/builder.md`](details/builder.md) | `builder` owner の Markdown 変換、静的 Web サイト出力、HTML / CSS / JavaScript、theme component、builder fixture、builder owner 追加機能。 | GitHub read、runner 状態更新、API endpoint、SDK method 実装、UI DOM 詳細、状態 schema、admin 静的配信、setup / release 手順、fixture 証跡責務。 |
| [`docs/details/runner.md`](details/runner.md) | `runner` owner の GitHub 監視、設定読取、状態ファイル更新呼び出し、pipeline、deploy、snapshot 作成トリガー、通知、runner fixture、runner owner 追加機能。 | API endpoint の認証・応答本文、SDK method 実装、UI DOM 詳細、builder の変換処理、admin 静的配信、security 主本文、状態 schema、setup / release 手順、fixture 証跡責務。 |
| [`docs/details/api.md`](details/api.md) | `api` owner の HTTP 共通契約、endpoint、request / response、状態ファイル read/write 呼び出し境界、認証連携、api owner 追加機能。 | SDK method 実装、UI DOM 詳細、runner の build 実行責務、builder の変換処理、admin 静的配信、security 主本文、状態 schema、setup / release 手順、fixture 証跡責務。 |
| [`docs/details/admin.md`](details/admin.md) | `admin` owner の管理 UI 静的ファイル配布物構成、配置、検証、HTTP 静的配信境界。 | UI DOM 詳細、SDK method 実装、API endpoint 実装、状態 schema、systemd 導入手順、release asset 取得手順、fixture 証跡責務。 |
| [`docs/details/sdk.md`](details/sdk.md) | `sdk` owner の SDK class、method、HTTP 対応、query / body 生成、error、stream、token 破棄。 | API endpoint 実装、API endpoint の状態ファイル更新責務、UI DOM 詳細、状態 schema、状態ファイル直接操作、admin 静的配信、setup / release 手順、fixture 証跡責務。 |
| [`docs/details/ui.md`](details/ui.md) | `ui` owner の DOM id、panel、操作、表示状態、SDK 呼び出し、秘密情報消去。 | SDK method 実装、API endpoint 実装、状態 schema、状態ファイル直接操作、admin 静的配信、setup / release 手順、fixture 証跡責務。 |
| [`docs/details/setup.md`](details/setup.md) | `setup` owner のバイナリ配布、配置、systemd、セットアップ、アップデート、リリース成果物検証。 | runner / api / sdk / ui / admin の個別機能本文、状態 schema、API endpoint、SDK method、UI DOM、fixture 証跡責務、外部依存追加。 |
| [`docs/details/statefile.md`](details/statefile.md) | `statefile` owner の状態ファイル共通仕様、lock、atomic write、JSON Lines、破損時処理、状態読取 adapter、主要 schema。 | API endpoint の request / response、runner の業務処理、SDK method 実装、UI 表示判断、setup / release 手順、fixture 証跡責務、個別 component の業務判断。 |
| [`docs/details/archive.md`](details/archive.md) | `archive` owner の build log archive / cleanup の実体処理、snapshot 保存形式、download tar.gz 生成安全性、snapshot delete 実体処理、rollback 転送実体処理。 | runner の通常 build 実行、snapshot 作成トリガー判定、API 共通 request / response、SDK method 実装、UI DOM 詳細、状態 schema、setup / release 手順、fixture 証跡責務。 |
| [`docs/details/commitstatus.md`](details/commitstatus.md) | `commitstatus` owner の GitHub Commit Status API payload、送信順、失敗時非反転、保存値、secret mask、検証条件。 | runner の build 実行判断、GitHub read、API endpoint、SDK method、UI DOM 詳細、状態 schema、setup / release 手順、fixture 証跡責務。 |
| [`docs/details/security.md`](details/security.md) | `security` owner の API token scope、API key、audit、session timeout、TOTP、rate limit、漏えい禁止、security 横断順序。 | API endpoint 共通処理、SDK method 実装、UI DOM 詳細、runner / builder の業務処理、状態 schema、setup / release 手順、fixture 証跡責務。 |
| [`docs/details/fixture.md`](details/fixture.md) | fixture manifest、assertion、fake、testdata、expected / effects、受け入れ fixture 共通契約、実装検証証跡テンプレート、acceptance checklist、差し戻し条件、実装検証証跡。 | 個別 component の通常処理本文、API endpoint 詳細、SDK method 実装、UI DOM 詳細、状態 schema、setup / release 実行手順。 |

[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6.3 は横断補足契約であり、横断処理順、同期禁止、成功後再取得、失敗時固定、横断受け入れ観点だけを扱う。個別機能本文は各 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を参照する。

すべての owner component 別の [`docs/details/*.md`](details/) 詳細本文責務は、冒頭に `## 0. 責務境界` を置き、以下の 4 項目を同じ意味で持つ。

| 項目 | 必須内容 |
|------|----------|
| owner component | そのファイルが主本文として扱う component を 1 件だけ書く。 |
| collaborator component | 呼び出し元、呼び出し先、schema、表示、検証の参照先を 0 件以上書く。owner component を含めてはならない。 |
| 持つ内容 | そのファイルだけが主本文として定義する内容を書く。 |
| 持たない内容 | 他 owner component へ委ねる内容を書く。 |

責務境界表の `持つ内容` と `持たない内容` が本文と矛盾する場合は、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.2a と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0 の判定を参照する。[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務では、責務境界表、本文、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b の詳細仕様参照表、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務の確認先を示す。

owner component 別の [`docs/details/*.md`](details/) 詳細本文責務の各節について、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務では以下を入口確認項目として扱う。

| 項目 | 入口確認 |
|------|----------|
| 節番号 | 既存の節番号を崩さない。番号の再採番は行わない。 |
| 参照 | 入口ファイルと owner component ファイルの参照先が一意に追跡できるよう、詳細仕様入口責務の対応表を更新する。 |
| owner | 各機能節に owner component を 1 件だけ明記する。 |
| collaborator | collaborator component は 0 件以上を明記し、owner component を含めない。 |
| 重複確認 | 同じ入力、出力、状態 schema、HTTP body、DOM id、fixture assertion の責務が 1 つに定まっている。 |
| 横断事項 | 横断する固定値は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務で確認し、横断共通基盤の扱いは [`docs/SPEC.md`](SPEC.md) 方針責務 §4.2 を参照する。 |
| 索引 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務に、owner component 別詳細本文責務の役割と責務範囲を反映する。 |

詳細仕様本文へ進む入口は、詳細仕様入口責務の「詳細仕様参照入口」に固定する。owner component 別の [`docs/details/*.md`](details/) 詳細本文責務は owner component の主本文として読む。collaborator component は、owner component 本文または詳細仕様入口責務の参照表で明示された呼び出し境界、schema、表示、security、setup、fixture、検証観点に限定して確認する。

責務整理、本文配置変更、参照先更新の整合確認では、以下を入口確認項目として扱う。

| 確認条件 | 確認先 |
|----------|------|
| 各本文が対応する owner component 別の [`docs/details/*.md`](details/) 詳細本文責務にだけ存在している。 | owner component 別の [`docs/details/*.md`](details/) 詳細本文責務 |
| [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務には入口、索引、共通固定値、対応表、owner / collaborator 境界管理、リポジトリ内ソース配置だけが残っている。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 |
| [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務、各 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務間の参照が矛盾していない。 | 該当責務文書 |
| 廃止済み節名、廃止済みファイル名、不要参照の取り残し。 | `rg` による検索証跡 |
| 実装ファイル、fixture、testdata の内容を文書配置整理だけで変更していない。 | `git diff --name-status` |

---

## 0c. 実装前確認項目

実装着手可否の判定は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.7 と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0d を参照する。実装前確認項目は、対象機能について実装前に確認する入口項目を示す。

| ゲート | 入口確認 |
|--------|----------|
| 対応表 | 対象機能が [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i の詳細節対応表に記載され、詳細仕様節と受け入れ条件が一意に示されている。 |
| テンプレート | 対象機能の詳細仕様が [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h の機能仕様テンプレートの必須項目を満たしている。 |
| 責務境界 | owner component、collaborator component、対象ファイル、呼び出し元、呼び出し先、変更対象として許可された状態ファイルが明記されている。 |
| 入出力 | すべての入力、出力、既定値、許容値、必須/任意、型、文字コード、時刻形式が明記されている。 |
| 状態管理 | 状態ファイルのパス、JSON 形式、更新タイミング、初期状態、破損時の扱い、権限が明記されている。 |
| 正常系 | 処理順序、分岐条件、ループ条件、成功条件、終了条件が明記されている。 |
| 異常系 | エラー条件、ログレベル、HTTP ステータス、戻り値、再試行有無、処理継続/中断条件が明記されている。 |
| 冪等性 | 同一リクエスト、再実行、途中失敗後の再開で二重実行・二重削除・状態破壊が発生しない条件が明記されている。 |
| 排他制御 | 同時実行、ロック、タイムアウト、ロック残存時の扱いが明記されている。 |
| セキュリティ | 秘密情報の保存禁止、マスク、ファイル権限、認証/認可、外部公開可否が明記されている。 |
| 検証 | 構文確認、単体確認、手動 API 確認、生成物確認、ログ確認、失敗系確認のいずれを行うかが明記されている。 |

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0c の実装前確認ゲートのいずれかが未充足の場合は、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.7、[`docs/SPEC.md`](SPEC.md) ポリシー責務 §0、[`docs/SPEC.md`](SPEC.md) ポリシー責務 §0d を参照する。[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務では、不足している入口項目と更新すべき対応表を特定する。

実装後の判定は、以下を入口として確認する。実装前確認項目は判定条件本文の正本ではない。Phase、状態分類、引き継ぎ契約は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、fixture、fake、実装検証証跡、acceptance checklist、差し戻し条件は [`docs/details/fixture.md`](details/fixture.md)、setup / release 実行条件は [`docs/details/setup.md`](details/setup.md) を正本として参照する。

1. 実装した機能が、該当 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務と詳細仕様入口責務の対応表に記載された入力、出力、状態、異常系、検証条件と一致する。
2. 対象機能が owner component 別の [`docs/details/*.md`](details/) 詳細本文責務で [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h の機能仕様テンプレートを満たし、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i の詳細節対応表の受け入れ条件を満たしている。
3. [`docs/details/fixture.md`](details/fixture.md) が要求する fixture、fake、expected / effects、実装検証証跡、対象外確認を満たしている。
4. 実装対象外に残す機能が [`docs/details/fixture.md`](details/fixture.md) の実装検証証跡に明記されている。
5. [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務、[`AGENTS.md`](../AGENTS.md) の参照リンクとファイル名が矛盾していない。
6. 実装ファイルを変更した場合、構文確認または実行確認の結果が記録できる。
7. 仕様外の挙動、暗黙の既定値、未記載の状態ファイル、未記載のエラー応答が存在しない。

---

## 0d. 共通固定値

共通固定値は、各コンポーネントで共通して使用する固定値を定義する。個別節に別の値が明記されていない限り、共通固定値を優先する。

| 項目 | 決定 |
|------|------|
| Go 最小バージョン | Go `1.22` 以上。標準ライブラリのみを使用し、外部 module は追加しない。 |
| 文字コード | 入力、出力、状態ファイル、HTTP body は UTF-8 固定。UTF-8 として読み取れない入力は処理を中断する。 |
| 改行 | 新規に書き出す text / JSON Lines ファイルは LF 固定。CRLF 入力は読み込み時に LF として扱う。 |
| 時刻 | 状態ファイル、API、ログの機械処理用時刻は UTC の ISO 8601 形式（例: `2026-09-16T09:00:00Z`）で保存する。ローカル時刻への変換は UI 表示に限定し、状態ファイル、API response、ログには保存しない。 |
| JSON | JSON object の未知キーは保存しない。読み込み時に未知キーを見つけた場合は無視し、次回保存時に除去する。 |
| atomic write | 状態ファイル更新手順は [`docs/details/statefile.md`](details/statefile.md) §22.0a を参照する。各 component は [`docs/details/statefile.md`](details/statefile.md) §22.0a の手順を使用し、独自更新手順を持たない。 |
| 権限 | 状態ファイル、秘密情報ファイル、ディレクトリの権限は [`docs/details/statefile.md`](details/statefile.md) §22.0a を参照する。 |
| ロック | 状態ファイル lock の作成、待機、解除、競合時応答は [`docs/details/statefile.md`](details/statefile.md) §22.0a を参照する。 |
| ログ秘密情報 | PAT、Webhook Secret、SMTP password、session token、API token は stdout、stderr、JSON log、API response、UI 表示へ平文出力しない。表示が必要な場合は `"***"` とする。 |
| 終了コード | CLI / runner は `0` 成功、`1` 一般エラー、`2` 入力・設定エラー、`3` 外部サービス・ネットワークエラー、`4` ロック競合を標準とする。個別節に明記がある場合もこの意味から外してはならない。 |
| 仕様外追加の扱い | 仕様にない環境変数、状態ファイル、HTTP endpoint、CLI option、外部依存の扱いは [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0 と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §4 を参照する。必要な場合は該当 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務または詳細仕様入口責務の対応表を確認する。 |

---

## 0e. 完全実装検証マトリクス

検証マトリクスは、実装確認へ進むための参照入口である。fixture 名、expected / effects、fake 動作、実装検証証跡項目、acceptance checklist、差し戻し条件は [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務を参照する。検証マトリクスは fixture 本文、期待値本文、差し戻し条件本文の正本ではない。

対象項目を確認済み扱いにできるかは、対象 owner component の詳細本文、関連 collaborator component の詳細本文、fixture 証跡責務、状態・計画責務を同時に確認する。確認済み扱いの可否は [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0a〜§0f、状態分類は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、fixture と実装検証証跡は [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務を参照する。

| 対象 | 詳細本文参照 | fixture / 証跡参照 |
|------|--------------|--------------------|
| `builder` | [`docs/details/builder.md`](details/builder.md) §2〜§8、[`docs/details/builder.md`](details/builder.md) §28 | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F、[`docs/details/fixture.md`](details/fixture.md) §28-F |
| `runner` | [`docs/details/runner.md`](details/runner.md) §10〜§20、[`docs/details/runner.md`](details/runner.md) §27.1〜§27.38 | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F、[`docs/details/fixture.md`](details/fixture.md) §27-F |
| `api` | [`docs/details/api.md`](details/api.md) §21〜§22、[`docs/details/api.md`](details/api.md) §25、[`docs/details/api.md`](details/api.md) §27 | [`docs/details/fixture.md`](details/fixture.md) §22-F、[`docs/details/fixture.md`](details/fixture.md) §27-F |
| `sdk` | [`docs/details/sdk.md`](details/sdk.md) §23 | [`docs/details/fixture.md`](details/fixture.md) §22-F、[`docs/details/fixture.md`](details/fixture.md) §27-F |
| `ui` | [`docs/details/ui.md`](details/ui.md) §24 | [`docs/details/fixture.md`](details/fixture.md) §22-F、[`docs/details/fixture.md`](details/fixture.md) §27-F |
| `setup` | [`docs/details/setup.md`](details/setup.md) §26 | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F、[`docs/details/fixture.md`](details/fixture.md) §27-F |
| `statefile` | [`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/statefile.md`](details/statefile.md) §22.0c | [`docs/details/fixture.md`](details/fixture.md) §22-F、[`docs/details/fixture.md`](details/fixture.md) §27-F |
| `security` | [`docs/details/security.md`](details/security.md) §27.42〜§27.47 | [`docs/details/fixture.md`](details/fixture.md) §27-F |
| `archive` | [`docs/details/archive.md`](details/archive.md) §27.7、[`docs/details/archive.md`](details/archive.md) §27.15 | [`docs/details/fixture.md`](details/fixture.md) §27-F |
| `commitstatus` | [`docs/details/commitstatus.md`](details/commitstatus.md) §27.1 | [`docs/details/fixture.md`](details/fixture.md) §27-F |
| `admin` | [`docs/details/admin.md`](details/admin.md) §0、A1〜A6 | [`docs/details/fixture.md`](details/fixture.md) §27-F |
| `fixture` | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F、[`docs/details/fixture.md`](details/fixture.md) §22-F、[`docs/details/fixture.md`](details/fixture.md) §27-F、[`docs/details/fixture.md`](details/fixture.md) §28-F | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F、[`docs/details/fixture.md`](details/fixture.md) §22-F、[`docs/details/fixture.md`](details/fixture.md) §27-F、[`docs/details/fixture.md`](details/fixture.md) §28-F |

状態分類、実装可否、Phase 判定条件は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を参照する。確認済み扱いの可否と仕様変更判定条件は [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0a〜§0f を参照する。

---

## 0f. 仕様策定完了チェック

仕様策定完了チェックは、Go 版初期実装へ進む前に確認する参照入口である。実装着手可否、判定方針は [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0a〜§0f、対象 Phase と状態分類は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §4 を参照する。仕様策定完了チェックは各 component の判定本文、fixture 本文の正本ではない。

| 対象 | 実装前に確認する詳細本文 | 確認する入口 |
|------|--------------------------|--------------|
| `builder` | [`docs/details/builder.md`](details/builder.md) §2〜§8、[`docs/details/builder.md`](details/builder.md) §28 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0e、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 |
| `runner` | [`docs/details/runner.md`](details/runner.md) §10〜§20、[`docs/details/runner.md`](details/runner.md) §27.1〜§27.38 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0e、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2 |
| `api` | [`docs/details/api.md`](details/api.md) §21〜§22、[`docs/details/api.md`](details/api.md) §25、[`docs/details/api.md`](details/api.md) §27 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0e、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.3 |
| `sdk` | [`docs/details/sdk.md`](details/sdk.md) §23 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0e、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.3 |
| `ui` | [`docs/details/ui.md`](details/ui.md) §24 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0e、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.3 |

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0f の仕様策定完了チェック表の対象外である `mcp`、MCP tools、MCP resources、MCP prompts、HTTP SSE transport、MCP audit / stats / config CRUD は、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務で将来計画として確認する。これらは、現時点の [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務および owner component 別の [`docs/details/*.md`](details/) 詳細本文責務に入出力、状態、起動手順、検証条件を持たない。

仕様策定完了チェックで未充足が見つかった場合の実装着手可否は、[`docs/SPEC.md`](SPEC.md) ポリシー責務 §0a〜§0f を参照する。詳細仕様入口責務では、改訂先の確認順を以下に示す。

1. 未充足項目が詳細仕様入口責務の記載対象外である場合は、先に [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務を確認する。
2. 未充足項目が入出力、状態ファイル、api、sdk、ui、処理順序、異常系、検証条件に関わる場合は、該当 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務または collaborator component 別の [`docs/details/*.md`](details/) 詳細本文責務を改訂する。
3. ファイル名、正本参照先、対象範囲が変わる場合は、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務の更新要否を確認する。
4. 対象項目の詳細節、参照先、受け入れ条件が変わる場合は、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i の詳細節対応表を更新する。
5. 改訂後、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b〜§0j、および [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務の Phase 条件を再確認する。

---


## 0g. Phase 詳細仕様参照

Phase の一覧、順序、対象 owner component、依存条件、判定条件、引き継ぎ契約、実装成果物チェックリストは [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を参照する。

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0g は、Phase を詳細仕様へ接続する入口だけを持つ。Phase の状態、順序、判定条件、将来計画の昇格判定は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を正本とする。

| Phase | owner component | 主な詳細仕様本文 | fixture / 証跡 |
|-------|-----------------|------------------|----------------|
| Phase 1 | `builder` | [`docs/details/builder.md`](details/builder.md) §1〜§9、[`docs/details/builder.md`](details/builder.md) §8a | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F |
| Phase 2 | `runner` | [`docs/details/runner.md`](details/runner.md) §10〜§20、[`docs/details/runner.md`](details/runner.md) §15a | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F |
| Phase 3 | `api` | [`docs/details/api.md`](details/api.md) §21〜§22、[`docs/details/api.md`](details/api.md) §21a、[`docs/details/api.md`](details/api.md) §25、[`docs/details/api.md`](details/api.md) §22.0f、[`docs/details/setup.md`](details/setup.md) §26 | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F、[`docs/details/fixture.md`](details/fixture.md) §22-F、[`docs/details/fixture.md`](details/fixture.md) §27-F |
| Phase 4 | `api` | [`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/api.md`](details/api.md) §22.0f | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F、[`docs/details/fixture.md`](details/fixture.md) §22-F、[`docs/details/fixture.md`](details/fixture.md) §27-F |
| Phase 5 | `sdk` | [`docs/details/sdk.md`](details/sdk.md) §23 | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F |
| Phase 6 | `ui` | [`docs/details/ui.md`](details/ui.md) §24 | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F |

対象 Phase の状態は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務で確認し、詳細本文は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0g の Phase 詳細仕様参照表に記載された owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を読む。

## 0h. 機能仕様テンプレート

対象項目を追加または改訂する場合は、該当する owner component 別の [`docs/details/*.md`](details/) 詳細本文責務の節に [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h の機能仕様テンプレート表の項目をすべて含める。既存節に含める場合も、実装者が [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h の機能仕様テンプレート表の項目を本文から一意に読み取れる状態にする。

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h の機能仕様テンプレート表の `入口判定` は、詳細仕様入口責務として不足項目を検出するための要約である。実装着手可否、確認済み扱い、禁止事項の判定は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.2a および [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0a〜§0f を参照する。

| 項目 | 必須内容 | 入口判定 |
|------|----------|----------------|
| 目的 | 何を解決する機能か、どの利用者または運用者のための機能か。 | 不足時は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.7 と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0d を参照する。 |
| owner / collaborator component | owner component と collaborator component を明記し、複数 component が関わる場合は責務境界を分けて書く。 | 不足時は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.2a と [`docs/SPEC.md`](SPEC.md) 方針責務 §4.7 を参照する。 |
| 入力 | CLI 引数、HTTP request、設定値、状態ファイル、環境変数、Markdown 入力、UI 操作などの入力元、型、必須/任意、既定値。 | 不足時は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.7 と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0d を参照する。 |
| 出力 | 生成ファイル、HTTP response、stdout/stderr、ログ、通知、UI 表示、終了コード。 | 不足時は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.7 と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0d を参照する。 |
| 状態 | 読み書きする状態ファイル、ディレクトリ、メモリ状態、ロック、更新責務、初期値、破損時の扱い。 | 不足時は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.7、[`docs/SPEC.md`](SPEC.md) ポリシー責務 §0d、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を参照する。 |
| 正常系 | 処理順序、分岐条件、成功条件、保存順序、外部コマンド呼び出し条件。 | 不足時は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.7 と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0d を参照する。 |
| 異常系 | エラー条件、継続/中断、HTTP status、終了コード、ログレベル、通知、リトライ有無。 | 不足時は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.7 と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0d を参照する。 |
| セキュリティ | 秘密情報、認証、認可、ファイル権限、外部公開可否、ログ出力制約。 | セキュリティ影響の扱いは [`docs/SPEC.md`](SPEC.md) ポリシー責務 §5、[`docs/SPEC.md`](SPEC.md) ポリシー責務 §11、[`docs/SPEC.md`](SPEC.md) ポリシー責務 §12 を参照する。 |
| 検証 | 必須テスト、手動確認、fixture、生成物確認、API 確認、異常系確認。 | 不足時は [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0a と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務を参照する。 |
| 判定条件 | どの検証が成功したら確認済みと扱うか。関連文書の更新要否。 | 不足時は [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0a、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、[`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務を参照する。 |

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h の機能仕様テンプレート表のいずれかが不足する対象項目の扱いは、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.7 と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0d を参照する。[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務では、不足している入口項目と参照先変更の有無だけを確認する。

---

## 0i. 詳細節対応表

詳細節対応表は、対象機能から該当する詳細仕様へ移動するための対応表である。対象機能の詳細仕様節と受け入れ条件の入口は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1〜§0i.4 の各対応表で確認する。

表の「参照 component」は参照先を探すための component 一覧である。owner component と collaborator component は、対象機能の詳細仕様節に記載された値を参照する。

表の「詳細仕様節」が複数ある場合は、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を主本文として読み、collaborator component 別の [`docs/details/*.md`](details/) 詳細本文責務は schema、呼び出し境界、表示、security、setup、fixture、検証観点の確認として読む。詳細仕様節には、参照先ファイル名と節番号を併記する。裸の節番号だけで参照先を解決してはならない。

詳細節対応表は owner component を置き換える表ではない。受け入れ条件が複数 component にまたがる場合でも、主本文は owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を基準とし、collaborator component 別の [`docs/details/*.md`](details/) 詳細本文責務は schema、呼び出し境界、表示、security、setup、fixture、検証観点の確認に限定する。owner component の入力、出力、状態、endpoint、SDK method、UI 操作の本文は owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を正本とする。

対象節に [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h の必須項目が不足している場合の実装着手可否は、[`docs/SPEC.md`](SPEC.md) ポリシー責務 §0a〜§0f を参照する。[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6 の追加仕様化機能参照で、owner、主本文、collaborator を確認する。

### 0i.1 Builder / 静的 Web サイト出力

[`docs/details/builder.md`](details/builder.md) §28.1〜§28.25 を実装する場合は、個別行の受け入れ条件に加えて、[`docs/details/builder.md`](details/builder.md) §28 の共通固定契約、CLI / 設定 / REPORT / 出力識別子固定契約、実装パイプライン固定契約、設定解決・終了コード固定契約、設定ファイル / 入力解決固定契約、REPORT 値型固定契約、stdout / stderr / REPORT 固定契約、atomic write / manifest / search index 副作用固定契約、ID / slug / search index / JS state 決定性固定契約、CSS / JS 出力固定契約、CSS / layout / print / visual 固定契約、browser runtime 固定契約、Markdown token / HTML node 変換固定契約、Markdown parser 優先順位固定契約、Markdown 構文文法固定契約、曖昧構文・機能併用固定契約、warning / error code 固定契約、既存出力互換・先取り実装禁止固定契約、個別固定補足契約、詳細実装確認ゲート固定契約、および [`docs/details/fixture.md`](details/fixture.md) §28-F を必ず読む。

| 機能 | 参照 component | 詳細仕様節 | 受け入れ条件 |
|------|-------------------|------------|--------------|
| 出力サイトサイズ警告閾値 | `builder` / `runner` / `api` | [`docs/details/builder.md`](details/builder.md) §8、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/api.md`](details/api.md) §22.0e | `OUTPUT_SIZE_WARN_MB`、`size_warn`、WARN ログ、API 表示が一致する。 |
| 出力サイトへのビルドメタ埋め込み | `builder` / `runner` / `api` | [`docs/details/builder.md`](details/builder.md) §2、[`docs/details/builder.md`](details/builder.md) §5、[`docs/details/builder.md`](details/builder.md) §8、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/builder.md`](details/builder.md) §27.4 | CLI/env 入力、HTML meta、REPORT、build log、`GET /api/output-meta` の値が一致する。 |
| 変換レポート出力 | `builder` / `runner` / `api` | [`docs/details/builder.md`](details/builder.md) §8、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/api.md`](details/api.md) §22.0e | `[REPORT]` stdout、runner 取り込み、`.build_logs`、`GET /api/output-meta` が一致する。 |
| シンタックスハイライト | `builder` | [`docs/details/builder.md`](details/builder.md) §7.8 | 対応言語、class 名、HTML escape、CSS 表示が一致する。 |
| 本文内全文検索 | `builder` | [`docs/details/builder.md`](details/builder.md) §7.9 | `assets/search-index.json`、検索 UI、ヒット遷移、対象テキストが一致する。 |
| アンカーリンク自動検証 | `builder` | [`docs/details/builder.md`](details/builder.md) §4.3、[`docs/details/builder.md`](details/builder.md) §8 | broken anchor 検出、`[WARN] BROKEN_LINK`、report 件数が一致する。 |
| コードブロックの折りたたみ | `builder` | [`docs/details/builder.md`](details/builder.md) §7.10 | 30 行超の初期折りたたみ、展開操作、印刷時展開が一致する。 |
| 印刷スタイル（`@media print`） | `builder` | [`docs/details/builder.md`](details/builder.md) §6 | `@media print` の非表示対象、コード展開、リンク URL 表示が一致する。 |
| 静的 Web サイト出力 | `builder` | [`docs/details/builder.md`](details/builder.md) §2、[`docs/details/builder.md`](details/builder.md) §5、[`docs/details/builder.md`](details/builder.md) §6、[`docs/details/builder.md`](details/builder.md) §7 | 入力ファイル/ディレクトリ、出力ファイル構成、asset、ページ生成が一致する。 |
| テーマコンポーネント | `builder` | [`docs/details/builder.md`](details/builder.md) §5、[`docs/details/builder.md`](details/builder.md) §6、[`docs/details/builder.md`](details/builder.md) §7 | `adlaire-default` の component、class、slot、asset 出力が一致する。 |
| 外部リンクの自動処理 | `builder` | [`docs/details/builder.md`](details/builder.md) §4.3 | `target="_blank"`、`rel="noopener noreferrer"`、内部リンクとの区別が一致する。 |
| 読み取り進捗バー | `builder` | [`docs/details/builder.md`](details/builder.md) §7.13 | 3px 固定表示、scroll 連動、初期/末尾状態が一致する。 |
| コードブロックのコピーボタン | `builder` | [`docs/details/builder.md`](details/builder.md) §7.6 | ボタン配置、コピー対象、成功/失敗時表示、アクセシビリティが一致する。 |
| 見出しアンカーリンクコピー | `builder` | [`docs/details/builder.md`](details/builder.md) §3、[`docs/details/builder.md`](details/builder.md) §7.11 | `.hn-link`、copy URL、重複 slug 連動が一致する。 |
| TOC 開閉状態の永続化 | `builder` | [`docs/details/builder.md`](details/builder.md) §7.3 | `localStorage` key、展開/折りたたみ、復元条件が一致する。 |
| 見出しスラグ重複解決 | `builder` | [`docs/details/builder.md`](details/builder.md) §4.5 | `-2`、`-3` の付与、TOC、検索、コピー URL との共通化が一致する。 |
| 前後章ナビゲーションボタン | `builder` | [`docs/details/builder.md`](details/builder.md) §4.5、[`docs/details/builder.md`](details/builder.md) §5、[`docs/details/builder.md`](details/builder.md) §7.15 | h2 単位の前後判定、章末尾配置、端の非表示条件が一致する。 |
| 内部リンク整合性チェック | `builder` | [`docs/details/builder.md`](details/builder.md) §4.3、[`docs/details/builder.md`](details/builder.md) §8 | `[label](#anchor)` 検証、WARN、`broken_links` が一致する。 |
| 見出し階層スキップ警告 | `builder` | [`docs/details/builder.md`](details/builder.md) §4.5、[`docs/details/builder.md`](details/builder.md) §8 | h1→h3 等の検出、WARN、`heading_skips` が一致する。 |
| 読了時間推計と表示 | `builder` | [`docs/details/builder.md`](details/builder.md) §4.5、[`docs/details/builder.md`](details/builder.md) §5、[`docs/details/builder.md`](details/builder.md) §6、[`docs/details/builder.md`](details/builder.md) §8 | 対象文字数、200文字/分、切り上げ、header 表示、report が一致する。 |
| テーブルのソート機能 | `builder` | [`docs/details/builder.md`](details/builder.md) §7.14 | クリック操作、昇順/降順、`aria-sort`、インジケーターが一致する。 |
| キーボードショートカット | `builder` | [`docs/details/builder.md`](details/builder.md) §7.12 | `/`、`Escape`、`t` の対象、フォーカス条件、入力中の無効化が一致する。 |
| ビルドキャッシュ | `builder` / `runner` | [`docs/details/builder.md`](details/builder.md) §5、[`docs/details/builder.md`](details/builder.md) §8、[`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/builder.md`](details/builder.md) §27.25 | `.build_cache.json`、入力 manifest、再利用条件、無効化条件、report counters が一致する。 |
| 依存ファイルトラッキング | `builder` / `runner` | [`docs/details/builder.md`](details/builder.md) §4.3、[`docs/details/builder.md`](details/builder.md) §5、[`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/builder.md`](details/builder.md) §27.28 | `.dependency_manifest.json`、依存抽出、関連 target 判定、破損時 full build が一致する。 |
| 差分ビルド | `builder` / `runner` | [`docs/details/builder.md`](details/builder.md) §28.1 | changed manifest、dependency 逆引き、未変更 page 維持、search index 再生成が一致する。 |
| 複数出力形式 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.2 | `html` のみ実出力し、予約 format は実行前に終了コード `2` で拒否する。 |
| Markdown 拡張記法サポート | `builder` | [`docs/details/builder.md`](details/builder.md) §28.3 | admonition、badge、escape、report count が一致する。 |
| コードブロック行番号表示 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.4 | 行番号、copy 対象除外、fold / highlight 併用が一致する。 |
| 見出しの自動採番 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.5 | 表示番号、TOC 番号、slug 不変、search index が一致する。 |
| セクション折りたたみ | `builder` | [`docs/details/builder.md`](details/builder.md) §28.6 | section 範囲、toggle、localStorage、印刷時展開が一致する。 |
| TOC 深さ制御 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.7 | heading filter、active tracking、範囲 validation が一致する。 |
| 最終更新日の自動埋め込み | `builder` | [`docs/details/builder.md`](details/builder.md) §28.8 | git/file timestamp、UTC 秒精度、footer、report が一致する。 |
| diff ハイライト | `builder` | [`docs/details/builder.md`](details/builder.md) §28.9 | diff fence の inserted/deleted/header/context class と copy 本文が一致する。 |
| 画像の遅延読み込み | `builder` | [`docs/details/builder.md`](details/builder.md) §28.10 | `loading="lazy"`、`decoding="async"`、path warning、alt escape が一致する。 |
| カスタムメタタグ注入 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.11 | meta validation、name/property、禁止 key、escape が一致する。 |
| ライトモード固定 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.12 | 生成 HTML が light 固定で、dark / auto / theme toggle / color scheme 永続化を出力しない。 |
| コードブロックのファイル名表示 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.13 | fence info parsing、title escape、copy 対象除外が一致する。 |
| テンプレート変数展開 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.14 | 変数 validation、code fence 内非置換、未定義時 warning/strict が一致する。 |
| HTML ミニファイ | `builder` | [`docs/details/builder.md`](details/builder.md) §28.15 | safe minify、pre/code 保持、必須 marker 検証、report が一致する。 |
| TOC ハイライト追従 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.16 | `.is-active`、`aria-current`、fallback、TOC depth 連動が一致する。 |
| Mermaid ダイアグラム描画 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.17 | 内製対応範囲、unsupported warning、外部 script 不使用が一致する。 |
| 脚注サポート | `builder` | [`docs/details/builder.md`](details/builder.md) §28.18 | 参照番号、footnotes、backlink、未定義時 warning/strict が一致する。 |
| インライン数式レンダリング | `builder` | [`docs/details/builder.md`](details/builder.md) §28.19 | inline/block math、escape、未閉鎖 delimiter、code 内非変換が一致する。 |
| ページ内ナビゲーション履歴 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.20 | hash pushState、back/forward、focus 移動、JS 無効時 fallback が一致する。 |
| 読み上げ対応（アクセシビリティ） | `builder` | [`docs/details/builder.md`](details/builder.md) §28.21 | landmark、aria-label、skip link、focus 順、重複 id 検出が一致する。 |
| 画像ライトボックス | `builder` | [`docs/details/builder.md`](details/builder.md) §28.22 | dialog、Escape、backdrop close、focus trap、alt warning が一致する。 |
| 印刷時 QR コード挿入 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.23 | URL validation、print-only SVG、長さ制限、外部 library 不使用が一致する。 |
| 定義リストサポート | `builder` | [`docs/details/builder.md`](details/builder.md) §28.24 | `<dl>/<dt>/<dd>` 出力、paragraph 境界、inline escape が一致する。 |
| タスクリストサポート | `builder` | [`docs/details/builder.md`](details/builder.md) §28.25 | disabled checkbox、checked 判定、nested list、aria が一致する。 |

### 0i.2 Runner / CI 実行

| 機能 | 参照 component | 詳細仕様節 | 受け入れ条件 |
|------|-------------------|------------|--------------|
| ビルドタイムアウト | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/api.md`](details/api.md) §22.0e | `build_timeout_seconds` の既定値、設定 API、`context.WithTimeout` の中断処理、終了コード、ログが一致する。 |
| ビルドログのファイル保存 | `runner` | [`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15 | `.build_logs/{id}.json` の schema、stdout/stderr、変換レポート、duration、権限が一致する。 |
| ネットワーク断時の再試行 | `runner` | [`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13 | `API_RETRY_MAX`、`API_RETRY_BASE_SECONDS`、指数バックオフ、失敗時ログが一致する。 |
| GitHub API レート制限自動待機 | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §13、[`docs/details/api.md`](details/api.md) §22.0e | `X-RateLimit-Remaining` と `X-RateLimit-Reset` の扱い、待機、API 表示が一致する。 |
| 転送後リモート整合性検証 | `runner` | [`docs/details/runner.md`](details/runner.md) §14a、[`docs/details/runner.md`](details/runner.md) §13 | SSH 転送後の SHA256 照合、不一致時の `.pending_transfers` 再投入、ログが一致する。 |
| マルチブランチビルド | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/api.md`](details/api.md) §22.0e | `BRANCH_TARGETS` と `.branch_config` の優先順位、順次処理、API 更新が一致する。 |
| ビルドログ世代管理 | `runner` | [`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15 | `LOG_KEEP_N` 超過時の削除順序、0 の扱い、削除ログが一致する。 |
| ビルド出力の外部転送 | `runner` | [`docs/details/runner.md`](details/runner.md) §14a、[`docs/details/runner.md`](details/runner.md) §13 | SSH 差分転送、複数ファイル処理、失敗時 pending、通知が一致する。 |
| ビルドクールダウン | `runner` | [`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13 | `BUILD_COOLDOWN_SECONDS` 内の起動スキップ、Webhook 二重トリガー抑止、ログが一致する。 |
| ビルド前の事前チェック | `runner` | [`docs/details/runner.md`](details/runner.md) §13、[`docs/details/setup.md`](details/setup.md) §26 | ディスク、`adlaire-ci-build`、pipeline 前提の確認、不足時の ERROR と通知が一致する。 |
| 定期強制ビルド | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/api.md`](details/api.md) §22.0e | `FORCE_BUILD_INTERVAL`、変更なし時の強制ビルド、設定 API が一致する。 |
| ビルド中重複スキップ | `runner` | [`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §13 | `.build_lock` の PID 判定、stale lock、競合時終了コードとログが一致する。 |
| GitHub PAT 有効期限の事前警告 | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §13、[`docs/details/api.md`](details/api.md) §22.0e | `GitHub-Authentication-Token-Expiration` の解析、7 日以内 WARN、API 表示が一致する。 |
| コミット情報のビルドログ記録 | `runner` | [`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15 | SHA、message、author、date を build id と同じログへ記録する。 |
| GitHub API 連続失敗によるサーキットブレーカー | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/api.md`](details/api.md) §22.0e | 閾値、open/close 状態、API reset、通知、状態ファイルが一致する。 |
| 設定ファイル起動時整合性チェック | `runner` | [`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/runner.md`](details/runner.md) §27.10 | 対象 JSON ファイル、検証順序、破損退避、初期化値、ログ、通知、終了コード、fixture が一致する。 |
| ビルドステータスファイル出力 | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.8 | `.build_status.json` の schema、更新タイミング、status/target_status、pending 件数、circuit 状態、API 参照元が一致する。 |
| ビルドトリガー種別の記録 | `runner` / `api` / `sdk` / `ui` | [`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/runner.md`](details/runner.md) §27.9 | `trigger` の有効値、判定条件、`.build_logs`、`.build_history`、`.build_status.json`、履歴 filter、UI 表示が一致する。 |
| GitHub Commit Status API | `runner` | [`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/commitstatus.md`](details/commitstatus.md) §27.1 | `commit_status_enabled`、context、target_url、pending/success/failure の送信条件、失敗時の扱い、build log 記録が一致する。 |
| ドライラン実行モード | `runner` | [`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/runner.md`](details/runner.md) §27.2 | `--dry-run` が状態ファイル、log、history、deploy、通知を変更せず、設定・GitHub・SHA 判定結果を固定 JSON で返す。 |
| ビルド失敗時の自動リトライ | `runner` | [`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/runner.md`](details/runner.md) §27.3 | retry 対象エラー、最大回数、backoff、attempt log、最終 status、SHA 更新禁止条件が一致する。 |
| Webhook 通知失敗リトライキュー | `runner` | [`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §16 | `.notify_pending` の schema、再送順序、失敗時保持が一致する。 |
| 複数ファイル監視 | `runner` / `builder` / `api` | [`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/runner.md`](details/runner.md) §27.21 | `target_files` の検証、対象別 SHA 差分、build target 決定、履歴・ログ・API 表示が一致する。 |
| ビルドパイプライン YAML 定義 | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.22 | `.pipeline.yml` の内製 subset parse、step 実行順、timeout、env、失敗時 status、API 保存が一致する。 |
| ローカルファイル監視モード | `runner` | [`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §27.23 | GitHub API を呼ばず、local snapshot の SHA-256 差分だけで変更検出し、trigger と status が一致する。 |
| タグ付きコミットのみビルド | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.24 | tag pattern、GitHub tags API、skip 条件、build log/history、設定 API が一致する。 |
| 並列マルチターゲットビルド | `runner` | [`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §14a、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/runner.md`](details/runner.md) §27.26 | worker 上限、target 別 status、pending transfer、最終 build status、ログ順序が一致する。 |
| ビルド前後フック | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.27 | `.hooks` schema、pre/post 実行、abort 条件、hook log、API CRUD が一致する。 |
| リモートビルド対応 | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §14a、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/runner.md`](details/runner.md) §27.29 | remote command、archive 取得、manifest 検証、状態記録、失敗時 rollback 不実行が一致する。 |
| ブランチ別環境変数 | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.31 | branch env schema、許可 key、secret mask、process env 注入、API 保存が一致する。 |
| ビルド通知連携 | `runner` / `api` / `sdk` / `ui` | [`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/runner.md`](details/runner.md) §16、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.32 | 通知 event、channel schema、送信順、retry、mask、notify log、API / UI 表示が一致する。 |
| ビルド時間トレンド記録 | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.33 | `.build_trends.json`、移動平均、中央値、p95、API response、破損時復旧が一致する。 |
| ビルド依存チェーン | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.34 | `.build_chain_config`、依存 DAG 検証、実行順、skip / failure status、chain log が一致する。 |
| ビルド優先度キュー | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.35 | queue priority、created_seq、同一優先度 FIFO、API 表示、cancel / clear が一致する。 |
| 失敗原因の自動分類 | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.36 | failure_category、evidence、分類優先順位、history / log / UI 表示が一致する。 |
| ビルド実行環境の記録 | `runner` | [`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/runner.md`](details/runner.md) §27.37 | build 開始時の environment snapshot、secret 非含有、log schema、検証 fixture が一致する。 |
| ビルド所要時間の異常検知 | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/runner.md`](details/runner.md) §16、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.38 | trend 基準、異常判定、WARN、history flag、通知 payload、設定値が一致する。 |

### 0i.3 API / SDK / UI

| 機能 | 参照 component | 詳細仕様節 | 受け入れ条件 |
|------|-------------------|------------|--------------|
| ポーリング間隔の動的変更 | `api` | [`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/setup.md`](details/setup.md) §26、[`docs/details/api.md`](details/api.md) §27.11 | `POST /api/schedule/interval` が systemd timer 設定を更新し、検証コマンドで反映を確認できる。 |
| GitHub Webhook 受信 | `api` / `runner` | [`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/api.md`](details/api.md) §22-W、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/api.md`](details/api.md) §27.12 | HMAC 検証、イベント記録、キュー投入またはビルドトリガー、エラー応答が一致する。 |
| 設定バリデーション API | `api` / `sdk` / `ui` | [`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/api.md`](details/api.md) §27.5 | `POST /api/config/validate` が状態を変更せず、正規化後設定、warnings、errors を返す。 |
| API アクセスログ | `api` / `sdk` / `ui` | [`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/api.md`](details/api.md) §27.6 | `.api_access_log` の schema、追記対象、マスク条件、一覧 API、UI 表示が一致する。 |
| Webhook イベントログ | `api` | [`docs/details/runner.md`](details/runner.md) §11、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/api.md`](details/api.md) §22-W、[`docs/details/api.md`](details/api.md) §27.13 | `.webhook_events.json` の JSON Lines schema と一覧 API が一致する。 |
| ヘルスチェックエンドポイント | `api` | [`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/api.md`](details/api.md) §27.16 | `GET /api/health` の稼働秒数、最終ビルド、最終転送、エラー応答が一致する。 |
| Webhook イベント一覧取得 API | `api` / `sdk` / `ui` | [`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/api.md`](details/api.md) §27.13 | `GET /api/webhook-events` の query、response、SDK method、UI 表示が一致する。 |
| ビルドログ重大度フィルター | `api` / `sdk` / `ui` | [`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/api.md`](details/api.md) §27.17 | `level=warn\|error` の query、検索結果、UI filter が一致する。 |
| ブランチ設定の動的変更 API | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/api.md`](details/api.md) §27.18 | `.branch_config`、GET/POST API、runner 再起動不要条件が一致する。 |
| 週次ビルドサマリー Webhook | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §16、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.19 | 週次判定、集計対象、通知 payload、手動送信 API が一致する。 |
| 設定変更の詳細 diff 記録 | `api` | [`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/api.md`](details/api.md) §27.20 | `.config_log` の diff 文字列、対象 API、マスク条件が一致する。 |
| ビルド承認フロー | `runner` / `api` / `sdk` / `ui` | [`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/runner.md`](details/runner.md) §16、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/api.md`](details/api.md) §27.30 | `.approval_queue`、承認/却下 API、通知、timeout、UI 操作、履歴 status が一致する。 |

### 0i.4 Archive / Artifact / Security

| 機能 | 参照 component | 詳細仕様節 | 受け入れ条件 |
|------|-------------------|------------|--------------|
| ビルドログのアーカイブ圧縮 | `runner` / `api` / `sdk` / `ui` | [`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/archive.md`](details/archive.md) §27.7 | gzip 形式、archive 先、参照順、cleanup/archive API、disk usage 集計、UI 表示が一致する。 |
| ビルド所要時間の記録と統計 API | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §15、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.14 | `started_at`、`finished_at`、`duration_seconds` と統計 API が一致する。 |
| ビルドアーティファクト世代管理 | `runner` / `api` | [`docs/details/runner.md`](details/runner.md) §14b、[`docs/details/api.md`](details/api.md) §22.0e | `.snapshots/` の保持世代、削除、rollback API が一致する。 |
| ビルドアーティファクト管理 | `api` / `ui` / `sdk` | [`docs/details/runner.md`](details/runner.md) §14b、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/archive.md`](details/archive.md) §27.15 | 一覧、ダウンロード、削除、ロールバックの API、SDK、UI が一致する。 |
| ビルドトリガー専用 API スコープ | `api` / `sdk` / `ui` | [`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/security.md`](details/security.md) §27.42 | `trigger` scope token が build 起動系だけを許可し、その他 API を拒否する。 |
| API キー管理 | `api` / `sdk` / `ui` | [`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/security.md`](details/security.md) §27.43 | API key 本体の一回表示、hash 保存、scope、期限、失効、監査が一致する。 |
| 監査ログ | `api` / `sdk` / `ui` | [`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/security.md`](details/security.md) §27.44 | `.audit_log` schema、対象操作、mask、検索 API、UI 表示が一致する。 |
| セッションタイムアウト変更設定 | `api` / `sdk` / `ui` | [`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/api.md`](details/api.md) §25、[`docs/details/security.md`](details/security.md) §27.45 | `session_timeout_seconds` の範囲、保存、既存 session の扱い、新規 session 期限が一致する。 |
| TOTP 二要素認証 | `api` / `sdk` / `ui` | [`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/api.md`](details/api.md) §25、[`docs/details/security.md`](details/security.md) §27.46 | RFC 6238 TOTP、二段階 login、secret 保存、確認、無効化、UI 操作が一致する。 |
| API レート制限 | `api` / `sdk` / `ui` | [`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/security.md`](details/security.md) §27.47 | IP / actor / endpoint group の固定窓制限、`429`、状態保存、設定 API、UI 表示が一致する。 |

---

## 詳細仕様セット構成

Adlaire CI は Go 版コンポーネントと JavaScript / HTML 管理ツールで構成する。

詳細仕様本文は、`builder`、`runner`、`api`、`admin`、`sdk`、`ui`、`setup`、`statefile`、`archive`、`commitstatus`、`security`、`fixture` の実装詳細を owner component 別に定義する。

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務は詳細仕様参照入口、索引、共通固定値、対応表、リポジトリ内ソース配置だけを持つ。各 component の入出力、状態、処理順序、異常系、検証条件の本文は、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を参照する。

`mcp` は将来計画であり、MCP 専用詳細仕様が新設されるまで、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務および owner component 別の [`docs/details/*.md`](details/) 詳細本文責務に入出力、状態、起動手順、検証条件の本文を持たない。

---

## 0j. リポジトリ内ソース配置

Adlaire CI の標準ディレクトリ構成 tree は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.3 を正本とする。[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0j は、同 tree の配置ごとの責務、標準配置の成立条件、標準外配置の禁止、将来追加予定 path、未作成 path の作成可否を定義する。

| パス | component | 役割 |
|------|-----------|------|
| `main.go` | `-` | 起動入口。現時点では実行ファイル名判定、引数受け取り、`builder` / `runner` 呼び出しを行う。 |
| `components/builder.go` | `builder` | Markdown / Markdown ディレクトリを静的 Web サイトへ変換する。 |
| `components/runner.go` | `runner` | GitHub polling、変更検出、ビルド起動、履歴、ログ、deploy を実行する。 |
| `components/api.go` | `api` | 管理 API サーバー、認証、状態ファイル操作を提供する。 |
| `components/*_test.go` | 対応 owner component | Go 実装ファイルに対応する fixture、hardening、endpoint、状態ファイル、判定テストを配置する。 |
| `admin/index.html` | `ui` | 標準管理ツール UI を提供する。 |
| `admin/adlaire-ci-sdk.js` | `sdk` | 管理 API 通信用 SDK を提供する。 |
| `admin/` | `ui` / `sdk` | 標準管理 UI の静的ファイルを配置する。 |
| `testdata/` | `fixture` | コンポーネント別 fixture を配置する。 |
| `docs/examples/` | `-` | 利用例、設定例、サンプル構成を配置する。 |

`main.go` は 1 ファイルとし、実装詳細を含めない。`components/` 配下は 1 コンポーネント = 1 Go ファイルとし、各ファイルは [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0j の標準リポジトリ内ソース配置表に記載された責務を実装する。

標準配置の成立条件は以下に固定する。

| 対象 | 成立条件 |
|------|--------------|
| `main.go` | repository root に 1 ファイルだけ存在し、実行ファイル名判定、引数受け取り、現行 owner component 呼び出しだけを持つ。Markdown 変換、CI 実行、HTTP handler、状態ファイル操作、archive、commitstatus、MCP の実装詳細を含まない。 |
| `components/*.go` | 実装対象 owner component ごとに 1 Go ファイルだけ存在する。現行 Go 実装は `builder` の `components/builder.go`、`runner` の `components/runner.go`、`api` の `components/api.go` とする。`admin` は `admin/` 配下の静的配布物であり、Go ファイルを持たない。`statefile`、`archive`、`commitstatus` は詳細仕様上の責務境界であり、単独 Go ファイルとして追加する場合は、追加対象 Phase または追加実装変更で仕様状態と索引を更新してから作成する。 |
| `components/*_test.go` | 対応する `components/*.go` の検証ファイルとして扱う。現行 Go 検証ファイルは `components/builder_test.go`、`components/runner_test.go`、`components/api_test.go` とする。新規 owner component の Go 実装を追加する場合は、対応する test file の要否を該当 owner component 別詳細本文責務と fixture 証跡責務で確定する。 |
| 標準外配置禁止 | `build_spec.go`、`runner.go`、`build_spec_test.go`、`runner_test.go`、`testdata/build_spec/` を標準配置として扱わない。同じ実装本文または同じ fixture を標準配置外に併存させない。 |
| testdata | 実装済みまたは仕様化済み・未実装の owner component ごとに `testdata/<component>/` を使用する。`testdata/mcp/` は MCP 専用詳細仕様が新設されるまで作成しない。 |
| admin | `admin/index.html` と `admin/adlaire-ci-sdk.js` は、それぞれ [`docs/details/ui.md`](details/ui.md) と [`docs/details/sdk.md`](details/sdk.md) の owner component 別詳細本文責務に従う。`admin/style.css` と `admin/app.js` は、[`docs/details/admin.md`](details/admin.md) A1 に定義された任意配布物として扱い、未定義の admin 静的ファイルを追加しない。 |
| 将来追加予定 path | `components/mcp.go` と MCP 用 fixture は未作成の将来追加予定 path であり、MCP 専用詳細仕様が新設され、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務で `仕様化済み・未実装` へ昇格するまで作成しない。 |

標準配置を変更する場合は、標準外配置と標準配置の両方が同じ実装実体として併存していないこと、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務の「実装ファイル一覧」、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務の状態分類、該当 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務、testdata 参照が同じ配置を指すことを確認する。

---

## owner component 別詳細本文責務索引

owner component 別の [`docs/details/*.md`](details/) 詳細本文責務の本文は、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務の「詳細仕様管理」と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i を入口として確認する。

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務は詳細仕様本文の正本ではない。対象機能の owner component は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b または [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i で特定し、該当する [`docs/details/*.md`](details/) 詳細本文責務を主本文とする。

| 対象範囲 | 主な参照先 |
|----------|------------|
| [`docs/details/builder.md`](details/builder.md) §1〜§9、[`docs/details/builder.md`](details/builder.md) §8a | [`docs/details/builder.md`](details/builder.md) |
| [`docs/details/runner.md`](details/runner.md) §10〜§20、[`docs/details/runner.md`](details/runner.md) §15a | [`docs/details/runner.md`](details/runner.md) |
| [`docs/details/api.md`](details/api.md) §21〜§22、[`docs/details/api.md`](details/api.md) §21a、[`docs/details/api.md`](details/api.md) §25 | [`docs/details/api.md`](details/api.md) |
| [`docs/details/sdk.md`](details/sdk.md) §23 | [`docs/details/sdk.md`](details/sdk.md) |
| [`docs/details/ui.md`](details/ui.md) §24 | [`docs/details/ui.md`](details/ui.md) |
| [`docs/details/setup.md`](details/setup.md) §26 | [`docs/details/setup.md`](details/setup.md) |
| [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6 追加仕様化機能 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6 と該当 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務 |
