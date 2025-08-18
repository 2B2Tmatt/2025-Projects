package session

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"
)

type Store struct {
	DB       *sql.DB
	Lifetime time.Duration
	Cookie   string
	Secure   bool
}

func randSID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *Store) Create(ctx context.Context, userID int64) (string, time.Time, error) {
	sid, err := randSID()
	if err != nil {
		return "", time.Time{}, err
	}
	exp := time.Now().Add(s.Lifetime)
	_, err = s.DB.ExecContext(ctx,
		`INSERT into sessions (sid, user_id, expires_at) VALUES ($1,$2,$3)`,
		sid, userID, exp)
	return sid, exp, err
}
