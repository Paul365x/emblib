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

var file string = "ATG12847.jef"

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
	origin_x := pay.Width / 2
	origin_y := pay.Height / 2
	cols := pay.Palette
	cmds = pay.Cmds
	col_idx := 0
	i := 0
	for p := range cmds {
		if p > 3000 {
			//		continue
		}
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
			if col_idx == 6 {
				goto out
			}
			col_idx = cmds[p].Color
			var r, g, b, a = cols[col_idx].RGBA()

			fmt.Printf("color chg: %d: %x %x %x %x\n", col_idx, uint8(r), uint8(g), uint8(b), uint8(a))
		case shared.Trim, shared.Jump: // jump without line/thread
			prev = fyne.NewPos(x+origin_x, y+origin_y)
			//prev = fyne.NewPos(x+prev.X, y+prev.Y)
			c_prev = cmds[p]
			c_prev.Dx = x
			c_prev.Dy = y

		default:
			if p > 5 {
				//		break
			}
			c_prev = cmds[p]
			c_prev.Dx = x
			c_prev.Dy = y
			l := canvas.NewLine(cols[col_idx])
			var r, g, b, a = cols[col_idx].RGBA()

			fmt.Printf("color draw: %d: %x %x %x %x\n", col_idx, uint8(r), uint8(g), uint8(b), uint8(a))
			l.Position1 = prev
			prev = fyne.NewPos(x+origin_x, y+origin_y)
			//prev = fyne.NewPos(x+prev.X, y+prev.Y)
			l.Position2 = prev
			l.StrokeWidth = 2
			//if col_idx == 5 {
			img.Add(l)
			//}
		}
		i++
	}
out:
	w.SetContent(img)

	w.Resize(fyne.NewSize(500, 500))
	w.ShowAndRun()

}
