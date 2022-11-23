package service

import (
	"context"

	"github.com/lukinairina90/crud_movies/internal/domain"
)

type MoviesRepository interface {
	List(ctx context.Context) (domain.ListMovie, error)
	Get(ctx context.Context, id int) (domain.Movie, error)
	Create(ctx context.Context, movie domain.Movie) (domain.Movie, error)
	Update(ctx context.Context, id int, movie domain.Movie) (domain.Movie, error)
	Delete(ctx context.Context, id int) error
}

type Movie struct {
	movieRepository MoviesRepository
}

func NewMovie(movieRepository MoviesRepository) *Movie {
	return &Movie{movieRepository: movieRepository}
}

func (m Movie) List(ctx context.Context) (domain.ListMovie, error) {
	return m.movieRepository.List(ctx)
}

func (m Movie) Get(ctx context.Context, id int) (domain.Movie, error) {
	return m.movieRepository.Get(ctx, id)
}

func (m Movie) Create(ctx context.Context, movie domain.Movie) (domain.Movie, error) {
	return m.movieRepository.Create(ctx, movie)
}

func (m Movie) Update(ctx context.Context, id int, movie domain.Movie) (domain.Movie, error) {
	return m.movieRepository.Update(ctx, id, movie)
}

func (m Movie) Delete(ctx context.Context, id int) error {
	return m.movieRepository.Delete(ctx, id)
}
