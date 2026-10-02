# Changelog

## V.0.0-dev

### Added

- Phase 13 implementation alignment and quality evidence package.
- Phase 13 required-check workflow.
- Release governance artifacts: `LICENSE`, `SECURITY.md`,
  `CONTRIBUTING.md`, `CODEOWNERS`, and `CHANGELOG.md`.

### Security

- Zero external production dependency policy is fixed to Go standard library
  and in-repository implementation only.
- Password KDF, token handling, secret masking, and statefile safety are tied
  to Phase 13 evidence records.

### Migration

- Phase 13 keeps the Phase 12 owner package five-file structure.
