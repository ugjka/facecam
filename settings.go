package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/gotk3/gotk3/gtk"
)

type settings struct {
	window   *gtk.Window
	deviceEn *gtk.Entry
	rateEn   *gtk.Entry
	heightEn *gtk.Entry
	start    *gtk.Button
	device   string
	height   int
	onStart  func()
}

func (s *settings) open() {
	s.initWidgets()
	s.window.Connect("destroy", s.close)
	s.start.Connect("pressed", s.onStart)
	s.fillInputs()
}

func (s *settings) close() {
	s.window.Destroy()
	gtk.MainQuit()
}

func (s *settings) hide() {
	s.window.SetVisible(false)
}

func (s *settings) unhide() {
	s.window.SetVisible(true)
}

func (s *settings) fillInputs() {
	s.deviceEn.SetText(s.device)
	s.heightEn.SetText(fmt.Sprintf("%d", s.height))
}

func (s *settings) validateInputs() (errtop error) {
	var err error
	devreg := regexp.MustCompile("^\\/dev\\/video\\d$")
	dev, err := s.deviceEn.GetText()
	fatal(err)
	if ok := devreg.MatchString(dev); ok {
		s.device = dev
	} else {
		err = fmt.Errorf("Invalid device")
		s.modalError(err)
		return err
	}
	s.height, err = strconv.Atoi(func() string {
		tmp, err := s.heightEn.GetText()
		fatal(err)
		return tmp
	}())
	if err != nil {
		s.modalError(err)
		return err
	}
	return nil
}

func (s *settings) initWidgets() {
	var err error
	s.window, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	fatal(err)
	s.window.SetPosition(gtk.WIN_POS_CENTER)
	s.window.SetTitle("Facecam")
	s.window.SetBorderWidth(6)
	s.window.SetSizeRequest(300, 100)
	s.window.SetIconFromFile("../share/icons/hicolor/256x256/apps/facecam.png")
	grid, err := gtk.GridNew()
	fatal(err)
	grid.SetColumnHomogeneous(true)
	grid.SetColumnSpacing(6)
	grid.SetRowSpacing(6)
	config, err := gtk.FrameNew("Configuration:")
	fatal(err)
	config.SetBorderWidth(6)
	grid.Attach(config, 0, 0, 1, 1)
	grid2, err := gtk.GridNew()
	fatal(err)
	grid2.SetColumnHomogeneous(true)
	grid2.SetColumnSpacing(6)
	grid2.SetRowSpacing(6)
	grid2.SetBorderWidth(6)
	grid2.Attach(labelNew("Video Device:"), 0, 0, 1, 1)
	s.deviceEn, err = gtk.EntryNew()
	fatal(err)
	grid2.Attach(s.deviceEn, 0, 1, 1, 1)
	grid2.Attach(labelNew("Height:"), 0, 4, 1, 1)
	s.heightEn, err = gtk.EntryNew()
	fatal(err)
	grid2.Attach(s.heightEn, 0, 5, 1, 1)
	config.Add(grid2)
	s.start, err = gtk.ButtonNew()
	fatal(err)
	s.start.SetLabel("Start")
	grid.Attach(s.start, 0, 1, 1, 1)
	s.window.Add(grid)
	s.window.ShowAll()
}

func (s *settings) modalError(err error) {
	msg := gtk.MessageDialogNew(s.window, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR,
		gtk.BUTTONS_CLOSE, "%s", err.Error())
	_, err = msg.Connect("response", func() {
		msg.Destroy()
	})
	fatal(err)
	msg.ShowAll()
}
