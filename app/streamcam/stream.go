package streamcam

import (
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"github.com/lazywei/go-opencv/opencv"
	"encoding/base64"
	"os"
	"io"
	"time"
)

func StreamVideo(ws *websocket.Conn, quit chan struct{}, capturechan chan *opencv.Capture) {
	revel.INFO.Printf("%s", "Access capture")
	capture := <-capturechan
	capturechan <- capture
	ticker := time.NewTicker(time.Millisecond * 10)
	for {
		select {
		case <- quit:
			revel.INFO.Printf("%s", "Goroutine closed")
			return
		case <- ticker.C:
			frame := capture.QueryFrame()
		
			//img := opencv.EncodeImage(".jpg", frame.ImageData(), 0)
			//
			//size := frame.ImageSize()
			//imgbuff := ((*[1 << 30]byte))(frame.ImageData())[:size]
			//fmt.Println(size)
		
			opencv.SaveImage("frame.jpg", frame, 0)
			img, _ := os.Open("frame.jpg")
			stats, _ := img.Stat()
			size := stats.Size()
			imgbuf := make([]byte, 0)
			chunkbuf := make([]byte, size)
			for {
				_, err := img.Read(chunkbuf)
				if err == io.EOF{
					img.Close()
					break
				}	
				imgbuf = append(imgbuf, chunkbuf...)	
			}

			result := "data:image/jpg;base64," + base64.StdEncoding.EncodeToString(imgbuf)
			err := websocket.Message.Send(ws, result)	
			if err != nil {
				ticker.Stop()
				revel.INFO.Printf("%s %s", "ERROR:", err)
			}
		}
	}
}