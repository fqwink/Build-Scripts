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
| §21〜§22 | `components/api.go` / `adlaire-ci-api` の詳細仕様 |
| §23 | `admin/adlaire-ci-sdk.js` の詳細仕様 |
| §24 | `admin/index.html` の詳細仕様 |
| §25 | 認証の実装仕様 |
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
| `components/api.go` | §21〜§22、§25、§26 | API 共通処理、endpoint、状態ファイル、認証、session、systemd。 |
| `admin/adlaire-ci-sdk.js` | §23 | SDK class、method、HTTP 対応、error、stream、token 破棄。 |
| `admin/index.html` | §24 | 画面構成、DOM id、panel、SDK 呼び出し、表示状態、秘密情報消去。 |
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
| `components/api.go` | §21〜§22、§25、§26 に API 共通契約、endpoint、状態ファイル schema、認証、認可、systemd、セットアップが定義されている。 | §22.0e にない endpoint、method、status code、response body、状態ファイル write を追加すること。 | §22.0e の全 endpoint が Request、Response、Errors、Read、Write、SDK、UI の対応表と一致する。 |
| `admin/adlaire-ci-sdk.js` | §23 に SDK class、method、引数、戻り値、HTTP endpoint 対応、error object、token 破棄条件が定義されている。 | SDK が §22.0e にない endpoint を呼ぶこと、body 禁止 endpoint に body を送ること、独自 error 形式を返すこと。 | 全 method が §22.0e と §23 の対応どおりに動作し、HTTP error を `AdlaireCIError` として扱う。 |
| `admin/index.html` | §24 に画面構成、panel、操作、成功表示、失敗表示、disabled、再取得、秘密情報消去が定義されている。 | SDK を介さず API を直接呼ぶこと、未定義の画面・操作・保存先を追加すること、秘密情報を DOM に残すこと。 | 全 UI 操作が §24 の表示条件と §23 の SDK method を満たし、秘密情報 field が指定条件で消去される。 |

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
| Phase 3 | `components/api.go` P0 / P1 | §21〜§22、§25、§26 のうち、認証、セッション、共通エラー、状態ファイル読み書き、ビルド操作、status、logs、history、queue、circuit breaker。 | Phase 2 が完了し、runner が書き込む状態ファイル schema が固定されている。 | §22.0f P0 / P1 の必須検証、§0e の `components/api.go` API 共通・状態ファイル検証、§0f の `components/api.go` 完了判定の該当範囲を満たす。 |
| Phase 4 | `components/api.go` P2〜P5 | §22.0f P2〜P5 の config、repo、branch、schedule、notify、snapshot、rollback、maintenance、access control、hooks、tokens 等。 | Phase 3 が完了し、API 共通処理と認証が固定されている。 | §22.0f P2〜P5 の必須検証と §0e の `components/api.go` endpoint 契約を満たす。 |
| Phase 5 | `admin/adlaire-ci-sdk.js` | §23 の SDK class、method、戻り値、HTTP error、token 破棄、query 生成。 | Phase 3 と Phase 4 が完了し、§22.0e の endpoint 契約が固定されている。 | §0e の SDK 契約と §0f の `admin/adlaire-ci-sdk.js` 完了判定を満たす。 |
| Phase 6 | `admin/index.html` | §24 の標準管理ツール UI、panel、操作、成功表示、失敗表示、disabled、再取得、秘密情報消去。 | Phase 5 が完了し、SDK method 契約が固定されている。 | §0e の UI 契約と §0f の `admin/index.html` 完了判定を満たす。 |

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
| 実装開始条件 | Phase 3 / Phase 4 が完了し、§22.0e の全 endpoint 契約が固定済みである。 |
| 入力 | constructor の `baseUrl`、method 引数、browser 標準 API、session token。 |
| 出力 | `Promise`、`AdlaireCIError`、`StreamHandle`、JSON object、SSE callback / stream reader、token 破棄。 |
| 正常系 | §23 の全 method が §22.0e の endpoint のみを呼び、query、body、HTTP method、response 変換、stream close を仕様どおり処理する。 |
| 異常系 | HTTP error、network error、timeout、unsupported browser、body 禁止 endpoint、認証失敗、`401` token 破棄を `AdlaireCIError` 契約に一致させる。 |
| セキュリティ | token を localStorage / sessionStorage へ保存しない。secret 入力値を console に出力しない。global 汚染を行わない。 |
| 実装対象外 | DOM 操作、admin panel 表示、API endpoint 新設、Node.js 専用 API、runtime 名別の互換分岐、bundler、polyfill、npm package。 |
| 必須検証 | §23 の method 契約、§0e の SDK 契約、§26.7 SDK 対象。API は fake fetch で成功 / error / timeout / stream を再現する。 |
| 完了条件 | 全 method が §22.0e と §23 の対応どおり動作し、HTTP error を `AdlaireCIError` として扱い、body 禁止 endpoint に body を送らない。 |

Phase 5 完了時は、Phase 6 へ引き継ぐ UI 契約として、SDK method 名、引数、戻り値、error object、loading / retry に必要な失敗情報、stream handle を固定する。

#### 0g.5.1 Phase 5 実装順序

Phase 5 は API 契約の薄い wrapper とし、表示判断を持たせない。実装は fake fetch による endpoint 対応確認から進める。

| 順序 | 実装単位 | 完了判定 |
|------|----------|----------|
| 1 | `AdlaireCI` constructor、base URL 正規化、token 保持、共通 request、query 生成を実装する。 | 末尾 slash、query encoding、Authorization header、body 禁止 endpoint が §23 と一致する。 |
| 2 | `AdlaireCIError`、HTTP error 変換、network error、timeout、unsupported browser 判定を実装する。 | UI が `status`、`message`、`details`、`requestId` を追加判断なしに表示できる。 |
| 3 | §23 の non-stream method を endpoint 対応表どおり実装する。 | 各 method の HTTP method、path、query、body、戻り値が §22.0e と一致する。 |
| 4 | `streamBuild()` と `StreamHandle.close()` を実装する。 | `fetch()` stream、AbortController、SSE frame parse、invalid frame、end event が §23 と一致する。 |
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
| 正常系 | §24 の panel、DOM id、表示条件、操作、成功表示、再取得、stream 表示、logout を仕様どおり処理する。 |
| 異常系 | login 失敗、session 期限切れ、API error、validation error、stream 切断、network error、maintenance、access block を UI 表示契約に一致させる。 |
| セキュリティ | secret 値を DOM に残さない。password / PAT / token / SMTP password / webhook secret は成功・失敗に関わらず指定条件で消去する。SDK を介さない直接 API 呼び出しは禁止。 |
| 実装対象外 | frontend framework、CSS framework、chart library、CDN script、build tool、外部 icon package、API endpoint 新設。 |
| 必須検証 | §24 の全 panel と主要操作、§0e の UI 契約、§26.7 UI/security 対象。SDK は fake implementation で成功 / 失敗 / loading / stream を再現する。 |
| 完了条件 | 全 UI 操作が §24 の表示条件と §23 の SDK method を満たし、成功表示、失敗表示、disabled、再取得、secret 消去が一致する。 |

Phase 6 完了時は、初期実装全体の完了判定として、Phase 1〜6 の引き継ぎ契約、§0e、§0f、§0g、§0i、§0j、§22.0f、§23、§24、§26.7 を再確認する。

#### 0g.6.1 Phase 6 実装順序

Phase 6 は、SDK 契約の利用者として UI を実装する。API 仕様の不足を UI 側で補完してはならない。

| 順序 | 実装単位 | 完了判定 |
|------|----------|----------|
| 1 | HTML shell、panel、DOM id、navigation、初期 disabled / loading 状態を実装する。 | §24 の DOM id と panel 構成が一致し、未認証時に保護操作が disabled になる。 |
| 2 | login、logout、session expiry、共通 error 表示、再 login 導線を実装する。 | `AdlaireCIError` の message / details が同一 panel 内に表示され、secret が残らない。 |
| 3 | status、history、logs、queue、manual build、force build、cancel、stream 表示を実装する。 | SDK method だけを呼び、loading、disabled、再取得、stream close が §24 と一致する。 |
| 4 | config、repo、branch、schedule、notify、snapshot、rollback、maintenance、access、hooks、tokens の panel 操作を実装する。 | 成功表示、失敗表示、secret field 消去、保存後再取得が各 panel の仕様と一致する。 |
| 5 | dashboard layout、notes、alert、tag rule、pipeline config の表示 / 保存を実装する。 | UI が未定義 field を追加せず、SDK response に基づく表示だけを行う。 |
| 6 | fake SDK で成功、失敗、loading、stream、session expiry、secret 消去を検証する。 | 初期実装全体の UI 受け入れ条件を満たし、直接 API 呼び出しが存在しない。 |

#### 0g.6.2 Phase 6 固定仕様

| 対象 | 固定仕様 |
|----------|----------|
| API 通信 | UI は必ず SDK method を呼ぶ。`fetch()`、`XMLHttpRequest`、直接 `ReadableStream` 生成は禁止する。 |
| framework | frontend framework、CSS framework、chart library、CDN script、外部 icon package、build tool は使用しない。 |
| secret field | password、PAT、token、SMTP password、Webhook Secret は成功 / 失敗に関わらず §24 の条件で消去する。 |
| unknown state | SDK response にない状態を UI が推測して成功扱いにしてはならない。不明状態は error または未取得として表示する。 |
| disabled | 実行中、未認証、maintenance、権限不足、入力不正時の disabled 条件は §24 に従う。 |
| UI 追加 | §24 にない panel、操作、保存先、設定項目を UI 都合で追加してはならない。 |

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

表の「詳細仕様節」が複数ある場合は、すべての節を同時に満たす。`components/builder.go` が責務 component に含まれる機能では、`ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` の同番号節を合わせて確認する。`components/runner.go` が責務 component に含まれる機能では、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` の同番号節を合わせて確認する。該当節に §0h の必須項目が不足している場合は、その項目を実装せず、先に詳細仕様を改訂する。

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

## 21. 管理ツール システム構成

```
systemd timer
  └─ components/runner.go（変更検出・ビルド起動）
       └─ SSH 転送

components/api.go（常駐 HTTP サーバー）

admin/index.html（標準管理ツール）
  └─ adlaire-ci-sdk.js（SDK）─── HTTP ───► components/api.go
```

owner component `api` は、Go 標準ライブラリ `net/http` で実装し、管理ツールからの API リクエストを受け付ける。`runner` とは独立して常駐する。

**`components/api.go` 設定値（スクリプト冒頭）：**

```go
Host              = "127.0.0.1"                               // バインドアドレス（外部公開禁止）
Port              = 8765                                      // リッスンポート
CredentialsFile   = "/opt/adlaire-builder/.admin_credentials" // 認証情報ファイル
OutputURL         = "https://example.com/" // 出力サイトの公開 URL
HistoryFile       = "/opt/adlaire-builder/.build_history"     // ビルド履歴ファイル
NotifyConfigFile  = "/opt/adlaire-builder/.notify_config"     // Webhook 通知設定
ServerConfigFile  = "/opt/adlaire-builder/.server_config"     // サーバー設定
AccessLogFile     = "/opt/adlaire-builder/.access_log"        // ログイン履歴
NotifyLogFile     = "/opt/adlaire-builder/.notify_log"        // Webhook 送信履歴
WebhookSecretFile = "/opt/adlaire-builder/.webhook_secret"    // GitHub Webhook HMAC-SHA256 Secret（→ §22）
SnapshotDir       = "/opt/adlaire-builder/.snapshots"         // スナップショット保存ディレクトリ（→ §14b）
LogLevel          = "INFO"
Owner             = "<GitHubオーナー名>"                       // 初期値。POST /api/repo-config で動的変更可能（.repo_config に保存）
Repo              = "<リポジトリ名>"                           // 初期値。POST /api/repo-config で動的変更可能（.repo_config に保存）
```

**systemd ユニット（常駐型、タイマー不要）：**

```ini
[Unit]
Description=Adlaire CI API Server
After=network.target

[Service]
Type=simple
User=deploy
ExecStart=/usr/local/bin/adlaire-ci-api
Restart=on-failure
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable --now adlaire-ci-api  # 登録・起動
sudo journalctl -u adlaire-ci-api -f        # ログ確認
```

---

## 22. バックエンド API 仕様

**ベース URL：** `http://localhost:{PORT}/api`
**認証：** `Authorization: Bearer {SESSION_TOKEN}`（`/api/login` で取得したセッショントークン）
**レスポンス形式：** JSON

### 22.0 API 共通契約

本節の API は `components/api.go` の対象仕様である。実装時は、エンドポイント固有仕様より先に以下の共通契約を満たす。

| 項目 | 仕様 |
|------|------|
| Go バージョン | Go `1.22` 以上。HTTP 実装は Go 標準ライブラリ `net/http` を使用する。 |
| bind | 既定値は `127.0.0.1:8765`。`--addr` で上書き可能。`--addr 0.0.0.0:<port>` を指定しても、`components/api.go` は TLS listener、origin 制限、IP allowlist、reverse proxy 設定生成を追加実行しない。 |
| 文字コード | リクエストボディ、レスポンスボディ、状態ファイルはいずれも UTF-8 とする。 |
| JSON レスポンス | JSON レスポンスには `Content-Type: application/json; charset=utf-8` を付与する。 |
| リクエスト body 上限 | JSON body は 1 MiB を上限とする。超過時は `413 Payload Too Large` と `{"error": "Payload too large"}` を返す。 |
| request body 禁止 | §22.0e で `Request` が `none` の endpoint に body がある場合は `400 Bad Request` と `{"error": "Request body is not allowed"}` を返す。 |
| 成功レスポンス | 各エンドポイント例に記載した JSON オブジェクトを返す。空レスポンスは使用しない。 |
| エラーレスポンス | エラー時は `{"error":"<message>"}` を返す。入力検証失敗時のみ `details` を配列 `[{ "field": "<field>", "message": "<reason>" }]` とし、複数エラーがある場合はリクエスト JSON の出現順、query、path parameter の順で並べる。入力検証以外の補足は `details` を使わず、`error` を実装者向けではない固定文言にする。 |
| 未知のパス | 定義されていない `/api/...` は `404 Not Found` と `{"error": "Not found"}` を返す。 |
| 未対応メソッド | パスは存在するがメソッドが異なる場合は `405 Method Not Allowed` と `{"error": "Method not allowed"}` を返す。 |
| JSON 不正 | JSON ボディのパースに失敗した場合は `400 Bad Request` と `{"error": "Invalid JSON"}` を返す。 |
| 入力検証失敗 | 型、必須キー、範囲、有効値が仕様と異なる場合は `422 Unprocessable Entity` と `{"error":"Validation failed","details":[...]}` を返す。`field` は JSON body key、query key、または path parameter 名とし、body 全体の形式不正は `field` を `"$"` とする。 |
| 認証なし | 認証必須エンドポイントで Bearer トークンがない、または無効な場合は `401 Unauthorized` と `{"error": "Unauthorized"}` を返す。 |
| 権限不足 | 認証済み API token の scope が不足する場合は `403 Forbidden` と `{"error":"Forbidden"}` を返す。管理 session は全 API 操作を許可する。API token は `read`、`trigger`、`operate`、`config`、`admin` の scope だけを許可し、token 作成時に指定された scope 外の endpoint は拒否する。 |
| 競合 | 現在状態と要求操作が両立しない場合は `409 Conflict` を返す。対象は、実行中ビルドへの二重開始、ビルド未実行時の cancel、停止済みスケジュールへの pause、稼働中スケジュールへの resume、lock 取得 10 秒超過、stale 判定不能な `.build_lock` である。 |
| 未設定機能 | endpoint の必須 secret、必須外部設定、必須状態ファイルが未設定で処理を開始できない場合は `501 Not Implemented` と `{"error":"Not configured"}` を返す。エンドポイント固有仕様で `422`、`503`、`500` を明記している場合のみ個別指定を優先する。 |
| 時刻形式 | API レスポンスと状態ファイルの機械処理用時刻は UTC ISO 8601 `YYYY-MM-DDTHH:MM:SSZ` とする。明示オフセット、timezone なし文字列、ミリ秒付き文字列は保存しない。外部 API から取得した時刻も保存前に UTC `Z` へ正規化する。 |
| GET の副作用 | `GET` エンドポイントは状態ファイルを書き換えない。診断 API が外部確認を行う場合も、結果保存は行わない。 |
| 状態ファイル更新 | JSON 状態ファイルの更新は同一ディレクトリの一時ファイルへ書き出してから `os.Rename` で置換する。秘密情報を含むファイルは作成後に mode `600` を設定する。rename 後は対象ファイルと親ディレクトリを `Sync` し、永続化失敗時は `500` を返す。 |
| 秘密情報 | PAT、Webhook Secret、セッショントークン、API トークンはログ、バックアップ、GET レスポンスへ平文出力しない。設定済み表示は `"***"` または boolean で返す。 |
| 並列更新 | 同一状態ファイルを更新する API は、ファイル単位のロックを取得してから読み込み、検証、書き込みを行う。ロック取得待ちは最大 10 秒とし、超過時は `409 Conflict` を返す。 |
| 監査ログ | 設定変更 API は、変更前後の値を `.config_log` に追記する。ただし秘密情報の値は変更前後とも `"***"` にマスクする。 |
| CORS | 既定では CORS ヘッダーを付与しない。標準管理ツールは同一 origin から配信する。`OPTIONS` preflight は定義しない。CORS を有効化する拡張は本ファイルで未定義とし、実装してはならない。 |
| セキュリティヘッダー | すべての API レスポンスに `Cache-Control: no-store`、`X-Content-Type-Options: nosniff` を付与する。SSE は `Cache-Control: no-store` と `X-Accel-Buffering: no` を付与する。 |
| 判定順 | path 解決 → method 検証 → body 可否/サイズ検証 → JSON parse → 認証 → 権限 → 入力検証 → 状態競合 → 処理実行の順に判定する。 |

エンドポイント例に記載されたフィールド名、型、有効値、HTTP ステータスは規範とする。API、SDK、標準管理ツールのいずれかを変更する場合は、§22、§23、§24 の対応関係を同時に確認する。

**API 共通エラー固定文言：**

| 条件 | HTTP status | body |
|------|-------------|------|
| 未知 path | `404` | `{"error":"Not found"}` |
| method 不一致 | `405` | `{"error":"Method not allowed"}` |
| body 禁止 endpoint に body あり | `400` | `{"error":"Request body is not allowed"}` |
| body 上限超過 | `413` | `{"error":"Payload too large"}` |
| JSON parse 失敗 | `400` | `{"error":"Invalid JSON"}` |
| 認証なし / 無効 token / 期限切れ session | `401` | `{"error":"Unauthorized"}` |
| scope 不足 / 管理操作不可 | `403` | `{"error":"Forbidden"}` |
| rate limit 超過 | `429` | `{"error":"Too many requests"}` |
| 入力検証失敗 | `422` | `{"error":"Validation failed","details":[...]}` |
| 状態競合 | `409` | endpoint 固有文言。未定義の場合は `{"error":"Conflict"}` |
| 必須設定なし | `501` | `{"error":"Not configured"}` |
| 内部処理失敗 | `500` | `{"error":"Internal server error"}` |

上表の body は空白差分を除いて固定とする。`500` の response body に Go error、path、secret、状態ファイル内容、外部 API response body を含めてはならない。内部原因は server log にだけ固定コード付きで出力する。

**API 実行順・副作用境界固定契約：**

全 endpoint は、下表の順序と副作用境界に従う。endpoint 固有節で異なる順序を明記していない限り、実装者判断で検証順、状態読取、状態書込、外部呼び出し、ログ追記の順序を入れ替えてはならない。

| 段階 | 処理 | 失敗時 status | 失敗時副作用 |
|------|------|---------------|--------------|
| 1 | path 解決。未知 path を判定する。 | `404` | 状態ファイル、外部 API、監査ログ、設定ログを変更しない。`.api_access_log` だけ §27.6 に従って記録する。 |
| 2 | method 検証。path が存在し method が不一致か判定する。 | `405` | endpoint 固有処理を開始しない。`.api_access_log` 以外を変更しない。 |
| 3 | body 禁止、body size、JSON parse を検証する。`Request=none` endpoint では JSON parse を行わない。 | `400` / `413` | 認証、状態読取、状態書込、外部呼び出しを行わない。 |
| 4 | access control、maintenance、認証不要 endpoint 判定、Bearer/session/API token 認証を実行する。 | `401` / `403` / `503` | 認証失敗時は endpoint 状態を読まない。認証処理で定義された session/ticket 更新とログ追記だけを許可する。 |
| 5 | scope、rate limit、権限を判定する。 | `403` / `429` | endpoint 固有の状態書込、外部呼び出し、コマンド実行を行わない。rate limit 状態更新と監査ログは §27.47 に限定する。 |
| 6 | query、path parameter、body schema、enum、範囲を検証する。 | `422` | 状態ファイルを変更しない。外部 API、systemd、runner、hook、通知を呼び出さない。 |
| 7 | endpoint 固有の read adapter を呼び、状態競合を判定する。 | `409` / `500` | write lock を取得していても target を変更しない。tmp file があれば削除する。 |
| 8 | endpoint 固有処理を実行し、必要な状態ファイルを §22.0d の Write 列順に更新する。 | endpoint 固有 | 途中失敗時の巻き戻しは、endpoint 固有節または §26 rollback 節に明記された範囲だけ行う。 |
| 9 | `.config_log`、`.audit_log`、`.access_log`、`.api_access_log` を仕様順に追記する。 | `500` | response を成功扱いにしない。既に確定済みの endpoint 状態は自動推測で再変更しない。 |
| 10 | response body と header を確定する。 | - | response 生成時に追加の状態読取、状態書込、外部呼び出しを行わない。 |

`GET` endpoint は段階 8 で状態ファイルを書き換えない。`POST`、`DELETE` endpoint でも、段階 6 までに失敗した場合は endpoint 固有の状態書込を一切行わない。外部 API 送信、systemd 操作、hook 実行、通知送信、snapshot 操作は、endpoint 固有節の処理順に現れる場合だけ実行する。実装者判断で「先に外部確認してから validation error を返す」処理にしてはならない。

### 22.0a 状態ファイル共通仕様

`components/api.go` および拡張後 `components/runner.go` が読み書きする状態ファイルは、下表の初期値、形式、更新責務に従う。表にない状態ファイルを追加してはならない。追加が必要な場合は、先に本節へパス、形式、初期値、更新責務、破損時の扱いを追記する。

| パス | 形式 | 初期値 | 更新責務 | 破損時の扱い |
|------|------|--------|----------|--------------|
| `.admin_credentials` | JSON object | `--init-credentials` で生成 | `components/api.go` | 起動時に ERROR ログを出し、HTTP サーバーを起動しない。 |
| `.totp_secret` | JSON object | `{"enabled":false,"secret_base32":null,"confirmed_at":null,"last_accepted_step":null}` | `components/api.go` | 読み込み不能時は TOTP 有効 login を `500` で拒否する。破損時は退避するが自動再生成で認証を弱めてはならない。 |
| `.audit_log` | JSON Lines | 空ファイル | `components/api.go` | 読み込み可能な行のみ返し、壊れた行は無視する。追記不能時は対象操作を失敗扱いにする。 |
| `.api_rate_state` | JSON object | `{"windows":{}}` | `components/api.go` | `.api_rate_state.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、空 window で再生成する。 |
| `.server_config` | JSON object | `{}` | `components/api.go` | `.server_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、空 object で再生成する。 |
| `.notify_config` | JSON object | `{"webhooks":[],"channels":[],"on":[],"summary":{"enabled":false,"interval":"weekly","hour":9,"day_of_week":1},"email":{"enabled":false,"to":[],"on":[]}}` | `components/runner.go` / `components/api.go` | `.notify_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、初期値で再生成する。 |
| `.notify_log` | JSON Lines | 空ファイル | `components/runner.go` | 読み込み可能な行のみ使用し、壊れた行は ERROR ログへ記録して無視する。 |
| `.notify_pending` | JSON array | `[]` | `components/runner.go` | `.notify_pending.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、`[]` で再生成する。 |
| `.pending_transfers` | JSON array | `[]` | `components/runner.go` | `.pending_transfers.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、`[]` で再生成する。 |
| `.build_history` | JSON Lines | 空ファイル | `components/runner.go` | 読み込み可能な行のみ使用し、壊れた行は ERROR ログへ記録して無視する。 |
| `.build_logs/{id}.json` | JSON object | ビルドごとに新規作成 | `components/runner.go` | 対象 ID の API は `500` を返し、既存ファイルは上書きしない。 |
| `.build_lock` | text | 不在 | `components/runner.go` | 内容は `pid={pid}\nstarted_at={UTC_ISO8601}\n` とする。PID が存在しない場合は stale lock として削除し、存在する場合は `409` 相当の実行中として扱う。形式不正または PID 判定不能の場合は上書きせず `409` を返す。 |
| `.branch_config` | JSON object | 不在 | `components/runner.go` / `components/api.go` | `.branch_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、再生成せず `BRANCH_TARGETS` デフォルトへフォールバックする。 |
| `.build_state` | JSON object | `{"running":false,"current_build_id":null,"queued":[],"last_started_at":null,"last_finished_at":null,"weekly_summary_last_sent_at":null,"weekly_summary_sent_date":null}` | `components/runner.go` / `components/api.go` | `.build_state.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、初期値で再生成する。 |
| `.build_status.json` | JSON object | `{"schema_version":1,"updated_at":null,"status":"none","running":false,"current_build_id":null,"last_build_id":null,"last_trigger":null,"last_target_status":null,"last_branch":null,"last_target_file":null,"last_blob_sha":null,"last_commit_sha":null,"last_started_at":null,"last_finished_at":null,"last_duration_seconds":null,"last_error":null,"last_deploy_status":null,"pending_transfers_count":0,"notify_pending_count":0,"circuit_open":false,"circuit_consecutive_failures":0,"output_sha256":null,"size_warn":false}` | `components/runner.go` | `.build_status.json.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、初期値で再生成する。 |
| `.build_circuit_state` | JSON object | `{"open":false,"consecutive_failures":0,"opened_at":null,"last_failure_at":null,"last_error":null}` | `components/runner.go` / `components/api.go` | `.build_circuit_state.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、初期値で再生成する。 |
| `.local_watch_state.json` | JSON object | `{"files":{}}` | `components/runner.go` | `.local_watch_state.json.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、full build 後に再生成する。 |
| `.build_cache.json` | JSON object | `{"schema_version":1,"entries":{}}` | `components/builder.go` | 破損時は WARN を出し、cache miss として扱い、成功後に再生成する。 |
| `.build_cache/pages/` | directory | 空ディレクトリ | `components/builder.go` | entry 不一致または読み取り不能 file は miss とし、他 entry は継続使用する。 |
| `.dependency_manifest.json` | JSON object | `{"pages":{}}` | `components/builder.go` / `components/runner.go` | 破損時は full build とし、成功後に再生成する。 |
| `.approval_queue` | JSON Lines | 空ファイル | `components/runner.go` / `components/api.go` | 読み込み可能な行のみ使用し、壊れた行は ERROR ログへ記録して無視する。 |
| `.build_trends.json` | JSON object | `{"schema_version":1,"samples":[],"summary":{"count":0,"avg_seconds":null,"median_seconds":null,"p95_seconds":null}}` | `components/runner.go` / `components/api.go` | `.build_trends.json.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、`.build_history` から再集計する。 |
| `.build_chain_config` | JSON object | `{"chains":[]}` | `components/runner.go` / `components/api.go` | `.build_chain_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、chain 無効として通常 build のみ継続する。 |
| `.repo_config` | JSON object | `{}` | `components/api.go` | `.repo_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、スクリプト定数へフォールバックする。 |
| `.config_log` | JSON Lines | 空ファイル | `components/api.go` | 読み込み可能な行のみ返し、壊れた行は無視する。 |
| `.access_log` | JSON Lines | 空ファイル | `components/api.go` | 読み込み可能な行のみ返し、壊れた行は無視する。 |
| `.api_access_log` | JSON Lines | 空ファイル | `components/api.go` | 読み込み可能な行のみ返し、壊れた行は無視する。秘密情報は記録しない。 |
| `.webhook_secret` | text | 不在 | `components/api.go` | 読み込み不能時は Webhook 受信を `501` で拒否する。 |
| `.webhook_events.json` | JSON Lines | 空ファイル | `components/api.go` | 読み込み可能な行のみ返し、壊れた行は無視する。 |
| `.access_control` | JSON object | `{"allow":[]}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.hooks` | JSON object | `{"hooks":[]}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.maintenance` | JSON object | `{"enabled":false,"reason":null,"since":null}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.api_tokens` | JSON object | `{"tokens":[]}` | `components/api.go` | `500` を返し、自動再生成しない。token 管理情報の消失による意図しない再許可を防ぐため、破損ファイルは上書きしない。 |
| `.alert_rules` | JSON object | `{"rules":[]}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.tag_rules` | JSON object | `{"rules":[]}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.pipeline_config` | JSON object | `{"extra_args":[],"env":{}}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.notes` | UTF-8 text | 空文字列 | `components/api.go` | 読み込み不能時は `500` を返し、自動上書きしない。 |
| `.smtp_config` | JSON object | SMTP 未設定値 | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.smtp_secret` | text | 不在 | `components/api.go` | 読み込み不能時は SMTP 送信を `422` で拒否する。 |
| `.dashboard_layout` | JSON object | `{"widgets":["status","stats","schedule","alerts","disk","rate_limit","snapshots","maintenance","queue"]}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |

`.build_logs/archive/` は gzip 圧縮済み build log の保存先ディレクトリである。初期値は空ディレクトリとし、`components/runner.go` または `POST /api/logs/archive` が必要時に作成する。圧縮済みファイル名は `{id}.json.gz` 固定とし、通常 `.build_logs/{id}.json` と同じ build id を表す。

JSON Lines ファイルは、1 行につき 1 JSON object とする。追記時は末尾に改行を必ず付ける。秘密情報を含む可能性のある `.admin_credentials`、`.totp_secret`、`.github_token`、`.webhook_secret`、`.smtp_secret` は mode `600` を必須とする。

**状態ファイル更新手順：**

1. 対象ファイルの `{name}.lock` を `O_CREATE|O_EXCL` で作成する。
2. ロック取得に失敗した場合は 100ms 間隔で最大 10 秒待つ。
3. 現在値を読み込み、schema と入力値を検証する。
4. 更新後 JSON を `{name}.tmp.{pid}` に UTF-8 / LF で書き出す。
5. ファイルを close し、通常状態ファイルは `0644`、秘密情報ファイルと lock file は `0600` に chmod する。
6. `os.Rename(tmp, target)` で置換する。
7. target file を open して `Sync` し、続けて親ディレクトリを open して `Sync` する。
8. ロックファイルを削除する。

手順 3〜7 の途中で失敗した場合は target を変更せず、tmp を削除し、ロックを削除して `500 Internal Server Error` を返す。`os.Rename` 後の `Sync` に失敗した場合は target を維持し、ERROR ログと `.config_log` へ失敗を記録して `500` を返す。複数ファイル更新 API は §22.0d の Write 列順にこの手順を実行し、途中失敗時は未処理ファイルを書き込まない。既に書き込んだファイルの自動ロールバックは行わず、`.config_log` に失敗内容を記録する。

**状態ファイル schema 厳格化契約：**

状態ファイルの読込、正規化、保存は以下に固定する。§22.0c または個別機能節で例外を明記していない限り、実装者判断で旧形式、未知 key、null、欠落配列を成功扱いにしてはならない。

| 対象 | 読込時 | 保存時 | 失敗時 |
|------|--------|--------|--------|
| 未知 key | JSON object に schema 未定義 key がある場合は破損扱いとする。例外は §22.0c で明記した旧形式正規化だけ。 | 未知 key を保存しない。既存未知 key を黙って削除して保存しない。 | read endpoint は `500 {"error":"State file is corrupted"}`。write endpoint は target を変更しない。 |
| 必須 key 不足 | 既定値補完が個別節で明記されていない場合は破損扱いとする。 | 必須 key はすべて明示保存する。 | 初期値再生成が §22.0a 表で指定されたファイルだけ再生成する。 |
| `null` | 型欄が `string/null`、`object/null`、`integer/null` 等で明示した key だけ許可する。 | nullable でない key に `null` を保存しない。 | validation error または破損扱い。 |
| 配列 | `[]` を既定値とする key は read adapter の戻り値で空配列を返してよい。 | 保存 API は配列 key を省略せず、空の場合も `[]` を明示する。 | 型不一致は `422` または `500`。 |
| 数値 | 整数 key は JSON number の整数だけ許可する。小数、指数表記由来の非整数、文字列数値は拒否する。 | 整数は JSON number として保存する。 | API 入力は `422`、状態ファイル読込は破損扱い。 |
| 時刻 | UTC ISO 8601 秒精度 `Z` だけ許可する。 | 保存前に UTC 秒精度へ丸める。ミリ秒、local timezone、offset 付き文字列を保存しない。 | API 入力は `422`、状態ファイル読込は破損扱い。 |
| mode | 秘密情報ファイルは `0600`、通常 JSON / JSON Lines は `0644`、directory は `0755` を標準とする。 | chmod 失敗時は成功扱いにしない。 | chmod 失敗は `500`。target を更新した後の chmod 失敗は ERROR ログに残す。 |
| 改行 | text / JSON / JSON Lines は LF で保存する。JSON object / array ファイルは末尾 LF 1 個を付ける。 | CRLF、BOM、末尾余分空白を新規保存しない。 | 入力 text が CRLF を含む場合の扱いは個別機能節に従う。 |

旧 schema からの正規化は、本ファイルに「旧 key」「変換後 key」「削除する key」「保存するか読み取り時だけか」を明記した場合だけ実装する。明記がない旧形式は破損扱いとし、黙って推測変換してはならない。

**API 状態読取アダプタ固定契約：**

`components/api.go` は、P0 / P1 endpoint の状態読取を下表の adapter 名と戻り値で実装する。各 adapter は Go 内部関数名として固定し、同じ状態ファイルを endpoint ごとに別ロジックで直接 parse してはならない。

| Adapter | 読取対象 | 正常戻り値 | 不在時 | 破損時 / 読込不能時 |
|---------|----------|------------|--------|---------------------|
| `readBuildStatus()` | `.build_status.json` | `BuildStatus` | `(nil, false, nil)` を返し、fallback 判定へ渡す。 | `(nil, true, ErrStateCorrupted)` または `ErrStateReadFailed`。 |
| `readBuildState()` | `.build_state` | `BuildState` | §22.0a の初期値を返す。 | `ErrStateCorrupted` または `ErrStateReadFailed`。 |
| `readBuildHistory()` | `.build_history` | `[]BuildHistoryEntry` | 空配列を返す。 | 行単位破損は除外し、ファイル読込不能だけ `ErrStateReadFailed`。 |
| `readBuildLog(id)` | `.build_logs/{id}.json`、`.build_logs/archive/{id}.json.gz` | `BuildLog` | 通常ログ不在時は archive を読む。両方不在は `ErrNotFound`。 | 対象 ID の通常ログまたは archive が破損している場合は `ErrStateCorrupted`。 |
| `readLatestBuildLogs(n,q)` | `.build_logs/`、`.build_logs/archive/` | `[]LogLine` | 空配列を返す。 | 個別ログ破損は除外し、`LOG_SKIP_CORRUPT` を server log へ記録する。ディレクトリ読込不能は `ErrStateReadFailed`。 |
| `readPendingTransfers()` | `.pending_transfers` | `[]PendingTransfer` | 空配列を返す。 | `ErrStateCorrupted` または `ErrStateReadFailed`。 |
| `readCircuitState()` | `.build_circuit_state` | `BuildCircuitState` | §22.0a の初期値を返す。 | `ErrStateCorrupted` または `ErrStateReadFailed`。 |
| `readBuildLock()` | `.build_lock` | `BuildLockState` | `running=false` を返す。 | 形式不正、PID 判定不能、OS 判定失敗は `running=true, stale=false, valid=false` として返し、build command は `409`。 |

`GET` endpoint は上表の adapter を read-only で呼び出し、状態ファイルの作成、削除、退避、chmod、正規化、再生成、破損行の除去書き戻しを行ってはならない。`GET` endpoint が `{name}.lock` を検出しても、`.build_lock` 以外の lock file は待機条件やエラー条件にせず、rename 済み target をそのまま読む。write endpoint は §22.0a の状態ファイル更新手順に従う。

**API 状態読取 priority：**

| Endpoint | 読取順 | 正常時 response 算出 | 不在時 | 破損時 / 読込不能時 |
|----------|--------|----------------------|--------|---------------------|
| `GET /api/status` | `readBuildStatus()` → 不在時だけ `readBuildHistory()`、`readBuildState()`、`readBuildLock()`、`readPendingTransfers()`、`readCircuitState()` | `.build_status.json` がある場合は同ファイルを正とし、`running` だけ `.build_lock` が valid running の場合に `true` へ上書きする。 | `.build_status.json` 不在時は fallback で `status`、`last_*`、`running`、`pending_transfers_count`、`circuit_*` を算出する。履歴なしは `status:"none"`。 | `.build_status.json` 破損は `500 {"error":"State file is corrupted"}`。fallback 中の必須読取破損も `500`。 |
| `GET /api/health` | `readBuildStatus()`、`readBuildLock()`、出力サイト確認 | API process、出力サイト、直近 build 状態を `ok` / `warn` / `error` で返す。 | `.build_status.json` 不在は `status:"degraded"`。 | `.build_status.json` 破損は例外的に `200` とし、`status:"degraded"`、該当 item を `error` にする。 |
| `GET /api/history` | `readBuildHistory()` | 有効行だけを新しい順に sort し、query filter 後に paging する。 | `total:0`、`pages:0`、`history:[]`。 | 行単位破損は除外し、`BUILD_HISTORY_SKIP_CORRUPT` を server log へ記録する。ファイル読込不能は `500 {"error":"State file read failed"}`。 |
| `GET /api/history/{id}/log` | `readBuildLog(id)` | 対象 ID の log object を返す。 | 通常ログと archive の両方が不在なら `404 {"error":"Not found"}`。 | 対象 ID の log 破損は `500 {"error":"State file is corrupted"}`。 |
| `GET /api/logs` | `readLatestBuildLogs(n,q)` | 最新 log line を時系列順へ正規化し、`q` 指定時は部分一致で絞り込む。 | `{"lines":[]}`。 | 個別 log 破損は除外する。ディレクトリ読込不能は `500 {"error":"State file read failed"}`。 |
| `GET /api/queue` | `readBuildState()`、必要時 `.server_config` | `queued`、`running`、`queue_max_size` を返す。 | `.build_state` 不在は `queued:[]`、`running:false`。`.server_config` 不在は既定 queue 上限。 | `.build_state` 破損は `500 {"error":"State file is corrupted"}`。 |
| `POST /api/build` / `POST /api/build/force` | `readCircuitState()` → `readBuildLock()` → `readBuildState()` | circuit closed かつ lock 非実行なら `.build_state` を更新し、build 開始または queue 追加を返す。 | `.build_state` / `.build_circuit_state` 不在は初期値。 | circuit / state 破損は `500`。`.build_lock` が valid running、形式不正、PID 判定不能の場合は `409 {"error":"Conflict"}`。 |
| `POST /api/cancel` | `readBuildLock()`、`readBuildState()` | 実行中 build を cancel request 状態へ更新する。 | lock 不在かつ running false は `409 {"error":"Conflict"}`。 | `.build_state` 破損は `500`。`.build_lock` 形式不正または PID 判定不能は `409`。 |
| `POST /api/circuit-breaker/reset` | `readCircuitState()` | `open:false`、`consecutive_failures:0`、`opened_at:null`、`last_error:null` を atomic write する。 | 不在は初期値から reset 後値を書き込む。 | 読取破損は `500 {"error":"State file is corrupted"}` とし、上書きしない。 |

`ErrStateCorrupted` は JSON parse 失敗、schema_version 不一致、必須 key 不足、型不一致、列挙値不一致、UTC 時刻形式不一致のいずれかで返す。`ErrStateReadFailed` は permission denied、通常ファイルではない path、gzip 読込失敗、I/O error で返す。API response body はそれぞれ `{"error":"State file is corrupted"}`、`{"error":"State file read failed"}` 固定とし、path、Go error、ファイル内容を含めない。

JSON Lines adapter は空行、JSON parse 失敗、JSON object 以外、必須 key 不足、型不一致の行を壊れた行として除外する。除外後に sort、filter、paging、`total`、`pages` を算出する。壊れた行の存在は response body に含めず、server log に固定コード、path、1 始まりの line number だけを記録する。

`.build_lock` の PID が存在しない場合、`readBuildLock()` は `running=false, stale=true, valid=true` を返す。read-only endpoint は stale lock を削除しない。build command は開始前に `.build_lock` を再読込し、同じ stale 判定なら `.build_lock` だけを削除してから新規 lock を作成する。削除失敗時は `409 {"error":"Conflict"}` とし、`.build_state` を変更しない。

### 22.0b 入力検証共通仕様

API 実装は以下の検証を共通で行う。違反時は、エンドポイント固有の指定がない限り `422 Unprocessable Entity` と `{"error":"Validation failed","details":[{"field":"<field>","message":"<reason>"}]}` を返す。

| 対象 | 検証条件 |
|------|----------|
| `id` パスパラメータ | `^[A-Za-z0-9_-]{1,64}$` に一致すること。`/`、`.`、空文字は禁止。 |
| `page` | 1 以上の整数。 |
| `per_page` | 1 以上 100 以下の整数。 |
| `limit` | 1 以上 200 以下の整数。 |
| `offset` | 0 以上の整数。 |
| `days` | 1 以上 366 以下の整数。 |
| `n` | 1 以上 1000 以下の整数。 |
| 日付 | `YYYY-MM-DD` 形式で、存在する暦日であること。 |
| 時刻 | 0 以上 23 以下の整数。 |
| URL | `http://` または `https://` で始まること。Webhook URL は本番用途では `https://` のみ許可し、`http://` はローカル検証用途のみ許可する。 |
| ファイルパス | 絶対パスのみ許可する。`..` を含むパス、NUL 文字、空文字は禁止。 |
| タグ | 1 件 1〜32 文字、最大 20 件。重複は除去して保存する。 |
| コメント | 最大 2000 文字。空文字 `""` はコメント削除として扱う。 |
| メールアドレス | `local@domain` 形式で、空白を含まないこと。 |
| CIDR | IPv4 アドレスまたは IPv4 CIDR として解釈できること。 |
| コマンド引数配列 | `string[]` とし、1 要素以上 32 要素以下。各要素は 1〜256 文字。実行は `/bin/sh -c` を使わず、Go 標準ライブラリ `os/exec` の `exec.CommandContext(args[0], args[1:]...)` とする。 |

**endpoint 別 query 検証上書き：**

共通検証値と endpoint 個別節の値が異なる場合は、下表を優先する。下表にない query は §22.0b の共通検証を使用する。

| Endpoint | query | 既定値 | 許容値 | 補足 |
|----------|-------|--------|--------|------|
| `GET /api/history` | `page` | `1` | 1 以上 | 整数文字列だけ許可する。 |
| `GET /api/history` | `per_page` | `20` | 1〜100 | `0`、負数、小数、指数表記は禁止。 |
| `GET /api/logs` | `n` | `100` | 1〜1000 | `q` は空文字を許可する。 |
| `GET /api/logs/search` | `q` | `""` | 0〜500 文字 | 空文字は全件検索ではなく level/from/to のみ検索として扱う。 |
| `GET /api/logs/search` | `from`, `to` | `""` | 空文字または `YYYY-MM-DD` | `from > to` は `422`。 |
| `GET /api/logs/search` | `level` | `""` | `info`, `warn`, `warning`, `error`, `debug`, 空文字 | 大文字小文字は区別しない。 |
| `GET /api/api-access-log` | `limit` | `100` | 1〜1000 | §27.6 を優先する。 |
| `GET /api/api-access-log` | `offset` | `0` | 0 以上 | 整数文字列だけ許可する。 |
| `GET /api/webhook-events` | `limit` | `50` | 1〜1000 | §27.13 を優先する。 |
| `GET /api/webhook-events` | `offset` | `0` | 0 以上 | 整数文字列だけ許可する。 |
| `GET /api/audit-log` | `limit` | `100` | 1〜200 | §27.44 を優先する。 |
| `GET /api/audit-log` | `offset` | `0` | 0 以上 | 整数文字列だけ許可する。 |
| `GET /api/stats` | `days` | `7` | 1〜366 | 整数文字列だけ許可する。 |
| `GET /api/stats/timeline` | `days` | `30` | 1〜366 | 整数文字列だけ許可する。 |
| `GET /api/stats/build-duration` | `n` | `10` | 1〜1000 | 整数文字列だけ許可する。 |
| `GET /api/stats/build-trends` | `n` | `100` | 1〜1000 | 整数文字列だけ許可する。 |

**入力検証 details 固定：**

| ケース | `details[].field` | `details[].message` |
|--------|-------------------|---------------------|
| 必須 body key 不足 | key 名 | `required` |
| 未知 body key | key 名 | `unknown field` |
| 型不一致 | key 名 | `invalid type` |
| 範囲外 | key 名または query 名 | `out of range` |
| enum 不一致 | key 名または query 名 | `invalid value` |
| path parameter 不正 | parameter 名 | `invalid path parameter` |
| body 全体が object でない | `$` | `object required` |
| query が整数でない | query 名 | `integer required` |
| 日付の暦日不正 | query 名または key 名 | `invalid date` |

複数エラーがある場合、body key は JSON object の出現順、query は URL query の出現順、path parameter は route 定義順で並べる。body、query、path にまたがる場合は body → query → path の順とする。SDK と UI は `details[].field` と `details[].message` をそのまま扱うため、実装者判断で文言を言い換えてはならない。

### 22.0c 主要状態ファイル schema

本節の schema は、API 実装、SDK 型、標準管理ツール表示、バックアップ/リストアの基準である。ここに定義したキー以外を保存してはならない。追加キーを追加する場合は、型、既定値、読み書き API、既存データの扱いを本節へ追記してから実装する。

**`.server_config` schema：**

| キー | 型 | 既定値 | 許容値 | 読み書き API | 説明 |
|------|----|--------|--------|--------------|------|
| `log_max_lines` | integer | `500` | 1〜10000 | `GET/POST /api/config` | `GET /api/logs` が返す最大行数。 |
| `history_max_count` | integer | `100` | 1〜10000 | `GET/POST /api/config` | `.build_history` の通常表示上限。削除処理の上限ではない。 |
| `build_timeout_seconds` | integer | `300` | 1〜86400 | `GET/POST /api/config` | 手動/自動ビルドのタイムアウト秒数。 |
| `log_retention_days` | integer | `30` | 0〜3650 | `GET/POST /api/config`, `POST /api/logs/cleanup` | `0` は自動削除なし。 |
| `log_level` | string | `"INFO"` | `"INFO"` / `"DEBUG"` / `"WARNING"` / `"ERROR"` | `GET/POST /api/config`, `POST /api/log-level` | `components/api.go` のランタイムログレベル。 |
| `pat_expires_at` | string/null | `null` | `YYYY-MM-DD` または `null` | `GET/POST /api/config` | PAT 期限表示・診断用。 |
| `snapshots_keep` | integer | `5` | 0〜100 | `GET/POST /api/config` | `0` はスナップショット保存無効。 |
| `queue_max_size` | integer | `3` | 0〜100 | `GET/POST /api/config`, `GET /api/queue` | `0` はキュー無効。 |
| `build_retry_max` | integer | `0` | 0〜10 | `GET/POST /api/config` | ビルド失敗時の自動リトライ最大回数。`0` は無効。 |
| `build_retry_base_seconds` | integer | `5` | 1〜3600 | `GET/POST /api/config` | 自動リトライ backoff 基底秒数。待機秒数は `base * attempt` とする。 |
| `commit_status_enabled` | boolean | `false` | `true` / `false` | `GET/POST /api/config` | GitHub Commit Status API 送信の有効/無効。 |
| `commit_status_context` | string | `"Adlaire CI"` | 1〜100 文字 | `GET/POST /api/config` | GitHub commit status の `context`。 |
| `commit_status_target_url` | string/null | `null` | `http://` または `https://` の URL、または `null` | `GET/POST /api/config` | Commit Status の `target_url`。`null` の場合は送信 payload から省略する。 |
| `log_archive_after_days` | integer | `0` | 0〜3650 | `GET/POST /api/config`, `POST /api/logs/archive` | `0` は archive 無効。指定日数より古い通常 build log を gzip 圧縮する。 |
| `build_trend_keep_count` | integer | `1000` | 10〜10000 | `GET/POST /api/config` | `.build_trends.json` に保持する trend sample 件数。 |
| `duration_anomaly` | object | `{"enabled":false,"min_samples":20,"avg_multiplier":2.0,"p95_multiplier":1.5}` | §27.38 | `GET/POST /api/config` | build 所要時間異常検知の設定。 |
| `force_build_interval_hours` | integer | `0` | 0〜8760 | `POST /api/schedule/force-interval`, `GET /api/schedule` | `0` は強制再ビルド無効。 |
| `build_cooldown_seconds` | integer | `0` | 0〜86400 | `POST /api/schedule/cooldown`, `GET /api/schedule` | `0` はクールダウン無効。 |
| `schedule_interval_seconds` | integer | `300` | 30〜86400 | `POST /api/schedule/interval`, `GET /api/schedule` | systemd timer 更新値。 |
| `schedule_paused` | boolean | `false` | `true` / `false` | `POST /api/schedule/pause`, `POST /api/schedule/resume`, `GET /api/schedule` | 自動ポーリング停止状態。 |
| `allowed_hours` | object/null | `null` | `{"from":0〜23,"to":0〜23}` または `null` | `POST /api/schedule/allowed-hours`, `GET /api/schedule` | UTC の自動ビルド許可時間帯。 |
| `session_timeout_seconds` | integer | `28800` | 300〜2592000 | `GET/POST /api/config` | 新規 session の有効期限秒数。既存 session の `expires_at` は変更しない。 |
| `api_rate_limit` | object | `{"enabled":true,"groups":{"login":{"window_seconds":60,"max_requests":10},"read":{"window_seconds":60,"max_requests":600},"trigger":{"window_seconds":60,"max_requests":60},"operate":{"window_seconds":60,"max_requests":120},"config":{"window_seconds":60,"max_requests":60},"admin":{"window_seconds":60,"max_requests":60}}}` | §27.47 | `GET /api/api-rate-limit`, `POST /api/api-rate-limit`, `GET/POST /api/config` | API rate limit の endpoint group 別固定窓設定。 |

`.server_config` の `POST /api/config` では `force_build_interval_hours`、`build_cooldown_seconds`、`schedule_interval_seconds`、`schedule_paused`、`allowed_hours` を直接更新してはならない。これらは専用スケジュール API からのみ更新する。

**`.notify_config` schema：**

| キー | 型 | 既定値 | 許容値 | 説明 |
|------|----|--------|--------|------|
| `webhooks` | object[] | `[]` | 下記 Webhook object | 互換通知先一覧。`channels` が空の場合、runner は `webhooks` を webhook channel として扱う。 |
| `channels` | object[] | `[]` | 下記 Channel object | 統一通知 channel 一覧。`channels` が存在する場合、runner は `channels` を優先し、`webhooks` / `email` は互換表示用として扱う。 |
| `on` | string[] | `[]` | `"start"`, `"success"`, `"failure"`, `"deploy_failure"`, `"weekly_summary"`, `"approval_required"`, `"duration_anomaly"`, `"config_corrupt"` | 通知イベント。重複は除去する。 |
| `summary` | object | 下記 Summary object | 下記 | 定期サマリー設定。 |
| `email` | object | 下記 Email object | 下記 | メール通知設定。SMTP 詳細は `.smtp_config` / `.smtp_secret` を正とする。 |

Channel object:

| キー | 型 | 既定値 | 許容値 | 説明 |
|------|----|--------|--------|------|
| `id` | string | 自動採番 | `n` + 数字、または 1〜64 文字の英数字 `_` `-` | channel 識別子。 |
| `type` | string | 必須 | `"webhook"` / `"email"` / `"command"` | 送信方式。 |
| `label` | string | `""` | 0〜64 文字 | 管理画面表示名。 |
| `enabled` | boolean | `true` | boolean | `false` の channel へは送信しない。 |
| `on` | string[] | `[]` | top-level `on` と同じ、または `"*"` | 空配列の場合は top-level `on` に従う。 |
| `config` | object | `{}` | type 別 schema | webhook url、email to、command_args 等。 |
| `retry_count` | integer | `2` | 0〜10 | retry 対象失敗時の追加試行回数。 |
| `retry_interval_seconds` | integer | `30` | 1〜3600 | 再試行間隔。 |

Webhook object:

| キー | 型 | 既定値 | 許容値 | 説明 |
|------|----|--------|--------|------|
| `url` | string | 必須 | URL 検証に従う | 送信先 URL。 |
| `label` | string | `""` | 0〜64 文字 | 管理画面表示名。 |
| `enabled` | boolean | `true` | boolean | `false` の宛先へは送信しない。 |
| `on` | string[] | `[]` | `"start"`, `"success"`, `"failure"`, `"deploy_failure"`, `"weekly_summary"`, `"approval_required"`, `"duration_anomaly"`, `"config_corrupt"`, `"*"` | この宛先が受け取るイベント。空配列の場合は top-level `on` に従う。 |
| `payload_template` | string/null | `null` | 0〜10000 文字または `null` | `null` は標準 payload。 |
| `retry_count` | integer | `2` | 0〜10 | 送信失敗時の追加試行回数。 |
| `retry_interval_seconds` | integer | `30` | 1〜3600 | 再試行間隔。 |
| `secret` | string/null | `null` | 1〜256 文字または `null` | 保存時は平文保存可。ただし GET/backup では `"***"` へマスクする。 |

Summary object:

| キー | 型 | 既定値 | 許容値 |
|------|----|--------|--------|
| `enabled` | boolean | `false` | boolean |
| `interval` | string | `"weekly"` | `"weekly"` 固定 |
| `hour` | integer | `9` | 0〜23 |
| `day_of_week` | integer | `1` | 0〜6 |

Email object:

| キー | 型 | 既定値 | 許容値 |
|------|----|--------|--------|
| `enabled` | boolean | `false` | boolean |
| `to` | string[] | `[]` | メールアドレス配列、最大 50 件 |
| `on` | string[] | `[]` | `"start"`, `"success"`, `"failure"`, `"duration_anomaly"` |

**`.branch_config` schema：**

```json
{
  "branch_targets": [
    {
      "branch": "main",
      "target_file": "docs",
      "sha_file": "/opt/adlaire-builder/.last_sha",
      "src": "/opt/adlaire-builder/repo/docs",
      "out": "/opt/adlaire-builder/dist/site",
      "deploy_targets": [
        { "host": "192.0.2.1", "user": "deploy", "dest_dir": "/var/www/html/" }
      ]
    }
  ]
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `branch_targets` | object[] | 必須 | 0〜50 件 | 永続ファイルの正本 key。空配列は `.branch_config` 削除と同義。 |
| `branch` | string | 必須 | 1〜128 文字、`refs/heads/` は含めない | GitHub branch 名。 |
| `target_file` | string | 必須 | 相対パス、`..` 禁止 | GitHub リポジトリ内の監視対象ファイル。 |
| `sha_file` | string | 必須 | 絶対パス | 対象 branch/file の SHA キャッシュ。 |
| `src` | string | 必須 | 絶対パス | blob 本文の書き出し先。 |
| `out` | string | 必須 | 絶対パス | ビルド成果物パス。 |
| `deploy_targets` | object[] | 必須 | 0〜20 件 | SSH 転送先。空配列は転送なし。 |
| `deploy_targets[].host` | string | 必須 | 1〜255 文字 | SSH host。 |
| `deploy_targets[].user` | string | 必須 | 1〜64 文字 | SSH user。 |
| `deploy_targets[].dest_dir` | string | 必須 | 絶対パス | 転送先ディレクトリ。 |

`.branch_config` が不在の場合、`GET /api/branch-config` は `source: "default"` と `BRANCH_TARGETS` の定数値を返す。`.branch_config` が存在する場合、`source: "file"` とファイル内容を返す。API request / response の表示名として `branches` を使う場合でも、永続ファイルへ保存する key は必ず `branch_targets` とする。API は `branches: []` または `branch_targets: []` を `.branch_config` の空配列保存として扱ってはならない。`POST /api/branch-config` で空配列を受け取った場合は `.branch_config` を削除し、default 復帰として扱う。

**`.repo_config` schema：**

| キー | 型 | 既定値 | 許容値 | 説明 |
|------|----|--------|--------|------|
| `owner` | string | スクリプト定数 `OWNER` | 1〜100 文字 | GitHub owner。 |
| `repo` | string | スクリプト定数 `REPO` | 1〜100 文字 | GitHub repository。 |
| `branch` | string | スクリプト定数 `BRANCH` | 1〜128 文字 | 単一ターゲット用 branch。 |
| `target_file` | string | スクリプト定数 `TARGET_FILE` | 相対パス、`..` 禁止 | 単一ターゲット用監視ファイル。 |
| `updated_at` | string | 更新時刻 | ISO 8601 | 最終更新日時。 |

`POST /api/repo-config` は指定されたキーのみ更新する。未指定キーは既存値を保持する。全キーが未指定の場合は `422` を返す。

**`.totp_secret` schema：**

```json
{
  "enabled": true,
  "secret_base32": "JBSWY3DPEHPK3PXP",
  "confirmed_at": "2026-09-15T10:00:00Z",
  "last_accepted_step": 59652320
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `enabled` | boolean | 必須 | boolean | TOTP 有効状態。 |
| `secret_base32` | string/null | 必須 | RFC 4648 base32、padding なし、16〜64 文字、または `null` | TOTP secret。API response、log、backup へ平文出力しない。 |
| `confirmed_at` | string/null | 必須 | ISO 8601 または `null` | TOTP 有効化完了日時。 |
| `last_accepted_step` | integer/null | 必須 | Unix time 30 秒 step または `null` | 同一 code 再利用防止。 |

**`.audit_log` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `timestamp` | string | 必須 | ISO 8601 | 発生日時。 |
| `request_id` | string | 必須 | 16 byte hex | API request 単位の識別子。 |
| `actor_type` | string | 必須 | `"admin"` / `"api_token"` / `"system"` / `"anonymous"` | 操作者種別。 |
| `actor_id` | string/null | 必須 | `"admin"`、token id、`"system"`、または `null` | 操作者。secret 本体は保存しない。 |
| `action` | string | 必須 | §27.44 | 操作種別。 |
| `target_type` | string | 必須 | §27.44 | 対象種別。 |
| `target_id` | string/null | 必須 | 対象 id または `null` | 対象識別子。 |
| `result` | string | 必須 | `"success"` / `"failure"` / `"denied"` | 結果。 |
| `remote_addr` | string/null | 必須 | IP 文字列または `null` | 接続元。 |
| `message` | string/null | 必須 | 0〜500 文字または `null` | 固定文言。secret、token、password は保存しない。 |

**`.api_rate_state` schema：**

```json
{
  "windows": {
    "ip:127.0.0.1:login": {
      "window_start": "2026-09-15T10:00:00Z",
      "count": 1
    }
  }
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `windows` | object | 必須 | key は `{dimension}:{id}:{endpoint_group}` | 固定窓状態。 |
| `window_start` | string | 必須 | ISO 8601 | 現在窓の開始時刻。 |
| `count` | integer | 必須 | 0 以上 | 現在窓内リクエスト数。 |

**`.api_tokens` schema：**

```json
{
  "tokens": [
    {
      "id": "tok000001",
      "label": "監視用",
      "scopes": ["read"],
      "token_hash": "<sha256_hex>",
      "created_at": "2026-09-15T10:00:00Z",
      "last_used_at": null,
      "expires_at": null,
      "revoked_at": null
    }
  ]
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `tokens` | object[] | 必須 | 0〜100 件 | 発行済み API token 一覧。 |
| `id` | string | 必須 | `tok` + 3 桁以上の数字 | token 識別子。 |
| `label` | string | 必須 | 1〜64 文字 | 表示名。 |
| `scopes` | string[] | 必須 | `read`, `trigger`, `operate`, `config`, `admin` の 1〜5 件 | token に許可する scope。 |
| `token_hash` | string | 必須 | SHA-256 hex | token 本体は保存しない。 |
| `created_at` | string | 必須 | ISO 8601 | 作成日時。 |
| `last_used_at` | string/null | 必須 | ISO 8601 または `null` | 最終使用日時。 |
| `expires_at` | string/null | 必須 | ISO 8601 または `null` | 有効期限。`null` は無期限。 |
| `revoked_at` | string/null | 必須 | ISO 8601 または `null` | 失効日時。`null` は有効。 |

`POST /api/tokens` は token 本体を `act_` + 32 byte 相当のランダム文字列として生成し、レスポンス時に 1 回だけ返す。保存する値は `token_hash` のみとする。`DELETE /api/tokens/{id}` は物理削除せず、`revoked_at` を現在時刻へ更新する。旧 `scope` 文字列が存在する場合は読み込み時に `scopes:[scope]` へ正規化して保存し直す。

**`.api_tokens` 実装固定値：**

| 項目 | 仕様 |
|------|------|
| token id 採番 | `tok` + 6 桁連番とする。既存最大番号が `tok000123` の場合、次は `tok000124` とする。連番抽出不能な id は衝突確認対象には含めるが、最大番号算出には使わない。 |
| token 本体 | `crypto/rand` 32 bytes を `encoding/base64.RawURLEncoding` で文字列化し、先頭に `act_` を付ける。保存前 hash は prefix を含む token 全体に対して `sha256` を計算する。 |
| `label` | 1〜64 文字。前後空白は保存前に除去する。除去後が空文字なら `422`。 |
| `scopes` | 1〜5 件。重複は除去し、保存値は `read`, `trigger`, `operate`, `config`, `admin` の順に正規化する。 |
| `expires_at` | `null` または現在時刻より後の UTC ISO 8601。過去または現在時刻は `422`。 |
| 一覧順 | `GET /api/tokens` は `created_at` 降順、同時刻は `id` 昇順で返す。 |
| 失効済み表示 | `GET /api/tokens` は失効済み token も返す。token 本体と `token_hash` は返さない。 |
| 認証時更新 | 有効 token 認証成功時だけ `last_used_at` を現在時刻へ更新する。期限切れ、失効済み、hash 不一致では更新しない。 |
| 破損行相当 | `.api_tokens.tokens` 内の個別 record が schema 不正の場合、認証と一覧は `500` を返し、自動補正しない。旧 `scope` から `scopes` への正規化だけは例外として許可する。 |

**`.maintenance` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `enabled` | boolean | 必須 | boolean | メンテナンス有効状態。 |
| `reason` | string/null | 必須 | 0〜500 文字または `null` | 理由。 |
| `since` | string/null | 必須 | ISO 8601 または `null` | 有効化日時。 |

`enabled: false` の場合、`reason` と `since` は `null` とする。`POST /api/maintenance/enable` は `enabled: true`、`reason`、`since` を同時に保存する。

**`.access_control` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `allow` | string[] | 必須 | IPv4 address または CIDR、0〜100 件 | 空配列は制限なし。保存時は入力順を保持し、重複は除去する。 |

`allow` の各要素は前後空白を除去してから検証する。空文字、IPv6、hostname、URL、CIDR prefix が 0〜32 以外、parse 不能な値は `422` とし、既存 `.access_control` を変更しない。

**`.hooks` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `hooks` | object[] | 必須 | 0〜50 件 | 登録済み hook。 |
| `id` | string | 必須 | §22.0e.2 の hook id | hook 識別子。 |
| `phase` | string | 必須 | `"pre"` / `"post"` | 実行 phase。 |
| `command_args` | string[] | 必須 | 1〜20 件 | shell を介さず `exec.CommandContext` に渡す引数配列。 |
| `enabled` | boolean | 必須 | boolean | `false` の hook は実行しない。 |
| `abort_on_failure` | boolean | 必須 | boolean | `pre` hook 失敗時だけ参照する。 |
| `timeout_seconds` | integer | 必須 | 1〜3600 | hook 単体の timeout。未指定作成時は `300`。 |

`command_args[0]` は 1〜256 文字、`command_args[1:]` の各要素は 1〜500 文字とし、NUL、改行、CR を禁止する。`command_args[0]` は絶対 path または PATH 解決可能なコマンド名に限定する。`phase` と `command_args` が既存 enabled hook と完全一致する場合、`POST /api/hooks` は `409 {"error":"Conflict"}` を返す。

**`.alert_rules` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `rules` | object[] | 必須 | 0〜100 件 | dashboard alert rule。 |
| `id` | string | 必須 | §22.0e.2 の alert rule id | rule 識別子。 |
| `metric` | string | 必須 | `success_rate_7d`, `avg_duration_seconds`, `last_build_age_hours`, `disk_usage_bytes` | 評価対象。 |
| `operator` | string | 必須 | `lt`, `gt`, `lte`, `gte` | 比較演算子。 |
| `threshold` | number | 必須 | 0 以上 | 比較値。 |
| `level` | string | 必須 | `info`, `warn`, `error` | alert severity。 |
| `message` | string | 必須 | 1〜200 文字 | UI 表示文。secret を含めない。 |

**`.tag_rules` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `rules` | object[] | 必須 | 0〜100 件 | 自動タグ付け rule。 |
| `id` | string | 必須 | §22.0e.2 の tag rule id | rule 識別子。 |
| `condition` | string | 必須 | §15B の条件式 grammar | 評価条件。 |
| `tags` | string[] | 必須 | 1〜20 件、各 1〜50 文字 | 付与するタグ。重複は除去する。 |

`condition` は `変数 空白 演算子 空白 値` の 1 条件だけを許可する。`&&`、`||`、括弧、関数呼び出し、正規表現、算術式は `422` とする。

**`.pipeline_config` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `extra_args` | string[] | 必須 | 0〜50 件 | `components/builder.go` に渡す追加 CLI 引数。 |
| `env` | object | 必須 | key/value は下記 | builder process に追加する環境変数。 |

`extra_args` は空文字、NUL、改行、CR を禁止し、`--src`、`--out`、`--build-id`、`--commit-sha`、`--build-at`、`--version`、`--help` を指定してはならない。`env` key は `^[A-Z_][A-Z0-9_]{0,63}$`、value は 0〜1000 文字とし、`PATH`、`HOME`、`SHELL`、`USER`、`GITHUB_TOKEN`、`ADLAIRE_TOKEN` は上書き禁止とする。

**`.dashboard_layout` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `widgets` | string[] | 必須 | §16D の widget id、1〜9 件 | 表示 widget 順序。重複禁止。 |

**`.smtp_config` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `host` | string/null | 必須 | hostname または IP、または `null` | SMTP server。 |
| `port` | integer | 必須 | 1〜65535 | SMTP port。 |
| `user` | string/null | 必須 | 0〜255 文字または `null` | SMTP user。 |
| `tls` | boolean | 必須 | boolean | STARTTLS または TLS 使用。 |
| `from` | string/null | 必須 | email address または `null` | From address。 |
| `to` | string[] | 必須 | email address、0〜50 件 | 送信先。 |
| `on` | string[] | 必須 | `start`, `success`, `failure`, `duration_anomaly` | 送信イベント。 |
| `enabled` | boolean | 必須 | boolean | メール通知有効状態。 |

`.smtp_secret` は UTF-8 text とし、末尾 LF なしで password 本体だけを保存する。mode は `0600` 固定。`POST /api/smtp-config` で `password` が未指定の場合、既存 `.smtp_secret` を変更しない。`password:null` は secret 削除を意味し、`.smtp_secret` が存在する場合だけ削除する。

**`.build_state` schema：**

```json
{
  "running": false,
  "current_build_id": null,
  "queued": [],
  "last_started_at": null,
  "last_finished_at": null,
  "weekly_summary_last_sent_at": null,
  "weekly_summary_sent_date": null
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `running` | boolean | 必須 | boolean | ビルド実行中状態。 |
| `current_build_id` | string/null | 必須 | build id または `null` | 実行中 build id。 |
| `queued` | object[] | 必須 | 0〜`queue_max_size` 件 | 待機中 build queue。 |
| `last_started_at` | string/null | 必須 | ISO 8601 または `null` | 最終開始日時。 |
| `last_finished_at` | string/null | 必須 | ISO 8601 または `null` | 最終完了日時。 |
| `weekly_summary_last_sent_at` | string/null | 必須 | ISO 8601 または `null` | 週次サマリー最終送信日時。 |
| `weekly_summary_sent_date` | string/null | 必須 | `YYYY-MM-DD` または `null` | 週次サマリー二重送信防止日付。 |

Queue entry:

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `id` | string | 必須 | `q` + 3 桁以上の数字 | queue id。 |
| `trigger` | string | 必須 | `"manual"`, `"webhook"` | 起動種別。強制実行は `"manual"` と `payload.force=true` で表す。 |
| `queued_at` | string | 必須 | ISO 8601 | queue 追加日時。 |
| `requested_by` | string | 必須 | `"api"`, `"webhook"` | queue 追加元。 |
| `payload` | object | 必須 | JSON object | force/webhook 等の追加情報。不要時は `{}`。 |

`running: false` の場合、`current_build_id` は `null` とする。`DELETE /api/queue` は `queued` を空配列へ置換し、`running` と `current_build_id` は変更しない。

**`.build_circuit_state` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `open` | boolean | 必須 | boolean | `true` の間は自動ポーリングを停止する。 |
| `consecutive_failures` | integer | 必須 | 0 以上 | 連続失敗回数。 |
| `opened_at` | string/null | 必須 | ISO 8601 または `null` | open に遷移した日時。 |
| `last_failure_at` | string/null | 必須 | ISO 8601 または `null` | 最終失敗日時。 |
| `last_error` | string/null | 必須 | 文字列または `null` | 最終失敗理由。 |

初期値は `{"open":false,"consecutive_failures":0,"opened_at":null,"last_failure_at":null,"last_error":null}` とする。`POST /api/circuit-breaker/reset` は初期値へ戻す。

**`.config_log` JSON Lines schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `at` | string | 必須 | ISO 8601 | 変更日時。 |
| `type` | string | 必須 | 状態ファイル種別 | 例: `"server_config"`, `"notify_config"`, `"repo_config"`。 |
| `action` | string | 必須 | `"create"`, `"update"`, `"delete"` | 変更種別。 |
| `diff` | object | 必須 | `{key:[before,after]}` | 変更前後。秘密情報は `"***"`。 |
| `diff_text` | string | 必須 | 1 文字以上 | 人間向け差分。秘密情報は `"***"`。 |

**`.access_log` JSON Lines schema：**

各行はログイン、ログアウト、API token 作成/失効、read token 認証の監査イベントを表す JSON object とする。秘密情報、セッショントークン、API token 本体を保存してはならない。

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `at` | string | 必須 | ISO 8601 | イベント日時。 |
| `action` | string | 必須 | `"login"`, `"logout"`, `"token_create"`, `"token_revoke"`, `"token_auth"` | 監査イベント種別。 |
| `result` | string | 必須 | `"success"`, `"failure"` | 成否。 |
| `session_id` | string/null | 必須 | 文字列または `null` | セッション識別用の短縮 ID。token 本体ではない。 |
| `token_id` | string/null | 必須 | API token id または `null` | API token 関連イベントの対象。 |
| `remote_addr` | string/null | 必須 | IP 文字列または `null` | 接続元。取得不能時は `null`。 |
| `reason` | string/null | 必須 | 文字列または `null` | 失敗理由。秘密情報を含めない。 |

**`.api_access_log` JSON Lines schema：**

各行は認証後 API、認証失敗 API、Webhook API の HTTP 呼び出し 1 件を表す JSON object とする。password、token、Webhook secret、SMTP password、request body の secret 値を保存してはならない。

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `at` | string | 必須 | ISO 8601 | response 送信直前の日時。 |
| `request_id` | string | 必須 | `req{YYYYMMDDHHmmss}-NNN` | API 呼び出し識別子。 |
| `method` | string | 必須 | HTTP method | `GET` / `POST` / `DELETE` 等。 |
| `path` | string | 必須 | `/api/...` | query を含まない path。 |
| `query` | object | 必須 | JSON object | 許可済み query key と値。秘密値は禁止。 |
| `status` | integer | 必須 | HTTP status code | response status。 |
| `duration_ms` | integer | 必須 | 0 以上 | handler 開始から response 確定までのミリ秒。 |
| `auth_type` | string | 必須 | `"session"`, `"api_token"`, `"webhook"`, `"none"` | 認証種別。 |
| `actor` | string/null | 必須 | `"admin"`、token id、`"webhook"`、または `null` | 操作者。token 本体は保存しない。 |
| `remote_addr` | string/null | 必須 | IP 文字列または `null` | 接続元。 |
| `user_agent` | string/null | 必須 | 文字列または `null` | 取得不能時は `null`。 |
| `error` | string/null | 必須 | エラーコードまたは `null` | 成功時は `null`。 |

**`.build_history` JSON Lines schema：**

各行は以下の JSON object とする。

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `id` | string | 必須 | `b{YYYYMMDDHHmmss}` | ビルド ID。 |
| `build_at` | string | 必須 | ISO 8601 | ビルド完了日時。 |
| `sha` | string/null | 必須 | Git SHA または `null` | 対象 blob / commit SHA。 |
| `status` | string | 必須 | `"success"`, `"failure"`, `"cancelled"`, `"hook_error"` | ビルド結果。 |
| `trigger` | string | 必須 | `"polling"`, `"force_interval"`, `"manual"`, `"webhook"`, `"retry_pending_transfer"`, `"startup_config_integrity"`, `"rollback"`, `"local_watch"`, `"approval"` | 起動種別。 |
| `output_size_bytes` | integer/null | 必須 | 0 以上または `null` | 成果物サイズ。 |
| `output_sha256` | string/null | 任意 | SHA-256 hex または `null` | 成果物チェックサム。 |
| `duration_seconds` | integer/null | 必須 | 0 以上または `null` | 所要時間。 |
| `retry_count` | integer | 任意 | 0 以上 | 最終成功または最終失敗までに実行した追加 retry 回数。未記録時は `0` と扱う。 |
| `commit_status_state` | string/null | 任意 | `"pending"`, `"success"`, `"failure"`, `"error"`, `null` | 最終 GitHub Commit Status 送信状態。未送信時は `null`。 |
| `flagged` | boolean | 必須 | boolean | 重要フラグ。 |
| `tags` | string[] | 必須 | タグ検証に従う | 手動/自動タグ。 |
| `comment` | string/null | 必須 | コメント検証に従う | コメント。 |
| `rollback_from` | string/null | 任意 | build id または `null` | rollback の元 build id。 |

**`.build_logs/{id}.json` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `id` | string | 必須 | build id | ファイル名 `{id}.json` と一致する。 |
| `started_at` | string | 必須 | ISO 8601 | ビルド開始日時。 |
| `finished_at` | string/null | 必須 | ISO 8601 または `null` | 完了前は `null`。 |
| `duration_seconds` | integer/null | 必須 | 0 以上または `null` | 完了前は `null`。 |
| `status` | string | 必須 | `"running"`, `"success"`, `"failure"`, `"cancelled"`, `"hook_error"` | 現在/最終状態。 |
| `trigger` | string | 必須 | `.build_history.trigger` と同じ | 起動種別。 |
| `branch` | string | 必須 | branch 名 | 対象 branch。 |
| `target_file` | string | 必須 | 相対パス | 対象ファイル。 |
| `sha` | string/null | 必須 | SHA または `null` | 対象 SHA。 |
| `commit_sha` | string/null | 必須 | SHA または `null` | commit SHA。 |
| `commit_message` | string/null | 必須 | 文字列または `null` | commit message。 |
| `commit_author` | string/null | 必須 | 文字列または `null` | commit author。 |
| `commit_at` | string/null | 必須 | ISO 8601 または `null` | commit 日時。 |
| `stdout` | string[] | 必須 | 0 件以上 | `pipeline.sh` stdout 行。 |
| `stderr` | string[] | 必須 | 0 件以上 | `pipeline.sh` stderr 行。 |
| `warnings` | string[] | 必須 | 0 件以上 | `[WARN]` 行または runner warning。 |
| `report` | object/null | 必須 | 下記 Report object または `null` | `[REPORT]` の解析結果。 |
| `output_size_bytes` | integer/null | 必須 | 0 以上または `null` | 成果物サイズ。 |
| `output_sha256` | string/null | 任意 | SHA-256 hex または `null` | 成果物チェックサム。 |
| `size_warn` | boolean | 必須 | boolean | サイズ警告。 |
| `transfer_verified` | boolean/null | 必須 | boolean または `null` | SSH 転送未実行時は `null`。 |
| `dry_run` | boolean | 必須 | boolean | 通常 build log では常に `false`。dry-run は build log を作成しないため、本 field が `true` の build log を新規作成してはならない。 |
| `attempts` | object[] | 必須 | 1 件以上 | build / deploy の試行履歴。dry-run は build log を作成しないため、本配列へ検証 attempt を保存しない。 |
| `retry_count` | integer | 必須 | 0 以上 | 追加 retry 回数。初回のみで終わった場合は `0`。 |
| `commit_status` | object/null | 必須 | CommitStatus object または `null` | GitHub Commit Status API 送信結果。無効時は `null`。 |
| `build_meta` | object | 必須 | BuildMeta object | 出力サイトへ埋め込んだ build metadata。 |
| `error` | string/null | 必須 | 文字列または `null` | 失敗理由。 |
| `comment` | string/null | 必須 | コメント検証に従う | コメント。 |
| `flagged` | boolean | 必須 | boolean | 重要フラグ。 |
| `tags` | string[] | 必須 | タグ検証に従う | タグ。 |

Report object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `headings` | integer | 必須 | 見出し数。 |
| `tables_count` | integer | 必須 | テーブル数。 |
| `code_blocks_count` | integer | 必須 | コードブロック数。 |
| `warnings_count` | integer | 必須 | 警告件数。 |
| `size_warn` | boolean | 必須 | 出力サイトサイズ警告。 |
| `broken_links` | integer | 必須 | 内部リンク不整合数。 |
| `heading_skips` | integer | 必須 | 見出しレベルスキップ数。 |
| `reading_time` | integer | 必須 | 推計読了時間。 |
| `theme` | string | 必須 | 使用 theme 名。 |
| `build_id` | string | 必須 | `[REPORT] build_id`。未指定時は空文字。 |
| `commit_sha` | string | 必須 | `[REPORT] commit_sha`。未指定時は空文字。 |
| `build_at` | string | 必須 | `[REPORT] build_at`。未指定時は空文字。 |

Attempt object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `attempt` | integer | 必須 | 初回は `1`。retry ごとに +1。 |
| `started_at` | string | 必須 | UTC ISO 8601。 |
| `finished_at` | string/null | 必須 | 完了時刻。dry-run は build log を作成しないため、本 field を保存しない。 |
| `stage` | string | 必須 | `"github"`, `"pipeline"`, `"deploy"` のいずれか。dry-run は build log を作成しないため `"dry_run"` stage を新規保存してはならない。 |
| `status` | string | 必須 | `"success"` または `"failure"`。 |
| `retryable` | boolean | 必須 | この失敗が retry 対象か。成功時は `false`。 |
| `error` | string/null | 必須 | 失敗理由。成功時は `null`。 |

CommitStatus object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `enabled` | boolean | 必須 | `.server_config.commit_status_enabled` の評価結果。 |
| `state` | string/null | 必須 | `"pending"`, `"success"`, `"failure"`, `"error"`、未送信時 `null`。 |
| `context` | string | 必須 | 送信 context。 |
| `target_url` | string/null | 必須 | 送信 target_url。省略時 `null`。 |
| `sent_at` | string/null | 必須 | 最終送信時刻。未送信時 `null`。 |
| `http_status` | integer/null | 必須 | GitHub API HTTP status。未送信時 `null`。 |
| `error` | string/null | 必須 | 送信失敗理由。成功時 `null`。 |

BuildMeta object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `build_id` | string | 必須 | HTML meta `adlaire-build-id` と同じ値。 |
| `commit_sha` | string | 必須 | HTML meta `adlaire-commit-sha` と同じ値。 |
| `build_at` | string | 必須 | HTML meta `adlaire-build-at` と同じ値。 |

### 22.0d API と状態ファイル対応表

API 実装では、下表の read/write 以外の状態ファイルを操作してはならない。複数ファイルを write する API は、表の順序で検証、バックアップ、書き込みを行い、途中失敗時は後続ファイルを書き込まない。

| API | Read | Write | 補足 |
|-----|------|-------|------|
| `POST /api/login` | `.admin_credentials`, `.totp_secret` | `.admin_credentials`, `.access_log`, `.audit_log` | TOTP 有効時は session を発行せず ticket を返す。 |
| `POST /api/login/totp` | `.admin_credentials`, `.totp_secret` | `.admin_credentials`, `.totp_secret`, `.access_log`, `.audit_log` | password verified ticket と TOTP code を検証し、成功時に session を発行する。 |
| `POST /api/logout` | メモリ上 session | メモリ上 session | ファイルは更新しない。 |
| `POST /api/change-password` | `.admin_credentials` | `.admin_credentials` | 現 session 以外をメモリから削除する。 |
| `GET /api/audit-log` | `.audit_log` | なし | 壊れた行は無視し、新しい順で返す。 |
| `GET /api/access-log` | `.access_log` | なし | 壊れた行は無視し、新しい順で返す。 |
| `GET /api/api-access-log` | `.api_access_log` | なし | 壊れた行は無視し、新しい順で返す。 |
| `GET /api/sessions` | メモリ上 session | なし | token 本体は返さない。 |
| `POST /api/sessions/revoke-all` | メモリ上 session | メモリ上 session, `.access_log` | 現 session 以外を削除する。 |
| `GET /api/auth/totp-status` | `.totp_secret` | なし | 単一 admin の TOTP 状態を返す。 |
| `POST /api/auth/totp-setup` | `.totp_secret` | なし | 仮 secret と otpauth URI をメモリで生成し、未確認のまま永続化しない。 |
| `POST /api/auth/totp-confirm` | `.totp_secret` | `.totp_secret`, `.audit_log` | 仮 secret を code 検証後に保存し、TOTP を有効化する。 |
| `DELETE /api/auth/totp` | `.totp_secret` | `.totp_secret`, `.audit_log` | code 検証後に TOTP を無効化する。 |
| `GET /api/status` | `.build_status.json`, `.build_state`, `.build_lock`, `.build_history` | なし | `.build_status.json` を第一参照元とする。不在時のみ `.build_state` と `.build_history` から後方互換の値を算出する。 |
| `POST /api/build` | `.server_config`, `.build_state`, `.build_lock`, `.maintenance`, `.build_circuit_state` | `.build_state` または queue | 実行中かつ queue 有効なら queue へ追加する。 |
| `POST /api/build/force` | `.server_config`, `.build_state`, `.build_lock`, `.maintenance`, `.build_circuit_state`, SHA cache | `.build_state`, SHA cache または queue | SHA reset と build trigger は同一ロック内で行う。 |
| `POST /api/build/cancel` | `.build_lock` | `.build_state`, `.build_logs/{id}.json` | 実行中でない場合は `409`。 |
| `GET /api/build/stream` | `.build_logs/{id}.json`, `.build_state` | なし | SSE 配信のみ。ログファイルは更新しない。 |
| `GET /api/logs` | `.build_logs/` | なし | 最新ログを読む。 |
| `GET /api/logs/search` | `.build_logs/` | なし | 横断検索のみ。 |
| `GET /api/logs/export` | `.build_logs/` | なし | JSON export。 |
| `POST /api/logs/cleanup` | `.server_config`, `.build_logs/` | `.build_logs/` | 削除対象のみ削除する。 |
| `POST /api/logs/archive` | `.server_config`, `.build_logs/` | `.build_logs/archive/`, `.build_logs/` | archive 対象を gzip 圧縮し、成功後に通常 log を削除する。 |
| `GET /api/history` | `.build_history` | なし | `page`、`per_page`、`trigger`、`tag`、`flagged` で絞り込み、ページングして返す。 |
| `GET /api/history/export` | `.build_history` | なし | 全件 export。 |
| `GET /api/history/{id}/log` | `.build_logs/{id}.json` | なし | ファイル破損時は `500`。 |
| `GET /api/history/{id}/comment` | `.build_logs/{id}.json` | なし | なし。 |
| `POST /api/history/{id}/comment` | `.build_logs/{id}.json` | `.build_logs/{id}.json`, `.config_log` | コメントだけ更新する。 |
| `POST /api/history/{id}/flag` | `.build_logs/{id}.json` | `.build_logs/{id}.json`, `.config_log` | flag だけ更新する。 |
| `POST /api/history/{id}/tags` | `.build_logs/{id}.json` | `.build_logs/{id}.json`, `.build_history`, `.config_log` | `.build_history` の同一 id にも反映する。 |
| `POST /api/history/{id}/rollback` | `.snapshots/{id}/`, `.server_config` | `.build_history`, `.build_logs/{new_id}.json` | rollback エントリを新規追加する。 |
| `GET /api/notify-config` | `.notify_config`, `.smtp_config` | なし | secrets はマスクする。 |
| `POST /api/notify-config` | `.notify_config` | `.notify_config`, `.config_log` | secret は GET で返さない。 |
| `GET /api/notify-log` | `.notify_log` | なし | 壊れた行は無視する。 |
| `POST /api/notify-test` | `.notify_config` | `.notify_log` | 送信結果を追記する。 |
| `POST /api/notify/weekly-summary` | `.notify_config`, `.build_logs/` | `.notify_log` | 宛先なしは `422`。 |
| `GET /api/config` | `.server_config` | なし | 既定値を merge して返す。 |
| `POST /api/config/validate` | `.server_config`, request body | なし | 保存せず検証結果だけ返す。`.config_log` も更新しない。 |
| `POST /api/config` | `.server_config` | `.server_config`, `.config_log` | 許可キーのみ更新する。 |
| `POST /api/log-level` | `.server_config` | `.server_config`, `.config_log` | `log_level` のみ更新する短縮 API。 |
| `GET /api/config-log` | `.config_log` | なし | 壊れた行は無視する。 |
| `GET /api/api-rate-limit` | `.server_config`, `.api_rate_state` | なし | rate limit 設定と現在 window summary を返す。 |
| `POST /api/api-rate-limit` | `.server_config` | `.server_config`, `.api_rate_state`, `.config_log`, `.audit_log` | policy を保存し、window state を初期化する。 |
| `GET /api/repo-info` | `.repo_config` | なし | 不在時は定数値を返す。 |
| `POST /api/repo-config` | `.repo_config` | `.repo_config`, `.config_log` | 未指定キーは保持する。 |
| `GET /api/branch-config` | `.branch_config` | なし | 不在時は default。 |
| `POST /api/branch-config` | `.branch_config` | `.branch_config`, `.config_log` | 空配列は `.branch_config` 削除。 |
| `GET /api/sysinfo` | 出力サイト, process start time | なし | 状態ファイルは更新しない。 |
| `GET /api/health` | `.build_history`, `.pending_transfers` | なし | 認証不要。 |
| `GET /api/stats` | `.build_history`, `.build_logs/` | なし | `days` の範囲を集計する。 |
| `GET /api/stats/timeline` | `.build_history` | なし | 日別集計のみ。 |
| `GET /api/stats/build-duration` | `.build_logs/` | なし | duration 集計のみ。 |
| `GET /api/output-meta` | `.build_history`, `.build_logs/`, `.build_logs/archive/`, 出力サイト | なし | 出力サイトと直近ログを集約する。通常 log 不在時は archive log を読む。 |
| `GET /api/pat-status` | `.github_token` | なし | 結果保存なし。 |
| `POST /api/pat-verify` | `.github_token` | なし | 結果保存なし。 |
| `POST /api/pat-update` | なし | `.github_token`, `.config_log` | token 値は `.config_log` でマスクする。 |
| `GET /api/backup` | `.server_config`, `.notify_config`, `.repo_config`, `.branch_config`, `.access_control`, `.hooks`, `.alert_rules`, `.tag_rules`, `.pipeline_config`, `.dashboard_layout`, `.smtp_config` | なし | secrets は `"***"` へマスクする。 |
| `POST /api/restore` | request body | `.server_config`, `.notify_config`, `.repo_config`, `.branch_config`, `.access_control`, `.hooks`, `.alert_rules`, `.tag_rules`, `.pipeline_config`, `.dashboard_layout`, `.smtp_config`, `.config_log` | restore 対象ファイルを検証後に表の順で置換する。 |
| `GET /api/schedule` | `.server_config` | なし | systemd 次回実行時刻は外部確認。 |
| `POST /api/schedule/interval` | `.server_config` | `.server_config`, `.config_log` | systemd timer 反映も行う。 |
| `POST /api/schedule/pause` | `.server_config` | `.server_config`, `.config_log` | 既に paused は `409`。 |
| `POST /api/schedule/resume` | `.server_config` | `.server_config`, `.config_log` | 稼働中は `409`。 |
| `POST /api/schedule/allowed-hours` | `.server_config` | `.server_config`, `.config_log` | `null` で解除。 |
| `POST /api/schedule/force-interval` | `.server_config` | `.server_config`, `.config_log` | `hours` を保存。 |
| `POST /api/schedule/cooldown` | `.server_config` | `.server_config`, `.config_log` | `seconds` を保存。 |
| `GET /api/dashboard` | `.build_status.json`, `.build_history`, `.build_state`, `.build_lock`, `.server_config`, `.alert_rules`, `.dashboard_layout`, 出力サイト, process start time | なし | 集約のみ。 |
| `GET /api/diagnostics` | `.github_token`, 出力サイト, systemd, `.notify_config` | なし | 診断結果は保存しない。 |
| `GET /api/rate-limit` | `.github_token` | なし | GitHub API 結果を返す。 |
| `GET /api/disk-usage` | `.build_logs/`, `.build_logs/archive/`, 出力サイト | なし | 集計のみ。 |
| `GET /api/webhook-events` | `.webhook_events.json` | なし | ページングして返す。 |
| `POST /api/webhook` | `.webhook_secret`, `.branch_config`, `.build_state`, `.maintenance`, `.build_circuit_state` | `.webhook_events.json`, `.build_state` または queue | 署名検証成功後のみイベント記録する。 |
| `GET /api/webhook-config` | `.webhook_secret` | なし | secret 本体は返さない。 |
| `POST /api/webhook-config` | なし | `.webhook_secret`, `.config_log` | secret 値は `.config_log` でマスクする。 |
| `POST /api/circuit-breaker/reset` | `.build_circuit_state` | `.build_circuit_state`, `.config_log` | 初期値へ戻す。冪等。 |
| `GET /api/snapshots` | `.snapshots/` | なし | 世代一覧を返す。 |
| `GET /api/snapshots/{id}/download` | `.snapshots/{id}/` | なし | バイナリを返す。 |
| `DELETE /api/snapshots/{id}` | `.snapshots/{id}/` | `.snapshots/`, `.config_log` | 対象 snapshot のみ削除する。 |
| `GET /api/maintenance` | `.maintenance` | なし | 不在時は disabled。 |
| `POST /api/maintenance/enable` | `.maintenance` | `.maintenance`, `.config_log` | `since` を現在時刻で保存する。 |
| `POST /api/maintenance/disable` | `.maintenance` | `.maintenance`, `.config_log` | disabled 状態を保存する。 |
| `GET /api/access-control` | `.access_control` | なし | 不在時は `allow: []`。 |
| `POST /api/access-control` | `.access_control` | `.access_control`, `.config_log` | allow 全体を置換する。 |
| `GET /api/hooks` | `.hooks` | なし | hook 一覧を返す。 |
| `POST /api/hooks` | `.hooks` | `.hooks`, `.config_log` | hook id を新規採番する。 |
| `DELETE /api/hooks/{id}` | `.hooks` | `.hooks`, `.config_log` | 対象 hook のみ削除する。 |
| `GET /api/hooks/{id}/log` | `.build_logs/{build_id}_hook_{id}.json` | なし | 直近 20 件を返す。 |
| `GET /api/alert-rules` | `.alert_rules` | なし | rule 一覧を返す。 |
| `POST /api/alert-rules` | `.alert_rules` | `.alert_rules`, `.config_log` | rule id を新規採番する。 |
| `DELETE /api/alert-rules/{id}` | `.alert_rules` | `.alert_rules`, `.config_log` | 対象 rule のみ削除する。 |
| `GET /api/tag-rules` | `.tag_rules` | なし | rule 一覧を返す。 |
| `POST /api/tag-rules` | `.tag_rules` | `.tag_rules`, `.config_log` | rule id を新規採番する。 |
| `DELETE /api/tag-rules/{id}` | `.tag_rules` | `.tag_rules`, `.config_log` | 対象 rule のみ削除する。 |
| `POST /api/verify-output` | `.build_history`, 出力サイト | なし | checksum 比較のみ。 |
| `GET /api/pipeline-config` | `.pipeline_config` | なし | 不在時は既定値。 |
| `POST /api/pipeline-config` | `.pipeline_config` | `.pipeline_config`, `.config_log` | config 全体を置換する。 |
| `GET /api/notes` | `.notes` | なし | 不在時は空文字。 |
| `POST /api/notes` | `.notes` | `.notes`, `.config_log` | content 全体を置換する。 |
| `GET /api/smtp-config` | `.smtp_config`, `.smtp_secret` | なし | password は返さず `password_set` だけ返す。 |
| `POST /api/smtp-config` | `.smtp_config`, `.smtp_secret` | `.smtp_config`, `.smtp_secret`, `.config_log` | password 指定時のみ `.smtp_secret` を更新する。 |
| `POST /api/smtp-test` | `.smtp_config`, `.smtp_secret` | `.notify_log` | 送信結果を記録する。 |
| `GET /api/queue` | `.build_state` | なし | queue 状態を返す。 |
| `DELETE /api/queue` | `.build_state` | `.build_state`, `.config_log` | 実行中 build は停止しない。 |
| `GET /api/dashboard-layout` | `.dashboard_layout` | なし | 不在時は既定 widget 順。 |
| `POST /api/dashboard-layout` | `.dashboard_layout` | `.dashboard_layout`, `.config_log` | widgets 全体を置換する。 |
| `GET /api/tokens` | `.api_tokens` | なし | token 本体は返さない。 |
| `POST /api/tokens` | `.api_tokens` | `.api_tokens`, `.access_log`, `.audit_log` | token 本体は作成時のみ返し、保存はハッシュのみ。 |
| `DELETE /api/tokens/{id}` | `.api_tokens` | `.api_tokens`, `.access_log`, `.audit_log` | 対象 token を失効する。 |

### 22.0e API 完全契約表

本表は API 実装、SDK 実装、標準管理ツール実装の契約インデックスである。実装者は endpoint を追加、削除、名称変更、body 変更、response 変更する前に本表を先に更新する。下表に存在しない endpoint は実装対象外とする。SHA reset 専用 endpoint とサマリー送信専用 endpoint は定義しない。

`Request` が `none` の場合、request body を受け付けない。空 JSON object `{}` も送信してはならない。`Response` は成功時 body の schema 名または最小 object を示す。詳細 schema は §22.0c、各 endpoint の個別例、§23 SDK 仕様、§24 UI 仕様を正とする。

`{message}` は `{"message": string}` を意味する。`{message,...}` 形式の response では `message` を必須キーとし、その他のキーも表記どおり必須とする。`?` が付いたキーだけを任意キーとする。成功時に空 body、`null` body、HTTP 204 は使用しない。

| Endpoint | Request | Response | Success | Errors | Read | Write | SDK | UI |
|----------|---------|----------|---------|--------|------|-------|-----|----|
| `POST /api/login` | `{password}` | `{token?,must_change,totp_required?,ticket?}` | `200` | `401`, `422`, `429`, `500` | `.admin_credentials`, `.totp_secret` | `.admin_credentials`, `.access_log`, `.audit_log` | `login(password)` | ログイン |
| `POST /api/login/totp` | `{ticket,code}` | `{token,must_change}` | `200` | `401`, `422`, `429`, `500` | `.admin_credentials`, `.totp_secret` | `.admin_credentials`, `.totp_secret`, `.access_log`, `.audit_log` | `loginTotp(ticket,code)` | ログイン |
| `POST /api/logout` | none | `{message}` | `200` | `401` | memory session | memory session | `logout()` | 全パネル共通 |
| `POST /api/change-password` | `{current_password,new_password}` | `{message}` | `200` | `401`, `422`, `500` | `.admin_credentials` | `.admin_credentials`, memory session | `changePassword()` | パスワード変更 |
| `GET /api/audit-log` | query `{limit,offset,actor?,action?,result?}` | `{log,total}` | `200` | `401`, `403`, `422`, `500` | `.audit_log` | none | `getAuditLog()` | 監査ログ |
| `GET /api/access-log` | query `{limit,offset}` | `{log}` | `200` | `401`, `422` | `.access_log` | none | `getAccessLog()` | アクセスログ |
| `GET /api/api-access-log` | query `{limit,offset,method?,path?,status?}` | `{log,total}` | `200` | `401`, `422`, `500` | `.api_access_log` | none | `getApiAccessLog()` | アクセスログ |
| `GET /api/sessions` | none | `{sessions}` | `200` | `401` | memory session | none | `getSessions()` | セッション管理 |
| `POST /api/sessions/revoke-all` | none | `{message,revoked_count}` | `200` | `401` | memory session | memory session, `.access_log` | `revokeAllSessions()` | セッション管理 |
| `GET /api/auth/totp-status` | none | `TotpStatus` | `200` | `401`, `500` | `.totp_secret` | none | `getTotpStatus()` | セキュリティ |
| `POST /api/auth/totp-setup` | none | `{secret,otpauth_uri}` | `200` | `401`, `403`, `409`, `500` | `.totp_secret` | none | `setupTotp()` | セキュリティ |
| `POST /api/auth/totp-confirm` | `{code}` | `TotpStatus` | `200` | `401`, `403`, `409`, `422`, `500` | `.totp_secret` | `.totp_secret`, `.audit_log` | `confirmTotp(code)` | セキュリティ |
| `DELETE /api/auth/totp` | `{code}` | `TotpStatus` | `200` | `401`, `403`, `409`, `422`, `500` | `.totp_secret` | `.totp_secret`, `.audit_log` | `disableTotp(code)` | セキュリティ |
| `GET /api/status` | none | `StatusObject` | `200` | `401`, `500` | `.build_status.json`, `.build_history`, `.build_state`, `.build_lock` | none | `getStatus()` | ステータス |
| `POST /api/build` | none | `{message,build_id?,queued?}` | `202` | `401`, `409`, `422`, `429`, `503` | `.server_config`, `.build_state`, `.build_lock`, `.maintenance`, `.build_circuit_state` | `.build_state` or queue | `triggerBuild()` | 手動実行 |
| `POST /api/build/force` | none | `{message,build_id?,queued?}` | `202` | `401`, `409`, `422`, `429`, `503` | `.server_config`, `.build_state`, `.build_lock`, `.maintenance`, `.build_circuit_state`, SHA cache | `.build_state`, SHA cache or queue | `buildForce()` | 手動実行 |
| `POST /api/build/cancel` | none | `{message}` | `200` | `401`, `404`, `409` | `.build_lock` | `.build_state`, `.build_logs/{id}.json` | `cancelBuild()` | 手動実行 |
| `GET /api/build/stream` | none | SSE `log/end` events | `200` | `401`, `404` | `.build_logs/{id}.json`, `.build_state` | none | `streamBuild()` | 手動実行 |
| `GET /api/logs` | query `{n,q}` | `{lines}` | `200` | `401`, `422`, `500` | `.build_logs/` | none | `getLogs()` | ログビューア |
| `GET /api/logs/search` | query `{q,from,to,level}` | `SearchResult` | `200` | `401`, `422` | `.build_logs/` | none | `searchLogs()` | ログビューア |
| `GET /api/logs/export` | none | `{exported_at,lines}` | `200` | `401` | `.build_logs/` | none | `exportLogs()` | ログビューア |
| `POST /api/logs/cleanup` | none | `{message,deleted_count}` | `200` | `401`, `500` | `.server_config`, `.build_logs/` | `.build_logs/` | `cleanupLogs()` | ログビューア, 設定 |
| `POST /api/logs/archive` | none | `{message,archived_count}` | `200` | `401`, `500` | `.server_config`, `.build_logs/` | `.build_logs/archive/`, `.build_logs/` | `archiveLogs()` | ログビューア, 設定 |
| `GET /api/history` | query `{page,per_page,trigger?,tag?,flagged?}` | `HistoryPageObject` | `200` | `401`, `422` | `.build_history` | none | `getHistory()` | ビルド履歴 |
| `GET /api/history/export` | none | `ExportObject` | `200` | `401` | `.build_history` | none | `exportHistory()` | ビルド履歴 |
| `GET /api/history/{id}/log` | path `{id}` | `HistoryLogObject` | `200` | `401`, `404`, `500` | `.build_logs/{id}.json` | none | `getHistoryLog(id)` | ビルド履歴 |
| `GET /api/history/{id}/comment` | path `{id}` | `CommentObject` | `200` | `401`, `404`, `500` | `.build_logs/{id}.json` | none | `getHistoryComment(id)` | ビルド履歴 |
| `POST /api/history/{id}/comment` | `{comment}` | `{message}` | `200` | `401`, `404`, `422`, `500` | `.build_logs/{id}.json` | `.build_logs/{id}.json`, `.config_log` | `setHistoryComment(id,comment)` | ビルド履歴 |
| `POST /api/history/{id}/flag` | `{flagged}` | `{message}` | `200` | `401`, `404`, `422`, `500` | `.build_logs/{id}.json` | `.build_logs/{id}.json`, `.config_log` | `setHistoryFlag(id,flagged)` | ビルド履歴 |
| `POST /api/history/{id}/tags` | `{tags}` | `{message}` | `200` | `401`, `404`, `422`, `500` | `.build_logs/{id}.json` | `.build_logs/{id}.json`, `.build_history`, `.config_log` | `setHistoryTags(id,tags)` | ビルド履歴 |
| `POST /api/history/{id}/rollback` | path `{id}` | `{message,build_id}` | `202` | `401`, `404`, `409`, `500` | `.snapshots/{id}/`, `.server_config` | `.build_history`, `.build_logs/{new_id}.json` | `rollbackHistory(id)` | ビルド履歴, スナップショット |
| `GET /api/sysinfo` | none | `SysinfoObject` | `200` | `401`, `500` | output file, process start time | none | `getSysinfo()` | システム情報 |
| `GET /api/health` | none | `HealthObject` | `200` | `500` | `.build_history`, `.pending_transfers` | none | `health()` | 死活監視 |
| `GET /api/schedule` | none | `ScheduleObject` | `200` | `401`, `500` | `.server_config`, systemd | none | `getSchedule()` | リポジトリ情報 |
| `POST /api/schedule/interval` | `{interval_seconds}` | `{message,interval_seconds}` | `200` | `401`, `409`, `422`, `500` | `.server_config` | `.server_config`, `.config_log`, systemd timer | `setScheduleInterval(seconds)` | リポジトリ情報 |
| `POST /api/schedule/pause` | none | `{message}` | `200` | `401`, `409`, `500` | `.server_config` | `.server_config`, `.config_log` | `pauseSchedule()` | リポジトリ情報 |
| `POST /api/schedule/resume` | none | `{message}` | `200` | `401`, `409`, `500` | `.server_config` | `.server_config`, `.config_log` | `resumeSchedule()` | リポジトリ情報 |
| `POST /api/schedule/allowed-hours` | `{from,to}` | `{message,allowed_hours}` | `200` | `401`, `422`, `500` | `.server_config` | `.server_config`, `.config_log` | `setAllowedHours()`, `clearAllowedHours()` | リポジトリ情報 |
| `POST /api/schedule/force-interval` | `{hours}` | `{message,hours}` | `200` | `401`, `422`, `500` | `.server_config` | `.server_config`, `.config_log` | `setForceInterval(hours)` | リポジトリ情報 |
| `POST /api/schedule/cooldown` | `{seconds}` | `{message,seconds}` | `200` | `401`, `422`, `500` | `.server_config` | `.server_config`, `.config_log` | `setBuildCooldown(seconds)` | リポジトリ情報 |
| `GET /api/notify-config` | none | `NotifyConfig` | `200` | `401`, `500` | `.notify_config`, `.smtp_config` | none | `getNotifyConfig()` | 通知設定 |
| `POST /api/notify-config` | `NotifyConfig` | `{message}` | `200` | `401`, `422`, `500` | `.notify_config` | `.notify_config`, `.config_log` | `setNotifyConfig(config)` | 通知設定 |
| `GET /api/notify-log` | query `{limit,offset}` | `{log}` | `200` | `401`, `422` | `.notify_log` | none | `getNotifyLog()` | 通知設定 |
| `POST /api/notify-test` | none | `{message,webhook_url}` | `200` | `401`, `422`, `500` | `.notify_config` | `.notify_log` | `notifyTest()` | 通知設定 |
| `POST /api/notify/weekly-summary` | none | `{message,period,success_count,failure_count,success_rate}` | `200` | `401`, `422`, `500` | `.notify_config`, `.build_logs/` | `.notify_log` | `notifyWeeklySummary()` | 通知設定 |
| `GET /api/config` | none | `ConfigObject` | `200` | `401`, `500` | `.server_config` | none | `getConfig()` | 設定 |
| `POST /api/config/validate` | partial `ConfigObject` | `ConfigValidationObject` | `200` | `401`, `422`, `500` | `.server_config` | none | `validateConfig(config)` | 設定 |
| `POST /api/config` | partial `ConfigObject` | `{message,config}` | `200` | `401`, `422`, `500` | `.server_config` | `.server_config`, `.config_log` | `setConfig(config)` | 設定 |
| `POST /api/log-level` | `{level}` | `{message,level}` | `200` | `401`, `422`, `500` | `.server_config` | `.server_config`, `.config_log` | `setLogLevel(level)` | 設定 |
| `GET /api/config-log` | query `{limit,offset}` | `{log}` | `200` | `401`, `422` | `.config_log` | none | `getConfigLog()` | 設定 |
| `GET /api/api-rate-limit` | none | `ApiRateLimitPolicy` | `200` | `401`, `403`, `500` | `.server_config`, `.api_rate_state` | none | `getApiRateLimit()` | セキュリティ |
| `POST /api/api-rate-limit` | `ApiRateLimitPolicy` | `ApiRateLimitPolicy` | `200` | `401`, `403`, `422`, `500` | `.server_config` | `.server_config`, `.api_rate_state`, `.config_log`, `.audit_log` | `setApiRateLimit(policy)` | セキュリティ |
| `GET /api/pat-status` | none | `PatStatusObject` | `200` | `401`, `501`, `500` | `.github_token` | none | `getPatStatus()` | システム情報 |
| `POST /api/pat-verify` | none | `PatVerifyObject` | `200` | `401`, `501`, `500` | `.github_token` | none | `patVerify()` | システム情報 |
| `POST /api/pat-update` | `{token}` | `{message}` | `200` | `401`, `422`, `500` | none | `.github_token`, `.config_log` | `updatePat(token)` | システム情報 |
| `GET /api/stats` | query `{days}` | `StatsObject` | `200` | `401`, `422` | `.build_history`, `.build_logs/` | none | `getStats(days)` | 統計 |
| `GET /api/stats/timeline` | query `{days}` | `TimelineObject` | `200` | `401`, `422` | `.build_history` | none | `getStatsTimeline(days)` | 統計 |
| `GET /api/stats/build-duration` | query `{n}` | `BuildDurationStats` | `200` | `401`, `422` | `.build_logs/` | none | `getStatsBuildDuration(n)` | 統計 |
| `GET /api/output-meta` | none | `OutputMetaObject` | `200` | `401`, `404`, `500` | `.build_history`, `.build_logs/`, `.build_logs/archive/`, output file | none | `getOutputMeta()` | システム情報 |
| `GET /api/repo-info` | none | `RepoInfoObject` | `200` | `401`, `500` | `.repo_config` | none | `getRepoInfo()` | リポジトリ情報 |
| `POST /api/repo-config` | partial `RepoInfoObject` | `{message}` | `200` | `401`, `422`, `500` | `.repo_config` | `.repo_config`, `.config_log` | `setRepoConfig(config)` | リポジトリ情報 |
| `GET /api/branch-config` | none | `{source,branches}` | `200` | `401`, `500` | `.branch_config` | none | `getBranchConfig()` | リポジトリ情報 |
| `POST /api/branch-config` | `{branches}` | `{message,branches_count}` | `200` | `401`, `422`, `500` | `.branch_config` | `.branch_config`, `.config_log` | `setBranchConfig(branches)` | リポジトリ情報 |
| `GET /api/backup` | none | `BackupObject` | `200` | `401`, `500` | config state files | none | `backup()` | 設定 |
| `POST /api/restore` | `BackupObject` | `{message}` | `200` | `401`, `422`, `500` | request body | config state files, `.config_log` | `restore(config)` | 設定 |
| `GET /api/dashboard` | none | `DashboardObject` | `200` | `401`, `500` | `.build_status.json`, `.build_history`, `.build_state`, `.build_lock`, `.server_config`, `.alert_rules`, `.dashboard_layout`, output file, process start time | none | `getDashboard()` | ステータス, システム診断 |
| `GET /api/diagnostics` | none | `DiagnosticsObject` | `200` | `401`, `500` | `.github_token`, output file, systemd, `.notify_config` | none | `getDiagnostics()` | システム診断 |
| `GET /api/rate-limit` | none | `RateLimitObject` | `200` | `401`, `501`, `500` | `.github_token` | none | `getRateLimit()` | システム情報 |
| `GET /api/disk-usage` | none | `DiskUsageObject` | `200` | `401`, `500` | `.build_logs/`, `.build_logs/archive/`, output file | none | `getDiskUsage()` | システム情報 |
| `GET /api/webhook-events` | query `{limit,offset}` | `{events,total}` | `200` | `401`, `422`, `500` | `.webhook_events.json` | none | `getWebhookEvents(limit,offset)` | システム診断 |
| `POST /api/webhook` | GitHub webhook body | `{message,ref?}` | `200` | `400`, `403`, `409`, `501`, `503` | `.webhook_secret`, `.branch_config`, `.build_state`, `.maintenance`, `.build_circuit_state` | `.webhook_events.json`, `.build_state` or queue | none | 外部 Webhook |
| `GET /api/webhook-config` | none | `{configured}` | `200` | `401`, `500` | `.webhook_secret` | none | `getWebhookConfig()` | 通知設定 |
| `POST /api/webhook-config` | `{secret}` | `{message}` | `200` | `401`, `422`, `500` | none | `.webhook_secret`, `.config_log` | `setWebhookConfig(secret)` | 通知設定 |
| `POST /api/circuit-breaker/reset` | none | `{message,open,consecutive_failures}` | `200` | `401`, `500` | `.build_circuit_state` | `.build_circuit_state`, `.config_log` | `resetCircuitBreaker()` | 手動実行, システム診断 |
| `GET /api/snapshots` | none | `{snapshots}` | `200` | `401`, `500` | `.snapshots/` | none | `getSnapshots()` | スナップショット |
| `GET /api/snapshots/{id}/download` | path `{id}` | binary | `200` | `401`, `404`, `500` | `.snapshots/{id}/` | none | `downloadSnapshot(id)` | スナップショット |
| `DELETE /api/snapshots/{id}` | path `{id}` | `{message}` | `200` | `401`, `404`, `500` | `.snapshots/{id}/` | `.snapshots/`, `.config_log` | `deleteSnapshot(id)` | スナップショット |
| `GET /api/maintenance` | none | `MaintenanceObject` | `200` | `401`, `500` | `.maintenance` | none | `getMaintenance()` | メンテナンス |
| `POST /api/maintenance/enable` | `{reason}` | `{message,since}` | `200` | `401`, `422`, `500` | `.maintenance` | `.maintenance`, `.config_log` | `enableMaintenance(reason)` | メンテナンス |
| `POST /api/maintenance/disable` | none | `{message}` | `200` | `401`, `500` | `.maintenance` | `.maintenance`, `.config_log` | `disableMaintenance()` | メンテナンス |
| `GET /api/access-control` | none | `{allow}` | `200` | `401`, `500` | `.access_control` | none | `getAccessControl()` | アクセス制御 |
| `POST /api/access-control` | `{allow}` | `{message,allow}` | `200` | `401`, `422`, `500` | `.access_control` | `.access_control`, `.config_log` | `setAccessControl(allowList)` | アクセス制御 |
| `GET /api/hooks` | none | `{hooks}` | `200` | `401`, `500` | `.hooks` | none | `getHooks()` | フック |
| `POST /api/hooks` | `{phase,command_args,abort_on_failure,timeout_seconds?}` | `HookRecord` | `201` | `401`, `409`, `422`, `500` | `.hooks` | `.hooks`, `.config_log` | `addHook()` | フック |
| `DELETE /api/hooks/{id}` | path `{id}` | `{message}` | `200` | `401`, `404`, `500` | `.hooks` | `.hooks`, `.config_log` | `deleteHook(id)` | フック |
| `GET /api/hooks/{id}/log` | path `{id}` | `{id,runs}` | `200` | `401`, `404`, `500` | `.build_logs/{build_id}_hook_{id}.json` | none | `getHookLog(id)` | フック |
| `GET /api/alert-rules` | none | `{rules}` | `200` | `401`, `500` | `.alert_rules` | none | `getAlertRules()` | 設定 |
| `POST /api/alert-rules` | `{metric,operator,threshold,level,message}` | `AlertRule` | `201` | `401`, `409`, `422`, `500` | `.alert_rules` | `.alert_rules`, `.config_log` | `addAlertRule()` | 設定 |
| `DELETE /api/alert-rules/{id}` | path `{id}` | `{message}` | `200` | `401`, `404`, `500` | `.alert_rules` | `.alert_rules`, `.config_log` | `deleteAlertRule(id)` | 設定 |
| `GET /api/tag-rules` | none | `{rules}` | `200` | `401`, `500` | `.tag_rules` | none | `getTagRules()` | 設定 |
| `POST /api/tag-rules` | `{condition,tags}` | `TagRule` | `201` | `401`, `409`, `422`, `500` | `.tag_rules` | `.tag_rules`, `.config_log` | `addTagRule()` | 設定 |
| `DELETE /api/tag-rules/{id}` | path `{id}` | `{message}` | `200` | `401`, `404`, `500` | `.tag_rules` | `.tag_rules`, `.config_log` | `deleteTagRule(id)` | 設定 |
| `POST /api/verify-output` | none | `{match,expected,actual}` | `200` | `401`, `404`, `500` | `.build_history`, output file | none | `verifyOutput()` | システム診断 |
| `GET /api/pipeline-config` | none | `PipelineConfig` | `200` | `401`, `500` | `.pipeline_config` | none | `getPipelineConfig()` | 設定 |
| `POST /api/pipeline-config` | `PipelineConfig` | `{message}` | `200` | `401`, `422`, `500` | `.pipeline_config` | `.pipeline_config`, `.config_log` | `setPipelineConfig(config)` | 設定 |
| `GET /api/notes` | none | `{content,updated_at}` | `200` | `401`, `500` | `.notes` | none | `getNotes()` | 運用ノート |
| `POST /api/notes` | `{content}` | `{message,updated_at}` | `200` | `401`, `422`, `500` | `.notes` | `.notes`, `.config_log` | `setNotes(content)` | 運用ノート |
| `GET /api/smtp-config` | none | `SmtpConfig` | `200` | `401`, `500` | `.smtp_config`, `.smtp_secret` | none | `getSmtpConfig()` | 通知設定 |
| `POST /api/smtp-config` | partial `SmtpConfig` with optional password | `{message}` | `200` | `401`, `422`, `500` | `.smtp_config`, `.smtp_secret` | `.smtp_config`, `.smtp_secret`, `.config_log` | `setSmtpConfig(config)` | 通知設定 |
| `POST /api/smtp-test` | none | `{result,message}` | `200` | `401`, `422`, `500` | `.smtp_config`, `.smtp_secret` | `.notify_log` | `smtpTest()` | 通知設定 |
| `GET /api/queue` | none | `{queued,max_size}` | `200` | `401`, `500` | `.build_state` | none | `getQueue()` | 手動実行 |
| `DELETE /api/queue` | none | `{message,cleared_count}` | `200` | `401`, `500` | `.build_state` | `.build_state`, `.config_log` | `clearQueue()` | 手動実行 |
| `GET /api/dashboard-layout` | none | `{widgets}` | `200` | `401`, `500` | `.dashboard_layout` | none | `getDashboardLayout()` | ステータス |
| `POST /api/dashboard-layout` | `{widgets}` | `{message}` | `200` | `401`, `422`, `500` | `.dashboard_layout` | `.dashboard_layout`, `.config_log` | `setDashboardLayout(widgets)` | ステータス |
| `GET /api/tokens` | none | `{tokens}` | `200` | `401`, `403`, `500` | `.api_tokens` | none | `getTokens()` | API トークン管理 |
| `POST /api/tokens` | `{label,scopes,expires_at?}` | `TokenCreateResult` | `201` | `401`, `403`, `422`, `500` | `.api_tokens` | `.api_tokens`, `.access_log`, `.audit_log` | `createToken(label,scopes,expiresAt)` | API トークン管理 |
| `DELETE /api/tokens/{id}` | path `{id}` | `{message}` | `200` | `401`, `403`, `404`, `500` | `.api_tokens` | `.api_tokens`, `.access_log`, `.audit_log` | `revokeToken(id)` | API トークン管理 |

**成果物 / archive / backup API 副作用固定契約：**

| API | 処理順序 | 成功時副作用 | 失敗時副作用 |
|-----|----------|--------------|--------------|
| `POST /api/logs/cleanup` | `.server_config` 読込 → cleanup 対象 log 算出 → 対象 file 削除 → response。 | 対象 `.build_logs/{id}.json` だけを削除する。archive 済み `.json.gz` は削除しない。 | 算出前失敗は差分なし。途中削除失敗は `500` とし、削除済み file は戻さない。 |
| `POST /api/logs/archive` | `.server_config` 読込 → 対象 log 算出 → `.build_logs/archive/{id}.json.gz.tmp.{pid}` 作成 → gzip 書込 → fsync → rename → 元 log 削除 → response。 | archive 成功した log だけ元 `.json` を削除する。gzip は 1 log 1 file。 | gzip 作成または rename 失敗時は元 log を残す。元 log 削除失敗は `500` とし、archive 済み `.json.gz` は残す。 |
| `GET /api/backup` | 対象設定 file 読込 → 不在 file に既定値適用 → secret mask → response。 | 状態ファイルを更新しない。 | 読込不能な必須 file は `500`。任意 file 不在は既定値で返す。 |
| `POST /api/restore` | request 検証 → secret mask `"***"` の既存値補完 → 全対象 payload 生成 → §22.0d の順に atomic write → `.config_log` 追記 → response。 | 設定系状態 file だけを置換する。履歴、ログ、snapshot、session、token 本体は復元しない。 | 検証失敗は差分なし。途中 write 失敗は未処理 file を書かず `500`。処理済み file は戻さない。 |
| `GET /api/snapshots/{id}/download` | id 検証 → snapshot directory 検証 → tar.gz stream 生成 → response。 | 状態ファイルを更新しない。 | 不正 id は `422`、不在は `404`、stream 中の読込失敗は接続を終了し状態差分なし。 |
| `POST /api/history/{id}/rollback` | id 検証 → running 確認 → snapshot 検証 → 新 build id 採番 → deploy 転送 → rollback log/history 保存 → response。 | 新規 rollback build log/history だけを追加する。元 snapshot、元 history、state 全体は巻き戻さない。 | 転送失敗は rollback build log/history を failure として保存し、元 snapshot は削除しない。running 中は差分なし `409`。 |

archive 対象 id と snapshot id は build id 形式だけを許可する。API は request path の URL decode 後に `/`、`\`、`..`、空文字、NUL を含む id を `422` とする。tar.gz へ格納する path は snapshot directory からの相対 path とし、絶対 path、`..`、symlink entry、hardlink entry、device entry を含めてはならない。

backup response に secret 原文を含めてはならない。`password`、`token`、`secret`、`smtp_password`、`webhook_secret`、`.github_token`、`.smtp_secret`、`.api_tokens` の hash 元値は `"***"` または `*_set:boolean` で表現する。`POST /api/restore` で `"***"` を受け取った secret は既存値保持を意味し、既存値がない場合は未設定として扱う。

**backup / restore 固定契約：**

| 項目 | 仕様 |
|------|------|
| backup 対象 | `.server_config`、`.notify_config`、`.repo_config`、`.branch_config`、`.access_control`、`.hooks`、`.alert_rules`、`.tag_rules`、`.pipeline_config`、`.dashboard_layout`、`.smtp_config`、`.webhook_secret`、`.smtp_secret`。 |
| backup 対象外 | `.admin_credentials`、`.api_tokens`、`.totp_secret`、session、`.build_history`、`.build_logs/`、`.snapshots/`、`.notify_log`、`.notify_pending`、`.webhook_events.json`、`.approval_queue`、`.build_state`。 |
| backup 不在値 | 任意設定 file 不在は schema 既定値で返す。secret file 不在は `*_set:false`。 |
| backup secret | `.webhook_secret` と `.smtp_secret` は本体を返さず、`webhook_secret_set` / `smtp_password_set` boolean だけ返す。 |
| restore 検証 | すべての対象 payload を先に schema 検証し、1 件でも不正なら書込を開始せず `422`。 |
| restore `"***"` | 対応する既存 secret がある場合だけ既存値保持。既存 secret がない場合は未設定として扱い、新規 secret 文字列として保存しない。 |
| restore secret 削除 | secret key が `null` の場合は削除。key 省略は既存保持。 |
| restore 書込順 | `.server_config` → `.notify_config` → `.repo_config` → `.branch_config` → `.access_control` → `.hooks` → `.alert_rules` → `.tag_rules` → `.pipeline_config` → `.dashboard_layout` → `.smtp_config` → `.webhook_secret` → `.smtp_secret` → `.config_log`。 |
| restore 途中失敗 | 未処理 file は書かない。処理済み file は巻き戻さない。response は `500`。 |
| restore no-op | 全対象が既存値と同一の場合は file と `.config_log` を変更せず `{ "message":"No changes" }`。 |

restore の `.config_log` は対象 file ごとの差分を 1 record にまとめ、secret はすべて `"***"` とする。backup / restore の response、server log、fixture expected に secret 平文を含めてはならない。

**backup / restore fixture 固定：**

| fixture | 入力 | 期待結果 |
|---------|------|----------|
| backup-mask | secret 設定済みで backup | secret 本体なし、`*_set:true`。 |
| restore-validate-fail | 1 file schema 不正 | `422`、全 file 差分なし。 |
| restore-secret-keep | `"***"` かつ既存 secret あり | 既存 secret 維持、平文出力なし。 |
| restore-secret-missing | `"***"` かつ既存 secret なし | secret 未設定のまま、`"***"` を保存しない。 |
| restore-secret-delete | secret `null` | secret file 削除。 |
| restore-write-failure | 中途 write 失敗 | 未処理 file は差分なし、処理済み file は維持、`500`。 |

**ビルド操作の競合優先順位：**

`POST /api/build`、`POST /api/build/force`、`POST /api/webhook`、`POST /api/history/{id}/rollback` は、以下の順に判定する。

1. 認証・権限を確認する。Webhook は署名検証を認証の代替とする。
2. `.maintenance.enabled == true` の場合は `503 {"error":"maintenance"}` を返す。
3. `.build_circuit_state.open == true` の場合は `409 {"error":"circuit_open"}` を返す。ただし `POST /api/circuit-breaker/reset` は対象外。
4. `.build_state.running == true` の場合:
   - queue 対応 API（`POST /api/build`、`POST /api/build/force`、署名検証済み `POST /api/webhook`）は `queue_max_size > 0` かつ空きがあれば queue に追加し、`202 {"message":"Build queued","queued":true}` を返す。
   - queue が無効または満杯の場合は `429 {"error":"queue_full"}` を返す。
   - rollback は queue へ積まず、`409 {"error":"Build is running"}` を返す。
5. 実行中でなければ `.build_state.running=true` と新しい `build_id` を保存し、`202 {"message":"Build started","build_id":"...","queued":false}` を返す。

`POST /api/build/cancel` は queue を削除しない。実行中 build のみを cancel 対象とし、実行中でない場合は `409 {"error":"No build is running"}` を返す。`DELETE /api/queue` は実行中 build を停止せず、待機 queue のみ削除する。

### 22.0e.1 API 機能別処理契約

集約 API は、下表の読取元、算出方法、空状態の戻り値に従う。下表にない読取元や推測値を使用してレスポンスを補完してはならない。

| 機能 | Endpoint | 読取元 | 算出方法 | 空状態 / 不足時 |
|------|----------|--------|----------|-----------------|
| 現在状態 | `GET /api/status` | `.build_status.json`, `.build_history`, `.build_state`, `.build_lock` | `.build_status.json` から `last_sha`、`last_build_at`、`last_build_status`、`last_trigger`、pending 件数、circuit 状態を返す。不在時のみ `.build_history` の最新行と `.build_state` から算出する。`last_sha` は `.build_status.json.last_blob_sha` があればその値、なければ `last_commit_sha`、両方なければ `null` とする。 | 履歴なしは `last_sha:null`, `last_build_at:null`, `last_build_status:"none"`, `last_trigger:null`, `output_url:null`。 |
| 手動ビルド開始 | `POST /api/build` | `.server_config`, `.build_state`, `.build_lock`, `.maintenance`, `.build_circuit_state` | §22.0e の競合優先順位に従い、開始または queue 追加を行う。SHA cache は変更しない。 | queue 無効または満杯は `429 {"error":"queue_full"}`。 |
| 強制ビルド開始 | `POST /api/build/force` | `.server_config`, `.build_state`, `.build_lock`, `.maintenance`, `.build_circuit_state`, SHA cache | 開始可能な場合のみ SHA cache を空 SHA に更新し、同一状態更新内で build を開始する。queue 追加時は queue entry の `payload.force=true` を保存する。 | queue 無効または満杯は `429 {"error":"queue_full"}`。 |
| ログ一覧 | `GET /api/logs` | `.build_logs/` | 最新 build log の `stdout`、`stderr`、`warnings` を時系列順に連結し、`n` 件に丸める。`q` が空でない場合は部分一致行だけを返す。 | ログなしは `{"lines":[]}`。 |
| ログ検索 | `GET /api/logs/search` | `.build_logs/` | 全 build log を新しい順に読み、`q`、`from`、`to`、`level` で絞り込む。`level` は行内の `[INFO]`、`[WARNING]`、`[ERROR]`、`[DEBUG]` に一致させる。 | 一致なしは `results:[]`。 |
| 履歴一覧 | `GET /api/history` | `.build_history` | JSON Lines を新しい順で読み、`trigger`、`tag`、`flagged` を指定時のみ完全一致で絞り込み、ページングする。壊れた行は無視して ERROR ログに記録する。`total` と `pages` は絞り込み後の有効行だけで算出する。 | 履歴なしまたは一致なしは `total:0`, `pages:0`, `history:[]`。 |
| 出力メタ | `GET /api/output-meta` | `.build_history`, `.build_logs/`, `.build_logs/archive/`, 出力サイト | 出力サイトの現在サイズと mtime、直近成功履歴の `output_sha256`、直近ログの `report`、直近ログまたは HTML meta の `build_id` / `commit_sha` / `build_at` を返す。通常ログを先に読み、不在時だけ archive を読む。 | 出力サイト不在は `404`。report 不在の数値は `null`、warning は `[]`。build meta 不在は空文字。 |
| ダッシュボード | `GET /api/dashboard` | `.build_status.json`, `.build_history`, `.server_config`, `.alert_rules`, `.dashboard_layout`, `.build_state`, `.build_lock`, 出力サイト, process start time | `status`、`sysinfo`、`stats(days=7)`、`schedule`、`alerts` を同一リクエスト時点で算出し、widget 順序は `.dashboard_layout.widgets` を使用する。 | `.dashboard_layout` 不在は既定 widget 順。alerts なしは `[]`。 |
| 診断 | `GET /api/diagnostics` | `.github_token`, 出力サイト, systemd, `.notify_config`, `.webhook_secret` | PAT、GitHub API、出力サイト、systemd、Webhook 設定を個別 item として返す。診断結果は保存しない。 | 各項目は `ok`、`warn`、`error` のいずれかを返す。 |
| キュー | `GET /api/queue` | `.build_state`, `.server_config` | `.build_state.queued` と `.server_config.queue_max_size` を返す。 | `.build_state` 不在は初期値で `queued:[]`。 |
| バックアップ | `GET /api/backup` | §22.0d の backup 対象状態ファイル | 設定状態だけを export し、secret 値は `"***"` または boolean にマスクする。履歴、ログ、snapshot、session は含めない。 | 不在の任意設定ファイルは初期値で返す。 |
| リストア | `POST /api/restore` | request body | 対象 state schema をすべて検証してから §22.0d の write 順に保存する。secret が `"***"` の場合は既存 secret を保持する。 | 検証失敗は書き込み前に `422`。途中失敗は未処理ファイルを書かない。 |

### 22.0e.2 API ID 採番契約

API が新規 ID を生成する機能は、下表の形式に従う。既存 ID と衝突した場合は、同一時刻内で末尾に 3 桁の連番を付け、最大 999 まで試行する。999 回衝突した場合は `500 Internal Server Error` を返す。

| 対象 | 形式 | 例 | 衝突時 |
|------|------|----|--------|
| build id | `b{YYYYMMDDHHmmss}` | `b20260915100500` | `b20260915100500-001` |
| queue id | `q{YYYYMMDDHHmmss}` | `q20260915100500` | `q20260915100500-001` |
| token id | `tok` + 6 桁連番 | `tok000001` | 既存最大番号 + 1。999999 超過時は `500` |
| hook id | `h{YYYYMMDDHHmmss}` | `h20260915100500` | `h20260915100500-001` |
| approval id | `appr{YYYYMMDDHHmmss}` | `appr20260915100500` | `appr20260915100500-001` |
| alert rule id | `r{YYYYMMDDHHmmss}` | `r20260915100500` | `r20260915100500-001` |
| tag rule id | `t{YYYYMMDDHHmmss}` | `t20260915100500` | `t20260915100500-001` |
| snapshot id | `snap{YYYYMMDDHHmmss}` | `snap20260915100500` | `snap20260915100500-001` |

### 22.0e.3 API Endpoint 種別別実装契約

API handler は endpoint ごとの個別処理へ入る前に、§22.0 の判定順を必ず適用する。共通判定後の endpoint 種別別処理は下表に従う。

| 種別 | 対象 endpoint | 処理順序 | 成功時 | 失敗時 |
|------|---------------|----------|--------|--------|
| read-only list | `GET /api/history`, `GET /api/notify-log`, `GET /api/config-log`, `GET /api/webhook-events` | query 検証 → 対象ファイル読込 → 壊れた行を除外 → sort / paging → response 生成 | `total` がある endpoint は除外後件数を返す。 | query 不正は `422`。ファイル読込不能は `500`。 |
| read-only aggregate | `GET /api/status`, `GET /api/dashboard`, `GET /api/health`, `GET /api/output-meta` | 必要ファイルを read-only で読込 → 不在時初期値適用 → 算出 → response 生成 | GET の副作用なし。`.build_status.json` 破損時も API が自動修復しない。 | 必須ファイル破損は endpoint 固有の `500`。任意ファイル不在は初期値。 |
| single-file update | `POST /api/config`, `POST /api/repo-config`, `POST /api/maintenance/*`, `POST /api/access-control`, `POST /api/dashboard-layout` | body 検証 → 現在値読込 → 差分生成 → atomic write → `.config_log` 追記 | 更新後値または `{message}` を返す。 | body 検証失敗は書込前に `422`。`.config_log` 失敗時は対象更新済みのまま `500`。 |
| multi-file update | `POST /api/restore`, `POST /api/smtp-config`, `POST /api/history/{id}/tags` | 全入力検証 → 全対象読込 → 書込計画生成 → §22.0d の Write 順に atomic write | 全対象の更新完了後に response を返す。 | 検証失敗は書込なし `422`。途中失敗は未処理ファイルを書かず `500`。 |
| secret update | `POST /api/pat-update`, `POST /api/webhook-config`, `POST /api/smtp-config` password あり | secret 入力検証 → secret ファイル atomic write mode `0600` → `.config_log` へ `"***"` で記録 | secret 本体を response に含めない。 | secret 書込失敗は `500`。ログ、response、stdout へ平文を出さない。 |
| build command | `POST /api/build`, `POST /api/build/force`, `POST /api/history/{id}/rollback` | 認証 → maintenance → circuit → running / queue 判定 → `.build_state` 更新 | `202` と開始または queue 結果を返す。 | running 競合は `409` または `429`。状態更新失敗は `500`。 |
| destructive delete | `DELETE /api/queue`, `DELETE /api/snapshots/{id}`, `DELETE /api/hooks/{id}`, `DELETE /api/alert-rules/{id}`, `DELETE /api/tag-rules/{id}`, `DELETE /api/tokens/{id}` | path / auth 検証 → 対象存在確認 → 削除または失効 → 対象 endpoint の契約に従い `.config_log` または `.audit_log` 追記 | `{message}` と件数がある場合は件数を返す。 | 対象不在は `404`。部分削除は禁止し、失敗時は `500`。 |
| external check | `POST /api/pat-verify`, `GET /api/rate-limit`, `GET /api/diagnostics`, `POST /api/smtp-test`, `POST /api/notify-test` | 設定読込 → timeout 付き外部確認 → 結果 response → 必要時 log 追記 | 確認結果を保存しない。ただし test 送信 log は仕様どおり追記する。 | 未設定は `501` または endpoint 固有 `422`。timeout は `500`。 |
| binary response | `GET /api/snapshots/{id}/download` | path 検証 → snapshot 存在確認 → archive stream | `Content-Type` と `Content-Disposition` を付与する。 | 不在は `404`。読込失敗は `500`。 |
| stream response | `GET /api/build/stream` | 認証 → 最新 / 実行中 log 特定 → SSE header → frame 送信 | `log` frame 後、必ず `end` frame を送って close する。 | log 不在は `404`。送信中断時は状態ファイルを更新しない。 |

`.config_log`、`.access_log`、`.notify_log` への追記は JSON Lines 1 行単位で行う。追記失敗時は対象 endpoint の副作用が既に完了している場合でも、失敗を `500` として返し、次回 GET で破損行を無視できる形式を維持する。追記行の末尾改行を書けなかった場合は、その行を破損行として扱う。

**JSON Lines 読取・破損行契約：**

| ファイル | GET / 集計時 | 追記失敗時 | server log 固定コード |
|----------|--------------|------------|-----------------------|
| `.build_history` | 壊れた行を無視し、有効行だけで算出する。 | runner build 処理は失敗扱い。 | `BUILD_HISTORY_SKIP_CORRUPT` |
| `.config_log` | 壊れた行を無視する。 | 設定変更 API は `500`。対象状態ファイル更新済みの場合は戻さない。 | `CONFIG_LOG_WRITE_FAILED` |
| `.access_log` | 壊れた行を無視する。 | auth / token 操作は `500`。logout だけは token 破棄を優先する。 | `ACCESS_LOG_WRITE_FAILED` |
| `.api_access_log` | 壊れた行を無視する。 | 本来の API response を優先し、`500` へ変更しない。 | `API_ACCESS_LOG_WRITE_FAILED` |
| `.audit_log` | 壊れた行を無視する。 | 監査対象操作は `500`。既に状態更新済みの場合は戻さない。 | `AUDIT_LOG_WRITE_FAILED` |
| `.notify_log` | 壊れた行を無視する。 | 通知送信結果だけ失敗扱いにし、build 成否は反転しない。 | `NOTIFY_LOG_WRITE_FAILED` |
| `.webhook_events.json` | 壊れた行を無視する。 | queue 追加前なら `500`、queue 追加後なら response に `event_log_failed:true` を含める。 | `WEBHOOK_EVENT_LOG_WRITE_FAILED` |
| `.approval_queue` | 壊れた行を無視する。 | approval entry 作成は失敗扱い。 | `APPROVAL_QUEUE_WRITE_FAILED` |

JSON Lines の壊れた行は、空行、JSON parse 失敗、JSON object 以外、必須 key 不足、型不一致のいずれかとする。壊れた行を response に含めてはならない。壊れた行を検出しても GET API が対象ファイルを自動修復、削除、上書きしてはならない。

### 22.0e.4 API レスポンス正規化契約

API response は、§22.0e の Response 列、§22.0c の schema、§23 の SDK 型定義表に一致させる。実装者は endpoint ごとに以下の正規化を行う。

| 対象 | 仕様 |
|------|------|
| object response | 必須 key をすべて含める。値が存在しない場合は、schema で nullable の key だけ `null` を使用する。 |
| array response | 配列 key は未取得、対象なし、空状態のいずれでも `[]` を返す。`null`、key 省略は禁止する。 |
| boolean | `true` / `false` だけを返す。`"true"`、`1`、`0` は使用しない。 |
| integer | JSON number の整数として返す。文字列化しない。 |
| float | JSON number とし、小数第 2 位までに丸める指定がある値だけ `math.Round(x*100)/100` 相当にする。 |
| timestamp | UTC ISO 8601 `YYYY-MM-DDTHH:MM:SSZ`。空状態は nullable key なら `null`、非 nullable key なら endpoint 固有の既定値。 |
| message | 成功 message は endpoint ごとの固定文言とし、入力値を連結しない。 |
| unknown key | response に schema 外 key を追加しない。互換目的の旧 key 追加も禁止する。 |

**ページング response：**

| Endpoint | 入力 | 出力 | 算出 |
|----------|------|------|------|
| `GET /api/history` | `page`, `per_page` | `total`, `page`, `per_page`, `pages`, `history` | `pages = ceil(total / per_page)`。`total=0` の場合 `pages=0`。 |
| `GET /api/webhook-events` | `limit`, `offset` | `events`, `total` | `total` は filter 後、limit/offset 前の件数。 |
| `GET /api/api-access-log` | `limit`, `offset` | `log`, `total` | `total` は filter 後、limit/offset 前の件数。 |
| `GET /api/audit-log` | `limit`, `offset` | `log`, `total` | `total` は filter 後、limit/offset 前の件数。 |

`offset >= total` の場合は空配列を返し、`404` にしない。`page > pages` の場合は `history:[]` を返し、`page` は request 値を保持する。`page < 1`、`per_page < 1`、`limit < 1`、`offset < 0` は `422`。

**部分更新 response：**

| Endpoint | 成功 response | 補足 |
|----------|---------------|------|
| `POST /api/config` | `{message:"Config updated",config}` または `{message:"No changes",config}` | `config` は更新後に既定値 merge 済みの `ConfigObject`。 |
| `POST /api/repo-config` | `{message:"Repo config updated"}` または `{message:"No changes"}` | 未指定 key は保持する。 |
| `POST /api/branch-config` | `{message:"Branch config updated",branches_count}` | 空配列で default 復帰した場合 `branches_count=0`。 |
| `POST /api/notify-config` | `{message:"Notify config updated"}` | secret は返さない。 |
| `POST /api/smtp-config` | `{message:"SMTP config updated"}` | password は返さない。 |
| `POST /api/dashboard-layout` | `{message:"Dashboard layout updated"}` | 保存後 widgets は再取得で確認する。 |

**no-op / 部分更新副作用契約：**

| Endpoint | no-op 判定 | no-op 時の副作用 | 更新時の副作用 |
|----------|------------|------------------|----------------|
| `POST /api/config` | 正規化後 config が既存値と一致 | `.server_config`、`.config_log`、`.audit_log` を変更しない。 | `.server_config` → `.config_log` → `.audit_log`。 |
| `POST /api/repo-config` | 指定 key の正規化後値が既存値と一致 | `.repo_config`、`.config_log` を変更しない。 | `.repo_config` → `.config_log`。 |
| `POST /api/branch-config` | 正規化後 `branch_targets` が既存値と一致 | `.branch_config`、`.config_log` を変更しない。 | `.branch_config` 作成/置換/削除 → `.config_log`。 |
| `POST /api/notify-config` | secret mask 適用後の比較で既存値と一致 | `.notify_config`、`.config_log` を変更しない。 | `.notify_config` → `.config_log`。 |
| `POST /api/smtp-config` | config と password 更新有無が既存値と一致 | `.smtp_config`、`.smtp_secret`、`.config_log` を変更しない。 | `.smtp_config` → 必要時 `.smtp_secret` → `.config_log`。 |
| `POST /api/dashboard-layout` | widgets 配列が既存値と一致 | `.dashboard_layout`、`.config_log` を変更しない。 | `.dashboard_layout` → `.config_log`。 |

no-op response は endpoint 固有の `No changes` が定義されている場合はその文言を返す。定義がない endpoint は通常成功文言を返してよいが、状態ファイル、JSON Lines、監査ログ、通知ログに差分を作ってはならない。部分更新では未指定 key を保持し、`null` が削除を意味する key は個別節に明記された key だけとする。

**削除 / 失効 response：**

| Endpoint | 成功 response | 不在時 |
|----------|---------------|--------|
| `DELETE /api/snapshots/{id}` | `{message:"Snapshot deleted"}` | `404` |
| `DELETE /api/hooks/{id}` | `{message:"Hook deleted"}` | `404` |
| `DELETE /api/alert-rules/{id}` | `{message:"Alert rule deleted"}` | `404` |
| `DELETE /api/tag-rules/{id}` | `{message:"Tag rule deleted"}` | `404` |
| `DELETE /api/tokens/{id}` | `{message:"Token revoked"}` | `404` |
| `DELETE /api/queue` | `{message:"Queue cleared",cleared_count}` | queue 空でも `200`、`cleared_count=0` |

**SSE frame 契約：**

`GET /api/build/stream` は `text/event-stream; charset=utf-8` を返し、各 frame は以下の JSON を `data:` 行に 1 件ずつ出力する。

```json
{"type":"log","line":"[INFO] build started","at":"2026-09-16T10:00:00Z"}
{"type":"end","status":"success","duration_seconds":12}
```

`type` は `"log"` または `"end"` だけとする。`log.line` は最大 4000 文字とし、超過分は末尾を切り捨てる。`end` frame は接続終了前に 1 回だけ送る。送信中に client が切断した場合、状態ファイル、history、log を変更しない。

**バイナリ response 契約：**

`GET /api/snapshots/{id}/download` は JSON error 以外では binary response とし、成功時に JSON body を返さない。`Content-Type: application/octet-stream`、`Content-Disposition: attachment; filename="{id}.tar.gz"` を付与する。`id` に `"`、`\`、改行を含む値は path 検証で `422` とする。

**レスポンス検証条件：**

| ケース | 期待結果 |
|--------|----------|
| schema 必須 key | 全 endpoint が必須 key を返す。 |
| 空配列 | 配列 key は `[]`、`null` や key 省略なし。 |
| 変更なし config | `No changes`、不要な `.config_log` 追記なし。 |
| delete queue 空 | `200`、`cleared_count=0`。 |
| SSE 正常終了 | `end` frame が 1 回だけ送信される。 |
| snapshot download | binary body、固定 header、JSON success body なし。 |

**API / SDK / UI / 状態ファイル 横断固定契約：**

下表の機能群は、API endpoint、SDK method、UI 操作、状態ファイル副作用を同じ実装単位でそろえる。API だけ、SDK だけ、UI だけを先行して仕様外の仮実装にしてはならない。UI が未実装の Phase では、UI 列は fixture の期待操作として固定し、実装完了扱いには含めない。

| 機能群 | API | SDK | UI 操作 | 状態ファイル副作用 | 成功後再取得 | 失敗時固定 |
|--------|-----|-----|---------|--------------------|--------------|------------|
| login / session | `POST /api/login`, `POST /api/login/totp`, `POST /api/logout`, `GET /api/sessions`, `POST /api/sessions/revoke-all` | `login()`, `loginTotp()`, `logout()`, `getSessions()`, `revokeAllSessions()` | login、TOTP 確認、logout、session 一括失効。 | `.admin_credentials`、`.access_log`、`.audit_log`、memory session。token 本体は永続化しない。 | login 成功後は `getDashboard()` → `getStatus()` → `getQueue()`。session 失効後は `getSessions()`。 | `401` は SDK token / ticket / secret field を破棄し、UI は `panel-login` のみ表示する。 |
| build control | `POST /api/build`, `POST /api/build/force`, `POST /api/build/cancel`, `GET /api/build/stream`, `GET /api/queue`, `DELETE /api/queue` | `triggerBuild()`, `buildForce()`, `cancelBuild()`, `streamBuild()`, `getQueue()`, `clearQueue()` | manual build、force build、cancel、stream 開始/停止、queue clear。 | `.build_state`、`.build_lock`、`.build_logs/{id}.json`、queue entry。force build 時だけ SHA cache 更新。 | `getStatus()` → `getQueue()`、stream end 後は `getLogs()` も実行。 | running は `409`、queue full は `429`、maintenance / circuit は `503` または仕様済み `409`。UI は同じ build request を自動再送しない。 |
| logs / history | `GET /api/logs`, `GET /api/logs/search`, `GET /api/history`, `GET /api/history/{id}/log`, `POST /api/history/{id}/comment`, `POST /api/history/{id}/flag`, `POST /api/history/{id}/tags` | `getLogs()`, `searchLogs()`, `getHistory()`, `getHistoryLog()`, `setHistoryComment()`, `setHistoryFlag()`, `setHistoryTags()` | log 表示、検索、history 表示、comment / flag / tag 保存。 | read-only GET は副作用なし。comment / flag / tag は `.build_logs/{id}.json` と必要時 `.build_history`、`.config_log`。 | 変更系は `getHistory()`、comment は `getHistoryComment(id)` も実行。 | `404` は選択解除または not found 表示。`422 details` は field error。壊れた JSON Lines は response に含めない。 |
| config / repo / branch / schedule | `GET/POST /api/config`, `POST /api/config/validate`, `GET/POST /api/repo-config`, `GET/POST /api/branch-config`, `GET /api/schedule`, `POST /api/schedule/*` | `getConfig()`, `setConfig()`, `validateConfig()`, `setRepoConfig()`, `getBranchConfig()`, `setBranchConfig()`, `getSchedule()`, schedule 系 method | config 保存、validate、repo 保存、branch 保存、schedule 変更。 | `.server_config`、`.repo_config`、`.branch_config`、`.config_log`。validate は保存なし。systemd 反映失敗時も保存済み値は戻さない。 | 保存系は対象 GET → `getConfigLog()`。validate は再取得なし。 | `422` は書込前停止。systemd 失敗 `500` は保存済み値を再取得して表示する。no-op は状態ファイルと log を変更しない。 |
| notify / SMTP / webhook | `GET/POST /api/notify-config`, `POST /api/notify-test`, `GET /api/notify-log`, `GET/POST /api/smtp-config`, `POST /api/smtp-test`, `GET/POST /api/webhook-config`, `GET /api/webhook-events`, `POST /api/notify/weekly-summary` | `getNotifyConfig()`, `setNotifyConfig()`, `notifyTest()`, `getNotifyLog()`, `getSmtpConfig()`, `setSmtpConfig()`, `smtpTest()`, `getWebhookConfig()`, `setWebhookConfig()`, `getWebhookEvents()`, `notifyWeeklySummary()` | notify 保存、test、SMTP 保存/test、webhook secret 保存、event 表示、weekly summary。 | `.notify_config`、`.smtp_config`、`.smtp_secret`、`.webhook_secret`、`.notify_log`、`.webhook_events.json`、`.config_log`。 | 保存系は対象 GET → `getConfigLog()`。test / summary は `getNotifyLog()`。 | secret は response / log / fixture へ平文出力しない。保存成功・失敗とも UI secret field を消去する。 |
| snapshots / rollback / maintenance | `GET /api/snapshots`, `GET /api/snapshots/{id}/download`, `DELETE /api/snapshots/{id}`, `POST /api/history/{id}/rollback`, `GET /api/maintenance`, `POST /api/maintenance/enable`, `POST /api/maintenance/disable` | `getSnapshots()`, `downloadSnapshot()`, `deleteSnapshot()`, `rollbackHistory()`, `getMaintenance()`, `enableMaintenance()`, `disableMaintenance()` | snapshot list/download/delete、rollback、maintenance enable/disable。 | `.snapshots/`、`.build_history`、`.build_logs/{new_id}.json`、`.maintenance`、`.config_log`。download は副作用なし。 | delete は `getSnapshots()`。rollback は `getHistory()` → `getStatus()`。maintenance は `getMaintenance()`。 | delete / rollback は確認必須。running rollback は `409`。maintenance enabled 中は build / rollback / 設定変更系を disabled。 |
| access / hooks / rules / pipeline / notes / layout | access、hooks、alert rules、tag rules、pipeline config、notes、dashboard layout の GET/POST/DELETE endpoint | 対応する §23 SDK method | 保存、追加、削除、notes 保存、dashboard layout 保存。 | `.access_control`、`.hooks`、`.alert_rules`、`.tag_rules`、`.pipeline_config`、`.notes`、`.dashboard_layout`、`.config_log`。 | 対象 GET → 必要時 `getConfigLog()`。layout は `getDashboardLayout()` → `getDashboard()`。 | duplicate `409` は競合表示。validation `422` は field error。削除対象不在は `404`。 |
| tokens / audit / access logs / rate limit | `GET /api/tokens`, `POST /api/tokens`, `DELETE /api/tokens/{id}`, `GET /api/audit-log`, `GET /api/access-log`, `GET /api/api-access-log`, `GET/POST /api/api-rate-limit` | `getTokens()`, `createToken()`, `revokeToken()`, `getAuditLog()`, `getAccessLog()`, `getApiAccessLog()`, `getApiRateLimit()`, `setApiRateLimit()` | token 発行/失効、audit / access log 表示、rate limit 保存。 | `.api_tokens`、`.audit_log`、`.access_log`、`.api_access_log`、`.api_rate_state`、`.server_config`、`.config_log`。token 本体は作成時 response のみ。 | token 操作は `getTokens()` → `getAuditLog()`。rate limit は `getApiRateLimit()`。 | token 本体は再取得不可。`403` は logout しない。`429` は rate limit 表示し、同一操作を自動 retry しない。 |

**横断処理順固定：**

| 処理種別 | 固定順序 |
|----------|----------|
| 認証必須 JSON API | method / path 判定 → body 禁止判定 → JSON parse → 認証 / scope → rate limit → endpoint 固有 validation → read → write 計画 → atomic write → JSON Lines 追記 → response。 |
| read-only API | method / path 判定 → body 禁止判定 → 認証 / scope → query validation → read → 壊れた任意行除外 → response。read-only API は状態ファイルを書き換えない。 |
| UI 変更操作 | panel error / success 消去 → UI 入力検証 → 対象操作 disabled → SDK 呼び出し → 成功後再取得 → success 表示 → secret 消去 → disabled 再評価。 |
| UI 取得操作 | panel error 消去 → 対象操作 disabled → SDK 呼び出し → DOM 更新 → empty state 判定 → disabled 再評価。success 表示は行わない。 |
| SDK request | 引数検証 → path / query / body 生成 → Authorization 付与 → timeout 設定 → fetch → status 判定 → response parse → token 変化適用 → return / throw。 |
| multi-file write | 全入力検証 → 全対象 read → 全 write payload 生成 → §22.0d の Write 順に atomic write → JSON Lines 追記 → response。途中失敗時は未処理ファイルを書かない。 |

**横断 fixture 固定：**

| fixture | 入力 | 必須確認 |
|---------|------|----------|
| cross auth expired | 任意の認証必須 API が `401`。 | SDK は token を破棄し、UI は全 secret field を消去して `panel-login` だけを表示する。対象 API の状態ファイル副作用なし。 |
| cross config no-op | `POST /api/config` に既存値と同一の正規化済み body。 | API は `No changes`、SDK は response をそのまま返し、UI は成功表示する。`.server_config`、`.config_log`、`.audit_log` に差分なし。 |
| cross validation details | 任意の保存 API が `422 details`。 | SDK は `AdlaireCIError.details` を保持し、UI は該当 field と panel summary に表示する。状態ファイルを書かない。 |
| cross refresh failure | 変更 API は成功し、成功後再取得の 2 件目が `500`。 | 変更副作用は維持し、同じ変更 API を再実行しない。UI は操作成功を `global-success`、再取得失敗を panel error に分けて表示する。 |
| cross secret failure | secret 保存 API が `500`。 | SDK error に secret 原文を含めず、UI は secret field を消去する。状態ファイル、log、fixture に secret 原文が残らない。 |
| cross stream invalid frame | `GET /api/build/stream` が parse 不能 frame を返す。 | SDK は `AdlaireCIError(status=0,message="Invalid SSE frame")`、UI は stream error を表示し、status / queue を再取得する。状態ファイルは変更しない。 |
| cross binary snapshot | `downloadSnapshot(id)` が binary success。 | API は binary header、SDK は `Blob`、UI は download 開始表示。JSON parse、success JSON body、状態ファイル更新なし。 |
| cross destructive cancel | 削除 / rollback 確認 dialog を cancel。 | SDK method 呼び出し 0 回、状態ファイル副作用なし、success / error 表示差分なし。 |

### 22.0f 実装優先度

対象項目の実装時は、下表の順に進める。上位の完了条件を満たす前に下位へ進んではならない。同一優先度内では、API、SDK、UI、状態ファイル、検証手順を同じ Pull Request で同期する。

| 優先度 | 対象 | 完了条件 |
|--------|------|----------|
| P0 | 認証、セッション、共通エラー、状態ファイル読み書き、`.access_log`、`.config_log` | `POST /api/login` から認証必須 API の共通処理までが §22.0〜§22.0e と一致し、秘密情報がログとレスポンスに出ない。 |
| P1 | ビルド操作、status、logs、history、queue、circuit breaker | 手動ビルド、強制ビルド、キャンセル、キュー、履歴、ログ取得が同一状態ファイル契約で動作する。 |
| P2 | config、repo、branch、schedule、PAT、diagnostics、dashboard | 設定変更が `.config_log` に残り、GET 系集約 API が状態ファイルを更新しない。 |
| P3 | notify、SMTP、webhook、webhook config、weekly summary | 通知送信責務が `components/runner.go`、設定責務が `components/api.go` に分離され、secret はマスクされる。 |
| P4 | snapshots、rollback、maintenance、access control、hooks | 運用系 API が `409`、`422`、`503` を仕様どおり返し、ロールバックは履歴に `trigger: "rollback"` を残す。 |
| P5 | alert rules、tag rules、pipeline config、notes、dashboard layout、tokens | 拡張設定が schema どおり保存され、SDK と UI の操作名が §22.0e と一致する。 |

各優先度の検証条件は以下とする。

| 優先度 | 必須検証 |
|--------|----------|
| P0 | 認証成功、認証失敗、期限切れ session、`401` 時 SDK token 破棄、`.access_log` 追記、秘密情報マスクを確認する。 |
| P1 | 手動 build、force build、running 中の queue、cancel、history/log 取得、`409`、`429`、`503` を確認する。 |
| P2 | config/repo/branch/schedule の保存、`.config_log` 追記、GET 系 API が状態ファイルを書き換えないことを確認する。 |
| P3 | Webhook test、weekly summary、SMTP test、webhook secret 保存、secret mask、通知失敗ログを確認する。 |
| P4 | snapshot list/download/delete、rollback、maintenance enable/disable、access control block、hook success/failure を確認する。 |
| P5 | rule 追加/削除、pipeline config 保存、notes 保存、dashboard layout 保存、token 発行/失効、token 本体が再取得不可であることを確認する。 |

**API P0 / P1 fixture 固定：**

P0 / P1 実装は、下表の fixture をすべて満たした場合だけ完了扱いにする。fixture は実装言語の test case 名または subtest 名へそのまま写せる粒度とし、期待 HTTP status、期待 body、状態ファイル副作用を同時に確認する。

| Fixture | 入力状態 / Request | 期待 response | 状態ファイル副作用 |
|---------|--------------------|---------------|--------------------|
| A1 common route errors | 未定義 `/api/unknown`、既存 path への未許可 method、body 禁止 endpoint への body 付き request、JSON 不正文を順に送る。 | `404 {"error":"Not found"}`、`405 {"error":"Method not allowed"}`、`400 {"error":"Request body is not allowed"}`、`400 {"error":"Invalid JSON"}`。 | 状態ファイルを作成、更新、削除しない。 |
| A2 auth errors | Bearer なし、無効 token、期限切れ session、scope 不足 API token で認証必須 endpoint を呼ぶ。 | `401 {"error":"Unauthorized"}` または `403 {"error":"Forbidden"}`。 | 秘密情報を response、`.access_log`、server log に出さない。 |
| A3 status primary | 正常な `.build_status.json`、正常な `.build_lock`、`.pending_transfers`、`.build_circuit_state` を置いて `GET /api/status`。 | `.build_status.json` の値を正とし、valid running lock がある場合だけ `running:true`。 | GET 副作用なし。mtime、mode、内容が変わらない。 |
| A4 status fallback | `.build_status.json` 不在、`.build_history` に成功履歴、`.build_state` に queue、`.build_lock` 不在で `GET /api/status`。 | `status`、`last_*`、`queued`、`running:false` を fallback 算出する。 | `.build_status.json` を生成しない。 |
| A5 status corrupted | `.build_status.json` を不正 JSON にして `GET /api/status`。 | `500 {"error":"State file is corrupted"}`。 | 退避ファイル作成、自動修復、初期値上書きを行わない。 |
| A6 history jsonl filtering | `.build_history` に有効行 2 件、空行 1 件、不正 JSON 1 件、必須 key 不足 1 件を置き `GET /api/history?page=1&per_page=10`。 | `total:2`、`pages:1`、`history` は有効行だけを新しい順で返す。 | 壊れた行を書き戻し削除しない。server log に `BUILD_HISTORY_SKIP_CORRUPT`。 |
| A7 build log lookup | `.build_logs/{id}.json` 正常、通常ログ不在で archive 正常、両方不在、破損ログをそれぞれ `GET /api/history/{id}/log`。 | 正常は log object、archive は gzip 展開結果、両方不在は `404 {"error":"Not found"}`、破損は `500 {"error":"State file is corrupted"}`。 | archive を通常ログへ復元しない。破損ログを上書きしない。 |
| A8 queue state | `.build_state` 不在、正常、破損の 3 状態で `GET /api/queue`。 | 不在は `queued:[]`、正常は保存値、破損は `500 {"error":"State file is corrupted"}`。 | 不在時も `.build_state` を生成しない。 |
| A9 build conflict | valid running `.build_lock`、形式不正 `.build_lock`、PID 判定不能 lock、queue 上限到達状態で `POST /api/build`。 | lock 系は `409 {"error":"Conflict"}`。queue 上限到達は endpoint 固有の `429`。 | 失敗時に `.build_state`、`.build_history`、`.build_status.json` を変更しない。 |
| A10 circuit breaker | `.build_circuit_state` が `open:true` の状態で `POST /api/build`、次に `POST /api/circuit-breaker/reset`。 | build は `503` または endpoint 固有 circuit open error。reset は `200` と固定 message。 | build 失敗では queue 追加なし。reset は `.build_circuit_state` だけを reset 後値へ atomic write。 |
| A11 write lock timeout | `.build_state.lock` を保持した状態で `.build_state` 更新 endpoint を呼ぶ。 | 10 秒経過後 `409 {"error":"Conflict"}`。 | tmp file を残さず、target を変更しない。 |
| A12 GET side-effect zero | `GET /api/status`、`GET /api/history`、`GET /api/logs`、`GET /api/queue` を連続実行する。 | 各 endpoint は入力状態に応じた正常 response または固定 error response。 | request 前後で対象状態ファイル一覧、mtime、mode、内容が一致する。 |

**API P2〜P5 fixture 固定：**

P2〜P5 実装は、下表の fixture をすべて満たした場合だけ完了扱いにする。fixture は既存 endpoint と既存状態ファイルだけを対象とし、§22.0e にない endpoint、§22.0a にない状態ファイル、§24 にない UI 操作を追加してはならない。

| Fixture | 優先度 | 入力状態 / Request | 期待 response | 状態ファイル副作用 |
|---------|--------|--------------------|---------------|--------------------|
| B1 config no-op | P2 | 既存 `.server_config` と同じ body で `POST /api/config`。 | `200` と `{message:"No changes",config}`。 | `.server_config`、`.config_log`、`.audit_log` を変更しない。 |
| B2 config validation failure | P2 | `queue_max_size=-1`、未知 enum、相対 path を含む `POST /api/config`。 | `422 {"error":"Validation failed","details":[...]}`。 | 状態ファイルを変更しない。 |
| B3 branch config save | P2 | 有効な branch target 2 件で `POST /api/branch-config`。 | `{message:"Branch config updated",branches_count:2}`。 | `.branch_config` を atomic write し、`.config_log` に差分を記録する。 |
| B4 schedule systemd failure | P2 | `.server_config` 保存成功後、systemd timer 更新を fake failure にする。 | `500`。 | `.server_config` は更新済み、`.config_log` に `systemd_update_failed` を記録し、未定義 rollback を行わない。 |
| B5 PAT update secret mask | P2 | `POST /api/pat-update` に token を送る。 | `{message:"PAT updated"}`。 | secret file は mode `600`。response、`.config_log`、`.audit_log`、server log に token 平文を出さない。 |
| B6 dashboard read-only | P2 | `.dashboard_layout`、`.build_state`、`.build_status.json`、`.alert_rules` を置き `GET /api/dashboard`。 | widget 順に dashboard object を返す。 | GET は対象状態ファイルを作成、修復、更新しない。 |
| C1 notify config mask | P3 | Webhook secret と SMTP password を含む通知設定保存後、GET / backup / log を確認する。 | secret は `"***"` または `*_set:true` だけを返す。 | secret 平文を状態表示、履歴、通知ログ、backup に残さない。 |
| C2 webhook receive signed | P3 | 正常署名の GitHub push payload を `POST /api/webhook`。 | `202` と queued 結果。 | `.webhook_events.json` 追記後、必要時 `.build_state.queued` へ `trigger:"webhook"` を追加する。 |
| C3 webhook invalid signature | P3 | 署名なし、不正 prefix、不一致署名。 | `401 {"error":"Unauthorized"}`。 | event log、queue、history を変更しない。 |
| C4 SMTP test disabled | P3 | SMTP disabled で `POST /api/smtp-test`。 | endpoint 固有の `422`。 | `.notify_log` へ成功扱いを残さず、secret を出力しない。 |
| D1 snapshot delete | P4 | 存在する snapshot id で `DELETE /api/snapshots/{id}`。 | `{message:"Snapshot deleted"}`。 | 対象 snapshot だけ削除し、`.config_log` または監査対象 log に削除を記録する。 |
| D2 rollback running conflict | P4 | `.build_state.running=true` で `POST /api/history/{id}/rollback`。 | `409 {"error":"Build is running"}`。 | queue、history、snapshot、deploy target を変更しない。 |
| D3 maintenance blocks build | P4 | maintenance enabled 状態で `POST /api/build`。 | `503` と endpoint 固有 maintenance error。 | build queue、history、log を変更しない。 |
| D4 access-control deny | P4 | allowlist に接続元が含まれない状態で任意認証必須 API。 | 認証判定前に `403 {"error":"Forbidden"}`。 | password / token 検証、access log 成功行、対象操作副作用を行わない。 |
| D5 hook timeout | P4 | `pre` hook が timeout。 | build status は `hook_error` または endpoint 固有の hook error。 | pipeline を実行せず、hook log と build log に timeout を固定値で記録する。 |
| E1 token issue once | P5 | `POST /api/tokens` で token 作成。 | token 本体を作成 response に 1 回だけ含める。 | `.api_tokens` には hash だけを保存し、再取得 API では token 本体を返さない。 |
| E2 token revoke missing | P5 | 存在しない token id を `DELETE /api/tokens/{id}`。 | `404 {"error":"Not found"}`。 | `.api_tokens`、`.audit_log` を変更しない。 |
| E3 alert/tag duplicate | P5 | 同一 alert rule または tag rule を 2 回作成。 | 2 回目は `409 {"error":"Conflict"}`。 | 2 回目は該当状態ファイルと `.config_log` を変更しない。 |
| E4 pipeline config reserved arg | P5 | `extra_args` に `--src`、`--out`、`--state-dir` を含める。 | `422 {"error":"Validation failed","details":[...]}`。 | `.pipeline_config` を変更しない。 |
| E5 notes same content | P5 | 同じ `content` を 2 回 `POST /api/notes`。 | 2 回目は `No changes`。 | 2 回目は `.notes`、`.config_log` を変更しない。 |
| E6 dashboard layout invalid | P5 | 重複 widget、未知 widget、空配列を `POST /api/dashboard-layout`。 | `422 {"error":"Validation failed","details":[...]}`。 | `.dashboard_layout` を変更しない。 |

| メソッド | パス | 認証 | 説明 |
|---------|------|------|------|
| `POST` | `/api/login` | 不要 | ログイン（セッショントークン返却） |
| `POST` | `/api/logout` | 要 | ログアウト（セッション破棄） |
| `POST` | `/api/change-password` | 要 | パスワード変更 |
| `GET` | `/api/access-log` | 要 | ログイン、ログアウト、API token 監査ログを返す |
| `GET` | `/api/sessions` | 要 | 有効セッション一覧を返す |
| `POST` | `/api/sessions/revoke-all` | 要 | 現セッション以外の全セッションを強制無効化する |
| `GET` | `/api/status` | 要 | 最終ビルド時刻・SHA・成否・実行中フラグを返す |
| `POST` | `/api/build` | 要 | 手動ビルドトリガー（`components/runner.go` を即時起動） |
| `POST` | `/api/build/force` | 要 | SHA リセットとビルドをアトミックに実行する（強制ビルド） |
| `POST` | `/api/build/cancel` | 要 | 実行中のビルドを強制停止する（`running: true` のときのみ有効） |
| `GET` | `/api/build/stream` | 要 | 実行中または直近ビルドログを SSE で配信する |
| `GET` | `/api/logs?n=100&q=<keyword>` | 要 | 最新ビルドログを n 行返す（`q` 省略時は全行） |
| `GET` | `/api/logs/search?q=<keyword>&from=<date>&to=<date>&level=<warn\|error>` | 要 | 日付範囲・重大度を指定して過去ビルドログ（`.build_logs/`）を横断検索する |
| `GET` | `/api/logs/export` | 要 | ビルドログ全件を JSON 形式でエクスポートする |
| `POST` | `/api/logs/cleanup` | 要 | 保持期間（`log_retention_days`）を超えた `.build_logs/` エントリを削除する |
| `GET` | `/api/history?page=<n>&per_page=<n>` | 要 | 過去ビルド履歴一覧をページ指定で返す（省略時: `page=1`, `per_page=20`） |
| `GET` | `/api/history/export` | 要 | ビルド履歴一覧を JSON 形式でエクスポートする |
| `GET` | `/api/history/{id}/log` | 要 | 指定ビルド ID のログを取得する |
| `GET` | `/api/history/{id}/comment` | 要 | 指定ビルドのコメントを取得する |
| `POST` | `/api/history/{id}/comment` | 要 | 指定ビルドにコメントを付与・更新する |
| `POST` | `/api/history/{id}/flag` | 要 | 指定ビルドに重要フラグを設定・解除する |
| `POST` | `/api/history/{id}/tags` | 要 | 指定ビルドのタグを置換する |
| `POST` | `/api/history/{id}/rollback` | 要 | 指定ビルド ID のスナップショットから SSH 転送を再実行する（→ §14b） |
| `GET` | `/api/sysinfo` | 要 | 出力サイトサイズ・更新日時・稼働時間を返す |
| `GET` | `/api/health`        | 不要 | 死活監視用ヘルスチェック |
| `GET` | `/api/schedule` | 要 | systemd timer の次回実行予定時刻を返す |
| `POST` | `/api/schedule/interval` | 要 | systemd タイマーのポーリング間隔を動的変更する |
| `POST` | `/api/schedule/pause` | 要 | ポーリングを一時停止する |
| `POST` | `/api/schedule/resume` | 要 | ポーリングを再開する |
| `POST` | `/api/schedule/allowed-hours` | 要 | 自動ビルド許可時間帯を設定・解除する |
| `POST` | `/api/schedule/force-interval` | 要 | `FORCE_BUILD_INTERVAL`（強制再ビルド間隔）を動的変更する |
| `POST` | `/api/schedule/cooldown` | 要 | `BUILD_COOLDOWN_SECONDS`（ビルドクールダウン秒数）を動的変更する |
| `GET` | `/api/notify-config` | 要 | Webhook 通知設定を返す |
| `POST` | `/api/notify-config` | 要 | Webhook 通知設定を更新する |
| `GET` | `/api/notify-log` | 要 | Webhook 送信履歴（日時・イベント・HTTP ステータス・成否）を返す |
| `POST` | `/api/notify-test` | 要 | Webhook にテスト通知を送信し疎通を確認する |
| `POST` | `/api/notify/weekly-summary` | 要 | 週次サマリー Webhook を即時手動送信する（過去 7 日間の統計を集計して送信） |
| `GET` | `/api/config` | 要 | サーバー設定を返す |
| `POST` | `/api/config` | 要 | サーバー設定を更新する |
| `POST` | `/api/log-level` | 要 | `components/api.go` の `log_level` を変更する |
| `GET` | `/api/config-log` | 要 | 設定変更履歴（変更日時・種別・変更前後の値）を返す |
| `GET` | `/api/pat-status` | 要 | GitHub PAT の有効性確認 |
| `POST` | `/api/pat-verify` | 要 | GitHub API を呼び出し PAT の有効性をリアルタイム検証する |
| `POST` | `/api/pat-update` | 要 | `.github_token` ファイルを更新し PAT を差し替える |
| `GET` | `/api/stats?days=7` | 要 | ビルド統計（成功率・回数・平均間隔）を返す |
| `GET` | `/api/stats/timeline?days=30` | 要 | 日別ビルド成功/失敗件数の時系列配列を返す |
| `GET` | `/api/stats/build-duration?n=20` | 要 | 過去 N 件のビルド所要時間統計（平均・最小・最大・直近リスト）を返す |
| `GET` | `/api/stats/build-trends?n=100` | 要 | `.build_trends.json` から所要時間 trend、中央値、p95、異常件数を返す |
| `GET` | `/api/output-meta` | 要 | 出力サイトのサイズ・見出し数・生成日時・前回比サイズ差分を返す |
| `GET` | `/api/repo-info` | 要 | リポジトリ設定（OWNER/REPO/BRANCH/TARGET_FILE）を返す |
| `POST` | `/api/repo-config` | 要 | リポジトリ監視設定（OWNER / REPO / BRANCH / TARGET_FILE）を更新する |
| `GET` | `/api/branch-config` | 要 | 現在有効なブランチターゲット設定を返す |
| `POST` | `/api/branch-config` | 要 | ブランチターゲット設定を `.branch_config` へ書き込む |
| `GET` | `/api/backup` | 要 | 全設定（通知設定・サーバー設定）を JSON 形式でエクスポートする |
| `POST` | `/api/restore` | 要 | JSON 形式の設定をインポートし全設定を上書き復元する |
| `GET` | `/api/dashboard` | 要 | ステータス・システム情報・統計・スケジュールを一括返却する |
| `GET` | `/api/diagnostics` | 要 | PAT・GitHub API・出力サイト・systemd・Webhook の一括自己診断結果を返す |
| `GET` | `/api/rate-limit` | 要 | GitHub API のレート制限残量・上限・リセット時刻を返す |
| `GET` | `/api/disk-usage` | 要 | ビルドログ合計・出力サイトのディスク使用量を返す |
| `GET` | `/api/webhook-events?limit=50&offset=0` | 要 | 受信 Webhook イベント一覧を新しい順にページネーション付きで返す（→ `.webhook_events.json`） |
| `POST` | `/api/webhook` | 不要（Secret 検証） | GitHub push Webhook を受信し、署名検証後にビルドをトリガーする（→ §22 Webhook 受信仕様） |
| `GET` | `/api/webhook-config` | 要 | Webhook Secret 設定状態を返す |
| `POST` | `/api/webhook-config` | 要 | Webhook Secret を設定する |
| `POST` | `/api/circuit-breaker/reset` | 要 | サーキットブレーカーをリセットする（`open: false`・`consecutive_failures: 0` に戻しポーリングを再開） |
| `GET` | `/api/snapshots` | 要 | ビルド成果物スナップショット一覧を返す |
| `GET` | `/api/snapshots/{id}/download` | 要 | 指定スナップショットをダウンロードする |
| `DELETE` | `/api/snapshots/{id}` | 要 | 指定スナップショットを削除する |
| `GET` | `/api/maintenance` | 要 | メンテナンスモード状態を返す |
| `POST` | `/api/maintenance/enable` | 要 | メンテナンスモードを有効化する |
| `POST` | `/api/maintenance/disable` | 要 | メンテナンスモードを無効化する |
| `GET` | `/api/access-control` | 要 | IP / CIDR 許可リストを返す |
| `POST` | `/api/access-control` | 要 | IP / CIDR 許可リストを置換する |
| `GET` | `/api/hooks` | 要 | ビルドフック一覧を返す |
| `POST` | `/api/hooks` | 要 | ビルドフックを追加する |
| `DELETE` | `/api/hooks/{id}` | 要 | ビルドフックを削除する |
| `GET` | `/api/hooks/{id}/log` | 要 | ビルドフック実行ログを返す |
| `GET` | `/api/alert-rules` | 要 | アラートルール一覧を返す |
| `POST` | `/api/alert-rules` | 要 | アラートルールを追加する |
| `DELETE` | `/api/alert-rules/{id}` | 要 | アラートルールを削除する |
| `GET` | `/api/tag-rules` | 要 | 自動タグ付けルール一覧を返す |
| `POST` | `/api/tag-rules` | 要 | 自動タグ付けルールを追加する |
| `DELETE` | `/api/tag-rules/{id}` | 要 | 自動タグ付けルールを削除する |
| `POST` | `/api/verify-output` | 要 | 現在の出力サイト checksum を検証する |
| `GET` | `/api/pipeline-config` | 要 | ビルドパイプライン設定を返す |
| `POST` | `/api/pipeline-config` | 要 | ビルドパイプライン設定を置換する |
| `GET` | `/api/build-chain-config` | 要 | ビルド依存チェーン設定を返す |
| `POST` | `/api/build-chain-config` | 要 | ビルド依存チェーン設定を置換する |
| `GET` | `/api/notes` | 要 | 運用ノートを返す |
| `POST` | `/api/notes` | 要 | 運用ノートを保存する |
| `GET` | `/api/smtp-config` | 要 | SMTP 設定を返す |
| `POST` | `/api/smtp-config` | 要 | SMTP 設定を保存する |
| `POST` | `/api/smtp-test` | 要 | SMTP テスト送信を行う |
| `GET` | `/api/queue` | 要 | ビルドキュー状態を返す |
| `DELETE` | `/api/queue` | 要 | 待機中ビルドキューを削除する |
| `GET` | `/api/approvals` | 要 | ビルド承認待ち entry を返す |
| `POST` | `/api/approvals/{id}/approve` | 要 | 指定承認 entry を承認し build queue へ投入する |
| `POST` | `/api/approvals/{id}/reject` | 要 | 指定承認 entry を却下する |
| `GET` | `/api/dashboard-layout` | 要 | ダッシュボードウィジェット設定を返す |
| `POST` | `/api/dashboard-layout` | 要 | ダッシュボードウィジェット設定を置換する |
| `GET` | `/api/tokens` | 要 | API token 一覧を返す。token 本体は返さない |
| `POST` | `/api/tokens` | 要 | API token を発行する。token 本体は作成時のみ返す |
| `DELETE` | `/api/tokens/{id}` | 要 | API token を失効する |

**`POST /api/login` リクエスト / レスポンス：**
```json
// リクエスト
{ "password": "admin" }

// レスポンス（TOTP 無効）
{ "token": "<session_token>", "must_change": "prompt" }

// レスポンス（TOTP 有効）
{ "totp_required": true, "ticket": "<login_ticket>", "must_change": "none" }
```

`must_change` の有効値：`"none"`（変更不要）| `"prompt"`（促す：初回ログイン時）| `"forced"`（強制：5 回目以降。変更完了まで管理画面の操作を制限）

**`POST /api/login/totp` リクエスト / レスポンス：**
```json
// リクエスト
{ "ticket": "<login_ticket>", "code": "123456" }

// レスポンス
{ "token": "<session_token>", "must_change": "none" }
```

**`POST /api/logout` リクエスト / レスポンス：**
```json
// リクエスト: なし（Bearer トークンのみ）
// レスポンス: 200
{ "message": "Logged out" }
```

**`POST /api/change-password` リクエスト / レスポンス：**
```json
// リクエスト
{ "current_password": "...", "new_password": "..." }
// レスポンス: 200
{ "message": "Password changed" }
```

**`GET /api/status` レスポンス例：**
```json
{
  "last_sha": "abc123",
  "last_build_at": "2026-09-14T10:00:00Z",
  "last_build_status": "success",
  "last_trigger": "polling",
  "last_deploy_status": "success",
  "pending_transfers_count": 0,
  "notify_pending_count": 0,
  "circuit_open": false,
  "output_url": "https://example.com/",
  "running": false
}
```

`last_build_status` の有効値：`"success"` | `"failure"` | `"success_deploy_pending"` | `"skipped_no_change"` | `"skipped_cooldown"` | `"circuit_open"` | `"config_recovered"` | `"config_error"` | `"lock_skipped"` | `"none"`（初回未実行時）
`last_trigger` の有効値：§13 の `trigger` 有効値または `null`
`running` の有効値：`true`（ビルド実行中）| `false`（待機中）
`running` の判定：`.build_state.running == true` または有効な `.build_lock` が存在する場合に `true` を返す。`.build_state.running == false` かつ `.build_lock` が存在しない場合は `false` を返す。形式不正または PID 判定不能な `.build_lock` が存在する場合は、状態競合として `running: true` を返し、API 側でロックを上書きしない。

**`GET /api/logs` レスポンス例：**
```json
{ "lines": ["2026-09-14T10:00:00Z [INFO] Build start", "..."] }
```

**`GET /api/logs/export` レスポンス例：**
```json
{
  "exported_at": "2026-09-15T10:00:00Z",
  "lines": ["2026-09-14T10:00:00Z [INFO] Build start", "..."]
}
```

**`GET /api/history` レスポンス例：**
```json
{
  "total": 42, "page": 1, "per_page": 20, "pages": 3,
  "history": [
    { "id": "b001", "build_at": "2026-09-15T10:00:00Z", "sha": "abc123", "status": "success", "output_size_bytes": 2048576, "trigger": "polling", "duration_seconds": 42, "flagged": false, "tags": ["release"] },
    { "id": "b002", "build_at": "2026-09-14T18:30:00Z", "sha": "def456", "status": "failure", "output_size_bytes": null,    "trigger": "manual", "duration_seconds": 7,  "flagged": true,  "tags": [] }
  ]
}
```

`page` は 1 始まり。`per_page` の最大値は 100。`trigger` は §13 の固定値のみ許可する。`tag` はタグ検証と同じ文字列制約を適用する。`flagged` は `"true"` または `"false"` のみ許可する。未知 query、範囲外の `page` / `per_page`、不正な `trigger` / `tag` / `flagged` は `422` とする。範囲外ページを指定した場合は `history: []` を返す。

`id` はビルド実行時に生成するユニーク識別子（形式：`b{YYYYMMDDHHmmss}`）。`.build_logs/{id}.json` に対応するログファイルが保存される。

**`POST /api/build` / `POST /api/build/force` レスポンス例：**
```json
// 即時開始
{ "message": "Build started", "build_id": "b20260915100500", "queued": false }

// 実行中のため queue へ追加
{ "message": "Build queued", "queued": true }
```

SHA キャッシュのクリアだけを行う専用 API は定義しない。強制再ビルドは必ず `POST /api/build/force` を使用し、SHA reset と build trigger を同一ロック内で実行する。

**`GET /api/sysinfo` レスポンス例：**
```json
{
  "output_size_bytes": 2048576,
  "output_mtime": "2026-09-15T10:00:00Z",
  "uptime_seconds": 86400
}
```

**`GET /api/schedule` レスポンス例：**
```json
{ "next_run_at": "2026-09-15T10:05:00Z", "interval": "5min", "paused": false, "allowed_hours": { "from": 9, "to": 18 } }
```

`paused` が `true` のとき、ポーリングは停止中で `next_run_at` は `null` を返す。

`allowed_hours`：自動ビルドを許可する時間帯（UTC）。`null` = 無制限。`from` 以上 `to` 未満の時刻のみビルドを実行する。許可時間帯外のポーリングでは変更を検出しても実行を保留し、次の許可時間帯に入った時点で実行する。

**`POST /api/schedule/allowed-hours` リクエスト / レスポンス：**
```json
// 設定
{ "from": 9, "to": 18 }
// 解除（無制限に戻す）
{ "from": null, "to": null }
// レスポンス: 200
{ "message": "Allowed hours updated", "allowed_hours": { "from": 9, "to": 18 } }
```

**`POST /api/schedule/pause` / `POST /api/schedule/resume` レスポンス例：**
```json
{ "message": "Schedule paused" }
{ "message": "Schedule resumed" }
```

既に一時停止中に `pause`、または稼働中に `resume` を呼び出した場合は `409 Conflict` を返す。

> **責務分担：** Webhook 通知の**送信責務は `components/runner.go`** にある。`components/runner.go` はビルド完了時に `.notify_config` を読み込んで Webhook を送信する。`components/api.go`（通知 API）は設定の読み書きのみを担い、自身では通知を送信しない。

**`GET /api/notify-config` レスポンス例：**
```json
{
  "webhooks": [
    { "url": "https://hooks.example.com/...", "label": "メイン", "enabled": true, "payload_template": null, "retry_count": 2, "retry_interval_seconds": 30, "secret": null }
  ],
  "channels": [
    { "id": "n001", "type": "webhook", "label": "メイン", "enabled": true, "on": ["failure"], "config": { "url": "https://hooks.example.com/..." }, "retry_count": 2, "retry_interval_seconds": 30 }
  ],
  "on": ["failure"],
  "summary": { "enabled": false, "interval": "weekly", "hour": 9, "day_of_week": 1 },
  "email": { "enabled": false, "to": [], "on": [] }
}
```

`secret`：Webhook 署名シークレット。未設定時は `null`、設定済み時は `"***"`（マスク）を返す（→ 16E 参照）。

`on` の有効値：`"start"`（ビルド開始時）| `"success"`（ビルド成功時）| `"failure"`（ビルド失敗時）| `"deploy_failure"`（転送失敗時）| `"weekly_summary"`（定期サマリー送信時）| `"approval_required"`（承認待ち発生時）| `"duration_anomaly"`（所要時間異常時）| `"config_corrupt"`（設定破損復旧時）。複数指定可。

`summary`：週次サマリー通知の設定。`enabled: true` のとき指定曜日・時刻で統計サマリーを Webhook 送信する。`interval` の有効値は `"weekly"` 固定。`hour` は 0〜23（UTC）。`day_of_week` は 0 = 日曜〜6 = 土曜。

**`POST /api/notify/weekly-summary` レスポンス例：**
```json
{ "message": "Weekly summary sent", "period": "2026-09-08/2026-09-14", "success_count": 12, "failure_count": 1, "success_rate": 92.3 }
```

即時週次サマリー送信。Webhook 未設定または無効時は `422` を返す。

`payload_template`：Webhook 送信 JSON ペイロードのテンプレート文字列。`null` = デフォルトペイロードを使用。テンプレート内で使用可能な変数は以下の通り。

| 変数 | 内容 |
|------|------|
| `{{id}}` | ビルド ID |
| `{{status}}` | ビルド結果（`success` / `failure`） |
| `{{sha}}` | 対象コミット SHA |
| `{{duration_seconds}}` | ビルド所要時間（秒） |
| `{{build_at}}` | ビルド実行日時（ISO 8601） |

`secret`（16E）：Webhook 送信時の HMAC-SHA256 署名用シークレット文字列。設定時はリクエストヘッダーに `X-Adlaire-Signature: sha256=<hmac>` を付与する。`null` = 署名なし。`GET /api/notify-config` で返却する際、設定済みの場合は `"***"` でマスクし、未設定の場合は `null` を返す。`POST /api/notify-config` で更新可能。

**`GET /api/config` レスポンス例：**
```json
{ "log_max_lines": 500, "history_max_count": 100, "build_timeout_seconds": 300, "log_retention_days": 30, "log_archive_after_days": 0, "log_level": "INFO", "pat_expires_at": null, "snapshots_keep": 5, "queue_max_size": 3, "build_retry_max": 0, "build_retry_base_seconds": 5, "commit_status_enabled": false, "commit_status_context": "Adlaire CI", "commit_status_target_url": null, "build_trend_keep_count": 1000, "duration_anomaly": { "enabled": false, "min_samples": 20, "avg_multiplier": 2.0, "p95_multiplier": 1.5 } }
```

`pat_expires_at`：PAT の有効期限日（`YYYY-MM-DD` 形式）。`null` = 未設定。`GET /api/diagnostics` の `pat` 項目で 7 日以内なら `"warn"`、期限当日以前なら `"error"` に変更。

`POST /api/config` で更新可能なキーは `log_max_lines`、`history_max_count`、`build_timeout_seconds`、`log_retention_days`、`log_archive_after_days`、`log_level`、`pat_expires_at`、`snapshots_keep`、`queue_max_size`、`build_retry_max`、`build_retry_base_seconds`、`commit_status_enabled`、`commit_status_context`、`commit_status_target_url`、`build_trend_keep_count`、`duration_anomaly` に限定する。未知キーを含む場合は `422` を返し、既存設定を変更しない。

**`GET /api/health` レスポンス例：**
```json
{
  "status": "ok",
  "last_build_at": "2026-09-15T10:00:00Z",
  "last_build_status": "success",
  "last_deploy_at": "2026-09-15T10:01:00Z",
  "last_deploy_status": "success",
  "pending_transfers": 0,
  "uptime_seconds": 86400
}
```

- `status`：常に `"ok"`（サーバーが応答している限り）
- `last_build_at`：最終ビルド完了日時（未実行時 `null`）
- `last_build_status`：`GET /api/status` の `last_build_status` と同じ値
- `last_deploy_at`：最終 SSH 転送完了日時（未実行時 `null`）
- `last_deploy_status`：`"success"` | `"failure"` | `"skipped"` | `"none"`
- `pending_transfers`：ペンディングキューのエントリ数
- `uptime_seconds`：`components/api.go` 起動からの経過秒数

**`GET /api/pat-status` レスポンス例：**
```json
{ "valid": true, "checked_at": "2026-09-15T10:00:00Z" }
```

**`GET /api/access-log` レスポンス例：**
```json
{ "log": [
    { "at": "2026-09-15T10:00:00Z", "result": "success" },
    { "at": "2026-09-15T09:00:00Z", "result": "failure" }
]}
```

**`GET /api/stats` レスポンス例：**
```json
{
  "days": 7,
  "total_builds": 42,
  "success_count": 40,
  "failure_count": 2,
  "success_rate": 0.952,
  "avg_interval_minutes": 240,
  "avg_duration_seconds": 38,
  "max_duration_seconds": 91
}
```

**`GET /api/repo-info` レスポンス例：**
```json
{
  "owner": "fqwink",
  "repo": "Adlaire-Design-System",
  "branch": "main",
  "target_file": "docs"
}
```

**`GET /api/backup` レスポンス例：**
```json
{
  "exported_at": "2026-09-15T10:00:00Z",
  "notify_config": {
    "webhooks": [{ "url": "https://hooks.example.com/...", "label": "メイン", "enabled": true, "payload_template": null, "retry_count": 2, "retry_interval_seconds": 30, "secret": null }],
    "on": ["failure"],
    "summary": { "enabled": false, "interval": "weekly", "hour": 9, "day_of_week": 1 },
    "email": { "enabled": false, "to": [], "on": [] }
  },
  "server_config": { "log_max_lines": 500, "history_max_count": 100 }
}
```

**`POST /api/notify-test` レスポンス例：**
```json
{ "message": "Test notification sent", "webhook_url": "https://hooks.example.com/..." }
```

**`POST /api/build/force` レスポンス例：**
```json
{ "message": "Build started", "build_id": "b20260915100500", "queued": false }
```

**`POST /api/pat-verify` レスポンス例：**
```json
{ "valid": true, "checked_at": "2026-09-15T10:05:00Z", "scopes": ["contents:read"] }
```

**`GET /api/history/{id}/log` レスポンス例：**
```json
{
  "id": "b001",
  "build_at": "2026-09-15T10:00:00Z",
  "sha": "abc123",
  "status": "success",
  "output_size_bytes": 2048576,
  "trigger": "polling",
  "comment": null,
  "flagged": false,
  "tags": ["release"],
  "lines": ["2026-09-15T10:00:00Z [INFO] Build start", "..."]
}
```

**`POST /api/restore` リクエスト / レスポンス：**
```json
// リクエスト（GET /api/backup と同一形式）
{
  "notify_config": {
    "webhooks": [{ "url": "https://hooks.example.com/...", "label": "メイン", "enabled": true, "payload_template": null, "retry_count": 2, "retry_interval_seconds": 30, "secret": null }],
    "on": ["failure"],
    "summary": { "enabled": false, "interval": "weekly", "hour": 9, "day_of_week": 1 },
    "email": { "enabled": false, "to": [], "on": [] }
  },
  "server_config": { "log_max_lines": 500, "history_max_count": 100 }
}
// レスポンス: 200
{ "message": "Restored" }
```

**`GET /api/notify-log` レスポンス例：**
```json
{ "log": [
    { "at": "2026-09-15T10:00:00Z", "event": "failure", "http_status": 200, "result": "success", "attempt": 1, "error": null },
    { "at": "2026-09-14T18:30:00Z", "event": "start",   "http_status": 500, "result": "failure", "attempt": 3, "error": "HTTP 500" }
]}
```

`event` の有効値：`"start"` | `"success"` | `"failure"` | `"weekly_summary"`（`GET /api/notify-config` の `on` と同一）。`result` の有効値：`"success"` | `"failure"`（Webhook 送信の成否）。

**`GET /api/sessions` レスポンス例：**
```json
{ "sessions": [
    { "created_at": "2026-09-15T09:00:00Z", "expires_at": "2026-09-15T17:00:00Z", "current": true },
    { "created_at": "2026-09-15T08:00:00Z", "expires_at": "2026-09-15T16:00:00Z", "current": false }
]}
```

**`POST /api/sessions/revoke-all` レスポンス例：**
```json
{ "message": "All other sessions revoked", "revoked_count": 1 }
```

**`POST /api/schedule/interval` リクエスト / レスポンス：**
```json
// リクエスト
{ "interval_seconds": 300 }
// レスポンス: 200
{ "message": "Interval updated", "interval_seconds": 300 }
```

**`POST /api/build/cancel` レスポンス例：**
```json
{ "message": "Build cancelled" }
```

`running: false` のときに呼び出した場合は `409 Conflict` → `{"error": "No build is running"}` を返す。

**`GET /api/build/stream` — SSE ストリーミング：**

`Content-Type: text/event-stream` で接続を維持し、ビルドログを逐次配信する。認証は他 API と同じ `Authorization: Bearer {SESSION_TOKEN}` ヘッダーで行う。セッショントークンを query parameter、Cookie、body で受け付けてはならない。

```
data: {"type": "log",  "line": "2026-09-15T10:00:01Z [INFO] Build start"}

data: {"type": "log",  "line": "2026-09-15T10:00:42Z [INFO] Build success"}

data: {"type": "end",  "status": "success", "duration_seconds": 42}
```

`type` の有効値：`"log"`（ログ行）| `"end"`（ビルド完了）。ビルドが未実行時に接続した場合は即時 `{"type": "end", "status": null}` を送信して切断する。

**`POST /api/log-level` リクエスト / レスポンス：**
```json
// リクエスト
{ "level": "DEBUG" }
// レスポンス: 200
{ "message": "Log level changed", "level": "DEBUG" }
```

`level` の有効値：`"DEBUG"` | `"INFO"` | `"WARNING"` | `"ERROR"`

**`POST /api/pat-update` リクエスト / レスポンス：**
```json
// リクエスト
{ "token": "github_pat_..." }
// レスポンス: 200
{ "message": "PAT updated" }
```

**`GET /api/dashboard` レスポンス例：**
```json
{
  "status": {
    "last_sha": "abc123",
    "last_build_at": "2026-09-15T10:00:00Z",
    "last_build_status": "success",
    "output_url": "https://example.com/",
    "running": false
  },
  "sysinfo": {
    "output_size_bytes": 2048576,
    "output_mtime": "2026-09-15T10:00:00Z",
    "uptime_seconds": 86400
  },
  "stats": {
    "days": 7,
    "total_builds": 42,
    "success_count": 40,
    "failure_count": 2,
    "success_rate": 0.952,
    "avg_interval_minutes": 240
  },
  "schedule": { "next_run_at": "2026-09-15T10:05:00Z", "interval": "5min", "paused": false },
  "alerts": [
    { "level": "warn", "message": "PAT expires in 5 days" }
  ]
}
```

`alerts`：診断異常（`warn` / `error`）および PAT 期限切れ間近の場合に項目を返す。異常なしのときは空配列 `[]`。`level` の有効値：`"warn"` | `"error"`。パネルのナビゲーション項目に `alerts` の最高深刻度（`error` > `warn`）のバッジを表示する。

**`GET /api/logs/search` レスポンス例：**
```json
{
  "query": "ERROR",
  "from": "2026-09-10",
  "to": "2026-09-15",
  "results": [
    { "id": "b20260915100000", "build_at": "2026-09-15T10:00:00Z", "lines": ["2026-09-15T10:00:01Z [ERROR] Build failed"] },
    { "id": "b20260912183000", "build_at": "2026-09-12T18:30:00", "lines": ["2026-09-12T18:30:05 [ERROR] Timeout"] }
  ]
}
```

`from` / `to` は `YYYY-MM-DD` 形式。省略時は全期間。`q` 省略時は全行返却。
`level=warn` で `[WARN]` 行のみ、`level=error` で `[ERROR]` 行のみを絞り込む。省略時は全レベルを返却する。

**`GET /api/output-meta` レスポンス例：**
```json
{
  "size_bytes": 2048576,
  "mtime": "2026-09-15T10:00:00Z",
  "heading_count": 342,
  "size_diff_bytes": 1024,
  "tables_count": 128,
  "code_blocks_count": 64,
  "build_warnings": ["未対応記法: admonition (3箇所)"],
  "size_warn": false,
  "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
  "build_id": "b20260915100000",
  "commit_sha": "abc1234",
  "build_at": "2026-09-15T10:00:00Z"
}
```

`size_diff_bytes`：前回ビルド時との差分（正＝増加、負＝減少、`null`＝比較不能）。
前回サイズは `.build_history` の直近エントリに記録された `output_size_bytes` フィールドから取得する。

`tables_count` / `code_blocks_count`：直近ビルドの変換レポート（§8）より取得。ビルド前は `null`。
`build_warnings`：直近ビルドで発生した警告メッセージの配列（§8 参照）。ビルド前は空配列 `[]`。
`build_id` / `commit_sha` / `build_at`：直近ビルドログの `build_meta` を優先し、不在の場合は出力 HTML の meta tag を読み取る。どちらにも存在しない場合は空文字を返す。
値は `components/runner.go` が `.build_logs/{id}.json` または `.build_logs/archive/{id}.json.gz` から最新エントリを読み取って返す。

**`GET /api/stats/timeline` レスポンス例：**
```json
{
  "days": 30,
  "timeline": [
    { "date": "2026-09-15", "success": 3, "failure": 0 },
    { "date": "2026-09-14", "success": 2, "failure": 1 },
    { "date": "2026-09-13", "success": 4, "failure": 0 }
  ]
}
```

日付降順。`days` 日分のうちビルドが 0 件の日はエントリなし。

**`GET /api/stats/build-duration` レスポンス例：**

クエリパラメータ `n`（デフォルト 20）で対象件数を指定する。`.build_logs/{id}.json` の `duration_seconds` フィールドを集計する。

```json
{
  "n": 20,
  "count": 18,
  "avg_seconds": 38.5,
  "min_seconds": 22,
  "max_seconds": 67,
  "recent": [
    { "id": "b20260915100000", "build_at": "2026-09-15T10:00:00Z", "duration_seconds": 42, "status": "success" },
    { "id": "b20260914183000", "build_at": "2026-09-14T18:30:00Z", "duration_seconds": 7,  "status": "failure" }
  ]
}
```

`count` は `duration_seconds` が記録されているビルドの件数（`n` 以下）。`recent` は新しい順。

**`GET /api/webhook-events` レスポンス例：**

クエリパラメータ `limit`（デフォルト 50、上限 200）と `offset` でページネーションする。`.webhook_events.json` を逆順（新しい順）で返す。

```json
{
  "total": 128,
  "offset": 0,
  "limit": 50,
  "events": [
    {
      "timestamp": "2026-09-15T10:00:00Z",
      "delivery_id": "abc-123-def",
      "event": "push",
      "ref": "refs/heads/main",
      "sha": "abc123def456",
      "build_triggered": true
    }
  ]
}
```

**`POST /api/circuit-breaker/reset` レスポンス例：**

```json
{ "message": "Circuit breaker reset", "open": false, "consecutive_failures": 0 }
```

サーキットブレーカーが既に閉じている（`open: false`）場合も同じレスポンスを返す（冪等）。

**`GET /api/branch-config` レスポンス例：**

```json
{
  "source": "file",
  "branches": [
    {
      "branch": "main",
      "target_file": "docs",
      "sha_file": "/opt/adlaire-builder/.last_sha",
      "src": "/opt/adlaire-builder/repo/docs",
      "out": "/opt/adlaire-builder/dist/site",
      "deploy_targets": [
        { "host": "192.0.2.1", "user": "deploy", "dest_dir": "/var/www/html/" }
      ]
    }
  ]
}
```

`source` は設定の出所を示す。`.branch_config` ファイルが存在する場合は `"file"`、存在しない場合（`BRANCH_TARGETS` デフォルト値を使用）は `"default"` を返す。

**`POST /api/branch-config` リクエスト / レスポンス：**

```json
// リクエスト（GET /api/branch-config の branches と同一形式）
{
  "branches": [
    {
      "branch": "main",
      "target_file": "docs",
      "sha_file": "/opt/adlaire-builder/.last_sha",
      "src": "/opt/adlaire-builder/repo/docs",
      "out": "/opt/adlaire-builder/dist/site",
      "deploy_targets": [
        { "host": "192.0.2.1", "user": "deploy", "dest_dir": "/var/www/html/" }
      ]
    }
  ]
}
// レスポンス: 200
{ "message": "Branch config updated", "branches_count": 1 }
```

`branches` が空配列 `[]` の場合は `.branch_config` ファイルを削除し、`BRANCH_TARGETS` のデフォルト値に戻す（`source: "default"` に戻る）。変更は次回ポーリング周回から反映される。

**`POST /api/notify/weekly-summary` レスポンス例：**

```json
{ "message": "Weekly summary sent", "period": "2026-09-08/2026-09-14", "success_count": 12, "failure_count": 1, "success_rate": 92.3 }
```

`on: ["weekly_summary"]` 設定の Webhook 宛先がない場合は `422 Unprocessable Entity` を返す。

**`.build_logs/{id}.json` 追加フィールド（ビルド所要時間・コミット情報・サイズ警告）：**

```json
{
  "id": "b20260915100000",
  "started_at": "2026-09-15T09:59:18Z",
  "finished_at": "2026-09-15T10:00:00Z",
  "duration_seconds": 42,
  "status": "success",
  "commit_sha": "abc123def456",
  "commit_message": "fix: typo in §4.3 description",
  "commit_author": "Kazuhiro Kurata",
  "commit_at": "2026-09-15T09:58:00Z",
  "size_warn": false,
  "broken_links": 0,
  "heading_skips": 0,
  "reading_time": 87
}
```

`started_at` は `pipeline.sh` 実行開始時刻、`finished_at` は完了（または失敗）時刻。`duration_seconds` は整数（小数点以下切り捨て）。
`commit_sha` / `commit_message` / `commit_author` / `commit_at` はコミット情報取得 API（§13）の結果を記録する。API 失敗時は `null`。
`size_warn` は出力サイト合計サイズが `OUTPUT_SIZE_WARN_MB` 超過時 `true`、それ以外 `false`。`OUTPUT_SIZE_WARN_MB = 0` の場合は常に `false`。

**`GET /api/diagnostics` レスポンス例：**
```json
{
  "checked_at": "2026-09-15T10:00:00Z",
  "items": [
    { "name": "pat",         "status": "ok",   "message": "PAT is valid" },
    { "name": "github_api",  "status": "ok",   "message": "GitHub API reachable" },
    { "name": "output_file", "status": "ok",   "message": "Output file exists (2.0 MB)" },
    { "name": "systemd",     "status": "ok",   "message": "adlaire-ci.service is active" },
    { "name": "webhook",     "status": "warn", "message": "Webhook URL not configured" }
  ]
}
```

`status` の有効値：`"ok"` | `"warn"` | `"error"`。

**`GET /api/rate-limit` レスポンス例：**
```json
{ "limit": 5000, "remaining": 4823, "reset_at": "2026-09-15T11:00:00Z", "used": 177 }
```

**`GET /api/disk-usage` レスポンス例：**
```json
{
  "build_logs_bytes": 10485760,
  "build_logs_count": 42,
  "build_logs_archive_bytes": 2097152,
  "build_logs_archive_count": 18,
  "output_file_bytes": 2048576,
  "total_bytes": 14631448
}
```

`build_logs_bytes` / `build_logs_count` は `.build_logs/{id}.json` のみを集計する。`build_logs_archive_bytes` / `build_logs_archive_count` は `.build_logs/archive/{id}.json.gz` のみを集計する。`total_bytes` は通常 build log、archive build log、出力サイトの合計 bytes とする。

**`GET /api/config-log` レスポンス例：**
```json
{ "log": [
    { "at": "2026-09-15T10:00:00Z", "type": "server_config", "diff": { "log_max_lines": [500, 1000] }, "diff_text": "- log_max_lines: 500\n+ log_max_lines: 1000" },
    { "at": "2026-09-14T18:00:00Z", "type": "notify_config", "diff": { "enabled": [false, true] },    "diff_text": "- enabled: false\n+ enabled: true" },
    { "at": "2026-09-13T12:00:00", "type": "repo_config",   "diff": { "branch": ["main", "develop"] }, "diff_text": "- branch: main\n+ branch: develop" }
]}
```

`type` の有効値：`"server_config"` | `"notify_config"` | `"repo_config"`
`diff` の形式：`{ フィールド名: [変更前, 変更後] }`。設定変更時に `.config_log` へ追記する。
`diff_text`：`diff` を `"- key: old_value\n+ key: new_value"` 形式の文字列に変換したフィールド。複数フィールド変更時は行を連結する。設定変更時に `diff` と同時に記録する。

**`GET /api/history/{id}/comment` レスポンス例：**
```json
{ "id": "b20260914183000", "comment": "ネットワーク障害による失敗。再ビルド済み。", "updated_at": "2026-09-14T19:00:00Z" }
```

コメント未設定時は `"comment": null`。コメントは `.build_logs/{id}.json` の `comment` フィールドに保存する。

**`POST /api/history/{id}/comment` リクエスト / レスポンス：**
```json
// リクエスト
{ "comment": "ネットワーク障害による失敗。再ビルド済み。" }
// レスポンス: 200
{ "message": "Comment saved" }
```

**`POST /api/history/{id}/flag` リクエスト / レスポンス：**
```json
// リクエスト
{ "flagged": true }
// レスポンス: 200
{ "message": "Flag updated" }
```

フラグは `.build_logs/{id}.json` の `flagged` フィールドに保存する。

**`POST /api/history/{id}/tags` リクエスト / レスポンス：**
```json
// リクエスト
{ "tags": ["release", "hotfix"] }
// レスポンス: 200
{ "message": "Tags updated" }
```

タグは `.build_logs/{id}.json` の `tags` フィールド（`string[]`）に保存する。空配列 `[]` を指定するとタグをすべて削除する。

**`GET /api/tokens` レスポンス例：**
```json
{ "tokens": [
    { "id": "tok000001", "label": "監視用", "scopes": ["read"], "created_at": "2026-09-15T10:00:00Z", "last_used_at": "2026-09-15T11:00:00Z", "expires_at": null, "revoked_at": null }
]}
```

**`POST /api/tokens` リクエスト / レスポンス：**
```json
// リクエスト
{ "label": "監視用", "scopes": ["read"], "expires_at": null }
// レスポンス: 201
{ "id": "tok000001", "token": "act_...", "label": "監視用", "scopes": ["read"], "created_at": "2026-09-15T10:00:00Z", "expires_at": null }
```

`token` はレスポンス時のみ返却し、以後は取得不可。`scopes` の有効値は `read`、`trigger`、`operate`、`config`、`admin` とする。管理 session は全 API 操作を許可し、API token は指定 scope の範囲だけを許可する。トークンは `Authorization: Bearer <token>` ヘッダーで送信する。

**`DELETE /api/tokens/{id}` レスポンス例：**
```json
{ "message": "Token revoked" }
```

---

### スナップショット（14A）

ビルド成功時に出力サイトを `.snapshots/` へ自動保存する。保持世代数は `GET /api/config` の `snapshots_keep`（デフォルト `5`、`0` = 機能無効）で制御し、超過した古い世代は自動削除する。

**`GET /api/snapshots` レスポンス例：**
```json
{ "snapshots": [
    { "id": "snap001", "build_id": "b20260915100000", "saved_at": "2026-09-15T10:00:00Z", "size_bytes": 2048576 },
    { "id": "snap002", "build_id": "b20260914183000", "saved_at": "2026-09-14T18:30:00Z", "size_bytes": 2031616 }
]}
```

**`GET /api/snapshots/{id}/download`**
バイナリレスポンス。`Content-Type: application/octet-stream`、`Content-Disposition: attachment; filename="site"` を付与する。

**`DELETE /api/snapshots/{id}` レスポンス例：**
```json
{ "message": "Snapshot deleted" }
```

**`POST /api/history/{id}/rollback` リクエスト / レスポンス：**
```json
// リクエスト: なし（パスパラメーターのみ）
// レスポンス: 202
{ "message": "Rollback started", "build_id": "b20260914183000" }
```

- `.snapshots/{id}/` が存在しない場合は `404 Not Found` を返す
- ロールバックは非同期で SSH 転送を実行する（`running: true` 中は `409 Conflict` を返す）
- 転送成功時は `.build_history` に `trigger: "rollback"` のエントリを追記する

**スナップショット保存・削除固定契約：**

| 処理 | 固定仕様 |
|------|----------|
| snapshot id | §22.0e.2 の `snap{YYYYMMDDHHmmss}` 形式。build id を使い回さない。 |
| 保存対象 | 出力サイトディレクトリ配下の通常ファイルだけ。symlink、socket、device、隠し一時ファイルは保存対象外。 |
| archive 形式 | `.snapshots/{snapshot_id}/site.tar.gz` と `.snapshots/{snapshot_id}/meta.json` を作成する。 |
| `meta.json` | `id`、`build_id`、`saved_at`、`size_bytes`、`file_count`、`output_sha256` を必須 key とする。 |
| 世代削除 | 新 snapshot 作成成功後に `saved_at` 昇順で超過分だけ削除する。削除失敗は build 成功を失敗へ反転しないが、WARN log に固定コード `SNAPSHOT_PRUNE_FAILED` を出す。 |
| rollback | 対象 snapshot の `site.tar.gz` を展開して転送し、新しい build id で `.build_logs/{id}.json` と `.build_history` を作成する。 |

rollback 開始時は `.build_lock` を取得し、取得できない場合は `409 {"error":"Build is running"}` を返す。`.build_lock` 取得後に `.build_state.running=true`、`current_build_id=<new_id>` を保存し、転送完了後に finalizer で `running=false` とする。rollback は queue に積まない。

---

### Webhook 受信仕様（22-W）

`POST /api/webhook` は GitHub からの push イベントを受信し、署名検証後にビルドをトリガーする。認証ヘッダー（`Authorization: Bearer`）は不要だが、`X-Hub-Signature-256` ヘッダーによる HMAC-SHA256 署名検証が必須である。

**署名検証：**
```go
// components/api.go の実装例
mac := hmac.New(sha256.New, []byte(secret))
mac.Write(body)
expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
if subtle.ConstantTimeCompare([]byte(expected), []byte(requestHeader.Get("X-Hub-Signature-256"))) != 1 {
    return 403
}
```

- Secret は `WEBHOOK_SECRET_FILE`（`/opt/adlaire-builder/.webhook_secret`）から読み込む
- Secret 未設定時（ファイル不在）は Webhook 受信を `501 Not Implemented` で拒否する

**`POST /api/webhook` リクエストヘッダー：**
```
X-GitHub-Event: push
X-Hub-Signature-256: sha256=<hmac_hex>
Content-Type: application/json
```

**`POST /api/webhook` レスポンス：**
```json
// 200: ビルドトリガー成功
{ "message": "Build triggered", "ref": "refs/heads/main" }

// 400: ペイロード不正（ref フィールド欠損等）
{ "error": "Invalid payload" }

// 403: 署名検証失敗
{ "error": "Invalid signature" }

// 409: ビルド実行中
{ "error": "Build already running" }

// 501: Secret 未設定
{ "error": "Webhook not configured" }
```

**ビルドトリガー条件：**
- `ref` フィールドが `BRANCH_TARGETS` のいずれかの `branch` と一致する push イベントのみビルドをトリガーする
- 一致する branch が存在しない場合は `200` + `{ "message": "No matching branch" }` を返す（ビルドはしない）

**イベントログ：**
署名検証成功後、受信イベントを `.webhook_events.json` に JSON Lines 形式（1行1イベント）で追記する。ビルドの実行可否によらず全受信イベントを記録する。

記録フォーマット（1行）：
```json
{"timestamp": "2026-09-15T10:00:00Z", "delivery_id": "abc-123-def", "event": "push", "ref": "refs/heads/main", "sha": "abc123def456", "build_triggered": true}
```

| フィールド | 型 | 説明 |
|-----------|-----|------|
| `timestamp` | string (ISO 8601) | イベント受信日時 |
| `delivery_id` | string | `X-GitHub-Delivery` ヘッダー値 |
| `event` | string | `X-GitHub-Event` ヘッダー値（`"push"` 等） |
| `ref` | string | ペイロードの `ref` フィールド |
| `sha` | string | ペイロードの `after` フィールド（push 後の HEAD SHA） |
| `build_triggered` | boolean | ビルドをトリガーしたか（`true` / `false`） |

**Webhook Secret 設定 API：**

| メソッド | パス | 認証 | 説明 |
|---------|------|------|------|
| `GET` | `/api/webhook-config` | 要 | Secret 設定状態（設定済み/未設定）を返す |
| `POST` | `/api/webhook-config` | 要 | `WEBHOOK_SECRET_FILE` を更新して Secret を設定する |

**`GET /api/webhook-config` レスポンス例：**
```json
{ "configured": true }
```

**`POST /api/webhook-config` リクエスト / レスポンス：**
```json
// リクエスト
{ "secret": "<新しいSecret文字列>" }
// レスポンス: 200
{ "message": "Webhook secret updated" }
```

---

### スケジュール強制再ビルド仕様（22-S）

**`POST /api/schedule/force-interval` リクエスト / レスポンス：**
```json
// リクエスト
{ "hours": 24 }
// 無効化
{ "hours": 0 }
// レスポンス: 200
{ "message": "Force build interval updated", "hours": 24 }
```

- `components/runner.go` 側の `FORCE_BUILD_INTERVAL` を動的変更する（`.server_config` に保存し、起動時に読み込む）
- `hours` は 0 以上の整数。0 で機能無効化

**`POST /api/schedule/cooldown` リクエスト / レスポンス：**
```json
// リクエスト
{ "seconds": 120 }
// 無効化
{ "seconds": 0 }
// レスポンス: 200
{ "message": "Build cooldown updated", "seconds": 120 }
```

- `components/runner.go` 側の `BUILD_COOLDOWN_SECONDS` を動的変更する（`.server_config` に保存し、起動時に読み込む）
- `seconds` は 0 以上の整数。0 で機能無効化

---

### メンテナンスモード（14B）

メンテナンスモード有効中は手動ビルド（`POST /api/build`・`POST /api/build/force`）・スケジュールビルドの両方を拒否し、`503 Service Unavailable` + `{ "error": "maintenance" }` を返す。`GET /api/health` は制限対象外とする。

**`GET /api/maintenance` レスポンス例：**
```json
{ "enabled": false, "reason": null, "since": null }
```
メンテナンス中は `{ "enabled": true, "reason": "定期メンテナンス", "since": "2026-09-15T10:00:00Z" }`。

**`POST /api/maintenance/enable` リクエスト / レスポンス：**
```json
// リクエスト
{ "reason": "定期メンテナンス" }
// レスポンス: 200
{ "message": "Maintenance mode enabled", "since": "2026-09-15T10:00:00Z" }
```

**`POST /api/maintenance/disable` レスポンス：**
```json
{ "message": "Maintenance mode disabled" }
```

**メンテナンス更新固定契約：**

| API | 入力検証 | 保存値 | no-op | 失敗時 |
|-----|----------|--------|-------|--------|
| `POST /api/maintenance/enable` | `reason` 必須、1〜500 文字、前後空白除去後空は禁止 | `enabled:true`、`reason`、`since=now` | 既に同一 reason で enabled の場合は状態を変更せず `{ "message":"No changes","since":"<existing since>" }` | 検証失敗 `422`、保存失敗 `500` |
| `POST /api/maintenance/disable` | body 禁止 | `enabled:false`、`reason:null`、`since:null` | 既に disabled の場合は状態を変更せず `{ "message":"No changes" }` | 保存失敗 `500` |

メンテナンス判定は §22.0e の共通判定順に従い、`POST /api/build`、`POST /api/build/force`、署名検証済み `POST /api/webhook`、`POST /api/history/{id}/rollback` を拒否対象とする。設定参照系 GET、認証、ログ参照、メンテナンス解除は拒否しない。

**メンテナンス判定・副作用固定契約：**

| 対象 | 判定 | 拒否時 |
|------|------|--------|
| `POST /api/build` | 認証、rate limit、入力検証後、`.build_state` 更新前に `.maintenance.enabled` を確認する。 | `503 {"error":"maintenance"}`。queue 追加、build id 採番、`.build_state` 更新を行わない。 |
| `POST /api/build/force` | SHA reset 前に `.maintenance.enabled` を確認する。 | `503 {"error":"maintenance"}`。SHA cache、queue、`.build_state` を変更しない。 |
| `POST /api/webhook` | 署名検証、payload 検証後、`.webhook_events.json` 追記前に確認する。 | `503 {"error":"maintenance"}`。event log と queue を変更しない。 |
| `POST /api/history/{id}/rollback` | snapshot 存在確認後、rollback build log 作成前に確認する。 | `503 {"error":"maintenance"}`。history、build log、pending transfer を変更しない。 |
| runner 定期起動 | lock 取得後、差分検出前に確認する。 | build せず `.build_status.json.status="skipped_maintenance"` を保存し、`.last_sha` を更新しない。 |

`POST /api/maintenance/enable` と `POST /api/maintenance/disable` の保存順は、`.maintenance` atomic write → `.config_log` 追記 → response とする。`.config_log` 追記失敗時は `500` を返し、保存済み `.maintenance` は巻き戻さない。同一状態 no-op では `.maintenance`、`.config_log`、`.audit_log` を変更しない。

**メンテナンス fixture 固定：**

| fixture | 入力 | 期待結果 |
|---------|------|----------|
| maintenance-build-deny | enabled 中に `POST /api/build` | `503`、queue 差分なし、build id なし。 |
| maintenance-force-deny | enabled 中に `POST /api/build/force` | `503`、SHA cache 差分なし。 |
| maintenance-webhook-deny | enabled 中に署名済み webhook | `503`、event log と queue 差分なし。 |
| maintenance-disable-noop | disabled 中に disable | `200 No changes`、状態ファイル差分なし。 |

---

### IP アクセス制限（14C）

`allow` に IPv4 アドレスまたは CIDR 表記のリストを設定する。空リストは制限なし（全接続許可）を意味する。制限に一致しない接続元からのリクエストは `403 Forbidden` を返す。`GET /api/health` は制限対象外とする。設定は `.access_control` に保存する。

**`GET /api/access-control` レスポンス例：**
```json
{ "allow": ["192.168.1.0/24", "10.0.0.1"] }
```
制限なしの場合: `{ "allow": [] }`

**`POST /api/access-control` リクエスト / レスポンス：**
```json
// リクエスト
{ "allow": ["192.168.1.0/24", "10.0.0.1"] }
// レスポンス: 200
{ "message": "Access control updated", "allow": ["192.168.1.0/24", "10.0.0.1"] }
```

**アクセス制限判定固定契約：**

1. 接続元 IP は `X-Forwarded-For`、`X-Real-IP` を使わず、`net/http.Request.RemoteAddr` から取得する。
2. `RemoteAddr` が parse 不能な場合は `403 {"error":"Forbidden"}` を返す。
3. `.access_control` 不在、または `allow:[]` は全許可とする。
4. CIDR は `net.ParseCIDR`、単一 IPv4 は `net.ParseIP` で判定する。
5. 判定失敗時は認証処理より前に `403` を返し、password や token 検証を実行しない。

`POST /api/access-control` は正規化後の `allow` 配列が既存値と一致する場合、`.access_control` と `.config_log` を変更せず `{ "message":"No changes","allow":[...] }` を返す。

**アクセス制御更新・拒否固定契約：**

| 項目 | 仕様 |
|------|------|
| 判定順 | path / method 判定後、body parse 前、認証前に実行する。拒否時は password、session token、API token、rate limit state を検証または更新しない。 |
| 対象外 | `GET /api/health` だけを対象外とする。静的 admin UI、SDK JS、その他 `/api/` 以外の配信は本契約の対象外。 |
| allow 正規化 | 重複除去、辞書順 sort、単一 IPv4 は canonical 文字列、CIDR は `IP/mask` 表記へ正規化する。 |
| IPv6 | 初期実装では保存不可。IPv6 literal または IPv6 CIDR は `422`。 |
| private / public | private address に限定しない。入力が IPv4 または IPv4 CIDR として妥当なら保存可能。 |
| 保存順 | `.access_control` atomic write → `.config_log` 追記 → response。 |
| `.config_log` 失敗 | `500`。保存済み `.access_control` は巻き戻さない。 |
| 破損時 | §22.0a に従い初期値で再生成し、制限なしとして扱う。 |

**アクセス制御 fixture 固定：**

| fixture | 入力 | 期待結果 |
|---------|------|----------|
| access-allow-empty | `.access_control.allow=[]` | 任意 IP の API が認証処理へ進む。 |
| access-deny-before-auth | allow 不一致 IP で `POST /api/login` | `403`、`.access_log`、`.audit_log`、rate state 差分なし。 |
| access-normalize | 重複 allow を保存 | sort / 重複除去後の配列を返し `.config_log` 記録。 |
| access-ipv6-reject | IPv6 literal を保存 | `422`、状態差分なし。 |

---

### ビルドフック（14E）

ビルド実行の直前（`pre`）・直後（`post`）に事前登録したコマンド引数配列を実行する。フック設定は `.hooks` に保存する。外部入力文字列をシェルへ渡す実装は禁止し、Go 標準ライブラリ `os/exec` の `exec.CommandContext(args[0], args[1:]...)` で実行する。

- `pre` フックが失敗（`exit_code != 0`）し `abort_on_failure: true` の場合、ビルドを中断しステータスを `hook_error` とする。
- `post` フックは `abort_on_failure` 設定に関わらずビルド結果（`success` / `failure`）を変更しない。
- フックの実行ログは `.build_logs/{build_id}_hook_{id}.json` に保存する。

**`GET /api/hooks` レスポンス例：**
```json
{ "hooks": [
    { "id": "h20260915100500", "phase": "pre",  "command_args": ["echo", "build start"], "enabled": true, "abort_on_failure": true, "timeout_seconds": 300, "created_at": "2026-09-15T10:05:00Z" },
    { "id": "h20260915100600", "phase": "post", "command_args": ["echo", "build end"],   "enabled": true, "abort_on_failure": false, "timeout_seconds": 300, "created_at": "2026-09-15T10:06:00Z" }
]}
```

**`POST /api/hooks` リクエスト / レスポンス：**
```json
// リクエスト
{ "phase": "pre", "command_args": ["echo", "build start"], "abort_on_failure": true }
// レスポンス: 201
{ "id": "h20260915100500", "phase": "pre", "command_args": ["echo", "build start"], "enabled": true, "abort_on_failure": true, "timeout_seconds": 300, "created_at": "2026-09-15T10:05:00Z" }
```

`phase` の有効値は `"pre"` または `"post"`。`command_args[0]` は絶対パス、または `PATH` 解決可能なコマンド名とする。`command_args` に空文字、NUL 文字、改行を含めてはならない。

**`DELETE /api/hooks/{id}` レスポンス：**
```json
{ "message": "Hook deleted" }
```

**`GET /api/hooks/{id}/log` レスポンス例：**
```json
{ "id": "h20260915100500", "runs": [
    { "build_id": "b20260915100000", "ran_at": "2026-09-15T10:00:00Z", "exit_code": 0, "output": "build start\n" },
    { "build_id": "b20260914183000", "ran_at": "2026-09-14T18:30:00Z", "exit_code": 1, "output": "Error: command not found\n" }
]}
```
`runs` は直近 20 件を返す（新しい順）。

**フック実行・保存固定契約：**

| 項目 | 仕様 |
|------|------|
| 実行順 | `phase` ごとに `.hooks.hooks` の配列順。`pre` は pipeline 前、`post` は pipeline/deploy/snapshot 後。 |
| disabled | `enabled:false` は読み飛ばし、hook log を作成しない。 |
| timeout | `timeout_seconds` 超過時は process を kill し、`exit_code:null`、`timed_out:true` として hook log を保存する。 |
| stdout/stderr | 最大各 10000 文字。超過分は末尾切り捨て、`truncated:true` を保存する。 |
| hook log | `.build_logs/{build_id}_hook_{hook_id}.json` に JSON object で保存し、同一 build/hook の再実行時は上書きせず `runs` へ追記する。 |
| pre abort | `pre` 失敗かつ `abort_on_failure:true` の場合、pipeline を実行せず build status を `hook_error` とする。 |
| secret mask | stdout、stderr、保存済み output、server log、通知 payload へ保存する前に runner の secret mask を適用する。 |
| log 保存失敗 | `pre` hook では build 本体を開始せず `hook_error`。`post` hook では build 結果を維持し、server log に `HOOK_LOG_WRITE_FAILED` を出す。 |

**`.hooks` record schema：**

| キー | 型 | 必須 | 許容値 |
|------|----|------|--------|
| `id` | string | 必須 | §22.0e.2 の hook id。 |
| `phase` | string | 必須 | `"pre"` または `"post"`。 |
| `command_args` | string[] | 必須 | 1〜20 件。各値は NUL、LF、CR 禁止。 |
| `enabled` | boolean | 必須 | boolean。作成時 `true` 固定。 |
| `abort_on_failure` | boolean | 必須 | boolean。 |
| `timeout_seconds` | integer | 必須 | 1〜3600。省略時 300。 |
| `created_at` | string | 必須 | UTC ISO 8601。 |

`.hooks` に未知 key、必須 key 不足、型不一致、不正 phase、不正 command、重複 id がある場合、`GET /api/hooks`、`POST /api/hooks`、`DELETE /api/hooks/{id}` は `500 {"error":"Internal server error"}` を返す。破損内容、command_args の secret らしき値、stdout/stderr は response と log に出さない。

**hooks API 更新順：**

| API | 更新順 | 失敗時 |
|-----|--------|--------|
| `POST /api/hooks` | 入力検証 → `.hooks` lock → id 採番 → record append → `.hooks` atomic write → `.config_log` 追記 → response | `.config_log` 失敗時は `500`。追加済み record は巻き戻さない。 |
| `DELETE /api/hooks/{id}` | path id 検証 → `.hooks` lock → 対象存在確認 → record 削除 → `.hooks` atomic write → `.config_log` 追記 → response | 対象不在は `404`。`.config_log` 失敗時は `500`、削除済み record は巻き戻さない。 |
| `GET /api/hooks/{id}/log` | path id 検証 → `.hooks` で存在確認 → `.build_logs/*_hook_{id}.json` を新しい順で最大 20 件読込 → response | hook 不在は `404`。個別 hook log 破損はその file を除外し、server log に固定コードを出す。 |

hook log JSON は `{ "hook_id", "build_id", "phase", "started_at", "finished_at", "duration_seconds", "exit_code", "timed_out", "stdout", "stderr", "truncated" }` を必須 key とする。`GET /api/hooks/{id}/log` の `output` は `stdout + stderr` をこの順で連結した表示用互換値とし、保存時点で secret mask 済みの値だけを返す。

**hooks fixture 固定：**

| fixture | 入力 | 期待結果 |
|---------|------|----------|
| hook-pre-success | pre hook exit 0 | pipeline 実行、hook log 保存、secret mask 済み。 |
| hook-pre-abort | pre hook exit 1 / abort true | pipeline 未実行、status `hook_error`、history に `failure_category:"hook_error"`。 |
| hook-pre-warn | pre hook exit 1 / abort false | build 継続、hook log に exit code。 |
| hook-post-failure | build success 後 post hook failure | build success 維持、hook log 保存。 |
| hook-timeout | timeout 超過 | process kill、`timed_out:true`、`exit_code:null`。 |
| hook-log-write-failure | pre hook log 保存失敗 | build 本体未実行、`hook_error`。 |
| post failure | build status を変更しない。WARN log と hook log だけを残す。 |

保存する hook log file は 1 実行 1 JSON object とし、`hook_id`、`build_id`、`phase`、`started_at`、`finished_at`、`duration_seconds`、`exit_code`、`timed_out`、`stdout`、`stderr`、`truncated` を必須 key とする。`GET /api/hooks/{id}/log` は複数 file を集約し、response の `runs[]` へ `build_id`、`ran_at`、`exit_code`、`output` を返す。

---

### カスタムアラートルール（15A）

`GET /api/dashboard` 取得時にルールを評価し、条件を満たすものを `alerts` 配列へ自動追加する。`.alert_rules` に保存する。

対応メトリクス：`success_rate_7d`（7日間成功率 %）/ `avg_duration_seconds`（7日間平均ビルド時間 秒）/ `last_build_age_hours`（最終ビルドからの経過時間 時間）/ `disk_usage_bytes`（ディスク使用量 バイト）
対応演算子：`lt`（未満）/ `gt`（超過）/ `lte`（以下）/ `gte`（以上）

**`GET /api/alert-rules` レスポンス例：**
```json
{ "rules": [
    { "id": "r001", "metric": "success_rate_7d", "operator": "lt", "threshold": 90, "level": "warn", "message": "7日間成功率が90%を下回っています" },
    { "id": "r002", "metric": "last_build_age_hours", "operator": "gt", "threshold": 48, "level": "error", "message": "48時間以上ビルドが実行されていません" }
]}
```

**`POST /api/alert-rules` リクエスト / レスポンス：**
```json
// リクエスト
{ "metric": "avg_duration_seconds", "operator": "gt", "threshold": 120, "level": "warn", "message": "平均ビルド時間が2分を超えています" }
// レスポンス: 201
{ "id": "r003", "metric": "avg_duration_seconds", "operator": "gt", "threshold": 120, "level": "warn", "message": "平均ビルド時間が2分を超えています" }
```

**`DELETE /api/alert-rules/{id}` レスポンス：**
```json
{ "message": "Alert rule deleted" }
```

**アラート評価固定契約：**

`.alert_rules` 不在時は `rules:[]` と扱う。`GET /api/dashboard` は rule 配列順で評価し、条件一致した rule だけ `alerts` へ追加する。alert object は `id`、`level`、`message`、`metric`、`value`、`threshold` を必須 key とする。評価に必要な metric が算出不能な rule は alert 化せず、WARN log `ALERT_METRIC_UNAVAILABLE` を出す。

`POST /api/alert-rules` は `metric`、`operator`、`threshold`、`level`、`message` の正規化後値が既存 rule と一致する場合、`409 {"error":"Conflict"}` を返す。`DELETE` は対象 id 不在時 `404` とし、部分削除は行わない。

---

### 自動タグ付けルール（15B）

ビルド完了時に条件式を評価し、マッチしたルールのタグを `.build_history` のエントリへ自動追記する（手動タグと共存する）。`.tag_rules` に保存する。

条件式で使用可能な変数：`status`（`"success"` / `"failure"`）/ `duration_seconds`（整数）/ `trigger`（`"polling"` / `"force_interval"` / `"manual"` / `"webhook"` / `"retry_pending_transfer"` / `"startup_config_integrity"` / `"rollback"` / `"local_watch"` / `"approval"`）
演算子：`==`・`!=`・`>`・`<`・`>=`・`<=`

**`GET /api/tag-rules` レスポンス例：**
```json
{ "rules": [
    { "id": "t001", "condition": "status == 'failure'", "tags": ["要確認"] },
    { "id": "t002", "condition": "duration_seconds > 120", "tags": ["低速"] },
    { "id": "t003", "condition": "trigger == 'force_interval'", "tags": ["強制実行"] }
]}
```

**`POST /api/tag-rules` リクエスト / レスポンス：**
```json
// リクエスト
{ "condition": "status == 'success'", "tags": ["green"] }
// レスポンス: 201
{ "id": "t004", "condition": "status == 'success'", "tags": ["green"] }
```

**`DELETE /api/tag-rules/{id}` レスポンス：**
```json
{ "message": "Tag rule deleted" }
```

**自動タグ評価固定契約：**

タグ評価は build log の最終 status / duration / trigger が確定した後、`.build_history` 追記前に行う。手動タグと自動タグが重複した場合は 1 件に正規化し、既存順を保持したうえで自動タグを末尾へ追加する。条件式 parse 失敗を含む破損 rule がある場合、その rule を無視せず build を `failure` にし、ERROR log `TAG_RULE_INVALID` を出す。

`POST /api/tag-rules` は同一 `condition` と同一 `tags` の rule が存在する場合、`409 {"error":"Conflict"}` を返す。

---

### 出力サイトチェックサム（15C）

ビルド成功時に出力サイト配下の全通常ファイルから manifest SHA-256 を算出し、`.build_history` の該当エントリに `output_sha256` として記録する。manifest は `relative_path + "\n" + file_sha256 + "\n"` を相対パス昇順で連結した文字列とし、その SHA-256 hex を `output_sha256` とする。

**`GET /api/output-meta` レスポンス変更（`sha256` / build meta フィールド追加）：**
```json
{
  "size_bytes": 2048576,
  "mtime": "2026-09-15T10:00:00Z",
  "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
  "build_id": "b20260915100000",
  "commit_sha": "abc1234",
  "build_at": "2026-09-15T10:00:00Z"
}
```

**`GET /api/history/{id}/log` レスポンス変更（`output_sha256` フィールド追加）：**
```json
{ "id": "b001", "build_at": "2026-09-15T10:00:00Z", "sha": "abc123", "status": "success",
  "output_size_bytes": 2048576, "output_sha256": "e3b0c44298fc1c149afbf4c8996fb924...",
  "trigger": "polling", "comment": null, "flagged": false, "tags": ["release"],
  "lines": ["2026-09-15T10:00:00Z [INFO] Build start", "..."] }
```

**`POST /api/verify-output` レスポンス例：**
```json
// 一致時
{ "match": true, "expected": "e3b0c44298fc1c149afbf4c8996fb924...", "actual": "e3b0c44298fc1c149afbf4c8996fb924..." }
// 不一致時
{ "match": false, "expected": "e3b0c44298fc1c149afbf4c8996fb924...", "actual": "f4a2d5591c8f3a742f902e3b6f7c1c3d..." }
```
`expected` は `.build_history` の最終成功エントリに記録された `output_sha256`。現在の出力サイトが存在しない場合は `404` を返す。

**checksum 算出固定契約：**

対象 path は `/` 区切りの相対 path とし、先頭 `/`、`..`、NUL、改行を含む path は manifest 算出前に build 失敗とする。ディレクトリ、symlink、device、socket は manifest に含めない。空ディレクトリの checksum は空文字列に対する SHA-256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` とする。

`POST /api/verify-output` は状態ファイルを変更しない。`.build_history` に成功履歴がない場合、`404 {"error":"Not Found"}` を返す。

---

### ビルドパイプライン設定（15D）

`components/runner.go` がビルド実行時に `.pipeline_config` を読み込み、`components/builder.go` の呼び出しに `extra_args`・`env` を適用する。`.pipeline_config` に保存する。

**`GET /api/pipeline-config` レスポンス例：**
```json
{ "extra_args": ["--verbose"], "env": { "DEBUG": "1" } }
```
初期値（未設定時）：`{ "extra_args": [], "env": {} }`

**`POST /api/pipeline-config` リクエスト / レスポンス：**
```json
// リクエスト
{ "extra_args": ["--verbose"], "env": { "DEBUG": "1" } }
// レスポンス: 200
{ "message": "Pipeline config updated" }
```

`POST /api/pipeline-config` は `.pipeline_config` 全体置換とし、部分更新を許可しない。`extra_args` または `env` のいずれかが欠ける場合は `422`。正規化後値が既存値と一致する場合は `.pipeline_config` と `.config_log` を変更せず `{ "message":"No changes" }` を返す。

runner は build 開始後、builder command 組み立て直前に `.pipeline_config` を 1 回だけ読む。読込不能または schema 不正は build を開始せず `failure` とし、SHA cache を更新しない。

---

### 運用ノート（15E）

システム全体の運用メモを Markdown テキストで保存・取得する。`.notes` に保存する。認証必須。

**`GET /api/notes` レスポンス例：**
```json
{ "content": "# 運用メモ\n定期メンテナンス: 毎週日曜 2:00〜4:00\nPAT 更新期限: 2026-12-01", "updated_at": "2026-09-15T10:00:00Z" }
```
初回（未作成）時：`{ "content": "", "updated_at": null }`

**`POST /api/notes` リクエスト / レスポンス：**
```json
// リクエスト
{ "content": "# 運用メモ\n定期メンテナンス: 毎週日曜 2:00〜4:00" }
// レスポンス: 200
{ "message": "Notes updated", "updated_at": "2026-09-15T10:00:00Z" }
```

`.notes` は UTF-8 text として保存し、JSON ではない。`POST /api/notes` の `content` は 0〜100000 文字、NUL 禁止とする。保存時は本文をそのまま `.notes` に atomic write し、更新時刻は `.config_log` の `at` を返す。既存本文と一致する場合は `.notes` と `.config_log` を変更せず `{ "message":"No changes","updated_at":null }` を返す。

---

### メール通知（SMTP）（16B）

Webhook に加えてメールでビルド結果を通知できる機能。SMTP 接続設定は `.smtp_config` に、パスワードは `.smtp_secret`（パーミッション 600）に分離して保存する。`GET /api/notify-config` のレスポンスに `email` セクションを追加する。

**`GET /api/smtp-config` レスポンス例：**
```json
{ "host": "smtp.example.com", "port": 587, "user": "notify@example.com", "tls": true, "from": "notify@example.com", "to": ["ops@example.com"], "on": ["failure"], "enabled": true, "password_set": true }
```
`password_set`：`.smtp_secret` が存在するかを真偽値で返す。パスワード本体は返却しない。
未設定時：`{ "host": null, "port": 587, "user": null, "tls": true, "from": null, "to": [], "on": [], "enabled": false, "password_set": false }`

**`POST /api/smtp-config` リクエスト / レスポンス：**
```json
// リクエスト（変更するフィールドのみ指定可。password フィールドは省略可能）
{ "host": "smtp.example.com", "port": 587, "user": "notify@example.com", "password": "s3cr3t", "tls": true, "from": "notify@example.com", "to": ["ops@example.com"], "on": ["failure"], "enabled": true }
// レスポンス: 200
{ "message": "SMTP config updated" }
```
`on` の有効値：`"start"` / `"success"` / `"failure"`

**`POST /api/smtp-test` レスポンス例：**
```json
// 成功時
{ "result": "success", "message": "Test email sent to ops@example.com" }
// 失敗時
{ "result": "failure", "message": "Connection refused: smtp.example.com:587" }
```
SMTP 未設定または `enabled: false` の場合は `422` を返す。

**SMTP 更新・送信固定契約：**

| 処理 | 仕様 |
|------|------|
| 更新 | `.smtp_config` → 必要時 `.smtp_secret` → `.config_log` の順に書く。途中失敗時は未処理ファイルを書かない。 |
| 削除 | `password:null` は `.smtp_secret` 削除。削除対象が不在なら no-op。 |
| GET | `.smtp_secret` の存在だけを `password_set` で返し、password 本体は返さない。 |
| 送信 timeout | 接続、TLS、送信全体を合計 30 秒で timeout する。 |
| 送信 log | `.notify_log` に `type:"smtp_test"`、`result`、`message`、`at` を追記する。password は記録しない。 |

`enabled:true` にする場合は `host`、`port`、`from`、`to` 1 件以上を必須とする。`POST /api/smtp-test` は `enabled:false`、宛先なし、secret 必須構成で `.smtp_secret` 不在のいずれも `422 {"error":"SMTP not configured"}` を返す。

**SMTP 更新詳細：**

| ケース | `.smtp_config` | `.smtp_secret` | `.config_log` | response |
|--------|----------------|----------------|---------------|----------|
| config のみ変更 | 保存 | 変更なし | mask 済み diff 追記 | `200 {"message":"SMTP config updated"}` |
| password 追加 / 変更 | 保存 | mode `0600` で atomic write | password は `"***"` で追記 | `200` |
| `password:null` | 保存 | 存在すれば削除 | password は `"***"` で追記 | `200` |
| 完全 no-op | 変更なし | 変更なし | 追記なし | `200 {"message":"No changes"}` |
| `.smtp_secret` 書込失敗 | 必要なら `.smtp_config` 保存済み | 失敗 | 追記なし | `500` |
| `.config_log` 追記失敗 | 保存済み | 保存または削除済み | 失敗 | `500` |

`POST /api/smtp-test` は `.smtp_config` と `.smtp_secret` を読み、送信成功 / 失敗のどちらも `.notify_log` へ追記してから response を返す。`.notify_log` 追記失敗時は `500` を返す。SMTP password、認証失敗時の server response に含まれる credential 断片、接続 URL の userinfo は `message` と log に含めず固定文言へ置換する。

**SMTP fixture 固定：**

| fixture | 入力 | 期待結果 |
|---------|------|----------|
| smtp-save-password | password 付き保存 | `.smtp_secret` mode `0600`、GET は `password_set:true`、log は `"***"`。 |
| smtp-delete-password | `password:null` | `.smtp_secret` 削除、password 平文なし。 |
| smtp-noop | 同一 config / password 未指定 | 状態差分なし、`.config_log` 追記なし。 |
| smtp-test-success | 設定済み test | `.notify_log` に success、response success。 |
| smtp-test-disabled | `enabled:false` | `422`、`.notify_log` 差分なし。 |
| smtp-log-failure | test 後 `.notify_log` 追記失敗 | `500`、password 平文なし。 |

**`GET /api/notify-config` への追加（`email` セクション）：**
```json
{
  "webhooks": [ { "url": "...", "label": "メイン", "enabled": true, "payload_template": null, "retry_count": 2, "retry_interval_seconds": 30, "secret": "***" } ],
  "channels": [ { "id": "n002", "type": "email", "label": "Ops", "enabled": true, "on": ["failure", "duration_anomaly"], "config": { "to": ["ops@example.com"] }, "retry_count": 2, "retry_interval_seconds": 30 } ],
  "on": ["failure"],
  "summary": { "enabled": false, "interval": "weekly", "hour": 9, "day_of_week": 1 },
  "email": { "enabled": true, "to": ["ops@example.com"], "on": ["failure"] }
}
```

---

### ビルドキューイング（16C）

ビルド実行中に `POST /api/build`・`POST /api/build/force` を受信した場合、キューに追加して順次実行する。キューの最大長は `GET /api/config` の `queue_max_size`（デフォルト `3`、`0` = キューなし）で制御する。キューが満杯の場合は `429 Too Many Requests` + `{ "error": "queue_full" }` を返す。

**`GET /api/queue` レスポンス例：**
```json
{ "queued": [
    { "id": "q20260915100100", "trigger": "manual", "queued_at": "2026-09-15T10:01:00Z", "requested_by": "admin", "priority": "normal", "created_seq": 1, "payload": { "force": false } },
    { "id": "q20260915100200", "trigger": "webhook", "queued_at": "2026-09-15T10:02:00Z", "requested_by": "webhook", "priority": "normal", "created_seq": 2, "payload": { "delivery_id": "delivery-1", "branch": "main", "sha": "0123456789abcdef0123456789abcdef01234567" } }
  ], "max_size": 3 }
```
キューが空の場合：`{ "queued": [], "max_size": 3 }`

**`DELETE /api/queue` レスポンス：**
```json
{ "message": "Queue cleared", "cleared_count": 2 }
```
実行中のビルドは停止しない（`POST /api/build/cancel` を別途使用する）。

`POST /api/config` に `queue_max_size`（整数、`0` = キューなし）追加。

**queue 処理固定契約：**

queue 追加は `.build_state` の atomic write で行い、id は §22.0e.2 の queue id とする。queue entry は FIFO を標準とし、§27.35 の優先度キューが有効な場合だけ priority を使用する。`DELETE /api/queue` は `.build_state.queued` だけを空配列にし、実行中 build、lock、history、log を変更しない。

`queue_max_size=0` の場合、実行中に受けた queue 対応 API は `429 {"error":"queue_full"}` を返す。`.build_state.queued` の件数が `queue_max_size` 以上の場合も同じ body とする。

**queue entry schema：**

| キー | 型 | 必須 | 許容値 |
|------|----|------|--------|
| `id` | string | 必須 | §22.0e.2 の queue id。 |
| `trigger` | string | 必須 | `"manual"`、`"webhook"`、`"approval"`。 |
| `queued_at` | string | 必須 | UTC ISO 8601。 |
| `requested_by` | string | 必須 | `"admin"`、`"webhook"`、`"approval"`、API token id。 |
| `priority` | string | 必須 | §27.35 の値。未指定作成時は `"normal"`。 |
| `created_seq` | integer | 必須 | 1 以上。既存最大 + 1。 |
| `payload` | object | 必須 | trigger ごとの固定 payload。未使用時は `{}`。 |

`payload` は trigger ごとに以下を許可する。未知 key は `422`、runner 読込時は queue entry 破損として当該 entry を処理せず ERROR ログに記録する。

| trigger | payload |
|---------|---------|
| `manual` | `{ "force": boolean }`。 |
| `webhook` | `{ "delivery_id": string, "branch": string, "sha": string }`。 |
| `approval` | `{ "approval_id": string, "branch": string, "sha": string, "target": string }`。 |

**queue 更新順：**

| 操作 | 更新順 | 失敗時 |
|------|--------|--------|
| API queue 追加 | `.build_state` lock → 最新 state 読込 → running / max_size / 重複確認 → id と created_seq 採番 → atomic write → response | write 失敗は `500`。queue 追加なし。 |
| webhook queue 追加 | event 検証 → `.webhook_events.json` 追記 → `.build_state` lock → queue append → response | event 追記前の失敗は queue なし。queue 追加失敗は event result を `error` にできる場合だけ追記し、response は `500`。 |
| approval approve queue 追加 | `.approval_queue` lock → pending 確認 → `.build_state` lock → queue append → approval status `approved` 追記 → response | queue full は `429`、approval は pending のまま。approval status 追記失敗時は `500`、queue 追加済み entry は巻き戻さない。 |
| runner 取り出し | `.build_state` lock → §27.35 の順で 1 件選択 → selected entry 削除 → `running=true` と `current_build_id` 設定 → atomic write | write 失敗は build 開始なし、lock を解放し終了コード `1`。 |
| queue clear | `.build_state` lock → `queued=[]` → atomic write → `.config_log` 追記 → response | `.config_log` 失敗時は `500`。cleared queue は巻き戻さない。 |

重複判定は `trigger` と `payload` の正規化 JSON が一致する waiting entry を対象とする。重複時は新規 entry を追加せず `200 {"message":"Already queued","queued":true,"queue_id":"<existing>"}` を返す。`force=true` の manual entry は `force=false` と別 entry として扱う。

**queue fixture 固定：**

| fixture | 入力 | 期待結果 |
|---------|------|----------|
| queue-add-running | running 中に manual build | queue append、created_seq 最大 + 1。 |
| queue-duplicate | 同一 manual payload を再投入 | 新規追加なし、既存 queue_id を返す。 |
| queue-full | max_size 到達 | `429 {"error":"queue_full"}`、差分なし。 |
| queue-clear | waiting 2 件で `DELETE /api/queue` | `cleared_count=2`、running/current_build_id 維持。 |
| queue-runner-take | urgent と normal が混在 | urgent を削除し running に設定、他 entry 維持。 |

---

### ダッシュボードウィジェットカスタマイズ（16D）

ダッシュボードに表示するウィジェットの種類・順序を設定できる機能。`.dashboard_layout` に保存する。

有効なウィジェット識別子：`status`（ビルド状態）/ `stats`（統計サマリー）/ `schedule`（次回実行）/ `alerts`（アラート）/ `disk`（ディスク使用量）/ `rate_limit`（GitHub API レート制限）/ `snapshots`（スナップショット件数）/ `maintenance`（メンテナンス状態）/ `queue`（キュー状態）

**`GET /api/dashboard-layout` レスポンス例：**
```json
{ "widgets": ["status", "alerts", "stats", "schedule", "disk"] }
```
未設定時はデフォルト順（全ウィジェット）を返す。

**`POST /api/dashboard-layout` リクエスト / レスポンス：**
```json
// リクエスト
{ "widgets": ["status", "alerts", "schedule", "stats"] }
// レスポンス: 200
{ "message": "Dashboard layout updated" }
```
`widgets` に未知の識別子が含まれる場合は `422` を返す。

**dashboard layout 固定契約：**

既定 widget 順は `["status","stats","schedule","alerts","disk","rate_limit","snapshots","maintenance","queue"]` とする。`POST /api/dashboard-layout` は `widgets` 全体置換のみ許可し、空配列、重複、未知 id は `422`。正規化後値が既存値と一致する場合は `.dashboard_layout` と `.config_log` を変更せず `{ "message":"No changes" }` を返す。

---

**`POST /api/logs/cleanup` レスポンス例：**
```json
{ "message": "Cleanup completed", "deleted_count": 12 }
```

`log_retention_days` が `0` の場合は削除せず `deleted_count: 0` を返す。

**`GET /api/history/export` レスポンス：**

`Content-Type: application/json` で返却される。

```json
{ "exported_at": "2026-09-15T10:00:00Z", "history": [
    { "id": "b20260915100000", "build_at": "2026-09-15T10:00:00Z", "sha": "abc123", "status": "success", "trigger": "polling", "duration_seconds": 42, "flagged": false, "tags": ["release"], "comment": "" },
    { "id": "b20260914183000", "build_at": "2026-09-14T18:30:00Z", "sha": "def456", "status": "failure", "trigger": "manual", "duration_seconds": 7, "flagged": true, "tags": [], "comment": "ネットワーク障害による失敗" }
]}
```

**`POST /api/repo-config` リクエスト / レスポンス：**
```json
// リクエスト（変更するフィールドのみ指定可）
{ "owner": "fqwink", "repo": "Adlaire-Design-System", "branch": "main", "target_file": "docs" }
// レスポンス: 200
{ "message": "Repo config updated" }
```

設定は `.repo_config` に保存し、`GET /api/repo-info` もこのファイルを参照する。

`trigger` の有効値：§13 の固定値（`"polling"`、`"force_interval"`、`"manual"`、`"webhook"`、`"retry_pending_transfer"`、`"startup_config_integrity"`、`"rollback"`、`"local_watch"`、`"approval"`）

**エラーレスポンス形式：**

| ステータス | 条件 | レスポンス |
|-----------|------|-----------|
| `401` | 認証失敗・セッション期限切れ | `{"error": "Unauthorized"}` |
| `404` | リソース不存在 | `{"error": "Not Found"}` |
| `409` | 状態競合（`POST /api/build/cancel` でビルド未実行時等） | `{"error": "No build is running"}` |
| `422` | Webhook 送信失敗（`POST /api/notify-test`） | `{"error": "Webhook delivery failed"}` |
| `429` | ビルドキューが満杯（`POST /api/build` / `POST /api/build/force`） | `{"error": "queue_full"}` |
| `500` | サーバー内部エラー | `{"error": "Internal Server Error"}` |
| `503` | メンテナンスモード中のリクエスト（`POST /api/build` 等） | `{"error": "maintenance"}` |

---

## 23. JavaScript SDK 仕様

本節は、`admin/adlaire-ci-sdk.js` に関する仕様である。

**ファイル：** `admin/adlaire-ci-sdk.js`（単一ファイル、外部依存なし）
**モジュール形式：** ES Module（`import` / `export`）

**SDK 実行環境契約：**

| 項目 | 仕様 |
|------|------|
| JavaScript | ECMAScript 2022 以上を前提とする。transpile、bundle、polyfill は標準仕様に含めない。 |
| module | `admin/adlaire-ci-sdk.js` は ES Module とし、`export { AdlaireCI, AdlaireCIError }` を必須 export とする。default export は定義しない。 |
| browser API | `fetch`、`AbortController`、`ReadableStream.getReader()`、`TextDecoder`、`URLSearchParams` が存在する browser を必須環境とする。いずれかが存在しない場合、`AdlaireCI` constructor は `TypeError("Unsupported browser runtime")` を投げる。 |
| 非 browser runtime | browser API 行の必須 API が存在しない実行環境では、runtime 名を判定分岐せず、`AdlaireCI` constructor が `TypeError("Unsupported browser runtime")` を投げる。Node.js 専用 API、npm package、bundler、polyfill による補完は行わない。 |
| 外部依存 | npm package、CDN script、framework、build tool を使用してはならない。 |
| global 汚染 | `window.AdlaireCI` 等の global 代入を行わない。標準管理ツールは ES Module import で SDK を読み込む。 |
| stream 前提 | `streamBuild()` は native `EventSource` を使用しない。Authorization header を付与できる `fetch` streaming を必須実装とする。 |

```js
class AdlaireCI {
  constructor({ baseUrl })
  // this._token でセッショントークンを管理。login() 後の全リクエストに自動付与

  login(password)                               // POST /api/login → {token?, must_change, totp_required?, ticket?}; token がある場合は this._token にセット
  loginTotp(ticket, code)                       // POST /api/login/totp → {token, must_change}; this._token にセット
  logout()                                      // POST /api/logout; this._token をクリア
  changePassword(currentPassword, newPassword)  // POST /api/change-password
  getAuditLog({ limit = 100, offset = 0, actor = null, action = null, result = null } = {}) // GET /api/audit-log → Promise<{log: AuditLogRecord[], total: number}>

  getStatus()               // GET /api/status               → Promise<StatusObject>
  triggerBuild()            // POST /api/build               → Promise<{message: string, build_id?: string, queued?: boolean}>
  getLogs(n = 100, q = '')  // GET /api/logs?n={n}&q={q}     → Promise<{lines: string[]}>
  getHistory({ page = 1, perPage = 20, trigger = null, tag = null, flagged = null } = {}) // GET /api/history?page={page}&per_page={perPage}&trigger={trigger}&tag={tag}&flagged={flagged} → Promise<HistoryPageObject>
  getSysinfo()              // GET /api/sysinfo              → Promise<SysinfoObject>
  getSchedule()             // GET /api/schedule             → Promise<ScheduleObject>
  getNotifyConfig()         // GET /api/notify-config        → Promise<NotifyConfig>
  setNotifyConfig(config)   // POST /api/notify-config       → Promise<{message: string}>
  getConfig()               // GET /api/config               → Promise<ConfigObject>
  validateConfig(config)    // POST /api/config/validate     → Promise<ConfigValidationObject>
  setConfig(config)         // POST /api/config              → Promise<{message: string, config: ConfigObject}>
  health()                  // GET /api/health               → Promise<{status: string}>
  getPatStatus()            // GET /api/pat-status           → Promise<PatStatusObject>
  getAccessLog()            // GET /api/access-log           → Promise<{log: AccessRecord[]}>
  getApiAccessLog({ limit = 100, offset = 0, method = null, path = null, status = null } = {}) // GET /api/api-access-log → Promise<{log: ApiAccessRecord[], total: number}>
  getStats(days = 7)        // GET /api/stats?days={days}    → Promise<StatsObject>
  exportLogs()              // GET /api/logs/export          → Promise<{exported_at: string, lines: string[]}>
  cleanupLogs()             // POST /api/logs/cleanup        → Promise<{message: string, deleted_count: number}>
  archiveLogs()             // POST /api/logs/archive        → Promise<{message: string, archived_count: number}>
  getRepoInfo()             // GET /api/repo-info            → Promise<RepoInfoObject>
  backup()                  // GET /api/backup               → Promise<BackupObject>
  restore(config)           // POST /api/restore             → Promise<{message: string}>
  notifyTest()              // POST /api/notify-test         → Promise<{message: string, webhook_url: string}>
  buildForce()              // POST /api/build/force         → Promise<{message: string, build_id?: string, queued?: boolean}>
  patVerify()               // POST /api/pat-verify          → Promise<PatVerifyObject>
  getHistoryLog(id)         // GET /api/history/{id}/log     → Promise<HistoryLogObject>
  cancelBuild()             // POST /api/build/cancel        → Promise<{message: string}>
  resetCircuitBreaker()     // POST /api/circuit-breaker/reset → Promise<{message: string, open: boolean, consecutive_failures: number}>
  streamBuild(onLine, onEnd) // GET /api/build/stream (SSE)  → Promise<StreamHandle>（onLine(line), onEnd({status, duration_seconds}) コールバック）
  setLogLevel(level)        // POST /api/log-level           → Promise<{message: string, level: string}>
  updatePat(token)          // POST /api/pat-update          → Promise<{message: string}>
  getDashboard()            // GET /api/dashboard            → Promise<DashboardObject>
  getNotifyLog()            // GET /api/notify-log           → Promise<{log: NotifyRecord[]}>
  getSessions()             // GET /api/sessions             → Promise<{sessions: SessionRecord[]}>
  revokeAllSessions()       // POST /api/sessions/revoke-all → Promise<{message: string, revoked_count: number}>
  getTotpStatus()           // GET /api/auth/totp-status     → Promise<TotpStatus>
  setupTotp()               // POST /api/auth/totp-setup     → Promise<{secret: string, otpauth_uri: string}>
  confirmTotp(code)         // POST /api/auth/totp-confirm   → Promise<TotpStatus>
  disableTotp(code)         // DELETE /api/auth/totp         → Promise<TotpStatus>
  setScheduleInterval(seconds) // POST /api/schedule/interval → Promise<{message: string, interval_seconds: number}>
  pauseSchedule()              // POST /api/schedule/pause   → Promise<{message: string}>
  resumeSchedule()             // POST /api/schedule/resume  → Promise<{message: string}>
  setAllowedHours(from, to)    // POST /api/schedule/allowed-hours {from, to} → Promise<{message: string, allowed_hours: {from: number, to: number}}>
  clearAllowedHours()          // POST /api/schedule/allowed-hours {from:null, to:null} → Promise<{message: string, allowed_hours: null}>
  setForceInterval(hours)      // POST /api/schedule/force-interval → Promise<{message: string, hours: number}>
  setBuildCooldown(seconds)    // POST /api/schedule/cooldown → Promise<{message: string, seconds: number}>
  searchLogs(q = '', from = '', to = '') // GET /api/logs/search?q={q}&from={from}&to={to} → Promise<SearchResult>
  getOutputMeta()              // GET /api/output-meta       → Promise<OutputMetaObject>
  getStatsTimeline(days = 30)  // GET /api/stats/timeline?days={days} → Promise<TimelineObject>
  getStatsBuildDuration(n = 10) // GET /api/stats/build-duration?n={n} → Promise<BuildDurationStats>
  getBuildTrends(n = 100)       // GET /api/stats/build-trends?n={n} → Promise<BuildTrendStats>
  getDiagnostics()             // GET /api/diagnostics       → Promise<DiagnosticsObject>
  getRateLimit()               // GET /api/rate-limit        → Promise<RateLimitObject>
  getDiskUsage()               // GET /api/disk-usage        → Promise<DiskUsageObject>
  getConfigLog()               // GET /api/config-log        → Promise<{log: ConfigLogRecord[]}>
  getApiRateLimit()            // GET /api/api-rate-limit    → Promise<ApiRateLimitPolicy>
  setApiRateLimit(policy)      // POST /api/api-rate-limit   → Promise<ApiRateLimitPolicy>
  getBranchConfig()            // GET /api/branch-config     → Promise<{source: string, branches: BranchTargetRecord[]}>
  setBranchConfig(branches)    // POST /api/branch-config    → Promise<{message: string, branches_count: number}>
  notifyWeeklySummary()        // POST /api/notify/weekly-summary → Promise<{message: string, period: string, success_count: number, failure_count: number, success_rate: number}>
  getWebhookEvents(limit = 50, offset = 0) // GET /api/webhook-events?limit={limit}&offset={offset} → Promise<{events: WebhookEventRecord[], total: number}>
  getWebhookConfig()          // GET /api/webhook-config    → Promise<{configured: boolean}>
  setWebhookConfig(secret)    // POST /api/webhook-config   → Promise<{message: string}>
  getBuildChainConfig()       // GET /api/build-chain-config → Promise<BuildChainConfig>
  setBuildChainConfig(chains) // POST /api/build-chain-config → Promise<{message: string, chains_count: number}>
  getApprovals()              // GET /api/approvals         → Promise<{approvals: ApprovalRecord[]}>
  approveBuild(id)            // POST /api/approvals/{id}/approve → Promise<{message: string, queued: boolean}>
  rejectBuild(id)             // POST /api/approvals/{id}/reject  → Promise<{message: string}>
  getHistoryComment(id)        // GET /api/history/{id}/comment  → Promise<CommentObject>
  setHistoryComment(id, comment) // POST /api/history/{id}/comment → Promise<{message: string}>
  setRepoConfig(config)        // POST /api/repo-config      → Promise<{message: string}>
  exportHistory()              // GET /api/history/export    → Promise<ExportObject>（JSON）
  setHistoryFlag(id, flagged)  // POST /api/history/{id}/flag → Promise<{message: string}>
  setHistoryTags(id, tags)     // POST /api/history/{id}/tags → Promise<{message: string}>
  rollbackHistory(id)          // POST /api/history/{id}/rollback → Promise<{message: string, build_id: string}>
  getTokens()                  // GET /api/tokens            → Promise<{tokens: TokenRecord[]}>
  createToken(label, scopes = ['read'], expiresAt = null) // POST /api/tokens → Promise<TokenCreateResult>
  revokeToken(id)              // DELETE /api/tokens/{id}    → Promise<{message: string}>
  // 14A スナップショット
  getSnapshots()               // GET /api/snapshots         → Promise<{snapshots: SnapshotRecord[]}>
  downloadSnapshot(id)         // GET /api/snapshots/{id}/download → Promise<Blob>
  deleteSnapshot(id)           // DELETE /api/snapshots/{id} → Promise<{message: string}>
  // 14B メンテナンスモード
  getMaintenance()             // GET /api/maintenance       → Promise<MaintenanceObject>
  enableMaintenance(reason)    // POST /api/maintenance/enable → Promise<{message: string, since: string}>
  disableMaintenance()         // POST /api/maintenance/disable → Promise<{message: string}>
  // 14C IP アクセス制限
  getAccessControl()           // GET /api/access-control    → Promise<{allow: string[]}>
  setAccessControl(allowList)  // POST /api/access-control   → Promise<{message: string, allow: string[]}>
  // 14E フック
  getHooks()                   // GET /api/hooks             → Promise<{hooks: HookRecord[]}>
  addHook(phase, commandArgs, abortOnFailure = true, timeoutSeconds = 300) // POST /api/hooks → Promise<HookRecord>
  deleteHook(id)               // DELETE /api/hooks/{id}     → Promise<{message: string}>
  getHookLog(id)               // GET /api/hooks/{id}/log    → Promise<{id: string, runs: HookRunRecord[]}>
  // 15A アラートルール
  getAlertRules()              // GET /api/alert-rules       → Promise<{rules: AlertRule[]}>
  addAlertRule(metric, operator, threshold, level, message) // POST /api/alert-rules → Promise<AlertRule>
  deleteAlertRule(id)          // DELETE /api/alert-rules/{id} → Promise<{message: string}>
  // 15B 自動タグ付けルール
  getTagRules()                // GET /api/tag-rules         → Promise<{rules: TagRule[]}>
  addTagRule(condition, tags)  // POST /api/tag-rules        → Promise<TagRule>
  deleteTagRule(id)            // DELETE /api/tag-rules/{id} → Promise<{message: string}>
  // 15C チェックサム
  verifyOutput()               // POST /api/verify-output    → Promise<{match: boolean, expected: string, actual: string}>
  // 15D パイプライン設定
  getPipelineConfig()          // GET /api/pipeline-config   → Promise<PipelineConfig>
  setPipelineConfig(config)    // POST /api/pipeline-config  → Promise<{message: string}>
  // 15E 運用ノート
  getNotes()                   // GET /api/notes             → Promise<{content: string, updated_at: string|null}>
  setNotes(content)            // POST /api/notes            → Promise<{message: string, updated_at: string|null}>
  // 16B メール通知
  getSmtpConfig()              // GET /api/smtp-config       → Promise<SmtpConfig>
  setSmtpConfig(config)        // POST /api/smtp-config      → Promise<{message: string}>
  smtpTest()                   // POST /api/smtp-test        → Promise<{result: string, message: string}>
  // 16C ビルドキュー
  getQueue()                   // GET /api/queue             → Promise<{queued: QueueEntry[], max_size: number}>
  clearQueue()                 // DELETE /api/queue          → Promise<{message: string, cleared_count: number}>
  // 16D ダッシュボードレイアウト
  getDashboardLayout()         // GET /api/dashboard-layout  → Promise<{widgets: string[]}>
  setDashboardLayout(widgets)  // POST /api/dashboard-layout → Promise<{message: string}>
}

export { AdlaireCI, AdlaireCIError };
```

全メソッドは `Promise` を返す。`streamBuild` は SSE 接続確立後に `StreamHandle` で resolve し、接続前エラーは `AdlaireCIError` で reject する。HTTP エラー（4xx / 5xx）は `AdlaireCIError` としてスローする。`401` 受信時はセッション期限切れとして `this._token` をクリアする。constructor、private method、helper 関数を除く public method は §22.0e の SDK 列と完全一致させる。

**SDK 共通実装契約：**

| 項目 | 仕様 |
|------|------|
| `baseUrl` | 末尾 `/` を除去して保持する。空文字、`null`、`undefined` は `TypeError`。 |
| URL 組み立て | パスは `/api/...` をそのまま連結し、クエリ値は `encodeURIComponent` でエンコードする。 |
| 認証ヘッダー | `this._token` が存在する場合のみ `Authorization: Bearer ${token}` を付与する。 |
| JSON 送信 | `POST` / `DELETE` で body を送る場合は `Content-Type: application/json` を付与し、`JSON.stringify` した body を送信する。 |
| JSON 受信 | `Content-Type` が JSON の場合のみ `response.json()` を呼ぶ。§22.0e で JSON response を定義した endpoint の成功時に空 body を受信した場合は protocol error として `AdlaireCIError(status=0, message="Empty JSON response")` を投げる。 |
| `AdlaireCIError` | `name="AdlaireCIError"`、`status`、`message`、`details`、`responseBody` を持つ `Error` 派生クラスとする。constructor は `new AdlaireCIError({status, message, details = null, responseBody = null})` とし、network error、timeout、protocol error は `status=0` とする。`message` は API error response の `error`、network error は `"Network error"`、timeout は `"Request timeout"`、protocol error は固定文言を使用する。 |
| `logout()` | API 呼び出しが失敗しても `finally` で `this._token` をクリアする。 |
| request timeout | 通常 API は 30 秒で abort し、`AdlaireCIError(status=0, message="Request timeout")` を投げる。`streamBuild()` は接続確立まで 30 秒、接続確立後は timeout なしとし、利用者が `StreamHandle.close()` で停止する。 |
| `streamBuild()` | token がない場合は接続前に `AdlaireCIError(status=401, message="Unauthorized")` を投げる。native `EventSource` は Authorization header を付与できないため使用禁止とする。SDK は `fetch()`、`AbortController`、`ReadableStream` reader を使用し、`Accept: text/event-stream` と `Authorization` header を付与して SSE frame を解析する。`data:` 行の JSON を parse し、`type="log"` は `onLine(line)`、`type="end"` は `onEnd({status, duration_seconds})` を呼んで reader を close する。parse 不能 frame は `AdlaireCIError(status=0, message="Invalid SSE frame")` として stream error にする。 |
| `StreamHandle` | `streamBuild()` の戻り値は `{ close(): void, closed: boolean }` とする。`close()` は AbortController を abort し、複数回呼んでも例外を投げない。`closed` は `end` 受信、error、または `close()` 後に `true` になる。 |
| Blob レスポンス | `downloadSnapshot(id)` のみ `response.blob()` を使用する。その他は JSON とする。 |
| メソッド引数検証 | SDK 側でも必須引数の空値、配列型、数値範囲を検証し、HTTP 送信前に `TypeError` を投げる。 |
| endpoint 対応 | SDK method は §22.0e の SDK 列に存在する endpoint だけを呼び出す。§22.0e にない endpoint を SDK 独自判断で追加してはならない。 |
| body なし endpoint | §22.0e の `Request` が `none` の場合、SDK は `fetch` に `body` を設定しない。`{}` も送信しない。 |
| token 保存 | セッショントークンはメモリ上の `this._token` のみに保持する。`localStorage`、`sessionStorage`、Cookie へ保存しない。 |
| 秘密情報引数 | `updatePat(token)`、`setWebhookConfig(secret)`、SMTP password、`createToken()` の返却 token は console 出力しない。 |
| query 生成 | `undefined`、`null`、空文字の任意 query は送信しない。ただし仕様上 `""` が意味を持つ `q`、`from`、`to` は空文字を送ってよい。 |
| 戻り値補完禁止 | API response にない値を SDK が推測して追加しない。表示用加工は UI 側で行う。 |
| retry | SDK は自動 retry を行わない。ユーザー操作による再実行、または UI の明示的な再取得のみを許可する。 |

**SDK transport / error 固定契約：**

SDK の内部 request helper は、すべての public method で下表の処理順に固定する。public method ごとに個別 fetch 処理を複製してはならない。

| 順序 | 処理 | 固定仕様 |
|------|------|----------|
| 1 | 引数検証 | 必須引数、型、空配列、数値範囲を検証する。失敗時は `TypeError` を投げ、`fetch()` を呼ばない。 |
| 2 | URL 生成 | `baseUrl + path + query` を生成する。query key は method 契約表の順序で追加する。 |
| 3 | body 生成 | `Request=none` では `body` と `Content-Type` を設定しない。JSON body ありの場合だけ `JSON.stringify()` する。 |
| 4 | header 生成 | `Accept`、必要時 `Content-Type`、必要時 `Authorization` を付与する。token が空の場合は `Authorization` header を付けない。 |
| 5 | timeout 設定 | 通常 request は `AbortController` で 30 秒 timeout。`streamBuild()` は接続確立まで 30 秒。 |
| 6 | `fetch()` 実行 | network error、abort、CORS 等の失敗はすべて `AdlaireCIError(status=0,message="Network error")`、timeout だけ `"Request timeout"` とする。 |
| 7 | response parse | 成功 / 失敗に関わらず JSON error body がある場合は parse する。parse 不能 error body は `responseBody` に text を保持し、`message` は HTTP status 固定文言とする。 |
| 8 | token 変化 | `401` の場合だけ `this._token = null`。`403`、`429`、`500`、network error、timeout では token を破棄しない。 |
| 9 | return / throw | `2xx` は endpoint の型で返す。`4xx` / `5xx` は `AdlaireCIError` を投げる。 |

HTTP status と SDK error の対応は下表に固定する。

| 条件 | `AdlaireCIError.status` | `message` | `details` | token |
|------|-------------------------|-----------|-----------|-------|
| API JSON error | HTTP status | response の `error` | response の `details` が配列なら保持 | `401` のみ破棄 |
| API error body なし | HTTP status | `HTTP {status}` | `null` | `401` のみ破棄 |
| API error JSON parse 不能 | HTTP status | `HTTP {status}` | `null` | `401` のみ破棄 |
| network error | `0` | `Network error` | `null` | 維持 |
| timeout | `0` | `Request timeout` | `null` | 維持 |
| empty success JSON | `0` | `Empty JSON response` | `null` | 維持 |
| invalid success JSON | `0` | `Invalid JSON response` | `null` | 維持 |
| invalid SSE frame | `0` | `Invalid SSE frame` | `null` | 維持 |

`429` は SDK で自動待機、自動再送、自動 refresh を行わない。binary response は `downloadSnapshot(id)` の `2xx` のみ `Blob` とし、`4xx` / `5xx` では可能な限り JSON error として parse して `AdlaireCIError` を投げる。`streamBuild()` は接続後の `close()` をユーザー停止として扱い、`AdlaireCIError` を投げない。接続後に network error または invalid frame が発生した場合は `StreamHandle.closed=true` にし、呼び出し側へ stream error として通知できる状態にする。

**SDK メソッド実装固定契約：**

| 項目 | 仕様 |
|------|------|
| public method 定義順 | class 内の public method は §23 の一覧順に定義する。追加 public method を末尾に置くことは禁止し、先に §22.0e と本一覧を更新する。 |
| private helper | private helper は `_request`, `_json`, `_query`, `_requireToken`, `_validateId`, `_clearTokenOn401` の範囲で定義してよい。helper を export しない。 |
| TypeError 文言 | SDK 側引数検証の `TypeError.message` は `"Invalid argument: <name>"` に固定する。複数不正がある場合は最初に検出した引数だけを返す。 |
| path parameter | `id` を path に入れる method は、SDK 側で `encodeURIComponent(id)` を必ず行う。`/`、`.`、`..`、空文字は送信前に `TypeError`。 |
| query parameter | query key は §23 SDK 引数変換契約の表記順で生成する。任意 query が未指定の場合、`?` 自体を付けない。 |
| body parameter | body object の key 順は §23 SDK 引数変換契約の送信値順とする。未知 key を SDK が追加しない。 |
| token mutation | `login()` と `loginTotp()` は response に `token` が存在する場合だけ `this._token` を更新する。`totp_required:true` かつ token なしの場合は既存 token を保持せず `null` にする。 |
| logout failure | `logout()` は network error、`401`、`500` のいずれでも `finally` で `this._token=null` にする。 |
| response passthrough | 成功 response は clone、整形、既定値 merge を行わず、そのまま返す。Blob と StreamHandle は例外とする。 |
| error details | API error response の `details` が配列なら `AdlaireCIError.details` に同じ配列を保持する。配列でなければ `null`。 |
| responseBody | JSON parse できた error body は object のまま、parse 不能 error body は先頭 4000 文字の string として `responseBody` に保持する。 |

**SDK 検証 fixture：**

| fixture | 入力 | 合格条件 |
|---------|------|----------|
| auth token flow | `login()`、`loginTotp()`、`logout()`、`401` response | token set / clear が仕様どおり。localStorage、sessionStorage、Cookie を使わない。 |
| request shape | 全 public method を fake fetch で呼ぶ | method、path、query、body、headers が §22.0e と §23 引数変換契約に一致する。 |
| error shape | `400`、`401`、`403`、`422 details`、`500`、network error、timeout | `AdlaireCIError` の `status`、`message`、`details`、`responseBody` が固定値になる。 |
| stream | log frame、end frame、invalid frame、client close | callback、closed、error が仕様どおり。EventSource を使用しない。 |
| binary | `downloadSnapshot(id)` | `Blob` を返し、JSON parse を試みない。 |

**SDK P0 / P1 操作固定契約：**

P0 / P1 実装では、下表の SDK method を最小運用範囲として固定する。SDK は API response を成功時に補完せず、失敗時はすべて `AdlaireCIError` へ変換する。UI が必要とする表示用既定値、並べ替え、ラベル変換は SDK で行わない。

| SDK method | HTTP | 成功時 | 失敗時 | 追加禁止事項 |
|------------|------|--------|--------|--------------|
| `login(password)` | `POST /api/login` | `token` がある場合だけ `this._token` へ保存する。`totp_required:true` の場合は token を保存せず response を返す。 | `401`、`429`、`500` は `AdlaireCIError`。`401` で既存 token を破棄する。 | password を console、error、responseBody 加工結果へ出さない。 |
| `logout()` | `POST /api/logout` | response に関わらず `finally` で token を破棄する。 | network error、`401`、`500` でも token 破棄後に error を投げる。 | logout 失敗を理由に token を保持しない。 |
| `getStatus()` | `GET /api/status` | `StatusObject` をそのまま返す。 | `500 State file is corrupted` / `State file read failed` を message として保持する。 | `.build_status.json` 不在時の fallback 値を SDK が推測しない。 |
| `triggerBuild()` | `POST /api/build` | `{message, build_id?, queued?}` を返す。 | `409`、`429`、`503` を UI が分岐できる `status` 付き error にする。 | running / queue / circuit を SDK 側で事前判定しない。 |
| `buildForce()` | `POST /api/build/force` | `{message, build_id?, queued?}` を返す。 | `409`、`429`、`503` を `AdlaireCIError`。 | force 可否を SDK 側で状態推測しない。 |
| `cancelBuild()` | `POST /api/build/cancel` | `{message}` を返す。 | running なしの `409` を `AdlaireCIError(status=409)`。 | cancel 後に SDK が自動 `getStatus()` を呼ばない。 |
| `getLogs(n,q)` | `GET /api/logs` | `{lines}` を返す。`lines` は API 順序を保持する。 | `500` は固定 error message を保持する。 | line を結合、trim、level 分類しない。 |
| `getHistory(options)` | `GET /api/history` | `{total,page,per_page,pages,history}` を返す。 | query 不正 `422` は `details` を保持する。 | `pages`、`total` を SDK 側で再計算しない。 |
| `getHistoryLog(id)` | `GET /api/history/{id}/log` | log object を返す。 | `404`、`500` を `AdlaireCIError`。 | archive fallback を SDK 側で再試行しない。 |
| `getQueue()` | `GET /api/queue` | `{queued,max_size}` を返す。`running` が response に含まれる場合も削除しない。 | `.build_state` 破損の `500` を固定 message で保持する。 | queue 並び替え、重複排除、上限補正をしない。 |
| `resetCircuitBreaker()` | `POST /api/circuit-breaker/reset` | `{message,open,consecutive_failures}` を返す。 | 破損状態 `500` を `AdlaireCIError`。 | reset 成功後に SDK が自動 build を開始しない。 |
| `streamBuild(onLine,onEnd)` | `GET /api/build/stream` | `StreamHandle` を返し、`log` frame を `onLine`、`end` frame を `onEnd` へ渡す。 | 接続前 `401`、`404`、timeout、invalid frame を `AdlaireCIError`。 | `EventSource`、自動 reconnect、log 永続化を行わない。 |

**SDK P0 / P1 fixture 固定：**

| fixture | fake fetch 入力 | 合格条件 |
|---------|-----------------|----------|
| sdk p1 status corrupted | `GET /api/status` が `500 {"error":"State file is corrupted"}` | `AdlaireCIError.status=500`、`message="State file is corrupted"`、token 維持。 |
| sdk p1 build conflict | `POST /api/build` が `409 {"error":"Conflict"}` | `AdlaireCIError.status=409`、自動 retry なし、自動 `getStatus()` 呼び出しなし。 |
| sdk p1 queue full | `POST /api/build` が `429 {"error":"queue_full"}` | `AdlaireCIError.status=429`、`message="queue_full"`、body 再送なし。 |
| sdk p1 history paging | `getHistory({page:2,perPage:20,trigger:"manual"})` | query は `page=2&per_page=20&trigger=manual`、`total` と `pages` は API 値をそのまま返す。 |
| sdk p1 log not found | `GET /api/history/{id}/log` が `404 {"error":"Not found"}` | `AdlaireCIError.status=404`、`id` は `encodeURIComponent` 済み。 |
| sdk p1 stream end | `log` frame 2 件、`end` frame 1 件 | `onLine` 2 回、`onEnd` 1 回、`StreamHandle.closed=true`。 |
| sdk p1 stream invalid | `data:` 行が JSON parse 不能 | `AdlaireCIError(status=0,message="Invalid SSE frame")`、`closed=true`。 |
| sdk p1 unauthorized | 任意 P0/P1 endpoint が `401` | `this._token=null`、次 request に Authorization header を付けない。 |

**SDK P2〜P5 操作固定契約：**

P2〜P5 SDK は、§22.0e の endpoint 契約と §23 SDK 引数変換契約だけに従う。SDK は保存前検証の一部を `TypeError` で行ってよいが、API response の補完、no-op 判定、secret mask 変換、状態ファイル由来値の再計算を行ってはならない。

| 機能群 | SDK method | 成功時 | 失敗時 | 追加禁止事項 |
|--------|------------|--------|--------|--------------|
| config / repo / branch | `getConfig()`, `setConfig(config)`, `validateConfig(config)`, `setRepoConfig(config)`, `getBranchConfig()`, `setBranchConfig(branches)` | API response をそのまま返す。`set*` は success message と count / config を保持する。 | `422 details` は `AdlaireCIError.details` に保持する。`500` は固定 message。 | SDK 側で未知 key を削除しない。既定値 merge しない。 |
| schedule | `setScheduleInterval()`, `pauseSchedule()`, `resumeSchedule()`, `setAllowedHours()`, `clearAllowedHours()`, `setForceInterval()`, `setBuildCooldown()` | response の interval / hours / seconds / allowed_hours をそのまま返す。 | systemd 更新失敗の `500` を `AdlaireCIError` にする。 | timer 状態を SDK が推測しない。 |
| diagnostics / dashboard | `getDiagnostics()`, `getDashboard()`, `getRateLimit()`, `getDiskUsage()`, `getOutputMeta()` | read-only response をそのまま返す。 | item 単位 warn/error は成功 response として返し、HTTP error だけ例外にする。 | read-only response を SDK が保存・集計しない。 |
| notify / SMTP / webhook | `getNotifyConfig()`, `setNotifyConfig()`, `notifyTest()`, `notifyWeeklySummary()`, `getWebhookEvents()`, `getWebhookConfig()`, `setWebhookConfig()`, `getSmtpConfig()`, `setSmtpConfig()`, `smtpTest()` | secret は API が mask した値だけ返す。 | 未設定 `422` / `501`、送信失敗 `500` を保持する。 | secret 平文を console、throw message、responseBody 加工結果へ出さない。 |
| snapshots / rollback | `getSnapshots()`, `downloadSnapshot(id)`, `deleteSnapshot(id)`, `rollbackHistory(id)` | list は API 順序、download は Blob、rollback は `{message,build_id}`。 | `404`、running `409` を `AdlaireCIError`。 | SDK が snapshot 存在確認や rollback 可否を事前推測しない。 |
| maintenance / access / hooks | `getMaintenance()`, `enableMaintenance()`, `disableMaintenance()`, `getAccessControl()`, `setAccessControl()`, `getHooks()`, `addHook()`, `deleteHook()`, `getHookLog()` | API response をそのまま返す。 | access deny `403`、validation `422`、hook timeout `500` を保持する。 | command args を文字列結合しない。CIDR を SDK 独自正規化しない。 |
| alert / tag / pipeline / notes / dashboard layout | `getAlertRules()`, `addAlertRule()`, `deleteAlertRule()`, `getTagRules()`, `addTagRule()`, `deleteTagRule()`, `getPipelineConfig()`, `setPipelineConfig()`, `getNotes()`, `setNotes()`, `getDashboardLayout()`, `setDashboardLayout()` | API response の rules / config / notes / widgets をそのまま返す。 | duplicate `409`、validation `422`、read failure `500` を保持する。 | rule 重複排除、widget 補完、notes trim を行わない。 |
| tokens / sessions / audit | `getTokens()`, `createToken()`, `revokeToken()`, `getSessions()`, `revokeAllSessions()`, `getAuditLog()`, `getApiAccessLog()` | `createToken()` の token 本体は response として 1 回だけ返す。 | `403`、`404`、`422`、`429` を status 付きで保持する。 | token 本体を保存しない。token list に作成時 token を合成しない。 |

**SDK P2〜P5 fixture 固定：**

| fixture | fake fetch 入力 | 合格条件 |
|---------|-----------------|----------|
| sdk p2 config validation | `POST /api/config` が `422 details` | `AdlaireCIError.status=422`、`details` 配列保持、送信 body の未知 key は削除されていない。 |
| sdk p2 schedule failure | `POST /api/schedule/interval` が `500 {"error":"Internal server error"}` | error を投げ、SDK が timer 再試行や rollback request を行わない。 |
| sdk p3 secret mask | `GET /api/notify-config` と `GET /api/smtp-config` が mask 値を返す | mask 値をそのまま返し、secret 平文を生成しない。 |
| sdk p3 webhook events paging | `getWebhookEvents(20,40)` | query は `limit=20&offset=40`、`total` は API 値をそのまま返す。 |
| sdk p4 snapshot binary | `downloadSnapshot(id)` が binary response | `Blob` を返し、JSON parse を試みない。 |
| sdk p4 rollback conflict | `rollbackHistory(id)` が `409 {"error":"Build is running"}` | `AdlaireCIError.status=409`、自動 `getStatus()` 呼び出しなし。 |
| sdk p5 token issue | `createToken()` が `{token:"..."}` を返す | token を response として返すだけで、SDK 内部保存、console 出力、token list 合成をしない。 |
| sdk p5 duplicate rule | `addAlertRule()` または `addTagRule()` が `409 Conflict` | `AdlaireCIError.status=409`、自動 retry なし。 |

**SDK 引数変換契約：**

SDK method は、下表の通りに引数を path、query、body へ変換する。下表にない引数、既定値、body key を追加してはならない。

| SDK method | 引数 | 変換先 | 送信値 |
|------------|------|--------|--------|
| `login(password)` | `password` | body | `{password}` |
| `loginTotp(ticket,code)` | `ticket`, `code` | body | `{ticket,code}` |
| `changePassword(currentPassword,newPassword)` | `currentPassword`, `newPassword` | body | `{current_password: currentPassword, new_password: newPassword}` |
| `getAuditLog({limit,offset,actor,action,result})` | `limit=100`, `offset=0`, `actor=null`, `action=null`, `result=null` | query | `limit`、`offset` は常に送信する。任意値は `null` の場合送信しない。 |
| `getLogs(n,q)` | `n=100`, `q=""` | query | `n`、`q`。`q` は空文字でも送信する。 |
| `getHistory({page,perPage,trigger,tag,flagged})` | `page=1`, `perPage=20`, `trigger=null`, `tag=null`, `flagged=null` | query | `page`、`per_page: perPage` は常に送信する。`trigger`、`tag`、`flagged` は `null` の場合は送信しない。 |
| `getApiAccessLog({limit,offset,method,path,status})` | `limit=100`, `offset=0`, `method=null`, `path=null`, `status=null` | query | `limit`、`offset` は常に送信する。`method`、`path`、`status` は `null` の場合は送信しない。 |
| `validateConfig(config)` | `config` | body | `config` をそのまま送信する。保存は API が行わない。 |
| `setNotifyConfig(config)` | `config` | body | `config` をそのまま送信する。 |
| `setConfig(config)` | `config` | body | `config` をそのまま送信する。未知 key は送信前に削除せず、API の `422` に委ねる。 |
| `updatePat(token)` | `token` | body | `{token}` |
| `setScheduleInterval(seconds)` | `seconds` | body | `{interval_seconds: seconds}` |
| `setAllowedHours(from,to)` | `from`, `to` | body | `{from,to}` |
| `clearAllowedHours()` | なし | body | `{from:null,to:null}` |
| `setForceInterval(hours)` | `hours` | body | `{hours}` |
| `setBuildCooldown(seconds)` | `seconds` | body | `{seconds}` |
| `searchLogs(q,from,to,level)` | `q=""`, `from=""`, `to=""`, `level=undefined` | query | `q`、`from`、`to` は空文字でも送信する。`level` は指定時のみ送信する。 |
| `getWebhookEvents(limit,offset)` | `limit=50`, `offset=0` | query | `limit`、`offset` |
| `setWebhookConfig(secret)` | `secret` | body | `{secret}` |
| `confirmTotp(code)` | `code` | body | `{code}` |
| `disableTotp(code)` | `code` | body | `{code}` |
| `setApiRateLimit(policy)` | `policy` | body | `policy` をそのまま送信する。 |
| `setHistoryComment(id,comment)` | `id`, `comment` | path/body | path `{id}`、body `{comment}` |
| `setHistoryFlag(id,flagged)` | `id`, `flagged` | path/body | path `{id}`、body `{flagged}` |
| `setHistoryTags(id,tags)` | `id`, `tags` | path/body | path `{id}`、body `{tags}` |
| `createToken(label,scopes,expiresAt)` | `label`, `scopes=["read"]`, `expiresAt=null` | body | `{label,scopes,expires_at: expiresAt}` |
| `addHook(phase,commandArgs,abortOnFailure,timeoutSeconds)` | `phase`, `commandArgs`, `abortOnFailure=true`, `timeoutSeconds=300` | body | `{phase,command_args: commandArgs, abort_on_failure: abortOnFailure, timeout_seconds: timeoutSeconds}` |
| `addAlertRule(metric,operator,threshold,level,message)` | 各引数 | body | `{metric,operator,threshold,level,message}` |
| `addTagRule(condition,tags)` | `condition`, `tags` | body | `{condition,tags}` |
| `setPipelineConfig(config)` | `config` | body | `config` をそのまま送信する。 |
| `setNotes(content)` | `content` | body | `{content}` |
| `setSmtpConfig(config)` | `config` | body | `config` をそのまま送信する。`password` が未指定の場合は送信しない。 |
| `setDashboardLayout(widgets)` | `widgets` | body | `{widgets}` |

path に入る `id` は `encodeURIComponent` したうえで 1 segment として連結する。`/`、`.`、空文字を含む `id` は HTTP 送信前に `TypeError` とする。

**SDK method 完全性検証契約：**

SDK 実装完了時は、§22.0e の SDK 列に記載された method 名と `AdlaireCI.prototype` の public method 名が一致しなければならない。constructor、private method、helper 関数、`AdlaireCIError` は比較対象外とする。

| 検証項目 | 合格条件 |
|----------|----------|
| endpoint coverage | §22.0e の SDK 列で `none` 以外の method がすべて `AdlaireCI.prototype` に存在する。 |
| extra method | §22.0e の SDK 列に存在しない public method がない。 |
| request shape | 各 method が §22.0e の Request と §23 SDK 引数変換契約どおりの path / query / body を生成する。 |
| response handling | JSON endpoint は JSON object を返し、binary endpoint は `Blob`、stream endpoint は `StreamHandle` を返す。 |
| error handling | 4xx / 5xx、network error、timeout、empty JSON、invalid SSE frame が `AdlaireCIError` になる。 |
| token handling | `login()` 成功で token を保持し、`logout()` と `401` で token を破棄する。 |

**SDK 型定義表：**

本表は SDK が返す object 型の正本である。`nullable` は `null` を許可することを示す。配列は未取得時でも `[]` を返し、`undefined` を返してはならない。API response に存在しないキーを SDK が補完してはならない。ただし `GET /api/config`、`GET /api/notify-config`、`GET /api/dashboard-layout` の既定値 merge は API 側の責務とする。

| 型名 | 必須キー | nullable キー | 配列キー | 対応 API |
|------|----------|---------------|----------|----------|
| `StatusObject` | `last_sha`, `last_build_at`, `last_build_status`, `last_trigger`, `last_deploy_status`, `pending_transfers_count`, `notify_pending_count`, `circuit_open`, `output_url`, `running` | `last_sha`, `last_build_at`, `last_trigger`, `last_deploy_status`, `output_url` | なし | `GET /api/status` |
| `HistoryPageObject` | `total`, `page`, `per_page`, `pages`, `history` | なし | `history: HistoryRecord[]` | `GET /api/history` |
| `HistoryRecord` | `id`, `build_at`, `sha`, `status`, `trigger`, `output_size_bytes`, `duration_seconds`, `flagged`, `tags` | `sha`, `output_size_bytes`, `duration_seconds`, `comment`, `output_sha256`, `rollback_from` | `tags: string[]` | `.build_history` |
| `HistoryLogObject` | `.build_logs/{id}.json` の全必須キー、`lines` | `.build_logs/{id}.json` の nullable キー | `stdout`, `stderr`, `warnings`, `tags`, `lines` | `GET /api/history/{id}/log` |
| `CommentObject` | `id`, `comment`, `updated_at` | `comment`, `updated_at` | なし | `GET /api/history/{id}/comment` |
| `ConfigObject` | `.server_config` schema の全キー | `pat_expires_at`, `allowed_hours` | なし | `GET /api/config` |
| `ConfigValidationObject` | `valid`, `config`, `errors`, `warnings` | なし | `errors`, `warnings` | `POST /api/config/validate` |
| `AuditLogRecord` | `.audit_log` schema の全必須キー | `actor_id`, `target_id`, `remote_addr`, `message` | なし | `GET /api/audit-log` |
| `TotpStatus` | `enabled`, `confirmed_at` | `confirmed_at` | なし | `GET /api/auth/totp-status`, `POST /api/auth/totp-confirm`, `DELETE /api/auth/totp` |
| `ApiRateLimitPolicy` | `enabled`, `groups`, `state_summary` | なし | `groups`, `state_summary` | `GET/POST /api/api-rate-limit` |
| `NotifyConfig` | `webhooks`, `channels`, `on`, `summary`, `email` | `webhooks[].payload_template`, `webhooks[].secret` | `webhooks`, `channels`, `on`, `email.to`, `email.on` | `GET /api/notify-config` |
| `RepoInfoObject` | `owner`, `repo`, `branch`, `target_file` | なし | なし | `GET /api/repo-info` |
| `BranchTargetRecord` | `branch`, `target_file`, `sha_file`, `src`, `out`, `deploy_targets` | なし | `deploy_targets` | `GET /api/branch-config` |
| `SysinfoObject` | `output_size_bytes`, `output_mtime`, `uptime_seconds` | `output_mtime` | なし | `GET /api/sysinfo` |
| `HealthObject` | `status`, `last_build_at`, `last_build_status`, `last_deploy_at`, `last_deploy_status`, `pending_transfers`, `uptime_seconds` | `last_build_at`, `last_deploy_at` | なし | `GET /api/health` |
| `ScheduleObject` | `next_run_at`, `interval`, `paused`, `allowed_hours` | `next_run_at`, `allowed_hours` | なし | `GET /api/schedule` |
| `StatsObject` | `days`, `total_builds`, `success_count`, `failure_count`, `success_rate`, `avg_interval_minutes`, `avg_duration_seconds`, `max_duration_seconds` | `avg_interval_minutes`, `avg_duration_seconds`, `max_duration_seconds` | なし | `GET /api/stats` |
| `TimelineObject` | `days`, `timeline` | なし | `timeline` | `GET /api/stats/timeline` |
| `BuildDurationStats` | `n`, `count`, `avg_seconds`, `min_seconds`, `max_seconds`, `recent` | `avg_seconds`, `min_seconds`, `max_seconds` | `recent` | `GET /api/stats/build-duration` |
| `OutputMetaObject` | `size_bytes`, `mtime`, `heading_count`, `size_diff_bytes`, `tables_count`, `code_blocks_count`, `build_warnings`, `size_warn`, `sha256`, `build_id`, `commit_sha`, `build_at` | `mtime`, `size_diff_bytes`, `tables_count`, `code_blocks_count`, `sha256` | `build_warnings` | `GET /api/output-meta` |
| `DashboardObject` | `status`, `sysinfo`, `stats`, `schedule`, `alerts` | なし | `alerts` | `GET /api/dashboard` |
| `DiagnosticsObject` | `checked_at`, `items` | なし | `items` | `GET /api/diagnostics` |
| `RateLimitObject` | `limit`, `remaining`, `reset_at`, `used` | なし | なし | `GET /api/rate-limit` |
| `DiskUsageObject` | `build_logs_bytes`, `build_logs_count`, `build_logs_archive_bytes`, `build_logs_archive_count`, `output_file_bytes`, `total_bytes` | なし | なし | `GET /api/disk-usage` |
| `AccessRecord` | `.access_log` schema の全必須キー | `session_id`, `token_id`, `remote_addr`, `reason` | なし | `GET /api/access-log` |
| `ApiAccessRecord` | `.api_access_log` schema の全必須キー | `actor`, `remote_addr`, `user_agent`, `error` | なし | `GET /api/api-access-log` |
| `ConfigLogRecord` | `.config_log` schema の全必須キー | なし | なし | `GET /api/config-log` |
| `NotifyRecord` | `at`, `event`, `http_status`, `result`, `attempt`, `error` | `http_status`, `error` | なし | `GET /api/notify-log` |
| `SessionRecord` | `created_at`, `expires_at`, `last_used_at`, `current` | `last_used_at` | なし | `GET /api/sessions` |
| `WebhookEventRecord` | `timestamp`, `delivery_id`, `event`, `ref`, `sha`, `build_triggered` | `delivery_id`, `sha` | なし | `GET /api/webhook-events` |
| `SnapshotRecord` | `id`, `build_id`, `saved_at`, `size_bytes` | なし | なし | `GET /api/snapshots` |
| `MaintenanceObject` | `enabled`, `reason`, `since` | `reason`, `since` | なし | `GET /api/maintenance` |
| `QueueEntry` | `id`, `trigger`, `queued_at`, `requested_by`, `payload` | なし | なし | `GET /api/queue` |
| `TokenRecord` | `id`, `label`, `scopes`, `created_at`, `last_used_at`, `expires_at`, `revoked_at` | `last_used_at`, `expires_at`, `revoked_at` | `scopes` | `GET /api/tokens` |
| `TokenCreateResult` | `id`, `token`, `label`, `scopes`, `created_at`, `expires_at` | `expires_at` | `scopes` | `POST /api/tokens` |
| `HookRecord` | `id`, `phase`, `command_args`, `enabled`, `abort_on_failure`, `timeout_seconds` | なし | `command_args` | `GET/POST /api/hooks` |
| `HookRunRecord` | `build_id`, `ran_at`, `exit_code`, `output` | なし | なし | `GET /api/hooks/{id}/log` |
| `AlertRule` | `id`, `metric`, `operator`, `threshold`, `level`, `message` | なし | なし | `GET/POST /api/alert-rules` |
| `TagRule` | `id`, `condition`, `tags` | なし | `tags` | `GET/POST /api/tag-rules` |
| `PipelineConfig` | `extra_args`, `env` | なし | `extra_args` | `GET /api/pipeline-config` |
| `BuildChainConfig` | `chains` | なし | `chains` | `GET /api/build-chain-config` |
| `BuildTrendStats` | `count`, `avg_seconds`, `median_seconds`, `p95_seconds`, `anomaly_count`, `samples` | `avg_seconds`, `median_seconds`, `p95_seconds` | `samples` | `GET /api/stats/build-trends` |
| `SmtpConfig` | `host`, `port`, `user`, `tls`, `from`, `to`, `on`, `enabled`, `password_set` | `host`, `user`, `from` | `to`, `on` | `GET /api/smtp-config` |
| `ApprovalRecord` | `id`, `status`, `branch`, `sha`, `target`, `created_at`, `expires_at` | なし | なし | `GET /api/approvals` |
| `BackupObject` | `exported_at`, `server_config`, `notify_config` | なし | 設定内容に従う | `GET /api/backup` |
| `ExportObject` | `exported_at`, `history` | なし | `history` | `GET /api/history/export` |

---

## 24. 標準管理ツール 仕様

本節は、`admin/index.html` に関する仕様である。

**ファイル構成：**
```
/opt/adlaire-builder/admin/
├── index.html          # 管理画面（単一ファイル完結）
└── adlaire-ci-sdk.js   # SDK（標準管理ツールに同梱）
```

**DOM / section / form field 命名契約：**

標準管理ツールは、下表の DOM id、`data-panel`、form field name を使用する。表にない主要パネル id、主要 form name、主要 button id を追加してはならない。表示・非表示は `hidden` 属性で制御し、DOM 要素の生成順は本表の順序とする。

| パネル | section id | data-panel | 主フォーム id | 主要 field name | 主要 button id |
|--------|------------|------------|---------------|-----------------|----------------|
| ログイン | `panel-login` | `login` | `form-login` | `password`, `totp_code` | `btn-login`, `btn-login-totp` |
| パスワード変更 | `panel-password` | `password` | `form-password` | `current_password`, `new_password` | `btn-change-password` |
| ステータス | `panel-status` | `status` | なし | なし | `btn-refresh-status`, `btn-save-dashboard-layout` |
| 手動実行 | `panel-build` | `build` | なし | なし | `btn-build`, `btn-build-force`, `btn-build-cancel`, `btn-stream-close`, `btn-queue-clear` |
| 承認待ち | `panel-approvals` | `approvals` | なし | なし | `btn-refresh-approvals` |
| ログビューア | `panel-logs` | `logs` | `form-log-search` | `n`, `q`, `from`, `to`, `level` | `btn-load-logs`, `btn-search-logs`, `btn-export-logs`, `btn-cleanup-logs`, `btn-archive-logs` |
| ビルド履歴 | `panel-history` | `history` | `form-history-filter` | `page`, `per_page`, `trigger`, `tag`, `flagged` | `btn-export-history` |
| システム情報 | `panel-system` | `system` | `form-pat` | `token`, `pat_expires_at` | `btn-pat-verify`, `btn-pat-update` |
| 通知設定 | `panel-notify` | `notify` | `form-notify` | `webhooks`, `channels`, `on`, `summary`, `email`, `secret`, `smtp_password` | `btn-save-notify`, `btn-notify-test`, `btn-weekly-summary`, `btn-save-webhook-secret`, `btn-save-smtp`, `btn-smtp-test` |
| 設定 | `panel-config` | `config` | `form-config` | `log_max_lines`, `history_max_count`, `build_timeout_seconds`, `log_retention_days`, `log_archive_after_days`, `log_level`, `queue_max_size`, `snapshots_keep`, `build_retry_max`, `build_retry_base_seconds`, `commit_status_enabled`, `commit_status_context`, `commit_status_target_url`, `build_trend_keep_count`, `duration_anomaly_enabled`, `duration_anomaly_min_samples`, `duration_anomaly_avg_multiplier`, `duration_anomaly_p95_multiplier` | `btn-save-config`, `btn-validate-config`, `btn-set-log-level` |
| セキュリティ | `panel-security` | `security` | `form-security` | `totp_code`, `api_rate_enabled`, `api_rate_group`, `api_rate_window_seconds`, `api_rate_max_requests`, `session_timeout_seconds` | `btn-load-security`, `btn-totp-setup`, `btn-totp-confirm`, `btn-totp-disable`, `btn-save-api-rate-limit`, `btn-save-session-timeout` |
| アクセスログ | `panel-access-log` | `access-log` | `form-api-access-log-filter` | `limit`, `offset`, `method`, `path`, `status` | `btn-load-access-log`, `btn-load-api-access-log` |
| 監査ログ | `panel-audit-log` | `audit-log` | `form-audit-log-filter` | `limit`, `offset`, `actor`, `action`, `result` | `btn-load-audit-log` |
| 統計 | `panel-stats` | `stats` | なし | なし | `btn-load-stats` |
| リポジトリ情報 | `panel-repo` | `repo` | `form-repo` | `owner`, `repo`, `branch`, `target_file`, `interval_seconds`, `allowed_from`, `allowed_to`, `maintenance_reason` | `btn-save-repo`, `btn-save-branch-config`, `btn-pause-schedule`, `btn-resume-schedule`, `btn-enable-maintenance`, `btn-disable-maintenance` |
| セッション管理 | `panel-sessions` | `sessions` | なし | なし | `btn-load-sessions`, `btn-revoke-sessions` |
| システム診断 | `panel-diagnostics` | `diagnostics` | なし | なし | `btn-run-diagnostics`, `btn-verify-output` |
| ビルド比較 | `panel-compare` | `compare` | `form-compare` | `left_build_id`, `right_build_id` | `btn-compare-builds` |
| API トークン管理 | `panel-tokens` | `tokens` | `form-token` | `label`, `scopes`, `expires_at` | `btn-create-token` |
| 運用ノート | `panel-notes` | `notes` | `form-notes` | `content` | `btn-save-notes` |
| スナップショット | `panel-snapshots` | `snapshots` | なし | なし | `btn-load-snapshots` |
| メンテナンス | `panel-maintenance` | `maintenance` | `form-maintenance` | `reason` | `btn-maintenance-enable`, `btn-maintenance-disable` |
| アクセス制御 | `panel-access-control` | `access-control` | `form-access-control` | `allow` | `btn-save-access-control` |
| フック | `panel-hooks` | `hooks` | `form-hook` | `phase`, `command_args`, `abort_on_failure` | `btn-add-hook` |

共通領域の DOM id は、`app-root`、`nav-panels`、`global-banner`、`global-error`、`global-success`、`maintenance-banner`、`build-log-stream`、`build-queue-summary`、`issued-token-once` とする。エラー表示要素は各 panel 内に `id="{section-id}-error"`、成功表示要素は `id="{section-id}-success"` を置く。`label[for]` と input `id` は `field-{field_name}` 形式で一致させる。複数行・配列入力は `textarea` または table row で表現し、保存直前に SDK 引数の型へ変換する。

**画面構成：**

| パネル | 表示内容 | 表示条件 |
|-------|---------|---------|
| ログイン | パスワード入力フォーム。TOTP が必要な場合は同じ panel 内で TOTP code 入力へ切り替える。 | 未ログイン時のみ |
| パスワード変更 | 現在・新パスワード入力フォーム | `must_change: "prompt"` または `"forced"` 時（`"forced"` 時は他パネル非表示） |
| ステータス | 最終ビルド時刻・SHA・成否・出力サイトリンク・ダッシュボードウィジェット編集モード（表示するウィジェットをチェックボックスで選択・並び替え・保存） | ログイン済み |
| 手動実行 | ビルドトリガーボタン・SHA リセットを含む強制ビルドボタン・キャンセルボタン（実行中のみ有効）・実行結果表示・リアルタイムログ表示エリア（SSE ストリーミング）・キュー状態表示（待機中件数・クリアボタン） | ログイン済み |
| 承認待ち | 承認待ち build の一覧（id・branch・sha・target・status・created_at・expires_at）・承認ボタン・却下ボタン・再取得ボタン。pending 以外の entry は操作ボタンを disabled にする。 | ログイン済み |
| ログビューア | 最新ビルドログ（n 行・キーワードフィルター・ログレベルフィルターボタン（INFO / WARNING / ERROR / DEBUG）・JSON エクスポートボタン・横断検索フォーム（期間指定）・検索結果一覧） | ログイン済み |
| ビルド履歴 | 過去ビルド一覧（日時・SHA・成否・トリガー種別・所要時間・重要フラグ列・タグ列・ログ表示リンク・コメント入力欄）・フラグ付きのみ表示フィルター・タグフィルター・ページネーション UI・JSON エクスポートボタン | ログイン済み |
| システム情報 | 出力サイトサイズ・更新日時・稼働時間・ディスク使用量（ログ合計・出力サイト）・PAT 即時検証ボタン・PAT 更新フォーム・PAT 有効期限表示（設定フォーム・期限切れ間近で警告表示）・GitHub API レート制限表示 | ログイン済み |
| 通知設定     | Webhook 一覧（追加/削除/ラベル/有効無効切り替え/リトライ回数・間隔設定/シークレット入力欄）・通知条件設定（ビルド開始時・成功時・失敗時）・各 Webhook ペイロードテンプレート編集フォーム（変数一覧表示）・テスト送信ボタン・定期サマリー設定（間隔・時刻・曜日・即時送信ボタン）・送信履歴（試行回数・エラー内容列含む）・メール通知セクション（SMTP 設定フォーム・宛先リスト・通知条件・テスト送信ボタン） | ログイン済み |
| 設定         | ログ保持行数・履歴保持件数の設定変更・ビルドタイムアウト設定・ログレベル変更（INFO / DEBUG）・ログ保持期間（日数、0 = 無制限）・スナップショット保持世代数設定・ビルドキュー最大長設定・手動クリーンアップボタン・設定変更履歴（変更日時・項目・変更前後の値）・IP アクセス制限セクション（許可 IP / CIDR 一覧・追加フォーム・削除ボタン）・フック設定セクション（pre / post フック一覧・command_args 入力フォーム・実行ログリンク・有効無効切り替え）・アラートルール設定セクション（メトリクス・演算子・しきい値・レベル・メッセージの入力フォーム・ルール一覧・削除ボタン）・自動タグ付けルールセクション（条件式・タグ入力フォーム・ルール一覧・削除ボタン）・パイプライン設定セクション（追加引数入力欄・環境変数テーブル） | ログイン済み |
| セキュリティ | 単一 admin の TOTP 状態、TOTP setup secret / otpauth URI の一回表示、TOTP 確認/無効化、API rate limit policy、session timeout 設定。QR code 生成は初期実装対象外。 | ログイン済み |
| アクセスログ | ログイン履歴（日時・成否）・API アクセスログ（method、path、status、duration、actor、remote_addr）・API アクセスログフィルター | ログイン済み |
| 監査ログ | 設定変更、認証、token、build trigger、権限拒否の監査イベント一覧と actor / action / result フィルター。secret、password、token 本体は表示しない。 | ログイン済み |
| 統計         | ビルド回数・成功率・平均間隔・平均・最大ビルド時間・中央値・p95・異常件数・日別時系列データ（グラフ表示対応） | ログイン済み |
| リポジトリ情報 | 監視対象リポジトリ・ブランチ・ファイルの確認・設定変更フォーム（OWNER / REPO / BRANCH / TARGET_FILE）・ポーリング間隔変更フォーム・ポーリング一時停止／再開ボタン・許可時間帯設定（from〜to、解除ボタン）・メンテナンスモード有効化フォーム（理由入力）・解除ボタン・現在の状態表示 | ログイン済み |
| セッション管理 | 有効セッション一覧・全セッション強制無効化ボタン | ログイン済み |
| システム診断   | PAT・GitHub API・出力サイト・systemd・Webhook の診断項目一覧（ok / warn / error）・診断実行ボタン・アラートバッジ（`GET /api/dashboard` の `alerts` に基づき warn / error を表示）・出力整合性チェック項目（`POST /api/verify-output` 結果表示）・メンテナンスモード中はバナーを全パネル上部に表示 | ログイン済み |
| ビルド比較     | ビルド履歴から 2 件を選択・ログ並列表示・差分ハイライト（クライアントサイド処理） | ログイン済み |
| API トークン管理 | 発行済みトークン一覧（ラベル・スコープ・作成日時・最終使用日時・有効期限・失効日時）・新規発行フォーム（ラベル入力・scope 複数選択・有効期限）・発行時のみトークン文字列を表示・失効ボタン | ログイン済み |
| 運用ノート     | Markdown レンダリング表示・編集モード切替・保存ボタン・最終更新日時表示 | ログイン済み |
| スナップショット | ビルド成果物の世代一覧（最大`snapshots_keep`件）・個別ダウンロード・削除 | ログイン済み |
| メンテナンス   | メンテナンスモードの有効化（理由テキスト付き）・無効化・状態・開始時刻表示 | ログイン済み |
| アクセス制御   | 許可 IP / CIDR 一覧・CIDR 追加フォーム・削除ボタン（ブロック時は 403 を返す） | ログイン済み |
| フック         | Pre/Post ビルドフック一覧・command_args 追加・削除・実行ログ（直近 N 件）確認 | ログイン済み |

**UI パネル初期取得契約：**

各パネルを表示する時は、下表の SDK method を上から順に呼び出す。表示済み panel へ再遷移した場合も、ユーザー操作で表示した時点で同じ順序で再取得する。空配列はエラーではなく空状態として表示する。

| パネル | 初期取得 SDK method | 空状態表示 |
|--------|---------------------|------------|
| ステータス | `getDashboard()`, `getDashboardLayout()`, `getQueue()` | alerts なし、queue なしを 1 行で表示する。 |
| 手動実行 | `getStatus()`, `getQueue()`, `getMaintenance()` | queue なしを 1 行で表示する。 |
| ログビューア | `getLogs(100,"")` | ログなしを 1 行で表示する。 |
| ビルド履歴 | `getHistory({page:1,perPage:20})` | 履歴なしを 1 行で表示する。 |
| システム情報 | `getSysinfo()`, `getPatStatus()`, `getRateLimit()`, `getDiskUsage()` | 取得不能項目は該当ブロックにエラー表示し、他ブロックは表示する。 |
| 通知設定 | `getNotifyConfig()`, `getWebhookConfig()`, `getSmtpConfig()`, `getNotifyLog()` | webhook / email / log なしをそれぞれ 1 行で表示する。 |
| 設定 | `getConfig()`, `getConfigLog()`, `getAccessControl()`, `getHooks()`, `getAlertRules()`, `getTagRules()`, `getPipelineConfig()` | 各一覧なしを 1 行で表示する。 |
| セキュリティ | `getTotpStatus()`, `getApiRateLimit()`, `getConfig()` | TOTP disabled、rate limit disabled は通常状態として表示する。admin 権限なしの `403` は rate limit と session timeout 編集欄だけ非表示にする。 |
| アクセスログ | `getAccessLog()`, `getApiAccessLog({limit:100,offset:0})` | ログなしを 1 行で表示する。 |
| 監査ログ | `getAuditLog({limit:100,offset:0})` | log なしを 1 行で表示する。admin 権限なしの `403` は panel を非表示にする。 |
| 統計 | `getStats(7)`, `getStatsTimeline(30)`, `getStatsBuildDuration(20)` | 統計対象なしを 1 行で表示する。 |
| リポジトリ情報 | `getRepoInfo()`, `getBranchConfig()`, `getSchedule()`, `getMaintenance()` | branch target なしを 1 行で表示する。 |
| セッション管理 | `getSessions()` | 現 session だけの場合も通常一覧として表示する。 |
| システム診断 | `getDiagnostics()` | 診断 item なしは `No diagnostics` と表示する。 |
| ビルド比較 | `getHistory({page:1,perPage:100})` | 比較対象 2 件未満は compare button を disabled にする。 |
| API トークン管理 | `getTokens()` | token なしを 1 行で表示する。 |
| 運用ノート | `getNotes()` | content 空文字は空 editor として表示する。 |
| スナップショット | `getSnapshots()` | snapshot なしを 1 行で表示する。 |
| メンテナンス | `getMaintenance()` | disabled 状態を通常表示する。 |
| アクセス制御 | `getAccessControl()` | `allow:[]` は制限なしとして表示する。 |
| フック | `getHooks()` | hook なしを 1 行で表示する。 |

UI は、初期取得で一部 API が失敗した場合、ログイン状態を維持し、該当 panel の error 領域に失敗を表示する。ただし `401` は全 panel 表示を中止してログイン画面へ戻す。

**パスワード変更フロー：**
- `must_change: "prompt"`: パスワード変更パネルを表示。他パネルも操作可能
- `must_change: "forced"`: パスワード変更パネルのみ表示。変更完了後に通常画面へ遷移

**UI 操作契約表：**

標準管理ツールは、下表の SDK method 以外を直接呼び出してはならない。ファイル操作、`fetch()` の直接呼び出し、`systemctl` 実行、`components/runner.go` 直接起動は禁止する。成功時表示は対象パネル内に 1 行で表示し、失敗時表示は `AdlaireCIError.message` と `details` を同じパネル内に表示する。

| パネル | 操作 | SDK method | 成功時表示 | 成功後再取得 | disabled 条件 |
|--------|------|------------|------------|--------------|---------------|
| ログイン | ログイン | `login(password)` | 通常画面へ遷移、または TOTP 入力へ切替 | `getDashboard()`, `getStatus()`, `getQueue()` | password 空欄、送信中 |
| ログイン | TOTP 確認 | `loginTotp(ticket,code)` | 通常画面へ遷移 | `getDashboard()`, `getStatus()`, `getQueue()` | ticket なし、code 空欄、送信中 |
| パスワード変更 | 変更保存 | `changePassword(currentPassword,newPassword)` | `Password changed` | なし | 入力不足、送信中 |
| 全パネル共通 | ログアウト | `logout()` | ログイン画面へ戻る | なし | 送信中 |
| ステータス | 再読み込み | `getDashboard()` | 最終更新時刻を表示 | `getDashboard()` | 読み込み中 |
| ステータス | ウィジェット保存 | `setDashboardLayout(widgets)` | `Dashboard layout updated` | `getDashboardLayout()`, `getDashboard()` | widgets 空配列、送信中 |
| 手動実行 | ビルド開始 | `triggerBuild()` | `Build queued` または `Build started` | `getStatus()`, `getQueue()` | running true、maintenance enabled、送信中 |
| 手動実行 | 強制ビルド | `buildForce()` | `Force build queued` または `Force build started` | `getStatus()`, `getQueue()` | maintenance enabled、送信中 |
| 手動実行 | キャンセル | `cancelBuild()` | `Build cancelled` | `getStatus()`, `getQueue()` | running false、送信中 |
| 手動実行 | リアルタイムログ開始 | `streamBuild(onLine,onEnd)` | ログ行を追記 | `getStatus()`, `getQueue()` on end | token なし、接続中 |
| 手動実行 | キュー削除 | `clearQueue()` | `Queue cleared` | `getQueue()`, `getStatus()` | queue 空、送信中 |
| ログビューア | 最新ログ取得 | `getLogs(n,q)` | ログ行を表示 | なし | 読み込み中 |
| ログビューア | 横断検索 | `searchLogs(q,from,to)` | 検索結果件数を表示 | なし | 日付不正、検索中 |
| ログビューア | JSON export | `exportLogs()` | export 完了表示 | なし | 処理中 |
| ログビューア | cleanup | `cleanupLogs()` | 削除件数を表示 | `getLogs()`, `getDiskUsage()` | 処理中 |
| ログビューア | archive | `archiveLogs()` | 圧縮件数を表示 | `getLogs()`, `getDiskUsage()` | 処理中 |
| ビルド履歴 | 一覧取得 | `getHistory({page,perPage,trigger,tag,flagged})` | 件数を表示 | なし | 読み込み中 |
| ビルド履歴 | コメント保存 | `setHistoryComment(id,comment)` | `Comment saved` | `getHistory()`, `getHistoryComment(id)` | id 未選択、送信中 |
| ビルド履歴 | 重要フラグ切替 | `setHistoryFlag(id,flagged)` | `Flag updated` | `getHistory()` | id 未選択、送信中 |
| ビルド履歴 | タグ保存 | `setHistoryTags(id,tags)` | `Tags updated` | `getHistory()` | id 未選択、送信中 |
| ビルド履歴 | rollback | `rollbackHistory(id)` | `Rollback started` | `getHistory()`, `getStatus()` | running true、snapshot なし、送信中 |
| システム情報 | PAT 検証 | `patVerify()` | 検証結果を表示 | `getPatStatus()`, `getDiagnostics()` | 送信中 |
| システム情報 | PAT 更新 | `updatePat(token)` | `PAT updated` | `getPatStatus()`, `getDiagnostics()` | token 空欄、送信中 |
| 通知設定 | 通知設定保存 | `setNotifyConfig(config)` | `Notify config updated` | `getNotifyConfig()` | 入力不正、送信中 |
| 通知設定 | Webhook test | `notifyTest()` | 送信結果を表示 | `getNotifyLog()` | webhook 未設定、送信中 |
| 通知設定 | 週次 summary 送信 | `notifyWeeklySummary()` | success/failure 件数を表示 | `getNotifyLog()` | 宛先なし、送信中 |
| 通知設定 | Webhook secret 保存 | `setWebhookConfig(secret)` | `Webhook secret updated` | `getWebhookConfig()` | secret 空欄、送信中 |
| 通知設定 | SMTP 保存 | `setSmtpConfig(config)` | `SMTP config updated` | `getSmtpConfig()`, `getNotifyConfig()` | 入力不正、送信中 |
| 通知設定 | SMTP test | `smtpTest()` | 送信結果を表示 | `getNotifyLog()` | SMTP disabled、送信中 |
| 設定 | サーバー設定保存 | `setConfig(config)` | `Config updated` | `getConfig()`, `getConfigLog()` | 入力不正、送信中 |
| 設定 | サーバー設定検証 | `validateConfig(config)` | 検証結果を表示 | なし | 入力不正、送信中 |
| 設定 | log level 変更 | `setLogLevel(level)` | `Log level changed` | `getConfig()` | level 未選択、送信中 |
| 設定 | access control 保存 | `setAccessControl(allowList)` | `Access control updated` | `getAccessControl()`, `getConfigLog()` | CIDR 不正、送信中 |
| 設定 | hook 追加/削除 | `addHook()`, `deleteHook(id)` | hook 件数を表示 | `getHooks()`, `getConfigLog()` | 入力不正、送信中 |
| 設定 | alert rule 追加/削除 | `addAlertRule()`, `deleteAlertRule(id)` | rule 件数を表示 | `getAlertRules()`, `getDashboard()` | 入力不正、送信中 |
| 設定 | tag rule 追加/削除 | `addTagRule()`, `deleteTagRule(id)` | rule 件数を表示 | `getTagRules()` | 入力不正、送信中 |
| 設定 | pipeline config 保存 | `setPipelineConfig(config)` | `Pipeline config updated` | `getPipelineConfig()`, `getConfigLog()` | 入力不正、送信中 |
| セキュリティ | TOTP setup | `setupTotp()` | secret と otpauth URI を一回表示 | なし | TOTP enabled、送信中 |
| セキュリティ | TOTP confirm | `confirmTotp(code)` | `TOTP enabled` | `getTotpStatus()`, `getAuditLog()` | code 空欄、送信中 |
| セキュリティ | TOTP disable | `disableTotp(code)` | `TOTP disabled` | `getTotpStatus()`, `getAuditLog()` | code 空欄、送信中 |
| セキュリティ | API rate limit 保存 | `setApiRateLimit(policy)` | `API rate limit updated` | `getApiRateLimit()`, `getAuditLog()` | 入力不正、送信中 |
| セキュリティ | session timeout 保存 | `setConfig({session_timeout_seconds})` | `Session timeout updated` | `getConfig()`, `getConfigLog()` | 範囲外、送信中 |
| アクセスログ | 取得 | `getAccessLog()` | 件数を表示 | なし | 読み込み中 |
| アクセスログ | API アクセスログ取得 | `getApiAccessLog({limit,offset,method,path,status})` | 件数を表示 | なし | 読み込み中 |
| 監査ログ | 取得 | `getAuditLog({limit,offset,actor,action,result})` | 件数を表示 | なし | 読み込み中 |
| 統計 | 取得 | `getStats()`, `getStatsTimeline()`, `getStatsBuildDuration()` | グラフと数値を更新 | なし | 読み込み中 |
| リポジトリ情報 | repo 保存 | `setRepoConfig(config)` | `Repo config updated` | `getRepoInfo()`, `getConfigLog()` | 入力不正、送信中 |
| リポジトリ情報 | branch config 保存 | `setBranchConfig(branches)` | `Branch config updated` | `getBranchConfig()`, `getConfigLog()` | 入力不正、送信中 |
| リポジトリ情報 | schedule 変更 | `setScheduleInterval()`, `pauseSchedule()`, `resumeSchedule()`, `setAllowedHours()`, `clearAllowedHours()`, `setForceInterval()`, `setBuildCooldown()` | schedule 変更結果を表示 | `getSchedule()`, `getConfigLog()` | 入力不正、送信中 |
| セッション管理 | 全 revoke | `revokeAllSessions()` | revoke 件数を表示 | `getSessions()` | 送信中 |
| システム診断 | 診断実行 | `getDiagnostics()` | 診断結果を表示 | なし | 読み込み中 |
| システム診断 | 出力 checksum 検証 | `verifyOutput()` | match 結果を表示 | なし | 送信中 |
| API トークン管理 | token 発行 | `createToken(label,scopes,expiresAt)` | token 本体を 1 回だけ表示 | `getTokens()`, `getAuditLog()` | label 空欄、scope 未選択、送信中 |
| API トークン管理 | token 失効 | `revokeToken(id)` | `Token revoked` | `getTokens()` | id 未選択、送信中 |
| 運用ノート | 保存 | `setNotes(content)` | `Notes updated` | `getNotes()` | 送信中 |
| スナップショット | 一覧取得 | `getSnapshots()` | 件数を表示 | なし | 読み込み中 |
| スナップショット | download | `downloadSnapshot(id)` | browser download を開始 | なし | id 未選択、送信中 |
| スナップショット | 削除 | `deleteSnapshot(id)` | `Snapshot deleted` | `getSnapshots()` | id 未選択、送信中 |
| メンテナンス | 有効化 | `enableMaintenance(reason)` | `Maintenance mode enabled` | `getMaintenance()`, `getDashboard()` | reason 空欄、送信中 |
| メンテナンス | 無効化 | `disableMaintenance()` | `Maintenance mode disabled` | `getMaintenance()`, `getDashboard()` | 送信中 |

**UI 共通動作契約：**

| 項目 | 仕様 |
|------|------|
| 初期表示 | `localStorage` から token を復元しない。画面読み込み時は未ログイン状態から開始する。 |
| API 経路 | UI は必ず `AdlaireCI` SDK method を呼び出す。`fetch()`、`XMLHttpRequest`、`EventSource`、`ReadableStream` reader の直接生成は禁止する。ただし SDK の `streamBuild()` が返した `StreamHandle.close()` を呼ぶ操作は許可する。 |
| API 呼び出し中 | 対象ボタンを disabled にし、同一操作の二重送信を防ぐ。完了または失敗後に元へ戻す。 |
| 成功表示 | 変更系操作は成功時にパネル内へ 1 行の成功メッセージを表示し、関連 GET API を再取得する。 |
| 失敗表示 | SDK が投げた `AdlaireCIError.message` をパネル内エラー領域に表示する。`details` が配列の場合は各 `field` のフォーム項目に `message` を紐付け、該当項目が存在しない場合はパネル内エラー領域へ箇条書きで表示する。`responseBody`、token、secret、PAT は表示しない。 |
| `401` | token を破棄し、ログインパネルへ戻す。直前の入力値のうち秘密情報は消去する。 |
| `403` | 操作権限なしとしてエラー表示し、ログアウトはしない。 |
| `409` | 状態競合としてエラー表示し、ステータス・キュー・スケジュールを再取得する。 |
| `422` | 入力エラーとして該当フォーム項目へエラーを紐付ける。 |
| `429` | rate limit または login lock としてエラー表示し、同一操作ボタンを 10 秒間 disabled にする。ログインロックの場合は password field を空にする。 |
| `503` | メンテナンスバナーを表示し、ビルド操作ボタンを disabled にする。 |
| 秘密情報入力 | password、TOTP code、TOTP secret、PAT、Webhook Secret、SMTP password、発行直後 token は画面遷移、成功表示、再取得後にフォーム値から消去する。 |
| 自動更新 | ステータス、キュー、SSE 以外のパネルは自動ポーリングしない。ユーザー操作または画面表示時に取得する。 |
| SSE 切断 | `streamBuild()` が error になった場合はリアルタイム表示を停止し、`getStatus()` と `getQueue()` を再取得する。ユーザーが停止した場合はエラー表示しない。 |
| フォーム保存 | 保存 API が成功するまで UI 上の表示値を確定表示にしない。失敗時は入力値を保持する。 |
| 入力検証 | UI は送信前に必須入力、数値範囲、配列空、URL、CIDR、日付形式を検証する。UI 検証に通っても API 側検証は省略しない。 |
| 破壊的操作 | snapshot 削除、token 失効、queue clear、session revoke all、rollback はクリック後に確認ダイアログを 1 回表示する。確認文には対象 ID または件数を含める。 |
| 表示時刻 | API から受け取った UTC ISO 8601 をブラウザのローカル時刻で表示してよい。ただし data 属性または title 属性に元の ISO 8601 文字列を保持する。 |
| 一覧の空状態 | 配列が空の場合は、空表ではなくパネル内に 1 行の空状態メッセージを表示する。空状態はエラーとして扱わない。 |
| focus / aria | `422` は最初の invalid field へ focus する。`401` はログイン password field へ focus する。SSE ログ領域は `aria-live="polite"` とし、エラー領域は `role="alert"` とする。 |

**UI 初期ロード / イベント処理順序：**

標準管理ツールは、`DOMContentLoaded` 後に以下の順で初期化する。順序を入れ替えてはならない。

1. `app-root`、`nav-panels`、`global-error`、`global-success`、各 `panel-*` の存在を検査する。欠落時は `global-error` に `UI initialization failed` を表示し、以降の API 呼び出しを行わない。
2. `window.AdlaireCI` 等の global 参照を使わず、`./adlaire-ci-sdk.js` から `AdlaireCI` と `AdlaireCIError` を ES Module import する。
3. `AdlaireCI` を `new AdlaireCI({baseUrl})` で 1 回だけ生成する。`baseUrl` は同一 origin の `/api` を既定値とし、外部 origin は標準仕様では許可しない。
4. すべての panel を `hidden=true` にし、`panel-login` だけを表示する。
5. form submit と button click の event listener を登録する。登録対象は §24 の DOM / section / form field 命名契約表の id に限定する。
6. `localStorage`、`sessionStorage`、Cookie から token を読み込まない。
7. `global-error`、`global-success`、各 panel error/success を空にする。
8. login password field へ focus する。

ログイン成功後の初期取得順は、`getDashboard()` → `getStatus()` → `getQueue()` → `getMaintenance()` → `getDashboardLayout()` とする。途中で `401` を受信した場合は残りの取得を中止してログイン画面へ戻す。`getMaintenance()` が `enabled=true` を返した場合は `maintenance-banner` を表示し、ビルド開始、強制ビルド、rollback、hook 追加、設定変更系ボタンを disabled にする。

イベント処理は、各操作につき以下の順で行う。

1. 対象 panel の error/success を空にする。
2. UI 側入力検証を行う。失敗時は SDK method を呼ばない。
3. 対象 button と同一操作グループを disabled にする。
4. SDK method を呼ぶ。
5. 成功時は成功メッセージを表示し、§24 UI 操作契約表の成功後再取得を左から順に実行する。
6. 失敗時は `AdlaireCIError` として表示する。`TypeError` は UI 実装エラーとして `global-error` に `Client error` を表示する。
7. 秘密情報 field を消去する。
8. disabled を解除する。ただし `401`、`503`、SSE 接続中、メンテナンス中、または仕様上 disabled 条件が継続する場合は解除しない。

秘密情報 field は、`password`、`current_password`、`new_password`、`totp_code`、`token`、`secret`、`smtp_password`、`issued-token-once`、`totp-secret-once` とする。これらは成功、失敗、画面遷移、`401`、`logout()`、`revokeAllSessions()` のいずれの場合も DOM 値を空にする。発行直後 token は `issued-token-once`、TOTP setup secret は `totp-secret-once` に 1 回だけ表示し、次の任意の user action で消去する。

**UI 操作完全性検証契約：**

標準管理ツールの実装完了時は、§24 の DOM / section / form field 命名契約表と UI 操作契約表を照合し、下表を満たす。

| 検証項目 | 合格条件 |
|----------|----------|
| panel coverage | DOM / section / form field 命名契約表の `section id` がすべて `index.html` に存在する。 |
| button coverage | UI 操作契約表の各操作に対応する button または form submit が存在し、event listener が 1 つだけ登録される。 |
| SDK only | UI 操作契約表の SDK method 以外を UI から呼び出していない。直接 `fetch()`、`XMLHttpRequest`、`EventSource` を使用していない。 |
| success refresh | 成功後再取得列に複数 method がある場合、左から順に await し、途中失敗時は残りを中止して error 表示する。 |
| disabled restore | 操作失敗時も、継続条件がない限り disabled を解除する。`401`、`503`、SSE 接続中、メンテナンス中は解除しない。 |
| secret clearing | §24 の秘密情報 field が、成功、失敗、画面遷移、`401`、logout、revoke all の全経路で空になる。 |
| empty state | UI パネル初期取得契約の空状態表示が、各 panel 内に 1 行で表示される。 |
| global error | 初期化失敗、SDK constructor 失敗、想定外 `TypeError` は `global-error` に固定文言 `Client error` または `UI initialization failed` を表示する。 |

**UI 状態遷移固定契約：**

| 状態 | 表示 / 処理 |
|------|-------------|
| 未ログイン | `panel-login` だけを表示し、nav item は disabled。API token、session token、TOTP ticket を永続化しない。 |
| password 成功 / TOTP 必須 | login password field を消去し、同じ login panel 内で TOTP field を表示する。`ticket` はメモリだけに保持する。 |
| login 完了 | `panel-status` を表示し、ログイン成功後の初期取得順を実行する。 |
| must_change prompt | パスワード変更 panel を表示するが、他 panel 操作を許可する。 |
| must_change forced | パスワード変更 panel 以外を hidden または disabled にする。 |
| maintenance enabled | `maintenance-banner` を表示し、build、force build、rollback、hook 追加、設定変更系操作を disabled にする。 |
| session expired | token と ticket を破棄し、秘密情報 field を消去し、`panel-login` に戻す。 |
| fatal UI init error | API 呼び出しを行わず、`global-error` に `UI initialization failed` を表示する。 |

**UI 成功後再取得失敗契約：**

成功後再取得列に複数 SDK method がある操作では、成功 message を先に表示せず、再取得がすべて成功した後に成功 message を表示する。途中の再取得が失敗した場合は、対象操作自体は成功済みとして扱い、`global-success` に操作成功の固定文言、該当 panel error に再取得失敗を表示する。再取得失敗を理由に同じ変更 API を自動再実行してはならない。

**UI 操作状態固定契約：**

標準管理ツールは、同一操作の多重実行、API 成功前の確定表示、秘密情報の残存を防ぐため、各操作を下表の状態で管理する。

| 状態 | 開始条件 | UI 表示 | 許可される遷移 |
|------|----------|---------|----------------|
| `idle` | 初期状態、前操作完了後 | 操作可能。継続 disabled 条件がある場合は対象だけ disabled。 | `validating` |
| `validating` | ユーザー操作直後 | panel error / success を空にする。field error を初期化する。 | `blocked` / `sending` |
| `blocked` | UI 入力検証失敗、確認 dialog cancel | SDK method を呼ばない。該当 field error または表示変更なし。 | `idle` |
| `sending` | SDK method 呼び出し開始 | 同一操作 button / submit だけ disabled。spinner または loading 表示は対象 panel 内だけ。 | `refreshing` / `failed` |
| `refreshing` | 変更 API 成功後、成功後再取得が必要 | success を未表示のまま、再取得を左から順に実行する。 | `succeeded` / `refresh_failed` |
| `succeeded` | 変更 API と必要な再取得が成功 | panel success を 1 行表示し、secret field を消去する。 | `idle` |
| `refresh_failed` | 変更 API は成功したが再取得に失敗 | `global-success` に操作成功、panel error に再取得失敗を表示する。secret field を消去する。 | `idle` |
| `failed` | SDK method が `AdlaireCIError` を投げる | panel error と field error を表示する。secret field を消去する。 | `idle` / `unauthorized` |
| `unauthorized` | `401` | token / ticket / secret field を消去し、全 panel を hidden、`panel-login` だけ表示する。 | `idle` |

同一 panel に複数操作がある場合でも、`sending` による disabled は同一操作グループに限定する。ただし §24 の disabled 優先順位で maintenance、forced password change、SSE 接続中、`401` が上位条件として残る場合は、その上位条件に従う。UI は `sending` または `refreshing` の間、同じ SDK method を再実行してはならない。

秘密情報 field は、`blocked` のうち確認 dialog cancel を除き、`succeeded`、`refresh_failed`、`failed`、`unauthorized`、panel 遷移、logout、revoke all のいずれでも空にする。API 成功前に入力欄以外の確定表示、一覧更新、badge 更新、設定値反映を行ってはならない。入力中の form 値は、secret を除き、`failed` と `refresh_failed` では保持する。

**UI DOM 更新契約：**

| 項目 | 仕様 |
|------|------|
| 一覧描画 | API response 配列の順序を保持する。UI 独自 sort は行わない。 |
| 件数表示 | `total` がある API は `total` を表示する。配列長を total の代替にしない。 |
| 日時表示 | 表示テキストはローカル時刻でよいが、`datetime` または `title` に元の UTC ISO 8601 を保持する。 |
| disabled | 送信中 disabled は操作単位で行う。同一 panel の無関係 button は disabled にしない。ただし maintenance、forced password change、SSE 接続中は仕様上の対象をまとめて disabled にする。 |
| error 領域 | panel ごとに 1 つの error 領域を使う。field error は該当 field に紐付け、panel error にも summary を 1 行表示する。 |
| success 領域 | 変更系操作成功時だけ更新する。GET 再読み込みだけでは success を表示しない。 |
| secret one-time 表示 | issued token と TOTP secret は専用領域に 1 回だけ表示し、任意の次 user action、panel 遷移、logout、`401` で消去する。 |

**破壊的操作確認文言：**

| 操作 | 確認文 |
|------|--------|
| snapshot 削除 | `Delete snapshot {id}?` |
| token 失効 | `Revoke token {id}?` |
| queue clear | `Clear {count} queued builds?` |
| session revoke all | `Revoke all other sessions?` |
| rollback | `Rollback from build {id}?` |

確認 dialog で cancel した場合は SDK method を呼ばず、success / error 表示を変更しない。

**UI fixture 固定：**

| fixture | 入力 | 合格条件 |
|---------|------|----------|
| login totp | `login()` が `totp_required:true` を返す | password 消去、TOTP field 表示、ticket は DOM に表示しない。 |
| forced password | `must_change:"forced"` | password panel 以外が操作不可。変更成功後に通常初期取得を行う。 |
| refresh failure | 変更 API 成功後の再取得 2 件目が失敗 | 変更成功は維持し、再取得失敗だけ panel error に表示する。 |
| destructive cancel | 確認 dialog cancel | SDK method 呼び出し 0 回、表示差分なし。 |
| secret clearing | token 発行、TOTP setup、PAT 更新、Webhook secret 保存 | 次 user action または遷移で秘密情報 field と一回表示が消える。 |

**UI P0 / P1 操作固定契約：**

P0 / P1 UI は、ビルド状態確認、手動ビルド、強制ビルド、キャンセル、SSE ログ表示、履歴、ログ、キュー、circuit breaker reset だけを最小運用操作として固定する。UI は SDK response に存在しない状態を推測せず、API / SDK の error status と message に基づいて表示を分岐する。

| 操作 | 使用 SDK method | 成功時表示 | 成功後再取得 | 失敗時表示 / disabled |
|------|-----------------|------------|--------------|------------------------|
| 初期状態取得 | `getStatus()`, `getQueue()` | status badge、last build、queue 件数を表示する。 | なし | `500` は status panel error。`401` は login panel へ戻す。 |
| 手動 build | `triggerBuild()` | `Build queued` または `Build started` を表示する。 | `getStatus()`, `getQueue()` | `409` は競合表示後に status/queue 再取得。`429` は build button を 10 秒 disabled。`503` は maintenance/circuit 表示。 |
| 強制 build | `buildForce()` | `Force build queued` または `Force build started` を表示する。 | `getStatus()`, `getQueue()` | `409`、`429`、`503` は手動 build と同じ扱い。 |
| cancel | `cancelBuild()` | `Build cancel requested` を表示する。 | `getStatus()`, `getQueue()` | `409` は「実行中 build なし」として表示し、status/queue を再取得する。 |
| SSE 表示開始 | `streamBuild(onLine,onEnd)` | 接続中は stream indicator を表示し、受信行を append する。 | `onEnd` 後に `getStatus()`, `getQueue()`, `getLogs()` | invalid frame / network error は stream error 表示後に status/queue 再取得。ユーザー停止は error 表示なし。 |
| ログ取得 | `getLogs(n,q)` | 行順を保持して log panel に表示する。空配列は空状態表示。 | なし | `500` は log panel error。検索条件は保持する。 |
| 履歴取得 | `getHistory(options)` | `total`、`pages`、history 行を API 順序で表示する。 | なし | `422` は該当 filter field、`500` は history panel error。 |
| 履歴 log 表示 | `getHistoryLog(id)` | 選択 build の log detail を表示する。 | なし | `404` は選択解除して not found 表示。`500` は detail error。 |
| キュー取得 | `getQueue()` | queue 件数、max size、各 entry を API 順序で表示する。 | なし | `500` は queue panel error。 |
| circuit reset | `resetCircuitBreaker()` | `Circuit breaker reset` を表示する。 | `getStatus()`, `getQueue()` | `500` は circuit error。成功後に build を自動開始しない。 |

P0 / P1 UI の disabled 条件は以下に固定する。

| 条件 | disabled 対象 | 解除条件 |
|------|---------------|----------|
| SDK method 実行中 | 同一操作 button と同一 form submit | 成功または失敗後。ただし下記継続条件がある場合は解除しない。 |
| `streamBuild()` 接続中 | stream 開始 button、manual build button、force build button | `end` frame、stream error、またはユーザー停止。 |
| `getStatus().running=true` | manual build button。force build は仕様上許可される場合のみ有効。rollback は disabled。 | 次回 `getStatus().running=false`。 |
| `429` | 同一操作 button | 10 秒経過後に status/queue 再取得し、継続条件がなければ解除する。 |
| `503` maintenance / circuit | build、force build、cancel 以外の状態変更操作。circuit reset は有効。 | maintenance disabled または circuit reset 成功後の再取得。 |
| `401` | 全 authenticated 操作 | login 成功後。 |

**UI P0 / P1 fixture 固定：**

| fixture | fake SDK 入力 | 合格条件 |
|---------|---------------|----------|
| ui p1 initial status error | `getStatus()` が `AdlaireCIError(status=500,message="State file is corrupted")` | status panel error に固定 message を表示し、build button を成功扱いにしない。 |
| ui p1 manual build conflict | `triggerBuild()` が `409 Conflict` | error 表示、`getStatus()` と `getQueue()` をこの順で再取得、同じ build request を再送しない。 |
| ui p1 queue full | `triggerBuild()` が `429 queue_full` | build button を 10 秒 disabled、password や secret field は変更しない。 |
| ui p1 stream success | `streamBuild()` が log 2 件と end 1 件を返す | log 行 2 件を append、end 後に status、queue、logs を順に再取得、stream indicator を消す。 |
| ui p1 stream user close | ユーザーが `StreamHandle.close()` を押す | error 表示なし、closed 表示、status/queue 再取得あり。 |
| ui p1 history validation | `getHistory()` が `422 details` を返す | 該当 filter field に message を紐付け、history rows を前回表示のまま維持する。 |
| ui p1 log not found | `getHistoryLog(id)` が `404 Not found` | detail panel に not found を表示し、履歴一覧は再取得しない。 |
| ui p1 circuit reset | `resetCircuitBreaker()` 成功 | circuit 表示を閉じ、status/queue を再取得し、build を自動開始しない。 |
| ui p1 unauthorized | 任意操作が `401` | token/ticket/secret field を消去し、`panel-login` だけ表示する。 |

**UI P2〜P5 操作固定契約：**

P2〜P5 UI は、§24 UI 操作契約表の SDK method だけを呼び出す。UI は API / SDK response の補完、状態ファイル直接操作、未定義 endpoint 呼び出し、保存成功前の確定表示を行ってはならない。

| 機能群 | 主操作 | 成功時表示 | 成功後再取得 | 失敗時表示 / disabled |
|--------|--------|------------|--------------|------------------------|
| config / repo / branch | config 保存、repo 保存、branch config 保存、config validate | API message を表示する。validate は `valid` と errors / warnings を表示する。 | 保存系は `getConfig()` または対象 GET と `getConfigLog()`。validate は再取得なし。 | `422` は field error。`500` は panel error。入力値は保持する。 |
| schedule | interval、pause、resume、allowed hours、force interval、cooldown | 変更後値を表示する。 | `getSchedule()`, `getConfigLog()` | systemd 失敗 `500` は schedule panel error とし、再取得で保存済み値を表示する。 |
| diagnostics / dashboard | diagnostics、dashboard、rate limit、disk、output meta 取得 | item ごとの ok/warn/error を表示する。 | なし | HTTP error は panel error。item warn/error を HTTP error として扱わない。 |
| notify / SMTP / webhook | notify config 保存、webhook secret 保存、SMTP 保存、test、weekly summary、webhook events 表示 | API message、送信結果、件数を表示する。 | 保存系は対象 GET と `getConfigLog()`。test / summary は `getNotifyLog()`。 | secret 入力は成功・失敗の両方で消去する。未設定 `422` / `501` は panel error。 |
| snapshots / rollback | snapshot list、download、delete、rollback | list 件数、download 開始、delete 完了、rollback 開始を表示する。 | delete は `getSnapshots()`。rollback は `getHistory()`, `getStatus()`。 | delete / rollback は確認 dialog 必須。running `409` は status 再取得。 |
| maintenance / access / hooks | maintenance enable/disable、access 保存、hook 追加/削除 | 固定成功文言と件数または状態を表示する。 | `getMaintenance()` / `getAccessControl()` / `getHooks()` と `getConfigLog()`。 | maintenance enabled 中は build / rollback / 設定変更系を disabled。hook 追加失敗時は command 入力を保持する。 |
| alert / tag / pipeline / notes / layout | rule 追加/削除、pipeline 保存、notes 保存、dashboard layout 保存 | 固定成功文言を表示する。 | 対象 GET、必要時 `getDashboard()` または `getConfigLog()`。 | duplicate `409` は競合表示。validation `422` は field error。no-op は成功表示のみ。 |
| tokens / sessions / audit | token 発行/失効、session revoke、audit/API access log 表示 | token 発行時は token 本体を一回表示する。失効/revoke は固定成功文言。 | token 操作は `getTokens()`, `getAuditLog()`。session revoke は `getSessions()`。 | token 本体は次 user action、panel 遷移、logout、`401` で消去する。`403` は logout しない。 |

P2〜P5 UI の秘密情報消去条件は以下に固定する。

| 対象 field / 表示 | 消去タイミング |
|-------------------|----------------|
| PAT、Webhook Secret、SMTP password | 保存成功、保存失敗、panel 遷移、logout、`401`。 |
| 発行直後 API token | 次 user action、copy button 押下後、panel 遷移、logout、`401`。 |
| TOTP secret / ticket / code | confirm 成功、confirm 失敗、panel 遷移、logout、`401`。 |
| password / current_password / new_password | login / change 成功、login / change 失敗、logout、`401`。 |

**UI P2〜P5 fixture 固定：**

| fixture | fake SDK 入力 | 合格条件 |
|---------|---------------|----------|
| ui p2 config validation | `setConfig()` が `422 details` | 該当 field に error、panel error summary 1 行、入力値保持、`getConfig()` を呼ばない。 |
| ui p2 schedule save failure | `setScheduleInterval()` が `500` | panel error 表示後に `getSchedule()` を 1 回呼び、保存済み値を表示する。 |
| ui p3 secret save failure | `setWebhookConfig()` または `setSmtpConfig()` が `500` | secret field を消去し、secret 平文を error 表示しない。 |
| ui p3 notify test | `notifyTest()` 成功 | 結果表示後に `getNotifyLog()` を呼び、通知設定を自動保存しない。 |
| ui p4 snapshot delete cancel | delete 確認 dialog cancel | SDK method 呼び出し 0 回、success / error 表示差分なし。 |
| ui p4 rollback conflict | `rollbackHistory()` が `409 Build is running` | error 表示、`getStatus()` を呼ぶ、rollback request を再送しない。 |
| ui p4 maintenance enabled | `getMaintenance()` が enabled | maintenance banner 表示、build / rollback / 設定変更系 disabled、disable maintenance は enabled。 |
| ui p5 token issue once | `createToken()` 成功 | token 本体を一回表示し、`getTokens()` 後の一覧には token 本体を表示しない。 |
| ui p5 duplicate rule | `addAlertRule()` が `409 Conflict` | 競合表示、rule list は前回表示を保持し、自動 retry しない。 |
| ui p5 layout invalid | `setDashboardLayout()` が `422 details` | 該当 widget field error、dashboard 表示順を変更しない。 |

**UI 表示データ固定契約：**

| 表示対象 | 表示元 | 表示順 | 補完禁止 |
|----------|--------|--------|----------|
| dashboard widget | `getDashboardLayout().widgets`、`getDashboard()` | layout の `widgets` 順。未知 widget は表示せず、panel error に `Unknown dashboard widget` を 1 行表示する。 | UI が widget を自動追加、並び替え、既定復元してはならない。 |
| build history | `getHistory()` | API response の `history` 配列順。 | UI 独自 sort、欠落 duration の算出、status 名の言い換えは禁止。 |
| build compare | `getHistory({page:1,perPage:100})`、選択後 `getHistoryLog(left)`, `getHistoryLog(right)` | 左選択、右選択の順。ログ行は各 response の行順。 | API にない diff 結果を保存しない。比較結果は DOM 上の一時表示だけとする。 |
| approvals | `getApprovals()` | API response の `approvals` 配列順。`pending` 以外は操作 button disabled。 | 期限切れ判定を UI 時刻だけで確定しない。API status を正とする。 |
| tokens | `getTokens()`、`createToken()` | 一覧は API 配列順。発行直後 token は `issued-token-once` だけへ表示する。 | `GET /api/tokens` の record に token 本体を合成しない。 |
| notes | `getNotes()` | `content` を editor へそのまま入れる。表示 preview は HTML escape 後の簡易 Markdown 表示に限定する。 | UI が保存前に trim、整形、Markdown 拡張を行わない。 |
| hooks / rules | `getHooks()`、`getAlertRules()`、`getTagRules()` | API 配列順。 | UI 側で重複排除、無効化推測、command 文字列結合を行わない。 |
| pipeline config | `getPipelineConfig()` | `extra_args` と `env` を response 順で表示する。 | reserved arg の削除、env key の補完、inline YAML の再整形を行わない。 |

**UI 入力正規化固定契約：**

| 入力 | UI 正規化 | SDK 送信値 | 禁止事項 |
|------|-----------|------------|----------|
| 数値 field | ASCII 数字だけを整数化する。空文字は未指定として扱う。 | number または未指定。 | `Number("") == 0` 扱い、範囲外丸め、自動既定値保存は禁止。 |
| checkbox 群 | checked の DOM 順で配列化する。 | string array または boolean。 | 未選択時に UI が勝手に全選択へ戻さない。 |
| comma separated tag / scope | `,` で分割し、各要素の前後 ASCII whitespace だけ除去し、空要素を捨てる。 | string array。 | 重複削除、大小文字変換、未知値削除は API に委ねる。 |
| URL / webhook / SMTP host | 前後 whitespace を除去する。 | string。 | scheme 補完、host 置換、password 埋め込みは禁止。 |
| CIDR / IP | 前後 whitespace を除去し、空行を捨てる。 | string array。 | UI 側で CIDR 正規化や範囲展開を行わない。 |
| command_args | 1 行 1 引数として配列化する。空行は捨てる。 | string array。 | shell 文字列結合、quote 展開、環境変数展開は禁止。 |
| notes content | 入力値をそのまま送る。 | `{content}`。 | trim、改行正規化、Markdown 整形は禁止。 |
| date / expires_at | 空文字は `null`、入力ありは browser が返す ISO 互換文字列を送る。 | string/null。 | UI が現在時刻を補完しない。 |

**UI error / disabled 優先順位固定：**

| 優先 | 条件 | UI 処理 |
|------|------|---------|
| 1 | fatal UI init error | 全 API 呼び出しを停止し、`global-error` に `UI initialization failed`。 |
| 2 | `401` | token / ticket / secret field を消去し、全 panel を hidden、`panel-login` だけ表示。 |
| 3 | must_change forced | password panel 以外を hidden または disabled。 |
| 4 | maintenance enabled | build、force build、rollback、hook 追加、設定変更系を disabled。maintenance disable は有効。 |
| 5 | SSE 接続中 | stream 開始、manual build、force build を disabled。stream close は有効。 |
| 6 | 対象 SDK method 実行中 | 同一操作 button / submit だけ disabled。 |
| 7 | `429` | 同一操作 button を 10 秒 disabled。10 秒後に必要な再取得を行い、上位条件がなければ解除。 |
| 8 | `422 details` | 最初の invalid field へ focus し、field error と panel summary を表示。 |
| 9 | `403` / `409` / `500` | panel error を表示し、仕様上の再取得だけ実行。logout や自動 retry は行わない。 |

上位条件が残っている場合、下位条件の解除処理で button を有効化してはならない。複数 error が同時に発生した場合は、最上位の条件だけを global 表示し、下位の詳細は対象 panel error に残す。

**UI 詳細 fixture 固定：**

| fixture | fake SDK 入力 | 合格条件 |
|---------|---------------|----------|
| ui dashboard unknown widget | `getDashboardLayout()` が `["status","unknown","stats"]` を返す。 | `status`、`stats` だけ表示し、順序保持。`unknown` は panel error 1 行。layout 保存を自動実行しない。 |
| ui compare two builds | history 2 件選択後、左右の `getHistoryLog()` が異なる stdout を返す。 | 左右ログを別 column で API 行順表示し、差分 class は DOM 一時表示だけ。状態保存 API を呼ばない。 |
| ui compare missing build | 右側 `getHistoryLog()` が `404`。 | compare panel error、左側表示は維持、history 再取得なし、選択値は保持。 |
| ui approvals expired | `getApprovals()` が `status:"expired"` を含む。 | approve / reject button disabled、期限切れ表示、UI が pending へ戻さない。 |
| ui approval approve conflict | `approveBuild(id)` が `409`。 | error 表示後に `getApprovals()` を 1 回呼び、同じ approve を再送しない。 |
| ui notes preserve content | notes に前後空白と連続改行を含めて保存。 | `setNotes(content)` へ入力値そのまま送信し、trim しない。 |
| ui hook command args | 3 行の command args を入力し、中央行が空。 | 空行を除いた配列を `addHook()` に渡し、shell 文字列を作らない。 |
| ui pipeline reserved arg | `setPipelineConfig()` が `422 details`。 | field error を表示し、入力値を保持し、`getPipelineConfig()` を呼ばない。 |
| ui token issued clear | `createToken()` が token 本体を返す。 | `issued-token-once` に 1 回表示し、次 user action で消去。`getTokens()` の一覧に token 本体を表示しない。 |
| ui disabled priority | maintenance enabled 中に `429` が発生し 10 秒経過。 | maintenance が継続する限り build / rollback / 設定変更系は disabled のまま。 |

**UI 設定値契約：**

| 設定値 | 取得元 | 既定値 | 仕様 |
|--------|--------|--------|------|
| SDK `baseUrl` | `index.html` 内の `data-api-base-url` 属性 | `/api` | 空文字の場合は `/api` を使用する。外部 origin の URL は標準仕様では使用しない。 |
| 初期表示 panel | 固定値 | `panel-login` | token 永続化を行わないため、画面読み込み直後は常にログイン panel を表示する。 |
| theme token | `:root` CSS custom property | §6 の値 | JavaScript は theme token を変更しない。UI 操作で theme 切替を実装しない。 |
| panel 表示制御 | `hidden` 属性 | 全 panel hidden、`panel-login` のみ表示 | DOM 削除ではなく `hidden` で切り替える。 |
| API 呼び出し経路 | `AdlaireCI` instance | 1 instance | panel ごとに SDK instance を作らず、画面全体で 1 つの `AdlaireCI` instance を共有する。 |

---

## 25. 認証 実装仕様

**初期認証情報ファイル形式（JSON）：**
```json
{
  "password_hash": "<sha256_iter_v1_hex>",
  "salt": "<hex>",
  "algorithm": "sha256_iter_v1",
  "iterations": 260000,
  "must_change": true,
  "login_count": 0,
  "last_login_at": null,
  "updated_at": "2026-09-15T10:00:00Z"
}
```

**`.admin_credentials` schema 固定：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `password_hash` | string | 必須 | 64 文字 lowercase hex | password 本体は保存しない。 |
| `salt` | string | 必須 | 64 文字 lowercase hex | 32 bytes salt。 |
| `algorithm` | string | 必須 | `"sha256_iter_v1"` 固定 | 他 algorithm は初期実装で拒否する。 |
| `iterations` | integer | 必須 | `260000` 固定 | 値が異なる場合は認証を `500` で拒否する。 |
| `must_change` | boolean | 必須 | boolean | 初期生成時 `true`、パスワード変更後 `false`。 |
| `login_count` | integer | 必須 | 0 以上 | session token 発行成功時だけ +1。TOTP ticket 発行時は増やさない。 |
| `last_login_at` | string/null | 必須 | UTC ISO 8601 または `null` | session token 発行成功時だけ更新する。 |
| `updated_at` | string | 必須 | UTC ISO 8601 | password hash 更新時刻。 |

`.admin_credentials` に未知 key がある場合は credentials 破損として扱い、自動削除しない。必須 key 不足、型不一致、hex 不正、`algorithm` 不一致、`iterations` 不一致もすべて credentials 破損とする。API 起動時検証で credentials 破損を検出した場合は、§22.0a に従って ERROR ログを出し、HTTP サーバーを起動しない。HTTP サーバー稼働中の読込時検証で credentials 破損を検出した場合、`POST /api/login` と `POST /api/change-password` は `500 {"error":"Internal server error"}` を返す。API response、`.access_log`、`.audit_log`、journal に破損内容、hash、salt を出してはならない。

**ハッシュアルゴリズム：** Go 標準ライブラリのみで実装する `sha256_iter_v1`

| 項目 | 仕様 |
|------|------|
| salt 生成 | `crypto/rand` で 32 bytes を生成し、`encoding/hex` で 64 文字の hex 文字列として保存する。 |
| 初回 digest | `sha256(salt_bytes || password_utf8_bytes)` |
| 反復 | `iterations = 260000`。2 回目以降は `sha256(previous_digest || salt_bytes || password_utf8_bytes)` を繰り返す。 |
| 保存値 | 最終 digest を lowercase hex 文字列で `password_hash` に保存する。 |
| 比較 | 入力パスワードから同一手順で digest を生成し、`crypto/subtle.ConstantTimeCompare` で比較する。 |

`golang.org/x/crypto/pbkdf2` 等の外部パッケージは使用しない。PBKDF2、bcrypt、Argon2 等は本ファイルでは実装値を定義しない。password hash は `.admin_credentials.password_hash` に保存し、API response へは出さない。

**セッショントークン生成：**
`crypto/rand` で 32 bytes を生成し、`encoding/hex` で 64 文字の lowercase hex 文字列へ変換する。

**セッション管理：** `components/api.go` 内のインメモリ辞書で管理。有効期限は新規発行時点の `.server_config.session_timeout_seconds` とする。設定不在時は 8 時間。再起動で全セッション破棄。単一 admin の複数同時セッションを許容する。辞書 key は token 本体ではなく `sha256(token)` の lowercase hex とし、API response、`.access_log`、`.audit_log`、サーバーログへ token 本体を出力してはならない。

**認証入力・保存禁止契約：**

| 項目 | 仕様 |
|------|------|
| 認証 header | `Authorization: Bearer {token}` だけを受け付ける。 |
| Cookie | session cookie、remember-me cookie、CSRF cookie は発行しない。受信しても認証に使わない。 |
| query token | `?token=`、`access_token`、`session` query は認証に使わず、存在しても無視する。 |
| body token | login / totp 以外の body token は認証に使わない。 |
| 永続化禁止 | session token、login ticket、setup 仮 secret、連続失敗回数はファイル保存しない。 |
| response 禁止 | password hash、salt、session token hash、ticket hash、TOTP secret 保存値は response に含めない。 |
| log 禁止 | password、current_password、new_password、token、ticket、hash、salt、TOTP code は `.access_log`、`.audit_log`、`.api_access_log`、journal に含めない。 |

**セッション期限切れ時：** `401 Unauthorized` を返す。クライアント（SDK）は `this._token` をクリアし、再ログインを促す。

**認証処理の副作用境界：**

| ケース | HTTP status | `.admin_credentials` | メモリ session / ticket | `.access_log` | `.audit_log` | 備考 |
|--------|-------------|----------------------|-------------------------|---------------|--------------|------|
| password 不一致 | `401` | 変更なし | 変更なし | `login_failure` を追記 | `login_failure` を追記 | 連続失敗回数はメモリ上で +1。 |
| 連続失敗ロック | `429` | 変更なし | 変更なし | `login_locked` を追記 | `permission_denied` を追記 | password hash 検証は実行しない。 |
| password 成功 / TOTP 無効 | `200` | `login_count`、`last_login_at` 更新 | session 追加 | `login_success` を追記 | `login_success` を追記 | session token は全永続ログに保存しない。 |
| password 成功 / TOTP 有効 | `200` | 変更なし | ticket 追加 | `totp_required` を追記 | `totp_required` を追記 | `login_count` と `last_login_at` は更新しない。 |
| TOTP code 不一致 | `401` | 変更なし | ticket 削除 | `totp_failure` を追記 | `totp_failure` を追記 | ticket は再利用不可。 |
| TOTP 成功 | `200` | `login_count`、`last_login_at` 更新 | ticket 削除、session 追加 | `login_success` を追記 | `login_success` を追記 | session 期限は成功時点の設定で決める。 |
| session 期限切れ | `401` | 変更なし | 対象 session 削除 | 追記しない | 追記しない | `.api_access_log` は通常 API request として記録する。 |
| logout | `200` | 変更なし | 対象 session 削除 | `logout` を追記 | `logout` を追記 | token 本体と token hash は保存しない。 |
| password 変更成功 | `200` | 新 salt / hash、`must_change:false`、`updated_at` 更新 | 現 session 以外削除 | `password_change` を追記 | `password_change` を追記 | 新旧 password、hash、salt は保存しない。 |

上表で `.access_log` または `.audit_log` の追記が必要な処理は、response 返却前に追記を完了する。追記失敗時は、個別節で別指定がない限り `500 {"error":"Internal server error"}` を返す。session token、login ticket、TOTP setup secret は、必要なログ追記がすべて成功するまで response に含めてはならない。

**認証副作用順序固定契約：**

認証関連 endpoint は、以下の順序で副作用を確定する。response を返した後に session、ticket、credentials、ログを遅延更新してはならない。

| 処理 | 固定順序 | 失敗時 |
|------|----------|--------|
| `POST /api/login` password 不一致 | credentials 読込 → lock 判定 → hash 比較 → 失敗回数更新 → `.access_log` → `.audit_log` → `401` response | ログ追記失敗は `500`。失敗回数は戻さない。session/ticket は作成しない。 |
| `POST /api/login` password 成功 / TOTP 無効 | credentials 読込 → hash 比較 → 失敗回数 reset → `.admin_credentials` 更新 → session token 生成 → `.access_log` → `.audit_log` → token response | credentials またはログ追記失敗は `500`。token は response しない。 |
| `POST /api/login` password 成功 / TOTP 有効 | credentials 読込 → hash 比較 → 失敗回数 reset → ticket 生成 → `.access_log` → `.audit_log` → ticket response | ログ追記失敗は `500`。ticket は保存しない。 |
| `POST /api/login/totp` 成功 | ticket 検証 → TOTP secret 読込 → code 検証 → `.totp_secret.last_accepted_step` 更新 → `.admin_credentials` 更新 → session token 生成 → `.access_log` → `.audit_log` → token response | 途中失敗時は token を response しない。ticket は成功/失敗いずれも再利用不可にする。 |
| `POST /api/logout` | token 認証 → 対象 session 削除 → `.access_log` → `.audit_log` → response | ログ追記失敗は `500`。削除済み session は戻さない。 |
| `POST /api/change-password` | token 認証 → current password 検証 → new password 検証 → salt/hash 生成 → `.admin_credentials` 更新 → 現 session 以外削除 → `.access_log` → `.audit_log` → response | credentials 更新失敗は session を変更しない。ログ追記失敗時は `500` だが更新済み credentials は戻さない。 |
| `POST /api/sessions/revoke-all` | token 認証 → 現 session 以外を削除 → `.access_log` → `.audit_log` → response | ログ追記失敗時は `500`。削除済み session は戻さない。 |

session token と login ticket は `crypto/rand` 成功後にだけ生成し、生成した値はメモリ上で hash 化して保持する。response body に含める token / ticket は、その request の成功 response 1 回だけに含める。`403`、`429`、`500`、network 切断検出時に、未送信 token をログや状態ファイルへ退避してはならない。

**パスワード入力制約：**

| 項目 | 仕様 |
|------|------|
| 最小長 | 8 文字 |
| 最大長 | 128 文字 |
| 許可文字 | UTF-8 文字列。NUL 文字は禁止。前後空白はトリムせず、入力値そのものを検証・ハッシュ化する。 |
| 初期パスワード | `--init-credentials` で `.admin_credentials` に生成する。初回ログイン時は `must_change: "prompt"` を返す。 |
| 変更時検証 | `new_password` が現在パスワードと同一の場合は `422` を返す。 |
| 失敗時応答 | パスワード不一致は `401` と `{"error":"Unauthorized"}` を返し、どの条件に失敗したかは返さない。 |
| 成功時保存 | `.admin_credentials` に新 salt、新 hash、`must_change:false`、`updated_at` を原子的に保存する。 |

**セッションレコード形式（メモリ上）：**

```json
{
  "token_hash": "<sha256_hex>",
  "session_id": "<16_byte_hex>",
  "created_at": "2026-09-15T10:00:00Z",
  "expires_at": "2026-09-15T18:00:00Z",
  "last_used_at": "2026-09-15T10:05:00Z"
}
```

`session_id` は `crypto/rand` で 16 bytes を生成し、lowercase hex とする。`GET /api/sessions` は `session_id` ではなく `current`、`created_at`、`expires_at`、`last_used_at` のみ返す。認証必須 API で有効 token を受信した場合、`last_used_at` を現在時刻へ更新する。期限切れ token は検出時にメモリから削除する。`POST /api/logout` は対象 token のみ削除する。`POST /api/sessions/revoke-all` は現在 token 以外を削除する。

**ログイン失敗制御：**

| 項目 | 仕様 |
|------|------|
| 失敗記録 | `components/api.go` はメモリ上で直近の連続ログイン失敗回数と最終失敗時刻を保持する。再起動で失敗回数はリセットされる。 |
| ロック条件 | 連続 10 回失敗した場合、最終失敗から 10 分間 `POST /api/login` を `429 Too Many Requests` と `{"error":"Too many attempts"}` で拒否する。 |
| 成功時 | ログイン成功時は連続失敗回数を 0 に戻す。 |
| 応答時間 | パスワード不一致、存在しない credentials、ロック中を除く検証失敗では、条件の詳細をレスポンスへ出さない。 |
| ログ | 成功、失敗、ロック拒否はいずれも `.access_log` へ追記する。password、token、hash、salt は記録しない。 |

`.access_log` または `.audit_log` 追記失敗時は、ログイン失敗では `500` を返し、失敗回数は増加済みのままとする。ログイン成功時は session token 発行前に `.admin_credentials`、`.access_log`、`.audit_log` を更新し、いずれかに失敗した場合は session token を発行しない。TOTP 有効時は ticket 発行前に `.access_log` と `.audit_log` を追記し、追記失敗時は ticket を発行しない。

**ログインフロー：**
```
POST /api/login
  ├─ 連続失敗ロック中 → 429
  └─ パスワードハッシュ検証
       ├─ 失敗 → 連続失敗回数 + 1 → .access_log 追記 → .audit_log 追記 → 401
       └─ 成功 → 連続失敗回数を 0 へリセット
                  ├─ must_change == true → must_change: "prompt"
                  ├─ TOTP 有効 → ticket 生成・返却
                  └─ TOTP 無効 → login_count / last_login_at 更新 → セッショントークン生成・返却
```

**パスワード変更時：** 新しい salt を生成しハッシュを更新し、`must_change:false` を保存する。変更完了後に現セッション以外のセッションを破棄。

**`--init-credentials` オプション：** `components/api.go` を `--init-credentials` 引数で起動した場合、初期パスワード `admin` で `.admin_credentials` を生成して終了する（HTTP サーバーは起動しない）。

`.admin_credentials` が既に存在する場合、`--init-credentials` は上書きせず `409` 相当の終了コード `2` で終了し、標準エラーへ `credentials already exist` を出力する。初期化成功時の終了コードは `0` とする。

**`--init-credentials` CLI 固定契約：**

| ケース | stdout | stderr | 終了コード | 副作用 |
|--------|--------|--------|------------|--------|
| 新規生成成功 | `credentials initialized` + LF | 空 | `0` | `.admin_credentials` を mode `0600` で作成する。 |
| 既存あり | 空 | `credentials already exist` + LF | `2` | 既存ファイルを変更しない。 |
| `--state-dir` 相対 path | 空 | `state directory must be absolute: {path}` + LF | `2` | ファイル作成なし。 |
| 書込失敗 | 空 | `credentials write failed` + LF | `1` | tmp を削除し、部分ファイルを残さない。 |
| rand 失敗 | 空 | `random source failed` + LF | `1` | ファイル作成なし。 |

生成手順は、state dir 検証 → 既存確認 → salt 生成 → hash 生成 → `{path}.tmp.{pid}` へ JSON + LF 書込 → mode `0600` → file sync → rename → parent directory sync の順に固定する。rename 後の sync に失敗した場合は `1` を返し、作成済みファイルは残る。実装者判断で初期パスワードを環境変数、対話入力、ランダム生成へ変更してはならない。

**認証 fixture 固定：**

| fixture | 入力 | 期待結果 |
|---------|------|----------|
| auth-password-failure | 誤 password で `POST /api/login` | `401`、session/ticket なし、失敗回数 +1、`.access_log` と `.audit_log` に secret なし。 |
| auth-login-lock | 連続 10 回失敗後の `POST /api/login` | `429`、password hash 検証なし、`.access_log` に `login_locked`、`.audit_log` に `permission_denied`。 |
| auth-session-issued | TOTP 無効で password 成功 | token は response のみ、`.admin_credentials.login_count` +1、ログに token/hash/salt なし。 |
| auth-session-expired | 期限切れ session で保護 API | `401`、対象 session 削除、`.access_log` と `.audit_log` は追記しない。 |
| auth-password-change | password 変更成功 | 新 salt/hash、現 session 以外削除、`password_change` ログ、password/hash/salt 平文なし。 |
| auth-log-write-failure | login 成功時に `.audit_log` 追記失敗 | `500`、session token を response しない。 |

---

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

`POST /api/config/validate` は、`POST /api/config` と同じ入力を受け取り、保存せずに検証結果を返す。

Request body は partial `ConfigObject` とする。未知 key を含む場合は `422 {"error":"Unknown config key","key":"..."}` を返す。型不一致、範囲外、URL 不正、commit status context 空文字、retry 設定不正は `200` で `valid=false` として返す。JSON 不正は `422` とする。

成功時 response:

```json
{
  "valid": false,
  "config": {"log_max_lines": 500},
  "errors": [
    {"field":"build_retry_max","code":"out_of_range","message":"build_retry_max must be 0..10"}
  ],
  "warnings": [
    {"field":"commit_status_target_url","code":"http_url","message":"https is recommended"}
  ]
}
```

`config` は既存 `.server_config` に request body を merge した正規化後の値を返す。ただし保存してはならない。`.config_log`、`.api_access_log` 以外の状態ファイルを更新してはならない。`.api_access_log` は通常 API 呼び出しとして記録する。

**validate 副作用固定契約：**

| 項目 | 仕様 |
|------|------|
| merge | request body は既存 config に shallow merge し、object 値は対象 object 単位で置換する。 |
| secret key | secret 風 key は未知 key と同じく `422`。値は response、access log、server log に出さない。 |
| `valid=false` | HTTP status は `200`。状態ファイルは更新しない。 |
| error 順 | request body の key 出現順で `errors[]` を返す。 |
| warning 順 | field 名昇順で返す。 |

検証条件:

| ケース | 期待結果 |
|--------|----------|
| 正常値 | `valid=true`、errors 空、状態差分なし。 |
| 範囲外 | `valid=false`、該当 field error。 |
| 未知 key | HTTP 422、状態差分なし。 |
| secret 風 key | HTTP 422、response と log に値を含めない。 |
| 複数 error | request key 出現順で返る。 |
| valid false | HTTP 200、保存差分なし。 |

### 27.6 API アクセスログ

`components/api.go` は全 `/api/` request について `.api_access_log` へ JSON Lines を追記する。`GET /api/health` も対象とする。静的 file 配信、admin HTML、SDK JS は対象外とする。

追記タイミングは response status 確定後とする。追記失敗時は、対象 API の本来の response を優先し、サーバーログに `API_ACCESS_LOG_WRITE_FAILED` を出す。access log 書き込み失敗を理由に API response を `500` へ変更してはならない。

`GET /api/api-access-log` は `limit`、`offset`、`method`、`path`、`status` query を受け付ける。`limit` は 1〜1000、既定値 100。`offset` は 0 以上、既定値 0。`method` は大文字 HTTP method 完全一致。`path` は prefix match。`status` は HTTP status 完全一致。壊れた行は無視し、新しい順で返す。

**access log record 固定契約：**

| key | 型 | 仕様 |
|-----|----|------|
| `timestamp` | string | response status 確定時刻。 |
| `request_id` | string | 16 byte hex。response header `X-Request-Id` と一致させる。 |
| `method` | string | 大文字 HTTP method。 |
| `path` | string | query を除いた path。 |
| `status` | integer | 実際に返した HTTP status。 |
| `actor` | string/null | admin、token id、webhook、または `null`。 |
| `auth_type` | string | `"none"`、`"session"`、`"api_token"`、`"webhook"`。 |

request body、query 全体、cookie、Authorization header、token、password、secret は保存しない。access log 追記失敗時も、既に確定した response status を変更しない。

検証条件:

| ケース | 期待結果 |
|--------|----------|
| 認証成功 API | actor `admin` または token id で記録される。 |
| 認証失敗 API | status 401、actor null、secret なし。 |
| Webhook | auth_type `webhook`、actor `webhook`。 |
| filter | method/path/status が完全に効く。 |
| 書込失敗 | 本来 response を維持し、server log に WARN。 |
| request id | response header と log の `request_id` が一致する。 |
| secret query | query 全体を保存しない。 |

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

本機能の目的は、管理 API から systemd timer の実行間隔を変更し、次回以降の runner 起動間隔を固定仕様どおり反映することである。

owner component は `api` とする。collaborator component は `runner`、`statefile` とする。`api` は schedule 設定と systemd timer 更新の責務を持つ。`runner` は本機能で systemd timer を変更しない。

**入力 / 出力：**

| 項目 | 仕様 |
|------|------|
| API | `POST /api/schedule/interval` |
| Request | `{ "interval_seconds": integer }` |
| 許容値 | 30〜86400 秒 |
| Response | `{ "message": "Schedule interval updated", "interval_seconds": N }` |
| 状態 | `.server_config.schedule_interval_seconds` を更新し、`.config_log` へ記録する。 |
| systemd | `/etc/systemd/system/adlaire-ci.timer` の `OnUnitActiveSec` を `N seconds` 相当へ更新する。 |

**処理順序：**

1. 認証、maintenance、入力型、範囲を検証する。
2. `.server_config` を atomic write で更新する。
3. `.config_log` に `schedule_interval_seconds` の diff を追記する。
4. systemd timer ファイルを書き換える。
5. `systemctl daemon-reload` を実行する。
6. `systemctl restart adlaire-ci.timer` を実行する。
7. `systemctl show adlaire-ci.timer -p OnUnitActiveSec` 相当で反映を確認する。
8. 成功 response を返す。

**systemd 更新固定契約：**

| 項目 | 仕様 |
|------|------|
| timer 書換 | 既存 unit を直接編集せず、管理対象 drop-in `/etc/systemd/system/adlaire-ci.timer.d/override.conf` を atomic write する。 |
| drop-in 内容 | `[Timer]`、`OnUnitActiveSec={N}s`、`Persistent=true` のみを書き込む。 |
| 反映確認 | `systemctl show` の値を秒へ正規化して request 値と一致確認する。 |
| rollback | systemd 失敗時に `.server_config` は巻き戻さない。失敗を `.config_log` に追加する。 |

**異常系：**

| 条件 | 応答 / 処理 |
|------|-------------|
| 入力が integer でない、または範囲外 | `422`。状態ファイル、systemd は変更しない。 |
| `.server_config` 書き込み失敗 | `500`。systemd は変更しない。 |
| systemd timer 書き換え失敗 | `500`。`.server_config` は更新済みのまま残し、`.config_log` に `systemd_update_failed` を記録する。 |
| daemon-reload / restart / show 失敗 | `500`。server log に `SCHEDULE_INTERVAL_APPLY_FAILED` を出す。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 正常値 300 | `.server_config`、timer、API response が 300 で一致する。 |
| 最小値 30 | 成功する。 |
| 範囲外 29 | `422`、差分なし。 |
| systemd 失敗 | `500`、`.config_log` に失敗記録。 |
| show 不一致 | `500`、server config は更新済み。 |

### 27.12 GitHub Webhook 受信

本機能の目的は、GitHub push event を HMAC-SHA256 署名検証したうえで受信し、定期 polling を待たずに build queue へ投入することである。

owner component は `api` とする。collaborator component は `runner`、`statefile` とする。`api` は署名検証、イベント記録、queue 投入を担当し、`runner` は queue entry を処理する。

**入力 / 出力：**

| 項目 | 仕様 |
|------|------|
| API | `POST /api/webhook` |
| 必須 header | `X-GitHub-Event`, `X-GitHub-Delivery`, `X-Hub-Signature-256` |
| 対象 event | `push` のみ |
| Secret | `.webhook_secret` |
| 成功 response | `{ "message": "Webhook accepted", "queued": true, "event_id": "..." }` |
| 状態 | `.webhook_events.json` へ追記し、必要に応じて `.build_state.queued` へ `trigger="webhook"` entry を追加する。 |

**署名検証：**

署名は `sha256=` prefix を含む lowercase hex とする。検証は raw request body に対して `HMAC-SHA256(secret, body)` を計算し、定数時間比較で行う。secret 不在、header 不在、prefix 不正、hex 不正、署名不一致はすべて `401` とし、queue へ投入しない。

**処理順序：**

1. method と body size を検証する。
2. `.webhook_secret` を読み込む。
3. HMAC 署名を検証する。
4. JSON body を parse する。
5. event が `push` であることを確認する。
6. `ref`、`after`、`repository.owner.login`、`repository.name` を抽出する。
7. `.branch_config` と照合し、対象 branch がある場合だけ queue へ追加する。
8. `.webhook_events.json` にイベント結果を JSON Lines で追記する。
9. response を返す。

**Webhook 副作用固定契約：**

| ケース | `.webhook_events.json` | `.build_state.queued` | response |
|--------|------------------------|-----------------------|----------|
| 署名不正 / secret 不在 | 変更なし | 変更なし | `401 {"error":"Unauthorized"}` |
| JSON 不正 | 変更なし | 変更なし | `422 {"error":"Validation failed",...}` |
| event が `push` 以外 | `ignored_event` を追記 | 変更なし | `202 {"message":"Webhook ignored","queued":false,"event_id":...}` |
| 対象 branch なし | `ignored_branch` を追記 | 変更なし | `202 {"message":"Webhook ignored","queued":false,"event_id":...}` |
| queue 追加成功 | `queued` を追記 | queue entry を追加 | `202 {"message":"Webhook accepted","queued":true,"event_id":...,"queue_id":...}` |
| queue full | `queue_full` を追記 | 変更なし | `429 {"error":"queue_full"}` |
| event log 追記失敗 / queue 追加前 | 変更なし | 変更なし | `500 {"error":"Internal server error"}` |
| event log 追記失敗 / queue 追加後 | 変更なし | queue entry は残す | `202 {"message":"Webhook accepted","queued":true,"event_log_failed":true,"queue_id":...}` |

queue entry は §16C の queue entry schema を使用し、`trigger:"webhook"`、`requested_by:"webhook"`、`priority:"normal"`、`payload.delivery_id`、`payload.branch`、`payload.sha` を保存する。`X-GitHub-Delivery` が既に pending queue に存在し、同一 branch / sha の場合は重複投入せず、event log に `result:"duplicate"`、既存 `queued_id` を記録し、`202 {"message":"Webhook already queued","queued":true,"queue_id":"<existing>"}` を返す。

**Webhook validation 固定契約：**

| 項目 | 仕様 |
|------|------|
| body size | 最大 1 MiB。超過時は署名検証後に `413`。 |
| branch 抽出 | `ref` が `refs/heads/{branch}` 形式でない場合は `ignored_branch`。 |
| sha | `after` が 40 文字 lowercase hex でない場合は `422`。 |
| repository | `owner.login` と `name` から `owner/name` を作る。欠落時は `422`。 |
| event id | `wh{YYYYMMDDHHmmss}`、衝突時 `-001`。 |
| secret | signature、secret、raw payload は event log、access log、audit log に保存しない。 |

**異常系：**

| 条件 | 応答 / 処理 |
|------|-------------|
| 署名不正 | `401`、イベントログ追記なし、queue なし。 |
| event が `push` 以外 | `202`、`queued=false`、イベントログには `ignored_event` として記録する。 |
| JSON 不正 | `422`、queue なし。 |
| 対象 branch なし | `202`、`queued=false`、イベントログに記録する。 |
| queue 上限 | `429`、イベントログに `queue_full` を記録する。 |
| イベントログ書き込み失敗 | `500`、queue 追加前なら queue しない。queue 追加後なら response に `queued=true` と `event_log_failed=true` を含める。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 正常 push | 署名検証成功、イベントログ追記、queue entry `trigger="webhook"`。 |
| 署名不一致 | `401`、状態差分なし。 |
| 対象外 branch | `202 queued=false`、イベントログのみ。 |
| queue full | `429`、queue 差分なし。 |
| duplicate | queue 追加なし、既存 queue id を返す。 |
| payload oversized | `413`、状態差分なし。 |

### 27.13 Webhook イベントログ / 一覧取得 API

本機能の目的は、受信した GitHub Webhook の監査情報を `.webhook_events.json` に保存し、管理 API、SDK、UI からページング参照できるようにすることである。

owner component は `api` とする。collaborator component は `sdk`、`ui`、`statefile` とする。

**保存 schema：**

`.webhook_events.json` は JSON Lines とし、1 行 1 event を追記する。mode は `600` とする。

| key | 型 | 必須 | 説明 |
|-----|----|------|------|
| `timestamp` | string | 必須 | ISO 8601 UTC。 |
| `delivery_id` | string/null | 必須 | `X-GitHub-Delivery`。 |
| `event` | string | 必須 | GitHub event 名。 |
| `ref` | string/null | 必須 | push ref。 |
| `branch` | string/null | 必須 | `refs/heads/` を除いた branch。 |
| `sha` | string/null | 必須 | push `after`。 |
| `repository` | string/null | 必須 | `owner/repo`。 |
| `build_triggered` | boolean | 必須 | queue 追加済みなら `true`。 |
| `queued_id` | string/null | 必須 | queue id または `null`。 |
| `result` | string | 必須 | `"queued"`, `"duplicate"`, `"ignored_event"`, `"ignored_branch"`, `"queue_full"`, `"error"`。 |

`delivery_id` は 1〜200 文字、`event` は 1〜100 文字、`repository` は `owner/repo` 形式、`sha` は `null` または 40 文字 lowercase hex とする。保存時に request header 全体、署名値、secret、payload 全体を保存してはならない。

**一覧 API：**

`GET /api/webhook-events` は `limit` と `offset` query を受け付ける。`limit` は 1〜1000、既定値 50。`offset` は 0 以上、既定値 0。新しい順で返す。壊れた行は無視し、server log に `WEBHOOK_EVENT_LOG_SKIP_CORRUPT` を出す。

Response は `{ "events": WebhookEventRecord[], "total": N }` とする。SDK `getWebhookEvents(limit,offset)` は `limit` と `offset` を常に query へ送信する。UI は件数、delivery id、event、branch、sha、result、queued id を表示する。

**webhook events 取得固定契約：**

| 項目 | 仕様 |
|------|------|
| 並び順 | `timestamp` 降順、同時刻は file 出現順の逆順。 |
| total | 壊れた行を除外した総件数。 |
| offset | filter 後、並び替え後に適用する。 |
| 壊れた行 | 内容を response、server log に含めない。固定コードだけ出す。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| event 追記 | JSON Lines へ schema 通り保存される。 |
| 一覧取得 | 新しい順、limit/offset が効く。 |
| 壊れた行 | API は継続し、壊れた行を返さない。 |
| duplicate delivery | queue 重複なし、event log は `duplicate`。 |
| total | 壊れた行を除外した件数。 |

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

本機能の目的は、認証不要の `GET /api/health` で、外部監視が Adlaire CI の最低限の稼働状態を確認できるようにすることである。

owner component は `api` とする。collaborator component は `statefile` とする。

**Response：**

```json
{
  "status": "ok",
  "last_build_at": "2026-09-16T00:00:00Z",
  "last_build_status": "success",
  "last_deploy_at": "2026-09-16T00:01:00Z",
  "last_deploy_status": "success",
  "pending_transfers": 0,
  "uptime_seconds": 3600,
  "checks": []
}
```

`status` は `"ok"`、`"degraded"`、`"error"` のいずれかとする。必須状態ファイル破損がある場合は `degraded`、API process が応答できるが重大な read error がある場合は `error` とする。HTTP status は、API 自体が response を生成できる限り `200` とし、JSON 生成不能などの場合だけ `500` とする。

**読み取り元：**

`.build_status.json` を第一参照元とし、不在時は `.build_history` と `.pending_transfers` から算出する。`.build_status.json` 破損時は自動修復せず、`checks[]` に `build_status_corrupt` を含める。

**health checks 固定契約：**

| check | 条件 |
|-------|------|
| `build_status_missing` | `.build_status.json` 不在で fallback 算出した。 |
| `build_status_corrupt` | `.build_status.json` parse 失敗。 |
| `pending_transfers_read_error` | `.pending_transfers` 読取失敗。 |
| `notify_pending_read_error` | `.notify_pending` 読取失敗。 |
| `runner_stale` | `running=true` かつ `last_started_at` が 24 時間より古い。 |

`checks[]` は上表の順で返す。`status` は checks が空なら `"ok"`、read error または stale があれば `"degraded"`、response 生成不能だけ `"error"` とする。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 正常 | HTTP 200、`status="ok"`。 |
| pending transfer あり | `pending_transfers` に件数、`status="degraded"`。 |
| build status 破損 | HTTP 200、`status="degraded"`、checks に記録。 |
| read error | HTTP 200 または 500 の条件が仕様通り。 |
| stale runner | checks に `runner_stale`。 |

### 27.17 ビルドログ重大度フィルター

本機能の目的は、`GET /api/logs/search` と UI ログビューアで重大度別に build log を絞り込めるようにすることである。

owner component は `api` とする。collaborator component は `sdk`、`ui`、`archive`、`statefile` とする。

**入力：**

`GET /api/logs/search` は既存 query に加えて `level` を受け付ける。`level` の許容値は `"info"`、`"warn"`、`"warning"`、`"error"`、`"debug"` とし、大文字小文字は区別しない。正規化後は `"INFO"`、`"WARNING"`、`"ERROR"`、`"DEBUG"` とする。`warn` は `"WARNING"` と同義とする。不正値は `422`。

**検索対象：**

`.build_logs/{id}.json.stdout`、`stderr`、`warnings`、`error`、archive log を対象とする。行頭が `[WARN]` または `[WARNING]` の行は WARNING、`[ERROR]` または stderr の非空行は ERROR、`[DEBUG]` は DEBUG、それ以外は INFO と分類する。

**ログ行分類固定契約：**

| 入力 | level |
|------|-------|
| stderr の非空行 | `ERROR` |
| stdout / warnings の `[ERROR]` prefix | `ERROR` |
| stdout / warnings の `[WARN]` または `[WARNING]` prefix | `WARNING` |
| stdout / warnings の `[DEBUG]` prefix | `DEBUG` |
| 上記以外 | `INFO` |

検索結果は `{build_id, level, source, line_number, message}` とし、`source` は `"stdout"`、`"stderr"`、`"warnings"`、`"error"` のいずれかとする。`line_number` は 1 始まり、配列項目は配列内 index + 1 とする。

**SDK / UI：**

SDK `searchLogs(q,from,to,level)` は `level` 指定時だけ query に送信する。UI は INFO / WARNING / ERROR / DEBUG の filter control を提供し、選択なしでは全件を表示する。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| `level=warn` | WARNING 行だけ返る。 |
| `level=error` | ERROR 行だけ返る。 |
| 不正 level | `422`。 |
| archive log | 通常 log と同じ分類で検索される。 |
| stderr | prefix なしでも ERROR。 |
| line number | 1 始まりで返る。 |

### 27.18 ブランチ設定の動的変更 API

本機能の目的は、監視対象 branch / target / deploy target を `.branch_config` で管理し、API 経由で変更できるようにすることである。

owner component は `api` とする。collaborator component は `runner`、`statefile` とする。

**API：**

`GET /api/branch-config` は `.branch_config` が存在する場合 `{"source":"file","branches":[...]}`、不在の場合 `{"source":"default","branches":[...]}` を返す。`POST /api/branch-config` は `{ "branches": BranchTargetRecord[] }` を受け取る。

**保存仕様：**

永続ファイルの key は必ず `branch_targets` とする。API request / response で `branches` を使う場合も保存前に `branch_targets` へ変換する。空配列を受け取った場合は `.branch_config` を削除し、default 復帰とする。

**branch config 正規化固定契約：**

| 項目 | 仕様 |
|------|------|
| branch 重複 | 同一 `branch` の重複は `422`。 |
| deploy target id | branch 内で一意。重複は `422`。 |
| target_files | 存在する場合は §27.21 の正規化を適用する。 |
| 保存順 | branch 名昇順、同一 branch 内 deploy target id 昇順で保存する。 |
| 削除 | POST empty で `.branch_config` を削除した後、`.config_log` に default 復帰を記録する。 |

**検証：**

各 branch target は `branch`、`target_file`、`sha_file`、`src`、`out`、`deploy_targets` を必須とする。`target_file` は相対パスで `..` 禁止、`sha_file` / `src` / `out` / `dest_dir` は絶対パス、deploy target は最大 20 件、branch target は最大 50 件とする。

**runner 取り込み：**

`components/runner.go` は起動ごとに `.branch_config` を読む。API 更新後、runner 再起動は不要だが、既に実行中の runner へは反映しない。次回起動から反映する。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| GET default | `.branch_config` 不在で `source="default"`。 |
| POST valid | `.branch_config.branch_targets` として保存、config log 追記。 |
| POST empty | `.branch_config` 削除、default 復帰。 |
| 相対 `sha_file` | `422`、状態差分なし。 |
| branch 重複 | `422`、状態差分なし。 |
| POST empty log failure | `.branch_config` は削除済み、response は `500`。 |

### 27.19 週次ビルドサマリー Webhook

本節の主本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27.19 を正とする。owner component は `runner`、collaborator component は `api`、`statefile` とする。

### 27.20 設定変更の詳細 diff 記録

本機能の目的は、設定変更 API が何を変更したかを `.config_log` に機械可読 diff と人間可読 diff の両方で残すことである。

owner component は `api` とする。collaborator component は `statefile` とする。

**対象 API：**

`.config_log` を write する全 API を対象とする。少なくとも `POST /api/config`、`POST /api/log-level`、`POST /api/notify-config`、`POST /api/repo-config`、`POST /api/branch-config`、`POST /api/webhook-config`、`POST /api/pat-update`、schedule 系 API、maintenance、access-control、hooks、alert-rules、tag-rules、pipeline-config、notes、smtp-config、dashboard-layout、snapshot delete を含む。

**ログ schema：**

各行は §22.0c `.config_log` schema に従う。`diff` は `{key:[before,after]}`、`diff_text` は 1 行以上の文字列とする。差分がない場合、対象 API は状態ファイルを書かず、`.config_log` も追記せず、response は `{ "message": "No changes" }` とする。

**マスク条件：**

キー名に `password`、`token`、`secret`、`pat`、`smtp_password` を含む値は before / after とも `"***"` に置換する。配列や object の内部 key も同じ規則で再帰的にマスクする。

**diff_text 形式：**

`{key}: {before} -> {after}` を key 名昇順で 1 行ずつ連結する。値は JSON 表現とし、secret は `"***"` とする。複数行値は `\n` escape した 1 行 JSON string とする。

**diff 生成固定契約：**

| 項目 | 仕様 |
|------|------|
| 比較対象 | 正規化後、既定値 merge 後の before / after object。保存対象外 key、response 専用 key、password 平文は比較対象に含めない。 |
| object | dot path で再帰比較する。例: `summary.hour`。 |
| array | index 比較ではなく配列全体を JSON 値として比較する。並び順が仕様上正規化される配列は正規化後に比較する。 |
| secret key | key path のいずれかに `password`、`token`、`secret`、`pat`、`smtp_password` を含む場合、before / after を `"***"` にする。 |
| no-op | diff が空の場合は状態ファイル、secret file、`.config_log` を変更しない。 |
| 追記順 | 対象状態ファイル保存後に `.config_log` を追記する。`.config_log` 失敗時は `500`、保存済み状態は巻き戻さない。 |

`.config_log` record の `target` は endpoint 固定名、`actor` は管理 session なら `"admin"`、API token なら token id とする。`request body` 全体、HTTP header、cookie、secret 平文を保存してはならない。

**diff 対象 endpoint 固定名：**

| endpoint | target |
|----------|--------|
| `POST /api/config` | `server_config` |
| `POST /api/notify-config` | `notify_config` |
| `POST /api/repo-config` | `repo_config` |
| `POST /api/branch-config` | `branch_config` |
| `POST /api/webhook-config` | `webhook_config` |
| `POST /api/smtp-config` | `smtp_config` |
| `POST /api/pipeline-config` | `pipeline_config` |
| その他 `.config_log` 対象 | method と path から `/api/` prefix を除き、`/` と `-` を `_` に置換した固定名。 |

diff 生成は状態保存前に memory 上で完了させる。diff 生成に失敗した場合は状態ファイルを書かない。`.config_log` 追記に失敗した場合は保存済み状態を巻き戻さず、response は `500` とする。

**異常系：**

| 条件 | 処理 |
|------|------|
| config log 追記失敗 | 対象状態ファイルの更新を失敗扱いにし、`500` を返す。 |
| diff 生成失敗 | 状態ファイルを書かず `500`。 |
| secret マスク漏れ検出 | 実装不合格。該当 API は完了扱いにしない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 通常変更 | `.config_log` に diff と diff_text が残る。 |
| secret 変更 | 値は `"***"` だけ保存される。 |
| 変更なし | 状態ファイルも config log も更新しない。 |
| 複数 key | key 昇順で diff_text を生成する。 |
| array 正規化 | 正規化後に同一なら no-op。 |
| config log 失敗 | 対象状態は保存済み、response は `500`。 |
| diff failure | 状態差分なしで `500`。 |
| target name | endpoint から固定 target 名が生成される。 |

**§27.21〜§27.38 / §27.42〜§27.47 機能別実装完全性固定契約：**

§27.21〜§27.38、§27.42〜§27.47 の各機能は、個別節の本文に加えて下表を満たした場合だけ実装完了とする。§27.38a は §27.21〜§27.38 の runner 拡張を横断検証する補足契約として扱う。

| 節 | 機能 | 入力 | 出力 | 状態ファイル / 外部副作用 | 失敗時副作用 | 必須 fixture |
|----|------|------|------|---------------------------|--------------|--------------|
| §27.21 | 複数ファイル監視 | `.branch_config.branch_targets[].target_files`、GitHub content SHA または local SHA。 | `changed_targets[]`、target 単位 SHA cache、build log。 | 成功時だけ該当 target SHA cache を更新する。 | SHA 部分失敗では build を開始せず、成功取得済み cache も更新しない。 | 1 件変更、複数変更、変更なし、不正 path、force build、SHA 部分失敗。 |
| §27.22 | pipeline YAML | `.pipeline.yml` または `.pipeline_config.inline_yaml`。 | `pipeline_steps[]`、step stdout/stderr、build status。 | step を定義順に実行し、build 終了時に定義順で保存する。 | parse / command 不正では build を開始しない。required step 失敗で後続 required step を実行しない。 | 成功、required 失敗、optional 失敗、禁止 YAML、secret mask、timeout。 |
| §27.23 | local watch | `.server_config.watch_mode`、local `src` 配下 Markdown。 | `.local_watch_state.json`、trigger `local_watch`。 | local mode では GitHub API / PAT を呼ばず、成功時だけ state を置換する。 | file read 失敗は build なし。dry-run は state を作成 / 更新しない。 | 初回、変更なし、1 file 変更、token 不在、file 削除、out 除外。 |
| §27.24 | tag filter | `.server_config.tag_filter`、GitHub tags refs。 | `matched_tags[]`、skip status。 | tag 一致時だけ build。tag 不一致 skip は `.build_status.json` だけ更新する。 | tags API 最終失敗は build なし。local mode 併用は終了コード `2`。 | tag 一致、不一致、patterns 空、API 失敗、prefix pattern、local mode 併用。 |
| §27.25 | build cache | `--cache-dir`、cache 設定、input / deps SHA。 | cache entry、page cache、`[REPORT]` cache counts。 | hit 時は変換結果を再利用し、miss 成功時だけ entry を書く。 | cache read/write 失敗は build を成功可能にし、WARN / report に残す。 | 2 回目 hit、1 file miss、theme 変更、cache 破損、entry 不一致、write failure。 |
| §27.26 | parallel deploy | `.server_config.deploy_parallelism`、deploy targets。 | `target_results[]`、`.pending_transfers`。 | target ごとに並列転送し、result は設定順で保存する。 | worker 内部失敗は該当 target failure。他 target は継続する。 | 3 target/2 worker、1 target 失敗、parallelism 1、result order、pending duplicate。 |
| §27.27 | build hooks | `.hooks`、pre/post hook command_args。 | hook log、build log hook result。 | pre は build 前、post は build 後に id 昇順で実行する。 | pre abort で build 本体を開始しない。post 失敗は build status を変更しない。 | pre success、pre abort、post failure、timeout、delete hook、secret mask。 |
| §27.28 | dependency tracking | builder dependency manifest、Markdown link / asset reference。 | `.dependency_manifest.json`、affected target 判定。 | build 成功時だけ manifest を更新する。 | manifest 破損は full build 扱い。dry-run は更新しない。 | direct dep、shared asset、deleted dep、manifest corrupt、dry-run no write。 |
| §27.29 | remote build | remote build 設定、ssh target、artifact path。 | remote artifact metadata、build log remote section。 | remote command と artifact fetch を行い、検証成功時だけ deploy / history へ進む。 | remote timeout / checksum mismatch は failure とし、secret / command credential を保存しない。 | remote success、timeout、checksum mismatch、artifact missing、secret mask。 |
| §27.30 | approval flow | approval policy、pending request、approve/reject API。 | `.approval_queue`、history `approval_*`、通知。 | pending 作成後、承認時だけ queue / build へ進む。 | expired / rejected は build を開始しない。history append 失敗時も approval record は残す。 | pending、approve、reject、expired、duplicate approve、notify failure。 |
| §27.31 | branch env | branch env 設定、pipeline / hook env。 | merged env、masked log。 | runner env → branch env → step/hook env の順で上書きする。 | env key 不正は保存不可 / runner 設定エラー。secret はログに出さない。 | merge order、invalid key、secret stdout、branch missing、override。 |
| §27.32 | build notification | notify config、build event、channel 設定。 | notify payload、`.notify_log`、`.notify_pending`。 | event ごとに payload を生成し、送信結果を記録する。 | retry 対象失敗は pending。secret は payload/log に含めない。 | success notify、failure notify、webhook 5xx、command timeout、disabled、secret mask。 |
| §27.33 | build trends | build duration samples、history/log。 | `.build_trends.json`、trend stats response。 | build 完了時に sample を追加し、上限件数で trim する。 | trend 保存失敗は build 成否を反転しない。破損時は再集計契約に従う。 | sample append、trim、median/p95、corrupt rebuild、save failure。 |
| §27.34 | build chain | chain config、upstream build result。 | chain execution log、queued child build。 | chain 条件一致時だけ次 build を queue する。 | chain config 破損は chain 無効として通常 build は継続する。 | chain success、condition mismatch、loop detect、queue full、config corrupt。 |
| §27.35 | priority queue | queue request、priority、created_seq。 | sorted queue、queue API response。 | queue 保存時に priority / created_seq を固定する。 | queue full は `429`、既存 queue を変更しない。 | high priority、same priority FIFO、queue full、clear queue、invalid priority。 |
| §27.36 | failure classification | build failure evidence、logs/status。 | `failure_category`、`failure_evidence[]`、API filter。 | finalizer で分類し、log/history/status に保存する。 | 分類不能でも build 結果は保持し、category `unknown` とする。 | pipeline timeout、api failure、deploy failure、unknown、filter。 |
| §27.37 | execution environment | runner host/runtime info、builder version。 | environment record、build log/status。 | build 開始時に取得し、log に保存する。 | 取得失敗は `unknown` を保存し build は継続する。secret / env 全量は保存しない。 | normal、hostname unavailable、builder version timeout、secret env absent。 |
| §27.38 | duration anomaly | duration samples、anomaly config。 | anomaly flag、alert / notify、stats。 | build 完了時に閾値判定し、該当時だけ alert/notify を作る。 | samples 不足では判定しない。通知失敗は build 成否を反転しない。 | normal、avg anomaly、p95 anomaly、insufficient samples、notify failure。 |
| §27.42 | trigger API scope | API token scopes、endpoint group。 | allow/deny decision、audit。 | scope 一致時だけ endpoint 実行。 | scope 不足は `403`、endpoint 固有処理なし。 | read allowed、trigger denied、admin allowed、health no scope、audit deny。 |
| §27.43 | API key management | token label/scopes/expires_at。 | `.api_tokens` record、token 本体 1 回 response。 | token 作成 / 失効 / 認証成功時に状態と audit を更新する。 | audit 失敗時は token 本体を返さない。token 破損 state は自動再生成しない。 | create、revoke、expired auth、scope auth、audit failure、secret non-persistence。 |
| §27.44 | audit log | 認証、権限、設定、token、build trigger event。 | `.audit_log` JSON Lines、audit API response。 | 監査対象操作の成否確定後に追記する。 | 追記失敗は対象操作を `500` 扱い。ただし保存済み状態は個別契約どおり戻さない。 | password change、permission denied、secret input、corrupt line、filter、append failure。 |
| §27.45 | session timeout | `.server_config.session_timeout_seconds`。 | 新規 session `expires_at`。 | session 発行直前の設定で期限を計算する。 | 範囲外は `422`。既存 session の期限は変更しない。 | 300 秒、既存 session 維持、範囲外、key 不在、TOTP login、秒精度。 |
| §27.46 | TOTP | setup secret、ticket、TOTP code。 | `.totp_secret`、ticket/session response、audit/access log。 | secret / ticket はメモリと 1 回 response に限定し、成功時だけ永続状態を更新する。 | code 不一致 / replay は token を返さない。audit 失敗時も secret 平文を出さない。 | setup、login、replay、secret one-time、ticket reuse、disable、window、全角 code。 |
| §27.47 | API rate limit | rate policy、remote addr、actor key。 | `.api_rate_state`、`429`、state summary。 | key 群を同一 lock で判定 / 更新する。 | 上限超過では count を増やさず endpoint 固有処理を行わない。audit 失敗時は `500`。 | under limit、over IP、over token、disabled、policy update、audit failure、window reset。 |

**§27.21〜§27.38 / §27.42〜§27.47 API / SDK / UI 連動固定契約：**

| 節 | API | SDK | UI |
|----|-----|-----|----|
| §27.21 | branch config API と status/history/log に反映する。 | `getBranchConfig()` / `setBranchConfig()` は `target_files` を削除しない。 | リポジトリ情報 panel で複数 target を表示 / 保存する。 |
| §27.22 | pipeline config API。 | `getPipelineConfig()` / `setPipelineConfig()`。 | 設定 panel で pipeline config を表示 / 保存し、shell 文字列へ変換しない。 |
| §27.23 | `GET/POST /api/config` の `watch_mode`。 | `getConfig()` / `setConfig()`。 | 設定 panel で `github` / `local` を選択する。 |
| §27.24 | `GET/POST /api/config` の `tag_filter`。 | `getConfig()` / `setConfig()`。 | 設定 panel で tag filter を表示 / 保存する。 |
| §27.25 | `GET/POST /api/config` の cache key と builder report。 | `getConfig()` / `setConfig()`。 | 設定 panel と build result 表示で cache counts を表示する。 |
| §27.26 | `GET/POST /api/config` の `deploy_parallelism` と status/history。 | `getConfig()` / `setConfig()`。 | 設定 panel で parallelism を表示 / 保存し、結果は履歴/logで表示する。 |
| §27.27 | hooks API。 | hook methods。 | フック panel で command_args を 1 行 1 引数として表示 / 保存する。 |
| §27.28 | endpoint 追加なし。manifest は runner/builder 内部状態。 | SDK method 追加なし。 | UI 操作追加なし。build result で依存情報を表示してよい。 |
| §27.29 | remote build config は config / pipeline 系 API に含める。 | `getConfig()` / `setConfig()` または pipeline method。 | 設定 panel で remote build 設定を表示 / 保存する。 |
| §27.30 | approvals API。 | approval methods。 | 承認待ち panel で approve / reject を操作する。 |
| §27.31 | branch config または config API に含める。 | branch/config method。 | リポジトリ情報 panel で branch env を表示 / 保存する。 |
| §27.32 | notify config / notify log API。 | notify methods。 | 通知設定 panel で channel、payload、履歴を表示する。 |
| §27.33 | stats/build trends API。 | `getBuildTrends()`。 | 統計 panel で trend を表示する。 |
| §27.34 | build chain config API。 | `getBuildChainConfig()` / `setBuildChainConfig()`。 | 設定またはリポジトリ情報 panel で chain を表示 / 保存する。 |
| §27.35 | queue API。 | `getQueue()` / `clearQueue()` と build trigger methods。 | 手動実行 panel で priority 反映後 queue を表示する。 |
| §27.36 | history/stats API の filter と response に含める。 | history / stats methods は category を保持する。 | 履歴 panel で failure category filter を表示する。 |
| §27.37 | output meta / status / history log に含める。 | `getOutputMeta()` / `getStatus()` / `getHistoryLog()`。 | システム情報または履歴 detail に表示する。 |
| §27.38 | config / stats / dashboard alerts に含める。 | config / stats / dashboard methods。 | 統計 panel と dashboard alert で anomaly を表示する。 |
| §27.42 | endpoint ごとの scope 判定。 | SDK は token scope を推測せず API error を返す。 | UI は `403` を権限不足として表示し、logout しない。 |
| §27.43 | token API。 | token methods。 | API token 管理 panel で発行 token を 1 回だけ表示する。 |
| §27.44 | audit log API。 | `getAuditLog()`。 | 監査ログ panel で filter 表示し、secret は表示しない。 |
| §27.45 | config API。 | `setConfig({session_timeout_seconds})`。 | セキュリティ panel で session timeout を表示 / 保存する。 |
| §27.46 | TOTP/auth API。 | TOTP/auth methods。 | セキュリティ/login panel で one-time secret/ticket flow を扱う。 |
| §27.47 | rate limit API。 | `getApiRateLimit()` / `setApiRateLimit()`。 | セキュリティ panel で policy と state summary を表示する。 |

API / SDK / UI のいずれも、上表に存在しない補完 endpoint、補完 method、補完 UI 操作を追加してはならない。個別節が endpoint 追加なしとする機能は、runner / builder の内部挙動または既存 response field の範囲で実装する。

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

本機能の目的は、本番向けなど approval_required な target の build / deploy を人間承認後にだけ実行することである。

owner component は `api` とする。collaborator component は `runner`、`sdk`、`ui`、`statefile` とする。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 key | `branch_targets[].approval_required` |
| 状態 | `.approval_queue` JSON Lines。mode `600`。 |
| API | `GET /api/approvals`、`POST /api/approvals/{id}/approve`、`POST /api/approvals/{id}/reject` |
| status | `"pending"`、`"approved"`、`"rejected"`、`"expired"`。 |
| timeout | `.server_config.approval_timeout_seconds`。既定値 86400。 |

**正常系：**

1. runner は差分検出後、approval_required target について build を開始せず approval entry を作成する。
2. `.notify_config` に従い approval request 通知を送信する。
3. API approve 後、queue entry に `trigger="approval"` を追加する。
4. reject 後は build せず `.build_history.status="approval_rejected"` を記録する。
5. timeout 超過 entry は runner 起動時に `expired` へ更新する。

**approval queue / notify 固定契約：**

| 項目 | 仕様 |
|------|------|
| pending id | `appr{YYYYMMDDHHmmss}`、同秒衝突時は `-001` から連番。既存 id は再利用しない。 |
| queue payload | approve で追加する queue entry は `trigger:"approval"`、`priority:"normal"`、`requested_by` は actor id、`payload.approval_id` を含める。 |
| 通知 payload | `{event:"approval_required", approval_id, branch, sha, target, expires_at}`。secret、token、path secret は含めない。 |
| 重複 pending | 同一 branch / sha / target の pending がある場合、新規通知は送らず既存 id を返す。 |
| expired 更新 | runner 起動時に期限超過 pending をすべて処理し、created_at 昇順で expired record を追記する。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| entry 不在 | API は `404`。 |
| pending 以外への approve/reject | `409`。 |
| 通知失敗 | approval entry は残し、`.notify_pending` に追記する。 |
| queue full | approve API は `429`。approval status は pending のまま。 |

**SDK / UI：**

SDK は `getApprovals()`、`approveBuild(id)`、`rejectBuild(id)` を提供する。UI は pending 件数、branch、sha、target、created_at、expires_at、approve/reject 操作を表示する。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| approval required | build せず pending 作成。 |
| approve | queue 追加、trigger approval。 |
| reject | build なし、history 記録。 |
| timeout | expired、build なし。 |

**`.approval_queue` record schema：**

| キー | 型 | 必須 | 許容値 |
|------|----|------|--------|
| `id` | string | 必須 | `appr{YYYYMMDDHHmmss}`、衝突時 `-001`。 |
| `status` | string | 必須 | `"pending"`、`"approved"`、`"rejected"`、`"expired"`。 |
| `branch` | string | 必須 | branch target 名。 |
| `sha` | string | 必須 | 40 文字 lowercase hex。 |
| `target` | string | 必須 | branch target id または target file。 |
| `created_at` | string | 必須 | UTC ISO 8601。 |
| `expires_at` | string | 必須 | UTC ISO 8601。 |
| `decided_at` | string/null | 必須 | approve / reject / expire 時刻。 |
| `decided_by` | string/null | 必須 | 管理 session は `"admin"`、API token は token id。 |
| `queue_id` | string/null | 必須 | approve で追加した queue id。 |
| `reason` | string/null | 必須 | reject 理由または expire 理由。 |

`.approval_queue` は JSON Lines append-only とする。同一 id の最新 record を有効状態として扱い、古い record は監査履歴として残す。`GET /api/approvals` は id ごとに最新 record だけを返し、`created_at` 降順、同時刻は id 昇順で並べる。壊れた行は無視し、response に含めない。

**approval 状態遷移固定契約：**

| 現在 status | 操作 | 次 status | 副作用 |
|-------------|------|-----------|--------|
| `pending` | approve | `approved` | `.build_state.queued[]` に `trigger:"approval"` entry を追加し、`queue_id` を保存する。 |
| `pending` | reject | `rejected` | `.build_history` に `status:"approval_rejected"` を追記する。queue は追加しない。 |
| `pending` | timeout | `expired` | `.build_history` に `status:"approval_expired"` を追記する。queue は追加しない。 |
| `approved` | approve / reject | 変更なし | `409 {"error":"Conflict"}`。 |
| `rejected` | approve / reject | 変更なし | `409 {"error":"Conflict"}`。 |
| `expired` | approve / reject | 変更なし | `409 {"error":"Conflict"}`。 |

**approval 更新順：**

| 操作 | 更新順 | 失敗時 |
|------|--------|--------|
| pending 作成 | `.approval_queue` lock → 重複確認 → pending record append → 通知送信 → `.notify_log` / `.notify_pending` 更新 | 通知失敗でも pending は残す。pending append 失敗時は build を開始せず runner failure。 |
| approve | `.approval_queue` lock → 最新 pending 確認 → `.build_state` lock → queue append → approved record append → response | queue full は `429`、approval は pending のまま。approved append 失敗時は `500`、queue 追加済み entry は巻き戻さない。 |
| reject | `.approval_queue` lock → 最新 pending 確認 → rejected record append → `.build_history` append → response | history append 失敗時は `500`。rejected record は巻き戻さない。 |
| timeout | runner 起動時に `.approval_queue` lock → expires_at 超過 pending を expired record append → `.build_history` append | history append 失敗時も expired record は残し、runner は ERROR を出して継続する。 |

pending 作成時の重複判定は `branch`、`sha`、`target` が同一で、最新 status が `pending` の record とする。重複時は新規 record を作成せず、既存 pending id を使用する。approve / reject API は body を受け付けない。reject reason は初期実装では固定 `"rejected"` とする。

**approval fixture 固定：**

| fixture | 入力 | 期待結果 |
|---------|------|----------|
| approval-create | approval_required target に差分 | build なし、pending record、通知成功または pending。 |
| approval-duplicate | 同一 branch/sha/target を再検出 | pending 重複作成なし。 |
| approval-approve | pending approve | queue 追加、approved record、queue_id 保存。 |
| approval-reject | pending reject | rejected record、history `approval_rejected`。 |
| approval-timeout | expires_at 超過 | expired record、history `approval_expired`。 |
| approval-queue-full | max_size 到達時 approve | `429`、status pending 維持。 |
| approval duplicate notify | 重複時は通知を送らない。 |
| approved append failure | queue は残り、API は `500`。 |

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

本機能の目的は、外部システムが最小権限で build を開始できる API token を発行できるようにすることである。

**scope：**

| scope | 許可 |
|-------|------|
| `trigger` | `POST /api/build`、`POST /api/build/force` のみ。 |
| `read` | 読み取り endpoint のみ。 |
| `operate` | 運用操作。 |
| `config` | 設定変更。 |
| `admin` | 管理操作。 |

`trigger` scope は cancel、queue clear、config、token、audit、PAT 更新を許可しない。

**endpoint group 対応：**

| scope | 許可 endpoint |
|-------|---------------|
| `trigger` | `POST /api/build`, `POST /api/build/force` |
| `read` | `GET /api/status`, `GET /api/logs`, `GET /api/logs/search`, `GET /api/logs/export`, `GET /api/history`, `GET /api/history/export`, `GET /api/history/{id}/log`, `GET /api/history/{id}/comment`, `GET /api/sysinfo`, `GET /api/health`, `GET /api/schedule`, `GET /api/notify-config`, `GET /api/notify-log`, `GET /api/config`, `GET /api/config-log`, `GET /api/pat-status`, `GET /api/stats`, `GET /api/stats/timeline`, `GET /api/stats/build-duration`, `GET /api/stats/build-trends`, `GET /api/output-meta`, `GET /api/repo-info`, `GET /api/branch-config`, `GET /api/backup`, `GET /api/dashboard`, `GET /api/diagnostics`, `GET /api/rate-limit`, `GET /api/disk-usage`, `GET /api/webhook-events`, `GET /api/webhook-config`, `GET /api/snapshots`, `GET /api/snapshots/{id}/download`, `GET /api/maintenance`, `GET /api/access-control`, `GET /api/hooks`, `GET /api/hooks/{id}/log`, `GET /api/alert-rules`, `GET /api/tag-rules`, `GET /api/pipeline-config`, `GET /api/build-chain-config`, `GET /api/notes`, `GET /api/smtp-config`, `GET /api/queue`, `GET /api/approvals`, `GET /api/dashboard-layout` |
| `operate` | `POST /api/build/cancel`, `GET /api/build/stream`, `POST /api/notify-test`, `POST /api/notify/weekly-summary`, `POST /api/pat-verify`, `POST /api/circuit-breaker/reset`, `DELETE /api/queue`, `POST /api/smtp-test`, `POST /api/verify-output`, `POST /api/history/{id}/rollback`, `POST /api/approvals/{id}/approve`, `POST /api/approvals/{id}/reject` |
| `config` | `POST /api/schedule/interval`, `POST /api/schedule/pause`, `POST /api/schedule/resume`, `POST /api/schedule/allowed-hours`, `POST /api/schedule/force-interval`, `POST /api/schedule/cooldown`, `POST /api/notify-config`, `POST /api/config/validate`, `POST /api/config`, `POST /api/log-level`, `POST /api/pat-update`, `POST /api/repo-config`, `POST /api/branch-config`, `POST /api/restore`, `POST /api/webhook-config`, `DELETE /api/snapshots/{id}`, `POST /api/maintenance/enable`, `POST /api/maintenance/disable`, `POST /api/access-control`, `POST /api/hooks`, `DELETE /api/hooks/{id}`, `POST /api/alert-rules`, `DELETE /api/alert-rules/{id}`, `POST /api/tag-rules`, `DELETE /api/tag-rules/{id}`, `POST /api/pipeline-config`, `POST /api/build-chain-config`, `POST /api/notes`, `POST /api/smtp-config`, `POST /api/dashboard-layout`, `POST /api/history/{id}/comment`, `POST /api/history/{id}/flag`, `POST /api/history/{id}/tags`, `POST /api/logs/cleanup`, `POST /api/logs/archive` |
| `admin` | `GET /api/access-log`, `GET /api/api-access-log`, `GET /api/audit-log`, `GET /api/api-rate-limit`, `POST /api/api-rate-limit`, `GET /api/sessions`, `POST /api/sessions/revoke-all`, `GET /api/auth/totp-status`, `POST /api/auth/totp-setup`, `POST /api/auth/totp-confirm`, `DELETE /api/auth/totp`, `GET /api/tokens`, `POST /api/tokens`, `DELETE /api/tokens/{id}` |

管理 session は上表に関係なく全 endpoint を許可する。API token が複数 scope を持つ場合は、いずれか 1 つの scope が endpoint に一致すれば許可する。`POST /api/login`、`POST /api/login/totp`、`POST /api/logout`、`POST /api/change-password` は API token scope 判定の対象外とし、API token では使用できない。`POST /api/webhook` は GitHub Webhook secret 検証専用であり、API token では使用できない。上表に存在しない endpoint は §22.0e と本表へ追加されるまで API token では許可してはならない。

**正常系：**

1. `POST /api/tokens` は `scopes:["trigger"]` を受け付ける。
2. 発行 token は `Authorization: Bearer <token>` で使用する。
3. `trigger` token による build は `.audit_log` と `.build_logs/{id}.json.trigger_actor` に token id を記録する。

**判定順：**

1. Bearer token の形式を検査する。`act_` prefix の場合は API token、64 文字 hex の場合は session token として扱う。どちらにも一致しない場合は `401`。
2. API token の hash、期限、失効状態を検証する。
3. endpoint group 対応表で scope を判定する。
4. 許可時だけ endpoint 処理へ進む。
5. 拒否時は `.audit_log` に `permission_denied`、`actor_type:"api_token"`、`actor_id:{token_id}`、`target_type:"endpoint"`、`target_id:{METHOD + " " + path}`、`result:"denied"` を追記してから `403` を返す。

`permission_denied` の監査ログ追記に失敗した場合は `500` を返し、対象 endpoint の処理は実行しない。

**認証・rate limit 組み合わせ順：**

| 段階 | 処理 |
|------|------|
| 1 | route match と method check を行う。未定義は `404` / `405` を返し、API token scope 判定は行わない。 |
| 2 | 認証不要 endpoint か判定する。`GET /api/health` と `POST /api/webhook` は API token scope 判定対象外。 |
| 3 | login endpoint は §27.47 の login rate limit を先に判定する。 |
| 4 | 認証必須 endpoint は Bearer token を検証し、actor を確定する。 |
| 5 | API token の場合だけ scope 判定を行う。管理 session は scope 表を参照しない。 |
| 6 | 認証後 endpoint の rate limit を §27.47 に従って判定する。 |
| 7 | endpoint 固有処理へ進む。 |

scope 判定前に endpoint 固有の request body parse、状態ファイル更新、外部送信を実行してはならない。

**scope 判定固定条件：**

| 条件 | 判定 |
|------|------|
| session token | scope 表を参照せず許可する。 |
| API token の `scopes` が空 | token record 破損として扱い `500`。 |
| API token の `scopes` に未知値 | token record 破損として扱い `500`。 |
| endpoint が複数 group に現れる | §27.42 の表で先に定義された group を採用する。 |
| path parameter 付き endpoint | 正規化 path pattern で scope 判定する。例: `DELETE /api/tokens/abc` は `DELETE /api/tokens/{id}`。 |
| query string | scope 判定には使用しない。method と path pattern だけで判定する。 |
| `GET /api/health` | 認証不要。API token scope 判定を行わない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| trigger token で `POST /api/build` | `202`。 |
| trigger token で `GET /api/status` | `403`。`read` を併用した場合のみ成功。 |
| trigger token で config | `403`。 |
| 未定義 endpoint | token scope に関係なく `404` または `405`。 |
| 複数 scope | いずれかに一致する endpoint だけ成功。 |
| path parameter endpoint | pattern 正規化後の scope で判定する。 |
| query 付き endpoint | query を除外して scope 判定する。 |
| scope 前 body | 権限不足時に request body validation や状態更新を行わない。 |

### 27.43 API キー管理

本機能の目的は、API key の発行、一覧、失効、期限、scope を実装し、key 本体を保存しないことである。

**API 権限：**

| endpoint | 必要認証 |
|----------|----------|
| `GET /api/tokens` | 管理 session または `admin` scope API token |
| `POST /api/tokens` | 管理 session または `admin` scope API token |
| `DELETE /api/tokens/{id}` | 管理 session または `admin` scope API token |

`admin` scope API token で新しい `admin` scope API token を発行してよい。発行者 token と発行対象 token は別 record とし、親子関係は保存しない。

**正常系：**

1. `POST /api/tokens` は `label`、`scopes`、`expires_at` を受け取る。
2. token 本体は `act_` + 32 byte 相当のランダム文字列とする。
3. `.api_tokens` には `sha256(token)` の lowercase hex だけを保存する。
4. response の `token` は作成時 1 回だけ返す。`GET /api/tokens` では返さない。
5. 認証時は hash 一致、`revoked_at == null`、`expires_at == null または now < expires_at` を満たす token だけ有効とする。
6. 作成、認証成功、失効、期限切れ拒否を `.audit_log` へ記録する。ただし token 本体は記録しない。

**`.api_tokens` record schema：**

| キー | 型 | 必須 | 許容値 |
|------|----|------|--------|
| `id` | string | 必須 | `tok` + 6 桁以上の数字。 |
| `label` | string | 必須 | 1〜64 文字。前後空白は保存前に除去し、空になれば `422`。 |
| `token_hash` | string | 必須 | `sha256(token)` の 64 文字 lowercase hex。 |
| `scopes` | array[string] | 必須 | `read`、`trigger`、`operate`、`config`、`admin` の 1〜5 件。保存時は重複除去し固定順へ正規化する。 |
| `created_at` | string | 必須 | UTC ISO 8601。 |
| `last_used_at` | string/null | 必須 | 認証成功後に更新する。 |
| `expires_at` | string/null | 必須 | `null` または UTC ISO 8601。 |
| `revoked_at` | string/null | 必須 | 失効時に UTC ISO 8601 を保存する。 |

`.api_tokens` に未知 key、必須 key 不足、型不一致、未知 scope、空 scopes、hash 形式不正、時刻形式不正がある場合は破損として扱い、token 認証、一覧、作成、失効をすべて `500 {"error":"Internal server error"}` で拒否する。ただし旧 `scope` 文字列から `scopes:[scope]` への正規化は §22.0c の例外に従う。破損内容、hash、token 本体は response、`.access_log`、`.audit_log`、journal に出力しない。

**API token 認証時の副作用境界：**

| ケース | HTTP status | `.api_tokens` | `.access_log` | `.audit_log` | 備考 |
|--------|-------------|---------------|---------------|--------------|------|
| hash 不一致 | `401` | 変更なし | 追記しない | 追記しない | token id が特定できないため監査対象外。 |
| 期限切れ | `401` | `last_used_at` 更新なし | `token_expired` を追記 | `token_expired` を追記 | token id が特定できた場合だけ記録する。 |
| 失効済み | `401` | `last_used_at` 更新なし | `token_revoked_reject` を追記 | `token_revoked_reject` を追記 | token 本体と hash は保存しない。 |
| scope 不足 | `403` | `last_used_at` 更新なし | `permission_denied` を追記 | `permission_denied` を追記 | 対象 endpoint は実行しない。 |
| 認証成功 | endpoint 固有 | `last_used_at` を現在時刻へ更新 | `token_auth` を追記 | `token_auth` を追記 | endpoint 実行前に更新する。 |

上表で `.api_tokens` 保存後に `.access_log` または `.audit_log` 追記へ失敗した場合、対象 API は `500` とし、endpoint 固有処理へ進まない。`last_used_at` 更新済み record は巻き戻さない。

**処理順：**

`POST /api/tokens` は以下の順で処理する。

1. 認証と `admin` scope を確認する。
2. `label`、`scopes`、`expires_at` を検証する。
3. `.api_tokens` のファイルロックを取得する。
4. 既存 record を読み込み、id を採番する。
5. token 本体を生成し、hash を算出する。
6. record を append して `.api_tokens` を原子的に保存する。
7. `.access_log` に token 作成成功を追記する。
8. `.audit_log` に `token_create` を追記する。
9. response に token 本体を 1 回だけ含めて返す。

`.api_tokens` 保存後に `.access_log` または `.audit_log` 追記へ失敗した場合、API は `500` を返す。作成済み token record は削除せず、再実行時は新しい token を発行する。

`DELETE /api/tokens/{id}` は以下の順で処理する。

1. 認証と `admin` scope を確認する。
2. path `id` を検証する。
3. `.api_tokens` をロックして対象 record を検索する。
4. 対象が存在しない、または `revoked_at != null` の場合は `404` を返す。
5. `revoked_at` を現在時刻に設定して保存する。
6. `.access_log` と `.audit_log` に失効成功を追記する。
7. `{ "message": "Token revoked" }` を返す。

認証に使用中の API token 自身を失効してよい。その場合、当該リクエストは成功し、次リクエストから `401` になる。

**token ID 採番・返却固定契約：**

| 項目 | 仕様 |
|------|------|
| 採番元 | `.api_tokens.tokens[].id` の数値 suffix 最大値。存在しない場合は `tok000001`。 |
| 衝突時 | 最大 suffix + 1 を採用する。削除済みや失効済み id は再利用しない。 |
| token 本体 | `act_` + `crypto/rand` 32 bytes を base64url padding なしで encode した文字列。 |
| 作成 response | `{ "id", "label", "scopes", "created_at", "expires_at", "revoked_at", "last_used_at", "token" }`。`token` はこの response だけに含める。 |
| 一覧 response | `tokens` 配列に `token_hash` と `token` を含めない。並び順は `created_at` 降順、同時刻は `id` 昇順。 |
| label 正規化 | 前後空白を除去し、内部空白は保持する。64 文字判定は Unicode code point 数ではなく UTF-8 byte 数で行う。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| scopes 空 | `422`。 |
| 未知 scope | `422`。 |
| label 空 | `422`。 |
| expires_at が現在以前 | `422`。 |
| token 上限 100 件超過 | `422`。失効済み record も件数に含める。 |
| 期限切れ token | `401`。 |
| 失効済み token 再失効 | `404` または冪等成功にせず `404` 固定。 |
| `.api_tokens` 破損 | `500`。自動再生成しない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 作成 | token 本体は response だけ。 |
| 一覧 | token hash と本体を返さない。 |
| 期限切れ | `401`、last_used_at 更新なし。 |
| 失効 | `revoked_at` 保存、以後 `401`。 |
| 自己失効 | 失効リクエストは `200`、以後同 token は `401`。 |
| log 追記失敗 | `500`。token record は残る。 |
| token record 破損 | `500`、自動再生成なし、secret 出力なし。 |
| scope 不足 | `403`、endpoint 実行なし、`permission_denied` 記録。 |
| 採番衝突 | 既存最大 suffix + 1 で作成される。 |
| 作成 response | token 本体は 1 回だけ含まれ、一覧では返らない。 |

### 27.44 監査ログ

本機能の目的は、認証、権限拒否、token、設定、build trigger などの重要操作を追跡できる JSON Lines 監査ログとして保存することである。

**対象 action：**

| action | target_type |
|--------|-------------|
| `login_success`, `login_failure`, `logout`, `totp_required`, `totp_failure`, `totp_enabled`, `totp_disabled` | `auth` |
| `password_change` | `auth` |
| `token_create`, `token_revoke`, `token_auth`, `token_expired`, `token_revoked_reject` | `api_token` |
| `permission_denied` | `endpoint` または `auth` |
| `build_trigger`, `build_force_trigger` | `build` |
| `config_update`, `rate_limit_update` | `config` |

**record 生成規則：**

| 項目 | 仕様 |
|------|------|
| `timestamp` | 操作結果が確定した UTC 時刻。 |
| `request_id` | リクエスト受付時に生成した 16 byte hex。同一 API 処理中に複数 log を書く場合は同じ値を使う。 |
| `actor_type` | 未認証 login は `"anonymous"`、管理 session は `"admin"`、API token は `"api_token"`、内部処理は `"system"`。 |
| `actor_id` | 管理 session は `"admin"`、API token は token id、未認証は `null`、内部処理は `"system"`。 |
| `target_id` | 対象 id がある場合は id。endpoint 拒否は `"{METHOD} {path}"`。対象なしは `null`。 |
| `result` | 成功は `"success"`、認証失敗や検証失敗は `"failure"`、権限拒否は `"denied"`。 |
| `message` | 固定文言のみ。入力値を連結しない。最大 500 文字。 |

監査ログへ保存する object は `.audit_log` schema のキーだけとする。未知キー、request body、query 全体、header 全体、cookie、secret、token、password、hash、salt、TOTP secret を保存してはならない。

**action 生成固定値：**

| 操作 | action | actor_type | target_type | target_id | result |
|------|--------|------------|-------------|-----------|--------|
| password login 成功 | `login_success` | `anonymous` | `auth` | `login` | `success` |
| password login 失敗 | `login_failure` | `anonymous` | `auth` | `login` | `failure` |
| login ロック拒否 | `permission_denied` | `anonymous` | `auth` | `login` | `denied` |
| TOTP 必須 ticket 発行 | `totp_required` | `anonymous` | `auth` | `login` | `success` |
| TOTP login 失敗 | `totp_failure` | `anonymous` | `auth` | `login_totp` | `failure` |
| logout | `logout` | `admin` または `api_token` | `auth` | `logout` | `success` |
| password 変更 | `password_change` | `admin` | `auth` | `password` | `success` |
| API token 作成 | `token_create` | `admin` または `api_token` | `api_token` | 作成 token id | `success` |
| API token 認証成功 | `token_auth` | `api_token` | `api_token` | token id | `success` |
| API token 期限切れ拒否 | `token_expired` | `api_token` | `api_token` | token id | `failure` |
| API token 失効済み拒否 | `token_revoked_reject` | `api_token` | `api_token` | token id | `failure` |
| API token 失効 | `token_revoke` | `admin` または `api_token` | `api_token` | 失効 token id | `success` |
| scope 不足 | `permission_denied` | `api_token` | `endpoint` | `{METHOD} {path}` | `denied` |
| rate limit 超過 | `permission_denied` | `anonymous`、`admin`、`api_token` | `endpoint` | `{METHOD} {path}` | `denied` |
| TOTP 有効化 | `totp_enabled` | `admin` | `auth` | `totp` | `success` |
| TOTP 無効化 | `totp_disabled` | `admin` | `auth` | `totp` | `success` |
| 設定変更 | `config_update` | `admin` または `api_token` | `config` | 変更 key | `success` |
| rate limit 設定変更 | `rate_limit_update` | `admin` または `api_token` | `config` | `api_rate_limit` | `success` |

同一操作で `.config_log` と `.audit_log` の両方を追記する場合、`.config_log` を先に追記する。`.config_log` 成功後に `.audit_log` が失敗した場合は `500` を返し、`.config_log` は巻き戻さない。`.audit_log` 追記失敗そのものを `.audit_log` に記録しようとしてはならない。

**監査 record 保存順・失敗契約：**

| 操作種別 | 保存順 | `.audit_log` 失敗時 |
|----------|--------|----------------------|
| 認証成功 | session または token 状態更新 → `.access_log` → `.audit_log` → response | token / session 状態は巻き戻さず `500`。response に token 本体を含めない。 |
| 認証失敗 | `.access_log` → `.audit_log` → response | `500`。失敗理由詳細は返さない。 |
| 権限拒否 | `.access_log` → `.audit_log` → `403` | `500`。対象 endpoint は実行しない。 |
| 設定変更 | 対象設定保存 → `.config_log` → `.audit_log` → response | 対象設定と `.config_log` は巻き戻さず `500`。 |
| token 作成 | `.api_tokens` 保存 → `.access_log` → `.audit_log` → response | 作成済み record は残し、token 本体は返さず `500`。 |

**正常系：**

1. 監査対象操作の成否が確定した後、`.audit_log` へ 1 行追記する。
2. 監査ログ追記に失敗した場合、対象操作は失敗扱いにし、`500` を返す。
3. secret、password、session token、API token 本体、hash、salt、TOTP secret は保存しない。
4. `GET /api/audit-log` は `limit`、`offset`、`actor`、`action`、`result` で絞り込み、新しい順で返す。

**取得仕様：**

| 項目 | 仕様 |
|------|------|
| 読み込み順 | ファイル先頭から全有効行を読み、フィルタ後に `timestamp` 降順、同時刻はファイル出現順の逆順で返す。 |
| `total` | フィルタ後、ページング前の有効 record 件数。 |
| `actor` filter | `actor_id` と完全一致。`anonymous` を指定した場合は `actor_type:"anonymous"` かつ `actor_id:null` に一致させる。 |
| `action` filter | `action` 完全一致。 |
| `result` filter | `success`、`failure`、`denied` の完全一致。 |
| 壊れた行 | 無視する。API response に壊れた行の内容を含めない。 |
| 返却上限 | `limit` は §22.0b の 1〜200。未指定時は 100。 |

`GET /api/audit-log` は監査ログ取得操作自体を `.audit_log` へ記録しない。通常の API request として `.api_access_log` には記録する。フィルタ値が未知 action、未知 result、200 文字超過 actor の場合は `422` とし、壊れた行の有無とは独立して判定する。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| password 変更 | `password_change` が記録される。 |
| 権限拒否 | `permission_denied` が記録される。 |
| secret 入力 | 監査ログに平文がない。 |
| 壊れた行 | API は無視して返す。 |
| filter | `actor`、`action`、`result` が完全一致で絞り込まれる。 |
| 追記失敗 | 対象操作は `500`。 |
| 取得操作 | `GET /api/audit-log` 自身は監査ログへ追記されない。 |
| 未知 filter | `422`。 |
| audit failure token create | token record は残るが token 本体は返らない。 |
| body secret | request body 全体が保存されない。 |

### 27.45 セッションタイムアウト変更設定

本機能の目的は、新規 session の有効期限を管理 API から変更可能にし、既存 session への影響を明確にすることである。

**仕様：**

| 項目 | 値 |
|------|----|
| 設定 key | `.server_config.session_timeout_seconds` |
| 既定値 | `28800` |
| 最小値 | `300` |
| 最大値 | `2592000` |
| 更新 API | `POST /api/config` または `setConfig({session_timeout_seconds})` |

設定変更は新規 session にだけ適用する。既存 session の `expires_at` は延長も短縮もしない。

**処理契約：**

| 項目 | 仕様 |
|------|------|
| 保存先 | `.server_config.session_timeout_seconds`。他 key と同時更新された場合も §22.0a の単一ファイル更新手順で保存する。 |
| 既定値 merge | `.server_config` に key がない場合、`GET /api/config` は `28800` を返す。ファイルへ暗黙保存しない。 |
| login 時適用 | session token 発行直前に `.server_config` を読み、当該時点の値で `expires_at` を計算する。 |
| TOTP login | TOTP 有効時は `POST /api/login/totp` の成功時点で値を読む。`POST /api/login` の password 成功時点では session を発行しない。 |
| 変更監査 | `POST /api/config` で値が変わった場合は `.config_log` に差分を記録する。`.audit_log` は `config_update`、`target_type:"config"`、`target_id:"session_timeout_seconds"` を記録する。 |
| 同値更新 | 同じ値の更新は `200` とし、`.server_config` の再保存は行ってよいが、差分なしとして `.config_log` と `.audit_log` には記録しない。 |

session timeout の値は session 発行時に秒単位で加算する。`expires_at = issued_at + session_timeout_seconds` とし、計算後の時刻は UTC ISO 8601 秒精度で保存する。ミリ秒、ナノ秒、local timezone は保存しない。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 300 秒設定 | 新規 session が 300 秒後に期限切れ。 |
| 既存 session | 設定変更後も元の `expires_at`。 |
| 範囲外 | `422`。 |
| key 不在 | `GET /api/config` は `28800`。 |
| TOTP login | `POST /api/login/totp` 成功時点の値で session 期限を決める。 |
| 秒精度 | `expires_at` は UTC ISO 8601 秒精度。 |

### 27.46 TOTP 二要素認証

本機能の目的は、外部ライブラリなしで RFC 6238 TOTP を検証し、password 漏えい時の管理 API 不正利用を抑止することである。

**アルゴリズム：**

| 項目 | 仕様 |
|------|------|
| secret | `crypto/rand` 20 bytes を RFC 4648 base32 padding なしで表示する。 |
| HMAC | `crypto/hmac` + `crypto/sha1`。 |
| step | 30 秒。 |
| digits | 6 桁。 |
| window | 現在 step の前後 1 step を許可する。 |
| replay 防止 | `.totp_secret.last_accepted_step` 以下の step は拒否する。 |
| otpauth URI | `otpauth://totp/Adlaire%20CI:admin?secret={secret}&issuer=Adlaire%20CI&algorithm=SHA1&digits=6&period=30`。 |

QR code 生成は初期実装対象外とする。UI は secret と otpauth URI を一回表示し、ユーザーが認証アプリへ手入力またはURI貼り付けできるようにする。

**メモリ上状態：**

| 状態 | 保存場所 | 期限 | 内容 |
|------|----------|------|------|
| setup 仮 secret | `components/api.go` のメモリ | 10 分 | `secret_base32`, `created_at`。サーバー再起動で破棄する。 |
| login ticket | `components/api.go` のメモリ | 5 分 | `ticket_hash`, `created_at`, `password_verified_at`。ticket 本体は hash 化して保持する。 |

setup 仮 secret と login ticket は永続ファイルへ保存しない。API response、UI 一回表示、メモリ上状態以外に secret/ticket 本体を残してはならない。

**正常系：**

1. `POST /api/auth/totp-setup` は仮 secret と otpauth URI を返すが、永続化しない。
2. `POST /api/auth/totp-confirm` は仮 secret と code を検証し、成功時だけ `.totp_secret.enabled=true` と secret 情報を保存する。
3. TOTP 有効時の `POST /api/login` は password 成功後に ticket を返し、session token は返さない。
4. `POST /api/login/totp` は ticket と code を検証し、成功時に session token を返す。
5. `DELETE /api/auth/totp` は code を検証して TOTP を無効化する。

**endpoint 処理詳細：**

| endpoint | 処理 |
|----------|------|
| `GET /api/auth/totp-status` | `.totp_secret` を読み、`{enabled,confirmed_at}` だけを返す。`secret_base32` と `last_accepted_step` は返さない。 |
| `POST /api/auth/totp-setup` | TOTP 有効時は `409`。無効時は仮 secret を生成し、既存の未確認仮 secret を上書きする。`.totp_secret` は書き込まない。 |
| `POST /api/auth/totp-confirm` | 仮 secret がない、または期限切れなら `409`。code 成功時に `.totp_secret` を保存し、仮 secret をメモリから削除する。 |
| `DELETE /api/auth/totp` | `.totp_secret.enabled == false` は `409`。code 成功時に `enabled:false`, `secret_base32:null`, `confirmed_at:null`, `last_accepted_step:null` を保存する。 |
| `POST /api/login` | password 成功かつ TOTP 有効なら `ticket` を `crypto/rand` 32 bytes の lowercase hex で生成し、`{must_change,totp_required:true,ticket}` を返す。 |
| `POST /api/login/totp` | ticket hash と code を検証し、成功時に ticket を削除して session token を返す。失敗時も ticket は削除する。 |

TOTP code は 6 桁の ASCII 数字のみ受け付ける。空文字、全角数字、空白付き文字列、6 桁以外は `422` とする。検証は `window` 内の step を古い順に試し、最初に一致した step を採用する。採用 step が `.totp_secret.last_accepted_step` 以下の場合は `401` とする。

TOTP 関連の成功、失敗、無効化、ticket 発行は `.audit_log` へ記録する。code 不一致、replay、ticket 不正は `result:"failure"` とし、code、secret、ticket 本体は保存しない。

**TOTP 計算固定契約：**

| 項目 | 仕様 |
|------|------|
| counter | `floor(unix_seconds / 30)` を 8 byte big-endian unsigned integer として HMAC 入力にする。 |
| truncation | RFC 4226 dynamic truncation を使い、31 bit integer を `10^6` で剰余する。 |
| 表示 | 6 桁未満は左ゼロ埋めする。 |
| base32 decode | 大文字 ASCII のみ保存する。入力確認時は空白を除去せず、保存値と同じ RFC 4648 padding なし形式だけを扱う。 |
| clock source | API server の現在時刻だけを使う。client 時刻は受け取らない。 |

**TOTP 状態更新順：**

| 操作 | 更新順 | 失敗時 |
|------|--------|--------|
| setup 開始 | 仮 secret 生成 → メモリ保存 → response | メモリ保存失敗時は `500`、secret を返さない。 |
| confirm 成功 | `.totp_secret` 保存 → 仮 secret 削除 → `.audit_log` 追記 → response | `.audit_log` 失敗時は `500`。保存済み `.totp_secret` と仮 secret 削除は巻き戻さない。 |
| login password 成功 / TOTP 有効 | ticket 生成 → メモリ保存 → `.access_log` 追記 → `.audit_log` 追記 → response | log 失敗時は ticket を削除し `500`。 |
| login TOTP 成功 | code step 採用 → `.totp_secret.last_accepted_step` 保存 → ticket 削除 → `.admin_credentials` 更新 → `.access_log` 追記 → `.audit_log` 追記 → session 追加 → response | session 追加前の失敗は token を返さない。ticket 削除後は同 ticket を再利用不可。 |
| login TOTP 失敗 | ticket 削除 → `.access_log` 追記 → `.audit_log` 追記 → response | log 失敗時は `500`。ticket は巻き戻さない。 |
| disable 成功 | `.totp_secret` を無効値で保存 → 未使用 ticket / 仮 secret 全削除 → `.audit_log` 追記 → response | `.audit_log` 失敗時は `500`。保存済み無効化は巻き戻さない。 |

`POST /api/login/totp` は、ticket が存在しない、期限切れ、hash 不一致、既に削除済みのいずれの場合も `401 {"error":"Unauthorized"}` を返す。ticket 不正の詳細、ticket hash、ticket 発行時刻は response と log に含めない。TOTP 有効化または無効化に成功した場合、既存 session は破棄しない。

**異常系：**

| 条件 | 処理 |
|------|------|
| code 不一致 | `401`。 |
| code 形式不正 | `422`。 |
| ticket 期限切れ | `401`。ticket 有効期限は 5 分。 |
| setup 未実行 confirm | `409`。 |
| TOTP 無効状態の disable | `409`。 |
| TOTP 有効状態の setup | `409`。 |
| `.totp_secret` 破損 | TOTP 有効 login、status、confirm、disable は `500`。自動無効化しない。 |
| `.totp_secret` 書き込み失敗 | `500`。response に secret、token、ticket を含めない。 |
| `.audit_log` 追記失敗 | `500`。保存済み `.totp_secret` は巻き戻さない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 正常 setup | secret は確認まで保存されない。 |
| 正常 login | password 後 ticket、TOTP 後 token。 |
| replay | 同一 step の再利用は `401`。 |
| secret 表示 | response と UI の一回表示以外に平文が残らない。 |
| setup 期限切れ | confirm は `409`。 |
| login ticket 再利用 | 1 回成功後または失敗後の同一 ticket は `401`。 |
| confirm audit 失敗 | `500`、`.totp_secret` は保存済み、secret 平文は log なし。 |
| disable 成功 | `.totp_secret` は無効値、既存 session は維持、ticket と仮 secret は削除。 |
| window 前後 | 現在 step の前後 1 step が成功し、同 step 再利用は失敗する。 |
| 全角 code | `422`。 |

### 27.47 API レート制限

本機能の目的は、login 総当たり、API token 濫用、外部連携の暴走を Go 標準ライブラリだけで抑止することである。

**固定窓 policy：**

| group | 既定 window | 既定 max | 対象 |
|-------|-------------|----------|------|
| `login` | 60 秒 | 10 | `POST /api/login`、`POST /api/login/totp`。 |
| `read` | 60 秒 | 600 | 読み取り endpoint。 |
| `trigger` | 60 秒 | 60 | build trigger endpoint。 |
| `operate` | 60 秒 | 120 | 運用操作 endpoint。 |
| `config` | 60 秒 | 60 | 設定変更 endpoint。 |
| `admin` | 60 秒 | 60 | token、audit、rate limit endpoint。 |

**判定キー：**

| 認証状態 | key |
|----------|-----|
| 認証前 | `ip:{remote_addr}:{group}` |
| session | `session:admin:{group}` と `ip:{remote_addr}:{group}` の両方 |
| API token | `token:{token_id}:{group}` と `ip:{remote_addr}:{group}` の両方 |

どちらか一方でも上限を超えた場合は `429 Too Many Requests` と `{"error":"Too many requests"}` を返す。

`remote_addr` は `net/http.Request.RemoteAddr` の host 部分を使用する。`X-Forwarded-For`、`X-Real-IP`、`Forwarded` header は標準では信用せず、key 生成に使用しない。IPv6 は `net.SplitHostPort` で host を抽出し、正規化済み文字列をそのまま key に入れる。host 抽出に失敗した場合は `ip:unknown:{group}` を使用する。

**policy object：**

```json
{
  "enabled": true,
  "groups": {
    "login": { "window_seconds": 60, "max_requests": 10 },
    "read": { "window_seconds": 60, "max_requests": 600 },
    "trigger": { "window_seconds": 60, "max_requests": 60 },
    "operate": { "window_seconds": 60, "max_requests": 120 },
    "config": { "window_seconds": 60, "max_requests": 60 },
    "admin": { "window_seconds": 60, "max_requests": 60 }
  },
  "state_summary": []
}
```

永続化先は `.server_config.api_rate_limit` とし、保存時は `enabled` と `groups` だけを保存する。`state_summary` は `GET /api/api-rate-limit` と `POST /api/api-rate-limit` response 用の算出値であり、永続化しない。

**`state_summary` item：**

| キー | 型 | 説明 |
|------|----|------|
| `key` | string | `.api_rate_state.windows` の key。 |
| `group` | string | endpoint group。 |
| `window_start` | string | window 開始時刻。 |
| `count` | integer | 現在 count。 |
| `reset_at` | string | `window_start + window_seconds`。 |

`state_summary` は `reset_at` 降順、同時刻は `key` 昇順で最大 100 件返す。期限切れ window は response 算出前に `.api_rate_state` から削除してよい。

**正常系：**

1. path 解決後、認証前に IP key の login 制限を確認する。
2. 認証後、endpoint group を判定し、actor key と IP key の count を更新する。
3. `GET /api/api-rate-limit` は policy と window summary を返す。
4. `POST /api/api-rate-limit` は policy を検証して保存し、`.api_rate_state.windows` を空にする。

**endpoint group 判定：**

| 条件 | group |
|------|-------|
| `POST /api/login`、`POST /api/login/totp` | `login` |
| §27.42 の `trigger` 許可 endpoint | `trigger` |
| §27.42 の `operate` 許可 endpoint | `operate` |
| §27.42 の `config` 許可 endpoint | `config` |
| §27.42 の `admin` 許可 endpoint | `admin` |
| §27.42 の `read` 許可 endpoint | `read` |
| `GET /api/health` | rate limit 対象外 |
| 未知 path / method 不一致 | rate limit 判定前に `404` / `405` |

endpoint が複数 group に現れる場合は、`login`、`admin`、`config`、`operate`、`trigger`、`read` の順で最初に一致した group を採用する。

**判定・更新手順：**

1. `.server_config.api_rate_limit.enabled == false` の場合、`.api_rate_state` を読まずに対象 API 処理へ進む。
2. endpoint group を決める。
3. `.api_rate_state` のファイルロックを取得する。
4. 対象 key ごとに `.api_rate_state.windows[key]` を確認する。
5. window が存在しない、または `now >= window_start + window_seconds` の場合、`window_start=now`, `count=0` で初期化する。
6. `count >= max_requests` の key が 1 つでもあれば、count を増やさず `.audit_log` に `permission_denied` を追記し、`429` を返す。
7. 上限未満の場合、対象 key すべての `count` を 1 増やして `.api_rate_state` を保存し、対象 API 処理へ進む。

rate limit の `429` は `.audit_log` に `permission_denied` として記録する。監査ログ追記に失敗した場合は `500` を返す。login group の認証前 `429` は `actor_type:"anonymous"`、`actor_id:null` とする。

認証後 endpoint の rate limit では、actor key と IP key の両方を同じ lock 内で判定・更新する。片方だけの count 更新に成功した状態を残してはならない。`.api_rate_state` 保存失敗時は対象 API を実行せず `500` を返す。rate limit 判定で `429` になる request は count を増やさない。

**rate limit 副作用固定契約：**

| ケース | `.api_rate_state` | `.access_log` | `.audit_log` | endpoint 固有処理 |
|--------|-------------------|---------------|--------------|-------------------|
| 上限未満 | count を増やす | response 確定後に通常追記 | endpoint が監査対象の場合だけ追記 | 実行する |
| 上限超過 | count を増やさない | `429` として追記 | `permission_denied` を追記 | 実行しない |
| `.api_rate_state` 保存失敗 | 部分更新を残さない | `500` として追記を試行 | 追記しない | 実行しない |
| `.audit_log` 失敗 | count を増やさない | `500` として追記を試行 | 失敗 | 実行しない |

`429` 判定時は `.audit_log` 追記を `.api_rate_state` 保存前に行う。`.audit_log` 追記に成功した場合だけ `429` を返す。`.audit_log` 追記に失敗した場合は `.api_rate_state` を変更せず `500` を返す。

**設定更新契約：**

| 項目 | 仕様 |
|------|------|
| 必須 group | `login`、`read`、`trigger`、`operate`、`config`、`admin` をすべて含める。 |
| 余分な group | `422`。 |
| `enabled` | boolean 必須。 |
| `window_seconds` | integer 必須、1〜86400。 |
| `max_requests` | integer 必須、1〜100000。 |
| 保存順 | `.server_config` 保存 → `.api_rate_state.windows` 空保存 → `.config_log` 追記 → `.audit_log` に `rate_limit_update` 追記 → response。 |
| 同値更新 | `200` とし、`.server_config` と `.api_rate_state` は変更しない。`.config_log` と `.audit_log` に追記しない。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| policy group 不足 | `422`。 |
| window_seconds 範囲外 | `422`。許容値は 1〜86400。 |
| max_requests 範囲外 | `422`。許容値は 1〜100000。 |
| `.api_rate_state` 書き込み失敗 | `500`。対象 API は実行しない。 |
| unknown group | `422`。 |
| `enabled` 欠落 | `422`。 |
| `.api_rate_state` 破損 | §22.0a に従って退避し、空 window で再生成する。 |
| `.audit_log` 追記失敗 | `500`。`429` response は返さない。 |
| lock 取得 10 秒超過 | `409`。対象 API は実行しない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| login 11 回目 | `429`。 |
| window 経過 | count reset、成功。 |
| session と IP | どちらか超過で `429`。 |
| disabled | `enabled:false` の場合は判定せず成功。 |
| policy 更新 | `.api_rate_state.windows` が空になる。 |
| state_summary | 最大 100 件、`reset_at` 降順で返る。 |
| 429 count | `429` になった request では count が増えない。 |
| 同値更新 | state と log を変更せず `200`。 |
| proxy header | `X-Forwarded-For` ではなく `RemoteAddr` host で key を作る。 |
| audit failure on 429 | count は増えず `500`。 |
| state save failure | endpoint 固有処理なし、部分 count 更新なし。 |
