package drive115

import (
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

type UploadInfo struct {
	UploadURL string `json:"upload_url"`
}

func (c *Client) GetUploadInfo() (*UploadInfo, error) {
	req, err := http.NewRequest("GET", apiGetUploadInfo, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", c.cookie)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get upload info failed with status: %s", resp.Status)
	}

	var info UploadInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	return &info, nil
}

type UploadResult struct {
	FileID string `json:"file_id"`
}

func (c *Client) Upload(reader io.Reader, uploadURL string) (*UploadResult, error) {
	req, err := http.NewRequest("POST", uploadURL, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/octet-stream")

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