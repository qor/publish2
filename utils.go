package publish2

import (
	"net/http"

	"github.com/qor/qor"
	"github.com/qor/qor/utils"
)

func GetScheduledTime(request *http.Request, writer http.ResponseWriter) string {
	var scheduledTime string
	if str := request.URL.Query().Get("publish_scheduled_time"); str != "" {
		scheduledTime = str
	}

	if cookie, err := request.Cookie("publish_scheduled_time"); err == nil {
		scheduledTime = cookie.Value
	}

	context := qor.Context{Request: request, Writer: writer}
	utils.SetCookie(http.Cookie{Name: "publish_scheduled_time", Value: scheduledTime}, &context)

	return scheduledTime
}
