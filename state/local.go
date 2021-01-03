package state

import (
	"encoding/json"
	"fmt"
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

// Get returns state from local source
func (s LocalBackend) Get() *State {
	if s.state.empty() {
		if _, err := os.Stat(s.Path); os.IsNotExist(err) {
			s.state = &State{}
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
	file, err := json.MarshalIndent(st, "", " ")
	if err != nil {
		log.Fatal("Failed to Load State for Writing")
	}
	fmt.Println(string(file))
	err = ioutil.WriteFile(s.Path, file, 0644)
	if err != nil {
		log.Fatal("Failed to write state to file")
	}
	return nil
}
