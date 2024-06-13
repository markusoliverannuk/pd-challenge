package sqlstore

import (
	"challenge/internal/models"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"gitlab.com/0x4149/logz"
)

type GistsRepository struct {
	store *Store
}

func (g *GistsRepository) CreateGist(gist models.Gist) (*models.Gist, error) {
	var username, description string
	query := `SELECT username, description FROM gists WHERE username = ? AND description = ?`
	err := g.store.Db.QueryRow(query, gist.Username, gist.Description).Scan(&username, &description)

	if err != nil {
		if err == sql.ErrNoRows {

			query = `INSERT INTO gists (id, username, description, seen) VALUES (?, ?, ?, ?)`

			gist.Id = uuid.New().String()
			_, err = g.store.Db.Exec(query, gist.Id, gist.Username, gist.Description, 0)
			if err != nil {
				return nil, err
			}

			return &gist, nil
		}
		return nil, err // return the error if it's not ErrNoRows
	}

	return nil, errors.New("nothing new to add")
}

func (g *GistsRepository) ChangeVisibility(gist models.Gist) error {
	logz.Info("Gist author:", gist.Username, "Gist:", gist.Description)
	query := `UPDATE gists SET seen = 1 WHERE username = ? AND description = ?`
	rows, err := g.store.Db.Exec(query, gist.Username, gist.Description)
	if err != nil {
		return err
	}
	r_affected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if r_affected == 0 {
		return errors.New("nothing on db changed")
	}

	return nil
}

func (g *GistsRepository) GetAllFiles(id string) ([]string, error) {
	query := `SELECT path FROM files WHERE id = ?`
	rows, err := g.store.Db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []string
	for rows.Next() {
		var file string
		if err := rows.Scan(&file); err != nil {
			return nil, err
		}

		files = append(files, file)
	}
	return files, nil
}

func (g *GistsRepository) CreateFile(file models.File) (*models.File, error) {
	query := `INSERT OR REPLACE INTO files (id, username, path) VALUES (?, ?, ?)`

	_, err := g.store.Db.Exec(query, file.Id, file.Username, file.Path)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (g *GistsRepository) GetUsersOld(username string) ([]models.Gist, error) {
	query := `SELECT id, username, description FROM gists WHERE username = ? AND seen = 1`
	rows, err := g.store.Db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gists []models.Gist
	// iterating over the rows
	for rows.Next() {
		var gist models.Gist
		if err := rows.Scan(&gist.Id, &gist.Username, &gist.Description); err != nil {
			return nil, err
		}

		allFiles, err := g.GetAllFiles(gist.Id)
		if err != nil {
			logz.Error(err)
		}
		gist.Files = allFiles

		gists = append(gists, gist)

	}

	return gists, nil
}

func (g *GistsRepository) GetUsersNew(username string) ([]models.Gist, error) {
	query := `SELECT id, username, description FROM gists WHERE username = ? AND seen = 0`
	rows, err := g.store.Db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gists []models.Gist
	// iterating over the rows
	for rows.Next() {
		var gist models.Gist
		if err := rows.Scan(&gist.Id, &gist.Username, &gist.Description); err != nil {
			return nil, err
		}

		allFiles, err := g.GetAllFiles(gist.Id)
		if err != nil {
			logz.Error(err)
		}
		gist.Files = allFiles

		gists = append(gists, gist)
	}

	for _, gistToChange := range gists {
		g.ChangeVisibility(gistToChange)
	}

	return gists, nil
}

func (g *GistsRepository) GetUniqueNames() ([]string, error) {
	query := "SELECT DISTINCT username FROM gists"
	rows, err := g.store.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uniqueNames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		uniqueNames = append(uniqueNames, name)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return uniqueNames, nil
}
