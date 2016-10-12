package version_test

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/qor/version"
)

type Wiki struct {
	gorm.Model
	Title string
	Body  string
	version.Version
}

func TestVersions(t *testing.T) {
	var wiki = Wiki{Title: "wiki 1", Body: "wiki 1"}
	DB.Create(&wiki)

	wiki.SetVersionName("v1")
	wiki.Body = "wiki 1 - v1"
	DB.Save(&wiki)

	wiki.SetVersionName("v2")
	wiki.Body = "wiki 1 - v2"
	DB.Save(&wiki)

	var count int
	DB.Model(&Wiki{}).Where("id = ?", wiki.ID).Count(&count)
	if count != 3 {
		t.Errorf("Should have %v versions for wiki", 3)
	}
}
