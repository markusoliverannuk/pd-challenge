package models

type Gist struct {
	Id          string   `db:"id" json:"id"`
	Username    string   `db:"username" json:"username"`
	Description string   `db:"description" json:"description"`
	Files       []string `json:"files"`
}

type File struct {
	Id       string `db:"id"`
	Username string `db:"username"`
	Path     string `db:"path"`
}
