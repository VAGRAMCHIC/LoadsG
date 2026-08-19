package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"loadsg/lib/dto"
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

func New(address string) *Client {
	return &Client{
		BaseURL: address,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) SetToken(t string) {
	c.Token = t
}

func (c *Client) do(method, path string, body any, out any) error {
	var req *http.Request
	var err error

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		req, err = http.NewRequest(method, c.BaseURL+path, bytes.NewReader(b))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, c.BaseURL+path, nil)
		if err != nil {
			return err
		}
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: %s", resp.Status)
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

// Ответ /v1/auth/login — handler отдаёт gin.H{"access_token": string}
type LoginResult struct {
	AccessToken string `json:"access_token"`
}

func (c *Client) Login(uid, token string) (*LoginResult, error) {
	var res LoginResult
	err := c.do("POST", "/v1/auth/login", dto.LoginRequest{
		UID:       uid,
		TokenHash: token,
	}, &res)
	return &res, err
}

// Ответ create-* — handler отдаёт gin.H{"job_id": loadJob.Id}
type CreateJobResult struct {
	JobID string `json:"job_id"`
}

func (c *Client) CreateConstantLoad(req dto.CreateConstantHTTPLoadRequest) (*CreateJobResult, error) {
	var res CreateJobResult
	err := c.do("POST", "/v1/manager/http/constant", req, &res)
	return &res, err
}

func (c *Client) CreateRampUpLoad(req dto.CreateRampUpHTTPLoadRequest) (*CreateJobResult, error) {
	var res CreateJobResult
	err := c.do("POST", "/v1/manager/http/ramp-up", req, &res)
	return &res, err
}

func (c *Client) CreateFixedLoad(req dto.CreateFixedHTTPLoadRequest) (*CreateJobResult, error) {
	var res CreateJobResult
	err := c.do("POST", "/v1/manager/http/fixed", req, &res)
	return &res, err
}

func (c *Client) CreateFakeLoad(req dto.CreateFakeHTTPLoadRequest) (*CreateJobResult, error) {
	var res CreateJobResult
	err := c.do("POST", "/v1/manager/http/fake", req, &res)
	return &res, err
}

func (c *Client) DeleteLoadJob(id string) error {
	return c.do("DELETE", "/v1/manager/http/constant/"+id, nil, nil)
}