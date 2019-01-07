package brahma

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	log "gopkg.in/src-d/go-log.v1"
)

type Client struct {
	url     *url.URL
	client  *http.Client
	storage string
}

func NewClient(server string) (*Client, error) {
	url, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	return &Client{
		url:     url,
		client:  new(http.Client),
		storage: "client",
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

func (c *Client) Download() error {
	for {
		repo, err := c.Repository()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		name := fmt.Sprintf("%s.siva", repo.ID)
		path := filepath.Join(c.storage, name)

		log.With(log.Fields{
			"id":   repo.ID,
			"url":  repo.URL,
			"siva": path},
		).Infof("downloading repository")

		err = Download(repo.URL, path)
		if err != nil {
			return err
		}
	}
}
