// region: packages

package main

import (
	"fmt"
	"log/syslog"
	"os"
	"regexp"
	"strings"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"

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

	// TODO: list subcommands

	re := regexp.MustCompile(`.*/`)
	cmd := re.ReplaceAllString(os.Args[0], "")

	help := func() {
		fmt.Println("one of these subcommands expected:")
		fmt.Println("	- instanceList: blah blah blah...")
		fmt.Println("	- rasaExec")
		fmt.Println("	- rasaInit")
		fmt.Println("	- rasaVersion")
		fmt.Println("")
		fmt.Println("use `" + cmd + " subcommand --help` for further info")
		os.Exit(1)
	}
	if len(os.Args) < 2 || strings.HasPrefix(os.Args[1], "-") {
		help()
	}

	// endregion: subcommand
	// region: flag set

	config := *cfg.NewConfig(cmd)
	fs := config.NewFlagSet(config.Name + " " + os.Args[1])
	fs.Arguments = os.Args[2:]
	fs.Entries = map[string]cfg.Entry{
		"consoleEnabled":   {Desc: "writes standard output to console", Type: "bool", Def: false},
		"logFileEnabled":   {Desc: "writes logs to file", Type: "bool", Def: false},
		"logFileLevel":     {Desc: "min severity for logs written into file", Type: "int", Def: 7},
		"logFileOutput":    {Desc: "logFileEnabled destination", Type: "string", Def: "./log"},
		"logStdoutEnabled": {Desc: "writes logs to stdout", Type: "bool", Def: false},
		"logStdoutLevel":   {Desc: "min severity for logs written to stdout", Type: "int", Def: 7},
		"logStderrEnabled": {Desc: "writes logs to stderr", Type: "bool", Def: true},
		"logStderrLevel":   {Desc: "min severity for logs written to stderr", Type: "int", Def: 7},
		"logSyslogEnabled": {Desc: "writes logs to local syslog", Type: "bool", Def: true},
		"logSyslogLevel":   {Desc: "min severity for logs written to syslog", Type: "int", Def: 7},
		"subArgs":          {Desc: "appended to the tail, use when you want to pass something begins with - (or use the -- separator)", Type: "string", Def: ""},
	}

	subCmd := strings.ToLower(os.Args[1])
	switch subCmd {
	case "instancelist":
		fs.Entries["instanceRoot"] = cfg.Entry{Desc: "directory where the instances are stored", Type: "string", Def: "/app/instances"}
	case "rasaexec":
		fs.Entries["rasaCmd"] = cfg.Entry{Desc: "rasa command", Type: "string", Def: "rasa"}
	// case "rasainit":
	// 	fs.Entries["instanceEnabled"] = cfg.Entry{Desc: "instance is enabled or not", Type: "bool", Def: true}
	// 	fs.Entries["instanceRoot"] = cfg.Entry{Desc: "directory where the instances are stored", Type: "string", Def: "/app/instances"}
	// 	fs.Entries["rasaCmd"] = cfg.Entry{Desc: "rasa command", Type: "string", Def: "rasa"}
	// 	fs.Entries["rasaPrompt"] = cfg.Entry{Desc: "choose default options for prompts and suppress warnings (DEPRECATED!)", Type: "bool", Def: false}
	case "rasaversion":
		fs.Entries["rasaCmd"] = cfg.Entry{Desc: "rasa command", Type: "string", Def: "rasa"}
	default:
		help()
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

	if config.Entries["logFileEnabled"].Value.(bool) {
		logLevel := syslog.Priority(config.Entries["logFileLevel"].Value.(int))
		file := config.Entries["logFileOutput"].Value.(string)
		_, err := logger.NewCh(log.ChConfig{Severity: &logLevel, Type: log.ChFile, File: file})
		if err != nil {
			fmt.Println(err)
		}
	}
	if config.Entries["logStdoutEnabled"].Value.(bool) {
		logLevel := syslog.Priority(config.Entries["logStdoutLevel"].Value.(int))
		_, err := logger.NewCh(log.ChConfig{Severity: &logLevel, Type: log.ChFile, File: os.Stdout})
		if err != nil {
			fmt.Println(err)
		}
	}
	if config.Entries["logStderrEnabled"].Value.(bool) {
		logLevel := syslog.Priority(config.Entries["logStderrLevel"].Value.(int))
		_, err := logger.NewCh(log.ChConfig{Severity: &logLevel, Type: log.ChFile, File: os.Stderr})
		if err != nil {
			fmt.Println(err)
		}
	}
	if config.Entries["logSyslogEnabled"].Value.(bool) {
		logLevel := syslog.Priority(config.Entries["logSyslogLevel"].Value.(int))
		_, err := logger.NewCh(log.ChConfig{Severity: &logLevel, Type: log.ChSyslog})
		if err != nil {
			fmt.Println(err)
		}
	}

	// Logger.Out(LOG_DEBUG, spew.Sdump(Config))

	// endregion: logger
	// region: init modules

	rasa.Config = config
	rasa.Logger = logger

	// endregion: init modules
	// region: routing

	// logger.Out(LOG_DEBUG, spew.Sdump(fs.FlagSet))

	var stdout []byte
	var stderr error

	switch subCmd {
	case "rasaexec":
		stdout, stderr = rasa.Exec(subArgs, nil)
	// case "rasainit":
	// 	stdout, stderr = rasa.Init()
	case "instancelist":
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
	case "rasaversion":
		stdout, stderr = rasa.Exec([]string{"--version"}, nil)
	default:
		msg := "invalid subcommand " + subCmd
		logger.Out(log.LOG_EMERG, msg)
		os.Exit(2)
	}

	// endregion: routing
	// region: outputs

	if stderr != nil {
		logger.Out(log.LOG_ERR, stderr)
	}
	if config.Entries["consoleEnabled"].Value.(bool) {
		fmt.Println(string(stdout))
	}

	// endregion: outputs

}
