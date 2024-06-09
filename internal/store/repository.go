package store

import "challenge/internal/models"

type GistsRepository interface {
	CreateGist(gist models.Gist) (*models.Gist, error)
	CreateFile(file models.File) (*models.File, error)
	GetUsers(username string) ([]models.Gist, error)
}
