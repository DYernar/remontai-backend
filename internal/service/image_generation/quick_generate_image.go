package imagegeneration

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/DYernar/remontai-backend/internal/domain"
)

func (s *service) QuickGenerateImage(
	ctx context.Context,
	userID string,
	imageFile multipart.File,
	imageHeader *multipart.FileHeader,
	roomType,
	styleID string,
) (domain.ImageGenerationModel, error) {
	imageURL, err := s.uploadToS3(ctx, imageFile, imageHeader)
	if err != nil {
		return domain.ImageGenerationModel{}, err
	}

	imageGen := domain.ImageGenerationModel{
		UserID:            userID,
		StyleID:           styleID,
		RoomType:          roomType,
		Prompt:            fmt.Sprintf("Generate %s design", roomType),
		ImageURL:          imageURL,
		GeneratedImageURL: "",
		Status:            domain.ImageGenerateStatusPending,
		ErrorMessage:      "",
	}

	result, err := s.repo.CreateImageGeneration(ctx, imageGen)
	if err != nil {
		return domain.ImageGenerationModel{}, err
	}

	return result, nil
}

func (s *service) uploadToS3(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	rand.Seed(time.Now().UnixNano())
	randomID := rand.Intn(1000000)
	extension := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("uploads/%d_%d_%s%s", timestamp, randomID, "image", extension)

	// Upload to Google Cloud Storage
	err = s.s3.UploadFile(ctx, fileBytes, filename)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to storage: %w", err)
	}

	// Return the public URL for Google Cloud Storage
	// Format: https://storage.googleapis.com/bucket-name/object-name
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", "remont_ai_media_storage", filename)

	return publicURL, nil
}
