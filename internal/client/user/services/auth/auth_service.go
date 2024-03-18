package auth

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/anoriar/gophkeeper/internal/client/user/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/user/dto/repository/request"
	"github.com/anoriar/gophkeeper/internal/client/user/repository/secret"
	"github.com/anoriar/gophkeeper/internal/client/user/repository/user"
)

type AuthService struct {
	userRepository   user.UserRepositoryInterface
	secretRepository secret.SecretRepositoryInterface
	logger           *zap.Logger
}

func NewAuthService(userRepository user.UserRepositoryInterface, secretRepository secret.SecretRepositoryInterface, logger *zap.Logger) *AuthService {
	return &AuthService{
		userRepository:   userRepository,
		secretRepository: secretRepository,
		logger:           logger,
	}
}

func (a *AuthService) Register(ctx context.Context, command command.RegisterCommand) error {
	token, err := a.userRepository.Register(ctx, request.RegisterRequest{
		Login:    command.UserName,
		Password: command.Password,
	})
	if err != nil {
		a.logger.Error("register error", zap.String("error", err.Error()))
		return fmt.Errorf("register error: %v", err.Error())
	}

	err = a.secretRepository.SaveAuthToken(token)
	if err != nil {
		a.logger.Error("save auth token error", zap.String("error", err.Error()))
		return fmt.Errorf("save auth token error: %v", err.Error())
	}

	err = a.secretRepository.SaveMasterPassword(command.MasterPassword)
	if err != nil {
		a.logger.Error("save master password error", zap.String("error", err.Error()))
		return fmt.Errorf("save master password error: %v", err.Error())
	}
	return nil
}

func (a *AuthService) Login(ctx context.Context, command command.LoginCommand) error {
	token, err := a.userRepository.Login(ctx, request.LoginRequest{
		Login:    command.UserName,
		Password: command.Password,
	})
	if err != nil {
		a.logger.Error("login error", zap.String("error", err.Error()))
		return fmt.Errorf("login error: %v", err.Error())
	}

	err = a.secretRepository.SaveAuthToken(token)
	if err != nil {
		a.logger.Error("save auth token error", zap.String("error", err.Error()))
		return fmt.Errorf("save auth token error: %v", err.Error())
	}
	err = a.secretRepository.SaveMasterPassword(command.MasterPassword)
	if err != nil {
		a.logger.Error("save master password error", zap.String("error", err.Error()))
		return fmt.Errorf("save master password error: %v", err.Error())
	}
	return nil
}
