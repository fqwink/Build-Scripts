# Adlaire CI — 試験本番運用詳細仕様

[`docs/details/production-validation.md`](production-validation.md) は、本番検証詳細本文責務として、ConoHa VPS 試験本番運用、本番運用前提の本番検証、本番環境同等テスト、VPS simulation、provider target、destructive operation 禁止境界、運用中バグ修正順序、完了証跡だけを定義する。状態、Phase、実装可否は [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務、方針とポリシーは [`docs/SPEC.md`](../SPEC.md) 方針責務・ポリシー責務、参照入口は [`docs/DETAIL_INDEX.md`](../DETAIL_INDEX.md) 詳細仕様入口責務、fixture schema と記録先は [`docs/details/fixture.md`](fixture.md) fixture 証跡責務、実在 path は [`docs/DOCUMENT_INDEX.md`](../DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務を正本とする。

## 0. 責務境界

本番検証詳細本文責務は Phase 17 の ConoHa VPS 試験本番運用に関する、運用対象、provider target、OS / plan 前提、検証対象、検証 mode、禁止操作、運用中バグ修正順序、fixture 証跡、完了条件を所有する。

本番検証詳細本文責務は、owner component の通常処理、API route、SDK method、UI DOM、runner pipeline、state schema、release asset format、setup install 処理、VPS provider API 操作実装、GitHub 設定、credential 管理方式を本文として定義しない。これらは該当する owner component 別詳細本文、[`docs/SPEC.md`](../SPEC.md) 方針責務・ポリシー責務、または [`AGENTS.md`](../../AGENTS.md) 作業ルールを参照する。

Phase 17 は ConoHa VPS 試験本番運用の Phase であり、正式本番運用、customer data を扱う実運用、ユーザー環境の破壊的変更、provider resource の作成・削除自動化を目的にしてはならない。

## 1. Phase 17 目的

Phase 17 は、ConoHa VPS 上に本番運用前提の試験本番環境を置き、実際に稼働させながら、本番検証、本番環境同等テスト、simulation、問題検出、仕様全般策定、バグ修正、再検証、安定化を反復する試験本番運用 Phase とする。

Phase 17 は以下を目的とする。

- ConoHa VPS を先行 target として、Ubuntu Server 24.04 LTS 64bit の最小採用 plan 上で、install、update、rollback、service lifecycle、runtime flow、state persistence、network、filesystem、recovery の本番運用前提を検証する。
- destructive operation を伴う failure class は、本番 VPS へ直接適用せず、本番同等 simulation で検証する。
- エックスサーバ VPS は今後の開発状況で判断する将来対象として保持し、Phase 17 の完了必須条件に含めない。
- 試験本番運用中に検出したバグ、仕様不整合、環境不整合は、実装修正前に仕様全般策定へ戻し、責務正本を改訂してから修正、再配置、再検証する。
- 本番検証の成功、未実行、対象外、将来判断、運用中修正、再検証結果を fixture 証跡に接続し、口頭報告、手作業メモ、VPS 画面確認だけを完了根拠にしない。

<a id="phase-17-provider-targets"></a>

## 2. 本番環境想定

Phase 17 の provider target は以下に固定する。

| provider target | 扱い | 完了条件への影響 |
|-----------------|------|------------------|
| `conoha-vps-primary` | ConoHa VPS 先行試験本番運用対象。 | Phase 17 完了に必須。 |
| `xserver-vps-future` | エックスサーバ VPS 将来判断対象。 | Phase 17 完了の阻害要因にしない。 |

`conoha-vps-primary` は、実 provider 上の試験本番運用対象である。Phase 17 実装時は、試験本番運用専用 VPS だけを使用し、正式本番の customer data、実運用 secret、手動運用中の本番 state を検証対象にしてはならない。

`xserver-vps-future` は、[`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務の将来判断行へ接続する。`xserver-vps-future` を Phase 17 の未完了 item、必須 check、failure、または blocker として扱ってはならない。

Phase 17 の ConoHa VPS 試験本番運用環境は、Ubuntu Server 24.04 LTS 64bit を固定 OS とし、Ubuntu 24.04 の要件を満たす ConoHa VPS 最小 plan を採用する。ConoHa VPS 上の最小採用 plan は RAM 1GB 以上とし、512MB plan は Phase 17 の対象外とする。Ubuntu 26.04、Ubuntu 22.04、Debian、AlmaLinux、Rocky Linux、CentOS Stream、Oracle Linux、FreeBSD、および application template は Phase 17 の標準 OS として扱わない。

<a id="phase-17-validation-modes"></a>

## 3. 検証 mode

Phase 17 の検証 mode は以下に固定する。

| mode | 実行場所 | 用途 |
|------|----------|------|
| `trial-production-conoha-vps` | ConoHa VPS の試験本番運用環境 | 継続稼働、install、update、rollback、service lifecycle、API health、admin CLI health、runner dry-run、state directory、network timeout、reboot recovery、運用中バグ修正ループを検証する。 |
| `production-equivalent-simulation` | local または container による本番同等 simulation | disk full、permission denied、short write、fsync failure、rename failure、process kill、network refusal、DNS failure、SSH failure、log write failure、destructive operation 境界を検証する。 |

`production-equivalent-simulation` は destructive operation と fault injection の検証に使用する。`production-equivalent-simulation` の成功だけで `trial-production-conoha-vps` の必須運用検証を代替してはならない。

`trial-production-conoha-vps` は provider 固有の lifecycle と継続稼働を確認するために使用するが、VPS deletion、disk rebuild、volume detach、firewall lockout、SSH lockout、秘密情報の出力、customer data の投入を実行してはならない。

<a id="phase-17-buildout-plan"></a>

## 3a. 構築プラン

Phase 17 の構築は、以下の work unit を順序固定で実行する。

| work unit | 完了条件 |
|-----------|----------|
| `provider-prerequisite` | ConoHa VPS 上に Ubuntu Server 24.04 LTS 64bit、RAM 1GB 以上の試験本番運用専用 VPS が存在し、512MB plan、別 OS、application template、正式本番 data が使われていないことを記録する。 |
| `access-baseline` | SSH 到達性、known_hosts、host key、管理用 user、sudo capability、時刻同期、hostname、network、DNS、firewall の初期状態を記録する。 |
| `system-baseline` | kernel、systemd、filesystem、available disk、required command、umask、locale、timezone、state directory parent の owner / mode を記録する。 |
| `release-acquisition` | GitHub Release asset、SHA256SUMS、version、binary 起動結果を取得し、checksum mismatch、取得失敗、version mismatch を失敗として固定する。 |
| `install-bootstrap` | state directory、config、credential 初期化、systemd unit、専用 user、最小権限、service start、health endpoint を構築し、stdout / stderr contract を記録する。 |
| `runtime-baseline` | API health、Admin CLI、SDK method、UI 疎通、runner dry-run、build dry-run、state write、audit log、secret mask を記録する。 |
| `operation-baseline` | 試験本番運用開始時点の service state、queue state、statefile digest、log digest、reachable endpoint、failure recovery baseline を記録する。 |
| `update-rollback-baseline` | update、rollback、service restart、state migration、rollback staging、rollback 後 version、failure cleanup を記録する。 |
| `simulation-baseline` | 本番 VPS 上で実行してはならない destructive / fault class を `production-equivalent-simulation` へ分離し、実 provider に適用していないことを記録する。 |
| `monitoring-baseline` | health、journal、audit、access、config log、state write、disk usage、credential leakage、known bug の監視項目を記録する。 |

ConoHa VPS resource の作成、削除、plan 変更、disk rebuild、volume 操作、firewall lockout は、Adlaire CI の Phase 17 実装機能として自動化しない。Phase 17 は、試験本番運用専用 VPS が存在することを preflight で検証し、その後の install、service、runtime、update、rollback、運用、simulation の結果を証跡化する。

禁止 provider operation token は `vps-create`、`vps-delete`、`plan-change`、`disk-rebuild`、`volume-create`、`volume-delete`、`volume-attach`、`volume-detach`、`firewall-lockout`、`ssh-lockout` に固定する。Phase 17 の fixture、negative boundary、required check、PR 証跡は、この token set 以外の名称で禁止 provider operation を表現してはならない。

## 4. 検証対象

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
| failure injection | disk full、permission denied、short write、fsync failure、rename failure、process kill、reboot、中断復旧。 |
| log / audit / secret | audit log、access log、config log の必須書込み失敗、secret mask、token / key 非出力。 |
| trial production operation | 継続稼働、問題検出、仕様全般策定、バグ修正、再配置、再検証、証跡記録の反復。 |
| document drift | Phase 17 状態、path、anchor、fixture root、workflow、required check、future target の drift。 |

各対象の詳細な owner 実装契約は、対象 owner component 別詳細本文を参照する。Phase 17 は、これらの実装契約を本番環境同等条件で検証する入口であり、owner 詳細本文の処理仕様を再定義しない。

## 5. 禁止事項

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

## 6. ConoHa VPS 先行試験本番運用契約

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

## 7. エックスサーバ VPS 将来判断契約

エックスサーバ VPS は Phase 17 時点では `xserver-vps-future` として扱う。

Phase 17 の証跡は、`xserver-vps-future` について以下を記録する。

| key | 固定値 |
|-----|--------|
| `provider_target` | `xserver-vps-future` |
| `classification` | `future_plan` |
| `future_ref` | [`docs/ROADMAP.md` 状態・計画責務 統合機能インベントリ](../ROADMAP.md#522-統合ロードマップ表) |
| `phase17_completion_blocker` | `false` |

`xserver-vps-future` に対して required check、fixture root、provider credential、real VPS execution、simulation 専用分岐を作成する場合は、Phase 17 の完了条件ではなく、将来 Phase または改訂予定 item として扱う。

<a id="phase-17-production-validation-completion"></a>

## 8. 証跡・完了条件

Phase 17 の fixture 証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 17 ConoHa VPS 試験本番運用証跡](fixture.md#phase-17-production-validation-evidence) を正本とする。

Phase 17 は、以下の closure counter がすべて `0` になるまで `実装済み` に遷移してはならない。

| counter | 完了値 | 未完了条件 |
|---------|--------|------------|
| `phase17_conoha_trial_operation_open_count` | `0` | ConoHa VPS 試験本番運用の必須 record が不足している。 |
| `phase17_buildout_open_count` | `0` | Phase 17 構築プランの work unit に未実行、順序不一致、証跡不足、または禁止 provider operation がある。 |
| `phase17_production_simulation_open_count` | `0` | 本番同等 simulation の必須 record が不足している。 |
| `phase17_install_update_rollback_open_count` | `0` | install、update、rollback のいずれかに未検証または未証跡がある。 |
| `phase17_systemd_lifecycle_open_count` | `0` | systemd lifecycle、専用 user、最小権限、書込み先制限に未完了がある。 |
| `phase17_runtime_flow_open_count` | `0` | API、Admin、SDK、UI、runner、build、deploy dry-run の runtime flow に未完了がある。 |
| `phase17_failure_injection_open_count` | `0` | failure injection class に未実行、silent success、または復旧未検証がある。 |
| `phase17_secret_leak_open_count` | `0` | secret、token、private key、IP 固有 secret が証跡へ出力される。 |
| `phase17_destructive_operation_open_count` | `0` | destructive operation が禁止境界を越えている、または検出証跡が不足している。 |
| `phase17_xserver_future_misclassified_count` | `0` | `xserver-vps-future` が必須検証、未完了 item、または blocker として扱われている。 |
| `phase17_document_drift_open_count` | `0` | Phase 17 の状態、path、anchor、fixture root、workflow、required check、future target に drift がある。 |
| `phase17_trial_operation_open_count` | `0` | 試験本番運用ループの継続稼働、問題検出、再配置、再検証、証跡記録に未完了がある。 |
| `phase17_bugfix_spec_gap_count` | `0` | 運用中に検出したバグまたは不整合に対し、仕様全般策定を先行していない修正がある。 |
| `phase17_known_bug_open_count` | `0` | Phase 17 完了時点で既知重大バグまたは未修正バグが残っている。 |
| `final_open_item_count` | `0` | 上記 counter または closure record に残件がある。 |

Phase 17 の実装 PR は、ConoHa VPS 試験本番運用 record、構築 record、本番同等 simulation record、failure injection record、secret leak 検査 record、destructive operation 境界 record、運用中バグ修正 record、document drift record を同一証跡 package に含める。
