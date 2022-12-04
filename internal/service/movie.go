package service

import (
	"context"
	"fmt"
	"time"

	"github.com/lukinairina90/crud_movies/internal/domain"
	"github.com/lukinairina90/in_memory_cache/generic_cache"
)

const movieKeyPattern = "movie:%d"

type MoviesRepository interface {
	List(ctx context.Context) (domain.ListMovie, error)
	Get(ctx context.Context, id int) (domain.Movie, error)
	Create(ctx context.Context, movie domain.Movie) (domain.Movie, error)
	Update(ctx context.Context, id int, movie domain.Movie) (domain.Movie, error)
	Delete(ctx context.Context, id int) error
}

type Cacher[K comparable, V any] interface {
	Set(key K, value V, ttl time.Duration) error
	Get(key K) (V, error)
	Delete(key K) error
}

type Movie struct {
	movieRepository MoviesRepository
	cache           Cacher[string, domain.Movie]
}

func NewMovie(movieRepository MoviesRepository, cacher Cacher[string, domain.Movie]) *Movie {
	return &Movie{
		movieRepository: movieRepository,
		cache:           cacher,
	}
}

func (m Movie) List(ctx context.Context) (domain.ListMovie, error) {
	return m.movieRepository.List(ctx)
}

func (m Movie) Get(ctx context.Context, id int) (domain.Movie, error) {
	res, err := m.cache.Get(fmt.Sprintf(movieKeyPattern, id))
	if err != nil && err != generic_cache.ErrKeyNotFound {
		return domain.Movie{}, err
	}

	if err == generic_cache.ErrKeyNotFound {
		res, err = m.movieRepository.Get(ctx, id)
		if err != nil {
			return domain.Movie{}, err
		}

		if err := m.cache.Set(fmt.Sprintf(movieKeyPattern, id), res, time.Minute*10); err != nil {
			return domain.Movie{}, err
		}
	}

	return res, nil
}

func (m Movie) Create(ctx context.Context, movie domain.Movie) (domain.Movie, error) {
	return m.movieRepository.Create(ctx, movie)
}

func (m Movie) Update(ctx context.Context, id int, movie domain.Movie) (domain.Movie, error) {
	movie, err := m.movieRepository.Update(ctx, id, movie)
	if err != nil {
		return domain.Movie{}, err
	}

	if err := m.cache.Delete(fmt.Sprintf(movieKeyPattern, id)); err != nil {
		return domain.Movie{}, err
	}

	return movie, err
}

func (m Movie) Delete(ctx context.Context, id int) error {
	if err := m.movieRepository.Delete(ctx, id); err != nil {
		return err
	}

	return m.cache.Delete(fmt.Sprintf(movieKeyPattern, id))
}
