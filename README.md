# minesweeper-go

A terminal Minesweeper written in Go. Keyboard-driven, transparent-friendly TUI, persistent leaderboard.

```
              ✦ MINESWEEPER ✦
            ✦  built in Go  ✦

       ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

                SELECT A DIFFICULTY

       [ 1 ]  Beginner       ·  9 × 9    ·  10 mines
       [ 2 ]  Intermediate   ·  16 × 16  ·  40 mines
       [ 3 ]  Expert         ·  30 × 16  ·  99 mines

       ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

       [ L ]  Leaderboard
       [ Q ]  Quit
```

## Features

- Three difficulty levels (Beginner / Intermediate / Expert)
- First-click safety — first reveal is guaranteed safe and opens a region
- Chord-reveal — auto-reveal neighbors of fully-flagged numbered tiles
- Live timer + flag counter
- Persistent top-5 leaderboard per difficulty
- Adaptive layout that resizes with the terminal
- Transparent-terminal-friendly rendering
- Mid-game difficulty switching, restart, and back-to-menu

## Requirements

- Go **1.26** or newer
- A terminal that supports Unicode + 24-bit color (most modern terminals: iTerm2, Alacritty, Kitty, WezTerm, GNOME Terminal, Windows Terminal)

## Install

### Option 1 — go install (recommended once published)

```sh
go install github.com/minesweeper-go@latest
```

Binary lands in `$GOBIN` (usually `~/go/bin`). Make sure that's on your `PATH`:

```sh
export PATH=$PATH:$(go env GOBIN):$(go env GOPATH)/bin
```

### Option 2 — clone and build

```sh
git clone https://github.com/minesweeper-go.git
cd minesweeper-go
go build -o minesweeper-go .
./minesweeper-go
```

### Option 3 — Makefile

```sh
make            # build into ./minesweeper-go
make install    # install into /usr/local/bin (may need sudo)
make uninstall  # remove from /usr/local/bin
make test       # run the test suite
make clean      # remove the binary
```

### Option 4 — run without installing

```sh
go run .
```

## Usage

```sh
minesweeper-go                # opens the menu
minesweeper-go beginner       # skip menu, jump into Beginner
minesweeper-go intermediate   # or short forms: b / i / e
minesweeper-go expert
```

## Controls

### In the menu

| Key | Action |
|---|---|
| `1` / `2` / `3` | start Beginner / Intermediate / Expert |
| `l` | open leaderboard |
| `q` / `Esc` | quit |

### In game

| Key | Action |
|---|---|
| `↑ ↓ ← →` | move cursor |
| `space` / `Enter` | reveal tile under cursor |
| `f` | toggle flag |
| `c` | chord — reveal neighbors of a fully-flagged numbered tile |
| `1` / `2` / `3` | restart at that difficulty |
| `r` | restart current difficulty |
| `m` | back to menu |
| `l` | open leaderboard |
| `q` / `Esc` | quit |

### After winning or losing

| Key | Action |
|---|---|
| `r` | restart |
| `m` | back to menu |
| `l` | open leaderboard |
| `q` | quit |

## Leaderboard

Top 5 best times per difficulty, persisted to:

```
~/.minesweeper-go-scores.json
```

Records update automatically on every win. Delete the file to reset.

## Development

```sh
go test ./...   # run the test suite (19 tests)
go vet ./...    # static analysis
gofmt -w .      # format
```

### Project structure

| File | Purpose |
|---|---|
| `main.go` | entry point, event loop, CLI arg parsing |
| `game.go` | `Game` state, mine placement, win/loss handling |
| `tile.go` | `Tile`, flood-fill reveal (BFS), chord-reveal |
| `bindings.go` | keyboard action map |
| `render.go` | board, status, HUD, palette, border |
| `menu.go` | start menu |
| `leaderboard.go` | persistent top-5 score storage |
| `leaderboard_view.go` | full-screen leaderboard view |
| `state.go` | `GameState` enum |
| `level.go` | `Difficulty` and board dimensions |
| `game_test.go` | unit tests |

## License

[MIT](LICENSE)
