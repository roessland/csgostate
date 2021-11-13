package main

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
	"sync"
)

type User struct {
	SteamID   string `json:"SteamID"`
	PushToken string `json:"PushToken"`
	NickName  string `json:"NickName"`
	AvatarURL string `json:"AvatarURL"`
}

type UserRepo interface {
	GetBySteamID(steamID string) (*User, error)
	GetByPushToken(pushToken string) (*User, error)
	Create(user *User) error
}

var _ UserRepo = &DBUserRepo{}

type DBUserRepo struct {
	db              *bolt.DB
	pushTokenSecret string
}

const usersBucket = "users"

func NewDBUserRepo(db *bolt.DB, pushTokenSecret string) (*DBUserRepo, error) {
	// Ensure bucket exists
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(usersBucket))
		return err
	})
	if err != nil {
		return nil, err
	}

	return &DBUserRepo{
		db:              db,
		pushTokenSecret: pushTokenSecret,
	}, nil
}

func (userRepo *DBUserRepo) GetBySteamID(steamID string) (*User, error) {
	var user *User
	err := userRepo.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(usersBucket))
		buf := bucket.Get([]byte(steamID))
		if buf == nil {
			return nil
		}
		var u User
		err := json.Unmarshal(buf, &u)
		if err != nil {
			return errors.Wrapf(err, "when decoding %s", string(buf))
		}
		user = &u
		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userRepo *DBUserRepo) GetByPushToken(pushToken string) (*User, error) {
	var user *User
	err := userRepo.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(usersBucket))
		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			var u User
			err := json.Unmarshal(value, &u)
			if err != nil {
				return errors.Wrapf(err, "when decoding %s", string(value))
			}
			if u.PushToken == pushToken {
				user = &u
				break
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepo *DBUserRepo) Create(user *User) error {
	// Generate a deterministic push token.
	user.PushToken = getPushToken(userRepo.pushTokenSecret, user.SteamID)

	// Store the user in the user bucket using the SteamID as the key.
	err := userRepo.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(usersBucket))
		encoded, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(user.SteamID), encoded)
	})
	return err
}

// verify that it satisfies interface
var _ UserRepo = &InMemoryUserRepo{}

type InMemoryUserRepo struct {
	sync.RWMutex
	users           map[string]*User
	pushTokenSecret string
}

func NewInMemoryUserRepo(pushTokenSecret string) *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users:           make(map[string]*User),
		pushTokenSecret: pushTokenSecret,
	}
}

func (userRepo *InMemoryUserRepo) GetBySteamID(steamID string) (*User, error) {
	userRepo.RLock()
	defer userRepo.RUnlock()
	return userRepo.users[steamID], nil
}

func (userRepo *InMemoryUserRepo) GetByPushToken(pushToken string) (*User, error) {
	userRepo.RLock()
	defer userRepo.RUnlock()
	for _, user := range userRepo.users {
		if user.PushToken == pushToken {
			return user, nil
		}
	}
	return nil, nil
}

func (userRepo *InMemoryUserRepo) Create(user *User) error {
	userRepo.Lock()
	defer userRepo.Unlock()
	user.PushToken = getPushToken(userRepo.pushTokenSecret, user.SteamID)
	u := *user
	userRepo.users[user.SteamID] = &u
	return nil
}

func getPushToken(secret string, steamID string) string {
	hash := sha256.New()
	hash.Write([]byte(secret))
	hash.Write([]byte(steamID))
	digest64 := hash.Sum(nil)
	pushTokenUuid, err := uuid.FromBytes(digest64[:16])
	if err != nil {
		panic(err)
	}
	return pushTokenUuid.String()
}
