package publish2_test

import (
	"github.com/qor/publish2"
	"github.com/qor/qor/test/utils"
)

var DB = utils.TestDB()

func init() {
	models := []interface{}{&Wiki{}, &Post{}, &Article{}, &Discount{}, &User{}, &Campaign{}, &Product{}}

	DB.DropTableIfExists(models...)
	DB.AutoMigrate(models...)
	publish2.RegisterCallbacks(DB)
}
