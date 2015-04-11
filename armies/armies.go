package armies

import (
	"errors"
	"fmt"
	"github.com/pgruenbacher/gotai/regions"
)

type armyId string

type Armies map[armyId]*Army

type Army struct {
	Name             string
	Id               armyId
	Morale           int
	Strength         int
	StartingRegion   regions.RegionId
	HomeRegion       regions.RegionId
	Region           *regions.Region
	Deployed         bool
	CutOff           bool
	Home             *regions.Region
	StartingHostiles []armyId
	Hostiles         []*Army
	parent           Armies
}

// Using config values to initialize the rest of the object
func (self Armies) Init(r regions.Regions) error {
	for _, army := range self {
		// declare starting regions
		if region, ok := r[army.StartingRegion]; ok {
			army.Region = region
		} else {
			return errors.New(fmt.Sprintf("problem initializing %v", army.Id))
		}
		// declare home region
		if region, ok := r[army.HomeRegion]; ok {
			army.Home = region
		} else {
			return errors.New(fmt.Sprintf("problem initializing %v", army.Id))
		}
		// make a pointer to parent for referencing other armies
		army.parent = self
		//  declare starting hostilities between armies for scenarios
		for _, hostileId := range army.StartingHostiles {
			if hostileArmy, ok := self[hostileId]; !ok {
				return errors.New(fmt.Sprintf("hostile army %v doesn't exist for %v", hostileId, army.Id))
			} else if ok = army.isStartingHostileArmy(hostileArmy); !ok {
				return errors.New(fmt.Sprintf("hostile army %v doesn't recipocrate hostility to %v", hostileId, army.Id))
			} else {
				army.Hostiles = append(army.Hostiles, hostileArmy)
			}
		}
	}
	return nil
}

func (self Armies) EvalSupplies(r regions.Regions) error {
	for _, army := range self {
		supplied := army.EvalSupplyRoute(r)
		if !supplied {
			army.CutOff = true
		}
	}
	return nil
}

func (self *Army) Deploy() {
	self.Deployed = true
}

func (self *Army) EvalSupplyRoute(r regions.Regions) bool {
	if self.Region == self.Home {
		return true
	}
	shortPath := r.Path(self.Region.Id, self.Home.Id, nil, nil)
	longPath := r.Path(self.Region.Id, self.Home.Id, hostileFilter, self)
	if len(longPath) != len(shortPath) {
		return false
	}
	return true
}

func (self *Army) March(to *regions.Edge, regs regions.Regions) error {
	if _, err := self.ValidateMarch(to); err != nil {
		return err
	}
	self.Region = regs[to.Dst]
	return nil
}

// returns true (valid region) if there are no hostiles in it
func hostileFilter(region *regions.Region, eval interface{}) bool {
	for _, army := range eval.(*Army).parent {
		if army.Id == eval.(*Army).Id {
			continue
		}
		if hostile := eval.(*Army).isHostileArmy(army); hostile {
			if army.Region.Id == region.Id {
				return false
			}
		}
	}
	return true
}

func (self *Army) isHostileArmy(a *Army) bool {
	for _, hostile := range a.Hostiles {
		if hostile.Id == self.Id {
			return true
		}
	}
	return false
}

func (self *Army) isStartingHostileArmy(a *Army) bool {
	for _, hostileId := range a.StartingHostiles {
		if hostileId == self.Id {
			return true
		}
	}
	return false
}

func (a *Army) ValidateMarch(edge *regions.Edge) (bool, error) {
	if !a.Deployed {
		return false, errors.New(fmt.Sprintf("march invalid: army %v not deployed", a.Name))
	}
	if a.Region.Id != edge.Src {
		return false, errors.New(fmt.Sprintf("march invalid: army %v not located at %v", a.Name, edge.Src))
	}
	return true, nil
}
