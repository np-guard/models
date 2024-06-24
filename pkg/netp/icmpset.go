// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package netp

import (
	"github.com/np-guard/models/pkg/interval"
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
	case DestinationUnreachable:
		return newDestinationUnreachable + code
	case Redirect:
		return newRedirect + code
	case TimeExceeded:
		return newTimeExceeded + code
	case Echo:
		return newEcho
	case EchoReply:
		return newEchoReply
	case SourceQuench:
		return newSourceQuench
	default:
		return t
	}
}

//lint:ignore U1000 should be used in the future
func mapToOld(newCode int) (ICMP, error) {
	t := newCode
	switch {
	case newCode < newRedirect:
		t = newDestinationUnreachable
	case newCode < newTimeExceeded:
		t = newRedirect
	case newCode < ParameterProblem:
		t = newTimeExceeded
	case newCode == newEcho:
		t = Echo
	case newCode == newEchoReply:
		t = EchoReply
	case newCode == newSourceQuench:
		t = SourceQuench
	}
	code := newCode - t
	return NewICMP(&ICMPTypeCode{Type: t, Code: &code})
}

type ICMPSet uint32

func (s ICMPSet) IsSubset(other ICMPSet) bool {
	return s|other == other
}

func (s ICMPSet) Union(other ICMPSet) ICMPSet {
	return s | other
}

func (s ICMPSet) Contains(i int) bool {
	return ((1 << i) & s) != 0
}

func (s ICMPSet) IntervalSet() *interval.CanonicalSet {
	res := interval.NewCanonicalSet()
	for i := 0; i <= last; i++ {
		if s.Contains(i) {
			res.AddInterval(interval.New(int64(i), int64(i)))
		}
	}
	return res
}

func (s ICMPSet) collect(old int) []ICMP {
	res := []ICMP{}
	for code := 0; code <= maxCodes[old]; code++ {
		if s.Contains(mapToNew(old, code)) {
			res = append(res, ICMP{&ICMPTypeCode{Type: old, Code: &code}})
		}
	}
	if len(res) == maxCodes[old]+1 {
		res = []ICMP{{&ICMPTypeCode{Type: old, Code: nil}}}
	}
	return res
}

func (s ICMPSet) ICMPList() []ICMP {
	if ICMPSet(all).IsSubset(s) {
		return []ICMP{{nil}}
	}
	res := []ICMP{}
	for t := range maxCodes {
		res = append(res, s.collect(t)...)
	}
	return res
}

func fromIndex(i int) ICMPSet {
	return 1 << i
}

const (
	allDestinationUnreachable = 0b00000000000000111111
	allRedirect               = 0b00000000001111000000
	allTimeExceeded           = 0b00000000110000000000
	allOther                  = 0b11111111000000000000
	all                       = allDestinationUnreachable | allRedirect | allTimeExceeded | allOther
)

func FromICMP(t ICMP) ICMPSet {
	if t.typeCode == nil {
		return all
	}
	return fromIndex(mapToNew(t.typeCode.Type, *t.typeCode.Code))
}

func FromIntervalSet(intervalSet *interval.CanonicalSet) ICMPSet {
	if intervalSet.IsEmpty() {
		return 0
	}
	var res ICMPSet
	for i := 0; i <= last; i++ {
		if intervalSet.Contains(int64(i)) {
			res |= res.Union(fromIndex(i))
		}
	}
	return res
}
