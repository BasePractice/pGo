package maps

import (
	"encoding/json"
	"fmt"
)

type Tail int

const (
	TailNone Tail = iota
	TailWall
	TailSpot
	TailBox
	TailPlayer
	TailBoxAndSpot
	TailPlayerAndSpot
)

type Map struct {
	Width  int      `json:"width"`
	Height int      `json:"height"`
	Values [][]Tail `json:"values"`
}

func Ctor(width, height int) *Map {
	skeleton := make([][]Tail, height)
	for i := range skeleton {
		skeleton[i] = make([]Tail, width)
	}
	return &Map{width, height, skeleton}
}

func (m *Map) Marshal() ([]byte, error) {
	return json.Marshal(*m)
}

func (m *Map) Unmarshal(raw []byte) error {
	return json.Unmarshal(raw, m)
}

func (m *Map) String() string {
	return fmt.Sprintf("Map{width: %d, height: %d}", m.Width, m.Height)
}

func (m *Map) At(x, y int) *Tail {
	return &m.Values[y][x]
}
