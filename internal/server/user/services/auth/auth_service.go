//go:build !test
// +build !test

package auth

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/anoriar/gophkeeper/internal/server/shared/services/uuid"

	"github.com/anoriar/gophkeeper/internal/server/shared/config"
	"github.com/anoriar/gophkeeper/internal/server/user/services/auth/internal/services/token/jwt"

	sharedErrors "github.com/anoriar/gophkeeper/internal/server/shared/errors"
	"github.com/anoriar/gophkeeper/internal/server/user/dto/auth"
	"github.com/anoriar/gophkeeper/internal/server/user/dto/requests/login"
	"github.com/anoriar/gophkeeper/internal/server/user/dto/requests/register"
	"github.com/anoriar/gophkeeper/internal/server/user/repository"
	"github.com/anoriar/gophkeeper/internal/server/user/services/auth/internal/factory/salt"
	user2 "github.com/anoriar/gophkeeper/internal/server/user/services/auth/internal/factory/user"
	"github.com/anoriar/gophkeeper/internal/server/user/services/auth/internal/services/password"
	"github.com/anoriar/gophkeeper/internal/server/user/services/auth/internal/services/token"
	"github.com/anoriar/gophkeeper/internal/server/user/services/auth/internal/services/token/tokenerrors"

	"go.uber.org/zap"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type AuthService struct {
	userRepository  repository.UserRepositoryInterface
	passwordService password.PasswordServiceInterface
	tokenService    token.TokenSerivceInterface
	userFactory     user2.UserFactoryInterface
	saltFactory     salt.SaltFactoryInterface
	logger          *zap.Logger
}

func NewAuthService(
	userRepository repository.UserRepositoryInterface,
	uuidGen uuid.UUIDGeneratorInterface,
	config *config.Config,
	logger *zap.Logger,
) *AuthService {
	return &AuthService{
		userRepository:  userRepository,
		passwordService: password.NewArgonPasswordService(),
		tokenService:    jwt.NewJWTTokenService(config.JwtSecretKey),
		userFactory:     user2.NewUserFactory(uuidGen),
		saltFactory:     salt.NewSaltFactory(),
		logger:          logger,
	}
}

func (service *AuthService) RegisterUser(ctx context.Context, registerUserDto register.RegisterUserRequestDto) (string, error) {
	salt, err := service.saltFactory.GenerateSalt()
	if err != nil {
		service.logger.Error(err.Error())
		return "", fmt.Errorf("%w: %v", sharedErrors.ErrInternalError, err)
	}
	hashedPassword := service.passwordService.GenerateHashedPassword(registerUserDto.Password, salt)

	newUser := service.userFactory.Create(registerUserDto.Login, hashedPassword, hex.EncodeToString(salt))
	err = service.userRepository.AddUser(ctx, newUser)
	if err != nil {
		if errors.Is(err, sharedErrors.ErrConflict) {
			return "", ErrUserAlreadyExists
		}
		service.logger.Error(err.Error())
		return "", fmt.Errorf("%w: %v", sharedErrors.ErrInternalError, err)
	}

	tokenString, err := service.tokenService.BuildTokenString(auth.UserClaims{UserID: newUser.ID})
	if err != nil {
		service.logger.Error(err.Error())
		return "", fmt.Errorf("%w: %v", sharedErrors.ErrInternalError, err)
	}
	return tokenString, nil
}

func (service *AuthService) LoginUser(ctx context.Context, dto login.LoginUserRequestDto) (string, error) {
	existedUser, err := service.userRepository.GetUserByLogin(ctx, dto.Login)
	if err != nil {
		if errors.Is(err, sharedErrors.ErrNotFound) {
			return "", sharedErrors.ErrUserUnauthorized
		}
		service.logger.Error(err.Error())
		return "", fmt.Errorf("%w: %v", sharedErrors.ErrInternalError, err)
	}
	saltInBytes, err := hex.DecodeString(existedUser.Salt)
	if err != nil {
		service.logger.Error(err.Error())
		return "", fmt.Errorf("%w: %v", sharedErrors.ErrInternalError, err)
	}
	hashedPasswordFromRequest := service.passwordService.GenerateHashedPassword(dto.Password, saltInBytes)

	if hashedPasswordFromRequest != existedUser.Password {
		return "", sharedErrors.ErrUserUnauthorized
	}

	tokenString, err := service.tokenService.BuildTokenString(auth.UserClaims{UserID: existedUser.ID})
	if err != nil {
		service.logger.Error(err.Error())
		return "", fmt.Errorf("%w: %v", sharedErrors.ErrInternalError, err)
	}
	return tokenString, nil
}

func (service *AuthService) ValidateToken(token string) (auth.UserClaims, error) {
	claims, err := service.tokenService.GetUserClaims(token)
	if err != nil {
		if errors.Is(err, tokenerrors.ErrTokenNotValid) {
			return auth.UserClaims{}, sharedErrors.ErrUserUnauthorized
		}
		service.logger.Error(err.Error())
		return auth.UserClaims{}, fmt.Errorf("%w: %v", sharedErrors.ErrInternalError, err)
	}

	return claims, nil
}
