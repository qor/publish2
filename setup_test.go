package version_test

import (
	"github.com/qor/qor/test/utils"
	"github.com/qor/version"
)

var DB = utils.TestDB()

func init() {
	models := []interface{}{&Wiki{}, &Post{}, &Article{}, &Discount{}}

	DB.DropTableIfExists(models...)
	DB.AutoMigrate(models...)
	version.RegisterCallbacks(DB)
}
