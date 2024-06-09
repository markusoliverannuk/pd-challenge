package models

type Gist struct {
	Username    string   `db:"username" json:"username"`
	Description string   `db:"description" json:"description"`
	Files       []string `json:"files"`
}

type File struct {
	Username string `db:"username"`
	Path     string `db:"path"`
}
