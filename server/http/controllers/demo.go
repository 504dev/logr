package controllers

import (
	"github.com/504dev/logr/libs/cachify"
	"github.com/504dev/logr/repo"
	"github.com/504dev/logr/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type DemoController struct {
	repos      *repo.Repos
	jwtService *types.JwtService
}

func NewDemoController(jwtService *types.JwtService, repos *repo.Repos) *DemoController {
	return &DemoController{
		repos:      repos,
		jwtService: jwtService,
	}
}

func (demo *DemoController) FreeToken(c *gin.Context) {
	usr, err := demo.repos.User.GetById(types.UserDemoId)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	tokenstring, err := cachify.Cachify("free-token", func() (interface{}, error) {
		claims := types.Claims{
			Id:       usr.Id,
			Role:     usr.Role,
			GihubId:  usr.GithubId,
			Username: usr.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			},
		}
		tokenstring, err := demo.jwtService.SignToken(&claims)
		if err != nil {
			return nil, err
		}
		return tokenstring, err
	}, 4*time.Minute)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, tokenstring)
}
