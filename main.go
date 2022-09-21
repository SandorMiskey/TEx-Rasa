// region: packages

package main

import (
	"bufio"
	"fmt"
	"log/syslog"
	"os"
	"os/exec"

	// "github.com/SandorMiskey/TEx-Rasa/instance"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"
	"github.com/davecgh/go-spew/spew"
)

// endregion: packages
// region: global variables

var (
	Config cfg.Config = *cfg.NewConfig(os.Args[0])
	Logger log.Logger = *log.NewLogger()
)

const (
	LOG_ERR    syslog.Priority = log.LOG_ERR
	LOG_NOTICE syslog.Priority = log.LOG_NOTICE
	LOG_INFO   syslog.Priority = log.LOG_INFO
	LOG_DEBUG  syslog.Priority = log.LOG_DEBUG
	LOG_EMERG  syslog.Priority = log.LOG_EMERG
)

// endregion: globals

func main() {

	// region: cli flags

	if len(os.Args) < 2 {
		fmt.Println("subcommand expected")
		os.Exit(1)
	}
	subCommand := os.Args[1]

	flagSet := Config.NewFlagSet(Config.Name + " " + subCommand)
	flagSet.Arguments = os.Args[2:]
	flagSet.Entries = map[string]cfg.Entry{
		"instanceRoot": {Desc: "directory where the instances are stored", Type: "string", Def: "/app/instances"},
		"logLevel":     {Desc: "Logger min severity", Type: "int", Def: 7},
	}

	switch subCommand {
	case "copy":
	case "destroy":
	case "init":
	case "list":
	case "version":
	default:
		fmt.Println("no such subcommand '" + subCommand + "', usage: " + Config.Name + " {init,destroy,list} [options] [args]")
		os.Exit(1)
	}

	err := flagSet.ParseCopy()
	if err != nil {
		panic(err)
	}
	// subArgs := flagSet.FlagSet.Args()

	// endregion: cli
	// region: logger

	defer Logger.Close()

	logLevel := syslog.Priority(Config.Entries["logLevel"].Value.(int))
	_, _ = Logger.NewCh(log.ChConfig{Severity: &logLevel})

	// Logger.Out(LOG_DEBUG, spew.Sdump(Config))

	// endregion: logger
	// region: routing

	Logger.Out(LOG_DEBUG, spew.Sdump(flagSet.FlagSet))

	switch subCommand {
	case "copy":

	case "destroy":
		/*
			trash directory from Config
		*/

	case "init":
		/*
			-v info
			-vv debug
			-quiet warning

			--init_dir (root/instance)
			--no_prompt

			rasaCMD
			instanceRoot
			instanceEnabled
			subArg = list of instances
		*/

	case "list":
	// files, err := ioutil.ReadDir(Config.Entries["instanceRoot"].Value.(string))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, f := range files {
	// 	fmt.Println(f.Name())
	// }
	// instance.List()
	case "version":
		cmd := exec.Command("rasa", "--version")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			panic(err)
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			panic(err)
		}
		cmd.Start()
		stdin.Close()
		scannerErr := bufio.NewScanner(stderr)
		for scannerErr.Scan() {
			fmt.Printf("commonGwExecuteActual() stderr: %s\n", scannerErr.Text())
		}
		if err := scannerErr.Err(); err != nil {
			panic(err)
		}
		scannerOut := bufio.NewScanner(stdout)
		for scannerOut.Scan() {
			fmt.Println(string(scannerOut.Bytes()))
		}
		if err := scannerOut.Err(); err != nil {
			panic(err)
		}
		if err := cmd.Wait(); err != nil {
			panic(err)
		}

	default:
		panic("no such subcommand '" + subCommand)
	}

	// endregion:

}
