package publish2

import (
	"fmt"

	"github.com/qor/admin"
)

type controller struct {
	Resource *admin.Resource
}

func (controller) Versions(context *admin.Context) {
	records := context.Resource.NewSlice()
	context.GetDB().Set(VersionMode, VersionMultipleMode).Set(ScheduleMode, ModeOff).Set(VisibleMode, ModeOff).Find(records, fmt.Sprintf("%v = ?", context.Resource.PrimaryDBName()), context.ResourceID)

	result := context.Render("publish2/versions", records)
	context.Writer.Write([]byte(result))
}
