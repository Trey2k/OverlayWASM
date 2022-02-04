package overlay

import (
	"fmt"
	"strconv"
	"syscall/js"
)

type AlertBox struct {
	ID     int
	Info   *ModuleInfo
	busy   bool
	events []*EventStruct
	elm    js.Value
	isNew  bool
}

func newAlertBox(id, top, left, width, height int, isNew bool) *AlertBox {
	return initAlertBox(&ModuleInfo{
		ID:     id,
		Type:   AlertBoxModule,
		Top:    top,
		Left:   left,
		Width:  width,
		Height: height,
		IsNew:  isNew,
	}, true)
}

func initAlertBox(moduleInfo *ModuleInfo, editorMode bool) *AlertBox {
	alertBox := &AlertBox{
		ID:    moduleInfo.ID,
		Info:  moduleInfo,
		isNew: moduleInfo.IsNew,
	}

	jquery := js.Global().Get("$")
	newClass := ""
	if alertBox.isNew {
		newClass = "new"
	}

	jquery.Invoke(".modules").Call("append", fmt.Sprintf(`
	<div class="module %[2]s" id=%[1]d>
		<div class="alertBox %[2]s" id=%[1]d>
		</div>
	</div>
`, alertBox.ID, newClass))

	alertBox.elm = jquery.Invoke(fmt.Sprintf(".module#%d", alertBox.ID))

	alertBox.elm.Call("css", "top", moduleInfo.Top)
	alertBox.elm.Call("css", "left", moduleInfo.Left)
	alertBox.elm.Call("css", "width", moduleInfo.Width)
	alertBox.elm.Call("css", "height", moduleInfo.Height)

	if editorMode {
		alertBox.elm.Call("prepend",
			fmt.Sprintf(`
				<span class="moduleTitle">%[2]s</span>
				<div class="ui-resizable-handle ui-resizable-nw nw%[1]d"></div>
				<div class="ui-resizable-handle ui-resizable-ne ne%[1]d"></div>
				<div class="ui-resizable-handle ui-resizable-sw sw%[1]d"></div>
				<div class="ui-resizable-handle ui-resizable-se se%[1]d"></div>
				<div class="ui-resizable-handle ui-resizable-n n%[1]d"></div>
				<div class="ui-resizable-handle ui-resizable-s s%[1]d"></div>
				<div class="ui-resizable-handle ui-resizable-e e%[1]d"></div>
				<div class="ui-resizable-handle ui-resizable-w w%[1]d"></div>
				`, alertBox.ID, "Alert Box"))

		alertBox.elm.Call("resizable", map[string]interface{}{
			"containment": ".overlay",
			"scroll":      false,
			"stop":        js.FuncOf(alertBox.onStop),
			"handles": map[string]interface{}{
				"nw": fmt.Sprintf(".nw%d", alertBox.ID),
				"ne": fmt.Sprintf(".ne%d", alertBox.ID),
				"sw": fmt.Sprintf(".sw%d", alertBox.ID),
				"se": fmt.Sprintf(".se%d", alertBox.ID),
				"n":  fmt.Sprintf(".n%d", alertBox.ID),
				"s":  fmt.Sprintf(".s%d", alertBox.ID),
				"e":  fmt.Sprintf(".e%d", alertBox.ID),
				"w":  fmt.Sprintf(".w%d", alertBox.ID),
			},
		})

		alertBox.elm.Call("draggable", map[string]interface{}{
			"stop":        js.FuncOf(alertBox.onStop),
			"containment": ".overlay",
		})
	}

	return alertBox
}

func (alertBox *AlertBox) onStop(this js.Value, args []js.Value) interface{} {
	top := alertBox.elm.Call("css", "top").String()
	alertBox.Info.Top, _ = strconv.Atoi(top[:len(top)-2])
	left := alertBox.elm.Call("css", "left").String()
	alertBox.Info.Left, _ = strconv.Atoi(left[:len(left)-2])
	width := alertBox.elm.Call("css", "width").String()
	alertBox.Info.Width, _ = strconv.Atoi(width[:len(width)-2])
	height := alertBox.elm.Call("css", "height").String()
	alertBox.Info.Height, _ = strconv.Atoi(height[:len(height)-2])
	return nil
}

func (alertBox *AlertBox) Update() {
	if len(alertBox.events) > 0 && !alertBox.busy {
		event := alertBox.events[0]
		alertBox.events = alertBox.events[1:]

		switch event.Type {
		case InvalidEvent:
		}
	}
}

func (alertBox *AlertBox) GetInfo() *ModuleInfo {
	return alertBox.Info
}

func (alertBox *AlertBox) SendEvent(event *EventStruct) {
	alertBox.events = append(alertBox.events, event)
}

func (alertBox *AlertBox) GetElement() js.Value {
	return alertBox.elm
}

func (alertBox *AlertBox) Destroy() {
	alertBox.elm.Call("remove")
	alertBox = nil
}

func (alertBox *AlertBox) UpdateInfo(info *ModuleInfo) {
	alertBox.Info = info
	alertBox.elm.Call("css", "top", info.Top)
	alertBox.elm.Call("css", "left", info.Left)
	alertBox.elm.Call("css", "width", info.Width)
	alertBox.elm.Call("css", "height", info.Height)
}
