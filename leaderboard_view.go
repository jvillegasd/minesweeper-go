package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func leaderboardRows(sw int) []menuRow {
	titleStyle := base.Foreground(accentYellow).Bold(true)
	tagline := base.Foreground(textSilver).Italic(true)
	dim := base.Foreground(textDim)
	silver := base.Foreground(textSilver)
	silverBold := base.Foreground(textSilver).Bold(true)
	rankColor := base.Foreground(accentGreen).Bold(true)

	difficultyStyle := map[Difficulty]tcell.Style{
		Beginner:     base.Foreground(accentGreen).Bold(true),
		Intermediate: base.Foreground(accentYellow).Bold(true),
		Expert:       base.Foreground(accentRed).Bold(true),
	}

	sepWidth := 60
	if sepWidth > sw-4 {
		sepWidth = sw - 4
	}
	if sepWidth < 10 {
		sepWidth = 10
	}

	rows := []menuRow{
		menuRow{segments: []segment{
			{"✦  ", titleStyle},
			{"LEADERBOARD", titleStyle},
			{"  ✦", titleStyle},
		}},
		plainRow("", silver),
		plainRow("✦  best times per difficulty  ✦", tagline),
		plainRow("", silver),
		plainRow(strings.Repeat("━", sepWidth), dim),
		plainRow("", silver),
	}

	for _, d := range []Difficulty{Beginner, Intermediate, Expert} {
		rows = append(rows, plainRow(strings.ToUpper(d.String()), difficultyStyle[d]))

		entries := leaderboard.entries(d)
		if len(entries) == 0 {
			rows = append(rows, menuRow{segments: []segment{
				{"    no records yet", dim},
			}})
		} else {
			for i, e := range entries {
				medal := "  "
				medalStyle := dim
				switch i {
				case 0:
					medal = "🥇"
					medalStyle = rankColor
				case 1:
					medal = "🥈"
				case 2:
					medal = "🥉"
				}
				rows = append(rows, menuRow{segments: []segment{
					{fmt.Sprintf("    %d. ", i+1), dim},
					{medal + " ", medalStyle},
					{fmt.Sprintf("%4ds", e.Seconds), silverBold},
					{"    ", silver},
					{e.Date.Format("2006-01-02"), dim},
				}})
			}
		}
		rows = append(rows, plainRow("", silver))
	}

	rows = append(rows,
		plainRow(strings.Repeat("━", sepWidth), dim),
		plainRow("", silver),
		plainRow("press any key to return", dim),
	)
	return rows
}

func drawLeaderboard(s tcell.Screen) {
	s.Clear()
	sw, sh := s.Size()

	rows := leaderboardRows(sw)
	startY := (sh - len(rows)) / 2
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

func runLeaderboard(screen tcell.Screen) {
	drawLeaderboard(screen)
	for {
		ev := screen.PollEvent()
		switch ev.(type) {
		case *tcell.EventKey:
			return
		case *tcell.EventResize:
			screen.Sync()
			drawLeaderboard(screen)
		}
	}
}
