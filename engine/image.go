package engine

import (
	"image/color"

	"path/filepath"
	"strings"

	"github.com/fogleman/gg"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

type ImgComposer struct {
	px       float32
	py       float32
	img      *gg.Context
	ox       float32
	oy       float32
	name     string
	path     string
	filetype RenderType
}

func NewJpgComposer() *ImgComposer {
	return &ImgComposer{
		px:       0.0,
		py:       0.0,
		img:      nil,
		ox:       0.0,
		oy:       0.0,
		name:     "",
		filetype: Jpg,
		path:     "",
	}
}

func NewPngComposer() *ImgComposer {
	return &ImgComposer{
		px:       0.0,
		py:       0.0,
		img:      nil,
		ox:       0.0,
		oy:       0.0,
		name:     "",
		filetype: Png,
		path:     "",
	}
}

func (c *ImgComposer) Setup(ox, oy float32, name string) {
	c.ox = ox
	c.oy = oy
	c.px = ox
	c.py = oy
	c.img = gg.NewContext(int(3.0*ox), int(3.0*oy))
	c.img.SetColor(color.White)
	c.img.Clear()
	c.name = filepath.Base(name)
	c.path = strings.Trim(name, c.name)
	c.name = strings.TrimSuffix(c.name, filepath.Ext(c.name))
}

func (c *ImgComposer) SetPos(x, y float32) {
	c.px = x + c.ox
	c.py = y + c.oy
}

func (c *ImgComposer) Line(ex, ey float32, col color.Color) {
	x2 := c.ox + ex
	y2 := c.oy + ey
	c.img.SetColor(col)
	c.img.SetLineWidth(2.0)
	c.img.DrawLine(float64(c.px), float64(c.py), float64(x2), float64(y2))
	c.img.Stroke()
	c.px = x2
	c.py = y2

}

func (c *ImgComposer) Get() any {
	return c.img.Image()
}

func (c *ImgComposer) Display() {
	a := app.New()
	w := a.NewWindow(c.name)
	//container = container.NewWithoutLayout()
	content := canvas.NewImageFromImage(c.img.Image())
	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
