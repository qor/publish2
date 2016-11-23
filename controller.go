package publish2

import (
	"fmt"
	"html/template"

	"github.com/qor/admin"
)

type controller struct {
	Resource *admin.Resource
}

func (ctr controller) Dashboard(context *admin.Context) {
	context.Action = "index"
	context.Execute("publish2/dashboard", ctr.Resource)
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
