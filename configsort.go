package main

import (
	"sort"

	"github.com/korandiz/v4l"
	"github.com/korandiz/v4l/fmt/yuyv"
)

type configs []v4l.DeviceConfig

func (c configs) createMap() (m map[uint32]configs) {
	m = make(map[uint32]configs)
	sort.Sort(c)
	for _, v := range c {
		if v.Format != yuyv.FourCC {
			continue
		}
		m[v.FPS.N/v.FPS.D] = append(m[v.FPS.N/v.FPS.D], v)
	}
	return m
}

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
