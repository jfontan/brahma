package brahma

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Client struct {
	url    *url.URL
	client *http.Client
}

func NewClient(server string) (*Client, error) {
	url, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	return &Client{
		url:    url,
		client: new(http.Client),
	}, nil
}

func (c *Client) path(path string) string {
	u := &url.URL{Path: path}
	return c.url.ResolveReference(u).String()
}

func (c *Client) Repository() (Repository, error) {
	res, err := c.client.Get(c.path("/repository"))
	if err != nil {
		return Repository{}, err
	}

	var repo Repository
	err = json.NewDecoder(res.Body).Decode(&repo)
	if err != nil {
		return Repository{}, err
	}

	return repo, nil
}
