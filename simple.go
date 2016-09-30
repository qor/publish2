package version

import "time"

type SimpleVersionMode struct {
	VersionName    string
	VersionStartAt *time.Time
	VersionEndAt   *time.Time
}
