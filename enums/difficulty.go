package enums

type Difficulty int

const (
	Beginner Difficulty = iota
	Intermediate
	Expert
)

var difficultyName = map[Difficulty]string{
	Beginner:     "beginner",
	Intermediate: "intermediate",
	Expert:       "expert",
}

func (d Difficulty) String() string {
	return difficultyName[d]
}

type Level struct {
	Width, Height, Mines int
}

var Levels = map[Difficulty]Level{
	Beginner:     {Width: 9, Height: 9, Mines: 10},
	Intermediate: {Width: 16, Height: 16, Mines: 40},
	Expert:       {Width: 30, Height: 16, Mines: 99},
}
