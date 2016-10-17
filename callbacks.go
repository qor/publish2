package version

import (
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/qor/utils"
)

func IsSchedulableModel(model interface{}) (ok bool) {
	if model != nil {
		_, ok = reflect.New(utils.ModelType(model)).Interface().(ScheduledInterface)
	}
	return
}

func RegisterCallbacks(db *gorm.DB) {
	db.Callback().Query().Before("gorm:query").Register("publish:query", queryCallback)
	db.Callback().RowQuery().Before("gorm:query").Register("publish:query", queryCallback)
}

func queryCallback(scope *gorm.Scope) {
	if IsSchedulableModel(scope.Value) {
		var scheduledTime *time.Time
		if v, ok := scope.Get("publish:scheduled_time"); ok {
			if t, ok := v.(*time.Time); ok {
				scheduledTime = t
			} else if t, ok := v.(time.Time); ok {
				scheduledTime = &t
			}
		}

		if scheduledTime == nil {
			now := time.Now()
			scheduledTime = &now
		}

		scope.Search.Where("(scheduled_start_at IS NULL OR scheduled_start_at < ?) AND (scheduled_end_at IS NULL OR scheduled_end_at < ?)", scheduledTime, scheduledTime)
	}
}
