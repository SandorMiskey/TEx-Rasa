// region: packages

package rasa

import (
	"errors"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"
)

// endregion: packages
// region: messages

var (
	ErrInvalidRasaPrompt = errors.New("invalid rasaPrompt")
)

// endregion: messages

func Init(c cfg.Config) (result []byte, err error) {

	/*
		config:
		- instanceEnabled
		- port
		- nlu only

		validate if exists
		validate instance name

		subArg = list of instances -> --init-dir
		create instance dir
		cd instance dir

		exec
	*/

	// region: validations

	// TODO

	// endregion: validations
	// region: no prompt

	// prompt := "--no-prompt"
	// if _, ok := Config.Entries["rasaPrompt"].Value.(bool); !ok {
	// 	err = ErrInvalidRasaPrompt
	// 	Logger.Out(log.LOG_ERR, err)
	// 	return
	// }
	// if Config.Entries["rasaPrompt"].Value.(bool) {
	// 	prompt = ""
	// }

	// endregion: no prompt
	// region: execute

	// if err = Wd(); err != nil {
	// 	Logger.Out(log.LOG_ERR, err)
	// 	return
	// }

	result, err = Exec(c, []string{"init", "--no-prompt", logLevel(c), "-h"}, nil)
	if err != nil {
		Logger.Out(log.LOG_ERR, "rasa.Init()", err)
	}
	return

	// endregion: execute

}
