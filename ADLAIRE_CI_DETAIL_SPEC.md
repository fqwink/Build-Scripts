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
5. §26 のセットアップ・アップデート手順と §26.7 の受け入れ条件に影響がある場合は、実装 PR の検証対象に含める。

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
| §26 | バイナリ配布前提のセットアップ、アップデート、受け入れ条件 |

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
| `components/api.go` | `ADLAIRE_CI_DETAIL_API_SPEC.md` §21〜§22、§25、§27.5〜§27.6、§27.11〜§27.13、§27.16〜§27.18、§27.20、§27.30、§27.42〜§27.47 | API 共通処理、endpoint、状態ファイル read/write、認証、session、API owner 追加機能。 |
| `admin/adlaire-ci-sdk.js` | `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 | SDK class、method、HTTP 対応、error、stream、token 破棄。 |
| `admin/index.html` | `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 | 画面構成、DOM id、panel、SDK 呼び出し、表示状態、秘密情報消去。 |
| `components/mcp.go` | 詳細仕様なし | 本ファイルでは実装可能な入出力、状態、起動手順、検証条件を定義しない。 |

---

## 0b.1 責務 component 別 詳細仕様ファイル分割仕様

本節は、`ADLAIRE_CI_DETAIL_SPEC.md` を責務 component 別に分割する場合の固定仕様である。分割は、仕様内容の移動と参照先更新だけを対象とし、機能追加、実装状態変更、実装可否変更、ロードマップ変更、方針・ポリシー追加を含めてはならない。

分割後の詳細仕様ファイルは以下に固定する。`COMMON`、`CORE`、`BASE`、`SHARED`、`FOUNDATION`、その他の横断共通基盤ファイルは作成しない。

| ファイル | 持つ内容 | 持たない内容 |
|----------|----------|--------------|
| `ADLAIRE_CI_DETAIL_SPEC.md` | 詳細仕様の入口、読み方、共通固定値、実装前確認項目、検証マトリクス、Phase、詳細節対応表、リポジトリ内ソース配置、責務 component 分割仕様。 | 各 component の詳細な処理本文、fixture 詳細、個別 endpoint 詳細、個別 UI 操作詳細。 |
| `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` | `builder` owner の Markdown 変換、静的 Web サイト出力、HTML / CSS / JavaScript、theme component、builder fixture。 | runner / API / SDK / UI の実行責務。 |
| `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` | `runner` owner の GitHub 監視、状態ファイル更新、pipeline、deploy、snapshot、通知、runner fixture。 | API endpoint の認証・応答本文、SDK method、UI DOM 詳細。 |
| `ADLAIRE_CI_DETAIL_API_SPEC.md` | `api` owner の HTTP 共通契約、endpoint、状態ファイル read/write、認証連携、API fixture。 | SDK 内部実装、UI DOM 詳細、runner の build 実行責務。 |
| `ADLAIRE_CI_DETAIL_SDK_SPEC.md` | `sdk` owner の SDK class、method、HTTP 対応、error、stream、token 破棄。 | API endpoint の状態ファイル更新責務、UI DOM 詳細。 |
| `ADLAIRE_CI_DETAIL_UI_SPEC.md` | `ui` owner の DOM id、panel、操作、表示状態、SDK 呼び出し、秘密情報消去。 | SDK method 実装、API endpoint 実装、状態ファイル直接操作。 |
| `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` | `setup` owner のバイナリ配布、配置、systemd、セットアップ、アップデート、リリース成果物検証。 | runner / API / SDK / UI の個別機能本文。 |
| `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` | fixture manifest、assertion、fake、testdata、受け入れ fixture 共通契約、PR 証跡テンプレート。 | 個別 component の通常処理本文。 |

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
5. fixture、fake、PR 証跡が必要な場合は `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` を読む。

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
| atomic write | JSON / text 状態ファイル更新は同一ディレクトリに一時ファイルを書き出し、`file.Sync()` と `file.Close()` の成功後に `os.Rename` で置換する。同一ファイルシステム外への一時ファイル作成は禁止する。 |
| 権限 | 秘密情報を含むファイルは `0600`、通常状態ファイルは `0644`、ディレクトリは `0755` を既定値とする。既存ファイル更新時も権限が緩い場合は既定値へ補正する。 |
| ロック | 共有状態ファイル更新は `{filename}.lock` を同一ディレクトリに作成して排他する。ロック取得待ちは runner では 0 秒、API では最大 10 秒。超過時は runner が ERROR ログで当該処理をスキップし、API は `409 Conflict` を返す。 |
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
| セットアップ | systemd | §26 の unit 名、`ExecStart`、配置パス、権限、起動確認コマンドが実際の導入手順と一致する。 |

検証結果は、実装 PR の本文または実装完了報告に、対象、実行コマンド、期待結果、実結果を対応付けて記録する。検証不能な項目がある場合は、その項目を完了扱いにしてはならない。

---

## 0f. 仕様策定完了チェック

本節は、Go 版初期実装へ進む前に仕様策定が完了しているかを判定するチェックである。実装者は、責務 component ごとに下表の必須条件を満たすまで実装を開始してはならない。

| 対象 | 実装着手条件 | 実装禁止条件 | 完了判定 |
|------|--------------|--------------|----------|
| `components/builder.go` | §2〜§8 に CLI option、入力 Markdown、出力サイトディレクトリ、終了コード、stderr、HTML 構造、テーマコンポーネント、JS/CSS、生成物確認が定義されている。 | §4〜§7 にない Markdown 記法、CSS class、JavaScript 機能、外部 asset、theme を追加すること。 | §0e の `components/builder.go` 必須検証をすべて満たし、生成サイトが §5〜§7 と一致する。 |
| `components/runner.go` | §10〜§20 に設定値、状態ファイル、GitHub API、SHA 比較、pipeline 起動、SSH 転送、snapshot、通知、ログ、systemd が定義されている。 | 未定義の環境変数、状態ファイル、queue 挙動、通知チャンネル、pipeline 形式を追加すること。 | §0e の `components/runner.go` 必須検証をすべて満たし、状態ファイル更新順序が §13、§22.0a、§22.0d と一致する。 |
| `components/api.go` | `ADLAIRE_CI_DETAIL_API_SPEC.md` §21〜§22、§25 と本ファイル §26 に API 共通契約、endpoint、状態ファイル schema、認証、認可、systemd、セットアップが定義されている。 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e にない endpoint、method、status code、response body、状態ファイル write を追加すること。 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e の全 endpoint が Request、Response、Errors、Read、Write、SDK、UI の対応表と一致する。 |
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
| Phase 3 | `components/api.go` P0 / P1 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §21〜§22、§25 と本ファイル §26 のうち、認証、セッション、共通エラー、状態ファイル読み書き、ビルド操作、status、logs、history、queue、circuit breaker。 | Phase 2 が完了し、runner が書き込む状態ファイル schema が固定されている。 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0f P0 / P1 の必須検証、§0e の `components/api.go` API 共通・状態ファイル検証、§0f の `components/api.go` 完了判定の該当範囲を満たす。 |
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
| 完了条件 | §0e の `components/builder.go` 必須検証、§0f の完了判定、§8a の全 fixture、§26.7 の build script 対象を満たす。 |

Phase 1 完了時は、次 Phase へ引き継ぐ CLI 契約として、`adlaire-ci-build` の終了コード、stdout 進捗行、`[WARN]` 行、`[REPORT]` 行、出力ディレクトリ構造を固定する。Phase 2 以降は、この契約を変更してはならない。変更が必要な場合は Phase 1 仕様改訂に戻る。

#### 0g.1.1 Phase 1 実装順序

Phase 1 は、以下の順序で実装する。順序を入れ替える場合は、入れ替え理由と影響が §8 の処理順序、§8a fixture、§26.7 build script 受け入れ条件に影響しないことを実装 PR 本文に記録する。

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
| 完了条件 | §0e の `components/runner.go` 必須検証、§0f の完了判定、§15a の全 fixture、§26.7 の runner 対象を満たす。 |

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
| 必須検証 | §22.0f P0 / P1、§25、§26.4.2、§26.7 API 対象。 |
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
| 6 | §22.0f P0 / P1、§25、§26.4.2、§26.7 API の確認を PR 本文に記録する。 | P0 / P1 を完了扱いにでき、Phase 4 が共通契約を再実装せず利用できる。 |

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
| 必須検証 | §22.0f P2〜P5 の検証条件、§22.0d の read/write 対応表、§26.7 API/security 対象。 |
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
| 必須検証 | `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 の method 契約、§0e の SDK 契約、§26.7 SDK 対象。API は fake fetch で成功 / error / timeout / stream を再現する。 |
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
| 必須検証 | `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 の全 panel と主要操作、§0e の UI 契約、§26.7 UI/security 対象。SDK は fake implementation で成功 / 失敗 / loading / stream を再現する。 |
| 完了条件 | 全 UI 操作が `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 の表示条件と `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 の SDK method を満たし、成功表示、失敗表示、disabled、再取得、secret 消去が一致する。 |

Phase 6 完了時は、初期実装全体の完了判定として、Phase 1〜6 の引き継ぎ契約、§0e、§0f、§0g、§0i、§0j、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0f、`ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23、`ADLAIRE_CI_DETAIL_UI_SPEC.md` §24、§26.7 を再確認する。

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

**Phase fixture / testdata 配置固定契約：**

| Phase | 必須配置 | 必須内容 | 禁止事項 |
|-------|----------|----------|----------|
| Phase 1 | `testdata/builder/single/`、`testdata/builder/site/`、`testdata/builder/empty-dir/`、`testdata/builder/strict/`、`testdata/builder/safe/`、各 fixture の `expected/`。 | 入力 Markdown、テーマ設定、asset 入力、期待 HTML / CSS / JS / search index、期待 stdout / stderr、期待終了コード。 | 実行環境ごとに変わる絶対 path、timestamp、乱数、外部 URL 取得結果を期待値へ含めてはならない。 |
| Phase 2 | `testdata/runner/r1/`〜`testdata/runner/r20/`、各 fixture の `state/`、`github/`、`pipeline/`、`ssh/`、`notify/`、`expected/`。 | GitHub fake response、状態ファイル初期値、lock 状態、pipeline fake 結果、deploy fake 結果、通知 fake 結果、期待 `.last_sha`、期待 queue / snapshot。 | 実 GitHub API、実 SSH、実通知先、実 remote branch 状態に依存して合否を決めてはならない。 |
| Phase 3 | `testdata/api/p0-p1/auth/`、`status/`、`history/`、`logs/`、`queue/`、`stream/`、`errors/`、各 fixture の `state/`、`requests/`、`responses/`、`expected/`。 | HTTP method / path / query / header / body、状態ファイル初期値、期待 response、期待 error body、SSE frame、状態 read/write 後の期待値。 | 仕様未定義 endpoint、P2 以降の endpoint、外部公開設定を fixture に含めてはならない。 |
| Phase 4 | `testdata/api/p2-p5/config/`、`notify/`、`snapshots/`、`maintenance/`、`hooks/`、`tokens/`、各 fixture の `state/`、`requests/`、`responses/`、`expected/`。 | config / notify / snapshot / rollback / maintenance / hook / token の正常系、validation error、secret mask、P0 / P1 互換確認の期待値。 | token 原文、secret 原文、mask 前 payload、再取得不可 token の復元値を fixture または expected に含めてはならない。 |
| Phase 5 | `testdata/sdk/request-shape/`、`error-shape/`、`stream/`、`binary/`、`p2-p5/`。 | fake fetch transcript、期待 request、期待 SDK return、期待 `AdlaireCIError`、期待 stream event、timeout / abort の期待結果。 | Node.js 専用 API、bundler、npm package、実 network、browser storage 依存を検証前提にしてはならない。 |
| Phase 6 | `testdata/ui/login/`、`status/`、`build/`、`config/`、`secret/`、`stream/`、`p2-p5/`。 | fake SDK script、入力 DOM 状態、操作手順、期待 DOM assertion、期待 SDK call、期待 disabled / loading / error / success 表示。 | 直接 `fetch()`、CDN、外部 framework、画像 snapshot だけの合否判定、secret 表示を含めてはならない。 |

**fake 実装固定契約：**

| fake | 対象 Phase | 必須動作 | 必須記録 |
|------|------------|----------|----------|
| fake GitHub server | Phase 2 | fixture の JSON 応答だけを返す。未定義 method / path は `404` とする。rate limit、`304`、`409`、`500` は fixture で明示された場合のみ返す。 | method、path、query、request body、認証 header の有無、呼び出し順。token 値は `***` に置換する。 |
| fake ssh executable | Phase 2 | fixture 指定の stdout、stderr、終了コード、timeout を返す。実 shell、実 SSH、実 file 転送は実行しない。 | argv、stdin 有無、環境変数名、終了コード、timeout 発生有無。secret 値は記録しない。 |
| fake notifier | Phase 2 / Phase 4 | fixture 指定の HTTP status、response body、timeout を返す。通知先へ送信しない。 | URL の host 部分、payload schema、mask 後 payload、retry 回数、最終結果。 |
| fake filesystem | Phase 1〜Phase 4 | atomic write 失敗、sync 失敗、lock 競合、JSON 破損、permission error を fixture 単位で再現する。通常 file I/O の代替にはしない。 | 対象 path、操作種別、注入した失敗、復旧後の状態。 |
| fake fetch | Phase 5 | `status`、`headers`、`body`、network error、timeout、abort、stream chunk を fixture どおり返す。実 network は使用しない。 | method、URL、query、headers、body 有無、abort 発生有無、呼び出し順。token 値は `***` に置換する。 |
| fake SDK | Phase 6 | SDK method ごとに固定 return、固定 throw、固定 stream event を返す。UI からの直接 API 呼び出しは受け付けない。 | method 名、引数、呼び出し順、throw した error code、stream unsubscribe 実行有無。 |

**実装 PR 証跡固定契約：**

| 証跡 | 必須記載 | 不足時の扱い |
|------|----------|--------------|
| 変更対象 | 対象 Phase、owner component、collaborator component、変更ファイル、追加 fixture / testdata path。 | 対象 Phase の成果物不足として未完了。 |
| 固定契約 | 追加または固定した CLI、状態 schema、HTTP API、SDK method、DOM id、fake 動作、終了コード、error body。 | 後続 Phase が参照できないため未完了。 |
| 検証 | 実行コマンド、fixture 名、期待結果、実結果、判定。 | 合否を再現できないため未完了。 |
| 未実装対象 | 対象 Phase 外の機能、将来計画、MCP、外部公開設定など実装していない範囲。 | 先取り実装または範囲不明として未完了。 |
| 後続 Phase への影響 | 後続 Phase が利用してよい contract と、利用してはならない未固定 contract。 | 次 Phase 着手条件未充足として未完了。 |
| secret 確認 | log、fixture、snapshot、UI 表示、PR 本文に secret / token / password 原文がないこと。 | security 不合格として未完了。 |

#### 0g.8.1 実装 PR 完了扱い禁止条件

以下のいずれかに該当する場合、実装 PR は merge 可能であっても Phase 完了として扱ってはならない。

| 条件 | 扱い |
|------|------|
| 対象 Phase の必須成果物が不足している。 | PR 本文に不足項目を明記し、Phase 未完了とする。 |
| fixture または testdata がなく、手動確認だけで合格としている。 | 仕様上 fixture が必須の Phase では未完了とする。 |
| 実装対象外の機能を先取りしている。 | 仕様違反として扱い、対象外機能を削除するか、先に仕様改訂する。 |
| 後続 Phase が依存する契約を PR 本文に固定していない。 | 後続 Phase の実装を開始してはならない。 |
| secret、token、password を log、fixture、snapshot、UI 表示へ平文出力している。 | security 不合格として Phase 未完了とする。 |
| §0g の順序、§0i の対応表、§0j のソース配置、§26.7 の受け入れ条件のいずれかと矛盾している。 | 仕様不整合として扱い、実装または仕様を修正する。 |

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

表の「詳細仕様節」が複数ある場合は、すべての節を同時に満たす。`components/builder.go` が責務 component に含まれる機能では、`ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` の同番号節を合わせて確認する。`components/runner.go` が責務 component に含まれる機能では、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` の同番号節を合わせて確認する。`components/api.go` が責務 component に含まれる機能では、`ADLAIRE_CI_DETAIL_API_SPEC.md` の同番号節を合わせて確認する。`admin/adlaire-ci-sdk.js` が責務 component に含まれる機能では、`ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 を合わせて確認する。`admin/index.html` が責務 component に含まれる機能では、`ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 を合わせて確認する。該当節に §0h の必須項目が不足している場合は、その項目を実装せず、先に詳細仕様を改訂する。

| 機能 | 責務 component | 詳細仕様節 | 受け入れ条件 |
|------|-------------------|------------|--------------|
| ビルドタイムアウト | `components/runner.go` / `components/api.go` | §12、§13、§22.0e | `build_timeout_seconds` の既定値、設定 API、`context.WithTimeout` の中断処理、終了コード、ログが一致する。 |
| ポーリング間隔の動的変更 | `components/api.go` | §22.0e、§26、§27.11 | `POST /api/schedule/interval` が systemd timer 設定を更新し、検証コマンドで反映を確認できる。 |
| ビルドログのファイル保存 | `components/runner.go` | §11、§13、§15 | `.build_logs/{id}.json` の schema、stdout/stderr、変換レポート、duration、権限が一致する。 |
| GitHub Webhook 受信 | `components/api.go` / `components/runner.go` | §22.0e、§22-W、§13、§27.12 | HMAC 検証、イベント記録、キュー投入またはビルドトリガー、エラー応答が一致する。 |
| ネットワーク断時の再試行 | `components/runner.go` | §12、§13 | `API_RETRY_MAX`、`API_RETRY_BASE_SECONDS`、指数バックオフ、失敗時ログが一致する。 |
| GitHub API レート制限自動待機 | `components/runner.go` / `components/api.go` | §13、§22.0e | `X-RateLimit-Remaining` と `X-RateLimit-Reset` の扱い、待機、API 表示が一致する。 |
| 転送後リモート整合性検証 | `components/runner.go` | §14a、§13 | SSH 転送後の SHA256 照合、不一致時の `.pending_transfers` 再投入、ログが一致する。 |
| マルチブランチビルド | `components/runner.go` / `components/api.go` | §12、§13、§22.0e | `BRANCH_TARGETS` と `.branch_config` の優先順位、順次処理、API 更新が一致する。 |
| ビルドログ世代管理 | `components/runner.go` | §12、§13、§15 | `LOG_KEEP_N` 超過時の削除順序、0 の扱い、削除ログが一致する。 |
| ビルド出力の外部転送 | `components/runner.go` | §14a、§13 | SSH 差分転送、複数ファイル処理、失敗時 pending、通知が一致する。 |
| ビルドクールダウン | `components/runner.go` | §12、§13 | `BUILD_COOLDOWN_SECONDS` 内の起動スキップ、Webhook 二重トリガー抑止、ログが一致する。 |
| ビルド前の事前チェック | `components/runner.go` | §13、§26 | ディスク、`adlaire-ci-build`、pipeline 前提の確認、不足時の ERROR と通知が一致する。 |
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
| `components/statefile.go` | `.build_history`、`.build_logs`、`.server_config` など状態ファイルの読み書きを扱う。 |
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

### — 管理ツール —

## 21. API 詳細仕様（分割済み）

`api` owner component の詳細仕様本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` を正とする。

| 移動済み範囲 | 移動先 |
|--------------|--------|
| §21〜§22 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §21〜§22 |
| §25 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §25 |
| api owner の §27 個別節 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §27 |

`ADLAIRE_CI_DETAIL_SPEC.md` §27.38a は横断補足契約として本ファイルに残す。

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

## 26. セットアップ・アップデート手順

本節は、Go 版 Adlaire CI のセットアップ手順を定義する。

### §26.1 要件

| 項目 | 要件 |
|------|------|
| Go 版バイナリ | `adlaire-ci-build`、`adlaire-ci-runner`。管理 API 導入時は `adlaire-ci-api` も配置する。 |
| 配布形式 | GitHub Release に添付された OS/arch 別の実行バイナリを標準とする。初期標準は Linux x86_64（`linux-amd64`）。 |
| Go toolchain | 利用環境には不要。リリースバイナリをそのまま配置し、利用環境で `go build` しない。 |
| checksum | Release 添付ファイルごとの SHA-256 checksum を取得し、配置前に必ず検証する。 |
| init システム | systemd（Linux） |
| バージョン管理 | GitHub Releases のタグ付き安定版を使用する。利用環境でリポジトリ checkout を更新経路にしない。 |
| ネットワーク | GitHub API への HTTPS 送信。SSH 転送機能を実装した場合のみデプロイ先への SSH 接続。 |

### §26.2 設定変数

スクリプト内で以下の変数をカスタマイズする。

| 変数 | デフォルト値 | 説明 |
|------|------------|------|
| `INSTALL_DIR` | `/opt/adlaire-builder` | インストール先ディレクトリ |
| `BIN_DIR` | `/usr/local/bin` | Go 版バイナリ配置先 |
| `SERVICE_USER` | `root` | systemd サービスの実行ユーザー |
| `VERSION` | —（必須） | セットアップ・アップデート対象の安定版タグ（例：`V.1.100`） |
| `OS_ARCH` | `linux-amd64` | 取得するリリースバイナリの OS/arch。初期標準は `linux-amd64` のみ |
| `DOWNLOAD_DIR` | `/tmp/adlaire-ci-release-$VERSION` | Release 添付ファイルの一時取得先 |

### §26.2a リリース成果物

セットアップ手順は、以下の Release 添付ファイルを取得対象とする。

| 成果物 | 取得タイミング | 説明 |
|--------|----------------|------|
| `adlaire-ci-build-$OS_ARCH` | 初回セットアップ、アップデート | `components/builder.go` から生成した Markdown → 静的 Web サイトビルドバイナリ。 |
| `adlaire-ci-runner-$OS_ARCH` | 初回セットアップ、アップデート | `components/runner.go` から生成した CI ランナーバイナリ。 |
| `adlaire-ci-api-$OS_ARCH` | 管理 API 導入手順、管理 API 導入後のアップデート | `components/api.go` から生成した管理 API サーバーバイナリ。 |
| `admin-ui.tar.gz` | 管理 API 導入手順、管理 API 導入後のアップデート | `admin/index.html` と `admin/adlaire-ci-sdk.js` を含む管理 UI 配布物。 |
| `SHA256SUMS` | Release 添付ファイル取得時 | Release 添付ファイルの SHA-256 checksum 一覧。 |

Release asset 名は上表の文字列と完全一致させる。`$OS_ARCH` は `linux-amd64` だけを初期標準とし、未知 OS/arch を指定した場合は取得前に `unsupported OS_ARCH: {OS_ARCH}` を stderr へ出力して終了コード `2` とする。`SHA256SUMS` は `"{sha256}  {filename}"` 形式の LF 区切り text とし、対象 filename が 1 回だけ出現することを必須とする。対象行が 0 件または 2 件以上の場合は checksum 検証失敗とする。

### §26.2b セットアップ・アップデート機能単位

セットアップ・アップデート実装は、以下の機能単位に分割する。各機能は前段の出力だけを入力として受け取り、失敗時は後続機能を実行しない。

| 機能 | 入力 | 出力 | 失敗条件 | 失敗時の終了状態 |
|------|------|------|----------|------------------|
| Release asset resolver | `VERSION`、`OS_ARCH`、取得対象成果物名、GitHub Release URL | `DOWNLOAD_DIR` 内の取得済みファイル | `VERSION` / `OS_ARCH` 空、HTTP status 非 2xx、取得ファイル 0 byte | 取得済みファイルを配置せず終了 |
| checksum verifier | `SHA256SUMS`、取得済み成果物 | 検証済み成果物一覧 | `SHA256SUMS` 不在、対象行不在、SHA-256 不一致 | バイナリ配置を実行せず終了 |
| binary installer | 検証済みバイナリ、`BIN_DIR` | `adlaire-ci-build`、`adlaire-ci-runner`、必要時 `adlaire-ci-api` | 入力バイナリ不在、実行権限付与失敗、`install` 失敗 | systemd 変更を実行せず終了 |
| secret initializer | PAT 入力、`INSTALL_DIR` | `.github_token` mode `0600` | PAT 空、書き込み失敗、mode 補正失敗 | systemd 変更を実行せず終了 |
| state initializer | `INSTALL_DIR` | `.last_sha`、必要時 `.build_logs/`、`.snapshots/` | 書き込み失敗、mode 補正失敗 | systemd 変更を実行せず終了 |
| systemd unit writer | unit 内容、`SERVICE_USER`、`INSTALL_DIR`、`BIN_DIR` | `/etc/systemd/system/adlaire-ci.service`、`adlaire-ci.timer`、必要時 `adlaire-ci-api.service` | unit 書き込み失敗、`systemctl daemon-reload` 失敗 | enable/start を実行せず終了 |
| service activator | systemd unit 名 | active な timer / service | `enable --now` 失敗、`is-active` 非 `active` | 直前の journal 確認コマンドを出力して終了 |
| admin UI installer | `admin-ui.tar.gz`、`INSTALL_DIR` | `$INSTALL_DIR/admin/index.html`、`$INSTALL_DIR/admin/adlaire-ci-sdk.js` | archive 不在、checksum 不一致、展開後必須ファイル不在 | API service 起動を実行せず終了 |
| rollback executor | `BACKUP_DIR`、`BIN_DIR`、再起動対象 unit | 旧バイナリ復元済み状態 | 旧バイナリ不在、復元失敗、復元後 restart 失敗 | 自動復旧を継続せず journal 確認対象を出力 |

**セットアップ / アップデート共通終了コード：**

| 終了コード | 条件 |
|------------|------|
| `0` | 全手順成功。 |
| `1` | 取得失敗、checksum 不一致、配置失敗、systemd 操作失敗、権限補正失敗、rollback 失敗。 |
| `2` | 変数不正、unsupported OS/arch、必須入力空、既存 credentials あり、実行前検証不合格。 |

各手順は失敗時に固定文言を stderr へ 1 行以上出力する。secret 値、PAT、token、password、Release URL に埋め込まれた認証情報を stderr/stdout に出してはならない。

**Release asset 取得・検証固定契約：**

セットアップ、管理 API 導入、アップデートはいずれも下表の順序で Release asset を扱う。順序を入れ替えてはならない。取得済みファイルは checksum 検証が成功するまで配置対象として扱わない。

| 手順 | 入力 | 成功条件 | 失敗時 |
|------|------|----------|--------|
| 1. 変数検証 | `VERSION`, `OS_ARCH`, `DOWNLOAD_DIR`, 取得対象 asset 名 | 空値なし、`OS_ARCH=linux-amd64`、`DOWNLOAD_DIR` が `/` でない。 | 終了コード `2`。directory 作成、download、配置を行わない。 |
| 2. download dir 作成 | `DOWNLOAD_DIR` | directory が存在し mode `0755` 以上で書込可能。 | 終了コード `1`。配置、systemd 操作を行わない。 |
| 3. asset 取得 | Release URL、asset 名 | HTTP 2xx、取得ファイル size > 0。 | 終了コード `1`。取得済み未検証ファイルを配置しない。 |
| 4. SHA256SUMS 取得 | Release URL | HTTP 2xx、size > 0、LF text。 | 終了コード `1`。asset を配置しない。 |
| 5. 対象行確認 | `SHA256SUMS`、asset 名 | 対象 filename が 1 回だけ出現する。 | 終了コード `1`。0 件、2 件以上はいずれも checksum failure。 |
| 6. checksum 検証 | 対象 asset、対象 SHA-256 | 実ファイル digest が一致する。 | 終了コード `1`。asset を配置しない。 |
| 7. 実行権限付与前確認 | 検証済み binary asset | 通常ファイルであり、directory / symlink ではない。 | 終了コード `1`。配置しない。 |

`admin-ui.tar.gz` は checksum 検証後に一時展開ディレクトリへ展開する。展開後に `index.html` と `adlaire-ci-sdk.js` が同一展開ルート直下に存在しない場合は失敗とし、既存 `$INSTALL_DIR/admin` を変更しない。archive 展開時に絶対 path、`..`、symlink、hardlink、device file を含む entry がある場合は失敗とする。

**配置・権限固定契約：**

| 対象 | 配置方法 | mode | 失敗時 |
|------|----------|------|--------|
| `adlaire-ci-build` | 検証済み asset を `install -m 0755` で `$BIN_DIR/adlaire-ci-build` へ配置する。 | `0755` | systemd restart を行わない。旧 binary がある場合は保持する。 |
| `adlaire-ci-runner` | 検証済み asset を `install -m 0755` で `$BIN_DIR/adlaire-ci-runner` へ配置する。 | `0755` | systemd restart を行わない。旧 binary がある場合は保持する。 |
| `adlaire-ci-api` | 検証済み asset を `install -m 0755` で `$BIN_DIR/adlaire-ci-api` へ配置する。 | `0755` | API service を restart / start しない。 |
| `.github_token` | 一時ファイルへ書込後、mode `0600`、rename、fsync。 | `0600` | systemd unit を変更しない。secret 平文を stderr/stdout に出さない。 |
| `.last_sha` | `{"sha":""}` + LF を一時ファイルへ書込後、mode `0600`、rename、fsync。 | `0600` | systemd unit を変更しない。 |
| `.admin_credentials` | `adlaire-ci-api --init-credentials` の固定手順で生成する。 | `0600` | API service を enable/start しない。 |
| `admin/` | checksum 検証済み archive を一時 directory へ展開し、必須ファイル確認後に差し替える。 | directory `0755`、file `0644` | 既存 `admin/` を変更しない。 |

通常ファイル配置では symlink を最終配置先として許可しない。既存配置先が symlink の場合は `1` で停止し、symlink の参照先を上書きしてはならない。`BIN_DIR`、`INSTALL_DIR`、`DOWNLOAD_DIR` が同一 path、親子関係で危険な組み合わせ、またはいずれかが `/` の場合は `2` で停止する。

**setup / update 失敗時の状態保持契約：**

| 失敗箇所 | 保持するもの | 変更してよいもの | 禁止事項 |
|----------|--------------|------------------|----------|
| download / checksum | 既存 binary、既存 systemd、既存 state、既存 admin UI | `DOWNLOAD_DIR` 内の取得済みファイル | 未検証 asset の配置、service restart。 |
| binary 配置前 | 既存 binary、既存 service 稼働状態 | `DOWNLOAD_DIR` | systemd unit 書換、state 書換。 |
| binary 配置後 / restart 前 | 配置済み新 binary または rollback 対象旧 binary | rollback executor が対象 binary だけ復元してよい。 | state、history、secret、admin UI の巻き戻し。 |
| runner restart 失敗 | `.github_token`、`.last_sha`、history、snapshot、admin UI | build / runner binary の旧版復元、runner restart 1 回 | API credentials や admin UI の変更。 |
| API setup 失敗 | runner binary、runner timer、runner state | API binary、admin 一時展開 directory | runner timer 停止、`.github_token` 変更。 |
| admin UI 差し替え失敗 | 旧 admin UI、API binary、runner state | admin 一時 directory / backup directory | API restart、credentials 変更。 |
| rollback 失敗 | 現在配置済み binary、state、secret | journal 確認対象の報告 | 追加 rollback 推測、state/history/secret 巻き戻し。 |

**セットアップ / アップデート副作用固定契約：**

セットアップ、管理 API 導入、アップデートは、下表の副作用境界を超えてはならない。実装者判断で部分成功を成功報告したり、secret、state、systemd、admin UI をまとめて巻き戻したりしてはならない。

| 段階 | 変更可能対象 | 成功確定条件 | 失敗時固定動作 |
|------|--------------|--------------|----------------|
| 変数検証 | なし | すべての必須変数が空でなく、危険 path でない。 | 終了コード `2`。directory 作成、download、systemd 操作を行わない。 |
| download | `DOWNLOAD_DIR` 配下だけ | 対象 asset と `SHA256SUMS` を取得し、size > 0。 | 既存 binary、state、secret、systemd、admin UI を変更しない。 |
| checksum | `DOWNLOAD_DIR` 配下だけ | 対象 filename が `SHA256SUMS` に 1 回だけ存在し、SHA-256 が一致する。 | 未検証 asset を配置しない。 |
| binary 配置 | `$BIN_DIR` の対象 binary だけ | 通常ファイルへ `0755` で配置し、`--version` が対象 version を返す。 | systemd restart を行わない。配置済み新 binary は rollback 表に従う。 |
| secret / state 初期化 | 対象 secret / state file だけ | 一時ファイル、mode、rename、fsync、親 directory sync が成功する。 | systemd unit を変更しない。secret 平文を出力しない。 |
| admin UI 展開 | admin 一時 directory、成功時のみ `$INSTALL_DIR/admin` | archive 安全検査、必須ファイル確認、差し替えがすべて成功する。 | 既存 admin UI を維持する。API restart を行わない。 |
| systemd unit 配置 | 対象 unit file だけ | unit 書込、mode、`systemctl daemon-reload` が成功する。 | enable / restart / start を行わない。 |
| service 起動 / 再起動 | 対象 unit だけ | `systemctl is-active` が `active`。API は health check も成功。 | rollback 表に従い、追加推測復旧を行わない。 |
| 最終確認 | なし | §26.3、§26.3b、§26.5 の確認項目がすべて成功。 | 成功報告しない。確認失敗箇所と journal 確認対象を出力する。 |

setup / update 実装は、各段階の開始と成功を stderr または stdout に固定文言で 1 行ずつ出してよいが、PAT、password、session token、API token、Webhook secret、SMTP password、Release URL の credential 部分は出力してはならない。secret file が既に存在する場合は、個別手順で上書きを明記している場合を除き、既存値を保持する。特に `.github_token`、`.admin_credentials`、`.webhook_secret`、`.smtp_secret` は、アップデートで自動上書きしない。

`systemctl daemon-reload` 成功だけではセットアップ成功と扱わない。`enable --now`、`restart`、`is-active`、API 導入時の `/api/health` 確認まで完了して初めて成功とする。確認コマンドが利用環境に存在しない場合は、同等確認を実装 PR の検証で実施し、未確認のまま成功扱いにしない。

### §26.3 Go 版初回セットアップ手順

対象は Go 版の `components/builder.go` と `components/runner.go` から生成した `adlaire-ci-build`、`adlaire-ci-runner`、`adlaire-ci.service`、`adlaire-ci.timer` とする。

初回セットアップは以下の停止条件に従う。各手順は直前の手順が成功した場合のみ実行する。失敗時に後続手順を継続してはならない。

| 手順 | 停止条件 | 失敗時の扱い |
|------|----------|--------------|
| 変数検証 | `INSTALL_DIR`、`BIN_DIR`、`VERSION`、`OS_ARCH`、`DOWNLOAD_DIR` が空、`INSTALL_DIR` が `/`、`BIN_DIR` が `/`、`DOWNLOAD_DIR` が `/` | 何も変更せず終了する。 |
| バイナリ取得 | Release バイナリまたは `SHA256SUMS` の取得に失敗、checksum 検証に失敗 | バイナリを配置せず、systemd 設定を変更せず終了する。 |
| バイナリ配置 | checksum 検証済みバイナリが存在しない、または `install` が失敗 | systemd 設定を変更せず終了する。 |
| secret 保存 | PAT が空 | `.github_token` を作成せず終了する。 |
| systemd 配置 | unit ファイル生成または `systemctl daemon-reload` が失敗 | timer を enable せず終了する。 |
| 起動確認 | `systemctl is-active adlaire-ci.timer` が `active` でない | 失敗として扱い、直前のログ確認コマンドを表示する。 |

初回セットアップが中断した場合、作成済みの通常ディレクトリと展開済みリリース資産は自動削除しない。秘密情報ファイルを作成した後に失敗した場合は、`.github_token` の mode が `0600` であることを確認し、mode 補正に失敗した場合はその場で停止する。

**初回セットアップ後の固定確認：**

| 確認 | コマンド | 合格条件 |
|------|----------|----------|
| build binary | `$BIN_DIR/adlaire-ci-build --version` | exit `0`、stdout が `adlaire-ci-build ADLAIRE_CI_SPEC` を含む。 |
| runner binary | `$BIN_DIR/adlaire-ci-runner --version` | exit `0`、stdout が `adlaire-ci-runner ADLAIRE_CI_SPEC` を含む。 |
| PAT file | `stat -c '%a' "$INSTALL_DIR/.github_token"` | `600`。 |
| SHA cache | `cat "$INSTALL_DIR/.last_sha"` | `{"sha":""}` + LF。 |
| timer | `systemctl is-active adlaire-ci.timer` | `active`。 |

確認のいずれかが失敗した場合、セットアップは失敗扱いとする。ただし自動削除や状態ファイル巻き戻しは行わない。

```bash
# ── 変数設定 ──────────────────────────────────────────
INSTALL_DIR="/opt/adlaire-builder"
BIN_DIR="/usr/local/bin"
VERSION="V.1.100"
OS_ARCH="linux-amd64"
DOWNLOAD_DIR="/tmp/adlaire-ci-release-$VERSION"
SERVICE_USER="root"

# ── 1. インストール先作成 ─────────────────────────────
mkdir -p "$INSTALL_DIR"
chmod 0755 "$INSTALL_DIR"

# ── 2. Release バイナリ取得・checksum 検証 ────────────
mkdir -p "$DOWNLOAD_DIR"
cd "$DOWNLOAD_DIR"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$VERSION/adlaire-ci-build-$OS_ARCH"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$VERSION/adlaire-ci-runner-$OS_ARCH"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$VERSION/SHA256SUMS"
grep "  adlaire-ci-build-$OS_ARCH$" SHA256SUMS | sha256sum -c -
grep "  adlaire-ci-runner-$OS_ARCH$" SHA256SUMS | sha256sum -c -

# ── 3. Go 版バイナリ配置 ──────────────────────────────
install -m 0755 "adlaire-ci-build-$OS_ARCH"  "$BIN_DIR/adlaire-ci-build"
install -m 0755 "adlaire-ci-runner-$OS_ARCH" "$BIN_DIR/adlaire-ci-runner"

# ── 4. GitHub PAT 保存 ────────────────────────────────
printf '%s\n' "<PAT>" > "$INSTALL_DIR/.github_token"
chmod 600 "$INSTALL_DIR/.github_token"

# ── 5. SHA キャッシュ初期化 ───────────────────────────
printf '%s\n' '{"sha":""}' > "$INSTALL_DIR/.last_sha"
chmod 600 "$INSTALL_DIR/.last_sha"

# ── 6. systemd サービスファイル配置 ───────────────────
# §26.4.1 のファイル内容を /etc/systemd/system/ に配置した上で:
systemctl daemon-reload

# ── 7. タイマー有効化・起動 ───────────────────────────
systemctl enable --now adlaire-ci.timer

# ── 8. 起動確認 ───────────────────────────────────────
systemctl status adlaire-ci.timer
```

Go 版初回セットアップでは以下を実行しない。

| 対象 | 理由 |
|------|------|
| `/usr/local/bin/adlaire-ci-api --init-credentials --state-dir "$INSTALL_DIR"` | 初回セットアップ対象は runner と build バイナリに限定し、API 認証情報生成は §26.3b で実行する。 |
| `systemctl enable --now adlaire-ci-api` | API service は §26.3b の API バイナリ配置、認証情報生成、unit 配置がすべて成功した後にのみ起動する。 |
| `.build_logs/` 作成 | runner 初期導入ではビルド実行時に必要な状態だけを初期化し、API が参照する履歴ディレクトリは §26.3b で作成する。 |
| `.snapshots/` 作成 | snapshot 参照・rollback API と組み合わせて使うため、§26.3b の管理 API 導入時に作成する。 |

### §26.3b 管理 API 導入後の追加セットアップ手順

`components/api.go`、`admin/index.html`、`admin/adlaire-ci-sdk.js` を実装した後にのみ本手順を実行する。

管理 API 導入手順は、runner の既存稼働状態を壊してはならない。`adlaire-ci-api` の配置、認証情報生成、systemd enable のいずれかが失敗した場合でも、`adlaire-ci.timer` は停止しない。`.admin_credentials` が既に存在する場合は `--init-credentials` を再実行せず、既存 credentials を維持する。

管理 API 導入手順は以下の停止条件に従う。

| 手順 | 停止条件 | 失敗時の扱い |
|------|----------|--------------|
| ディレクトリ作成 | `$INSTALL_DIR/.build_logs`、`$INSTALL_DIR/.snapshots`、`$INSTALL_DIR/admin` の作成に失敗 | runner timer を変更せず終了する。 |
| API バイナリ取得 | `adlaire-ci-api-$OS_ARCH` または `SHA256SUMS` の取得、checksum 検証に失敗 | API バイナリを配置せず終了する。 |
| 管理 UI 取得 | `admin-ui.tar.gz` の取得、checksum 検証、展開に失敗 | API service を起動せず終了する。 |
| 管理 UI 必須ファイル確認 | `$INSTALL_DIR/admin/index.html` または `$INSTALL_DIR/admin/adlaire-ci-sdk.js` が存在しない | API service を起動せず終了する。 |
| API バイナリ配置 | checksum 検証済み API バイナリ不在、または `install` 失敗 | API service を起動せず終了する。 |
| 認証情報生成 | `.admin_credentials` 新規生成に失敗。ただし既存ファイルがある場合は成功扱い | API service を起動せず終了する。 |
| systemd 配置 | unit 書き込みまたは `systemctl daemon-reload` 失敗 | API service を enable/start せず終了する。 |
| 起動確認 | `systemctl is-active adlaire-ci-api` が `active` でない | runner timer を停止せず、API の journal 確認コマンドを出力して終了する。 |

**管理 API 導入後の固定確認：**

| 確認 | コマンド | 合格条件 |
|------|----------|----------|
| API binary | `$BIN_DIR/adlaire-ci-api --version` | exit `0`、stdout が `adlaire-ci-api ADLAIRE_CI_SPEC` を含む。 |
| credentials | `stat -c '%a' "$INSTALL_DIR/.admin_credentials"` | `600`。 |
| admin UI | `test -f "$INSTALL_DIR/admin/index.html"` / `test -f "$INSTALL_DIR/admin/adlaire-ci-sdk.js"` | 両方成功。 |
| API service | `systemctl is-active adlaire-ci-api` | `active`。 |
| local health | `curl -fsS http://127.0.0.1:8765/api/health` | HTTP `200`、JSON object。 |

`curl` が利用できない環境では、Go 実装 PR の検証で `net/http` client または同等のローカル HTTP 確認を行う。未確認のまま API 導入完了扱いにしてはならない。

```bash
# ── 1. 拡張用ディレクトリ作成 ─────────────────────────
mkdir -p "$INSTALL_DIR/.build_logs"
mkdir -p "$INSTALL_DIR/.snapshots"
mkdir -p "$INSTALL_DIR/admin"

# ── 2. Release バイナリ取得・checksum 検証 ────────────
mkdir -p "$DOWNLOAD_DIR"
cd "$DOWNLOAD_DIR"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$VERSION/adlaire-ci-api-$OS_ARCH"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$VERSION/admin-ui.tar.gz"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$VERSION/SHA256SUMS"
grep "  adlaire-ci-api-$OS_ARCH$" SHA256SUMS | sha256sum -c -
grep "  admin-ui.tar.gz$" SHA256SUMS | sha256sum -c -

# ── 3. Go 版 API バイナリ配置 ─────────────────────────
install -m 0755 "adlaire-ci-api-$OS_ARCH" "$BIN_DIR/adlaire-ci-api"

# ── 4. 管理 UI 配布物展開 ────────────────────────────
tar -xzf admin-ui.tar.gz -C "$INSTALL_DIR/admin"
test -f "$INSTALL_DIR/admin/index.html"
test -f "$INSTALL_DIR/admin/adlaire-ci-sdk.js"

# ── 5. 初期認証情報生成（初期パスワード: admin）────────
/usr/local/bin/adlaire-ci-api --init-credentials --state-dir "$INSTALL_DIR"
chmod 600 "$INSTALL_DIR/.admin_credentials"

# ── 6. 管理 API systemd サービス配置 ─────────────────
# §26.4.2 のファイル内容を /etc/systemd/system/adlaire-ci-api.service に配置した上で:
systemctl daemon-reload

# ── 7. サービス有効化・起動 ───────────────────────────
systemctl enable --now adlaire-ci-api

# ── 8. 起動確認 ───────────────────────────────────────
systemctl status adlaire-ci-api
```

### §26.4 systemd サービスファイル

#### §26.4.1 Go 版 runner の systemd ファイル

**`/etc/systemd/system/adlaire-ci.service`**（`components/runner.go`）：

```ini
[Unit]
Description=Adlaire CI Runner

[Service]
Type=oneshot
User=root
WorkingDirectory=/opt/adlaire-builder
ExecStart=/usr/local/bin/adlaire-ci-runner --state-dir /opt/adlaire-builder
```

**`/etc/systemd/system/adlaire-ci.timer`**（`components/runner.go` 定期起動タイマー）：

```ini
[Unit]
Description=Adlaire CI Runner Timer

[Timer]
OnBootSec=1min
OnUnitActiveSec=5min
Unit=adlaire-ci.service

[Install]
WantedBy=timers.target
```

#### §26.4.2 管理 API 導入後の systemd ファイル

**`/etc/systemd/system/adlaire-ci-api.service`**（`components/api.go`）：

```ini
[Unit]
Description=Adlaire CI API Server
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/adlaire-builder
ExecStart=/usr/local/bin/adlaire-ci-api --addr 127.0.0.1:8765 --state-dir /opt/adlaire-builder
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

`User` / `WorkingDirectory` / `ExecStart` のパスは §26.2 の設定変数に合わせて変更する。

systemd unit は上記キー以外を初期標準で追加しない。`Environment=`、`EnvironmentFile=`、`ExecStartPre=`、`ExecStartPost=` を追加する場合は、先に本節へ対象変数、secret 扱い、失敗時挙動を定義する。API service は `127.0.0.1:8765` bind を標準とし、外部公開 bind は本ファイルで未定義のため設定しない。

### §26.5 アップデート手順

`git pull`、利用環境での `go build`、開発ブランチ checkout は使用しない。タグ付き安定版のリリースバイナリを配置し、サービスを再起動する。管理 API を導入していない構成では、管理 API サービスは再起動対象に含めない。

アップデートは以下の順序で実行し、途中失敗時は表の rollback 条件に従う。

| 手順 | 成功条件 | 失敗時 rollback / 停止条件 |
|------|----------|-----------------------------|
| 現在版記録 | 既存バイナリを退避し、退避先パスを保持する。 | 更新を開始しない。 |
| バイナリ取得 | 新 Release バイナリと `SHA256SUMS` の取得、checksum 検証が成功する。 | 退避済み旧バイナリを維持して終了する。 |
| バイナリ更新 | checksum 検証済みの新バイナリを `install -m 0755` で配置できる。 | 退避済み旧バイナリを元へ戻し、サービスを再起動しない。 |
| runner 再起動 | `systemctl restart adlaire-ci.timer` と `systemctl is-active adlaire-ci.timer` が成功する。 | 旧バイナリを戻し、再度 `systemctl restart adlaire-ci.timer` を 1 回だけ実行する。 |
| API 再起動 | API 導入済みの場合のみ `systemctl restart adlaire-ci-api` と `systemctl is-active adlaire-ci-api` が成功する。 | 旧バイナリを戻し、runner と API の再起動を 1 回だけ実行する。 |
| 管理 UI 更新 | API 導入済みの場合のみ `admin-ui.tar.gz` の取得、checksum 検証、一時ディレクトリへの展開、必須ファイル確認、旧 `admin/` との差し替えが成功する。 | 旧 `admin/` を維持または退避先から復元し、API 再起動を実行しない。 |

rollback 後も service が active にならない場合は、自動復旧を継続せず、`journalctl -u adlaire-ci.service -n 100`、API 導入済みなら `journalctl -u adlaire-ci-api -n 100` を確認対象として報告する。rollback はバイナリ差し戻しと service restart のみを行い、状態ファイル、履歴、ログ、secret を巻き戻してはならない。

**アップデート rollback 固定契約：**

| 失敗箇所 | rollback 対象 | rollback 後に実行する確認 | 禁止事項 |
|----------|---------------|----------------------------|----------|
| checksum 検証前 | なし | 旧 service active 確認のみ | 取得済み未検証ファイルを配置しない。 |
| build / runner 配置失敗 | 配置に成功した新バイナリだけ旧版へ戻す。 | `adlaire-ci-build --version`、`adlaire-ci-runner --version` | systemd restart しない。 |
| runner restart 失敗 | build / runner 旧版復元 | `systemctl is-active adlaire-ci.timer` | state、history、secret を戻さない。 |
| API binary 配置失敗 | API 旧版復元。runner は戻さない。 | `adlaire-ci-api --version` | runner service を restart しない。 |
| API restart 失敗 | API 旧版復元 | `systemctl is-active adlaire-ci-api` | `.admin_credentials`、admin UI を戻さない。ただし UI 更新前に失敗した場合。 |
| admin UI 展開失敗 | 旧 `admin/` 維持 | 必須ファイル確認 | API restart しない。 |
| admin UI 差し替え後 API restart 失敗 | 旧 `admin/` 復元、API 旧版復元 | API service active 確認 | runner state を戻さない。 |

rollback は 1 回だけ実行する。rollback 自体が失敗した場合は、追加の推測復旧を行わず、失敗箇所、退避先、現在配置済みファイル、journal 確認コマンドを報告対象として固定する。

**アップデート実行判定固定契約：**

アップデート実装は、以下の判定順で対象を決定する。判定結果を実装者判断で短絡、統合、または省略してはならない。

1. 既存 `$BIN_DIR/adlaire-ci-build` と `$BIN_DIR/adlaire-ci-runner` の存在を確認する。どちらかが不在の場合は終了コード `2` とし、更新を開始しない。
2. API 導入済み判定は `$BIN_DIR/adlaire-ci-api` が通常ファイルとして存在し、`systemctl is-enabled adlaire-ci-api` が `enabled` または `static` を返す場合だけ `true` とする。
3. API 導入済みでない場合、`adlaire-ci-api-$OS_ARCH` と `admin-ui.tar.gz` は取得しない。
4. API 導入済みの場合、build / runner / API binary と admin UI を同じ `NEW_VERSION` の asset から取得する。version 混在は禁止する。
5. すべての対象 asset の checksum 検証が成功するまで、既存 binary、既存 admin UI、systemd unit を変更しない。
6. binary 配置後の version 確認に失敗した場合は、その binary を配置失敗として rollback 対象に含める。
7. runner restart が失敗した場合、API restart と admin UI 更新へ進まない。
8. API restart が失敗した場合、admin UI 更新へ進まない。
9. admin UI 差し替え後に API restart が失敗した場合、admin UI と API binary だけ rollback 対象とし、build / runner binary と runner timer は戻さない。

**アップデート後確認固定契約：**

| 確認 | API 未導入 | API 導入済み |
|------|------------|--------------|
| binary version | `adlaire-ci-build --version`、`adlaire-ci-runner --version` が `NEW_VERSION` を含む。 | 左記に加え `adlaire-ci-api --version` が `NEW_VERSION` を含む。 |
| service | `systemctl is-active adlaire-ci.timer` が `active`。 | 左記に加え `systemctl is-active adlaire-ci-api` が `active`。 |
| admin UI | 確認しない。 | `$INSTALL_DIR/admin/index.html` と `$INSTALL_DIR/admin/adlaire-ci-sdk.js` が存在する。 |
| local API | 確認しない。 | `GET /api/health` が HTTP `200` JSON object を返す。 |
| state preservation | `.github_token`、`.last_sha`、`.build_state`、`.build_history` の mtime と内容が更新対象操作と無関係に変わっていない。 | 左記に加え `.admin_credentials` が存在する場合は mode `600` と内容が保持される。 |

確認失敗時はアップデート失敗として扱う。binary 配置や restart が成功していても、確認失敗を成功報告してはならない。local API 確認で `curl` がない場合は Go 実装 PR の検証で `net/http` client による同等確認を行い、未確認のまま合格扱いにしない。

```bash
# ── 変数設定 ──────────────────────────────────────────
BIN_DIR="/usr/local/bin"
NEW_VERSION="V.2.102"
OS_ARCH="linux-amd64"
DOWNLOAD_DIR="/tmp/adlaire-ci-release-$NEW_VERSION"
BACKUP_DIR="/tmp/adlaire-ci-bin-backup-${NEW_VERSION}"
ADMIN_BACKUP_DIR="/tmp/adlaire-ci-admin-backup-${NEW_VERSION}"
ADMIN_TMP_DIR="/tmp/adlaire-ci-admin-new-${NEW_VERSION}"

# ── 1. 既存バイナリ退避 ──────────────────────────────
mkdir -p "$BACKUP_DIR"
cp "$BIN_DIR/adlaire-ci-build"  "$BACKUP_DIR/adlaire-ci-build"
cp "$BIN_DIR/adlaire-ci-runner" "$BACKUP_DIR/adlaire-ci-runner"
if [ -d "/opt/adlaire-builder/admin" ]; then
  cp -a "/opt/adlaire-builder/admin" "$ADMIN_BACKUP_DIR"
fi

# ── 2. Release バイナリ取得・checksum 検証 ────────────
mkdir -p "$DOWNLOAD_DIR"
cd "$DOWNLOAD_DIR"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$NEW_VERSION/adlaire-ci-build-$OS_ARCH"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$NEW_VERSION/adlaire-ci-runner-$OS_ARCH"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$NEW_VERSION/SHA256SUMS"
grep "  adlaire-ci-build-$OS_ARCH$" SHA256SUMS | sha256sum -c -
grep "  adlaire-ci-runner-$OS_ARCH$" SHA256SUMS | sha256sum -c -

# ── 3. Go 版バイナリ更新 ──────────────────────────────
install -m 0755 "adlaire-ci-build-$OS_ARCH"  "$BIN_DIR/adlaire-ci-build"
install -m 0755 "adlaire-ci-runner-$OS_ARCH" "$BIN_DIR/adlaire-ci-runner"

# ── 4. runner サービス再起動 ─────────────────────────
systemctl restart adlaire-ci.timer

# ── 5. 起動確認 ───────────────────────────────────────
systemctl status adlaire-ci.timer
```

管理 API 導入後は、追加で `adlaire-ci-api` を再起動する。

```bash
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$NEW_VERSION/adlaire-ci-api-$OS_ARCH"
curl -fLO "https://github.com/<owner>/<repo>/releases/download/$NEW_VERSION/admin-ui.tar.gz"
grep "  adlaire-ci-api-$OS_ARCH$" SHA256SUMS | sha256sum -c -
grep "  admin-ui.tar.gz$" SHA256SUMS | sha256sum -c -
install -m 0755 "adlaire-ci-api-$OS_ARCH" "$BIN_DIR/adlaire-ci-api"
mkdir -p "$ADMIN_TMP_DIR"
tar -xzf admin-ui.tar.gz -C "$ADMIN_TMP_DIR"
test -f "$ADMIN_TMP_DIR/index.html"
test -f "$ADMIN_TMP_DIR/adlaire-ci-sdk.js"
if [ -d "/opt/adlaire-builder/admin" ]; then
  mv "/opt/adlaire-builder/admin" "$ADMIN_BACKUP_DIR"
fi
mv "$ADMIN_TMP_DIR" "/opt/adlaire-builder/admin"
systemctl restart adlaire-ci-api
systemctl status adlaire-ci-api
```

### §26.6 サービス操作リファレンス

#### Go 版 runner

| 操作 | コマンド |
|------|---------|
| 状態確認 | `systemctl status adlaire-ci.timer` |
| 起動 | `systemctl start adlaire-ci.timer` |
| 停止 | `systemctl stop adlaire-ci.timer` |
| 再起動 | `systemctl restart adlaire-ci.timer` |
| ログ確認（runner） | `journalctl -u adlaire-ci.service -f` |

#### 管理 API 導入後

| 操作 | コマンド |
|------|---------|
| 状態確認 | `systemctl status adlaire-ci-api` |
| 起動 | `systemctl start adlaire-ci-api` |
| 停止 | `systemctl stop adlaire-ci-api` |
| 再起動 | `systemctl restart adlaire-ci-api` |
| ログ確認（API） | `journalctl -u adlaire-ci-api -f` |

### §26.7 実装受け入れ条件

実装 PR は、下表の受け入れ条件をすべて満たすまで完了扱いにしてはならない。実装対象外のコンポーネントは「未実装」として明記し、合格扱いにしない。

| 対象 | 必須コマンド / 確認 | 合格条件 |
|------|---------------------|----------|
| Go 共通 | `gofmt -l <実装対象Goファイル>` | 実装対象 Go ファイルが存在する場合、出力が空。未作成ファイルはコマンド対象に含めない。 |
| Go test | `go test ./...` | Go module が存在する場合に成功する。Go module が存在しない場合は、その理由を実装完了報告に明記する。 |
| build script | `adlaire-ci-build --src <sample.md> --out <tmp.html>` | exit code `0`、HTML 出力あり、`[REPORT]` の `status` が `success`。 |
| runner | `adlaire-ci-runner --state-dir <tmp-state>` | 必須 secret 未設定時の exit code / ERROR log が §12 と一致し、`.build_lock` が残らない。 |
| API | `POST /api/login`、`GET /api/status`、未知 path、body 禁止 endpoint、JSON 不正、認証なし | §22.0 / §22.0e の status code と body に一致する。 |
| SDK | browser runtime で `login()`、`getStatus()`、`streamBuild()`、HTTP error、timeout を確認する。 | `AdlaireCIError`、`StreamHandle`、token 破棄、timeout が §23 と一致する。 |
| UI | login、manual build、SSE 表示、config 保存、token 発行、logout を確認する。 | §24 の DOM id、disabled、成功表示、失敗表示、再取得、秘密情報消去に一致する。 |
| setup | §26.3 または §26.3b の手順を fresh 環境で実行する。 | unit 配置、権限、`systemctl is-active`、secret mode が仕様どおり。 |
| update | §26.5 の手順を前版バイナリから新 tag のリリースバイナリへ実行する。 | 旧バイナリ退避、新バイナリ配置、restart、失敗時 rollback 条件が仕様どおり。 |
| security | secret 値を含む入力後、stdout、stderr、journal、API response、UI 表示を確認する。 | PAT、Webhook Secret、SMTP password、session token、API token 本体が平文で出ない。 |

**認証 / セットアップ fixture 固定：**

| fixture | 入力 | 合格条件 |
|---------|------|----------|
| credentials init success | 空 state dir で `adlaire-ci-api --init-credentials --state-dir <abs>` | `.admin_credentials` mode `600`、schema 全 key、`must_change=true`、stdout 固定文言。 |
| credentials init existing | `.admin_credentials` 既存 | exit `2`、stderr `credentials already exist`、既存ファイル差分なし。 |
| login success | 初期 password `admin` | `must_change:"prompt"`、TOTP 無効時 token 発行、hash/salt 非表示、`.access_log` 成功行。 |
| login failure lock | password 連続 10 回失敗 | 10 回目後、10 分間 `429 {"error":"Too many attempts"}`、password 詳細非表示。 |
| change password | current 正、新 password 有効 | `.admin_credentials` の salt/hash 更新、`must_change=false`、現 session 以外破棄。 |
| session restart | token 発行後に API process restart | 旧 token は `401`、session file は存在しない。 |
| setup checksum mismatch | Release asset と `SHA256SUMS` 不一致 | バイナリ配置なし、systemd 変更なし、終了コード `1`。 |
| update restart failure | 新バイナリ配置後に service restart 失敗 | 旧バイナリ復元を 1 回だけ行い、state/history/secret は巻き戻さない。 |

**セットアップ / アップデート fixture 固定：**

| fixture | 入力 | 合格条件 |
|---------|------|----------|
| setup variable invalid | `INSTALL_DIR=/`、`BIN_DIR=/`、または `OS_ARCH=darwin-arm64` | 終了コード `2`、download / directory 作成 / 配置なし。 |
| setup download failure | build binary 取得が HTTP `404` | 終了コード `1`、runner binary 取得済みでも配置なし、systemd 変更なし。 |
| setup checksum duplicate | `SHA256SUMS` に同一 asset 行が 2 件 | 終了コード `1`、checksum failure、配置なし。 |
| setup symlink target | `$BIN_DIR/adlaire-ci-build` が symlink | 終了コード `1`、symlink 参照先を上書きしない。 |
| setup pat empty | PAT 入力が空 | 終了コード `2`、`.github_token` 作成なし、systemd 変更なし。 |
| setup success | 正常 asset、正常 PAT、fresh 環境 | binary mode `755`、`.github_token`/`.last_sha` mode `600`、timer active、固定確認すべて成功。 |
| api setup credentials existing | `.admin_credentials` 既存で API 導入 | `--init-credentials` を再実行せず既存 credentials を保持し、API service 起動確認まで進む。 |
| api setup admin archive unsafe | `admin-ui.tar.gz` に `../x` または symlink entry | 終了コード `1`、既存 admin UI 維持、API service start なし。 |
| api setup health failure | API service active だが `/api/health` が非 200 | セットアップ失敗扱い、runner timer は停止しない、journal 確認対象を出力。 |
| update api absent | `adlaire-ci-api` 未導入環境で update | build / runner asset だけ取得・検証・配置し、API asset と admin UI を取得しない。 |
| update checksum before change | update 対象 asset の checksum 不一致 | 既存 binary、admin UI、systemd、state に差分なし。 |
| update runner restart failure | runner restart fake failure | build / runner 旧版復元を 1 回だけ実行し、API restart と admin UI 更新へ進まない。 |
| update api restart failure | API restart fake failure | API binary 旧版復元、runner は戻さない、admin UI 更新へ進まない。 |
| update admin restart failure | admin UI 差し替え後の API restart fake failure | 旧 admin UI と API 旧版を復元し、runner state / history / secret は巻き戻さない。 |
| update rollback failure | 旧 binary 不在で rollback executor 失敗 | 追加復旧せず、失敗箇所、退避先、現在配置済みファイル、journal 確認コマンドを報告対象にする。 |

**運用 API fixture 固定：**

| fixture | 入力 | 合格条件 |
|---------|------|----------|
| snapshot save and prune | `snapshots_keep=2` で build success を 3 回実行 | 最新 2 世代だけ残り、各 snapshot に `site.tar.gz` と `meta.json` が存在する。 |
| snapshot rollback running | `.build_state.running=true` で `POST /api/history/{id}/rollback` | `409 {"error":"Build is running"}`、queue 追加なし、history 追記なし。 |
| maintenance enable no-op | 同一 reason で enable を 2 回実行 | 2 回目は `No changes`、`.config_log` 追記なし。 |
| access-control deny | allow に接続元以外を設定して API 呼び出し | 認証前に `403 {"error":"Forbidden"}`、password/token 検証なし。 |
| hook pre abort | `pre` hook が exit `1`、`abort_on_failure=true` | pipeline 未実行、build status `hook_error`、hook log 保存。 |
| alert duplicate | 同一 alert rule を 2 回作成 | 2 回目は `409 {"error":"Conflict"}`、`.alert_rules` 差分なし。 |
| tag rule invalid | 破損 condition を含む `.tag_rules` で build | build failure、ERROR log `TAG_RULE_INVALID`、SHA cache 更新なし。 |
| verify-output no history | 成功履歴なしで `POST /api/verify-output` | `404 {"error":"Not Found"}`、状態ファイル変更なし。 |
| pipeline config reserved arg | `extra_args:["--src","x"]` | `422`、`.pipeline_config` 差分なし。 |
| notes no-op | 同一 content を 2 回保存 | 2 回目は `No changes`、`.config_log` 追記なし。 |
| smtp secret mask | password 付き `POST /api/smtp-config` 後に GET / backup / log 確認 | password 本体は返らず、mask または `password_set:true` だけ表示。 |
| queue disabled | `queue_max_size=0`、build running 中に `POST /api/build` | `429 {"error":"queue_full"}`、`.build_state.queued` は空。 |
| dashboard duplicate widget | widgets に重複 id を指定 | `422`、`.dashboard_layout` 差分なし。 |

**§22〜§26 API / SDK / UI / 認証 / セットアップ 実装完全性固定契約：**

§22〜§26 のコンポーネントは、各節の本文、endpoint 表、SDK method 表、UI 操作契約、fixture に加えて下表を満たした場合だけ実装完了とする。下表は既存機能の詳細実装を固めるものであり、未定義 endpoint、未定義 UI、未定義認証方式、将来計画機能を追加する根拠にしてはならない。

| 節 | 機能 | 入力 | 出力 | 状態ファイル / 外部副作用 | 失敗時副作用 | 必須 fixture |
|----|------|------|------|---------------------------|--------------|--------------|
| §22.0 | API 共通 | HTTP method/path/header/body、remote addr。 | 固定 status、固定 error body、security header。 | `.api_access_log` 以外は endpoint 契約に従う。 | path/method/body/JSON/auth/scope/validation 失敗時は endpoint 固有処理を開始しない。 | unknown path、method mismatch、body 禁止、JSON 不正、401、403、422、500 mask。 |
| §22.0a〜§22.0c | 状態ファイル / schema | state dir、JSON / JSON Lines / text state。 | typed adapter result、固定初期値、破損時 error。 | atomic write、lock、chmod、fsync、corrupt backup。 | read-only API は状態を修復しない。write 失敗は target を部分更新しない。 | corrupt JSON、unknown key、nullable 違反、lock timeout、chmod failure、GET no write。 |
| §22.0d〜§22.0e | endpoint 契約 | endpoint ごとの request/query/body/path。 | endpoint ごとの response、SDK/UI 対応。 | Read / Write 列に明記された状態だけ扱う。 | 個別 status と共通 error 優先順位に従う。未定義 endpoint を追加しない。 | P0/P1、P2〜P5、pagination、SSE、binary、no-op、partial failure。 |
| §23 | SDK | public method 引数、token、fake fetch response。 | Promise return、`AdlaireCIError`、`StreamHandle`、Blob。 | token は memory のみ。DOM / state file / storage を変更しない。 | `401` だけ token 破棄。`403`、`429`、`500`、network、timeout は token 維持。 | request shape、error shape、timeout、invalid JSON、invalid SSE、body 禁止、401 purge。 |
| §24 | 標準管理 UI | DOM event、form value、SDK return/error。 | DOM 表示、disabled/loading、success/error、secret 消去。 | SDK method だけを呼ぶ。直接 API、状態ファイル、systemd を触らない。 | API 成功前に確定表示しない。失敗時は secret を消し、非 secret 入力は保持する。 | login/TOTP、manual build、stream、refresh failure、secret clearing、disabled priority、direct fetch absence。 |
| §25 | 認証 | password、TOTP code、session token、API token。 | session token、ticket、auth error、access/audit log。 | `.admin_credentials`、`.totp_secret`、memory session/ticket、logs。 | token/ticket は必要ログ成功まで返さない。hash/salt/secret/token 本体を保存しない。 | init、login success/failure、lock、change password、session restart、TOTP replay、audit failure。 |
| §26 | setup/update | release asset、checksum、INSTALL_DIR、BIN_DIR、systemd。 | binary 配置、admin UI、unit、service active、health。 | 検証済み asset だけ配置。secret/state/systemd は段階順に変更する。 | checksum/download/unsafe archive/restart 失敗で後続段階に進まない。rollback は定義範囲だけ 1 回。 | setup invalid、download failure、checksum duplicate、symlink target、API setup、update rollback、health failure。 |

**§22〜§26 横断失敗時副作用固定契約：**

| ケース | 固定結果 |
|--------|----------|
| API validation failure | 状態ファイル、外部 API、systemd、runner、hook、通知を変更しない。`.api_access_log` だけ通常記録対象とする。 |
| API write success / log failure | 個別節が巻き戻しを明記していない限り、保存済み状態は巻き戻さず `500` を返す。 |
| SDK network / timeout | `AdlaireCIError(status=0)` とし、自動 retry、自動 refresh、token 破棄を行わない。 |
| UI refresh failure after success | 操作成功は保持し、再取得失敗だけ panel error に表示する。同じ変更 API を自動再実行しない。 |
| auth log failure before token response | session token、login ticket、API token 本体を response しない。 |
| setup partial failure | 既存 binary、state、secret、admin UI、systemd を、表で許可した対象以外は変更しない。 |
| update rollback failure | 追加推測復旧を行わず、失敗箇所、退避先、現在配置済みファイル、journal 確認対象を報告する。 |

**§22〜§26 実装前・実装後確認固定契約：**

| 段階 | 確認 | 合格条件 |
|------|------|----------|
| 実装前 | endpoint / SDK / UI 対応 | §22.0e の API、§23 の SDK method、§24 の UI 操作が同一機能でそろっている。欠落時は先に仕様改訂する。 |
| 実装前 | state schema | 使用する状態ファイルが §22.0a / §22.0c にあり、型、初期値、破損時処理、更新責務が定義済み。 |
| 実装前 | secret handling | secret 値の保存先、mask、response 禁止、log 禁止、UI 消去条件が定義済み。 |
| 実装後 | common error | unknown path、method mismatch、body 禁止、JSON 不正、401、403、422、500 が固定 body と一致する。 |
| 実装後 | state side effect | 成功、validation failure、conflict、write failure、log failure の状態差分が fixture expected と一致する。 |
| 実装後 | client behavior | SDK error、UI disabled、success/error、refresh、secret 消去、direct fetch 不在が固定契約どおり。 |
| 実装後 | setup/update | checksum、unsafe archive、restart failure、rollback failure、health failure が §26 fixture と一致する。 |

Phase 別の実装受け入れ条件は以下とする。

| Phase | 必須確認 | 合格条件 |
|-------|----------|----------|
| Phase 1 | §0g.1、§8a、§26.7 build script | `adlaire-ci-build` の CLI、出力構造、report、fixture、冪等性が固定され、Phase 2 が追加判断なしに呼び出せる。 |
| Phase 2 | §0g.2、§15a、§26.7 runner | runner 状態ファイル、履歴、ログ、pending queue、snapshot、通知、lock、終了コードが固定され、Phase 3 が状態ファイルを追加判断なしに読める。 |
| Phase 3 | §0g.3、§22.0f P0/P1、§25、§26.7 API | API 共通契約、認証、session、P0/P1 endpoint、error body、状態ファイル read/write が固定される。 |
| Phase 4 | §0g.4、§22.0f P2〜P5、§22.0e | 全 endpoint 契約が SDK 実装に渡せる粒度で固定され、P0/P1 の互換を壊していない。 |
| Phase 5 | §0g.5、§23、§26.7 SDK | 全 SDK method が固定済み endpoint 契約に一致し、Phase 6 が UI 実装に使える error / stream / return 契約を持つ。 |
| Phase 6 | §0g.6、§24、§26.7 UI/security | 全 UI 操作が SDK 経由で成立し、成功/失敗/disabled/loading/secret 消去が仕様どおりである。 |

Phase 別の実装検証記録は以下の単位で行う。

| Phase | 必須記録 | 失敗時の扱い |
|-------|----------|--------------|
| Phase 1 | CLI 引数、fixture A〜D、生成物一覧、`[REPORT]`、冪等性、strict / non-strict の結果。 | `components/builder.go` を完了扱いにせず、Phase 2 着手禁止。 |
| Phase 2 | secret 不足、lock、GitHub fake、pipeline fake、deploy fake、snapshot、notify、状態ファイル schema の結果。 | `components/runner.go` を完了扱いにせず、Phase 3 着手禁止。 |
| Phase 3 | P0 / P1 endpoint、認証、session、error body、SSE、状態 read/write、systemd service の結果。 | API P0 / P1 を完了扱いにせず、Phase 4 着手禁止。 |
| Phase 4 | P2〜P5 endpoint、secret mask、rollback、maintenance、hook、token、P0 / P1 互換確認の結果。 | endpoint 契約を固定扱いにせず、Phase 5 着手禁止。 |
| Phase 5 | method 対応表、fake fetch、HTTP error、timeout、stream、`401` token 破棄、body 禁止の結果。 | SDK 契約を固定扱いにせず、Phase 6 着手禁止。 |
| Phase 6 | DOM id、panel、SDK 呼び出し、success/error/loading/disabled、stream、secret 消去、直接 API 呼び出し不存在の結果。 | 初期実装完了扱いにせず、未充足 UI 仕様を解消する。 |

Phase 別受け入れ条件のいずれかが未実行、失敗、または環境都合で省略された場合、その Phase を完了扱いにしてはならない。後続 Phase の実装 PR を開始する前に、先行 Phase の未充足条件を仕様または実装で解消する。

**Phase 完了判定 fixture 記録固定：**

| Phase | fixture 記録 | fake 記録 | 期待値記録 |
|-------|--------------|-----------|------------|
| Phase 1 | 使用した `testdata/builder/...` path、入力 Markdown、テーマ設定、asset 入力。 | fake filesystem を使った場合のみ、注入した失敗と対象 path。 | 生成物一覧、期待 HTML / CSS / JS / search index、stdout / stderr、終了コード、冪等性結果。 |
| Phase 2 | 使用した `testdata/runner/r*` path、初期 `state/`、GitHub 応答、pipeline / deploy / notify fixture。 | fake GitHub server、fake ssh executable、fake notifier、fake filesystem の呼び出し記録。 | `.last_sha`、queue、snapshot、通知結果、retry / circuit breaker、状態ファイル更新有無。 |
| Phase 3 | 使用した `testdata/api/p0-p1/...` path、request、初期 `state/`。 | API handler への request 記録、SSE frame 記録、状態 read/write 記録。 | HTTP status、response body、error body、SSE event、状態ファイル更新後 expected。 |
| Phase 4 | 使用した `testdata/api/p2-p5/...` path、request、secret mask fixture。 | hook / notify / token / maintenance の fake 呼び出し記録。 | HTTP status、response body、mask 後 payload、token 再取得不可結果、P0 / P1 回帰結果。 |
| Phase 5 | 使用した `testdata/sdk/...` path、fake fetch transcript、stream chunk。 | fake fetch の method、URL、headers、body 有無、abort / timeout 記録。 | SDK return、`AdlaireCIError`、stream event、token 破棄、body 禁止 endpoint の request 不成立。 |
| Phase 6 | 使用した `testdata/ui/...` path、初期 DOM、操作手順。 | fake SDK の method、引数、呼び出し順、stream unsubscribe 記録。 | DOM assertion、success / error / loading / disabled、secret 消去、直接 `fetch()` 不存在。 |

**未実行検証の代替条件：**

| 未実行対象 | 代替として認める確認 | 完了扱い |
|------------|----------------------|----------|
| 実 GitHub API、実 SSH、実通知先 | §0g.8 の fake GitHub server、fake ssh executable、fake notifier の fixture 検証。 | 認める。実外部接続は Phase 完了条件に含めない。 |
| systemd 実起動 | service file の内容確認、起動 command / user / working directory / environment の静的確認、API handler の Go test。 | 実 target host での起動確認が仕様化された Phase では `未実行`。仕様化されていない場合は静的確認で可。 |
| browser 実機操作 | Vanilla JS が動作する browser runtime または同等 DOM 環境での fake SDK / DOM assertion。 | 認める。ただし直接 `fetch()` 不存在、secret 消去、loading / disabled は必須。 |
| `go test` または `gofmt -l` が実行不能 | 実行不能理由、未実行 command、再実行条件の記録。 | 認めない。該当 Phase は `未実行` を含むため完了扱い不可。 |
| fixture expected の欠落 | 欠落 path、期待値を定義できない理由の記録。 | 認めない。fixture expected を追加するまで完了扱い不可。 |

**fixture 期待値更新固定：**

| 更新対象 | 更新条件 | 必須確認 |
|----------|----------|----------|
| `expected/` 内の生成物 | 仕様本文の出力契約、schema、error body、DOM id、SDK return のいずれかが変更された場合のみ更新する。 | 仕様本文の該当節と fixture expected が同じ値を示すこと。 |
| fake transcript | 外部依存の呼び出し method、path、payload、retry、timeout の仕様が変更された場合のみ更新する。 | secret / token / password 原文が transcript に存在しないこと。 |
| DOM assertion | UI の DOM id、panel、表示文言、disabled / loading / success / error 条件が変更された場合のみ更新する。 | SDK method 対応表と DOM assertion が一致すること。 |
| error expected | HTTP status、exit code、`AdlaireCIError.code`、stderr prefix が変更された場合のみ更新する。 | 正常系 fixture と異常系 fixture の両方で期待値が固定されていること。 |

Phase 完了判定テンプレートは以下とする。実装 PR 本文では、対象 Phase ごとに本テンプレートの項目を埋める。

```text
## Phase 完了判定

- 対象 Phase:
- owner component:
- collaborator component:
- 実装対象ファイル:
- 追加 fixture / testdata:
- 固定契約:
- 実装対象外:
- 後続 Phase への影響:

| 対象 | コマンド / 確認 | 期待結果 | 実結果 | 判定 |
|------|------------------|----------|--------|------|
| <対象> | <実行内容> | <仕様上の期待結果> | <実際の結果> | PASS / FAIL / 未実行 |
```

`判定` が `FAIL` または `未実行` の行を含む場合、その Phase は完了扱いにしてはならない。環境制約により確認できない項目がある場合も `未実行` とし、完了扱いにするには代替検証を仕様化してから再実行する。

受け入れ結果は、実装 PR 本文に `対象 / コマンド / 期待結果 / 実結果 / 判定` の形式で記録する。失敗、未実行、環境都合で省略した項目がある場合、そのコンポーネントを完了扱いにしてはならない。

---

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
| `setup` | release binary 配置、systemd、初期化、update、rollback、既存 secret 保持を固定する。 | §26 | `runner`、`api`、`statefile` | runtime 機能追加、状態 schema の暗黙変更、外部依存追加。 |

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

**§27 fixture 配置・命名固定契約：**

§27 の fixture は、実装者が実行順や期待値を推測しないように、下表の単位で配置する。実装ファイル名やテストフレームワークは本節では固定しないが、fixture の入力、期待出力、期待副作用は機能単位で分離する。

| 機能範囲 | fixture 配置単位 | 必須内容 | 禁止事項 |
|----------|------------------|----------|----------|
| §27.1〜§27.4 | `status/`、`dry-run/`、`retry/`、`output-meta/` 相当の機能別単位。 | runner 入力、GitHub fake response、builder 入力、期待 log/history/status/report。 | 実際の GitHub API、実時刻依存の期待値、secret 平文。 |
| §27.5〜§27.11 | `config/`、`access-log/`、`archive/`、`schedule/` 相当の API 機能別単位。 | HTTP request、初期状態、期待 response、期待状態差分、期待 log。 | API response だけの検証、状態差分未確認。 |
| §27.12〜§27.20 | `webhook/`、`stats/`、`snapshot/`、`health/`、`branch-config/`、`summary/`、`config-log/` 相当の機能別単位。 | 署名、payload、query、snapshot 入力、破損行、期待 paging、期待 rollback。 | 署名検証省略、破損行の黙殺仕様未確認。 |
| §27.21〜§27.38 | `runner-extensions/` 配下の機能別単位。 | target、pipeline、cache、queue、hook、remote、approval、notification、trend の正常/異常/部分失敗。 | 実行完了順依存、外部 shell 展開、未定義状態ファイル。 |
| §27.42〜§27.47 | `security/` 配下の scope、token、audit、session、totp、rate-limit 単位。 | route 判定、body 未評価、token hash、audit failure、window reset、secret mask。 | token 本体保存、Authorization header 保存、監査なし権限拒否。 |

fixture 名は `success-*`、`failure-*`、`partial-*`、`noop-*`、`security-*` のいずれかで始める。fixture 名に実行時刻、乱数、環境依存 path、実 token 値を含めてはならない。期待時刻は固定値を使い、現在時刻依存の検証では fake clock を fixture 入力に含める。

**§27 fixture カタログ固定契約：**

§27 の実装 PR は、対象機能について下表の fixture を作成する。複数節を同一 PR で実装する場合は、各節の fixture を省略せず作成する。fixture は入力、期待 response、期待 stdout/stderr、期待状態差分、期待 log/history/audit/notify、secret 非表示確認のうち該当するものを含める。

| 節 | fixture 名 | 入力 fixture | 期待出力 / 期待副作用 |
|----|------------|--------------|------------------------|
| §27.1 | `success-commit-status-pending-success`、`failure-commit-status-unavailable-sha`、`failure-commit-status-api-error`、`noop-commit-status-disabled` | server config、commit SHA、fake GitHub response、build result。 | GitHub payload、build log `commit_status`、history state、token mask、disabled 時呼び出し 0。 |
| §27.2 | `success-dry-run-changed`、`noop-dry-run-unchanged`、`failure-dry-run-github-error`、`security-dry-run-secret-mask` | CLI args、state dir、fake GitHub response。 | stdout JSON、終了コード、状態差分なし、lock なし、secret mask。 |
| §27.3 | `success-retry-after-rate-limit`、`success-retry-after-timeout`、`failure-retry-limit-exceeded`、`noop-retry-nonretryable` | retry config、attempt sequence、pipeline / deploy fake result。 | attempts 配列、retry_count、最終 status、SHA 更新有無、未実行 attempt 不作成。 |
| §27.4 | `success-output-meta-html-report-api`、`success-output-meta-empty-values`、`failure-output-meta-invalid-sha`、`failure-output-meta-invalid-time` | builder args、Markdown 入力、build meta 値。 | HTML meta、REPORT、output-meta response、終了コード、既存出力保護。 |
| §27.5 | `success-config-validate-valid`、`success-config-validate-invalid`、`failure-config-validate-unknown-key`、`security-config-validate-secret-mask` | config JSON、既存 `.server_config`。 | valid/errors/warnings、状態差分なし、unknown key `422`、secret 非表示。 |
| §27.6 | `success-api-access-log-authenticated`、`success-api-access-log-unauthorized`、`failure-api-access-log-append`、`security-api-access-log-body-mask` | HTTP request、actor、target endpoint result。 | `.api_access_log` 追記、status/duration、body 未保存、append failure 挙動。 |
| §27.7 | `success-log-archive`、`noop-log-archive-empty`、`failure-log-archive-gzip`、`success-log-cleanup` | build logs、retention config、archive target。 | gzip archive、元 log 維持/削除条件、cleanup response、破損 log skip。 |
| §27.8 | `success-build-status-running`、`success-build-status-final`、`failure-build-status-write`、`failure-build-status-corrupt-api` | runner state、build result、queue/pending/circuit state。 | `.build_status.json`、status API、finalizer、write failure 時 log/history 維持。 |
| §27.9 | `success-trigger-manual`、`success-trigger-webhook`、`success-trigger-approval`、`failure-trigger-filter-invalid` | trigger source、history query、build log。 | trigger enum 保存、history filter、未知 trigger 不保存、`422`。 |
| §27.10 | `success-startup-integrity-clean`、`success-startup-integrity-recovered`、`failure-startup-integrity-unrecoverable`、`security-startup-integrity-secret-mode` | startup state files、mode、corrupt files。 | recovery record、WARN/ERROR、build 開始可否、secret file mode 補正。 |
| §27.11 | `success-schedule-interval`、`success-schedule-pause-resume`、`failure-schedule-systemd-update`、`noop-schedule-same-value` | schedule request、server config、fake systemd。 | `.server_config`、systemd result、config log、保存済み config の再取得値。 |
| §27.12 | `success-webhook-push-queued`、`failure-webhook-invalid-signature`、`noop-webhook-duplicate-delivery`、`failure-webhook-queue-full` | GitHub headers、raw body、secret、queue state。 | webhook event、queue entry、状態差分なし条件、重複防止。 |
| §27.13 | `success-webhook-events-page`、`success-webhook-events-empty`、`partial-webhook-events-corrupt-line`、`failure-webhook-events-invalid-query` | `.webhook_events.json`、query。 | events/total、破損行除外、secret 非表示、`422`。 |
| §27.14 | `success-stats-summary`、`success-stats-timeline`、`partial-stats-corrupt-log-skip`、`failure-stats-invalid-query` | history、build logs、archive logs、query。 | stats response、rounding、WARN、read-only 差分なし。 |
| §27.15 | `success-snapshot-list-download`、`success-snapshot-delete`、`success-snapshot-rollback`、`failure-snapshot-running-conflict` | `.snapshots/`、history id、running state。 | tar.gz download、delete 差分、rollback history/log、不正 entry 防止。 |
| §27.16 | `success-health-ok`、`success-health-degraded`、`failure-health-read-error`、`noop-health-readonly` | status/pending/notify/process state。 | health JSON、HTTP status、degraded 判定、状態差分なし。 |
| §27.17 | `success-log-search-level`、`success-log-search-archive`、`partial-log-search-corrupt-skip`、`failure-log-search-invalid-level` | build logs、archive logs、query。 | search results、line_number、破損除外、`422`。 |
| §27.18 | `success-branch-config-get-default`、`success-branch-config-post`、`failure-branch-config-invalid-path`、`partial-branch-config-log-failure` | branch config request、existing config。 | branch_targets、config log、validation 差分なし、保存済み状態維持。 |
| §27.19 | `success-weekly-summary-auto`、`success-weekly-summary-manual`、`noop-weekly-summary-same-day`、`failure-weekly-summary-send` | notify config、history、fake webhook。 | payload、notify log、sent date 更新条件、失敗時 pending。 |
| §27.20 | `success-config-diff-simple`、`success-config-diff-nested`、`noop-config-diff-same-value`、`security-config-diff-secret-mask` | before/after config、request。 | diff/diff_text、config log、no-op 差分なし、secret mask。 |
| §27.21 | `success-multi-file-one-change`、`success-multi-file-many-change`、`noop-multi-file-all-skip`、`failure-multi-file-path-traversal` | branch target、target files、SHA cache。 | build 対象集合、対象別 SHA 更新、重複排除、path error。 |
| §27.22 | `success-yaml-pipeline-file-priority`、`success-yaml-pipeline-inline`、`failure-yaml-pipeline-parse`、`security-yaml-pipeline-secret-mask` | `.pipeline.yml`、inline YAML、step fake result。 | pipeline_source、step log、終了コード、deploy/SHA 更新なし条件。 |
| §27.23 | `success-local-watch-change`、`noop-local-watch-no-change`、`failure-local-watch-state-corrupt`、`failure-local-watch-tag-filter-conflict` | local files、watch state、server config。 | local watch state、trigger、GitHub call 0、終了コード `2`。 |
| §27.24 | `success-tag-filter-match`、`noop-tag-filter-unmatched`、`failure-tag-filter-api`、`failure-tag-filter-pattern` | tag refs、patterns、SHA cache。 | matched tags、SHA 更新条件、skip 挙動、不正 pattern。 |
| §27.25 | `success-build-cache-hit`、`success-build-cache-miss`、`partial-build-cache-byte-mismatch`、`failure-build-cache-save` | Markdown、cache index/pages、theme/version。 | byte 一致、cache atomic save、hit 破棄、build 成否維持。 |
| §27.26 | `success-parallel-targets-all`、`partial-parallel-targets-some-fail`、`failure-parallel-targets-all-fail`、`success-parallel-targets-order-stable` | target list、parallel result sequence。 | 設定順 result、overall status、timeout、history/status。 |
| §27.27 | `success-hook-pre-post`、`failure-hook-pre-abort`、`partial-hook-post-fail`、`security-hook-shell-denied` | hook config、command_args、fake command result。 | hook log、secret mask、pre abort、shell 展開禁止。 |
| §27.28 | `success-dependency-manifest`、`success-dependency-missing`、`failure-dependency-build-keeps-old`、`security-dependency-path-normalize` | Markdown refs、existing manifest、build result。 | dependency manifest、成功時置換、失敗時旧 manifest 維持。 |
| §27.29 | `success-remote-build-artifact`、`failure-remote-build-auth`、`failure-remote-build-checksum`、`security-remote-build-unsafe-archive` | remote config、artifact archive、manifest/checksum。 | 一時展開、検証後 deploy、既存出力保護、unsafe entry 拒否。 |
| §27.30 | `success-approval-approve`、`noop-approval-reject`、`noop-approval-expire`、`failure-approval-double-approve` | approval queue、API action、clock。 | 状態 enum、queue 連携、物理削除なし、audit。 |
| §27.31 | `success-branch-env-match`、`noop-branch-env-default`、`security-branch-env-secret-mask`、`failure-branch-env-unknown-key` | branch config env、build command fake。 | child env、log key 名、secret mask、未定義 key 除外。 |
| §27.32 | `success-notify-build-success`、`success-notify-build-failure`、`partial-notify-send-failure`、`security-notify-secret-mask` | notify config、build result、fake channel。 | notify log/pending、build 成否維持、payload mask。 |
| §27.33 | `success-trend-append`、`success-trend-replace-same-id`、`partial-trend-retention`、`failure-trend-save` | build history/log、trend state。 | sample upsert、summary、retention、終了コード最低 `1`。 |
| §27.34 | `success-chain-linear`、`success-chain-parallel`、`failure-chain-cycle`、`partial-chain-required-fail` | chain config、job result sequence。 | chain_run_id、job log、skipped 条件、validation `422`。 |
| §27.35 | `success-priority-queue-order`、`success-priority-queue-fifo-tie`、`failure-priority-queue-full`、`partial-priority-queue-dequeue-failure` | queued entries、runner lock、API action。 | dequeue atomic、running 設定、entry 残存、`429`。 |
| §27.36 | `success-failure-category-timeout`、`success-failure-category-network`、`success-failure-category-deploy`、`noop-failure-category-success-null` | exit code、stderr、API error、status。 | category enum、unknown fallback、success null、history/log。 |
| §27.37 | `success-environment-full`、`partial-environment-null-fields`、`security-environment-secret-path`、`failure-environment-version-missing` | OS/env/version/state path inputs。 | environment object、null field、basename 保存、secret path 非表示。 |
| §27.38 | `success-duration-anomaly-avg`、`success-duration-anomaly-p95`、`noop-duration-anomaly-sample-shortage`、`partial-duration-anomaly-notify-failure` | trend summary、duration config、build result。 | anomaly tag、history flag、通知 event、重複なし。 |
| §27.42 | `security-scope-trigger-allowed`、`security-scope-read-denied`、`security-scope-path-param`、`failure-scope-audit-failure` | route/method、token record、request body。 | body 未評価、`403`、audit、Authorization 非保存。 |
| §27.43 | `security-token-create-once`、`security-token-list-mask`、`success-token-revoke`、`failure-token-expired-auth` | token request、token state、clock。 | token 1 回表示、hash 保存、revoke、期限切れ `401`。 |
| §27.44 | `success-audit-operation`、`success-audit-denied`、`failure-audit-append`、`security-audit-secret-mask` | actor、target、result、operation input。 | JSON Lines、mask、必須 audit failure、paging。 |
| §27.45 | `success-session-active`、`failure-session-expired`、`success-session-timeout-update`、`success-session-revoke-all` | session state、clock、timeout config。 | last_seen/expires_at、`401`、revoke all、token 非表示。 |
| §27.46 | `security-totp-setup-once`、`success-totp-confirm`、`failure-totp-code-reuse`、`success-totp-disable` | setup ticket、TOTP code、clock、secret state。 | secret 有効保存条件、ticket 一回使用、window、disable。 |
| §27.47 | `security-rate-limit-login`、`success-rate-limit-window-reset`、`security-rate-limit-ip-actor`、`failure-rate-limit-state-save` | policy、rate state、RemoteAddr、actor。 | count、`429`、audit 成功条件、部分 count 更新なし。 |

**§27 fixture ファイルセット固定契約：**

各 fixture は、下表のファイルセットを持つ。該当しない入出力は `not-applicable.txt` を置くのではなく、`manifest.json` の `not_applicable` 配列に理由付きで記録する。実装者は fixture ごとに必要ファイルを推測してはならない。

| ファイル | 必須 | 内容 | 禁止事項 |
|----------|------|------|----------|
| `manifest.json` | 必須 | fixture 名、対象節、機能名、分類、fake clock、owner component、collaborator component、参照仕様節、not_applicable 理由。 | 実行環境依存 path、乱数、実 secret。 |
| `input/request.json` | API / SDK / UI fixture で必須 | HTTP method、path、query、headers、body、SDK method、SDK args、UI action。 | Authorization header の実 token、secret 平文。 |
| `input/cli.json` | runner / builder fixture で必須 | binary 名、argv、env key、cwd、stdin、fake clock。 | 実 home path、実 credential path。 |
| `input/state/` | 状態参照 fixture で必須 | 実行前状態ファイル一式。存在しない状態は `manifest.json` の `missing_state` に列挙する。 | 期待状態を混ぜること、実 secret。 |
| `input/files/` | builder / artifact / hook / remote fixture で必須 | Markdown、YAML、archive、snapshot、hook 入力などの対象ファイル。 | 実外部サービスから取得した未固定ファイル。 |
| `input/fakes.json` | 外部 API / command fixture で必須 | fake GitHub、fake SSH、fake SMTP、fake webhook、fake systemd、fake command の応答順。 | 実ネットワーク呼び出し前提。 |
| `expected/response.json` | HTTP / SDK fixture で必須 | HTTP status、headers、JSON body、SDK return/error。 | 未定義 key、順序非決定配列。 |
| `expected/stdout.txt` | CLI stdout がある fixture で必須 | stdout 完全一致。stdout なしは空ファイル。 | 現在時刻、絶対環境 path。 |
| `expected/stderr.txt` | CLI stderr がある fixture で必須 | stderr 完全一致。stderr なしは空ファイル。 | secret、実 token、実 URL credential。 |
| `expected/state/` | 状態差分がある fixture で必須 | 実行後状態ファイル一式、または `state-diff.json`。 | 期待しないファイルの混入。 |
| `expected/logs/` | log / history / audit / notify fixture で必須 | build log、history、access log、audit log、notify log の期待値。 | secret 平文、実 Authorization header。 |
| `expected/effects.json` | 必須 | 外部 API 呼び出し、command 実行、通知送信、download/stream 中断、呼び出し 0 件の期待値。 | 呼び出し順未指定、実外部送信。 |
| `expected/security.json` | secret / auth / rate limit fixture で必須 | secret 非表示確認対象、禁止文字列、token hash 検証、scope 判定、rate count。 | secret を検証用に平文保存すること。 |

`expected/state/` は、fixture が検証対象とする状態ファイルだけを含める。変更してはならない状態ファイルは `expected/effects.json` の `unchanged_paths` に列挙する。削除されるべきファイルは `expected/effects.json` の `deleted_paths` に列挙し、空 directory の存在可否も明記する。

**§27 fixture manifest schema 固定契約：**

`manifest.json` は次の schema に従う。未知 key は禁止する。

```json
{
  "name": "success-example",
  "section": "27.1",
  "feature": "commit_status",
  "category": "success",
  "owner_component": "commitstatus",
  "collaborator_components": ["runner", "statefile"],
  "components": ["commitstatus", "runner", "statefile"],
  "references": ["§27.1", "§22.0a"],
  "fake_clock": "2026-09-16T00:00:00Z",
  "not_applicable": [
    { "path": "input/request.json", "reason": "CLI fixture" }
  ],
  "missing_state": [
    ".build_status.json"
  ],
  "assertions": [
    "state",
    "logs",
    "effects",
    "secret-mask"
  ]
}
```

| key | 型 | 必須 | 許容値 |
|-----|----|------|--------|
| `name` | string | 必須 | §27 fixture カタログ固定契約に記載された fixture 名。 |
| `section` | string | 必須 | `27.1`〜`27.38`、`27.42`〜`27.47`。 |
| `feature` | string | 必須 | lowercase snake_case。 |
| `category` | string | 必須 | `success`、`failure`、`partial`、`noop`、`security`。fixture 名 prefix と一致する。 |
| `owner_component` | string | 必須 | `builder`、`runner`、`api`、`sdk`、`ui`、`statefile`、`archive`、`commitstatus`、`setup`、`security` のいずれか 1 件。 |
| `collaborator_components` | array[string] | 必須 | owner 以外の component。該当なしは空配列。 |
| `components` | array[string] | 必須 | `owner_component` と `collaborator_components` を重複なしで含む配列。許容値は `builder`、`runner`、`api`、`sdk`、`ui`、`statefile`、`archive`、`commitstatus`、`setup`、`security`。 |
| `references` | array[string] | 必須 | 参照仕様節。対象 §27.x と関連 §22 / §23 / §24 / §25 / §26 を含める。 |
| `fake_clock` | string/null | 必須 | UTC ISO 8601 または `null`。時刻依存 fixture は `null` 禁止。 |
| `not_applicable` | array[object] | 必須 | 該当しない必須候補ファイルと理由。空配列可。 |
| `missing_state` | array[string] | 必須 | 実行前に存在しないことを期待する状態ファイル。空配列可。 |
| `assertions` | array[string] | 必須 | `response`、`stdout`、`stderr`、`state`、`logs`、`effects`、`secret-mask`、`order`、`idempotency`、`no-write` の 1 件以上。 |

**§27 fixture 合否判定固定契約：**

| 判定 | 合格条件 |
|------|----------|
| response | status、headers、body、error details、request id が期待値と一致する。 |
| stdout / stderr | 改行を含め完全一致する。時刻や path は fake 値だけを使う。 |
| state | 期待対象状態ファイルが byte 等価または `state-diff.json` と一致する。未列挙状態ファイルに差分がない。 |
| logs | JSON Lines は行順、key、値、末尾改行が一致する。破損行 fixture では破損行を修復しない。 |
| effects | 外部 API、外部 command、通知、download、stream の呼び出し回数、順序、payload が一致する。 |
| secret-mask | 禁止文字列が response、stdout/stderr、state、logs、effects、UI DOM に存在しない。 |
| order | 複数状態更新は仕様の保存順と一致する。途中失敗 fixture は失敗地点以降の副作用がない。 |
| idempotency | 同一 fixture を 2 回適用した場合、2 回目の差分が no-op 仕様と一致する。 |
| no-write | read-only / dry-run fixture で状態、logs、effects の差分がない。 |

上表のうち `manifest.json.assertions` に含まれる判定が 1 つでも失敗した場合、その fixture は失敗とする。対象機能の必須 fixture が 1 件でも存在しない、または skip された場合、その機能の実装 PR は未完了とする。

**§27 fixture assertion 選択固定契約：**

fixture の `manifest.json.assertions` は、実装者が任意に減らしてはならない。fixture 名 prefix、owner component、collaborator component に応じて、下表の assertion を必ず含める。個別節で追加検証が必要な場合は、下表へ追加してから fixture を作成する。

| 条件 | 必須 assertion |
|------|----------------|
| `success-*` | `state`、`logs`、`effects`。HTTP / SDK fixture では `response` も必須。CLI fixture では `stdout` または `stderr` の少なくとも一方を必須とする。 |
| `failure-*` | `response` または `stdout` / `stderr`、`state`、`effects`、`order`。状態差分なしを期待する場合は `no-write` も必須。 |
| `partial-*` | `state`、`logs`、`effects`、`order`。どの副作用まで完了し、どこから未実行かを `expected/effects.json` で明示する。 |
| `noop-*` | `response` または `stdout`、`no-write`、`idempotency`、`effects`。外部呼び出し 0 件を `expected/effects.json` に明記する。 |
| `security-*` | `response` または `stdout` / `stderr`、`secret-mask`、`effects`。認証 / scope / rate limit / TOTP / token fixture では `state` も必須。 |
| `components` に `api` を含む | `response`、`state`、`effects`。read-only API は `no-write` も必須。 |
| `components` に `sdk` を含む | `response`、`effects`。`401` fixture では token 破棄の期待値を `expected/state/` または `expected/effects.json` に含める。 |
| `components` に `ui` を含む | `response`、`effects`、`secret-mask`。DOM 期待値または UI action 後の field 消去期待を含める。 |
| `components` に `runner` を含む | `state`、`logs`、`effects`、`order`。dry-run は `no-write` を必須とする。 |
| `components` に `builder` を含む | `stdout` または `state`、`effects`。出力ファイル byte 比較または `state-diff.json` を必須とする。 |
| 外部 API / command / 通知を扱う | `effects`、`secret-mask`。呼び出し回数、順序、payload、mask 済み値を必須とする。 |

**§27 expected/effects.json schema 固定契約：**

`expected/effects.json` は次の schema に従う。未知 key は禁止する。

```json
{
  "external_calls": [
    {
      "type": "github_status",
      "order": 1,
      "method": "POST",
      "target": "/repos/{owner}/{repo}/statuses/{sha}",
      "payload": {},
      "result": "success"
    }
  ],
  "commands": [],
  "notifications": [],
  "downloads": [],
  "streams": [],
  "unchanged_paths": [],
  "deleted_paths": [],
  "created_paths": [],
  "write_order": [],
  "forbidden_writes": [],
  "forbidden_calls": []
}
```

| key | 型 | 必須 | 仕様 |
|-----|----|------|------|
| `external_calls` | array[object] | 必須 | GitHub、SMTP、webhook、SSH、remote build など process 外呼び出し。呼び出しなしは空配列。 |
| `commands` | array[object] | 必須 | pipeline、hook、systemd、archive、setup/update など local command 実行。実行なしは空配列。 |
| `notifications` | array[object] | 必須 | 通知送信、pending 化、retry 対象。通知なしは空配列。 |
| `downloads` | array[object] | 必須 | snapshot / artifact download の byte size、content type、中断有無。該当なしは空配列。 |
| `streams` | array[object] | 必須 | SSE / fetch stream の event、close、error。該当なしは空配列。 |
| `unchanged_paths` | array[string] | 必須 | 実行後に変更があってはならない状態ファイル、出力ファイル、log。 |
| `deleted_paths` | array[string] | 必須 | 実行後に削除される path。削除なしは空配列。 |
| `created_paths` | array[string] | 必須 | 実行後に新規作成される path。作成なしは空配列。 |
| `write_order` | array[string] | 必須 | 書き込み順。書き込みなしは空配列。複数状態更新 fixture では空配列禁止。 |
| `forbidden_writes` | array[string] | 必須 | 書き込み禁止 path。read-only、dry-run、validation failure fixture では対象状態ファイルを必ず列挙する。 |
| `forbidden_calls` | array[string] | 必須 | 呼び出し禁止の外部 API / command / notification。呼び出し禁止なしは空配列。 |

`external_calls[]`、`commands[]`、`notifications[]` は `order` を持つ。並列処理 fixture で完了順が非決定の場合でも、期待保存順は `write_order` に固定する。secret を含む payload は、fixture 内でも平文を保存せず `"***"` を使う。

**§27 fixture 不足時 未完了判定固定契約：**

| 不足 | 未完了理由 |
|------|------------|
| fixture 名が §27 fixture カタログに存在しない。 | カタログ外 fixture のため未完了。 |
| 必須 fixture が存在しない。 | 機能の正常 / 異常 / no-op / security / partial coverage 不足。 |
| `manifest.json.assertions` が §27 fixture assertion 選択固定契約を満たさない。 | 合否判定不足。 |
| `expected/effects.json` が存在しない、または必須 key が欠落する。 | 副作用検証不足。 |
| `unchanged_paths` または `forbidden_writes` が空で、fixture が failure / noop / read-only / dry-run / validation error のいずれかである。 | 無変更保証不足。 |
| 外部 API / command / notification を扱う fixture で `forbidden_calls` が空、かつ禁止対象なしの理由が `manifest.json.not_applicable` にない。 | 外部副作用境界不足。 |
| secret fixture で `expected/security.json` が存在しない。 | secret mask 検証不足。 |
| JSON Lines を扱う fixture で末尾改行、行順、破損行保持/除外条件を期待値に含めていない。 | log / audit / history 検証不足。 |
| idempotency fixture で 1 回目と 2 回目の期待差分を分離していない。 | 再実行検証不足。 |
| partial fixture で失敗地点より後の `forbidden_writes` / `forbidden_calls` を列挙していない。 | 部分失敗境界不足。 |

上表に 1 件でも該当する場合、対象機能は実装済みとして扱わない。fixture が多くても、期待副作用、禁止副作用、secret mask、保存順、再実行差分が明示されていなければ、バグ修正ゼロ化の検証を満たさない。

**§27 fixture 相互整合固定契約：**

fixture 内の `manifest.json`、`input/*`、`expected/*` は相互に矛盾してはならない。下表の不整合が 1 件でもある場合、その fixture は失敗とし、対象機能の実装 PR は未完了とする。

| 整合対象 | 固定条件 | 不整合時の扱い |
|----------|----------|----------------|
| `manifest.json.name` と directory 名 | directory 名は `manifest.json.name` と完全一致する。 | fixture 名不一致として失敗。 |
| `manifest.json.category` と fixture 名 prefix | `success-*` は `success`、`failure-*` は `failure`、`partial-*` は `partial`、`noop-*` は `noop`、`security-*` は `security` とする。 | 分類不一致として失敗。 |
| `manifest.json.section` と fixture カタログ | section と fixture 名の組み合わせは §27 fixture カタログ固定契約に存在する組み合わせだけ許可する。 | カタログ外として失敗。 |
| `manifest.json.owner_component` と `components` | `owner_component` は `components` に必ず含める。 | owner 責務不一致として失敗。 |
| `manifest.json.collaborator_components` と `components` | `collaborator_components` は `components` にすべて含め、`owner_component` を含めてはならない。 | collaborator 責務不一致として失敗。 |
| `manifest.json.components` と入力ファイル | `runner` / `builder` を含む場合は `input/cli.json`、`api` / `sdk` / `ui` を含む場合は `input/request.json` を置く。該当しない場合は `not_applicable` に理由を置く。 | 入力責務不一致として失敗。 |
| `manifest.json.assertions` と期待値ファイル | `response` は `expected/response.json`、`stdout` は `expected/stdout.txt`、`stderr` は `expected/stderr.txt`、`state` は `expected/state/`、`logs` は `expected/logs/`、`effects` は `expected/effects.json`、`secret-mask` は `expected/security.json` を要求する。 | 期待値不足として失敗。 |
| `expected/effects.json.write_order` と `expected/state/` | `write_order` に列挙された path は `expected/state/` または `expected/logs/` に期待値を持つ。 | 保存順だけの空検証として失敗。 |
| `expected/effects.json.unchanged_paths` と `created_paths` / `deleted_paths` | 同一 path を `unchanged_paths` と `created_paths` または `deleted_paths` に同時に含めない。 | 副作用矛盾として失敗。 |
| `expected/effects.json.forbidden_writes` と `write_order` | 同一 path を `forbidden_writes` と `write_order` に同時に含めない。 | 禁止書込矛盾として失敗。 |
| `expected/effects.json.forbidden_calls` と `external_calls` / `commands` / `notifications` | 禁止対象と同じ call target を実行期待に含めない。 | 禁止呼び出し矛盾として失敗。 |
| `expected/security.json` と期待値全体 | `expected/security.json.forbidden_values` に含まれる文字列は、fixture directory 内の `input/` を除く全 expected file に出現してはならない。 | secret mask 不足として失敗。 |
| `input/fakes.json` と `expected/effects.json` | fake response を消費する外部 call / command は `expected/effects.json` に同数、同順で記録する。未消費 fake がある場合は `manifest.json.not_applicable` に理由を置く。 | fake / effect 不一致として失敗。 |
| `missing_state` と `input/state/` | `manifest.json.missing_state` に列挙した path は `input/state/` に存在してはならない。 | 初期状態矛盾として失敗。 |

**§27 component 別検証責務固定契約：**

実装 PR は、対象 component ごとに下表の責務を満たす。複数 component を含む機能では、owner component と collaborator component を fixture manifest に分けて記録し、各 component の責務をすべて満たすまで完了扱いにしてはならない。

| component | owner 時の必須検証責務 | collaborator 時の必須検証責務 | 完了判定 | 禁止越境 |
|-----------|------------------------|-------------------------------|----------|----------|
| `builder` | CLI 入力、Markdown 入力、出力 site / HTML / REPORT、asset、cache、dependency、meta、終了コードを fixture で固定する。 | runner / API に渡す REPORT、meta、dependency manifest の key 名と nullable 条件を固定する。 | 出力 file の byte 比較または構造化 expected が存在し、既存出力保護と失敗時 no-write が検証済み。 | GitHub read、状態ファイル直接更新、通知送信。 |
| `runner` | CLI、設定、lock、SHA cache、build log、history、status、queue、external call、notification、終了コードを fixture で固定する。 | builder output、statefile schema、commitstatus payload を仕様どおり消費し、補完しない。 | 成功、失敗、skip、partial、dry-run の状態差分と write order が検証済み。 | API endpoint 追加、SDK method 追加、UI 操作追加。 |
| `api` | method/path/query/body/header、auth/scope/rate limit、response、状態 read/write、access/audit/config log を fixture で固定する。 | SDK / UI が補完しなくても扱える response schema、HTTP status、error body を返す。 | read-only は no-write、write API は保存順、validation failure は forbidden_writes が検証済み。 | runner CLI 処理、UI DOM 操作、未定義状態ファイル作成。 |
| `sdk` | method、args、query 生成、error 変換、token 破棄、binary / stream handling を fixture で固定する。 | UI が API 詳細を知らずに扱える戻り値と error をそのまま伝播する。 | API response の補完がなく、`401` / `403` / `429` / network error が区別される。 | API response 推測補完、未定義 endpoint 呼び出し、状態ファイル直接操作。 |
| `ui` | SDK method 呼び出し、DOM 表示、disabled/loading/error、secret field 消去、再取得順を fixture で固定する。 | SDK 戻り値だけを表示し、API / 状態ファイルの内部構造を再解釈しない。 | 直接 API 呼び出し、状態ファイル操作、secret DOM 残存がない。 | 直接 `fetch()`、状態ファイル操作、外部 command 実行。 |
| `statefile` | schema、atomic write、JSON Lines、lock、破損時処理、保存順、no-write / forbidden write を fixture で固定する。 | runner / API / archive の保存対象ごとに `write_order`、`unchanged_paths`、`forbidden_writes` を固定する。 | read-only、dry-run、validation failure、partial failure の副作用境界が検証済み。 | component 固有の業務判断、UI 表示判断、API response 補完。 |
| `archive` | log archive、snapshot、download、delete、rollback の保存 / 取得 / 削除境界を fixture で固定する。 | API download / rollback response と runner history への影響を statefile 契約に合わせる。 | 元 log 保持、archive 破損時挙動、安全でない entry 拒否、rollback 失敗時 no-write が検証済み。 | build 成否の反転、元 build log の改変、未検証 archive 展開。 |
| `commitstatus` | GitHub Commit Status payload、pending / final 送信順、失敗時非反転、secret mask を fixture で固定する。 | runner の build 結果を受け取り、build 成否を変更せず送信結果だけを返す。 | pending 失敗、final 失敗、commit SHA なし、無効時呼び出し 0 が検証済み。 | build 実行判断、dry-run での GitHub write、status 失敗による build 成否反転。 |
| `security` | token hash、session、TOTP、scope、audit、rate limit、secret mask、forbidden call/write を fixture で固定する。 | API / SDK / UI / runner の secret 表示、認証失敗、副作用境界を検証する。 | 認証失敗、権限拒否、rate limit、audit failure の副作用境界が検証済み。 | 業務処理代行、認可前状態更新、secret 平文保存。 |
| `setup` | binary 配置、service 更新、rollback、secret 既存値保持、stdout/stderr mask、終了コードを fixture で固定する。 | runner / API の初期状態と既存 secret を壊さないことを effects で固定する。 | 部分失敗時の復元対象と復元禁止対象が `expected/effects.json` に明記済み。 | runtime 機能追加、状態 schema 暗黙変更、外部依存追加。 |

component 責務を別 PR へ分割する場合でも、分割先 PR が満たすべき owner component、collaborator component、fixture 名、期待ファイル、禁止副作用を PR 本文に明記する。責務の所在が不明な場合は、その機能を実装完了扱いにしてはならない。

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

**§27 PR 別必須記録固定契約：**

各 §27 実装 PR は、本文に下表を記録する。記録がない PR は、コードと fixture が存在しても未完了とする。

| 記録項目 | 必須内容 |
|----------|----------|
| wave | 対象 wave、対象 §27.x、先行 wave 完了 commit または PR 番号。 |
| 実装対象 | 実装する機能名、owner component、collaborator component、変更ファイル、追加 fixture path。 |
| 実装対象外 | 同じ wave 内で今回実装しない §27.x、後続 wave、MCP、外部公開構成、未定義 endpoint / UI / 状態ファイル。 |
| fixture | §27 fixture カタログの fixture 名、manifest / effects / security の検証結果。 |
| acceptance | §27 実装 PR acceptance checklist の各項目の pass / fail / 未実行。 |
| 後続影響 | 後続 PR が利用してよい contract、利用してはならない未固定 contract。 |

**§27 実装 PR 最終受け入れゲート：**

§27 の機能実装 PR は、下表の全 gate を満たした場合だけ「完了」と判定する。1 件でも未達がある場合は「未完了」、仕様逸脱または secret 漏えいリスクがある場合は「差し戻し」とする。

| gate | 完了条件 | 未完了条件 | 差し戻し条件 |
|------|----------|------------|--------------|
| scope | 実装対象が §27 fixture カタログ固定契約に存在する機能だけである。 | 実装対象節の記載が PR 本文にない。 | 未定義 endpoint、未定義 UI、未定義状態ファイル、MCP、外部公開構成を追加している。 |
| fixture | 対象 §27.x の必須 fixture がすべて存在し、skip されていない。 | 必須 fixture が不足、または fixture 名が不一致。 | fixture が実装挙動に合わせて期待値を緩めている。 |
| manifest | 全 fixture の `manifest.json` が schema、assertion 選択、owner / collaborator component 責務を満たす。 | assertion、references、owner_component、collaborator_components、components、fake_clock のいずれかが不足。 | unknown key、実 secret、実環境 path、乱数依存を含む。 |
| expected | `expected/response.json`、`expected/state/`、`expected/logs/`、`expected/effects.json`、`expected/security.json` が assertion と一致する。 | assertion に対応する expected file が不足。 | expected と manifest / effects / security が矛盾する。 |
| side effect | `write_order`、`unchanged_paths`、`forbidden_writes`、`forbidden_calls` が対象機能の成功 / 失敗 / no-op / partial を説明できる。 | 禁止副作用または無変更保証が不足。 | 失敗時に未許可状態を書き換える、外部呼び出しを行う。 |
| secret | secret 平文が expected、logs、effects、UI DOM、stdout/stderr に存在しない。 | secret 検証対象が不足。 | token、password、TOTP secret、PAT、Authorization header が平文で残る。 |
| component | builder / runner / api / sdk / ui / statefile / archive / commitstatus / security / setup の該当責務が全て fixture に紐づく。 | owner / collaborator component の所在が不明。 | SDK / UI が API response を推測補完、または UI が直接 API / 状態ファイルを操作する。 |
| repeatability | fake clock、fake external response、固定 path により、同じ fixture が同じ結果を再現する。 | idempotency / no-op の 2 回目期待値が不足。 | 現在時刻、実ネットワーク、実 OS 差分に依存する。 |

**§27 実装 PR acceptance checklist：**

実装 PR 本文または検証ログには、下表の項目を記録する。記録がない項目は未検証として扱い、対象機能を完了扱いにしてはならない。

| 項目 | 記録内容 |
|------|----------|
| 対象仕様 | 実装した §27.x、関連 §22 / §23 / §24 / §25 / §26、owner component、collaborator component。 |
| 対象 fixture | 作成または更新した fixture 名一覧。fixture 名は §27 fixture カタログ固定契約と一致させる。 |
| 実行結果 | fixture ごとの pass / fail、実行コマンド、終了コード。 |
| 状態差分 | 作成、更新、削除、変更禁止の path。`expected/effects.json` と一致させる。 |
| 外部副作用 | GitHub、SSH、SMTP、webhook、systemd、hook、remote build、notification の呼び出し回数と順序。 |
| secret 検証 | 禁止文字列、mask 対象、平文が残らないことを確認した出力範囲。 |
| 部分失敗 | partial / failure fixture の失敗地点、完了済み副作用、禁止副作用。 |
| 再実行 | idempotency / no-op fixture の 1 回目と 2 回目の差分。 |
| 対象外確認 | 未定義 endpoint、未定義 UI、未定義状態ファイル、MCP、外部公開構成を追加していないこと。 |

**§27 差し戻し固定条件：**

次のいずれかに該当する PR は、fixture が pass していても差し戻しとする。

| 条件 | 理由 |
|------|------|
| 仕様にない endpoint、SDK method、UI 操作、状態ファイルを追加している。 | 仕様外実装。 |
| fixture の期待値が実装都合に合わせて仕様より弱い。 | 検証の形骸化。 |
| secret 平文、Authorization header、TOTP secret、PAT、password が expected または log に残る。 | secret 漏えい。 |
| read-only / dry-run / validation failure で状態、log、外部 call が変化する。 | 副作用違反。 |
| partial failure で失敗地点以降の write / call が発生する。 | 部分失敗境界違反。 |
| SDK が API response を補完し、UI が SDK を迂回し、runner / builder が未定義状態ファイルを作成する。 | component 責務違反。 |
| 実ネットワーク、実時刻、実ユーザー環境、実 secret に依存する fixture だけで合格している。 | 再現性不足。 |

**§27 部分失敗・再実行固定契約：**

| ケース | 固定挙動 |
|--------|----------|
| 状態保存前の validation 失敗 | 状態ファイル、JSON Lines、外部 API、外部 command、通知を実行しない。 |
| 主状態保存後の log 追記失敗 | 個別節が rollback を明記しない限り主状態は戻さず、response は `500` とし、次回 GET は保存済み主状態を返す。 |
| log 追記後の audit 失敗 | 個別節で必須 audit と定義された操作は `500`。任意 audit の場合は対象操作成功を維持し、WARN を保存する。 |
| 外部 API 送信成功後の状態保存失敗 | 外部送信の再送を自動実行しない。状態保存失敗を `500` または runner failure として記録し、再実行時は個別節の重複防止 key で判定する。 |
| 通知送信失敗 | build / config / security の主結果を反転しない。`.notify_pending` または個別節の失敗記録だけを更新する。 |
| download / stream 中断 | サーバー側状態を成功/失敗へ変更しない。access log は中断 status を記録してよいが、history と build log は変更しない。 |
| 再実行 no-op | 同一入力で差分がない保存 API は、個別節の no-op response を返し、状態、config log、audit log、notify log に新規差分を作らない。 |
| 破損 JSON Lines | read API は破損行を除外し、破損内容を response に出さない。write API は既存破損行を修復、削除、並べ替えしない。 |

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

runner と `POST /api/logs/archive` は、`.server_config.log_archive_after_days > 0` の場合、対象日数より古い `.build_logs/{id}.json` を gzip 圧縮し、`.build_logs/archive/{id}.json.gz` へ保存する。圧縮成功後、元の `.build_logs/{id}.json` を削除する。`.build_logs/archive/` 内のファイルを再圧縮してはならない。

gzip は Go 標準ライブラリ `compress/gzip` を使用し、mtime は元ファイル mtime ではなく圧縮実行時刻でよい。圧縮前 JSON を読み込めないファイルは archive 対象外とし、WARN `LOG_ARCHIVE_SKIP_CORRUPT: id=<id>` を出す。実行中 build の `current_build_id` と一致する log は対象外とする。

ログ参照 API は通常ファイルを先に探し、存在しない場合に archive を探す。archive を読む場合は gzip 展開後に通常 `.build_logs/{id}.json` と同じ schema として扱う。`GET /api/logs/search`、`GET /api/history/{id}/log`、`GET /api/output-meta`、`GET /api/disk-usage` は archive を含めて動作する。

`POST /api/logs/cleanup` は、archive 済みファイルも `log_retention_days` の削除対象に含める。`POST /api/logs/archive` は `{ "message": "Logs archived", "archived_count": N }` を返す。

**archive / cleanup 固定契約：**

| 項目 | 仕様 |
|------|------|
| archive 対象判定 | build log JSON の `finished_at` を基準にする。欠落時は file mtime を使わず対象外。 |
| archive id | file 名 `{id}.json` の id と JSON 内 `id` が一致する場合だけ対象。 |
| gzip path | `.build_logs/archive/{id}.json.gz`。既存 archive がある場合は上書きせず skip する。 |
| cleanup 順 | 通常 log 削除 → archive log 削除 → 空 archive directory 削除試行。 |
| 削除失敗 | 処理継続し、response に `failed_count` を含める。 |
| response | archive は `archived_count`、cleanup は `deleted_count` と `failed_count` を返す。 |

検証条件:

| ケース | 期待結果 |
|--------|----------|
| 対象ログあり | `.json.gz` 作成、元 `.json` 削除、API 参照可。 |
| 実行中ログ | archive しない。 |
| 破損ログ | archive しない、WARN、処理継続。 |
| archive API | 件数を返し、disk usage に archive bytes を含める。 |
| cleanup | 通常 log と archive log の両方を保持期間で削除する。 |
| archive 既存 | 上書きせず skip。 |
| cleanup 一部失敗 | `failed_count` に計上し処理継続。 |

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

本機能の目的は、`.snapshots/` に保存された build artifact を API、SDK、UI から一覧、download、削除、rollback できるようにすることである。

owner component は `archive` とする。collaborator component は `api`、`sdk`、`ui`、`runner`、`statefile` とする。snapshot 作成は `runner` の §14b を正とする。

**API 契約：**

| API | 処理 |
|-----|------|
| `GET /api/snapshots` | `.snapshots/{id}/` を新しい順で一覧する。 |
| `GET /api/snapshots/{id}/download` | 対象 snapshot を tar.gz として streaming download する。 |
| `DELETE /api/snapshots/{id}` | 対象 snapshot だけを削除し、`.config_log` に記録する。 |
| `POST /api/history/{id}/rollback` | 対象 snapshot を deploy target へ再転送し、新規 rollback build log/history を作成する。 |

`id` は build id と一致するものだけ許可する。`/`、`..`、空文字、URL decode 後に path separator を含む値は `422` とする。

**download tar.gz 生成契約：**

| 項目 | 仕様 |
|------|------|
| root | `.snapshots/{id}/` を root とし、root 外を参照しない。 |
| entry path | root からの相対 path。`/` 始まり、`..`、空 segment、NUL、Windows drive prefix は禁止。 |
| entry 種別 | 通常ファイルと directory だけを含める。symlink、hardlink、device、socket、fifo は含めない。 |
| header | `Content-Type: application/octet-stream`、`Content-Disposition: attachment; filename="{id}.tar.gz"`。 |
| 順序 | directory、file とも相対 path 辞書順。 |
| mtime | snapshot 内 file の mtime を使用してよい。存在しない場合は build log の `finished_at`。 |
| secret 除外 | `.github_token`、`.admin_credentials`、`.api_tokens`、`.smtp_secret`、`.webhook_secret`、runner 状態ファイル名は検出時点で `500` とし、download を中止する。 |

**Rollback 仕様：**

rollback は新しい build id を採番し、`.build_history.trigger="rollback"`、`rollback_from=<元id>` を保存する。元 snapshot は変更しない。rollback 中に別 build が running の場合は `409` とする。転送失敗時は rollback build log を `failure` とし、元 snapshot は削除しない。

rollback は snapshot 内の成果物を deploy target へ再転送する操作であり、以下を行ってはならない。

| 禁止対象 | 理由 |
|----------|------|
| `.last_sha` 更新 | rollback は監視対象 SHA の処理完了ではない。 |
| `.server_config`、`.branch_config`、`.notify_config` の復元 | 設定 rollback ではない。 |
| `.build_history` の過去行書き換え | rollback は新規履歴として追記する。 |
| `.build_logs/{元id}.json` の変更 | 元 build の証跡を保持する。 |
| 元 snapshot の削除または上書き | rollback 成否に関係なく元成果物を保持する。 |
| secret / token / credentials の復元 | snapshot に secret を含めないため復元対象外。 |

rollback build log は `target_status="success"` または `failure_build` とし、`trigger="rollback"`、`rollback_from=<元id>`、`snapshot_id=<元id>` を含める。rollback 転送で pending が発生した場合は `success_deploy_pending` とし、`.pending_transfers` に rollback 用 entry を追加する。

**snapshot 一覧・削除固定契約：**

| 項目 | 仕様 |
|------|------|
| 一覧対象 | `.snapshots/{id}/manifest.json` が存在する directory だけ。 |
| size | directory 配下の通常ファイル size 合計。symlink は size 集計前に異常扱い。 |
| delete 順 | id validation → running check → snapshot directory 確認 → delete → `.config_log` 追記 → response。 |
| delete log 失敗 | snapshot 削除済みのまま `500`。削除は巻き戻さない。 |
| rollback pending | pending entry には `rollback_from`、`snapshot_id`、deploy target を保存する。 |

**UI / SDK：**

SDK は `getSnapshots()`、`downloadSnapshot(id)`、`deleteSnapshot(id)`、`rollbackHistory(id)` を提供する。UI は snapshot 一覧に id、saved_at、size_bytes、download、delete、rollback 操作を表示する。delete と rollback は実行中 build がある場合 disabled とする。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 一覧 | snapshot id、build id、保存日時、size が返る。 |
| download | tar.gz を返し、snapshot 外のファイルを含まない。 |
| delete | 対象 id だけ削除、config log 追記。 |
| rollback 成功 | 新規 build id、trigger rollback、rollback_from 保存。 |
| 不正 id | `422`、状態差分なし。 |
| download symlink | symlink entry を含めず、secret 名検出時は `500`。 |
| rollback pending | 新規 rollback log/history、pending entry、元 snapshot 維持。 |
| rollback running | `409`、状態差分なし。 |
| delete log failure | snapshot は削除済み、response は `500`。 |

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
