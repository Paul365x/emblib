package engine

import (
	"image/color"

	"github.com/emblib/adapters/shared"
)

/*
** Interface composer
 */

/*
** Setup: does whatever is needed to setup this dingle
** Line: draws a line between two points, takes x1,y1,x2,y2 and a color.Color
** Get: returns the object we are working with - image, container etc
** Display: displays the image on a screen - jpg, fyne img etc
 */

type Composer interface {
	Setup(ox, oy float32)
	SetPos(ox, oy float32)
	Line(ex, ey float32, c color.Color)
	Get() any
	Display()
}

/*
** Engine code
 */

type RenderType int

const (
	Fyne RenderType = iota + 1
)

type Engine struct {
	RType RenderType
	Pay   *shared.Payload
	Comp  Composer
}

func NewEngine() *Engine {
	return &Engine{
		RType: 0,
		Pay:   nil,
		Comp:  nil,
	}
}

func (e *Engine) Setup(t RenderType, p *shared.Payload) {
	e.RType = t
	e.Pay = p
	switch t {
	case Fyne:
		e.Comp = NewFyneComposer()
	}
}

func (e *Engine) Run() {

	c_prev := shared.PCommand{
		Command1: shared.Stitch,
		Command2: shared.Stitch,
		Dx:       0,
		Dy:       0,
		Color:    0,
	}

	cols := e.Pay.Palette
	cmds := e.Pay.Cmds
	col_idx := 0
	i := 0
	comp := e.Comp
	comp.Setup(e.Pay.Width/2, e.Pay.Height/2) // all stitches are offset from centre

	for p := range cmds {
		if p == 0 && cmds[p].Dx == 0 && cmds[p].Dy == 0 {
			continue
		}
		if p == 1 {
			c_prev = cmds[p]
			comp.SetPos(cmds[p].Dx, cmds[p].Dy)
			continue
		}
		x := cmds[p].Dx + c_prev.Dx
		y := cmds[p].Dy + c_prev.Dy
		switch cmds[p].Command1 {
		case shared.ColorChg:
			col_idx = cmds[p].Color
		case shared.Trim, shared.Jump: // jump without line/thread
			comp.SetPos(x, y)
			c_prev = cmds[p]
			c_prev.Dx = x
			c_prev.Dy = y
		default:
			c_prev = cmds[p]
			c_prev.Dx = x
			c_prev.Dy = y
			comp.Line(x, y, cols[col_idx])
		}
		i++
	}
}

func (e *Engine) Get() any {
	return e.Comp.Get()
}

func (e *Engine) Display() {
	e.Comp.Display()
}
