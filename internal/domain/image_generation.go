package domain

import "time"

// Domain Model
type ImageGenerationModel struct {
	ID                string              `json:"id"`
	UserID            string              `json:"user_id"`
	StyleID           string              `json:"style_id"`
	RoomType          string              `json:"room_type"`
	Prompt            string              `json:"prompt"`
	ImageURL          string              `json:"image_url"`
	GeneratedImageURL string              `json:"generated_image_url"`
	Status            ImageGenerateStatus `json:"status"`
	ErrorMessage      string              `json:"error_message"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
}

type ImageGenerateStatus string

const (
	ImageGenerateStatusPending   ImageGenerateStatus = "pending"
	ImageGenerateStatusCompleted ImageGenerateStatus = "completed"
	ImageGenerateStatusFailed    ImageGenerateStatus = "failed"
)
