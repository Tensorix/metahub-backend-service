package authpage

import (
	"context"
	"log"

	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
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
