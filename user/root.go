package user

import (
	"image/color"
	"strconv"

	// バージョンはv2で揃える！！
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

const (
	// Element Size
	WINDOW_WIDTH     = 1200
	WINDOW_HEIGHT    = 800
	TAB_WIDTH        = 500
	TAB_HEIGHT       = 500
	CONTAINER_WIDTH  = 30
	CONTAINER_HEIGHT = 300
	CONTAINER_POS_X  = 100
	CONTAINER_POS_Y  = 300
	TEXT_HEIGHT      = 30
)

var (
	Red            = color.NRGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
	Orangered      = color.NRGBA{R: 0xff, G: 0x45, B: 0x00, A: 0xff}
	Orange         = color.NRGBA{R: 0xff, G: 0xa5, B: 0x00, A: 0xff}
	Yellow         = color.NRGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff}
	Lime           = color.NRGBA{R: 0x00, G: 0xff, B: 0x00, A: 0xff}
	Springgreen    = color.NRGBA{R: 0x00, G: 0xff, B: 0x7f, A: 0xff}
	Aqua           = color.NRGBA{R: 0x00, G: 0xff, B: 0xff, A: 0xff}
	Dodgerblue     = color.NRGBA{R: 0x1e, G: 0x90, B: 0xff, A: 0xff}
	Blue           = color.NRGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff}
	Lightsteelblue = color.NRGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff}
	White          = color.White
)

type bar struct {
	size   fyne.Size
	canvas fyne.CanvasObject
	window fyne.Window
}

func RunApp(topTenKey []string, topTenValue []int) {
	// アプリ起動
	myApp := app.New()
	myWindow := myApp.NewWindow("Canvas")

	rect := canvas.NewRectangle(chooseColor(10))
	rect.Resize(fyne.NewSize(100, 100))
	rect.Move(fyne.NewPos(50, 50))

	rect2 := canvas.NewRectangle(chooseColor(0))

	rect2.Resize(fyne.NewSize(50, 125))

	window1Month := create1MonthGraph(topTenKey, topTenValue)
	// window1Week := create1WeekGraph()

	tab := container.NewAppTabs(
		container.NewTabItem("1Month", window1Month),
		// container.NewTabItem("1Week", window1Week),
		// container.NewTabItem("1Month", window1Month),
		// container.NewTabItem("3Month", window3Month),
	)

	tab.Resize(fyne.NewSize(TAB_WIDTH, TAB_HEIGHT))
	myWindow.Resize(fyne.NewSize(WINDOW_WIDTH, WINDOW_WIDTH))
	myWindow.SetContent(tab)
	myWindow.ShowAndRun()

	/////////  ここのコード読む！！！！！！！！！！！！！！
	// https://github.com/fyne-io/examples/
	// https://github.com/fyne-io/calculator/
}

func create1MonthGraph(topTenKey []string, topTenValue []int) fyne.CanvasObject {
	const barDistance = 50
	keyLength := len(topTenKey)

	containers := make([]fyne.CanvasObject, keyLength)
	for i := 0; i < keyLength; i++ {
		containers[i] = createBarChart(
			chooseColor(10),
			chooseColor(keyLength-(i+1)),
			strconv.Itoa(topTenValue[keyLength-(i+1)]),
			topTenKey[keyLength-(i+1)],
			float32(barDistance*i),
			float32(topTenValue[keyLength-(i+1)]),
		)
	}

	window1Month := container.NewWithoutLayout(containers...)

	return window1Month
}

// func create1WeekGraph() fyne.CanvasObject {

// 	container1 := createBarChart(chooseColor(10), chooseColor(4), "80", "Youtube", 0, 80)

// 	container2 := createBarChart(chooseColor(10), chooseColor(3), "120", "Udemy", 50, 120)

// 	container3 := createBarChart(chooseColor(10), chooseColor(2), "100", "github", 100, 100)

// 	container4 := createBarChart(chooseColor(10), chooseColor(1), "50", "github", 150, 50)

// 	container5 := createBarChart(chooseColor(10), chooseColor(0), "160", "github", 200, 160)

// 	window1Week := container.NewWithoutLayout(
// 		container1,
// 		container2,
// 		container3,
// 		container4,
// 		container5,
// 	)

// 	return window1Week

// }

func chooseColor(colors int) color.Color {
	switch colors {
	case 0:
		return Red
	case 1:
		return Orangered
	case 2:
		return Orange
	case 3:
		return Yellow
	case 4:
		return Lime
	case 5:
		return Springgreen
	case 6:
		return Aqua
	case 7:
		return Dodgerblue
	case 8:
		return Blue
	case 9:
		return Lightsteelblue
	default:
		return color.White
	}
}

// 棒グラフと項目名をコンテナとして作成
func createBarChart(textColor, rectColor color.Color, ontime, textContent string, duration, barHeight float32) *fyne.Container {
	time := canvas.NewText(ontime+"h", textColor)

	rect := canvas.NewRectangle(rectColor)

	text := canvas.NewText(textContent, textColor)

	containerBarText := container.NewWithoutLayout(
		time,
		rect,
		text,
	)

	containerBarText.Resize(fyne.NewSize(CONTAINER_WIDTH, CONTAINER_HEIGHT))
	time.Resize(fyne.NewSize(containerBarText.Size().Width, TEXT_HEIGHT))
	rect.Resize(fyne.NewSize(containerBarText.Size().Width, barHeight))
	text.Resize(fyne.NewSize(containerBarText.Size().Width, TEXT_HEIGHT))

	containerBarText.Move(fyne.NewPos(CONTAINER_POS_X+duration, CONTAINER_POS_Y))
	time.Move(fyne.NewPos(containerBarText.Position().X, containerBarText.Size().Height-(TEXT_HEIGHT+rect.Size().Height)))
	rect.Move(fyne.NewPos(containerBarText.Position().X, containerBarText.Size().Height-rect.Size().Height))
	text.Move(fyne.NewPos(containerBarText.Position().X, containerBarText.Size().Height))

	return containerBarText
}
