package update

import (
	"strconv"
	"strings"
)

// Version represents a semantic version (major.minor.patch)
type Version struct {
	Major int
	Minor int
	Patch int
	Raw   string
}

// ParseVersion parses a version string like "v1.2.3" or "1.2.3"
func ParseVersion(s string) Version {
	raw := s
	s = strings.TrimPrefix(s, "v")
	parts := strings.Split(s, ".")

	v := Version{Raw: raw}

	if len(parts) >= 1 {
		v.Major, _ = strconv.Atoi(parts[0])
	}
	if len(parts) >= 2 {
		v.Minor, _ = strconv.Atoi(parts[1])
	}
	if len(parts) >= 3 {
		patchStr := parts[2]
		if idx := strings.IndexAny(patchStr, "-+"); idx != -1 {
			patchStr = patchStr[:idx]
		}
		v.Patch, _ = strconv.Atoi(patchStr)
	}

	return v
}

// IsZero returns true if the version is unset
func (v Version) IsZero() bool {
	return v.Major == 0 && v.Minor == 0 && v.Patch == 0
}

// String returns the version as a string
func (v Version) String() string {
	if v.Raw != "" {
		return v.Raw
	}
	return "v" + strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Patch)
}

// Compare returns -1 if v < other, 0 if v == other, 1 if v > other
func (v Version) Compare(other Version) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}
	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}
	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}
	return 0
}

// LessThan returns true if v < other
func (v Version) LessThan(other Version) bool {
	return v.Compare(other) < 0
}

// IsDev returns true if the version appears to be a development version
func (v Version) IsDev() bool {
	raw := strings.ToLower(v.Raw)
	return raw == "dev" || raw == "" || strings.Contains(raw, "dev")
}
