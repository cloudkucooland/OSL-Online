package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

const staticDir = "/www/oo/static"

func getServeMux() *httprouter.Router {
	m := httprouter.New()
	m.HandleOPTIONS = true
	m.GlobalOPTIONS = http.HandlerFunc(headers)

	if _, err := os.Stat(staticDir); err != nil {
		slog.Error(err.Error())
		panic(err.Error())
	}
	m.ServeFiles("/static/*filepath", http.Dir(staticDir))
	appDir := fmt.Sprintf("%s/_app", staticDir)
	m.ServeFiles("/_app/*filepath", http.Dir(appDir))

	m.NotFound = http.HandlerFunc(notFound)

	// URL to login, returns the JWT to pass in to authenticated endpoints
	m.POST("/api/v1/getJWT", login)

	m.GET("/api/v1/member/:id", authMW(getMember, AuthLevelView))

	m.GET("/api/v1/subscriber/:id", authMW(getSubscriber, AuthLevelView))

	return m
}

func headers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS, HEAD, DELETE, PATCH")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Accept, If-Modified-Since, If-Match, If-None-Match, Authorization")

	w.Header().Set("Content-Type", jsonType)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	// default: redirect to webui
	if r.URL.String() == "" || r.URL.String() == "/" {
		http.Redirect(w, r, "/static/index.html", http.StatusMovedPermanently)
		return
	}

	// if static, but not found, treat it as extra info on the webui
	if strings.HasPrefix(r.URL.String(), "/static") {
		slog.Info("unknown static URL requested", "url", r.URL.String())
		url := "member"
		newLoc := fmt.Sprintf("/static/index.html?u=%s", url)

		http.Redirect(w, r, newLoc, http.StatusMovedPermanently)
		return
	}

	// something unexpected (not /static or /_app) requested, look under /static for it
	newLoc := fmt.Sprintf("/static/index.html?u=%s", r.URL)
	slog.Debug("not found, redirecting", "request", r.URL.String(), "new", newLoc, "method", r.Method)
	http.Redirect(w, r, newLoc, http.StatusMovedPermanently)
}
