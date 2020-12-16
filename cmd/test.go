// (c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>

package main

import (
	"github.com/gonyyi/alog"
	"github.com/gonyyi/alog/alogc"
	"github.com/gonyyi/cache"
	"os"
)

func main() {
	// create new config file
	// err := cache.cache.CreateNewConfig("./config.json", true)

	c := cache.New()
	l := alogc.New(os.Stderr, "", alog.F_STD)
	c.SetLogger(l)

	err := c.Open("./config.json")
	if err != nil {
		println(err.Error())
	}

	c.CachePullAll()
	c.Save()
}
