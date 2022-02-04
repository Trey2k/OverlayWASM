package overlay

import "syscall/js"

type OverlayModuleType int

const (
	InvalidModule = OverlayModuleType(iota)
	AlertBoxModule
)

type ModuleInfo struct {
	ID     int
	Type   OverlayModuleType
	Top    int
	Left   int
	Width  int
	Height int
	IsNew  bool
}

type Module interface {
	GetInfo() *ModuleInfo
	SendEvent(event *EventStruct)
	GetElement() js.Value
	Destroy()
	Update()
	UpdateInfo(*ModuleInfo)
}

func NewModule(moduleInfo *ModuleInfo, editorMode bool) Module {
	switch moduleInfo.Type {
	case AlertBoxModule:
		return initAlertBox(moduleInfo, editorMode)
	default:
		return nil
	}
}
