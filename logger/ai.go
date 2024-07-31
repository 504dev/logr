package logger

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	lgc "github.com/504dev/logr-go-client"
	"github.com/504dev/logr/config"
	"net/http"
	"regexp"
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

var (
	urlRegex = regexp.MustCompile(`^(https?://[^/]+)/(.*)$`)
	genres   = []string{
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
	chaptersN        = 5
	promptInterval   = 10 * time.Second
	chapterWordCount = 300
)

func ai(conf *lgc.Config) {
	log, _ := conf.NewLogger("ai.log")
	log.Body = "[{version}] {message}"

	ollamaUrl, ollamaModel, err := parseLLMURL(config.Get().DemoDash.Llm)
	if err != nil {
		log.Error(err)
		return
	}

	for {
		if err := generateBook(log, ollamaUrl, ollamaModel); err != nil {
			log.Error(err)
		}
		time.Sleep(promptInterval)
	}
}

func parseLLMURL(url string) (string, string, error) {
	matches := urlRegex.FindStringSubmatch(url)
	if len(matches) != 3 {
		return "", "", errors.New("invalid llm url")
	}
	return matches[1], matches[2], nil
}

func generateBook(log *lgc.Logger, ollamaUrl, ollamaModel string) error {
	genre := genres[time.Now().Nanosecond()%len(genres)]
	prompt := fmt.Sprintf(`Imagine that you are a writer in the %v genre.
Think of the title of a book about a monitoring service called "Logr", which was developed by a 30-year-old developer from Saint-Petersburg named Dima.
Then state the genre of the book.
Then make a table of contents of %v chapters.
Then write a 100-word summary of the book.`, genre, chaptersN)
	history := ChatHistory{
		{Role: "user", Content: prompt},
	}

	onToken := func(token string) {
		log.Inc("tokens", 1)
		token = strings.TrimSpace(token)
		if token == "Dima" {
			log.Inc(token, 1)
		}
	}

	answer, err := Prompt(ollamaUrl, ollamaModel, history, func(s string) { log.Notice(s) }, onToken)
	if err != nil {
		return err
	}

	history = append(history, answer)

	log.Info("")
	for i := 1; i <= chaptersN; i++ {
		prompt := fmt.Sprintf(`Give me a %v-word chapter %v. Start with the title of the chapter.`, chapterWordCount, i)
		if i == chaptersN {
			prompt += " This is the last chapter, ending the book epically."
		}
		history = append(history, &ChatHistoryItem{Role: "user", Content: prompt})

		answer, err := Prompt(ollamaUrl, ollamaModel, history, func(s string) { log.Info(s) }, onToken)

		if err != nil {
			return err
		}

		history = append(history, answer)
		log.Info("")
	}

	return nil
}

func Prompt(ollamaUrl string, ollamaModel string, history ChatHistory, onSentence func(string), onToken func(string)) (*ChatHistoryItem, error) {

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(RequestBody{
		Model:    ollamaModel,
		Messages: history,
		Stream:   true,
	}); err != nil {
		return nil, err
	}

	resp, err := http.Post(ollamaUrl+"/api/chat", "application/json", &buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	answer := ChatHistoryItem{Role: "assistant"}
	tmp := ""

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var item ResponseBody
		if err := json.Unmarshal(scanner.Bytes(), &item); err != nil {
			return nil, err
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
