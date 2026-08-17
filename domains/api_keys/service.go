package apikeys

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"

	"idiom-api-services/packages/crypto"
)

var (
	ErrNameRequired = errors.New("api key name is required")
	ErrNotFound     = errors.New("api key not found")
	ErrInvalidKey   = errors.New("invalid api key")
)

type CreateInput struct {
	Name   string   `json:"name"`
	Scopes []string `json:"scopes"`
}

type CreateResult struct {
	APIKey *APIKey `json:"api_key"`
	Key    string  `json:"key"`
}

func Verify(
	ctx context.Context,
	repo *Repository,
	projectID string,
	key string,
) (*APIKey, error) {
	if projectID == "" || key == "" {
		return nil, ErrInvalidKey
	}

	prefix := keyPrefix(key)
	if prefix == "" {
		return nil, ErrInvalidKey
	}

	apiKey, err := repo.GetActiveByPrefix(ctx, projectID, prefix)
	if err != nil {
		return nil, err
	}

	if apiKey == nil {
		return nil, ErrInvalidKey
	}

	if !crypto.CheckPasswordHash(key, apiKey.KeyHash) {
		return nil, ErrInvalidKey
	}

	return apiKey, nil
}

func keyPrefix(key string) string {
	if len(key) < 11 {
		return ""
	}

	return key[:11]
}

func ListByProjectID(
	ctx context.Context,
	repo *Repository,
	projectID string,
	identityID string,
) ([]APIKey, error) {
	return repo.ListByProjectID(ctx, projectID, identityID)
}

func Create(
	ctx context.Context,
	repo *Repository,
	projectID string,
	identityID string,
	input CreateInput,
) (*CreateResult, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrNameRequired
	}

	key, prefix, hash, err := generateKey()
	if err != nil {
		return nil, err
	}

	apiKey := &APIKey{
		ID:        crypto.GenerateID("ak_"),
		Name:      name,
		ProjectID: projectID,
		KeyPrefix: prefix,
		KeyHash:   hash,
		Scopes:    input.Scopes,
		IsActive:  true,
		CreatedBy: identityID,
	}

	if err := repo.Create(ctx, apiKey, identityID); err != nil {
		return nil, err
	}

	return &CreateResult{
		APIKey: apiKey,
		Key:    key,
	}, nil
}

func Rotate(
	ctx context.Context,
	repo *Repository,
	projectID string,
	identityID string,
	apiKeyID string,
) (*CreateResult, error) {
	key, prefix, hash, err := generateKey()
	if err != nil {
		return nil, err
	}

	apiKey, err := repo.Rotate(
		ctx,
		apiKeyID,
		projectID,
		identityID,
		prefix,
		hash,
	)
	if err != nil {
		return nil, err
	}

	return &CreateResult{
		APIKey: apiKey,
		Key:    key,
	}, nil
}

func Revoke(
	ctx context.Context,
	repo *Repository,
	projectID string,
	identityID string,
	apiKeyID string,
) error {
	return repo.Revoke(
		ctx,
		apiKeyID,
		projectID,
		identityID,
	)
}

func generateKey() (string, string, string, error) {
	raw := make([]byte, 32)

	if _, err := rand.Read(raw); err != nil {
		return "", "", "", err
	}

	key := "ik_" + base64.RawURLEncoding.EncodeToString(raw)

	// The prefix is only used to identify/search for the key.
	// Keep the actual secret out of it.
	prefix := key[:11]

	hash, err := crypto.HashPassword(key)
	if err != nil {
		return "", "", "", err
	}

	return key, prefix, hash, nil
}
