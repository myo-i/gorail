package user

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

func RunApp() {
	// アプリ起動
	a := app.New()
	w := a.NewWindow("Gorail")

	baseState := widget.NewLabel("Hello!!")
	w.SetContent(container.NewVBox(
		baseState,
		widget.NewButton("Next", func() {
			baseState.SetText("World!!")
		}),
	))

	w.ShowAndRun()
}
