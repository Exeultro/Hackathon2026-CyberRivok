package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Type string

const (
	TypeAccess  Type = "access"
	TypeRefresh Type = "refresh"
)

var ErrInvalid = errors.New("token: invalid")

type Claims struct {
	Typ Type `json:"typ"`
	jwt.RegisteredClaims
}

func Parse(secret []byte, issuer, raw string, expected Type) (*Claims, error) {
	c := &Claims{}
	parsed, err := jwt.ParseWithClaims(raw, c, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalid
		}
		return secret, nil
	},
		jwt.WithIssuer(issuer),
		jwt.WithValidMethods([]string{"HS256"}),
	)
	if err != nil || !parsed.Valid {
		return nil, ErrInvalid
	}
	if c.Typ != expected || c.Subject == "" {
		return nil, ErrInvalid
	}
	return c, nil
}

func NewAccess(subject, issuer string, now time.Time, ttl time.Duration) Claims {
	return Claims{
		Typ: TypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
}

func NewRefresh(subject, issuer, jti string, now time.Time, ttl time.Duration) Claims {
	return Claims{
		Typ: TypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			ID:        jti,
		},
	}
}
