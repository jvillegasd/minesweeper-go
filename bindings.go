package main

import "github.com/gdamore/tcell/v2"

type Action func(*Game) bool

var metaKeyBindings = map[tcell.Key]Action{
	tcell.KeyEscape: func(g *Game) bool { return true },
}

var metaRuneBindings = map[rune]Action{
	'q': func(g *Game) bool { return true },
	'r': func(g *Game) bool { g.Restart(); return false },
}

var gameplayKeyBindings = map[tcell.Key]Action{
	tcell.KeyUp:    func(g *Game) bool { g.MoveCursor(-1, 0); return false },
	tcell.KeyDown:  func(g *Game) bool { g.MoveCursor(+1, 0); return false },
	tcell.KeyLeft:  func(g *Game) bool { g.MoveCursor(0, -1); return false },
	tcell.KeyRight: func(g *Game) bool { g.MoveCursor(0, +1); return false },
	tcell.KeyEnter: func(g *Game) bool { g.RevealAtCursor(); return false },
}

var gameplayRuneBindings = map[rune]Action{
	'f': func(g *Game) bool { g.FlagAtCursor(); return false },
	' ': func(g *Game) bool { g.RevealAtCursor(); return false },
}

func dispatch(ev *tcell.EventKey, keyMap map[tcell.Key]Action, runeMap map[rune]Action) (Action, bool) {
	if a, ok := keyMap[ev.Key()]; ok {
		return a, true
	}
	if ev.Key() == tcell.KeyRune {
		if a, ok := runeMap[ev.Rune()]; ok {
			return a, true
		}
	}
	return nil, false
}
