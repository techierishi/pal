package tui

import (
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
)

func Modal() {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := setKeybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("modal", maxX/4, maxY/4, 3*maxX/4, 3*maxY/4, 1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Modal Title"
		v.Frame = true

		if _, err := g.SetView("textbox", maxX/4+1, maxY/4+1, 3*maxX/4-1, maxY/2-1, 1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			// You can customize the textbox here
			v.Editable = true
			v.Wrap = true
		}

		if _, err := g.SetView("button", maxX/2-5, 3*maxY/4-3, maxX/2+5, 3*maxY/4-1, 1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.FgColor = gocui.ColorWhite
			v.BgColor = gocui.ColorGreen
			fmt.Fprintln(v, "OK")
		}
	}
	return nil
}

func setKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("button", gocui.KeyEnter, gocui.ModNone, submit); err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func submit(g *gocui.Gui, v *gocui.View) error {
	textView, err := g.View("textbox")
	if err != nil {
		return err
	}
	text := textView.Buffer()

	// Handle the submitted text here

	// For now, just print the text
	fmt.Println("Submitted Text:", text)

	return quit(g, v)
}
