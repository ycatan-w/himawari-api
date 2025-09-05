package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"github.com/ycatan-w/himawari-api/internal/db"
	"github.com/ycatan-w/himawari-api/internal/output"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	output.PrintBox(fmt.Sprintf("\033[0;32mhimawari-server\033[0m version \033[0;33m%s\033[0m (\033[2;93m%s\033[0m) %s", version, commit, date))

	// Flags
	portFlag := flag.Int("port", 9740, "Port to run the server on")
	initFlag := flag.Bool("init-db", false, "Initialize database and exit")
	flag.Parse()
	port := *portFlag

	// Init DB
	if *initFlag {
		output.PrintHeader("Initialize Database")
		if err := db.InitDB(); err != nil {
			output.PrintFail("Database initialization failure")
			log.Fatalf("DB initialization failed: %v", err)
		}
		output.PrintSuccess("Database initialized")
		return
	}

	// Connect to DB for server execution
	if err := db.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	// Start server
	if runtime.GOOS == "windows" {
		runWindowsService(port)
	} else {
		runUnixService(port)
	}
}
