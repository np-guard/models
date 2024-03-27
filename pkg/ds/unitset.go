// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ds

type UnitSet struct {
	empty bool
}

func NewUnitSet() *UnitSet {
	return &UnitSet{empty: true}
}

func (m *UnitSet) Equal(other *UnitSet) bool {
	return m.empty == other.empty
}

func (m *UnitSet) Copy() *UnitSet {
	return &UnitSet{empty: m.empty}
}

func (m *UnitSet) Hash() int {
	return m.Size()
}

func (m *UnitSet) IsEmpty() bool {
	return m.empty
}

func (m *UnitSet) IsSubset(other *UnitSet) bool {
	return m.empty || !other.empty
}

func (m *UnitSet) Size() int {
	if m.empty {
		return 0
	}
	return 1
}

func (m *UnitSet) Union(other *UnitSet) *UnitSet {
	return &UnitSet{empty: m.empty && other.empty}
}

func (m *UnitSet) Intersect(other *UnitSet) *UnitSet {
	return &UnitSet{empty: m.empty || other.empty}
}

func (m *UnitSet) Subtract(other *UnitSet) *UnitSet {
	return &UnitSet{empty: m.IsSubset(other)}
}
