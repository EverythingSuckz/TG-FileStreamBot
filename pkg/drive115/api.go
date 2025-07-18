package drive115

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	apiGetUploadInfo = "https://uplb.115.com/3.0/getuploadinfo.php"
)


type Client struct {
	cookie     string
	httpClient *http.Client
}

func NewClient(cookie string) *Client {
	return &Client{
		cookie:     cookie,
		httpClient: &http.Client{},
	}
}

func (c *Client) GetUploadInfo() (string, error) {
	// This is a placeholder. The actual implementation will be more complex.
	// It needs to make a request to apiGetUploadInfo and parse the response.
	return "upload_url", nil
}

func (c *Client) Upload(reader io.Reader, uploadURL string) (*UploadResult, error) {
	// This is a placeholder. The actual implementation will be more complex.
	// It needs to stream the data from the reader to the uploadURL.
	req, err := http.NewRequest("POST", uploadURL, reader)
	if err != nil {
		return nil, err
	}
	// Add necessary headers, such as Content-Type.

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upload failed with status: %s", resp.Status)
	}

	var result UploadResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}