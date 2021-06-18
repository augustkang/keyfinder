package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type KeyFinder struct {
	Hours        int
	URL          string
	Client       *iam.Client
	SlackChannel string
}

type SlackMessage struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
}

func NewKeyFinder(hours int, url string, channel string) (kf *KeyFinder) {
	return &KeyFinder{
		Hours:        hours,
		URL:          url,
		SlackChannel: channel,
	}
}

func (kf *KeyFinder) SetIAMClient() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("[ERROR] Failed to load config")
		panic(err.Error())
	}
	kf.Client = iam.NewFromConfig(cfg)
}

func (kf *KeyFinder) GetUserNames() (userNames []string) {

	var names []string
	paginator := iam.NewListUsersPaginator(kf.Client, &iam.ListUsersInput{
		PathPrefix: aws.String("/"),
	})

	for paginator.HasMorePages() {
		listUserOutput, err := paginator.NextPage(context.TODO())
		if err != nil {
			panic(err.Error())
		}
		for _, user := range listUserOutput.Users {
			names = append(names, *user.UserName)
		}
	}
	return names
}

func (kf *KeyFinder) GetKeyList(userNames []string) (keyList []types.AccessKeyMetadata) {

	var userKeyList []types.AccessKeyMetadata

	for _, name := range userNames {
		keyPaginator := iam.NewListAccessKeysPaginator(kf.Client, &iam.ListAccessKeysInput{
			UserName: aws.String(name),
		})
		for keyPaginator.HasMorePages() {
			output, err := keyPaginator.NextPage(context.TODO())
			if err != nil {
				fmt.Println("[ERROR] Failed to get NextPage of keyPaginator in listKeys")
				panic(err.Error())
			}
			userKeyList = append(userKeyList, output.AccessKeyMetadata...)
		}

	}
	return userKeyList
}

func (kf *KeyFinder) CheckKeyAges(allKeyList []types.AccessKeyMetadata) {

	var count int
	for _, key := range allKeyList {
		expireTime := time.Now().UTC().Add(time.Hour * time.Duration(-kf.Hours))

		// If
		if expireTime.Sub(*key.CreateDate) > 0 {
			count += 1
			kf.PostSlackMessage(key)
			fmt.Println("[RESULT] Sent Slack Message regarding IAM User using expired Access Key Pair : ", *key.UserName)
		}
	}
	if count == 0 {
		fmt.Println("[RESULT] There was no expired Access Keys.")
	}
}

func (kf *KeyFinder) PostSlackMessage(key types.AccessKeyMetadata) {

	msg := fmt.Sprintf(`Access key expired!
	IAM User : %s
	Access Key ID : %s
	Create Date : %s`, *key.UserName, *key.AccessKeyId, *key.CreateDate)

	message, err := json.Marshal(SlackMessage{
		Text:    msg,
		Channel: kf.SlackChannel,
	})

	if err != nil {
		fmt.Println("[ERROR] Failed to marshal slack message")
		panic(err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, kf.URL, bytes.NewBuffer(message))
	if err != nil {
		fmt.Println("[ERROR] Failed to create http request")
		panic(err.Error())
	}

	req.Header.Add("Content-Type", "application/json")

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("[ERROR] Failed to get response from Slack")
		panic(err.Error())
	}

	buf := new(bytes.Buffer)

	buf.ReadFrom(resp.Body)

	if buf.String() != "ok" {
		panic(err.Error())
	}

}
