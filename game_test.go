package main

import "testing"

func TestZeroValueOfGameState(t *testing.T) {
	var s GameState
	if s != StatePlaying {
		t.Errorf("zero value should be StatePlaying, got %v", s)
	}
}

func TestNewGameDimensions(t *testing.T) {
	cases := []struct {
		d             Difficulty
		height, width int
		mines         int
	}{
		{Beginner, 9, 9, 10},
		{Intermediate, 16, 16, 40},
		{Expert, 16, 30, 99},
	}
	for _, c := range cases {
		g := NewGame(c.d)
		if len(g.grid) != c.height {
			t.Errorf("%s: height %d, want %d", c.d, len(g.grid), c.height)
		}
		if len(g.grid[0]) != c.width {
			t.Errorf("%s: width %d, want %d", c.d, len(g.grid[0]), c.width)
		}
		if g.totalMines != c.mines {
			t.Errorf("%s: totalMines %d, want %d", c.d, g.totalMines, c.mines)
		}
		if g.minesPlaced {
			t.Errorf("%s: minesPlaced should start false", c.d)
		}
	}
}

func TestPlaceMinesCount(t *testing.T) {
	g := NewGame(Intermediate)
	g.placeMines(8, 8)

	count := 0
	for i := range g.grid {
		for j := range g.grid[i] {
			if g.grid[i][j].isMine {
				count++
			}
		}
	}
	if count != g.totalMines {
		t.Errorf("placed %d mines, want %d", count, g.totalMines)
	}
	if !g.minesPlaced {
		t.Error("minesPlaced should be true after placeMines")
	}
}

func TestPlaceMinesSafeZone(t *testing.T) {
	for trial := 0; trial < 50; trial++ {
		g := NewGame(Intermediate)
		safeI, safeJ := 5, 5
		g.placeMines(safeI, safeJ)

		for di := -1; di <= 1; di++ {
			for dj := -1; dj <= 1; dj++ {
				if g.grid[safeI+di][safeJ+dj].isMine {
					t.Errorf("trial %d: mine inside safe zone at (%d,%d)",
						trial, safeI+di, safeJ+dj)
				}
			}
		}
	}
}

func TestPlaceMinesAdjacency(t *testing.T) {
	g := NewGame(Beginner)
	g.placeMines(4, 4)

	for i := range g.grid {
		for j := range g.grid[i] {
			if g.grid[i][j].isMine {
				continue
			}
			expected := 0
			for _, d := range directions {
				ni, nj := i+d[0], j+d[1]
				if !g.isValidCoord(ni, nj) {
					continue
				}
				if g.grid[ni][nj].isMine {
					expected++
				}
			}
			if g.grid[i][j].adjMines != expected {
				t.Errorf("(%d,%d): adjMines %d, want %d",
					i, j, g.grid[i][j].adjMines, expected)
			}
		}
	}
}

func TestFirstClickSafety(t *testing.T) {
	for trial := 0; trial < 100; trial++ {
		g := NewGame(Beginner)
		g.cursorI, g.cursorJ = 0, 0
		g.RevealAtCursor()
		if g.grid[0][0].isMine {
			t.Errorf("trial %d: first click landed on a mine", trial)
		}
		if g.state == StateLost {
			t.Errorf("trial %d: first click ended game", trial)
		}
	}
}

func TestRevealTileMineLost(t *testing.T) {
	g := NewGame(Beginner)
	g.grid[0][0].isMine = true
	g.minesPlaced = true

	g.revealTile(0, 0)

	if g.state != StateLost {
		t.Errorf("state %v, want StateLost", g.state)
	}
	if !g.grid[0][0].isRevealed {
		t.Error("clicked mine should be revealed")
	}
}

func TestRevealTileFloodFromZero(t *testing.T) {
	// Construct a 3x3 grid with all zero-tiles (no mines).
	g := &Game{
		grid: make([][]Tile, 3),
	}
	for i := range g.grid {
		g.grid[i] = make([]Tile, 3)
	}
	g.minesPlaced = true

	g.revealTile(1, 1)

	for i := range g.grid {
		for j := range g.grid[i] {
			if !g.grid[i][j].isRevealed {
				t.Errorf("(%d,%d) should be revealed by flood", i, j)
			}
		}
	}
}

func TestRevealTileFloodStopsAtNumbers(t *testing.T) {
	// 3x3: mine at (0,2). (0,1) and (1,2) become numbers.
	// Click (2,0): zero, flood expands, stops at numbers.
	g := &Game{
		grid: make([][]Tile, 3),
	}
	for i := range g.grid {
		g.grid[i] = make([]Tile, 3)
	}
	g.minesPlaced = true

	g.grid[0][2].isMine = true
	g.grid[0][1].adjMines = 1
	g.grid[1][1].adjMines = 1
	g.grid[1][2].adjMines = 1

	g.revealTile(2, 0)

	if g.grid[0][2].isRevealed {
		t.Error("mine should not be revealed by flood")
	}
	if !g.grid[0][1].isRevealed || !g.grid[1][1].isRevealed || !g.grid[1][2].isRevealed {
		t.Error("number-border tiles should be revealed")
	}
	if !g.grid[0][0].isRevealed || !g.grid[2][2].isRevealed {
		t.Error("zero-tiles should be revealed")
	}
}

func TestRevealTileSkipsFlagged(t *testing.T) {
	g := NewGame(Beginner)
	g.minesPlaced = true
	g.grid[0][0].isFlagged = true

	g.revealTile(0, 0)

	if g.grid[0][0].isRevealed {
		t.Error("flagged tile must not be revealed")
	}
}

func TestCheckWinFreshFalse(t *testing.T) {
	g := NewGame(Beginner)
	g.placeMines(0, 0)

	if g.checkWin() {
		t.Error("fresh game should not be won")
	}
}

func TestCheckWinAllNonMinesRevealed(t *testing.T) {
	g := NewGame(Beginner)
	g.placeMines(0, 0)
	for i := range g.grid {
		for j := range g.grid[i] {
			if !g.grid[i][j].isMine {
				g.grid[i][j].isRevealed = true
			}
		}
	}

	if !g.checkWin() {
		t.Error("all non-mines revealed should win")
	}
}

func TestToggleFlagCount(t *testing.T) {
	g := NewGame(Beginner)

	g.toggleFlag(0, 0)
	if g.flagsPlaced != 1 || !g.grid[0][0].isFlagged {
		t.Errorf("after flag: flagsPlaced=%d, flagged=%v", g.flagsPlaced, g.grid[0][0].isFlagged)
	}

	g.toggleFlag(0, 0)
	if g.flagsPlaced != 0 || g.grid[0][0].isFlagged {
		t.Errorf("after unflag: flagsPlaced=%d, flagged=%v", g.flagsPlaced, g.grid[0][0].isFlagged)
	}
}

func TestToggleFlagOnRevealedNoop(t *testing.T) {
	g := NewGame(Beginner)
	g.grid[0][0].isRevealed = true

	g.toggleFlag(0, 0)

	if g.grid[0][0].isFlagged || g.flagsPlaced != 0 {
		t.Error("toggle on revealed tile must be a no-op")
	}
}

func TestChordRevealRevealsNeighbors(t *testing.T) {
	// 3x3 with mine at (1,2), correctly flagged.
	// Center (1,1) is a "1". Chord on (1,1) should reveal all unflagged neighbors.
	g := &Game{
		grid: make([][]Tile, 3),
	}
	for i := range g.grid {
		g.grid[i] = make([]Tile, 3)
	}
	g.minesPlaced = true

	g.grid[1][2].isMine = true
	g.grid[1][2].isFlagged = true
	g.flagsPlaced = 1

	g.grid[0][1].adjMines = 1
	g.grid[0][2].adjMines = 1
	g.grid[1][1].adjMines = 1
	g.grid[1][1].isRevealed = true
	g.grid[2][1].adjMines = 1
	g.grid[2][2].adjMines = 1

	g.chordReveal(1, 1)

	for _, p := range []coord{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {2, 0}, {2, 1}, {2, 2}} {
		if !g.grid[p.i][p.j].isRevealed {
			t.Errorf("(%d,%d) should be revealed by chord", p.i, p.j)
		}
	}
	if g.grid[1][2].isRevealed {
		t.Error("flagged mine must not be revealed by chord")
	}
}

func TestChordRevealNoopWhenFlagsMismatch(t *testing.T) {
	g := &Game{
		grid: make([][]Tile, 3),
	}
	for i := range g.grid {
		g.grid[i] = make([]Tile, 3)
	}
	g.minesPlaced = true

	g.grid[1][2].isMine = true
	g.grid[1][1].adjMines = 1
	g.grid[1][1].isRevealed = true

	g.chordReveal(1, 1)

	for i := range g.grid {
		for j := range g.grid[i] {
			if i == 1 && j == 1 {
				continue
			}
			if g.grid[i][j].isRevealed {
				t.Errorf("(%d,%d) should not be revealed when flag count mismatches", i, j)
			}
		}
	}
}

func TestRestartResetsState(t *testing.T) {
	g := NewGame(Beginner)
	g.RevealAtCursor()
	g.toggleFlag(5, 5)

	g.Restart()

	if g.minesPlaced {
		t.Error("minesPlaced should reset to false")
	}
	if g.flagsPlaced != 0 {
		t.Errorf("flagsPlaced %d, want 0", g.flagsPlaced)
	}
	if g.state != StatePlaying {
		t.Errorf("state %v, want StatePlaying", g.state)
	}
	if !g.startTime.IsZero() {
		t.Error("startTime should reset")
	}
}

func TestSetDifficultySwitches(t *testing.T) {
	g := NewGame(Beginner)
	g.SetDifficulty(Expert)

	if g.difficulty != Expert {
		t.Errorf("difficulty %v, want Expert", g.difficulty)
	}
	if len(g.grid) != 16 || len(g.grid[0]) != 30 {
		t.Errorf("grid %dx%d, want 16x30", len(g.grid), len(g.grid[0]))
	}
}

func TestMoveCursorClampsToGrid(t *testing.T) {
	g := NewGame(Beginner)

	g.MoveCursor(-1, 0)
	if g.cursorI != 0 {
		t.Errorf("cursorI %d, want 0 (clamped)", g.cursorI)
	}

	g.cursorI, g.cursorJ = 8, 8
	g.MoveCursor(1, 1)
	if g.cursorI != 8 || g.cursorJ != 8 {
		t.Errorf("cursor (%d,%d), want (8,8) (clamped)", g.cursorI, g.cursorJ)
	}
}
