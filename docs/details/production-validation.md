# Adlaire CI — 本番検証・VPS Pull Bootstrap 詳細仕様

[`docs/details/production-validation.md`](production-validation.md) は、本番検証詳細本文責務として、Phase 17 の ConoHa VPS 試験本番運用と Phase 18 の VPS Pull Bootstrap に関する、provider target、minimum plan class、opaque environment identity、credential / SSH 入力境界、Phase 17 試験本番運用窓、Phase 18 bootstrap session 停止条件、destructive operation 禁止境界、運用中バグ修正順序、bootstrap session、bootstrap 証跡、fixture 証跡への接続条件だけを定義する。状態、Phase、実装可否は [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務、方針とポリシーは [`docs/SPEC.md`](../SPEC.md) 方針責務・ポリシー責務、参照入口は [`docs/DETAIL_INDEX.md`](../DETAIL_INDEX.md) 詳細仕様入口責務、fixture schema と記録先は [`docs/details/fixture.md`](fixture.md) fixture 証跡責務、実在 path は [`docs/DOCUMENT_INDEX.md`](../DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務を正本とする。

<a id="phase-17-responsibility-boundary"></a>
**0. 責務境界：**

本番検証詳細本文責務は Phase 17 の ConoHa VPS 試験本番運用と Phase 18 の VPS Pull Bootstrap に関する、運用対象、provider target、OS / plan 前提、minimum plan class、opaque environment identity、検証対象、検証 mode、credential / SSH 入力境界、Phase 17 試験本番運用窓、Phase 18 bootstrap session 停止条件、禁止操作、運用中バグ修正順序、bootstrap session、fixture 証跡への接続条件を所有する。

本番検証詳細本文責務は、owner component の通常処理、API route、SDK method、UI DOM、runner pipeline、state schema、release asset format、setup install 処理、VPS provider API 操作実装、GitHub 設定、credential 管理方式を本文として定義しない。これらは該当する owner component 別詳細本文、[`docs/SPEC.md`](../SPEC.md) 方針責務・ポリシー責務、または [`AGENTS.md`](../../AGENTS.md) 作業ルールを参照する。ただし、Phase 17 と Phase 18 の検証で secret を証跡へ混入させないための入力境界、許可入力経路、禁止入力経路、証跡化禁止事項は本番検証詳細本文責務が所有する。これは Adlaire CI の credential 管理方式ではなく、本番検証の boundary contract として扱う。

Phase 17 は ConoHa VPS 試験本番運用の Phase であり、正式本番運用、customer data を扱う実運用、ユーザー環境の破壊的変更、provider resource の作成・削除自動化を目的にしてはならない。

<a id="phase-17-purpose"></a>
**1. Phase 17 目的：**

Phase 17 は、ConoHa VPS 上に本番運用前提の試験本番環境を置くための初期本番検証基盤を完成させ、試験本番VPS 作成後に同じ証跡契約で本番検証、本番環境同等テスト、simulation、問題検出、仕様全般策定、バグ修正、再検証、安定化を反復できる状態にする Phase とする。

Phase 17 は以下を目的とする。

- ConoHa VPS を先行 target として、Ubuntu Server 24.04 LTS 64bit の最小採用 plan 上で、install、update、rollback、service lifecycle、runtime flow、state persistence、network、filesystem、recovery の本番運用前提を検証するための入力、証跡、counter、required check を固定する。
- destructive operation を伴う failure class は、本番 VPS へ直接適用せず、本番同等 simulation で検証する。
- エックスサーバ VPS は今後の開発状況で判断する将来対象として保持し、Phase 17 の完了必須条件に含めない。
- 試験本番運用中に検出したバグ、仕様不整合、環境不整合は、実装修正前に仕様全般策定へ戻し、責務正本を改訂してから修正、再配置、再検証する。
- 本番検証の成功、未実行、対象外、将来判断、運用中修正、再検証結果を fixture 証跡に接続し、口頭報告、手作業メモ、VPS 画面確認だけを完了根拠にしない。

<a id="phase-17-provider-targets"></a>

**2. 本番環境想定：**

Phase 17 の provider target は以下に固定する。

| provider target | 扱い | 完了条件への影響 |
|-----------------|------|------------------|
| `conoha-vps-primary` | ConoHa VPS 先行試験本番運用対象。 | Phase 17 完了に必須。 |
| `xserver-vps-future` | エックスサーバ VPS 将来判断対象。 | Phase 17 完了の阻害要因にしない。 |

`conoha-vps-primary` は、実 provider 上の試験本番運用対象である。Phase 17 実装時は、試験本番VPS を使用し、正式本番の customer data、実運用 secret、手動運用中の本番 state を検証対象にしてはならない。

`xserver-vps-future` は、[`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務の将来判断行へ接続する。`xserver-vps-future` を Phase 17 の未完了 item、必須 check、failure、または blocker として扱ってはならない。

Phase 17 の ConoHa VPS 試験本番運用環境は、Ubuntu Server 24.04 LTS 64bit を固定 OS とし、`minimum_plan_class=conoha-vps-1gb-memory-class` を固定最小 plan class とする。`minimum_plan_class` は、RAM が 1024 MB 以上で Ubuntu Server 24.04 LTS 64bit を選択できる ConoHa VPS の最小 memory class を表す repository 内固定 token であり、provider 画面や API の SKU 表示名を正本にしない。provider 側の plan 表示名、plan id、region、instance 名は `provider_plan_label` 等の opaque label としてだけ fixture 証跡に記録し、`minimum_plan_class` の値を変更してはならない。512MB plan、Ubuntu 26.04、Ubuntu 22.04、Debian、AlmaLinux、Rocky Linux、CentOS Stream、Oracle Linux、FreeBSD、および application template は Phase 17 の標準 OS / plan として扱わない。

`Build-Scripts-vps-2026` は、`conoha-vps-primary` の非秘密 name tag として使用する。この tag は試験本番運用・開発検証兼用 VPS の opaque environment identity であり、Phase 17 の試験本番運用、Phase 18 の VPS Pull Bootstrap、開発検証、運用中バグ修正後の再検証、update / rollback 確認に使用できる。`Build-Scripts-vps-2026` を正式本番 VPS、customer data 用 VPS、Phase 専用 VPS、または provider resource id として扱ってはならない。

ConoHa VPS 試験本番運用環境の provider plan、region、VPS instance、public endpoint は、fixture 証跡では opaque label としてだけ記録する。provider account id、provider resource id、グローバル IP address、hostname、FQDN、credential file path、secret 値、secret hash を fixture、expected、record、Pull Request body、log、stdout、stderr に保存してはならない。実 provider の識別が必要な場合は、operator が管理する repository 外の対応表で照合し、repository 内の証跡には `metadata_policy=opaque-non-secret-labels` を記録する。

<a id="phase-17-validation-modes"></a>

**3. 検証 mode：**

Phase 17 の検証 mode は以下に固定する。

| mode | 実行場所 | 用途 |
|------|----------|------|
| `trial-production-conoha-vps` | ConoHa VPS の試験本番運用環境 | 継続稼働、install、update、rollback、service lifecycle、API health、admin CLI health、runner dry-run、state directory、network timeout、reboot recovery、運用中バグ修正ループを検証する。 |
| `production-equivalent-simulation` | local または container による本番同等 simulation | [Phase 17 failure class 固定 set](#phase-17-failure-class-set) と destructive operation 境界を検証する。 |

`production-equivalent-simulation` は destructive operation と fault injection の検証に使用する。`production-equivalent-simulation` の成功だけで `trial-production-conoha-vps` の必須運用検証を代替してはならない。

`trial-production-conoha-vps` は provider 固有の lifecycle と継続稼働を確認するために使用するが、VPS deletion、disk rebuild、volume detach、firewall lockout、SSH lockout、秘密情報の出力、customer data の投入を実行してはならない。

<a id="phase-17-failure-class-set"></a>

Phase 17 failure class 固定 set は `disk full`、`permission denied`、`short write`、`fsync failure`、`rename failure`、`process kill`、`network refusal`、`DNS failure`、`SSH failure`、`log write failure`、`reboot recovery`、`interrupted recovery` の順序固定 array とする。Phase 17 の simulation、failure injection、required check、fixture、record、PR 証跡は、この固定 set 以外の failure class 名、欠落、追加、重複、順序変更、`reboot` / `中断復旧` などの別表記を使用してはならない。

<a id="phase-17-buildout-plan"></a>

**3a. 構築プラン：**

Phase 17 の構築は、以下の work unit を順序固定で実行する。

| work unit | 完了条件 |
|-----------|----------|
| `provider-prerequisite` | ConoHa VPS 上に Ubuntu Server 24.04 LTS 64bit、`minimum_plan_class=conoha-vps-1gb-memory-class` の試験本番VPS が存在し、512MB plan、別 OS、application template、正式本番 data が使われていないことを記録する。provider plan、region、VPS instance、public endpoint は opaque label だけで記録し、provider account id、resource id、IP address、hostname、FQDN を証跡化しない。 |
| `access-baseline` | SSH 到達性、known_hosts、host key、管理用 user、sudo capability、時刻同期、opaque host label、network、DNS、firewall の初期状態を記録する。hostname、IP address、FQDN は証跡へ保存しない。 |
| `system-baseline` | kernel、systemd、filesystem、available disk、required command、umask、locale、timezone、state directory parent の owner / mode を記録する。 |
| `release-acquisition` | GitHub Release asset、SHA256SUMS、version、binary 起動結果を取得し、checksum mismatch、取得失敗、version mismatch を失敗として固定する。 |
| `install-bootstrap` | state directory、config、credential 初期化、systemd unit、専用 user、最小権限、service start、health endpoint を構築し、stdout / stderr contract を記録する。 |
| `runtime-baseline` | API health、Admin CLI、SDK method、UI 疎通、runner dry-run、build dry-run、state write、audit log、secret mask を記録する。 |
| `operation-baseline` | 試験本番運用開始時点の service state、queue state、statefile digest、log digest、reachable endpoint、failure recovery baseline を記録する。 |
| `update-rollback-baseline` | update、rollback、service restart、state migration、rollback staging、rollback 後 version、failure cleanup を記録する。 |
| `simulation-baseline` | 本番 VPS 上で実行してはならない destructive / fault class を `production-equivalent-simulation` へ分離し、実 provider に適用していないことを記録する。 |
| `monitoring-baseline` | health、journal、audit、access、config log、state write、disk usage、credential leakage、known bug の監視項目を記録する。 |

ConoHa VPS resource の作成、削除、plan 変更、disk rebuild、volume 操作、firewall lockout は、Adlaire CI の Phase 17 実装機能として自動化しない。Phase 17 は、試験本番VPS 作成後に preflight で検証する証跡契約を固定し、その後の install、service、runtime、update、rollback、運用、simulation の結果を証跡化できる状態を完了条件とする。

禁止 provider operation token は `vps-create`、`vps-delete`、`plan-change`、`disk-rebuild`、`volume-create`、`volume-delete`、`volume-attach`、`volume-detach`、`firewall-lockout`、`ssh-lockout` に固定する。Phase 17 の fixture、negative boundary、required check、PR 証跡は、この token set 以外の名称で禁止 provider operation を表現してはならない。

<a id="phase-17-secret-input-boundary"></a>

**3b. credential / SSH 入力境界：**

Phase 17 の試験本番運用検証で使用する provider API token、admin token、SSH private key、IP 固有 secret、host 固有 secret は、repository、fixture、expected、record、workflow、log、stdout、stderr、Pull Request body、issue comment、chat message、screen shot に保存してはならない。

Phase 17 の許可 secret 入力経路は以下に固定する。

| input channel | 許可条件 |
|---------------|----------|
| `operator-local-file` | repository root 外の operator local file だけを使用し、証跡には file path を保存せず `secret_class`、`input_channel`、`redacted_reference` だけを保存する。 |
| `operator-stdin` | secret 値を record、stdout、stderr、shell history、log に残さない。 |
| `ssh-agent` | secret 値を取り出して記録しない。 |
| `systemd-credential` | secret 値を取り出して記録しない。 |
| `not-required` | 対象 check が secret を必要としない場合だけ使用する。 |

Phase 17 の禁止 secret 入力経路は以下に固定する。

| input channel | 禁止条件 |
|---------------|----------|
| `repository-file` | repository 内 file へ secret を置くことを禁止する。 |
| `fixture-file` | fixture、expected、record へ secret を置くことを禁止する。 |
| `command-argument` | shell command argument として secret を渡すことを禁止する。 |
| `environment-variable` | 環境変数として secret を渡すことを禁止する。 |
| `pr-body` | Pull Request body へ secret を保存することを禁止する。 |
| `issue-comment` | issue comment へ secret を保存することを禁止する。 |
| `chat-message` | chat message へ secret を保存することを禁止する。 |
| `stdout` | stdout へ secret を出すことを禁止する。 |
| `stderr` | stderr へ secret を出すことを禁止する。 |
| `log` | log へ secret を保存することを禁止する。 |
| `screen-shot` | screen shot へ secret を写すことを禁止する。 |

禁止入力経路を secret 入力として成功扱いしてはならない。禁止入力経路の検査 record は、禁止経路を遮断した境界証跡として [`docs/details/fixture.md` fixture 証跡責務 Phase 17 ConoHa VPS 試験本番運用証跡](fixture.md#phase-17-production-validation-evidence) の `records/security.jsonl` へ記録し、`boundary_result=blocked` と secret boundary counter へ接続する。

Phase 17 の証跡に保存できる secret 関連情報は、`secret_class`、`input_channel`、`redacted_reference`、`secret_reference_policy=metadata-only`、`secret_value_present=false`、`scan_target`、`secret_scan_result`、`mask_result`、`credential_storage_result`、`destructive_operation_result`、`boundary_result` だけとする。`scan_target` は検査対象分類だけを表し、host、IP、path、secret 値、credential file path を含めてはならない。secret 値、secret 値の hash、private key fingerprint、provider account id、credential file path、host 固有 secret、IP 固有 secret、実 token の prefix / suffix を保存してはならない。

<a id="phase-17-trial-operation-window"></a>

**3c. 試験本番運用窓・停止条件：**

Phase 17 の `trial-production-conoha-vps` の試験本番運用窓は `initial-production-validation` 固定とする。運用窓は `operation-baseline` の成功後に開始し、service start、API health、Admin CLI health、runner dry-run、state write、audit log、secret mask の初回成功と、停止条件が発火していないことを 1 件の baseline monitor sample として記録する。

Phase 17 の初期本番検証監視間隔は `0` 秒固定とし、必須監視 sample 数は `1` 件固定とする。監視 sample は health、systemd state、state write、audit / access / config log write、disk usage、network reachability、secret leakage boundary、known bug status を同一 sample として記録する。監視 sample の時刻、index、結果、停止条件判定、再検証参照は [`docs/details/fixture.md` fixture 証跡責務 Phase 17 ConoHa VPS 試験本番運用証跡](fixture.md#phase-17-production-validation-evidence) の `records/operation.jsonl` へ接続する。長時間継続運用監視は Phase 17 の完了 blocker にせず、試験本番VPS 作成後の後続運用記録として扱う。

Phase 17 の初期本番検証 failure 許容数は `0` 件とする。`failed`、`unknown`、記録欠落、または baseline monitor sample の必須観測値欠落は、試験本番運用の未完了として fixture 証跡へ接続する。failure を検出した場合は、仕様全般策定先行、修正、再配置、再検証、再監視 sample の証跡が揃うまで Phase 17 を完了扱いにしてはならない。

Phase 17 の停止条件は `secret-boundary-failure`、`destructive-operation-boundary-failure`、`state-corruption`、`service-unrecoverable`、`known-critical-bug` に固定する。停止条件が発火した場合は、試験本番運用を成功扱いせず、停止理由、直前 sample、仕様全般策定先、修正対象、再検証条件、復旧可否を証跡化する。停止条件発火後の修正は、[`docs/SPEC.md`](../SPEC.md) 方針責務・ポリシー責務と本ファイルの試験本番運用ループに従い、仕様全般策定を先行しなければならない。

<a id="phase-17-validation-targets"></a>

**4. 検証対象：**

Phase 17 は以下を検証対象に含める。

| 対象 | 必須確認 |
|------|----------|
| setup install | release asset 取得、checksum、install path、state directory 作成、systemd unit 配置、stdout / stderr contract。 |
| update / rollback | version 遷移、rollback staging、rollback 後 service 状態、state file 復旧、失敗時 cleanup。 |
| release asset | binary 起動、version 出力、SHA256SUMS 照合、署名または署名対象外理由、SBOM または対象外理由。 |
| systemd lifecycle | enable、start、restart、stop、status、journal、専用 user、最小権限、書込み先制限。 |
| API / Admin / SDK / UI | health、credential 初期化、admin CLI 接続、SDK method、UI から API への疎通、timeout、body limit。 |
| runner / build / deploy | dry-run、queue 生成、build 実行、deploy simulation、commit status 対象外または成功証跡。 |
| statefile | state directory `0700`、owner / mode 検証、lock、atomic write、JSON Lines、migration、復旧。 |
| SSH / remote path | known_hosts、host key、timeout、remote path validation、private IP / SSRF 境界。 |
| network | DNS、IPv4 / IPv6 境界、redirect 禁止、timeout、firewall 影響、provider outbound / inbound 境界。 |
| failure injection | [Phase 17 failure class 固定 set](#phase-17-failure-class-set) 全件。 |
| log / audit / secret | audit log、access log、config log の必須書込み失敗、secret mask、token / key 非出力。 |
| trial production operation | `initial-production-validation`、1 baseline sample、停止条件、継続稼働前提、問題検出、仕様全般策定、バグ修正、再配置、再検証、証跡記録の反復。 |
| document drift | Phase 17 状態、path、anchor、fixture root、workflow、required check、future target の drift。 |

各対象の詳細な owner 実装契約は、対象 owner component 別詳細本文を参照する。Phase 17 は、これらの実装契約を本番環境同等条件で検証する入口であり、owner 詳細本文の処理仕様を再定義しない。

<a id="phase-17-forbidden-items"></a>
**5. 禁止事項：**

Phase 17 では以下を禁止する。

- customer data、実運用 credential、実運用 secret、private key、個人情報を repository、fixture、log、PR body、artifact に保存する。
- VPS deletion、disk rebuild、volume detach、firewall lockout、SSH lockout、account lockout を自動検証または fixture 実行で行う。
- provider API token、admin token、SSH private key、IP 固有 secret、hostname 固有 secret を repository に保存する。
- `xserver-vps-future` を Phase 17 の必須検証未完了 item として扱う。
- `production-equivalent-simulation` のみを根拠に `trial-production-conoha-vps` 検証を完了扱いにする。
- 手作業画面確認、スクリーンショット、口頭報告だけを完了証跡にする。
- 本番 VPS 上で destructive failure injection を直接実行する。
- 試験本番運用中のバグ、仕様不整合、環境不整合を、仕様全般策定なしに hotfix、直接修正、口頭判断、または実装者判断だけで修正する。

<a id="phase-17-conoha-contract"></a>

**6. ConoHa VPS 先行試験本番運用契約：**

`trial-production-conoha-vps` は `conoha-vps-primary` に対して実行する。

ConoHa VPS 試験本番運用検証は以下の順序で実行する。

1. `preflight` で OS、kernel、systemd、filesystem、available disk、network、DNS、user、sudo capability、state directory parent、time sync、required command availability を記録する。
2. `install` で release asset、checksum、binary 起動、state directory、systemd unit、service start、health endpoint を検証する。
3. `runtime-flow` で Admin CLI、API health、credential 初期化、runner dry-run、build dry-run、state write、audit log、secret mask を検証する。
4. `update` で version 変更、service restart、state migration、rollback staging を検証する。
5. `rollback` で rollback commit、service 復旧、state 整合、旧 version 出力、failure cleanup を検証する。
6. `reboot-recovery` で VPS reboot 後の service 自動復旧、lock 復旧、active queue 復旧、log continuation を検証する。
7. `cleanup` で検証用 user、temporary file、temporary state、temporary service override を削除または隔離し、cleanup 結果を記録する。

ConoHa VPS 試験本番運用検証では、削除系 provider operation を実行せず、VPS 外部からの操作は read-only または service lifecycle に限定する。provider resource の作成や削除を行う場合は、Phase 17 実装 PR の範囲外とし、別途仕様策定と承認を必要とする。

Phase 17 の試験本番運用ループは `operation-start`、`monitor`、`issue-detect`、`spec-general-update`、`implementation-fix`、`redeploy`、`revalidate`、`evidence-record`、`known-bug-zero-check` の順序に固定する。`spec-general-update` を通過しない `implementation-fix` を禁止し、仕様全般策定の責務正本、変更範囲、再検証条件、証跡接続が未確定のまま修正してはならない。

<a id="phase-17-xserver-future-contract"></a>

**7. エックスサーバ VPS 将来判断契約：**

エックスサーバ VPS は Phase 17 時点では `xserver-vps-future` として扱う。

Phase 17 の fixture 証跡へ接続する future target 分類は、`xserver-vps-future` について以下を記録する。

| key | 固定値 |
|-----|--------|
| `provider_target` | `xserver-vps-future` |
| `classification` | `future_plan` |
| `future_ref` | [`docs/ROADMAP.md` 状態・計画責務 統合機能インベントリ](../ROADMAP.md#522-統合ロードマップ表) |
| `completion_blocker` | `false` |

`xserver-vps-future` に対して required check、fixture root、provider credential、VPS execution record、simulation 専用分岐を作成する場合は、Phase 17 の完了条件ではなく、将来 Phase または改訂予定 item として扱う。

<a id="phase-17-production-validation-evidence-connection"></a>
<a id="phase-17-production-validation-completion"></a>

**8. 証跡接続条件：**

Phase 17 の fixture 証跡、closure counter、required check、required check workflow、checker 実行入口、negative boundary、document drift record は [`docs/details/fixture.md` fixture 証跡責務 Phase 17 ConoHa VPS 試験本番運用証跡](fixture.md#phase-17-production-validation-evidence) を正本とする。

Phase 17 は、[`docs/details/fixture.md` fixture 証跡責務 Phase 17 ConoHa VPS 試験本番運用証跡](fixture.md#phase-17-production-validation-evidence) の closure counter、required check、required check workflow、checker 実行入口、negative boundary、document drift record がすべて完了条件を満たすまで `実装済み` に遷移してはならない。本番検証詳細本文責務では、closure counter の key、完了値、未完了条件、required check 名、record schema を再掲しない。

Phase 17 の実装 PR は、[`docs/details/fixture.md` fixture 証跡責務 Phase 17 ConoHa VPS 試験本番運用証跡](fixture.md#phase-17-production-validation-evidence) が要求する全 record を同一証跡 package に含める。

<a id="phase-18-vps-pull-bootstrap-contract"></a>

**9. Phase 18 VPS Pull Bootstrap 契約：**

Phase 18 は、作成済み VPS に `VPS Pull Bootstrap` を導入し、Adlaire CI の CI/CD 実行基盤を VPS 側へ入れる Phase とする。`VPS Pull Bootstrap` は shell 製 bootstrap script による初回導入方式であり、GitHub Actions、GitHub Secrets、GitHub self-hosted runner、provider 固有の startup script 機能、または provider API を Phase 18 の必須実行基盤として扱ってはならない。

Phase 18 の provider target は `conoha-vps-primary` を先行対象、`xserver-vps-future` を将来判断対象に固定する。`xserver-vps-future` は Phase 18 の未完了 item、required check、failure、または blocker として扱わない。Phase 18 は provider 非依存の bootstrap 契約を固定する Phase であり、ConoHa 固有機能を正本条件として埋め込んではならない。

Phase 18 は VPS 作成、VPS 削除、plan 変更、disk rebuild、volume 操作、firewall lockout、SSH lockout、provider account 設定変更、GitHub 設定変更を実装対象にしない。VPS の作成と bootstrap script の初回実行は operator が実施する前提であり、Phase 18 の仕様は、その後 VPS 内で Adlaire CI 実行基盤が自律的に pull / update / rollback / evidence 記録できる状態を固定する。

Phase 18 の実行環境前提は Phase 17 と同じく、Ubuntu Server 24.04 LTS 64bit、`minimum_plan_class=conoha-vps-1gb-memory-class`、`minimum_ram_mb=1024`、`metadata_policy=opaque-non-secret-labels` に固定する。VPS の用途 label は `試験本番VPS` とする。ただし Phase 18 自体を試験本番運用 Phase と呼んではならず、Phase 18 は `VPS Pull Bootstrap` による実行基盤導入 Phase として扱う。

Phase 18 の validation mode は `vps-pull-bootstrap` 固定とする。`production-equivalent-simulation` は Phase 17 の destructive / fault class 検証 mode であり、Phase 18 の bootstrap 契約を代替してはならない。

Phase 18 の bootstrap session は `operator-approved-bootstrap-session` 固定とする。Phase 18 実装完了には 1 件以上の bootstrap session を必須とし、24 時間継続監視は完了 blocker にしない。bootstrap session 内では、開始時、bootstrap runtime flow 完了後、終了時の 3 sample を最小証跡として記録する。各 sample は service state、API health、Admin CLI health、runner state、statefile digest、audit / access / config log write、secret leakage boundary、known bug status、document drift status を含む。

Phase 18 の topology role は `ci-cd`、`site`、`single-node` に固定する。標準運用は Linux CI/CD server と静的サイト配信 server の 2 台構成である。最小運用は 1 台 VPS の `single-node` として、CI/CD 実行基盤と静的サイト配信を同一 VPS 内に同居できる。`single-node` は最小構成であり、2 台構成の置換ではなく bootstrap 契約上の許可構成として扱う。

Phase 18 の source channel は `integration-head` と `stable-release` を区別する。Phase 18 では `integration-head` を bootstrap 検証用の取得経路、`stable-release` を将来の本番運用用取得経路として定義する。`stable-release` の正式 CD、配信 channel、release selection、rollback policy は Phase 19 以降の対象であり、Phase 18 の完了条件へ混入してはならない。

Phase 18 の bootstrap script artifact は配布物名 `adlaire-ci-vps-pull-bootstrap.sh`、runtime `/bin/sh`、呼出し境界 `operator-runs-sh-script-on-vps` に固定する。bootstrap script は provider 固有の startup script 機能ではなく、operator が作成済み VPS 内で初回実行する shell artifact である。bootstrap script は Go 標準実装で生成または配布される実行基盤を VPS 内へ導入する入口に限定し、provider API、GitHub Actions、GitHub Secrets、GitHub self-hosted runner、SSH key 管理、VPS 作成、VPS 削除、plan 変更を内包してはならない。

Phase 18 の bootstrap script は `/bin/sh` で構文解釈できる範囲に固定し、Bash 固有構文、外部 shell framework、外部 package manager helper、外部 library、Git clone 必須化、repository hosting provider 固有 API 必須化を禁止する。HTTP 取得が必要な場合は、VPS に既に存在する OS 標準 command の可用性を bootstrap artifact health で記録し、取得 command が存在しない場合は bootstrap session を成功扱いせず `bootstrap-script-http-client-missing` として issue triage へ接続する。取得 command を暗黙に install して成功扱いにしてはならない。

Phase 18 の bootstrap script 入力は、`integration-head` または `stable-release` の source channel、`ci-cd` / `site` / `single-node` の topology role、取得元識別子、取得対象 digest、install directory、bin directory、state directory、service user に限定する。取得元識別子形式は `absolute-https-url-or-absolute-file-path`、digest algorithm は `sha256`、install directory は `/opt/adlaire-builder`、bin directory は `/usr/local/bin`、state directory は `/opt/adlaire-builder`、service user は `root` 固定とする。入力値は command argument として secret を渡してはならず、secret を必要とする取得経路は Phase 18 の完了条件に含めない。source channel、topology role、bootstrap script artifact、source identifier 形式、digest algorithm、install / bin / state directory、service user、secret 入力禁止、state model は [`docs/details/fixture.md` fixture 証跡責務 Phase 18 VPS Pull Bootstrap](fixture.md#phase-18-vps-pull-bootstrap-evidence) の `input/operation_scope.json` と `records/health.jsonl` に接続しなければならない。

Phase 18 の bootstrap script state model は `staging-verify-commit-rollback` 固定とする。bootstrap script は staging directory へ取得し、digest 検証、実行権限検証、version 検証、systemd unit 検証、state directory 検証、health 検証を完了してから commit し、検証前の失敗では既存 binary、既存 service、既存 state、既存 credential を変更してはならない。commit 後の health 失敗では rollback 証跡を作成し、復旧可否、直前 version、失敗段階、再実行条件を bootstrap update / rollback drill と issue triage へ接続する。

Phase 18 の work unit は以下の順序固定とする。

| work unit | 完了条件 |
|-----------|----------|
| `phase17-handover` | Phase 17 の fixture 証跡、required check、状態、VPS 前提、未残条件が Phase 18 の開始条件として到達可能である。 |
| `bootstrap-session-open` | `operator-approved-bootstrap-session` を開始し、operator 承認、session id、開始時 sample、secret 境界、禁止 provider operation 境界を記録する。 |
| `bootstrap-artifact-health` | bootstrap script、取得 channel、binary version、systemd service、health endpoint、state directory、log write が成功し、失敗または未実施を open item へ接続する。 |
| `bootstrap-runtime-flow` | API、Admin CLI、SDK、UI の疎通、認証、timeout、body limit、secret mask、状態再取得が成功する。 |
| `pull-runner-foundation` | runner queue、pull source、build、状態更新、audit log、commit status 対象外または成功、deploy simulation または承認済み deploy target が証跡化される。 |
| `bootstrap-update-rollback-drill` | update、rollback、service restart、state migration、rollback 後 version、failure cleanup が成功する。 |
| `bootstrap-issue-triage` | bootstrap 中の issue、仕様不整合、環境不整合、known bug を検出し、分類、責務正本、修正要否、再検証条件を記録する。 |
| `spec-first-fix` | 修正が必要な issue は、実装修正前に仕様全般策定、責務正本改訂、検証条件固定を完了する。 |
| `rebootstrap-and-revalidate` | 修正後の再配置、service 復旧、対象 required check 再実行、再検証結果、known bug 0 判定を記録する。 |
| `bootstrap-session-close` | 終了時 sample、open item 0、secret 漏えい 0、destructive operation 0、document drift 0、cleanup 結果を記録して session を閉じる。 |

Phase 18 の bootstrap session は、`bootstrap-session-open` で開始 sample を記録し、`bootstrap-runtime-flow`、`pull-runner-foundation`、`bootstrap-update-rollback-drill` の完了後に runtime sample を記録し、`rebootstrap-and-revalidate` と cleanup 判定後に終了 sample を記録する。開始 sample を runtime 実行後に後付けで作成すること、runtime sample を update / rollback 前の状態で代替すること、終了 sample を open item の有無を確認せず作成することを禁止する。

Phase 18 の `bootstrap-issue-triage` は、bootstrap 中の issue を検出した場合だけでなく、issue が検出されなかった場合も `no_issue_detected` として記録する。`spec-first-fix` は `fix_required`、`environment_only`、`not_required` のいずれかを明示し、`fix_required` の場合は仕様全般策定、責務正本改訂、実装修正、再検証を同一 bootstrap session の証跡へ接続する。`environment_only` または `not_required` の場合も、修正不要理由、責務正本、再検証条件を記録し、作業省略を口頭判断または暗黙の成功として扱ってはならない。

Phase 18 の `rebootstrap-and-revalidate` は、実装修正があった場合は再 bootstrap 後の service 復旧と対象再検証を記録し、実装修正がない場合は no-op 再検証として、bootstrap runtime flow、停止条件、known bug 0、document drift 0 の再判定を記録する。実装修正がないことを理由に `rebootstrap-and-revalidate` の記録を省略してはならない。

Phase 18 の bootstrap closure package は、Phase 17 引継ぎ、開始 sample、runtime sample、終了 sample、bootstrap artifact health、API / Admin / SDK / UI runtime、pull runner foundation、bootstrap update / rollback drill、issue triage、仕様先行修正判定、再検証、secret 境界、destructive operation 境界、document drift 境界、cleanup 判定を同一 session id に接続する。同一 session id で接続できない証跡、別 session の結果を寄せ集めた証跡、session id を持たない証跡、または fixture counter だけの完了報告を Phase 18 の完了根拠にしてはならない。

Phase 18 の bootstrap 完了証跡は、VPS 内の実測結果から作成する。[`docs/details/fixture.md` fixture 証跡責務 Phase 18 VPS Pull Bootstrap](fixture.md#phase-18-vps-pull-bootstrap-evidence) の `testdata/phase18/vps-pull-bootstrap/records/*.jsonl` が `evidence_origin=checker-acceptance-fixture` と `execution_environment=repository-fixture` を持つ場合、その record は checker acceptance fixture としてだけ扱い、Phase 18 の bootstrap 完了証跡として扱ってはならない。

Phase 18 の bootstrap 中バグ修正 loop は `issue-detect`、`spec-general-update`、`implementation-fix`、`redeploy`、`revalidate`、`evidence-record`、`known-bug-zero-check` の順序に固定する。`spec-general-update` を通過しない `implementation-fix` を禁止する。

Phase 18 の停止条件は `secret-boundary-failure`、`destructive-operation-boundary-failure`、`state-corruption`、`service-unrecoverable`、`known-critical-bug`、`document-drift-blocker` に固定する。停止条件が発火した場合は、bootstrap session を成功扱いせず、停止理由、直前 sample、責務正本、修正条件、再検証条件、復旧可否を証跡化する。

Phase 18 では以下を禁止する。

- customer data、実運用 credential、実運用 secret、private key、個人情報を repository、fixture、log、PR body、artifact に保存する。
- provider account id、provider resource id、IP address、hostname、FQDN、credential file path、secret 値、secret hash、token prefix / suffix を証跡化する。
- VPS 作成、VPS 削除、plan 変更、disk rebuild、volume 操作、firewall lockout、SSH lockout を Phase 18 の成功操作として扱う。
- GitHub Actions、GitHub Secrets、GitHub self-hosted runner、provider startup script 機能、または provider API を Phase 18 の必須実行面として扱う。
- `xserver-vps-future` を Phase 18 の必須検証未完了 item として扱う。
- 24 時間監視がないことだけを理由に Phase 18 を未完了扱いにする。
- スクリーンショット、口頭報告、手作業メモ、provider 画面確認だけを完了証跡にする。
- 仕様全般策定なしの hotfix、実装者判断だけの修正、または fixture counter だけの完了報告を成功扱いにする。

Phase 18 の fixture 証跡、closure counter、required check、required check workflow、checker 実行入口、negative boundary、document drift record は [`docs/details/fixture.md` fixture 証跡責務 Phase 18 VPS Pull Bootstrap](fixture.md#phase-18-vps-pull-bootstrap-evidence) を正本とする。本番検証詳細本文責務では、Phase 18 closure counter の key、完了値、record schema、required check 名を再掲しない。
