package publish2

import (
	"path"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/utils"
)

func init() {
	admin.RegisterViewPath("github.com/qor/publish2/views")
}

func (Version) ConfigureQorResource(res resource.Resourcer) {
	enablePublishMode(res)
}

func (Schedule) ConfigureQorResource(res resource.Resourcer) {
	enablePublishMode(res)
}

func (Visible) ConfigureQorResource(res resource.Resourcer) {
	enablePublishMode(res)
}

func enablePublishMode(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		if res.GetTheme("publish2") == nil {
			res.UseTheme("publish2")

			if IsPublishReadyableModel(res.Value) {
				res.IndexAttrs(res.IndexAttrs(), "-PublishReady")
				res.EditAttrs(res.EditAttrs(), "PublishReady")
				res.NewAttrs(res.NewAttrs(), "PublishReady")
			}

			if IsSchedulableModel(res.Value) {
				res.Meta(&admin.Meta{
					Name: "ScheduleEventID",
					Type: "hidden",
				})

				res.IndexAttrs(res.IndexAttrs(), "-ScheduleEventID")
				res.EditAttrs(res.EditAttrs(), "ScheduledStartAt", "ScheduledEndAt", "ScheduleEventID")
				res.NewAttrs(res.NewAttrs(), "ScheduledStartAt", "ScheduledEndAt", "ScheduleEventID")

				res.Scope(&admin.Scope{
					Default: true,
					Handle: func(tx *gorm.DB, context *qor.Context) *gorm.DB {
						if startAt := context.Request.URL.Query().Get("schedule_start_at"); startAt != "" {
							if t, err := utils.ParseTime(startAt, context); err == nil {
								tx = tx.Set(ScheduleStart, t)
							}
						}

						if endAt := context.Request.URL.Query().Get("schedule_end_at"); endAt != "" {
							if t, err := utils.ParseTime(endAt, context); err == nil {
								tx = tx.Set(ScheduleEnd, t)
							}
						}

						return tx
					},
				})
			}

			if IsVersionableModel(res.Value) {
				if IsSchedulableModel(res.Value) {
					res.IndexAttrs(res.IndexAttrs(), "-ScheduledStartAt", "-ScheduledEndAt")
				}

				res.Meta(&admin.Meta{
					Name: "VersionPriority",
					Type: "hidden",
				})

				res.Meta(&admin.Meta{
					Name: "VersionName",
					Type: "hidden",
				})

				res.Meta(&admin.Meta{
					Name:  "Versions",
					Label: "Versions",
					Type:  "publish_versions",
					Valuer: func(interface{}, *qor.Context) interface{} {
						return ""
					},
				})

				res.Scope(&admin.Scope{
					Default: true,
					Handle: func(tx *gorm.DB, context *qor.Context) *gorm.DB {
						if versionName := context.Request.URL.Query().Get("version_name"); versionName != "" {
							tx = tx.Set(VersionNameMode, versionName)
						}
						return tx
					},
				})

				router := res.GetAdmin().GetRouter()
				ctr := controller{Resource: res}
				router.Get(path.Join(res.ToParam(), res.ParamIDName(), "versions"), ctr.Versions, admin.RouteConfig{Resource: res})

				res.IndexAttrs(res.IndexAttrs(), "Versions", "-VersionName", "-VersionPriority")
				res.EditAttrs(res.EditAttrs(), "-Versions", "-VersionPriority", "VersionName")
				res.NewAttrs(res.NewAttrs(), "-Versions", "-VersionPriority", "VersionName")
			}

			res.GetAdmin().RegisterFuncMap("get_schedule_event", func(record interface{}, context *admin.Context) interface{} {
				if scheduledInterface, ok := record.(ScheduledInterface); ok {
					var scheduleEvent ScheduleEvent
					if scheduledInterface.GetScheduleEventID() != nil {
						context.GetDB().First(&scheduleEvent, "id = ?", scheduledInterface.GetScheduleEventID())
						return scheduleEvent
					}
				}
				return nil
			})
		}
	}
}

type Publish struct {
}

func (Publish) ConfigureQorResourceBeforeInitialize(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		if res.Config.Name == "" {
			res.Name = "Schedule"
		}

		if len(res.Config.Menu) == 0 {
			res.Config.Menu = []string{"Publishing"}
		}

		Admin := res.GetAdmin()
		//TODO make it configable
		Admin.AddResource(&ScheduleEvent{}, &admin.Config{Name: "Event", Menu: res.Config.Menu, Priority: -1})
		Admin.Config.DB.AutoMigrate(&ScheduleEvent{})

		ctr := controller{Resource: res}
		Admin.GetRouter().Get(res.ToParam(), ctr.Dashboard)
	}
}
