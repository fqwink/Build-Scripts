package setup

type CICDConfig struct {
	SchemaVersion        int                      `json:"schema_version"`
	Profile              string                   `json:"profile"`
	Role                 string                   `json:"role"`
	ID                   string                   `json:"id"`
	AllowedProjects      []string                 `json:"allowed_projects"`
	AllowedEnvironments  []CICDEnvironmentBinding `json:"allowed_environments"`
	User                 string                   `json:"user"`
	StateDir             string                   `json:"state_dir"`
	ReleaseID            string                   `json:"release_id"`
	Listen               string                   `json:"listen"`
	TLSCertFile          *string                  `json:"tls_cert_file"`
	TLSKeyFile           *string                  `json:"tls_key_file"`
	AllowHTTP            bool                     `json:"allow_http"`
	WorkerSlots          int                      `json:"worker_slots"`
	Toolchains           []CICDToolchain          `json:"toolchains"`
	SSHAuthorizedKeyFile *string                  `json:"ssh_authorized_key_file"`
}

type CICDEnvironmentBinding struct {
	ProjectID     string `json:"project_id"`
	EnvironmentID string `json:"environment_id"`
	Host          string `json:"host"`
	Port          int    `json:"port"`
}

type CICDToolchain struct {
	Capability string              `json:"capability"`
	Version    string              `json:"version"`
	Root       string              `json:"root"`
	Files      []CICDToolchainFile `json:"files"`
}

type CICDToolchainFile struct {
	Path       string `json:"path"`
	SHA256     string `json:"sha256"`
	Executable bool   `json:"executable"`
}

type ownerFileContract struct {
	Owner string
	Files []string
}

func newOwnerFileContract() ownerFileContract {
	return ownerFileContract{
		Owner: "setup",
		Files: []string{"setup.go", "model.go", "validate.go", "execute.go", "setup_test.go"},
	}
}
