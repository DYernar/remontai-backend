package postgres

import (
	"context"

	"github.com/DYernar/remontai-backend/internal/domain"
)

// ListStyles retrieves all styles ordered by created_at DESC
func (r *repository) ListStyles(ctx context.Context) ([]domain.StyleModel, error) {
	query := `
		SELECT id, name, description, image, created_at, updated_at 
		FROM styles 
		ORDER BY created_at DESC`

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	styles := []domain.StyleModel{}
	for rows.Next() {
		var style domain.StyleModel
		err := rows.Scan(&style.ID, &style.Name, &style.Description,
			&style.Image, &style.CreatedAt, &style.UpdatedAt)
		if err != nil {
			return nil, err
		}
		styles = append(styles, style)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return styles, nil
}
