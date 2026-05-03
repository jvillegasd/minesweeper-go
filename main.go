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
			ni, nj := i+d[0], j+d[1]
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
		grid: grid,
		state: enums.StatePlaying,
		totalMines: level.Mines,
		totalTiles: totalTiles,
	}
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Println("failed to create screen:", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Println("failed to init screen:", err)
		os.Exit(1)
	}
	defer screen.Fini()

	game := NewGame(enums.Beginner)
	_ = game
}
