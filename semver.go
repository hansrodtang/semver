// Package semver provides a Semantic Versioning library for Go. It allows you to parse and compare semver version strings.
// Covers version 2.0.0 of the semver specification.
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
	dot        = "."
	hyphen     = "-"
	plus       = "+"
	delimiters = dot + hyphen + plus
)

// Version is the container for semver version data.
type Version struct {
	major      uint64
	minor      uint64
	patch      uint64
	prerelease *prereleases
	metadata   []string
}

type prereleases struct {
	values  []string
	numbers []int
}

// Build accepts version numbers in uint64 and optional prerelease and metadata information in a string array.
// Build circumvents error checking, and is mostly used for testing.
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

// New accepts a valid semver version string and returns a Version struct.
// Returns error if the supplied string is an invalid semver version.
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
		prerelease := strings.SplitN(version, hyphen, 2)
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

		num, err := strconv.ParseUint(partial, 10, 0)
		if err != nil {
			return nil, errors.New(fmt.Sprint("expected unsigned integer: ", partial))
		}
		if hasLeadingZero(partial) {
			return nil, errors.New(fmt.Sprint("leading zeroes in version number: ", partial))
		}
		versionNumbers[i] = num
	}

	result.major = versionNumbers[0]
	result.minor = versionNumbers[1]
	result.patch = versionNumbers[2]

	return result, nil
}

// String returns a valid semver string based on the data contained in Version.
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

// Major returns the major version.
func (v Version) Major() uint64 {
	return v.major
}

// SetMajor accepts a uint64 to change the currently set major version.
func (v *Version) SetMajor(major uint64) {
	v.major = major
}

// IncrementMajor increases Major version by 1.
func (v *Version) IncrementMajor() {
	v.major++
}

// DecrementMajor decreases Major version by 1.
func (v *Version) DecrementMajor() {
	v.major--
}

// Minor returns the minor version.
func (v Version) Minor() uint64 {
	return v.minor
}

// SetMinor accepts a uint64 to change the currently set minor version.
func (v *Version) SetMinor(minor uint64) {
	v.minor = minor
}

// IncrementMinor increases Minor version by 1.
func (v *Version) IncrementMinor() {
	v.minor++
}

// DecrementMinor decreases Minor version by 1.
func (v *Version) DecrementMinor() {
	v.minor--
}

// Patch returns the patch version.
func (v Version) Patch() uint64 {
	return v.patch
}

// SetPatch accepts a uint64 to change the currently set patch version.
func (v *Version) SetPatch(patch uint64) {
	v.patch = patch
}

// IncrementPatch increases Patch version by 1.
func (v *Version) IncrementPatch() {
	v.patch++
}

// DecrementPatch decreases Patch version by 1.
func (v *Version) DecrementPatch() {
	v.patch--
}

// Prerelease returns the prerelease identifiers as a dot seperated string.
func (v Version) Prerelease() string {
	return strings.Join(v.prerelease.values, dot)
}

// SetPrerelease accepts a series of strings to form the prerelease identifiers.
// Returns error if any of the supplied strings aren't a valid prerelease identifier.
func (v *Version) SetPrerelease(identifiers ...string) error {
	var result []string
	var numbers []int

	for _, ident := range identifiers {
		if len(ident) < 1 {
			return errors.New("identifier is empty")
		}

		if num, err := strconv.ParseUint(ident, 10, 0); err == nil {
			if hasLeadingZero(ident) {
				return errors.New(fmt.Sprint("leading zeroes in numerical identifier: ", ident))
			}
			numbers = append(numbers, int(num))
			result = append(result, ident)
		} else {
			if !containsOnly(ident, alphanumeric) {
				return errors.New(fmt.Sprint("not alphanumerical: ", ident))
			}
			numbers = append(numbers, -1)
			result = append(result, ident)
		}
	}
	pre := &prereleases{result, numbers}
	v.prerelease = pre
	return nil
}

// Metadata returns the metadata identifiers as a dot seperated string.
func (v Version) Metadata() string {
	return strings.Join(v.metadata, dot)
}

// SetMetadata accepts a series of strings to form the metadata identifiers.
// Returns error if any of the supplied strings aren't a valid metadata identifier.
func (v *Version) SetMetadata(identifiers ...string) error {
	var result []string

	for _, ident := range identifiers {
		if len(ident) < 1 {
			return errors.New("identifier is empty")
		}

		if !containsOnly(ident, alphanumeric) {
			return errors.New(fmt.Sprint("not alphanumerical: ", ident))
		}
		result = append(result, ident)
	}

	v.metadata = result
	return nil
}

// Satifies accepts a set of comparators and version numbers as a string.
// The syntax for the string is documented here: https://www.npmjs.org/doc/misc/semver.html
// Returns true if Version matches the comparators, false if it does not.
func (v Version) Satifies(requirements string) (bool, error) {
	return false, nil
}
