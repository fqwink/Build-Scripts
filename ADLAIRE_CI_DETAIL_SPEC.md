# Adlaire CI — 詳細仕様

本ファイルは `ADLAIRE_CI_SPEC.md` の Part 3 詳細仕様であり、実装の具体的詳細に関する正本である。

本ファイルに、方針、ポリシー、実装状態、正本関係、ロードマップ状態、実装可否の上位判断を記載してはならない。これらは `ADLAIRE_CI_SPEC.md` を正とする。

---

# Part 3 — 仕様
> 実装の具体的詳細を定める。「どのように動作・実装するか」に答える。

---

## 詳細仕様の読み方

本ファイルは、実装者が実装時に参照する詳細仕様だけを扱う。方針、ポリシー、成熟度定義、ロードマップ状態、実装可否、PR 分割判断は `ADLAIRE_CI_SPEC.md` を正とし、本ファイルで再定義しない。

実装者は、対象機能ごとに以下の順で読む。

1. `ADLAIRE_CI_SPEC.md` の実装状態、Part 1 §12、§13 で、対象が実装対象であることを確認する。
2. 本ファイル §0i で、対象機能に対応する詳細仕様節と受け入れ条件を特定する。
3. 本ファイル §0a〜§0h で、詳細仕様の記載基準、共通固定値、実装前確認項目、検証条件、Phase 順序を確認する。
4. 責務 component の詳細節を読み、owner component、collaborator component、入力、出力、状態、正常系、異常系、セキュリティ、検証条件を確認する。
5. `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 のセットアップ・アップデート手順と `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.7 の受け入れ条件に影響がある場合は、実装 PR の検証対象に含める。

詳細仕様節に §0h の必須項目が不足している場合は、実装判断で補完してはならない。先に本ファイルを改訂し、`ADLAIRE_CI_SPEC.md` の対象範囲と整合させる。

| 範囲 | 役割 |
|------|------|
| §0〜§0j | 詳細仕様の記載基準、実装前確認項目、共通固定値、検証、Phase、詳細節対応表、リポジトリ内ソース配置 |
| §1〜§9 | `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md`。`components/builder.go` / `adlaire-ci-build` の詳細仕様 |
| §10〜§20 | `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md`。`components/runner.go` / `adlaire-ci-runner` の詳細仕様 |
| §21〜§22 | `ADLAIRE_CI_DETAIL_API_SPEC.md`。`components/api.go` / `adlaire-ci-api` の詳細仕様 |
| §23 | `ADLAIRE_CI_DETAIL_SDK_SPEC.md`。`admin/adlaire-ci-sdk.js` の詳細仕様 |
| §24 | `ADLAIRE_CI_DETAIL_UI_SPEC.md`。`admin/index.html` の詳細仕様 |
| §25 | `ADLAIRE_CI_DETAIL_API_SPEC.md`。認証の実装仕様 |
| §26 | `ADLAIRE_CI_DETAIL_SETUP_SPEC.md`。バイナリ配布前提のセットアップ、アップデート、受け入れ条件 |

本ファイルを分割する場合は、§0b.1 の責務 component 分割仕様に従う。分割後も `ADLAIRE_CI_DETAIL_SPEC.md` は入口、索引、共通固定値、責務 component 対応表を持つ。分割先ファイルは、それぞれの owner component と collaborator component の詳細仕様だけを持つ。

## 0a. 詳細仕様の記載基準

本ファイルの仕様項目は、実装者が追加の設計判断や推測を行わずに実装できる粒度で記載する。

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

未確定の内容は、実装可能な詳細仕様として記載してはならない。未確定の場合は、本ファイルへ推測で具体値を記載せず、`ADLAIRE_CI_SPEC.md` で状態を確認する。

対象範囲の内容は、実装ファイルが存在しなくても、本節の基準に従って実装可能な粒度まで具体化する。

---

## 0b. 詳細仕様参照表

本節は、責務 component ごとに参照する詳細仕様節を示す。実装状態、実装可否、ロードマップ状態は `ADLAIRE_CI_SPEC.md` を確認する。

| 責務 component | 詳細仕様節 | 主な確認対象 |
|--------------------|------------|--------------|
| `components/builder.go` | `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §1〜§9、§8a、§27.4、§27.25、§27.28 | CLI、入力 Markdown、出力サイト、HTML / CSS / JavaScript、変換 report、fixture、builder owner 追加機能。 |
| `components/runner.go` | `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §10〜§20、§15a、§27.1〜§27.3、§27.8〜§27.10、§27.14、§27.19、§27.21〜§27.24、§27.26〜§27.27、§27.29、§27.31〜§27.38 | 設定、状態ファイル、GitHub API、pipeline、転送、snapshot、通知、systemd、fixture、runner owner 追加機能。 |
| `components/api.go` | `ADLAIRE_CI_DETAIL_API_SPEC.md` §21〜§22、§25、§27.5〜§27.6、§27.11〜§27.13、§27.16〜§27.18、§27.20、§27.30、§27.42〜§27.47、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a、§22.0c | API 共通処理、endpoint、状態ファイル read/write、認証、session、API owner 追加機能。 |
| `admin/adlaire-ci-sdk.js` | `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 | SDK class、method、HTTP 対応、error、stream、token 破棄。 |
| `admin/index.html` | `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 | 画面構成、DOM id、panel、SDK 呼び出し、表示状態、秘密情報消去。 |
| `setup` | `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 | バイナリ配布、配置、systemd、セットアップ、アップデート、リリース成果物検証。 |
| `components/statefile.go` | `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a、§22.0c | 状態ファイル共通仕様、lock、atomic write、JSON Lines、破損時処理、状態読取 adapter、主要 schema。 |
| `components/archive.go` | `ADLAIRE_CI_DETAIL_ARCHIVE_SPEC.md` §27.7、§27.15 | build log archive、snapshot、download、delete、rollback、cleanup。 |
| `components/mcp.go` | 詳細仕様なし | 本ファイルでは実装可能な入出力、状態、起動手順、検証条件を定義しない。 |

---

## 0b.1 責務 component 別 詳細仕様ファイル分割仕様

本節は、`ADLAIRE_CI_DETAIL_SPEC.md` を責務 component 別に分割する場合の固定仕様である。分割は、仕様内容の移動と参照先更新だけを対象とし、機能追加、実装状態変更、実装可否変更、ロードマップ変更、方針・ポリシー追加を含めてはならない。

分割後の詳細仕様ファイルは以下に固定する。`COMMON`、`CORE`、`BASE`、`SHARED`、`FOUNDATION`、その他の横断共通基盤ファイルは作成しない。

| ファイル | 持つ内容 | 持たない内容 |
|----------|----------|--------------|
| `ADLAIRE_CI_DETAIL_SPEC.md` | 詳細仕様の入口、読み方、共通固定値、実装前確認項目、検証マトリクス、Phase、詳細節対応表、リポジトリ内ソース配置、責務 component 分割仕様。 | 各 component の詳細な処理本文、fixture 詳細、状態ファイル schema 詳細、個別 endpoint 詳細、個別 UI 操作詳細。 |
| `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` | `builder` owner の Markdown 変換、静的 Web サイト出力、HTML / CSS / JavaScript、theme component、builder fixture。 | runner / API / SDK / UI の実行責務。 |
| `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` | `runner` owner の GitHub 監視、状態ファイル更新、pipeline、deploy、snapshot、通知、runner fixture。 | API endpoint の認証・応答本文、SDK method、UI DOM 詳細。 |
| `ADLAIRE_CI_DETAIL_API_SPEC.md` | `api` owner の HTTP 共通契約、endpoint、状態ファイル read/write、認証連携、API fixture。 | SDK 内部実装、UI DOM 詳細、runner の build 実行責務。 |
| `ADLAIRE_CI_DETAIL_SDK_SPEC.md` | `sdk` owner の SDK class、method、HTTP 対応、error、stream、token 破棄。 | API endpoint の状態ファイル更新責務、UI DOM 詳細。 |
| `ADLAIRE_CI_DETAIL_UI_SPEC.md` | `ui` owner の DOM id、panel、操作、表示状態、SDK 呼び出し、秘密情報消去。 | SDK method 実装、API endpoint 実装、状態ファイル直接操作。 |
| `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` | `setup` owner のバイナリ配布、配置、systemd、セットアップ、アップデート、リリース成果物検証。 | runner / API / SDK / UI の個別機能本文。 |
| `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` | `statefile` owner の状態ファイル共通仕様、lock、atomic write、JSON Lines、破損時処理、状態読取 adapter、主要 schema。 | API endpoint の request / response、runner の業務処理、UI 表示判断。 |
| `ADLAIRE_CI_DETAIL_ARCHIVE_SPEC.md` | `archive` owner の build log archive、snapshot、download、delete、rollback、cleanup。 | runner の build 実行、API 共通 request / response、SDK method 実装、UI DOM 詳細。 |
| `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` | fixture manifest、assertion、fake、testdata、受け入れ fixture 共通契約、PR 証跡テンプレート。 | 個別 component の通常処理本文。 |

`ADLAIRE_CI_DETAIL_SPEC.md` §27.38a は、runner、builder、API、SDK、UI にまたがる横断補足契約であり、責務 component 別の分割先へ移動しない。§27.21〜§27.38 を実装する場合は、owner component の分割先詳細仕様ファイルと §27.38a を同時に満たす。

分割時の移動単位は、owner component を第一基準とする。複数 component が関わる機能は、owner component のファイルに主本文を置き、collaborator component のファイルには参照リンク、禁止事項、受け入れ観点だけを置く。主本文を複数ファイルへ重複定義してはならない。

分割先ファイルへ移動する各節は、移動後も以下を満たす。

| 項目 | 必須条件 |
|------|----------|
| 節番号 | 既存の節番号を維持する。番号の再採番は行わない。 |
| 参照 | 旧参照先と新参照先が一意に追跡できるよう、`ADLAIRE_CI_DETAIL_SPEC.md` の対応表を更新する。 |
| owner | 各機能節に owner component を 1 件だけ明記する。 |
| collaborator | collaborator component は 0 件以上を明記し、owner component を含めない。 |
| 重複禁止 | 同じ入力、出力、状態 schema、HTTP body、DOM id、fixture assertion を複数ファイルで重複定義しない。 |
| 横断事項 | 横断する固定値は `ADLAIRE_CI_DETAIL_SPEC.md` に置く。横断共通基盤を component として扱わない。 |
| 索引 | `DOCUMENT_INDEX.md` に、分割後ファイルの役割と正本範囲を反映する。 |

分割後に実装者が詳細仕様を読む順序は以下に固定する。

1. `ADLAIRE_CI_SPEC.md` で実装対象、実装状態、実装可否を確認する。
2. `ADLAIRE_CI_DETAIL_SPEC.md` で共通固定値、責務 component、詳細節対応表を確認する。
3. owner component の分割先詳細仕様ファイルを読む。
4. collaborator component がある場合は、該当する分割先詳細仕様ファイルの参照節を読む。
5. 状態ファイルの読み書き、lock、atomic write、schema を扱う場合は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` を読む。
6. fixture、fake、PR 証跡が必要な場合は `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` を読む。

分割作業は、以下の完了条件をすべて満たすまで完了扱いにしてはならない。

| 完了条件 | 判定 |
|----------|------|
| 旧ファイル内の移動対象本文が対応する分割先に移動している。 | 必須 |
| `ADLAIRE_CI_DETAIL_SPEC.md` には入口、索引、共通固定値、対応表、分割仕様だけが残っている。 | 必須 |
| `ADLAIRE_CI_SPEC.md`、`DOCUMENT_INDEX.md`、各分割先ファイル間の参照が矛盾していない。 | 必須 |
| `rg` で旧節名、旧ファイル名、移動前参照の取り残しを確認している。 | 必須 |
| 実装ファイル、fixture、testdata の内容を分割作業だけで変更していない。 | 必須 |

---

## 0c. 実装前確認項目

実装者は、対象機能について以下の条件をすべて満たすまで実装を開始してはならない。

| ゲート | 合格条件 |
|--------|----------|
| 対応表 | 対象機能が §0i の詳細節対応表に記載され、詳細仕様節と受け入れ条件が一意に示されている。 |
| テンプレート | 対象機能の詳細仕様が §0h の機能仕様テンプレートの必須項目を満たしている。 |
| 責務境界 | owner component、collaborator component、対象ファイル、呼び出し元、呼び出し先、変更してよい状態ファイルが明記されている。 |
| 入出力 | すべての入力、出力、既定値、許容値、必須/任意、型、文字コード、時刻形式が明記されている。 |
| 状態管理 | 状態ファイルのパス、JSON 形式、更新タイミング、初期状態、破損時の扱い、権限が明記されている。 |
| 正常系 | 処理順序、分岐条件、ループ条件、成功条件、終了条件が明記されている。 |
| 異常系 | エラー条件、ログレベル、HTTP ステータス、戻り値、再試行有無、処理継続/中断条件が明記されている。 |
| 冪等性 | 同一リクエスト、再実行、途中失敗後の再開で二重実行・二重削除・状態破壊が発生しない条件が明記されている。 |
| 排他制御 | 同時実行、ロック、タイムアウト、ロック残存時の扱いが明記されている。 |
| セキュリティ | 秘密情報の保存禁止、マスク、ファイル権限、認証/認可、外部公開可否が明記されている。 |
| 検証 | 構文確認、単体確認、手動 API 確認、生成物確認、ログ確認、失敗系確認のいずれを行うかが明記されている。 |

上記ゲートのいずれかが未充足の場合、実装判断で補完してはならない。先に本ファイルまたは `ADLAIRE_CI_SPEC.md` を改訂し、未充足項目を仕様として確定する。

実装後の完了条件は以下とする。

1. 実装した機能が、本ファイルに記載された入力、出力、状態、異常系、検証条件と一致する。
2. 対象機能が §0h の機能仕様テンプレートを満たし、§0i の詳細節対応表の受け入れ条件を満たしている。
3. 実装対象外に残す機能が PR 本文に明記されている。
4. `ADLAIRE_CI_SPEC.md`、`ADLAIRE_CI_DETAIL_SPEC.md`、`DOCUMENT_INDEX.md`、`AGENTS.md` のファイル名参照が矛盾していない。
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
| 時刻 | 状態ファイル、API、ログの機械処理用時刻は UTC の ISO 8601 形式（例: `2026-09-16T09:00:00Z`）で保存する。UI 表示のみローカル時刻へ変換してよい。 |
| JSON | JSON object の未知キーは保存しない。読み込み時に未知キーを見つけた場合は無視し、次回保存時に除去する。 |
| atomic write | 状態ファイル更新手順は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a を正とする。各 component は同節の手順を使用し、独自更新手順を持たない。 |
| 権限 | 状態ファイル、秘密情報ファイル、ディレクトリの権限は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a を正とする。 |
| ロック | 状態ファイル lock の作成、待機、解除、競合時応答は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a を正とする。 |
| ログ秘密情報 | PAT、Webhook Secret、SMTP password、session token、API token は stdout、stderr、JSON log、API response、UI 表示へ平文出力しない。表示が必要な場合は `"***"` とする。 |
| 終了コード | CLI / runner は `0` 成功、`1` 一般エラー、`2` 入力・設定エラー、`3` 外部サービス・ネットワークエラー、`4` ロック競合を標準とする。個別節に明記がある場合もこの意味から外してはならない。 |
| 禁止事項 | 仕様にない環境変数、状態ファイル、HTTP endpoint、CLI option、外部依存を実装者判断で追加してはならない。必要な場合は先に本仕様を改訂する。 |

---

## 0e. 完全実装検証マトリクス

対象項目を完了扱いにする場合は、責務 component ごとに下表の検証を満たす。実装ファイルが存在しても、本表の必須検証が未完了の場合は完了扱いにしない。

| 対象 | 必須検証 | 合格条件 |
|------|----------|----------|
| `components/builder.go` | CLI 正常系 | `adlaire-ci-build --src <valid.md-or-dir> --out <site-dir>` が終了コード `0` で終了し、静的 Web サイトと `[REPORT]` を生成する。 |
| `components/builder.go` | CLI 異常系 | 入力不存在、UTF-8 不正、未知引数、出力不可ディレクトリ、未知 theme で §2・§8 の終了コードと stderr が一致する。 |
| `components/builder.go` | Markdown 変換 | 見出し、重複 slug、内部リンク警告、脚注、表、引用、リスト、コードフェンス、未閉鎖フェンス、HTML escape が §4 の出力構造と一致する。 |
| `components/builder.go` | 生成物 | 出力サイトディレクトリに `index.html`、ページ HTML、`assets/style.css`、`assets/app.js`、`assets/search-index.json` が生成され、§5〜§7 の ID / class / JS 機能を含む。 |
| `components/runner.go` | 設定検証 | `--state-dir`、`BRANCH_TARGETS`、必須ファイル不足、未知設定キーで §12 のログ・終了コード・採用優先順位が一致する。 |
| `components/runner.go` | 状態更新 | 成功、ビルド失敗、GitHub API 失敗、転送失敗、lock 競合、JSON 破損で §13 と §22.0a の更新順序・未更新条件が一致する。 |
| `components/runner.go` | 冪等性 | 同一 SHA 再実行、pending retry 再実行、通知 pending 再実行、stale lock 復旧で二重履歴・二重 snapshot・状態破壊が発生しない。 |
| `components/api.go` | API 共通 | 未知 path、未対応 method、body 禁止、JSON 不正、body 上限、認証なし、権限不足、入力検証失敗、ロック競合が §22.0 の status と body を返す。 |
| `components/api.go` | 状態ファイル | 全 write API が §22.0a / §22.0d の対象ファイルだけを atomic write し、秘密情報を平文出力しない。 |
| `components/api.go` | endpoint 契約 | §22.0e の全 endpoint について Request、Response、Success、Errors、Read、Write、SDK、UI の対応が実装と一致する。 |
| `admin/adlaire-ci-sdk.js` | SDK 契約 | 全 method が §22.0e の endpoint のみを呼び、body なし endpoint に body を送らず、HTTP error を `AdlaireCIError` として返す。 |
| `admin/index.html` | UI 契約 | 全操作が §24 の SDK method 経由で動作し、成功表示、失敗表示、disabled、再取得、秘密情報消去が一致する。 |
| セットアップ | systemd | `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 の unit 名、`ExecStart`、配置パス、権限、起動確認コマンドが実際の導入手順と一致する。 |

検証結果は、実装 PR の本文または実装完了報告に、対象、実行コマンド、期待結果、実結果を対応付けて記録する。検証不能な項目がある場合は、その項目を完了扱いにしてはならない。

---

## 0f. 仕様策定完了チェック

本節は、Go 版初期実装へ進む前に仕様策定が完了しているかを判定するチェックである。実装者は、責務 component ごとに下表の必須条件を満たすまで実装を開始してはならない。

| 対象 | 実装着手条件 | 実装禁止条件 | 完了判定 |
|------|--------------|--------------|----------|
| `components/builder.go` | §2〜§8 に CLI option、入力 Markdown、出力サイトディレクトリ、終了コード、stderr、HTML 構造、テーマコンポーネント、JS/CSS、生成物確認が定義されている。 | §4〜§7 にない Markdown 記法、CSS class、JavaScript 機能、外部 asset、theme を追加すること。 | §0e の `components/builder.go` 必須検証をすべて満たし、生成サイトが §5〜§7 と一致する。 |
| `components/runner.go` | §10〜§20 に設定値、状態ファイル、GitHub API、SHA 比較、pipeline 起動、SSH 転送、snapshot、通知、ログ、systemd が定義されている。 | 未定義の環境変数、状態ファイル、queue 挙動、通知チャンネル、pipeline 形式を追加すること。 | §0e の `components/runner.go` 必須検証をすべて満たし、状態ファイル更新順序が §13、§22.0a、§22.0d と一致する。 |
| `components/api.go` | `ADLAIRE_CI_DETAIL_API_SPEC.md` §21〜§22、§25 と `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 に API 共通契約、endpoint、状態ファイル schema、認証、認可、systemd、セットアップが定義されている。 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e にない endpoint、method、status code、response body、状態ファイル write を追加すること。 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e の全 endpoint が Request、Response、Errors、Read、Write、SDK、UI の対応表と一致する。 |
| `admin/adlaire-ci-sdk.js` | `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 に SDK class、method、引数、戻り値、HTTP endpoint 対応、error object、token 破棄条件が定義されている。 | SDK が `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e にない endpoint を呼ぶこと、body 禁止 endpoint に body を送ること、独自 error 形式を返すこと。 | 全 method が `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e と `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 の対応どおりに動作し、HTTP error を `AdlaireCIError` として扱う。 |
| `admin/index.html` | `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 に画面構成、panel、操作、成功表示、失敗表示、disabled、再取得、秘密情報消去が定義されている。 | SDK を介さず API を直接呼ぶこと、未定義の画面・操作・保存先を追加すること、秘密情報を DOM に残すこと。 | 全 UI 操作が `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 の表示条件と `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 の SDK method を満たし、秘密情報 field が指定条件で消去される。 |

上表の対象外である `components/mcp.go`、MCP tools、MCP resources、MCP prompts、HTTP SSE transport、MCP audit / stats / config CRUD は、初期実装では実装しない。これらは、本ファイル内に入出力、状態、起動手順、検証条件を定義しない。

仕様策定完了チェックで未充足が見つかった場合は、実装を開始せず、以下の順で仕様を補完する。

1. 未充足項目が本ファイルの記載対象外である場合は、先に `ADLAIRE_CI_SPEC.md` を確認する。
2. 未充足項目が入出力、状態ファイル、API、SDK、UI、処理順序、異常系、検証条件に関わる場合は、本ファイルの該当節を改訂する。
3. ファイル名、正本関係、対象範囲が変わる場合は、`DOCUMENT_INDEX.md` の更新要否を確認する。
4. 対象項目の詳細節または受け入れ条件が変わる場合は、§0i の詳細節対応表を更新する。
5. 補完後、§0b、§0c、§0e、本節、§0g、§0h、§0i、§0j の条件を再確認する。

---

## 0g. 初期実装 Phase 分割

Go 版初期実装は、`ADLAIRE_CI_SPEC.md` §0e の対象範囲を一括実装せず、下表の Phase 順に進める。上位 Phase の完了判定を満たす前に、下位 Phase の実装 PR を開始してはならない。

| Phase | 対象 | 実装範囲 | 依存条件 | 完了条件 |
|-------|------|----------|----------|----------|
| Phase 1 | `components/builder.go` | §2〜§9 の CLI、Markdown 変換、静的 Web サイト出力、テーマコンポーネント、生成物確認。 | なし。 | §0e の `components/builder.go` 必須検証と §0f の `components/builder.go` 完了判定を満たす。 |
| Phase 2 | `components/runner.go` | §10〜§20 の CI ランナー、GitHub API 連携、SHA キャッシュ、pipeline 起動、SSH 転送、snapshot、通知、ログ、systemd。 | Phase 1 が完了し、`adlaire-ci-build` の CLI 契約が固定されている。 | §0e の `components/runner.go` 必須検証と §0f の `components/runner.go` 完了判定を満たす。 |
| Phase 3 | `components/api.go` P0 / P1 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §21〜§22、§25 と `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 のうち、認証、セッション、共通エラー、状態ファイル読み書き、ビルド操作、status、logs、history、queue、circuit breaker。 | Phase 2 が完了し、runner が書き込む状態ファイル schema が固定されている。 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0f P0 / P1 の必須検証、§0e の `components/api.go` API 共通・状態ファイル検証、§0f の `components/api.go` 完了判定の該当範囲を満たす。 |
| Phase 4 | `components/api.go` P2〜P5 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0f P2〜P5 の config、repo、branch、schedule、notify、snapshot、rollback、maintenance、access control、hooks、tokens 等。 | Phase 3 が完了し、API 共通処理と認証が固定されている。 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0f P2〜P5 の必須検証と §0e の `components/api.go` endpoint 契約を満たす。 |
| Phase 5 | `admin/adlaire-ci-sdk.js` | `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 の SDK class、method、戻り値、HTTP error、token 破棄、query 生成。 | Phase 3 と Phase 4 が完了し、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e の endpoint 契約が固定されている。 | §0e の SDK 契約と §0f の `admin/adlaire-ci-sdk.js` 完了判定を満たす。 |
| Phase 6 | `admin/index.html` | `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 の標準管理ツール UI、panel、操作、成功表示、失敗表示、disabled、再取得、秘密情報消去。 | Phase 5 が完了し、SDK method 契約が固定されている。 | §0e の UI 契約と §0f の `admin/index.html` 完了判定を満たす。 |

### 0g.1 Phase 1 完全仕様ゲート（`components/builder.go`）

Phase 1 は、Markdown 入力から静的 Web サイト出力までを `adlaire-ci-build` 単体で完結させる。実装者は、本 Phase で CI runner、管理 API、SDK、UI、MCP、GitHub API、SSH 転送を実装してはならない。

| 項目 | 固定仕様 |
|------|----------|
| 実装開始条件 | §2〜§8、§8a、§0e、§0f、§0h、§0i、§0j を確認済みである。 |
| 入力 | `--src` で指定された UTF-8 Markdown ファイルまたは Markdown ディレクトリ、`--out`、`--title`、`--theme`、`--base-dir`、`--strict`。 |
| 出力 | `index.html`、必要なページ HTML、`assets/style.css`、`assets/app.js`、`assets/search-index.json`、stdout の進捗行、`[WARN]`、`[REPORT]`。 |
| 状態 | 実行中のメモリ状態だけを使用する。状態ファイル、cache、lock、network、secret は使用しない。 |
| 正常系 | §8 の 12 手順どおりに、入力収集、見出し収集、変換、検索 index、site assembly、atomic output、report 出力を行う。 |
| 異常系 | CLI 不正、入力不存在、UTF-8 不正、theme 不正、出力失敗、未閉鎖 fence、broken link strict failure を §2、§4、§8、§8a の終了コードと stderr に一致させる。 |
| セキュリティ | 生 HTML は pass-through せず escape する。外部 asset、外部 font、CDN、外部 JavaScript を読み込まない。 |
| 実装対象外 | GitHub API、runner 状態ファイル、deploy、snapshot、通知、API server、SDK、admin UI、MCP。 |
| 必須 fixture | §8a Fixture A〜D。 |
| 完了条件 | §0e の `components/builder.go` 必須検証、§0f の完了判定、§8a の全 fixture、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.7 の build script 対象を満たす。 |

Phase 1 完了時は、次 Phase へ引き継ぐ CLI 契約として、`adlaire-ci-build` の終了コード、stdout 進捗行、`[WARN]` 行、`[REPORT]` 行、出力ディレクトリ構造を固定する。Phase 2 以降は、この契約を変更してはならない。変更が必要な場合は Phase 1 仕様改訂に戻る。

#### 0g.1.1 Phase 1 実装順序

Phase 1 は、以下の順序で実装する。順序を入れ替える場合は、入れ替え理由と影響が §8 の処理順序、§8a fixture、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.7 build script 受け入れ条件に影響しないことを実装 PR 本文に記録する。

| 順序 | 実装単位 | 完了判定 |
|------|----------|----------|
| 1 | CLI 引数解析、既定値、未知引数、必須値検証、終了コードを実装する。 | §2 の CLI 契約どおりに成功 / 失敗し、未知引数と空 `--title` が終了コード `2` になる。 |
| 2 | 入力収集、UTF-8 読み込み、Markdown file / directory 判定、page list 生成を実装する。 | 単一ファイルと複数ファイルで page 順序が固定され、入力不存在と UTF-8 不正が仕様どおり失敗する。 |
| 3 | block parser、inline parser、heading slug、TOC、footnote、code fence、table、link validation を実装する。 | §3〜§7 の変換結果が deterministic で、未閉鎖 fence と broken link strict failure が fixture で再現できる。 |
| 4 | theme component、layout、`assets/style.css`、`assets/app.js`、検索 index を実装する。 | 外部 asset を参照せず、DOM id、asset path、search index schema が §7〜§8 と一致する。 |
| 5 | output directory atomic 更新、既存出力保護、stdout progress、`[WARN]`、`[REPORT]` を実装する。 | 成功時だけ出力が置換され、失敗時に既存正常出力が変更されない。 |
| 6 | §8a Fixture A〜D を追加し、生成物、report、strict / non-strict、冪等性を確認する。 | 全 fixture が通り、Phase 2 が `adlaire-ci-build` を追加判断なしに呼び出せる。 |

#### 0g.1.2 Phase 1 固定仕様

| 対象 | 固定仕様 |
|----------|----------|
| CLI 名 | 実行バイナリ名は `adlaire-ci-build` とする。実装ファイル名や package 名から別名を推測してはならない。 |
| 出力形式 | 初期実装は静的 Web サイト出力であり、単一 HTML 専用実装へ戻してはならない。 |
| theme | theme は §5 の component 名から選択する。CSS framework、CDN、外部 icon package を追加してはならない。 |
| Markdown 差異 | Go 標準ライブラリと内製 parser で処理する。外部 Markdown library を追加してはならない。 |
| report | `[REPORT]` は runner が読む契約であるため、field 名、status 値、stdout 出力位置を実装者判断で変更してはならない。 |
| URL / path | 生成 URL、asset URL、slug、relative link は §7〜§8 の規則に従う。環境依存の絶対 URL を混入してはならない。 |

### 0g.2 Phase 2 完全仕様ゲート（`components/runner.go`）

Phase 2 は、`adlaire-ci-build` を呼び出す自己ホスト型 CI runner を実装する。実装者は、本 Phase で管理 API、SDK、admin UI の endpoint や画面を実装してはならない。ただし、後続 API が読む状態ファイル schema は本 Phase で固定する。

| 項目 | 固定仕様 |
|------|----------|
| 実装開始条件 | Phase 1 が完了し、`adlaire-ci-build` の CLI 契約、`[REPORT]` 形式、出力構造が固定済みである。 |
| 入力 | `--state-dir`、`.github_token`、`.last_sha`、`.branch_config` または `BRANCH_TARGETS`、GitHub API response、`.ci/pipeline.sh`、SSH 転送設定、通知設定。 |
| 出力 | `.last_sha`、`.build_state`、`.build_status.json`、`.build_history`、`.build_logs/{id}.json`、`.pending_transfers`、`.notify_log`、`.notify_pending`、`.build_circuit_state`、`.snapshots/`、stdout/stderr slog。 |
| 状態更新順序 | lock 取得、設定整合性チェック、`.build_status.json` running 更新、`.build_state.running=true`、GitHub SHA 確認、必要時 pipeline、log/history、deploy/snapshot/notify、成功時 SHA 更新、`.build_status.json` 最終更新、`.build_state.running=false`、lock 削除の順で行う。 |
| 正常系 | 変更なし skip、変更あり build 成功、deploy なし成功、deploy 成功、pending retry 成功、通知成功、snapshot 世代管理を §13〜§16 に一致させる。 |
| 異常系 | GitHub API 失敗、rate limit、PAT 不足、pipeline 失敗、deploy 失敗、通知失敗、JSON 破損、lock 競合、stale lock、disk 不足を §11〜§16、§20、§22.0a に一致させる。 |
| セキュリティ | PAT、Webhook Secret、SMTP password を stdout、stderr、JSON log、pending queue に平文出力しない。`.github_token` は `0600`。 |
| 実装対象外 | HTTP server、API endpoint、browser SDK、admin UI、API token、session UI、MCP。 |
| 必須 fixture | §15a Fixture R1〜R7。実 GitHub / 実 SSH 接続ではなく fake server / fake executable で再現する。 |
| 完了条件 | §0e の `components/runner.go` 必須検証、§0f の完了判定、§15a の全 fixture、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.7 の runner 対象を満たす。 |

Phase 2 完了時は、Phase 3 以降へ引き継ぐ状態契約として、`.build_state`、`.build_status.json`、`.build_history`、`.build_logs/{id}.json`、`.pending_transfers`、`.notify_*`、`.branch_config`、`.build_circuit_state` の schema、権限、atomic write 条件、破損時復旧条件を固定する。Phase 3 以降は、これらの schema を暗黙に変更してはならない。

#### 0g.2.1 Phase 2 実装順序

Phase 2 は、runner が後続 API の状態正本になるため、状態ファイル schema を先に固定してから外部操作を接続する。

| 順序 | 実装単位 | 完了判定 |
|------|----------|----------|
| 1 | `--state-dir`、設定読み込み、secret 読み込み、権限確認、branch target 正規化を実装する。 | 不足 secret、権限不正、branch 設定不正が §10〜§12 の exit code / log と一致する。 |
| 2 | lock 取得、stale lock 判定、`.build_state` 初期化、atomic write 共通処理を実装する。 | 同時起動、stale lock、JSON 破損時に lock と state が仕様どおり復旧または停止する。 |
| 3 | GitHub API polling、SHA 比較、変更なし skip、rate limit / API error 処理を実装する。 | 変更なしでは pipeline を起動せず、API 異常が `.last_sha` を更新しない。 |
| 4 | pipeline 起動、`adlaire-ci-build` 呼び出し、stdout/stderr 収集、`[REPORT]` parse を実装する。 | build 成功 / 失敗が history、log、state に同一 build id で記録される。 |
| 5 | deploy、pending transfer retry、snapshot、rollback 前提データ、通知、circuit breaker を実装する。 | fake ssh / fake notifier / fake server で成功、失敗、retry、circuit open が再現できる。 |
| 6 | §15a Fixture R1〜R7 を追加し、Phase 3 が読む状態ファイルを fixture 結果で固定する。 | `.build_state`、`.build_status.json`、`.build_history`、`.build_logs/{id}.json`、queue、notify、snapshot が受け入れ条件と一致する。 |

#### 0g.2.2 Phase 2 固定仕様

| 対象 | 固定仕様 |
|----------|----------|
| `.last_sha` 更新 | GitHub SHA 確認、pipeline、deploy、snapshot、通知処理の仕様上必要な記録が成功した後にのみ更新する。失敗時更新は禁止する。 |
| build id | runner 内で一意に生成し、state、history、log、pending、notify の関連付けに同じ値を使用する。file 名と JSON field の値を一致させる。 |
| lock | lock 取得失敗時は同時起動として扱う。lock を無視して二重実行する fallback は禁止する。 |
| external command | pipeline、ssh、hook は shell 文字列ではなく引数配列で実行する。外部入力を shell に渡してはならない。 |
| secret | PAT、SSH 秘密情報、SMTP password、Webhook Secret は log、history、pending queue、notification body に平文保存しない。 |
| API 連携 | Phase 2 では HTTP endpoint を持たない。API 用に状態を固定するだけで、API server を起動しない。 |

### 0g.3 Phase 3 完全仕様ゲート（`components/api.go` P0 / P1）

Phase 3 は、管理 API の最小運用範囲を実装する。対象は §22.0f の P0 / P1 に限定し、認証、session、共通エラー、状態ファイル読み書き、ビルド操作、status、logs、history、queue、circuit breaker を固定する。

| 項目 | 固定仕様 |
|------|----------|
| 実装開始条件 | Phase 2 が完了し、runner が書き込む状態ファイル schema と lock 条件が固定済みである。 |
| 入力 | HTTP request、Authorization header、JSON body、path parameter、query parameter、`.admin_credentials`、Phase 2 の状態ファイル。 |
| 出力 | HTTP status、JSON response、SSE stream、`.access_log`、`.config_log`、必要な状態ファイル更新、systemd service log。 |
| API 範囲 | §22.0f P0 / P1 の endpoint のみ。§22.0e にない endpoint は `404`、未許可 method は `405`。 |
| 正常系 | login、logout、change password、session revoke、status、manual build、force build、cancel、stream、logs、history、queue、circuit breaker reset を §22.0〜§22.0e に一致させる。 |
| 異常系 | 認証なし、期限切れ session、権限不足、JSON 不正、body 上限超過、validation error、lock 競合、maintenance、running conflict を共通 error body で返す。 |
| セキュリティ | `GET /api/health` 以外は認証必須。session token は file 保存しない。password hash と token は response / log に平文出力しない。 |
| 実装対象外 | P2〜P5 endpoint、SDK 実装、admin UI 実装、MCP、外部 reverse proxy 設定。 |
| 必須検証 | §22.0f P0 / P1、§25、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.4.2、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.7 API 対象。 |
| 完了条件 | §0e の API 共通・状態ファイル検証、§0f の `components/api.go` 完了判定の P0 / P1 範囲、§22.0f P0 / P1 の検証を満たす。 |

Phase 3 完了時は、Phase 4〜6 へ引き継ぐ API 共通契約として、認証 header、session expiry、error body、pagination、lock error、validation error、SSE event 形式、状態ファイル read/write 境界を固定する。

#### 0g.3.1 Phase 3 実装順序

Phase 3 は、API 共通契約を後続 endpoint の土台として固定するため、endpoint 個別実装より先に共通処理を完成させる。

| 順序 | 実装単位 | 完了判定 |
|------|----------|----------|
| 1 | HTTP server 起動、routing、method 判定、body 上限、JSON decode、共通 error response を実装する。 | 未知 path、未許可 method、body 禁止、JSON 不正、body 上限超過が §22.0 と一致する。 |
| 2 | `.admin_credentials`、password hash、login、session 発行、session 期限、logout、revoke を実装する。 | `GET /api/health` 以外が認証必須になり、token / hash が response と log に出ない。 |
| 3 | runner 状態ファイルの read adapter、破損時 error、pagination、log read 境界を実装する。 | Phase 2 の状態 schema を変更せずに status、history、logs、queue を返せる。 |
| 4 | manual build、force build、cancel、stream、running conflict、lock conflict を実装する。 | 同時実行を拒否し、SSE の log / end event と cancel 結果が §22.0e と一致する。 |
| 5 | circuit breaker reset、access log、config log、systemd service 動作確認を実装する。 | 操作結果が状態ファイルと log に残り、secret を含まない。 |
| 6 | §22.0f P0 / P1、§25、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.4.2、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.7 API の確認を PR 本文に記録する。 | P0 / P1 を完了扱いにでき、Phase 4 が共通契約を再実装せず利用できる。 |

#### 0g.3.2 Phase 3 固定仕様

| 対象 | 固定仕様 |
|----------|----------|
| auth header | 認証は §22.0 の `Authorization` header 契約に従う。cookie 認証、query token、local file token は追加しない。 |
| session | session token はメモリ管理とし、状態ファイルへ保存しない。process restart 後の再 login は正常仕様とする。 |
| error body | 全 error は共通 error body に統一する。endpoint ごとに独自 error schema を返してはならない。 |
| state write | Phase 3 は P0 / P1 に必要な状態更新だけを行う。P2〜P5 の config 保存や token 発行を先取りしない。 |
| SSE | browser SDK が `fetch()` stream で読む前提に固定する。EventSource 専用仕様へ変更してはならない。 |
| external exposure | 初期 binding は `127.0.0.1` を標準とする。外部公開、TLS 終端、reverse proxy 設定は Phase 3 の対象外とする。 |

### 0g.4 Phase 4 完全仕様ゲート（`components/api.go` P2〜P5）

Phase 4 は、Phase 3 の API 共通処理を変更せず、§22.0f P2〜P5 の拡張 endpoint を追加する。Phase 4 では SDK と admin UI を実装しないが、SDK/UI が利用する endpoint 契約を最終固定する。

| 項目 | 固定仕様 |
|------|----------|
| 実装開始条件 | Phase 3 が完了し、API 共通契約、認証、session、error body、状態ファイル更新方式が固定済みである。 |
| 入力 | P2〜P5 endpoint の path、query、JSON body、状態ファイル、systemd timer 操作、通知設定、snapshot、hook、token 設定。 |
| 出力 | JSON response、状態ファイル更新、`.config_log`、`.notify_log`、snapshot archive response、hook log、token response。 |
| API 範囲 | §22.0f P2、P3、P4、P5 の endpoint。既存 P0 / P1 の response 互換を壊してはならない。 |
| 正常系 | config/repo/branch/schedule、PAT、diagnostics、dashboard、notify、SMTP、webhook、snapshot、rollback、maintenance、access control、hooks、alert rules、tag rules、pipeline config、notes、dashboard layout、tokens。 |
| 異常系 | validation error、secret mask error、systemd 操作失敗、snapshot 不存在、rollback 競合、maintenance 中実行禁止、access control block、hook timeout、token 再取得禁止を §22.0e に一致させる。 |
| セキュリティ | secret は保存時も response 時も mask 条件に従う。API token 本体は作成時 response のみ返し、再取得不可。 |
| 実装対象外 | SDK class、admin UI DOM、MCP、外部通知 channel の本ファイルで未定義の拡張。 |
| 必須検証 | §22.0f P2〜P5 の検証条件、§22.0d の read/write 対応表、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.7 API/security 対象。 |
| 完了条件 | §0e の endpoint 契約、§0f の `components/api.go` 完了判定、§22.0e の全 endpoint 契約を満たす。 |

Phase 4 完了時は、Phase 5 へ引き継ぐ SDK 契約として、全 endpoint の method、path、query、request body、response body、error status、error body、認証要否、body 禁止条件を固定する。

#### 0g.4.1 Phase 4 実装順序

Phase 4 は、P0 / P1 の互換を保持したまま endpoint 面を完成させる。既存共通処理の書き換えが必要な場合は、Phase 3 契約への影響を先に仕様化する。

| 順序 | 実装単位 | 完了判定 |
|------|----------|----------|
| 1 | §22.0d の read / write 対応表を実装対象 endpoint ごとに再確認する。 | endpoint が読む状態ファイル、書く状態ファイル、secret mask 条件が一意に決まる。 |
| 2 | config、repo、branch、schedule、PAT、diagnostics、dashboard endpoint を実装する。 | validation、atomic write、mask、systemd timer 操作の成功 / 失敗が §22.0e と一致する。 |
| 3 | notify、SMTP、webhook、snapshot、rollback、maintenance endpoint を実装する。 | snapshot 不存在、rollback 競合、maintenance 中実行禁止、通知失敗が仕様どおり返る。 |
| 4 | access control、hooks、alert rules、tag rules、pipeline config、notes、dashboard layout endpoint を実装する。 | allow / block、hook timeout、tag rule、notes、layout 保存が状態 schema と一致する。 |
| 5 | API token endpoint を実装し、作成時だけ token 本体を返す。 | token 再取得不可、token hash 保存、失効、権限不足 error が §22.0e と一致する。 |
| 6 | 全 endpoint の method、path、query、body、response、error、auth を SDK 実装用契約として PR 本文に記録する。 | Phase 5 が追加判断なしに SDK method を実装できる。 |

#### 0g.4.2 Phase 4 固定仕様

| 対象 | 固定仕様 |
|----------|----------|
| P0 / P1 互換 | Phase 4 で P0 / P1 の response field 名、status code、error body、auth 条件を変更してはならない。 |
| secret response | secret は保存成功時も mask 表示のみ返す。平文 secret を返せるのは API token 作成直後の token 本体だけとする。 |
| rollback | rollback は snapshot から出力物を戻す操作であり、history、secret、state 全体を巻き戻してはならない。 |
| maintenance | maintenance 中に禁止する操作と許可する read 操作は §22.0e に従う。実装者判断で追加禁止しない。 |
| hook | hook は登録済み引数配列のみ実行する。shell 展開、環境変数補完、任意 script 文字列実行は禁止する。 |
| endpoint 追加 | SDK/UI の都合で §22.0e にない endpoint を追加してはならない。必要な場合は先に仕様改訂する。 |

### 0g.5 Phase 5 完全仕様ゲート（`admin/adlaire-ci-sdk.js`）

Phase 5 は、固定済み API 契約に対する browser SDK を単一 ES Module として実装する。SDK は API 通信抽象化だけを責務とし、UI 表示、DOM 操作、状態ファイル直接操作を行わない。

| 項目 | 固定仕様 |
|------|----------|
| 実装開始条件 | Phase 3 / Phase 4 が完了し、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e の全 endpoint 契約が固定済みである。 |
| 入力 | constructor の `baseUrl`、method 引数、browser 標準 API、session token。 |
| 出力 | `Promise`、`AdlaireCIError`、`StreamHandle`、JSON object、SSE callback / stream reader、token 破棄。 |
| 正常系 | `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 の全 method が `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e の endpoint のみを呼び、query、body、HTTP method、response 変換、stream close を仕様どおり処理する。 |
| 異常系 | HTTP error、network error、timeout、unsupported browser、body 禁止 endpoint、認証失敗、`401` token 破棄を `AdlaireCIError` 契約に一致させる。 |
| セキュリティ | token を localStorage / sessionStorage へ保存しない。secret 入力値を console に出力しない。global 汚染を行わない。 |
| 実装対象外 | DOM 操作、admin panel 表示、API endpoint 新設、Node.js 専用 API、runtime 名別の互換分岐、bundler、polyfill、npm package。 |
| 必須検証 | `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 の method 契約、§0e の SDK 契約、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.7 SDK 対象。API は fake fetch で成功 / error / timeout / stream を再現する。 |
| 完了条件 | 全 method が `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e と `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 の対応どおり動作し、HTTP error を `AdlaireCIError` として扱い、body 禁止 endpoint に body を送らない。 |

Phase 5 完了時は、Phase 6 へ引き継ぐ UI 契約として、SDK method 名、引数、戻り値、error object、loading / retry に必要な失敗情報、stream handle を固定する。

#### 0g.5.1 Phase 5 実装順序

Phase 5 は API 契約の薄い wrapper とし、表示判断を持たせない。実装は fake fetch による endpoint 対応確認から進める。

| 順序 | 実装単位 | 完了判定 |
|------|----------|----------|
| 1 | `AdlaireCI` constructor、base URL 正規化、token 保持、共通 request、query 生成を実装する。 | 末尾 slash、query encoding、Authorization header、body 禁止 endpoint が `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 と一致する。 |
| 2 | `AdlaireCIError`、HTTP error 変換、network error、timeout、unsupported browser 判定を実装する。 | UI が `status`、`message`、`details`、`requestId` を追加判断なしに表示できる。 |
| 3 | `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 の non-stream method を endpoint 対応表どおり実装する。 | 各 method の HTTP method、path、query、body、戻り値が `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e と一致する。 |
| 4 | `streamBuild()` と `StreamHandle.close()` を実装する。 | `fetch()` stream、AbortController、SSE frame parse、invalid frame、end event が `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 と一致する。 |
| 5 | `401` token 破棄、logout、revoke、token 再設定、secret 非保存を実装する。 | token が storage に保存されず、`401` 後に認証付き request が送られない。 |
| 6 | fake fetch で成功、HTTP error、network error、timeout、stream、body 禁止を検証する。 | Phase 6 が SDK method だけで全 UI 操作を実装できる。 |

#### 0g.5.2 Phase 5 固定仕様

| 対象 | 固定仕様 |
|----------|----------|
| module 形式 | 単一 ES Module とする。npm package、bundler、transpiler、Node.js 専用 API は使用しない。 |
| storage | token、password、PAT、SMTP password、Webhook Secret を localStorage、sessionStorage、IndexedDB、cookie に保存しない。 |
| response 補完 | API response にない値を SDK が推測して追加しない。表示用加工、既定ラベル、並べ替えは UI 側の責務とする。 |
| stream | native EventSource は Authorization header を付与できないため使用しない。`fetch()` stream を標準とする。 |
| direct DOM | SDK は DOM を読まない、書かない。UI 状態、message、focus、disabled を操作してはならない。 |
| endpoint drift | §22.0e と一致しない endpoint、method、body、query を SDK 都合で追加してはならない。 |

### 0g.6 Phase 6 完全仕様ゲート（`admin/index.html`）

Phase 6 は、標準管理ツール UI を単一 HTML と Vanilla JavaScript で実装する。UI は `admin/adlaire-ci-sdk.js` 経由でのみ API と通信し、API を直接 `fetch` してはならない。

| 項目 | 固定仕様 |
|------|----------|
| 実装開始条件 | Phase 5 が完了し、SDK method 契約、error object、stream handle が固定済みである。 |
| 入力 | ユーザー操作、form input、SDK response、SDK error、SSE / stream event、same-origin 配信環境。 |
| 出力 | DOM 表示、成功/失敗表示、disabled 状態、loading 状態、再取得、secret field 消去、SDK method 呼び出し。 |
| 正常系 | `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 の panel、DOM id、表示条件、操作、成功表示、再取得、stream 表示、logout を仕様どおり処理する。 |
| 異常系 | login 失敗、session 期限切れ、API error、validation error、stream 切断、network error、maintenance、access block を UI 表示契約に一致させる。 |
| セキュリティ | secret 値を DOM に残さない。password / PAT / token / SMTP password / webhook secret は成功・失敗に関わらず指定条件で消去する。SDK を介さない直接 API 呼び出しは禁止。 |
| 実装対象外 | frontend framework、CSS framework、chart library、CDN script、build tool、外部 icon package、API endpoint 新設。 |
| 必須検証 | `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 の全 panel と主要操作、§0e の UI 契約、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.7 UI/security 対象。SDK は fake implementation で成功 / 失敗 / loading / stream を再現する。 |
| 完了条件 | 全 UI 操作が `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 の表示条件と `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 の SDK method を満たし、成功表示、失敗表示、disabled、再取得、secret 消去が一致する。 |

Phase 6 完了時は、初期実装全体の完了判定として、Phase 1〜6 の引き継ぎ契約、§0e、§0f、§0g、§0i、§0j、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0f、`ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23、`ADLAIRE_CI_DETAIL_UI_SPEC.md` §24、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.7 を再確認する。

#### 0g.6.1 Phase 6 実装順序

Phase 6 は、SDK 契約の利用者として UI を実装する。API 仕様の不足を UI 側で補完してはならない。

| 順序 | 実装単位 | 完了判定 |
|------|----------|----------|
| 1 | HTML shell、panel、DOM id、navigation、初期 disabled / loading 状態を実装する。 | `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 の DOM id と panel 構成が一致し、未認証時に保護操作が disabled になる。 |
| 2 | login、logout、session expiry、共通 error 表示、再 login 導線を実装する。 | `AdlaireCIError` の message / details が同一 panel 内に表示され、secret が残らない。 |
| 3 | status、history、logs、queue、manual build、force build、cancel、stream 表示を実装する。 | SDK method だけを呼び、loading、disabled、再取得、stream close が `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 と一致する。 |
| 4 | config、repo、branch、schedule、notify、snapshot、rollback、maintenance、access、hooks、tokens の panel 操作を実装する。 | 成功表示、失敗表示、secret field 消去、保存後再取得が各 panel の仕様と一致する。 |
| 5 | dashboard layout、notes、alert、tag rule、pipeline config の表示 / 保存を実装する。 | UI が未定義 field を追加せず、SDK response に基づく表示だけを行う。 |
| 6 | fake SDK で成功、失敗、loading、stream、session expiry、secret 消去を検証する。 | 初期実装全体の UI 受け入れ条件を満たし、直接 API 呼び出しが存在しない。 |

#### 0g.6.2 Phase 6 固定仕様

| 対象 | 固定仕様 |
|----------|----------|
| API 通信 | UI は必ず SDK method を呼ぶ。`fetch()`、`XMLHttpRequest`、直接 `ReadableStream` 生成は禁止する。 |
| framework | frontend framework、CSS framework、chart library、CDN script、外部 icon package、build tool は使用しない。 |
| secret field | password、PAT、token、SMTP password、Webhook Secret は成功 / 失敗に関わらず `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 の条件で消去する。 |
| unknown state | SDK response にない状態を UI が推測して成功扱いにしてはならない。不明状態は error または未取得として表示する。 |
| disabled | 実行中、未認証、maintenance、権限不足、入力不正時の disabled 条件は `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 に従う。 |
| UI 追加 | `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 にない panel、操作、保存先、設定項目を UI 都合で追加してはならない。 |

### 0g.7 Phase 間引き継ぎ契約

各 Phase の完了時は、次 Phase が依存する契約を変更不可として扱う。後続 Phase で変更が必要になった場合は、後続 Phase の実装で吸収せず、契約を定義した Phase の仕様改訂へ戻す。

| 引き継ぎ元 | 引き継ぎ先 | 固定する契約 |
|------------|------------|--------------|
| Phase 1 | Phase 2 | `adlaire-ci-build` CLI、終了コード、stdout / stderr、`[REPORT]`、出力サイト構造。 |
| Phase 2 | Phase 3 | runner 状態ファイル schema、lock、history/log、pending queue、circuit breaker、snapshot、通知ログ。 |
| Phase 3 | Phase 4 | API 共通契約、認証、session、error body、validation、lock error、SSE 基本形式。 |
| Phase 4 | Phase 5 | 全 endpoint の method、path、query、request、response、error、認証要否。 |
| Phase 5 | Phase 6 | SDK method 名、引数、戻り値、error object、stream handle、token 破棄条件。 |
| Phase 6 | 初期実装完了 | UI 操作、表示状態、secret 消去、SDK 経由通信、実装完了検証結果。 |

引き継ぎ契約は、実装 PR 本文に `固定契約 / 参照節 / 検証結果 / 後続 Phase への影響` の形式で記録する。

### 0g.8 Phase 別 実装 PR 成果物チェックリスト

各 Phase の実装 PR は、実装コードだけでなく、仕様どおりに実装したことを再現できる成果物を含める。下表の必須成果物が不足する PR は、その Phase を完了扱いにしてはならない。

| Phase | 変更対象ファイル | fixture / testdata | PR 本文に固定する契約 | 必須検証 | 実装対象外として明記するもの |
|-------|------------------|--------------------|------------------------|----------|------------------------------|
| Phase 1 | `components/builder.go`、必要な Go test、`testdata/builder/`。 | §8a Fixture A〜D の入力、期待出力、失敗系入力。 | CLI option、終了コード、stdout / stderr、`[WARN]`、`[REPORT]`、出力サイト構造、asset path、search index schema。 | `gofmt -l`、Go test、fixture A〜D、冪等性、strict / non-strict、外部 asset 不存在。 | runner、GitHub API、SSH、API server、SDK、admin UI、MCP。 |
| Phase 2 | `components/runner.go`、必要な Go test、`testdata/runner/`、fake GitHub server、fake ssh / notifier。 | §15a Fixture R1〜R7 の状態ディレクトリ、API 応答、pipeline 結果、deploy / notify 結果。 | 状態ファイル schema、lock、build id、`.last_sha` 更新条件、pending queue、snapshot、notify、circuit breaker。 | `gofmt -l`、Go test、secret 不足、lock 競合、変更なし skip、build 成功 / 失敗、deploy retry、通知失敗、JSON 破損。 | HTTP API、SDK、admin UI、API token、session UI、MCP。 |
| Phase 3 | `components/api.go`、必要な Go test、API fixture、状態ファイル fixture、systemd service 確認資料。 | P0 / P1 endpoint の request / response、認証あり / なし、SSE、状態 read / write fixture。 | API 共通 error body、auth header、session expiry、pagination、SSE event、P0 / P1 endpoint 契約、状態 read/write 境界。 | `gofmt -l`、Go test、login/logout、認証なし、未知 path、body 禁止、JSON 不正、status/history/logs/queue、manual build、cancel、stream。 | P2〜P5 endpoint、SDK、admin UI、MCP、外部公開設定。 |
| Phase 4 | `components/api.go`、必要な Go test、P2〜P5 API fixture、secret mask fixture。 | config、repo、branch、schedule、notify、snapshot、rollback、maintenance、access、hooks、tokens の request / response。 | 全 endpoint の method、path、query、request、response、error、auth、secret mask、token 再取得不可条件、P0 / P1 互換。 | `gofmt -l`、Go test、P2〜P5 endpoint、validation、secret mask、rollback、maintenance、hook timeout、token 発行 / 失効、P0 / P1 回帰確認。 | SDK class、admin UI DOM、MCP、本ファイルで未定義通知 channel。 |
| Phase 5 | `admin/adlaire-ci-sdk.js`、SDK test fixture、fake fetch / stream fixture。 | 成功 response、HTTP error、network error、timeout、SSE frame、invalid frame、body 禁止 endpoint。 | SDK method 名、引数、戻り値、`AdlaireCIError`、`StreamHandle`、token 破棄、query 生成、body 禁止。 | browser runtime または同等環境で fake fetch 検証、HTTP error、timeout、stream、`401` token 破棄、storage 不使用確認。 | DOM 操作、admin UI、API endpoint 新設、Node.js 専用 API、bundler、npm package。 |
| Phase 6 | `admin/index.html`、UI test fixture、fake SDK、必要な静的 asset。 | fake SDK の成功、失敗、loading、stream、session expiry、secret 入力 fixture。 | DOM id、panel、SDK method 対応、success / error 表示、disabled、loading、再取得、secret 消去、直接 API 呼び出し禁止。 | UI 操作確認、fake SDK 成功 / 失敗、loading、stream、session expiry、secret 消去、直接 `fetch()` 不存在。 | API endpoint 新設、SDK 契約変更、frontend framework、CSS framework、CDN、build tool、MCP。 |

各 Phase の実装 PR は、本文に `成果物 / 固定契約 / 検証 / 未実装対象 / 後続 Phase への影響` を記録する。未実行の検証がある場合は、環境理由だけで合格扱いにせず、未完了として扱う。

**Phase fixture / testdata / fake / PR 証跡契約：**

Phase fixture / testdata 配置、fake 実装、実装 PR 証跡の詳細は `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §0g.8-F を正とする。本節には Phase 成果物チェックリストだけを置き、fixture manifest、assertion、fake、PR 証跡の本文を重複定義しない。

#### 0g.8.1 実装 PR 完了扱い禁止条件

以下のいずれかに該当する場合、実装 PR は merge 可能であっても Phase 完了として扱ってはならない。

| 条件 | 扱い |
|------|------|
| 対象 Phase の必須成果物が不足している。 | PR 本文に不足項目を明記し、Phase 未完了とする。 |
| fixture または testdata がなく、手動確認だけで合格としている。 | 仕様上 fixture が必須の Phase では未完了とする。 |
| 実装対象外の機能を先取りしている。 | 仕様違反として扱い、対象外機能を削除するか、先に仕様改訂する。 |
| 後続 Phase が依存する契約を PR 本文に固定していない。 | 後続 Phase の実装を開始してはならない。 |
| secret、token、password を log、fixture、snapshot、UI 表示へ平文出力している。 | security 不合格として Phase 未完了とする。 |
| §0g の順序、§0i の対応表、§0j のソース配置、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.7 の受け入れ条件のいずれかと矛盾している。 | 仕様不整合として扱い、実装または仕様を修正する。 |

---

## 0h. 機能仕様テンプレート

対象項目を追加または改訂する場合は、該当する詳細仕様節に以下の項目をすべて含める。既存節に含める場合も、実装者が下表の項目を本文から一意に読み取れる状態にする。

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

上表のいずれかが不足する対象項目は、実装者判断で補完してはならない。不足を見つけた場合は、実装 PR ではなく仕様改訂 PR として本ファイルを先に更新する。

---

## 0i. 詳細節対応表

本節は、対象機能から本ファイル内の実装詳細へ移動するための対応表である。実装者は対象機能を実装する前に、下表の「詳細仕様節」と「受け入れ条件」を確認する。

表の「責務 component」は参照先を探すための component 一覧である。owner component と collaborator component は、対象機能の詳細仕様節に記載された値を正とする。

表の「詳細仕様節」が複数ある場合は、すべての節を同時に満たす。`components/builder.go` が責務 component に含まれる機能では、`ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` の同番号節を合わせて確認する。`components/runner.go` が責務 component に含まれる機能では、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` の同番号節を合わせて確認する。`components/api.go` が責務 component に含まれる機能では、`ADLAIRE_CI_DETAIL_API_SPEC.md` の同番号節を合わせて確認する。`admin/adlaire-ci-sdk.js` が責務 component に含まれる機能では、`ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 を合わせて確認する。`admin/index.html` が責務 component に含まれる機能では、`ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 を合わせて確認する。セットアップ、アップデート、バイナリ配布、systemd、リリース成果物検証に関わる機能では、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 を合わせて確認する。該当節に §0h の必須項目が不足している場合は、その項目を実装せず、先に詳細仕様を改訂する。

| 機能 | 責務 component | 詳細仕様節 | 受け入れ条件 |
|------|-------------------|------------|--------------|
| ビルドタイムアウト | `components/runner.go` / `components/api.go` | §12、§13、§22.0e | `build_timeout_seconds` の既定値、設定 API、`context.WithTimeout` の中断処理、終了コード、ログが一致する。 |
| ポーリング間隔の動的変更 | `components/api.go` | §22.0e、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26、§27.11 | `POST /api/schedule/interval` が systemd timer 設定を更新し、検証コマンドで反映を確認できる。 |
| ビルドログのファイル保存 | `components/runner.go` | §11、§13、§15 | `.build_logs/{id}.json` の schema、stdout/stderr、変換レポート、duration、権限が一致する。 |
| GitHub Webhook 受信 | `components/api.go` / `components/runner.go` | §22.0e、§22-W、§13、§27.12 | HMAC 検証、イベント記録、キュー投入またはビルドトリガー、エラー応答が一致する。 |
| ネットワーク断時の再試行 | `components/runner.go` | §12、§13 | `API_RETRY_MAX`、`API_RETRY_BASE_SECONDS`、指数バックオフ、失敗時ログが一致する。 |
| GitHub API レート制限自動待機 | `components/runner.go` / `components/api.go` | §13、§22.0e | `X-RateLimit-Remaining` と `X-RateLimit-Reset` の扱い、待機、API 表示が一致する。 |
| 転送後リモート整合性検証 | `components/runner.go` | §14a、§13 | SSH 転送後の SHA256 照合、不一致時の `.pending_transfers` 再投入、ログが一致する。 |
| マルチブランチビルド | `components/runner.go` / `components/api.go` | §12、§13、§22.0e | `BRANCH_TARGETS` と `.branch_config` の優先順位、順次処理、API 更新が一致する。 |
| ビルドログ世代管理 | `components/runner.go` | §12、§13、§15 | `LOG_KEEP_N` 超過時の削除順序、0 の扱い、削除ログが一致する。 |
| ビルド出力の外部転送 | `components/runner.go` | §14a、§13 | SSH 差分転送、複数ファイル処理、失敗時 pending、通知が一致する。 |
| ビルドクールダウン | `components/runner.go` | §12、§13 | `BUILD_COOLDOWN_SECONDS` 内の起動スキップ、Webhook 二重トリガー抑止、ログが一致する。 |
| ビルド前の事前チェック | `components/runner.go` | §13、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 | ディスク、`adlaire-ci-build`、pipeline 前提の確認、不足時の ERROR と通知が一致する。 |
| 定期強制ビルド | `components/runner.go` / `components/api.go` | §12、§13、§22.0e | `FORCE_BUILD_INTERVAL`、変更なし時の強制ビルド、設定 API が一致する。 |
| ビルド中重複スキップ | `components/runner.go` | §11、§13 | `.build_lock` の PID 判定、stale lock、競合時終了コードとログが一致する。 |
| GitHub PAT 有効期限の事前警告 | `components/runner.go` / `components/api.go` | §13、§22.0e | `GitHub-Authentication-Token-Expiration` の解析、7 日以内 WARN、API 表示が一致する。 |
| コミット情報のビルドログ記録 | `components/runner.go` | §13、§15 | SHA、message、author、date を build id と同じログへ記録する。 |
| GitHub API 連続失敗によるサーキットブレーカー | `components/runner.go` / `components/api.go` | §11、§12、§13、§22.0e | 閾値、open/close 状態、API reset、通知、状態ファイルが一致する。 |
| 出力サイトサイズ警告閾値 | `components/builder.go` / `components/runner.go` / `components/api.go` | §8、§12、§13、§22.0e | `OUTPUT_SIZE_WARN_MB`、`size_warn`、WARN ログ、API 表示が一致する。 |
| 設定ファイル起動時整合性チェック | `components/runner.go` | §11、§12、§13、§22.0a、§22.0c、§27.10 | 対象 JSON ファイル、検証順序、破損退避、初期化値、ログ、通知、終了コード、fixture が一致する。 |
| ビルドステータスファイル出力 | `components/runner.go` / `components/api.go` | §11、§13、§15、§22.0a、§22.0c、§22.0e、§27.8 | `.build_status.json` の schema、更新タイミング、status/target_status、pending 件数、circuit 状態、API 参照元が一致する。 |
| ビルドトリガー種別の記録 | `components/runner.go` / `components/api.go` / `admin/adlaire-ci-sdk.js` / `admin/index.html` | §13、§15、§22.0c、§22.0e、§23、§24、§27.9 | `trigger` の有効値、判定条件、`.build_logs`、`.build_history`、`.build_status.json`、履歴 filter、UI 表示が一致する。 |
| GitHub Commit Status API | `components/runner.go` | §12、§13、§15、§22.0c、§27.1 | `commit_status_enabled`、context、target_url、pending/success/failure の送信条件、失敗時の扱い、build log 記録が一致する。 |
| ドライラン実行モード | `components/runner.go` | §11、§12、§13、§15、§27.2 | `--dry-run` が状態ファイル、log、history、deploy、通知を変更せず、設定・GitHub・SHA 判定結果を固定 JSON で返す。 |
| ビルド失敗時の自動リトライ | `components/runner.go` | §12、§13、§15、§22.0c、§27.3 | retry 対象エラー、最大回数、backoff、attempt log、最終 status、SHA 更新禁止条件が一致する。 |
| 出力サイトへのビルドメタ埋め込み | `components/builder.go` / `components/runner.go` / `components/api.go` | §2、§5、§8、§13、§22.0e、§27.4 | CLI/env 入力、HTML meta、REPORT、build log、`GET /api/output-meta` の値が一致する。 |
| 設定バリデーション API | `components/api.go` / `admin/adlaire-ci-sdk.js` / `admin/index.html` | §22.0c、§22.0e、§23、§24、§27.5 | `POST /api/config/validate` が状態を変更せず、正規化後設定、warnings、errors を返す。 |
| API アクセスログ | `components/api.go` / `admin/adlaire-ci-sdk.js` / `admin/index.html` | §22.0a、§22.0c、§22.0e、§23、§24、§27.6 | `.api_access_log` の schema、追記対象、マスク条件、一覧 API、UI 表示が一致する。 |
| ビルドログのアーカイブ圧縮 | `components/runner.go` / `components/api.go` / `admin/adlaire-ci-sdk.js` / `admin/index.html` | §12、§13、§15、§22.0c、§22.0e、§23、§24、§27.7 | gzip 形式、archive 先、参照順、cleanup/archive API、disk usage 集計、UI 表示が一致する。 |
| Webhook イベントログ | `components/api.go` | §11、§22.0e、§22-W、§27.13 | `.webhook_events.json` の JSON Lines schema と一覧 API が一致する。 |
| ビルド所要時間の記録と統計 API | `components/runner.go` / `components/api.go` | §15、§22.0e、§27.14 | `started_at`、`finished_at`、`duration_seconds` と統計 API が一致する。 |
| ビルドアーティファクト世代管理 | `components/runner.go` / `components/api.go` | §14b、§22.0e | `.snapshots/` の保持世代、削除、rollback API が一致する。 |
| ビルドアーティファクト管理 | `components/api.go` / `admin/index.html` / `admin/adlaire-ci-sdk.js` | §14b、§22.0e、§23、§24、§27.15 | 一覧、ダウンロード、削除、ロールバックの API、SDK、UI が一致する。 |
| ヘルスチェックエンドポイント | `components/api.go` | §22.0e、§27.16 | `GET /api/health` の稼働秒数、最終ビルド、最終転送、エラー応答が一致する。 |
| Webhook イベント一覧取得 API | `components/api.go` / `admin/adlaire-ci-sdk.js` / `admin/index.html` | §22.0e、§23、§24、§27.13 | `GET /api/webhook-events` の query、response、SDK method、UI 表示が一致する。 |
| ビルドログ重大度フィルター | `components/api.go` / `admin/adlaire-ci-sdk.js` / `admin/index.html` | §22.0e、§23、§24、§27.17 | `level=warn\|error` の query、検索結果、UI filter が一致する。 |
| 変換レポート出力 | `components/builder.go` / `components/runner.go` / `components/api.go` | §8、§13、§15、§22.0e | `[REPORT]` stdout、runner 取り込み、`.build_logs`、`GET /api/output-meta` が一致する。 |
| シンタックスハイライト | `components/builder.go` | §7.8 | 対応言語、class 名、HTML escape、CSS 表示が一致する。 |
| 本文内全文検索 | `components/builder.go` | §7.9 | `assets/search-index.json`、検索 UI、ヒット遷移、対象テキストが一致する。 |
| アンカーリンク自動検証 | `components/builder.go` | §4.3、§8 | broken anchor 検出、`[WARN] BROKEN_LINK`、report 件数が一致する。 |
| コードブロックの折りたたみ | `components/builder.go` | §7.10 | 30 行超の初期折りたたみ、展開操作、印刷時展開が一致する。 |
| 印刷スタイル（`@media print`） | `components/builder.go` | §6 | `@media print` の非表示対象、コード展開、リンク URL 表示が一致する。 |
| 静的 Web サイト出力 | `components/builder.go` | §2、§5、§6、§7 | 入力ファイル/ディレクトリ、出力ファイル構成、asset、ページ生成が一致する。 |
| テーマコンポーネント | `components/builder.go` | §5、§6、§7 | `adlaire-default` の component、class、slot、asset 出力が一致する。 |
| 外部リンクの自動処理 | `components/builder.go` | §4.3 | `target="_blank"`、`rel="noopener noreferrer"`、内部リンクとの区別が一致する。 |
| 読み取り進捗バー | `components/builder.go` | §7.13 | 3px 固定表示、scroll 連動、初期/末尾状態が一致する。 |
| コードブロックのコピーボタン | `components/builder.go` | §7.6 | ボタン配置、コピー対象、成功/失敗時表示、アクセシビリティが一致する。 |
| 見出しアンカーリンクコピー | `components/builder.go` | §3、§7.11 | `.hn-link`、copy URL、重複 slug 連動が一致する。 |
| TOC 開閉状態の永続化 | `components/builder.go` | §7.3 | `localStorage` key、展開/折りたたみ、復元条件が一致する。 |
| 見出しスラグ重複解決 | `components/builder.go` | §4.5 | `-2`、`-3` の付与、TOC、検索、コピー URL との共通化が一致する。 |
| 前後章ナビゲーションボタン | `components/builder.go` | §4.5、§5、§7.15 | h2 単位の前後判定、章末尾配置、端の非表示条件が一致する。 |
| 内部リンク整合性チェック | `components/builder.go` | §4.3、§8 | `[label](#anchor)` 検証、WARN、`broken_links` が一致する。 |
| 見出し階層スキップ警告 | `components/builder.go` | §4.5、§8 | h1→h3 等の検出、WARN、`heading_skips` が一致する。 |
| 読了時間推計と表示 | `components/builder.go` | §4.5、§5、§6、§8 | 対象文字数、200文字/分、切り上げ、header 表示、report が一致する。 |
| Webhook 通知失敗リトライキュー | `components/runner.go` | §11、§13、§16 | `.notify_pending` の schema、再送順序、失敗時保持が一致する。 |
| ブランチ設定の動的変更 API | `components/runner.go` / `components/api.go` | §11、§12、§22.0e、§27.18 | `.branch_config`、GET/POST API、runner 再起動不要条件が一致する。 |
| 週次ビルドサマリー Webhook | `components/runner.go` / `components/api.go` | §12、§13、§16、§22.0e、§27.19 | 週次判定、集計対象、通知 payload、手動送信 API が一致する。 |
| 設定変更の詳細 diff 記録 | `components/api.go` | §22.0a、§22.0e、§27.20 | `.config_log` の diff 文字列、対象 API、マスク条件が一致する。 |
| テーブルのソート機能 | `components/builder.go` | §7.14 | クリック操作、昇順/降順、`aria-sort`、インジケーターが一致する。 |
| キーボードショートカット | `components/builder.go` | §7.12 | `/`、`Escape`、`t` の対象、フォーカス条件、入力中の無効化が一致する。 |
| 複数ファイル監視 | `components/runner.go` / `components/builder.go` / `components/api.go` | §11、§12、§13、§15、§27.21 | `target_files` の検証、対象別 SHA 差分、build target 決定、履歴・ログ・API 表示が一致する。 |
| ビルドパイプライン YAML 定義 | `components/runner.go` / `components/api.go` | §12、§13、§15、§22.0e、§27.22 | `.pipeline.yml` の内製 subset parse、step 実行順、timeout、env、失敗時 status、API 保存が一致する。 |
| ローカルファイル監視モード | `components/runner.go` | §11、§12、§13、§27.23 | GitHub API を呼ばず、local snapshot の SHA-256 差分だけで変更検出し、trigger と status が一致する。 |
| タグ付きコミットのみビルド | `components/runner.go` / `components/api.go` | §12、§13、§15、§22.0e、§27.24 | tag pattern、GitHub tags API、skip 条件、build log/history、設定 API が一致する。 |
| ビルドキャッシュ | `components/builder.go` / `components/runner.go` | §5、§8、§11、§13、§27.25 | `.build_cache.json`、入力 manifest、再利用条件、無効化条件、report counters が一致する。 |
| 並列マルチターゲットビルド | `components/runner.go` | §12、§13、§14a、§15、§27.26 | worker 上限、target 別 status、pending transfer、最終 build status、ログ順序が一致する。 |
| ビルド前後フック | `components/runner.go` / `components/api.go` | §13、§15、§22.0e、§27.27 | `.hooks` schema、pre/post 実行、abort 条件、hook log、API CRUD が一致する。 |
| 依存ファイルトラッキング | `components/builder.go` / `components/runner.go` | §4.3、§5、§11、§13、§27.28 | `.dependency_manifest.json`、依存抽出、関連 target 判定、破損時 full build が一致する。 |
| リモートビルド対応 | `components/runner.go` / `components/api.go` | §12、§13、§14a、§15、§27.29 | remote command、archive 取得、manifest 検証、状態記録、失敗時 rollback 不実行が一致する。 |
| ビルド承認フロー | `components/runner.go` / `components/api.go` / `admin/adlaire-ci-sdk.js` / `admin/index.html` | §11、§13、§15、§16、§22.0e、§27.30 | `.approval_queue`、承認/却下 API、通知、timeout、UI 操作、履歴 status が一致する。 |
| ブランチ別環境変数 | `components/runner.go` / `components/api.go` | §12、§13、§15、§22.0e、§27.31 | branch env schema、許可 key、secret mask、process env 注入、API 保存が一致する。 |
| ビルド通知連携 | `components/runner.go` / `components/api.go` / `admin/adlaire-ci-sdk.js` / `admin/index.html` | §13、§15、§16、§22.0e、§27.32 | 通知 event、channel schema、送信順、retry、mask、notify log、API / UI 表示が一致する。 |
| ビルド時間トレンド記録 | `components/runner.go` / `components/api.go` | §13、§15、§22.0e、§27.33 | `.build_trends.json`、移動平均、中央値、p95、API response、破損時復旧が一致する。 |
| ビルド依存チェーン | `components/runner.go` / `components/api.go` | §11、§13、§15、§22.0e、§27.34 | `.build_chain_config`、依存 DAG 検証、実行順、skip / failure status、chain log が一致する。 |
| ビルド優先度キュー | `components/runner.go` / `components/api.go` | §11、§13、§22.0e、§27.35 | queue priority、created_seq、同一優先度 FIFO、API 表示、cancel / clear が一致する。 |
| 失敗原因の自動分類 | `components/runner.go` / `components/api.go` | §13、§15、§22.0e、§27.36 | failure_category、evidence、分類優先順位、history / log / UI 表示が一致する。 |
| ビルド実行環境の記録 | `components/runner.go` | §13、§15、§27.37 | build 開始時の environment snapshot、secret 非含有、log schema、検証 fixture が一致する。 |
| ビルド所要時間の異常検知 | `components/runner.go` / `components/api.go` | §13、§15、§16、§22.0e、§27.38 | trend 基準、異常判定、WARN、history flag、通知 payload、設定値が一致する。 |
| ビルドトリガー専用 API スコープ | `components/api.go` / `admin/adlaire-ci-sdk.js` / `admin/index.html` | §22.0a、§22.0c、§22.0e、§23、§24、§27.42 | `trigger` scope token が build 起動系だけを許可し、その他 API を拒否する。 |
| API キー管理 | `components/api.go` / `admin/adlaire-ci-sdk.js` / `admin/index.html` | §22.0a、§22.0c、§22.0e、§23、§24、§27.43 | API key 本体の一回表示、hash 保存、scope、期限、失効、監査が一致する。 |
| 監査ログ | `components/api.go` / `admin/adlaire-ci-sdk.js` / `admin/index.html` | §22.0a、§22.0c、§22.0e、§23、§24、§27.44 | `.audit_log` schema、対象操作、mask、検索 API、UI 表示が一致する。 |
| セッションタイムアウト変更設定 | `components/api.go` / `admin/adlaire-ci-sdk.js` / `admin/index.html` | §22.0c、§22.0e、§23、§24、§25、§27.45 | `session_timeout_seconds` の範囲、保存、既存 session の扱い、新規 session 期限が一致する。 |
| TOTP 二要素認証 | `components/api.go` / `admin/adlaire-ci-sdk.js` / `admin/index.html` | §22.0a、§22.0c、§22.0e、§23、§24、§25、§27.46 | RFC 6238 TOTP、二段階 login、secret 保存、確認、無効化、UI 操作が一致する。 |
| API レート制限 | `components/api.go` / `admin/adlaire-ci-sdk.js` / `admin/index.html` | §22.0a、§22.0c、§22.0e、§23、§24、§27.47 | IP / actor / endpoint group の固定窓制限、`429`、状態保存、設定 API、UI 表示が一致する。 |

---

## 0. システム概要

Adlaire CI は Go 版 3 コンポーネントと JavaScript/HTML 管理ツールで構成する。

本ファイルは、`components/builder.go`、`components/runner.go`、`components/api.go`、標準管理ツール `admin/index.html`、JavaScript SDK `admin/adlaire-ci-sdk.js` の実装詳細を定義する。`components/mcp.go` の入出力、状態、起動手順、検証条件は本ファイルでは定義しない。

**`components/builder.go`（ビルドスクリプト）**
GitHub リポジトリ上またはローカル上の Markdown ファイルまたは Markdown ディレクトリを静的 Web サイトに変換してローカルディレクトリへ出力する。標準実行バイナリ名は `adlaire-ci-build` とする。

**`components/runner.go`（CI ランナー）**
GitHub の Git Trees API / Git Blobs API を使用し、対象ファイルの blob SHA 変更を検出する。変更があった場合のみ Markdown 本文を書き出し、`adlaire-ci-build` を起動し、成功時に SHA キャッシュを更新する。systemd タイマーで定期実行する oneshot 設計。

SSH 転送、ペンディングキュー、スナップショット、Webhook 通知、マルチブランチ、ビルドログ保存、サーキットブレーカーは Go 版 `components/runner.go` の対象機能である。

**`components/api.go`（管理 API サーバー）**
Go 標準ライブラリ `net/http` を使用する常駐 HTTP サーバー。管理ツールからの API リクエストを受け付け、認証・状態取得・手動ビルドトリガーを処理する。`adlaire-ci-api.service` として systemd に登録し、`components/runner.go` とは独立して常駐する。

**Go 版実行フロー：**
```
systemd timer
  └─ adlaire-ci-runner（components/runner.go, oneshot）
       ├─ 変更なし → スキップ
       └─ 変更あり → adlaire-ci-build（components/builder.go）→ 静的 Web サイト生成
```

**管理 API を含む想定フロー：**
```
adlaire-ci-runner
  └─ SSH 転送 / スナップショット / Webhook 通知 / ビルドログ保存

adlaire-ci-api.service（常駐）
  └─ adlaire-ci-api（components/api.go）→ SDK → 管理ツール
```

Go 標準ライブラリと GitHub PAT（`contents: read`）を基本要件とする。TLS 終端に nginx 等のリバースプロキシを使う場合でも、Adlaire CI 本体は HTTP サーバーとして実装する。

---

## 0j. リポジトリ内ソース配置

Adlaire CI の標準リポジトリ内ソース配置は以下とする。

本節は、移行後の標準配置を定義する。標準配置への実装移行が完了するまでは、現行リポジトリに `build_spec.go`、`runner.go`、`build_spec_test.go`、`runner_test.go`、`testdata/build_spec/` が残る場合がある。実装移行 PR では、本節の配置へそろえ、移行後に旧配置を残してはならない。

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

| パス | 役割 |
|------|------|
| `main.go` | 起動入口。サブコマンド判定、引数受け取り、責務 component 呼び出しを行う。 |
| `components/builder.go` | Markdown / Markdown ディレクトリを静的 Web サイトへ変換する。 |
| `components/runner.go` | GitHub polling、変更検出、ビルド起動、履歴、ログ、deploy を実行する。 |
| `components/api.go` | 管理 API サーバー、認証、状態ファイル操作を提供する。 |
| `components/admin.go` | 管理 UI 静的ファイルの配布・配置を扱う。 |
| `components/statefile.go` | `.build_history`、`.build_logs`、`.server_config` など状態ファイルの読み書きを扱う。詳細は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` を正とする。 |
| `components/archive.go` | ビルドログ圧縮、snapshot、配布アーカイブを扱う。 |
| `components/commitstatus.go` | GitHub Commit Status API 送信を扱う。 |
| `components/mcp.go` | MCP 接続を扱う。 |
| `admin/` | 標準管理 UI の静的ファイルを配置する。 |
| `testdata/` | コンポーネント別 fixture を配置する。 |
| `docs/examples/` | 利用例、設定例、サンプル構成を配置する。 |

`main.go` は 1 ファイルとし、実装詳細を含めない。`components/` 配下は 1 コンポーネント = 1 Go ファイルとし、各ファイルは上表の責務を実装する。

---

## 1. Builder 詳細仕様（分割済み）

`builder` owner component の詳細仕様本文は `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` を正とする。

| 移動済み範囲 | 移動先 |
|--------------|--------|
| §1〜§9 | `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §1〜§9 |
| §8a | `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §8a |
| §27.4 | `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §27.4 |
| §27.25 | `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §27.25 |
| §27.28 | `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §27.28 |

### — CI ランナー —

## 10. Runner 詳細仕様（分割済み）

`runner` owner component の詳細仕様本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` を正とする。

| 移動済み範囲 | 移動先 |
|--------------|--------|
| §10〜§20 | `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §10〜§20 |
| §15a | `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §15a |
| runner owner の §27 個別節 | `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27 |

`ADLAIRE_CI_DETAIL_SPEC.md` §27.38a は横断補足契約として本ファイルに残す。

§27.21〜§27.38 の runner 拡張機能を実装する場合は、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` の個別節を正本とし、横断する処理順、状態ファイル保存責務、API / SDK / UI 連動条件、受け入れ fixture は本ファイル §27.38a を同時に確認する。

### — 管理ツール —

## 21. API 詳細仕様（分割済み）

`api` owner component の詳細仕様本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` を正とする。

| 移動済み範囲 | 移動先 |
|--------------|--------|
| §21〜§22 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §21〜§22 |
| §25 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §25 |
| api owner の §27 個別節 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §27 |

`ADLAIRE_CI_DETAIL_SPEC.md` §27.38a は横断補足契約として本ファイルに残す。

§27.21〜§27.38 に関わる API / SDK / UI 連動条件は、各 owner component の分割先詳細仕様ファイルと本ファイル §27.38a を同時に満たす。§27.38a の内容を API / SDK / UI の分割先へ重複定義してはならない。

## 23. SDK 詳細仕様（分割済み）

`sdk` owner component の詳細仕様本文は `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 を正とする。

| 移動済み範囲 | 移動先 |
|--------------|--------|
| §23 | `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 |

## 24. UI 詳細仕様（分割済み）

`ui` owner component の詳細仕様本文は `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 を正とする。

| 移動済み範囲 | 移動先 |
|--------------|--------|
| §24 | `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 |

## 25. API 認証詳細仕様（分割済み）

`api` owner component の認証詳細仕様本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` §25 を正とする。

## 26. Setup 詳細仕様（分割済み）

`setup` owner component の詳細仕様本文は `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 を正とする。

| 移動済み範囲 | 移動先 |
|--------------|--------|
| §26 | `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 |

## 27. 追加仕様化機能 詳細仕様

本節は、`ADLAIRE_CI_DETAIL_SPEC.md` に定義済みの追加機能の詳細仕様である。本節に定義された機能は、§0i、§12、§13、§15、§22、§23、§24 と同時に満たす。

### 27.0 追加仕様化機能 共通実装契約

本節の §27.1〜§27.38、§27.42〜§27.47 は、各個別節に別指定がない限り、以下の共通実装契約を満たす。

**実装対象判定：**

| 項目 | 仕様 |
|------|------|
| 対象節 | §27.1〜§27.38、§27.42〜§27.47。 |
| 対象外 | 本ファイルに詳細節が存在しない機能、`ADLAIRE_CI_SPEC.md` で実装不可に分類される機能、MCP 専用機能。 |
| 実装単位 | 個別節単位で実装する。ただし API、SDK、UI、状態ファイル schema、検証条件が同一機能に含まれる場合は同一実装 PR 内でそろえる。 |
| 補完禁止 | 個別節に存在しない endpoint、状態ファイル、設定 key、UI 操作、SDK method を実装判断で追加してはならない。 |

**共通処理順序：**

| 順序 | 処理 | 失敗時 |
|------|------|--------|
| 1 | CLI / HTTP path / method / body size / JSON parse を検証する。 | §8、§13、§22.0 の終了コードまたは HTTP status を返す。 |
| 2 | 認証、scope、maintenance、rate limit、lock 競合を判定する。 | 対象処理を開始せず固定エラーを返す。 |
| 3 | 入力値を schema と個別節の許容値で検証する。 | 状態ファイル、外部 API、外部 command を変更しない。 |
| 4 | 読取対象状態ファイルを読み、破損時処理を個別節または §22.0a に従って行う。 | 自動復旧が定義されていない場合は `500` または終了コード `1`。 |
| 5 | 外部呼び出し前に build id、request id、queue id などの id を採番する。 | 採番不能は `500` または終了コード `1`。 |
| 6 | 個別節の主要処理を実行する。 | 個別節の異常系に従う。 |
| 7 | 状態ファイル、ログ、履歴、通知、監査を個別節の保存順で更新する。 | 保存順が未定義の場合は下表の保存順を使う。 |
| 8 | API response、stdout、UI 表示用値を返す。 | response 生成不能は `500`。 |

**保存順が未定義の場合の既定保存順：**

| 種別 | 保存順 |
|------|--------|
| runner build 成功 | `.build_logs/{id}.json` → `.build_history` → `.build_status.json` → SHA cache / local watch state → notification。 |
| runner build 失敗 | `.build_logs/{id}.json` → `.build_history` → `.build_status.json` → notification。SHA cache は更新しない。 |
| runner skip | `.build_status.json` のみ。個別節が明示しない限り `.build_logs/{id}.json` と `.build_history` は作成しない。 |
| API 設定変更 | 対象状態ファイル → `.config_log` → `.audit_log` → response。 |
| API token / auth / security | 対象状態ファイル → `.access_log` → `.audit_log` → response。 |
| Webhook / queue | イベントログ → `.build_state.queued` → response。個別節で queue 先行が明記される場合は個別節を優先する。 |

状態更新は、§22.0a の lock、atomic write、fsync、rename、親ディレクトリ fsync を使用する。JSON Lines の追記は 1 行 1 object、UTF-8、LF、末尾改行必須とする。追記対象ファイルの親ディレクトリが存在しない場合は、個別節で作成可と明記されている場合だけ作成する。

**共通エラー優先順位：**

| 優先 | 条件 | HTTP / 終了コード |
|------|------|-------------------|
| 1 | path 不存在、method 不一致 | `404` / `405` |
| 2 | body 禁止、body size 超過、JSON parse 失敗 | `400` / `413` |
| 3 | 認証なし、token 不正、session 期限切れ | `401` |
| 4 | scope 不足、管理操作不可 | `403` |
| 5 | rate limit 超過、queue full | `429` |
| 6 | 入力 schema、範囲、enum、path 検証失敗 | `422` |
| 7 | 状態競合、二重実行、未準備状態 | `409` |
| 8 | 必須外部設定なし | `501`、個別節が `422` または `503` を指定する場合は個別節優先 |
| 9 | 状態ファイル read/write、外部 command、外部 API の処理失敗 | `500` または個別節の終了コード |

複数条件が同時に成立する場合は、上表の上位を返す。ただし `POST /api/webhook` の署名検証失敗は、情報漏えいを避けるため JSON parse より前に `401` を返してよい。

**共通データ制約：**

| 対象 | 仕様 |
|------|------|
| 時刻 | UTC ISO 8601 `YYYY-MM-DDTHH:MM:SSZ`。保存時にミリ秒と timezone offset は使わない。 |
| id | 個別節に定義がない場合は `^[A-Za-z0-9_-]{1,64}$`。path separator、`.`、`..`、空文字は禁止。 |
| path | API request の path 値は、個別節に絶対 path と明記したもの以外は相対 path とし、`..`、NUL、改行を禁止する。 |
| 並び順 | API 一覧は個別節に明記がない限り、新しい順、同時刻は id 昇順。 |
| 文字列上限 | message、reason、error、description は個別節に指定がない限り 500 文字。改行は `\n` 文字列へ escape する。 |
| secret | `password`、`token`、`secret`、`pat`、`smtp_password` を key 名に含む値は、response、log、history、audit、backup、UI 表示で `"***"` に mask する。 |

**API / SDK / UI 同期契約：**

| 項目 | 仕様 |
|------|------|
| API | §22.0e に endpoint が存在する機能だけを実装対象とする。endpoint 追加が必要な場合は、先に §22.0d、§22.0e、§23、§24 を同時に更新する。 |
| SDK | SDK method は §22.0e の SDK 列と §23 の引数変換契約だけに従う。API response を推測補完しない。 |
| UI | UI 操作は §24 UI 操作契約表に存在する SDK method だけを呼ぶ。直接 `fetch()`、状態ファイル操作、外部 command 実行を行わない。 |
| 状態ファイル | §22.0a と §22.0c に存在しない状態ファイルを作成しない。必要な場合は本ファイル内で schema と破損時処理を先に定義する。 |

**§27 component 責務ベース読解契約：**

§27 の各機能は、機能番号ではなく component 責務を先に確定してから読む。実装者は、対象機能の個別節を読む前に下表で owner component、参照節、越境禁止事項を確認する。owner component が未確定のまま実装、fixture、PR 分割を開始してはならない。

| component | 正本責務 | 主な参照節 | collaborator | 越境禁止 |
|-----------|----------|------------|--------------|----------|
| `builder` | Markdown 入力から静的 Web サイト、asset、REPORT、HTML meta、dependency / cache 出力を生成する。 | §2〜§8、§27.4、§27.25、§27.28 | `runner`、`api` | GitHub read、状態ファイル直接更新、API response 補完、通知送信。 |
| `runner` | CLI、設定、GitHub / local read、差分判定、lock、build 実行、dry-run、retry、queue、status / log / history / notification を制御する。 | §10〜§16、§27.1〜§27.4、§27.21〜§27.38 | `builder`、`statefile`、`commitstatus` | API endpoint 追加、UI 操作追加、SDK method 追加、状態 schema の暗黙追加。 |
| `api` | HTTP endpoint、auth / scope、validation、状態 read/write、response schema、access / audit / config log を提供する。 | §22、§25、§27.5〜§27.20、§27.42〜§27.47 | `statefile`、`sdk`、`security` | runner 専用 CLI 処理、UI DOM 操作、SDK 側補完前提の response 欠落。 |
| `sdk` | API endpoint を仕様どおり呼び出す JavaScript method、query / body 生成、error 伝播、binary / stream 処理を提供する。 | §23、§27 API / SDK / UI 連動固定契約 | `api`、`ui` | API response の推測補完、未定義 endpoint 呼び出し、状態ファイル直接操作。 |
| `ui` | SDK method を使った管理画面表示、入力、loading / disabled / error / secret field 消去、再取得順を提供する。 | §24、§27 API / SDK / UI 連動固定契約 | `sdk`、`api` | 直接 `fetch()`、状態ファイル操作、外部 command 実行、secret DOM 残存。 |
| `statefile` | 状態 schema、atomic write、JSON Lines、lock、破損時処理、保存順、no-write / forbidden write を固定する。 | §22.0a、§22.0c、§27 fixture / effects 契約 | `runner`、`api`、`archive` | component 固有の業務判断、API response 補完、UI 表示判断。 |
| `archive` | build log archive、snapshot、download、delete、rollback の保存 / 取得 / 削除境界を固定する。 | §14b、§27.7、§27.15 | `runner`、`api`、`statefile` | build 成否の反転、未検証 archive 展開、元 build log の改変。 |
| `commitstatus` | GitHub Commit Status の pending / final payload、送信順、失敗時非反転、secret mask を固定する。 | §27.1、§13、§15 | `runner`、`statefile` | build 実行判断、GitHub write API の dry-run 実行、status 失敗による build 成否反転。 |
| `security` | API token、session、TOTP、scope、audit、rate limit、secret mask、漏えい禁止値を固定する。 | §25、§27.42〜§27.47 | `api`、`sdk`、`ui`、`statefile` | runner / builder の業務処理代行、secret 平文保存、認可前の状態更新。 |
| `setup` | release binary 配置、systemd、初期化、update、rollback、既存 secret 保持を固定する。 | `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 | `runner`、`api`、`statefile` | runtime 機能追加、状態 schema の暗黙変更、外部依存追加。 |

component 責務の読み順は、`component 正本責務` → `個別機能節` → `状態 / API / SDK / UI 連動表` → `fixture / effects / security expected` とする。PR 本文、fixture manifest、受け入れ checklist はこの順序で component を記録する。component 間の共通処理は `statefile`、`security`、`setup` のいずれかの責務として明記し、横断共通基盤や core という別 component を作って扱ってはならない。

**共通検証完了条件：**

| 検証 | 合格条件 |
|------|----------|
| scope | 変更対象が `ADLAIRE_CI_DETAIL_SPEC.md` の既存詳細機能に限定されている。 |
| schema | 状態ファイル schema、API response、SDK 型、UI 表示が同じ key 名と nullable 条件で一致する。 |
| success | 正常系 fixture が、保存順、response、ログ、履歴、通知、監査の期待値をすべて満たす。 |
| failure | 異常系 fixture が、状態差分なしまたは定義済み部分更新だけで終了する。 |
| secret | token、password、secret、PAT、TOTP secret、session token が response、log、history、audit、backup、UI に平文で出ない。 |
| idempotency | 同じ GET、同じ dry-run、変更なし保存、同じ cleanup/archive 対象なしの再実行で追加差分が出ない。 |
| conflict | lock 競合、queue full、running build、破損 state の結果が固定 status / 終了コードで再現できる。 |

**§27 実装 PR 完了判定固定契約：**

§27 の機能を実装した PR は、下表をすべて満たした場合だけ完了扱いにする。下表は §27 に既に定義済みの機能だけへ適用し、将来計画、MCP、外部公開構成、未定義 endpoint、未定義 UI、未定義状態ファイルを追加する根拠にしてはならない。

| 判定項目 | 必須条件 | 未達時の扱い |
|----------|----------|--------------|
| 節対応 | 実装対象の §27.x、関連する §22.0d、§22.0e、§23、§24、§25、§26 の参照箇所を PR 本文または検証ログで明示する。 | 対応仕様不明として未完了。 |
| 入力固定 | CLI 引数、HTTP method/path/query/body/header、SDK 引数、UI field、状態ファイル入力、環境変数入力のいずれか該当するものを fixture に含める。 | 入力契約不足として未完了。 |
| 出力固定 | stdout/stderr、HTTP response、SDK return/error、UI 表示、状態ファイル差分、ログ、履歴、通知、監査の該当出力を期待値として保存する。 | 出力契約不足として未完了。 |
| 正常系 | 最小正常値、既定値、境界内最大値、no-op、重複実行、既存状態ありの正常系を少なくとも 1 件ずつ該当範囲で検証する。 | 正常系不足として未完了。 |
| 異常系 | validation error、認証/権限 error、lock/queue conflict、状態破損、外部 API/command 失敗、保存失敗の該当ケースを検証する。 | 異常系不足として未完了。 |
| 副作用 | 失敗時に書き換えてよいファイル、書き換えてはならないファイル、外部送信有無、再実行時の残留状態を fixture で比較する。 | 副作用不明として未完了。 |
| secret | request、response、stdout/stderr、build log、history、audit、access log、notify log、backup、UI DOM に secret 平文が存在しないことを検証する。 | secret 漏えいリスクとして未完了。 |
| 順序 | 複数状態を更新する機能は、§27.0 の保存順または個別節の保存順を fixture 名または検証ログに明記する。 | 部分失敗時挙動不明として未完了。 |
| 回帰 | 実装対象外 endpoint、実装対象外 SDK method、実装対象外 UI 操作、MCP、外部公開 bind が追加されていないことを確認する。 | 範囲逸脱として未完了。 |

**§27 fixture / PR 証跡詳細契約：**

§27 の fixture 配置、fixture カタログ、manifest schema、assertion、expected/effects、相互整合、component 別検証責務、PR 証跡、受け入れゲート、差し戻し条件、部分失敗・再実行契約は `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §27-F を正とする。本節には §27 の実装順序、PR 分割、個別機能の処理本文だけを置き、fixture 共通契約を重複定義しない。

**§27 実装順序・PR 分割固定契約：**

§27 の機能は、下表の wave 順に実装する。後続 wave は、先行 wave の対象機能が §27 実装 PR 最終受け入れゲートを満たすまで開始してはならない。下表は §27 に既に記載済みの機能だけを対象とし、将来計画、MCP、外部公開構成、未定義 endpoint、未定義 UI、未定義状態ファイルを追加する根拠にしてはならない。

| wave | 対象節 | 実装単位 | 依存条件 | PR 上限 | 完了時固定契約 |
|------|--------|----------|----------|---------|----------------|
| W1 | §27.1〜§27.4 | runner / builder の基礎拡張。commit status、dry-run、retry、output meta。 | §8〜§20 の builder / runner 基本契約が完了している。 | 1 PR につき 1〜2 機能。 | build log、history、status、REPORT、HTML meta、外部 GitHub call の fixture が固定される。 |
| W2 | §27.5〜§27.11 | API / runner 運用基盤。config validation、access log、archive、status、trigger、startup integrity、schedule。 | W1 完了。API 共通契約 §22.0 が固定済み。 | 1 PR につき 1〜3 機能。 | 設定、状態 read/write、access/config log、systemd fake の fixture が固定される。 |
| W3 | §27.12〜§27.20 | webhook、stats、snapshot、health、search、branch config、summary、config diff。 | W2 完了。状態ファイルと API read/write 境界が固定済み。 | 1 PR につき 1〜3 機能。 | webhook event、queue、snapshot、stats、branch config、notify summary、config diff の fixture が固定される。 |
| W4 | §27.21〜§27.29 | runner / builder 拡張前半。複数監視、YAML pipeline、local watch、tag filter、cache、parallel targets、hooks、dependency、remote build。 | W1〜W3 完了。runner 状態、API 設定、fixture fake が固定済み。 | 1 PR につき 1〜2 機能。 | target selection、pipeline、cache、hook、dependency、remote artifact の成功 / 失敗 / partial fixture が固定される。 |
| W5 | §27.30〜§27.38 | runner / operation 拡張後半。approval、branch env、notification、trend、chain、priority、failure category、environment、duration anomaly。 | W4 完了。queue、notification、history、trend の状態境界が固定済み。 | 1 PR につき 1〜2 機能。 | approval、env、notify、trend、chain、queue、category、environment、anomaly の保存順と再実行 fixture が固定される。 |
| W6 | §27.42〜§27.47 | security / auth / audit / rate limit。trigger scope、API key、audit、session timeout、TOTP、rate limit。 | W2〜W5 完了。API endpoint、audit/access log、secret mask 境界が固定済み。 | 1 PR につき 1〜2 機能。 | scope、token、audit、session、TOTP、rate limit の security fixture と forbidden side effect が固定される。 |

**§27 W1-PR1 実装対象固定契約：**

W1 の最初の実装 PR は `§27.2 ドライラン実行モード` だけを対象とする。`§27.1`、`§27.3`、`§27.4` は W1 内の後続 PR とし、W1-PR1 で同時実装してはならない。W1-PR1 は副作用なしで runner の設定読取、GitHub 読取、差分判定、終了コード、stdout JSON、secret mask を固定するための PR とする。

| 項目 | 固定内容 |
|------|----------|
| 対象節 | `§27.2` のみ。 |
| owner component | `runner`。必要な fake GitHub response と fixture harness は含める。 |
| collaborator component | `statefile`。状態ファイル、lock、log、history、SHA cache、notification、deploy、snapshot の no-write / forbidden write を検証する。 |
| 対象外 | `§27.1` commit status、`§27.3` retry、`§27.4` output meta、API、SDK、UI、MCP、外部公開構成、状態ファイル schema 新設。 |
| 許可される副作用 | なし。dry-run は状態ファイル、lock、SHA cache、build log、history、notification、deploy、snapshot、GitHub Commit Status を作成、更新、削除してはならない。 |
| 許可される外部呼び出し | fake GitHub read だけ。実 GitHub API、GitHub write API、SSH、pipeline、notification、systemd、hook は呼び出さない。 |
| 完了後に固定される契約 | dry-run stdout JSON schema、終了コード、secret mask、GitHub read failure、設定破損時 no-write、複数 target 表示、cooldown 表示。 |

**§27 W1-PR1 必須 fixture 固定契約：**

W1-PR1 は、下表の fixture をすべて含める。1 件でも不足、skip、期待副作用不足、secret mask 不足がある場合、W1-PR1 は未完了とする。

| fixture 名 | 必須 assertion | 入力 | 期待結果 |
|------------|----------------|------|----------|
| `success-dry-run-changed` | `stdout`、`state`、`logs`、`effects`、`no-write`、`secret-mask` | `--dry-run`、変更あり fake GitHub response、既存 SHA cache、対象 branch / target。 | stdout JSON に `would_build=true`、対象 commit/blob、`would_write` 一覧、状態 / logs 差分なし、fake GitHub read 以外の外部 call 0。 |
| `noop-dry-run-unchanged` | `stdout`、`state`、`logs`、`effects`、`no-write`、`idempotency` | `--dry-run`、変更なし fake GitHub response、既存 SHA cache。 | stdout JSON に `would_build=false`、`reason="no_change"`、状態 / logs 差分なし、2 回実行して差分なし。 |
| `failure-dry-run-github-error` | `stdout`、`stderr`、`state`、`logs`、`effects`、`no-write`、`order`、`secret-mask` | `--dry-run`、fake GitHub `500` または rate limit failure。 | 終了コード `3`、状態 / logs 差分なし、pipeline / deploy / notification 呼び出し 0、error detail は secret を含まない。 |
| `security-dry-run-secret-mask` | `stdout`、`stderr`、`state`、`logs`、`effects`、`secret-mask`、`no-write` | token、webhook URL、notification URL、branch config secret を含む入力。 | stdout/stderr/effects/log 期待値に secret 平文が存在せず、mask 後値だけを含み、状態 / logs 差分なし。 |

**§27 W1-PR1 acceptance checklist：**

W1-PR1 の PR 本文には、下表を記録する。記録がない項目は未検証として扱い、W1-PR1 を完了扱いにしてはならない。

| 項目 | 必須記録 |
|------|----------|
| 対象仕様 | `§27.2`、関連する `§12`、`§13`、`§15`、`§22.0a`、`§27 fixture` 契約。 |
| 対象外 | `§27.1`、`§27.3`、`§27.4`、API、SDK、UI、MCP、外部公開構成、状態ファイル schema 新設。 |
| fixture | `success-dry-run-changed`、`noop-dry-run-unchanged`、`failure-dry-run-github-error`、`security-dry-run-secret-mask`。 |
| 副作用確認 | 状態ファイル、lock、SHA cache、build log、history、notification、deploy、snapshot、commit status に差分がないこと。 |
| effects 確認 | `expected/effects.json` の `external_calls` には fake GitHub read だけを記録し、`commands`、`notifications`、`downloads`、`streams` は空配列にする。`unchanged_paths` と `forbidden_writes` には状態ファイル、lock、SHA cache、build log、history、pending、snapshot、notify を列挙する。 |
| 外部呼び出し確認 | fake GitHub read 以外の呼び出しが 0 件であること。`forbidden_calls` には GitHub write API、SSH、pipeline、notification、systemd、hook、commit status を列挙する。 |
| secret 確認 | token、password、secret、PAT、Authorization header が stdout/stderr/effects/expected に平文で存在しないこと。 |
| 後続影響 | W1 後続 PR が利用してよい dry-run JSON schema、終了コード、no-write 契約。 |

**§27 W1-PR1 manifest / effects 整合固定契約：**

W1-PR1 の全 fixture は、§27 fixture manifest schema 固定契約と §27 expected/effects.json schema 固定契約に加えて、下表を満たす。

| 対象 | 固定内容 |
|------|----------|
| `manifest.json.section` | `27.2` に固定する。 |
| `manifest.json.feature` | `dry_run` に固定する。 |
| `manifest.json.owner_component` | `runner` に固定する。 |
| `manifest.json.collaborator_components` | `["statefile"]` に固定する。dry-run no-write を検証するために statefile 責務を collaborator として記録する。 |
| `manifest.json.components` | `["runner","statefile"]` に固定する。API、SDK、UI、builder、archive、commitstatus、setup、security を含めてはならない。 |
| `manifest.json.references` | `§27.2`、`§12`、`§13`、`§15`、`§22.0a`、`§27 fixture` を含める。 |
| `manifest.json.not_applicable` | `input/request.json`、`expected/response.json`、UI DOM、SDK return、download、stream が該当しない理由を記録する。 |
| `manifest.json.missing_state` | fixture ごとに実行前に存在しない状態ファイルだけを列挙する。存在する状態を missing として扱ってはならない。 |
| `expected/effects.json.external_calls` | fake GitHub read のみを記録する。変更なし fixture でも GitHub read を実行する場合は 1 件記録する。GitHub read を行わない fixture は理由を `manifest.json.not_applicable` に記録する。 |
| `expected/effects.json.commands` | 常に空配列。pipeline、ssh、systemd、hook、archive、setup/update を含めてはならない。 |
| `expected/effects.json.notifications` | 常に空配列。通知送信、pending 化、retry 対象を含めてはならない。 |
| `expected/effects.json.write_order` | 常に空配列。dry-run は状態、log、history、SHA cache、lock を書かない。 |
| `expected/effects.json.created_paths` / `deleted_paths` | 常に空配列。dry-run は file / directory を作成、削除しない。 |
| `expected/effects.json.unchanged_paths` | 状態ファイル、lock、SHA cache、build log、history、pending、snapshot、notify、deploy 出力を列挙する。 |
| `expected/effects.json.forbidden_writes` | `unchanged_paths` と同じ対象に加え、`.build_status.json`、`.build_state`、`.build_history`、`.build_logs/`、`.last_sha`、`.notify_pending`、`.pending_transfers` を列挙する。 |
| `expected/effects.json.forbidden_calls` | GitHub write API、GitHub Commit Status、SSH、pipeline、deploy、snapshot、notification、systemd、hook、remote build を列挙する。 |
| `expected/security.json` | 禁止文字列として token、Authorization header 値、PAT、password、secret、webhook URL credential、notification URL credential を含める。 |

W1-PR1 fixture の `manifest.json`、`expected/effects.json`、`expected/security.json` が上表を満たさない場合、fixture は存在していても未完了とする。

**§27 W1-PR1 stdout / exit code 整合固定契約：**

W1-PR1 の dry-run stdout は、`§27.2 ドライラン実行モード` の JSON schema と完全一致させる。fixture expected が `§27.2` と異なる key 名、reason 値、終了コード、no-write 判定を持つ場合は、fixture ではなく本仕様を修正対象とする。

| 対象 | 固定内容 |
|------|----------|
| stdout JSON | JSON object 1 件と末尾 LF だけを出力する。通常ログ、進捗行、secret 平文、追加 JSON 行を混在させてはならない。 |
| `would_write` | 実書込結果ではなく、非 dry-run なら発生する論理書込予定だけを列挙する。dry-run 実行中の実ファイル差分は常に 0 件でなければならない。 |
| `expected/effects.json.write_order` | dry-run の実書込順序であるため常に空配列にする。`stdout.would_write` と混同してはならない。 |
| GitHub read 失敗 | fake GitHub read の最終失敗、rate limit、network failure は終了コード `3`、`reason="github_error"`、`errors[]` 出力、状態 / log / 通知 / status / deploy 差分なしに固定する。 |
| 設定破損 | 起動時設定検証で破損を検出した場合は終了コード `2`、`reason="config_error"`、`errors[]` 出力、backup / 初期化 / 正規化 / quarantine / rewrite なしに固定する。 |
| 後続 PR 制約 | W1 後続 PR は `§27.2` の stdout key、reason 値、終了コード、no-write、fake GitHub read、secret mask を変更してはならない。 |

W1-PR1 完了後、W1 内の後続 PR は `§27.1`、`§27.3`、`§27.4` のいずれか 1〜2 機能を対象にできる。ただし W1-PR1 の dry-run no-write 契約、stdout JSON schema、終了コード、secret mask 契約、fake GitHub read 契約を変更してはならない。変更が必要な場合は、W1-PR1 の仕様改訂として本節と `§27.2` を先に更新する。

**§27 PR 分割禁止条件：**

| 条件 | 扱い |
|------|------|
| 1 PR で 4 機能以上を実装する。 | レビュー不能として未完了。W2 / W3 の軽量 read-only API を含む場合でも最大 3 機能までとする。 |
| 先行 wave の未完了機能に依存する後続 wave を実装する。 | 順序違反として未完了。 |
| API だけ、SDK だけ、UI だけを先行し、状態 schema または fixture を同一 PR で固定しない。 | component 責務不足として未完了。 |
| fixture カタログにない機能を便宜的に同梱する。 | 仕様外実装として差し戻し。 |
| 既存 fixture の期待値を弱めて新機能を通す。 | 検証の形骸化として差し戻し。 |
| 複数 wave にまたがる横断 refactor を主目的にする。 | §27 機能実装 PR として扱わず、先に仕様改訂が必要。 |

**§27.1〜§27.20 機能別実装完全性固定契約：**

§27.1〜§27.20 の各機能は、個別節の本文に加えて下表を満たした場合だけ実装完了とする。下表は実装対象の入力、出力、状態、失敗時副作用、fixture を固定するための補助契約であり、個別節と矛盾する場合は個別節のより具体的な値を優先する。

| 節 | 機能 | 入力 | 出力 | 状態ファイル / 外部副作用 | 失敗時副作用 | 必須 fixture |
|----|------|------|------|---------------------------|--------------|--------------|
| §27.1 | GitHub Commit Status API | `.server_config.commit_status_*`、commit SHA、build 結果。 | GitHub status payload、`.build_logs.{commit_status}`、`.build_history.commit_status_state`。 | GitHub Status API 送信、build log / history 追記。 | status 送信失敗で build 成否を反転しない。token、response body 全体を保存しない。 | pending→success、pending→failure、commit SHA なし、pending 失敗、final 失敗、無効時呼び出し 0。 |
| §27.2 | dry-run | CLI option、runner 設定、fake GitHub read 結果。 | stdout JSON 1 object、終了コード。 | 状態ファイル、lock、通知、deploy、status API を変更しない。 | 設定破損でも退避 / 再生成しない。fake GitHub read 失敗は終了コード `3`。 | SHA 差分、差分なし、cooldown、fake GitHub read 失敗、設定破損、複数 target、secret 非表示。 |
| §27.3 | 自動リトライ | retry 設定、失敗種別、attempt 結果。 | `.build_logs.attempts[]`、`.build_history.retry_count`。 | retry 対象だけ再試行。最終成功時だけ SHA / snapshot / deploy success を確定。 | 非 retry 対象では再試行しない。中断時は未実行 attempt を作らない。 | API 429 後成功、pipeline timeout 後成功、pipeline exit 1、deploy checksum mismatch、上限到達、中断。 |
| §27.4 | build meta | builder CLI 引数または環境変数。 | HTML meta、`[REPORT]`、build log、`GET /api/output-meta`。 | 出力 HTML と report 生成。 | 不正値は builder 終了コード `2`、既存出力を成功扱いしない。 | 値あり、空値、複数ページ、不正 SHA、HTML / REPORT / API 一致。 |
| §27.5 | config validate API | partial config JSON。 | `valid`、正規化 config、`errors[]`、`warnings[]`。 | `.api_access_log` 以外を変更しない。 | unknown key は `422`。validation error は `200 valid=false`。 | valid true、valid false、unknown key、secret key、破損 config、no write。 |
| §27.6 | API access log | 全 `/api/` request の method/path/status/duration/actor。 | `.api_access_log` JSON Lines、参照 API response。 | request ごとに追記。secret と body は記録しない。 | access log 追記失敗は対象 API の個別規定に従い、secret を出さない。 | health、認証失敗、認証成功、422、500、secret body、pagination。 |
| §27.7 | log archive / cleanup | retention 設定、対象 build log。 | gzip archive、cleanup / archive response。 | 対象 log の圧縮 / 削除、disk usage 更新対象。 | gzip 失敗時は元 log を維持。削除対象なしは no-op。 | archive 成功、対象なし、gzip 失敗、破損 log skip、cleanup 0 日。 |
| §27.8 | build status file | runner 状態、直近 build、queue、pending、circuit。 | `.build_status.json`、`GET /api/status` 第一参照値。 | runner finalizer で atomic write。 | status write 失敗は build log/history の結果を壊さない。 | running、success、failure、pending transfer、circuit open、破損時 API 500。 |
| §27.9 | trigger 種別 | runner 起動理由、API / webhook / approval / rollback。 | `.build_logs.trigger`、`.build_history.trigger`、API filter。 | trigger を固定 enum で保存。 | 未知 trigger は保存しない。filter 不正は `422`。 | polling、manual、webhook、rollback、approval、unknown filter。 |
| §27.10 | 起動時整合性チェック | runner 起動時の設定 / 状態ファイル群。 | recovery record、WARN / ERROR、必要時 status。 | 個別節で許可された復旧だけ実行。 | 認証 / API が読む状態を勝手に初期化しない。復旧不能は build 開始しない。 | 正常、破損復旧、復旧不能、secret mode 補正、補正失敗。 |
| §27.11 | 動的 schedule | interval / pause / resume / allowed hours API。 | schedule response、`.server_config`、systemd timer 更新結果。 | config 保存後に systemd 更新、config log 追記。 | systemd 失敗時は `500`。保存済み config は巻き戻さず再取得で確定値を返す。 | interval 成功、pause/resume、allowed hours、systemd 失敗、no-op。 |
| §27.12 | Webhook 受信 | GitHub headers、署名、payload。 | queue entry、event log、HTTP response。 | 署名検証後に event 保存、必要時 queue 追加。 | 署名失敗は状態差分なし。重複 delivery は二重 queue しない。 | valid push、invalid signature、duplicate delivery、unsupported event、branch mismatch、queue full。 |
| §27.13 | Webhook event log | `.webhook_events.json`、limit/offset。 | `events[]`、`total`。 | GET は read-only。Webhook 受信時だけ追記。 | 壊れた JSON Lines 行は除外し、response に破損内容を出さない。 | paging、empty、corrupt line skip、limit 不正、secret 非表示。 |
| §27.14 | duration / stats | build history、build logs、days/n query。 | stats、timeline、duration stats response。 | GET は read-only。 | 破損 log は除外し WARN。query 不正は `422`。 | count 0、成功/失敗混在、archive 含む、duplicate id、avg rounding、query invalid。 |
| §27.15 | artifact snapshot | `.snapshots/{id}`、history id、running 状態。 | snapshot list、tar.gz download、delete response、rollback history/log。 | download read-only、delete は対象 snapshot 削除、rollback は新規履歴作成。 | 不正 id / running は差分なし。delete log 失敗は削除済みのまま `500`。 | list、download、安全でない entry、delete、rollback success、rollback pending、running conflict。 |
| §27.16 | health | `.build_status.json`、pending / notify state、process uptime。 | HTTP 200 health JSON。 | read-only。 | 破損状態は自動修復せず degraded。response 生成不能だけ `500`。 | ok、pending、status missing、status corrupt、read error、stale runner。 |
| §27.17 | log severity filter | q/from/to/level、build logs。 | search results with level/source/line_number。 | read-only。 | 不正 level は `422`。破損 log は除外。 | warn、error、debug、archive、stderr、line number、不正 level。 |
| §27.18 | branch config API | `branches[]` request、既定 branch target。 | `.branch_config.branch_targets`、GET response、config log。 | POST valid は保存、POST empty は削除。runner は次回起動から反映。 | validation 失敗は差分なし。log 失敗時は保存済み状態を戻さない。 | GET default、POST valid、POST empty、重複 branch、path 不正、log failure。 |
| §27.19 | weekly summary | `.notify_config.summary`、history、channel 設定。 | webhook payload、notify log、手動 API response。 | 自動成功時だけ sent date 更新。手動は sent date 更新なし。 | 宛先なしは自動 WARN / 手動 `422`。送信失敗は sent date 更新なし。 | 条件一致、同日二重、宛先なし、手動送信、送信失敗、skipped 除外。 |
| §27.20 | config diff log | 設定変更 API の before / after。 | `.config_log` diff / diff_text。 | 対象状態保存後に config log 追記。 | no-op は状態 / log 差分なし。config log 失敗は保存済み状態を戻さず `500`。 | simple diff、nested diff、array diff、secret mask、no-op、log failure、target mapping。 |

**§27.1〜§27.20 API / SDK / UI 連動固定契約：**

| 節 | API | SDK | UI |
|----|-----|-----|----|
| §27.1 | `GET/POST /api/config` の commit status key と runner 結果に反映する。 | `getConfig()` / `setConfig()` で key を削除しない。 | 設定 panel で有効/無効、context、target URL を表示 / 保存する。 |
| §27.2 | API endpoint は追加しない。 | SDK method は追加しない。 | UI 操作は追加しない。dry-run は CLI 検証だけ。 |
| §27.3 | `GET/POST /api/config` の retry key と build log / history に反映する。 | `getConfig()` / `setConfig()` が retry key をそのまま扱う。 | 設定 panel で retry 最大回数と base 秒数を表示 / 保存する。 |
| §27.4 | `GET /api/output-meta` で HTML meta と同値を返す。 | `getOutputMeta()` が response を補完せず返す。 | システム情報 panel に build id、commit、build at を表示する。 |
| §27.5 | `POST /api/config/validate` は保存しない。 | `validateConfig(config)` は `valid/errors/warnings` をそのまま返す。 | validate 結果を保存成功と混同せず表示する。 |
| §27.6 | `GET /api/api-access-log` で paging/filter する。 | `getApiAccessLog()` は query を仕様順で送る。 | アクセスログ panel に API access log filter を表示する。 |
| §27.7 | `POST /api/logs/archive`、`POST /api/logs/cleanup`。 | `archiveLogs()` / `cleanupLogs()`。 | ログビューアで件数を表示し、disk usage を再取得する。 |
| §27.8 | `GET /api/status` の第一参照元にする。 | `getStatus()` は status object をそのまま返す。 | ステータス panel と手動実行 panelで running / last result を表示する。 |
| §27.9 | `GET /api/history?trigger=` を固定 enum で filter する。 | `getHistory({trigger})` は値を変換しない。 | 履歴 panel の trigger filter は固定 enum だけを選択肢にする。 |
| §27.10 | endpoint は追加しない。 | SDK method は追加しない。 | UI 操作は追加しない。runner 起動時検証だけ。 |
| §27.11 | schedule 系 API を使用する。 | schedule 系 method を使用する。 | リポジトリ情報 panel で interval/pause/resume/allowed hours を操作する。 |
| §27.12 | `POST /api/webhook`、webhook config API。 | config 表示 / 保存は `getWebhookConfig()` / `setWebhookConfig()`。 | 通知設定 panel で secret 設定状態を表示 / 保存する。 |
| §27.13 | `GET /api/webhook-events`。 | `getWebhookEvents(limit,offset)`。 | 通知設定または webhook event 一覧で delivery / result を表示する。 |
| §27.14 | stats API 群。 | stats method 群。 | 統計 panel で数値と時系列を表示する。 |
| §27.15 | snapshot / rollback API 群。 | snapshot / rollback method 群。 | snapshot panel と履歴 panel で download/delete/rollback を操作する。 |
| §27.16 | `GET /api/health` は認証不要。 | `health()`。 | UI 初期化の必須前提にはしない。診断 panel で表示してよい。 |
| §27.17 | `GET /api/logs/search?level=`。 | `searchLogs(q,from,to,level)`。 | ログビューアに INFO/WARNING/ERROR/DEBUG filter を置く。 |
| §27.18 | branch config API。 | `getBranchConfig()` / `setBranchConfig()`。 | リポジトリ情報 panel で branch target を表示 / 保存する。 |
| §27.19 | notify weekly summary API。 | `notifyWeeklySummary()`。 | 通知設定 panel で手動送信結果を表示する。 |
| §27.20 | `.config_log` を返す API。 | `getConfigLog()`。 | 設定 panel に diff_text を表示し、secret は `"***"` のまま表示する。 |

UI は、上表に存在しない §27.1〜§27.20 の SDK method を呼んではならない。SDK は、上表および §23 に存在しない API endpoint を呼んではならない。API は、上表で endpoint 追加なしとした機能に endpoint を追加してはならない。

### 27.1 GitHub Commit Status API

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.1 を正とする。owner component は `runner`、collaborator component は `commitstatus`、`statefile` とする。

### 27.2 ドライラン実行モード

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.2 を正とする。owner component は `runner`、collaborator component は `statefile` とする。

### 27.3 ビルド失敗時の自動リトライ

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.3 を正とする。owner component は `runner`、collaborator component は `statefile` とする。

### 27.4 出力サイトへのビルドメタ埋め込み

本節の主本文は `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §27.4 を正とする。owner component は `builder`、collaborator component は `runner`、`api`、`statefile` とする。

### 27.5 設定バリデーション API

本節の主本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` §27.5 を正とする。owner component は `api`、collaborator component は `sdk`、`ui`、`statefile` とする。

### 27.6 API アクセスログ

本節の主本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` §27.6 を正とする。owner component は `api`、collaborator component は `sdk`、`ui`、`statefile` とする。

### 27.7 ビルドログのアーカイブ圧縮

本節の主本文は `ADLAIRE_CI_DETAIL_ARCHIVE_SPEC.md` §27.7 を正とする。owner component は `archive`、collaborator component は `runner`、`api`、`statefile` とする。

### 27.8 ビルドステータスファイル出力

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.8 を正とする。owner component は `runner`、collaborator component は `api`、`statefile` とする。

### 27.9 ビルドトリガー種別の記録

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.9 を正とする。owner component は `runner`、collaborator component は `api`、`sdk`、`ui`、`statefile` とする。

### 27.10 設定ファイル起動時整合性チェック

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.10 を正とする。owner component は `runner`、collaborator component は `statefile` とする。

### 27.11 ポーリング間隔の動的変更

本節の主本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` §27.11 を正とする。owner component は `api`、collaborator component は `runner`、`statefile` とする。

### 27.12 GitHub Webhook 受信

本節の主本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` §27.12 を正とする。owner component は `api`、collaborator component は `runner`、`statefile` とする。

### 27.13 Webhook イベントログ / 一覧取得 API

本節の主本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` §27.13 を正とする。owner component は `api`、collaborator component は `sdk`、`ui`、`statefile` とする。

### 27.14 ビルド所要時間の記録と統計 API

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.14 を正とする。owner component は `runner`、collaborator component は `api`、`statefile`、`archive` とする。

### 27.15 ビルドアーティファクト管理

本節の主本文は `ADLAIRE_CI_DETAIL_ARCHIVE_SPEC.md` §27.15 を正とする。owner component は `archive`、collaborator component は `api`、`sdk`、`ui`、`runner`、`statefile` とする。snapshot 作成は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §14b を正とする。

### 27.16 ヘルスチェックエンドポイント

本節の主本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` §27.16 を正とする。owner component は `api`、collaborator component は `statefile` とする。

### 27.17 ビルドログ重大度フィルター

本節の主本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` §27.17 を正とする。owner component は `api`、collaborator component は `sdk`、`ui`、`archive`、`statefile` とする。

### 27.18 ブランチ設定の動的変更 API

本節の主本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` §27.18 を正とする。owner component は `api`、collaborator component は `runner`、`statefile` とする。

### 27.19 週次ビルドサマリー Webhook

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.19 を正とする。owner component は `runner`、collaborator component は `api`、`statefile` とする。

### 27.20 設定変更の詳細 diff 記録

本節の主本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` §27.20 を正とする。owner component は `api`、collaborator component は `statefile` とする。

### 27.21 複数ファイル監視

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.21 を正とする。owner component は `runner`、collaborator component は `builder`、`api`、`statefile` とする。

### 27.22 ビルドパイプライン YAML 定義

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.22 を正とする。owner component は `runner`、collaborator component は `api`、`statefile` とする。

### 27.23 ローカルファイル監視モード

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.23 を正とする。owner component は `runner`、collaborator component は `statefile` とする。

### 27.24 タグ付きコミットのみビルド

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.24 を正とする。owner component は `runner`、collaborator component は `api`、`statefile` とする。

### 27.25 ビルドキャッシュ

本節の主本文は `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §27.25 を正とする。owner component は `builder`、collaborator component は `runner`、`statefile` とする。

### 27.26 並列マルチターゲットビルド

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.26 を正とする。owner component は `runner`、collaborator component は `statefile` とする。

### 27.27 ビルド前後フック

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.27 を正とする。owner component は `runner`、collaborator component は `api`、`statefile` とする。

### 27.28 依存ファイルトラッキング

本節の主本文は `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §27.28 を正とする。owner component は `builder`、collaborator component は `runner`、`statefile` とする。

### 27.29 リモートビルド対応

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.29 を正とする。owner component は `runner`、collaborator component は `api`、`archive`、`statefile` とする。

### 27.30 ビルド承認フロー

本節の主本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` §27.30 を正とする。owner component は `api`、collaborator component は `runner`、`sdk`、`ui`、`statefile` とする。

### 27.31 ブランチ別環境変数

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.31 を正とする。owner component は `runner`、collaborator component は `api`、`statefile` とする。

### 27.32 ビルド通知連携

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.32 を正とする。owner component は `runner`、collaborator component は `api`、`sdk`、`ui`、`statefile` とする。

### 27.33 ビルド時間トレンド記録

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.33 を正とする。owner component は `runner`、collaborator component は `api`、`statefile` とする。

### 27.34 ビルド依存チェーン

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.34 を正とする。owner component は `runner`、collaborator component は `api`、`statefile` とする。

### 27.35 ビルド優先度キュー

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.35 を正とする。owner component は `runner`、collaborator component は `api`、`statefile` とする。

### 27.36 失敗原因の自動分類

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.36 を正とする。owner component は `runner`、collaborator component は `api`、`statefile` とする。

### 27.37 ビルド実行環境の記録

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.37 を正とする。owner component は `runner`、collaborator component は `statefile` とする。

### 27.38 ビルド所要時間の異常検知

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.38 を正とする。owner component は `runner`、collaborator component は `api`、`statefile` とする。

### 27.38a Runner 拡張機能 実装補足契約

本節は、責務 component 別詳細仕様へ分割しない。本節の内容は、runner、builder、API、SDK、UI、statefile、archive の境界をまたぐ横断補足契約であり、特定 owner component の主本文として扱ってはならない。

本節は §27.21〜§27.38 の runner / builder 拡張機能に共通する補足契約である。各個別節と矛盾する場合は個別節を優先する。

**設定 key と保存責務：**

| 機能 | 設定 key / 状態 | 更新責務 | 読取責務 |
|------|------------------|----------|----------|
| 複数ファイル監視 | `.branch_config.branch_targets[].target_files`、`.sha_cache/` | API は `.branch_config`、runner は SHA cache | runner、API |
| YAML pipeline | `.pipeline.yml`、`.pipeline_config` | API は `.pipeline_config`、runner はログのみ | runner |
| local watch | `.server_config.watch_mode`、`.local_watch_state.json` | API は `.server_config`、runner は local watch state | runner |
| tag filter | `.server_config.tag_filter` | API | runner |
| build cache | `.server_config.build_cache_enabled`、`.build_cache.json`、`.build_cache/pages/` | API は `.server_config`、builder は cache | builder、runner |
| deploy parallelism | `.server_config.deploy_parallelism` | API | runner |
| hooks | `.hooks`、`.build_logs/{build_id}_hook_{hook_id}.json` | API は `.hooks`、runner は hook log | runner、API |
| dependency manifest | `.dependency_manifest.json` | builder | runner、builder |
| remote build | `.server_config.remote_build` | API | runner |
| approvals | `.approval_queue`、`.build_state.queued` | runner と API | runner、API |
| branch env | `.branch_config.branch_targets[].env` | API | runner |
| notifications | `.notify_config`、`.notify_log`、`.notify_pending` | API は config、runner は log/pending | runner、API |
| build trends | `.build_trends.json` | runner | runner、API |
| build chain | `.build_chain_config` | API | runner |
| priority queue | `.build_state.queued[]` | API と runner | runner、API |
| failure category | `.build_logs/{id}.json.failure_category`、`.build_history.failure_category` | runner | API、UI |
| environment | `.build_logs/{id}.json.environment` | runner | API、UI |
| duration anomaly | `.server_config.duration_anomaly`、`.build_trends.json`、`.build_history` | API は config、runner は判定結果 | runner、API |

上表にない状態ファイルへ保存してはならない。複数機能が同じ状態ファイルを更新する場合、§22.0a のファイル lock を共有し、読み込み直後の最新内容に対して差分を適用する。

**runner 起動時の拡張機能処理順：**

1. CLI 引数、`--dry-run`、`--state-dir`、`--config` を検証する。
2. `--dry-run` の場合は `.build_lock` を取得せず、状態 directory 作成、backup、初期化、正規化、log / history / status 保存を行わない。設定、既存状態、fake GitHub read、差分判定だけを `§27.2` の dry-run flow で評価し、stdout JSON を出力して終了する。
3. `.build_lock` を取得する。
4. §27.10 の起動時整合性チェックを実行する。
5. `.server_config`、`.branch_config`、`.notify_config`、`.build_chain_config`、`.hooks` を読む。
6. queue entry がある場合は §27.35 の順序で 1 件だけ選ぶ。
7. §27.30 approval timeout を更新する。
8. watch mode、branch target、target files、tag filter、dependency manifest から build 対象を決定する。
9. build id、trigger、environment、start time を確定する。
10. pre hook、remote build または local builder、pipeline、dependency manifest、cache、deploy、post hook を個別節の順序で実行する。
11. failure category、duration、trend、duration anomaly、status、history、notification を保存する。
12. `.build_lock` を削除し、`.build_status.json.running=false` を finalizer として保存する。

手順 9 以降で異常終了した場合でも、可能な限り `.build_status.json.running=false` を保存する。finalizer 保存に失敗した場合は ERROR ログを出し、終了コードを最低 `1` にする。

**機能別の不変条件：**

| 機能 | 不変条件 |
|------|----------|
| 複数ファイル監視 | `target_file` と `target_files` が同時に存在する場合、API response は両方を返してよいが、runner 内部では `target_files` に正規化して処理する。重複 target は 1 回だけ build 対象にする。 |
| YAML pipeline | `.pipeline.yml` と `.pipeline_config.inline_yaml` が両方存在する場合、`.pipeline.yml` を優先する。どちらを使用したかを `.build_logs/{id}.json.pipeline_source` に保存する。 |
| local watch | local mode では GitHub Commit Status、GitHub rate limit、GitHub tag refs を呼ばない。tag filter が有効な場合は設定不整合として終了コード `2`。 |
| tag filter | skip 時は SHA cache を更新しないため、次回も同じ commit を評価する。matched tag 名は最大 100 件まで build log に保存する。 |
| build cache | cache hit の出力は通常変換結果と byte 単位で一致しなければならない。不一致検出時は hit を破棄し miss として再変換する。 |
| deploy parallelism | target result は設定順で保存する。実行完了順で保存してはならない。 |
| hooks | hook stdout/stderr に secret mask を適用してから保存する。hook log の保存失敗は runner log に ERROR を出し、pre hook の場合は build を中断する。 |
| dependency manifest | manifest 生成は build 成功後だけ確定保存する。失敗 build の途中 manifest で既存 manifest を上書きしない。 |
| remote build | remote artifact 展開先は state dir 配下の一時 directory とし、既存 output directory へ直接展開しない。検証成功後に deploy 処理へ渡す。 |
| approvals | approval entry は pending の間だけ approve/reject 可能。approved、rejected、expired を物理削除せず、API 一覧で状態を返す。 |
| branch env | secret key の値は child process には渡すが、build log には key 名だけ保存する。mask は stdout/stderr、hook log、notification payload に適用する。 |
| notifications | 通知送信は build 成否を反転させない。通知失敗は `.notify_log` と `.notify_pending` だけで表現する。 |
| build trends | 同一 build id の sample が既にある場合は append せず置換する。再実行や rollback で別 build id なら別 sample とする。 |
| build chain | chain job ごとに独立した build log を作成する。同一 chain 内の job は `chain_run_id` を共有する。 |
| priority queue | queue 取り出し時に対象 entry を `.build_state.queued` から削除し、`.build_state.running=true` と同一 lock 内で保存する。 |
| failure category | failure 以外の status では `failure_category` を `null` とする。過去互換で値がある成功行は API response に warning を付ける。 |
| environment | `state_dir` は保存してよいが、home directory 内の secret file path は保存しない。該当する場合は basename のみ保存する。 |
| duration anomaly | anomaly 判定は trend 更新前の summary で行う。同じ build id の再判定で tag を重複追加しない。 |

**API 更新時の補足検証：**

| API | 追加検証 |
|-----|----------|
| `POST /api/branch-config` | `target_file` と `target_files` の少なくとも一方が必須。両方ある場合は正規化後に同一 target 集合になること。 |
| `POST /api/pipeline-config` | `inline_yaml` は最大 64 KiB。保存前に YAML subset parse を実行し、parse 不能なら `422`。 |
| `POST /api/config` | `watch_mode`, `tag_filter`, `build_cache_enabled`, `deploy_parallelism`, `remote_build`, `duration_anomaly`, `build_trend_keep_count`, `approval_timeout_seconds` を検証する。 |
| `POST /api/hooks` | `command_args[0]` が空、相対 path かつ PATH 解決不能、または 256 文字超の場合は `422`。 |
| `POST /api/build-chain-config` | job id 重複、循環、未定義依存、disabled job への required 依存を `422`。 |
| `POST /api/notify-config` | channel id 重複、未知 event、secret 平文の GET response 混入を禁止する。 |

**受け入れ fixture 固定：**

| fixture | 入力 | 合格条件 |
|---------|------|----------|
| runner extension dry-run | §27.21〜§27.38 の各設定を 1 つずつ有効化した dry-run | 状態ファイル差分なし、`would_call` と `would_write` が個別節どおり。 |
| runner extension success | local fixture で build 成功 | build log、history、status、trend、notification の保存順と値が一致する。 |
| runner extension failure | pipeline timeout、deploy failure、hook abort、config error | failure_category、history status、status file、通知 event が固定値になる。 |
| API config validation | 各設定 key の正常値、境界値、範囲外、不正型 | 正常値は保存、範囲外と不正型は `422`、状態差分なし。 |
| secret masking | env secret、webhook secret、SMTP password、API token を含む build / API 操作 | response、log、history、audit、notify log、UI 表示に平文がない。 |

**§27.21〜§27.38 Runner 拡張機能別実装完全性固定契約：**

§27.21〜§27.38 は、runner、builder、API、SDK、UI にまたがる拡張機能である。各機能は個別節と §27.38a に加えて下表を満たした場合だけ実装完了とする。下表は既存詳細機能の実装完了判定であり、未定義 endpoint、未定義状態ファイル、未定義 UI、外部公開構成、将来計画機能を追加する根拠にしてはならない。

| 節 | 機能 | 実装入口 | 正規化 / 判定 | 成功時副作用 | 失敗時副作用 | 必須 fixture |
|----|------|----------|---------------|--------------|--------------|--------------|
| §27.21 | 複数ファイル監視 | runner target selection、branch config API。 | `target_file` と `target_files` を重複除去済み配列へ正規化する。 | 変更あり target だけ build し、対象別 SHA cache を更新する。 | 不正 path は build 前に停止し、SHA cache を更新しない。 | 1 file change、複数 change、重複 target、path traversal、全 target skip。 |
| §27.22 | YAML pipeline | runner pipeline resolution、pipeline config API。 | `.pipeline.yml` を優先し、未設定時だけ inline YAML を使う。 | 使用元、step 結果、終了コードを build log に保存する。 | parse 失敗は build 開始前 failure、deploy / SHA 更新なし。 | file 優先、inline、parse error、step failure、secret mask。 |
| §27.23 | ローカルファイル監視 | runner watch mode。 | GitHub 入力を使わず local mtime / checksum で差分判定する。 | local watch state を更新し、trigger=`local_watch` を保存する。 | tag filter 併用など不整合は終了コード `2`、state 更新なし。 | change、no change、deleted file、state corrupt、GitHub call 0。 |
| §27.24 | タグ付きコミットのみビルド | runner tag refs 判定。 | 許可 pattern に一致する tag だけ build 対象にする。 | matched tag を最大 100 件 build log に保存する。 | tag API 失敗または不正 pattern では SHA cache を更新しない。 | matched、unmatched、`*`、prefix、API failure、不正 pattern。 |
| §27.25 | ビルドキャッシュ | builder page conversion。 | 入力 hash、theme、renderer version、関連依存を cache key に含める。 | cache hit は通常出力と byte 一致し、miss は cache を atomic 保存する。 | 不一致 hit は破棄して miss 扱い。保存失敗は build failure にしない。 | hit、miss、stale、corrupt、byte mismatch、save failure。 |
| §27.26 | 並列マルチターゲットビルド | runner target executor。 | 設定順を canonical order とし、並列完了順に依存しない。 | target result を設定順で log/history/status に保存する。 | 一部失敗は個別 target failure と全体 status を固定規則で保存する。 | all success、partial failure、all failure、timeout、order stable。 |
| §27.27 | ビルド前後フック | runner hook executor、hooks API。 | `command_args` 配列だけを実行し、shell 文字列展開を禁止する。 | hook log を secret mask 後に保存し、post hook は build 結果を反転しない。 | pre hook 失敗は build 中断。hook log 保存失敗時の終了条件を固定する。 | pre success、pre fail、post fail、timeout、mask、log failure。 |
| §27.28 | 依存ファイルトラッキング | builder dependency collector。 | 入力 Markdown、参照画像、link、include 相当の依存を相対 path で保存する。 | build 成功時だけ dependency manifest を置換保存する。 | build 失敗時は既存 manifest を上書きしない。 | no deps、multi deps、missing dep、path normalize、failure keeps old。 |
| §27.29 | リモートビルド対応 | runner remote build executor。 | remote artifact を一時 directory へ取得し、checksum / manifest 検証後だけ採用する。 | 検証済み artifact を deploy / snapshot へ渡す。 | remote failure、checksum 不一致、展開失敗は既存出力を変更しない。 | success、auth fail、checksum mismatch、unsafe archive、timeout。 |
| §27.30 | ビルド承認フロー | queue、approval API、runner dequeue。 | approval entry は pending / approved / rejected / expired の固定 enum とする。 | approved entry だけ build queue へ進め、監査ログを残す。 | rejected / expired は build せず物理削除しない。 | approve、reject、expire、double approve、running conflict。 |
| §27.31 | ブランチ別環境変数 | runner process env、branch config API。 | branch target に一致した env だけ child process へ渡す。 | build log には key 名と mask 済み値だけ保存する。 | secret 漏えい検出時は保存前に mask し、未定義 key は渡さない。 | branch match、default、secret mask、unknown key、child env only。 |
| §27.32 | ビルド通知連携 | notify config、runner notifier。 | event、channel、retry policy を正規化して送信対象を決める。 | notify log / pending を更新する。build 成否は反転しない。 | 送信失敗は pending 化し、secret を payload / log に残さない。 | success、failure、retry、disabled、missing channel、mask。 |
| §27.33 | ビルド時間トレンド記録 | runner finalizer、stats API。 | build id 単位で sample を upsert する。 | trend summary と history duration を更新する。 | trend 保存失敗は build 結果を維持し、終了コードを最低 `1` にする。 | append、replace same id、retention、corrupt、save failure。 |
| §27.34 | ビルド依存チェーン | chain config API、runner scheduler。 | job id、依存、循環、disabled dependency を build 前に検証する。 | job ごとの build log と共通 `chain_run_id` を保存する。 | 必須依存失敗時は後続 required job を skipped として固定保存する。 | linear、parallel、cycle、required fail、optional fail、disabled dependency。 |
| §27.35 | ビルド優先度キュー | queue API、runner dequeue。 | priority、created_at、id で安定順序を決める。 | 取り出しと running 設定を同一 lock 内で保存する。 | queue full は `429`、取り出し失敗は queue entry を残す。 | priority order、FIFO tie、queue full、dequeue atomic、cancel。 |
| §27.36 | 失敗原因の自動分類 | runner failure mapper。 | exit code、stderr pattern、API error、timeout を固定 enum へ分類する。 | failure build だけ log/history に category を保存する。 | 分類不能は `unknown`。成功 build は `null`。 | timeout、validation、network、deploy、hook、unknown、success null。 |
| §27.37 | ビルド実行環境の記録 | runner environment collector。 | OS、arch、binary version、working dir、state dir を secret 除外して保存する。 | build log の environment object に保存する。 | 取得不能 field は `null`。secret path は basename だけ保存する。 | full、partial null、secret path、version missing。 |
| §27.38 | 所要時間異常検知 | trend summary、duration anomaly config。 | 判定は trend 更新前 summary で行う。 | anomaly tag、history flag、通知 event を重複なく保存する。 | sample 不足、failure build、通知失敗では build 成否を反転しない。 | sample不足、平均超過、p95正常、failure除外、通知失敗、tag重複なし。 |

**§27.21〜§27.38 API / SDK / UI 連動固定契約：**

| 機能群 | API | SDK | UI |
|--------|-----|-----|----|
| 監視 / target / trigger | branch config、queue、approval、history、stats API だけを使う。 | §23 にある method だけを呼び、target 正規化を重複実装しない。 | 固定 enum と API response の選択肢だけを表示する。 |
| pipeline / hook / remote build | config 保存 API と log 参照 API だけを使う。 | request body を仕様 key のまま送る。shell 文字列化しない。 | secret field は保存後に空にし、command_args は配列 UI として扱う。 |
| cache / dependency / artifact | build log、output meta、snapshot、artifact API だけを使う。 | binary download は response body を変換しない。 | download / delete / rollback 成功後は status、history、snapshot を再取得する。 |
| notification / trend / anomaly | notify、stats、history API だけを使う。 | retry / pending を成功扱いに変換しない。 | 送信失敗を build 失敗として表示しない。 |
| security interaction | API token、rate limit、audit、access log の共通契約を通す。 | `401` では token を破棄し、`403` と区別する。 | `401` で login panel に戻し、secret field を消去する。 |

**§27.42〜§27.47 認証・監査・制限機能 実装完全性固定契約：**

§27.42〜§27.47 は、API token、監査、session、TOTP、rate limit に関する安全機能である。各機能は個別節に加えて下表を満たした場合だけ実装完了とする。

| 節 | 機能 | 判定入口 | 成功時副作用 | 失敗時副作用 | 漏えい禁止値 | 必須 fixture |
|----|------|----------|--------------|--------------|--------------|--------------|
| §27.42 | API token scope | route / method 確定後、body parse 前。 | 許可 endpoint だけ処理し、必要時 audit に actor を残す。 | 権限不足は対象処理を実行せず `403`。audit 失敗時は `500`。 | token 本体、Authorization header。 | trigger allowed、read denied、multi scope、path param、body 未評価、audit failure。 |
| §27.43 | API key 管理 | admin session または admin scope。 | token hash だけ保存し、作成時だけ token 本体を返す。 | validation 失敗は保存差分なし。失効済み token は再有効化しない。 | token 本体、token hash の不要露出。 | create、list mask、revoke、expired、duplicate label、admin token create。 |
| §27.44 | 監査ログ | security / config / operation event 確定時。 | 1 event 1 JSON Lines で追記し、actor / target / result を保存する。 | 必須 audit 失敗は対象処理を `500` にする。任意 audit は個別節優先。 | secret、password、token、TOTP secret、raw request body。 | success、denied、failure、mask、append failure、pagination。 |
| §27.45 | session timeout | login、authenticated request、timeout config API。 | session の last_seen / expires_at を固定規則で更新する。 | timeout session は `401`、対象 endpoint は実行しない。 | session token。 | active、expired、sliding update、config update、revoke all、clock boundary。 |
| §27.46 | TOTP | setup、confirm、login/totp、disable。 | secret は confirm 成功後だけ有効保存し、ticket は一回だけ使う。 | ticket 再利用、期限切れ、code 不正は対象状態を変更しない。 | TOTP secret、backup code 相当値、ticket token。 | setup、confirm、login success、code reuse、disable、audit failure。 |
| §27.47 | API rate limit | route / auth / scope 判定の定義済み位置。 | 上限未満だけ count を増やし endpoint 処理へ進む。 | `429` は count を増やさず、audit 成功時だけ返す。 | API token、session token、request body。 | login 11 回目、window reset、IP+actor、disabled、policy update、state save failure。 |

**§27.42〜§27.47 セキュリティ機能 横断順序固定契約：**

| 順序 | 処理 | 固定条件 |
|------|------|----------|
| 1 | route / method を確定する。 | 未定義 route は認証、rate limit、body parse より前に `404` / `405`。 |
| 2 | 認証不要 endpoint を判定する。 | `GET /api/health` と `POST /api/webhook` は個別契約を優先する。 |
| 3 | login rate limit を判定する。 | login group は認証前 IP key で判定する。 |
| 4 | 認証情報を検証する。 | session と API token を混同しない。形式不一致は `401`。 |
| 5 | API token scope を判定する。 | scope 不足時は body validation と状態更新を行わない。 |
| 6 | 認証後 rate limit を判定する。 | actor key と IP key を同一 lock 内で判定・更新する。 |
| 7 | endpoint 固有 validation を行う。 | 失敗時は対象状態、外部 API、外部 command を変更しない。 |
| 8 | endpoint 固有処理を実行する。 | 成功時だけ個別節の保存順で状態、access log、audit log を確定する。 |

上表の順序を変更してはならない。個別節が別順序を明記する場合は、セキュリティ上の漏えいを増やさない範囲で個別節を優先する。順序変更が必要な場合は、先に本表と該当個別節を同時に改訂する。


### 27.42 ビルドトリガー専用 API スコープ

本節の主本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` §27.42 を正とする。owner component は `api`、collaborator component は `sdk`、`ui`、`security`、`statefile` とする。

### 27.43 API キー管理

本節の主本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` §27.43 を正とする。owner component は `api`、collaborator component は `sdk`、`ui`、`security`、`statefile` とする。

### 27.44 監査ログ

本節の主本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` §27.44 を正とする。owner component は `api`、collaborator component は `security`、`statefile` とする。

### 27.45 セッションタイムアウト変更設定

本節の主本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` §27.45 を正とする。owner component は `api`、collaborator component は `sdk`、`ui`、`security`、`statefile` とする。

### 27.46 TOTP 二要素認証

本節の主本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` §27.46 を正とする。owner component は `api`、collaborator component は `sdk`、`ui`、`security`、`statefile` とする。

### 27.47 API レート制限

本節の主本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` §27.47 を正とする。owner component は `api`、collaborator component は `sdk`、`ui`、`security`、`statefile` とする。
