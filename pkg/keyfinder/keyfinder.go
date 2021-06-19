package keyfinder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"

	"augustkang.com/keyfinder/pkg/slack"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type KeyFinder struct {
	Hours        int
	URL          string
	Client       *iam.Client
	SlackChannel string
}

// NewKeyFinder returns KeyFinder struct
func NewKeyFinder(hours int, url string, client *iam.Client, channel string) (kf *KeyFinder) {
	return &KeyFinder{
		Hours:        hours,
		URL:          url,
		Client:       client,
		SlackChannel: channel,
	}
}

// GetUserNames retrieves all IAM Users and return user names as a slice
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

// GetKeyList retrieves all Access Key Pairs with given IAM User's name, then return AccessKeyMetadata as a slice
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

// CheckKeyAges compare AccessKeyMetadata.CreateDate to kf.Hours
func (kf *KeyFinder) CheckKeyAges(allKeyList []types.AccessKeyMetadata) {
	var count int
	for _, key := range allKeyList {
		expireTime := time.Now().UTC().Add(time.Hour * time.Duration(-kf.Hours))

		// If key expired then send message to Slack Webhook
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

// PostSlackMessage sends HTTP POST request to kf.URL
func (kf *KeyFinder) PostSlackMessage(key types.AccessKeyMetadata) {
	msg := fmt.Sprintf(`Access key expired!
	IAM User : %s
	Access Key ID : %s
	Create Date : %s`, *key.UserName, *key.AccessKeyId, *key.CreateDate)

	message, err := json.Marshal(slack.SlackMessage{
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
