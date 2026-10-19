package texts

import (
	"context"
	"errors"
)

var ErrValidation = errors.New("texts: validation")

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

type Request struct {
	Mode         string
	Lang         string
	Difficulty   int
	Count        int
	CodeLanguage string
	UserID       string
}

func (s *Service) Get(ctx context.Context, req Request) ([]Item, error) {
	if err := validate(&req); err != nil {
		return nil, err
	}

	switch req.Mode {
	case "words", "symbols":
		return s.generate(req), nil
	case "sentences", "code":
		items, err := s.repo.Pick(ctx, req.Mode, req.Lang, req.Difficulty, req.Count, req.CodeLanguage)
		if err != nil {
			return nil, err
		}
		if len(items) == 0 && (req.Mode == "sentences" || req.Mode == "code") {
			for diff := 1; diff <= 5; diff++ {
				if diff == req.Difficulty {
					continue
				}
				items, err = s.repo.Pick(ctx, req.Mode, req.Lang, diff, req.Count, req.CodeLanguage)
				if err != nil {
					return nil, err
				}
				if len(items) > 0 {
					break
				}
			}
		}
		return items, nil
	}
	return nil, ErrValidation
}

func (s *Service) generate(req Request) []Item {
	switch req.Mode {
	case "words":
		return GenerateWords(req.Lang, req.Difficulty, req.Count)
	case "symbols":
		return GenerateSymbols(req.Lang, req.Difficulty, req.Count)
	}
	return nil
}

func validate(r *Request) error {
	switch r.Mode {
	case "words", "sentences", "code", "symbols":
	default:
		return ErrValidation
	}
	if r.Lang != "ru" && r.Lang != "en" {
		return ErrValidation
	}
	if r.Difficulty < 1 || r.Difficulty > 5 {
		return ErrValidation
	}
	if r.Count < 1 || r.Count > 50 {
		return ErrValidation
	}
	if r.Mode == "code" {
		switch r.CodeLanguage {
		case "python", "javascript", "html", "css":
		default:
			return ErrValidation
		}
	}
	return nil
}
