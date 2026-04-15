package locker

import (
	"fmt"
	"image/color"
	"os/exec"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"NekoSleep/internal/config"
)


// ==============================
//  КАСТОМНАЯ КНОПКА ВЫКЛЮЧЕНИЯ
// ==============================
type customButton struct {
	widget.BaseWidget
	bg         *canvas.Rectangle
	content    *fyne.Container
	baseColor  color.Color
	hoverColor color.Color
	onTapped   func()
}

func newCustomButton(text string, iconRes fyne.Resource, action func()) *customButton {
	b := &customButton{
		baseColor:  color.RGBA{R: 130, G: 200, B: 130, A: 255}, 
		hoverColor: color.RGBA{R: 100, G: 170, B: 100, A: 255}, 
		onTapped:   action,
	}

	b.bg = canvas.NewRectangle(b.baseColor)
	b.bg.CornerRadius = 6 

	lbl := canvas.NewText(text, color.Black)
	lbl.TextStyle.Bold = true
	lbl.Alignment = fyne.TextAlignCenter

	ico := canvas.NewImageFromResource(iconRes)
	ico.FillMode = canvas.ImageFillContain
	ico.SetMinSize(fyne.NewSize(20, 20))

	b.content = container.NewCenter(container.NewHBox(ico, lbl))

	b.ExtendBaseWidget(b)
	return b
}

func (b *customButton) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(b.bg, b.content))
}

func (b *customButton) Tapped(_ *fyne.PointEvent) {
	if b.onTapped != nil {
		b.onTapped()
	}
}

func (b *customButton) TappedSecondary(_ *fyne.PointEvent) {}

func (b *customButton) MouseIn(_ *desktop.MouseEvent) {
	b.bg.FillColor = b.hoverColor
	b.bg.Refresh()
}

func (b *customButton) MouseMoved(_ *desktop.MouseEvent) {}

func (b *customButton) MouseOut() {
	b.bg.FillColor = b.baseColor
	b.bg.Refresh()
}


// ==========================
// ЛОГИКА ЭКРАНА БЛОКИРОВКИ
// ==========================
var focusTicker *time.Ticker

func Show(sleepImg fyne.Resource, onProlong func()) {
	a := fyne.CurrentApp()
	if a == nil {
		a = app.New()
	}

	w := a.NewWindow("NekoSleep - SleepinTime")

	w.SetFullScreen(true)
	w.SetCloseIntercept(func() {})

	bgColor := color.RGBA{R: 26, G: 26, B: 46, A: 255}
	bg := canvas.NewRectangle(bgColor)

	img := canvas.NewImageFromResource(sleepImg)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(250, 250))

	sleepTitle := widget.NewLabel("Shhh, the kitten is already asleep.")
	sleepTitle.TextStyle = fyne.TextStyle{Bold: true}

	sleepSubtitle := widget.NewLabel("You're tired. Take a break.")

	data, err := config.Load()
	if err != nil { 
		return 
	}

	cyclesLeft, _ := strconv.Atoi(data.Cycles)

	info := widget.NewLabel(fmt.Sprintf("Prolongation left: %d", cyclesLeft))

	var unlockBtn *widget.Button

    unlockBtn = widget.NewButton("I need a minute...", func() {

        if cyclesLeft > 0 {

            cyclesLeft--

            data.Cycles = strconv.Itoa(cyclesLeft)

            config.Save(data)

			if focusTicker != nil {
                focusTicker.Stop()
            }
			
			StopKeyboardBlocker()

            if onProlong != nil { onProlong()}

            w.Close()  
        }
    })

	if cyclesLeft <= 0 {
		unlockBtn.Disable()
		unlockBtn.SetText("Prolongations are over 😿")
	}

	shutdownBtn := newCustomButton("Shut down PC", theme.LogoutIcon(), func() {	
		exec.Command("shutdown", "/s", "/t", "0").Run()
	})

	sizedUnlockBtn := container.NewGridWrap(fyne.NewSize(250, 40), unlockBtn)
	sizedShutdownBtn := container.NewGridWrap(fyne.NewSize(250, 40), shutdownBtn)

	contentBox := container.NewVBox(
		container.NewCenter(img),
		container.NewCenter(sleepTitle),
		container.NewCenter(sleepSubtitle),

		widget.NewLabel(""), 

		container.NewCenter(info),
		container.NewCenter(sizedUnlockBtn),
		
		widget.NewLabel(""), 
		
		container.NewCenter(sizedShutdownBtn),
	)

	content := container.NewCenter(contentBox)
	w.SetContent(container.NewStack(bg, content))

	w.Show()

	w.RequestFocus() 

    focusTicker = time.NewTicker(1 * time.Second)
    go func() {
        for range focusTicker.C {
            w.RequestFocus()
        }
    }()

	go StartKeyboardBlocker()
}