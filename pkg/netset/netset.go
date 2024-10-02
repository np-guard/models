/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// Package netset implements types for network connection sets objects and operations.
// Types defined in this package:
// IPBlock - captures a set of IP ranges
// TCPUDPSet - captures sets of protocols (within TCP,UDP only) and ports (source and destination)
// ICMPSet - captures sets of type,code values for ICMP protocol
// TransportSet - captures union of elements from TCPUDPSet, ICMPSet
// EndpointsTrafficSet - captures a set of traffic attribute for tuples of (source IP range, destination IP range, TransportSet)
package netset
