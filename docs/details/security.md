# Adlaire CI — Security 詳細仕様

<a id="phase-19-security-contract"></a>

**Phase 19 identity・認可・secret契約：**

ownerは `security`、実装主体は [components/security/](../../components/security/)。適用profileは [`docs/SPEC.md` ポリシー責務 §0i](../SPEC.md#policy-inhouse-cicd)、schemaは [`docs/details/statefile.md` 詳細本文責務 schema](statefile.md#phase-19-schema-registry)。legacy login/sessionを黙って複数userへ拡張しない。

**認証：** cicd loginはuser ID/password/nullable totpを受け、既存credentialのKDF・定時間比較・TOTP・lockout規則へ接続する。id=adminだけcredential=nullを許し既存.admin_credentialsを唯一の資格情報とする。他userはCredential必須。無効/未知user/不正passwordは同じP19_AUTH、存在を漏らすdetail禁止。TOTP未設定ならnull、設定済みなら6桁必須。新user初期passwordはcrypto/rand32 bytesのbase64url、OperationResult.one_timeで1回だけ返し、credentialへ既存KDFで保存する。新規user作成・passwordChange時は既存sessionを全失効。admin以外のTOTPは現profileではnull固定（外部認証/TOTP拡張は別仕様であり暗黙実装しない）。session tokenは既存方式のmemory-only、userIDへ紐付け、再起動で失効。毎requestでenabled/grant/expiryを再検査する。cicd sessionはlegacy sessionとは別namespaceで検証し、逆方向へBearerを流用しない。named userは既存must_change/login_count算出を用い、Root内Credential.login_countは安全整数上限で飽和加算する。forced sessionはpasswordChange/logout以外を403 P19_FORBIDDENで拒否し、promptは操作を妨げない。password変更成功では全sessionを失効して再loginを要求する。login rate limitはIPとuserごと10回/60秒、超過P19_RATE、失敗カウンタ上限4096（期限切れから回収、満杯は新login拒否）。

service tokenは32 random bytes、base64url paddedなし43文字、DBにはSHA256だけ、期限は発行から1秒〜30日、responseで1回だけ返す。Bearerでsessionと同じ入口へ渡し定時間比較、role/grantを毎回再検査する。serviceの承認・user/grant操作は禁止。serviceへのplatform-admin grantは保存拒否し、session tokenとservice tokenは型別に長さ/encodingを検査してから照合する。worker SSH identityはAPI Bearerではない。

<a id="phase-19-auth-runtime-boundary"></a>
**認証runtime境界：** cicd専用のsession/rate/lockout辞書はapi ownerのmemoryに保持し、legacy辞書を参照・更新しない。sessionは4096件以下、期限切れだけを回収し、満杯はP19_QUOTA。token生成は32 bytesのlowercase hex、hash衝突は3回まで再生成しP19_ENTROPY、entropy失敗も同codeとする。候補tokenはcredentials更新前に生成し、credentials確定→必須access/audit→最新認可再取得のすべてが成功した後だけmemoryへ登録する。途中失敗でtokenを返さず、既に保存したlogin_count/TOTP stepを巻き戻さない。user revisionまたはcredential fingerprintが変わったsessionは失効し、login_count/last_login_atの更新だけでは他sessionを失効させない。adminは既存credential/TOTPの現在fingerprintも照合する。全sessionは再起動で消失する。

loginのIP keyはRemoteAddrのIPだけ（forwarded header不使用、不正値はunknown）、user keyは検証済Idとする。IPの10回/60秒をbody decode前、userの同limitをdecode後/credential検査前に適用する。各fixed windowは初回受付時から60秒、now>=終端で更新し、時計逆行は拒否する。rate辞書とpassword失敗辞書はそれぞれ4096 key以下、期限切れだけを回収し、満杯はP19_RATE。password不一致の10回目でIPを600秒lockし、その後はKDF未実行でP19_RATE。未lockの失敗entryは最終失敗から600秒で回収する。password一致で失敗entryを削除し、TOTP不一致/JSON不正/rate拒否をpassword失敗へ加算しない。未知/disabled userも同じ認証拒否とし、不存在入力では固定dummy KDFを実行する。

service Bearerはcanonicalなbase64url無padding43文字/decoded32 bytesのみ、session Bearerはlowercase hex64文字のみとして分離する。service照合は全保存hashとの定時間比較を行い、複数一致は状態破損として拒否する。毎requestで現在enabled/期限/hash/grantを読み直す。logoutは現在sessionを削除し、serviceの場合は現在hashを期限までmemory拒否集合へ追加する（上限4096、満杯P19_QUOTA）。これは永続service無効化ではなく、service拒否集合も再起動で消失する。資格情報を永続失効させる場合はservice更新のenabled=falseを使用する。logout後の必須log失敗でも削除/拒否を戻さない。認証・認可失敗はOperation/Root監査を作らず、既存access/auditの固定fieldだけへ記録し、認証前actorはsystemとする。

<a id="phase-19-permission-matrix"></a>
**権限行列：**

deny-by-default。User/Service IDのsystem/worker/siteは予約済みで作成拒否、adminは初期Userだけ、user/service間のID重複も拒否する。enabledかつexpires_at=nullまたはnow<expiryのGrantだけを採用し、roleの和集合を許可する。platform-adminだけproject_id/environment_id両方null、他roleはproject必須。approverはenvironment必須、他roleのenvironmentはnull。対象外scope組合せは保存拒否。viewer/operator/project-adminの包含関係以外を暗黙包含しない。

| role | 許可 |
|---|---|
| viewer | 自projectの設定（secret値除く）、run/log/artifact/環境/統計の読取、plan preview |
| operator | viewer + run開始/retry/resume/cancel、手動pin、配備plan/開始/cancel/rollback |
| approver | 対象projectのID/nameと対象environmentの配備plan読取、approve/rejectだけ。run/log/secretの読取権限は与えない。 |
| project-admin | operator + definition/template/environment/secretの変更、project内user/serviceへのviewer/operator/approver付与。platform-adminや別project権限を付与できない。 |
| platform-admin | project/worker/target/user/service作成、project-admin付与、global設定・recover。project data操作は対象projectの明示grantも必要。 |

disabled projectは新規run/deploy開始・retry/resumeを拒否するが、初期構成のため設定変更/読取/cancel/reconcileは各role行列の範囲で許可する。project削除は非対応でenabled=falseによる停止だけ、user削除は無効化のみ、最後の有効platform-adminを失効/無効化できない。subject自身の有効grantはAPI Authorization projectionへ渡し、他subjectや無効/期限切れgrantを含めない。1subjectのprojectionがAPI byte上限を超えるgrant変更はP19_QUOTAで拒否する。quota値を現在利用量より下げる変更は409。worker/target変更はdrainかつactive lease/deploymentなしが条件。認証→scope認可→resource lookup→body decodeの順にし、他project IDは403（存在detailなし）。

承認対象はDeployPlan digest、承認有効期限は発行から3600秒、必要人数はEnvironment.approvals（0〜2）。production=trueは1以上かつallow_self_approval=false。同一subject重複不可、service不可。approver grantが失効/無効化、plan内容/環境revision/期待世代変更で承認無効。self禁止時requester本人は拒否。rejectはdeployment.failed/P19_APPROVALとして終端化、承認だけで実体配備はしない。

<a id="phase-19-automatic-subject"></a>
**自動triggerの内部主体：** watch/scheduleのRun.requesterは予約Id systemを使用し、User/Service/Grant/credential/session/tokenを新設しない。systemはcoordinator内部の登録済project trigger専用であり、API login/Bearer/要求bodyから選択できず、platform-admin権限を持つuserの代替でもない。project.enabledと現在参照Definitionの整合を受付commitで確認する。自動配備でもEnvironmentのapproval・production・freeze・target認可境界を省略せず、内部主体だけで承認人数を満たさない。再実行を人間またはserviceが要求した場合は元system主体を継承せず、現在の要求subjectを認可・記録する。保存表現は [docs/details/statefile.md 詳細本文責務 Plan保存整合](statefile.md#phase-19-plan-state-invariants)、trigger受付は [docs/details/runner.md 詳細本文責務](runner.md#phase-19-trigger-admission) を参照する。

<a id="phase-19-authorization-go-boundary"></a>
**認可Go境界：** `AuthorizeCICD(principal,snapshot,access,now)` は、既に資格情報を検証したprincipalと、同一Root世代の現在Subject/Grant projectionだけを使用する純粋関数である。credential検証やsession発行は行わない。principalの期限/must_changeと現在Subjectのenabled/service期限を毎回確認し、非nil失効checkを呼ぶstatefileの [操作transaction境界](statefile.md#phase-19-operation-go-boundary) から使用する。checkの未接続を全許可へ変換しない。未知operation/role/subject種別、不正scope、無効/期限切れgrantは拒否する。permissionと包含関係は [権限行列](#phase-19-permission-matrix)、operation別要求roleは [API catalog](api.md#phase-19-operation-catalog) だけを正本とする。

戻り値はprincipal本人の有効Grantの独立copy、元snapshotのgeneration、principal/serviceの早い方の期限であり、credential/token/他subjectのgrantを含めない。grantPutは変更先Grantのscope/roleを照合し、operationGetは照会対象Operationのrequester/scopeから判定する。入力projectionを変えない。保存済み状態の意味検査と認可を混同せず、認可関数が不正Rootを修復することもない。HTTP login/middlewareの未実装をこの純粋関数の成功で完了扱いにしない。

<a id="phase-19-secret-crypto"></a>
**secret暗号・投入・mask：**

master.jsonのencryption_key/cursor_keyは各crypto/rand32 bytes、相互一致拒否、file0600/parent0700。AES-256-GCM（Go crypto/aes + cipher.NewGCM、tag16 bytes）を使用し、nonceはsecret versionごと12 random bytes、同一keyで既存nonce一致時3回再生成後P19_ENTROPY。AADはUTF-8 `p19-secret-v1\0PROJECT\0ID\0VERSION`（VERSIONは10進）とする。ciphertextはSealのtag付きbytesをbase64。復号失敗はP19_SECRET、平文や暗号文をerrorへ含めない。master自動再生成/backup同梱/鍵なしfallbackは禁止。明示初期化以外はmaster不在で起動拒否。

secret平文はUTF-8一行、1〜4096 bytes、CR/LF/NULなし。入力bodyだけで受け、argv/URL/環境file/状態snapshotに保存しない。versionはRoot.secrets内の既存最大+1、旧版はactive参照がある間削除不可。暗号fileをO_EXCL staging→fsync→非上書きpublish→親fsync後に、SecretMetaと監査を同じRoot transactionへcommitして初めて参照可能にする。secret一覧/復号はRoot.secretsを入口とし、directory scanで未commitの暗号fileを採用しない。Root commit失敗のfileは未参照として新規Put前のGCで回収し、同versionの上書きや既存暗号文の別入力への流用をしない。GCもlock下で全Plan/active参照を再照合してmetadata除去後にfileを削除する。jobのsecret参照はversion固定、secret数16/合計16384 bytes以内。復号はdurable start直前、worker stdinでだけ渡し、利用後byte sliceを上書きする（GCによる全copy消去を保証しない）。deploy key/SSH/masterはjobに渡さない。

mask patternはraw値、standard padded base64、base64url（paddedと無paddingの双方）、url.QueryEscape、url.PathEscapeを重複排除し、byte長降順→byte辞書順で最長一致を `[REDACTED]` に置換する。長さ1のsecretもmaskする。chunk境界処理はrunner正本。mask後だけlog/notification/diagnosticへ渡す。任意変換・暗号化・分割された悪意ある出力の完全検出は保証せず、untrusted userへsecret使用権を付与しない。Node禁止をmaskで代替しない。

auditはOperation受付/拒否/終端、user/grant/secret/worker変更、approve/reject/recoverを同じstate transactionに記録する。秘密値・password・token・nonce・秘密hash・request本文を含めず、subject/scope/action/target/codeだけを残す。監査保存失敗は変更/実行を止める。読取失敗と403はrate-limited診断、認証前subjectはsystemに固定する。公開read modelからcredential/totp/token_sha256/SSH pathを必ず除く。

<a id="phase-19-mask-go-boundary"></a>
**マスク候補のGo境界：** `CICDMaskPatterns(secrets [][]byte)` は [secret暗号・投入・mask](#phase-19-secret-crypto) の値と変換集合だけをbyte列で返す純粋処理とする。最大16値/合計16384 bytes、各値は同正本の平文制約を検査し、空集合は有効、空値/不正UTF-8/禁止改行/容量超過はP19_SECRET。候補は重複を除きbyte長降順/byte辞書順で返し、入力sliceおよび候補同士と共有しない。変換候補の利用/最長一致/未確定suffix/UTF-8処理はrunner所有であり、本関数はlogを書かず置換済み出力を再走査しない。呼出元は元のsecretと候補の保持をjob寿命内へ限定する。

検証は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 横断テスト証跡共通入口](../DETAIL_INDEX.md#cross-test-evidence-route) を参照する。

---

<a id="phase-19-security-implementation-classification"></a>
**実装分類：** RBAC・secret投入/mask・認証保存順は `INHOUSE`、使用primitiveはcrypto/aes、crypto/cipher、crypto/hmac、crypto/rand、既存KDF（`STD`）。入力/状態/異常系は[`docs/details/security.md` 詳細本文責務](#phase-19-secret-crypto)、test vector・fixture・mutationは[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 横断テスト証跡共通入口](../DETAIL_INDEX.md#cross-test-evidence-route)から本ownerのPhase19 caseへ接続する。保守ownerは `security`。削除はcicd-v1を無効化し、active処理と当該schema参照が0、保持データを検証済みexportした場合だけとし、legacy profileへの暗黙移行は行わない。依存採否は[`docs/SPEC.md` ポリシー責務 §4](../SPEC.md#policy-dependencies)を参照する。

<a id="0-責務境界"></a>

**0. 責務境界：**

| 項目 | 内容 |
|------|------|
| owner component | `security` |
| 実装主体 | [`components/security/`](../../components/security/)。認証、認可、token、session、TOTP、rate limit、audit の API 境界は [`components/api/`](../../components/api/) と接続する。 |
| 持つ内容 | `security` owner が主本文として定義する API token scope、API key、audit、session timeout、TOTP、rate limit、漏えい禁止、security 横断順序。 |
| 検証接続 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 owner 詳細本文 検証接続共通入口](../DETAIL_INDEX.md#owner-detail-verification-route) を参照する。 |

---

<a id="対象範囲"></a>
**対象範囲：**

| 範囲 | 内容 |
|------|------|
| [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) | ビルドトリガー専用 API スコープ。 |
| [`docs/details/security.md` 詳細本文責務 §27.43](security.md#sec-27-43) | API キー管理。 |
| [`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) | 監査ログ。 |
| [`docs/details/security.md` 詳細本文責務 §27.45](security.md#sec-27-45) | セッションタイムアウト変更設定。 |
| [`docs/details/security.md` 詳細本文責務 §27.46](security.md#sec-27-46) | TOTP 二要素認証。 |
| [`docs/details/security.md` 詳細本文責務 §27.47](security.md#sec-27-47) | API レート制限。 |

---

<a id="横断固定契約"></a>
**横断固定契約：**

<a id="sec-27-42"></a>
**[`docs/details/security.md` 詳細本文責務 §27.42〜§27.47 認証・監査・制限機能 実装確認固定契約](security.md#sec-27-42)：**

[`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) は、API token、監査、session、TOTP、rate limit に関する安全機能である。各機能の詳細実装確認では、対象の [§27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の機能契約に加えて、この認証・監査・制限機能実装確認固定表を満たす。

| 節 | 機能 | 判定入口 | 成功時副作用 | 失敗時副作用 | 漏えい禁止値 | fixture 証跡入口 |
|----|------|----------|--------------|--------------|--------------|----------------|
| [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) | API token scope | route / method 確定後、body parse 前。 | 許可 endpoint だけ処理し、[`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) の監査対象 event に該当する場合は audit に actor を残す。 | 権限不足は対象処理を実行せず `403`。audit 失敗時は `500`。 | token 本体、Authorization header。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/security.md` 詳細本文責務 §27.43](security.md#sec-27-43) | API key 管理 | admin session または admin scope。 | token hash だけ保存し、作成時だけ token 本体を返す。 | validation 失敗は保存差分なし。失効済み token は再有効化しない。 | token 本体、token hash の不要露出。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) | 監査ログ | security / config / operation event 確定時。 | 1 event 1 JSON Lines で追記し、actor / target / result を保存する。 | 対象 action の audit 失敗は API では `500`、runner 内部 event では runner failure とする。 | secret、password、token、TOTP secret、raw request body。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/security.md` 詳細本文責務 §27.45](security.md#sec-27-45) | session timeout | login、authenticated request、timeout config API。 | session 発行時に `expires_at` を固定し、有効な認証 request ごとに `last_used_at` だけを更新する。 | timeout session は `401`、対象 endpoint は実行しない。 | session token。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/security.md` 詳細本文責務 §27.46](security.md#sec-27-46) | TOTP | setup、confirm、login/totp、disable。 | secret は confirm 成功後だけ有効保存し、ticket は一回だけ使う。 | ticket 再利用、期限切れ、code 不正は対象状態を変更しない。 | TOTP secret、backup code 相当値、ticket token。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/security.md` 詳細本文責務 §27.47](security.md#sec-27-47) | API rate limit | route / auth / scope 判定の定義済み位置。 | 上限未満だけ count を増やし endpoint 処理へ進む。 | `429` は count を増やさず、audit 成功時だけ返す。 | API token、session token、request body。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |

<a id="sec-27-42-2"></a>
**[`docs/details/security.md` 詳細本文責務 §27.42〜§27.47 セキュリティ機能 横断順序固定契約](security.md#sec-27-42-2)：**

| 順序 | 処理 | 固定条件 |
|------|------|----------|
| 1 | route / method を確定する。 | 未定義 route は認証、rate limit、body parse より前に `404` / `405`。 |
| 2 | IP access control を判定する。 | `GET /api/health` だけは省略する。読取、拒否、破損時の副作用は [`docs/details/api.md` 詳細本文責務 §22.0c.1](api.md#sec-22-0c-1) を参照する。 |
| 3 | 認証不要 endpoint を判定する。 | `GET /api/health` と `POST /api/webhook` は個別契約を優先する。Webhook は [`docs/details/api.md` 詳細本文責務 §27.12](api.md#sec-27-12) の body size 判定と署名検証を認証代替とする。 |
| 4 | login rate limit を判定する。 | login group は credential body の parse と認証より前に IP key で判定する。login 以外ではこの段階を省略する。 |
| 5 | 認証情報を検証する。 | 認証必須 endpoint は session と API token を混同しない。形式不一致は `401`。login は検証済み body 上限内で credential body を parse し、個別認証契約を適用する。 |
| 6 | 認可 gate を判定する。 | API token は endpoint scope を判定し、scope 不足は `403 {"error":"Forbidden"}` とする。管理 session は session 発行時の `password_change_required` を判定し、`true` の場合は `POST /api/change-password` と `POST /api/logout` 以外を `403 {"error":"Password change required"}` とする。拒否時は endpoint 固有 body の parse / validation、認証後 rate limit 更新、endpoint 固有状態更新を行わない。認証不要 endpoint では省略する。 |
| 7 | 認証後 rate limit を判定する。 | actor key と IP key を同一 lock 内で判定・更新する。`POST /api/webhook` は署名検証成功後に `trigger` group の IP key だけを判定・更新する。その他の認証不要 endpoint は個別契約で対象 group が明記された場合だけ実行する。 |
| 8 | endpoint 固有 body / query / path validation を行う。 | 失敗時は対象状態、外部 API、外部 command を変更しない。前段階までに完了した security 副作用は保持する。 |
| 9 | endpoint 固有処理を実行する。 | 成功時だけ [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint に固定された保存順で状態、access log、audit log を確定する。 |

[`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) 認証・監査・制限機能実装確認固定表の順序を変更してはならない。対象の [§27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の機能契約はこの順序を上書きできず、異なる順序の記載がある場合は仕様不整合として実装しない。例外が必要な場合は、先に横断順序固定表と対象の [§27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の機能契約を同一仕様変更で改訂し、例外の endpoint、挿入位置、前後の副作用、失敗時の応答を固定表に明記する。その改訂が完了するまで例外は存在しない。


<a id="additional-management-security-order"></a>
**追加管理profileの横断順序適用：** 上表のstep1（route/method）とstep2（IP access control）を実行した後、次表に明示したrouteだけはstep3〜8の代わりに各責務正本を使用する。request ID生成とendpoint実行後の共通ログはAPI正本を維持する。IP拒否は認証/body parse/rate/業務stateへ進まず既存403/503、Origin/schema/認証拒否とその前段副作用は参照先に固定する。この例外をlegacy既存routeやcicd-v1へ拡張しない。

| 対象 | step2後の唯一の順序正本 / 後続境界 |
|---|---|
| multi-user-v1のPOST /api/auth/oidc/start・finish | [OIDC認証確定](#oidc-authentication-commit)。rate後body、flow消費/外部callの失敗境界を同正本で判定する。 |
| multi-user-v1のGET /api/auth/me | [本人HTTP責務](api.md#auth-self-http-contract)。本人sessionとread rateだけ、API tokenは拒否する。 |
| §27.48〜§27.70の追加管理route（OIDC以外、profileに登録済みだけ） | [API追加管理validation順](api.md#additional-management-validation-order)。そのstep6認証・認可は [RBAC](#sec-27-58) / [user token結合](#multi-user-token-totp-contract)、認証後rateはstep6成功直後/step7前でpermission対応scopeと同名groupへ接続する。 |
| version / openapi / auth options / share statusの認証不要route | 上記API validationの例外なし。認証・rate・業務write免除は各HTTP正本だけを適用し、IP access control免除を推測しない。 |

追加管理の認証後scope/rate groupは必要permissionから一意に決定し、user:admin/system:admin→admin、config:write→config、build:write→operate、status:read/config:read→readとする。HTTP GETだからreadへ降格する等の推測は禁止。streamは接続受付1回だけrate countを増やし、keepaliveの再認証でcountを増やさない。forced password sessionのgateは参照先の認証判定内で従来どおり適用する。

---

<a id="認証共通詳細"></a>
**認証共通詳細：**

[`docs/details/security.md` 詳細本文責務 認証共通詳細](security.md#認証共通詳細) は、password 認証、session、login ticket、認証ログ、`--init-credentials` の [`docs/details/security.md` 詳細本文責務](security.md) である。HTTP endpoint の method、path、request、response、status は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) および [`docs/details/api.md` 詳細本文責務 §25](api.md#25-認証-実装仕様) を参照する。`.admin_credentials` schema は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) を参照する。

**password hash 固定契約：**

| 項目 | 仕様 |
|------|------|
| algorithm | `INHOUSE` 分類の `pbkdf2_hmac_sha256_v1`。Go 標準ライブラリ primitive だけで内製実装する。 |
| salt 生成 | `crypto/rand` で 32 bytes を生成し、`encoding/hex` で 64 文字の hex 文字列として保存する。 |
| PBKDF2 PRF | HMAC-SHA256。`crypto/hmac` と `crypto/sha256` だけを使用する。 |
| block 入力 | 32 bytes tag のため block index `1` だけを使用する。`U_1 = HMAC(password_utf8_bytes, salt_bytes || 0x00000001)`、`U_n = HMAC(password_utf8_bytes, U_{n-1})` とし、`U_1` から `U_iterations` までを byte ごと XOR した 32 bytes を tag とする。 |
| 反復 | `iterations = 600000` 固定。 |
| 保存値 | tag 32 bytes を lowercase hex 64 文字列で `.admin_credentials.password_hash` に保存する。 |
| 比較 | 入力 password から同一手順で tag を生成し、hex decode 後の 32 bytes 同士を `crypto/subtle.ConstantTimeCompare` で比較する。 |
| 実装境界 | `crypto/hmac`、`crypto/sha256`、`crypto/rand`、`crypto/subtle`、`encoding/hex` だけで成立する。`golang.org/x/*`、外部 password hash / KDF library、外部 test helper、外部 vector generator を使用してはならない。 |

**password 入力制約：**

| 項目 | 仕様 |
|------|------|
| 最小長 | 8 文字。 |
| 最大長 | 128 文字。 |
| 許可文字 | UTF-8 文字列。NUL 文字は禁止。前後空白はトリムせず、入力値そのものを検証・ハッシュ化する。 |
| 初期 password | `--init-credentials` の非 terminal 標準入力から取得し、平文を保存または出力せず `.admin_credentials` に hash、`must_change:true`、`login_count:0` を生成する。session 発行時の response は session / ticket 固定契約の `must_change` 算出を使用する。 |
| 変更時検証 | `new_password` が現在 password と同一の場合は `422`。 |
| 失敗時応答 | password 不一致は `401 {"error":"Unauthorized"}`。どの条件に失敗したかは返さない。 |
| 成功時保存 | `.admin_credentials` に新 salt、新 hash、`must_change:false`、`updated_at` を原子的に保存する。 |

**session / ticket 固定契約：**

| 項目 | 仕様 |
|------|------|
| session token 生成 | `crypto/rand` で 32 bytes を生成し、`encoding/hex` で 64 文字の lowercase hex 文字列へ変換する。 |
| session 保存場所 | `api` 内のインメモリ辞書。再起動で全 session を破棄する。 |
| session key | token 本体ではなく `sha256(token)` の lowercase hex。 |
| session 期限 | 新規発行時点の `.server_config.session_timeout_seconds`。設定不在時は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の既定値を使用する。 |
| session_id | `crypto/rand` で 16 bytes を生成し、lowercase hex とする。 |
| `must_change` 算出 | session token 発行と同じ credentials 更新で `login_count` を飽和加算した後に算出する。`.admin_credentials.must_change=false` は `"none"`、`true` かつ増加後 `login_count` が 1〜4 は `"prompt"`、`true` かつ増加後 `login_count` が 5 以上は `"forced"` とする。 |
| password 変更 gate | session record は発行 response の `must_change` が `"forced"` の場合だけ `password_change_required:true` をメモリ上に保持する。password 変更成功時に現 session の値を `false` にし、その他の session は失効する。 |
| login ticket | TOTP 有効 login の password 成功時だけ生成し、メモリ上で hash 化して保持する。 |
| login ticket 保持値 | password 検証済みであること、ticket hash、発行時刻、期限、credentials fingerprint だけを保持する。`must_change`、`login_count`、`last_login_at`、`password_change_required` は保持しない。 |
| credentials fingerprint | ticket 発行時の `password_hash`、`salt`、`algorithm`、`iterations` の 10 進 ASCII、`updated_at` をこの順で LF 1 byte で連結した UTF-8 byte 列の SHA-256 lowercase hex。TOTP 成功処理で最新 credentials から再算出し、`crypto/subtle.ConstantTimeCompare` で一致しない場合は ticket を消費して `401 {"error":"Unauthorized"}` とする。fingerprint は response と log に出力しない。 |
| 永続化禁止 | session token、login ticket、setup 仮 secret、連続失敗回数はファイル保存しない。 |
| response 禁止 | password hash、salt、session token hash、ticket hash、TOTP secret 保存値は response に含めない。 |
| log 禁止 | password、current_password、new_password、token、ticket、hash、salt、TOTP code は `.access_log`、`.audit_log`、`.api_access_log`、journal に含めない。 |

`GET /api/sessions` は `session_id` ではなく `current`、`created_at`、`expires_at`、`last_used_at` のみ返す。認証必須 API で有効 token を受信した場合、`last_used_at` を現在時刻へ更新する。期限切れ token は検出時にメモリから削除する。`POST /api/logout` は対象 token のみ削除する。`POST /api/sessions/revoke-all` は現在 token 以外を削除する。

**login 失敗制御：**

| 項目 | 仕様 |
|------|------|
| key | [`docs/details/security.md` 詳細本文責務 §27.47](security.md#sec-27-47) の IP key と同じ正規化を使用する。`RemoteAddr` を `net.SplitHostPort` で host に正規化し、parse 失敗時は同節の固定 fallback key を使用する。 |
| 失敗記録 | `POST /api/login` の password hash 不一致だけが、IP key ごとの連続失敗回数と最終失敗時刻をメモリ上で更新する。JSON / body 検証失敗、pre-auth rate limit 拒否、TOTP 失敗はこの回数を更新しない。再起動で全 entry を破棄する。 |
| lock 成立 | 10 回目の連続 password 不一致は `401 {"error":"Unauthorized"}` を返し、`locked_until = last_failure_at + 10分` を記録する。その後 `now < locked_until` の request は password hash 検証を行わず `429 {"error":"Too many attempts"}` で拒否する。 |
| lock 解除 | `now >= locked_until` であれば対象 entry を削除してから password hash 検証を行う。password hash 一致時も、TOTP の有効・無効にかかわらず対象 entry を削除する。 |
| rate limit 優先順位 | pre-auth login rate limit を login lock より前に判定する。両方が拒否条件を満たす場合は `429 {"error":"Too many requests"}` を返し、login lock 判定、password hash 検証、失敗回数更新を行わない。 |
| 応答時間 | password 不一致、存在しない credentials、lock 中を除く検証失敗では、条件の詳細を response へ出さない。 |
| log | 成功、失敗、lock 拒否はいずれも `.access_log` へ追記する。password、token、hash、salt は記録しない。 |

**認証処理の副作用境界：**

| ケース | HTTP status | `.admin_credentials` | メモリ session / ticket | `.access_log` | `.audit_log` | 備考 |
|--------|-------------|----------------------|-------------------------|---------------|--------------|------|
| password 不一致 | `401` | 変更なし | 変更なし | `login_failure` を追記 | `login_failure` を追記 | 連続失敗回数はメモリ上で +1。 |
| 連続失敗 lock | `429` | 変更なし | 変更なし | `login_locked` を追記 | `permission_denied` を追記 | password hash 検証は実行しない。 |
| password 成功 / TOTP 無効 | `200` | `login_count`、`last_login_at` 更新 | session 追加 | `login_success` を追記 | `login_success` を追記 | session token は全永続ログに保存しない。 |
| password 成功 / TOTP 有効 | `200` | 変更なし | ticket 追加 | `totp_required` を追記 | `totp_required` を追記 | `login_count` と `last_login_at` は更新しない。 |
| TOTP code 不一致 | `401` | 変更なし | ticket 削除 | `totp_failure` を追記 | `totp_failure` を追記 | ticket は再利用不可。 |
| TOTP 成功 | `200` | `login_count`、`last_login_at` 更新 | ticket 削除、session 追加 | `login_success` を追記 | `login_success` を追記 | session 期限は成功時点の設定で決める。 |
| session 期限切れ | `401` | 変更なし | 対象 session 削除 | 追記しない | 追記しない | `.api_access_log` は通常 API request として記録する。 |
| logout | `200` | 変更なし | 対象 session 削除 | `logout` を追記 | `logout` を追記 | token 本体と token hash は保存しない。 |
| password 変更成功 | `200` | 新 salt / hash、`must_change:false`、`updated_at` 更新 | 現 session の `password_change_required=false`、現 session 以外削除 | `password_change` を追記 | `password_change` を追記 | 新旧 password、hash、salt は保存しない。 |
| 他 session 一括失効 | `200` | 変更なし | 現 session 以外削除 | `session_revoke_all` を追記 | `session_revoke_all` を追記 | 現 session と実行中 response は維持する。 |

[`docs/details/security.md` 詳細本文責務 認証共通詳細](security.md#認証共通詳細) の固定表で `.access_log` または `.audit_log` の追記が必要な処理は、response 返却前に追記を完了する。追記失敗時は、[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint または対象の [§27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) 機能契約で別の応答が定義されていない限り `500 {"error":"Internal server error"}` を返す。session token、login ticket、TOTP setup secret は、必要なログ追記がすべて成功するまで response に含めてはならない。

**認証副作用順序固定契約：**

| 処理 | 固定順序 | 失敗時 |
|------|----------|--------|
| `POST /api/login` password 不一致 | pre-auth rate limit → IP key の lock entry 期限判定 → credentials 読込 → hash 比較 → 失敗回数更新 → `.access_log` → `.audit_log` → `401` response | 10 回目も `401`。ログ追記失敗は `500`。失敗回数は戻さない。session/ticket は作成しない。 |
| `POST /api/login` password 成功 / TOTP 無効 | pre-auth rate limit → IP key の lock entry 期限判定 → credentials lock 取得 → credentials 読込 → hash 比較 → IP key の失敗 entry 削除 → session token 生成 → `login_count=min(login_count+1,9223372036854775807)` と `last_login_at` 更新 → `must_change` 算出 → credentials 保存 → `.access_log` → `.audit_log` →同値の `password_change_required` を持つ session record 追加 → token response | token 生成失敗は credentials 更新前の `500`。credentials またはログ追記失敗も `500` とし、token を response せず session record を追加しない。ログ失敗前に保存済みの credentials は戻さない。 |
| `POST /api/login` password 成功 / TOTP 有効 | pre-auth rate limit → IP key の lock entry 期限判定 → credentials 読込 → hash 比較 → IP key の失敗 entry 削除 → password 検証済み ticket 生成 → `.access_log` → `.audit_log` → ticket response | `.admin_credentials`、`login_count`、`last_login_at` は更新せず、ticket に `must_change` を保存しない。ログ追記失敗は `500`。ticket は保存しない。 |
| `POST /api/login/totp` 成功 | ticket 検証・即時無効化 → TOTP secret 読込 → code 検証 → credentials lock 取得 → credentials 再読込 → credentials fingerprint 一致判定 → session token 生成 → `.totp_secret.last_accepted_step` 更新 → `login_count=min(login_count+1,9223372036854775807)` と `last_login_at` 更新 → `must_change` 算出 → credentials 保存 → `.access_log` → `.audit_log` →同値の `password_change_required` を持つ session record 追加 → token response | fingerprint 不一致は状態更新前の `401`。その他の途中失敗時は token を response せず session record を追加しない。ticket は成功/失敗いずれも再利用不可とし、ticket 発行時の credentials 値から `must_change` を推測しない。失敗地点より前に永続化済みの `.totp_secret` または credentials は戻さない。 |
| `POST /api/logout` | token 認証 → 対象 session 削除 → `.access_log` → `.audit_log` → response | ログ追記失敗は `500`。削除済み session は戻さない。 |
| `POST /api/change-password` | token 認証 → current password 検証 → new password 検証 → salt/hash 生成 → `.admin_credentials` 更新 → 現 session の `password_change_required=false` → 現 session 以外削除 → `.access_log` → `.audit_log` → response | credentials 更新失敗は session を変更しない。ログ追記失敗時は `500` だが更新済み credentials と session 変更は戻さない。 |
| `POST /api/sessions/revoke-all` | token 認証 → 現 session 以外を削除 → `.access_log` → `.audit_log` → response | ログ追記失敗時は `500`。削除済み session は戻さない。 |

<a id="security-auth-concurrency-contract"></a>
**認証 transaction・並行更新固定契約：**

security owner は、credentials、TOTP、session、login ticket、TOTP setup 仮 secret、login 失敗 entry をまたぐ認証状態遷移を process 内の単一 auth transaction coordinator で直列化する。対象は `POST /api/login`、`POST /api/login/totp`、`POST /api/logout`、`POST /api/change-password`、`GET /api/sessions`、`POST /api/sessions/revoke-all`、`GET /api/auth/totp-status`、`POST /api/auth/totp-setup`、`POST /api/auth/totp-confirm`、`DELETE /api/auth/totp`、session 認証時の期限切れ削除と `last_used_at` 更新である。同じ request から coordinator を再帰取得してはならない。

auth transaction coordinator の取得待ちは最大 60 秒とする。待機開始を elapsed `0` とし、elapsed `<60s` でだけ取得成功を許可する。elapsed `60s` 到達時は timeout を優先し、同時刻の解放を取得成功扱いにしない。timeout 時は endpoint 固有 body parse、credentials / TOTP read、entropy 取得、memory state 変更、状態 write、`.access_log`、`.audit_log` を開始せず、`409 {"error":"Conflict"}` を返す。request context が取得前に cancel された場合も副作用なしで終了する。取得後は、対象処理の認証副作用順序固定契約が完了するまで coordinator を保持し、response 値を確定した後に解放する。

auth transaction coordinator は、statefile locked update adapter 内の file lock と memory auth lock のどちらよりも先に取得し、どちらを解放した後にも最後まで保持する。file lock と memory auth lock を同時に保持してはならず、一方を保持中に他方を取得してはならない。同じ transaction で両方が必要な場合は、memory auth lock で必要値を read-copy して解放 → statefile locked update adapter を完了して file lock を解放 → memory auth lock で確定済み mutation を適用して解放、の順に固定する。coordinator が他 request の認証状態遷移を排他するため、各区間の間に同じ auth 状態を別 request が変更することはない。filesystem、network、command、log write、password hash、TOTP 計算、entropy 読取は memory auth lock 外で実行する。memory auth lock 内では session、ticket、pending TOTP、login failure / rate entry の read-copy または確定済み mutation だけを行う。auth transaction coordinator は statefile I/O と log write をまたいで保持できるが、別 request の response writer、network client、systemd command を待ってはならない。

| 競合処理 | 直列化後の固定結果 |
|----------|--------------------|
| 同時 password login 成功 | 先に coordinator を取得した request から 1 件ずつ最新 credentials を読み、各成功ごとに `login_count` を飽和加算する。lost update を禁止する。 |
| login と password change | password change の credentials rename が先なら旧 password login は `401`。login の session 発行が先なら password change はその session を使用でき、変更成功時に他 session を失効する。 |
| TOTP ticket と credentials 変更 | ticket の credentials fingerprint と成功時再読込値が一致しない場合は `401`。旧 revision の ticket から session を発行しない。 |
| 同時 TOTP setup | 後に coordinator を取得した setup が pending secret を置換する。confirm は request 開始時ではなく coordinator 取得後の最新 pending secret だけを検証し、置換前 secret の code は `409`。 |
| TOTP confirm と disable | 先に完了した transaction の永続状態を後続 request が再読込する。後続 request は stale memory snapshot を使用しない。 |
| revoke-all と通常認証 | revoke-all が先なら削除対象 token の後続認証は `401`。通常認証が先ならその request だけ完了でき、revoke-all 完了後の次 request は `401`。 |
| session 期限切れと logout | 最初の transaction が対象 session を削除し、後続 request は `401`。同じ session を二重削除、二重 audit しない。 |

credentials、TOTP、API token の read-modify-write は [`docs/details/statefile.md` 詳細本文責務 lock 内 read-modify-write adapter 固定契約](statefile.md#statefile-locked-update-adapter-contract) を使用する。transaction 開始前の read-only snapshot に基づく保存、file lock 外の `login_count` 加算、TOTP replay step 更新、token `last_used_at` 更新を禁止する。

session token と login ticket は `crypto/rand` 成功後にだけ生成し、生成した値はメモリ上で hash 化して保持する。response body に含める token / ticket は、その request の成功 response 1 回だけに含める。`403`、`429`、`500`、network 切断検出時に、未送信 token をログや状態ファイルへ退避してはならない。

<a id="init-credentials-cli-contract"></a>
**`--init-credentials` CLI 固定契約：**

初期 password は非 terminal の標準入力から exactly 1 行で受け取る。入力 byte 列は UTF-8 password byte 列と末尾 LF 1 byte だけで構成し、末尾 LF より後の byte、途中の LF、CR、NUL、UTF-8 不正を禁止する。末尾 LF は password に含めず、前後空白を trim しない。UTF-8 code point 数は 8〜128 とする。実装は最大 514 bytes まで読み、513 bytes 以下で exactly 1 個の末尾 LFを確認してから UTF-8 と code point 数を検証する。標準入力が terminal の場合は password を読み取らない。password byte 列は hash 入力以外に複製、保存、log 出力せず、処理終了前に保持 buffer を上書きする。

`--state-dir` の検証は、security 処理を開始する前に [`docs/details/api.md` 詳細本文責務 `api` CLI 固定契約](api.md#api-cli-contract) と [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 CLI state directory 共通固定契約](../DETAIL_INDEX.md#common-state-dir-contract) に従って完了させる。security owner は state directory の parse、検証順、エラー出力を再定義しない。

| init-credentials CLI ケース | stdout | stderr | 終了コード | 副作用 |
|----------------------------|--------|--------|------------|--------|
| 新規生成成功 | `credentials initialized` + LF | 空 | `0` | `.admin_credentials` を mode `0600` で作成する。 |
| 既存あり | 空 | `credentials already exist` + LF | `2` | 既存ファイルを変更しない。 |
| 標準入力が terminal | 空 | `initial password stdin must not be terminal` + LF | `2` | password 読取とファイル作成なし。 |
| password 入力が形式不正 | 空 | `invalid initial password` + LF | `2` | salt / hash 生成とファイル作成なし。 |
| password 読取失敗 | 空 | `password input failed` + LF | `1` | salt / hash 生成とファイル作成なし。 |
| 書込失敗 | 空 | `credentials write failed` + LF | `1` | tmp を削除し、部分ファイルを残さない。 |
| rand 失敗 | 空 | `random source failed` + LF | `1` | ファイル作成なし。 |

生成手順は、API CLI の state dir 検証 → `.admin_credentials` の予備存在確認 → 標準入力の terminal 判定 → password 読取・検証 → salt 生成 → hash 生成 → [`docs/details/statefile.md` 詳細本文責務 §22.0a 状態ファイル更新手順](statefile.md#statefile-update-procedure) の create-only mode 呼出し、の順に固定する。予備確認で既存 file を検出した場合は標準入力を読まず「既存あり」で終了する。create-only mode は lock 取得後に存在を再確認するため、予備確認後の競合でも既存 `.admin_credentials` を上書きしない。`ErrStateAlreadyExists` は「既存あり」の終了コード `2`、lock / tmp / write / chmod / file sync / rename / parent sync / cleanup failure は「書込失敗」の終了コード `1` へ写像する。rename 後の partial failure では作成済みファイルを残す。security owner は tmp 名、lock 名、rename、cleanup を再定義しない。初期 password を process 引数、環境変数、設定ファイル、固定値、terminal 対話、ランダム生成、stdout / stderr 出力から取得してはならない。

**認証共通実装確認ゲート：**

| 観点 | 合格条件 | 禁止条件 |
|------|----------|----------|
| one-time response | session token、login ticket、API token 本体、TOTP setup secret、otpauth URI は、該当成功 response 1 回だけに含める。 | `500`、`401`、`403`、`409`、`422`、`429` response、log、状態ファイル、fixture expected へ平文を残すこと。 |
| memory-only state | session、login ticket、TOTP setup 仮 secret、login 失敗回数は process memory だけに保持し、再起動で破棄する。 | session、login ticket、TOTP setup 仮 secret、login 失敗回数を保存する未定義永続ファイルの作成、ticket / session / 仮 secret の backup / restore 対象化。 |
| required security log before token | token / ticket / secret を response に含める前に、endpoint 固有契約が必須とする `.access_log` と `.audit_log` の追記を完了する。`.api_access_log` は [`docs/details/api.md` 詳細本文責務 §27.6](api.md#sec-27-6) の best-effort 記録とし、one-time 値の返却 gate にしない。 | 必須 `.access_log` または `.audit_log` の追記失敗時に token / ticket / secret を response へ含めること。 |
| hash-only storage | password、session token、login ticket、API token は保存時に hash 化し、平文を保存しない。 | hash 算出入力の平文、token 本体、ticket 本体、password 本体を expected / log に保存すること。 |
| fixed error body | 認証失敗、権限不足、rate limit、validation failure は固定 error body だけを返す。 | password 不一致理由、token record 詳細、scope 一覧、TOTP step、rate limit key を response に出すこと。 |
| no endpoint side effect | 認証、scope、rate limit、body validation のいずれかで失敗した場合、endpoint 固有処理を開始しない。 | 状態ファイル更新、外部 API / command 実行、通知送信、snapshot / archive / deploy 操作。 |
| audit dependency | 必須 audit が定義された操作は、audit 追記成功まで完了 response を返さない。 | 必須 audit 失敗時に成功 response を返すこと、audit failure を audit に再帰記録すること。 |
| secret scan | response、stdout、stderr、journal、`.access_log`、`.audit_log`、`.api_access_log`、`.config_log`、fixture expected に禁止値が存在しないことを検証する。 | 目視確認だけで secret 非表示を合格扱いにすること。 |

**one-time response 固定契約：**

| 値 | 返却 endpoint | 保存先 | 成功 response 以外の扱い |
|----|---------------|--------|--------------------------|
| session token | `POST /api/login`、`POST /api/login/totp` | memory に `sha256(token)`。 | response 前の必須 `.access_log` / `.audit_log` 失敗時は token を破棄し、`500`。 |
| login ticket | TOTP 有効時の `POST /api/login` | memory に `sha256(ticket)`。 | response 前の必須 `.access_log` / `.audit_log` 失敗時は ticket を破棄し、`500`。期限切れ / 失敗後は再利用不可。 |
| API token 本体 | `POST /api/tokens` | `.api_tokens.token_hash` のみ。 | `.api_tokens` 保存後の必須 `.access_log` / `.audit_log` 失敗時は record を残し、token 本体は返さず `500`。 |
| TOTP setup secret | `POST /api/auth/totp-setup` | memory の仮 secret。confirm 成功後だけ `.totp_secret`。 | setup response 以外へ出さない。audit / access / server log / status response へ出さない。 |
| otpauth URI | `POST /api/auth/totp-setup` | 永続保存しない。 | UI 一回表示以外へ出さない。fixture expected では secret 部分を `***` として扱う。 |

認証共通 fixture 名、入力、expected、合格条件、実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) を唯一の正本とする。[`docs/details/security.md`](security.md) 詳細本文責務では、one-time response、memory-only state、hash-only storage、audit dependency、secret scan の実装契約だけを扱う。

---

<a id="sec-27-42-3"></a>
**27.42 ビルドトリガー専用 API スコープ：**
[`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42-3) の境界は owner component `security`、collaborator component `api`、`sdk`、`ui`、`statefile` とする。


本機能の目的は、外部システムによる build 開始操作を `trigger` scope の API token と build 開始 endpoint だけに限定することである。

<a id="api-token-scope-contract"></a>
**scope：**

| scope | 許可 |
|-------|------|
| `trigger` | `POST /api/build`、`POST /api/build/force` のみ。 |
| `read` | 読み取り endpoint のみ。 |
| `operate` | 運用操作。 |
| `config` | 設定変更。 |
| `admin` | 管理操作。 |

`trigger` scope は cancel、queue clear、config、token、audit、PAT 更新を許可しない。

**endpoint group 対応：**

| scope | 許可 endpoint |
|-------|---------------|
| `trigger` | `POST /api/build`, `POST /api/build/force` |
| `read` | `GET /api/status`, `GET /api/logs`, `GET /api/logs/search`, `GET /api/logs/export`, `GET /api/history`, `GET /api/history/export`, `GET /api/history/{id}/log`, `GET /api/history/{id}/comment`, `GET /api/sysinfo`, `GET /api/health`, `GET /api/schedule`, `GET /api/notify-config`, `GET /api/notify-log`, `GET /api/config`, `GET /api/config-log`, `GET /api/pat-status`, `GET /api/stats`, `GET /api/stats/timeline`, `GET /api/stats/build-duration`, `GET /api/stats/build-trends`, `GET /api/output-meta`, `GET /api/repo-info`, `GET /api/branch-config`, `GET /api/backup`, `GET /api/dashboard`, `GET /api/diagnostics`, `GET /api/rate-limit`, `GET /api/disk-usage`, `GET /api/webhook-events`, `GET /api/webhook-config`, `GET /api/snapshots`, `GET /api/snapshots/{id}/download`, `GET /api/maintenance`, `GET /api/access-control`, `GET /api/hooks`, `GET /api/hooks/{id}/log`, `GET /api/alert-rules`, `GET /api/tag-rules`, `GET /api/pipeline-config`, `GET /api/build-chain-config`, `GET /api/notes`, `GET /api/smtp-config`, `GET /api/queue`, `GET /api/approvals`, `GET /api/dashboard-layout` |
| `operate` | `POST /api/build/cancel`, `GET /api/build/stream`, `POST /api/notify-test`, `POST /api/notify/weekly-summary`, `POST /api/pat-verify`, `POST /api/circuit-breaker/reset`, `DELETE /api/queue`, `POST /api/smtp-test`, `POST /api/verify-output`, `POST /api/history/{id}/rollback`, `POST /api/approvals/{id}/approve`, `POST /api/approvals/{id}/reject` |
| `config` | `POST /api/schedule/interval`, `POST /api/schedule/pause`, `POST /api/schedule/resume`, `POST /api/schedule/allowed-hours`, `POST /api/schedule/force-interval`, `POST /api/schedule/cooldown`, `POST /api/notify-config`, `POST /api/config/validate`, `POST /api/config`, `POST /api/log-level`, `POST /api/pat-update`, `POST /api/repo-config`, `POST /api/branch-config`, `POST /api/restore`, `POST /api/webhook-config`, `DELETE /api/snapshots/{id}`, `POST /api/maintenance/enable`, `POST /api/maintenance/disable`, `POST /api/access-control`, `POST /api/hooks`, `DELETE /api/hooks/{id}`, `POST /api/alert-rules`, `DELETE /api/alert-rules/{id}`, `POST /api/tag-rules`, `DELETE /api/tag-rules/{id}`, `POST /api/pipeline-config`, `POST /api/build-chain-config`, `POST /api/notes`, `POST /api/smtp-config`, `POST /api/dashboard-layout`, `POST /api/history/{id}/comment`, `POST /api/history/{id}/flag`, `POST /api/history/{id}/tags`, `POST /api/logs/cleanup`, `POST /api/logs/archive` |
| `admin` | `GET /api/access-log`, `GET /api/api-access-log`, `GET /api/audit-log`, `GET /api/api-rate-limit`, `POST /api/api-rate-limit`, `GET /api/sessions`, `POST /api/sessions/revoke-all`, `GET /api/auth/totp-status`, `POST /api/auth/totp-setup`, `POST /api/auth/totp-confirm`, `DELETE /api/auth/totp`, `GET /api/tokens`, `POST /api/tokens`, `DELETE /api/tokens/{id}` |

multi-user-v1追加管理routeのscopeは [user token結合契約](#multi-user-token-totp-contract) の明示対応だけを使用し、以下のlegacy endpoint閉集合との混在を避ける。既存routeは同対応契約で明示した例外以外、この閉集合を維持する。

管理 session は API token scope 固定表の対象外とするが、[`docs/details/security.md` 詳細本文責務 §27.42〜§27.47 セキュリティ機能 横断順序固定契約](security.md#sec-27-42-2) の `password_change_required` 判定を常に適用する。API token が複数 scope を持つ場合は、いずれか 1 つの scope が endpoint に一致すれば許可する。`POST /api/login`、`POST /api/login/totp`、`POST /api/logout`、`POST /api/change-password` は API token scope 判定の対象外とし、API token では使用できない。`POST /api/webhook` は GitHub Webhook secret 検証専用であり、API token では使用できない。[`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の scope 固定表に存在しない endpoint は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) と [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の scope 固定表へ追加されるまで API token では許可してはならない。

**正常系：**

1. `POST /api/tokens` は `scopes:["trigger"]` を受け付ける。
2. 発行 token は `Authorization: Bearer <token>` で使用する。
3. `trigger` token による build は `.audit_log` と `.build_logs/{id}.json.trigger_actor` に token id を記録する。

**判定順：**

1. Bearer token の形式を検査する。`act_` prefix の場合は API token、64 文字 hex の場合は session token として扱う。どちらにも一致しない場合は `401`。
2. API token の hash、期限、失効状態を検証する。
3. endpoint group 対応表で scope を判定する。
4. 許可時だけ endpoint 処理へ進む。
5. 拒否時は `.audit_log` に `permission_denied`、`actor_type:"api_token"`、`actor_id:{token_id}`、`target_type:"endpoint"`、`target_id:{METHOD + " " + path}`、`result:"denied"` を追記してから `403` を返す。

`permission_denied` の監査ログ追記に失敗した場合は `500` を返し、対象 endpoint の処理は実行しない。

**認証・rate limit 組み合わせ順：**

認証、API token scope、login rate limit、認証後 rate limit、endpoint 固有 validation の実行順は [`docs/details/security.md` 詳細本文責務 §27.42〜§27.47 セキュリティ機能 横断順序固定契約](security.md#sec-27-42-2) を唯一の正本とする。scope 判定前に endpoint 固有の request body parse、状態ファイル更新、外部送信を実行してはならない。

**scope 判定固定条件：**

| 条件 | 判定 |
|------|------|
| session token | API token scope 表は参照せず、[`docs/details/security.md` 詳細本文責務 §27.42〜§27.47 セキュリティ機能 横断順序固定契約](security.md#sec-27-42-2) の `password_change_required` 判定後に許可する。 |
| API token の `scopes` が空 | token record 破損として扱い `500`。 |
| API token の `scopes` に未知値 | token record 破損として扱い `500`。 |
| endpoint が複数 group に現れる | 仕様不整合とし、group の優先順を推測せず実装しない。[`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の endpoint 割当ては method と正規化 path pattern の組で正確に 1 group だけに属しなければならない。 |
| path parameter 付き endpoint | 正規化 path pattern で scope 判定する。例: `DELETE /api/tokens/abc` は `DELETE /api/tokens/{id}`。 |
| query string | scope 判定には使用しない。method と path pattern だけで判定する。 |
| `GET /api/health` | 認証不要。API token scope 判定を行わない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| trigger token で `POST /api/build` | `202`。 |
| trigger token で `GET /api/status` | `403`。`read` を併用した場合のみ成功。 |
| trigger token で config | `403`。 |
| 未定義 endpoint | token scope に関係なく `404` または `405`。 |
| 複数 scope | いずれかに一致する endpoint だけ成功。 |
| path parameter endpoint | pattern 正規化後の scope で判定する。 |
| query 付き endpoint | query を除外して scope 判定する。 |
| scope 前 body | 権限不足時に request body validation や状態更新を行わない。 |

<a id="sec-27-43"></a>
**27.43 API キー管理：**
[`docs/details/security.md` 詳細本文責務 §27.43](security.md#sec-27-43) の境界は owner component `security`、collaborator component `api`、`sdk`、`ui`、`statefile` とする。


本機能の目的は、API key の発行、一覧、失効、期限、scope を実装し、key 本体を保存しないことである。

**API 権限：**

| endpoint | 必要認証 |
|----------|----------|
| `GET /api/tokens` | 管理 session または `admin` scope API token |
| `POST /api/tokens` | 管理 session または `admin` scope API token |
| `DELETE /api/tokens/{id}` | 管理 session または `admin` scope API token |

`admin` scope API token は、`POST /api/tokens` の認証と scope 判定を通過した場合だけ新しい `admin` scope API token を発行する。発行者 token と発行対象 token は別 record とし、親子関係は保存しない。

**正常系：**

1. `POST /api/tokens` は `label`、`scopes`、`expires_at` を受け取る。
2. token 本体は `io.ReadFull(crypto/rand.Reader, buffer)` で取得した 32 bytes を `base64.RawURLEncoding.EncodeToString` で符号化し、`act_` を前置した文字列とする。
3. `.api_tokens` には `sha256(token)` の lowercase hex だけを保存する。
4. response の `token` は作成時 1 回だけ返す。`GET /api/tokens` では返さない。
5. 認証時は hash 一致、`revoked_at == null`、`expires_at == null または now < expires_at` を満たす token だけ有効とする。
6. 作成、認証成功、失効、期限切れ拒否を `.audit_log` へ記録する。ただし token 本体は記録しない。

**`.api_tokens` schema 参照：**

`.api_tokens` の top-level object、record のキー、型、必須性、許容値、正規化例外、採番固定値は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.api_tokens` schema を唯一の正本とする。[`docs/details/security.md`](security.md) 詳細本文責務は token の生成、hash、認証、scope 判定、作成・失効・認証時の副作用と漏えい禁止だけを定義し、schema を再掲しない。

`.api_tokens` が同 schema に違反する場合は破損として扱い、token 認証、一覧、作成、失効をすべて `500 {"error":"Internal server error"}` で拒否する。破損内容、hash、token 本体は response、`.access_log`、`.audit_log`、journal に出力しない。

**API token 認証時の副作用境界：**

| ケース | HTTP status | `.api_tokens` | `.access_log` | `.audit_log` | 備考 |
|--------|-------------|---------------|---------------|--------------|------|
| hash 不一致 | `401` | 変更なし | 追記しない | 追記しない | token id が特定できないため監査対象外。 |
| 期限切れ | `401` | `last_used_at` 更新なし | `token_expired` を追記 | `token_expired` を追記 | token id が特定できた場合だけ記録する。 |
| 失効済み | `401` | `last_used_at` 更新なし | `token_revoked_reject` を追記 | `token_revoked_reject` を追記 | token 本体と hash は保存しない。 |
| scope 不足 | `403` | `last_used_at` 更新なし | `permission_denied` を追記 | `permission_denied` を追記 | 対象 endpoint は実行しない。 |
| 認証成功 | endpoint 固有 | `last_used_at` を現在時刻へ更新 | `token_auth` を追記 | `token_auth` を追記 | endpoint 実行前に更新する。 |

[`docs/details/security.md` 詳細本文責務 §27.43](security.md#sec-27-43) の固定表で `.api_tokens` 保存後に `.access_log` または `.audit_log` 追記へ失敗した場合、対象 API は `500` とし、endpoint 固有処理へ進まない。`last_used_at` 更新済み record は巻き戻さない。

**処理順：**

`POST /api/tokens` は以下の順で処理する。

1. 認証と `admin` scope を確認する。
2. `label`、`scopes`、`expires_at` を検証する。
3. `.api_tokens` のファイルロックを取得する。
4. 既存 record を読み込み、id を採番する。
5. token 本体を生成し、hash を算出する。
6. record を append して `.api_tokens` を原子的に保存する。
7. `.access_log` に `token_create` を追記する。
8. `.audit_log` に `token_create` を追記する。
9. response に token 本体を 1 回だけ含めて返す。

`.api_tokens` 保存後に `.access_log` または `.audit_log` 追記へ失敗した場合、API は `500` を返す。作成済み token record は削除せず、再実行時は新しい token を発行する。

`DELETE /api/tokens/{id}` は以下の順で処理する。

1. 認証と `admin` scope を確認する。
2. path `id` を検証する。
3. `.api_tokens` をロックして対象 record を検索する。
4. 対象が存在しない、または `revoked_at != null` の場合は `404` を返す。
5. `revoked_at` を現在時刻に設定して保存する。
6. `.access_log` と `.audit_log` に `token_revoke` を追記する。
7. `{ "message": "Token revoked" }` を返す。

認証に使用中の API token 自身を path `id` に指定した場合、その token を失効対象にする。当該リクエストは成功し、次リクエストから `401` になる。

**token ID 採番・返却固定契約：**

| 項目 | 仕様 |
|------|------|
| 採番元 | `.api_tokens.tokens[].id` の数値 suffix 最大値。存在しない場合は `tok000001`。 |
| 衝突時 | 最大 suffix + 1 を採用する。削除済みや失効済み id は再利用しない。 |
| token 本体 | `io.ReadFull(crypto/rand.Reader, buffer)` で 32 bytes を完全に取得し、`base64.RawURLEncoding.EncodeToString` で padding なしの 43 ASCII 文字へ変換し、`act_` を前置する。完成 token は合計 47 bytes、本体 alphabet は `[A-Za-z0-9_-]` に固定する。 |
| 作成 response | `{ "id", "label", "scopes", "created_at", "expires_at", "revoked_at", "last_used_at", "token" }`。`token` はこの response だけに含める。 |
| 一覧 response | `tokens` 配列に `token_hash` と `token` を含めない。並び順は `created_at` 降順、同時刻は `id` 昇順。 |
| label 正規化 | 前後空白を除去し、内部空白は保持する。64 文字判定は Unicode code point 数ではなく UTF-8 byte 数で行う。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| scopes 空 | `422`。 |
| 未知 scope | `422`。 |
| label 空 | `422`。 |
| expires_at が現在以前 | `422`。 |
| token 上限 100 件超過 | `422`。失効済み record も件数に含める。 |
| 期限切れ token | `401`。 |
| 失効済み token 再失効 | `404` または冪等成功にせず `404` 固定。 |
| `.api_tokens` 破損 | `500`。自動再生成しない。 |
| `crypto/rand.Reader` が 32 bytes を取得できない | `500 {"error":"Internal server error"}`。token、hash、record、access log、audit log を作成せず、lock を解放する。短い読み取りを成功扱いしない。 |
| 生成した `sha256(token)` が既存 record の `token_hash` と一致 | `500 {"error":"Internal server error"}`。自動再生成せず、既存 record と全 log を変更せず lock を解放する。token または hash の値を出力しない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 作成 | token 本体は response だけ。 |
| 一覧 | token hash と本体を返さない。 |
| 期限切れ | `401`、last_used_at 更新なし。 |
| 失効 | `revoked_at` 保存、以後 `401`。 |
| 自己失効 | 失効リクエストは `200`、以後同 token は `401`。 |
| log 追記失敗 | `500`。token record は残る。 |
| token record 破損 | `500`、自動再生成なし、secret 出力なし。 |
| scope 不足 | `403`、endpoint 実行なし、`permission_denied` 記録。 |
| 採番衝突 | 既存最大 suffix + 1 で作成される。 |
| 作成 response | token 本体は 1 回だけ含まれ、一覧では返らない。 |

<a id="sec-27-44"></a>
**27.44 監査ログ：**
owner component は `security` とする。collaborator component は `api`、`runner`、`statefile` とする。


本機能の目的は、認証、権限拒否、token 作成・失効、設定変更、build trigger、session revoke、TOTP、approval の状態変更を追跡できる JSON Lines 監査ログとして保存することである。

**record 生成規則：**

| 項目 | 仕様 |
|------|------|
| `timestamp` | 操作結果が確定した UTC 時刻。 |
| `request_id` | API event は [`docs/details/api.md` 詳細本文責務 request ID 固定契約](api.md#api-request-id-contract) で生成した値を使用し、同一 API 処理中に複数 log を書く場合は同じ値を使う。`runner` 内部 event は `null`。 |
| `actor_type` | 未認証 login は `"anonymous"`、管理 session は `"admin"`、API token は `"api_token"`、署名検証済み Webhook は `"webhook"`、内部処理は `"system"`。 |
| `actor_id` | 管理sessionはlegacyで `admin`、multi-user-v1では認証済UserRecord.id。API tokenはtoken id、署名検証済みWebhookは `webhook`、未認証はnull、内部処理は `system`。 |
| `target_id` | 対象 id がある場合は id。endpoint 拒否は `"{METHOD} {path}"`。対象なしは `null`。 |
| `result` | 成功は `"success"`、認証失敗や検証失敗は `"failure"`、権限拒否は `"denied"`。 |
| `remote_addr` | [`docs/details/api.md` 詳細本文責務 §27.6](api.md#sec-27-6) の `remote_addr` 導出規則に従う。API request に紐づかない内部処理は `null`。 |
| `message` | 固定文言のみ。入力値を連結しない。最大 500 文字。 |

監査 record のキー、型、必須性、許容値は [`docs/details/statefile.md` 詳細本文責務 `.audit_log` schema](statefile.md#audit-log-schema) を唯一の正本とする。監査ログへ保存する object はその schema のキーだけとし、未知キー、request body、query 全体、header 全体、cookie、secret、token、password、hash、salt、TOTP secret を保存してはならない。

multi-user-v1のuser/role操作だけは [user/role確定契約](#multi-user-record-lifecycle) のaction/messageを使用し、以下のlegacy操作表へ重複列挙しない。

**action 生成固定値：**

| 操作 | action | actor_type | target_type | target_id | result |
|------|--------|------------|-------------|-----------|--------|
| password login 成功 | `login_success` | `anonymous` | `auth` | `login` | `success` |
| password login 失敗 | `login_failure` | `anonymous` | `auth` | `login` | `failure` |
| login ロック拒否 | `permission_denied` | `anonymous` | `auth` | `login` | `denied` |
| TOTP 必須 ticket 発行 | `totp_required` | `anonymous` | `auth` | `login` | `success` |
| TOTP login 失敗 | `totp_failure` | `anonymous` | `auth` | `login_totp` | `failure` |
| logout | `logout` | `admin` または `api_token` | `auth` | `logout` | `success` |
| 他 session 一括失効 | `session_revoke_all` | `admin` または `api_token` | `auth` | `sessions` | `success` |
| password 変更 | `password_change` | `admin` | `auth` | `password` | `success` |
| API token 作成 | `token_create` | `admin` または `api_token` | `api_token` | 作成 token id | `success` |
| API token 認証成功 | `token_auth` | `api_token` | `api_token` | token id | `success` |
| API token 期限切れ拒否 | `token_expired` | `api_token` | `api_token` | token id | `failure` |
| API token 失効済み拒否 | `token_revoked_reject` | `api_token` | `api_token` | token id | `failure` |
| API token 失効 | `token_revoke` | `admin` または `api_token` | `api_token` | 失効 token id | `success` |
| scope 不足 | `permission_denied` | `api_token` | `endpoint` | `{METHOD} {path}` | `denied` |
| rate limit 超過 | `permission_denied` | `anonymous`、`admin`、`api_token`、`webhook` | `endpoint` | `{METHOD} {path}` | `denied` |
| TOTP setup 開始 | `totp_setup` | `admin` または `api_token` | `auth` | `totp_setup` | `success` |
| TOTP confirm / disable code 失敗 | `totp_failure` | `admin` または `api_token` | `auth` | `totp_confirm` または `totp_disable` | `failure` |
| TOTP 有効化 | `totp_enabled` | `admin` または `api_token` | `auth` | `totp` | `success` |
| TOTP 無効化 | `totp_disabled` | `admin` または `api_token` | `auth` | `totp` | `success` |
| 通常 build 開始または queue 追加 | `build_trigger` | `admin` または `api_token` | `build` | build id または queue id | `success` |
| force build 開始または queue 追加 | `build_force_trigger` | `admin` または `api_token` | `build` | build id または queue id | `success` |
| `.config_log` 対象設定変更成功 | `config_update` | `admin` または `api_token` | `config` | `.config_log.type` と同じ固定名 | `success` |
| `.config_log` 対象設定変更部分失敗 | `config_update` | `admin` または `api_token` | `config` | `.config_log.type` と同じ固定名 | `failure` |
| rate limit 設定変更成功 | `rate_limit_update` | `admin` または `api_token` | `config` | `api_rate_limit` | `success` |
| rate limit 設定変更部分失敗 | `rate_limit_update` | `admin` または `api_token` | `config` | `api_rate_limit` | `failure` |
| approval pending 作成 | `approval_pending` | `system` | `approval` | approval id | `success` |
| approval 承認 | `approval_approved` | `admin` または `api_token` | `approval` | approval id | `success` |
| approval 却下 | `approval_rejected` | `admin` または `api_token` | `approval` | approval id | `success` |
| approval 期限切れ | `approval_expired` | `system` | `approval` | approval id | `success` |
| share link 作成 | `share_link_create` | `admin` または `api_token` | `share_link` | 作成 share link id | `success` |
| share link 失効 | `share_link_revoke` | `admin` または `api_token` | `share_link` | 失効 share link id | `success` |

<a id="sec-27-44-config-audit-order"></a>
**設定変更の保存順・失敗固定契約：**

同一操作で `.config_log` と `.audit_log` の両方を追記する場合、対象状態保存 → `.config_log` → `.audit_log` → response の順に固定する。`.config_log` 追記失敗時は server log に `CONFIG_LOG_WRITE_FAILED` だけを記録し、`.audit_log` を試行せず `500` を返す。`.audit_log` 追記失敗時は server log に `AUDIT_LOG_WRITE_FAILED` だけを記録し、`500` を返す。いずれの失敗でも保存済み対象状態と追記済み先行 log を巻き戻さず、入力値、diff、Go error、path、secret を server log に含めない。`.audit_log` 追記失敗そのものを `.audit_log` に再帰記録してはならない。

**監査 record 保存順・失敗契約：**

| 操作種別 | 保存順 | `.audit_log` 失敗時 |
|----------|--------|----------------------|
| 認証成功 | session または token 状態更新 → `.access_log` → `.audit_log` → response | token / session 状態は巻き戻さず `500`。response に token 本体を含めない。 |
| 認証失敗 | `.access_log` → `.audit_log` → response | `500`。失敗理由詳細は返さない。 |
| 権限拒否 | `.access_log` → `.audit_log` → `403` | `500`。対象 endpoint は実行しない。 |
| 設定変更 | [設定変更の保存順・失敗固定契約](#sec-27-44-config-audit-order) を適用する。 | 同契約を適用する。 |
| token 作成 | `.api_tokens` 保存 → `.access_log` → `.audit_log` → response | 作成済み record は残し、token 本体は返さず `500`。 |
| build trigger | build / queue 状態保存 → `.audit_log` → response | 保存済み build / queue 状態は巻き戻さず `500`。 |
| session revoke / TOTP | 対象状態更新 → `.access_log` が必要な場合は追記 → `.audit_log` → response | 保存済み状態は巻き戻さず `500`。one-time secret または token は返さない。 |
| approval API | queue / approval / history を [`docs/details/api.md` 詳細本文責務 §27.30](api.md#sec-27-30) の順で更新 → `.audit_log` → response | 保存済み状態は巻き戻さず `500`。 |
| approval runner event | approval / history を [`docs/details/runner.md` 詳細本文責務 §27.30](runner.md#sec-27-30) の順で更新 → `.audit_log` | 保存済み状態は巻き戻さず runner failure とし、後続通知は実行しない。 |
| share link 作成 | `.share_links` 保存 → `.audit_log` → `.admin_events` → token response | 保存済み share link record は巻き戻さず、token 本体は返さず `500`。 |
| share link 失効 | `.share_links` 保存 → `.audit_log` → `.admin_events` → response | 保存済み失効状態は巻き戻さず `500`。 |

**取得仕様：**

| 項目 | 仕様 |
|------|------|
| 読み込み順 | ファイル先頭から全有効行を読み、フィルタ後に `timestamp` 降順、同時刻はファイル出現順の逆順で返す。 |
| `total` | フィルタ後、ページング前の有効 record 件数。 |
| `actor` filter | `actor_id` と完全一致。`anonymous` を指定した場合は `actor_type:"anonymous"` かつ `actor_id:null` に一致させる。 |
| `action` filter | `action` 完全一致。 |
| `result` filter | `success`、`failure`、`denied` の完全一致。 |
| 壊れた行 | 無視する。API response に壊れた行の内容を含めない。 |
| 返却上限 | `limit` は [`docs/details/api.md` 詳細本文責務 §22.0b](api.md#sec-22-0b) の 1〜200。未指定時は 100。 |

`GET /api/audit-log` は監査ログ取得操作自体を `.audit_log` へ記録しない。通常の API request として `.api_access_log` には記録する。フィルタ値が未知 action、未知 result、200 Unicode scalar values 超過の actor の場合は `422` とし、壊れた行の有無とは独立して判定する。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| password 変更 | `password_change` が記録される。 |
| 権限拒否 | `permission_denied` が記録される。 |
| secret 入力 | 監査ログに平文がない。 |
| 壊れた行 | API は無視して返す。 |
| filter | `actor`、`action`、`result` が完全一致で絞り込まれる。 |
| 追記失敗 | 対象操作は `500`。 |
| 取得操作 | `GET /api/audit-log` 自身は監査ログへ追記されない。 |
| 未知 filter | `422`。 |
| audit failure token create | token record は残るが token 本体は返らない。 |
| body secret | request body 全体が保存されない。 |

<a id="sec-27-45"></a>
**27.45 セッションタイムアウト変更設定：**
[`docs/details/security.md` 詳細本文責務 §27.45](security.md#sec-27-45) の境界は owner component `security`、collaborator component `api`、`sdk`、`ui`、`statefile` とする。


本機能の目的は、新規 session の有効期限を管理 API から更新し、既存 session への影響を明確にすることである。

**設定参照：**

| 項目 | 値 |
|------|----|
| 設定 schema | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.server_config.session_timeout_seconds`。型、既定値、許容範囲は同 schema を正とする。 |
| 更新 API | [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の `POST /api/config`。SDK 呼出は [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) の `setConfig(config)` を正とする。 |

設定変更は新規 session にだけ適用する。既存 session の `expires_at` は延長も短縮もしない。

**処理契約：**

| 項目 | 仕様 |
|------|------|
| 保存先 | `.server_config.session_timeout_seconds`。他 key と同時更新された場合も [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の単一ファイル更新手順で保存する。 |
| 既定値 merge | `.server_config` に key がない場合、`GET /api/config` は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の schema 既定値を返す。ファイルへ暗黙保存しない。 |
| login 時適用 | session token 発行直前に `.server_config` を読み、当該時点の値で `expires_at` を計算する。 |
| TOTP login | TOTP 有効時は `POST /api/login/totp` の成功時点で値を読む。`POST /api/login` の password 成功時点では session を発行しない。 |
| 変更監査 | `POST /api/config` で値が変わった場合は `.config_log` に差分を記録する。`.audit_log` は `config_update`、`target_type:"config"`、`target_id:"session_timeout_seconds"` を記録する。 |
| 同値更新 | 同じ値の更新は `200` とし、`.server_config`、`.config_log`、`.audit_log`、既存 session を変更しない。 |

session timeout の値は session 発行時に秒単位で加算する。`expires_at = issued_at + session_timeout_seconds` とし、計算後の時刻は UTC ISO 8601 秒精度で保存する。ミリ秒、ナノ秒、local timezone は保存しない。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| schema 最小値の設定 | 新規 session が設定値の秒数後に期限切れ。 |
| 既存 session | 設定変更後も元の `expires_at`。 |
| 範囲外 | `422`。 |
| key 不在 | `GET /api/config` は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の schema 既定値を返す。 |
| TOTP login | `POST /api/login/totp` 成功時点の値で session 期限を決める。 |
| 秒精度 | `expires_at` は UTC ISO 8601 秒精度。 |

<a id="sec-27-46"></a>
**27.46 TOTP 二要素認証：**
[`docs/details/security.md` 詳細本文責務 §27.46](security.md#sec-27-46) の境界は owner component `security`、collaborator component `api`、`sdk`、`ui`、`statefile` とする。


本機能の目的は、外部ライブラリなしで RFC 6238 TOTP を検証し、password 漏えい時の管理 API 不正利用を抑止することである。

**アルゴリズム：**

| 項目 | 仕様 |
|------|------|
| secret | `crypto/rand` 20 bytes を RFC 4648 base32 padding なしで表示する。 |
| HMAC | `crypto/hmac` + `crypto/sha1`。 |
| step | 30 秒。 |
| digits | 6 桁。 |
| window | 現在 step の前後 1 step を許可する。 |
| replay 防止 | `.totp_secret.last_accepted_step` 以下の step は拒否する。 |
| otpauth URI | `otpauth://totp/Adlaire%20CI:admin?secret={secret}&issuer=Adlaire%20CI&algorithm=SHA1&digits=6&period=30`。 |

QR code 生成はこの詳細本文責務では仕様化しない。UI は secret と otpauth URI を一回表示し、ユーザーの認証アプリ登録手段は手入力または URI 貼り付けに限定する。

**メモリ上状態：**

| 状態 | 保存場所 | 期限 | 内容 |
|------|----------|------|------|
| setup 仮 secret | `api` のメモリ | 10 分 | `secret_base32`, `created_at`。サーバー再起動で破棄する。 |
| login ticket | `api` のメモリ | 5 分 | `ticket_hash`, `created_at`, `password_verified_at`, `credentials_fingerprint`。ticket 本体は hash 化して保持する。 |

setup 仮 secret と login ticket は永続ファイルへ保存しない。API response、UI 一回表示、メモリ上状態以外に secret/ticket 本体を残してはならない。

**正常系：**

1. `POST /api/auth/totp-setup` は仮 secret と otpauth URI を返すが、永続化しない。
2. `POST /api/auth/totp-confirm` は仮 secret と code を検証し、成功時だけ `.totp_secret.enabled=true` と secret 情報を保存する。
3. TOTP 有効時の `POST /api/login` は password 成功後に ticket を返し、session token は返さない。
4. `POST /api/login/totp` は ticket と code を検証し、成功時に session token を返す。
5. `DELETE /api/auth/totp` は code を検証して TOTP を無効化する。

**endpoint 処理詳細：**

| endpoint | 処理 |
|----------|------|
| `GET /api/auth/totp-status` | `.totp_secret` を読み、`{enabled,confirmed_at}` だけを返す。`secret_base32` と `last_accepted_step` は返さない。 |
| `POST /api/auth/totp-setup` | TOTP 有効時は `409`。無効時は仮 secret を生成し、既存の未確認仮 secret を上書きする。`.totp_secret` は書き込まない。 |
| `POST /api/auth/totp-confirm` | 仮 secret がない、または期限切れなら `409`。code 成功時に `.totp_secret` を保存し、仮 secret をメモリから削除する。 |
| `DELETE /api/auth/totp` | `.totp_secret.enabled == false` は `409`。code 成功時に `enabled:false`, `secret_base32:null`, `confirmed_at:null`, `last_accepted_step:null` を保存する。 |
| `POST /api/login` | password 成功かつ TOTP 有効なら `ticket` を `crypto/rand` 32 bytes の lowercase hex で生成する。credentials を更新せず、ticket に `must_change` を保持しない。HTTP response は [`docs/details/api.md` 詳細本文責務 `POST /api/login`](api.md#25-認証-実装仕様) の `LoginResult` を参照する。 |
| `POST /api/login/totp` | ticket hash と code を検証し、成功時に ticket を削除する。最新 credentials を再読込して session 発行時の `must_change` を算出する。失敗時も ticket は削除する。HTTP response は [`docs/details/api.md` 詳細本文責務 `POST /api/login/totp`](api.md#25-認証-実装仕様) を参照する。 |

TOTP code は 6 桁の ASCII 数字のみ受け付ける。空文字、全角数字、空白付き文字列、6 桁以外は `422` とする。検証は `window` 内の step を古い順に試し、最初に一致した step を採用する。採用 step が `.totp_secret.last_accepted_step` 以下の場合は `401` とする。

TOTP setup、ticket 発行、login TOTP 成功 / 失敗、confirm 成功 / 失敗、disable 成功 / 失敗は `.audit_log` へ記録する。code 不一致、replay、ticket 不正は `result:"failure"` とし、code、secret、ticket 本体は保存しない。

**TOTP 計算固定契約：**

| 項目 | 仕様 |
|------|------|
| counter | `floor(unix_seconds / 30)` を 8 byte big-endian unsigned integer として HMAC 入力にする。 |
| truncation | RFC 4226 dynamic truncation を使い、31 bit integer を `10^6` で剰余する。 |
| 表示 | 6 桁未満は左ゼロ埋めする。 |
| base32 decode | 大文字 ASCII のみ保存する。入力確認時は空白を除去せず、保存値と同じ RFC 4648 padding なし形式だけを扱う。 |
| clock source | API server の現在時刻だけを使う。client 時刻は受け取らない。 |

**TOTP 状態更新順：**

| 操作 | 更新順 | 失敗時 |
|------|--------|--------|
| setup 開始 | 仮 secret 生成 → メモリ保存 → `.audit_log` 追記 → response | メモリ保存または audit 追記失敗時は仮 secret を削除して `500`、secret を返さない。 |
| confirm 成功 | `.totp_secret` 保存 → 仮 secret 削除 → `.audit_log` 追記 → response | `.audit_log` 失敗時は `500`。保存済み `.totp_secret` と仮 secret 削除は巻き戻さない。 |
| login password 成功 / TOTP 有効 | ticket 生成 → メモリ保存 → `.access_log` 追記 → `.audit_log` 追記 → response | log 失敗時は ticket を削除し `500`。 |
| login TOTP 成功 | code step 採用 → `.totp_secret.last_accepted_step` 保存 → ticket 削除 → `.admin_credentials` 更新 → `.access_log` 追記 → `.audit_log` 追記 → session 追加 → response | session 追加前の失敗は token を返さない。ticket 削除後は同 ticket を再利用不可。 |
| login TOTP 失敗 | ticket 削除 → `.access_log` 追記 → `.audit_log` 追記 → response | log 失敗時は `500`。ticket は巻き戻さない。 |
| disable 成功 | `.totp_secret` を無効値で保存 → 未使用 ticket / 仮 secret 全削除 → `.audit_log` 追記 → response | `.audit_log` 失敗時は `500`。保存済み無効化は巻き戻さない。 |

`POST /api/login/totp` は、ticket が存在しない、期限切れ、hash 不一致、既に削除済みのいずれの場合も `401 {"error":"Unauthorized"}` を返す。ticket 不正の詳細、ticket hash、ticket 発行時刻は response と log に含めない。TOTP 有効化または無効化に成功した場合、既存 session は破棄しない。

**異常系：**

| 条件 | 処理 |
|------|------|
| code 不一致 | `401`。 |
| code 形式不正 | `422`。 |
| ticket 期限切れ | `401`。ticket 有効期限は 5 分。 |
| setup 未実行 confirm | `409`。 |
| TOTP 無効状態の disable | `409`。 |
| TOTP 有効状態の setup | `409`。 |
| `.totp_secret` 破損 | TOTP 有効 login、status、confirm、disable は `500`。自動無効化しない。 |
| `.totp_secret` 書き込み失敗 | `500`。response に secret、token、ticket を含めない。 |
| `.audit_log` 追記失敗 | `500`。保存済み `.totp_secret` は巻き戻さない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 正常 setup | secret は確認まで保存されない。 |
| 正常 login | password 後 ticket、TOTP 後 token。 |
| replay | 同一 step の再利用は `401`。 |
| secret 表示 | response と UI の一回表示以外に平文が残らない。 |
| setup 期限切れ | confirm は `409`。 |
| login ticket 再利用 | 1 回成功後または失敗後の同一 ticket は `401`。 |
| confirm audit 失敗 | `500`、`.totp_secret` は保存済み、secret 平文は log なし。 |
| disable 成功 | `.totp_secret` は無効値、既存 session は維持、ticket と仮 secret は削除。 |
| window 前後 | 現在 step の前後 1 step が成功し、同 step 再利用は失敗する。 |
| 全角 code | `422`。 |

<a id="sec-27-47"></a>
**27.47 API レート制限：**
[`docs/details/security.md` 詳細本文責務 §27.47](security.md#sec-27-47) の境界は owner component `security`、collaborator component `api`、`sdk`、`ui`、`statefile` とする。


本機能の目的は、login 総当たり、API token 濫用、外部連携の暴走を Go 標準ライブラリだけで抑止することである。

**固定窓 policy：**

| group | 既定 window | 既定 max | 対象 |
|-------|-------------|----------|------|
| `login` | 60 秒 | 10 | [endpoint group判定](#security-rate-limit-group-map) のlogin集合。 |
| `read` | 60 秒 | 600 | 読み取り endpoint。 |
| `trigger` | 60 秒 | 60 | build trigger endpoint と署名検証済み `POST /api/webhook`。 |
| `operate` | 60 秒 | 120 | 運用操作 endpoint。 |
| `config` | 60 秒 | 60 | 設定変更 endpoint。 |
| `admin` | 60 秒 | 60 | token、audit、rate limit endpoint。 |

**判定キー：**

| 認証状態 | key |
|----------|-----|
| 認証前 | `ip:{remote_addr}:{group}` |
| session | `session:admin:{group}` と `ip:{remote_addr}:{group}` の両方 |
| API token | `token:{token_id}:{group}` と `ip:{remote_addr}:{group}` の両方 |
| 署名検証済み Webhook | `ip:{remote_addr}:trigger` だけ |

どちらか一方でも上限を超えた場合は `429 Too Many Requests` と `{"error":"Too many requests"}` を返す。

`remote_addr` は `net/http.Request.RemoteAddr` の host 部分を使用する。`X-Forwarded-For`、`X-Real-IP`、`Forwarded` header は標準では信用せず、key 生成に使用しない。IPv6 は `net.SplitHostPort` で host を抽出し、正規化済み文字列をそのまま key に入れる。host 抽出に失敗した場合は `ip:unknown:{group}` を使用する。

**policy / response schema 参照：**

永続化する `.server_config.api_rate_limit` のキー、型、必須 group、範囲、既定値は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の ApiRateLimitPolicy object を唯一の正本とする。`GET /api/api-rate-limit` と `POST /api/api-rate-limit` の request / response、response 専用 `state_summary` のキー、並び順、件数は [`docs/details/api.md` 詳細本文責務 security owner 機能 API 側境界](api.md#sec-security-owner-api-boundary) を正とする。`state_summary` は永続化しない。

**正常系：**

1. path 解決後、認証前に IP key の login 制限を確認する。
2. 認証後、endpoint group を判定し、actor key と IP key の count を更新する。署名検証済み Webhook は `trigger` group の IP key だけを更新する。
3. `GET /api/api-rate-limit` は policy と window summary を返す。
4. `POST /api/api-rate-limit` は policy を検証して保存し、`.api_rate_state.windows` を空にする。

<a id="security-rate-limit-group-map"></a>
**endpoint group 判定：**

| 条件 | group |
|------|-------|
| `POST /api/login`、`POST /api/login/totp`、multi-user-v1のOIDC start/finish | `login` |
| 署名検証済み `POST /api/webhook` | `trigger` |
| [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の `trigger` 許可 endpoint | `trigger` |
| [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の `operate` 許可 endpoint | `operate` |
| [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の `config` 許可 endpoint | `config` |
| [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の `admin` 許可 endpoint | `admin` |
| [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の `read` 許可 endpoint | `read` |
| `GET /api/health` | rate limit 対象外 |
| 未知 path / method 不一致 | rate limit 判定前に `404` / `405` |

rate limit group は [endpoint group 判定表](#security-rate-limit-group-map) と [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の一意な endpoint 割当てから決定する。同じ method と正規化 path pattern が複数 group に存在する場合は仕様不整合であり、実装者判断の優先順で解決してはならない。

**判定・更新手順：**

1. `.server_config.api_rate_limit.enabled == false` の場合、`.api_rate_state` を読まずに対象 API 処理へ進む。
2. endpoint group を決める。
3. `.api_rate_state` のファイルロックを取得する。
4. 対象 key ごとに `.api_rate_state.windows[key]` を確認する。
5. window が存在しない、または `now >= window_start + window_seconds` の場合、`window_start=now`, `count=0` で初期化する。
6. `count >= max_requests` の key が 1 つでもあれば、count を増やさず `.audit_log` に `permission_denied` を追記し、`429` を返す。
7. 上限未満の場合、対象 key すべての `count` を 1 増やして `.api_rate_state` を保存し、対象 API 処理へ進む。

rate limit の `429` は `.audit_log` に `permission_denied` として記録する。監査ログ追記に失敗した場合は `500` を返す。login group の認証前 `429` は `actor_type:"anonymous"`、`actor_id:null`、署名検証済み Webhook の `429` は `actor_type:"webhook"`、`actor_id:"webhook"` とする。

認証後 endpoint の rate limit では、actor key と IP key の両方を同じ lock 内で判定・更新する。片方だけの count 更新に成功した状態を残してはならない。`.api_rate_state` 保存失敗時は対象 API を実行せず `500` を返す。rate limit 判定で `429` になる request は count を増やさない。

**rate limit 副作用固定契約：**

| ケース | `.api_rate_state` | `.api_access_log` | `.audit_log` | endpoint 固有処理 |
|--------|-------------------|---------------|--------------|-------------------|
| 上限未満 | count を増やす | response 確定後に通常追記 | endpoint が監査対象の場合だけ追記 | 実行する |
| 上限超過 | count を増やさない | `429` として追記 | `permission_denied` を追記 | 実行しない |
| `.api_rate_state` 保存失敗 | 部分更新を残さない | `500` として追記を試行 | 追記しない | 実行しない |
| `.audit_log` 失敗 | count を増やさない | `500` として追記を試行 | 失敗 | 実行しない |

`429` 判定時は `.audit_log` 追記を `.api_rate_state` 保存前に行う。`.audit_log` 追記に成功した場合だけ `429` を返す。`.audit_log` 追記に失敗した場合は `.api_rate_state` を変更せず `500` を返す。

**設定更新契約：**

| 項目 | 仕様 |
|------|------|
| policy validation | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の ApiRateLimitPolicy object schema に完全一致する場合だけ保存する。 |
| 保存順 | `.server_config` 保存 → `.api_rate_state.windows` 空保存 → `.config_log` 追記 → `.audit_log` に `rate_limit_update` 追記 → response。 |
| 同値更新 | `200` とし、`.server_config` と `.api_rate_state` は変更しない。`.config_log` と `.audit_log` に追記しない。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| policy schema 不一致 | `422`。`.server_config`、`.api_rate_state`、`.config_log`、`.audit_log` を変更しない。 |
| `.api_rate_state` 書き込み失敗 | `500`。対象 API は実行しない。 |
| `.api_rate_state` 破損 | [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) に従って退避し、空 window で再生成する。 |
| `.audit_log` 追記失敗 | `500`。`429` response は返さない。 |
| lock 取得 10 秒超過 | `409`。対象 API は実行しない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| login 11 回目 | `429`。 |
| window 経過 | count reset、成功。 |
| session と IP | どちらか超過で `429`。 |
| disabled | `enabled:false` の場合は判定せず成功。 |
| policy 更新 | `.api_rate_state.windows` が空になる。 |
| state_summary | 最大 100 件、`reset_at` 降順で返る。 |
| 429 count | `429` になった request では count が増えない。 |
| 同値更新 | state と log を変更せず `200`。 |
| proxy header | `X-Forwarded-For` ではなく `RemoteAddr` host で key を作る。 |
| audit failure on 429 | count は増えず `500`。 |
| state save failure | endpoint 固有処理なし、部分 count 更新なし。 |

<a id="notify-hmac-signature-contract"></a>
**通知HMACの計算：** Webhook config.secretのUTF-8 byteをkey、実際に送信するHTTP body byteをmessageとし、Go標準crypto/hmac・sha256でHMAC-SHA256を計算する。header値はASCII sha256=とlowercase hex64桁の連結だけ、secret=nullは署名header無し。JSONを再encodeした別byte、payload_sha256、受信用.webhook_secretをmessage/keyへ代用しない。header名は呼出元のAPI test/runner wire契約が選択し、同時に両headerを送信しない。signature/raw secretはmemoryだけ、response/log/event/expectedへの実値転写を禁止する。

<a id="sec-27-48"></a>
**27.48 マルチユーザー対応：**

マルチユーザー対応の owner は `security` とする。

`.users` と `.roles` の schema は [`docs/details/statefile.md` 詳細本文責務 §22.0d](statefile.md#sec-22-0d) を正本とする。

既定 user は初期セットアップ時に 1 件だけ作成する。既定 user は `admin` role を持つ。

session には `user_id`、`role_ids`、`permissions_hash`、`issued_at`、`expires_at` を保存する。

role または permission が変更された場合、既存 session の `permissions_hash` が不一致になった request は `401` とし、再 login を要求する。

最後の active admin user を disabled、locked、role 剥奪してはならない。

さらにactive adminかつlocal credential非nullのuserを最低1件維持し、外部provider障害で全管理者が復旧不能になる変更を409とする。user createでpassword省略かつexternal_subjects空は422。初期移行userはdisplay_name=username、external_subjects=[]、role_ids=["admin"]、created_at/updated_at=journal.started_at、status=activeに固定する。

user/roleの更新と最後のadmin判定は同一authorization更新lock内で最新users/rolesを読み、更新後全体を検証してから保存する。最後の2管理者への並行無効化でも更新後の最低1件を維持する。lock取得時にactorの最新session/user/permissionも再検証し、actor失効なら401、存続actorによる最後のadmin変更なら409とする。user/role参照をread途中で混在させない。管理APIによるpassword resetまたは資格情報削除は主state確定直後にそのuserの全memory sessionを失効させ、audit失敗でも再有効化しない。本人POST /api/change-passwordだけは既存の現在session維持契約を使い、markerによるtoken失効完了と必須log成功後に現sessionのpassword_change_required=falseを確定する。その他のsession/ticketは主state確定直後に失効、途中失敗では現sessionも失効させる。旧credential fileは変更しない。permissions_hashは閉じたobject `{permissions:string[],role_ids:string[]}`（両配列重複なしASCII昇順）のcanonical JSONに対するSHA-256とし、保存済みsession値と毎requestで比較する。

**監査対象：**

| 操作 | audit action | audit/event message |
|------|--------------|---------------------|
| user create | `user_create` | `User created` |
| user update | `user_update` | `User updated` |
| user disable | `user_disable` | `User disabled` |
| role assignment change | `user_role_update` | `User roles updated` |
| role create | `role_create` | `Role created` |
| role update | `role_update` | `Role updated` |
| role delete | `role_delete` | `Role deleted` |

<a id="multi-user-record-lifecycle"></a>
**user/role管理の確定・no-op：** UserCreateの外部subject省略は[]、TOTPは既存のdisabled完全object、credentialはpassword省略ならnull。local passwordのcreate/admin resetは毎回新saltで既存KDFを使い、must_change=true/login_count=0/last_login_at=null/updated_at=現在時刻とする。既存passwordと同じ平文でもadmin resetはno-opにしない。通常の本人password変更は既存認証契約を維持する。

解決済permissionに*を持たないactorは、既に*を持つuserの変更/disable/password resetと、新たに*を与えるuser create/role assignmentを403で拒否する。user:adminは非superadminのuser/role割当を管理する明示権限であり、system adminの自己作成や初期adminの奪取を許可しない。

PATCHは最新recordへmerge後に全体検証し、上のactor保護403→admin保護409を先に判定、その後active userがcredential=nullかつexternal_subjects=[]なら422とする。user/role参照不在は422、別userが同じprovider/subjectを持つ場合は409。同じusernameはdisabled/lockedを含め409。createのusers1000件/roles128件到達は409、既存recordのno-op/PATCHを容量だけで拒否しない。

passwordを含まないPATCHの正規化後recordが同値ならupdated_atも変更せず、主state/audit/event/cache世代更新0。DELETE済disabled userへの再DELETEも200/no-op。差分確定時だけupdated_atを更新する。user PATCHのaudit actionはstatusをdisabledへ変える場合user_disable、そうでなくrole_ids差分ならuser_role_update、その他user_updateの優先順で1件だけ選ぶ。lock状態への変更はuser_update。roleのno-opは同じ扱いで、system role制約をno-op判定より先に検証する。

上表のauditはactor_type=admin|api_token、actor_idは認証済id、target_type=user|role、target_idは対象id、result=success、timestamp=確定時刻、request_id/remote_addrは共通規則。admin eventはtype=user、actor=同id、target_type/idとmessageは同表、severity=info、idは共通時刻採番とする。主state→audit→eventの途中失敗では後続を書かず500、成功済み主stateを戻さない。user失効処理の永続marker順は [user token失効状態責務](statefile.md#multi-user-revocation-state) を先に適用し、失効中の認証を許可しない。

<a id="multi-user-migration-contract"></a>
**単一管理者からの明示移行：** owner=`security`、分類=`STD`（既存KDF/TOTP primitiveは各正本の分類を継承）。legacy-v1とcicd-v1の認証を混在させず、追加管理機能の認証profileは `multi-user-v1` とする。profile状態とlockは [statefile詳細本文責務](statefile.md#multi-user-profile-state)。通常起動で自動移行しない。

移行CLIは `adlaire-ci-api --migrate-users --state-dir ABS --username NAME`。他operationとの併用・未知/重複optionはexit2、stdoutは成功時だけ `{"profile":"multi-user-v1","migrated":true}` + LF、stderrは `AUTH_MIGRATION_INVALID` / `AUTH_MIGRATION_BUSY` / `AUTH_MIGRATION_FAILED` の固定code + LF、exitは順に2/3/4。secretの入力・再表示を行わない。argument/schema/profile/事前条件不適合はexit2、lifecycle lock競合は3、I/O/entropy/journal不一致/中断再開の整合不適合は4とする。成功以外のstdoutは空、stderrは分類code1行だけ。API停止をlifecycle flockで確認してから、旧credential/TOTP/tokenの完全schema/権限、未作成のusers/roles/profile、usernameを検証する。旧TOTP不在は無効値、旧tokens不在は空配列。旧credential不在・旧algorithm・破損・非対応profileは拒否する。

乱数User idと開始時刻を確定しprepared journalをcreate-only保存→system roles作成→初期active admin user作成（credential/TOTPは元の全keyを変更せずcopy）→全旧tokenをそのuserへ紐付け、revoked_atがnullのものを開始時刻で失効→全出力をread-back検証→auth_profile atomic作成→journalをcommittedへ更新→journal削除/parent syncの順。旧credential/TOTP fileは消さず0600で保持するが、新profileでは認証に使用しない。移行後のlegacyへのdowngradeは禁止。初期管理者を新password無しで作る別経路は設けない。

preparedからの再実行は保存済みid/時刻/source digestだけを使用する。旧credential/TOTPは元file byte digestとの一致を必須とし、tokensだけは元byte digestまたは予定出力のcanonical digest一致を許可する（移行自身の確定済みwriteを外部改変と誤認しない）。その他の差分は停止する。既存出力はjournal.output_sha256とcanonical digestが一致するものだけ再利用し、差分は上書きしない。tokens書換済みの再開は全recordのuser_id/revoked_atが予定値に一致する場合だけ許可する。profile未commitでjournalがあればHTTP起動を拒否する。profile commit後のcleanup失敗は新profileを維持し、同CLIでcleanupだけを再開する。journal無し・profile存在の再移行はexit2。移行前APIを再開するためにjournalやprofileを手作業削除する復旧を許可しない。

<a id="multi-user-token-totp-contract"></a>
**user資格情報・token/TOTP結合：** loginはmulti-user-v1だけusername必須（正規化はユーザー管理契約）、passwordは該当user credentialを照合する。不在/locked/disabled/credential=nullは同じ401、username存在性を返さずrate limitを適用する。must_change、login_count、last_login_atは当該credentialだけを更新する。session/user id、TOTP ticket/setup仮secretは同じuserへ束縛し、他userへの流用・role変更後のticket消費を拒否する。

TOTPのalgorithm/window/replay/confirmは [既存TOTP詳細本文責務](#sec-27-46) を再利用し、保存先だけUserRecord.totpへ変更する。setup/confirm/deleteは本人sessionと現passwordを必要としAPI token/OIDC-only userは実行不可。無効化は当該user全session/ticket失効、last_accepted_step更新はusers lock内のcompare-and-updateとする。totp_enabledだけの保存・全user共有secret・UserPatch経由のTOTP解除は禁止。OIDC loginもlocal TOTP有効userには同じTOTP ticketを発行し、完了前にsessionを発行しない。

token発行user_idは本人sessionから確定しrequestでは指定させない。list/revokeは本人のactive sessionだけが本人分を扱い、解決済permissionに `*` を持つsessionだけが全user分を扱える。API tokenからGET/POST/DELETEのtoken管理routeはmulti-user-v1で403とし、legacyのadmin-scope発行連鎖を引き継がない。tokenのownerをrequestで指定する経路は作らない。*以外のsessionによる他user token idへのDELETEは不在と同じ404とし、所有者情報を返さない。target userのrevocation markerが未処理のPATCH/DELETEは409とし、restartによる既存失効復旧が完了する前に再有効化/別変更で上書きしない。毎requestでtoken hash/期限/revoke→owner user active→現在permission→scopeの順に検証する。追加管理routeのscope対応は [横断順序適用表](#additional-management-security-order) を使用し、role permissionとの積集合だけを許可する。既存scopeを新permissionへ自動昇格しない。発行時も現在permissionに対応しないscopeを403で拒否し、readはstatus:read/config:read、trigger/operateはbuild:write、configはconfig:write、adminはuser:admin/system:adminのいずれか（*は全件）を必要とする。全体100record容量到達は409、既存recordを黙って削除しない。user無効化/lock/password変更/credential削除時に当該userの有効tokenを全失効する。users commit後tokens write失敗はAPI500・当該userのtoken認証をprocess内で停止し、次起動も未完了revocationを検出するまで受付しない。復旧の永続revocation markerは [statefile詳細本文責務](statefile.md#multi-user-revocation-state) を参照する。

multi-user-v1の既存routeは、追加管理permission表にない場合system adminのみに制限する。ただし本人tokenのGET/POST `/api/tokens` とDELETE `/api/tokens/{id}` は本人active session（他user対象はsystem admin）で許可する。既存read routeは [明示許可表](#multi-user-legacy-read-authorization)、POST `/api/build`、`/api/build/force` はbuild:writeで許可する（tokenは既存trigger scopeも必要）。login/logout/change-password/TOTPは本人認証契約、署名Webhook/認証不要routeは従来契約を適用する。この閉じた例外以外にprefixやmethodから権限を推測しない。


<a id="multi-user-legacy-read-authorization"></a>
**既存read routeの明示許可：** multi-user-v1で下表だけを既存routeのsystem admin fallbackから除外する。route登録そのものを新設せず、APIの既存schema/状態read/secret maskを維持する。scope/rateは必要permissionに対応する既存read/adminを使用し、sessionでもpermission不足は403。表にないrouteは従来の閉じたfallbackを維持する。

| exact route template集合 | 必要permission |
|---|---|
| GET /api/status、GET /api/history、GET /api/dashboard、GET /api/queue、GET /api/maintenance、GET /api/dashboard-layout、GET /api/snapshots | status:read |
| GET /api/config、GET /api/config-log | config:read |
| GET /api/notify-log | system:admin |
| GET /api/audit-log | user:adminまたはsystem:adminのいずれか。*は共通規則で許可。masked audit recordだけを返す。 |

<a id="sec-27-50"></a>
**27.50 外部認証連携：**

外部認証連携の owner は `security` とする。

対応 provider type は `oidc` だけとする。

OIDC discovery URL は `issuer + "/.well-known/openid-configuration"` とする。

外部認証連携は local password 認証を暗黙に無効化しない。

provider secret は `.external_auth_config` に保存せず、`client_secret_ref` だけを保存する。

callbackとclaim判定の本文は [OIDC開始・完了と鍵検証](#oidc-code-flow-contract) を参照する。

検証失敗は `401` とし、失敗理由の詳細、token、claim 全文を response、server log、audit log に出してはならない。

<a id="oidc-code-flow-contract"></a>
**OIDC開始・完了と鍵検証：** `INHOUSE`、owner=`security`。規範は [OIDC Core code flow](https://openid.net/specs/openid-connect-core-1_0.html#CodeFlowAuth)、[RFC7636](https://www.rfc-editor.org/rfc/rfc7636)、[RFC7517](https://www.rfc-editor.org/rfc/rfc7517)。Go標準crypto/rsa・sha256・rand・subtle・encoding/json/base64・net/http/urlだけを使用し外部実装を移植しない。code flow + PKCE S256 + RS256だけを採用し、implicit/hybrid/none/HS*/EC/暗号化JWT/dynamic registration/UserInfo/emailによる自動user作成は不採用。

provider idは `[a-z][a-z0-9_-]{0,63}`、最大16provider、issuerはquery/fragment/userinfo/末尾slashなしのHTTPS URL。scopesは `["openid"]` 固定。redirect_uriは管理UIの検証済みHTTPS origin + `/` で、開始requestのOriginとexact一致を必要とする。client_secret_ref=nullならpublic client、非nullならclient_secret_basicだけを採用する。client secretは既存secret file adapterで0600保存、rawをconfig/responseへ入れない。provider削除/issuer変更はenabled値を問わず、紐付くuserがいれば409。provider更新は未完了flowを全失効する。

discovery/JWKS/token通信は各10秒・response最大262144byte・redirect禁止、TLS検証必須、環境proxy禁止。全接続でDNS解決結果すべてを公開global unicastに限定し、検証済みIPへdial（TLS ServerNameは元host）、private/loopback/link-local/multicastを拒否する。discovery issuerは設定とbyte一致、authorization/token/jwksは同じissuer originのHTTPS URLに限定する。metadataのcode/S256/RS256対応を確認し非対応なら拒否。metadata/JWKSはflow開始時取得しmemoryへ5分保持、完了時は同じsnapshotだけを使用、未知kidで外部再取得・任意jku/x5u追従をしない。

開始ごとにstate/nonce/browser binding/code_verifierを別々のCSPRNG32byteからpadding無しbase64url生成する。verifier43文字、challenge=base64url(SHA256(verifier ASCII))。最大128flow、TTL300秒、単一process memoryだけ、再起動で全失効。browser bindingはHttpOnly/Secure/SameSite=Lax/Path=/の `__Host-adlaire-oidc` cookie、Max-Age=300。開始は既存flow cookieを置換して当該flowを失効し、満杯は409、entropy失敗は500。authorization URLへresponse_type=code/client_id/redirect_uri/scope/state/nonce/code_challenge/code_challenge_method=S256だけをpercent encodeして付与する。

完了はstate・cookieをconstant-time照合し、有効flowをmutex内で一回だけ消費してからcodeを交換する。cookie不一致ではflowを消費しない。codeは1..2048byte、NUL/CR/LF禁止。token formはgrant_type=authorization_code/code/redirect_uri/code_verifier/client_id、confidentialだけBasic headerを加える。[RFC6749 §2.3.1](https://www.rfc-editor.org/rfc/rfc6749#section-2.3.1) に従い、client_idとsecretをそれぞれGo標準url.QueryEscapeでform符号化してからASCII colonで連結し、標準padding付きbase64を使用する。raw値をそのままBasicへ連結せず、client_secretをform/query/URLへ追加しない。200のJSONからid_tokenを取得し、access/refresh tokenは保存/使用しない。JWTは16384byte以下、3segmentのstrict unpadded base64url、JSON duplicate key/invalid UTF-8/未知critを拒否。header alg=RS256、kid必須、jku/x5u/x5c/crit禁止。JWKS最大32key、kid一意、RSA n=2048..4096bit・e=65537、useがあればsig、key_opsがあればverifyだけ、algがあればRS256。秘密key parameterを拒否し、一致するkey1件だけでrsa.VerifyPKCS1v15(SHA256)を行う。

署名成功後iss exact、audはclient_idと同じstring、またはclient_idを含む1件以上の空でないstringの重複なし配列とする。複数audではazpを必須とし、単独audでもazpを指定した場合はstringかつclient_idとのexact一致を必要とする。nonce exact、subは保存型、exp/iatと任意nbfはsigned 64bitの整数Unix秒（小数/指数表記/null禁止）を確認する。now < exp+30、iat <= now+30、iat >= flow.started_at-30、exp > iat、exp-iat <=3600、nbfがあればnow+30 >= nbf。加減算はchecked arithmeticまたはoverflowしない比較順を用い、signed overflowを境界合格へ読み替えない。未知claimは無視するがclaimをuser属性へcopyしない。全検証後のuser/role/TOTP再検証、発行順、provider更新との競合は [OIDC認証確定契約](#oidc-authentication-commit) を適用する。流出や不具合時はprovider無効化・未完了flow/session失効を行い、当該provider参照userが0になった場合だけ削除できる。

<a id="oidc-origin-contract"></a>
**OIDC URL/originの固定比較：** issuer/redirect_uri/metadata endpointはURL parserで検証し、schemeをASCII lowercaseでhttpsと判定、userinfo/fragment/queryを禁止する。DNS hostはASCIIのlabel1..63byte・全体253byte以下、各labelは英数字始終/中間英数hyphen、末尾dot/空labelを拒否する。IDNA/Unicode hostnameはv1で不採用。hostのpercent escape/IPv6 zone/IPv4-mapped IPv6は拒否する。数字とdotだけのhostはnet.ParseIPでIPv4として成功する場合だけ許可し、127.1/整数IPv4/8進表記をDNS名として再解釈しない。IP literalは標準net.ParseIPで解釈し、接続の可否は以下のpublic IP境界で別判定する。portは省略または10進1..65535、0/符号/前後空白を拒否する。

originの比較keyはhttps/lowercase DNS hostまたはcanonical IP/実効port（省略は443）のtuple。公開issuer_originはhttps://host、非443だけ:portを付加しIPv6を[]で囲む。redirect_uriの許可origin、metadata各endpointのsame-origin、UI authorization_url確認はこの同じ意味で比較する。issuer/sub/JWT iss自体はURI正規化して照合せず、設定のissuer文字列とのbyte一致を維持する。metadata URLの任意queryを認証requestへ持ち込まない。生成済みauthorization_urlだけは開始手順の固定queryを付加した別の出力型であり、provider設定URLのquery禁止を適用して拒否しない。各HTTP接続のDNS全結果/public IP/pinned dial/TLS ServerNameの既存境界をorigin一致だけで省略しない。

**OIDC接続IPの閉じた許可集合：** このv1のserver側discovery/JWKS/token接続だけに適用する。DNS結果0件/1件でも不許可なら全体拒否、検証済み集合からdial時の再DNSを行わない。IPv4はIsGlobalUnicastかつIsPrivate=falseを必要とし、さらに0.0.0.0/8、100.64.0.0/10、192.0.0.0/24、192.0.2.0/24、192.88.99.0/24、198.18.0.0/15、198.51.100.0/24、203.0.113.0/24、240.0.0.0/4を拒否する。IPv6は2000::/3内だけを許可し、2001::/23、2001:db8::/32、2002::/16、3fff::/20を拒否する。mapped IPv4はURL/DNS両方で拒否する。private/loopback/link-local/multicast/unspecifiedをpublic扱いするIsGlobalUnicast単独判定は禁止する。TCP443以外もURLの検証済みportだけ、HTTP Host/TLS ServerNameは元host、10秒deadlineはDNSからbody完了までを含む。fixtureのdocumentation address許可はfake adapter内の特別扱いだけでproductionへ伝播しない。

<a id="oidc-authentication-commit"></a>
**OIDC認証確定・失敗時境界：** 起動時にexternal auth全体enabledとprovider enabledの両方を確認する。startはOrigin/HTTP header検証→login group rate gate→JSON body検証/API token拒否→設定read→capacity事前確認→discovery→JWKS→4entropy取得→flow mutex内で期限切れ回収/再capacity確認→旧cookie flow置換/新flow保存→cookie/responseの順とする。upstream/entropy失敗は既存flowを変更しない。新flow保存後response書込失敗は当該新flowだけを失効し、旧flowを復活させない。confidential/publicのtoken endpoint認証方式はmetadataが列挙する場合にclient_secret_basic/noneがそれぞれ含まれることを確認し、列挙省略時はconfidentialだけclient_secret_basicとして扱いpublicは拒否する。metadata/JWKSの不整合・通信失敗は502でflowを作らない。

finishのtoken responseはUTF-8/duplicate keyなし/trailing JSONなしのJSON objectで、id_tokenが1..16384byteのstringであることを必要とする。上流HTTP失敗と同じくresponse形式不正は502、取り出したJWTの不正は401。token_type/access_token/refresh_token/expires_in等は認証判断に使わず保存しない。AuthorizationにAPI tokenを渡したstart/finishは401、sessionやBearerを新flowのuser指定へ利用しない。finishも同じOrigin/header→login rate→body/API token検証をflow照合より前に行う。multi-user-v1の認証後rate actor keyはsession:{user_id}:{group}、API tokenは既存token id keyとし、cicd-v1/legacyのkeyを変更しない。code交換後の再試行は全失敗とも新startから行う。

code交換中のprovider更新を通過させないため、flowは開始時providerのcanonical config digest（secret refを含みraw secretを含めない）を保持する。署名検証後にauthorization mutex内で最新provider有効性とdigest、user active/一意subject、revocation marker、role/permissions/TOTPを再検証する。provider不一致、対応user不在/disabled/lockedは401、状態破損/読取不能は500。ここでの照合と発行の間に更新を挟まず、外部通信をmutex内へ持ち込まない。

TOTP無効ならsession entropy生成→credential非nullのlogin_count/last_login_at更新→access log→audit log→session追加→responseとし、credential=nullなら更新を省略・must_change=noneとする。TOTP有効ならticket entropy生成→totp_requiredのaccess/audit→ticket追加→responseとし、login_count/last_login_at/sessionを更新しない。ticketはuser id/credential fingerprint（非nullは既存連結規則、nullはUTF-8 byteの `null` 4byteをSHA256）/permissions_hash/provider id/digest/subjectに束縛し、login/totpで最新値再検証後だけ既存TOTP確定順へ接続する。credential=nullでもTOTP stepのcompare-and-updateとログは必要、local資格情報を自動生成しない。必須logのaction/target/actorは既存login_success/totp_required/login_failure規則へ対応させ、OIDC専用のclaim診断は作らない。拒否401はlogin_failureのaccess/audit、上流502も同じfailure記録とし、log失敗なら500 Internal server errorを優先する。Origin/schema/capacity/entropyによる発行前拒否では認証結果logを追加しない（既存rate拒否監査はその正本を使用）。

multi-user memory sessionは既存保持値にauth_source=local|oidc、provider_id=null|string、provider_digest=null|Digest、provider_subject=null|stringを加え、localは3値null、oidcは3値必須とする。毎requestでuser/permissionsの既存照合に加えてOIDC providerの有効性/digestとuser.external_subjects内の同じprovider/subjectを確認し、不一致は401・session削除。provider更新/無効化はそのproviderの未完了flow/ticket/sessionを失効させ、external auth全体の無効化は全OIDC分を失効させる。userのexternal_subjects変更はそのuserのOIDC session/ticketを失効させ、local sessionと独立API tokenの自動失効へ拡張しない。session/ticket追加後response書込失敗は発行済みmemory recordを削除し、成功済みcredential/logは戻さず、消費済みflowを復活させない。発行途中entropy/credential/log失敗はsession/ticketを追加せず500、完了前に確定した永続更新は維持する。新しいsession fieldを永続fileまたはHTTP responseへ追加しない。

<a id="sec-27-58"></a>
**27.58 ロールベースアクセス制御：**

ロールベースアクセス制御の owner は `security` とする。

permission 名は以下に固定する。

| permission | 許可範囲 |
|------------|----------|
| `*` | 全操作。system admin role だけが保持できる。 |
| `status:read` | status、history、metrics、badge の read。 |
| `build:write` | build trigger、queue reorder、queue cancel。 |
| `config:read` | config、snapshot、template の read。 |
| `config:write` | config、snapshot create / restore、template apply。 |
| `user:admin` | user、role、external auth、share link 管理。 |
| `system:admin` | datastore、cache、retention、event stream、webhook resend。 |

API endpoint group と必要 permission は以下に固定する。ここにない追加管理 API endpoint は `403` とし、実装判断で既存 permission に割り当ててはならない。route、request、response は [`docs/details/api.md` 詳細本文責務 §27.48](api.md#sec-27-48)〜[§27.70](api.md#sec-27-70) を参照する。

| endpoint 群 | 必要 permission |
|--------------|----------------|
| `/api/users`, `/api/users/{id}`, `/api/roles`, `/api/roles/{id}`, `/api/external-auth-config`, `/api/external-auth/test`, `/api/share-links`, `/api/share-links/{id}`, `/api/share/{token}/status` | `user:admin`。ただし `/api/share/{token}/status` は share token 検証成功時だけ session permission 不要。 |
| `/api/datastore`, `/api/datastore/switch`, `/api/cache-policy`, `/api/response-cache`, `/api/events`, `/api/events/stream`, `/api/notify-log/{delivery_id}/resend` | `system:admin` |
| `/api/config-snapshots`, `/api/config-snapshots/{id}`, `/api/config-snapshots/{id}/restore`, `/api/config-snapshots/{left_id}/diff/{right_id}`, `/api/config-templates`, `/api/config-templates/{id}`, `/api/config-templates/{id}/apply` | read は `config:read`、write / restore / apply / delete は `config:write` |
| `/api/projects`, `/api/projects/{id}` | read は `config:read`、write / archive は `config:write` |
| `/api/queue/{queue_id}`, `/api/queue/reorder` | `build:write` |
| `/api/history/retention`, `/api/history/retention/run` | read は `status:read`、write / run は `system:admin` |
| `/api/stats/export`, `/api/metrics`, `/api/badge/status.svg`, `/api/snapshots/{left_id}/diff/{right_id}` | `status:read` |
| `/api/auth/me` | 有効な本人sessionのみ、permission不要/API token403。詳細は [本人HTTP責務](api.md#auth-self-http-contract)。 |
| `/api/version`, `/api/openapi.json`, `/api/auth/options` | 認証不要。個別permission不要。 |

`*` 以外の permission は暗黙に他 permission を含めない。

追加管理 API の認証・認可判定は以下の順に固定する。route、request、response、HTTP status の一覧は [`docs/details/api.md` 詳細本文責務 §27.48](api.md#sec-27-48)〜[§27.70](api.md#sec-27-70) を参照する。

| 順序 | 判定 | 失敗時 |
|------|------|--------|
| 1 | version/openapi/auth optionsは認証判定を行わない。 | 各HTTP正本のread-only副作用だけ。 |
| 2 | `/api/share/{token}/status` は share token を検証する。 | token 不正は `401`、不在または revoke 済みは `404`、期限切れは `410`。session を作成しない。 |
| 3 | その他の追加管理 API は session または API token を検証する。 | 不在、不正、期限切れは `401`。対象状態は読まない。 |
| 4 | user status を確認する。 | `disabled` または `locked` は `401`。対象状態は読まない。 |
| 5 | permission を判定する。 | 不足は `403`。対象状態は読まず、permission denied audit だけを許可する。 |
| 6 | endpoint 固有の user / role / share link 制約を確認する。 | 最後の admin 保護、system role 変更、role 使用中、self lock は `409`。 |

初期 system role は以下に固定する。

| role id | name | permissions | system |
|---------|------|-------------|--------|
| `admin` | `Administrator` | `*` | `true` |
| `operator` | `Operator` | `status:read`, `build:write`, `config:read` | `true` |
| `viewer` | `Viewer` | `status:read` | `true` |

system role は削除禁止、`id` 変更禁止、`system:false` への変更禁止とする。system roleへのPATCHでname/permissionsを指定するrequestは同値でも409とし、状態を変更しない。入力禁止id/system keyの422はAPI schema判定を優先する。

role削除時はstatusを問わず対象roleを参照するuserが1件でもあれば409とする。disabled/locked userの参照をdanglingにしてから再有効化する経路を作らない。上表はroute templateのexact集合であり、prefix wildcard、末尾slash、未知subrouteにpermissionを継承しない。version aliasはAPIが確定したcanonical route templateで判定する。

<a id="sec-27-61"></a>
**27.61 ユーザー管理 API：**

ユーザー管理 API の owner は `security` とする。

username は case-sensitive ではなく、保存前に lowercase に正規化する。

同一usernameはstatusを問わず全UserRecordで一意とし、locked userとの重複も409とする。ASCII uppercaseだけをlowercaseへ変換し、Unicode case fold、前後空白除去、全角変換は行わない。正規化後の保存型を満たさない値は422。

password 更新時は平文を保存せず、`pbkdf2_hmac_sha256_v1` だけを使用する。

UserRecord.credential=nullのuserはlocal password loginを禁止し、外部認証またはAPI tokenだけを許可する。資格情報の型は [`docs/details/statefile.md` 詳細本文責務 UserRecord](statefile.md#sec-22-0d)、公開投影は [`docs/details/api.md` 詳細本文責務 §27.48](api.md#sec-27-48) を参照する。

disabled user の session は次 request で `401` とする。

locked user の password login は `401` とし、API token と外部認証も拒否する。

user更新の最後のadmin/local credential保護と競合時判定は [マルチユーザー認証責務](#sec-27-48) を参照する。

<a id="sec-27-67"></a>
**27.67 読み取り専用共有リンク：**

読み取り専用共有リンクの owner は `security` とする。

share token は 32 byte 以上の乱数を URL-safe Base64 で表現する。

token 本体は作成 response で 1 回だけ返し、`.share_links` には SHA-256 hash だけを保存する。

share link scope は `status`、`history`、`snapshot_diff` のいずれかとする。

share link は read-only であり、build trigger、config update、queue 操作、user 操作、webhook resend を実行してはならない。

期限切れは `410`、revoke 済みは `404` とする。

share link request は通常 session を作成しない。

share link 作成時の `expires_at` は `null` または現在時刻より後の UTC ISO 8601 秒精度とする。過去時刻、現在時刻と同一秒、local timezone、offset 付き時刻は `422` とする。revoke は `revoked_at` が `null` の record だけを対象とし、revoke 済み record への再 revoke は `404` とする。

share token 生成は CSPRNG から 32 byte を取得し、padding なし URL-safe Base64 で表現する。生成 token の文字集合は `A-Z`、`a-z`、`0-9`、`-`、`_` だけとし、`=`、`+`、`/`、空白、改行を含めてはならない。CSPRNG 失敗時は `.share_links`、audit、admin event を変更せず `500` とする。token hash は token 文字列の UTF-8 byte 列ではなく、Base64 decode 後の 32 byte 以上の raw token byte に対する SHA-256 lowercase hex とする。

share token 検証は、(1) path parameter を percent decode する、(2) URL-safe Base64 padding なしとして decode する、(3) decode 後 byte length が 32 以上であることを確認する、(4) SHA-256 lowercase hex を算出する、(5) `.share_links.links[].token_hash` と constant-time 比較する、(6) `revoked_at`、(7) `expires_at` の順に行う。decode 不能または 32 byte 未満では `.share_links` を読まない。hash 不一致、不在、revoke 済みでは同じ `404` 結果とし、どの条件だったかを response、audit、admin event、access log で区別できる値として出さない。

share link 管理 API の `.audit_log` は [`docs/details/statefile.md` 詳細本文責務 `.audit_log` schema](statefile.md#audit-log-schema) の key だけを使用し、`action` は `share_link_create` または `share_link_revoke`、`target_type` は `share_link`、`target_id` は share link id、`actor_type` と `actor_id` は管理 session または API token とする。share link 管理 API の `.admin_events` は [`docs/details/statefile.md` 詳細本文責務 AdminEventRecord](statefile.md#sec-22-0d) の key だけを使用し、`type` は `security`、`target_type` は `share_link`、`target_id` は share link id、`actor` は管理 session user id または API token id、`message` は固定文言だけとする。`GET /api/share/{token}/status` は `.audit_log` と `.admin_events` を書き込まない。`.audit_log`、`.admin_events`、`.access_log`、`.api_access_log`、server log には token 本体、token hash、token prefix、token length、Authorization header、Cookie、raw path を保存してはならない。作成 response 返却後の token 再表示、list response への token 追加、SDK / UI による token 復元、log からの token 復元を禁止する。

share link の `scope` は `status`、`history`、`snapshot_diff` の単一値だけを許可する。複数 scope、wildcard、空配列、将来 scope 名、permission 名、API path を受け付けてはならない。`scope` ごとの response composition は [`docs/details/api.md` 詳細本文責務 §27.67](api.md#sec-27-67) を参照し、security owner は token、hash、期限、revoke、漏えい禁止だけを正本として持つ。

<a id="phase-13-security-alignment-contract"></a>
**Phase 13 security 実装整合契約：**

[`docs/details/security.md`](security.md) 詳細本文責務は、Phase 13 で password hash、token、署名、secret mask、file safety、credential rotation、必須 security log 境界を所有する。Phase 13 の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](../DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry)、完了証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 13 実装整合・品質改善証跡](fixture.md#phase-13-implementation-alignment-quality-evidence) を参照する。

| 対象 | 固定契約 |
|------|----------|
| password hash | 管理 credential の新規作成、password 変更、credential rotation では `pbkdf2_hmac_sha256_v1` だけを使用する。`pbkdf2_hmac_sha256_v1` は [`docs/SPEC.md` ポリシー責務 §4](../SPEC.md#policy-dependencies) の `INHOUSE` 分類であり、Go 標準ライブラリ primitive だけで実装する。独自 SHA-256 反復 KDF、単純 SHA-256、salt 連結 hash、外部 password hash / KDF library を新規保存形式として使用してはならない。 |
| migration | 既存 credential を読み込む場合、`pbkdf2_hmac_sha256_v1` 以外の algorithm を成功扱いにして温存しない。旧 `sha256_iter_v1` または未知 algorithm は credentials schema 不正として扱い、session、token、audit success、password change success を発生させない。 |
| token input | Admin CLI と MCP の token は token file または stdin から取得する。token の command argv 取得、process list に残る token、shell history に残る token を許可しない。 |
| token storage | token 本体は作成 response または stdin / file の入力 byte 以外へ保存しない。状態 file、log、fixture expected、release notes、PR 証跡には token hash、prefix、suffix、length を含めない。 |
| secret mask | PAT、Webhook secret、SMTP password、API token、session token、MCP token、share token、credential URL は読込直後に mask 登録し、stdout、stderr、server log、build log、audit、access、config log、fixture expected に平文を残さない。 |
| signature | Webhook、notification、release checksum、archive digest に関わる署名または digest の canonical input、algorithm、比較順、失敗時 response は owner 詳細本文から security owner へ到達できる。 |
| file safety | secret file は symlink 非追従、通常 file、mode `0600`、owner 一致、UTF-8 または byte 契約一致を満たす場合だけ読む。directory、device、FIFO、socket、hardlink count 2 以上は拒否する。 |
| required logs | security event の audit / access log が owner 詳細本文で必須書込みと定義される場合、書込み失敗を無視して成功 response を返してはならない。 |
| credential rotation | rotation は旧 credential 検証、new credential 生成、statefile atomic write、audit、旧 token revoke、rollback 不可範囲の明示を同一 operation として扱う。途中失敗時の session / token / audit の状態は fixture で固定する。 |

`pbkdf2_hmac_sha256_v1` credential hash の Phase 13 固定値は以下とする。下表の値を実装者判断で変更してはならない。値を変更する場合は [`docs/SPEC.md` ポリシー責務 §4](../SPEC.md#policy-dependencies) と本文節を先に改訂する。

| 項目 | 値 |
|------|----|
| classification | `INHOUSE` |
| algorithm | `pbkdf2_hmac_sha256_v1` |
| PRF | HMAC-SHA256 |
| iterations | `600000` |
| block count | `1` |
| salt length | 32 bytes |
| tag length | 32 bytes |
| salt source | `crypto/rand.Reader` |
| saved salt format | lowercase hex 64 文字。`.admin_credentials.salt` に保存する。 |
| saved tag format | lowercase hex 64 文字。`.admin_credentials.password_hash` に保存する。 |
| password input | UTF-8 byte sequence、Unicode normalization は行わない |
| maximum KDF input length | 1024 bytes。password validation の最大長は [password 入力制約](#認証共通詳細) の 128 文字を先に適用する。 |
| verification compare | constant-time compare |
| allowed primitives | `crypto/hmac`、`crypto/sha256`、`crypto/rand`、`crypto/subtle`、`encoding/hex` |
| forbidden dependency | `golang.org/x/*`、外部 password hash / KDF library、外部 PBKDF2 library、外部 test helper |
| legacy hash | `sha256_iter_v1` の新規生成と成功扱いを禁止する。既存検出時は credentials schema 不正として扱う。 |

`pbkdf2_hmac_sha256_v1` の実装確認 test vector は以下を固定値とする。implementation test は下表の `password` と `salt` を UTF-8 byte sequence として扱い、tag 32 bytes の lowercase hex が完全一致しなければならない。test vector 生成に外部 library、外部 command、外部 service を使用してはならない。

| case id | password | salt | iterations | expected password_hash |
|---------|----------|------|------------|------------------------|
| `pbkdf2-sha256-rfc-short-1` | `password` | `salt` | `1` | `120fb6cffcf8b32c43e7225256c4f837a86548c92ccc35480805987cb70be17b` |
| `pbkdf2-sha256-rfc-short-2` | `password` | `salt` | `2` | `ae4d0c95af6b46d32d0adff928f06dd02a303f8ef3c251dfd6e2d85a95474c43` |
| `pbkdf2-sha256-production-600000` | `password` | `salt` | `600000` | `669cfe52482116fda1aa2cbe409b2f56c8e4563752b7a28f6eaab614ee005178` |

Phase 13 の security evidence は `input/security_inventory.json` に `hash_vectors`、`token_sources`、`masked_outputs`、`file_safety_cases`、`rotation_cases` を持たせる。`hash_vectors` は raw password を保存せず、case id、salt hex、password_hash hex、verify input label、expected result だけを保存する。raw secret を fixture に保存した場合は `phase13_security_kdf_open_count` と `phase13_required_log_write_ignore_count` の両方へ未完了として計上する。

Phase 13 の `security` 実装は、[`docs/details/fixture.md` fixture 証跡責務 Phase 13 実装整合・品質改善証跡](fixture.md#phase-13-implementation-alignment-quality-evidence) の `phase13_security_kdf_open_count=0`、`phase13_token_arg_open_count=0`、`phase13_required_log_write_ignore_count=0` を満たすまで完了扱いにしてはならない。
