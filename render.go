package main

import "github.com/gdamore/tcell/v2"

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

func (g *Game) draw(s tcell.Screen) {
	s.Clear()

	base := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)

	for i := range g.grid {
		for j := range g.grid[i] {
			tile := g.grid[i][j]
			x, y := j*cellWidth, i

			glyph := ' '
			style := base

			hiddenBg := tcell.ColorDarkSlateGray
			if (i+j)%2 == 1 {
				hiddenBg = tcell.ColorDimGray
			}

			switch {
			case !tile.isRevealed && tile.isFlagged:
				glyph = 'F'
				style = style.Foreground(tcell.ColorRed).Background(hiddenBg)
			case !tile.isRevealed:
				glyph = ' '
				style = style.Background(hiddenBg)
			case tile.isMine:
				glyph = '*'
				style = style.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)
			case tile.adjMines == 0:
				glyph = ' '
			default:
				glyph = rune('0' + tile.adjMines)
				if c, ok := numberColors[tile.adjMines]; ok {
					style = style.Foreground(c)
				}
			}

			if i == g.cursorI && j == g.cursorJ {
				style = style.Reverse(true)
			}

			s.SetContent(x,   y, ' ',   nil, style)
			s.SetContent(x+1, y, glyph, nil, style)
			s.SetContent(x+2, y, ' ',   nil, style)
		}
	}

	s.Show()
}
