package main

import "github.com/ycatan-w/himawari-api/internal/server"

// -------------------- Linux/macOS --------------------
func runUnixService(port int) {
	srv := server.New()
	srv.Port = port
	srv.Run()
}
