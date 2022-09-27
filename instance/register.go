package instance

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"
)

func Register(c cfg.Config, i []string) (instances *Instances, err error) {

	me := "instance.Register()"

	// region: existing instances

	instances, err = List(c)
	if err != nil {
		Logger.Out(log.LOG_ERR, me, err)
		return
	}

	// endregion: existing instances
	// region: instance name

	name, ok := c.Entries["instanceName"].Value.(string)
	if !ok {
		err = errors.New("invalid instance name")
		Logger.Out(log.LOG_ERR, me, name, err)
		return
	}

	if len(name) == 0 {
		if len(i) > 0 {
			name = i[0]
			if len(i) > 1 {
				Logger.Out(log.LOG_DEBUG, me, "extra parameters will be ignored", i[1:])
			}
		}
	} else {
		Logger.Out(log.LOG_DEBUG, me, "extra parameters will be ignored", i)
	}

	if len(name) == 0 {
		err = errors.New("you have to name your instance to be registered")
		Logger.Out(log.LOG_ERR, me, err)
		return
	}

	if !validateName(name) {
		err = errors.New("invalid instance name")
		Logger.Out(log.LOG_ERR, me, err)
		return
	}

	if instances.Instances[name] != nil {
		err = errors.New("instance " + name + " already exists")
		Logger.Out(log.LOG_ERR, me, err)
		return
	}

	Logger.Out(log.LOG_DEBUG, me, "instance name", name)

	// endregion: instance name
	// region: engine config

	var confInstance Instance
	var confByteArray []byte

	confInstance.Enabled, ok = c.Entries["instanceEnabled"].Value.(bool)
	if !ok {
		err = errors.New("invalid instanceEnabled")
		Logger.Out(log.LOG_ERR, me, err, c.Entries["instanceEnabled"].Value)
		return
	}

	confInstance.NLU, ok = c.Entries["instanceNLU"].Value.(bool)
	if !ok {
		err = errors.New("invalid instanceNLU")
		Logger.Out(log.LOG_ERR, me, err, c.Entries["instanceNLU"].Value)
		return
	}

	confInstance.Port, ok = c.Entries["instancePort"].Value.(int)
	if !ok {
		err = errors.New("invalid instancePort")
		Logger.Out(log.LOG_ERR, me, err, c.Entries["instancePort"].Value)
		return
	}
	for _, v := range instances.Instances {
		if v.Port == confInstance.Port && v.Enabled {
			Logger.Out(log.LOG_WARNING, me, "conflicting instance found, new instance forced to be disabled")
			confInstance.Enabled = false
		}
	}

	confByteArray, err = json.Marshal(confInstance)
	if err != nil {
		Logger.Out(log.LOG_ERR, err)
		return
	}

	// endregion: engine config
	// region: lock

	if err = Lock(c); err != nil {
		Logger.Out(log.LOG_ERR, me, err)
		return
	}

	// endregion: lock
	// region: create

	if err = os.Mkdir(instances.Root+"/"+name, os.ModePerm); err != nil {
		Logger.Out(log.LOG_ERR, me, err)
		Unlock(c)
		return
	}

	Logger.Out(log.LOG_DEBUG, me, "instance directory created")

	err = os.WriteFile(instances.Root+"/"+name+"/engine.json", confByteArray, 0644)
	if err != nil {
		Logger.Out(log.LOG_ERR, me, err)
		Unlock(c)
		return
	}

	Logger.Out(log.LOG_DEBUG, me, "instance config created")

	// endregion: create

	Unlock(c)
	instances, err = List(c)
	if err != nil {
		Logger.Out(log.LOG_ERR, me, err)
		return
	}
	return

}
