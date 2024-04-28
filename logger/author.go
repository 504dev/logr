package logger

import (
	"fmt"
	lgc "github.com/504dev/logr-go-client"
	"github.com/504dev/logr/config"
	"github.com/go-resty/resty/v2"
	"strings"
	"time"
)

type ChatHistoryItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatHistory []*ChatHistoryItem
type RequestBody struct {
	Model    string      `json:"model"`
	Messages ChatHistory `json:"messages"`
	Stream   bool        `json:"stream"`
}
type ResponseBody struct {
	Model   string          `json:"model"`
	Message ChatHistoryItem `json:"message"`
	Done    bool            `json:"done"`
}

func author(conf *lgc.Config) {
	defer func() {
		<-time.After(10 * time.Second)
		author(conf)
	}()
	OLLAMA_MODEL := config.Get().DemoDash.Model
	OLLAMA_CHAT_URL := config.Get().DemoDash.Url
	log, _ := conf.NewLogger("author.log")
	n := 5
	genres := []string{
		"science fiction",
		"drama",
		"horror",
		"mysticism",
		"phantasmagoria",
		"detective",
		"novel",
		"fantasy",
		"adventure",
		"comedy",
	}
	genre := genres[time.Now().Nanosecond()%len(genres)]
	prompt := fmt.Sprintf(`Imagine that you are a writer in the %s genre.
Think of the title of a book about a monitoring service called logr, which was developed by a 30-year-old developer from St. Petersburg named Dima.
Then state the genre of the book.
Then make a table of contents of %s short chapter titles.
Then write a 100-word summary of the book.`, genre, n)
	history := ChatHistory{
		{Role: "user", Content: prompt},
	}

	var body ResponseBody

	client := resty.New()
	_, err := client.R().
		SetBody(RequestBody{
			Model:    OLLAMA_MODEL,
			Messages: history,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&body).
		Post(OLLAMA_CHAT_URL)

	if err != nil {
		log.Error(err)
		return
	}

	history = append(history, &body.Message)

	log.Notice(body.Message.Content)

	for i := 1; i <= n; i++ {
		prompt := fmt.Sprintf(`Write chapter %v, as long as you can.`, i)
		history = append(history, &ChatHistoryItem{Role: "user", Content: prompt})

		var body ResponseBody
		_, err := client.R().
			SetBody(RequestBody{
				Model:    OLLAMA_MODEL,
				Messages: history,
			}).
			SetHeader("Accept", "application/json").
			SetResult(&body).
			Post(OLLAMA_CHAT_URL)
		if err != nil {
			log.Error(err)
			return
		}
		history = append(history, &body.Message)

		for _, chunk := range splitIntoSentences(body.Message.Content) {
			log.Info(chunk)
			<-time.After(time.Second * 3)
		}
	}
}
func splitIntoSentences(text string) []string {
	var sentences []string
	var sentence strings.Builder
	for i, r := range text {
		switch r {
		case '.', '?', '!':
			sentence.WriteRune(r)
			if i+1 < len(text) && text[i+1] == ' ' {
				sentences = append(sentences, strings.Trim(sentence.String(), " "))
				sentence.Reset()
			}
		case '\n':
			if sentence.Len() > 0 {
				sentences = append(sentences, strings.Trim(sentence.String(), " "))
				sentence.Reset()
			}
		default:
			sentence.WriteRune(r)
		}
	}
	if sentence.Len() > 0 {
		sentences = append(sentences, strings.Trim(sentence.String(), " "))
	}
	return sentences
}
