package model

import (
	"context"
	"log/slog"
	"strings"

	"google.golang.org/api/admin/directory/v1"
)

func (id ChapterID) ChapterSync(ctx context.Context) error {
	// assumes that GOOGLE_APPLICATION_CREDENTIALS enviornment is set

	adminService, err := admin.NewService(ctx)
	if err != nil {
		return err
	}

	addr, err := id.emailaddress()
	if err != nil || addr == "" {
		return nil
	}

	call := adminService.Members.List(addr)
	known := make(map[string]bool, 0)
	_ = call.Pages(ctx, func(members *admin.Members) error {
		for _, m := range members.Members {
			e := strings.ToLower(m.Email)

			if strings.Contains(e, "@saint-luke.net") { // leave these alone
				known[e] = true
				continue
			}

			ok, err := checkChapter(ctx, e, id)
			if err != nil {
				slog.Error("chapter sync", "error", err.Error())
				continue
			}
			if !ok {
				slog.Info("chapter sync", "message", "removing", "user", e)
				if err := adminService.Members.Delete(addr, e).Do(); err != nil {
					slog.Error("chapter", "error", err.Error())
					continue
				}
			}
			known[e] = true
		}
		return nil
	})

	c, err := id.Load()
	if err != nil {
		return err
	}
	toadd, err := c.Members()
	if err != nil {
		return err
	}
	for _, m := range toadd {
		e := strings.ToLower(m.PrimaryEmail)
		if e == "" {
			continue
		}
		if _, ok := known[e]; !ok {
			slog.Info("chapter sync", "message", "adding user", "email", e)
			if _, err := adminService.Members.Insert(addr, &admin.Member{Email: e}).Do(); err != nil {
				slog.Error("chapter sync", "error", err.Error())
				continue
			}
		}
	}
	return nil
}

func checkChapter(ctx context.Context, email string, chapterid ChapterID) (bool, error) {
	found, err := SearchEmail(ctx, email, true)
	if err != nil {
		return false, err
	}
	if len(found) == 0 { // should be != 1 but Br Dan and Sr. Mary-O share an address
		return false, nil
	}
	/* member, err := found[0].ID.Get()
	if err != nil {
		return false, err
	} */

	// XXX get chapter membership

	return true, nil
}

func (id MemberID) UnsubscribeChapter(ctx context.Context, chapterid ChapterID) error {
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

	addr, err := chapterid.emailaddress()
	if err != nil || addr == "" {
		return nil
	}

	slog.Info("chapter email", "message", "removing", "user", m.PrimaryEmail)
	if err := adminService.Members.Delete(addr, m.PrimaryEmail).Do(); err != nil {
		slog.Error("chapter email", "error", err.Error())
		return err
	}
	return nil
}

func (id MemberID) SubscribeChapter(ctx context.Context, chapterid ChapterID) error {
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

	addr, err := chapterid.emailaddress()
	if err != nil || addr == "" {
		return nil
	}

	slog.Info("chapter email", "message", "adding user", "email", m.PrimaryEmail)
	if _, err := adminService.Members.Insert(addr, &admin.Member{Email: m.PrimaryEmail}).Do(); err != nil {
		slog.Error("chapter email", "error", err.Error())
		return err
	}
	return nil
}

func (chapter ChapterID) emailaddress() (string, error) {
	var e string
	err := db.QueryRow("SELECT email FROM chapters WHERE id = ?", chapter).Scan(&e)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}
	return e, nil
}

func (id MemberID) UnsubscribeAllChapters(ctx context.Context) error {
	if _, err := db.ExecContext(ctx, "DELETE FROM `chaptermembers` WHERE `member` = ?", id); err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func (id MemberID) SubscribeAllChapters(ctx context.Context) error {
	chapters, err := id.GetChapters()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, c := range chapters {
		if err := id.SubscribeChapter(ctx, ChapterID(c)); err != nil {
			return err
		}
	}
	return nil
}
