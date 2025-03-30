package game

import (
	"testing"

	"sokoban/game/maps"
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
		t.Fatal("values[2][1] is not a tail player")
	}
	game.Right()
	if values[2][2] != maps.TailPlayer {
		t.Fatal("values[2][2] is not a tail player")
	} else if values[2][1] != maps.TailSpot {
		t.Fatal("values[2][1] is not a tail spot")
	}
	game.Right()
	if values[2][3] != maps.TailPlayer {
		t.Fatal("values[2][3] is not a tail player")
	} else if values[2][4] != maps.TailBox {
		t.Fatal("values[2][4] is not a tail box")
	}
	game.Up()
	game.Right()
	game.Right()
	game.Down()
	game.Down()
	game.Down()
	game.Down()
	game.Left()
	game.Down()
	if values[7][4] != maps.TailBoxAndSpot {
		t.Fatal("values[8][5] is not a tail box and spot")
	}
}
