// region: packages

package instance

import (
	"regexp"

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

func NewInstances() *Instances {
	return &Instances{
		Locked:    false,
		Root:      "",
		Instances: make(map[string]*Instance),
	}

}

func validateName(name string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(name)
}
