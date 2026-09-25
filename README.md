# Build-Scripts

Build-Scripts は、Markdown から静的 Web サイトを生成し、Go で実装する Adlaire CI の仕様と運用を管理するためのリポジトリです。

## 概要

本リポジトリでは、Markdown ファイルまたは Markdown ディレクトリから静的 Web サイトを生成する Go 版ビルドスクリプト、GitHub 上の対象 Markdown 変更を検出してビルドを実行する Go 版 CI ランナー、管理 API、JavaScript SDK、管理 UI の仕様を管理します。

文書構造、正本参照先、実装ファイルの所在は [`docs/DOCUMENT_INDEX.md`](docs/DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務を参照します。

## 最初に読む文書

1. [`AGENTS.md`](AGENTS.md): 作業ルール、承認、Git 運用を確認する。
2. [`docs/SPEC.md`](docs/SPEC.md) 方針責務・ポリシー責務: 仕様、方針、ポリシー、禁止事項、実装着手可否を確認する。
3. [`docs/DOCUMENT_INDEX.md`](docs/DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務: 目的別の参照先を確認する。
4. [`docs/ROADMAP.md`](docs/ROADMAP.md) 状態・計画責務: 各 component と各機能の現在状態、Phase、機能インベントリ、将来計画を確認する。状態語彙と実装可否は [`docs/SPEC.md` ポリシー責務 §0a](docs/SPEC.md#0a-仕様成熟度ポリシー) を参照する。
5. [`docs/DETAIL_INDEX.md`](docs/DETAIL_INDEX.md) 詳細仕様入口責務: owner component 別の詳細仕様参照先と fixture 証跡責務の参照先を確認する。
6. [`docs/DESIGN.md`](docs/DESIGN.md) デザイン責務: 生成静的 Web サイトのデザイン関係を確認する。

## リリース形式

リリース形式の正本は [`docs/SPEC.md` ポリシー責務 §1](docs/SPEC.md#1-バージョン管理)、セットアップ詳細は [`docs/details/setup.md` 詳細本文責務 §26](docs/details/setup.md#26-セットアップアップデート手順) を参照します。

## 注意

作業ルールは [`AGENTS.md`](AGENTS.md) を基準とします。
