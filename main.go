package main

import (
	"fmt"
	"net/http"
	"strconv"

	"augustkang.com/keyfinder/pkg"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
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

	keyFinder := pkg.NewKeyFinder(intHours, url, channel)

	keyFinder.SetIAMClient()

	keyNames := keyFinder.GetUserNames()

	keyList := keyFinder.GetKeyList(keyNames)

	keyFinder.CheckKeyAges(keyList)

	fmt.Fprintf(w, "Result sent to Slack %s Channel!\n", channel)

}

func main() {
	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(":8080", nil)
}
