package ai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type OllamaChat struct {
	Url   string
	Model string
}

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

func NewOllamaChat(url string) (*OllamaChat, error) {
	matches := urlRegex.FindStringSubmatch(url)
	if len(matches) != 3 {
		return nil, errors.New("invalid llm url")
	}
	return &OllamaChat{
		Url:   matches[1],
		Model: matches[2],
	}, nil
}

func (ollama *OllamaChat) Prompt(
	history ChatHistory,
	onSentence func(string),
	onToken func(string),
) (*ChatHistoryItem, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(RequestBody{
		Model:    ollama.Model,
		Messages: history,
		Stream:   true,
	}); err != nil {
		return nil, err
	}

	resp, err := http.Post(ollama.Url+"/api/chat", "application/json", &buf)
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
