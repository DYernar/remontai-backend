package openai

import (
	"context"
	"fmt"
	"io"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	client := openai.NewClient(apiKey)
	return &OpenAIClient{client: client}
}

func (c *OpenAIClient) GenerateRoomDesign(ctx context.Context, style, roomType string, image io.Reader) (string, error) {
	prompt := c.getPrompt(style, roomType)

	imgReq := openai.ImageEditRequest{
		Image:          openai.WrapReader(image, "init_image.png", "image/png"),
		Prompt:         prompt,
		Model:          openai.CreateImageModelGptImage1,
		N:              1,
		Size:           openai.CreateImageSize1024x1024,
		ResponseFormat: openai.CreateImageResponseFormatURL,
	}

	imgResp, err := c.client.CreateEditImage(ctx, imgReq)
	if err != nil {
		return "", fmt.Errorf("CreateEditImage error: %w", err)
	}
	if len(imgResp.Data) == 0 {
		return "", fmt.Errorf("CreateEditImage returned empty data")
	}

	// 5) Return the resulting image URL
	return imgResp.Data[0].URL, nil
}
