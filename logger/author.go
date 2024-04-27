package logger

import (
	"fmt"
	lgc "github.com/504dev/logr-go-client"
	"github.com/504dev/logr/config"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"time"
)

const LLMAPIURL = "https://api.coze.com/open_api/v2/chat"

type HistoryItem struct {
	Role    string `json:"role"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (hi *HistoryItem) String() string {
	return hi.Content
}

type ChatHistory []*HistoryItem
type RequestBody struct {
	ConversationId string      `json:"conversation_id"`
	BotId          string      `json:"bot_id"`
	User           string      `json:"user"`
	Query          string      `json:"query"`
	Stream         bool        `json:"stream"`
	ChatHistory    ChatHistory `json:"chat_history"`
}
type ResponseBody struct {
	Messages       ChatHistory `json:"messages"`
	ConversationId string      `json:"conversation_id"`
	Code           int         `json:"code"`
	Msg            string      `json:"msg"`
}

func (r *ResponseBody) Answer() *HistoryItem {
	return r.Messages[0]
}

func author(conf *lgc.Config) {
	defer func() {
		<-time.After(time.Second)
		author(conf)
	}()
	log, _ := conf.NewLogger("author.log")
	n := 5
	genres := []string{
		"фантастика",
		"драма",
		"ужасы",
		"мистика",
		"фантасмагория",
		"детектив",
		"роман",
		"фэнтези",
		"приключения",
	}
	history := ChatHistory{}
	prompt := fmt.Sprintf(`Представь что ты писатель в жанре %s. Пиши на русском.
Придумай название книги про сервис мониторинга под названием logr, который разработал 30 летний разработчик из Санкт-Петербурга по имени Дима.
Затем укажи жанр книги.
Затем составь оглавление из %s коротких названий глав.
Затем напиши краткое описание книги на 500 символов.`, genres[rand.Intn(len(genres))], n)

	var body ResponseBody

	client := resty.New()
	_, err := client.R().
		SetBody(RequestBody{
			BotId:       config.Get().DemoDash.BotId,
			User:        config.Get().DemoDash.UserId,
			Stream:      false,
			Query:       prompt,
			ChatHistory: history,
		}).
		SetAuthToken(config.Get().DemoDash.ApiKey).
		SetHeader("Accept", "application/json").
		SetResult(&body).
		Post(LLMAPIURL)

	if err != nil {
		log.Error(err)
		return
	}

	history = append(history, &HistoryItem{
		Role:    "user",
		Type:    "text",
		Content: prompt,
	})
	history = append(history, body.Answer())

	log.Notice(body.Answer())

	for i := 1; i <= n; i++ {
		prompt := fmt.Sprintf("Напиши в одном сообщении Главу %v книги", i)
		var body ResponseBody
		_, err := client.R().
			SetBody(RequestBody{
				BotId:       config.Get().DemoDash.BotId,
				User:        config.Get().DemoDash.UserId,
				Stream:      false,
				Query:       prompt,
				ChatHistory: history,
			}).
			SetAuthToken(config.Get().DemoDash.ApiKey).
			SetHeader("Accept", "application/json").
			SetResult(&body).
			Post(LLMAPIURL)
		if err != nil {
			log.Error(err)
			return
		}
		history = append(history, body.Answer())
		log.Info(body.Answer())
	}
}
