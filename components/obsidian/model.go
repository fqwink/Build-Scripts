package obsidian

const ownerName = "obsidian"
const obsidianInputMode = "obsidian-vault"
const obsidianMapSchemaVersion = "obsidian-map-v1"
const obsidianSyncStateSchemaVersion = "obsidian-sync-state-v1"
const obsidianSyncPlanSchemaVersion = "obsidian-sync-plan-v1"
const obsidianSyncRollbackSchemaVersion = "obsidian-sync-rollback-v1"

var obsidianBinaryVersion = "V.0.0-dev"

type obsidianConfig struct {
	InputMode   string
	VaultRoot   string
	EntryPath   string
	FilterFile  string
	ReportFile  string
	OutputRoot  string
	Strict      bool
	BuilderArgs []string
}

type obsidianSyncConfig struct {
	Action       string
	VaultRoot    string
	ProjectRoot  string
	StateDir     string
	PlanFile     string
	PlanHash     string
	RollbackFile string
	Direction    string
	DeletePolicy string
	ConflictDir  string
	TombstoneDir string
	OpenURI      bool
	OpenURISet   bool
}

type obsidianError struct {
	Code   string
	Exit   int
	Target string
	Line   int
	Column int
}

func (e obsidianError) Error() string {
	return e.Code
}

type obsidianFilter struct {
	IncludeNotes  []string
	ExcludeNotes  []string
	IncludeAssets []string
	ExcludeAssets []string
}

type vaultFile struct {
	SourcePath string
	AbsPath    string
	Size       int64
	SHA256     string
	Kind       string
}

type noteState struct {
	SourcePath     string
	NormalizedPath string
	RawText        string
	NormalizedText string
	Title          string
	SHA256         string
	TrailingLF     bool
	Headings       map[string]bool
	OutgoingLinks  []obsidianLink
	AssetEmbeds    []obsidianAsset
	Tags           []string
	Diagnostics    []obsidianDiagnostic
}

type obsidianMap struct {
	SchemaVersion       string               `json:"schema_version"`
	VaultDigest         string               `json:"vault_digest"`
	EntrySourcePath     string               `json:"entry_source_path"`
	EntryNormalizedPath string               `json:"entry_normalized_path"`
	Notes               []obsidianMapNote    `json:"notes"`
	Assets              []obsidianAsset      `json:"assets"`
	Diagnostics         []obsidianDiagnostic `json:"diagnostics"`
}

type obsidianMapNote struct {
	SourcePath       string               `json:"source_path"`
	NormalizedPath   string               `json:"normalized_path"`
	Title            string               `json:"title"`
	SourceSHA256     string               `json:"source_sha256"`
	SourceTrailingLF bool                 `json:"source_trailing_lf"`
	OutgoingLinks    []obsidianLink       `json:"outgoing_links"`
	AssetEmbeds      []obsidianAsset      `json:"asset_embeds"`
	Tags             []string             `json:"tags"`
	Diagnostics      []obsidianDiagnostic `json:"diagnostics"`
}

type obsidianLink struct {
	Raw        string `json:"raw"`
	Target     string `json:"target"`
	TargetPath string `json:"target_path"`
	Heading    string `json:"heading"`
	Alias      string `json:"alias"`
	Status     string `json:"status"`
	ErrorCode  string `json:"error_code"`
	Line       int    `json:"line"`
	Column     int    `json:"column"`
}

type obsidianAsset struct {
	Raw            string `json:"raw"`
	SourcePath     string `json:"source_path"`
	NormalizedPath string `json:"normalized_path"`
	SHA256         string `json:"sha256"`
	Size           int64  `json:"size"`
	MediaType      string `json:"media_type"`
	Status         string `json:"status"`
	ErrorCode      string `json:"error_code"`
	Line           int    `json:"line"`
	Column         int    `json:"column"`
}

type obsidianDiagnostic struct {
	Code       string `json:"code"`
	Severity   string `json:"severity"`
	SourcePath string `json:"source_path"`
	Line       int    `json:"line"`
	Column     int    `json:"column"`
	Target     string `json:"target"`
}

type obsidianSyncTree struct {
	RootDigest string
	Files      map[string]obsidianSyncFile
}

type obsidianSyncFile struct {
	Path   string
	Abs    string
	Digest string
	Size   int64
}

type obsidianSyncState struct {
	SchemaVersion     string                   `json:"schema_version"`
	ProjectRootDigest string                   `json:"project_root_digest"`
	VaultRootDigest   string                   `json:"vault_root_digest"`
	Entries           []obsidianSyncStateEntry `json:"entries"`
	Tombstones        []obsidianSyncTombstone  `json:"tombstones"`
	Conflicts         []obsidianSyncConflict   `json:"conflicts"`
	LastApplyID       string                   `json:"last_apply_id"`
}

type obsidianSyncStateEntry struct {
	Path           string `json:"path"`
	ProjectDigest  string `json:"project_digest"`
	VaultDigest    string `json:"vault_digest"`
	LastSyncDigest string `json:"last_sync_digest"`
	LastSyncUnix   int64  `json:"last_sync_unix"`
}

type obsidianSyncPlan struct {
	SchemaVersion     string                  `json:"schema_version"`
	PlanID            string                  `json:"plan_id"`
	CreatedAtUnix     int64                   `json:"created_at_unix"`
	Direction         string                  `json:"direction"`
	DeletePolicy      string                  `json:"delete_policy"`
	ConflictDir       string                  `json:"conflict_dir"`
	TombstoneDir      string                  `json:"tombstone_dir"`
	ProjectRootDigest string                  `json:"project_root_digest"`
	VaultRootDigest   string                  `json:"vault_root_digest"`
	StateDigest       string                  `json:"state_digest"`
	Operations        []obsidianSyncOperation `json:"operations"`
	Conflicts         []obsidianSyncConflict  `json:"conflicts"`
	Tombstones        []obsidianSyncTombstone `json:"tombstones"`
}

type obsidianSyncOperation struct {
	Op               string `json:"op"`
	Direction        string `json:"direction"`
	Path             string `json:"path"`
	Source           string `json:"source"`
	Destination      string `json:"destination"`
	BeforeDigest     string `json:"before_digest"`
	AfterDigest      string `json:"after_digest"`
	ConflictID       string `json:"conflict_id"`
	TombstoneID      string `json:"tombstone_id"`
	RollbackRequired bool   `json:"rollback_required"`
}

type obsidianSyncConflict struct {
	ConflictID    string `json:"conflict_id"`
	Path          string `json:"path"`
	Reason        string `json:"reason"`
	ProjectDigest string `json:"project_digest"`
	VaultDigest   string `json:"vault_digest"`
	StateDigest   string `json:"state_digest"`
	Resolution    string `json:"resolution"`
}

type obsidianSyncTombstone struct {
	TombstoneID   string `json:"tombstone_id"`
	Path          string `json:"path"`
	Digest        string `json:"digest"`
	Direction     string `json:"direction"`
	CreatedAtUnix int64  `json:"created_at_unix"`
	TombstonePath string `json:"tombstone_path"`
}

type obsidianSyncRollback struct {
	SchemaVersion     string                          `json:"schema_version"`
	ApplyID           string                          `json:"apply_id"`
	PlanHash          string                          `json:"plan_hash"`
	StartedAtUnix     int64                           `json:"started_at_unix"`
	CompletedAtUnix   int64                           `json:"completed_at_unix"`
	Operations        []obsidianSyncRollbackOperation `json:"operations"`
	StateBeforeDigest string                          `json:"state_before_digest"`
	StateAfterDigest  string                          `json:"state_after_digest"`
}

type obsidianSyncRollbackOperation struct {
	Op             string `json:"op"`
	Path           string `json:"path"`
	Destination    string `json:"destination"`
	PreviousDigest string `json:"previous_digest"`
	NewDigest      string `json:"new_digest"`
	BackupPath     string `json:"backup_path"`
	Status         string `json:"status"`
	ErrorCode      string `json:"error_code"`
}

type obsidianSyncPlanResponse struct {
	Command    string `json:"command"`
	PlanFile   string `json:"plan_file"`
	PlanHash   string `json:"plan_hash"`
	Operations int    `json:"operations"`
	Conflicts  int    `json:"conflicts"`
	Tombstones int    `json:"tombstones"`
	Applied    bool   `json:"applied"`
}

type obsidianSyncApplyResponse struct {
	Command           string `json:"command"`
	PlanFile          string `json:"plan_file"`
	PlanHash          string `json:"plan_hash"`
	OperationsApplied int    `json:"operations_applied"`
	Conflicts         int    `json:"conflicts"`
	Tombstones        int    `json:"tombstones"`
	RollbackFile      string `json:"rollback_file"`
	StateDigest       string `json:"state_digest"`
}

type obsidianSyncRollbackResponse struct {
	Command              string `json:"command"`
	RollbackFile         string `json:"rollback_file"`
	OperationsRolledBack int    `json:"operations_rolled_back"`
	Conflicts            int    `json:"conflicts"`
	StateDigest          string `json:"state_digest"`
}

type ownerFileContract struct {
	Owner string
	Files []string
}

func newOwnerFileContract() ownerFileContract {
	return ownerFileContract{
		Owner: ownerName,
		Files: []string{"obsidian.go", "model.go", "validate.go", "execute.go", "obsidian_test.go"},
	}
}
