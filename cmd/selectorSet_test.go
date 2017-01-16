package cmd

import (
	"reflect"
	"strings"
	"testing"
)

func TestSelectorSetSet(t *testing.T) {
	cases := []struct {
		flags    []string
		selector selectorSet
		want     string
	}{
		{
			[]string{"env=test", "region=us2"},
			selectorSet{},
			"env=test,region=us2",
		},
		{
			[]string{"env=test", "region=eu1", "region=us2"},
			selectorSet{},
			"env=test,region=eu1,region=us2",
		},
	}

	for _, c := range cases {
		for _, f := range c.flags {
			c.selector.Set(f)
		}
		if strings.Compare(c.selector.String(), c.want) != 0 {
			t.Errorf("Expected %s, got %s", c.want, c.selector.String())
		}
	}
}

func TestSelectorSetType(t *testing.T) {
	want := "string"
	selector := selectorSet{}
	if strings.Compare(selector.Type(), want) != 0 {
		t.Errorf("Expected %s, got %s", want, selector.Type())
	}
}

func TestSelectorSetToMap(t *testing.T) {
	cases := []struct {
		flags    []string
		selector selectorSet
		want     map[string]string
	}{
		{
			[]string{"env=test", "region=eu1"},
			selectorSet{},
			map[string]string{"env": "test", "region": "eu1"},
		},
		{
			[]string{"env=test", "region=eu1", "region=us2"},
			selectorSet{},
			map[string]string{"env": "test", "region": "us2"},
		},
		{
			[]string{"env", "region=us2"},
			selectorSet{},
			map[string]string{"env": "", "region": "us2"},
		},
	}

	for _, c := range cases {
		for _, f := range c.flags {
			c.selector.Set(f)
		}
		if !reflect.DeepEqual(c.want, c.selector.ToMap()) {
			t.Errorf("Expected %s, got %s", c.want, c.selector.ToMap())
		}
	}
}
