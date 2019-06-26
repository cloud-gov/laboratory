package libheader

import (
	"encoding/json"
	"io"
	"sync"
)

type Intent struct {
	lock sync.Mutex
	// Headers we have and are setting via Load()
	Have map[string][]string
}

func NewIntent() *Intent {
	return &Intent{}
}

// Reads a JSON file into the comparer.
func (i *Intent) Load(targetHeaders io.Reader) (bool, error) {
	if err := json.NewDecoder(targetHeaders).Decode(&i.Have); err != nil {
		return false, err
	}
	return true, nil
}
