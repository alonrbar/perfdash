package dashboard

import (
	"log"

	"github.com/alonrbar/perfdash/internal"
	"github.com/alonrbar/perfdash/internal/ui/widgets/cpu"
	"github.com/jroimartin/gocui"
)

// Dashboard is the main UI element
type Dashboard struct {
	gui       *gocui.Gui
	cpuWidget *cpu.Widget
}

//
// Constructors
//

// New - create new dashboard element
func New() (*Dashboard, error) {
	log.Println("Creating new dashboard")

	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalln("Failed to create GUI: ", err)
		return nil, err
	}

	dash := &Dashboard{
		gui:       gui,
		cpuWidget: cpu.NewWidget(gui),
	}

	return dash, nil
}

//
// Public methods
//

// Open the dashboard
func (dash *Dashboard) Open() {

	log.Println("Opening the dashboard")

	// Configure GUI widgets
	gui := dash.gui
	dash.cpuWidget.Redraw(internal.Point{X: 0, Y: 0})

	gui.SetManagerFunc(layout(dash))

	// Set key bindings
	err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		dash.Close()
		return gocui.ErrQuit
	})
	if err != nil {
		log.Println("Could not set key binding: ", err)
		return
	}

	// Start the main UI loop
	err = gui.MainLoop()
	if err != nil {
		if err == gocui.ErrQuit {
			log.Println("Bye")
			return
		}

		log.Fatalln("Failed to start main GUI loop: ", err)
		return
	}
}

// Close - close the ui
func (dash *Dashboard) Close() {
	dash.gui.Close()
}

//
// Private functions
//

// Layout handler re-calculates view sizes when the terminal window resizes
func layout(dash *Dashboard) layoutFunc {
	return func(gui *gocui.Gui) error {

		if err := dash.cpuWidget.Redraw(internal.Point{X: 0, Y: 0}); err != nil {
			return err
		}

		return nil
	}
}

//
// Private types
//

type layoutFunc func(gui *gocui.Gui) error
