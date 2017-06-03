package hosts

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {
	testCases := []struct {
		input    string
		expected Entry
	}{
		{" ", Entry{}},
		{"  ## line with comment only", Entry{Comment: "## line with comment only"}},
		{"ip_address hostname", Entry{Address: "ip_address", Hostname: "hostname"}},
		{"ip_address hostname #comment", Entry{"ip_address", "hostname", nil, "#comment"}},
		{"ip_address hostname alias1 alias2", Entry{"ip_address", "hostname", []string{"alias1", "alias2"}, ""}},
		{"ip_address hostname alias #comment", Entry{"ip_address", "hostname", []string{"alias"}, "#comment"}},
	}

	for _, testCase := range testCases {
		actual := parseLineEntry(testCase.input)
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("Actual value %v is not the expected %v", actual, testCase.expected)
		}
	}
}
