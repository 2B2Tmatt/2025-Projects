package dbconn

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq" // or pgx stdlib if you use it
)

func OpenDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	log.Printf("DATABASE_URL = %q", dsn) // should include host=db

	if strings.TrimSpace(dsn) == "" {
		// Fallback that works in Compose
		dsn = "host=db port=5432 user=postgres password=Thefakestpoggers123! dbname=blog_db01 sslmode=disable"
		log.Printf("DATABASE_URL missing; using Compose fallback DSN")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// Retry while Postgres starts up
	for i := 0; i < 30; i++ {
		if err := db.Ping(); err == nil {
			log.Printf("database connection established")
			return db, nil
		}
		time.Sleep(1 * time.Second)
	}
	return nil, fmt.Errorf("database not ready after retries")
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

-- helpful indexes
CREATE INDEX IF NOT EXISTS idx_posts_user_id_created_at ON posts (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sessions_uid             ON sessions (uid);        -- 👈 fixed
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at      ON sessions (expires_at);

COMMIT;
`
	_, err := db.Exec(schema)
	return err
}
