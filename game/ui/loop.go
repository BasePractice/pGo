package ui

import (
	"image/color"
	"log"
	_ "log"

	"sokoban/game"
	"sokoban/game/maps"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

type Loop struct {
	game          *game.Game
	bord          [][]maps.Tail
	width, height int
}

func Ctor(g *game.Game) *Loop {
	return &Loop{game: g, bord: [][]maps.Tail{}}
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

func (d *Loop) UpdateMap(v [][]maps.Tail) {
	d.bord = make([][]maps.Tail, len(v))
	for i := range d.bord {
		d.bord[i] = make([]maps.Tail, len(v))
		copy(d.bord[i], v[i])
	}
	d.width, d.height = len(v[0]), len(v)
}

func (d *Loop) Update() error {
	if repeatingKeyPressed(ebiten.KeyA) || repeatingKeyPressed(ebiten.KeyLeft) {
		d.game.Left()
	}
	if repeatingKeyPressed(ebiten.KeyD) || repeatingKeyPressed(ebiten.KeyRight) {
		d.game.Right()
	}
	if repeatingKeyPressed(ebiten.KeyW) || repeatingKeyPressed(ebiten.KeyUp) {
		d.game.Up()
	}
	if repeatingKeyPressed(ebiten.KeyS) || repeatingKeyPressed(ebiten.KeyDown) {
		d.game.Down()
	}
	return nil
}

func (d *Loop) Draw(screen *ebiten.Image) {
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

func (d *Loop) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (d *Loop) Start() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sokoban")
	if err := ebiten.RunGame(d); err != nil {
		log.Fatal(err)
	}
}

func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}
