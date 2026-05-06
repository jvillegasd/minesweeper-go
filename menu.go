package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

const (
	bigOptionsW = 41
	compactMinW = 19
	compactMinH = 10
)

type segment struct {
	text  string
	style tcell.Style
}

type menuRow struct {
	segments []segment
}

func (r menuRow) width() int {
	n := 0
	for _, s := range r.segments {
		n += len([]rune(s.text))
	}
	return n
}

func plainRow(text string, style tcell.Style) menuRow {
	return menuRow{segments: []segment{{text, style}}}
}

func drawSegments(s tcell.Screen, x, y int, segs []segment) {
	cursor := x
	for _, seg := range segs {
		for _, r := range []rune(seg.text) {
			s.SetContent(cursor, y, r, nil, seg.style)
			cursor++
		}
	}
}

func menuRows(sw int) []menuRow {
	bigOptions := sw >= bigOptionsW

	titleStyle := base.Foreground(tcell.NewRGBColor(255, 100, 100)).Bold(true)

	tagline := base.Foreground(textSilver).Italic(true)
	header := base.Foreground(accentYellow).Bold(true)
	keyBracket := base.Foreground(accentYellow).Bold(true)
	silver := base.Foreground(textSilver)
	dim := base.Foreground(textDim)
	green := base.Foreground(accentGreen).Bold(true)
	amber := base.Foreground(accentYellow).Bold(true)
	red := base.Foreground(accentRed).Bold(true)

	rows := make([]menuRow, 0, 16)

	rows = append(rows,
		menuRow{segments: []segment{
			{"✦ ", titleStyle},
			{"MINESWEEPER", titleStyle},
			{" ✦", titleStyle},
		}},
	)

	rows = append(rows,
		plainRow("", silver),
		plainRow("✦  built in Go  ✦", tagline),
		plainRow("", silver),
	)

	sepWidth := 60
	if sepWidth > sw-4 {
		sepWidth = sw - 4
	}
	if sepWidth < 10 {
		sepWidth = 10
	}
	rows = append(rows,
		plainRow(strings.Repeat("━", sepWidth), dim),
		plainRow("", silver),
		plainRow("SELECT A DIFFICULTY", header),
		plainRow("", silver),
	)

	bestText := func(d Difficulty) (string, tcell.Style) {
		if best, ok := leaderboard.best(d); ok {
			return fmt.Sprintf("  ·  best %ds", best), base.Foreground(accentGreen)
		}
		return "", silver
	}

	if bigOptions {
		bBest, bBestStyle := bestText(Beginner)
		iBest, iBestStyle := bestText(Intermediate)
		eBest, eBestStyle := bestText(Expert)

		rows = append(rows,
			menuRow{segments: []segment{
				{"[ 1 ]  ", keyBracket},
				{"Beginner       ", green},
				{"·  9 × 9    ·  10 mines", silver},
				{bBest, bBestStyle},
			}},
			menuRow{segments: []segment{
				{"[ 2 ]  ", keyBracket},
				{"Intermediate   ", amber},
				{"·  16 × 16  ·  40 mines", silver},
				{iBest, iBestStyle},
			}},
			menuRow{segments: []segment{
				{"[ 3 ]  ", keyBracket},
				{"Expert         ", red},
				{"·  30 × 16  ·  99 mines", silver},
				{eBest, eBestStyle},
			}},
		)
	} else {
		rows = append(rows,
			menuRow{segments: []segment{
				{"[1] ", keyBracket},
				{"Beginner", green},
			}},
			menuRow{segments: []segment{
				{"[2] ", keyBracket},
				{"Intermediate", amber},
			}},
			menuRow{segments: []segment{
				{"[3] ", keyBracket},
				{"Expert", red},
			}},
		)
	}

	rows = append(rows,
		plainRow("", silver),
		plainRow(strings.Repeat("━", sepWidth), dim),
		plainRow("", silver),
		menuRow{segments: []segment{
			{"[ L ]  ", keyBracket},
			{"Leaderboard", silver},
		}},
		menuRow{segments: []segment{
			{"[ Q ]  ", keyBracket},
			{"Quit", silver},
		}},
	)
	return rows
}

func drawMenu(s tcell.Screen) {
	s.Clear()
	sw, sh := s.Size()

	if sw < compactMinW || sh < compactMinH {
		style := base.Foreground(accentRed).Bold(true)
		drawString(s, 0, 0, "terminal too small — resize", style)
		s.Show()
		return
	}

	rows := menuRows(sw)
	needH := len(rows)
	startY := (sh - needH) / 2
	if startY < 0 {
		startY = 0
	}

	for i, r := range rows {
		x := (sw - r.width()) / 2
		if x < 0 {
			x = 0
		}
		drawSegments(s, x, startY+i, r.segments)
	}
	s.Show()
}

func runMenu(screen tcell.Screen) (Difficulty, bool) {
	drawMenu(screen)
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				return 0, false
			}
			switch ev.Rune() {
			case 'q', 'Q':
				return 0, false
			case '1':
				return Beginner, true
			case '2':
				return Intermediate, true
			case '3':
				return Expert, true
			case 'l', 'L':
				runLeaderboard(screen)
				drawMenu(screen)
			}
		case *tcell.EventResize:
			screen.Sync()
			drawMenu(screen)
		}
	}
}
