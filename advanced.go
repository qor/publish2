package version

import "fmt"

type AdvancedMode struct {
	VersionID string `gorm:"primary_key"`
}

func (mode *AdvancedMode) SetVersion(versionID uint) {
	mode.VersionID = fmt.Sprint(versionID)
}
