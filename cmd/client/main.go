//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"
)

///link: https://stdiopt.github.io/gowasm-experiments/bouncy/
///link: https://stdiopt.github.io/gowasm-experiments/hexy/
///link: https://github.com/stdiopt/gowasm-experiments
///link: https://github.com/justinclift/tinygo_canvas2/blob/master/wasm.go

//go:wasmexport access
func access(token string) {
	append(fmt.Sprintf("Access token is '%s'", token))
}

func append(text string) {
	document := js.Global().Get("document")
	body := document.Call("querySelector", "body")
	div := document.Call("createElement", "div")
	div.Set("innerHTML", text)
	body.Call("appendChild", div)
}

func main() {
	append("Started")
	<-make(chan bool)
}
