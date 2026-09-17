# Adlaire CI — 詳細仕様

本ファイルは `ADLAIRE_CI_SPEC.md` の Part 3 詳細仕様の入口であり、索引、共通固定値、実装前確認項目、検証マトリクス、Phase、詳細節対応表、リポジトリ内ソース配置、横断補足契約を持つ正本である。

各 owner component の具体的な入出力、状態、処理順序、異常系、セキュリティ制約、検証条件の本文は、責務 component 別の `ADLAIRE_CI_DETAIL_*_SPEC.md` を正本とする。

本ファイルに、方針、ポリシー、実装状態、正本関係、ロードマップ状態、実装可否の上位判断を記載してはならない。これらは `ADLAIRE_CI_SPEC.md` を正とする。

---

# Part 3 — 仕様
> Part 3 詳細仕様セットとして、実装の具体的詳細を定める。「どのように動作・実装するか」に答える。

---

## 詳細仕様の読み方

本ファイルは、実装者が実装時に最初に参照する詳細仕様入口だけを扱う。方針、ポリシー、成熟度定義、ロードマップ状態、実装可否、PR 分割判断は `ADLAIRE_CI_SPEC.md` を正とし、本ファイルで再定義しない。

実装者は、対象機能ごとに以下の順で読む。

1. `ADLAIRE_CI_SPEC.md` の実装状態、Part 1 §12、§13 で、対象が実装対象であることを確認する。
2. 本ファイル §0i で、対象機能に対応する詳細仕様節と受け入れ条件を特定する。
3. 本ファイル §0a〜§0h で、詳細仕様の記載基準、共通固定値、実装前確認項目、検証条件、Phase 順序を確認する。
4. owner component の `ADLAIRE_CI_DETAIL_*_SPEC.md` を読み、owner component、collaborator component、入力、出力、状態、正常系、異常系、セキュリティ、検証条件を確認する。
5. `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 のセットアップ・アップデート手順と `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.7 の受け入れ条件に影響がある場合は、実装 PR の検証対象に含める。

詳細仕様節に §0h の必須項目が不足している場合は、実装判断で補完してはならない。先に該当 owner component の詳細仕様ファイルまたは本ファイルの対応表を改訂し、`ADLAIRE_CI_SPEC.md` の対象範囲と整合させる。

| 範囲 | 役割 |
|------|------|
| §0〜§0j | 詳細仕様の記載基準、実装前確認項目、共通固定値、検証、Phase、詳細節対応表、リポジトリ内ソース配置 |
| §1〜§9 | `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md`。`builder` / `adlaire-ci-build` の詳細仕様 |
| §10〜§20 | `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md`。`runner` / `adlaire-ci-runner` の詳細仕様 |
| §21〜§22 | `ADLAIRE_CI_DETAIL_API_SPEC.md`。`api` / `adlaire-ci-api` の詳細仕様。§21a は管理 API サーバー制限を定義する。 |
| §23 | `ADLAIRE_CI_DETAIL_SDK_SPEC.md`。`sdk` の詳細仕様 |
| §24 | `ADLAIRE_CI_DETAIL_UI_SPEC.md`。`ui` の詳細仕様 |
| §25 | `ADLAIRE_CI_DETAIL_API_SPEC.md`。認証の実装仕様 |
| §26 | `ADLAIRE_CI_DETAIL_SETUP_SPEC.md`。バイナリ配布前提のセットアップ、アップデート、受け入れ条件 |
| `admin` | `ADLAIRE_CI_DETAIL_ADMIN_SPEC.md`。管理 UI 静的ファイルの配布物構成、配置、HTTP 静的配信境界 |

責務 component 別詳細仕様ファイルの管理は、§0b.1 に従う。`ADLAIRE_CI_DETAIL_SPEC.md` は入口、索引、共通固定値、責務 component 対応表、横断補足契約だけを持つ。責務 component 別詳細仕様ファイルは、それぞれの owner component と collaborator component の詳細仕様だけを持つ。

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

未確定の内容は、実装可能な詳細仕様として記載してはならない。未確定の場合は、本ファイルまたは責務 component 別詳細仕様ファイルへ推測で具体値を記載せず、`ADLAIRE_CI_SPEC.md` で状態を確認する。

対象範囲の内容は、実装ファイルが存在しなくても、本節の基準に従って責務 component 別詳細仕様ファイルへ実装可能な粒度まで具体化する。

---

## 0b. 詳細仕様参照表

本節は、責務 component ごとに参照する詳細仕様節を示す。実装状態、実装可否、ロードマップ状態は `ADLAIRE_CI_SPEC.md` を確認する。

| 責務 component | 詳細仕様節 | 主な確認対象 |
|--------------------|------------|--------------|
| `builder` | `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §1〜§9、§8a、§27.4、§27.25、§27.28 | CLI、入力 Markdown、出力サイト、HTML / CSS / JavaScript、変換 report、fixture、builder owner 追加機能。 |
| `runner` | `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §10〜§20、§15a、§27.2〜§27.3、§27.8〜§27.10、§27.14、§27.19、§27.21〜§27.24、§27.26〜§27.27、§27.29、§27.31〜§27.38、`ADLAIRE_CI_DETAIL_COMMITSTATUS_SPEC.md` §27.1、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 | 設定、状態ファイル、GitHub API 読取、pipeline、転送、snapshot、通知、fixture、runner owner 追加機能、Commit Status 呼び出し境界、systemd / setup 参照境界。 |
| `api` | `ADLAIRE_CI_DETAIL_API_SPEC.md` §21〜§22、§21a、§25、§27.5〜§27.6、§27.11〜§27.13、§27.16〜§27.18、§27.20、§27.30、§27.42〜§27.47、`ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` §27.42〜§27.47、`ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a、§22.0c | API 共通処理、API server 制限、endpoint、状態ファイル read/write 呼び出し境界、認証連携、security 呼び出し境界、API owner 追加機能。 |
| `admin` | `ADLAIRE_CI_DETAIL_ADMIN_SPEC.md` §0、A1〜A5 | 管理 UI 静的ファイルの配布物構成、配置、検証、HTTP 静的配信境界。 |
| `sdk` | `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 | SDK class、method、HTTP 対応、error、stream、token 破棄。 |
| `ui` | `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 | 画面構成、DOM id、panel、SDK 呼び出し、表示状態、秘密情報消去。 |
| `setup` | `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 | バイナリ配布、配置、systemd、セットアップ、アップデート、リリース成果物検証。 |
| `statefile` | `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a、§22.0c | 状態ファイル共通仕様、lock、atomic write、JSON Lines、破損時処理、状態読取 adapter、主要 schema。 |
| `archive` | `ADLAIRE_CI_DETAIL_ARCHIVE_SPEC.md` §27.7、§27.15 | build log archive、snapshot、download、delete、rollback、cleanup。 |
| `commitstatus` | `ADLAIRE_CI_DETAIL_COMMITSTATUS_SPEC.md` §27.1 | GitHub Commit Status API payload、送信順、失敗時非反転、保存値、secret mask。 |
| `security` | `ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` §27.42〜§27.47 | API token scope、API key、audit、session timeout、TOTP、rate limit、漏えい禁止、security 横断順序。 |
| `mcp` | 詳細仕様なし | 将来計画。現時点では実装可能な入出力、状態、起動手順、ツール定義、検証条件を定義しない。 |

上表の `詳細仕様節` は参照入口であり、主本文の owner component を変更しない。複数ファイルを参照する行では、対象機能の owner component のファイルを主本文とし、他ファイルは collaborator の境界、schema、fixture、security、setup、受け入れ条件を確認するために読む。参照先に同じ HTTP body、状態 schema、DOM id、SDK method、fixture assertion を重複定義してはならない。

---

## 0b.1 責務 component 別 詳細仕様ファイル管理仕様

本節は、責務 component 別に分割済みの詳細仕様ファイルを維持するための固定仕様である。責務境界の変更、仕様内容の移動、参照先更新を行う場合も、機能追加、実装状態変更、実装可否変更、ロードマップ変更、方針・ポリシー追加を含めてはならない。

詳細仕様ファイルは以下に固定する。`COMMON`、`CORE`、`BASE`、`SHARED`、`FOUNDATION`、その他の横断共通基盤ファイルは作成しない。

| ファイル | 持つ内容 | 持たない内容 |
|----------|----------|--------------|
| `ADLAIRE_CI_DETAIL_SPEC.md` | 詳細仕様の入口、読み方、共通固定値、実装前確認項目、検証マトリクス、Phase、詳細節対応表、リポジトリ内ソース配置、責務 component 別詳細仕様ファイル管理仕様。 | 各 component の詳細な処理本文、fixture 詳細、状態ファイル schema 詳細、個別 endpoint 詳細、個別 UI 操作詳細。 |
| `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` | `builder` owner の Markdown 変換、静的 Web サイト出力、HTML / CSS / JavaScript、theme component、builder fixture。 | runner / api / sdk / ui の実行責務。 |
| `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` | `runner` owner の GitHub 監視、状態ファイル更新、pipeline、deploy、snapshot、通知、runner fixture。 | API endpoint の認証・応答本文、SDK method、UI DOM 詳細。 |
| `ADLAIRE_CI_DETAIL_API_SPEC.md` | `api` owner の HTTP 共通契約、endpoint、状態ファイル read/write、認証連携。 | SDK 内部実装、UI DOM 詳細、runner の build 実行責務、fixture 詳細。 |
| `ADLAIRE_CI_DETAIL_ADMIN_SPEC.md` | `admin` owner の管理 UI 静的ファイル配布物構成、配置、検証、HTTP 静的配信境界。 | UI DOM 詳細、SDK method 実装、API endpoint 実装、systemd 導入手順。 |
| `ADLAIRE_CI_DETAIL_SDK_SPEC.md` | `sdk` owner の SDK class、method、HTTP 対応、error、stream、token 破棄。 | API endpoint の状態ファイル更新責務、UI DOM 詳細。 |
| `ADLAIRE_CI_DETAIL_UI_SPEC.md` | `ui` owner の DOM id、panel、操作、表示状態、SDK 呼び出し、秘密情報消去。 | SDK method 実装、API endpoint 実装、状態ファイル直接操作。 |
| `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` | `setup` owner のバイナリ配布、配置、systemd、セットアップ、アップデート、リリース成果物検証。 | runner / api / sdk / ui の個別機能本文。 |
| `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` | `statefile` owner の状態ファイル共通仕様、lock、atomic write、JSON Lines、破損時処理、状態読取 adapter、主要 schema。 | API endpoint の request / response、runner の業務処理、UI 表示判断。 |
| `ADLAIRE_CI_DETAIL_ARCHIVE_SPEC.md` | `archive` owner の build log archive、snapshot、download、delete、rollback、cleanup。 | runner の build 実行、API 共通 request / response、SDK method 実装、UI DOM 詳細。 |
| `ADLAIRE_CI_DETAIL_COMMITSTATUS_SPEC.md` | `commitstatus` owner の GitHub Commit Status API payload、送信順、失敗時非反転、保存値、secret mask。 | runner の build 実行判断、GitHub read、API endpoint、SDK method、UI DOM 詳細。 |
| `ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` | `security` owner の API token scope、API key、audit、session timeout、TOTP、rate limit、漏えい禁止、security 横断順序。 | API endpoint 共通処理、SDK method 実装、UI DOM 詳細、runner / builder の業務処理。 |
| `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` | fixture manifest、assertion、fake、testdata、API P0〜P5 fixture、受け入れ fixture 共通契約、PR 証跡テンプレート。 | 個別 component の通常処理本文。 |

`ADLAIRE_CI_DETAIL_SPEC.md` §27.38a は、runner、builder、api、sdk、ui、statefile、archive にまたがる横断補足契約であり、責務 component 別の分割先へ移動しない。§27.21〜§27.38 または api / sdk / ui / statefile の横断連動を実装する場合は、owner component の分割先詳細仕様ファイルを正本とし、§27.38a は横断処理順、同期禁止、成功後再取得、失敗時固定、横断受け入れ観点の確認として読む。

詳細仕様の配置単位は、owner component を第一基準とする。複数 component が関わる機能は、owner component のファイルに主本文を置き、collaborator component のファイルには参照リンク、禁止事項、受け入れ観点だけを置く。主本文を複数ファイルへ重複定義してはならない。

すべての責務 component 別詳細仕様ファイルは、冒頭に `## 0. 責務境界` を置き、以下の 4 項目を同じ意味で持つ。

| 項目 | 必須内容 |
|------|----------|
| owner component | そのファイルが主本文として扱う component を 1 件だけ書く。 |
| collaborator component | 呼び出し元、呼び出し先、schema 参照先、表示参照先、検証参照先を 0 件以上書く。owner component を含めてはならない。 |
| 持つ内容 | そのファイルだけが主本文として定義する入出力、状態、処理、異常系、検証条件を書く。 |
| 持たない内容 | 他 owner component へ委ねる処理、状態、API、SDK、UI、fixture、setup、security を書く。 |

責務境界表の `持つ内容` と `持たない内容` が本文と矛盾する場合は、本文を実装判断に使ってはならない。先に責務境界表、本文、§0b の詳細仕様参照表、`DOCUMENT_INDEX.md` を同時に整合させる。

責務 component 別詳細仕様ファイルの各節は、以下を満たす。

| 項目 | 必須条件 |
|------|----------|
| 節番号 | 既存の節番号を維持する。番号の再採番は行わない。 |
| 参照 | 入口ファイルと owner component ファイルの参照先が一意に追跡できるよう、`ADLAIRE_CI_DETAIL_SPEC.md` の対応表を更新する。 |
| owner | 各機能節に owner component を 1 件だけ明記する。 |
| collaborator | collaborator component は 0 件以上を明記し、owner component を含めない。 |
| 重複禁止 | 同じ入力、出力、状態 schema、HTTP body、DOM id、fixture assertion を複数ファイルで重複定義しない。 |
| 横断事項 | 横断する固定値は `ADLAIRE_CI_DETAIL_SPEC.md` に置く。横断共通基盤を component として扱わない。 |
| 索引 | `DOCUMENT_INDEX.md` に、分割後ファイルの役割と正本範囲を反映する。 |

実装者が詳細仕様を読む順序は以下に固定する。

1. `ADLAIRE_CI_SPEC.md` で実装対象、実装状態、実装可否を確認する。
2. `ADLAIRE_CI_DETAIL_SPEC.md` で共通固定値、責務 component、詳細節対応表を確認する。
3. owner component の分割先詳細仕様ファイルを読む。
4. collaborator component がある場合は、該当する分割先詳細仕様ファイルの参照節を読む。
5. 状態ファイルの読み書き、lock、atomic write、schema を扱う場合は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` を読む。
6. 認証、scope、token、audit、session、TOTP、rate limit、漏えい禁止を扱う場合は `ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` を読む。
7. fixture、fake、PR 証跡が必要な場合は `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` を読む。

詳細仕様ファイルの責務整理、節移動、参照先更新は、以下の完了条件をすべて満たすまで完了扱いにしてはならない。

| 完了条件 | 判定 |
|----------|------|
| 旧ファイル内の移動対象本文が対応する分割先に移動している。 | 必須 |
| `ADLAIRE_CI_DETAIL_SPEC.md` には入口、索引、共通固定値、対応表、管理仕様、横断補足契約だけが残っている。 | 必須 |
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

上記ゲートのいずれかが未充足の場合、実装判断で補完してはならない。先に該当 owner component の詳細仕様ファイル、本ファイルの対応表、または `ADLAIRE_CI_SPEC.md` を改訂し、未充足項目を仕様として確定する。

実装後の完了条件は以下とする。

1. 実装した機能が、該当 owner component の詳細仕様ファイルと本ファイルの対応表に記載された入力、出力、状態、異常系、検証条件と一致する。
2. 対象機能が owner component の詳細仕様本文で §0h の機能仕様テンプレートを満たし、§0i の詳細節対応表の受け入れ条件を満たしている。
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
| 禁止事項 | 仕様にない環境変数、状態ファイル、HTTP endpoint、CLI option、外部依存を実装者判断で追加してはならない。必要な場合は先に該当 owner component の詳細仕様ファイルまたは本ファイルの対応表を改訂する。 |

---

## 0e. 完全実装検証マトリクス

対象項目を完了扱いにする場合は、責務 component ごとに下表の検証を満たす。実装ファイルが存在しても、本表の必須検証が未完了の場合は完了扱いにしない。

本表の `builder`、`runner`、`api`、`sdk`、`ui`、`setup` は Phase の主対象 component である。`statefile`、`security`、`archive`、`commitstatus`、`admin`、`fixture` は、主対象 component の collaborator component として完了判定に参加する。collaborator component の検証が失敗する場合、主対象 component の実装も完了扱いにしてはならない。

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
| `sdk` | SDK 契約 | 全 method が §22.0e の endpoint のみを呼び、body なし endpoint に body を送らず、HTTP error を `AdlaireCIError` として返す。 |
| `ui` | UI 契約 | 全操作が §24 の SDK method 経由で動作し、成功表示、失敗表示、disabled、再取得、秘密情報消去が一致する。 |
| `setup` | systemd | `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 の unit 名、`ExecStart`、配置パス、権限、起動確認コマンドが実際の導入手順と一致する。 |
| `statefile` | 状態ファイル契約 | 状態ファイルの schema、lock、atomic write、JSON Lines、破損時処理、権限、秘密情報マスクが `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a、§22.0c と一致する。 |
| `security` | 認証・認可・漏えい禁止 | scope、API key、audit、session、TOTP、rate limit、秘密情報非表示、失敗時副作用が `ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` §27.42〜§27.47 と一致する。 |
| `archive` | artifact / log archive | gzip archive、snapshot、download、delete、rollback、cleanup の実体処理が `ADLAIRE_CI_DETAIL_ARCHIVE_SPEC.md` §27.7、§27.15 と一致し、API / SDK / UI の応答契約を上書きしない。 |
| `commitstatus` | GitHub Commit Status | payload、送信順、失敗時非反転、保存値、secret mask が `ADLAIRE_CI_DETAIL_COMMITSTATUS_SPEC.md` §27.1 と一致し、runner の build 実行判断を上書きしない。 |
| `admin` | 静的配布境界 | admin 配布物、archive validation、HTTP 静的配信、setup 連携が `ADLAIRE_CI_DETAIL_ADMIN_SPEC.md` §0、A1〜A5 と一致し、UI / SDK の本文を重複定義しない。 |
| `fixture` | fixture / fake / 証跡 | Phase 別 fixture、fake、assertion、expected / effects、PR 証跡が `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §0g.8-F、§22-F、§27-F と一致する。 |

検証結果は、実装 PR の本文または実装完了報告に、対象、実行コマンド、期待結果、実結果を対応付けて記録する。検証不能な項目がある場合は、その項目を完了扱いにしてはならない。

---

## 0f. 仕様策定完了チェック

本節は、Go 版初期実装へ進む前に仕様策定が完了しているかを判定するチェックである。実装者は、責務 component ごとに下表の必須条件を満たすまで実装を開始してはならない。

| 対象 | 実装着手条件 | 実装禁止条件 | 完了判定 |
|------|--------------|--------------|----------|
| `builder` | §2〜§8 に CLI option、入力 Markdown、出力サイトディレクトリ、終了コード、stderr、HTML 構造、テーマコンポーネント、JS/CSS、生成物確認が定義されている。 | §4〜§7 にない Markdown 記法、CSS class、JavaScript 機能、外部 asset、theme を追加すること。 | §0e の `builder` 必須検証をすべて満たし、生成サイトが §5〜§7 と一致する。 |
| `runner` | §10〜§20 に設定値、状態ファイル、GitHub API、SHA 比較、pipeline 起動、SSH 転送、snapshot、通知、ログ、systemd / setup 参照境界が定義されている。systemd unit 本文とセットアップ手順は `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 を正とする。 | 未定義の環境変数、状態ファイル、queue 挙動、通知チャンネル、pipeline 形式を追加すること。systemd unit file の生成、配置、更新、enable、restart を runner に追加すること。 | §0e の `runner` 必須検証をすべて満たし、状態ファイル更新順序が §13、§22.0a、§22.0d と一致する。 |
| `api` | `ADLAIRE_CI_DETAIL_API_SPEC.md` §21〜§22、§21a、§25 と `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 に API 共通契約、API server 制限、endpoint、状態ファイル schema、認証、認可、systemd、セットアップが定義されている。 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e にない endpoint、method、status code、response body、状態ファイル write を追加すること。 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e の全 endpoint が Request、Response、Errors、Read、Write、SDK、UI の対応表と一致する。 |
| `sdk` | `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 に SDK class、method、引数、戻り値、HTTP endpoint 対応、error object、token 破棄条件が定義されている。 | SDK が `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e にない endpoint を呼ぶこと、body 禁止 endpoint に body を送ること、独自 error 形式を返すこと。 | 全 method が `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e と `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 の対応どおりに動作し、HTTP error を `AdlaireCIError` として扱う。 |
| `ui` | `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 に画面構成、panel、操作、成功表示、失敗表示、disabled、再取得、秘密情報消去が定義されている。 | SDK を介さず API を直接呼ぶこと、未定義の画面・操作・保存先を追加すること、秘密情報を DOM に残すこと。 | 全 UI 操作が `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 の表示条件と `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 の SDK method を満たし、秘密情報 field が指定条件で消去される。 |

上表の対象外である `mcp`、MCP tools、MCP resources、MCP prompts、HTTP SSE transport、MCP audit / stats / config CRUD は、初期実装では実装しない。これらは、Part 3 詳細仕様セット内に入出力、状態、起動手順、検証条件を定義しない。

仕様策定完了チェックで未充足が見つかった場合は、実装を開始せず、以下の順で仕様を補完する。

1. 未充足項目が本ファイルの記載対象外である場合は、先に `ADLAIRE_CI_SPEC.md` を確認する。
2. 未充足項目が入出力、状態ファイル、api、sdk、ui、処理順序、異常系、検証条件に関わる場合は、該当 owner component の詳細仕様ファイルまたは collaborator の詳細仕様ファイルを改訂する。
3. ファイル名、正本関係、対象範囲が変わる場合は、`DOCUMENT_INDEX.md` の更新要否を確認する。
4. 対象項目の詳細節、参照先、受け入れ条件が変わる場合は、§0i の詳細節対応表を更新する。
5. 補完後、§0b、§0c、§0e、本節、§0g、§0h、§0i、§0j の条件を再確認する。

---

## 0g. 初期実装 Phase 分割

Go 版初期実装は、`ADLAIRE_CI_SPEC.md` §0e の対象範囲を一括実装せず、下表の Phase 順に進める。上位 Phase の完了判定を満たす前に、下位 Phase の実装 PR を開始してはならない。

各 Phase の `対象` は、その Phase の owner component を示す。状態ファイル、security、archive、commitstatus、admin、fixture、setup が関わる場合も、それらは collaborator component として該当 Phase の完了条件に含める。collaborator component の詳細仕様に未充足がある場合は、owner component の実装で補完せず、先に該当する責務 component 別詳細仕様ファイルを改訂する。

| Phase | 対象 | 実装範囲 | 依存条件 | 完了条件 |
|-------|------|----------|----------|----------|
| Phase 1 | `builder` | §2〜§9 の CLI、Markdown 変換、静的 Web サイト出力、テーマコンポーネント、生成物確認。 | なし。 | §0e の `builder` 必須検証と §0f の `builder` 完了判定を満たす。 |
| Phase 2 | `runner` | §10〜§20 の CI ランナー、GitHub API 連携、SHA キャッシュ、pipeline 起動、SSH 転送、snapshot、通知、ログ、systemd / setup 参照境界。systemd unit 本文と配置手順は `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 を正とする。 | Phase 1 が完了し、`adlaire-ci-build` の CLI 契約が固定されている。 | §0e の `runner` 必須検証と §0f の `runner` 完了判定を満たす。 |
| Phase 3 | `api` P0 / P1 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §21〜§22、§21a、§25 と `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 のうち、認証、セッション、共通エラー、状態ファイル読み書き、ビルド操作、status、logs、history、queue、circuit breaker。 | Phase 2 が完了し、runner が書き込む状態ファイル schema が固定されている。 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0f P0 / P1 の必須検証、§0e の `api` API 共通・状態ファイル検証、§0f の `api` 完了判定の該当範囲を満たす。 |
| Phase 4 | `api` P2〜P5 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0f P2〜P5 の config、repo、branch、schedule、notify、snapshot、rollback、maintenance、access control、hooks、tokens 等。 | Phase 3 が完了し、API 共通処理と認証が固定されている。 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0f P2〜P5 の必須検証と §0e の `api` endpoint 契約を満たす。 |
| Phase 5 | `sdk` | `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 の SDK class、method、戻り値、HTTP error、token 破棄、query 生成。 | Phase 3 と Phase 4 が完了し、`ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e の endpoint 契約が固定されている。 | §0e の `sdk` 契約と §0f の `sdk` 完了判定を満たす。 |
| Phase 6 | `ui` | `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 の標準管理ツール UI、panel、操作、成功表示、失敗表示、disabled、再取得、秘密情報消去。 | Phase 5 が完了し、SDK method 契約が固定されている。 | §0e の `ui` 契約と §0f の `ui` 完了判定を満たす。 |

### 0g.1 Phase 1 完全仕様ゲート（`builder`）

Phase 1 の実装詳細本文は `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` を正とする。親ファイルでは、Phase 1 の対象、依存条件、完了条件、後続 Phase への引き継ぎ確認だけを扱う。

| 確認 | 参照先 |
|------|--------|
| CLI、Markdown 変換、静的 Web サイト出力、theme component、生成物確認 | `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` §1〜§9、§8a |
| Phase 1 fixture、testdata、PR 証跡 | `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §0g.8-F |
| release / setup 受け入れ条件 | `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.7 |

### 0g.2 Phase 2 完全仕様ゲート（`runner`）

Phase 2 の実装詳細本文は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` を正とする。親ファイルでは、Phase 2 が Phase 1 の `adlaire-ci-build` 契約に依存し、後続 API が読む runner 状態契約を固定することだけを扱う。

| 確認 | 参照先 |
|------|--------|
| CI runner、GitHub API 連携、SHA cache、pipeline、deploy、snapshot、通知、systemd | `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §10〜§20、§15a |
| runner fixture、fake GitHub、fake ssh / notifier、PR 証跡 | `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §0g.8-F |
| release / setup 受け入れ条件 | `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.7 |

### 0g.3 Phase 3 完全仕様ゲート（`api` P0 / P1）

Phase 3 の実装詳細本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` と `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` を正とする。親ファイルでは、API 共通契約、認証、状態 read/write、P0 / P1 endpoint が Phase 4〜6 の前提になることだけを扱う。

| 確認 | 参照先 |
|------|--------|
| API 共通処理、API server 制限、認証、P0 / P1 endpoint、状態 read/write | `ADLAIRE_CI_DETAIL_API_SPEC.md` §21〜§22、§21a、§25、§22.0f |
| 状態ファイル schema、lock、atomic write | `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a、§22.0c |
| 管理 API 導入、systemd、release 受け入れ条件 | `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 |

### 0g.4 Phase 4 完全仕様ゲート（`api` P2〜P5）

Phase 4 の実装詳細本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` を正とする。親ファイルでは、P2〜P5 endpoint が SDK / UI の最終入力契約になることだけを扱う。

| 確認 | 参照先 |
|------|--------|
| P2〜P5 endpoint、request / response、error、auth、secret mask | `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e、§22.0f |
| security 連携 | `ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` §27.42〜§27.47 |
| API fixture、endpoint 証跡 | `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §0g.8-F、§22-F、§27-F |

### 0g.5 Phase 5 完全仕様ゲート（`sdk`）

Phase 5 の実装詳細本文は `ADLAIRE_CI_DETAIL_SDK_SPEC.md` を正とする。親ファイルでは、SDK が固定済み API endpoint だけを呼び、UI 表示判断を持たないことだけを扱う。

| 確認 | 参照先 |
|------|--------|
| SDK class、method、HTTP 対応、error、stream、token 破棄 | `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 |
| API endpoint 対応 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e |
| SDK fixture、fake fetch / stream | `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §0g.8-F |

### 0g.6 Phase 6 完全仕様ゲート（`ui`）

Phase 6 の実装詳細本文は `ADLAIRE_CI_DETAIL_UI_SPEC.md` を正とする。親ファイルでは、UI が SDK 経由だけで API と通信し、秘密情報を DOM に残さないことだけを扱う。

| 確認 | 参照先 |
|------|--------|
| DOM id、panel、操作、表示状態、SDK 呼び出し、秘密情報消去 | `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 |
| SDK method 契約 | `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 |
| UI fixture、fake SDK、直接 API 呼び出し禁止確認 | `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §0g.8-F |

### 0g.7 Phase 間引き継ぎ契約

各 Phase の完了時は、次 Phase が依存する契約を変更不可として扱う。後続 Phase で変更が必要になった場合は、後続 Phase の実装で吸収せず、契約を定義した owner 詳細仕様へ戻す。

| 引き継ぎ元 | 引き継ぎ先 | 固定する契約 | 正本 |
|------------|------------|--------------|------|
| Phase 1 | Phase 2 | `adlaire-ci-build` CLI、終了コード、stdout / stderr、`[REPORT]`、出力サイト構造。 | `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` |
| Phase 2 | Phase 3 | runner 状態ファイル schema、lock、history/log、pending queue、circuit breaker、snapshot、通知ログ。 | `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` / `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` |
| Phase 3 | Phase 4 | API 共通契約、認証、session、error body、validation、lock error、SSE 基本形式。 | `ADLAIRE_CI_DETAIL_API_SPEC.md` |
| Phase 4 | Phase 5 | 全 endpoint の method、path、query、request、response、error、認証要否。 | `ADLAIRE_CI_DETAIL_API_SPEC.md` |
| Phase 5 | Phase 6 | SDK method 名、引数、戻り値、error object、stream handle、token 破棄条件。 | `ADLAIRE_CI_DETAIL_SDK_SPEC.md` |
| Phase 6 | 初期実装完了 | UI 操作、表示状態、secret 消去、SDK 経由通信、実装完了検証結果。 | `ADLAIRE_CI_DETAIL_UI_SPEC.md` |

### 0g.8 Phase 別 実装 PR 成果物チェックリスト

Phase fixture / testdata 配置、fake 実装、実装 PR 証跡の詳細は `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §0g.8-F を正とする。親ファイルでは、Phase ごとの成果物参照先だけを保持する。

| Phase | 実装対象 | 成果物・fixture 正本 | 受け入れ条件 |
|-------|----------|----------------------|--------------|
| Phase 1 | `builder` | `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §0g.8-F | `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` と `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.7 を満たす。 |
| Phase 2 | `runner` | `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §0g.8-F | `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` と `ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26.7 を満たす。 |
| Phase 3 | `api` P0 / P1 | `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §0g.8-F、§27-F | `ADLAIRE_CI_DETAIL_API_SPEC.md` P0 / P1 と setup API 導入条件を満たす。 |
| Phase 4 | `api` P2〜P5 | `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §0g.8-F、§27-F | `ADLAIRE_CI_DETAIL_API_SPEC.md` P2〜P5 を満たす。 |
| Phase 5 | `sdk` | `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §0g.8-F | `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 を満たす。 |
| Phase 6 | `ui` | `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §0g.8-F | `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 を満たす。 |

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

上表のいずれかが不足する対象項目は、実装者判断で補完してはならない。不足を見つけた場合は、実装 PR ではなく仕様改訂 PR として該当 owner component の詳細仕様ファイルまたは本ファイルの対応表を先に更新する。

---

## 0i. 詳細節対応表

本節は、対象機能から該当する詳細仕様へ移動するための対応表である。実装者は対象機能を実装する前に、下表の「詳細仕様節」と「受け入れ条件」を確認する。

表の「責務 component」は参照先を探すための component 一覧である。owner component と collaborator component は、対象機能の詳細仕様節に記載された値を正とする。

表の「詳細仕様節」が複数ある場合は、owner component の詳細仕様ファイルを主本文として読み、collaborator component の詳細仕様ファイルは schema、呼び出し境界、表示、security、setup、fixture、検証観点の確認として読む。ファイル名を伴わない裸の節番号は、同じ行の「責務 component」から該当する owner component または collaborator component の詳細仕様ファイルへ解決する。`builder` が責務 component に含まれる機能では、`ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` の同番号節を合わせて確認する。`runner` が責務 component に含まれる機能では、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` の同番号節を合わせて確認する。`api` が責務 component に含まれる機能では、`ADLAIRE_CI_DETAIL_API_SPEC.md` の同番号節を合わせて確認する。`sdk` が責務 component に含まれる機能では、`ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23 を合わせて確認する。`ui` が責務 component に含まれる機能では、`ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 を合わせて確認する。`statefile`、`security`、`archive`、`commitstatus`、`admin`、`fixture`、`setup` が責務 component に含まれる場合は、それぞれの責務 component 別詳細仕様ファイルを合わせて確認する。該当節に §0h の必須項目が不足している場合は、その項目を実装せず、先に詳細仕様を改訂する。

詳細節対応表は owner component を置き換える表ではない。受け入れ条件が複数 component にまたがる場合でも、主本文は owner component の詳細仕様ファイルを正とし、collaborator component の詳細仕様は schema、呼び出し境界、表示、security、setup、fixture、検証観点だけを補完する。

| 機能 | 責務 component | 詳細仕様節 | 受け入れ条件 |
|------|-------------------|------------|--------------|
| ビルドタイムアウト | `runner` / `api` | §12、§13、§22.0e | `build_timeout_seconds` の既定値、設定 API、`context.WithTimeout` の中断処理、終了コード、ログが一致する。 |
| ポーリング間隔の動的変更 | `api` | §22.0e、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26、§27.11 | `POST /api/schedule/interval` が systemd timer 設定を更新し、検証コマンドで反映を確認できる。 |
| ビルドログのファイル保存 | `runner` | §11、§13、§15 | `.build_logs/{id}.json` の schema、stdout/stderr、変換レポート、duration、権限が一致する。 |
| GitHub Webhook 受信 | `api` / `runner` | §22.0e、§22-W、§13、§27.12 | HMAC 検証、イベント記録、キュー投入またはビルドトリガー、エラー応答が一致する。 |
| ネットワーク断時の再試行 | `runner` | §12、§13 | `API_RETRY_MAX`、`API_RETRY_BASE_SECONDS`、指数バックオフ、失敗時ログが一致する。 |
| GitHub API レート制限自動待機 | `runner` / `api` | §13、§22.0e | `X-RateLimit-Remaining` と `X-RateLimit-Reset` の扱い、待機、API 表示が一致する。 |
| 転送後リモート整合性検証 | `runner` | §14a、§13 | SSH 転送後の SHA256 照合、不一致時の `.pending_transfers` 再投入、ログが一致する。 |
| マルチブランチビルド | `runner` / `api` | §12、§13、§22.0e | `BRANCH_TARGETS` と `.branch_config` の優先順位、順次処理、API 更新が一致する。 |
| ビルドログ世代管理 | `runner` | §12、§13、§15 | `LOG_KEEP_N` 超過時の削除順序、0 の扱い、削除ログが一致する。 |
| ビルド出力の外部転送 | `runner` | §14a、§13 | SSH 差分転送、複数ファイル処理、失敗時 pending、通知が一致する。 |
| ビルドクールダウン | `runner` | §12、§13 | `BUILD_COOLDOWN_SECONDS` 内の起動スキップ、Webhook 二重トリガー抑止、ログが一致する。 |
| ビルド前の事前チェック | `runner` | §13、`ADLAIRE_CI_DETAIL_SETUP_SPEC.md` §26 | ディスク、`adlaire-ci-build`、pipeline 前提の確認、不足時の ERROR と通知が一致する。 |
| 定期強制ビルド | `runner` / `api` | §12、§13、§22.0e | `FORCE_BUILD_INTERVAL`、変更なし時の強制ビルド、設定 API が一致する。 |
| ビルド中重複スキップ | `runner` | §11、§13 | `.build_lock` の PID 判定、stale lock、競合時終了コードとログが一致する。 |
| GitHub PAT 有効期限の事前警告 | `runner` / `api` | §13、§22.0e | `GitHub-Authentication-Token-Expiration` の解析、7 日以内 WARN、API 表示が一致する。 |
| コミット情報のビルドログ記録 | `runner` | §13、§15 | SHA、message、author、date を build id と同じログへ記録する。 |
| GitHub API 連続失敗によるサーキットブレーカー | `runner` / `api` | §11、§12、§13、§22.0e | 閾値、open/close 状態、API reset、通知、状態ファイルが一致する。 |
| 出力サイトサイズ警告閾値 | `builder` / `runner` / `api` | §8、§12、§13、§22.0e | `OUTPUT_SIZE_WARN_MB`、`size_warn`、WARN ログ、API 表示が一致する。 |
| 設定ファイル起動時整合性チェック | `runner` | §11、§12、§13、§22.0a、§22.0c、§27.10 | 対象 JSON ファイル、検証順序、破損退避、初期化値、ログ、通知、終了コード、fixture が一致する。 |
| ビルドステータスファイル出力 | `runner` / `api` | §11、§13、§15、§22.0a、§22.0c、§22.0e、§27.8 | `.build_status.json` の schema、更新タイミング、status/target_status、pending 件数、circuit 状態、API 参照元が一致する。 |
| ビルドトリガー種別の記録 | `runner` / `api` / `sdk` / `ui` | §13、§15、§22.0c、§22.0e、§23、§24、§27.9 | `trigger` の有効値、判定条件、`.build_logs`、`.build_history`、`.build_status.json`、履歴 filter、UI 表示が一致する。 |
| GitHub Commit Status API | `runner` | §12、§13、§15、§22.0c、§27.1 | `commit_status_enabled`、context、target_url、pending/success/failure の送信条件、失敗時の扱い、build log 記録が一致する。 |
| ドライラン実行モード | `runner` | §11、§12、§13、§15、§27.2 | `--dry-run` が状態ファイル、log、history、deploy、通知を変更せず、設定・GitHub・SHA 判定結果を固定 JSON で返す。 |
| ビルド失敗時の自動リトライ | `runner` | §12、§13、§15、§22.0c、§27.3 | retry 対象エラー、最大回数、backoff、attempt log、最終 status、SHA 更新禁止条件が一致する。 |
| 出力サイトへのビルドメタ埋め込み | `builder` / `runner` / `api` | §2、§5、§8、§13、§22.0e、§27.4 | CLI/env 入力、HTML meta、REPORT、build log、`GET /api/output-meta` の値が一致する。 |
| 設定バリデーション API | `api` / `sdk` / `ui` | §22.0c、§22.0e、§23、§24、§27.5 | `POST /api/config/validate` が状態を変更せず、正規化後設定、warnings、errors を返す。 |
| API アクセスログ | `api` / `sdk` / `ui` | §22.0a、§22.0c、§22.0e、§23、§24、§27.6 | `.api_access_log` の schema、追記対象、マスク条件、一覧 API、UI 表示が一致する。 |
| ビルドログのアーカイブ圧縮 | `runner` / `api` / `sdk` / `ui` | §12、§13、§15、§22.0c、§22.0e、§23、§24、§27.7 | gzip 形式、archive 先、参照順、cleanup/archive API、disk usage 集計、UI 表示が一致する。 |
| Webhook イベントログ | `api` | §11、§22.0e、§22-W、§27.13 | `.webhook_events.json` の JSON Lines schema と一覧 API が一致する。 |
| ビルド所要時間の記録と統計 API | `runner` / `api` | §15、§22.0e、§27.14 | `started_at`、`finished_at`、`duration_seconds` と統計 API が一致する。 |
| ビルドアーティファクト世代管理 | `runner` / `api` | §14b、§22.0e | `.snapshots/` の保持世代、削除、rollback API が一致する。 |
| ビルドアーティファクト管理 | `api` / `ui` / `sdk` | §14b、§22.0e、§23、§24、§27.15 | 一覧、ダウンロード、削除、ロールバックの API、SDK、UI が一致する。 |
| ヘルスチェックエンドポイント | `api` | §22.0e、§27.16 | `GET /api/health` の稼働秒数、最終ビルド、最終転送、エラー応答が一致する。 |
| Webhook イベント一覧取得 API | `api` / `sdk` / `ui` | §22.0e、§23、§24、§27.13 | `GET /api/webhook-events` の query、response、SDK method、UI 表示が一致する。 |
| ビルドログ重大度フィルター | `api` / `sdk` / `ui` | §22.0e、§23、§24、§27.17 | `level=warn\|error` の query、検索結果、UI filter が一致する。 |
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
| Webhook 通知失敗リトライキュー | `runner` | §11、§13、§16 | `.notify_pending` の schema、再送順序、失敗時保持が一致する。 |
| ブランチ設定の動的変更 API | `runner` / `api` | §11、§12、§22.0e、§27.18 | `.branch_config`、GET/POST API、runner 再起動不要条件が一致する。 |
| 週次ビルドサマリー Webhook | `runner` / `api` | §12、§13、§16、§22.0e、§27.19 | 週次判定、集計対象、通知 payload、手動送信 API が一致する。 |
| 設定変更の詳細 diff 記録 | `api` | §22.0a、§22.0e、§27.20 | `.config_log` の diff 文字列、対象 API、マスク条件が一致する。 |
| テーブルのソート機能 | `builder` | §7.14 | クリック操作、昇順/降順、`aria-sort`、インジケーターが一致する。 |
| キーボードショートカット | `builder` | §7.12 | `/`、`Escape`、`t` の対象、フォーカス条件、入力中の無効化が一致する。 |
| 複数ファイル監視 | `runner` / `builder` / `api` | §11、§12、§13、§15、§27.21 | `target_files` の検証、対象別 SHA 差分、build target 決定、履歴・ログ・API 表示が一致する。 |
| ビルドパイプライン YAML 定義 | `runner` / `api` | §12、§13、§15、§22.0e、§27.22 | `.pipeline.yml` の内製 subset parse、step 実行順、timeout、env、失敗時 status、API 保存が一致する。 |
| ローカルファイル監視モード | `runner` | §11、§12、§13、§27.23 | GitHub API を呼ばず、local snapshot の SHA-256 差分だけで変更検出し、trigger と status が一致する。 |
| タグ付きコミットのみビルド | `runner` / `api` | §12、§13、§15、§22.0e、§27.24 | tag pattern、GitHub tags API、skip 条件、build log/history、設定 API が一致する。 |
| ビルドキャッシュ | `builder` / `runner` | §5、§8、§11、§13、§27.25 | `.build_cache.json`、入力 manifest、再利用条件、無効化条件、report counters が一致する。 |
| 並列マルチターゲットビルド | `runner` | §12、§13、§14a、§15、§27.26 | worker 上限、target 別 status、pending transfer、最終 build status、ログ順序が一致する。 |
| ビルド前後フック | `runner` / `api` | §13、§15、§22.0e、§27.27 | `.hooks` schema、pre/post 実行、abort 条件、hook log、API CRUD が一致する。 |
| 依存ファイルトラッキング | `builder` / `runner` | §4.3、§5、§11、§13、§27.28 | `.dependency_manifest.json`、依存抽出、関連 target 判定、破損時 full build が一致する。 |
| リモートビルド対応 | `runner` / `api` | §12、§13、§14a、§15、§27.29 | remote command、archive 取得、manifest 検証、状態記録、失敗時 rollback 不実行が一致する。 |
| ビルド承認フロー | `runner` / `api` / `sdk` / `ui` | §11、§13、§15、§16、§22.0e、§27.30 | `.approval_queue`、承認/却下 API、通知、timeout、UI 操作、履歴 status が一致する。 |
| ブランチ別環境変数 | `runner` / `api` | §12、§13、§15、§22.0e、§27.31 | branch env schema、許可 key、secret mask、process env 注入、API 保存が一致する。 |
| ビルド通知連携 | `runner` / `api` / `sdk` / `ui` | §13、§15、§16、§22.0e、§27.32 | 通知 event、channel schema、送信順、retry、mask、notify log、API / UI 表示が一致する。 |
| ビルド時間トレンド記録 | `runner` / `api` | §13、§15、§22.0e、§27.33 | `.build_trends.json`、移動平均、中央値、p95、API response、破損時復旧が一致する。 |
| ビルド依存チェーン | `runner` / `api` | §11、§13、§15、§22.0e、§27.34 | `.build_chain_config`、依存 DAG 検証、実行順、skip / failure status、chain log が一致する。 |
| ビルド優先度キュー | `runner` / `api` | §11、§13、§22.0e、§27.35 | queue priority、created_seq、同一優先度 FIFO、API 表示、cancel / clear が一致する。 |
| 失敗原因の自動分類 | `runner` / `api` | §13、§15、§22.0e、§27.36 | failure_category、evidence、分類優先順位、history / log / UI 表示が一致する。 |
| ビルド実行環境の記録 | `runner` | §13、§15、§27.37 | build 開始時の environment snapshot、secret 非含有、log schema、検証 fixture が一致する。 |
| ビルド所要時間の異常検知 | `runner` / `api` | §13、§15、§16、§22.0e、§27.38 | trend 基準、異常判定、WARN、history flag、通知 payload、設定値が一致する。 |
| ビルドトリガー専用 API スコープ | `api` / `sdk` / `ui` | §22.0a、§22.0c、§22.0e、§23、§24、§27.42 | `trigger` scope token が build 起動系だけを許可し、その他 API を拒否する。 |
| API キー管理 | `api` / `sdk` / `ui` | §22.0a、§22.0c、§22.0e、§23、§24、§27.43 | API key 本体の一回表示、hash 保存、scope、期限、失効、監査が一致する。 |
| 監査ログ | `api` / `sdk` / `ui` | §22.0a、§22.0c、§22.0e、§23、§24、§27.44 | `.audit_log` schema、対象操作、mask、検索 API、UI 表示が一致する。 |
| セッションタイムアウト変更設定 | `api` / `sdk` / `ui` | §22.0c、§22.0e、§23、§24、§25、§27.45 | `session_timeout_seconds` の範囲、保存、既存 session の扱い、新規 session 期限が一致する。 |
| TOTP 二要素認証 | `api` / `sdk` / `ui` | §22.0a、§22.0c、§22.0e、§23、§24、§25、§27.46 | RFC 6238 TOTP、二段階 login、secret 保存、確認、無効化、UI 操作が一致する。 |
| API レート制限 | `api` / `sdk` / `ui` | §22.0a、§22.0c、§22.0e、§23、§24、§27.47 | IP / actor / endpoint group の固定窓制限、`429`、状態保存、設定 API、UI 表示が一致する。 |

---

## 0. システム概要

Adlaire CI は Go 版コンポーネントと JavaScript/HTML 管理ツールで構成する。

Part 3 詳細仕様セットは、`builder`、`runner`、`api`、`admin`、`sdk`、`ui`、`setup`、`statefile`、`archive`、`commitstatus`、`security`、`fixture` の実装詳細を責務 component 別に定義する。`mcp` の入出力、状態、起動手順、検証条件は定義しない。

**`builder`（ビルドスクリプト）**
GitHub リポジトリ上またはローカル上の Markdown ファイルまたは Markdown ディレクトリを静的 Web サイトに変換してローカルディレクトリへ出力する。標準実行バイナリ名は `adlaire-ci-build` とする。

**`runner`（CI ランナー）**
GitHub の Git Trees API / Git Blobs API を使用し、対象ファイルの blob SHA 変更を検出する。変更があった場合のみ Markdown 本文を書き出し、`adlaire-ci-build` を起動し、成功時に SHA キャッシュを更新する。systemd タイマーで定期実行する oneshot 設計。

SSH 転送、ペンディングキュー、スナップショット、Webhook 通知、マルチブランチ、ビルドログ保存、サーキットブレーカーは `runner` の対象機能である。

**`api`（管理 API サーバー）**
Go 標準ライブラリ `net/http` を使用する常駐 HTTP サーバー。管理ツールからの API リクエストを受け付け、認証・状態取得・手動ビルドトリガーを処理する。`adlaire-ci-api.service` として systemd に登録し、`runner` とは独立して常駐する。

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

標準配置への移行完了前は、`components/builder.go` の現行実装実体を `build_spec.go`、`components/runner.go` の現行実装実体を `runner.go` として扱う。詳細仕様内で `components/builder.go` または `components/runner.go` から生成すると記載するバイナリは、移行完了前に限り、それぞれ `build_spec.go` または `runner.go` から生成する同等バイナリとして扱う。

本節の tree は標準配置の最終形を示す。現時点で `将来計画` または `仕様化済み・未実装` の path は、該当 owner component が実装対象になった PR で追加する。標準配置図に含まれていることだけを理由に、未実装ファイル、将来計画ファイル、空ディレクトリ、placeholder を作成してはならない。

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

| パス | component | 役割 |
|------|-----------|------|
| `main.go` | `-` | 起動入口。サブコマンド判定、引数受け取り、責務 component 呼び出しを行う。 |
| `components/builder.go` | `builder` | Markdown / Markdown ディレクトリを静的 Web サイトへ変換する。 |
| `components/runner.go` | `runner` | GitHub polling、変更検出、ビルド起動、履歴、ログ、deploy を実行する。 |
| `components/api.go` | `api` | 管理 API サーバー、認証、状態ファイル操作を提供する。 |
| `components/admin.go` | `admin` | 管理 UI 静的ファイルの配布物構成、配置、検証、HTTP 静的配信境界を扱う。詳細は `ADLAIRE_CI_DETAIL_ADMIN_SPEC.md` を正とする。 |
| `components/statefile.go` | `statefile` | `.build_history`、`.build_logs`、`.server_config` など状態ファイルの読み書きを扱う。詳細は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` を正とする。 |
| `components/archive.go` | `archive` | ビルドログ圧縮、snapshot、配布アーカイブを扱う。 |
| `components/commitstatus.go` | `commitstatus` | GitHub Commit Status API 送信を扱う。 |
| `components/mcp.go` | `mcp` | 将来計画。現時点では実装可能な詳細仕様を持たず、MCP 専用詳細仕様が新設されるまで実装対象ではない。 |
| `admin/` | `ui` / `sdk` | 標準管理 UI の静的ファイルを配置する。 |
| `testdata/` | `fixture` | コンポーネント別 fixture を配置する。 |
| `docs/examples/` | `-` | 利用例、設定例、サンプル構成を配置する。 |

`main.go` は 1 ファイルとし、実装詳細を含めない。`components/` 配下は 1 コンポーネント = 1 Go ファイルとし、各ファイルは上表の責務を実装する。

標準配置への移行完了条件は以下に固定する。

| 対象 | 移行完了条件 |
|------|--------------|
| `main.go` | repository root に 1 ファイルだけ存在し、サブコマンド判定、引数受け取り、owner component 呼び出しだけを持つ。Markdown 変換、CI 実行、HTTP handler、状態ファイル操作、archive、commitstatus、MCP の実装詳細を含まない。 |
| `components/*.go` | 実装対象 owner component ごとに 1 Go ファイルだけ存在する。`builder` は `components/builder.go`、`runner` は `components/runner.go`、`api` は `components/api.go`、`admin` は `components/admin.go`、`statefile` は `components/statefile.go`、`archive` は `components/archive.go`、`commitstatus` は `components/commitstatus.go` とする。 |
| 標準移行前ファイル | `build_spec.go`、`runner.go`、`build_spec_test.go`、`runner_test.go`、`testdata/build_spec/` は、対応する標準配置へ移動済みであり、同じ実装本文または同じ fixture が旧配置に残っていない。 |
| testdata | 実装済みまたは仕様化済み・未実装の owner component ごとに `testdata/<component>/` を使用する。`testdata/mcp/` は MCP 専用詳細仕様が新設されるまで作成しない。 |
| admin | `admin/index.html` と `admin/adlaire-ci-sdk.js` は、それぞれ `ui` と `sdk` の owner 詳細仕様に従う。`admin/style.css` と `admin/app.js` は、`ADLAIRE_CI_DETAIL_ADMIN_SPEC.md` A1 に定義された任意配布物として扱い、未定義の admin 静的ファイルを追加しない。 |
| 将来計画 | `components/mcp.go` と MCP 用 fixture は、MCP 専用詳細仕様が新設され、`ADLAIRE_CI_SPEC.md` で `仕様化済み・未実装` へ昇格するまで作成しない。 |

標準配置へ移行する PR は、旧配置名と標準配置名の両方が同じ実装実体として併存していないこと、`DOCUMENT_INDEX.md` の Specified Components、`ADLAIRE_CI_SPEC.md` の実装状態、該当 owner component の詳細仕様、testdata 参照が同じ配置を指すことを確認する。

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
| runner owner の §27 個別節（§27.1 を除く） | `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` §27 |

`ADLAIRE_CI_DETAIL_SPEC.md` §27.38a は横断補足契約として本ファイルに残す。

§27.21〜§27.38 の runner 拡張機能を実装する場合は、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` の個別節を正本とし、横断する処理順、状態ファイル保存責務、api / sdk / ui 連動条件、受け入れ fixture の同期確認として本ファイル §27.38a を同時に確認する。

### — 管理ツール —

## 21. API 詳細仕様（分割済み）

`api` owner component の詳細仕様本文は `ADLAIRE_CI_DETAIL_API_SPEC.md` を正とする。

| 移動済み範囲 | 移動先 |
|--------------|--------|
| §21〜§22 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §21〜§22、§21a |
| §25 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §25 |
| api owner の §27 個別節 | `ADLAIRE_CI_DETAIL_API_SPEC.md` §27 |

`ADLAIRE_CI_DETAIL_SPEC.md` §27.38a は横断補足契約として本ファイルに残す。

§27.21〜§27.38 に関わる api / sdk / ui 連動条件は、各 owner component の分割先詳細仕様ファイルを正本とし、本ファイル §27.38a は横断同期確認として同時に確認する。§27.38a の内容を api / sdk / ui の分割先へ重複定義してはならない。

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

本節は、追加仕様化機能の詳細仕様参照である。各機能の主本文は、owner component の分割先詳細仕様ファイルを正とする。本節に定義された機能は、§0i と §27.0 を入口として、owner 詳細仕様の主本文と必要な collaborator 詳細仕様の確認項目を組み合わせて実装可否を判定する。

### 27.0 追加仕様化機能 共通実装契約

§27 の主本文は、owner component の分割先詳細仕様ファイルを正とする。親ファイルでは、§27 の実装時に共通して確認する参照順、越境禁止、PR 証跡の入口だけを定義する。

| 確認 | 固定内容 |
|------|----------|
| 実装対象判定 | `ADLAIRE_CI_SPEC.md` で実装状態と実装可否を確認し、将来計画、実装不可、未仕様化、MCP 専用機能を実装対象にしない。 |
| owner 確定 | §0b と §0i で owner component を 1 件に確定し、主本文は owner の分割先詳細仕様ファイルで確認する。 |
| collaborator 確認 | api / sdk / ui / statefile / security / archive / commitstatus などの collaborator がある場合は、該当分割先ファイルの参照節を schema、呼び出し境界、表示、security、setup、fixture、検証観点として読む。 |
| 補完禁止 | 個別節または分割先詳細仕様に存在しない endpoint、状態ファイル、設定 key、UI 操作、SDK method、外部依存を実装判断で追加しない。 |
| 状態更新 | 状態ファイル更新は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0a、§22.0c を正とし、lock、atomic write、JSON Lines、破損時処理を独自定義しない。 |
| security | secret mask、token、session、scope、audit、rate limit は `ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` を正とし、平文保存・平文表示を行わない。 |
| fixture / PR 証跡 | fixture manifest、expected/effects、assertion、PR 証跡、受け入れゲートは `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §27-F を正とする。 |
| api / sdk / ui 同期 | API endpoint、SDK method、UI 操作が同一機能に関わる場合は、endpoint は `ADLAIRE_CI_DETAIL_API_SPEC.md`、SDK method は `ADLAIRE_CI_DETAIL_SDK_SPEC.md`、UI 操作は `ADLAIRE_CI_DETAIL_UI_SPEC.md` をそれぞれ正本とし、名称、引数、response、error、表示、成功後再取得、失敗時固定が食い違わないことを確認する。 |

§27 の機能を実装した PR は、対象節、owner 詳細仕様、collaborator 詳細仕様、fixture、secret mask、失敗時副作用、実装対象外を PR 本文に記録する。記録が不足する場合は、実装完了として扱わない。

§27 の PR 分割、dry-run 固定契約、fixture 完了条件の詳細は、owner 詳細仕様と `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §27-F を正とする。親ファイルに同じ fixture schema、expected/effects、個別機能本文を重複定義しない。

### 27.1 GitHub Commit Status API

本節の主本文は `ADLAIRE_CI_DETAIL_COMMITSTATUS_SPEC.md` §27.1 を正とする。owner component は `commitstatus`、collaborator component は `runner`、`statefile` とする。runner の build 実行、commit SHA 確定、build id 採番、pipeline / deploy / snapshot / history の最終結果確定は `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` を正とする。

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

### 27.38a 横断連動・Runner 拡張機能 実装補足契約

本節は、責務 component 別詳細仕様へ分割しない。§27.21〜§27.38 および api / sdk / ui / statefile の横断連動は runner、builder、api、sdk、ui、statefile、archive にまたがるため、実装者は owner 詳細仕様を主本文とし、本節を横断確認として同時に確認する。

本節の正本範囲は、横断確認、同期禁止、横断処理順、成功後再取得、失敗時固定、実装完了時の横断受け入れ観点に限定する。本節は、個別機能の処理本文、入力、出力、状態 schema、fixture schema、endpoint 詳細、SDK method 詳細、UI DOM 詳細を持たない。これらは各 owner / collaborator の分割先詳細仕様ファイルを正とする。

本節と owner component 別詳細仕様ファイルの内容が矛盾する場合は、個別機能の入出力、状態、処理、異常系、endpoint、SDK、UI、fixture は owner component 別詳細仕様ファイルを正とし、横断処理順、API / SDK / UI / statefile 同期、成功後再取得、失敗時固定だけを本節で確認する。§27.38a を理由に、owner 詳細仕様に存在しない endpoint、SDK method、UI 操作、状態ファイル、設定 key、fixture を追加してはならない。

| 確認 | 固定内容 |
|------|----------|
| runner 起点 | §27.21〜§27.38 の多くは runner の build 実行、queue、history、log、notification に影響するため、`ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` の該当 §27 節を先に確認する。 |
| builder 連携 | cache、dependency、output meta、生成物に関わる場合は `ADLAIRE_CI_DETAIL_BUILDER_SPEC.md` の該当 §27 節を同時に確認する。 |
| API 連携 | 設定保存、queue、approval、history、stats、snapshot、rollback、hook、notify、search に関わる場合は `ADLAIRE_CI_DETAIL_API_SPEC.md` の endpoint / state read-write 契約を同時に確認する。 |
| SDK / UI 連携 | API を管理画面から操作する機能は、SDK method は `ADLAIRE_CI_DETAIL_SDK_SPEC.md` §23、DOM / 表示条件は `ADLAIRE_CI_DETAIL_UI_SPEC.md` §24 を正本として確認する。 |
| statefile | 状態 schema、lock、atomic write、JSON Lines、破損時処理、保存順は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` を正とする。 |
| archive | snapshot、artifact、download、delete、rollback、log archive は `ADLAIRE_CI_DETAIL_ARCHIVE_SPEC.md` を正とする。 |
| fixture | §27.21〜§27.38 の受け入れ fixture、secret mask、effects、PR 証跡は `ADLAIRE_CI_DETAIL_FIXTURE_SPEC.md` §27-F を正とする。 |

**api / sdk / ui / statefile 横断連動契約：**

下表は横断確認表であり、新しい API endpoint、SDK method、UI 操作、状態ファイル副作用を定義する表ではない。下表の機能群は、各 owner component 別詳細仕様ファイルに定義済みの API endpoint、SDK method、UI 操作、状態ファイル副作用を同じ実装単位でそろえる。API だけ、SDK だけ、UI だけを先行して仕様外の仮実装にしてはならない。UI が未実装の Phase では、UI 列は fixture の期待操作として固定し、実装完了扱いには含めない。

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

**横断処理順契約：**

| 処理種別 | 固定順序 |
|----------|----------|
| 認証必須 JSON API | method / path 判定 → body 禁止判定 → JSON parse → 認証 / scope → rate limit → endpoint 固有 validation → read → write 計画 → atomic write → JSON Lines 追記 → response。 |
| read-only API | method / path 判定 → body 禁止判定 → 認証 / scope → query validation → read → 壊れた任意行除外 → response。read-only API は状態ファイルを書き換えない。 |
| UI 変更操作 | panel error / success 消去 → UI 入力検証 → 対象操作 disabled → SDK 呼び出し → 成功後再取得 → success 表示 → secret 消去 → disabled 再評価。 |
| UI 取得操作 | panel error 消去 → 対象操作 disabled → SDK 呼び出し → DOM 更新 → empty state 判定 → disabled 再評価。success 表示は行わない。 |
| SDK request | 引数検証 → path / query / body 生成 → Authorization 付与 → timeout 設定 → fetch → status 判定 → response parse → token 変化適用 → return / throw。 |
| multi-file write | 全入力検証 → 全対象 read → 全 write payload 生成 → `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0d の Write 順に atomic write → JSON Lines 追記 → response。途中失敗時は未処理ファイルを書かない。 |

§27.21〜§27.38 の実装では、owner 詳細仕様にない状態ファイル、endpoint、SDK method、UI 操作、外部公開構成を追加してはならない。追加が必要な場合は、親ファイルではなく、該当 owner / collaborator の分割先詳細仕様を先に改訂する。

### 27.42 ビルドトリガー専用 API スコープ

本節の主本文は `ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` §27.42 を正とする。owner component は `security`、collaborator component は `api`、`sdk`、`ui`、`statefile` とする。

### 27.43 API キー管理

本節の主本文は `ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` §27.43 を正とする。owner component は `security`、collaborator component は `api`、`sdk`、`ui`、`statefile` とする。

### 27.44 監査ログ

本節の主本文は `ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` §27.44 を正とする。owner component は `security`、collaborator component は `api`、`statefile` とする。

### 27.45 セッションタイムアウト変更設定

本節の主本文は `ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` §27.45 を正とする。owner component は `security`、collaborator component は `api`、`sdk`、`ui`、`statefile` とする。

### 27.46 TOTP 二要素認証

本節の主本文は `ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` §27.46 を正とする。owner component は `security`、collaborator component は `api`、`sdk`、`ui`、`statefile` とする。

### 27.47 API レート制限

本節の主本文は `ADLAIRE_CI_DETAIL_SECURITY_SPEC.md` §27.47 を正とする。owner component は `security`、collaborator component は `api`、`sdk`、`ui`、`statefile` とする。
