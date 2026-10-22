package security

import (
	"encoding/base64"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestCICDBearerAndCredentialBinding(t *testing.T) {
	service := base64.RawURLEncoding.EncodeToString(make([]byte, 32))
	for _, test := range []struct{ token, kind string }{
		{strings.Repeat("a", 64), "user"}, {service, "service"},
		{strings.Repeat("A", 64), ""}, {"acs_" + strings.Repeat("a", 64), ""},
		{service + "=", ""}, {service[:42] + "B", ""}, {" " + service, ""},
		{service + "\n", ""}, {strings.Repeat("a", 63), ""}, {"", ""},
	} {
		if CICDBearerKind(test.token) != test.kind {
			t.Fatal("security.phase19-auth.bearer.kind")
		}
	}
	record, err := NewPasswordRecord("Fixture-password-2026")
	if err != nil {
		t.Fatal("security.phase19-auth.credential.fixture")
	}
	credential := CICDCredential{PasswordHash: record.PasswordHash, Salt: record.Salt, Algorithm: record.Algorithm, Iterations: record.Iterations, UpdatedAt: "2026-10-22T00:00:00Z"}
	if !VerifyCICDPassword(credential, "Fixture-password-2026") || VerifyCICDPassword(credential, "wrong-password") {
		t.Fatal("security.phase19-auth.credential.verify")
	}
	binding := CICDCredentialFingerprint(credential)
	at := "2026-10-22T00:01:00Z"
	credential.LoginCount, credential.LastLoginAt = 5, &at
	if CICDCredentialFingerprint(credential) != binding {
		t.Fatal("security.phase19-auth.credential.login-stable")
	}
	credential.MustChange = true
	if CICDCredentialFingerprint(credential) == binding {
		t.Fatal("security.phase19-auth.credential.changed")
	}
	credential.Iterations = 1
	if VerifyCICDPassword(credential, "Fixture-password-2026") {
		t.Fatal("security.phase19-auth.credential.invalid")
	}
}

func cicdAuthorizationFixture() (CICDPrincipal, CICDAuthorizationSnapshot, time.Time) {
	now := time.Date(2026, 10, 22, 0, 0, 0, 0, time.UTC)
	principal := CICDPrincipal{Kind: "user", ID: "alice", ExpiresAt: now.Add(time.Hour), MustChange: "none"}
	snapshot := CICDAuthorizationSnapshot{Generation: 7, Subject: CICDSubject{Kind: "user", ID: "alice", Enabled: true}, ProjectExists: true, ProjectEnabled: true, EnvironmentExists: true}
	return principal, snapshot, now
}

func TestCICDAuthorizationMatrix(t *testing.T) {
	t.Parallel()
	for _, group := range []struct{ operations, roles, scope string }{
		{"projectPut workerPut targetPut userPut servicePut workerDrain workerResume workersList targetsList usersList servicesList", "platform-admin", "global"},
		{"passwordChange logout", "viewer operator project-admin approver platform-admin", "global"},
		{"projectsList", "viewer operator project-admin approver platform-admin", "optional"},
		{"grantPut grantsList", "project-admin platform-admin", "project"},
		{"plan definitionsList templatesList environmentsList runsList runGet logsGet artifactsList artifactGet artifactDownload artifactCompare cachesList environmentStatus deploymentsList metricsGet", "viewer operator project-admin", "project"},
		{"runStart runRetry runResume runCancel artifactPin deploymentPlan deploymentStart deploymentCancel deploymentRollback", "operator project-admin", "project"},
		{"definitionPut templatePut environmentPut secretPut secretsList cacheDelete auditsList operationGet", "project-admin", "project"},
		{"deploymentGet", "viewer operator project-admin approver", "environment"},
		{"deploymentApprove deploymentReject", "approver", "environment"},
		{"deploymentReconcile", "platform-admin", "environment"},
	} {
		for _, operation := range strings.Fields(group.operations) {
			for _, role := range []string{"viewer", "operator", "project-admin", "approver", "platform-admin"} {
				t.Run(operation+"/"+role, func(t *testing.T) {
					principal, snapshot, now := cicdAuthorizationFixture()
					grant := CICDGrant{SubjectKind: "user", SubjectID: "alice", Role: role, Enabled: true, ProjectID: "p"}
					if role == "platform-admin" {
						grant.ProjectID = ""
					}
					if role == "approver" {
						grant.EnvironmentID = "prod"
					}
					snapshot.Grants = []CICDGrant{grant}
					access := CICDAccess{Operation: operation, ResourceOwner: "other"}
					if group.scope == "project" || group.scope == "environment" {
						access.ProjectID = "p"
					}
					if group.scope == "environment" {
						access.EnvironmentID = "prod"
					}
					if operation == "passwordChange" {
						access.ResourceOwner = "alice"
					}
					if operation == "grantPut" {
						access.ProposedGrant = &CICDGrant{SubjectKind: "user", SubjectID: "bob", ProjectID: "p", Role: "viewer", Enabled: true}
					}
					want := strings.Contains(" "+group.roles+" ", " "+role+" ")
					got, err := AuthorizeCICD(principal, snapshot, access, now)
					if (err == nil) != want || !want && (err != ErrCICDForbidden || !reflect.DeepEqual(got, CICDAuthorization{})) {
						t.Fatal("security.phase19-access.matrix.result")
					}
					if want && (got.SubjectID != principal.ID || got.Generation != 7 || len(got.Grants) != 1) {
						t.Fatal("security.phase19-access.matrix.projection")
					}
				})
			}
		}
	}
}

func TestCICDAuthorizationBoundaries(t *testing.T) {
	for _, test := range []struct {
		id     string
		change func(*CICDPrincipal, *CICDAuthorizationSnapshot, *CICDAccess, time.Time)
		want   error
	}{
		{"expired-session", func(p *CICDPrincipal, _ *CICDAuthorizationSnapshot, _ *CICDAccess, now time.Time) { p.ExpiresAt = now }, ErrCICDAuth},
		{"disabled-subject", func(_ *CICDPrincipal, s *CICDAuthorizationSnapshot, _ *CICDAccess, _ time.Time) {
			s.Subject.Enabled = false
		}, ErrCICDAuth},
		{"unknown-subject", func(_ *CICDPrincipal, s *CICDAuthorizationSnapshot, _ *CICDAccess, _ time.Time) { s.Subject.ID = "bob" }, ErrCICDAuth},
		{"reserved-subject", func(p *CICDPrincipal, s *CICDAuthorizationSnapshot, _ *CICDAccess, _ time.Time) {
			p.ID, s.Subject.ID = "system", "system"
		}, ErrCICDAuth},
		{"forced", func(p *CICDPrincipal, _ *CICDAuthorizationSnapshot, _ *CICDAccess, _ time.Time) {
			p.MustChange = "forced"
		}, ErrCICDForbidden},
		{"unknown-must-change", func(p *CICDPrincipal, _ *CICDAuthorizationSnapshot, _ *CICDAccess, _ time.Time) {
			p.MustChange = "unknown"
		}, ErrCICDAuth},
		{"expired-grant", func(_ *CICDPrincipal, s *CICDAuthorizationSnapshot, _ *CICDAccess, now time.Time) {
			s.Grants[0].ExpiresAt = &now
		}, ErrCICDForbidden},
		{"disabled-grant", func(_ *CICDPrincipal, s *CICDAuthorizationSnapshot, _ *CICDAccess, _ time.Time) {
			s.Grants[0].Enabled = false
		}, ErrCICDForbidden},
		{"other-project", func(_ *CICDPrincipal, _ *CICDAuthorizationSnapshot, a *CICDAccess, _ time.Time) { a.ProjectID = "q" }, ErrCICDForbidden},
		{"other-subject-grant", func(_ *CICDPrincipal, s *CICDAuthorizationSnapshot, _ *CICDAccess, _ time.Time) {
			s.Grants[0].SubjectID = "bob"
		}, ErrCICDForbidden},
		{"unknown-role", func(_ *CICDPrincipal, s *CICDAuthorizationSnapshot, _ *CICDAccess, _ time.Time) {
			s.Grants[0].Role = "owner"
		}, ErrCICDForbidden},
		{"unknown-operation", func(_ *CICDPrincipal, _ *CICDAuthorizationSnapshot, a *CICDAccess, _ time.Time) {
			a.Operation = "executeAnything"
		}, ErrCICDForbidden},
		{"missing-project", func(_ *CICDPrincipal, s *CICDAuthorizationSnapshot, _ *CICDAccess, _ time.Time) {
			s.ProjectExists = false
		}, ErrCICDForbidden},
		{"missing-scope", func(_ *CICDPrincipal, _ *CICDAuthorizationSnapshot, a *CICDAccess, _ time.Time) { a.ProjectID = "" }, ErrCICDForbidden},
		{"disabled-project-start", func(_ *CICDPrincipal, s *CICDAuthorizationSnapshot, a *CICDAccess, _ time.Time) {
			s.ProjectEnabled = false
			a.Operation = "runStart"
		}, ErrCICDForbidden},
		{"disabled-project-read", func(_ *CICDPrincipal, s *CICDAuthorizationSnapshot, _ *CICDAccess, _ time.Time) {
			s.ProjectEnabled = false
		}, nil},
		{"own-operation", func(_ *CICDPrincipal, _ *CICDAuthorizationSnapshot, a *CICDAccess, _ time.Time) {
			a.Operation, a.ResourceOwner = "operationGet", "alice"
		}, nil},
		{"other-operation", func(_ *CICDPrincipal, _ *CICDAuthorizationSnapshot, a *CICDAccess, _ time.Time) {
			a.Operation, a.ResourceOwner = "operationGet", "bob"
		}, ErrCICDForbidden},
		{"approver-other-environment", func(_ *CICDPrincipal, s *CICDAuthorizationSnapshot, a *CICDAccess, _ time.Time) {
			s.Grants[0].Role, s.Grants[0].EnvironmentID = "approver", "staging"
			a.Operation, a.EnvironmentID = "deploymentApprove", "prod"
		}, ErrCICDForbidden},
		{"service-approve", func(p *CICDPrincipal, s *CICDAuthorizationSnapshot, a *CICDAccess, now time.Time) {
			p.Kind, s.Subject.Kind, s.Grants[0].SubjectKind = "service", "service", "service"
			expiry := now.Add(time.Hour)
			s.Subject.ExpiresAt = &expiry
			s.Grants[0].Role, s.Grants[0].EnvironmentID = "approver", "prod"
			a.Operation, a.EnvironmentID = "deploymentApprove", "prod"
		}, ErrCICDForbidden},
		{"expired-service", func(p *CICDPrincipal, s *CICDAuthorizationSnapshot, _ *CICDAccess, now time.Time) {
			p.Kind, s.Subject.Kind, s.Grants[0].SubjectKind = "service", "service", "service"
			s.Subject.ExpiresAt = &now
		}, ErrCICDAuth},
	} {
		t.Run(test.id, func(t *testing.T) {
			p, s, now := cicdAuthorizationFixture()
			s.Grants = []CICDGrant{{SubjectKind: "user", SubjectID: "alice", ProjectID: "p", Role: "operator", Enabled: true}}
			a := CICDAccess{Operation: "runGet", ProjectID: "p"}
			test.change(&p, &s, &a, now)
			got, err := AuthorizeCICD(p, s, a, now)
			if err != test.want || err != nil && !reflect.DeepEqual(got, CICDAuthorization{}) {
				t.Fatal("security.phase19-access.boundary.result")
			}
		})
	}
}

func TestCICDAuthorizationGrantChanges(t *testing.T) {
	for _, test := range []struct {
		id, role, project, environment string
		valid                          bool
	}{
		{"viewer", "viewer", "p", "", true}, {"operator", "operator", "p", "", true}, {"approver", "approver", "p", "prod", true}, {"project-admin", "project-admin", "p", "", false}, {"global-admin", "platform-admin", "", "", false}, {"other-project", "viewer", "q", "", false},
	} {
		t.Run(test.id, func(t *testing.T) {
			p, s, now := cicdAuthorizationFixture()
			s.Grants = []CICDGrant{{SubjectKind: "user", SubjectID: "alice", ProjectID: "p", Role: "project-admin", Enabled: true}}
			a := CICDAccess{Operation: "grantPut", ProjectID: test.project, EnvironmentID: test.environment, ProposedGrant: &CICDGrant{SubjectKind: "user", SubjectID: "bob", Role: test.role, ProjectID: test.project, EnvironmentID: test.environment, Enabled: true}}
			_, err := AuthorizeCICD(p, s, a, now)
			if (err == nil) != test.valid {
				t.Fatal("security.phase19-access.grant-change.result")
			}
		})
	}
}

func TestCICDAuthorizationProjection(t *testing.T) {
	p, s, now := cicdAuthorizationFixture()
	expiry := now.Add(time.Minute)
	s.Grants = []CICDGrant{{SubjectKind: "user", SubjectID: "alice", ProjectID: "p", Role: "viewer", Enabled: true, ExpiresAt: &expiry}, {SubjectKind: "user", SubjectID: "bob", ProjectID: "q", Role: "viewer", Enabled: true}}
	got, err := AuthorizeCICD(p, s, CICDAccess{Operation: "runGet", ProjectID: "p"}, now)
	if err != nil || len(got.Grants) != 1 || got.Grants[0].SubjectID != "alice" {
		t.Fatal("security.phase19-access.projection.scope")
	}
	got.Grants[0].SubjectID = "tampered"
	*got.Grants[0].ExpiresAt = now
	if s.Grants[0].SubjectID != "alice" || !expiry.Equal(now.Add(time.Minute)) {
		t.Fatal("security.phase19-access.projection.alias")
	}
}

func TestOwnerFileContract(t *testing.T) {
	if Owner() != "security" {
		t.Fatalf("unexpected owner: %s", Owner())
	}
	if !executeOwnerFileContract(newOwnerFileContract()) {
		t.Fatalf("security owner file contract failed")
	}
}

func TestPBKDF2HMACSHA256Vectors(t *testing.T) {
	cases := []struct {
		name       string
		password   string
		wantHash   string
		iterations int
	}{
		{
			name:       "short-1",
			password:   "password",
			wantHash:   "120fb6cffcf8b32c43e7225256c4f837a86548c92ccc35480805987cb70be17b",
			iterations: 1,
		},
		{
			name:       "short-2",
			password:   "password",
			wantHash:   "ae4d0c95af6b46d32d0adff928f06dd02a303f8ef3c251dfd6e2d85a95474c43",
			iterations: 2,
		},
		{
			name:       "production",
			password:   "password",
			wantHash:   "669cfe52482116fda1aa2cbe409b2f56c8e4563752b7a28f6eaab614ee005178",
			iterations: CredentialIterations,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hash := pbkdf2HMACSHA256([]byte(tc.password), []byte("salt"), tc.iterations, CredentialTagLength)
			if got := toLowerHex(hash); got != tc.wantHash {
				t.Fatalf("hash mismatch: got %s want %s", got, tc.wantHash)
			}
		})
	}
}

func TestPasswordRecordContract(t *testing.T) {
	record, err := NewPasswordRecord("password")
	if err != nil {
		t.Fatalf("NewPasswordRecord failed: %v", err)
	}
	if err := ValidatePasswordRecord(record); err != nil {
		t.Fatalf("record validation failed: %v", err)
	}
	ok, err := VerifyPassword(record, "password")
	if err != nil || !ok {
		t.Fatalf("password verification failed: ok=%v err=%v", ok, err)
	}
	ok, err = VerifyPassword(record, "different-password")
	if err != nil {
		t.Fatalf("unexpected verification error: %v", err)
	}
	if ok {
		t.Fatalf("different password verified")
	}
}

func toLowerHex(data []byte) string {
	const digits = "0123456789abcdef"
	out := make([]byte, len(data)*2)
	for i, b := range data {
		out[i*2] = digits[b>>4]
		out[i*2+1] = digits[b&0x0f]
	}
	return string(out)
}
