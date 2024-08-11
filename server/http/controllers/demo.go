package controllers

import (
	"github.com/504dev/logr/cachify"
	"github.com/504dev/logr/models/user"
	"github.com/504dev/logr/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type DemoController struct {
	jwtService *types.JwtService
}

func NewDemoController(jwtService *types.JwtService) *DemoController {
	return &DemoController{jwtService: jwtService}
}

func (d *DemoController) FreeToken(c *gin.Context) {
	usr, err := user.GetById(types.UserDemoId)
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
		tokenstring, err := d.jwtService.SignToken(&claims)
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
