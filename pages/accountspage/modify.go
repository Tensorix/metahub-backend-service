package accountspage

import (
	"context"

	accounts "github.com/Tensorix/metahub-backend-service/gen/proto/v1/accounts"
	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	"github.com/Tensorix/metahub-backend-service/onebot"
	"github.com/Tensorix/metahub-backend-service/pages/authpage"
)

func (server *server) ModifyAccount(_ context.Context, in *accounts.ModifyAccountRequest) (*accounts.ModifyAccountResponse, error) {
	var user authpage.User
	var acc Account
	token := in.Token.Token
	response := &accounts.ModifyAccountResponse{
		Result: &auth.CheckResponse{
			Result: auth.CheckResult_CHECK_RESULT_FAILED,
		},
		ModifyResult: accounts.ModifyAccountResult_MODIFY_ACCOUNT_RESULT_FAILED,
	}
	username := authpage.GetUsername(token)
	if username == "" {
		return response, nil
	}
	response.Result.Result = auth.CheckResult_CHECK_RESULT_SUCCESS
	if err := onebot.DB.First(&user, "username = ?", username).Error; err != nil {
		response.ModifyResult = accounts.ModifyAccountResult_MODIFY_ACCOUNT_RESULT_FAILED
		return response, nil
	}
	if err := onebot.DB.First(&acc, "account_tag = ? AND user_id = ?", in.Data.AccountTag, user.Id).Error; err != nil {
		response.ModifyResult = accounts.ModifyAccountResult_MODIFY_ACCOUNT_RESULT_NOT_EXISTS
		return response, nil
	}
	acc.AccountTag = in.Data.AccountTag
	acc.IP = in.Data.Ip
	acc.Port = in.Data.Port
	if err := onebot.DB.Save(&acc).Error; err != nil {
		return response, nil
	}
	response.Result.Result = auth.CheckResult_CHECK_RESULT_SUCCESS
	response.ModifyResult = accounts.ModifyAccountResult_MODIFY_ACCOUNT_RESULT_SUCCESS
	return response, nil
}
