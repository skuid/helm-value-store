package cmd

import (
	"strings"
)

type selectorSet []string

func (ss *selectorSet) String() string {
	return strings.Join([]string(*ss), ",")
}

func (ss *selectorSet) Set(selector string) error {
	for _, pair := range strings.Split(selector, ",") {
		if len(pair) > 0 {
			*ss = append(*ss, pair)
		}
	}
	return nil
}

func (ss *selectorSet) Type() string {
	return "string"
}

func (ss selectorSet) ToMap() map[string]string {
	response := map[string]string{}
	for _, k := range ss {
		parts := strings.SplitN(k, "=", 2)
		if len(parts) > 1 {
			response[parts[0]] = parts[1]
		} else {
			response[parts[0]] = ""
		}
	}
	return response
}
