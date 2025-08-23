//go:build js && wasm

package main

import (
	"log"
	"syscall/js"

	"github.com/omid3699/god_says/internal"
)

var g *internal.God

func initGod(this js.Value, args []js.Value) any {
	amount := internal.DefaultAmount
	if len(args) > 0 && args[0].Type() == js.TypeNumber {
		amount = args[0].Int()
	}
	ng, err := internal.NewGod(amount)
	if err != nil {
		return map[string]any{"ok": false, "error": err.Error()}
	}
	g = ng
	return map[string]any{"ok": true}
}

func speak(this js.Value, args []js.Value) any {
	if g == nil {
		return map[string]any{"ok": false, "error": "god not initialized"}
	}
	return map[string]any{"ok": true, "message": g.Speak()}
}

func main() {
	js.Global().Set("initGod", js.FuncOf(initGod))
	js.Global().Set("speak", js.FuncOf(speak))

	log.Println("GodSays WASM ready")
	select {} // keep running
}

// TODO: Improve WASM
