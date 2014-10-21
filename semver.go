package semver

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	numbers    string = "0123456789"
	letters           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-"
	alphanum          = letters + numbers
	dot               = "."
	hyphen            = "-"
	plus              = "+"
	delimiters        = dot + hyphen + plus
	allchars          = alphanum + delimiters
)

type VersionGetter interface {
	Major() uint64
	Minor() uint64
	Patch() uint64
	Prerelease() string
	PrerelaseIdentifiers() []string
	Metadata() string
	MetadataIdentifiers() []string
}

type VersionSetter interface {
	SetPrerelease(...string) error
	SetMetadata(...string) error
	SetMajor(uint64)
	SetMinor(uint64)
	SetPatch(uint64)
}

type VersionGetterSetter interface {
	VersionGetter
	VersionSetter
}

type Version struct {
	major      uint64
	minor      uint64
	patch      uint64
	prerelease *prereleases
	metadata   []string
}

type prereleases struct {
	values  []string
	numbers map[int]uint64
}

func Build(major, minor, patch uint64, extra ...[]string) *Version {
	if len(extra) == 1 {
		ver := &Version{major, minor, patch, nil, nil}
		ver.SetPrerelease(extra[0]...)
		return ver
	}
	if len(extra) > 1 {
		ver := &Version{major, minor, patch, nil, extra[1]}
		ver.SetPrerelease(extra[0]...)
		return ver
	}
	return &Version{major, minor, patch, nil, nil}
}

func New(version string) (*Version, error) {
	var versions []string
	var prereleases []string
	var metadatas []string

	result := new(Version)

	if strings.Contains(version, plus) {
		metadata := strings.Split(version, plus)
		metadatas = strings.Split(metadata[1], dot)

		if err := result.SetMetadata(metadatas...); err != nil {
			return nil, err
		}

		version = metadata[0]
	}

	if strings.Contains(version, hyphen) {
		prerelease := strings.Split(version, hyphen)
		prereleases = strings.Split(prerelease[1], dot)

		if err := result.SetPrerelease(prereleases...); err != nil {
			return nil, err
		}

		version = prerelease[0]
	}

	versions = strings.Split(version, dot)
	if len(versions) != 3 {
		return nil, errors.New("major.minor.patch pattern not found")
	}

	var versionNumbers [3]uint64
	for i, partial := range versions {

		if num, err := strconv.ParseUint(partial, 10, 0); err != nil {
			return nil, errors.New(fmt.Sprint("expected unsigned integer: ", partial))
		} else {
			if hasLeadingZero(partial) {
				return nil, errors.New(fmt.Sprint("leading zeroes in version number: ", partial))
			}
			versionNumbers[i] = num
		}
	}

	result.major = versionNumbers[0]
	result.minor = versionNumbers[1]
	result.patch = versionNumbers[2]

	return result, nil
}

func (v Version) String() string {
	var buffer bytes.Buffer
	w := bufio.NewWriter(&buffer)

	fmt.Fprintf(w, "%d.%d.%d", v.major, v.minor, v.patch)

	if v.prerelease != nil {
		fmt.Fprintf(w, "%v%v", hyphen, v.Prerelease())
	}

	if v.metadata != nil {
		fmt.Fprintf(w, "%v%v", plus, v.Metadata())
	}

	w.Flush()
	return buffer.String()
}

func (v Version) Major() uint64 {
	return v.major
}

func (v *Version) SetMajor(major uint64) {
	v.major = major
}

func (v Version) Minor() uint64 {
	return v.minor
}

func (v *Version) SetMinor(minor uint64) {
	v.minor = minor
}

func (v Version) Patch() uint64 {
	return v.patch
}

func (v *Version) SetPatch(patch uint64) {
	v.patch = patch
}

func (v Version) Prerelease() string {
	return strings.Join(v.prerelease.values, dot)
}

func (v *Version) SetPrerelease(identifiers ...string) error {
	var result []string
	numbers := make(map[int]uint64)

	for i, ident := range identifiers {
		if len(ident) < 1 {
			return errors.New("identifier is empty")
		}

		if num, err := strconv.ParseUint(ident, 10, 0); err == nil {
			if hasLeadingZero(ident) {
				return errors.New(fmt.Sprint("leading zeroes in numerical identifier: ", ident))
			}
			numbers[i] = num
			result = append(result, ident)
		} else {
			if !containsOnly(ident, alphanum) {
				return errors.New(fmt.Sprint("not alphanumerical: ", ident))
			}
			result = append(result, ident)
		}
	}
	pre := &prereleases{result, numbers}
	v.prerelease = pre
	return nil
}

func (v Version) Metadata() string {
	return strings.Join(v.metadata, dot)
}

func (v *Version) SetMetadata(identifiers ...string) error {
	var result []string

	for _, ident := range identifiers {
		if len(ident) < 1 {
			return errors.New("identifier is empty")
		}

		if !containsOnly(ident, alphanum) {
			return errors.New(fmt.Sprint("not alphanumerical: ", ident))
		}
		result = append(result, ident)
	}

	v.metadata = result
	return nil
}

func (v Version) Satifies(other string) bool {
	return true
}
