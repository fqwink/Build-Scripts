# Contributing

Contributions must follow `AGENTS.md` and `docs/SPEC.md`.

## Required Order

1. Confirm the responsible source of truth in `docs/SPEC.md`.
2. Update the responsible specification before implementation when behavior,
   policy, state, evidence, or release criteria change.
3. Implement only the active Phase scope recorded in `docs/ROADMAP.md`.
4. Keep implementation and verification evidence aligned with
   `docs/DETAIL_INDEX.md` and the relevant `docs/details/*.md` file.

## Pull Requests

Pull requests are Phase-bound. Do not split one Phase across multiple pull
requests, and do not mix future Phase implementation into the current Phase.

Every pull request must include meaningful verification evidence. Skipped,
unavailable, or environment-limited checks must be recorded explicitly.

## Tests

Tests must verify the specification contract from observable inputs, outputs,
state changes, logs, and evidence records. Passing coverage numbers alone are
not sufficient completion evidence.
