package contract

import "time"

// Request

type CreateEnvVarRequest struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	RepoName string `json:"repo_name"`
}

// Response

type EnvVarResponse struct {
	Id        string    `json:"id"`
	Key       string    `json:"key"`
	RepoName  string    `json:"repo_name"`
	CreatedAt time.Time `json:"created_at"`
}

type EnvVarsResponse []EnvVarResponse

type EnvVarForAgentResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EnvVarsForAgentResponse []EnvVarForAgentResponse
