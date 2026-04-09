package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// App — наша главная структура, которая хранит состояние приложения
type App struct {
	fyneApp fyne.App
	window  fyne.Window
}

// NewApp — конструктор нашего приложения
func NewApp(font fyne.Resource, icon fyne.Resource, kitten_greet fyne.Resource) *App {
	// 1. Создаем базовое приложение Fyne
	a := app.New()

	// 2. Устанавливаем кастомную тему (функция NewCustomTheme лежит в theme.go)
	a.Settings().SetTheme(NewCustomTheme(font, icon))

	// 3. Создаем главное окно
	w := a.NewWindow("NekoSleep")
	w.Resize(fyne.NewSize(640, 400))
	w.SetFixedSize(true)

	// 4. Получаем интерфейс из layout.go и вставляем его в окно
	content := buildMainLayout(kitten_greet)
	w.SetContent(content)

	return &App{
		fyneApp: a,
		window:  w,
	}
}

// Run запускает главный цикл приложения
func (a *App) Run() {
	a.window.ShowAndRun()
}