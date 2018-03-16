package auth

import (
	"context"

	"github.com/ryutah/firebase-auth-example/server/common"

	firebase "firebase.google.com/go"
)

type Authenticator interface {
	Auth() (userID string, err error)
}

type firebaseAuthenticator struct {
	ctx     context.Context
	logger  common.Logger
	idToken string
}

func (f *firebaseAuthenticator) Auth() (string, error) {
	app, err := firebase.NewApp(f.ctx, nil)
	if err != nil {
		f.logger.Error("failed to create firebase app : %v", err)
		return "", err
	}
	cli, err := app.Auth(f.ctx)
	if err != nil {
		f.logger.Error("failed to create auth client : %v", err)
		return "", err
	}

	token, err := cli.VerifyIDToken(f.idToken)
	if err != nil {
		f.logger.Warn("failed to verify id token : %v", err)
		return "", ErrAuthenticate
	}
	f.logger.Debug("Token : %#v", token)
	f.logger.Debug("Claims : %v", token.Claims)

	return token.UID, nil
}
