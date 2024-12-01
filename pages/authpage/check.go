package authpage

import (
	"context"
	"log"

	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	"github.com/golang-jwt/jwt/v5"
)

func (s *server) Check(_ context.Context, in *auth.CheckRequest) (*auth.CheckResponse, error) {
	var user User
	result := &auth.CheckResponse{
		Result: auth.CheckResult_CHECK_RESULT_SUCCESS,
	}
	t := in.GetToken()
	username := GetUsername(t)

	if err := db.First(&user, "username = ?", username).Error; err != nil {
		log.Println(err)
	}
	registerBot(user.Id, username)

	return result, nil
}

func verifyToken(tokenString string) (string, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return "", false
	}
	username, err := token.Claims.GetIssuer()
	return username, err == nil
}
