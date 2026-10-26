package github

import (
	"context"
	"net/http"
	"strconv"
)

func (c *Client) GetInstallation(ctx context.Context, installationID int64) (*Installation, error) {
	appJWT, err := c.generateJWT()
	if err != nil {
		return nil, err
	}

	var installation Installation

	path := "/app/installations/" + strconv.FormatInt(installationID, 10)

	if err := c.doJSON(
		ctx,
		http.MethodGet,
		path,
		appJWT,
		nil,
		http.StatusOK,
		&installation,
	); err != nil {
		return nil, err
	}

	return &installation, nil
}

func (c *Client) CreateInstallationToken(ctx context.Context, installationID int64, repositoryIDs []int64) (*InstallationToken, error) {
	appJWT, err := c.generateJWT()
	if err != nil {
		return nil, err
	}

	input := createInstallationTokenRequest{
		RepositoryIDs: repositoryIDs,
		Permissions: map[string]string{
			"contents": "read",
		},
	}

	var result InstallationToken
	path := "/app/installations/" + strconv.FormatInt(installationID, 10) + "/access_tokens"
	if err := c.doJSON(
		ctx,
		http.MethodPost,
		path,
		appJWT,
		input,
		http.StatusCreated,
		&result,
	); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) ListRepositories(ctx context.Context, installationToken string) (*RepositoriesResponse, error) {
	result := RepositoriesResponse{Repositories: make([]Repository, 0)}

	for page := 1; ; page++ {
		var current RepositoriesResponse
		path := "/installation/repositories?per_page=100&page=" + strconv.Itoa(page)

		if err := c.doJSON(
			ctx,
			http.MethodGet,
			path,
			installationToken,
			nil,
			http.StatusOK,
			&current,
		); err != nil {
			return nil, err
		}

		result.TotalCount = current.TotalCount
		result.Repositories = append(result.Repositories, current.Repositories...)
		if len(current.Repositories) < 100 {
			break
		}
	}

	return &result, nil
}
