// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package fconn

import (
	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
)

const (
	newDestinationUnreachable = 0
	newRedirect               = 6
	newTimeExceeded           = 10
	newEcho                   = 17
	newEchoReply              = 18
	newSourceQuench           = 19
	last                      = 19
)

func mapToNew(t, code int) int {
	switch t {
	case netp.DestinationUnreachable:
		return newDestinationUnreachable + code
	case netp.Redirect:
		return newRedirect + code
	case netp.TimeExceeded:
		return newTimeExceeded + code
	case netp.Echo:
		return newEcho
	case netp.EchoReply:
		return newEchoReply
	case netp.SourceQuench:
		return newSourceQuench
	default:
		return t
	}
}

//lint:ignore U1000 should be used in the future
func mapToOld(newCode int) (netp.ICMP, error) {
	t := newCode
	switch {
	case newCode < newRedirect:
		t = newDestinationUnreachable
	case newCode < newTimeExceeded:
		t = newRedirect
	case newCode < netp.ParameterProblem:
		t = newTimeExceeded
	case newCode == newEcho:
		t = netp.Echo
	case newCode == newEchoReply:
		t = netp.EchoReply
	case newCode == newSourceQuench:
		t = netp.SourceQuench
	}
	code := newCode - t
	return netp.NewICMP(&netp.ICMPTypeCode{Type: t, Code: &code})
}

type ICMPSet uint32

func (s *ICMPSet) ContainedIn(other *ICMPSet) bool {
	return (*s)|(*other) == (*other)
}

func (s *ICMPSet) Union(other *ICMPSet) *ICMPSet {
	var res = (*s) | (*other)
	return &res
}

func (s *ICMPSet) Intersect(other *ICMPSet) *ICMPSet {
	var res = (*s) & (*other)
	return &res
}

func (s *ICMPSet) Subtract(other *ICMPSet) *ICMPSet {
	var res = (*s) & ^(*other)
	return &res
}

func (s *ICMPSet) Equal(other *ICMPSet) bool {
	return *s == *other
}

func (s *ICMPSet) Copy() *ICMPSet {
	var res = *s
	return &res
}

func (s *ICMPSet) Hash() int {
	return int(*s)
}

func (s *ICMPSet) Size() int {
	res := 0
	for i := 0; i <= last; i++ {
		if s.Contains(i) {
			res++
		}
	}
	return res
}

func (s *ICMPSet) IsEmpty() bool {
	return s.Equal(EmptyICMPSet())
}

func (s *ICMPSet) Contains(i int) bool {
	return ((1 << i) & (*s)) != 0
}

func (s *ICMPSet) IntervalSet() *interval.CanonicalSet {
	res := interval.NewCanonicalSet()
	for i := 0; i <= last; i++ {
		if s.Contains(i) {
			res.AddInterval(interval.New(int64(i), int64(i)))
		}
	}
	return res
}

func (s *ICMPSet) collect(old int) []netp.ICMP {
	res := []netp.ICMP{}
	for code := 0; code <= netp.MaxCodes[old]; code++ {
		if s.Contains(mapToNew(old, code)) {
			res = append(res, netp.ICMP{TypeCode: &netp.ICMPTypeCode{Type: old, Code: &code}})
		}
	}
	if len(res) == netp.MaxCodes[old]+1 {
		res = []netp.ICMP{{TypeCode: &netp.ICMPTypeCode{Type: old, Code: nil}}}
	}
	return res
}

func (s *ICMPSet) ICMPList() []netp.ICMP {
	all := ICMPSet(allCodes)
	if s.ContainedIn(&all) {
		return []netp.ICMP{{TypeCode: nil}}
	}
	res := []netp.ICMP{}
	for t := range netp.MaxCodes {
		res = append(res, s.collect(t)...)
	}
	return res
}

func fromIndex(i int) *ICMPSet {
	var res ICMPSet = 1 << i
	return &res
}

const (
	allDestinationUnreachable = 0b00000000000000111111
	allRedirect               = 0b00000000001111000000
	allTimeExceeded           = 0b00000000110000000000
	allOther                  = 0b11111111000000000000
	allCodes                  = allDestinationUnreachable | allRedirect | allTimeExceeded | allOther
)

func EmptyICMPSet() *ICMPSet {
	var res ICMPSet = 0
	return &res
}

func NewICMPSet(t netp.ICMP) *ICMPSet {
	if t.TypeCode == nil {
		return EmptyICMPSet()
	}
	return fromIndex(mapToNew(t.TypeCode.Type, *t.TypeCode.Code))
}

func FromIntervalSet(intervalSet *interval.CanonicalSet) *ICMPSet {
	if intervalSet.IsEmpty() {
		return EmptyICMPSet()
	}
	var res = EmptyICMPSet()
	for i := 0; i <= last; i++ {
		if intervalSet.Contains(int64(i)) {
			res = res.Union(fromIndex(i))
		}
	}
	return res
}
