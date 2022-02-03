package overlay

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

type MessageHandler func(msg *MessageStruct) error

type MessageType int

type WebsocketConn struct {
	ws         js.Value
	MsgHandler MessageHandler
	Ready      bool
}

type MessageStruct struct {
	Type    MessageType
	Overlay *OverlayStruct
}

const (
	InvalidMessage = MessageType(iota)
	GetOverlay
	Return
)

func newWebsocket(conURL string) *WebsocketConn {
	conn := &WebsocketConn{}
	conn.ws = js.Global().Get("WebSocket").New(conURL)
	conn.ws.Call("addEventListener", "open", js.FuncOf(conn.Open))
	conn.ws.Call("addEventListener", "message", js.FuncOf(conn.Message))

	return conn
}

func (conn *WebsocketConn) MessageHandler(handler MessageHandler) {
	conn.MsgHandler = handler
}

func (conn *WebsocketConn) Open(this js.Value, args []js.Value) interface{} {
	fmt.Println("WS Connection Opened")
	conn.Ready = true
	toSend := &MessageStruct{
		Type: GetOverlay,
	}

	jsonBytes, err := json.Marshal(toSend)
	if err != nil {
		fmt.Println("Error marshalling message:", err)
		return nil
	}

	conn.Send(jsonBytes)
	return nil
}

func (conn *WebsocketConn) Error(this js.Value, args []js.Value) interface{} {
	fmt.Printf("WS Error: %s\n", args[0].String())
	return nil
}

func (conn *WebsocketConn) Message(this js.Value, args []js.Value) interface{} {
	event := args[0]
	data := event.Get("data").String()

	msg := &MessageStruct{}
	err := json.Unmarshal([]byte(data), msg)
	if err != nil {
		fmt.Println("Error unmarshalling message:", err)
		return nil
	}

	if conn.MsgHandler != nil {
		err = conn.MsgHandler(msg)
		if err != nil {
			fmt.Println("Error handling message:", err)
			return nil
		}
	}
	return nil
}

func (ws *WebsocketConn) Send(data []byte) error {
	ws.ws.Call("send", string(data))
	return nil
}

func (conn *WebsocketConn) Close() {
	conn.ws.Call("close")
}

func (ws *WebsocketConn) Write(data []byte) error {
	return nil
}
