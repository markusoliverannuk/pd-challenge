package sqlstore

import (
	"challenge/internal/store"
	"database/sql"
)

type Store struct {
	Db             *sql.DB
	gistRepository *GistsRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		Db: db,
	}
}

func (s *Store) Gists() store.GistsRepository {
	if s.gistRepository != nil {
		return s.gistRepository
	}

	s.gistRepository = &GistsRepository{
		store: s,
	}

	return s.gistRepository
}
