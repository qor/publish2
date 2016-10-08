package version_test

import (
	"github.com/qor/qor/test/utils"
	"github.com/qor/version"
)

var DB = utils.TestDB()

func init() {
	models := []interface{}{&Blog{}, &Product{}, &version.QorVersion{}}

	DB.DropTableIfExists(models...)
	DB.AutoMigrate(models...)
}
