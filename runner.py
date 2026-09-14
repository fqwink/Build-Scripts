#!/usr/bin/env python3
"""
Adlaire CI Runner
仕様: build_spec_v3_spec.md Part 3 §10–§20
"""

VERSION = "1.0"

import base64
import json
import logging
import os
import subprocess
import sys
import urllib.error
import urllib.request

# ── 設定値（§12） ─────────────────────────────────────────
TOKEN_FILE   = "/opt/adlaire-builder/.github_token"   # GitHub PAT
OWNER        = "<GitHubオーナー名>"                    # リポジトリオーナー
REPO         = "<リポジトリ名>"                        # リポジトリ名
BRANCH       = "main"                                  # 対象ブランチ
TARGET_FILE  = "adlaire-db-spec.md"                   # 監視対象ファイル
SHA_FILE     = "/opt/adlaire-builder/.last_sha"       # blob SHA キャッシュ
SRC          = "/opt/adlaire-builder/repo/adlaire-db-spec.md"  # 書き出し先
BUILD_SCRIPT = "/opt/adlaire-builder/build_spec_v3.py"
LOG_LEVEL    = "INFO"
# ─────────────────────────────────────────────────────────


def setup_logging() -> None:
    logging.basicConfig(
        format="%(asctime)s [%(levelname)s] %(message)s",
        datefmt="%Y-%m-%dT%H:%M:%S",
        level=getattr(logging, LOG_LEVEL, logging.INFO),
        stream=sys.stdout,
    )


def read_token() -> str:
    """GitHub PAT を読み込む。不在・空の場合は起動失敗（§13）。"""
    if not os.path.exists(TOKEN_FILE):
        logging.error("Token file not found: %s", TOKEN_FILE)
        sys.exit(1)
    with open(TOKEN_FILE, encoding="utf-8") as f:
        token = f.read().strip()
    if not token:
        logging.error("Token file is empty: %s", TOKEN_FILE)
        sys.exit(1)
    return token


def github_get(path: str, token: str) -> dict:
    """GitHub API への GET リクエスト。HTTP エラー・通信エラー時は例外を送出。"""
    url = f"https://api.github.com{path}"
    req = urllib.request.Request(
        url,
        headers={
            "Authorization": f"Bearer {token}",
            "Accept": "application/vnd.github+json",
            "X-GitHub-Api-Version": "2022-11-28",
        },
    )
    logging.debug("GET %s", url)
    with urllib.request.urlopen(req) as resp:
        return json.load(resp)


def get_blob_sha(token: str) -> str:
    """Step 1: Git Trees API で TARGET_FILE の blob SHA を取得する（§13）。"""
    path = f"/repos/{OWNER}/{REPO}/git/trees/{BRANCH}?recursive=1"
    try:
        tree = github_get(path, token)
    except urllib.error.HTTPError as e:
        logging.error("Git Trees API HTTP error: %s %s", e.code, e.reason)
        sys.exit(1)
    except urllib.error.URLError as e:
        logging.error("Git Trees API failed: %s", e.reason)
        sys.exit(1)

    for item in tree.get("tree", []):
        if item.get("path") == TARGET_FILE:
            sha = item["sha"]
            logging.debug("blob SHA: %s", sha)
            return sha

    logging.error("Target file not found in tree: %s", TARGET_FILE)
    sys.exit(1)


def read_last_sha() -> str:
    """SHA キャッシュファイルから前回 SHA を読み込む。なければ空文字を返す（§13）。"""
    if not os.path.exists(SHA_FILE):
        return ""
    with open(SHA_FILE, encoding="utf-8") as f:
        return f.read().strip()


def fetch_blob(sha: str, token: str) -> bytes:
    """Step 2: Git Blobs API でファイル本文を取得し Base64 デコードして返す（§13）。"""
    path = f"/repos/{OWNER}/{REPO}/git/blobs/{sha}"
    try:
        blob = github_get(path, token)
    except urllib.error.HTTPError as e:
        logging.error("Git Blobs API HTTP error: %s %s", e.code, e.reason)
        sys.exit(1)
    except urllib.error.URLError as e:
        logging.error("Git Blobs API failed: %s", e.reason)
        sys.exit(1)

    encoding = blob.get("encoding", "base64")
    if encoding != "base64":
        logging.error("Unexpected blob encoding: %s", encoding)
        sys.exit(1)

    # GitHub API は content にニューラインを含む場合があるため除去してからデコード
    raw = blob.get("content", "").replace("\n", "")
    return base64.b64decode(raw)


def write_src(content: bytes) -> None:
    """デコード済み内容を SRC パスへ書き出す（§13）。"""
    os.makedirs(os.path.dirname(SRC), exist_ok=True)
    with open(SRC, "wb") as f:
        f.write(content)
    logging.debug("Written: %s (%d bytes)", SRC, len(content))


def run_pipeline() -> bool:
    """pipeline.sh を実行する。成功なら True、失敗なら False を返す（§13, §14）。"""
    pipeline = os.path.join(os.path.dirname(SRC), ".ci", "pipeline.sh")
    if not os.path.exists(pipeline):
        logging.error("pipeline.sh not found: %s", pipeline)
        return False
    logging.info("Build start: %s", pipeline)
    result = subprocess.run(
        ["bash", pipeline],
        cwd=os.path.dirname(SRC),
    )
    if result.returncode == 0:
        logging.info("Build succeeded (exit 0)")
        return True
    else:
        logging.error("Build failed (exit %d)", result.returncode)
        return False


def save_sha(sha: str) -> None:
    """SHA_FILE を新 SHA で更新する。ビルド成功後にのみ呼ぶ（§13, Part 2 §5）。"""
    with open(SHA_FILE, "w", encoding="utf-8") as f:
        f.write(sha)
    logging.info("SHA updated: %s", sha)


def main() -> None:
    setup_logging()
    logging.info("Adlaire CI Runner v%s start", VERSION)

    # トークン読み込み（不在の場合は起動失敗）
    token = read_token()

    # Step 1: blob SHA 取得
    current_sha = get_blob_sha(token)
    last_sha    = read_last_sha()

    # 変更なし → スキップして正常終了
    if current_sha == last_sha:
        logging.info("No changes detected (SHA: %s). Skip.", current_sha)
        return

    logging.info("Change detected: %s → %s", last_sha or "(none)", current_sha)

    # Step 2: blob 取得・SRC へ書き出し
    content = fetch_blob(current_sha, token)
    write_src(content)
    logging.info("Source written: %s (%d bytes)", SRC, len(content))

    # ビルド実行
    if not run_pipeline():
        # ビルド失敗時は SHA を更新しない（次回起動時に再試行）
        sys.exit(1)

    # ビルド成功後にのみ SHA を更新
    save_sha(current_sha)
    logging.info("Adlaire CI Runner done.")


if __name__ == "__main__":
    main()
