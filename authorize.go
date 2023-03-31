package pocketbase

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/singleflight"
)

type authorizer interface {
	authorize() error
	refresh() error
}

type authorizeNoOp struct{}

func (a authorizeNoOp) authorize() error {
	return nil
}

func (a authorizeNoOp) refresh() error {
	return nil
}

type authorizeEmailPassword struct {
	email       string
	password    string
	token       string
	tokenValid  time.Time
	client      *resty.Client
	url         string
	tokenSingle singleflight.Group
}

func newAuthorizeEmailPassword(c *resty.Client, url string, email string, password string) authorizer {
	return &authorizeEmailPassword{
		client:      c,
		email:       email,
		password:    password,
		url:         url,
		tokenSingle: singleflight.Group{},
	}
}

func newAuthorizationRefresh(c *resty.Client, url string, token string) authorizer {
	return &authorizeEmailPassword{
		client:      c,
		token:       token,
		url:         url,
		tokenSingle: singleflight.Group{},
	}
}

func (a *authorizeEmailPassword) authorize() error {
	type authResponse struct {
		Token string `json:"token"`
	}

	_, err, _ := a.tokenSingle.Do("auth", func() (interface{}, error) {
		if time.Now().Before(a.tokenValid) {
			return nil, nil
		}

		resp, err := a.client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]interface{}{
				"identity": a.email,
				"password": a.password,
			}).
			SetResult(&authResponse{}).
			SetHeader("Authorization", "").
			Post(a.url)

		if err != nil {
			return nil, fmt.Errorf("[auth] can't send request to pocketbase %w", err)
		}

		if resp.IsError() {
			return nil, fmt.Errorf("[auth] pocketbase returned status: %d, msg: %s, err %w",
				resp.StatusCode(),
				resp.String(),
				ErrInvalidResponse,
			)
		}

		auth := *resp.Result().(*authResponse)
		a.token = auth.Token
		a.client.SetHeader("Authorization", auth.Token)
		a.tokenValid = time.Now().Add(60 * time.Minute)

		return nil, nil
	})
	return err
}

func (a *authorizeEmailPassword) refresh() error {
	type authResponse struct {
		Token string `json:"token"`
	}
	_, err, _ := a.tokenSingle.Do("auth-refresh", func() (interface{}, error) {
		resp, err := a.client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", a.token).
			SetResult(&authResponse{}).
			Post(a.url)

		if err != nil {
			return nil, fmt.Errorf("[auth-refresh] can't send request to pocketbase %w", err)
		}

		if resp.IsError() {
			return nil, fmt.Errorf("[auth-refresh] pocketbase returned status: %d, msg: %s, err %w",
				resp.StatusCode(),
				resp.String(),
				ErrInvalidResponse,
			)
		}

		auth := *resp.Result().(*authResponse)
		a.token = auth.Token
		a.client.SetHeader("Authorization", auth.Token)
		a.tokenValid = time.Now().Add(60 * time.Minute)
		return nil, nil
	})
	return err
}
