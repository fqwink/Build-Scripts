# Build-Scripts — 文書索引

このファイルは、Build-Scripts リポジトリ内の文書・実装ファイルの参照先と役割を整理する索引である。

仕様・詳細仕様・デザイン正本・本索引は `docs/` 配下に集約する。ルールブック [`AGENTS.md`](../AGENTS.md) と入口文書 [`README.md`](../README.md) はリポジトリ root に置く。

## 仕様構造管理

仕様構造の管理は、文書の所在、正本範囲、記載先、変更手順、完了条件を分離して扱う。

| 管理領域 | 確認する節 | 判断する内容 |
|----------|------------|--------------|
| 構造地図 | `仕様構造マップ` | 読みたい目的から参照先を選ぶための地図。 |
| 階層定義 | `仕様構造` | 文書階層、文書の位置付け、役割の定義。 |
| 判断 | `仕様判断フロー` | 目的ごとに最初に読む文書、次に確認する文書。 |
| 記載先 | `仕様記載先マトリクス` | 何をどの文書に書き、どこに書かないか。 |
| 変更手順 | `仕様変更手順` | 文書構造変更時の作業順序。 |
| 完了条件 | `仕様構造完了条件` | 完了報告前に満たすべき確認項目。 |

本節は、仕様構造を変更するための統制入口である。本節自体は仕様本文、詳細仕様本文、実装状態、実装可否、ロードマップを定義しない。

## 仕様構造マップ

本節は、読みたい目的から参照先を選ぶための地図である。文書階層そのものの定義は `仕様構造` を正とする。

| 読みたいもの | 参照先 | 参照理由 |
|--------------|--------|----------|
| 作業してよい条件 | [`AGENTS.md`](../AGENTS.md) | 承認、Git、PR、文書整合の最上位ルールを確認する。 |
| 初見向け概要 | [`README.md`](../README.md) | 最小限の入口と主要参照先を確認する。 |
| 文書構造 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) | 文書の役割、所在、正本範囲、実装ファイル所在を確認する。 |
| 方針・ポリシー | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | 方針、ポリシー、正本関係、禁止事項、リリース判断を確認する。 |
| 状態・計画 | [`docs/ROADMAP.md`](ROADMAP.md) | 実装状態、実装可否、Phase、機能インベントリ、将来計画、追加仕様化機能参照、横断補足契約を確認する。 |
| 詳細仕様入口 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) | 詳細仕様の読み方、共通固定値、対応表、ソース配置を確認する。 |
| 詳細仕様本文 | [`docs/details/*.md`](details/) | owner component の入出力、状態、処理順序、異常系、検証条件を確認する。 |
| デザイン責務 | [`docs/DESIGN.md`](DESIGN.md) | 生成静的 Web サイトのデザイン方針と視覚仕様を確認する。 |
| 実装所在 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) の `仕様化済みコンポーネント` | 実装ファイル、テスト、fixture の所在と状態を確認する。 |

## 仕様構造

本節は、仕様文書の階層と位置付けを定義する。目的別の参照先選択は `仕様構造マップ` を入口とする。

仕様文書の構造は、作業ルール、入口、索引、方針・ポリシー正本、デザイン正本、状態・計画正本、詳細入口、責務 component 別本文に分ける。

| 階層 | ファイル | 位置付け | 主な用途 |
|------|----------|----------|----------|
| 作業ルール | [`AGENTS.md`](../AGENTS.md) | 最上位ルールブック | 承認、Git 運用、仕様書管理、文書整合ルールを判断する。 |
| 入口 | [`README.md`](../README.md) | 初見向け入口 | リポジトリ概要、読む順番、主要ファイルを把握する。 |
| 索引 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) | 文書・実装ファイル索引 | 文書配置、正本関係、実装ファイル所在を確認する。 |
| マスター仕様 | [`docs/SPEC.md`](SPEC.md) | 方針責務・ポリシー責務正本 | 方針、ポリシー、正本関係、禁止事項、リリース判断を確認する。 |
| ロードマップ | [`docs/ROADMAP.md`](ROADMAP.md) | 状態・計画正本 | 実装状態、実装可否、Phase、機能インベントリ、将来計画、昇格手順、追加仕様化機能参照、横断補足契約を判断する。 |
| 詳細仕様入口 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) | 詳細仕様入口責務 | 詳細仕様の読み方、共通固定値、対応表、ソース配置を確認する。 |
| 詳細仕様本文 | [`docs/details/*.md`](details/) | owner component 別本文 | 対象 component の入出力、状態、処理順序、異常系、検証条件を実装単位で確認する。 |
| デザイン正本 | [`docs/DESIGN.md`](DESIGN.md) | デザイン責務正本 | 生成静的 Web サイトのデザイン方針、レイアウト、色、タイポグラフィ、TOC、コードブロック等の視覚仕様を判断する。 |

仕様判断では、上位階層が下位階層を置き換えるのではなく、各階層の正本範囲だけを参照する。[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) は所在と役割の索引であり、仕様本文を定義しない。

## 仕様判断フロー

目的別の参照先は以下に固定する。迷った場合は、最初に本表で参照先を決め、参照先の正本範囲だけを確認する。

| 目的 | 最初に読む文書 | 次に確認する文書 | 判断内容 |
|------|----------------|------------------|----------|
| 作業ルール、承認、Git 運用を確認したい | [`AGENTS.md`](../AGENTS.md) | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) | 作業開始可否、変更承認、PR 作成、文書整合の手順を判断する。 |
| リポジトリ全体の文書構造を把握したい | [`README.md`](../README.md) | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) | どの文書が何を持つか、どの順番で読むかを判断する。 |
| 方針、ポリシー、禁止事項を判断したい | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | [`docs/ROADMAP.md`](ROADMAP.md) | 対象領域に適用する原則と制約を判断する。 |
| 実装状態、実装可否、Phase、将来計画を判断したい | [`docs/ROADMAP.md`](ROADMAP.md) | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) | 対象が未仕様化、将来計画、改訂予定、仕様化済み・未実装、実装中・検証未完了、実装済みのどれかを判断する。 |
| 詳細仕様本文を探したい | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) | 対象 owner component の [`docs/details/*.md`](details/) | 対象機能の owner component、参照節、受け入れ条件を判断する。 |
| 実装ファイル、テスト、fixture の所在を確認したい | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0j | 現行実装ファイル、将来追加予定 path、標準配置の扱いを判断する。 |
| component の入出力、状態、処理順序、異常系、検証条件を確認したい | 対象 owner component の [`docs/details/*.md`](details/) | collaborator component の [`docs/details/*.md`](details/) | 実装時に従う具体仕様と collaborator 境界を判断する。 |
| 生成静的 Web サイトの見た目を確認したい | [`docs/DESIGN.md`](DESIGN.md) | [`docs/details/builder.md`](details/builder.md) | レイアウト、色、タイポグラフィ、TOC、コードブロック等の視覚仕様を判断する。 |

本表は文書選択の判断フローであり、各文書の正本範囲を拡張しない。

## 仕様判断ルール

仕様判断では、目的を先に確定し、目的に対応する正本だけを読む。

| 判断ルール | 内容 |
|------------|------|
| 作業ルール優先 | 作業可否、承認、Git 操作、PR 作成は常に [`AGENTS.md`](../AGENTS.md) を正とする。 |
| 状態判断優先 | 未仕様化、将来計画、改訂予定、仕様化済み・未実装、実装中・検証未完了、実装済みの判断は [`docs/ROADMAP.md`](ROADMAP.md) を正とする。 |
| 詳細本文優先 | 入出力、状態、処理順序、異常系、検証条件は owner component の [`docs/details/*.md`](details/) を正とする。 |
| 索引限定 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) は所在と役割の索引であり、仕様本文を定義しない。 |
| 入口限定 | [`README.md`](../README.md) は入口であり、詳細ルールや詳細仕様本文を重複定義しない。 |
| デザイン責務限定 | [`docs/DESIGN.md`](DESIGN.md) は生成静的 Web サイトのデザイン関係の正本であり、機能仕様、運用仕様、API 仕様の正本ではない。 |

## 仕様記載先マトリクス

仕様、詳細仕様、デザイン正本、索引を改訂する場合は、下表に従って記載先を選ぶ。

| 書く内容 | 書く場所 | 書いてはいけない場所 |
|----------|----------|----------------------|
| 作業ルール、承認条件、Git 運用、PR 作成、文書整合ルール | [`AGENTS.md`](../AGENTS.md) | [`README.md`](../README.md)、[`docs/SPEC.md`](SPEC.md)、[`docs/ROADMAP.md`](ROADMAP.md)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、[`docs/details/*.md`](details/) |
| 生成 HTML のデザイン関係を除く方針、ポリシー、禁止事項、リリース判断、正本関係 | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | [`docs/ROADMAP.md`](ROADMAP.md)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、[`docs/details/*.md`](details/)、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、[`README.md`](../README.md) |
| 実装状態、実装可否、Phase、機能インベントリ、将来計画、昇格手順、追加仕様化機能参照、横断補足契約 | [`docs/ROADMAP.md`](ROADMAP.md) | [`docs/SPEC.md`](SPEC.md)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、[`docs/details/*.md`](details/)、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、[`README.md`](../README.md) |
| 詳細仕様の入口、読み順、共通固定値、実装前確認項目、検証マトリクス、詳細節対応表、リポジトリ内ソース配置 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) | [`docs/SPEC.md`](SPEC.md)、[`docs/ROADMAP.md`](ROADMAP.md)、[`docs/details/*.md`](details/)、[`README.md`](../README.md) |
| owner component の入出力、状態、処理順序、異常系、セキュリティ制約、検証条件 | owner component 別の [`docs/details/*.md`](details/) | [`docs/SPEC.md`](SPEC.md)、[`docs/ROADMAP.md`](ROADMAP.md)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、[`README.md`](../README.md) |
| 文書配置、読む順番、正本関係、実装ファイル所在、仕様構造の索引 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) | [`docs/SPEC.md`](SPEC.md)、[`docs/ROADMAP.md`](ROADMAP.md)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、[`docs/details/*.md`](details/) |
| 生成静的 Web サイトのデザイン方針、レイアウト、色、タイポグラフィ、TOC、コードブロック等の視覚仕様 | [`docs/DESIGN.md`](DESIGN.md) | [`docs/SPEC.md`](SPEC.md)、[`docs/ROADMAP.md`](ROADMAP.md)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、[`docs/details/*.md`](details/) |
| 初見向け概要、最小限の読む順番、主要ファイル案内 | [`README.md`](../README.md) | [`docs/SPEC.md`](SPEC.md)、[`docs/ROADMAP.md`](ROADMAP.md)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、[`docs/details/*.md`](details/) |

[`README.md`](../README.md) は入口であり、詳細ルール、詳細仕様本文、実装状態、ロードマップ、API 仕様、状態 schema、検証 matrix を重複定義しない。[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) は索引であり、仕様本文、詳細仕様本文、実装可否、ロードマップ状態を定義しない。

## 仕様記載ルール

仕様構造を崩さないため、記載先は以下の順で確定する。

1. 書く内容が作業ルールか、仕様本文か、詳細仕様本文か、デザイン本文か、索引かを分類する。
2. `仕様記載先マトリクス` で書く場所を確定する。
3. 書いてはいけない場所に同じ意味の本文が残る場合は、重複として整理する。
4. 参照だけで足りる場合は、本文を複製せず、正本への参照に留める。
5. 実装状態、実装可否、Phase、将来計画を動かす場合は、[`docs/ROADMAP.md`](ROADMAP.md) の正本範囲として扱う。

## 仕様変更手順

仕様構造、文書配置、ファイル名、参照先、正本関係を変更する場合は、以下の順で作業する。

1. 変更目的を `仕様判断フロー` で分類する。
2. 記載先を `仕様記載先マトリクス` で確定する。
3. 確定した正本文書だけを編集する。
4. 文書名、節名、正本範囲、実装ファイル所在に影響がある場合は、必要な索引と参照だけを更新する。
5. 廃止済みファイル名、廃止済み節名、不要参照、削除済み文書名が残っていないことを `rg` で確認する。
6. 文書構造整理だけの作業では、実装ファイル、testdata、fixture を変更しない。
7. `仕様構造完了条件` をすべて満たしてから完了扱いにする。

上記手順は、仕様本文の意味、実装状態、ロードマップ状態を変更する許可ではない。仕様本文の意味を変更する場合は、変更内容に対応する正本文書のルールに従う。

## 仕様変更ルール

仕様構造変更では、変更の種類ごとに確認対象を固定する。

| 変更の種類 | 必ず確認する文書 | 確認内容 |
|------------|------------------|----------|
| ファイル名変更 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、[`README.md`](../README.md)、[`AGENTS.md`](../AGENTS.md)、[`docs/SPEC.md`](SPEC.md)、[`docs/ROADMAP.md`](ROADMAP.md)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) | 廃止済みファイル名参照が残っていないこと。 |
| 正本関係変更 | [`AGENTS.md`](../AGENTS.md)、[`docs/SPEC.md`](SPEC.md)、[`docs/DESIGN.md`](DESIGN.md)、[`docs/ROADMAP.md`](ROADMAP.md)、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) | 作業ルール、仕様正本、デザイン正本、索引の記載が矛盾しないこと。 |
| 詳細仕様責務整理 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、[`docs/details/*.md`](details/)、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) | owner component、collaborator、参照表、本文配置が一致すること。 |
| README 整理 | [`README.md`](../README.md)、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) | README が入口に留まり、詳細ルールを重複定義していないこと。 |
| 実装所在整理 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0j | 実装ファイル、テスト、fixture の所在と状態が一致すること。 |

## 仕様構造完了条件

仕様構造、文書配置、ファイル名、参照先、正本関係を変更する作業は、以下をすべて満たすまで完了扱いにしない。

| 確認項目 | 完了条件 |
|----------|----------|
| 記載先 | 変更した内容が `仕様記載先マトリクス` の `書く場所` に一致している。 |
| 重複禁止 | [`README.md`](../README.md) に詳細ルール、詳細仕様本文、実装状態、ロードマップ、API 仕様、状態 schema、検証 matrix を重複定義していない。 |
| 索引境界 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) に仕様本文、詳細仕様本文、実装可否、ロードマップ状態を定義していない。 |
| 詳細入口境界 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) に方針、ポリシー、実装状態、実装可否、ロードマップ状態を定義していない。 |
| 詳細本文境界 | [`docs/details/*.md`](details/) に owner component 以外の主本文を混在させていない。 |
| 参照整合 | 廃止済みファイル名、廃止済み節名、不要参照、削除済み文書名が `rg` で残っていない。 |
| 変更範囲 | 文書構造整理だけの作業で、実装ファイル、testdata、fixture を変更していない。 |
| PR 範囲 | 同一目的の文書構造変更が既存 PR に集約され、並行 PR と同一ファイル編集を発生させていない。 |

上記のいずれかを満たせない場合は、完了報告せず、該当文書の正本範囲、記載先、参照先を再整備する。

## 仕様完了報告ルール

仕様構造整理の完了報告では、以下を満たす。

| 完了報告に含める内容 | 内容 |
|----------------------|------|
| 変更対象 | 変更した文書名を明記する。 |
| 変更内容 | 追加、削除、移動、簡潔化、参照更新の内容を明記する。 |
| 正本範囲 | 仕様本文を変えたのか、索引を変えたのか、入口を変えたのかを明記する。 |
| 検証 | `git diff --check`、廃止参照検索、変更範囲確認を明記する。 |
| 未実施 | docs-only で実装検証を未実施にした場合は理由を明記する。 |

## リポジトリ文書索引

本節以降は、リポジトリ内文書と実装ファイル所在の索引である。仕様判断の本文は各正本を参照する。

| 索引領域 | 確認する節 | 判断する内容 |
|----------|------------|--------------|
| 読む順番 | `読む順番` | 作業開始から詳細仕様本文までの確認順序。 |
| 文書一覧 | `文書一覧` | 各文書の役割と所在。 |
| 実装所在入口 | `実装ファイル索引` | 実装ファイル所在の判断原則。 |
| 詳細仕様管理 | `詳細仕様管理` | owner component 別詳細仕様ファイルの配置と状態。 |
| 実装ファイル一覧 | `仕様化済みコンポーネント` | 実装ファイル、テスト、fixture の所在と状態。 |
| 正本関係 | `正本関係` | 判断対象ごとの正本。 |
| 整合ガード | `整合ガードレール` | 文書整合を壊さないための禁止事項。 |
| 整合メモ | `整合メモ` | 現行状態に関する注意点。 |

## 読む順番

| 順序 | ファイル | 目的 |
|------|----------|------|
| 1 | [`AGENTS.md`](../AGENTS.md) | 作業ルール、承認、Git 運用、文書整合ルールを確認する。 |
| 2 | [`README.md`](../README.md) | 初見向けの概要、読む順番、主要ファイルを確認する。 |
| 3 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) | 文書と実装ファイルの役割、正本関係、配置を確認する。 |
| 4 | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | 方針、ポリシー、禁止事項、リリース判断を確認する。 |
| 5 | [`docs/DESIGN.md`](DESIGN.md) | 生成 HTML のデザイン関係を確認する。 |
| 6 | [`docs/ROADMAP.md`](ROADMAP.md) | 実装状態、実装可否、Phase、将来計画、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6 の追加仕様化機能参照、横断補足契約を確認する。 |
| 7 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) | 詳細仕様の入口、読み順、共通固定値、対応表、ソース配置を確認する。 |
| 8 | [`docs/details/*.md`](details/) | 対象 owner component の入出力、状態、処理順序、異常系、検証条件を確認する。 |

上記の順序は、文書整理、仕様改訂、実装、検証、PR 作成のすべてで共通とする。[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) は索引であり、仕様判断の正本ではない。

## 文書一覧

| ファイル | 役割 |
|---------|------|
| [`README.md`](../README.md) | 初見向け入口。概要、読む順番、主要ファイルを示す。 |
| [`docs/SPEC.md`](SPEC.md) | 方針責務・ポリシー責務として、方針、ポリシー、正本関係、禁止事項、リリース判断を示す。 |
| [`docs/ROADMAP.md`](ROADMAP.md) | 実装状態、実装可否、Phase、機能インベントリ、将来計画、昇格手順、追加仕様化機能参照、横断補足契約を示す。 |
| [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) | 詳細仕様の入口、読み順、共通固定値、対応表、リポジトリ内ソース配置を示す。 |
| [`docs/details/builder.md`](details/builder.md) | `builder` owner component の詳細仕様。 |
| [`docs/details/runner.md`](details/runner.md) | `runner` owner component の詳細仕様。 |
| [`docs/details/api.md`](details/api.md) | `api` owner component の詳細仕様。 |
| [`docs/details/admin.md`](details/admin.md) | `admin` owner component の詳細仕様。 |
| [`docs/details/sdk.md`](details/sdk.md) | `sdk` owner component の詳細仕様。 |
| [`docs/details/ui.md`](details/ui.md) | `ui` owner component の詳細仕様。 |
| [`docs/details/setup.md`](details/setup.md) | `setup` owner component の詳細仕様。 |
| [`docs/details/statefile.md`](details/statefile.md) | `statefile` owner component の詳細仕様。 |
| [`docs/details/archive.md`](details/archive.md) | `archive` owner component の詳細仕様。 |
| [`docs/details/commitstatus.md`](details/commitstatus.md) | `commitstatus` owner component の詳細仕様。 |
| [`docs/details/security.md`](details/security.md) | `security` owner component の詳細仕様。 |
| [`docs/details/fixture.md`](details/fixture.md) | `fixture` owner component の詳細仕様。 |
| [`docs/DESIGN.md`](DESIGN.md) | 生成静的 Web サイトのデザイン関係の正本。デザイン方針、レイアウト、色、タイポグラフィ、TOC、コードブロック等の視覚仕様を示す。 |
| [`AGENTS.md`](../AGENTS.md) | エージェント作業ルールブック。承認、仕様書管理、実装管理、Git 運用、文書整合の最上位ルール。 |
| [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) | 本索引。文書・実装ファイルの役割と所在を示す。仕様本文を定義しない。 |

## 実装ファイル索引

実装ファイルの所在と状態は、`仕様化済みコンポーネント` を正とする。実装ファイルが存在することだけで、仕様化済み、実装可、完了済みとは判断しない。

## 詳細仕様管理

詳細仕様は、責務 component 別の [`docs/details/*.md`](details/) を本文として管理する。[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) は、詳細仕様本文を集約せず、入口、読み順、共通固定値、対応表、リポジトリ内ソース配置だけを持つ。[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6 の追加仕様化機能参照と横断補足契約を正とする。

下表は、詳細仕様本文の配置先を示す索引である。各ファイルの責務境界、持つ内容、持たない内容、owner / collaborator の扱いは [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b および [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b.1 を正とする。

| ファイル | 状態 | 役割 |
|----------|------|------|
| [`docs/details/builder.md`](details/builder.md) | 責務化済み | `builder` owner の Markdown 変換、静的 Web サイト出力、HTML / CSS / JavaScript、theme component、builder fixture、builder owner 追加機能。 |
| [`docs/details/runner.md`](details/runner.md) | 責務化済み | `runner` owner の GitHub 監視、設定読取、状態ファイル更新呼び出し、pipeline、deploy、snapshot 作成トリガー、通知、runner fixture、runner owner 追加機能。 |
| [`docs/details/api.md`](details/api.md) | 責務化済み | `api` owner の HTTP 共通契約、endpoint、request / response、状態ファイル read/write 呼び出し境界、認証連携、api owner 追加機能。API fixture は [`docs/details/fixture.md`](details/fixture.md) §22-F。 |
| [`docs/details/admin.md`](details/admin.md) | 責務化済み | `admin` owner の管理 UI 静的ファイル配布物構成、配置、検証、HTTP 静的配信境界、A6 fixture 固定契約、setup/admin/release 連動 fixture 参照。 |
| [`docs/details/sdk.md`](details/sdk.md) | 責務化済み | `sdk` owner の SDK class、method、HTTP 対応、query / body 生成、error、stream、token 破棄。 |
| [`docs/details/ui.md`](details/ui.md) | 責務化済み | `ui` owner の DOM id、panel、操作、表示状態、SDK 呼び出し、秘密情報消去。 |
| [`docs/details/setup.md`](details/setup.md) | 責務化済み | `setup` owner のバイナリ配布、配置、systemd、セットアップ、アップデート、リリース成果物検証、admin 配布・rollback・secret 保持の連動 fixture。 |
| [`docs/details/statefile.md`](details/statefile.md) | 責務化済み | `statefile` owner の状態ファイル共通仕様、lock、atomic write、JSON Lines、破損時処理、状態読取 adapter、主要 schema。 |
| [`docs/details/archive.md`](details/archive.md) | 責務化済み | `archive` owner の build log archive、snapshot、download、delete、rollback、cleanup。 |
| [`docs/details/commitstatus.md`](details/commitstatus.md) | 責務化済み | `commitstatus` owner の GitHub Commit Status API payload、送信順、失敗時非反転、保存値、secret mask、検証条件。 |
| [`docs/details/security.md`](details/security.md) | 責務化済み | `security` owner の API token scope、API key、audit、session timeout、TOTP、rate limit、漏えい禁止、security 横断順序。 |
| [`docs/details/fixture.md`](details/fixture.md) | 責務化済み | fixture manifest、assertion、fake、testdata、expected / effects、Phase 3 / Phase 4 API fixture、api / sdk / ui / statefile cross fixture、setup/admin/release 連動 fixture、受け入れ fixture 共通契約、PR 証跡テンプレート、acceptance checklist、差し戻し条件、実装 PR 完了証跡。 |

詳細仕様を改訂する場合は、[`docs/ROADMAP.md`](ROADMAP.md) で実装状態と実装可否を確認し、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) の対応表から owner component を特定し、該当する [`docs/details/*.md`](details/) を本文として更新する。[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) は配置と役割の索引に限定し、仕様本文、詳細仕様本文、実装状態の最終判断を定義しない。

## 仕様化済みコンポーネント

[`docs/ROADMAP.md`](ROADMAP.md) と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) では、以下のコンポーネントも仕様化されている。

リポジトリ内ソース配置は、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.3 のディレクトリ構成と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0j を参照する。

下表は、現行リポジトリに存在する実装ファイルと、将来追加予定 path を区別して示す。`main.go`、`components/*.go`、`admin/` 配下の静的 UI ファイル、`testdata/<component>/` を現行配置として扱う。

標準配置図に含まれる未作成 path は、将来追加予定 path として扱い、該当 owner component が実装対象になった PR で追加する。標準配置図に含まれていることだけを理由に、未実装ファイル、将来追加予定 path、空ディレクトリ、placeholder を作成しない。

| パス | component | 状態 | 役割 |
|------|-----------|------|------|
| `main.go` | `-` | 実装済み | 起動入口。現時点では実行ファイル名に応じて `builder` または `runner` component を呼び出す。 |
| `components/builder.go` | `builder` | 実装済み | Go 版静的 Web サイトビルドスクリプト。`adlaire-ci-build` バイナリとして実行する。 |
| `components/builder_test.go` | `builder` | 実装済み | `components/builder.go` の Phase 1 fixture テスト。 |
| `go.mod` | `-` | 実装済み | Go module 定義。外部 module は追加しない。 |
| `testdata/builder/` | `builder` | 実装済み | Phase 1 の受け入れ fixture 入力。 |
| `components/runner.go` | `runner` | 実装済み | Go 版 CI ランナー。`adlaire-ci-runner` バイナリとして実行する。Phase 2 完了判定パスを対象とする。 |
| `components/runner_test.go` | `runner` | 実装済み | `components/runner.go` の Phase 2 fixture、hardening、完了判定パステスト。 |
| `components/api.go` | `api` | 実装済み | 管理 API サーバー。常駐 HTTP サーバーとして Adlaire CI の状態確認・操作 API を提供する。 |
| `admin/adlaire-ci-sdk.js` | `sdk` | 実装済み | 管理ツール用 JavaScript SDK。管理 API 通信を抽象化する。 |
| `admin/index.html` | `ui` | 実装済み | 標準管理ツール UI。SDK 経由で API と通信する。 |
| `components/mcp.go` | `mcp` | 将来追加予定 path | MCP サーバー。現時点では未作成であり、実装可能な詳細仕様を持たず、MCP 専用詳細仕様が新設されるまで実装対象ではない。 |

## 正本関係

| 判断対象 | 正本 |
|----------|------|
| 作業ルール、承認、Git 運用、文書整合 | [`AGENTS.md`](../AGENTS.md) |
| 方針、ポリシー、正本関係、禁止事項、リリース判断 | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 |
| 実装状態、実装可否、Phase、機能インベントリ、将来計画、昇格手順、追加仕様化機能参照、横断補足契約 | [`docs/ROADMAP.md`](ROADMAP.md) |
| 詳細仕様の入口、索引、共通固定値、実装前確認項目、検証マトリクス、詳細節対応表、リポジトリ内ソース配置 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) |
| 各 component の入出力、状態、処理順序、異常系、検証条件の本文 | owner component 別の [`docs/details/*.md`](details/) |
| 生成静的 Web サイトのデザイン関係 | [`docs/DESIGN.md`](DESIGN.md) |
| 文書・実装ファイルの参照先と役割 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) |
| 実装ファイル、テスト、fixture の所在 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) の「仕様化済みコンポーネント」 |

[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) は索引であり、仕様本文、デザイン本文、実装可否判断の正本ではない。仕様を変更する場合は、先に該当する正本仕様書を更新し、その内容に基づいて実装ファイルを更新する。

## 整合ガードレール

文書整合を保つため、以下を守る。

| Guardrail | 内容 |
|-----------|------|
| 廃止参照禁止 | リネーム済みファイル、廃止済み節、削除済み文書名を残さない。 |
| 二重正本禁止 | 同じ判断対象を複数文書で正本として定義しない。 |
| README 肥大化禁止 | README に詳細仕様、詳細ルール、状態表、ロードマップを戻さない。 |
| DETAIL 方針化禁止 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) と [`docs/details/*.md`](details/) に方針、ポリシー、ロードマップ状態の正本本文を持たせない。 |
| 索引本文化禁止 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) に機能仕様、API 仕様、状態 schema、処理本文を持たせない。 |
| 実装混入禁止 | 文書構造整理だけの PR で実装ファイル、testdata、fixture を変更しない。 |

## 整合メモ

現時点では、[`docs/ROADMAP.md`](ROADMAP.md)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、owner component 別の [`docs/details/*.md`](details/) に記載された一部機能は仕様化済みだが、リポジトリ内に実装コードが存在しない。

仕様化済みだが未実装の内容は、実装済み機能として扱わない。

現行実装実体は `main.go`、`components/*.go`、`admin/` 配下の静的 UI ファイル、`testdata/<component>/` である。`components/builder.go`、`components/runner.go`、`components/api.go`、`admin/adlaire-ci-sdk.js`、`admin/index.html` は [`docs/ROADMAP.md`](ROADMAP.md) の実装状態と本ファイルの `仕様化済みコンポーネント` に従って実装済みとして扱う。Go toolchain による `gofmt` と `go test` の検証対象は Go ファイルとする。

`build_spec.go`、`runner.go`、`build_spec_test.go`、`runner_test.go`、`testdata/build_spec/` は標準外配置であり、現行実装実体として扱わない。標準配置と現行実体の判断は、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.3、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0j、本ファイルの `仕様化済みコンポーネント` を同時に確認する。
