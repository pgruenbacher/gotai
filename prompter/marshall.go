package prompter

import (
	"fmt"
	// "log"
	// "strings"

	// "github.com/jroimartin/gocui"
	"github.com/pgruenbacher/gotai/armies"
)

type Marshall struct {
	id          string
	ArmyManager *armies.ArmiesManager
}

type marchOption struct {
	order armies.MarchOrder
}

func (self marchOption) Option() option {
	return option{
		txt:   fmt.Sprintf("march to %v", self.order.Dst),
		id:    "asfd",
		order: self.order,
	}
}

type deployOption struct {
	order armies.DeployOrder
}

func (self deployOption) Option() option {
	return option{
		txt:   fmt.Sprintf("deploy"),
		id:    "asdf",
		order: self.order,
	}
}

func (self *Marshall) GiveArmyOptions() (options []validOption) {
	orders := self.ArmyManager.GivePossibleOrders("army1")
	for _, order := range orders {
		switch t := order.(type) {
		default:
			// give no option
		case armies.MarchOrder:
			opt := marchOption{
				order: t,
			}
			options = append(options, opt)
		case armies.DeployOrder:
			opt := deployOption{
				order: t,
			}
			options = append(options, opt)
		}
	}

	return options
}
