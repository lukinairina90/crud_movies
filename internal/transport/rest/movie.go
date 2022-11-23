package rest

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	restErrors "github.com/lukinairina90/crud_movies/pkg/rest/errors"

	"github.com/gin-gonic/gin"
	"github.com/lukinairina90/crud_movies/internal/domain"
)

type MovieService interface {
	List(ctx context.Context) (domain.ListMovie, error)
	Get(ctx context.Context, id int) (domain.Movie, error)
	Create(ctx context.Context, movie domain.Movie) (domain.Movie, error)
	Update(ctx context.Context, id int, movie domain.Movie) (domain.Movie, error)
	Delete(ctx context.Context, id int) error
}

type Movie struct {
	movieService MovieService
}

func NewMovie(movieService MovieService) *Movie {
	return &Movie{movieService: movieService}
}

func (m Movie) List(ctx *gin.Context) {
	movies, err := m.movieService.List(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, restErrors.NewInternalServerErr())
		return
	}

	ctx.JSON(http.StatusOK, movies)
}

func (m Movie) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		fields := map[string]string{"id": "should be an integer"}
		ctx.JSON(http.StatusBadRequest, restErrors.NewBadRequestErr("validation error", fields))
		return
	}

	movie, err := m.movieService.Get(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(http.StatusNotFound, restErrors.NewNotFoundErr("movie not found"))
		default:
			ctx.JSON(http.StatusInternalServerError, restErrors.NewInternalServerErr())
		}

		return
	}

	ctx.JSON(http.StatusOK, movie)
}

func (m Movie) Create(ctx *gin.Context) {
	var movie domain.Movie
	if err := ctx.BindJSON(&movie); err != nil {
		ctx.JSON(http.StatusBadRequest, restErrors.NewBadRequestErr("cannot parse body", nil))
		return
	}

	createdMovie, err := m.movieService.Create(ctx, movie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, restErrors.NewInternalServerErr())
		return
	}

	ctx.JSON(http.StatusCreated, createdMovie)
}

func (m Movie) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		fields := map[string]string{"id": "should be an integer"}
		ctx.JSON(http.StatusBadRequest, restErrors.NewBadRequestErr("validation error", fields))
		return
	}

	var movie domain.Movie
	if err := ctx.BindJSON(&movie); err != nil {
		ctx.JSON(http.StatusBadRequest, restErrors.NewBadRequestErr("cannot parse body", nil))
		return
	}

	updatedMovie, err := m.movieService.Update(ctx, id, movie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, restErrors.NewInternalServerErr())
		return
	}

	ctx.JSON(http.StatusOK, updatedMovie)
}

func (m Movie) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		fields := map[string]string{"id": "should be an integer"}
		ctx.JSON(http.StatusBadRequest, restErrors.NewBadRequestErr("validation error", fields))
		return
	}

	if err := m.movieService.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, restErrors.NewInternalServerErr())
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
