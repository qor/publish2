package version_test

import (
	"testing"

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
}
