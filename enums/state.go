
package enums

type GameState int

const (
	StateLost GameState = iota,
	StateWon,
	StatePlaying
)

var stateName = map[GameState]string {
	StateLost: "lost",
	StateWon: "won",
	StatePlaying: "playing",
}

func (gs GameState) String() string {
	return stateName[gs]
}

