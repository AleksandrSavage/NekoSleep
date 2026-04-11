package main

import (
	"NekoSleep/internal/ui"
)

func main() {
	
	App := ui.NewApp(
		resourceNunitoRegularTtf,
		resourceNunitoBoldTtf,
		resourceKitteniconIco,
	 	resourceKittengreetPng,
		
	)

	// Запускаем
	App.Run()
}