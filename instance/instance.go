// region: packages

package instance

import "regexp"

// endregion: packages
// region: types and constants

type Config struct {
	Enabled bool   `json:"enabled"`
	NLU     bool   `json:"nlu"`
	Port    int    `json:"port"`
	Path    string `json:"path"`
}

// endregion: types and constants

func ValidName(name string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(name)
}
