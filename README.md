# Build-Scripts

Build-Scripts は、Markdown から静的 Web サイトを生成し、Go で実装する Adlaire CI の仕様と運用を管理するためのリポジトリです。

## 概要

本リポジトリでは、Markdown ファイルまたは Markdown ディレクトリから静的 Web サイトを生成する Go 版ビルドスクリプトと、GitHub 上の対象 Markdown 変更を検出してビルドを実行する Go 版 CI ランナーの仕様を管理します。

方針・ポリシー・実装状態・実装可否・ロードマップの正本は `docs/SPEC.md` です。詳細仕様の入口、索引、共通固定値、実装前確認項目、検証マトリクス、Phase、詳細節対応表、リポジトリ内ソース配置、横断補足契約は `docs/DETAIL_INDEX.md` を参照します。各 owner component の入出力、状態、処理順序、異常系、検証条件の本文は `docs/details/*.md` を正本として扱います。

## 読む順番

1. `AGENTS.md`: 作業ルール、承認、Git 運用を確認する。
2. `docs/DOCUMENT_INDEX.md`: 文書と実装ファイルの役割、正本関係、配置を確認する。
3. `docs/SPEC.md`: 方針、ポリシー、実装状態、実装可否、ロードマップを確認する。
4. `docs/DETAIL_INDEX.md`: 詳細仕様の入口、読み順、共通固定値、対応表、ソース配置を確認する。
5. `docs/details/*.md`: 対象 owner component の入出力、状態、処理順序、異常系、検証条件を確認する。

## 主要ファイル

- `docs/SPEC.md`: Adlaire CI のマスター仕様書。
- `docs/DETAIL_INDEX.md`: Adlaire CI の Part 3 詳細仕様の入口。読み順、共通固定値、対応表、リポジトリ内ソース配置、横断補足契約を持つ。
- `docs/details/*.md`: owner component 別の入出力、状態、処理順序、異常系、検証条件の詳細仕様本文。
- `docs/DESIGN.md`: 生成静的 Web サイトのデザイン仕様を整理する補助文書。
- `docs/DOCUMENT_INDEX.md`: 文書・実装ファイルの役割を整理する索引。
- `AGENTS.md`: 本リポジトリにおけるエージェント作業ルール。

## 実装ファイル

標準ソース配置と実装状態は `docs/DOCUMENT_INDEX.md` の Specified Components を参照します。実装判断では、`docs/SPEC.md` の実装状態と `docs/DETAIL_INDEX.md` のリポジトリ内ソース配置を同時に確認します。

## リリース形式

Adlaire CI の標準リリース形式は、GitHub Releases に添付する OS/arch 別の Go 実行バイナリです。初期標準は Linux x86_64 とし、`adlaire-ci-build`、`adlaire-ci-runner`、管理 API 導入後の `adlaire-ci-api` を SHA-256 checksum 検証後に `/usr/local/bin/` へ配置します。

## 注意

仕様や挙動を変更する場合は、先に `docs/SPEC.md`、`docs/DETAIL_INDEX.md`、または該当する owner component 別の `docs/details/*.md` を改訂し、その内容に基づいて実装ファイルと関連文書を整合させます。
