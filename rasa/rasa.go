// region: packages

package rasa

import (
	"log/syslog"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"
)

var (
	Logger log.Logger
)

func init() {
}

func logLevel(c cfg.Config) (flag string) {

	level, ok := c.Entries["logLevel"].Value.(int)
	if !ok {
		Logger.Out(log.LOG_WARNING, "invalid loglevel, falling back to -vv")
		level = 7
	}

	switch logLevel := syslog.Priority(level); {
	case logLevel == log.LOG_DEBUG:
		flag = "-vv"
	case logLevel >= log.LOG_INFO:
		flag = "-v"
	case logLevel <= log.LOG_WARNING:
		flag = "-quiet"
	default:
		flag = ""
	}

	Logger.Out(log.LOG_DEBUG, "rasaCmd LogLevel() set to", flag)
	return
}
