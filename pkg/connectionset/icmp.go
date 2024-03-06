package connectionset

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
	echoReply              = 0
	destinationUnreachable = 3
	sourceQuench           = 4
	redirect               = 5
	echo                   = 8
	timeExceeded           = 11
	parameterProblem       = 12
	timestamp              = 13
	timestampReply         = 14
	informationRequest     = 15
	informationReply       = 16

	undefinedICMP = -2
)

// inverseICMPType returns the reply type for request type and vice versa.
// When there is no inverse, returns undefinedICMP
func inverseICMPType(t int) int {
	switch t {
	case echo:
		return echoReply
	case echoReply:
		return echo

	case timestamp:
		return timestampReply
	case timestampReply:
		return timestamp

	case informationRequest:
		return informationReply
	case informationReply:
		return informationRequest

	case destinationUnreachable, sourceQuench, redirect, timeExceeded, parameterProblem:
		return undefinedICMP
	default:
		log.Panicf("Impossible ICMP type: %v", t)
	}
	return undefinedICMP
}

//nolint:revive // magic numbers are fine here
func ValidateICMP(t, c int) error {
	maxCodes := map[int]int{
		echoReply:              0,
		destinationUnreachable: 5,
		sourceQuench:           0,
		redirect:               3,
		echo:                   0,
		timeExceeded:           1,
		parameterProblem:       0,
		timestamp:              0,
		timestampReply:         0,
		informationRequest:     0,
		informationReply:       0,
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
