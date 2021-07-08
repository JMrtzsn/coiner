package notification

import (
	"fmt"
	"net/http"
	"net/url"

)

type Telegram struct {
	ID        string
	Key       string
	ChannelID string

}

func NewTelegram(id string, key string, channel string) Telegram {
	return Telegram{
		ID:        id,
		Key:       key,
		ChannelID: channel,
	}
}

func (t Telegram) Notify(text string) error {
	baseURL, err := url.Parse(fmt.Sprintf("https://api.telegram.org/bot%s:%s/sendMessage", t.ID, t.Key))
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Add("chat_id", t.ChannelID)
	params.Add("text", text)

	baseURL.RawQuery = params.Encode()
	resp, err := http.Get(baseURL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("notification/telegram: %d - %s", resp.StatusCode, resp.Status)
	}
	return nil
}

func (t Telegram) OnError(err error) {
	title := "🛑 ERROR"
	message := fmt.Sprintf("%s\n-----\n%s", title, err)
	t.Notify(message)
}
