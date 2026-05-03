package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/minesweeper-go/enums"
)


type Tile struct {
	isMine bool
	isRevealed bool
	isFlagged bool
	isHover bool
	adjMines int
}

type Game struct {
	grid [][]Tile
	state enums.GameState
	totalMines int
	totalTiles int
}

var directions = [8][2]int{
	{-1, -1}, {-1, 0}, {-1, 1},
	{ 0, -1},          { 0, 1},
	{ 1, -1}, { 1, 0}, { 1, 1}, 
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
	height := len(g.grid)
	if height == 0 {
		return
	}

	width := len(g.grid[0])
	if i < 0 || i >= height || j < 0 || j >= width {
		return
	}

	if g.grid[i][j].isRevealed {
		return
	}
	
	g.grid[i][j].isRevealed = true
	if g.grid[i][j].isMine {
		g.state = enums.StateLost
		return
	}

	if g.grid[i][j].adjMines == 0 {
		for _, d := range directions {
			ni, nj := i + d[0], j + d[1]
			g.revealTile(ni, nj)
		}
	}
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

func main() {
	fmt.Println("Hello, World!")
}
