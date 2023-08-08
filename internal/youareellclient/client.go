package youareellclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const URLBase = "http://zipcode.rocks:8085"

type Client struct {
	Debug   bool
	BaseURL string
}

func (c *Client) Raw(method string, location string, input interface{}, output interface{}) error {
	baseURL := c.BaseURL
	if baseURL == "" {
		baseURL = URLBase
	}

	targetURL := baseURL + location

	if c.Debug {
		fmt.Printf("[%s %s]\n", method, targetURL)
	}

	var body io.Reader
	{
		if input != nil {
			contents, err := json.Marshal(input)
			if err != nil {
				return err
			}
			body = bytes.NewReader(contents)
			if c.Debug {
				fmt.Printf("[Body: %s]\n", contents)
			}
		}
	}
	request, err := http.NewRequest(method, targetURL, body)
	if err != nil {
		return err
	}

	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if output != nil {
		request.Header.Set("Accept", "application/json")
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	if c.Debug {
		fmt.Printf("[Status: %d]\n", response.StatusCode)
	}

	contents, _ := io.ReadAll(response.Body)
	if c.Debug {
		fmt.Printf("[Output: %s]\n", contents)
	}

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		if output != nil {
			err = json.Unmarshal(contents, output)
			if err != nil {
				return err
			}
		}
		return nil
	}

	return fmt.Errorf("http status: %d", response.StatusCode)
}
