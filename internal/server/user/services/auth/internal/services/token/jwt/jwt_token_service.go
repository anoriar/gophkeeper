package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/anoriar/gophkeeper/internal/server/user/dto/auth"
	"github.com/anoriar/gophkeeper/internal/server/user/services/auth/internal/services/token/tokenerrors"
)

const tokenExpired = time.Hour * 3

type JWTTokenService struct {
	secretKey string
}

func NewJWTTokenService(secretKey string) *JWTTokenService {
	return &JWTTokenService{secretKey: secretKey}
}

func (service *JWTTokenService) BuildTokenString(userClaims auth.UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{ // когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpired))},
		UserID: userClaims.UserID,
	})

	tokenString, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return tokenString, nil
}

func (service *JWTTokenService) GetUserClaims(tokenString string) (auth.UserClaims, error) {
	jwtClaims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, jwtClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.secretKey), nil
	})
	if err != nil {
		var validationError *jwt.ValidationError
		if errors.As(err, &validationError) {
			return auth.UserClaims{}, tokenerrors.ErrTokenNotValid
		}
		return auth.UserClaims{}, fmt.Errorf("%w", err)
	}
	if !token.Valid {
		return auth.UserClaims{}, tokenerrors.ErrTokenNotValid
	}

	return auth.UserClaims{
		UserID: jwtClaims.UserID,
	}, nil
}
