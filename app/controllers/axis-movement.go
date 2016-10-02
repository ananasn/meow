package controllers

import (
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"strings"
)

func (c WebSockApp) AxisMovementWidgetWebSocket(user string, ws *websocket.Conn) revel.Result {
	revel.INFO.Printf("%s", "WS_AXIS_MOV request recieved")
	defer ws.Close()

	newmessages := make(chan string)
	//quit := make(chan struct{})
	go func() {
		var msg string
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			close(newmessages)
			return
		}
		newmessages <- msg
	}()
	for {
		select {
		case msg := <-newmessages:
			//TODO:Change printf on cmds to virtual keyboard
			if strings.HasPrefix(msg, "X") {
				revel.INFO.Printf("%s", "X_AXIS cmd recieved")
			}
			if strings.HasPrefix(msg, "Y") {
				revel.INFO.Printf("%s", "Y_AXIS cmd recieved")
			}
			if strings.HasPrefix(msg, "Z") {
				revel.INFO.Printf("%s", "Z_AXIS cmd recieved")
			}
		}
	}
}
