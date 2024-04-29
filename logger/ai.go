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

func ai(conf *lgc.Config) {
	defer func() {
		<-time.After(10 * time.Second)
		ai(conf)
	}()

	log, _ := conf.NewLogger("ai.log")
	log.Body = "[{version}] {message}"

	n := 5
	genres := []string{
		"Romance",
		"Biographies & Memoirs",
		"Literary Fiction",
		"Science Fiction",
		"Drama",
		"Horror",
		"Mysticism",
		"Phantasmagoria",
		"Detective",
		"Novel",
		"Fantasy",
		"Adventure",
		"Comedy",
		"Manga and anime",
		"Young Adult",
		"Self-Help",
	}
	genre := genres[time.Now().Nanosecond()%len(genres)]
	prompt := fmt.Sprintf(`Imagine that you are a writer in the %v genre.
Think of the title of a book about a monitoring service called "Logr", which was developed by a 30-year-old developer from Saint-Petersburg named Dima.
Then state the genre of the book.
Then make a table of contents of %v chapters.
Then write a 100-word summary of the book.`, genre, n)
	history := ChatHistory{
		{Role: "user", Content: prompt},
	}

	onToken := func(t string) { log.Inc("tokens", 1) }

	answer, err := Prompt(history, func(s string) { log.Notice(s) }, onToken)
	if err != nil {
		log.Error(err)
		return
	}

	history = append(history, answer)

	log.Info("")
	for i := 1; i <= n; i++ {
		prompt := fmt.Sprintf(`Write a chapter %v that is 500 words long.`, i)
		if i == n {
			prompt += "This is the last chapter, ending the book epically."
		}
		history = append(history, &ChatHistoryItem{Role: "user", Content: prompt})

		answer, err := Prompt(history, func(s string) { log.Info(s) }, onToken)

		if err != nil {
			log.Error(err)
			return
		}

		history = append(history, answer)
		log.Info("")
	}
}

func Prompt(history ChatHistory, onSentence func(string), onToken func(string)) (*ChatHistoryItem, error) {
	OLLAMA_MODEL := config.Get().DemoDash.Model
	OLLAMA_CHAT_URL := config.Get().DemoDash.Url + "/api/chat"

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
		if onToken != nil {
			onToken(item.Message.Content)
		}
		answer.Content += item.Message.Content
		tmp += item.Message.Content

		if splitted := splitIntoSentences(tmp); len(splitted) > 1 {
			onSentence(splitted[0])
			tmp = splitted[1]
		}
	}
	onSentence(tmp)

	return &answer, nil
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
