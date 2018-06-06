package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dist/index.html")
	})
	r.Path("/static/{rest}").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("dist/"))),
	)

	r.Path("/signin").Methods("POST").HandlerFunc(signin)
	r.Path("/signout").Methods("POST").HandlerFunc(signout)
	r.Path("/verify").Methods("GET").HandlerFunc(verifySession)

	log.Println("Start server on port 8080...")
	http.ListenAndServe(":8080", r)
}

func signin(w http.ResponseWriter, r *http.Request) {
	idToken := r.FormValue("idToken")
	s, err := createSession(idToken, 1*time.Minute)
	if err != nil {
		log.Printf("failed to create session: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	resp, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Printf("failed to write json: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "session_state",
		Value: s.SessionState,
	})
	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func signout(w http.ResponseWriter, r *http.Request) {
	state, err := r.Cookie("session_state")
	if err != nil {
		log.Println("session is not set")
		http.Error(w, "no session exists", 403)
		return
	}

	deleteSession(state.Value)
	state.Expires = time.Now().Add(-1 * time.Second)
	http.SetCookie(w, state)
	fmt.Fprintln(w, "Succes sign out")
}

func verifySession(w http.ResponseWriter, r *http.Request) {
	state, err := r.Cookie("session_state")
	if err != nil {
		log.Println("session is not set")
		http.Error(w, "no session exists", 403)
		return
	}

	s, err := getSession(state.Value)
	if err != nil {
		log.Printf("failed to get session %v", err)
		http.Error(w, "invalid session state", 403)
		return
	}
	if err := s.valid(); err != nil {
		log.Println("session expired")
		http.Error(w, "session expired", 403)
		return
	}
	fmt.Fprintln(w, "session is valid")
}
