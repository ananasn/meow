package controllers

import (
	"fmt"
	"meow/app"
	"meow/app/streamcam"
	"meow/app/webterm"
	"strconv"

	"github.com/revel/revel"
	"golang.org/x/net/websocket"
)

type WebSockApp struct {
	*revel.Controller
}

func (c WebSockApp) Index() revel.Result {
	revel.INFO.Printf("%s", "GET request recieved")
	host := app.AddrHTTP
	port := app.PortHTTP
	return c.Render(host, port)
}

func (c WebSockApp) StreamWidget() revel.Result {
	revel.INFO.Printf("%s", "GET request recieved")
	host := app.AddrHTTP
	port := app.PortHTTP
	return c.Render(host, port)
}

func (c WebSockApp) AxisMovementsWidget() revel.Result {
	revel.INFO.Printf("%s", "GET request recieved")
	host := app.AddrHTTP
	port := app.PortHTTP
	return c.Render(host, port)
}

func (c WebSockApp) CoordInfoWidget() revel.Result {
	revel.INFO.Printf("%s", "GET request recieved")
	host := app.AddrHTTP
	port := app.PortHTTP
	return c.Render(host, port)
}
func (c WebSockApp) GCodeWidget() revel.Result {
	revel.INFO.Printf("%s", "GET request recieved")
	host := app.AddrHTTP
	port := app.PortHTTP
	return c.Render(host, port)
}

func (c WebSockApp) WebSockHandler(user string, ws *websocket.Conn) revel.Result {
	defer ws.Close()
	revel.INFO.Printf("%s", "WS request recieved")
	newmessages := make(chan string)
	quit := make(chan struct{})
	go streamcam.StreamVideo(ws, quit, app.GetCapture())

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
		select {
		case msg := <-newmessages:
			if msg == "" {
				close(quit)
				revel.INFO.Printf("%s", "WS close connection")
				return nil
			}
			revel.INFO.Printf("WS message '%s'", msg)
		}
	}
}

func (c WebSockApp) TerminalHendler(user string, ws *websocket.Conn) revel.Result {
	defer ws.Close()
	revel.INFO.Printf("%s", "WS terminal opened")
	//init html pattern (terminal table)
	if webterm.Doc == nil {
		revel.INFO.Printf("%s", "Create tree")
		doc, err := webterm.HTMLToNodesTree(revel.ViewsPath + "/WebSockApp/index.html")
		if err != nil {
			revel.INFO.Printf("%s", "Err during tree creation")
			return nil
		}
		webterm.Doc = doc
	}
	newmessages := make(chan int)
	quit := make(chan struct{})
	go func() {
		var msg string
		for {
			revel.INFO.Printf("WS-terminal listen in goroutine")
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				close(newmessages)
				return
			}
			int_msg, _ := strconv.Atoi(msg)
			newmessages <- int_msg
		}
	}()
	go func() { //TODO: Make refactoring
		for {
			select {
			case <-webterm.NewLineChan:
				revel.INFO.Printf("Move to new line")
				webterm.SerialWrite(byte(10), app.PortTerminal)
			}
		}
	}()
	go func() {
		for {
			select {
			case <-quit:
				revel.INFO.Printf("%s", "Goroutine closed")
				return
			default:
				str, err := webterm.SerialRead(app.PortTerminal)
				if err == nil {
					res := webterm.EscStringToHTML(str)
					fmt.Printf(res[len(res)-50 : len(res)-1])
					websocket.Message.Send(ws, res)
				}
			}
		}
	}()
	for {
		select {
		case msg := <-newmessages:
			if msg == 0 {
				revel.INFO.Printf("%s", "WS-terminal close connection")
				close(quit)
				return nil
			}
			revel.INFO.Printf("WS-terminal recieve message %d ", msg)
			err := webterm.SerialWrite(byte(msg), app.PortTerminal)
			if err != nil {
				revel.INFO.Printf("WS-terminal error while send message %d ", msg)
			}
		}
	}
}
