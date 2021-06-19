package slack

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type SlackMessage struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
}

func GetKST(t time.Time) (kst time.Time) {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		fmt.Println("[ERROR] Failed to load location")
		panic(err.Error())
	}
	return t.In(loc)
}

func SetText(key types.AccessKeyMetadata, channel string) (b []byte) {
	text := fmt.Sprintf(`Access key expired!
	IAM User : %s
	Access Key ID : %s
	Create Date : %s`, *key.UserName, *key.AccessKeyId, GetKST(*key.CreateDate))

	serializedMessage, err := json.Marshal(SlackMessage{
		Text:    text,
		Channel: channel,
	})
	if err != nil {
		fmt.Println("[ERROR] Failed to marshal slack message")
		panic(err.Error())
	}
	return serializedMessage
}
