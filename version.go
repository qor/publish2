package version

import (
	"time"

	"github.com/jinzhu/gorm"
)

type QorVersion struct {
	gorm.Model
	Name    string
	StartAt *time.Time
	EndAt   *time.Time
}
