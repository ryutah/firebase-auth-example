package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/option"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

func main() {
	http.HandleFunc("/getuser", getUserData)
	http.HandleFunc("/idtoken", verifyIDToken)
	http.HandleFunc("/revoke", revokeRefreshToken)
	http.ListenAndServe(":8080", nil)
}

func getUserData(w http.ResponseWriter, r *http.Request) {
	client, err := createAuthClient()
	if err != nil {
		log.Printf("failed to create client %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	uid := r.FormValue("uid")
	log.Printf("UID: %v", uid)
	u, err := client.GetUser(context.Background(), uid)
	if err != nil {
		log.Printf("failed to get user %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	log.Printf("User: %+v", u)
	uJSON, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		log.Printf("failed to marshal user: %#v", err.Error())
	}
	w.Write(uJSON)
}

func verifyIDToken(w http.ResponseWriter, r *http.Request) {
	client, err := createAuthClient()
	if err != nil {
		log.Printf("failed to create client %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	tkn, err := client.VerifyIDTokenAndCheckRevoked(context.Background(), r.FormValue("idtoken"))
	if err != nil {
		log.Printf("failed to verify id token: %#v", err)
		http.Error(w, err.Error(), 401)
		return
	}
	tJSON, _ := json.MarshalIndent(tkn, "", "  ")
	w.Write(tJSON)
}

func revokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	client, err := createAuthClient()
	if err != nil {
		log.Printf("failed to create client %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	if err := client.RevokeRefreshTokens(context.Background(), r.FormValue("uid")); err != nil {
		log.Printf("failed to revoke refresh token %#v", err)
		http.Error(w, err.Error(), 500)
	}
	fmt.Fprintln(w, "Revoke refresh token")
}

func createAuthClient() (*auth.Client, error) {
	opts := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: "[PROJECT_ID]",
	}, opts)
	if err != nil {
		return nil, err
	}
	return app.Auth(context.Background())
}
