package auth

import (
	"context"
	"encoding/json"

	"google.golang.org/appengine"

	"github.com/ryutah/firebase-auth-example/webui-example/server/common"

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
	app, err := firebase.NewApp(f.ctx, &firebase.Config{
		ProjectID: appengine.AppID(f.ctx),
	})
	if err != nil {
		f.logger.Error("failed to create firebase app : %v", err)
		return "", err
	}
	cli, err := app.Auth(f.ctx)
	if err != nil {
		f.logger.Error("failed to create auth client : %v", err)
		return "", err
	}

	token, err := cli.VerifyIDToken(f.ctx, f.idToken)
	if err != nil {
		f.logger.Warn("failed to verify id token : %v", err)
		return "", ErrAuthenticate
	}
	tokenJSON, _ := json.MarshalIndent(token, "", "  ")
	claimsJSON, _ := json.MarshalIndent(token.Claims, "", "  ")
	f.logger.Debug("\n%v", string(tokenJSON))
	f.logger.Debug("\n%v", string(claimsJSON))

	return token.UID, nil
}
