package cache

import (
	"fmt"
	"strings"
)

// Key represents a key in the format executor_version_protocol.
type Key struct {
	Value string

	Major  int
	Minor  int
	Patch  int
	Suffix string

	Version  string
	Executor string
	Protocol string
}

// ByVersion is a custom type for sorting keys by version.
type ByVersion []Key

func (a ByVersion) Len() int      { return len(a) }
func (a ByVersion) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByVersion) Less(i, j int) bool {
	return a[i].IsLessThan(a[j])
}

func (key Key) IsLessThan(other Key) bool {
	if key.Executor != other.Executor {
		return key.Executor < other.Executor
	}

	if key.Protocol != other.Protocol {
		return key.Protocol < other.Protocol
	}

	if key.Major != other.Major {
		return key.Major < other.Major
	}

	if key.Minor != other.Minor {
		return key.Minor < other.Minor
	}

	if key.Patch != other.Patch {
		return key.Patch < other.Patch
	}

	return key.Suffix < other.Suffix
}

func (key Key) IsGreaterThan(other Key) bool {
	if key.Executor != other.Executor {
		return key.Executor > other.Executor
	}

	if key.Protocol != other.Protocol {
		return key.Protocol > other.Protocol
	}

	if key.Major != other.Major {
		return key.Major > other.Major
	}

	if key.Minor != other.Minor {
		return key.Minor > other.Minor
	}

	if key.Patch != other.Patch {
		return key.Patch > other.Patch
	}

	return key.Suffix > other.Suffix
}

// ParseKey parses a key into its components.
func ParseKey(key string) (Key, error) {
	parts := strings.Split(key, "_")
	if len(parts) != 3 {
		return Key{}, fmt.Errorf("invalid key format: %s", key)
	}

	versionParts := strings.Split(parts[1], ".")
	if len(versionParts) != 3 {
		return Key{}, fmt.Errorf("invalid version format: %s", parts[1])
	}

	var major, minor, patch int
	var suffix string
	if _, err := fmt.Sscanf(versionParts[0], "%d", &major); err != nil {
		return Key{}, err
	}
	if _, err := fmt.Sscanf(versionParts[1], "%d", &minor); err != nil {
		return Key{}, err
	}
	if _, err := fmt.Sscanf(versionParts[2], "%d-%s", &patch, &suffix); err != nil {
		if _, err := fmt.Sscanf(versionParts[2], "%d", &patch); err != nil {
			return Key{}, err
		}
	}

	return Key{
		Value: key,

		Major:  major,
		Minor:  minor,
		Patch:  patch,
		Suffix: parts[1],

		Executor: parts[0],
		Version:  parts[1],
		Protocol: parts[2],
	}, nil
}
