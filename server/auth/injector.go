package auth

import (
	"context"

	"github.com/ryutah/firebase-auth-example/server/common"
)

type Injector interface {
	FirebaseAuthenticator(ctx context.Context, idToken string) Authenticator
}

type injector struct {
	commonInjector common.Injector
}

func NewInjector() Injector {
	return &injector{
		commonInjector: common.NewInjector(),
	}
}

func (i *injector) FirebaseAuthenticator(ctx context.Context, idToken string) Authenticator {
	return &firebaseAuthenticator{
		logger:  i.commonInjector.GAELogger(ctx),
		ctx:     ctx,
		idToken: idToken,
	}
}
