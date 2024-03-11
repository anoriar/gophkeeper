package auth

import (
	"context"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/anoriar/gophkeeper/internal/server/user/dto/requests/login"
	"github.com/anoriar/gophkeeper/internal/server/user/services/auth/internal/services/token/tokenerrors"

	"github.com/anoriar/gophkeeper/internal/server/shared/app/logger"
	sharedErrors "github.com/anoriar/gophkeeper/internal/server/shared/errors"
	"github.com/anoriar/gophkeeper/internal/server/user/dto/auth"
	"github.com/anoriar/gophkeeper/internal/server/user/dto/requests/register"
	"github.com/anoriar/gophkeeper/internal/server/user/entity"
	"github.com/anoriar/gophkeeper/internal/server/user/repository/mock_user_repository"
	"github.com/anoriar/gophkeeper/internal/server/user/services/auth/internal/factory/salt/mock_salt_factory"
	"github.com/anoriar/gophkeeper/internal/server/user/services/auth/internal/factory/user/mock_user_factory"
	"github.com/anoriar/gophkeeper/internal/server/user/services/auth/internal/services/password/mock_password_service"
	"github.com/anoriar/gophkeeper/internal/server/user/services/auth/internal/services/token/mock_token_service"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAuthService_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mock_user_repository.NewMockUserRepositoryInterface(ctrl)
	saltFactoryMock := mock_salt_factory.NewMockSaltFactoryInterface(ctrl)
	tokenServiceMock := mock_token_service.NewMockTokenSerivceInterface(ctrl)
	passwordServiceMock := mock_password_service.NewMockPasswordServiceInterface(ctrl)
	userFactoryMock := mock_user_factory.NewMockUserFactoryInterface(ctrl)
	loggerMock, err := logger.Initialize("info")
	require.NoError(t, err)

	type args struct {
		ctx             context.Context
		registerUserDto register.RegisterUserRequestDto
	}
	tests := []struct {
		name          string
		mockBehaviour func()
		args          args
		want          string
		err           error
	}{
		{
			name: "success",
			mockBehaviour: func() {
				salt := []byte{2, 3, 24, 32, 11}
				hashedPassword := "dashdiasdss"
				encodedSalt := hex.EncodeToString(salt)
				userID := "ab161651-ec2c-4cbb-a6c6-f8cf414e503d"
				login := "test"
				userPassword := "user-password"
				userMock := entity.User{
					ID:       userID,
					Login:    login,
					Password: hashedPassword,
					Salt:     encodedSalt,
				}
				saltFactoryMock.EXPECT().GenerateSalt().Return(salt, nil).Times(1)
				passwordServiceMock.EXPECT().GenerateHashedPassword(userPassword, salt).Return(hashedPassword).Times(1)
				userFactoryMock.EXPECT().Create(login, hashedPassword, encodedSalt).Return(userMock).Times(1)
				userRepositoryMock.EXPECT().AddUser(gomock.Any(), userMock).Return(nil).Times(1)
				tokenServiceMock.EXPECT().BuildTokenString(auth.UserClaims{UserID: userID}).Return("result-token", nil).Times(1)
			},
			args: args{
				ctx: context.Background(),
				registerUserDto: register.RegisterUserRequestDto{
					Login:    "test",
					Password: "user-password",
				},
			},
			want: "result-token",
		},
		{
			name: "conflict error add user",
			mockBehaviour: func() {
				salt := []byte{2, 3, 24, 32, 11}
				hashedPassword := "dashdiasdss"
				encodedSalt := hex.EncodeToString(salt)
				userID := "ab161651-ec2c-4cbb-a6c6-f8cf414e503d"
				login := "test"
				userPassword := "user-password"
				userMock := entity.User{
					ID:       userID,
					Login:    login,
					Password: hashedPassword,
					Salt:     encodedSalt,
				}
				saltFactoryMock.EXPECT().GenerateSalt().Return(salt, nil).Times(1)
				passwordServiceMock.EXPECT().GenerateHashedPassword(userPassword, salt).Return(hashedPassword).Times(1)
				userFactoryMock.EXPECT().Create(login, hashedPassword, encodedSalt).Return(userMock).Times(1)
				userRepositoryMock.EXPECT().AddUser(gomock.Any(), userMock).Return(sharedErrors.ErrConflict).Times(1)
				tokenServiceMock.EXPECT().BuildTokenString(auth.UserClaims{UserID: userID}).Times(0)
			},
			args: args{
				ctx: context.Background(),
				registerUserDto: register.RegisterUserRequestDto{
					Login:    "test",
					Password: "user-password",
				},
			},
			want: "",
			err:  ErrUserAlreadyExists,
		},
		{
			name: "internal server error factory",
			mockBehaviour: func() {
				saltFactoryMock.EXPECT().GenerateSalt().Return([]byte{}, errors.New("internal server error")).Times(1)
			},
			args: args{
				ctx: context.Background(),
				registerUserDto: register.RegisterUserRequestDto{
					Login:    "test",
					Password: "user-password",
				},
			},
			want: "",
			err:  sharedErrors.ErrInternalError,
		},
		{
			name: "internal server error add user",
			mockBehaviour: func() {
				salt := []byte{2, 3, 24, 32, 11}
				hashedPassword := "dashdiasdss"
				encodedSalt := hex.EncodeToString(salt)
				userID := "ab161651-ec2c-4cbb-a6c6-f8cf414e503d"
				login := "test"
				userPassword := "user-password"
				userMock := entity.User{
					ID:       userID,
					Login:    login,
					Password: hashedPassword,
					Salt:     encodedSalt,
				}
				saltFactoryMock.EXPECT().GenerateSalt().Return(salt, nil).Times(1)
				passwordServiceMock.EXPECT().GenerateHashedPassword(userPassword, salt).Return(hashedPassword).Times(1)
				userFactoryMock.EXPECT().Create(login, hashedPassword, encodedSalt).Return(userMock).Times(1)
				userRepositoryMock.EXPECT().AddUser(gomock.Any(), userMock).Return(sharedErrors.ErrInternalError).Times(1)
			},
			args: args{
				ctx: context.Background(),
				registerUserDto: register.RegisterUserRequestDto{
					Login:    "test",
					Password: "user-password",
				},
			},
			want: "",
			err:  sharedErrors.ErrInternalError,
		},
		{
			name: "internal server error build token",
			mockBehaviour: func() {
				salt := []byte{2, 3, 24, 32, 11}
				hashedPassword := "dashdiasdss"
				encodedSalt := hex.EncodeToString(salt)
				userID := "ab161651-ec2c-4cbb-a6c6-f8cf414e503d"
				login := "test"
				userPassword := "user-password"
				userMock := entity.User{
					ID:       userID,
					Login:    login,
					Password: hashedPassword,
					Salt:     encodedSalt,
				}
				saltFactoryMock.EXPECT().GenerateSalt().Return(salt, nil).Times(1)
				passwordServiceMock.EXPECT().GenerateHashedPassword(userPassword, salt).Return(hashedPassword).Times(1)
				userFactoryMock.EXPECT().Create(login, hashedPassword, encodedSalt).Return(userMock).Times(1)
				userRepositoryMock.EXPECT().AddUser(gomock.Any(), userMock).Return(nil).Times(1)
				tokenServiceMock.EXPECT().BuildTokenString(auth.UserClaims{UserID: userID}).Return("", sharedErrors.ErrInternalError).Times(1)
			},
			args: args{
				ctx: context.Background(),
				registerUserDto: register.RegisterUserRequestDto{
					Login:    "test",
					Password: "user-password",
				},
			},
			want: "",
			err:  sharedErrors.ErrInternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			service := &AuthService{
				userRepository:  userRepositoryMock,
				passwordService: passwordServiceMock,
				tokenService:    tokenServiceMock,
				userFactory:     userFactoryMock,
				saltFactory:     saltFactoryMock,
				logger:          loggerMock,
			}
			got, err := service.RegisterUser(tt.args.ctx, tt.args.registerUserDto)
			if tt.err != nil {
				if !errors.Is(err, tt.err) {
					t.Errorf("RegisterUser() error expectation: got = %v, want %v", err, tt.err)
				}
			}
			if got != tt.want {
				t.Errorf("RegisterUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mock_user_repository.NewMockUserRepositoryInterface(ctrl)
	tokenServiceMock := mock_token_service.NewMockTokenSerivceInterface(ctrl)
	passwordServiceMock := mock_password_service.NewMockPasswordServiceInterface(ctrl)
	loggerMock, err := logger.Initialize("info")
	require.NoError(t, err)

	type args struct {
		ctx                 context.Context
		loginUserRequestDto login.LoginUserRequestDto
	}
	tests := []struct {
		name          string
		mockBehaviour func()
		args          args
		want          string
		err           error
	}{
		{
			name: "success",
			mockBehaviour: func() {
				salt := []byte{2, 3, 24, 32, 11}
				hashedPassword := "dashdiasdss"
				encodedSalt := hex.EncodeToString(salt)
				userID := "ab161651-ec2c-4cbb-a6c6-f8cf414e503d"
				login := "test"
				userPassword := "user-password"
				userMock := entity.User{
					ID:       userID,
					Login:    login,
					Password: hashedPassword,
					Salt:     encodedSalt,
				}
				userRepositoryMock.EXPECT().GetUserByLogin(gomock.Any(), login).Return(userMock, nil)
				passwordServiceMock.EXPECT().GenerateHashedPassword(userPassword, salt).Return(hashedPassword).Times(1)
				tokenServiceMock.EXPECT().BuildTokenString(auth.UserClaims{UserID: userID}).Return("result-token", nil).Times(1)
			},
			args: args{
				ctx: context.Background(),
				loginUserRequestDto: login.LoginUserRequestDto{
					Login:    "test",
					Password: "user-password",
				},
			},
			want: "result-token",
		},
		{
			name: "get user by login internal error",
			mockBehaviour: func() {
				login := "test"
				userRepositoryMock.EXPECT().GetUserByLogin(gomock.Any(), login).Return(entity.User{}, errors.New("err"))
			},
			args: args{
				ctx: context.Background(),
				loginUserRequestDto: login.LoginUserRequestDto{
					Login:    "test",
					Password: "user-password",
				},
			},
			want: "",
			err:  sharedErrors.ErrInternalError,
		},
		{
			name: "get user by login not found",
			mockBehaviour: func() {
				login := "test"
				userRepositoryMock.EXPECT().GetUserByLogin(gomock.Any(), login).Return(entity.User{}, sharedErrors.ErrNotFound)
			},
			args: args{
				ctx: context.Background(),
				loginUserRequestDto: login.LoginUserRequestDto{
					Login:    "test",
					Password: "user-password",
				},
			},
			want: "",
			err:  sharedErrors.ErrUserUnauthorized,
		},
		{
			name: "salt error",
			mockBehaviour: func() {
				hashedPassword := "dashdiasdss"
				encodedSalt := "h8rhg04tg"
				userID := "ab161651-ec2c-4cbb-a6c6-f8cf414e503d"
				login := "test"
				userMock := entity.User{
					ID:       userID,
					Login:    login,
					Password: hashedPassword,
					Salt:     encodedSalt,
				}
				userRepositoryMock.EXPECT().GetUserByLogin(gomock.Any(), login).Return(userMock, nil)
			},
			args: args{
				ctx: context.Background(),
				loginUserRequestDto: login.LoginUserRequestDto{
					Login:    "test",
					Password: "user-password",
				},
			},
			want: "",
			err:  sharedErrors.ErrInternalError,
		},
		{
			name: "passwords not similar",
			mockBehaviour: func() {
				salt := []byte{2, 3, 24, 32, 11}
				hashedPassword := "dashdiasdss"
				encodedSalt := hex.EncodeToString(salt)
				userID := "ab161651-ec2c-4cbb-a6c6-f8cf414e503d"
				login := "test"
				userPassword := "user-password"
				userMock := entity.User{
					ID:       userID,
					Login:    login,
					Password: hashedPassword,
					Salt:     encodedSalt,
				}
				userRepositoryMock.EXPECT().GetUserByLogin(gomock.Any(), login).Return(userMock, nil)
				passwordServiceMock.EXPECT().GenerateHashedPassword(userPassword, salt).Return("12345").Times(1)
			},
			args: args{
				ctx: context.Background(),
				loginUserRequestDto: login.LoginUserRequestDto{
					Login:    "test",
					Password: "user-password",
				},
			},
			want: "",
			err:  sharedErrors.ErrUserUnauthorized,
		},
		{
			name: "build token internal error",
			mockBehaviour: func() {
				salt := []byte{2, 3, 24, 32, 11}
				hashedPassword := "dashdiasdss"
				encodedSalt := hex.EncodeToString(salt)
				userID := "ab161651-ec2c-4cbb-a6c6-f8cf414e503d"
				login := "test"
				userPassword := "user-password"
				userMock := entity.User{
					ID:       userID,
					Login:    login,
					Password: hashedPassword,
					Salt:     encodedSalt,
				}
				userRepositoryMock.EXPECT().GetUserByLogin(gomock.Any(), login).Return(userMock, nil)
				passwordServiceMock.EXPECT().GenerateHashedPassword(userPassword, salt).Return(hashedPassword).Times(1)
				tokenServiceMock.EXPECT().BuildTokenString(auth.UserClaims{UserID: userID}).Return("", errors.New("err")).Times(1)
			},
			args: args{
				ctx: context.Background(),
				loginUserRequestDto: login.LoginUserRequestDto{
					Login:    "test",
					Password: "user-password",
				},
			},
			want: "",
			err:  sharedErrors.ErrInternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			service := &AuthService{
				userRepository:  userRepositoryMock,
				passwordService: passwordServiceMock,
				tokenService:    tokenServiceMock,
				logger:          loggerMock,
			}
			got, err := service.LoginUser(tt.args.ctx, tt.args.loginUserRequestDto)

			if tt.err != nil {
				if !errors.Is(err, tt.err) {
					t.Errorf("LoginUser() error expectation: got = %v, want %v", err, tt.err)
				}
			}
			if got != tt.want {
				t.Errorf("LoginUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mock_user_repository.NewMockUserRepositoryInterface(ctrl)
	tokenServiceMock := mock_token_service.NewMockTokenSerivceInterface(ctrl)
	passwordServiceMock := mock_password_service.NewMockPasswordServiceInterface(ctrl)
	loggerMock, err := logger.Initialize("info")
	require.NoError(t, err)

	type args struct {
		token string
	}
	tests := []struct {
		name          string
		mockBehaviour func()
		args          args
		want          auth.UserClaims
		err           error
	}{
		{
			name: "success",
			mockBehaviour: func() {
				userID := "ab161651-ec2c-4cbb-a6c6-f8cf414e503d"
				userClaimsMock := auth.UserClaims{
					UserID: userID,
				}
				tokenServiceMock.EXPECT().GetUserClaims("rweawr").Return(userClaimsMock, nil)
			},
			args: args{
				token: "rweawr",
			},
			want: auth.UserClaims{
				UserID: "ab161651-ec2c-4cbb-a6c6-f8cf414e503d",
			},
		},
		{
			name: "token not valid error",
			mockBehaviour: func() {
				userClaimsMock := auth.UserClaims{}
				tokenServiceMock.EXPECT().GetUserClaims("rweawr").Return(userClaimsMock, tokenerrors.ErrTokenNotValid)
			},
			args: args{
				token: "rweawr",
			},
			want: auth.UserClaims{},
			err:  sharedErrors.ErrUserUnauthorized,
		},
		{
			name: "internal error",
			mockBehaviour: func() {
				userClaimsMock := auth.UserClaims{}
				tokenServiceMock.EXPECT().GetUserClaims("rweawr").Return(userClaimsMock, errors.New("err"))
			},
			args: args{
				token: "rweawr",
			},
			want: auth.UserClaims{},
			err:  sharedErrors.ErrInternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			service := &AuthService{
				userRepository:  userRepositoryMock,
				passwordService: passwordServiceMock,
				tokenService:    tokenServiceMock,
				logger:          loggerMock,
			}
			got, err := service.ValidateToken(tt.args.token)
			if tt.err != nil {
				if !errors.Is(err, tt.err) {
					t.Errorf("ValidateToken() error expectation: got = %v, want %v", err, tt.err)
				}
			}
			if got != tt.want {
				t.Errorf("ValidateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
