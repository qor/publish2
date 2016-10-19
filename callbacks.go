package version

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/qor/utils"
)

const (
	VersionMode         = "publish:version:mode"
	VersionMultipleMode = "multiple"

	ScheduleMode    = "publish:schedule:mode"
	ScheduleCurrent = "publish:schedule:current"
	ScheduleStart   = "publish:schedule:start"
	ScheduleEnd     = "publish:schedule:end"

	VisibleMode = "publish:visible:mode"
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

func IsPublishReadyableModel(model interface{}) (ok bool) {
	if model != nil {
		_, ok = reflect.New(utils.ModelType(model)).Interface().(PublishReadyInterface)
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
	var (
		scheduledTime      *time.Time
		isSchedulable      = IsSchedulableModel(scope.Value)
		isVersionable      = IsVersionableModel(scope.Value)
		isPublishReadyable = IsPublishReadyableModel(scope.Value)
		conditions         []string
		conditionValues    []interface{}
	)

	if isSchedulable {
		switch mode, _ := scope.DB().Get(ScheduleMode); mode {
		case "all":
		default:
			if v, ok := scope.Get(ScheduleCurrent); ok {
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

			conditions = append(conditions, "(scheduled_start_at IS NULL OR scheduled_start_at <= ?) AND (scheduled_end_at IS NULL OR scheduled_end_at >= ?)")
			conditionValues = append(conditionValues, scheduledTime, scheduledTime)
		}
	}

	if isPublishReadyable {
		switch mode, _ := scope.DB().Get(VisibleMode); mode {
		case "all":
		default:
			conditions = append(conditions, "publish_ready = ?")
			conditionValues = append(conditionValues, true)
		}
	}

	if isVersionable {
		switch mode, _ := scope.DB().Get(VersionMode); mode {
		case VersionMultipleMode:
			scope.Search.Where(strings.Join(conditions, " AND "), conditionValues...).Order("version_priority DESC")
		default:
			var sql string
			if len(conditions) == 0 {
				sql = fmt.Sprintf("(id, version_priority) IN (SELECT id, MAX(version_priority) FROM %v GROUP BY id)", scope.QuotedTableName())
			} else {
				sql = fmt.Sprintf("(id, version_priority) IN (SELECT id, MAX(version_priority) FROM %v WHERE %v GROUP BY id)", scope.QuotedTableName(), strings.Join(conditions, " AND "))
			}

			scope.Search.Where(sql, conditionValues...).Order("version_priority DESC")
		}
	} else {
		scope.Search.Where(strings.Join(conditions, " AND "), conditionValues...)
	}
}

func createCallback(scope *gorm.Scope) {
	if IsVersionableModel(scope.Value) {
		if field, ok := scope.FieldByName("VersionName"); ok {
			field.IsBlank = false
		}

		updateVersionPriority(scope)
	}
}

func updateCallback(scope *gorm.Scope) {
	if IsVersionableModel(scope.Value) {
		updateVersionPriority(scope)
	}
}

func updateVersionPriority(scope *gorm.Scope) {
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
