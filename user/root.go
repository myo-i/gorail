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
	// Color
	RED = iota
	ORANGERED
	ORANGE
	YELLOW
	LIME
	SPRINGGREEN
	AQUA
	DODGERBLUE
	BLUE
	LIGHTSTEELBLUE
	WHITE

	// Element Size
	WINDOW_WIDTH     = 800
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
	red            = chooseColor(RED)
	orangered      = chooseColor(ORANGERED)
	orange         = chooseColor(ORANGE)
	yellow         = chooseColor(YELLOW)
	lime           = chooseColor(LIME)
	springgreen    = chooseColor(SPRINGGREEN)
	aqua           = chooseColor(AQUA)
	dodgerblue     = chooseColor(DODGERBLUE)
	blue           = chooseColor(BLUE)
	lightsteelblue = chooseColor(LIGHTSTEELBLUE)
	white          = chooseColor(WHITE)
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

	rect := canvas.NewRectangle(chooseColor(WHITE))
	rect.Resize(fyne.NewSize(100, 100))
	rect.Move(fyne.NewPos(50, 50))

	rect2 := canvas.NewRectangle(chooseColor(RED))

	rect2.Resize(fyne.NewSize(50, 125))

	window1Month := create1MonthGraph(topTenKey, topTenValue)
	window1Week := create1WeekGraph()

	tab := container.NewAppTabs(
		container.NewTabItem("1Month", window1Month),
		container.NewTabItem("1Week", window1Week),
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

	container1 := createBarChart(white, lightsteelblue, strconv.Itoa(topTenValue[9]), topTenKey[9], 0, float32(topTenValue[9]))
	container2 := createBarChart(white, blue, strconv.Itoa(topTenValue[8]), topTenKey[8], 50, float32(topTenValue[8]))
	container3 := createBarChart(white, dodgerblue, strconv.Itoa(topTenValue[7]), topTenKey[7], 100, float32(topTenValue[7]))
	container4 := createBarChart(white, aqua, strconv.Itoa(topTenValue[6]), topTenKey[6], 150, float32(topTenValue[6]))
	container5 := createBarChart(white, springgreen, strconv.Itoa(topTenValue[5]), topTenKey[5], 200, float32(topTenValue[5]))
	container6 := createBarChart(white, lime, strconv.Itoa(topTenValue[4]), topTenKey[4], 250, float32(topTenValue[4]))
	container7 := createBarChart(white, yellow, strconv.Itoa(topTenValue[3]), topTenKey[3], 300, float32(topTenValue[3]))
	container8 := createBarChart(white, orange, strconv.Itoa(topTenValue[2]), topTenKey[2], 350, float32(topTenValue[2]))
	container9 := createBarChart(white, orangered, strconv.Itoa(topTenValue[1]), topTenKey[1], 400, float32(topTenValue[1]))
	container10 := createBarChart(white, red, strconv.Itoa(topTenValue[0]), topTenKey[0], 450, float32(topTenValue[0]))

	window1Month := container.NewWithoutLayout(
		container1,
		container2,
		container3,
		container4,
		container5,
		container6,
		container7,
		container8,
		container9,
		container10,
	)

	return window1Month
}

func create1WeekGraph() fyne.CanvasObject {

	container1 := createBarChart(white, lime, "80", "Youtube", 0, 80)

	container2 := createBarChart(white, yellow, "120", "Udemy", 50, 120)

	container3 := createBarChart(white, orange, "100", "github", 100, 100)

	container4 := createBarChart(white, orangered, "50", "github", 150, 50)

	container5 := createBarChart(white, red, "160", "github", 200, 160)

	window1Week := container.NewWithoutLayout(
		container1,
		container2,
		container3,
		container4,
		container5,
	)

	return window1Week

}

func chooseColor(colors int) color.Color {
	switch colors {
	case 0:
		// Red
		return color.NRGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
	case 1:
		// orangered
		return color.NRGBA{R: 0xff, G: 0x45, B: 0x00, A: 0xff}
	case 2:
		// orange
		return color.NRGBA{R: 0xff, G: 0xa5, B: 0x00, A: 0xff}
	case 3:
		// yellow
		return color.NRGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff}
	case 4:
		// lime
		return color.NRGBA{R: 0x00, G: 0xff, B: 0x00, A: 0xff}
	case 5:
		// springgreen
		return color.NRGBA{R: 0x00, G: 0xff, B: 0x7f, A: 0xff}
	case 6:
		// aqua
		return color.NRGBA{R: 0x00, G: 0xff, B: 0xff, A: 0xff}
	case 7:
		// dodgerblue
		return color.NRGBA{R: 0x1e, G: 0x90, B: 0xff, A: 0xff}
	case 8:
		// blue
		return color.NRGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff}
	case 9:
		// lightsteelblue
		return color.NRGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff}
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
