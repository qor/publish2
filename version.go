package version

type Version struct {
	VersionName string `gorm:"primary_key"`
}

func (version *Version) SetVersionName(name string) {
	version.VersionName = name
}
