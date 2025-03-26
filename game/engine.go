package game

import "sokoban/game/maps"

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
	d       Updater
	x       int
	y       int
}

type Updater interface {
	Update(v [][]maps.Tail)
}

type EmptyUpdater struct{}

func (e EmptyUpdater) Update(v [][]maps.Tail) {

}

type Mover interface {
	Left()
	Right()
	Up()
	Down()
}

func Ctor(mapName string, drawer Updater) (*Game, error) {
	m := maps.Map{}
	err := FileUnmarshal(mapName, &m)
	if err != nil {
		return nil, err
	}
	x := 0
	y := 0
	for i := 0; i < m.Height; i++ {
		for j := 0; j < m.Width; j++ {
			v := *m.At(i, j)
			if v == maps.TailPlayer {
				x = j
				y = i
				break
			}
		}
	}
	return &Game{0, mapName, m.Values, false, false, drawer, x, y}, nil
}

func (g *Game) moviePlayerTo(x, y int) {
	if g.v[g.x][g.y] == maps.TailPlayer {
		if g.v[x][y] == maps.TailSpot {
			g.v[x][y] = maps.TailPlayerAndSpot
		} else {
			g.v[x][y] = maps.TailPlayer
		}
		g.v[g.x][g.y] = maps.TailNone
	} else if g.v[g.x][g.y] == maps.TailPlayerAndSpot {
		if g.v[x][y] == maps.TailSpot {
			g.v[x][y] = maps.TailPlayerAndSpot
		} else {
			g.v[x][y] = maps.TailPlayer
		}
		g.v[g.x][g.y] = maps.TailSpot
	}
}

func (g *Game) moving(to mover) {
	nX := g.x + to.dX
	nY := g.x + to.dY
	if g.mayMoving(nX, nY) {
		g.movePlayer(nX, nY)
	} else if g.isBox(nX, nY) && g.mayMoving(nX+to.dX, nY+to.dY) {
		g.moveBox(nX+to.dX, nY+to.dY)
		g.movePlayer(nX, nY)
	}
}

func (g *Game) movePlayer(x, y int) {
	if g.x != x || g.y != y {
		g.d.Update(g.v)
		g.moviePlayerTo(x, y)
		g.Steps++
	}
	g.x = x
	g.y = y
}

func (g *Game) moveBox(x, y int) {
	//TODO: Реализовать перемещение
	if g.x != x || g.y != y {
		g.d.Update(g.v)
		g.moviePlayerTo(x, y)
		g.Steps++
	}
	g.x = x
	g.y = y
}

func (g *Game) mayMoving(x, y int) bool {
	if g.v[y][x] == maps.TailNone || g.v[x][y] == maps.TailSpot {
		return true
	}
	return false
}

func (g *Game) isBox(x, y int) bool {
	if g.v[y][x] == maps.TailBox || g.v[x][y] == maps.TailBoxAndSpot {
		return true
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
