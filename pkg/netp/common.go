/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package netp

type Protocol interface {
	// InverseDirection returns the response expected for a request made using this protocol
	InverseDirection() Protocol
}

type AnyProtocol struct{}

func (t AnyProtocol) InverseDirection() Protocol { return AnyProtocol{} }

type ProtocolString string

const (
	ProtocolStringTCP  ProtocolString = "TCP"
	ProtocolStringUDP  ProtocolString = "UDP"
	ProtocolStringICMP ProtocolString = "ICMP"
)
