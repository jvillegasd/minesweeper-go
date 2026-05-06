package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

func parseDifficulty(arg string) (Difficulty, bool) {
	switch strings.ToLower(arg) {
	case "beginner", "b":
		return Beginner, true
	case "intermediate", "i":
		return Intermediate, true
	case "expert", "e":
		return Expert, true
	}
	return 0, false
}

func main() {
	var (
		startDifficulty Difficulty
		skipMenu        bool
	)

	if len(os.Args) > 1 {
		d, ok := parseDifficulty(os.Args[1])
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown difficulty: %s\n", os.Args[1])
			fmt.Fprintln(os.Stderr, "valid: beginner, intermediate, expert")
			os.Exit(1)
		}
		startDifficulty = d
		skipMenu = true
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Println("failed to create screen:", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Println("failed to init screen:", err)
		os.Exit(1)
	}
	screen.SetStyle(tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite))
	defer screen.Fini()

	if !skipMenu {
		d, ok := runMenu(screen)
		if !ok {
			return
		}
		startDifficulty = d
	}

	game := NewGame(startDifficulty)
	game.draw(screen)

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			screen.PostEvent(tcell.NewEventInterrupt(nil))
		}
	}()

	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if action, ok := dispatch(ev, metaKeyBindings, metaRuneBindings); ok {
				if action(game) {
					return
				}
				game.draw(screen)
				continue
			}

			if game.state != StatePlaying {
				continue
			}

			if action, ok := dispatch(ev, gameplayKeyBindings, gameplayRuneBindings); ok {
				action(game)
				game.draw(screen)
			}
		case *tcell.EventResize:
			screen.Sync()
			game.draw(screen)
		case *tcell.EventInterrupt:
			game.draw(screen)
		}
	}
}
