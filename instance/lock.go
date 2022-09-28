package instance

import (
	"errors"
	"os"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"
)

func Lock(c cfg.Config) (err error) {

	// skip lock

	lock, ok := c.Entries["instanceLock"].Value.(bool)
	if !ok {
		Logger.Out(log.LOG_WARNING, "cannot validate instanceLock")
		lock = true
	}
	if !lock {
		Logger.Out(log.LOG_NOTICE, "locking skipped due to instanceLock != true")
		return
	}

	// is instanceRoot valid and locked

	var root = new()

	if root, err = Root(c); err != nil {
		Logger.Out(log.LOG_ERR, err)
		return
	}

	if root.Locked {
		err = errors.New("already locked")
		Logger.Out(log.LOG_ERR, err)
		return
	}

	// lock

	err = os.WriteFile(root.Root+"/locked", nil, 0644)
	if err != nil {
		Logger.Out(log.LOG_ERR, err)
		return
	}

	// out

	Logger.Out(log.LOG_DEBUG, "instanceRoot locked")
	return
}

func Unlock(c cfg.Config) (err error) {

	// skip lock

	lock, ok := c.Entries["instanceLock"].Value.(bool)
	if !ok {
		Logger.Out(log.LOG_WARNING, "cannot validate instanceLock, locking set")
		lock = true
	}
	if !lock {
		Logger.Out(log.LOG_NOTICE, "locking skipped due to instanceLock != true")
		return
	}

	// is instanceRoot valid and locked

	var root = new()

	if root, err = Root(c); err != nil {
		Logger.Out(log.LOG_ERR, err)
		return
	}

	if !root.Locked {
		Logger.Out(log.LOG_WARNING, "instanceRoot is not locked actually")
		return nil
	}

	// remove lock

	if err = os.Remove(root.Root + "/locked"); err != nil {
		Logger.Out(log.LOG_ERR, err)
		return
	}

	// out

	Logger.Out(log.LOG_DEBUG, "instanceRoot lock removed")
	return
}
