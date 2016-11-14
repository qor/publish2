package publish2

import (
	"fmt"

	"github.com/qor/admin"
)

type controller struct {
	Resource *admin.Resource
}

func (ctr controller) Publish(context *admin.Context) {
	context.Execute("publish2/dashboard", ctr.Resource)
}

func (ctr controller) Versions(context *admin.Context) {
	records := context.Resource.NewSlice()
	context.GetDB().Set(VersionMode, VersionMultipleMode).Set(ScheduleMode, ModeOff).Set(VisibleMode, ModeOff).Find(records, fmt.Sprintf("%v = ?", context.Resource.PrimaryDBName()), context.ResourceID)

	result := context.Render("publish2/versions", records)
	context.Writer.Write([]byte(result))
}
