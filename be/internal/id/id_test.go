package id

import "fmt"

type Seq struct {
	N int
}

func (s *Seq) NewID() string {
	s.N++
	return fmt.Sprintf("id-%d", s.N)
}
