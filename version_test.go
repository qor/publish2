package version_test

import (
	"testing"
	"time"

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

type Post struct {
	gorm.Model
	Title string
	Body  string
	version.Version
	version.Schedule
}

func TestVersionsWithSchedule(t *testing.T) {
	var post = Post{Title: "post 1", Body: "post 1"}
	DB.Create(&post)

	post.SetVersionName("v1")
	post.Body = "post 1 - v1"
	now := time.Now()
	post.ScheduleStartAt = &now
	oneDayLater := now.Add(24 * time.Hour)
	post.ScheduleEndAt = &oneDayLater
	DB.Save(&post)

	post.SetVersionName("v2")
	post.Body = "post 1 - v2"
	post.ScheduleStartAt = &oneDayLater
	post.ScheduleEndAt = nil
	DB.Save(&post)

	var count int
	DB.Model(&Post{}).Where("id = ?", post.ID).Count(&count)
	if count != 3 {
		t.Errorf("Should have %v versions for post", 3)
	}
	// TODO test with scheduled time
}
