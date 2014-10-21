package semver

import (
	"encoding/json"
	"strings"
)

func (v *Version) UnmarshalJSON(b []byte) error {
	input := strings.Trim(string(b), "\"")
	ver, err := New(input)
	if err != nil {
		return err
	}
	*v = *ver
	return nil
}

func (p *Version) MarshalJSON() ([]byte, error) {
	output, err := json.Marshal(p.String())
	if err != nil {
		return nil, err
	}
	return output, nil
}
