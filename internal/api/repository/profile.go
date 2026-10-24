package repository

import (
	"context"

	"github.com/SudiMango/podlet/internal/contract"
	db "github.com/SudiMango/podlet/internal/sql/database"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProfileRepository struct {
	store *db.Queries
}

func NewProfileRepository(store *db.Queries) *ProfileRepository {
	return &ProfileRepository{
		store: store,
	}
}

func (r *ProfileRepository) GetPodletProfileByUserId(ctx context.Context, userId string) (*contract.PodletProfileResponse, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	user, err := r.store.GetPodletProfileByUserId(ctx, pgUserId)
	if err != nil {
		return nil, err
	}

	return &contract.PodletProfileResponse{
		UserId:    user.UserID.String(),
		FirstName: user.FirstName,
		LastName:  &user.LastName.String,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}
