#!/bin/sh
set -u

script_name="adlaire-ci-vps-pull-bootstrap.sh"
bootstrap_runtime="/bin/sh"
source_identifier_format="absolute-https-url-or-absolute-file-path"
source_channel_choices="integration-head|stable-release"
topology_role_choices="ci-cd|site|single-node"
sha256_format="64-lowercase-hex"
secret_input_policy="no-command-argument-secret"
state_model="staging-verify-commit-rollback"
cli_contract="fixed-token-flag-interface"
stdout_contract="single-json-object-secret-safe"
stderr_contract="diagnostic-json-lines-on-failure"
exit_code_contract="common-cli-exit-code-contract"
bootstrap_session_id="phase18.session.vps-pull-bootstrap"
runner_service_unit_path="/etc/systemd/system/adlaire-ci.service"
runner_timer_unit_path="/etc/systemd/system/adlaire-ci.timer"
api_service_unit_path="/etc/systemd/system/adlaire-ci-api.service"
source_channel=""
topology_role=""
source_identifier=""
target_sha256=""
install_dir=""
bin_dir=""
state_dir=""
service_user=""
source_kind=""
stage_dir=""
backup_dir=""
api_token=""
admin_password=""
temporary_password=""
fresh_credentials="0"
commit_started="0"
rollback_action="not_required"
cleanup_result="not_started"

cleanup_stage() {
	if [ -n "$stage_dir" ] && [ -d "$stage_dir" ]; then
		rm -rf "$stage_dir" >/dev/null 2>&1 || cleanup_result="failed"
	fi
	if [ "$cleanup_result" = "not_started" ]; then
		cleanup_result="passed"
	fi
}

backup_path() {
	printf '%s/%s' "$backup_dir" "$1"
}

backup_file() {
	backup_source="$1"
	backup_name="$2"
	backup_mode="$3"
	if [ -L "$backup_source" ]; then
		return 1
	fi
	if [ -e "$backup_source" ]; then
		if [ ! -f "$backup_source" ]; then
			return 1
		fi
		install -m "$backup_mode" "$backup_source" "$(backup_path "$backup_name")" >/dev/null 2>&1 || return 1
	else
		printf 'missing\n' >"$(backup_path "$backup_name.missing")" || return 1
	fi
	return 0
}

restore_file() {
	restore_target="$1"
	restore_name="$2"
	restore_mode="$3"
	if [ -f "$(backup_path "$restore_name")" ]; then
		install -m "$restore_mode" "$(backup_path "$restore_name")" "$restore_target.phase18-restore" >/dev/null 2>&1 || return 1
		mv "$restore_target.phase18-restore" "$restore_target" >/dev/null 2>&1 || return 1
	elif [ -f "$(backup_path "$restore_name.missing")" ]; then
		rm -f "$restore_target" "$restore_target.phase18-restore" >/dev/null 2>&1 || return 1
	fi
	return 0
}

capture_service_state() {
	capture_unit="$1"
	if systemctl is-enabled "$capture_unit" >/dev/null 2>&1; then
		printf 'enabled\n' >"$(backup_path "$capture_unit.enabled")" || return 1
	else
		printf 'disabled\n' >"$(backup_path "$capture_unit.enabled")" || return 1
	fi
	if systemctl is-active "$capture_unit" >/dev/null 2>&1; then
		printf 'active\n' >"$(backup_path "$capture_unit.active")" || return 1
	else
		printf 'inactive\n' >"$(backup_path "$capture_unit.active")" || return 1
	fi
	return 0
}

file_first_line() {
	file_first_line_value=""
	IFS= read -r file_first_line_value <"$1" || true
	printf '%s' "$file_first_line_value"
}

file_has_text() {
	file_has_text_path="$1"
	file_has_text_want="$2"
	if [ ! -f "$file_has_text_path" ] || [ -L "$file_has_text_path" ]; then
		return 1
	fi
	while IFS= read -r file_has_text_line || [ -n "$file_has_text_line" ]; do
		case "$file_has_text_line" in
		*"$file_has_text_want"*) return 0 ;;
		esac
	done <"$file_has_text_path"
	return 1
}

file_sha256() {
	file_sha256_result=$(sha256sum "$1" 2>/dev/null) || return 1
	printf '%s' "${file_sha256_result%% *}"
}

restore_service_state() {
	restore_unit="$1"
	if [ "$(file_first_line "$(backup_path "$restore_unit.enabled")")" = "enabled" ]; then
		systemctl enable "$restore_unit" >/dev/null 2>&1 || return 1
	else
		systemctl disable "$restore_unit" >/dev/null 2>&1 || true
	fi
	if [ "$(file_first_line "$(backup_path "$restore_unit.active")")" = "active" ]; then
		systemctl start "$restore_unit" >/dev/null 2>&1 || return 1
	else
		systemctl stop "$restore_unit" >/dev/null 2>&1 || true
	fi
	return 0
}

rollback_changes() {
	rollback_ok="1"
	systemctl stop adlaire-ci-api.service adlaire-ci.timer adlaire-ci.service >/dev/null 2>&1 || true
	for rollback_binary in adlaire-ci-build adlaire-ci-runner adlaire-ci-api adlaire-ci-setup adlaire-ci-admin adlaire-ci-mcp adlaire-ci-obsidian; do
		rm -f "$bin_dir/$rollback_binary.phase18-tmp" "$bin_dir/$rollback_binary.phase18-restore" >/dev/null 2>&1 || rollback_ok="0"
		restore_file "$bin_dir/$rollback_binary" "$rollback_binary" 0755 || rollback_ok="0"
	done
	for rollback_unit_path in "$runner_service_unit_path" "$runner_timer_unit_path" "$api_service_unit_path"; do
		rm -f "$rollback_unit_path.phase18-tmp" "$rollback_unit_path.phase18-restore" >/dev/null 2>&1 || rollback_ok="0"
	done
	restore_file "$runner_service_unit_path" "adlaire-ci.service" 0644 || rollback_ok="0"
	restore_file "$runner_timer_unit_path" "adlaire-ci.timer" 0644 || rollback_ok="0"
	restore_file "$api_service_unit_path" "adlaire-ci-api.service" 0644 || rollback_ok="0"
	for rollback_state_name in .admin_credentials .server_config .pipeline_config .branch_config .last_sha .pending_transfers .notify_pending; do
		restore_file "$state_dir/$rollback_state_name" "state-$rollback_state_name" 0600 || rollback_ok="0"
	done
	restore_file "$install_dir/admin/index.html" "admin-index.html" 0644 || rollback_ok="0"
	restore_file "$install_dir/admin/adlaire-ci-sdk.js" "admin-sdk.js" 0644 || rollback_ok="0"
	if [ -f "$(backup_path "phase18-source-dir.missing")" ]; then
		rm -rf "$state_dir/phase18-source" >/dev/null 2>&1 || rollback_ok="0"
	fi
	if [ -f "$(backup_path "phase18-site-dir.missing")" ]; then
		rm -rf "$state_dir/phase18-site" >/dev/null 2>&1 || rollback_ok="0"
	fi
	if [ -f "$(backup_path "phase18-evidence-dir.missing")" ]; then
		rm -rf "$state_dir/.phase18-evidence" >/dev/null 2>&1 || rollback_ok="0"
	fi
	systemctl daemon-reload >/dev/null 2>&1 || rollback_ok="0"
	for rollback_unit in adlaire-ci.service adlaire-ci.timer adlaire-ci-api.service; do
		restore_service_state "$rollback_unit" || rollback_ok="0"
	done
	if [ "$rollback_ok" = "1" ]; then
		return 0
	fi
	return 1
}

json_fail() {
	failure_code="$1"
	failure_stage="$2"
	failure_exit_code="$3"
	api_token=""
	admin_password=""
	temporary_password=""
	if [ "$commit_started" = "1" ]; then
		if rollback_changes; then
			rollback_action="restored"
			cleanup_result="passed"
		else
			failure_code="bootstrap-script-rollback-failed"
			rollback_action="restore_failed"
			cleanup_result="failed"
			failure_exit_code="1"
		fi
	fi
	cleanup_stage
	printf '{"status":"failed","failure_code":"%s","stage":"%s","source_channel":"%s","topology_role":"%s","source_identifier_kind":"%s","target_digest_sha256":"%s","rollback_action":"%s","cleanup_result":"%s"}\n' \
		"$failure_code" "$failure_stage" "$source_channel" "$topology_role" "$source_kind" "$target_sha256" "$rollback_action" "$cleanup_result" >&2
	exit "$failure_exit_code"
}

require_command() {
	command -v "$1" >/dev/null 2>&1 || json_fail "bootstrap-script-systemd-unit-invalid" "required-command" 4
}

require_exact_args() {
	if [ "$#" -ne 16 ]; then
		json_fail "bootstrap-script-source-unreachable" "parse" 2
	fi
	if [ "$1" != "--source-channel" ] || [ "$3" != "--topology-role" ] || [ "$5" != "--source" ] || [ "$7" != "--sha256" ] || [ "$9" != "--install-dir" ] || [ "${11}" != "--bin-dir" ] || [ "${13}" != "--state-dir" ] || [ "${15}" != "--service-user" ]; then
		json_fail "bootstrap-script-source-unreachable" "parse" 2
	fi
	source_channel="$2"
	topology_role="$4"
	source_identifier="$6"
	target_sha256="$8"
	install_dir="${10}"
	bin_dir="${12}"
	state_dir="${14}"
	service_user="${16}"
}

require_fixed_inputs() {
	case "$source_channel" in
	integration-head | stable-release) ;;
	*) json_fail "bootstrap-script-source-unreachable" "source-channel" 2 ;;
	esac
	case "$topology_role" in
	ci-cd | site | single-node) ;;
	*) json_fail "bootstrap-script-source-unreachable" "topology-role" 2 ;;
	esac
	case "$source_identifier" in
	https://*) source_kind="absolute-https-url" ;;
	/*) source_kind="absolute-file-path" ;;
	*) json_fail "bootstrap-script-source-unreachable" "source-identifier" 2 ;;
	esac
	if [ "${#target_sha256}" -ne 64 ]; then
		json_fail "bootstrap-script-digest-mismatch" "digest-format" 2
	fi
	case "$target_sha256" in
	*[!0123456789abcdef]*) json_fail "bootstrap-script-digest-mismatch" "digest-format" 2 ;;
	esac
	if [ "$install_dir" != "/opt/adlaire-builder" ] || [ "$bin_dir" != "/usr/local/bin" ] || [ "$state_dir" != "/opt/adlaire-builder" ] || [ "$service_user" != "root" ]; then
		json_fail "bootstrap-script-state-directory-invalid" "fixed-path" 2
	fi
	for fixed_path in "$install_dir" "$bin_dir" "$state_dir"; do
		if [ -L "$fixed_path" ]; then
			json_fail "bootstrap-script-state-directory-invalid" "symlink-boundary" 1
		fi
	done
	for protected_path in "$install_dir/admin" "$state_dir/phase18-source" "$state_dir/phase18-site" "$state_dir/.phase18-evidence"; do
		if [ -L "$protected_path" ]; then
			json_fail "bootstrap-script-state-directory-invalid" "symlink-boundary" 1
		fi
	done
}

read_admin_password() {
	if [ -t 0 ]; then
		json_fail "bootstrap-script-health-failed" "admin-password-stdin" 2
	fi
	if ! IFS= read -r admin_password; then
		json_fail "bootstrap-script-health-failed" "admin-password-stdin" 2
	fi
	if IFS= read -r extra_password_line; then
		admin_password=""
		json_fail "bootstrap-script-health-failed" "admin-password-stdin" 2
	fi
	if [ "${#admin_password}" -lt 12 ] || [ "${#admin_password}" -gt 128 ]; then
		admin_password=""
		json_fail "bootstrap-script-health-failed" "admin-password-policy" 2
	fi
	case "$admin_password" in
	*[!A-Za-z0-9._%+,:=@]*)
		admin_password=""
		json_fail "bootstrap-script-health-failed" "admin-password-policy" 2
		;;
	esac
	temporary_password="P18Tmp_"
	set -- $(od -An -N32 -tx1 /dev/urandom 2>/dev/null) || {
		admin_password=""
		temporary_password=""
		json_fail "bootstrap-script-health-failed" "temporary-password-entropy" 1
	}
	for entropy_word in "$@"; do
		temporary_password="$temporary_password$entropy_word"
	done
	if [ "${#temporary_password}" -ne 71 ]; then
		admin_password=""
		temporary_password=""
		json_fail "bootstrap-script-health-failed" "temporary-password-entropy" 1
	fi
	if [ "$admin_password" = "$temporary_password" ]; then
		admin_password=""
		temporary_password=""
		json_fail "bootstrap-script-health-failed" "admin-password-policy" 2
	fi
}

download_source() {
	source_payload="$stage_dir/adlaire-ci-source"
	if [ "$source_kind" = "absolute-file-path" ]; then
		if [ ! -f "$source_identifier" ] || [ -L "$source_identifier" ]; then
			json_fail "bootstrap-script-source-unreachable" "source-copy" 3
		fi
		install -m 0755 "$source_identifier" "$source_payload" >/dev/null 2>&1 || json_fail "bootstrap-script-source-unreachable" "source-copy" 3
	else
		curl -fsSL --proto '=https' --tlsv1.2 --max-time 120 -o "$source_payload" "$source_identifier" >/dev/null 2>&1 || json_fail "bootstrap-script-source-unreachable" "source-download" 3
		chmod 0755 "$source_payload" >/dev/null 2>&1 || json_fail "bootstrap-script-executable-permission-invalid" "source-permission" 1
	fi
	actual_sha256=$(file_sha256 "$source_payload") || json_fail "bootstrap-script-digest-mismatch" "digest" 1
	if [ "$actual_sha256" != "$target_sha256" ]; then
		json_fail "bootstrap-script-digest-mismatch" "digest" 1
	fi
}

verify_binary_version() {
	verify_name="$1"
	verify_path="$stage_dir/bin/$verify_name"
	install -m 0755 "$stage_dir/adlaire-ci-source" "$verify_path" >/dev/null 2>&1 || json_fail "bootstrap-script-executable-permission-invalid" "binary-stage" 1
	version_line=$("$verify_path" --version 2>/dev/null) || json_fail "bootstrap-script-version-mismatch" "binary-version" 1
	case "$version_line" in
	"$verify_name "*)
		case "$source_channel:$version_line" in
		stable-release:*" V.0.0-dev "*) json_fail "bootstrap-script-version-mismatch" "binary-version" 1 ;;
		esac
		;;
	*) json_fail "bootstrap-script-version-mismatch" "binary-version" 1 ;;
	esac
}

write_runner_units() {
	runner_unit="$stage_dir/units/adlaire-ci.service"
	runner_timer="$stage_dir/units/adlaire-ci.timer"
	api_unit="$stage_dir/units/adlaire-ci-api.service"
	printf '[Unit]\nDescription=Adlaire CI Runner\nAfter=network.target\n\n[Service]\nType=oneshot\nUser=%s\nWorkingDirectory=%s\nExecStart=%s --state-dir %s --once\n' \
		"$service_user" "$install_dir" "$bin_dir/adlaire-ci-runner" "$state_dir" >"$runner_unit" || json_fail "bootstrap-script-systemd-unit-invalid" "unit-render" 1
	printf '[Unit]\nDescription=Adlaire CI Runner Timer\n\n[Timer]\nOnBootSec=1min\nOnUnitActiveSec=5min\nUnit=adlaire-ci.service\n\n[Install]\nWantedBy=timers.target\n' >"$runner_timer" || json_fail "bootstrap-script-systemd-unit-invalid" "unit-render" 1
	printf '[Unit]\nDescription=Adlaire CI API Server\nAfter=network.target\n\n[Service]\nType=simple\nUser=%s\nWorkingDirectory=%s\nExecStart=%s --addr 127.0.0.1:8765 --state-dir %s\nRestart=always\nRestartSec=10\nKillSignal=SIGTERM\nTimeoutStopSec=15s\n\n[Install]\nWantedBy=multi-user.target\n' \
		"$service_user" "$install_dir" "$bin_dir/adlaire-ci-api" "$state_dir" >"$api_unit" || json_fail "bootstrap-script-systemd-unit-invalid" "unit-render" 1
	for unit_path in "$runner_unit" "$runner_timer" "$api_unit"; do
		unit_first_line=$(file_first_line "$unit_path")
		if [ "$unit_first_line" != "[Unit]" ]; then
			json_fail "bootstrap-script-systemd-unit-invalid" "unit-validate" 1
		fi
	done
}

write_staged_state() {
	printf '{"extra_args":[],"env":{}}\n' >"$stage_dir/state/.pipeline_config" || json_fail "bootstrap-script-state-directory-invalid" "state-stage" 1
	printf '{"watch_mode":"local","build_cooldown_seconds":0,"snapshots_keep":2,"commit_status_enabled":false,"schedule_paused":false}\n' >"$stage_dir/state/.server_config" || json_fail "bootstrap-script-state-directory-invalid" "state-stage" 1
	printf '{"branch_targets":[{"branch":"main","target_file":"phase18","target_files":[],"sha_file":"%s/.last_sha","src":"%s/phase18-source","out":"%s/phase18-site","approval_required":false,"env":{},"deploy_targets":[]}]}\n' \
		"$state_dir" "$state_dir" "$state_dir" >"$stage_dir/state/.branch_config" || json_fail "bootstrap-script-state-directory-invalid" "state-stage" 1
	printf '{"sha":"phase18-bootstrap-initial"}\n' >"$stage_dir/state/.last_sha" || json_fail "bootstrap-script-state-directory-invalid" "state-stage" 1
	printf '[]\n' >"$stage_dir/state/.pending_transfers" || json_fail "bootstrap-script-state-directory-invalid" "state-stage" 1
	printf '[]\n' >"$stage_dir/state/.notify_pending" || json_fail "bootstrap-script-state-directory-invalid" "state-stage" 1
	printf '# Phase 18 VPS Pull Bootstrap\n\nThis page is generated by the live local runner validation flow.\n' >"$stage_dir/state/source.md" || json_fail "bootstrap-script-state-directory-invalid" "state-stage" 1
	for staged_state_file in .pipeline_config .server_config .branch_config .last_sha .pending_transfers .notify_pending; do
		chmod 0600 "$stage_dir/state/$staged_state_file" >/dev/null 2>&1 || json_fail "bootstrap-script-state-directory-invalid" "state-mode" 1
	done
	chmod 0644 "$stage_dir/state/source.md" >/dev/null 2>&1 || json_fail "bootstrap-script-state-directory-invalid" "state-mode" 1
}

stage_admin_assets() {
	if ! "$stage_dir/bin/adlaire-ci-setup" phase18-materialize-admin --install-dir "$stage_dir/install" >"$stage_dir/setup.stdout" 2>"$stage_dir/setup.stderr"; then
		json_fail "bootstrap-script-health-failed" "admin-assets-stage" 1
	fi
	if [ "$(file_first_line "$stage_dir/setup.stdout")" != "setup: success phase18-materialize-admin" ] || [ -s "$stage_dir/setup.stderr" ]; then
		json_fail "bootstrap-script-health-failed" "admin-assets-stage" 1
	fi
	if [ ! -f "$stage_dir/install/admin/index.html" ] || [ ! -f "$stage_dir/install/admin/adlaire-ci-sdk.js" ]; then
		json_fail "bootstrap-script-health-failed" "admin-assets-stage" 1
	fi
}

stage_credentials() {
	if [ -f "$state_dir/.admin_credentials" ]; then
		if [ -L "$state_dir/.admin_credentials" ]; then
			json_fail "bootstrap-script-state-directory-invalid" "credentials-symlink" 1
		fi
		fresh_credentials="0"
		return
	fi
	fresh_credentials="1"
	if ! printf '%s\n' "$temporary_password" | "$stage_dir/bin/adlaire-ci-api" --state-dir "$stage_dir/state" --init-credentials >"$stage_dir/credentials.stdout" 2>"$stage_dir/credentials.stderr"; then
		json_fail "bootstrap-script-health-failed" "credentials-stage" 1
	fi
	if [ "$(file_first_line "$stage_dir/credentials.stdout")" != "credentials initialized" ] || [ -s "$stage_dir/credentials.stderr" ] || [ ! -f "$stage_dir/state/.admin_credentials" ]; then
		json_fail "bootstrap-script-health-failed" "credentials-stage" 1
	fi
	chmod 0600 "$stage_dir/state/.admin_credentials" >/dev/null 2>&1 || json_fail "bootstrap-script-state-directory-invalid" "credentials-mode" 1
}

prepare_backups() {
	backup_dir="$stage_dir/backup"
	mkdir -p "$backup_dir" >/dev/null 2>&1 || json_fail "bootstrap-script-rollback-failed" "backup" 1
	for backup_binary in adlaire-ci-build adlaire-ci-runner adlaire-ci-api adlaire-ci-setup adlaire-ci-admin adlaire-ci-mcp adlaire-ci-obsidian; do
		backup_file "$bin_dir/$backup_binary" "$backup_binary" 0755 || json_fail "bootstrap-script-rollback-failed" "backup" 1
	done
	backup_file "$runner_service_unit_path" "adlaire-ci.service" 0644 || json_fail "bootstrap-script-rollback-failed" "backup" 1
	backup_file "$runner_timer_unit_path" "adlaire-ci.timer" 0644 || json_fail "bootstrap-script-rollback-failed" "backup" 1
	backup_file "$api_service_unit_path" "adlaire-ci-api.service" 0644 || json_fail "bootstrap-script-rollback-failed" "backup" 1
	for backup_unit in adlaire-ci.service adlaire-ci.timer adlaire-ci-api.service; do
		capture_service_state "$backup_unit" || json_fail "bootstrap-script-rollback-failed" "backup-service-state" 1
	done
	for backup_state_name in .admin_credentials .server_config .pipeline_config .branch_config .last_sha .pending_transfers .notify_pending; do
		backup_file "$state_dir/$backup_state_name" "state-$backup_state_name" 0600 || json_fail "bootstrap-script-rollback-failed" "backup-state" 1
	done
	backup_file "$install_dir/admin/index.html" "admin-index.html" 0644 || json_fail "bootstrap-script-rollback-failed" "backup-admin" 1
	backup_file "$install_dir/admin/adlaire-ci-sdk.js" "admin-sdk.js" 0644 || json_fail "bootstrap-script-rollback-failed" "backup-admin" 1
	for backup_dir_name in phase18-source phase18-site .phase18-evidence; do
		if [ -L "$state_dir/$backup_dir_name" ]; then
			json_fail "bootstrap-script-state-directory-invalid" "backup-directory-symlink" 1
		fi
		if [ -d "$state_dir/$backup_dir_name" ]; then
			printf 'exists\n' >"$(backup_path "$backup_dir_name-dir.exists")" || json_fail "bootstrap-script-rollback-failed" "backup-directory" 1
		elif [ -e "$state_dir/$backup_dir_name" ]; then
			json_fail "bootstrap-script-state-directory-invalid" "backup-directory-type" 1
		else
			printf 'missing\n' >"$(backup_path "$backup_dir_name-dir.missing")" || json_fail "bootstrap-script-rollback-failed" "backup-directory" 1
		fi
	done
}

atomic_install() {
	atomic_source="$1"
	atomic_target="$2"
	atomic_mode="$3"
	if [ -L "$atomic_target" ]; then
		json_fail "bootstrap-script-executable-permission-invalid" "atomic-install-symlink" 1
	fi
	install -m "$atomic_mode" "$atomic_source" "$atomic_target.phase18-tmp" >/dev/null 2>&1 || json_fail "bootstrap-script-executable-permission-invalid" "atomic-install" 1
	mv "$atomic_target.phase18-tmp" "$atomic_target" >/dev/null 2>&1 || json_fail "bootstrap-script-executable-permission-invalid" "atomic-commit" 1
}

commit_install() {
	prepare_backups
	mkdir -p "$install_dir" "$bin_dir" "$state_dir" "$install_dir/admin" "$state_dir/.build_logs" "$state_dir/.snapshots" "$state_dir/phase18-source" "$state_dir/.phase18-evidence" >/dev/null 2>&1 || json_fail "bootstrap-script-state-directory-invalid" "directory-create" 1
	chmod 0700 "$state_dir" >/dev/null 2>&1 || json_fail "bootstrap-script-state-directory-invalid" "state-mode" 1
	chmod 0755 "$install_dir/admin" "$state_dir/phase18-source" >/dev/null 2>&1 || json_fail "bootstrap-script-state-directory-invalid" "directory-mode" 1
	chmod 0700 "$state_dir/.build_logs" "$state_dir/.snapshots" "$state_dir/.phase18-evidence" >/dev/null 2>&1 || json_fail "bootstrap-script-state-directory-invalid" "directory-mode" 1
	commit_started="1"
	systemctl stop adlaire-ci.timer adlaire-ci.service >/dev/null 2>&1 || true
	for install_binary in adlaire-ci-build adlaire-ci-runner adlaire-ci-api adlaire-ci-setup adlaire-ci-admin adlaire-ci-mcp adlaire-ci-obsidian; do
		atomic_install "$stage_dir/bin/$install_binary" "$bin_dir/$install_binary" 0755
		installed_version=$("$bin_dir/$install_binary" --version 2>/dev/null) || json_fail "bootstrap-script-version-mismatch" "binary-installed-version" 1
		case "$installed_version" in
		"$install_binary "*) ;;
		*) json_fail "bootstrap-script-version-mismatch" "binary-installed-version" 1 ;;
		esac
	done
	atomic_install "$stage_dir/install/admin/index.html" "$install_dir/admin/index.html" 0644
	atomic_install "$stage_dir/install/admin/adlaire-ci-sdk.js" "$install_dir/admin/adlaire-ci-sdk.js" 0644
	for install_state_name in .server_config .pipeline_config .branch_config .last_sha .pending_transfers .notify_pending; do
		atomic_install "$stage_dir/state/$install_state_name" "$state_dir/$install_state_name" 0600
	done
	if [ "$fresh_credentials" = "1" ]; then
		atomic_install "$stage_dir/state/.admin_credentials" "$state_dir/.admin_credentials" 0600
	fi
	atomic_install "$stage_dir/state/source.md" "$state_dir/phase18-source/index.md" 0644
	atomic_install "$stage_dir/units/adlaire-ci.service" "$runner_service_unit_path" 0644
	atomic_install "$stage_dir/units/adlaire-ci.timer" "$runner_timer_unit_path" 0644
	atomic_install "$stage_dir/units/adlaire-ci-api.service" "$api_service_unit_path" 0644
	systemctl daemon-reload >/dev/null 2>&1 || json_fail "bootstrap-script-systemd-unit-invalid" "daemon-reload" 1
	for verify_unit in adlaire-ci.service adlaire-ci.timer adlaire-ci-api.service; do
		systemctl cat "$verify_unit" >/dev/null 2>&1 || json_fail "bootstrap-script-systemd-unit-invalid" "unit-health" 1
	done
	systemctl enable adlaire-ci-api.service adlaire-ci.timer >/dev/null 2>&1 || json_fail "bootstrap-script-health-failed" "service-enable" 1
	systemctl restart adlaire-ci-api.service >/dev/null 2>&1 || json_fail "bootstrap-script-health-failed" "api-service-start" 1
}

curl_request() {
	curl_method="$1"
	curl_path="$2"
	curl_body="$3"
	curl_auth="$4"
	curl_output="$5"
	curl_config="$stage_dir/curl.config"
	{
		printf 'silent\n'
		printf 'show-error\n'
		printf 'fail-with-body\n'
		printf 'max-time = 30\n'
		printf 'request = "%s"\n' "$curl_method"
		printf 'url = "http://127.0.0.1:8765%s"\n' "$curl_path"
		printf 'header = "Accept: application/json"\n'
		if [ "$curl_auth" = "auth" ]; then
			printf 'header = "Authorization: Bearer %s"\n' "$api_token"
		fi
		if [ "$curl_body" != "none" ]; then
			printf 'header = "Content-Type: application/json"\n'
			printf 'data-binary = "@%s"\n' "$curl_body"
		fi
		printf 'output = "%s"\n' "$curl_output"
	} >"$curl_config" || return 1
	chmod 0600 "$curl_config" >/dev/null 2>&1 || return 1
	curl --config "$curl_config" >/dev/null 2>&1
}

wait_for_api() {
	wait_attempt="0"
	while [ "$wait_attempt" -lt 30 ]; do
		if curl_request GET /api/health none noauth "$stage_dir/health.json"; then
			if file_has_text "$stage_dir/health.json" '"status"'; then
				return 0
			fi
		fi
		wait_attempt=$((wait_attempt + 1))
		sleep 1
	done
	return 1
}

extract_api_token() {
	token_json="$1"
	case "$token_json" in
	*'"token":"'*)
		token_rest=${token_json#*\"token\":\"}
		api_token=${token_rest%%\"*}
		;;
	*) api_token="" ;;
	esac
	case "$api_token" in
	"" | *'
'*) return 1 ;;
	esac
	return 0
}

api_login() {
	login_password="$1"
	printf '{"password":"%s"}\n' "$login_password" >"$stage_dir/login.json" || return 1
	chmod 0600 "$stage_dir/login.json" >/dev/null 2>&1 || return 1
	curl_request POST /api/login "$stage_dir/login.json" noauth "$stage_dir/login-response.json" || return 1
	login_response=$(file_first_line "$stage_dir/login-response.json")
	extract_api_token "$login_response" || return 1
	printf '%s\n' "$api_token" >"$stage_dir/admin-token" || return 1
	chmod 0600 "$stage_dir/admin-token" >/dev/null 2>&1 || return 1
	return 0
}

change_initial_password() {
	printf '{"current_password":"%s","new_password":"%s"}\n' "$temporary_password" "$admin_password" >"$stage_dir/change-password.json" || return 1
	chmod 0600 "$stage_dir/change-password.json" >/dev/null 2>&1 || return 1
	curl_request POST /api/change-password "$stage_dir/change-password.json" auth "$stage_dir/change-password-response.json" || return 1
	file_has_text "$stage_dir/change-password-response.json" '"message":"Password changed"' || return 1
	temporary_password=""
	return 0
}

admin_status_check() {
	if ! "$bin_dir/adlaire-ci-admin" --api-url http://127.0.0.1:8765 --token-file "$stage_dir/admin-token" --json status >"$stage_dir/admin-status.json" 2>"$stage_dir/admin-status.stderr"; then
		return 1
	fi
	if [ -s "$stage_dir/admin-status.stderr" ]; then
		return 1
	fi
	file_has_text "$stage_dir/admin-status.json" '"last_build_status"' || return 1
	file_has_text "$stage_dir/admin-status.json" '"running"' || return 1
	return 0
}

verify_ui_sdk() {
	curl_request GET / none noauth "$stage_dir/ui.html" || return 1
	curl_request GET /admin/adlaire-ci-sdk.js none noauth "$stage_dir/sdk.js" || return 1
	ui_digest=$(file_sha256 "$stage_dir/ui.html") || return 1
	ui_expected_digest=$(file_sha256 "$install_dir/admin/index.html") || return 1
	sdk_digest=$(file_sha256 "$stage_dir/sdk.js") || return 1
	sdk_expected_digest=$(file_sha256 "$install_dir/admin/adlaire-ci-sdk.js") || return 1
	[ "$ui_digest" = "$ui_expected_digest" ] || return 1
	[ "$sdk_digest" = "$sdk_expected_digest" ] || return 1
	return 0
}

write_open_sample() {
	open_state_digest=$(file_sha256 "$state_dir/.server_config") || json_fail "bootstrap-script-health-failed" "open-sample-digest" 1
	printf '{"session_id":"%s","sample":"open","service_state":"active","api_health":"passed","admin_cli_health":"passed","runner_state":"ready","statefile_digest":"%s","log_write":"pending-runtime","secret_boundary":"passed","known_bug_status":"none","document_drift_status":"no_drift"}\n' \
		"$bootstrap_session_id" "$open_state_digest" >"$state_dir/.phase18-evidence/open.json" || json_fail "bootstrap-script-health-failed" "open-sample-write" 1
	chmod 0600 "$state_dir/.phase18-evidence/open.json" >/dev/null 2>&1 || json_fail "bootstrap-script-health-failed" "open-sample-mode" 1
}

exercise_api_configuration() {
	printf '{"seconds":2}\n' >"$stage_dir/cooldown-enable.json" || return 1
	printf '{"seconds":0}\n' >"$stage_dir/cooldown-disable.json" || return 1
	curl_request POST /api/schedule/cooldown "$stage_dir/cooldown-enable.json" auth "$stage_dir/cooldown-enable-response.json" || return 1
	curl_request POST /api/schedule/cooldown "$stage_dir/cooldown-disable.json" auth "$stage_dir/cooldown-disable-response.json" || return 1
	file_has_text "$stage_dir/cooldown-disable-response.json" '"seconds":0' || return 1
	return 0
}

run_local_runner() {
	rm -rf "$state_dir/phase18-site" >/dev/null 2>&1 || return 1
	if ! "$bin_dir/adlaire-ci-runner" --state-dir "$state_dir" --once >"$stage_dir/runner.stdout" 2>"$stage_dir/runner.stderr"; then
		return 1
	fi
	if [ -s "$stage_dir/runner.stderr" ]; then
		return 1
	fi
	if [ ! -f "$state_dir/phase18-site/index.html" ] || [ ! -f "$state_dir/phase18-site/.dependency_manifest.json" ]; then
		return 1
	fi
	file_has_text "$state_dir/.build_history" '"status":"success"' || return 1
	file_has_text "$state_dir/.build_status.json" '"status":"success"' || return 1
	return 0
}

verify_runtime_logs() {
	for runtime_log in .access_log .audit_log .api_access_log .config_log; do
		if [ ! -s "$state_dir/$runtime_log" ] || [ -L "$state_dir/$runtime_log" ]; then
			return 1
		fi
	done
	return 0
}

restart_api_and_wait() {
	systemctl restart adlaire-ci-api.service >/dev/null 2>&1 || return 1
	wait_for_api
}

run_update_rollback_drill() {
	mkdir -p "$stage_dir/drill" >/dev/null 2>&1 || return 1
	credential_digest_before=$(file_sha256 "$state_dir/.admin_credentials") || return 1
	history_digest_before=$(file_sha256 "$state_dir/.build_history") || return 1
	for drill_binary in adlaire-ci-build adlaire-ci-runner adlaire-ci-api adlaire-ci-setup adlaire-ci-admin adlaire-ci-mcp adlaire-ci-obsidian; do
		install -m 0755 "$bin_dir/$drill_binary" "$stage_dir/drill/$drill_binary" >/dev/null 2>&1 || return 1
		install -m 0755 "$stage_dir/adlaire-ci-source" "$bin_dir/$drill_binary.phase18-tmp" >/dev/null 2>&1 || return 1
		mv "$bin_dir/$drill_binary.phase18-tmp" "$bin_dir/$drill_binary" >/dev/null 2>&1 || return 1
		"$bin_dir/$drill_binary" --version >/dev/null 2>&1 || return 1
	done
	restart_api_and_wait || return 1
	for drill_binary in adlaire-ci-build adlaire-ci-runner adlaire-ci-api adlaire-ci-setup adlaire-ci-admin adlaire-ci-mcp adlaire-ci-obsidian; do
		install -m 0755 "$stage_dir/drill/$drill_binary" "$bin_dir/$drill_binary.phase18-tmp" >/dev/null 2>&1 || return 1
		mv "$bin_dir/$drill_binary.phase18-tmp" "$bin_dir/$drill_binary" >/dev/null 2>&1 || return 1
		"$bin_dir/$drill_binary" --version >/dev/null 2>&1 || return 1
	done
	restart_api_and_wait || return 1
	for drill_binary in adlaire-ci-build adlaire-ci-runner adlaire-ci-api adlaire-ci-setup adlaire-ci-admin adlaire-ci-mcp adlaire-ci-obsidian; do
		install -m 0755 "$stage_dir/adlaire-ci-source" "$bin_dir/$drill_binary.phase18-tmp" >/dev/null 2>&1 || return 1
		mv "$bin_dir/$drill_binary.phase18-tmp" "$bin_dir/$drill_binary" >/dev/null 2>&1 || return 1
		installed_digest=$(file_sha256 "$bin_dir/$drill_binary") || return 1
		[ "$installed_digest" = "$target_sha256" ] || return 1
	done
	restart_api_and_wait || return 1
	systemctl restart adlaire-ci.timer >/dev/null 2>&1 || return 1
	credential_digest_after=$(file_sha256 "$state_dir/.admin_credentials") || return 1
	history_digest_after=$(file_sha256 "$state_dir/.build_history") || return 1
	[ "$credential_digest_before" = "$credential_digest_after" ] || return 1
	[ "$history_digest_before" = "$history_digest_after" ] || return 1
	return 0
}

write_runtime_sample() {
	runtime_state_digest=$(file_sha256 "$state_dir/.build_history") || json_fail "bootstrap-script-health-failed" "runtime-sample-digest" 1
	printf '{"session_id":"%s","sample":"runtime","service_state":"active","api_health":"passed","admin_cli_health":"passed","runner_state":"success","statefile_digest":"%s","audit_log_write":"passed","access_log_write":"passed","config_log_write":"passed","update_rollback":"passed","secret_boundary":"passed","known_bug_status":"none","document_drift_status":"no_drift"}\n' \
		"$bootstrap_session_id" "$runtime_state_digest" >"$state_dir/.phase18-evidence/runtime.json" || json_fail "bootstrap-script-health-failed" "runtime-sample-write" 1
	chmod 0600 "$state_dir/.phase18-evidence/runtime.json" >/dev/null 2>&1 || json_fail "bootstrap-script-health-failed" "runtime-sample-mode" 1
}

write_close_sample() {
	close_state_digest=$(file_sha256 "$state_dir/.build_history") || json_fail "bootstrap-script-health-failed" "close-sample-digest" 1
	printf '{"session_id":"%s","sample":"close","service_state":"active","api_health":"passed","admin_cli_health":"passed","runner_state":"success","statefile_digest":"%s","secret_boundary":"passed","destructive_operation_result":"blocked_zero","known_bug_status":"none","document_drift_status":"no_drift","cleanup_result":"passed","final_open_item_result":0}\n' \
		"$bootstrap_session_id" "$close_state_digest" >"$state_dir/.phase18-evidence/close.json" || json_fail "bootstrap-script-health-failed" "close-sample-write" 1
	chmod 0600 "$state_dir/.phase18-evidence/close.json" >/dev/null 2>&1 || json_fail "bootstrap-script-health-failed" "close-sample-mode" 1
}

logout_api() {
	curl_request POST /api/logout none auth "$stage_dir/logout-response.json" || return 1
	file_has_text "$stage_dir/logout-response.json" '"message":"Logged out"' || return 1
	api_token=""
	return 0
}

require_exact_args "$@"
require_fixed_inputs
for command_name in sh mktemp mkdir chmod install mv rm sha256sum systemctl sleep od; do
	require_command "$command_name"
done
command -v curl >/dev/null 2>&1 || json_fail "bootstrap-script-http-client-missing" "http-client" 4
read_admin_password
stage_dir=$(mktemp -d "${TMPDIR:-/tmp}/adlaire-ci-bootstrap.XXXXXX") || json_fail "bootstrap-script-state-directory-invalid" "staging" 1
mkdir -p "$stage_dir/bin" "$stage_dir/units" "$stage_dir/state" "$stage_dir/install" >/dev/null 2>&1 || json_fail "bootstrap-script-state-directory-invalid" "staging" 1
chmod 0700 "$stage_dir" "$stage_dir/state" >/dev/null 2>&1 || json_fail "bootstrap-script-state-directory-invalid" "staging-mode" 1
download_source
for stage_binary in adlaire-ci-build adlaire-ci-runner adlaire-ci-api adlaire-ci-setup adlaire-ci-admin adlaire-ci-mcp adlaire-ci-obsidian; do
	verify_binary_version "$stage_binary"
done
write_runner_units
write_staged_state
stage_admin_assets
stage_credentials
commit_install
wait_for_api || json_fail "bootstrap-script-health-failed" "api-health-open" 1
if [ "$fresh_credentials" = "1" ]; then
	api_login "$temporary_password" || json_fail "bootstrap-script-health-failed" "api-login-initial" 1
	change_initial_password || json_fail "bootstrap-script-health-failed" "api-password-change" 1
else
	api_login "$admin_password" || json_fail "bootstrap-script-health-failed" "api-login-existing" 1
fi
admin_status_check || json_fail "bootstrap-script-health-failed" "admin-cli-open" 1
verify_ui_sdk || json_fail "bootstrap-script-health-failed" "ui-sdk-open" 1
write_open_sample
exercise_api_configuration || json_fail "bootstrap-script-health-failed" "api-config-runtime" 1
run_local_runner || json_fail "bootstrap-script-health-failed" "runner-runtime" 1
admin_status_check || json_fail "bootstrap-script-health-failed" "admin-cli-runtime" 1
verify_ui_sdk || json_fail "bootstrap-script-health-failed" "ui-sdk-runtime" 1
verify_runtime_logs || json_fail "bootstrap-script-health-failed" "runtime-log-write" 1
run_update_rollback_drill || json_fail "bootstrap-script-health-failed" "update-rollback-drill" 1
api_login "$admin_password" || json_fail "bootstrap-script-health-failed" "api-login-rebootstrap" 1
admin_status_check || json_fail "bootstrap-script-health-failed" "admin-cli-revalidate" 1
verify_ui_sdk || json_fail "bootstrap-script-health-failed" "ui-sdk-revalidate" 1
verify_runtime_logs || json_fail "bootstrap-script-health-failed" "runtime-log-revalidate" 1
write_runtime_sample
write_close_sample
logout_api || json_fail "bootstrap-script-health-failed" "api-logout" 1
systemctl is-active adlaire-ci-api.service >/dev/null 2>&1 || json_fail "bootstrap-script-health-failed" "api-service-close" 1
systemctl is-active adlaire-ci.timer >/dev/null 2>&1 || json_fail "bootstrap-script-health-failed" "runner-timer-close" 1
admin_password=""
temporary_password=""
commit_started="0"
cleanup_stage
if [ "$cleanup_result" != "passed" ]; then
	json_fail "bootstrap-script-health-failed" "cleanup" 1
fi
printf '{"status":"ok","evidence_origin":"live-vps-bootstrap","execution_environment":"vps","provider_target":"conoha-vps-primary","vps_label":"試験本番VPS","validation_mode":"vps-pull-bootstrap","bootstrap_session_id":"%s","bootstrap_session_kind":"operator-approved-bootstrap-session","source_channel":"%s","topology_role":"%s","bootstrap_script_artifact":"%s","bootstrap_runtime":"%s","source_identifier_kind":"%s","source_digest_sha256":"%s","install_dir":"%s","bin_dir":"%s","state_dir":"%s","service_user":"%s","operation_samples":["open","runtime","close"],"work_unit_result":"all_completed","api_admin_sdk_ui_result":"passed","pull_runner_foundation_result":"passed","update_rollback_result":"passed","bugfix_loop_result":"no_issue_detected","secret_boundary_result":"passed","destructive_operation_result":"blocked_zero","known_bug_status":"none","document_drift_status":"no_drift","cleanup_result":"%s","final_open_item_result":0}\n' \
	"$bootstrap_session_id" "$source_channel" "$topology_role" "$script_name" "$bootstrap_runtime" "$source_kind" "$target_sha256" "$install_dir" "$bin_dir" "$state_dir" "$service_user" "$cleanup_result"
