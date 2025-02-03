package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func main() {
	pay := read_pes("Clef.pes")
	cmds := pay.Cmds

	myApp := app.New()
	w := myApp.NewWindow("Lines")
	c_prev := PCommand{
		Command1: Stitch,
		Command2: Stitch,
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
	col_idx := 0
	for p := range cmds {
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
		case ColorChg:
			col_idx = cmds[p].Color - 1 // 1 indexed in file
		case Trim: // jump without line/thread
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
			l.Position1 = prev
			prev = fyne.NewPos(x+origin_x, y+origin_y)
			l.Position2 = prev
			l.StrokeWidth = 2
			img.Add(l)
		}
	}

	w.SetContent(img)

	w.Resize(fyne.NewSize(100, 100))
	w.ShowAndRun()

}
