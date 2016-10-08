package version_test

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/version"
)

type Blog struct {
	gorm.Model
	Title   string
	Content string
	version.SimpleMode
}

func TestSimpleMode(t *testing.T) {
	var blog = Blog{Title: "article 1", Content: "article 1"}
	DB.Create(&blog)

	blog.Content = "article 1 - v2"
	blog.SetVersion("v2")
	DB.Save(&blog)

	blog.Content = "article 1 - v3"
	blog.SetVersion("v3")
	DB.Save(&blog)

	var count int
	DB.Model(&Blog{}).Where("id = ?", blog.ID).Count(&count)
	if count != 3 {
		t.Errorf("Should have %v versions for blog", 3)
	}
}

func TestSimpleModeWithScheduleTime(t *testing.T) {
	var blog = Blog{Title: "article 2", Content: "article 2"}
	DB.Create(&blog)

	blog.Content = "article 2 - v2"
	blog.SetVersion("v2")
	startAt := time.Now().Add(-time.Hour)
	endAt := time.Now().Add(time.Hour)
	blog.VersionStartAt = &startAt
	blog.VersionEndAt = &endAt
	DB.Save(&blog)

	blog.Content = "article 2 - v3"
	blog.SetVersion("v3")
	blog.VersionStartAt = nil
	blog.VersionEndAt = nil
	DB.Save(&blog)

	var count int
	DB.Model(&Blog{}).Where("id = ?", blog.ID).Count(&count)
	if count != 3 {
		t.Errorf("Should have %v versions for blog", 3)
	}
}
