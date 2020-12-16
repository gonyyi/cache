// (c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>

package cache

type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}
