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
	m.GET("/api/v1/member/:id/chapters", authMW(getMemberChapters, AuthLevelView))
	m.POST("/api/v1/member/:id", authMW(setMember, AuthLevelManager)) // should use PUT
	m.PUT("/api/v1/member/:id", authMW(setMember, AuthLevelManager))
	m.GET("/api/v1/member/:id/vcard", authMW(getMemberVcard, AuthLevelView))
	m.POST("/api/v1/member", authMW(createMember, AuthLevelAdmin))
	m.PUT("/api/v1/member/:id/chapters", authMW(setMemberChapters, AuthLevelAdmin))

	// manage giving records
	m.GET("/api/v1/giving/:id", authMW(getMemberGiving, AuthLevelAdmin))
	m.POST("/api/v1/giving/:id", authMW(postMemberGiving, AuthLevelAdmin))
	m.GET("/api/v1/changelog/:id", authMW(getMemberChangelog, AuthLevelAdmin))

	// self-service (not complete)
	m.GET("/api/v1/me", authMW(getMe, AuthLevelView))
	m.GET("/api/v1/me/chapters", authMW(getMeChapters, AuthLevelView))
	m.POST("/api/v1/me", authMW(setMe, AuthLevelView)) // should use PUT
	m.PUT("/api/v1/me", authMW(setMe, AuthLevelView))
	m.PUT("/api/v1/me/chapters", authMW(setMeChapters, AuthLevelView))

	// for instituions and individuals who subscribe to Doxology but aren't vowed
	m.GET("/api/v1/subscriber/:id", authMW(getSubscriber, AuthLevelView))
	m.POST("/api/v1/subscriber/:id", authMW(setSubscriber, AuthLevelAdmin))
	// m.POST("/api/v1/subscriber", authMW(createSubscriber, AuthLevelAdmin))

	// search lists
	m.POST("/api/v1/search", authMW(postSearch, AuthLevelView))           // members
	m.POST("/api/v1/searchemail", authMW(postEmailSearch, AuthLevelView)) // members by email address
	m.POST("/api/v1/subsearch", authMW(postSubSearch, AuthLevelView))     // subscribers

	// reports
	m.GET("/api/v1/report/:report", authMW(reports, AuthLevelManager))

	m.POST("/api/v1/email", authMW(postEmail, AuthLevelAdmin))

	// manage chapters
	m.GET("/api/v1/chapter", getChapters) // public
	m.PUT("/api/v1/chapter/:id", authMW(postChapter, AuthLevelAdmin))
	m.GET("/api/v1/chapter/:id", authMW(getChapterMembers, AuthLevelView))

	m.GET("/api/v1/localities", getLocalities) // public
	m.GET("/api/v1/locality/:joint", authMW(getLocalityMembers, AuthLevelView))

	// leadership
	m.GET("/api/v1/leaders/:category", authMW(getLeadership, AuthLevelView))

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
