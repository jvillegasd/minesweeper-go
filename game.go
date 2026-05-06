package main

import (
	"github.com/minesweeper-go/enums"
	"math/rand/v2"
)

type Game struct {
	grid        [][]Tile
	state       enums.GameState
	totalMines  int
	totalTiles  int
	flagsPlaced int
	cursorI     int
	cursorJ     int
	minesPlaced bool
	difficulty  enums.Difficulty
}

func NewGame(d enums.Difficulty) *Game {
	level, ok := enums.Levels[d]
	if !ok {
		level = enums.Levels[enums.Beginner]
	}

	grid := make([][]Tile, level.Height)
	for i := range grid {
		grid[i] = make([]Tile, level.Width)
	}

	return &Game{
		grid:       grid,
		state:      enums.StatePlaying,
		totalMines: level.Mines,
		totalTiles: level.Height * level.Width,
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
	}
	g.revealTile(g.cursorI, g.cursorJ)
	if g.state == enums.StateLost {
		g.revealAllMines()
		return
	}
	if g.checkWin() {
		g.state = enums.StateWon
	}
}

func (g *Game) FlagAtCursor() {
	g.toggleFlag(g.cursorI, g.cursorJ)
}

func (g *Game) Restart() {
	*g = *NewGame(g.difficulty)
}
