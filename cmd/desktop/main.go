//go:build windows || darwin

package main

import (
	"sokoban/game"
	"sokoban/game/ui"
)

func main() {
	g, _ := game.Ctor("resources/demo.json", nil)
	var desktop = ui.Ctor(g)
	g.Updater(desktop)
	desktop.Start()
}
