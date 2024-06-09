package sqlstore

import (
	"database/sql"
	"challenge/internal/store"
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
