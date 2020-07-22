package cpu

import (
	"fmt"
	"time"

	"github.com/alonrbar/perfdash/internal/lib/geo"
	"github.com/alonrbar/perfdash/internal/ui"
	"github.com/jroimartin/gocui"
)

// -------------------- //
//   Public constants
// -------------------- //

// ViewName = "CPU"
const ViewName = "CPU"

// ------------------ //
//   Public types
// ------------------ //

// Widget is the UI widget for the CPU meter
type Widget struct {
	caption string
	gui     *gocui.Gui
}

// ------------------ //
//    Constructors
// ------------------ //

// NewWidget - Create new CPU widget
func NewWidget(gui *gocui.Gui) *Widget {
	return &Widget{
		gui: gui,
	}
}

// ------------------ //
//   Public methods
// ------------------ //

// Start the widget redraw loop
func (widget *Widget) Start(topLeft geo.Point) {
	go func() {
		for {
			widget.gui.Update(func(g *gocui.Gui) error {
				// updates to the UI must happen inside a gui.Update method
				return widget.Redraw(topLeft)
			})
			time.Sleep(time.Second)
		}
	}()
}

// Redraw the CPU widget
func (widget *Widget) Redraw(topLeft geo.Point) error {

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
	view.FgColor = gocui.ColorCyan
	view.SelFgColor = gocui.ColorCyan

	widget.caption = time.Now().Format(time.RFC3339)
	_, err = fmt.Fprintln(view, termWidth, termHeight, widget.caption)
	if err != nil {
		return err
	}

	return nil
}
