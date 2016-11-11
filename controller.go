package publish2

import "github.com/qor/admin"

type controller struct {
	Resource *admin.Resource
}

func (controller) Versions(context *admin.Context) {
	result := context.Render("publish2/versions", nil)
	context.Writer.Write([]byte(result))
}
