package send_telegram

import (
	"fmt"
	"log"
	"net/http"
)

func SendTelegram(bot_token string, chat_id string, text string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", bot_token, chat_id, text)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
}
