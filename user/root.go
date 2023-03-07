package user

import (
	"image/color"

	// バージョンはv2で揃える！！
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func RunApp() {
	// アプリ起動
	a := app.New()
	w := a.NewWindow("Gorail")

	// baseState := widget.NewLabel("Hello!!")

	// dayCanvas := w.Canvas()
	// aaa := create1DayGraph()
	// dayCanvas.SetContent(aaa)

	// タブで作ろうとしたやつ
	// window1Day := container.NewVBox(
	// 	widget.NewLabel("1Day"),
	// 	create1DayGraph(),
	// )
	// window1Week := container.NewVBox(
	// 	widget.NewLabel("1Week"),
	// )
	// window1Month := container.NewVBox(
	// 	widget.NewLabel("1Month"),
	// )
	// window3Month := container.NewVBox(
	// 	widget.NewLabel("3Month"),
	// )
	tab := container.NewAppTabs(
		container.NewTabItem("1Day", create1DayGraph()),
		container.NewTabItem("1Week", create1WeekGraph()),
		// container.NewTabItem("1Month", window1Month),
		// container.NewTabItem("3Month", window3Month),
	)

	w.SetContent(tab)
	w.Resize(fyne.NewSize(1000, 1000))

	w.ShowAndRun()
}

func create1DayGraph() fyne.CanvasObject {
	// DBから一日分を取得
	line := canvas.NewLine(color.White)
	line.StrokeWidth = 5

	return line
}

func create1WeekGraph() fyne.CanvasObject {
	// DBから一日分を取得
	line := canvas.NewLine(color.Opaque)
	line.StrokeWidth = 5

	return line
}
