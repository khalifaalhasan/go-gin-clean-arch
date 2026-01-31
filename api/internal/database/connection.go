package database

import (
	"database/sql"

	_ "github.com/lib/pq" 
)

// NewConnection membuka koneksi ke database PostgreSQL
func NewConnection(dbSource string) (*sql.DB, error) {
	// Buka koneksi menggunakan driver "postgres"
	db, err := sql.Open("postgres", dbSource)
	if err != nil {
		return nil, err
	}

	// Cek apakah koneksi benar-benar hidup (Ping)
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}