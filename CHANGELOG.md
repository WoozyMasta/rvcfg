# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog][],
and this project adheres to [Semantic Versioning][].

<!--
## Unreleased

### Added
### Changed
### Removed
-->

## [0.3.0][] - 2026-03-29

### Added

* Lint subsystem based on `lintkit` with rule catalog
  and generated diagnostics snapshot (`rules.yaml`).
* Public diagnostics catalog/docs updates for lint
  and parser/runtime diagnostics.

### Changed

* Lexer and token internals were reworked for a new token/position model.
* Lint architecture and rule wiring were redesigned;
  lint-related APIs and contracts were updated.
* Diagnostic model, codes, and related function contracts were revised
  for new lint/parse pipeline.
* Parser/process and source-map flows were updated
  to match new lexer/lint contracts.

[0.3.0]: https://github.com/WoozyMasta/rvcfg/compare/v0.2.0...v0.3.0

## [0.2.0][] - 2026-03-17

### Added

* `PreprocessMode` profiles: `strict`, `compat`, `extended`.
* `ParseTokens(...)` API for parsing already tokenized input
  without a second lexing pass.
* Extended preprocessor intrinsics for deterministic templating workflows:
  `__PATH_NORM(...)`, `__STR_*`, `__FILES_*`, `__FOR_RANGE_RENDER(...)`,
  `__FOR_EACH_RENDER(...)`.
* Preprocessor documentation in `PREPROCESSOR.md`.

### Changed

* Preprocessor defaults are now strict and side-effect free by default:
  `__EXEC`/`__EVAL`, dynamic, and filename intrinsics stay opt-in unless
  compatibility profile is selected.
* Parser/preprocessor benchmark paths were optimized to reduce runtime memory
  usage and speed up large config processing.

### Fixed

* Macro expansion parity for edge cases: malformed function-like calls,
  argument mismatch handling, token-paste/stringify corner cases, and
  expansion safety inside strings/comments.
* Include and directive handling parity for practical corpus cases:
  angle/quote include forms, no-space directive variants,
  and directive-tail behavior.

### Removed

* Legacy `MACROS.md` document in favor of `PREPROCESSOR.md`.

[0.2.0]: https://github.com/WoozyMasta/rvcfg/compare/v0.1.1...v0.2.0

## [0.1.1][] - 2026-02-26

### Added

* `RenderFile()` and `RenderFileWithOptions()` helpers
  for AST-to-text rendering.
* `ArrayWrapByName` option for per-assignment multiline wrapping.

[0.1.1]: https://github.com/WoozyMasta/rvcfg/compare/v0.1.0...v0.1.1

## [0.1.0][] - 2026-02-08

### Added

* First public release

[0.1.0]: https://github.com/WoozyMasta/rvcfg/tree/v0.1.0

<!--links-->
[Keep a Changelog]: https://keepachangelog.com/en/1.1.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html
