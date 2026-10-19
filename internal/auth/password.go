package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"

	"typing/internal/config"
)

var (
	ErrInvalidHash         = errors.New("auth: invalid password hash format")
	ErrIncompatibleVersion = errors.New("auth: incompatible argon2 version")
	ErrPasswordMismatch    = errors.New("auth: password mismatch")
)

type Hasher struct {
	cfg config.Argon2Config
}

func NewHasher(cfg config.Argon2Config) *Hasher {
	return &Hasher{cfg: cfg}
}

func (h *Hasher) Hash(password string) (string, error) {
	salt := make([]byte, h.cfg.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("auth: read salt: %w", err)
	}

	key := argon2.IDKey(
		[]byte(password),
		salt,
		h.cfg.Iterations,
		h.cfg.MemoryKiB,
		h.cfg.Parallelism,
		h.cfg.KeyLen,
	)

	b64 := base64.RawStdEncoding
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		h.cfg.MemoryKiB,
		h.cfg.Iterations,
		h.cfg.Parallelism,
		b64.EncodeToString(salt),
		b64.EncodeToString(key),
	), nil
}

func (h *Hasher) Verify(password, encoded string) error {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return ErrInvalidHash
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return ErrInvalidHash
	}
	if version != argon2.Version {
		return ErrIncompatibleVersion
	}

	var memory, iterations uint32
	var parallelism uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); err != nil {
		return ErrInvalidHash
	}

	b64 := base64.RawStdEncoding
	salt, err := b64.DecodeString(parts[4])
	if err != nil {
		return ErrInvalidHash
	}
	want, err := b64.DecodeString(parts[5])
	if err != nil {
		return ErrInvalidHash
	}

	got := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(want)))
	if subtle.ConstantTimeCompare(got, want) != 1 {
		return ErrPasswordMismatch
	}
	return nil
}
