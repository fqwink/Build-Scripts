# Adlaire CI — Statefile 詳細仕様

[`docs/details/statefile.md`](statefile.md) は `statefile` owner component の詳細本文責務として、`statefile` が主本文として持つ実装契約だけを扱う。

owner / collaborator 境界管理は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0b.1](../DETAIL_INDEX.md#0b1-owner-component-別-owner-collaborator-境界管理) に従う。`statefile` owner component の主本文であり、collaborator component の仕様は読み書き境界、業務処理、表示、security、検証観点として参照する。fixture、expected、fake、実装検証証跡は [`docs/details/fixture.md`](fixture.md) fixture 証跡責務を参照する。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `statefile` |
| 実装主体 | 単独の Go artifact は持たない。runner が所有する状態の読み書きは [`components/runner.go`](../../components/runner.go)、API が所有する状態の読み書きは [`components/api.go`](../../components/api.go) に内包する。 |
| 持つ内容 | `statefile` owner が主本文として定義する状態ファイル共通仕様、lock、atomic write、JSON Lines、破損時処理、状態読取 adapter、主要 schema。 |
| 持たない内容 | API endpoint の request / response、runner の業務処理、SDK method 実装、UI 表示判断、setup / update 手順、release 生成・公開手順、fixture 証跡責務、個別 component の業務判断。 |

---

<a id="対象範囲"></a>
**対象範囲：**

| 範囲 | 内容 |
|------|------|
| [§22.0a](statefile.md#sec-22-0a) | 状態ファイル共通仕様、更新手順、schema 厳格化、状態読取 adapter。 |
| [§22.0c](statefile.md#sec-22-0c) | 主要状態ファイル schema。 |

<a id="sec-22-0a"></a>
**22.0a 状態ファイル共通仕様：**

`api` および拡張後 `runner` が読み書きする状態ファイルは、[`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の状態ファイル固定表の初期値、形式、更新責務に従う。表にない状態ファイルを追加してはならない。追加が必要な場合は、先に [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) へパス、形式、初期値、更新責務、破損時の扱いを追記する。

| パス | 形式 | 初期値 | 更新責務 | 破損時の扱い |
|------|------|--------|----------|--------------|
| `.admin_credentials` | JSON object | `--init-credentials` で生成 | `api` | 起動時に ERROR ログを出し、HTTP サーバーを起動しない。 |
| `.github_token` | UTF-8 text | 不在 | `setup` / `api` | 自動生成、自動修復、response への内容出力を行わない。読取不能は PAT または GitHub API caller を失敗させる。 |
| `.totp_secret` | JSON object | `{"enabled":false,"secret_base32":null,"confirmed_at":null,"last_accepted_step":null}` | `api` | 読み込み不能時は TOTP 有効 login caller へ認証失敗を返す。破損時は退避するが自動再生成で認証を弱めてはならない。 |
| `.audit_log` | JSON Lines | 空ファイル | `api` / `runner` | 読み込み可能な行のみ返し、壊れた行は無視する。追記不能時は対象操作を失敗扱いにする。 |
| `.api_rate_state` | JSON object | `{"windows":{}}` | `api` | `.api_rate_state.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、空 window で再生成する。 |
| `.server_config` | JSON object | `{}` | `api` | `.server_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、空 object で再生成する。 |
| `.notify_config` | JSON object | `{"webhooks":[],"channels":[],"on":[],"summary":{"enabled":false,"interval":"weekly","hour":9,"day_of_week":1},"email":{"enabled":false,"to":[],"on":[]}}` | `runner` / `api` | `.notify_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、初期値で再生成する。 |
| `.notify_log` | JSON Lines | 空ファイル | `runner` / `api` | 読み込み可能な行のみ使用し、壊れた行は ERROR ログへ記録して無視する。 |
| `.notify_pending` | JSON array | `[]` | `runner` | `.notify_pending.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、`[]` で再生成する。 |
| `.pending_transfers` | JSON array | `[]` | `runner` | `.pending_transfers.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、`[]` で再生成する。 |
| `.build_history` | JSON Lines | 空ファイル | `runner` / `api` | 読み込み可能な行のみ使用し、壊れた行は ERROR ログへ記録して無視する。 |
| `.build_logs/{id}.json` | JSON object | ビルドごとに新規作成 | `runner` / `api` | 対象 ID の read adapter は `ErrStateCorrupted` を返し、既存ファイルは上書きしない。 |
| `.build_logs/archive/{id}.json.gz` | gzip 圧縮 JSON object | archive 時に新規作成 | `archive` | runner / api owner は archive 処理を呼び出すか read adapter 経由で読むだけとする。展開不能または展開後 schema 不一致は破損として除外し、既存 archive を上書きしない。展開後 object は通常 build log schema と同型とする。 |
| `.build_logs/{build_id}_hook_{hook_id}.json` | JSON object | hook 実行ごとに新規作成 | `runner` | 個別破損 file は hook log 一覧から除外して固定 ERROR code を記録し、自動修復または上書きしない。 |
| `.build_lock` | text | 不在 | `runner` | 内容は `pid={pid}\nstarted_at={UTC_ISO8601}\n` とする。PID が存在しない場合は `readBuildLock()` が stale として返し、read-only caller は削除しない。runner owner の build / rollback / snapshot delete 排他 coordinator だけが開始前再読取で同じ stale 判定を確認後に削除できる。PID が存在する場合は実行中 conflict、形式不正または PID 判定不能は上書きせず conflict failure とする。 |
| `.last_sha` / `BranchTarget.SHAFile` | JSON object | `{"sha":""}` | `runner` | JSON 破損、object 以外、`sha` key 不在、`sha` 型不一致は当該 target の decode failure とし、成功時まで更新しない。 |
| `.sha_cache/{branch_safe}/{target_hash}.sha` | UTF-8 text | 不在 | `runner` | 40 または 64 文字 lowercase hex と末尾 LF 以外は cache miss として WARN を記録し、build 成功時だけ置換する。 |
| `.branch_config` | JSON object | 不在 | `runner` / `api` | `.branch_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、再生成せず `BRANCH_TARGETS` デフォルトへフォールバックする。 |
| `.build_state` | JSON object | `{"running":false,"current_build_id":null,"active_queue_entry":null,"queued":[],"last_started_at":null,"last_finished_at":null,"weekly_summary_last_sent_at":null,"weekly_summary_sent_date":null}` | `runner` / `api` | `.build_state.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、初期値で再生成する。 |
| `.build_status.json` | JSON object | `{"schema_version":1,"updated_at":null,"status":"none","running":false,"current_build_id":null,"last_build_id":null,"last_trigger":null,"last_target_status":null,"last_branch":null,"last_target_file":null,"last_blob_sha":null,"last_commit_sha":null,"last_started_at":null,"last_finished_at":null,"last_duration_seconds":null,"last_error":null,"last_deploy_at":null,"last_deploy_status":null,"pending_transfers_count":0,"notify_pending_count":0,"circuit_open":false,"circuit_consecutive_failures":0,"output_sha256":null,"size_warn":false}` | `runner` | `.build_status.json.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、初期値で再生成する。 |
| `.build_circuit_state` | JSON object | `{"open":false,"consecutive_failures":0,"opened_at":null,"last_failure_at":null,"last_error":null}` | `runner` / `api` | `.build_circuit_state.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、初期値で再生成する。 |
| `.local_watch_state.json` | JSON object | `{"files":{}}` | `runner` | `.local_watch_state.json.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、full build 後に再生成する。 |
| `.build_cache.json` | JSON object | `{"schema_version":1,"entries":{}}` | `builder` | 破損時は WARN を出し、cache miss として扱い、成功後に再生成する。 |
| `.build_cache/pages/` | directory | 空ディレクトリ | `builder` | entry 不一致または読み取り不能 file は miss とし、他 entry は継続使用する。 |
| `.dependency_manifest.json` | JSON object | 不在 | `builder` / `runner` | 不在、JSON 破損、schema 不一致は full build とし、成功後だけ再生成する。失敗時は既存 file を変更しない。 |
| `.snapshots/{build_id}/` | directory | snapshot 成功時に新規作成 | `archive` | runner owner は save / rollback を呼び出し、api owner は list / download / delete を呼び出すだけとする。`meta.json` または `site.tar.gz` 不在・不一致は当該 snapshot を破損扱いとし、自動修復、部分利用、上書きを行わない。保存形式は [`docs/details/archive.md` 詳細本文責務 §27.15](archive.md#sec-27-15) を正とする。 |
| `.remote_artifacts/{build_id}/` | temporary directory | remote build artifact 取得時だけ作成 | `runner` | 同じ build id の既存 directory は上書きせず remote build を失敗させる。成功・失敗後に削除を試み、削除失敗は WARN として残す。 |
| `.approval_queue` | JSON Lines | 空ファイル | `runner` / `api` | 読み込み可能な行のみ使用し、壊れた行は ERROR ログへ記録して無視する。 |
| `.build_trends.json` | JSON object | `{"schema_version":1,"samples":[],"summary":{"count":0,"avg_seconds":null,"median_seconds":null,"p95_seconds":null,"anomaly_count":0}}` | `runner` | `.build_trends.json.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、`.build_history` から再集計する。 |
| `.build_chain_config` | JSON object | `{"chains":[]}` | `api` | `.build_chain_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、chain 無効として通常 build のみ継続する。 |
| `.repo_config` | JSON object | 不在 | `runner` / `api` | read-only caller は退避、再生成、上書きを行わず、`ErrStateCorrupted` または `ErrStateReadFailed` を受け取る。不在時だけ `readRepoConfig()` が固定既定値を memory 上で返す。 |
| `.config_log` | JSON Lines | 空ファイル | `api` | 読み込み可能な行のみ返し、壊れた行は無視する。 |
| `.access_log` | JSON Lines | 空ファイル | `api` | 読み込み可能な行のみ返し、壊れた行は無視する。 |
| `.api_access_log` | JSON Lines | 空ファイル | `api` | 読み込み可能な行のみ返し、壊れた行は無視する。秘密情報は記録しない。 |
| `.webhook_secret` | text | 不在 | `api` | 読み込み不能時は Webhook caller へ secret read failure を返す。 |
| `.webhook_events.json` | JSON Lines | 空ファイル | `api` | 読み込み可能な行のみ返し、壊れた行は無視する。 |
| `.access_control` | JSON object | `{"allow":[]}` | `api` | 破損または読取不能時は自動修復せず、read adapter は caller へ判別可能な失敗を返す。API の拒否契約は [`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) を参照する。 |
| `.hooks` | JSON object | `{"hooks":[]}` | `api` | 初期値で再生成し、ERROR ログを記録する。 |
| `.maintenance` | JSON object | `{"enabled":false,"reason":null,"since":null}` | `api` | 破損または読取不能時は自動修復せず、read adapter は caller へ判別可能な失敗を返す。API の拒否判定は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e)、runner の拒否判定は [`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) を参照する。 |
| `.api_tokens` | JSON object | `{"tokens":[]}` | `api` | token caller へ破損失敗を返し、自動再生成しない。token 管理情報の消失による意図しない再許可を防ぐため、破損ファイルは上書きしない。 |
| `.alert_rules` | JSON object | `{"rules":[]}` | `api` | 初期値で再生成し、ERROR ログを記録する。 |
| `.tag_rules` | JSON object | `{"rules":[]}` | `api` | 初期値で再生成し、ERROR ログを記録する。 |
| `.pipeline_config` | JSON object | `{"extra_args":[],"env":{},"inline_yaml":null}` | `api` | 初期値で再生成し、ERROR ログを記録する。 |
| `.notes` | UTF-8 text | 空文字列 | `api` | 読み込み不能時は read failure を返し、自動上書きしない。 |
| `.smtp_config` | JSON object | `{"host":null,"port":587,"user":null,"tls":true,"from":null,"to":[],"on":[],"enabled":false}` | `api` | 初期値で再生成し、ERROR ログを記録する。 |
| `.smtp_secret` | text | 不在 | `api` | 読み込み不能時は SMTP caller へ secret read failure を返す。 |
| `.dashboard_layout` | JSON object | `{"widgets":["status","stats","schedule","alerts","disk","rate_limit","snapshots","maintenance","queue"]}` | `api` | 初期値で再生成し、ERROR ログを記録する。 |

状態ファイル固定表の「破損時の扱い」に記載した退避、再生成、初期値復旧、default fallback は、更新 caller または起動時整合性回復が明示的に回復処理を呼び出した場合だけ実行する。read-only caller は必ず [状態読取 adapter 固定契約](#statefile-read-adapter-contract) を優先し、作成、削除、退避、再生成、書き戻しを行わず、破損または読取不能を caller へ返す。

`.build_logs/archive/` は gzip 圧縮済み build log の保存先ディレクトリである。初期値は空ディレクトリとし、runner または `POST /api/logs/archive` から呼び出された archive owner が書き込み前に存在確認し、不在の場合だけ作成する。runner と api owner は directory または gzip file を直接作成しない。圧縮済みファイル名は `{id}.json.gz` 固定とし、通常 `.build_logs/{id}.json` と同じ build id を表す。

JSON Lines ファイルは、1 行につき 1 JSON object とする。追記時は末尾に改行を必ず付ける。runtime 状態ファイルは内容分類にかかわらず mode `0600`、statefile 責務が作成する runtime 状態ディレクトリは mode `0700` を必須とする。公開成果物ディレクトリ、インストール先 root、admin 静的配布物の mode は statefile 責務に含めない。

<a id="statefile-update-procedure"></a>
**状態ファイル更新手順：**

状態ファイル更新手順は、既存 target を置換できる通常 mode と、target 不在時だけ作成できる create-only mode を持つ。caller は owner 詳細本文で mode を固定し、省略時は通常 mode とする。statefile owner が caller の業務条件から mode を推測してはならない。

1. 対象ファイルの `{name}.lock` を `O_CREATE|O_EXCL|O_WRONLY`、mode `0600` で作成する。
2. ロック取得に失敗した場合は 100ms 間隔で最大 10 秒待つ。
3. 通常 mode は現在値を読み込み、schema と入力値を検証する。create-only mode は `Lstat` で target の存在だけを再確認し、通常 file、directory、symlink、その他の file type のいずれでも存在する場合は内容を読まず `ErrStateAlreadyExists` とする。target 不在時だけ入力値を検証して手順 4 へ進む。
4. 状態ファイル固定表が対象に指定する JSON、JSON Lines、または UTF-8 text の更新後 payload を、対象と同じ親ディレクトリの `{name}.tmp.{pid}` に `O_CREATE|O_EXCL|O_WRONLY`、mode `0600`、UTF-8 / LF で全 byte 書き出す。JSON object / array は末尾 LF 1 個、JSON Lines は各 object の末尾 LF 1 個、UTF-8 text は末尾 LF 1 個だけを持つ。既存 tmp を truncate または再利用しない。
5. tmp file に `Sync` を実行し、mode `0600` を確定してから close する。`Sync`、chmod、close のいずれかが失敗した場合は rename しない。
6. `os.Rename(tmp, target)` で置換する。
7. 親ディレクトリを open して `Sync` し、rename された directory entry を永続化する。
8. ロックファイルを削除する。

`os.Rename` が成功する前に手順 3〜6 が失敗した場合は target を変更せず、当該 write 呼び出しが作成した tmp がある場合は tmp、当該 write 呼び出しが取得した lock がある場合は lock の順で、それぞれ `os.Remove` を 1 回だけ実行する。削除成功と `os.IsNotExist` は cleanup 成功とする。それ以外の tmp 削除失敗は `STATE_TMP_CLEANUP_FAILED`、lock 削除失敗は `STATE_LOCK_CLEANUP_FAILED` を ERROR で server log に記録し、code と対象 basename 以外を記録しない。cleanup は再試行せず、cleanup 失敗後も最初の write failure を呼び出し元へ返し、残存 tmp または lock を別 process のものと推測して削除しない。lock 取得前の失敗では tmp と lock の cleanup を行わない。

create-only mode で `ErrStateAlreadyExists` が確定した場合は tmp を作成せず、この呼び出しが取得した lock を 1 回削除する。削除成功または `os.IsNotExist` の場合は `ErrStateAlreadyExists` を返す。lock 削除失敗は `STATE_LOCK_CLEANUP_FAILED` を ERROR で記録し、target を変更せず cleanup failure を返す。caller は `ErrStateAlreadyExists` と cleanup failure を同じ公開結果へ丸めてはならない。

`os.Rename` 成功後に手順 7 または 8 が失敗した場合は、rename 済みの新しい target を維持し、server log に `STATE_WRITE_AFTER_RENAME_FAILED` を ERROR で記録し、write caller へ post-rename partial write failure を返す。statefile owner は caller 固有の `.config_log`、`.audit_log`、その他の業務 log を直接追記しない。write caller は自身の詳細本文責務に明記された場合だけ partial failure を記録し、失敗した atomic write を同じ payload で再試行しない。owner component の詳細本文責務が lifecycle 最終化または補償の発動条件、対象、固定値、書込順、最大回数を明記する場合に限り、最初の業務 write とは別の補償 atomic write をその契約どおり実行できる。この補償は元の payload の retry として扱わず、statefile owner が発動判断、補償値、順序、回数を補完してはならない。複数ファイル更新 caller は、呼び出し元が定義する Write 列順にこの手順を実行し、途中失敗時は未処理ファイルを書き込まない。API 固有の Write 列順は [`docs/details/api.md` 詳細本文責務 §22.0d](api.md#sec-22-0d) 以降を参照する。既に書き込んだファイルの暗黙のロールバックは行わない。

<a id="statefile-schema-strictness-contract"></a>
**状態ファイル schema 厳格化契約：**

状態ファイルの読込、正規化、保存は以下に固定する。[`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) または個別機能節で例外を明記していない限り、実装者判断で旧形式、未知 key、null、欠落配列を成功扱いにしてはならない。

| 対象 | 読込時 | 保存時 | 失敗時 |
|------|--------|--------|--------|
| 未知 key | JSON object に schema 未定義 key がある場合は破損扱いとする。例外は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) で明記した旧形式正規化だけ。 | 未知 key を保存しない。既存未知 key を黙って削除して保存しない。 | read adapter は `ErrStateCorrupted` を返す。write 呼び出しは target を変更しない。 |
| 必須 key 不足 | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の対象 schema に「欠落時に適用する既定値」と「保存するか読み取り時だけか」が明記されていない場合は破損扱いとする。 | 必須 key はすべて明示保存する。 | 初期値再生成が [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) 表で指定されたファイルだけ再生成する。 |
| `null` | 型欄が `string/null`、`object/null`、`integer/null` 等で明示した key だけ許可する。 | nullable でない key に `null` を保存しない。 | validation error または破損扱い。 |
| 配列 | `[]` を既定値とする key は read adapter の戻り値で空配列を返す。 | 保存呼び出しは配列 key を省略せず、空の場合も `[]` を明示する。 | 型不一致は caller 固有の validation error または状態ファイル破損扱い。 |
| 数値 | 整数 key は JSON number の整数だけ許可する。小数、指数表記由来の非整数、文字列数値は拒否する。 | 整数は JSON number として保存する。 | 書込入力は caller 固有の validation error、状態ファイル読込は破損扱い。 |
| 時刻 | UTC ISO 8601 秒精度 `YYYY-MM-DDTHH:MM:SSZ` だけ許可する。 | 保存前に UTC 秒精度へ丸める。ミリ秒、ナノ秒、local timezone、UTC 以外の offset 付き文字列を保存しない。 | 書込入力は caller 固有の validation error、状態ファイル読込は破損扱い。 |
| mode | runtime 状態ファイルと lock file は `0600`、statefile 責務が作成する runtime 状態 directory は `0700` に固定する。 | tmp と lock の chmod は rename 前に完了する。chmod 失敗時は rename せず成功扱いにしない。 | chmod failure を返し、旧 target を byte 単位で維持し、tmp と lock は状態ファイル更新手順の rename 前 cleanup 固定契約で削除する。 |
| 改行 | text / JSON / JSON Lines は LF で保存する。JSON object / array ファイルは末尾 LF 1 個を付ける。 | CRLF、BOM、末尾余分空白を新規保存しない。 | CRLF を含む保存入力は validation error、CRLF を含む既存 runtime 状態ファイルは破損扱いとする。 |

旧 schema からの正規化は、[`docs/details/statefile.md`](statefile.md) 詳細本文責務に「旧 key」「変換後 key」「削除する key」「保存するか読み取り時だけか」を明記した場合だけ実装する。明記がない旧形式は破損扱いとし、黙って推測変換してはならない。

以降の schema 表で `ISO 8601`、`UTC ISO 8601`、`UTC ISO 8601 秒精度` と略記する時刻文字列は、すべて [状態ファイル schema 厳格化契約](#statefile-schema-strictness-contract) の `YYYY-MM-DDTHH:MM:SSZ` 固定契約を意味する。個別 schema が `YYYY-MM-DD` を明記する日付値はこの時刻契約の対象外とする。

<a id="statefile-read-adapter-contract"></a>
**状態読取 adapter 固定契約：**

[`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の状態読取 adapter 固定契約は、statefile owner component が提供する状態読取 adapter 名、読取対象、正常戻り値、不在時、破損時 / 読込不能時の固定契約である。API endpoint ごとの読取順、response 算出、HTTP status は [`docs/details/api.md` 詳細本文責務 §22.0c.1](api.md#sec-22-0c-1) と [`docs/details/api.md` 詳細本文責務 §22.0d](api.md#sec-22-0d) 以降を参照する。各 component は同じ状態ファイルを endpoint ごとに別ロジックで直接 parse してはならない。

| Adapter | 読取対象 | 正常戻り値 | 不在時 | 破損時 / 読込不能時 |
|---------|----------|------------|--------|---------------------|
| `readBuildStatus()` | `.build_status.json` | `BuildStatus` | `(nil, false, nil)` を返し、fallback 判定へ渡す。 | `(nil, true, ErrStateCorrupted)` または `ErrStateReadFailed`。 |
| `readBuildState()` | `.build_state` | `BuildState` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の初期値を返す。 | `ErrStateCorrupted` または `ErrStateReadFailed`。 |
| `readServerConfig()` | `.server_config` | 既定値を merge し、全 top-level key を持つ `ServerConfig` | 全 key が既定値の `ServerConfig` を返す。 | `ErrStateCorrupted` または `ErrStateReadFailed`。read-only caller は退避、再生成、書き戻しを行わない。 |
| `readBuildHistory()` | `.build_history` | 物理行順を保持した、schema-valid かつ id 重複除外済みの `[]BuildHistoryEntry` | 空配列を返す。 | 行単位破損と同一 id の後続行は除外し、ファイル読込不能だけ `ErrStateReadFailed`。 |
| `readBuildLog(id)` | `.build_logs/{id}.json`、`.build_logs/archive/{id}.json.gz` | `BuildLog` | 通常ログ不在時は archive を読む。両方不在は `ErrNotFound`。 | 対象 ID の通常ログまたは archive が破損している場合は `ErrStateCorrupted`。 |
| `readLatestBuildLogs(n,q)` | `.build_logs/`、`.build_logs/archive/` | `[]LogLine` | 空配列を返す。 | 個別ログ破損は除外し、`LOG_SKIP_CORRUPT` を server log へ記録する。ディレクトリ読込不能は `ErrStateReadFailed`。 |
| `readBranchConfig()` | `.branch_config` | 旧形式正規化例外を適用した `[]BranchTarget` | `(nil, false, nil)` を返し、caller の default 判定へ渡す。 | `(nil, true, ErrStateCorrupted)` または `ErrStateReadFailed`。read-only caller は退避、削除、default 書込みを行わない。 |
| `readRepoConfig()` | `.repo_config` | `RepoConfig` | `{owner:"fqwink",repo:"Build-Scripts",updated_at:null}` を memory 上で返し、file を作成しない。 | `ErrStateCorrupted` または `ErrStateReadFailed`。保存値の一部だけを既定値へ差し替えない。 |
| `readPendingTransfers()` | `.pending_transfers` | `[]PendingTransfer` | 空配列を返す。 | `ErrStateCorrupted` または `ErrStateReadFailed`。 |
| `readNotifyPending()` | `.notify_pending` | `[]NotificationPending` | 空配列を返す。 | `ErrStateCorrupted` または `ErrStateReadFailed`。 |
| `readCircuitState()` | `.build_circuit_state` | `BuildCircuitState` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の初期値を返す。 | `ErrStateCorrupted` または `ErrStateReadFailed`。 |
| `readMaintenance()` | `.maintenance` | `MaintenanceState` | `{enabled:false,reason:null,since:null}` を memory 上で返し、file を作成しない。 | `ErrStateCorrupted` または `ErrStateReadFailed`。無効と推測せず caller へ失敗を返す。 |
| `readAccessControl()` | `.access_control` | `AccessControlState` | `{allow:[]}` を memory 上で返し、file を作成しない。 | `ErrStateCorrupted` または `ErrStateReadFailed`。全許可と推測せず caller へ失敗を返す。 |
| `readBuildLock()` | `.build_lock` | `BuildLockState` | `running=false` を返す。 | 形式不正、PID 判定不能、OS 判定失敗は `running=true, stale=false, valid=false` として返し、呼び出し元は conflict failure として扱う。 |

read-only 呼び出しでは、[`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の固定表の adapter を使用し、状態ファイルの作成、削除、退避、chmod、正規化、再生成、破損行の除去書き戻しを行ってはならない。API endpoint 固有の適用条件は [`docs/details/api.md` 詳細本文責務 §22.0c.1](api.md#sec-22-0c-1) を参照する。`{name}.lock` を検出しても、`.build_lock` 以外の lock file は待機条件やエラー条件にせず、rename 済み target をそのまま読む。write 呼び出しは [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の状態ファイル更新手順に従う。

`ErrStateCorrupted` は JSON parse 失敗、schema_version 不一致、必須 key 不足、型不一致、列挙値不一致、UTC 時刻形式不一致のいずれかで返す。`ErrStateReadFailed` は permission denied、通常ファイルではない path、gzip 読込失敗、I/O error で返す。API の公開応答は [`docs/details/api.md` 詳細本文責務 §22.0c.1](api.md#sec-22-0c-1) と [`docs/details/api.md` 詳細本文責務 §22.0d](api.md#sec-22-0d) 以降を参照する。statefile adapter の error は、path、Go error、ファイル内容を呼び出し元へ公開する response 値として含めない。

JSON Lines adapter は空行、JSON parse 失敗、JSON object 以外、必須 key 不足、型不一致の行を壊れた行として除外する。除外後に sort、filter、paging、`total`、`pages` を算出する。壊れた行の存在は呼び出し元の公開値に含めず、server log に固定コード、path、1 始まりの line number だけを記録する。

`readBuildHistory()` は schema-valid な行を物理的な file 先頭から読み、最初に出現した id の行だけを採用する。同じ id の後続行は `BUILD_HISTORY_DUPLICATE_ID`、path、1 始まりの line number、id だけを server log へ記録して除外し、file を書き換えない。`.build_history` で対象機能契約に並び順がない「最新」「直近」「一つ前」を判定する場合は、`finished_at` 降順、同時刻は `id` 降順とし、「一つ前」はその並びで直後の行とする。status summary 対象行は、[runner 結果値 schema](#runner-result-schema) で `status summary` が `許可` の詳細結果だけとし、`approval_rejected`、`approval_expired`、`skipped_dependency_failed` を含む history-only 行は除外する。対象機能契約が別の sort、filter、集計対象を明記する場合は、その契約を適用する。

`.build_lock` の PID が存在しない場合、`readBuildLock()` は `running=false, stale=true, valid=true` を返す。read-only caller は stale lock を削除しない。runner owner の build / rollback / snapshot delete 排他 coordinator は開始直前に `.build_lock` を再読込し、同じ stale 判定なら `.build_lock` だけを削除してから新規 lock を作成する。削除失敗時は conflict failure とし、`.build_state` を変更しない。API owner、archive owner、read-only caller は `.build_lock` を直接作成、削除、上書きしない。

<a id="sec-22-0c"></a>
**22.0c 主要状態ファイル schema：**

[`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の schema は、API 実装、SDK 型、標準管理ツール表示、バックアップ/リストアの基準である。[`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) に定義したキー以外を保存してはならない。追加キーを追加する場合は、型、既定値、読み書き API、既存データの扱いを [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) へ追記してから実装する。

<a id="server-config-schema"></a>
**`.server_config` schema：**

| キー | 型 | 既定値 | 許容値 | 読み書き API | 説明 |
|------|----|--------|--------|--------------|------|
| `log_max_lines` | integer | `500` | 1〜10000 | `GET/POST /api/config` | `GET /api/logs` が返す最大行数。 |
| `history_max_count` | integer | `100` | 1〜10000 | `GET/POST /api/config` | `.build_history` の通常表示上限。削除処理の上限ではない。 |
| `build_timeout_seconds` | integer | `300` | 1〜86400 | `GET/POST /api/config` | 手動/自動ビルドのタイムアウト秒数。 |
| `log_retention_days` | integer | `30` | 0〜3650 | `GET/POST /api/config`, `POST /api/logs/cleanup` | `0` は自動削除なし。 |
| `log_level` | string | `"INFO"` | `"INFO"` / `"DEBUG"` / `"WARNING"` / `"ERROR"` | `GET/POST /api/config`, `POST /api/log-level` | `api` のランタイムログレベル。 |
| `pat_expires_at` | string/null | `null` | `YYYY-MM-DD` または `null` | `GET/POST /api/config` | PAT 期限表示・診断用。 |
| `snapshots_keep` | integer | `5` | 0〜100 | `GET/POST /api/config` | `0` はスナップショット保存無効。 |
| `queue_max_size` | integer | `3` | 0〜100 | `GET/POST /api/config`, `GET /api/queue` | running 中に保持できる待機 entry 上限。`0` は running 中の追加待機を無効にする。idle 時の最初の dispatch entry 1 件は別途許可する。 |
| `build_retry_max` | integer | `0` | 0〜10 | `GET/POST /api/config` | ビルド失敗時の自動リトライ最大回数。`0` は無効。 |
| `build_retry_base_seconds` | integer | `5` | 1〜3600 | `GET/POST /api/config` | 自動リトライ backoff 基底秒数。待機秒数は `base * attempt` とする。 |
| `commit_status_enabled` | boolean | `false` | `true` / `false` | `GET/POST /api/config` | GitHub Commit Status API 送信の有効/無効。 |
| `commit_status_context` | string | `"Adlaire CI"` | 1〜100 文字 | `GET/POST /api/config` | GitHub commit status の `context`。 |
| `commit_status_target_url` | string/null | `null` | `http://` または `https://` の URL、または `null` | `GET/POST /api/config` | Commit Status の `target_url`。`null` の場合は送信 payload から省略する。 |
| `log_archive_after_days` | integer | `0` | 0〜3650 | `GET/POST /api/config`, `POST /api/logs/archive` | `0` は archive 無効。指定日数より古い通常 build log を gzip 圧縮する。 |
| `build_trend_keep_count` | integer | `1000` | 10〜10000 | `GET/POST /api/config` | `.build_trends.json` に保持する trend sample 件数。 |
| `duration_anomaly` | object | `{"enabled":false,"min_samples":20,"avg_multiplier":2.0,"p95_multiplier":1.5}` | 次の DurationAnomaly object | `GET/POST /api/config` | build 所要時間異常検知の設定。判定処理は [`docs/details/runner.md` 詳細本文責務 §27.38](runner.md#sec-27-38) を参照する。 |
| `watch_mode` | string | `"github"` | `"github"` / `"local"` | `GET/POST /api/config` | 監視 source。`"local"` の処理契約は [`docs/details/runner.md` 詳細本文責務 §27.23](runner.md#sec-27-23) を参照する。 |
| `tag_filter` | object | `{"enabled":false,"patterns":[]}` | 次の TagFilter object | `GET/POST /api/config` | tag 付き commit の build 条件。 |
| `build_cache_enabled` | boolean | `false` | `true` / `false` | `GET/POST /api/config` | runner が builder へ cache directory を渡す条件。 |
| `deploy_parallelism` | integer | `1` | 1〜16 | `GET/POST /api/config` | deploy worker の最大並列数。 |
| `remote_build` | object | `{"enabled":false,"host":null,"user":null,"work_dir":null,"command_args":[],"artifact_path":null}` | 次の RemoteBuildConfig object | `GET/POST /api/config` | remote build の接続先と command。 |
| `approval_timeout_seconds` | integer | `86400` | 60〜2592000 | `GET/POST /api/config` | approval pending の有効期間。 |
| `force_build_interval_hours` | integer | `0` | 0〜8760 | `POST /api/schedule/force-interval`, `GET /api/schedule` | `0` は強制再ビルド無効。 |
| `build_cooldown_seconds` | integer | `0` | 0〜86400 | `POST /api/schedule/cooldown`, `GET /api/schedule` | `0` はクールダウン無効。 |
| `schedule_interval_seconds` | integer | `300` | 30〜86400 | `POST /api/schedule/interval`, `GET /api/schedule` | systemd timer 更新値。 |
| `schedule_paused` | boolean | `false` | `true` / `false` | `POST /api/schedule/pause`, `POST /api/schedule/resume`, `GET /api/schedule` | 自動ポーリング停止状態。 |
| `allowed_hours` | object/null | `null` | `{"from":0〜23,"to":0〜23}` または `null` | `POST /api/schedule/allowed-hours`, `GET /api/schedule` | UTC の自動ビルド許可時間帯。 |
| `session_timeout_seconds` | integer | `28800` | 300〜2592000 | `GET/POST /api/config` | 新規 session の有効期限秒数。既存 session の `expires_at` は変更しない。 |
| `api_rate_limit` | object | `{"enabled":true,"groups":{"login":{"window_seconds":60,"max_requests":10},"read":{"window_seconds":60,"max_requests":600},"trigger":{"window_seconds":60,"max_requests":60},"operate":{"window_seconds":60,"max_requests":120},"config":{"window_seconds":60,"max_requests":60},"admin":{"window_seconds":60,"max_requests":60}}}` | 次の ApiRateLimitPolicy object | `GET/POST /api/api-rate-limit` | API rate limit の endpoint group 別固定窓設定。判定処理は [`docs/details/security.md` 詳細本文責務 §27.47](security.md#sec-27-47) を参照する。 |

`.server_config` の on-disk object は [`.server_config` schema](#server-config-schema) の top-level key だけを 0 件以上持つ sparse object を許可する。`{}` は有効であり、省略した top-level key には同 schema の既定値を読取時の memory 上だけで適用する。これは [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の必須 key 不足契約に対する `.server_config` 限定の例外である。未知の top-level key は破損とする。存在する key は型、範囲、列挙値を検証し、存在する nested object は各 object 表の必須 key を過不足なく持たなければならない。

`readServerConfig()` は sparse object と既定値を deep copy した memory 上の object へ top-level key 単位で merge し、全 top-level key を持つ `ServerConfig` を返す。nested object の部分 merge は行わない。読取時は `.server_config` の作成、完全形への書き戻し、key 順の変更、mtime の変更を行わない。既存値と更新値の no-op 判定は、どちらも `readServerConfig()` と同じ規則で正規化した全 key の値で比較する。

`POST /api/config`、`POST /api/log-level`、`POST /api/api-rate-limit`、および `POST /api/schedule/*` が no-op でない `.server_config` 更新を確定した場合は、更新後の正規化値の全 top-level key を省略せず atomic write する。起動時整合性回復が初期値 `{}` を再生成する場合と、backup / restore が sparse object を復元する場合だけはこの完全形保存の例外とし、backup / restore の入出力と no-op 判定は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の backup / restore 固定契約を参照する。

TagFilter object:

| キー | 型 | 必須 | 許容値 / 説明 |
|------|----|------|---------------|
| `enabled` | boolean | 必須 | `true` / `false`。 |
| `patterns` | string[] | 必須 | 0〜100 件。各値は 1〜255 文字の完全一致文字列、末尾 `*` だけを使う prefix pattern、または `*` 単体。重複は禁止する。 |

`patterns` は入力順を保持して保存する。空文字、前後空白、制御文字、中間 `*`、複数の `*`、正規表現構文は validation failure とする。`enabled=true` かつ空配列の場合は任意の tag 1 件以上を一致条件とする。build 判定は [`docs/details/runner.md` 詳細本文責務 §27.24](runner.md#sec-27-24) を参照する。

RemoteBuildConfig object:

| キー | 型 | 必須 | 許容値 / 説明 |
|------|----|------|---------------|
| `enabled` | boolean | 必須 | `true` / `false`。 |
| `host` | string/null | 必須 | `^[A-Za-z0-9._-]{1,255}$` または `null`。 |
| `user` | string/null | 必須 | `^[A-Za-z0-9._-]{1,64}$` または `null`。 |
| `work_dir` | string/null | 必須 | remote host 上の正規化済み絶対 POSIX path または `null`。 |
| `command_args` | string[] | 必須 | 0〜50 件。各値は 1〜4096 bytes の有効な UTF-8 文字列とし、NUL、CR、LF、その他の制御文字を禁止する。 |
| `artifact_path` | string/null | 必須 | remote host 上の正規化済み絶対 POSIX path または `null`。 |

`work_dir` と `artifact_path` は `/` で開始し、有効な UTF-8 文字列とし、NUL、CR、LF、その他の制御文字、および正規化前の `..` segment を禁止する。入力を Go `path.Clean` で正規化し、root 以外の末尾 `/` を除いた値を保存する。`artifact_path="/"` は file を指さないため validation failure とする。`command_args[0]` は `/` で開始する正規化済み絶対 POSIX path、または `^[A-Za-z0-9._-]{1,256}$` に一致する command 名とする。絶対 path の場合は `work_dir` と同じ path 検証と正規化を適用し、`/` 単体を禁止する。

`enabled=false` では `host`、`user`、`work_dir`、`artifact_path` は `null`、`command_args` は空配列を許可する。`enabled=true` ではこれら 4 文字列を非 `null`、`command_args` を 1 件以上とする。SSH への実行引数化、artifact 取得、検証、cleanup は [`docs/details/runner.md` 詳細本文責務 §27.29](runner.md#sec-27-29) を参照する。

DurationAnomaly object:

| キー | 型 | 必須 | 許容値 / 説明 |
|------|----|------|---------------|
| `enabled` | boolean | 必須 | `true` / `false`。 |
| `min_samples` | integer | 必須 | 1〜10000。 |
| `avg_multiplier` | number | 必須 | 1.0〜100.0。 |
| `p95_multiplier` | number | 必須 | 1.0〜100.0。 |

未知 key、必須 key 欠落、型不一致、範囲外の値は validation failure とする。異常判定、history 更新、通知、trend 更新の順序は [`docs/details/runner.md` 詳細本文責務 §27.38](runner.md#sec-27-38) を参照する。

ApiRateLimitPolicy object:

| キー | 型 | 必須 | 許容値 / 説明 |
|------|----|------|---------------|
| `enabled` | boolean | 必須 | `true` / `false`。 |
| `groups` | object | 必須 | `login`、`read`、`trigger`、`operate`、`config`、`admin` の 6 key を過不足なく持つ。 |
| `groups.<group>.window_seconds` | integer | 必須 | 1〜86400。 |
| `groups.<group>.max_requests` | integer | 必須 | 1〜100000。 |

未知 key、group の不足または追加、必須 key 欠落、型不一致、範囲外の値は validation failure とする。固定窓の判定、actor / IP key、count 更新、監査は [`docs/details/security.md` 詳細本文責務 §27.47](security.md#sec-27-47) を参照する。HTTP request / response は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) を参照する。

`.server_config` の `POST /api/config` では `force_build_interval_hours`、`build_cooldown_seconds`、`schedule_interval_seconds`、`schedule_paused`、`allowed_hours` を直接更新してはならない。これらは専用スケジュール API からのみ更新する。

**共通 field 定義：**

| 共通 field | 型 | 許容値 | 説明 |
|------------|----|--------|------|
| notification `label` | string | 0〜64 文字 | 管理画面表示名。 |
| notification `retry_interval_seconds` | integer | 1〜3600 | 再試行間隔。 |
| build `output_size_bytes` | integer/null | 0 以上または `null` | 成果物サイズ。 |
| build `output_sha256` | string/null | SHA-256 hex または `null` | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の出力成果物 manifest SHA-256。 |
| build `flagged` | boolean | boolean | 重要フラグ。 |
| build `comment` | string/null | 0〜2000 文字または `null` | コメント。 |
| request `remote_addr` | string/null | IP 文字列または `null` | 接続元。 |

**`.notify_config` schema：**

| キー | 型 | 既定値 | 許容値 | 説明 |
|------|----|--------|--------|------|
| `webhooks` | object[] | `[]` | 次の Webhook object | 互換通知先一覧。`channels` が空の場合、runner は `webhooks` を webhook channel として扱う。 |
| `channels` | object[] | `[]` | 次の Channel object | 統一通知 channel 一覧。`channels` が存在する場合、runner は `channels` を優先し、`webhooks` / `email` は互換表示用として扱う。 |
| `on` | string[] | `[]` | `"start"`, `"success"`, `"failure"`, `"deploy_failure"`, `"weekly_summary"`, `"approval_required"`, `"duration_anomaly"`, `"config_corrupt"` | 通知イベント。重複は除去する。 |
| `summary` | object | 次の Summary object | 次の | 定期サマリー設定。 |
| `email` | object | 次の Email object | 次の | メール通知設定。SMTP 詳細は `.smtp_config` / `.smtp_secret` を基準とする。 |

Channel object:

| キー | 型 | 既定値 | 許容値 | 説明 |
|------|----|--------|--------|------|
| `id` | string | 自動採番 | `n` + 数字、または 1〜64 文字の英数字 `_` `-` | channel 識別子。 |
| `type` | string | 必須 | `"webhook"` / `"email"` / `"command"` | 送信方式。 |
| `label` | string | `""` | Channel object は共通 field `notification label` | 管理画面表示名。 |
| `enabled` | boolean | `true` | boolean | `false` の channel へは送信しない。 |
| `on` | string[] | `[]` | top-level `on` と同じ、または `"*"` | 空配列の場合は top-level `on` の event 集合を使用する。 |
| `config` | object | `{}` | type 別 schema | webhook `url`、email `to`、command `command_args`。 |
| `retry_count` | integer | `2` | 0〜10 | retry 対象失敗時の追加試行回数。 |
| `retry_interval_seconds` | integer | `30` | Channel object は共通 field `notification retry_interval_seconds` | 再試行間隔。 |

Webhook object:

| キー | 型 | 既定値 | 許容値 | 説明 |
|------|----|--------|--------|------|
| `url` | string | 必須 | `http` / `https`、host 必須、userinfo / fragment 禁止 | 送信先 URL。 |
| `label` | string | `""` | Webhook object は共通 field `notification label` | 管理画面表示名。 |
| `enabled` | boolean | `true` | boolean | `false` の宛先へは送信しない。 |
| `on` | string[] | `[]` | `"start"`, `"success"`, `"failure"`, `"deploy_failure"`, `"weekly_summary"`, `"approval_required"`, `"duration_anomaly"`, `"config_corrupt"`, `"*"` | この宛先が受け取るイベント。空配列の場合は top-level `on` の event 集合を使用する。 |
| `payload_template` | string/null | `null` | 0〜10000 文字または `null` | `null` は標準 payload。 |
| `retry_count` | integer | `2` | 0〜10 | 送信失敗時の追加試行回数。 |
| `retry_interval_seconds` | integer | `30` | Webhook object は共通 field `notification retry_interval_seconds` | 再試行間隔。 |
| `secret` | string/null | `null` | 1〜256 文字または `null` | 保存時は平文保存可。ただし GET/backup では `"***"` へマスクする。 |

Summary object:

| キー | 型 | 既定値 | 許容値 |
|------|----|--------|--------|
| `enabled` | boolean | `false` | boolean |
| `interval` | string | `"weekly"` | `"weekly"` 固定 |
| `hour` | integer | `9` | 0〜23 |
| `day_of_week` | integer | `1` | 0〜6 |

Email object:

| キー | 型 | 既定値 | 許容値 |
|------|----|--------|--------|
| `enabled` | boolean | `false` | boolean |
| `to` | string[] | `[]` | メールアドレス配列、最大 50 件 |
| `on` | string[] | `[]` | `"start"`, `"success"`, `"failure"`, `"duration_anomaly"` |

**`.notify_log` JSON Lines schema：**

`.notify_log` は通常通知または手動 test の channel 処理結果を記録する。外部送信を開始した場合は 1 attempt につき 1 record だけを追記し、最終 retry 失敗を `failure` と `dropped` の 2 record に重複記録してはならない。送信開始前の SMTP 未設定判定と、最終 retry 失敗による破棄判定も、それぞれ 1 record だけを追記する。payload 本文、送信先 URL、email address、command 引数、stdout、stderr、secret は保存しない。

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `id` | string | 必須 | base は `ntfy{YYYYMMDDHHmmss}` | 衝突処理は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の時刻ベース ID 契約に従う。 |
| `at` | string | 必須 | UTC ISO 8601 秒精度 | attempt 結果確定時刻。 |
| `event` | string | 必須 | `.notify_config.on` の列挙値、`"notify_test"`、`"smtp_test"` | 通知 event。test 値は対応する手動 test API のみ使用する。 |
| `channel_id` | string | 必須 | Channel object の `id`、または `event="smtp_test"` の場合だけ `"smtp-test"` | 正規化後 channel id。 |
| `channel_type` | string | 必須 | `"webhook"`, `"email"`, `"command"` | 送信方式。 |
| `payload_sha256` | string | 必須 | 64 文字 lowercase hex | secret mask 後 NotificationPayload object、または API 詳細本文で固定した test payload の canonical JSON byte に対する SHA-256。 |
| `result` | string | 必須 | `"success"`, `"failure"`, `"not_configured"`, `"dropped"` | channel 処理結果。 |
| `attempt` | integer | 必須 | 1 以上 | 初回の送信または送信前判定を `1` とし、保存に成功した retry 結果ごとに 1 を加える。retry の log 保存失敗時は pending の `attempts` を変更せず、次回に同じ番号を再使用する。`dropped` は最後に実行した retry の番号を使い、別番号を採番しない。 |
| `http_status` | integer/null | 必須 | 100〜599 または `null` | HTTP response を受信した webhook だけ整数。 |
| `error_code` | string/null | 必須 | `"http_1xx"`, `"http_3xx"`, `"http_4xx"`, `"http_5xx"`, `"timeout"`, `"network_error"`, `"smtp_not_configured"`, `"smtp_error"`, `"command_error"`, `"retry_exhausted"`, `null` | 送信失敗分類。 |
| `error` | string/null | 必須 | 500 文字以下または `null` | secret mask 後の短い失敗理由。 |

`result="success"` では `error_code` と `error` を `null` とする。それ以外の result では `error_code` と `error` を非 `null` とする。`result="failure"` の固定値は次のとおりとする。

| 条件 | `error_code` | `error` | `http_status` |
|------|--------------|---------|---------------|
| Webhook HTTP `100`〜`199` | `"http_1xx"` | `"HTTP {status}"` | 受信した status |
| Webhook HTTP `300`〜`399` | `"http_3xx"` | `"HTTP {status}"` | 受信した status |
| Webhook HTTP `400`〜`499` | `"http_4xx"` | `"HTTP {status}"` | 受信した status |
| Webhook HTTP `500`〜`599` | `"http_5xx"` | `"HTTP {status}"` | 受信した status |
| Webhook response 受信前の接続失敗 | `"network_error"` | `"network error"` | `null` |
| Webhook、SMTP、command の timeout | `"timeout"` | `"timeout"` | `null` |
| SMTP 接続・認証・送信失敗 | `"smtp_error"` | `"SMTP error"` | `null` |
| command process 起動失敗 | `"command_error"` | `"command start failed"` | `null` |
| command 非 0 終了 | `"command_error"` | `"command exit {exit_code}"` | `null` |

`{status}` は 3 桁の decimal HTTP status、`{exit_code}` は符号付き decimal process exit code とし、前後空白を付けない。`result="not_configured"` は必須 SMTP 設定がないため通常の email channel 送信を開始しなかった場合だけ使用し、`error_code="smtp_not_configured"`、`error="SMTP not configured"`、`http_status=null` とする。`result="dropped"` は最後に許可された retry が `http_5xx` または `timeout` で失敗し、pending entry を破棄する場合だけ使用し、`error_code="retry_exhausted"`、`error="retry exhausted"` とする。最後の失敗が HTTP 5xx なら受信 status を `http_status` に保存し、timeout なら `null` とする。最終 retry では同じ attempt 番号の `failure` record を別途追記しない。

`event="notify_test"` は `POST /api/notify-test` の 1 回の Webhook 送信だけに使い、`channel_id` は API が選択した正規化後 Webhook channel id、`channel_type="webhook"`、`attempt=1` とする。`event="smtp_test"` は `POST /api/smtp-test` の 1 回の SMTP 送信だけに使い、`channel_id="smtp-test"`、`channel_type="email"`、`attempt=1`、`http_status=null` とする。両 test event は `.notify_config.on`、Channel object の `on`、`.notify_pending`、NotificationPayload object の event 値として使用しない。

**`.notify_pending` schema：**

`.notify_pending` は NotificationPending object の JSON array とする。array は `created_at` 昇順、同時刻は `id` 昇順で保存する。

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `id` | string | 必須 | base は `np{YYYYMMDDHHmmss}` | 衝突処理は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の時刻ベース ID 契約に従う。 |
| `event` | string | 必須 | `.notify_config.on` の列挙値 | payload の `event` と一致させる。 |
| `channel_id` | string | 必須 | `type="webhook"` の Channel object の `id` | retry 対象 Webhook channel。 |
| `channel_type` | string | 必須 | `"webhook"` | retry 対象は Webhook の HTTP 5xx / timeout だけとする。 |
| `payload` | object | 必須 | NotificationPayload object | secret mask 後 payload。 |
| `attempts` | integer | 必須 | 1 以上 | 初回送信を含む実行済み attempt 数。 |
| `next_attempt_at` | string | 必須 | UTC ISO 8601 秒精度 | 次回 retry を許可する最早時刻。 |
| `last_error` | string | 必須 | 1〜500 文字 | secret mask 後の最終失敗理由。 |
| `created_at` | string | 必須 | UTC ISO 8601 秒精度 | pending 初回作成時刻。 |

NotificationPayload object は全 event 共通で `event`、`build_id`、`status`、`branch`、`trigger`、`created_at` を必須 key とし、取得できない値は `null` とする。`event` は string、その他の共通 key は string/null とする。event 別に追加できる key は以下だけとし、それ以外の未知 key は禁止する。

| event | 追加 key |
|-------|----------|
| `start`, `success`, `failure` | なし。 |
| `deploy_failure` | `target_id` string/null、`reason` string/null。 |
| `weekly_summary` | `period_days` integer、`period_from` string、`period_to` string、`success_count` integer、`failure_count` integer、`success_rate` number、`avg_duration_seconds` number/null、`max_duration_seconds` number/null、`max_duration_build_id` string/null。 |
| `approval_required` | `approval_id` string、`sha` string、`target` string、`requested_trigger` string、`requested_force` boolean、`expires_at` string。値は対応する pending approval record から変更せずコピーする。 |
| `duration_anomaly` | `duration_seconds` number、`avg_seconds` number/null、`p95_seconds` number/null、`threshold_source` は `"avg"`, `"p95"`, `"avg_and_p95"` のいずれか。 |
| `config_corrupt` | `files` string[]、`recovered` boolean。 |

NotificationPayload object は canonical JSON 生成時に key を ASCII 昇順で並べ、配列は処理責務で定めた順序を保持する。Webhook secret、SMTP password、API token、session token、repository token、Authorization header、command environment value を保存しない。

**`.pending_transfers` schema：**

`.pending_transfers` は PendingTransfer object の JSON array とする。`branch_idx`、`deploy_idx`、`attempt` は旧 key としても受け付けず、保存しない。

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `build_id` | string | 必須 | build id | 転送対象成果物の build id。 |
| `trigger` | string | 必須 | `"deploy"`, `"rollback"` | pending transfer の発生経路。 |
| `source_kind` | string | 必須 | `"output"`, `"snapshot"` | 再送元の種別。 |
| `rollback_from` | string/null | 必須 | snapshot 元 build id または `null` | `trigger="rollback"` の場合だけ非 `null`。 |
| `snapshot_id` | string/null | 必須 | snapshot id または `null` | `trigger="rollback"` の場合だけ非 `null`。 |
| `branch` | string | 必須 | branch target 名 | 転送元 branch。 |
| `target_id` | string | 必須 | Deploy object の `target_id` | 設定順に依存しない deploy target 識別子。 |
| `out` | string/null | 必須 | 絶対 path または `null` | `source_kind="output"` の local output directory。`snapshot` では `null`。 |
| `host` | string | 必須 | `^[A-Za-z0-9._-]{1,255}$` | deploy target host。 |
| `user` | string | 必須 | `^[A-Za-z0-9._-]{1,64}$` | deploy target user。 |
| `dest_dir` | string | 必須 | 絶対 path | remote destination directory。 |
| `output_sha256` | string | 必須 | 64 文字 lowercase hex | pending 投入時の成果物 manifest SHA-256。 |
| `failed_at` | string | 必須 | UTC ISO 8601 秒精度 | 直近失敗時刻。 |
| `retry_count` | integer | 必須 | 0 以上 | pending 投入後に実行した retry 回数。初期値は `0`。 |
| `last_error` | string | 必須 | 1〜500 文字 | secret mask 後の最終失敗理由。 |

`trigger="deploy"` は `source_kind="output"`、`out`非 `null`、`rollback_from=null`、`snapshot_id=null` とする。`trigger="rollback"` は `source_kind="snapshot"`、`out=null`、`rollback_from`と `snapshot_id` を非 `null` の同一 snapshot id とする。それ以外の組合せは schema 不正とする。`output_sha256` は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の固定アルゴリズムで算出した値とする。

array は投入順を保持する。転送先識別子は `branch`、`target_id`、`host`、`user`、`dest_dir` の完全一致とする。同一転送先の新しい build が pending になった場合は、既存 entry の array 位置を保ったまま `build_id`、`trigger`、`source_kind`、`rollback_from`、`snapshot_id`、`out`、`output_sha256`、`failed_at`、`last_error` を新しい値へ置換し、`retry_count=0` とする。同一転送先への新しい build が転送成功した場合は、対応する既存 pending entry を削除する。再送元の検証、snapshot 再展開、転送、cleanup は [`docs/details/runner.md` 詳細本文責務 §14a](runner.md#14a-ssh-サイト転送) と [`docs/details/archive.md` 詳細本文責務 §27.15](archive.md#sec-27-15) を参照する。

<a id="branch-config-schema"></a>
**`.branch_config` schema：**

```json
{
  "branch_targets": [
    {
      "branch": "main",
      "target_file": "docs",
      "target_files": ["docs"],
      "sha_file": "/opt/adlaire-builder/.last_sha",
      "src": "/opt/adlaire-builder/repo/docs",
      "out": "/opt/adlaire-builder/dist/site",
      "approval_required": false,
      "env": {},
      "deploy_targets": [
        { "id": "primary", "host": "192.0.2.1", "user": "deploy", "dest_dir": "/var/www/html/" }
      ]
    }
  ]
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `branch_targets` | object[] | 必須 | 0〜50 件 | 永続ファイルの保存 key。空配列は `.branch_config` 削除と同義。 |
| `branch` | string | 必須 | 1〜128 文字、`refs/heads/` は含めない、`branch_targets` 内で一意 | GitHub branch 名。 |
| `target_file` | string | 必須 | 相対パス、`..` 禁止 | GitHub リポジトリ内の監視対象ファイル。 |
| `target_files` | string[] | 任意 | 1〜100 件の相対パス | 複数監視対象。省略時は `[target_file]`。保存時は `target_file` を必ず含める。 |
| `sha_file` | string | 必須 | 絶対パス | 対象 branch/file の SHA キャッシュ。 |
| `src` | string | 必須 | 絶対パス | blob 本文の書き出し先。 |
| `out` | string | 必須 | 絶対パス | ビルド成果物パス。 |
| `approval_required` | boolean | 任意 | boolean | 省略時は `false`。 |
| `env` | object | 任意 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Process environment entry 共通固定契約](../DETAIL_INDEX.md#process-environment-entry-contract) | 省略時は `{}`。branch target 固有の process environment。 |
| `deploy_targets` | object[] | 必須 | 0〜20 件 | SSH 転送先。空配列は転送なし。 |
| `deploy_targets[].id` | string | 任意 | `^[A-Za-z0-9_-]{1,64}$`、同一 branch target 内で一意 | 省略時は正規化前の `{branch_index}-{target_index}`。 |
| `deploy_targets[].host` | string | 必須 | `^[A-Za-z0-9._-]{1,255}$` | SSH host。 |
| `deploy_targets[].user` | string | 必須 | `^[A-Za-z0-9._-]{1,64}$` | SSH user。 |
| `deploy_targets[].dest_dir` | string | 必須 | 正規化済み絶対 POSIX path | 転送先ディレクトリ。 |

`target_files` は各 path を `/` 区切りへ正規化し、`.` segment を除去し、大文字小文字を区別して重複排除した後、辞書順で保存する。空文字、絶対 path、`..` segment、NUL、改行、CR を含む path は許可しない。`target_file` は正規化済み `target_files` に必ず含める。複数 target の差分判定と SHA cache は [`docs/details/runner.md` 詳細本文責務 §27.21](runner.md#sec-27-21) を参照する。

`deploy_targets[].dest_dir` は `/` で開始し、有効な UTF-8 文字列とし、NUL、CR、LF、その他の制御文字、および正規化前の `..` segment を禁止する。入力を Go `path.Clean` で正規化し、root 以外の末尾 `/` を除いた値を保存する。この正規化後の `dest_dir` を deploy target の一致判定と SSH 転送のみに使用する。

`env` の process 注入と secret mask は [`docs/details/runner.md` 詳細本文責務 §27.31](runner.md#sec-27-31) を参照する。

`target_files`、`approval_required`、`env`、`deploy_targets[].id` の欠落は `.branch_config` に限る旧形式正規化例外とする。読取時は [`.branch_config` schema](#branch-config-schema) の既定値または導出値を memory 上で補い、読取だけではファイルを書き換えない。次回の API 保存時は正規化後の全 key を明示保存する。`deploy_targets[].id` の導出に使う index は request の配列順を基準とする。正規化後は `branch_targets` を branch 名の byte 昇順、各 `deploy_targets` を id の byte 昇順で保存する。branch 重複または同一 branch target 内の deploy target id 重複は schema 不正とし、API 入力は `422`、保存済み状態の読取は `ErrStateCorrupted` とする。

`.branch_config` の永続 key は `branch_targets` に固定する。API の表示名、default 復帰、空配列入力時の挙動は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) および [`docs/details/api.md` 詳細本文責務 §27.18](api.md#sec-27-18) を参照する。statefile は `branches` を永続 key として保存してはならない。

<a id="last-sha-schema"></a>
**`.last_sha` / `BranchTarget.SHAFile` schema：**

```json
{"sha":""}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `sha` | string | 必須 | 空文字または Git blob SHA | 空文字は初回実行扱い。 |

SHA cache は target ごとの処理済み Git blob SHA を保存する JSON object とする。legacy text 形式は自動変換しない。

| 状態 | 読込時の扱い | 保存時の扱い |
|------|--------------|--------------|
| 不在 | 初回実行として `previous_blob_sha=""` を扱う。 | build 成功時に `{"sha":"<new_sha>"}` を atomic write する。 |
| `{"sha":""}` | 初回実行として扱う。 | build 成功時に新 SHA で上書きする。 |
| `{"sha":"<value>"}` | `<value>` を前回 SHA として比較する。 | build 成功時に新 SHA で上書きする。 |
| JSON 破損 | 当該 target の decode failure として扱う。 | 更新しない。 |
| object 以外 | JSON 破損と同じ扱い。 | 更新しない。 |
| `sha` key 不在または string 以外 | JSON 破損と同じ扱い。 | 更新しない。 |
| 未知 key あり | schema 破損として当該 target の decode failure とする。 | 更新しない。 |

SHA cache の更新タイミング、skip / failure 時の更新可否、複数 target 時の個別更新は [`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) の SHA cache 読み書き契約を参照する。

**`.sha_cache/{branch_safe}/{target_hash}.sha` schema：**

file 内容は Git object SHA の lowercase hexadecimal 40 文字または 64 文字と、末尾 LF 1 個だけとする。先頭空白、末尾空白、CRLF、複数行、大文字 hex、`0x` prefix、末尾 LF なしは不正とする。不在または不正値は cache miss とし、不正な既存 file を読取時に自動更新しない。path の `branch_safe` / `target_hash` 導出、比較、成功時更新、skip / failure 時の更新禁止は [`docs/details/runner.md` 詳細本文責務 §27.21](runner.md#sec-27-21) を参照する。

**`.local_watch_state.json` schema：**

```json
{"files":{"docs/index.md":{"sha256":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","mtime_unix":0,"size":0}}}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `files` | object | 必須 | 0 件以上 | key は `src` からの正規化済み相対 path。 |
| `files.<path>.sha256` | string | 必須 | 64 文字 lowercase hex | file 内容の SHA-256。 |
| `files.<path>.mtime_unix` | integer | 必須 | 0 以上 | scan 時の Unix 秒。差分判定の正本は `sha256` とする。 |
| `files.<path>.size` | integer | 必須 | 0 以上 | scan 時の byte 数。 |

`files` の path は `/` 区切り、辞書順とし、空文字、先頭 `/`、空 segment、`.`、`..`、NUL、改行、CR を禁止する。未知 key を保存してはならない。走査対象、除外 path、差分判定、破損時 full build、成功時置換、dry-run の更新禁止は [`docs/details/runner.md` 詳細本文責務 §27.23](runner.md#sec-27-23) を参照する。

**`.build_cache.json` / `.build_cache/pages/{cache_key}.json` schema：**

`.build_cache.json` は cache index として次の JSON object に固定する。

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `schema_version` | integer | 必須 | `1` | 未知 version は全 entry cache miss。 |
| `entries` | object | 必須 | string:string | key は入力 base からの Markdown 相対 path、value は 64 文字 lowercase hex の cache key。 |

`entries` の key は `/` 区切り、ASCII 昇順で保存し、空文字、絶対 path、`.` / `..` segment、NUL、改行、CR を禁止する。同一 input path に複数 cache key を保存しない。

`.build_cache/pages/{cache_key}.json` は次の CachePage object に固定する。

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `cache_key` | string | 必須 | file 名と同じ 64 文字 lowercase hex | cache entry 識別子。 |
| `input_path` | string | 必須 | index の input relative path | cache index key と一致させる。 |
| `input_sha256` | string | 必須 | 64 文字 lowercase hex | 入力 Markdown byte SHA-256。 |
| `builder_version` | string | 必須 | 非空文字 | builder version。 |
| `theme` | string | 必須 | builder の theme 列挙値 | 変換時 theme。 |
| `build_config_hash` | string | 必須 | 64 文字 lowercase hex | 正規化済み build config SHA-256。 |
| `deps` | object | 必須 | string:string | key は dependency 相対 path、value は 64 文字 lowercase hex。 |
| `html_fragment` | string | 必須 | UTF-8 string | page body の変換済み HTML fragment。 |
| `metadata` | object | 必須 | CachePageMetadata object | site assembly と search index 再生成に必要な page metadata。 |
| `created_at` | string | 必須 | UTC ISO 8601 秒精度 | cache entry 作成時刻。 |

CachePageMetadata object は `title` string、`slug` string、`output_path` string、`headings` object[]、`warnings` string[]、`reading_time` integer を必須 key とする。`headings[]` は `level` integer 1〜6、`id` string、`text` string の 3 key だけを持つ。`output_path` と `deps` の key は input base または output root からの `/` 区切り相対 path とし、絶対 path、user home path、temporary path を保存しない。

cache key の導出、hit / miss 判定、byte 等価検証、save タイミング、WARN / REPORT は [`docs/details/builder.md` 詳細本文責務 §27.25](builder.md#sec-27-25) を参照する。

**`.dependency_manifest.json` schema：**

`.dependency_manifest.json` は次の DependencyManifest object に固定する。不在、JSON 破損、必須 key 不足、未知 key、version 不一致は schema 不一致とし、旧 `{"pages":{}}` 形式を自動変換しない。

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `version` | integer | 必須 | `1` | manifest schema version。 |
| `builder` | object | 必須 | DependencyManifestBuilder object | manifest 生成条件。 |
| `pages` | object | 必須 | string:DependencyManifestPage | key は input base からの Markdown 相対 path。 |
| `generated_outputs` | object | 必須 | string:string | key は output root からの相対 path、value は page key または `"asset"`, `"search-index"`, `"manifest"`。 |

DependencyManifestBuilder object:

| キー | 型 | 必須 | 許容値 |
|------|----|------|--------|
| `name` | string | 必須 | `"adlaire-ci-build"` |
| `spec_section` | string | 必須 | `"28.1"` |
| `config_hash` | string | 必須 | 64 文字 lowercase hex |

DependencyManifestPage object:

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `input_sha256` | string | 必須 | 64 文字 lowercase hex | 入力 Markdown byte SHA-256。 |
| `output_path` | string | 必須 | 正規化済み相対 path | 対応する HTML output path。 |
| `dependencies` | string[] | 必須 | ASCII 昇順、重複なし | input base 内の参照 file path。 |
| `dependency_sha256` | object | 必須 | string:string | key は `dependencies` の各 path、value は 64 文字 lowercase hex。 |
| `config_hash` | string | 必須 | 64 文字 lowercase hex | 当該 page の正規化済み設定 SHA-256。 |

`pages`、`generated_outputs`、`dependency_sha256` の object key は ASCII 昇順で保存する。すべての path は `/` 区切り相対 path とし、空文字、絶対 path、`.` / `..` segment、NUL、改行、CR を禁止する。timestamp、host path、user name、temporary path、random value、secret を保存しない。dependency 抽出、broken dependency 処理、build 対象判定、保存タイミングは [`docs/details/builder.md` 詳細本文責務 §27.28](builder.md#sec-27-28) と [§28.1](builder.md#sec-28-1) を参照する。

<a id="repo-config-schema"></a>
**`.repo_config` schema：**

```json
{
  "owner": "fqwink",
  "repo": "Build-Scripts",
  "updated_at": "2026-09-15T10:00:00Z"
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `owner` | string | 必須 | `^[A-Za-z0-9](?:[A-Za-z0-9-]{0,98}[A-Za-z0-9])?$` | GitHub owner。URL path segment としてそのまま使用できる ASCII だけを許可する。 |
| `repo` | string | 必須 | `^[A-Za-z0-9](?:[A-Za-z0-9._-]{0,98}[A-Za-z0-9])?$` | GitHub repository 名。`/`、query、fragment、空白、制御文字は許可しない。 |
| `updated_at` | string/null | 必須 | UTC ISO 8601 秒精度または `null` | `null` は `.repo_config` 不在時に `readRepoConfig()` が返す固定既定値だけで使用する。API 保存時は request 完了時刻を保存する。 |

`.repo_config` は GitHub repository 識別子だけを保存する正本とする。branch、監視対象、source path、output path、deploy target を保存してはならず、これらは [`.branch_config` schema](#branch-config-schema) を唯一の正本とする。

repo config write caller は request の `owner` または `repo` のうち指定された key だけを、`readRepoConfig()` の正規化済み現在値へ適用する。未指定 key は現在値を保持し、全 key 未指定は validation failure とする。差分がある場合だけ `owner`、`repo`、`updated_at` の 3 key をすべて保存する。未知 key、`branch`、`target_file` を含む request は validation failure とし、既存 file を変更しない。

**`.totp_secret` schema：**

```json
{
  "enabled": true,
  "secret_base32": "JBSWY3DPEHPK3PXP",
  "confirmed_at": "2026-09-15T10:00:00Z",
  "last_accepted_step": 59652320
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `enabled` | boolean | 必須 | boolean | TOTP 有効状態。 |
| `secret_base32` | string/null | 必須 | RFC 4648 base32、padding なし、16〜64 文字、または `null` | TOTP secret。呼び出し元の公開値、log、backup へ平文出力しない。 |
| `confirmed_at` | string/null | 必須 | ISO 8601 または `null` | TOTP 有効化完了日時。 |
| `last_accepted_step` | integer/null | 必須 | Unix time 30 秒 step または `null` | 同一 code 再利用防止。 |

**`.admin_credentials` schema：**

```json
{
  "password_hash": "<sha256_iter_v1_hex>",
  "salt": "<hex>",
  "algorithm": "sha256_iter_v1",
  "iterations": 260000,
  "must_change": true,
  "login_count": 0,
  "last_login_at": null,
  "updated_at": "2026-09-15T10:00:00Z"
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `password_hash` | string | 必須 | 64 文字 lowercase hex | password 本体は保存しない。 |
| `salt` | string | 必須 | 64 文字 lowercase hex | 32 bytes salt。 |
| `algorithm` | string | 必須 | `"sha256_iter_v1"` 固定 | 他 algorithm は初期実装で拒否する。 |
| `iterations` | integer | 必須 | `260000` 固定 | 値が異なる場合は認証 caller へ破損失敗を返す。 |
| `must_change` | boolean | 必須 | boolean | 初期生成時 `true`、パスワード変更後 `false`。 |
| `login_count` | integer | 必須 | 0 以上 | session token 発行成功時だけ +1。TOTP ticket 発行時は増やさない。 |
| `last_login_at` | string/null | 必須 | UTC ISO 8601 または `null` | session token 発行成功時だけ更新する。 |
| `updated_at` | string | 必須 | UTC ISO 8601 | password hash 更新時刻。 |

`.admin_credentials` に未知 key がある場合は credentials 破損として扱い、自動削除しない。必須 key 不足、型不一致、hex 不正、`algorithm` 不一致、`iterations` 不一致もすべて credentials 破損とする。API 起動時検証、login / password change の公開応答、認証ログ、監査ログ、漏えい禁止値は [`docs/details/security.md` 詳細本文責務 認証共通詳細](security.md#認証共通詳細) および [`docs/details/security.md` 詳細本文責務 §27.45](security.md#sec-27-45)〜[§27.46](security.md#sec-27-46) を参照する。statefile は破損内容、hash、salt を呼び出し元の公開値として返してはならない。

**`.audit_log` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `timestamp` | string | 必須 | ISO 8601 | 発生日時。 |
| `request_id` | string/null | 必須 | 32 文字 lowercase hex または `null` | API request は [`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) で生成する識別子。API request に紐づかない `runner` 内部 event は `null`。 |
| `actor_type` | string | 必須 | `"admin"` / `"api_token"` / `"webhook"` / `"system"` / `"anonymous"` | 操作者種別。 |
| `actor_id` | string/null | 必須 | `"admin"`、token id、`"webhook"`、`"system"`、または `null` | 操作者。secret 本体は保存しない。 |
| `action` | string | 必須 | [`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) | 操作種別。 |
| `target_type` | string | 必須 | [`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) | 対象種別。 |
| `target_id` | string/null | 必須 | 対象 id または `null` | 対象識別子。 |
| `result` | string | 必須 | `"success"` / `"failure"` / `"denied"` | 結果。 |
| `remote_addr` | string/null | 必須 | IP 文字列または `null` | [`docs/details/api.md` 詳細本文責務 §27.6](api.md#sec-27-6) と同じ `RemoteAddr` 導出規則で得た接続元。forwarded header は参照しない。 |
| `message` | string/null | 必須 | 0〜500 文字または `null` | 固定文言。secret、token、password は保存しない。 |

**`.api_rate_state` schema：**

```json
{
  "windows": {
    "ip:127.0.0.1:login": {
      "window_start": "2026-09-15T10:00:00Z",
      "count": 1
    }
  }
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `windows` | object | 必須 | key は `{dimension}:{id}:{endpoint_group}`。`dimension` は `session`, `token`, `ip`、`endpoint_group` は `login`, `read`, `trigger`, `operate`, `config`, `admin` | 固定窓状態。 |
| `window_start` | string | 必須 | ISO 8601 | 現在窓の開始時刻。 |
| `count` | integer | 必須 | 0 以上 | 現在窓内リクエスト数。 |

**`.api_tokens` schema：**

```json
{
  "tokens": [
    {
      "id": "tok000001",
      "label": "監視用",
      "scopes": ["read"],
      "token_hash": "<sha256_hex>",
      "created_at": "2026-09-15T10:00:00Z",
      "last_used_at": null,
      "expires_at": null,
      "revoked_at": null
    }
  ]
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `tokens` | object[] | 必須 | 0〜100 件 | 発行済み API token 一覧。 |
| `id` | string | 必須 | `tok` + 6 桁以上の数字 | token 識別子。 |
| `label` | string | 必須 | 1〜64 文字 | 表示名。 |
| `scopes` | string[] | 必須 | `read`, `trigger`, `operate`, `config`, `admin` の 1〜5 件 | token に許可する scope。 |
| `token_hash` | string | 必須 | SHA-256 hex | token 本体は保存しない。 |
| `created_at` | string | 必須 | ISO 8601 | 作成日時。 |
| `last_used_at` | string/null | 必須 | ISO 8601 または `null` | 最終使用日時。 |
| `expires_at` | string/null | 必須 | ISO 8601 または `null` | 有効期限。`null` は無期限。 |
| `revoked_at` | string/null | 必須 | ISO 8601 または `null` | 失効日時。`null` は有効。 |

`.api_tokens` の永続化では token 本体を保存せず、`token_hash` だけを保存する。旧 `scope` 文字列が存在する場合は、旧形式の唯一の読取例外として `scopes:[scope]` へ正規化できる。それ以外の個別 record が schema 不正の場合、read adapter は破損失敗を返し、自動補正、部分除外、再生成を行わない。token 生成、id 採番、label / scope / 期限の正規化、一覧順、失効、認証、`last_used_at` 更新、HTTP response は [`docs/details/security.md` 詳細本文責務 §27.43](security.md#sec-27-43) を正本とする。

**`.maintenance` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `enabled` | boolean | 必須 | boolean | メンテナンス有効状態。 |
| `reason` | string/null | 必須 | 0〜500 文字または `null` | 理由。 |
| `since` | string/null | 必須 | ISO 8601 または `null` | 有効化日時。 |

`enabled: false` の場合、`reason` と `since` は `null` とする。遷移条件、HTTP request / response、副作用は [`docs/details/api.md`](api.md) 詳細本文責務のメンテナンスモード契約を参照する。

**`.access_control` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `allow` | string[] | 必須 | IPv4 address または CIDR、0〜100 件 | 空配列は制限なし。保存時は入力順を保持し、重複は除去する。 |

`allow` の各要素は前後空白を除去してから検証する。空文字、IPv6、hostname、URL、CIDR prefix が 0〜32 以外、parse 不能な値は validation failure とし、既存 `.access_control` を変更しない。

**`.hooks` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `hooks` | object[] | 必須 | 0〜50 件 | 登録済み hook。 |
| `id` | string | 必須 | [`docs/details/api.md` 詳細本文責務 §22.0e.2](api.md#sec-22-0e-2) の hook id | hook 識別子。 |
| `phase` | string | 必須 | `"pre"` / `"post"` | 実行 phase。 |
| `command_args` | string[] | 必須 | 1〜20 件 | shell を介さず `exec.CommandContext` に渡す引数配列。 |
| `enabled` | boolean | 必須 | boolean | `false` の hook は実行しない。 |
| `abort_on_failure` | boolean | 必須 | boolean | `pre` hook 失敗時だけ参照する。 |
| `timeout_seconds` | integer | 必須 | 1〜3600 | hook 単体の timeout。未指定作成時は `300`。 |
| `created_at` | string | 必須 | UTC ISO 8601 | hook 作成日時。 |

`command_args[0]` は 1〜256 文字、`command_args[1:]` の各要素は 1〜500 文字とし、NUL、改行、CR を禁止する。`command_args[0]` は絶対 path または PATH 解決可能なコマンド名に限定する。`phase` と `command_args` が既存 enabled hook と完全一致する場合の重複時の扱いは [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) を参照する。

**`.alert_rules` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `rules` | object[] | 必須 | 0〜100 件 | dashboard alert rule。 |
| `id` | string | 必須 | [`docs/details/api.md` 詳細本文責務 §22.0e.2](api.md#sec-22-0e-2) の alert rule id | rule 識別子。 |
| `metric` | string | 必須 | `success_rate_7d`, `avg_duration_seconds`, `last_build_age_hours`, `disk_usage_bytes` | 評価対象。 |
| `operator` | string | 必須 | `lt`, `gt`, `lte`, `gte` | 比較演算子。 |
| `threshold` | number | 必須 | 0 以上 | 比較値。 |
| `level` | string | 必須 | `info`, `warn`, `error` | alert severity。 |
| `message` | string | 必須 | 1〜200 文字 | UI 表示文。secret を含めない。 |

**`.tag_rules` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `rules` | object[] | 必須 | 0〜100 件 | 自動タグ付け rule。 |
| `id` | string | 必須 | [`docs/details/api.md` 詳細本文責務 §22.0e.2](api.md#sec-22-0e-2) の tag rule id | rule 識別子。 |
| `condition` | string | 必須 | [`docs/details/api.md` 詳細本文責務 自動タグ付けルール](api.md#tag-rule-api) の条件式 grammar | 評価条件。 |
| `tags` | string[] | 必須 | 1〜20 件、各 1〜50 文字 | 付与するタグ。重複は除去する。 |

`condition` は `変数 空白 演算子 空白 値` の 1 条件だけを許可する。`&&`、`||`、括弧、関数呼び出し、正規表現、算術式は validation failure とする。

<a id="pipeline-config-schema"></a>
**`.pipeline_config` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `extra_args` | string[] | 必須 | 0〜50 件 | `builder` に渡す追加 CLI 引数。 |
| `env` | object | 必須 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Process environment entry 共通固定契約](../DETAIL_INDEX.md#process-environment-entry-contract) | builder process に追加する環境変数。 |
| `inline_yaml` | string/null | 必須 | UTF-8 で 0〜262144 bytes、または `null` | `.pipeline.yml` 不在時に使用する内製 YAML subset。空文字は pipeline 未指定として扱う。 |

<a id="pipeline-config-reserved-builder-options"></a>
**`.pipeline_config.extra_args` 予約 builder option 固定契約：**

`extra_args` は空文字、NUL、改行、CR を禁止し、`--src`、`--out`、`--build-id`、`--commit-sha`、`--build-at`、`--cache-dir`、`--version`、`--help` およびこれらの `--name=value` 形式を指定してはならない。

`inline_yaml` は UTF-8 不正、BOM、NUL、CR を拒否し、LF は保持する。YAML grammar、source 優先順位、parse、step 実行は [`docs/details/runner.md` 詳細本文責務 §27.22](runner.md#sec-27-22) を参照する。

**`.dashboard_layout` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `widgets` | string[] | 必須 | [`docs/details/api.md` 詳細本文責務 ダッシュボードウィジェットカスタマイズ](api.md#dashboard-layout-api) の widget id、1〜9 件 | 表示 widget 順序。重複禁止。 |

**`.smtp_config` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `host` | string/null | 必須 | hostname または IP、または `null` | SMTP server。 |
| `port` | integer | 必須 | 1〜65535 | SMTP port。 |
| `user` | string/null | 必須 | 0〜255 文字または `null` | SMTP user。 |
| `tls` | boolean | 必須 | boolean | STARTTLS または TLS 使用。 |
| `from` | string/null | 必須 | email address または `null` | From address。 |
| `to` | string[] | 必須 | email address、0〜50 件 | 送信先。 |
| `on` | string[] | 必須 | `start`, `success`, `failure`, `duration_anomaly` | 送信イベント。 |
| `enabled` | boolean | 必須 | boolean | メール通知有効状態。 |

`.smtp_secret` は UTF-8 text とし、末尾 LF なしで password 本体だけを保存する。mode は `0600` 固定。`password` を `.smtp_config` に保存してはならない。password 未指定、`null`、非空文字の更新操作は [`docs/details/api.md`](api.md) 詳細本文責務の `POST /api/smtp-config` 契約を正本とする。

**`.build_state` schema：**

```json
{
  "running": false,
  "current_build_id": null,
  "active_queue_entry": null,
  "queued": [],
  "last_started_at": null,
  "last_finished_at": null,
  "weekly_summary_last_sent_at": null,
  "weekly_summary_sent_date": null
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `running` | boolean | 必須 | boolean | ビルド実行中状態。 |
| `current_build_id` | string/null | 必須 | build id または `null` | 実行中 build id。 |
| `active_queue_entry` | object/null | 必須 | Queue entry schema または `null` | runner が待機列から取得し、結果未確定の queue entry。 |
| `queued` | object[] | 必須 | running 中は 0〜`queue_max_size` 件、idle 中は 0〜`max(1, queue_max_size)` 件 | durable build dispatch queue。 |
| `last_started_at` | string/null | 必須 | ISO 8601 または `null` | 最終開始日時。 |
| `last_finished_at` | string/null | 必須 | ISO 8601 または `null` | 最終完了日時。 |
| `weekly_summary_last_sent_at` | string/null | 必須 | ISO 8601 または `null` | 週次サマリー最終送信日時。 |
| `weekly_summary_sent_date` | string/null | 必須 | `YYYY-MM-DD` または `null` | 週次サマリー二重送信防止日付。 |

Queue entry schema:

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `id` | string | 必須 | [`docs/details/api.md` 詳細本文責務 §22.0e.2](api.md#sec-22-0e-2) の queue id | queue id。 |
| `trigger` | string | 必須 | `"manual"`, `"webhook"`, `"approval"` | 起動種別。 |
| `queued_at` | string | 必須 | UTC ISO 8601 | queue 追加日時。 |
| `requested_by` | string | 必須 | `"admin"`, `"webhook"`, `"approval"`, API token id | queue 追加元。 |
| `priority` | string | 必須 | `"low"`, `"normal"`, `"high"`, `"urgent"` | 優先度。未指定作成時は `"normal"`。処理順は [`docs/details/runner.md` 詳細本文責務 §27.35](runner.md#sec-27-35) を参照する。 |
| `created_seq` | integer | 必須 | 1 以上 | 既存最大 + 1。 |
| `payload` | object | 必須 | trigger ごとの固定 payload | 不要時は `{}`。 |

`created_seq` の欠落だけは旧 queue entry の読取例外とする。API の read-only 表示では末尾扱いとし、状態ファイルを書き換えない。runner の取り出しまたは queue 更新前に `.build_state` lock 内で既存最大値 + 1 を保存し、正規化保存に失敗した場合は build を開始しない。具体的な順序は [`docs/details/runner.md` 詳細本文責務 §27.35](runner.md#sec-27-35) を参照する。その他の必須 key 欠落は queue entry 破損とする。

Queue entry `payload` は trigger ごとに以下を許可する。未知 key は API 追加時 validation failure、runner 読込時は queue entry 破損として当該 entry を処理せず ERROR ログに記録する。

| trigger | payload |
|---------|---------|
| `manual` | `{ "force": boolean }`。 |
| `webhook` | `{ "delivery_id": string, "branch": string, "sha": string }`。 |
| `approval` | `{ "approval_id": string, "branch": string, "sha": string, "target": string, "requested_trigger": string, "requested_force": boolean, "delivery_id": string/null }`。`requested_trigger` は `"polling"`, `"force_interval"`, `"manual"`, `"webhook"`, `"local_watch"` のいずれか。`requested_force` は SHA 一致時も build を実行する要求の場合だけ `true`。`delivery_id` は webhook 由来だけ非 `null`。 |

`running:true` の場合、`current_build_id` は非 `null` の build id とする。`running:false` の場合、`current_build_id` は `null` とする。この対応が崩れた object は `.build_state` 破損とし、read adapter は値を補完しない。`active_queue_entry` は runner が結果を確定できず再実行を必要とする場合に限り `running:false` でも非 `null` を許可する。`active_queue_entry` と `queued[]` に同じ queue id が同時に存在してはならない。waiting queue の表示・追加・clear は [`docs/details/api.md`](api.md) 詳細本文責務の queue API 契約、active 選択・取得・完了遷移は [`docs/details/runner.md` 詳細本文責務 §27.35](runner.md#sec-27-35) を参照する。

`.build_state` の key ごとの書込責務は次に固定する。`api` は `queued` の append / clear / 旧 entry 正規化だけを書き込み、`runner` は `active_queue_entry`、`running`、`current_build_id`、`last_started_at`、`last_finished_at`、`weekly_summary_last_sent_at`、`weekly_summary_sent_date` を書き込む。`api` は runner 起動要求だけで runner 所有 field を変更せず、`runner` は API 要求の waiting clear を実行しない。書込条件と遷移順は [`docs/details/api.md`](api.md) 詳細本文責務の queue / build API 契約と [`docs/details/runner.md` 詳細本文責務 §27.19](runner.md#sec-27-19) / [`docs/details/runner.md` 詳細本文責務 §27.35](runner.md#sec-27-35) を参照する。

queue 保存上限は、valid running `.build_lock`、`.build_state.running=true`、または `active_queue_entry != null` のいずれかが成立する場合に `queue_max_size`、すべて成立しない場合に `max(1, queue_max_size)` とする。上限判定、重複判定、append、clear、runner の active 移動は `.build_state` lock 内で行う。`queue_max_size=0` かつ active のない idle 状態の最初の entry は runner 起動要求を durable に渡す dispatch slot であり、2 件目を許可しない。

**`.build_circuit_state` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `open` | boolean | 必須 | boolean | circuit open 状態。runner の開始判定は [`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー)、API の拒否判定は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) を参照する。 |
| `consecutive_failures` | integer | 必須 | 0 以上 | 連続失敗回数。 |
| `opened_at` | string/null | 必須 | ISO 8601 または `null` | open に遷移した日時。 |
| `last_failure_at` | string/null | 必須 | ISO 8601 または `null` | 最終失敗日時。 |
| `last_error` | string/null | 必須 | 文字列または `null` | 最終失敗理由。 |

初期値は `{"open":false,"consecutive_failures":0,"opened_at":null,"last_failure_at":null,"last_error":null}` とする。reset の HTTP 条件、保存順、no-op、監査副作用は [`docs/details/api.md`](api.md) 詳細本文責務の `POST /api/circuit-breaker/reset` 契約を参照する。

**`.config_log` JSON Lines schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `at` | string | 必須 | ISO 8601 | 変更日時。 |
| `type` | string | 必須 | `^[a-z][a-z0-9_]{0,63}$` | [`docs/details/api.md` 詳細本文責務 §27.20](api.md#sec-27-20) の endpoint 固定名。 |
| `action` | string | 必須 | `"create"`, `"update"`, `"delete"` | 変更種別。 |
| `actor` | string | 必須 | `"admin"` または API token id | 操作者。token 本体は保存しない。 |
| `request_id` | string | 必須 | 32 文字 lowercase hex | [`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) で生成した同一 request の識別子。 |
| `endpoint` | string | 必須 | `{METHOD} {path_template}` | path parameter の実値と query を含めない。 |
| `result` | string | 必須 | `"success"`, `"partial_failure"` | 主状態変更と必須後続処理の結果。主状態変更後の後続処理失敗だけ `"partial_failure"` とする。 |
| `error` | string/null | 必須 | `null` または `^[a-z][a-z0-9_]{0,63}$` | `result="success"` では `null`。`result="partial_failure"` では caller 固有節が定める固定 error code。 |
| `diff` | object | 必須 | `{key:[before,after]}` | 変更前後。秘密情報は `"***"`。 |
| `diff_text` | string | 必須 | 1 文字以上 | 人間向け差分。秘密情報は `"***"`。 |

`diff` は `{key:[before,after]}`、`diff_text` は 1 行以上の文字列とする。両 field の生成、key 順、JSON 表現、再帰的 secret mask は [`docs/details/api.md` 詳細本文責務 §27.20](api.md#sec-27-20) を正本とする。statefile owner は API から受け取った mask 済み `diff` と `diff_text` の schema を検証して保存するだけとし、差分生成、mask 対象推測、再 mask、平文への復元を行わない。

**`.access_log` JSON Lines schema：**

各行は認証、ログアウト、password 変更、session 一括失効、API token 認証・作成・失効の security event を表す JSON object とする。秘密情報、session token、API token 本体を保存してはならない。

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `at` | string | 必須 | ISO 8601 | イベント日時。 |
| `action` | string | 必須 | `"login_failure"`, `"login_locked"`, `"login_success"`, `"totp_required"`, `"totp_failure"`, `"logout"`, `"password_change"`, `"session_revoke_all"`, `"token_create"`, `"token_revoke"`, `"token_auth"`, `"token_expired"`, `"token_revoked_reject"`, `"permission_denied"` | security event 種別。操作との対応は [`docs/details/security.md` 詳細本文責務 認証共通詳細](security.md#認証共通詳細) と [§27.43](security.md#sec-27-43) を参照する。 |
| `result` | string | 必須 | `"success"`, `"failure"` | 成否。 |
| `session_id` | string/null | 必須 | 文字列または `null` | セッション識別用の短縮 ID。token 本体ではない。 |
| `token_id` | string/null | 必須 | API token id または `null` | API token 関連イベントの対象。 |
| `remote_addr` | string/null | 必須 | IP 文字列または `null` | 接続元。取得不能時は `null`。 |
| `reason` | string/null | 必須 | 文字列または `null` | 失敗理由。秘密情報を含めない。 |

**`.api_access_log` JSON Lines schema：**

各行は認証後 API、認証失敗 API、Webhook API の HTTP 呼び出し 1 件を表す JSON object とする。password、token、Webhook secret、SMTP password、request body の secret 値を保存してはならない。

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `at` | string | 必須 | ISO 8601 | response 送信直前の日時。 |
| `request_id` | string | 必須 | 32 文字 lowercase hex | [`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) で生成する API 呼び出し識別子。 |
| `method` | string | 必須 | HTTP method | `GET` / `POST` / `PUT` / `PATCH` / `DELETE`。 |
| `path` | string | 必須 | `/api/...` | query を含まない path。 |
| `query` | object | 必須 | JSON object | 許可済み query key と値。秘密値は禁止。 |
| `status` | integer | 必須 | HTTP status code | response status。 |
| `duration_ms` | integer | 必須 | 0 以上 | handler 開始から response 確定までのミリ秒。 |
| `auth_type` | string | 必須 | `"session"`, `"api_token"`, `"webhook"`, `"none"` | 認証種別。 |
| `actor` | string/null | 必須 | `"admin"`、token id、`"webhook"`、または `null` | 操作者。token 本体は保存しない。 |
| `remote_addr` | string/null | 必須 | access log は共通 field `request remote_addr` | 接続元。取得不能時は `null`。 |
| `user_agent` | string/null | 必須 | 文字列または `null` | 取得不能時は `null`。 |
| `error` | string/null | 必須 | エラーコードまたは `null` | 成功時は `null`。 |

**`.webhook_events.json` JSON Lines schema：**

`.webhook_events.json` は 1 行 1 event を追記する。mode は `0600` とする。request header 全体、署名値、Webhook secret、payload 全体を保存してはならない。

| key | 型 | 必須 | 許容値 / 説明 |
|-----|----|------|---------------|
| `id` | string | 必須 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の時刻ベース ID 契約に従う `wh` prefix の event id。 |
| `timestamp` | string | 必須 | UTC ISO 8601。 |
| `delivery_id` | string | 必須 | `X-GitHub-Delivery`。1〜200 文字。 |
| `event` | string | 必須 | GitHub event 名。1〜100 文字。 |
| `ref` | string/null | 必須 | push ref または `null`。 |
| `branch` | string/null | 必須 | `refs/heads/` を除いた branch または `null`。 |
| `sha` | string/null | 必須 | push `after`。40 文字 lowercase hex または `null`。 |
| `repository` | string/null | 必須 | `owner/repo` 形式または `null`。 |
| `build_triggered` | boolean | 必須 | queue 追加済みなら `true`。 |
| `queued_id` | string/null | 必須 | queue id または `null`。 |
| `result` | string | 必須 | `"queued"`, `"duplicate"`, `"ignored_event"`, `"ignored_branch"`, `"queue_full"`。 |
| `error_code` | string/null | 必須 | `result="queue_full"` では `"queue_full"`、その他の `result` では `null`。secret、署名、payload 断片を含めない。 |

**`.approval_queue` JSON Lines schema：**

`.approval_queue` は append-only とし、1 行 1 record を追記する。同一 id について schema-valid な物理行のうち file 末尾に最も近い record を最新 record として有効状態にし、それ以前の record は監査履歴として残す。壊れた行と空行は最新判定から除外し、`created_at`、`decided_at`、status の順序から最新状態を推測してはならない。

| キー | 型 | 必須 | 許容値 / 説明 |
|------|----|------|---------------|
| `id` | string | 必須 | base は `appr{YYYYMMDDHHmmss}`。衝突処理は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の時刻ベース ID 契約に従う。 |
| `status` | string | 必須 | `"pending"`、`"approved"`、`"rejected"`、`"expired"`。 |
| `branch` | string | 必須 | branch target 名。 |
| `sha` | string | 必須 | 40 文字 lowercase hex。 |
| `target` | string | 必須 | branch target id または target file。 |
| `requested_trigger` | string | 必須 | `"polling"`, `"force_interval"`, `"manual"`, `"webhook"`, `"local_watch"`。 |
| `requested_force` | boolean | 必須 | SHA 一致時も build を実行する要求なら `true`。`force_interval` と手動 force は `true`、その他は `false`。 |
| `requested_by` | string | 必須 | `"system"`, `"admin"`, `"webhook"`, API token id。 |
| `delivery_id` | string/null | 必須 | webhook delivery id または `null`。 |
| `created_at` | string | 必須 | UTC ISO 8601。 |
| `expires_at` | string | 必須 | UTC ISO 8601。 |
| `decided_at` | string/null | 必須 | approve / reject / expire 時刻または `null`。 |
| `decided_by` | string/null | 必須 | 管理 session は `"admin"`、API token は token id、未決定時は `null`。 |
| `queue_id` | string/null | 必須 | approve で追加した queue id または `null`。 |
| `reason` | string/null | 必須 | reject 理由、expire 理由、または `null`。 |

同一 `id` の後続 record は `branch`、`sha`、`target`、`requested_trigger`、`requested_force`、`requested_by`、`delivery_id`、`created_at`、`expires_at` を pending record から変更せずコピーする。`pending` では `decided_at=null`、`decided_by=null`、`queue_id=null`、`reason=null` とする。`approved` では `decided_at`、`decided_by`、`queue_id` を非 `null`、`reason=null` とする。`rejected` では `decided_at`、`decided_by` を非 `null`、`queue_id=null`、`reason="rejected"` とする。`expired` では `decided_at` を非 `null`、`decided_by="system"`、`queue_id=null`、`reason="expired"` とする。

**`.build_trends.json` schema：**

```json
{
  "schema_version": 1,
  "samples": [
    {
      "build_id": "b20260915100500",
      "finished_at": "2026-09-15T10:05:00Z",
      "branch": "main",
      "trigger": "polling",
      "duration_seconds": 12,
      "status": "success",
      "target_status": "success",
      "anomaly": false
    }
  ],
  "summary": {
    "count": 1,
    "avg_seconds": 12,
    "median_seconds": 12,
    "p95_seconds": 12,
    "anomaly_count": 0
  }
}
```

| キー | 型 | 必須 | 許容値 / 説明 |
|------|----|------|---------------|
| `schema_version` | integer | 必須 | `1` 固定。 |
| `samples` | object[] | 必須 | 0〜`build_trend_keep_count` 件。`finished_at` 昇順。同一 `build_id` は 1 件だけ許可する。 |
| `samples[].build_id` | string | 必須 | build id。 |
| `samples[].finished_at` | string | 必須 | UTC ISO 8601。 |
| `samples[].branch` | string | 必須 | build 対象 branch。 |
| `samples[].trigger` | string | 必須 | `.build_history.trigger` の許容値。 |
| `samples[].duration_seconds` | number | 必須 | 0 以上。 |
| `samples[].status` | string | 必須 | 集計用の正規化結果。`"success"`、`"failure"`、`"cancelled"`。 |
| `samples[].target_status` | string | 必須 | 原因を保持する詳細結果。[runner 結果値 schema](#runner-result-schema) で build log が許可された値だけを許可する。 |
| `samples[].anomaly` | boolean | 必須 | 所要時間異常として確定済みの場合だけ `true`。 |
| `summary` | object | 必須 | `samples` 全件から再計算した集計。部分更新は禁止する。 |
| `summary.count` | integer | 必須 | `samples` 件数と一致する 0 以上の整数。 |
| `summary.avg_seconds` | number/null | 必須 | sample 0 件では `null`。それ以外は小数第 2 位まで。 |
| `summary.median_seconds` | number/null | 必須 | sample 0 件では `null`。それ以外は小数第 2 位まで。 |
| `summary.p95_seconds` | number/null | 必須 | sample 0 件では `null`。 |
| `summary.anomaly_count` | integer | 必須 | `samples[].anomaly == true` の件数。 |

`target_status` から `status` への写像は、`success` と `success_deploy_pending` を `success`、`failure_` prefix と `hook_error` を `failure`、`cancelled` を `cancelled` とする。build log を作成しない approval、skip、config、lock、circuit の詳細結果は trend sample に含めない。

統計算出、同一 build id の置換、保持件数 trim、破損復旧は [`docs/details/runner.md` 詳細本文責務 §27.33](runner.md#sec-27-33) を参照する。`GET /api/stats/build-trends` の query、選択順、response schema、部分集合 summary、不在・破損時の HTTP 応答は [`docs/details/api.md`](api.md) 詳細本文責務の `BuildTrendStats` 契約を正本とする。statefile owner は schema-valid object と read error の返却だけを担当し、API response を組み立てない。

**`.build_chain_config` schema：**

```json
{
  "chains": [
    {
      "id": "build-docs",
      "branch": "main",
      "target_file": "docs/index.md",
      "depends_on": [],
      "required": true,
      "enabled": true
    }
  ]
}
```

| キー | 型 | 必須 | 許容値 / 説明 |
|------|----|------|---------------|
| `chains` | object[] | 必須 | 0〜100 件。配列順は同順位 job の実行順として保持する。 |
| `chains[].id` | string | 必須 | `^[a-zA-Z0-9_-]{1,64}$`。config 内で一意。 |
| `chains[].branch` | string | 必須 | 1〜255 文字の branch 名。前後空白、制御文字、空文字は禁止する。 |
| `chains[].target_file` | string | 必須 | repository root 起点の `/` 区切り相対 path。絶対 path、`..` segment、空文字は禁止する。 |
| `chains[].depends_on` | string[] | 必須 | 0〜20 件。同一 config 内の有効な job id だけを重複なく指定する。自分自身は指定できない。 |
| `chains[].required` | boolean | 必須 | dependency failure 時に当該 job を skip する場合は `true`。 |
| `chains[].enabled` | boolean | 必須 | 実行対象は `true`。enabled job から disabled job への依存は禁止する。 |

未知 key、型不一致、件数超過、id 重複、参照先不在、自己依存、disabled job 依存、循環依存は schema validation failure とし、write adapter は既存 `.build_chain_config` を変更しない。HTTP `422`、`.config_log`、response details は [`docs/details/api.md`](api.md) 詳細本文責務の `POST /api/build-chain-config` 契約を参照する。topological order、required / optional 境界、chain summary は [`docs/details/runner.md` 詳細本文責務 §27.34](runner.md#sec-27-34) を参照する。

<a id="runner-result-schema"></a>
**runner 結果値 schema：**

runner 結果値は保存先ごとに意味を分離する。`.build_logs/{id}.json.status` と `.build_status.json.status` は正規化状態、`.build_logs/{id}.json.target_status` と `.build_history.status` と `.build_status.json.last_target_status` は詳細結果である。同じ field 名へ正規化状態と詳細結果を混在させてはならない。

| 詳細結果 | build log | history | status summary | 意味 |
|----------|-----------|---------|----------------|------|
| `success` | 許可 | 許可 | 許可 | build と必須保存が成功した。 |
| `success_deploy_pending` | 許可 | 許可 | 許可 | build は成功し、1 件以上の deploy が pending である。 |
| `failure_api` | 許可 | 許可 | 許可 | GitHub API の最終失敗。 |
| `failure_decode` | 許可 | 許可 | 許可 | blob decode または source materialize 失敗。 |
| `failure_precheck` | 許可 | 許可 | 許可 | build 前検証失敗。 |
| `failure_build` | 許可 | 許可 | 許可 | builder または必須 pipeline step の失敗。 |
| `failure_state_write` | 許可 | 許可 | 許可 | 必須状態保存失敗。 |
| `failure_target_missing` | 許可 | 許可 | 許可 | 必須 target 不在。 |
| `failure_pipeline_config` | 許可 | 許可 | 許可 | pipeline 設定不正または読取不能。 |
| `failure_timeout` | 許可 | 許可 | 許可 | build、step、remote build の timeout。 |
| `failure_remote_build` | 許可 | 許可 | 許可 | remote build または artifact 検証失敗。 |
| `failure_tag_rule` | 許可 | 許可 | 許可 | 保存済み自動 tag rule の parse または評価失敗。 |
| `hook_error` | 許可 | 許可 | 許可 | pre hook により build を中止した。 |
| `cancelled` | 許可 | 許可 | 許可 | 実行中 build を明示的に中止した。 |
| `approval_rejected` | 禁止 | 許可 | 禁止 | approval record の却下。build log は作成しない。 |
| `approval_expired` | 禁止 | 許可 | 禁止 | approval record の期限切れ。build log は作成しない。 |
| `skipped_dependency_failed` | 禁止 | 許可 | 禁止 | chain dependency failure による未実行 job。build log は作成しない。 |
| `skipped_no_change` | 禁止 | 禁止 | 許可 | SHA 差分なし。 |
| `skipped_cooldown` | 禁止 | 禁止 | 許可 | cooldown 中。 |
| `skipped_tag_filter` | 禁止 | 禁止 | 許可 | tag filter 不一致。 |
| `skipped_maintenance` | 禁止 | 禁止 | 許可 | maintenance mode 中。 |
| `circuit_open` | 禁止 | 禁止 | 許可 | circuit breaker open。 |
| `config_recovered` | 禁止 | 禁止 | 許可 | 起動時設定復旧後に build 未実行。 |
| `config_error` | 禁止 | 禁止 | 許可 | 起動時設定不正により build 未実行。 |
| `lock_skipped` | 禁止 | 禁止 | 許可 | 有効な実行中 lock により起動を skip。 |

`failure` のような詳細段階を失う総称を `.build_history.status` または `.build_logs/{id}.json.target_status` へ保存してはならない。正規化が必要な API、集計、UI は [runner 結果値 schema](#runner-result-schema) の詳細結果から `success`、`failure`、`cancelled`、`skipped` へ写像する。

<a id="build-history-schema"></a>
**`.build_history` JSON Lines schema：**

各行は [`.build_history` JSON Lines schema](#build-history-schema) の JSON object とする。未知 key は禁止する。build log を伴う行は同じ `id` の log の要約とし、approval または chain skip の行は [`.build_history` JSON Lines schema](#build-history-schema) の nullable 値を使用する。

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `id` | string | 必須 | build id または approval id | `.build_history` の schema-valid 行全体で一意な build、approval、chain job ID。 |
| `status` | string | 必須 | [runner 結果値 schema](#runner-result-schema) で history が許可された値 | 詳細結果。 |
| `trigger` | string | 必須 | `"polling"`, `"force_interval"`, `"manual"`, `"webhook"`, `"retry_pending_transfer"`, `"rollback"`, `"local_watch"`, `"approval"` | 起動種別。 |
| `branch` | string/null | 必須 | branch 名または `null` | approval / build target の branch。取得不能時だけ `null`。 |
| `target_file` | string/null | 必須 | 相対 path または `null` | 対象 file または directory。 |
| `started_at` | string/null | 必須 | UTC ISO 8601 または `null` | build 未実行の approval event は `null`。 |
| `finished_at` | string | 必須 | UTC ISO 8601 | build 完了または event 確定日時。API の `build_at` はこの値から算出する。 |
| `duration_seconds` | integer/null | 必須 | 0 以上または `null` | build 未実行 event は `null`。 |
| `commit_sha` | string/null | 必須 | 40 文字 lowercase hex または `null` | 対象 commit SHA。 |
| `blob_sha` | string/null | 必須 | 40 文字 lowercase hex または `null` | 対象 blob SHA。API の `sha` は `commit_sha`、`blob_sha` の順で算出する。 |
| `pages` | integer/null | 必須 | 0 以上または `null` | report 不在は `null`。 |
| `warnings` | integer | 必須 | 0 以上 | warning 件数。 |
| `error` | string/null | 必須 | 500 文字以内または `null` | 固定エラー文言。 |
| `output_size_bytes` | integer/null | 必須 | 0 以上または `null` | 成果物サイズ。 |
| `output_sha256` | string/null | 必須 | 64 文字 lowercase hex または `null` | 成果物 manifest SHA-256。 |
| `size_warn` | boolean | 必須 | boolean | サイズ警告。 |
| `retry_count` | integer | 必須 | 0 以上 | 初回以後の追加 retry 回数。 |
| `commit_status_state` | string/null | 必須 | `"pending"`, `"success"`, `"failure"`, `"error"`, `null` | 最終 GitHub Commit Status。 |
| `flagged` | boolean | 必須 | boolean | 重要フラグ。 |
| `tags` | string[] | 必須 | 0〜100 件、重複禁止 | 手動 / 自動タグ。 |
| `comment` | string/null | 必須 | 2000 文字以内または `null` | 管理コメント。 |
| `rollback_from` | string/null | 必須 | build id または `null` | rollback 元 build id。 |
| `snapshot_id` | string/null | 必須 | build id または `null` | 使用または生成した snapshot id。 |
| `failure_category` | string/null | 必須 | 読取時は FailureCategory 値、前方互換 category id、または `null`。保存時は FailureCategory 値または `null`。 | `status="success_deploy_pending"` は `"deploy_failure"`。失敗結果は runner の分類値。`status="success"`、cancel、skip、approval event は `null`。 |
| `chain_run_id` | string/null | 必須 | chain run id または `null` | chain 外は `null`。 |

`.build_history.failure_category` の前方互換 category id は正規表現 `^[a-z][a-z0-9_]{0,63}$` に一致し、現行 FailureCategory 値に一致しない保存済み文字列とする。`readBuildHistory()` はこの値を破損として除外せず原値のまま返し、書き換え、現行値への推測変換、`"unknown"` への正規化を行わない。write adapter と runner は前方互換 category id を新規保存してはならず、現行 FailureCategory 値または `null` だけを保存する。この読取例外は `.build_history.failure_category` だけに適用し、他 field、`.build_logs/{id}.json.failure_category`、未知 key の許容へ拡張しない。API warning と filter は [`docs/details/api.md`](api.md) 詳細本文責務の `HistoryPageObject` 契約を参照する。

build log を伴わない history record は次の値を固定する。表にない field は `.build_history` schema の型を満たす空値に固定し、実装者判断で省略してはならない。

| record | `id` | `trigger` | `started_at` | `finished_at` | `duration_seconds` | SHA | 出力 / report | error | chain |
|--------|------|-----------|--------------|---------------|--------------------|-----|-----------------|-------|-------|
| `approval_rejected` | approval id | `approval` | `null` | `decided_at` | `null` | `commit_sha:null`, `blob_sha:approval.sha` | `pages:null`, `warnings:0`, output 値 `null`, `size_warn:false` | `"approval rejected"` | `null` |
| `approval_expired` | approval id | `approval` | `null` | expire 確定時刻 | `null` | `commit_sha:null`, `blob_sha:approval.sha` | `pages:null`, `warnings:0`, output 値 `null`, `size_warn:false` | `"approval expired"` | `null` |
| `skipped_dependency_failed` | 新規 build id | 元の起動 trigger | `null` | skip 確定時刻 | `null` | 取得済み値または `null` | `pages:null`, `warnings:0`, output 値 `null`, `size_warn:false` | `"dependency failed"` | `chain_run_id` 必須 |

`approval_rejected`、`approval_expired`、`skipped_dependency_failed` の `retry_count` は `0`、`commit_status_state`、`rollback_from`、`snapshot_id`、`failure_category`、`comment` は `null`、`flagged` は `false`、`tags` は空配列とする。approval record の `branch` と `target_file` は `.approval_queue` の `branch` と `target`、chain skip の値は ChainJob から取得する。対応する `.build_logs/{id}.json` がないことは正常である。

**`.build_logs/{id}.json` schema：**

最終保存時は以下の全必須 key を 1 object に含める。未知 key、history 専用結果、summary 専用結果は保存してはならない。

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `id` | string | 必須 | build id | ファイル名 `{id}.json` と一致する。 |
| `status` | string | 必須 | `"running"`, `"success"`, `"failure"`, `"cancelled"` | 正規化状態。 |
| `target_status` | string/null | 必須 | [runner 結果値 schema](#runner-result-schema) で build log が許可された値または `null` | 詳細結果。途中保存だけ `null`、最終保存では必須 string。 |
| `trigger` | string | 必須 | `.build_history.trigger` と同じ | 起動種別。 |
| `trigger_actor` | string/null | 必須 | token id、`"admin"`、`"system"`、または `null` | 起動主体。secret 値は禁止する。 |
| `branch` | string | 必須 | branch 名 | 対象 branch。 |
| `target_file` | string | 必須 | 相対 path | 対象 file または directory。 |
| `started_at` | string | 必須 | UTC ISO 8601 | 開始日時。 |
| `finished_at` | string/null | 必須 | UTC ISO 8601 または `null` | 完了前は `null`。 |
| `duration_seconds` | integer/null | 必須 | 0 以上または `null` | 完了前は `null`。 |
| `blob_sha` | string/null | 必須 | 40 文字 lowercase hex または `null` | 対象 blob SHA。 |
| `previous_blob_sha` | string/null | 必須 | 40 文字 lowercase hex または `null` | 初回または取得不能は `null`。 |
| `commit_sha` | string/null | 必須 | 40 文字 lowercase hex または `null` | commit SHA。 |
| `commit_message` | string/null | 必須 | 1000 文字以内または `null` | commit message。 |
| `commit_author` | string/null | 必須 | 255 文字以内または `null` | commit author。 |
| `commit_at` | string/null | 必須 | UTC ISO 8601 または `null` | commit 日時。 |
| `pipeline` | object | 必須 | Pipeline object | 標準 builder command または pipeline の実行結果。未実行時も null field を持つ object。 |
| `pipeline_steps` | object[] | 必須 | PipelineStep object、定義順 | 標準 builder command だけを使用した場合は空配列。 |
| `report` | object/null | 必須 | Report object または `null` | `[REPORT]` 解析結果。 |
| `warnings` | string[] | 必須 | 0 件以上 | warning code / message。 |
| `deploy` | object[] | 必須 | Deploy object、設定順 | deploy 未実行は空配列。 |
| `target_results` | object[] | 必須 | TargetResult object、設定順 | parallel deploy 未使用は空配列。 |
| `changed_targets` | object[] | 必須 | ChangedTarget object、target path 辞書順 | 単一 target も同じ形式で保存する。 |
| `matched_tags` | string[] | 必須 | 最大 100 件 | tag filter 無効または一致 tag なしは空配列。 |
| `skipped_hook_ids` | string[] | 必須 | hook id 昇順 | disabled hook なしは空配列。 |
| `attempts` | object[] | 必須 | Attempt object、実行順 | retry attempt が未開始なら空配列。 |
| `retry_count` | integer | 必須 | 0 以上 | 追加 retry 回数。 |
| `remote_build` | object/null | 必須 | RemoteBuild object または `null` | local build は `null`。 |
| `chain` | object/null | 必須 | Chain object または `null` | chain 外は `null`。 |
| `chain_summary` | object/null | 必須 | ChainSummary object または `null` | chain 最終 job 以外は `null`。 |
| `snapshot_id` | string/null | 必須 | build id または `null` | snapshot 未作成は `null`。 |
| `rollback_from` | string/null | 必須 | build id または `null` | rollback 以外は `null`。 |
| `output_size_bytes` | integer/null | 必須 | 0 以上または `null` | 成果物サイズ。 |
| `output_sha256` | string/null | 必須 | 64 文字 lowercase hex または `null` | 成果物 manifest SHA-256。 |
| `size_warn` | boolean | 必須 | boolean | サイズ警告。 |
| `transfer_verified` | boolean/null | 必須 | boolean または `null` | deploy 未実行は `null`。 |
| `commit_status` | object/null | 必須 | CommitStatus object または `null` | GitHub Commit Status 無効時は `null`。 |
| `build_meta` | object/null | 必須 | BuildMeta object または `null` | builder 出力未生成は `null`。 |
| `failure_category` | string/null | 必須 | FailureCategory 値または `null` | `target_status="success_deploy_pending"` は `"deploy_failure"`。`status="failure"` は runner の分類値。その他は `null`。 |
| `failure_evidence` | object[] | 必須 | FailureEvidence object、最大 10 件 | `status="failure"` または `target_status="success_deploy_pending"` で evidence がある場合だけ保存し、その他または evidence なしは空配列。 |
| `environment` | object | 必須 | Environment object | build id 採番直後の実行環境。 |
| `error` | string/null | 必須 | 500 文字以内または `null` | 固定エラー文言。 |
| `comment` | string/null | 必須 | 2000 文字以内または `null` | 管理コメント。 |
| `flagged` | boolean | 必須 | boolean | 重要フラグ。 |
| `tags` | string[] | 必須 | 0〜100 件、重複禁止 | 手動 / 自動タグ。 |

途中保存は `status:"running"`、`target_status:null`、`finished_at:null`、`duration_seconds:null` の組合せに限る。最終保存は `status` を `"success"`、`"failure"`、`"cancelled"` のいずれかに確定し、`target_status` を非 `null`、`finished_at >= started_at`、`duration_seconds >= 0` とする。これら以外の status / nullable 組合せは build log 破損とする。dry-run は build log 自体を作成しない。

Report object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `pages` | integer | 必須 | 出力ページ数。 |
| `headings` | integer | 必須 | 見出し数。 |
| `tables_count` | integer | 必須 | テーブル数。 |
| `code_blocks_count` | integer | 必須 | コードブロック数。 |
| `warnings_count` | integer | 必須 | 警告件数。 |
| `size_warn` | boolean | 必須 | 出力サイトサイズ警告。 |
| `broken_links` | integer | 必須 | 内部リンク不整合数。 |
| `heading_skips` | integer | 必須 | 見出しレベルスキップ数。 |
| `reading_time` | integer | 必須 | 推計読了時間。 |
| `theme` | string | 必須 | 使用 theme 名。 |
| `build_id` | string | 必須 | `[REPORT] build_id`。未指定時は空文字。 |
| `commit_sha` | string | 必須 | `[REPORT] commit_sha`。未指定時は空文字。 |
| `build_at` | string | 必須 | `[REPORT] build_at`。未指定時は空文字。 |

Attempt object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `attempt` | integer | 必須 | 初回は `1`。retry ごとに +1。 |
| `started_at` | string | 必須 | UTC ISO 8601。 |
| `finished_at` | string/null | 必須 | 完了時刻。dry-run は build log を作成しないため、本 field を保存しない。 |
| `stage` | string | 必須 | `"github"`, `"pipeline"`, `"remote_build"`, `"deploy"` のいずれか。dry-run は build log を作成しないため `"dry_run"` stage を新規保存してはならない。 |
| `status` | string | 必須 | `"success"` または `"failure"`。 |
| `retryable` | boolean | 必須 | この失敗が retry 対象か。成功時は `false`。 |
| `error` | string/null | 必須 | 失敗理由。成功時は `null`。 |

CommitStatus object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `enabled` | boolean | 必須 | `.server_config.commit_status_enabled` の評価結果。 |
| `state` | string/null | 必須 | `"pending"`, `"success"`, `"failure"`, `"error"`、未送信時 `null`。 |
| `context` | string | 必須 | 送信 context。 |
| `target_url` | string/null | 必須 | 送信 target_url。省略時 `null`。 |
| `sent_at` | string/null | 必須 | 最終送信時刻。未送信時 `null`。 |
| `http_status` | integer/null | 必須 | GitHub API HTTP status。未送信時 `null`。 |
| `error` | string/null | 必須 | 送信失敗理由。成功時 `null`。 |

BuildMeta object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `build_id` | string | 必須 | HTML meta `adlaire-build-id` と同じ値。 |
| `commit_sha` | string | 必須 | HTML meta `adlaire-commit-sha` と同じ値。 |
| `build_at` | string | 必須 | HTML meta `adlaire-build-at` と同じ値。 |

Pipeline object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `exit_code` | integer/null | 必須 | process 未実行または timeout は `null`。 |
| `stdout` | string | 必須 | LF 正規化、secret mask 済み、最大 64 KiB。 |
| `stderr` | string | 必須 | LF 正規化、secret mask 済み、最大 64 KiB。 |
| `stdout_truncated` | boolean | 必須 | stdout 切り詰め時だけ `true`。 |
| `stderr_truncated` | boolean | 必須 | stderr 切り詰め時だけ `true`。 |

PipelineStep object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `index` | integer | 必須 | `.pipeline.yml` または inline YAML の 0 始まり定義順。 |
| `name` | string | 必須 | 1〜100 文字。 |
| `required` | boolean | 必須 | 必須 step なら `true`。 |
| `status` | string | 必須 | `"success"`, `"failure"`, `"optional_failed"`, `"timeout"`, `"not_run"`。 |
| `started_at` | string/null | 必須 | 未実行は `null`。 |
| `finished_at` | string/null | 必須 | 未実行は `null`。 |
| `duration_seconds` | integer/null | 必須 | 未実行は `null`、それ以外は 0 以上。 |
| `exit_code` | integer/null | 必須 | timeout / 未実行は `null`。 |
| `stdout` | string | 必須 | secret mask 済み、最大 64 KiB。未実行は空文字。 |
| `stderr` | string | 必須 | secret mask 済み、最大 64 KiB。未実行は空文字。 |
| `truncated` | boolean | 必須 | stdout または stderr を切り詰めた場合だけ `true`。 |
| `error` | string/null | 必須 | 固定エラー文言または `null`。 |

HookLog object:

| キー | 型 | 必須 | 許容値 / 説明 |
|------|----|------|---------------|
| `hook_id` | string | 必須 | `^[A-Za-z0-9_-]{1,64}$`。file 名の `{hook_id}` と一致する。 |
| `build_id` | string | 必須 | build id。file 名の `{build_id}` と一致する。 |
| `phase` | string | 必須 | `"pre"` または `"post"`。 |
| `status` | string | 必須 | `"success"`、`"failure"`、`"timeout"` のいずれか。 |
| `started_at` | string | 必須 | UTC ISO 8601。 |
| `finished_at` | string | 必須 | UTC ISO 8601。`finished_at >= started_at`。 |
| `duration_seconds` | integer | 必須 | 0 以上。負値になる clock drift は `0` に丸める。 |
| `stdout` | string | 必須 | LF 正規化、secret mask 済み、UTF-8 byte 数で最大 65536 bytes。 |
| `stderr` | string | 必須 | LF 正規化、secret mask 済み、UTF-8 byte 数で最大 65536 bytes。 |
| `exit_code` | integer/null | 必須 | process が終了コードを返した場合は integer。timeout、起動失敗、signal 終了は `null`。 |
| `timed_out` | boolean | 必須 | timeout の場合だけ `true`。`true` の場合は `status="timeout"`、`exit_code=null`。 |
| `truncated` | boolean | 必須 | stdout または stderr を 65536 bytes で切り詰めた場合だけ `true`。 |

HookLog object は未知 key を禁止する。secret mask を先に適用し、その後に UTF-8 code point を分断しない位置で stdout / stderr をそれぞれ 65536 bytes 以下へ切り詰める。hook の実行順、process kill、pre abort、post failure、保存タイミングは [`docs/details/runner.md` 詳細本文責務 §27.27](runner.md#sec-27-27) を参照する。

Deploy object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `target_id` | string | 必須 | 設定上の target id。未指定時は `{branch_index}-{target_index}`。 |
| `host` | string | 必須 | deploy host。credential を含めない。 |
| `user` | string | 必須 | deploy user。 |
| `dest_dir` | string | 必須 | deploy destination。 |
| `status` | string | 必須 | `"success"`, `"failure"`, `"pending"`。 |
| `transfer_verified` | boolean | 必須 | remote checksum 検証成功時だけ `true`。 |
| `files_total` | integer | 必須 | 0 以上。 |
| `files_uploaded` | integer | 必須 | 0 以上かつ `files_total` 以下。 |
| `files_skipped` | integer | 必須 | 0 以上かつ `files_total` 以下。 |
| `bytes_uploaded` | integer | 必須 | 0 以上。 |
| `error` | string/null | 必須 | 固定エラー文言または `null`。 |

ChangedTarget object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `target_file` | string | 必須 | 正規化済み相対 path。 |
| `before_sha` | string/null | 必須 | 変更前 SHA。不明は `null`。 |
| `after_sha` | string/null | 必須 | 変更後 SHA。missing / error は `null`。 |
| `source` | string | 必須 | `"github"` または `"local"`。 |
| `result` | string | 必須 | `"changed"`, `"unchanged"`, `"missing"`, `"error"`。 |

TargetResult object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `target_id` | string | 必須 | Deploy object と同じ target id。 |
| `status` | string | 必須 | `"success"` または `"failure"`。 |
| `started_at` | string | 必須 | UTC ISO 8601。 |
| `finished_at` | string | 必須 | UTC ISO 8601。 |
| `error_code` | string/null | 必須 | `"deploy_timeout"`, `"deploy_ssh_error"`, `"deploy_checksum_error"`, `"deploy_internal_error"`, `null`。`status="success"` では `null`、`status="failure"` では非 `null` を必須とする。 |
| `error` | string/null | 必須 | `error_code` が `null` なら `null`。非 `null` ではそれぞれ `"deploy timeout"`, `"deploy ssh failed"`, `"deploy checksum failed"`, `"deploy internal error"` とする。 |

RemoteBuild object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `host` | string | 必須 | remote host。 |
| `user` | string | 必須 | remote user。 |
| `work_dir_basename` | string | 必須 | remote work directory の basename だけを保存する。 |
| `command_name` | string | 必須 | argv 先頭要素の basename。引数は保存しない。 |
| `exit_code` | integer/null | 必須 | timeout / 接続失敗は `null`。 |
| `duration_seconds` | integer | 必須 | 0 以上。 |
| `artifact_size_bytes` | integer/null | 必須 | artifact 未取得は `null`。 |
| `manifest_file_count` | integer/null | 必須 | manifest 未検証は `null`。 |
| `status` | string | 必須 | `"success"` または `"failure"`。 |
| `error` | string/null | 必須 | 固定エラー文言または `null`。 |

Chain object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `chain_run_id` | string | 必須 | base は `chain{YYYYMMDDHHmmss}`。衝突処理は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の時刻ベース ID 契約に従う。 |
| `job_id` | string | 必須 | `.build_chain_config.chains[].id`。 |
| `depends_on` | string[] | 必須 | 設定値と同じ順序。 |
| `chain_index` | integer | 必須 | topological 実行順の 0 始まり index。 |

ChainSummary object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `chain_run_id` | string | 必須 | Chain object と同じ値。 |
| `total_jobs` | integer | 必須 | enabled job 件数。 |
| `success_count` | integer | 必須 | `success` / `success_deploy_pending` 件数。 |
| `failure_count` | integer | 必須 | failure / hook / cancel 件数。 |
| `skipped_count` | integer | 必須 | `skipped_dependency_failed` 件数。 |

FailureCategory 値:

`"github_api"`, `"pipeline_timeout"`, `"pipeline_exit"`, `"deploy_failure"`, `"hook_error"`, `"config_error"`, `"resource_error"`, `"unknown"` のいずれかとする。

runner と write adapter が保存できる値はこの固定列挙だけとする。`.build_history` の前方互換 category id は保存許可値ではなく、既存履歴を消失させず読み取るための明示的な read-only 例外である。

FailureEvidence object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `source` | string | 必須 | `"github_api"`, `"pipeline"`, `"deploy"`, `"hook"`, `"config"`, `"resource"`, `"runner"`。 |
| `code` | string | 必須 | 同じ build log の `failure_category` と完全一致する FailureCategory 値。 |
| `message` | string | 必須 | secret と入力値連結を含まない固定文言、最大 300 文字。 |
| `at` | string | 必須 | UTC ISO 8601。 |

Environment object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `os` | string | 必須 | `runtime.GOOS`。 |
| `arch` | string | 必須 | `runtime.GOARCH`。 |
| `go_version` | string | 必須 | `runtime.Version()`。 |
| `runner_version` | string | 必須 | build info main version または `"unknown"`。 |
| `builder_version` | string | 必須 | builder version token または `"unknown"`。 |
| `hostname` | string | 必須 | 最大 255 文字。取得不能は `"unknown"`。 |
| `state_dir` | string | 必須 | home 配下は basename、それ以外は絶対 path。 |
| `disk_free_bytes` | integer/null | 必須 | 0 以上または取得不能時 `null`。 |
| `captured_at` | string | 必須 | UTC ISO 8601。 |
| `env_keys` | string[] | 必須 | branch / pipeline から注入した key の ASCII 昇順。値は保存禁止。 |

<a id="build-status-schema"></a>
**`.build_status.json` schema：**

`.build_status.json` は runner の現在状態と直近結果を 1 ファイルで読むための要約 schema である。API endpoint ごとの読取順と status 不在時集約値の算出は [`docs/details/api.md` 詳細本文責務 §22.0c.1](api.md#sec-22-0c-1) を参照する。UI 表示は [`docs/details/ui.md` 詳細本文責務 §24](ui.md#24-標準管理ツール-仕様) を参照する。MCP の現在状態は [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務を参照する。

```json
{
  "schema_version": 1,
  "updated_at": "2026-09-16T01:00:12Z",
  "status": "success",
  "running": false,
  "current_build_id": null,
  "last_build_id": "b20260916010000",
  "last_trigger": "polling",
  "last_target_status": "success",
  "last_branch": "main",
  "last_target_file": "docs",
  "last_blob_sha": "0123456789abcdef0123456789abcdef01234567",
  "last_commit_sha": "89abcdef0123456789abcdef0123456789abcdef",
  "last_started_at": "2026-09-16T01:00:00Z",
  "last_finished_at": "2026-09-16T01:00:12Z",
  "last_duration_seconds": 12,
  "last_error": null,
  "last_deploy_at": "2026-09-16T01:00:12Z",
  "last_deploy_status": "success",
  "pending_transfers_count": 0,
  "notify_pending_count": 0,
  "circuit_open": false,
  "circuit_consecutive_failures": 0,
  "output_sha256": null,
  "size_warn": false
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `schema_version` | integer | 必須 | `1` 固定 | schema version。 |
| `updated_at` | string/null | 必須 | UTC ISO 8601 または `null` | 初期値のみ `null`。 |
| `status` | string | 必須 | `"none"`, `"running"`, `"success"`, `"failure"`, `"cancelled"`, `"skipped"`, `"warning"` | 正規化状態。`"none"` は初期値のみ。 |
| `running` | boolean | 必須 | boolean | runner が build 処理中なら `true`。 |
| `current_build_id` | string/null | 必須 | build id または `null` | 実行中 build id。実行中でなければ `null`。 |
| `last_build_id` | string/null | 必須 | build id または `null` | 最後に status summary 対象の build log / history を確定した build id。history-only 行の追記では更新しない。未実行なら `null`。 |
| `last_trigger` | string/null | 必須 | `.build_history.trigger` の値、`"startup_config_integrity"`、または `null` | 起動種別。startup integrity だけで終了した場合も保存できる。 |
| `last_target_status` | string/null | 必須 | [runner 結果値 schema](#runner-result-schema) で status summary が許可された値または `null` | 詳細結果。初回のみ `null`。 |
| `last_branch` | string/null | 必須 | branch 名または `null` | 最終対象 branch。 |
| `last_target_file` | string/null | 必須 | 相対 path または `null` | 最終対象 file。 |
| `last_blob_sha` | string/null | 必須 | Git blob SHA または `null` | 取得不能時は `null`。 |
| `last_commit_sha` | string/null | 必須 | Git commit SHA または `null` | 取得不能時は `null`。 |
| `last_started_at` | string/null | 必須 | UTC ISO 8601 または `null` | 最終開始時刻。 |
| `last_finished_at` | string/null | 必須 | UTC ISO 8601 または `null` | 最終終了時刻。 |
| `last_duration_seconds` | integer/null | 必須 | 0 以上または `null` | 最終所要時間。 |
| `last_error` | string/null | 必須 | 文字列または `null` | 成功・通常 skip は `null`。 |
| `last_deploy_at` | string/null | 必須 | UTC ISO 8601 または `null` | 最後に deploy attempt の結果が確定した時刻。deploy 未実行の初期値は `null`。 |
| `last_deploy_status` | string/null | 必須 | `"success"`, `"failure"`, `"pending"`, `"skipped"`, `"none"`, `null` | 最終 deploy 状態。 |
| `pending_transfers_count` | integer | 必須 | 0 以上 | pending transfer 件数。 |
| `notify_pending_count` | integer | 必須 | 0 以上 | pending notification 件数。 |
| `circuit_open` | boolean | 必須 | boolean | `.build_circuit_state.open` と一致する。 |
| `circuit_consecutive_failures` | integer | 必須 | 0 以上 | `.build_circuit_state.consecutive_failures` と一致する。 |
| `output_sha256` | string/null | 必須 | SHA-256 hex または `null` | 直近成功成果物の manifest SHA-256。 |
| `size_warn` | boolean | 必須 | boolean | 直近 report の size warning。 |

`.build_status.json` の更新タイミング、各 `status` の選択条件、書き込み失敗時の runner 終了コードは [`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) の build status 更新契約を参照する。

<a id="sec-22-0s"></a>
**[§22.0s 状態ファイル実装確認固定契約](statefile.md#sec-22-0s)：**

`statefile` owner component は、[`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a)〜[§22.0c](statefile.md#sec-22-0c) の schema、adapter、更新手順、破損時処理を実装単位として扱う。状態ファイルごとの暗黙処理は追加せず、[`docs/details/statefile.md` 詳細本文責務 §22.0s](statefile.md#sec-22-0s) の固定表の共通契約を満たす。

| 観点 | 入力 | 必須処理 | 成功時出力 | 失敗時出力 / 副作用 |
|------|------|----------|------------|---------------------|
| path 解決 | state dir、[`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の状態ファイル path | state dir 外へ出る path、absolute user input、`..`、symlink 経由の secret 参照を拒否する。 | 正規化済み target path。 | 呼び出し元へ path failure を返す。target を作成・変更しない。 |
| schema 読取 | JSON object / JSON array / JSON Lines / text | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の key、型、nullable、列挙値、UTC 時刻形式を検証する。 | typed value。 | `ErrStateCorrupted` または行単位 skip。未知 key を削除して成功扱いにしない。 |
| 初期値 | file 不在 | [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の「不在時」または初期値を返す。 | 初期 typed value。 | read-only 呼び出しでは file を作成しない。write 呼び出しだけ更新手順で作成する。 |
| atomic write | 更新後 JSON / text | `{name}.lock`、tmp、chmod、rename、file sync、parent sync、lock 削除を順に行う。 | target が完全な新内容へ置換される。 | rename 前失敗は旧 target を維持し、状態ファイル更新手順の rename 前 cleanup 固定契約を適用する。rename 後の file sync、parent sync、lock 削除失敗は新 target を維持し、`STATE_WRITE_AFTER_RENAME_FAILED` と post-rename partial write failure を返す。 |
| lock timeout | 既存 `{name}.lock` | 100ms 間隔、最大 10 秒待つ。 | lock 取得後に更新継続。 | 呼び出し元へ conflict failure を返す。target を変更しない。 |
| chmod | target file / directory | runtime 状態 file と lock file に `0600`、statefile 責務の runtime 状態 directory に `0700` を適用する。 | mode が固定値に一致する。 | chmod 失敗は成功扱いにしない。rename 前なら target 変更なし、rename 後なら ERROR ログへ記録する。 |
| corrupt backup | [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) で退避指定された破損 file | `{name}.corrupt.{YYYYMMDDHHMMSS}.bak` へ同一 directory 内で rename する。 | backup file と再生成初期値。 | backup 失敗時は再生成せず backup failure を返す。secret 内容を log / 呼び出し元の公開値に含めない。 |
| JSON Lines | JSON Lines file | 空行、JSON object 以外、必須 key 不足、型不一致行を除外する。 | 有効行だけの配列。 | 壊れた行は server log に固定 code、path、line number だけ記録する。呼び出し元の公開値に壊れた行数を含めない。 |
| 複数ファイル更新 | 複数 state 書込 caller | 呼び出し元が定義する Write 列順に 1 file ずつ atomic write する。API 固有の順序は [`docs/details/api.md` 詳細本文責務 §22.0d](api.md#sec-22-0d) 以降を参照する。 | 全対象が順に更新される。 | 未処理 file は変更しない。更新済み file は自動 rollback しない。caller 固有 log は caller の詳細本文責務に定義がある場合だけ追記し、statefile owner は代行追記しない。 |

**read / write 境界固定：**

| 呼び出し種別 | 許可する処理 | 禁止する処理 |
|--------------|--------------|--------------|
| read-only caller | 既存 target の読取、typed value 変換、JSON Lines の有効行抽出、fallback 値算出。 | file 作成、chmod、corrupt backup、旧形式保存、lock 待機、lock 削除、tmp 作成。 |
| write caller | 入力 validation 後の atomic write、[`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) で定義された初期値作成、定義済み corrupt backup。 | validation 前の状態変更、未定義 file 作成、未知 key 保存、secret 平文 log。 |
| runner write | build / rollback lifecycle に必要な state 更新、build / rollback / snapshot delete 排他の `.build_lock` 取得・stale 再確認後の削除・所有確認付き解放。 | API 専用 credentials / token / auth state の直接変更。 |
| setup write | 初期配置に必要な `.github_token`、`.last_sha`、admin directory の配置。 | API runtime state、history、build log、session、token の生成。 |

**runner 連動状態更新固定ゲート：**

runner / archive / commitstatus / security / api が同じ実装変更で状態更新を組み合わせる場合でも、statefile owner の契約は [`docs/details/statefile.md` 詳細本文責務 §22.0s](statefile.md#sec-22-0s) の固定表で固定する。呼び出し元 component は業務判断を持ち、statefile は path、schema、lock、atomic write、JSON Lines、破損時処理だけを担当する。

| ゲート | statefile 側の固定処理 | 呼び出し元が渡す値 | 失敗時境界 |
|--------|------------------------|-------------------|------------|
| schema precheck | 保存前に [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の key、型、nullable、enum、UTC 時刻、配列要素 schema を検証する。 | 保存済みとして確定した typed value。 | schema 不一致は target 変更なしで `ErrStateCorrupted` または validation error を返す。 |
| unknown key rejection | 既存 file と新規 value の両方で未知 key を拒否する。例外は [`docs/details/statefile.md` 詳細本文責務 §22.0s](statefile.md#sec-22-0s) に明記済みの旧形式正規化だけ。 | 表示用 key、SDK 用 key、fixture 用 key を含まない object。 | 未知 key を削除して保存しない。既存未知 key も暗黙修復しない。 |
| write order evidence | 複数 file 更新では呼び出し元が決めた順に 1 file ずつ atomic write し、fixture の `write_order` と一致させる。 | 順序付き write plan。 | 失敗地点以降は実行しない。成功済み file は statefile が rollback しない。 |
| JSON Lines append | append 対象は 1 行 1 JSON object とし、末尾 newline を固定する。 | 1 record の typed value。 | append 失敗は対象操作へ返し、既存行の rewrite、sort、修復をしない。 |
| no mutation read | read-only adapter は fallback 値を返すだけで、file 作成、chmod、backup、lock 削除、旧形式保存を行わない。 | 読取対象 path と fallback 条件。 | 読取失敗は typed error を返し、filesystem 差分なし。 |
| corrupt handling | [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) に再生成指定がある file だけ backup → 初期値再生成を許可する。 | 破損判定結果と対象 path。 | backup 失敗時は再生成しない。再生成指定がない file は変更しない。 |
| state path mode | runtime 状態 file と lock file は `0600`、statefile 責務の runtime 状態 directory は `0700` に固定する。 | path の状態責務への帰属判定。 | chmod 失敗を成功扱いにせず、secret 内容を log / 呼び出し元の公開値 / fixture expected に出さない。 |

状態ファイル fixture 名、初期状態、操作、expected、合格条件、実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) statefile owner fixture 固定契約を正本とする。[`docs/details/statefile.md`](statefile.md) 詳細本文責務では、schema、atomic write、lock、JSON Lines、破損時処理、保存順、read-only no mutation の実装契約だけを扱う。
