package cache

import (
	"github.com/gonyyi/alog"
	"io"
)

type cacheLogger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

// =====================================================================================================================
// CUSTOM -- ALOG adopter
// =====================================================================================================================

// Create custom level
const (
	LvDEBUG alog.Level = 1 << iota
	LvINFO
	LvWARN
	LvERROR
	LvFATAL
	LvALL = LvDEBUG | LvINFO | LvWARN | LvERROR | LvFATAL
)

func NewCacheLog(out io.Writer) *cacheLog {
	c := cacheLog{}
	c.SetPrefix("cache ")
	c.SetFlag(alog.F_STD)
	c.SetOutput(out)
	c.LvOverride(LvALL)
	// c.LvDisable(LvDEBUG)
	return &c
}

type cacheLog struct {
	alog.ALogger
}

func (c *cacheLog) Debugf(fmt string, a ...interface{}) {
	c.Printfl(LvDEBUG, "[DEBUG] "+fmt, a...)
}
func (c *cacheLog) Infof(fmt string, a ...interface{}) {
	c.Printfl(LvINFO, "[INFO]  "+fmt, a...)
}
func (c *cacheLog) Warnf(fmt string, a ...interface{}) {
	c.Printfl(LvWARN, "[WARN]  "+fmt, a...)
}
func (c *cacheLog) Errorf(fmt string, a ...interface{}) {
	c.Printfl(LvERROR, "[ERROR] "+fmt, a...)
}
func (c *cacheLog) Fatalf(fmt string, a ...interface{}) {
	c.Printfl(LvERROR, "[FATAL] "+fmt, a...)
	// Terminate program
	// os.Exit(1)
}
