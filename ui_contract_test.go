package main

import (
	"strings"
	"testing"
)

func TestUIRequiredDOMContract(t *testing.T) {
	ui := readTextFile(t, "admin/index.html")

	required := []string{
		`id="app-root"`,
		`data-api-base-url=""`,
		`id="nav-panels"`,
		`id="global-banner"`,
		`id="global-error"`,
		`id="global-success"`,
		`id="maintenance-banner"`,
		`id="build-log-stream"`,
		`id="build-queue-summary"`,
		`id="issued-token-once"`,
		`id="totp-secret-once"`,
		`id="totp-otpauth-once"`,
		`name="failure_category"`,
		`name="backup_file"`,
		`id="btn-backup"`,
		`id="btn-restore"`,
		`id="btn-copy-issued-token"`,
		`id="btn-copy-totp-secret"`,
		`id="btn-copy-totp-otpauth"`,
		`id="btn-load-webhook-events"`,
		`id="btn-circuit-reset"`,
		`id="btn-approve-build"`,
		`id="btn-reject-build"`,
		`id="btn-save-history-comment"`,
		`id="btn-save-history-flag"`,
		`id="btn-save-history-tags"`,
		`id="btn-rollback-history"`,
		`id="btn-save-pipeline-config"`,
		`id="btn-load-build-chain"`,
		`id="btn-save-build-chain"`,
		`id="btn-add-alert-rule"`,
		`id="btn-delete-alert-rule"`,
		`id="btn-add-tag-rule"`,
		`id="btn-delete-tag-rule"`,
		`id="btn-set-schedule-interval"`,
		`id="btn-set-allowed-hours"`,
		`id="btn-clear-allowed-hours"`,
		`id="btn-set-force-interval"`,
		`id="btn-set-build-cooldown"`,
		`id="btn-revoke-token"`,
		`id="btn-download-snapshot"`,
		`id="btn-delete-snapshot"`,
		`id="btn-delete-hook"`,
		`id="btn-load-hook-log"`,
		`:focus-visible`,
		`.field-error`,
		`error.className = "field-error";`,
	}
	for _, pattern := range required {
		if !strings.Contains(ui, pattern) {
			t.Fatalf("UI DOM contract missing %q", pattern)
		}
	}
}

func TestUIUsesSDKBoundaryAndBaseURLContract(t *testing.T) {
	ui := readTextFile(t, "admin/index.html")

	required := []string{
		`import { AdlaireCI, AdlaireCIError } from "./adlaire-ci-sdk.js";`,
		`function resolveBaseUrl(root)`,
		`return window.location.origin;`,
		`new URL(raw)`,
		`url.origin !== window.location.origin`,
		`url.username !== ""`,
		`url.password !== ""`,
		`url.pathname !== "/"`,
		`url.search !== ""`,
		`url.hash !== ""`,
		`state.client = new AdlaireCI({ baseUrl });`,
	}
	for _, pattern := range required {
		if !strings.Contains(ui, pattern) {
			t.Fatalf("UI SDK/base URL contract missing %q", pattern)
		}
	}

	forbidden := []string{
		`data-api-base-url="/api"`,
		`fetch(`,
		`XMLHttpRequest`,
		`EventSource`,
		`ReadableStream`,
		`window.AdlaireCI`,
		`window.AdlaireCIError`,
		`localStorage`,
		`sessionStorage`,
	}
	for _, pattern := range forbidden {
		if strings.Contains(ui, pattern) {
			t.Fatalf("UI must not contain %q", pattern)
		}
	}
}

func TestUIOperationContractSnippets(t *testing.T) {
	ui := readTextFile(t, "admin/index.html")

	required := []string{
		`getDashboard();`,
		`getStatus();`,
		`getQueue();`,
		`getMaintenance();`,
		`getDashboardLayout();`,
		`state.streamHandle.done.then`,
		`summary === null`,
		`state.client.getLogs()`,
		`state.client.getOutputMeta()`,
		`state.client.getBuildTrends(100)`,
		`state.client.getWebhookEvents(50, 0)`,
		`state.client.resetCircuitBreaker()`,
		`state.client.approveBuild(trimmedValue("form-approval", "approval_id"))`,
		`state.client.rejectBuild(trimmedValue("form-approval", "approval_id"))`,
		`state.client.setHistoryComment(trimmedValue("form-history-operation", "history_id"), textValue("form-history-operation", "comment"))`,
		`state.client.setHistoryFlag(trimmedValue("form-history-operation", "history_id"), Boolean(formValue("form-history-operation", "flagged")))`,
		`state.client.setHistoryTags(trimmedValue("form-history-operation", "history_id"), listValue("form-history-operation", "tags"))`,
		`state.client.rollbackHistory(id)`,
		`state.client.setNotifyConfig(notifyConfigFromForm())`,
		`state.client.setRepoConfig(repoConfigFromForm())`,
		`state.client.setBranchConfig(branchConfigFromForm())`,
		`state.client.getAccessLog()`,
		`state.client.backup()`,
		`state.client.restore(jsonValue("form-config", "backup_file", {}))`,
		`state.client.setPipelineConfig(jsonValue("form-config", "pipeline_config", {}))`,
		`state.client.getBuildChainConfig()`,
		`state.client.setBuildChainConfig(jsonValue("form-config", "build_chains", []))`,
		`state.client.addAlertRule(trimmedValue("form-config", "alert_metric"), trimmedValue("form-config", "alert_operator"), numberValue("form-config", "alert_threshold"), trimmedValue("form-config", "alert_level"), textValue("form-config", "alert_message"))`,
		`state.client.deleteAlertRule(trimmedValue("form-config", "alert_rule_id"))`,
		`state.client.addTagRule(jsonValue("form-config", "tag_condition", {}), listValue("form-config", "tag_tags"))`,
		`state.client.deleteTagRule(trimmedValue("form-config", "tag_rule_id"))`,
		`state.client.setScheduleInterval(intValue("form-repo", "interval_seconds"))`,
		`state.client.setAllowedHours(intValue("form-repo", "allowed_from"), intValue("form-repo", "allowed_to"))`,
		`state.client.clearAllowedHours()`,
		`state.client.setForceInterval(intValue("form-repo", "force_interval_hours"))`,
		`state.client.setBuildCooldown(intValue("form-repo", "build_cooldown_seconds"))`,
		`state.client.revokeToken(id)`,
		`state.client.downloadSnapshot(id)`,
		`state.client.deleteSnapshot(id)`,
		`state.client.deleteHook(trimmedValue("form-hook", "hook_id"))`,
		`state.client.getHookLog(trimmedValue("form-hook", "hook_id"))`,
		`duration_anomaly = {`,
		`failureCategory: trimmedValue("form-history-filter", "failure_category") || null`,
		`showOneTimeSecrets({ token: result.token })`,
		`totpOtpauth: result.otpauth_uri || ""`,
		`navigator.clipboard.writeText(value)`,
	}
	for _, pattern := range required {
		if !strings.Contains(ui, pattern) {
			t.Fatalf("UI operation contract missing %q", pattern)
		}
	}

	forbidden := []string{
		`name="webhooks"`,
		`name="email"`,
		`webhooks: listValue`,
		`email: listValue`,
		`state.client.setBranchConfig([repoConfigFromForm()])`,
		`state.client.getAccessLog,`,
	}
	for _, pattern := range forbidden {
		if strings.Contains(ui, pattern) {
			t.Fatalf("UI operation contract must not contain %q", pattern)
		}
	}
}

func TestUIErrorContractSnippets(t *testing.T) {
	ui := readTextFile(t, "admin/index.html")

	required := []string{
		`error.status === 401`,
		`error.status === 403`,
		`error.status === 409`,
		`error.status === 422`,
		`error.status === 429`,
		`error.status === 503`,
		`refetchAfterConflict(panelId)`,
		`disableTemporarily(button, 10000)`,
		`showFieldError(panelId, field, message)`,
		`firstInvalid.focus()`,
		`$("field-password").value = "";`,
		`setMessage("maintenance-banner"`,
		`state.rateLimited.set(button.id, until)`,
	}
	for _, pattern := range required {
		if !strings.Contains(ui, pattern) {
			t.Fatalf("UI error contract missing %q", pattern)
		}
	}
}
