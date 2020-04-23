package util

import (
	"net/http"
	"strings"
)

func SendNotify(msg string) bool {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://notify-api.line.me/api/notify", strings.NewReader(string("message=")+msg))
	if err != nil {
		return false
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer ClriI9Me3mXuMcl4Z9gb0R51JkhdYkMoI9HTos58fTt")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false
	}
	return true
}
