package version_test

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/qor/version"
)

type User struct {
	gorm.Model
	Name string
	version.Visible
}

func TestPublishReady(t *testing.T) {
	user := User{Name: "user"}
	DB.Create(&user)

	if !DB.First(&User{}, "id = ?", user.ID).RecordNotFound() {
		t.Errorf("Should not able to find created record when publish not ready")
	}

	if DB.Set(version.VisibleMode, version.ModeOff).First(&User{}, "id = ?", user.ID).RecordNotFound() {
		t.Errorf("Should be able to find created record with visible mode `all`")
	}

	user.SetPublishReady(true)
	DB.Save(&user)

	if DB.First(&User{}, "id = ?", user.ID).RecordNotFound() {
		t.Errorf("Should be able to find created record when publish is ready")
	}
}
