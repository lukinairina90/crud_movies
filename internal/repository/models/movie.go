package models

import "github.com/lukinairina90/crud_movies/internal/domain"

type Movie struct {
	ID             int64  `db:"id"`
	Name           string `db:"name"`
	Description    string `db:"description"`
	ProductionYear int    `db:"production_year"`
	Poster         string `db:"poster"`
	Actors         string `db:"actors"`
	Genre          string `db:"genre"`
}

func (m Movie) ToDomain() domain.Movie {
	return domain.Movie{
		ID:             m.ID,
		Name:           m.Name,
		Description:    m.Description,
		ProductionYear: m.ProductionYear,
		Poster:         m.Poster,
		Actors:         m.Actors,
		Genre:          m.Genre,
	}
}
