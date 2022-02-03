package overlay

type OverlayModuleType int

const (
	InvalidModule = OverlayModuleType(iota)
	AlertBoxModule
)

type EventType int

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
}

type EventStruct struct {
	Type EventType
	Data TwitchEventStruct
}

const (
	InvalidEvent = EventType(iota)
	TestEvent
	TwitchMessageEvent
	TwitchFollow
)

type TwitchEventStruct struct {
	Channel        string
	DisplayName    string
	ProfilePicture string
	UserID         string
	MessageContent string
}

func NewModule(moduleInfo *ModuleInfo, editorMode bool) Module {
	switch moduleInfo.Type {
	case AlertBoxModule:
		return initAlertBox(moduleInfo, editorMode)
	default:
		return nil
	}
}
