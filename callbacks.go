package version

import (
	"fmt"
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

func IsVersionableModel(model interface{}) (ok bool) {
	if model != nil {
		_, ok = reflect.New(utils.ModelType(model)).Interface().(VersionableInterface)
	}
	return
}

func RegisterCallbacks(db *gorm.DB) {
	db.Callback().Query().Before("gorm:query").Register("publish:query", queryCallback)
	db.Callback().RowQuery().Before("gorm:query").Register("publish:query", queryCallback)
}

func queryCallback(scope *gorm.Scope) {
	var scheduledTime *time.Time
	var isSchedulable, isVersionable bool

	if IsSchedulableModel(scope.Value) {
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

		isSchedulable = true
	}

	if IsVersionableModel(scope.Value) {
		isVersionable = true
	}

	switch {
	case isSchedulable && isVersionable:
		sql := fmt.Sprintf("(version_name = '' AND id NOT IN (SELECT id FROM %v WHERE version_name <> '' AND (scheduled_start_at IS NULL OR scheduled_start_at <= ?) AND (scheduled_end_at IS NULL OR scheduled_end_at >= ?))) OR (version_name <> '' AND (scheduled_start_at IS NULL OR scheduled_start_at <= ?) AND (scheduled_end_at IS NULL OR scheduled_end_at >= ?))", scope.QuotedTableName())
		scope.Search.Where(sql, scheduledTime, scheduledTime, scheduledTime, scheduledTime).Order("scheduled_start_at DESC")
	}
}
