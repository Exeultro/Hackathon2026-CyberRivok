package user

import (
	"context"
	"strings"
	"unicode/utf8"

	"typing/internal/settings"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetMe(ctx context.Context, id string) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

type UpdateInput struct {
	Nickname *string
	Settings *settings.Patch
}

func (s *Service) UpdateMe(ctx context.Context, id string, in UpdateInput) (*User, error) {
	if in.Nickname != nil {
		n := strings.TrimSpace(*in.Nickname)
		in.Nickname = &n
	}
	return s.repo.Update(ctx, id, in.Nickname, in.Settings)
}

func ValidateUpdate(in UpdateInput) map[string]any {
	details := map[string]any{}

	if in.Nickname != nil {
		l := utf8.RuneCountInString(*in.Nickname)
		if l < 2 || l > 24 {
			details["nickname"] = "length must be 2..24"
		}
	}

	if p := in.Settings; p != nil {
		if p.Lang != nil && *p.Lang != "ru" && *p.Lang != "en" {
			details["settings.lang"] = "must be one of: ru, en"
		}
		if p.DisplayUnit != nil && *p.DisplayUnit != "word" && *p.DisplayUnit != "char" {
			details["settings.display_unit"] = "must be one of: word, char"
		}
		if p.Theme != nil && utf8.RuneCountInString(*p.Theme) > 64 {
			details["settings.theme"] = "max length 64"
		}
	}

	return details
}
