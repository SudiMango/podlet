package service

import (
	"context"

	"github.com/SudiMango/podlet/internal/api/dto"
	"github.com/SudiMango/podlet/internal/api/repository"
)

type ProfileService struct {
	profileRepository *repository.ProfileRepository
}

func NewProfileService(profileRepository *repository.ProfileRepository) *ProfileService {
	return &ProfileService{
		profileRepository: profileRepository,
	}
}

func (s *ProfileService) GetPodletProfileByUserId(ctx context.Context, userId string) (*dto.PodletProfileResponse, error) {
	user, err := s.profileRepository.GetPodletProfileByUserId(ctx, userId)
	if err != nil || user == nil {
		return nil, err
	}

	return user, nil
}
