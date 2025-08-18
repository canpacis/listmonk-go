package listmonkgo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

type Client struct {
	config *ClientConfig
}

func (c *Client) auth() string {
	return fmt.Sprintf("token %s:%s", c.config.APIUser, c.config.Token)
}

func (c *Client) do(ctx context.Context, method, path string, data any) (*http.Response, error) {
	endpoint, err := url.JoinPath(c.config.BaseURL, path)
	if err != nil {
		return nil, err
	}

	if data != nil {
		q, err := query.Values(data)
		if err != nil {
			return nil, err
		}
		encoded := q.Encode()
		if len(encoded) > 0 {
			endpoint = fmt.Sprintf("%s?%s", endpoint, q.Encode())
		}
	}

	// This has to be io.Reader otherwise http.NewRequest call panics with nil values
	var body io.Reader
	if (method == "POST" || method == "PUT") && data != nil {
		body = new(bytes.Buffer)
		encoder := json.NewEncoder(body.(*bytes.Buffer))
		if err := encoder.Encode(data); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.auth())
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) multipart(ctx context.Context, path string, fields map[string]string, files map[string]io.Reader) (*http.Response, error) {
	endpoint, err := url.JoinPath(c.config.BaseURL, path)
	if err != nil {
		return nil, err
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	for key, r := range files {
		w, err := writer.CreateFormFile(key, key)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(w, r); err != nil {
			return nil, err
		}
	}

	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.auth())
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type Response[T any] struct {
	Data T `json:"data"`
}

func decode[T any](resp *http.Response) (*T, error) {
	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		data := new(ErrorResponse)
		if err := decoder.Decode(data); err != nil {
			return nil, err
		}
		return nil, errors.New(data.Message)
	}

	data := new(T)
	if err := decoder.Decode(data); err != nil {
		return nil, err
	}
	return data, nil
}

func request[T any](client *Client, ctx context.Context, method, path string, data any) (*T, error) {
	resp, err := client.do(ctx, method, path, data)
	if err != nil {
		return nil, err
	}
	return decode[T](resp)
}

type ClientConfig struct {
	BaseURL    string
	APIUser    string
	Token      string
	HTTPClient *http.Client
}

func WithBaseURL(baseUrl string) func(*ClientConfig) {
	return func(cc *ClientConfig) {
		cc.BaseURL = baseUrl
	}
}

func WithAPIUser(apiUser string) func(*ClientConfig) {
	return func(cc *ClientConfig) {
		cc.APIUser = apiUser
	}
}

func WithToken(token string) func(*ClientConfig) {
	return func(cc *ClientConfig) {
		cc.Token = token
	}
}

func WithHTTPClient(client *http.Client) func(*ClientConfig) {
	return func(cc *ClientConfig) {
		cc.HTTPClient = client
	}
}

type ConfigOption func(*ClientConfig)

func New(options ...ConfigOption) *Client {
	config := &ClientConfig{
		HTTPClient: http.DefaultClient,
	}

	for _, option := range options {
		option(config)
	}

	return NewWithConfig(config)
}

func NewWithConfig(config *ClientConfig) *Client {
	return &Client{
		config: config,
	}
}
