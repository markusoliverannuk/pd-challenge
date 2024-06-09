package sqlstore

import "challenge/internal/models"

type GistsRepository struct {
	store *Store
}

func (g *GistsRepository) CreateGist(gist models.Gist) (*models.Gist, error) {
	query := `INSERT OR REPLACE INTO gists (username, description) VALUES (?, ?)`

	_, err := g.store.Db.Exec(query, gist.Username, gist.Description)
	if err != nil {
		return nil, err
	}

	return &gist, nil
}

func (g *GistsRepository) CreateFile(file models.File) (*models.File, error) {
	query := `INSERT OR REPLACE INTO files (username, path) VALUES (?, ?)`

	_, err := g.store.Db.Exec(query, file.Username, file.Path)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (g *GistsRepository) GetUsers(username string) ([]models.Gist, error) {
	query := `SELECT g.username, g.description, f.path FROM gists g LEFT JOIN files f ON g.username = f.username WHERE f.username = ?`
	rows, err := g.store.Db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map to store gists and their files
	gistsMap := make(map[string]*models.Gist)

	// Iterate over the rows
	for rows.Next() {
		var username, description, path string
		if err := rows.Scan(&username, &description, &path); err != nil {
			return nil, err
		}

		// Check if the gist already exists in the map
		if _, exists := gistsMap[username]; !exists {
			gistsMap[username] = &models.Gist{
				Username:    username,
				Description: description,
				Files:       []string{},
			}
		}

		// Add the file to the gist
		if path != "" { // Ignore rows without a file
			gistsMap[username].Files = append(gistsMap[username].Files, path)
		}
	}

	// Convert the map to a slice
	gists := []models.Gist{}
	for _, gist := range gistsMap {
		gists = append(gists, *gist)
	}

	return gists, nil
}
