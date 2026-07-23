package sudoku

import (
	"fmt"
	"math/rand/v2"
)

// Difficulty represents the difficulty level of a generated Sudoku puzzle.
type Difficulty int

const (
	Easy   Difficulty = iota
	Medium Difficulty = iota
	Hard   Difficulty = iota
)

// Grid is a 9×9 Sudoku board; 0 represents an empty cell.
type Grid [9][9]int

// Puzzle holds the puzzle (with empty cells) and its unique solution.
type Puzzle struct {
	Grid     Grid
	Solution Grid
}

// cellsToRemove maps difficulty to the number of cells removed from the solution.
//
//   - Easy:   36 removed → 45 givens
//   - Medium: 46 removed → 35 givens
//   - Hard:   52 removed → 29 givens
var cellsToRemove = map[Difficulty]int{
	Easy:   36,
	Medium: 46,
	Hard:   52,
}

// base is a pre-verified valid Sudoku solution used as the starting template.
// All row/column/box constraints are satisfied.
var base = Grid{
	{1, 2, 3, 4, 5, 6, 7, 8, 9},
	{4, 5, 6, 7, 8, 9, 1, 2, 3},
	{7, 8, 9, 1, 2, 3, 4, 5, 6},
	{2, 3, 4, 5, 6, 7, 8, 9, 1},
	{5, 6, 7, 8, 9, 1, 2, 3, 4},
	{8, 9, 1, 2, 3, 4, 5, 6, 7},
	{3, 4, 5, 6, 7, 8, 9, 1, 2},
	{6, 7, 8, 9, 1, 2, 3, 4, 5},
	{9, 1, 2, 3, 4, 5, 6, 7, 8},
}

// Generate returns a new Sudoku puzzle at the requested difficulty level.
// Returns an error if the difficulty is not Easy, Medium, or Hard.
func Generate(d Difficulty) (Puzzle, error) {
	count, ok := cellsToRemove[d]
	if !ok {
		return Puzzle{}, fmt.Errorf("invalid difficulty %d", d)
	}

	rng := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())) //nolint:gosec
	solution := shuffle(base, rng)
	grid := removeCells(solution, count, rng)

	return Puzzle{Grid: grid, Solution: solution}, nil
}

// Validate returns true if g is a valid completed Sudoku:
// digits 1-9 appear exactly once in every row, column, and 3×3 box.
func Validate(g Grid) bool {
	// Check rows.
	for r := range 9 {
		if !validSet(g[r][:]) {
			return false
		}
	}
	// Check columns.
	for c := range 9 {
		col := make([]int, 9)
		for r := range 9 {
			col[r] = g[r][c]
		}
		if !validSet(col) {
			return false
		}
	}
	// Check 3×3 boxes.
	for boxR := range 3 {
		for boxC := range 3 {
			cells := make([]int, 0, 9)
			for r := range 3 {
				for c := range 3 {
					cells = append(cells, g[boxR*3+r][boxC*3+c])
				}
			}
			if !validSet(cells) {
				return false
			}
		}
	}
	return true
}

// validSet returns true if cells contains exactly the digits 1–9 with no repeats.
func validSet(cells []int) bool {
	seen := [10]bool{}
	for _, v := range cells {
		if v < 1 || v > 9 || seen[v] {
			return false
		}
		seen[v] = true
	}
	return true
}

// shuffle produces a new valid solution by applying five structure-preserving
// transformations to the base grid.
//
//  1. Digit remapping — replace each digit with a random permutation.
//  2. Row shuffle within bands — permute rows inside each band {0-2, 3-5, 6-8}.
//  3. Column shuffle within stacks — permute columns inside each stack {0-2, 3-5, 6-8}.
//  4. Band shuffle — permute the three row-bands.
//  5. Stack shuffle — permute the three column-stacks.
func shuffle(g Grid, rng *rand.Rand) Grid {
	// Step 1: digit remapping.
	perm := rng.Perm(9)
	for r := range 9 {
		for c := range 9 {
			g[r][c] = perm[g[r][c]-1] + 1
		}
	}

	// Step 2: shuffle rows within each band.
	for band := range 3 {
		rows := rng.Perm(3)
		base := band * 3
		var tmp [3][9]int
		for i, row := range rows {
			tmp[i] = g[base+row]
		}
		for i := range 3 {
			g[base+i] = tmp[i]
		}
	}

	// Step 3: shuffle columns within each stack.
	for stack := range 3 {
		cols := rng.Perm(3)
		base := stack * 3
		for r := range 9 {
			var tmp [3]int
			for i, col := range cols {
				tmp[i] = g[r][base+col]
			}
			for i := range 3 {
				g[r][base+i] = tmp[i]
			}
		}
	}

	// Step 4: shuffle bands.
	bandOrder := rng.Perm(3)
	var tmpGrid Grid
	for i, b := range bandOrder {
		for j := range 3 {
			tmpGrid[i*3+j] = g[b*3+j]
		}
	}
	g = tmpGrid

	// Step 5: shuffle stacks.
	stackOrder := rng.Perm(3)
	for r := range 9 {
		var tmpRow [9]int
		for i, s := range stackOrder {
			for j := range 3 {
				tmpRow[i*3+j] = g[r][s*3+j]
			}
		}
		g[r] = tmpRow
	}

	return g
}

// removeCells blanks out count cells from the solution to create the puzzle.
// Cells are chosen in a random order.
func removeCells(solution Grid, count int, rng *rand.Rand) Grid {
	puzzle := solution
	positions := rng.Perm(81)
	for _, pos := range positions[:count] {
		puzzle[pos/9][pos%9] = 0
	}
	return puzzle
}
