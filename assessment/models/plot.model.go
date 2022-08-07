package models

type Plot struct {
	IMDB_ID string `json:"imdbID"`
	Title   string `json:"Title"`
	Plot    string `json:"Plot"`
}
