# Build-Scripts

Build-Scripts は、Markdown から静的 Web サイトを生成し、Go で実装する Adlaire CI の仕様と運用を管理するためのリポジトリです。

## 概要

本リポジトリでは、Markdown ファイルまたは Markdown ディレクトリから静的 Web サイトを生成する Go 版ビルドスクリプトと、GitHub 上の対象 Markdown 変更を検出してビルドを実行する Go 版 CI ランナーの仕様を管理します。

文書構造、正本参照先、実装ファイルの所在は [`docs/DOCUMENT_INDEX.md`](docs/DOCUMENT_INDEX.md) を参照します。

## 最初に読む文書

1. [`AGENTS.md`](AGENTS.md): 作業ルール、承認、Git 運用を確認する。
2. [`docs/SPEC.md`](docs/SPEC.md) 方針責務・ポリシー責務: 仕様、方針、ポリシー、禁止事項、実装着手可否を確認する。
3. [`docs/DOCUMENT_INDEX.md`](docs/DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務: 目的別の参照先を確認する。

## リリース形式

リリース形式の正本は [`docs/SPEC.md`](docs/SPEC.md) ポリシー責務 §1、セットアップ詳細は [`docs/details/setup.md`](docs/details/setup.md) §26 を参照します。

## 注意

作業ルールは [`AGENTS.md`](AGENTS.md) を基準とします。
