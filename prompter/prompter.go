// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prompter

import (
	"fmt"
	"log"
	// "strings"

	"github.com/jroimartin/gocui"
)

const delta = 2
const optionHeight = 5
const helpWidth = 22
const helpHeight = 6
const promptHeight = 8

var (
	curOption = -1
	options   = []string{}
	idxOption = 0
)

func main() {
	var err error

	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)
	if err := initKeybindings(g); err != nil {
		log.Panicln(err)
	}
	if err := newOption(g, promptHeight); err != nil {
		log.Panicln(err)
	}

	err = g.MainLoop()
	if err != nil && err != gocui.Quit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	newPrompt(g)
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

func newPrompt(g *gocui.Gui) error {
	maxX, _ := g.Size()
	v, err := g.SetView("prompt", 0, 0, maxX-helpWidth, promptHeight)
	if err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, "Prompt")
	}
	return nil
}

func initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return newOption(g, idxOption*optionHeight+promptHeight)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return nextOption(g, true)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return nextOption(g, false)
		}); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.Quit
}

func newOption(g *gocui.Gui, offset int) error {
	maxX, _ := g.Size()
	name := fmt.Sprintf("v%v", idxOption)
	v, err := g.SetView(name, 2, offset, maxX/2+optionHeight, offset+optionHeight)
	if err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Wrap = true
		fmt.Fprintln(v, "option", idxOption, " ", len(options))
	}
	if err := g.SetCurrentView(name); err != nil {
		return err
	}
	v.BgColor = gocui.ColorRed

	if curOption >= 0 {
		cv, err := g.View(options[curOption])
		if err != nil {
			return err
		}
		cv.BgColor = g.BgColor
	}

	options = append(options, name)
	curOption = len(options) - 1
	idxOption += 1
	return nil
}

func delOption(g *gocui.Gui) error {
	if len(options) <= 1 {
		return nil
	}

	if err := g.DeleteView(options[curOption]); err != nil {
		return err
	}
	options = append(options[:curOption], options[curOption+1:]...)

	return nextOption(g, false)
}

func nextOption(g *gocui.Gui, upward bool) error {
	if len(options) <= 1 {
		return nil
	}
	next := curOption
	if upward {
		next = curOption + 1
	} else {
		next = curOption - 1
	}
	if next > len(options)-1 {
		next = 0
	} else if next < 0 {
		next = len(options) - 1
	}
	nv, err := g.View(options[next])
	if err != nil {
		return err
	}
	if err := g.SetCurrentView(options[next]); err != nil {
		return err
	}
	nv.BgColor = gocui.ColorRed
	if len(options) > 1 {
		cv, err := g.View(options[curOption])
		if err != nil {
			return err
		}
		cv.BgColor = g.BgColor
	}
	curOption = next
	return nil
}
