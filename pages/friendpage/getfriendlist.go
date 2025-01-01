package friendpage

import (
	"context"
	"log"

	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	friend "github.com/Tensorix/metahub-backend-service/gen/proto/v1/friend"
	"github.com/Tensorix/metahub-backend-service/onebot"
	"github.com/Tensorix/metahub-backend-service/pages/authpage"
)

func (s *server) GetFriendList(_ context.Context, in *friend.FriendListRequest) (*friend.FriendListResponse, error) {
	response := &friend.FriendListResponse{
		Result: auth.CheckResult_CHECK_RESULT_FAILED,
	}
	token := in.Token
	username := authpage.GetUsername(token)
	if username == "" {
		return response, nil
	}
	response.Result = auth.CheckResult_CHECK_RESULT_SUCCESS
	for _, bot := range onebot.Bots {
		if bot.Username != username {
			continue
		}
		fl := &friend.FriendList{
			AccountTag: bot.AccountTag,
		}
		friends, err := bot.GetFriendList()
		if err != nil {
			log.Println(err)
			continue
		}
		for _, f := range friends {
			fl.Friends = append(fl.Friends, &friend.Friend{
				UserId:   f.Id,
				Uid:      f.UID,
				Nickname: f.Nickname,
				Remark:   f.Remark,
			})
		}
		response.FriendList = append(response.FriendList, fl)
	}
	return response, nil
}
