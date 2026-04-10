package qdrant

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"unicode"
)

const fullVersionParts = 3
const reducedVersionParts = 2
const unknownVersion = "Unknown"

type Version struct {
	Major int
	Minor int
}

func removeLeadingNonNumeric(versionStr string) string {
	return strings.TrimLeftFunc(versionStr, func(r rune) bool {
		return !unicode.IsDigit(r)
	})
}

// ParseVersion converts a version string "x.y[.z]" into a Version struct.
func ParseVersion(versionStr string) (*Version, error) {
	cleanedVersionStr := removeLeadingNonNumeric(versionStr)
	parts := strings.SplitN(cleanedVersionStr, ".", fullVersionParts)
	if len(parts) < reducedVersionParts {
		return nil, fmt.Errorf("unable to parse version, expected format: x.y[.z], found: %s", cleanedVersionStr)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse major version: %w", err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse minor version: %w", err)
	}

	return &Version{
		Major: major,
		Minor: minor,
	}, nil
}

func IsCompatible(clientVersion, serverVersion string) bool {
	if clientVersion == serverVersion {
		return true
	}

	client, err := ParseVersion(clientVersion)
	if err != nil {
		return false
	}

	server, err := ParseVersion(serverVersion)
	if err != nil {
		return false
	}

	if client.Major != server.Major {
		return false
	}

	diff := client.Minor - server.Minor
	return diff <= 1 && diff >= -1
}

func getClientVersion() string {
	packageName := "github.com/hanzoai/vector-go"
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return unknownVersion
	}

	if bi.Main.Path == packageName {
		return bi.Main.Version
	}

	for _, dep := range bi.Deps {
		if dep.Path == packageName {
			return dep.Version
		}
	}

	return unknownVersion
}
