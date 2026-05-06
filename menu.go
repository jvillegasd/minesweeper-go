package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

var titleArt = []string{
	"███╗   ███╗██╗███╗   ██╗███████╗███████╗██╗    ██╗███████╗███████╗██████╗ ███████╗██████╗ ",
	"████╗ ████║██║████╗  ██║██╔════╝██╔════╝██║    ██║██╔════╝██╔════╝██╔══██╗██╔════╝██╔══██╗",
	"██╔████╔██║██║██╔██╗ ██║█████╗  ███████╗██║ █╗ ██║█████╗  █████╗  ██████╔╝█████╗  ██████╔╝",
	"██║╚██╔╝██║██║██║╚██╗██║██╔══╝  ╚════██║██║███╗██║██╔══╝  ██╔══╝  ██╔═══╝ ██╔══╝  ██╔══██╗",
	"██║ ╚═╝ ██║██║██║ ╚████║███████╗███████║╚███╔███╔╝███████╗███████╗██║     ███████╗██║  ██║",
	"╚═╝     ╚═╝╚═╝╚═╝  ╚═╝╚══════╝╚══════╝ ╚══╝╚══╝ ╚══════╝╚══════╝╚═╝     ╚══════╝╚═╝  ╚═╝",
}

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

func menuRows() []menuRow {
	titleStyles := []tcell.Style{
		base.Foreground(tcell.NewRGBColor(255, 100, 100)).Bold(true),
		base.Foreground(tcell.NewRGBColor(235, 80, 80)).Bold(true),
		base.Foreground(tcell.NewRGBColor(215, 60, 60)).Bold(true),
		base.Foreground(tcell.NewRGBColor(190, 45, 45)).Bold(true),
		base.Foreground(tcell.NewRGBColor(165, 30, 30)).Bold(true),
		base.Foreground(tcell.NewRGBColor(140, 20, 20)).Bold(true),
	}

	tagline := base.Foreground(textSilver).Italic(true)
	header := base.Foreground(accentYellow).Bold(true)
	keyBracket := base.Foreground(accentYellow).Bold(true)
	silver := base.Foreground(textSilver)
	dim := base.Foreground(textDim)
	green := base.Foreground(accentGreen).Bold(true)
	amber := base.Foreground(accentYellow).Bold(true)
	red := base.Foreground(accentRed).Bold(true)

	separator := strings.Repeat("━", 60)

	rows := make([]menuRow, 0, len(titleArt)+16)
	for i, line := range titleArt {
		rows = append(rows, plainRow(line, titleStyles[i]))
	}

	rows = append(rows,
		plainRow("", silver),
		plainRow("✦  built in Go  ✦", tagline),
		plainRow("", silver),
		plainRow(separator, dim),
		plainRow("", silver),
		plainRow("SELECT A DIFFICULTY", header),
		plainRow("", silver),
		menuRow{segments: []segment{
			{"[ 1 ]  ", keyBracket},
			{"Beginner       ", green},
			{"·  9 × 9    ·  10 mines", silver},
		}},
		menuRow{segments: []segment{
			{"[ 2 ]  ", keyBracket},
			{"Intermediate   ", amber},
			{"·  16 × 16  ·  40 mines", silver},
		}},
		menuRow{segments: []segment{
			{"[ 3 ]  ", keyBracket},
			{"Expert         ", red},
			{"·  30 × 16  ·  99 mines", silver},
		}},
		plainRow("", silver),
		plainRow(separator, dim),
		plainRow("", silver),
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

	rows := menuRows()

	needW := 0
	for _, r := range rows {
		if w := r.width(); w > needW {
			needW = w
		}
	}
	needH := len(rows)

	if sw < needW || sh < needH {
		style := base.Foreground(accentRed).Bold(true)
		drawString(s, 0, 0, "terminal too small for menu — resize", style)
		s.Show()
		return
	}

	startY := (sh - needH) / 2
	for i, r := range rows {
		x := (sw - r.width()) / 2
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
			}
		case *tcell.EventResize:
			screen.Sync()
			drawMenu(screen)
		}
	}
}
