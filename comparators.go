package semver

// Compare accepts a Version and compares itself against it, returning >0 for greater than, 0 for equals and <0 for less than.
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

	if v.prerelease == nil {
		if other.prerelease == nil {
			return 0
		}
		return 1
	}
	if other.prerelease == nil {
		return -1
	}

	return v.prerelease.compare(other.prerelease)
}

func (p *prereleases) compare(other *prereleases) int {
	for i := 0; i < len(p.values) && i < len(other.values); i++ {
		if p.values[i] != other.values[i] {
			if val1 := p.numbers[i]; val1 > 0 {
				if val2 := other.numbers[i]; val2 > 0 {
					if val1 > val2 {
						return 1
					}
					return -1
				}
			}
			if p.values[i] > other.values[i] {
				return 1
			}
			return -1
		}
	}

	if len(p.values) == len(other.values) {
		return 0
	} else if len(p.values) < len(other.values) {
		return -1
	}
	return 1
}
