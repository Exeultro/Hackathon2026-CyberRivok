package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"typing/internal/config"
	"typing/internal/token"
)

type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresIn  int
	RefreshExpiresAt time.Time
	RefreshJTI       string
}

type TokenManager struct {
	cfg config.JWTConfig
}

func NewTokenManager(cfg config.JWTConfig) *TokenManager {
	return &TokenManager{cfg: cfg}
}

func (m *TokenManager) IssuePair(userID string) (TokenPair, error) {
	now := time.Now().UTC()

	jti, err := randomID()
	if err != nil {
		return TokenPair{}, err
	}

	access := token.NewAccess(userID, m.cfg.Issuer, now, m.cfg.AccessTTL)
	refresh := token.NewRefresh(userID, m.cfg.Issuer, jti, now, m.cfg.RefreshTTL)

	accessStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, access).SignedString(m.cfg.Secret)
	if err != nil {
		return TokenPair{}, fmt.Errorf("auth: sign access: %w", err)
	}
	refreshStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refresh).SignedString(m.cfg.Secret)
	if err != nil {
		return TokenPair{}, fmt.Errorf("auth: sign refresh: %w", err)
	}

	return TokenPair{
		AccessToken:      accessStr,
		RefreshToken:     refreshStr,
		AccessExpiresIn:  int(m.cfg.AccessTTL.Seconds()),
		RefreshExpiresAt: now.Add(m.cfg.RefreshTTL),
		RefreshJTI:       jti,
	}, nil
}

func (m *TokenManager) Parse(raw string, expected token.Type) (*token.Claims, error) {
	return token.Parse(m.cfg.Secret, m.cfg.Issuer, raw, expected)
}

func HashRefresh(raw string) []byte {
	sum := sha256.Sum256([]byte(raw))
	return sum[:]
}

func randomID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("auth: random id: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
