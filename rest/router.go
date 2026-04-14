package rest

import (
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"
)

func GetServeMux() http.Handler {
	mux := http.NewServeMux()

	// Public Routes
	mux.HandleFunc("POST /api/v1/getJWT", login)
	mux.HandleFunc("GET /api/vi/refreshJWT", refresh)
	mux.HandleFunc("POST /api/v1/register", postRegister)
	mux.HandleFunc("GET /api/v1/commemorations", getCommemorations)
	mux.HandleFunc("GET /api/v1/chapter", getChapters)
	mux.HandleFunc("GET /api/v1/localities", getLocalities)

	// Member Management
	mux.HandleFunc("GET /api/v1/member/{id}", authMW(getMember, model.AuthLevelView))
	mux.HandleFunc("GET /api/v1/member/{id}/chapters", authMW(getMemberChapters, model.AuthLevelView))
	mux.HandleFunc("PUT /api/v1/member/{id}", authMW(setMember, model.AuthLevelManager))
	mux.HandleFunc("PUT /api/v1/member/{id}/chapters", authMW(setMemberChapters, model.AuthLevelManager))
	mux.HandleFunc("POST /api/v1/member", authMW(createMember, model.AuthLevelAdmin))
	mux.HandleFunc("GET /api/v1/member/{id}/vcard", authMW(getMemberVcard, model.AuthLevelView))

	// Giving & Changelog
	mux.HandleFunc("GET /api/v1/giving/{id}", authMW(getMemberGiving, model.AuthLevelAdmin))
	mux.HandleFunc("POST /api/v1/giving/{id}", authMW(postMemberGiving, model.AuthLevelAdmin))
	mux.HandleFunc("GET /api/v1/changelog/{id}", authMW(getMemberChangelog, model.AuthLevelFullView))

	// Notes
	mux.HandleFunc("GET /api/v1/member/{id}/notes", authMW(getNotes, model.AuthLevelFullView))
	mux.HandleFunc("POST /api/v1/member/{id}/notes", authMW(postNote, model.AuthLevelManager))
	mux.HandleFunc("DELETE /api/v1/member/{id}/notes/{noteid}", authMW(deleteNote, model.AuthLevelManager))

	// Self-Service
	mux.HandleFunc("GET /api/v1/me", authMW(getMe, model.AuthLevelView))
	mux.HandleFunc("PUT /api/v1/me", authMW(setMe, model.AuthLevelView))
	mux.HandleFunc("GET /api/v1/me/chapters", authMW(getMeChapters, model.AuthLevelView))
	mux.HandleFunc("PUT /api/v1/me/chapters", authMW(setMeChapters, model.AuthLevelView))
	mux.HandleFunc("GET /api/v1/me/giving", authMW(getMeGiving, model.AuthLevelView))

	// Subscribers & Search
	mux.HandleFunc("GET /api/v1/subscriber/{id}", authMW(getSubscriber, model.AuthLevelView))
	mux.HandleFunc("POST /api/v1/subscriber/{id}", authMW(setSubscriber, model.AuthLevelAdmin))
	mux.HandleFunc("POST /api/v1/search", authMW(postSearch, model.AuthLevelView))
	mux.HandleFunc("POST /api/v1/searchemail", authMW(postEmailSearch, model.AuthLevelView))
	mux.HandleFunc("POST /api/v1/subsearch", authMW(postSubSearch, model.AuthLevelView))

	// Reports & Dashboard
	mux.HandleFunc("GET /api/v1/report/{report}", authMW(reports, model.AuthLevelFullView))
	mux.HandleFunc("GET /api/v1/dashboard", authMW(getDashboard, model.AuthLevelView))
	mux.HandleFunc("GET /api/v1/necrology", authMW(getNecrology, model.AuthLevelView))
	mux.HandleFunc("POST /api/v1/email", authMW(postEmail, model.AuthLevelAdmin))

	// Chapters & Localities
	mux.HandleFunc("PUT /api/v1/chapter/{id}", authMW(putChapter, model.AuthLevelAdmin))
	mux.HandleFunc("GET /api/v1/chapter/{id}", authMW(getChapterMembers, model.AuthLevelView))
	mux.HandleFunc("GET /api/v1/locality/{joint}", authMW(getLocalityMembers, model.AuthLevelView))
	mux.HandleFunc("GET /api/v1/leaders/{category}", authMW(getLeadership, model.AuthLevelView))

	// Prayer Wall
	mux.HandleFunc("GET /api/v1/prayers", getPublicPrayers);
	mux.HandleFunc("GET /api/v1/member/{id}/prayers", authMW(getPrayers, model.AuthLevelView));
	mux.HandleFunc("POST /api/v1/prayers", authMW(addPrayer, model.AuthLevelView));
	mux.HandleFunc("DELETE /api/v1/prayers/{id}", authMW(deletePrayer, model.AuthLevelView));

	// Wrap the entire mux in global middleware (Headers/CORS)
	return globalMW(mux)
}

// globalMW handles headers and CORS for every single request
func globalMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "OSL-Member-Manager")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS, HEAD, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		// Handle Preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

/*
func notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
} */
