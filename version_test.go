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
	now := time.Now()
	oneDayLater := now.Add(24 * time.Hour)

	post := Post{Title: "post 1", Body: "post 1"}
	DB.Create(&post)

	post.SetVersionName("v1")
	post.Body = "post 1 - v1"
	post.SetScheduledStartAt(&now)
	post.SetScheduledEndAt(&oneDayLater)
	DB.Save(&post)

	post.SetVersionName("v2")
	post.Body = "post 1 - v2"
	post.SetScheduledStartAt(&oneDayLater)
	post.SetScheduledEndAt(nil)
	DB.Save(&post)

	var count int
	DB.Model(&Post{}).Where("id = ?", post.ID).Count(&count)
	if count != 1 {
		t.Errorf("Should have one available post", 1)
	}

	var post1, post2, post3 Post
	DB.Set("publish:scheduled_time", now.Add(-24*time.Hour)).Model(&Post{}).Where("id = ?", post.ID).First(&post1)
	if post1.Body != "post 1" {
		t.Errorf("should find default version")
	}

	DB.Set("publish:scheduled_time", now.Add(6*time.Hour)).Model(&Post{}).Where("id = ?", post.ID).First(&post2)
	if post3.Body != "post 1 - v1" {
		t.Errorf("should find first version")
	}

	DB.Set("publish:scheduled_time", now.Add(25*time.Hour)).Model(&Post{}).Where("id = ?", post.ID).First(&post3)
	if post3.Body != "post 1 - v2" {
		t.Errorf("should find second version")
	}
}
