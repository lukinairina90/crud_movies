package domain

type ListMovie []Movie

type Movie struct {
	ID             int64  `json:"id" db:"id" swaggerignore:"true"`
	Name           string `json:"name" db:"name"`
	Description    string `json:"description" db:"description"`
	ProductionYear int    `json:"production_year" db:"production_year"`
	Poster         string `json:"poster" db:"poster"`
	Actors         string `json:"actors" db:"actors"`
	Genre          string `json:"genre" db:"genre"`
}
