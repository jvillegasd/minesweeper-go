package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

const (
	cellWidth = 3
	belowRows = 4 // gap + separator + status + hud
)

var base = tcell.StyleDefault.Background(tcell.ColorBlack)

var (
	boardBg     = tcell.NewRGBColor(22, 24, 38)
	tileBgDark  = tcell.NewRGBColor(72, 130, 88)
	tileBgLight = tcell.NewRGBColor(108, 170, 120)
	mineBg      = tcell.NewRGBColor(180, 30, 30)
	flagFg      = tcell.NewRGBColor(255, 100, 100)

	accentYellow = tcell.ColorYellow
	accentGreen  = tcell.ColorLightGreen
	accentRed    = tcell.ColorRed
	textSilver   = tcell.ColorSilver
	textDim      = tcell.NewRGBColor(120, 120, 140)
)

var numberPalette = map[int]tcell.Color{
	1: tcell.NewRGBColor(100, 160, 255),
	2: tcell.NewRGBColor(100, 220, 110),
	3: tcell.NewRGBColor(255, 110, 110),
	4: tcell.NewRGBColor(190, 120, 230),
	5: tcell.NewRGBColor(230, 160, 60),
	6: tcell.NewRGBColor(80, 220, 200),
	7: tcell.ColorSilver,
	8: tcell.ColorWhite,
}

func drawString(s tcell.Screen, x, y int, msg string, style tcell.Style) {
	for i, r := range []rune(msg) {
		s.SetContent(x+i, y, r, nil, style)
	}
}

func tileGlyph(t Tile, parity int) (rune, tcell.Style) {
	hiddenBg := tileBgDark
	if parity == 1 {
		hiddenBg = tileBgLight
	}

	revealed := base.Background(boardBg)

	switch {
	case !t.isRevealed && t.isFlagged:
		return '⚑', base.Foreground(flagFg).Background(hiddenBg).Bold(true)
	case !t.isRevealed:
		return ' ', base.Background(hiddenBg)
	case t.isMine:
		return '✸', base.Foreground(tcell.ColorWhite).Background(mineBg).Bold(true)
	case t.adjMines == 0:
		return ' ', revealed
	default:
		style := revealed.Bold(true)
		if c, ok := numberPalette[t.adjMines]; ok {
			style = style.Foreground(c)
		}
		return rune('0' + t.adjMines), style
	}
}

func (g *Game) drawBoard(s tcell.Screen, offsetX, offsetY int) {
	for i := range g.grid {
		for j := range g.grid[i] {
			glyph, style := tileGlyph(g.grid[i][j], (i+j)%2)
			if i == g.cursorI && j == g.cursorJ {
				style = style.Reverse(true)
			}

			x := offsetX + j*cellWidth
			y := offsetY + i
			s.SetContent(x, y, ' ', nil, style)
			s.SetContent(x+1, y, glyph, nil, style)
			s.SetContent(x+2, y, ' ', nil, style)
		}
	}
}

func (g *Game) drawSeparator(s tcell.Screen, sw, y int) {
	boardW := len(g.grid[0]) * cellWidth
	width := boardW
	if width < 40 {
		width = 40
	}
	if width > sw {
		width = sw
	}
	x := (sw - width) / 2
	drawString(s, x, y, strings.Repeat("━", width), base.Foreground(textDim))
}

func centeredX(sw, segWidth int) int {
	x := (sw - segWidth) / 2
	if x < 0 {
		return 0
	}
	return x
}

func segmentsWidth(segs []segment) int {
	n := 0
	for _, seg := range segs {
		n += len([]rune(seg.text))
	}
	return n
}

func (g *Game) drawStatus(s tcell.Screen, sw, y int) {
	flagColor := accentYellow
	if g.flagsPlaced > g.totalMines {
		flagColor = accentRed
	}

	dot := segment{"  ·  ", base.Foreground(textDim)}
	seconds := int(g.elapsed().Seconds())

	segs := []segment{
		{"⚑ ", base.Foreground(flagColor)},
		{fmt.Sprintf("%d", g.flagsPlaced), base.Foreground(flagColor).Bold(true)},
		{fmt.Sprintf("/%d", g.totalMines), base.Foreground(textSilver)},
		dot,
		{g.difficulty.String(), base.Foreground(textSilver).Bold(true)},
		dot,
		{fmt.Sprintf("%ds", seconds), base.Foreground(textSilver)},
	}
	drawSegments(s, centeredX(sw, segmentsWidth(segs)), y, segs)
}

func (g *Game) drawHUD(s tcell.Screen, sw, y int) {
	var segs []segment
	switch g.state {
	case StateWon:
		segs = []segment{
			{"YOU WIN", base.Foreground(accentGreen).Bold(true)},
			{"  ·  press ", base.Foreground(textSilver)},
			{"r", base.Foreground(accentYellow).Bold(true)},
			{" to restart  ·  ", base.Foreground(textSilver)},
			{"q", base.Foreground(accentYellow).Bold(true)},
			{" to quit", base.Foreground(textSilver)},
		}
	case StateLost:
		segs = []segment{
			{"BOOM", base.Foreground(accentRed).Bold(true)},
			{"  ·  press ", base.Foreground(textSilver)},
			{"r", base.Foreground(accentYellow).Bold(true)},
			{" to restart  ·  ", base.Foreground(textSilver)},
			{"q", base.Foreground(accentYellow).Bold(true)},
			{" to quit", base.Foreground(textSilver)},
		}
	case StatePlaying:
		segs = []segment{
			{"arrows", base.Foreground(accentYellow).Bold(true)},
			{":move  ", base.Foreground(textSilver)},
			{"space", base.Foreground(accentYellow).Bold(true)},
			{":reveal  ", base.Foreground(textSilver)},
			{"f", base.Foreground(accentYellow).Bold(true)},
			{":flag  ", base.Foreground(textSilver)},
			{"c", base.Foreground(accentYellow).Bold(true)},
			{":chord  ", base.Foreground(textSilver)},
			{"q", base.Foreground(accentYellow).Bold(true)},
			{":quit", base.Foreground(textSilver)},
		}
	}
	drawSegments(s, centeredX(sw, segmentsWidth(segs)), y, segs)
}

func (g *Game) draw(s tcell.Screen) {
	s.Clear()

	sw, sh := s.Size()
	boardW := len(g.grid[0]) * cellWidth
	boardH := len(g.grid)
	needW := boardW
	needH := boardH + belowRows

	if sw < needW || sh < needH {
		msg := fmt.Sprintf("terminal too small — need %d×%d, got %d×%d",
			needW, needH, sw, sh)
		drawString(s, 0, 0, msg, base.Foreground(accentRed).Bold(true))
		s.Show()
		return
	}

	offsetX := (sw - boardW) / 2
	offsetY := (sh - needH) / 2
	if offsetY < 0 {
		offsetY = 0
	}

	g.drawBoard(s, offsetX, offsetY)

	sepY := offsetY + boardH + 1
	g.drawSeparator(s, sw, sepY)
	g.drawStatus(s, sw, sepY+1)
	g.drawHUD(s, sw, sepY+2)

	s.Show()
}
