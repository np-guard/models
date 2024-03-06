package connectionset

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTCPConn(t *testing.T, srcMinP, srcMaxP, dstMinP, dstMaxP int64) *ConnectionSet {
	t.Helper()
	res := NewConnectionSet(false)
	res.AddTCPorUDPConn(ProtocolStringTCP, srcMinP, srcMaxP, dstMinP, dstMaxP)
	return res
}

func newUDPConn(t *testing.T, srcMinP, srcMaxP, dstMinP, dstMaxP int64) *ConnectionSet {
	t.Helper()
	res := NewConnectionSet(false)
	res.AddTCPorUDPConn(ProtocolStringUDP, srcMinP, srcMaxP, dstMinP, dstMaxP)
	return res
}

func newICMPconn(t *testing.T) *ConnectionSet {
	t.Helper()
	res := NewConnectionSet(false)
	res.AddICMPConnection(MinICMPtype, MaxICMPtype, MinICMPcode, MaxICMPcode)
	return res
}

func allButTCP(t *testing.T) *ConnectionSet {
	t.Helper()
	res := NewConnectionSet(true)
	tcpOnly := res.tcpConn()
	return res.Subtract(tcpOnly)
}

type statefulnessTest struct {
	name     string
	srcToDst *ConnectionSet
	dstToSrc *ConnectionSet
	// expectedIsStateful represents the expected IsStateful computed value for srcToDst,
	// which should be either StatefulTrue or StatefulFalse, given the input dstToSrc connection.
	// the computation applies only to the TCP protocol within those connections.
	expectedIsStateful int
	// expectedStatefulConn represents the subset from srcToDst which is not related to the "non-stateful" mark (*) on the srcToDst connection,
	// the stateless part for TCP is srcToDst.Subtract(statefulConn)
	expectedStatefulConn *ConnectionSet
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
			srcToDst:             newTCPConn(t, MinPort, MaxPort, MinPort, MaxPort), // TCP all ports
			dstToSrc:             newTCPConn(t, MinPort, MaxPort, MinPort, MaxPort), // TCP all ports
			expectedIsStateful:   StatefulTrue,
			expectedStatefulConn: newTCPConn(t, MinPort, MaxPort, MinPort, MaxPort), // TCP all ports
		},
		{
			name:     "first_all_cons_second_tcp_with_ports",
			srcToDst: NewConnectionSet(true),                  // all connections
			dstToSrc: newTCPConn(t, 80, 80, MinPort, MaxPort), // TCP , src-ports: 80, dst-ports: all

			// there is a subset of the tcp connection which is not stateful
			expectedIsStateful: StatefulFalse,

			// TCP src-ports: all, dst-port: 80 , union: all non-TCP conns
			expectedStatefulConn: allButTCP(t).Union(newTCPConn(t, MinPort, MaxPort, 80, 80)),
		},
		{
			name:                 "first_all_conns_second_no_tcp",
			srcToDst:             NewConnectionSet(true), // all connections
			dstToSrc:             newICMPconn(t),         // ICMP
			expectedIsStateful:   StatefulFalse,
			expectedStatefulConn: allButTCP(t), // UDP, ICMP (all TCP is considered stateless here)
		},
		{
			name:                 "tcp_with_ports_both_directions_exact_match",
			srcToDst:             newTCPConn(t, 80, 80, 443, 443),
			dstToSrc:             newTCPConn(t, 443, 443, 80, 80),
			expectedIsStateful:   StatefulTrue,
			expectedStatefulConn: newTCPConn(t, 80, 80, 443, 443),
		},
		{
			name:                 "tcp_with_ports_both_directions_partial_match",
			srcToDst:             newTCPConn(t, 80, 100, 443, 443),
			dstToSrc:             newTCPConn(t, 443, 443, 80, 80),
			expectedIsStateful:   StatefulFalse,
			expectedStatefulConn: newTCPConn(t, 80, 80, 443, 443),
		},
		{
			name:                 "tcp_with_ports_both_directions_no_match",
			srcToDst:             newTCPConn(t, 80, 100, 443, 443),
			dstToSrc:             newTCPConn(t, 80, 80, 80, 80),
			expectedIsStateful:   StatefulFalse,
			expectedStatefulConn: NewConnectionSet(false),
		},
		{
			name:                 "udp_and_tcp_with_ports_both_directions_no_match",
			srcToDst:             newTCPConn(t, 80, 100, 443, 443).Union(newUDPConn(t, 80, 100, 443, 443)),
			dstToSrc:             newTCPConn(t, 80, 80, 80, 80).Union(newUDPConn(t, 80, 80, 80, 80)),
			expectedIsStateful:   StatefulFalse,
			expectedStatefulConn: newUDPConn(t, 80, 100, 443, 443),
		},
		{
			name:                 "no_tcp_in_first_direction",
			srcToDst:             newUDPConn(t, 80, 100, 443, 443),
			dstToSrc:             newTCPConn(t, 80, 80, 80, 80).Union(newUDPConn(t, 80, 80, 80, 80)),
			expectedIsStateful:   StatefulTrue,
			expectedStatefulConn: newUDPConn(t, 80, 100, 443, 443),
		},
		{
			name:                 "empty_conn_in_first_direction",
			srcToDst:             NewConnectionSet(false),
			dstToSrc:             newTCPConn(t, 80, 80, 80, 80).Union(newUDPConn(t, MinPort, MaxPort, MinPort, MaxPort)),
			expectedIsStateful:   StatefulTrue,
			expectedStatefulConn: NewConnectionSet(false),
		},
		{
			name:     "only_udp_icmp_in_first_direction_and_empty_second_direction",
			srcToDst: newUDPConn(t, MinPort, MaxPort, MinPort, MaxPort).Union(newICMPconn(t)),
			dstToSrc: NewConnectionSet(false),
			// stateful analysis does not apply to udp/icmp, thus considered in the result as "stateful"
			// (to avoid marking it as stateless in the output)
			expectedIsStateful:   StatefulTrue,
			expectedStatefulConn: newUDPConn(t, MinPort, MaxPort, MinPort, MaxPort).Union(newICMPconn(t)),
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
	fmt.Println("done")
}
