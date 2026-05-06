package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
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

	game := NewGame(Beginner)
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
