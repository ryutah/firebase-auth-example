package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/ryutah/firebase-auth-example/webui-example/server/auth"
	"github.com/ryutah/firebase-auth-example/webui-example/server/common"

	"google.golang.org/appengine"
)

type sampleHandler struct {
	ctx    context.Context
	w      http.ResponseWriter
	r      *http.Request
	auth   auth.Injector
	common common.Injector
}

func NewSampleHandler() http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost}),
	)(new(sampleHandler))
}

func SampleHandlerMethods() []string {
	return []string{http.MethodGet, http.MethodPost, http.MethodOptions}
}

func (s *sampleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.ctx, s.w, s.r = appengine.NewContext(r), w, r
	s.auth = auth.NewInjector()
	s.common = common.NewInjector()

	switch r.Method {
	case http.MethodGet:
		s.get()
	case http.MethodPost:
		s.post()
	case http.MethodOptions:
		break
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *sampleHandler) get() {
	s.w.Header().Set("Content-Type", "application/json")

	idToken := s.r.Header.Get("Authorization")
	if idToken == "" {
		s.common.GAELogger(s.ctx).Warn("Authorization header is blank")
		s.responseAuthResult(false, "")
		return
	}

	userID, err := s.auth.FirebaseAuthenticator(s.ctx, idToken).Auth()
	if err == auth.ErrAuthenticate {
		s.responseAuthResult(false, "")
		return
	} else if err != nil {
		http.Error(s.w, "Unauthorized", http.StatusInternalServerError)
		return
	}

	s.common.GAELogger(s.ctx).Info("User ID : %v", userID)
	s.responseAuthResult(true, userID)
}

func (s *sampleHandler) post() {
	resp := map[string]string{"foo": "this is post"}
	s.w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(s.w).Encode(resp); err != nil {
		s.common.GAELogger(s.ctx).Error("failed to write json : %v", err)
		http.Error(s.w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *sampleHandler) responseAuthResult(verify bool, userID string) {
	resp := map[string]interface{}{
		"verify": verify,
		"userID": userID,
	}
	if !verify {
		s.w.WriteHeader(http.StatusUnauthorized)
	}
	if err := json.NewEncoder(s.w).Encode(resp); err != nil {
		s.common.GAELogger(s.ctx).Error("failed to write json : %v", err)
		http.Error(s.w, err.Error(), http.StatusInternalServerError)
	}
}
