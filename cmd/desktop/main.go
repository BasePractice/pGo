//go:build windows || darwin

package main

import (
	"sokoban/game"
	"sokoban/game/ui"
)

func main() {
	g, _ := game.CtorFile("resources/demo.json", nil)
	var desktop = ui.Ctor(g)
	desktop.Start()
}
