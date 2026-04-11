package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"

	"NekoSleep/internal/config"
)

func buildMainLayout(kitten_greet fyne.Resource, w fyne.Window) fyne.CanvasObject {

	// --- 1 ---
	helloText := widget.NewRichTextFromMarkdown("# Good night")
	
	img := canvas.NewImageFromResource(kitten_greet)
	img.FillMode = canvas.ImageFillContain

	img.SetMinSize(fyne.NewSize(80, 80)) 
	
	headerRow := container.NewCenter(container.NewHBox(helloText, img))


	// --- 2 ---
	var hours, minutes []string
	for i := 0; i < 24; i++ {
		hours = append(hours, fmt.Sprintf("%02d", i))
	}
	for i := 0; i < 60; i++ {
		minutes = append(minutes, fmt.Sprintf("%02d", i))
	}


	// --- 3 ---
	currentHour, currentMinute, currentCycle := "00", "00", "1"


	// --- 4 ---
	hourSelect := widget.NewSelect(hours, func(selected string) {
		currentHour = selected
	})
	hourSelect.SetSelected("00")

	minuteSelect := widget.NewSelect(minutes, func(selected string) {
		currentMinute = selected
	})
	minuteSelect.SetSelected("00")

	questionText := widget.NewLabel("when to sleep?")
	timeSelectionRow := container.NewCenter(container.NewHBox(questionText, hourSelect, minuteSelect))


	// --- 5 ---
	calcButton := widget.NewButton("save", func() {
		
		data := &config.SleepData{
			Hour:   currentHour,
			Minute: currentMinute,
			Cycles: currentCycle,
		}

		err := config.Save(data)
		
		if err != nil {
			fmt.Println("❌ Ошибка сохранения:", err)
		} else {
			fmt.Println("✅ Настройки успешно сохранены в config.json!")
		}
	})

	sizedButton := container.NewGridWrap(fyne.NewSize(160, 45), calcButton)
    buttonRow := container.NewCenter(sizedButton)
	

	// --- 6 ---
	var cycleOptions []string
	for i := 1; i <= 5; i++ {
		cycleOptions = append(cycleOptions, fmt.Sprint(i))
	}
	
	cycleSelect := widget.NewSelect(cycleOptions, func(selected string) {
		currentCycle = selected
	})
	cycleSelect.SetSelected("1") 

	smallSelectWrapper := container.NewGridWrap(fyne.NewSize(70, 35), cycleSelect)
	
	infoIcon := newHoverIcon(
        theme.InfoIcon(), 
        "How many times you can\n unlock your screen.", 
        w.Canvas(),
    )
    
    infoWrapper := container.NewGridWrap(fyne.NewSize(24, 24), infoIcon)

	bottomRightRow := container.NewHBox(layout.NewSpacer(),infoWrapper, smallSelectWrapper)


	// --- 7 ---
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