package store

type Store interface {
	Gists() GistsRepository
}
