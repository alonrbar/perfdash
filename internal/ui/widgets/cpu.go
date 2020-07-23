package widgets

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/alonrbar/perfdash/internal/lib/geo"
	"github.com/alonrbar/perfdash/internal/lib/num"
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
)

const cpuWidgetName = "CPU"

// ------------------ //
//   Public types
// ------------------ //

// CPUWidget is the UI widget for the CPU meter
type CPUWidget struct {
	gui         *gocui.Gui
	view        *gocui.View
	topLeft     geo.Point
	bottomRight geo.Point
	cpuValues   []int
}

// ------------------ //
//    Constructors
// ------------------ //

// NewCPUWidget creates a new CPU widget
func NewCPUWidget(gui *gocui.Gui, topLeft geo.Point, bottomRight geo.Point) (*CPUWidget, error) {

	widget := &CPUWidget{
		gui:         gui,
		topLeft:     topLeft,
		bottomRight: bottomRight,
	}

	gui.Update(func(*gocui.Gui) error {

		view, err := widget.setView()
		if err != nil {
			return err
		}
		widget.view = view

		view.Title = "CPU"
		view.Highlight = true
		view.FgColor = gocui.ColorCyan
		view.SelFgColor = gocui.ColorCyan

		return nil
	})

	log.Println("New CPU widget created")

	return widget, nil
}

// ------------------ //
//   Public methods
// ------------------ //

// Start the widget redraw loop
func (widget *CPUWidget) Start(topLeft geo.Point) {
	go func() {
		for {
			widget.addCPUSample()
			widget.gui.Update(func(g *gocui.Gui) error {
				return widget.Redraw()
			})
			time.Sleep(time.Second / 2)
		}
	}()
}

// Resize the CPU widget and return it's view
func (widget *CPUWidget) Resize(topLeft geo.Point, bottomRight geo.Point) (*gocui.View, error) {
	widget.topLeft = topLeft
	widget.bottomRight = bottomRight
	return widget.setView()
}

// Redraw the CPU widget content.
// Notice: This method must be called inside a gui.Update call.
func (widget *CPUWidget) Redraw() error {

	widget.view.Clear()

	widgetWidth := widget.bottomRight.X - widget.topLeft.X
	graphWidth := widgetWidth - 3
	maxValue := num.Max(widget.cpuValues...)
	startIndex := num.Max(0, len(widget.cpuValues)-graphWidth)

	log.Printf("Redrawing CPU widget. Width: %d", widgetWidth)

	// Draw the CPU graph line by line, from top to bottom
	builder := strings.Builder{}
	for row := maxValue; row > 0; row-- {

		// Y axis label
		_, err := builder.WriteString(fmt.Sprintf("%2v ", row))
		if err != nil {
			return errors.Wrap(err, "failed to write Y axis label")
		}

		// Build current graph row
		for col := startIndex; col < len(widget.cpuValues); col++ {
			var char string
			if widget.cpuValues[col] >= row {
				char = "\u2588"
			} else {
				char = " "
			}
			_, err := builder.WriteString(char)
			if err != nil {
				return errors.Wrap(err, "failed to write graph line to builder")
			}
		}

		// Emit row to the screen
		_, err = fmt.Fprintln(widget.view, builder.String())
		if err != nil {
			return errors.Wrap(err, "failed to write graph line to screen")
		}
		builder.Reset()
	}

	return nil
}

// ------------------ //
//   Private methods
// ------------------ //

// setView is a convenient wrapper for gocui.Gui.SetView
func (widget *CPUWidget) setView() (*gocui.View, error) {
	view, err := widget.gui.SetView(cpuWidgetName, widget.topLeft.X, widget.topLeft.Y, widget.bottomRight.X, widget.bottomRight.Y)
	if err != nil && err != gocui.ErrUnknownView {
		// ErrUnknownView is not a real error condition.
		// It just says that the view did not exist before and needs initialization.
		return nil, errors.Wrap(err, "failed to set view")
	}
	return view, nil
}

func (widget *CPUWidget) addCPUSample() error {
	cpuVal, err := getCPU()
	if err != nil {
		return errors.Wrap(err, "failed to get CPU usage data")
	}
	widget.cpuValues = append(widget.cpuValues, cpuVal)
	return nil
}

type cpuQueryResult struct {
	PercentProcessorTime int
}

func getCPU() (int, error) {

	query := `
		SELECT PercentProcessorTime 
		FROM Win32_PerfFormattedData_PerfOS_Processor
		WHERE Name = "_Total"
	`
	var queryResult []cpuQueryResult
	err := wmi.Query(query, &queryResult)
	if err != nil {
		return 0, errors.Wrap(err, "WMI CPU query failed")
	}
	if len(queryResult) != 1 {
		return 0, errors.Errorf("invalid query result length: %v", queryResult)
	}

	return queryResult[0].PercentProcessorTime, nil
}
