package api

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/fqwink/build-scripts/components/security"
	"github.com/fqwink/build-scripts/components/setup"
	"github.com/fqwink/build-scripts/components/statefile"
)

func parseCICDManagementArgs(args []string) (apiCLIConfig, bool, error) {
	cfg := apiCLIConfig{CICDManagement: true}
	reject := func() (apiCLIConfig, bool, error) {
		return apiCLIConfig{}, false, exitError{Code: 2, Msg: "P19_INPUT"}
	}
	for _, arg := range args {
		if !validCLIArgToken(arg) {
			return reject()
		}
	}
	seen := map[string]bool{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if seen[arg] {
			return reject()
		}
		seen[arg] = true
		switch arg {
		case "--cicd-management":
		case "--config", "--state-dir":
			if i+1 >= len(args) || !cicdManagementPath(args[i+1]) {
				return reject()
			}
			i++
			if arg == "--config" {
				cfg.CICDConfigPath = args[i]
			} else {
				cfg.StateDir = args[i]
			}
		default:
			return reject()
		}
	}
	if !seen["--cicd-management"] || !seen["--config"] || !seen["--state-dir"] {
		return reject()
	}
	return cfg, false, nil
}

func cicdManagementPath(path string) bool {
	return filepath.IsAbs(path) && filepath.Clean(path) == path && len(path) <= 4096 && !strings.Contains(path, "\\")
}

func prepareCICDManagement(cfg apiCLIConfig, readConfig func(string) (setup.CICDConfig, error)) (*http.Server, error) {
	config, err := readConfig(cfg.CICDConfigPath)
	if err != nil {
		return nil, err
	}
	if config.Role != "coordinator" || config.StateDir != cfg.StateDir || config.Listen != "127.0.0.1:9321" {
		return nil, errors.New("P19_SCHEMA")
	}
	store, err := statefile.OpenCICDIdentityStore(cfg.StateDir)
	if err != nil {
		return nil, err
	}
	if _, err := os.Lstat(filepath.Join(cfg.StateDir, ".auth_profile")); !errors.Is(err, os.ErrNotExist) {
		return nil, statefile.ErrCICDCorrupt
	}
	if err := cicdManagementLegacyFiles(cfg.StateDir); err != nil {
		return nil, err
	}
	legacy, err := NewAPIServer(APIConfig{Addr: config.Listen, StateDir: cfg.StateDir})
	if err != nil {
		return nil, statefile.ErrCICDCorrupt
	}
	if _, err := legacy.readMergedConfig(); err != nil {
		return nil, statefile.ErrCICDCorrupt
	}
	totp, err := readTOTPSecret(filepath.Join(cfg.StateDir, ".totp_secret"))
	if err != nil || validateTOTPSecretValue(totp) != nil {
		return nil, statefile.ErrCICDCorrupt
	}
	handler, err := newCICDAuthHTTP(legacy, store)
	if err != nil {
		return nil, err
	}
	server := legacy.HTTPServer()
	server.Handler = cicdManagementHandler(handler)
	return server, nil
}

func cicdManagementLegacyFiles(path string) error {
	dir, err := os.OpenFile(path, os.O_RDONLY|syscall.O_DIRECTORY|syscall.O_NOFOLLOW|syscall.O_CLOEXEC, 0)
	if err != nil {
		return statefile.ErrCICDCorrupt
	}
	defer dir.Close()
	for _, name := range []string{".admin_credentials", ".totp_secret", ".server_config"} {
		file, err := os.OpenFile(fmt.Sprintf("/proc/self/fd/%d/%s", dir.Fd(), name), os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_CLOEXEC|syscall.O_NONBLOCK, 0)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) && name != ".admin_credentials" {
				continue
			}
			return statefile.ErrCICDCorrupt
		}
		info, err := file.Stat()
		if err != nil {
			_ = file.Close()
			return statefile.ErrCICDCorrupt
		}
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok || !info.Mode().IsRegular() || info.Mode() != 0600 || stat.Uid != uint32(os.Geteuid()) || stat.Nlink != 1 || info.Size() > 1048576 {
			_ = file.Close()
			return statefile.ErrCICDCorrupt
		}
		data, readErr := io.ReadAll(io.LimitReader(file, 1048577))
		closeErr := file.Close()
		if readErr != nil || closeErr != nil || int64(len(data)) != info.Size() || len(data) > 1048576 {
			return statefile.ErrCICDCorrupt
		}
		if _, err := statefile.CanonicalCICDJSON(data); err != nil {
			return statefile.ErrCICDCorrupt
		}
		kind := ""
		if name == ".admin_credentials" {
			kind = "Credential"
		}
		if name == ".totp_secret" {
			kind = "TOTP"
		}
		if kind != "" && statefile.ValidateCICDValue(kind, data) != nil {
			return statefile.ErrCICDCorrupt
		}
	}
	return nil
}

func cicdManagementHandler(handler *cicdAuthHTTP) http.Handler {
	legacy := handler.legacy.Handler()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/cicd" || strings.HasPrefix(r.URL.Path, "/api/cicd/") {
			handler.ServeHTTP(w, r)
			return
		}
		legacy.ServeHTTP(w, r)
	})
}

func executeCICDManagement(cfg apiCLIConfig, stderr io.Writer) int {
	server, err := prepareCICDManagement(cfg, setup.ReadCICDConfig)
	if err != nil {
		code := err.Error()
		switch code {
		case "P19_SCHEMA", "P19_CAPABILITY", "P19_BUSY", "P19_STATE_CORRUPT":
		default:
			code = "P19_STATE_CORRUPT"
		}
		fmt.Fprintln(stderr, code)
		return 3
	}
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(signals)
	return listenCICDManagement(server, signals, stderr)
}

func listenCICDManagement(server *http.Server, signals <-chan os.Signal, stderr io.Writer) int {
	listener, err := net.Listen("tcp4", server.Addr)
	if err != nil {
		fmt.Fprintln(stderr, "P19_STORAGE")
		return 4
	}
	return serveCICDManagement(server, listener, signals, stderr)
}

func serveCICDManagement(server *http.Server, listener net.Listener, signals <-chan os.Signal, stderr io.Writer) int {
	errCh := make(chan error, 1)
	go func() { errCh <- server.Serve(listener) }()
	select {
	case <-errCh:
		_ = server.Close()
		fmt.Fprintln(stderr, "P19_STORAGE")
		return 4
	case <-signals:
	}
	done := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		done <- server.Shutdown(ctx)
	}()
	var shutdownErr error
	select {
	case shutdownErr = <-done:
	case <-signals:
		_ = server.Close()
		<-done
		shutdownErr = errors.New("P19_SHUTDOWN")
	}
	serveErr := <-errCh
	if shutdownErr != nil || !errors.Is(serveErr, http.ErrServerClosed) {
		_ = server.Close()
		fmt.Fprintln(stderr, "P19_SHUTDOWN")
		return 4
	}
	return 0
}

func newCICDAuthHTTP(legacy *APIServer, store cicdIdentityBackend) (*cicdAuthHTTP, error) {
	if legacy == nil || store == nil {
		return nil, statefile.ErrCICDCorrupt
	}
	return &cicdAuthHTTP{legacy: legacy, store: store, gate: make(chan struct{}, 1), entropy: rand.Reader, sessions: make(map[string]cicdAuthSession), revoked: make(map[string]time.Time), windows: make(map[string]cicdAuthWindow), failures: make(map[string]cicdAuthFailure)}, nil
}

func (h *cicdAuthHTTP) random(n int) (string, error) {
	h.entropyMu.Lock()
	defer h.entropyMu.Unlock()
	data := make([]byte, n)
	if _, err := io.ReadFull(h.entropy, data); err != nil {
		return "", errors.New("P19_ENTROPY")
	}
	return hex.EncodeToString(data), nil
}

func (h *cicdAuthHTTP) enter(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if ctx.Err() != nil {
		return statefile.ErrCICDBusy
	}
	select {
	case h.gate <- struct{}{}:
		if ctx.Err() != nil {
			<-h.gate
			return statefile.ErrCICDBusy
		}
		return nil
	case <-ctx.Done():
		return statefile.ErrCICDBusy
	}
}

func (h *cicdAuthHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := h.random(16)
	if err != nil {
		cicdAuthWrite(w, nil, nil, nil, err)
		return
	}
	id = "o-" + id
	r = r.WithContext(context.WithValue(r.Context(), apiRequestIDContextKey{}, id))
	if r.URL.Path != "/api/cicd/login" && r.URL.Path != "/api/cicd/logout" {
		cicdAuthWrite(w, id, nil, nil, errors.New("P19_NOT_FOUND"))
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		cicdAuthWrite(w, id, nil, nil, errors.New("P19_METHOD"))
		return
	}
	if r.URL.RawQuery != "" || r.URL.ForceQuery || r.URL.RawPath != "" {
		cicdAuthWrite(w, id, nil, nil, errors.New("P19_INPUT"))
		return
	}
	if err := h.enter(r.Context()); err != nil {
		cicdAuthWrite(w, id, nil, nil, err)
		return
	}
	defer func() { <-h.gate }()
	now := h.legacy.cfg.Now().UTC().Truncate(time.Second)
	if now.Year() < 1970 || now.Year() > 9999 {
		cicdAuthWrite(w, id, nil, nil, statefile.ErrCICDCorrupt)
		return
	}
	h.expire(now)
	if r.URL.Path == "/api/cicd/login" {
		h.login(w, r, id, now)
		return
	}
	token := cicdAuthBearer(r)
	authorization, err := h.authorizeLocked(token, security.CICDAccess{Operation: "logout"}, now)
	if err != nil {
		h.authError(w, r, id, authorization, err)
		return
	}
	if _, err := cicdAuthBody(r, false); err != nil {
		cicdAuthWrite(w, id, nil, &authorization, err)
		return
	}
	if authorization.SubjectKind == "service" {
		if len(h.revoked) >= 4096 {
			cicdAuthWrite(w, id, nil, &authorization, security.ErrCICDQuota)
			return
		}
		h.revoked[tokenHash(token)] = authorization.ExpiresAt
	} else {
		delete(h.sessions, tokenHash(token))
	}
	if err := h.logAuth(r, "logout", authorization.SubjectID, true); err != nil {
		cicdAuthWrite(w, id, nil, nil, err)
		return
	}
	cicdAuthWrite(w, id, map[string]any{}, &authorization, nil)
}

func cicdAuthIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		if ip := net.ParseIP(host); ip != nil {
			return ip.String()
		}
	}
	return "unknown"
}

func (h *cicdAuthHTTP) expire(now time.Time) {
	for key, session := range h.sessions {
		if !now.Before(session.principal.ExpiresAt) {
			delete(h.sessions, key)
		}
	}
	for key, until := range h.revoked {
		if !now.Before(until) {
			delete(h.revoked, key)
		}
	}
	for key, window := range h.windows {
		if !now.Before(window.start.Add(60 * time.Second)) {
			delete(h.windows, key)
		}
	}
	for key, failure := range h.failures {
		until := failure.until
		if until.IsZero() {
			until = failure.at.Add(600 * time.Second)
		}
		if !now.Before(until) {
			delete(h.failures, key)
		}
	}
}

func (h *cicdAuthHTTP) rate(key string, now time.Time) error {
	window, exists := h.windows[key]
	if exists && (now.Before(window.start) || window.count >= 10) || !exists && len(h.windows) >= 4096 {
		return errors.New("P19_RATE")
	}
	if !exists {
		window.start = now
	}
	window.count++
	h.windows[key] = window
	return nil
}

func (h *cicdAuthHTTP) login(w http.ResponseWriter, r *http.Request, id string, now time.Time) {
	ip := cicdAuthIP(r)
	if err := h.rate("ip:"+ip, now); err != nil {
		h.authError(w, r, id, security.CICDAuthorization{}, err)
		return
	}
	body, err := cicdAuthBody(r, true)
	if err != nil {
		cicdAuthWrite(w, id, nil, nil, err)
		return
	}
	userID, password := body["user_id"].(string), body["password"].(string)
	if err := h.rate("user:"+userID, now); err != nil {
		h.authError(w, r, id, security.CICDAuthorization{}, err)
		return
	}
	failure, exists := h.failures[ip]
	if exists && (now.Before(failure.at) || !failure.until.IsZero() && now.Before(failure.until)) || !exists && len(h.failures) >= 4096 {
		h.authError(w, r, id, security.CICDAuthorization{}, errors.New("P19_RATE"))
		return
	}
	if len(h.sessions) >= 4096 {
		cicdAuthWrite(w, id, nil, nil, security.ErrCICDQuota)
		return
	}
	token := ""
	for attempt := 0; attempt < 3; attempt++ {
		token, err = h.random(32)
		if err != nil {
			break
		}
		if _, used := h.sessions[tokenHash(token)]; !used {
			break
		}
		token = ""
	}
	if err != nil || token == "" {
		cicdAuthWrite(w, id, nil, nil, errors.New("P19_ENTROPY"))
		return
	}
	config, err := h.legacy.readMergedConfig()
	if err != nil {
		cicdAuthWrite(w, id, nil, nil, statefile.ErrCICDCorrupt)
		return
	}
	expires := now.Add(time.Duration(config.SessionTimeoutSeconds) * time.Second)
	if expires.Year() > 9999 {
		cicdAuthWrite(w, id, nil, nil, statefile.ErrCICDCorrupt)
		return
	}
	if userID != "admin" && body["totp"] != nil {
		h.authError(w, r, id, security.CICDAuthorization{}, security.ErrCICDAuth)
		return
	}
	adminBinding := ""
	result, err := h.store.Login(userID, password, id, now, func() (string, bool, error) {
		return h.adminLogin(r.Context(), password, body["totp"], now, &adminBinding)
	})
	if result.PasswordVerified {
		delete(h.failures, ip)
	} else if err == security.ErrCICDAuth {
		failure.count++
		failure.at = now
		if failure.count >= 10 {
			failure.until = now.Add(600 * time.Second)
		}
		h.failures[ip] = failure
	}
	if err != nil {
		h.authError(w, r, id, security.CICDAuthorization{}, err)
		return
	}
	principal := security.CICDPrincipal{Kind: "user", ID: userID, ExpiresAt: expires, MustChange: result.MustChange}
	session := cicdAuthSession{principal: principal, binding: cicdAuthIdentityBinding(result.Identity)}
	if userID == "admin" {
		session.adminBinding = adminBinding
	}
	if err := h.logAuth(r, "login_success", userID, true); err != nil {
		cicdAuthWrite(w, id, nil, nil, err)
		return
	}
	authorization, err := h.store.Authorize(principal, security.CICDAccess{Operation: "logout"}, now, func(identity statefile.CICDIdentity) error { return h.checkSession(session, identity) })
	if err != nil {
		cicdAuthWrite(w, id, nil, nil, err)
		return
	}
	h.sessions[tokenHash(token)] = session
	cicdAuthWrite(w, id, map[string]any{"token": token, "user_id": userID, "expires_at": expires.Format(time.RFC3339), "must_change": result.MustChange}, &authorization, nil)
}

func cicdAuthIdentityBinding(identity statefile.CICDIdentity) string {
	fingerprint := ""
	if identity.Credential != nil {
		fingerprint = security.CICDCredentialFingerprint(*identity.Credential)
	}
	data, _ := json.Marshal([]any{identity.Kind, identity.ID, identity.Revision, fingerprint})
	return tokenHash(string(data))
}

func (h *cicdAuthHTTP) checkSession(session cicdAuthSession, identity statefile.CICDIdentity) error {
	if !tokenHashEqual(session.binding, cicdAuthIdentityBinding(identity)) {
		return security.ErrCICDAuth
	}
	if identity.ID == "admin" {
		current, err := h.adminBinding()
		if err != nil {
			return err
		}
		if !tokenHashEqual(session.adminBinding, current) {
			return security.ErrCICDAuth
		}
	}
	return nil
}

func (h *cicdAuthHTTP) authorizeLocked(token string, access security.CICDAccess, now time.Time) (security.CICDAuthorization, error) {
	var principal security.CICDPrincipal
	var check func(statefile.CICDIdentity) error
	switch security.CICDBearerKind(token) {
	case "user":
		session, exists := h.sessions[tokenHash(token)]
		if !exists || !now.Before(session.principal.ExpiresAt) {
			delete(h.sessions, tokenHash(token))
			return security.CICDAuthorization{}, security.ErrCICDAuth
		}
		principal = session.principal
		check = func(identity statefile.CICDIdentity) error { return h.checkSession(session, identity) }
	case "service":
		if _, revoked := h.revoked[tokenHash(token)]; revoked {
			return security.CICDAuthorization{}, security.ErrCICDAuth
		}
		identity, err := h.store.ReadService(tokenHash(token))
		if err != nil {
			return security.CICDAuthorization{}, err
		}
		principal = security.CICDPrincipal{Kind: "service", ID: identity.ID, ExpiresAt: identity.ExpiresAt, MustChange: "none"}
		check = func(current statefile.CICDIdentity) error {
			if !tokenHashEqual(current.TokenSHA256, tokenHash(token)) {
				return security.ErrCICDAuth
			}
			return nil
		}
	default:
		return security.CICDAuthorization{}, security.ErrCICDAuth
	}
	authorization, err := h.store.Authorize(principal, access, now, check)
	if err == security.ErrCICDForbidden {
		current, authErr := h.store.Authorize(principal, security.CICDAccess{Operation: "logout"}, now, check)
		if authErr != nil {
			return security.CICDAuthorization{}, authErr
		}
		return current, err
	}
	if err == security.ErrCICDAuth && principal.Kind == "user" {
		delete(h.sessions, tokenHash(token))
	}
	return authorization, err
}

func cicdAuthBearer(r *http.Request) string {
	values := r.Header.Values("Authorization")
	if len(values) != 1 || !strings.HasPrefix(values[0], "Bearer ") {
		return ""
	}
	token := strings.TrimPrefix(values[0], "Bearer ")
	if security.CICDBearerKind(token) == "" {
		return ""
	}
	return token
}

func cicdAuthBody(r *http.Request, login bool) (map[string]any, error) {
	media, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || media != "application/json" || len(r.Header.Values("Content-Type")) != 1 || len(params) > 1 || len(params) == 1 && strings.ToLower(params["charset"]) != "utf-8" || len(r.Header.Values("Content-Encoding")) != 0 {
		return nil, errors.New("P19_MEDIA")
	}
	data, err := io.ReadAll(io.LimitReader(r.Body, 1048577))
	if err != nil {
		return nil, errors.New("P19_JSON")
	}
	if len(data) > 1048576 {
		return nil, errors.New("P19_BODY")
	}
	decoder, depth := json.NewDecoder(bytes.NewReader(data)), 0
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.New("P19_JSON")
		}
		if delim, ok := token.(json.Delim); ok {
			if delim == '{' || delim == '[' {
				depth++
			} else {
				depth--
			}
			if depth > 32 {
				return nil, errors.New("P19_JSON")
			}
		}
	}
	canonical, err := statefile.CanonicalCICDJSON(data)
	if err != nil {
		return nil, errors.New("P19_JSON")
	}
	var object map[string]any
	if json.Unmarshal(canonical, &object) != nil || object == nil {
		return nil, errors.New("P19_SCHEMA")
	}
	if !login {
		if len(object) != 0 {
			return nil, errors.New("P19_SCHEMA")
		}
		return object, nil
	}
	if len(object) != 3 {
		return nil, errors.New("P19_SCHEMA")
	}
	id, idOK := object["user_id"].(string)
	password, passwordOK := object["password"].(string)
	value, totpOK := object["totp"]
	idJSON, _ := json.Marshal(id)
	if !idOK || statefile.ValidateCICDValue("Id", idJSON) != nil || !passwordOK || len(password) < 1 || len(password) > 4096 || !totpOK {
		return nil, errors.New("P19_SCHEMA")
	}
	if value != nil {
		text, ok := value.(string)
		if !ok || len(text) != 6 {
			return nil, errors.New("P19_SCHEMA")
		}
		for _, c := range text {
			if c < '0' || c > '9' {
				return nil, errors.New("P19_SCHEMA")
			}
		}
	}
	return object, nil
}

func (h *cicdAuthHTTP) adminLogin(ctx context.Context, password string, code any, now time.Time, binding *string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	select {
	case h.legacy.authTransactions <- struct{}{}:
	case <-ctx.Done():
		return "", false, statefile.ErrCICDBusy
	}
	defer func() { <-h.legacy.authTransactions }()
	path := filepath.Join(h.legacy.cfg.StateDir, ".admin_credentials")
	credential, err := readCredentialsLocked(path)
	if err != nil {
		return "", false, statefile.ErrCICDCorrupt
	}
	if !credentialsPasswordEqual(credential, password) {
		return "", false, security.ErrCICDAuth
	}
	totpPath := filepath.Join(h.legacy.cfg.StateDir, ".totp_secret")
	totp, err := readTOTPSecretLocked(totpPath)
	if err != nil {
		return "", true, statefile.ErrCICDCorrupt
	}
	if !totp.Enabled && code != nil || totp.Enabled && code == nil {
		return "", true, security.ErrCICDAuth
	}
	if totp.Enabled {
		step, ok := verifyTOTPCodeStep(*totp.SecretBase32, code.(string), now)
		if !ok || totp.LastAcceptedStep != nil && step <= *totp.LastAcceptedStep {
			return "", true, security.ErrCICDAuth
		}
		_, err = updateTOTPSecretLocked(totpPath, func(current apiTOTPSecret) (apiTOTPSecret, bool, error) {
			if !current.Enabled || current.SecretBase32 == nil || !passwordHashEqual(*current.SecretBase32, *totp.SecretBase32) || current.LastAcceptedStep != nil && step <= *current.LastAcceptedStep {
				return current, false, errAPITOTPInvalid
			}
			current.LastAcceptedStep = &step
			return current, true, nil
		})
		if errors.Is(err, errAPITOTPInvalid) {
			return "", true, security.ErrCICDAuth
		}
		if err != nil {
			return "", true, statefile.ErrCICDStorage
		}
	}
	updated, err := updateCredentialsLocked(path, func(current apiCredentials) (apiCredentials, bool, error) {
		if !passwordHashEqual(credentialsFingerprint(current), credentialsFingerprint(credential)) {
			return current, false, errAPIPasswordMismatch
		}
		if current.LastLoginAt != nil && now.Format(time.RFC3339) < *current.LastLoginAt {
			return current, false, statefile.ErrCICDCorrupt
		}
		advanceCredentialsLogin(&current, now.Format(time.RFC3339))
		return current, true, nil
	})
	if errors.Is(err, errAPIPasswordMismatch) {
		return "", true, security.ErrCICDAuth
	}
	if err != nil {
		return "", true, statefile.ErrCICDStorage
	}
	*binding = cicdAdminBinding(updated, totp)
	return credentialsMustChange(&updated), true, nil
}

func (h *cicdAuthHTTP) adminBinding() (string, error) {
	credential, err := readCredentialsLocked(filepath.Join(h.legacy.cfg.StateDir, ".admin_credentials"))
	if err != nil {
		return "", statefile.ErrCICDCorrupt
	}
	totp, err := readTOTPSecretLocked(filepath.Join(h.legacy.cfg.StateDir, ".totp_secret"))
	if err != nil {
		return "", statefile.ErrCICDCorrupt
	}
	return cicdAdminBinding(credential, totp), nil
}

func cicdAdminBinding(credential apiCredentials, totp apiTOTPSecret) string {
	data, _ := json.Marshal([]any{credentialsFingerprint(credential), credential.MustChange, totp.Enabled, totp.SecretBase32, totp.ConfirmedAt})
	return tokenHash(string(data))
}

func (h *cicdAuthHTTP) logAuth(r *http.Request, action, subject string, success bool) error {
	result := "failure"
	if success {
		result = "success"
	}
	if h.legacy.appendAuthAccessLog(r, action, result, "") != nil || h.legacy.appendSecurityAudit(r, action, "system", subject, "auth", "cicd", result) != nil {
		return statefile.ErrCICDStorage
	}
	return nil
}

func (h *cicdAuthHTTP) authError(w http.ResponseWriter, r *http.Request, id string, auth security.CICDAuthorization, err error) {
	if err == security.ErrCICDAuth || err == security.ErrCICDForbidden || err.Error() == "P19_RATE" {
		action, subject := "permission_denied", "system"
		if err == security.ErrCICDAuth {
			action = "login_failure"
		}
		if auth.SubjectID != "" {
			subject = auth.SubjectID
		}
		if logErr := h.logAuth(r, action, subject, false); logErr != nil {
			err = logErr
		}
	}
	if auth.SubjectID == "" {
		cicdAuthWrite(w, id, nil, nil, err)
	} else {
		cicdAuthWrite(w, id, nil, &auth, err)
	}
}

func cicdAuthWrite(w http.ResponseWriter, id any, data any, auth *security.CICDAuthorization, failure error) {
	var authorization any
	if auth != nil {
		grants := make([]any, 0, len(auth.Grants))
		for _, grant := range auth.Grants {
			var project, environment any
			if grant.ProjectID != "" {
				project = grant.ProjectID
			}
			if grant.EnvironmentID != "" {
				environment = grant.EnvironmentID
			}
			grants = append(grants, map[string]any{"subject_kind": grant.SubjectKind, "subject_id": grant.SubjectID, "project_id": project, "environment_id": environment, "role": grant.Role, "enabled": grant.Enabled, "expires_at": grant.ExpiresAt})
		}
		authorization = map[string]any{"subject_kind": auth.SubjectKind, "subject_id": auth.SubjectID, "generation": auth.Generation, "expires_at": auth.ExpiresAt, "grants": grants}
	}
	status := http.StatusOK
	response := map[string]any{"data": data, "request_id": id, "authorization": authorization}
	if failure != nil {
		code, message := failure.Error(), "operation failed"
		switch code {
		case "P19_JSON", "P19_INPUT":
			status, message = 400, "invalid request"
		case "P19_AUTH":
			status, message = 401, "authentication required"
		case "P19_FORBIDDEN":
			status, message = 403, "permission denied"
		case "P19_NOT_FOUND":
			status, message = 404, "not found"
		case "P19_METHOD":
			status, message = 405, "method not allowed"
		case "P19_BUSY", "P19_REVISION":
			status, message = 409, "state conflict"
		case "P19_BODY":
			status, message = 413, "request too large"
		case "P19_MEDIA":
			status, message = 415, "unsupported media type"
		case "P19_SCHEMA":
			status, message = 422, "validation failed"
		case "P19_RATE", "P19_QUOTA":
			status, message = 429, "limit exceeded"
		case "P19_STATE_CORRUPT", "P19_STATE_UNCERTAIN":
			status, message = 503, "unavailable"
		case "P19_STORAGE":
			status, message = 507, "storage unavailable"
		case "P19_ENTROPY":
			status = 500
		default:
			status, code = 500, "P19_INTERNAL"
		}
		response = map[string]any{"error": message, "code": code, "request_id": id, "details": []any{}, "authorization": authorization}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

var runnerSleep = time.Sleep

type exitError struct {
	Code int
	Msg  string
}

func (e exitError) Error() string { return e.Msg }

func validCLIArgToken(arg string) bool {
	if !utf8.ValidString(arg) {
		return false
	}
	for _, r := range arg {
		if r == 0 || r == '\n' || r == '\r' || r < 0x20 || r == 0x7f {
			return false
		}
	}
	return true
}

func acquireStateFileLock(path string) (func(), error) {
	now := time.Now()
	sleep := func(duration time.Duration) {
		runnerSleep(duration)
		now = now.Add(duration)
	}
	lock, err := statefile.AcquireLock(path+".lock", statefile.LockOptions{
		Owner:          "api-statefile",
		StaleAfter:     30 * time.Second,
		AcquireTimeout: 10 * time.Second,
		RetryInterval:  100 * time.Millisecond,
		Clock:          func() time.Time { return now },
		Sleep:          sleep,
	})
	if err != nil {
		return func() {}, err
	}
	return func() { _ = lock.Release() }, nil
}

func executeOwnerFileContract(contract ownerFileContract) bool {
	return validateOwnerFileContract(contract)
}
