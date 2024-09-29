/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package netp

import (
	"fmt"
	"log"
	"slices"
)

// general non-strict ICMP type, code ranges
const (
	MinICMPType int64 = 0
	MaxICMPType int64 = 254
	MinICMPCode int64 = 0
	MaxICMPCode int64 = 255
)

type ICMPTypeCode struct {
	// ICMP type allowed.
	Type int

	// ICMP code allowed. If omitted, any code is allowed
	Code *int
}

type ICMP struct {
	TypeCode *ICMPTypeCode
}

func NewICMP(typeCode *ICMPTypeCode) (ICMP, error) {
	err := ValidateICMP(typeCode)
	if err != nil {
		return ICMP{}, err
	}
	if typeCode == nil {
		return ICMP{TypeCode: nil}, nil
	}
	res := &ICMPTypeCode{Type: typeCode.Type}
	if HasSingleCode(typeCode.Type) {
		res.Code = nil
	} else {
		res.Code = typeCode.Code
	}
	return ICMP{TypeCode: res}, nil
}

func ICMPFromTypeAndCode(icmpType, icmpCode *int) (ICMP, error) {
	if icmpType == nil && icmpCode != nil {
		return ICMP{}, fmt.Errorf("cannot specify ICMP code without ICMP type")
	}
	if icmpType != nil {
		return NewICMP(&ICMPTypeCode{Type: *icmpType, Code: icmpCode})
	}
	return NewICMP(nil)
}

func int64ToInt(i *int64) *int {
	if i == nil {
		return nil
	}
	res := int(*i)
	return &res
}

func ICMPFromTypeAndCode64(icmpType, icmpCode *int64) (ICMP, error) {
	return ICMPFromTypeAndCode(int64ToInt(icmpType), int64ToInt(icmpCode))
}

func (t ICMP) ICMPTypeCode() *ICMPTypeCode {
	if t.TypeCode == nil {
		return nil
	}
	if t.TypeCode.Code == nil {
		return t.TypeCode
	}
	// avoid aliasing and mutation by someone else
	code := *t.TypeCode.Code
	return &ICMPTypeCode{Type: t.TypeCode.Type, Code: &code}
}

// InverseDirection returns the ICMP message that is the reply to this ICMP message.
func (t ICMP) InverseDirection() Protocol {
	if t.TypeCode == nil {
		return ICMP{TypeCode: nil}
	}

	if invType := inverseICMPType(t.TypeCode.Type); invType != undefinedICMP {
		// TODO: is this well defined?
		return ICMP{TypeCode: &ICMPTypeCode{Type: invType, Code: t.TypeCode.Code}}
	}
	return nil
}

// Based on https://datatracker.ietf.org/doc/html/rfc792

const (
	EchoReply              = 0
	DestinationUnreachable = 3
	SourceQuench           = 4
	Redirect               = 5
	Echo                   = 8
	TimeExceeded           = 11
	ParameterProblem       = 12
	Timestamp              = 13
	TimestampReply         = 14
	InformationRequest     = 15
	InformationReply       = 16

	undefinedICMP = -2
)

// inverseICMPType returns the reply type for request type and vice versa.
// When there is no inverse, returns undefinedICMP
func inverseICMPType(t int) int {
	switch t {
	case Echo:
		return EchoReply
	case EchoReply:
		return Echo

	case Timestamp:
		return TimestampReply
	case TimestampReply:
		return Timestamp

	case InformationRequest:
		return InformationReply
	case InformationReply:
		return InformationRequest

	case DestinationUnreachable, SourceQuench, Redirect, TimeExceeded, ParameterProblem:
		return undefinedICMP
	default:
		log.Panicf("Impossible ICMP type: %v", t)
	}
	return undefinedICMP
}

// maxCodes is a map from ICMP type to the maximum code allowed for that type.
// All the values between 0 and the maximum code are allowed.
var maxCodes = map[int]int{
	EchoReply:              0,
	DestinationUnreachable: 5,
	SourceQuench:           0,
	Redirect:               3,
	Echo:                   0,
	TimeExceeded:           1,
	ParameterProblem:       0,
	Timestamp:              0,
	TimestampReply:         0,
	InformationRequest:     0,
	InformationReply:       0,
}

func MaxCode(t int) int {
	return maxCodes[t]
}

var types = []int{
	EchoReply,
	DestinationUnreachable,
	SourceQuench,
	Redirect,
	Echo,
	TimeExceeded,
	ParameterProblem,
	Timestamp,
	TimestampReply,
	InformationRequest,
	InformationReply,
}

func Types() []int {
	return slices.Clone(types)
}

func ValidateICMP(typeCode *ICMPTypeCode) error {
	if typeCode == nil {
		return nil
	}
	maxCode, ok := maxCodes[typeCode.Type]
	if !ok {
		return fmt.Errorf("invalid ICMP type %v", typeCode.Type)
	}
	if typeCode.Code != nil && *typeCode.Code > maxCode {
		return fmt.Errorf("ICMP code %v is invalid for ICMP type %v", *typeCode.Code, typeCode.Type)
	}
	return nil
}

func HasSingleCode(t int) bool {
	return maxCodes[t] == 0
}

func (t ICMP) ProtocolString() ProtocolString {
	return ProtocolStringICMP
}
