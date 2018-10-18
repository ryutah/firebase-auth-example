package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"firebase.google.com/go"
	"firebase.google.com/go/auth"

	"cloud.google.com/go/firestore"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type user struct {
	UID   string `json:"uid" firestore:"-"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func init() {
	r := mux.NewRouter()

	r.Path("/users").Methods("POST").HandlerFunc(createUser)
	r.Path("/users/me").Methods("PATCH").HandlerFunc(patchUser)
	r.Path("/users/me").Methods("GET").HandlerFunc(getMe)

	http.Handle("/", r)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	newUser := new(user)
	if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
		log.Warningf(ctx, "failed to create firestore clietn: %#v", err)
		http.Error(w, "invalid request", 400)
		return
	}

	client, err := firestore.NewClient(ctx, appengine.AppID(ctx))
	if err != nil {
		log.Errorf(ctx, "failed to create firestore clietn: %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	if _, err := client.Collection("users").Doc(newUser.UID).Set(ctx, newUser); err != nil {
		log.Errorf(ctx, "failed to create user: %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	jsonResponse(ctx, w, 201, newUser)
}

func patchUser(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	authHeader := r.Header.Get("Authorization")
	headers := strings.SplitN(authHeader, " ", 2)
	if len(headers) != 2 || strings.HasPrefix(strings.ToLower(headers[0]), "bearer") {
		http.Error(w, "unauthorized_request", 401)
		return
	}

	idToken := strings.TrimSpace(headers[1])

	auth, err := authClient(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to create firebase auth client: %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	token, err := auth.VerifyIDTokenAndCheckRevoked(ctx, idToken)
	if err != nil {
		log.Warningf(ctx, "failed to verify id token: %#v", err)
		http.Error(w, "unauthorized_request", 401)
		return
	}
	log.Infof(ctx, "%v", token.Claims)

	uid := token.UID
	newUser := new(user)
	if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
		log.Warningf(ctx, "failed to create firestore clietn: %#v", err)
		http.Error(w, "invalid request", 400)
		return
	}

	client, err := firestore.NewClient(ctx, appengine.AppID(ctx))
	if err != nil {
		log.Errorf(ctx, "failed to create firestore clietn: %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	doc := client.Collection("users").Doc(uid)
	snap, err := doc.Get(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get user: %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	oldUsr := new(user)
	if err := snap.DataTo(oldUsr); err != nil {
		log.Errorf(ctx, "failed to parse user data: %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	newUser.Email = oldUsr.Email

	if err := client.RunTransaction(ctx, func(tc context.Context, tx *firestore.Transaction) error {
		if err := tx.Set(doc, newUser); err != nil {
			return err
		}
		return auth.SetCustomUserClaims(ctx, uid, map[string]interface{}{
			"finished": true,
		})
	}); err != nil {
		log.Errorf(ctx, "failed to update user data: %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	jsonResponse(ctx, w, 200, newUser)
}

func getMe(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	authHeader := r.Header.Get("Authorization")
	headers := strings.SplitN(authHeader, " ", 2)
	if len(headers) != 2 || strings.HasPrefix(strings.ToLower(headers[0]), "bearer") {
		http.Error(w, "unauthorized_request", 401)
		return
	}

	idToken := strings.TrimSpace(headers[1])

	auth, err := authClient(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to create firebase auth client: %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	token, err := auth.VerifyIDTokenAndCheckRevoked(ctx, idToken)
	if err != nil {
		log.Warningf(ctx, "failed to verify id token: %#v", err)
		http.Error(w, "unauthorized_request", 401)
		return
	}
	log.Infof(ctx, "%v", token.Claims)

	client, err := firestore.NewClient(ctx, appengine.AppID(ctx))
	if err != nil {
		log.Warningf(ctx, "failed to create firestore clietn: %#v", err)
		http.Error(w, "invalid request", 400)
		return
	}

	snap, err := client.Collection("users").Doc(token.UID).Get(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get user: %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	usr := new(user)
	if err := snap.DataTo(usr); err != nil {
		log.Errorf(ctx, "failed to parse user data: %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	usr.UID = snap.Ref.ID
	jsonResponse(ctx, w, 200, usr)
}

func jsonResponse(ctx context.Context, w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Errorf(ctx, "failed to render json: %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func authClient(ctx context.Context) (*auth.Client, error) {
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: appengine.AppID(ctx),
	})
	if err != nil {
		return nil, err
	}
	return app.Auth(ctx)
}
