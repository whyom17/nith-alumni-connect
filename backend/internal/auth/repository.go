package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/nith-alumni-connect/backend/internal/platform/db"
)

// Repository defines data access operations required by the auth service.
type Repository interface {
	GetUserByLinkedInSub(ctx context.Context, sub string) (*db.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*db.User, error)
	CreateUser(ctx context.Context, linkedinSub, email, role, status string) (*db.User, error)
	CreateRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) (*db.RefreshToken, error)
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*db.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error
}

type sqlRepository struct {
	queries *db.Queries
}

// NewRepository creates a new database-backed auth repository.
func NewRepository(queries *db.Queries) Repository {
	return &sqlRepository{queries: queries}
}

func (r *sqlRepository) GetUserByLinkedInSub(ctx context.Context, sub string) (*db.User, error) {
	u, err := r.queries.GetUserByLinkedInSub(ctx, sub)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *sqlRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*db.User, error) {
	var pgUUID pgtype.UUID
	if err := pgUUID.Scan(id.String()); err != nil {
		return nil, err
	}
	u, err := r.queries.GetUserByID(ctx, pgUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *sqlRepository) CreateUser(ctx context.Context, linkedinSub, email, role, status string) (*db.User, error) {
	u, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		LinkedinSub: linkedinSub,
		Email:       email,
		Role:        role,
		Status:      status,
	})
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *sqlRepository) CreateRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) (*db.RefreshToken, error) {
	var pgUUID pgtype.UUID
	if err := pgUUID.Scan(userID.String()); err != nil {
		return nil, err
	}
	t, err := r.queries.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:    pgUUID,
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *sqlRepository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*db.RefreshToken, error) {
	t, err := r.queries.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *sqlRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	return r.queries.RevokeRefreshToken(ctx, tokenHash)
}

func (r *sqlRepository) RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	var pgUUID pgtype.UUID
	if err := pgUUID.Scan(userID.String()); err != nil {
		return err
	}
	return r.queries.RevokeAllUserRefreshTokens(ctx, pgUUID)
}
