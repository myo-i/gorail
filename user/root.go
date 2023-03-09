package user

import (
	"fmt"
	"image/color"

	// バージョンはv2で揃える！！
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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

	rect := canvas.NewRectangle(color.White)
	rect.Resize(fyne.NewSize(100, 100))
	rect.Move(fyne.NewPos(50, 50))

	fmt.Println("2")

	red := color.NRGBA{R: 0xff, G: 0x33, B: 0x33, A: 0xff}
	rect2 := canvas.NewRectangle(red)

	green := color.NRGBA{R: 0x43, G: 0xff, B: 0x64, A: 0xd9}
	rect3 := canvas.NewRectangle(green)

	// time.AfterFunc(time.Second*3, func() {
	// 	rect2.Resize(fyne.NewSize(100, 100))
	// 	rect2.Move(fyne.NewPos(50, 50))
	// 	myWindow.Canvas().Refresh(rect2)
	// 	time.AfterFunc(time.Second*2, func() {
	// 		rect2.Resize(fyne.NewSize(125, 125))
	// 		rect2.Move(fyne.NewPos(45, 45))
	// 		myWindow.Canvas().Refresh(rect2)
	// 		fmt.Println("3")
	// 	})
	// })
	fmt.Println("4")

	// タブで作ろうとしたやつ
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

	hour := &canvas.Line{StrokeColor: color.White, StrokeWidth: 5}
	minute := &canvas.Line{StrokeColor: color.White, StrokeWidth: 3}
	second := &canvas.Line{StrokeColor: color.White, StrokeWidth: 1}

	lineWhite := canvas.NewLine(color.White)
	lineWhite.Position1.X = float32(50)
	lineWhite.Position1.Y = float32(50)
	lineWhite.Position2.X = float32(50)
	lineWhite.Position2.Y = float32(50)

	/////////////     やりたいレイアウトはこれ！！！！   /////////////
	// 後はポジションを動的にする!!!!
	labelWeek := widget.NewLabel("1Week")
	labelDay := widget.NewLabel("1Day")
	labelWeek.Move(fyne.NewPos(30, 30))

	rect4 := canvas.NewRectangle(green)
	rect4.Resize(fyne.NewSize(20, 20))
	rect4.Move(fyne.NewPos(100, 100))

	window1Week := container.NewWithoutLayout(
		// 同じCanvasObjectを複数個入れても同一のものとみなされる
		labelWeek,
		labelDay,
		hour,
		minute,
		second,
		lineWhite,
		rect4,
	)
	/////////////////////////////////////////////////////////////
	window1Week.Resize(fyne.NewSize(450, 300))
	rect2.Resize(fyne.NewSize(window1Week.Size().Width-50, window1Week.Size().Height-50))

	tab := container.NewAppTabs(
		container.NewTabItem("1Day", window1Day),
		container.NewTabItem("1Week", window1Week),
		// container.NewTabItem("1Month", window1Month),
		// container.NewTabItem("3Month", window3Month),
	)

	// w.Resize(fyne.NewSize(1000, 1000))

	// w.ShowAndRun()
	tab.Resize(fyne.NewSize(400, 400))
	myWindow.SetContent(tab)
	myWindow.Resize(fyne.NewSize(500, 500))
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

func create1WeekGraph() fyne.CanvasObject {
	// DBから一日分を取得
	line := canvas.NewLine(color.Opaque)
	line.StrokeWidth = 5

	return line
}
