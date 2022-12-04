package rest

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lukinairina90/crud_movies/internal/domain"
)

type Movies interface {
	List(ctx context.Context) (domain.ListMovie, error)
	Get(ctx context.Context, id int) (domain.Movie, error)
	Create(ctx context.Context, movie domain.Movie) (domain.Movie, error)
	Update(ctx context.Context, id int, movie domain.Movie) (domain.Movie, error)
	Delete(ctx context.Context, id int) error
}

func NewMovie(movieService Movies) *Movie {
	return &Movie{movieService: movieService}
}

type Movie struct {
	movieService Movies
}

func (m Movie) InjectRoutes(r *gin.Engine, middlewares ...gin.HandlerFunc) {
	movies := r.Group("/movies").Use(middlewares...)
	{
		movies.GET("/", m.getAllMovies)
		movies.GET("/:id", m.getMovie)
		movies.POST("/", m.createMovie)
		movies.PUT("/:id", m.updateMovie)
		movies.DELETE("/:id", m.deleteMovie)
	}
}

// @Summary Get All Movies
// @Security ApiKeyAuth
// @Tags movies
// @Description get all movies
// @ID get-all-movies
// @Accept  json
// @Produce  json
// @Success 200 {object} domain.ListMovie
// @Failure 400,404 {object} BadRequestErr
// @Failure 500 {object} InternalServerErr
// @Failure default {object} InternalServerErr
// @Router /movies [get]
func (m Movie) getAllMovies(ctx *gin.Context) {
	movies, err := m.movieService.List(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewInternalServerErr("transport | movieService.List error"))
		return
	}

	ctx.JSON(http.StatusOK, movies)
}

// @Summary Get  Movie By ID
// @Security ApiKeyAuth
// @Tags movies
// @Description get movie by id
// @ID get-movie
// @Accept  json
// @Produce  json
// @Param id path int true "Movie ID"
// @Success 200 {object} domain.Movie
// @Failure 400,404 {object} BadRequestErr
// @Failure 500 {object} InternalServerErr
// @Failure default {object} InternalServerErr
// @Router /movies/{id} [get]
func (m Movie) getMovie(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		fields := map[string]string{"id": "should be an integer"}
		ctx.JSON(http.StatusBadRequest, NewBadRequestErr("validation error", fields))
		return
	}

	movie, err := m.movieService.Get(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(http.StatusNotFound, NewNotFoundErr("movie not found"))
		default:
			ctx.JSON(http.StatusInternalServerError, NewInternalServerErr("transport | movieService.Get error"))
		}

		return
	}

	ctx.JSON(http.StatusOK, movie)
}

// @Summary Create Movie
// @Security ApiKeyAuth
// @Tags movies
// @Description create movie
// @ID create-movie
// @Accept  json
// @Produce  json
// @Param input body domain.Movie true "movie description"
// @Success 200 {object} domain.Movie
// @Failure 400,404 {object} BadRequestErr
// @Failure 500 {object} InternalServerErr
// @Failure default {object} InternalServerErr
// @Router /movies/ [post]
func (m Movie) createMovie(ctx *gin.Context) {
	var movie domain.Movie
	if err := ctx.BindJSON(&movie); err != nil {
		ctx.JSON(http.StatusBadRequest, NewBadRequestErr("cannot parse body", nil))
		return
	}

	createdMovie, err := m.movieService.Create(ctx, movie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewInternalServerErr("transport | movieService.Create error"))
		return
	}

	ctx.JSON(http.StatusCreated, createdMovie)
}

// @Summary Update Movie By ID
// @Security ApiKeyAuth
// @Tags movies
// @Description update movie by id
// @ID update-movie
// @Accept  json
// @Produce  json
// @Param id path int true "Movie ID"
// @Param input body domain.Movie true "movie description"
// @Success 200 {object} domain.Movie
// @Failure 400,404 {object} BadRequestErr
// @Failure 500 {object} InternalServerErr
// @Failure default {object} InternalServerErr
// @Router /movies/{id} [put]
func (m Movie) updateMovie(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		fields := map[string]string{"id": "should be an integer"}
		ctx.JSON(http.StatusBadRequest, NewBadRequestErr("validation error", fields))
		return
	}

	var movie domain.Movie
	if err := ctx.BindJSON(&movie); err != nil {
		ctx.JSON(http.StatusBadRequest, NewBadRequestErr("cannot parse body", nil))
		return
	}

	updatedMovie, err := m.movieService.Update(ctx, id, movie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewInternalServerErr("transport | movieService.Update error"))
		return
	}

	ctx.JSON(http.StatusOK, updatedMovie)
}

// @Summary Delete  Movie By ID
// @Security ApiKeyAuth
// @Tags movies
// @Description delete movie by id
// @ID delete-movie
// @Accept  json
// @Produce  json
// @Param id path int true "Movie ID"
// @Success 200 {object} domain.Movie
// @Failure 400,404 {object} BadRequestErr
// @Failure 500 {object} InternalServerErr
// @Failure default {object} InternalServerErr
// @Router /movies/{id} [delete]
func (m Movie) deleteMovie(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		fields := map[string]string{"id": "should be an integer"}
		ctx.JSON(http.StatusBadRequest, NewBadRequestErr("validation error", fields))
		return
	}

	if err := m.movieService.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, NewInternalServerErr("transport | movieService.Delete error"))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
