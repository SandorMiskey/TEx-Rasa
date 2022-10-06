// region: packages

package engine

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"

	"github.com/SandorMiskey/TEx-Rasa/instance"
	"github.com/SandorMiskey/TEx-Rasa/rasa"
	"github.com/SandorMiskey/TEx-kit/log"
	"go.uber.org/multierr"
)

// endregion: packages
// region: types

type buildTemplate struct {
	Instances map[string]*instance.Config `json:"instances"`
	Locked    bool                        `json:"locked"`
	Logger    *log.Logger                 `json:"-"`
	Path      string                      `json:"path"`
	Rasa      rasa.Rasa                   `json:"rasa"`
}
type Build struct {
	Err error `json:"err"`
	buildTemplate
}

// endregion: types
// region: constants

const (
	config    = "/engine.json"
	instances = "/instances"
	lockfile  = "/locked"

	LOG_DEBUG = log.LOG_DEBUG
	LOG_ERR   = log.LOG_ERR
)

// endregion: constants
// region: init engine

func New() Build {
	return Build{
		Err: nil,
		buildTemplate: buildTemplate{
			Instances: map[string]*instance.Config{},
			Locked:    true,
			Logger:    log.NewLogger(),
			Path:      "",
			Rasa:      rasa.Rasa{},
		},
	}
}

func Init(b *Build) (Build, error) {

	// region: check logger

	if b.Logger == nil {
		b.Logger = log.NewLogger()
		defer b.Logger.Close()
	}

	// endregion: logger
	// region: check b.Path

	b.Path = strings.TrimSuffix(b.Path, "/")
	if _, err := os.Stat(b.Path); os.IsNotExist(err) {
		b.Err = multierr.Append(b.Err, err)
		b.Logger.Out(LOG_ERR, b.Path, err)
		return *b, err
	}
	b.Logger.Out(LOG_DEBUG, b.Path)

	// endregion: path
	// region: is locked

	b.Locked = true
	if _, err := os.Stat(b.Path + lockfile); os.IsNotExist(err) {
		b.Locked = false
	}
	b.Logger.Out(LOG_DEBUG, b.Locked)

	// endregion: lock
	// region: open and read root

	dir, err := os.Open(b.Path + instances)
	if err != nil {
		b.Err = multierr.Append(b.Err, err)
		b.Logger.Out(LOG_ERR, err)
		return *b, err
	}

	files, err := dir.Readdir(0)
	if err != nil {
		b.Logger.Out(log.LOG_ERR, err)
		return *b, err
	}

	// endregion: open and read root
	// region: fetch instances

	if b.Instances == nil {
		b.Instances = map[string]*instance.Config{}
	}

	for k, v := range files {
		if v.IsDir() {
			if !instance.ValidName(v.Name()) {
				b.Logger.Out(LOG_DEBUG, v.Name(), false)
				continue
			}
			b.Logger.Out(LOG_DEBUG, k, v.Name())

			path := b.Path + instances + "/" + v.Name()
			jsonPath := path + config
			jsonContent, err := os.ReadFile(jsonPath)
			if err != nil {
				b.Logger.Out(LOG_ERR, err)
				continue
			}
			b.Logger.Out(log.LOG_DEBUG, jsonPath, string(jsonContent))

			var inst instance.Config
			json.Unmarshal(jsonContent, &inst)
			inst.Path = path
			b.Instances[v.Name()] = &inst
			b.Logger.Out(log.LOG_DEBUG, inst)
		}
	}

	// endregion: fetch instances
	// region: set rasa

	// TODO: init if b.Rasa == nil
	// rasa := rasa.New()
	// rasa.Init()

	// endregion: rasa

	return *b, nil
}

func (b *Build) Init() (Build, error) {
	return Init(b)
}

// endregion: init engine
// region: json marshal

func (b Build) MarshalJSON() ([]byte, error) {
	var errorMsg string
	if b.Err != nil {
		errorMsg = b.Err.Error()
	}

	type anonStruct struct {
		Err string `json:"err"`
		buildTemplate
	}
	anon := anonStruct{
		Err: errorMsg,
	}

	values := reflect.ValueOf(b.buildTemplate)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		if types.Field(i).Name == "Err" {
			continue
		}
		reflect.ValueOf(&anon).Elem().FieldByName(types.Field(i).Name).Set(values.Field(i))
	}

	return json.Marshal(anon)
}

// endregion: json marshal
