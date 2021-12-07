package staterepo

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/roessland/csgostate/pkg/csgostate"
	bolt "go.etcd.io/bbolt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type StateRepo interface {
	GetLatest() (*csgostate.State, error)
	GetLatestForPlayer(steamID string) (*csgostate.State, error)
	Push(state *csgostate.State) error
	DebugJsonForPlayer(steamID string) error
	GetAllForPlayer(steamID string) ([]csgostate.State, error)
	GetAll() ([]csgostate.State, error)
	ArchiveMonth(year, month int, w io.Writer) error
	DeleteUserDBStuff() error
}

var _ StateRepo = &DBStateRepo{}

type DBStateRepo struct {
	userDB  *bolt.DB
	stateDB *bolt.DB
}

var statesBucketName = []byte("states")

func getStatesBucketNameForUser(steamID string) []byte {
	return []byte(fmt.Sprintf("states-%s", steamID))
}

func NewDBStateRepo(userDB *bolt.DB, stateDB *bolt.DB) (*DBStateRepo, error) {
	return &DBStateRepo{
		userDB:  userDB,
		stateDB: stateDB,
	}, nil
}

// GetLatest gets the latest state in the DB.
func (stateRepo *DBStateRepo) GetLatest() (*csgostate.State, error) {
	var state *csgostate.State
	err := stateRepo.stateDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(statesBucketName)
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

// GetLatestForPlayer gets the latest state in the DB for a certain player.
func (stateRepo *DBStateRepo) GetLatestForPlayer(steamID string) (*csgostate.State, error) {
	var state *csgostate.State
	err := stateRepo.stateDB.View(func(tx *bolt.Tx) error {
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
			return errors.Wrapf(err, "GetLatestForPlayer: error when decoding %s", string(val))
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

// Push adds a state event to a user's bucket and to the global bucket.
// Key is an autoincrement.
// Value is the raw json.
func (stateRepo *DBStateRepo) Push(state *csgostate.State) error {
	// Store the user in the user bucket using the SteamID as the key.
	err := stateRepo.stateDB.Update(func(tx *bolt.Tx) error {
		// states
		{
			bucket, err := tx.CreateBucketIfNotExists(statesBucketName)
			if err != nil {
				return err
			}
			id, err := bucket.NextSequence()
			if err != nil {
				return err
			}
			err = bucket.Put(itob(id), state.RawJson)
			if err != nil {
				return err
			}
		}
		// states-by-user
		{
			userBucketName := getStatesBucketNameForUser(state.Provider.SteamID)
			userBucket, err := tx.CreateBucketIfNotExists(userBucketName)
			if err != nil {
				return err
			}
			id, err := userBucket.NextSequence()
			if err != nil {
				return err
			}
			err = userBucket.Put(itob(id), state.RawJson)
			if err != nil {
				return err
			}
		}
		return nil
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
	return stateRepo.stateDB.View(func(tx *bolt.Tx) error {
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

// GetAll gets all states in DB.
func (stateRepo *DBStateRepo) GetAll() ([]csgostate.State, error) {
	var states []csgostate.State
	err := stateRepo.stateDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(statesBucketName)
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

// GetAllForPlayer gets all states in DB for a certain player.
func (stateRepo *DBStateRepo) GetAllForPlayer(steamID string) ([]csgostate.State, error) {
	var states []csgostate.State
	err := stateRepo.stateDB.View(func(tx *bolt.Tx) error {
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

// ArchiveMonth writes a month of data to JSON.
func (stateRepo *DBStateRepo) ArchiveMonth(year, month int, w io.Writer) error {
	encoder := json.NewEncoder(w)
	return stateRepo.stateDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(statesBucketName)
		if bucket == nil {
			return nil
		}

		cursor := bucket.Cursor()
		for key, val := cursor.First(); key != nil; key, val = cursor.Next() {
			// Only events from a certain month
			s := csgostate.State{}
			err := json.Unmarshal(val, &s)
			if err != nil {
				return err
			}
			t := time.Unix(int64(s.Provider.Timestamp), 0)
			stateYear, stateMonth, _ := t.Date()
			if stateYear != year || int(stateMonth) != month {
				continue
			}

			// Decode/encode to ensure we only write a single line to output
			var v interface{}
			err = json.Unmarshal(val, &v)
			if err != nil {
				return err
			}
			err = encoder.Encode(v)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteUserDBStuff deletes all state related info from user DB.
func (stateRepo *DBStateRepo) DeleteUserDBStuff() error {
	return stateRepo.userDB.Update(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			fmt.Println("found", string(name))
			if string(name) == "states" || strings.HasPrefix(string(name), "states-") {
				return tx.DeleteBucket(name)
			}
			return nil
		})
	})
}
