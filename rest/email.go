package rest

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/email"
	"github.com/cloudkucooland/OSL-Online/model"
)

func postEmail(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	if err := r.ParseMultipartForm(1024 * 4); err != nil {
		slog.Warn(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	message := r.PostFormValue("message")
	if message == "" {
		sendError(w, fmt.Errorf("message not set"), http.StatusNotAcceptable)
		return
	}

	whom := r.PostFormValue("whom")
	if whom == "" {
		sendError(w, fmt.Errorf("whom not set"), http.StatusNotAcceptable)
		return
	}
	subject := r.PostFormValue("subject")
	if subject == "" {
		sendError(w, fmt.Errorf("subject not set"), http.StatusNotAcceptable)
		return
	}

	var ids []model.MemberID
	var err error
	switch whom {
	case "all":
		ids, err = model.ActiveMemberIDs(r.Context())
	case "annual":
		ids, err = model.AnnualMemberIDs(r.Context())
	case "life":
		ids, err = model.LifeMemberIDs(r.Context())
	case "friends":
		ids, err = model.FriendIDs(r.Context())
	default:
		ids, err = model.TestMemberIDs()
	}
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	slog.Info("sending bulk email", "whom", whom, "subject", subject, "from", getUser(r))

	if err := email.SendGeneric(r.Context(), ids, subject, message); err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}
