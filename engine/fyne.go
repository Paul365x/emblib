package engine

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type FyneComposer struct {
	app    fyne.App
	window fyne.Window
	prev   fyne.Position
	img    *fyne.Container
	ox     float32
	oy     float32
}

func NewFyneComposer() *FyneComposer {
	return &FyneComposer{
		app:    nil,
		window: nil,
		prev:   fyne.NewPos(0.0, 0.0),
		img:    nil,
		ox:     0.0,
		oy:     0.0,
	}
}

func (c *FyneComposer) Setup(ox, oy float32) {
	c.ox = ox
	c.oy = oy
	c.prev = fyne.NewPos(c.ox, c.oy)
	c.img = container.NewWithoutLayout()
}

func (c *FyneComposer) SetPos(x, y float32) {
	c.prev = fyne.NewPos(x+c.ox, y+c.oy)
}

func (c *FyneComposer) Line(ex, ey float32, col color.Color) {
	l := canvas.NewLine(col)
	l.Position1 = c.prev
	c.prev = fyne.NewPos(ex+c.ox, ey+c.oy)
	l.Position2 = c.prev
	l.StrokeWidth = 2
	c.img.Add(l)
}

func (c *FyneComposer) Get() any {
	return c.img
}

func (c *FyneComposer) Display() {
	c.app = app.New()
	c.window = c.app.NewWindow("Lines")
	c.window.SetContent(c.img)
	c.window.Resize(fyne.NewSize(800, 600))
	c.window.ShowAndRun()
}
