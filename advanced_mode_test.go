package version_test

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/qor/version"
)

type Product struct {
	gorm.Model
	Title   string
	Content string
	version.AdvancedMode
}

func TestAdvancedMode(t *testing.T) {
	var product = Product{Title: "article 1", Content: "article 1"}
	DB.Create(&product)
}
