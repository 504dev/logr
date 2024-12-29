package tests

import (
	"fmt"
	logr "github.com/504dev/logr-go-client"
	"github.com/504dev/logr/server"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLogServer(t *testing.T) {
	cfg := logr.Config{}
	logger, _ := cfg.NewLogger("test.log")
	srv, err := server.NewLogServer(
		"",
		"",
		"",
		"",
		func() string {
			return ""
		},
		nil,
		logger,
	)

	assert.Nil(t, err, "should not return an error")

	go srv.Run()

	time.Sleep(time.Second)
	fmt.Println(srv.Stop())
}
