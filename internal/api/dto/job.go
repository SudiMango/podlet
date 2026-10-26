package dto

import (
	"time"

	"github.com/SudiMango/podlet/internal/contract"
)

// Request

type CreateJobRequest struct {
	JobId          string
	ServerId       string
	JobType        contract.JobType
	RepoName       string
	RepositoryId   int64
	InstallationId int64
}

// Response

type JobResponse struct {
	Id             string
	UserId         string
	ServerId       string
	JobStatus      contract.JobStatus
	JobType        contract.JobType
	ErrorMessage   *string
	RepoName       string
	RepositoryId   int64
	InstallationId int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
