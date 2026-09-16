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
4. 対象コンポーネントの詳細節を読み、入力、出力、状態、正常系、異常系、セキュリティ、検証条件を確認する。
5. §26 のセットアップ・アップデート手順と §26.7 の受け入れ条件に影響がある場合は、実装 PR の検証対象に含める。

詳細仕様節に §0h の必須項目が不足している場合は、実装判断で補完してはならない。先に本ファイルを改訂し、`ADLAIRE_CI_SPEC.md` の対象範囲と整合させる。

| 範囲 | 役割 |
|------|------|
| §0〜§0j | 詳細仕様の記載基準、実装前確認項目、共通固定値、検証、Phase、詳細節対応表、リポジトリ内ソース配置 |
| §1〜§9 | `components/builder.go` / `adlaire-ci-build` の詳細仕様 |
| §10〜§20 | `components/runner.go` / `adlaire-ci-runner` の詳細仕様 |
| §21〜§22 | `components/api.go` / `adlaire-ci-api` の詳細仕様 |
| §23 | `admin/adlaire-ci-sdk.js` の詳細仕様 |
| §24 | `admin/index.html` の詳細仕様 |
| §25 | 認証の実装仕様 |
| §26 | バイナリ配布前提のセットアップ、アップデート、受け入れ条件 |

## 0a. 詳細仕様の記載基準

本ファイルの仕様項目は、実装者が追加の設計判断や推測を行わずに実装できる粒度で記載する。

仕様項目を追加または改訂する場合は、対象範囲に応じて以下を明記する。

| 項目 | 記載する内容 |
|------|-------------|
| 対象 | 対象コンポーネント、対象ファイル、対象機能、責務境界 |
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

本節は、対象コンポーネントごとに参照する詳細仕様節を示す。実装状態、実装可否、ロードマップ状態は `ADLAIRE_CI_SPEC.md` を確認する。

| 対象コンポーネント | 詳細仕様節 | 主な確認対象 |
|--------------------|------------|--------------|
| `components/builder.go` | §1〜§9、§8a | CLI、入力 Markdown、出力サイト、HTML / CSS / JavaScript、変換 report、fixture。 |
| `components/runner.go` | §10〜§20、§15a、§26 | 設定、状態ファイル、GitHub API、pipeline、転送、snapshot、通知、systemd、fixture。 |
| `components/api.go` | §21〜§22、§25、§26 | API 共通処理、endpoint、状態ファイル、認証、session、systemd。 |
| `admin/adlaire-ci-sdk.js` | §23 | SDK class、method、HTTP 対応、error、stream、token 破棄。 |
| `admin/index.html` | §24 | 画面構成、DOM id、panel、SDK 呼び出し、表示状態、秘密情報消去。 |
| `components/mcp.go` | 詳細仕様なし | 本ファイルでは実装可能な入出力、状態、起動手順、検証条件を定義しない。 |

---

## 0c. 実装前確認項目

実装者は、対象機能について以下の条件をすべて満たすまで実装を開始してはならない。

| ゲート | 合格条件 |
|--------|----------|
| 対応表 | 対象機能が §0i の詳細節対応表に記載され、詳細仕様節と受け入れ条件が一意に示されている。 |
| テンプレート | 対象機能の詳細仕様が §0h の機能仕様テンプレートの必須項目を満たしている。 |
| 責務境界 | 対象コンポーネント、対象ファイル、呼び出し元、呼び出し先、変更してよい状態ファイルが明記されている。 |
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

対象項目を完了扱いにする場合は、対象コンポーネントごとに下表の検証を満たす。実装ファイルが存在しても、本表の必須検証が未完了の場合は完了扱いにしない。

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

本節は、Go 版初期実装へ進む前に仕様策定が完了しているかを判定するチェックである。実装者は、対象コンポーネントごとに下表の必須条件を満たすまで実装を開始してはならない。

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
| 対象コンポーネント | `components/builder.go`、`components/runner.go`、`components/api.go`、`admin/adlaire-ci-sdk.js`、`admin/index.html` のいずれが責務を持つか。複数の場合は責務境界を分けて書く。 | 実装不可。 |
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

表の「詳細仕様節」が複数ある場合は、すべての節を同時に満たす。該当節に §0h の必須項目が不足している場合は、その項目を実装せず、先に詳細仕様を改訂する。

| 機能 | 対象コンポーネント | 詳細仕様節 | 受け入れ条件 |
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
| `main.go` | 起動入口。サブコマンド判定、引数受け取り、対象コンポーネント呼び出しを行う。 |
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

## 1. 要件

| 項目 | 内容 |
|------|------|
| Go バージョン | Go `1.22` 以上。 |
| 外部依存 | なし。Go 標準ライブラリのみを使用する。外部依存が必要になった場合は実装せず、先に `ADLAIRE_CI_SPEC.md` Part 1 §4.1 と Part 2 §4 に従って仕様改訂する。 |
| 入力 | UTF-8 エンコードの Markdown ファイル、または Markdown ファイルを含むディレクトリ |
| 出力 | 静的 Web サイトディレクトリ（HTML / CSS / JavaScript / search index） |

---

## 2. ファイルパス設定

Go 版 `components/builder.go` は、以下の既定値を持つ設定構造体で入出力パスを管理する。

```go
type BuildConfig struct {
    Src     string
    Out     string
    Title   string
    Theme   string
    BaseDir string
    Strict  bool
}

var DefaultBuildConfig = BuildConfig{
    Src: "/opt/adlaire-builder/repo/docs",
    Out: "/opt/adlaire-builder/dist/site",
    Title: "Adlaire Documentation",
    Theme: "adlaire-default",
    BaseDir: "",
    Strict: false,
}
```

別の環境で実行する場合は、この既定値を CLI 引数で上書きする。Go 版 `components/builder.go` は設定ファイルを読み込まない。

**CLI 引数仕様：**

| 引数 | 必須 | 既定値 | 説明 |
|------|------|--------|------|
| `--src <path>` | 任意 | `DefaultBuildConfig.Src` | 入力 Markdown ファイルまたは Markdown ディレクトリの絶対パスまたは相対パス。相対パスはカレントディレクトリ基準で解決する。 |
| `--out <path>` | 任意 | `DefaultBuildConfig.Out` | 出力サイトディレクトリの絶対パスまたは相対パス。存在しない場合は作成する。 |
| `--title <text>` | 任意 | `DefaultBuildConfig.Title` | サイト名、`index.html` の `<title>`、header 表示名に使用する。空文字は禁止。 |
| `--theme <name>` | 任意 | `DefaultBuildConfig.Theme` | 初期仕様では `adlaire-default` のみ許可する。 |
| `--base-dir <path>` | 任意 | `DefaultBuildConfig.BaseDir` | 相対リンク・画像解決の基準ディレクトリ。空の場合は `--src` がファイルなら親ディレクトリ、ディレクトリなら `--src` 自身を使用する。 |
| `--strict` | 任意 | `false` | 警告をビルド失敗として扱う。警告が 1 件以上ある場合は終了コード `2` とする。 |
| `--build-id <id>` | 任意 | 空文字 | 出力 HTML の `<head>` に `adlaire-build-id` として埋め込む。空文字の場合も空 content の meta を出力する。 |
| `--commit-sha <sha>` | 任意 | 空文字 | 出力 HTML の `<head>` に `adlaire-commit-sha` として埋め込む。空文字の場合も空 content の meta を出力する。 |
| `--build-at <iso8601>` | 任意 | 空文字 | 出力 HTML の `<head>` に `adlaire-build-at` として埋め込む。値がある場合は UTC ISO 8601 のみ許可する。 |
| `--version` | 任意 | なし | バイナリ名、仕様名、Go build 情報を 1 行で標準出力へ表示して終了する。 |
| `--help` | 任意 | なし | 引数一覧を標準出力へ表示して終了する。 |

**CLI 引数の異常系：**

| 条件 | 終了コード | 出力 |
|------|------------|------|
| 未知の引数 | `2` | stderr に `unknown option: <name>` |
| `--src` / `--out` / `--title` / `--theme` / `--base-dir` / `--build-id` / `--commit-sha` / `--build-at` の値欠落 | `2` | stderr に `missing value: <name>` |
| `--src` が存在しない | `2` | stderr に `source not found: <path>` |
| `--src` が Markdown ファイルでもディレクトリでもない | `2` | stderr に `source is not markdown file or directory: <path>` |
| `--src` 内の Markdown が UTF-8 として読めない | `2` | stderr に `source is not valid UTF-8: <path>` |
| `--title` が空文字 | `2` | stderr に `title must not be empty` |
| `--theme` が `adlaire-default` 以外 | `2` | stderr に `unknown theme: <name>` |
| `--build-id` が空文字以外で `b{YYYYMMDDHHmmss}` または `b{YYYYMMDDHHmmss}-NNN` 形式でない | `2` | stderr に `invalid build id: <value>` |
| `--commit-sha` が空文字以外で 7〜40 文字の lowercase hex でない | `2` | stderr に `invalid commit sha: <value>` |
| `--build-at` が空文字以外で UTC ISO 8601 でない | `2` | stderr に `invalid build at: <value>` |
| `--out` ディレクトリ作成失敗 | `1` | stderr に `cannot create output directory: <path>` |
| `--out` が既存ファイル | `1` | stderr に `output path is not directory: <path>` |
| `--out` 書き込み失敗 | `1` | stderr に `cannot write output: <path>` |

`--help` と `--version` は他の引数より優先し、成功時は終了コード `0` とする。

**CLI パース固定仕様：**

- 引数は `flag` package 互換の `--name value` と `--name=value` の両方を許可する。
- 短縮オプション（例：`-s`、`-o`）は禁止する。指定された場合は未知の引数として扱う。
- 同一引数が複数回指定された場合は最後の値を採用する。ただし `--strict` は 1 回以上指定されれば `true` とする。
- `--src`、`--out`、`--base-dir` の相対パスは `os.Getwd()` の戻り値を基準に `filepath.Abs()` で絶対パスへ変換する。
- `--base-dir` が空ではない場合、存在するディレクトリでなければならない。存在しない場合は終了コード `2`、stderr に `base directory not found: <path>` を出力する。
- `--base-dir` がファイルの場合は終了コード `2`、stderr に `base path is not directory: <path>` を出力する。
- stderr のエラー行は末尾に改行 1 つを付ける。複数エラーをまとめて出力せず、最初に検出したエラー 1 件で終了する。

**固定出力：**

| 条件 | stdout |
|------|--------|
| `--help` | `Usage: adlaire-ci-build [--src path] [--out path] [--title text] [--theme name] [--base-dir path] [--strict] [--build-id id] [--commit-sha sha] [--build-at iso8601] [--version] [--help]` |
| `--version` | `adlaire-ci-build ADLAIRE_CI_SPEC go=<runtime.Version()>` |

`--help` と `--version` の stdout は 1 行固定とし、末尾に改行 1 つを付ける。`--help` または `--version` を指定した場合、`--src` の存在確認、`--theme` 検証、出力ディレクトリ作成は行わない。

---

## 2a. 入力収集・出力パス決定

`components/builder.go` は、`--src` がファイルかディレクトリかで入力収集方法を切り替える。

| `--src` 種別 | 入力収集 | 出力 |
|--------------|----------|------|
| Markdown ファイル | 指定ファイル 1 件のみを入力とする。拡張子は `.md` または `.markdown` のみ許可する。 | `index.html` 1 件と `assets/` を出力する。 |
| ディレクトリ | 配下の `.md` / `.markdown` ファイルを再帰収集する。隠しディレクトリ、`.git`、`.ci`、`node_modules`、`vendor`、`dist` は収集対象外とする。 | `index.html` をサイト目次、各 Markdown を `pages/{slug}.html` として出力する。 |

入力ファイルの並び順は、`--base-dir` からの相対パスを `/` 区切りに正規化した文字列の昇順とする。OS やファイルシステムの列挙順に依存してはならない。

Markdown ディレクトリ入力で Markdown ファイルが 0 件の場合は終了コード `2` とし、stderr に `no markdown files found: <path>` を出力する。

`PageData.Title` は、各 Markdown ファイルの最初の h1 見出しを使用する。h1 が存在しない場合は、拡張子を除いたファイル名をタイトル化して使用する。ファイル名タイトル化では、`-` と `_` を空白に置換し、前後空白を除去する。空になった場合は `Untitled` とする。

ページ slug は以下の順で決定する。

1. `--base-dir` からの相対パスから拡張子を除く。
2. パス区切り `/`、空白、`_` を `-` に置換する。
3. §4.1 `slugify` と同じ文字種ルールで正規化する。
4. 空になった場合は `page` とする。
5. 同一 slug が重複した場合は、2 件目以降に `-2`、`-3` のように連番を付与する。

出力パスは以下とする。

| 条件 | `PageData.OutputPath` |
|------|-----------------------|
| 単一 Markdown 入力 | `index.html` |
| ディレクトリ入力のサイト目次 | `index.html` |
| ディレクトリ入力の各 Markdown | `pages/{pageSlug}.html` |

`assets/style.css`、`assets/app.js`、`assets/search-index.json` は常に出力する。`assets/` 配下へ Markdown 由来ファイルを出力してはならない。

**リンク・画像・相対パス解決：**

| 入力 | 処理 |
|------|------|
| `http://` / `https://` URL | 外部 URL としてそのまま出力する。 |
| `mailto:` / `tel:` URL | 内部リンク検証対象外とし、`href` はそのまま出力する。 |
| `#anchor` | 同一ページ内 anchor として、ページ内の一意化後 slug と照合する。 |
| `./doc.md` / `../dir/doc.md` / `dir/doc.markdown#x` | `--base-dir` 基準で Markdown 入力ファイルへ解決し、該当ページの HTML パスへ変換する。anchor がある場合は `#x` を維持し、出力先ページの slug と照合する。 |
| `.md` / `.markdown` 以外の相対リンク | ファイルをコピーせず、元の相対 URL を保持する。存在確認は警告対象外。 |
| 画像 `![alt](path)` | ファイルコピーを行わず、`src` は元 URL を `esc()` して出力する。`assets/` へ画像を複製してはならない。 |

Markdown 間リンクの解決に失敗した場合、HTML は元 URL のまま出力し、`[WARN] BROKEN_PAGE_LINK: <url> (in: <source>)` を出力する。`--strict` が `true` の場合、警告出力後に終了コード `2` とする。ページ間リンク解決で使用するパス比較は、絶対パス化、`filepath.Clean()`、パス区切り `/` 正規化を行った文字列で比較する。

**入力収集・出力生成の機能単位契約：**

`components/builder.go` は、以下の機能単位を順番に実行する。各機能単位は前段の出力だけを入力とし、失敗時は後続機能を実行しない。

| 機能単位 | 入力 | 出力 | 失敗条件 | 失敗時の状態 |
|----------|------|------|----------|--------------|
| CLI resolver | `os.Args`、`DefaultBuildConfig` | `BuildConfig` | 未知引数、値欠落、空 title、未知 theme | ファイル読込・出力作成を行わず終了コード `2`。 |
| source collector | `BuildConfig.Src`、`BuildConfig.BaseDir` | `[]PageInput` | source 不在、Markdown 0 件、UTF-8 不正、許可外拡張子 | 出力ディレクトリを変更せず終了コード `2`。 |
| page planner | `[]PageInput` | `[]PageData`、slug map、output path map | slug / output path 衝突を連番解決後も一意化できない | 出力ディレクトリを変更せず終了コード `1`。 |
| markdown renderer | `PageInput`、heading map、footnote map | page HTML fragment、warning list、report counters | 未閉鎖 fence、broken link などの警告があり `--strict=true` | 出力ディレクトリを変更せず終了コード `2`。 |
| site assembler | `[]PageData`、theme component | `index.html`、`pages/*.html`、`assets/*` のメモリ上生成物 | theme component 欠落、template 合成失敗 | 出力ディレクトリを変更せず終了コード `1`。 |
| atomic output writer | メモリ上生成物、`BuildConfig.Out` | 出力サイトディレクトリ | tmp 作成失敗、書き込み失敗、rename 失敗、sync 失敗 | 既存正常出力を維持し終了コード `1`。 |
| report emitter | report counters、warning list | stdout `[WARN]`、`[REPORT]` | stdout 書き込み失敗 | 終了コード `1`。出力済みファイルは巻き戻さない。 |

`PageInput` は `{source_path, relative_path, raw_text}` を持つメモリ上構造とする。`PageData` は `{title, slug, source_path, output_path, html, headings, warnings, reading_time}` を持つメモリ上構造とする。これらの構造は状態ファイルとして保存しない。

**atomic output writer 詳細：**

1. `BuildConfig.Out` と同じ親ディレクトリに `{basename}.tmp.{pid}` を作成する。
2. tmp 内に `index.html`、必要な `pages/`、`assets/style.css`、`assets/app.js`、`assets/search-index.json` をすべて書き込む。
3. 全ファイルを close し、通常ファイルは `0644`、ディレクトリは `0755` に補正する。
4. tmp 配下の各ファイルを `Sync` し、tmp ディレクトリを `Sync` する。
5. 既存 `BuildConfig.Out` が存在する場合は `{basename}.prev.{pid}` へ rename する。
6. tmp を `BuildConfig.Out` へ rename する。
7. 親ディレクトリを `Sync` する。
8. 手順 6 まで成功した後に prev を削除する。prev 削除失敗は `[WARN] OUTPUT_PREV_CLEANUP_FAILED: path={path}` とし、終了コードは `0` のままとする。

手順 1〜4 で失敗した場合は tmp を削除し、既存 `BuildConfig.Out` を変更しない。手順 5 成功後から手順 6 失敗までの間に失敗した場合は、prev を `BuildConfig.Out` へ戻す。復元に失敗した場合は stderr に `cannot restore previous output: <path>` を出力し終了コード `1` とする。

---

## 3. 処理パイプライン

```
┌────────────────────────────────────────────────────────────────┐
│ 1. MDファイル読み込み（raw_lines）                              │
├────────────────────────────────────────────────────────────────┤
│ 2. 見出し抽出パス（フェンス内を除外）                           │
│    → headings: []Heading（slug は重複解決済み）                 │
│    → slugByLine: map[int]string                                │
├────────────────────────────────────────────────────────────────┤
│ 2b. 脚注定義収集パス                                            │
│    → ctx.FootnoteDefs: map[string]string                       │
│    ※ ctx.FootnoteOrder は convert() 実行中に inline() が更新    │
├────────────────────────────────────────────────────────────────┤
│ 3. TOC HTML 生成（buildTOC(headings)）                          │
│    → tocHTML: string                                           │
├────────────────────────────────────────────────────────────────┤
│ 4. MD → HTML 変換（convert(lines, headings, ctx)）              │
│    → ConvertResult.HTML                                        │
├────────────────────────────────────────────────────────────────┤
│ 5. サイトテンプレート合成                                      │
│    theme component、CSS、JS、search index を生成                 │
├────────────────────────────────────────────────────────────────┤
│ 6. サイトディレクトリ書き出し（OUT）                            │
└────────────────────────────────────────────────────────────────┘
```

---

## 4. 関数リファレンス

### 4.1 `slugify(text string) string`

見出しテキストから HTML アンカー用のスラグ（`id` 属性値）を生成する。

**責務：** `slugify` は入力テキストから重複解決前の base slug だけを返す。重複解決は §4.5 の見出し収集処理で `slugCount map[string]int` を使って行う。`slugify` 内で状態を保持してはならない。

**処理手順：**
1. Markdown 記法文字（`` ` * _ ~ [ ] ``）を除去
2. 文字を 1 文字ずつ走査し、空白・ハイフン・ドット・アンダースコアは `-` に変換、英数字・Unicode 文字（カテゴリ `L*`、`N*`）はそのまま保持、それ以外は除去
3. 連続する `-` を 1 つに正規化、前後の `-` をトリム
4. 空文字列になった場合は `section` にフォールバック

**入力：**

| 引数 | 型 | 条件 |
|------|----|------|
| `text` | `string` | UTF-8 文字列。空文字を許可する。 |

**出力：**

| 戻り値 | 型 | 条件 |
|--------|----|------|
| `baseSlug` | `string` | 空文字は禁止。入力が空または除去後に空になる場合は `section`。 |

**エラー：** 戻り値エラーは持たない。入力が UTF-8 不正になる可能性は §2 の入力ファイル読み込み段階で排除する。

**禁止事項：**
- package 変数、グローバル map、呼び出し間で残る状態を使用しない。
- 重複時の `-2` 付与を `slugify` 内で行わない。
- URL encode は行わない。`id` と `href` には本関数の戻り値または一意化後 slug をそのまま使用する。

---

### 4.2 `esc(s string) string`

Go 標準ライブラリ `html.EscapeString(s)` 相当の処理を行う。HTML 特殊文字（`<`、`>`、`&`、`"`、`'`）をエスケープする。

**実装契約：**
- Go 実装では `html.EscapeString` を使用する。
- 追加の独自置換、Markdown 記法変換、改行変換を行わない。
- 空文字は空文字を返す。

---

### 4.3 `inline(text string, ctx *RenderContext) string`

インライン Markdown 記法を HTML に変換する。コードスパンを先にキャラクターレベルのスキャンで処理することで、後続の正規表現がコードスパン内の記法を誤って変換するのを防ぐ。

**関連型：**

```go
type RenderContext struct {
    FootnoteDefs      map[string]string
    FootnoteOrder     []string
    FootnoteSeen      map[string]bool
    InternalLinkRefs  map[string]bool
    BrokenLinks       []string
    HeadingSkipCount  int
    CharCount         int
}
```

`ctx` は `nil` 禁止。呼び出し元は `convert()` 実行前に `RenderContext` を作成し、脚注定義を `FootnoteDefs` に収集しておく。

**処理手順：**

1. **コードスパン抽出（文字列スキャン）**
   - バッククォート（`` ` ``）の位置を順次検索
   - `` `` `` で始まる場合はダブルバッククォートコードスパンとして処理
   - 単一バッククォートの場合はシングルバッククォートコードスパンとして処理
   - コードスパン内テキストは `esc()` した上で `<code class="ic">` でラップ
   - コードスパン外テキストは `esc()` してそのままセグメントに追加

2. **インライン記法の正規表現置換**（コードスパン処理後のテキストに適用）

| パターン | 出力 |
|---------|------|
| `***text***` | `<strong><em>text</em></strong>` |
| `**text**` | `<strong>text</strong>` |
| `*text*`（単独 `*`） | `<em>text</em>` |
| `__text__` | `<strong>text</strong>` |
| `_text_`（単独 `_`） | `<em>text</em>` |
| `~~text~~` | `<del>text</del>` |
| `![alt](url)` | `<img src="url" alt="alt" style="max-width:100%">` |
| `[label](url)`（`url` が `http://` または `https://` で始まる場合） | `<a href="url" target="_blank" rel="noopener noreferrer">label</a>` |
| `[label](url)`（上記以外 — 内部リンク・アンカー） | `<a href="url">label</a>` |
| `[^id]` | `<sup><a href="#fn-id" id="fnref-id" class="fn-ref">[N]</a></sup>`（N は参照順の番号） |

**処理順の注意：** 画像記法（`![alt](url)`）はリンク記法（`[label](url)`）より先にマッチングする。リンク記法はマッチング後に `url` が `http://` または `https://` で始まるかを判定し、外部リンクと内部リンクを区別する。脚注参照（`[^id]`）はリンク置換後に適用する。

**内部リンク整合性チェック：**

`url` が `#` で始まる内部アンカーリンクを処理する際、アンカー部分（`#` 以降）を `ctx.InternalLinkRefs[anchor] = true` として記録する。`convert()` 末尾で生成済み slug の集合と照合し、一致しないアンカーを次のように報告する：

```
[WARN] BROKEN_LINK: #anchor-text  (in: [label](#anchor-text))
```

不一致件数は `[REPORT]` の `broken_links` フィールドに反映される。照合は変換完了後に行うため、ドキュメント内の順序（前方参照・後方参照）を問わず検証できる。

**脚注・内部リンク状態：**

| 変数 | 型 | 用途 |
|------|-----|------|
| `ctx.FootnoteDefs` | `map[string]string` | 脚注 ID → テキストの対応（`convert()` 呼び出し前に収集） |
| `ctx.FootnoteOrder` | `[]string` | 本文中での参照順脚注 ID リスト（`inline()` が参照時に追記） |
| `ctx.FootnoteSeen` | `map[string]bool` | 同一脚注 ID の重複登録防止 |
| `ctx.InternalLinkRefs` | `map[string]bool` | 本文中に出現した内部アンカー参照（`#` 除いた文字列）の集合 |

**脚注番号付与：**
- 初回参照時のみ `FootnoteOrder` へ ID を append する。
- 2 回目以降の同一 ID 参照では既存番号を再利用する。
- `FootnoteDefs` に存在しない ID でも HTML 参照は出力し、脚注本文は空文字として扱わず、末尾脚注出力時に `[WARN] MISSING_FOOTNOTE: id` を出す。

**インライン記法のネスト仕様：**

| 入力 | 出力 |
|------|------|
| `***text***` | `<strong><em>text</em></strong>` |
| `**_text_**` | `<strong><em>text</em></strong>` |
| `__*text*__` | `<strong><em>text</em></strong>` |
| `*__text__*` | `<em><strong>text</strong></em>` |
| `_**text**_` | `<em><strong>text</strong></em>` |

上表以外のネストした強調・削除・リンクの組み合わせは追加変換しない。未対応ネストは、先にマッチした外側または内側の単一記法だけを変換し、残った Markdown 記号は `esc()` 済みテキストとして出力する。実装者判断で CommonMark 全互換のネスト処理を追加してはならない。

---

### 4.4 `buildTOC(headings []Heading) string`

見出しリストから TOC（目次）の HTML を生成する。h1〜h3 のみを対象とし（h4 は除外）、子見出しを持つ見出しはグループとしてアコーディオン形式に構築する。

**関連型：**

```go
type Heading struct {
    Level      int
    Text       string
    Slug       string
    LineNumber int
}
```

**引数：** `headings []Heading`

**出力：** TOC の `<li>` 要素群の HTML 文字列（`<ul>` ルートは HTML テンプレート側で定義）

**グループ（`.tg`）とリーフ（`.ti`）の判定：**
現在の見出しの次の見出しレベルが現在より深い場合、グループとして扱い、展開ボタン（`.tg-btn`）付きの `<ul>` をネストする。それ以外はリーフ（`.ti`）として `<li>` 1 個を出力する。

**スタック管理：** 内部スタックは `[]tocStackItem` とする。

```go
type tocStackItem struct {
    Level int
    Kind  string // "group" または "item"
}
```

`closeTo(targetLevel int)` で不要になった `</ul></li>` を閉じる。スタックにはグループ（`"group"`）だけでなく、レベル管理のためリーフ（`"item"`）も積む。`closeTo()` はリーフをサイレントに捨て、グループのみ `</ul></li>` を出力する。

**生成 HTML 構造（グループの場合）：**
```html
<li class="tg">
  <div class="tg-row">
    <a href="#slug" class="tl lv1" data-slug="slug">見出しテキスト</a>
    <button class="tg-btn" aria-expanded="false" data-target="tg-slug" aria-label="展開">
      <svg>…</svg>
    </button>
  </div>
  <ul id="tg-slug" class="tc" hidden>
    …子アイテム…
  </ul>
</li>
```

**生成 HTML 構造（リーフの場合）：**
```html
<li class="ti">
  <a href="#slug" class="tl lv1" data-slug="slug">見出しテキスト</a>
</li>
```

**エラー：** 戻り値エラーは持たない。`Level` が 1〜4 以外の `Heading` は無視し、`[WARN] INVALID_HEADING_LEVEL: line={LineNumber}` を stdout へ出力する。警告件数は `[REPORT]` の `warnings` に含める。

---

### 4.5 `convert(lines []string, headings []Heading, ctx *RenderContext) ConvertResult`

Markdown の行リストを走査し、HTML コンテンツ文字列を生成するメインコンバーター。

**関連型：**

```go
type ConvertResult struct {
    HTML               string
    ReadingTimeMinutes int
    BrokenLinks        int
    HeadingSkips       int
    Warnings           []string
}
```

**引数：**
- `lines []string` — Markdown の全行。末尾改行は含めない。
- `headings []Heading` — 見出し収集済みリスト。各 `Heading.Slug` は重複解決済みとする。
- `ctx *RenderContext` — `inline()` と共有する変換状態。`nil` 禁止。

**戻り値：**
- `HTML` — 本文 HTML 文字列。
- `ReadingTimeMinutes` — §4.5 の読了時間算出結果。
- `BrokenLinks` — 内部リンク不一致件数。
- `HeadingSkips` — 見出し階層スキップ件数。
- `Warnings` — `[WARN]` として stdout 出力した警告本文の一覧。

**エラー：** 戻り値エラーは持たない。入力ファイル不存在、UTF-8 不正、書き込み失敗などの異常は §8 の実行方法と終了コードで扱う。Markdown 構文上の不足は警告またはフォールバック出力で処理する。

**内部バッファと状態変数：**

| 変数 | 型 | 用途 |
|------|-----|------|
| `out` | `strings.Builder` | 本文 HTML 断片の蓄積 |
| `fence_active` | `bool` | コードフェンス内かどうか |
| `fence_lang` | `string` | コードフェンスの言語識別子 |
| `fence_buf` | `[]string` | フェンス内の行バッファ |
| `fence_marker` | `string` | フェンス開始マーカー（`` ``` `` or `~~~`） |
| `para_buf` | `[]string` | 段落テキストの行バッファ |
| `list_stack` | `[]listStackItem` | リストのネスト状態スタック |
| `table_buf` | `[]string` | テーブル行のバッファ |
| `slugByLine` | `map[int]string` | `headings` から生成する行番号 → slug の対応 |

**内部ヘルパー関数：**

- `flushPara()` — `para_buf` を `<p class="mp">` として出力し、バッファをクリア
- `flushList()` — `list_stack` を巻き戻し、すべての `</ul>` / `</ol>` を閉じる
- `flushTable()` — `table_buf` をパースし `<div class="tw"><table class="mt">` として出力
- `emitCode(&out, fence_lang, fence_buf)` — `fence_buf` を `<div class="cb-wrap">` 構造として出力

**ブロック要素の検出優先順位（1 行ずつ処理）：**

1. フェンスコードブロック開始・終了（`` ``` `` / `~~~` で始まる行）
2. 見出し（`#` で始まる行）
3. 水平線（`---`、`***`、`___`）
4. テーブル（`|` で始まる行、連続する `|` 行をまとめて処理）
5. 引用（`>` で始まる行、連続する `>` でネスト可能）
6. 定義リスト（`: 定義` 形式の行かつ `para_buf` に用語がある場合）
7. リスト項目（`- * +` または `1.` 形式、`[ ]`/`[x]` プレフィックスでタスクリスト）
8. 空行（バッファのフラッシュトリガー）
9. 脚注定義行（`[^id]:` で始まる行、`ctx.FootnoteDefs` 収集済みのためスキップ）
10. 段落（上記以外の非空行、連続行を 1 つの `<p>` にまとめる）。先読みループは次のいずれかに該当する行で停止する：`#`（見出し）、`|`（テーブル）、`` ` ``×3以上（フェンス）、`~`×3以上（フェンス）、`>`（引用）、リストマーカー（`[-*+]` または `\d+[.)]`）、`: `（定義リストマーカー）、水平線（`---+`・`***+`・`___+`）

**テーブル変換の詳細：**
セパレーター行（`:---:`、`---` などで構成された行）のインデックスを自動検出し、セパレーター行より前の行をヘッダー（`<th>`）、それ以降を本文（`<td>`）として出力する。セパレーター行自体は出力しない。

**引用ネストの詳細：**
`>` で始まる連続行をまとめて収集し、`renderBlockquote(lines []string, ctx *RenderContext) string` が再帰的にネストを処理する。1 レベル分の `>` を剥いた後、内側行を先頭から走査し、`>` で始まる連続する行は `renderBlockquote()` を再帰呼び出し、それ以外の行は `inline(text, ctx)` でレンダリングして結合する。これにより、単一行・複数行・混在ネスト（同一ブロック内で `>` 行と `>>` 行が混在する場合）をすべて正しく処理する。例：`>> text` → `<blockquote class="mbq"><blockquote class="mbq">text</blockquote></blockquote>`。

**タスクリストの詳細：**
リスト項目のコンテンツが正規表現 `^\[([ xX])\]\s+` に一致する場合、`<li class="ml-task">` として出力する。チェック済み（`[x]` / `[X]`）は `checked` 属性付き、未チェック（`[ ]`）は属性なしの `<input type="checkbox" disabled>` を先頭に配置する。

**定義リストの詳細：**
`: 定義` 行（`strings.HasPrefix(line, ": ")`、コロン＋スペース1文字）を検出したとき `para_buf` に内容があれば、`para_buf` の末尾要素を用語（`<dt>`）として取り出し、`<dl class="mdl"><dt>用語</dt><dd>定義</dd></dl>` を出力する。連続する `: ` 行は同一 `<dl>` 内の追加 `<dd>` としてまとめて処理し、その後に `</dl>` を閉じる。

**脚注定義行のスキップ：**
`^\[\^[^\]]+\]:` に一致する行は `ctx.FootnoteDefs` への収集が完了しているためスキップし、本文への出力を行わない。

**主要ブロックの固定 HTML 断片：**

| Markdown | HTML |
|----------|------|
| `# Title` | `<h1 id="title" class="mh h1">Title<button class="hn-link" data-href="#title" aria-label="リンクをコピー">¶</button></h1>` |
| `text` | `<p class="mp">text</p>` |
| `---` | `<hr class="mr">` |
| `> quote` | `<blockquote class="mbq">quote</blockquote>` |
| `- item` | `<ul class="ml"><li>item</li></ul>` |
| `1. item` | `<ol class="ml"><li>item</li></ol>` |
| `- [x] done` | `<ul class="ml"><li class="ml-task"><input type="checkbox" disabled checked>done</li></ul>` |
| `term` + 次行 `: desc` | `<dl class="mdl"><dt>term</dt><dd>desc</dd></dl>` |

上表の属性順、class 名、button 文言、`checked` 属性の位置は固定する。テストでは空白の連続を 1 つへ正規化して比較してよいが、タグ名、属性名、属性値、親子構造は完全一致させる。

**脚注セクションの末尾出力：**
`convert()` 末尾で `ctx.FootnoteOrder` が非空の場合、`<section class="fn-section">` 内に参照順番号付きの脚注リスト（`<ol class="fn-list">`）を出力する。各脚注には本文への戻りリンク（`<a class="fn-back">↩</a>`）を付与する。

**フェンスコードブロックの未閉鎖フォールバック：**
ファイル末尾まで読んだ時点で `fence_active` が `true` のままの場合（閉じる `` ``` `` がない場合）、`fence_buf` にコンテンツがあれば `emitCode(&out, fence_lang, fence_buf)` を呼び出して強制出力する。`fence_buf` が空（フェンス開始直後に EOF）の場合は何も出力しない。未閉鎖フェンスは `[WARN] UNCLOSED_FENCE: line={startLine}` を出力し、`Warnings` に追加する。

**見出し階層スキップ警告：**

`convert()` 内で `prevHeadingLevel int = 0` をローカル変数として保持する。見出し行（`#` で始まる行）を処理するたびに現在レベルと前回レベルを比較し、2 段以上の降順スキップ（例：h1→h3、h2→h4）を検出した場合に次の形式で `[WARN]` を出力する：

```
[WARN] HEADING_SKIP: h1→h3 "見出しテキスト"
```

- 昇順への復帰（例：h3→h1）はスキップに該当しない（章の区切りとして正常）
- 同レベルの連続（h2→h2）・1段降順（h2→h3）もスキップに該当しない
- 件数は `[REPORT]` の `heading_skips` フィールドに反映される

**読了時間集計：**

`convert()` 内で `ctx.CharCount = 0` に初期化する。段落・リスト・引用テキストを `inline(text, ctx)` 処理する直前に、元の Markdown テキスト文字数（スペース・改行を含む）を加算する。コードブロック・フェンス内テキスト・見出しテキスト・テーブルは集計対象外とする。

読了時間の算出：
```go
readingTimeMinutes := int(math.Ceil(float64(ctx.CharCount) / 200.0)) // 200文字/分、切り上げ
```

算出した `readingTimeMinutes` は `ConvertResult.ReadingTimeMinutes` として返し、呼び出し元が `[REPORT]` 行と `.build_logs/{id}.json` に記録する。また HTML ヘッダーへの静的埋め込み（§5）にも使用する。

**見出し出力 HTML 構造：**
`#` で始まる行を `h1`〜`h4` に変換する際、末尾に `.hn-link` ボタンを付与する。

```html
<h2 id="slug" class="mh h2">見出しテキスト<button class="hn-link" data-href="#slug" aria-label="リンクをコピー">¶</button></h2>
```

- `data-href` 属性：`#` + 一意化後 slug
- `¶`（U+00B6 PILCROW SIGN）を使用
- CSS で通常時 `opacity: 0`、親見出し要素のホバー時に `opacity: 1` に変化する

**見出しスラグ重複解決：**

見出し収集処理では `slugCount map[string]int` をローカル変数として保持し、同一 base slug が複数の見出しに割り当てられる場合に一意化する。

| 条件 | スラグ |
|------|--------|
| 初出 | `{slug}` そのまま |
| 2 回目 | `{slug}-2` |
| 3 回目 | `{slug}-3` |
| N 回目 | `{slug}-{N}` |

```go
slugCount := map[string]int{}

uniqueSlug := func(base string) string {
    n := slugCount[base] + 1
    slugCount[base] = n
    if n == 1 {
        return base
    }
    return fmt.Sprintf("%s-%d", base, n)
}
```

一意化後のスラグは `id` 属性・`data-href` 属性・TOC リンク `href`（§4.4）・`¶` ボタン・全文検索インデックス（§7.9）・前後章ナビゲーション（§7.15）のすべてで共通使用する。

---

### 4.6 `emitCode(out *strings.Builder, lang string, buf []string)`

フェンスコードブロックを HTML に変換して `out` に追記する。`convert()` 内クロージャではなく、Go ファイル内の非公開補助関数として実装する。

**入力：**

| 引数 | 型 | 条件 |
|------|----|------|
| `out` | `*strings.Builder` | `nil` 禁止。呼び出し元の本文 HTML 出力先。 |
| `lang` | `string` | フェンス開始行の言語識別子。未指定時は空文字。 |
| `buf` | `[]string` | フェンス内本文。空配列を許可する。 |

**出力 HTML 構造：**
```html
<div class="cb-wrap" data-lang="rust">
  <div class="cb-meta">
    <span class="cl">rust</span>
    <button class="cb-copy" aria-label="コピー">コピー</button>
  </div>
  <pre class="cb"><code>…エスケープ済みコード…</code></pre>
</div>
```

- `data-lang` 属性：言語識別子（なければ空文字列）
- `.cl`（言語ラベル）：言語識別子がない場合は出力しない
- コード内容：`strings.Join(buf, "\n")` を `esc()` でエスケープしてから出力
- `buf` が空の場合でも空の `<pre class="cb"><code></code></pre>` を出力する

**エラー：** 戻り値エラーは持たない。`out == nil` はプログラム不備として `panic("nil output builder")` を許可する。

---

## 5. 静的 Web サイト出力構造

サイト全体の合成は `assembleSite(site SiteData, theme Theme) error` が担当する。`convert()` は各 Markdown ページの本文 HTML を生成し、`assembleSite()` はページ HTML、共通 CSS、共通 JavaScript、検索インデックス、テーマコンポーネントを出力サイトディレクトリへ書き出す。

**関連型：**

```go
type SiteData struct {
    Title          string
    Pages          []PageData
    SearchIndex    []SearchIndexEntry
    GeneratedAtUTC string
}

type PageData struct {
    Title              string
    SourcePath         string
    OutputPath         string
    Layout             string
    TocHTML            string
    BodyHTML           string
    ReadingTimeMinutes int
}

type Theme struct {
    Name       string
    Components ThemeComponents
    CSS        string
    JS         string
}

type ThemeComponents struct {
    Header     string
    Sidebar    string
    Breadcrumb string
    TOC        string
    Search     string
    Footer     string
    CodeBlock  string
    Table      string
    Pagination string
}

type SearchIndexEntry struct {
    URL   string
    ID    string
    Title string
    Body  string
}
```

**入力契約：**

| フィールド | 型 | 条件 |
|------------|----|------|
| `SiteData.Title` | `string` | 空は禁止。HTML 出力時は `esc()` する。 |
| `SiteData.Pages` | `[]PageData` | 1 件以上。単一 Markdown 入力の場合も 1 ページのサイトとして扱う。 |
| `SiteData.SearchIndex` | `[]SearchIndexEntry` | サイト内検索用。未生成時は空配列。 |
| `SearchIndexEntry.URL` | `string` | `PageData.OutputPath` または `PageData.OutputPath + "#" + headingSlug`。空は禁止。 |
| `SearchIndexEntry.ID` | `string` | 見出し slug。ページ単位エントリの場合は空文字を許可する。 |
| `SearchIndexEntry.Title` | `string` | 検索結果に表示するタイトル。空は禁止。 |
| `SearchIndexEntry.Body` | `string` | 段落先頭 200 文字以内の説明文。該当本文がない場合は空文字を許可する。 |
| `Theme.Name` | `string` | 初期仕様では `adlaire-default` 固定。 |
| `PageData.Title` | `string` | 空の場合は `SiteData.Title` を使用する。HTML 出力時は `esc()` する。 |
| `PageData.SourcePath` | `string` | 入力 Markdown の絶対パスまたは `--base-dir` からの相対パス。 |
| `PageData.OutputPath` | `string` | `--out` からの相対 HTML パス。トップページは `index.html`。 |
| `PageData.Layout` | `string` | 初期仕様では `document` または `index`。未知値は禁止。 |
| `TocHTML` | `string` | `buildTOC(headings)` の戻り値。`<ul id="toc-root">` の内側へ挿入する。 |
| `BodyHTML` | `string` | `ConvertResult.HTML`。`<div class="ci">` の内側へ挿入する。 |
| `ReadingTimeMinutes` | `int` | 1 以上。0 以下の場合は `1` として表示する。 |
| `GeneratedAtUTC` | `string` | UTC ISO 8601。空の場合は生成日時 meta を出力しない。 |

`assembleSite()` は上記フィールドを結合してファイルを書き出すだけとし、Markdown 変換、slug 生成、TOC 生成、検索インデックス抽出、警告集計を行ってはならない。

**出力ディレクトリ構造：**

```text
{out}/
  index.html
  pages/
    {slug}.html
  assets/
    style.css
    app.js
    search-index.json
```

単一 Markdown 入力の場合、本文ページを `index.html` として出力し、`pages/` は作成しなくてよい。Markdown ディレクトリ入力の場合、`index.html` はサイト目次ページとし、各 Markdown ファイルを `pages/{slug}.html` として出力する。

**相対 root 算出：**

`relativeRoot` は、各 HTML ファイルから `assets/` へ到達するための相対 prefix とする。

| `PageData.OutputPath` | `relativeRoot` |
|-----------------------|----------------|
| `index.html` | `""` |
| `pages/example.html` | `"../"` |
| `pages/dir-example.html` | `"../"` |

初期仕様ではページ HTML を `pages/` 直下に平坦化するため、`relativeRoot` は上表の 2 種類のみとする。`pages/dir/page.html` のような階層出力を追加する場合は、先に本表、リンク解決、検索 index URL、breadcrumb 仕様を改訂する。

**出力更新手順：**

`assembleSite()` は `--out` を直接途中更新してはならない。以下の順で一時ディレクトリへ完全生成してから置換する。

1. `--out` と同じ親ディレクトリに `{basename}.tmp.{pid}` を作成する。
2. 一時ディレクトリ配下へ `index.html`、ページ HTML、`assets/` をすべて書き出す。
3. すべてのファイルについて `file.Sync()` と `file.Close()` を完了する。
4. 一時ディレクトリ内の必須ファイル存在を検証する。
5. 既存 `--out` が存在する場合は `{basename}.prev.{pid}` へ `os.Rename` する。
6. 一時ディレクトリを `--out` へ `os.Rename` する。
7. 置換成功後、旧 `{basename}.prev.{pid}` を削除する。

手順 1〜4 で失敗した場合は一時ディレクトリを削除し、既存 `--out` を変更してはならない。手順 6 で失敗した場合は `{basename}.prev.{pid}` が存在するか確認し、存在する場合は `{basename}.prev.{pid}` を `--out` へ戻す復旧 rename を必ず 1 回試行する。復旧 rename が失敗した場合は stderr に `cannot restore previous output directory: <path>` を出力した後、続けて `cannot replace output directory: <path>` を出力し、終了コード `1` とする。`{basename}.prev.{pid}` が存在しない場合は `cannot replace output directory: <path>` のみを出力し、終了コード `1` とする。

`--out` が既存ファイルでディレクトリではない場合は終了コード `1` とし、stderr に `output path is not directory: <path>` を出力する。

**初期テーマコンポーネント：**

初期実装の theme は `adlaire-default` のみとする。theme component は Go コード内の内製テンプレートとして保持し、外部テンプレートファイルを読み込んではならない。

| コンポーネント | 責務 |
|----------------|------|
| `header` | サイト名、現在ページ名、読了時間を表示する。 |
| `sidebar` | サイト内ページ一覧と現在ページの TOC を表示する。 |
| `breadcrumb` | `index.html` から現在ページまでの階層を表示する。 |
| `toc` | 現在ページの見出し TOC を表示する。 |
| `search` | `assets/search-index.json` を読み込み、サイト内検索 UI を提供する。 |
| `footer` | 生成時刻とサイト名を表示する。 |
| `codeblock` | 言語ラベル、コピー操作、折りたたみ状態を提供する。 |
| `table` | 横スクロール wrapper と列 sort を提供する。 |
| `pagination` | 前後ページへのリンクを表示する。 |

カスタムテーマ、外部テンプレート、テーマパッケージ、theme component 差し替え、複数 theme 同梱の入出力、状態、検証条件は本ファイルでは定義しない。

**ページ HTML の必須 DOM 構造：**

```html
<!DOCTYPE html>
<html lang="ja">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="adlaire-build-id" content="{BuildMeta.BuildID}">
  <meta name="adlaire-commit-sha" content="{BuildMeta.CommitSHA}">
  <meta name="adlaire-build-at" content="{BuildMeta.BuildAt}">
  <title>{PageData.Title}</title>
  <link rel="stylesheet" href="{relativeRoot}assets/style.css">
</head>
<body>
  <div id="progress-bar"></div>   <!-- 読み取り進捗バー（ページ上端固定、高さ 3px、幅 = スクロール率 % → §7.13） -->
  <header id="hdr">        <!-- 固定ヘッダー（高さ 52px、背景 --adlaire-surface-accent） -->
    <span id="doc-title">{PageData.Title}</span>
    <span id="reading-time">約 {PageData.ReadingTimeMinutes} 分</span>  <!-- 読了時間（ビルド時に静的埋め込み → §4.5・§6） -->
  </header>
  <div id="lay">           <!-- フレックスコンテナ -->
    <nav id="sb">          <!-- サイドバー（幅 260px、固定） -->
      <div class="sb-search-wrap">
        <input id="sb-search" type="search" autocomplete="off" aria-label="セクション検索">  <!-- TOC 検索 -->
      </div>
      <div id="sb-toc">
        <ul class="tr" id="toc-root">{PageData.TocHTML}</ul>  <!-- TOC リスト -->
        <p id="sb-none" class="sb-none" hidden>一致するセクションはありません</p>
      </div>
    </nav>
    <main id="ct">         <!-- コンテンツエリア -->
      <div class="ci">     <!-- 最大幅 760px センタリングコンテナ -->
        {PageData.BodyHTML}
      </div>
    </main>
  </div>
  <button id="btt" type="button" aria-label="トップへ戻る">↑</button>        <!-- トップへ戻るボタン -->
  <script src="{relativeRoot}assets/app.js"></script>
</body>
</html>
```

**ディレクトリ入力時の `index.html` 仕様：**

`index.html` はサイト目次ページとして生成し、Markdown 本文を持たない。`PageData.Layout` は `index` とする。

| 項目 | 仕様 |
|------|------|
| `<title>` | `SiteData.Title` |
| `#doc-title` | `SiteData.Title` |
| `.ci` 内 | `<h1 id="site-index" class="mh h1">{SiteData.Title}<button class="hn-link" data-href="#site-index" aria-label="リンクをコピー">¶</button></h1>` の後に `<ul class="site-page-list">` を出力する。 |
| ページ一覧 | `SiteData.Pages` のうち `Layout == "document"` のページを入力ソート順で `<li><a href="{OutputPath}">{Title}</a></li>` として出力する。 |
| TOC | サイト目次ページの TOC は `site-index` 1 件だけを含める。 |
| 検索 index | サイト目次ページのエントリを 1 件追加し、`url` は `index.html#site-index`、`id` は `site-index`、`title` は `SiteData.Title`、`body` は空文字とする。 |

**必須 DOM ID / class 契約：**

| セレクター | 個数 | 使用者 | 変更可否 |
|------------|------|--------|----------|
| `#progress-bar` | 各 HTML 1 個 | `assets/app.js` §7.13 | 変更禁止 |
| `#hdr` | 各 HTML 1 個 | CSS / layout | 変更禁止 |
| `#doc-title` | 各 HTML 1 個 | header | 変更禁止 |
| `#reading-time` | 各 HTML 1 個 | header | 変更禁止 |
| `#lay` | 各 HTML 1 個 | CSS / sidebar layout | 変更禁止 |
| `#sb` | 各 HTML 1 個 | sidebar JS | 変更禁止 |
| `#sb-search` | 各 HTML 1 個 | TOC/search JS | 変更禁止 |
| `#toc-root` | 各 HTML 1 個 | TOC JS | 変更禁止 |
| `#ct` | 各 HTML 1 個 | content layout | 変更禁止 |
| `.ci` | 各 HTML 1 個 | content width | 変更禁止 |
| `#btt` | 各 HTML 1 個 | top button JS | 変更禁止 |

実装は上表のセレクターを追加、削除、リネームしてはならない。UI 改善で新しいセレクターが必要な場合は、先に本表へ追加し、§6 と §7 の CSS / JavaScript 契約を同時に更新する。

**サイト合成の禁止事項：**
- `PageData.BodyHTML`、`PageData.TocHTML` はすでに HTML として生成済みのため、`assembleSite()` 内で再エスケープしない。
- `assets/search-index.json` は `encoding/json` の出力だけを受け付け、文字列連結で JSON を自作しない。
- `<header id="hdr">`、`<nav id="sb">`、`<main id="ct">`、`<div class="ci">` の id / class を変更しない。
- `assets/style.css`、`assets/app.js` 以外の CSS / JavaScript を生成してはならない。
- 外部 CSS、外部 JavaScript、外部フォント参照を追加しない。

---

## 6. CSS クラス一覧

### コンテンツ要素

| クラス | 要素 | 説明 |
|--------|------|------|
| `.mh` | `h1`〜`h4` | 見出し共通スタイル |
| `.h1`〜`.h4` | `h1`〜`h4` | 見出しレベル別スタイル |
| `.mp` | `p` | 本文段落（`max-width: 68ch`） |
| `.mr` | `hr` | 水平線 |
| `.mbq` | `blockquote` | 引用ブロック |
| `.ic` | `code` | インラインコード |
| `.ml` | `ul`/`ol` | リスト |

### コードブロック

| クラス | 説明 |
|--------|------|
| `.cb-wrap` | コードブロック外枠（`position: relative`） |
| `.cb-meta` | 言語ラベル＋コピーボタンのオーバーレイ（`position: absolute; top: 8px; right: 10px`） |
| `.cl` | 言語ラベル（`var(--adlaire-font-family-mono)`、大文字、`var(--adlaire-font-size-xs)`） |
| `.cb-copy` | コピーボタン（ホバーまで非表示） |
| `.cb-copy.copied` | コピー完了状態（1.8 秒間） |
| `.cb` | `<pre>` 要素（背景 `var(--adlaire-surface-soft-strong)`） |

### タスクリスト

| クラス | 説明 |
|--------|------|
| `.ml-task` | タスクリスト `<li>`（`list-style: none`、チェックボックス付き） |
| `.ml-task input[type="checkbox"]` | チェックボックス（`disabled`、`accent-color: var(--adlaire-color-primary)`） |

### 定義リスト

| クラス | 説明 |
|--------|------|
| `.mdl` | `<dl>` 要素 |
| `.mdl dt` | 定義用語（`font-weight: semibold`） |
| `.mdl dd` | 定義本文（`margin-left: 1.5rem`、`color: var(--adlaire-surface-text-muted)`） |

### テーブル

| クラス | 説明 |
|--------|------|
| `.tw` | テーブルラッパー（`overflow-x: auto`） |
| `.mt` | `<table>` 要素 |

### 脚注

| クラス | 説明 |
|--------|------|
| `.fn-ref` | 脚注参照リンク（`<sup>` 内、`font-size: var(--adlaire-font-size-xs)`） |
| `.fn-section` | 脚注セクション全体（本文末尾、`border-top` で区切り） |
| `.fn-list` | 脚注 `<ol>`（`font-size: var(--adlaire-font-size-sm)`） |
| `.fn-item` | 脚注 `<li>` |
| `.fn-back` | 本文への戻りリンク（`↩`） |

### TOC（サイドバー）

| クラス | 説明 |
|--------|------|
| `.tr` | TOC ルート `<ul>` |
| `.tg` | グループ（子を持つ見出し） |
| `.tg-row` | グループの行（リンク＋展開ボタン） |
| `.tg-btn` | グループ展開トグルボタン |
| `.tc` | グループの子 `<ul>`（`hidden` で折りたたみ） |
| `.ti` | リーフ（子を持たない見出し） |
| `.tl` | TOC リンク（`<a>` 要素） |
| `.tl.lv1`〜`.tl.lv3` | TOC リンクのレベル別スタイル |
| `.tl.active` | アクティブな TOC リンク |
| `.sb-none` | 検索結果なしメッセージ |

### シンタックスハイライト CSS クラス（§7.8）

| クラス | 対象トークン |
|--------|------------|
| `.hl-kw` | キーワード |
| `.hl-str` | 文字列リテラル |
| `.hl-num` | 数値リテラル |
| `.hl-cmt` | コメント |
| `.hl-key` | JSON オブジェクトキー |
| `.hl-op` | diff 追加行（`+`）・削除行（`-`） |

### 見出しアンカーリンクコピー CSS クラス（§7.11）

| クラス | 要素 | 説明 |
|--------|------|------|
| `.hn-link` | `<button>` | 見出し末尾に付与する `¶` ボタン（通常時 `opacity: 0`、ホバー時 `opacity: 1`） |

### 進捗バー・テーブルソート CSS クラス（§7.13・§7.14）

| クラス / セレクター | 要素 | 説明 |
|--------------------|------|------|
| `#progress-bar` | `<div>` | ページ上端固定（`position: fixed; top: 0; left: 0`）、高さ 3px、幅は JS で設定、`background: var(--adlaire-color-primary)`、`z-index: 1000`、`transition: width 0.1s linear` |
| `.mt th[data-sort]` | `<th>` | ソート可能列ヘッダー（`cursor: pointer; user-select: none`） |
| `.mt th[aria-sort="ascending"]::after` | `::after` 疑似要素 | `content: " ▲"` |
| `.mt th[aria-sort="descending"]::after` | `::after` 疑似要素 | `content: " ▼"` |

### 読了時間表示 CSS クラス（§4.5）

| クラス / セレクター | 要素 | 説明 |
|--------------------|------|------|
| `#reading-time` | `<span>` | 固定ヘッダー右端に表示（`margin-left: auto`）。`color: var(--adlaire-text-secondary)`、`font-size: var(--adlaire-font-size-sm)`、`white-space: nowrap` |

### 前後章ナビゲーション CSS クラス（§7.15）

| クラス / セレクター | 要素 | 説明 |
|--------------------|------|------|
| `.ch-nav` | `<nav>` | 章末尾ナビゲーション（`display: flex; justify-content: space-between; padding: 1rem 0; margin-top: 2rem; border-top: 1px solid var(--adlaire-border-default)`） |
| `.ch-prev` | `<a>` | 前の章リンク（`color: var(--adlaire-color-primary)`、テキスト装飾なし） |
| `.ch-next` | `<a>` | 次の章リンク（`color: var(--adlaire-color-primary)`、テキスト装飾なし） |

### 印刷スタイル（§6 @media print）

`@media print` ブロックでの主な規則：

| 対象 | 印刷時の処理 |
|------|------------|
| `#hdr`（固定ヘッダー）、`#sb`（サイドバー）、`.cb-copy`、`.hn-link`、`#btt`、`.expand-code`、`#progress-bar`、`.ch-nav` | `display: none` |
| `#lay`、`#main` | ブロック表示・幅 100%・余白 0 |
| `pre.cb[data-collapsible]` | `max-height: none`（折りたたみ解除） |
| `a[href^="http"]::after`、`a[href^="https"]::after` | `content: " (" attr(href) ")"` で URL を末尾に表示 |
| `h2`、`h3` | `page-break-before: avoid` |

---

## 7. JavaScript 機能

### 7.1 テーマ切り替え（廃止）

ADS 採用により、ダークモードおよびテーマトグルボタンは廃止。
出力サイトはライトモード固定（`prefers-color-scheme` 非対応）。

### 7.2 サイドバー開閉

`localStorage` キー `adb-sb` に `"1"`（開）または `"0"`（閉）を保存する。

- デスクトップ（`> 768px`）：`#sb.closed` クラスと JS インラインスタイル（`ct.style.marginLeft`）で幅を制御。CSS の `#sb.closed ~ #ct { margin-left: 0 }` ルールは初期レンダリング時のみ効く。以降の開閉操作はすべて JS インラインスタイルが CSS クラスより優先する
- モバイル（`≤ 768px`）：`#sb.open` / `transform: translateX` で画面外から引き出す
- モバイルでは TOC リンククリック時に自動的にサイドバーを閉じる

### 7.3 TOC グループ展開

`.tg-btn` クリックで対応する `<ul id="tg-{slug}">` の `hidden` 属性をトグルし、`aria-expanded` 属性を更新する。`data-target` 属性で対象 `<ul>` の ID を指定する。

**開閉状態の永続化：**
展開操作のたびに現在展開中のグループのスラグ配列を `localStorage` キー `adlaire-toc-state` に JSON 文字列で保存する。ページ読み込み時（`DOMContentLoaded`）に同キーを読み込み、保存済みスラグのグループを展開状態で描画する。`localStorage` アクセスはすべて `try/catch` で保護し、失敗時はデフォルト状態（初期展開なし）にフォールバックする。

### 7.4 TOC 検索フィルター

`#sb-search` の `input` イベントで TOC 項目をリアルタイムフィルタリングする。

- 各 `.tl` のテキストコンテンツと検索クエリ（小文字化）を照合
- 非一致の `<li>` は `hidden = true` で非表示
- グループは子に一致項目があれば表示を維持し、`<ul>` を強制展開
- 入力がクリアされた場合は全項目を復元し、展開状態を `aria-expanded` に従って復元
- 一致なしの場合は `#sb-none` メッセージを表示

**状態変数 `searchActive`：** 検索中はアクティブ TOC リンクの自動スクロールを抑制する。

### 7.5 アクティブ見出し追跡

`IntersectionObserver` で `.mh[id]` 要素のビューポート内への進入を監視する。

- `rootMargin: "-8% 0px -78% 0px"` — 画面上部 8% 〜 22% の帯域内の見出しを「アクティブ」とみなす
- 進入した見出しの `id` に対応する TOC リンクに `.active` クラスを付与
- 親グループが折りたたまれている場合は自動展開し `aria-expanded="true"` を設定
- 検索中でない場合は対応 TOC リンクを `scrollIntoView` で可視範囲にスクロール

**進捗バー連動：** `scroll` イベントリスナー（`passive: true`）を同一リスナーで共有し、スクロールのたびに読み取り進捗バーの幅を更新する（→ §7.13）。

### 7.6 コピーボタン

各 `.cb-copy` ボタンのクリックで、親 `.cb-wrap` 内の `<code>` の `innerText` をクリップボードにコピーする。

```
navigator.clipboard.writeText() → 成功: done()
                                 → 失敗: fallbackCopy() → done()
fallbackCopy(): <textarea> を一時生成して execCommand('copy') を実行
done(): ボタンテキストを "✓ 完了" に変更、.copied クラス付与、1.8 秒後に "コピー" に戻す
```

### 7.7 トップへ戻るボタン

`scroll` イベント（`passive: true`）で `scrollY > 400` の場合に `#btt.visible` クラスを付与する。クリックで `window.scrollTo({ top: 0, behavior: 'smooth' })` を実行する。

### 7.8 シンタックスハイライト

コードブロックに言語別の色分けを `assets/app.js` で適用する。外部ライブラリ不要。

**対応言語：** `python` / `bash` / `json` / `sql` / `ini` / `diff`

**実装方式：**
各言語ごとにトークン正規表現パターンを定義し、`<code>` 要素のテキストに対して順次マッチを走らせる。マッチしたトークンを `<span class="hl-{type}">` でラップしてから `innerHTML` に書き戻す。

| CSS クラス | 対象トークン |
|---|---|
| `hl-kw` | キーワード（`def`, `class`, `if`, `SELECT` 等） |
| `hl-str` | 文字列リテラル（シングル／ダブルクォート） |
| `hl-num` | 数値リテラル |
| `hl-cmt` | コメント（`#`・`//`・`--` 行コメント） |
| `hl-key` | オブジェクトキー（JSON の `"key":` パターン） |
| `hl-op` | diff の追加行（`+`）・削除行（`-`） |

**適用タイミング：** `DOMContentLoaded` 後に全 `.cb-wrap` 要素を走査し、`data-lang` 属性の値を参照して言語を判定する（`data-lang` は `.cb-wrap` div に付与されており、`<code>` 要素には付与されない）。ハイライト処理は配下の `<code>` 要素のテキストに対して行う。`data-lang` が未設定または対応外の場合はハイライトをスキップする。

### 7.9 本文内全文検索

ビルド時に `assets/search-index.json` を生成し、`assets/app.js` の検索 UI から読み込む。TOC 検索フィルター（§7.4）と検索 UI を統合し、本文ヒット箇所へのジャンプを提供する。

**インデックス生成仕様（components/builder.go）：**
ビルド時に全ページの見出しと各段落の先頭 200 文字を抽出し、以下の配列形式で `assets/search-index.json` に書き出す。

```json
[
  { "url": "pages/example.html#anchor-slug", "id": "anchor-slug", "title": "見出しテキスト", "body": "段落先頭200文字..." },
  ...
]
```

**検索 index 固定契約：**

- JSON の top-level は配列とする。object wrapper は使用しない。
- entry の key 順は `url`、`id`、`title`、`body` とする。
- `encoding/json` の標準エスケープを使用し、独自 pretty print は行わない。
- entry の並び順は、ページ入力順、同一ページ内では見出し出現順とする。ページ単位エントリは各ページの先頭に 1 件だけ追加する。
- `body` は HTML タグ除去後のプレーンテキストを Unicode rune 単位で最大 200 文字とする。200 文字を超える場合は 200 文字で切り、三点リーダーを追加しない。
- 空白、タブ、改行の連続は半角スペース 1 つへ正規化し、前後空白を削除する。
- `title` と `body` に HTML entity を残してはならない。検索 index は表示前に JS 側で `textContent` として挿入し、`innerHTML` へ直接挿入しない。

**検索 UI の配置：**
§7.4 の TOC 検索フィルター入力欄を兼用する。入力値が 2 文字以上になった時点でインデックスに対して部分一致検索を実行する。

**ヒット箇所ハイライト：**
一致したエントリの見出しを検索結果として表示する。現在ページ内のヒットは TOC 内でハイライト（`.toc-hit` クラス付与）し、クリックで対象アンカーへスクロールする。別ページのヒットは `url` へ遷移する。現在ページ内では `<mark>` 要素でヒット文字列をページ内マーキングする（外部依存なし・標準 DOM 操作のみ）。

**クリア：**
入力欄を空にすると TOC ハイライトおよびページ内マーキングをすべて解除する。

### 7.10 コードブロック折りたたみ

30 行超のコードブロックを初期折りたたみ状態でレンダリングし、ユーザー操作で全行表示に切り替える。

**閾値：** 30 行（空行を含む総行数）。30 行以下のブロックは折りたたみ UI を生成しない。

**初期状態：**
折りたたみ対象の `<pre>` に `data-collapsible="true"` と `data-total-lines="{N}"` を付与する。CSS で `max-height` を制限（先頭 10 行相当）し、下端をグラデーションフェードでマスクする。

**展開リンク：**
ブロック末尾に `<button class="expand-code">全 {N} 行を表示</button>` を配置する。クリックで `max-height` を解除し、ボタンを非表示にする（折りたたみへの再折りたたみ機能は提供しない）。

**コピーボタンとの共存：**
§7.6 のコピーボタンは折りたたみ状態でも常時表示する。コピー操作は全行テキストを対象とする（表示行のみではない）。

### 7.11 見出しアンカーリンクコピー

見出しにホバーすると表示される `¶` ボタン（`.hn-link`）をクリックすると、その見出しのアンカー URL をクリップボードにコピーする。

**コピー対象 URL：**
`window.location.origin + window.location.pathname + button.dataset.href`

**実装：**
§7.6 と同じ Clipboard API（`navigator.clipboard.writeText()`）を使用し、同一のフォールバック（`execCommand('copy')`）を流用する。コピー完了フィードバックは `aria-label` を `"コピーしました"` に一時変更し、1.8 秒後に `"リンクをコピー"` へ復元する（§7.6 のコピーボタンと同じタイミング）。

**イベント登録：**
`DOMContentLoaded` 後に `document.querySelectorAll('.hn-link')` を走査して `click` リスナーを登録する。

### 7.12 キーボードショートカット

`keydown` イベントで以下のショートカットを処理する。フォーカスが `<input>`・`<textarea>`・`<select>` にある場合は `/` と `t` を無効にする（`Escape` のみ有効）。

| キー | 動作 |
|------|------|
| `/` | TOC 検索欄（`#sb-search`）にフォーカスを移動し、ページスクロールを抑止（`event.preventDefault()`） |
| `Escape` | TOC 検索欄の内容をクリアし、検索結果をリセットする（§7.4 の検索リセット処理と同等） |
| `t` | ページ先頭へスクロール（`window.scrollTo({ top: 0, behavior: 'smooth' })`） |

**イベント登録：** `document.addEventListener('keydown', handler)` を `DOMContentLoaded` 後に登録する。

### 7.13 読み取り進捗バー

ページ上端に高さ 3px の進捗バー（`<div id="progress-bar">`）を固定表示する（→ §5 HTML 出力構造）。

**幅の計算：**
```js
const scrolled = document.documentElement.scrollTop;
const total    = document.documentElement.scrollHeight
                 - document.documentElement.clientHeight;
const pct      = total > 0 ? (scrolled / total * 100) : 0;
document.getElementById('progress-bar').style.width = pct + '%';
```

**スクロールイベント共有：** §7.5 のアクティブ見出し追跡が登録する `scroll` イベントリスナー（`passive: true`）内で処理する。リスナーを別途登録しない。

**CSS：**

```css
#progress-bar {
  position: fixed;
  top: 0;
  left: 0;
  height: 3px;
  width: 0%;
  background: var(--adlaire-color-primary);
  z-index: 1000;
  transition: width 0.1s linear;
}
```

### 7.14 テーブル列ソート

`<th>` クリックで列を昇順／降順ソートする。

**HTML 構造：** `flushTable()` が生成する全 `<th>` に `data-sort="{col_index}"` 属性（0始まりの列インデックス）と `aria-sort="none"` を付与する。

```html
<th data-sort="0" aria-sort="none">列名</th>
```

**ソートアルゴリズム：**

1. クリックされた `<th>` の `aria-sort` を確認し、`"ascending"` → `"descending"`、それ以外 → `"ascending"` に決定
2. 同テーブル内の他の全 `<th>` の `aria-sort` を `"none"` にリセット
3. 対象列のセルテキスト（`textContent.trim()`）を取得し、`Number()` で数値変換可能ならば数値比較、それ以外は文字列比較（`localeCompare`）でソート
4. `<tbody>` の子 `<tr>` を並び替えて再挿入
5. クリックされた `<th>` の `aria-sort` を決定した値に更新（CSS `::after` でインジケーター表示）

**イベント登録：** `DOMContentLoaded` 後に `document.querySelectorAll('.mt th[data-sort]')` を走査して `click` リスナーを登録する。

**スコープ：** 同一テーブル内のソートのみ。複数列ソートは対象外。

---

### 7.15 前後章ナビゲーションボタン

h2 見出し単位で「← 前の章」「次の章 →」ボタンを各章末尾に静的生成する（→ §4.5・§5・§6 CSS）。

**生成方法：** `injectChapterNavigation(bodyHTML string, headings []Heading) string` として実装する。`convert()` が `ConvertResult.HTML` を返した後、呼び出し元が h2 見出しの `Slug` と `Text` を `headings` から抽出し、各 h2 章の末尾（次の h2 の直前、または文書末）に `<nav class="ch-nav">` を挿入する。

`injectChapterNavigation()` は `Heading.Level == 2` の見出しだけを対象とする。対象 h2 が 0 件または 1 件の場合、章ナビゲーションを挿入しない。

**HTML 構造：**

```html
<!-- 先頭章（.ch-prev なし） -->
<nav class="ch-nav">
  <span></span>
  <a class="ch-next" href="#next-slug">次の章タイトル →</a>
</nav>

<!-- 中間章 -->
<nav class="ch-nav">
  <a class="ch-prev" href="#prev-slug">← 前の章タイトル</a>
  <a class="ch-next" href="#next-slug">次の章タイトル →</a>
</nav>

<!-- 最終章（.ch-next なし） -->
<nav class="ch-nav">
  <a class="ch-prev" href="#prev-slug">← 前の章タイトル</a>
  <span></span>
</nav>
```

**スラグの参照：** `href` に使用するスラグは §4.5 のスラグ重複解決後の一意スラグを使用する。

**スコープ：** h2 レベルの見出しのみ。h3 以下の小節には生成しない。

**挿入失敗時の扱い：** 対象 h2 の HTML 位置を特定できない場合、本文 HTML を変更せず `[WARN] CHAPTER_NAV_SKIPPED: slug={slug}` を出力し、`[REPORT] warnings` に含める。

**印刷時：** `@media print` で `.ch-nav { display: none }` とする（§6 CSS 参照）。

---

## 8. 実行方法

```bash
/usr/local/bin/adlaire-ci-build
```

`adlaire-ci-build` の実行は、§2 の CLI 引数仕様に従う。引数なしの場合は `DefaultBuildConfig` の `Src`、`Out`、`Title`、`Theme`、`BaseDir`、`Strict` を使用する。

**実行順序契約：**

1. CLI 引数を検証する。`--help` / `--version` はここで処理し、Markdown 読み込みを行わない。
2. `--theme` を検証し、`adlaire-default` の `Theme` を選択する。
3. `--src` がファイルの場合は 1 ページ、ディレクトリの場合は配下の `.md` ファイルを辞書順に収集する。
4. 各 Markdown を UTF-8 として読み込み、行配列 `lines []string` を作成する。
5. 各ページの見出しを収集し、`[]Heading` と `slugByLine` を作成する。
6. 各ページの脚注定義を収集し、`RenderContext` を初期化する。
7. `buildTOC(headings)` で `PageData.TocHTML` を作成する。
8. `convert(lines, headings, ctx)` で `ConvertResult` を作成する。
9. `injectChapterNavigation(result.HTML, headings)` を適用し、`PageData.BodyHTML` を確定する。
10. サイト全体の検索インデックスを `encoding/json` で生成し、`SiteData.SearchIndex` を確定する。
11. `assembleSite(siteData, theme)` で `index.html`、ページ HTML、`assets/style.css`、`assets/app.js`、`assets/search-index.json` を生成する。
12. 出力サイトディレクトリのファイル数と合計サイズを取得し、stdout に完了行と `[REPORT]` 行を出力する。

途中で失敗した場合は、失敗段階以降を実行しない。一時ファイルが存在する場合は削除してから終了する。

**終了コード：**

| 終了コード | 条件 | 後続処理 |
|------------|------|----------|
| `0` | 静的 Web サイト生成に成功し、`[REPORT]` 行を出力した。 | `components/runner.go` は成功として扱う。 |
| `1` | 出力ディレクトリ作成、HTML / CSS / JavaScript / search index 書き込み、テンプレート合成など処理中の一般エラー。 | `components/runner.go` はビルド失敗として扱い、SHA を更新しない。 |
| `2` | CLI 引数不正、入力ファイル不存在、入力 UTF-8 不正。 | `components/runner.go` は設定または入力エラーとして扱い、SHA を更新しない。 |

終了コード `0` の場合、stdout には必ず `Collecting Markdown...`、`Converting MD...`、`Building site...`、`Writing assets...`、`Done → ...`、`[REPORT] ...` をこの順序で出力する。警告がある場合は `[REPORT]` の直前に `[WARN] ...` を 1 件 1 行で出力する。

終了コード `1` または `2` の場合、stderr に原因を 1 行以上出力し、`[REPORT]` 行は出力しない。途中まで作成した出力サイトは公開用パスへ残してはならず、一時ディレクトリを削除して終了する。

**標準出力：**
```
Collecting Markdown...
Converting MD...
Building site...
Writing assets...
Done → /opt/adlaire-builder/dist/site  (pages=12 files=15 bytes=1713731)
[REPORT] pages=12 headings=342 tables=128 code_blocks=64 warnings=3 size_warn=false broken_links=1 heading_skips=0 reading_time=87 theme=adlaire-default
```

**変換レポート行（`[REPORT]` プレフィックス）：**
ビルド完了直後に 1 行で出力する。フィールドはスペース区切りの `key=value` 形式で固定順。

| フィールド | 内容 |
|---|---|
| `pages` | 出力した HTML ページ数 |
| `headings` | 出力サイト内の見出し要素（`h1`〜`h6`）の総数 |
| `tables` | 変換したテーブルの総数 |
| `code_blocks` | 変換したコードブロックの総数 |
| `warnings` | ビルド中に発生した警告件数 |
| `size_warn` | 出力サイト合計サイズが `OUTPUT_SIZE_WARN_MB` を超えた場合 `true`、それ以外 `false`（`OUTPUT_SIZE_WARN_MB = 0` の場合は常に `false`） |
| `broken_links` | 参照先スラグが存在しない内部アンカーリンク（`[label](#anchor)`）の件数 |
| `heading_skips` | 見出しレベルが 2 段以上の降順スキップとなった件数 |
| `reading_time` | 推計読了時間（分、切り上げ）。200文字/分で算出 |
| `theme` | 使用した theme 名。初期仕様では `adlaire-default` |
| `build_id` | `--build-id` または環境変数 `ADLAIRE_BUILD_ID` の値。未指定時は空文字 |
| `commit_sha` | `--commit-sha` または環境変数 `ADLAIRE_COMMIT_SHA` の値。未指定時は空文字 |
| `build_at` | `--build-at` または環境変数 `ADLAIRE_BUILD_AT` の値。未指定時は空文字 |

固定順は以下とし、未使用フィールドの省略は禁止する。

```text
pages headings tables code_blocks warnings size_warn broken_links heading_skips reading_time theme build_id commit_sha build_at
```

`warnings` は出力した `[WARN]` 行数と一致しなければならない。`reading_time` は `ConvertResult.ReadingTimeMinutes`、`broken_links` は `ConvertResult.BrokenLinks`、`heading_skips` は `ConvertResult.HeadingSkips` を使用する。

警告が発生した場合、`[REPORT]` 行の直前に `[WARN] {メッセージ}` 形式で 1 件ずつ出力する。

**`components/runner.go` による取り込み：**
Go 版 CI ランナーでは、`components/runner.go` が `pipeline.sh` の標準出力から `[REPORT]` 行と `[WARN]` 行を抽出し、パースした結果を `.build_logs/{id}.json` のビルドログエントリに追記する。

```json
{
  "build_id": "b20260915100000",
  "status": "success",
  "report": {
    "headings": 342,
    "tables_count": 128,
    "code_blocks_count": 64,
    "warnings_count": 3,
    "size_warn": false,
    "broken_links": 1,
    "heading_skips": 0,
    "reading_time": 87,
    "theme": "adlaire-default",
    "build_id": "b20260915100000",
    "commit_sha": "abc1234",
    "build_at": "2026-09-15T10:00:00Z"
  }
}
```

> **フィールド名の対応：** stdout の `[REPORT]` 行は `tables=` / `code_blocks=` の短縮キーを使用するが、`.build_logs/{id}.json` への保存時および `GET /api/output-meta` レスポンスでは `tables_count` / `code_blocks_count` に変換する（→ §22）。

**再実行時の注意：** スラグ重複カウンタ、脚注参照順、脚注定義は `adlaire-ci-build` の 1 実行内で初期化する。通常の `/usr/local/bin/adlaire-ci-build` 実行では複数回実行しても出力は同一になる。Go 版では変換状態をパッケージグローバル変数として共有せず、変換処理ごとに専用の状態構造体を生成する。

---

## 8a. `components/builder.go` 受け入れ fixture

Go 版 `components/builder.go` の初期実装は、本節の fixture をすべて満たすまで完了として扱わない。fixture ファイルは実装 PR で `testdata/builder/` 配下へ追加する。仕様 PR では fixture の期待値を本節で固定する。

### Fixture A: 単一 Markdown 入力

**入力ファイル：** `testdata/builder/single/source.md`

~~~markdown
# Title

Intro paragraph with [self](#title).

## Install

```bash
echo hello
```

- [x] done
- [ ] todo

| Name | Value |
| ---- | ----- |
| A | 1 |

[^n]: note body

See footnote[^n].
~~~

**実行：**

```bash
adlaire-ci-build --src testdata/builder/single/source.md --out /tmp/adlaire-ci-fixture-single --title "Fixture Site"
```

**期待結果：**

- 終了コード `0`。
- `/tmp/adlaire-ci-fixture-single/index.html`、`assets/style.css`、`assets/app.js`、`assets/search-index.json` が存在する。
- `pages/` は存在しない。
- `index.html` に `<h1 id="title" class="mh h1">Title<button class="hn-link" data-href="#title" aria-label="リンクをコピー">¶</button></h1>` を含む。
- `index.html` に `<a href="#title">self</a>` を含み、`BROKEN_LINK` 警告を出さない。
- `index.html` に `<div class="cb-wrap" data-lang="bash">` と `<span class="cl">bash</span>` を含む。
- `index.html` に `<li class="ml-task"><input type="checkbox" disabled checked>done</li>` と `<li class="ml-task"><input type="checkbox" disabled>todo</li>` を含む。
- `[REPORT]` は `pages=1`、`theme=adlaire-default` を含む。

### Fixture B: ディレクトリ Markdown 入力

**入力ファイル：**

```text
testdata/builder/site/docs/intro.md
testdata/builder/site/docs/guide/setup.md
testdata/builder/site/docs/guide/setup_copy.md
```

`intro.md`:

```markdown
# Intro

Go to [setup](guide/setup.md#setup).
```

`guide/setup.md`:

```markdown
# Setup

Body.
```

`guide/setup_copy.md`:

```markdown
# Setup

Second.
```

**実行：**

```bash
adlaire-ci-build --src testdata/builder/site/docs --out /tmp/adlaire-ci-fixture-site --title "Docs"
```

**期待結果：**

- 終了コード `0`。
- `/tmp/adlaire-ci-fixture-site/index.html` がサイト目次ページである。
- document ページは `pages/guide-setup.html`、`pages/guide-setup-copy.html`、`pages/intro.html` として出力される。
- 入力順は `guide/setup.md`、`guide/setup_copy.md`、`intro.md` の辞書順とし、サイト目次の表示順も同一とする。
- `intro.html` 内の `guide/setup.md#setup` は、同じ `pages/` ディレクトリ内のページ間リンクとして `guide-setup.html#setup` に変換する。`guide/setup.html#setup`、`pages/guide-setup.html#setup`、元の `guide/setup.md#setup` のまま出力してはならない。
- 2 つの `# Setup` 見出しはページごとの slug 空間でそれぞれ `setup` とする。ページをまたいだ見出し slug に `-2` を付けてはならない。
- `assets/search-index.json` は `index.html#site-index`、`pages/guide-setup.html#setup`、`pages/guide-setup-copy.html#setup`、`pages/intro.html#intro` の entry を含む。

### Fixture C: 異常系

| 実行 | 終了コード | stderr |
|------|------------|--------|
| `adlaire-ci-build --theme unknown` | `2` | `unknown theme: unknown` |
| `adlaire-ci-build --src /path/not-found.md` | `2` | `source not found: /path/not-found.md` |
| `adlaire-ci-build --src testdata/builder/empty-dir` | `2` | `no markdown files found: testdata/builder/empty-dir` |
| `adlaire-ci-build --title ""` | `2` | `title must not be empty` |

異常系 fixture では `[REPORT]` を stdout へ出力してはならない。`--out` に既存の正常出力がある場合でも、異常系実行で既存出力を変更してはならない。

### Fixture D: 冪等性

同一入力、同一 CLI 引数で 2 回連続実行した場合、`GeneratedAtUTC` を含む meta 行を除き、全出力ファイルの内容が一致しなければならない。比較対象から除外できるのは、HTML 内の `name="adlaire-generated-at"` meta と footer の生成時刻表示だけとする。検索 index、ページ HTML の本文、CSS、JS、REPORT の数値は一致必須とする。

---

## 9. 既知の制限

| 制限 | 詳細 |
|------|------|
| インデント付き閉じフェンス | CommonMark では `  ``` ` のようなインデント付き閉じフェンスが有効だが、本実装は `raw.strip() == fence_marker` で完全一致を要求するため非対応。インデント付き閉じフェンスはフェンス内の行として取り込まれる |
| 生 HTML のパススルー | Markdown 中の生 HTML（`<div>`、`<span>` 等）は `esc()` でエスケープされテキストとして出力される。HTML をそのまま通過させる機能はない |

---

### — CI ランナー —

## 10. CI ランナー 要件

| 項目 | 内容 |
|------|------|
| Go バージョン | Go `1.22` 以上。 |
| 外部依存 | **なし** — Go 標準ライブラリ（`net/http`、`encoding/json`、`os`、`os/exec`、`log/slog` 等）を使用する |
| 対象 OS | Linux（systemd 対応環境） |
| ネットワーク | サーバーから `api.github.com` への HTTPS 送信のみ |

---

## 10a. CI ランナー 実装対象

本節は、Go 版 `components/runner.go` として実装する CI ランナー機能を定義する。

`components/runner.go` は `adlaire-ci-runner` バイナリとして実行する。起動形式は systemd timer から呼び出される oneshot 実行とし、1 回の起動で対象ブランチ設定を読み込み、変更検出、ビルド起動、ログ保存、通知、転送、後処理を完了して終了する。

実装時は、対象項目ごとに §0c の実装前確認項目を満たしていることを確認する。未充足の項目が 1 つでもある場合は、実装を開始せず、先に本ファイルの該当節を改訂する。

| 項目 | 関連節 | 実装内容 |
|------|--------|------------|
| JSON 形式の SHA キャッシュ | §11〜§13 | 各ターゲットの `sha_file` を `{"sha": "..."}` JSON 形式で読み書きする。 |
| `BRANCH_TARGETS` | §12〜§13 | 複数ブランチ、複数 target file、複数出力先を 1 つの設定リストとして処理する。 |
| GitHub API リトライ | §12〜§13 | `API_RETRY_MAX`、`API_RETRY_BASE_SECONDS`、指数バックオフ、レート制限待機を実装する。 |
| ビルドロック | §11〜§13 | `.build_lock` による多重起動防止を実装する。 |
| ビルドクールダウン | §12〜§13 | `BUILD_COOLDOWN_SECONDS` による起動抑制を実装する。 |
| 強制再ビルド間隔 | §12〜§13 | `FORCE_BUILD_INTERVAL` による変更なし時の定期強制ビルドを実装する。 |
| コミット情報記録 | §13 | ビルドトリガー commit の SHA、message、author、date をビルドログへ記録する。 |
| 事前チェック | §13 | ディスク空き容量、`adlaire-ci-build` 実行可否、Go 版ビルドバイナリ配置を確認する。 |
| Webhook 通知 | §13 | `.notify_config` 読み込み、成功/失敗/転送失敗/週次サマリー通知、`.notify_pending` 再送を実装する。 |
| ビルドログ保存 | §11〜§15 | `.build_logs/{id}.json` へ stdout/stderr、変換レポート、所要時間を保存する。 |
| SSH 転送 | §14a | SHA256 差分検出、stdin パイプ転送、転送後整合性検証、ペンディングキューを実装する。 |
| スナップショット | §14b | `.snapshots/` への成果物保存、世代管理、ロールバック連携を実装する。 |
| サーキットブレーカー | §13 | `.build_circuit_state` による連続失敗停止を実装する。 |
| ログ世代管理 | §13 | `LOG_KEEP_N` による `.build_logs/` 削除を実装する。 |
| 出力サイズ警告 | §12〜§13 | `OUTPUT_SIZE_WARN_MB` による WARN ログと `size_warn` 記録を実装する。 |

### 初期実装対象外の連携範囲

§10〜§20 には、`components/runner.go` 単体の責務ではなく管理 API、標準管理ツール、追加の運用機能と結合して成立する項目が含まれる。これらは、API・SDK・UI の対象節に、呼び出し元、呼び出し先、状態ファイル、失敗時応答、検証条件が定義されるまで `components/runner.go` 単体で実装しない。

| 項目 | 理由 |
|------|------|
| API 経由の動的ブランチ設定 | `components/runner.go` 単体では設定 API を持たないため、`components/api.go` 実装と合わせて扱う。 |
| API 経由のロールバック | `POST /api/history/{id}/rollback` は `components/api.go` のエンドポイント実装が前提となる。 |
| 管理画面からのスケジュール操作 | systemd timer の変更 API と標準管理ツール UI が前提となる。 |

---

## 11. CI ランナー ファイル構成

本節のファイル構成は、Go 版 CI ランナーで使用するファイル、管理 API / SDK / UI 側のファイル、出力先、systemd ファイルを分離して示す。

### Go 版 CI ランナーで使用するファイル

| パス | 用途 |
|------|------|
| `/usr/local/bin/adlaire-ci-runner` | `components/runner.go` から生成する CI ランナーバイナリ。 |
| `/usr/local/bin/adlaire-ci-build` | `components/builder.go` から生成する Markdown → 静的 Web サイトビルドバイナリ。 |
| `/opt/adlaire-builder/.github_token` | GitHub PAT。Go 版 `components/runner.go` が読み込む。 |
| `/opt/adlaire-builder/.last_sha` | 前回取得した blob SHA。JSON 形式で保存する。 |
| `/opt/adlaire-builder/repo/docs/` | GitHub Blobs API から取得した Markdown の書き出し先。単一 Markdown の場合も本ディレクトリ内へ保存する。 |
| `/opt/adlaire-builder/repo/.ci/pipeline.sh` | `components/runner.go` が `bash` で起動するビルド手順。 |

```
/opt/adlaire-builder/
├── .github_token
├── .last_sha
└── repo/
    ├── docs/
    └── .ci/
        └── pipeline.sh
```

### CI ランナー状態ファイル

以下は §10a の実装対象に対応するファイルである。Go 版 `components/runner.go` は本仕様に従って作成・読み書きする。

| パス | 用途 |
|------|------|
| `/opt/adlaire-builder/.build_history` | ビルド履歴。 |
| `/opt/adlaire-builder/.notify_config` | Webhook 通知設定。runner は起動時整合性チェックと送信読み込みを行い、管理 API は設定更新を行う。 |
| `/opt/adlaire-builder/.notify_log` | Webhook 送信履歴。 |
| `/opt/adlaire-builder/.notify_pending` | Webhook 通知失敗時の再送キュー。 |
| `/opt/adlaire-builder/.pending_transfers` | SSH 転送失敗時の再送キュー。 |
| `/opt/adlaire-builder/.build_lock` | 実行中ビルドの PID ロック。 |
| `/opt/adlaire-builder/.branch_config` | ブランチターゲット設定。永続 JSON key は `branch_targets` とする。 |
| `/opt/adlaire-builder/.build_state` | ビルド実行状態、週次サマリー送信日等。 |
| `/opt/adlaire-builder/.build_status.json` | runner 現在状態と直近結果の要約。API / UI / MCP の read-only 参照元。 |
| `/opt/adlaire-builder/.build_circuit_state` | サーキットブレーカー状態。 |
| `/opt/adlaire-builder/.build_logs/` | ビルドごとの個別ログ。 |
| `/opt/adlaire-builder/.snapshots/` | ビルド成果物スナップショット。 |

### 管理 API / SDK / UI 側ファイル

以下は `components/api.go`、`admin/adlaire-ci-sdk.js`、`admin/index.html` の仕様に属する。CI ランナー拡張と連携するものを含むが、Go 版 `components/runner.go` 単体の実装対象範囲には含めない。

```
/opt/adlaire-builder/
├── .admin_credentials   # 認証情報ファイル（JSON、パーミッション 600）
├── .server_config       # サーバー設定（JSON）
├── .access_log          # ログイン履歴（JSON）
├── .repo_config         # リポジトリ監視設定（JSON）
├── .config_log          # 設定変更履歴（JSON）
├── .access_control      # IP アクセス制限設定（JSON）
├── .hooks               # ビルド前後フック設定（JSON）
├── .maintenance         # メンテナンスモード状態（JSON）
├── .api_tokens          # API トークン管理（JSON、トークン本体はハッシュのみ保存）
├── .alert_rules         # カスタムアラートルール（JSON）
├── .tag_rules           # 自動タグ付けルール（JSON）
├── .pipeline_config     # パイプライン設定（JSON）
├── .notes               # 運用ノート（テキスト）
├── .smtp_config         # SMTP 設定（JSON、パスワード除く）
├── .smtp_secret         # SMTP パスワード（プレーンテキスト、パーミッション 600）
├── .dashboard_layout    # ダッシュボードウィジェットレイアウト（JSON）
├── .webhook_events.json # Webhook 受信イベントログ（JSON Lines 形式、1行1イベント）
└── admin/
    ├── index.html           # 管理画面（単一ファイル完結）
    └── adlaire-ci-sdk.js    # JavaScript SDK（管理画面に同梱）
```

管理 API サーバーの実行ファイルは `/usr/local/bin/adlaire-ci-api` とし、`components/api.go` から生成する。

### 出力先・配信先

```
/opt/adlaire-builder/dist/
└── site/                # 静的 Web サイト出力先

/var/www/html/           # SSH 転送先
```

### systemd

```
/etc/systemd/system/
├── adlaire-ci.service      # systemd ユニット（oneshot）
├── adlaire-ci.timer        # systemd タイマー（定期実行）
└── adlaire-ci-api.service  # 管理 API サーバー（常駐）
```

### リポジトリ側

```
<repo>/
├── docs   # ソース Markdown（GitHub 上のマスター）
└── .ci/
    └── pipeline.sh      # ビルド手順定義（実行権限付き）
```

---

## 12. 設定値（`components/runner.go`）

Go 版 `components/runner.go` は本節の設定値を正とする。設定値は Go 構造体の既定値、設定ファイル、または CLI 引数で与える。どの入力経路を採用する場合でも、内部表現は本節のキー名・型・既定値に従う。

**関連型：**

```go
type RunnerConfig struct {
    StateDir                    string
    PendingFile                 string
    APIRetryMax                 int
    APIRetryBaseSeconds         int
    BuildCooldownSeconds        int
    HistoryKeepN                int
    ForceBuildIntervalHours     int
    LogKeepN                    int
    APICircuitBreakerThreshold  int
    OutputSizeWarnMB            int
    WeeklySummaryEnabled        bool
    WeeklySummaryDay            int
    WeeklySummaryHour           int
    BranchTargets               []BranchTarget
}

type BranchTarget struct {
    Branch        string
    TargetFile    string
    SHAFile       string
    Src           string
    Out           string
    DeployTargets []DeployTarget
}

type DeployTarget struct {
    Host    string
    User    string
    DestDir string
}
```

**設定値検証：**

| 項目 | 条件 | 不正時 |
|------|------|--------|
| `StateDir` | 絶対パス。空文字不可。 | 終了コード `2` |
| `PendingFile` | 絶対パス。空文字不可。 | 終了コード `2` |
| `APIRetryMax` | 0 以上。 | 終了コード `2` |
| `APIRetryBaseSeconds` | 0 以上。 | 終了コード `2` |
| `BuildCooldownSeconds` | 0 以上。 | 終了コード `2` |
| `HistoryKeepN` | 0 以上。 | 終了コード `2` |
| `ForceBuildIntervalHours` | 0 以上。 | 終了コード `2` |
| `LogKeepN` | 0 以上。 | 終了コード `2` |
| `APICircuitBreakerThreshold` | 0 以上。 | 終了コード `2` |
| `OutputSizeWarnMB` | 0 以上。 | 終了コード `2` |
| `WeeklySummaryDay` | 0〜6。 | 終了コード `2` |
| `WeeklySummaryHour` | 0〜23。 | 終了コード `2` |
| `BranchTargets` | 1 件以上。 | 終了コード `2` |
| `BranchTarget.Branch` | 空文字不可。`..`、`~`、制御文字禁止。 | 終了コード `2` |
| `BranchTarget.TargetFile` | 相対パス。絶対パス、`..`、先頭 `/` 禁止。 | 終了コード `2` |
| `BranchTarget.SHAFile` / `Src` / `Out` | 絶対パス。空文字不可。 | 終了コード `2` |
| `DeployTarget.Host` / `User` / `DestDir` | 空文字不可。`DestDir` は絶対パス。 | 当該 deploy target を無効として ERROR ログ。全 target 無効なら終了コード `2` |

`.server_config`、`.branch_config`、CLI 引数から読み込んだ値は `RunnerConfig` に正規化してから使用する。正規化後の `RunnerConfig` にないキーを処理フローで直接参照してはならない。

**設定入力の優先順位：**

1. CLI 引数
2. `.server_config` / `.branch_config` など状態ファイルの保存値
3. 本節の既定値

同一キーが複数の入力経路に存在する場合は、上位の値だけを採用する。採用しなかった値を混合してはならない。未知キーは WARN ログ `CONFIG_UNKNOWN_KEY: key={key}` を出して無視する。

**runner CLI 引数仕様：**

| 引数 | 必須 | 既定値 | 説明 |
|------|------|--------|------|
| `--state-dir <path>` | 任意 | `/opt/adlaire-builder` | 状態ファイル、repo、dist、admin の基準ディレクトリ。相対パスは禁止。 |
| `--once` | 任意 | `true` | 1 回だけ実行して終了する。Go 版 runner は oneshot 固定のため、指定してもしなくても同じ挙動とする。 |
| `--dry-run` | 任意 | `false` | 設定、状態、GitHub target、SHA 差分、起動可否だけを検証し、ビルド、deploy、通知、状態ファイル更新を行わず終了する。 |
| `--version` | 任意 | なし | バイナリ名、仕様名、Go build 情報を 1 行で標準出力へ表示して終了する。 |
| `--help` | 任意 | なし | 引数一覧を標準出力へ表示して終了する。 |

未知引数、値欠落、相対 `--state-dir` は終了コード `2` とし、ビルド処理を開始しない。

**runner CLI パース固定仕様：**

- 引数は `flag` package 互換の `--name value` と `--name=value` の両方を許可する。
- 短縮オプション（例：`-s`、`-o`）は禁止する。指定された場合は未知の引数として扱う。
- 同一引数が複数回指定された場合は最後の値を採用する。ただし `--once` は指定有無にかかわらず `true` として扱う。`--dry-run` は 1 回以上指定されれば `true` とする。
- `--help` と `--version` は他の引数より優先し、`.github_token` 読み込み、lock 作成、状態ファイル読み込み、GitHub API 呼び出しを行わない。
- stderr のエラー行は末尾に改行 1 つを付ける。複数エラーをまとめて出力せず、最初に検出したエラー 1 件で終了する。

**固定出力：**

| 条件 | stdout |
|------|--------|
| `--help` | `Usage: adlaire-ci-runner [--state-dir path] [--once] [--dry-run] [--version] [--help]` |
| `--version` | `adlaire-ci-runner ADLAIRE_CI_SPEC go=<runtime.Version()>` |

**CLI 異常系：**

| 条件 | 終了コード | stderr |
|------|------------|--------|
| 未知の引数 | `2` | `unknown option: <name>` |
| `--state-dir` 値欠落 | `2` | `missing value: --state-dir` |
| `--state-dir` が空文字 | `2` | `state directory must not be empty` |
| `--state-dir` が相対パス | `2` | `state directory must be absolute: <path>` |
| `--state-dir` が存在しない | `2` | `state directory not found: <path>` |
| `--state-dir` がディレクトリではない | `2` | `state path is not directory: <path>` |

以下の `BRANCH_TARGETS`、`PENDING_FILE`、`API_RETRY_MAX`、`BUILD_COOLDOWN_SECONDS`、`HISTORY_KEEP_N`、`FORCE_BUILD_INTERVAL`、`LOG_KEEP_N`、`API_CIRCUIT_BREAKER_THRESHOLD`、`OUTPUT_SIZE_WARN_MB`、`WEEKLY_SUMMARY_*` を Go 版 runner の標準設定とする。

```text
PENDING_FILE           = "/opt/adlaire-builder/.pending_transfers"   # SSH 転送ペンディングキュー（JSON）
API_RETRY_MAX          = 5    # GitHub API 失敗時の最大再試行回数（指数バックオフ）
API_RETRY_BASE_SECONDS = 1    # バックオフ基底秒数（1→2→4→8→16 秒。0 = リトライ無効）
BUILD_COOLDOWN_SECONDS = 60   # 前回ビルド完了から次ビルドまでの最小間隔（秒。0 = 無効）→ §13
HISTORY_KEEP_N         = 10   # スナップショット保持世代数（0 = 無制限）→ §14b
FORCE_BUILD_INTERVAL   = 0    # 強制再ビルド間隔（時間。0 = 無効）→ §13
LOG_KEEP_N             = 50   # ビルドログ保持件数（0 = 無制限）→ §13
API_CIRCUIT_BREAKER_THRESHOLD = 3    # 全ブランチ連続失敗の許容周回数（0 = 無効）→ §13
OUTPUT_SIZE_WARN_MB           = 5    # 出力サイト合計サイズ警告閾値（MB。0 = 無効）→ §13・§8
WEEKLY_SUMMARY_ENABLED        = true # 週次サマリー Webhook の有効/無効 → §13
WEEKLY_SUMMARY_DAY            = 0    # 送信曜日（0=月曜〜6=日曜） → §13
WEEKLY_SUMMARY_HOUR           = 9    # 送信時刻（0〜23、ローカル時刻） → §13

# ブランチターゲット設定（→ §14a）
# 複数エントリを定義した場合はリスト順に順次処理する（並列処理は対象外）
BRANCH_TARGETS = [
    {
        "branch":      "main",                                        # 監視対象ブランチ
        "target_file": "docs",                                        # 監視対象 Markdown ファイルまたはディレクトリ
        "sha_file":    "/opt/adlaire-builder/.last_sha",             # blob SHA キャッシュ（JSON 形式: {"sha": "..."}）
        "src":         "/opt/adlaire-builder/repo/docs",              # Blobs API 書き出し先
        "out":         "/opt/adlaire-builder/dist/site",              # ビルド成果物ディレクトリ
        "deploy_targets": [
            {
                "host":     "<配信サーバーIP>",                       # SSH 転送先ホスト
                "user":     "deploy",                                 # SSH 接続ユーザー
                "dest_dir": "/var/www/html",                          # 転送先ディレクトリ
            }
        ],
    }
]
```

`BRANCH_TARGETS` が空の場合、`components/runner.go` は ERROR ログを出力し、ビルドを実行せず終了コード `2` で終了する。

**runner 終了コード：**

| 終了コード | 条件 |
|------------|------|
| `0` | 起動、ロック確認、対象処理が完了した。変更なし、クールダウン、既存ロックによるスキップも正常終了に含める。 |
| `1` | 1 件以上のビルドまたは転送が失敗したが、runner 自体は最後まで処理できた。 |
| `2` | 設定不正、必須ファイル不足、CLI 引数不正。 |
| `3` | GitHub API など外部サービスへの全再試行が失敗し、全ターゲットが処理不能。 |
| `4` | `.build_lock` 作成に失敗し、既存 PID の実行中確認もできない。 |

systemd timer からの再実行を妨げないため、終了コード `1` と `3` でもロック削除、ログ保存、通知キュー保存を試行してから終了する。

**atomic write 共通契約：**

runner が JSON object、JSON array、SHA cache、`.build_state`、`.build_circuit_state`、`.notify_pending`、`.pending_transfers`、`.server_config`、`.branch_config` を更新する場合は、以下の順序で atomic write を行う。

1. 対象ファイルと同じディレクトリに `.{basename}.tmp.{pid}` を作成する。
2. JSON は末尾改行付き UTF-8 として書き込む。
3. `file.Sync()` を実行してから close する。
4. `os.Rename(tmp, target)` で置換する。
5. 親ディレクトリを open できる場合は directory sync を実行する。directory sync が `EINVAL` 等で未対応の場合は WARN ログ `DIR_SYNC_UNSUPPORTED: path={dir}` を出し、処理は成功扱いとする。

atomic write 失敗時は対象ファイルを更新済みとして扱わない。tmp ファイルが残った場合は削除を試行し、削除失敗時は WARN ログ `TMP_CLEANUP_FAILED: path={tmp}` を出す。

**lock ファイル契約：**

`.build_lock` は UTF-8 text で、内容は `pid={pid}\nstarted_at={UTC_ISO8601}\n` とする。

| 状態 | 処理 |
|------|------|
| lock なし | `os.OpenFile(path, O_CREATE|O_EXCL|O_WRONLY, 0600)` で作成する。 |
| lock あり・PID 実行中 | INFO ログ `BUILD_SKIP: already running (PID {pid})` を出し終了コード `0`。 |
| lock あり・PID 不在 | WARN ログ `STALE_LOCK: pid={pid}` を出し、lock を削除してから再作成する。 |
| lock あり・PID 解析不能 | ERROR ログ `LOCK_CORRUPT: path={path}` を出し終了コード `4`。 |
| lock 作成失敗 | ERROR ログ `LOCK_CREATE_FAILED: {reason}` を出し終了コード `4`。 |

PID 実行中判定は Linux の `/proc/{pid}` 存在確認で行う。`/proc` を読めない場合は PID 実行中確認不能として終了コード `4` とする。

**runner 状態ファイル schema（実装固定）：**

`.last_sha` / `BranchTarget.SHAFile`:

```json
{"sha":""}
```

| キー | 型 | 必須 | 許容値 |
|------|----|------|--------|
| `sha` | string | 必須 | 空文字または Git blob SHA。空文字は初回実行扱い。 |

`.branch_config`:

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
        {"host": "192.0.2.1", "user": "deploy", "dest_dir": "/var/www/html"}
      ]
    }
  ]
}
```

`.branch_config` が存在しない場合は §12 の `BRANCH_TARGETS` 既定値を使用する。存在する場合は `branch_targets` を必須キーとし、未知キーは `CONFIG_UNKNOWN_KEY: key={key}` を出して無視する。`branch_targets` が空配列の場合は終了コード `2` とする。

`.build_state`:

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

| キー | 型 | 必須 | 更新責務 |
|------|----|------|----------|
| `running` | boolean | 必須 | lock 取得後 `true`、終了直前 `false`。 |
| `current_build_id` | string/null | 必須 | 実行中 build id。終了後は `null`。 |
| `queued` | array | 必須 | API 側の手動キュー用。runner は読み取りのみとし、初期実装では変更しない。 |
| `last_started_at` | string/null | 必須 | UTC ISO 8601。 |
| `last_finished_at` | string/null | 必須 | UTC ISO 8601。 |
| `weekly_summary_last_sent_at` | string/null | 必須 | UTC ISO 8601。 |
| `weekly_summary_sent_date` | string/null | 必須 | `YYYY-MM-DD`。二重送信防止に使用する。 |

`.build_circuit_state`:

```json
{
  "open": false,
  "consecutive_failures": 0,
  "opened_at": null,
  "last_failure_at": null,
  "last_error": null
}
```

`.notify_pending` entry:

```json
{
  "event": "success",
  "url": "https://example.com/hook",
  "payload": {},
  "queued_at": "2026-09-16T00:00:00Z",
  "retry_count": 1,
  "last_error": "connection refused"
}
```

`.pending_transfers` entry は §14a の形式を正とする。JSON array 内の entry は投入順を保持し、再試行も投入順で処理する。重複統合は `out`、`host`、`user`、`dest_dir` の 4 項目完全一致で判定する。

**設定ファイル起動時整合性チェック：**

本機能の目的は、runner 起動時に状態ファイルの破損、型不一致、必須 key 不足、権限不備を検出し、ビルド処理開始前に復旧または停止することである。対象コンポーネントは `components/runner.go` のみとし、管理 API の HTTP endpoint、SDK、UI は本機能の実装対象に含めない。

対象ファイルは次の 6 件に固定する。実装者判断で対象ファイルを追加または除外してはならない。

| 順序 | ファイル | 必須性 | 正常時の扱い | 不在時の扱い |
|------|----------|--------|--------------|--------------|
| 1 | `.branch_config` | 任意 | JSON object として検証し、`branch_targets` を `RunnerConfig.BranchTargets` へ正規化する。 | ファイルを作成せず、§12 の `BRANCH_TARGETS` 既定値を使用する。 |
| 2 | `.notify_config` | 任意 | JSON object として検証し、通知送信時の設定として使用する。 | 初期値を atomic write で作成する。 |
| 3 | `.build_state` | 必須 | JSON object として検証し、`queued`、`last_finished_at`、週次サマリー状態を保持する。 | 初期値を atomic write で作成する。 |
| 4 | `.build_circuit_state` | 必須 | JSON object として検証し、circuit breaker 状態を保持する。 | 初期値を atomic write で作成する。 |
| 5 | `.pending_transfers` | 必須 | JSON array として検証し、entry の投入順を保持する。 | `[]` を atomic write で作成する。 |
| 6 | `.notify_pending` | 必須 | JSON array として検証し、entry の投入順を保持する。 | `[]` を atomic write で作成する。 |

実行順序は固定とする。

1. CLI 引数を検証し、`--help` または `--version` の場合は本チェックを実行しない。
2. `StateDir` が絶対パスかつ既存ディレクトリであることを確認する。
3. `.build_lock` を取得する。実行中 PID がある場合は本チェックを実行せず終了コード `0` で終了する。
4. 本チェックを上表の順序で実行する。
5. 本チェックが復旧可能な問題のみで完了した場合は、`.github_token` 読み込み、pending retry、cooldown、target 処理へ進む。
6. 本チェックが停止条件に該当した場合は、`.build_state.running` を `true` にせず、`.build_logs/{id}.json` と `.build_history` を作成せず、`.build_lock` を削除して終了する。

判定分類は次のとおり固定する。

| 分類 | 条件 | 処理 |
|------|------|------|
| `missing_optional` | `.branch_config` が存在しない。 | ファイルを作成せず、`BRANCH_TARGETS` 既定値を採用する。ログは出さない。 |
| `missing_required` | `.notify_config`、`.build_state`、`.build_circuit_state`、`.pending_transfers`、`.notify_pending` が存在しない。 | 初期値を atomic write で作成し、WARN ログ `CONFIG_INIT: path={path}` を出す。 |
| `empty_file` | ファイルサイズ 0 byte。 | `parse_error` と同じ扱い。 |
| `parse_error` | UTF-8 として読めない、または JSON parse に失敗する。 | corrupt backup へ退避し、ファイル別初期化を行う。 |
| `top_level_type_mismatch` | object 必須のファイルが array/string/null、array 必須のファイルが object/string/null。 | corrupt backup へ退避し、ファイル別初期化を行う。 |
| `required_key_missing` | 必須 key が存在しない。 | corrupt backup へ退避し、ファイル別初期化を行う。 |
| `required_key_type_mismatch` | 必須 key の型が schema と異なる。 | corrupt backup へ退避し、ファイル別初期化を行う。 |
| `invalid_value` | 値が許容範囲外。例: 負の `retry_count`、不正な ISO 8601、相対 `sha_file`。 | corrupt backup へ退避し、ファイル別初期化を行う。 |
| `unknown_key` | schema にない key が存在する。 | backup せず、未知 key を除去した正規化 JSON を atomic write し、WARN ログ `CONFIG_UNKNOWN_KEY: path={path} key={key}` を出す。 |
| `permission_error` | 読み込み、rename、chmod、親ディレクトリ sync、lock 作成のいずれかが権限エラー。 | 自動退避せず、ERROR ログ `CONFIG_PERMISSION_ERROR: path={path} op={op}` を出し、終了コード `2`。 |
| `io_error` | 権限以外の read/write/rename/sync 失敗。 | 自動退避せず、ERROR ログ `CONFIG_IO_ERROR: path={path} op={op}` を出し、終了コード `1`。 |

corrupt backup のファイル名は `{original}.corrupt.{YYYYMMDDHHMMSS}.bak` とする。timestamp は UTC、秒単位、ゼロ埋め固定とする。同一秒内に同じファイルの backup 名が衝突した場合は、2 件目以降を `{original}.corrupt.{YYYYMMDDHHMMSS}.{n}.bak` とし、`n` は `2` から始める。

ファイル別の復旧結果は次のとおり固定する。

| ファイル | 復旧時の書き込み内容 | 復旧後の処理継続 | 備考 |
|----------|----------------------|------------------|------|
| `.branch_config` | 書き込まない。破損元を backup した後、`.branch_config` は不在状態にする。 | 継続する。 | §12 の `BRANCH_TARGETS` 既定値へ fallback する。 |
| `.notify_config` | §22.0a の初期値。 | 継続する。 | このファイル自体が破損していた場合、`config_corrupt` 通知は送信しない。 |
| `.build_state` | §22.0a の初期値。 | 継続する。 | `queued` は失われるため、ERROR ログと通知対象に含める。 |
| `.build_circuit_state` | §22.0a の初期値。 | 継続する。 | circuit open 状態は解除されるため、ERROR ログと通知対象に含める。 |
| `.pending_transfers` | `[]`。 | 継続する。 | 未再送転送は失われるため、ERROR ログと通知対象に含める。 |
| `.notify_pending` | `[]`。 | 継続する。 | 未送信通知は失われるため、ERROR ログを出す。通知 queue 自体が失われるため追加 queue は行わない。 |

復旧時ログは次の 3 種に固定する。

| 条件 | ログ level | メッセージ |
|------|------------|------------|
| 初期作成 | WARN | `CONFIG_INIT: path={path}` |
| corrupt backup 成功 | ERROR | `CONFIG_CORRUPT_BACKUP: path={path} backup={backup_path} reason={reason}` |
| 正規化書き戻し | WARN | `CONFIG_NORMALIZED: path={path}` |

復旧通知は `.notify_config` の復旧と検証が完了した後に 1 回だけ送信する。送信条件は、復旧対象に `.notify_config` と `.notify_pending` 以外のファイルが 1 件以上含まれ、かつ `.notify_config.webhooks[]` のうち `enabled=true` で `on` に `"config_corrupt"` または `"*"` を含む宛先が存在する場合とする。payload は以下の JSON object に固定する。

```json
{
  "event": "config_corrupt",
  "reason": "config_corrupt",
  "recovered_at": "2026-09-16T00:00:00Z",
  "files": [
    {
      "path": "/opt/adlaire-builder/.build_state",
      "action": "backup_and_init",
      "backup": "/opt/adlaire-builder/.build_state.corrupt.20260916000000.bak",
      "reason": "parse_error"
    }
  ]
}
```

通知送信に失敗した場合は、`.notify_pending` が本チェックで破損していない場合に限り、`event="config_corrupt"` の entry を `.notify_pending` へ追記する。`.notify_pending` が本チェックで初期化された場合は、二重消失を避けるため追記せず、ERROR ログ `NOTIFY_FAILED: event=config_corrupt url={url}` のみ出す。

schema 検証では次を必須とする。

- `.branch_config` は `branch_targets` のみを永続 key とし、`branches` だけを持つファイルは `required_key_missing` として扱う。API request / response の alias は永続ファイルへ保存する前に `branch_targets` へ変換する。
- `.notify_config.webhooks[].on` は `start`、`success`、`failure`、`deploy_failure`、`weekly_summary`、`config_corrupt`、`*` のみ許可する。
- `.build_state` は `weekly_summary_sent_date` を必須 key とする。値は `null` または `YYYY-MM-DD` とする。
- `.pending_transfers[]` は §14a の pending entry schema と一致すること。1 件でも不正 entry がある場合はファイル全体を corrupt として扱う。
- `.notify_pending[]` は `event`、`url`、`payload`、`queued_at`、`retry_count`、`last_error` を必須 key とする。`payload` は JSON object、`retry_count` は 0 以上の integer とする。

本チェックは冪等でなければならない。初期化または正規化済みの状態で runner を再起動した場合、追加 backup、追加通知、追加 WARN/ERROR は発生しない。同一破損ファイルが復旧失敗後に残っている場合だけ、次回起動時に再度同じ判定を行う。

検証 fixture は以下を必須とする。

| fixture | 入力状態 | 期待結果 |
|---------|----------|----------|
| `config-startup/missing-required` | 対象必須ファイルが存在しない。 | 初期値が作成され、終了コード `0`、ビルド処理へ進む。 |
| `config-startup/corrupt-build-state` | `.build_state` が `{bad json`。 | `.build_state.corrupt.{timestamp}.bak` へ退避、初期値作成、`config_corrupt` 通知対象、終了コード `0`。 |
| `config-startup/corrupt-branch-config` | `.branch_config` が JSON array。 | backup 後 `.branch_config` は不在、既定 `BRANCH_TARGETS` 採用、終了コード `0`。 |
| `config-startup/unknown-key` | `.build_circuit_state` に未知 key がある。 | backup なしで未知 key を除去、`CONFIG_NORMALIZED`、終了コード `0`。 |
| `config-startup/invalid-pending-entry` | `.pending_transfers` に必須 key 不足 entry がある。 | backup 後 `[]` 作成、終了コード `0`。 |
| `config-startup/permission-error` | 対象ファイルが読み込み不可。 | 自動退避なし、終了コード `2`、`.build_state.running` 未変更。 |
| `config-startup/help-version-skip` | `--help` または `--version`。 | 対象ファイルを読まず、変更しない。 |

完了条件は、上記 fixture を Go test で検証し、`ADLAIRE_CI_SPEC.md` の状態表、`DOCUMENT_INDEX.md` の実装状態、PR 本文の検証結果が一致していることとする。

**build id 契約：**

runner が生成する build id は UTC 時刻ベースの `b{YYYYMMDDHHmmss}` とする。同一秒内に複数 target のビルドログが必要な場合は、2 件目以降を `b{YYYYMMDDHHmmss}-2`、`-3` とする。build id は `.build_logs/{id}.json`、`.build_history`、`.snapshots/{id}/` で同一値を使用する。

---

## 13. 処理フロー

本節の処理フローは、Go 版 `components/runner.go` の標準フローである。

**状態更新順序の規範：**

1. `.build_lock` を作成する。
2. 設定ファイル起動時整合性チェックを実行する。
3. `.build_status.json` に `status="running"`、`trigger`、`started_at`、`current_build_id` を atomic write で保存する。
4. `.build_state.running=true`、`current_build_id`、`last_started_at` を atomic write で保存する。
5. GitHub API、Blob 書き出し、事前チェック、`pipeline.sh` 実行を行う。
6. ビルド成功時のみ `sha_file` を新 SHA に更新する。
7. `.build_logs/{id}.json` を作成し、stdout/stderr、`[REPORT]`、警告、転送結果、`trigger` を保存する。
8. `.build_history` に同じ `id` の要約行を JSON Lines で追記する。
9. 転送成功後に `.snapshots/` を更新する。
10. `.build_status.json` に最終状態を atomic write で保存する。
11. `.build_state.running=false`、`last_finished_at` を保存する。
12. `.build_lock` を削除する。

途中失敗時は、失敗が発生した段階以降の成功前提更新を行わない。例えば `pipeline.sh` 失敗時は `sha_file`、snapshot、転送成功履歴を更新しない。ただし `.build_logs/{id}.json`、`.build_history`、`.build_state.running=false`、通知 pending は失敗記録として保存する。

**ターゲット結果分類：**

runner は `BRANCH_TARGETS` の各 entry について、最終的に次のいずれか 1 つの `target_status` を確定する。

| `target_status` | 条件 | runner 終了コードへの影響 |
|-----------------|------|---------------------------|
| `skipped_no_change` | SHA 一致かつ強制ビルド条件未達。 | 失敗扱いしない。 |
| `skipped_cooldown` | クールダウンで全体処理を開始しない。 | `0`。 |
| `success` | pipeline 成功、SHA 更新成功、deploy target がないか全 deploy target が成功。 | `0` 候補。 |
| `success_deploy_pending` | pipeline 成功、SHA 更新成功、1 件以上の deploy が pending。 | `1`。 |
| `failure_api` | GitHub API が全再試行失敗。 | 全 target がこれなら `3`、一部なら `1`。 |
| `failure_decode` | Blob Base64 decode または Markdown 書き出し失敗。 | `1`。 |
| `failure_precheck` | 事前チェック失敗。 | `1`。 |
| `failure_build` | `pipeline.sh` が非 0 終了または timeout。 | `1`。 |
| `failure_state_write` | SHA、ログ、履歴、状態ファイルの必須更新に失敗。 | `1`。 |

`BRANCH_TARGETS` が複数ある場合、runner は設定不正を除き、1 target の失敗で全体処理を中断しない。全 target 処理後、最も重い終了コードを採用する。重さは `4 > 3 > 2 > 1 > 0` とする。

**更新可否マトリクス：**

| 状況 | `sha_file` | `.build_logs/{id}.json` | `.build_history` | deploy | snapshot |
|------|------------|-------------------------|------------------|--------|----------|
| 変更なし | 更新しない | 作成しない | 追記しない | 実行しない | 作成しない |
| 強制ビルド成功 | 新 SHA または既存 SHA を保存 | 作成する | 追記する | 実行する | deploy 成功時のみ作成 |
| API 失敗 | 更新しない | 作成する | 追記する | 実行しない | 作成しない |
| Blob decode 失敗 | 更新しない | 作成する | 追記する | 実行しない | 作成しない |
| 事前チェック失敗 | 更新しない | 作成する | 追記する | 実行しない | 作成しない |
| pipeline 失敗 | 更新しない | 作成する | 追記する | 実行しない | 作成しない |
| pipeline 成功・deploy 成功 | 更新する | 作成する | 追記する | 実行する | 作成する |
| pipeline 成功・deploy pending | 更新する | 作成する | 追記する | pending へ投入 | 作成しない |

`sha_file` は pipeline 成功後、deploy 実行前に更新する。理由は、ビルド成果物生成が成功した時点で入力 SHA の処理は完了しており、deploy 失敗は `.pending_transfers` の責務で再試行するためである。

**runner 機能単位契約：**

`components/runner.go` は、下表の機能単位で状態を更新する。各機能単位は、Write 列にない状態ファイルを更新してはならない。

| 機能単位 | Read | Write | 成功条件 | 失敗時更新 |
|----------|------|-------|----------|------------|
| lock acquisition | `.build_lock` | `.build_lock` | lock file を `O_CREATE|O_EXCL` で作成し、PID と started_at を書き込む。 | 実行中 PID がある場合は終了コード `0`。形式不正 lock は終了コード `4`、上書き禁止。 |
| startup config integrity | `.branch_config`, `.notify_config`, `.build_state`, `.build_circuit_state`, `.pending_transfers`, `.notify_pending` | 同左、corrupt backup | 対象 JSON の parse、top-level type、必須 key、型、値範囲、未知 key 正規化が完了する。 | 復旧可能な破損は backup と初期化後に継続。権限エラーは終了コード `2`、IO エラーは終了コード `1`。 |
| status start | `.build_state`, `.build_circuit_state`, `.pending_transfers`, `.notify_pending` | `.build_status.json` | `status="running"` と当該起動の `trigger` を保存する。 | 書き込み失敗は終了コード `1`。`.build_state.running=true` へ進まない。 |
| state start | `.build_state` | `.build_state` | `running=true`、`current_build_id`、`last_started_at` を保存する。 | `.build_lock` を削除し、終了コード `1`。 |
| pending transfer retry | `.pending_transfers` | `.pending_transfers`, `.build_logs/{id}.json` | 成功 entry を削除し、失敗 entry は `retry_count` を +1 して保持する。 | queue ファイル破損時は §22.0a に従い退避して `[]` で再生成する。 |
| notify pending retry | `.notify_pending`, `.notify_config` | `.notify_pending`, `.notify_log` | HTTP 2xx の entry を削除し、失敗 entry は `retry_count` を +1 して保持する。 | `.notify_config` 不在時は送信せず entry を保持し、ERROR ログを記録する。 |
| cooldown gate | `.build_state`, `.server_config` | none | cooldown 範囲外なら target 処理へ進む。 | cooldown 中は `target_status=skipped_cooldown`、終了コード `0`。 |
| GitHub tree resolve | `.github_token`, branch target | `.build_logs/{id}.json` | target file の blob SHA を取得する。 | 全 retry 失敗は `failure_api` を記録し、SHA を更新しない。 |
| SHA decision | `sha_file`, `.server_config` | none | 変更あり、force 条件成立、または webhook/force queue payload により build 対象を確定する。 | 変更なしは `skipped_no_change`。状態ファイルを更新しない。 |
| blob materializer | GitHub blob API, branch target | `src` | Base64 decode 後、対象 Markdown を atomic write する。 | decode/write 失敗は `failure_decode`。SHA を更新しない。 |
| precheck | branch target, output path, build binary | `.build_logs/{id}.json` | disk、binary、version がすべて合格する。 | `failure_precheck` を記録し、pipeline を起動しない。 |
| pipeline executor | `src`, `.pipeline_config` | `.build_logs/{id}.json`, `.build_history` | timeout 前に exit code `0`。 | 非 0 / timeout は `failure_build`。SHA、deploy、snapshot を更新しない。 |
| report importer | pipeline stdout | `.build_logs/{id}.json` | `[REPORT]` を 1 行だけ検出し schema へ変換する。 | `[REPORT]` なしは warning として `report:null` を保存する。build 成否は pipeline exit code に従う。 |
| SHA updater | `sha_file`, new SHA | `sha_file` | `{"sha":"<new_sha>"}` を atomic write する。 | `failure_state_write`。deploy、snapshot を実行しない。 |
| deploy executor | `deploy_targets`, output site | `.pending_transfers`, `.build_logs/{id}.json` | 全 deploy target の転送と検証が成功する。 | 失敗 target を `.pending_transfers` に保存し、`target_status=success_deploy_pending`。 |
| snapshot writer | output site, `.server_config` | `.snapshots/`, `.build_logs/{id}.json` | `snapshots_keep > 0` の場合に新 snapshot を保存し、超過世代を削除する。 | snapshot 失敗は build 成功を取り消さず WARN として記録する。 |
| status finalizer | target results, `.build_state`, `.build_circuit_state`, `.pending_transfers`, `.notify_pending` | `.build_status.json` | runner 起動 1 回の最終要約を保存する。 | 書き込み失敗は ERROR ログを出し、runner 終了コードを最低 `1` にする。 |
| finalizer | target results, `.build_state` | `.build_state`, `.build_lock` | `running=false`、`last_finished_at` を保存し、lock を削除する。 | lock 削除失敗は WARN。`.build_state.running=false` 保存失敗は終了コード `1`。 |

`failure_api`、`failure_decode`、`failure_precheck`、`failure_build`、`failure_state_write` は、同じ build id の `.build_logs/{id}.json.status` では `"failure"` として保存し、詳細理由は `error` に固定文言で保存する。`target_status` は runner 内部分類および `.build_logs/{id}.json.error` の詳細判定に使用し、API response の `status` 値としては返さない。

**ビルドトリガー種別契約：**

`trigger` は runner が処理を開始した原因を表す固定文字列であり、`.build_logs/{id}.json`、`.build_history`、`.build_status.json`、API response、SDK 型、UI 表示で同じ値を使用する。実装者判断で `auto`、`force`、`scheduled` など別名を追加してはならない。

| `trigger` | 発生条件 | build log | history | status |
|-----------|----------|-----------|---------|--------|
| `polling` | systemd timer 等の通常起動で SHA 差分により build する。 | 記録する | 記録する | 記録する |
| `force_interval` | SHA 差分なしだが `FORCE_BUILD_INTERVAL` 超過により build する。 | 記録する | 記録する | 記録する |
| `manual` | `POST /api/build` により `.build_state.queued` へ投入された entry を処理する。 | 記録する | 記録する | 記録する |
| `webhook` | GitHub Webhook 受信により `.build_state.queued` へ投入された entry を処理する。 | 記録する | 記録する | 記録する |
| `retry_pending_transfer` | `.pending_transfers` の再送のみを実行する。 | 再送専用 log を作成する場合に記録する | 再送履歴を追記する場合に記録する | 記録する |
| `startup_config_integrity` | 設定ファイル起動時整合性チェックで破損復旧または停止が発生する。 | 作成しない | 追記しない | 記録する |
| `rollback` | `POST /api/history/{id}/rollback` により snapshot を再転送する。 | 記録する | 記録する | 記録する |

queue entry の `trigger` は `"manual"` または `"webhook"` のみ許可する。`"force"` は使用せず、強制実行 API は queue 保存時に `"manual"` と `payload.force=true` を保存する。`force_interval` は runner が `.build_state.last_finished_at` と設定値から内部判定する場合のみ使用する。

**`.build_status.json` schema：**

`.build_status.json` は runner の現在状態と直近結果を 1 ファイルで読むための要約である。API / UI / MCP は現在状態を表示する場合、`.build_status.json` を第一参照元とし、存在しない場合のみ `.build_state`、`.build_history`、`.build_lock` から後方互換の値を算出してよい。

```json
{
  "schema_version": 1,
  "updated_at": "2026-09-16T01:00:12Z",
  "status": "success",
  "running": false,
  "current_build_id": null,
  "last_build_id": "b20260916010000",
  "last_trigger": "polling",
  "last_target_status": "success",
  "last_branch": "main",
  "last_target_file": "docs",
  "last_blob_sha": "012345",
  "last_commit_sha": "abcdef",
  "last_started_at": "2026-09-16T01:00:00Z",
  "last_finished_at": "2026-09-16T01:00:12Z",
  "last_duration_seconds": 12,
  "last_error": null,
  "last_deploy_status": "success",
  "pending_transfers_count": 0,
  "notify_pending_count": 0,
  "circuit_open": false,
  "circuit_consecutive_failures": 0,
  "output_sha256": null,
  "size_warn": false
}
```

| キー | 型 | 必須 | 許容値 |
|------|----|------|--------|
| `schema_version` | integer | 必須 | `1` 固定。 |
| `updated_at` | string/null | 必須 | UTC ISO 8601。初期値のみ `null`。 |
| `status` | string | 必須 | `"none"`, `"running"`, `"success"`, `"failure"`, `"skipped_no_change"`, `"skipped_cooldown"`, `"success_deploy_pending"`, `"circuit_open"`, `"config_recovered"`, `"config_error"`, `"lock_skipped"`。`"none"` は初期値のみ。 |
| `running` | boolean | 必須 | runner が build 処理中なら `true`。 |
| `current_build_id` | string/null | 必須 | 実行中 build id。実行中でなければ `null`。 |
| `last_build_id` | string/null | 必須 | 最後に build log / history を作成した build id。未実行なら `null`。 |
| `last_trigger` | string/null | 必須 | 上記 `trigger` 有効値または `null`。 |
| `last_target_status` | string/null | 必須 | §13 の `target_status` または `null`。 |
| `last_branch` / `last_target_file` | string/null | 必須 | 最終対象。対象なしの skip / config 系では `null`。 |
| `last_blob_sha` / `last_commit_sha` | string/null | 必須 | 取得不能時は `null`。 |
| `last_started_at` / `last_finished_at` | string/null | 必須 | UTC ISO 8601 または `null`。 |
| `last_duration_seconds` | integer/null | 必須 | 0 以上または `null`。 |
| `last_error` | string/null | 必須 | 成功・通常 skip は `null`。 |
| `last_deploy_status` | string/null | 必須 | `"success"`, `"pending"`, `"failed"`, `"skipped"`, `null`。 |
| `pending_transfers_count` / `notify_pending_count` | integer | 必須 | 0 以上。 |
| `circuit_open` | boolean | 必須 | `.build_circuit_state.open` と一致する。 |
| `circuit_consecutive_failures` | integer | 必須 | 0 以上。 |
| `output_sha256` | string/null | 必須 | 直近成功成果物の manifest SHA-256 または `null`。 |
| `size_warn` | boolean | 必須 | 直近 report の size warning。 |

`.build_status.json` の更新タイミングは次に固定する。

| タイミング | `status` | `trigger` | 備考 |
|------------|----------|-----------|------|
| lock 取得後、target 処理前 | `running` | 判定済み trigger | `current_build_id` を設定する。 |
| 実行中 lock 検出 | `lock_skipped` | `polling` | `.build_state` と build log は変更しない。 |
| cooldown skip | `skipped_cooldown` | `polling` | `last_build_id` は既存値を保持する。 |
| SHA 変更なし | `skipped_no_change` | `polling` | build log / history は作成しない。 |
| force interval build 成功 | `success` | `force_interval` | build log / history と同じ build id を保存する。 |
| build 成功・deploy 成功 | `success` | 実行 trigger | `last_deploy_status="success"`。 |
| build 成功・deploy pending | `success_deploy_pending` | 実行 trigger | `pending_transfers_count` を更新する。 |
| build 失敗 | `failure` | 実行 trigger | `last_error` に固定エラー文言を保存する。 |
| circuit open skip | `circuit_open` | `polling` | GitHub API 呼び出し前に保存する。 |
| startup config 復旧あり | `config_recovered` | `startup_config_integrity` | build log / history は作成しない。復旧後に通常 build へ進む場合、次の `running` 更新で上書きしてよい。 |
| startup config 停止 | `config_error` | `startup_config_integrity` | `.build_state.running=true` へ進まない。 |
| pending transfer retry のみ | `success` または `failure` | `retry_pending_transfer` | build が発生しない場合でも pending 件数を更新する。 |

`.build_status.json` は atomic write 対象であり、書き込み失敗時に部分ファイルを残してはならない。未知 key を追加してはならない。API は GET request で `.build_status.json` を自動修復してはならない。

**queue entry 実行契約：**

`.build_state.queued` に entry がある場合、runner は通常ポーリング対象の前に queue を FIFO で 1 件だけ取り出して処理する。queue entry 処理が成功または失敗として `.build_history` に記録された場合、その entry を queue から削除する。runner 起動 1 回で複数 queue entry を連続処理してはならない。queue entry の `trigger` が `"manual"` かつ `payload.force=true` の場合は SHA 比較を行わず build を実行する。`trigger` が `"webhook"` の場合は payload の `ref` と `sha` を優先し、branch target に一致しない entry は `failure_api` として記録した後に queue から削除する。

**状態ファイル破損時の処理：**

runner が読み込む JSON object / JSON array の状態ファイルが破損している場合は、§22.0a の破損時の扱いに従う。JSON Lines は壊れた行だけを無視し、ファイル全体を破棄してはならない。破損退避ファイル名は `{original}.corrupt.{YYYYMMDDHHMMSS}.bak` とする。

```
components/runner.go 起動（systemd タイマーから呼び出し）
    │
    ├─ .github_token 読み込み（不在、空、改行除去後 1 文字未満の場合は ERROR ログ、終了コード 2）
    │
    ├─ [重複チェック] .build_lock が存在する場合
    │   ├─ ファイル内 PID が実行中 → INFO ログ（`BUILD_SKIP: already running (PID N)`）、正常終了
    │   └─ PID が存在しない（前回の異常終了） → .build_lock を削除して続行
    │
    ├─ .build_lock に自プロセスの PID と started_at を書き込み
    │   （以降、正常終了・例外終了いずれの場合も defer で .build_lock を削除）
    │
    ├─ [設定ファイル起動時整合性チェック]
    │   ├─ .branch_config / .notify_config / .build_state / .build_circuit_state / .pending_transfers / .notify_pending を固定順序で検証
    │   ├─ 復旧可能な破損 → corrupt backup、初期化または default fallback、必要時 config_corrupt 通知
    │   ├─ unknown key のみ → backup せず正規化書き戻し
    │   └─ permission error / IO error → .build_state.running=true にせず終了
    │
    ├─ [ペンディングキュー再試行] PENDING_FILE が存在する場合（→ §14a）
    │   └─ ペンディングエントリごとに SSH 転送を再試行
    │       ├─ 成功 → エントリを PENDING_FILE から削除
    │       └─ 失敗 → ERROR ログ、エントリを保持（次回起動時に再試行）
    │
    ├─ [Webhook 通知ペンディング再試行] .notify_pending が存在する場合
    │   └─ ペンディングエントリごとに Webhook 送信を再試行
    │       ├─ 成功（HTTP 2xx）→ INFO ログ、エントリを .notify_pending から削除
    │       └─ 失敗（HTTP エラー・接続失敗）→ ERROR ログ、エントリを保持（次回起動時に再試行）
    │
    ├─ [クールダウンチェック] BUILD_COOLDOWN_SECONDS > 0 の場合
    │   └─ .build_state.last_finished_at から BUILD_COOLDOWN_SECONDS 秒未満
    │       → INFO ログ（`COOLDOWN: skip, last_build N秒前`）、正常終了
    │
    ├─ BRANCH_TARGETS の各エントリを順次処理：
    │   │
    │   ├─ Step 1: Git Trees API
    │   │   GET /repos/{OWNER}/{REPO}/git/trees/{branch}?recursive=1
    │   │   → target_file の blob SHA を取得
    │   │   └─ API 失敗時：API_RETRY_MAX 回まで指数バックオフ（API_RETRY_BASE_SECONDS × 2^n 秒）で再試行
    │   │        ├─ 全試行失敗時：ERROR ログ、このエントリを failure(api_error) として記録し、SHA を更新せず次エントリへ進む
    │   │        └─ レスポンスヘッダー X-RateLimit-Remaining = 0 の場合
    │   │             → X-RateLimit-Reset（Unix 時刻）まで待機してから再試行
    │   │               INFO ログ（`RATE_LIMIT: waiting until {reset_time}`）
    │   │   [PAT 有効期限チェック] レスポンスヘッダー GitHub-Authentication-Token-Expiration が存在する場合
    │   │        ├─ 残日数 ≤ 7 日 → WARN ログ（`PAT_EXPIRY_WARN: expires_at={date} remaining_days={N}`）
    │   │        └─ 残日数 > 7 日 → 処理継続（チェックのみ）
    │   │
    │   ├─ SHA 比較（sha_file の前回 SHA と比較）
    │   │   ├─ 一致（変更なし）かつ FORCE_BUILD_INTERVAL = 0 → INFO ログ、このエントリをスキップ
    │   │   ├─ 一致（変更なし）かつ FORCE_BUILD_INTERVAL > 0 → 前回ビルドから指定時間以上経過していれば強制ビルド続行
    │   │   └─ 不一致（変更あり）→ 続行
    │   │
    │   ├─ Step 2: Git Blobs API
    │   │   GET /repos/{OWNER}/{REPO}/git/blobs/{sha}
    │   │   → Base64 デコード → src パスへ書き出し
    │   │   └─ API 失敗時：API_RETRY_MAX 回まで指数バックオフで再試行
    │   │        ├─ 全試行失敗時：ERROR ログ、このエントリを failure(api_error) として記録し、SHA を更新せず次エントリへ進む
    │   │        └─ レスポンスヘッダー GitHub-Authentication-Token-Expiration が存在する場合は Step 1 と同様に PAT 有効期限チェックを行う
    │   │
    │   ├─ [コミット情報取得] ビルドトリガーとなったコミット情報を取得し、ビルドログへ記録する
    │   │   GET /repos/{OWNER}/{REPO}/commits?path={BRANCH_TARGET.src}&sha={branch}&per_page=1
    │   │   → 先頭エントリから取得：
    │   │       commit_sha    = commit.sha
    │   │       commit_message = commit.commit.message（先頭1行のみ）
    │   │       commit_author  = commit.commit.author.name
    │   │       commit_at      = commit.commit.author.date
    │   │   └─ API 失敗時：各フィールドを null として記録し、処理続行（ビルドは妨げない）
    │   │
    │   ├─ [事前チェック] pipeline.sh 実行前に以下を確認し、不足時は ERROR ログ＋deploy_failure Webhook 通知、このエントリをスキップ
    │   │   ├─ ディスク空き容量 ≥ max(出力サイト推定サイズ × 3, 64MiB)。取得は `syscall.Statfs(outDir)` を使用する
    │   │   ├─ `adlaire-ci-build` が存在し実行可能であること（`os.Stat` と mode bit）
    │   │   └─ `/usr/local/bin/adlaire-ci-build --version` が終了コード 0 で、stdout に `adlaire-ci-build` と `ADLAIRE_CI_SPEC` を含むこと
    │   │
    │   ├─ pipeline.sh 実行（bash {src の親ディレクトリ}/.ci/pipeline.sh）
    │   │   ├─ 成功（exit 0）：INFO ログ
    │   │   │   └─ [Webhook 通知送信] on: ["success"] 設定時
    │   │   │       → .notify_config の Webhook 宛先へ POST（ペイロード: event="success", branch, build_id 等）
    │   │   │       → 送信失敗（HTTP エラー・接続失敗・タイムアウト）の場合：
    │   │   │           ERROR ログ（`NOTIFY_FAIL: url={url} status={code}`）
    │   │   │           .notify_pending へ `{"event":"success","url":"...","payload":{...},"queued_at":"<ISO8601>"}` を追記
    │   │   └─ 失敗（exit ≠ 0）：ERROR ログ、sha_file 更新せず、このエントリを failure(build_failed) として記録し、次エントリへ進む
    │   │       └─ [Webhook 通知送信] on: ["failure"] 設定時
    │   │           → .notify_config の Webhook 宛先へ POST（ペイロード: event="failure", branch, build_id 等）
    │   │           → 送信失敗の場合：ERROR ログ、.notify_pending へキューイング（success と同一形式）
    │   │
    │   └─ sha_file を新 SHA で更新（`{"sha": "<new_sha>"}` を JSON 書き込み）
    │        └─ SSH サイト転送（deploy_targets リストの各エントリへ転送 → §14a）
    │             ├─ [転送後整合性検証] ssh user@host "sha256sum /dest/<relative-path>" でリモート SHA を取得
    │             │   ├─ 全ファイルのローカル sha256 と一致 → 転送成功
    │             │   └─ 不一致またはコマンド失敗 → ERROR ログ、ペンディングキューへ再投入（§14a）
    │             └─ 整合性検証成功後 → スナップショット保存（→ §14b）
    │
    │        [出力サイズチェック] OUTPUT_SIZE_WARN_MB > 0 の場合
    │        出力サイト配下の通常ファイル合計サイズを取得し、閾値と比較：
    │            size_mb = total_site_bytes / (1024 * 1024)
    │            size_mb > OUTPUT_SIZE_WARN_MB の場合：
    │            → WARN ログ（`OUTPUT_SIZE_WARN: size={size_mb:.1f}MB threshold={OUTPUT_SIZE_WARN_MB}MB`）
    │            → ビルドログの size_warn フィールドを true に設定（§8）
    │
    └─ [サーキットブレーカー判定] API_CIRCUIT_BREAKER_THRESHOLD > 0 の場合
        全ブランチの今周回結果を集計し、全ブランチが失敗（API エラー・スキップを除くビルド失敗）の場合：
        → .build_circuit_state（JSON）の consecutive_failures を +1
        いずれか成功した場合：
        → consecutive_failures を 0 にリセット
        consecutive_failures ≥ API_CIRCUIT_BREAKER_THRESHOLD の場合：
        → ERROR ログ（`CIRCUIT_OPEN: consecutive_failures={N} threshold={API_CIRCUIT_BREAKER_THRESHOLD}`）
        → deploy_failure Webhook 通知（reason="circuit_open"）
        → .build_circuit_state に `{"open": true, "consecutive_failures": N, "opened_at": "<ISO8601>"}` を書き込み
        → 以降のポーリング周回はビルドをスキップ（SHA チェックも行わない）
        → POST /api/circuit-breaker/reset でリセット可能（open: false、consecutive_failures: 0 に戻す）

    ├─ [ログクリーンアップ] LOG_KEEP_N > 0 の場合
    │   .build_logs/ 内の {id}.json を mtime 昇順でソートし、
    │   件数が LOG_KEEP_N を超えた分を古いものから削除
    │
    └─ [週次サマリー判定] WEEKLY_SUMMARY_ENABLED = true かつ on: ["weekly_summary"] 設定の Webhook 宛先が存在する場合
        現在の曜日が WEEKLY_SUMMARY_DAY かつ現在時刻が WEEKLY_SUMMARY_HOUR:00±30分以内の場合：
        ├─ 二重送信防止チェック：.build_state の weekly_summary_sent_date と当日の日付（YYYY-MM-DD）を比較
        │   → 一致（本日送信済み）→ スキップ（INFO ログ: `WEEKLY_SUMMARY_SKIP: already sent today`）
        │   → 不一致（未送信）→ 以下を実行
        ├─ 過去 7 日間の .build_logs/{id}.json を集計：
        │       成功件数・失敗件数・成功率（%）・平均ビルド時間（秒）・最長ビルド（秒・ID）
        ├─ on: ["weekly_summary"] 設定の Webhook 宛先へ POST
        │       ペイロード: event="weekly_summary", period="7d", success_count, failure_count, success_rate, avg_duration_seconds, max_duration_seconds, max_duration_build_id
        ├─ 送信成功 → .build_state に `weekly_summary_sent_date: "YYYY-MM-DD"` を記録（INFO ログ）
        └─ 送信失敗（HTTP エラー・接続失敗）→ ERROR ログ、.notify_pending へキューイング（他の Webhook 失敗と同一形式）
```

---

## 14. `pipeline.sh`

リポジトリの `.ci/pipeline.sh` にビルド手順を記述する。

- 実行権限（`chmod +x`）が必要
- 先頭の shebang は `#!/usr/bin/env bash` または `#!/bin/bash` とする
- `set -euo pipefail` を shebang 直後に記述する
- 終了コード `0` で成功、`0` 以外で失敗とみなす
- runner は `bash {src の親ディレクトリ}/.ci/pipeline.sh` として実行し、作業ディレクトリは `src` の親ディレクトリに設定する
- 環境変数として `ADLAIRE_CI_SRC`、`ADLAIRE_CI_OUT`、`ADLAIRE_CI_BRANCH`、`ADLAIRE_CI_BUILD_ID` を渡す
- `pipeline.sh` は最終的に `ADLAIRE_CI_OUT` のパスへ HTML を生成しなければならない

**例：**
```bash
#!/usr/bin/env bash
set -euo pipefail
/usr/local/bin/adlaire-ci-build --src "$ADLAIRE_CI_SRC" --out "$ADLAIRE_CI_OUT"
```

ビルド実行コマンドは `pipeline.sh` 内に直接記述する（`components/runner.go` は参照しない）。`adlaire-ci-build` は `components/builder.go` から生成した Go 版バイナリである。

**runner からの実行契約：**

| 項目 | 仕様 |
|------|------|
| command | `bash {srcParent}/.ci/pipeline.sh` |
| working directory | `{srcParent}` |
| timeout | `.server_config.build_timeout_seconds` が存在すればその値、なければ `300` 秒。 |
| stdout/stderr | それぞれ最大 1 MiB までメモリに保持し、超過分は末尾 1 MiB を保存する。超過時は `stdout_truncated` / `stderr_truncated` を `true` にする。 |
| 環境変数 | 親 process の環境を引き継ぎ、下表の値で上書きする。 |

| 環境変数 | 値 |
|----------|----|
| `ADLAIRE_CI_SRC` | `BranchTarget.Src` |
| `ADLAIRE_CI_OUT` | `BranchTarget.Out` |
| `ADLAIRE_CI_BRANCH` | `BranchTarget.Branch` |
| `ADLAIRE_CI_BUILD_ID` | build id |
| `ADLAIRE_CI_TARGET_FILE` | `BranchTarget.TargetFile` |
| `ADLAIRE_CI_STATE_DIR` | `RunnerConfig.StateDir` |

`pipeline.sh` が timeout した場合、runner は process group を終了し、`target_status="failure_build"`、`error="pipeline timeout"` として記録する。timeout 時も stdout/stderr の取得済み内容は `.build_logs/{id}.json` に保存する。

**GitHub API 固定契約：**

| 項目 | 仕様 |
|------|------|
| Base URL | `https://api.github.com` |
| User-Agent | `adlaire-ci-runner` |
| Authorization | `Bearer {token}` |
| Accept | `application/vnd.github+json` |
| API version header | `X-GitHub-Api-Version: 2022-11-28` |
| timeout | 30 秒 |
| retry 対象 | network error、timeout、HTTP `429`、`500`、`502`、`503`、`504` |
| retry 対象外 | HTTP `400`、`401`、`403`（rate limit を除く）、`404`、`422` |
| backoff | `API_RETRY_BASE_SECONDS * 2^attempt` 秒。attempt は 0 始まり。 |

HTTP `401` は `failure_api` とし、ERROR ログ `GITHUB_AUTH_FAILED` を出す。HTTP `404` は `target_file` または branch 設定不正として `failure_api` とし、ERROR ログ `GITHUB_NOT_FOUND: branch={branch} target={target_file}` を出す。HTTP `403` で `X-RateLimit-Remaining: 0` の場合のみ rate limit として reset まで待機する。

---

## 14a. SSH サイト転送

本節は、Go 版 CI ランナーの SSH 転送標準仕様である。

Go 版 `components/runner.go` は、`pipeline.sh` 成功後に、出力サイトディレクトリを SSH 経由で静的コンテンツ配信サーバーへ転送する。本節を SSH 転送の正本仕様とする。

`components/runner.go` は `pipeline.sh` 成功後に、出力サイトディレクトリ配下の全ファイルを SSH 経由で静的コンテンツ配信サーバーへ転送する。scp・rsync は使用しない。SSH コマンドは `ssh` バイナリを `exec.CommandContext` で直接起動し、`/bin/sh -c` を使わない。

### 設定値

`BRANCH_TARGETS` 各エントリの `deploy_targets` リスト内で管理する（→ §12）。

| フィールド | 説明 | 例 |
|-----------|------|-----|
| `host` | 配信サーバーのホスト名 / IP | `"192.0.2.1"` |
| `user` | SSH 接続ユーザー | `"deploy"` |
| `dest_dir` | 配信サーバー上の転送先ディレクトリ | `"/var/www/html"` |

転送対象は `BRANCH_TARGETS[n]["out"]` のディレクトリ配下にある通常ファイルすべてとする。`deploy_targets` に複数エントリを定義した場合は全ての転送先へ順次転送する。

### 差分検出

転送前にリモートサーバーで対象ファイルごとの SHA256 ハッシュを取得し、ローカルファイルのハッシュと比較する。

```bash
# components/runner.go が os/exec 経由で実行
ssh <user>@<host> sha256sum <dest_dir>/<relative-path>
```

- ハッシュが一致 → 当該ファイルをスキップ（`SKIP` ログを記録）
- ハッシュが不一致、またはリモートにファイルが存在しない → 当該ファイルを転送する

### 転送

stdin パイプ経由で SSH 転送する。

```bash
# components/runner.go が os/exec（StdinPipe）経由で実行
ssh <user>@<host> 'mkdir -p <dest_dir>/<relative-dir> && tee <dest_dir>/<relative-path>'
```

runner は local file を開き、SSH process の stdin へ `io.Copy` で送る。リモート側 stdout は破棄してよいが、stderr は失敗理由として `.build_logs/{id}.json.error` と ERROR ログへ記録する。

### ペンディングキュー

転送失敗時（接続エラー・認証失敗等）は `PENDING_FILE`（JSON）へエントリを追記する。

```json
[
  {
    "branch_idx": 0,
    "deploy_idx": 0,
    "out": "/opt/adlaire-builder/dist/site",
    "host": "192.0.2.1",
    "user": "deploy",
    "dest_dir": "/var/www/html",
    "failed_at": "2026-09-15T10:00:00Z",
    "retry_count": 1
  }
]
```

- `components/runner.go` 起動時（`BRANCH_TARGETS` 処理前）に `PENDING_FILE` を読み込み、エントリごとに再試行する（→ §13 処理フロー）
- 再試行成功時にエントリを削除する。失敗時は `retry_count` をインクリメントして保持する
- SSH 転送失敗 Webhook 通知（`deploy_failure` イベント）を送信する（on: `["deploy_failure"]` 設定時）
- 同一 `out`、`host`、`user`、`dest_dir` の pending エントリが既に存在する場合は新規追記せず、既存エントリの `retry_count` を +1 し、`failed_at` を最新時刻へ更新する

### 転送後整合性検証

SSH 転送完了後に、リモートファイルの SHA-256 チェックサムをローカルのものと照合する。

**検証コマンド：**
```
ssh {user}@{host} sha256sum {dest_dir}/{filename}
```
出力形式 `{hash}  {filename}` の最初のフィールドをローカル `crypto/sha256` の hex digest と比較する。

| 項目 | 仕様 |
|---|---|
| タイムアウト | 30 秒（SSH 転送タイムアウトとは独立） |
| 検証失敗時 | ERROR ログ＋ペンディングキューへ再投入。スナップショット保存はしない |
| ログフィールド | `transfer_verified: false`（`.build_logs/{id}.json` に記録） |
| 正常時 | `transfer_verified: true`（`.build_logs/{id}.json` に記録） |

### ログ

| 状態 | ログレベル | メッセージ例 |
|------|-----------|------------|
| スキップ（差分なし） | `INFO` | `SKIP site: no change` |
| 転送成功 | `INFO` | `DEPLOY site → 192.0.2.1` |
| 転送失敗→キューイング | `ERROR` | `DEPLOY FAILED site: <reason> (queued)` |
| ペンディング再試行成功 | `INFO` | `PENDING RETRY OK site → 192.0.2.1` |
| ペンディング再試行失敗 | `ERROR` | `PENDING RETRY FAILED site: <reason>` |
| 整合性検証失敗→再投入 | `ERROR` | `VERIFY FAILED site → 192.0.2.1: checksum mismatch (queued)` |

---

## 14b. スナップショット管理

本節は、Go 版 CI ランナーのスナップショット標準仕様である。

Go 版 `components/runner.go` は、SSH 転送成功後に `.snapshots/` ディレクトリへ成果物を保存する。本節をスナップショット保存、世代管理、ロールバック連携の正本仕様とする。

`components/runner.go` は SSH 転送成功後に、ビルド成果物を `.snapshots/` ディレクトリへアーカイブする。`HISTORY_KEEP_N = 0` の場合はスナップショット機能を無効化する。

### ディレクトリ構造

```
/opt/adlaire-builder/
└── .snapshots/
    ├── b20260915100000/           # build_id = b{YYYYMMDDHHmmss}
    │   └── site  # ビルド成果物のコピー
    ├── b20260914180000/
    │   └── site
    └── ...
```

- ディレクトリ名は `b{YYYYMMDDHHmmss}` 形式のビルド ID（`.build_history` の `id` と一致する）
- `BRANCH_TARGETS` に複数エントリがある場合は、同一ビルド ID ディレクトリ内に各エントリの成果物をまとめて保存する

### 世代管理

- スナップショット保存後、`.snapshots/` 内のディレクトリ数が `HISTORY_KEEP_N` を超えた場合、最古のディレクトリから順に削除する
- 削除対象ディレクトリの特定は作成日時降順ソートで行う（ディレクトリ名の辞書順 = 時系列順）

### ロールバック

`POST /api/history/{id}/rollback`（→ §22）で指定ビルド ID のスナップショットから SSH 転送を再実行する。

- ロールバック API は `components/api.go` の実装を前提とする。`components/api.go` が実装されるまでは、API 経由のロールバックはとして扱う
- `.snapshots/{id}/` が存在しない場合は `404` を返す
- 転送成功時は `.build_history` に rollback エントリを追記する

### ログ

| 状態 | ログレベル | メッセージ例 |
|------|-----------|------------|
| スナップショット保存 | `INFO` | `SNAPSHOT b20260915100000 saved` |
| 古世代削除 | `INFO` | `SNAPSHOT b20260910000000 pruned (keep_n=10)` |
| ロールバック成功 | `INFO` | `ROLLBACK b20260914180000 → 192.0.2.1 OK` |
| ロールバック失敗 | `ERROR` | `ROLLBACK b20260914180000: <reason>` |

---

## 15. ログ

本節は、Go 版 `components/runner.go` の stdout ログと構造化ビルドログを定義する。

### stdout ログ

| レベル | 出力条件 |
|--------|---------|
| `INFO` | 起動、変更なしスキップ、ビルド開始・完了、SHA 更新 |
| `WARNING` | — |
| `ERROR` | トークン読み込み失敗、API 失敗、ビルド失敗 |
| `DEBUG` | API レスポンス詳細等（`LOG_LEVEL = "DEBUG"` 時のみ） |

stdout は Go 標準ライブラリ `log/slog` で出力し、systemd が journald に転送する。独自 logger 実装を使用してはならない。

### 構造化ビルドログ

Go 版 `components/runner.go` は、ビルドごとに `.build_logs/{id}.json` を作成する。

| 項目 | 内容 |
|------|------|
| ビルドログファイル | ビルドごとに `.build_logs/{id}.json` を作成する。 |
| stdout / stderr 保存 | `pipeline.sh` の標準出力・標準エラーをビルドログへ保存する。 |
| 変換レポート取り込み | `components/builder.go` が出力する `[REPORT]` 行をパースし、`tables_count`、`code_blocks_count` 等へ変換して保存する。 |
| 警告取り込み | `[WARN]` 行を配列として保存し、`warnings` 件数と整合させる。 |
| ビルド所要時間 | `started_at`、`finished_at`、`duration_seconds` を保存する。 |
| コミット情報 | ビルド対象 commit の SHA、message、author、date を保存する。 |
| 転送検証結果 | SSH 転送後整合性検証の結果として `transfer_verified` を保存する。 |
| 出力サイズ警告 | `OUTPUT_SIZE_WARN_MB` 超過時に `size_warn: true` を保存する。 |
| ログ世代管理 | `LOG_KEEP_N` を超過した `.build_logs/{id}.json` を古いものから削除する。 |

これらのログ項目を実装対象に含める時点で、§10a の実装対象、§12 の設定値、§13 の処理フロー、§22 の API レスポンス仕様と整合させる。

**`.build_logs/{id}.json` 完全 schema：**

```json
{
  "id": "b20260916010000",
  "branch": "main",
  "target_file": "docs",
  "target_status": "success",
  "trigger": "polling",
  "started_at": "2026-09-16T01:00:00Z",
  "finished_at": "2026-09-16T01:00:12Z",
  "duration_seconds": 12,
  "commit": {
    "sha": "abcdef",
    "message": "Update docs",
    "author": "example",
    "date": "2026-09-16T00:59:00Z"
  },
  "blob_sha": "012345",
  "previous_blob_sha": "",
  "pipeline": {
    "exit_code": 0,
    "stdout": "...",
    "stderr": "",
    "stdout_truncated": false,
    "stderr_truncated": false
  },
  "report": {
    "pages": 1,
    "headings": 1,
    "tables_count": 0,
    "code_blocks_count": 0,
    "warnings_count": 0,
    "size_warn": false,
    "broken_links": 0,
    "heading_skips": 0,
    "reading_time": 1,
    "theme": "adlaire-default"
  },
  "warnings": [],
  "deploy": [
    {
      "host": "192.0.2.1",
      "user": "deploy",
      "dest_dir": "/var/www/html",
      "status": "success",
      "transfer_verified": true,
      "files_total": 3,
      "files_uploaded": 3,
      "files_skipped": 0,
      "bytes_uploaded": 1234,
      "error": null
    }
  ],
  "snapshot_id": "b20260916010000",
  "error": null
}
```

| キー | 型 | 必須 | 条件 |
|------|----|------|------|
| `id` | string | 必須 | build id。 |
| `branch` | string | 必須 | `BranchTarget.Branch`。 |
| `target_file` | string | 必須 | `BranchTarget.TargetFile`。 |
| `target_status` | string | 必須 | §13 の分類値。 |
| `trigger` | string | 必須 | §13 の `trigger` 有効値。 |
| `started_at` / `finished_at` | string | 必須 | UTC ISO 8601。 |
| `duration_seconds` | integer | 必須 | 0 以上。 |
| `commit` | object | 必須 | 取得失敗時は各値を `null` にする。 |
| `blob_sha` | string/null | 必須 | GitHub Trees API で検出した SHA。API 失敗時は `null`。 |
| `previous_blob_sha` | string | 必須 | sha_file 読み込み値。未設定時は空文字。 |
| `pipeline` | object | 必須 | pipeline 未実行時も `exit_code:null`、stdout/stderr 空文字で保存する。 |
| `report` | object/null | 必須 | `[REPORT]` が存在しない場合は `null`。 |
| `warnings` | string[] | 必須 | `[WARN]` 行から prefix を除いた文字列配列。 |
| `deploy` | object[] | 必須 | deploy target がない場合は空配列。 |
| `snapshot_id` | string/null | 必須 | snapshot 未作成時は `null`。 |
| `error` | string/null | 必須 | 成功時 `null`。失敗時は固定文言を保存する。 |

`.build_logs/{id}.json` は `encoding/json` で生成し、末尾改行を付ける。未知キーを追加してはならない。`report.tables_count` と `report.code_blocks_count` は stdout `[REPORT]` の `tables`、`code_blocks` から変換して保存する。

**`.build_history` JSON Lines schema：**

`.build_history` は 1 行 1 JSON object とし、各行は次の schema を満たす。

```json
{"id":"b20260916010000","branch":"main","target_file":"docs","status":"success","trigger":"polling","started_at":"2026-09-16T01:00:00Z","finished_at":"2026-09-16T01:00:12Z","duration_seconds":12,"commit_sha":"abcdef","blob_sha":"012345","pages":1,"warnings":0,"size_warn":false,"output_sha256":null,"rollback_from":null}
```

`status` は `target_status` と同じ値を保存する。`trigger` は §13 の有効値を保存する。`output_sha256` は出力サイト全体 manifest の SHA-256 hex とし、manifest 生成に失敗した場合のみ `null` を許可する。JSON Lines 追記は `O_APPEND|O_CREATE|O_WRONLY` で行い、1 行全体を書き込んでから file sync する。

**固定エラー文言：**

| 条件 | `error` |
|------|---------|
| GitHub 認証失敗 | `github authentication failed` |
| GitHub target 不在 | `github target not found` |
| GitHub API 全再試行失敗 | `github api failed` |
| Blob decode 失敗 | `blob decode failed` |
| Markdown 書き出し失敗 | `source write failed` |
| 事前チェック失敗 | `precheck failed` |
| pipeline timeout | `pipeline timeout` |
| pipeline 非 0 | `pipeline failed` |
| deploy pending | `deploy pending` |
| 状態ファイル書き込み失敗 | `state write failed` |

---

## 15a. `components/runner.go` 受け入れ fixture

Go 版 `components/runner.go` の初期実装は、本節の fixture をすべて満たすまで完了として扱わない。fixture ファイルは実装 PR で `testdata/runner/` 配下へ追加する。外部 GitHub API と SSH サーバーへ実接続するテストは初期 fixture に含めず、HTTP test server と fake `ssh` executable で再現する。

### Fixture R1: CLI 異常系

| 実行 | 終了コード | stdout | stderr |
|------|------------|--------|--------|
| `adlaire-ci-runner --help` | `0` | `Usage: adlaire-ci-runner [--state-dir path] [--once] [--version] [--help]` | 空 |
| `adlaire-ci-runner --state-dir relative` | `2` | 空 | `state directory must be absolute: relative` |
| `adlaire-ci-runner --unknown` | `2` | 空 | `unknown option: --unknown` |

### Fixture R2: 変更なし skip

**前提状態：**

```text
testdata/runner/r2/state/.github_token
testdata/runner/r2/state/.last_sha
testdata/runner/r2/state/.branch_config
```

`.last_sha`:

```json
{"sha":"blob-1"}
```

GitHub Trees API fake response は `target_file=docs` の SHA として `blob-1` を返す。

**期待結果：**

- 終了コード `0`。
- `.last_sha` は変更しない。
- `.build_logs/` に新規 log を作成しない。
- `.build_history` に追記しない。
- stdout slog に INFO ログ `NO_CHANGE: branch=main target=docs sha=blob-1` を 1 件出力する。

### Fixture R3: 変更あり build 成功 deploy なし

**前提状態：**

- `.last_sha` は `{"sha":"old-blob"}`。
- GitHub Trees API fake response は `new-blob` を返す。
- GitHub Blobs API fake response は UTF-8 Markdown を Base64 で返す。
- `deploy_targets` は空配列。
- fake `.ci/pipeline.sh` は終了コード `0` で、stdout に `adlaire-ci-build` の `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default` を出力する。

**期待結果：**

- 終了コード `0`。
- `.last_sha` は `{"sha":"new-blob"}` に atomic write される。
- `.build_logs/{id}.json` が作成され、`target_status="success"`、`pipeline.exit_code=0`、`report.pages=1`、`report.tables_count=0`、`deploy=[]`、`error=null` を含む。
- `.build_history` に同じ `id` の JSON Lines が 1 行追記される。
- `.build_state.running` は終了時 `false`、`current_build_id` は `null`。
- `.build_lock` は終了時に存在しない。

### Fixture R4: pipeline 失敗

**前提状態：**

- GitHub fake response は変更ありを返す。
- fake `.ci/pipeline.sh` は終了コード `7`、stdout `before fail`、stderr `failed` を出力する。

**期待結果：**

- 終了コード `1`。
- `.last_sha` は旧 SHA のまま。
- `.build_logs/{id}.json` は `target_status="failure_build"`、`pipeline.exit_code=7`、`pipeline.stdout="before fail"`、`pipeline.stderr="failed"`、`error="pipeline failed"` を含む。
- `.build_history` に `status="failure_build"` の行を追記する。
- deploy、snapshot は実行しない。

### Fixture R5: deploy pending

**前提状態：**

- pipeline は成功する。
- `deploy_targets` は 1 件。
- fake `ssh` は転送時に終了コード `255` を返す。

**期待結果：**

- 終了コード `1`。
- `.last_sha` は新 SHA に更新する。
- `.pending_transfers` に `out`、`host`、`user`、`dest_dir`、`failed_at`、`retry_count=1` を持つ entry を 1 件追加する。
- `.build_logs/{id}.json` は `target_status="success_deploy_pending"`、`deploy[0].status="pending"`、`deploy[0].transfer_verified=false`、`error="deploy pending"` を含む。
- snapshot は作成しない。

### Fixture R6: lock 競合

`.build_lock` が存在し、`pid` が `/proc/{pid}` に存在する実行中 PID を指す場合、runner は終了コード `0` で終了し、`.build_state`、`.build_logs/`、`.build_history` を変更しない。stdout slog に `BUILD_SKIP: already running (PID {pid})` を出力する。

### Fixture R7: 状態破損

`.notify_pending` が JSON として壊れている場合、runner は `.notify_pending.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、`.notify_pending` を `[]` で再生成する。その後、通常処理を継続する。退避ファイル名の timestamp は UTC とし、秒単位で固定する。

---

## 16. systemd タイマー

### `adlaire-ci.service`（oneshot）

```ini
[Unit]
Description=Adlaire CI Runner

[Service]
Type=oneshot
User=deploy
ExecStart=/usr/local/bin/adlaire-ci-runner --state-dir /opt/adlaire-builder
StandardOutput=journal
StandardError=journal
```

### `adlaire-ci.timer`（5分ごと定期実行）

```ini
[Unit]
Description=Adlaire CI Runner Timer
After=network-online.target
Wants=network-online.target

[Timer]
OnBootSec=1min
OnUnitActiveSec=5min

[Install]
WantedBy=timers.target
```

```bash
sudo systemctl enable --now adlaire-ci.timer  # タイマー登録・起動
sudo systemctl list-timers adlaire-ci          # 次回実行時刻確認
sudo journalctl -u adlaire-ci -f               # ログ確認
```

---

## 17. GitHub 側設定

| 項目 | 内容 |
|------|------|
| PAT スコープ | `contents: read`（読み取り専用）のみ |
| PAT の種類 | Fine-grained PAT（特定リポジトリのみ許可）を使用する。 |
| Webhook 設定（ポーリング方式） | **不要**（デフォルト。`BRANCH_TARGETS` によるポーリングのみ使用する場合） |
| Webhook 設定（受信方式） | GitHub リポジトリ設定 → Webhooks → Add webhook で `POST /api/webhook` の URL・Secret を設定する（→ §22）。イベントは `push` のみ選択する。**外部公開エンドポイントが必要**（リバースプロキシ経由） |

---

## 18. 初回セットアップ手順

本節は Go 版 CI ランナー導入手順である。`adlaire-ci-build`、`adlaire-ci-runner`、管理 API 導入後の `adlaire-ci-api` の各バイナリを配置する。

```bash
# 1. deploy ユーザー作成
sudo useradd -m -s /bin/bash deploy

# 2. 作業ディレクトリ作成
sudo mkdir -p /opt/adlaire-builder/repo
sudo chown -R deploy:deploy /opt/adlaire-builder

# 3. GitHub PAT を保存（Fine-grained PAT、contents: read のみ）
echo "<PAT>" | sudo -u deploy tee /opt/adlaire-builder/.github_token
sudo chmod 600 /opt/adlaire-builder/.github_token

# 3b. CI サーバー → 配信サーバー SSH 鍵設定
#     SSH 転送機能を使用する場合のみ実行する。
#     deploy ユーザーの SSH 鍵を生成（既存鍵がある場合はスキップ）
sudo -u deploy ssh-keygen -t ed25519 -f /home/deploy/.ssh/id_ed25519 -N ""
#     公開鍵を配信サーバーへ登録（配信サーバー側で実行）
#     cat /home/deploy/.ssh/id_ed25519.pub >> ~/.ssh/authorized_keys
#     初回接続時の known_hosts 登録
sudo -u deploy ssh-keyscan -H <配信サーバーIP> >> /home/deploy/.ssh/known_hosts

# 4. SHA キャッシュファイルを初期化
printf '%s\n' '{"sha":""}' | sudo -u deploy tee /opt/adlaire-builder/.last_sha
sudo chmod 600 /opt/adlaire-builder/.last_sha

# 5. Go 版バイナリを配置
sudo install -m 0755 adlaire-ci-build /usr/local/bin/adlaire-ci-build
sudo install -m 0755 adlaire-ci-runner /usr/local/bin/adlaire-ci-runner

# 6. systemd ユニットを登録・タイマー起動
sudo cp adlaire-ci.service /etc/systemd/system/
sudo cp adlaire-ci.timer   /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now adlaire-ci.timer

# 以降はの管理 API サーバー導入手順
# 7. adlaire-ci-api を配置
sudo install -m 0755 adlaire-ci-api /usr/local/bin/adlaire-ci-api

# 8. 認証情報ファイルを初期化（初期パスワード: admin）
sudo -u deploy /usr/local/bin/adlaire-ci-api --init-credentials --state-dir /opt/adlaire-builder
sudo chmod 600 /opt/adlaire-builder/.admin_credentials

# 9. systemd ユニットを登録・起動
sudo cp adlaire-ci-api.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now adlaire-ci-api
```

---

## 19. 管理 API サーバー 既知の制限

| 制限 | 詳細 |
|------|------|
| セッションはインメモリ管理 | 再起動で全セッションが消去される |
| HTTPS 非対応 | `components/api.go` は TLS listener、証明書読み込み、HTTPS redirect を実装しない。HTTP listener のみ起動する |
| シングルユーザー専用 | `POST /api/login` は `password` のみを受け取り、`username`、`user_id`、role、group を受け取らない。該当 field を含む request は `422` を返す |
| 並列リクエストの制限 | `components/api.go` は Go 標準ライブラリ `net/http` の標準サーバーで処理し、独自の接続数上限、IP 単位 rate limit、worker pool を実装しない |

---

## 20. CI ランナー 既知の制限

### 20.1 Go 版 runner の制限

| 制限 | 詳細 |
|------|------|
| ポーリング遅延 | 変更検出は systemd timer の実行間隔に依存する。即時反応が必要な場合は `POST /api/webhook` を併用する。 |
| `pipeline.sh` 起動 | ビルド起動は `src` と同じディレクトリ配下の `.ci/pipeline.sh` を標準とする。YAML 形式のパイプライン定義は本ファイルでは定義しない。 |
| `BRANCH_TARGETS` 直列処理 | 複数エントリはリスト順に順次処理する。並列処理は行わない。1 件の処理が失敗しても、失敗をログと `.build_logs/{id}.json` に記録した上で次エントリへ進む。 |
| GitHub API リトライ | GitHub API 失敗時は `API_RETRY_MAX` 回まで指数バックオフで再試行する。全試行失敗時は ERROR ログを記録し、当該ターゲットのビルドをスキップする。SHA は更新しない。 |
| ビルド失敗時の扱い | `pipeline.sh` が非 0 で終了した場合は ERROR ログを出し、SHA を更新しない。次回実行では同じ blob SHA を再検出して再度ビルド対象になる。 |

### 20.2 標準機能の制限

以下は管理 API または追加拡張と連携する場合の制限である。

| 制限 | 詳細 |
|------|------|
| Webhook 受信の外部公開 | `POST /api/webhook` は `components/api.go`（`127.0.0.1` バインド）で受信するため、GitHub から直接受信する構成ではリバースプロキシと TLS 終端が必要。 |
| ペンディングキュー | ペンディング再試行が失敗した場合、`retry_count` を 1 増やしてエントリを保持する。runner による自動放棄は行わない。削除は転送成功時、または管理 API / 手動運用で明示的に削除する場合に限定する。 |
| ペンディングキュー肥大化 | `queue_max_size` を超えた新規投入は ERROR ログを記録し、新規エントリを追加しない。既存エントリは削除しない。 |
| サーキットブレーカー | 連続失敗回数が `API_CIRCUIT_BREAKER_THRESHOLD` 以上になった場合はポーリングを停止し、`POST /api/circuit-breaker/reset` でのみ復帰する。 |

---

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

対象コンポーネント `components/api.go` は、Go 標準ライブラリ `net/http` で実装し、管理ツールからの API リクエストを受け付ける。`components/runner.go` とは独立して常駐する。

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
| 権限不足 | 認証済み token の scope が不足する場合は `403 Forbidden` と `{"error":"Forbidden"}` を返す。初期仕様の API token scope は `read` のみであり、API token は `GET` のみ許可する。`POST`、`DELETE`、設定変更、ビルド起動、token 発行、secret 更新は `/api/login` で発行された管理セッション token のみ許可する。 |
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

### 22.0a 状態ファイル共通仕様

`components/api.go` および拡張後 `components/runner.go` が読み書きする状態ファイルは、下表の初期値、形式、更新責務に従う。表にない状態ファイルを追加してはならない。追加が必要な場合は、先に本節へパス、形式、初期値、更新責務、破損時の扱いを追記する。

| パス | 形式 | 初期値 | 更新責務 | 破損時の扱い |
|------|------|--------|----------|--------------|
| `.admin_credentials` | JSON object | `--init-credentials` で生成 | `components/api.go` | 起動時に ERROR ログを出し、HTTP サーバーを起動しない。 |
| `.server_config` | JSON object | `{}` | `components/api.go` | `.server_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、空 object で再生成する。 |
| `.notify_config` | JSON object | `{"webhooks":[],"on":[],"summary":{"enabled":false,"interval":"weekly","hour":9,"day_of_week":1},"email":{"enabled":false,"to":[],"on":[]}}` | `components/runner.go` / `components/api.go` | `.notify_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、初期値で再生成する。 |
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
| `.repo_config` | JSON object | `{}` | `components/api.go` | `.repo_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、スクリプト定数へフォールバックする。 |
| `.config_log` | JSON Lines | 空ファイル | `components/api.go` | 読み込み可能な行のみ返し、壊れた行は無視する。 |
| `.access_log` | JSON Lines | 空ファイル | `components/api.go` | 読み込み可能な行のみ返し、壊れた行は無視する。 |
| `.api_access_log` | JSON Lines | 空ファイル | `components/api.go` | 読み込み可能な行のみ返し、壊れた行は無視する。秘密情報は記録しない。 |
| `.webhook_secret` | text | 不在 | `components/api.go` | 読み込み不能時は Webhook 受信を `501` で拒否する。 |
| `.webhook_events.json` | JSON Lines | 空ファイル | `components/api.go` | 読み込み可能な行のみ返し、壊れた行は無視する。 |
| `.access_control` | JSON object | `{"allow":[]}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.hooks` | JSON object | `{"hooks":[]}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.maintenance` | JSON object | `{"enabled":false,"reason":null,"since":null}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.api_tokens` | JSON object | `{"tokens":[]}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.alert_rules` | JSON object | `{"rules":[]}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.tag_rules` | JSON object | `{"rules":[]}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.pipeline_config` | JSON object | `{"extra_args":[],"env":{}}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.notes` | UTF-8 text | 空文字列 | `components/api.go` | 読み込み不能時は `500` を返し、自動上書きしない。 |
| `.smtp_config` | JSON object | SMTP 未設定値 | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.smtp_secret` | text | 不在 | `components/api.go` | 読み込み不能時は SMTP 送信を `422` で拒否する。 |
| `.dashboard_layout` | JSON object | `{"widgets":["status","stats","schedule","alerts","disk","rate_limit","snapshots","maintenance","queue"]}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |

`.build_logs/archive/` は gzip 圧縮済み build log の保存先ディレクトリである。初期値は空ディレクトリとし、`components/runner.go` または `POST /api/logs/archive` が必要時に作成する。圧縮済みファイル名は `{id}.json.gz` 固定とし、通常 `.build_logs/{id}.json` と同じ build id を表す。

JSON Lines ファイルは、1 行につき 1 JSON object とする。追記時は末尾に改行を必ず付ける。秘密情報を含む可能性のある `.admin_credentials`、`.github_token`、`.webhook_secret`、`.smtp_secret` は mode `600` を必須とする。

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
| `force_build_interval_hours` | integer | `0` | 0〜8760 | `POST /api/schedule/force-interval`, `GET /api/schedule` | `0` は強制再ビルド無効。 |
| `build_cooldown_seconds` | integer | `0` | 0〜86400 | `POST /api/schedule/cooldown`, `GET /api/schedule` | `0` はクールダウン無効。 |
| `schedule_interval_seconds` | integer | `300` | 30〜86400 | `POST /api/schedule/interval`, `GET /api/schedule` | systemd timer 更新値。 |
| `schedule_paused` | boolean | `false` | `true` / `false` | `POST /api/schedule/pause`, `POST /api/schedule/resume`, `GET /api/schedule` | 自動ポーリング停止状態。 |
| `allowed_hours` | object/null | `null` | `{"from":0〜23,"to":0〜23}` または `null` | `POST /api/schedule/allowed-hours`, `GET /api/schedule` | UTC の自動ビルド許可時間帯。 |

`.server_config` の `POST /api/config` では `force_build_interval_hours`、`build_cooldown_seconds`、`schedule_interval_seconds`、`schedule_paused`、`allowed_hours` を直接更新してはならない。これらは専用スケジュール API からのみ更新する。

**`.notify_config` schema：**

| キー | 型 | 既定値 | 許容値 | 説明 |
|------|----|--------|--------|------|
| `webhooks` | object[] | `[]` | 下記 Webhook object | 通知先一覧。 |
| `on` | string[] | `[]` | `"start"`, `"success"`, `"failure"`, `"deploy_failure"`, `"weekly_summary"`, `"config_corrupt"` | 通知イベント。重複は除去する。 |
| `summary` | object | 下記 Summary object | 下記 | 定期サマリー設定。 |
| `email` | object | 下記 Email object | 下記 | メール通知設定。SMTP 詳細は `.smtp_config` / `.smtp_secret` を正とする。 |

Webhook object:

| キー | 型 | 既定値 | 許容値 | 説明 |
|------|----|--------|--------|------|
| `url` | string | 必須 | URL 検証に従う | 送信先 URL。 |
| `label` | string | `""` | 0〜64 文字 | 管理画面表示名。 |
| `enabled` | boolean | `true` | boolean | `false` の宛先へは送信しない。 |
| `on` | string[] | `[]` | `"start"`, `"success"`, `"failure"`, `"deploy_failure"`, `"weekly_summary"`, `"config_corrupt"`, `"*"` | この宛先が受け取るイベント。空配列の場合は top-level `on` に従う。 |
| `payload_template` | string/null | `null` | 0〜10000 文字または `null` | `null` は標準 payload。 |
| `retry_count` | integer | `2` | 0〜10 | 送信失敗時の追加試行回数。 |
| `retry_interval_seconds` | integer | `30` | 1〜3600 | 再試行間隔。 |
| `secret` | string/null | `null` | 1〜256 文字または `null` | 保存時は平文保存可。ただし GET/backup では `"***"` へマスクする。 |

Summary object:

| キー | 型 | 既定値 | 許容値 |
|------|----|--------|--------|
| `enabled` | boolean | `false` | boolean |
| `interval` | string | `"weekly"` | `"daily"` / `"weekly"` |
| `hour` | integer | `9` | 0〜23 |
| `day_of_week` | integer | `1` | 0〜6 |

Email object:

| キー | 型 | 既定値 | 許容値 |
|------|----|--------|--------|
| `enabled` | boolean | `false` | boolean |
| `to` | string[] | `[]` | メールアドレス配列、最大 50 件 |
| `on` | string[] | `[]` | `"start"`, `"success"`, `"failure"` |

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

**`.api_tokens` schema：**

```json
{
  "tokens": [
    {
      "id": "tok001",
      "label": "監視用",
      "scope": "read",
      "token_hash": "<sha256_hex>",
      "created_at": "2026-09-15T10:00:00Z",
      "last_used_at": null,
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
| `scope` | string | 必須 | `"read"` | 初期仕様では read のみ。 |
| `token_hash` | string | 必須 | SHA-256 hex | token 本体は保存しない。 |
| `created_at` | string | 必須 | ISO 8601 | 作成日時。 |
| `last_used_at` | string/null | 必須 | ISO 8601 または `null` | 最終使用日時。 |
| `revoked_at` | string/null | 必須 | ISO 8601 または `null` | 失効日時。`null` は有効。 |

`POST /api/tokens` は token 本体を `act_` + 32 byte 相当のランダム文字列として生成し、レスポンス時に 1 回だけ返す。保存する値は `token_hash` のみとする。`DELETE /api/tokens/{id}` は物理削除せず、`revoked_at` を現在時刻へ更新する。

**`.maintenance` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `enabled` | boolean | 必須 | boolean | メンテナンス有効状態。 |
| `reason` | string/null | 必須 | 0〜500 文字または `null` | 理由。 |
| `since` | string/null | 必須 | ISO 8601 または `null` | 有効化日時。 |

`enabled: false` の場合、`reason` と `since` は `null` とする。`POST /api/maintenance/enable` は `enabled: true`、`reason`、`since` を同時に保存する。

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
| `trigger` | string | 必須 | `"polling"`, `"force_interval"`, `"manual"`, `"webhook"`, `"retry_pending_transfer"`, `"startup_config_integrity"`, `"rollback"` | 起動種別。 |
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
| `dry_run` | boolean | 必須 | boolean | dry-run log の場合のみ `true`。通常 build は `false`。 |
| `attempts` | object[] | 必須 | 1 件以上 | build / deploy の試行履歴。dry-run は空配列ではなく検証 attempt 1 件を保存する。 |
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
| `finished_at` | string/null | 必須 | 完了時刻。dry-run 検証のみでも設定する。 |
| `stage` | string | 必須 | `"dry_run"`, `"github"`, `"pipeline"`, `"deploy"` のいずれか。 |
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
| `POST /api/login` | `.admin_credentials` | `.admin_credentials`, `.access_log` | 成功時のみ `login_count` を更新する。 |
| `POST /api/logout` | メモリ上 session | メモリ上 session | ファイルは更新しない。 |
| `POST /api/change-password` | `.admin_credentials` | `.admin_credentials` | 現 session 以外をメモリから削除する。 |
| `GET /api/access-log` | `.access_log` | なし | 壊れた行は無視し、新しい順で返す。 |
| `GET /api/api-access-log` | `.api_access_log` | なし | 壊れた行は無視し、新しい順で返す。 |
| `GET /api/sessions` | メモリ上 session | なし | token 本体は返さない。 |
| `POST /api/sessions/revoke-all` | メモリ上 session | メモリ上 session, `.access_log` | 現 session 以外を削除する。 |
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
| `POST /api/tokens` | `.api_tokens` | `.api_tokens`, `.access_log` | token 本体は作成時のみ返し、保存はハッシュのみ。 |
| `DELETE /api/tokens/{id}` | `.api_tokens` | `.api_tokens`, `.access_log` | 対象 token を失効する。 |

### 22.0e API 完全契約表

本表は API 実装、SDK 実装、標準管理ツール実装の契約インデックスである。実装者は endpoint を追加、削除、名称変更、body 変更、response 変更する前に本表を先に更新する。下表に存在しない endpoint は実装対象外とする。SHA reset 専用 endpoint とサマリー送信専用 endpoint は定義しない。

`Request` が `none` の場合、request body を受け付けない。空 JSON object `{}` も送信してはならない。`Response` は成功時 body の schema 名または最小 object を示す。詳細 schema は §22.0c、各 endpoint の個別例、§23 SDK 仕様、§24 UI 仕様を正とする。

`{message}` は `{"message": string}` を意味する。`{message,...}` 形式の response では `message` を必須キーとし、その他のキーも表記どおり必須とする。`?` が付いたキーだけを任意キーとする。成功時に空 body、`null` body、HTTP 204 は使用しない。

| Endpoint | Request | Response | Success | Errors | Read | Write | SDK | UI |
|----------|---------|----------|---------|--------|------|-------|-----|----|
| `POST /api/login` | `{password}` | `{token,must_change}` | `200` | `401`, `422`, `429`, `500` | `.admin_credentials` | `.admin_credentials`, `.access_log` | `login()` | ログイン |
| `POST /api/logout` | none | `{message}` | `200` | `401` | memory session | memory session | `logout()` | 全パネル共通 |
| `POST /api/change-password` | `{current_password,new_password}` | `{message}` | `200` | `401`, `422`, `500` | `.admin_credentials` | `.admin_credentials`, memory session | `changePassword()` | パスワード変更 |
| `GET /api/access-log` | query `{limit,offset}` | `{log}` | `200` | `401`, `422` | `.access_log` | none | `getAccessLog()` | アクセスログ |
| `GET /api/api-access-log` | query `{limit,offset,method?,path?,status?}` | `{log,total}` | `200` | `401`, `422`, `500` | `.api_access_log` | none | `getApiAccessLog()` | アクセスログ |
| `GET /api/sessions` | none | `{sessions}` | `200` | `401` | memory session | none | `getSessions()` | セッション管理 |
| `POST /api/sessions/revoke-all` | none | `{message,revoked_count}` | `200` | `401` | memory session | memory session, `.access_log` | `revokeAllSessions()` | セッション管理 |
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
| `POST /api/hooks` | `{phase,command_args,abort_on_failure}` | `HookRecord` | `201` | `401`, `422`, `500` | `.hooks` | `.hooks`, `.config_log` | `addHook()` | フック |
| `DELETE /api/hooks/{id}` | path `{id}` | `{message}` | `200` | `401`, `404`, `500` | `.hooks` | `.hooks`, `.config_log` | `deleteHook(id)` | フック |
| `GET /api/hooks/{id}/log` | path `{id}` | `{id,runs}` | `200` | `401`, `404`, `500` | `.build_logs/{build_id}_hook_{id}.json` | none | `getHookLog(id)` | フック |
| `GET /api/alert-rules` | none | `{rules}` | `200` | `401`, `500` | `.alert_rules` | none | `getAlertRules()` | 設定 |
| `POST /api/alert-rules` | `{metric,operator,threshold,level,message}` | `AlertRule` | `201` | `401`, `422`, `500` | `.alert_rules` | `.alert_rules`, `.config_log` | `addAlertRule()` | 設定 |
| `DELETE /api/alert-rules/{id}` | path `{id}` | `{message}` | `200` | `401`, `404`, `500` | `.alert_rules` | `.alert_rules`, `.config_log` | `deleteAlertRule(id)` | 設定 |
| `GET /api/tag-rules` | none | `{rules}` | `200` | `401`, `500` | `.tag_rules` | none | `getTagRules()` | 設定 |
| `POST /api/tag-rules` | `{condition,tags}` | `TagRule` | `201` | `401`, `422`, `500` | `.tag_rules` | `.tag_rules`, `.config_log` | `addTagRule()` | 設定 |
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
| `GET /api/tokens` | none | `{tokens}` | `200` | `401`, `500` | `.api_tokens` | none | `getTokens()` | API トークン管理 |
| `POST /api/tokens` | `{label,scope}` | `TokenCreateResult` | `201` | `401`, `422`, `500` | `.api_tokens` | `.api_tokens`, `.access_log` | `createToken(label,scope)` | API トークン管理 |
| `DELETE /api/tokens/{id}` | path `{id}` | `{message}` | `200` | `401`, `404`, `500` | `.api_tokens` | `.api_tokens`, `.access_log` | `revokeToken(id)` | API トークン管理 |

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
| token id | `tok{YYYYMMDDHHmmss}` | `tok20260915100500` | `tok20260915100500-001` |
| hook id | `h{YYYYMMDDHHmmss}` | `h20260915100500` | `h20260915100500-001` |
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
| destructive delete | `DELETE /api/queue`, `DELETE /api/snapshots/{id}`, `DELETE /api/hooks/{id}`, `DELETE /api/alert-rules/{id}`, `DELETE /api/tag-rules/{id}`, `DELETE /api/tokens/{id}` | path / auth 検証 → 対象存在確認 → 削除または失効 → audit log | `{message}` と件数がある場合は件数を返す。 | 対象不在は `404`。部分削除は禁止し、失敗時は `500`。 |
| external check | `POST /api/pat-verify`, `GET /api/rate-limit`, `GET /api/diagnostics`, `POST /api/smtp-test`, `POST /api/notify-test` | 設定読込 → timeout 付き外部確認 → 結果 response → 必要時 log 追記 | 確認結果を保存しない。ただし test 送信 log は仕様どおり追記する。 | 未設定は `501` または endpoint 固有 `422`。timeout は `500`。 |
| binary response | `GET /api/snapshots/{id}/download` | path 検証 → snapshot 存在確認 → archive stream | `Content-Type` と `Content-Disposition` を付与する。 | 不在は `404`。読込失敗は `500`。 |
| stream response | `GET /api/build/stream` | 認証 → 最新 / 実行中 log 特定 → SSE header → frame 送信 | `log` frame 後、必ず `end` frame を送って close する。 | log 不在は `404`。送信中断時は状態ファイルを更新しない。 |

`.config_log`、`.access_log`、`.notify_log` への追記は JSON Lines 1 行単位で行う。追記失敗時は対象 endpoint の副作用が既に完了している場合でも、失敗を `500` として返し、次回 GET で破損行を無視できる形式を維持する。追記行の末尾改行を書けなかった場合は、その行を破損行として扱う。

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
| `GET` | `/api/notes` | 要 | 運用ノートを返す |
| `POST` | `/api/notes` | 要 | 運用ノートを保存する |
| `GET` | `/api/smtp-config` | 要 | SMTP 設定を返す |
| `POST` | `/api/smtp-config` | 要 | SMTP 設定を保存する |
| `POST` | `/api/smtp-test` | 要 | SMTP テスト送信を行う |
| `GET` | `/api/queue` | 要 | ビルドキュー状態を返す |
| `DELETE` | `/api/queue` | 要 | 待機中ビルドキューを削除する |
| `GET` | `/api/dashboard-layout` | 要 | ダッシュボードウィジェット設定を返す |
| `POST` | `/api/dashboard-layout` | 要 | ダッシュボードウィジェット設定を置換する |
| `GET` | `/api/tokens` | 要 | API token 一覧を返す。token 本体は返さない |
| `POST` | `/api/tokens` | 要 | API token を発行する。token 本体は作成時のみ返す |
| `DELETE` | `/api/tokens/{id}` | 要 | API token を失効する |

**`POST /api/login` リクエスト / レスポンス：**
```json
// リクエスト
{ "password": "admin" }

// レスポンス
{ "token": "<session_token>", "must_change": "prompt" }
```

`must_change` の有効値：`"none"`（変更不要）| `"prompt"`（促す：初回ログイン時）| `"forced"`（強制：5 回目以降。変更完了まで管理画面の操作を制限）

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
  "on": ["failure"],
  "summary": { "enabled": false, "interval": "weekly", "hour": 9, "day_of_week": 1 },
  "email": { "enabled": false, "to": [], "on": [] }
}
```

`secret`：Webhook 署名シークレット。未設定時は `null`、設定済み時は `"***"`（マスク）を返す（→ 16E 参照）。

`on` の有効値：`"start"`（ビルド開始時）| `"success"`（ビルド成功時）| `"failure"`（ビルド失敗時）| `"weekly_summary"`（定期サマリー送信時）。複数指定可。

`summary`：定期サマリー通知の設定。`enabled: true` のとき指定スケジュールで統計サマリーを Webhook 送信する。`interval` の有効値：`"daily"` | `"weekly"`。`hour` は 0〜23（UTC）。`day_of_week` は `"weekly"` 時のみ有効（0 = 日曜〜6 = 土曜）。

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
{ "log_max_lines": 500, "history_max_count": 100, "build_timeout_seconds": 300, "log_retention_days": 30, "log_archive_after_days": 0, "log_level": "INFO", "pat_expires_at": null, "snapshots_keep": 5, "queue_max_size": 3, "build_retry_max": 0, "build_retry_base_seconds": 5, "commit_status_enabled": false, "commit_status_context": "Adlaire CI", "commit_status_target_url": null }
```

`pat_expires_at`：PAT の有効期限日（`YYYY-MM-DD` 形式）。`null` = 未設定。`GET /api/diagnostics` の `pat` 項目で 7 日以内なら `"warn"`、期限当日以前なら `"error"` に変更。

`POST /api/config` で更新可能なキーは `log_max_lines`、`history_max_count`、`build_timeout_seconds`、`log_retention_days`、`log_archive_after_days`、`log_level`、`pat_expires_at`、`snapshots_keep`、`queue_max_size`、`build_retry_max`、`build_retry_base_seconds`、`commit_status_enabled`、`commit_status_context`、`commit_status_target_url` に限定する。未知キーを含む場合は `422` を返し、既存設定を変更しない。

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
    { "id": "tok001", "label": "監視用", "scope": "read", "created_at": "2026-09-15T10:00:00Z", "last_used_at": "2026-09-15T11:00:00Z" }
]}
```

**`POST /api/tokens` リクエスト / レスポンス：**
```json
// リクエスト
{ "label": "監視用", "scope": "read" }
// レスポンス: 201
{ "id": "tok001", "token": "act_...", "label": "監視用", "scope": "read", "created_at": "2026-09-15T10:00:00Z" }
```

`token` はレスポンス時のみ返却し、以後は取得不可。`scope` の有効値：`"read"`（読み取り専用）。読み取り専用トークンは `GET` 系エンドポイントのみ許可し、`POST` / `DELETE` 系は `403 Forbidden` を返す。トークンは `Authorization: Bearer <token>` ヘッダーで送信する。

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

---

### ビルドフック（14E）

ビルド実行の直前（`pre`）・直後（`post`）に事前登録したコマンド引数配列を実行する。フック設定は `.hooks` に保存する。外部入力文字列をシェルへ渡す実装は禁止し、Go 標準ライブラリ `os/exec` の `exec.CommandContext(args[0], args[1:]...)` で実行する。

- `pre` フックが失敗（`exit_code != 0`）し `abort_on_failure: true` の場合、ビルドを中断しステータスを `hook_error` とする。
- `post` フックは `abort_on_failure` 設定に関わらずビルド結果（`success` / `failure`）を変更しない。
- フックの実行ログは `.build_logs/{build_id}_hook_{id}.json` に保存する。

**`GET /api/hooks` レスポンス例：**
```json
{ "hooks": [
    { "id": "h001", "phase": "pre",  "command_args": ["echo", "build start"], "enabled": true, "abort_on_failure": true },
    { "id": "h002", "phase": "post", "command_args": ["echo", "build end"],   "enabled": true, "abort_on_failure": false }
]}
```

**`POST /api/hooks` リクエスト / レスポンス：**
```json
// リクエスト
{ "phase": "pre", "command_args": ["echo", "build start"], "abort_on_failure": true }
// レスポンス: 201
{ "id": "h001", "phase": "pre", "command_args": ["echo", "build start"], "enabled": true, "abort_on_failure": true }
```

`phase` の有効値は `"pre"` または `"post"`。`command_args[0]` は絶対パス、または `PATH` 解決可能なコマンド名とする。`command_args` に空文字、NUL 文字、改行を含めてはならない。

**`DELETE /api/hooks/{id}` レスポンス：**
```json
{ "message": "Hook deleted" }
```

**`GET /api/hooks/{id}/log` レスポンス例：**
```json
{ "id": "h001", "runs": [
    { "build_id": "b20260915100000", "ran_at": "2026-09-15T10:00:00Z", "exit_code": 0, "output": "build start\n" },
    { "build_id": "b20260914183000", "ran_at": "2026-09-14T18:30:00Z", "exit_code": 1, "output": "Error: command not found\n" }
]}
```
`runs` は直近 20 件を返す（新しい順）。

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

---

### 自動タグ付けルール（15B）

ビルド完了時に条件式を評価し、マッチしたルールのタグを `.build_history` のエントリへ自動追記する（手動タグと共存する）。`.tag_rules` に保存する。

条件式で使用可能な変数：`status`（`"success"` / `"failure"`）/ `duration_seconds`（整数）/ `trigger`（`"polling"` / `"force_interval"` / `"manual"` / `"webhook"` / `"retry_pending_transfer"` / `"startup_config_integrity"` / `"rollback"`）
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

**`GET /api/notify-config` への追加（`email` セクション）：**
```json
{
  "webhooks": [ { "url": "...", "label": "メイン", "enabled": true, "payload_template": null, "retry_count": 2, "retry_interval_seconds": 30, "secret": "***" } ],
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
    { "id": "q001", "trigger": "manual", "queued_at": "2026-09-15T10:01:00Z" },
    { "id": "q002", "trigger": "webhook", "queued_at": "2026-09-15T10:02:00Z" }
  ], "max_size": 3 }
```
キューが空の場合：`{ "queued": [], "max_size": 3 }`

**`DELETE /api/queue` レスポンス：**
```json
{ "message": "Queue cleared", "cleared_count": 2 }
```
実行中のビルドは停止しない（`POST /api/build/cancel` を別途使用する）。

`POST /api/config` に `queue_max_size`（整数、`0` = キューなし）追加。

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

`trigger` の有効値：§13 の固定値（`"polling"`、`"force_interval"`、`"manual"`、`"webhook"`、`"retry_pending_transfer"`、`"startup_config_integrity"`、`"rollback"`）

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

  login(password)                               // POST /api/login → {token, must_change}; this._token にセット
  logout()                                      // POST /api/logout; this._token をクリア
  changePassword(currentPassword, newPassword)  // POST /api/change-password

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
  getDiagnostics()             // GET /api/diagnostics       → Promise<DiagnosticsObject>
  getRateLimit()               // GET /api/rate-limit        → Promise<RateLimitObject>
  getDiskUsage()               // GET /api/disk-usage        → Promise<DiskUsageObject>
  getConfigLog()               // GET /api/config-log        → Promise<{log: ConfigLogRecord[]}>
  getBranchConfig()            // GET /api/branch-config     → Promise<{source: string, branches: BranchTargetRecord[]}>
  setBranchConfig(branches)    // POST /api/branch-config    → Promise<{message: string, branches_count: number}>
  notifyWeeklySummary()        // POST /api/notify/weekly-summary → Promise<{message: string, period: string, success_count: number, failure_count: number, success_rate: number}>
  getWebhookEvents(limit = 50, offset = 0) // GET /api/webhook-events?limit={limit}&offset={offset} → Promise<{events: WebhookEventRecord[], total: number}>
  getWebhookConfig()          // GET /api/webhook-config    → Promise<{configured: boolean}>
  setWebhookConfig(secret)    // POST /api/webhook-config   → Promise<{message: string}>
  getHistoryComment(id)        // GET /api/history/{id}/comment  → Promise<CommentObject>
  setHistoryComment(id, comment) // POST /api/history/{id}/comment → Promise<{message: string}>
  setRepoConfig(config)        // POST /api/repo-config      → Promise<{message: string}>
  exportHistory()              // GET /api/history/export    → Promise<ExportObject>（JSON）
  setHistoryFlag(id, flagged)  // POST /api/history/{id}/flag → Promise<{message: string}>
  setHistoryTags(id, tags)     // POST /api/history/{id}/tags → Promise<{message: string}>
  rollbackHistory(id)          // POST /api/history/{id}/rollback → Promise<{message: string, build_id: string}>
  getTokens()                  // GET /api/tokens            → Promise<{tokens: TokenRecord[]}>
  createToken(label, scope = 'read') // POST /api/tokens     → Promise<TokenCreateResult>
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
  addHook(phase, commandArgs, abortOnFailure = true) // POST /api/hooks → Promise<HookRecord>
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
  setNotes(content)            // POST /api/notes            → Promise<{message: string, updated_at: string}>
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

全メソッドは `Promise` を返す。`streamBuild` は SSE 接続確立後に `StreamHandle` で resolve し、接続前エラーは `AdlaireCIError` で reject する。HTTP エラー（4xx / 5xx）は `AdlaireCIError` としてスローする。`401` 受信時はセッション期限切れとして `this._token` をクリアする。`constructor` を除く合計は 96 メソッド。

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

**SDK 引数変換契約：**

SDK method は、下表の通りに引数を path、query、body へ変換する。下表にない引数、既定値、body key を追加してはならない。

| SDK method | 引数 | 変換先 | 送信値 |
|------------|------|--------|--------|
| `login(password)` | `password` | body | `{password}` |
| `changePassword(currentPassword,newPassword)` | `currentPassword`, `newPassword` | body | `{current_password: currentPassword, new_password: newPassword}` |
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
| `setHistoryComment(id,comment)` | `id`, `comment` | path/body | path `{id}`、body `{comment}` |
| `setHistoryFlag(id,flagged)` | `id`, `flagged` | path/body | path `{id}`、body `{flagged}` |
| `setHistoryTags(id,tags)` | `id`, `tags` | path/body | path `{id}`、body `{tags}` |
| `createToken(label,scope)` | `label`, `scope="read"` | body | `{label,scope}` |
| `addHook(phase,commandArgs,abortOnFailure)` | `phase`, `commandArgs`, `abortOnFailure=true` | body | `{phase,command_args: commandArgs, abort_on_failure: abortOnFailure}` |
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
| `NotifyConfig` | `webhooks`, `on`, `summary`, `email` | `webhooks[].payload_template`, `webhooks[].secret` | `webhooks`, `on`, `email.to`, `email.on` | `GET /api/notify-config` |
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
| `TokenRecord` | `id`, `label`, `scope`, `created_at`, `last_used_at` | `last_used_at` | なし | `GET /api/tokens` |
| `TokenCreateResult` | `id`, `token`, `label`, `scope`, `created_at` | なし | なし | `POST /api/tokens` |
| `HookRecord` | `id`, `phase`, `command_args`, `enabled`, `abort_on_failure` | なし | `command_args` | `GET/POST /api/hooks` |
| `HookRunRecord` | `build_id`, `ran_at`, `exit_code`, `output` | なし | なし | `GET /api/hooks/{id}/log` |
| `AlertRule` | `id`, `metric`, `operator`, `threshold`, `level`, `message` | なし | なし | `GET/POST /api/alert-rules` |
| `TagRule` | `id`, `condition`, `tags` | なし | `tags` | `GET/POST /api/tag-rules` |
| `PipelineConfig` | `extra_args`, `env` | なし | `extra_args` | `GET /api/pipeline-config` |
| `SmtpConfig` | `host`, `port`, `user`, `tls`, `from`, `to`, `on`, `enabled`, `password_set` | `host`, `user`, `from` | `to`, `on` | `GET /api/smtp-config` |
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
| ログイン | `panel-login` | `login` | `form-login` | `password` | `btn-login` |
| パスワード変更 | `panel-password` | `password` | `form-password` | `current_password`, `new_password` | `btn-change-password` |
| ステータス | `panel-status` | `status` | なし | なし | `btn-refresh-status`, `btn-save-dashboard-layout` |
| 手動実行 | `panel-build` | `build` | なし | なし | `btn-build`, `btn-build-force`, `btn-build-cancel`, `btn-stream-close`, `btn-queue-clear` |
| ログビューア | `panel-logs` | `logs` | `form-log-search` | `n`, `q`, `from`, `to`, `level` | `btn-load-logs`, `btn-search-logs`, `btn-export-logs`, `btn-cleanup-logs`, `btn-archive-logs` |
| ビルド履歴 | `panel-history` | `history` | `form-history-filter` | `page`, `per_page`, `trigger`, `tag`, `flagged` | `btn-export-history` |
| システム情報 | `panel-system` | `system` | `form-pat` | `token`, `pat_expires_at` | `btn-pat-verify`, `btn-pat-update` |
| 通知設定 | `panel-notify` | `notify` | `form-notify` | `webhooks`, `on`, `summary`, `email`, `secret`, `smtp_password` | `btn-save-notify`, `btn-notify-test`, `btn-weekly-summary`, `btn-save-webhook-secret`, `btn-save-smtp`, `btn-smtp-test` |
| 設定 | `panel-config` | `config` | `form-config` | `log_max_lines`, `history_max_count`, `build_timeout_seconds`, `log_retention_days`, `log_archive_after_days`, `log_level`, `queue_max_size`, `snapshots_keep`, `build_retry_max`, `build_retry_base_seconds`, `commit_status_enabled`, `commit_status_context`, `commit_status_target_url` | `btn-save-config`, `btn-validate-config`, `btn-set-log-level` |
| アクセスログ | `panel-access-log` | `access-log` | `form-api-access-log-filter` | `limit`, `offset`, `method`, `path`, `status` | `btn-load-access-log`, `btn-load-api-access-log` |
| 統計 | `panel-stats` | `stats` | なし | なし | `btn-load-stats` |
| リポジトリ情報 | `panel-repo` | `repo` | `form-repo` | `owner`, `repo`, `branch`, `target_file`, `interval_seconds`, `allowed_from`, `allowed_to`, `maintenance_reason` | `btn-save-repo`, `btn-save-branch-config`, `btn-pause-schedule`, `btn-resume-schedule`, `btn-enable-maintenance`, `btn-disable-maintenance` |
| セッション管理 | `panel-sessions` | `sessions` | なし | なし | `btn-load-sessions`, `btn-revoke-sessions` |
| システム診断 | `panel-diagnostics` | `diagnostics` | なし | なし | `btn-run-diagnostics`, `btn-verify-output` |
| ビルド比較 | `panel-compare` | `compare` | `form-compare` | `left_build_id`, `right_build_id` | `btn-compare-builds` |
| API トークン管理 | `panel-tokens` | `tokens` | `form-token` | `label`, `scope` | `btn-create-token` |
| 運用ノート | `panel-notes` | `notes` | `form-notes` | `content` | `btn-save-notes` |
| スナップショット | `panel-snapshots` | `snapshots` | なし | なし | `btn-load-snapshots` |
| メンテナンス | `panel-maintenance` | `maintenance` | `form-maintenance` | `reason` | `btn-maintenance-enable`, `btn-maintenance-disable` |
| アクセス制御 | `panel-access-control` | `access-control` | `form-access-control` | `allow` | `btn-save-access-control` |
| フック | `panel-hooks` | `hooks` | `form-hook` | `phase`, `command_args`, `abort_on_failure` | `btn-add-hook` |

共通領域の DOM id は、`app-root`、`nav-panels`、`global-banner`、`global-error`、`global-success`、`maintenance-banner`、`build-log-stream`、`build-queue-summary`、`issued-token-once` とする。エラー表示要素は各 panel 内に `id="{section-id}-error"`、成功表示要素は `id="{section-id}-success"` を置く。`label[for]` と input `id` は `field-{field_name}` 形式で一致させる。複数行・配列入力は `textarea` または table row で表現し、保存直前に SDK 引数の型へ変換する。

**画面構成：**

| パネル | 表示内容 | 表示条件 |
|-------|---------|---------|
| ログイン | パスワード入力フォーム | 未ログイン時のみ |
| パスワード変更 | 現在・新パスワード入力フォーム | `must_change: "prompt"` または `"forced"` 時（`"forced"` 時は他パネル非表示） |
| ステータス | 最終ビルド時刻・SHA・成否・出力サイトリンク・ダッシュボードウィジェット編集モード（表示するウィジェットをチェックボックスで選択・並び替え・保存） | ログイン済み |
| 手動実行 | ビルドトリガーボタン・SHA リセットを含む強制ビルドボタン・キャンセルボタン（実行中のみ有効）・実行結果表示・リアルタイムログ表示エリア（SSE ストリーミング）・キュー状態表示（待機中件数・クリアボタン） | ログイン済み |
| ログビューア | 最新ビルドログ（n 行・キーワードフィルター・ログレベルフィルターボタン（INFO / WARNING / ERROR / DEBUG）・JSON エクスポートボタン・横断検索フォーム（期間指定）・検索結果一覧） | ログイン済み |
| ビルド履歴 | 過去ビルド一覧（日時・SHA・成否・トリガー種別・所要時間・重要フラグ列・タグ列・ログ表示リンク・コメント入力欄）・フラグ付きのみ表示フィルター・タグフィルター・ページネーション UI・JSON エクスポートボタン | ログイン済み |
| システム情報 | 出力サイトサイズ・更新日時・稼働時間・ディスク使用量（ログ合計・出力サイト）・PAT 即時検証ボタン・PAT 更新フォーム・PAT 有効期限表示（設定フォーム・期限切れ間近で警告表示）・GitHub API レート制限表示 | ログイン済み |
| 通知設定     | Webhook 一覧（追加/削除/ラベル/有効無効切り替え/リトライ回数・間隔設定/シークレット入力欄）・通知条件設定（ビルド開始時・成功時・失敗時）・各 Webhook ペイロードテンプレート編集フォーム（変数一覧表示）・テスト送信ボタン・定期サマリー設定（間隔・時刻・曜日・即時送信ボタン）・送信履歴（試行回数・エラー内容列含む）・メール通知セクション（SMTP 設定フォーム・宛先リスト・通知条件・テスト送信ボタン） | ログイン済み |
| 設定         | ログ保持行数・履歴保持件数の設定変更・ビルドタイムアウト設定・ログレベル変更（INFO / DEBUG）・ログ保持期間（日数、0 = 無制限）・スナップショット保持世代数設定・ビルドキュー最大長設定・手動クリーンアップボタン・設定変更履歴（変更日時・項目・変更前後の値）・IP アクセス制限セクション（許可 IP / CIDR 一覧・追加フォーム・削除ボタン）・フック設定セクション（pre / post フック一覧・command_args 入力フォーム・実行ログリンク・有効無効切り替え）・アラートルール設定セクション（メトリクス・演算子・しきい値・レベル・メッセージの入力フォーム・ルール一覧・削除ボタン）・自動タグ付けルールセクション（条件式・タグ入力フォーム・ルール一覧・削除ボタン）・パイプライン設定セクション（追加引数入力欄・環境変数テーブル） | ログイン済み |
| アクセスログ | ログイン履歴（日時・成否）・API アクセスログ（method、path、status、duration、actor、remote_addr）・API アクセスログフィルター | ログイン済み |
| 統計         | ビルド回数・成功率・平均間隔・平均・最大ビルド時間・日別時系列データ（グラフ表示対応） | ログイン済み |
| リポジトリ情報 | 監視対象リポジトリ・ブランチ・ファイルの確認・設定変更フォーム（OWNER / REPO / BRANCH / TARGET_FILE）・ポーリング間隔変更フォーム・ポーリング一時停止／再開ボタン・許可時間帯設定（from〜to、解除ボタン）・メンテナンスモード有効化フォーム（理由入力）・解除ボタン・現在の状態表示 | ログイン済み |
| セッション管理 | 有効セッション一覧・全セッション強制無効化ボタン | ログイン済み |
| システム診断   | PAT・GitHub API・出力サイト・systemd・Webhook の診断項目一覧（ok / warn / error）・診断実行ボタン・アラートバッジ（`GET /api/dashboard` の `alerts` に基づき warn / error を表示）・出力整合性チェック項目（`POST /api/verify-output` 結果表示）・メンテナンスモード中はバナーを全パネル上部に表示 | ログイン済み |
| ビルド比較     | ビルド履歴から 2 件を選択・ログ並列表示・差分ハイライト（クライアントサイド処理） | ログイン済み |
| API トークン管理 | 発行済みトークン一覧（ラベル・スコープ・作成日時・最終使用日時）・新規発行フォーム（ラベル入力・スコープ選択）・発行時のみトークン文字列を表示・失効ボタン | ログイン済み |
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
| アクセスログ | `getAccessLog()`, `getApiAccessLog({limit:100,offset:0})` | ログなしを 1 行で表示する。 |
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
| ログイン | ログイン | `login(password)` | 通常画面へ遷移 | `getDashboard()`, `getStatus()`, `getQueue()` | password 空欄、送信中 |
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
| アクセスログ | 取得 | `getAccessLog()` | 件数を表示 | なし | 読み込み中 |
| アクセスログ | API アクセスログ取得 | `getApiAccessLog({limit,offset,method,path,status})` | 件数を表示 | なし | 読み込み中 |
| 統計 | 取得 | `getStats()`, `getStatsTimeline()`, `getStatsBuildDuration()` | グラフと数値を更新 | なし | 読み込み中 |
| リポジトリ情報 | repo 保存 | `setRepoConfig(config)` | `Repo config updated` | `getRepoInfo()`, `getConfigLog()` | 入力不正、送信中 |
| リポジトリ情報 | branch config 保存 | `setBranchConfig(branches)` | `Branch config updated` | `getBranchConfig()`, `getConfigLog()` | 入力不正、送信中 |
| リポジトリ情報 | schedule 変更 | `setScheduleInterval()`, `pauseSchedule()`, `resumeSchedule()`, `setAllowedHours()`, `clearAllowedHours()`, `setForceInterval()`, `setBuildCooldown()` | schedule 変更結果を表示 | `getSchedule()`, `getConfigLog()` | 入力不正、送信中 |
| セッション管理 | 全 revoke | `revokeAllSessions()` | revoke 件数を表示 | `getSessions()` | 送信中 |
| システム診断 | 診断実行 | `getDiagnostics()` | 診断結果を表示 | なし | 読み込み中 |
| システム診断 | 出力 checksum 検証 | `verifyOutput()` | match 結果を表示 | なし | 送信中 |
| API トークン管理 | token 発行 | `createToken(label,scope)` | token 本体を 1 回だけ表示 | `getTokens()` | label 空欄、送信中 |
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
| 秘密情報入力 | PAT、Webhook Secret、SMTP password、発行直後 token は画面遷移、成功表示、再取得後にフォーム値から消去する。 |
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

秘密情報 field は、`password`、`current_password`、`new_password`、`token`、`secret`、`smtp_password`、`issued-token-once` とする。これらは成功、失敗、画面遷移、`401`、`logout()`、`revokeAllSessions()` のいずれの場合も DOM 値を空にする。発行直後 token は `issued-token-once` に 1 回だけ表示し、次の任意の user action で消去する。

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

**認証情報ファイル形式（JSON）：**
```json
{
  "password_hash": "<sha256_iter_v1_hex>",
  "salt": "<hex>",
  "algorithm": "sha256_iter_v1",
  "iterations": 260000,
  "login_count": 0,
  "updated_at": "2026-09-15T10:00:00Z"
}
```

**ハッシュアルゴリズム：** Go 標準ライブラリのみで実装する `sha256_iter_v1`

| 項目 | 仕様 |
|------|------|
| salt 生成 | `crypto/rand` で 32 bytes を生成し、`encoding/hex` で 64 文字の hex 文字列として保存する。 |
| 初回 digest | `sha256(salt_bytes || password_utf8_bytes)` |
| 反復 | `iterations = 260000`。2 回目以降は `sha256(previous_digest || salt_bytes || password_utf8_bytes)` を繰り返す。 |
| 保存値 | 最終 digest を lowercase hex 文字列で `password_hash` に保存する。 |
| 比較 | 入力パスワードから同一手順で digest を生成し、`crypto/subtle.ConstantTimeCompare` で比較する。 |

`golang.org/x/crypto/pbkdf2` 等の外部パッケージは使用しない。PBKDF2、bcrypt、Argon2 等は本ファイルでは実装値を定義しない。

**セッショントークン生成：**
`crypto/rand` で 32 bytes を生成し、`encoding/hex` で 64 文字の lowercase hex 文字列へ変換する。

**セッション管理：** `components/api.go` 内のインメモリ辞書で管理。有効期限 8 時間。再起動で全セッション破棄。同一ユーザーの複数同時セッションを許容する。辞書 key は token 本体ではなく `sha256(token)` の lowercase hex とし、API response、`.access_log`、サーバーログへ token 本体を出力してはならない。

**セッション期限切れ時：** `401 Unauthorized` を返す。クライアント（SDK）は `this._token` をクリアし、再ログインを促す。

**パスワード入力制約：**

| 項目 | 仕様 |
|------|------|
| 最小長 | 8 文字 |
| 最大長 | 128 文字 |
| 許可文字 | UTF-8 文字列。NUL 文字は禁止。前後空白はトリムせず、入力値そのものを検証・ハッシュ化する。 |
| 初期パスワード | `admin`。初回ログイン時は `must_change: "prompt"` を返す。 |
| 変更時検証 | `new_password` が現在パスワードと同一の場合は `422` を返す。 |
| 失敗時応答 | パスワード不一致は `401` と `{"error":"Unauthorized"}` を返し、どの条件に失敗したかは返さない。 |
| 成功時保存 | 新 salt、新 hash、`login_count: 0`、`updated_at` を原子的に保存する。 |

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

**ログインフロー：**
```
POST /api/login
  ├─ 連続失敗ロック中 → 429
  └─ パスワードハッシュ検証
       ├─ 失敗 → 連続失敗回数 + 1 → .access_log 追記 → 401
       └─ 成功 → login_count + 1 → ファイル更新
                  → 連続失敗回数を 0 へリセット
                  ├─ login_count == 1 → must_change: "prompt"（促す）
                  ├─ login_count >= 5 → must_change: "forced"（強制）
                  └─ それ以外      → must_change: "none"
                  → セッショントークン生成・返却
```

**パスワード変更時：** `login_count` を 0 にリセット。新しい salt を生成しハッシュを更新。変更完了後に現セッション以外のセッションを破棄。

**`--init-credentials` オプション：** `components/api.go` を `--init-credentials` 引数で起動した場合、初期パスワード `admin` で `.admin_credentials` を生成して終了する（HTTP サーバーは起動しない）。

`.admin_credentials` が既に存在する場合、`--init-credentials` は上書きせず `409` 相当の終了コード `2` で終了し、標準エラーへ `credentials already exist` を出力する。初期化成功時の終了コードは `0` とする。

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

Phase 完了判定テンプレートは以下とする。実装 PR 本文では、対象 Phase ごとに本テンプレートの項目を埋める。

```text
## Phase 完了判定

- 対象 Phase:
- 対象コンポーネント:
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

本節は、将来計画から実装可能仕様へ昇格した機能の詳細仕様である。本節に定義された機能は、§0i、§12、§13、§15、§22、§23、§24 と同時に満たす。

### 27.1 GitHub Commit Status API

runner は `.server_config.commit_status_enabled == true` の場合、対象 commit に GitHub Commit Status API を送信する。送信先は `POST /repos/{owner}/{repo}/statuses/{sha}` とし、`sha` は対象 target の commit SHA を使用する。blob SHA しか取得できない場合は status を送信せず、`.build_logs/{id}.json.commit_status.error="commit sha unavailable"` を保存する。

送信 payload は次に固定する。

| key | 値 |
|-----|----|
| `state` | 開始時 `"pending"`、成功時 `"success"`、失敗時 `"failure"`、設定/状態書込など CI 自体の異常時 `"error"`。 |
| `context` | `.server_config.commit_status_context`。既定値 `"Adlaire CI"`。 |
| `description` | 140 文字以内。開始時 `Build started`、成功時 `Build succeeded`、失敗時 `Build failed: <target_status>`、pending deploy 時 `Build succeeded with deploy pending`。 |
| `target_url` | `.server_config.commit_status_target_url` が `null` でなければ送信する。 |

処理順序は以下とする。

1. commit SHA を確定する。
2. build id を採番する。
3. `.build_status.json status="running"` を書く前に GitHub status `"pending"` を送信する。
4. pipeline / deploy / snapshot / history の最終結果確定後、GitHub status の最終 state を送信する。
5. `.build_logs/{id}.json.commit_status` に最終送信結果を保存する。
6. `.build_history.commit_status_state` に最終 state を保存する。

Commit Status 送信失敗は build 成否を反転させない。送信失敗時は WARN ログ `COMMIT_STATUS_FAILED: status=<http_status> error=<reason>` を出し、`.build_logs/{id}.json.commit_status.state` を最後に送信しようとした state、`error` を固定文言で保存する。GitHub API 認証失敗、403、404、5xx、network error はすべて送信失敗として扱う。

検証条件:

| ケース | 期待結果 |
|--------|----------|
| 有効・成功 build | fake GitHub server に `pending` → `success` の順で 2 回送信される。 |
| 有効・pipeline 失敗 | `pending` → `failure` の順で送信される。 |
| deploy pending | 最終 state は `success`、description は deploy pending を含む。 |
| commit SHA なし | status 送信なし、build は継続、log に `commit sha unavailable`。 |
| GitHub status 送信失敗 | build 成否は維持、WARN と `commit_status.error` を保存する。 |
| 無効 | status API 呼び出し 0 回、`commit_status.enabled=false`。 |

### 27.2 ドライラン実行モード

`adlaire-ci-runner --dry-run` は、実行計画を検証する読み取り専用モードである。dry-run は `.build_lock`、`.build_state`、`.build_status.json`、`.build_history`、`.build_logs/`、`.pending_transfers`、`.notify_*`、`.snapshots/`、GitHub Commit Status、deploy 先を変更してはならない。

dry-run の stdout は JSON object 1 件と末尾改行に固定する。

```json
{
  "dry_run": true,
  "state_dir": "/opt/adlaire-builder",
  "targets": [
    {
      "branch": "main",
      "target_file": "docs",
      "previous_sha": "old",
      "current_blob_sha": "new",
      "current_commit_sha": "abcdef",
      "would_build": true,
      "trigger": "polling",
      "reason": "sha_changed"
    }
  ],
  "would_write": [],
  "would_call": ["github_tree", "github_blob"],
  "errors": [],
  "warnings": []
}
```

`targets[].reason` は `"sha_changed"`、`"no_change"`、`"cooldown"`、`"circuit_open"`、`"config_error"`、`"github_error"`、`"precheck_error"` のいずれかとする。`would_write` は常に空配列でなければならない。

終了コードは、検証が完了した場合 `0`、設定不正または入力不正 `2`、GitHub API 全再試行失敗 `3` とする。`--help` / `--version` と同時指定された場合は `--help` / `--version` を優先し、dry-run JSON を出力しない。

検証条件:

| ケース | 期待結果 |
|--------|----------|
| SHA 差分あり | `would_build=true`、状態ファイル差分なし。 |
| SHA 差分なし | `would_build=false`、`reason="no_change"`。 |
| cooldown | `would_build=false`、`reason="cooldown"`。 |
| GitHub API 失敗 | 終了コード `3`、`errors[]` に理由、状態ファイル差分なし。 |
| 設定破損 | 終了コード `2`、破損ファイルの退避も再生成も行わない。 |

### 27.3 ビルド失敗時の自動リトライ

runner は `.server_config.build_retry_max > 0` の場合、retry 対象失敗だけを同一 build id 内で最大 `build_retry_max` 回追加試行する。総試行回数は `1 + build_retry_max` とする。

retry 対象は以下に限定する。

| 対象 | 条件 |
|------|------|
| GitHub API | network error、HTTP 429、HTTP 500〜599。 |
| pipeline | timeout。exit code 非 0 は retry しない。 |
| deploy | SSH 接続失敗、検証用 checksum 取得失敗、network timeout。checksum mismatch は retry しない。 |

retry 待機秒数は `build_retry_base_seconds * attempt` とする。初回失敗後の retry 1 回目は `base * 1`、retry 2 回目は `base * 2`。待機中に process が終了した場合、未完了 retry を再開してはならない。

attempt ごとの結果は `.build_logs/{id}.json.attempts[]` に必ず保存する。最終 attempt が成功した場合、`.build_history.status` は `"success"` とし、`retry_count` に追加 retry 回数を保存する。全 attempt 失敗時は `"failure"` とする。SHA 更新、snapshot、deploy 成功記録は最終成功時だけ行う。

検証条件:

| ケース | 期待結果 |
|--------|----------|
| API 429 後成功 | attempts 2 件、history success、retry_count 1。 |
| pipeline timeout 後成功 | attempts 2 件、SHA は最終成功後のみ更新。 |
| pipeline exit 1 | retry なし、history failure、retry_count 0。 |
| deploy checksum mismatch | retry なし、pending transfer 記録。 |
| retry 上限到達 | history failure、attempts は `1 + build_retry_max` 件。 |

### 27.4 出力サイトへのビルドメタ埋め込み

`adlaire-ci-build` は `--build-id`、`--commit-sha`、`--build-at` を受け取り、全 HTML ページの `<head>` に次の meta を必ず出力する。

```html
<meta name="adlaire-build-id" content="b20260916010000">
<meta name="adlaire-commit-sha" content="abcdef1">
<meta name="adlaire-build-at" content="2026-09-16T01:00:00Z">
```

値が空文字の場合も meta tag は出力し、`content=""` とする。HTML escape は attribute escape とし、`"`、`&`、`<`、`>` を escape する。runner が pipeline を起動する場合は、同じ値を CLI 引数または環境変数 `ADLAIRE_BUILD_ID`、`ADLAIRE_COMMIT_SHA`、`ADLAIRE_BUILD_AT` のいずれかで渡す。CLI 引数を使える場合は CLI 引数を優先する。

`[REPORT]` には `build_id`、`commit_sha`、`build_at` を追加する。`.build_logs/{id}.json.build_meta`、`GET /api/output-meta` は HTML meta と同じ値を返す。

検証条件:

| ケース | 期待結果 |
|--------|----------|
| 値あり | 全 HTML に 3 meta が同値で出力される。 |
| 値なし | 3 meta が `content=""` で出力される。 |
| 複数ページ | index と全ページで同じ build meta。 |
| output-meta | HTML meta、REPORT、build log、API response が一致する。 |

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

検証条件:

| ケース | 期待結果 |
|--------|----------|
| 正常値 | `valid=true`、errors 空、状態差分なし。 |
| 範囲外 | `valid=false`、該当 field error。 |
| 未知 key | HTTP 422、状態差分なし。 |
| secret 風 key | HTTP 422、response と log に値を含めない。 |

### 27.6 API アクセスログ

`components/api.go` は全 `/api/` request について `.api_access_log` へ JSON Lines を追記する。`GET /api/health` も対象とする。静的 file 配信、admin HTML、SDK JS は対象外とする。

追記タイミングは response status 確定後とする。追記失敗時は、対象 API の本来の response を優先し、サーバーログに `API_ACCESS_LOG_WRITE_FAILED` を出す。access log 書き込み失敗を理由に API response を `500` へ変更してはならない。

`GET /api/api-access-log` は `limit`、`offset`、`method`、`path`、`status` query を受け付ける。`limit` は 1〜1000、既定値 100。`offset` は 0 以上、既定値 0。`method` は大文字 HTTP method 完全一致。`path` は prefix match。`status` は HTTP status 完全一致。壊れた行は無視し、新しい順で返す。

検証条件:

| ケース | 期待結果 |
|--------|----------|
| 認証成功 API | actor `admin` または token id で記録される。 |
| 認証失敗 API | status 401、actor null、secret なし。 |
| Webhook | auth_type `webhook`、actor `webhook`。 |
| filter | method/path/status が完全に効く。 |
| 書込失敗 | 本来 response を維持し、server log に WARN。 |

### 27.7 ビルドログのアーカイブ圧縮

runner と `POST /api/logs/archive` は、`.server_config.log_archive_after_days > 0` の場合、対象日数より古い `.build_logs/{id}.json` を gzip 圧縮し、`.build_logs/archive/{id}.json.gz` へ保存する。圧縮成功後、元の `.build_logs/{id}.json` を削除する。`.build_logs/archive/` 内のファイルを再圧縮してはならない。

gzip は Go 標準ライブラリ `compress/gzip` を使用し、mtime は元ファイル mtime ではなく圧縮実行時刻でよい。圧縮前 JSON を読み込めないファイルは archive 対象外とし、WARN `LOG_ARCHIVE_SKIP_CORRUPT: id=<id>` を出す。実行中 build の `current_build_id` と一致する log は対象外とする。

ログ参照 API は通常ファイルを先に探し、存在しない場合に archive を探す。archive を読む場合は gzip 展開後に通常 `.build_logs/{id}.json` と同じ schema として扱う。`GET /api/logs/search`、`GET /api/history/{id}/log`、`GET /api/output-meta`、`GET /api/disk-usage` は archive を含めて動作する。

`POST /api/logs/cleanup` は、archive 済みファイルも `log_retention_days` の削除対象に含める。`POST /api/logs/archive` は `{ "message": "Logs archived", "archived_count": N }` を返す。

検証条件:

| ケース | 期待結果 |
|--------|----------|
| 対象ログあり | `.json.gz` 作成、元 `.json` 削除、API 参照可。 |
| 実行中ログ | archive しない。 |
| 破損ログ | archive しない、WARN、処理継続。 |
| archive API | 件数を返し、disk usage に archive bytes を含める。 |
| cleanup | 通常 log と archive log の両方を保持期間で削除する。 |

### 27.8 ビルドステータスファイル出力

本機能の目的は、runner の現在状態と直近結果を `.build_status.json` に集約し、API、SDK、UI、将来 MCP が同じ read-only 情報を参照できるようにすることである。

対象コンポーネントは `components/runner.go` と `components/api.go` とする。`components/runner.go` は `.build_status.json` の唯一の通常更新責務を持つ。`components/api.go` は `GET /api/status`、`GET /api/dashboard`、`GET /api/health` で read-only 参照する。API は `.build_status.json` を自動修復してはならない。

**入力：**

| 入力 | 説明 |
|------|------|
| runner 起動状態 | lock 取得、起動時整合性チェック結果、target 処理結果、pending transfer / notify 件数、circuit 状態。 |
| `.build_state` | `running`、`current_build_id`、`last_started_at`、`last_finished_at`、`queued`。 |
| `.build_history` | 直近 build id、status、trigger、duration。 |
| `.pending_transfers` | pending transfer 件数。 |
| `.notify_pending` | pending notify 件数。 |
| `.build_circuit_state` | circuit open 状態、連続失敗回数。 |

**出力：**

`.build_status.json` は §22.0a / §22.0c の schema に従う JSON object とする。文字コードは UTF-8、改行は末尾 1 つ、ファイル mode は `600` とする。更新は同一ディレクトリ一時ファイルへの書き込み、`fsync`、`os.Rename`、親ディレクトリ `fsync` の順で atomic write する。

**更新タイミング：**

| タイミング | 必須値 |
|------------|--------|
| 起動時整合性チェックで復旧または停止が発生した直後 | `status="warning"` または `"failure"`、`last_trigger="startup_config_integrity"`、`running=false`。 |
| build 開始前 | `status="running"`、`running=true`、`current_build_id`、`last_trigger`、`last_started_at` を保存する。 |
| 変更なし skip | `status="skipped"`、`last_target_status="skipped_no_change"`、`running=false`。build id は更新しない。 |
| cooldown skip | `status="skipped"`、`last_target_status="skipped_cooldown"`、`running=false`。 |
| build 成功 | `status="success"`、`last_build_id`、`last_finished_at`、`last_duration_seconds`、`output_sha256` を保存する。 |
| deploy pending | `status="success"`、`last_deploy_status="pending"`、`pending_transfers_count` を保存する。 |
| build 失敗 | `status="failure"`、`last_error`、`last_target_status`、`last_finished_at` を保存する。 |
| runner finalizer | `running=false`、`current_build_id=null` を必ず保存する。 |

`status` の許容値は `"none"`、`"running"`、`"success"`、`"failure"`、`"skipped"`、`"warning"` に固定する。`last_target_status` は §13 の `target_status` 値、または `null` とする。`last_deploy_status` は `"success"`、`"failure"`、`"pending"`、`"skipped"`、`"none"`、`null` のいずれかとする。

**異常系：**

| 条件 | 処理 |
|------|------|
| `.build_status.json` 書き込み失敗 | ERROR ログ `BUILD_STATUS_WRITE_FAILED: path={path} error={reason}` を出し、runner 終了コードを最低 `1` にする。build 成功後に発生した場合も終了コードは `1` とする。 |
| `.build_status.json` 破損を API が検出 | `GET /api/status` と `GET /api/dashboard` は `500` を返す。`GET /api/health` は `status="degraded"` を返し、破損を `checks[]` に含める。 |
| `.build_status.json` 不在 | API は `.build_state`、`.build_history`、`.build_lock` から後方互換値を算出して返してよい。ただしファイル作成はしない。 |
| pending 件数読取失敗 | 件数を `null` にせず `0` として返してはならない。status 書き込み時はエラー扱いにし、`last_error` に固定文言を保存する。 |

**セキュリティ：**

`.build_status.json` に GitHub PAT、Webhook URL secret、SMTP password、API token、session token、request body を保存してはならない。`last_error` は最大 500 文字に切り詰め、改行は `\n` 文字列へ escape する。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| build 開始 | `status="running"`、`running=true`、`current_build_id` が保存される。 |
| build 成功 | `status="success"`、`running=false`、`last_build_id` と `last_duration_seconds` が保存される。 |
| 変更なし | build log / history を作らず、`.build_status.json` は `skipped_no_change` を保持する。 |
| 起動時整合性復旧 | `last_trigger="startup_config_integrity"`、復旧内容が `last_error` または warning として確認できる。 |
| 書込失敗 | runner 終了コードが最低 `1`、ERROR ログが出る。 |
| API read | `GET /api/status`、`GET /api/dashboard` が `.build_status.json` を第一参照元にする。 |

### 27.9 ビルドトリガー種別の記録

本機能の目的は、runner がなぜ build または関連処理を開始したかを、履歴、ログ、状態、API、SDK、UI で同一の固定値として扱うことである。

対象コンポーネントは `components/runner.go`、`components/api.go`、`admin/adlaire-ci-sdk.js`、`admin/index.html` とする。`components/runner.go` は trigger の確定と永続化を担当し、API / SDK / UI は既存値の表示と filter のみを担当する。

**trigger 許容値：**

| 値 | 発生条件 | 補足 |
|----|----------|------|
| `polling` | systemd timer 等の通常起動で SHA 差分がある。 | 既定の自動ビルド。 |
| `force_interval` | SHA 差分なし、かつ `force_build_interval_hours` 条件を満たす。 | 手動 force には使わない。 |
| `manual` | `POST /api/build` または `POST /api/build/force` 由来の queue entry を処理する。 | force は `payload.force=true` で表す。 |
| `webhook` | `POST /api/webhook` 由来の queue entry を処理する。 | 署名検証成功済み event のみ。 |
| `retry_pending_transfer` | `.pending_transfers` の再送のみを実行する。 | 通常 build とは別 trigger。 |
| `startup_config_integrity` | 起動時整合性チェックで復旧、正規化、または停止が発生する。 | build log / history は作成しない。 |
| `rollback` | `POST /api/history/{id}/rollback` により snapshot を再転送する。 | 新しい build id を作成する。 |

上表以外の値を保存、返却、表示してはならない。特に `"auto"`、`"force"`、`"scheduled"`、`"timer"` は使用禁止とする。

**保存先：**

| 保存先 | 必須条件 |
|--------|----------|
| `.build_logs/{id}.json.trigger` | build log を作成する全処理で必須。 |
| `.build_history.trigger` | build history を追記する全処理で必須。 |
| `.build_status.json.last_trigger` | build、skip、復旧、rollback の最終 trigger を保存する。 |
| Queue entry `trigger` | `"manual"` または `"webhook"` のみ許可する。 |
| API response | `StatusObject.last_trigger`、`HistoryRecord.trigger`、`DashboardObject.status.last_trigger` で同じ値を返す。 |

**判定順序：**

1. 起動引数が `--help` または `--version` の場合、trigger を確定しない。
2. 起動時整合性チェックで復旧、正規化、停止が発生した場合、`startup_config_integrity` を `.build_status.json` に記録する。
3. queue entry が存在する場合、entry の `trigger` を採用する。
4. pending transfer の再送だけで終了する起動は `retry_pending_transfer` とする。
5. SHA 差分がある通常起動は `polling` とする。
6. SHA 差分がなく force interval 条件を満たす場合は `force_interval` とする。
7. rollback API が作成する処理は `rollback` とする。

複数条件が同時に成立した場合は、上記順序で最初に該当した trigger を採用する。1 回の runner 起動で複数 branch target を処理する場合、target ごとに同じ trigger を保存する。ただし queue entry が target を指定する場合は、対象 target のみにその trigger を適用する。

**API / SDK / UI：**

`GET /api/history` の `trigger` query は上表の値だけを受け付ける。不正値は `422` を返す。SDK `getHistory({trigger})` は値を変換せず送信する。UI は filter の選択肢を上表の 7 件に固定し、未知 trigger を受け取った場合は `Unknown` へ丸めず、該当行に `invalid trigger` エラーを表示する。

**異常系：**

| 条件 | 処理 |
|------|------|
| queue entry の trigger が不正 | queue entry を処理せず ERROR ログ `INVALID_TRIGGER: id={id} trigger={value}`、HTTP API 由来なら queue 作成時に `422`。 |
| 既存 history に未知 trigger がある | API はその行を返すが、`warnings[]` に `unknown_trigger` を含める。新規保存では未知値を禁止する。 |
| build log と history の trigger 不一致 | API は `500` を返し、server log に `TRIGGER_MISMATCH: id={id}` を出す。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| polling build | log、history、status が `polling` で一致する。 |
| manual force | queue は `manual`、payload は `force=true`、保存 trigger は `manual`。 |
| webhook | 署名検証成功時だけ `webhook` が保存される。 |
| force interval | SHA 差分なしで `force_interval` が保存される。 |
| filter | `GET /api/history?trigger=manual` が manual のみ返す。 |
| 不正 trigger | queue 作成または history filter が `422`。 |

### 27.10 設定ファイル起動時整合性チェック

本機能の目的は、runner が build 処理に入る前に、runner が読む状態ファイルの破損、型不一致、必須 key 不足、権限不備を検出し、規定どおり復旧または停止することである。

対象コンポーネントは `components/runner.go` のみとする。管理 API、SDK、UI は本機能の実行責務を持たない。API が同じ状態ファイルを読む場合も、起動時整合性チェックを代行してはならない。

**対象ファイル：**

| 順序 | ファイル | 不在時 | 破損時 | unknown key |
|------|----------|--------|--------|-------------|
| 1 | `.branch_config` | 作成せず default 採用 | backup 後、不在扱い | 除去して正規化 |
| 2 | `.notify_config` | 初期値作成 | backup 後、初期値作成 | 除去して正規化 |
| 3 | `.build_state` | 初期値作成 | backup 後、初期値作成 | 除去して正規化 |
| 4 | `.build_circuit_state` | 初期値作成 | backup 後、初期値作成 | 除去して正規化 |
| 5 | `.pending_transfers` | `[]` 作成 | backup 後、`[]` 作成 | 除去して正規化 |
| 6 | `.notify_pending` | `[]` 作成 | backup 後、`[]` 作成 | 除去して正規化 |

**実行順序：**

1. CLI 引数を検証する。
2. `--help` または `--version` の場合は本チェックを実行しない。
3. `--dry-run` の場合は検証だけ行い、backup、初期化、正規化、通知、状態更新を行わない。
4. `StateDir` が絶対パスかつ既存ディレクトリであることを確認する。
5. `.build_lock` を取得する。
6. 上表の順序で対象ファイルを検証する。
7. 復旧可能な問題は backup、初期化、正規化を行う。
8. 復旧結果に応じて `.build_status.json.last_trigger="startup_config_integrity"` を保存する。
9. 復旧通知条件を満たす場合は `config_corrupt` 通知を 1 回だけ送信する。
10. 停止条件がなければ pending retry、cooldown、target 処理へ進む。

**判定分類と停止条件：**

| 分類 | 処理 | runner 終了コード |
|------|------|------------------|
| `missing_optional` | `.branch_config` を作らず default 採用。 | 継続、最終結果に従う |
| `missing_required` | 初期値を atomic write。 | 継続、最終結果に従う |
| `parse_error` | corrupt backup 後、ファイル別復旧。 | 継続、最終結果に従う |
| `top_level_type_mismatch` | corrupt backup 後、ファイル別復旧。 | 継続、最終結果に従う |
| `required_key_missing` | corrupt backup 後、ファイル別復旧。 | 継続、最終結果に従う |
| `required_key_type_mismatch` | corrupt backup 後、ファイル別復旧。 | 継続、最終結果に従う |
| `invalid_value` | corrupt backup 後、ファイル別復旧。 | 継続、最終結果に従う |
| `unknown_key` | backup せず未知 key を除去して atomic write。 | 継続、最終結果に従う |
| `permission_error` | 自動復旧しない。`.build_state.running` を変更しない。 | `2` |
| `io_error` | 自動復旧しない。`.build_state.running` を変更しない。 | `1` |

backup 名は `{original}.corrupt.{YYYYMMDDHHMMSS}.bak` とする。UTC 秒単位で衝突する場合は `{original}.corrupt.{YYYYMMDDHHMMSS}.{n}.bak` とし、`n` は `2` から始める。

**状態更新禁止事項：**

本チェックだけで `.build_logs/{id}.json` と `.build_history` を作成してはならない。`startup_config_integrity` は `.build_status.json` の `last_trigger` にだけ記録する。ただし、本チェック後に通常 build が発生する場合、通常 build の log / history は実際の build trigger を保存する。

`--dry-run` では、破損検出結果を dry-run JSON の `errors[]` または `warnings[]` に出力するだけとし、backup、初期化、正規化、通知、`.build_status.json` 更新を行わない。

**復旧通知：**

復旧通知は `.notify_config` の検証完了後、復旧対象に `.notify_config` と `.notify_pending` 以外のファイルが 1 件以上ある場合だけ送信する。送信イベントは `config_corrupt` とする。`.notify_pending` が破損復旧された場合、通知失敗時の pending 追記は行わない。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 必須ファイル不在 | 初期値作成、終了コード `0`、追加 backup なし。 |
| `.branch_config` 破損 | backup 後に `.branch_config` 不在、default 採用。 |
| `.build_state` 破損 | backup、初期値作成、`startup_config_integrity` 記録。 |
| unknown key | backup なしで正規化、2 回目起動では追加 WARN なし。 |
| permission error | 自動復旧なし、終了コード `2`、running 未変更。 |
| dry-run | 差分なし、backup なし、dry-run JSON に検出結果。 |
| 通知失敗 | `.notify_pending` が正常な場合だけ pending 追記。 |

### 27.11 ポーリング間隔の動的変更

本機能の目的は、管理 API から systemd timer の実行間隔を変更し、次回以降の runner 起動間隔を固定仕様どおり反映することである。

対象コンポーネントは `components/api.go` とする。`components/runner.go` は本機能で systemd timer を変更しない。

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

### 27.12 GitHub Webhook 受信

本機能の目的は、GitHub push event を HMAC-SHA256 署名検証したうえで受信し、定期 polling を待たずに build queue へ投入することである。

対象コンポーネントは `components/api.go` と `components/runner.go` とする。`components/api.go` は署名検証、イベント記録、queue 投入を担当し、`components/runner.go` は queue entry を処理する。

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

### 27.13 Webhook イベントログ / 一覧取得 API

本機能の目的は、受信した GitHub Webhook の監査情報を `.webhook_events.json` に保存し、管理 API、SDK、UI からページング参照できるようにすることである。

対象コンポーネントは `components/api.go`、`admin/adlaire-ci-sdk.js`、`admin/index.html` とする。

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
| `result` | string | 必須 | `"queued"`, `"ignored_event"`, `"ignored_branch"`, `"queue_full"`, `"error"`。 |

**一覧 API：**

`GET /api/webhook-events` は `limit` と `offset` query を受け付ける。`limit` は 1〜1000、既定値 50。`offset` は 0 以上、既定値 0。新しい順で返す。壊れた行は無視し、server log に `WEBHOOK_EVENT_LOG_SKIP_CORRUPT` を出す。

Response は `{ "events": WebhookEventRecord[], "total": N }` とする。SDK `getWebhookEvents(limit,offset)` は `limit` と `offset` を常に query へ送信する。UI は件数、delivery id、event、branch、sha、result、queued id を表示する。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| event 追記 | JSON Lines へ schema 通り保存される。 |
| 一覧取得 | 新しい順、limit/offset が効く。 |
| 壊れた行 | API は継続し、壊れた行を返さない。 |

### 27.14 ビルド所要時間の記録と統計 API

本機能の目的は、build ごとの開始・終了・所要時間を構造化ログへ保存し、統計 API で直近 N 件の平均、最小、最大を返すことである。

対象コンポーネントは `components/runner.go` と `components/api.go` とする。

**記録仕様：**

`components/runner.go` は `.build_logs/{id}.json` に `started_at`、`finished_at`、`duration_seconds` を必ず保存する。`started_at` は build id 採番直後、`finished_at` は最終 target status 確定直後とする。`duration_seconds` は `finished_at - started_at` を秒単位で切り上げず整数化し、1 秒未満は `0` とする。

`.build_history.duration_seconds` は `.build_logs/{id}.json.duration_seconds` と同じ値にする。失敗、deploy pending、rollback でも記録する。変更なし skip で build log を作らない場合は記録しない。

**統計 API：**

`GET /api/stats/build-duration?n=N` は `.build_logs/` と `.build_logs/archive/` を読み、完了済み build log の `duration_seconds` が `null` でない最新 N 件を集計する。`N` は 1〜1000、既定値 20 とする。

Response は `BuildDurationStats` とし、`count=0` の場合は `avg_seconds`、`min_seconds`、`max_seconds` を `null`、`recent` を `[]` とする。

**異常系：**

| 条件 | 処理 |
|------|------|
| `n` 不正 | `422`。 |
| build log 破損 | 対象 log を除外し、server log に WARN。 |
| archive gzip 展開失敗 | 対象 log を除外し、server log に WARN。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 成功 build | log/history に同一 duration。 |
| 失敗 build | duration を記録する。 |
| 統計対象なし | `count=0`、平均/最小/最大 `null`。 |
| archive 含む | 通常 log と archive log を横断集計する。 |

### 27.15 ビルドアーティファクト管理

本機能の目的は、`.snapshots/` に保存された build artifact を API、SDK、UI から一覧、download、削除、rollback できるようにすることである。

対象コンポーネントは `components/api.go`、`admin/adlaire-ci-sdk.js`、`admin/index.html` とする。snapshot 作成は `components/runner.go` の §14b を正とする。

**API 契約：**

| API | 処理 |
|-----|------|
| `GET /api/snapshots` | `.snapshots/{id}/` を新しい順で一覧する。 |
| `GET /api/snapshots/{id}/download` | 対象 snapshot を tar.gz として streaming download する。 |
| `DELETE /api/snapshots/{id}` | 対象 snapshot だけを削除し、`.config_log` に記録する。 |
| `POST /api/history/{id}/rollback` | 対象 snapshot を deploy target へ再転送し、新規 rollback build log/history を作成する。 |

`id` は build id と一致するものだけ許可する。`/`、`..`、空文字、URL decode 後に path separator を含む値は `422` とする。

**Rollback 仕様：**

rollback は新しい build id を採番し、`.build_history.trigger="rollback"`、`rollback_from=<元id>` を保存する。元 snapshot は変更しない。rollback 中に別 build が running の場合は `409` とする。転送失敗時は rollback build log を `failure` とし、元 snapshot は削除しない。

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

### 27.16 ヘルスチェックエンドポイント

本機能の目的は、認証不要の `GET /api/health` で、外部監視が Adlaire CI の最低限の稼働状態を確認できるようにすることである。

対象コンポーネントは `components/api.go` とする。

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

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 正常 | HTTP 200、`status="ok"`。 |
| pending transfer あり | `pending_transfers` に件数、`status="degraded"`。 |
| build status 破損 | HTTP 200、`status="degraded"`、checks に記録。 |
| read error | HTTP 200 または 500 の条件が仕様通り。 |

### 27.17 ビルドログ重大度フィルター

本機能の目的は、`GET /api/logs/search` と UI ログビューアで重大度別に build log を絞り込めるようにすることである。

対象コンポーネントは `components/api.go`、`admin/adlaire-ci-sdk.js`、`admin/index.html` とする。

**入力：**

`GET /api/logs/search` は既存 query に加えて `level` を受け付ける。`level` の許容値は `"info"`、`"warn"`、`"warning"`、`"error"`、`"debug"` とし、大文字小文字は区別しない。正規化後は `"INFO"`、`"WARNING"`、`"ERROR"`、`"DEBUG"` とする。`warn` は `"WARNING"` と同義とする。不正値は `422`。

**検索対象：**

`.build_logs/{id}.json.stdout`、`stderr`、`warnings`、`error`、archive log を対象とする。行頭が `[WARN]` または `[WARNING]` の行は WARNING、`[ERROR]` または stderr の非空行は ERROR、`[DEBUG]` は DEBUG、それ以外は INFO と分類する。

**SDK / UI：**

SDK `searchLogs(q,from,to,level)` は `level` 指定時だけ query に送信する。UI は INFO / WARNING / ERROR / DEBUG の filter control を提供し、選択なしでは全件を表示する。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| `level=warn` | WARNING 行だけ返る。 |
| `level=error` | ERROR 行だけ返る。 |
| 不正 level | `422`。 |
| archive log | 通常 log と同じ分類で検索される。 |

### 27.18 ブランチ設定の動的変更 API

本機能の目的は、監視対象 branch / target / deploy target を `.branch_config` で管理し、API 経由で変更できるようにすることである。

対象コンポーネントは `components/api.go` と `components/runner.go` とする。

**API：**

`GET /api/branch-config` は `.branch_config` が存在する場合 `{"source":"file","branches":[...]}`、不在の場合 `{"source":"default","branches":[...]}` を返す。`POST /api/branch-config` は `{ "branches": BranchTargetRecord[] }` を受け取る。

**保存仕様：**

永続ファイルの key は必ず `branch_targets` とする。API request / response で `branches` を使う場合も保存前に `branch_targets` へ変換する。空配列を受け取った場合は `.branch_config` を削除し、default 復帰とする。

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

### 27.19 週次ビルドサマリー Webhook

本機能の目的は、過去 7 日間の build 結果を指定曜日・時刻に集計し、Webhook へ定期通知することである。

対象コンポーネントは `components/runner.go`、`components/api.go` とする。runner は自動送信、API は設定表示・手動送信を担当する。

**設定：**

`.notify_config.summary.enabled=true` の場合だけ有効とする。`interval` は `"weekly"`、`hour` は 0〜23、`day_of_week` は 0〜6 とする。タイムゾーンは UTC 固定。

**自動送信条件：**

runner 起動時に、現在 UTC の曜日と時が設定値に一致し、`.build_state.weekly_summary_sent_date` が当日でない場合に送信する。送信成功時だけ `weekly_summary_last_sent_at` と `weekly_summary_sent_date` を更新する。

**集計対象：**

`.build_history` のうち、現在時刻から過去 7 日以内の行を対象とする。`status="success"` を成功、`"failure"`、`"cancelled"`、`"hook_error"` を失敗として数える。所要時間は `duration_seconds != null` の行だけ平均対象にする。

**通知 payload：**

```json
{
  "event": "weekly_summary",
  "period_days": 7,
  "success_count": 10,
  "failure_count": 2,
  "success_rate": 83.33,
  "avg_duration_seconds": 42
}
```

**手動送信 API：**

`POST /api/notify/weekly-summary` は同じ集計を即時送信する。手動送信は `weekly_summary_sent_date` を更新しない。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 条件一致 | Webhook 送信、`.notify_log` 追記、sent date 更新。 |
| 同日二重起動 | 2 回目は送信しない。 |
| 宛先なし | 自動送信は WARN、手動 API は `422`。 |
| 手動送信 | payload を返し、sent date は変更しない。 |

### 27.20 設定変更の詳細 diff 記録

本機能の目的は、設定変更 API が何を変更したかを `.config_log` に機械可読 diff と人間可読 diff の両方で残すことである。

対象コンポーネントは `components/api.go` とする。

**対象 API：**

`.config_log` を write する全 API を対象とする。少なくとも `POST /api/config`、`POST /api/log-level`、`POST /api/notify-config`、`POST /api/repo-config`、`POST /api/branch-config`、`POST /api/webhook-config`、`POST /api/pat-update`、schedule 系 API、maintenance、access-control、hooks、alert-rules、tag-rules、pipeline-config、notes、smtp-config、dashboard-layout、snapshot delete を含む。

**ログ schema：**

各行は §22.0c `.config_log` schema に従う。`diff` は `{key:[before,after]}`、`diff_text` は 1 行以上の文字列とする。差分がない場合、対象 API は状態ファイルを書かず、`.config_log` も追記せず、response は `{ "message": "No changes" }` とする。

**マスク条件：**

キー名に `password`、`token`、`secret`、`pat`、`smtp_password` を含む値は before / after とも `"***"` に置換する。配列や object の内部 key も同じ規則で再帰的にマスクする。

**diff_text 形式：**

`{key}: {before} -> {after}` を key 名昇順で 1 行ずつ連結する。値は JSON 表現とし、secret は `"***"` とする。複数行値は `\n` escape した 1 行 JSON string とする。

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
