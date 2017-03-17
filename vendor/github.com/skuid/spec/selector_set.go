package spec

import (
	"strings"
)

// SelectorSet is a custom flag that takes key-value pairs joined by equal
// signs and separated by commas. The flag may be specified multiple times and
// the values will be merged.
//
// Ex:
//
//     mycmd -flag key=value,key2=value2 -flag key3=value3
//
// Calling ToMap() results in:
//
//    map[string]string{
//        "key": "value",
//        "key2": "value2",
//        "key3": "value3",
//    }
//
type SelectorSet []string

// String satisfies the flag.Value interface
func (ss *SelectorSet) String() string {
	return strings.Join([]string(*ss), ",")
}

// Set satisfies the flag.Value interface
func (ss *SelectorSet) Set(selector string) error {
	for _, pair := range strings.Split(selector, ",") {
		if len(pair) > 0 {
			*ss = append(*ss, pair)
		}
	}
	return nil
}

// Type returns the type of the flag as a string
func (ss *SelectorSet) Type() string {
	return "string"
}

// ToMap returns a map representation of the SelectorSet
func (ss SelectorSet) ToMap() map[string]string {
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
