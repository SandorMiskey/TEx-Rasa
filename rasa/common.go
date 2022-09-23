// region: packages

package rasa

import (
	"errors"
	"log/syslog"
	"os"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"
)

var (
	Config cfg.Config
	Logger log.Logger

	ErrInvalidInstanceRoot = errors.New("invalid instance root dir")
)

func init() {
}

func LogLevel() (flag string) {
	switch logLevel := syslog.Priority(Config.Entries["logLevel"].Value.(int)); {
	case logLevel == log.LOG_DEBUG:
		flag = "-vv"
	case logLevel >= log.LOG_INFO:
		flag = "-v"
	case logLevel <= log.LOG_WARNING:
		flag = "-quiet"
	default:
		flag = ""
	}

	Logger.Out(log.LOG_EMERG, "rasaCmd LogLevel() set to", flag)
	return
}

func Wd() (err error) {
	wd, ok := Config.Entries["instanceRoot"].Value.(string)
	if !ok {
		err = ErrInvalidInstanceRoot
		Logger.Out(log.LOG_ERR, err)
		return
	}
	if err = os.Chdir(wd); err != nil {
		Logger.Out(log.LOG_ERR, err)
		return
	}
	if _, err = os.Getwd(); err != nil {
		Logger.Out(log.LOG_ERR, err)
		return
	}
	Logger.Out(log.LOG_DEBUG, "rasa.Wd() directory set to", wd)
	return
}
