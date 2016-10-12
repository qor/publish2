package version

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Schedule struct {
	ScheduleStartAt *time.Time
	ScheduleEndAt   *time.Time
}

type ScheduleEvent struct {
	gorm.Model
	Name            string
	ScheduleStartAt *time.Time
	ScheduleEndAt   *time.Time
}
