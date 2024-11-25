package authpage

import (
	"context"
	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	"github.com/golang-jwt/jwt/v5"
)

func (s *login) Check(_ context.Context, in *auth.CheckRequest) (*auth.CheckResponse, error) {
	t := in.GetToken()
	valid := verifyToken(t)
	result := auth.CheckResult_CHECK_RESULT_SUCCESS
	if !valid {
		result = auth.CheckResult_CHECK_RESULT_FAILED
	}
	return &auth.CheckResponse{
		Result: result,
	}, nil
}

func verifyToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return false
	}
	_, err = token.Claims.GetIssuer()
	return err == nil
}
