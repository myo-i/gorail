package user

import (
	"image/color"

	// バージョンはv2で揃える！！
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	// Color
	RED = iota
	GREEN
	YELLOW
	RIGHTBLUE
	WHITE

	// Element Size
	WINDOW_WIDTH     = 800
	WINDOW_HEIGHT    = 800
	TAB_WIDTH        = 500
	TAB_HEIGHT       = 500
	CONTAINER_WIDTH  = 30
	CONTAINER_HEIGHT = 150
	CONTAINER_POS_X  = 100
	CONTAINER_POS_Y  = 100
	TEXT_HEIGHT      = 30
)

type bar struct {
	size   fyne.Size
	canvas fyne.CanvasObject
	window fyne.Window
}

func RunApp() {
	// アプリ起動
	myApp := app.New()
	myWindow := myApp.NewWindow("Canvas")
	// myCanvas := myWindow.Canvas()

	// blue := color.NRGBA{R: 0, G: 0, B: 180, A: 255}
	// rect := canvas.NewRectangle(blue)
	// create1DayGraph(myCanvas)

	rect := canvas.NewRectangle(chooseColor(WHITE))
	rect.Resize(fyne.NewSize(100, 100))
	rect.Move(fyne.NewPos(50, 50))

	rect2 := canvas.NewRectangle(chooseColor(RED))

	// windowやappのサイズが変更されたら動的にサイズを取得する！！
	rect2.Resize(fyne.NewSize(50, 125))

	window1Day := create1DayGraph()
	window1Week := create1WeekGraph()

	tab := container.NewAppTabs(
		container.NewTabItem("1Day", window1Day),
		container.NewTabItem("1Week", window1Week),
		// container.NewTabItem("1Month", window1Month),
		// container.NewTabItem("3Month", window3Month),
	)

	// w.Resize(fyne.NewSize(1000, 1000))

	// w.ShowAndRun()
	tab.Resize(fyne.NewSize(TAB_WIDTH, TAB_HEIGHT))
	myWindow.Resize(fyne.NewSize(WINDOW_WIDTH, WINDOW_WIDTH))
	myWindow.SetContent(tab)
	myWindow.ShowAndRun()

	/////////  ここのコード読む！！！！！！！！！！！！！！
	// https://github.com/fyne-io/examples/
	// https://github.com/fyne-io/calculator/
}

func create1DayGraph() fyne.CanvasObject {
	green := chooseColor(GREEN)
	yellow := chooseColor(YELLOW)
	white := chooseColor(WHITE)
	rightBlue := chooseColor(RIGHTBLUE)

	container2 := createBarChart(white, green, "Google", 50, 200)

	container3 := createBarChart(white, yellow, "Brave", 100, 150)

	container4 := createBarChart(white, rightBlue, "FireFox", 150, 50)

	// border
	labelWeek := widget.NewLabel("1Week")
	labelDay := widget.NewLabel("1Day")
	labelWeek.Move(fyne.NewPos(30, 30))

	window1Day := container.NewWithoutLayout(
		labelWeek,
		labelDay,
		container2,
		container3,
		container4,
	)

	return window1Day
}

func create1WeekGraph() fyne.CanvasObject {
	green := chooseColor(GREEN)
	yellow := chooseColor(YELLOW)
	white := chooseColor(WHITE)
	rightBlue := chooseColor(RIGHTBLUE)

	container1 := createBarChart(white, green, "Youtube", 0, 80)

	container2 := createBarChart(white, green, "Udemy", 50, 120)

	container3 := createBarChart(white, yellow, "github", 100, 100)

	container4 := createBarChart(white, rightBlue, "github", 150, 50)

	// border
	labelWeek := widget.NewLabel("1Week")
	labelDay := widget.NewLabel("1Day")
	labelWeek.Move(fyne.NewPos(30, 30))

	window1Week := container.NewWithoutLayout(
		labelWeek,
		labelDay,
		container1,
		container2,
		container3,
		container4,
	)

	return window1Week

}

func chooseColor(colors int) color.Color {
	switch colors {
	case 0:
		// Red
		return color.NRGBA{R: 0xff, G: 0x33, B: 0x33, A: 0xff}
	case 1:
		// Green
		return color.NRGBA{R: 0x43, G: 0xff, B: 0x64, A: 0xd9}
	case 2:
		// Yellow
		return color.NRGBA{R: 0xf2, G: 0xff, B: 0x00, A: 0xff}
	case 3:
		// Right Blue
		return color.NRGBA{R: 0x00, G: 0xbb, B: 0xff, A: 0xff}
	case 4:
		return color.White
	}
	return color.White
}

// 棒グラフと項目名をコンテナとして作成
func createBarChart(textColor, rectColor color.Color, textContent string, duration, barHeight float32) *fyne.Container {
	rect := canvas.NewRectangle(rectColor)

	text := canvas.NewText(textContent, textColor)

	containerBarText := container.NewWithoutLayout(
		rect,
		text,
	)

	containerBarText.Resize(fyne.NewSize(CONTAINER_WIDTH, CONTAINER_HEIGHT))
	rect.Resize(fyne.NewSize(containerBarText.Size().Width, barHeight))
	text.Resize(fyne.NewSize(containerBarText.Size().Width, TEXT_HEIGHT))

	containerBarText.Move(fyne.NewPos(CONTAINER_POS_X+duration, CONTAINER_POS_Y))
	rect.Move(fyne.NewPos(containerBarText.Position().X, containerBarText.Size().Height-rect.Size().Height))
	text.Move(fyne.NewPos(containerBarText.Position().X, containerBarText.Size().Height))

	return containerBarText
}
