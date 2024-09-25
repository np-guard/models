/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package ipblock

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/np-guard/models/pkg/interval"
)

const (
	// CidrAll represents the CIDR for all addresses "0.0.0.0/0"
	CidrAll = "0.0.0.0/0"

	// internal const  below
	ipByte         = 0xff
	ipShift0       = 24
	ipShift1       = 16
	ipShift2       = 8
	ipBase         = 10
	ipMask         = 0xffffffff
	maxIPv4Bits    = 32
	cidrSeparator  = "/"
	bitSize64      = 64
	commaSeparator = ", "
	dash           = "-"
)

// IPBlock captures a set of IP ranges
type IPBlock struct {
	ipRange *interval.CanonicalSet
}

// New returns a new IPBlock object
func New() *IPBlock {
	return &IPBlock{
		ipRange: interval.NewCanonicalSet(),
	}
}

// ToIPRanges returns a string of the ip ranges in the current IPBlock object
func (b *IPBlock) ToIPRanges() string {
	return strings.Join(b.toIPRangesList(), commaSeparator)
}

// toIPRange returns a string of the ip range of a single interval
func toIPRange(i interval.Interval) string {
	startIP := intToIP4(i.Start())
	endIP := intToIP4(i.End())
	return rangeIPstr(startIP, endIP)
}

// toIPRangesList: returns a list of the ip-ranges strings in the current IPBlock object
func (b *IPBlock) toIPRangesList() []string {
	IPRanges := make([]string, b.ipRange.NumIntervals())
	for index, span := range b.ipRange.Intervals() {
		IPRanges[index] = toIPRange(span)
	}
	return IPRanges
}

// ContainedIn checks if this IP block is contained within another IP block.
func (b *IPBlock) ContainedIn(other *IPBlock) bool {
	if b == other {
		return true
	}
	return b.ipRange.ContainedIn(other.ipRange)
}

// Intersect returns a new IPBlock from intersection of this IPBlock with input IPBlock
func (b *IPBlock) Intersect(c *IPBlock) *IPBlock {
	if b == c {
		return b.Copy()
	}
	return &IPBlock{
		ipRange: b.ipRange.Intersect(c.ipRange),
	}
}

// Equal returns true if this IPBlock equals the input IPBlock
func (b *IPBlock) Equal(c *IPBlock) bool {
	if b == c {
		return true
	}
	return b.ipRange.Equal(c.ipRange)
}

// Subtract returns a new IPBlock from subtraction of input IPBlock from this IPBlock
func (b *IPBlock) Subtract(c *IPBlock) *IPBlock {
	if b == c {
		return New()
	}
	return &IPBlock{
		ipRange: b.ipRange.Subtract(c.ipRange),
	}
}

// Union returns a new IPBlock from union of input IPBlock with this IPBlock
func (b *IPBlock) Union(c *IPBlock) *IPBlock {
	if b == c {
		return b.Copy()
	}
	return &IPBlock{
		ipRange: b.ipRange.Union(c.ipRange),
	}
}

// Overlap returns whether the two IPBlocks have at least one IP address in common
func (b *IPBlock) Overlap(c *IPBlock) bool {
	return !b.Intersect(c).IsEmpty()
}

// IsEmpty returns true if this IPBlock is empty
func (b *IPBlock) IsEmpty() bool {
	return b.ipRange.IsEmpty()
}

func rangeIPstr(start, end string) string {
	return fmt.Sprintf("%v-%v", start, end)
}

// Copy returns a new copy of IPBlock object
func (b *IPBlock) Copy() *IPBlock {
	return &IPBlock{ipRange: b.ipRange.Copy()}
}

func (b *IPBlock) IPCount() int {
	return int(b.ipRange.CalculateSize())
}

// Split returns a set of IpBlock objects, each with a single range of ips
func (b *IPBlock) Split() []*IPBlock {
	res := make([]*IPBlock, b.ipRange.NumIntervals())
	for index, span := range b.ipRange.Intervals() {
		res[index] = &IPBlock{
			ipRange: span.ToSet(),
		}
	}
	return res
}

// intToIP4 returns a string of an ip address from an input integer ip value
func intToIP4(ipInt int64) string {
	if ipInt < 0 || ipInt > math.MaxUint32 {
		return "0.0.0.0"
	}
	var d [4]byte
	binary.BigEndian.PutUint32(d[:], uint32(ipInt))
	return net.IPv4(d[0], d[1], d[2], d[3]).String()
}

// DisjointIPBlocks returns an IPBlock of disjoint ip ranges from 2 input IPBlock objects
func DisjointIPBlocks(set1, set2 []*IPBlock) []*IPBlock {
	ipbList := []*IPBlock{}
	for _, ipb := range set1 {
		ipbList = append(ipbList, ipb.Copy())
	}
	for _, ipb := range set2 {
		ipbList = append(ipbList, ipb.Copy())
	}
	// sort ipbList by ip_count per ipblock
	sort.Slice(ipbList, func(i, j int) bool {
		return ipbList[i].IPCount() < ipbList[j].IPCount()
	})
	// making sure the resulting list does not contain overlapping ipBlocks
	res := []*IPBlock{}
	for _, ipb := range ipbList {
		res = addIntervalToList(ipb, res)
	}

	if len(res) == 0 {
		res = []*IPBlock{GetCidrAll()}
	}
	return res
}

// addIntervalToList is used for computation of DisjointIPBlocks
func addIntervalToList(ipbNew *IPBlock, ipbList []*IPBlock) []*IPBlock {
	toAdd := []*IPBlock{}
	for idx, ipb := range ipbList {
		if !ipb.ipRange.Overlap(ipbNew.ipRange) {
			continue
		}
		intersection := ipb.Intersect(ipbNew)
		ipbNew = ipbNew.Subtract(intersection)
		if !ipb.Equal(intersection) {
			toAdd = append(toAdd, intersection)
			ipbList[idx] = ipbList[idx].Subtract(intersection)
		}
		if ipbNew.IsEmpty() {
			break
		}
	}
	ipbList = append(ipbList, ipbNew.Split()...)
	ipbList = append(ipbList, toAdd...)
	return ipbList
}

// FromCidr returns a new IPBlock object from input CIDR string
func FromCidr(cidr string) (*IPBlock, error) {
	span, err := cidrToInterval(cidr)
	if err != nil {
		return nil, err
	}
	return &IPBlock{
		ipRange: span.ToSet(),
	}, nil
}

// PairCIDRsToIPBlocks returns two IPBlock objects from two input CIDR strings
func PairCIDRsToIPBlocks(cidr1, cidr2 string) (ipb1, ipb2 *IPBlock, err error) {
	ipb1, err1 := FromCidr(cidr1)
	ipb2, err2 := FromCidr(cidr2)
	if err1 != nil || err2 != nil {
		return nil, nil, errors.Join(err1, err2)
	}
	return ipb1, ipb2, nil
}

// FromCidrOrAddress returns a new IPBlock object from input string of CIDR or IP address
func FromCidrOrAddress(s string) (*IPBlock, error) {
	if strings.Contains(s, cidrSeparator) {
		return FromCidr(s)
	}
	return FromIPAddress(s)
}

// FromCidrList returns IPBlock object from multiple CIDRs given as list of strings
func FromCidrList(cidrsList []string) (*IPBlock, error) {
	res := New()
	for _, cidr := range cidrsList {
		block, err := FromCidr(cidr)
		if err != nil {
			return nil, err
		}
		res = res.Union(block)
	}
	return res, nil
}

// ExceptCidrs returns a new IPBlock with all cidr ranges removed
func (b *IPBlock) ExceptCidrs(exceptions ...string) (*IPBlock, error) {
	holes := interval.NewCanonicalSet()
	for i := range exceptions {
		intervalHole, err := cidrToInterval(exceptions[i])
		if err != nil {
			return nil, err
		}
		holes.AddInterval(intervalHole)
	}
	return &IPBlock{ipRange: b.ipRange.Subtract(holes)}, nil
}

// FromIPAddress returns an IPBlock object from input IP address string
func FromIPAddress(ipAddress string) (*IPBlock, error) {
	ipNum, err := parseIP(ipAddress)
	if err != nil {
		return nil, err
	}
	return &IPBlock{
		ipRange: interval.New(ipNum, ipNum).ToSet(),
	}, nil
}

func cidrToInterval(cidr string) (interval.Interval, error) {
	// convert string to IPNet struct
	_, ipv4Net, err := net.ParseCIDR(cidr)
	if err != nil {
		return interval.Interval{}, err
	}

	// convert IPNet struct mask and address to uint32
	// network is BigEndian
	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	startNum := binary.BigEndian.Uint32(ipv4Net.IP)
	// find the final address
	endNum := (startNum & mask) | (mask ^ ipMask)
	return interval.New(int64(startNum), int64(endNum)), nil
}

// ToCidrList returns a list of CIDR strings for this IPBlock object
func (b *IPBlock) ToCidrList() []string {
	cidrList := []string{}
	for _, ipRange := range b.ipRange.Intervals() {
		cidrList = append(cidrList, intervalToCidrList(ipRange)...)
	}
	return cidrList
}

// ToCidrListString returns a string with all CIDRs within the IPBlock object
func (b *IPBlock) ToCidrListString() string {
	return strings.Join(b.ToCidrList(), commaSeparator)
}

// ListToPrint: returns a uniform to print list s.t. each element contains either a single cidr or an ip range
// todo - remove this func, and put its code under ToRangesList()
func (b *IPBlock) ListToPrint() []string {
	cidrsIPRangesList := []string{}
	for _, ipRange := range b.ipRange.Intervals() {
		cidr := intervalToCidrList(ipRange)
		if len(cidr) == 1 {
			cidrsIPRangesList = append(cidrsIPRangesList, cidr[0])
		} else {
			cidrsIPRangesList = append(cidrsIPRangesList, toIPRange(ipRange))
		}
	}
	return cidrsIPRangesList
}

func (b *IPBlock) ToRangesList() []string {
	return b.ListToPrint()
}

func (b *IPBlock) ToRangesListString() string {
	return strings.Join(b.ToRangesList(), commaSeparator)
}

// ToIPAdressString returns the IP Address string for this IPBlock
func (b *IPBlock) ToIPAddressString() string {
	if b.ipRange.IsSingleNumber() {
		return b.FirstIPAddress()
	}
	return ""
}

// FirstIPAddress() returns the first IP Address string for this IPBlock
func (b *IPBlock) FirstIPAddress() string {
	return intToIP4(b.ipRange.Min())
}

func intervalToCidrList(ipRange interval.Interval) []string {
	start := ipRange.Start()
	end := ipRange.End()
	res := []string{}
	for end >= start {
		maxSize := maxIPv4Bits
		for maxSize > 0 {
			s := maxSize - 1
			mask := int64(math.Round(math.Pow(2, maxIPv4Bits) - math.Pow(2, float64(maxIPv4Bits)-float64(s))))
			maskBase := start & mask
			if maskBase != start {
				break
			}
			maxSize--
		}
		x := math.Log(float64(end)-float64(start)+1) / math.Log(2)
		maxDiff := byte(maxIPv4Bits - math.Floor(x))
		if maxSize < int(maxDiff) {
			maxSize = int(maxDiff)
		}
		ip := intToIP4(start)
		res = append(res, fmt.Sprintf("%s/%d", ip, maxSize))
		start += int64(math.Pow(2, maxIPv4Bits-float64(maxSize)))
	}
	return res
}

func parseIP(ip string) (int64, error) {
	startIP := net.ParseIP(ip)
	if startIP == nil {
		return 0, fmt.Errorf("%v is not a valid ipv4", ip)
	}
	return int64(binary.BigEndian.Uint32(startIP.To4())), nil
}

// FromIPRangeStr returns IPBlock object from input IP range string (example: "169.255.0.0-172.15.255.255")
func FromIPRangeStr(ipRangeStr string) (*IPBlock, error) {
	ipAddresses := strings.Split(ipRangeStr, dash)
	if len(ipAddresses) != 2 {
		return nil, errors.New("unexpected ipRange str")
	}
	startIPNum, err0 := parseIP(ipAddresses[0])
	endIPNum, err1 := parseIP(ipAddresses[1])
	if err0 != nil || err1 != nil {
		return nil, errors.Join(err0, err1)
	}
	res := &IPBlock{
		ipRange: interval.New(startIPNum, endIPNum).ToSet(),
	}
	return res, nil
}

// GetCidrAll returns IPBlock object of the entire range 0.0.0.0/0
func GetCidrAll() *IPBlock {
	res, _ := FromCidr(CidrAll)
	return res
}

// PrefixLength returns the cidr's prefix length, assuming the ipBlock is exactly one cidr.
// Prefix length specifies the number of bits in the IP address that are to be used as the subnet mask.
func (b *IPBlock) PrefixLength() (int64, error) {
	cidrs := b.ToCidrList()
	if len(cidrs) != 1 {
		return 0, errors.New("prefixLength err: ipBlock is not a single cidr")
	}
	cidrStr := cidrs[0]
	lenStr := strings.Split(cidrStr, cidrSeparator)[1]
	return strconv.ParseInt(lenStr, ipBase, bitSize64)
}

// String returns an IPBlock's string -- either single IP address, or list of CIDR strings
func (b *IPBlock) String() string {
	if b.ipRange.IsSingleNumber() {
		return b.FirstIPAddress()
	}
	return b.ToCidrListString()
}
