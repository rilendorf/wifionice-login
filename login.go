package main

import (
	"gopkg.in/resty.v1"

	"errors"
	"net/url"

	"regexp"
)

var (
	ErrInvalidResponse = errors.New("Invalid response :(")
)

var CookieExp = regexp.MustCompile("csrf=([0-9a-f]*);")

type IndexRequest struct{}

func (ir *IndexRequest) Send() (res *IndexResponse, err error) {
	resp, err := resty.R().
		Get("https://login.wifionice.de/en/")
	if err != nil {
		return nil, err
	}

	c := resp.Header().Get("Set-Cookie")

	r := CookieExp.FindStringSubmatch(c)
	if len(r) != 2 {
		return nil, ErrInvalidResponse
	}

	return &IndexResponse{
		CSRFToken: r[1],
	}, nil
}

type IndexResponse struct {
	CSRFToken string
}

type LoginRequest struct {
	Login     bool
	Logout    bool
	CSRFToken string
}

type LoginResponse struct {
	StatusCode int
}

func (lr *LoginRequest) Send() (res *LoginResponse, err error) {
	v := url.Values(map[string][]string{
		"CSRFToken": []string{lr.CSRFToken},
	})

	if lr.Login {
		v["login"] = []string{"true"}
	}

	if lr.Logout {
		v["logout"] = []string{"true"}
	}

	resp, err := resty.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Cookie", "csrf="+lr.CSRFToken).
		SetBody(v.Encode()).
		Post("https://login.wifionice.de/en/")
	if err != nil {
		if UnwrapErr(err).Error() == "auto redirect is disabled" {
			err = nil
		} else {
			return
		}
	}

	return &LoginResponse{resp.StatusCode()}, nil
}

func UnwrapErr(e error) error {
	err, ok := e.(*url.Error)
	if !ok {
		return e
	}

	return err.Err
}
