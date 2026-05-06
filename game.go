package main

import (
	"math/rand/v2"
	"time"
)

type Game struct {
	grid               [][]Tile
	state              GameState
	totalMines         int
	flagsPlaced        int
	cursorI            int
	cursorJ            int
	minesPlaced        bool
	difficulty         Difficulty
	startTime          time.Time
	endTime            time.Time
	requestMenu        bool
	requestLeaderboard bool
	lastWinSeconds     int
	lastWinRank        int
	lastWinIsBest      bool
}

func NewGame(d Difficulty) *Game {
	level, ok := Levels[d]
	if !ok {
		level = Levels[Beginner]
	}

	grid := make([][]Tile, level.Height)
	for i := range grid {
		grid[i] = make([]Tile, level.Width)
	}

	return &Game{
		grid:       grid,
		totalMines: level.Mines,
		difficulty: d,
	}
}

func (g *Game) placeMines(safeI, safeJ int) {
	height := len(g.grid)
	width := len(g.grid[0])

	inSafeZone := func(i, j int) bool {
		di := i - safeI
		dj := j - safeJ
		if di < 0 {
			di = -di
		}
		if dj < 0 {
			dj = -dj
		}
		return di <= 1 && dj <= 1
	}

	placed := 0
	for placed < g.totalMines {
		idx := rand.IntN(height * width)
		i, j := idx/width, idx%width
		if g.grid[i][j].isMine {
			continue
		}
		if inSafeZone(i, j) {
			continue
		}

		g.grid[i][j].isMine = true
		placed++
	}

	for i := range g.grid {
		for j := range g.grid[i] {
			if g.grid[i][j].isMine {
				continue
			}

			count := 0
			for _, d := range directions {
				ni, nj := i+d[0], j+d[1]
				if !g.isValidCoord(ni, nj) {
					continue
				}
				if g.grid[ni][nj].isMine {
					count++
				}
			}

			g.grid[i][j].adjMines = count
		}
	}

	g.minesPlaced = true
}

func (g *Game) MoveCursor(di, dj int) {
	ni := g.cursorI + di
	nj := g.cursorJ + dj
	if ni >= 0 && ni < len(g.grid) {
		g.cursorI = ni
	}
	if nj >= 0 && nj < len(g.grid[0]) {
		g.cursorJ = nj
	}
}

func (g *Game) RevealAtCursor() {
	if !g.minesPlaced {
		g.placeMines(g.cursorI, g.cursorJ)
		g.startTime = time.Now()
	}
	g.revealTile(g.cursorI, g.cursorJ)
	g.settleEndState()
}

func (g *Game) FlagAtCursor() {
	g.toggleFlag(g.cursorI, g.cursorJ)
}

func (g *Game) ChordAtCursor() {
	g.chordReveal(g.cursorI, g.cursorJ)
	g.settleEndState()
}

func (g *Game) settleEndState() {
	if g.state == StateLost {
		g.revealAllMines()
		if g.endTime.IsZero() {
			g.endTime = time.Now()
		}
		return
	}
	if g.checkWin() {
		g.state = StateWon
		if g.endTime.IsZero() {
			g.endTime = time.Now()
		}
		seconds := int(g.elapsed().Seconds())
		g.lastWinSeconds = seconds
		rank, isBest := leaderboard.record(g.difficulty, seconds)
		g.lastWinRank = rank
		g.lastWinIsBest = isBest
		saveLeaderboard(leaderboard)
	}
}

func (g *Game) elapsed() time.Duration {
	if g.startTime.IsZero() {
		return 0
	}
	end := g.endTime
	if end.IsZero() {
		end = time.Now()
	}
	return end.Sub(g.startTime)
}

func (g *Game) Restart() {
	*g = *NewGame(g.difficulty)
}

func (g *Game) SetDifficulty(d Difficulty) {
	*g = *NewGame(d)
}

func (g *Game) RequestMenu() {
	g.requestMenu = true
}

func (g *Game) RequestLeaderboard() {
	g.requestLeaderboard = true
}
