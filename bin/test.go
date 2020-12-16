package main

import (
	"github.com/gonyyi/cache"
)

func main() {
	// create new config file
	// err := cache.NewConfig("./config.json", true)

	c, err := cache.New("./config.json")
	if err != nil {
		println(err.Error())
	}
	c.CachePull("test")
	c.Save()

	b, _ := c.GetCacheData("test")
	println(string(b))
}
