package main

import "github.com/minesweeper-go/enums"

type Tile struct {
	isMine bool
	isRevealed bool
	isFlagged bool
	adjMines int
}

var directions = [8][2]int{
	{-1, -1}, {-1, 0}, {-1, 1},
	{ 0, -1},          { 0, 1},
	{ 1, -1}, { 1, 0}, { 1, 1},
}

func (g *Game) isValidCoord(i, j int) bool {
	height := len(g.grid)
	if height == 0 {
		return false
	}

	width := len(g.grid[0])
	if i < 0 || i >= height || j < 0 || j >= width {
		return false
	}

	return true
}

func (g *Game) revealAllMines() {
	for i := range g.grid {
		for j := range g.grid[i] {
			if g.grid[i][j].isMine {
				g.grid[i][j].isRevealed = true
			}
		}
	}
}

func (g *Game) revealTile(i, j int) {
	if !g.isValidCoord(i, j) {
		return
	}

	if g.grid[i][j].isFlagged || g.grid[i][j].isRevealed {
		return
	}

	g.grid[i][j].isRevealed = true
	if g.grid[i][j].isMine {
		g.state = enums.StateLost
		return
	}

	if g.grid[i][j].adjMines == 0 {
		for _, d := range directions {
			ni, nj := i+d[0], j+d[1]
			g.revealTile(ni, nj)
		}
	}
}

func (g *Game) toggleFlag(i, j int) {
	if !g.isValidCoord(i, j) {
		return
	}

	if g.grid[i][j].isRevealed {
		return
	}

	g.grid[i][j].isFlagged = !g.grid[i][j].isFlagged
}

func (g *Game) checkWin() bool {
	for i := range g.grid {
		for j := range g.grid[i] {
			if !g.grid[i][j].isMine && !g.grid[i][j].isRevealed {
				return false
			}
		}
	}
	return true
}
