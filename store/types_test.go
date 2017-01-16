package store_test

import (
	"github.com/skuid/helm-value-store/store"
	"testing"
)

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
			t.Errorf("Failed %#v.MatchesSelector(%v): Expected %s, got %s", c.release, c.selector, c.want, got)

		}
	}
}
