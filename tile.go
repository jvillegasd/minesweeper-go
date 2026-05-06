package main

type Tile struct {
	isMine     bool
	isRevealed bool
	isFlagged  bool
	adjMines   int
}

type coord struct {
	i, j int
}

var directions = [8][2]int{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
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

	if g.grid[i][j].isMine {
		g.grid[i][j].isRevealed = true
		g.state = StateLost
		return
	}

	queue := []coord{{i, j}}
	for len(queue) > 0 {
		c := queue[0]
		queue = queue[1:]

		if !g.isValidCoord(c.i, c.j) {
			continue
		}
		if g.grid[c.i][c.j].isFlagged || g.grid[c.i][c.j].isRevealed {
			continue
		}
		if g.grid[c.i][c.j].isMine {
			continue
		}

		g.grid[c.i][c.j].isRevealed = true

		if g.grid[c.i][c.j].adjMines == 0 {
			for _, d := range directions {
				queue = append(queue, coord{c.i + d[0], c.j + d[1]})
			}
		}
	}
}

func (g *Game) chordReveal(i, j int) {
	if !g.isValidCoord(i, j) {
		return
	}
	t := g.grid[i][j]
	if !t.isRevealed || t.adjMines == 0 {
		return
	}

	flagCount := 0
	for _, d := range directions {
		ni, nj := i+d[0], j+d[1]
		if !g.isValidCoord(ni, nj) {
			continue
		}
		if g.grid[ni][nj].isFlagged {
			flagCount++
		}
	}
	if flagCount != t.adjMines {
		return
	}

	for _, d := range directions {
		ni, nj := i+d[0], j+d[1]
		if !g.isValidCoord(ni, nj) {
			continue
		}
		if g.grid[ni][nj].isFlagged || g.grid[ni][nj].isRevealed {
			continue
		}
		g.revealTile(ni, nj)
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

	if g.grid[i][j].isFlagged {
		g.flagsPlaced++
	} else {
		g.flagsPlaced--
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
