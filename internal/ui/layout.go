package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"NekoSleep/internal/config"
	"NekoSleep/internal/monitor"
)

func buildMainLayout(kitten_greet fyne.Resource, w fyne.Window) fyne.CanvasObject {

	mainWrapper := container.NewStack()

	var makeEditScreen func() fyne.CanvasObject
	var makeSavedScreen func(data *config.SleepData) fyne.CanvasObject

	makeSavedScreen = func(data *config.SleepData) fyne.CanvasObject {
		timeText := widget.NewLabel(fmt.Sprintf("Sleep time: %s:%s", data.Hour, data.Minute))
		cyclesText := widget.NewLabel(fmt.Sprintf("Prolongation left: %s", data.Cycles))

		timeText.Alignment = fyne.TextAlignCenter
		cyclesText.Alignment = fyne.TextAlignCenter

		go func() {
            ticker := time.NewTicker(2 * time.Second)
            for range ticker.C {
                freshData, err := config.Load()
                
                if err == nil {
                    fyne.Do(func() {
                        timeText.SetText(fmt.Sprintf("Sleep time: %s:%s", freshData.Hour, freshData.Minute))
                        cyclesText.SetText(fmt.Sprintf("Prolongation left: %s", freshData.Cycles))
                    })
                } else {
                    ticker.Stop()  
                    monitor.Stop() 
                    
                    fyne.Do(func() {
                        mainWrapper.Objects = []fyne.CanvasObject{makeEditScreen()}
                        mainWrapper.Refresh()
                    })
                    
                    return 
                }
            }
        }()

		denyBtn := widget.NewButtonWithIcon("deny", theme.CancelIcon(), func() {
			monitor.Stop()
			config.Delete()
			mainWrapper.Objects = []fyne.CanvasObject{makeEditScreen()}
			mainWrapper.Refresh()
		})

		sizedDenyBtn := container.NewGridWrap(fyne.NewSize(160, 45), denyBtn)

		content := container.NewVBox(
			layout.NewSpacer(),
			timeText,
			cyclesText,
			widget.NewLabel(""),
			container.NewCenter(sizedDenyBtn),
			layout.NewSpacer(),
		)

		return container.NewPadded(content)
	}

	makeEditScreen = func() fyne.CanvasObject {
		helloText := widget.NewRichTextFromMarkdown("# Good night")
		img := canvas.NewImageFromResource(kitten_greet)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(80, 80))
		headerRow := container.NewCenter(container.NewHBox(helloText, img))

		var hours, minutes []string
		for i := 0; i < 24; i++ {
			hours = append(hours, fmt.Sprintf("%02d", i))
		}
		for i := 0; i < 60; i++ {
			minutes = append(minutes, fmt.Sprintf("%02d", i))
		}

		currentHour, currentMinute, currentCycle := "00", "00", "1"

		hourSelect := widget.NewSelect(hours, func(selected string) { currentHour = selected })
		hourSelect.SetSelected("00")

		minuteSelect := widget.NewSelect(minutes, func(selected string) { currentMinute = selected })
		minuteSelect.SetSelected("00")

		questionText := widget.NewLabel("when to sleep?")
		timeSelectionRow := container.NewCenter(container.NewHBox(questionText, hourSelect, minuteSelect))

		calcButton := widget.NewButton("save", func() {
			data := &config.SleepData{
				Hour:   currentHour,
				Minute: currentMinute,
				Cycles: currentCycle,
			}

			err := config.Save(data)
			if err == nil {
				monitor.Start()
				mainWrapper.Objects = []fyne.CanvasObject{makeSavedScreen(data)}
				mainWrapper.Refresh()
			} else {
				fmt.Println("❌ Saving error:", err)
			}
		})

		sizedButton := container.NewGridWrap(fyne.NewSize(160, 45), calcButton)
		buttonRow := container.NewCenter(sizedButton)

		var cycleOptions []string
		for i := 1; i <= 5; i++ {
			cycleOptions = append(cycleOptions, fmt.Sprint(i))
		}

		cycleSelect := widget.NewSelect(cycleOptions, func(selected string) { currentCycle = selected })
		cycleSelect.SetSelected("1")

		smallSelectWrapper := container.NewGridWrap(fyne.NewSize(70, 35), cycleSelect)
		infoIcon := newHoverIcon(theme.InfoIcon(), "How many times can you prolong\n your session by 1 hour", w.Canvas())
		infoWrapper := container.NewGridWrap(fyne.NewSize(24, 24), infoIcon)
		bottomRightRow := container.NewHBox(layout.NewSpacer(), infoWrapper, smallSelectWrapper)

		content := container.NewVBox(
			layout.NewSpacer(),
			headerRow,
			timeSelectionRow,
			widget.NewLabel(""),
			buttonRow,
			layout.NewSpacer(),
			bottomRightRow,
		)

		return container.NewPadded(content)
	}

	data, err := config.Load()
	if err == nil {
		mainWrapper.Objects = []fyne.CanvasObject{makeSavedScreen(data)}
	} else {
		mainWrapper.Objects = []fyne.CanvasObject{makeEditScreen()}
	}

	return mainWrapper
}