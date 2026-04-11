package main

import (
	"NekoSleep/internal/ui"
	"NekoSleep/internal/locker"
)

func main() {
	
	App := ui.NewApp(
		resourceNunitoRegularTtf,
		resourceNunitoBoldTtf,
		resourceKitteniconIco,
	 	resourceKittengreetPng,
		
	)
	// Тестируем локер, передавая ему ту же картинку котенка
	locker.Show(resourceKittenasleepPng)
	
	App.Run()
}