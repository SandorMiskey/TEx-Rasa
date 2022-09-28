package instance

import (
	"errors"
	"os"
	"strings"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"
)

func Root(c cfg.Config) (inst *Instances, err error) {

	inst = new()

	// region: validate instanceRoot type

	dir, ok := c.Entries["instanceRoot"].Value.(string)
	if !ok {
		err = errors.New("invalid instance root type")
		Logger.Out(log.LOG_ERR, err)
		return
	}
	inst.Root = strings.TrimSuffix(dir, "/")

	Logger.Out(log.LOG_DEBUG, "instanceRoot", inst.Root)

	// endregion: validate instanceRoot type
	// region: check if exists

	if _, err = os.Stat(inst.Root); os.IsNotExist(err) {
		Logger.Out(log.LOG_ERR, err)
		return
	}

	Logger.Out(log.LOG_DEBUG, "instanceRoot exists")

	// endregion: check if exists
	// region: is locked

	inst.Locked = true
	if _, e := os.Stat(inst.Root + "/locked"); os.IsNotExist(e) {
		inst.Locked = false
	}

	Logger.Out(log.LOG_DEBUG, "instanceRoot locked", inst.Locked)

	// endregion: is locked
	// out

	return

}
