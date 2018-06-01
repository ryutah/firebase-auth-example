package main

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
)

type myClaim struct {
	*jwt.StandardClaims
	AuthTime int64 `json:"auth_time,omitempty"`
}

const jwtKeysURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/publicKeys"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "./dist/index.html")
		} else {
			http.ServeFile(w, r, "./dist/"+r.URL.Path[1:])
		}
	})
	http.HandleFunc("/createsession", createSessionCookie)
	http.HandleFunc("/verifysession", verifySessionCookie)
	http.ListenAndServe(":8080", nil)
}

func createSessionCookie(w http.ResponseWriter, r *http.Request) {
	client, err := createClient()
	if err != nil {
		log.Fatalf("failed to create client %#v", err)
	}
	rbody, _ := ioutil.ReadAll(r.Body)
	reqmap := make(map[string]interface{})
	json.Unmarshal(rbody, &reqmap)
	log.Printf("ID TOKEN : %v", reqmap["idToken"])

	rmap := map[string]interface{}{
		"idToken":       reqmap["idToken"],
		"validDuration": 5 * time.Minute / time.Second,
	}
	reqBody, _ := json.Marshal(rmap)
	req, err := http.NewRequest(
		"POST",
		"https://www.googleapis.com/identitytoolkit/v3/relyingparty/createSessionCookie",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		log.Printf("failed to create request: %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to send request: %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response: %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	log.Println(string(respBody))

	var respMap map[string]interface{}
	if err := json.Unmarshal(respBody, &respMap); err != nil {
		log.Printf("failed to unmarshal json body: %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	sessionCookie := respMap["sessionCookie"].(string)
	cookie := http.Cookie{
		Name:    "session",
		Value:   sessionCookie,
		Expires: time.Now().Add(30 * time.Minute),
	}
	http.SetCookie(w, &cookie)
	w.Write(respBody)
}

func verifySessionCookie(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		log.Printf("failed to get cookie %v", err)
		http.Error(w, err.Error(), 403)
		return
	}

	token, err := jwt.ParseWithClaims(sessionCookie.Value, new(myClaim), func(tkn *jwt.Token) (interface{}, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tkn.Header["alg"])
		}
		keys, err := getRSAKeys()
		if err != nil {
			return nil, err
		}
		key, ok := keys[tkn.Header["kid"].(string)]
		if !ok {
			return nil, fmt.Errorf("unexpected kid: %v", tkn.Header["kid"])
		}
		keyBlock, _ := pem.Decode([]byte(key))
		if keyBlock == nil {
			return nil, errors.New("failed to decode public key")
		}
		certificate, err := x509.ParseCertificate(keyBlock.Bytes)
		if err != nil {
			return nil, err
		}
		return certificate.PublicKey, nil
	})
	if err != nil {
		log.Printf("failed parse jwt token: %#v, %v", err, err.Error())
		http.Error(w, err.Error(), 403)
		return
	}
	if token.Valid {
		log.Printf("Token: %v", token)
	}
	claim := token.Claims.(*myClaim)
	cJSON, _ := json.MarshalIndent(claim, "", "  ")
	log.Printf(string(cJSON))
	log.Printf("%+v", token.Header)

	aclient, _ := createAuthClient()
	usr, err := aclient.GetUser(context.Background(), claim.Subject)
	if err != nil {
		log.Printf("failed to get user: %#v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	if claim.IssuedAt*1000 < usr.TokensValidAfterMillis {
		log.Println("revoked token")
		http.Error(w, "revoked session", 403)
		return
	}

	usrJSON, _ := json.MarshalIndent(usr, "", "  ")
	log.Println(string(usrJSON))
	fmt.Fprintln(w, "Success verify!!")
}

func createClient() (*http.Client, error) {
	client, _, err := transport.NewHTTPClient(
		context.Background(),
		option.WithScopes(
			"https://www.googleapis.com/auth/cloud-platform",
			"https://www.googleapis.com/auth/datastore",
			"https://www.googleapis.com/auth/devstorage.full_control",
			"https://www.googleapis.com/auth/firebase",
			"https://www.googleapis.com/auth/identitytoolkit",
			"https://www.googleapis.com/auth/userinfo.email",
		),
		option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")),
	)
	return client, err
}

func createAuthClient() (*auth.Client, error) {
	opts := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: os.Getenv("PROJECT_ID"),
	}, opts)
	if err != nil {
		log.Fatalf("faild to create firebase app %#v", err)
	}
	return app.Auth(context.Background())
}

func getRSAKeys() (map[string]string, error) {
	resp, err := http.Get(jwtKeysURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}
