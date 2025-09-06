package postgres

import (
	"context"
	"errors"

	"github.com/DYernar/remontai-backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (r *repository) CreateUser(ctx context.Context, user domain.UserModel) (domain.UserModel, error) {
	query := `
		INSERT INTO users (email, name, image, token, push_token) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, email, name, image, token, push_token, created_at, updated_at`

	var createdUser domain.UserModel
	err := r.conn.QueryRow(ctx, query, user.Email, user.Name, user.Image, user.Token, user.PushToken).
		Scan(&createdUser.ID, &createdUser.Email, &createdUser.Name, &createdUser.Image,
			&createdUser.Token, &createdUser.PushToken, &createdUser.CreatedAt, &createdUser.UpdatedAt)

	if err != nil {
		return domain.UserModel{}, err
	}

	return createdUser, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (domain.UserModel, error) {
	query := `
		SELECT id, email, name, image, token, push_token, created_at, updated_at 
		FROM users 
		WHERE email = $1`

	var user domain.UserModel
	err := r.conn.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.Name, &user.Image,
			&user.Token, &user.PushToken, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserModel{}, domain.ErrUserNotFound
		}
		return domain.UserModel{}, err
	}

	return user, nil
}

func (r *repository) GetUserByToken(ctx context.Context, token string) (domain.UserModel, error) {
	query := `
		SELECT id, email, name, image, token, push_token, created_at, updated_at 
		FROM users 
		WHERE token = $1`

	var user domain.UserModel
	err := r.conn.QueryRow(ctx, query, token).
		Scan(&user.ID, &user.Email, &user.Name, &user.Image,
			&user.Token, &user.PushToken, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserModel{}, domain.ErrUserNotFound
		}
		return domain.UserModel{}, err
	}

	return user, nil
}

func (r *repository) UpsertUser(ctx context.Context, user domain.UserModel) (domain.UserModel, error) {
	query := `
		INSERT INTO users (email, name, image, token, push_token) 
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (email) 
		DO UPDATE SET 
			name = EXCLUDED.name,
			image = EXCLUDED.image,
			token = EXCLUDED.token,
			push_token = EXCLUDED.push_token,
			updated_at = NOW()
		RETURNING id, email, name, image, token, push_token, created_at, updated_at`

	var upsertedUser domain.UserModel
	err := r.conn.QueryRow(ctx, query, user.Email, user.Name, user.Image, user.Token, user.PushToken).
		Scan(&upsertedUser.ID, &upsertedUser.Email, &upsertedUser.Name, &upsertedUser.Image,
			&upsertedUser.Token, &upsertedUser.PushToken, &upsertedUser.CreatedAt, &upsertedUser.UpdatedAt)

	if err != nil {
		return domain.UserModel{}, err
	}

	return upsertedUser, nil
}

func (r *repository) DeleteUser(ctx context.Context, userID string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.conn.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *repository) GetUserByID(ctx context.Context, userID string) (domain.UserModel, error) {
	query := `
		SELECT id, email, name, image, token, push_token, created_at, updated_at 
		FROM users 
		WHERE id = $1`

	var user domain.UserModel
	err := r.conn.QueryRow(ctx, query, userID).
		Scan(&user.ID, &user.Email, &user.Name, &user.Image,
			&user.Token, &user.PushToken, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserModel{}, domain.ErrUserNotFound
		}
		return domain.UserModel{}, err
	}

	return user, nil
}
