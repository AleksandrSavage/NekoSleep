package main

import (
	"NekoSleep/internal/ui"
	"NekoSleep/internal/monitor"
	"NekoSleep/internal/config"
)

func main() {
	
	App := ui.NewApp(
		resourceNunitoRegularTtf,
		resourceNunitoBoldTtf,
		resourceKitteniconIco,
	 	resourceKittengreetPng,
		
	)
	monitor.Init(resourceKittenasleepPng)

	if _, err := config.Load(); err == nil {
		monitor.Start()
	}
	
	App.Run()
}