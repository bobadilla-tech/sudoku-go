# What this is

Single-package Go library (`github.com/bobadilla-tech/sudoku-go`, package `sudoku`) for generating and validating Sudoku puzzles. Two source files total: `sudoku.go` (implementation) and `sudoku_test.go` (tests). No CLI, no subpackages, no external service code — this repo IS the library.

It's consumed by the Sudoku API on requiems.xyz (a separate repo/service) — treat the public API (`Generate`, `Validate`, `Difficulty`, `Grid`, `Puzzle`) as a versioned contract; breaking changes affect that consumer.

## Commands

```bash
go test ./...                                    # run all tests
go test ./... -run TestGenerate_AllDifficulties   # run a single test
go test ./... -covermode=atomic -coverprofile=coverage.out  # coverage (matches CI)
go vet ./...
gofmt -l .                                        # check formatting
```

No Makefile, no lint config beyond `go vet`. CI (`.github/workflows/ci.yml`) runs tests with coverage and uploads to Codecov on every push/PR to `main`.

## Architecture

Everything lives in `sudoku.go`. The generation algorithm, in order:

1. **`base`** — a single hardcoded, pre-verified valid 9×9 solution (simple diagonal-shift pattern), used as the seed for every puzzle.
2. **`shuffle`** — applies five structure-preserving transformations to `base` to produce a fresh random solution: digit remapping, row shuffle within bands, column shuffle within stacks, band shuffle, stack shuffle. Each step preserves Sudoku validity by construction, so the result never needs re-validation.
3. **`removeCells`** — blanks a difficulty-dependent number of cells (`cellsToRemove` map: Easy=36, Medium=46, Hard=52 removed) at random positions to produce the puzzle from the solution.
4. **`Generate(d Difficulty)`** ties these together, returns `Puzzle{Grid, Solution}`.
5. **`Validate(g Grid)`** independently checks rows/columns/3×3 boxes for exactly digits 1-9 with no repeats — used both for testing and as the public validity checker.

Known limitation (documented in README): puzzle **solution uniqueness is not enforced** — cell removal is random and doesn't verify the resulting puzzle has only one solution. This is a known gap, not an oversight; don't "fix" it silently if touching this code without flagging the tradeoff (uniqueness checking is nontrivial to add and was intentionally deferred).

`rand/v2` with `rand.NewPCG` is used for all randomness (not crypto-safe, and doesn't need to be — puzzles aren't security-sensitive).
