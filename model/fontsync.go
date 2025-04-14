package model

import (
	"context"
	"log/slog"
	"strings"

	"google.golang.org/api/admin/directory/v1"
)

func FontSync(ctx context.Context) error {
	// assumes that GOOGLE_APPLICATION_CREDENTIALS enviornment is set

	adminService, err := admin.NewService(ctx)
	if err != nil {
		return err
	}

	call := adminService.Members.List("font@saint-luke.net")
	known := make(map[string]bool, 0)
	_ = call.Pages(ctx, func(members *admin.Members) error {
		for _, m := range members.Members {
			e := strings.ToLower(m.Email)

			if strings.Contains(e, "@saint-luke.net") { // leave these alone
				known[e] = true
				continue
			}

			ok, err := checkFont(e)
			if err != nil {
				slog.Error("font", "error", err.Error())
				continue
			}
			if !ok {
				slog.Info("font", "message", "removing", "user", e)
				if err := adminService.Members.Delete("font@saint-luke.net", e).Do(); err != nil {
					slog.Error("font", "error", err.Error())
					continue
				}
			}
			known[e] = true
		}
		return nil
	})

	toadd, err := FontEmailedDirect()
	if err != nil {
		return err
	}
	for _, m := range toadd {
		m = strings.ToLower(m)
		if _, ok := known[m]; !ok {
			slog.Info("font", "message", "adding user", "email", m)
			if _, err := adminService.Members.Insert("font@saint-luke.net", &admin.Member{Email: m}).Do(); err != nil {
				slog.Error("font", "error", err.Error())
				continue
			}
		}
	}
	return nil
}

func checkFont(email string) (bool, error) {
	found, err := SearchEmail(email, true)
	if err != nil {
		return false, err
	}
	if len(found) == 0 { // should be != 1 but Br Dan and Sr. Mary-O share an address
		return false, nil
	}
	member, err := found[0].ID.Get()
	if err != nil {
		return false, err
	}
	if member.Newsletter == NONE {
		return false, nil
	}

	return true, nil
}

func (id MemberID) UnsubscribeFont(ctx context.Context) error {
	// assumes GOOGLE_APPLICATION_CREDENTIALS enviornment is set.
	adminService, err := admin.NewService(ctx)
	if err != nil {
		return err
	}

	m, err := id.Get()
	if err != nil {
		return err
	}
	if m.PrimaryEmail == "" {
		return nil
	}

	slog.Info("font", "message", "removing", "user", m.PrimaryEmail)
	if err := adminService.Members.Delete("font@saint-luke.net", m.PrimaryEmail).Do(); err != nil {
		slog.Error("font", "error", err.Error())
		return err
	}
	return nil
}

func (id MemberID) SubscribeFont(ctx context.Context) error {
	// assumes GOOGLE_APPLICATION_CREDENTIALS enviornment is set.
	adminService, err := admin.NewService(ctx)
	if err != nil {
		return err
	}

	m, err := id.Get()
	if err != nil {
		return err
	}
	if m.PrimaryEmail == "" {
		return nil
	}

	slog.Info("font", "message", "adding user", "email", m.PrimaryEmail)
	if _, err := adminService.Members.Insert("font@saint-luke.net", &admin.Member{Email: m.PrimaryEmail}).Do(); err != nil {
		slog.Error("font", "error", err.Error())
		return err
	}
	return nil
}
