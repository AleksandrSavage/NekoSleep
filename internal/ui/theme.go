package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// customTheme реализует интерфейс fyne.Theme.
// Сама структура скрыта (с маленькой буквы), мы будем отдавать её через функцию.
type customTheme struct {
	defaultTheme   fyne.Theme
	customFont     fyne.Resource
	customFontBold fyne.Resource
	customIcon     fyne.Resource
}

// NewCustomTheme — конструктор, который принимает сгенерированные ресурсы из main
func NewCustomTheme(font fyne.Resource, fontBold fyne.Resource, icon fyne.Resource) fyne.Theme {
	return &customTheme{
		defaultTheme: theme.DefaultTheme(),
		customFont:   	font,
		customFontBold: fontBold,
		customIcon:   	icon,
	}
}

func (c *customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return c.defaultTheme.Color(name, variant)
}

func (c *customTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Bold { return c.customFontBold }
	return c.customFont
}

func (c *customTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	// Подменяем иконку Home на твоего котенка
	if name == theme.IconNameHome { return c.customIcon }
	return c.defaultTheme.Icon(name)
}

func (c *customTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText { return 20 }
	if name == theme.SizeNameHeadingText { return 60 }
	return c.defaultTheme.Size(name)
}