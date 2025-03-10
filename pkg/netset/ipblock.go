/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package netset

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

	FirstIPAddressString = "0.0.0.0"
	LastIPAddressString  = "255.255.255.255"

	// internal const  below
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

// NewIPBlock returns a new IPBlock object
func NewIPBlock() *IPBlock {
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
	startIP := int64ToIP4(i.Start())
	endIP := int64ToIP4(i.End())
	return rangeIPstr(startIP, endIP)
}

// toIPRangesList: returns a list of the ip-ranges strings in the current IPBlock object
func (b *IPBlock) toIPRangesList() []string {
	intervals := b.ipRange.Intervals()
	ipRanges := make([]string, len(intervals))
	for index, span := range intervals {
		ipRanges[index] = toIPRange(span)
	}
	return ipRanges
}

// IsSubset checks if this IP block is contained within another IP block.
func (b *IPBlock) IsSubset(other *IPBlock) bool {
	if b == other {
		return true
	}
	return b.ipRange.IsSubset(other.ipRange)
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

func (b *IPBlock) Hash() int {
	return b.ipRange.Hash()
}

func (b *IPBlock) Size() int {
	return b.ipRange.Size()
}

// Subtract returns a new IPBlock from subtraction of input IPBlock from this IPBlock
func (b *IPBlock) Subtract(c *IPBlock) *IPBlock {
	if b == c {
		return NewIPBlock()
	}
	return &IPBlock{
		ipRange: b.ipRange.Subtract(c.ipRange),
	}
}

// Overlap returns whether the two IPBlocks have at least one IP address in common
func (b *IPBlock) Overlap(c *IPBlock) bool {
	return !b.Intersect(c).IsEmpty()
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

func (b *IPBlock) ipCount() int {
	return int(b.ipRange.CalculateSize())
}

// IsSingleIPAddress returns true if this ipblock is a single IP address
func (b *IPBlock) IsSingleIPAddress() bool {
	return b.ipRange.IsSingleNumber()
}

// Compare returns -1 if this<other, 1 if this>other, 0 o.w.
func (b *IPBlock) Compare(other *IPBlock) int {
	switch {
	case b.ipRange.Min() < other.ipRange.Min():
		return -1
	case b.ipRange.Min() > other.ipRange.Min():
		return 1
	case b.ipRange.Max() < other.ipRange.Max():
		return -1
	case b.ipRange.Max() > other.ipRange.Max():
		return 1
	default:
		return 0
	}
}

// Split returns a set of IPBlock objects, each with a single range of ips
func (b *IPBlock) Split() []*IPBlock {
	intervals := b.ipRange.Intervals()
	res := make([]*IPBlock, len(intervals))
	for index, span := range intervals {
		res[index] = &IPBlock{
			ipRange: span.ToSet(),
		}
	}
	return res
}

// SplitToCidrs returns a slice of IPBlocks, each representing a single CIDR
func (b *IPBlock) SplitToCidrs() []*IPBlock {
	cidrs := make([]*IPBlock, 0)
	for _, ipRange := range b.ipRange.Intervals() {
		for _, cidrString := range intervalToCidrList(ipRange) {
			ipblock, _ := IPBlockFromCidr(cidrString)
			cidrs = append(cidrs, ipblock)
		}
	}
	return cidrs
}

// int64ToIP4 returns a string of an ip address from an input integer ip value
func int64ToIP4(ipInt int64) string {
	if ipInt < 0 || ipInt > math.MaxUint32 {
		return "0.0.0.0"
	}
	//nolint:gosec // overflow is not possible
	ipUint32 := uint32(uint64(ipInt) & ipMask)
	var d [4]byte
	binary.BigEndian.PutUint32(d[:], ipUint32)
	return net.IPv4(d[0], d[1], d[2], d[3]).String()
}

// DisjointIPBlocks returns an IPBlock of disjoint ip ranges from 2 input IPBlock objects
func DisjointIPBlocks(set1, set2 []*IPBlock) []*IPBlock {
	ipbList := make([]*IPBlock, len(set1)+len(set2))
	for i, ipb := range set1 {
		ipbList[i] = ipb.Copy()
	}
	for i, ipb := range set2 {
		ipbList[len(set1)+i] = ipb.Copy()
	}
	// sort ipbList by ip_count per netset
	sort.Slice(ipbList, func(i, j int) bool {
		return ipbList[i].ipCount() < ipbList[j].ipCount()
	})
	// making sure the resulting list does not contain overlapping ipBlocks
	var res []*IPBlock
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
	var toAdd []*IPBlock
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

// IPBlockFromCidr returns a new IPBlock object from input CIDR string
func IPBlockFromCidr(cidr string) (*IPBlock, error) {
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
	ipb1, err1 := IPBlockFromCidr(cidr1)
	ipb2, err2 := IPBlockFromCidr(cidr2)
	if err1 != nil || err2 != nil {
		return nil, nil, errors.Join(err1, err2)
	}
	return ipb1, ipb2, nil
}

// IPBlockFromCidrOrAddress returns a new IPBlock object from input string of CIDR or IP address
func IPBlockFromCidrOrAddress(s string) (*IPBlock, error) {
	if strings.Contains(s, cidrSeparator) {
		return IPBlockFromCidr(s)
	}
	return IPBlockFromIPAddress(s)
}

// IPBlockFromIPRange returns a new IPBlock object that contains startIP-endIP
func IPBlockFromIPRange(startIP, endIP *IPBlock) (*IPBlock, error) {
	if !startIP.IsSingleIPAddress() || !endIP.IsSingleIPAddress() {
		return nil, fmt.Errorf("both startIP and endIP should be a single IP address")
	}
	return &IPBlock{
		ipRange: interval.New(startIP.ipRange.Min(), endIP.ipRange.Min()).ToSet(),
	}, nil
}

// IPBlockFromCidrList returns IPBlock object from multiple CIDRs given as list of strings
func IPBlockFromCidrList(cidrsList []string) (*IPBlock, error) {
	res := NewIPBlock()
	for _, cidr := range cidrsList {
		block, err := IPBlockFromCidr(cidr)
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

// IPBlockFromIPAddress returns an IPBlock object from input IP address string
func IPBlockFromIPAddress(ipAddress string) (*IPBlock, error) {
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

// AsCidr returns the CIDR string of this IPBlock object, if it contains exactly one CIDR,
// otherwise it returns an error
func (b *IPBlock) AsCidr() (string, error) {
	cidrList := b.ToCidrList()
	if len(cidrList) != 1 {
		return "", fmt.Errorf("ipblock contains %d cidrs", len(cidrList))
	}
	return cidrList[0], nil
}

// ToCidrList returns a list of CIDR strings for this IPBlock object
func (b *IPBlock) ToCidrList() []string {
	var cidrList []string
	for _, ipRange := range b.ipRange.Intervals() {
		cidrList = append(cidrList, intervalToCidrList(ipRange)...)
	}
	return cidrList
}

// ToCidrListString returns a string with all CIDRs within the IPBlock object
func (b *IPBlock) ToCidrListString() string {
	return strings.Join(b.ToCidrList(), commaSeparator)
}

// ListToPrint returns a uniform to print list s.t. each element contains either a single cidr or an ip range
func (b *IPBlock) ListToPrint() []string {
	var cidrsIPRangesList []string
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

// ToIPAddressString returns the IP Address string for this IPBlock
func (b *IPBlock) ToIPAddressString() string {
	if b.ipRange.IsSingleNumber() {
		return b.FirstIPAddress()
	}
	return ""
}

// FirstIPAddress returns the first IP Address string for this IPBlock
func (b *IPBlock) FirstIPAddress() string {
	return int64ToIP4(b.ipRange.Min())
}

// FirstIPAddressObject returns the first IP Address for this IPBlock
func (b *IPBlock) FirstIPAddressObject() *IPBlock {
	ipNum := b.ipRange.Min()
	return &IPBlock{
		ipRange: interval.New(ipNum, ipNum).ToSet(),
	}
}

// LastIPAddress returns the last IP Address string for this IPBlock
func (b *IPBlock) LastIPAddress() string {
	return int64ToIP4(b.ipRange.Max())
}

// LastIPAddressObject returns the last IP Address for this IPBlock
func (b *IPBlock) LastIPAddressObject() *IPBlock {
	ipNum := b.ipRange.Max()
	return &IPBlock{
		ipRange: interval.New(ipNum, ipNum).ToSet(),
	}
}

// NextIP returns the next ip address after this IPBlock
func (b *IPBlock) NextIP() (*IPBlock, error) {
	if GetLastIPAddress().IsSubset(b) {
		return nil, fmt.Errorf("%s is contained in ipblock", LastIPAddressString)
	}
	lastIP := b.LastIPAddressObject()
	ipNum := lastIP.ipRange.Min() + 1
	return &IPBlock{
		ipRange: interval.New(ipNum, ipNum).ToSet(),
	}, nil
}

// PreviousIP returns the previous ip address before this IPBlock
func (b *IPBlock) PreviousIP() (*IPBlock, error) {
	if GetFirstIPAddress().IsSubset(b) {
		return nil, fmt.Errorf("%s is contained in IPBlock", FirstIPAddressString)
	}
	firstIP := b.FirstIPAddressObject()
	ipNum := firstIP.ipRange.Min() - 1
	return &IPBlock{
		ipRange: interval.New(ipNum, ipNum).ToSet(),
	}, nil
}

func intervalToCidrList(ipRange interval.Interval) []string {
	start := ipRange.Start()
	end := ipRange.End()
	var res []string
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
		ip := int64ToIP4(start)
		res = append(res, fmt.Sprintf("%s/%d", ip, maxSize))
		start += int64(math.Pow(2, maxIPv4Bits-float64(maxSize)))
	}
	return res
}

func parseIP(ip string) (int64, error) {
	startIP := net.ParseIP(ip)
	if startIP == nil || startIP.To4() == nil {
		return 0, fmt.Errorf("%v is not a valid ipv4", ip)
	}
	return int64(binary.BigEndian.Uint32(startIP.To4())), nil
}

// IPBlockFromIPRangeStr returns IPBlock object from input IP range string (example: "169.255.0.0-172.15.255.255")
func IPBlockFromIPRangeStr(ipRangeStr string) (*IPBlock, error) {
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
	res, _ := IPBlockFromCidr(CidrAll)
	return res
}

// GetFirstIPAddress returns IPBlock object of 0.0.0.0
func GetFirstIPAddress() *IPBlock {
	res, _ := IPBlockFromIPAddress(FirstIPAddressString)
	return res
}

// GetLastIPAddress returns IPBlock object of 255.255.255.255
func GetLastIPAddress() *IPBlock {
	res, _ := IPBlockFromIPAddress(LastIPAddressString)
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

// TouchingIPRanges returns true if this and other ipblocks objects are touching.
// assumption: both IPBlocks represent a single IP range
func (b *IPBlock) TouchingIPRanges(other *IPBlock) (bool, error) {
	if b.ipRange.NumIntervals() != 1 || other.ipRange.NumIntervals() != 1 {
		return false, fmt.Errorf("both ipblocks should be a single IP range")
	}
	return !b.Overlap(other) && b.Union(other).ipRange.NumIntervals() == 1, nil
}

func (b *IPBlock) Complementary() *IPBlock {
	allIPBlock := GetCidrAll()
	return allIPBlock.Subtract(b)
}
