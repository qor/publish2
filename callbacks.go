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

	db.Callback().Create().Before("gorm:begin_transaction").Register("publish:versions", createCallback)
	db.Callback().Update().Before("gorm:begin_transaction").Register("publish:versions", updateCallback)
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
		sql := fmt.Sprintf("(id, version_priority) IN (SELECT id, MAX(version_priority) FROM %v WHERE (scheduled_start_at IS NULL OR scheduled_start_at <= ?) AND (scheduled_end_at IS NULL OR scheduled_end_at >= ?) GROUP BY id)", scope.QuotedTableName())
		scope.Search.Where(sql, scheduledTime, scheduledTime).Order("version_priority DESC")
	}
}

func createCallback(scope *gorm.Scope) {
	if field, ok := scope.FieldByName("VersionName"); ok {
		field.IsBlank = false
	}

	if field, ok := scope.FieldByName("VersionPriority"); ok {
		var scheduledTime *time.Time
		if scheduled, ok := scope.Value.(ScheduledInterface); ok {
			scheduledTime = scheduled.GetScheduledStartAt()
		}
		if scheduledTime == nil {
			unix := time.Unix(0, 0)
			scheduledTime = &unix
		}

		priority := fmt.Sprintf("%v_%v", scheduledTime.UTC().Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339Nano))
		field.Set(priority)
	}
}

func updateCallback(scope *gorm.Scope) {
	if field, ok := scope.FieldByName("VersionPriority"); ok {
		var scheduledTime *time.Time
		if scheduled, ok := scope.Value.(ScheduledInterface); ok {
			scheduledTime = scheduled.GetScheduledStartAt()
		}
		if scheduledTime == nil {
			unix := time.Unix(0, 0)
			scheduledTime = &unix
		}

		priority := fmt.Sprintf("%v_%v", scheduledTime.UTC().Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339Nano))
		field.Set(priority)
	}
}
