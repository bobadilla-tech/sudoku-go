# sudoku-go

A Go library for generating and validating Sudoku puzzles.

## Install

```bash
go get github.com/bobadilla-tech/sudoku-go
```

## Usage

### Generate a puzzle

```go
import sudoku "github.com/bobadilla-tech/sudoku-go"

puzzle, err := sudoku.Generate(sudoku.Easy)
if err != nil {
    log.Fatal(err)
}

// puzzle.Grid     — 9×9 board with 0s for empty cells
// puzzle.Solution — 9×9 completed solution
```

### Validate a completed grid

```go
valid := sudoku.Validate(puzzle.Solution)
fmt.Println(valid) // true
```

## Difficulty levels

| Constant        | Givens | Empty cells |
|-----------------|--------|-------------|
| `sudoku.Easy`   | 45     | 36          |
| `sudoku.Medium` | 35     | 46          |
| `sudoku.Hard`   | 29     | 52          |

## Algorithm

1. Start from a pre-verified valid base grid.
2. Remap digits with a random permutation.
3. Shuffle rows within each band (rows 0–2, 3–5, 6–8).
4. Shuffle columns within each stack (columns 0–2, 3–5, 6–8).
5. Shuffle the three bands.
6. Shuffle the three stacks.
7. Remove cells according to difficulty.

All transformations preserve Sudoku validity, so the result is always a valid puzzle with a known unique solution.

## License

Apache 2.0 — see [LICENSE](LICENSE).