package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/lukinairina90/crud_audit_log/pkg/domain/audit"
	"github.com/lukinairina90/crud_movies/internal/transport/rest"
	"github.com/sirupsen/logrus"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lukinairina90/crud_movies/internal/domain"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UsersRepository interface {
	Create(ctx context.Context, user domain.User) (int64, error)
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
}

//type InMemoryCache[K comparable, V any] interface {
//	Set(key K, value V, ttl time.Duration) error
//}

type SessionRepository interface {
	Create(ctx context.Context, token domain.RefreshSession) error
	Get(ctx context.Context, token string) (domain.RefreshSession, error)
}

type AuditClient interface {
	SendLogRequest(ctx context.Context, req audit.LogItem) error
}

type Users struct {
	repo        UsersRepository
	sessionRepo SessionRepository
	hasher      PasswordHasher

	hmacSecret []byte
	tokenTtl   time.Duration

	auditClient AuditClient
}

func NewUsers(repo UsersRepository, sessionRepo SessionRepository, auditClient AuditClient, hasher PasswordHasher, hmacSecret []byte, tokenTtl time.Duration) *Users {
	return &Users{
		repo:        repo,
		sessionRepo: sessionRepo,
		hasher:      hasher,
		hmacSecret:  hmacSecret,
		tokenTtl:    tokenTtl,
		auditClient: auditClient,
	}
}

func (s *Users) SignUp(ctx context.Context, inp domain.SignUpInput) error {
	var id int64
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	id, err = s.repo.Create(ctx, user) // Сделай чтобы репа возвращала объект юзера уже с заполненным ID. Подсказка: RETURNING в SQL запросе
	if err != nil {
		return err
	}

	if err := s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Entity:    audit.ENTITY_USER,
		Actions:   audit.ACTION_REGISTER,
		EntityID:  id,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Users SignUp",
		}).Error("failed to send log request:", err)
	}

	return nil
}

func (s *Users) SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error) {
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return "", "", err
	}

	user, err := s.repo.GetByCredentials(ctx, inp.Email, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", domain.ErrUserNotFound
		}
		return "", "", err
	}

	return s.generateTokens(ctx, user.ID)
}

func (s *Users) ParseToken(_ context.Context, token string) (int64, error) {

	fmt.Println(token)

	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpecting signing method %v", token.Header["alg"])
		}
		return s.hmacSecret, nil
	})
	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return int64(id), nil
}

func (s *Users) generateTokens(ctx context.Context, userId int64) (string, string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.Itoa(int(userId)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenTtl)),
	})

	accessToken, err := t.SignedString(s.hmacSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := s.sessionRepo.Create(ctx, domain.RefreshSession{
		UserID:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (s *Users) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	session, err := s.sessionRepo.Get(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", rest.ErrRefreshTokenExpired
	}

	return s.generateTokens(ctx, session.UserID)
}
