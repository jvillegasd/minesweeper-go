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
	isClicked bool
	isFlagged bool
	isHover bool
	adjMines int
}

type Grid struct {
	grid [][]Tile
}

var gameState string
var totalMines, totalTiles int


func main() {
	fmt.Println("Hello, World!")
}
