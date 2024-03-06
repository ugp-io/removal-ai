package removalai

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
)

type Client struct {
	APIKey string
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
	}
}

type ImageRemovalRequest struct {
	ImageURL  string
	Crop      *string
	Ecom      *string
	GetBase64 *string
}

type ImageRemovalResponse struct {
	Status         *int    `json:"status,omitempty"`
	Demo           *string `json:"demo,omitempty"`
	PreviewDemo    *string `json:"preview_demo,omitempty"`
	URL            *string `json:"url,omitempty"`
	HighRes        *string `json:"high_resolution,omitempty"`
	LowRes         *string `json:"low_resolution,omitempty"`
	Base64         *string `json:"base64,omitempty"`
	OriginalWidth  *int    `json:"original_width,omitempty"`
	OriginalHeight *int    `json:"original_height,omitempty"`
	PreviewWidth   *int    `json:"preview_width,omitempty"`
	PreviewHeight  *int    `json:"preview_height,omitempty"`
	Extra          *string `json:"extra,omitempty"`
}

func (c *Client) BackgroundRemoval(request ImageRemovalRequest) (*ImageRemovalResponse, error) {

	url := "https://api.removal.ai/3.0/remove"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("image_url", request.ImageURL)

	if request.Crop != nil {
		writer.WriteField("crop", *request.Crop)
	}

	if request.Ecom != nil {
		writer.WriteField("ecom", *request.Ecom)
	}

	if request.GetBase64 != nil {
		writer.WriteField("get_base64", *request.GetBase64)
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("accept", "application/json")
	req.Header.Set("Rm-Token", c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var responseData ImageRemovalResponse
	if err := json.Unmarshal(responseBody, &responseData); err != nil {
		return nil, err
	}

	return &responseData, nil
}
