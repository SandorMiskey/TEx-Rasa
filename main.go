// region: packages

package main

import (
	"fmt"
	"log/syslog"
	"os"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"
	"github.com/davecgh/go-spew/spew"

	"github.com/SandorMiskey/TEx-Rasa/rasaCmd"
)

// endregion: packages
// region: global variables

// var (
// 	config cfg.Config = *cfg.NewConfig(os.Args[0])
// 	logger log.Logger = *log.NewLogger()
// )

const (
	LOG_EMERG  syslog.Priority = log.LOG_EMERG
	LOG_ERR    syslog.Priority = log.LOG_ERR
	LOG_NOTICE syslog.Priority = log.LOG_NOTICE
	LOG_INFO   syslog.Priority = log.LOG_INFO
	LOG_DEBUG  syslog.Priority = log.LOG_DEBUG
)

// endregion: globals

func main() {

	// region: cli flags

	// region: subcommand

	if len(os.Args) < 2 {
		fmt.Println("subcommand expected")
		os.Exit(1)
	}
	subCommand := os.Args[1]

	// endregion: subcommand
	// region: flag set

	config := *cfg.NewConfig(os.Args[0])
	flagSet := config.NewFlagSet(config.Name + " " + subCommand)
	flagSet.Arguments = os.Args[2:]
	flagSet.Entries = map[string]cfg.Entry{
		// "instanceRoot": {Desc: "directory where the instances are stored", Type: "string", Def: "/app/instances"},
		"logLevel": {Desc: "Logger min severity", Type: "int", Def: 7},
		"subArgs":  {Desc: "appended to the tail, use when you want to pass something begins with -", Type: "string", Def: ""},
	}

	switch subCommand {
	case "copy":
	case "destroy":
	case "exec":
		flagSet.Entries["rasaCmd"] = cfg.Entry{Desc: "rasa command", Type: "string", Def: "rasa"}
	case "init":
		flagSet.Entries["rasaCmd"] = cfg.Entry{Desc: "rasa command", Type: "string", Def: "rasa"}
	case "list":
	case "version":
		flagSet.Entries["rasaCmd"] = cfg.Entry{Desc: "rasa command", Type: "string", Def: "rasa"}
	default:
		fmt.Println("no such subcommand '" + subCommand + "', usage: " + config.Name + " {init,destroy,list...} [options] [args]")
		os.Exit(1)
	}

	err := flagSet.ParseCopy()
	if err != nil {
		panic(err)
	}

	// endregion: flag set
	// region: tail

	subArgs := flagSet.FlagSet.Args()
	if subSubArgs := config.Entries["subArgs"].Value.(string); len(subSubArgs) != 0 {
		subArgs = append(subArgs, subSubArgs)
	}

	// endregion: tail

	// endregion: cli
	// region: logger

	logger := *log.NewLogger()
	defer logger.Close()

	logLevel := syslog.Priority(config.Entries["logLevel"].Value.(int))
	_, _ = logger.NewCh(log.ChConfig{Severity: &logLevel})

	// Logger.Out(LOG_DEBUG, spew.Sdump(Config))

	// endregion: logger
	// region: init modules

	rasaCmd.Config = config
	rasaCmd.Logger = logger

	// endregion: init modules
	// region: routing

	logger.Out(LOG_DEBUG, spew.Sdump(flagSet.FlagSet))

	switch subCommand {
	case "copy":
	case "destroy":
		/*
			trash directory from Config
		*/

	case "exec":
		rasaCmd.Exec(subArgs, nil)
	case "init":
		/*
			--init_dir (root/instance)
			--no_prompt

			rasaCMD
			instanceRoot
			instanceEnabled
			subArg = list of instances
		*/

		rasaCmd.Exec([]string{"init", "-h", rasaCmd.LogLevel()}, nil)
	case "list":
		/*
			files, err := ioutil.ReadDir(Config.Entries["instanceRoot"].Value.(string))
			if err != nil {
				log.Fatal(err)
			}

			for _, f := range files {
				fmt.Println(f.Name())
			}
			instance.List()
		*/
	case "version":
		rasaCmd.Exec([]string{"--version"}, nil)
	default:
		panic("no such subcommand '" + subCommand)
	}

	// endregion:

}
