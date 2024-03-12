package hypercube

import (
	"errors"
	"sort"
	"strings"

	"github.com/np-guard/models/pkg/interval"
)

// CanonicalSet is a canonical representation for set of n-dimensional cubes, from integer intervals
type CanonicalSet struct {
	layers     map[*interval.CanonicalIntervalSet]*CanonicalSet
	dimensions int
}

// NewCanonicalSet returns a new empty CanonicalSet with n dimensions
func NewCanonicalSet(n int) *CanonicalSet {
	return &CanonicalSet{
		layers:     map[*interval.CanonicalIntervalSet]*CanonicalSet{},
		dimensions: n,
	}
}

// Equal return true if c equals other (same canonical form)
func (c *CanonicalSet) Equal(other *CanonicalSet) bool {
	if c == other {
		return true
	}
	if c.dimensions != other.dimensions {
		return false
	}
	if len(c.layers) != len(other.layers) {
		return false
	}
	if len(c.layers) == 0 {
		return true
	}
	mapByString := map[string]*CanonicalSet{}
	for k, v := range c.layers {
		mapByString[k.String()] = v
	}
	for k, v := range other.layers {
		if w, ok := mapByString[k.String()]; !ok || !v.Equal(w) {
			return false
		}
	}
	return true
}

// Union returns a new CanonicalSet object that results from union of c with other
func (c *CanonicalSet) Union(other *CanonicalSet) *CanonicalSet {
	if c.dimensions != other.dimensions {
		return nil
	}
	res := NewCanonicalSet(c.dimensions)
	remainingFromOther := map[*interval.CanonicalIntervalSet]*interval.CanonicalIntervalSet{}
	for k := range other.layers {
		kCopy := k.Copy()
		remainingFromOther[k] = &kCopy
	}
	for k, v := range c.layers {
		remainingFromSelf := copyIntervalSet(k)
		for otherKey, otherVal := range other.layers {
			commonElem := copyIntervalSet(k)
			commonElem.Intersect(*otherKey)
			if commonElem.IsEmpty() {
				continue
			}
			remainingFromOther[otherKey].Subtraction(*commonElem)
			remainingFromSelf.Subtraction(*commonElem)
			if c.dimensions == 1 {
				res.layers[commonElem] = NewCanonicalSet(0)
				continue
			}
			newSubElem := v.Union(otherVal)
			res.layers[commonElem] = newSubElem
		}
		if !remainingFromSelf.IsEmpty() {
			res.layers[remainingFromSelf] = v.Copy()
		}
	}
	for k, v := range remainingFromOther {
		if !v.IsEmpty() {
			res.layers[v] = other.layers[k].Copy()
		}
	}
	res.applyElementsUnionPerLayer()
	return res
}

// IsEmpty returns true if c is empty
func (c *CanonicalSet) IsEmpty() bool {
	return len(c.layers) == 0
}

// Intersect returns a new CanonicalSet object that results from intersection of c with other
func (c *CanonicalSet) Intersect(other *CanonicalSet) *CanonicalSet {
	if c.dimensions != other.dimensions {
		return nil
	}
	res := NewCanonicalSet(c.dimensions)
	for k, v := range c.layers {
		for otherKey, otherVal := range other.layers {
			commonELem := copyIntervalSet(k)
			commonELem.Intersect(*otherKey)
			if commonELem.IsEmpty() {
				continue
			}
			if c.dimensions == 1 {
				res.layers[commonELem] = NewCanonicalSet(0)
				continue
			}
			newSubElem := v.Intersect(otherVal)
			if !newSubElem.IsEmpty() {
				res.layers[commonELem] = newSubElem
			}
		}
	}
	res.applyElementsUnionPerLayer()
	return res
}

// Subtraction returns a new CanonicalSet object that results from subtraction other from c
func (c *CanonicalSet) Subtraction(other *CanonicalSet) *CanonicalSet {
	if c.dimensions != other.dimensions {
		return nil
	}
	res := NewCanonicalSet(c.dimensions)
	for k, v := range c.layers {
		remainingFromSelf := copyIntervalSet(k)
		for otherKey, otherVal := range other.layers {
			commonELem := copyIntervalSet(k)
			commonELem.Intersect(*otherKey)
			if commonELem.IsEmpty() {
				continue
			}
			remainingFromSelf.Subtraction(*commonELem)
			if c.dimensions == 1 {
				continue
			}
			newSubElem := v.Subtraction(otherVal)
			if !newSubElem.IsEmpty() {
				res.layers[commonELem] = newSubElem
			}
		}
		if !remainingFromSelf.IsEmpty() {
			res.layers[remainingFromSelf] = v.Copy()
		}
	}
	res.applyElementsUnionPerLayer()
	return res
}

func (c *CanonicalSet) getIntervalSetUnion() *interval.CanonicalIntervalSet {
	res := interval.NewCanonicalIntervalSet()
	for k := range c.layers {
		res.Union(*k)
	}
	return res
}

// ContainedIn returns true ic other contained in c
func (c *CanonicalSet) ContainedIn(other *CanonicalSet) (bool, error) {
	if c.dimensions != other.dimensions {
		return false, errors.New("ContainedIn mismatch between num of dimensions for input args")
	}
	if c.dimensions == 1 {
		if len(c.layers) != 1 || len(other.layers) != 1 {
			return false, errors.New("unexpected object of dimension size 1")
		}
		cInterval := c.getIntervalSetUnion()
		otherInterval := other.getIntervalSetUnion()
		return cInterval.ContainedIn(*otherInterval), nil
	}

	isSubsetCount := 0
	for k, v := range c.layers {
		currentLayer := copyIntervalSet(k)
		for otherKey, otherVal := range other.layers {
			commonKey := copyIntervalSet(currentLayer)
			commonKey.Intersect(*otherKey)
			remaining := copyIntervalSet(currentLayer)
			remaining.Subtraction(*commonKey)
			if !commonKey.IsEmpty() {
				subContainment, err := v.ContainedIn(otherVal)
				if !subContainment || err != nil {
					return subContainment, err
				}
				if !remaining.IsEmpty() {
					currentLayer = remaining
				} else {
					isSubsetCount += 1
					break
				}
			}
		}
	}
	return isSubsetCount == len(c.layers), nil
}

// Copy returns a new CanonicalSet object, copied from c
func (c *CanonicalSet) Copy() *CanonicalSet {
	res := NewCanonicalSet(c.dimensions)
	for k, v := range c.layers {
		newKey := k.Copy()
		res.layers[&newKey] = v.Copy()
	}
	return res
}

func getCubeStr(cube []*interval.CanonicalIntervalSet) string {
	strList := []string{}
	for _, v := range cube {
		strList = append(strList, "("+v.String()+")")
	}
	return "[" + strings.Join(strList, ",") + "]"
}

// String returns a string representation of c
func (c *CanonicalSet) String() string {
	cubesList := c.GetCubesList()
	strList := []string{}
	for _, cube := range cubesList {
		strList = append(strList, getCubeStr(cube))
	}
	sort.Strings(strList)
	return strings.Join(strList, "; ")
}

// GetCubesList returns the list of cubes in c, each cube as a slice of CanonicalIntervalSet
func (c *CanonicalSet) GetCubesList() [][]*interval.CanonicalIntervalSet {
	res := [][]*interval.CanonicalIntervalSet{}
	if c.dimensions == 1 {
		for k := range c.layers {
			res = append(res, []*interval.CanonicalIntervalSet{k})
		}
		return res
	}
	for k, v := range c.layers {
		subRes := v.GetCubesList()
		for _, subList := range subRes {
			cube := []*interval.CanonicalIntervalSet{k}
			cube = append(cube, subList...)
			res = append(res, cube)
		}
	}
	return res
}

func (c *CanonicalSet) applyElementsUnionPerLayer() {
	type pair struct {
		hc *CanonicalSet                    // hypercube set object
		is []*interval.CanonicalIntervalSet // interval-set list
	}
	equivClasses := map[string]*pair{}
	for k, v := range c.layers {
		if _, ok := equivClasses[v.String()]; ok {
			equivClasses[v.String()].is = append(equivClasses[v.String()].is, k)
		} else {
			equivClasses[v.String()] = &pair{hc: v, is: []*interval.CanonicalIntervalSet{k}}
		}
	}
	newLayers := map[*interval.CanonicalIntervalSet]*CanonicalSet{}
	for _, p := range equivClasses {
		newVal := p.hc
		newKey := p.is[0]
		for i := 1; i < len(p.is); i += 1 {
			newKey.Union(*p.is[i])
		}
		newLayers[newKey] = newVal
	}
	c.layers = newLayers
}

// FromCube returns a new CanonicalSet created from a single input cube
// the input cube is a slice of CanonicalIntervalSet, treated as ordered list of dimension values
func FromCube(cube []*interval.CanonicalIntervalSet) *CanonicalSet {
	if len(cube) == 0 {
		return nil
	}
	if len(cube) == 1 {
		res := NewCanonicalSet(1)
		cubeVal := cube[0].Copy()
		res.layers[&cubeVal] = NewCanonicalSet(0)
		return res
	}
	res := NewCanonicalSet(len(cube))
	cubeVal := cube[0].Copy()
	res.layers[&cubeVal] = FromCube(cube[1:])
	return res
}

func FromCubeAsIntervals(values ...*interval.CanonicalIntervalSet) *CanonicalSet {
	return FromCube(values)
}

// FromCubeShort returns a new CanonicalSet created from a single input cube
// the input cube is given as an ordered list of integer values, where each two values
// represent the range (start,end) for a dimension value
func FromCubeShort(values ...int64) *CanonicalSet {
	cube := []*interval.CanonicalIntervalSet{}
	for i := 0; i < len(values); i += 2 {
		cube = append(cube, interval.FromInterval(values[i], values[i+1]))
	}
	return FromCube(cube)
}

func copyIntervalSet(a *interval.CanonicalIntervalSet) *interval.CanonicalIntervalSet {
	res := a.Copy()
	return &res
}
