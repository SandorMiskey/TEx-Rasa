// region: packages

package main

import (
	"log/syslog"
	"os"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"
)

// endregion: packages
// region: global variables

var (
	Config cfg.Config
	Logger log.Logger
)

// endregion: globals

func main() {

	// region: config and cli flags

	Config = *cfg.NewConfig(os.Args[0])
	fs := Config.NewFlagSet(os.Args[0])
	fs.Entries = map[string]cfg.Entry{
		"httpEnable":  {Desc: "enable http", Type: "bool", Def: true},
		"httpPort":    {Desc: "http port", Type: "int", Def: 8080},
		"httpsEnable": {Desc: "enable http", Type: "bool", Def: true},
		"httpsPort":   {Desc: "http port", Type: "int", Def: 8081},

		"logLevel": {Desc: "Logger min severity", Type: "int", Def: 5},
	}

	err := fs.ParseCopy()
	if err != nil {
		panic(err)
	}

	// endregion: cli flags
	// region: logger

	logLevel := syslog.Priority(Config.Entries["logLevel"].Value.(int))

	Logger = *log.NewLogger()
	defer Logger.Close()
	_, _ = Logger.NewCh(log.ChConfig{Severity: &logLevel})

	_ = Logger.Out(logLevel, *log.ChDefaults.Mark)

	// endregion: logger
	// region:

	// endregion:

}
