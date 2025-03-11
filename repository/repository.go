package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Repository struct {
	Db *sql.DB
}

type NewRepositoryOptions struct {
	Dsn string
}

func NewRepository(opts NewRepositoryOptions) *Repository {
	db, err := sql.Open("postgres", opts.Dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic("failed: " + err.Error())
	}

	return &Repository{
		Db: db,
	}
}
