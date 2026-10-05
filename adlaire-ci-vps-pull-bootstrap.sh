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
commit_started="0"
rollback_action="not_required"
cleanup_result="not_started"

json_fail() {
	code="$1"
	stage="$2"
	exit_code="$3"
	if [ "$commit_started" = "1" ]; then
		if rollback_changes; then
			rollback_action="restored"
			cleanup_result="passed"
		else
			code="bootstrap-script-rollback-failed"
			rollback_action="restore_failed"
			cleanup_result="failed"
			exit_code="1"
		fi
	fi
	cleanup_stage
	printf '{"status":"failed","failure_code":"%s","stage":"%s","source_channel":"%s","topology_role":"%s","source_identifier_kind":"%s","target_digest_sha256":"%s","rollback_action":"%s","cleanup_result":"%s"}\n' \
		"$code" "$stage" "$source_channel" "$topology_role" "$source_kind" "$target_sha256" "$rollback_action" "$cleanup_result" >&2
	exit "$exit_code"
}

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
	source_path="$1"
	backup_name="$2"
	if [ -L "$source_path" ]; then
		return 1
	fi
	if [ -e "$source_path" ]; then
		install -m 0644 "$source_path" "$(backup_path "$backup_name")" >/dev/null 2>&1 || return 1
	else
		printf 'missing\n' >"$(backup_path "$backup_name.missing")" || return 1
	fi
	return 0
}

restore_file() {
	target_path="$1"
	backup_name="$2"
	mode="$3"
	if [ -f "$(backup_path "$backup_name")" ]; then
		install -m "$mode" "$(backup_path "$backup_name")" "$target_path" >/dev/null 2>&1 || return 1
	elif [ -f "$(backup_path "$backup_name.missing")" ]; then
		rm -f "$target_path" >/dev/null 2>&1 || return 1
	fi
	return 0
}

rollback_changes() {
	rollback_ok="1"
	for binary in adlaire-ci-build adlaire-ci-runner adlaire-ci-api adlaire-ci-setup adlaire-ci-admin adlaire-ci-mcp adlaire-ci-obsidian; do
		rm -f "$bin_dir/$binary.phase18-tmp" >/dev/null 2>&1 || rollback_ok="0"
		restore_file "$bin_dir/$binary" "$binary" 0755 || rollback_ok="0"
	done
	for unit in adlaire-ci.service adlaire-ci.timer adlaire-ci-api.service; do
		rm -f "/etc/systemd/system/$unit.phase18-tmp" >/dev/null 2>&1 || rollback_ok="0"
	done
	restore_file "/etc/systemd/system/adlaire-ci.service" "adlaire-ci.service" 0644 || rollback_ok="0"
	restore_file "/etc/systemd/system/adlaire-ci.timer" "adlaire-ci.timer" 0644 || rollback_ok="0"
	restore_file "/etc/systemd/system/adlaire-ci-api.service" "adlaire-ci-api.service" 0644 || rollback_ok="0"
	systemctl daemon-reload >/dev/null 2>&1 || rollback_ok="0"
	if [ "$rollback_ok" = "1" ]; then
		return 0
	fi
	return 1
}

require_command() {
	command -v "$1" >/dev/null 2>&1 || json_fail "bootstrap-script-systemd-unit-invalid" "required-command" 4
}

require_exact_args() {
	if [ "$#" -ne 16 ]; then
		json_fail "bootstrap-script-source-unreachable" "parse" 2
	fi
	if [ "$1" != "--source-channel" ] || [ "$3" != "--topology-role" ] || [ "$5" != "--source" ] || [ "$7" != "--sha256" ] || [ "$9" != "--install-dir" ] || [ "$11" != "--bin-dir" ] || [ "$13" != "--state-dir" ] || [ "$15" != "--service-user" ]; then
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
	for path in "$install_dir" "$bin_dir" "$state_dir"; do
		if [ -L "$path" ]; then
			json_fail "bootstrap-script-state-directory-invalid" "symlink-boundary" 1
		fi
	done
}

download_source() {
	payload="$stage_dir/adlaire-ci-source"
	if [ "$source_kind" = "absolute-file-path" ]; then
		if [ ! -f "$source_identifier" ] || [ -L "$source_identifier" ]; then
			json_fail "bootstrap-script-source-unreachable" "source-copy" 3
		fi
		install -m 0755 "$source_identifier" "$payload" >/dev/null 2>&1 || json_fail "bootstrap-script-source-unreachable" "source-copy" 3
	else
		if command -v curl >/dev/null 2>&1; then
			curl -fsSL --proto '=https' --tlsv1.2 -o "$payload" "$source_identifier" >/dev/null 2>&1 || json_fail "bootstrap-script-source-unreachable" "source-download" 3
		elif command -v wget >/dev/null 2>&1; then
			wget -q -O "$payload" "$source_identifier" >/dev/null 2>&1 || json_fail "bootstrap-script-source-unreachable" "source-download" 3
		else
			json_fail "bootstrap-script-http-client-missing" "http-client" 3
		fi
		chmod 0755 "$payload" >/dev/null 2>&1 || json_fail "bootstrap-script-executable-permission-invalid" "source-permission" 1
	fi
	actual_sha256=$(sha256sum "$payload" 2>/dev/null) || json_fail "bootstrap-script-digest-mismatch" "digest" 1
	actual_sha256=${actual_sha256%% *}
	if [ "$actual_sha256" != "$target_sha256" ]; then
		json_fail "bootstrap-script-digest-mismatch" "digest" 1
	fi
}

verify_binary_version() {
	verify_name="$1"
	verify_path="$stage_dir/$verify_name"
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
	runner_unit="$stage_dir/adlaire-ci.service"
	runner_timer="$stage_dir/adlaire-ci.timer"
	api_unit="$stage_dir/adlaire-ci-api.service"
	printf '[Unit]\nDescription=Adlaire CI Runner\n\n[Service]\nType=oneshot\nUser=%s\nWorkingDirectory=%s\nExecStart=%s --state-dir %s\n' \
		"$service_user" "$install_dir" "$bin_dir/adlaire-ci-runner" "$state_dir" >"$runner_unit" || json_fail "bootstrap-script-systemd-unit-invalid" "unit-render" 1
	printf '[Unit]\nDescription=Adlaire CI Runner Timer\n\n[Timer]\nOnBootSec=1min\nOnUnitActiveSec=5min\nUnit=adlaire-ci.service\n\n[Install]\nWantedBy=timers.target\n' >"$runner_timer" || json_fail "bootstrap-script-systemd-unit-invalid" "unit-render" 1
	printf '[Unit]\nDescription=Adlaire CI API Server\nWants=adlaire-ci.timer\nAfter=network.target adlaire-ci.timer\n\n[Service]\nType=simple\nUser=%s\nWorkingDirectory=%s\nExecStart=%s --addr 127.0.0.1:8765 --state-dir %s\nRestart=always\nRestartSec=10\nKillSignal=SIGTERM\nTimeoutStopSec=15s\n\n[Install]\nWantedBy=multi-user.target\n' \
		"$service_user" "$install_dir" "$bin_dir/adlaire-ci-api" "$state_dir" >"$api_unit" || json_fail "bootstrap-script-systemd-unit-invalid" "unit-render" 1
}

commit_install() {
	mkdir -p "$install_dir" "$bin_dir" "$state_dir" "$state_dir/.build_logs" "$state_dir/.snapshots" >/dev/null 2>&1 || json_fail "bootstrap-script-state-directory-invalid" "directory-create" 1
	chmod 0700 "$state_dir" >/dev/null 2>&1 || json_fail "bootstrap-script-state-directory-invalid" "state-mode" 1
	backup_dir="$stage_dir/backup"
	mkdir -p "$backup_dir" >/dev/null 2>&1 || json_fail "bootstrap-script-rollback-failed" "backup" 1
	for binary in adlaire-ci-build adlaire-ci-runner adlaire-ci-api adlaire-ci-setup adlaire-ci-admin adlaire-ci-mcp adlaire-ci-obsidian; do
		if ! backup_file "$bin_dir/$binary" "$binary"; then
			json_fail "bootstrap-script-rollback-failed" "backup" 1
		fi
	done
	for unit in adlaire-ci.service adlaire-ci.timer adlaire-ci-api.service; do
		if ! backup_file "/etc/systemd/system/$unit" "$unit"; then
			json_fail "bootstrap-script-rollback-failed" "backup" 1
		fi
	done
	commit_started="1"
	for binary in adlaire-ci-build adlaire-ci-runner adlaire-ci-api adlaire-ci-setup adlaire-ci-admin adlaire-ci-mcp adlaire-ci-obsidian; do
		if [ -L "$bin_dir/$binary" ]; then
			json_fail "bootstrap-script-executable-permission-invalid" "binary-symlink" 1
		fi
		install -m 0755 "$stage_dir/adlaire-ci-source" "$bin_dir/$binary.phase18-tmp" >/dev/null 2>&1 || json_fail "bootstrap-script-executable-permission-invalid" "binary-install" 1
		mv "$bin_dir/$binary.phase18-tmp" "$bin_dir/$binary" >/dev/null 2>&1 || json_fail "bootstrap-script-executable-permission-invalid" "binary-commit" 1
		installed_version=$("$bin_dir/$binary" --version 2>/dev/null) || json_fail "bootstrap-script-version-mismatch" "binary-installed-version" 1
		case "$installed_version" in
		"$binary "*) ;;
		*) json_fail "bootstrap-script-version-mismatch" "binary-installed-version" 1 ;;
		esac
	done
	for unit in adlaire-ci.service adlaire-ci.timer adlaire-ci-api.service; do
		install -m 0644 "$stage_dir/$unit" "/etc/systemd/system/$unit.phase18-tmp" >/dev/null 2>&1 || json_fail "bootstrap-script-systemd-unit-invalid" "unit-install" 1
		mv "/etc/systemd/system/$unit.phase18-tmp" "/etc/systemd/system/$unit" >/dev/null 2>&1 || json_fail "bootstrap-script-systemd-unit-invalid" "unit-commit" 1
	done
	printf '{"event":"bootstrap","kind":"audit"}\n' >>"$state_dir/.phase18-bootstrap-audit.jsonl" || json_fail "bootstrap-script-health-failed" "log-write" 1
	printf '{"event":"bootstrap","kind":"access"}\n' >>"$state_dir/.phase18-bootstrap-access.jsonl" || json_fail "bootstrap-script-health-failed" "log-write" 1
	printf '{"event":"bootstrap","kind":"config"}\n' >>"$state_dir/.phase18-bootstrap-config.jsonl" || json_fail "bootstrap-script-health-failed" "log-write" 1
	chmod 0600 "$state_dir/.phase18-bootstrap-audit.jsonl" "$state_dir/.phase18-bootstrap-access.jsonl" "$state_dir/.phase18-bootstrap-config.jsonl" >/dev/null 2>&1 || json_fail "bootstrap-script-health-failed" "log-mode" 1
	systemctl daemon-reload >/dev/null 2>&1 || json_fail "bootstrap-script-systemd-unit-invalid" "daemon-reload" 1
	systemctl enable --now adlaire-ci.timer >/dev/null 2>&1 || json_fail "bootstrap-script-health-failed" "runner-timer" 1
	systemctl cat adlaire-ci.service >/dev/null 2>&1 || json_fail "bootstrap-script-systemd-unit-invalid" "runner-unit-health" 1
	systemctl cat adlaire-ci.timer >/dev/null 2>&1 || json_fail "bootstrap-script-systemd-unit-invalid" "timer-unit-health" 1
	systemctl cat adlaire-ci-api.service >/dev/null 2>&1 || json_fail "bootstrap-script-systemd-unit-invalid" "api-unit-health" 1
	api_service_state="not_started_credentials_missing"
	if [ -f "$state_dir/.admin_credentials" ]; then
		systemctl enable --now adlaire-ci-api >/dev/null 2>&1 || json_fail "bootstrap-script-health-failed" "api-service" 1
		systemctl is-active adlaire-ci-api >/dev/null 2>&1 || json_fail "bootstrap-script-health-failed" "api-service-health" 1
		api_service_state="started"
	fi
}

require_exact_args "$@"
require_fixed_inputs
for command_name in sh mktemp mkdir chmod install mv rm sha256sum systemctl; do
	require_command "$command_name"
done

stage_dir=$(mktemp -d "${TMPDIR:-/tmp}/adlaire-ci-bootstrap.XXXXXX") || json_fail "bootstrap-script-state-directory-invalid" "staging" 1
download_source
verify_binary_version "adlaire-ci-setup"
write_runner_units
commit_install
cleanup_stage
printf '{"status":"ok","source_channel":"%s","topology_role":"%s","source_identifier_kind":"%s","source_digest_sha256":"%s","install_dir":"%s","bin_dir":"%s","state_dir":"%s","service_user":"%s","installed_binaries":7,"runner_timer":"enabled","api_unit":"installed","api_service":"%s","rollback_action":"not_required","cleanup_result":"%s"}\n' \
	"$source_channel" "$topology_role" "$source_kind" "$target_sha256" "$install_dir" "$bin_dir" "$state_dir" "$service_user" "$api_service_state" "$cleanup_result"
