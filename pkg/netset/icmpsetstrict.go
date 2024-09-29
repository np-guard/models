/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package netset

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/np-guard/models/pkg/netp"
)

// ICMPSetStrict is a set of _valid_ (by RFC) ICMP values, encoded as a bitset
type ICMPSetStrict uint32

// Encoding for ICMP types and codes, enumerating the possible pairs of values.
// For example:
// * 0 is the pair (type=DestinationUnreachable, code=0).
// * 2 is the pair (type=DestinationUnreachable, code=2).
// * 7 is the pair (type=Redirect, code=1).
// The idea is to use a simple bitset for the set of _valid_ ICMP values.
const (
	encodedDestinationUnreachable = 0
	encodedRedirect               = 6
	encodedTimeExceeded           = 10
	encodedParameterProblem       = 12
	encodedTimestamp              = 13
	encodedTimestampReply         = 14
	encodedInformationRequest     = 15
	encodedInformationReply       = 16
	encodedEcho                   = 17
	encodedEchoReply              = 18
	encodedSourceQuench           = 19
	last                          = 19
)

func encode(t, code int) int {
	switch t {
	case netp.DestinationUnreachable:
		return encodedDestinationUnreachable + code
	case netp.Redirect:
		return encodedRedirect + code
	case netp.TimeExceeded:
		return encodedTimeExceeded + code
	case netp.ParameterProblem:
		return encodedParameterProblem
	case netp.Timestamp:
		return encodedTimestamp
	case netp.TimestampReply:
		return encodedTimestampReply
	case netp.InformationRequest:
		return encodedInformationRequest
	case netp.InformationReply:
		return encodedInformationReply
	case netp.Echo:
		return encodedEcho
	case netp.EchoReply:
		return encodedEchoReply
	case netp.SourceQuench:
		return encodedSourceQuench
	default:
		log.Panicf("Invalid ICMP type %v", t)
		return t
	}
}

//lint:ignore U1000 should be used in the future
func decode(encodedCode int) (netp.ICMP, error) {
	t := encodedCode
	switch {
	case encodedCode < encodedRedirect:
		t = encodedDestinationUnreachable
	case encodedCode < encodedTimeExceeded:
		t = encodedRedirect
	case encodedCode < netp.ParameterProblem:
		t = encodedTimeExceeded
	case encodedCode == encodedEcho:
		t = netp.Echo
	case encodedCode == encodedEchoReply:
		t = netp.EchoReply
	case encodedCode == encodedSourceQuench:
		t = netp.SourceQuench
	}
	code := encodedCode - t
	return netp.NewICMP(&netp.ICMPTypeCode{Type: t, Code: &code})
}

func (s *ICMPSetStrict) IsSubset(other *ICMPSetStrict) bool {
	return ((*s) | (*other)) == (*other)
}

func (s *ICMPSetStrict) Union(other *ICMPSetStrict) *ICMPSetStrict {
	var res = (*s) | (*other)
	return &res
}

func (s *ICMPSetStrict) Intersect(other *ICMPSetStrict) *ICMPSetStrict {
	var res = (*s) & (*other)
	return &res
}

func (s *ICMPSetStrict) Subtract(other *ICMPSetStrict) *ICMPSetStrict {
	var res = (*s) & ^(*other)
	return &res
}

func (s *ICMPSetStrict) Equal(other *ICMPSetStrict) bool {
	return *s == *other
}

func (s *ICMPSetStrict) Copy() *ICMPSetStrict {
	var res = *s
	return &res
}

func (s *ICMPSetStrict) Hash() int {
	return int(*s)
}

func (s *ICMPSetStrict) Size() int {
	res := 0
	for i := 0; i <= last; i++ {
		if s.Contains(i) {
			res++
		}
	}
	return res
}

func (s *ICMPSetStrict) IsEmpty() bool {
	return s.Equal(EmptyICMPSetStrict())
}

func (s *ICMPSetStrict) Contains(i int) bool {
	return ((1 << i) & (*s)) != 0
}

// collect returns a list of ICMP values for a given type, collecting into a single ICMP value with nil Code if all codes are present.
func (s *ICMPSetStrict) collect(old int) []netp.ICMP {
	var res []netp.ICMP
	for code := 0; code <= netp.MaxCode(old); code++ {
		if s.Contains(encode(old, code)) {
			icmp, err := netp.NewICMP(&netp.ICMPTypeCode{Type: old, Code: &code})
			if err != nil {
				log.Panicf("collection failed for type %v, code %v", old, &code)
			}
			res = append(res, icmp)
		}
	}
	if len(res) == netp.MaxCode(old)+1 {
		res = []netp.ICMP{{TypeCode: &netp.ICMPTypeCode{Type: old, Code: nil}}}
	}
	return res
}

// Partitions returns a list of ICMP values.
// if all codes for a given type are present, it adds a single ICMP value with nil Code.
// If all ICMP values are present, a single ICMP value with nil TypeCode is returned.
func (s *ICMPSetStrict) Partitions() []netp.ICMP {
	all := ICMPSetStrict(allCodes)
	if all.IsSubset(s) {
		return []netp.ICMP{{TypeCode: nil}}
	}
	var res []netp.ICMP
	for _, t := range netp.Types() {
		res = append(res, s.collect(t)...)
	}
	return res
}

func fromIndex(i int) *ICMPSetStrict {
	var res ICMPSetStrict = 1 << i
	return &res
}

func (s *ICMPSetStrict) IsAll() bool {
	return s.Equal(AllICMPSetStrict())
}

// constants for sets of ICMP codes, grouped by types.
// For example, allDestinationUnreachable is the set of all ICMP codes for DestinationUnreachable type.
const (
	allDestinationUnreachable = 0b00000000000000111111
	allRedirect               = 0b00000000001111000000
	allTimeExceeded           = 0b00000000110000000000
	allOther                  = 0b11111111000000000000
	allCodes                  = allDestinationUnreachable | allRedirect | allTimeExceeded | allOther
)

func EmptyICMPSetStrict() *ICMPSetStrict {
	var res ICMPSetStrict = 0
	return &res
}

func AllICMPSetStrict() *ICMPSetStrict {
	res := ICMPSetStrict(allCodes)
	return &res
}

func NewICMPSetStrict(t netp.ICMP) *ICMPSetStrict {
	if t.TypeCode == nil {
		return AllICMPSetStrict()
	}
	if t.TypeCode.Code != nil {
		return fromIndex(encode(t.TypeCode.Type, *t.TypeCode.Code))
	}
	var res ICMPSetStrict
	switch t.TypeCode.Type {
	case netp.DestinationUnreachable:
		res = allDestinationUnreachable
	case netp.Redirect:
		res = allRedirect
	case netp.TimeExceeded:
		res = allTimeExceeded
	default:
		res = *fromIndex(encode(t.TypeCode.Type, 0))
	}
	return &res
}

func getICMPCubeStr(cube netp.ICMP) string {
	tc := cube.ICMPTypeCode()
	if tc == nil {
		return ""
	}
	if tc.Code == nil {
		if netp.HasSingleCode(tc.Type) {
			return fmt.Sprintf("icmp-type: %d icmp-code: 0", tc.Type)
		}
		return fmt.Sprintf("icmp-type: %d", tc.Type)
	}
	return fmt.Sprintf("icmp-type: %d icmp-code: %d", tc.Type, *tc.Code)
}

func (s *ICMPSetStrict) String() string {
	if s.IsEmpty() {
		return ""
	}
	cubes := s.Partitions()
	var resStrings = make([]string, len(cubes))
	for i, cube := range cubes {
		resStrings[i] = getICMPCubeStr(cube)
	}
	sort.Strings(resStrings)
	str := string(netp.ProtocolStringICMP)
	last := strings.Join(resStrings, semicolon)
	if last != "" {
		str += " " + last
	}
	return str
}
