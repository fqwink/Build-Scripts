package obsidian

const ownerName = "obsidian"
const obsidianInputMode = "obsidian-vault"
const obsidianMapSchemaVersion = "obsidian-map-v1"

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
