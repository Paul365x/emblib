package engine

import (
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
	"github.com/skratchdot/open-golang/open"
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
	c.name = filepath.Base(name)
	c.path = strings.Trim(name, c.name)
	c.name = strings.TrimSuffix(c.name, filepath.Ext(c.name))
}

func (c *ImgComposer) SetPos(x, y float32) {
	c.px = x
	c.py = y
}

func (c *ImgComposer) Line(ex, ey float32, col color.Color) {
	x2 := c.ox + ex
	y2 := c.oy + ey
	c.img.SetColor(col)
	c.img.DrawLine(float64(c.px), float64(c.py), float64(x2), float64(y2))
	c.px = x2
	c.py = y2

}

func (c *ImgComposer) Get() any {
	return c.img.Image()
}

func (c *ImgComposer) Display() {
	var tmp string
	switch c.filetype {
	case Jpg:
		tmp = filepath.Join(os.TempDir(), c.name+".jpg")
		gg.SaveJPG(tmp, c.img.Image(), 255)
	case Png:
		tmp = filepath.Join(os.TempDir(), c.name+".png")
		gg.SavePNG(tmp, c.img.Image())
	}
	open.Run(tmp)
	os.Remove(tmp)
}
