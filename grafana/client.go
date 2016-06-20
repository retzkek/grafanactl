package gapi

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Client struct {
	key     string
	baseURL url.URL
	*http.Client
}

//New creates a new grafana client
//auth can be in user:pass format, or it can be an api key
func New(auth, baseURL string) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	key := ""
	if strings.Contains(auth, ":") {
		split := strings.Split(auth, ":")
		u.User = url.UserPassword(split[0], split[1])
	} else if auth != "" {
		key = fmt.Sprintf("Bearer %s", auth)
	}
	return &Client{
		key,
		*u,
		&http.Client{},
	}, nil
}

func (c *Client) newRequest(method, uri string, body io.Reader) (*http.Request, error) {
	url := c.baseURL
	url.Path = path.Join(url.Path, uri)
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return req, err
	}
	if c.key != "" {
		req.Header.Add("Authorization", c.key)
	}
	if body == nil {
		log.WithFields(log.Fields{
			"url": url.String(),
		}).Debug("request")
	} else {
		log.WithFields(log.Fields{
			"url":  url.String(),
			"body": body.(*bytes.Buffer).String(),
		}).Debug("request")
	}
	req.Header.Add("Content-Type", "application/json")
	return req, err
}

func (c *Client) DoRead(req *http.Request) ([]byte, error) {
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
