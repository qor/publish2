package version_test

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/qor/version"
)

type Product struct {
	gorm.Model
	Title       string
	Description string
	version.AdvancedMode
}

var globalVersion1 = version.QorVersion{Name: "global_v1"}
var globalVersion2 = version.QorVersion{Name: "global_v2"}

func init() {
	DB.Save(&globalVersion1)
	DB.Save(&globalVersion2)
}

func TestAdvancedMode(t *testing.T) {
	var product = Product{Title: "product 1", Description: "product 1"}
	DB.Create(&product)

	product.SetVersion(globalVersion1.ID)
	product.Description = "product 1 - v1"
	DB.Save(&product)

	product.SetVersion(globalVersion2.ID)
	product.Description = "product 1 - v2"
	DB.Save(&product)

	var count int
	DB.Model(&Product{}).Where("id = ?", product.ID).Count(&count)
	if count != 3 {
		t.Errorf("Should have %v versions for product", 3)
	}
}
