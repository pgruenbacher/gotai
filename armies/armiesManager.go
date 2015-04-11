package armies

import (
	"errors"
	"fmt"
	"github.com/pgruenbacher/gotai/actions"
	"github.com/pgruenbacher/gotai/regions"
	"github.com/pgruenbacher/gotai/utils"
)

// Orders
// Generic Army order has army id
type ArmyOrder struct {
	actions.Order
	ArmyId armyId
}

// specific orders
type MarchOrder struct {
	ArmyOrder
	Src regions.RegionId
	Dst regions.RegionId
}

type DeployOrder struct {
	ArmyOrder
}

type EncampOrder struct {
	ArmyOrder
}

// manager structeure
type ArmiesManager struct {
	Armies  Armies
	regions regions.Regions
	Config  Config
}

type Config struct {
	Penalties []Penalty
}

type Penalty struct {
	Boundary        regions.EdgeBoundary
	CrossingPenalty int
}

func (self *ArmiesManager) Init(r regions.Regions) error {
	fmt.Println(self.Config.Penalties[2].Boundary == regions.Wall)
	self.regions = r
	err := utils.ReadDir("armies", &self.Armies)
	if err != nil {
		return err
	}
	self.Armies.Init(r)
	return nil
}

func (self *ArmiesManager) EvaluateArmies() error {
	if err := self.Armies.EvalSupplies(self.regions); err != nil {
		return err
	}
	return nil

}

func (self *ArmiesManager) ReadOrders(order interface{}) (err error) {
	err = nil
	switch t := order.(type) {
	default:
		// do nothing
	case []MarchOrder:
		err = self.marchOrders(t)

	case []DeployOrder:
		err = self.deployOrders(t)
	}
	return err
}

/*
 * individual order handling
 *
 */

func (self *ArmiesManager) marchOrders(orders []MarchOrder) error {
	if valid, err := self.validateMarchOrders(orders); err != nil {
		return err
	} else if !valid {
		return nil
	}
	for _, order := range orders {
		fmt.Println("march", self.Armies[order.ArmyId])
	}
	return nil
}

func (self *ArmiesManager) deployOrders(orders []DeployOrder) error {
	for _, order := range orders {
		fmt.Println(self.Armies[order.ArmyId])
	}
	return nil
}

/*
 * Validation Section
 *
 */

func (self *ArmiesManager) validateMarchOrders(orders []MarchOrder) (bool, error) {
	for _, order := range orders {
		// validate the army id
		if err := self.validateArmyOrder(order.ArmyOrder); err != nil {
			return false, err
		}

		army := self.Armies[order.ArmyId]

		// validate src region
		if err := self.validateSrc(army, order); err != nil {
			return false, err
		}
		// validate destination region
		if err := self.validateDestination(army, order); err != nil {
			return false, err
		}

	}
	return true, nil
}

func (self *ArmiesManager) validateArmyOrder(order ArmyOrder) error {
	if _, ok := self.Armies[order.ArmyId]; !ok {
		return errors.New(fmt.Sprintf("order %v had invalid armyId %v", order.Id, order.ArmyId))
	}
	return nil
}

func (self *ArmiesManager) validateSrc(army *Army, order MarchOrder) error {
	if _, ok := self.regions[order.Src]; !ok {
		return errors.New(fmt.Sprintf("invalid src id %v", order.Src))
	}
	if army.Region.Id != order.Src {
		return errors.New(fmt.Sprintf("army region %v doesn't match src %v", army.Region.Id, order.Src))
	}
	return nil
}

func (self *ArmiesManager) validateDestination(army *Army, order MarchOrder) error {
	valid := false
	destination, ok := self.regions[order.Dst]
	if !ok {
		return errors.New(fmt.Sprintf("invalid destination id %v", order.Dst))
	}
	for _, edge := range army.Region.Edges {
		if edge.Dst == order.Dst {
			valid = true
		}
	}
	if !valid {
		return errors.New(fmt.Sprintf("none of army region  %v edges don't match army destination %v", army.Region.Id, order.Dst))
	}

	if !hostileFilter(destination, army) {
		return errors.New(fmt.Sprintf("hostile army present in %v must attack order not march", destination.Id))
	}

	return nil
}
