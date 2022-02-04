package overlay

import (
	"fmt"
	"strconv"
	"syscall/js"

	"github.com/Trey2k/OpenStreaming/overlayWASM/common"
)

type OverlayStruct struct {
	ID         int
	UserID     int
	Key        string
	ModuleInfo map[int]*ModuleInfo
	Modules    map[int]Module
	NewModules map[int]Module
	Context    *ContextMenu

	started    bool
	editorMode bool
	ws         *WebsocketConn
	interval   js.Value
}

func NewOverlay(host string, key string, editorMode bool) *OverlayStruct {
	overlay := &OverlayStruct{}
	overlay.Key = key
	overlay.started = false
	overlay.editorMode = editorMode
	connStr := fmt.Sprintf("wss://%s/api/overlay/websocket?token=%s", host, overlay.Key)
	overlay.ws = newWebsocket(connStr)
	overlay.ws.MessageHandler(overlay.HandleMessage)

	return overlay
}

func (overlay *OverlayStruct) HandleMessage(msg *MessageStruct) error {
	switch msg.Type {
	// Return is called after we call getOverlay
	case OverlayInfo:
		if !overlay.started {
			overlay.Start(msg.Overlay)
			break
		}
		overlay.UpdateInfo(msg.Overlay)

		return nil
	default:
		return fmt.Errorf("unknown or message type: %d\n%v", msg.Type, msg)
	}
	return nil
}

func (overlay *OverlayStruct) UpdateInfo(info *OverlayStruct) {
	overlay.ModuleInfo = info.ModuleInfo

	// Update all module info and delete any modules that need to be deleted
	for k, v := range overlay.Modules {
		moduleInfo, ok := overlay.ModuleInfo[k]
		if ok {
			v.UpdateInfo(moduleInfo)
			continue
		}

		v.Destroy()
		delete(overlay.Modules, k)

	}

	// Create any new modules
	for k, v := range info.ModuleInfo {
		if _, ok := overlay.Modules[k]; ok {
			continue
		}

		overlay.Modules[k] = NewModule(v, overlay.editorMode)
	}

	// Call destroy on the old modules
	for _, v := range overlay.NewModules {
		v.Destroy()
	}

	// Reset NewMOdules slice
	overlay.NewModules = make(map[int]Module)

	// If we are loading stop as clicking save will start the loading screen
	common.StopLoading()
}

func (overlay *OverlayStruct) Start(info *OverlayStruct) {
	overlay.started = true
	fmt.Println("Starting overlay")
	overlay.ModuleInfo = info.ModuleInfo

	overlay.ID = info.ID
	overlay.UserID = info.UserID
	overlay.Key = info.Key

	overlay.Modules = make(map[int]Module)
	overlay.NewModules = make(map[int]Module)
	for k, v := range overlay.ModuleInfo {
		overlay.Modules[k] = NewModule(v, overlay.editorMode)
	}

	if overlay.editorMode {
		overlay.startEditor()
	}

	common.StopLoading()

	overlay.interval = js.Global().Call("setInterval", js.FuncOf(overlay.Update), 50)
}

func (overlay *OverlayStruct) Update(this js.Value, args []js.Value) interface{} {
	for _, v := range overlay.Modules {
		v.Update()
	}

	for _, v := range overlay.Modules {
		v.Update()
	}

	return nil

}

func (overlay *OverlayStruct) deleteModule(this js.Value, args []js.Value) interface{} {
	jquery := js.Global().Get("$")
	elm := jquery.Invoke(overlay.Context.elem)
	idString := elm.Call("attr", "id").String()
	id, err := strconv.Atoi(idString)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return nil
	}

	isNew := elm.Call("is", ".new").Bool()
	var module Module
	if isNew {
		module = overlay.NewModules[id]
		delete(overlay.NewModules, id)
	} else {
		module = overlay.Modules[id]
		delete(overlay.Modules, id)
	}
	delete(overlay.ModuleInfo, id)

	module.Destroy()
	return nil
}
