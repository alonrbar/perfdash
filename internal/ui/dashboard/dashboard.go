package dashboard

import (
	"log"

	"github.com/alonrbar/perfdash/internal/lib/geo"
	"github.com/alonrbar/perfdash/internal/ui/widgets/cpu"
	"github.com/jroimartin/gocui"
)

// ------------------ //
//   Public types
// ------------------ //

// Dashboard is the main UI element
type Dashboard struct {
	gui       *gocui.Gui
	cpuWidget *cpu.Widget
}

// ------------------ //
//   Private types
// ------------------ //

type layoutFunc func(gui *gocui.Gui) error

// ------------------ //
//    Constructors
// ------------------ //

// New - create new dashboard element
func New() (*Dashboard, error) {
	log.Println("Creating new dashboard")

	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	dash := &Dashboard{
		gui:       gui,
		cpuWidget: cpu.NewWidget(gui),
	}

	return dash, nil
}

// ------------------ //
//   Public methods
// ------------------ //

// Open the dashboard
func (dash *Dashboard) Open() error {

	log.Println("Opening the dashboard")

	// Configure GUI widgets
	gui := dash.gui
	gui.SetManagerFunc(layout(dash))

	// Set key bindings
	err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		dash.Close()
		return gocui.ErrQuit
	})
	if err != nil {
		return err
	}

	// Start widget loops
	dash.cpuWidget.Start(geo.Origin)

	// Start the main UI loop
	err = gui.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

// Close - close the ui
func (dash *Dashboard) Close() {
	dash.gui.Close()
}

// --------------------- //
//   Private functions
// --------------------- //

// Layout handler re-calculates view sizes when the terminal window resizes
func layout(dash *Dashboard) layoutFunc {
	return func(gui *gocui.Gui) error {

		err := dash.cpuWidget.Redraw(geo.Origin)
		if err != nil {
			return err
		}

		gui.Cursor = false

		return nil
	}
}
