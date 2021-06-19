package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"augustkang.com/keyfinder/pkg/awsclient"
	"augustkang.com/keyfinder/pkg/keyfinder"
)

// RootHandler handles request to root(/)
func RootHandler(w http.ResponseWriter, r *http.Request) {

	hours := r.URL.Query().Get("hours")
	url := r.URL.Query().Get("url")
	channel := r.URL.Query().Get("channel")

	if url == "" {
		fmt.Fprintf(w, "Parameter 'url' is missing. Please specify slack webhook url.\nURL Request example\nlocalhost:8080/?hours=10&url=https://hooks.slack.com/services/abcdefg/hijklmnop/abcdefghijklmnopqrstuvwxyz")
		return
	}
	if hours == "" {
		fmt.Fprintf(w, "Parameter 'hours' is missing. Please specify slack webhook url.\nURL Request example\nlocalhost:8080/?hours=10&url=https://hooks.slack.com/services/abcdefg/hijklmnop/abcdefghijklmnopqrstuvwxyz")
		return
	}
	if channel == "" {
		channel = "#example"
	}

	intHours, err := strconv.Atoi(hours)
	if err != nil {
		fmt.Println("[FAIL] Failed to convert hours query string to int")
		fmt.Println(err.Error())
	}
	client := awsclient.GetIAMClient()

	kf := keyfinder.NewKeyFinder(intHours, url, client, channel)
	keyNames := kf.GetUserNames()
	keyList := kf.GetKeyList(keyNames)
	kf.CheckKeyAges(keyList)
	fmt.Fprintf(w, "Result sent to Slack %s Channel!\n", channel)

}
