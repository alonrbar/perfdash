package main

import (
	"log"
	"os"
	"os/user"

	"github.com/alonrbar/perfdash/internal/ui"
)

func main() {

	// Configure logs
	logfile := createLogFile()
	defer logfile.Close()
	log.SetOutput(logfile)

	log.Println("App started")

	// Start the UI
	dash, err := ui.NewDashboard()
	if err != nil {
		log.Fatalf("Failed to init the dashboard: %+v\n", err)
	}
	err = dash.Open()
	if err != nil {
		log.Fatalf("Failed to open the dashboard: %+v\n", err)
	}
	log.Println("Bye")
}

func createLogFile() *os.File {

	usr, err := user.Current()
	if err != nil {
		log.Fatalln("Failed to get user info:", err)
	}

	err = os.MkdirAll(usr.HomeDir+"/.perfdash", os.ModeDir)
	if err != nil {
		log.Fatalln("Failed to create logs dir:", err)
	}

	logfile, err := os.OpenFile(usr.HomeDir+"/.perfdash/applog.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}

	return logfile
}
