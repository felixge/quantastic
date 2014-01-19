package version

import (
	"fmt"
)

// @TODO Add time, os, go version, build number, etc.
func NewVersion(release, commit string) Version {
	return Version{release, commit}
}

type Version struct {
	release string
	commit  string
}

func (v Version) String() string {
	if v.release == "" {
		return "N/A"
	} else if v.commit == "" {
		return v.release
	}
	return fmt.Sprintf("%s (%s)", v.release, v.commit)
}
