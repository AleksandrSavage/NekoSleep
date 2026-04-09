package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"NekoSleep/internal/config"
)

func buildMainLayout(kitten_greet fyne.Resource) fyne.CanvasObject {

	// --- 1. Шапка: Текст + Картинка (Мы восстановили размер котенка 40x40) ---
	helloText := widget.NewRichTextFromMarkdown("# Good night")
	
	img := canvas.NewImageFromResource(kitten_greet)
	img.FillMode = canvas.ImageFillContain
	// Устанавливаем размер, как у исходника. Чтобы текст встал вровень,
	// тебе нужно увеличить шрифт через настройки темы в theme.go.
	img.SetMinSize(fyne.NewSize(80, 80)) 
	
	headerRow := container.NewCenter(container.NewHBox(helloText, img))


	// --- 2. Данные для списков времени ---
	var hours, minutes []string
	for i := 0; i < 24; i++ {
		hours = append(hours, fmt.Sprintf("%02d", i))
	}
	for i := 0; i < 60; i++ {
		minutes = append(minutes, fmt.Sprintf("%02d", i))
	}


	// --- 3. Текст результата (Сделаем его жирным и по центру) ---
	resultLabel := widget.NewLabel("Выбранное время: 00:00")
	resultLabel.TextStyle = fyne.TextStyle{Bold: true}
	resultLabel.Alignment = fyne.TextAlignCenter 

	currentHour, currentMinute, currentCycle := "00", "00", "1"
	updateTimeDisplay := func() {
		resultLabel.SetText(fmt.Sprintf("Выбранное время: %s:%s", currentHour, currentMinute))
	}


	// --- 4. Списки времени (УБРАЛИ СЕПАРАТОРЫ И ДВОЕТОЧИЕ) ---
	hourSelect := widget.NewSelect(hours, func(selected string) {
		currentHour = selected
		updateTimeDisplay()
	})
	hourSelect.SetSelected("00")

	minuteSelect := widget.NewSelect(minutes, func(selected string) {
		currentMinute = selected
		updateTimeDisplay()
	})
	minuteSelect.SetSelected("00")

	questionText := widget.NewLabel("when to sleep?")
	timeSelectionRow := container.NewCenter(container.NewHBox(questionText, hourSelect, minuteSelect))


	// --- 5. КНОПКА (По центру) ---
	calcButton := widget.NewButton("Рассчитать", func() {
		
		// 1. Собираем выбранные данные в структуру из пакета config
		data := &config.SleepData{
			Hour:   currentHour,
			Minute: currentMinute,
			Cycles: currentCycle,
		}

		// 2. Отправляем на сохранение
		err := config.Save(data)
		
		// 3. Проверяем результат
		if err != nil {
			fmt.Println("❌ Ошибка сохранения:", err)
		} else {
			fmt.Println("✅ Настройки успешно сохранены в config.json!")
		}
	})
	buttonRow := container.NewCenter(calcButton)


	// --- 6. СПИСОК (От 1 до 5 в правом нижнем углу) ---
	var cycleOptions []string
	for i := 1; i <= 5; i++ {
		cycleOptions = append(cycleOptions, fmt.Sprint(i))
	}
	
	cycleSelect := widget.NewSelect(cycleOptions, func(selected string) {
		currentCycle = selected
	})
	cycleSelect.SetSelected("1") 

	// ВОТ ОНА, МАГИЯ:
	// Оборачиваем список в контейнер с жестко заданным размером.
	// 70 - это ширина, 35 - высота (можешь поиграть с этими цифрами!)
	smallSelectWrapper := container.NewGridWrap(fyne.NewSize(70, 35), cycleSelect)
	
	// И теперь кладем в правый нижний угол нашу "коробку", а не сам список
	bottomRightRow := container.NewHBox(layout.NewSpacer(), smallSelectWrapper)


	// --- 7. ГЛАВНАЯ СБОРКА (УБРАЛИ СЕПАРАТОРЫ) ---
	content := container.NewVBox(
		layout.NewSpacer(),      // 1. Верхняя пружина (толкает всё вниз)
		
		headerRow,               // 2. Шапка
		timeSelectionRow,        // 3. Вопрос со списками 
		resultLabel,             // 4. Текст результата 
		buttonRow,               // 5. Кнопка расчета
		
		layout.NewSpacer(),      // 6. Нижняя пружина (вместе с верхней держит блок 2-5 по центру)
		
		bottomRightRow,          // 7. Так как это стоит ПОСЛЕ нижней пружины, оно прилипнет к самому низу окна
	)

	// Добавляем отступы по краям окна, чтобы правый нижний список не прилипал вплотную к рамке
	return container.NewPadded(content)
}