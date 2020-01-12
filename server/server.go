package server

import (
	"github.com/504dev/kidlog/config"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

type LogHandler struct{}

func (t LogHandler) Write(b []byte) (int, error) {
	return len(b), nil
}

func Init() {
	gin.ForceConsoleColor()

	gin.DefaultWriter = io.MultiWriter(os.Stdout, LogHandler{})

	r := NewRouter()

	r.Run(config.Get().Bind.Http)
}
