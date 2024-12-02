package friendpage

import (
	"context"
	"log"

	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	friend "github.com/Tensorix/metahub-backend-service/gen/proto/v1/friend"
	"github.com/Tensorix/metahub-backend-service/pages/authpage"
)

func (s *server) GetFriendList(_ context.Context, in *friend.FriendListRequest) (*friend.FriendListResponse, error) {
	friendList := []*friend.FriendList{}
	response := &friend.FriendListResponse{
		Result: auth.CheckResult_CHECK_RESULT_FAILED,
	}
	token := in.Token
	username := authpage.GetUsername(token)
	if username == "" {
		return response, nil
	}
	response.Result = auth.CheckResult_CHECK_RESULT_SUCCESS
	for _, bot := range bots {
		if bot == nil {
			break
		}
		if bot.Username != username {
			continue
		}
		fl, err := bot.GetFriendList()
		if err != nil {
			log.Println(err)
			break
		}
		friends := &friend.FriendList{
			AccountTag: bot.AccountTag,
			Friends:    []*friend.Friend{},
		}
		for _, f := range fl.Data {
			friends.Friends = append(friends.Friends, &friend.Friend{
				UserId:   f.UserID,
				Nickname: f.Nickname,
				Remark:   f.Remark,
			})
		}
		friendList = append(friendList, friends)
	}
	response.FriendList = friendList
	return response, nil
}
