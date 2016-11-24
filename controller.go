package publish2

import (
	"fmt"
	"html/template"

	"github.com/qor/admin"
	"github.com/qor/qor"
)

type controller struct {
	Resource *admin.Resource
}

type visiblePublishResourceInterface interface {
	VisiblePublishResource(*qor.Context) bool
}

func (ctr controller) Dashboard(context *admin.Context) {
	type result struct {
		Resource *admin.Resource
		Results  interface{}
	}

	var results = []result{}

	for _, res := range context.Admin.GetResources() {
		if IsSchedulableModel(res.Value) {
			if visibleInterface, ok := res.Value.(visiblePublishResourceInterface); ok {
				if !visibleInterface.VisiblePublishResource(context.Context) {
					continue
				}
			} else if res.Config.Invisible {
				continue
			}

			db := context.GetDB()

			data := res.NewSlice()
			if db.Set(VersionMode, VersionMultipleMode).Find(data).RowsAffected > 0 {
				results = append(results, result{
					Resource: res,
					Results:  data,
				})
			}
		}
	}

	context.Action = "index"
	context.Execute("publish2/dashboard", results)
}

func (ctr controller) Versions(context *admin.Context) {
	records := context.Resource.NewSlice()
	context.GetDB().Set(VersionMode, VersionMultipleMode).Set(ScheduleMode, ModeOff).Set(VisibleMode, ModeOff).Find(records, fmt.Sprintf("%v = ?", context.Resource.PrimaryDBName()), context.ResourceID)

	result := context.Funcs(template.FuncMap{
		"version_metas": func() (metas []*admin.Meta) {
			for _, name := range []string{"VersionName", "ScheduledStartAt", "ScheduledEndAt", "PublishReady"} {
				if meta := ctr.Resource.GetMeta(name); meta != nil {
					metas = append(metas, meta)
				}
			}
			return
		},
	}).Render("publish2/versions", records)
	context.Writer.Write([]byte(result))
}
