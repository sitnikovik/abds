# AGENTS.md

Guidance for AI agents working in this repository.

## Project overview

- Course repository for big data labs and seminars, based on the ClickHouse model of an
  energy company ("EnergoSbyt").
- `lab/` contains student assignments. Each lab has a `README.md` (in Russian) that is the student-facing spec.
- `prac/` contains seminar materials. Do not change `prac/` unless the task says so.
- `autograder/` is a Go module that checks student solutions in GitHub Actions.

## Language

- Student-facing text (lab READMEs, FAQ) is in Russian.
- Code, comments, GoDoc, commit messages and this file are in English.

## Code style

- Follow `STYLE.md` for all Go code.
- Declarative style: constructors with functional options, immutable structs, no setters, no getters.
- Do not check values for `nil`. Return `error` values.
- Fields are unexported. Types are exported and created through constructors.

## Documentation

- Every package, type, function, method and field has a GoDoc comment.
- Write and check GoDoc with the `/godoc` skill.
- Do not write comments inside function bodies.

## Autograder

- Module: `autograder/`, Go 1.25.7 (see `autograder/go.mod`).
- Every check inside an autograder step MUST use `buggles`
  (`github.com/sitnikovik/buggles/expect/...`), not hand-written assertions.
- `buggles` is also the check library in `_test.go` files.
- Compose checks with `expect.All`, `expect.Any` and `expect.AtLeast`,
  and return the combined error from `Run`.
- Do not write `if err != nil { return ... }` for assertions.
- Each step is a package under `internal/step/...` with a `Run(ctx, vars) error` method.
- Pin the `buggles` version in `autograder/go.mod`. The framework is not stable yet.
- ClickHouse access uses `github.com/ClickHouse/clickhouse-go/v2` over the native protocol.
  - Declarative repo layer: `internal/repo/<entity>/<operation>/`.
    Each operation has a constructor `NewX(db, opts...)` and a method `Run(ctx) error`.
  - Inserts go through batches: `PrepareBatch`, `Append`, `Send`.
  - Values are passed as query arguments, never concatenated into SQL.
  - Table and column names live in `consts.go`.

## Tests

- Express every expectation with `buggles`. Use testify only as the `t` binding via `assert.NoError`.
- A test with a single expectation calls `assert.NoError(t, <buggles check>)` directly.
- A test with more than one expectation combines them with `expect.All`, `expect.Any`
  or `expect.AtLeast` and passes the result to one `assert.NoError(t, ...)`.
  Do not wrap a single check in `expect.All`.
- Call `t.Parallel()` wherever it is safe.
- No table-driven tests. Write one `t.Run` block per case.
- Put tests in `_test` packages.
- Get the checked result directly: call the constructor and the method
  in one expression instead of storing one-time intermediate objects
  in variables (see STYLE.md, "Direct result in tests").
- Do not pass a literal `nil` for a non-error parameter. Declare a variable
  with the default zero value named after the parameter and pass it.
  A literal `nil` stays fine for `error` parameters (see STYLE.md,
  "Named defaults instead of nil").
- Keep coverage at or above `COVERAGE_THRESHOLD` from `.github/ci.env`.

## Commands

- `make -C autograder check`: unit tests, integration tests, coverage and lint.
  Must pass before a PR that touches `autograder/`.
- `make -C autograder lint`: lint only.
- `make -C autograder unit-test`: unit tests only.
