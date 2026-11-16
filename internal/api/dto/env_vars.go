package dto

// Request

type CreateEnvVarRequest struct {
	Id             string
	Key            string
	ValueEncrypted string
	RepositoryId   int64
	RepoName       string
}

type GetEnvVarsRequest struct {
	RepoName     string
	RepositoryId int64
}
