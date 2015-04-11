package main

import (
	//"fmt"
	//"github.com/BurntSushi/toml"
	"github.com/pgruenbacher/gotai/actions"
	"github.com/pgruenbacher/gotai/armies"
	"github.com/pgruenbacher/gotai/characters"
	"github.com/pgruenbacher/gotai/regions"
	"github.com/pgruenbacher/gotai/utils"
	"github.com/pgruenbacher/log"
)

func main() {

	var characters characters.Characters
	err := utils.ReadDir("characters", &characters)

	log.Info("characters: %v", len(characters))
	if err != nil {
		log.Error("error", err)
	}

	var regions regions.Regions
	err = utils.ReadDir("regions", &regions)

	log.Info("regions: %v", len(regions))
	if err != nil {
		log.Error("error", err)
	}
	err = regions.ConnectAll()
	if err != nil {
		log.Error("%v", err)
	}

	var armyManager armies.ArmiesManager
	err = utils.ReadFile("./armies/manager.toml", &armyManager)
	if err != nil {
		log.Error("%v", err)
	} else {
		log.Info("%v", armyManager)
	}
	err = armyManager.Init(regions)

	order := actions.Order{Id: "asdf"}
	order1 := armies.ArmyOrder{
		Order:  order,
		ArmyId: "army1",
	}
	// order2 := armies.DeployOrder{order1}

	order3 := armies.MarchOrder{
		ArmyOrder: order1,
		Src:       "region3",
		Dst:       "region2",
	}
	orders := []armies.MarchOrder{order3}
	err = armyManager.ReadOrders(orders)
	// log.Info("%v", order2)

	if err != nil {
		log.Error("%v", err)
	}
}
