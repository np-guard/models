/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package netset

/*
package netset implements types for network connection sets objects and operations.

Types defined in this package:
--------------------------------------------------------------------------------------

IPBlock - captures a set of IP ranges

TCPUDPSet - captures sets of protocols (within TCP,UDP only) and ports (source and destinaion)

ICMPSet - captures sets of types,codes for ICMP protocol

TransportSet - captures connection-sets for protocols from {TCP, UDP, ICMP}

ConnectionSet - captures a set of connections for tuples of (src IP range, dst IP range, TransportSet)

*/
