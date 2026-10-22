#!/usr/bin/env bash
set -euo pipefail
umask 077
export LC_ALL=C
export GOENV=off GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off DENO_NO_PACKAGE_JSON=1 DENO_NO_UPDATE_CHECK=1 DENO_NO_PROMPT=1

fail() {
  printf 'ci-setup:%s\n' "$1" >&2
  exit 1
}

[[ $# == 1 ]] || fail arguments
[[ "${RUNNER_TEMP:-}" == /* && -d "$RUNNER_TEMP" && ! -L "$RUNNER_TEMP" ]] || fail temp-boundary
[[ "$RUNNER_TEMP" != *[$'\n\r\\']* ]] || fail temp-boundary

if [[ "$1" == cleanup ]]; then
  root=${ADLAIRE_CI_TOOLS:-}
  if [[ -n "$root" ]]; then
    suffix=${root#"$RUNNER_TEMP"/adlaire-tools.}
    [[ "$root" == "$RUNNER_TEMP"/adlaire-tools.* && "$suffix" =~ ^[A-Za-z0-9]{8}$ && ! -L "$root" ]] || fail cleanup-boundary
    rm -rf -- "$root" 2>/dev/null || fail cleanup
  fi
  exit 0
fi

[[ "$1" == go || "$1" == deno ]] || fail runtime
system=$(uname -s) || fail platform
arch=$(uname -m) || fail platform
[[ "$system" == Linux && "$arch" == x86_64 ]] || fail platform
for file in "${GITHUB_PATH:-}" "${GITHUB_ENV:-}"; do
  [[ "$file" == /* && -f "$file" && ! -L "$file" && "$file" != *[$'\n\r\\']* ]] || fail output-boundary
done

root=$(mktemp -d "$RUNNER_TEMP/adlaire-tools.XXXXXXXX") || fail staging
keep=0
finish() {
  status=$?
  trap - EXIT
  if [[ "$keep" != 1 ]]; then
    if ! rm -rf -- "$root" 2>/dev/null; then
      printf 'ci-setup:cleanup\n' >&2
      status=1
    fi
  fi
  exit "$status"
}
trap finish EXIT

if [[ "$1" == go ]]; then
  version=1.22.12
  url=https://go.dev/dl/go1.22.12.linux-amd64.tar.gz
  expected=4fa4f869b0f7fc6bb1eb2660e74657fbf04cdd290b5aef905585c86051b34d43
  bin="$root/go/bin"
else
  version=2.9.7
  url=https://github.com/denoland/deno/releases/download/v2.9.7/deno-x86_64-unknown-linux-gnu.zip
  expected=c6527f24f4b16031d3ae4fa9f658d5f11534c8d84ce7dc8502420280919c3490
  bin="$root/deno"
fi
archive="$root/archive"
for attempt in 1 2 3; do
  rm -f -- "$archive" || fail cleanup
  curl_exit=0
  http_status=$(curl --disable --fail --silent --location --proto '=https' --proto-redir '=https' --tlsv1.2 --connect-timeout 15 --max-time 50 --retry 0 --max-filesize 200000000 --write-out '%{http_code}' --output "$archive" "$url" 2>/dev/null) || curl_exit=$?
  [[ "$http_status" =~ ^[0-9]{3}$ ]] || http_status=000
  printf 'ci-download:runtime=%s attempt=%s curl_exit=%s http_status=%s\n' "$1" "$attempt" "$curl_exit" "$http_status" >&2
  if [[ "$curl_exit" == 0 && "$http_status" == 200 ]]; then
    break
  fi
  retry=0
  case "$curl_exit" in
    5|6|7|18|28|35|52|55|56|92) retry=1 ;;
    22) case "$http_status" in 408|429|5[0-9][0-9]) retry=1 ;; esac ;;
  esac
  [[ "$retry" == 1 && "$attempt" != 3 ]] || fail download
  sleep 1 || fail download
done
digest_output=$(sha256sum -- "$archive" 2>/dev/null) || fail digest-read
[[ "${digest_output%% *}" == "$expected" ]] || fail digest

if [[ "$1" == go ]]; then
  tar -tzf "$archive" --quoting-style=literal >"$root/list" 2>/dev/null || fail archive-list
  [[ -s "$root/list" ]] || fail archive-entry
  while IFS= read -r entry; do
    case "$entry" in
      go/*) ;;
      *) fail archive-entry ;;
    esac
    case "/${entry%/}/" in
      *'/../'*|*'/./'*|*'//'*|*'\'*) fail archive-entry ;;
    esac
  done <"$root/list"
  grep -Fxq go/bin/go "$root/list" || fail archive-entry
  tar -tvzf "$archive" >"$root/types" 2>/dev/null || fail archive-list
  while IFS= read -r entry; do
    [[ "$entry" == -* || "$entry" == d* ]] || fail archive-type
  done <"$root/types"
  tar --extract --gzip --file "$archive" --directory "$root" --no-same-owner --no-same-permissions >/dev/null 2>&1 || fail extract
else
  entries=$(unzip -Z1 "$archive" 2>/dev/null) || fail archive-list
  [[ "$entries" == deno ]] || fail archive-entry
  mkdir "$bin" || fail staging
  unzip -q "$archive" -d "$bin" >/dev/null 2>&1 || fail extract
fi
[[ -f "$bin/$1" && ! -L "$bin/$1" ]] || fail executable
chmod 0755 "$bin/$1" || fail executable
if [[ "$1" == go ]]; then
  actual=$("$bin/go" version 2>/dev/null) || fail version
  [[ "$actual" == 'go version go1.22.12 linux/amd64' ]] || fail version
else
  actual=$("$bin/deno" --version 2>/dev/null) || fail version
  [[ "${actual%%$'\n'*}" == 'deno 2.9.7 (stable, release, x86_64-unknown-linux-gnu)' ]] || fail version
fi
rm -f -- "$archive" "$root/list" "$root/types" || fail cleanup
printf '%s\n' "$bin" >>"$GITHUB_PATH" || fail publish
printf 'ADLAIRE_CI_TOOLS=%s\nGOCACHE=%s/gocache\nGOMODCACHE=%s/gomodcache\nDENO_DIR=%s/denocache\nGOENV=off\nGOTOOLCHAIN=local\nGOPROXY=off\nGOSUMDB=off\nDENO_NO_PACKAGE_JSON=1\nDENO_NO_UPDATE_CHECK=1\nDENO_NO_PROMPT=1\n' "$root" "$root" "$root" "$root" >>"$GITHUB_ENV" || fail publish
printf '{"record_type":"development-actions-bootstrap","runtime":"%s","version":"%s","archive_sha256":"%s"}\n' "$1" "$version" "$expected"
keep=1
