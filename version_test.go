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
		t.Errorf("Should have one available post")
	}

	var post1, post2, post3 Post
	DB.Set("publish:scheduled_time", now.Add(-24*time.Hour)).Model(&Post{}).Where("id = ?", post.ID).First(&post1)
	if post1.Body != "post 1" {
		t.Errorf("should find default version, but got %v", post1.Body)
	}

	DB.Set("publish:scheduled_time", now.Add(6*time.Hour)).Model(&Post{}).Where("id = ?", post.ID).First(&post2)
	if post2.Body != "post 1 - v1" {
		t.Errorf("should find first version, but got %v", post2.Body)
	}

	DB.Set("publish:scheduled_time", now.Add(25*time.Hour)).Model(&Post{}).Where("id = ?", post.ID).First(&post3)
	if post3.Body != "post 1 - v2" {
		t.Errorf("should find second version, but got %v", post3.Body)
	}
}

func TestVersionsWithOverlappedSchedule(t *testing.T) {
	now := time.Now()
	post := prepareOverlappedPost("post 2")

	var post1, post2, post3 Post
	DB.Set("publish:scheduled_time", now.Add(-36*time.Hour)).Model(&Post{}).Where("id = ?", post.ID).First(&post1)
	if post1.Body != "post 2" {
		t.Errorf("should find default version, but got %v", post1.Body)
	}

	DB.Set("publish:scheduled_time", now.Add(6*time.Hour)).Model(&Post{}).Where("id = ?", post.ID).First(&post2)
	if post2.Body != "post 2 - v2" {
		t.Errorf("should find first version, but got %v", post2.Body)
	}

	DB.Set("publish:scheduled_time", now.Add(25*time.Hour)).Model(&Post{}).Where("id = ?", post.ID).First(&post3)
	if post3.Body != "post 2 - v1" {
		t.Errorf("should find second version, but got %v", post3.Body)
	}
}

func TestVersionsWithOverlappedSchedules(t *testing.T) {
	now := time.Now()
	postV1 := prepareOverlappedPost("post 3 - 1")
	postV2 := prepareOverlappedPost("post 3 - 2")
	postIDs := []uint{postV1.ID, postV2.ID}

	var count uint
	DB.Set("publish:scheduled_time", now.Add(-36*time.Hour)).Model(&Post{}).Where("id IN (?)", postIDs).Count(&count)
	if count != 2 {
		t.Errorf("should only find 2 valid versions, but got %v", count)
	}

	DB.Set("publish:scheduled_time", now.Add(6*time.Hour)).Model(&Post{}).Where("id IN (?)", postIDs).Count(&count)
	if count != 2 {
		t.Errorf("should only find 2 valid versions, but got %v", count)
	}

	DB.Set("publish:scheduled_time", now.Add(25*time.Hour)).Model(&Post{}).Where("id IN (?)", postIDs).Count(&count)
	if count != 2 {
		t.Errorf("should only find 2 valid versions, but got %v", count)
	}
}

func prepareOverlappedPost(name string) *Post {
	now := time.Now()
	oneDayAgo := now.Add(-24 * time.Hour)
	oneDayLater := now.Add(24 * time.Hour)

	post := Post{Title: name, Body: name}
	DB.Create(&post)

	post.SetVersionName("v1")
	post.Body = name + " - v1"
	post.SetScheduledStartAt(&oneDayAgo)
	post.SetScheduledEndAt(nil)
	DB.Save(&post)

	post.SetVersionName("v2")
	post.Body = name + " - v2"
	post.SetScheduledStartAt(&now)
	post.SetScheduledEndAt(&oneDayLater)
	DB.Save(&post)

	return &post
}
