package overlay

import (
	"fmt"
	"syscall/js"
)

type ContextMenu struct {
	locationX int
	locationY int
	elem      js.Value
	menu      js.Value
	addMenu   js.Value
	newIndex  int
}

func (overlay *OverlayStruct) initContextMenu() {
	jquery := js.Global().Get("$")
	overlay.Context = &ContextMenu{
		newIndex: -1000, // Make sure ids do not conflict
	}

	js.Global().Set("deleteContextModule", js.FuncOf(overlay.deleteModule))

	jquery.Invoke("#addAlertBox").Call("click",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			overlay.NewModules[overlay.Context.newIndex] = newAlertBox(overlay.Context.newIndex,
				overlay.Context.locationY-175, overlay.Context.locationX-100, 200, 350, true)

			info := overlay.NewModules[overlay.Context.newIndex].GetInfo()
			overlay.ModuleInfo[info.ID] = info

			overlay.Context.newIndex++
			return nil
		}))

	overlay.Context.menu = jquery.Invoke(".context-trigger").Call("dropdown",
		js.ValueOf(map[string]interface{}{
			"constain_width": true,
			"belowOrigin":    true,
			"alignment":      "left",
		}))

	overlay.Context.addMenu = jquery.Invoke(".addMenu-trigger").Call("dropdown",
		map[string]interface{}{
			"inDuration":      300,
			"outDuration":     225,
			"constrain_width": true,
			"hover":           true,
			"belowOrigin":     true,
			"alignment":       "left",
		})

	js.Global().Get("document").Call("addEventListener", "contextmenu",
		js.FuncOf(overlay.Context.OpenContextMenu))
}

func (context *ContextMenu) OpenContextMenu(this js.Value, args []js.Value) interface{} {
	event := args[0]
	jquery := js.Global().Get("$")
	context.menu.Call("dropdown", "close")
	context.hideModuleContext()

	if jquery.Invoke(event.Get("target")).Call("is", ".module").Bool() ||
		jquery.Invoke(event.Get("target")).Call("is", ".alertBox").Bool() {
		context.elem = jquery.Invoke(event.Get("target"))
		context.showModuleContext()
	} else if !context.elem.Equal(js.Undefined()) {
		context.elem = js.Undefined()
	}

	context.locationX = event.Get("clientX").Int() + 5
	context.locationY = event.Get("clientY").Int() + 5

	context.menu.Call("css", "left", fmt.Sprintf("%dpx", context.locationX))
	context.menu.Call("css", "top", fmt.Sprintf("%dpx", context.locationY))

	context.menu.Call("dropdown", "open")
	event.Call("preventDefault")
	return nil
}

func (context *ContextMenu) showModuleContext() {
	jquery := js.Global().Get("$")
	jquery.Invoke("#contextMenu").Call("prepend",
		`<li class="contextModule" id="delete" onclick="deleteContextModule()"><a class="red-text text-darken-1"><i class="material-icons">cancel</i>Delete</a></li>`)
}

func (context *ContextMenu) hideModuleContext() {
	jquery := js.Global().Get("$")
	jquery.Invoke("#contextMenu").Call("children").Call("remove", ".contextModule")
}
