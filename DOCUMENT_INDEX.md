# ASB Document Index

このファイルは、Adlaire-Static-Base（ASB）関連ドキュメントの参照先と役割を整理する索引である。

## Documents

| ファイル | 役割 |
|---------|------|
| `ASB-spec.md` | ASB の仕様正本。仕様判断の最上位基準。 |
| `ASB-spec.html` | `ASB-spec.md` から生成する HTML 版仕様書。仕様改訂の入力元ではない。 |
| `IMPLEMENTATION_TASKS.md` | `ASB-spec.md` に基づく実装フェーズ別タスクリスト。フェーズ番号、優先度、開発版バージョン、実装順序、実装タスク、フェーズ別完了条件を管理する。 |
| `AGENTS.md` | エージェント作業ルールブック。承認、Git運用、仕様書管理の最上位ルール。 |
| `README.md` | プロジェクト概要と主要文書への入口。仕様正本ではない。 |
| `templates/docs/` | 仕様節、実装フェーズ、README節、変更履歴の雛形。仕様正本ではない。 |
| `templates/rulebook/` | 新規・派生リポジトリ向けルールブック雛形。現行 `AGENTS.md` の代替ではない。 |

## Source Of Truth

仕様判断では `ASB-spec.md` を正とする。

`ASB-spec.html`、`IMPLEMENTATION_TASKS.md`、`README.md` に仕様上の不整合がある場合は、`ASB-spec.md` に合わせて更新する。

仕様改訂は `ASB-spec.md` のみで行う。

`ASB-spec.html` は `ASB-spec.md` に基づいて更新・再生成する。

`IMPLEMENTATION_TASKS.md` は `ASB-spec.md` で確定した仕様を実装フェーズ単位へ展開する文書であり、仕様を追加または変更する文書ではない。

P0、P1、P2 などのフェーズ詳細は `IMPLEMENTATION_TASKS.md` で管理し、`ASB-spec.md` へ記載しない。

`templates/` 配下のファイルは雛形であり、仕様判断、作業ルール判断、実装フェーズ判断の正本ではない。

`templates/rulebook/` 配下のファイルと現行 `AGENTS.md` が矛盾する場合は、現行 `AGENTS.md` を正とする。
