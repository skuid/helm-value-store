package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/kubernetes/helm/cmd/helm/downloader"
	"k8s.io/helm/cmd/helm/helmpath"
	"k8s.io/helm/pkg/helm"
	rls "k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/strvals"
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

func (r Release) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s", r.UniqueID, r.Name, r.Chart, r.Version)
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
func (r Release) Upgrade(chartLocation string, dryRun bool, timeout int64) (*rls.UpdateReleaseResponse, error) {
	return client.UpdateRelease(
		r.Name,
		chartLocation,
		helm.UpdateValueOverrides([]byte(r.Values)),
		helm.UpgradeDryRun(dryRun),
		helm.UpgradeTimeout(timeout),
	)
}

// Install creates an new release in a cluster
func (r Release) Install(chartLocation string, dryRun bool, timeout int64) (*rls.InstallReleaseResponse, error) {
	return client.InstallRelease(
		chartLocation,
		r.Namespace,
		helm.ValueOverrides([]byte(r.Values)),
		helm.ReleaseName(r.Name),
		helm.InstallDryRun(dryRun),
		helm.InstallTimeout(timeout),
	)
}

// Download gets the release from an index server
func (r Release) Download() (string, error) {
	dl := downloader.ChartDownloader{
		HelmHome: helmpath.Home(os.Getenv("HELM_HOME")),
		Out:      os.Stdout,
	}

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}
	filename, _, err := dl.DownloadTo(r.Chart, r.Version, tmpDir)

	if err == nil {
		lname, err := filepath.Abs(filename)
		if err != nil {
			return filename, err
		}
		return lname, nil
	}

	return filename, fmt.Errorf("file %q not found", r.Chart)
}

// Get the release content from Tiller
func (r Release) Get() (*rls.GetReleaseContentResponse, error) {
	return client.ReleaseContent(r.Name)
}

// MergeValues parses string values and then merges them into the
// existing Values for a release.
// Adopted from kubernetes/helm/cmd/helm/install.go
func (r *Release) MergeValues(values []string) error {
	base := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(r.Values), &base); err != nil {
		return fmt.Errorf("Error parsing values for release %s: %s", r.Name, err)
	}

	// User specified a value via --set
	for _, value := range values {
		if err := strvals.ParseInto(value, base); err != nil {
			return fmt.Errorf("failed parsing --set data: %s", err)
		}
	}

	mergedValues, err := yaml.Marshal(base)
	if err != nil {
		return fmt.Errorf("Error parsing merged values for release %s: %s", r.Name, err)
	}
	r.Values = string(mergedValues)

	return nil
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
	Delete(uniqueID string) error

	List(selector map[string]string) (Releases, error)
	Load(Releases) error
	Setup() error
}
