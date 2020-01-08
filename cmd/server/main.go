package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/504dev/kidlog/clickhouse"
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/models/dashboard"
	"github.com/504dev/kidlog/models/log"
	"github.com/504dev/kidlog/models/user"
	"github.com/504dev/kidlog/mysql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type LogHandler struct{}

func (t LogHandler) Write(b []byte) (int, error) {
	return len(b), nil
}

func main() {
	config.Init()
	clickhouse.Init()
	mysql.Init()

	gin.ForceConsoleColor()

	gin.DefaultWriter = io.MultiWriter(os.Stdout, LogHandler{})

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/logs", func(c *gin.Context) {
		dashid, _ := strconv.Atoi(c.Query("dash_id"))
		if dashid == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "day required"})
			return
		}

		logname := c.Query("logname")
		hostname := c.Query("hostname")
		message := c.Query("message")
		level, _ := strconv.Atoi(c.Query("level"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		offset, _ := strconv.ParseInt(c.Query("offset"), 10, 64)

		from, _ := strconv.ParseInt(c.Query("timestamp.from"), 10, 64)
		to, _ := strconv.ParseInt(c.Query("timestamp.to"), 10, 64)

		where := log.Filter{
			Timestamp: [2]int64{from, to},
			DashId:    dashid,
			Logname:   logname,
			Hostname:  hostname,
			Level:     level,
			Message:   message,
			Offset:    offset,
			Limit:     limit,
		}
		fmt.Println(where)

		logs, err := log.GetAll(where)
		fmt.Println(err)
		c.JSON(200, logs)
	})
	r.GET("/dashboards", func(c *gin.Context) {
		dashboards, _ := dashboard.GetAll()
		c.JSON(200, dashboards)
	})
	r.GET("/dashboard/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		dash, _ := dashboard.GetById(id)
		c.JSON(200, dash)
	})
	r.GET("/users", func(c *gin.Context) {
		users, _ := user.GetAll()
		c.JSON(200, users)
	})
	r.GET("/user/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		usr, _ := user.GetById(id)
		c.JSON(200, usr)
	})

	{
		r.GET("/me", func(c *gin.Context) {
			usr, _ := user.GetById(1)
			c.JSON(200, usr)
		})

		r.GET("/my/dashboards", func(c *gin.Context) {
			dashboards, _ := dashboard.GetAll()
			c.JSON(200, dashboards)
		})
	}

	{
		CLIENT_ID := "a3e0eabef800cd0e7a84"
		CLIENT_SECRET := "95344c1682df6e82e71652398dcf9f44b1c6ed8d"
		SCOPE := "user"
		STATE := "secretstring"
		REDIRECT_URL := "http://localhost:8080/"

		r.GET("/oauth/signin", func(c *gin.Context) {
			params := url.Values{}
			params.Add("client_id", CLIENT_ID)
			params.Add("state", STATE)
			params.Add("scope", SCOPE)
			redirectUrl := "https://github.com/login/oauth/authorize?" + params.Encode()
			c.Redirect(http.StatusMovedPermanently, redirectUrl)
			c.Abort()
		})
		r.GET("/oauth/callback", func(c *gin.Context) {
			state := c.Query("state")
			code := c.Query("code")
			requestBody, _ := json.Marshal(map[string]string{
				"client_id":     CLIENT_ID,
				"client_secret": CLIENT_SECRET,
				"code":          code,
				"state":         state,
			})

			url := "https://github.com/login/oauth/access_token"
			resp, _ := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(body))

			c.Redirect(http.StatusMovedPermanently, REDIRECT_URL)
			c.Abort()
		})
	}

	r.Run(config.Get().Bind.Http)
}
