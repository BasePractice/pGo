//go:build js && wasm

package main

import (
	"encoding/json"
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
			type ms struct {
				Width  int    `json:"width"`
				Height int    `json:"height"`
				Data   string `json:"data"`
			}
			if len(args) > 0 && args[0].Type() == js.TypeObject {
				var d ms
				data := args[0].Get("data").String()
				json.Unmarshal([]byte(data), &d)
				if d.Width != 0 || d.Height != 0 {
					c.game.UpdateLine("demo", d.Data, d.Width, d.Height)
					c.loop.Refresh()
				}
				c.SetStatus(fmt.Sprintf("%+v", d))
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
