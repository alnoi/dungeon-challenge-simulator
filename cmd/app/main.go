package main

import (
	"fmt"
	"os"

	"dungeon-event-processor/internal/config"
	"dungeon-event-processor/internal/engine"
	"dungeon-event-processor/internal/output"
	"dungeon-event-processor/internal/parser"
)

func main() {
	appConfig, err := config.Load(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cfg, err := parser.ParseConfig(appConfig.ConfigPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	eventsFile, err := os.Open(appConfig.EventsPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer eventsFile.Close()

	events, err := parser.ParseEvents(eventsFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	processor := engine.New(cfg)
	for _, event := range events {
		processor.Handle(event)
	}

	for _, event := range processor.Logs() {
		fmt.Println(output.FormatEvent(event))
	}

	fmt.Println(output.FormatReport(processor.Report()))
}
