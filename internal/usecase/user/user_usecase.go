package user

import (
	"context"
	"errors"
	"time"

	"github.com/mashurimansur/goCMS/internal/domain/user"
	"github.com/mashurimansur/goCMS/internal/utils/token"
	"golang.org/x/crypto/bcrypt"
)

type UseCase interface {
	Register(ctx context.Context, u *user.User, password string) error
	Login(ctx context.Context, email, password string) (string, *user.User, error)
	GetProfile(ctx context.Context, id string) (*user.User, error)
	UpdateProfile(ctx context.Context, u *user.User) error
	ListUsers(ctx context.Context, limit, offset int) ([]*user.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type userUseCase struct {
	userRepo      user.Repository
	tokenMaker    token.Maker
	tokenDuration time.Duration
}

func NewUserUseCase(userRepo user.Repository, tokenMaker token.Maker, tokenDuration time.Duration) UseCase {
	return &userUseCase{
		userRepo:      userRepo,
		tokenMaker:    tokenMaker,
		tokenDuration: tokenDuration,
	}
}

func (uc *userUseCase) Register(ctx context.Context, u *user.User, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	if u.Role == "" {
		u.Role = "user"
	}
	if u.Status == "" {
		u.Status = "active"
	}
	return uc.userRepo.Create(ctx, u)
}

func (uc *userUseCase) Login(ctx context.Context, email, password string) (string, *user.User, error) {
	u, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}
	if u == nil {
		return "", nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	accessToken, _, err := uc.tokenMaker.CreateToken(u.ID, uc.tokenDuration)
	if err != nil {
		return "", nil, err
	}

	return accessToken, u, nil
}

func (uc *userUseCase) GetProfile(ctx context.Context, id string) (*user.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

func (uc *userUseCase) UpdateProfile(ctx context.Context, u *user.User) error {
	return uc.userRepo.Update(ctx, u)
}

func (uc *userUseCase) ListUsers(ctx context.Context, limit, offset int) ([]*user.User, error) {
	return uc.userRepo.List(ctx, limit, offset)
}

func (uc *userUseCase) DeleteUser(ctx context.Context, id string) error {
	return uc.userRepo.Delete(ctx, id)
}
