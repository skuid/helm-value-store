package store_test

import (
	"reflect"
	"testing"

	"cloud.google.com/go/datastore"
	"github.com/skuid/helm-value-store/store"
)

func TestReleaseLoad(t *testing.T) {
	cases := []struct {
		name           string
		release        store.Release
		properties     []datastore.Property
		expectedLabels map[string]string
		hasError       bool
	}{
		{
			"Contains labels",
			store.Release{},
			[]datastore.Property{{Name: "labels", Value: []byte(`{"region":"us"}`), NoIndex: true}},
			map[string]string{"region": "us"},
			false,
		},
		{
			"Empty label map",
			store.Release{},
			[]datastore.Property{{Name: "labels", Value: []byte(`{}`), NoIndex: true}},
			map[string]string{},
			false,
		},
		{
			"Expect json parse error",
			store.Release{},
			[]datastore.Property{{Name: "labels", Value: []byte(``), NoIndex: true}},
			map[string]string{},
			true,
		},
	}

	for _, c := range cases {
		err := c.release.Load(c.properties)
		if (err != nil) != c.hasError {
			t.Errorf("Got unexpected result: err = %v, expected error: %t", err, c.hasError)
		}
		if !reflect.DeepEqual(c.release.Labels, c.expectedLabels) {
			if len(c.release.Labels) == 0 && len(c.expectedLabels) == 0 {
				continue
			}
			t.Errorf("Labels do not match. Expected %v, got %v", c.expectedLabels, c.release.Labels)
		}

	}

}

func TestReleaseSave(t *testing.T) {
	cases := []struct {
		name             string
		release          store.Release
		expectedProperty datastore.Property
		hasError         bool
	}{
		{
			"Contains labels",
			store.Release{Labels: map[string]string{"region": "us"}},
			datastore.Property{Name: "ReleaseLabels", Value: []byte(`{"region":"us"}`), NoIndex: true},
			false,
		},
		{
			"Empty label map",
			store.Release{Labels: map[string]string{}},
			datastore.Property{Name: "ReleaseLabels", Value: []byte(`{}`), NoIndex: true},
			false,
		},
		{
			"nil label map",
			store.Release{Labels: nil},
			datastore.Property{Name: "ReleaseLabels", Value: []byte(`null`), NoIndex: true},
			false,
		},
		{
			"ReleaseLabels exist",
			store.Release{ReleaseLabels: []byte(`{"key": "value"}`)},
			datastore.Property{Name: "ReleaseLabels", Value: []byte(`null`), NoIndex: true},
			false,
		},
	}

	for _, c := range cases {
		props, err := c.release.Save()
		if (err != nil) != c.hasError {
			t.Errorf("Test '%s' got unexpected result: err = %v, expected error: %t", c.name, err, c.hasError)
		}
		if len(props) == 0 {
			t.Errorf("Test %s got no properties", c.name)
			continue
		}
		if !reflect.DeepEqual(props[len(props)-1], c.expectedProperty) {
			t.Errorf("Test '%s': Properties do not match. Expected %v, got %v", c.name, c.expectedProperty, props[len(props)-1])
		}

	}

}

func TestMatchesSelector(t *testing.T) {
	cases := []struct {
		release  store.Release
		selector map[string]string
		want     bool
	}{
		{
			store.Release{Labels: map[string]string{"region": "us", "environment": "test"}},
			map[string]string{"region": "us", "environment": "test"},
			true,
		},
		{
			store.Release{Labels: map[string]string{"region": "us", "environment": "test"}},
			map[string]string{"region": "us", "environment": ""},
			true,
		},
		{
			store.Release{Labels: map[string]string{"region": "us", "environment": ""}},
			map[string]string{"region": "us", "environment": "test"},
			false,
		},
		{
			store.Release{Labels: map[string]string{"region": "us", "environment": "test"}},
			map[string]string{"region": "us", "environment": "test"},
			true,
		},
		{
			store.Release{Labels: map[string]string{}},
			map[string]string{"region": "us", "environment": "test"},
			false,
		},
		{
			store.Release{},
			map[string]string{"region": "us", "environment": "test"},
			false,
		},
		{
			store.Release{Labels: map[string]string{"region": "us", "environment": "test"}},
			map[string]string{},
			true,
		},
	}

	for _, c := range cases {
		got := c.release.MatchesSelector(c.selector)
		if got != c.want {
			t.Errorf("Failed %#v.MatchesSelector(%v): Expected %t, got %t", c.release, c.selector, c.want, got)

		}
	}
}

func TestMergeValues(t *testing.T) {
	cases := []struct {
		release   store.Release
		setValues []string
		want      string
	}{
		{
			store.Release{Values: "foo: 42"},
			[]string{"foo=24"},
			"foo: 24\n",
		},
		{
			store.Release{Values: "bar: value1\nfoo: value2"},
			[]string{"bar=value2", "foo=value1"},
			"bar: value2\nfoo: value1\n",
		},
		{
			store.Release{Values: "bar: value1\nfoo: value2"},
			[]string{"bar=value2,foo=value1"},
			"bar: value2\nfoo: value1\n",
		},
		{
			store.Release{Values: "foo:\n  bar: value1"},
			[]string{"foo.bar=value2"},
			"foo:\n  bar: value2\n",
		},
		{
			store.Release{Values: "foo:\n  bar: value1\n  baz: value2"},
			[]string{"foo.bar=value2"},
			"foo:\n  bar: value2\n  baz: value2\n",
		},
	}

	for _, c := range cases {
		err := c.release.MergeValues(c.setValues)
		if err != nil {
			t.Errorf("Error running %#v.MergeValues(%v): %s", c.release, c.setValues, err)
		}
		got := c.release.Values
		if got != c.want {
			t.Errorf("Failed %#v.MergeValues(%v): Expected %s, got %s", c.release, c.setValues, c.want, got)
		}
	}
}
