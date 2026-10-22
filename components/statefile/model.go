package statefile

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/fqwink/build-scripts/components/security"
)

var ErrCICDSchema = errors.New("P19_SCHEMA")

var (
	ErrCICDCorrupt     = errors.New("P19_STATE_CORRUPT")
	ErrCICDBusy        = errors.New("P19_BUSY")
	ErrCICDRevision    = errors.New("P19_REVISION")
	ErrCICDStorage     = errors.New("P19_STORAGE")
	ErrCICDUncertain   = errors.New("P19_STATE_UNCERTAIN")
	ErrCICDCleanup     = errors.New("P19_CLEANUP")
	ErrCICDCapability  = errors.New("P19_CAPABILITY")
	ErrCICDIdempotency = errors.New("P19_IDEMPOTENCY")
	ErrCICDStale       = errors.New("P19_STALE")
	ErrCICDNotFound    = errors.New("P19_NOT_FOUND")
	ErrCICDInternal    = errors.New("P19_INTERNAL")
)

type cicdOperationRequest struct {
	id, digest, requestID string
	principal             security.CICDPrincipal
	access                security.CICDAccess
	check                 func() error
}

type cicdOperationMutation struct {
	targetID      string
	targetAttempt int
	result        json.RawMessage
	oneTime       string
}

type cicdOperationReply struct {
	operation  map[string]any
	generation int64
	replayed   bool
	result     json.RawMessage
	oneTime    string
}

const cicdSnapshotLimit = 33554432

// API admission modes, not a registry of executable business handlers.
var cicdOperationKinds = map[string]bool{
	"passwordChange": true, "projectPut": true, "workerPut": true, "targetPut": true, "userPut": true, "servicePut": true, "definitionPut": true, "templatePut": true, "environmentPut": true, "grantPut": true, "secretPut": true, "workerDrain": true, "artifactPin": true,
	"workerResume": false, "runStart": false, "runRetry": false, "runResume": false, "runCancel": false, "cacheDelete": false, "deploymentStart": false, "deploymentApprove": false, "deploymentReject": false, "deploymentCancel": false, "deploymentRollback": false, "deploymentReconcile": false,
}

type cicdHead struct {
	SchemaVersion int    `json:"schema_version"`
	Generation    int64  `json:"generation"`
	SHA256        string `json:"sha256"`
}

type cicdSnapshotStore struct {
	path       string
	observed   sync.Mutex
	generation int64
	digest     string
	uncertain  bool
	fault      func(string) error
	diagnostic func(error)
	clock      func() time.Time
}

// All stores share the process side of the filesystem lock, including aliases.
var cicdSnapshotMutex sync.Mutex

// Fixed storage types; the runtime never parses documentation.
var cicdStorageTypes = map[string]any{
	"Id":                "id",
	"JobId":             "jobid",
	"Digest":            "digest",
	"SHA":               "sha",
	"Time":              "utc",
	"RelPath":           "relpath",
	"AbsPath":           "abspath",
	"Glob":              "glob",
	"Ref":               "ref",
	"URL":               "url",
	"Scalar":            "scalar",
	"Bytes":             "base64",
	"Credential":        "legacy(AdminCredentials)",
	"TOTP":              "legacy(TOTPSecret)",
	"Quota":             map[string]string{"parallel": "int(1,16)", "queue": "int(1,512)", "artifact_mb": "int(128,102400)", "cache_mb": "int(0,10240)", "retention_days": "int(1,30)"},
	"Resources":         map[string]string{"cpu_percent": "int(10,400)", "memory_mb": "int(256,8192)", "disk_mb": "int(128,2048)", "tasks": "int(8,256)", "timeout_seconds": "int(1,3600)", "output_mb": "int(1,16)"},
	"Input":             map[string]string{"kind": "enum(string,integer,boolean)", "required": "bool", "default": "nullable(Scalar)", "choices": "list(Scalar,0,32)", "min": "int(-2147483648,2147483647)", "max": "int(-2147483648,2147483647)"},
	"Condition":         map[string]string{"branches": "list(Glob,0,32)", "tags": "list(Glob,0,32)", "paths": "list(Glob,0,32)", "input_equals": "map(Scalar,0,32)"},
	"Matrix":            map[string]string{"axes": "map(ScalarList,0,4)", "exclude": "list(ScalarMap,0,64)", "include": "list(ScalarMap,0,64)"},
	"ScalarList":        "list(Scalar,1,32)",
	"ScalarMap":         "map(Scalar,0,32)",
	"Retry":             map[string]string{"max_attempts": "int(1,3)", "backoff_seconds": "int(1,30)"},
	"CachePolicy":       map[string]string{"enabled": "bool", "paths": "list(RelPath,0,16)", "ttl_seconds": "int(60,2592000)"},
	"SecretRef":         map[string]string{"id": "Id", "version": "int(1,2147483647)"},
	"BuilderTask":       map[string]string{"kind": "enum(builder)", "source": "RelPath", "output": "RelPath"},
	"GoTestTask":        map[string]string{"kind": "enum(go-test)", "packages": "list(RelPath,1,64)", "race": "const(false)"},
	"GoVetTask":         map[string]string{"kind": "enum(go-vet)", "packages": "list(RelPath,1,64)"},
	"GoBuildTask":       map[string]string{"kind": "enum(go-build)", "package": "RelPath", "output": "RelPath"},
	"DenoTask":          map[string]string{"kind": "enum(deno-check)", "files": "list(RelPath,1,64)"},
	"VerifyTask":        map[string]string{"kind": "enum(artifact-verify)", "artifact": "Digest"},
	"DeployTask":        map[string]string{"kind": "enum(deploy)", "environment_id": "Id", "artifact_from": "JobId"},
	"Task":              "oneof(BuilderTask,GoTestTask,GoVetTask,GoBuildTask,DenoTask,VerifyTask,DeployTask)",
	"Job":               map[string]string{"id": "JobId", "needs": "list(JobId,0,32)", "condition": "Condition", "task": "Task", "matrix": "Matrix", "resources": "Resources", "retry": "Retry", "cache": "CachePolicy", "env": "envmap(EnvValue,0,64)", "secrets": "envmap(SecretRef,0,16)", "optional": "bool"},
	"EnvValue":          "text(0,4096)",
	"Watch":             map[string]string{"id": "Id", "ref_kind": "enum(branch,tag)", "ref": "Ref", "interval_seconds": "int(60,3600)"},
	"Schedule":          map[string]string{"id": "Id", "ref_kind": "enum(branch,tag)", "ref": "Ref", "anchor": "Time", "every_seconds": "int(60,604800)", "catch_up": "enum(skip,latest)", "enabled": "bool"},
	"TemplateUse":       map[string]string{"id": "Id", "revision": "int(1,2147483647)", "prefix": "JobId", "values": "ScalarMap"},
	"Definition":        map[string]string{"project_id": "Id", "id": "Id", "revision": "int(1,2147483647)", "inputs": "map(Input,0,32)", "manual": "bool", "watches": "list(Watch,0,32)", "schedules": "list(Schedule,0,32)", "templates": "list(TemplateUse,0,16)", "jobs": "list(Job,1,128)", "fail_fast": "bool"},
	"Template":          map[string]string{"project_id": "Id", "id": "Id", "revision": "int(1,2147483647)", "inputs": "map(Input,0,32)", "jobs": "list(Job,1,128)"},
	"Source":            map[string]string{"owner": "text(1,100)", "repo": "text(1,100)", "ref_kind": "enum(branch,tag)", "ref": "Ref", "sha": "SHA", "tree_digest": "Digest", "files": "list(SourceFile,0,10000)"},
	"PlannedDeployTask": map[string]string{"kind": "enum(deploy)", "environment_id": "Id", "artifact_from": "Id"},
	"PlannedTask":       "oneof(BuilderTask,GoTestTask,GoVetTask,GoBuildTask,DenoTask,VerifyTask,PlannedDeployTask)",
	"PlannedJob":        map[string]string{"id": "Id", "base_job_id": "JobId", "selected": "bool", "matrix_values": "ScalarMap", "needs": "list(Id,0,256)", "task": "PlannedTask", "resources": "Resources", "retry": "Retry", "cache": "CachePolicy", "env": "envmap(EnvValue,0,64)", "secrets": "envmap(SecretRef,0,16)", "optional": "bool"},
	"PlanData":          map[string]string{"project_id": "Id", "definition_id": "Id", "definition_revision": "int(1,2147483647)", "source": "Source", "values": "ScalarMap", "jobs": "list(PlannedJob,1,256)", "fail_fast": "bool", "toolchain_digest": "Digest"},
	"Plan":              map[string]string{"digest": "Digest", "data": "PlanData"},
	"Project":           map[string]string{"id": "Id", "name": "text(1,128)", "revision": "int(1,2147483647)", "enabled": "bool", "source_owner": "text(1,100)", "source_repo": "text(1,100)", "definition_id": "nullable(Id)", "definition_revision": "int(0,2147483647)", "quota": "Quota", "toolchain_digest": "nullable(Digest)"},
	"Stage":             map[string]string{"target_ids": "list(Id,1,32)", "percent": "int(1,100)", "observe_seconds": "int(10,3600)", "max_error_percent": "int(0,100)", "min_requests": "int(0,10000)"},
	"Freeze":            map[string]string{"from": "Time", "until": "Time"},
	"Environment":       map[string]string{"project_id": "Id", "id": "Id", "revision": "int(1,2147483647)", "production": "bool", "approvals": "int(0,2)", "allow_self_approval": "bool", "targets": "list(Id,1,32)", "strategy": "enum(rolling,blue-green,canary)", "stages": "list(Stage,1,16)", "freeze": "list(Freeze,0,32)", "health_timeout_seconds": "int(1,60)", "retained_successes": "int(2,30)", "health_path": "enum(index.html)"},
	"Grant":             map[string]string{"subject_kind": "enum(user,service)", "subject_id": "Id", "project_id": "nullable(Id)", "environment_id": "nullable(Id)", "role": "enum(viewer,operator,approver,project-admin,platform-admin)", "enabled": "bool", "expires_at": "nullable(Time)"},
	"User":              map[string]string{"id": "Id", "name": "text(1,128)", "revision": "int(1,2147483647)", "enabled": "bool", "credential": "nullable(Credential)", "totp": "nullable(TOTP)", "created_at": "Time"},
	"Service":           map[string]string{"id": "Id", "revision": "int(1,2147483647)", "token_sha256": "Digest", "enabled": "bool", "expires_at": "Time"},
	"Lease":             map[string]string{"project_id": "Id", "run_id": "Id", "attempt": "int(1,100)", "job_id": "Id", "generation": "int(1,9007199254740991)", "until": "Time"},
	"Worker":            map[string]string{"id": "Id", "revision": "int(1,2147483647)", "enabled": "bool", "draining": "bool", "host": "text(1,253)", "port": "int(1,65535)", "user": "Id", "allowed_projects": "list(Id,1,64)", "capabilities": "list(text(1,64),1,16)", "slots": "int(1,4)", "last_seen": "nullable(Time)", "leases": "list(Lease,0,4)", "toolchain_digest": "Digest", "total_memory_mb": "int(896,1048576)"},
	"JobResult":         map[string]string{"id": "Id", "execution_id": "nullable(Id)", "state": "enum(queued,assigned,running,cancel_requested,succeeded,failed,cancelled,skipped,recovery_required)", "try_number": "int(0,3)", "worker_id": "nullable(Id)", "lease_generation": "int(0,9007199254740991)", "started_at": "nullable(Time)", "finished_at": "nullable(Time)", "duration_ms": "nullable(int(0,9007199254740991))", "code": "nullable(text(1,64))", "artifact": "nullable(Digest)", "reused_from_attempt": "nullable(int(1,100))", "log_gap": "bool"},
	"Run":               map[string]string{"project_id": "Id", "id": "Id", "revision": "int(1,2147483647)", "attempt": "int(1,100)", "event_key": "Digest", "requester": "Id", "plan": "Plan", "state": "enum(queued,running,cancel_requested,succeeded,failed,cancelled,recovery_required)", "created_at": "Time", "started_at": "nullable(Time)", "finished_at": "nullable(Time)", "jobs": "list(JobResult,1,256)", "cancel_requested": "bool", "trigger": "enum(manual,watch,schedule,retry,resume)"},
	"Approval":          map[string]string{"subject": "Id", "at": "Time", "plan_digest": "Digest", "expires_at": "Time"},
	"DeployPlan":        map[string]string{"project_id": "Id", "environment_id": "Id", "environment_revision": "int(1,2147483647)", "expected_generation": "int(0,9007199254740991)", "artifact": "Digest", "config_digest": "Digest", "strategy": "enum(rolling,blue-green,canary)", "stages": "list(Stage,1,16)", "rollback_of": "nullable(Id)"},
	"Deployment":        map[string]string{"id": "Id", "revision": "int(1,2147483647)", "plan": "DeployPlan", "plan_digest": "Digest", "requester": "Id", "state": "enum(waiting_approval,queued,staging,switching,observing,succeeded,failed,rolling_back,rolled_back,cancel_requested,cancelled,recovery_required)", "generation": "int(1,9007199254740991)", "created_at": "Time", "finished_at": "nullable(Time)", "approvals": "list(Approval,0,2)", "stage_index": "int(0,16)", "previous_artifact": "nullable(Digest)", "code": "nullable(text(1,64))", "recovery_code": "nullable(text(1,64))", "targets": "list(TargetResult,1,32)", "observations": "list(StageObservation,0,512)"},
	"Artifact":          map[string]string{"project_id": "Id", "digest": "Digest", "manifest_digest": "Digest", "size_bytes": "int(1,1073741824)", "created_at": "Time", "source_sha": "SHA", "run_id": "Id", "attempt": "int(1,100)", "job_id": "Id", "pins": "list(Id,0,1024)", "deleting": "bool"},
	"CacheEntry":        map[string]string{"project_id": "Id", "key": "Digest", "digest": "Digest", "size_bytes": "int(1,1073741824)", "created_at": "Time", "last_used_at": "Time", "expires_at": "Time", "in_use": "list(Id,0,256)", "deleting": "bool"},
	"Operation":         map[string]string{"id": "Id", "project_id": "nullable(Id)", "subject": "Id", "request_digest": "Digest", "kind": "text(1,64)", "state": "enum(accepted,succeeded,failed,recovery_required)", "target_id": "nullable(text(1,64))", "target_attempt": "nullable(int(1,100))", "created_at": "Time", "expires_at": "Time", "code": "nullable(text(1,64))"},
	"StageObservation":  map[string]string{"target_id": "Id", "stage_index": "int(1,16)", "generation": "int(1,9007199254740991)", "step": "int(1,128)", "from": "Time", "until": "Time", "observed_until": "Time", "requests": "int(0,9007199254740991)", "errors": "int(0,9007199254740991)", "state": "enum(observing,passed,failed,unknown)", "code": "nullable(text(1,64))"},
	"Audit":             map[string]string{"operation_id": "Id", "subject": "Id", "project_id": "nullable(Id)", "environment_id": "nullable(Id)", "action": "text(1,64)", "phase": "enum(accepted,finished,rejected)", "at": "Time", "target_id": "nullable(text(1,64))", "code": "nullable(text(1,64))"},
	"Route":             map[string]string{"project_id": "Id", "environment_id": "Id", "generation": "int(0,9007199254740991)", "operation_id": "Id", "stable": "nullable(Digest)", "candidate": "nullable(Digest)", "percent": "int(0,100)", "config_digest": "Digest"},
	"Log":               map[string]string{"schema_version": "const(1)", "project_id": "Id", "run_id": "Id", "attempt": "int(1,100)", "job_id": "Id", "stream": "enum(stdout,stderr)", "first_offset": "int(0,9007199254740991)", "next_offset": "int(0,9007199254740991)", "data": "text(0,1048576)", "truncated": "bool", "gap": "bool", "terminal": "bool"},
	"SecretMeta":        map[string]string{"project_id": "Id", "id": "Id", "version": "int(1,2147483647)", "created_at": "Time"},
	"Root":              map[string]string{"schema_version": "const(1)", "generation": "int(1,9007199254740991)", "projects": "list(Project,0,64)", "definitions": "list(Definition,0,1024)", "templates": "list(Template,0,1024)", "environments": "list(Environment,0,1024)", "users": "list(User,1,256)", "services": "list(Service,0,256)", "secrets": "list(SecretMeta,0,10000)", "grants": "list(Grant,1,4096)", "workers": "list(Worker,0,32)", "runs": "list(Run,0,512)", "deployments": "list(Deployment,0,1024)", "artifacts": "list(Artifact,0,10000)", "caches": "list(CacheEntry,0,10000)", "operations": "list(Operation,0,10000)", "audits": "list(Audit,0,20000)", "targets": "list(Target,0,1024)", "environment_states": "list(EnvironmentState,0,1024)", "notifications": "list(Notification,0,10000)", "trigger_cursors": "list(TriggerCursor,0,4096)", "status_deliveries": "list(StatusDelivery,0,1024)"},
	"Head":              map[string]string{"schema_version": "const(1)", "generation": "int(1,9007199254740991)", "sha256": "Digest"},
	"Master":            map[string]string{"schema_version": "const(1)", "encryption_key": "Bytes", "cursor_key": "Bytes"},
	"SecretCipher":      map[string]string{"schema_version": "const(1)", "project_id": "Id", "id": "Id", "version": "int(1,2147483647)", "nonce": "Bytes", "ciphertext": "Bytes", "created_at": "Time"},
	"Target":            map[string]string{"id": "Id", "revision": "int(1,2147483647)", "project_id": "Id", "host": "text(1,253)", "port": "int(1,65535)", "user": "Id", "base_url": "URL", "enabled": "bool"},
	"EnvironmentState":  map[string]string{"project_id": "Id", "environment_id": "Id", "generation": "int(0,9007199254740991)", "artifact": "nullable(Digest)", "active_deployment": "nullable(Id)", "blocked": "bool"},
	"TargetResult":      map[string]string{"target_id": "Id", "generation": "int(0,9007199254740991)", "artifact": "nullable(Digest)", "phase": "enum(pending,staged,switched,healthy,restored,failed,unknown)", "next_step": "int(1,128)", "previous_route": "nullable(Route)", "code": "nullable(text(1,64))"},
	"Notification":      map[string]string{"id": "Id", "project_id": "Id", "event": "enum(queued,started,finished,cancelled,waiting_approval,stopped)", "target_id": "Id", "target_attempt": "int(1,100)", "source_revision": "int(1,2147483647)", "attempt": "int(1,3)", "state": "enum(queued,sent,failed)", "next_at": "Time", "code": "nullable(text(1,64))", "at": "Time"},
	"SourceFile":        map[string]string{"path": "RelPath", "blob_sha": "SHA", "sha256": "Digest", "size": "int(0,10485760)", "executable": "bool"},
	"TriggerCursor":     map[string]string{"project_id": "Id", "definition_id": "Id", "definition_revision": "int(1,2147483647)", "trigger_id": "Id", "kind": "enum(watch,schedule)", "last_sha": "nullable(SHA)", "last_due": "nullable(Time)"},
	"StatusDelivery":    map[string]string{"project_id": "Id", "run_id": "Id", "attempt": "int(1,100)", "state": "enum(pending,success,failure,error)", "sent": "bool", "code": "nullable(text(1,64))", "at": "Time"},
	"SiteResult":        map[string]string{"route": "Route", "verified": "bool", "requests": "int(0,9007199254740991)", "errors": "int(0,9007199254740991)", "code": "nullable(text(1,64))"},
	"SiteOperation":     map[string]string{"operation_id": "Id", "step": "int(1,128)", "type": "enum(stage,switch,finalize,restore,delete)", "request_digest": "Digest", "committed_at": "Time", "result": "SiteResult"},
	"SiteWindow":        map[string]string{"from": "Time", "until": "Time", "requests": "int(0,9007199254740991)", "errors": "int(0,9007199254740991)"},
	"SiteCounters":      map[string]string{"step": "int(0,128)", "generation": "int(0,9007199254740991)", "windows": "list(SiteWindow,0,8640)"},
	"WorkerExecution":   map[string]string{"operation_id": "Id", "assignment_bytes": "int(1,2097152)", "assignment_digest": "Digest", "toolchain_digest": "Digest", "project_id": "Id", "run_id": "Id", "attempt": "int(1,100)", "job_id": "Id", "lease_generation": "int(1,9007199254740991)", "lease_until": "Time", "unit_names": "list(text(1,128),0,65)", "resources": "Resources", "job": "JobResult", "updated_at": "Time", "cleanup_ok": "bool"},
	"SiteState":         map[string]string{"schema_version": "const(1)", "route": "Route", "active_operation": "nullable(Id)", "previous_route": "nullable(Route)", "operations": "list(SiteOperation,0,1024)", "counters": "SiteCounters"},
	"SetupState":        map[string]string{"schema_version": "const(1)", "trust_revision": "int(1,2147483647)", "trust_digest": "Digest", "revoked_key_ids": "list(Digest,0,32)", "current_release": "text(1,128)", "previous_releases": "list(text(1,128),0,2)"},
}

type ownerFileContract struct {
	Owner string
	Files []string
}

type LockOptions struct {
	Owner          string
	StaleAfter     time.Duration
	RetryInterval  time.Duration
	AcquireTimeout time.Duration
	Clock          func() time.Time
	Sleep          func(time.Duration)
}

type LockHandle struct {
	path string
	file string
}

func newOwnerFileContract() ownerFileContract {
	return ownerFileContract{
		Owner: Owner(),
		Files: []string{"statefile.go", "model.go", "validate.go", "execute.go", "statefile_test.go"},
	}
}
