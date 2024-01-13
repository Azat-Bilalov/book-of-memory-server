package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/redis"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/repository"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type InterfaceAuthUsecase interface {
	Login(login string, passwd string) (*ds.LoginResponse, error)
	Register(firstName string, lastName string, login string, passwd string) (*ds.RegisterResponse, error)
	Logout() error
}

type AuthUsecase struct {
	userRepository repository.InterfaceUserRepository
	redis          *redis.Client
}

func NewAuthUsecase(userRepository repository.InterfaceUserRepository, redis *redis.Client) *AuthUsecase {
	return &AuthUsecase{userRepository: userRepository, redis: redis}
}

func (u *AuthUsecase) Login(login string, passwd string) (*ds.LoginResponse, error) {
	user, err := u.userRepository.FindByEmail(login)
	if err != nil {
		return nil, err
	}
	if user.Passwd != passwd {
		return nil, nil
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, ds.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "book-of-memory",
		},
		User_id: user.User_id,
		Role:    user.Role,
	})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}
	return &ds.LoginResponse{
		ExpiresIn:   int(time.Hour) * 24,
		AccessToken: tokenString,
		TokenType:   "Bearer",
		Role:        user.Role,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		User_id:     user.User_id,
	}, nil
}

func (u *AuthUsecase) Register(firstName string, lastName string, email string, passwd string) (*ds.RegisterResponse, error) {
	if firstName == "" || lastName == "" || email == "" || passwd == "" {
		return nil, errors.New("не все поля заполнены")
	}
	user := &ds.User{
		User_id:   uuid.New().String(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Passwd:    passwd,
		Role:      ds.USER_ROLE_USER,
	}
	user, err := u.userRepository.Store(user)
	if err != nil {
		return nil, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, ds.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "book-of-memory",
		},
		User_id: user.User_id,
		Role:    user.Role,
	})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}
	return &ds.RegisterResponse{
		ExpiresIn:   int(time.Hour) * 24,
		AccessToken: tokenString,
		TokenType:   "Bearer",
		Role:        user.Role,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		User_id:     user.User_id,
	}, nil
}

func (u *AuthUsecase) Logout(ctx context.Context, tokenString string) error {
	err := u.redis.WriteJWTToBlacklist(ctx, tokenString, time.Hour*24)
	if err != nil {
		return err
	}
	return nil
}
