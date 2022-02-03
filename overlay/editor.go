package overlay

import (
	"syscall/js"
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

	jquery.Invoke("#addAlertBox").Call("click",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			overlay.NewModules = append(overlay.NewModules, newAlertBox(0, 0, 0, 350, 200))
			return nil
		}))

	overlay.initContextMenu()

}
