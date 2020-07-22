package cpu

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/alonrbar/perfdash/internal/lib/geo"
	"github.com/alonrbar/perfdash/internal/lib/num"
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
	values []int
	gui    *gocui.Gui
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
			time.Sleep(time.Second * 2)
		}
	}()
}

// Redraw the CPU widget
func (widget *Widget) Redraw(topLeft geo.Point) error {

	const graphWidth = 10

	gui := widget.gui

	_, termHeight := gui.Size()

	view, err := gui.SetView(ViewName, topLeft.X, topLeft.Y, graphWidth+3, termHeight-ui.MarginBottom)
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

	if cpuVal, err := getCPU(); err != nil {
		log.Panicln("Failed to get cpu", err)
	} else {
		widget.values = append(widget.values, cpuVal)
	}

	printGraph(widget.values, view, graphWidth)

	return nil
}

func printGraph(values []int, view *gocui.View, graphWidth int) {

	max := 0
	for _, val := range values {
		if val > max {
			max = val
		}
	}

	start := num.Max(0, len(values)-graphWidth)

	builder := strings.Builder{}
	for row := max; row > 0; row-- {

		builder.WriteString(fmt.Sprintf("%2v ", row))

		for col := start; col < len(values); col++ {
			if values[col] >= row {
				_, err := builder.WriteString("\u2588")
				if err != nil {
					log.Panicln(err)
				}
			} else {
				_, err := builder.WriteString(" ")
				if err != nil {
					log.Panicln(err)
				}
			}
		}
		_, err := fmt.Fprintln(view, builder.String())
		if err != nil {
			log.Panicln(err)
		}
		builder.Reset()
	}
}

func getCPU() (int, error) {
	buf := bytes.Buffer{}

	cmd := exec.Command("wmic", "cpu", "get", "loadpercentage")
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		return 0, err
	}

	parts := strings.Fields(buf.String())
	if len(parts) < 2 {
		return 0, nil
	}
	cpuStr := strings.TrimSpace(parts[1])
	cpuInt, err := strconv.ParseInt(cpuStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(cpuInt), nil
}
