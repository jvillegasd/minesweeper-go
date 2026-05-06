package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

const (
	cellWidth  = 3
	belowRows  = 4 // gap + separator + status + hud
	borderRows = 2 // top + bottom border
	borderCols = 2 // left + right border
)

var base = tcell.StyleDefault

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

func drawBorder(s tcell.Screen, x, y, w, h int, style tcell.Style) {
	s.SetContent(x, y, '╭', nil, style)
	s.SetContent(x+w-1, y, '╮', nil, style)
	s.SetContent(x, y+h-1, '╰', nil, style)
	s.SetContent(x+w-1, y+h-1, '╯', nil, style)

	for i := 1; i < w-1; i++ {
		s.SetContent(x+i, y, '─', nil, style)
		s.SetContent(x+i, y+h-1, '─', nil, style)
	}
	for i := 1; i < h-1; i++ {
		s.SetContent(x, y+i, '│', nil, style)
		s.SetContent(x+w-1, y+i, '│', nil, style)
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

	if best, ok := leaderboard.best(g.difficulty); ok {
		segs = append(segs,
			dot,
			segment{"best ", base.Foreground(textDim)},
			segment{fmt.Sprintf("%ds", best), base.Foreground(accentGreen).Bold(true)},
		)
	}
	drawSegments(s, centeredX(sw, segmentsWidth(segs)), y, segs)
}

func (g *Game) drawHUD(s tcell.Screen, sw, y int) {
	yellow := base.Foreground(accentYellow).Bold(true)
	silver := base.Foreground(textSilver)

	var segs []segment
	switch g.state {
	case StateWon:
		segs = []segment{
			{"YOU WIN", base.Foreground(accentGreen).Bold(true)},
			{" in ", silver},
			{fmt.Sprintf("%ds", g.lastWinSeconds), base.Foreground(accentGreen).Bold(true)},
		}
		if g.lastWinIsBest {
			segs = append(segs, segment{"  ★ NEW BEST", base.Foreground(accentYellow).Bold(true)})
		} else if g.lastWinRank > 0 {
			segs = append(segs, segment{
				fmt.Sprintf("  ·  rank #%d", g.lastWinRank),
				silver,
			})
		}
		segs = append(segs,
			segment{"  ·  ", base.Foreground(textDim)},
			segment{"r", yellow}, segment{":restart  ", silver},
			segment{"m", yellow}, segment{":menu  ", silver},
			segment{"l", yellow}, segment{":leaderboard  ", silver},
			segment{"q", yellow}, segment{":quit", silver},
		)
	case StateLost:
		segs = []segment{
			{"BOOM", base.Foreground(accentRed).Bold(true)},
			{"  ·  ", base.Foreground(textDim)},
			{"r", yellow}, {":restart  ", silver},
			{"m", yellow}, {":menu  ", silver},
			{"l", yellow}, {":leaderboard  ", silver},
			{"q", yellow}, {":quit", silver},
		}
	case StatePlaying:
		segs = []segment{
			{"arrows", yellow}, {":move  ", silver},
			{"space", yellow}, {":reveal  ", silver},
			{"f", yellow}, {":flag  ", silver},
			{"c", yellow}, {":chord  ", silver},
			{"m", yellow}, {":menu  ", silver},
			{"q", yellow}, {":quit", silver},
		}
	}
	drawSegments(s, centeredX(sw, segmentsWidth(segs)), y, segs)
}

func (g *Game) draw(s tcell.Screen) {
	s.Clear()

	sw, sh := s.Size()
	boardW := len(g.grid[0]) * cellWidth
	boardH := len(g.grid)
	needW := boardW + borderCols
	needH := boardH + borderRows + belowRows

	if sw < needW || sh < needH {
		msg := fmt.Sprintf("terminal too small — need %d×%d, got %d×%d",
			needW, needH, sw, sh)
		drawString(s, 0, 0, msg, base.Foreground(accentRed).Bold(true))
		s.Show()
		return
	}

	boardX := (sw - boardW) / 2
	boardY := (sh-needH)/2 + 1
	if boardY < 1 {
		boardY = 1
	}

	drawBorder(s, boardX-1, boardY-1, boardW+borderCols, boardH+borderRows,
		base.Foreground(tcell.ColorWhite))
	g.drawBoard(s, boardX, boardY)

	sepY := boardY + boardH + 2
	g.drawSeparator(s, sw, sepY)
	g.drawStatus(s, sw, sepY+1)
	g.drawHUD(s, sw, sepY+2)

	s.Show()
}
