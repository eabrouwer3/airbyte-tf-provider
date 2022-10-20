package apiclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const BASE_URL = "api/v1"

type ApiClient struct {
	HostURL    string `http://localhost:8000`
	HTTPClient *http.Client
}

type HealthCheckResponse struct {
	available bool
}

func (c *ApiClient) check() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/health", c.HostURL, BASE_URL), nil)
	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	hcr := HealthCheckResponse{}
	err = json.Unmarshal(body, &hcr)
	if err != nil {
		return err
	}

	return nil
}

func (c *ApiClient) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
