package prompter

import (
	"github.com/jroimartin/gocui"
	"github.com/pgruenbacher/gotai/state"
	"log"
)

type validOption interface {
	Display() option
}

type option struct {
	id       string
	txt      string
	order    interface{}
	limitKey limitedOpt
}

type limitedOpt string

const (
	armyMovementOrder limitedOpt = "armyMovementOrder"
)

type Prompt struct {
	g         *gocui.Gui
	options   []option
	curOption int
	orderList map[limitedOpt]option
	game      *state.State
	marshall  *Marshall
}

func (self *Prompt) InitiatePrompt(game *state.State) {

	self.game = game
	self.curOption = -1
	marshall := Marshall{
		ArmyManager: game.ArmiesManager,
	}
	self.marshall = &marshall

	var err error
	self.g = gocui.NewGui()
	if err := self.g.Init(); err != nil {
		log.Panicln(err)
	}
	defer self.g.Close()
	// create the layout of the prompt, help, an
	self.g.SetLayout(layout)
	// initiate key bindings
	if err := self.initKeybindings(); err != nil {
		log.Panicln(err)
	}
	// if err := newOption(g, promptHeight); err != nil {
	// 	log.Panicln(err)
	// }

	err = self.g.MainLoop()
	if err != nil && err != gocui.Quit {
		log.Panicln(err)
	}

}

// func (self *Prompt) EditPrompt(text string) {
// 	editView(text, promptName, self.g)
// }

// func (self *Prompt) NewOptions(opts []option)

// func (self *Prompt) EditOption(opt option) {

// 	editView(text)
// }

func (self *Prompt) newPromptandOptions() (err error) {
	orders := make([]interface{}, len(self.orderList))

	for _, order := range self.orderList {
		orders = append(orders, order.order)
	}
	err = self.game.Next(orders)
	if err != nil {
		log.Panicln(err)
	}
	err = self.delOptions()
	if err != nil {
		return err
	}
	err = self.delOrderList()
	if err != nil {
		return err
	}
	opts := self.marshall.GiveArmyOptions()
	for _, opt := range opts {
		err = self.newOption(opt)
		if err != nil {
			return err
		}
	}
	return nil
}
