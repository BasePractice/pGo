package game

import (
	"log"
	"sokoban/game/maps"
	"testing"
)

func TestEngine_Moving(t *testing.T) {
	values := [][]maps.Tail{
		{0, 0, 1, 1, 1, 1, 1, 0},
		{1, 1, 1, 0, 0, 0, 1, 0},
		{1, 2, 4, 3, 0, 0, 1, 0},
		{1, 1, 1, 0, 3, 2, 1, 0},
		{1, 2, 1, 1, 3, 0, 1, 0},
		{1, 0, 1, 0, 2, 0, 1, 1},
		{1, 3, 0, 5, 3, 3, 2, 1},
		{1, 0, 0, 0, 2, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1},
	}
	game := &Game{0, "demo", values, false, false, EmptyUpdater{}, 2, 2}
	game.Left()
	if values[2][1] != maps.TailPlayerAndSpot {
		log.Fatal("values[2][3] is not a tail player")
	}
}
