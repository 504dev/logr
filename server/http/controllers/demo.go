package controllers

import (
	"github.com/504dev/logr/libs/cachify"
	"github.com/504dev/logr/repo"
	"github.com/504dev/logr/types"
	"github.com/504dev/logr/types/jwtservice"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type DemoController struct {
	repos      *repo.Repos
	jwtService *jwtservice.JwtService
}

func NewDemoController(jwtService *jwtservice.JwtService, repos *repo.Repos) *DemoController {
	return &DemoController{
		repos:      repos,
		jwtService: jwtService,
	}
}

func (demo *DemoController) FreeToken(c *gin.Context) {
	usr, err := demo.repos.User.GetById(types.USER_DEMO_ID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	const tokenLifeTime = 15 * time.Minute
	const cacheTime = 4 * time.Minute
	tokenstring, err := cachify.Cachify("free-token", func() (interface{}, error) {
		claims := jwtservice.Claims{
			Id:       usr.Id,
			Role:     usr.Role,
			GihubId:  usr.GithubId,
			Username: usr.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifeTime)),
			},
		}
		tokenstring, err := demo.jwtService.SignToken(&claims)
		if err != nil {
			return nil, err
		}

		return tokenstring, err
	}, cacheTime)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, tokenstring)
}
