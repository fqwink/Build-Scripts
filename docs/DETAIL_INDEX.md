# Adlaire CI — 詳細仕様

本ファイルは `docs/SPEC.md` の Part 3 詳細仕様の入口であり、索引、共通固定値、実装前確認項目、検証マトリクス、Phase、詳細節対応表、リポジトリ内ソース配置、横断補足契約を持つ正本である。

各 owner component の具体的な入出力、状態、処理順序、異常系、セキュリティ制約、検証条件の本文は、責務 component 別の `docs/details/*.md` を正本とする。

本ファイルに、方針、ポリシー、実装状態、正本関係、ロードマップ状態、実装可否の上位判断を記載してはならない。これらは `docs/SPEC.md` を正とする。

本ファイルは詳細仕様本文を集約する場所ではない。個別 component の処理本文、endpoint 詳細、SDK method、UI DOM、状態 schema、fixture assertion、setup 手順を追記する場合は、該当する owner component 別詳細仕様ファイルを更新し、本ファイルは索引または参照先だけを更新する。

---

# Part 3 — 仕様
> Part 3 詳細仕様セットとして、実装の具体的詳細を定める。「どのように動作・実装するか」に答える。

---

## 詳細仕様の読み方

本ファイルは、実装者が実装時に最初に参照する詳細仕様入口だけを扱う。方針、ポリシー、成熟度定義、ロードマップ状態、実装可否、PR 分割判断は `docs/SPEC.md` を正とし、本ファイルで再定義しない。

詳細仕様を読む順番は、`README.md` と `docs/DOCUMENT_INDEX.md` の Reading Order と同じである。本ファイルから読み始めた場合でも、先に `AGENTS.md`、`docs/DOCUMENT_INDEX.md`、`docs/SPEC.md` を確認済みでなければならない。

実装者は、対象機能ごとに以下の順で読む。

1. `docs/SPEC.md` の実装状態、Part 1 §12、§13 で、対象が実装対象であることを確認する。
2. 本ファイル §0i で、対象機能に対応する詳細仕様節と受け入れ条件を特定する。
3. 本ファイル §0a〜§0h で、詳細仕様の記載基準、共通固定値、実装前確認項目、検証条件、Phase 順序を確認する。
4. owner component の `docs/details/*.md` を主本文として読み、入力、出力、状態、正常系、異常系、セキュリティ、検証条件を確認する。
5. collaborator component がある場合は、該当する `docs/details/*.md` を呼び出し境界、schema、表示、security、setup、fixture、検証観点として確認する。
6. `docs/details/setup.md` §26 のセットアップ・アップデート手順と `docs/details/setup.md` §26.7 の受け入れ条件に影響がある場合は、実装 PR の検証対象に含める。

詳細仕様節に §0h の必須項目が不足している場合は、実装判断で補完してはならない。先に該当 owner component の詳細仕様ファイルを主本文として改訂し、本ファイルの対応表と `docs/SPEC.md` の対象範囲を整合させる。

| 範囲 | 役割 |
|------|------|
| §0〜§0j | 詳細仕様の記載基準、実装前確認項目、共通固定値、検証、Phase、詳細節対応表、リポジトリ内ソース配置 |
| §1〜§9 | `docs/details/builder.md`。`builder` / `adlaire-ci-build` の詳細仕様 |
| §10〜§20 | `docs/details/runner.md`。`runner` / `adlaire-ci-runner` の詳細仕様 |
| §21〜§22 | `docs/details/api.md`。`api` / `adlaire-ci-api` の詳細仕様。§21a は管理 API サーバー制限を定義する。 |
| §23 | `docs/details/sdk.md`。`sdk` の詳細仕様 |
| §24 | `docs/details/ui.md`。`ui` の詳細仕様 |
| §25 | `docs/details/api.md`。認証の実装仕様 |
| §26 | `docs/details/setup.md`。バイナリ配布前提のセットアップ、アップデート、受け入れ条件 |
| `admin` | `docs/details/admin.md`。管理 UI 静的ファイルの配布物構成、配置、HTTP 静的配信境界 |

責務 component 別詳細仕様ファイルの管理は、§0b.1 に従う。`docs/DETAIL_INDEX.md` は入口、索引、共通固定値、責務 component 対応表、横断補足契約だけを持つ。責務 component 別詳細仕様ファイルは、それぞれの owner component と collaborator component の詳細仕様だけを持つ。

## 0a. 詳細仕様の記載基準

Part 3 の詳細仕様項目は、実装者が追加の設計判断や推測を行わずに実装できる粒度で記載する。

仕様項目を追加または改訂する場合は、対象範囲に応じて以下を明記する。

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

未確定の内容は、実装可能な詳細仕様として記載してはならない。未確定の場合は、本ファイルまたは責務 component 別詳細仕様ファイルへ推測で具体値を記載せず、`docs/SPEC.md` で状態を確認する。

対象範囲の内容は、実装ファイルが存在しなくても、本節の基準に従って責務 component 別詳細仕様ファイルへ実装可能な粒度まで具体化する。

---

## 0b. 詳細仕様参照表

本節は、責務 component ごとに参照する詳細仕様節を示す。実装状態、実装可否、ロードマップ状態は `docs/SPEC.md` を確認する。

| 責務 component | 詳細仕様節 | 主な確認対象 |
|--------------------|------------|--------------|
| `builder` | `docs/details/builder.md` §1〜§9、§8a、§27.4、§27.25、§27.28 | CLI、入力 Markdown、静的 Web サイト出力、HTML / CSS / JavaScript、theme component、変換 report、builder fixture、builder owner 追加機能。 |
| `runner` | `docs/details/runner.md` §10〜§20、§15a、§27.2〜§27.3、§27.8〜§27.10、§27.14、§27.19、§27.21〜§27.24、§27.26〜§27.27、§27.29、§27.31〜§27.38、`docs/details/commitstatus.md` §27.1、`docs/details/setup.md` §26 | GitHub 監視、設定読取、状態ファイル更新呼び出し、pipeline、deploy、snapshot 作成トリガー、通知、runner fixture、runner owner 追加機能、Commit Status 呼び出し境界、systemd / setup 参照境界。 |
| `api` | `docs/details/api.md` §21〜§22、§21a、§25、§27.5〜§27.6、§27.11〜§27.13、§27.16〜§27.18、§27.20、§27.30、§27.42〜§27.47、`docs/details/security.md` §27.42〜§27.47、`docs/details/statefile.md` §22.0a、§22.0c | HTTP 共通契約、API server 制限、endpoint、request / response、状態ファイル read/write 呼び出し境界、認証連携、security 呼び出し境界、api owner 追加機能。 |
| `admin` | `docs/details/admin.md` §0、A1〜A5 | 管理 UI 静的ファイルの配布物構成、配置、検証、HTTP 静的配信境界。 |
| `sdk` | `docs/details/sdk.md` §23 | SDK class、method、HTTP 対応、query / body 生成、error、stream、token 破棄。 |
| `ui` | `docs/details/ui.md` §24 | 画面構成、DOM id、panel、SDK 呼び出し、表示状態、秘密情報消去。 |
| `setup` | `docs/details/setup.md` §26 | バイナリ配布、配置、systemd、セットアップ、アップデート、リリース成果物検証。 |
| `statefile` | `docs/details/statefile.md` §22.0a、§22.0c | 状態ファイル共通仕様、lock、atomic write、JSON Lines、破損時処理、状態読取 adapter、主要 schema。 |
| `archive` | `docs/details/archive.md` §27.7、§27.15 | build log archive、snapshot、download、delete、rollback、cleanup。 |
| `commitstatus` | `docs/details/commitstatus.md` §27.1 | GitHub Commit Status API payload、送信順、失敗時非反転、保存値、secret mask、検証条件。 |
| `security` | `docs/details/security.md` §27.42〜§27.47 | API token scope、API key、audit、session timeout、TOTP、rate limit、漏えい禁止、security 横断順序。 |
| `mcp` | 詳細仕様なし | 将来計画。現時点では実装可能な入出力、状態、起動手順、ツール定義、検証条件を定義しない。 |

上表の `詳細仕様節` は参照入口であり、主本文の owner component を変更しない。複数ファイルを参照する行では、対象機能の owner component のファイルを主本文とし、他ファイルは collaborator の境界、schema、fixture、security、setup、受け入れ条件を確認するために読む。参照先に同じ HTTP body、状態 schema、DOM id、SDK method、fixture assertion を重複定義してはならない。

---

## 0b.1 責務 component 別 詳細仕様ファイル管理仕様

本節は、責務 component 別に分割済みの詳細仕様ファイルを維持するための固定仕様である。責務境界の変更、仕様内容の移動、参照先更新を行う場合も、機能追加、実装状態変更、実装可否変更、ロードマップ変更、方針・ポリシー追加を含めてはならない。

詳細仕様ファイルは以下に固定する。`COMMON`、`CORE`、`BASE`、`SHARED`、`FOUNDATION`、その他の横断共通基盤ファイルは作成しない。

| ファイル | 持つ内容 | 持たない内容 |
|----------|----------|--------------|
| `docs/DETAIL_INDEX.md` | 詳細仕様の入口、読み方、共通固定値、実装前確認項目、検証マトリクス、Phase、詳細節対応表、リポジトリ内ソース配置、責務 component 別詳細仕様ファイル管理仕様。 | 各 component の詳細な処理本文、fixture / PR 証跡正本、状態 schema、API endpoint 詳細、SDK method 実装、UI DOM 詳細、setup / release 手順。 |
| `docs/details/builder.md` | `builder` owner の Markdown 変換、静的 Web サイト出力、HTML / CSS / JavaScript、theme component、builder fixture、builder owner 追加機能。 | GitHub read、runner 状態更新、API endpoint、SDK method 実装、UI DOM 詳細、状態 schema、admin 静的配信、setup / release 手順、fixture / PR 証跡正本。 |
| `docs/details/runner.md` | `runner` owner の GitHub 監視、設定読取、状態ファイル更新呼び出し、pipeline、deploy、snapshot 作成トリガー、通知、runner fixture、runner owner 追加機能。 | API endpoint の認証・応答本文、SDK method 実装、UI DOM 詳細、builder の変換処理、admin 静的配信、security 主本文、状態 schema、setup / release 手順、fixture / PR 証跡正本。 |
| `docs/details/api.md` | `api` owner の HTTP 共通契約、endpoint、request / response、状態ファイル read/write 呼び出し境界、認証連携、api owner 追加機能。 | SDK method 実装、UI DOM 詳細、runner の build 実行責務、builder の変換処理、admin 静的配信、security 主本文、状態 schema、setup / release 手順、fixture / PR 証跡正本。 |
| `docs/details/admin.md` | `admin` owner の管理 UI 静的ファイル配布物構成、配置、検証、HTTP 静的配信境界。 | UI DOM 詳細、SDK method 実装、API endpoint 実装、状態 schema、systemd 導入手順、release asset 取得手順、fixture / PR 証跡正本。 |
| `docs/details/sdk.md` | `sdk` owner の SDK class、method、HTTP 対応、query / body 生成、error、stream、token 破棄。 | API endpoint 実装、API endpoint の状態ファイル更新責務、UI DOM 詳細、状態 schema、状態ファイル直接操作、admin 静的配信、setup / release 手順、fixture / PR 証跡正本。 |
| `docs/details/ui.md` | `ui` owner の DOM id、panel、操作、表示状態、SDK 呼び出し、秘密情報消去。 | SDK method 実装、API endpoint 実装、状態 schema、状態ファイル直接操作、admin 静的配信、setup / release 手順、fixture / PR 証跡正本。 |
| `docs/details/setup.md` | `setup` owner のバイナリ配布、配置、systemd、セットアップ、アップデート、リリース成果物検証。 | runner / api / sdk / ui / admin の個別機能本文、状態 schema、API endpoint、SDK method、UI DOM、fixture / PR 証跡正本、外部依存追加。 |
| `docs/details/statefile.md` | `statefile` owner の状態ファイル共通仕様、lock、atomic write、JSON Lines、破損時処理、状態読取 adapter、主要 schema。 | API endpoint の request / response、runner の業務処理、SDK method 実装、UI 表示判断、setup / release 手順、fixture / PR 証跡正本、個別 component の業務判断。 |
| `docs/details/archive.md` | `archive` owner の build log archive / cleanup の実体処理、snapshot 保存形式、download tar.gz 生成安全性、snapshot delete 実体処理、rollback 転送実体処理。 | runner の通常 build 実行、snapshot 作成トリガー判定、API 共通 request / response、SDK method 実装、UI DOM 詳細、状態 schema、setup / release 手順、fixture / PR 証跡正本。 |
| `docs/details/commitstatus.md` | `commitstatus` owner の GitHub Commit Status API payload、送信順、失敗時非反転、保存値、secret mask、検証条件。 | runner の build 実行判断、GitHub read、API endpoint、SDK method、UI DOM 詳細、状態 schema、setup / release 手順、fixture / PR 証跡正本。 |
| `docs/details/security.md` | `security` owner の API token scope、API key、audit、session timeout、TOTP、rate limit、漏えい禁止、security 横断順序。 | API endpoint 共通処理、SDK method 実装、UI DOM 詳細、runner / builder の業務処理、状態 schema、setup / release 手順、fixture / PR 証跡正本。 |
| `docs/details/fixture.md` | fixture manifest、assertion、fake、testdata、expected / effects、受け入れ fixture 共通契約、PR 証跡テンプレート、acceptance checklist、差し戻し条件、実装 PR 完了証跡。 | 個別 component の通常処理本文、API endpoint 詳細、SDK method 実装、UI DOM 詳細、状態 schema、setup / release 実行手順。 |

`docs/DETAIL_INDEX.md` §27.38a は、runner、builder、api、sdk、ui、statefile、archive にまたがる横断補足契約であり、責務 component 別の分割先へ移動しない。§27.21〜§27.38 または api / sdk / ui / statefile の横断連動を実装する場合は、owner component の分割先詳細仕様ファイルを正本とし、§27.38a は横断処理順、同期禁止、成功後再取得、失敗時固定、横断受け入れ観点の確認として読む。

詳細仕様の配置単位は、owner component を第一基準とする。複数 component が関わる機能は、owner component のファイルに主本文を置き、collaborator component のファイルには参照リンク、禁止事項、受け入れ観点だけを置く。主本文を複数ファイルへ重複定義してはならない。

すべての責務 component 別詳細仕様ファイルは、冒頭に `## 0. 責務境界` を置き、以下の 4 項目を同じ意味で持つ。

| 項目 | 必須内容 |
|------|----------|
| owner component | そのファイルが主本文として扱う component を 1 件だけ書く。 |
| collaborator component | 呼び出し元、呼び出し先、schema 参照先、表示参照先、検証参照先を 0 件以上書く。owner component を含めてはならない。 |
| 持つ内容 | そのファイルだけが主本文として定義する入出力、状態、処理、異常系、検証条件を書く。 |
| 持たない内容 | 他 owner component へ委ねる処理、状態、API、SDK、UI、fixture、setup、security を書く。 |

責務境界表の `持つ内容` と `持たない内容` が本文と矛盾する場合は、本文を実装判断に使ってはならない。先に責務境界表、本文、§0b の詳細仕様参照表、`docs/DOCUMENT_INDEX.md` を同時に整合させる。

責務 component 別詳細仕様ファイルの各節は、以下を満たす。

| 項目 | 必須条件 |
|------|----------|
| 節番号 | 既存の節番号を維持する。番号の再採番は行わない。 |
| 参照 | 入口ファイルと owner component ファイルの参照先が一意に追跡できるよう、`docs/DETAIL_INDEX.md` の対応表を更新する。 |
| owner | 各機能節に owner component を 1 件だけ明記する。 |
| collaborator | collaborator component は 0 件以上を明記し、owner component を含めない。 |
| 重複禁止 | 同じ入力、出力、状態 schema、HTTP body、DOM id、fixture assertion を複数ファイルで重複定義しない。 |
| 横断事項 | 横断する固定値は `docs/DETAIL_INDEX.md` に置く。横断共通基盤を component として扱わない。 |
| 索引 | `docs/DOCUMENT_INDEX.md` に、分割後ファイルの役割と正本範囲を反映する。 |

実装者が詳細仕様を読む順序は以下に固定する。

1. `docs/SPEC.md` で実装対象、実装状態、実装可否を確認する。
2. `docs/DETAIL_INDEX.md` §0〜§0j で共通固定値、責務 component、詳細節対応表、リポジトリ内ソース配置を確認する。
3. owner component の分割先詳細仕様ファイルを主本文として読む。
4. collaborator component がある場合は、該当する分割先詳細仕様ファイルの参照節を呼び出し境界、schema、表示、security、setup、fixture、検証観点として読む。
5. 状態ファイルの読み書き、lock、atomic write、schema を扱う場合は `docs/details/statefile.md` を読む。
6. 認証、scope、token、audit、session、TOTP、rate limit、漏えい禁止を扱う場合は `docs/details/security.md` を読む。
7. fixture、fake、PR 証跡が必要な場合は `docs/details/fixture.md` を読む。

詳細仕様ファイルの責務整理、節移動、参照先更新は、以下の完了条件をすべて満たすまで完了扱いにしてはならない。

| 完了条件 | 判定 |
|----------|------|
| 旧ファイル内の移動対象本文が対応する分割先に移動している。 | 必須 |
| `docs/DETAIL_INDEX.md` には入口、索引、共通固定値、対応表、管理仕様、横断補足契約だけが残っている。 | 必須 |
| `docs/SPEC.md`、`docs/DOCUMENT_INDEX.md`、各分割先ファイル間の参照が矛盾していない。 | 必須 |
| `rg` で旧節名、旧ファイル名、移動前参照の取り残しを確認している。 | 必須 |
| 実装ファイル、fixture、testdata の内容を分割作業だけで変更していない。 | 必須 |

---

## 0c. 実装前確認項目

実装者は、対象機能について以下の条件をすべて満たすまで実装を開始してはならない。

| ゲート | 合格条件 |
|--------|----------|
| 対応表 | 対象機能が §0i の詳細節対応表に記載され、詳細仕様節と受け入れ条件が一意に示されている。 |
| テンプレート | 対象機能の詳細仕様が §0h の機能仕様テンプレートの必須項目を満たしている。 |
| 責務境界 | owner component、collaborator component、対象ファイル、呼び出し元、呼び出し先、変更対象として許可された状態ファイルが明記されている。 |
| 入出力 | すべての入力、出力、既定値、許容値、必須/任意、型、文字コード、時刻形式が明記されている。 |
| 状態管理 | 状態ファイルのパス、JSON 形式、更新タイミング、初期状態、破損時の扱い、権限が明記されている。 |
| 正常系 | 処理順序、分岐条件、ループ条件、成功条件、終了条件が明記されている。 |
| 異常系 | エラー条件、ログレベル、HTTP ステータス、戻り値、再試行有無、処理継続/中断条件が明記されている。 |
| 冪等性 | 同一リクエスト、再実行、途中失敗後の再開で二重実行・二重削除・状態破壊が発生しない条件が明記されている。 |
| 排他制御 | 同時実行、ロック、タイムアウト、ロック残存時の扱いが明記されている。 |
| セキュリティ | 秘密情報の保存禁止、マスク、ファイル権限、認証/認可、外部公開可否が明記されている。 |
| 検証 | 構文確認、単体確認、手動 API 確認、生成物確認、ログ確認、失敗系確認のいずれを行うかが明記されている。 |

上記ゲートのいずれかが未充足の場合、実装判断で補完してはならない。先に該当 owner component の詳細仕様ファイルを改訂し、必要に応じて collaborator 詳細仕様、本ファイルの対応表、`docs/SPEC.md` を同じ仕様 PR で整合させ、未充足項目を仕様として確定する。

実装後の完了条件は以下とする。

1. 実装した機能が、該当 owner component の詳細仕様ファイルと本ファイルの対応表に記載された入力、出力、状態、異常系、検証条件と一致する。
2. 対象機能が owner component の詳細仕様本文で §0h の機能仕様テンプレートを満たし、§0i の詳細節対応表の受け入れ条件を満たしている。
3. 実装対象外に残す機能が PR 本文に明記されている。
4. `docs/SPEC.md`、`docs/DETAIL_INDEX.md`、`docs/DOCUMENT_INDEX.md`、`AGENTS.md` のファイル名参照が矛盾していない。
5. 実装ファイルを変更した場合、構文確認または実行確認の結果が記録できる。
6. 仕様外の挙動、暗黙の既定値、未記載の状態ファイル、未記載のエラー応答が存在しない。

---

## 0d. 共通固定値

本節は、各コンポーネントで共通して使用する固定値を定義する。個別節に別の値が明記されていない限り、本節を優先する。

| 項目 | 決定 |
|------|------|
| Go 最小バージョン | Go `1.22` 以上。標準ライブラリのみを使用し、外部 module は追加しない。 |
| 文字コード | 入力、出力、状態ファイル、HTTP body は UTF-8 固定。UTF-8 として読み取れない入力は処理を中断する。 |
| 改行 | 新規に書き出す text / JSON Lines ファイルは LF 固定。CRLF 入力は読み込み時に LF として扱う。 |
| 時刻 | 状態ファイル、API、ログの機械処理用時刻は UTC の ISO 8601 形式（例: `2026-09-16T09:00:00Z`）で保存する。ローカル時刻への変換は UI 表示に限定し、状態ファイル、API response、ログには保存しない。 |
| JSON | JSON object の未知キーは保存しない。読み込み時に未知キーを見つけた場合は無視し、次回保存時に除去する。 |
| atomic write | 状態ファイル更新手順は `docs/details/statefile.md` §22.0a を正とする。各 component は同節の手順を使用し、独自更新手順を持たない。 |
| 権限 | 状態ファイル、秘密情報ファイル、ディレクトリの権限は `docs/details/statefile.md` §22.0a を正とする。 |
| ロック | 状態ファイル lock の作成、待機、解除、競合時応答は `docs/details/statefile.md` §22.0a を正とする。 |
| ログ秘密情報 | PAT、Webhook Secret、SMTP password、session token、API token は stdout、stderr、JSON log、API response、UI 表示へ平文出力しない。表示が必要な場合は `"***"` とする。 |
| 終了コード | CLI / runner は `0` 成功、`1` 一般エラー、`2` 入力・設定エラー、`3` 外部サービス・ネットワークエラー、`4` ロック競合を標準とする。個別節に明記がある場合もこの意味から外してはならない。 |
| 禁止事項 | 仕様にない環境変数、状態ファイル、HTTP endpoint、CLI option、外部依存を実装者判断で追加してはならない。必要な場合は先に該当 owner component の詳細仕様ファイルまたは本ファイルの対応表を改訂する。 |

---

## 0e. 完全実装検証マトリクス

対象項目を完了扱いにする場合は、責務 component ごとに下表の検証を満たす。実装ファイルが存在しても、本表の必須検証が未完了の場合は完了扱いにしない。

本表の `builder`、`runner`、`api`、`sdk`、`ui`、`setup` は Phase の主対象 component である。`statefile`、`security`、`archive`、`commitstatus`、`admin`、`fixture` は、主対象 component の collaborator component として完了判定に参加する。collaborator component の検証が失敗する場合、主対象 component の実装も完了扱いにしてはならない。

実装完了判定は、機能実装、fixture / testdata、fake、検証結果、PR 証跡を 1 組として扱う。本ファイルは完了判定の入口と参照順だけを示し、fixture 名、expected / effects、fake 動作、PR 証跡項目、acceptance checklist、差し戻し条件は `docs/details/fixture.md` §0g.8-F、§22-F、§27-F を正とする。コードが仕様どおりに見える場合でも、同ファイルに定義された対象 fixture、expected / effects、実行コマンド、実結果、secret 確認、対象外確認が不足する場合は完了扱いにしない。

| 対象 | 必須検証 | 合格条件 |
|------|----------|----------|
| `builder` | CLI 正常系 | `adlaire-ci-build --src <valid.md-or-dir> --out <site-dir>` が終了コード `0` で終了し、静的 Web サイトと `[REPORT]` を生成する。 |
| `builder` | CLI 異常系 | 入力不存在、UTF-8 不正、未知引数、出力不可ディレクトリ、未知 theme で §2・§8 の終了コードと stderr が一致する。 |
| `builder` | Markdown 変換 | 見出し、重複 slug、内部リンク警告、脚注、表、引用、リスト、コードフェンス、未閉鎖フェンス、HTML escape が §4 の出力構造と一致する。 |
| `builder` | 生成物 | 出力サイトディレクトリに `index.html`、ページ HTML、`assets/style.css`、`assets/app.js`、`assets/search-index.json` が生成され、§5〜§7 の ID / class / JS 機能を含む。 |
| `runner` | 設定検証 | `--state-dir`、`BRANCH_TARGETS`、必須ファイル不足、未知設定キーで §12 のログ・終了コード・採用優先順位が一致する。 |
| `runner` | 状態更新 | 成功、ビルド失敗、GitHub API 失敗、転送失敗、lock 競合、JSON 破損で §13 と §22.0a の更新順序・未更新条件が一致する。 |
| `runner` | 冪等性 | 同一 SHA 再実行、pending retry 再実行、通知 pending 再実行、stale lock 復旧で二重履歴・二重 snapshot・状態破壊が発生しない。 |
| `api` | API 共通 | 未知 path、未対応 method、body 禁止、JSON 不正、body 上限、認証なし、権限不足、入力検証失敗、ロック競合が §22.0 の status と body を返す。 |
| `api` | 状態ファイル | 全 write API が §22.0a / §22.0d の対象ファイルだけを atomic write し、秘密情報を平文出力しない。 |
| `api` | endpoint 契約 | §22.0e の全 endpoint について Request、Response、Success、Errors、Read、Write、SDK、UI の対応が実装と一致する。 |
| `sdk` | SDK 契約 | 全 method が §22.0e の endpoint のみを呼び、query / body 生成、body なし endpoint の body 禁止、HTTP error の `AdlaireCIError` 変換が §23 と一致する。 |
| `ui` | UI 契約 | 全操作が §24 の SDK method 経由で動作し、成功表示、失敗表示、disabled、再取得、秘密情報消去が一致する。 |
| `setup` | systemd | `docs/details/setup.md` §26 の unit 名、`ExecStart`、配置パス、権限、起動確認コマンドが実際の導入手順と一致する。 |
| `statefile` | 状態ファイル契約 | 状態ファイルの schema、lock、atomic write、JSON Lines、破損時処理、権限、秘密情報マスクが `docs/details/statefile.md` §22.0a、§22.0c と一致する。 |
| `security` | 認証・認可・漏えい禁止 | scope、API key、audit、session、TOTP、rate limit、秘密情報非表示、失敗時副作用が `docs/details/security.md` §27.42〜§27.47 と一致する。 |
| `archive` | artifact / log archive | gzip archive、snapshot、download、delete、rollback、cleanup の実体処理が `docs/details/archive.md` §27.7、§27.15 と一致し、API / SDK / UI の応答契約を上書きしない。 |
| `commitstatus` | GitHub Commit Status | payload、送信順、失敗時非反転、保存値、secret mask、検証条件が `docs/details/commitstatus.md` §27.1 と一致し、runner の build 実行判断を上書きしない。 |
| `admin` | 静的配布境界 | admin 配布物、archive validation、HTTP 静的配信、setup 連携が `docs/details/admin.md` §0、A1〜A5 と一致し、UI / SDK の本文を重複定義しない。 |
| `fixture` | fixture / fake / 証跡 | Phase 別 fixture、fake、assertion、expected / effects、PR 証跡、acceptance checklist、差し戻し条件が `docs/details/fixture.md` §0g.8-F、§22-F、§27-F と一致する。 |

検証結果は、実装 PR の本文または実装完了報告に、対象、実行コマンド、fixture 名、期待結果、実結果、状態差分、外部副作用、secret 確認、対象外確認を対応付けて記録する。記録形式、必須項目、不足時の扱いは `docs/details/fixture.md` の PR 証跡契約を正とする。検証不能な項目がある場合は、その項目を完了扱いにしてはならない。

---

## 0f. 仕様策定完了チェック

本節は、Go 版初期実装へ進む前に仕様策定が完了しているかを判定するチェックである。実装者は、責務 component ごとに下表の必須条件を満たすまで実装を開始してはならない。

| 対象 | 実装着手条件 | 実装禁止条件 | 完了判定 |
|------|--------------|--------------|----------|
| `builder` | §2〜§8 に CLI option、入力 Markdown、出力サイトディレクトリ、終了コード、stderr、HTML 構造、テーマコンポーネント、JS/CSS、生成物確認が定義されている。 | §4〜§7 にない Markdown 記法、CSS class、JavaScript 機能、外部 asset、theme を追加すること。 | §0e の `builder` 必須検証をすべて満たし、生成サイトが §5〜§7 と一致する。 |
| `runner` | §10〜§20 に設定値、状態ファイル、GitHub API、SHA 比較、pipeline 起動、SSH 転送、snapshot、通知、ログ、systemd / setup 参照境界が定義されている。systemd unit 本文とセットアップ手順は `docs/details/setup.md` §26 を正とする。 | 未定義の環境変数、状態ファイル、queue 挙動、通知チャンネル、pipeline 形式を追加すること。systemd unit file の生成、配置、更新、enable、restart を runner に追加すること。 | §0e の `runner` 必須検証をすべて満たし、状態ファイル更新順序が §13、§22.0a、§22.0d と一致する。 |
| `api` | `docs/details/api.md` §21〜§22、§21a、§25 と `docs/details/setup.md` §26 に API 共通契約、API server 制限、endpoint、状態ファイル schema、認証、認可、systemd、セットアップが定義されている。 | `docs/details/api.md` §22.0e にない endpoint、method、status code、response body、状態ファイル write を追加すること。 | `docs/details/api.md` §22.0e の全 endpoint が Request、Response、Errors、Read、Write、SDK、UI の対応表と一致する。 |
| `sdk` | `docs/details/sdk.md` §23 に SDK class、method、引数、戻り値、HTTP endpoint 対応、error object、token 破棄条件が定義されている。 | SDK が `docs/details/api.md` §22.0e にない endpoint を呼ぶこと、body 禁止 endpoint に body を送ること、独自 error 形式を返すこと。 | 全 method が `docs/details/api.md` §22.0e と `docs/details/sdk.md` §23 の対応どおりに動作し、HTTP error を `AdlaireCIError` として扱う。 |
| `ui` | `docs/details/ui.md` §24 に画面構成、panel、操作、成功表示、失敗表示、disabled、再取得、秘密情報消去が定義されている。 | SDK を介さず API を直接呼ぶこと、未定義の画面・操作・保存先を追加すること、秘密情報を DOM に残すこと。 | 全 UI 操作が `docs/details/ui.md` §24 の表示条件と `docs/details/sdk.md` §23 の SDK method を満たし、秘密情報 field が指定条件で消去される。 |

上表の対象外である `mcp`、MCP tools、MCP resources、MCP prompts、HTTP SSE transport、MCP audit / stats / config CRUD は、初期実装では実装しない。これらは、Part 3 詳細仕様セット内に入出力、状態、起動手順、検証条件を定義しない。

仕様策定完了チェックで未充足が見つかった場合は、実装を開始せず、以下の順で仕様を改訂する。

1. 未充足項目が本ファイルの記載対象外である場合は、先に `docs/SPEC.md` を確認する。
2. 未充足項目が入出力、状態ファイル、api、sdk、ui、処理順序、異常系、検証条件に関わる場合は、該当 owner component の詳細仕様ファイルまたは collaborator の詳細仕様ファイルを改訂する。
3. ファイル名、正本関係、対象範囲が変わる場合は、`docs/DOCUMENT_INDEX.md` の更新要否を確認する。
4. 対象項目の詳細節、参照先、受け入れ条件が変わる場合は、§0i の詳細節対応表を更新する。
5. 改訂後、§0b、§0c、§0e、本節、§0g、§0h、§0i、§0j の条件を再確認する。

---

## 0g. 初期実装 Phase 分割

Go 版初期実装は、`docs/SPEC.md` §0e の対象範囲を一括実装せず、下表の Phase 順に進める。上位 Phase の完了判定を満たす前に、下位 Phase の実装 PR を開始してはならない。

実装単位、実装 PR 単位、完了判定単位は Phase のみとする。`P0`、`P1`、`P2〜P5` などの優先度ラベル、抽象段階、API 内部段階名を実装単位として使ってはならない。API の実装範囲は、Phase 3 を「API 基盤・認証・状態 read/write・運用基本操作」、Phase 4 を「API 拡張運用操作」として扱う。

各 Phase の `対象` は、その Phase の owner component を示す。状態ファイル、security、archive、commitstatus、admin、fixture、setup が関わる場合も、それらは collaborator component として該当 Phase の完了条件に含める。collaborator component の詳細仕様に未充足がある場合は、owner component の実装で補完せず、先に該当する責務 component 別詳細仕様ファイルを改訂する。

| Phase | 対象 | 実装範囲 | 依存条件 | 完了条件 |
|-------|------|----------|----------|----------|
| Phase 1 | `builder` | §2〜§9 の CLI、Markdown 変換、静的 Web サイト出力、テーマコンポーネント、生成物確認。 | なし。 | §0e の `builder` 必須検証と §0f の `builder` 完了判定を満たす。 |
| Phase 2 | `runner` | §10〜§20 の CI ランナー、GitHub API 連携、SHA キャッシュ、pipeline 起動、SSH 転送、snapshot、通知、ログ、systemd / setup 参照境界。systemd unit 本文と配置手順は `docs/details/setup.md` §26 を正とする。 | Phase 1 が完了し、`adlaire-ci-build` の CLI 契約が固定されている。 | §0e の `runner` 必須検証と §0f の `runner` 完了判定を満たす。 |
| Phase 3 | `api` | `docs/details/api.md` §21〜§22、§21a、§25 と `docs/details/setup.md` §26 のうち、認証、セッション、共通エラー、状態ファイル読み書き、ビルド操作、status、logs、history、queue、circuit breaker。 | Phase 2 が完了し、runner が書き込む状態ファイル schema が固定されている。 | API 基盤・認証・状態 read/write・運用基本操作の必須検証、§0e の `api` API 共通・状態ファイル検証、§0f の `api` 完了判定の該当範囲を満たす。 |
| Phase 4 | `api` | `docs/details/api.md` §22.0e、§22.0f のうち、config、repo、branch、schedule、notify、snapshot、rollback、maintenance、access control、hooks、tokens 等の拡張運用操作。 | Phase 3 が完了し、API 共通処理と認証が固定されている。 | API 拡張運用操作の必須検証と §0e の `api` endpoint 契約を満たす。 |
| Phase 5 | `sdk` | `docs/details/sdk.md` §23 の SDK class、method、戻り値、HTTP error、token 破棄、query / body 生成。 | Phase 3 と Phase 4 が完了し、`docs/details/api.md` §22.0e の endpoint 契約が固定されている。 | §0e の `sdk` 契約と §0f の `sdk` 完了判定を満たす。 |
| Phase 6 | `ui` | `docs/details/ui.md` §24 の標準管理ツール UI、DOM id、panel、操作、SDK 呼び出し、成功表示、失敗表示、disabled、再取得、秘密情報消去。 | Phase 5 が完了し、SDK method 契約が固定されている。 | §0e の `ui` 契約と §0f の `ui` 完了判定を満たす。 |

### 0g.1 Phase 1 完全仕様ゲート（`builder`）

Phase 1 の実装詳細本文は `docs/details/builder.md` を正とする。親ファイルでは、Phase 1 の対象、依存条件、完了条件、後続 Phase への引き継ぎ確認だけを扱う。

| 確認 | 参照先 |
|------|--------|
| CLI、Markdown 変換、静的 Web サイト出力、theme component、生成物確認 | `docs/details/builder.md` §1〜§9、§8a |
| Phase 1 fixture、testdata、PR 証跡 | `docs/details/fixture.md` §0g.8-F |
| release / setup 受け入れ条件 | `docs/details/setup.md` §26.7 |

### 0g.2 Phase 2 完全仕様ゲート（`runner`）

Phase 2 の実装詳細本文は `docs/details/runner.md` を正とする。親ファイルでは、Phase 2 が Phase 1 の `adlaire-ci-build` 契約に依存し、後続 API が読む runner 状態契約を固定することだけを扱う。

| 確認 | 参照先 |
|------|--------|
| CI runner、GitHub API 連携、SHA cache、pipeline、deploy、snapshot、通知、systemd | `docs/details/runner.md` §10〜§20、§15a |
| runner fixture、fake GitHub、fake ssh / notifier、PR 証跡 | `docs/details/fixture.md` §0g.8-F |
| release / setup 受け入れ条件 | `docs/details/setup.md` §26.7 |

### 0g.3 Phase 3 完全仕様ゲート（`api`）

Phase 3 の実装詳細本文は `docs/details/api.md` と `docs/details/setup.md` を正とする。親ファイルでは、API 共通契約、認証、状態 read/write、運用基本 endpoint が Phase 4〜6 の前提になることだけを扱う。

| 確認 | 参照先 |
|------|--------|
| API 共通処理、API server 制限、認証、運用基本 endpoint、状態 read/write | `docs/details/api.md` §21〜§22、§21a、§25、§22.0f |
| 状態ファイル schema、lock、atomic write | `docs/details/statefile.md` §22.0a、§22.0c |
| 管理 API 導入、systemd、release 受け入れ条件 | `docs/details/setup.md` §26 |

### 0g.4 Phase 4 完全仕様ゲート（`api`）

Phase 4 の実装詳細本文は `docs/details/api.md` を正とする。親ファイルでは、拡張運用 endpoint が SDK / UI の最終入力契約になることだけを扱う。

| 確認 | 参照先 |
|------|--------|
| 拡張運用 endpoint、request / response、error、auth、secret mask | `docs/details/api.md` §22.0e、§22.0f |
| security 連携 | `docs/details/security.md` §27.42〜§27.47 |
| API fixture、endpoint 証跡 | `docs/details/fixture.md` §0g.8-F、§22-F、§27-F |

### 0g.5 Phase 5 完全仕様ゲート（`sdk`）

Phase 5 の実装詳細本文は `docs/details/sdk.md` を正とする。親ファイルでは、SDK が固定済み API endpoint だけを呼び、UI 表示判断を持たないことだけを扱う。

| 確認 | 参照先 |
|------|--------|
| SDK class、method、HTTP 対応、query / body 生成、error、stream、token 破棄 | `docs/details/sdk.md` §23 |
| API endpoint 対応 | `docs/details/api.md` §22.0e |
| SDK fixture、fake fetch / stream | `docs/details/fixture.md` §0g.8-F |

### 0g.6 Phase 6 完全仕様ゲート（`ui`）

Phase 6 の実装詳細本文は `docs/details/ui.md` を正とする。親ファイルでは、UI が SDK 経由だけで API と通信し、秘密情報を DOM に残さないことだけを扱う。

| 確認 | 参照先 |
|------|--------|
| DOM id、panel、操作、表示状態、SDK 呼び出し、秘密情報消去 | `docs/details/ui.md` §24 |
| SDK method 契約 | `docs/details/sdk.md` §23 |
| UI fixture、fake SDK、直接 API 呼び出し禁止確認 | `docs/details/fixture.md` §0g.8-F |

### 0g.7 Phase 間引き継ぎ契約

各 Phase の完了時は、次 Phase が依存する契約を変更不可として扱う。後続 Phase で変更が必要になった場合は、後続 Phase の実装で吸収せず、契約を定義した owner 詳細仕様へ戻す。

| 引き継ぎ元 | 引き継ぎ先 | 固定する契約 | 正本 |
|------------|------------|--------------|------|
| Phase 1 | Phase 2 | `adlaire-ci-build` CLI、終了コード、stdout / stderr、`[REPORT]`、出力サイト構造。 | `docs/details/builder.md` |
| Phase 2 | Phase 3 | runner 状態ファイル schema、lock、history/log、pending queue、circuit breaker、snapshot、通知ログ。 | `docs/details/runner.md` / `docs/details/statefile.md` |
| Phase 3 | Phase 4 | API 共通契約、認証、session、error body、validation、lock error、SSE 基本形式。 | `docs/details/api.md` |
| Phase 4 | Phase 5 | 全 endpoint の method、path、query、request、response、error、認証要否。 | `docs/details/api.md` |
| Phase 5 | Phase 6 | SDK method 名、引数、戻り値、error object、stream handle、token 破棄条件。 | `docs/details/sdk.md` |
| Phase 6 | 初期実装完了 | UI 操作、表示状態、secret 消去、SDK 経由通信、実装完了検証結果。 | `docs/details/ui.md` |

### 0g.8 Phase 別 実装 PR 成果物チェックリスト

Phase fixture / testdata 配置、fake 実装、実装 PR 証跡の詳細は `docs/details/fixture.md` §0g.8-F を正とする。親ファイルでは、Phase ごとの成果物参照先だけを保持し、fixture 名、expected / effects、fake 動作、PR 証跡項目を重複定義しない。

| Phase | 実装対象 | 成果物・fixture 正本 | 受け入れ条件 |
|-------|----------|----------------------|--------------|
| Phase 1 | `builder` | `docs/details/fixture.md` §0g.8-F | `docs/details/builder.md` と `docs/details/setup.md` §26.7 を満たす。 |
| Phase 2 | `runner` | `docs/details/fixture.md` §0g.8-F | `docs/details/runner.md` と `docs/details/setup.md` §26.7 を満たす。 |
| Phase 3 | `api` | `docs/details/fixture.md` §0g.8-F、§27-F | `docs/details/api.md` の API 基盤・認証・状態 read/write・運用基本操作と setup API 導入条件を満たす。 |
| Phase 4 | `api` | `docs/details/fixture.md` §0g.8-F、§27-F | `docs/details/api.md` の API 拡張運用操作を満たす。 |
| Phase 5 | `sdk` | `docs/details/fixture.md` §0g.8-F | `docs/details/sdk.md` §23 を満たす。 |
| Phase 6 | `ui` | `docs/details/fixture.md` §0g.8-F | `docs/details/ui.md` §24 を満たす。 |

---

## 0h. 機能仕様テンプレート

対象項目を追加または改訂する場合は、該当する owner component の詳細仕様節に以下の項目をすべて含める。既存節に含める場合も、実装者が下表の項目を本文から一意に読み取れる状態にする。

| 項目 | 必須内容 | 未記載時の扱い |
|------|----------|----------------|
| 目的 | 何を解決する機能か、どの利用者または運用者のための機能か。 | 実装不可。 |
| 責務 component | owner component と collaborator component を明記し、複数 component が関わる場合は責務境界を分けて書く。 | 実装不可。 |
| 入力 | CLI 引数、HTTP request、設定値、状態ファイル、環境変数、Markdown 入力、UI 操作などの入力元、型、必須/任意、既定値。 | 実装不可。 |
| 出力 | 生成ファイル、HTTP response、stdout/stderr、ログ、通知、UI 表示、終了コード。 | 実装不可。 |
| 状態 | 読み書きする状態ファイル、ディレクトリ、メモリ状態、ロック、更新責務、初期値、破損時の扱い。 | 状態を持つ実装は禁止。 |
| 正常系 | 処理順序、分岐条件、成功条件、保存順序、外部コマンド呼び出し条件。 | 実装不可。 |
| 異常系 | エラー条件、継続/中断、HTTP status、終了コード、ログレベル、通知、リトライ有無。 | 実装不可。 |
| セキュリティ | 秘密情報、認証、認可、ファイル権限、外部公開可否、ログ出力禁止事項。 | セキュリティ影響がある機能は実装不可。 |
| 検証 | 必須テスト、手動確認、fixture、生成物確認、API 確認、異常系確認。 | 完了扱い不可。 |
| 完了条件 | どの検証が成功したら実装完了と扱うか。関連文書の更新要否。 | 完了扱い不可。 |

上表のいずれかが不足する対象項目は、実装者判断で補完してはならない。不足を見つけた場合は、実装 PR ではなく仕様改訂 PR として該当 owner component の詳細仕様ファイルを主本文として先に更新し、本ファイルの対応表は参照先変更がある場合だけ更新する。

---

## 0i. 詳細節対応表

本節は、対象機能から該当する詳細仕様へ移動するための対応表である。実装者は対象機能を実装する前に、下表の「詳細仕様節」と「受け入れ条件」を確認する。

表の「責務 component」は参照先を探すための component 一覧である。owner component と collaborator component は、対象機能の詳細仕様節に記載された値を正とする。

表の「詳細仕様節」が複数ある場合は、owner component の詳細仕様ファイルを主本文として読み、collaborator component の詳細仕様ファイルは schema、呼び出し境界、表示、security、setup、fixture、検証観点の確認として読む。ファイル名を伴わない裸の節番号は、同じ行の「責務 component」から該当する owner component または collaborator component の詳細仕様ファイルへ解決する。

詳細節対応表は owner component を置き換える表ではない。受け入れ条件が複数 component にまたがる場合でも、主本文は owner component の詳細仕様ファイルを正とし、collaborator component の詳細仕様は schema、呼び出し境界、表示、security、setup、fixture、検証観点の確認に限定する。collaborator component は、owner component の入力、出力、状態、endpoint、SDK method、UI 操作を追加定義しない。

該当節に §0h の必須項目が不足している場合は、その項目を実装せず、先に詳細仕様を改訂する。§27.1〜§27.47 の owner、主本文、collaborator は §27.1〜§27.47 追加仕様化機能参照インデックスを確認する。

### 0i.1 Builder / 静的 Web サイト出力

| 機能 | 責務 component | 詳細仕様節 | 受け入れ条件 |
|------|-------------------|------------|--------------|
| 出力サイトサイズ警告閾値 | `builder` / `runner` / `api` | §8、§12、§13、§22.0e | `OUTPUT_SIZE_WARN_MB`、`size_warn`、WARN ログ、API 表示が一致する。 |
| 出力サイトへのビルドメタ埋め込み | `builder` / `runner` / `api` | §2、§5、§8、§13、§22.0e、§27.4 | CLI/env 入力、HTML meta、REPORT、build log、`GET /api/output-meta` の値が一致する。 |
| 変換レポート出力 | `builder` / `runner` / `api` | §8、§13、§15、§22.0e | `[REPORT]` stdout、runner 取り込み、`.build_logs`、`GET /api/output-meta` が一致する。 |
| シンタックスハイライト | `builder` | §7.8 | 対応言語、class 名、HTML escape、CSS 表示が一致する。 |
| 本文内全文検索 | `builder` | §7.9 | `assets/search-index.json`、検索 UI、ヒット遷移、対象テキストが一致する。 |
| アンカーリンク自動検証 | `builder` | §4.3、§8 | broken anchor 検出、`[WARN] BROKEN_LINK`、report 件数が一致する。 |
| コードブロックの折りたたみ | `builder` | §7.10 | 30 行超の初期折りたたみ、展開操作、印刷時展開が一致する。 |
| 印刷スタイル（`@media print`） | `builder` | §6 | `@media print` の非表示対象、コード展開、リンク URL 表示が一致する。 |
| 静的 Web サイト出力 | `builder` | §2、§5、§6、§7 | 入力ファイル/ディレクトリ、出力ファイル構成、asset、ページ生成が一致する。 |
| テーマコンポーネント | `builder` | §5、§6、§7 | `adlaire-default` の component、class、slot、asset 出力が一致する。 |
| 外部リンクの自動処理 | `builder` | §4.3 | `target="_blank"`、`rel="noopener noreferrer"`、内部リンクとの区別が一致する。 |
| 読み取り進捗バー | `builder` | §7.13 | 3px 固定表示、scroll 連動、初期/末尾状態が一致する。 |
| コードブロックのコピーボタン | `builder` | §7.6 | ボタン配置、コピー対象、成功/失敗時表示、アクセシビリティが一致する。 |
| 見出しアンカーリンクコピー | `builder` | §3、§7.11 | `.hn-link`、copy URL、重複 slug 連動が一致する。 |
| TOC 開閉状態の永続化 | `builder` | §7.3 | `localStorage` key、展開/折りたたみ、復元条件が一致する。 |
| 見出しスラグ重複解決 | `builder` | §4.5 | `-2`、`-3` の付与、TOC、検索、コピー URL との共通化が一致する。 |
| 前後章ナビゲーションボタン | `builder` | §4.5、§5、§7.15 | h2 単位の前後判定、章末尾配置、端の非表示条件が一致する。 |
| 内部リンク整合性チェック | `builder` | §4.3、§8 | `[label](#anchor)` 検証、WARN、`broken_links` が一致する。 |
| 見出し階層スキップ警告 | `builder` | §4.5、§8 | h1→h3 等の検出、WARN、`heading_skips` が一致する。 |
| 読了時間推計と表示 | `builder` | §4.5、§5、§6、§8 | 対象文字数、200文字/分、切り上げ、header 表示、report が一致する。 |
| テーブルのソート機能 | `builder` | §7.14 | クリック操作、昇順/降順、`aria-sort`、インジケーターが一致する。 |
| キーボードショートカット | `builder` | §7.12 | `/`、`Escape`、`t` の対象、フォーカス条件、入力中の無効化が一致する。 |
| ビルドキャッシュ | `builder` / `runner` | §5、§8、§11、§13、§27.25 | `.build_cache.json`、入力 manifest、再利用条件、無効化条件、report counters が一致する。 |
| 依存ファイルトラッキング | `builder` / `runner` | §4.3、§5、§11、§13、§27.28 | `.dependency_manifest.json`、依存抽出、関連 target 判定、破損時 full build が一致する。 |

### 0i.2 Runner / CI 実行

| 機能 | 責務 component | 詳細仕様節 | 受け入れ条件 |
|------|-------------------|------------|--------------|
| ビルドタイムアウト | `runner` / `api` | §12、§13、§22.0e | `build_timeout_seconds` の既定値、設定 API、`context.WithTimeout` の中断処理、終了コード、ログが一致する。 |
| ビルドログのファイル保存 | `runner` | §11、§13、§15 | `.build_logs/{id}.json` の schema、stdout/stderr、変換レポート、duration、権限が一致する。 |
| ネットワーク断時の再試行 | `runner` | §12、§13 | `API_RETRY_MAX`、`API_RETRY_BASE_SECONDS`、指数バックオフ、失敗時ログが一致する。 |
| GitHub API レート制限自動待機 | `runner` / `api` | §13、§22.0e | `X-RateLimit-Remaining` と `X-RateLimit-Reset` の扱い、待機、API 表示が一致する。 |
| 転送後リモート整合性検証 | `runner` | §14a、§13 | SSH 転送後の SHA256 照合、不一致時の `.pending_transfers` 再投入、ログが一致する。 |
| マルチブランチビルド | `runner` / `api` | §12、§13、§22.0e | `BRANCH_TARGETS` と `.branch_config` の優先順位、順次処理、API 更新が一致する。 |
| ビルドログ世代管理 | `runner` | §12、§13、§15 | `LOG_KEEP_N` 超過時の削除順序、0 の扱い、削除ログが一致する。 |
| ビルド出力の外部転送 | `runner` | §14a、§13 | SSH 差分転送、複数ファイル処理、失敗時 pending、通知が一致する。 |
| ビルドクールダウン | `runner` | §12、§13 | `BUILD_COOLDOWN_SECONDS` 内の起動スキップ、Webhook 二重トリガー抑止、ログが一致する。 |
| ビルド前の事前チェック | `runner` | §13、`docs/details/setup.md` §26 | ディスク、`adlaire-ci-build`、pipeline 前提の確認、不足時の ERROR と通知が一致する。 |
| 定期強制ビルド | `runner` / `api` | §12、§13、§22.0e | `FORCE_BUILD_INTERVAL`、変更なし時の強制ビルド、設定 API が一致する。 |
| ビルド中重複スキップ | `runner` | §11、§13 | `.build_lock` の PID 判定、stale lock、競合時終了コードとログが一致する。 |
| GitHub PAT 有効期限の事前警告 | `runner` / `api` | §13、§22.0e | `GitHub-Authentication-Token-Expiration` の解析、7 日以内 WARN、API 表示が一致する。 |
| コミット情報のビルドログ記録 | `runner` | §13、§15 | SHA、message、author、date を build id と同じログへ記録する。 |
| GitHub API 連続失敗によるサーキットブレーカー | `runner` / `api` | §11、§12、§13、§22.0e | 閾値、open/close 状態、API reset、通知、状態ファイルが一致する。 |
| 設定ファイル起動時整合性チェック | `runner` | §11、§12、§13、§22.0a、§22.0c、§27.10 | 対象 JSON ファイル、検証順序、破損退避、初期化値、ログ、通知、終了コード、fixture が一致する。 |
| ビルドステータスファイル出力 | `runner` / `api` | §11、§13、§15、§22.0a、§22.0c、§22.0e、§27.8 | `.build_status.json` の schema、更新タイミング、status/target_status、pending 件数、circuit 状態、API 参照元が一致する。 |
| ビルドトリガー種別の記録 | `runner` / `api` / `sdk` / `ui` | §13、§15、§22.0c、§22.0e、§23、§24、§27.9 | `trigger` の有効値、判定条件、`.build_logs`、`.build_history`、`.build_status.json`、履歴 filter、UI 表示が一致する。 |
| GitHub Commit Status API | `runner` | §12、§13、§15、§22.0c、§27.1 | `commit_status_enabled`、context、target_url、pending/success/failure の送信条件、失敗時の扱い、build log 記録が一致する。 |
| ドライラン実行モード | `runner` | §11、§12、§13、§15、§27.2 | `--dry-run` が状態ファイル、log、history、deploy、通知を変更せず、設定・GitHub・SHA 判定結果を固定 JSON で返す。 |
| ビルド失敗時の自動リトライ | `runner` | §12、§13、§15、§22.0c、§27.3 | retry 対象エラー、最大回数、backoff、attempt log、最終 status、SHA 更新禁止条件が一致する。 |
| Webhook 通知失敗リトライキュー | `runner` | §11、§13、§16 | `.notify_pending` の schema、再送順序、失敗時保持が一致する。 |
| 複数ファイル監視 | `runner` / `builder` / `api` | §11、§12、§13、§15、§27.21 | `target_files` の検証、対象別 SHA 差分、build target 決定、履歴・ログ・API 表示が一致する。 |
| ビルドパイプライン YAML 定義 | `runner` / `api` | §12、§13、§15、§22.0e、§27.22 | `.pipeline.yml` の内製 subset parse、step 実行順、timeout、env、失敗時 status、API 保存が一致する。 |
| ローカルファイル監視モード | `runner` | §11、§12、§13、§27.23 | GitHub API を呼ばず、local snapshot の SHA-256 差分だけで変更検出し、trigger と status が一致する。 |
| タグ付きコミットのみビルド | `runner` / `api` | §12、§13、§15、§22.0e、§27.24 | tag pattern、GitHub tags API、skip 条件、build log/history、設定 API が一致する。 |
| 並列マルチターゲットビルド | `runner` | §12、§13、§14a、§15、§27.26 | worker 上限、target 別 status、pending transfer、最終 build status、ログ順序が一致する。 |
| ビルド前後フック | `runner` / `api` | §13、§15、§22.0e、§27.27 | `.hooks` schema、pre/post 実行、abort 条件、hook log、API CRUD が一致する。 |
| リモートビルド対応 | `runner` / `api` | §12、§13、§14a、§15、§27.29 | remote command、archive 取得、manifest 検証、状態記録、失敗時 rollback 不実行が一致する。 |
| ブランチ別環境変数 | `runner` / `api` | §12、§13、§15、§22.0e、§27.31 | branch env schema、許可 key、secret mask、process env 注入、API 保存が一致する。 |
| ビルド通知連携 | `runner` / `api` / `sdk` / `ui` | §13、§15、§16、§22.0e、§27.32 | 通知 event、channel schema、送信順、retry、mask、notify log、API / UI 表示が一致する。 |
| ビルド時間トレンド記録 | `runner` / `api` | §13、§15、§22.0e、§27.33 | `.build_trends.json`、移動平均、中央値、p95、API response、破損時復旧が一致する。 |
| ビルド依存チェーン | `runner` / `api` | §11、§13、§15、§22.0e、§27.34 | `.build_chain_config`、依存 DAG 検証、実行順、skip / failure status、chain log が一致する。 |
| ビルド優先度キュー | `runner` / `api` | §11、§13、§22.0e、§27.35 | queue priority、created_seq、同一優先度 FIFO、API 表示、cancel / clear が一致する。 |
| 失敗原因の自動分類 | `runner` / `api` | §13、§15、§22.0e、§27.36 | failure_category、evidence、分類優先順位、history / log / UI 表示が一致する。 |
| ビルド実行環境の記録 | `runner` | §13、§15、§27.37 | build 開始時の environment snapshot、secret 非含有、log schema、検証 fixture が一致する。 |
| ビルド所要時間の異常検知 | `runner` / `api` | §13、§15、§16、§22.0e、§27.38 | trend 基準、異常判定、WARN、history flag、通知 payload、設定値が一致する。 |

### 0i.3 API / SDK / UI

| 機能 | 責務 component | 詳細仕様節 | 受け入れ条件 |
|------|-------------------|------------|--------------|
| ポーリング間隔の動的変更 | `api` | §22.0e、`docs/details/setup.md` §26、§27.11 | `POST /api/schedule/interval` が systemd timer 設定を更新し、検証コマンドで反映を確認できる。 |
| GitHub Webhook 受信 | `api` / `runner` | §22.0e、§22-W、§13、§27.12 | HMAC 検証、イベント記録、キュー投入またはビルドトリガー、エラー応答が一致する。 |
| 設定バリデーション API | `api` / `sdk` / `ui` | §22.0c、§22.0e、§23、§24、§27.5 | `POST /api/config/validate` が状態を変更せず、正規化後設定、warnings、errors を返す。 |
| API アクセスログ | `api` / `sdk` / `ui` | §22.0a、§22.0c、§22.0e、§23、§24、§27.6 | `.api_access_log` の schema、追記対象、マスク条件、一覧 API、UI 表示が一致する。 |
| Webhook イベントログ | `api` | §11、§22.0e、§22-W、§27.13 | `.webhook_events.json` の JSON Lines schema と一覧 API が一致する。 |
| ヘルスチェックエンドポイント | `api` | §22.0e、§27.16 | `GET /api/health` の稼働秒数、最終ビルド、最終転送、エラー応答が一致する。 |
| Webhook イベント一覧取得 API | `api` / `sdk` / `ui` | §22.0e、§23、§24、§27.13 | `GET /api/webhook-events` の query、response、SDK method、UI 表示が一致する。 |
| ビルドログ重大度フィルター | `api` / `sdk` / `ui` | §22.0e、§23、§24、§27.17 | `level=warn\|error` の query、検索結果、UI filter が一致する。 |
| ブランチ設定の動的変更 API | `runner` / `api` | §11、§12、§22.0e、§27.18 | `.branch_config`、GET/POST API、runner 再起動不要条件が一致する。 |
| 週次ビルドサマリー Webhook | `runner` / `api` | §12、§13、§16、§22.0e、§27.19 | 週次判定、集計対象、通知 payload、手動送信 API が一致する。 |
| 設定変更の詳細 diff 記録 | `api` | §22.0a、§22.0e、§27.20 | `.config_log` の diff 文字列、対象 API、マスク条件が一致する。 |
| ビルド承認フロー | `runner` / `api` / `sdk` / `ui` | §11、§13、§15、§16、§22.0e、§27.30 | `.approval_queue`、承認/却下 API、通知、timeout、UI 操作、履歴 status が一致する。 |

### 0i.4 Archive / Artifact / Security

| 機能 | 責務 component | 詳細仕様節 | 受け入れ条件 |
|------|-------------------|------------|--------------|
| ビルドログのアーカイブ圧縮 | `runner` / `api` / `sdk` / `ui` | §12、§13、§15、§22.0c、§22.0e、§23、§24、§27.7 | gzip 形式、archive 先、参照順、cleanup/archive API、disk usage 集計、UI 表示が一致する。 |
| ビルド所要時間の記録と統計 API | `runner` / `api` | §15、§22.0e、§27.14 | `started_at`、`finished_at`、`duration_seconds` と統計 API が一致する。 |
| ビルドアーティファクト世代管理 | `runner` / `api` | §14b、§22.0e | `.snapshots/` の保持世代、削除、rollback API が一致する。 |
| ビルドアーティファクト管理 | `api` / `ui` / `sdk` | §14b、§22.0e、§23、§24、§27.15 | 一覧、ダウンロード、削除、ロールバックの API、SDK、UI が一致する。 |
| ビルドトリガー専用 API スコープ | `api` / `sdk` / `ui` | §22.0a、§22.0c、§22.0e、§23、§24、§27.42 | `trigger` scope token が build 起動系だけを許可し、その他 API を拒否する。 |
| API キー管理 | `api` / `sdk` / `ui` | §22.0a、§22.0c、§22.0e、§23、§24、§27.43 | API key 本体の一回表示、hash 保存、scope、期限、失効、監査が一致する。 |
| 監査ログ | `api` / `sdk` / `ui` | §22.0a、§22.0c、§22.0e、§23、§24、§27.44 | `.audit_log` schema、対象操作、mask、検索 API、UI 表示が一致する。 |
| セッションタイムアウト変更設定 | `api` / `sdk` / `ui` | §22.0c、§22.0e、§23、§24、§25、§27.45 | `session_timeout_seconds` の範囲、保存、既存 session の扱い、新規 session 期限が一致する。 |
| TOTP 二要素認証 | `api` / `sdk` / `ui` | §22.0a、§22.0c、§22.0e、§23、§24、§25、§27.46 | RFC 6238 TOTP、二段階 login、secret 保存、確認、無効化、UI 操作が一致する。 |
| API レート制限 | `api` / `sdk` / `ui` | §22.0a、§22.0c、§22.0e、§23、§24、§27.47 | IP / actor / endpoint group の固定窓制限、`429`、状態保存、設定 API、UI 表示が一致する。 |

---

## 0. システム概要

Adlaire CI は Go 版コンポーネントと JavaScript / HTML 管理ツールで構成する。

Part 3 詳細仕様セットは、`builder`、`runner`、`api`、`admin`、`sdk`、`ui`、`setup`、`statefile`、`archive`、`commitstatus`、`security`、`fixture` の実装詳細を責務 component 別に定義する。

本ファイルは詳細仕様の入口、索引、共通固定値、対応表、リポジトリ内ソース配置、横断補足契約だけを持つ。各 component の入出力、状態、処理順序、異常系、検証条件の本文は、責務 component 別の `docs/details/*.md` を正とする。

`mcp` は将来計画であり、MCP 専用詳細仕様が新設されるまで、本ファイルおよび責務 component 別詳細仕様ファイルでは入出力、状態、起動手順、検証条件を定義しない。

---

## 0j. リポジトリ内ソース配置

Adlaire CI の標準リポジトリ内ソース配置は以下とする。

本節は、移行後の標準配置を定義する。標準配置への実装移行は完了済みであり、旧配置の `build_spec.go`、`runner.go`、`build_spec_test.go`、`runner_test.go`、`testdata/build_spec/` を現行実装実体として扱わない。

本節の tree は標準配置の最終形を示す。現時点で `将来計画` または `仕様化済み・未実装` の path は、該当 owner component が実装対象になった PR で追加する。標準配置図に含まれていることだけを理由に、未実装ファイル、将来計画ファイル、空ディレクトリ、placeholder を作成してはならない。

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

| パス | component | 役割 |
|------|-----------|------|
| `main.go` | `-` | 起動入口。サブコマンド判定、引数受け取り、責務 component 呼び出しを行う。 |
| `components/builder.go` | `builder` | Markdown / Markdown ディレクトリを静的 Web サイトへ変換する。 |
| `components/runner.go` | `runner` | GitHub polling、変更検出、ビルド起動、履歴、ログ、deploy を実行する。 |
| `components/api.go` | `api` | 管理 API サーバー、認証、状態ファイル操作を提供する。 |
| `admin/index.html` | `ui` | 標準管理ツール UI を提供する。 |
| `admin/adlaire-ci-sdk.js` | `sdk` | 管理 API 通信用 SDK を提供する。 |
| `admin/` | `ui` / `sdk` | 標準管理 UI の静的ファイルを配置する。 |
| `testdata/` | `fixture` | コンポーネント別 fixture を配置する。 |
| `docs/examples/` | `-` | 利用例、設定例、サンプル構成を配置する。 |

`main.go` は 1 ファイルとし、実装詳細を含めない。`components/` 配下は 1 コンポーネント = 1 Go ファイルとし、各ファイルは上表の責務を実装する。

標準配置への移行完了条件は以下に固定する。

| 対象 | 移行完了条件 |
|------|--------------|
| `main.go` | repository root に 1 ファイルだけ存在し、サブコマンド判定、引数受け取り、owner component 呼び出しだけを持つ。Markdown 変換、CI 実行、HTTP handler、状態ファイル操作、archive、commitstatus、MCP の実装詳細を含まない。 |
| `components/*.go` | 実装対象 owner component ごとに 1 Go ファイルだけ存在する。現行 Go 実装は `builder` の `components/builder.go`、`runner` の `components/runner.go`、`api` の `components/api.go` とする。`admin` は `admin/` 配下の静的配布物であり、Go ファイルを持たない。`statefile`、`archive`、`commitstatus` は詳細仕様上の責務境界であり、単独 Go ファイルとして追加する場合は、追加対象 Phase または追加実装 PR で仕様状態と索引を更新してから作成する。 |
| 標準移行前ファイル | `build_spec.go`、`runner.go`、`build_spec_test.go`、`runner_test.go`、`testdata/build_spec/` は、対応する標準配置へ移動済みであり、同じ実装本文または同じ fixture が旧配置に残っていない。 |
| testdata | 実装済みまたは仕様化済み・未実装の owner component ごとに `testdata/<component>/` を使用する。`testdata/mcp/` は MCP 専用詳細仕様が新設されるまで作成しない。 |
| admin | `admin/index.html` と `admin/adlaire-ci-sdk.js` は、それぞれ `ui` と `sdk` の owner 詳細仕様に従う。`admin/style.css` と `admin/app.js` は、`docs/details/admin.md` A1 に定義された任意配布物として扱い、未定義の admin 静的ファイルを追加しない。 |
| 将来計画 | `components/mcp.go` と MCP 用 fixture は、MCP 専用詳細仕様が新設され、`docs/SPEC.md` で `仕様化済み・未実装` へ昇格するまで作成しない。 |

標準配置を変更する PR は、旧配置名と標準配置名の両方が同じ実装実体として併存していないこと、`docs/DOCUMENT_INDEX.md` の Specified Components、`docs/SPEC.md` の実装状態、該当 owner component の詳細仕様、testdata 参照が同じ配置を指すことを確認する。

---

## 1〜26. 分割済み詳細仕様インデックス

本節は、旧 Part 3 本文から責務 component 別詳細仕様ファイルへ移動済みの範囲を示す索引である。各範囲の主本文は下表の移動先を正とし、本ファイルでは本文を再定義しない。

`docs/DETAIL_INDEX.md` §27.38a は横断補足契約として本ファイルに残す。§27.21〜§27.38 または api / sdk / ui / statefile の横断連動を実装する場合は、owner component の分割先詳細仕様ファイルを主本文とし、§27.38a を横断同期確認として同時に確認する。

| 移動済み範囲 | owner component | 主本文 |
|--------------|-----------------|--------|
| §1〜§9、§8a | `builder` | `docs/details/builder.md` §1〜§9、§8a |
| §10〜§20、§15a | `runner` | `docs/details/runner.md` §10〜§20、§15a |
| §21〜§22、§21a | `api` | `docs/details/api.md` §21〜§22、§21a |
| §23 | `sdk` | `docs/details/sdk.md` §23 |
| §24 | `ui` | `docs/details/ui.md` §24 |
| §25 | `api` | `docs/details/api.md` §25 |
| §26 | `setup` | `docs/details/setup.md` §26 |
| §27.1 | `commitstatus` | `docs/details/commitstatus.md` §27.1 |
| §27.2〜§27.3、§27.8〜§27.10、§27.14、§27.19、§27.21〜§27.24、§27.26〜§27.27、§27.29〜§27.38 | `runner` | `docs/details/runner.md` §27 |
| §27.4、§27.25、§27.28 | `builder` | `docs/details/builder.md` §27 |
| §27.5〜§27.6、§27.11〜§27.13、§27.16〜§27.18、§27.20、§27.30 | `api` | `docs/details/api.md` §27 |
| §27.7、§27.15 | `archive` | `docs/details/archive.md` §27 |
| §27.42〜§27.47 | `security` | `docs/details/security.md` §27 |

## 27. 追加仕様化機能 詳細仕様

本節は、追加仕様化機能の詳細仕様参照である。各機能の主本文は、owner component の分割先詳細仕様ファイルを正とする。本節に定義された機能は、§0i と §27.0 を入口として、owner 詳細仕様の主本文と必要な collaborator 詳細仕様の確認項目を組み合わせて実装可否を判定する。

### 27.0 追加仕様化機能 共通実装契約

§27 の主本文は、owner component の分割先詳細仕様ファイルを正とする。親ファイルでは、§27 の実装時に共通して確認する参照順、越境禁止、PR 証跡の入口だけを定義する。

| 確認 | 固定内容 |
|------|----------|
| 実装対象判定 | `docs/SPEC.md` で実装状態と実装可否を確認し、将来計画、実装不可、未仕様化、MCP 専用機能を実装対象にしない。 |
| owner 確定 | §0b と §0i で owner component を 1 件に確定し、主本文は owner の分割先詳細仕様ファイルで確認する。 |
| collaborator 確認 | collaborator がある場合は、§27.1〜§27.47 の参照インデックスに列挙された component の分割先ファイルを schema、呼び出し境界、表示、security、setup、fixture、検証観点として読む。 |
| 補完禁止 | 個別節または分割先詳細仕様に存在しない endpoint、状態ファイル、設定 key、UI 操作、SDK method、外部依存を実装判断で追加しない。追加が必要な場合は owner component の詳細仕様、関連 collaborator 詳細仕様、fixture catalog、必要な対応表を先に更新する。 |
| 状態更新 | 状態ファイル更新は `docs/details/statefile.md` §22.0a、§22.0c を正とし、lock、atomic write、JSON Lines、破損時処理を独自定義しない。 |
| security | secret mask、token、session、scope、audit、rate limit は `docs/details/security.md` を正とし、平文保存・平文表示を行わない。 |
| fixture / PR 証跡 | fixture manifest、expected/effects、assertion、PR 証跡、受け入れゲートは `docs/details/fixture.md` §27-F を正とする。 |
| api / sdk / ui 同期 | API endpoint、SDK method、UI 操作が同一機能に関わる場合は、endpoint は `docs/details/api.md`、SDK method は `docs/details/sdk.md`、UI 操作は `docs/details/ui.md` をそれぞれ正本とし、名称、引数、response、error、表示、成功後再取得、失敗時固定が食い違わないことを確認する。 |

§27 の機能を実装した PR は、対象節、owner 詳細仕様、collaborator 詳細仕様、fixture、secret mask、失敗時副作用、実装対象外を PR 本文に記録する。記録が不足する場合は、実装完了として扱わない。

§27 の PR 分割、dry-run 固定契約、fixture 完了条件の詳細は、owner 詳細仕様と `docs/details/fixture.md` §27-F を正とする。親ファイルに同じ fixture schema、expected/effects、個別機能本文を重複定義しない。

### 27.1〜27.47 追加仕様化機能 参照インデックス

本節は、§27 機能の参照先を一覧化するインデックスである。個別機能の入力、出力、状態、処理順序、異常系、endpoint、SDK method、UI DOM、fixture は下表の「主本文」に記載された owner component 詳細仕様を正とする。親ファイルは、下表に記載された主本文、owner component、collaborator component を置き換えない。

| 節 | 機能 | owner | 主本文 | collaborator | 親ファイル側の扱い |
|----|------|-------|--------|--------------|--------------------|
| §27.1 | GitHub Commit Status API | `commitstatus` | `docs/details/commitstatus.md` §27.1 | `runner`、`statefile` | Commit Status payload と送信順の参照先だけを示す。runner の build 実行、commit SHA 確定、build id 採番、pipeline / deploy / snapshot / history の最終結果確定は `docs/details/runner.md` を正とする。 |
| §27.2 | ドライラン実行モード | `runner` | `docs/details/runner.md` §27.2 | `statefile` | 状態ファイル非更新、ログ、history、deploy、通知の扱いを親ファイルで再定義しない。 |
| §27.3 | ビルド失敗時の自動リトライ | `runner` | `docs/details/runner.md` §27.3 | `statefile` | retry 対象、回数、backoff、SHA 更新禁止条件を親ファイルで再定義しない。 |
| §27.4 | 出力サイトへのビルドメタ埋め込み | `builder` | `docs/details/builder.md` §27.4 | `runner`、`api`、`statefile` | HTML meta、REPORT、API 表示、状態反映の境界だけを確認する。 |
| §27.5 | 設定バリデーション API | `api` | `docs/details/api.md` §27.5 | `sdk`、`ui`、`statefile` | validate の保存禁止、response、SDK/UI 対応を親ファイルで再定義しない。 |
| §27.6 | API アクセスログ | `api` | `docs/details/api.md` §27.6 | `sdk`、`ui`、`statefile` | `.api_access_log` schema と一覧 API の主本文を親ファイルへ複製しない。 |
| §27.7 | ビルドログのアーカイブ圧縮 | `archive` | `docs/details/archive.md` §27.7 | `runner`、`api`、`statefile` | gzip archive、cleanup、参照順の実体処理を親ファイルで再定義しない。 |
| §27.8 | ビルドステータスファイル出力 | `runner` | `docs/details/runner.md` §27.8 | `api`、`statefile` | `.build_status.json` schema と更新タイミングを親ファイルで再定義しない。 |
| §27.9 | ビルドトリガー種別の記録 | `runner` | `docs/details/runner.md` §27.9 | `api`、`sdk`、`ui`、`statefile` | trigger 有効値、判定条件、UI 表示の同期確認だけを扱う。 |
| §27.10 | 設定ファイル起動時整合性チェック | `runner` | `docs/details/runner.md` §27.10 | `statefile` | JSON 破損、退避、初期化、終了コードを親ファイルで再定義しない。 |
| §27.11 | ポーリング間隔の動的変更 | `api` | `docs/details/api.md` §27.11 | `runner`、`statefile` | systemd timer 反映の導入・検証手順は `docs/details/setup.md` §26 を確認する。 |
| §27.12 | GitHub Webhook 受信 | `api` | `docs/details/api.md` §27.12 | `runner`、`statefile` | HMAC、event 記録、queue 投入、エラー応答を親ファイルで再定義しない。 |
| §27.13 | Webhook イベントログ / 一覧取得 API | `api` | `docs/details/api.md` §27.13 | `sdk`、`ui`、`statefile` | webhook event schema、一覧 API、SDK/UI 対応を親ファイルで再定義しない。 |
| §27.14 | ビルド所要時間の記録と統計 API | `runner` | `docs/details/runner.md` §27.14 | `api`、`statefile`、`archive` | duration 計測、統計値、archive 連携の境界だけを確認する。 |
| §27.15 | ビルドアーティファクト管理 | `archive` | `docs/details/archive.md` §27.15 | `api`、`sdk`、`ui`、`runner`、`statefile` | snapshot 作成トリガーは `docs/details/runner.md` §14b を正とする。 |
| §27.16 | ヘルスチェックエンドポイント | `api` | `docs/details/api.md` §27.16 | `statefile` | health response とエラー応答を親ファイルで再定義しない。 |
| §27.17 | ビルドログ重大度フィルター | `api` | `docs/details/api.md` §27.17 | `sdk`、`ui`、`archive`、`statefile` | log query、検索結果、UI filter の同期確認だけを扱う。 |
| §27.18 | ブランチ設定の動的変更 API | `api` | `docs/details/api.md` §27.18 | `runner`、`statefile` | `.branch_config`、GET/POST API、runner 反映条件を親ファイルで再定義しない。 |
| §27.19 | 週次ビルドサマリー Webhook | `runner` | `docs/details/runner.md` §27.19 | `api`、`statefile` | 週次集計、通知 payload、手動送信 API の境界だけを確認する。 |
| §27.20 | 設定変更の詳細 diff 記録 | `api` | `docs/details/api.md` §27.20 | `statefile` | `.config_log` diff 形式と mask 条件を親ファイルで再定義しない。 |
| §27.21 | 複数ファイル監視 | `runner` | `docs/details/runner.md` §27.21 | `builder`、`api`、`statefile` | target_files、SHA 差分、build target、API 表示の同期確認だけを扱う。 |
| §27.22 | ビルドパイプライン YAML 定義 | `runner` | `docs/details/runner.md` §27.22 | `api`、`statefile` | pipeline subset、step 実行順、timeout、env を親ファイルで再定義しない。 |
| §27.23 | ローカルファイル監視モード | `runner` | `docs/details/runner.md` §27.23 | `statefile` | GitHub API 非使用条件と local snapshot 差分検出を親ファイルで再定義しない。 |
| §27.24 | タグ付きコミットのみビルド | `runner` | `docs/details/runner.md` §27.24 | `api`、`statefile` | tag pattern、skip 条件、history/log 反映を親ファイルで再定義しない。 |
| §27.25 | ビルドキャッシュ | `builder` | `docs/details/builder.md` §27.25 | `runner`、`statefile` | cache manifest、再利用条件、無効化条件を親ファイルで再定義しない。 |
| §27.26 | 並列マルチターゲットビルド | `runner` | `docs/details/runner.md` §27.26 | `statefile` | worker 上限、target 別 status、ログ順序を親ファイルで再定義しない。 |
| §27.27 | ビルド前後フック | `runner` | `docs/details/runner.md` §27.27 | `api`、`statefile` | hook schema、pre/post 実行、abort 条件を親ファイルで再定義しない。 |
| §27.28 | 依存ファイルトラッキング | `builder` | `docs/details/builder.md` §27.28 | `runner`、`statefile` | dependency manifest、関連 target 判定、full build 条件を親ファイルで再定義しない。 |
| §27.29 | リモートビルド対応 | `runner` | `docs/details/runner.md` §27.29 | `api`、`archive`、`statefile` | remote command、archive 取得、manifest 検証を親ファイルで再定義しない。 |
| §27.30 | ビルド承認フロー | `api` | `docs/details/api.md` §27.30 | `runner`、`sdk`、`ui`、`statefile` | approval queue、承認/却下 API、通知、UI 操作の同期確認だけを扱う。 |
| §27.31 | ブランチ別環境変数 | `runner` | `docs/details/runner.md` §27.31 | `api`、`statefile` | branch env schema、許可 key、secret mask、process env 注入を親ファイルで再定義しない。 |
| §27.32 | ビルド通知連携 | `runner` | `docs/details/runner.md` §27.32 | `api`、`sdk`、`ui`、`statefile` | 通知 event、channel schema、retry、notify log の境界だけを確認する。 |
| §27.33 | ビルド時間トレンド記録 | `runner` | `docs/details/runner.md` §27.33 | `api`、`statefile` | `.build_trends.json`、移動平均、中央値、p95 を親ファイルで再定義しない。 |
| §27.34 | ビルド依存チェーン | `runner` | `docs/details/runner.md` §27.34 | `api`、`statefile` | DAG 検証、実行順、skip / failure status を親ファイルで再定義しない。 |
| §27.35 | ビルド優先度キュー | `runner` | `docs/details/runner.md` §27.35 | `api`、`statefile` | queue priority、created_seq、FIFO を親ファイルで再定義しない。 |
| §27.36 | 失敗原因の自動分類 | `runner` | `docs/details/runner.md` §27.36 | `api`、`statefile` | failure_category、evidence、分類優先順位を親ファイルで再定義しない。 |
| §27.37 | ビルド実行環境の記録 | `runner` | `docs/details/runner.md` §27.37 | `statefile` | environment snapshot と secret 非含有条件を親ファイルで再定義しない。 |
| §27.38 | ビルド所要時間の異常検知 | `runner` | `docs/details/runner.md` §27.38 | `api`、`statefile` | trend 基準、異常判定、通知 payload、設定値を親ファイルで再定義しない。 |
| §27.42 | ビルドトリガー専用 API スコープ | `security` | `docs/details/security.md` §27.42 | `api`、`sdk`、`ui`、`statefile` | scope 判定、拒否条件、API/SDK/UI 表示を親ファイルで再定義しない。 |
| §27.43 | API キー管理 | `security` | `docs/details/security.md` §27.43 | `api`、`sdk`、`ui`、`statefile` | API key の一回表示、hash 保存、scope、期限、失効、監査を親ファイルで再定義しない。 |
| §27.44 | 監査ログ | `security` | `docs/details/security.md` §27.44 | `api`、`statefile` | `.audit_log` schema、対象操作、mask、検索 API、UI 表示を親ファイルで再定義しない。 |
| §27.45 | セッションタイムアウト変更設定 | `security` | `docs/details/security.md` §27.45 | `api`、`sdk`、`ui`、`statefile` | timeout 範囲、保存、既存 session、新規 session 期限を親ファイルで再定義しない。 |
| §27.46 | TOTP 二要素認証 | `security` | `docs/details/security.md` §27.46 | `api`、`sdk`、`ui`、`statefile` | TOTP、二段階 login、secret 保存、確認、無効化、UI 操作を親ファイルで再定義しない。 |
| §27.47 | API レート制限 | `security` | `docs/details/security.md` §27.47 | `api`、`sdk`、`ui`、`statefile` | 固定窓制限、`429`、状態保存、設定 API、UI 表示を親ファイルで再定義しない。 |

### 27.38a 横断連動・Runner 拡張機能 実装補足契約

本節は、責務 component 別詳細仕様へ分割しない。§27.21〜§27.38 および api / sdk / ui / statefile の横断連動は runner、builder、api、sdk、ui、statefile、archive にまたがるため、実装者は owner 詳細仕様を主本文とし、本節を横断確認として同時に確認する。

本節の正本範囲は、横断確認、同期禁止、横断処理順、成功後再取得、失敗時固定、実装完了時の横断受け入れ観点に限定する。本節は、個別機能の処理本文、入力、出力、状態 schema、fixture schema、endpoint 詳細、SDK method 詳細、UI DOM 詳細を持たない。これらは各 owner / collaborator の分割先詳細仕様ファイルを正とする。

本節と owner component 別詳細仕様ファイルの内容が矛盾する場合は、個別機能の入出力、状態、処理、異常系、endpoint、SDK、UI、fixture は owner component 別詳細仕様ファイルを正とし、横断処理順、API / SDK / UI / statefile 同期、成功後再取得、失敗時固定だけを本節で確認する。§27.38a を理由に、owner 詳細仕様に存在しない endpoint、SDK method、UI 操作、状態ファイル、設定 key、fixture を追加してはならない。

| 確認 | 固定内容 |
|------|----------|
| runner 起点 | §27.21〜§27.38 の多くは runner の build 実行、queue、history、log、notification に影響するため、`docs/details/runner.md` の該当 §27 節を先に確認する。 |
| builder 連携 | cache、dependency、output meta、生成物に関わる場合は `docs/details/builder.md` の該当 §27 節を同時に確認する。 |
| API 連携 | 設定保存、queue、approval、history、stats、snapshot、rollback、hook、notify、search に関わる場合は `docs/details/api.md` の endpoint / state read-write 契約を同時に確認する。 |
| SDK / UI 連携 | API を管理画面から操作する機能は、SDK method は `docs/details/sdk.md` §23、DOM / 表示条件は `docs/details/ui.md` §24 を正本として確認する。 |
| statefile | 状態 schema、lock、atomic write、JSON Lines、破損時処理、保存順は `docs/details/statefile.md` を正とする。 |
| archive | snapshot、artifact、download、delete、rollback、log archive は `docs/details/archive.md` を正とする。 |
| fixture | §27.21〜§27.38 の受け入れ fixture、secret mask、effects、PR 証跡は `docs/details/fixture.md` §27-F を正とする。 |

**api / sdk / ui / statefile 横断連動契約：**

下表は横断確認表であり、新しい API endpoint、SDK method、UI 操作、状態ファイル副作用を定義する表ではない。下表の機能群は、各 owner component 別詳細仕様ファイルに定義済みの API endpoint、SDK method、UI 操作、状態ファイル副作用を同じ実装単位でそろえる。API だけ、SDK だけ、UI だけを先行して仕様外の仮実装にしてはならない。UI が未実装の Phase では、UI 列は fixture の期待操作として固定し、実装完了扱いには含めない。

| 機能群 | API | SDK | UI 操作 | 状態ファイル副作用 | 成功後再取得 | 失敗時固定 |
|--------|-----|-----|---------|--------------------|--------------|------------|
| login / session | `POST /api/login`, `POST /api/login/totp`, `POST /api/logout`, `GET /api/sessions`, `POST /api/sessions/revoke-all` | `login()`, `loginTotp()`, `logout()`, `getSessions()`, `revokeAllSessions()` | login、TOTP 確認、logout、session 一括失効。 | `.admin_credentials`、`.access_log`、`.audit_log`、memory session。token 本体は永続化しない。 | login 成功後は `getDashboard()` → `getStatus()` → `getQueue()`。session 失効後は `getSessions()`。 | `401` は SDK token / ticket / secret field を破棄し、UI は `panel-login` のみ表示する。 |
| build control | `POST /api/build`, `POST /api/build/force`, `POST /api/build/cancel`, `GET /api/build/stream`, `GET /api/queue`, `DELETE /api/queue` | `triggerBuild()`, `buildForce()`, `cancelBuild()`, `streamBuild()`, `getQueue()`, `clearQueue()` | manual build、force build、cancel、stream 開始/停止、queue clear。 | `.build_state`、`.build_lock`、`.build_logs/{id}.json`、queue entry。force build 時だけ SHA cache 更新。 | `getStatus()` → `getQueue()`、stream end 後は `getLogs()` も実行。 | running は `409`、queue full は `429`、maintenance / circuit は `503` または仕様済み `409`。UI は同じ build request を自動再送しない。 |
| logs / history | `GET /api/logs`, `GET /api/logs/search`, `GET /api/history`, `GET /api/history/{id}/log`, `POST /api/history/{id}/comment`, `POST /api/history/{id}/flag`, `POST /api/history/{id}/tags` | `getLogs()`, `searchLogs()`, `getHistory()`, `getHistoryLog()`, `setHistoryComment()`, `setHistoryFlag()`, `setHistoryTags()` | log 表示、検索、history 表示、comment / flag / tag 保存。 | read-only GET は副作用なし。comment / flag / tag は `.build_logs/{id}.json` と、owner 詳細仕様で履歴・設定ログ更新が定義された場合に限り `.build_history`、`.config_log`。 | 変更系は `getHistory()`、comment は `getHistoryComment(id)` も実行。 | `404` は選択解除または not found 表示。`422 details` は field error。壊れた JSON Lines は response に含めない。 |
| config / repo / branch / schedule | `GET/POST /api/config`, `POST /api/config/validate`, `GET/POST /api/repo-config`, `GET/POST /api/branch-config`, `GET /api/schedule`, `POST /api/schedule/*` | `getConfig()`, `setConfig()`, `validateConfig()`, `setRepoConfig()`, `getBranchConfig()`, `setBranchConfig()`, `getSchedule()`, schedule 系 method | config 保存、validate、repo 保存、branch 保存、schedule 変更。 | `.server_config`、`.repo_config`、`.branch_config`、`.config_log`。validate は保存なし。systemd 反映失敗時も保存済み値は戻さない。 | 保存系は対象 GET → `getConfigLog()`。validate は再取得なし。 | `422` は書込前停止。systemd 失敗 `500` は保存済み値を再取得して表示する。no-op は状態ファイルと log を変更しない。 |
| notify / SMTP / webhook | `GET/POST /api/notify-config`, `POST /api/notify-test`, `GET /api/notify-log`, `GET/POST /api/smtp-config`, `POST /api/smtp-test`, `GET/POST /api/webhook-config`, `GET /api/webhook-events`, `POST /api/notify/weekly-summary` | `getNotifyConfig()`, `setNotifyConfig()`, `notifyTest()`, `getNotifyLog()`, `getSmtpConfig()`, `setSmtpConfig()`, `smtpTest()`, `getWebhookConfig()`, `setWebhookConfig()`, `getWebhookEvents()`, `notifyWeeklySummary()` | notify 保存、test、SMTP 保存/test、webhook secret 保存、event 表示、weekly summary。 | `.notify_config`、`.smtp_config`、`.smtp_secret`、`.webhook_secret`、`.notify_log`、`.webhook_events.json`、`.config_log`。 | 保存系は対象 GET → `getConfigLog()`。test / summary は `getNotifyLog()`。 | secret は response / log / fixture へ平文出力しない。保存成功・失敗とも UI secret field を消去する。 |
| snapshots / rollback / maintenance | `GET /api/snapshots`, `GET /api/snapshots/{id}/download`, `DELETE /api/snapshots/{id}`, `POST /api/history/{id}/rollback`, `GET /api/maintenance`, `POST /api/maintenance/enable`, `POST /api/maintenance/disable` | `getSnapshots()`, `downloadSnapshot()`, `deleteSnapshot()`, `rollbackHistory()`, `getMaintenance()`, `enableMaintenance()`, `disableMaintenance()` | snapshot list/download/delete、rollback、maintenance enable/disable。 | `.snapshots/`、`.build_history`、`.build_logs/{new_id}.json`、`.maintenance`、`.config_log`。download は副作用なし。 | delete は `getSnapshots()`。rollback は `getHistory()` → `getStatus()`。maintenance は `getMaintenance()`。 | delete / rollback は確認必須。running rollback は `409`。maintenance enabled 中は build / rollback / 設定変更系を disabled。 |
| access / hooks / rules / pipeline / notes / layout | access、hooks、alert rules、tag rules、pipeline config、notes、dashboard layout の GET/POST/DELETE endpoint | 対応する §23 SDK method | 保存、追加、削除、notes 保存、dashboard layout 保存。 | `.access_control`、`.hooks`、`.alert_rules`、`.tag_rules`、`.pipeline_config`、`.notes`、`.dashboard_layout`、`.config_log`。 | 対象 GET → 変更系で config log 対象の場合は `getConfigLog()`。layout は `getDashboardLayout()` → `getDashboard()`。 | duplicate `409` は競合表示。validation `422` は field error。削除対象不在は `404`。 |
| tokens / audit / access logs / rate limit | `GET /api/tokens`, `POST /api/tokens`, `DELETE /api/tokens/{id}`, `GET /api/audit-log`, `GET /api/access-log`, `GET /api/api-access-log`, `GET/POST /api/api-rate-limit` | `getTokens()`, `createToken()`, `revokeToken()`, `getAuditLog()`, `getAccessLog()`, `getApiAccessLog()`, `getApiRateLimit()`, `setApiRateLimit()` | token 発行/失効、audit / access log 表示、rate limit 保存。 | `.api_tokens`、`.audit_log`、`.access_log`、`.api_access_log`、`.api_rate_state`、`.server_config`、`.config_log`。token 本体は作成時 response のみ。 | token 操作は `getTokens()` → `getAuditLog()`。rate limit は `getApiRateLimit()`。 | token 本体は再取得不可。`403` は logout しない。`429` は rate limit 表示し、同一操作を自動 retry しない。 |

**横断処理順契約：**

| 処理種別 | 固定順序 |
|----------|----------|
| 認証必須 JSON API | method / path 判定 → body 禁止判定 → JSON parse → 認証 / scope → rate limit → endpoint 固有 validation → read → write 計画 → atomic write → JSON Lines 追記 → response。 |
| read-only API | method / path 判定 → body 禁止判定 → 認証 / scope → query validation → read → 壊れた任意行除外 → response。read-only API は状態ファイルを書き換えない。 |
| UI 変更操作 | panel error / success 消去 → UI 入力検証 → 対象操作 disabled → SDK 呼び出し → 成功後再取得 → success 表示 → secret 消去 → disabled 再評価。 |
| UI 取得操作 | panel error 消去 → 対象操作 disabled → SDK 呼び出し → DOM 更新 → empty state 判定 → disabled 再評価。success 表示は行わない。 |
| SDK request | 引数検証 → path / query / body 生成 → Authorization 付与 → timeout 設定 → fetch → status 判定 → response parse → token 変化適用 → return / throw。 |
| multi-file write | 全入力検証 → 全対象 read → 全 write payload 生成 → `docs/details/api.md` §22.0d の Write 順に atomic write → JSON Lines 追記 → response。途中失敗時は未処理ファイルを書かない。 |

§27.21〜§27.38 の実装では、owner 詳細仕様にない状態ファイル、endpoint、SDK method、UI 操作、外部公開構成を追加してはならない。追加が必要な場合は、親ファイルではなく、該当 owner / collaborator の分割先詳細仕様を先に改訂する。
