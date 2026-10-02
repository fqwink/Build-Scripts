# Security Policy

## Supported Versions

Adlaire CI is pre-release software. Security fixes are accepted only for the
current `main` branch and active release-candidate work unless a release branch
is explicitly documented in `docs/ROADMAP.md`.

## Reporting

Do not publish secrets, exploit payloads, private infrastructure details, or
vulnerability reproduction data in public issues.

Report a vulnerability through GitHub private vulnerability reporting when it
is available for this repository. If private reporting is unavailable, open a
public issue that contains only a non-sensitive summary and request a private
handoff channel from the repository owner.

## Secret Leakage

If a report includes leaked credentials, tokens, private keys, session values,
or deployment paths, the report must identify the affected artifact and the
required rotation boundary without reposting the secret value.

## Response Policy

Security reports are triaged against `docs/SPEC.md`, `docs/DETAIL_INDEX.md`,
and the relevant owner detail file. A fix is complete only when the affected
specification, implementation, tests, evidence records, and release notes are
consistent.
