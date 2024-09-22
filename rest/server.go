package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var srv *http.Server
var sk jwk.Set

const sessionName string = "OSL-Online"
const jsonType = "application/json; charset=UTF-8"
const jsonTypeShort = "application/json"
const jsonStatusOK = `{"status":"ok"}`
const BcryptRounds = 14

type authLevel uint8

const (
	AuthLevelView    authLevel = iota // view  members/subscribers
	AuthLevelManager                  // change/add members/subscribers
	AuthLevelAdmin                    // add users
)

// Start launches the HTTP server which is responsible for the frontend and the HTTP API.
func Start(ctx context.Context) {
	srv = &http.Server{
		Handler:           getServeMux(),
		Addr:              ":8443", // make config file...
		WriteTimeout:      (30 * time.Second),
		ReadTimeout:       (30 * time.Second),
		ReadHeaderTimeout: (2 * time.Second),
	}

	// creates the keys if needed
	sk = getJWSigningKeys()

	slog.Info("Starting up REST server", "on", ":8443")
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			slog.Error(err.Error())
			panic(err.Error())
		}
	}()

	<-ctx.Done()
	slog.Info("Shutting down REST server")
	if err := srv.Shutdown(context.Background()); err != nil {
		slog.Error(err.Error())
	}
}

func jsonError(e error) string {
	return fmt.Sprintf(`{"status":"error","error":"%s"}`, e.Error())
}

func contentTypeIs(req *http.Request, check string) bool {
	contentType := strings.Split(strings.Replace(req.Header.Get("Content-Type"), " ", "", -1), ";")[0]
	return strings.EqualFold(contentType, check)
}

func authMW(h httprouter.Handle, requiredlevel authLevel) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token, err := jwt.ParseRequest(r,
			jwt.WithKeySet(sk, jws.WithInferAlgorithmFromKey(true), jws.WithUseDefault(true)),
			jwt.WithValidate(true),
			jwt.WithAudience(sessionName),
			jwt.WithAcceptableSkew(20*time.Second),
		)
		if err != nil {
			slog.Error("token parse/validate failed", "error", err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		username := string(token.Subject())
		claim, ok := token.Get("level")
		if !ok {
			err := fmt.Errorf("no level claim in token")
			slog.Error(err.Error(), "user", username, "claim", claim, "type", fmt.Sprintf("%T", claim))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		checklevel, ok := claim.(float64) // why does this come across as float64?
		if !ok {
			err := fmt.Errorf("authlevel type assertion failed")
			slog.Error(err.Error(), "user", username, "claim", claim, "type", fmt.Sprintf("%T", claim))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if authLevel(checklevel) < requiredlevel {
			err := fmt.Errorf("access level too low")
			slog.Warn(err.Error(), "wanted", requiredlevel, "got", checklevel, "username", username)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		h(w, r, ps)
	}
}

func getUser(r *http.Request) string {
	token, err := jwt.ParseRequest(r,
		jwt.WithKeySet(sk, jws.WithInferAlgorithmFromKey(true), jws.WithUseDefault(true)),
		jwt.WithValidate(true),
		jwt.WithAudience(sessionName),
		jwt.WithAcceptableSkew(20*time.Second),
	)
	if err != nil {
		slog.Error("token parse/validate failed", "error", err.Error())
		return ""
	}

	username := string(token.Subject())
	return username
}
