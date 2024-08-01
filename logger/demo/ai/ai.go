package ai

import (
	"fmt"
	lgc "github.com/504dev/logr-go-client"
	"github.com/504dev/logr/config"
	"regexp"
	"strings"
	"time"
)

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

func Run(conf *lgc.Config) {
	log, _ := conf.NewLogger("ai.log")
	log.Body = "[{version}] {message}"

	ollama, err := NewOllamaChat(config.Get().DemoDash.Llm)
	if err != nil {
		log.Error(err)
		return
	}

	for {
		if err := generateBook(log, ollama); err != nil {
			log.Error(err)
		}
		time.Sleep(promptInterval)
	}
}

func generateBook(log *lgc.Logger, ollama *OllamaChat) error {
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

	answer, err := ollama.Prompt(history, func(s string) { log.Notice(s) }, onToken)
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

		answer, err := ollama.Prompt(history, func(s string) { log.Info(s) }, onToken)

		if err != nil {
			return err
		}

		history = append(history, answer)
		log.Info("")
	}

	return nil
}
