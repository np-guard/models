// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package netp

import (
	"fmt"
	"log"
	"math"

	"github.com/np-guard/models/pkg/interval"
)

const (
	newDestinationUnreachable = 0
	newRedirect               = 6
	newTimeExceeded           = 10
	newEcho                   = 17
	newEchoReply              = 18
	newSourceQuench           = 19
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
func mapToOld(newCode int) (t, code int) {
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
	default:
		t = newCode
	}
	code = newCode - t
	return
}

type ICMPSet uint32

func (s ICMPSet) IsSubset(other ICMPSet) bool {
	return s|other == other
}

func (s ICMPSet) Union(other ICMPSet) ICMPSet {
	return s | other
}

const (
	allDestinationUnreachable = 0b00000000000000111111
	allRedirect               = 0b00000000001111000000
	allTimeExceeded           = 0b00000000110000000000
	allOther                  = 0b11111111000000000000
)

func FromICMP(t ICMP) ICMPSet {
	if t.typeCode == nil {
		return allDestinationUnreachable | allRedirect | allTimeExceeded | allOther
	}
	if t.typeCode == nil {
		return math.MaxUint32
	}
	return 1 << mapToNew(t.typeCode.Type, *t.typeCode.Code)
}
