package semver

import (
	"encoding/json"
	"strings"
)

func (v *Version) UnmarshalJSON(b []byte) error {
	input := strings.Trim(string(b), "\"")
	ver, err := New(input)
	*v = *ver
	return err
}

func (p *Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

type Versions []*Version

func (v Versions) Len() int {
	return len(v)
}

func (v Versions) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v Versions) Less(i, j int) bool {
	return v[i].Compare(v[j]) < 0
}
