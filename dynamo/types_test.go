package dynamo

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/skuid/helm-value-store/store"
)

func TestMarshalRelease(t *testing.T) {
	cases := []struct {
		avm  attributeValueMap
		want *store.Release
	}{
		{
			attributeValueMap{
				"UniqueID":  {S: aws.String("abc123")},
				"Name":      {S: aws.String("prom1")},
				"Chart":     {S: aws.String("skuid/prometheus")},
				"Namespace": {S: aws.String("default")},
				"Version":   {S: aws.String("0.1.3")},
				"Values":    {S: aws.String("image: prometheus\ntag: latest")},
				"Labels": {M: attributeValueMap{
					"region":      &dynamodb.AttributeValue{S: aws.String("us")},
					"environment": &dynamodb.AttributeValue{S: aws.String("test")}}},
				"Unused": {S: aws.String("nothing")},
			},
			&store.Release{
				UniqueID:  "abc123",
				Labels:    map[string]string{"region": "us", "environment": "test"},
				Name:      "prom1",
				Chart:     "skuid/prometheus",
				Namespace: "default",
				Version:   "0.1.3",
				Values:    "image: prometheus\ntag: latest",
			},
		},
		{
			attributeValueMap{
				"UniqueID": {S: aws.String("abc123")},
				"Name":     {S: aws.String("prom1")},
			},
			&store.Release{
				UniqueID: "abc123",
				Name:     "prom1",
			},
		},
	}

	for _, c := range cases {
		got, err := c.avm.MarshalRelease()
		if err != nil {
			t.Errorf("Errored marshaling release %s", err.Error())
		} else if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Failed marshaling release: Expected \n\t%#v\ngot: \n\t%#v", c.want, got)
		}
	}
}

func TestUnmarshalRelease(t *testing.T) {
	cases := []struct {
		release *store.Release
		want    attributeValueMap
	}{
		{
			&store.Release{
				UniqueID:  "abc123",
				Labels:    map[string]string{"region": "us", "environment": "test"},
				Name:      "prom1",
				Chart:     "skuid/prometheus",
				Namespace: "default",
				Version:   "0.1.3",
				Values:    "image: prometheus\ntag: latest",
			},
			attributeValueMap{
				"UniqueID":  {S: aws.String("abc123")},
				"Name":      {S: aws.String("prom1")},
				"Chart":     {S: aws.String("skuid/prometheus")},
				"Namespace": {S: aws.String("default")},
				"Version":   {S: aws.String("0.1.3")},
				"Values":    {S: aws.String("image: prometheus\ntag: latest")},
				"Labels": {M: attributeValueMap{
					"region":      &dynamodb.AttributeValue{S: aws.String("us")},
					"environment": &dynamodb.AttributeValue{S: aws.String("test")}}},
			},
		},
		{
			&store.Release{
				UniqueID: "abc123",
				Name:     "prom1",
			},
			attributeValueMap{
				"UniqueID": {S: aws.String("abc123")},
				"Name":     {S: aws.String("prom1")},
			},
		},
	}

	for _, c := range cases {
		got := attributeValueMap{}
		err := got.UnmarshalRelease(*c.release)
		if err != nil {
			t.Errorf("Errored marshaling release %s", err.Error())
		} else if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Failed unmarshaling release: Expected \n\t%#v\ngot: \n\t%#v", c.want, got)
		}
	}
}

func TestMarshalReleases(t *testing.T) {
	cases := []struct {
		avm  attributeValueMaps
		want store.Releases
	}{
		{
			attributeValueMaps{{
				"UniqueID":  {S: aws.String("abc123")},
				"Name":      {S: aws.String("prom1")},
				"Chart":     {S: aws.String("skuid/prometheus")},
				"Namespace": {S: aws.String("default")},
				"Version":   {S: aws.String("0.1.3")},
				"Values":    {S: aws.String("image: prometheus\ntag: latest")},
				"Labels": {M: attributeValueMap{
					"region":      &dynamodb.AttributeValue{S: aws.String("us")},
					"environment": &dynamodb.AttributeValue{S: aws.String("test")}}},
				"Unused": {S: aws.String("nothing")},
			}},
			store.Releases{{
				UniqueID:  "abc123",
				Labels:    map[string]string{"region": "us", "environment": "test"},
				Name:      "prom1",
				Chart:     "skuid/prometheus",
				Namespace: "default",
				Version:   "0.1.3",
				Values:    "image: prometheus\ntag: latest",
			}},
		},
		{
			attributeValueMaps{{
				"UniqueID": {S: aws.String("abc123")},
				"Name":     {S: aws.String("prom1")},
			}},
			store.Releases{{
				UniqueID: "abc123",
				Name:     "prom1",
			}},
		},
	}

	for _, c := range cases {
		got, err := c.avm.MarshalReleases()
		if err != nil {
			t.Errorf("Errored marshaling releases %s", err.Error())
		} else if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Failed marshaling releases: Expected \n%#v\ngot: \n%#v", c.want, got)
		}
	}
}

func TestUnmarshalReleases(t *testing.T) {
	cases := []struct {
		releases store.Releases
		want     attributeValueMaps
	}{
		{
			store.Releases{{
				UniqueID: "abc123",
				Name:     "prom1",
			}},
			attributeValueMaps{{
				"UniqueID": {S: aws.String("abc123")},
				"Name":     {S: aws.String("prom1")},
			}},
		},
	}

	for _, c := range cases {
		got := attributeValueMaps{}
		err := got.UnmarshalReleases(c.releases)
		if err != nil {
			t.Errorf("Errored marshaling releases %s", err.Error())
		}
		for i := range got {
			if !reflect.DeepEqual(got[i], c.want[i]) {
				t.Errorf("Failed unmarshaling releases: Expected \n%#v\ngot: \n%#v", c.want, got)
			}
		}

	}
}
