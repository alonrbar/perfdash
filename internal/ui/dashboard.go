package ui

import (
	"log"

	"github.com/alonrbar/perfdash/internal/lib/geo"
	"github.com/alonrbar/perfdash/internal/ui/widgets"
	"github.com/jroimartin/gocui"
)

// Dashboard is the main UI element
type Dashboard struct {
	gui       *gocui.Gui
	cpuWidget *widgets.CPUWidget
}

type layoutFunc func(gui *gocui.Gui) error

// ------------------ //
//    Constructors
// ------------------ //

// NewDashboard - create new dashboard element
func NewDashboard() (*Dashboard, error) {
	log.Println("Creating new dashboard")

	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	cpuWidget, err := widgets.NewCPUWidget(gui, geo.Origin, geo.Point{X: 13, Y: 20})
	if err != nil {
		return nil, err
	}

	dash := &Dashboard{
		gui:       gui,
		cpuWidget: cpuWidget,
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
	gui.SetManagerFunc(dashboardLayout(dash))

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
func dashboardLayout(dash *Dashboard) layoutFunc {
	return func(gui *gocui.Gui) error {
		return nil
	}
}
