package gemini

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/genai"
)

type Client struct {
	gemini *genai.Client
}

func NewClient(ctx context.Context, apiKey string) (*Client, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, err
	}
	return &Client{gemini: client}, nil
}

func (c *Client) GenerateImage(imgData []byte, styleName, roomType string) ([]byte, error) {
	ctx := context.Background()

	parts := []*genai.Part{
		genai.NewPartFromText(c.getPrompt(styleName, roomType)),
		{
			InlineData: &genai.Blob{
				MIMEType: "image/png",
				Data:     imgData,
			},
		},
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := c.gemini.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-image-preview",
		contents,
		nil,
	)

	if err != nil {
		return nil, err
	}

	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			fmt.Println(part.Text)
		} else if part.InlineData != nil {
			return part.InlineData.Data, nil
		}
	}

	return nil, errors.New("no image generated")
}
