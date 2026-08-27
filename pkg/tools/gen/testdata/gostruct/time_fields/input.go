package time_fields

import "time"

// Event is a time-typed record.
type Event struct {
	Started  time.Time     `kcl:"name=started,type=str"`
	Duration time.Duration `kcl:"name=duration,type=int"`
}
