package state

import (
	"github.com/pgruenbacher/gotai/armies"
	// "github.com/pgruenbacher/gotai/characters"
	// "github.com/pgruenbacher/gotai/prompter"
	"github.com/pgruenbacher/gotai/regions"
	"github.com/pgruenbacher/gotai/utils"
	"github.com/pgruenbacher/log"
	"reflect"
)

type State struct {
	ArmiesManager *armies.ArmiesManager
}

type allOrders struct {
	MarchOrders  []armies.MarchOrder
	DeployOrders []armies.DeployOrder
}

func (self *State) InitiateGame() {
	// var characters characters.Characters
	// if err := utils.ReadDir("characters", &characters); err != nil {
	// 	log.Error("%v", err)
	// }

	var regions regions.Regions
	if err := utils.ReadDir("regions", &regions); err != nil {
		log.Error("%v", err)
	}

	if err := regions.ConnectAll(); err != nil {
		log.Error("%v", err)
	}

	var armiesManager armies.ArmiesManager
	if err := armiesManager.Init(regions); err != nil {
		log.Error("%v", err)
	}
	self.ArmiesManager = &armiesManager
}

func (self *State) Next(orders []interface{}) error {
	org := self.organizeOrders(orders)
	v := reflect.ValueOf(org)
	for i := 0; i < v.NumField(); i++ {
		if _, err := self.ArmiesManager.ReadOrders(v.Field(i).Interface()); err != nil {
			log.Error("%v", err)
		}
	}

	// if _, err := self.ArmiesManager.ReadOrders(org.marchOrders); err != nil {
	// 	log.Error("%v", err)
	// }
	return nil
}

func (self *State) organizeOrders(orders []interface{}) allOrders {
	var all allOrders
	for _, order := range orders {
		switch t := order.(type) {
		default:
			// do nothing with orders we don't recognize
		case armies.MarchOrder:
			all.MarchOrders = append(all.MarchOrders, t)
		case armies.DeployOrder:
			all.DeployOrders = append(all.DeployOrders, t)
		}
	}
	return all
}
