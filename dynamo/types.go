package dynamo

import (
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/skuid/helm-value-store/store"
)

type attributeValueMap map[string]*dynamodb.AttributeValue

func (avm *attributeValueMap) MarshalRelease() (*store.Release, error) {
	r := &store.Release{}
	for k, v := range *avm {
		switch strings.ToLower(k) {
		case "uniqueid":
			r.UniqueID = *v.S
		case "name":
			r.Name = *v.S
		case "chart":
			r.Chart = *v.S
		case "namespace":
			r.Namespace = *v.S
		case "version":
			r.Version = *v.S
		case "values":
			r.Values = *v.S
		case "labels":
			labels := map[string]string{}
			for label, value := range v.M {
				labels[label] = *value.S
			}
			r.Labels = labels
		}
	}
	return r, nil
}

func (avm *attributeValueMap) UnmarshalRelease(r store.Release) error {
	response := attributeValueMap{}

	st := reflect.TypeOf(r)
	v := reflect.ValueOf(r)
	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)
		fieldType := st.Field(i)
		if strings.Compare(fieldType.Name, "Labels") == 0 {
			labels := attributeValueMap{}
			for _, k := range fieldVal.MapKeys() {
				labels[k.String()] = &dynamodb.AttributeValue{S: aws.String(fieldVal.MapIndex(k).String())}
			}
			if len(labels) > 0 {
				response[st.Field(i).Name] = &dynamodb.AttributeValue{M: labels}
			}
		} else if fieldVal.Len() > 0 {
			response[st.Field(i).Name] = &dynamodb.AttributeValue{S: aws.String(fieldVal.String())}
		}
	}
	*avm = response

	return nil
}

type attributeValueMaps []attributeValueMap

func (avms *attributeValueMaps) MarshalReleases() (store.Releases, error) {
	rs := store.Releases{}
	for _, avm := range *avms {
		r, err := avm.MarshalRelease()
		if err != nil {
			return nil, err
		}
		rs = append(rs, *r)
	}
	return rs, nil
}

func (avms *attributeValueMaps) UnmarshalReleases(rs store.Releases) error {
	output := attributeValueMaps{}
	for _, r := range rs {
		avm := attributeValueMap{}
		err := avm.UnmarshalRelease(r)
		if err != nil {
			return err
		}
		output = append(output, avm)
	}
	*avms = output
	return nil
}
