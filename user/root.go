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

	rect3 := canvas.NewRectangle(chooseColor(GREEN))

	window1Day := container.NewGridWithColumns(
		2,
		container.NewGridWithColumns(3,
			rect,
			rect2,
			rect3,
		),
		container.NewGridWithColumns(3,
			rect,
			rect2,
			rect3,
		),
	)

	rect2.Resize(fyne.NewSize(50, 125))

	window1Week := Democreate1WeekGraph()

	/////////////     やりたいレイアウトはこれ！！！！   /////////////
	// 後はポジションを動的にする!!!!

	/////////////////////////////////////////////////////////////

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

func setContentToCircle(c fyne.Canvas) {
	red := color.NRGBA{R: 0xff, G: 0x33, B: 0x33, A: 0xff}
	circle := canvas.NewCircle(color.White)
	circle.StrokeWidth = 4
	circle.StrokeColor = red
	circle.Resize(fyne.Size{Height: 10, Width: 10})
	c.SetContent(circle)
}

func setContentToText(c fyne.Canvas) {
	green := color.NRGBA{R: 0, G: 180, B: 0, A: 255}
	text := canvas.NewText("Text", green)
	text.TextStyle.Bold = true
	// text.TextSize = 100
	text.Resize(fyne.NewSize(100, 50))
	c.SetContent(text)
}

func create1DayGraph(c fyne.Canvas) {
	// yellow := color.NRGBA{R: 226, G: 231, B: 17, A: 1}
	rect := canvas.NewRectangle(color.White)
	rect.SetMinSize(fyne.NewSize(100, 100))
	rect.Move(fyne.NewPos(50, 50))
	c.SetContent(rect)

	// line := canvas.NewLine(color.White)
	// line.Position1 = fyne.Position{X: 100, Y: 30}
	// line.Position2 = fyne.Position{X: 300 - 30, Y: 100}
	// line.Resize(fyne.Size{Height: line.Position2.X, Width: line.Position2.Y})
	// // line.Move(line.Position())
	// fmt.Println(line.Position())
	// c.SetContent(line)
}

func Democreate1WeekGraph() fyne.CanvasObject {
	green := chooseColor(GREEN)
	yellow := chooseColor(YELLOW)
	white := chooseColor(WHITE)
	rightBlue := chooseColor(RIGHTBLUE)

	container1 := createBarChart(green, green, "Youtube", 0, 80)

	container2 := createBarChart(green, white, "Udemy", 50, 120)

	container3 := createBarChart(yellow, yellow, "github", 100, 100)

	container4 := createBarChart(rightBlue, rightBlue, "github", 150, 50)

	// border
	labelWeek := widget.NewLabel("1Week")
	labelDay := widget.NewLabel("1Day")
	labelWeek.Move(fyne.NewPos(30, 30))

	window1Week := container.NewWithoutLayout(
		// 同じCanvasObjectを複数個入れても同一のものとみなされる
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

func create1WeekGraph() fyne.CanvasObject {
	// DBから一日分を取得
	line := canvas.NewLine(color.Opaque)
	line.StrokeWidth = 5

	return line
}
