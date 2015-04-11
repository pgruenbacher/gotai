package main

import (
	//"fmt"
	//"github.com/BurntSushi/toml"
	// "github.com/pgruenbacher/gotai/actions"
	// "github.com/pgruenbacher/gotai/armies"
	// "github.com/pgruenbacher/gotai/characters"
	"github.com/pgruenbacher/gotai/prompter"
	// "github.com/pgruenbacher/gotai/regions"
	"github.com/pgruenbacher/gotai/state"
	// "github.com/pgruenbacher/gotai/utils"
	// "github.com/pgruenbacher/log"
)

func main() {
	var game state.State
	game.InitiateGame()
	var prompt prompter.Prompt
	prompt.InitiatePrompt(&game)
}

func mockOrders() {
	// order := actions.Order{Id: "asdf"}
	// order1 := armies.ArmyOrder{
	// 	Order:  order,
	// 	ArmyId: "army1",
	// }
	// order2 := armies.DeployOrder{order1}

	// order3 := armies.MarchOrder{
	// 	ArmyOrder: order1,
	// 	Src:       "region3",
	// 	Dst:       "region2",
	// }
	// // order4 := armies.MarchOrder{
	// // 	ArmyOrder: order1,
	// // 	Src:       "region2",
	// // 	Dst:       "region1",
	// // }
	// orders1 := []armies.DeployOrder{order2}
	// orders2 := []armies.MarchOrder{order3}
	// // orders3 := []armies.MarchOrder{order4}
}
