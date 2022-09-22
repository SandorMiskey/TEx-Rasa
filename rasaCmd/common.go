// region: packages

package rasaCmd

import (
	"log/syslog"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"
)

var (
	Config cfg.Config
	Logger log.Logger
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
