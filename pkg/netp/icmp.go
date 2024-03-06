package netp

import (
	"fmt"
	"log"
)

type ICMPCodeType struct {
	// ICMP type allowed.
	Type int

	// ICMP code allowed. If omitted, any code is allowed
	Code *int
}

type ICMP struct {
	*ICMPCodeType
}

func (t ICMP) InverseDirection() Protocol {
	if t.ICMPCodeType == nil {
		return nil
	}

	if invType := inverseICMPType(t.Type); invType != undefinedICMP {
		return ICMP{ICMPCodeType: &ICMPCodeType{Type: invType, Code: t.Code}}
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

//nolint:revive // magic numbers are fine here
func ValidateICMP(t, c int) error {
	maxCodes := map[int]int{
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
	maxCode, ok := maxCodes[t]
	if !ok {
		return fmt.Errorf("invalid ICMP type %v", t)
	}
	if c > maxCode {
		return fmt.Errorf("ICMP code %v is invalid for ICMP type %v", c, t)
	}
	return nil
}
