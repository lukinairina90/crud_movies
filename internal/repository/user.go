package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/lukinairina90/crud_movies/internal/domain"
)

type Users struct {
	db *sqlx.DB
}

func NewUsers(db *sqlx.DB) *Users {
	return &Users{db: db}
}

func (r *Users) Create(ctx context.Context, user domain.User) (int64, error) {
	var id int64
	err := r.db.QueryRowxContext(ctx, "INSERT INTO users (name, email, password, registered_at) values ($1, $2, $3, $4) RETURNING id", user.Name, user.Email, user.Password, user.RegisteredAt).Scan(&id)
	//	if err := m.db.QueryRowxContext(ctx, "UPDATE movie SET name=$1, description=$2, production_year=$3, genre=$4, actors=$5, poster=$6 WHERE id=$7 RETURNING *", mMovie.Name, mMovie.Description, mMovie.ProductionYear, mMovie.Genre, mMovie.Actors, mMovie.Poster, id).StructScan(&mMovie); err != nil {
	return id, err
}

func (r *Users) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, name, email, registered_at FROM users WHERE email=$1 AND password=$2", email, password).Scan(&user.ID, &user.Name, &user.Email, &user.RegisteredAt)

	return user, err
}
