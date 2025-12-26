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
	_ = call.Pages(ctx, func(members *admin.Members) error {
		for _, m := range members.Members {
			e := strings.ToLower(m.Email)

			if strings.Contains(e, "@saint-luke.net") {
				known[e] = true
				continue
			}

			ok, err := checkDoxology(ctx, e)
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
			}
			known[e] = true
		}
		return nil
	})

	toadd, err := DoxologyEmailedDirect(ctx)
	if err != nil {
		return err
	}
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

func (id MemberID) UnsubscribeDoxology(ctx context.Context) error {
	// assumes GOOGLE_APPLICATION_CREDENTIALS enviornment is set.
	adminService, err := admin.NewService(ctx)
	if err != nil {
		return err
	}

	m, err := id.Get(ctx)
	if err != nil {
		return err
	}
	if m.PrimaryEmail == "" {
		return nil
	}

	slog.Info("doxology", "message", "removing", "user", m.PrimaryEmail)
	if err := adminService.Members.Delete("doxology@saint-luke.net", m.PrimaryEmail).Do(); err != nil {
		slog.Error("doxology", "error", err.Error())
		return err
	}
	return nil
}

func (id MemberID) SubscribeDoxology(ctx context.Context) error {
	// assumes GOOGLE_APPLICATION_CREDENTIALS enviornment is set.
	adminService, err := admin.NewService(ctx)
	if err != nil {
		return err
	}

	m, err := id.Get(ctx)
	if err != nil {
		return err
	}
	if m.PrimaryEmail == "" {
		return nil
	}

	slog.Info("doxology", "message", "adding user", "email", m.PrimaryEmail)
	if _, err := adminService.Members.Insert("doxology@saint-luke.net", &admin.Member{Email: m.PrimaryEmail}).Do(); err != nil {
		slog.Error("doxology", "error", err.Error())
		return err
	}
	return nil
}

func checkDoxology(ctx context.Context, email string) (bool, error) {
	found, err := SearchEmail(ctx, email, true)
	if err != nil {
		return false, err
	}
	if len(found) == 0 { // should be != 1 but Br Dan and Sr. Mary-O share an address
		return checkDoxologySubscriber(ctx, email)
	}
	member, err := found[0].ID.Get(ctx)
	if err != nil {
		return false, err
	}
	if member.Doxology == NONE {
		return false, nil
	}

	return true, nil
}

func checkDoxologySubscriber(ctx context.Context, email string) (bool, error) {
	found, err := SubscriberSearchEmail(ctx, email)
	if err != nil {
		return false, err
	}
	if len(found) == 0 {
		return false, nil
	}
	subscriber, err := found[0].ID.Get(ctx)
	if err != nil {
		return false, err
	}
	if subscriber.Doxology == NONE {
		return false, nil
	}

	return true, nil
}
