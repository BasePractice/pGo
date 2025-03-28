package main

import (
	"image/color"
	"log"

	"sokoban/game"
	"sokoban/game/maps"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 400
	screenHeight = 400
)

var (
	gray   = color.RGBA{R: 128, G: 128, B: 128, A: 255}
	orange = color.RGBA{R: 235, G: 131, B: 52, A: 255}
	red    = color.RGBA{R: 235, G: 52, B: 52, A: 255}
	yellow = color.RGBA{R: 235, G: 217, B: 52, A: 255}
	blue   = color.RGBA{R: 52, G: 52, B: 235, A: 255}
	pink   = color.RGBA{R: 235, G: 52, B: 226, A: 255}
	none   = color.RGBA{R: 128, G: 52, B: 235, A: 255}
)

type Desktop struct {
	game          *game.Game
	bord          [][]maps.Tail
	width, height int
}

func getColor(v maps.Tail) color.Color {
	switch v {
	case maps.TailNone:
		return gray
	case maps.TailWall:
		return orange
	case maps.TailSpot:
		return red
	case maps.TailBox:
		return yellow
	case maps.TailPlayer:
		return blue
	case maps.TailBoxAndSpot:
		return pink
	case maps.TailPlayerAndSpot:
		return none
	default:
		return gray
	}
}

func (d *Desktop) UpdateMap(v [][]maps.Tail) {
	d.bord = make([][]maps.Tail, len(v))
	for i := range d.bord {
		d.bord[i] = make([]maps.Tail, len(v))
		copy(d.bord[i], v[i])
	}
	d.width, d.height = len(v[0]), len(v)
}

func (d *Desktop) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		d.game.Left()
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		d.game.Right()
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		d.game.Up()
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		d.game.Down()
	}
	return nil
}

func (d *Desktop) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 30, G: 30, B: 30, A: 255})

	if d.width == 0 || d.height == 0 {
		return
	}
	var size = screen.Bounds().Dx() / d.width
	for y, yv := range d.bord {
		for x, xv := range yv {
			color := getColor(xv)
			vector.DrawFilledRect(screen, float32(x*size), float32(y*size),
				float32(size), float32(size), color, true)
		}
	}
}

func (d *Desktop) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g, _ := game.Ctor("resources/demo.json", nil)
	var desktop = Desktop{game: g, bord: [][]maps.Tail{}}
	g.Updater(&desktop)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sokoban")
	if err := ebiten.RunGame(&desktop); err != nil {
		log.Fatal(err)
	}
}
