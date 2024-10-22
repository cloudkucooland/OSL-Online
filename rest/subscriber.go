package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/julienschmidt/httprouter"
)

func getSubscriber(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	sid := model.SubscriberID(id)
	m, err := sid.Get()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

func setSubscriber(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := r.ParseMultipartForm(1024 * 64); err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	field := r.PostFormValue("field")
	if field == "" {
		err := fmt.Errorf("field not set")
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	value := r.PostFormValue("value")

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	sid := model.SubscriberID(id)
	if err := sid.SetField(field, value); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	headers(w, r)
	fmt.Fprint(w, jsonStatusOK)
}
