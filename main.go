package main

import (
	"fmt"
	"syscall/js"

	"github.com/Trey2k/OpenStreaming/overlayWASM/overlay"
)

func main() {
	fmt.Println("Overlay Client Starting2")

	js.Global().Set("startOverlay", js.FuncOf(startOverlay))

	// hold
	select {}

}

func startOverlay(this js.Value, args []js.Value) interface{} {
	if len(args) != 3 {
		fmt.Println("Invalid number of arguments for startOverlay(), expected 3, got", len(args))
		return nil
	}

	fmt.Println("Starting overlay")

	overlay := overlay.NewOverlay(args[0].String(), args[1].String(), args[2].Bool())

	fmt.Printf("Overlay ID: %d\n", overlay.ID)
	return nil
}
