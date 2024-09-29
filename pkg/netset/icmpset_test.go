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

func TestBasicICMPSet(t *testing.T) {
	// create ICMPSet objects
	i1 := int(8)
	all := netset.AllICMPSet()
	obj1, err := netp.ICMPFromTypeAndCode(&i1, nil)
	require.Nil(t, err)
	icmpset := netset.NewICMPSet(obj1)

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
