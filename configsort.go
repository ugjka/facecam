package main

import (
	"github.com/korandiz/v4l"
)

type configs []v4l.DeviceConfig

func (c configs) Len() int {
	return len(c)
}

func (c configs) Less(i, j int) bool {
	if c[i].Height == c[j].Height {
		if c[i].FPS.Cmp(c[j].FPS) == -1 {
			return true
		}
		return false
	}
	return c[i].Height < c[j].Height
}

func (c configs) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
