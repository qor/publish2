package version

import "time"

type Schedule struct {
	VersionStartAt *time.Time
	VersionEndAt   *time.Time
}
