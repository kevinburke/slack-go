package slack

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/kevinburke/rest"
)

const Version = "0.1"

var userAgent string

func init() {
	userAgent = fmt.Sprintf("slack-go/%s", Version)
}

func New(token string) *Client {
	restClient := rest.NewBearerClient(token, baseURL)
	restClient.UploadType = rest.FormURLEncoded
	c := &Client{
		Client:  restClient,
		baseURL: baseURL,
	}
	c.Chat = &ChatService{client: c}
	c.Users = &UserService{client: c}
	return c
}

const baseURL = "https://api.slack.com"

type Client struct {
	*rest.Client
	baseURL string

	Chat  *ChatService
	Users *UserService
}

// CreateResource makes a POST request to the given resource.
func (c *Client) CreateResource(ctx context.Context, pathPart string, data url.Values, v interface{}) error {
	return c.MakeRequest(ctx, "POST", pathPart, data, v)
}

func (c *Client) ListResource(ctx context.Context, pathPart string, data url.Values, v interface{}) error {
	return c.MakeRequest(ctx, "GET", pathPart, data, v)
}

// Make a request to the Slack API.
func (c *Client) MakeRequest(ctx context.Context, method string, pathPart string, data url.Values, v interface{}) error {
	rb := new(strings.Reader)
	if data != nil && (method == "POST" || method == "PUT") {
		rb = strings.NewReader(data.Encode())
	}
	if method == "GET" && data != nil {
		pathPart = pathPart + "?" + data.Encode()
	}
	req, err := c.NewRequest(method, pathPart, rb)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	if ua := req.Header.Get("User-Agent"); ua == "" {
		req.Header.Set("User-Agent", userAgent)
	} else {
		req.Header.Set("User-Agent", userAgent+" "+ua)
	}
	return c.Do(req, &v)
}