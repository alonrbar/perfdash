package cpu

import (
	"fmt"

	"github.com/alonrbar/perfdash/internal"
	"github.com/alonrbar/perfdash/internal/ui"
	"github.com/jroimartin/gocui"
)

// ViewName = "CPU"
const ViewName = "CPU"

// Widget is the UI widget for the CPU meter
type Widget struct {
	gui *gocui.Gui
}

// NewWidget - Create new CPU widget
func NewWidget(gui *gocui.Gui) *Widget {
	return &Widget{
		gui,
	}
}

// Redraw the CPU widget
func (widget *Widget) Redraw(topLeft internal.Point) error {

	gui := widget.gui

	termWidth, termHeight := gui.Size()

	view, err := gui.SetView(ViewName, topLeft.X, topLeft.Y, termWidth/2, termHeight-ui.MarginBottom)
	if err != nil && err != gocui.ErrUnknownView {
		// ErrUnknownView is not a real error condition.
		// It just says that the view did not exist before and needs initialization.
		return err
	}

	view.Clear()

	view.Title = ViewName
	view.Highlight = true
	view.BgColor = gocui.ColorBlack
	view.FgColor = gocui.ColorCyan

	_, err = fmt.Fprintln(view, termWidth, termHeight)
	if err != nil {
		return err
	}

	return nil
}
