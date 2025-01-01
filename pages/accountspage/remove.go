package accountspage

import (
	"context"
	"log"

	accounts "github.com/Tensorix/metahub-backend-service/gen/proto/v1/accounts"
	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	"github.com/Tensorix/metahub-backend-service/onebot"
	"github.com/Tensorix/metahub-backend-service/pages/authpage"
)

func (server *server) RemoveAccount(_ context.Context, in *accounts.RemoveAccountRequest) (*accounts.RemoveAccountResponse, error) {
	var user authpage.User
	var acc Account
	token := in.Token.Token
	response := &accounts.RemoveAccountResponse{
		Result: &auth.CheckResponse{
			Result: auth.CheckResult_CHECK_RESULT_FAILED,
		},
		RemoveResult: accounts.RemoveAccountResult_REMOVE_ACCOUNT_RESULT_FAILED,
	}
	username := authpage.GetUsername(token)
	if username == "" {
		return response, nil
	}
	response.Result.Result = auth.CheckResult_CHECK_RESULT_SUCCESS
	if err := onebot.DB.First(&user, "username = ?", username).Error; err != nil {
		log.Println(err)
		response.RemoveResult = accounts.RemoveAccountResult_REMOVE_ACCOUNT_RESULT_FAILED
		return response, nil
	}
	if err := onebot.DB.First(&acc, "id = ? AND user_id = ?", in.Id, user.Id).Error; err != nil {
		response.RemoveResult = accounts.RemoveAccountResult_REMOVE_ACCOUNT_RESULT_NOT_EXISTS
		return response, nil
	}
	if err := onebot.DB.Delete(&acc).Error; err != nil {
		response.RemoveResult = accounts.RemoveAccountResult_REMOVE_ACCOUNT_RESULT_FAILED
		return response, nil
	}
	response.RemoveResult = accounts.RemoveAccountResult_REMOVE_ACCOUNT_RESULT_SUCCESS
	newbots := make([]*onebot.Onebot, 0)
	for _, bot := range onebot.Bots {
		if bot.AccountID == acc.ID {
			bot.Shutdown()
			continue
		}
		newbots = append(newbots, bot)
	}
	onebot.Bots = newbots
	return response, nil
}
