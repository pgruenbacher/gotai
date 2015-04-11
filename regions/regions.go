package regions

import (
	"errors"
	"fmt"
)

type Regions map[RegionId]*Region

type RegionId string

/*
Nodes contain units, and are connected by edges.
*/
type Region struct {
	// Id is the id of the region
	Id RegionId
	// Size is the number of units supported by the region.
	Size int
	// Supply, if army(s) exceed supply limit and have no supply train, will begin to starve
	Supply int
	// Edges go from this region to others.
	Edges map[RegionId]Edge
	// neighbors
	Neighbors []RegionId
	Rivers    []RegionId
	Mountains []RegionId
	Walls     []RegionId
}

type EdgeBoundary string

const (
	River      EdgeBoundary = "River"
	Mountain   EdgeBoundary = "Mountain"
	Wall       EdgeBoundary = "Wall"
	NoBoundary EdgeBoundary = "None"
)

/*
Edge goes from one node to another.
*/
type Edge struct {
	// Src is the id of the source node.
	Src RegionId
	// Dst is the id of the destination node.
	Dst RegionId
	// Whether there is a river/mountain/wall between the regions.
	Boundary EdgeBoundary
}

func (self Regions) initializeAll() bool {
	for _, region := range self {
		region.Edges = make(map[RegionId]Edge, len(region.Neighbors))
	}
	return true
}

func (self Regions) ConnectAll() error {
	if ok := self.initializeAll(); !ok {
		return errors.New("couldn't initialize")
	}
	for regionId, region := range self {
		if regionId != region.Id {
			return errors.New(fmt.Sprintf("region id %v doesn't match region %v", region.Id, regionId))
		}
		if len(region.Neighbors) == 0 {
			return errors.New(fmt.Sprintf("region %v has no neighbors", regionId))
		}
		for _, neighborId := range region.Neighbors {
			if neighbor, ok := self[neighborId]; ok {
				if neighbor.Id == region.Id {
					return errors.New(fmt.Sprintf("neighbor %v can't neighbor itself", regionId))
				}
				err := validateNeighbor(region, neighbor)
				if err != nil {
					return err
				}
				err = region.Connect(neighbor)
				if err != nil {
					return err
				}
			} else {
				return errors.New(fmt.Sprintf("neighbor %v does not exist", neighbor))
			}
		}
	}
	return nil
}

func (self *Region) Connect(region *Region) error {
	bound, err := validateBoundaries(self, region)
	if err != nil {
		return err
	}
	away := &Edge{
		Src:      self.Id,
		Dst:      region.Id,
		Boundary: bound,
	}
	here := &Edge{
		Src:      region.Id,
		Dst:      self.Id,
		Boundary: bound,
	}
	self.Edges[region.Id] = *away
	region.Edges[self.Id] = *here
	return nil
}

func validateNeighbor(a, b *Region) error {
	if valid := checkNeighbors(a.Id, b.Neighbors); !valid {
		return errors.New(fmt.Sprintf("neighbor %v doesn't reference region %v", a.Id, b.Id))
	}
	return nil
}

func validateBoundaries(a, b *Region) (EdgeBoundary, error) {
	for _, neighborId := range a.Rivers {
		if neighborId != b.Id {
			continue
		}
		if valid := checkNeighbors(a.Id, b.Rivers); !valid {
			return NoBoundary, errors.New(fmt.Sprintf("neighbor %v river doesn't reference region %v", a.Id, b.Id))
		}
		return River, nil
	}
	for _, neighborId := range a.Mountains {
		if neighborId != b.Id {
			continue
		}
		if valid := checkNeighbors(a.Id, b.Mountains); !valid {
			return NoBoundary, errors.New(fmt.Sprintf("neighbor %v mountains doesn't reference region %v", a.Id, b.Id))
		}
		return Mountain, nil
	}
	for _, neighborId := range a.Walls {
		if neighborId != b.Id {
			continue
		}
		if valid := checkNeighbors(a.Id, b.Walls); !valid {
			return NoBoundary, errors.New(fmt.Sprintf("neighbor %v walls doesn't reference region %v", a.Id, b.Id))
		}
		return Wall, nil
	}
	return NoBoundary, nil
}

func checkNeighbors(a RegionId, list []RegionId) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
