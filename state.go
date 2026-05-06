package main

type GameState int

const (
	StatePlaying GameState = iota
	StateWon
	StateLost
)

var stateName = map[GameState]string{
	StatePlaying: "playing",
	StateWon:     "won",
	StateLost:    "lost",
}

func (gs GameState) String() string {
	return stateName[gs]
}
