package testutil

import (
	"fmt"
	"time"

	"personal-budgeting/be/internal/clock"
	"personal-budgeting/be/internal/id"
)

type FixedClock struct {
	T time.Time
}

func (f FixedClock) Now() time.Time { return f.T }

var _ clock.Clock = FixedClock{}

type SeqID struct {
	N int
}

func (s *SeqID) NewID() string {
	s.N++
	return fmt.Sprintf("id-%d", s.N)
}

var _ id.Generator = (*SeqID)(nil)


