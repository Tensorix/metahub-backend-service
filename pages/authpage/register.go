package authpage

import (
	"context"

	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
)

func (s *login) Register(_ context.Context, in *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	var err error
	username := in.GetUsername()
	password := in.GetPassword()
	response := &auth.RegisterResponse{
		Result: auth.RegisterResult_REGISTER_RESULT_SUCCESS,
	}
	if username == "" || password == "" {
		response.Result = auth.RegisterResult_REGISTER_RESULT_VALUE_NULL
		return response, nil
	}
	var user User
	err = db.First(&user, "username = ?", username).Error
	if err == nil {
		response.Result = auth.RegisterResult_REGISTER_RESULT_EXISTS
		return response, nil
	}
	user = User{Username: username, Pwd: password}
	err = db.Create(&user).Error
	if err != nil {
		response.Result = auth.RegisterResult_REGISTER_RESULT_UNSPECIFIED
	}
	return response, nil
}
