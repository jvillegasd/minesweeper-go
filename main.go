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
	for r := range g.grid {
		for c := range g.grid[r] {
			if g.grid[r][c].isMine {
				g.grid[r][c].isRevealed = true
			}
		}
	}
}


func (g *Game) revealTile(r, c int) {
	height := len(g.grid)
	if height == 0 {
		return
	}

	width := len(g.grid[0])
	if r < 0 || r >= height || c < 0 || c >= width {
		return
	}

	if g.grid[r][c].isRevealed {
		return
	}
	
	g.grid[r][c].isRevealed = true
	if g.grid[r][c].isMine {
		g.state = enums.StateLost
		return
	}

	if g.grid[r][c].adjMines == 0 {
		for _, d := range directions {
			nr, nc := r + d[0], c + d[1]
			g.revealTile(nr, nc)
		}
	}
}


func main() {
	fmt.Println("Hello, World!")
}
