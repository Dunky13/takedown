package releaseflow

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var versionPattern = regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)

// Version is a SemVer version with a core compatible with native installers.
type Version struct {
	Major, Minor, Patch int
	Prerelease, Build   string
}

func ParseVersion(value string) (Version, error) {
	matches := versionPattern.FindStringSubmatch(value)
	if matches == nil {
		return Version{}, fmt.Errorf("invalid SemVer %q (omit the v prefix)", value)
	}
	var core [3]int
	for i := range core {
		n, err := strconv.Atoi(matches[i+1])
		if err != nil || n > 65535 {
			return Version{}, fmt.Errorf("version components must be between 0 and 65535: %q", value)
		}
		core[i] = n
	}
	for _, part := range strings.Split(matches[4], ".") {
		if numeric(part) && len(part) > 1 && part[0] == '0' {
			return Version{}, fmt.Errorf("numeric prerelease identifiers cannot have leading zeros: %q", value)
		}
	}
	return Version{core[0], core[1], core[2], matches[4], matches[5]}, nil
}

func (v Version) Core() string { return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch) }

func (v Version) String() string {
	value := v.Core()
	if v.Prerelease != "" {
		value += "-" + v.Prerelease
	}
	if v.Build != "" {
		value += "+" + v.Build
	}
	return value
}

// Compare follows SemVer precedence; build metadata does not affect ordering.
func (v Version) Compare(other Version) int {
	for _, pair := range [][2]int{{v.Major, other.Major}, {v.Minor, other.Minor}, {v.Patch, other.Patch}} {
		if pair[0] < pair[1] {
			return -1
		}
		if pair[0] > pair[1] {
			return 1
		}
	}
	if v.Prerelease == other.Prerelease {
		return 0
	}
	if v.Prerelease == "" {
		return 1
	}
	if other.Prerelease == "" {
		return -1
	}
	a, b := strings.Split(v.Prerelease, "."), strings.Split(other.Prerelease, ".")
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] == b[i] {
			continue
		}
		an, bn := numeric(a[i]), numeric(b[i])
		if an && !bn {
			return -1
		}
		if !an && bn {
			return 1
		}
		if an && bn && len(a[i]) != len(b[i]) {
			if len(a[i]) < len(b[i]) {
				return -1
			}
			return 1
		}
		return strings.Compare(a[i], b[i])
	}
	if len(a) < len(b) {
		return -1
	}
	return 1
}

func numeric(value string) bool {
	if value == "" {
		return false
	}
	for _, c := range value {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func bump(base Version, kind string) (Version, error) {
	next := Version{Major: base.Major, Minor: base.Minor, Patch: base.Patch}
	switch kind {
	case "major":
		next.Major++
		next.Minor, next.Patch = 0, 0
	case "minor":
		next.Minor++
		next.Patch = 0
	case "patch", "nightly":
		next.Patch++
	default:
		return Version{}, fmt.Errorf("unknown bump kind %q", kind)
	}
	return ParseVersion(next.String())
}
