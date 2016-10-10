package version

import "time"

type Schedule struct {
	ScheduleStartAt *time.Time
	ScheduleEndAt   *time.Time
}
