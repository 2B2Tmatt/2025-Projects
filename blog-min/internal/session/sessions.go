package sessions

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"time"
)

type Session struct {
	Sid        string
	Uid        int64
	Created_at time.Time
	Expires_at time.Time
}

func generateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v\n", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func CreateSession(w http.ResponseWriter, db *sql.DB, uid int64, duration time.Duration) (*Session, error) {
	sid := generateToken(32)
	_, err := db.Exec(`
	INSERT INTO sessions (sid, uid, expires_at) VALUES ($1, $2, $3)`,
		sid, uid, (time.Now().Add(duration)))
	if err != nil {
		return nil, err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sid,
		Path:     "/",
		Expires:  time.Now().Add(duration),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	session := Session{sid, uid, time.Now(), time.Now().Add(duration)}
	return &session, nil
}

func EndSession(w http.ResponseWriter, db *sql.DB, sid string) error {
	_, err := db.Exec(`DELETE FROM SESSIONS WHERE sid=$1`, sid)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

func CheckSession(db *sql.DB, sid string) (int64, error) {
	var uid int64
	err := db.QueryRow(`
        SELECT uid FROM sessions
        WHERE sid = $1 AND expires_at > now()
    `, sid).Scan(&uid)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil // no session
	}
	if err != nil {
		return 0, err // real DB error
	}
	return uid, nil // valid session
}
