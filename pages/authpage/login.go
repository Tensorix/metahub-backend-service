package authpage

import (
	"context"
	"errors"
	"time"

	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func (s *server) Login(_ context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {
	result := &auth.LoginResponse{
		Result: auth.LoginResult_LOGIN_RESULT_SUCCESS,
	}
	var err error
	username := in.GetUsername()
	password := in.GetPassword()
	if username == "" || password == "" {
		result.Result = auth.LoginResult_LOGIN_RESULT_FAILED
		return result, nil
	}
	var user User
	err = db.First(&user, "username = ? AND pwd = ?", username, password).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		result.Result = auth.LoginResult_LOGIN_RESULT_FAILED
		return result, nil
	case err != nil:
		result.Result = auth.LoginResult_LOGIN_RESULT_UNSPECIFIED
		return result, nil
	}
	t, err := createToken(username)
	if err != nil {
		result.Result = auth.LoginResult_LOGIN_RESULT_UNSPECIFIED
		return result, nil
	}
	result.Token = t
	return result, nil
}

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.RegisteredClaims{
			Issuer:    username,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireAt)),
		})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
