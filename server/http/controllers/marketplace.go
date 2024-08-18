package controllers

import (
	"encoding/json"
	"errors"
	"github.com/504dev/logr/config"
	. "github.com/504dev/logr/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httputil"
	"time"
)

type MarketplaceController struct{}

func (*MarketplaceController) Webhook(ctx *gin.Context) {
	requestDump, err := httputil.DumpRequest(ctx.Request, true)
	if err != nil {
		Logger.Error(err)
	}
	Logger.Notice(string(requestDump))
	ctx.AbortWithStatus(http.StatusOK)
}

func (*MarketplaceController) Support(c *gin.Context) {
	var data struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Message string `json:"message"`
		Token   string `json:"recaptchaToken,omitempty"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	verifyData, err := checkRecaptcha(c, config.Get().RecaptchaSecret, data.Token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	data.Token = ""
	payload, _ := json.Marshal(data)
	Logger.Notice("%v %v", string(payload), verifyData)
	c.AbortWithStatus(http.StatusOK)
}

type siteVerifyResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

func checkRecaptcha(ctx context.Context, secret, response string) (*siteVerifyResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://www.google.com/recaptcha/api/siteverify", nil)
	if err != nil {
		return nil, err
	}

	// Add necessary request parameters.
	q := req.URL.Query()
	q.Add("secret", secret)
	q.Add("response", response)
	req.URL.RawQuery = q.Encode()

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response.
	var body siteVerifyResponse
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	// Check recaptcha verification success.
	if !body.Success {
		return &body, errors.New("unsuccessful recaptcha verify request")
	}

	return &body, nil
}
