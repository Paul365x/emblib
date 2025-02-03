package adapters

import (
	"image/color"
	"os"
)

var fh *os.File = os.Stdout

// ColorSub stores the color structure
type ColorSub struct {
	CodeLen  uint8
	Code     []byte
	Color    color.Color
	u1       uint8
	ColType  uint32
	DescLen  uint8
	Desc     string
	BrandLen uint8
	Brand    string
	ChartLen uint8
	Chart    string
	count    uint32
}

const (
	Stitch = iota + 1
	Jump
	Trim
	ColorChg
	End
)
