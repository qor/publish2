package version_test

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/version"
)

type Discount struct {
	gorm.Model
	Name string
	version.Schedule
}

func TestSchedule(t *testing.T) {
	oneDayAgo := time.Now().Add(-24 * time.Hour)
	oneDayLater := time.Now().Add(24 * time.Hour)

	discount := Discount{Name: "discount1"}
	DB.Create(&discount)

	discount.SetScheduledEndAt(&oneDayAgo)
	DB.Save(&discount)

	if !DB.First(&Discount{}, "id = ?", discount.ID).RecordNotFound() {
		t.Errorf("Should not find records that not in scheduled")
	}

	discount.SetScheduledEndAt(&oneDayLater)
	DB.Save(&discount)

	if DB.First(&Discount{}, "id = ?", discount.ID).RecordNotFound() {
		t.Errorf("Should find records that in scheduled")
	}
}
