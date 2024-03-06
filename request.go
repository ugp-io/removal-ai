package removalai

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	Crop      string
	Ecom      string
	GetBase64 string
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

	fmt.Println("HERE internal")
	// Prepare the body of the POST request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("image_url", request.ImageURL)
	_ = writer.WriteField("crop", request.Crop)
	_ = writer.WriteField("ecom", request.Ecom)
	_ = writer.WriteField("get_base64", request.GetBase64)
	err := writer.Close()
	if err != nil {
		return nil, err
	}

	// Create a new request using http
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	// Add the required headers to the request
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("accept", "application/json")
	req.Header.Set("Rm-Token", c.APIKey)

	// Send the request using a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Unmarshal the JSON response into the ResponseData struct
	var responseData ImageRemovalResponse
	var response interface{}
	if err := json.Unmarshal(responseBody, &responseData); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		panic(err)
	}
	fmt.Println(response)

	return &responseData, nil
}
