package sudoku

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerate_AllDifficulties verifies that Generate succeeds for every
// supported difficulty and returns a puzzle with the correct number of givens.
func TestGenerate_AllDifficulties(t *testing.T) {
	t.Parallel()

	cases := []struct {
		difficulty    Difficulty
		name          string
		expectedEmpty int
	}{
		{Easy, "easy", 36},
		{Medium, "medium", 46},
		{Hard, "hard", 52},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p, err := Generate(tc.difficulty)
			require.NoError(t, err)

			// Count empty cells in puzzle.
			empty := 0
			for r := range 9 {
				for c := range 9 {
					if p.Grid[r][c] == 0 {
						empty++
					}
				}
			}
			assert.Equal(t, tc.expectedEmpty, empty, "wrong number of empty cells")
		})
	}
}

// TestGenerate_SolutionIsValid verifies that the solution returned by Generate
// is a valid completed Sudoku for every difficulty.
func TestGenerate_SolutionIsValid(t *testing.T) {
	t.Parallel()

	for _, d := range []Difficulty{Easy, Medium, Hard} {
		p, err := Generate(d)
		require.NoError(t, err)
		assert.True(t, Validate(p.Solution), "solution must be a valid completed Sudoku")
	}
}

// TestGenerate_PuzzleMatchesSolution verifies that non-empty cells in the
// puzzle match the corresponding cells in the solution.
func TestGenerate_PuzzleMatchesSolution(t *testing.T) {
	t.Parallel()

	p, err := Generate(Easy)
	require.NoError(t, err)

	for r := range 9 {
		for c := range 9 {
			if p.Grid[r][c] != 0 {
				assert.Equal(t, p.Solution[r][c], p.Grid[r][c],
					"cell [%d][%d] in puzzle does not match solution", r, c)
			}
		}
	}
}

// TestGenerate_InvalidDifficulty verifies that an out-of-range Difficulty
// value returns an error and a zero Puzzle.
func TestGenerate_InvalidDifficulty(t *testing.T) {
	t.Parallel()

	p, err := Generate(Difficulty(99))
	require.Error(t, err)
	assert.Equal(t, Puzzle{}, p)
}

// TestValidate_ValidGrid verifies that Validate returns true for a known
// correct completed Sudoku.
func TestValidate_ValidGrid(t *testing.T) {
	t.Parallel()
	assert.True(t, Validate(base))
}

// TestValidate_DuplicateInRow verifies that a grid with a duplicate in a row
// is rejected.
func TestValidate_DuplicateInRow(t *testing.T) {
	t.Parallel()
	g := base
	g[0][0] = g[0][1] // duplicate in first row
	assert.False(t, Validate(g))
}

// TestValidate_DuplicateInColumn verifies that a grid with a duplicate in a
// column is rejected.
func TestValidate_DuplicateInColumn(t *testing.T) {
	t.Parallel()
	g := base
	g[1][0] = g[0][0] // duplicate in first column
	assert.False(t, Validate(g))
}

// TestValidate_DuplicateInBox verifies that a grid with a duplicate in a 3×3
// box is rejected.
func TestValidate_DuplicateInBox(t *testing.T) {
	t.Parallel()
	g := base
	g[2][2] = g[0][0] // duplicate in top-left box
	assert.False(t, Validate(g))
}

// TestValidate_ZeroCell verifies that a grid containing a 0 (empty cell) is
// rejected as incomplete.
func TestValidate_ZeroCell(t *testing.T) {
	t.Parallel()
	g := base
	g[0][0] = 0
	assert.False(t, Validate(g))
}

// TestShuffle_PreservesValidity verifies that the shuffle function always
// produces a valid completed Sudoku from the base grid.
func TestShuffle_PreservesValidity(t *testing.T) {
	t.Parallel()
	// Run several times to catch any non-deterministic failure.
	for range 20 {
		p, err := Generate(Easy)
		require.NoError(t, err)
		assert.True(t, Validate(p.Solution), "shuffled solution must be valid")
	}
}
