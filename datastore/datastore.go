package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/datastore"
	"github.com/skuid/helm-value-store/store"
	"google.golang.org/api/option"
)

const kind = "hvsRelease"

// ReleaseStore stores and retrieves releases from a GCP Datastore table
type ReleaseStore struct {
	client *datastore.Client
}

type serviceAccount struct {
	ProjectID string `json:"project_id"`
}

// NewReleaseStore creates a new ReleaseStore
func NewReleaseStore(serviceAccountFile string) (*ReleaseStore, error) {
	data, err := ioutil.ReadFile(serviceAccountFile)
	if err != nil {
		return nil, fmt.Errorf("Error opening service account file: %q", err)
	}
	sa := &serviceAccount{}
	err = json.Unmarshal(data, sa)
	if err != nil {
		return nil, fmt.Errorf("Error parsing service account file: %q", err)
	}

	client, err := datastore.NewClient(context.Background(), sa.ProjectID, option.WithServiceAccountFile(serviceAccountFile))
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %q", err)
	}

	rs := &ReleaseStore{client: client}
	return rs, nil
}

// Get gets a release by it's UniqueID
func (rs ReleaseStore) Get(ctx context.Context, uniqueID string) (*store.Release, error) {
	response := &store.Release{}
	key := datastore.NameKey(kind, uniqueID, nil)
	if err := rs.client.Get(ctx, key, response); err != nil {
		return nil, fmt.Errorf("Error getting release: %q", err)
	}
	return response, nil
}

// Delete deletes a release by it's UniqueID
func (rs ReleaseStore) Delete(ctx context.Context, uniqueID string) error {
	key := datastore.NameKey(kind, uniqueID, nil)
	if err := rs.client.Delete(ctx, key); err != nil {
		return fmt.Errorf("Error deleting release: %q", err)
	}
	return nil
}

// Put creates or updates a release
func (rs ReleaseStore) Put(ctx context.Context, r store.Release) error {
	key := datastore.NameKey(kind, r.UniqueID, nil)
	if _, err := rs.client.Put(ctx, key, &r); err != nil {
		return fmt.Errorf("Error putting release: %q", err)
	}
	return nil
}

// List returns releases
func (rs ReleaseStore) List(ctx context.Context, selector map[string]string) (store.Releases, error) {

	releases := &store.Releases{}
	query := datastore.NewQuery(kind)
	_, err := rs.client.GetAll(ctx, query, releases)
	if err != nil {
		return nil, fmt.Errorf("Error getting releases: %q", err)
	}

	response := store.Releases{}
	for _, r := range *releases {
		if r.MatchesSelector(selector) {
			response = append(response, r)
		}
	}
	return response, nil
}

// Load bulk-writes releases to datastore
func (rs ReleaseStore) Load(ctx context.Context, releases store.Releases) error {
	keys := []*datastore.Key{}
	input := []*store.Release{}
	for i, r := range releases {
		keys = append(keys, datastore.NameKey(kind, r.UniqueID, nil))
		input = append(input, &releases[i])
	}

	if _, err := rs.client.PutMulti(context.Background(), keys, input); err != nil {
		return fmt.Errorf("Error loading releases: %q", err)
	}
	return nil
}

// Setup satisfies the RelaseStore interface. No action is required
func (rs ReleaseStore) Setup(ctx context.Context) error { return nil }
