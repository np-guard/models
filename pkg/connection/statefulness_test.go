// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package connection_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/connection"
	"github.com/np-guard/models/pkg/netp"
)

func newTCPConn(t *testing.T, srcMinP, srcMaxP, dstMinP, dstMaxP int64) *connection.Set {
	t.Helper()
	return connection.TCPorUDPConnection(netp.ProtocolStringTCP, srcMinP, srcMaxP, dstMinP, dstMaxP)
}

func newUDPConn(t *testing.T, srcMinP, srcMaxP, dstMinP, dstMaxP int64) *connection.Set {
	t.Helper()
	return connection.TCPorUDPConnection(netp.ProtocolStringUDP, srcMinP, srcMaxP, dstMinP, dstMaxP)
}

func newICMPconn(t *testing.T) *connection.Set {
	t.Helper()
	return connection.ICMPConnection(
		connection.MinICMPtype, connection.MaxICMPtype,
		connection.MinICMPcode, connection.MaxICMPcode)
}

func newTCPUDPSet(t *testing.T, p netp.ProtocolString) *connection.Set {
	t.Helper()
	return connection.TCPorUDPConnection(p,
		connection.MinPort, connection.MaxPort,
		connection.MinPort, connection.MaxPort)
}

type statefulnessTest struct {
	name     string
	srcToDst *connection.Set
	dstToSrc *connection.Set
	// expectedIsStateful represents the expected IsStateful computed value for srcToDst,
	// which should be either StatefulTrue or StatefulFalse, given the input dstToSrc connection.
	// the computation applies only to the TCP protocol within those connections.
	expectedIsStateful connection.StatefulState
	// expectedStatefulConn represents the subset from srcToDst which is not related to the "non-stateful" mark (*) on the srcToDst connection,
	// the stateless part for TCP is srcToDst.Subtract(statefulConn)
	expectedStatefulConn *connection.Set
}

func (tt statefulnessTest) runTest(t *testing.T) {
	t.Helper()
	statefulConn := tt.srcToDst.ConnectionWithStatefulness(tt.dstToSrc)
	require.Equal(t, tt.expectedIsStateful, tt.srcToDst.IsStateful)
	require.True(t, tt.expectedStatefulConn.Equal(statefulConn))
}

func TestAll(t *testing.T) {
	var testCasesStatefulness = []statefulnessTest{
		{
			name:                 "tcp_all_ports_on_both_directions",
			srcToDst:             newTCPUDPSet(t, netp.ProtocolStringTCP), // TCP all ports
			dstToSrc:             newTCPUDPSet(t, netp.ProtocolStringTCP), // TCP all ports
			expectedIsStateful:   connection.StatefulTrue,
			expectedStatefulConn: newTCPUDPSet(t, netp.ProtocolStringTCP), // TCP all ports
		},
		{
			name:     "first_all_cons_second_tcp_with_ports",
			srcToDst: connection.All(),                                              // all connections
			dstToSrc: newTCPConn(t, 80, 80, connection.MinPort, connection.MaxPort), // TCP , src-ports: 80, dst-ports: all

			// there is a subset of the tcp connection which is not stateful
			expectedIsStateful: connection.StatefulFalse,

			// TCP src-ports: all, dst-port: 80 , union: all non-TCP conns
			expectedStatefulConn: connection.All().Subtract(newTCPUDPSet(t, netp.ProtocolStringTCP)).Union(
				newTCPConn(t, connection.MinPort, connection.MaxPort, 80, 80)),
		},
		{
			name:               "first_all_conns_second_no_tcp",
			srcToDst:           connection.All(), // all connections
			dstToSrc:           newICMPconn(t),   // ICMP
			expectedIsStateful: connection.StatefulFalse,
			// UDP, ICMP (all TCP is considered stateless here)
			expectedStatefulConn: connection.All().Subtract(newTCPUDPSet(t, netp.ProtocolStringTCP)),
		},
		{
			name:                 "tcp_with_ports_both_directions_exact_match",
			srcToDst:             newTCPConn(t, 80, 80, 443, 443),
			dstToSrc:             newTCPConn(t, 443, 443, 80, 80),
			expectedIsStateful:   connection.StatefulTrue,
			expectedStatefulConn: newTCPConn(t, 80, 80, 443, 443),
		},
		{
			name:                 "tcp_with_ports_both_directions_partial_match",
			srcToDst:             newTCPConn(t, 80, 100, 443, 443),
			dstToSrc:             newTCPConn(t, 443, 443, 80, 80),
			expectedIsStateful:   connection.StatefulFalse,
			expectedStatefulConn: newTCPConn(t, 80, 80, 443, 443),
		},
		{
			name:                 "tcp_with_ports_both_directions_no_match",
			srcToDst:             newTCPConn(t, 80, 100, 443, 443),
			dstToSrc:             newTCPConn(t, 80, 80, 80, 80),
			expectedIsStateful:   connection.StatefulFalse,
			expectedStatefulConn: connection.None(),
		},
		{
			name:                 "udp_and_tcp_with_ports_both_directions_no_match",
			srcToDst:             newTCPConn(t, 80, 100, 443, 443).Union(newUDPConn(t, 80, 100, 443, 443)),
			dstToSrc:             newTCPConn(t, 80, 80, 80, 80).Union(newUDPConn(t, 80, 80, 80, 80)),
			expectedIsStateful:   connection.StatefulFalse,
			expectedStatefulConn: newUDPConn(t, 80, 100, 443, 443),
		},
		{
			name:                 "no_tcp_in_first_direction",
			srcToDst:             newUDPConn(t, 70, 100, 443, 443),
			dstToSrc:             newTCPConn(t, 70, 80, 80, 80).Union(newUDPConn(t, 70, 80, 80, 80)),
			expectedIsStateful:   connection.StatefulTrue,
			expectedStatefulConn: newUDPConn(t, 70, 100, 443, 443),
		},
		{
			name:                 "empty_conn_in_first_direction",
			srcToDst:             connection.None(),
			dstToSrc:             newTCPConn(t, 80, 80, 80, 80).Union(newTCPUDPSet(t, netp.ProtocolStringUDP)),
			expectedIsStateful:   connection.StatefulTrue,
			expectedStatefulConn: connection.None(),
		},
		{
			name:     "only_udp_icmp_in_first_direction_and_empty_second_direction",
			srcToDst: newTCPUDPSet(t, netp.ProtocolStringUDP).Union(newICMPconn(t)),
			dstToSrc: connection.None(),
			// stateful analysis does not apply to udp/icmp, thus considered in the result as "stateful"
			// (to avoid marking it as stateless in the output)
			expectedIsStateful:   connection.StatefulTrue,
			expectedStatefulConn: newTCPUDPSet(t, netp.ProtocolStringUDP).Union(newICMPconn(t)),
		},
	}
	t.Parallel()
	// explainTests is the list of tests to run
	for testIdx := range testCasesStatefulness {
		tt := testCasesStatefulness[testIdx]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.runTest(t)
		})
	}
}
