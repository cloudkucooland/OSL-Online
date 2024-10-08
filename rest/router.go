package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func getServeMux() *httprouter.Router {
	m := httprouter.New()
	m.HandleOPTIONS = true
	m.GlobalOPTIONS = http.HandlerFunc(headers)

	m.NotFound = http.HandlerFunc(notFound)

	// URL to login, returns the JWT to pass in to authenticated endpoints
	m.POST("/api/v1/getJWT", login)
	// takes an email address, returns an "OK" after the password message is sent
	m.POST("/api/v1/register", postRegister)

	// manage individual members
	m.GET("/api/v1/member/:id", authMW(getMember, AuthLevelView))
	m.POST("/api/v1/member/:id", authMW(setMember, AuthLevelManager))
	m.POST("/api/v1/member", authMW(createMember, AuthLevelAdmin))

	// manage giving records
	m.GET("/api/v1/giving/:id", authMW(getMemberGiving, AuthLevelAdmin))
	m.POST("/api/v1/giving/:id", authMW(postMemberGiving, AuthLevelAdmin))
	m.GET("/api/v1/changelog/:id", authMW(getMemberChangelog, AuthLevelAdmin))

	// self-service (not complete)
	m.GET("/api/v1/me", authMW(getMe, AuthLevelView))
	m.POST("/api/v1/me", authMW(setMe, AuthLevelView))

	// for instituions and individuals who subscribe to Doxology but aren't vowed
	m.GET("/api/v1/subscriber/:id", authMW(getSubscriber, AuthLevelView))
	m.POST("/api/v1/subscriber/:id", authMW(setSubscriber, AuthLevelAdmin))
	// m.POST("/api/v1/subscriber", authMW(createSubscriber, AuthLevelAdmin))

	// search lists
	m.POST("/api/v1/search", authMW(postSearch, AuthLevelView))
	m.POST("/api/v1/subsearch", authMW(postSubSearch, AuthLevelView))

	// reports - ad-hoc, but works
	m.GET("/api/v1/report/notrenewed", authMW(reportNotrenewed, AuthLevelManager))
	m.GET("/api/v1/report/expired", authMW(reportExpired, AuthLevelManager))
	m.GET("/api/v1/report/email", authMW(reportEmail, AuthLevelManager))
	m.GET("/api/v1/report/annual", authMW(reportAnnual, AuthLevelManager))
	m.GET("/api/v1/report/life", authMW(reportLife, AuthLevelManager))
	return m
}

func headers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "OSL-Member-Manager")
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
	headers(w, r)
	http.Error(w, "Not Found", http.StatusNotFound)
}
