package shared

import (
	"fmt"
	"image/color"
	"os"
)

var fh *os.File = os.Stdout

// ColorSub stores the color structure
type ColorSub struct {
	CodeLen  uint8
	Code     []byte
	Color    color.Color
	U1       uint8
	ColType  uint32
	DescLen  uint8
	Desc     string
	BrandLen uint8
	Brand    string
	ChartLen uint8
	Chart    string
	Count    uint32
}

// Dump writes out this Struct
func (p ColorSub) Dump() {
	fmt.Printf("\t\tColorSub:\n")
	fmt.Printf("\t\t\tCodeLen: %d 0x%X\n", p.CodeLen, p.CodeLen)
	fmt.Printf("\t\t\tCode: 0x%X\n", p.Code)
	rgba := color.RGBAModel.Convert(p.Color).(color.RGBA)
	fmt.Printf("\t\t\tColor: (%d, %d, %d, %d)\n", rgba.R, rgba.G, rgba.B, rgba.A)
	fmt.Printf("\t\t\tu1: %d 0x%X\n", p.U1, p.U1)
	fmt.Printf("\t\t\tDescLen: %d 0x%X\n", p.DescLen, p.DescLen)
	fmt.Printf("\t\t\tDesc: %s\n", p.Desc)
	fmt.Printf("\t\t\tBrandLen: %d 0x%X\n", p.BrandLen, p.BrandLen)
	fmt.Printf("\t\t\tBrand: %s\n", p.Brand)
	fmt.Printf("\t\t\tChartLen: %d 0x%X\n", p.ChartLen, p.ChartLen)
	fmt.Printf("\t\t\tChart: %s\n", p.Chart)
	fmt.Printf("\t\t\tcount: %d 0x%X\n", p.Count, p.Count)
}

const (
	Stitch = iota + 1
	Jump
	Trim
	ColorChg
	End
)
