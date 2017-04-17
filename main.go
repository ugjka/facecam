//
//Facecam for poor man's stream setup
//
package main

import (
	"sync"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

var icon = "./facecam.png"

func main() {
	var vid video
	vid.stop = make(chan bool)
	vid.wg = sync.WaitGroup{}
	vid.ticker = time.NewTicker(time.Second * 5)
	var set settings
	set.device = "/dev/video0"
	set.height = 192
	set.onStart = func() {
		err := set.validateInputs()
		if err != nil {
			return
		}
		vid.height = set.height
		//vid.width = int(float64(set.height) * 1.333333333)
		vid.devicePath = set.device
		vid.icon, _ = set.window.GetIcon()
		set.hide()
		vid.open()
	}
	vid.onclose = func() {
		set.unhide()
	}
	gtk.Init(nil)
	set.open()
	gtk.Main()
}
