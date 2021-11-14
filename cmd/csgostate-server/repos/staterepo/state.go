package staterepo

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/roessland/csgostate/csgostate"
	bolt "go.etcd.io/bbolt"
	"io/ioutil"
	"os"
)

type StateRepo interface {
	GetLatest(steamID string) (*csgostate.State, error)
	Push(state *csgostate.State) error
	DebugJsonForPlayer(steamID string) error
	GetAllForPlayer(steamID string) ([]csgostate.State, error)
}

var _ StateRepo = &DBStateRepo{}

type DBStateRepo struct {
	db *bolt.DB
}

func getStatesBucketNameForUser(steamID string) []byte {
	return []byte(fmt.Sprintf("states-%s", steamID))
}

func NewDBStateRepo(db *bolt.DB) (*DBStateRepo, error) {
	return &DBStateRepo{
		db: db,
	}, nil
}

func (stateRepo *DBStateRepo) GetLatest(steamID string) (*csgostate.State, error) {
	var state *csgostate.State
	err := stateRepo.db.View(func(tx *bolt.Tx) error {
		bucketName := getStatesBucketNameForUser(steamID)
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return nil
		}

		cursor := bucket.Cursor()
		key, val := cursor.Last()
		if key == nil {
			return nil
		}

		s := &csgostate.State{}
		err := json.Unmarshal(val, &state)
		if err != nil {
			return errors.Wrapf(err, "GetLatest: error when decoding %s", string(val))
		}
		state = s
		state.RawJson = make([]byte, len(val))
		copy(state.RawJson, val)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return state, nil
}

// Push adds a state event to a user's bucket.
// Key is an autoincrement.
// Value is the raw json.
func (stateRepo *DBStateRepo) Push(state *csgostate.State) error {
	// Store the user in the user bucket using the SteamID as the key.
	err := stateRepo.db.Update(func(tx *bolt.Tx) error {
		bucketName := getStatesBucketNameForUser(state.Provider.SteamID)
		bucket, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		id, err := bucket.NextSequence()
		if err != nil {
			return err
		}
		fmt.Println("pushing rawjson of length ", len(state.RawJson), "bucket ", string(bucketName))
		return bucket.Put(itob(id), state.RawJson)
	})
	return err
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func (stateRepo *DBStateRepo) DebugJsonForPlayer(steamID string) error {
	return stateRepo.db.View(func(tx *bolt.Tx) error {
		bucketName := getStatesBucketNameForUser(steamID)
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return nil
		}

		cursor := bucket.Cursor()
		for key, val := cursor.First(); key != nil; key, val = cursor.Next() {
			s := &csgostate.State{}
			err := json.Unmarshal(val, &s)
			if err != nil {
				ioutil.WriteFile("whendecoding.txt", []byte(val), 0666)
				return errors.Wrapf(err, "when decoding %s", string(val))
			}
			a := prettyPrintRawJson(val)
			b := prettyPrintState(s)
			if a != b {
				ioutil.WriteFile("rawjson.txt", []byte(a), 0666)
				ioutil.WriteFile("state.txt", []byte(b), 0666)
				os.Exit(1337)
			}
		}
		return nil
	})
}

func prettyPrintRawJson(rawJson []byte) string {
	var v map[string]interface{}
	err := json.Unmarshal(rawJson, &v)
	if err != nil {
		panic(err)
	}
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(buf)
}

func prettyPrintState(state *csgostate.State) string {
	buf, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(buf)
}

func (stateRepo *DBStateRepo) GetAllForPlayer(steamID string) ([]csgostate.State, error) {
	var states []csgostate.State
	err := stateRepo.db.View(func(tx *bolt.Tx) error {
		bucketName := getStatesBucketNameForUser(steamID)
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return nil
		}

		cursor := bucket.Cursor()
		for key, val := cursor.First(); key != nil; key, val = cursor.Next() {
			s := csgostate.State{}
			err := json.Unmarshal(val, &s)
			if err != nil {
				return errors.Wrapf(err, "when decoding %s", string(val))
			}
			s.RawJson = make([]byte, len(val))
			copy(s.RawJson, val)
			states = append(states, s)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return states, nil
}
