package slacknotifications

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "time"
)

type SlackRequestBody struct {
    Text string `json:"text"`
}

// method for sending the message to slack given the webhookurl
func SendSlackNotification(webhookUrl string, msg string) error {
  // create body for the request and trigger call to webhookurl
    slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
    req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
    if err != nil {
        return err
    }

    req.Header.Add("Content-Type", "application/json")

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }

    // fetch response and retrive error if any. 
    buf := new(bytes.Buffer)
    buf.ReadFrom(resp.Body)
    if buf.String() != "ok" {
        return errors.New("Non-ok response returned from Slack")
    }
    return nil
}
