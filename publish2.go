package publish2

import (
	"path"

	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/utils"
)

func init() {
	admin.RegisterViewPath("github.com/qor/publish2/views")
}

func (SharedVersion) ConfigureQorResource(res resource.Resourcer) {
	enablePublishMode(res)
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

				if res.GetAdmin().GetResource(utils.ModelType(&Publish{}).Name()) == nil {
					res.GetAdmin().AddResource(&Publish{})
				}
			}

			if IsVersionableModel(res.Value) {
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

				res.Action(&admin.Action{
					Name:        "Create New Version",
					Method:      "GET",
					URLOpenType: "slideout",
					URL: func(record interface{}, context *admin.Context) string {
						if versionable, ok := record.(VersionableInterface); ok {
							url := context.URLFor(record) + "?new_version=true&"
							if versionName := versionable.GetVersionName(); versionName != "" {
								url = url + "&version_name=" + versionName
							}
							return url
						}
						return ""
					},
					Modes: []string{"menu_item"},
				})

				router := res.GetAdmin().GetRouter()
				ctr := controller{Resource: res}
				router.Get(path.Join(res.ToParam(), res.ParamIDName(), "versions"), ctr.Versions, admin.RouteConfig{Resource: res})

				res.IndexAttrs(res.IndexAttrs(), "Versions", "-VersionPriority")
				res.EditAttrs(res.EditAttrs(), "-Versions", "-VersionPriority", "VersionName")
				res.NewAttrs(res.NewAttrs(), "-Versions", "-VersionPriority", "VersionName")
			}

			if IsShareableVersionModel(res.Value) {
				res.Meta(&admin.Meta{
					Name: "ShareableVersion",
					Type: "string",
					Valuer: func(interface{}, *qor.Context) interface{} {
						return ""
					},
					Setter: func(resource interface{}, metaValue *resource.MetaValue, context *qor.Context) {
					},
				})
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

			res.GetAdmin().RegisterFuncMap("get_schedule_events", func(context *admin.Context) interface{} {
				res := context.Admin.GetResource("ScheduleEvent")
				scheduleEvents := res.NewSlice()
				context.GetDB().Find(scheduleEvents)
				return map[string]interface{}{"Resource": res, "Events": scheduleEvents}
			})
		}
	}
}

type Publish struct {
}

func (Publish) ConfigureQorResourceBeforeInitialize(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		res.Meta(&admin.Meta{
			Name: "ScheduledStartAt",
			Valuer: func(interface{}, *qor.Context) interface{} {
				return ""
			},
		})

		res.UseTheme("publish2")

		if res.Config.Name == "" {
			res.Name = "Schedule"
		}

		if len(res.Config.Menu) == 0 {
			res.Config.Menu = []string{"Publishing"}
		}

		Admin := res.GetAdmin()

		if res.GetAdmin().GetResource("ScheduleEvent") == nil {
			Admin.AddResource(&ScheduleEvent{}, &admin.Config{Menu: res.Config.Menu, Priority: -1})
			Admin.Config.DB.AutoMigrate(&ScheduleEvent{})
		}

		Admin.GetRouter().Use(&admin.Middleware{
			Name: "publish2",
			Handler: func(context *admin.Context, middleware *admin.Middleware) {
				tx := context.GetDB()

				if startAt := context.Request.URL.Query().Get("schedule_start_at"); startAt != "" {
					if t, err := utils.ParseTime(startAt, context.Context); err == nil {
						tx = tx.Set(ScheduleStart, t).Set(VersionMode, VersionMultipleMode)
					}
				}

				if endAt := context.Request.URL.Query().Get("schedule_end_at"); endAt != "" {
					if t, err := utils.ParseTime(endAt, context.Context); err == nil {
						tx = tx.Set(ScheduleEnd, t).Set(VersionMode, VersionMultipleMode)
					}
				}

				if versionName := context.Request.URL.Query().Get("version_name"); versionName != "" {
					tx = tx.Set(VersionNameMode, versionName)
				}

				context.SetDB(tx)
				middleware.Next(context)
			},
		})

		ctr := controller{Resource: res}
		Admin.GetRouter().Get(res.ToParam(), ctr.Dashboard, admin.RouteConfig{Resource: res})
	}
}
