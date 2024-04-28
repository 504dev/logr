package logger

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	lgc "github.com/504dev/logr-go-client"
	"github.com/504dev/logr/config"
	"net/http"
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

func Prompt(history ChatHistory, onSentence func(string)) (*ChatHistoryItem, error) {
	OLLAMA_MODEL := config.Get().DemoDash.Model
	OLLAMA_CHAT_URL := config.Get().DemoDash.Url

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(RequestBody{
		Model:    OLLAMA_MODEL,
		Messages: history,
		Stream:   true,
	}); err != nil {
		return nil, err
	}

	resp, err := http.Post(OLLAMA_CHAT_URL, "application/json", &buf)
	if err != nil {
		return nil, err
	}

	answer := ChatHistoryItem{Role: "assistant"}
	tmp := ""

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var item ResponseBody
		if err := json.Unmarshal(scanner.Bytes(), &item); err != nil {
			panic(err)
		}
		answer.Content += item.Message.Content
		tmp += item.Message.Content
		splitted := splitIntoSentences(tmp)
		if len(splitted) > 1 {
			onSentence(splitted[0])
			tmp = splitted[1]
		}
	}
	onSentence(tmp)

	return &answer, nil
}

func author(conf *lgc.Config) {
	defer func() {
		<-time.After(10 * time.Second)
		author(conf)
	}()

	log, _ := conf.NewLogger("author.log")
	log.Body = "[{version}, pid={pid}] {message}"

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

	answer, err := Prompt(history, func(s string) {
		log.Notice(s)
	})
	if err != nil {
		log.Error(err)
		return
	}

	history = append(history, answer)

	for i := 1; i <= n; i++ {
		prompt := fmt.Sprintf(`Write chapter %v, as long as you can.`, i)
		history = append(history, &ChatHistoryItem{Role: "user", Content: prompt})

		answer, err := Prompt(history, func(s string) {
			log.Info(s)
		})

		if err != nil {
			log.Error(err)
			return
		}

		history = append(history, answer)
	}
}

func splitIntoSentences(text string) []string {
	var sentences []string
	var sentence strings.Builder
	skip := false
	for i, r := range text {
		if skip {
			skip = false
			continue
		}
		var next byte
		if i+1 < len(text) {
			next = text[i+1]
		}
		switch r {
		case '.', '?', '!':
			sentence.WriteRune(r)
			if (next == ' ' || next == '\n') && sentence.Len() > 32 {
				sentences = append(sentences, sentence.String())
				sentence.Reset()
				skip = true
			}
		case '\n':
			sentences = append(sentences, sentence.String())
			sentence.Reset()
		default:
			sentence.WriteRune(r)
		}
	}
	if sentence.Len() > 0 {
		sentences = append(sentences, sentence.String())
	}
	return sentences
}
