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
		- params and subArgs
		- check if instance exists (instance.List())
		- instance.Register() in doesn't exist (flag if it's allowed or not) (not sure if this is a good idea...)
		- set --init-dir (instanceRoot + instanceName)
		- init (rasaExec())
	*/

	result, err = Exec(c, []string{"init", "--no-prompt", logLevel(c), "-h"}, nil)
	if err != nil {
		Logger.Out(log.LOG_ERR, "rasa.Init()", err)
	}
	return

	// endregion: execute

}
