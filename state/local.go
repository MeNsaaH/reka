package state

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

// LocalBackend is an implementation of State that performs all operations
// locally. This is the "default" state type
type LocalBackend struct {
	state *State
	Path  string
}

// GetState returns state from local source
func (s LocalBackend) GetState() *State {
	if s.state.Empty() {
		if _, err := os.Stat(s.Path); os.IsNotExist(err) {
			log.Debugf("State file not found, using empty state")
			return s.state
		}
		stateFile, err := os.Open(s.Path)
		if err != nil {
			log.Fatalf("Could not open state file %s", s.Path)
		}
		defer stateFile.Close()
		byteValue, err := ioutil.ReadAll(stateFile)
		if err != nil {
			log.Fatalf("Failed to read data from state file %s", s.Path)
		}
		json.Unmarshal(byteValue, &s.state)
	}
	return s.state
}

// WriteState writes state to local path
func (s LocalBackend) WriteState(st *State) error {
	log.Debugf("Writing state to %s\n", s.Path)
	data, err := json.MarshalIndent(st, "", " ")
	if err != nil {
		log.Fatal("Failed to Load State for Writing")
	}
	err = ioutil.WriteFile(s.Path, data, 0644)
	if err != nil {
		log.Fatal("Failed to write state to file")
	}
	return nil
}
