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
	// gets a new JWT if the current one is about to expire
	m.GET("/api/vi/refreshJWT", refresh)
	// takes an email address, returns an "OK" after the password message is sent
	m.POST("/api/v1/register", postRegister)

	// public -- pulled by WADO
	m.GET("/api/v1/commemorations", getCommemorations)

	// manage individual members
	m.GET("/api/v1/member/:id", authMW(getMember, AuthLevelView))
	m.GET("/api/v1/member/:id/chapters", authMW(getMemberChapters, AuthLevelView))
	m.PUT("/api/v1/member/:id", authMW(setMember, AuthLevelManager))
	m.GET("/api/v1/member/:id/vcard", authMW(getMemberVcard, AuthLevelView))
	m.POST("/api/v1/member", authMW(createMember, AuthLevelAdmin))
	m.PUT("/api/v1/member/:id/chapters", authMW(setMemberChapters, AuthLevelManager))

	// manage giving records
	m.GET("/api/v1/giving/:id", authMW(getMemberGiving, AuthLevelAdmin))
	m.POST("/api/v1/giving/:id", authMW(postMemberGiving, AuthLevelAdmin))
	// m.DELETE("/api/v1/giving/:id", authMW(deleteMemberGiving, AuthLevelAdmin))

	m.GET("/api/v1/changelog/:id", authMW(getMemberChangelog, AuthLevelManager))

	m.POST("/api/v1/member/:id/notes", authMW(postNote, AuthLevelManager))
	m.GET("/api/v1/member/:id/notes", authMW(getNotes, AuthLevelManager))
	m.DELETE("/api/v1/member/:id/notes/:noteid", authMW(deleteNote, AuthLevelManager))

	// self-service (not complete)
	m.GET("/api/v1/me", authMW(getMe, AuthLevelView))
	m.PUT("/api/v1/me", authMW(setMe, AuthLevelView))
	m.GET("/api/v1/me/chapters", authMW(getMeChapters, AuthLevelView))
	m.PUT("/api/v1/me/chapters", authMW(setMeChapters, AuthLevelView))
	m.GET("/api/v1/me/giving", authMW(getMeGiving, AuthLevelView))

	// for instituions who subscribe to Doxology
	m.GET("/api/v1/subscriber/:id", authMW(getSubscriber, AuthLevelView))
	m.POST("/api/v1/subscriber/:id", authMW(setSubscriber, AuthLevelAdmin))
	// m.DELETE("/api/v1/subscriber/:id", authMW(deleteSubscriber, AuthLevelAdmin))
	// m.POST("/api/v1/subscriber", authMW(createSubscriber, AuthLevelAdmin))

	// search lists
	m.POST("/api/v1/search", authMW(postSearch, AuthLevelView))           // members
	m.POST("/api/v1/searchemail", authMW(postEmailSearch, AuthLevelView)) // members by email address
	m.POST("/api/v1/subsearch", authMW(postSubSearch, AuthLevelView))     // subscribers

	// reports
	m.GET("/api/v1/report/:report", authMW(reports, AuthLevelManager))
	m.GET("/api/v1/dashboard", authMW(getDashboard, AuthLevelView))
	m.GET("/api/v1/necrology", authMW(getNecrology, AuthLevelView))

	m.POST("/api/v1/email", authMW(postEmail, AuthLevelAdmin))

	// manage chapters
	m.GET("/api/v1/chapter", getChapters) // public
	m.PUT("/api/v1/chapter/:id", authMW(putChapter, AuthLevelAdmin))
	m.GET("/api/v1/chapter/:id", authMW(getChapterMembers, AuthLevelView))
	// m.DELETE("/api/v1/chapter/:id", authMW(deleteChapter, AuthLevelAdmin))

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
