package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"

	"github.com/julienschmidt/httprouter"
)

func postSubSearch(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if err := req.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		http.Error(res, jsonError(err), http.StatusNotAcceptable)
		return
	}

	query := req.PostFormValue("query")
	if query == "" {
		err := fmt.Errorf("query not set")
		slog.Error(err.Error())
		http.Error(res, jsonError(err), http.StatusNotAcceptable)
		return
	}

	// XXX min length or other checks?

	result, err := model.SubscriberSearch(query)
	if err != nil {
		slog.Warn(err.Error())
		http.Error(res, jsonError(err), http.StatusInternalServerError)
		return
	}

	headers(res, req)
	if err := json.NewEncoder(res).Encode(result); err != nil {
		slog.Warn(err.Error())
		http.Error(res, jsonError(err), http.StatusInternalServerError)
		return
	}
}
