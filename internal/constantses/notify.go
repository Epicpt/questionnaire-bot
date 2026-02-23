package constantses

import "time"

const (
	RemindTime = 3 * 24 * time.Hour
	NotRemind  = 0 // dont send notification
	Remind     = 1
)
