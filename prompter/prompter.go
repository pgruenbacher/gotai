// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prompter

import (
	"fmt"
	// "log"
	// "strings"

	"github.com/jroimartin/gocui"
)

const delta = 2
const optionHeight = 5
const helpWidth = 22
const helpHeight = 6
const promptHeight = 8
const promptName = "prompt"

const orderListName = "orderList"
const orderListHeight = 16
const orderListWidth = 22

func layout(g *gocui.Gui) error {
	newPrompt(g)
	newOrderList(g)
	maxX, _ := g.Size()
	v, err := g.SetView("legend", maxX-helpWidth, 0, maxX, helpHeight)
	if err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, "KEYBINDINGS")
		fmt.Fprintln(v, "Space: Enter")
		fmt.Fprintln(v, "Backspace: Back")
		fmt.Fprintln(v, "← ↑ → ↓: Move View")
		fmt.Fprintln(v, "ctrl-C: Exit")
	}
	return nil
}

// func editView(text string, id string, g *gocui.Gui) error {
// 	view, err := g.View(id)
// 	if err != nil {
// 		return err
// 	}
// 	view.Clear()
// 	view.Wrap = true
// 	fmt.Fprintln(view, text)
// 	return nil
// }

func newPrompt(g *gocui.Gui) error {
	maxX, _ := g.Size()
	v, err := g.SetView(orderListName, maxX-orderListWidth, helpHeight, maxX, helpHeight+orderListHeight)
	if err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, "This is the prompt")
	}
	return nil
}

func newOrderList(g *gocui.Gui) error {
	maxX, _ := g.Size()
	v, err := g.SetView(promptName, 0, 0, maxX-helpWidth, promptHeight)
	if err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, "This is the order list")
	}
	return nil
}

func (self *Prompt) initKeybindings() error {
	if err := self.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := self.g.SetKeybinding("", gocui.KeySpace, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return self.newPromptandOptions()
		}); err != nil {
		return err
	}
	if err := self.g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return self.selectOption()
		}); err != nil {
		return err
	}
	if err := self.g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return self.nextOption(true)
		}); err != nil {
		return err
	}
	if err := self.g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return self.nextOption(false)
		}); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.Quit
}

func (self *Prompt) newOption(o validOption) error {
	opt := o.Option()
	offset := len(self.options)*optionHeight + promptHeight
	maxX, _ := self.g.Size()
	v, err := self.g.SetView(opt.id, 2, offset, maxX/2+optionHeight, offset+optionHeight)
	if err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Wrap = true
		fmt.Fprintln(v, opt.txt)
	}
	// if err := self.g.SetCurrentView(opt.id); err != nil {
	// 	return err
	// }
	v.BgColor = gocui.ColorRed

	if self.curOption >= 0 {
		cv, err := self.g.View(self.options[self.curOption].id)
		if err != nil {
			return err
		}
		cv.BgColor = self.g.BgColor
	}

	self.options = append(self.options, opt)
	self.curOption = len(self.options) - 1
	return nil
}

func (self *Prompt) selectOption() error {
	if len(self.options) <= 0 {
		fmt.Println("self.options <= 1")
		return nil
	}
	self.pushOrderList(self.options[self.curOption])
	err := self.updateOrderList()
	if err != nil {
		return err
	}
	return nil
}

func (self *Prompt) pushOrderList(opt option) error {
	self.orderList = append(self.orderList, opt)
	return nil
}

func (self *Prompt) updateOrderList() error {
	v, err := self.g.View(orderListName)
	if err != nil {
		return err
	}
	v.Clear()
	for _, order := range self.orderList {
		fmt.Fprintln(v, order.txt)
	}
	return nil
}

func (self *Prompt) delOptions() error {
	if len(self.options) <= 0 {
		return nil
	}

	for _, option := range self.options {
		if err := self.g.DeleteView(option.id); err != nil {
			return err
		}
	}
	self.options = make([]option, 0)
	self.curOption = -1
	return nil
}

func (self *Prompt) delOrderList() error {
	self.orderList = make([]option, 0)
	v, err := self.g.View(orderListName)
	if err != nil {
		return err
	}
	v.Clear()
	return nil
}

func (self *Prompt) nextOption(upward bool) error {
	if len(self.options) <= 1 {
		return nil
	}
	next := self.curOption
	if upward {
		next = self.curOption + 1
	} else {
		next = self.curOption - 1
	}
	if next > len(self.options)-1 {
		next = 0
	} else if next < 0 {
		next = len(self.options) - 1
	}
	nv, err := self.g.View(self.options[next].id)
	if err != nil {
		return err
	}
	// if err := self.g.SetCurrentView(self.options[next].id); err != nil {
	// 	return err
	// }
	nv.BgColor = gocui.ColorRed
	if len(self.options) > 1 {
		cv, err := self.g.View(self.options[self.curOption].id)
		if err != nil {
			return err
		}
		cv.BgColor = self.g.BgColor
	}
	self.curOption = next
	return nil
}
