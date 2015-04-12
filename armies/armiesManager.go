package armies

import (
	"errors"
	"fmt"
	"github.com/pgruenbacher/gotai/actions"
	"github.com/pgruenbacher/gotai/events"
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

// Events
// Generic army event has army id
type ArmyEvent struct {
	events.Event
	ArmyId armyId
}

// specific events
type MarchEvent struct {
	ArmyEvent
	Src regions.RegionId
	Dst regions.RegionId
}

type DeployEvent struct {
	ArmyEvent
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

// Armies Manager methods
func (self *ArmiesManager) Init(r regions.Regions) error {
	if err := utils.ReadFile("./armies/manager.toml", self); err != nil {
		return err
	}
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

func (self *ArmiesManager) ReadOrders(orders interface{}) (events interface{}, err error) {
	err = nil

	switch t := orders.(type) {
	default:
		// do nothing
	case []MarchOrder:
		events, err = self.marchOrders(t)

	case []DeployOrder:
		events, err = self.deployOrders(t)
	}
	return events, err
}

func (self *ArmiesManager) GivePossibleOrders(id armyId) (orders []interface{}) {
	army := self.Armies[id]
	// if army is not deployed, give movement orders
	if !army.Deployed {
		deployOrder := DeployOrder{
			ArmyOrder: newArmyOrder(army.Id),
		}
		orders = append(orders, deployOrder)
	} else {
		for _, to := range army.AvailableEdges() {
			moveOrder := MarchOrder{
				ArmyOrder: newArmyOrder(army.Id),
				Src:       to.Src.Id,
				Dst:       to.Dst.Id,
			}
			orders = append(orders, moveOrder)
		}
	}
	return orders
}

/*
 * individual order handling
 *
 */

func (self *ArmiesManager) marchOrders(orders []MarchOrder) (e []MarchEvent, err error) {
	if err = self.validateMarchOrders(orders); err != nil {
		return e, err
	}
	for _, order := range orders {
		// get the army to perform on
		army := self.Armies[order.ArmyId]
		// Perform the march
		edge := army.Region.Edges[order.Dst]
		if err = army.March(&edge); err != nil {
			return e, err
		}
		// make the event
		event := MarchEvent{
			ArmyEvent: newArmyEvent(army.Id),
			Src:       edge.Src.Id,
			Dst:       edge.Dst.Id,
		}
		e = append(e, event)

	}
	return e, nil
}

func (self *ArmiesManager) deployOrders(orders []DeployOrder) (e []DeployEvent, err error) {
	for _, order := range orders {
		if err = self.Armies[order.ArmyId].Deploy(); err != nil {
			return e, err
		}
		event := DeployEvent{
			ArmyEvent: newArmyEvent(order.ArmyId),
		}
		e = append(e, event)
	}
	return e, nil
}

/*
 * Logic Section
 */

func (self *ArmiesManager) enactBoundaryPenalty(interface{}) {

}

/*
 * Validation Section
 *
 */

func (self *ArmiesManager) validateMarchOrders(orders []MarchOrder) error {
	for _, order := range orders {
		// validate the army id
		if err := self.validateArmyOrder(order.ArmyOrder); err != nil {
			return err
		}

		army := self.Armies[order.ArmyId]

		// validate src region
		if err := self.validateSrc(army, order); err != nil {
			return err
		}
		// validate destination region
		if err := self.validateDestination(army, order); err != nil {
			return err
		}

	}
	return nil
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
		if edge.Dst.Id == order.Dst {
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

/*
 *
 * Utilities
 *
 */

func newArmyEvent(id armyId) ArmyEvent {
	return ArmyEvent{
		Event:  events.NewEvent(),
		ArmyId: id,
	}
}

func newArmyOrder(id armyId) ArmyOrder {
	return ArmyOrder{
		Order:  actions.NewOrder(),
		ArmyId: id,
	}
}
