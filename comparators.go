package semver

import "strconv"

func (v *Version) Compare(other *Version) int {
	if v.major != other.major {
		if v.major > other.major {
			return 1
		}
		return -1
	}

	if v.minor != other.minor {
		if v.minor > other.minor {
			return 1
		}
		return -1
	}

	if v.patch != other.patch {
		if v.patch > other.patch {
			return 1
		}
		return -1
	}

	if v.prerelease != nil || other.prerelease != nil {
		return comparePreRelease(v, other)
	}
	return 0
}

func comparePreRelease(main, other *Version) int {

	if len(main.prerelease) != 0 {
		if len(other.prerelease) != 0 {
			for i := 0; i < len(main.prerelease) && i < len(other.prerelease); i++ {
				if main.prerelease[i] != other.prerelease[i] {
					if x, err := strconv.ParseUint(main.prerelease[i], 10, 0); err == nil {
						if y, err := strconv.ParseUint(other.prerelease[i], 10, 0); err == nil {
							if x > y {
								return 1
							} else {
								return -1
							}
						}
					}
					if main.prerelease[i] > other.prerelease[i] {
						return 1
					} else {
						return -1
					}
				}
			}
		}
		return -1
	}
	if len(other.prerelease) != 0 {
		return 1
	}

	return 0
}
