package app

import (
	"fmt"
	"meow/app/webterm"

	"github.com/lazywei/go-opencv/opencv"
	"github.com/revel/revel"
	"github.com/tarm/serial"
)

var Capture = make(chan *opencv.Capture, 1)

var PortTerminal *serial.Port
var AddrHTTP string
var PortHTTP string

func GetCapture() chan *opencv.Capture {
	return Capture
}

func init() {

	//Initialize camera
	fmt.Println("Get capture")
	Capture <- opencv.NewCameraCapture(0)
	fmt.Println("Capture done")

	//Open serial connection
	port, err := webterm.StartSerial("/dev/cu.SLAB_USBtoUART", 115200)
	if err != nil {
		fmt.Println("Can'n open serial connection!")
	} else {
		fmt.Println("Serial connection is opened")
		PortTerminal = port
	}
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	// register startup functions with OnAppStart
	// ( order dependent )
	// revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)
	revel.OnAppStart(func() {
		var found bool = false
		AddrHTTP, found = revel.Config.String("http.addr")
		PortHTTP, found = revel.Config.String("http.port")
		if !found {
			panic("http.addr or http.port are not defined in the config section")
		}
	})
}

// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}
