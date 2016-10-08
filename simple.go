package version

import "time"

type SimpleMode struct {
	VersionName    string `gorm:"primary_key"`
	VersionStartAt *time.Time
	VersionEndAt   *time.Time
}

func (mode *SimpleMode) SetVersion(name string) {
	mode.VersionName = name
}
