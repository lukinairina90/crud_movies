package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lukinairina90/crud_movies/internal/domain"
	"github.com/lukinairina90/crud_movies/internal/repository/models"
)

type Movie struct {
	db *sqlx.DB
}

func NewMovie(db *sqlx.DB) *Movie {
	return &Movie{db: db}
}

func (m Movie) List(ctx context.Context) (domain.ListMovie, error) {
	var list []models.Movie
	if err := m.db.SelectContext(ctx, &list, "SELECT * FROM movie"); err != nil {
		return nil, err
	}

	dlist := make(domain.ListMovie, 0, len(list))
	for _, movie := range list {
		dlist = append(dlist, movie.ToDomain())
	}

	return dlist, nil
}

func (m Movie) Get(ctx context.Context, id int) (domain.Movie, error) {
	var movie models.Movie
	if err := m.db.GetContext(ctx, &movie, "SELECT * FROM  movie WHERE id=$1", id); err != nil {
		return domain.Movie{}, err
	}

	return movie.ToDomain(), nil
}

func (m Movie) Create(ctx context.Context, movie domain.Movie) (domain.Movie, error) {
	mMovie := models.Movie{
		Name:           movie.Name,
		Description:    movie.Description,
		ProductionYear: movie.ProductionYear,
		Poster:         movie.Poster,
		Actors:         movie.Actors,
		Genre:          movie.Genre,
	}

	if err := m.db.QueryRowxContext(ctx, "INSERT INTO movie (name, description, production_year, genre, actors, poster) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *", mMovie.Name, mMovie.Description, mMovie.ProductionYear, mMovie.Genre, mMovie.Actors, mMovie.Poster).StructScan(&mMovie); err != nil {
		return domain.Movie{}, err
	}

	return mMovie.ToDomain(), nil
}

func (m Movie) Update(ctx context.Context, id int, movie domain.Movie) (domain.Movie, error) {
	mMovie := models.Movie{
		Name:           movie.Name,
		Description:    movie.Description,
		ProductionYear: movie.ProductionYear,
		Poster:         movie.Poster,
		Actors:         movie.Actors,
		Genre:          movie.Genre,
	}

	if err := m.db.QueryRowxContext(ctx, "UPDATE movie SET name=$1, description=$2, production_year=$3, genre=$4, actors=$5, poster=$6 WHERE id=$7 RETURNING *", mMovie.Name, mMovie.Description, mMovie.ProductionYear, mMovie.Genre, mMovie.Actors, mMovie.Poster, id).StructScan(&mMovie); err != nil {
		return domain.Movie{}, err
	}

	return mMovie.ToDomain(), nil
}

func (m Movie) Delete(ctx context.Context, id int) error {
	if _, err := m.db.ExecContext(ctx, "DELETE FROM movie WHERE id=$1", id); err != nil {
		return err
	}
	return nil
}

//const movieKeyPattern = "movie:%d"
//
//type CachedMovie struct {
//	repo  *Movie
//	cache *generic_cache.Cache[string, domain.Movie]
//}
//
//func NewCachedMovie(repo *Movie, cache *generic_cache.Cache[string, domain.Movie]) *CachedMovie {
//	return &CachedMovie{
//		repo:  repo,
//		cache: cache,
//	}
//}
//
//func (c CachedMovie) List(ctx context.Context) (domain.ListMovie, error) {
//	return c.repo.List(ctx)
//}
//
//func (c CachedMovie) Get(ctx context.Context, id int) (domain.Movie, error) {
//	res, err := c.cache.Get(fmt.Sprintf(movieKeyPattern, id))
//	if err != nil && err != generic_cache.ErrKeyNotFound {
//		return domain.Movie{}, err
//	}
//
//	if err == generic_cache.ErrKeyNotFound {
//		res, err = c.repo.Get(ctx, id)
//		if err != nil {
//			return domain.Movie{}, err
//		}
//
//		if err := c.cache.Set(fmt.Sprintf(movieKeyPattern, id), res, time.Minute); err != nil {
//			return domain.Movie{}, err
//		}
//	}
//
//	return res, nil
//}
//
//func (c CachedMovie) Create(ctx context.Context, movie domain.Movie) (domain.Movie, error) {
//	return c.repo.Create(ctx, movie)
//}
//
//func (c CachedMovie) Update(ctx context.Context, id int, movie domain.Movie) (domain.Movie, error) {
//	return c.repo.Update(ctx, id, movie)
//}
//
//func (c CachedMovie) Delete(ctx context.Context, id int) error {
//	return c.repo.Delete(ctx, id)
//}
