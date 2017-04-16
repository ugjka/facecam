package main

import "github.com/gotk3/gotk3/gtk"

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}

func labelNew(s string) *gtk.Label {
	l, err := gtk.LabelNew(s)
	fatal(err)
	l.SetHAlign(gtk.ALIGN_START)
	return l
}
