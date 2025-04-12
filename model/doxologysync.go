package model

import (
	"context"
	"log/slog"
	"strings"

	"google.golang.org/api/admin/directory/v1"
)

func DoxologySync(ctx context.Context) error {
	// assumes GOOGLE_APPLICATION_CREDENTIALS enviornment is set.
	adminService, err := admin.NewService(ctx)
	if err != nil {
		return err
	}

	call := adminService.Members.List("doxology@saint-luke.net")
	known := make(map[string]bool, 0)
	err = call.Pages(ctx, func(members *admin.Members) error {
		for _, m := range members.Members {
			e := strings.ToLower(m.Email)
			ok, err := checkDoxology(e)
			if err != nil {
				slog.Error("doxology", "error", err.Error())
				continue
			}
			if !ok {
				slog.Info("doxology", "message", "removing", "user", e)
				if err := adminService.Members.Delete("doxology@saint-luke.net", e).Do(); err != nil {
					slog.Error("doxology", "error", err.Error())
					continue
				}
			} else {
				known[e] = true
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	toadd, err := DoxologyEmailedDirect()
	for _, m := range toadd {
		m = strings.ToLower(m)
		if _, ok := known[m]; !ok {
			slog.Info("doxology", "message", "adding user", "email", m)
			if _, err := adminService.Members.Insert("doxology@saint-luke.net", &admin.Member{Email: m}).Do(); err != nil {
				slog.Error("doxology", "error", err.Error())
				continue
			}
		}
	}
	return nil
}

func checkDoxology(email string) (bool, error) {
	found, err := SearchEmail(email, true)
	if err != nil {
		return false, err
	}
	if len(found) == 0 { // should be != 1 but Br Dan and Sr. Mary-O share an address
		return checkDoxologySubscriber(email)
	}
	member, err := found[0].ID.Get()
	if err != nil {
		return false, err
	}
	if strings.ToLower(member.Doxology) == "none" {
		return false, nil
	}

	return true, nil
}

func checkDoxologySubscriber(email string) (bool, error) {
	found, err := SubscriberSearchEmail(email)
	if err != nil {
		return false, err
	}
	if len(found) == 0 {
		return false, nil
	}
	subscriber, err := found[0].ID.Get()
	if err != nil {
		return false, err
	}
	if strings.ToLower(subscriber.Doxology) == "none" {
		return false, nil
	}

	return true, nil
}
