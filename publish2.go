package publish2

import (
	"path"

	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
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

			res.Meta(&admin.Meta{
				Name:  "Versions",
				Label: "Versions",
				Type:  "publish_versions",
				Valuer: func(interface{}, *qor.Context) interface{} {
					return ""
				},
			})

			router := res.GetAdmin().GetRouter()
			ctr := controller{Resource: res}
			router.Get(path.Join(res.ToParam(), res.ParamIDName(), "versions"), ctr.Versions, admin.RouteConfig{Resource: res})

			res.IndexAttrs(res.IndexAttrs(), "Versions")
		}
	}
}
