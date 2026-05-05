package main

import (
	"math/rand/v2"

	"github.com/minesweeper-go/enums"
)

type Game struct {
	grid        [][]Tile
	state       enums.GameState
	totalMines  int
	totalTiles  int
	flagsPlaced int
	cursorI     int
	cursorJ     int
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

	placed := 0
	totalTiles := level.Height * level.Width
	for placed < level.Mines {
		idx := rand.IntN(totalTiles)
		i, j := idx/level.Width, idx%level.Width
		if grid[i][j].isMine {
			continue
		}

		grid[i][j].isMine = true
		placed++
	}

	for i := range grid {
		for j := range grid[i] {
			if grid[i][j].isMine {
				continue
			}

			count := 0
			for _, d := range directions {
				ni, nj := i+d[0], j+d[1]
				if ni < 0 || ni >= level.Height || nj < 0 || nj >= level.Width {
					continue
				}
				if grid[ni][nj].isMine {
					count++
				}
			}

			grid[i][j].adjMines = count
		}
	}

	return &Game{
		grid:       grid,
		state:      enums.StatePlaying,
		totalMines: level.Mines,
		totalTiles: totalTiles,
		cursorI:    0,
		cursorJ:    0,
		difficulty: d,
	}
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
