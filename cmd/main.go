package main

import (
	"fmt"
	"log"

	"github.com/lukinairina90/crud_movies/internal/transport/grpc"

	"github.com/gin-gonic/gin"
	"github.com/lukinairina90/crud_movies/internal/domain"
	"github.com/lukinairina90/crud_movies/internal/repository"
	"github.com/lukinairina90/crud_movies/internal/service"
	"github.com/lukinairina90/crud_movies/internal/transport/rest"
	"github.com/lukinairina90/crud_movies/pkg/config"
	"github.com/lukinairina90/crud_movies/pkg/database"
	"github.com/lukinairina90/crud_movies/pkg/hash"
	"github.com/lukinairina90/in_memory_cache/generic_cache"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/lukinairina90/crud_movies/docs"
)

// ENV
// DB_HOST=localhost;DB_NAME=movies;DB_PASS=goLANGninja;DB_PORT=5432;DB_SSL_MODE=false;DB_USER=postgres;PORT=8080;TOKEN_TTL=24h;CACHE_TTL=10m

// @title CRUD_movies
// @version 1.0
// @description API Server Movies Application

// @host localhost:8080
// BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.Parse()
	if err != nil {
		logrus.Fatalf("error psring config: %s", err.Error())
	}

	// init db
	db, err := database.CreateConn(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.SSLMode)
	if err != nil {
		logrus.Fatalf("failed to connection db: %s", err.Error())
	}

	tokenSecret := []byte("sample secret")

	defer db.Close()

	// init deps
	hasher := hash.NewMD5Hasher("salt")
	movieCache := generic_cache.New[string, domain.Movie]()

	movieRepository := repository.NewMovie(db)
	//cachedMovieRepo := repository.NewCachedMovie(movieRepository, movieCache)

	movieService := service.NewMovie(movieRepository, movieCache)
	moviesTransport := rest.NewMovie(movieService)

	usersRepository := repository.NewUsers(db)
	tokensRepository := repository.NewTokens(db)

	auditClient, err := grpc.NewClient(9000)
	if err != nil {
		log.Fatal(err)
	}

	usersService := service.NewUsers(usersRepository, tokensRepository, auditClient, hasher, tokenSecret, cfg.TokenTTL)

	authTransport := rest.NewAuth(usersService)

	// init routes
	g := gin.New()
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	g.Use(rest.LoggingMiddleware())
	authTransport.InjectRoutes(g)
	moviesTransport.InjectRoutes(g, authTransport.AuthMiddleware())

	fmt.Println("Server run...")
	if err := g.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		logrus.Fatalf("error occured while running http server %s", err.Error())
	}
}
