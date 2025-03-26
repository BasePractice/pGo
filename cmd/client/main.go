//go:build js && wasm

package main

import "syscall/js"

func main() {
	document := js.Global().Get("document")
	body := document.Call("querySelector", "body")
	div := document.Call("createElement", "div")
	div.Set("innerHTML", "Hello, WASM")
	body.Call("appendChild", div)
	<-make(chan bool)
}
