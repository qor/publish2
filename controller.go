package publish2

import "github.com/qor/admin"

type controller struct {
	Resource *admin.Resource
}

func (controller) Versions(context *admin.Context) {
	if record, err := context.FindOne(); err == nil {
		result := context.Render("publish2/versions", record)
		context.Writer.Write([]byte(result))
	}
}
