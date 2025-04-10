package game

import (
	"strconv"
	"strings"

	"sokoban/game/maps"
)

var (
	MoveLeft  = mover{dX: -1, dY: 0}
	MoveRight = mover{dX: 1, dY: 0}
	MoveUp    = mover{dX: 0, dY: -1}
	MoveDown  = mover{dX: 0, dY: 1}
)

type mover struct {
	dX int
	dY int
}

type Game struct {
	Steps   int    `json:"steps"`
	MapName string `json:"map_name"`
	v       [][]maps.Tail
	changed bool
	done    bool
	u       Updater
	x       int
	y       int
}

type Updater interface {
	UpdateMap(v [][]maps.Tail)
}

type EmptyUpdater struct{}

func (e EmptyUpdater) UpdateMap(v [][]maps.Tail) {

}

type Mover interface {
	Left()
	Right()
	Up()
	Down()
}

func CtorNone() *Game {
	return &Game{0, "none", make([][]maps.Tail, 0), false, false, nil, 0, 0}
}

func CtorFile(filename string, updater Updater) (*Game, error) {
	m := maps.Map{}
	err := FileUnmarshal(filename, &m)
	if err != nil {
		return nil, err
	}
	xPlayer := 0
	yPlayer := 0
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			v := *m.At(x, y)
			if v == maps.TailPlayer {
				xPlayer = x
				yPlayer = y
				break
			}
		}
	}
	return &Game{0, filename, m.Values, false, false, updater, xPlayer, yPlayer}, nil
}

func (g *Game) UpdateLine(name, line string, width, height int) {
	parts := strings.Split(line, ",")
	g.x = 0
	g.y = 0
	g.v = make([][]maps.Tail, height)
	for y := 0; y < height; y++ {
		g.v[y] = make([]maps.Tail, width)
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			v, _ := strconv.Atoi(parts[(y*width)+x])
			g.v[y][x] = maps.Tail(v)
			if g.v[y][x] == maps.TailPlayer {
				g.x = x
				g.y = y
			}
		}
	}
}

func (g *Game) moviePlayerTo(x, y int) {
	if g.v[g.y][g.x] == maps.TailPlayer {
		if g.v[y][x] == maps.TailSpot {
			g.v[y][x] = maps.TailPlayerAndSpot
		} else {
			g.v[y][x] = maps.TailPlayer
		}
		g.v[g.y][g.x] = maps.TailNone
	} else if g.v[g.y][g.x] == maps.TailPlayerAndSpot {
		if g.v[y][x] == maps.TailSpot {
			g.v[y][x] = maps.TailPlayerAndSpot
		} else {
			g.v[y][x] = maps.TailPlayer
		}
		g.v[g.y][g.x] = maps.TailSpot
	}
}

func (g *Game) moving(to mover) {
	nX := g.x + to.dX
	nY := g.y + to.dY
	if g.mayMoving(nX, nY) {
		g.movePlayer(nX, nY)
	} else if g.isBox(nX, nY) && g.mayMoving(nX+to.dX, nY+to.dY) {
		g.moveBox(nX+to.dX, nY+to.dY)
		if g.v[nY][nX] == maps.TailBox {
			g.v[nY][nX] = maps.TailNone
		} else if g.v[nY][nX] == maps.TailBoxAndSpot {
			g.v[nY][nX] = maps.TailSpot
		}

		g.movePlayer(nX, nY)
	}
}

func (g *Game) movePlayer(x, y int) {
	if g.x != x || g.y != y {
		g.moviePlayerTo(x, y)
		g.u.UpdateMap(g.v)
		g.Steps++
	}
	g.x = x
	g.y = y
}

func (g *Game) moveBox(x, y int) {
	if g.v[y][x] == maps.TailSpot {
		g.v[y][x] = maps.TailBoxAndSpot
	} else {
		g.v[y][x] = maps.TailBox
	}
}

func (g *Game) mayMoving(x, y int) bool {
	if len(g.v) > x && x >= 0 && len(g.v[x]) > y && y >= 0 {
		if g.v[y][x] == maps.TailNone || g.v[y][x] == maps.TailSpot {
			return true
		}
	}
	return false
}

func (g *Game) isBox(x, y int) bool {
	if len(g.v) > x && x >= 0 && len(g.v[x]) > y && y >= 0 {
		if g.v[y][x] == maps.TailBox || g.v[y][x] == maps.TailBoxAndSpot {
			return true
		}
	}
	return false
}

func (g *Game) Left() {
	g.moving(MoveLeft)
}
func (g *Game) Right() {
	g.moving(MoveRight)
}
func (g *Game) Up() {
	g.moving(MoveUp)
}
func (g *Game) Down() {
	g.moving(MoveDown)
}

func (g *Game) IsWon() bool {
	for y := 0; y < len(g.v); y++ {
		for x := 0; x < len(g.v[y]); x++ {
			if g.v[y][x] == maps.TailSpot {
				return false
			}
		}
	}
	return true
}

func (g *Game) Updater(u Updater) {
	g.u = u
	if u != nil {
		u.UpdateMap(g.v)
	}
}
