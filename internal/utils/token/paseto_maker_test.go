package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(RandomString(32))
	require.NoError(t, err)

	username := RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

func TestNewPasetoMaker_InvalidKeySize(t *testing.T) {
	maker, err := NewPasetoMaker("short-key")
	require.Error(t, err)
	require.Nil(t, maker)
}

func TestPasetoMaker_InvalidToken(t *testing.T) {
	maker, err := NewPasetoMaker(RandomString(32))
	require.NoError(t, err)

	payload, err := maker.VerifyToken("invalid-token-string")
	require.Error(t, err)
	require.Nil(t, payload)
}

func TestPayload_Valid(t *testing.T) {
	payload, err := NewPayload("test-user", time.Minute)
	require.NoError(t, err)
	require.NotNil(t, payload)

	err = payload.Valid()
	require.NoError(t, err)
}

func TestPayload_Expired(t *testing.T) {
	payload, err := NewPayload("test-user", -time.Minute)
	require.NoError(t, err)
	require.NotNil(t, payload)

	err = payload.Valid()
	require.Error(t, err)
	require.Equal(t, ErrExpiredToken, err)
}
