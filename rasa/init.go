// region: packages

package rasa

import (
	"errors"

	"github.com/SandorMiskey/TEx-Rasa/instance"
	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"

	"go.uber.org/multierr"
)

// endregion: packages
// region: messages

var (
	ErrInvalidRasaPrompt = errors.New("invalid rasaPrompt")
)

// endregion: messages

func Init(c cfg.Config, i []string) (result []byte, err error) {

	// region: instance name and root

	err = instance.NameParse(&c, &i)
	if err != nil {
		Logger.Out(log.LOG_ERR, "unable to parse instance name", err)
		return
	}
	name := c.Entries["instanceName"].Value.(string)
	Logger.Out(log.LOG_DEBUG, "instance name", name)

	_, err = instance.Root(c)
	if err != nil {
		Logger.Out(log.LOG_ERR, "failed to validate instanceRoot", err)
		return
	}
	root := c.Entries["instanceRoot"].Value.(string)
	Logger.Out(log.LOG_DEBUG, "instance root", root)

	// endregion: instance name
	// region: lock

	if err = instance.Lock(c); err != nil {
		Logger.Out(log.LOG_ERR, err)
		return
	}

	// endregion: lock
	// region: exists

	exists, err := instance.Exists(c)
	if err != nil {
		instance.Unlock(c)
		Logger.Out(log.LOG_ERR, err)
		return
	}
	if !exists {
		instance.Unlock(c)
		err = errors.New("instance " + name + " doesn't exist")
		Logger.Out(log.LOG_ERR, err)
		return
	}
	Logger.Out(log.LOG_DEBUG, "instance exists")

	// endregion: existing instances
	// region: execute

	subCmd := []string{"init", logLevel(c), "--no-prompt", "--init-dir", root + "/" + name}
	subCmd = append(subCmd, i...)
	Logger.Out(log.LOG_DEBUG, "rasa.Init() subCmd", subCmd)

	lock := c.Entries["instanceLock"]
	c.Entries["instanceLock"] = cfg.Entry{Value: false}

	result, err = Exec(c, subCmd, nil)

	c.Entries["instanceLock"] = lock
	err = multierr.Append(err, instance.Unlock(c))
	if err != nil {
		Logger.Out(log.LOG_ERR, "rasa.Init() exec and unlock", err)
	}

	return

	// endregion: execute

}
