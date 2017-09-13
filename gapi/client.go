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

type ApiError struct {
	ResponseCode   int
	ResponseStatus string
	Message        string
}

func (e ApiError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.ResponseStatus
}

type Client struct {
	key     string
	headers map[string]string
	baseURL url.URL
	*http.Client
}

//New creates a new grafana client
//auth can be in user:pass format, or it can be an api key
//headers should be a comma-separated list of extra headers to send with each
//request, e.g. "X-User:foo,X-Grafana-Org-Id:1"
func New(auth, headers, baseURL string) (*Client, error) {
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
	hdr := make(map[string]string)
	if headers != "" {
		split := strings.Split(headers, ",")
		for _, header := range split {
			kv := strings.Split(header, ":")
			hdr[kv[0]] = kv[1]
		}
	}
	return &Client{
		key:     key,
		headers: hdr,
		baseURL: *u,
		Client:  &http.Client{},
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
	for k, v := range c.headers {
		req.Header.Add(k, v)
	}
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
