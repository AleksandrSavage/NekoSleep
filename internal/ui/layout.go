package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Обрати внимание: мы добавили аргумент helloImage (ресурс твоей картинки)
func buildMainLayout(kitten_greet fyne.Resource) fyne.CanvasObject {

	// --- 1. Шапка: Текст + Картинка ---
	helloText := widget.NewLabel("Привет!")
	
	// Создаем виджет картинки из ресурса
	img := canvas.NewImageFromResource(kitten_greet)
	img.FillMode = canvas.ImageFillContain // Сохраняем пропорции
	img.SetMinSize(fyne.NewSize(40, 40))   // Обязательно задаем размер, иначе в HBox она исчезнет
	
	// Выстраиваем их в ряд по горизонтали
	headerHBox := container.NewHBox(helloText, img)


	// --- 2. Генерация данных для списков (00..23 и 00..59) ---
	var hours []string
	for i := 0; i < 24; i++ {
		hours = append(hours, fmt.Sprintf("%02d", i)) // %02d делает из "5" -> "05"
	}

	var minutes []string
	for i := 0; i < 60; i++ {
		minutes = append(minutes, fmt.Sprintf("%02d", i))
	}


	// --- 3. Динамический текст ---
	// Создаем текст с дефолтным значением
	resultLabel := widget.NewLabel("Выбранное время: 00:00")
	// Можно сделать его жирным
	resultLabel.TextStyle = fyne.TextStyle{Bold: true} 

	// Переменные для хранения текущего выбора
	currentHour := "00"
	currentMinute := "00"

	// Функция, которая склеивает часы и минуты и обновляет Label
	updateTimeDisplay := func() {
		resultLabel.SetText(fmt.Sprintf("Выбранное время: %s:%s", currentHour, currentMinute))
	}


	// --- 4. Выпадающие списки (Select) ---
	hourSelect := widget.NewSelect(hours, func(selected string) {
		currentHour = selected // Запоминаем выбор
		updateTimeDisplay()    // Обновляем текст на экране
	})
	hourSelect.SetSelected("00") // Ставим значение по умолчанию

	minuteSelect := widget.NewSelect(minutes, func(selected string) {
		currentMinute = selected
		updateTimeDisplay()
	})
	minuteSelect.SetSelected("00")


	// --- 5. Строка вопроса: Текст + Списки ---
	questionText := widget.NewLabel("Когда спать?")
	// Ставим двоеточие между часами и минутами просто для красоты
	colon := widget.NewLabel(":") 
	timeSelectionHBox := container.NewHBox(questionText, hourSelect, colon, minuteSelect)


	// --- 6. Собираем всё в финальный вертикальный столбец ---
	content := container.NewVBox(
		headerHBox,             // Шапка
		widget.NewSeparator(),  // Горизонтальная линия-разделитель (для эстетики)
		timeSelectionHBox,      // Вопрос и списки
		widget.NewSeparator(), 
		resultLabel,            // Динамический результат
	)

	// Добавляем отступы по краям, чтобы не прилипало к окну
	paddedContent := container.NewPadded(content)

	return paddedContent
}