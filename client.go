package zenfo

import (
	"net/http"
	"time"
)

// Client wraps http.Client for custom UA string, Getting, etc.
// TODO Add rate limit per worker
type Client struct {
	client *http.Client
}

// NewClient returns new Client object
func NewClient() *Client {
	c := new(Client)
	c.client = &http.Client{
		Timeout: time.Second * 10,
	}
	return c
}

// Get wraps core http GET functionality with custom UA string
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "zenfo.info/1.0")
	return c.client.Do(req)
}
