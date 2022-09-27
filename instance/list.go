package instance

import (
	"encoding/json"
	"os"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"
)

func List(c cfg.Config) (result *Instances, err error) {

	// region: validate instanceRoot

	result, err = Root(c)
	if err != nil {
		Logger.Out(log.LOG_ERR, err)
		return
	}

	// endregion: validate root
	// region: open and read root

	dir, err := os.Open(result.Root)
	if err != nil {
		Logger.Out(log.LOG_ERR, err)
		return
	}

	files, err := dir.Readdir(0)
	if err != nil {
		Logger.Out(log.LOG_ERR, err)
		return
	}

	// endregion: open and read root
	// region: fetch instances

	for _, v := range files {
		if v.IsDir() {
			if !validateName(v.Name()) {
				Logger.Out(log.LOG_NOTICE, "instance.List()", v.Name(), "invalid instance name")
				continue
			}
			Logger.Out(log.LOG_DEBUG, "instance.List() instance", v.Name())

			jsonFile, err := os.ReadFile(result.Root + "/" + v.Name() + "/engine.json")
			if err != nil {
				Logger.Out(log.LOG_ERR, "instance.List() unable to open engine.json", err)
				continue
			}
			Logger.Out(log.LOG_DEBUG, "instance.List() got config", string(jsonFile))

			var inst Instance
			json.Unmarshal(jsonFile, &inst)
			result.Instances[v.Name()] = &inst
			Logger.Out(log.LOG_DEBUG, "instance.List() parsed config", inst)
		}
	}

	// endregion: fetch instances

	return
}
