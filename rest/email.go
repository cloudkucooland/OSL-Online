package rest

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/email"
	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/julienschmidt/httprouter"
)

func postEmail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)
	if err := r.ParseMultipartForm(1024 * 4); err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	message := r.PostFormValue("message")
	if message == "" {
		err := fmt.Errorf("message not set")
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	whom := r.PostFormValue("whom")
	if whom == "" {
		err := fmt.Errorf("whom not set")
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	var ids []model.MemberID
	var err error
	switch whom {
	case "all":
		ids, err = model.ActiveMemberIDs()
	case "annual":
		ids, err = model.AnnualMemberIDs()
	case "life":
		ids, err = model.LifeMemberIDs()
	case "friends":
		ids, err = model.FriendIDs()
	default:
		ids = make([]model.MemberID, 0)
	}
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	if err := email.SendGeneric(ids, message); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}
