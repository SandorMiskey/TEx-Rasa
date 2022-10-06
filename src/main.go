// region: packages

package main

import (
	"encoding/json"
	"fmt"
	"log/syslog"
	"os"
	"regexp"
	"strings"

	"github.com/SandorMiskey/TEx-Rasa/engine"
	"github.com/SandorMiskey/TEx-Rasa/rasa"
	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/multierr"
	"github.com/SandorMiskey/TEx-Rasa/instance"
)

// endregion: packages
// region: types and globals

type Command struct {
	cmd  string
	conf cfg.Config
	fs   cfg.FlagSet
	sub  string
	tail []string
}

const (
	LOG_DEBUG   syslog.Priority = log.LOG_DEBUG
	LOG_EMERG   syslog.Priority = log.LOG_EMERG
	LOG_ERR     syslog.Priority = log.LOG_ERR
	LOG_INFO    syslog.Priority = log.LOG_INFO
	LOG_WARNING syslog.Priority = log.LOG_WARNING

	// LOG_NOTICE syslog.Priority = log.LOG_NOTICE
	// LOG_INFO   syslog.Priority = log.LOG_INFO
)

// endregion: globals

func main() {

	// region: cli flags

	// region: subcommand

	// TODO: list subcommands

	c := Command{}
	re := regexp.MustCompile(`.*/`)
	c.cmd = re.ReplaceAllString(os.Args[0], "")

	help := func() {
		fmt.Println("one of these subcommands expected:")
		fmt.Println("	- b|build: init and dump engine")
		fmt.Println("	- r|rebuild: rebuild and dump engine")
		fmt.Println("	- ...")
		fmt.Println("")
		fmt.Println("use `" + c.cmd + " subcommand --help` for further info")
		os.Exit(1)
	}

	if len(os.Args) < 2 || strings.HasPrefix(os.Args[1], "-") {
		help()
	}
	c.sub = strings.ToLower(os.Args[1])

	// endregion: subcommand
	// region: flag set

	c.conf = *cfg.NewConfig(c.cmd)
	c.fs = *c.conf.NewFlagSet(c.conf.Name + " " + c.sub)
	c.fs.Arguments = os.Args[2:]
	c.fs.Entries = map[string]cfg.Entry{
		"forceBuild":   {Desc: "ignore engine validation errors", Type: "bool", Def: false},
		"forceRasa":    {Desc: "ignore rasa validation errors", Type: "bool", Def: false},
		"indent":       {Desc: "json indentation", Type: "string", Def: "	"},
		"logFile":      {Desc: "logFileEnabled destination", Type: "string", Def: "./log"},
		"logLevel":     {Desc: "min severity for logs written into file", Type: "int", Def: 7},
		"logStderr":    {Desc: "writes logs to stderr", Type: "bool", Def: true},
		"logStdout":    {Desc: "writes logs to stdout", Type: "bool", Def: false},
		"logSyslog":    {Desc: "writes logs to local syslog", Type: "bool", Def: false},
		"rasaCmd":      {Desc: "rasa command", Type: "string", Def: "rasa"},
		"rasaLogLevel": {Desc: "rasa log level (--quiet|--verbose|--debug), calculated form -logLevel if -", Type: "string", Def: ""},
		"root":         {Desc: "directory where instances are stored", Type: "string", Def: "/app"},
		"stdout":       {Desc: "writes output to stdout", Type: "bool", Def: true},
		"stderr":       {Desc: "writes errors to stderr", Type: "bool", Def: true},
	}

	commonEntriesApply := []string{}
	commonEntries := map[string]cfg.Entry{
		// "instanceEnabled": {Desc: "enable or disable instance", Type: "bool", Def: true},
		// "instanceLock":    {Desc: "lock instances", Type: "bool", Def: true},
		// "instanceName":    {Desc: "instance name to be registered or initiated, can be omitted like: cmd {instanceRegister|rasaInit} -instanceRoot /foo/bar instance-name-to-be-used)", Type: "string", Def: ""},
		// "instancePort":    {Desc: "listening port, if there is a enabled instance with this port, then -instanceEnabled will be forced to false", Type: "int", Def: 5005},
		// "rasaCmd":         {Desc: "rasa command", Type: "string", Def: "rasa"},
		// "rasaLogLevel":    {Desc: "min severity for logs written to syslog", Type: "int", Def: 7},
		"subArgs": {Desc: "appended to the tail, use when you want to pass something begins with - (or use the -- separator)", Type: "string", Def: ""},
	}

	switch c.sub {
	case "b", "build":
		c.sub = "build"
	case "r", "rebuild":
		c.sub = "rebuild"
	// case "instancelist":
	// 	commonEntriesApply = []string{"instanceRoot"}
	// case "instanceregister":
	// commonEntriesApply = []string{
	// 	"instanceEnabled",
	// 	"instanceLock",
	// 	"instanceName",
	// 	"instancePort",
	// 	"instanceRoot",
	// }
	// fs.Entries["instanceNLU"] = cfg.Entry{Desc: "enable or disable nlu (only) mode", Type: "bool", Def: false}
	// case "instanceroot":
	// 	commonEntriesApply = []string{"instanceRoot"}
	// case "rasaexec":
	// 	commonEntriesApply = []string{
	// 		"instanceLock",
	// 		"instanceRoot",
	// 		"rasaCmd",
	// 		"rasaLogLevel",
	// 		"subArgs",
	// 	}
	// case "rasainit":
	// 	commonEntriesApply = []string{
	// 		"instanceLock",
	// 		"instanceName",
	// 		"instanceRoot",
	// 		"rasaCmd",
	// 		"rasaLogLevel",
	// 		"subArgs",
	// 	}
	// 	fs.Entries["rasaPrompt"] = cfg.Entry{Desc: "choose default options for prompts and suppress warnings (DEPRECATED!)", Type: "bool", Def: false}
	// case "rasaversion":
	// commonEntriesApply = []string{"rasaCmd"}
	default:
		help()
	}

	for _, v := range commonEntriesApply {
		c.fs.Entries[v] = commonEntries[v]
	}

	errParse := c.fs.ParseCopy()
	if errParse != nil {
		panic(errParse)
	}

	// endregion: flag set
	// region: tail

	c.tail = c.fs.FlagSet.Args()
	if add, ok := c.conf.Entries["subArgs"].Value.(string); ok && len(add) > 0 {
		c.tail = append(c.tail, add)
	}

	// endregion: tail

	// endregion: cli
	// region: logger

	logLevel := syslog.Priority(c.conf.Entries["logLevel"].Value.(int))
	logger := *log.NewLogger()
	defer logger.Close()

	if c.conf.Entries["logFile"].Value.(string) != "" {
		_, err := logger.NewCh(log.ChConfig{Severity: &logLevel, Type: log.ChFile, File: c.conf.Entries["logFile"].Value.(string)})
		if err != nil {
			println(err)
		}
	}
	if c.conf.Entries["logStdout"].Value.(bool) {
		_, err := logger.NewCh(log.ChConfig{Severity: &logLevel, Type: log.ChFile, File: os.Stdout})
		if err != nil {
			println(err)
		}
	}
	if c.conf.Entries["logStderr"].Value.(bool) {
		_, err := logger.NewCh(log.ChConfig{Severity: &logLevel, Type: log.ChFile, File: os.Stderr})
		if err != nil {
			println(err)
		}
	}
	if c.conf.Entries["logSyslog"].Value.(bool) {
		_, err := logger.NewCh(log.ChConfig{Severity: &logLevel, Type: log.ChSyslog})
		if err != nil {
			println(err)
		}
	}

	logger.Out(LOG_DEBUG, "command", c)

	// endregion: logger
	// region: init rasa

	rasa := rasa.New()
	rasa.Command = c.conf.Entries["rasaCmd"].Value.(string)
	rasa.LogLevel = c.conf.Entries["rasaLogLevel"].Value.(string)

	if rasa.LogLevel == "-" {
		switch level := syslog.Priority(c.conf.Entries["logLevel"].Value.(int)); {
		case level == LOG_DEBUG:
			rasa.LogLevel = "--debug"
		case level >= LOG_INFO:
			rasa.LogLevel = "--verbose"
		case level <= LOG_WARNING:
			rasa.LogLevel = "--quiet"
		default:
			rasa.LogLevel = ""
		}
	}
	rasa.Init()
	if rasa.Err != nil {
		logger.Out(log.LOG_ERR, "rasa init failed", rasa.Err)
		if !c.conf.Entries["forceRasa"].Value.(bool) {
			c.sub = "build"
		} else {
			logger.Out(LOG_WARNING, "rasa init failed, but execution forced")
		}
	}

	logger.Out(LOG_DEBUG, "rasa", rasa)

	// endregion: rasa
	// region: init build

	build := engine.New()
	build.Path = c.conf.Entries["root"].Value.(string)
	build.Logger = &logger
	build.Rasa = rasa
	build.Init()
	if build.Err != nil {
		logger.Out(log.LOG_ERR, "build failed", build.Err)
		if !c.conf.Entries["forceBuild"].Value.(bool) {
			c.sub = "build"
		} else {
			logger.Out(LOG_WARNING, "build failed, but execution forced")
		}
	}
	build.Err = multierr.Append(build.Err, rasa.Err)
	logger.Out(LOG_DEBUG, "build", build)

	// endregion: builds
	// region: routing

	var out interface{}
	var err error

	switch c.sub {
	case "build":
		out = build
		err = build.Err
	case "rebuild":
		rasa, _ = rasa.Init()
		out, _ = build.Init()
		err = multierr.Append(rasa.Err, build.Err)
	// case "instancelist":
	// 	// out, err = instance.List(config)
	// case "instanceregister":
	// 	// out, err = instance.Register(config, subArgs)
	// case "instanceroot":
	// 	// out, err = instance.Root(config)
	// case "rasaexec":
	// 	// out, err = rasa.Exec(config, subArgs, nil)
	// case "rasainit":
	// 	// out, err = rasa.Init(config, subArgs)
	// case "rasaversion":
	// c.conf.Entries["instanceLock"] = cfg.Entry{Value: false}
	// out, err = rasa.Exec(config, []string{"--version"}, nil)
	default:
		logger.Out(LOG_EMERG, "lost subcommand")
		help()
	}

	logger.Out(LOG_DEBUG, "err", err)
	logger.Out(LOG_DEBUG, "out", out)

	// endregion: routing
	// region: outputs

	var stdout string
	if c.conf.Entries["stdout"].Value.(bool) {
		switch out := out.(type) {

		case []byte:
			stdout = string(out)

		case engine.Build, *engine.Build:
			// b, e := json.Marshal(out)
			b, e := json.MarshalIndent(out, "", c.conf.Entries["indent"].Value.(string))
			if e != nil {
				logger.Out(LOG_ERR, e)
				err = multierr.Append(err, e)
			}
			stdout = string(b)

		default:
			stdout = spew.Sdump(out)

		}
	}

	if err != nil {
		logger.Out(LOG_ERR, "stderr", err)
		if c.conf.Entries["stderr"].Value.(bool) {
			println(fmt.Sprintf("%+v", err))
		}
	}
	if out != nil {
		logger.Out(LOG_DEBUG, "stdout", stdout)
		if c.conf.Entries["stdout"].Value.(bool) {
			fmt.Println(stdout)
		}
	}

	// endregion: outputs

}
