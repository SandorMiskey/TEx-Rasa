// region: packages

package instance

import (
	"errors"
	"regexp"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"
)

// endregion: packages
// region: types and constants

type Instance struct {
	Enabled bool `json:"enabled"`
	NLU     bool `json:"nlu"`
	Port    int  `json:"port"`
}

type Instances struct {
	Instances map[string]*Instance
	Locked    bool
	Root      string
}

var (
	Logger log.Logger
)

// endregion: types and constants

func new() *Instances {
	return &Instances{
		Locked:    false,
		Root:      "",
		Instances: make(map[string]*Instance),
	}

}

func Exists(c cfg.Config) (exists bool, err error) {

	name := c.Entries["instanceName"].Value.(string)
	if !NameValid(name) {
		err = errors.New("invalid instance name")
		Logger.Out(log.LOG_ERR, err)
		return
	}

	instances, err := List(c)
	if err != nil {
		Logger.Out(log.LOG_ERR, err)
		return
	}

	exists = false
	if instances.Instances[name] != nil {
		exists = true
		Logger.Out(log.LOG_DEBUG, "instance exists")
		return
	}

	Logger.Out(log.LOG_DEBUG, "instance doesn't exist")
	return
}

func NameValid(name string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(name)
}

func NameParse(c *cfg.Config, i *[]string) (err error) {
	name, ok := c.Entries["instanceName"].Value.(string)
	if !ok {
		err = errors.New("invalid instance name")
		Logger.Out(log.LOG_ERR, name, err)
		return
	}

	if len(name) == 0 {
		if len(*i) > 0 {
			name = (*i)[0]
			Logger.Out(log.LOG_DEBUG, "instance name grabbed from the tail", name)

			*i = (*i)[1:]
			Logger.Out(log.LOG_NOTICE, "extra parameters will be shifted", *i)
		}
	} else {
		Logger.Out(log.LOG_NOTICE, "extra parameters will be ignored", i)
	}
	c.Entries["instanceName"] = cfg.Entry{Value: name}

	if len(name) == 0 {
		err = errors.New("you have to name your instance to be registered or initialized")
		Logger.Out(log.LOG_ERR, err)
		return
	}

	if !NameValid(name) {
		err = errors.New("invalid instance name")
		Logger.Out(log.LOG_ERR, err)
		return
	}

	Logger.Out(log.LOG_DEBUG, "instance name", name)
	return
}
