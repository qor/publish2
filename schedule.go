package version

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Schedule struct {
	ScheduledStartAt *time.Time
	ScheduledEndAt   *time.Time
}

func (schedule *Schedule) SetScheduledStartAt(t *time.Time) {
	schedule.ScheduledStartAt = t
}

func (schedule *Schedule) SetScheduledEndAt(t *time.Time) {
	schedule.ScheduledEndAt = t
}

type ScheduledInterface interface {
	SetScheduledStartAt(t *time.Time)
	SetScheduledEndAt(*time.Time)
}

type ScheduleEvent struct {
	gorm.Model
	Name            string
	ScheduleStartAt *time.Time
	ScheduleEndAt   *time.Time
}
