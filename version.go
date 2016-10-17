package version

type Version struct {
	VersionName string `gorm:"primary_key"`
}

func (version *Version) SetVersionName(name string) {
	version.VersionName = name
}

func (version *Version) GetVersionName() string {
	return version.VersionName
}

type VersionableInterface interface {
	SetVersionName(string)
	GetVersionName() string
}
