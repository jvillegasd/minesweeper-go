package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/minesweeper-go/enums"
)

const cellWidth = 3

var numberColors = map[int]tcell.Color{
	1: tcell.ColorBlue,
	2: tcell.ColorGreen,
	3: tcell.ColorRed,
	4: tcell.ColorPurple,
	5: tcell.ColorMaroon,
	6: tcell.ColorTeal,
	7: tcell.ColorSilver,
	8: tcell.ColorGray,
}

func drawString(s tcell.Screen, x, y int, msg string, style tcell.Style) {
	for i, r := range []rune(msg) {
		s.SetContent(x+i, y, r, nil, style)
	}
}

func tileGlyph(t Tile, parity int) (rune, tcell.Style) {
	base := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)

	hiddenBg := tcell.ColorDarkSlateGray
	if parity == 1 {
		hiddenBg = tcell.ColorDimGray
	}

	switch {
	case !t.isRevealed && t.isFlagged:
		return 'F', base.Foreground(tcell.ColorRed).Background(hiddenBg)
	case !t.isRevealed:
		return ' ', base.Background(hiddenBg)
	case t.isMine:
		return '*', base.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)
	case t.adjMines == 0:
		return ' ', base
	default:
		style := base
		if c, ok := numberColors[t.adjMines]; ok {
			style = style.Foreground(c)
		}
		return rune('0' + t.adjMines), style
	}
}

func (g *Game) drawBoard(s tcell.Screen) {
	for i := range g.grid {
		for j := range g.grid[i] {
			glyph, style := tileGlyph(g.grid[i][j], (i+j)%2)
			if i == g.cursorI && j == g.cursorJ {
				style = style.Reverse(true)
			}

			x, y := j*cellWidth, i
			s.SetContent(x, y, ' ', nil, style)
			s.SetContent(x+1, y, glyph, nil, style)
			s.SetContent(x+2, y, ' ', nil, style)
		}
	}
}

func (g *Game) drawHUD(s tcell.Screen) {
	boardWidth := len(g.grid[0]) * cellWidth
	boardBottom := len(g.grid)

	var msg string
	msgStyle := tcell.StyleDefault.Background(tcell.ColorBlack)

	switch g.state {
	case enums.StateWon:
		msg = "you win! press r to restart, q to quit"
		msgStyle = msgStyle.Foreground(tcell.ColorGreen).Bold(true)
	case enums.StateLost:
		msg = "boom! press r to restart, q to quit"
		msgStyle = msgStyle.Foreground(tcell.ColorRed).Bold(true)
	case enums.StatePlaying:
		msg = "arrows: move  space/enter: reveal  f: flag  q: quit"
		msgStyle = msgStyle.Foreground(tcell.ColorSilver)
	}

	x := (boardWidth - len(msg)) / 2
	x = max(x, 0)
	y := boardBottom + 1
	drawString(s, x, y, msg, msgStyle)
}

func (g *Game) drawStatus(s tcell.Screen) {
	msg := fmt.Sprintf("flags: %d / %d", g.flagsPlaced, g.totalMines)

	fg := tcell.ColorYellow
	if g.flagsPlaced > g.totalMines {
		fg = tcell.ColorRed
	}
	msgStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(fg)

	boardWidth := len(g.grid[0]) * cellWidth
	boardBottom := len(g.grid)
	x := (boardWidth - len(msg)) / 2
	x = max(x, 0)
	drawString(s, x, boardBottom, msg, msgStyle)
}

func (g *Game) draw(s tcell.Screen) {
	s.Clear()
	g.drawBoard(s)
	g.drawStatus(s)
	g.drawHUD(s)
	s.Show()
}
