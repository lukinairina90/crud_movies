package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/lukinairina90/crud_movies/internal/repository"
	"github.com/lukinairina90/crud_movies/internal/service"
	"github.com/lukinairina90/crud_movies/internal/transport/rest"
	"github.com/lukinairina90/crud_movies/pkg/config"
	"github.com/lukinairina90/crud_movies/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	srv := gin.New()

	cfg, err := config.Parse()
	if err != nil {
		logrus.Fatalf("error psring config: %s", err.Error())
	}

	db, err := database.CreateConn(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.SSLMode)
	if err != nil {
		logrus.Fatalf("failed to connection db: %s", err.Error())
	}

	movieRepository := repository.NewMovie(db)
	movieService := service.NewMovie(movieRepository)
	movieTransport := rest.NewMovie(movieService)

	srv.GET("/movies", movieTransport.List)
	srv.GET("/movie/:id", movieTransport.Get)
	srv.POST("/movie", movieTransport.Create)
	srv.PUT("/movie/:id", movieTransport.Update)
	srv.DELETE("/movie/:id", movieTransport.Delete)

	if err := srv.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		logrus.Fatalf("error occured while running http server %s", err.Error())
	}
}
