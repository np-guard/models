package connection

/*func TestICMPTypeOutOfRange(t *testing.T) {
	i1 := int64(135)
	i2 := int64(136)
	c1, _ := ICMPConnection(&i1, nil)
	c2, _ := ICMPConnection(&i2, nil)
	c3 := netset.AllICMPSet()
	c4 := c3.Union(c1.Union(c2))
	require.Equal(t, c4, c3)
}*/

/*
func TestICMPTypeOutOfRange(t *testing.T) {
	c1 := connection.ICMPConnection(135, 136, connection.MinICMPCode, connection.MaxICMPCode)
	c2 := connection.ICMPConnection(connection.MinICMPType, connection.MaxICMPType, connection.MinICMPCode, connection.MaxICMPCode)
	union := c1.Union(c2)
	fmt.Println(c1.String())
	fmt.Println(c2.String())
	fmt.Println(union.String())
	require.Equal(t, "protocol: ICMP icmp-type: 135-136", c1.String())
	require.Equal(t, "protocol: ICMP", c2.String())
	require.Equal(t, "protocol: ICMP", union.String())
}

*/
