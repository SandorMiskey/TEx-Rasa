// region: packages

package rasa

import (
	"encoding/json"
	"errors"
	"reflect"

	"go.uber.org/multierr"
)

// endregion: packages
// region: types and constants

type rasaTemplate struct {
	Command  string `json:"command"`
	LogLevel string `json:"logLevel"`
	Version  string `json:"version"`
}
type Rasa struct {
	Err error `json:"err"`
	rasaTemplate
}

// endregion: types and constants
// region: init

func New() Rasa {
	return (Rasa{
		Err: nil,
		rasaTemplate: rasaTemplate{
			Command:  "",
			LogLevel: "",
			Version:  "",
		},
	})
}

func Init(r *Rasa) (Rasa, error) {
	if len(r.Command) == 0 {
		r.Err = multierr.Append(r.Err, errors.New("empty command"))
	}
	if r.LogLevel != "" && r.LogLevel != "--debug" && r.LogLevel != "--verbose" && r.LogLevel != "--quiet" {
		r.Err = multierr.Append(r.Err, errors.New("invalid log level"))
		r.LogLevel = "--debug"
	}

	// TODO: Exec version
	r.Version = "TODO!"

	return *r, nil
}

func (r *Rasa) Init() (Rasa, error) {
	return Init(r)
}

// endregion: init
// region: json marshal

func (r Rasa) MarshalJSON() ([]byte, error) {
	var errorMsg string
	if r.Err != nil {
		errorMsg = r.Err.Error()
	}

	type anonStruct struct {
		Err string `json:"err"`
		rasaTemplate
	}
	anon := anonStruct{
		Err: errorMsg,
	}

	values := reflect.ValueOf(r.rasaTemplate)
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
