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
