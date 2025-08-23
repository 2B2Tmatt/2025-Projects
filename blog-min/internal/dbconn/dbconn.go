package dbconn

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func OpenDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL1")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is required (put it in .env for dev or pass via Compose)")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	for ctx.Err() == nil {
		if err := db.PingContext(ctx); err == nil {
			log.Printf("database connection established")
			return db, nil
		}
		time.Sleep(time.Second)
	}
	return nil, fmt.Errorf("database not ready")
}

func InitSchema(db *sql.DB) error {
	const schema = `
BEGIN;

CREATE TABLE IF NOT EXISTS users (
  id             SERIAL PRIMARY KEY,
  display_name   VARCHAR(100) NOT NULL UNIQUE,
  email          VARCHAR(255) NOT NULL UNIQUE,
  password_hash  TEXT NOT NULL,
  created_at     TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS posts (
  id         BIGSERIAL PRIMARY KEY,
  user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title      TEXT NOT NULL,
  body       TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sessions (
  sid        TEXT PRIMARY KEY,
  uid        INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_posts_user_id_created_at ON posts (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sessions_uid             ON sessions (uid);       
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at      ON sessions (expires_at);

COMMIT;
`
	_, err := db.Exec(schema)
	return err
}
