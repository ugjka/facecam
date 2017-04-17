package main

import (
	"image"
	"runtime"
	"sync"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/korandiz/v4l"
	"github.com/korandiz/v4l/fmt/yuyv"
)

type video struct {
	window     *gtk.Window
	image      *gtk.Image
	pixbuff    *gdk.Pixbuf
	eventbox   *gtk.EventBox
	width      int
	height     int
	configs    configs
	config     v4l.DeviceConfig
	device     *v4l.Device
	devicePath string
	onclose    func()
	stop       chan bool
	wg         sync.WaitGroup
	tmpbuff    *gdk.Pixbuf
	ticker     *time.Ticker
	img        *yuyv.Image
	rgba       *image.RGBA
	icon       *gdk.Pixbuf
}

func (w *video) open() {
	var err error
	if err := w.getDevice(); err != nil {
		w.modalError(err)
		w.close()
		return
	}
	w.width = int(float64(w.height) * (float64(w.config.Width) / float64(w.config.Height)))
	w.initWidgets()
	w.wg.Add(1)
	w.startStream()
	_, err = w.window.Connect("destroy", w.close)
	fatal(err)
	_, err = w.eventbox.Connect("button-press-event", func(b *gtk.EventBox, e *gdk.Event) {
		if gdk.EventKeyNewFromEvent(e).Type() == gdk.EVENT_3BUTTON_PRESS {
			w.window.Destroy()
		}
	})
	fatal(err)
}

func (w *video) openDevice() error {
	return w.device.TurnOn()
}

func (w *video) startStream() {
	binfo, err := w.device.BufferInfo()
	fatal(err)
	w.img = &yuyv.Image{
		Pix:    make([]byte, binfo.BufferSize),
		Stride: binfo.ImageStride,
		Rect:   image.Rect(0, 0, 640, 480),
	}
	w.rgba = &image.RGBA{
		Pix:    w.pixbuff.GetPixels(),
		Stride: w.pixbuff.GetRowstride(),
		Rect:   image.Rect(0, 0, w.pixbuff.GetWidth(), w.pixbuff.GetHeight()),
	}
	go func(w *video) {
		defer w.wg.Done()
		for {
			select {
			case <-w.stop:
				return
			case <-w.ticker.C:
				runtime.GC()
			default:
			}
			buffer, err := w.device.Capture()
			fatal(err)
			buffer.ReadAt(w.img.Pix, 0)
			yuyv.ToRGBA(w.rgba, w.rgba.Rect, w.img, w.img.Rect.Min)
			_, err = glib.IdleAdd(func() bool {
				w.update()
				return true
			})
			fatal(err)

		}
	}(w)
}

func (w *video) getDevice() error {
	var err error
	w.device, err = v4l.Open(w.devicePath)
	if err != nil {
		return err
	}
	w.configs, err = w.device.ListConfigs()
	if err != nil {
		return err
	}
	w.config = w.configs[0]
	if err := w.device.SetConfig(w.config); err != nil {
		return err
	}
	if err := w.openDevice(); err != nil {
		return err
	}
	return nil
}

func (w *video) close() {
	if w.window != nil {
		w.window.Destroy()
	}
	if w.device != nil {
		w.stop <- true
		w.wg.Wait()
		w.device.TurnOff()
		w.device.Close()
	}
	w.onclose()
}

func (w *video) initWidgets() {
	var err error
	w.window, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	fatal(err)
	w.window.SetKeepAbove(true)
	w.window.SetDecorated(false)
	w.window.SetSizeRequest(w.width, w.height)
	w.window.SetPosition(gtk.WIN_POS_CENTER)
	w.window.SetResizable(true)
	w.window.SetTitle("facecam")
	w.window.SetIcon(w.icon)
	w.pixbuff, err = gdk.PixbufNew(gdk.COLORSPACE_RGB, true, 8, w.config.Width, w.config.Height)
	fatal(err)
	w.image, err = gtk.ImageNew()
	fatal(err)
	w.eventbox, err = gtk.EventBoxNew()
	fatal(err)
	w.eventbox.Add(w.image)
	w.window.Add(w.eventbox)
	w.window.ShowAll()
	screen, err := gdk.ScreenGetDefault()
	fatal(err)
	w.window.Move(screen.GetWidth(), screen.GetHeight())
}

func (w *video) update() {
	var err error
	w.tmpbuff, err = w.pixbuff.ScaleSimple(w.width, w.height, gdk.INTERP_NEAREST)
	fatal(err)
	w.image.SetFromPixbuf(w.tmpbuff)
}

func (w *video) modalError(err error) {
	msg := gtk.MessageDialogNew(w.window, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR,
		gtk.BUTTONS_CLOSE, "%s", err.Error())
	_, err = msg.Connect("response", func() {
		msg.Destroy()
	})
	fatal(err)
	msg.SetKeepAbove(true)
	msg.ShowAll()
}
