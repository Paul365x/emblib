package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/emblib/adapters/jef"
	"github.com/emblib/adapters/pes_pec"
	"github.com/emblib/adapters/shared"
)

var file string = "2024SewNow.JEF"

func main() {
	var pay *shared.Payload
	var cmds []shared.PCommand

	file_type := strings.ToLower(filepath.Ext(file))
	switch file_type {
	case ".pes":
		pay = pes_pec.Read_pes(file)
	case ".jef":
		pay = jef.Read_jef(file)

	}

	myApp := app.New()
	w := myApp.NewWindow("Lines")
	c_prev := shared.PCommand{
		Command1: shared.Stitch,
		Command2: shared.Stitch,
		Dx:       0,
		Dy:       0,
		Color:    0,
	}

	prev := fyne.NewPos(0, 0)
	img := container.NewWithoutLayout()

	// all stitches are offset from centre
	origin_x := pay.Width / 2.0
	origin_y := pay.Height / 2.0
	cols := pay.Palette
	cmds = pay.Cmds
	col_idx := 0
	i := 0
	for p := range cmds {
		/*
			if p >= 3 {
				fmt.Printf("breaking on command %d\n", p)
				break
			}*/
		if p == 0 && cmds[p].Dx == 0 && cmds[p].Dy == 0 {
			continue
		}
		if p == 1 {
			c_prev = cmds[p]
			prev = fyne.NewPos(cmds[p].Dx+origin_x, cmds[p].Dy+origin_y)
			continue
		}
		x := cmds[p].Dx + c_prev.Dx
		y := cmds[p].Dy + c_prev.Dy
		switch cmds[p].Command1 {
		case shared.ColorChg:
			col_idx = cmds[p].Color - 1 // 1 indexed in file
			fmt.Printf("color chg: %d\n", col_idx)
		case shared.Trim, shared.Jump: // jump without line/thread
			c_prev = cmds[p]
			c_prev.Dx = x
			c_prev.Dy = y
			prev = fyne.NewPos(x+origin_x, y+origin_y)
		default:
			if p > 5 {
				//		break
			}
			c_prev = cmds[p]
			c_prev.Dx = x
			c_prev.Dy = y
			l := canvas.NewLine(cols[col_idx])
			l.Position1 = prev
			prev = fyne.NewPos(x+origin_x, y+origin_y)
			l.Position2 = prev
			l.StrokeWidth = 2
			//if col_idx == 5 {
			img.Add(l)
			//}
		}
		i++
	}

	w.SetContent(img)

	w.Resize(fyne.NewSize(100, 100))
	w.ShowAndRun()

}
