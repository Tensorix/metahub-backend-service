package accountspage

import (
	"context"

	accounts "github.com/Tensorix/metahub-backend-service/gen/proto/v1/accounts"
	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	"github.com/Tensorix/metahub-backend-service/onebot"
	"github.com/Tensorix/metahub-backend-service/pages/authpage"
)

type Account struct {
	ID         int32
	AccountTag string
	UserID     int32
	IP         string
	Port       int32
}

func (server *server) QueryAccount(_ context.Context, in *accounts.QueryAccountRequest) (*accounts.QueryAccountResponse, error) {
	var user authpage.User
	var accountls []Account
	token := in.Token.Token
	response := &accounts.QueryAccountResponse{
		Result: &auth.CheckResponse{
			Result: auth.CheckResult_CHECK_RESULT_FAILED,
		},
		QueryResult: accounts.QueryAccountResult_QUERY_ACCOUNT_RESULT_FAILED,
		Data:        []*accounts.AccountData{},
	}
	username := authpage.GetUsername(token)
	if username == "" {
		return response, nil
	}
	response.Result.Result = auth.CheckResult_CHECK_RESULT_SUCCESS
	if err := onebot.DB.First(&user, "username = ?", username).Error; err != nil {
		return response, nil
	}
	if err := onebot.DB.Where("user_id = ?", user.Id).Find(&accountls).Error; err != nil {
		return response, nil
	}
	for _, acc := range accountls {
		response.Data = append(response.Data, &accounts.AccountData{
			Id:         acc.ID,
			AccountTag: acc.AccountTag,
			// UserId:     user.Id,
			Ip:   acc.IP,
			Port: acc.Port,
		})
	}
	response.QueryResult = accounts.QueryAccountResult_QUERY_ACCOUNT_RESULT_SUCCESS
	return response, nil
}
