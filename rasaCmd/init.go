// region: packages

package rasaCmd

import (
	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"
)

var (
	Config cfg.Config
	Logger log.Logger
)

func init() {
}

// region: loglevel

// switch logLevel := syslog.Priority(Config.Entries["logLevel"].Value.(int)); {
// case logLevel <= log.LOG_WARNING:
// 	// set -quiet
// case logLevel >= log.LOG_INFO:
// 	// set -v
// case logLevel == log.LOG_DEBUG:
// 	// set -vv
// default:
// 	// set nothing
// }

// endregion: loglevel
