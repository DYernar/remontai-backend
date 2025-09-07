package postgres

import (
	"context"

	"github.com/DYernar/remontai-backend/internal/domain"
)

// CreateImageGeneration inserts a new image generation
func (r *repository) CreateImageGeneration(ctx context.Context, imageGen domain.ImageGenerationModel) (domain.ImageGenerationModel, error) {
	query := `
		INSERT INTO image_generations (user_id, style_id, room_type, prompt, image_url, generated_image_url, status, error_message, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, user_id, style_id, room_type, prompt, image_url, generated_image_url, status, error_message, created_at, updated_at`

	var result domain.ImageGenerationModel
	err := r.conn.QueryRow(ctx, query,
		imageGen.UserID,
		imageGen.StyleID,
		imageGen.RoomType,
		imageGen.Prompt,
		imageGen.ImageURL,
		imageGen.GeneratedImageURL,
		imageGen.Status,
		imageGen.ErrorMessage,
	).Scan(
		&result.ID,
		&result.UserID,
		&result.StyleID,
		&result.RoomType,
		&result.Prompt,
		&result.ImageURL,
		&result.GeneratedImageURL,
		&result.Status,
		&result.ErrorMessage,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return domain.ImageGenerationModel{}, err
	}

	return result, nil
}

// GetImageGenerationsByUserId retrieves all image generations for a specific user ordered by created_at DESC
func (r *repository) GetImageGenerationsByUserId(ctx context.Context, userID string) ([]domain.ImageGenerationModel, error) {
	query := `
		SELECT id, user_id, style_id, room_type, prompt, image_url, generated_image_url, status, error_message, created_at, updated_at
		FROM image_generations 
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.conn.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	imageGenerations := []domain.ImageGenerationModel{}
	for rows.Next() {
		var imageGen domain.ImageGenerationModel
		err := rows.Scan(
			&imageGen.ID,
			&imageGen.UserID,
			&imageGen.StyleID,
			&imageGen.RoomType,
			&imageGen.Prompt,
			&imageGen.ImageURL,
			&imageGen.GeneratedImageURL,
			&imageGen.Status,
			&imageGen.ErrorMessage,
			&imageGen.CreatedAt,
			&imageGen.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		imageGenerations = append(imageGenerations, imageGen)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return imageGenerations, nil
}
