package imagegeneration

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/DYernar/remontai-backend/internal/domain"
	"github.com/google/uuid"
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

	style, err := s.repo.GetStyleByID(ctx, styleID)
	if err != nil {
		return domain.ImageGenerationModel{}, err
	}

	// imageLink, err := s.openAIClient.GenerateRoomDesign(
	// 	ctx,
	// 	style.Name,
	// 	roomType,
	// 	imageFile,
	// )
	// if err != nil {
	// 	s.logger.Errorf("Error generating image", "error", err, "userid", userID)
	// 	return domain.ImageGenerationModel{}, err
	// }

	// imageGen.GeneratedImageURL = imageLink

	imageResp, err := s.fluxClient.GenerateImage(imageURL, style.Name, roomType)
	if err != nil {
		s.logger.Errorf("Error generating image", "error", err, "userid", userID)
		return domain.ImageGenerationModel{}, err
	}

	if imageResp.Processing == nil || len(imageResp.Processing.FutureLinks) == 0 {
		s.logger.Errorf("Error generating image: no future links", "userid", userID)
		return domain.ImageGenerationModel{}, fmt.Errorf("no future links in flux response")
	}

	imageGen.GeneratedImageURL = imageResp.Processing.FutureLinks[0]
	imageGen.Status = domain.ImageGenerateStatusCompleted

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
	// store to local folder

	extension := filepath.Ext(header.Filename)
	imageName := uuid.New().String()
	filename := fmt.Sprintf("uploads/%s_%s%s", imageName, "image", extension)

	// Upload to Google Cloud Storage
	err = s.s3.UploadFile(ctx, fileBytes, filename)
	if err != nil {
		fmt.Printf("Failed to upload file to storage: %v\n", err)
		return "", fmt.Errorf("failed to upload file to storage: %w", err)
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", "remont_ai_media_storage", filename)

	return publicURL, nil
}
