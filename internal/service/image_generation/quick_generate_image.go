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
	fileBytes, err := s.getImageBytes(imageFile)
	if err != nil {
		return domain.ImageGenerationModel{}, err
	}

	imageURL, err := s.uploadToS3(ctx, fileBytes, imageHeader)
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

	// imageResp, err := s.fluxClient.GenerateImage(imageURL, style.Name, roomType)
	// if err != nil {
	// 	s.logger.Errorf("Error generating image", "error", err, "userid", userID)
	// 	return domain.ImageGenerationModel{}, err
	// }

	// if imageResp.Processing == nil || len(imageResp.Processing.FutureLinks) == 0 {
	// 	s.logger.Errorf("Error generating image: no future links", "userid", userID)
	// 	return domain.ImageGenerationModel{}, fmt.Errorf("no future links in flux response")
	// }

	// imageGen.GeneratedImageURL = imageResp.Processing.FutureLinks[0]

	genImageResp, err := s.geminiClient.GenerateImage(fileBytes, style.Name, roomType)
	if err != nil {
		s.logger.Errorf("Error generating image", "error", err, "userid", userID)
		return domain.ImageGenerationModel{}, err
	}

	// upload generated image to s3
	uploadedImageURL, err := s.uploadByteImageToS3(ctx, genImageResp)
	if err != nil {
		s.logger.Errorf("Error uploading generated image to S3", "error", err, "userid", userID)
		return domain.ImageGenerationModel{}, err
	}

	imageGen.GeneratedImageURL = uploadedImageURL
	imageGen.Status = domain.ImageGenerateStatusCompleted

	result, err := s.repo.CreateImageGeneration(ctx, imageGen)
	if err != nil {
		return domain.ImageGenerationModel{}, err
	}

	return result, nil
}

func (s *service) uploadByteImageToS3(ctx context.Context, imageBytes []byte) (string, error) {
	// extension as png
	filename := fmt.Sprintf("%s.png", uuid.New().String())
	err := s.s3.UploadFile(ctx, imageBytes, filename)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to storage: %w", err)
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", "remont_ai_media_storage", filename)

	return publicURL, nil
}

func (s *service) getImageBytes(image multipart.File) ([]byte, error) {
	fileBytes, err := io.ReadAll(image)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return fileBytes, nil
}

func (s *service) uploadToS3(ctx context.Context, fileBytes []byte, header *multipart.FileHeader) (string, error) {
	extension := filepath.Ext(header.Filename)
	imageName := uuid.New().String()
	filename := fmt.Sprintf("uploads/%s_%s%s", imageName, "image", extension)

	// Upload to Google Cloud Storage
	err := s.s3.UploadFile(ctx, fileBytes, filename)
	if err != nil {
		fmt.Printf("Failed to upload file to storage: %v\n", err)
		return "", fmt.Errorf("failed to upload file to storage: %w", err)
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", "remont_ai_media_storage", filename)

	return publicURL, nil
}
