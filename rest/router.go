package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/cloudkucooland/OSL-Online/model"
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
	m.GET("/api/v1/member/:id", authMW(getMember, model.AuthLevelView))
	m.GET("/api/v1/member/:id/chapters", authMW(getMemberChapters, model.AuthLevelView))
	m.PUT("/api/v1/member/:id", authMW(setMember, model.AuthLevelManager))
	m.GET("/api/v1/member/:id/vcard", authMW(getMemberVcard, model.AuthLevelView))
	m.POST("/api/v1/member", authMW(createMember, model.AuthLevelAdmin))
	m.PUT("/api/v1/member/:id/chapters", authMW(setMemberChapters, model.AuthLevelManager))

	// manage giving records
	m.GET("/api/v1/giving/:id", authMW(getMemberGiving, model.AuthLevelAdmin))
	m.POST("/api/v1/giving/:id", authMW(postMemberGiving, model.AuthLevelAdmin))
	// m.DELETE("/api/v1/giving/:id", authMW(deleteMemberGiving, model.AuthLevelAdmin))

	m.GET("/api/v1/changelog/:id", authMW(getMemberChangelog, model.AuthLevelFullView))

	m.POST("/api/v1/member/:id/notes", authMW(postNote, model.AuthLevelManager))
	m.GET("/api/v1/member/:id/notes", authMW(getNotes, model.AuthLevelFullView))
	m.DELETE("/api/v1/member/:id/notes/:noteid", authMW(deleteNote, model.AuthLevelManager))

	// self-service (not complete)
	m.GET("/api/v1/me", authMW(getMe, model.AuthLevelView))
	m.PUT("/api/v1/me", authMW(setMe, model.AuthLevelView))
	m.GET("/api/v1/me/chapters", authMW(getMeChapters, model.AuthLevelView))
	m.PUT("/api/v1/me/chapters", authMW(setMeChapters, model.AuthLevelView))
	m.GET("/api/v1/me/giving", authMW(getMeGiving, model.AuthLevelView))

	// for instituions who subscribe to Doxology
	m.GET("/api/v1/subscriber/:id", authMW(getSubscriber, model.AuthLevelView))
	m.POST("/api/v1/subscriber/:id", authMW(setSubscriber, model.AuthLevelAdmin))
	// m.DELETE("/api/v1/subscriber/:id", authMW(deleteSubscriber, model.AuthLevelAdmin))
	// m.POST("/api/v1/subscriber", authMW(createSubscriber, model.AuthLevelAdmin))

	// search lists
	m.POST("/api/v1/search", authMW(postSearch, model.AuthLevelView))           // members
	m.POST("/api/v1/searchemail", authMW(postEmailSearch, model.AuthLevelView)) // members by email address
	m.POST("/api/v1/subsearch", authMW(postSubSearch, model.AuthLevelView))     // subscribers

	// reports
	m.GET("/api/v1/report/:report", authMW(reports, model.AuthLevelFullView))
	m.GET("/api/v1/dashboard", authMW(getDashboard, model.AuthLevelView))
	m.GET("/api/v1/necrology", authMW(getNecrology, model.AuthLevelView))

	m.POST("/api/v1/email", authMW(postEmail, model.AuthLevelAdmin))

	// manage chapters
	m.GET("/api/v1/chapter", getChapters) // public
	m.PUT("/api/v1/chapter/:id", authMW(putChapter, model.AuthLevelAdmin))
	m.GET("/api/v1/chapter/:id", authMW(getChapterMembers, model.AuthLevelView))
	// m.DELETE("/api/v1/chapter/:id", authMW(deleteChapter, model.AuthLevelAdmin))

	m.GET("/api/v1/localities", getLocalities) // public
	m.GET("/api/v1/locality/:joint", authMW(getLocalityMembers, model.AuthLevelView))

	// leadership
	m.GET("/api/v1/leaders/:category", authMW(getLeadership, model.AuthLevelView))

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
