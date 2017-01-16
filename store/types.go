package store

import (
	"fmt"
	"os"

	"k8s.io/helm/pkg/helm"
	rls "k8s.io/helm/pkg/proto/hapi/services"
	"strings"
)

var client *helm.Client

func init() {
	client = helm.NewClient(helm.Host(os.Getenv("TILLER_HOST")))
}

// A Release contains metadata about a release of a healm chart
type Release struct {
	UniqueID  string            `json:"unique_id"`
	Labels    map[string]string `json:"labels"`
	Name      string            `json:"name"`
	Chart     string            `json:"chart"`
	Namespace string            `json:"namespace"`
	Version   string            `json:"version"`
	Values    string            `json:"values"`
}

// MatchesSelector checks if the specified release contains all the key/value pairs in it's Labels
func (r Release) MatchesSelector(selector map[string]string) bool {
	if (r.Labels == nil || len(r.Labels) == 0) && len(selector) > 0 {
		return false
	}
	for selectorK, selectorV := range selector {
		labelValue, ok := r.Labels[selectorK]
		if !ok || (strings.Compare(labelValue, selectorV) != 0 && len(selectorV) > 0) {
			return false
		}
	}
	return true
}

// ReleaseUnmarshaler is an interface for unmarshaling a release
type ReleaseUnmarshaler interface {
	UnmarshalRelease(Release) error
}

// ReleaseMarshaler is an interface for marshaling a release
type ReleaseMarshaler interface {
	MarshalRelease() (*Release, error)
}

// Upgrade sends an update to an existing release in a cluster
func (r Release) Upgrade(dryRun bool, timeout int64) error {
	resp, err := client.UpdateRelease(
		r.Name,
		r.Chart,
		helm.UpdateValueOverrides([]byte(r.Values)),
		helm.UpgradeDryRun(dryRun),
		helm.UpgradeTimeout(timeout),
	)
	defer fmt.Println(resp)
	return err
}

// Install creates an new release in a cluster
func (r Release) Install(dryRun bool, timeout int64) error {
	resp, err := client.InstallRelease(
		r.Chart,
		r.Namespace,
		helm.ValueOverrides([]byte(r.Values)),
		helm.ReleaseName(r.Name),
		helm.InstallDryRun(dryRun),
		helm.InstallTimeout(timeout),
	)
	fmt.Println(resp)
	return err
}

// Get the release content from Tiller
func (r Release) Get() (*rls.GetReleaseContentResponse, error) {
	return client.ReleaseContent(r.Name)
}

// Releases is a slice of release
type Releases []Release

// ReleasesUnmarshaler is an interface for unmarshaling slices of release
type ReleasesUnmarshaler interface {
	UnmarshalReleases(Releases) error
}

// ReleasesMarshaler is an interface for marshaling slices of release
type ReleasesMarshaler interface {
	MarshalReleases() (Releases, error)
}

// A ReleaseStore is a backend that stores releases
type ReleaseStore interface {
	Get(uniqueID string) (*Release, error)
	Put(Release) error

	List(selector map[string]string) (Releases, error)
	Load(Releases) error
}
