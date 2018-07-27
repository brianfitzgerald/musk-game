package model

type Room struct {
	Players   int
	MuskDrawn bool
	Code      string
	Disaster  string
}

const (
	TableName = "musk-game-rooms"
)
