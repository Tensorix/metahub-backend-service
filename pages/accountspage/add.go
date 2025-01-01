package accountspage

import (
	"context"

	accounts "github.com/Tensorix/metahub-backend-service/gen/proto/v1/accounts"
	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	"github.com/Tensorix/metahub-backend-service/onebot"
	"github.com/Tensorix/metahub-backend-service/pages/authpage"
)

func (server *server) AddAccount(_ context.Context, in *accounts.AddAccountRequest) (*accounts.AddAccountResponse, error) {
	var user authpage.User
	var acc Account
	token := in.Token.Token
	response := &accounts.AddAccountResponse{
		Result: &auth.CheckResponse{
			Result: auth.CheckResult_CHECK_RESULT_FAILED,
		},
		AddResult: accounts.AddAccountResult_ADD_RESULT_FAILED,
	}
	username := authpage.GetUsername(token)
	if username == "" {
		return response, nil
	}
	response.Result.Result = auth.CheckResult_CHECK_RESULT_SUCCESS
	if err := onebot.DB.First(&user, "username = ?", username).Error; err != nil {
		response.AddResult = accounts.AddAccountResult_ADD_RESULT_FAILED
		return response, nil
	}
	if err := onebot.DB.First(&acc, "(account_tag = ? AND user_id = ?) OR port = ?", in.Data.AccountTag, user.Id, in.Data.Port).Error; err == nil {
		response.AddResult = accounts.AddAccountResult_ADD_RESULT_EXISTS
		return response, nil
	}
	acc = Account{
		AccountTag: in.Data.AccountTag,
		UserID:     user.Id,
		IP:         in.Data.Ip,
		Port:       in.Data.Port,
	}
	if err := onebot.DB.Create(&acc).Error; err != nil {
		return response, nil
	}
	response.AddResult = accounts.AddAccountResult_ADD_RESULT_SUCCESS
	response.Id = acc.ID
	response.Result.Result = auth.CheckResult_CHECK_RESULT_SUCCESS
	onebot.NewOnebot(username, in.Data.AccountTag, in.Data.Ip, int(in.Data.Port), user.Id, acc.ID)
	return response, nil
}
