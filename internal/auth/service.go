package auth

import (
	"context"
	"errors"
	"strings"

	"typing/internal/config"
	"typing/internal/token"
)

var (
	ErrInvalidCredentials = errors.New("auth: invalid credentials")
)

type Service struct {
	repo   *Repo
	hasher *Hasher
	tokens *TokenManager
}

func NewService(repo *Repo, cfg *config.Config) *Service {
	return &Service{
		repo:   repo,
		hasher: NewHasher(cfg.Argon2),
		tokens: NewTokenManager(cfg.JWT),
	}
}

type AuthResult struct {
	User         *User
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

func (s *Service) Register(ctx context.Context, email, password, nickname string) (*AuthResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	nickname = strings.TrimSpace(nickname)

	hash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	u, err := s.repo.CreateUser(ctx, email, nickname, hash)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailTaken):
			return nil, ErrEmailTaken
		case errors.Is(err, ErrNicknameTaken):
			return nil, ErrNicknameTaken
		default:
			return nil, err
		}
	}
	return s.issue(ctx, u)
}

func (s *Service) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	u, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if err := s.hasher.Verify(password, u.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}
	return s.issue(ctx, u)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*AuthResult, error) {
	if _, err := s.tokens.Parse(refreshToken, token.TypeRefresh); err != nil {
		return nil, ErrInvalidCredentials
	}
	hash := HashRefresh(refreshToken)
	rt, err := s.repo.GetActiveRefresh(ctx, hash)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	u, err := s.repo.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if err := s.repo.RevokeRefresh(ctx, hash); err != nil {
		return nil, err
	}
	return s.issue(ctx, u)
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	if _, err := s.tokens.Parse(refreshToken, token.TypeRefresh); err != nil {
		return nil
	}
	return s.repo.RevokeRefresh(ctx, HashRefresh(refreshToken))
}

func (s *Service) issue(ctx context.Context, u *User) (*AuthResult, error) {
	pair, err := s.tokens.IssuePair(u.ID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.InsertRefreshToken(ctx, u.ID, HashRefresh(pair.RefreshToken), pair.RefreshExpiresAt); err != nil {
		return nil, err
	}
	return &AuthResult{
		User:         u,
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    pair.AccessExpiresIn,
	}, nil
}
