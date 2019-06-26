package libheader

import (
	"encoding/json"
	"io"
	"sync"
)

type HeaderComparer struct {
	lock sync.Mutex
	Have map[string][]string
	Want map[string][]string
}

func NewComparer() *HeaderComparer {
	return &HeaderComparer{}
}

// Compares two header targets and returns the diff.
func (h *HeaderComparer) Compare(have map[string][]string) []DiffStruct {
	// lock so we can compare at a given point in time.
	h.lock.Lock()
	h.Have = have
	if diff := Equal(h.Have, h.Want); diff != nil {
		h.lock.Unlock()
		return diff
	}
	h.lock.Unlock()
	return []DiffStruct{}
}

// Reads a JSON file into the comparer.
func (h *HeaderComparer) Load(targetHeaders io.Reader) (bool, error) {
	if err := json.NewDecoder(targetHeaders).Decode(&h.Want); err != nil {
		return false, err
	}
	return true, nil
}
