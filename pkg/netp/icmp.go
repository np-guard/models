/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package netp

import (
	"fmt"
	"log"
)

type ICMPTypeCode struct {
	// ICMP type allowed.
	Type int

	// ICMP code allowed. If omitted, any code is allowed
	Code *int
}

type ICMP struct {
	typeCode *ICMPTypeCode
}

func NewICMP(typeCode *ICMPTypeCode) (ICMP, error) {
	err := ValidateICMP(typeCode)
	if err != nil {
		return ICMP{}, err
	}
	return ICMP{typeCode: typeCode}, nil
}

func (t ICMP) ICMPTypeCode() *ICMPTypeCode {
	if t.typeCode == nil {
		return nil
	}
	if t.typeCode.Code == nil {
		return t.typeCode
	}
	// avoid aliasing and mutation by someone else
	code := *t.typeCode.Code
	return &ICMPTypeCode{Type: t.typeCode.Type, Code: &code}
}

func (t ICMP) InverseDirection() Protocol {
	if t.typeCode == nil {
		return nil
	}

	if invType := inverseICMPType(t.typeCode.Type); invType != undefinedICMP {
		return ICMP{typeCode: &ICMPTypeCode{Type: invType, Code: t.typeCode.Code}}
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

func ValidateICMP(typeCode *ICMPTypeCode) error {
	if typeCode == nil {
		return nil
	}
	maxCode, ok := maxCodes[typeCode.Type]
	if !ok {
		return fmt.Errorf("invalid ICMP type %v", typeCode.Type)
	}
	if *typeCode.Code > maxCode {
		return fmt.Errorf("ICMP code %v is invalid for ICMP type %v", *typeCode.Code, typeCode.Type)
	}
	return nil
}

func (t ICMP) ProtocolString() ProtocolString {
	return ProtocolStringICMP
}
