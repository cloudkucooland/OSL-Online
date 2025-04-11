package main

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"google.golang.org/api/admin/directory/v1"

	"github.com/cloudkucooland/OSL-Online/model"
)

func main() {
	ctx := context.Background()

	dbpath := os.Getenv("OO_DB")
	if dbpath == "" {
		panic("OO_DB enviornment var not set. e.g. oo:password@unix(/var/lib/mysql/mysql.sock)/oo")
	}

	if err := model.Connect(ctx, dbpath); err != nil {
		slog.Error("startup", "message", "Error connecting to database", "error", err.Error())
		panic(err)
	}

	gac := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if gac == "" {
		panic("GOOGLE_APPLICATION_CREDENTIALS enviornment var not set.")
	}

	adminService, err := admin.NewService(ctx)
	if err != nil {
		panic(err.Error())
	}

	call := adminService.Members.List("font@saint-luke.net")
	known := make(map[string]bool, 0)
	err = call.Pages(ctx, func(members *admin.Members) error {
		for _, m := range members.Members {
			e := strings.ToLower(m.Email)
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
			} else {
				known[e] = true
			}
		}
		return nil
	})
	if err != nil {
		panic(err.Error())
	}

	toadd, err := model.FontEmailedDirect()
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
}

func checkFont(email string) (bool, error) {
	found, err := model.SearchEmail(email, true)
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
	if strings.ToLower(member.Newsletter) == "none" {
		return false, nil
	}

	return true, nil
}
