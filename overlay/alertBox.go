package overlay

import (
	"fmt"
	"syscall/js"
)

type AlertBox struct {
	ID        int
	Info      *ModuleInfo
	busy      bool
	eventChan chan *EventStruct
	elm       js.Value
}

func newAlertBox(id, top, left, width, height int) *AlertBox {
	return initAlertBox(&ModuleInfo{
		ID:     id,
		Type:   AlertBoxModule,
		Top:    top,
		Left:   left,
		Width:  width,
		Height: height,
		IsNew:  true,
	}, true)
}

func initAlertBox(moduleInfo *ModuleInfo, editorMode bool) *AlertBox {
	alertBox := &AlertBox{
		ID:        moduleInfo.ID,
		Info:      moduleInfo,
		eventChan: make(chan *EventStruct, 1),
	}

	jquery := js.Global().Get("$")
	if jquery.Equal(js.Undefined()) {
		fmt.Println("jquery not found")
		return nil
	}

	jquery.Invoke(".modules").Call("append", fmt.Sprintf(`
	<div class="module" id=%[1]d>
		<div class="alertBox" id=%[1]d>
		</div>
	</div>
`, alertBox.ID))

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

		alertBox.elm.Call("draggable", map[string]interface{}{"containment": ".overlay"})
	}

	return alertBox
}

func (alertBox *AlertBox) EventHandler() *ModuleInfo {
	for {
		if alertBox.busy {
			continue
		}
		event := <-alertBox.eventChan
		switch event.Type {
		case TestEvent:

		}
	}
}

func (alertBox *AlertBox) GetInfo() *ModuleInfo {
	return alertBox.Info
}

func (alertBox *AlertBox) SendEvent(event *EventStruct) {
	alertBox.eventChan <- event
}
