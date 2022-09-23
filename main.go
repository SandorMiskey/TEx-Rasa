// region: packages

package main

import (
	"fmt"
	"log/syslog"
	"os"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"
	"github.com/davecgh/go-spew/spew"

	"github.com/SandorMiskey/TEx-Rasa/rasa"
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
	fs := config.NewFlagSet(config.Name + " " + subCommand)
	fs.Arguments = os.Args[2:]
	fs.Entries = map[string]cfg.Entry{
		"logLevel": {Desc: "Logger min severity", Type: "int", Def: 7},
		"subArgs":  {Desc: "appended to the tail, use when you want to pass something begins with -", Type: "string", Def: ""},
	}

	switch subCommand {
	case "copy":
	case "destroy":
	case "exec":
		fs.Entries["rasaCmd"] = cfg.Entry{Desc: "rasa command", Type: "string", Def: "rasa"}
	case "init":
		fs.Entries["instanceEnabled"] = cfg.Entry{Desc: "instance is enabled or not", Type: "bool", Def: true}
		fs.Entries["instanceRoot"] = cfg.Entry{Desc: "directory where the instances are stored", Type: "string", Def: "/app/instances"}
		fs.Entries["rasaCmd"] = cfg.Entry{Desc: "rasa command", Type: "string", Def: "rasa"}
		fs.Entries["rasaPrompt"] = cfg.Entry{Desc: "choose default options for prompts and suppress warnings", Type: "bool", Def: false}
	case "list":
	case "version":
		fs.Entries["rasaCmd"] = cfg.Entry{Desc: "rasa command", Type: "string", Def: "rasa"}
	default:
		fmt.Println("no such subcommand '" + subCommand + "', usage: " + config.Name + " {init,destroy,list...} [options] [args]")
		os.Exit(1)
	}

	err := fs.ParseCopy()
	if err != nil {
		panic(err)
	}

	// endregion: flag set
	// region: tail

	subArgs := fs.FlagSet.Args()
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

	rasa.Config = config
	rasa.Logger = logger

	// endregion: init modules
	// region: routing

	logger.Out(LOG_DEBUG, spew.Sdump(fs.FlagSet))

	switch subCommand {
	case "copy":
	case "destroy":
		/*
			trash directory from Config
		*/

	case "exec":
		rasa.Exec(subArgs, nil)
	case "init":
		rasa.Init()
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
		rasa.Exec([]string{"--version"}, nil)
	default:
		panic("no such subcommand '" + subCommand)
	}

	// endregion:

}
