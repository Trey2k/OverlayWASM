package overlay

import (
	"fmt"

	"github.com/Trey2k/OpenStreaming/overlayWASM/common"
)

type OverlayStruct struct {
	ID         int
	UserID     int
	Key        string
	ModuleInfo map[int]*ModuleInfo
	Modules    map[int]Module
	NewModules []Module
	Context    *ContextMenu
	editorMode bool
	ws         *WebsocketConn
}

func NewOverlay(host string, key string, editorMode bool) *OverlayStruct {
	overlay := &OverlayStruct{}
	overlay.Key = key
	overlay.editorMode = editorMode
	connStr := fmt.Sprintf("wss://%s/api/overlay/websocket?token=%s", host, overlay.Key)
	overlay.ws = newWebsocket(connStr)
	overlay.ws.MessageHandler(overlay.HandleMessage)

	return overlay
}

func (overlay *OverlayStruct) HandleMessage(msg *MessageStruct) error {
	switch msg.Type {
	// Return is called after we call getOverlay
	case Return:
		newOverlay := msg.Overlay
		overlay.ModuleInfo = msg.Overlay.ModuleInfo

		overlay.ID = newOverlay.ID
		overlay.UserID = newOverlay.UserID
		overlay.Key = newOverlay.Key
		overlay.Start()
		return nil
	case InvalidMessage:
	default:
		return fmt.Errorf("unknown or message type: %d\n%v", msg.Type, msg)
	}
	return nil
}

func (overlay *OverlayStruct) Start() {
	overlay.Modules = make(map[int]Module)
	for k, v := range overlay.ModuleInfo {
		overlay.Modules[k] = NewModule(v, overlay.editorMode)
	}

	if overlay.editorMode {
		overlay.startEditor()
	}
	common.StopLoading()
}
