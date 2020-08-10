package log

import (
	"fmt"
	"io"
	"learn/mywebcrawler/helper/log/base"
	"sync"
)

type LoggerCreator func(
	level base.LogLevel,
	format base.LogFormat,
	writer io.Writer,
	options []base.Option) base.MyLogger

var loggerCreatorMap = map[base.LoggerType]LoggerCreator{}

var rwm sync.RWMutex

func RegisterLogger(
	loggerType base.LoggerType,
	creator LoggerCreator,
	cover bool) error {
	if loggerType == "" {
		return fmt.Errorf("logger register error: invaild logger type")
	}
	if creator == nil {
		return fmt.Errorf("logger register error: invaild logger creator (logger type:%s)", loggerType)
	}
	rwm.Lock()
	defer rwm.Unlock()
	if _, ok := loggerCreatorMap[loggerType]; ok || !cover {
		return fmt.Errorf("logger register error: already existing logger for type %v", loggerType)
	}
	loggerCreatorMap[loggerType] = creator
	return nil
}
