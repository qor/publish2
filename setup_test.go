package version_test

import "github.com/qor/qor/test/utils"

var DB = utils.TestDB()

func init() {
	models := []interface{}{&Wiki{}, &Post{}}

	DB.DropTableIfExists(models...)
	DB.AutoMigrate(models...)
}
