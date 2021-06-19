package slack

type SlackMessage struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
}
