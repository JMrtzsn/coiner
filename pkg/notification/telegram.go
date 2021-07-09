package notification

import (
	"fmt"
	"net/http"
	"net/url"
)

type Telegram struct {
	Bot       string
	ChannelID string
}

func NewTelegram(key string, channel string) Telegram {
	return Telegram{
		Bot:       key,
		ChannelID: channel,
	}
}

func (t Telegram) Notify(text string) error {
	baseURL, err := url.Parse(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.Bot))
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

func (t Telegram) OnError(err error) error {
	title := "🛑 ERROR"
	message := fmt.Sprintf("%s\n-----\n%s", title, err)
	err = t.Notify(message)
	if err != nil {
		return err
	}
	return nil
}
