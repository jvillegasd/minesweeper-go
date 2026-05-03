package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
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
	state string
	totalMines int
	totalTiles int
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
	if r < 0 || r > height || c < width || c > width {
		return
	}

	if g.grid[r][c].isRevealed {
		return
	}


}


func main() {
	fmt.Println("Hello, World!")
}
