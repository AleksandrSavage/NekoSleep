package main

import (
	"NekoSleep/internal/ui"
)

func main() {
	
	App := ui.NewApp(resourceNunitoRegularTtf, resourceKitteniconIco, resourceKittengreetPng)

	// Запускаем
	App.Run()
}