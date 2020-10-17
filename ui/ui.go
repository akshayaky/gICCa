package ui

import (
	"log"

	"github.com/jroimartin/gocui"
)

var (
	viewArr = []string{"messageBox", "channels", "members"}
	active  = 0
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	if nextIndex == 0 || nextIndex == 3 {
		g.Cursor = true
	} else {
		g.Cursor = false
	}

	active = nextIndex
	return nil
}

func nextChannel(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cursorx, cursory := v.Cursor()
		if err := v.SetCursor(cursorx, cursory+1); err != nil {
			originx, originy := v.Origin()
			if err := v.SetOrigin(originx, originy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func previousChannel(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cursorx, cursory := v.Cursor()
		if cursory != 0 {
			if err := v.SetCursor(cursorx, cursory-1); err != nil {
				originx, originy := v.Origin()
				if err := v.SetOrigin(originx, originy-1); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func getChannel(toName *string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		var channel string
		var err error
		if v != nil {
			_, cursory := v.Cursor()
			if channel, err = v.Line(cursory); err != nil {
				channel = ""
			}
			*toName = channel
			ChangeTitle(toName, "mainView", g)
			setCurrentViewOnTop(g, "messageBox")
			active = 0
			clearView("mainView", g)

		}

		return nil
	}
}

func getText(chan2 chan string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if err := v.SetCursor(0, 0); err != nil {
			return err
		}
		if err := v.SetOrigin(0, 0); err != nil {
			return err
		}
		chan2 <- v.Buffer()
		v.Clear()

		return nil
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("channels", 0, 0, maxX/6-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "channels"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		// v.Wrap = true

	}
	if v, err := g.SetView("members", 5*maxX/6, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "members"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		// v.Wrap = true

	}
	if v, err := g.SetView("mainView", maxX/6, 0, 5*maxX/6-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Messages"

		v.Wrap = true
		v.Autoscroll = true

	}
	if v, err := g.SetView("messageBox", maxX/6, maxY-3, 5*maxX/6-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = ""
		v.Editable = true
		v.Wrap = true

		if _, err := setCurrentViewOnTop(g, "messageBox"); err != nil {
			return err
		}

	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

//ChangeTitle changes the title of a particular view
func ChangeTitle(title *string, viewName string, g *gocui.Gui) {
	if v, err := g.View(viewName); err == nil {
		v.Title = *title
	}
}

func clearView(viewName string, g *gocui.Gui) {
	if v, err := g.View(viewName); err == nil {
		v.Clear()
	}
}

//Control sets the UI of the app
func Control(g *gocui.Gui, chan2 chan string, toName *string) error {

	for {
		g.SetManagerFunc(layout)
		if err := g.SetKeybinding("messageBox", gocui.KeyEnter, gocui.ModNone,
			getText(chan2)); err != nil {
			return err
		}
		if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding("channels", gocui.KeyArrowDown, gocui.ModNone, nextChannel); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding("channels", gocui.KeyArrowUp, gocui.ModNone, previousChannel); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding("channels", gocui.KeyEnter, gocui.ModNone, getChannel(toName)); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
			log.Panicln(err)
		}

		if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
			log.Panicln(err)
		}
	}
	// fmt.Println("Exiting")
	return nil
}
