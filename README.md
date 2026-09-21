# Build-Scripts

Build-Scripts は、Markdown から静的 Web サイトを生成し、Go で実装する Adlaire CI の仕様と運用を管理するためのリポジトリです。

## 概要

本リポジトリでは、Markdown ファイルまたは Markdown ディレクトリから静的 Web サイトを生成する Go 版ビルドスクリプトと、GitHub 上の対象 Markdown 変更を検出してビルドを実行する Go 版 CI ランナーの仕様を管理します。

文書構造、読む順番、正本参照先、実装ファイルの所在は [`docs/DOCUMENT_INDEX.md`](docs/DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務を参照します。

## 最初に読む文書

1. [`AGENTS.md`](AGENTS.md): 作業ルール、承認、Git 運用を確認する。
2. [`docs/SPEC.md`](docs/SPEC.md) 方針責務・ポリシー責務: 仕様、方針、ポリシー、禁止事項、実装着手可否を確認する。
3. [`docs/DOCUMENT_INDEX.md`](docs/DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務: 文書構造、正本参照先、実装ファイル所在を確認する。
4. [`docs/DESIGN.md`](docs/DESIGN.md) デザイン責務: 生成 HTML のデザイン関係を確認する。
5. [`docs/ROADMAP.md`](docs/ROADMAP.md) 状態・計画責務: 実装状態、Phase、将来計画を確認する。
6. [`docs/DETAIL_INDEX.md`](docs/DETAIL_INDEX.md) 詳細仕様入口責務: 詳細仕様入口と owner component 参照表を確認する。

## 参照先

仕様判断の本文は、対象に応じて [`docs/SPEC.md`](docs/SPEC.md)、[`docs/ROADMAP.md`](docs/ROADMAP.md)、[`docs/DETAIL_INDEX.md`](docs/DETAIL_INDEX.md)、owner component 別の [`docs/details/*.md`](docs/details/) 詳細本文責務を参照します。README は入口であり、仕様本文、詳細仕様本文、実装状態、ロードマップ、API 仕様、状態 schema、検証マトリクスを定義しません。

## 実装ファイル

標準ディレクトリ構成は [`docs/SPEC.md`](docs/SPEC.md) 方針責務 §4.3、実装ファイル所在は [`docs/DOCUMENT_INDEX.md`](docs/DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務、実装状態と Phase は [`docs/ROADMAP.md`](docs/ROADMAP.md) 状態・計画責務を参照します。

## リリース形式

標準リリース形式は Go 実行バイナリです。詳細は [`docs/SPEC.md`](docs/SPEC.md) ポリシー責務 §1 のリリース方針と [`docs/details/setup.md`](docs/details/setup.md) 詳細本文責務 §26 のセットアップ詳細を参照します。

## 注意

作業ルールは [`AGENTS.md`](AGENTS.md) を基準とします。
