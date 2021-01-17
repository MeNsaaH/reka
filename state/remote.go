package state

import (
	"context"
	"encoding/json"
	"fmt"
	"gocloud.dev/blob"
	"io/ioutil"
	// Import the blob packages we want to be able to open.
	_ "gocloud.dev/blob/azureblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"

	log "github.com/sirupsen/logrus"
)

// RemoteBackend is an implementation of State that performs all operations
// in a cloud storage container e.g s3
type RemoteBackend struct {
	state    *State
	Path     string
	Bucket   string
	BlobType string
}

// GetState returns state from remote source
func (s RemoteBackend) GetState() *State {
	bucket, ctx := s.getBucket()
	defer bucket.Close()
	if s.state.Empty() {
		if exist, _ := bucket.Exists(ctx, s.Path); !exist {
			log.Debug("State file not found, using empty state")
			return s.state
		}
		stateFile, err := bucket.NewReader(ctx, s.Path, nil)
		if err != nil {
			log.Fatalf("Could not open state bucket: %s", err)
		}
		defer stateFile.Close()
		byteValue, err := ioutil.ReadAll(stateFile)
		if err != nil {
			log.Fatalf("Failed to read data from remote state file %s", s.Path)
		}
		json.Unmarshal(byteValue, &s.state)
	}
	return s.state
}

// WriteState writes state to remote path
func (s RemoteBackend) WriteState(st *State) error {
	log.Debugf("Writing state to remote %s\n", s.Path)
	data, err := json.MarshalIndent(st, "", " ")
	if err != nil {
		log.Fatal("Failed to Load State for Writing")
	}
	bucket, ctx := s.getBucket()
	defer bucket.Close()
	stateFile, err := bucket.NewWriter(ctx, s.Path, nil)
	if err != nil {
		log.Fatalf("Could not open state bucket: %s", err)
	}
	defer stateFile.Close()
	_, writeErr := stateFile.Write(data)
	if writeErr != nil {
		log.Fatal("Failed to write state to file")
	}
	return nil
}

func (s RemoteBackend) getBucket() (*blob.Bucket, context.Context) {
	ctx := context.Background()
	// Open a connection to the bucket.
	bucketURL := fmt.Sprintf("%s://%s",
		s.BlobType, s.Bucket)
	bucket, err := blob.OpenBucket(ctx, bucketURL)
	if err != nil {
		log.Fatalf("Failed to setup bucket: %s", err)
	}
	return bucket, ctx
}
