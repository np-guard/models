/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package netset_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netset"
)

func TestDiscreteTrafficSetBasicOperations(t *testing.T) {
	s1 := interval.New(0, 1).ToSet()
	s2 := interval.New(2, 3).ToSet()
	s3 := s1.Union(s2)
	s4 := interval.New(1, 1).ToSet()
	s5 := interval.New(0, 0).ToSet()
	fmt.Println(s3.String())

	conn1 := netset.NewDiscreteEndpointsTrafficSet(s1, s2, netset.AllTCPTransport())
	conn2 := netset.NewDiscreteEndpointsTrafficSet(s4, s2, netset.AllTCPTransport())
	conn3 := netset.NewDiscreteEndpointsTrafficSet(s5, s2, netset.AllTCPTransport())
	require.Equal(t, conn3.Union(conn2), conn1)
	require.Equal(t, conn1.Subtract(conn2), conn3)
	require.Equal(t, conn1.Subtract(conn3), conn2)
	require.True(t, conn1.Subtract(conn2).Subtract(conn3).IsEmpty())
	fmt.Println(conn1.String())
	fmt.Println(conn2.String())
	fmt.Println((conn1.Subtract(conn2)).String())
	fmt.Println((conn2.Subtract(conn1)).String())

	fmt.Println("done")
}
