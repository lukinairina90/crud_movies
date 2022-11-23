package domain

type ListMovie []Movie

type Movie struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	ProductionYear int    `json:"production_year"`
	Poster         string `json:"poster"`
	Actors         string `json:"actors"`
	Genre          string `json:"genre"`
}
