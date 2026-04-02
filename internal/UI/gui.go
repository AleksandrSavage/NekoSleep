package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Run() {
	// 1. Создание приложения
    a := app.New()
    // 2. Создание окна
    w := a.NewWindow("NekoHib")
	w.Resize(fyne.NewSize(500, 400))
	w.SetFixedSize(true)
    // 3. Создание виджетов
    label := widget.NewLabel("Привет, Fyne!")
    button := widget.NewButton("Нажми меня", func() {
        label.SetText("Кнопка нажата")
    })

    // 4. Размещение виджетов в контейнере (вертикальный список)
    content := container.NewVBox(
        label,
        button,
    )

    // 5. Установка содержимого и запуск
    w.SetContent(content)
    w.ShowAndRun()
}