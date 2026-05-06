package main

import "github.com/gdamore/tcell/v2"

var titleArt = []string{
	"███╗   ███╗██╗███╗   ██╗███████╗███████╗██╗    ██╗███████╗███████╗██████╗ ███████╗██████╗ ",
	"████╗ ████║██║████╗  ██║██╔════╝██╔════╝██║    ██║██╔════╝██╔════╝██╔══██╗██╔════╝██╔══██╗",
	"██╔████╔██║██║██╔██╗ ██║█████╗  ███████╗██║ █╗ ██║█████╗  █████╗  ██████╔╝█████╗  ██████╔╝",
	"██║╚██╔╝██║██║██║╚██╗██║██╔══╝  ╚════██║██║███╗██║██╔══╝  ██╔══╝  ██╔═══╝ ██╔══╝  ██╔══██╗",
	"██║ ╚═╝ ██║██║██║ ╚████║███████╗███████║╚███╔███╔╝███████╗███████╗██║     ███████╗██║  ██║",
	"╚═╝     ╚═╝╚═╝╚═╝  ╚═╝╚══════╝╚══════╝ ╚══╝╚══╝ ╚══════╝╚══════╝╚═╝     ╚══════╝╚═╝  ╚═╝",
}

type menuLine struct {
	text  string
	style tcell.Style
}

func menuLines() []menuLine {
	titleStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorRed).
		Bold(true)
	taglineStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorGray).
		Italic(true)
	headerStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorYellow).
		Bold(true)
	optionStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorSilver)
	hintStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorGray)

	lines := make([]menuLine, 0, len(titleArt)+10)
	for _, t := range titleArt {
		lines = append(lines, menuLine{t, titleStyle})
	}
	lines = append(lines,
		menuLine{"", optionStyle},
		menuLine{"— built in Go —", taglineStyle},
		menuLine{"", optionStyle},
		menuLine{"SELECT A DIFFICULTY", headerStyle},
		menuLine{"", optionStyle},
		menuLine{"[ 1 ]  Beginner       9 × 9    10 mines", optionStyle},
		menuLine{"[ 2 ]  Intermediate   16 × 16  40 mines", optionStyle},
		menuLine{"[ 3 ]  Expert         30 × 16  99 mines", optionStyle},
		menuLine{"", optionStyle},
		menuLine{"[ Q ]  Quit", hintStyle},
	)
	return lines
}

func drawMenu(s tcell.Screen) {
	s.Clear()
	sw, sh := s.Size()

	lines := menuLines()

	needW := 0
	for _, ln := range lines {
		w := len([]rune(ln.text))
		if w > needW {
			needW = w
		}
	}
	needH := len(lines)

	if sw < needW || sh < needH {
		style := tcell.StyleDefault.
			Background(tcell.ColorBlack).
			Foreground(tcell.ColorRed).
			Bold(true)
		drawString(s, 0, 0, "terminal too small for menu — resize", style)
		s.Show()
		return
	}

	startY := (sh - needH) / 2
	for i, ln := range lines {
		runes := len([]rune(ln.text))
		x := (sw - runes) / 2
		drawString(s, x, startY+i, ln.text, ln.style)
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
