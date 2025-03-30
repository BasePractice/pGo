//go:build js && wasm

package main

import (
	"encoding/json"
	"image/color"
	"log"
	"strconv"
	"syscall/js"

	"github.com/lucasb-eyer/go-colorful"

	"github.com/stdiopt/gowasm-experiments/arty/painter"
)

///link: https://stdiopt.github.io/gowasm-experiments/bouncy/
///link: https://stdiopt.github.io/gowasm-experiments/hexy/
///link: https://github.com/stdiopt/gowasm-experiments
///link: https://github.com/justinclift/tinygo_canvas2/blob/master/wasm.go

var token string = ""

type pos struct {
	x, y float64
}

type WasmContext struct {
	done    chan struct{}
	painter *painter.BufPainter
	addr    string

	doc      js.Value
	canvasEl js.Value
	ctx      js.Value
	ws       js.Value
	im       js.Value

	colorHex  string
	lineWidth float64
	textOff   pos
	lastPos   pos

	byteArray js.Value
	width     float64
	height    float64
}

func CtorContext(addr string) (*WasmContext, error) {
	done := make(chan struct{})

	painter, err := painter.New()
	if err != nil {
		return nil, err
	}
	return &WasmContext{
		done:      done,
		painter:   painter,
		addr:      addr,
		lineWidth: 10,
	}, nil
}

func (c *WasmContext) Start() {
	c.initCanvas()
	c.initFrameUpdate()
	c.initConnection()
	<-c.done
}

func (c *WasmContext) Close() {
	close(c.done)
}

func (c *WasmContext) initCanvas() {
	c.doc = js.Global().Get("document")
	c.canvasEl = c.doc.Call("getElementById", "game_id")
	c.width = c.canvasEl.Get("width").Float()
	c.height = c.canvasEl.Get("height").Float()
	c.ctx = c.canvasEl.Call("getContext", "2d")
	c.im = c.ctx.Call("createImageData", 1, 1)
	c.byteArray = js.Global().Get("Uint8Array").New(1 * 4)
	c.painter.OnInit = func(m painter.InitOP) {
		c.im = c.ctx.Call("createImageData", m.Width, m.Height)
		c.byteArray = js.Global().Get("Uint8Array").New(m.Width * m.Height * 4)
		c.SetStatus("connected")
		c.initEvents()
	}
}

func (c *WasmContext) initFrameUpdate() {
	// Hold the callbacks without blocking
	go func() {
		var renderFrame js.Func
		renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			c.draw()
			js.Global().Call("requestAnimationFrame", renderFrame)
			return nil
		})
		defer renderFrame.Release()
		js.Global().Call("requestAnimationFrame", renderFrame)
		<-c.done
	}()
}

func (c *WasmContext) initConnection() {
	go func() {
		c.SetStatus("connecting...")
		c.ws = js.Global().Get("WebSocket").New(c.addr)
		onopen := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			c.SetStatus("receiving... (it takes some time)")
			return nil
		})
		defer onopen.Release()
		onmessage := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			c.painter.HandleRaw([]byte(args[0].Get("data").String()))
			return nil
		})
		defer onmessage.Release()
		c.ws.Set("onopen", onopen)
		c.ws.Set("onmessage", onmessage)

		<-c.done
	}()
}

func (c *WasmContext) draw() {
	// golang buffer
	// Needs to be a Uint8Array while image data have Uint8ClampedArray
	js.CopyBytesToJS(c.byteArray, c.painter.ImageData())
	c.im.Get("data").Call("set", c.byteArray)
	c.ctx.Call("putImageData", c.im, 0, 0)
}

func (c *WasmContext) initEvents() {
	go func() {
		// DOM events
		colorEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			e := args[0]
			c.colorHex = e.Get("target").Get("value").String()
			return nil
		})
		szEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			e := args[0]
			v, _ := strconv.ParseFloat(e.Get("target").Get("value").String(), 64)
			c.lineWidth = v
			return nil
		})
		defer szEvt.Release()

		c.doc.Call("getElementById", "color").Call("addEventListener", "change", colorEvt)
		c.doc.Call("getElementById", "size").Call("addEventListener", "change", szEvt)

		keyPressEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			e := args[0]
			e.Call("preventDefault")
			key := e.Get("key").String()
			c.SetStatus(key)
			col, _ := colorful.Hex(c.colorHex) // Ignore error
			op := painter.TextOP{
				color.RGBA{uint8(col.R * 255), uint8(col.G * 255), uint8(col.B * 255), 255},
				c.lineWidth + 6,
				c.lastPos.x + c.textOff.x, c.lastPos.y + c.textOff.y,
				key,
			}
			c.textOff.x += (c.lineWidth + 10) * 0.6

			c.painter.HandleOP(op)
			buf, err := json.Marshal(painter.Message{op})
			if err != nil {
				return nil
			}
			c.ws.Call("send", string(buf))
			return nil

		})
		defer keyPressEvt.Release()
		c.doc.Call("addEventListener", "keypress", keyPressEvt)

		<-c.done
	}()
}

func (c *WasmContext) SetStatus(txt string) {
	c.doc.Call("getElementById", "status").Set("innerHTML", txt)
}

func setToken(this js.Value, args []js.Value) interface{} {
	if len(args) >= 1 {
		token = args[0].String()
	}
	return nil
}

func main() {
	js.Global().Set("token", js.FuncOf(setToken))
	c, err := CtorContext("ws:/localhost:9090/ws")
	if err != nil {
		log.Fatal("could not start", err)
	}
	defer c.Close()
	c.Start()
}
