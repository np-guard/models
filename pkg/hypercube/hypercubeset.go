package hypercube

import (
	"errors"
	"sort"
	"strings"

	"github.com/np-guard/models/pkg/interval"
)

// CanonicalSet is a canonical representation for set of n-dimensional cubes, from integer intervals
type CanonicalSet struct {
	layers     map[*interval.CanonicalSet]*CanonicalSet
	dimensions int
}

// NewCanonicalSet returns a new empty CanonicalSet with n dimensions
func NewCanonicalSet(n int) *CanonicalSet {
	return &CanonicalSet{
		layers:     map[*interval.CanonicalSet]*CanonicalSet{},
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
	if c == other {
		return c.Copy()
	}
	if c.dimensions != other.dimensions {
		return nil
	}
	remainingFromOther := map[*interval.CanonicalSet]*interval.CanonicalSet{}
	for otherKey := range other.layers {
		remainingFromOther[otherKey] = otherKey.Copy()
	}
	layers := map[*interval.CanonicalSet]*CanonicalSet{}
	for k, v := range c.layers {
		remainingFromSelf := k.Copy()
		for otherKey, otherVal := range other.layers {
			commonElem := k.Intersect(otherKey)
			if commonElem.IsEmpty() {
				continue
			}
			remainingFromOther[otherKey] = remainingFromOther[otherKey].Subtract(commonElem)
			remainingFromSelf = remainingFromSelf.Subtract(commonElem)
			newSubElem := NewCanonicalSet(0)
			if c.dimensions != 1 {
				newSubElem = v.Union(otherVal)
			}
			layers[commonElem] = newSubElem
		}
		if !remainingFromSelf.IsEmpty() {
			layers[remainingFromSelf] = v.Copy()
		}
	}
	for k, v := range remainingFromOther {
		if !v.IsEmpty() {
			layers[v] = other.layers[k].Copy()
		}
	}
	return &CanonicalSet{
		layers:     getElementsUnionPerLayer(layers),
		dimensions: c.dimensions,
	}
}

// IsEmpty returns true if c is empty
func (c *CanonicalSet) IsEmpty() bool {
	return len(c.layers) == 0
}

// Intersect returns a new CanonicalSet object that results from intersection of c with other
func (c *CanonicalSet) Intersect(other *CanonicalSet) *CanonicalSet {
	if c == other {
		return c.Copy()
	}
	if c.dimensions != other.dimensions {
		return nil
	}

	layers := map[*interval.CanonicalSet]*CanonicalSet{}
	for k, v := range c.layers {
		for otherKey, otherVal := range other.layers {
			commonELem := k.Intersect(otherKey)
			if commonELem.IsEmpty() {
				continue
			}
			if c.dimensions == 1 {
				layers[commonELem] = NewCanonicalSet(0)
				continue
			}
			newSubElem := v.Intersect(otherVal)
			if !newSubElem.IsEmpty() {
				layers[commonELem] = newSubElem
			}
		}
	}
	return &CanonicalSet{
		layers:     getElementsUnionPerLayer(layers),
		dimensions: c.dimensions,
	}
}

// Subtract returns a new CanonicalSet object that results from subtraction other from c
func (c *CanonicalSet) Subtract(other *CanonicalSet) *CanonicalSet {
	if c == other {
		return NewCanonicalSet(c.dimensions)
	}
	if c.dimensions != other.dimensions {
		return nil
	}
	layers := map[*interval.CanonicalSet]*CanonicalSet{}
	for k, v := range c.layers {
		remainingFromSelf := k.Copy()
		for otherKey, otherVal := range other.layers {
			commonElem := k.Intersect(otherKey)
			if commonElem.IsEmpty() {
				continue
			}
			remainingFromSelf = remainingFromSelf.Subtract(commonElem)
			if c.dimensions == 1 {
				continue
			}
			newSubElem := v.Subtract(otherVal)
			if !newSubElem.IsEmpty() {
				layers[commonElem] = newSubElem
			}
		}
		if !remainingFromSelf.IsEmpty() {
			layers[remainingFromSelf] = v.Copy()
		}
	}
	return &CanonicalSet{
		layers:     getElementsUnionPerLayer(layers),
		dimensions: c.dimensions,
	}
}

func (c *CanonicalSet) getIntervalSetUnion() *interval.CanonicalSet {
	res := interval.NewCanonicalSet()
	for k := range c.layers {
		res.Union(k)
	}
	return res
}

// ContainedIn returns true ic other contained in c
func (c *CanonicalSet) ContainedIn(other *CanonicalSet) (bool, error) {
	if c == other {
		return true, nil
	}
	if c.dimensions != other.dimensions {
		return false, errors.New("ContainedIn mismatch between num of dimensions for input args")
	}
	if c.dimensions == 1 {
		if len(c.layers) != 1 || len(other.layers) != 1 {
			return false, errors.New("unexpected object of dimension size 1")
		}
		cInterval := c.getIntervalSetUnion()
		otherInterval := other.getIntervalSetUnion()
		return cInterval.ContainedIn(otherInterval), nil
	}

	isSubsetCount := 0
	for k, v := range c.layers {
		currentLayer := k.Copy()
		for otherKey, otherVal := range other.layers {
			commonKey := currentLayer.Intersect(otherKey)
			remaining := currentLayer.Subtract(commonKey)
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
		res.layers[k.Copy()] = v.Copy()
	}
	return res
}

func getCubeStr(cube []*interval.CanonicalSet) string {
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

// GetCubesList returns the list of cubes in c, each cube as a slice of CanonicalSet
func (c *CanonicalSet) GetCubesList() [][]*interval.CanonicalSet {
	res := [][]*interval.CanonicalSet{}
	if c.dimensions == 1 {
		for k := range c.layers {
			res = append(res, []*interval.CanonicalSet{k})
		}
		return res
	}
	for k, v := range c.layers {
		subRes := v.GetCubesList()
		for _, subList := range subRes {
			cube := []*interval.CanonicalSet{k}
			cube = append(cube, subList...)
			res = append(res, cube)
		}
	}
	return res
}

func getElementsUnionPerLayer(layers map[*interval.CanonicalSet]*CanonicalSet) map[*interval.CanonicalSet]*CanonicalSet {
	type pair struct {
		hc *CanonicalSet            // hypercube set object
		is []*interval.CanonicalSet // interval-set list
	}
	equivClasses := map[string]*pair{}
	for k, v := range layers {
		if _, ok := equivClasses[v.String()]; ok {
			equivClasses[v.String()].is = append(equivClasses[v.String()].is, k)
		} else {
			equivClasses[v.String()] = &pair{hc: v, is: []*interval.CanonicalSet{k}}
		}
	}
	newLayers := map[*interval.CanonicalSet]*CanonicalSet{}
	for _, p := range equivClasses {
		newVal := p.hc
		newKey := p.is[0]
		for i := 1; i < len(p.is); i += 1 {
			newKey = newKey.Union(p.is[i])
		}
		newLayers[newKey] = newVal
	}
	return newLayers
}

// FromCube returns a new CanonicalSet created from a single input cube
// the input cube is a slice of CanonicalSet, treated as ordered list of dimension values
func FromCube(cube []*interval.CanonicalSet) *CanonicalSet {
	if len(cube) == 0 {
		return nil
	}
	if len(cube) == 1 {
		res := NewCanonicalSet(1)
		res.layers[cube[0].Copy()] = NewCanonicalSet(0)
		return res
	}
	res := NewCanonicalSet(len(cube))
	res.layers[cube[0].Copy()] = FromCube(cube[1:])
	return res
}
