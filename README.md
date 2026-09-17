# Build-Scripts

Build-Scripts は、Markdown から静的 Web サイトを生成し、Go で実装する Adlaire CI の仕様と運用を管理するためのリポジトリです。

## 概要

本リポジトリでは、Markdown ファイルまたは Markdown ディレクトリから静的 Web サイトを生成する Go 版ビルドスクリプトと、GitHub 上の対象 Markdown 変更を検出してビルドを実行する Go 版 CI ランナーの仕様を管理します。

文書構造、読む順番、正本関係、実装ファイルの所在は `docs/DOCUMENT_INDEX.md` を参照します。

## 最初に読む文書

1. `AGENTS.md`: 作業ルール、承認、Git 運用を確認する。
2. `docs/DOCUMENT_INDEX.md`: 文書構造、読む順番、正本関係、実装ファイル所在を確認する。

## 仕様正本

仕様判断では、`docs/DOCUMENT_INDEX.md` の Reading Order と Source Of Truth に従って正本を確認します。README は入口であり、仕様本文、詳細仕様本文、実装状態、ロードマップ、API 仕様、状態 schema、検証 matrix を定義しません。

## 実装ファイル

標準ソース配置と実装状態は `docs/DOCUMENT_INDEX.md` の Specified Components を参照します。

## リリース形式

標準リリース形式は Go 実行バイナリです。詳細は `docs/SPEC.md` のリリース方針と `docs/details/setup.md` のセットアップ詳細を参照します。

## 注意

作業ルールは `AGENTS.md` を正とします。
