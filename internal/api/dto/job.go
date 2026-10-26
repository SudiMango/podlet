package dto

import "github.com/SudiMango/podlet/internal/contract"

type CreateJobRequest struct {
	JobId          string
	ServerId       string
	JobType        contract.JobType
	RepoName       string
	RepositoryId   int64
	InstallationId int64
}
