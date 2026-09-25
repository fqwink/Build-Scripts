# Adlaire CI — 詳細仕様入口

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務は、対象機能から owner component、collaborator component、詳細本文、fixture 証跡へ到達するための入口、共通固定値、対応表だけを管理する。方針・ポリシー・状態定義・着手可否は [`docs/SPEC.md`](SPEC.md)、現在状態と実装計画は [`docs/ROADMAP.md`](ROADMAP.md)、実在所在は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) を正本とする。

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務は owner component の処理本文、HTTP body、状態 schema、SDK method、UI DOM、fixture assertion、現在状態を再定義しない。

## 詳細仕様管理

| 確認対象 | 正本 |
|----------|------|
| 方針、ポリシー、状態語彙、遷移条件 | [`docs/SPEC.md`](SPEC.md) |
| 実装 artifact / 機能の現在状態と実装計画 | [`docs/ROADMAP.md`](ROADMAP.md) |
| owner component 別の詳細本文 | [`docs/details/`](details/) |
| fixture、expected、fake、実装検証証跡 | [`docs/details/fixture.md`](details/fixture.md) |
| 生成 HTML のデザイン | [`docs/DESIGN.md`](DESIGN.md) |
| 文書、実装、testdata の実在所在 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) |

## 詳細仕様入口責務

対象機能は [詳細節対応表](#0i-詳細節対応表) で唯一の owner component を確定し、[詳細仕様参照表](#0b-詳細仕様参照表) からその owner component の主本文へ移動する。

## 詳細仕様参照入口

1. [`docs/ROADMAP.md`](ROADMAP.md) で対象機能の現在状態と実装計画上の割当を確認する。
2. [詳細節対応表](#0i-詳細節対応表) で機能の owner component を一件に確定し、対応する詳細節を開く。
3. [詳細仕様参照表](#0b-詳細仕様参照表) で owner component の主本文が正しいことを確認する。
4. collaborator がある場合だけ、その component の詳細本文を境界確認として読む。
5. fixture と実装検証証跡は [`docs/details/fixture.md`](details/fixture.md) を読む。

## 0a. 詳細仕様の記載基準

記載義務、必須項目、不足時の扱いは [`docs/SPEC.md` 方針責務 §4.4](SPEC.md#sec-4-4) と [`docs/SPEC.md` ポリシー責務 §0 詳細仕様必須項目](SPEC.md#detail-contract-required-fields) を参照する。

## 0b.0 詳細仕様選択フロー

[詳細仕様カテゴリ](#0b01-詳細仕様カテゴリ) は検索入口、[詳細節対応表](#0i-詳細節対応表) は機能 owner と詳細節の対応、[詳細仕様参照表](#0b-詳細仕様参照表) は owner component と主本文の対応を担当する。カテゴリ名を owner component として扱わない。

## 0b.0.1 詳細仕様カテゴリ

| カテゴリ | owner 候補 | 詳細本文 |
|----------|------------|----------|
| Build | `builder` | [`docs/details/builder.md`](details/builder.md) |
| CI / 運用 | `runner`、`commitstatus`、`archive` | [`docs/details/runner.md`](details/runner.md)、[`docs/details/commitstatus.md`](details/commitstatus.md)、[`docs/details/archive.md`](details/archive.md) |
| 管理 | `api`、`sdk`、`ui`、`admin` | [`docs/details/api.md`](details/api.md)、[`docs/details/sdk.md`](details/sdk.md)、[`docs/details/ui.md`](details/ui.md)、[`docs/details/admin.md`](details/admin.md) |
| 状態 / 安全 | `statefile`、`security` | [`docs/details/statefile.md`](details/statefile.md)、[`docs/details/security.md`](details/security.md) |
| 配布 | `setup`、`release` | [`docs/details/setup.md`](details/setup.md)。`release` 専用詳細本文は未作成。 |
| 検証証跡 | fixture 証跡責務 | [`docs/details/fixture.md`](details/fixture.md) |
| MCP | `mcp` | 専用詳細仕様未作成。現在状態は [`docs/ROADMAP.md`](ROADMAP.md) を参照する。 |

## 0b. 詳細仕様参照表

| owner component | 主本文 | 主な責務 |
|-----------------|--------|----------|
| `builder` | [`docs/details/builder.md`](details/builder.md) | Markdown 入力から静的 Web サイトを生成する。 |
| `runner` | [`docs/details/runner.md`](details/runner.md) | GitHub 監視、build 実行、deploy、通知を調整する。 |
| `api` | [`docs/details/api.md`](details/api.md) | 管理 HTTP API の request / response と副作用境界を持つ。 |
| `admin` | [`docs/details/admin.md`](details/admin.md) | 管理 UI 静的配布物と配信境界を持つ。 |
| `sdk` | [`docs/details/sdk.md`](details/sdk.md) | JavaScript SDK の公開 method と HTTP 変換を持つ。 |
| `ui` | [`docs/details/ui.md`](details/ui.md) | 管理画面の DOM、操作、表示状態を持つ。 |
| `setup` | [`docs/details/setup.md`](details/setup.md) | バイナリ配置、systemd、更新、Release asset 受け入れを持つ。 |
| `release` | 専用詳細本文なし | GitHub Release 成果物の生成・公開前検証・公開の詳細本文は未作成。現在状態は [`docs/ROADMAP.md`](ROADMAP.md) を参照する。 |
| `statefile` | [`docs/details/statefile.md`](details/statefile.md) | 状態 schema、lock、atomic write、破損処理を持つ。 |
| `archive` | [`docs/details/archive.md`](details/archive.md) | log archive、snapshot、download、delete、rollback 実体を持つ。 |
| `commitstatus` | [`docs/details/commitstatus.md`](details/commitstatus.md) | GitHub Commit Status の payload と送信契約を持つ。 |
| `security` | [`docs/details/security.md`](details/security.md) | token、scope、session、TOTP、audit、rate limit を持つ。 |
| fixture 証跡責務 | [`docs/details/fixture.md`](details/fixture.md) | fixture、expected、fake、assertion、実装検証証跡を持つ。 |
| `mcp` | 専用詳細仕様なし | 現在状態は [`docs/ROADMAP.md`](ROADMAP.md) を参照する。 |

## 0b.1 owner component 別 owner / collaborator 境界管理

owner / collaborator 境界の規則は [`docs/SPEC.md` 方針責務 §4.2a](SPEC.md#sec-4-2a) を参照する。[詳細節対応表](#0i-詳細節対応表) の `owner` 列は機能から owner component を特定する入口、[詳細仕様参照表](#0b-詳細仕様参照表) は owner component から主本文を特定する入口とする。

## 0c. 実装前確認項目

実装着手可否は [`docs/SPEC.md` 方針責務 §4.7](SPEC.md#sec-4-7) と [`docs/SPEC.md` ポリシー責務 §0d](SPEC.md#0d-仕様凍結ポリシー) を正本とする。この入口では、対象機能が [詳細節対応表](#0i-詳細節対応表) に存在し、owner 詳細本文と fixture 証跡へ到達できることだけを確認する。

## 0d. 共通固定値

| 項目 | 固定値 |
|------|--------|
| Go 最小バージョン | Go `1.22` 以上。 |
| 文字コード | 入力、出力、状態ファイル、HTTP body は UTF-8。 |
| 改行 | 新規 text / JSON Lines は LF。CRLF 入力は読み込み時に LF として扱う。 |
| 機械処理時刻 | UTC の ISO 8601 秒精度 `YYYY-MM-DDTHH:MM:SSZ`。ミリ秒、ナノ秒、UTC 以外の offset、local timezone の保存を禁止する。ローカル時刻は UI 表示だけで使用する。 |
| CLI 終了コード | `0` 成功、`1` 一般エラー、`2` 入力・設定エラー、`3` 外部サービス・ネットワークエラー、`4` lock 形式不正などの継続不能な lock 異常。実行中 lock による通常 skip は `0`。 |
| CLI 共通 option | `--help` と `--version`。短縮 option は使用しない。 |
| 時刻ベース ID | prefix と UTC `YYYYMMDDHHmmss` を連結する。未衝突 ID に suffix は付けない。衝突時は `-001` から `-999` まで 3 桁連番を順に試し、上限到達時は既存 ID を上書きせず失敗とする。 |
| 出力成果物 manifest SHA-256 | 出力 root 配下の通常 file だけを entry とし、`/` 区切りの相対 path を UTF-8 byte 辞書順に並べる。各 file の SHA-256 を lowercase hex で算出し、各 entry の `relative_path + "\n" + file_sha256 + "\n"` を順に連結した byte 列全体の SHA-256 lowercase hex を `output_sha256` とする。directory は走査だけに使用し entry に含めない。symlink、link count 2 以上の hardlink、device、socket、FIFO を 1 件でも検出した場合は除外継続せず算出失敗とする。出力 root は symlink でない directory、全 path は valid UTF-8 とし、先頭 `/`、空 segment、`.`、`..`、backslash、NUL、CR、LF を含む相対 path は算出失敗とする。file は no-follow open 後の identity / type と読取前後の size / mtime が列挙時から不変の場合だけ採用し、走査中の追加・削除・置換・変更は算出失敗とする。通常 file が 0 件の出力 root は空 byte 列の SHA-256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` とする。 |

状態ファイルの lock、atomic write、権限、JSON 処理は [`docs/details/statefile.md`](details/statefile.md)、秘密情報は [`docs/details/security.md`](details/security.md)、外部依存とデータ交換形式は [`docs/SPEC.md`](SPEC.md) を正本とする。

## 0e. 完全実装検証マトリクス

| 対象 | 詳細本文 | fixture / 証跡 |
|------|----------|----------------|
| `builder` | [`docs/details/builder.md`](details/builder.md) | [`docs/details/fixture.md` §8a-F](details/fixture.md#8a-f-builder-初期受け入れ-fixture-契約)、[`docs/details/fixture.md` §28-F](details/fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) |
| `runner` | [`docs/details/runner.md`](details/runner.md) | [`docs/details/fixture.md` §15a-F](details/fixture.md#15a-f-runner-初期受け入れ-fixture-契約)、[`docs/details/fixture.md` §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| `api` | [`docs/details/api.md`](details/api.md) | [`docs/details/fixture.md` §22-F](details/fixture.md#22-f-api-fixture-契約)、[`docs/details/fixture.md` §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| `sdk` | [`docs/details/sdk.md`](details/sdk.md) | [`docs/details/fixture.md` §22-F](details/fixture.md#22-f-api-fixture-契約)、[`docs/details/fixture.md` §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| `ui` | [`docs/details/ui.md`](details/ui.md) | [`docs/details/fixture.md` §22-F](details/fixture.md#22-f-api-fixture-契約)、[`docs/details/fixture.md` §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| `admin` | [`docs/details/admin.md`](details/admin.md) | [`docs/details/fixture.md` §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| `statefile` | [`docs/details/statefile.md`](details/statefile.md) | [`docs/details/fixture.md` §22-F](details/fixture.md#22-f-api-fixture-契約)、[`docs/details/fixture.md` §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| `security` | [`docs/details/security.md`](details/security.md) | [`docs/details/fixture.md` §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| `archive` | [`docs/details/archive.md`](details/archive.md) | [`docs/details/fixture.md` §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| `commitstatus` | [`docs/details/commitstatus.md`](details/commitstatus.md) | [`docs/details/fixture.md` §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| `setup` | [`docs/details/setup.md`](details/setup.md) | [`docs/details/fixture.md`](details/fixture.md) |

完了判定と状態遷移は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー)、現在状態は [`docs/ROADMAP.md`](ROADMAP.md) を参照する。

## 0f. 仕様策定完了チェック

仕様策定の完了条件は [`docs/SPEC.md` 方針責務 §4.4](SPEC.md#sec-4-4)、[`docs/SPEC.md` 方針責務 §4.7](SPEC.md#sec-4-7)、[`docs/SPEC.md` ポリシー責務 §0d](SPEC.md#0d-仕様凍結ポリシー) を正本とする。この入口では、[詳細仕様参照表](#0b-詳細仕様参照表)、[`docs/SPEC.md` ポリシー責務 §0 詳細仕様必須項目](SPEC.md#detail-contract-required-fields)、[詳細節対応表](#0i-詳細節対応表)、[完全実装検証マトリクス](#0e-完全実装検証マトリクス) の参照が揃っていることだけを確認する。

## 0i. 詳細節対応表

[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i 詳細節対応表](DETAIL_INDEX.md#0i-詳細節対応表) は対象機能から唯一の owner と関連詳細本文へ移動するための対応表である。`owner` 列だけが機能 owner の正本であり、collaborator の接続境界と担当処理は owner の該当詳細節を参照する。現在状態は [`docs/ROADMAP.md`](ROADMAP.md)、受け入れ assertion は [`docs/details/fixture.md`](details/fixture.md) を正本とし、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i 詳細節対応表](DETAIL_INDEX.md#0i-詳細節対応表) では再掲しない。

<a id="0i1-builder--静的-web-サイト出力"></a>
**0i.1 Builder / 静的 Web サイト出力：**

| 機能 | owner | 詳細本文 |
|------|-------|----------|
| 出力サイトサイズ警告閾値 | `builder` | [`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法)、[`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e) |
| 出力サイトへのビルドメタ埋め込み | `builder` | [`docs/details/builder.md` 詳細本文責務 §2](details/builder.md#2-ファイルパス設定)、[`docs/details/builder.md` 詳細本文責務 §5](details/builder.md#5-静的-web-サイト出力構造)、[`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/builder.md` 詳細本文責務 §27.4](details/builder.md#sec-27-4) |
| 変換レポート出力 | `builder` | [`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e) |
| シンタックスハイライト | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.8](details/builder.md#sec-7-8) |
| 本文内全文検索 | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.9](details/builder.md#sec-7-9) |
| アンカーリンク自動検証 | `builder` | [`docs/details/builder.md` 詳細本文責務 §4.3](details/builder.md#sec-4-3)、[`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法) |
| コードブロックの折りたたみ | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.10](details/builder.md#sec-7-10) |
| 印刷スタイル（`@media print`） | `builder` | [`docs/details/builder.md` 詳細本文責務 §6](details/builder.md#6-css-クラス一覧)、[`docs/DESIGN.md` デザイン責務 Builder 拡張コンポーネント視覚契約](DESIGN.md#builder-拡張コンポーネント視覚契約) |
| 静的 Web サイト出力 | `builder` | [`docs/details/builder.md` 詳細本文責務 §2](details/builder.md#2-ファイルパス設定)、[`docs/details/builder.md` 詳細本文責務 §5](details/builder.md#5-静的-web-サイト出力構造)、[`docs/details/builder.md` 詳細本文責務 §6](details/builder.md#6-css-クラス一覧)、[`docs/details/builder.md` 詳細本文責務 §7](details/builder.md#7-javascript-機能)、[`docs/DESIGN.md` デザイン責務](DESIGN.md) |
| テーマコンポーネント | `builder` | [`docs/details/builder.md` 詳細本文責務 §5](details/builder.md#5-静的-web-サイト出力構造)、[`docs/details/builder.md` 詳細本文責務 §6](details/builder.md#6-css-クラス一覧)、[`docs/details/builder.md` 詳細本文責務 §7](details/builder.md#7-javascript-機能)、[`docs/DESIGN.md` デザイン責務](DESIGN.md) |
| 外部リンクの自動処理 | `builder` | [`docs/details/builder.md` 詳細本文責務 §4.3](details/builder.md#sec-4-3) |
| 読み取り進捗バー | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.13](details/builder.md#sec-7-13) |
| コードブロックのコピーボタン | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.6](details/builder.md#sec-7-6) |
| 見出しアンカーリンクコピー | `builder` | [`docs/details/builder.md` 詳細本文責務 §3](details/builder.md#3-処理パイプライン)、[`docs/details/builder.md` 詳細本文責務 §7.11](details/builder.md#sec-7-11) |
| TOC 開閉状態の永続化 | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.3](details/builder.md#sec-7-3) |
| 見出しスラグ重複解決 | `builder` | [`docs/details/builder.md` 詳細本文責務 §4.5](details/builder.md#sec-4-5) |
| 前後章ナビゲーションボタン | `builder` | [`docs/details/builder.md` 詳細本文責務 §4.5](details/builder.md#sec-4-5)、[`docs/details/builder.md` 詳細本文責務 §5](details/builder.md#5-静的-web-サイト出力構造)、[`docs/details/builder.md` 詳細本文責務 §7.15](details/builder.md#sec-7-15) |
| 内部リンク整合性チェック | `builder` | [`docs/details/builder.md` 詳細本文責務 §4.3](details/builder.md#sec-4-3)、[`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法) |
| 見出し階層スキップ警告 | `builder` | [`docs/details/builder.md` 詳細本文責務 §4.5](details/builder.md#sec-4-5)、[`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法) |
| 読了時間推計と表示 | `builder` | [`docs/details/builder.md` 詳細本文責務 §4.5](details/builder.md#sec-4-5)、[`docs/details/builder.md` 詳細本文責務 §5](details/builder.md#5-静的-web-サイト出力構造)、[`docs/details/builder.md` 詳細本文責務 §6](details/builder.md#6-css-クラス一覧)、[`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法) |
| テーブルのソート機能 | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.14](details/builder.md#sec-7-14) |
| キーボードショートカット | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.12](details/builder.md#sec-7-12) |
| ビルドキャッシュ | `builder` | [`docs/details/builder.md` 詳細本文責務 §5](details/builder.md#5-静的-web-サイト出力構造)、[`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法)、[`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/builder.md` 詳細本文責務 §27.25](details/builder.md#sec-27-25) |
| 依存ファイルトラッキング | `builder` | [`docs/details/builder.md` 詳細本文責務 §4.3](details/builder.md#sec-4-3)、[`docs/details/builder.md` 詳細本文責務 §5](details/builder.md#5-静的-web-サイト出力構造)、[`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/builder.md` 詳細本文責務 §27.28](details/builder.md#sec-27-28) |
| 差分ビルド | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.1](details/builder.md#sec-28-1) |
| 複数出力形式 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.2](details/builder.md#sec-28-2) |
| Markdown 拡張記法サポート | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.3](details/builder.md#sec-28-3) |
| コードブロック行番号表示 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.4](details/builder.md#sec-28-4) |
| 見出しの自動採番 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.5](details/builder.md#sec-28-5) |
| セクション折りたたみ | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.6](details/builder.md#sec-28-6) |
| TOC 深さ制御 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.7](details/builder.md#sec-28-7) |
| 最終更新日の自動埋め込み | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.8](details/builder.md#sec-28-8) |
| diff ハイライト | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.9](details/builder.md#sec-28-9) |
| 画像の遅延読み込み | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.10](details/builder.md#sec-28-10) |
| カスタムメタタグ注入 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.11](details/builder.md#sec-28-11) |
| ライトモード固定 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.12](details/builder.md#sec-28-12) |
| コードブロックのファイル名表示 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.13](details/builder.md#sec-28-13) |
| テンプレート変数展開 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.14](details/builder.md#sec-28-14) |
| HTML ミニファイ | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.15](details/builder.md#sec-28-15) |
| TOC ハイライト追従 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.16](details/builder.md#sec-28-16) |
| Mermaid ダイアグラム描画 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.17](details/builder.md#sec-28-17) |
| 脚注サポート | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.18](details/builder.md#sec-28-18) |
| インライン数式レンダリング | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.19](details/builder.md#sec-28-19) |
| ページ内ナビゲーション履歴 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.20](details/builder.md#sec-28-20) |
| 読み上げ対応（アクセシビリティ） | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.21](details/builder.md#sec-28-21) |
| 画像ライトボックス | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.22](details/builder.md#sec-28-22) |
| 印刷時 QR コード挿入 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.23](details/builder.md#sec-28-23) |
| 定義リストサポート | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.24](details/builder.md#sec-28-24) |
| タスクリストサポート | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.25](details/builder.md#sec-28-25) |

<a id="0i2-runner--ci-実行"></a>
**0i.2 Runner / CI 実行：**

| 機能 | owner | 詳細本文 |
|------|-------|----------|
| ビルドタイムアウト | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e) |
| ビルドログのファイル保存 | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ) |
| ネットワーク断時の再試行 | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー) |
| GitHub API レート制限自動待機 | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e) |
| 転送後リモート整合性検証 | `runner` | [`docs/details/runner.md` 詳細本文責務 §14a](details/runner.md#14a-ssh-サイト転送)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー) |
| マルチブランチビルド | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e) |
| ビルドログ世代管理 | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ) |
| ビルド出力の外部転送 | `runner` | [`docs/details/runner.md` 詳細本文責務 §14a](details/runner.md#14a-ssh-サイト転送)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー) |
| ビルドクールダウン | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー) |
| ビルド前の事前チェック | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/setup.md` 詳細本文責務 §26](details/setup.md#26-セットアップアップデート手順) |
| 定期強制ビルド | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e) |
| ビルド中重複スキップ | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー) |
| GitHub PAT 有効期限の事前警告 | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e) |
| コミット情報のビルドログ記録 | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ) |
| GitHub API 連続失敗によるサーキットブレーカー | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e) |
| 設定ファイル起動時整合性チェック | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/runner.md` 詳細本文責務 §27.10](details/runner.md#sec-27-10) |
| ビルドステータスファイル出力 | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.8](details/runner.md#sec-27-8) |
| ビルドトリガー種別の記録 | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/runner.md` 詳細本文責務 §27.9](details/runner.md#sec-27-9) |
| GitHub Commit Status API | `commitstatus` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/commitstatus.md` 詳細本文責務 §27.1](details/commitstatus.md#sec-27-1) |
| ドライラン実行モード | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/runner.md` 詳細本文責務 §27.2](details/runner.md#sec-27-2) |
| ビルド失敗時の自動リトライ | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/runner.md` 詳細本文責務 §27.3](details/runner.md#sec-27-3) |
| Webhook 通知失敗リトライキュー | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §16](details/runner.md#16-systemd-タイマー参照) |
| 複数ファイル監視 | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/runner.md` 詳細本文責務 §27.21](details/runner.md#sec-27-21) |
| ビルドパイプライン YAML 定義 | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.22](details/runner.md#sec-27-22) |
| ローカルファイル監視モード | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §27.23](details/runner.md#sec-27-23) |
| タグ付きコミットのみビルド | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.24](details/runner.md#sec-27-24) |
| 並列マルチターゲットビルド | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §14a](details/runner.md#14a-ssh-サイト転送)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/runner.md` 詳細本文責務 §27.26](details/runner.md#sec-27-26) |
| ビルド前後フック | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.27](details/runner.md#sec-27-27) |
| リモートビルド対応 | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §14a](details/runner.md#14a-ssh-サイト転送)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/runner.md` 詳細本文責務 §27.29](details/runner.md#sec-27-29) |
| ブランチ別環境変数 | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.31](details/runner.md#sec-27-31) |
| ビルド通知連携 | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/runner.md` 詳細本文責務 §16](details/runner.md#16-systemd-タイマー参照)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.32](details/runner.md#sec-27-32) |
| ビルド時間トレンド記録 | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.33](details/runner.md#sec-27-33) |
| ビルド依存チェーン | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.34](details/runner.md#sec-27-34) |
| ビルド優先度キュー | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.35](details/runner.md#sec-27-35) |
| 失敗原因の自動分類 | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.36](details/runner.md#sec-27-36) |
| ビルド実行環境の記録 | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/runner.md` 詳細本文責務 §27.37](details/runner.md#sec-27-37) |
| ビルド所要時間の異常検知 | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/runner.md` 詳細本文責務 §16](details/runner.md#16-systemd-タイマー参照)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.38](details/runner.md#sec-27-38) |

<a id="0i3-api--sdk--ui"></a>
**0i.3 API / SDK / UI：**

| 機能 | owner | 詳細本文 |
|------|-------|----------|
| ポーリング間隔の動的変更 | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/setup.md` 詳細本文責務 §26](details/setup.md#26-セットアップアップデート手順)、[`docs/details/api.md` 詳細本文責務 §27.11](details/api.md#sec-27-11) |
| GitHub Webhook 受信 | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/api.md` 詳細本文責務 Webhook 受信境界](details/api.md#webhook-receive-overview)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §27.12](details/api.md#sec-27-12) |
| 設定バリデーション API | `api` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/api.md` 詳細本文責務 §27.5](details/api.md#sec-27-5) |
| API アクセスログ | `api` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/api.md` 詳細本文責務 §27.6](details/api.md#sec-27-6) |
| Webhook イベントログ | `api` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/api.md` 詳細本文責務 Webhook 受信境界](details/api.md#webhook-receive-overview)、[`docs/details/api.md` 詳細本文責務 §27.13](details/api.md#sec-27-13) |
| ヘルスチェックエンドポイント | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/api.md` 詳細本文責務 §27.16](details/api.md#sec-27-16) |
| Webhook イベント一覧取得 API | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/api.md` 詳細本文責務 §27.13](details/api.md#sec-27-13) |
| ビルドログ重大度フィルター | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/api.md` 詳細本文責務 §27.17](details/api.md#sec-27-17) |
| ブランチ設定の動的変更 API | `api` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/api.md` 詳細本文責務 §27.18](details/api.md#sec-27-18) |
| 週次ビルドサマリー Webhook | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §16](details/runner.md#16-systemd-タイマー参照)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.19](details/runner.md#sec-27-19) |
| 設定変更の詳細 diff 記録 | `api` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/api.md` 詳細本文責務 §27.20](details/api.md#sec-27-20) |
| ビルド承認フロー | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/runner.md` 詳細本文責務 §16](details/runner.md#16-systemd-タイマー参照)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/api.md` 詳細本文責務 §27.30](details/api.md#sec-27-30) |
| 通知チャンネル管理 | `runner` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.32](details/runner.md#sec-27-32)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様) |
| ビルドログの保存済み有限 SSE 配信 | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様) |
| ビルド統計ダッシュボード | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様) |
| IP アドレス制限 | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様) |
| ビルドキュー可視化 | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様) |
| メンテナンスモード | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様) |
| アラート閾値設定 | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様) |
| バックアップ／リストア | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様) |

<a id="0i4-archive--artifact--security"></a>
**0i.4 Archive / Artifact / Security：**

| 機能 | owner | 詳細本文 |
|------|-------|----------|
| ビルドログのアーカイブ圧縮 | `archive` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/archive.md` 詳細本文責務 §27.7](details/archive.md#sec-27-7) |
| ビルド所要時間の記録と統計 API | `runner` | [`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.14](details/runner.md#sec-27-14) |
| ビルドアーティファクト世代管理 | `archive` | [`docs/details/runner.md` 詳細本文責務 §14b](details/runner.md#14b-スナップショット管理)、[`docs/details/archive.md` 詳細本文責務 §27.15](details/archive.md#sec-27-15)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e) |
| ビルドアーティファクト管理 | `archive` | [`docs/details/runner.md` 詳細本文責務 §14b](details/runner.md#14b-スナップショット管理)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/archive.md` 詳細本文責務 §27.15](details/archive.md#sec-27-15) |
| ビルドトリガー専用 API スコープ | `security` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/security.md` 詳細本文責務 §27.42](details/security.md#sec-27-42) |
| API キー管理 | `security` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/security.md` 詳細本文責務 §27.43](details/security.md#sec-27-43) |
| 監査ログ | `security` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.30](details/runner.md#sec-27-30)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/security.md` 詳細本文責務 §27.44](details/security.md#sec-27-44) |
| セッションタイムアウト変更設定 | `security` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/api.md` 詳細本文責務 §25](details/api.md#25-認証-実装仕様)、[`docs/details/security.md` 詳細本文責務 §27.45](details/security.md#sec-27-45) |
| TOTP 二要素認証 | `security` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/api.md` 詳細本文責務 §25](details/api.md#25-認証-実装仕様)、[`docs/details/security.md` 詳細本文責務 §27.46](details/security.md#sec-27-46) |
| API レート制限 | `security` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/security.md` 詳細本文責務 §27.47](details/security.md#sec-27-47) |

<a id="0i5-setup--release"></a>
**0i.5 Setup / Release：**

| 機能 | owner | 詳細本文 |
|------|-------|----------|
| 初回セットアップ | `setup` | [`docs/details/setup.md` 詳細本文責務 §26.1](details/setup.md#sec-26-1)〜[§26.3](details/setup.md#sec-26-3) |
| 管理 API 導入 | `setup` | [`docs/details/setup.md` 詳細本文責務 §26.3b](details/setup.md#sec-26-3b) |
| バイナリアップデート | `setup` | [`docs/details/setup.md` 詳細本文責務 §26.4](details/setup.md#sec-26-4) |
| アップデート rollback | `setup` | [`docs/details/setup.md` 詳細本文責務 §26.4](details/setup.md#sec-26-4) |
| GitHub Release 成果物生成・公開前検証・公開 | `release` | 専用詳細本文未作成。成果物の受け入れ側契約は [`docs/details/setup.md` 詳細本文責務 §26.2a](details/setup.md#sec-26-2a) を参照する。 |

## 詳細仕様セット構成

詳細仕様セットは [詳細仕様参照表](#0b-詳細仕様参照表) に列挙した owner component 詳細本文と [`docs/details/fixture.md`](details/fixture.md) で構成する。実在するファイルの一覧は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) を参照する。

## 0j. リポジトリ内ソース配置

標準ディレクトリ構成は [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3)、実在する文書・実装・testdata・未作成 path は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) を正本とする。[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務では tree、現在配置、将来 path を再掲しない。

## owner component 別詳細本文責務索引

[詳細仕様参照表](#0b-詳細仕様参照表) を owner component 別詳細本文責務の唯一の索引とする。
