package github

import (
	"context"
	"fmt"
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
		return nil, fmt.Errorf("get GitHub installation: %w", err)
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
		return nil, fmt.Errorf("create GitHub installation token: %w", err)
	}

	return &result, nil
}

func (c *Client) ListRepositories(ctx context.Context, installationToken string) (*RepositoriesResponse, error) {
	var result RepositoriesResponse

	if err := c.doJSON(
		ctx,
		http.MethodGet,
		"/installation/repositories?per_page=100",
		installationToken,
		nil,
		http.StatusOK,
		&result,
	); err != nil {
		return nil, fmt.Errorf("list GitHub repositories: %w", err)
	}

	return &result, nil
}
