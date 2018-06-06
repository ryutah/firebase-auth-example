package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"math/big"
	"time"
)

var (
	errSessionExpired      = errors.New("session expired")
	errInvalidSessionState = errors.New("invalid session state")
)

var sessions = make(map[string]*session)

type session struct {
	SessionState string    `json:"session_state"`
	IDToken      string    `json:"-"`
	ExpiredAt    time.Time `json:"expired_at"`
}

func (s *session) valid() error {
	if s.ExpiredAt.IsZero() {
		return nil
	}

	now := time.Now().UTC()
	if s.ExpiredAt.Before(now) {
		return errSessionExpired
	}
	return nil
}

func getSession(state string) (*session, error) {
	s, ok := sessions[state]
	if !ok {
		return nil, errInvalidSessionState
	}
	return s, nil
}

func deleteSession(state string) {
	delete(sessions, state)
}

func createSession(idToken string, expired time.Duration) (*session, error) {
	state, err := createSessionState()
	log.Println(state)
	if err != nil {
		return nil, err
	}
	s := &session{
		SessionState: state,
		IDToken:      idToken,
	}
	if expired != 0 {
		s.ExpiredAt = time.Now().Add(expired)
	}
	sessions[state] = s
	return s, nil
}

func createSessionState() (string, error) {
	var result string
	for len(result) < 32 {
		i, err := rand.Int(rand.Reader, big.NewInt(10000))
		if err != nil {
			return "", err
		}
		result += string(i.Int64())
	}
	return base64.StdEncoding.EncodeToString([]byte(result)), nil
}
