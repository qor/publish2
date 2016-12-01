package publish2

import (
	"net/http"

	"github.com/qor/qor"
	"github.com/qor/qor/utils"
)

func GetScheduledTime(request *http.Request, writer http.ResponseWriter) string {
	var scheduledTime string
	if values, ok := request.URL.Query()["publish_scheduled_time"]; ok {
		if len(values) > 0 && values[0] != "" {
			scheduledTime = values[0]
		}
	} else if cookie, err := request.Cookie("publish_scheduled_time"); err == nil {
		scheduledTime = cookie.Value
	}

	context := qor.Context{Request: request, Writer: writer}
	utils.SetCookie(http.Cookie{Name: "publish_scheduled_time", Value: scheduledTime}, &context)

	return scheduledTime
}
