package userrepo

import (
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
	"os"
	"testing"
)

func TestDBUserRepo(t *testing.T) {
	os.Remove("test.db")
	db, err := bolt.Open("test.db", 0666, nil)
	require.NoError(t, err)

	userRepo, err := NewDBUserRepo(db, "asdf")
	require.NoError(t, err)

	userMissingSteamID := &User{
		SteamID:   "",
	}
	err = userRepo.Create(userMissingSteamID)
	require.Error(t, err)

	user1 := &User{
		SteamID:   "asdf1234",
		NickName:  "Andy",
		AvatarURL: "None",
	}
	err = userRepo.Create(user1)
	require.NoError(t, err)
	require.True(t, len(user1.PushToken) > 10, user1.PushToken)

	user2, err := userRepo.GetBySteamID("asdf1234")
	require.NoError(t, err)
	require.Equal(t, user1.SteamID, user2.SteamID)
	require.Equal(t, user1.PushToken, user2.PushToken)

	userNil, err := userRepo.GetBySteamID("nonexistent")
	require.NoError(t, err)
	require.Nil(t, userNil)

	user3, err := userRepo.GetByPushToken(user1.PushToken)
	require.NoError(t, err)
	require.Equal(t, user1.SteamID, user3.SteamID)
	require.Equal(t, user1.PushToken, user3.PushToken)
}
