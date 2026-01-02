package clock

import "time"

type Fixed struct {
	T time.Time
}

func (f Fixed) Now() time.Time { return f.T }
