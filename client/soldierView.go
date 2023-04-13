package client

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"strconv"
	"time"
)

type soldierView struct {
	Grid        *ui.Grid
	Leader      *widgets.Paragraph
	Victim      *widgets.Paragraph
	Status      *widgets.Paragraph
	DDOSMode    *widgets.Paragraph
	ControlDash *widgets.List
	Banner      *widgets.Paragraph
	Soldiers    *widgets.Table
	Logs        *widgets.List
	TotalSpeed  *widgets.SparklineGroup
}

func newSoldierView() *soldierView {
	var v soldierView
	v.Leader = widgets.NewParagraph()
	v.Leader.Title = "Leader"
	v.Leader.Text = "LEADER"
	v.Leader.BorderStyle.Fg = ui.ColorMagenta
	v.Leader.TitleStyle.Fg = ui.ColorBlue

	v.Victim = widgets.NewParagraph()
	v.Victim.Title = "Victim"
	v.Victim.Text = "#.#.#.#"
	v.Victim.BorderStyle.Fg = ui.ColorMagenta
	v.Victim.TitleStyle.Fg = ui.ColorBlue

	v.Status = widgets.NewParagraph()
	v.Status.Title = "Status"
	v.Status.Text = "NOTHING"
	v.Status.BorderStyle.Fg = ui.ColorMagenta
	v.Status.TitleStyle.Fg = ui.ColorBlue

	v.DDOSMode = widgets.NewParagraph()
	v.DDOSMode.Title = "DDOS Mode"
	v.DDOSMode.Text = "NOTHING"
	v.DDOSMode.BorderStyle.Fg = ui.ColorMagenta
	v.DDOSMode.TitleStyle.Fg = ui.ColorBlue

	v.ControlDash = widgets.NewList()
	v.ControlDash.Title = "CONTROL"
	v.ControlDash.Rows = []string{
		"* EXIT\n",
	}
	v.ControlDash.TextStyle = ui.NewStyle(ui.ColorYellow)
	v.ControlDash.SelectedRowStyle = ui.NewStyle(ui.ColorGreen, ui.ColorClear, ui.ModifierBold)

	v.ControlDash.WrapText = false
	v.ControlDash.PaddingTop = 1
	v.ControlDash.BorderStyle.Fg = ui.ColorMagenta
	v.ControlDash.TitleStyle.Fg = ui.ColorBlue

	v.Banner = widgets.NewParagraph()
	v.Banner.Text = ""
	v.Banner.TextStyle = ui.NewStyle(ui.ColorCyan, ui.ColorClear, ui.ModifierBold)
	v.Banner.BorderStyle.Fg = ui.ColorMagenta
	v.Banner.Border = true

	v.Soldiers = widgets.NewTable()
	v.Soldiers.Title = "Soldiers"
	v.Soldiers.Rows = [][]string{
		{"", "", "", ""},
	}
	v.Soldiers.BorderStyle.Fg = ui.ColorMagenta
	v.Soldiers.TitleStyle.Fg = ui.ColorBlue
	v.Soldiers.RowSeparator = true
	v.Soldiers.TextStyle = ui.NewStyle(ui.ColorCyan, ui.ColorClear, ui.ModifierBold)
	v.Soldiers.PaddingTop = 1
	v.Logs = widgets.NewList()
	v.Logs.Title = "Logs"
	v.Logs.Rows = []string{}
	v.Logs.PaddingTop = 1
	v.Logs.PaddingLeft = 2
	v.Logs.TextStyle = ui.NewStyle(ui.ColorRed)
	v.Logs.SelectedRowStyle = ui.NewStyle(ui.ColorBlack, ui.ColorRed)
	v.Logs.BorderStyle.Fg = ui.ColorMagenta
	v.Logs.TitleStyle.Fg = ui.ColorBlue

	sl0 := widgets.NewSparkline()
	sl0.Data = []float64{0}
	sl0.LineColor = ui.ColorCyan

	v.TotalSpeed = widgets.NewSparklineGroup(sl0)
	v.TotalSpeed.Title = "Total Speed"
	v.TotalSpeed.SetRect(0, 0, 20, 10)
	v.TotalSpeed.BorderStyle.Fg = ui.ColorMagenta
	v.TotalSpeed.TitleStyle.Fg = ui.ColorBlue

	return &v
}

func (v *soldierView) ResetSize() {
	termWidth, termHeight := ui.TerminalDimensions()
	if termWidth > 20 {
		v.Grid.SetRect(0, 0, termWidth, termHeight)
	}
}

func (v *soldierView) Render() {
	ui.Render(v.Grid)
}
func (v *soldierView) Init() {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	v.Grid = ui.NewGrid()
	v.ResetSize()
	v.Grid.Set(
		ui.NewCol(3.0/4,
			ui.NewRow(0.2,
				ui.NewCol(1.0/4, v.Leader),
				ui.NewCol(1.0/4, v.Victim),
				ui.NewCol(1.0/4, v.Status),
				ui.NewCol(1.0/4, v.DDOSMode),
			),
			ui.NewRow(0.8,
				ui.NewCol(1.0/4,
					ui.NewRow(2.0/3, v.ControlDash),
					ui.NewRow(1.0/3, v.Banner),
				),
				ui.NewCol(3.0/4, v.Soldiers),
			),
		),
		ui.NewCol(1.0/4,
			ui.NewRow(1.0/2, v.Logs),
			ui.NewRow(1.0/2, v.TotalSpeed),
		),
	)
	go func() {
		for {
			frames := []string{frame0, frame1, frame2, frame3, frame4}
			for f := range frames {
				v.Banner.Text = frames[f]
				v.Render()
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
}

func (v *soldierView) updateDataForSoldier(data CampAPI) {
	v.Leader.Text = data.Leader.Name
	v.Victim.Text = data.Settings.VictimServer
	v.Status.Text = data.Settings.Status
	v.DDOSMode.Text = data.Settings.DDOSType
	v.Soldiers.Rows = [][]string{}
	for _, soldier := range data.Soldiers {
		v.Soldiers.Rows = append(v.Soldiers.Rows, []string{soldier.Name, soldier.Ip, strconv.Itoa(soldier.Speed) + " req/sec"})
	}
	if len(data.Soldiers) == 0 {
		v.Soldiers.Rows = append(v.Soldiers.Rows, []string{"", "", ""})
	}

	if data.Settings.Status == "stopped" {
		v.TotalSpeed.Sparklines[0].Title = "0 req/sec"

		v.TotalSpeed.Sparklines[0].Data = []float64{0}
	} else {
		if len(v.TotalSpeed.Sparklines[0].Data) >= 100 {
			v.TotalSpeed.Sparklines[0].Data = v.TotalSpeed.Sparklines[0].Data[1:]
		}
		v.TotalSpeed.Sparklines[0].Title = strconv.Itoa(data.TotalSpeed) + " req/sec"
		v.TotalSpeed.Sparklines[0].Data = append(v.TotalSpeed.Sparklines[0].Data, float64(data.TotalSpeed))
	}

}

func (v *soldierView) addLog(log string) {
	if log == "" {
		return
	}
	if len(v.Logs.Rows) > 10 {
		v.Logs.Rows = []string{}
	}
	v.Logs.Rows = append(v.Logs.Rows, log)
}
func StartSoldierView(changedDataChan chan CampAPI, logChan chan string) {
	v := newSoldierView()
	v.Init()
	defer ui.Close()

	v.Render()

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				v.ResetSize()
				v.Render()
			case "<Up>", "k":
				v.ControlDash.ScrollUp()
			case "<Down>", "j":
				v.ControlDash.ScrollDown()

			case "<MouseWheelDown>", "h":
				v.Logs.ScrollDown()
				v.Render()

			case "<MouseWheelUp>", "l":
				v.Logs.ScrollUp()
				v.Render()
			}

		case data := <-changedDataChan:
			go func() {
				v.updateDataForSoldier(data)
				v.Render()
			}()
		case log := <-logChan:
			go func() {
				v.addLog(log)
				v.Render()
			}()
		}
	}
}
