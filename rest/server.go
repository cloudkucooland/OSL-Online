package rest

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/cloudkucooland/OSL-Online/model"
)

var srv *http.Server
var sk jwk.Set

const sessionName string = "OSL-Online"
const jsonType = "application/json; charset=UTF-8"
const jsonStatusOK = `{"status":"ok"}`
const BcryptRounds = 14

func Start(ctx context.Context) {
	srv = &http.Server{
		Handler:           GetServeMux(),
		Addr:              ":8443",
		WriteTimeout:      (30 * time.Second),
		ReadTimeout:       (30 * time.Second),
		ReadHeaderTimeout: (2 * time.Second),
		ErrorLog:          newServerErrorLog(),
	}

	cert := "/etc/letsencrypt/live/saint-luke.net/fullchain.pem"
	key := "/etc/letsencrypt/live/saint-luke.net/privkey.pem"

	sk = getJWSigningKeys()

	slog.Info("Starting up REST server", "on", ":8443")
	go func() {
		if err := srv.ListenAndServeTLS(cert, key); err != http.ErrServerClosed {
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

func authMW(next http.HandlerFunc, requiredlevel model.AuthLevel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := parsetoken(r)
		if err != nil {
			slog.Error("token parse/validate failed", "error", err.Error())
			http.Error(w, jsonError(err), http.StatusUnauthorized)
			return
		}

		claim, ok := token.Get("level")
		if !ok {
			http.Error(w, `{"error":"no level claim"}`, http.StatusInternalServerError)
			return
		}

		checklevel, ok := claim.(float64)
		if !ok || model.AuthLevel(checklevel) < requiredlevel {
			slog.Warn("access level too low", "user", token.Subject())
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		username := model.Authname(token.Subject())
		uid, err := username.GetID()
		if err != nil {
			slog.Error("failed to get UID from token subject", "sub", token.Subject())
			http.Error(w, jsonError(err), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), model.CtxKeyID, model.MemberID(uid))
		ctx = context.WithValue(ctx, model.CtxKeyLevel, model.AuthLevel(checklevel))

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func parsetoken(r *http.Request) (jwt.Token, error) {
	return jwt.ParseRequest(r,
		jwt.WithHeaderKey("Authorization"),
		jwt.WithCookieKey("jwt"),
		jwt.WithKeySet(sk, jws.WithInferAlgorithmFromKey(true), jws.WithUseDefault(true)),
		jwt.WithValidate(true),
		jwt.WithAudience(sessionName),
		jwt.WithAcceptableSkew(20*time.Second),
	)
}

func getUser(r *http.Request) string {
	if changer, ok := r.Context().Value(model.CtxKeyID).(model.MemberID); ok {
		return fmt.Sprintf("%d", changer)
	}

	token, err := parsetoken(r)
	if err != nil {
		return ""
	}
	return string(token.Subject())
}

func jsonError(e error) string {
	return fmt.Sprintf(`{"status":"error","error":"%s"}`, e.Error())
}

type serverErrorLogWriter struct{}

func (*serverErrorLogWriter) Write(p []byte) (int, error) {
	m := string(p)
	if strings.HasPrefix(m, "http: TLS handshake error") {
		slog.Debug(m)
	} else {
		slog.Error(m)
	}
	return len(p), nil
}

func newServerErrorLog() *log.Logger {
	return log.New(&serverErrorLogWriter{}, "", 0)
}
