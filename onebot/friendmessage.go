package onebot

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type BotFriendMessage struct {
	MessageID int64 `json:"message_id"`
	UserID    int64 `json:"user_id"`
	Message   []struct {
		Type string `json:"type"`
		Data struct {
			Text     string `json:"text"`
			File     string `json:"file"`
			Filename string `json:"filename"`
			URL      string `json:"url"`
			Summary  string `json:"summary"`
			SubType  int    `json:"subType"`
		} `json:"data"`
	} `json:"message"`
	Time uint32 `json:"time"`
}

type FriendMessage struct {
	ID          uint
	MessageID   int64
	FriendID    int64
	MessageTS   uint32
	SelfMessage bool
	ReadMark    bool
	Hide        bool
	Revoke      bool
}

type FriendSubMessage struct {
	ID              uint
	FriendMessageID uint
	IsText          bool
	Message         string
}

func (bot *Onebot) friendMessage() error {
	var botmsg BotFriendMessage
	err := json.Unmarshal(bot.message, &botmsg)
	if err != nil {
		log.Println(err)
		return nil
	}
	var friend Friend
	if err := DB.First(&friend, "uid = ? AND account_id = ?", botmsg.UserID, bot.AccountID).Error; err != nil {
		return err
	}
	friendMessage := FriendMessage{
		MessageID:   botmsg.MessageID,
		FriendID:    friend.Id,
		MessageTS:   botmsg.Time,
		SelfMessage: false,
		ReadMark:    false,
		Hide:        false,
		Revoke:      false,
	}
	DB.Create(&friendMessage)
	var friendSubMessages []FriendSubMessage
	for i := 0; i < len(botmsg.Message); i++ {
		istext := true
		msg := botmsg.Message[i].Data.Text
		switch botmsg.Message[i].Type {
		case "text":
		case "image":
			istext = false
			filename, err := downloadFile(botmsg.Message[i].Data.URL)
			if err != nil {
				log.Println(err)
			}
			msg = filename
		}
		friendSubMessages = append(friendSubMessages, FriendSubMessage{
			FriendMessageID: friendMessage.ID,
			IsText:          istext,
			Message:         msg,
		})
	}
	DB.Create(&friendSubMessages)
	select {
	case <-bot.FriendMessage:
	default:
	}
	go func() {
		bot.FriendMessage <- struct{}{}
	}()
	return nil
}

func downloadFile(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:132.0) Gecko/20100101 Firefox/132.0")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	file, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	md5sum := md5.Sum(file)
	filename := hex.EncodeToString(md5sum[:])
	_, err = os.Stat(filename)
	if !os.IsNotExist(err){
		return filename, nil
	}
	err = os.WriteFile("cache/images/"+filename, file, 0644)
	if err != nil {
		return "", err
	}
	return filename, nil
}
