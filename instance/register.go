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

	// region: instance name

	err = NameParse(&c, &i)
	if err != nil {
		Logger.Out(log.LOG_ERR, "unable to parse instance name", err)
		return
	}
	name := c.Entries["instanceName"].Value.(string)
	Logger.Out(log.LOG_DEBUG, me, "instance name", name)

	// endregion: instance name
	// region: lock

	if err = Lock(c); err != nil {
		Logger.Out(log.LOG_ERR, me, err)
		return
	}

	// endregion: lock
	// region: existing instances

	exists, err := Exists(c)
	if err != nil {
		Unlock(c)
		Logger.Out(log.LOG_ERR, me, err)
		return
	}
	if exists {
		Unlock(c)
		err = errors.New("instance " + name + " already exists")
		Logger.Out(log.LOG_ERR, me, err)
		return
	}
	Logger.Out(log.LOG_DEBUG, me, "instance doesn't exist")

	// endregion: existing instances
	// region: engine config

	var confInstance Instance
	var confByteArray []byte
	var confOk bool

	confInstance.Enabled, confOk = c.Entries["instanceEnabled"].Value.(bool)
	if !confOk {
		err = errors.New("invalid instanceEnabled")
		Logger.Out(log.LOG_ERR, me, err, c.Entries["instanceEnabled"].Value)
		return
	}

	confInstance.NLU, confOk = c.Entries["instanceNLU"].Value.(bool)
	if !confOk {
		err = errors.New("invalid instanceNLU")
		Logger.Out(log.LOG_ERR, me, err, c.Entries["instanceNLU"].Value)
		return
	}

	confInstance.Port, confOk = c.Entries["instancePort"].Value.(int)
	if !confOk {
		err = errors.New("invalid instancePort")
		Logger.Out(log.LOG_ERR, me, err, c.Entries["instancePort"].Value)
		return
	}

	instances, err = List(c)
	if err != nil {
		Unlock(c)
		Logger.Out(log.LOG_ERR, me, err)
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

	Logger.Out(log.LOG_DEBUG, me, "instance config", confInstance)

	// endregion: engine config
	// region: create

	if err = os.Mkdir(instances.Root+"/"+name, os.ModePerm); err != nil {
		Unlock(c)
		Logger.Out(log.LOG_ERR, me, err)
		return
	}

	Logger.Out(log.LOG_DEBUG, me, "instance directory created")

	err = os.WriteFile(instances.Root+"/"+name+"/engine.json", confByteArray, 0644)
	if err != nil {
		Unlock(c)
		Logger.Out(log.LOG_ERR, me, err)
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
