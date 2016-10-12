package version_test

import "github.com/qor/qor/test/utils"

var DB = utils.TestDB()

func init() {
	models := []interface{}{}

	DB.DropTableIfExists(models...)
	DB.AutoMigrate(models...)
}
