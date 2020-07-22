package main

import (
	"log"
	"os"
	"os/user"

	"github.com/alonrbar/perfdash/internal/ui/dashboard"
)

func main() {

	// Configure logs
	logfile := createLogFile()
	defer logfile.Close()
	log.SetOutput(logfile)

	log.Println("App started")

	// Start the UI
	dash, err := dashboard.New()
	if err != nil {
		log.Fatalln("Failed to init the dashboard", err)
	}
	err = dash.Open()
	if err != nil {
		log.Fatalln("Failed to open dashboard", err)
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