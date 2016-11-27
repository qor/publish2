package publish2

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/qor/utils"
)

const (
	ModeOff             = "off"
	VersionMode         = "publish:version:mode"
	VersionNameMode     = "publish:version:name"
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
		isSchedulable      = IsSchedulableModel(scope.Value)
		isVersionable      = IsVersionableModel(scope.Value)
		isPublishReadyable = IsPublishReadyableModel(scope.Value)
		conditions         []string
		conditionValues    []interface{}
	)

	if isSchedulable {
		var scheduledStartTime, scheduledEndTime, scheduledCurrentTime *time.Time
		var mode, _ = scope.DB().Get(ScheduleMode)

		if v, ok := scope.Get(ScheduleStart); ok {
			if t, ok := v.(*time.Time); ok {
				scheduledStartTime = t
			} else if t, ok := v.(time.Time); ok {
				scheduledStartTime = &t
			}

			if scheduledStartTime != nil {
				conditions = append(conditions, "(scheduled_end_at IS NULL OR scheduled_end_at >= ?)")
				conditionValues = append(conditionValues, scheduledStartTime)
			}
		}

		if v, ok := scope.Get(ScheduleEnd); ok {
			if t, ok := v.(*time.Time); ok {
				scheduledEndTime = t
			} else if t, ok := v.(time.Time); ok {
				scheduledEndTime = &t
			}

			if scheduledEndTime != nil {
				conditions = append(conditions, "(scheduled_start_at IS NULL OR scheduled_start_at <= ?)")
				conditionValues = append(conditionValues, scheduledEndTime)
			}
		}

		if len(conditions) == 0 && mode != ModeOff {
			if v, ok := scope.Get(ScheduleCurrent); ok {
				if t, ok := v.(*time.Time); ok {
					scheduledCurrentTime = t
				} else if t, ok := v.(time.Time); ok {
					scheduledCurrentTime = &t
				}
			}

			if scheduledCurrentTime == nil {
				now := time.Now()
				scheduledCurrentTime = &now
			}

			conditions = append(conditions, "(scheduled_start_at IS NULL OR scheduled_start_at <= ?) AND (scheduled_end_at IS NULL OR scheduled_end_at >= ?)")
			conditionValues = append(conditionValues, scheduledCurrentTime, scheduledCurrentTime)
		}
	}

	if isPublishReadyable {
		switch mode, _ := scope.DB().Get(VisibleMode); mode {
		case ModeOff:
		default:
			conditions = append(conditions, "publish_ready = ?")
			conditionValues = append(conditionValues, true)
		}
	}

	if isVersionable {
		switch mode, _ := scope.DB().Get(VersionMode); mode {
		case VersionMultipleMode:
			scope.Search.Where(strings.Join(conditions, " AND "), conditionValues...)
		default:
			if versionName, ok := scope.DB().Get(VersionNameMode); ok && versionName != "" {
				scope.Search.Where("version_name = ?", versionName)
			} else {
				var sql string
				var primaryKeys []string

				for _, primaryField := range scope.PrimaryFields() {
					if primaryField.DBName != "version_name" {
						primaryKeys = append(primaryKeys, scope.Quote(primaryField.DBName))
					}
				}

				primaryKeyCondition := strings.Join(primaryKeys, ",")
				if len(conditions) == 0 {
					sql = fmt.Sprintf("(%v, version_priority) IN (SELECT %v, MAX(version_priority) FROM %v GROUP BY %v)", primaryKeyCondition, primaryKeyCondition, scope.QuotedTableName(), primaryKeyCondition)
				} else {
					sql = fmt.Sprintf("(%v, version_priority) IN (SELECT %v, MAX(version_priority) FROM %v WHERE %v GROUP BY %v)", primaryKeyCondition, primaryKeyCondition, scope.QuotedTableName(), strings.Join(conditions, " AND "), primaryKeyCondition)
				}

				scope.Search.Where(sql, conditionValues...)
			}
		}

		scope.Search.Order("version_priority DESC")
	} else {
		scope.Search.Where(strings.Join(conditions, " AND "), conditionValues...)
	}
}

func createCallback(scope *gorm.Scope) {
	if IsVersionableModel(scope.Value) {
		if field, ok := scope.FieldByName("VersionName"); ok {
			if field.IsBlank {
				field.Set(DefaultVersionName)
			}
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
