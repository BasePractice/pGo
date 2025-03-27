//go:build js && wasm

package main

import "syscall/js"

///link: https://stdiopt.github.io/gowasm-experiments/bouncy/
///link: https://stdiopt.github.io/gowasm-experiments/hexy/
///link: https://github.com/stdiopt/gowasm-experiments
///link: https://github.com/justinclift/tinygo_canvas2/blob/master/wasm.go

func main() {
	document := js.Global().Get("document")
	body := document.Call("querySelector", "body")
	div := document.Call("createElement", "div")
	div.Set("innerHTML", "Hello, WASM")
	body.Call("appendChild", div)
	<-make(chan bool)
}
