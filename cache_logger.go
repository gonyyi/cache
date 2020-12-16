// (c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>

package cache

type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

// =====================================================================================================================
// DEFAULT DUMMY LOGGER
// =====================================================================================================================
type dummyLog struct {}

func(d *dummyLog) Debugf(a string, b ...interface{}){}
func(d *dummyLog) Infof(a string, b ...interface{}){}
func(d *dummyLog) Warnf(a string, b ...interface{}){}
func(d *dummyLog) Errorf(a string, b ...interface{}){}
func(d *dummyLog) Fatalf(a string, b ...interface{}){}

func newDummyLogger() *dummyLog {
	return &dummyLog{}
}

