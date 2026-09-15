# Build-Scripts

Build-Scripts は、Adlaire DB 仕様書を HTML 化し、Adlaire CI の仕様と運用を管理するためのリポジトリです。

## 概要

本リポジトリでは、Adlaire DB 仕様書 Markdown から単一 HTML を生成するビルドスクリプトと、GitHub 上の対象 Markdown 変更を検出してビルドを実行する最小構成の CI ランナーを管理します。

機能仕様・運用仕様・API 仕様・CI 仕様の正本は `ADLAIRE_CI_SPEC.md` です。

## 主要ファイル

- `ADLAIRE_CI_SPEC.md`: Adlaire CI のマスター仕様書。
- `build_spec_v3.py`: Adlaire DB 仕様書 Markdown を単一 HTML へ変換するビルドスクリプト。
- `runner.py`: GitHub API で対象 Markdown の変更を検出し、ビルドパイプラインを実行する CI ランナー。
- `DESIGN.md`: 生成 HTML のデザイン仕様を整理する補助文書。
- `DOCUMENT_INDEX.md`: 文書・実装ファイルの役割を整理する索引。
- `AGENTS.md`: 本リポジトリにおけるエージェント作業ルール。

## 現在の実装範囲

現時点で実装済みの主要コンポーネントは `build_spec_v3.py` と `runner.py` です。

`api_server.py`、`adlaire-ci-sdk.js`、`admin/index.html`、`mcp_server.py` は仕様化済みまたは将来計画のコンポーネントであり、実装済みとして扱いません。

## 注意

仕様や挙動を変更する場合は、先に `ADLAIRE_CI_SPEC.md` を改訂し、その内容に基づいて実装ファイルと関連文書を整合させます。
