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
