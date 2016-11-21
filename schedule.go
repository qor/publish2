package publish2

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Schedule struct {
	ScheduledStartAt *time.Time `gorm:"index"`
	ScheduledEndAt   *time.Time `gorm:"index"`
	ScheduleEventID  *uint
}

func (schedule *Schedule) GetScheduledStartAt() *time.Time {
	return schedule.ScheduledStartAt
}

func (schedule *Schedule) SetScheduledStartAt(t *time.Time) {
	schedule.ScheduledStartAt = t
}

func (schedule *Schedule) GetScheduledEndAt() *time.Time {
	return schedule.ScheduledEndAt
}

func (schedule *Schedule) SetScheduledEndAt(t *time.Time) {
	schedule.ScheduledEndAt = t
}

func (schedule *Schedule) GetScheduleEventID() *uint {
	return schedule.ScheduleEventID
}

type ScheduledInterface interface {
	GetScheduledStartAt() *time.Time
	SetScheduledStartAt(*time.Time)
	GetScheduledEndAt() *time.Time
	SetScheduledEndAt(*time.Time)
	GetScheduleEventID() *uint
}

type ScheduleEvent struct {
	gorm.Model
	Name            string
	ScheduleStartAt *time.Time
	ScheduleEndAt   *time.Time
}

func (ScheduleEvent) AfterUpdate(tx *gorm.DB) {
	// sync time changes
}
