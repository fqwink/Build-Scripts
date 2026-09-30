class AdlaireCIError extends Error {
  constructor({ status, message, details = null, responseBody = null }) {
    super(message);
    this.name = "AdlaireCIError";
    this.status = status;
    this.details = Array.isArray(details) ? details : null;
    this.responseBody = responseBody;
  }
}

class AdlaireCI {
  constructor(options) {
    if (!runtimeSupported()) {
      throw new TypeError("Unsupported browser runtime");
    }
    if (options === null || typeof options !== "object" || Array.isArray(options)) {
      throw new TypeError("Invalid argument: baseUrl");
    }
    this.baseUrl = normalizeBaseUrl(options.baseUrl);
    this._token = null;
  }

  async login(password) {
    requireString(password, "password");
    const result = await this._request("/api/login", { method: "POST", body: { password } });
    if (result && typeof result.token === "string" && result.token !== "") {
      this._token = result.token;
    } else if (result && result.totp_required === true) {
      this._token = null;
    }
    return result;
  }

  async loginTotp(ticket, code) {
    requireString(ticket, "ticket");
    requireString(code, "code");
    const result = await this._request("/api/login/totp", { method: "POST", body: { ticket, code } });
    if (result && typeof result.token === "string" && result.token !== "") {
      this._token = result.token;
    }
    return result;
  }

  async logout() {
    try {
      return await this._request("/api/logout", { method: "POST" });
    } finally {
      this._token = null;
    }
  }

  changePassword(currentPassword, newPassword) {
    requireString(currentPassword, "currentPassword");
    requireString(newPassword, "newPassword");
    return this._request("/api/change-password", {
      method: "POST",
      body: { current_password: currentPassword, new_password: newPassword },
    });
  }

  getAuditLog({ limit = 100, offset = 0, actor = null, action = null, result = null } = {}) {
    return this._request("/api/audit-log", { query: this._query({ limit, offset, actor, action, result }) });
  }

  getStatus() {
    return this._request("/api/status");
  }

  triggerBuild() {
    return this._request("/api/build", { method: "POST" });
  }

  getLogs(n = 100, q = "") {
    return this._request("/api/logs", { query: this._query({ n, q }, ["q"]) });
  }

  getHistory({ page = 1, perPage = 20, trigger = null, tag = null, flagged = null, failureCategory = null } = {}) {
    requireNumber(page, "page");
    requireNumber(perPage, "perPage");
    if (flagged !== null && typeof flagged !== "boolean") {
      throw new TypeError("Invalid argument: flagged");
    }
    return this._request("/api/history", {
      query: this._query({ page, per_page: perPage, trigger, tag, flagged, failure_category: failureCategory }),
    });
  }

  getSysinfo() {
    return this._request("/api/sysinfo");
  }

  getSchedule() {
    return this._request("/api/schedule");
  }

  getNotifyConfig() {
    return this._request("/api/notify-config");
  }

  setNotifyConfig(config) {
    requireObject(config, "config");
    return this._request("/api/notify-config", { method: "POST", body: config });
  }

  getConfig() {
    return this._request("/api/config");
  }

  validateConfig(config) {
    requireObject(config, "config");
    return this._request("/api/config/validate", { method: "POST", body: config });
  }

  setConfig(config) {
    requireObject(config, "config");
    return this._request("/api/config", { method: "POST", body: config });
  }

  health() {
    return this._request("/api/health");
  }

  getPatStatus() {
    return this._request("/api/pat-status");
  }

  getAccessLog({ limit = 100, offset = 0 } = {}) {
    requireNumber(limit, "limit");
    requireNumber(offset, "offset");
    return this._request("/api/access-log", { query: this._query({ limit, offset }) });
  }

  getApiAccessLog({ limit = 100, offset = 0, method = null, path = null, status = null } = {}) {
    return this._request("/api/api-access-log", {
      query: this._query({ limit, offset, method, path, status }),
    });
  }

  getStats(days = 7) {
    return this._request("/api/stats", { query: this._query({ days }) });
  }

  exportLogs() {
    return this._request("/api/logs/export");
  }

  cleanupLogs() {
    return this._request("/api/logs/cleanup", { method: "POST" });
  }

  archiveLogs() {
    return this._request("/api/logs/archive", { method: "POST" });
  }

  getRepoInfo() {
    return this._request("/api/repo-info");
  }

  backup() {
    return this._request("/api/backup");
  }

  restore(config) {
    requireObject(config, "config");
    return this._request("/api/restore", { method: "POST", body: config });
  }

  notifyTest() {
    return this._request("/api/notify-test", { method: "POST" });
  }

  buildForce() {
    return this._request("/api/build/force", { method: "POST" });
  }

  patVerify() {
    return this._request("/api/pat-verify", { method: "POST" });
  }

  getHistoryLog(id) {
    return this._request(`/api/history/${this._validateId(id)}/log`);
  }

  cancelBuild() {
    return this._request("/api/build/cancel", { method: "POST" });
  }

  resetCircuitBreaker() {
    return this._request("/api/circuit-breaker/reset", { method: "POST" });
  }

  async streamBuild(onLine, onEnd) {
    if (typeof onLine !== "function") {
      throw new TypeError("Invalid argument: onLine");
    }
    if (typeof onEnd !== "function") {
      throw new TypeError("Invalid argument: onEnd");
    }
    this._requireToken();
    const controller = new AbortController();
    let timedOut = false;
    const timeout = setTimeout(() => {
      timedOut = true;
      controller.abort();
    }, 30000);
    let response;
    try {
      response = await fetch(this.baseUrl + "/api/build/stream", {
        method: "GET",
        headers: { Accept: "text/event-stream", Authorization: `Bearer ${this._token}` },
        signal: controller.signal,
      });
      if (!response.ok) {
        return await this._json(response);
      }
      if (!isEventStreamMediaType(response.headers.get("Content-Type") || "")) {
        throw new AdlaireCIError({ status: 0, message: "Invalid SSE response" });
      }
      if (!response.body || typeof response.body.getReader !== "function") {
        throw new AdlaireCIError({ status: 0, message: "Invalid SSE response" });
      }
    } catch (error) {
      if (error instanceof AdlaireCIError) {
        throw error;
      }
      if (timedOut && isAbortError(error)) {
        throw new AdlaireCIError({ status: 0, message: "Request timeout" });
      }
      throw new AdlaireCIError({ status: 0, message: "Network error" });
    } finally {
      clearTimeout(timeout);
    }
    const reader = response.body.getReader();
    const terminal = createStreamTerminal(controller);
    terminal.setReader(reader);
    readBuildSse(reader, onLine, onEnd, terminal);
    return terminal.handle;
  }

  setLogLevel(level) {
    requireString(level, "level");
    return this._request("/api/log-level", { method: "POST", body: { level } });
  }

  updatePat(token) {
    requireString(token, "token");
    return this._request("/api/pat-update", { method: "POST", body: { token } });
  }

  getDashboard() {
    return this._request("/api/dashboard");
  }

  getNotifyLog({ limit = 100, offset = 0 } = {}) {
    requireNumber(limit, "limit");
    requireNumber(offset, "offset");
    return this._request("/api/notify-log", { query: this._query({ limit, offset }) });
  }

  getSessions() {
    return this._request("/api/sessions");
  }

  revokeAllSessions() {
    return this._request("/api/sessions/revoke-all", { method: "POST" });
  }

  getTotpStatus() {
    return this._request("/api/auth/totp-status");
  }

  setupTotp() {
    return this._request("/api/auth/totp-setup", { method: "POST" });
  }

  confirmTotp(code) {
    requireString(code, "code");
    return this._request("/api/auth/totp-confirm", { method: "POST", body: { code } });
  }

  disableTotp(code) {
    requireString(code, "code");
    return this._request("/api/auth/totp", { method: "DELETE", body: { code } });
  }

  setScheduleInterval(seconds) {
    requireNumber(seconds, "seconds");
    return this._request("/api/schedule/interval", { method: "POST", body: { interval_seconds: seconds } });
  }

  pauseSchedule() {
    return this._request("/api/schedule/pause", { method: "POST" });
  }

  resumeSchedule() {
    return this._request("/api/schedule/resume", { method: "POST" });
  }

  setAllowedHours(from, to) {
    requireNumber(from, "from");
    requireNumber(to, "to");
    return this._request("/api/schedule/allowed-hours", { method: "POST", body: { from, to } });
  }

  clearAllowedHours() {
    return this._request("/api/schedule/allowed-hours", { method: "POST", body: { from: null, to: null } });
  }

  setForceInterval(hours) {
    requireNumber(hours, "hours");
    return this._request("/api/schedule/force-interval", { method: "POST", body: { hours } });
  }

  setBuildCooldown(seconds) {
    requireNumber(seconds, "seconds");
    return this._request("/api/schedule/cooldown", { method: "POST", body: { seconds } });
  }

  searchLogs(q = "", from = "", to = "", level = undefined) {
    return this._request("/api/logs/search", {
      query: this._query({ q, from, to, level }, ["q", "from", "to"]),
    });
  }

  getOutputMeta() {
    return this._request("/api/output-meta");
  }

  getStatsTimeline(days = 30) {
    return this._request("/api/stats/timeline", { query: this._query({ days }) });
  }

  getStatsBuildDuration(n = 10) {
    return this._request("/api/stats/build-duration", { query: this._query({ n }) });
  }

  getBuildTrends(n = 100) {
    return this._request("/api/stats/build-trends", { query: this._query({ n }) });
  }

  getDiagnostics() {
    return this._request("/api/diagnostics");
  }

  getRateLimit() {
    return this._request("/api/rate-limit");
  }

  getDiskUsage() {
    return this._request("/api/disk-usage");
  }

  getConfigLog({ limit = 100, offset = 0 } = {}) {
    requireNumber(limit, "limit");
    requireNumber(offset, "offset");
    return this._request("/api/config-log", { query: this._query({ limit, offset }) });
  }

  getApiRateLimit() {
    return this._request("/api/api-rate-limit");
  }

  setApiRateLimit(policy) {
    requireObject(policy, "policy");
    return this._request("/api/api-rate-limit", { method: "POST", body: policy });
  }

  getBranchConfig() {
    return this._request("/api/branch-config");
  }

  setBranchConfig(branches) {
    requireArray(branches, "branches");
    return this._request("/api/branch-config", { method: "POST", body: { branches } });
  }

  notifyWeeklySummary() {
    return this._request("/api/notify/weekly-summary", { method: "POST" });
  }

  getWebhookEvents(limit = 50, offset = 0) {
    return this._request("/api/webhook-events", { query: this._query({ limit, offset }) });
  }

  getWebhookConfig() {
    return this._request("/api/webhook-config");
  }

  setWebhookConfig(secret) {
    requireString(secret, "secret");
    return this._request("/api/webhook-config", { method: "POST", body: { secret } });
  }

  getBuildChainConfig() {
    return this._request("/api/build-chain-config");
  }

  setBuildChainConfig(chains) {
    requireArray(chains, "chains");
    return this._request("/api/build-chain-config", { method: "POST", body: { chains } });
  }

  getApprovals() {
    return this._request("/api/approvals");
  }

  approveBuild(id) {
    return this._request(`/api/approvals/${this._validateId(id)}/approve`, { method: "POST" });
  }

  rejectBuild(id) {
    return this._request(`/api/approvals/${this._validateId(id)}/reject`, { method: "POST" });
  }

  getHistoryComment(id) {
    return this._request(`/api/history/${this._validateId(id)}/comment`);
  }

  setHistoryComment(id, comment) {
    requireString(comment, "comment", true);
    return this._request(`/api/history/${this._validateId(id)}/comment`, { method: "POST", body: { comment } });
  }

  setRepoConfig({ owner = undefined, repo = undefined } = {}) {
    const patch = {};
    if (owner !== undefined) {
      patch.owner = owner;
    }
    if (repo !== undefined) {
      patch.repo = repo;
    }
    if (owner === undefined && repo === undefined) {
      throw new TypeError("Invalid argument: config");
    }
    return this._request("/api/repo-config", { method: "POST", body: patch });
  }

  exportHistory() {
    return this._request("/api/history/export");
  }

  setHistoryFlag(id, flagged) {
    if (typeof flagged !== "boolean") {
      throw new TypeError("Invalid argument: flagged");
    }
    return this._request(`/api/history/${this._validateId(id)}/flag`, { method: "POST", body: { flagged } });
  }

  setHistoryTags(id, tags) {
    requireArray(tags, "tags");
    return this._request(`/api/history/${this._validateId(id)}/tags`, { method: "POST", body: { tags } });
  }

  rollbackHistory(id) {
    return this._request(`/api/history/${this._validateId(id)}/rollback`, { method: "POST" });
  }

  getTokens() {
    return this._request("/api/tokens");
  }

  createToken(label, scopes = ["read"], expiresAt = null) {
    requireString(label, "label");
    requireArray(scopes, "scopes");
    return this._request("/api/tokens", {
      method: "POST",
      body: { label, scopes, expires_at: expiresAt },
    });
  }

  revokeToken(id) {
    return this._request(`/api/tokens/${this._validateId(id)}`, { method: "DELETE" });
  }

  getSnapshots() {
    return this._request("/api/snapshots");
  }

  downloadSnapshot(id) {
    return this._request(`/api/snapshots/${this._validateId(id)}/download`, { binary: true });
  }

  deleteSnapshot(id) {
    return this._request(`/api/snapshots/${this._validateId(id)}`, { method: "DELETE" });
  }

  getMaintenance() {
    return this._request("/api/maintenance");
  }

  enableMaintenance(reason) {
    requireString(reason, "reason");
    return this._request("/api/maintenance/enable", { method: "POST", body: { reason } });
  }

  disableMaintenance() {
    return this._request("/api/maintenance/disable", { method: "POST" });
  }

  getAccessControl() {
    return this._request("/api/access-control");
  }

  setAccessControl(allowList) {
    requireArray(allowList, "allowList");
    return this._request("/api/access-control", { method: "POST", body: { allow: allowList } });
  }

  getHooks() {
    return this._request("/api/hooks");
  }

  addHook(phase, commandArgs, abortOnFailure = true, timeoutSeconds = 300) {
    requireString(phase, "phase");
    requireArray(commandArgs, "commandArgs");
    if (typeof abortOnFailure !== "boolean") {
      throw new TypeError("Invalid argument: abortOnFailure");
    }
    requireNumber(timeoutSeconds, "timeoutSeconds");
    return this._request("/api/hooks", {
      method: "POST",
      body: { phase, command_args: commandArgs, abort_on_failure: abortOnFailure, timeout_seconds: timeoutSeconds },
    });
  }

  deleteHook(id) {
    return this._request(`/api/hooks/${this._validateId(id)}`, { method: "DELETE" });
  }

  getHookLog(id) {
    return this._request(`/api/hooks/${this._validateId(id)}/log`);
  }

  getAlertRules() {
    return this._request("/api/alert-rules");
  }

  addAlertRule(metric, operator, threshold, level, message) {
    requireString(metric, "metric");
    requireString(operator, "operator");
    requireNumber(threshold, "threshold");
    requireString(level, "level");
    requireString(message, "message");
    return this._request("/api/alert-rules", { method: "POST", body: { metric, operator, threshold, level, message } });
  }

  deleteAlertRule(id) {
    return this._request(`/api/alert-rules/${this._validateId(id)}`, { method: "DELETE" });
  }

  getTagRules() {
    return this._request("/api/tag-rules");
  }

  addTagRule(condition, tags) {
    requireObject(condition, "condition");
    requireArray(tags, "tags");
    return this._request("/api/tag-rules", { method: "POST", body: { condition, tags } });
  }

  deleteTagRule(id) {
    return this._request(`/api/tag-rules/${this._validateId(id)}`, { method: "DELETE" });
  }

  verifyOutput() {
    return this._request("/api/verify-output", { method: "POST" });
  }

  getPipelineConfig() {
    return this._request("/api/pipeline-config");
  }

  setPipelineConfig(config) {
    requireObject(config, "config");
    return this._request("/api/pipeline-config", { method: "POST", body: config });
  }

  getNotes() {
    return this._request("/api/notes");
  }

  setNotes(content) {
    requireString(content, "content", true);
    return this._request("/api/notes", { method: "POST", body: { content } });
  }

  getSmtpConfig() {
    return this._request("/api/smtp-config");
  }

  setSmtpConfig(config) {
    requireObject(config, "config");
    const body = {};
    for (const [key, value] of Object.entries(config)) {
      if (key === "password" && value === undefined) {
        continue;
      }
      body[key] = value;
    }
    return this._request("/api/smtp-config", { method: "POST", body });
  }

  smtpTest() {
    return this._request("/api/smtp-test", { method: "POST" });
  }

  getQueue() {
    return this._request("/api/queue");
  }

  clearQueue() {
    return this._request("/api/queue", { method: "DELETE" });
  }

  getDashboardLayout() {
    return this._request("/api/dashboard-layout");
  }

  setDashboardLayout(widgets) {
    requireArray(widgets, "widgets");
    return this._request("/api/dashboard-layout", { method: "POST", body: { widgets } });
  }

  async _request(path, { method = "GET", query = null, body = undefined, binary = false } = {}) {
    const url = this.baseUrl + path + (query ? `?${query}` : "");
    const headers = { Accept: binary ? "application/octet-stream, application/json" : "application/json" };
    const options = { method, headers };
    if (this._token) {
      headers.Authorization = `Bearer ${this._token}`;
    }
    if (body !== undefined) {
      headers["Content-Type"] = "application/json";
      options.body = JSON.stringify(body);
    }
    const controller = new AbortController();
    let timedOut = false;
    const timeout = setTimeout(() => {
      timedOut = true;
      controller.abort();
    }, 30000);
    options.signal = controller.signal;
    try {
      const response = await fetch(url, options);
      if (binary && response.ok) {
        return await readBinaryResponse(response);
      }
      return await this._json(response);
    } catch (error) {
      if (error instanceof AdlaireCIError) {
        throw error;
      }
      if (timedOut && isAbortError(error)) {
        throw new AdlaireCIError({ status: 0, message: "Request timeout" });
      }
      throw new AdlaireCIError({ status: 0, message: "Network error" });
    } finally {
      clearTimeout(timeout);
    }
  }

  async _json(response) {
    const contentType = response.headers.get("Content-Type") || "";
    const isJson = isJsonMediaType(contentType);
    const text = await response.text();
    if (response.ok) {
      if (!isJson) {
        throw new AdlaireCIError({ status: 0, message: "Invalid JSON response" });
      }
      if (text === "") {
        throw new AdlaireCIError({ status: 0, message: "Empty JSON response" });
      }
      try {
        return JSON.parse(text);
      } catch (_error) {
        throw new AdlaireCIError({ status: 0, message: "Invalid JSON response" });
      }
    }
    let message = `HTTP ${response.status}`;
    let details = null;
    let responseBody = null;
    if (isJson && text !== "") {
      try {
        responseBody = JSON.parse(text);
        if (responseBody && typeof responseBody.error === "string") {
          message = responseBody.error;
        }
        if (responseBody && Array.isArray(responseBody.details)) {
          details = responseBody.details;
        }
      } catch (_error) {
        responseBody = text.slice(0, 4000);
      }
    } else if (text !== "") {
      responseBody = text.slice(0, 4000);
    }
    if (response.status === 401) {
      this._clearTokenOn401(response.status);
    }
    throw new AdlaireCIError({ status: response.status, message, details, responseBody });
  }

  _query(values, keepEmpty = []) {
    const params = [];
    for (const [key, value] of Object.entries(values)) {
      if (value === undefined || value === null) {
        continue;
      }
      if (value === "" && !keepEmpty.includes(key)) {
        continue;
      }
      params.push(`${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`);
    }
    return params.join("&");
  }

  _requireToken() {
    if (!this._token) {
      throw new AdlaireCIError({ status: 401, message: "Unauthorized" });
    }
  }

  _validateId(id) {
    requireString(id, "id");
    if (!/^[A-Za-z0-9_-]{1,64}$/.test(id)) {
      throw new TypeError("Invalid argument: id");
    }
    return encodeURIComponent(id);
  }

  _clearTokenOn401(status) {
    if (status === 401) {
      this._token = null;
    }
  }
}

function runtimeSupported() {
  return typeof fetch === "function" &&
    typeof AbortController === "function" &&
    typeof ReadableStream !== "undefined" &&
    typeof TextDecoder === "function";
}

function normalizeBaseUrl(baseUrl) {
  if (typeof baseUrl !== "string" || baseUrl.trim() === "") {
    throw new TypeError("Invalid argument: baseUrl");
  }
  let url;
  try {
    url = new URL(baseUrl);
  } catch (_error) {
    throw new TypeError("Invalid argument: baseUrl");
  }
  if ((url.protocol !== "http:" && url.protocol !== "https:") ||
      url.username !== "" ||
      url.password !== "" ||
      url.pathname !== "/" ||
      url.search !== "" ||
      url.hash !== "") {
    throw new TypeError("Invalid argument: baseUrl");
  }
  return url.origin;
}

function requireString(value, name, allowEmpty = false) {
  if (typeof value !== "string" || (!allowEmpty && value === "")) {
    throw new TypeError(`Invalid argument: ${name}`);
  }
}

function requireNumber(value, name) {
  if (typeof value !== "number" || !Number.isFinite(value)) {
    throw new TypeError(`Invalid argument: ${name}`);
  }
}

function requireArray(value, name) {
  if (!Array.isArray(value)) {
    throw new TypeError(`Invalid argument: ${name}`);
  }
}

function requireObject(value, name) {
  if (value === null || typeof value !== "object" || Array.isArray(value)) {
    throw new TypeError(`Invalid argument: ${name}`);
  }
}

function isAbortError(error) {
  return error && error.name === "AbortError";
}

function mediaType(contentType) {
  return contentType.split(";")[0].trim().toLowerCase();
}

function isJsonMediaType(contentType) {
  const type = mediaType(contentType);
  if (type === "application/json") {
    return true;
  }
  if (!type.startsWith("application/")) {
    return false;
  }
  const subtype = type.slice("application/".length);
  const suffixIndex = subtype.lastIndexOf("+json");
  return suffixIndex > 0 && suffixIndex === subtype.length - "+json".length;
}

function isEventStreamMediaType(contentType) {
  return mediaType(contentType) === "text/event-stream";
}

function isOctetStreamMediaType(contentType) {
  return mediaType(contentType) === "application/octet-stream";
}

async function readBinaryResponse(response) {
  if (!isOctetStreamMediaType(response.headers.get("Content-Type") || "")) {
    throw new AdlaireCIError({ status: 0, message: "Invalid binary response" });
  }
  return await response.blob();
}

function createStreamTerminal(controller) {
  let reader = null;
  let closed = false;
  let error = null;
  let settled = false;
  let resolveDone;
  let rejectDone;
  const done = new Promise((resolve, reject) => {
    resolveDone = resolve;
    rejectDone = reject;
  });
  const terminal = {
    handle: {
      close() {
        if (settled) {
          return;
        }
        settled = true;
        closed = true;
        error = null;
        try {
          controller.abort();
        } catch (_error) {
          // close() is intentionally idempotent.
        }
        if (reader && typeof reader.cancel === "function") {
          reader.cancel().catch(() => {});
        }
        resolveDone(null);
      },
      get closed() {
        return closed;
      },
      get error() {
        return error;
      },
      done,
    },
    setReader(value) {
      reader = value;
    },
    isClosed() {
      return closed;
    },
    succeed(summary) {
      if (settled) {
        return;
      }
      settled = true;
      closed = true;
      error = null;
      resolveDone(summary);
    },
    fail(value) {
      if (settled) {
        return;
      }
      settled = true;
      closed = true;
      error = value;
      try {
        controller.abort();
      } catch (_error) {
        // The stream may already be closed.
      }
      rejectDone(value);
    },
  };
  return terminal;
}

async function readBuildSse(reader, onLine, onEnd, terminal) {
  const decoder = new TextDecoder("utf-8", { fatal: true });
  let buffer = "";
  let summary = null;

  try {
    while (true) {
      const { value, done } = await reader.read();
      if (done) {
        break;
      }
      let text;
      try {
        text = decoder.decode(value, { stream: true });
      } catch (_error) {
        throw invalidSseFrame();
      }
      if (text.includes("\r")) {
        throw invalidSseFrame();
      }
      if (summary !== null && text.length > 0) {
        throw invalidSseFrame();
      }
      buffer += text;

      let frameEnd = buffer.indexOf("\n\n");
      while (frameEnd !== -1) {
        const frame = buffer.slice(0, frameEnd);
        buffer = buffer.slice(frameEnd + 2);
        if (summary !== null) {
          throw invalidSseFrame();
        }
        const payload = parseBuildSseFrame(frame);
        if (payload.type === "log") {
          try {
            onLine(payload.line);
          } catch (_error) {
            throw streamCallbackFailed();
          }
        } else {
          summary = { status: payload.status, duration_seconds: payload.duration_seconds };
        }
        frameEnd = buffer.indexOf("\n\n");
      }
    }

    let tail;
    try {
      tail = decoder.decode();
    } catch (_error) {
      throw invalidSseFrame();
    }
    if (tail.includes("\r")) {
      throw invalidSseFrame();
    }
    buffer += tail;
    if (buffer !== "" || summary === null) {
      throw invalidSseFrame();
    }
    try {
      onEnd(summary);
    } catch (_error) {
      throw streamCallbackFailed();
    }
    terminal.succeed(summary);
  } catch (error) {
    if (terminal.isClosed()) {
      return;
    }
    if (error instanceof AdlaireCIError) {
      terminal.fail(error);
      return;
    }
    terminal.fail(new AdlaireCIError({ status: 0, message: "Network error" }));
  } finally {
    if (typeof reader.releaseLock === "function") {
      try {
        reader.releaseLock();
      } catch (_error) {
        // The stream may already be released by close().
      }
    }
  }
}

function parseBuildSseFrame(frame) {
  if (frame === "" || frame.includes("\n") || frame.includes("\r") || !frame.startsWith("data: ")) {
    throw invalidSseFrame();
  }
  let payload;
  try {
    payload = JSON.parse(frame.slice("data: ".length));
  } catch (_error) {
    throw invalidSseFrame();
  }
  if (payload === null || typeof payload !== "object" || Array.isArray(payload)) {
    throw invalidSseFrame();
  }
  if (payload.type === "log") {
    if (!hasExactKeys(payload, ["type", "line", "at"]) ||
        typeof payload.line !== "string" ||
        !isUtcSecond(payload.at)) {
      throw invalidSseFrame();
    }
    return payload;
  }
  if (payload.type === "end") {
    if (!hasExactKeys(payload, ["type", "status", "duration_seconds"]) ||
        !["running", "success", "failure", "cancelled"].includes(payload.status) ||
        !(payload.duration_seconds === null || (typeof payload.duration_seconds === "number" && Number.isFinite(payload.duration_seconds)))) {
      throw invalidSseFrame();
    }
    return payload;
  }
  throw invalidSseFrame();
}

function hasExactKeys(value, keys) {
  const actual = Object.keys(value);
  return actual.length === keys.length && keys.every((key) => actual.includes(key));
}

function isUtcSecond(value) {
  return typeof value === "string" && /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$/.test(value);
}

function invalidSseFrame() {
  return new AdlaireCIError({ status: 0, message: "Invalid SSE frame" });
}

function streamCallbackFailed() {
  return new AdlaireCIError({ status: 0, message: "Stream callback failed" });
}

export { AdlaireCI, AdlaireCIError };
