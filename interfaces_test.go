package semver

import (
	"encoding/json"
	"testing"
)

type container struct {
	Ver *Version `json:"name"`
}

var jsonString = []byte(`{"name":"1.2.3"}`)
var textString = "1.2.3"

func TestUnMarshal(t *testing.T) {

	ver := new(container)

	if err := json.Unmarshal(jsonString, &ver); err != nil {
		t.Fatalf("UnmarshalJSON(%q) => %q, want %v", jsonString, err, textString)
	}
	if ver.Ver.String() != textString {
		t.Errorf("UnmarshalJSON(%q) => %v, want %v", jsonString, ver, textString)
	}
}

func TestMarshal(t *testing.T) {
	ver := new(container)
	ver.Ver, _ = New(textString)

	result, err := json.Marshal(ver)
	if err != nil {
		t.Fatalf("MarshalJSON(%q) => %q, want %v", ver, err, string(jsonString))
	}

	if string(result) != string(jsonString) {
		t.Errorf("MarshalJSON(%q) => %v, want %v", ver, string(result), string(jsonString))
	}
}
