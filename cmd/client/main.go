//go:build js && wasm

package main

import (
	"fmt"
	"log"
	"sokoban/game"
	"sokoban/game/ui"
	"syscall/js"
)

type WasmContext struct {
	done chan struct{}
	addr string

	ws  js.Value
	doc js.Value

	game *game.Game
	loop *ui.Loop
}

func CtorContext(addr string) (*WasmContext, error) {
	done := make(chan struct{})

	return &WasmContext{
		done: done,
		addr: addr,
		game: game.CtorNone(),
	}, nil
}

func (c *WasmContext) Start() {
	c.initConnection()
	c.doc = js.Global().Get("document")
	c.loop = ui.Ctor(c.game)
	c.loop.Start()
	<-c.done
}

func (c *WasmContext) Close() {
	close(c.done)
}

func (c *WasmContext) initConnection() {
	go func() {
		c.ws = js.Global().Get("WebSocket").New(c.addr)
		onopen := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			c.SetStatus("receiving... (it takes some time)")
			return nil
		})
		defer onopen.Release()
		onmessage := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) > 0 && args[0].Type() == js.TypeObject {
				data := args[0].Get("data").String()
				width := args[0].Get("width").Type() == js.TypeNumber
				height := args[0].Get("height").Type() == js.TypeNumber
				if width && height {
					width := args[0].Get("width").Int()
					height := args[0].Get("height").Int()
					c.game.UpdateLine("demo", data, width, height)
				}
				c.SetStatus(fmt.Sprintf("%+v", args[0]))
			} else {
				c.SetStatus(fmt.Sprintf("Illegal arguments: %+v", args))
			}
			return nil
		})
		defer onmessage.Release()
		c.ws.Set("onopen", onopen)
		c.ws.Set("onmessage", onmessage)

		<-c.done
	}()
}

func (c *WasmContext) SetStatus(txt string) {
	c.doc.Call("getElementById", "status").Set("innerHTML", txt)
}

func main() {
	c, err := CtorContext("ws:/localhost:9090/ws")
	if err != nil {
		log.Fatal("could not start", err)
	}
	defer c.Close()
	c.Start()
}
