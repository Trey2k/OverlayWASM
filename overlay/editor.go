package overlay

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/Trey2k/OpenStreaming/overlayWASM/common"
)

func (overlay *OverlayStruct) startEditor() {
	jquery := js.Global().Get("$")
	jquery.Invoke(".sidenav").Call("sidenav")

	elms := js.Global().Get("document").Call("querySelectorAll", ".fixed-action-btn")
	js.Global().Get("M").Get("FloatingActionButton").Call("init", elms,
		map[string]interface{}{
			"direction":    "left",
			"hoverEnabled": true,
		})

	overlay.initContextMenu()

	jquery.Invoke("#save").Call("click", js.FuncOf(overlay.Save))
}

func (overlay *OverlayStruct) Save(this js.Value, args []js.Value) interface{} {
	msg := &EventStruct{
		Type:    SaveOverlay,
		Overlay: overlay,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Error marshalling message:", err)
		return nil
	}

	err = overlay.ws.Send(msgBytes)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return nil
	}

	// Start the loading screen, once server responds with overlayInfo update it will be stopepd
	common.StartLoading()

	return nil
}
