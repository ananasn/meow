package streamcam

import(
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
    "github.com/lazywei/go-opencv/opencv"
	"encoding/base64"
	"os"
	"io"
	"time"
)

func StreamVideo(ws *websocket.Conn, quit chan struct{}) {
	revel.INFO.Printf("%s", "IN GOROUTINE")
	ticker := time.NewTicker(time.Millisecond * 10).C
	for{
		select {
		case <- quit:
			revel.INFO.Printf("%s", "STOP GOROUTINE")
			return
		case <- ticker:
			capture := opencv.NewCameraCapture(0)
			defer capture.Release()
			revel.INFO.Printf("%s", "CAPTURE DONE")
			frame := capture.QueryFrame()
		
			//эта штука не работает
			//img := opencv.EncodeImage(".jpg", frame.ImageData(), 0)
			//и эта
			//size := frame.ImageSize()
			//imgbuff := ((*[1 << 30]byte))(frame.ImageData())[:size]
			//fmt.Println(size)
		
			opencv.SaveImage("frame.jpg", frame, 0)
			img, _:= os.Open("frame.jpg")
			stats, _ := img.Stat()
			size := stats.Size()
			imgbuf := make([]byte, 0)
			chunkbuf := make([]byte, size)
			for{
				_, err := img.Read(chunkbuf)
				if err == io.EOF{
					img.Close()
					break
				}	
				imgbuf = append(imgbuf, chunkbuf...)	
			}

			result := "data:image/jpg;base64," + base64.StdEncoding.EncodeToString(imgbuf)
			websocket.Message.Send(ws, &result)
			
		}
	}
}