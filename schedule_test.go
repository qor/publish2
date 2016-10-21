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

	if DB.Set(version.ScheduleCurrent, oneDayAgo.Add(-time.Hour)).First(&Discount{}, "id = ?", discount.ID).RecordNotFound() {
		t.Errorf("Should find records that in scheduled with set schedule mode")
	}

	if DB.Set(version.ScheduleMode, version.ModeOff).First(&Discount{}, "id = ?", discount.ID).RecordNotFound() {
		t.Errorf("Should find records that not in scheduled with all mode")
	}

	discount.SetScheduledEndAt(&oneDayLater)
	DB.Save(&discount)

	if DB.First(&Discount{}, "id = ?", discount.ID).RecordNotFound() {
		t.Errorf("Should find records that in scheduled")
	}

	if DB.Set(version.ScheduleMode, version.ModeOff).First(&Discount{}, "id = ?", discount.ID).RecordNotFound() {
		t.Errorf("Should find records that in scheduled with all mode")
	}
}

func TestScheduleWithStartAndEnd(t *testing.T) {
	now := time.Now()
	oneDayAgo := time.Now().Add(-24 * time.Hour)
	oneDayLater := time.Now().Add(24 * time.Hour)

	discountV1 := Discount{Name: "discount2 - 1"}
	discountV1.SetScheduledStartAt(&oneDayAgo)
	discountV1.SetScheduledEndAt(&now)
	DB.Create(&discountV1)

	discountV2 := Discount{Name: "discount2 - 2"}
	oneHourLater := now.Add(time.Hour)
	discountV2.SetScheduledStartAt(&oneHourLater)
	discountV2.SetScheduledEndAt(&oneDayLater)
	DB.Create(&discountV2)

	var count uint
	DB.Set(version.ScheduleCurrent, now.Add(-time.Hour)).Model(&Discount{}).Where("id IN (?)", []uint{discountV1.ID, discountV2.ID}).Count(&count)
	if count != 1 {
		t.Errorf("Should find one discount with scheduled now, but got %v", count)
	}

	DB.Set(version.ScheduleStart, now.Add(-time.Hour)).Set(version.ScheduleEnd, oneDayLater).Model(&Discount{}).Where("id IN (?)", []uint{discountV1.ID, discountV2.ID}).Count(&count)
	if count != 2 {
		t.Errorf("Should find two discounts with scheduled time range, but got %v", count)
	}
}

type Campaign struct {
	gorm.Model
	Name string
	version.Schedule
	version.Visible
}

func TestScheduleWithVisible(t *testing.T) {
	oneDayAgo := time.Now().Add(-24 * time.Hour)
	oneDayLater := time.Now().Add(24 * time.Hour)

	campaign := Campaign{Name: "campaign 1"}
	campaign.SetScheduledStartAt(&oneDayAgo)
	campaign.SetScheduledEndAt(&oneDayLater)
	DB.Create(&campaign)

	if !DB.First(&Campaign{}, "id = ?", campaign.ID).RecordNotFound() {
		t.Errorf("Should not be able to find created record that not visible")
	}

	if DB.Set(version.VisibleMode, version.ModeOff).First(&Campaign{}, "id = ?", campaign.ID).RecordNotFound() {
		t.Errorf("Should be able to find created record that not visible when with visible mode on")
	}

	campaign.SetPublishReady(true)
	DB.Save(&campaign)

	if DB.First(&Campaign{}, "id = ?", campaign.ID).RecordNotFound() {
		t.Errorf("Should be able to find created record that visible")
	}

	if !DB.Set(version.ScheduleCurrent, oneDayLater.Add(time.Hour)).First(&Campaign{}, "id = ?", campaign.ID).RecordNotFound() {
		t.Errorf("Should not be able to find created record that not in schedule")
	}
}
