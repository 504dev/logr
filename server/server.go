package server

import (
	"github.com/504dev/logr/config"
	. "github.com/504dev/logr/logger"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func ListenHTTP() error {
	gin.ForceConsoleColor()

	gin.DefaultWriter = io.MultiWriter(os.Stdout, GinWritter)

	r := NewRouter()

	return r.Run(config.Get().Bind.Http)
}
