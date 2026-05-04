package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/minesweeper-go/enums"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Println("failed to create screen:", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Println("failed to init screen:", err)
		os.Exit(1)
	}
	defer screen.Fini()

	game := NewGame(enums.Beginner)
	game.draw(screen)

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

			if game.state != enums.StatePlaying {
				continue
			}

			if action, ok := dispatch(ev, gameplayKeyBindings, gameplayRuneBindings); ok {
				action(game)
				game.draw(screen)
			}
		case *tcell.EventResize:
			screen.Sync()
			game.draw(screen)
		}
	}
}
