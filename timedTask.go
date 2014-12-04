package timedTask

import (
	"sort"
	"time"
)

type Stamp struct {
	time time.Time
	do   func()
}

type byTime []*Stamp

func (s byTime) Len() int      { return len(s) }
func (s byTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byTime) Less(i, j int) bool {
	return s[i].time.Before(s[j].time)
}

type Sequence struct {
	list    []*Stamp
	element chan *Stamp
}

func NewSequence() *Sequence {
	return &Sequence{
		list:    nil,
		element: make(chan *Stamp),
	}
}

func (s *Sequence) Append(t time.Time, d func()) {
	s.element <- &Stamp{
		time: t,
		do:   d,
	}
}

func (s *Sequence) Run() {
	var (
		now       time.Time = time.Now()
		effective time.Time
	)
	for {
		sort.Sort(byTime(s.list))
		if len(s.list) == 0 {
			effective = now.AddDate(10, 0, 0)
		} else {
			effective = s.list[0].time
		}
		select {
		case now = <-time.After(effective.Sub(now)):
			for k, v := range s.list {
				if !effective.Equal(v.time) {
					break
				}
				v.do()
				s.list = s.list[k+1:]
			}
			continue

		case element := <-s.element:
			s.list = append(s.list, element)
		}
		now = time.Now()
	}
}
