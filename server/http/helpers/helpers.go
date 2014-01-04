package helpers

import (
	"fmt"
	"time"
)

var Helpers = map[string]interface{} {
	"duration": func(val time.Duration) string {
		minutes := int(val.Minutes())
		hours := minutes / 60
		minutes = minutes % 60
		return fmt.Sprintf("%02d:%02d", hours, minutes)
	},
	"time": func(val time.Time) string {
		return val.Format("15:04")
	},
}
