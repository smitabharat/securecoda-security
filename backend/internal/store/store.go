package store

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct{ db *sql.DB }

type Alert struct {
	ID        int64
	DocID     string
	Title     string
	Type      string
	Detail    string
	CreatedAt time.Time
}

func NewSQLiteStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) migrate() error {
	sqlStmt := `CREATE TABLE IF NOT EXISTS alerts (
id INTEGER PRIMARY KEY AUTOINCREMENT,
doc_id TEXT,
title TEXT,
type TEXT,
detail TEXT,
created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`
	_, err := s.db.Exec(sqlStmt)
	return err
}

func (s *Store) CreateAlert(a *Alert) error {
	_, err := s.db.Exec(`INSERT INTO alerts(doc_id,title,type,detail) VALUES(?,?,?,?)`, a.DocID, a.Title, a.Type, a.Detail)
	return err
}

func (s *Store) ListAlerts() ([]*Alert, error) {
	rows, err := s.db.Query(`SELECT id,doc_id,title,type,detail,created_at FROM alerts ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*Alert
	for rows.Next() {
		a := &Alert{}
		var ts string
		if err := rows.Scan(&a.ID, &a.DocID, &a.Title, &a.Type, &a.Detail, &ts); err != nil {
			return nil, err
		}
		a.CreatedAt, _ = time.Parse(time.RFC3339, ts)
		out = append(out, a)
	}
	return out, nil
}
