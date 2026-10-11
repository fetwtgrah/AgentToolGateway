package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Client struct {
	BaseURL string
}

func NewClient(url string) *Client {

	return &Client{
		BaseURL: url,
	}

}
func (c *Client) CallTool(call any) error {
	data, _ := json.Marshal(call)
	resp, err := http.Post(
		c.BaseURL+"/toolCalls",
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
