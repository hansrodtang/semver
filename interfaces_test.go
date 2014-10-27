package semver

import (
	"encoding/json"
	"sort"
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

var unsortedVersions = []string{
	"1.2.3-alpha.2+123.456",
	"1.2.3-alpha.1",
	"1.4.3+123.456",
	"1.3.3-alpha.b-eta+123.b-uild",
	"1.2.3+123.b-uild",
	"1.2.3-alpha.b-eta",
	"0.0.1-alpha.preview+123.456",
}

var sortedVersions = []string{
	"0.0.1-alpha.preview+123.456",
	"1.2.3-alpha.1",
	"1.2.3-alpha.2+123.456",
	"1.2.3-alpha.b-eta",
	"1.2.3+123.b-uild",
	"1.3.3-alpha.b-eta+123.b-uild",
	"1.4.3+123.456",
}

func TestSorter(t *testing.T) {
	var unsortedlist Versions
	for _, v := range unsortedVersions {
		ver, _ := New(v)
		unsortedlist = append(unsortedlist, ver)
	}

	sort.Sort(unsortedlist)

	var sortedlist Versions
	for _, v := range sortedVersions {
		ver, _ := New(v)
		sortedlist = append(sortedlist, ver)
	}

	for i, r := range unsortedlist {
		result := r.String()
		expected := sortedlist[i].String()

		if result != expected {
			t.Errorf("sort.Sort() => %q, want %v", result, expected)
		}
	}

}
