package netset_test

import (
	"testing"

	"github.com/np-guard/models/pkg/connection"
	"github.com/np-guard/models/pkg/netset"
	"github.com/stretchr/testify/require"
)

func TestConnectionSet(t *testing.T) {
	src1, _ := netset.IPBlockFromCidr("10.240.10.0/24")
	dst1, _ := netset.IPBlockFromCidr("10.240.10.0/32")
	dst2 := src1.Subtract(dst1)
	conn1 := netset.ConnectionSetFrom(src1, dst1, connection.NewTCPSet())
	conn2 := netset.ConnectionSetFrom(src1, dst2, connection.NewTCPSet())

	unionConn := conn1.Union(conn2)
	conn3 := netset.ConnectionSetFrom(src1, src1, connection.NewTCPSet())

	require.True(t, unionConn.Equal(conn3))
	require.True(t, conn3.Equal(unionConn))

}
