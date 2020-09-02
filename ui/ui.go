package ui

import (
	"log"

	"github.com/jroimartin/gocui"
)

var (
	viewArr = []string{"mainView", "channels", "messageBox"}
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

		v.Wrap = true

	}
	if v, err := g.SetView("mainView", maxX/6, 0, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Messages"

		v.Wrap = true
		v.Autoscroll = true

	}
	if v, err := g.SetView("messageBox", maxX/6, maxY-3, maxX-1, maxY-1); err != nil {
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

func ChangeTitle(title string, viewName string, g *gocui.Gui) {
	if v, err := g.View(viewName); err == nil {
		v.Title = title
	}
}

//Control sets the UI of the app
func Control(g *gocui.Gui, chan2 chan string) error {

	for {
		g.SetManagerFunc(layout)
		if err := g.SetKeybinding("messageBox", gocui.KeyEnter, gocui.ModNone,
			getText(chan2)); err != nil {
			return err
		}
		if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
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
