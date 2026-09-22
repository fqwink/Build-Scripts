# Adlaire CI — Roadmap

[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務は、Adlaire CI の状態分類、実装可否、Phase、機能インベントリ、将来計画、昇格手順、追加仕様化機能参照、横断補足契約を管理するロードマップ正本である。

方針、ポリシー、正本参照先、禁止事項、リリース判断は [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務、生成 HTML のデザイン関係は [`docs/DESIGN.md`](DESIGN.md) デザイン責務、詳細仕様参照入口は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務、owner component 別の詳細本文は [`docs/details/*.md`](details/) 詳細本文責務を参照する。

HTTP response schema、状態ファイル schema、SDK method の実装詳細、UI DOM、fixture assertion、具体的な処理順序の正本ではない。

**状態・計画責務 共通参照先：**

[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務で実装詳細、HTTP endpoint、schema、DOM、SDK method、fixture、検証証跡に触れる場合、実装詳細は owner component 別の [`docs/details/*.md`](details/) 詳細本文責務、fixture / expected / fake / 実装検証証跡は [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務、文書・実装ファイル所在は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務を共通参照先とする。状態・計画責務では状態、実装可否、Phase、将来計画、昇格入口だけを本文として扱う。

## 1. ロードマップ責務

| 管理対象 | 状態・計画責務で扱う内容 | 詳細参照先 |
|----------|----------------------|--------------|
| 状態分類 | component ごとの `未仕様化`、`将来計画`、`改訂予定`、`仕様化済み・未実装`、`実装中・検証未完了`、`実装済み` の分類。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 共通参照先に従う。 |
| Phase | Phase 順序、対象 owner component、依存条件、判定条件、引き継ぎ契約。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 共通参照先に従う。 |
| 機能インベントリ | 機能の分類、実装可否、詳細仕様参照先。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 共通参照先に従う。 |
| 将来計画 | 実装不可の構想と昇格入口。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 共通参照先に従う。 |
| 追加仕様化機能参照・横断補足契約 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6 の追加仕様化機能 owner、主本文、collaborator、横断受け入れ観点。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 共通参照先に従う。 |

## 2. 状態分類

| 状態 | 実装可否 | 意味 | 実装者の扱い |
|------|----------|------|--------------|
| 実装済み | 完了済み | ソースコード実装と必須検証が完了した項目。 | 関連する実装ファイル、検証結果、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務を確認する。 |
| 実装中・検証未完了 | 検証待ち | 実装に着手済みだが、必須検証または証跡が未完了の項目。 | 未完了の確認対象は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 共通参照先に従う。 |
| 仕様化済み・未実装 | 実装可 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務と owner component 別の [`docs/details/*.md`](details/) 詳細本文責務に実装可能な詳細が揃っている項目。 | Phase と詳細仕様参照先を確認する。 |
| 改訂予定 | 実装不可 | 将来計画から格上げ済みだが、詳細仕様作成中の項目。 | 詳細仕様の作成先と不足項目を確認する。 |
| 将来計画 | 実装不可 | 構想として保持するが、実装契約がない項目。 | 昇格手順は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 を確認する。 |
| 未仕様化 | 実装不可 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務と詳細仕様に存在しない項目。 | 状態分類の最終判定は [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0a を確認する。 |

## 3. コンポーネント状態分類

コンポーネント状態分類では、仕様化済みの内容と実装済みの内容を区別して扱う。

[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §3 のコンポーネント状態分類表のコンポーネント名は、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.3 のディレクトリ構成に基づく。`main.go`、`components/*.go`、`admin/` 配下の管理 UI 静的ファイル、`testdata/<component>/` を現行配置として扱う。

| コンポーネント | 状態 | 備考 |
|---------------|------|------|
| `components/builder.go` | 実装済み | Go 版 Markdown → 静的 Web サイトビルドスクリプトとして Phase 1 の `gofmt` と `go test` 検証済み。 |
| `components/runner.go` | 実装済み | Go 版 CI ランナーとして Phase 2 判定パスの `gofmt` と `go test` 検証済み。 |
| `components/api.go` | 実装済み | Go 版管理 API サーバー。API 完全契約表の endpoint、認証、状態ファイル、管理操作、検証テストを実装済み。 |
| `admin/adlaire-ci-sdk.js` | 実装済み | 管理ツール用 JavaScript SDK。Phase 5 の SDK class、method、HTTP error、token 破棄、query / body 生成を実装済み。 |
| `admin/index.html` | 実装済み | 標準管理ツール UI。Phase 6 の DOM id、panel、SDK 呼び出し、表示状態、秘密情報消去を実装済み。 |
| `components/mcp.go` | 将来計画 | Go 版 MCP サーバー。将来計画として管理し、実装済みとは扱わない。 |

実装済み判定の最終条件は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.8 と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0a を参照する。

## 4. Phase 実装計画

## 4.1 初期実装 Phase 単位

Go 版初期実装の対象範囲と Phase 単位の扱いは [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0e と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0f を参照する。Phase 実装計画は Phase 順、対象、依存条件、判定条件を管理する。

実装単位、実装変更単位、判定単位に関する禁止事項は [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0f を参照する。API の実装範囲は、Phase 3 を「API 基盤・認証・状態 read/write・運用基本操作」、Phase 4 を「API 拡張運用操作」として管理する。

各 Phase の `対象` は、その Phase の owner component を示す。状態ファイル、security、archive、commitstatus、admin、fixture、setup が関わる場合も、それらは collaborator component として該当 Phase の判定条件に含める。collaborator component の詳細本文責務に未充足がある場合の扱いは [`docs/SPEC.md`](SPEC.md) 方針責務 §4.7 と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0d を参照する。

| Phase | 対象 | 実装範囲 | 依存条件 | 判定条件 |
|-------|------|----------|----------|----------|
| Phase 1 | `builder` | [`docs/details/builder.md`](details/builder.md) §2〜§9 の CLI、Markdown 変換、静的 Web サイト出力、テーマコンポーネント、生成物確認。 | なし。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0e の `builder` 必須検証と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0f の `builder` 判定入口を満たす。 |
| Phase 2 | `runner` | [`docs/details/runner.md`](details/runner.md) §10〜§20 の CI ランナー、GitHub API 連携、SHA キャッシュ、pipeline 起動、SSH 転送、snapshot、通知、ログ、systemd / setup 参照境界。systemd unit 本文と配置手順は [`docs/details/setup.md`](details/setup.md) §26 を参照する。 | Phase 1 が完了し、`adlaire-ci-build` の CLI 契約が固定されている。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0e の `runner` 必須検証と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0f の `runner` 判定入口を満たす。 |
| Phase 3 | `api` | [`docs/details/api.md`](details/api.md) §21〜§22、[`docs/details/api.md`](details/api.md) §21a、[`docs/details/api.md`](details/api.md) §25 と [`docs/details/setup.md`](details/setup.md) §26 のうち、認証、セッション、共通エラー、状態ファイル読み書き、ビルド操作、status、logs、history、queue、circuit breaker。 | Phase 2 が完了し、runner が書き込む状態ファイル schema が固定されている。 | API 基盤・認証・状態 read/write・運用基本操作の必須検証、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0e の `api` API 共通・状態ファイル検証、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0f の `api` 判定入口の該当範囲を満たす。 |
| Phase 4 | `api` | [`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/api.md`](details/api.md) §22.0f のうち、config、repo、branch、schedule、notify、snapshot、rollback、maintenance、access control、hooks、tokens 等の拡張運用操作。 | Phase 3 が完了し、API 共通処理と認証が固定されている。 | API 拡張運用操作の必須検証と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0e の `api` endpoint 契約を満たす。 |
| Phase 5 | `sdk` | [`docs/details/sdk.md`](details/sdk.md) §23 の SDK class、method、戻り値、HTTP error、token 破棄、query / body 生成。 | Phase 3 と Phase 4 が完了し、[`docs/details/api.md`](details/api.md) §22.0e の endpoint 契約が固定されている。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0e の `sdk` 契約と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0f の `sdk` 判定入口を満たす。 |
| Phase 6 | `ui` | [`docs/details/ui.md`](details/ui.md) §24 の標準管理ツール UI、DOM id、panel、操作、SDK 呼び出し、成功表示、失敗表示、disabled、再取得、秘密情報消去。 | Phase 5 が完了し、SDK method 契約が固定されている。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0e の `ui` 契約と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0f の `ui` 判定入口を満たす。 |

### 4.1.1 Phase 1 完全仕様ゲート（`builder`）

Phase 完全仕様ゲートは、対象 Phase の対象、依存条件、判定条件、後続 Phase への引き継ぎ確認だけを扱う。各 Phase の実装詳細本文は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 共通参照先と下表の参照先を正とする。

Phase 1 完全仕様ゲートでは、Phase 1 の対象、依存条件、判定条件、後続 Phase への引き継ぎ確認だけを扱う。

| 確認 | 参照先 |
|------|--------|
| CLI、Markdown 変換、静的 Web サイト出力、theme component、生成物確認 | [`docs/details/builder.md`](details/builder.md) §1〜§9、[`docs/details/builder.md`](details/builder.md) §8a |
| Phase 1 fixture、testdata、実装検証証跡 | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F |
| release / setup 受け入れ条件 | [`docs/details/setup.md`](details/setup.md) §26.7 |

### 4.1.2 Phase 2 完全仕様ゲート（`runner`）

Phase 2 完全仕様ゲートでは、Phase 2 が Phase 1 の `adlaire-ci-build` 契約に依存し、後続 API が読む runner 状態契約を固定することだけを扱う。

| 確認 | 参照先 |
|------|--------|
| CI runner、GitHub API 連携、SHA cache、pipeline、deploy、snapshot、通知、systemd | [`docs/details/runner.md`](details/runner.md) §10〜§20、[`docs/details/runner.md`](details/runner.md) §15a |
| runner fixture、fake GitHub、fake ssh / notifier、実装検証証跡 | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F |
| release / setup 受け入れ条件 | [`docs/details/setup.md`](details/setup.md) §26.7 |

### 4.1.3 Phase 3 完全仕様ゲート（`api`）

Phase 3 完全仕様ゲートでは、API 共通契約、認証、状態 read/write、運用基本 endpoint が Phase 4〜6 の前提になることだけを扱う。

| 確認 | 参照先 |
|------|--------|
| API 共通処理、API server 制限、認証、運用基本 endpoint、状態 read/write | [`docs/details/api.md`](details/api.md) §21〜§22、[`docs/details/api.md`](details/api.md) §21a、[`docs/details/api.md`](details/api.md) §25、[`docs/details/api.md`](details/api.md) §22.0f |
| 状態ファイル schema、lock、atomic write | [`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/statefile.md`](details/statefile.md) §22.0c |
| 管理 API 導入、systemd、release 受け入れ条件 | [`docs/details/setup.md`](details/setup.md) §26 |

### 4.1.4 Phase 4 完全仕様ゲート（`api`）

Phase 4 完全仕様ゲートでは、拡張運用 endpoint が SDK / UI の最終入力契約になることだけを扱う。

| 確認 | 参照先 |
|------|--------|
| 拡張運用 endpoint、request / response、error、auth、secret mask | [`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/api.md`](details/api.md) §22.0f |
| security 連携 | [`docs/details/security.md`](details/security.md) §27.42〜§27.47 |
| API fixture、endpoint 証跡 | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F、[`docs/details/fixture.md`](details/fixture.md) §22-F、[`docs/details/fixture.md`](details/fixture.md) §27-F |

### 4.1.5 Phase 5 完全仕様ゲート（`sdk`）

Phase 5 完全仕様ゲートでは、SDK が固定済み API endpoint だけを呼び、UI 表示決定を持たないことだけを扱う。

| 確認 | 参照先 |
|------|--------|
| SDK class、method、HTTP 対応、query / body 生成、error、stream、token 破棄 | [`docs/details/sdk.md`](details/sdk.md) §23 |
| API endpoint 対応 | [`docs/details/api.md`](details/api.md) §22.0e |
| SDK fixture、fake fetch / stream | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F |

### 4.1.6 Phase 6 完全仕様ゲート（`ui`）

Phase 6 完全仕様ゲートでは、UI が SDK 経由だけで API と通信し、秘密情報を DOM に残さないことだけを扱う。

| 確認 | 参照先 |
|------|--------|
| DOM id、panel、操作、表示状態、SDK 呼び出し、秘密情報消去 | [`docs/details/ui.md`](details/ui.md) §24 |
| SDK method 契約 | [`docs/details/sdk.md`](details/sdk.md) §23 |
| UI fixture、fake SDK、直接 API 呼び出し禁止確認 | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F |

### 4.1.7 Phase 間引き継ぎ契約

各 Phase の完了時は、次 Phase が依存する契約を変更不可として扱う。後続 Phase で変更が必要になった場合は、後続 Phase の実装で吸収せず、契約を定義した owner component 別の [`docs/details/*.md`](details/) 詳細本文責務へ戻す。

| 引き継ぎ元 | 引き継ぎ先 | 固定する契約 | 参照先 |
|------------|------------|--------------|------|
| Phase 1 | Phase 2 | `adlaire-ci-build` CLI、終了コード、stdout / stderr、`[REPORT]`、出力サイト構造。 | [`docs/details/builder.md`](details/builder.md) |
| Phase 2 | Phase 3 | runner 状態ファイル schema、lock、history/log、pending queue、circuit breaker、snapshot、通知ログ。 | [`docs/details/runner.md`](details/runner.md) / [`docs/details/statefile.md`](details/statefile.md) |
| Phase 3 | Phase 4 | API 共通契約、認証、session、error body、validation、lock error、SSE 基本形式。 | [`docs/details/api.md`](details/api.md) |
| Phase 4 | Phase 5 | 全 endpoint の method、path、query、request、response、error、認証要否。 | [`docs/details/api.md`](details/api.md) |
| Phase 5 | Phase 6 | SDK method 名、引数、戻り値、error object、stream handle、token 破棄条件。 | [`docs/details/sdk.md`](details/sdk.md) |
| Phase 6 | 初期実装確認 | UI 操作、表示状態、secret 消去、SDK 経由通信、実装確認結果。 | [`docs/details/ui.md`](details/ui.md) |

### 4.1.8 Phase 別 実装成果物チェックリスト

Phase fixture / testdata 配置、fake 実装、実装検証証跡の詳細は [`docs/details/fixture.md`](details/fixture.md) §0g.8-F を参照する。Phase 別実装成果物チェックリストは、Phase ごとの成果物参照先だけを保持し、fixture 名、expected / effects、fake 動作、実装検証証跡項目の正本ではない。

| Phase | 実装対象 | 成果物・fixture 確認先 | 受け入れ条件 |
|-------|----------|----------------------|--------------|
| Phase 1 | `builder` | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F | [`docs/details/builder.md`](details/builder.md) と [`docs/details/setup.md`](details/setup.md) §26.7 を満たす。 |
| Phase 2 | `runner` | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F | [`docs/details/runner.md`](details/runner.md) と [`docs/details/setup.md`](details/setup.md) §26.7 を満たす。 |
| Phase 3 | `api` | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F、[`docs/details/fixture.md`](details/fixture.md) §27-F | [`docs/details/api.md`](details/api.md) の API 基盤・認証・状態 read/write・運用基本操作と setup API 導入条件を満たす。 |
| Phase 4 | `api` | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F、[`docs/details/fixture.md`](details/fixture.md) §27-F | [`docs/details/api.md`](details/api.md) の API 拡張運用操作を満たす。 |
| Phase 5 | `sdk` | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F | [`docs/details/sdk.md`](details/sdk.md) §23 を満たす。 |
| Phase 6 | `ui` | [`docs/details/fixture.md`](details/fixture.md) §0g.8-F | [`docs/details/ui.md`](details/ui.md) §24 を満たす。 |

---


## 5. 機能インベントリと統合ロードマップ

## 5.1 機能一覧

[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.1 は、状態・計画責務が管理する機能インベントリである。各機能の仕様詳細は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務と owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を参照する。状態分類は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §2 と [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §3、および [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0a の仕様成熟度ポリシーに従って判定する。

### ビルド・CI ランナー（components/runner.go）

**Go 版で実装済みの初期範囲：**

- `--state-dir`、`--once`、`--version`、`--help` の CLI 契約
- `.branch_config` と `.last_sha` による branch target / SHA cache 読み込み
- GitHub Trees API / Blobs API による対象 Markdown 取得
- SHA 一致時の変更なし skip
- `pipeline.sh` 起動、stdout/stderr 収集、`[REPORT]` / `[WARN]` 取り込み
- `.build_logs/{id}.json`、`.build_history`、`.build_status.json`、`.build_state`、`.build_lock` の作成・更新
- deploy 失敗時の `.pending_transfers` 追加
- `.notify_pending` 破損時の退避と `[]` 再生成
- GitHub API retry / rate limit 待機
- `adlaire-ci-build` 実行可否と disk 空き容量の precheck
- `.snapshots/{build_id}/site` 保存と世代 pruning
- `.notify_pending` の HTTP 再送と成功時削除
- `.build_circuit_state` による circuit open skip と失敗回数記録
- `.pending_transfers` の起動時再試行と成功時削除
- SSH 転送時の checksum 比較、未変更ファイル skip、転送後 checksum 検証
- 複数 branch target の順次処理
- `BUILD_COOLDOWN_SECONDS` による cooldown skip
- `FORCE_BUILD_INTERVAL` による変更なし時の定期強制ビルド
- GitHub commits API によるトリガー commit 情報の build log 記録
- GitHub PAT 有効期限ヘッダーの 7 日以内 WARN ログ
- `.notify_config` に基づく成功・失敗・転送失敗 Webhook 通知送信
- `OUTPUT_SIZE_WARN_MB` による出力サイトサイズ警告
- Phase 2 fixture R1〜R7 と判定パステスト

**Go 版で仕様化済みの全体範囲：**

- GitHub リポジトリの対象ファイルを定期ポーリング（systemd timer）
- blob SHA による差分検出（変更なし時はビルドをスキップ）
- Markdown → 静的 Web サイト変換（`adlaire-ci-build` を呼び出し）
- ビルド成功後に SHA キャッシュを更新する
- ビルド失敗時は SHA キャッシュを更新せず、次回起動時に再試行可能な状態を残す
- systemd oneshot ユニットとして動作（`adlaire-ci.service`）

- ビルド結果を `.build_history` に記録（ID 形式：`b{YYYYMMDDHHmmss}`）
- ビルドごとのログを `.build_logs/{id}.json` に保存
- ビルド成功・失敗時に Webhook 通知を送信（`.notify_config` を読み込み送信。送信責務は `components/runner.go`。通知 API は設定の読み書きのみ）
- ビルド成功後、出力サイトディレクトリを SSH 経由で静的コンテンツ配信サーバーへ転送する
- SSH 転送失敗時は `.pending_transfers` へキューイングし、次回起動時に自動再試行する
- 転送成功後、出力サイトディレクトリを `.snapshots/` へアーカイブし `HISTORY_KEEP_N` 世代を超過分から自動削除する
- SSH 転送失敗時に Webhook 通知を送信する（`on: ["deploy_failure"]` 設定時）
- `BRANCH_TARGETS` リストで複数ブランチを順次ポーリング・ビルドする
- `FORCE_BUILD_INTERVAL` 設定時、前回ビルドから指定時間経過で変更なしでも強制ビルドする

### 管理 API エンドポイント（components/api.go）

| カテゴリ | エンドポイント |
|---------|--------------|
| 認証 | `POST /api/login` / `POST /api/logout` / `POST /api/change-password` |
| 死活監視 | `GET /api/health` |
| ビルド操作 | `POST /api/build` / `POST /api/build/force` / `POST /api/build/cancel` / `GET /api/build/stream` / `POST /api/circuit-breaker/reset` |
| ステータス | `GET /api/status` / `GET /api/dashboard` |
| ログ | `GET /api/logs` / `GET /api/logs/export` / `GET /api/logs/search` / `POST /api/logs/cleanup` |
| ビルド履歴 | `GET /api/history` / `GET /api/history/{id}/log` / `GET /api/history/{id}/comment` / `POST /api/history/{id}/comment` / `GET /api/history/export` / `POST /api/history/{id}/flag` / `POST /api/history/{id}/tags` / `POST /api/history/{id}/rollback` |
| スケジュール | `GET /api/schedule` / `POST /api/schedule/interval` / `POST /api/schedule/pause` / `POST /api/schedule/resume` / `POST /api/schedule/allowed-hours` / `POST /api/schedule/force-interval` / `POST /api/schedule/cooldown` |
| 通知 | `GET /api/notify-config` / `POST /api/notify-config` / `POST /api/notify-test` / `GET /api/notify-log` / `POST /api/notify/weekly-summary` |
| システム情報 | `GET /api/sysinfo` / `GET /api/output-meta` / `GET /api/diagnostics` / `GET /api/rate-limit` / `GET /api/disk-usage` |
| 統計 | `GET /api/stats` / `GET /api/stats/timeline` / `GET /api/stats/build-duration` |
| リポジトリ | `GET /api/repo-info` / `POST /api/repo-config` / `GET /api/branch-config` / `POST /api/branch-config` |
| PAT 管理 | `GET /api/pat-status` / `POST /api/pat-verify` / `POST /api/pat-update` |
| 設定 | `GET /api/config` / `POST /api/config` / `POST /api/log-level` / `GET /api/config-log` |
| アクセスログ | `GET /api/access-log` |
| バックアップ | `GET /api/backup` / `POST /api/restore` |
| セッション管理 | `GET /api/sessions` / `POST /api/sessions/revoke-all` |
| API トークン | `GET /api/tokens` / `POST /api/tokens` / `DELETE /api/tokens/{id}` |
| スナップショット | `GET /api/snapshots` / `GET /api/snapshots/{id}/download` / `DELETE /api/snapshots/{id}` |
| メンテナンス | `GET /api/maintenance` / `POST /api/maintenance/enable` / `POST /api/maintenance/disable` |
| アクセス制御 | `GET /api/access-control` / `POST /api/access-control` |
| フック | `GET /api/hooks` / `POST /api/hooks` / `DELETE /api/hooks/{id}` / `GET /api/hooks/{id}/log` |
| アラートルール | `GET /api/alert-rules` / `POST /api/alert-rules` / `DELETE /api/alert-rules/{id}` |
| Webhook 受信 | `POST /api/webhook` / `GET /api/webhook-events` / `GET /api/webhook-config` / `POST /api/webhook-config` |
| 自動タグ付け | `GET /api/tag-rules` / `POST /api/tag-rules` / `DELETE /api/tag-rules/{id}` |
| パイプライン | `GET /api/pipeline-config` / `POST /api/pipeline-config` / `POST /api/verify-output` |
| 運用ノート | `GET /api/notes` / `POST /api/notes` |
| メール通知 | `GET /api/smtp-config` / `POST /api/smtp-config` / `POST /api/smtp-test` |
| ビルドキュー | `GET /api/queue` / `DELETE /api/queue` |
| ダッシュボードレイアウト | `GET /api/dashboard-layout` / `POST /api/dashboard-layout` |

### SDK メソッド（adlaire-ci-sdk.js）

| カテゴリ | メソッド |
|---------|--------|
| 認証 | `login()` / `logout()` / `changePassword()` |
| ビルド操作 | `triggerBuild()` / `buildForce()` / `cancelBuild()` / `streamBuild()` / `resetCircuitBreaker()` |
| ステータス | `getStatus()` / `getDashboard()` |
| ログ | `getLogs()` / `exportLogs()` / `searchLogs()` / `cleanupLogs()` |
| ビルド履歴 | `getHistory()` / `getHistoryLog()` / `getHistoryComment()` / `setHistoryComment()` / `exportHistory()` / `setHistoryFlag()` / `setHistoryTags()` / `rollbackHistory()` |
| スケジュール | `getSchedule()` / `setScheduleInterval()` / `pauseSchedule()` / `resumeSchedule()` / `setAllowedHours()` / `clearAllowedHours()` / `setForceInterval()` / `setBuildCooldown()` |
| Webhook 受信 | `getWebhookEvents()` / `getWebhookConfig()` / `setWebhookConfig()` |
| 通知 | `getNotifyConfig()` / `setNotifyConfig()` / `notifyTest()` / `getNotifyLog()` / `notifyWeeklySummary()` |
| システム情報 | `getSysinfo()` / `getOutputMeta()` / `getDiagnostics()` / `getRateLimit()` / `getDiskUsage()` |
| 統計 | `getStats()` / `getStatsTimeline()` / `getStatsBuildDuration()` |
| リポジトリ | `getRepoInfo()` / `setRepoConfig()` / `getBranchConfig()` / `setBranchConfig()` |
| PAT 管理 | `getPatStatus()` / `patVerify()` / `updatePat()` |
| 設定 | `getConfig()` / `setConfig()` / `setLogLevel()` / `getConfigLog()` |
| アクセスログ | `getAccessLog()` |
| バックアップ | `backup()` / `restore()` |
| セッション管理 | `getSessions()` / `revokeAllSessions()` |
| API トークン | `getTokens()` / `createToken()` / `revokeToken()` |
| スナップショット | `getSnapshots()` / `downloadSnapshot()` / `deleteSnapshot()` |
| メンテナンス | `getMaintenance()` / `enableMaintenance()` / `disableMaintenance()` |
| アクセス制御 | `getAccessControl()` / `setAccessControl()` |
| フック | `getHooks()` / `addHook()` / `deleteHook()` / `getHookLog()` |
| アラートルール | `getAlertRules()` / `addAlertRule()` / `deleteAlertRule()` |
| 自動タグ付け | `getTagRules()` / `addTagRule()` / `deleteTagRule()` |
| パイプライン | `getPipelineConfig()` / `setPipelineConfig()` / `verifyOutput()` |
| 運用ノート | `getNotes()` / `setNotes()` |
| メール通知 | `getSmtpConfig()` / `setSmtpConfig()` / `smtpTest()` |
| ビルドキュー | `getQueue()` / `clearQueue()` |
| ダッシュボードレイアウト | `getDashboardLayout()` / `setDashboardLayout()` |
| 死活監視 | `health()` |

ES Module・外部依存なし。全メソッドは `Promise` を返す。`streamBuild` は SSE 接続確立後に `Promise<StreamHandle>` として resolve し、`StreamHandle` は `{ close(): void, closed: boolean }` を持つ。`constructor` を除く合計は 96 メソッド。

### 標準管理ツール パネル（admin/index.html）

| パネル | 主な機能 |
|-------|---------|
| ログイン | パスワード認証 |
| パスワード変更 | 強制変更フロー対応（5 回目以降は他パネルを非表示） |
| ステータス | 最終ビルド情報・出力サイトリンク・メンテナンスバナー表示（モード中） |
| 手動実行 | ビルド起動・強制ビルド・キャンセル・キュー状態表示・キューのクリア |
| ログビューア | ログ閲覧・キーワードフィルター・ログレベルフィルター・JSON エクスポート・横断検索（期間指定） |
| ビルド履歴 | 過去ビルド一覧・タグ列・フラグ列・ページネーション・タグ/フラグフィルター・JSON エクスポート・個別ログ参照 |
| システム情報 | ファイルサイズ・稼働時間・ディスク使用量・PAT 検証・PAT 更新フォーム・PAT 有効期限表示・GitHub API レート制限表示 |
| 通知設定 | 複数 Webhook 設定・ペイロードテンプレート編集・テスト送信・定期サマリー設定・送信履歴・メール通知設定（SMTP連携）・Webhook 署名シークレット設定 |
| 設定 | 保持行数・件数・タイムアウト・ログレベル変更・ログ保持期間（日数）・手動クリーンアップ・設定変更履歴・キュー最大サイズ設定・スナップショット保持世代数設定 |
| アクセスログ | ログイン履歴（日時・成否） |
| 統計 | 成功率・平均/最大ビルド時間・時系列グラフ |
| リポジトリ情報 | 監視設定確認・ポーリング間隔変更フォーム・ポーリング一時停止/再開・許可時間帯設定・強制再ビルド間隔設定・Webhook 受信 Secret 設定・ブランチターゲット設定（`.branch_config` 編集） |
| セッション管理 | セッション一覧・全セッション強制無効化 |
| システム診断 | PAT・GitHub API・ファイル・systemd・Webhook・出力整合性 一括診断・アラートバッジ表示・メンテナンスバナー表示 |
| ビルド比較 | ビルド履歴から 2 件を選択してログを並列差分表示 |
| API トークン管理 | 読み取り専用トークンの発行・一覧・失効 |
| スナップショット | ビルド成果物の世代一覧・ダウンロード・削除・ロールバック |
| メンテナンス | メンテナンスモードの有効化・解除・状態表示 |
| アクセス制御 | 許可 IP / CIDR 一覧・追加・削除 |
| フック | Pre/Post ビルドフック設定・実行ログ確認 |
| 運用ノート | Markdown 記述の運用メモ閲覧・編集 |

---

## 5.2 拡張ポイント・将来計画

Adlaire CI の実装済み項目、実装中・検証未完了項目、仕様化済み・未実装項目、改訂予定項目、将来計画項目を統合管理する。
仕様化する際は [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務の対応表、該当 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務への追記を先行させる。

将来計画は、実装対象ではない。将来計画内の「検討」「予定」「候補」「推奨」は、実装可能な仕様を意味しない。将来計画を実装対象にする場合は、先に対象項目を `改訂予定` へ昇格し、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務の対応表と該当 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務に実装可能な詳細仕様を追加したうえで `仕様化済み` とする。

**担当領域：** `CI ランナー` / `管理ツール・API` / `MCP サーバー` / `ビルドスクリプト`

### 5.2.1 統合ロードマップ参照入口

拡張ポイント・将来計画の項目は、単一の統合ロードマップ表で管理する。状態分類と実装可否の定義は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §2 を参照する。実装着手可否の判定は [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0a と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0d を参照する。

---

### 5.2.2 統合ロードマップ表

[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.2 の統合ロードマップ表は、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2 の全項目を状態別に統合した唯一の一覧である。項目を追加、削除、昇格、実装済みにする場合は、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.2 の統合ロードマップ表の `状態`、`実装可否`、`次アクション` を同時に更新する。

MCP サーバー領域の行は、現時点ではすべて将来構想例であり、実装契約、API 契約、状態ファイル契約、起動手順、検証条件の正本ではない。`components/mcp.go`、MCP tools、MCP resources、MCP prompts、HTTP SSE transport、MCP audit / stats / config CRUD は、MCP 専用詳細仕様を新設し、`改訂予定` を経て `仕様化済み・未実装` へ昇格するまで実装してはならない。

[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.2 の統合ロードマップ表の `概要` は、状態・計画責務としての状態要約である。関数、処理順序、endpoint、状態 schema、fixture、UI DOM、SDK method の本文は owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を参照する。実装済み行は、実装済み状態と検証済み証跡の要約だけを示す。

[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.2 の統合ロードマップ表の `確認先 / 次アクション` は、状態に応じた確認先または次作業だけを示す。実装済み行では検証済み証跡の要約を示し、fixture の入力、expected、fake、assertion 本文は [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務を参照する。未実装行では実装前に読む責務文書を示し、詳細本文の正本ではない。

| 状態 | 実装可否 | 担当領域 | 機能 | 概要 | 確認先 / 次アクション |
|------|----------|----------|------|------|--------------|
| 実装済み | 検証済み | CI ランナー | Phase 2 判定パス | [`docs/details/runner.md`](details/runner.md) 詳細本文責務の Phase 2 範囲が実装済み。 | Go test で Phase 2 fixture、hardening、判定パスを検証済み。 |
| 改訂予定 | 実装不可 | 全領域 | （なし） | 現時点で、将来計画から格上げ済みの仕様作成中項目はない。 | 格上げ時に元状態、格上げ日、整理順序（実装単位ではない）、詳細仕様作成先を概要へ記録する。 |
| 実装済み | 完了済み | CI ランナー | ビルドタイムアウト | [`docs/details/runner.md`](details/runner.md) 詳細本文責務のビルドタイムアウト契約が実装済み。API 経由の `build_timeout_seconds` 動的変更は管理 API 実装対象として残す。 | Go test と runner 回帰検証で pipeline 起動経路を検証済み。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ポーリング間隔の動的変更 | [`docs/details/api.md`](details/api.md) 詳細本文責務と [`docs/details/setup.md`](details/setup.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.3、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/setup.md`](details/setup.md) §26、[`docs/details/api.md`](details/api.md) §27.11 に従って実装する。 |
| 実装済み | 完了済み | CI ランナー | ビルドログのファイル保存 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務のビルドログ保存契約が実装済み。 | Go test で build log 生成と pipeline stdout/stderr 記録を検証済み。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | GitHub Webhook 受信 | [`docs/details/api.md`](details/api.md) 詳細本文責務と [`docs/details/runner.md`](details/runner.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.3、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/api.md`](details/api.md) §22-W、[`docs/details/api.md`](details/api.md) §27.12 に従って実装する。 |
| 実装済み | 完了済み | CI ランナー | ネットワーク断時の再試行 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務の GitHub API retry 契約が実装済み。 | Go test で fake GitHub 一時失敗からの retry 成功を検証済み。 |
| 実装済み | 完了済み | CI ランナー | GitHub API レート制限自動待機 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務の rate limit 待機契約が実装済み。 | Go test 対象の retry 経路と同じ GitHub API retry 実装で検証済み。 |
| 実装済み | 完了済み | CI ランナー | 転送後リモート整合性検証 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務の転送後整合性検証契約が実装済み。 | Go test で転送失敗時 pending 化と pending 再試行成功を検証済み。 |
| 実装済み | 完了済み | CI ランナー | マルチブランチビルド | [`docs/details/runner.md`](details/runner.md) 詳細本文責務の branch target 契約が実装済み。 | Go test で branch target 設定に基づく処理経路を検証済み。 |
| 実装済み | 完了済み | CI ランナー | ビルドログ世代管理 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務のログ世代管理契約が実装済み。 | `go test ./...` で runner 回帰検証済み。 |
| 実装済み | 完了済み | CI ランナー | ビルド出力の外部転送 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務の外部転送契約が実装済み。 | Go test で SSH command 差し替えによる pending / retry 経路を検証済み。 |
| 実装済み | 完了済み | CI ランナー | ビルドクールダウン | [`docs/details/runner.md`](details/runner.md) 詳細本文責務の cooldown 契約が実装済み。 | Go test で cooldown 中の build skip を検証済み。 |
| 実装済み | 完了済み | CI ランナー | ビルド前の事前チェック | [`docs/details/runner.md`](details/runner.md) 詳細本文責務のビルド前事前チェック契約が実装済み。 | Go test で事前チェック失敗経路を検証済み。 |
| 実装済み | 完了済み | CI ランナー | 定期強制ビルド | [`docs/details/runner.md`](details/runner.md) 詳細本文責務の force interval 契約が実装済み。API 経由の動的変更は管理 API 実装対象として残す。 | Go test で同一 SHA かつ force interval 超過時の build 実行を検証済み。 |
| 実装済み | 完了済み | CI ランナー | ビルド中重複スキップ | [`docs/details/runner.md`](details/runner.md) 詳細本文責務の build lock 契約が実装済み。 | Go test で lock conflict 時の build skip を検証済み。 |
| 実装済み | 完了済み | CI ランナー | GitHub PAT 有効期限の事前警告 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務の PAT 期限警告契約が実装済み。 | Go test で GitHub API response header 処理を検証済み。 |
| 実装済み | 完了済み | CI ランナー | コミット情報のビルドログ記録 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務の commit info 記録契約が実装済み。 | Go test で fake commits API からの commit info 記録を検証済み。 |
| 実装済み | 完了済み | CI ランナー | GitHub API 連続失敗によるサーキットブレーカー | [`docs/details/runner.md`](details/runner.md) 詳細本文責務の circuit breaker 契約が実装済み。 | Go test で circuit open 時の polling skip を検証済み。 |
| 実装済み | 完了済み | CI ランナー | 出力サイトサイズ警告閾値 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/builder.md`](details/builder.md) 詳細本文責務の出力サイズ警告契約が実装済み。 | Go test で閾値超過時の `size_warn=true` と WARN 記録を検証済み。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | Webhook イベントログ | [`docs/details/api.md`](details/api.md) 詳細本文責務と [`docs/details/runner.md`](details/runner.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.3、[`docs/details/runner.md`](details/runner.md) §11、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/api.md`](details/api.md) §22-W、[`docs/details/api.md`](details/api.md) §27.13 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド所要時間の記録と統計 API | [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/api.md`](details/api.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.4、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.14 に従って実装する。 |
| 実装済み | 完了済み | CI ランナー | ビルドアーティファクト世代管理 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/archive.md`](details/archive.md) 詳細本文責務の snapshot 世代管理契約が実装済み。`POST /api/history/{id}/rollback` による再転送は API 実装対象として残す。 | Go test で build 成功時の `.snapshots/{id}/site` 作成を検証済み。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | ビルドアーティファクト管理 | [`docs/details/archive.md`](details/archive.md)、[`docs/details/api.md`](details/api.md)、[`docs/details/sdk.md`](details/sdk.md)、[`docs/details/ui.md`](details/ui.md) の詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.4、[`docs/details/runner.md`](details/runner.md) §14b、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/archive.md`](details/archive.md) §27.15 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | ヘルスチェックエンドポイント | [`docs/details/api.md`](details/api.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.3、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/api.md`](details/api.md) §27.16 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | Webhook イベント一覧取得 API | [`docs/details/api.md`](details/api.md)、[`docs/details/sdk.md`](details/sdk.md)、[`docs/details/ui.md`](details/ui.md) の詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.3、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/api.md`](details/api.md) §27.13 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | ビルドログ重大度フィルター | [`docs/details/api.md`](details/api.md)、[`docs/details/sdk.md`](details/sdk.md)、[`docs/details/ui.md`](details/ui.md) の詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.3、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/api.md`](details/api.md) §27.17 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | 変換レポート出力 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務の変換レポート契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | ビルドスクリプト | シンタックスハイライト | [`docs/details/builder.md`](details/builder.md) 詳細本文責務のシンタックスハイライト契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | ビルドスクリプト | 本文内全文検索 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務の全文検索契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | ビルドスクリプト | アンカーリンク自動検証 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務のアンカーリンク検証契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | ビルドスクリプト | コードブロックの折りたたみ | [`docs/details/builder.md`](details/builder.md) 詳細本文責務のコードブロック折りたたみ契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | ビルドスクリプト | 印刷スタイル（`@media print`） | [`docs/details/builder.md`](details/builder.md) 詳細本文責務の印刷スタイル契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | ビルドスクリプト | 静的 Web サイト出力 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務の静的 Web サイト出力契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | ビルドスクリプト | テーマコンポーネント | [`docs/details/builder.md`](details/builder.md) 詳細本文責務のテーマコンポーネント契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | ビルドスクリプト | 外部リンクの自動処理 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務の外部リンク処理契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | ビルドスクリプト | 読み取り進捗バー | [`docs/details/builder.md`](details/builder.md) 詳細本文責務の読み取り進捗バー契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | ビルドスクリプト | コードブロックのコピーボタン | [`docs/details/builder.md`](details/builder.md) 詳細本文責務のコードブロックコピーボタン契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | ビルドスクリプト | 見出しアンカーリンクコピー | [`docs/details/builder.md`](details/builder.md) 詳細本文責務の見出しアンカーリンクコピー契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | ビルドスクリプト | TOC 開閉状態の永続化 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務の TOC 開閉状態永続化契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | ビルドスクリプト | 見出しスラグ重複解決 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務の見出しスラグ重複解決契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | ビルドスクリプト | 前後章ナビゲーションボタン | [`docs/details/builder.md`](details/builder.md) 詳細本文責務の前後章ナビゲーション契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | ビルドスクリプト | 内部リンク整合性チェック | [`docs/details/builder.md`](details/builder.md) 詳細本文責務の内部リンク整合性チェック契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | ビルドスクリプト | 見出し階層スキップ警告 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務の見出し階層スキップ警告契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | ビルドスクリプト | 読了時間推計と表示 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務の読了時間推計契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | CI ランナー | Webhook 通知失敗リトライキュー | [`docs/details/runner.md`](details/runner.md) 詳細本文責務の Webhook 通知失敗リトライキュー契約が実装済み。 | Go test で通知再送成功経路を検証済み。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ブランチ設定の動的変更 API | [`docs/details/api.md`](details/api.md) 詳細本文責務と [`docs/details/runner.md`](details/runner.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.3、[`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/api.md`](details/api.md) §27.18 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | 週次ビルドサマリー Webhook | [`docs/details/runner.md`](details/runner.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.3、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §16、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.19 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | 設定変更の詳細 diff 記録 | [`docs/details/api.md`](details/api.md) 詳細本文責務と [`docs/details/statefile.md`](details/statefile.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.3、[`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/api.md`](details/api.md) §27.20 に従って実装する。 |
| 実装済み | 完了済み | ビルドスクリプト | テーブルのソート機能 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務のテーブルソート契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 実装済み | 完了済み | ビルドスクリプト | キーボードショートカット | [`docs/details/builder.md`](details/builder.md) 詳細本文責務のキーボードショートカット契約が実装済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1 の確認済み項目。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | 複数ファイル監視 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/runner.md`](details/runner.md) §27.21 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | GitHub Commit Status API | [`docs/details/commitstatus.md`](details/commitstatus.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/commitstatus.md`](details/commitstatus.md) §27.1 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルドパイプライン YAML 定義 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.22 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ドライラン実行モード | [`docs/details/runner.md`](details/runner.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/runner.md`](details/runner.md) §27.2 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルドログのアーカイブ圧縮 | [`docs/details/archive.md`](details/archive.md) 詳細本文責務と [`docs/details/runner.md`](details/runner.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.4、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/archive.md`](details/archive.md) §27.7 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ローカルファイル監視モード | [`docs/details/runner.md`](details/runner.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §27.23 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | タグ付きコミットのみビルド | [`docs/details/runner.md`](details/runner.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.24 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルドキャッシュ | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/runner.md`](details/runner.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §5、[`docs/details/builder.md`](details/builder.md) §8、[`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/builder.md`](details/builder.md) §27.25 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド通知連携 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/api.md`](details/api.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/runner.md`](details/runner.md) §16、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.32 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド時間トレンド記録 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/api.md`](details/api.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.33 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド失敗時の自動リトライ | [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/statefile.md`](details/statefile.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/runner.md`](details/runner.md) §27.3 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | 並列マルチターゲットビルド | [`docs/details/runner.md`](details/runner.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §14a、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/runner.md`](details/runner.md) §27.26 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド前後フック | [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/api.md`](details/api.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.27 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | 依存ファイルトラッキング | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/runner.md`](details/runner.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §4.3、[`docs/details/builder.md`](details/builder.md) §5、[`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/builder.md`](details/builder.md) §27.28 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | リモートビルド対応 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §14a、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/runner.md`](details/runner.md) §27.29 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルドステータスファイル出力 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/statefile.md`](details/statefile.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.8 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド承認フロー | [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/api.md`](details/api.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.3、[`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/runner.md`](details/runner.md) §16、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.30 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ブランチ別環境変数 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/api.md`](details/api.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.31 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド依存チェーン | [`docs/details/runner.md`](details/runner.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.34 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド優先度キュー | [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/api.md`](details/api.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.35 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | 失敗原因の自動分類 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/api.md`](details/api.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.36 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド実行環境の記録 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/runner.md`](details/runner.md) §27.37 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルドトリガー種別の記録 | [`docs/details/runner.md`](details/runner.md)、[`docs/details/statefile.md`](details/statefile.md)、[`docs/details/api.md`](details/api.md)、[`docs/details/sdk.md`](details/sdk.md)、[`docs/details/ui.md`](details/ui.md) の詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/runner.md`](details/runner.md) §27.9、[`docs/details/runner.md`](details/runner.md) §27.23、[`docs/details/runner.md`](details/runner.md) §27.30 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | ビルド所要時間の異常検知 | [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/api.md`](details/api.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/runner.md`](details/runner.md) §15、[`docs/details/runner.md`](details/runner.md) §16、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/runner.md`](details/runner.md) §27.38 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | CI ランナー | 設定ファイル起動時整合性チェック | [`docs/details/runner.md`](details/runner.md) 詳細本文責務と [`docs/details/statefile.md`](details/statefile.md) 詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.2、[`docs/details/runner.md`](details/runner.md) §11、[`docs/details/runner.md`](details/runner.md) §12、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/runner.md`](details/runner.md) §27.10 に従って実装する。 |
| 将来計画 | 実装不可 | 管理ツール・API | マルチユーザー対応 | 将来計画候補。詳細本文、ユーザー model、認証境界、監査境界は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 通知先の拡張 | 将来計画候補。詳細本文、対象通知先、API / UI 契約は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | データストア切り替え | 将来計画候補。詳細本文、採用データストア、移行契約は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 外部認証連携 | 将来計画候補。詳細本文、認証方式、連携境界は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | TOTP 二要素認証 | [`docs/details/security.md`](details/security.md)、[`docs/details/api.md`](details/api.md)、[`docs/details/sdk.md`](details/sdk.md)、[`docs/details/ui.md`](details/ui.md) の詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.4、[`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/security.md`](details/security.md) §27.46 に従って実装する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 統計データの CSV エクスポート | 将来計画候補。詳細本文、出力形式、endpoint 契約は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | キュー内個別エントリのキャンセル | 将来計画候補。詳細本文、endpoint 契約、対象 queue、失敗時副作用は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | 設定バリデーション API | [`docs/details/api.md`](details/api.md)、[`docs/details/statefile.md`](details/statefile.md)、[`docs/details/sdk.md`](details/sdk.md)、[`docs/details/ui.md`](details/ui.md) の詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.3、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/api.md`](details/api.md) §27.5 に従って実装する。 |
| 将来計画 | 実装不可 | 管理ツール・API | Prometheus メトリクスエンドポイント | 将来計画候補。詳細本文、endpoint 契約、出力形式は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | CLI 管理クライアント | 将来計画候補。詳細本文、binary 名、command、API 呼び出し契約は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 設定の自動スナップショット | 将来計画候補。詳細本文、保存契約、副作用境界は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ステータスバッジ生成 | 将来計画候補。詳細本文、endpoint 契約、出力形式は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ビルド履歴の自動削除設定 | 将来計画候補。詳細本文、保持期間 key、削除対象、副作用境界は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | セッションタイムアウト変更設定 | [`docs/details/security.md`](details/security.md)、[`docs/details/api.md`](details/api.md)、[`docs/details/sdk.md`](details/sdk.md)、[`docs/details/ui.md`](details/ui.md) の詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.4、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/security.md`](details/security.md) §27.45 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | ビルドトリガー専用 API スコープ | [`docs/details/security.md`](details/security.md)、[`docs/details/api.md`](details/api.md)、[`docs/details/sdk.md`](details/sdk.md)、[`docs/details/ui.md`](details/ui.md) の詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.4、[`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/security.md`](details/security.md) §27.42 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | 監査ログ | [`docs/details/security.md`](details/security.md)、[`docs/details/statefile.md`](details/statefile.md)、[`docs/details/api.md`](details/api.md)、[`docs/details/sdk.md`](details/sdk.md)、[`docs/details/ui.md`](details/ui.md) の詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.4、[`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/security.md`](details/security.md) §27.44 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | API レート制限 | [`docs/details/security.md`](details/security.md)、[`docs/details/statefile.md`](details/statefile.md)、[`docs/details/api.md`](details/api.md)、[`docs/details/sdk.md`](details/sdk.md)、[`docs/details/ui.md`](details/ui.md) の詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.4、[`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/security.md`](details/security.md) §27.47 に従って実装する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ロールベースアクセス制御 | 将来計画候補。詳細本文、role 種別、認可 matrix、監査境界は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 成果物ダウンロード API | 将来計画候補。詳細本文、endpoint 契約、archive 形式、認可境界は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 設定スナップショット差分表示 | 将来計画候補。詳細本文、diff 形式、snapshot 対象、表示契約は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 複数プロジェクト管理 | 将来計画候補。詳細本文、project model、切替境界、状態分離は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | API キー管理 | [`docs/details/security.md`](details/security.md)、[`docs/details/statefile.md`](details/statefile.md)、[`docs/details/api.md`](details/api.md)、[`docs/details/sdk.md`](details/sdk.md)、[`docs/details/ui.md`](details/ui.md) の詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.4、[`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/security.md`](details/security.md) §27.43 に従って実装する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ビルドログのリアルタイム配信 | 将来計画候補。詳細本文、transport、endpoint 契約は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 設定のインポート／エクスポート | 将来計画候補。詳細本文、対象設定、形式、検証条件、失敗時処理は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ビルド統計ダッシュボード | 将来計画候補。詳細本文、集計項目、表示契約、更新条件は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ユーザー管理 API | 将来計画候補。詳細本文、user model、endpoint 契約、認可境界は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | IP アドレス制限 | 将来計画候補。詳細本文、許可範囲形式、評価順、拒否時 response は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | API バージョニング | 将来計画候補。詳細本文、version path、互換方針、routing 境界は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | Webhook 署名検証 | 将来計画候補。詳細本文、署名検証方式、拒否条件は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | API ドキュメント自動生成 | 将来計画候補。詳細本文、schema 形式、生成範囲、公開方法は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 通知チャンネル管理 | 将来計画候補。詳細本文、対象通知先、API / UI 契約は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ビルドキューの手動並び替え | 将来計画候補。詳細本文、queue 操作、並び替え制約、競合処理は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 設定テンプレート | 将来計画候補。詳細本文、template 形式、保存先、適用条件は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 仕様化済み・未実装 | 実装可 | 管理ツール・API | API アクセスログ | [`docs/details/api.md`](details/api.md)、[`docs/details/statefile.md`](details/statefile.md)、[`docs/details/sdk.md`](details/sdk.md)、[`docs/details/ui.md`](details/ui.md) の詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.3、[`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/statefile.md`](details/statefile.md) §22.0c、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/sdk.md`](details/sdk.md) §23、[`docs/details/ui.md`](details/ui.md) §24、[`docs/details/api.md`](details/api.md) §27.6 に従って実装する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 管理者向けイベントフィード | 将来計画候補。詳細本文、event 種別、配信方式、表示契約は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | ビルドキュー可視化 | 将来計画候補。詳細本文、queue 状態、表示形式、更新条件は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | メンテナンスモード | 将来計画候補。詳細本文、切替条件、停止範囲、復帰条件は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | 読み取り専用共有リンク | 将来計画候補。詳細本文、link 形式、有効期限、認可境界は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | アラート閾値設定 | 将来計画候補。詳細本文、閾値項目、評価周期、通知先は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | バックアップ／リストア | 将来計画候補。詳細本文、対象データ、archive 形式、復元手順は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | API レスポンスキャッシュ制御 | 将来計画候補。詳細本文、対象 endpoint、TTL、無効化条件は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | スナップショット間サイト差分 API | 将来計画候補。詳細本文、endpoint 契約、差分形式は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 将来計画 | 実装不可 | 管理ツール・API | Webhook 送信履歴の手動再送 API | 将来計画候補。詳細本文、endpoint 契約、再送境界は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、詳細仕様を追加する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | 差分ビルド | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.1、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | 複数出力形式 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.2、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | Markdown 拡張記法サポート | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.3、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | コードブロック行番号表示 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.4、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | 見出しの自動採番 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.5、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | セクション折りたたみ | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.6、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | TOC 深さ制御 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.7、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | 最終更新日の自動埋め込み | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.8、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | diff ハイライト | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.9、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | 画像の遅延読み込み | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.10、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | カスタムメタタグ注入 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.11、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | ライトモード固定 | [`docs/DESIGN.md`](DESIGN.md) デザイン責務、[`docs/details/builder.md`](details/builder.md) 詳細本文責務、[`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DESIGN.md`](DESIGN.md) デザイン責務、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.12、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | コードブロックのファイル名表示 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.13、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | テンプレート変数展開 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.14、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | HTML ミニファイ | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.15、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | TOC ハイライト追従 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.16、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | Mermaid ダイアグラム描画 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.17、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | 脚注サポート | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.18、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | インライン数式レンダリング | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.19、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | ページ内ナビゲーション履歴 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.20、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | 読み上げ対応（アクセシビリティ） | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.21、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | 画像ライトボックス | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.22、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | 出力サイトへのビルドメタ埋め込み | [`docs/details/builder.md`](details/builder.md)、[`docs/details/runner.md`](details/runner.md)、[`docs/details/api.md`](details/api.md) の詳細本文責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §2、[`docs/details/builder.md`](details/builder.md) §5、[`docs/details/builder.md`](details/builder.md) §8、[`docs/details/runner.md`](details/runner.md) §13、[`docs/details/api.md`](details/api.md) §22.0e、[`docs/details/builder.md`](details/builder.md) §27.4 に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | 印刷時 QR コード挿入 | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.23、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | 定義リストサポート | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.24、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 仕様化済み・未実装 | 実装可 | ビルドスクリプト | タスクリストサポート | [`docs/details/builder.md`](details/builder.md) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) fixture 証跡責務に実装契約が定義済み。 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1、[`docs/details/builder.md`](details/builder.md) §28.25、[`docs/details/fixture.md`](details/fixture.md) §28-F に従って実装する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP サーバー実装 | 将来計画候補。MCP 専用詳細仕様が新設されるまで、実装 path、transport、認証、接続設定は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP ツール・リソース公開 | 将来計画候補。MCP 専用詳細仕様が新設されるまで、tool 名、resource URI、scope、戻り値は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | AI 支援ビルドエラー分析 | 将来計画候補。MCP 専用詳細仕様が新設されるまで、分析開始条件、参照 tool、保存先、責務境界は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP Prompts 定義 | 将来計画候補。MCP 専用詳細仕様が新設されるまで、prompt 名、入力、参照 tool、出力は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP Sampling によるビルドログ自動分析 | 将来計画候補。MCP 専用詳細仕様が新設されるまで、sampling 契約、保存先、runner 連携、副作用境界は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP Notifications（イベントプッシュ） | 将来計画候補。MCP 専用詳細仕様が新設されるまで、イベント種別、queue、通知 payload、再送条件は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP HTTP SSE transport 対応 | 将来計画候補。MCP 専用詳細仕様が新設されるまで、transport、port、認証、同時接続条件は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP ツールスコープ細分化 | 将来計画候補。MCP 専用詳細仕様が新設されるまで、scope 種別、検証条件、token 保存先は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP ツール呼び出し監査ログ | 将来計画候補。MCP 専用詳細仕様が新設されるまで、監査項目、保存先、参照 API、mask 条件は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP リソース購読（Resource Subscriptions） | 将来計画候補。MCP 専用詳細仕様が新設されるまで、購読対象、更新判定、通知 payload は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP クライアント情報ログ | 将来計画候補。MCP 専用詳細仕様が新設されるまで、記録項目、保存先、保持期間、表示先は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP ツール実行統計 | 将来計画候補。MCP 専用詳細仕様が新設されるまで、統計項目、集計単位、保存先、参照 API は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP ツール実行タイムアウト設定 | 将来計画候補。MCP 専用詳細仕様が新設されるまで、timeout 設定 key、既定値、超過時 error は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP 設定 CRUD ツール | 将来計画候補。MCP 専用詳細仕様が新設されるまで、設定対象、tool 名、認可、検証条件は未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |
| 将来計画 | 実装不可 | MCP サーバー | MCP Elicitation による副作用操作の確認 | 将来計画候補。MCP 専用詳細仕様が新設されるまで、副作用 tool、確認 payload、承認結果の扱いは未定義。 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.2.3 の手順で `改訂予定` へ昇格し、MCP 専用詳細仕様を新設する。 |

### 5.2.3 状態変更・昇格手順

統合ロードマップ表の `将来計画` の項目を実装対象にする場合は、以下の順で進める。

1. 対象項目の `状態` を `改訂予定` に変更し、元状態、格上げ日、担当領域、整理順序（実装単位ではない）、ステータスを同じ行の概要または次アクションへ記録する。
2. [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務の対象節と、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §3 のコンポーネント状態分類、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §5.1 の機能一覧を改訂する。
3. [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i.1〜§0i.4 の詳細節対応表に、対象機能、対象コンポーネント、詳細仕様節、受け入れ条件を追加する。
4. 該当 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務に、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0h の機能仕様テンプレートを満たす目的、責務、入出力、状態、処理順序、異常系、セキュリティ、検証条件を追加する。MCP サーバー領域を昇格する場合は、MCP 専用詳細仕様を新設し、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務の owner component 参照表へ追加する。
5. [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務と [`AGENTS.md`](../AGENTS.md) の更新要否を確認する。
6. 仕様凍結条件を満たした後、対象項目の `状態` を `仕様化済み・未実装`、`実装可否` を `実装可` に変更する。
7. 実装と検証が完了した後、対象項目の `状態` を `実装済み`、`実装可否` を `完了済み` に変更する。

[`docs/ROADMAP.md`](ROADMAP.md) §5.2.3 の昇格手順を完了していない項目の実装着手可否は、[`docs/SPEC.md`](SPEC.md) ポリシー責務 §0a と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0d を参照する。

## 6. 追加仕様化機能実装参照

追加仕様化機能実装参照は、追加仕様化機能の状態、実装可否、owner、主本文、collaborator を確認するための参照索引である。各機能の主本文は、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を参照する。参照先の特定は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i と [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6 で行う。

### 6.1 追加仕様化機能 共通実装契約

[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6.2 に列挙する追加仕様化機能の主本文は、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を参照する。追加仕様化機能共通実装契約は、実装時に共通して確認する参照順、越境確認、実装検証証跡の入口だけを示す。

**追加仕様化機能 横断共通参照先：**

[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6 で API endpoint、SDK method、UI 操作、状態ファイル副作用、fixture、実装検証証跡に触れる場合、API 契約は [`docs/details/api.md`](details/api.md)、SDK 契約は [`docs/details/sdk.md`](details/sdk.md)、UI 契約は [`docs/details/ui.md`](details/ui.md)、状態ファイル契約は [`docs/details/statefile.md`](details/statefile.md)、fixture / 実装検証証跡は [`docs/details/fixture.md`](details/fixture.md) を共通参照先とする。状態・計画責務 §6 では、追加仕様化機能の owner、主本文、collaborator、横断確認観点だけを本文として扱う。

| 確認 | 固定内容 |
|------|----------|
| 実装対象判定 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務で状態分類と実装可否を確認する。将来計画、実装不可、未仕様化、MCP 専用詳細仕様がない状態の MCP 機能の扱いは、[`docs/SPEC.md`](SPEC.md) ポリシー責務 §0a と [`docs/SPEC.md`](SPEC.md) ポリシー責務 §0d を参照する。 |
| owner 確定 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0i で owner component を 1 件に確定し、主本文は owner component 別の [`docs/details/*.md`](details/) 詳細本文責務で確認する。 |
| collaborator 確認 | collaborator がある場合は、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6.2 の追加仕様化機能参照索引に列挙された collaborator component 別の [`docs/details/*.md`](details/) 詳細本文責務を schema、呼び出し境界、表示、security、setup、fixture、検証観点として読む。 |
| 補完確認 | endpoint、状態ファイル、設定 key、UI 操作、SDK method、外部依存が、個別節または owner component 別の [`docs/details/*.md`](details/) 詳細本文責務に存在することを確認する。存在しない項目の扱いは、[`docs/SPEC.md`](SPEC.md) ポリシー責務 §0、[`docs/SPEC.md`](SPEC.md) ポリシー責務 §4、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務、関連 collaborator component 別の [`docs/details/*.md`](details/) 詳細本文責務、fixture catalog、必要な対応表を参照する。 |
| 状態更新 | 状態ファイル更新は [`docs/details/statefile.md`](details/statefile.md) §22.0a、[`docs/details/statefile.md`](details/statefile.md) §22.0c を基準とする。lock、atomic write、JSON Lines、破損時処理は [`docs/details/statefile.md`](details/statefile.md) 詳細本文責務を正本とする。 |
| security | secret mask、token、session、scope、audit、rate limit は [`docs/details/security.md`](details/security.md) を基準とし、平文保存・平文表示を行わない。 |
| fixture / 実装検証証跡 | 追加仕様化機能の fixture manifest、expected/effects、assertion、実装検証証跡、受け入れゲートは [`docs/details/fixture.md`](details/fixture.md) §27-F、builder 拡張 fixture は [`docs/details/fixture.md`](details/fixture.md) §28-F を参照する。 |
| api / sdk / ui 同期 | API endpoint、SDK method、UI 操作が同一機能に関わる場合は、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6 追加仕様化機能 横断共通参照先に従い、名称、引数、response、error、表示、成功後再取得、失敗時固定が食い違わないことを確認する。 |

[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6.2 の追加仕様化機能を実装した変更は、対象節、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務、collaborator component 別の [`docs/details/*.md`](details/) 詳細本文責務、fixture、secret mask、失敗時副作用、実装対象外を [`docs/details/fixture.md`](details/fixture.md) の実装検証証跡に記録する。記録が不足する場合は、実装済みとして扱わない。

追加仕様化機能の責務分離、dry-run 固定契約、fixture 判定条件の詳細は、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務と [`docs/details/fixture.md`](details/fixture.md) §27-F を参照する。builder 拡張 fixture 判定条件は、[`docs/details/builder.md`](details/builder.md) §28 と [`docs/details/fixture.md`](details/fixture.md) §28-F を参照する。[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務は同じ fixture schema、expected/effects、個別機能本文の正本ではない。

### 6.2 追加仕様化機能 参照索引

[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6.2 の追加仕様化機能参照索引は、追加仕様化機能の owner、主本文、collaborator を一覧化するインデックスである。個別機能本文、状態 schema、endpoint、SDK method、UI DOM、fixture schema、横断処理順は、同索引の「主本文」に記載された owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を正本とする。[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務は、同索引に記載された主本文、owner component、collaborator component を置き換えない。

| 節 | 機能 | owner | 主本文 | collaborator |
|----|------|-------|--------|--------------|
| [`docs/details/commitstatus.md`](details/commitstatus.md) §27.1 | GitHub Commit Status API | `commitstatus` | [`docs/details/commitstatus.md`](details/commitstatus.md) §27.1 | `runner`、`statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.2 | ドライラン実行モード | `runner` | [`docs/details/runner.md`](details/runner.md) §27.2 | `statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.3 | ビルド失敗時の自動リトライ | `runner` | [`docs/details/runner.md`](details/runner.md) §27.3 | `statefile` |
| [`docs/details/builder.md`](details/builder.md) §27.4 | 出力サイトへのビルドメタ埋め込み | `builder` | [`docs/details/builder.md`](details/builder.md) §27.4 | `runner`、`api`、`statefile` |
| [`docs/details/api.md`](details/api.md) §27.5 | 設定バリデーション API | `api` | [`docs/details/api.md`](details/api.md) §27.5 | `sdk`、`ui`、`statefile` |
| [`docs/details/api.md`](details/api.md) §27.6 | API アクセスログ | `api` | [`docs/details/api.md`](details/api.md) §27.6 | `sdk`、`ui`、`statefile` |
| [`docs/details/archive.md`](details/archive.md) §27.7 | ビルドログのアーカイブ圧縮 | `archive` | [`docs/details/archive.md`](details/archive.md) §27.7 | `runner`、`api`、`statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.8 | ビルドステータスファイル出力 | `runner` | [`docs/details/runner.md`](details/runner.md) §27.8 | `api`、`statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.9 | ビルドトリガー種別の記録 | `runner` | [`docs/details/runner.md`](details/runner.md) §27.9 | `api`、`sdk`、`ui`、`statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.10 | 設定ファイル起動時整合性チェック | `runner` | [`docs/details/runner.md`](details/runner.md) §27.10 | `statefile` |
| [`docs/details/api.md`](details/api.md) §27.11 | ポーリング間隔の動的変更 | `api` | [`docs/details/api.md`](details/api.md) §27.11 | `runner`、`statefile` |
| [`docs/details/api.md`](details/api.md) §27.12 | GitHub Webhook 受信 | `api` | [`docs/details/api.md`](details/api.md) §27.12 | `runner`、`statefile` |
| [`docs/details/api.md`](details/api.md) §27.13 | Webhook イベントログ / 一覧取得 API | `api` | [`docs/details/api.md`](details/api.md) §27.13 | `sdk`、`ui`、`statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.14 | ビルド所要時間の記録と統計 API | `runner` | [`docs/details/runner.md`](details/runner.md) §27.14 | `api`、`statefile`、`archive` |
| [`docs/details/archive.md`](details/archive.md) §27.15 | ビルドアーティファクト管理 | `archive` | [`docs/details/archive.md`](details/archive.md) §27.15 | `api`、`sdk`、`ui`、`runner`、`statefile` |
| [`docs/details/api.md`](details/api.md) §27.16 | ヘルスチェックエンドポイント | `api` | [`docs/details/api.md`](details/api.md) §27.16 | `statefile` |
| [`docs/details/api.md`](details/api.md) §27.17 | ビルドログ重大度フィルター | `api` | [`docs/details/api.md`](details/api.md) §27.17 | `sdk`、`ui`、`archive`、`statefile` |
| [`docs/details/api.md`](details/api.md) §27.18 | ブランチ設定の動的変更 API | `api` | [`docs/details/api.md`](details/api.md) §27.18 | `runner`、`statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.19 | 週次ビルドサマリー Webhook | `runner` | [`docs/details/runner.md`](details/runner.md) §27.19 | `api`、`statefile` |
| [`docs/details/api.md`](details/api.md) §27.20 | 設定変更の詳細 diff 記録 | `api` | [`docs/details/api.md`](details/api.md) §27.20 | `statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.21 | 複数ファイル監視 | `runner` | [`docs/details/runner.md`](details/runner.md) §27.21 | `builder`、`api`、`statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.22 | ビルドパイプライン YAML 定義 | `runner` | [`docs/details/runner.md`](details/runner.md) §27.22 | `api`、`statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.23 | ローカルファイル監視モード | `runner` | [`docs/details/runner.md`](details/runner.md) §27.23 | `statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.24 | タグ付きコミットのみビルド | `runner` | [`docs/details/runner.md`](details/runner.md) §27.24 | `api`、`statefile` |
| [`docs/details/builder.md`](details/builder.md) §27.25 | ビルドキャッシュ | `builder` | [`docs/details/builder.md`](details/builder.md) §27.25 | `runner`、`statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.26 | 並列マルチターゲットビルド | `runner` | [`docs/details/runner.md`](details/runner.md) §27.26 | `statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.27 | ビルド前後フック | `runner` | [`docs/details/runner.md`](details/runner.md) §27.27 | `api`、`statefile` |
| [`docs/details/builder.md`](details/builder.md) §27.28 | 依存ファイルトラッキング | `builder` | [`docs/details/builder.md`](details/builder.md) §27.28 | `runner`、`statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.29 | リモートビルド対応 | `runner` | [`docs/details/runner.md`](details/runner.md) §27.29 | `api`、`archive`、`statefile` |
| [`docs/details/api.md`](details/api.md) §27.30 | ビルド承認フロー | `api` | [`docs/details/api.md`](details/api.md) §27.30 | `runner`、`sdk`、`ui`、`statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.31 | ブランチ別環境変数 | `runner` | [`docs/details/runner.md`](details/runner.md) §27.31 | `api`、`statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.32 | ビルド通知連携 | `runner` | [`docs/details/runner.md`](details/runner.md) §27.32 | `api`、`sdk`、`ui`、`statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.33 | ビルド時間トレンド記録 | `runner` | [`docs/details/runner.md`](details/runner.md) §27.33 | `api`、`statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.34 | ビルド依存チェーン | `runner` | [`docs/details/runner.md`](details/runner.md) §27.34 | `api`、`statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.35 | ビルド優先度キュー | `runner` | [`docs/details/runner.md`](details/runner.md) §27.35 | `api`、`statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.36 | 失敗原因の自動分類 | `runner` | [`docs/details/runner.md`](details/runner.md) §27.36 | `api`、`statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.37 | ビルド実行環境の記録 | `runner` | [`docs/details/runner.md`](details/runner.md) §27.37 | `statefile` |
| [`docs/details/runner.md`](details/runner.md) §27.38 | ビルド所要時間の異常検知 | `runner` | [`docs/details/runner.md`](details/runner.md) §27.38 | `api`、`statefile` |
| [`docs/details/security.md`](details/security.md) §27.42 | ビルドトリガー専用 API スコープ | `security` | [`docs/details/security.md`](details/security.md) §27.42 | `api`、`sdk`、`ui`、`statefile` |
| [`docs/details/security.md`](details/security.md) §27.43 | API キー管理 | `security` | [`docs/details/security.md`](details/security.md) §27.43 | `api`、`sdk`、`ui`、`statefile` |
| [`docs/details/security.md`](details/security.md) §27.44 | 監査ログ | `security` | [`docs/details/security.md`](details/security.md) §27.44 | `api`、`statefile` |
| [`docs/details/security.md`](details/security.md) §27.45 | セッションタイムアウト変更設定 | `security` | [`docs/details/security.md`](details/security.md) §27.45 | `api`、`sdk`、`ui`、`statefile` |
| [`docs/details/security.md`](details/security.md) §27.46 | TOTP 二要素認証 | `security` | [`docs/details/security.md`](details/security.md) §27.46 | `api`、`sdk`、`ui`、`statefile` |
| [`docs/details/security.md`](details/security.md) §27.47 | API レート制限 | `security` | [`docs/details/security.md`](details/security.md) §27.47 | `api`、`sdk`、`ui`、`statefile` |
| [`docs/details/builder.md`](details/builder.md) §28.1 | 差分ビルド | `builder` | [`docs/details/builder.md`](details/builder.md) §28.1 | `runner`、`statefile` |
| [`docs/details/builder.md`](details/builder.md) §28.2 | 複数出力形式 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.2 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.3 | Markdown 拡張記法サポート | `builder` | [`docs/details/builder.md`](details/builder.md) §28.3 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.4 | コードブロック行番号表示 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.4 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.5 | 見出しの自動採番 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.5 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.6 | セクション折りたたみ | `builder` | [`docs/details/builder.md`](details/builder.md) §28.6 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.7 | TOC 深さ制御 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.7 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.8 | 最終更新日の自動埋め込み | `builder` | [`docs/details/builder.md`](details/builder.md) §28.8 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.9 | diff ハイライト | `builder` | [`docs/details/builder.md`](details/builder.md) §28.9 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.10 | 画像の遅延読み込み | `builder` | [`docs/details/builder.md`](details/builder.md) §28.10 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.11 | カスタムメタタグ注入 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.11 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.12 | ライトモード固定 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.12 | [`docs/DESIGN.md`](DESIGN.md) デザイン責務 |
| [`docs/details/builder.md`](details/builder.md) §28.13 | コードブロックのファイル名表示 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.13 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.14 | テンプレート変数展開 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.14 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.15 | HTML ミニファイ | `builder` | [`docs/details/builder.md`](details/builder.md) §28.15 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.16 | TOC ハイライト追従 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.16 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.17 | Mermaid ダイアグラム描画 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.17 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.18 | 脚注サポート | `builder` | [`docs/details/builder.md`](details/builder.md) §28.18 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.19 | インライン数式レンダリング | `builder` | [`docs/details/builder.md`](details/builder.md) §28.19 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.20 | ページ内ナビゲーション履歴 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.20 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.21 | 読み上げ対応（アクセシビリティ） | `builder` | [`docs/details/builder.md`](details/builder.md) §28.21 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.22 | 画像ライトボックス | `builder` | [`docs/details/builder.md`](details/builder.md) §28.22 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.23 | 印刷時 QR コード挿入 | `builder` | [`docs/details/builder.md`](details/builder.md) §28.23 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.24 | 定義リストサポート | `builder` | [`docs/details/builder.md`](details/builder.md) §28.24 | なし |
| [`docs/details/builder.md`](details/builder.md) §28.25 | タスクリストサポート | `builder` | [`docs/details/builder.md`](details/builder.md) §28.25 | なし |

補足確認は以下に限定する。[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6.3 の補足確認表は実装本文を追加せず、主本文を読む際の境界確認だけを示す。

| 対象 | 補足確認 |
|------|----------|
| [`docs/details/commitstatus.md`](details/commitstatus.md) §27.1 | Commit Status payload と送信順は `commitstatus`、build 実行と最終結果確定は `runner` を参照する。 |
| [`docs/details/api.md`](details/api.md) §27.11 | systemd timer 反映の導入・検証手順は [`docs/details/setup.md`](details/setup.md) §26 を同時に確認する。 |
| [`docs/details/archive.md`](details/archive.md) §27.15 | snapshot 作成トリガーは [`docs/details/runner.md`](details/runner.md) §14b、artifact 操作は `archive` を参照する。 |
| [`docs/details/api.md`](details/api.md) §27.17、[`docs/details/runner.md`](details/runner.md) §27.21、[`docs/details/api.md`](details/api.md) §27.30、[`docs/details/runner.md`](details/runner.md) §27.32 | API、SDK、UI、状態ファイルが連動するため、該当 collaborator component 別の [`docs/details/*.md`](details/) 詳細本文責務と [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6.3 を同時に確認する。 |
| [`docs/details/security.md`](details/security.md) §27.42〜§27.47 | scope、token、audit、session、TOTP、rate limit は `security` を基準とし、API / SDK / UI は呼び出し境界と表示だけを担当する。 |

### 6.3 横断連動・Runner 拡張機能 実装補足契約

横断連動・Runner 拡張機能実装補足契約は、追加仕様化機能の横断補足契約である。[`docs/details/runner.md`](details/runner.md) §27.21〜§27.38 および api / sdk / ui / statefile の横断連動では、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を主本文とし、横断補足契約は横断確認として参照する。

横断補足契約の責務範囲は、横断確認、同期禁止、横断処理順、成功後再取得、失敗時固定、実装確認時の横断受け入れ観点に限定する。個別機能本文は各 owner / collaborator component 別の [`docs/details/*.md`](details/) 詳細本文責務を参照する。

横断補足契約と owner component 別の [`docs/details/*.md`](details/) 詳細本文責務の内容が矛盾する場合は、個別機能の入出力、状態、処理、異常系、endpoint、SDK、UI、fixture は owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を基準とし、横断処理順、API / SDK / UI / statefile 同期、成功後再取得、失敗時固定だけを横断補足契約で確認する。横断補足契約にだけ存在する endpoint、SDK method、UI 操作、状態ファイル、設定 key、fixture の扱いは、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.2a と該当 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を参照する。

| 確認 | 固定内容 |
|------|----------|
| runner 起点 | [`docs/details/runner.md`](details/runner.md) §27.21〜§27.38 の多くは runner の build 実行、queue、history、log、notification に影響するため、[`docs/details/runner.md`](details/runner.md) の該当 §27 節を先に確認する。 |
| builder 連携 | cache、dependency、output meta、生成物に関わる場合は [`docs/details/builder.md`](details/builder.md) の該当 §27 節を同時に確認する。 |
| API 連携 | 設定保存、queue、approval、history、stats、snapshot、rollback、hook、notify、search に関わる場合は [`docs/details/api.md`](details/api.md) の endpoint / state read-write 契約を同時に確認する。 |
| SDK / UI 連携 | API を管理画面から操作する機能は、SDK method は [`docs/details/sdk.md`](details/sdk.md) §23、DOM / 表示条件は [`docs/details/ui.md`](details/ui.md) §24 を参照して確認する。 |
| statefile | 状態 schema、lock、atomic write、JSON Lines、破損時処理、保存順は [`docs/details/statefile.md`](details/statefile.md) を参照する。 |
| archive | snapshot、artifact、download、delete、rollback、log archive は [`docs/details/archive.md`](details/archive.md) を参照する。 |
| fixture | [`docs/details/runner.md`](details/runner.md) §27.21〜§27.38 の受け入れ fixture、secret mask、effects、実装検証証跡は [`docs/details/fixture.md`](details/fixture.md) §27-F を参照する。 |

**api / sdk / ui / statefile 横断連動契約：**

[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6.4 の横断確認表は、新しい API endpoint、SDK method、UI 操作、状態ファイル副作用を定義する表ではない。同表は、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6 追加仕様化機能 横断共通参照先と該当 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務に従い、同一機能群の横断確認観点だけをそろえる。API、SDK、UI の実装順と仕様外仮実装の扱いは、[`docs/SPEC.md`](SPEC.md) ポリシー責務 §0f を参照する。UI が未実装の Phase では、UI 列は fixture の期待操作として確認し、実装済み扱いは [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §2 と [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §3 を参照する。

| 機能群 | API | SDK | UI 操作 | 状態ファイル副作用 | 成功後再取得 | 失敗時固定 |
|--------|-----|-----|---------|--------------------|--------------|------------|
| login / session | `POST /api/login`, `POST /api/login/totp`, `POST /api/logout`, `GET /api/sessions`, `POST /api/sessions/revoke-all` | `login()`, `loginTotp()`, `logout()`, `getSessions()`, `revokeAllSessions()` | login、TOTP 確認、logout、session 一括失効。 | `.admin_credentials`、`.access_log`、`.audit_log`、memory session。token 本体は永続化しない。 | login 成功後は `getDashboard()` → `getStatus()` → `getQueue()`。session 失効後は `getSessions()`。 | `401` は SDK token / ticket / secret field を破棄し、UI は `panel-login` のみ表示する。 |
| build control | `POST /api/build`, `POST /api/build/force`, `POST /api/build/cancel`, `GET /api/build/stream`, `GET /api/queue`, `DELETE /api/queue` | `triggerBuild()`, `buildForce()`, `cancelBuild()`, `streamBuild()`, `getQueue()`, `clearQueue()` | manual build、force build、cancel、stream 開始/停止、queue clear。 | `.build_state`、`.build_lock`、`.build_logs/{id}.json`、queue entry。force build 時だけ SHA cache 更新。 | `getStatus()` → `getQueue()`、stream end 後は `getLogs()` も実行。 | running は `409`、queue full は `429`、maintenance / circuit は `503` または仕様済み `409`。UI は同じ build request を自動再送しない。 |
| logs / history | `GET /api/logs`, `GET /api/logs/search`, `GET /api/history`, `GET /api/history/{id}/log`, `POST /api/history/{id}/comment`, `POST /api/history/{id}/flag`, `POST /api/history/{id}/tags` | `getLogs()`, `searchLogs()`, `getHistory()`, `getHistoryLog()`, `setHistoryComment()`, `setHistoryFlag()`, `setHistoryTags()` | log 表示、検索、history 表示、comment / flag / tag 保存。 | read-only GET は副作用なし。comment / flag / tag は `.build_logs/{id}.json` と、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務で履歴・設定ログ更新が定義された場合に限り `.build_history`、`.config_log`。 | 変更系は `getHistory()`、comment は `getHistoryComment(id)` も実行。 | `404` は選択解除または not found 表示。`422 details` は field error。壊れた JSON Lines は response に含めない。 |
| config / repo / branch / schedule | `GET/POST /api/config`, `POST /api/config/validate`, `GET/POST /api/repo-config`, `GET/POST /api/branch-config`, `GET /api/schedule`, `POST /api/schedule/*` | `getConfig()`, `setConfig()`, `validateConfig()`, `setRepoConfig()`, `getBranchConfig()`, `setBranchConfig()`, `getSchedule()`, schedule 系 method | config 保存、validate、repo 保存、branch 保存、schedule 変更。 | `.server_config`、`.repo_config`、`.branch_config`、`.config_log`。validate は保存なし。systemd 反映失敗時も保存済み値は戻さない。 | 保存系は対象 GET → `getConfigLog()`。validate は再取得なし。 | `422` は書込前停止。systemd 失敗 `500` は保存済み値を再取得して表示する。no-op は状態ファイルと log を変更しない。 |
| notify / SMTP / webhook | `GET/POST /api/notify-config`, `POST /api/notify-test`, `GET /api/notify-log`, `GET/POST /api/smtp-config`, `POST /api/smtp-test`, `GET/POST /api/webhook-config`, `GET /api/webhook-events`, `POST /api/notify/weekly-summary` | `getNotifyConfig()`, `setNotifyConfig()`, `notifyTest()`, `getNotifyLog()`, `getSmtpConfig()`, `setSmtpConfig()`, `smtpTest()`, `getWebhookConfig()`, `setWebhookConfig()`, `getWebhookEvents()`, `notifyWeeklySummary()` | notify 保存、test、SMTP 保存/test、webhook secret 保存、event 表示、weekly summary。 | `.notify_config`、`.smtp_config`、`.smtp_secret`、`.webhook_secret`、`.notify_log`、`.webhook_events.json`、`.config_log`。 | 保存系は対象 GET → `getConfigLog()`。test / summary は `getNotifyLog()`。 | secret は response / log / fixture へ平文出力しない。保存成功・失敗とも UI secret field を消去する。 |
| snapshots / rollback / maintenance | `GET /api/snapshots`, `GET /api/snapshots/{id}/download`, `DELETE /api/snapshots/{id}`, `POST /api/history/{id}/rollback`, `GET /api/maintenance`, `POST /api/maintenance/enable`, `POST /api/maintenance/disable` | `getSnapshots()`, `downloadSnapshot()`, `deleteSnapshot()`, `rollbackHistory()`, `getMaintenance()`, `enableMaintenance()`, `disableMaintenance()` | snapshot list/download/delete、rollback、maintenance enable/disable。 | `.snapshots/`、`.build_history`、`.build_logs/{new_id}.json`、`.maintenance`、`.config_log`。download は副作用なし。 | delete は `getSnapshots()`。rollback は `getHistory()` → `getStatus()`。maintenance は `getMaintenance()`。 | delete / rollback は確認必須。running rollback は `409`。maintenance enabled 中は build / rollback / 設定変更系を disabled。 |
| access / hooks / rules / pipeline / notes / layout | access、hooks、alert rules、tag rules、pipeline config、notes、dashboard layout の GET/POST/DELETE endpoint | 対応する [`docs/details/sdk.md`](details/sdk.md) §23 SDK method | 保存、追加、削除、notes 保存、dashboard layout 保存。 | `.access_control`、`.hooks`、`.alert_rules`、`.tag_rules`、`.pipeline_config`、`.notes`、`.dashboard_layout`、`.config_log`。 | 対象 GET → 変更系で config log 対象の場合は `getConfigLog()`。layout は `getDashboardLayout()` → `getDashboard()`。 | duplicate `409` は競合表示。validation `422` は field error。削除対象不在は `404`。 |
| tokens / audit / access logs / rate limit | `GET /api/tokens`, `POST /api/tokens`, `DELETE /api/tokens/{id}`, `GET /api/audit-log`, `GET /api/access-log`, `GET /api/api-access-log`, `GET/POST /api/api-rate-limit` | `getTokens()`, `createToken()`, `revokeToken()`, `getAuditLog()`, `getAccessLog()`, `getApiAccessLog()`, `getApiRateLimit()`, `setApiRateLimit()` | token 発行/失効、audit / access log 表示、rate limit 保存。 | `.api_tokens`、`.audit_log`、`.access_log`、`.api_access_log`、`.api_rate_state`、`.server_config`、`.config_log`。token 本体は作成時 response のみ。 | token 操作は `getTokens()` → `getAuditLog()`。rate limit は `getApiRateLimit()`。 | token 本体は再取得不可。`403` は logout しない。`429` は rate limit 表示し、同一操作を自動 retry しない。 |

**横断処理順契約：**

| 処理種別 | 固定順序 |
|----------|----------|
| 認証必須 JSON API | method / path 判定 → body 禁止判定 → JSON parse → 認証 / scope → rate limit → endpoint 固有 validation → read → write 計画 → atomic write → JSON Lines 追記 → response。 |
| read-only API | method / path 判定 → body 禁止判定 → 認証 / scope → query validation → read → 壊れた任意行除外 → response。read-only API は状態ファイルを書き換えない。 |
| UI 変更操作 | panel error / success 消去 → UI 入力検証 → 対象操作 disabled → SDK 呼び出し → 成功後再取得 → success 表示 → secret 消去 → disabled 再評価。 |
| UI 取得操作 | panel error 消去 → 対象操作 disabled → SDK 呼び出し → DOM 更新 → empty state 判定 → disabled 再評価。success 表示は行わない。 |
| SDK request | 引数検証 → path / query / body 生成 → Authorization 付与 → timeout 設定 → fetch → status 判定 → response parse → token 変化適用 → return / throw。 |
| multi-file write | 全入力検証 → 全対象 read → 全 write payload 生成 → [`docs/details/api.md`](details/api.md) §22.0d の Write 順に atomic write → JSON Lines 追記 → response。途中失敗時は未処理ファイルを書かない。 |

[`docs/details/runner.md`](details/runner.md) §27.21〜§27.38 の実装確認では、状態ファイル、endpoint、SDK method、UI 操作、外部公開構成が owner component 別の [`docs/details/*.md`](details/) 詳細本文責務に存在することを確認する。存在しない項目の追加可否と改訂先は、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.2a、[`docs/SPEC.md`](SPEC.md) ポリシー責務 §0、該当 owner / collaborator component 別の [`docs/details/*.md`](details/) 詳細本文責務を参照する。
