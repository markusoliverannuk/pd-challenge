package store

import "challenge/internal/models"

type GistsRepository interface {
	CreateGist(gist models.Gist) (*models.Gist, error)
	CreateFile(file models.File) (*models.File, error)
	GetUsersOld(username string) ([]models.Gist, error)
	GetUsersNew(username string) ([]models.Gist, error)
	GetUniqueNames() ([]string, error) // New method to fetch unique names
}
