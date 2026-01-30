package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var currentReq Request

func Parse() (Request, error) {
	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		return Request{}, err
	}

	if err := json.Unmarshal(inputData, &currentReq); err != nil {
		return Request{}, err
	}
	return currentReq, nil
}

func Log(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(os.Stderr, "USER_LOG: %s\n", msg)
}

func SendJSON(data interface{}) {
	body, _ := json.Marshal(data)
	send(200, "application/json", string(body))
}

func SendHTML(html string) {
	send(200, "text/html; charset=utf-8", html)
}

func SendError(status int, message string) {
	send(status, "application/json", fmt.Sprintf(`{"error": "%s"}`, message))
}

func send(status int, contentType string, body string) {
	resp := Response{
		Status: status,
		Headers: map[string][]string{
			"Content-Type": {contentType},
			"X-Powered-By": {"Gojinn SDK"},
		},
		Body: body,
	}
	json.NewEncoder(os.Stdout).Encode(resp)
}
