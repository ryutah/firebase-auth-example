package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ryutah/firebase-auth-example/server"
)

func init() {
	r := mux.NewRouter()

	r.Handle("/sample", server.NewSampleHandler()).Methods(server.SampleHandlerMethods()...)

	http.Handle("/", r)
}
