package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	driver    = "sqlite3"
	uploadDir = "./static/assets/pics/"
)

type store struct {
	db *sql.DB
}

func New(dsn string) (*store, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	s := &store{db}
	if err = s.init(); err != nil {
		return nil, fmt.Errorf("failed to init store: %w", err)
	}

	return s, nil
}

func (s *store) Close() error {
	return s.db.Close()
}

const initSQL = `
	BEGIN;

	CREATE TABLE IF NOT EXISTS posts (
		id         TEXT     NOT NULL PRIMARY KEY,
		caption    TEXT     NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NULL
	);

	CREATE TABLE IF NOT EXISTS images (
		filename   TEXT     NOT NULL PRIMARY KEY,
		post_id    TEXT     NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NULL
	);

	CREATE INDEX IF NOT EXISTS images_post_id_index ON images(post_id);

	END;
`

func (s *store) init() error {
	_, err := s.db.Exec(initSQL)
	if err != nil {
		return fmt.Errorf("failed to exec init statement: %w", err)
	}

	// make sure upload dir exists before handling requests
	err = os.MkdirAll(uploadDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to make upload dir: %w", err)
	}

	return nil
}
