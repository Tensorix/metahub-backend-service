package authpage

import (
	"context"
	"log"

	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	"github.com/Tensorix/metahub-backend-service/onebot"
)

func (s *server) Check(_ context.Context, in *auth.CheckRequest) (*auth.CheckResponse, error) {
	var user User
	result := &auth.CheckResponse{
		Result: auth.CheckResult_CHECK_RESULT_SUCCESS,
	}
	t := in.GetToken()
	username := GetUsername(t)

	if username == "" {
		result.Result = auth.CheckResult_CHECK_RESULT_FAILED
		return result, nil
	}

	if err := onebot.DB.First(&user, "username = ?", username).Error; err != nil {
		result.Result = auth.CheckResult_CHECK_RESULT_FAILED
		log.Println(err)
		return result, nil
	}

	return result, nil
}
