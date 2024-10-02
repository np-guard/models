/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package netset_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"
)

func TestBasicICMPSetStrict(t *testing.T) {
	// create ICMPSet objects
	i1 := 8
	all := netset.AllICMPSetStrict()
	obj1, err := netp.ICMPFromTypeAndCode(&i1, nil)
	require.Nil(t, err)
	icmpset := netset.NewICMPSetStrict(obj1)

	// test basic functions, operations
	fmt.Println(icmpset) // ICMP icmp-type: 8 icmp-code: 0
	fmt.Println(all)     // ICMP
	res := icmpset.Union(all)
	fmt.Println(res) // ICMP
	require.True(t, res.Equal(all))
	require.True(t, all.Equal(res))
	require.True(t, icmpset.IsSubset(all))
	fmt.Println("done")
}

func TestBasicICMPSet(t *testing.T) {
	icmpset := netset.NewICMPSet(8, 8, 0, 255) // ICMP type: 8
	icmpset1 := netset.NewICMPSet(8, 8, 0, 0)  // ICMP type: 8 code: 0
	fmt.Println(icmpset)
	fmt.Println(icmpset1)

	require.True(t, icmpset1.IsSubset(icmpset))
	require.False(t, icmpset.IsSubset(icmpset1))
	require.True(t, icmpset1.Union(icmpset).Equal(icmpset))

	require.False(t, icmpset.IsAll())
	require.False(t, icmpset.IsEmpty())

	require.False(t, icmpset1.IsAll())
	require.False(t, icmpset1.IsEmpty())

	fmt.Println("done")
}
