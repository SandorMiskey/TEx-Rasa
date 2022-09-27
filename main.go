// region: packages

package main

import (
	"encoding/json"
	"fmt"
	"log/syslog"
	"os"
	"regexp"
	"strings"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"
	"github.com/davecgh/go-spew/spew"

	"github.com/SandorMiskey/TEx-Rasa/instance"
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
		"consoleStdout":    {Desc: "writes output to stdout", Type: "bool", Def: true},
		"consoleStderr":    {Desc: "writes errors to stderr", Type: "bool", Def: true},
		"logFileEnabled":   {Desc: "writes logs to file", Type: "bool", Def: false},
		"logFileLevel":     {Desc: "min severity for logs written into file", Type: "int", Def: 7},
		"logFileOutput":    {Desc: "logFileEnabled destination", Type: "string", Def: "./log"},
		"logStdoutEnabled": {Desc: "writes logs to stdout", Type: "bool", Def: false},
		"logStdoutLevel":   {Desc: "min severity for logs written to stdout", Type: "int", Def: 7},
		"logStderrEnabled": {Desc: "writes logs to stderr", Type: "bool", Def: true},
		"logStderrLevel":   {Desc: "min severity for logs written to stderr", Type: "int", Def: 7},
		"logSyslogEnabled": {Desc: "writes logs to local syslog", Type: "bool", Def: true},
		"logSyslogLevel":   {Desc: "min severity for logs written to syslog", Type: "int", Def: 7},
	}

	commonEntriesApply := []string{}
	commonEntries := map[string]cfg.Entry{
		"instanceEnabled": {Desc: "enable or disable instance", Type: "bool", Def: true},
		"instanceLock":    {Desc: "lock instances", Type: "bool", Def: true},
		"instanceName":    {Desc: "instance name to be registered or initiated, can be omitted like: cmd {instanceRegister|rasaInit} -instanceRoot /foo/bar instance-name-to-be-used)", Type: "string", Def: ""},
		"instanceRoot":    {Desc: "directory where instances are stored", Type: "string", Def: "/app/instances"},
		"rasaCmd":         {Desc: "rasa command", Type: "string", Def: "rasa"},
	}

	subCmd := strings.ToLower(os.Args[1])
	switch subCmd {
	case "instancelist":
		commonEntriesApply = []string{"instanceRoot"}
	case "instanceregister":
		commonEntriesApply = []string{
			"instanceEnabled",
			"instanceLock",
			"instanceName",
			"instanceRoot",
		}
		fs.Entries["instanceNLU"] = cfg.Entry{Desc: "enable or disable nlu (only) mode", Type: "bool", Def: false}
		fs.Entries["instancePort"] = cfg.Entry{Desc: "listening port, if there is a enabled instance with this port, then -instanceEnabled will be forced to false", Type: "int", Def: 5005}
	case "instanceroot":
		commonEntriesApply = []string{"instanceRoot"}
	case "rasaexec":
		commonEntriesApply = []string{
			"instanceLock",
			"instanceRoot",
			"rasaCmd",
		}
		fs.Entries["subArgs"] = cfg.Entry{Desc: "appended to the tail, use when you want to pass something begins with - (or use the -- separator)", Type: "string", Def: ""}
	case "rasainit":
		// commonEntriesApply = []string{
		// 	"instanceEnabled",
		// 	"instanceLock",
		// 	"instanceName",
		// 	"instanceRoot",
		// 	"rasaCmd",
		// }
		// fs.Entries["rasaPrompt"] = cfg.Entry{Desc: "choose default options for prompts and suppress warnings (DEPRECATED!)", Type: "bool", Def: false}
	case "rasaversion":
		commonEntriesApply = []string{"rasaCmd"}
	default:
		help()
	}

	for _, v := range commonEntriesApply {
		fs.Entries[v] = commonEntries[v]
	}

	errParse := fs.ParseCopy()
	if errParse != nil {
		panic(errParse)
	}

	// endregion: flag set
	// region: tail

	subArgs := fs.FlagSet.Args()
	if add, ok := config.Entries["subArgs"].Value.(string); ok && len(add) > 0 {
		subArgs = append(subArgs, add)
	}
	fmt.Println(len(subArgs))

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

	rasa.Logger = logger
	instance.Logger = logger

	// Logger.Out(LOG_DEBUG, spew.Sdump(Config))

	// endregion: logger
	// region: routing

	// logger.Out(LOG_DEBUG, spew.Sdump(fs.FlagSet))

	var out interface{}
	var err error

	switch subCmd {
	case "instancelist":
		out, err = instance.List(config)
	case "instanceregister":
		out, err = instance.Register(config, subArgs)
	case "instanceroot":
		out, err = instance.Root(config)
	case "rasaexec":
		out, err = rasa.Exec(config, subArgs, nil)
	case "rasainit":
		// 	out, stderr = rasa.Init()
	case "rasaversion":
		config.Entries["instanceLock"] = cfg.Entry{Value: false}
		out, err = rasa.Exec(config, []string{"--version"}, nil)
	default:
		msg := "invalid subcommand " + subCmd
		logger.Out(log.LOG_EMERG, msg)
		os.Exit(2)
	}

	// endregion: routing
	// region: outputs

	var stdout string
	if config.Entries["consoleStdout"].Value.(bool) {
		switch out := out.(type) {
		case []byte:
			stdout = string(out)
		case instance.Instances, *instance.Instances:
			b, e := json.Marshal(out)
			if e != nil {
				logger.Out(log.LOG_ERR, e)
				err = fmt.Errorf("%s\n%s)", err, e)
			}
			stdout = string(b)
		case nil:

		default:
			stdout = spew.Sdump(out)
		}
	}

	if err != nil {
		logger.Out(log.LOG_ERR, "stderr", err)
		if config.Entries["consoleStderr"].Value.(bool) {
			println(err.Error())
		}
	}
	if out != nil {
		logger.Out(log.LOG_DEBUG, "stdout", stdout)
		if config.Entries["consoleStderr"].Value.(bool) {
			fmt.Println(stdout)
		}
	}

	// endregion: outputs

}
