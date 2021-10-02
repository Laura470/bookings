package dbrepo

import (
	"database/sql"

	"github.com/Laura470/bookings/internal/config"
	"github.com/Laura470/bookings/internal/repository"
)

//it is the repository itself
type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB //chiamo una funzione che in realtà non ha dietro un db
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
