package main

import (
	"fmt"
	"reflect"
)

type Order struct {
	Id string
}

type ArmyOrder struct {
	Order
	ArmyId string
}

type DeployOrder struct {
	ArmyOrder
	Dst string
}

func main() {
	ab := DeployOrder{
		Dst: "destinsat",
		ArmyOrder: ArmyOrder{
			ArmyId: "armyId",
			Order: Order{
				Id: "an ID",
			},
		},
	}

	as := reflect.TypeOf(ab)
	_, ok := as.FieldByName("Order")
	fmt.Println(as, ok)
}
