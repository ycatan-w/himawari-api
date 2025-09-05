package main

import (
	"log"

	"github.com/kardianos/service"
	"github.com/ycatan-w/himawari-api/internal/output"
	"github.com/ycatan-w/himawari-api/internal/server"
)

// -------------------- Windows service --------------------
func runWindowsService(port int) {
	svcConfig := &service.Config{
		Name:        "HimawariServer",
		DisplayName: "Himawari Server",
		Description: "Himawari API Service on Windows",
	}

	prg := &winProgram{port: port}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}

type winProgram struct {
	port int
}

func (p *winProgram) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *winProgram) run() {
	srv := server.New()
	srv.Port = p.port
	srv.Run()
}

func (p *winProgram) Stop(s service.Service) error {
	output.PrintSubHeader("Stopping Himawari Windows service")
	return nil
}
