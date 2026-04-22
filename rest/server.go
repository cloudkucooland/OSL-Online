package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
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
			sendError(w, err, http.StatusUnauthorized)
			return
		}

		claim, ok := token.Get("level")
		if !ok {
			sendError(w, fmt.Errorf("no level claim"), http.StatusInternalServerError)
			return
		}

		checklevel, ok := claim.(float64)
		if !ok || model.AuthLevel(checklevel) < requiredlevel {
			slog.Warn("access level too low", "user", token.Subject())
			sendError(w, fmt.Errorf("forbidden"), http.StatusForbidden)
			return
		}

		username := model.Authname(token.Subject())
		uid, err := username.GetID(r.Context())
		if err != nil {
			slog.Error("failed to get UID from token subject", "sub", token.Subject())
			sendError(w, err, http.StatusInternalServerError)
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
	token, err := parsetoken(r)
	if err != nil {
		return ""
	}
	return string(token.Subject())
}

func parseID(r *http.Request, key string) (int, error) {
	val := r.PathValue(key)
	id, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return id, nil
}

func parseUintID(r *http.Request, key string) (uint64, error) {
	val := r.PathValue(key)
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return id, nil
}

func parseIDFromString(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func sendJSON(w http.ResponseWriter, v any) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("json encode failed", "error", err)
		sendError(w, err, http.StatusInternalServerError)
	}
}

func sendError(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, jsonError(err))
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
