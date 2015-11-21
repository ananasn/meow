package controllers

import (
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
    "meow/app/streamcam"
	"runtime"
)

type WebSockApp struct {
	*revel.Controller
}

func (c WebSockApp) Index() revel.Result {
	revel.INFO.Printf("%s", "GET request recieved")
	host := "127.0.0.1"
	port := 9000
	return c.Render(host, port)
}

func (c WebSockApp) WebSockHandler(user string, ws *websocket.Conn) revel.Result {
	defer ws.Close()
	revel.INFO.Printf("Maxprocs: %d", runtime.GOMAXPROCS(0))
	revel.INFO.Printf("%s", "WS request recieved")
	newmessages := make(chan string)
	quit := make(chan struct{})	
	go streamcam.StreamVideo(ws, quit)
	
	go func() {
		var msg string
		for {
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				close(newmessages)
				return
			}
			newmessages <- msg
		}
	}()
	
	for {
		select{
		case msg := <- newmessages:
			if msg == "" {
				close(quit)
				revel.INFO.Printf("%s", "WS close connection")
				return nil
			}
			revel.INFO.Printf("WS message '%s'", msg)	
		}
	}
}
